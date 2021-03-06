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

type Slice struct {
	Base ArrayValue
	Len  int64
	Cap  int64
}

type SliceValue interface {
	vm.Value
	Get(*vm.Thread) Slice
	Set(*vm.Thread, Slice)
}

// slice ///////////////////////////////////////////////////////////////////////

type SliceV struct {
	Slice
}

func (v *SliceV) String() string {
	if v.Base == nil {
		return "<nil>"
	}
	return v.Base.Sub(0, v.Len).String()
}

func (v *SliceV) Assign(t *vm.Thread, o vm.Value) {
	v.Slice = o.(SliceValue).Get(t)
}

func (v *SliceV) Get(*vm.Thread) Slice {
	return v.Slice
}

func (v *SliceV) Set(t *vm.Thread, x Slice) {
	v.Slice = x
}
