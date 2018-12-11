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

var sliceTypes = make(map[vm.Type]*SliceType)

type SliceType struct {
	commonType
	Elem vm.Type
}

func NewSliceType(elem vm.Type) *SliceType {
	t, ok := sliceTypes[elem]
	if !ok {
		t = &SliceType{
			commonType: commonType{},
			Elem:       elem,
		}
		sliceTypes[elem] = t
	}
	return t
}

// Type interface //////////////////////////////////////////////////////////////

func (t *SliceType) Compat(o vm.Type, conv bool) bool { panic("NOT IMPLEMENTED") }
func (t *SliceType) Lit() vm.Type                     { panic("NOT IMPLEMENTED") }
func (t *SliceType) IsBoolean() bool                  { panic("NOT IMPLEMENTED") }
func (t *SliceType) IsInteger() bool                  { panic("NOT IMPLEMENTED") }
func (t *SliceType) IsFloat() bool                    { panic("NOT IMPLEMENTED") }
func (t *SliceType) IsIdeal() bool                    { panic("NOT IMPLEMENTED") }
func (t *SliceType) Zero() vm.Value                   { panic("NOT IMPLEMENTED") }
func (t *SliceType) String() string                   { panic("NOT IMPLEMENTED") }
