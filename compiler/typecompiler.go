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
	"go/ast"

	"github.com/xav/go-script/context"
	"github.com/xav/go-script/types"
	"github.com/xav/go-script/vm"
)

type typeCompiler struct {
	*Compiler
	block     *context.Block
	lateCheck func() bool // Check to be performed after a type declaration is compiled.
}

func noLateCheck() bool {
	return true
}

func (tc *typeCompiler) compileType(x ast.Expr, allowRec bool) vm.Type {
	switch x := x.(type) {
	case *ast.BadExpr:
		tc.SilentErrors++ // Error already reported by parser
		return nil

	case *ast.ArrayType:
		return tc.compileArrayType(x, allowRec)

	case *ast.ChanType:
		return tc.compileChanType(x, allowRec)

	case *ast.Ellipsis:
		tc.errorAt(x.Pos(), "illegal use of ellipsis")
		return nil

	case *ast.FuncType:
		fd := tc.compileFuncType(x, allowRec)
		if fd == nil {
			return nil
		}
		return fd.Type

	case *ast.Ident:
		return tc.compileIdent(x, allowRec)

	case *ast.InterfaceType:
		return tc.compileInterfaceType(x, allowRec)

	case *ast.MapType:
		return tc.compileMapType(x)

	case *ast.ParenExpr:
		return tc.compileType(x.X, allowRec)

	case *ast.StarExpr:
		return tc.compilePtrType(x)

	case *ast.StructType:
		return tc.compileStructType(x, allowRec)
	}

	tc.errorAt(x.Pos(), "expression used as type")
	return nil

}

func (tc *typeCompiler) compileArrayType(x *ast.ArrayType, allowRec bool) vm.Type {
	panic("NOT IMPLEMENTED")
}

func (tc *typeCompiler) compileChanType(x *ast.ChanType, allowRec bool) vm.Type {
	panic("NOT IMPLEMENTED")
}

func (tc *typeCompiler) compileFuncType(x *ast.FuncType, allowRec bool) *types.FuncDecl {
	panic("NOT IMPLEMENTED")
}

func (tc *typeCompiler) compileIdent(x *ast.Ident, allowRec bool) vm.Type {
	panic("NOT IMPLEMENTED")
}

func (tc *typeCompiler) compileInterfaceType(x *ast.InterfaceType, allowRec bool) *types.InterfaceType {
	panic("NOT IMPLEMENTED")
}

func (tc *typeCompiler) compileMapType(x *ast.MapType) vm.Type {
	panic("NOT IMPLEMENTED")
}

func (tc *typeCompiler) compilePtrType(x *ast.StarExpr) vm.Type {
	panic("NOT IMPLEMENTED")
}

func (tc *typeCompiler) compileStructType(x *ast.StructType, allowRec bool) vm.Type {
	panic("NOT IMPLEMENTED")
}
