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
	"math/big"

	"github.com/xav/go-script/values"
	"github.com/xav/go-script/vm"
)

var (
	idealZero = big.NewInt(0)
	idealOne  = big.NewInt(1)
)

type IdealIntType struct {
	commonType
}

// Type interface //////////////////////////////////////////////////////////////

// Compat returns whether this type is compatible with another type.
func (t *IdealIntType) Compat(o vm.Type, conv bool) bool {
	_, ok := o.Lit().(*IdealIntType)
	return ok
}

// Lit returns this type's literal.
func (t *IdealIntType) Lit() vm.Type {
	return t
}

// IsInteger returns true if this is an integer type.
func (t *IdealIntType) IsInteger() bool {
	return true
}

// IsIdeal returns true if this represents an ideal value.
func (t *IdealIntType) IsIdeal() bool {
	return true
}

// Zero returns a new zero value of this type.
func (t *IdealIntType) Zero() vm.Value {
	return &values.IdealIntV{V: idealZero}
}

// String returns the string representation of this type.
func (t *IdealIntType) String() string {
	return "ideal integer"
}
