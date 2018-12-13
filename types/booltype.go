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

// BoolType is used to represent boolean types (https://golang.org/ref/spec#Boolean_types).
type BoolType struct {
	commonType
}

// Type interface //////////////////////////////////////////////////////////////

// Compat returns whether this type is compatible with another type.
func (t *BoolType) Compat(o vm.Type, conv bool) bool {
	_, ok := o.Lit().(*BoolType)
	return ok
}

// Lit returns this type's literal.
func (t *BoolType) Lit() vm.Type {
	return t
}

// IsBoolean returns true if this is a boolean type.
func (t *BoolType) IsBoolean() bool {
	return true
}

// Zero returns a new zero value of this type.
func (t *BoolType) Zero() vm.Value {
	res := values.BoolV(false)
	return &res
}

// String returns the string representation of this type.
func (t *BoolType) String() string {
	return "<bool>"
}
