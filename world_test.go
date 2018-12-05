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
	"go/parser"
	"go/token"
	"os"
	"testing"

	"github.com/go-test/deep"
)

func TestMain(m *testing.M) {
	// ou := universe
	// um := &mocks.ScopeHandler{}
	// um.On("EnterChildScope").Return(&context.Scope{
	// 	Block: &context.Block{},
	// })
	// universe = um

	retCode := m.Run()

	// universe = ou

	os.Exit(retCode)
}

func TestNewWorld(t *testing.T) {
	tests := []struct {
		name string
		want *World
	}{
		{"default", &World{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewWorld()
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("NewWorld(): %v", diff)
			}
			if !got.scope.Global {
				t.Error("NewWorld(): Scope is not global")
			}
		})
	}
}

func TestWorld_CompilePackage(t *testing.T) {
	fset := token.NewFileSet()
	files := parseCode(t, fset, `
package main

func main() {
}`)

	w := NewWorld()
	w.CompilePackage(fset, files, "main")
}

func parseCode(t *testing.T, fset *token.FileSet, src string) []*ast.File {
	t.Helper()
	file, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		t.Error(err.Error())
		return nil
	}

	return []*ast.File{file}
}
