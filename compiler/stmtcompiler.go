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

	"github.com/xav/go-script/builtins"
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

func (sc *stmtCompiler) error(format string, args ...interface{}) {
	sc.errorAt(sc.pos, format, args...)
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

func (sc *stmtCompiler) compileTypeDecl(decl *ast.GenDecl) bool {
	name := decl.Specs[0].(*ast.TypeSpec).Name.Name
	_, level, dup := sc.Block.Lookup(name)
	if dup != nil && level == 0 {
		sc.error("type %s redeclared in this block, previous declaration at %s", name, sc.FSet.Position(dup.Pos()))
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
	}

	return ok
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
	nerr := sc.NumError()

	//////////////////////////////////////
	// Compile right side first so we have the types when compiling the left side and
	// so we don't see definitions made on the left side.

	rs := make([]*Expr, len(rhs))
	for i, re := range rhs {
		rs[i] = sc.CompileExpr(sc.Block, false, re)
	}

	// Check if compileExpr failed on any expr
	for _, r := range rs {
		if r == nil {
			return
		}
	}

	errOp := "assignment"
	if tok == token.DEFINE || tok == token.VAR {
		errOp = "declaration"
	}

	ac, ok := sc.checkAssign(sc.pos, rs, errOp, "value")
	ac.allowMapForms(len(lhs))

	// If this is a definition and the LHS is larger than the RHS, we won't be able to produce
	// the usual error message because we can't begin to infer the types of the LHS.
	if (tok == token.DEFINE || tok == token.VAR) && len(lhs) > len(ac.rmt.Elems) {
		sc.error("not enough values for definition")
	}

	// Compile left type if there is one
	var declType vm.Type
	if declTypeExpr != nil {
		declType = sc.compileType(sc.Block, declTypeExpr)
	}

	//////////////////////////////////////
	// Compile left side

	ls := make([]*Expr, len(lhs))
	nDefs := 0
	for i, le := range lhs {
		// If this is a definition, get the identifier and its type
		var ident *ast.Ident
		var lt vm.Type
		switch tok {
		case token.DEFINE:
			// Check that it's an identifier
			ident, ok = le.(*ast.Ident)
			if !ok {
				sc.errorAt(le.Pos(), "left side of := must be a name")
				nDefs++ // Suppress new definitions errors
				continue
			}

			// Is this simply an assignment?
			if _, ok := sc.Block.Defs[ident.Name]; ok {
				ident = nil
				break
			}
			nDefs++

		case token.VAR:
			ident = le.(*ast.Ident)
		}

		// If it's a definition, get or infer its type.
		if ident != nil {
			// Compute the identifier's type from the RHS type.
			// We use the computed MultiType so we don't have to worry about unpacking.
			switch {
			case declTypeExpr != nil:
				// We have a declaration type, use it.
				// If declType is nil, we gave an error when we compiled it.
				lt = declType

			case i >= len(ac.rmt.Elems):
				lt = nil // We already gave the "not enough" error above.
			case ac.rmt.Elems[i] == nil:
				lt = nil // We already gave the error when we compiled the RHS.

			case ac.rmt.Elems[i].IsIdeal():
				// If the type is absent and the corresponding expression is a constant expression of
				// ideal integer or ideal float type, the type of the declared variable is int or float respectively.
				switch {
				case ac.rmt.Elems[i].IsInteger():
					lt = builtins.IntType
				case ac.rmt.Elems[i].IsFloat():
					lt = builtins.Float64Type
				default:
					logger.Panic().
						Str("type", fmt.Sprintf("%v", rs[i].ExprType)).
						Msg("unexpected ideal type")
				}

			default:
				lt = ac.rmt.Elems[i]
			}

			// define the identifier
			if sc.defineVar(ident, lt) == nil {
				continue
			}
		}

		// Compile LHS
		ls[i] = sc.CompileExpr(sc.Block, false, le)
		if ls[i] == nil {
			continue
		}

		if ls[i].evalMapValue != nil {
			// Map indexes are not generally addressable, but they are assignable.
			sub := ls[i]
			ls[i] = ls[i].newExpr(sub.ExprType, sub.desc)
			ls[i].evalMapValue = sub.evalMapValue
			mvf := sub.evalMapValue
			et := sub.ExprType
			ls[i].evalAddr = func(t *vm.Thread) vm.Value {
				m, k := mvf(t)
				e := m.Elem(t, k)
				if e == nil {
					e = et.Zero()
					m.SetElem(t, k, e)
				}
				return e
			}
		} else if ls[i].evalAddr == nil {
			ls[i].error("cannot assign to %s", ls[i].desc)
			continue
		}
	}

	// A short variable declaration may redeclare variables provided they were originally
	// declared in the same block with the same type, and at least one of the variables is new.
	if tok == token.DEFINE && nDefs == 0 {
		sc.error("at least one new variable must be declared")
		return
	}

	// If there have been errors, our arrays are full of nil's so get out of here now.
	if nerr != sc.NumError() {
		return
	}

	// Check for 'a[x] = r, ok'
	if len(ls) == 1 && len(rs) == 2 && ls[0].evalMapValue != nil {
		sc.error("a[x] = r, ok form not implemented")
		return
	}

	//////////////////////////////////////
	// Create assigner

	var lt vm.Type
	n := len(lhs)
	if n == 1 {
		lt = ls[0].ExprType
	} else {
		lts := make([]vm.Type, len(ls))
		for i, l := range ls {
			if l != nil {
				lts[i] = l.ExprType
			}
		}
		lt = types.NewMultiType(lts)
	}

	bc := sc.enterChild()
	defer bc.exit()
	assign := ac.compile(bc.Block, lt)
	if assign == nil {
		return
	}

	//////////////////////////////////////
	// Compile

	switch {
	case n == 1:
		// Don't need temporaries and can avoid []Value.
		lf := ls[0].evalAddr
		sc.Push(func(t *vm.Thread) {
			assign(lf(t), t)
		})

	case tok == token.VAR || (tok == token.DEFINE && nDefs == n):
		// Don't need temporaries
		lfs := make([]func(*vm.Thread) vm.Value, n)
		for i, l := range ls {
			lfs[i] = l.evalAddr
		}

		sc.Push(func(t *vm.Thread) {
			dest := make([]vm.Value, n)
			for i, lf := range lfs {
				dest[i] = lf(t)
			}
			assign(values.MultiV(dest), t)
		})

	default:
		// Need temporaries
		lmt := lt.(*types.MultiType)
		lfs := make([]func(*vm.Thread) vm.Value, n)
		for i, l := range ls {
			lfs[i] = l.evalAddr
		}

		sc.Push(func(t *vm.Thread) {
			temp := lmt.Zero().(values.MultiV)
			assign(temp, t)
			// Copy to destination
			for i := 0; i < n; i++ {
				// TODO: Need to evaluate LHS before RHS
				lfs[i](t).Assign(t, temp[i])
			}
		})
	}
}

////////////////////////////////////////////////////////////////////////////////

// default importer of packages
var defaultImporter gotypes.Importer = importer.Default()

// srcImporter implements the ast.Importer signature.
func srcImporter(typesImporter gotypes.Importer, path string) (pkg *gotypes.Package, err error) {
	return typesImporter.Import(path)
}
