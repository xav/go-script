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

type PackageValue interface {
	vm.Value
	Get(*vm.Thread) PackageValue
	Ident(*vm.Thread, int) vm.Value
}

// package /////////////////////////////////////////////////////////////////////

type PackageV struct {
	Name   string
	Idents []vm.Value
}

func (v *PackageV) String() string                  { panic("NOT IMPLEMENTED") }
func (v *PackageV) Assign(t *vm.Thread, o vm.Value) { panic("NOT IMPLEMENTED") }