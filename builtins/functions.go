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

import "github.com/xav/go-script/types"

var (
	// CloseType is the definition of the close(c) built-in function (https://golang.org/ref/spec#Close)
	CloseType = &types.FuncType{Builtin: "close"}

	// LenType is the definition of the len(s) built-in function (https://golang.org/ref/spec#Length_and_capacity)
	LenType = &types.FuncType{Builtin: "len"}
	// CapType is the definition of the cap(s) built-in function (https://golang.org/ref/spec#Length_and_capacity)
	CapType = &types.FuncType{Builtin: "cap"}

	// NewType is the definition of the new built-in function (https://golang.org/ref/spec#Allocation)
	NewType = &types.FuncType{Builtin: "new"}

	// MakeType is the definition of the make built-in function (https://golang.org/ref/spec#Making_slices_maps_and_channels)
	MakeType = &types.FuncType{Builtin: "make"}

	// AppendType is the definition of the append built-in function (https://golang.org/ref/spec#Appending_and_copying_slices)
	AppendType = &types.FuncType{Builtin: "append"}
	// CopyType is the definition of the copy built-in function (https://golang.org/ref/spec#Appending_and_copying_slices)
	CopyType = &types.FuncType{Builtin: "copy"}

	// DeleteType is the definition of the delete built-in function (https://golang.org/ref/spec#Deletion_of_map_elements)
	DeleteType = &types.FuncType{Builtin: "delete"}

	// ComplexType is the definition of the complex built-in function (https://golang.org/ref/spec#Complex_numbers)
	ComplexType = &types.FuncType{Builtin: "complex"}
	// RealType is the definition of the real built-in function (https://golang.org/ref/spec#Complex_numbers)
	RealType = &types.FuncType{Builtin: "real"}
	// ImagType is the definition of the imag built-in function (https://golang.org/ref/spec#Complex_numbers)
	ImagType = &types.FuncType{Builtin: "imag"}

	// PanicType is the definition of the panic built-in function (https://golang.org/ref/spec#Handling_panics)
	PanicType = &types.FuncType{Builtin: "panic"}
	// RecoverType is the definition of the recover built-in function (https://golang.org/ref/spec#Handling_panics)
	RecoverType = &types.FuncType{Builtin: "recover"}

	// PrintType is the definition of the print built-in function (https://golang.org/ref/spec#Bootstrapping)
	PrintType = &types.FuncType{Builtin: "print"}
	// PrintlnType is the definition of the println built-in function (https://golang.org/ref/spec#Bootstrapping)
	PrintlnType = &types.FuncType{Builtin: "println"}
)
