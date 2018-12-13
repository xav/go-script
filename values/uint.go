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

type UintValue interface {
	vm.Value
	Get(*vm.Thread) uint64
	Set(*vm.Thread, uint64)
}

// uint8 ///////////////////////////////////////////////////////////////////////

type Uint8V uint8

func (v *Uint8V) String() string {
	return fmt.Sprint(*v)
}

func (v *Uint8V) Assign(t *vm.Thread, o vm.Value) {
	*v = Uint8V(o.(UintValue).Get(t))
}

func (v *Uint8V) Get(*vm.Thread) uint64 {
	return uint64(*v)
}

func (v *Uint8V) Set(t *vm.Thread, x uint64) {
	*v = Uint8V(x)
}

// uint16 //////////////////////////////////////////////////////////////////////

type Uint16V uint16

func (v *Uint16V) String() string {
	return fmt.Sprint(*v)
}

func (v *Uint16V) Assign(t *vm.Thread, o vm.Value) {
	*v = Uint16V(o.(UintValue).Get(t))
}

func (v *Uint16V) Get(*vm.Thread) uint64 {
	return uint64(*v)
}

func (v *Uint16V) Set(t *vm.Thread, x uint64) {
	*v = Uint16V(x)
}

// uint32 //////////////////////////////////////////////////////////////////////

type Uint32V uint32

func (v *Uint32V) String() string {
	return fmt.Sprint(*v)
}

func (v *Uint32V) Assign(t *vm.Thread, o vm.Value) {
	*v = Uint32V(o.(UintValue).Get(t))
}

func (v *Uint32V) Get(*vm.Thread) uint64 {
	return uint64(*v)
}

func (v *Uint32V) Set(t *vm.Thread, x uint64) {
	*v = Uint32V(x)
}

// uint64 //////////////////////////////////////////////////////////////////////

type Uint64V uint64

func (v *Uint64V) String() string {
	return fmt.Sprint(*v)
}

func (v *Uint64V) Assign(t *vm.Thread, o vm.Value) {
	*v = Uint64V(o.(UintValue).Get(t))
}

func (v *Uint64V) Get(*vm.Thread) uint64 {
	return uint64(*v)
}

func (v *Uint64V) Set(t *vm.Thread, x uint64) {
	*v = Uint64V(x)
}

// uint ////////////////////////////////////////////////////////////////////////

type UintV uint

func (v *UintV) String() string {
	return fmt.Sprint(*v)
}

func (v *UintV) Assign(t *vm.Thread, o vm.Value) {
	*v = UintV(o.(UintValue).Get(t))
}

func (v *UintV) Get(*vm.Thread) uint64 {
	return uint64(*v)
}

func (v *UintV) Set(t *vm.Thread, x uint64) {
	*v = UintV(x)
}

// uintptr /////////////////////////////////////////////////////////////////////

type UintptrV uintptr

func (v *UintptrV) String() string {
	return fmt.Sprint(*v)
}

func (v *UintptrV) Assign(t *vm.Thread, o vm.Value) {
	*v = UintptrV(o.(UintValue).Get(t))
}

func (v *UintptrV) Get(*vm.Thread) uint64 {
	return uint64(*v)
}

func (v *UintptrV) Set(t *vm.Thread, x uint64) {
	*v = UintptrV(x)
}
