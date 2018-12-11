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
	"go/token"
	"math/big"
	"strconv"
	"strings"

	"github.com/xav/go-script/builtins"
	"github.com/xav/go-script/context"
	"github.com/xav/go-script/types"
	"github.com/xav/go-script/values"
	"github.com/xav/go-script/vm"
)

var (
	unaryOpDescs = make(map[token.Token]string)
	binOpDescs   = make(map[token.Token]string)
)

// ExprInfo stores information needed to compile any expression node.
// Each expr also stores its exprInfo so further expressions can be compiled from it.
type ExprInfo struct {
	*Compiler
	pos token.Pos
}

func (xi *ExprInfo) newExpr(t vm.Type, desc string) *Expr {
	return &Expr{
		ExprInfo: xi,
		ExprType: t,
		desc:     desc,
	}
}

func (xi *ExprInfo) exprFromType(t vm.Type) *Expr {
	if t == nil {
		return nil
	}
	expr := xi.newExpr(nil, "type")
	expr.valType = t
	return expr
}

func (xi *ExprInfo) error(format string, args ...interface{}) {
	xi.errorAt(xi.pos, format, args...)
}

func (xi *ExprInfo) errorOpType(op token.Token, vt vm.Type) {
	xi.error("illegal operand type for '%v' operator\n\t%v", op, vt)
}

func (xi *ExprInfo) errorOpTypes(op token.Token, lt vm.Type, rt vm.Type) {
	xi.error("illegal operand types for '%v' operator\n\t%v\n\t%v", op, lt, rt)
}

////////////////////////////////////////////////////////////////////////////////

func (xi *ExprInfo) compileIntLit(lit string) *Expr {
	i, _ := new(big.Int).SetString(lit, 0)
	return xi.compileIdealInt(i, "integer literal")
}

func (xi *ExprInfo) compileFloatLit(lit string) *Expr {
	f, ok := new(big.Rat).SetString(lit)
	if !ok {
		logger.Panic().
			Str("pos", fmt.Sprintf("%v", xi.pos)).
			Msgf("malformed float literal %s passed parser", lit)
	}

	expr := xi.newExpr(IdealFloatType, "float literal")
	expr.eval = func() *big.Rat { return f }
	return expr
}

func (xi *ExprInfo) compileCharLit(lit string) *Expr {
	if lit[0] != '\'' {
		xi.SilentErrors++ // Caught by parser
		return nil
	}

	v, _, tail, err := strconv.UnquoteChar(lit[1:], '\'')
	if err != nil || tail != "'" {
		xi.SilentErrors++ // Caught by parser
		return nil
	}

	return xi.compileIdealInt(big.NewInt(int64(v)), "character literal")
}

func (xi *ExprInfo) compileStringLit(lit string) *Expr {
	s, err := strconv.Unquote(lit)
	if err != nil {
		xi.error("illegal string literal, %v", err)
		return nil
	}
	return xi.compileString(s)
}

func (xi *ExprInfo) compileIdealInt(i *big.Int, desc string) *Expr {
	expr := xi.newExpr(IdealIntType, desc)
	expr.eval = func() *big.Int { return i }
	return expr
}

func (xi *ExprInfo) compileCompositeLit(c *Expr, ifKeys []interface{}, vals []*Expr) *Expr {
	if c == nil {
		xi.error("invalid composite expression")
		return nil
	}

	for i, elmt := range vals {
		if elmt == nil {
			xi.error("nil argument (#%d)", i+1)
			return nil
		}
	}

	comp := xi.newExpr(c.valType, "composite literal")

	switch ty := c.valType.Lit().(type) {
	case *types.NamedType:
		comp = nil
	case *types.StructType:
		return xi.compileCompositeStructType(comp, ty, vals, ifKeys)
	case *types.ArrayType:
		return xi.compileCompositeArrayType(comp, ty, vals)
	case *types.SliceType:
		return xi.compileCompositeSliceType(comp, ty, vals)
	case *types.MapType:
		return xi.compileCompositeMapType(comp, ty, vals, ifKeys)
	default:
		comp = nil
	}

	if comp == nil {
		xi.error("composite literal not impleemented for type [%s]\n", c.valType.Lit().String())
	}
	return comp
}

func (xi *ExprInfo) compileCompositeStructType(comp *Expr, ty *types.StructType, elts []*Expr, ifKeys []interface{}) *Expr {
	keys := ifKeys
	sz := len(ty.Elems)

	if len(elts) > sz {
		xi.error("given too many elements (%d) (expected %d)", len(elts), len(ty.Elems))
		return nil
	}
	if len(elts) < sz {
		xi.error("too few values in struct initializer")
		return nil
	}
	if len(elts) != len(keys) {
		if len(keys) > 0 && len(keys) < len(elts) {
			xi.error("mixture of field:value and value initializers")
			return nil
		}
	}

	checkAndApplyConversion := func(elts []*Expr) bool {
		for i := 0; i < sz; i++ {
			if !ty.Elems[i].Type.IsIdeal() && elts[i].ExprType.IsIdeal() {
				elt := elts[i].resolveIdeal(ty.Elems[i].Type)
				if elt == nil {
					xi.error("cannot convert literal #%d (type %s) to type %s", i+1, elts[i].ExprType.String(), ty.Elems[i].Type.String())
					return false
				} else {
					elts[i] = elt
				}
			}
		}
		return true
	}

	var evalFct func(t *vm.Thread) vm.Value
	if len(keys) > 0 {
		if _, ok := keys[0].(string); !ok {
			xi.error("invalid key type '%T' (expected string)", keys[0])
			return nil
		}
		if len(ty.Elems) > len(keys) {
			xi.error("too few values in struct initializer")
			return nil
		}
		if len(ty.Elems) < len(keys) {
			xi.error("too many values in struct initializer")
			return nil
		}

		indices := make([]int, len(keys))
		for i, iv := range keys {
			name := iv.(string)
			indices[i] = -1
			for jj, jfield := range ty.Elems {
				if jfield.Name == name {
					indices[i] = jj
					break
				}
			}
			if indices[i] == -1 {
				xi.error("unknown field %s", name)
				return nil
			}
		}

		tmp := append([]*Expr{}, elts...)
		// reorder expressions to match struct's fields order
		for idx, ival := range indices {
			tmp[ival] = elts[idx]
		}
		if !checkAndApplyConversion(tmp) {
			return nil
		}
		elts = tmp

		evalFct = func(t *vm.Thread) vm.Value {
			out := ty.Zero().(values.StructValue)
			for i := 0; i < sz; i++ {
				out.Field(t, i).Assign(t, elts[i].asValue()(t))
			}
			return out
		}
	} else {
		if !checkAndApplyConversion(elts) {
			return nil
		}

		evalFct = func(t *vm.Thread) vm.Value {
			out := ty.Zero().(values.StructValue)
			for i := 0; i < sz; i++ {
				out.Field(t, i).Assign(t, elts[i].asValue()(t))
			}
			return out
		}
	}

	comp.genValue(evalFct)
	return comp
}

func (xi *ExprInfo) compileCompositeArrayType(comp *Expr, ty *types.ArrayType, elts []*Expr) *Expr {
	sz := len(elts)
	if int64(sz) > ty.Len {
		xi.error("given too many elements (%d) (expected %d)", sz, ty.Len)
		return nil
	}
	if !xi.massageLitIdeal(ty.Elem, elts) {
		return nil
	}
	evalFct := func(t *vm.Thread) vm.Value {
		base := ty.Zero().(values.ArrayValue)
		for i := 0; i < sz; i++ {
			base.Elem(t, int64(i)).Assign(t, elts[i].asValue()(t))
		}
		return base
	}
	comp.genValue(evalFct)
	return comp
}

func (xi *ExprInfo) compileCompositeSliceType(comp *Expr, ty *types.SliceType, elts []*Expr) *Expr {
	sz := len(elts)
	if !xi.massageLitIdeal(ty.Elem, elts) {
		return nil
	}
	arrt := types.NewArrayType(int64(sz), ty.Elem)
	evalFct := func(t *vm.Thread) vm.Value {
		base := arrt.Zero().(values.ArrayValue)
		for i := 0; i < sz; i++ {
			base.Elem(t, int64(i)).Assign(t, elts[i].asValue()(t))
		}
		return &values.SliceV{values.Slice{base, int64(sz), int64(sz)}}
	}
	comp.genValue(evalFct)
	return comp
}

func (xi *ExprInfo) compileCompositeMapType(comp *Expr, ty *types.MapType, elts []*Expr, ifKeys []interface{}) *Expr {
	if len(elts) != len(ifKeys) {
		xi.error("internal logic error")
		return nil
	}

	sz := len(elts)
	keys := make([]*Expr, len(ifKeys))
	for i := 0; i < sz; i++ {
		k, ok := ifKeys[i].(*Expr)
		if !ok {
			xi.error("key #%d isn't a *Expr! (got %T)", ifKeys[i], ifKeys[i])
			return nil
		}
		keys[i] = k
	}
	if !xi.massageLitIdeal(ty.Key, keys) || !xi.massageLitIdeal(ty.Elem, elts) {
		return nil
	}
	evalFct := func(t *vm.Thread) vm.Value {
		m := EvalMap{}
		for i := 0; i < sz; i++ {
			k := keys[i].asInterface()
			v := elts[i].asValue()
			m.SetElem(t, k(t), v(t))
		}
		out := ty.Zero().(values.MapValue)
		out.Set(t, &m)
		return out
	}
	comp.genValue(evalFct)
	return comp
}

// helper function to handle literals
func (xi *ExprInfo) massageLitIdeal(ty vm.Type, elts []*Expr) bool {
	ok := true
	if !ty.IsIdeal() {
		for i := 0; i < len(elts); i++ {
			if elts[i].ExprType.IsIdeal() {
				elt := elts[i].resolveIdeal(ty)
				if elt == nil {
					xi.error("cannot convert literal %d (type %s) to type %s", i+1, elts[i].ExprType.String(), ty.String())
					return false
				} else {
					elts[i] = elt
				}
			}
		}
	}
	return ok
}

func (xi *ExprInfo) compileFuncLit(decl *types.FuncDecl, fn func(*vm.Thread) values.Func) *Expr {
	expr := xi.newExpr(decl.Type, "function literal")
	expr.eval = fn
	return expr
}

////////////////////////////////////////////////////////////////////////////////

func (xi *ExprInfo) compileGlobalVariable(v *types.Variable) *Expr {
	if v.Type == nil {
		xi.SilentErrors++ // Placeholder definition from an earlier error
		return nil
	}

	if v.Init == nil {
		v.Init = v.Type.Zero()
	}

	expr := xi.newExpr(v.Type, "variable")
	val := v.Init
	expr.genValue(func(t *vm.Thread) vm.Value {
		return val
	})
	return expr
}

func (xi *ExprInfo) compileVariable(level int, v *types.Variable) *Expr {
	if v.Type == nil {
		xi.SilentErrors++ // Placeholder definition from an earlier error
		return nil
	}

	expr := xi.newExpr(v.Type, "variable")
	expr.genIdentOp(level, v.Index)
	return expr
}

func (xi *ExprInfo) compileString(s string) *Expr {
	expr := xi.newExpr(builtins.StringType, "string literal")
	expr.eval = func(*vm.Thread) string { return s }
	return expr
}

func (xi *ExprInfo) compileStringList(list []*Expr) *Expr {
	ss := make([]string, len(list))
	for i, s := range list {
		ss[i] = s.asString()(nil)
	}
	return xi.compileString(strings.Join(ss, ""))
}

////////////////////////////////////////////////////////////////////////////////

func (xi *ExprInfo) compilePackageImport(name string, pkg *context.PkgIdent, constant, callCtx bool) *Expr {
	panic("NOT IMPLEMENTED")
}

func (xi *ExprInfo) compileBinaryExpr(op token.Token, l, r *Expr) *Expr {
	panic("NOT IMPLEMENTED")
}

func (xi *ExprInfo) compileBuiltinCallExpr(b *context.Block, ft *types.FuncType, as []*Expr) *Expr {
	panic("NOT IMPLEMENTED")
}

func (xi *ExprInfo) compileCallExpr(b *context.Block, l *Expr, as []*Expr) *Expr {
	panic("NOT IMPLEMENTED")
}

func (xi *ExprInfo) compileIdent(b *context.Block, constant bool, callCtx bool, name string) *Expr {
	panic("NOT IMPLEMENTED")
}

func (xi *ExprInfo) compileIndexExpr(l, r *Expr) *Expr {
	panic("NOT IMPLEMENTED")
}

func (xi *ExprInfo) compileSliceExpr(arr, lo, hi *Expr) *Expr {
	panic("NOT IMPLEMENTED")
}

func (xi *ExprInfo) compileKeyValueExpr(key, val *Expr) *Expr {
	panic("NOT IMPLEMENTED")
}

func (xi *ExprInfo) compileSelectorExpr(v *Expr, name string) *Expr {
	panic("NOT IMPLEMENTED")
}

func (xi *ExprInfo) compileStarExpr(v *Expr) *Expr {
	panic("NOT IMPLEMENTED")
}

func (xi *ExprInfo) compileUnaryExpr(op token.Token, v *Expr) *Expr {
	panic("NOT IMPLEMENTED")
}
