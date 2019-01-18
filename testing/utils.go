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

package testing

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	script "github.com/xav/go-script"
)

func parseCode(t *testing.T, fset *token.FileSet, src string) []*ast.File {
	t.Helper()
	file, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		t.Error(err.Error())
		return nil
	}

	return []*ast.File{file}
}

func compile(t *testing.T, code string) error {
	t.Helper()
	fset := token.NewFileSet()
	files := parseCode(t, fset, code)
	w := script.NewWorld()
	_, err := w.CompilePackage(fset, files, "main")
	return err
}
