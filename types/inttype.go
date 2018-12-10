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

type IntType struct {
	commonType
	Bits uint // 0 for architecture-dependent types
	Name string
}

func MakeIntType(bits uint, name string) *IntType {
	t := IntType{commonType{}, bits, name}
	return &t
}

// Type interface //////////////////////////////////////////////////////////////

func (t *IntType) Compat(o vm.Type, conv bool) bool { panic("NOT IMPLEMENTED") }
func (t *IntType) Lit() vm.Type                     { panic("NOT IMPLEMENTED") }
func (t *IntType) IsBoolean() bool                  { panic("NOT IMPLEMENTED") }
func (t *IntType) IsInteger() bool                  { panic("NOT IMPLEMENTED") }
func (t *IntType) IsFloat() bool                    { panic("NOT IMPLEMENTED") }
func (t *IntType) IsIdeal() bool                    { panic("NOT IMPLEMENTED") }
func (t *IntType) Zero() vm.Value                   { panic("NOT IMPLEMENTED") }
func (t *IntType) String() string                   { panic("NOT IMPLEMENTED") }
