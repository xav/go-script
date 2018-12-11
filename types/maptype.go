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

var mapTypes = make(map[vm.Type]map[vm.Type]*MapType)

type MapType struct {
	commonType
	Key  vm.Type
	Elem vm.Type
}

// Type interface //////////////////////////////////////////////////////////////

func (t *MapType) Compat(o vm.Type, conv bool) bool { panic("NOT IMPLEMENTED") }
func (t *MapType) Lit() vm.Type                     { panic("NOT IMPLEMENTED") }
func (t *MapType) IsBoolean() bool                  { panic("NOT IMPLEMENTED") }
func (t *MapType) IsInteger() bool                  { panic("NOT IMPLEMENTED") }
func (t *MapType) IsFloat() bool                    { panic("NOT IMPLEMENTED") }
func (t *MapType) IsIdeal() bool                    { panic("NOT IMPLEMENTED") }
func (t *MapType) Zero() vm.Value                   { panic("NOT IMPLEMENTED") }
func (t *MapType) String() string                   { panic("NOT IMPLEMENTED") }
