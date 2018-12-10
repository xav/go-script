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

func (x *Expr) resolveIdeal(t vm.Type) *Expr { panic("NOT IMPLEMENTED") }

////////////////////////////////////////////////////////////////////////////////
// "As" functions retrieve evaluator functions from an expr, panicking if
// the requested evaluator has the wrong type.

// asValue returns a closure around a Value, according to the type of the underlying expression
func (x *Expr) asValue() func(t *vm.Thread) vm.Value      { panic("NOT IMPLEMENTED") }
func (x *Expr) asInterface() func(*vm.Thread) interface{} { panic("NOT IMPLEMENTED") }

func (x *Expr) asPackage() func(*vm.Thread) values.PackageValue {
	return x.eval.(func(*vm.Thread) values.PackageValue)
}
func (x *Expr) asBool() func(*vm.Thread) bool {
	return x.eval.(func(*vm.Thread) bool)
}
func (x *Expr) asUint() func(*vm.Thread) uint64 {
	return x.eval.(func(*vm.Thread) uint64)
}
func (x *Expr) asInt() func(*vm.Thread) int64 {
	return x.eval.(func(*vm.Thread) int64)
}
func (x *Expr) asFloat() func(*vm.Thread) float64 {
	return x.eval.(func(*vm.Thread) float64)
}
func (x *Expr) asString() func(*vm.Thread) string {
	return x.eval.(func(*vm.Thread) string)
}
func (x *Expr) asArray() func(*vm.Thread) values.ArrayValue {
	return x.eval.(func(*vm.Thread) values.ArrayValue)
}
func (x *Expr) asStruct() func(*vm.Thread) values.StructValue {
	return x.eval.(func(*vm.Thread) values.StructValue)
}
func (x *Expr) asPtr() func(*vm.Thread) vm.Value {
	return x.eval.(func(*vm.Thread) vm.Value)
}
func (x *Expr) asFunc() func(*vm.Thread) values.Func {
	return x.eval.(func(*vm.Thread) values.Func)
}
func (x *Expr) asSlice() func(*vm.Thread) values.Slice {
	return x.eval.(func(*vm.Thread) values.Slice)
}
func (x *Expr) asMap() func(*vm.Thread) values.Map {
	return x.eval.(func(*vm.Thread) values.Map)
}
func (x *Expr) asMulti() func(*vm.Thread) []vm.Value {
	return x.eval.(func(*vm.Thread) []vm.Value)
}

////////////////////////////////////////////////////////////////////////////////
