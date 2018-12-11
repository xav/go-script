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
				}

				elts[i] = elt
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
		return &values.SliceV{
			Slice: values.Slice{
				Base: base,
				Len:  int64(sz),
				Cap:  int64(sz),
			},
		}
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
				}
				elts[i] = elt
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
	// Save the original types of l.Type and r.Type for error messages.
	origlt := l.ExprType
	origrt := r.ExprType

	if op != token.SHL && op != token.SHR {
		// Except in shift expressions, if one operand has numeric type and the other operand is an ideal number,
		// the ideal number is converted to match the type of the other operand.
		if (l.ExprType.IsInteger() || l.ExprType.IsFloat()) && !l.ExprType.IsIdeal() && r.ExprType.IsIdeal() {
			r = r.resolveIdeal(l.ExprType)
		} else if (r.ExprType.IsInteger() || r.ExprType.IsFloat()) && l.ExprType.IsIdeal() && !r.ExprType.IsIdeal() {
			l = l.resolveIdeal(r.ExprType)
		}
		if l == nil || r == nil {
			return nil
		}

		// Except in shift expressions, if both operands are ideal numbers and one is an ideal float,
		// the other is converted to ideal float.
		if l.ExprType.IsIdeal() && r.ExprType.IsIdeal() {
			if l.ExprType.IsInteger() && r.ExprType.IsFloat() {
				l = l.resolveIdeal(r.ExprType)
			} else if l.ExprType.IsFloat() && r.ExprType.IsInteger() {
				r = r.resolveIdeal(l.ExprType)
			}
			if l == nil || r == nil {
				return nil
			}
		}

	}

	// Useful type predicates
	var (
		compat   = func() bool { return l.ExprType.Compat(r.ExprType, false) }
		integers = func() bool { return l.ExprType.IsInteger() && r.ExprType.IsInteger() }
		floats   = func() bool { return l.ExprType.IsFloat() && r.ExprType.IsFloat() }
		strings  = func() bool { return l.ExprType == builtins.StringType && r.ExprType == builtins.StringType }
		booleans = func() bool { return l.ExprType.IsBoolean() && r.ExprType.IsBoolean() }
	)

	// Type check
	var t vm.Type
	switch op {
	case token.ADD:
		if !compat() || (!integers() && !floats() && !strings()) {
			xi.errorOpTypes(op, origlt, origrt)
			return nil
		}
		t = l.ExprType
	case token.SUB, token.MUL, token.QUO:
		if !compat() || (!integers() && !floats()) {
			xi.errorOpTypes(op, origlt, origrt)
			return nil
		}
		t = l.ExprType
	case token.REM, token.AND, token.OR, token.XOR, token.AND_NOT:
		if !compat() || !integers() {
			xi.errorOpTypes(op, origlt, origrt)
			return nil
		}
		t = l.ExprType
	case token.SHL, token.SHR:
		if !l.ExprType.IsInteger() || !(r.ExprType.IsInteger() || r.ExprType.IsIdeal()) {
			xi.errorOpTypes(op, origlt, origrt)
			return nil
		}

		// The right operand in a shift operation must be always be of unsigned integer type or an ideal
		// number that can be safely converted into an unsigned integer type.
		if r.ExprType.IsIdeal() {
			r2 := r.resolveIdeal(builtins.UintType)
			if r2 == nil {
				return nil
			}

			// If the left operand is not ideal, convert the right to not ideal.
			if !l.ExprType.IsIdeal() {
				r = r2
			}

			// If both are ideal, but the right side isn't an ideal int, convert it to simplify things.
			if l.ExprType.IsIdeal() && !r.ExprType.IsInteger() {
				r = r.resolveIdeal(IdealIntType)
				if r == nil {
					logger.Panic().Msgf("conversion to uint succeeded, but conversion to ideal int failed")
				}
			}
		} else if _, ok := r.ExprType.Lit().(*types.UintType); !ok {
			xi.error("right operand of shift must be unsigned")
			return nil
		}

		if l.ExprType.IsIdeal() && !r.ExprType.IsIdeal() {
			l = l.resolveIdeal(builtins.IntType)
			if l == nil {
				return nil
			}
		}

		// At this point, we should have one of three cases:
		// 1) uint SHIFT uint
		// 2) int SHIFT uint
		// 3) ideal int SHIFT ideal int

		t = l.ExprType

	case token.LOR, token.LAND:
		if !booleans() {
			return nil
		}
		t = builtins.BoolType
	case token.ARROW:
		// The operands in channel sends differ in type: one is always a channel and the other
		// is a variable or value of the channel's element type.
		logger.Panic().Msg("Binary op <- not implemented")
		t = builtins.BoolType
	case token.LSS, token.GTR, token.LEQ, token.GEQ:
		if !compat() || (!integers() && !floats() && !strings()) {
			xi.errorOpTypes(op, origlt, origrt)
			return nil
		}
		t = builtins.BoolType
	case token.EQL, token.NEQ:
		if !compat() {
			xi.errorOpTypes(op, origlt, origrt)
			return nil
		}
		t = builtins.BoolType
	default:
		logger.Panic().
			Str("operator", fmt.Sprintf("%v", op)).
			Msg("unknown binary operator")
	}

	desc, ok := binOpDescs[op]
	if !ok {
		desc = op.String() + " expression"
		binOpDescs[op] = desc
	}

	// Check for ideal divide by zero
	if op == token.QUO || op == token.REM {
		if r.ExprType.IsIdeal() {
			if (r.ExprType.IsInteger() && r.asIdealInt()().Sign() == 0) || (r.ExprType.IsFloat() && r.asIdealFloat()().Sign() == 0) {
				xi.error("divide by zero")
				return nil
			}
		}
	}

	//////////////////////////////////////
	// Compile
	expr := xi.newExpr(t, desc)
	switch op {
	case token.ADD: // +
		expr.genBinOpAdd(l, r)
	case token.SUB: // -
		expr.genBinOpSub(l, r)
	case token.MUL: // *
		expr.genBinOpMul(l, r)
	case token.QUO: // /
		expr.genBinOpQuo(l, r)
	case token.REM: // %
		expr.genBinOpRem(l, r)

	case token.AND: // &
		expr.genBinOpAnd(l, r)
	case token.OR: // |
		expr.genBinOpOr(l, r)
	case token.XOR: // ^
		expr.genBinOpXor(l, r)
	case token.AND_NOT: // &^
		expr.genBinOpAndNot(l, r)
	case token.SHL: // <<
		if l.ExprType.IsIdeal() {
			lv := l.asIdealInt()()
			rv := r.asIdealInt()()
			const maxShift = 99999
			if rv.Cmp(big.NewInt(maxShift)) > 0 {
				xi.error("left shift by %v; exceeds implementation limit of %v", rv, maxShift)
				expr.ExprType = nil
				return nil
			}
			val := new(big.Int).Lsh(lv, uint(rv.Int64()))
			expr.eval = func() *big.Int {
				return val
			}
		} else {
			expr.genBinOpShl(l, r)
		}
	case token.SHR: // >>
		if l.ExprType.IsIdeal() {
			lv := l.asIdealInt()()
			rv := r.asIdealInt()()
			val := new(big.Int).Rsh(lv, uint(rv.Int64()))
			expr.eval = func() *big.Int {
				return val
			}
		} else {
			expr.genBinOpShr(l, r)
		}

	case token.LAND: // &&
		expr.genBinOpLogAnd(l, r)
	case token.LOR: // ||
		expr.genBinOpLogOr(l, r)

	case token.ARROW: // <-
		panic("Binary op <- not implemented")

	case token.EQL: // ==
		expr.genBinOpEql(l, r)
	case token.LSS: // <
		expr.genBinOpLss(l, r)
	case token.GTR: // >
		expr.genBinOpGtr(l, r)
	case token.LEQ: // <=
		expr.genBinOpLeq(l, r)
	case token.GEQ: // >=
		expr.genBinOpGeq(l, r)
	case token.NEQ: // !=
		expr.genBinOpNeq(l, r)
	default:
		logger.Panic().
			Str("operator", fmt.Sprintf("%v", op)).
			Msgf("Compilation of binary op %v not implemented", op)

	}

	return expr
}

func (xi *ExprInfo) compileBuiltinCallExpr(b *context.Block, ft *types.FuncType, as []*Expr) *Expr {
	panic("NOT IMPLEMENTED")
}

func (xi *ExprInfo) compileCallExpr(b *context.Block, l *Expr, as []*Expr) *Expr {
	// TODO: variadic functions

	// Type check
	lt, ok := l.ExprType.Lit().(*types.FuncType)
	if !ok {
		xi.error("cannot call non-function type %v", l.ExprType)
		return nil
	}

	// The arguments must be single-valued expressions assignment compatible with the parameters of F,
	// or a single multi-valued expression.
	nin := len(lt.In)
	assign := xi.compileAssign(xi.pos, b, types.NewMultiType(lt.In), as, "function call", "argument")
	if assign == nil {
		return nil
	}

	var t vm.Type
	nout := len(lt.Out)
	switch nout {
	case 0:
		t = types.EmptyType
	case 1:
		t = lt.Out[0]
	default:
		t = types.NewMultiType(lt.Out)
	}
	expr := xi.newExpr(t, "function call")

	// Gather argument and out types to initialize frame variables
	vts := make([]vm.Type, nin+nout)
	copy(vts, lt.In)
	copy(vts[nin:], lt.Out)

	// Compile
	lf := l.asFunc()
	call := func(t *vm.Thread) []vm.Value {
		fun := lf(t)
		fr := fun.NewFrame()
		for i, t := range vts {
			fr.Vars[i] = t.Zero()
		}
		assign(values.MultiV(fr.Vars[0:nin]), t)

		oldf := t.Frame
		t.Frame = fr
		fun.Call(t)
		t.Frame = oldf

		return fr.Vars[nin : nin+nout]
	}
	expr.genFuncCall(call)
	return expr
}

func (xi *ExprInfo) compileIdent(b *context.Block, constant bool, callCtx bool, name string) *Expr {
	bl, level, def := b.Lookup(name)
	if def == nil {
		xi.error("%s: undefined", name)
		return nil
	}

	switch def := def.(type) {
	case *types.Constant:
		expr := xi.newExpr(def.Type, "constant")
		if ft, ok := def.Type.(*types.FuncType); ok && ft.Builtin != "" {
			if !callCtx {
				xi.error("built-in function %s cannot be used as a value", ft.Builtin)
				return nil
			}
		} else {
			expr.genConstant(def.Value)
		}
		return expr

	case *types.Variable:
		if constant {
			xi.error("variable %s used in constant expression", name)
			return nil
		}
		if bl.Global {
			return xi.compileGlobalVariable(def)
		}
		return xi.compileVariable(level, def)

	case vm.Type:
		if callCtx {
			return xi.exprFromType(def)
		}
		xi.error("type %v used as expression", name)
		return nil

	case *context.PkgIdent:
		return xi.compilePackageImport(name, def, constant, true)
	}

	logger.Panic().
		Str("symbol", name).
		Str("type", fmt.Sprintf("%T", def)).
		Msg("symbol has unknown type")
	panic("unreachable")
}

func (xi *ExprInfo) compileIndexExpr(l, r *Expr) *Expr {
	panic("NOT IMPLEMENTED")
}

func (xi *ExprInfo) compileSliceExpr(arr, lo, hi *Expr) *Expr {
	arr = arr.derefArray()

	// Type check object
	var at vm.Type
	var maxIndex int64 = -1
	switch lt := arr.ExprType.Lit().(type) {
	case *types.ArrayType:
		at = types.NewSliceType(lt.Elem)
		maxIndex = lt.Len
	case *types.SliceType:
		at = lt
	case *types.StringType:
		at = lt
	default:
		xi.error("cannot slice %v", arr.ExprType)
		return nil
	}

	// Type check index and convert to int
	lo = lo.convertToInt(maxIndex, "slice", "slice")
	hi = hi.convertToInt(maxIndex, "slice", "slice")
	if lo == nil || hi == nil {
		return nil
	}

	//////////////////////////////////////
	// Compile

	expr := xi.newExpr(at, "slice expression")
	lof := lo.asInt()
	hif := hi.asInt()

	switch lt := arr.ExprType.Lit().(type) {
	case *types.ArrayType:
		arrf := arr.asArray()
		bound := lt.Len
		expr.eval = func(t *vm.Thread) values.Slice {
			arr, lo, hi := arrf(t), lof(t), hif(t)
			if lo > hi || hi > bound || lo < 0 {
				t.Abort(vm.SliceError{
					Lo:  lo,
					Hi:  hi,
					Cap: bound,
				})
			}
			return values.Slice{
				Base: arr.Sub(lo, bound-lo),
				Len:  hi - lo,
				Cap:  bound - lo,
			}
		}

	case *types.SliceType:
		arrf := arr.asSlice()
		expr.eval = func(t *vm.Thread) values.Slice {
			arr, lo, hi := arrf(t), lof(t), hif(t)
			if lo > hi || hi > arr.Cap || lo < 0 {
				t.Abort(vm.SliceError{
					Lo:  lo,
					Hi:  hi,
					Cap: arr.Cap,
				})
			}
			return values.Slice{
				Base: arr.Base.Sub(lo, arr.Cap-lo),
				Len:  hi - lo,
				Cap:  arr.Cap - lo,
			}
		}

	case *types.StringType:
		arrf := arr.asString()
		expr.eval = func(t *vm.Thread) string {
			arr, lo, hi := arrf(t), lof(t), hif(t)
			if lo > hi || hi > int64(len(arr)) || lo < 0 {
				t.Abort(vm.SliceError{
					Lo:  lo,
					Hi:  hi,
					Cap: int64(len(arr)),
				})
			}
			return arr[lo:hi]
		}

	default:
		logger.Panic().
			Str("type", fmt.Sprintf("%T", arr.ExprType.Lit())).
			Msg("unexpected left operand type")
	}

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
