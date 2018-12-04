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

import "github.com/xav/go-script/vm"

type CodeBuf struct {
	instrs vm.Code
}

func NewCodeBuf() *CodeBuf {
	return &CodeBuf{make(vm.Code, 0, 16)}
}

func (b *CodeBuf) Push(instr func(*vm.Thread)) {
	b.instrs = append(b.instrs, instr)
}

func (b *CodeBuf) NextPC() uint {
	return uint(len(b.instrs))
}

func (b *CodeBuf) Get() vm.Code {
	// Freeze this buffer into an array of exactly the right size
	a := make(vm.Code, len(b.instrs))
	copy(a, b.instrs)
	return vm.Code(a)
}
