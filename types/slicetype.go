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

var sliceTypes = make(map[vm.Type]*SliceType)

type SliceType struct {
	commonType
	Elem vm.Type
}

// NewSliceType creates a new slice type with the specified element type.
func NewSliceType(elem vm.Type) *SliceType {
	t, ok := sliceTypes[elem]
	if !ok {
		t = &SliceType{
			commonType: commonType{},
			Elem:       elem,
		}
		sliceTypes[elem] = t
	}
	return t
}

// Type interface //////////////////////////////////////////////////////////////

// Compat returns whether this type is compatible with another type.
func (t *SliceType) Compat(o vm.Type, conv bool) bool {
	t2, ok := o.Lit().(*SliceType)
	if !ok {
		return false
	}
	return t.Elem.Compat(t2.Elem, conv)
}

// Lit returns this type's literal.
func (t *SliceType) Lit() vm.Type {
	return t
}

// Zero returns a new zero value of this type.
func (t *SliceType) Zero() vm.Value {
	// The value of an uninitialized slice is nil.
	// The length and capacity of a nil slice are 0.
	return &values.SliceV{Slice: values.Slice{Base: nil, Len: 0, Cap: 0}}
}

// String returns the string representation of this type.
func (t *SliceType) String() string {
	return "[]" + t.Elem.String()
}
