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

func (xi *ExprInfo) newExpr(t vm.Type, desc string) *Expr { panic("NOT IMPLEMENTED") }
func (xi *ExprInfo) exprFromType(t vm.Type) *Expr         { panic("NOT IMPLEMENTED") }

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

func (xi *ExprInfo) compileCompositeLit(c *Expr, ikeys []interface{}, vals []*Expr) *Expr {
	panic("NOT IMPLEMENTED")
}

func (xi *ExprInfo) compileFuncLit(decl *types.FuncDecl, fn func(*vm.Thread) values.Func) *Expr {
	panic("NOT IMPLEMENTED")
}

////////////////////////////////////////////////////////////////////////////////

func (xi *ExprInfo) compileGlobalVariable(v *types.Variable) *Expr      { panic("NOT IMPLEMENTED") }
func (xi *ExprInfo) compileVariable(level int, v *types.Variable) *Expr { panic("NOT IMPLEMENTED") }
func (xi *ExprInfo) compileString(s string) *Expr                       { panic("NOT IMPLEMENTED") }
func (xi *ExprInfo) compileStringList(list []*Expr) *Expr               { panic("NOT IMPLEMENTED") }

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
