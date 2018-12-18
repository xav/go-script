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

	"github.com/xav/go-script/values"
	"github.com/xav/go-script/vm"
)

var (
	funcTypes         = newTypeArrayMap()
	variadicFuncTypes = newTypeArrayMap()
)

// Two function types are identical if they have the same number of parameters and result values,
// and if corresponding parameter and result types are identical.
// All "..." parameters have identical type.
// Parameter and result names are not required to match.

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

// NewFuncType creates a new function type with the specified input and output types.
func NewFuncType(in []vm.Type, variadic bool, out []vm.Type) *FuncType {
	inMap := funcTypes
	if variadic {
		inMap = variadicFuncTypes
	}

	outMapI := inMap.Get(in)
	if outMapI == nil {
		outMapI = inMap.Put(in, newTypeArrayMap())
	}
	outMap := outMapI.(typeArrayMap)

	tI := outMap.Get(out)
	if tI != nil {
		return tI.(*FuncType)
	}

	t := &FuncType{commonType{}, in, variadic, out, ""}
	outMap.Put(out, t)
	return t
}

// Type interface //////////////////////////////////////////////////////////////

// Compat returns whether this type is compatible with another type.
func (t *FuncType) Compat(o vm.Type, conv bool) bool {
	t2, ok := o.Lit().(*FuncType)
	if !ok {
		return false
	}
	if len(t.In) != len(t2.In) || t.Variadic != t2.Variadic || len(t.Out) != len(t2.Out) {
		return false
	}
	for i := range t.In {
		if !t.In[i].Compat(t2.In[i], conv) {
			return false
		}
	}
	for i := range t.Out {
		if !t.Out[i].Compat(t2.Out[i], conv) {
			return false
		}
	}
	return true
}

// Lit returns this type's literal.
func (t *FuncType) Lit() vm.Type {
	return t
}

// Zero returns a new zero value of this type.
func (t *FuncType) Zero() vm.Value {
	return &values.FuncV{
		Target: nil,
	}
}

// String returns the string representation of this type.
func (t *FuncType) String() string {
	if t.Builtin != "" {
		return "built-in function " + t.Builtin
	}
	args := typeListString(t.In, nil)
	if t.Variadic {
		if len(args) > 0 {
			args += ", "
		}
		args += "..."
	}
	s := "func(" + args + ")"
	if len(t.Out) > 0 {
		s += " (" + typeListString(t.Out, nil) + ")"
	}
	return s
}

func typeListString(ts []vm.Type, ns []*ast.Ident) string {
	s := ""
	for i, t := range ts {
		if i > 0 {
			s += ", "
		}
		if ns != nil && ns[i] != nil {
			s += ns[i].Name + " "
		}
		if t == nil {
			// Some places use nil types to represent errors
			s += "<none>"
		} else {
			s += t.String()
		}
	}
	return s
}

////////////////////////////////////////////////////////////////////////////////

func (t *FuncDecl) String() string {
	s := "func"
	if t.Name != nil {
		s += " " + t.Name.Name
	}
	s += funcTypeString(t.Type, t.InNames, t.OutNames)
	return s
}

func funcTypeString(ft *FuncType, ins []*ast.Ident, outs []*ast.Ident) string {
	s := "("
	s += typeListString(ft.In, ins)
	if ft.Variadic {
		if len(ft.In) > 0 {
			s += ", "
		}
		s += "..."
	}
	s += ")"
	if len(ft.Out) > 0 {
		s += " (" + typeListString(ft.Out, outs) + ")"
	}
	return s
}
