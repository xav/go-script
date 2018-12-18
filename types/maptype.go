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

var mapTypes = make(map[vm.Type]map[vm.Type]*MapType)

type MapType struct {
	commonType
	Key  vm.Type
	Elem vm.Type
}

// NewMapType creates a new map type of the specified type.
func NewMapType(key vm.Type, elem vm.Type) *MapType {
	ts, ok := mapTypes[key]
	if !ok {
		ts = make(map[vm.Type]*MapType)
		mapTypes[key] = ts
	}

	t, ok := ts[elem]
	if !ok {
		t = &MapType{commonType{}, key, elem}
		ts[elem] = t
	}

	return t
}

// Type interface //////////////////////////////////////////////////////////////

// Compat returns whether this type is compatible with another type.
func (t *MapType) Compat(o vm.Type, conv bool) bool {
	t2, ok := o.Lit().(*MapType)
	if !ok {
		return false
	}
	return t.Elem.Compat(t2.Elem, conv) && t.Key.Compat(t2.Key, conv)
}

// Lit returns this type's literal.
func (t *MapType) Lit() vm.Type {
	return t
}

// Zero returns a new zero value of this type.
func (t *MapType) Zero() vm.Value {
	// The value of an uninitialized map is nil.
	return &values.MapV{Target: nil}
}

// String returns the string representation of this type.
func (t *MapType) String() string {
	return "map[" + t.Key.String() + "] " + t.Elem.String()
}
