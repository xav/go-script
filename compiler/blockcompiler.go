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
	"go/ast"

	"github.com/xav/go-script/context"
)

// BlockCompiler captures information used throughout the compilation of a single block within a function.
type BlockCompiler struct {
	*FuncCompiler                // The FuncCompiler for the function enclosing this block.
	Parent        *BlockCompiler // The BlockCompiler for the block enclosing this one, or nil for a function-level block.
	Block         *context.Block // The block definition
	Label         *Label         // The label of this block, used for finding break and continue labels.
}

// CompileStmt compiles the specified statement within the block.
func (bc *BlockCompiler) CompileStmt(stmt ast.Stmt) {
	sc := &stmtCompiler{
		BlockCompiler: bc,
		pos:           stmt.Pos(),
		stmtLabel:     nil,
	}
	sc.compile(stmt)
}

func (bc *BlockCompiler) compileStmts(block *ast.BlockStmt) {
	if block == nil || block.List == nil {
		return
	}

	for _, sub := range block.List {
		bc.CompileStmt(sub)
	}
}

func (bc *BlockCompiler) enterChild() *BlockCompiler {
	block := bc.Block.EnterChild()
	return &BlockCompiler{
		FuncCompiler: bc.FuncCompiler,
		Block:        block,
		Parent:       bc,
	}
}
func (bc *BlockCompiler) exit() {
	bc.Block.Exit()
}
