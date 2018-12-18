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
	"github.com/xav/go-script/values"
	"github.com/xav/go-script/vm"
)

var structTypes = newTypeArrayMap()

type StructField struct {
	Name      string
	Type      vm.Type
	Anonymous bool
}

type StructType struct {
	commonType
	Elems []StructField
}

// NewStructType creates a new struct type with the specified fields.
func NewStructType(fields []StructField) *StructType {
	// Start by looking up just the types
	fts := make([]vm.Type, len(fields))
	for i, f := range fields {
		fts[i] = f.Type
	}
	tMapI := structTypes.Get(fts)
	if tMapI == nil {
		tMapI = structTypes.Put(fts, make(map[string]*StructType))
	}
	tMap := tMapI.(map[string]*StructType)

	// Construct key for field names
	key := ""
	for _, f := range fields {
		if f.Anonymous {
			key += "!"
		}
		key += f.Name + " "
	}

	t, ok := tMap[key]
	if !ok {
		// Create new struct type
		t = &StructType{commonType{}, fields}
		tMap[key] = t
	}

	return t
}

// Type interface //////////////////////////////////////////////////////////////

// Compat returns whether this type is compatible with another type.
func (t *StructType) Compat(o vm.Type, conv bool) bool {
	t2, ok := o.Lit().(*StructType)
	if !ok {
		return false
	}

	if len(t.Elems) != len(t2.Elems) {
		return false
	}

	for i, e := range t.Elems {
		e2 := t2.Elems[i]
		// An anonymous and a non-anonymous field are neither identical nor compatible.
		if e.Anonymous != e2.Anonymous || (!e.Anonymous && e.Name != e2.Name) || !e.Type.Compat(e2.Type, conv) {
			return false
		}
	}

	return true
}

// Lit returns this type's literal.
func (t *StructType) Lit() vm.Type {
	return t
}

// Zero returns a new zero value of this type.
func (t *StructType) Zero() vm.Value {
	res := values.StructV(make([]vm.Value, len(t.Elems)))
	for i, f := range t.Elems {
		res[i] = f.Type.Zero()
	}
	return &res
}

// String returns the string representation of this type.
func (t *StructType) String() string {
	s := "struct {"
	for i, f := range t.Elems {
		if i > 0 {
			s += "; "
		}
		if !f.Anonymous {
			s += f.Name + " "
		}
		s += f.Type.String()
	}
	return s + "}"
}
