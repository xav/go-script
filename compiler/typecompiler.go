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
	// TODO: Variadic function types
	in, inNames, _, inBad := tc.compileFields(x.Params, allowRec)
	out, outNames, _, outBad := tc.compileFields(x.Results, allowRec)

	if inBad || outBad {
		return nil
	}
	return &types.FuncDecl{
		Type:     types.NewFuncType(in, false, out),
		Name:     nil,
		InNames:  inNames,
		OutNames: outNames,
	}

}

func (tc *typeCompiler) compileIdent(x *ast.Ident, allowRec bool) vm.Type {
	_, _, def := tc.block.Lookup(x.Name)
	if def == nil {
		tc.errorAt(x.Pos(), "%s: undefined", x.Name)
		return nil
	}

	switch def := def.(type) {
	case *types.Constant:
		tc.errorAt(x.Pos(), "constant %v used as type", x.Name)
		return nil
	case *types.Variable:
		tc.errorAt(x.Pos(), "variable %v used as type", x.Name)
		return nil

	case *types.NamedType:
		if !allowRec && def.Incomplete {
			tc.errorAt(x.Pos(), "illegal recursive type")
			return nil
		}
		if !def.Incomplete && def.Def == nil {
			return nil // Placeholder type from an earlier error
		}
		return def

	case vm.Type:
		return def
	}

	logger.Panic().
		Str("symbol", x.Name).
		Str("type", fmt.Sprintf("%T", def)).
		Msg("symbol has unknown type")
	return nil
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

func (tc *typeCompiler) compileFields(fields *ast.FieldList, allowRec bool) ([]vm.Type, []*ast.Ident, []token.Pos, bool) {
	n := fields.NumFields()
	ts := make([]vm.Type, n)
	ns := make([]*ast.Ident, n)
	ps := make([]token.Pos, n)
	bad := false

	if fields != nil {
		i := 0
		for _, f := range fields.List {
			t := tc.compileType(f.Type, allowRec)
			if t == nil {
				bad = true
			}

			if f.Names == nil {
				ns[i] = nil
				ts[i] = t
				ps[i] = f.Type.Pos()
				i++
				continue
			}

			for _, n := range f.Names {
				ns[i] = n
				ts[i] = t
				ps[i] = n.Pos()
				i++
			}
		}
	}

	return ts, ns, ps, bad
}
