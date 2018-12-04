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

import (
	"go/token"

	"github.com/xav/go-script/vm"
)

type Variable struct {
	VarPos token.Pos

	// Index of this variable in the Frame
	Index int

	// Static type of this variable
	Type vm.Type

	// Value of this variable.
	// This is only used by Scope.NewFrame; therefore, it is useful for Global scopes
	// but cannot be used in function scopes.
	init vm.Value
}

func (v *Variable) Pos() token.Pos {
	return v.VarPos
}
