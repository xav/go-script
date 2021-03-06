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

import "github.com/xav/go-script/vm"

type Interface struct {
	Type  vm.Type
	Value vm.Value
}

type InterfaceValue interface {
	vm.Value
	Get(*vm.Thread) Interface
	Set(*vm.Thread, Interface)
}

// interface ///////////////////////////////////////////////////////////////////

type InterfaceV struct {
	Interface
}

func (v *InterfaceV) String() string {
	if v.Type == nil || v.Value == nil {
		return "<nil>"
	}
	return v.Value.String()
}

func (v *InterfaceV) Assign(t *vm.Thread, o vm.Value) {
	v.Interface = o.(InterfaceValue).Get(t)
}

func (v *InterfaceV) Get(*vm.Thread) Interface {
	return v.Interface
}

func (v *InterfaceV) Set(t *vm.Thread, x Interface) {
	v.Interface = x
}
