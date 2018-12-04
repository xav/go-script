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
	"github.com/xav/go-script/values"
	"github.com/xav/go-script/vm"
)

// Expr is the result of compiling an expression.
// It stores the type of the expression and its evaluator function.
type Expr struct {
	*ExprInfo

	// The type of the expression
	ExprType vm.Type

	// Evaluate this node as the given type.
	eval interface{}

	// Map index expressions permit special forms of assignment, for which we need to know the Map and key.
	evalMapValue func(t *vm.Thread) (values.Map, interface{})

	// Evaluate to the "address" of this value; that is, the settable Value object.
	// nil for expressions whose address cannot be taken.
	evalAddr func(t *vm.Thread) vm.Value

	// Execute this expression as a statement.
	// Only expressions that are valid expression statements should set this.
	Exec vm.CodeInstruction

	// If this expression is a type, this is its compiled type.
	// This is only permitted in the function position of a call expression. In this case, exprType should be nil.
	valType vm.Type

	// A short string describing this expression for error messages.
	desc string
}
