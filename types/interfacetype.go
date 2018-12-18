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
	"sort"
	"strings"

	"github.com/xav/go-script/values"
	"github.com/xav/go-script/vm"
)

var interfaceTypes = newTypeArrayMap()

// TODO: Interface values, types, and type compilation are implemented,
// but none of the type checking or semantics of interfaces are.

type InterfaceType struct {
	commonType
	Methods map[string]*FuncType
}

// NewInterfaceType creates a new interface type with the specified methods.
func NewInterfaceType(methods map[string]*FuncType, embeds []*InterfaceType) *InterfaceType {
	// Count methods of embedded interfaces
	nMethods := len(methods)
	for _, e := range embeds {
		nMethods += len(e.Methods)
	}

	// Combine methods
	allMethods := make(map[string]*FuncType, nMethods)
	for _, e := range embeds {
		for k, m := range e.Methods {
			allMethods[k] = m
		}
	}
	for k, m := range methods {
		allMethods[k] = m
	}

	// Sort methods
	var ks []string
	for k := range allMethods {
		ks = append(ks, k)
	}
	sort.Strings(ks)

	// Create the type signature
	mts := make([]vm.Type, 0, len(allMethods))
	for _, k := range ks {
		mts = append(mts, allMethods[k])
	}

	tMapI := interfaceTypes.Get(mts)
	if tMapI == nil {
		tMapI = interfaceTypes.Put(mts, make(map[string]*InterfaceType))
	}
	tMap := tMapI.(map[string]*InterfaceType)

	key := strings.Join(ks, " ")
	t, ok := tMap[key]
	if !ok {
		t = &InterfaceType{commonType{}, allMethods}
		tMap[key] = t
	}
	return t
}

// implementedBy tests if o implements t, returning nil, true if it does.
// Otherwise, it returns a method of t that o is missing and false.
func (t *InterfaceType) implementedBy(o vm.Type) (*FuncType, bool) {
	if len(t.Methods) == 0 {
		return nil, true
	}

	// The methods of a named interface types are those of the underlying type.
	if it, ok := o.Lit().(*InterfaceType); ok {
		o = it
	}

	switch o := o.(type) {
	case *NamedType:
		for name, tm := range t.Methods {
			sm, ok := o.Methods[name]
			if !ok || sm.decl.Type != tm {
				return tm, false
			}
		}
		return nil, true

	case *InterfaceType:
		for name, tm := range t.Methods {
			sm, ok := o.Methods[name]
			if !ok || tm != sm {
				return tm, false
			}
		}

		return nil, true
	}

	for _, tm := range t.Methods {
		return tm, false
	}
	return nil, false
}

// Type interface //////////////////////////////////////////////////////////////

func (t *InterfaceType) Compat(o vm.Type, conv bool) bool {
	t2, ok := o.Lit().(*InterfaceType)
	if !ok {
		return false
	}
	if len(t.Methods) != len(t2.Methods) {
		return false
	}
	for name, e := range t.Methods {
		e2 := t2.Methods[name]
		if !e.Compat(e2, conv) {
			return false
		}
	}
	return true
}

func (t *InterfaceType) Lit() vm.Type {
	return t
}

func (t *InterfaceType) Zero() vm.Value {
	return &values.InterfaceV{}
}

func (t *InterfaceType) String() string {
	// TODO: Instead of showing embedded interfaces, this shows their methods.
	fs := make([]string, 0, len(t.Methods))
	for name, m := range t.Methods {
		fs = append(fs, name+funcTypeString(m, nil, nil))
	}

	return "interface {" + strings.Join(fs, ";") + "}"
}
