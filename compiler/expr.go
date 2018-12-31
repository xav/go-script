// Copyright © 2018 Xavier Basty <xbasty@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package compiler

import (
	"fmt"
	"math/big"

	"github.com/xav/go-script/builtins"
	"github.com/xav/go-script/types"
	"github.com/xav/go-script/values"
	"github.com/xav/go-script/vm"
	"github.com/xav/horus/warden/vm/value"
)

// Expr is the result of compiling an expression.
// It stores the type of the expression and its evaluator function.
type Expr struct {
	*ExprInfo

	// The type of the expression
	ExprType vm.Type

	// Evaluate this node as the given type.
	eval interface{}

	// Map index expressions permit special forms of assignment, for which we need to know the Map and key.
	evalMapValue func(t *vm.Thread) (values.Map, interface{})

	// Evaluate to the "address" of this value; that is, the settable Value object.
	// nil for expressions whose address cannot be taken.
	evalAddr func(t *vm.Thread) vm.Value

	// Execute this expression as a statement.
	// Only expressions that are valid expression statements should set this.
	Exec vm.CodeInstruction

	// If this expression is a type, this is its compiled type.
	// This is only permitted in the function position of a call expression. In this case, exprType should be nil.
	valType vm.Type

	// A short string describing this expression for error messages.
	desc string
}

// resolveIdeal converts the value of the analyzed expression x, which must be a constant ideal number,
// to a new analyzed expression with a constant value of type t.
func (x *Expr) resolveIdeal(t vm.Type) *Expr {
	if !x.ExprType.IsIdeal() {
		logger.Panic().
			Str("type", fmt.Sprintf("%v", x.ExprType)).
			Msgf("attempted to convert from %v, expected ideal", x.ExprType)
	}

	var rat *big.Rat

	// It is an error to assign a value with a non-zero fractional part to an integer,
	// or if the assignment would overflow or underflow, or in general if the value
	// cannot be represented by the type of the variable.

	switch x.ExprType {
	case IdealFloatType:
		rat = x.asIdealFloat()()
		if t.IsInteger() && !rat.IsInt() {
			x.error("constant %v truncated to integer", rat.FloatString(6))
			return nil
		}
	case IdealIntType:
		i := x.asIdealInt()()
		rat = new(big.Rat).SetInt(i)
	default:
		logger.Panic().
			Str("type", fmt.Sprintf("%v", x.ExprType)).
			Msgf("unexpected ideal type %v", x.ExprType)
	}

	// Check bounds
	if t, ok := t.Lit().(vm.BoundedType); ok {
		if rat.Cmp(t.MinVal()) < 0 {
			x.error("constant %v underflows %v", rat.FloatString(6), t)
			return nil
		}
		if rat.Cmp(t.MaxVal()) > 0 {
			x.error("constant %v overflows %v", rat.FloatString(6), t)
			return nil
		}
	}

	// Convert rat to type t.
	res := x.newExpr(t, x.desc)
	switch t := t.Lit().(type) {
	case *types.UintType:
		n, d := rat.Num(), rat.Denom()
		f := new(big.Int).Quo(n, d)
		f = f.Abs(f)
		v := uint64(f.Int64())
		res.eval = func(*vm.Thread) uint64 { return v }

	case *types.IntType:
		n, d := rat.Num(), rat.Denom()
		f := new(big.Int).Quo(n, d)
		v := f.Int64()
		res.eval = func(*vm.Thread) int64 { return v }

	case *types.IdealIntType:
		n, d := rat.Num(), rat.Denom()
		f := new(big.Int).Quo(n, d)
		res.eval = func() *big.Int { return f }

	case *types.FloatType:
		n, d := rat.Num(), rat.Denom()
		v := float64(n.Int64()) / float64(d.Int64())
		res.eval = func(*vm.Thread) float64 { return v }

	case *types.IdealFloatType:
		res.eval = func() *big.Rat { return rat }

	default:
		logger.Panic().Msgf("cannot convert to type %T", t)
	}

	return res
}

// convertToInt converts this expression to an integer, if possible, or produces an error if not.
// It accepts big ints, uints, and ints.
// If max is not -1, produces an error if the value exceeds max.
// If negErr is not "", produces an error if the value is negative.
func (x *Expr) convertToInt(max int64, negErr string, errOp string) *Expr {
	switch x.ExprType.Lit().(type) {
	case *types.IdealIntType:
		val := x.asIdealInt()()
		if negErr != "" && val.Sign() < 0 {
			x.error("negative %s: %s", negErr, val)
			return nil
		}
		bound := max
		if negErr == "slice" {
			bound++
		}
		if max != -1 && val.Cmp(big.NewInt(bound)) >= 0 {
			x.error("index %s exceeds length %d", val, max)
			return nil
		}
		return x.resolveIdeal(builtins.IntType)

	case *types.UintType:
		// Convert to int
		na := x.newExpr(builtins.IntType, x.desc)
		af := x.asUint()
		na.eval = func(t *vm.Thread) int64 { return int64(af(t)) }
		return na

	case *types.IntType:
		// Good as is
		return x
	}

	x.error("illegal operand type for %s: %v", errOp, x.ExprType)
	return nil
}

// derefArray returns an expression of array type if the given expression is a *array type.
// Otherwise, returns the given expression.
func (x *Expr) derefArray() *Expr {
	if pt, ok := x.ExprType.Lit().(*types.PtrType); ok {
		if _, ok := pt.Elem.Lit().(*types.ArrayType); ok {
			deref := x.compileStarExpr(x)
			if deref == nil {
				logger.Panic().Msg("failed to dereference *array")
			}
			return deref
		}
	}
	return x
}

////////////////////////////////////////////////////////////////////////////////
// "As" functions retrieve evaluator functions from an expr, panicking if
// the requested evaluator has the wrong type.

// asValue returns a closure around a Value, according to the type of the underlying expression
func (x *Expr) asValue() func(t *vm.Thread) vm.Value {
	var fct func(t *vm.Thread) vm.Value
	switch ty := x.ExprType.Lit().(type) {
	case *types.BoolType:
		fct = func(t *vm.Thread) vm.Value {
			b := values.BoolV(x.asBool()(t))
			var val vm.Value = &b
			return val
		}
	case *types.UintType:
		v := x.asUint()
		switch ty.Bits {
		case 0:
			switch ty.Ptr {
			case true:
				fct = func(t *vm.Thread) vm.Value {
					vv := values.UintptrV(v(t))
					var val vm.Value = &vv
					return val
				}
			case false:
				fct = func(t *vm.Thread) vm.Value {
					vv := values.UintV(v(t))
					var val vm.Value = &vv
					return val
				}
			}
		case 8:
			fct = func(t *vm.Thread) vm.Value {
				vv := values.Uint8V(v(t))
				return &vv
			}
		case 16:
			fct = func(t *vm.Thread) vm.Value {
				vv := values.Uint16V(v(t))
				return &vv
			}
		case 32:
			fct = func(t *vm.Thread) vm.Value {
				vv := values.Uint32V(v(t))
				return &vv
			}
		case 64:
			fct = func(t *vm.Thread) vm.Value {
				vv := values.Uint64V(v(t))
				return &vv
			}
		}
	case *types.IntType:
		v := x.asInt()
		switch ty.Bits {
		case 0:
			fct = func(t *vm.Thread) vm.Value {
				vv := values.IntV(v(t))
				return &vv
			}
		case 8:
			fct = func(t *vm.Thread) vm.Value {
				vv := values.Int8V(v(t))
				return &vv
			}
		case 16:
			fct = func(t *vm.Thread) vm.Value {
				vv := values.Int16V(v(t))
				return &vv
			}
		case 32:
			fct = func(t *vm.Thread) vm.Value {
				vv := values.Int32V(v(t))
				return &vv
			}
		case 64:
			fct = func(t *vm.Thread) vm.Value {
				vv := values.Int64V(v(t))
				return &vv
			}
		}
	case *types.FloatType:
		v := x.asFloat()
		switch ty.Bits {
		case 32:
			fct = func(t *vm.Thread) vm.Value {
				vv := values.Float32V(v(t))
				return &vv
			}
		case 64:
			fct = func(t *vm.Thread) vm.Value {
				vv := values.Float64V(v(t))
				return &vv
			}
		}
	case *types.IdealIntType:
		fct = func(t *vm.Thread) vm.Value {
			v := x.asIdealInt()()
			vv := values.IdealIntV{V: v}
			return &vv
		}
	case *types.IdealFloatType:
		fct = func(t *vm.Thread) vm.Value {
			v := x.asIdealFloat()()
			vv := values.IdealFloatV{V: v}
			return &vv
		}
	case *types.StringType:
		fct = func(t *vm.Thread) vm.Value {
			v := x.asString()
			vv := values.StringV(v(t))
			return &vv
		}
	case *types.ArrayType:
		fct = func(t *vm.Thread) vm.Value {
			v := x.asArray()
			vv := v(t).Get(t)
			return vv
		}
	case *types.StructType:
		fct = func(t *vm.Thread) vm.Value {
			v := x.asStruct()
			vv := v(t).Get(t)
			return vv
		}
	case *types.PtrType:
		fct = func(t *vm.Thread) vm.Value {
			return x.asPtr()(t)
		}
	case *types.PackageType:
		fct = func(t *vm.Thread) vm.Value {
			return x.asPackage()(t)
		}
	case *types.FuncType:
		x.error("unhandled type: %v", ty.String())
	case *types.InterfaceType:
		x.error("unhandled type: %v", ty.String())
	case *types.SliceType:
		x.error("unhandled type: %v", ty.String())
	case *types.MapType:
		x.error("unhandled type: %v", ty.String())
	case *types.ChanType:
		x.error("unhandled type: %v", ty.String())
	case *types.NamedType:
		x.error("unhandled type: %v", ty.String())
	case *types.MultiType:
		x.error("unhandled type: %v", ty.String())
	default:
		x.error("unhandled type: %v", ty.String())
	}
	return fct
}

func (x *Expr) asInterface() func(*vm.Thread) interface{} {
	switch sf := x.eval.(type) {
	case func(t *vm.Thread) bool:
		return func(t *vm.Thread) interface{} { return sf(t) }
	case func(t *vm.Thread) uint64:
		return func(t *vm.Thread) interface{} { return sf(t) }
	case func(t *vm.Thread) int64:
		return func(t *vm.Thread) interface{} { return sf(t) }
	case func() *big.Int:
		return func(*vm.Thread) interface{} { return sf() }
	case func(t *vm.Thread) float64:
		return func(t *vm.Thread) interface{} { return sf(t) }
	case func() *big.Rat:
		return func(*vm.Thread) interface{} { return sf() }
	case func(t *vm.Thread) string:
		return func(t *vm.Thread) interface{} { return sf(t) }
	case func(t *vm.Thread) value.ArrayValue:
		return func(t *vm.Thread) interface{} { return sf(t) }
	case func(t *vm.Thread) value.StructValue:
		return func(t *vm.Thread) interface{} { return sf(t) }
	case func(t *vm.Thread) vm.Value:
		return func(t *vm.Thread) interface{} { return sf(t) }
	case func(t *vm.Thread) value.Func:
		return func(t *vm.Thread) interface{} { return sf(t) }
	case func(t *vm.Thread) value.Slice:
		return func(t *vm.Thread) interface{} { return sf(t) }
	case func(t *vm.Thread) value.Map:
		return func(t *vm.Thread) interface{} { return sf(t) }
	default:
		logger.Panic().Msgf("unexpected expression node type %T at %v", x.eval, x.pos)
	}

	panic("fail")
}

func (x *Expr) asPackage() func(*vm.Thread) values.PackageValue {
	return x.eval.(func(*vm.Thread) values.PackageValue)
}
func (x *Expr) asIdealInt() func() *big.Int {
	return x.eval.(func() *big.Int)
}
func (x *Expr) asIdealFloat() func() *big.Rat {
	return x.eval.(func() *big.Rat)
}
func (x *Expr) asBool() func(*vm.Thread) bool {
	return x.eval.(func(*vm.Thread) bool)
}
func (x *Expr) asUint() func(*vm.Thread) uint64 {
	return x.eval.(func(*vm.Thread) uint64)
}
func (x *Expr) asInt() func(*vm.Thread) int64 {
	return x.eval.(func(*vm.Thread) int64)
}
func (x *Expr) asFloat() func(*vm.Thread) float64 {
	return x.eval.(func(*vm.Thread) float64)
}
func (x *Expr) asString() func(*vm.Thread) string {
	return x.eval.(func(*vm.Thread) string)
}
func (x *Expr) asArray() func(*vm.Thread) values.ArrayValue {
	return x.eval.(func(*vm.Thread) values.ArrayValue)
}
func (x *Expr) asStruct() func(*vm.Thread) values.StructValue {
	return x.eval.(func(*vm.Thread) values.StructValue)
}
func (x *Expr) asPtr() func(*vm.Thread) vm.Value {
	return x.eval.(func(*vm.Thread) vm.Value)
}
func (x *Expr) asFunc() func(*vm.Thread) values.Func {
	return x.eval.(func(*vm.Thread) values.Func)
}
func (x *Expr) asSlice() func(*vm.Thread) values.Slice {
	return x.eval.(func(*vm.Thread) values.Slice)
}
func (x *Expr) asMap() func(*vm.Thread) values.Map {
	return x.eval.(func(*vm.Thread) values.Map)
}
func (x *Expr) asMulti() func(*vm.Thread) []vm.Value {
	return x.eval.(func(*vm.Thread) []vm.Value)
}

////////////////////////////////////////////////////////////////////////////////
