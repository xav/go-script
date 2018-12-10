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

type IdealFloatType struct {
	commonType
}

// Type interface //////////////////////////////////////////////////////////////

func (t *IdealFloatType) Compat(o vm.Type, conv bool) bool { panic("NOT IMPLEMENTED") }
func (t *IdealFloatType) Lit() vm.Type                     { panic("NOT IMPLEMENTED") }
func (t *IdealFloatType) IsBoolean() bool                  { panic("NOT IMPLEMENTED") }
func (t *IdealFloatType) IsInteger() bool                  { panic("NOT IMPLEMENTED") }
func (t *IdealFloatType) IsFloat() bool                    { panic("NOT IMPLEMENTED") }
func (t *IdealFloatType) IsIdeal() bool                    { panic("NOT IMPLEMENTED") }
func (t *IdealFloatType) Zero() vm.Value                   { panic("NOT IMPLEMENTED") }
func (t *IdealFloatType) String() string                   { panic("NOT IMPLEMENTED") }
