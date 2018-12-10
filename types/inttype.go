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

type IntType struct {
	commonType
	Bits uint // 0 for architecture-dependent types
	Name string
}

func MakeIntType(bits uint, name string) *IntType {
	t := IntType{commonType{}, bits, name}
	return &t
}

// Type interface //////////////////////////////////////////////////////////////

// Compat returns whether this type is compatible with another type.
func (t *IntType) Compat(o vm.Type, conv bool) bool {
	t2, ok := o.Lit().(*IntType)
	return ok && t == t2
}

// Lit returns this type's literal.
func (t *IntType) Lit() vm.Type {
	return t
}

// IsInteger returns true if this is an integer type.
func (t *IntType) IsInteger() bool {
	return true
}

// Zero returns a new zero value of this type.
func (t *IntType) Zero() vm.Value {
	switch t.Bits {
	case 8:
		res := values.Int8V(0)
		return &res
	case 16:
		res := values.Int16V(0)
		return &res
	case 32:
		res := values.Int32V(0)
		return &res
	case 64:
		res := values.Int64V(0)
		return &res

	case 0:
		res := values.IntV(0)
		return &res
	}
	panic("unexpected int bit count")
}

// String returns the string representation of this type.
func (t *IntType) String() string {
	return "<" + t.Name + ">"
}
