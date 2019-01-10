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
	"github.com/xav/go-script/context"
	"github.com/xav/go-script/types"
)

var (
	// BoolType is the definition of the 'bool' built-in type
	BoolType = context.Universe.DefineType("bool", universePos, &types.BoolType{})
)

var (
	// IntType is the definition of the 'int' built-in type
	IntType = context.Universe.DefineType("int", universePos, types.MakeIntType(0, "int"))

	// Int8Type is the definition of the 'int8' built-in type
	Int8Type = context.Universe.DefineType("int8", universePos, types.MakeIntType(8, "int8"))
	// Int16Type is the definition of the 'int16' built-in type
	Int16Type = context.Universe.DefineType("int16", universePos, types.MakeIntType(16, "int16"))
	// Int32Type is the definition of the 'int32' built-in type
	Int32Type = context.Universe.DefineType("int32", universePos, types.MakeIntType(32, "int32"))
	// Int64Type is the definition of the 'int64' built-in type
	Int64Type = context.Universe.DefineType("int64", universePos, types.MakeIntType(64, "int64"))
)

var (
	// UintType is the definition of the 'uint' built-in type
	UintType = context.Universe.DefineType("uint", universePos, types.MakeUIntType(0, false, "uint"))
	// UintptrType is the definition of the 'uintptr' built-in type
	UintptrType = context.Universe.DefineType("uintptr", universePos, types.MakeUIntType(0, true, "uintptr"))

	// Uint8Type is the definition of the 'uint8' built-in type
	Uint8Type = context.Universe.DefineType("uint8", universePos, types.MakeUIntType(8, false, "uint8"))
	// Uint16Type is the definition of the 'uint16' built-in type
	Uint16Type = context.Universe.DefineType("uint16", universePos, types.MakeUIntType(16, false, "uint16"))
	// Uint32Type is the definition of the 'uint32' built-in type
	Uint32Type = context.Universe.DefineType("uint32", universePos, types.MakeUIntType(32, false, "uint32"))
	// Uint64Type is the definition of the 'uint64' built-in type
	Uint64Type = context.Universe.DefineType("uint64", universePos, types.MakeUIntType(64, false, "uint64"))
)

var (
	// Float32Type is the definition of the 'float32' built-in type
	Float32Type = context.Universe.DefineType("float32", universePos, types.MakeFloatType(32, "float32"))
	// Float64Type is the definition of the 'float64' built-in type
	Float64Type = context.Universe.DefineType("float64", universePos, types.MakeFloatType(64, "float64"))
)

var (
	// StringType is the definition of the 'string' built-in type
	StringType = context.Universe.DefineType("string", universePos, &types.StringType{})
)
