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
	"go/ast"
	"go/token"

	"github.com/xav/go-script/builtins"
	"github.com/xav/go-script/context"
	"github.com/xav/go-script/types"
	"github.com/xav/go-script/vm"
)

// ExprCompiler stores information used throughout the compilation of a single expression.
// It does not embed funcCompiler because expressions can appear at top level.
type ExprCompiler struct {
	*Compiler
	block    *context.Block // The block this expression is being compiled in.
	constant bool           // Whether this expression is used in a constant context.
}

// compile compiles an expression AST.
// callCtx should be true if this AST is in the function position of a function call node;
// it allows the returned expression to be a type or a built-in function (which otherwise result in errors).
func (xc *ExprCompiler) compile(x ast.Expr, callCtx bool) *Expr {
	switch x := x.(type) {
	// Literals
	case *ast.BasicLit:
		return xc.compileBasicLit(x)
	case *ast.CompositeLit:
		return xc.compileCompositeLit(x, callCtx)
	case *ast.FuncLit:
		return xc.compileFuncLit(x)

	// Types
	case *ast.ArrayType:
		return xc.compileTypeExpr(x, callCtx)
	case *ast.ChanType:
		return xc.compileTypeExpr(x, callCtx)
	case *ast.Ellipsis:
		return xc.compileTypeExpr(x, callCtx)
	case *ast.FuncType:
		return xc.compileTypeExpr(x, callCtx)
	case *ast.InterfaceType:
		return xc.compileTypeExpr(x, callCtx)
	case *ast.MapType:
		return xc.compileTypeExpr(x, callCtx)

	// Remaining expressions
	case *ast.BadExpr:
		xc.SilentErrors++ // Error already reported by parser
		return nil
	case *ast.BinaryExpr:
		return xc.compileBinaryExpr(x)
	case *ast.CallExpr:
		return xc.compileCallExpr(x)
	case *ast.Ident:
		return xc.compileIdent(x, callCtx)
	case *ast.IndexExpr:
		return xc.compileIndexExpr(x)
	case *ast.SliceExpr:
		return xc.compileSliceExpr(x)
	case *ast.KeyValueExpr:
		return xc.compileKeyValueExpr(x)
	case *ast.ParenExpr:
		return xc.compile(x.X, callCtx)
	case *ast.SelectorExpr:
		return xc.compileSelectorExpr(x)
	case *ast.StarExpr:
		return xc.compileStarExpr(x, callCtx)
	case *ast.StructType:
		xc.errorAt(x.Pos(), "%T expression node not implemented", x)
		return nil
	case *ast.TypeAssertExpr:
		xc.errorAt(x.Pos(), "%T expression node not implemented", x)
		return nil
	case *ast.UnaryExpr:
		return xc.compileUnaryExpr(x)
	}

	logger.Panic().
		Str("type", fmt.Sprintf("%T", x)).
		Msg("unexpected ast node type")
	panic("unreachable")
}

// Literals ////////////////////////////////////////////////////////////////////

func (xc *ExprCompiler) compileBasicLit(x *ast.BasicLit) *Expr {
	ei := &ExprInfo{
		Compiler: xc.Compiler,
		pos:      x.Pos(),
	}

	switch x.Kind {
	case token.INT:
		return ei.compileIntLit(string(x.Value))
	case token.FLOAT:
		return ei.compileFloatLit(string(x.Value))
	case token.CHAR:
		return ei.compileCharLit(string(x.Value))
	case token.STRING:
		return ei.compileStringLit(string(x.Value))
	default:
		logger.Panic().
			Str("type", fmt.Sprintf("%v", x.Kind)).
			Msg("unexpected basic literal type")
		panic("unexpected basic literal type")
	}
}

func (xc *ExprCompiler) compileCompositeLit(x *ast.CompositeLit, callCtx bool) *Expr {
	ei := &ExprInfo{
		Compiler: xc.Compiler,
		pos:      x.Pos(),
	}

	var ty *Expr
	if x.Type != nil {
		// special case for array literal: [...]T{x,y,z}
		// we only know the number of elements here, so create the array type just now.
		switch xt := x.Type.(type) {
		case *ast.ArrayType:
			if _, ok := xt.Len.(*ast.Ellipsis); ok {
				elmtType := xc.compileType(xc.block, xt.Elt)
				ty = ei.exprFromType(types.NewArrayType(int64(len(x.Elts)), elmtType))
			} else {
				ty = xc.compile(x.Type, true)
			}
		default:
			ty = xc.compile(x.Type, true)
		}
	}

	keys := make([]interface{}, 0, len(x.Elts))
	vals := make([]*Expr, len(x.Elts))
	for i := range x.Elts {
		switch x.Elts[i].(type) {
		case *ast.KeyValueExpr:
			kv := x.Elts[i].(*ast.KeyValueExpr)
			switch kk := kv.Key.(type) {
			case *ast.Ident:
				switch x.Type.(type) {
				case *ast.Ident:
					keys = append(keys, kk.Name)
				default:
					keys = append(keys, xc.compile(kv.Key, callCtx))
				}
			default:
				ekey := xc.compile(kv.Key, callCtx)
				keys = append(keys, ekey)
			}
		default:
			vals[i] = xc.compile(x.Elts[i], callCtx)
		}
	}

	return ei.compileCompositeLit(ty, keys, vals)
}

func (xc *ExprCompiler) compileFuncLit(x *ast.FuncLit) *Expr {
	ei := &ExprInfo{
		Compiler: xc.Compiler,
		pos:      x.Pos(),
	}

	decl := ei.compileFuncType(xc.block, x.Type)
	if decl == nil {
		// TODO: Try compiling the body, perhaps with dummy argument definitions
		return nil
	}

	fn := ei.compileFunc(xc.block, decl, x.Body)
	if fn == nil {
		return nil
	}

	if xc.constant {
		xc.errorAt(x.Pos(), "function literal used in constant expression")
		return nil
	}

	return ei.compileFuncLit(decl, fn)
}

// Types ///////////////////////////////////////////////////////////////////////

func (xc *ExprCompiler) compileArrayType(x ast.Expr, callCtx bool) *Expr {
	switch x.(*ast.ArrayType).Len.(type) {
	case *ast.Ellipsis:
		xc.errorAt(x.Pos(), "array literal with ellipsis is not implemented")
		return nil
	default:
		// TODO: Use a multi-type case
		return xc.compileTypeExpr(x, callCtx)
	}
}

func (xc *ExprCompiler) compileTypeExpr(x ast.Expr, callCtx bool) *Expr {
	if !callCtx {
		xc.errorAt(x.Pos(), "type used as expression")
		return nil
	}

	ei := &ExprInfo{
		Compiler: xc.Compiler,
		pos:      x.Pos(),
	}

	return ei.exprFromType(xc.compileType(xc.block, x))
}

// Remaining expressions ///////////////////////////////////////////////////////

func (xc *ExprCompiler) compileBinaryExpr(x *ast.BinaryExpr) *Expr {
	l := xc.compile(x.X, false)
	r := xc.compile(x.Y, false)
	if l == nil || r == nil {
		return nil
	}

	ei := &ExprInfo{
		Compiler: xc.Compiler,
		pos:      x.Pos(),
	}

	return ei.compileBinaryExpr(x.Op, l, r)
}

func (xc *ExprCompiler) compileCallExpr(x *ast.CallExpr) *Expr {
	l := xc.compile(x.Fun, true)
	args := make([]*Expr, len(x.Args))

	bad := false
	for i, arg := range x.Args {
		if i == 0 && l != nil && (l.ExprType == vm.Type(builtins.MakeType) || l.ExprType == vm.Type(builtins.NewType)) {
			eiArg := &ExprInfo{
				Compiler: xc.Compiler,
				pos:      arg.Pos(),
			}
			args[i] = eiArg.exprFromType(xc.compileType(xc.block, arg))
		} else {
			args[i] = xc.compile(arg, false)
		}

		if args[i] == nil {
			bad = true
		}
	}
	if bad || l == nil {
		return nil
	}

	if xc.constant {
		xc.errorAt(x.Pos(), "function call in constant context")
		return nil
	}

	if l.valType != nil {
		xc.errorAt(x.Pos(), "type conversions not implemented")
		return nil
	}

	ei := &ExprInfo{
		Compiler: xc.Compiler,
		pos:      x.Pos(),
	}

	if ft, ok := l.ExprType.(*types.FuncType); ok && ft.Builtin != "" {
		return ei.compileBuiltinCallExpr(xc.block, ft, args)
	}

	return ei.compileCallExpr(xc.block, l, args)
}

func (xc *ExprCompiler) compileIdent(x *ast.Ident, callCtx bool) *Expr {
	ei := &ExprInfo{
		Compiler: xc.Compiler,
		pos:      x.Pos(),
	}
	return ei.compileIdent(xc.block, xc.constant, callCtx, x.Name)
}

func (xc *ExprCompiler) compileIndexExpr(x *ast.IndexExpr) *Expr {
	l := xc.compile(x.X, false)
	r := xc.compile(x.Index, false)
	if l == nil || r == nil {
		return nil
	}

	ei := &ExprInfo{
		Compiler: xc.Compiler,
		pos:      x.Pos(),
	}
	return ei.compileIndexExpr(l, r)
}

func (xc *ExprCompiler) compileSliceExpr(x *ast.SliceExpr) *Expr {
	var lo, hi *Expr
	arr := xc.compile(x.X, false)

	ei := &ExprInfo{
		Compiler: xc.Compiler,
		pos:      x.Pos(),
	}

	// If beginning was omitted, we need to provide it
	if x.Low == nil {
		lo = ei.compileIntLit("0")
	} else {
		lo = xc.compile(x.Low, false)
	}

	// If end was omitted, we need to compute len(x.X)
	if x.High == nil {
		hi = ei.compileBuiltinCallExpr(xc.block, builtins.LenType, []*Expr{arr})
	} else {
		hi = xc.compile(x.High, false)
	}

	if arr == nil || lo == nil || hi == nil {
		return nil
	}

	return ei.compileSliceExpr(arr, lo, hi)
}

func (xc *ExprCompiler) compileKeyValueExpr(x *ast.KeyValueExpr) *Expr {
	ei := &ExprInfo{
		Compiler: xc.Compiler,
		pos:      x.Pos(),
	}

	key := xc.compile(x.Key, true)
	val := xc.compile(x.Value, true)

	if key == nil {
		xc.errorAt(x.Key.Pos(), "could not compile 'key' expression")
		return nil
	}

	if val == nil {
		xc.errorAt(x.Value.Pos(), "could not compile 'value' expression")
		return nil
	}

	return ei.compileKeyValueExpr(key, val)
}

func (xc *ExprCompiler) compileSelectorExpr(x *ast.SelectorExpr) *Expr {
	ei := &ExprInfo{
		Compiler: xc.Compiler,
		pos:      x.Pos(),
	}

	v := xc.compile(x.X, false)
	if v == nil {
		return nil
	}

	return ei.compileSelectorExpr(v, x.Sel.Name)
}

func (xc *ExprCompiler) compileStarExpr(x *ast.StarExpr, callCtx bool) *Expr {
	ei := &ExprInfo{
		Compiler: xc.Compiler,
		pos:      x.Pos(),
	}

	// We pass down our call context because this could be a pointer type (and thus a type conversion).
	v := xc.compile(x.X, callCtx)
	if v == nil {
		return nil
	}

	if v.valType != nil {
		// Turns out this was a pointer type, not a dereference
		return ei.exprFromType(types.NewPtrType(v.valType))
	}

	return ei.compileStarExpr(v)
}

func (xc *ExprCompiler) compileUnaryExpr(x *ast.UnaryExpr) *Expr {
	ei := &ExprInfo{
		Compiler: xc.Compiler,
		pos:      x.Pos(),
	}

	v := xc.compile(x.X, false)
	if v == nil {
		return nil
	}

	return ei.compileUnaryExpr(x.Op, v)
}
