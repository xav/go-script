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

var multiTypes = newTypeArrayMap()

var EmptyType vm.Type = NewMultiType([]vm.Type{})

// MultiType is a special type used for multi-valued expressions, akin to a tuple type.
// It's not generally accessible within the language.
type MultiType struct {
	commonType
	Elems []vm.Type
}

func NewMultiType(elems []vm.Type) *MultiType {
	if t := multiTypes.Get(elems); t != nil {
		return t.(*MultiType)
	}

	t := &MultiType{
		commonType: commonType{},
		Elems:      elems,
	}
	multiTypes.Put(elems, t)
	return t
}

// Type interface //////////////////////////////////////////////////////////////

// Compat returns whether this type is compatible with another type.
func (t *MultiType) Compat(o vm.Type, conv bool) bool {
	t2, ok := o.Lit().(*MultiType)
	if !ok {
		return false
	}
	if len(t.Elems) != len(t2.Elems) {
		return false
	}
	for i := range t.Elems {
		if !t.Elems[i].Compat(t2.Elems[i], conv) {
			return false
		}
	}
	return true
}

// Lit returns this type's literal.
func (t *MultiType) Lit() vm.Type {
	return t
}

// Zero returns a new zero value of this type.
func (t *MultiType) Zero() vm.Value {
	res := make([]vm.Value, len(t.Elems))
	for i, t := range t.Elems {
		res[i] = t.Zero()
	}
	return values.MultiV(res)
}

// String returns the string representation of this type.
func (t *MultiType) String() string {
	if len(t.Elems) == 0 {
		return "<none>"
	}
	return typeListString(t.Elems, nil)
}
