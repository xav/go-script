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

import "github.com/xav/go-script/vm"

type UintType struct {
	commonType
	Bits uint // 0 for architecture-dependent types
	Ptr  bool // true for uintptr, false for all others
	Name string
}

func (t *UintType) Compat(o vm.Type, conv bool) bool { panic("NOT IMPLEMENTED") }
func (t *UintType) Lit() vm.Type                     { panic("NOT IMPLEMENTED") }
func (t *UintType) IsIdeal() bool                    { panic("NOT IMPLEMENTED") }
func (t *UintType) Zero() vm.Value                   { panic("NOT IMPLEMENTED") }
func (t *UintType) String() string                   { panic("NOT IMPLEMENTED") }

func MakeUIntType(bits uint, ptr bool, name string) *UintType {
	t := UintType{commonType{}, bits, ptr, name}
	return &t
}
