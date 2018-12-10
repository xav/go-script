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

var arrayTypes = make(map[int64]map[vm.Type]*ArrayType)

type ArrayType struct {
	commonType
	Len  int64
	Elem vm.Type
}

// Two array types are identical if they have identical element types and the same array length.

func NewArrayType(len int64, elem vm.Type) *ArrayType {
	ts, ok := arrayTypes[len]
	if !ok {
		ts = make(map[vm.Type]*ArrayType)
		arrayTypes[len] = ts
	}

	t, ok := ts[elem]
	if !ok {
		t = &ArrayType{commonType{}, len, elem}
		ts[elem] = t
	}

	return t
}

// Type interface //////////////////////////////////////////////////////////////

func (t *ArrayType) Compat(o vm.Type, conv bool) bool { panic("NOT IMPLEMENTED") }
func (t *ArrayType) Lit() vm.Type                     { panic("NOT IMPLEMENTED") }
func (t *ArrayType) IsBoolean() bool                  { panic("NOT IMPLEMENTED") }
func (t *ArrayType) IsInteger() bool                  { panic("NOT IMPLEMENTED") }
func (t *ArrayType) IsFloat() bool                    { panic("NOT IMPLEMENTED") }
func (t *ArrayType) IsIdeal() bool                    { panic("NOT IMPLEMENTED") }
func (t *ArrayType) Zero() vm.Value                   { panic("NOT IMPLEMENTED") }
func (t *ArrayType) String() string                   { panic("NOT IMPLEMENTED") }
