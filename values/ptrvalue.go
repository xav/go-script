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

type PtrValue interface {
	vm.Value
	Get(*vm.Thread) vm.Value
	Set(*vm.Thread, vm.Value)
}

// ptr /////////////////////////////////////////////////////////////////////////

type PtrV struct {
	Target vm.Value // nil if the pointer is nil
}

func (v *PtrV) String() string {
	if v.Target == nil {
		return "<nil>"
	}
	return "&" + v.Target.String()
}

func (v *PtrV) Assign(t *vm.Thread, o vm.Value) {
	v.Target = o.(PtrValue).Get(t)
}

func (v *PtrV) Get(*vm.Thread) vm.Value {
	return v.Target
}

func (v *PtrV) Set(t *vm.Thread, x vm.Value) {
	v.Target = x
}
