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
	"github.com/xav/go-script/vm"
)

func genAssign(lt vm.Type, r *Expr) func(lv vm.Value, t *vm.Thread) { panic("NOT IMPLEMENTED") }

func (x *Expr) genBinOpAdd(l, r *Expr)                         { panic("NOT IMPLEMENTED") }
func (x *Expr) genBinOpAnd(l, r *Expr)                         { panic("NOT IMPLEMENTED") }
func (x *Expr) genBinOpAndNot(l, r *Expr)                      { panic("NOT IMPLEMENTED") }
func (x *Expr) genBinOpEql(l, r *Expr)                         { panic("NOT IMPLEMENTED") }
func (x *Expr) genBinOpGeq(l, r *Expr)                         { panic("NOT IMPLEMENTED") }
func (x *Expr) genBinOpGtr(l, r *Expr)                         { panic("NOT IMPLEMENTED") }
func (x *Expr) genBinOpLeq(l, r *Expr)                         { panic("NOT IMPLEMENTED") }
func (x *Expr) genBinOpLogAnd(l, r *Expr)                      { panic("NOT IMPLEMENTED") }
func (x *Expr) genBinOpLogOr(l, r *Expr)                       { panic("NOT IMPLEMENTED") }
func (x *Expr) genBinOpLss(l, r *Expr)                         { panic("NOT IMPLEMENTED") }
func (x *Expr) genBinOpMul(l, r *Expr)                         { panic("NOT IMPLEMENTED") }
func (x *Expr) genBinOpNeq(l, r *Expr)                         { panic("NOT IMPLEMENTED") }
func (x *Expr) genBinOpOr(l, r *Expr)                          { panic("NOT IMPLEMENTED") }
func (x *Expr) genBinOpQuo(l, r *Expr)                         { panic("NOT IMPLEMENTED") }
func (x *Expr) genBinOpRem(l, r *Expr)                         { panic("NOT IMPLEMENTED") }
func (x *Expr) genBinOpShl(l, r *Expr)                         { panic("NOT IMPLEMENTED") }
func (x *Expr) genBinOpShr(l, r *Expr)                         { panic("NOT IMPLEMENTED") }
func (x *Expr) genBinOpSub(l, r *Expr)                         { panic("NOT IMPLEMENTED") }
func (x *Expr) genBinOpXor(l, r *Expr)                         { panic("NOT IMPLEMENTED") }
func (x *Expr) genConstant(v vm.Value)                         { panic("NOT IMPLEMENTED") }
func (x *Expr) genFuncCall(call func(t *vm.Thread) []vm.Value) { panic("NOT IMPLEMENTED") }
func (x *Expr) genIdentOp(level, index int)                    { panic("NOT IMPLEMENTED") }
func (x *Expr) genUnaryOpNeg(v *Expr)                          { panic("NOT IMPLEMENTED") }
func (x *Expr) genUnaryOpNot(v *Expr)                          { panic("NOT IMPLEMENTED") }
func (x *Expr) genUnaryOpXor(v *Expr)                          { panic("NOT IMPLEMENTED") }
func (x *Expr) genValue(vf func(*vm.Thread) vm.Value)          { panic("NOT IMPLEMENTED") }
