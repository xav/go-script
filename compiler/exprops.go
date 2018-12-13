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

	"github.com/xav/go-script/types"
	"github.com/xav/go-script/values"
	"github.com/xav/go-script/vm"
)

// genAssign returns a function that will assign the value of r to lv.
func genAssign(lt vm.Type, r *Expr) func(lv vm.Value, t *vm.Thread) {
	switch lt.Lit().(type) {
	case *types.BoolType:
		rf := r.asBool()
		return func(lv vm.Value, t *vm.Thread) {
			lv.(values.BoolValue).Set(t, rf(t))
		}
	case *types.UintType:
		rf := r.asUint()
		return func(lv vm.Value, t *vm.Thread) {
			lv.(values.UintValue).Set(t, rf(t))
		}
	case *types.IntType:
		rf := r.asInt()
		return func(lv vm.Value, t *vm.Thread) {
			lv.(values.IntValue).Set(t, rf(t))
		}
	case *types.FloatType:
		rf := r.asFloat()
		return func(lv vm.Value, t *vm.Thread) {
			lv.(values.FloatValue).Set(t, rf(t))
		}
	case *types.StringType:
		rf := r.asString()
		return func(lv vm.Value, t *vm.Thread) {
			lv.(values.StringValue).Set(t, rf(t))
		}
	case *types.ArrayType:
		rf := r.asArray()
		return func(lv vm.Value, t *vm.Thread) {
			lv.Assign(t, rf(t))
		}
	case *types.StructType:
		rf := r.asStruct()
		return func(lv vm.Value, t *vm.Thread) {
			lv.Assign(t, rf(t))
		}
	case *types.PtrType:
		rf := r.asPtr()
		return func(lv vm.Value, t *vm.Thread) {
			lv.(values.PtrValue).Set(t, rf(t))
		}
	case *types.FuncType:
		rf := r.asFunc()
		return func(lv vm.Value, t *vm.Thread) {
			lv.(values.FuncValue).Set(t, rf(t))
		}
	case *types.SliceType:
		rf := r.asSlice()
		return func(lv vm.Value, t *vm.Thread) {
			lv.(values.SliceValue).Set(t, rf(t))
		}
	case *types.MapType:
		rf := r.asMap()
		return func(lv vm.Value, t *vm.Thread) {
			lv.(values.MapValue).Set(t, rf(t))
		}
	default:
		logger.Panic().
			Str("type", fmt.Sprintf("%v", lt)).
			Str("pos", fmt.Sprintf("%v", r.pos)).
			Msg("unexpected left operand type")
	}
	panic("unreachable")
}

////////////////////////////////////////////////////////////////////////////////

// genBinOpAdd sets x.eval to a function that will return the result of 'l + r'.
func (x *Expr) genBinOpAdd(l, r *Expr) {
	switch t := l.ExprType.Lit().(type) {
	case *types.UintType:
		lf := l.asUint()
		rf := r.asUint()
		switch t.Bits {
		case 8:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l + r
				return uint64(uint8(ret))
			}
		case 16:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l + r
				return uint64(uint16(ret))
			}
		case 32:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l + r
				return uint64(uint32(ret))
			}
		case 64:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l + r
				return uint64(uint64(ret))
			}
		case 0:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l + r
				return uint64(uint(ret))
			}
		default:
			logger.Panic().
				Str("type", fmt.Sprintf("%v", t)).
				Str("pos", fmt.Sprintf("%v", x.pos)).
				Msgf("unexpected size %d", t.Bits)
		}
	case *types.IntType:
		lf := l.asInt()
		rf := r.asInt()
		switch t.Bits {
		case 8:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l + r
				return int64(int8(ret))
			}
		case 16:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l + r
				return int64(int16(ret))
			}
		case 32:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l + r
				return int64(int32(ret))
			}
		case 64:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l + r
				return int64(int64(ret))
			}
		case 0:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l + r
				return int64(int(ret))
			}
		default:
			logger.Panic().
				Str("type", fmt.Sprintf("%v", t)).
				Str("pos", fmt.Sprintf("%v", x.pos)).
				Msgf("unexpected size %d", t.Bits)
		}
	case *types.FloatType:
		lf := l.asFloat()
		rf := r.asFloat()
		switch t.Bits {
		case 32:
			x.eval = func(t *vm.Thread) float64 {
				l, r := lf(t), rf(t)
				var ret float64
				ret = l + r
				return float64(float32(ret))
			}
		case 64:
			x.eval = func(t *vm.Thread) float64 {
				l, r := lf(t), rf(t)
				var ret float64
				ret = l + r
				return float64(float64(ret))
			}
		default:
			logger.Panic().
				Str("type", fmt.Sprintf("%v", t)).
				Str("pos", fmt.Sprintf("%v", x.pos)).
				Msgf("unexpected size %d", t.Bits)
		}
	case *types.StringType:
		lf := l.asString()
		rf := r.asString()
		x.eval = func(t *vm.Thread) string {
			l, r := lf(t), rf(t)
			return l + r
		}
	case *types.IdealIntType:
		l := l.asIdealInt()()
		r := r.asIdealInt()()
		val := l.Add(l, r)
		x.eval = func() *big.Int {
			return val
		}
	case *types.IdealFloatType:
		l := l.asIdealFloat()()
		r := r.asIdealFloat()()
		val := l.Add(l, r)
		x.eval = func() *big.Rat {
			return val
		}
	default:
		logger.Panic().
			Str("type", fmt.Sprintf("%v", t)).
			Str("pos", fmt.Sprintf("%v", x.pos)).
			Msg("unexpected type")
	}
}

// genBinOpSub sets x.eval to a function that will return the result of 'l - r'.
func (x *Expr) genBinOpSub(l, r *Expr) {
	switch t := l.ExprType.Lit().(type) {
	case *types.UintType:
		lf := l.asUint()
		rf := r.asUint()
		switch t.Bits {
		case 8:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l - r
				return uint64(uint8(ret))
			}
		case 16:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l - r
				return uint64(uint16(ret))
			}
		case 32:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l - r
				return uint64(uint32(ret))
			}
		case 64:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l - r
				return uint64(uint64(ret))
			}
		case 0:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l - r
				return uint64(uint(ret))
			}
		default:
			logger.Panic().
				Str("type", fmt.Sprintf("%v", t)).
				Str("pos", fmt.Sprintf("%v", x.pos)).
				Msgf("unexpected size %d", t.Bits)
		}
	case *types.IntType:
		lf := l.asInt()
		rf := r.asInt()
		switch t.Bits {
		case 8:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l - r
				return int64(int8(ret))
			}
		case 16:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l - r
				return int64(int16(ret))
			}
		case 32:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l - r
				return int64(int32(ret))
			}
		case 64:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l - r
				return int64(int64(ret))
			}
		case 0:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l - r
				return int64(int(ret))
			}
		default:
			logger.Panic().
				Str("type", fmt.Sprintf("%v", t)).
				Str("pos", fmt.Sprintf("%v", x.pos)).
				Msgf("unexpected size %d", t.Bits)
		}
	case *types.FloatType:
		lf := l.asFloat()
		rf := r.asFloat()
		switch t.Bits {
		case 32:
			x.eval = func(t *vm.Thread) float64 {
				l, r := lf(t), rf(t)
				var ret float64
				ret = l - r
				return float64(float32(ret))
			}
		case 64:
			x.eval = func(t *vm.Thread) float64 {
				l, r := lf(t), rf(t)
				var ret float64
				ret = l - r
				return float64(float64(ret))
			}
		default:
			logger.Panic().
				Str("type", fmt.Sprintf("%v", t)).
				Str("pos", fmt.Sprintf("%v", x.pos)).
				Msgf("unexpected size %d", t.Bits)
		}
	case *types.IdealIntType:
		l := l.asIdealInt()()
		r := r.asIdealInt()()
		val := l.Sub(l, r)
		x.eval = func() *big.Int {
			return val
		}
	case *types.IdealFloatType:
		l := l.asIdealFloat()()
		r := r.asIdealFloat()()
		val := l.Sub(l, r)
		x.eval = func() *big.Rat {
			return val
		}
	default:
		logger.Panic().
			Str("type", fmt.Sprintf("%v", t)).
			Str("pos", fmt.Sprintf("%v", x.pos)).
			Msg("unexpected type")
	}
}

// genBinOpMul sets x.eval to a function that will return the result of 'l * r'.
func (x *Expr) genBinOpMul(l, r *Expr) {
	switch t := l.ExprType.Lit().(type) {
	case *types.UintType:
		lf := l.asUint()
		rf := r.asUint()
		switch t.Bits {
		case 8:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l * r
				return uint64(uint8(ret))
			}
		case 16:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l * r
				return uint64(uint16(ret))
			}
		case 32:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l * r
				return uint64(uint32(ret))
			}
		case 64:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l * r
				return uint64(uint64(ret))
			}
		case 0:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l * r
				return uint64(uint(ret))
			}
		default:
			logger.Panic().
				Str("type", fmt.Sprintf("%v", t)).
				Str("pos", fmt.Sprintf("%v", x.pos)).
				Msgf("unexpected size %d", t.Bits)
		}
	case *types.IntType:
		lf := l.asInt()
		rf := r.asInt()
		switch t.Bits {
		case 8:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l * r
				return int64(int8(ret))
			}
		case 16:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l * r
				return int64(int16(ret))
			}
		case 32:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l * r
				return int64(int32(ret))
			}
		case 64:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l * r
				return int64(int64(ret))
			}
		case 0:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l * r
				return int64(int(ret))
			}
		default:
			logger.Panic().
				Str("type", fmt.Sprintf("%v", t)).
				Str("pos", fmt.Sprintf("%v", x.pos)).
				Msgf("unexpected size %d", t.Bits)
		}
	case *types.FloatType:
		lf := l.asFloat()
		rf := r.asFloat()
		switch t.Bits {
		case 32:
			x.eval = func(t *vm.Thread) float64 {
				l, r := lf(t), rf(t)
				var ret float64
				ret = l * r
				return float64(float32(ret))
			}
		case 64:
			x.eval = func(t *vm.Thread) float64 {
				l, r := lf(t), rf(t)
				var ret float64
				ret = l * r
				return float64(float64(ret))
			}
		default:
			logger.Panic().
				Str("type", fmt.Sprintf("%v", t)).
				Str("pos", fmt.Sprintf("%v", x.pos)).
				Msgf("unexpected size %d", t.Bits)
		}
	case *types.IdealIntType:
		l := l.asIdealInt()()
		r := r.asIdealInt()()
		val := l.Mul(l, r)
		x.eval = func() *big.Int {
			return val
		}
	case *types.IdealFloatType:
		l := l.asIdealFloat()()
		r := r.asIdealFloat()()
		val := l.Mul(l, r)
		x.eval = func() *big.Rat {
			return val
		}
	default:
		logger.Panic().
			Str("type", fmt.Sprintf("%v", t)).
			Str("pos", fmt.Sprintf("%v", x.pos)).
			Msg("unexpected type")
	}
}

// genBinOpQuo sets x.eval to a function that will return the result of 'l / r'.
func (x *Expr) genBinOpQuo(l, r *Expr) {
	switch t := l.ExprType.Lit().(type) {
	case *types.UintType:
		lf := l.asUint()
		rf := r.asUint()
		switch t.Bits {
		case 8:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				if r == 0 {
					t.Abort(vm.DivByZeroError{})
				}
				ret = l / r
				return uint64(uint8(ret))
			}
		case 16:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				if r == 0 {
					t.Abort(vm.DivByZeroError{})
				}
				ret = l / r
				return uint64(uint16(ret))
			}
		case 32:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				if r == 0 {
					t.Abort(vm.DivByZeroError{})
				}
				ret = l / r
				return uint64(uint32(ret))
			}
		case 64:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				if r == 0 {
					t.Abort(vm.DivByZeroError{})
				}
				ret = l / r
				return uint64(uint64(ret))
			}
		case 0:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				if r == 0 {
					t.Abort(vm.DivByZeroError{})
				}
				ret = l / r
				return uint64(uint(ret))
			}
		default:
			logger.Panic().
				Str("type", fmt.Sprintf("%v", t)).
				Str("pos", fmt.Sprintf("%v", x.pos)).
				Msgf("unexpected size %d", t.Bits)
		}
	case *types.IntType:
		lf := l.asInt()
		rf := r.asInt()
		switch t.Bits {
		case 8:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				if r == 0 {
					t.Abort(vm.DivByZeroError{})
				}
				ret = l / r
				return int64(int8(ret))
			}
		case 16:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				if r == 0 {
					t.Abort(vm.DivByZeroError{})
				}
				ret = l / r
				return int64(int16(ret))
			}
		case 32:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				if r == 0 {
					t.Abort(vm.DivByZeroError{})
				}
				ret = l / r
				return int64(int32(ret))
			}
		case 64:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				if r == 0 {
					t.Abort(vm.DivByZeroError{})
				}
				ret = l / r
				return int64(int64(ret))
			}
		case 0:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				if r == 0 {
					t.Abort(vm.DivByZeroError{})
				}
				ret = l / r
				return int64(int(ret))
			}
		default:
			logger.Panic().
				Str("type", fmt.Sprintf("%v", t)).
				Str("pos", fmt.Sprintf("%v", x.pos)).
				Msgf("unexpected size %d", t.Bits)
		}
	case *types.FloatType:
		lf := l.asFloat()
		rf := r.asFloat()
		switch t.Bits {
		case 32:
			x.eval = func(t *vm.Thread) float64 {
				l, r := lf(t), rf(t)
				var ret float64
				if r == 0 {
					t.Abort(vm.DivByZeroError{})
				}
				ret = l / r
				return float64(float32(ret))
			}
		case 64:
			x.eval = func(t *vm.Thread) float64 {
				l, r := lf(t), rf(t)
				var ret float64
				if r == 0 {
					t.Abort(vm.DivByZeroError{})
				}
				ret = l / r
				return float64(float64(ret))
			}
		default:
			logger.Panic().
				Str("type", fmt.Sprintf("%v", t)).
				Str("pos", fmt.Sprintf("%v", x.pos)).
				Msgf("unexpected size %d", t.Bits)
		}
	case *types.IdealIntType:
		l := l.asIdealInt()()
		r := r.asIdealInt()()
		val := l.Quo(l, r)
		x.eval = func() *big.Int {
			return val
		}
	case *types.IdealFloatType:
		l := l.asIdealFloat()()
		r := r.asIdealFloat()()
		val := l.Quo(l, r)
		x.eval = func() *big.Rat {
			return val
		}
	default:
		logger.Panic().
			Str("type", fmt.Sprintf("%v", t)).
			Str("pos", fmt.Sprintf("%v", x.pos)).
			Msg("unexpected type")
	}
}

// genBinOpRem sets x.eval to a function that will return the result of 'l % r'.
func (x *Expr) genBinOpRem(l, r *Expr) {
	switch t := l.ExprType.Lit().(type) {
	case *types.UintType:
		lf := l.asUint()
		rf := r.asUint()
		switch t.Bits {
		case 8:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				if r == 0 {
					t.Abort(vm.DivByZeroError{})
				}
				ret = l % r
				return uint64(uint8(ret))
			}
		case 16:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				if r == 0 {
					t.Abort(vm.DivByZeroError{})
				}
				ret = l % r
				return uint64(uint16(ret))
			}
		case 32:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				if r == 0 {
					t.Abort(vm.DivByZeroError{})
				}
				ret = l % r
				return uint64(uint32(ret))
			}
		case 64:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				if r == 0 {
					t.Abort(vm.DivByZeroError{})
				}
				ret = l % r
				return uint64(uint64(ret))
			}
		case 0:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				if r == 0 {
					t.Abort(vm.DivByZeroError{})
				}
				ret = l % r
				return uint64(uint(ret))
			}
		default:
			logger.Panic().
				Str("type", fmt.Sprintf("%v", t)).
				Str("pos", fmt.Sprintf("%v", x.pos)).
				Msgf("unexpected size %d", t.Bits)
		}
	case *types.IntType:
		lf := l.asInt()
		rf := r.asInt()
		switch t.Bits {
		case 8:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				if r == 0 {
					t.Abort(vm.DivByZeroError{})
				}
				ret = l % r
				return int64(int8(ret))
			}
		case 16:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				if r == 0 {
					t.Abort(vm.DivByZeroError{})
				}
				ret = l % r
				return int64(int16(ret))
			}
		case 32:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				if r == 0 {
					t.Abort(vm.DivByZeroError{})
				}
				ret = l % r
				return int64(int32(ret))
			}
		case 64:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				if r == 0 {
					t.Abort(vm.DivByZeroError{})
				}
				ret = l % r
				return int64(int64(ret))
			}
		case 0:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				if r == 0 {
					t.Abort(vm.DivByZeroError{})
				}
				ret = l % r
				return int64(int(ret))
			}
		default:
			logger.Panic().
				Str("type", fmt.Sprintf("%v", t)).
				Str("pos", fmt.Sprintf("%v", x.pos)).
				Msgf("unexpected size %d", t.Bits)
		}
	case *types.IdealIntType:
		l := l.asIdealInt()()
		r := r.asIdealInt()()
		val := l.Rem(l, r)
		x.eval = func() *big.Int {
			return val
		}
	default:
		logger.Panic().
			Str("type", fmt.Sprintf("%v", t)).
			Str("pos", fmt.Sprintf("%v", x.pos)).
			Msg("unexpected type")
	}
}

// genBinOpAnd sets x.eval to a function that will return the result of 'l & r'.
func (x *Expr) genBinOpAnd(l, r *Expr) {
	switch t := l.ExprType.Lit().(type) {
	case *types.UintType:
		lf := l.asUint()
		rf := r.asUint()
		switch t.Bits {
		case 8:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l & r
				return uint64(uint8(ret))
			}
		case 16:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l & r
				return uint64(uint16(ret))
			}
		case 32:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l & r
				return uint64(uint32(ret))
			}
		case 64:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l & r
				return uint64(uint64(ret))
			}
		case 0:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l & r
				return uint64(uint(ret))
			}
		default:
			logger.Panic().
				Str("type", fmt.Sprintf("%v", t)).
				Str("pos", fmt.Sprintf("%v", x.pos)).
				Msgf("unexpected size %d", t.Bits)
		}
	case *types.IntType:
		lf := l.asInt()
		rf := r.asInt()
		switch t.Bits {
		case 8:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l & r
				return int64(int8(ret))
			}
		case 16:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l & r
				return int64(int16(ret))
			}
		case 32:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l & r
				return int64(int32(ret))
			}
		case 64:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l & r
				return int64(int64(ret))
			}
		case 0:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l & r
				return int64(int(ret))
			}
		default:
			logger.Panic().
				Str("type", fmt.Sprintf("%v", t)).
				Str("pos", fmt.Sprintf("%v", x.pos)).
				Msgf("unexpected size %d", t.Bits)
		}
	case *types.IdealIntType:
		l := l.asIdealInt()()
		r := r.asIdealInt()()
		val := l.And(l, r)
		x.eval = func() *big.Int {
			return val
		}
	default:
		logger.Panic().
			Str("type", fmt.Sprintf("%v", t)).
			Str("pos", fmt.Sprintf("%v", x.pos)).
			Msg("unexpected type")
	}
}

// genBinOpOr sets x.eval to a function that will return the result of 'l | r'.
func (x *Expr) genBinOpOr(l, r *Expr) {
	switch t := l.ExprType.Lit().(type) {
	case *types.UintType:
		lf := l.asUint()
		rf := r.asUint()
		switch t.Bits {
		case 8:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l | r
				return uint64(uint8(ret))
			}
		case 16:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l | r
				return uint64(uint16(ret))
			}
		case 32:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l | r
				return uint64(uint32(ret))
			}
		case 64:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l | r
				return uint64(uint64(ret))
			}
		case 0:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l | r
				return uint64(uint(ret))
			}
		default:
			logger.Panic().
				Str("type", fmt.Sprintf("%v", t)).
				Str("pos", fmt.Sprintf("%v", x.pos)).
				Msgf("unexpected size %d", t.Bits)
		}
	case *types.IntType:
		lf := l.asInt()
		rf := r.asInt()
		switch t.Bits {
		case 8:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l | r
				return int64(int8(ret))
			}
		case 16:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l | r
				return int64(int16(ret))
			}
		case 32:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l | r
				return int64(int32(ret))
			}
		case 64:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l | r
				return int64(int64(ret))
			}
		case 0:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l | r
				return int64(int(ret))
			}
		default:
			logger.Panic().
				Str("type", fmt.Sprintf("%v", t)).
				Str("pos", fmt.Sprintf("%v", x.pos)).
				Msgf("unexpected size %d", t.Bits)
		}
	case *types.IdealIntType:
		l := l.asIdealInt()()
		r := r.asIdealInt()()
		val := l.Or(l, r)
		x.eval = func() *big.Int {
			return val
		}
	default:
		logger.Panic().
			Str("type", fmt.Sprintf("%v", t)).
			Str("pos", fmt.Sprintf("%v", x.pos)).
			Msg("unexpected type")
	}
}

// genBinOpXor sets x.eval to a function that will return the result of 'l ^ r'.
func (x *Expr) genBinOpXor(l, r *Expr) {
	switch t := l.ExprType.Lit().(type) {
	case *types.UintType:
		lf := l.asUint()
		rf := r.asUint()
		switch t.Bits {
		case 8:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l ^ r
				return uint64(uint8(ret))
			}
		case 16:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l ^ r
				return uint64(uint16(ret))
			}
		case 32:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l ^ r
				return uint64(uint32(ret))
			}
		case 64:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l ^ r
				return uint64(uint64(ret))
			}
		case 0:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l ^ r
				return uint64(uint(ret))
			}
		default:
			logger.Panic().
				Str("type", fmt.Sprintf("%v", t)).
				Str("pos", fmt.Sprintf("%v", x.pos)).
				Msgf("unexpected size %d", t.Bits)
		}
	case *types.IntType:
		lf := l.asInt()
		rf := r.asInt()
		switch t.Bits {
		case 8:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l ^ r
				return int64(int8(ret))
			}
		case 16:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l ^ r
				return int64(int16(ret))
			}
		case 32:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l ^ r
				return int64(int32(ret))
			}
		case 64:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l ^ r
				return int64(int64(ret))
			}
		case 0:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l ^ r
				return int64(int(ret))
			}
		default:
			logger.Panic().
				Str("type", fmt.Sprintf("%v", t)).
				Str("pos", fmt.Sprintf("%v", x.pos)).
				Msgf("unexpected size %d", t.Bits)
		}
	case *types.IdealIntType:
		l := l.asIdealInt()()
		r := r.asIdealInt()()
		val := l.Xor(l, r)
		x.eval = func() *big.Int {
			return val
		}
	default:
		logger.Panic().
			Str("type", fmt.Sprintf("%v", t)).
			Str("pos", fmt.Sprintf("%v", x.pos)).
			Msg("unexpected type")
	}
}

// genBinOpAndNot sets x.eval to a function that will return the result of 'l &^ r'.
func (x *Expr) genBinOpAndNot(l, r *Expr) {
	switch t := l.ExprType.Lit().(type) {
	case *types.UintType:
		lf := l.asUint()
		rf := r.asUint()
		switch t.Bits {
		case 8:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l &^ r
				return uint64(uint8(ret))
			}
		case 16:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l &^ r
				return uint64(uint16(ret))
			}
		case 32:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l &^ r
				return uint64(uint32(ret))
			}
		case 64:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l &^ r
				return uint64(uint64(ret))
			}
		case 0:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l &^ r
				return uint64(uint(ret))
			}
		default:
			logger.Panic().
				Str("type", fmt.Sprintf("%v", t)).
				Str("pos", fmt.Sprintf("%v", x.pos)).
				Msgf("unexpected size %d", t.Bits)
		}
	case *types.IntType:
		lf := l.asInt()
		rf := r.asInt()
		switch t.Bits {
		case 8:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l &^ r
				return int64(int8(ret))
			}
		case 16:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l &^ r
				return int64(int16(ret))
			}
		case 32:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l &^ r
				return int64(int32(ret))
			}
		case 64:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l &^ r
				return int64(int64(ret))
			}
		case 0:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l &^ r
				return int64(int(ret))
			}
		default:
			logger.Panic().
				Str("type", fmt.Sprintf("%v", t)).
				Str("pos", fmt.Sprintf("%v", x.pos)).
				Msgf("unexpected size %d", t.Bits)
		}
	case *types.IdealIntType:
		l := l.asIdealInt()()
		r := r.asIdealInt()()
		val := l.AndNot(l, r)
		x.eval = func() *big.Int {
			return val
		}
	default:
		logger.Panic().
			Str("type", fmt.Sprintf("%v", t)).
			Str("pos", fmt.Sprintf("%v", x.pos)).
			Msg("unexpected type")
	}
}

// genBinOpShl sets x.eval to a function that will return the result of 'l << r'.
func (x *Expr) genBinOpShl(l, r *Expr) {
	switch t := l.ExprType.Lit().(type) {
	case *types.UintType:
		lf := l.asUint()
		rf := r.asUint()
		switch t.Bits {
		case 8:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l << r
				return uint64(uint8(ret))
			}
		case 16:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l << r
				return uint64(uint16(ret))
			}
		case 32:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l << r
				return uint64(uint32(ret))
			}
		case 64:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l << r
				return uint64(uint64(ret))
			}
		case 0:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l << r
				return uint64(uint(ret))
			}
		default:
			logger.Panic().
				Str("type", fmt.Sprintf("%v", t)).
				Str("pos", fmt.Sprintf("%v", x.pos)).
				Msgf("unexpected size %d", t.Bits)
		}
	case *types.IntType:
		lf := l.asInt()
		rf := r.asUint()
		switch t.Bits {
		case 8:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l << r
				return int64(int8(ret))
			}
		case 16:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l << r
				return int64(int16(ret))
			}
		case 32:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l << r
				return int64(int32(ret))
			}
		case 64:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l << r
				return int64(int64(ret))
			}
		case 0:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l << r
				return int64(int(ret))
			}
		default:
			logger.Panic().
				Str("type", fmt.Sprintf("%v", t)).
				Str("pos", fmt.Sprintf("%v", x.pos)).
				Msgf("unexpected size %d", t.Bits)
		}
	default:
		logger.Panic().
			Str("type", fmt.Sprintf("%v", t)).
			Str("pos", fmt.Sprintf("%v", x.pos)).
			Msg("unexpected type")
	}
}

// genBinOpShr sets x.eval to a function that will return the result of 'l >> r'.
func (x *Expr) genBinOpShr(l, r *Expr) {
	switch t := l.ExprType.Lit().(type) {
	case *types.UintType:
		lf := l.asUint()
		rf := r.asUint()
		switch t.Bits {
		case 8:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l >> r
				return uint64(uint8(ret))
			}
		case 16:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l >> r
				return uint64(uint16(ret))
			}
		case 32:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l >> r
				return uint64(uint32(ret))
			}
		case 64:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l >> r
				return uint64(uint64(ret))
			}
		case 0:
			x.eval = func(t *vm.Thread) uint64 {
				l, r := lf(t), rf(t)
				var ret uint64
				ret = l >> r
				return uint64(uint(ret))
			}
		default:
			logger.Panic().
				Str("type", fmt.Sprintf("%v", t)).
				Str("pos", fmt.Sprintf("%v", x.pos)).
				Msgf("unexpected size %d", t.Bits)

		}
	case *types.IntType:
		lf := l.asInt()
		rf := r.asUint()
		switch t.Bits {
		case 8:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l >> r
				return int64(int8(ret))
			}
		case 16:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l >> r
				return int64(int16(ret))
			}
		case 32:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l >> r
				return int64(int32(ret))
			}
		case 64:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l >> r
				return int64(int64(ret))
			}
		case 0:
			x.eval = func(t *vm.Thread) int64 {
				l, r := lf(t), rf(t)
				var ret int64
				ret = l >> r
				return int64(int(ret))
			}
		default:
			logger.Panic().
				Str("type", fmt.Sprintf("%v", t)).
				Str("pos", fmt.Sprintf("%v", x.pos)).
				Msgf("unexpected size %d", t.Bits)

		}
	default:
		logger.Panic().
			Str("type", fmt.Sprintf("%v", t)).
			Str("pos", fmt.Sprintf("%v", x.pos)).
			Msg("unexpected type")
	}
}

// genBinOpLogAnd sets x.eval to a function that will return the result of 'l && r'.
func (x *Expr) genBinOpLogAnd(l, r *Expr) {
	lf := l.asBool()
	rf := r.asBool()
	x.eval = func(t *vm.Thread) bool {
		return lf(t) && rf(t)
	}
}

// genBinOpLogOr sets x.eval to a function that will return the result of 'l || r'.
func (x *Expr) genBinOpLogOr(l, r *Expr) {
	lf := l.asBool()
	rf := r.asBool()
	x.eval = func(t *vm.Thread) bool {
		return lf(t) || rf(t)
	}
}

// genBinOpLogOr sets x.eval to a function that will return the result of 'l <- r'.
func (x *Expr) genBinSend(l, r *Expr) {
	panic("Binary op <- not implemented")
}

// genBinOpEql sets x.eval to a function that will return the result of 'l == r'.
func (x *Expr) genBinOpEql(l, r *Expr) {
	switch l.ExprType.Lit().(type) {
	case *types.BoolType:
		lf := l.asBool()
		rf := r.asBool()
		x.eval = func(t *vm.Thread) bool {
			l, r := lf(t), rf(t)
			return l == r
		}
	case *types.UintType:
		lf := l.asUint()
		rf := r.asUint()
		x.eval = func(t *vm.Thread) bool {
			l, r := lf(t), rf(t)
			return l == r
		}
	case *types.IntType:
		lf := l.asInt()
		rf := r.asInt()
		x.eval = func(t *vm.Thread) bool {
			l, r := lf(t), rf(t)
			return l == r
		}
	case *types.FloatType:
		lf := l.asFloat()
		rf := r.asFloat()
		x.eval = func(t *vm.Thread) bool {
			l, r := lf(t), rf(t)
			return l == r
		}
	case *types.IdealIntType:
		l := l.asIdealInt()()
		r := r.asIdealInt()()
		val := l.Cmp(r) == 0
		x.eval = func(t *vm.Thread) bool {
			return val
		}
	case *types.IdealFloatType:
		l := l.asIdealFloat()()
		r := r.asIdealFloat()()
		val := l.Cmp(r) == 0
		x.eval = func(t *vm.Thread) bool {
			return val
		}
	case *types.StringType:
		lf := l.asString()
		rf := r.asString()
		x.eval = func(t *vm.Thread) bool {
			l, r := lf(t), rf(t)
			return l == r
		}
	case *types.PtrType:
		lf := l.asPtr()
		rf := r.asPtr()
		x.eval = func(t *vm.Thread) bool {
			l, r := lf(t), rf(t)
			return l == r
		}
	case *types.FuncType:
		lf := l.asFunc()
		rf := r.asFunc()
		x.eval = func(t *vm.Thread) bool {
			l, r := lf(t), rf(t)
			return l == r
		}
	case *types.MapType:
		lf := l.asMap()
		rf := r.asMap()
		x.eval = func(t *vm.Thread) bool {
			l, r := lf(t), rf(t)
			return l == r
		}
	default:
		logger.Panic().
			Str("type", fmt.Sprintf("%v", l.ExprType)).
			Str("pos", fmt.Sprintf("%v", x.pos)).
			Msg("unexpected type")
	}
}

// genBinOpLss sets x.eval to a function that will return the result of 'l < r'.
func (x *Expr) genBinOpLss(l, r *Expr) {
	switch l.ExprType.Lit().(type) {
	case *types.UintType:
		lf := l.asUint()
		rf := r.asUint()
		x.eval = func(t *vm.Thread) bool {
			l, r := lf(t), rf(t)
			return l < r
		}
	case *types.IntType:
		lf := l.asInt()
		rf := r.asInt()
		x.eval = func(t *vm.Thread) bool {
			l, r := lf(t), rf(t)
			return l < r
		}
	case *types.FloatType:
		lf := l.asFloat()
		rf := r.asFloat()
		x.eval = func(t *vm.Thread) bool {
			l, r := lf(t), rf(t)
			return l < r
		}
	case *types.IdealIntType:
		l := l.asIdealInt()()
		r := r.asIdealInt()()
		val := l.Cmp(r) < 0
		x.eval = func(t *vm.Thread) bool {
			return val
		}
	case *types.IdealFloatType:
		l := l.asIdealFloat()()
		r := r.asIdealFloat()()
		val := l.Cmp(r) < 0
		x.eval = func(t *vm.Thread) bool {
			return val
		}
	case *types.StringType:
		lf := l.asString()
		rf := r.asString()
		x.eval = func(t *vm.Thread) bool {
			l, r := lf(t), rf(t)
			return l < r
		}
	default:
		logger.Panic().
			Str("type", fmt.Sprintf("%v", l.ExprType)).
			Str("pos", fmt.Sprintf("%v", x.pos)).
			Msg("unexpected type")
	}
}

// genBinOpGtr sets x.eval to a function that will return the result of 'l > r'.
func (x *Expr) genBinOpGtr(l, r *Expr) {
	switch l.ExprType.Lit().(type) {
	case *types.UintType:
		lf := l.asUint()
		rf := r.asUint()
		x.eval = func(t *vm.Thread) bool {
			l, r := lf(t), rf(t)
			return l > r
		}
	case *types.IntType:
		lf := l.asInt()
		rf := r.asInt()
		x.eval = func(t *vm.Thread) bool {
			l, r := lf(t), rf(t)
			return l > r
		}
	case *types.FloatType:
		lf := l.asFloat()
		rf := r.asFloat()
		x.eval = func(t *vm.Thread) bool {
			l, r := lf(t), rf(t)
			return l > r
		}
	case *types.IdealIntType:
		l := l.asIdealInt()()
		r := r.asIdealInt()()
		val := l.Cmp(r) > 0
		x.eval = func(t *vm.Thread) bool {
			return val
		}
	case *types.IdealFloatType:
		l := l.asIdealFloat()()
		r := r.asIdealFloat()()
		val := l.Cmp(r) > 0
		x.eval = func(t *vm.Thread) bool {
			return val
		}
	case *types.StringType:
		lf := l.asString()
		rf := r.asString()
		x.eval = func(t *vm.Thread) bool {
			l, r := lf(t), rf(t)
			return l > r
		}
	default:
		logger.Panic().
			Str("type", fmt.Sprintf("%v", l.ExprType)).
			Str("pos", fmt.Sprintf("%v", x.pos)).
			Msg("unexpected type")
	}
}

// genBinOpLeq sets x.eval to a function that will return the result of 'l <= r'.
func (x *Expr) genBinOpLeq(l, r *Expr) {
	switch l.ExprType.Lit().(type) {
	case *types.UintType:
		lf := l.asUint()
		rf := r.asUint()
		x.eval = func(t *vm.Thread) bool {
			l, r := lf(t), rf(t)
			return l <= r
		}
	case *types.IntType:
		lf := l.asInt()
		rf := r.asInt()
		x.eval = func(t *vm.Thread) bool {
			l, r := lf(t), rf(t)
			return l <= r
		}
	case *types.FloatType:
		lf := l.asFloat()
		rf := r.asFloat()
		x.eval = func(t *vm.Thread) bool {
			l, r := lf(t), rf(t)
			return l <= r
		}
	case *types.IdealIntType:
		l := l.asIdealInt()()
		r := r.asIdealInt()()
		val := l.Cmp(r) <= 0
		x.eval = func(t *vm.Thread) bool { return val }
	case *types.IdealFloatType:
		l := l.asIdealFloat()()
		r := r.asIdealFloat()()
		val := l.Cmp(r) <= 0
		x.eval = func(t *vm.Thread) bool { return val }
	case *types.StringType:
		lf := l.asString()
		rf := r.asString()
		x.eval = func(t *vm.Thread) bool {
			l, r := lf(t), rf(t)
			return l <= r
		}
	default:
		logger.Panic().
			Str("type", fmt.Sprintf("%v", l.ExprType)).
			Str("pos", fmt.Sprintf("%v", x.pos)).
			Msg("unexpected type")
	}
}

// genBinOpEql sets x.eval to a function that will return the result of 'l >= r'.
func (x *Expr) genBinOpGeq(l, r *Expr) {
	switch l.ExprType.Lit().(type) {
	case *types.UintType:
		lf := l.asUint()
		rf := r.asUint()
		x.eval = func(t *vm.Thread) bool {
			l, r := lf(t), rf(t)
			return l >= r
		}
	case *types.IntType:
		lf := l.asInt()
		rf := r.asInt()
		x.eval = func(t *vm.Thread) bool {
			l, r := lf(t), rf(t)
			return l >= r
		}
	case *types.FloatType:
		lf := l.asFloat()
		rf := r.asFloat()
		x.eval = func(t *vm.Thread) bool {
			l, r := lf(t), rf(t)
			return l >= r
		}
	case *types.IdealIntType:
		l := l.asIdealInt()()
		r := r.asIdealInt()()
		val := l.Cmp(r) >= 0
		x.eval = func(t *vm.Thread) bool {
			return val
		}
	case *types.IdealFloatType:
		l := l.asIdealFloat()()
		r := r.asIdealFloat()()
		val := l.Cmp(r) >= 0
		x.eval = func(t *vm.Thread) bool {
			return val
		}
	case *types.StringType:
		lf := l.asString()
		rf := r.asString()
		x.eval = func(t *vm.Thread) bool {
			l, r := lf(t), rf(t)
			return l >= r
		}
	default:
		logger.Panic().
			Str("type", fmt.Sprintf("%v", l.ExprType)).
			Str("pos", fmt.Sprintf("%v", x.pos)).
			Msg("unexpected type")
	}
}

// genBinOpNeq sets x.eval to a function that will return the result of 'l != r'.
func (x *Expr) genBinOpNeq(l, r *Expr) {
	switch l.ExprType.Lit().(type) {
	case *types.BoolType:
		lf := l.asBool()
		rf := r.asBool()
		x.eval = func(t *vm.Thread) bool {
			l, r := lf(t), rf(t)
			return l != r
		}
	case *types.UintType:
		lf := l.asUint()
		rf := r.asUint()
		x.eval = func(t *vm.Thread) bool {
			l, r := lf(t), rf(t)
			return l != r
		}
	case *types.IntType:
		lf := l.asInt()
		rf := r.asInt()
		x.eval = func(t *vm.Thread) bool {
			l, r := lf(t), rf(t)
			return l != r
		}
	case *types.FloatType:
		lf := l.asFloat()
		rf := r.asFloat()
		x.eval = func(t *vm.Thread) bool {
			l, r := lf(t), rf(t)
			return l != r
		}
	case *types.IdealIntType:
		l := l.asIdealInt()()
		r := r.asIdealInt()()
		val := l.Cmp(r) != 0
		x.eval = func(t *vm.Thread) bool {
			return val
		}
	case *types.IdealFloatType:
		l := l.asIdealFloat()()
		r := r.asIdealFloat()()
		val := l.Cmp(r) != 0
		x.eval = func(t *vm.Thread) bool {
			return val
		}
	case *types.StringType:
		lf := l.asString()
		rf := r.asString()
		x.eval = func(t *vm.Thread) bool {
			l, r := lf(t), rf(t)
			return l != r
		}
	case *types.PtrType:
		lf := l.asPtr()
		rf := r.asPtr()
		x.eval = func(t *vm.Thread) bool {
			l, r := lf(t), rf(t)
			return l != r
		}
	case *types.FuncType:
		lf := l.asFunc()
		rf := r.asFunc()
		x.eval = func(t *vm.Thread) bool {
			l, r := lf(t), rf(t)
			return l != r
		}
	case *types.MapType:
		lf := l.asMap()
		rf := r.asMap()
		x.eval = func(t *vm.Thread) bool {
			l, r := lf(t), rf(t)
			return l != r
		}
	default:
		logger.Panic().
			Str("type", fmt.Sprintf("%v", l.ExprType)).
			Str("pos", fmt.Sprintf("%v", x.pos)).
			Msg("unexpected type")
	}
}

////////////////////////////////////////////////////////////////////////////////

// genUnaryOpNeg sets x.eval to a function that will return the result of '-v'.
func (x *Expr) genUnaryOpNeg(v *Expr) {
	switch x.ExprType.Lit().(type) {
	case *types.UintType:
		vf := v.asUint()
		x.eval = func(t *vm.Thread) uint64 {
			v := vf(t)
			return -v
		}
	case *types.IntType:
		vf := v.asInt()
		x.eval = func(t *vm.Thread) int64 {
			v := vf(t)
			return -v
		}
	case *types.FloatType:
		vf := v.asFloat()
		x.eval = func(t *vm.Thread) float64 {
			v := vf(t)
			return -v
		}
	case *types.IdealIntType:
		val := v.asIdealInt()()
		val.Neg(val)
		x.eval = func() *big.Int {
			return val
		}
	case *types.IdealFloatType:
		val := v.asIdealFloat()()
		val.Neg(val)
		x.eval = func() *big.Rat {
			return val
		}
	default:
		logger.Panic().
			Str("type", fmt.Sprintf("%v", x.ExprType)).
			Str("pos", fmt.Sprintf("%v", x.pos)).
			Msg("unexpected type")
	}
}

// genUnaryOpNot sets x.eval to a function that will return the result of '!v'.
func (x *Expr) genUnaryOpNot(v *Expr) {
	switch x.ExprType.Lit().(type) {
	case *types.BoolType:
		vf := v.asBool()
		x.eval = func(t *vm.Thread) bool {
			v := vf(t)
			return !v
		}
	default:
		logger.Panic().
			Str("type", fmt.Sprintf("%v", x.ExprType)).
			Str("pos", fmt.Sprintf("%v", x.pos)).
			Msg("unexpected type")
	}
}

// genUnaryOpXor sets x.eval to a function that will return the result of '^v'.
func (x *Expr) genUnaryOpXor(v *Expr) {
	switch x.ExprType.Lit().(type) {
	case *types.UintType:
		vf := v.asUint()
		x.eval = func(t *vm.Thread) uint64 {
			v := vf(t)
			return ^v
		}
	case *types.IntType:
		vf := v.asInt()
		x.eval = func(t *vm.Thread) int64 {
			v := vf(t)
			return ^v
		}
	case *types.IdealIntType:
		val := v.asIdealInt()()
		val.Not(val)
		x.eval = func() *big.Int {
			return val
		}
	default:
		logger.Panic().
			Str("type", fmt.Sprintf("%v", x.ExprType)).
			Str("pos", fmt.Sprintf("%v", x.pos)).
			Msg("unexpected type")
	}
}

// genUnaryOpXor sets x.eval to a function that will return the result of '*v'.
func (x *Expr) genUnaryOpAddr(v *Expr) {
	vf := v.evalAddr
	x.eval = func(t *vm.Thread) vm.Value {
		return vf(t)
	}
}

// genUnaryOpXor sets x.eval to a function that will return the result of '<-v'.
func (x *Expr) genUnaryOpReceive(v *Expr) {
	panic("Unary op <- not implemented")
}

////////////////////////////////////////////////////////////////////////////////

// genConstant sets x.eval to a function that will return the value of v.
func (x *Expr) genConstant(v vm.Value) {
	switch x.ExprType.Lit().(type) {
	case *types.BoolType:
		x.eval = func(t *vm.Thread) bool {
			return v.(values.BoolValue).Get(t)
		}
	case *types.UintType:
		x.eval = func(t *vm.Thread) uint64 {
			return v.(values.UintValue).Get(t)
		}
	case *types.IntType:
		x.eval = func(t *vm.Thread) int64 {
			return v.(values.IntValue).Get(t)
		}
	case *types.IdealIntType:
		val := v.(values.IdealIntValue).Get()
		x.eval = func() *big.Int {
			return val
		}
	case *types.FloatType:
		x.eval = func(t *vm.Thread) float64 {
			return v.(values.FloatValue).Get(t)
		}
	case *types.IdealFloatType:
		val := v.(values.IdealFloatValue).Get()
		x.eval = func() *big.Rat {
			return val
		}
	case *types.StringType:
		x.eval = func(t *vm.Thread) string {
			return v.(values.StringValue).Get(t)
		}
	case *types.ArrayType:
		x.eval = func(t *vm.Thread) values.ArrayValue {
			return v.(values.ArrayValue).Get(t)
		}
	case *types.StructType:
		x.eval = func(t *vm.Thread) values.StructValue {
			return v.(values.StructValue).Get(t)
		}
	case *types.PtrType:
		x.eval = func(t *vm.Thread) vm.Value {
			return v.(values.PtrValue).Get(t)
		}
	case *types.FuncType:
		x.eval = func(t *vm.Thread) values.Func {
			return v.(values.FuncValue).Get(t)
		}
	case *types.SliceType:
		x.eval = func(t *vm.Thread) values.Slice {
			return v.(values.SliceValue).Get(t)
		}
	case *types.MapType:
		x.eval = func(t *vm.Thread) values.Map {
			return v.(values.MapValue).Get(t)
		}
	default:
		logger.Panic().
			Str("type", fmt.Sprintf("%v", x.ExprType)).
			Str("pos", fmt.Sprintf("%v", x.pos)).
			Msg("unexpected constant type")
	}
}

// genIdentOp sets x.eval to a function that will return the value of
// the variable specified by level and index in the current frame.
func (x *Expr) genIdentOp(level, index int) {
	x.evalAddr = func(t *vm.Thread) vm.Value {
		return t.Frame.Get(level, index)
	}

	switch x.ExprType.Lit().(type) {
	case *types.BoolType:
		x.eval = func(t *vm.Thread) bool {
			return t.Frame.Get(level, index).(values.BoolValue).Get(t)
		}
	case *types.UintType:
		x.eval = func(t *vm.Thread) uint64 {
			return t.Frame.Get(level, index).(values.UintValue).Get(t)
		}
	case *types.IntType:
		x.eval = func(t *vm.Thread) int64 {
			return t.Frame.Get(level, index).(values.IntValue).Get(t)
		}
	case *types.FloatType:
		x.eval = func(t *vm.Thread) float64 {
			return t.Frame.Get(level, index).(values.FloatValue).Get(t)
		}
	case *types.StringType:
		x.eval = func(t *vm.Thread) string {
			return t.Frame.Get(level, index).(values.StringValue).Get(t)
		}
	case *types.ArrayType:
		x.eval = func(t *vm.Thread) values.ArrayValue {
			return t.Frame.Get(level, index).(values.ArrayValue).Get(t)
		}
	case *types.StructType:
		x.eval = func(t *vm.Thread) values.StructValue {
			return t.Frame.Get(level, index).(values.StructValue).Get(t)
		}
	case *types.PtrType:
		x.eval = func(t *vm.Thread) vm.Value {
			return t.Frame.Get(level, index).(values.PtrValue).Get(t)
		}
	case *types.FuncType:
		x.eval = func(t *vm.Thread) values.Func {
			return t.Frame.Get(level, index).(values.FuncValue).Get(t)
		}
	case *types.SliceType:
		x.eval = func(t *vm.Thread) values.Slice {
			return t.Frame.Get(level, index).(values.SliceValue).Get(t)
		}
	case *types.MapType:
		x.eval = func(t *vm.Thread) values.Map {
			return t.Frame.Get(level, index).(values.MapValue).Get(t)
		}
	default:
		logger.Panic().
			Str("type", fmt.Sprintf("%v", x.ExprType)).
			Str("pos", fmt.Sprintf("%v", x.pos)).
			Msg("unexpected identifier type")
	}
}

func (x *Expr) genValue(vf func(*vm.Thread) vm.Value) {
	x.evalAddr = vf
	switch x.ExprType.Lit().(type) {
	case *types.BoolType:
		x.eval = func(t *vm.Thread) bool {
			return vf(t).(values.BoolValue).Get(t)
		}
	case *types.UintType:
		x.eval = func(t *vm.Thread) uint64 {
			return vf(t).(values.UintValue).Get(t)
		}
	case *types.IntType:
		x.eval = func(t *vm.Thread) int64 {
			return vf(t).(values.IntValue).Get(t)
		}
	case *types.FloatType:
		x.eval = func(t *vm.Thread) float64 {
			return vf(t).(values.FloatValue).Get(t)
		}
	case *types.StringType:
		x.eval = func(t *vm.Thread) string {
			return vf(t).(values.StringValue).Get(t)
		}
	case *types.ArrayType:
		x.eval = func(t *vm.Thread) values.ArrayValue {
			return vf(t).(values.ArrayValue).Get(t)
		}
	case *types.StructType:
		x.eval = func(t *vm.Thread) values.StructValue {
			return vf(t).(values.StructValue).Get(t)
		}
	case *types.PtrType:
		x.eval = func(t *vm.Thread) vm.Value {
			return vf(t).(values.PtrValue).Get(t)
		}
	case *types.FuncType:
		x.eval = func(t *vm.Thread) values.Func {
			return vf(t).(values.FuncValue).Get(t)
		}
	case *types.SliceType:
		x.eval = func(t *vm.Thread) values.Slice {
			return vf(t).(values.SliceValue).Get(t)
		}
	case *types.MapType:
		x.eval = func(t *vm.Thread) values.Map {
			return vf(t).(values.MapValue).Get(t)
		}
	default:
		logger.Panic().
			Str("type", fmt.Sprintf("%v", x.ExprType)).
			Str("pos", fmt.Sprintf("%v", x.pos)).
			Msg("unexpected result type")
	}
}

// genFuncCall sets x.eval to a function that will execute 'call' and return its result.
func (x *Expr) genFuncCall(call func(t *vm.Thread) []vm.Value) {
	x.Exec = func(t *vm.Thread) { call(t) }
	switch x.ExprType.Lit().(type) {
	case *types.BoolType:
		x.eval = func(t *vm.Thread) bool {
			return call(t)[0].(values.BoolValue).Get(t)
		}
	case *types.UintType:
		x.eval = func(t *vm.Thread) uint64 {
			return call(t)[0].(values.UintValue).Get(t)
		}
	case *types.IntType:
		x.eval = func(t *vm.Thread) int64 {
			return call(t)[0].(values.IntValue).Get(t)
		}
	case *types.FloatType:
		x.eval = func(t *vm.Thread) float64 {
			return call(t)[0].(values.FloatValue).Get(t)
		}
	case *types.StringType:
		x.eval = func(t *vm.Thread) string {
			return call(t)[0].(values.StringValue).Get(t)
		}
	case *types.ArrayType:
		x.eval = func(t *vm.Thread) values.ArrayValue {
			return call(t)[0].(values.ArrayValue).Get(t)
		}
	case *types.StructType:
		x.eval = func(t *vm.Thread) values.StructValue {
			return call(t)[0].(values.StructValue).Get(t)
		}
	case *types.PtrType:
		x.eval = func(t *vm.Thread) vm.Value {
			return call(t)[0].(values.PtrValue).Get(t)
		}
	case *types.FuncType:
		x.eval = func(t *vm.Thread) values.Func {
			return call(t)[0].(values.FuncValue).Get(t)
		}
	case *types.SliceType:
		x.eval = func(t *vm.Thread) values.Slice {
			return call(t)[0].(values.SliceValue).Get(t)
		}
	case *types.MapType:
		x.eval = func(t *vm.Thread) values.Map {
			return call(t)[0].(values.MapValue).Get(t)
		}
	case *types.MultiType:
		x.eval = func(t *vm.Thread) []vm.Value {
			return call(t)
		}
	default:
		logger.Panic().
			Str("type", fmt.Sprintf("%v", x.ExprType)).
			Str("pos", fmt.Sprintf("%v", x.pos)).
			Msg("unexpected result type")
	}
}
