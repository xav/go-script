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
	"fmt"
	"go/ast"
	"go/scanner"
	"go/token"

	"github.com/xav/go-script/context"
)

// Compiler captures information used throughout a package compilation.
type Compiler struct {
	FSet         *token.FileSet
	Errors       *scanner.ErrorList
	NumErrors    int
	SilentErrors int
}

// Reports a compilation error att the specified position
func (cc *Compiler) errorAt(pos token.Pos, format string, args ...interface{}) {
	cc.Errors.Add(cc.FSet.Position(pos), fmt.Sprintf(format, args...))
	cc.NumErrors++
}

// NumError returns the total number of errors detected yet
func (cc *Compiler) NumError() int {
	return cc.NumErrors + cc.SilentErrors
}

func (cc *Compiler) CompileExpr(b *context.Block, constant bool, expr ast.Expr) *Expr {
	panic("NOT IMPLEMENTED")
}
