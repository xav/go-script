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

package builtins

import (
	"go/token"

	"github.com/xav/go-script/context"
	"github.com/xav/go-script/values"
)

var universePos = token.NoPos

var (
	// TrueV represents the 'true' built-in constant value.
	TrueV = values.BoolV(true)
	// FalseV represents the 'false' built-in constant value.
	FalseV = values.BoolV(false)
)

func init() {
	context.Universe.Defs["true"] = &context.Constant{ConstPos: universePos, Type: BoolType, Value: &TrueV}
	context.Universe.Defs["false"] = &context.Constant{ConstPos: universePos, Type: BoolType, Value: &FalseV}
}
