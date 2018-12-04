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

package script

import (
	"github.com/xav/go-script/compiler"
	"github.com/xav/go-script/vm"
)

type ExprCode struct {
	world *World
	expr  *compiler.Expr
	eval  func(vm.Value, *vm.Thread)
}

func (e *ExprCode) Type() vm.Type {
	return e.expr.ExprType
}

func (e *ExprCode) Run() (vm.Value, error) {
	t := new(vm.Thread)
	t.Frame = e.world.scope.NewFrame(nil)

	switch e.expr.ExprType.(type) {
	// case *types.BigIntType:
	// 	return &values.BigIntV{e.expr.AsBigInt()()}, nil
	// case *types.BigFloatType:
	// 	return &values.BigFloatV{e.expr.AsBigFloat()()}, nil
	}

	v := e.expr.ExprType.Zero()
	eval := e.eval
	err := t.Try(func(t *vm.Thread) {
		eval(v, t)
	})

	return v, err
}
