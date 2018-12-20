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

type IdealFloatType struct {
	commonType
}

// Type interface //////////////////////////////////////////////////////////////

// Compat returns whether this type is compatible with another type.
func (t *IdealFloatType) Compat(o vm.Type, conv bool) bool {
	_, ok := o.Lit().(*IdealFloatType)
	return ok
}

// Lit returns this type's literal.
func (t *IdealFloatType) Lit() vm.Type {
	return t
}

// IsFloat returns true if this is a floating type.
func (t *IdealFloatType) IsFloat() bool {
	return true
}

// IsIdeal returns true if this represents an ideal value.
func (t *IdealFloatType) IsIdeal() bool {
	return true
}

// Zero returns a new zero value of this type.
func (t *IdealFloatType) Zero() vm.Value {
	return &values.IdealFloatV{V: big.NewRat(0, 1)}
}

// String returns the string representation of this type.
func (t *IdealFloatType) String() string {
	return "ideal float"
}
