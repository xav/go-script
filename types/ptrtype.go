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

var ptrTypes = make(map[vm.Type]*PtrType)

type PtrType struct {
	commonType
	Elem vm.Type
}

// Two pointer types are identical if they have identical base types.

func NewPtrType(elem vm.Type) *PtrType {
	t, ok := ptrTypes[elem]
	if !ok {
		t = &PtrType{
			commonType: commonType{},
			Elem:       elem,
		}
		ptrTypes[elem] = t
	}
	return t
}

// Type interface //////////////////////////////////////////////////////////////

func (t *PtrType) Compat(o vm.Type, conv bool) bool { panic("NOT IMPLEMENTED") }
func (t *PtrType) Lit() vm.Type                     { panic("NOT IMPLEMENTED") }
func (t *PtrType) IsBoolean() bool                  { panic("NOT IMPLEMENTED") }
func (t *PtrType) IsInteger() bool                  { panic("NOT IMPLEMENTED") }
func (t *PtrType) IsFloat() bool                    { panic("NOT IMPLEMENTED") }
func (t *PtrType) IsIdeal() bool                    { panic("NOT IMPLEMENTED") }
func (t *PtrType) Zero() vm.Value                   { panic("NOT IMPLEMENTED") }
func (t *PtrType) String() string                   { panic("NOT IMPLEMENTED") }
