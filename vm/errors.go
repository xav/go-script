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

package vm

import (
	"fmt"
)

// DivByZeroError is used to abort a thread when a division by zero occurs.
type DivByZeroError struct{}

func (DivByZeroError) Error() string {
	return "divide by zero"
}

////////////////////////////////////////////////////////////////////////////////

// NilPointerError is used to abort a thread when a nil pointer is dereferenced.
type NilPointerError struct{}

func (NilPointerError) Error() string {
	return "nil pointer dereference"
}

////////////////////////////////////////////////////////////////////////////////

// IndexError is used to abort a thread when an indexed expression is out of bounds.
type IndexError struct {
	Idx, Len int64
}

func (e IndexError) Error() string {
	if e.Idx < 0 {
		return fmt.Sprintf("negative index: %d", e.Idx)
	}
	return fmt.Sprintf("index %d exceeds length %d", e.Idx, e.Len)
}

////////////////////////////////////////////////////////////////////////////////

// SliceError is used to abort a thread when invalid indices are used for
// a slice manipulation of a Slice, Array or String.
type SliceError struct {
	Lo, Hi, Cap int64
}

func (e SliceError) Error() string {
	return fmt.Sprintf("slice [%d:%d]; cap %d", e.Lo, e.Hi, e.Cap)
}

////////////////////////////////////////////////////////////////////////////////

// KeyError is used to abort a thread when invalid key indices are used in a map expression.
type KeyError struct {
	Key interface{}
}

func (e KeyError) Error() string {
	return fmt.Sprintf("key '%v' not found in map", e.Key)
}

////////////////////////////////////////////////////////////////////////////////

// NegativeLengthError is used to abort a thread when a negative length is used when creating a Slice.
type NegativeLengthError struct {
	Len int64
}

func (e NegativeLengthError) Error() string {
	return fmt.Sprintf("negative length: %d", e.Len)
}

////////////////////////////////////////////////////////////////////////////////

// NegativeCapacityError is used to abort a thread when a negative capacity is used when creating a Slice.
type NegativeCapacityError struct {
	Len int64
}

func (e NegativeCapacityError) Error() string {
	return fmt.Sprintf("negative capacity: %d", e.Len)
}
