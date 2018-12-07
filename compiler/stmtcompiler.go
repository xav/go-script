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
	"go/importer"
	"go/token"
	"strconv"

	gotypes "go/types"

	"github.com/pkg/errors"
	"github.com/xav/go-script/context"
	"github.com/xav/go-script/types"
	"github.com/xav/go-script/values"
	"github.com/xav/go-script/vm"
)

const (
	returnPC = ^uint(0)
	badPC    = ^uint(1)
)

// stmtCompiler is used for the compilation of individual statements.
type stmtCompiler struct {
	*BlockCompiler           // The BlockCompiler for the block enclosing this statement.
	pos            token.Pos // The position of this statement.
	stmtLabel      *Label    // This statement's label, or nil if it is not labeled.
}

func (sc *stmtCompiler) error(format string, args ...interface{}) error {
	return sc.errorAt(sc.pos, format, args...)
}

func (sc *stmtCompiler) compile(s ast.Stmt) {
	if sc.Block.Inner != nil {
		logger.Panic().Msg("Child scope still entered")
	}

	notImplemented := false
	switch s := s.(type) {
	case *ast.BadStmt:
		sc.SilentErrors++ // Error already reported by parser.

	case *ast.AssignStmt:
		notImplemented = true

	case *ast.BlockStmt:
		notImplemented = true

	case *ast.BranchStmt:
		notImplemented = true

	case *ast.CaseClause:
		// sc.error("case clause outside switch")

	case *ast.CommClause:
		// sc.error("case clause outside select")
		notImplemented = true

	case *ast.DeclStmt:
		sc.compileDeclStmt(s)

	case *ast.DeferStmt:
		notImplemented = true

	case *ast.EmptyStmt:
		// Do nothing.

	case *ast.ExprStmt:
		notImplemented = true

	case *ast.ForStmt:
		notImplemented = true

	case *ast.GoStmt:
		notImplemented = true

	case *ast.IfStmt:
		notImplemented = true

	case *ast.IncDecStmt:
		notImplemented = true

	case *ast.LabeledStmt:
		notImplemented = true

	case *ast.RangeStmt:
		notImplemented = true

	case *ast.ReturnStmt:
		notImplemented = true

	case *ast.SelectStmt:
		notImplemented = true

	case *ast.SendStmt:
		notImplemented = true

	case *ast.SwitchStmt:
		notImplemented = true

	case *ast.TypeSwitchStmt:
		notImplemented = true

	default:
		logger.Panic().
			Str("type", fmt.Sprintf("%T", s)).
			Msg("Unexpected ast node type")
	}

	if notImplemented {
		sc.error("%T statement node not implemented", s)
	}

	if sc.Block.Inner != nil {
		logger.Panic().Msg("Forgot to exit child scope")
	}
}

// Statements //////////////////////////////////////////////////////////////////

func (sc *stmtCompiler) compileDeclStmt(stmt *ast.DeclStmt) {
	switch decl := stmt.Decl.(type) {
	case *ast.BadDecl:
		sc.SilentErrors++ // Error already reported by parser.
		return

	case *ast.FuncDecl:
		sc.compileFuncDecl(decl)

	case *ast.GenDecl:
		sc.compileGenDecl(decl)

	default:
		logger.Panic().
			Str("type", fmt.Sprintf("%T", stmt.Decl)).
			Msg("Unexpected Decl type")
	}
}

// Declarations ////////////////////////////////////////////////////////////////

func (sc *stmtCompiler) compileFuncDecl(fd *ast.FuncDecl) {
	if !sc.Block.Global {
		logger.Panic().Msg("FuncDecl at statement level")
	}

	ftd := sc.compileFuncType(sc.Block, fd.Type)
	if ftd == nil {
		return
	}

	// Declare and initialize v before compiling func so that body can refer to itself.
	c, prev := sc.Block.DefineConst(fd.Name.Name, sc.pos, ftd.Type, ftd.Type.Zero())
	if prev != nil {
		pos := prev.Pos()
		if pos.IsValid() {
			sc.errorAt(fd.Name.Pos(), "identifier %s redeclared in this block\n\tprevious declaration at %s", fd.Name.Name, sc.FSet.Position(pos))
		} else {
			sc.errorAt(fd.Name.Pos(), "identifier %s redeclared in this block", fd.Name.Name)
		}
	}

	fn := sc.compileFunc(sc.Block, ftd, fd.Body)
	if c == nil || fn == nil {
		// when the compilation failed, remove the func identifier from the block definitions.
		sc.Block.Undefine(fd.Name.Name)
		return
	}

	var zeroThread vm.Thread
	c.Value.(values.FuncValue).Set(nil, fn(&zeroThread))
}

func (sc *stmtCompiler) compileGenDecl(gd *ast.GenDecl) {
	switch gd.Tok {
	case token.CONST:
		sc.compileConstDecl(gd)
	case token.IMPORT:
		sc.compileImportDecl(gd)
	case token.TYPE:
		sc.compileTypeDecl(gd)
	case token.VAR:
		sc.compileVarDecl(gd)
	default:
		logger.Panic().
			Str("token", fmt.Sprintf("%c", gd.Tok)).
			Msg("Unexpected GenDecl token: %v")
	}
}

func (sc *stmtCompiler) compileConstDecl(decl *ast.GenDecl) {
	panic("NOT IMPLEMENTED")
}

func (sc *stmtCompiler) compileImportDecl(decl *ast.GenDecl) {
	if !sc.Block.Global {
		logger.Panic().Msg("import at statement level")
	}

	for _, spec := range decl.Specs {
		spec := spec.(*ast.ImportSpec)
		path, _ := strconv.Unquote(spec.Path.Value)
		id := path
		if spec.Name != nil {
			id = spec.Name.Name
		}

		pkg, err := srcImporter(defaultImporter, path)
		if err != nil {
			sc.errorAt(spec.Pos(), "could not import package [%s]: %v", path, err)
			continue
		}

		if spec.Name != nil {
			sc.definePkg(spec.Name, id, path)
		} else {
			id = pkg.Name()
			sc.definePkg(spec.Path, id, path)
		}
	}
}

func (sc *stmtCompiler) compileTypeDecl(decl *ast.GenDecl) error {
	name := decl.Specs[0].(*ast.TypeSpec).Name.Name
	_, level, dup := sc.Block.Lookup(name)
	if dup != nil && level == 0 {
		return sc.error("type %s redeclared in this block, previous declaration at %s", name, sc.FSet.Position(dup.Pos()))
	}

	ok := true
	for _, spec := range decl.Specs {
		spec := spec.(*ast.TypeSpec)

		// Create incomplete type
		nt := sc.Block.DefineType(spec.Name.Name, spec.Name.Pos(), nil)
		if nt != nil {
			nt.(*types.NamedType).Incomplete = true
		}

		// Compile type
		tc := &typeCompiler{
			Compiler:  sc.BlockCompiler.Compiler,
			block:     sc.Block,
			lateCheck: noLateCheck,
		}
		t := tc.compileType(spec.Type, false)
		if t == nil {
			// Create a placeholder type
			ok = false
		}

		// Fill incomplete type
		nt.(*types.NamedType).Complete(t)

		// Perform late type checking with complete type
		if !tc.lateCheck() {
			ok = false
			if nt != nil {
				nt.(*types.NamedType).Def = nil // Make the type a placeholder
			}
		}
	}

	if !ok {
		sc.Block.Undefine(name)
		return errors.Errorf("error compiling type %s.", name)
	}

	return nil
}

func (sc *stmtCompiler) compileVarDecl(decl *ast.GenDecl) error {
	for _, spec := range decl.Specs {
		spec := spec.(*ast.ValueSpec)
		if spec.Values == nil { // Declaration without assignment
			if spec.Type == nil {
				// The parser should have caught that
				logger.Panic().Msg("Type and Values nil in variable declaration")
			}

			t := sc.compileType(sc.Block, spec.Type)
			if t != nil {
				// If type compilation succeeded, define placeholders
				for _, n := range spec.Names {
					sc.defineVar(n, t)
				}
			}
		} else { // Declaration with assignment
			lhs := make([]ast.Expr, len(spec.Names))
			for i, n := range spec.Names {
				lhs[i] = n
			}

			sc.doAssign(lhs, spec.Values, decl.Tok, spec.Type)
		}
	}

	return nil
}

// Statement generation helpers ////////////////////////////////////////////////

func (sc *stmtCompiler) defineVar(ident *ast.Ident, t vm.Type) *types.Variable {
	v, prev := sc.Block.DefineVar(ident.Name, ident.Pos(), t)
	if prev != nil {
		if prev.Pos().IsValid() {
			sc.errorAt(ident.Pos(), "variable %s redeclared in this block\n\tprevious declaration at %s", ident.Name, sc.FSet.Position(prev.Pos()))
		} else {
			sc.errorAt(ident.Pos(), "variable %s redeclared in this block", ident.Name)
		}
		return nil
	}

	// Initialize the variable
	index := v.Index
	if v.Index >= 0 {
		sc.Push(func(v *vm.Thread) {
			v.Frame.Vars[index] = t.Zero()
		})
	}
	return v
}

func (sc *stmtCompiler) definePkg(ident ast.Node, id, path string) *context.PkgIdent {
	v, prev := sc.Block.DefinePackage(id, path, ident.Pos())
	if prev != nil {
		sc.errorAt(ident.Pos(), "%s redeclared as imported package name\n\tprevious declaration at %s", id, sc.FSet.Position(prev.Pos()))
		return nil
	}
	return v
	panic("NOT IMPLEMENTED")
}

////////////////////////////////////////////////////////////////////////////////

func (sc *stmtCompiler) doAssign(lhs []ast.Expr, rhs []ast.Expr, tok token.Token, declTypeExpr ast.Expr) {
	panic("NOT IMPLEMENTED")
}

////////////////////////////////////////////////////////////////////////////////

// default importer of packages
var defaultImporter gotypes.Importer = importer.Default()

// srcImporter implements the ast.Importer signature.
func srcImporter(typesImporter gotypes.Importer, path string) (pkg *gotypes.Package, err error) {
	return typesImporter.Import(path)
}
