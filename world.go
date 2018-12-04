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
	"go/scanner"
	"go/token"
	"strconv"

	"github.com/pkg/errors"

	"github.com/xav/go-script/compiler"
	"github.com/xav/go-script/context"
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

		importFiles, err := findPkgFiles(path)
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

	// We treat each package as the "main" one during compilation
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

func (w *World) compileDeclList(fset *token.FileSet, decls []ast.Decl) (Runnable, error) {
	stmts := make([]ast.Stmt, len(decls))
	for i, d := range decls {
		stmts[i] = &ast.DeclStmt{
			Decl: d,
		}
	}
	return w.CompileStmtList(fset, stmts)
}

func (w *World) CompileStmtList(fset *token.FileSet, stmts []ast.Stmt) (Runnable, error) {
	if len(stmts) == 1 {
		if s, ok := stmts[0].(*ast.ExprStmt); ok {
			return w.CompileExpr(fset, s.X)
		}
	}

	errors := new(scanner.ErrorList)
	pc := &compiler.PackageCompiler{fset, errors, 0, 0}
	cb := compiler.NewCodeBuf()
	fc := &compiler.FuncCompiler{
		PackageCompiler: pc,
		FnType:          nil,
		OutVarsNamed:    false,
		CodeBuf:         cb,
		Flow:            compiler.NewFlowBuf(cb),
		Labels:          make(map[string]*compiler.Label),
	}
	bc := &compiler.BlockCompiler{
		FuncCompiler: fc,
		Block:        w.scope.Block,
	}
	nerr := pc.NumError()

	for _, stmt := range stmts {
		bc.CompileStmt(stmt)
	}
	fc.CheckLabels()

	if nerr != pc.NumError() {
		errors.Sort()
		return nil, errors.Err()
	}

	return &stmtCode{
		world: w,
		code:  fc.Get(),
	}, nil

}

func (w *World) CompileExpr(fset *token.FileSet, e ast.Expr) (Runnable, error) {
	panic("NOT IMPLEMENTED")
}

func (w *World) Compile(fset *token.FileSet, text string) (Runnable, error) {
	panic("NOT IMPLEMENTED")
}

////////////////////////////////////////////////////////////////////////////////

func findPkgFiles(path string) ([]*ast.File, error) {
	panic("NOT IMPLEMENTED")
}
