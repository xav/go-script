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
	"unsafe"

	"github.com/xav/go-script/values"
	"github.com/xav/go-script/vm"
)

var (
	minIntVal   *big.Rat
	minInt8Val  *big.Rat
	minInt16Val *big.Rat
	minInt32Val *big.Rat
	minInt64Val *big.Rat
	maxIntVal   *big.Rat
	maxInt8Val  *big.Rat
	maxInt16Val *big.Rat
	maxInt32Val *big.Rat
	maxInt64Val *big.Rat
)

func initIntMinMax() {
	intBits := uint(8 * unsafe.Sizeof(int(0)))
	one := big.NewInt(1)

	num := big.NewInt(-1)
	num.Lsh(num, intBits-1)
	minIntVal = new(big.Rat).SetInt(num)

	num.SetInt64(-1)
	num.Lsh(num, 7)
	minInt8Val = new(big.Rat).SetInt(num)

	num.Lsh(num, 8)
	minInt16Val = new(big.Rat).SetInt(num)

	num.Lsh(num, 8)
	minInt32Val = new(big.Rat).SetInt(num)

	num.Lsh(num, 8)
	minInt64Val = new(big.Rat).SetInt(num)

	num.SetInt64(1)
	num.Lsh(num, intBits-1)
	num.Sub(num, one)
	maxIntVal = new(big.Rat).SetInt(num)

	num.SetInt64(1)
	num.Lsh(num, 7)
	num.Sub(num, one)
	maxInt8Val = new(big.Rat).SetInt(num)

	num.SetInt64(1)
	num.Lsh(num, 15)
	num.Sub(num, one)
	maxInt16Val = new(big.Rat).SetInt(num)

	num.SetInt64(1)
	num.Lsh(num, 31)
	num.Sub(num, one)
	maxInt32Val = new(big.Rat).SetInt(num)

	num.SetInt64(1)
	num.Lsh(num, 63)
	num.Sub(num, one)
	maxInt64Val = new(big.Rat).SetInt(num)
}

// IntType is used to represent signed integer number types (https://golang.org/ref/spec#Numeric_types).
type IntType struct {
	commonType
	Bits uint   // 0 for architecture-dependent types
	Name string // The name of the int type
}

// MakeIntType creates a new integer type of the specified bit size.
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
	logger.Panic().Msgf("unexpected int bit count: %d", t.Bits)
	panic("unreachable")
}

// String returns the string representation of this type.
func (t *IntType) String() string {
	return "<" + t.Name + ">"
}

// BoundedType interface ///////////////////////////////////////////////////////

// MinVal returns the smallest value of this type.
func (t *IntType) MinVal() *big.Rat {
	switch t.Bits {
	case 8:
		return minInt8Val
	case 16:
		return minInt16Val
	case 32:
		return minInt32Val
	case 64:
		return minInt64Val
	case 0:
		return minIntVal
	}
	logger.Panic().Msgf("unexpected int bit count: %d", t.Bits)
	panic("unreachable")
}

// MaxVal returns the largest value of this type.
func (t *IntType) MaxVal() *big.Rat {
	switch t.Bits {
	case 8:
		return maxInt8Val
	case 16:
		return maxInt16Val
	case 32:
		return maxInt32Val
	case 64:
		return maxInt64Val
	case 0:
		return maxIntVal
	}
	logger.Panic().Msgf("unexpected int bit count: %d", t.Bits)
	panic("unreachable")
}
