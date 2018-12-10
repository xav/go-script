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

package compiler

import (
	"go/token"

	"github.com/xav/go-script/builtins"
	"github.com/xav/go-script/context"
	"github.com/xav/go-script/types"
	"github.com/xav/go-script/vm"
)

// AssignCompiler compiles assignment operations.
// Anything other than short declarations should use the compileAssign wrapper.
//
// There are three valid types of assignment:
// 1) T = T
//    Assigning a single expression with single-valued type to a
//    single-valued type.
// 2) MT = T, T, ...
//    Assigning multiple expressions with single-valued types to a
//    multi-valued type.
// 3) MT = MT
//    Assigning a single expression with multi-valued type to a
//    multi-valued type.
type AssignCompiler struct {
	*Compiler
	pos         token.Pos
	rs          []*Expr          // The RHS expressions.  This may include nil's for expressions that failed to compile.
	rmt         *types.MultiType // The (possibly unary) MultiType of the RHS.
	isUnpack    bool             // Whether this is an unpack assignment (case 3).
	allowMap    bool             // Whether map special assignment forms are allowed.
	isMapUnpack bool             // Whether this is a "r, ok = a[x]" assignment.
	errOp       string           // The operation name to use in error messages, such as "assignment" or "function call".
	errPosName  string           // The name to use for positions in error messages, such as "argument".
}

func (ac *AssignCompiler) compile(b *context.Block, lt vm.Type) func(vm.Value, *vm.Thread) {
	panic("NOT IMPLEMENTED")
}

func (ac *AssignCompiler) allowMapForms(nls int) {
	ac.allowMap = true

	// Update unpacking info if this is 'r, ok = a[x]'
	if nls == 2 && len(ac.rs) == 1 && ac.rs[0] != nil && ac.rs[0].evalMapValue != nil {
		ac.isUnpack = true
		ac.isMapUnpack = true
		ac.rmt = types.NewMultiType([]vm.Type{
			ac.rs[0].ExprType,
			builtins.BoolType,
		})
	}
}
