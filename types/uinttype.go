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
	minUIntVal    *big.Rat
	maxUIntVal    *big.Rat
	maxUIntPtrVal *big.Rat
	maxUInt8Val   *big.Rat
	maxUInt16Val  *big.Rat
	maxUInt32Val  *big.Rat
	maxUInt64Val  *big.Rat
)

func initUintMinMax() {
	uintBits := uint(8 * unsafe.Sizeof(uint(0)))
	uintptrBits := uint(8 * unsafe.Sizeof(uintptr(0)))
	one := big.NewInt(1)

	minUIntVal = big.NewRat(0, 1)

	num := big.NewInt(1)
	num.Lsh(num, uintBits)
	num.Sub(num, one)
	maxUIntVal = new(big.Rat).SetInt(num)

	num.SetInt64(1)
	num.Lsh(num, uintptrBits)
	num.Sub(num, one)
	maxUIntPtrVal = new(big.Rat).SetInt(num)

	num.SetInt64(1)
	num.Lsh(num, 8)
	num.Sub(num, one)
	maxUInt8Val = new(big.Rat).SetInt(num)

	num.SetInt64(1)
	num.Lsh(num, 16)
	num.Sub(num, one)
	maxUInt16Val = new(big.Rat).SetInt(num)

	num.SetInt64(1)
	num.Lsh(num, 32)
	num.Sub(num, one)
	maxUInt32Val = new(big.Rat).SetInt(num)

	num.SetInt64(1)
	num.Lsh(num, 64)
	num.Sub(num, one)
	maxUInt64Val = new(big.Rat).SetInt(num)
}

// UintType is used to represent unsigned integer number types (https://golang.org/ref/spec#Numeric_types).
type UintType struct {
	commonType
	Bits uint   // 0 for architecture-dependent types
	Ptr  bool   // true for uintptr, false for all others
	Name string // The name of the uint type
}

// MakeUIntType creates a new unsigned integer type of the specified bit size.
func MakeUIntType(bits uint, ptr bool, name string) *UintType {
	t := UintType{commonType{}, bits, ptr, name}
	return &t
}

// Type interface //////////////////////////////////////////////////////////////

// Compat returns whether this type is compatible with another type.
func (t *UintType) Compat(o vm.Type, conv bool) bool {
	t2, ok := o.Lit().(*UintType)
	return ok && t == t2
}

// Lit returns this type's literal.
func (t *UintType) Lit() vm.Type {
	return t
}

// IsInteger returns true if this is an integer type.
func (t *UintType) IsInteger() bool {
	return true
}

// Zero returns a new zero value of this type.
func (t *UintType) Zero() vm.Value {
	switch t.Bits {
	case 0:
		if t.Ptr {
			res := values.UintptrV(0)
			return &res
		}
		res := values.UintV(0)
		return &res
	case 8:
		res := values.Uint8V(0)
		return &res
	case 16:
		res := values.Uint16V(0)
		return &res
	case 32:
		res := values.Uint32V(0)
		return &res
	case 64:
		res := values.Uint64V(0)
		return &res
	}
	logger.Panic().Msgf("unexpected uint bit count: %d", t.Bits)
	panic("unreachable")
}

// String returns the string representation of this type.
func (t *UintType) String() string {
	return "<" + t.Name + ">"
}

// BoundedType interface ///////////////////////////////////////////////////////

// MinVal returns the smallest value of this type.
func (t *UintType) MinVal() *big.Rat {
	return minUIntVal
}

// MaxVal returns the largest value of this type.
func (t *UintType) MaxVal() *big.Rat {
	bits := t.Bits
	switch bits {
	case 8:
		return maxUInt8Val
	case 16:
		return maxUInt16Val
	case 32:
		return maxUInt32Val
	case 64:
		return maxUInt64Val
	case 0:
		if t.Ptr {
			return maxUIntPtrVal
		}
		return maxUIntVal
	}
	logger.Panic().Msgf("unexpected uint bit count: %d", t.Bits)
	panic("unreachable")
}
