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

package values

import (
	"fmt"

	"github.com/xav/go-script/vm"
)

type IntValue interface {
	vm.Value
	Get(*vm.Thread) int64
	Set(*vm.Thread, int64)
}

// int8 ////////////////////////////////////////////////////////////////////////

type Int8V int8

func (v *Int8V) String() string {
	return fmt.Sprint(*v)
}

func (v *Int8V) Assign(t *vm.Thread, o vm.Value) {
	*v = Int8V(o.(IntValue).Get(t))
}

func (v *Int8V) Get(*vm.Thread) int64 {
	return int64(*v)
}

func (v *Int8V) Set(t *vm.Thread, x int64) {
	*v = Int8V(x)
}

// int16 ///////////////////////////////////////////////////////////////////////

type Int16V int16

func (v *Int16V) String() string {
	return fmt.Sprint(*v)
}

func (v *Int16V) Assign(t *vm.Thread, o vm.Value) {
	*v = Int16V(o.(IntValue).Get(t))
}

func (v *Int16V) Get(*vm.Thread) int64 {
	return int64(*v)
}

func (v *Int16V) Set(t *vm.Thread, x int64) {
	*v = Int16V(x)
}

// int32 ///////////////////////////////////////////////////////////////////////

type Int32V int32

func (v *Int32V) String() string {
	return fmt.Sprint(*v)
}

func (v *Int32V) Assign(t *vm.Thread, o vm.Value) {
	*v = Int32V(o.(IntValue).Get(t))
}

func (v *Int32V) Get(*vm.Thread) int64 {
	return int64(*v)
}

func (v *Int32V) Set(t *vm.Thread, x int64) {
	*v = Int32V(x)
}

// int64 ///////////////////////////////////////////////////////////////////////

type Int64V int64

func (v *Int64V) String() string {
	return fmt.Sprint(*v)
}

func (v *Int64V) Assign(t *vm.Thread, o vm.Value) {
	*v = Int64V(o.(IntValue).Get(t))
}

func (v *Int64V) Get(*vm.Thread) int64 {
	return int64(*v)
}

func (v *Int64V) Set(t *vm.Thread, x int64) {
	*v = Int64V(x)
}

// int /////////////////////////////////////////////////////////////////////////

type IntV int

func (v *IntV) String() string {
	return fmt.Sprint(*v)
}

func (v *IntV) Assign(t *vm.Thread, o vm.Value) {
	*v = IntV(o.(IntValue).Get(t))
}

func (v *IntV) Get(*vm.Thread) int64 {
	return int64(*v)
}

func (v *IntV) Set(t *vm.Thread, x int64) {
	*v = IntV(x)
}
