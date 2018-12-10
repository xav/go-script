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

type ArrayValue interface {
	vm.Value
	// TODO: Get() is here for uniformity, but is completely useless.
	// If a lot of other types have similarly useless Get methods, just special-case these uses.
	Get(*vm.Thread) ArrayValue         //
	Elem(*vm.Thread, int64) vm.Value   //
	Sub(i int64, len int64) ArrayValue // Sub returns an ArrayValue backed by the same array that starts from element i and has length len.
}

// array ///////////////////////////////////////////////////////////////////////

type ArrayV []vm.Value

func (v *ArrayV) String() string                  { panic("NOT IMPLEMENTED") }
func (v *ArrayV) Assign(t *vm.Thread, o vm.Value) { panic("NOT IMPLEMENTED") }
