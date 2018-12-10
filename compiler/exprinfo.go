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
	"go/token"

	"github.com/xav/go-script/context"
	"github.com/xav/go-script/types"
	"github.com/xav/go-script/values"
	"github.com/xav/go-script/vm"
)

// ExprInfo stores information needed to compile any expression node.
// Each expr also stores its exprInfo so further expressions can be compiled from it.
type ExprInfo struct {
	*Compiler
	pos token.Pos
}

func (xi *ExprInfo) compileIntLit(lit string) *Expr    { panic("NOT IMPLEMENTED") }
func (xi *ExprInfo) compileFloatLit(lit string) *Expr  { panic("NOT IMPLEMENTED") }
func (xi *ExprInfo) compileCharLit(lit string) *Expr   { panic("NOT IMPLEMENTED") }
func (xi *ExprInfo) compileStringLit(lit string) *Expr { panic("NOT IMPLEMENTED") }

func (xi *ExprInfo) compileCompositeLit(c *Expr, ikeys []interface{}, vals []*Expr) *Expr {
	panic("NOT IMPLEMENTED")
}

func (xi *ExprInfo) compileFuncLit(decl *types.FuncDecl, fn func(*vm.Thread) values.Func) *Expr {
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

////////////////////////////////////////////////////////////////////////////////

func (xi *ExprInfo) exprFromType(t vm.Type) *Expr { panic("NOT IMPLEMENTED") }
