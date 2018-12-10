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
	AppendType = &types.FuncType{Builtin: "append"}
	CapType    = &types.FuncType{Builtin: "cap"}
	CloseType  = &types.FuncType{Builtin: "close"}
	// TODO: removed closed?
	ClosedType = &types.FuncType{Builtin: "closed"}
	CopyType   = &types.FuncType{Builtin: "copy"}
	// TODO: delete
	// TODO: imag
	LenType     = &types.FuncType{Builtin: "len"}
	MakeType    = &types.FuncType{Builtin: "make"}
	NewType     = &types.FuncType{Builtin: "new"}
	PanicType   = &types.FuncType{Builtin: "panic"}
	PrintType   = &types.FuncType{Builtin: "print"}
	PrintlnType = &types.FuncType{Builtin: "println"}
	// TODO: real
	// TODO: recover
)
