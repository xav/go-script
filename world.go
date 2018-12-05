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

package script

import (
	"go/ast"
	"go/build"
	"go/parser"
	"go/scanner"
	"go/token"
	"path/filepath"
	"strconv"

	"github.com/pkg/errors"

	"github.com/xav/go-script/compiler"
	"github.com/xav/go-script/context"
	"github.com/xav/go-script/types"
	"github.com/xav/go-script/vm"
)

// Status for package visiting map
type packageStatus int

const (
	unvisited packageStatus = iota
	visiting
	done
)

// Track the status of each package we visit (unvisited/visiting/done)
var pkgStatus = make(map[string]packageStatus)

type World struct {
	scope *context.Scope
	frame *vm.Frame
}

// NewWorld creates a new World by asking the universe for a child scope.
func NewWorld() *World {
	w := &World{
		scope: universe.EnterChildScope(),
	}

	// this block's vars allocate directly
	w.scope.Global = true

	return w
}

func (w *World) CompilePackage(fset *token.FileSet, files []*ast.File, pkgPath string) (Runnable, error) {
	packageFiles := make(map[string]*ast.File)
	for _, f := range files {
		packageFiles[f.Name.Name] = f
	}

	switch pkgStatus[pkgPath] {
	case done:
		return &stmtCode{
			world: w,
			code:  make(vm.Code, 0),
		}, nil

	case visiting:
		logger.Warn().
			Str("package", pkgPath).
			Msg("package dependency cycle")
		return nil, errors.Errorf("package dependency cycle: [%v]", pkgPath)
	}

	pkgStatus[pkgPath] = visiting

	packageImports := []*ast.ImportSpec{}
	for _, f := range files {
		packageImports = append(packageImports, f.Imports...)
	}

	// Compile the imports before taking care of the package code
	for _, imp := range packageImports {
		path, _ := strconv.Unquote(imp.Path.Value)
		if _, ok := universe.Pkgs[path]; ok {
			// package was already compiled
			continue
		}

		importFiles, err := loadPkgFiles(path)
		if err != nil {
			return nil, errors.Wrapf(err, "could not retrieve files for package [%s]", path)
		}

		packageCode, err := w.CompilePackage(fset, importFiles, path)
		if err != nil {
			return nil, errors.Wrapf(err, "could not compile package [%s]", path)
		}

		_, err = packageCode.Run()
		if err != nil {
			return nil, errors.Wrapf(err, "error initializing package [%s]", path)
		}
	}

	// We treat each package as the root of the world during compilation
	prevScope := w.scope
	w.scope = w.scope.EnterChildScope()
	w.scope.Global = true
	defer func() {
		pkgStatus[pkgPath] = done
		universe.Pkgs[pkgPath] = w.scope
		w.scope.Exit()
		if pkgPath != "main" {
			w.scope = prevScope
		}
	}()

	packageDecls := make([]ast.Decl, 0)
	for _, f := range files {
		packageDecls = append(packageDecls, f.Decls...)
	}
	packageCode, err := w.compileDeclList(fset, packageDecls)
	if err != nil {
		return nil, errors.Wrapf(err, "could not compile declarations in package [%s]", pkgPath)
	}
	_, err = packageCode.Run()
	if err != nil {
		return nil, errors.Wrapf(err, "could not initialize declarations in package [%s]", pkgPath)
	}

	// TODO: make sure all the imports are used at this point.

	return packageCode, nil
}

// compileDeclList compiles a list of declaration nodes
func (w *World) compileDeclList(fset *token.FileSet, decls []ast.Decl) (Runnable, error) {
	stmts := make([]ast.Stmt, len(decls))
	for i, d := range decls {
		stmts[i] = &ast.DeclStmt{
			Decl: d,
		}
	}
	return w.compileStmtList(fset, stmts)
}

// compileStmtList compiles a list of statement nodes
func (w *World) compileStmtList(fset *token.FileSet, stmts []ast.Stmt) (Runnable, error) {
	if len(stmts) == 1 {
		if s, ok := stmts[0].(*ast.ExprStmt); ok {
			return w.compileExpr(fset, s.X)
		}
	}

	errs := new(scanner.ErrorList)
	cb := compiler.NewCodeBuf()
	cc := &compiler.Compiler{
		FSet:         fset,
		Errors:       errs,
		NumErrors:    0,
		SilentErrors: 0,
	}
	fc := &compiler.FuncCompiler{
		Compiler:     cc,
		FnType:       nil,
		OutVarsNamed: false,
		CodeBuf:      cb,
		Flow:         compiler.NewFlowBuf(cb),
		Labels:       make(map[string]*compiler.Label),
	}
	bc := &compiler.BlockCompiler{
		FuncCompiler: fc,
		Block:        w.scope.Block,
	}
	nerr := cc.NumError()

	for _, stmt := range stmts {
		bc.CompileStmt(stmt)
	}
	fc.CheckLabels()

	if nerr != cc.NumError() {
		errs.Sort()
		return nil, errs.Err()
	}

	return &stmtCode{
		world: w,
		code:  fc.Get(),
	}, nil

}

// compileExpr compiles expression nodes
func (w *World) compileExpr(fset *token.FileSet, e ast.Expr) (Runnable, error) {
	errs := new(scanner.ErrorList)
	cc := &compiler.Compiler{
		FSet:         fset,
		Errors:       errs,
		NumErrors:    0,
		SilentErrors: 0,
	}

	ec := cc.CompileExpr(w.scope.Block, false, e)
	if ec == nil {
		errs.Sort()
		return nil, errs.Err()
	}

	var eval func(vm.Value, *vm.Thread)
	switch t := ec.ExprType.(type) {
	// case *types.IdealIntType:
	// 	// nothing
	// case *types.FloatType:
	// 	// nothing
	default:
		if tm, ok := t.(*types.MultiType); ok && len(tm.Elems) == 0 {
			return &stmtCode{
				world: w,
				code:  vm.Code{ec.Exec},
			}, nil
		}
		// eval = compiler.GenAssign(ec.Type, ec)
	}
	return &ExprCode{
		world: w,
		expr:  ec,
		eval:  eval,
	}, nil
}

// loadPkgFiles finds and loads the files for the specified package.
func loadPkgFiles(path string) ([]*ast.File, error) {
	pkg, err := build.Import(path, "", 0)
	if err != nil {
		return nil, err
	}

	fset := token.NewFileSet()
	var files []*ast.File
	for i := range pkg.GoFiles {
		pkgfile := filepath.Join(pkg.Dir, pkg.GoFiles[i])
		f, err := parser.ParseFile(fset, pkgfile, nil, 0)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}

	return files, nil
}
