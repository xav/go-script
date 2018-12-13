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
	maxFloat32Val *big.Rat
	maxFloat64Val *big.Rat
	minFloat32Val *big.Rat
	minFloat64Val *big.Rat
)

func initFloatMinMax() {
	num := big.NewInt(0xffffff)
	num.Lsh(num, 127-23)
	maxFloat32Val = new(big.Rat).SetInt(num)

	num.SetInt64(0x1fffffffffffff)
	num.Lsh(num, 1023-52)
	maxFloat64Val = new(big.Rat).SetInt(num)

	minFloat32Val = new(big.Rat).Neg(maxFloat32Val)
	minFloat64Val = new(big.Rat).Neg(maxFloat64Val)
}

// FloatType is used to represent floating-point number types (https://golang.org/ref/spec#Numeric_types).
type FloatType struct {
	commonType
	Bits uint   // 0 for architecture-dependent type
	name string // The name of the float type
}

// MakeFloatType creates a new floating-point type of the specified bit size.
func MakeFloatType(bits uint, name string) *FloatType {
	t := FloatType{commonType{}, bits, name}
	return &t
}

// Type interface //////////////////////////////////////////////////////////////

// Compat returns whether this type is compatible with another type.
func (t *FloatType) Compat(o vm.Type, conv bool) bool {
	t2, ok := o.Lit().(*FloatType)
	return ok && t == t2
}

// Lit returns this type's literal.
func (t *FloatType) Lit() vm.Type {
	return t
}

// IsFloat returns true if this is a floating type.
func (t *FloatType) IsFloat() bool {
	return true
}

// Zero returns a new zero value of this type.
func (t *FloatType) Zero() vm.Value {
	switch t.Bits {
	case 32:
		res := values.Float32V(0)
		return &res
	case 64:
		res := values.Float64V(0)
		return &res
	}
	logger.Panic().Msgf("unexpected floating point bit count: %d", t.Bits)
	panic("unreachable")
}

// String returns the string representation of this type.
func (t *FloatType) String() string {
	return "<" + t.name + ">"
}

// BoundedType interface ///////////////////////////////////////////////////////

// MinVal returns the smallest value of this type.
func (t *FloatType) MinVal() *big.Rat {
	switch t.Bits {
	case 32:
		return minFloat32Val
	case 64:
		return minFloat64Val
	}
	logger.Panic().Msgf("unexpected floating point bit count: %d", t.Bits)
	panic("unreachable")
}

// MaxVal returns the largest value of this type.
func (t *FloatType) MaxVal() *big.Rat {
	switch t.Bits {
	case 32:
		return maxFloat32Val
	case 64:
		return maxFloat64Val
	}
	logger.Panic().Msgf("unexpected floating point bit count: %d", t.Bits)
	panic("unreachable")
}
