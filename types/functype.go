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

package types

import (
	"go/ast"

	"github.com/xav/go-script/vm"
)

type FuncDecl struct {
	Type     *FuncType
	Name     *ast.Ident   // nil for function literals
	InNames  []*ast.Ident // InNames will be one longer than Type.In if this function is variadic.
	OutNames []*ast.Ident
}

type FuncType struct {
	commonType
	In       []vm.Type
	Variadic bool
	Out      []vm.Type
	Builtin  string
}

func (t *FuncType) Compat(o vm.Type, conv bool) bool { panic("NOT IMPLEMENTED") }
func (t *FuncType) Lit() vm.Type                     { panic("NOT IMPLEMENTED") }
func (t *FuncType) IsIdeal() bool                    { panic("NOT IMPLEMENTED") }
func (t *FuncType) Zero() vm.Value                   { panic("NOT IMPLEMENTED") }
func (t *FuncType) String() string                   { panic("NOT IMPLEMENTED") }
