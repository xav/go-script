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

var arrayTypes = make(map[int64]map[vm.Type]*ArrayType)

// ArrayType is used to represent array types (https://golang.org/ref/spec#Array_types).
type ArrayType struct {
	commonType
	Len  int64   // The length of the array
	Elem vm.Type // The type of the array elements
}

// Two array types are identical if they have identical element types and the same array length.

// NewArrayType creates a new array type with the specified length and element type.
func NewArrayType(len int64, elem vm.Type) *ArrayType {
	ts, ok := arrayTypes[len]
	if !ok {
		ts = make(map[vm.Type]*ArrayType)
		arrayTypes[len] = ts
	}

	t, ok := ts[elem]
	if !ok {
		t = &ArrayType{commonType{}, len, elem}
		ts[elem] = t
	}

	return t
}

// Type interface //////////////////////////////////////////////////////////////

// Compat returns whether this type is compatible with another type.
func (t *ArrayType) Compat(o vm.Type, conv bool) bool {
	t2, ok := o.Lit().(*ArrayType)
	if !ok {
		return false
	}
	return t.Len == t2.Len && t.Elem.Compat(t2.Elem, conv)
}

// Lit returns this type's literal.
func (t *ArrayType) Lit() vm.Type {
	return t
}

// Zero returns a new zero value of this type.
func (t *ArrayType) Zero() vm.Value {
	res := values.ArrayV(make([]vm.Value, t.Len))
	// TODO: It's unfortunate that each element is separately heap allocated.
	// We could add ZeroArray to everything, though that doesn't help with multidimensional arrays.
	// Or we could do something unsafe. We'll have this same problem with structs.
	for i := int64(0); i < t.Len; i++ {
		res[i] = t.Elem.Zero()
	}
	return &res
}

// String returns the string representation of this type.
func (t *ArrayType) String() string {
	return "[]" + t.Elem.String()
}
