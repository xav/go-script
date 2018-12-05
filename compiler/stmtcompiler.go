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
	"go/token"
)

const (
	returnPC = ^uint(0)
	badPC    = ^uint(1)
)

// stmtCompiler is used for the compilation of individual statements.
type stmtCompiler struct {
	*BlockCompiler           // The BlockClompiler for the block enclosing this statement.
	pos            token.Pos // The position of this statement.
	stmtLabel      *Label    // This statement's label, or nil if it is not labeled.
}

func (sc *stmtCompiler) error(format string, args ...interface{}) {
	sc.errorAt(sc.pos, format, args...)
}

func (sc *stmtCompiler) compile(s ast.Stmt) {
	if sc.Block.Inner != nil {
		logger.Panic().Msg("Child scope still entered")
	}

	notImplemented := false
	switch s := s.(type) {
	case *ast.BadStmt:
		sc.SilentErrors++ // Error already reported by parser.

	case *ast.AssignStmt:
		notImplemented = true

	case *ast.BlockStmt:
		notImplemented = true

	case *ast.BranchStmt:
		notImplemented = true

	case *ast.CaseClause:
		// sc.error("case clause outside switch")

	case *ast.CommClause:
		// sc.error("case clause outside select")
		notImplemented = true

	case *ast.DeclStmt:
		notImplemented = true

	case *ast.DeferStmt:
		notImplemented = true

	case *ast.EmptyStmt:
		// Do nothing.

	case *ast.ExprStmt:
		notImplemented = true

	case *ast.ForStmt:
		notImplemented = true

	case *ast.GoStmt:
		notImplemented = true

	case *ast.IfStmt:
		notImplemented = true

	case *ast.IncDecStmt:
		notImplemented = true

	case *ast.LabeledStmt:
		notImplemented = true

	case *ast.RangeStmt:
		notImplemented = true

	case *ast.ReturnStmt:
		notImplemented = true

	case *ast.SelectStmt:
		notImplemented = true

	case *ast.SendStmt:
		notImplemented = true

	case *ast.SwitchStmt:
		notImplemented = true

	case *ast.TypeSwitchStmt:
		notImplemented = true

	default:
		logger.Panic().
			Str("type", fmt.Sprintf("%T", s)).
			Msg("Unexpected ast node type")
	}

	if notImplemented {
		sc.error("%T statement node not implemented", s)
	}

	if sc.Block.Inner != nil {
		logger.Panic().Msg("Forgot to exit child scope")
	}
}
