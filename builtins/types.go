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
	"github.com/xav/go-script/types"
)

var universePos = token.NoPos

var (
	BoolType = context.Universe.DefineType("bool", universePos, &types.BoolType{})
)

var (
	IntType = context.Universe.DefineType("int", universePos, types.MakeIntType(0, "int"))

	Int8Type  = context.Universe.DefineType("int8", universePos, types.MakeIntType(8, "int8"))
	Int16Type = context.Universe.DefineType("int16", universePos, types.MakeIntType(16, "int16"))
	Int32Type = context.Universe.DefineType("int32", universePos, types.MakeIntType(32, "int32"))
	Int64Type = context.Universe.DefineType("int64", universePos, types.MakeIntType(64, "int64"))
)

var (
	UintType    = context.Universe.DefineType("uint", universePos, types.MakeUIntType(0, false, "uint"))
	UintptrType = context.Universe.DefineType("uintptr", universePos, types.MakeUIntType(0, true, "uintptr"))

	Uint8Type  = context.Universe.DefineType("uint8", universePos, types.MakeUIntType(8, false, "uint8"))
	Uint16Type = context.Universe.DefineType("uint16", universePos, types.MakeUIntType(16, false, "uint16"))
	Uint32Type = context.Universe.DefineType("uint32", universePos, types.MakeUIntType(32, false, "uint32"))
	Uint64Type = context.Universe.DefineType("uint64", universePos, types.MakeUIntType(64, false, "uint64"))
)

var (
	Float32Type = context.Universe.DefineType("float32", universePos, types.MakeFloatType(32, "float32"))
	Float64Type = context.Universe.DefineType("float64", universePos, types.MakeFloatType(64, "float64"))
)

var (
	StringType = context.Universe.DefineType("string", universePos, &types.StringType{})
)
