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

type StringType struct {
	commonType
}

// Type interface //////////////////////////////////////////////////////////////

// Compat returns whether this type is compatible with another type.
func (t *StringType) Compat(o vm.Type, conv bool) bool {
	_, ok := o.Lit().(*StringType)
	return ok
}

// Lit returns this type's literal.
func (t *StringType) Lit() vm.Type {
	return t
}

// Zero returns a new zero value of this type.
func (t *StringType) Zero() vm.Value {
	res := values.StringV("")
	return &res
}

// String returns the string representation of this type.
func (t *StringType) String() string {
	return "<string>"
}
