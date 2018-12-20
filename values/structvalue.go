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

type StructValue interface {
	vm.Value
	// TODO: Get() is here for uniformity, but is completely useless.
	// If a lot of other types have similarly useless Get methods, just special-case these uses.
	Get(*vm.Thread) StructValue
	Field(*vm.Thread, int) vm.Value
}

// struct //////////////////////////////////////////////////////////////////////

type StructV []vm.Value

func (v *StructV) String() string {
	res := "{"
	for i, v := range *v {
		if i > 0 {
			res += ", "
		}
		res += v.String()
	}
	return res + "}"
}

func (v *StructV) Assign(t *vm.Thread, o vm.Value) {
	oa := o.(StructValue)
	l := len(*v)
	for i := 0; i < l; i++ {
		(*v)[i].Assign(t, oa.Field(t, i))
	}
}

func (v *StructV) Get(*vm.Thread) StructValue {
	return v
}

func (v *StructV) Field(t *vm.Thread, i int) vm.Value {
	return (*v)[i]
}
