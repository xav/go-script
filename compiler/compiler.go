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
	"github.com/xav/go-script/types"
	"github.com/xav/go-script/values"
	"github.com/xav/go-script/vm"
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
	ec := &ExprCompiler{cc, b, constant}
	nerr := cc.NumError()

	e := ec.compile(expr, false)
	if e == nil && nerr == cc.NumError() {
		logger.Panic().Msg("expression compilation failed without reporting errors")
	}

	return e
}

func (cc *Compiler) compileType(b *context.Block, typ ast.Expr) vm.Type {
	tc := &typeCompiler{
		Compiler:  cc,
		block:     b,
		lateCheck: noLateCheck,
	}
	t := tc.compileType(typ, false)
	if !tc.lateCheck() {
		return nil
	}

	return t
}

func (cc *Compiler) compileFuncType(b *context.Block, typ *ast.FuncType) *types.FuncDecl {
	tc := &typeCompiler{
		Compiler:  cc,
		block:     b,
		lateCheck: noLateCheck,
	}
	fd := tc.compileFuncType(typ, false)
	if fd != nil {
		if !tc.lateCheck() {
			fd = nil
		}
	}
	return fd
}

func (cc *Compiler) compileFunc(b *context.Block, decl *types.FuncDecl, body *ast.BlockStmt) func(*vm.Thread) values.Func {
	// Create body scope
	// The scope of a parameter or result is the body of the corresponding function.
	bodyScope := b.EnterChildScope()
	defer bodyScope.Exit()

	for i, t := range decl.Type.In {
		if decl.InNames[i] != nil {
			bodyScope.DefineVar(decl.InNames[i].Name, decl.InNames[i].Pos(), t)
		} else {
			bodyScope.DefineTemp(t)
		}
	}
	for i, t := range decl.Type.Out {
		if decl.OutNames[i] != nil {
			bodyScope.DefineVar(decl.OutNames[i].Name, decl.OutNames[i].Pos(), t)
		} else {
			bodyScope.DefineTemp(t)
		}
	}

	// Create block context
	cb := NewCodeBuf()
	fc := &FuncCompiler{
		Compiler:     cc,
		FnType:       decl.Type,
		OutVarsNamed: len(decl.OutNames) > 0 && decl.OutNames[0] != nil,
		CodeBuf:      cb,
		Flow:         NewFlowBuf(cb),
		Labels:       make(map[string]*Label),
	}
	bc := &BlockCompiler{
		FuncCompiler: fc,
		Block:        bodyScope.Block,
	}

	// Compile body
	nerr := cc.NumError()
	bc.compileStmts(body)
	fc.CheckLabels()
	if nerr != cc.NumError() {
		return nil
	}

	// Check that the body returned if necessary.
	if len(decl.Type.Out) > 0 && fc.Flow.reachesEnd(0) {
		cc.errorAt(body.Rbrace, "function ends without a return statement")
		return nil
	}

	code := fc.Get()
	maxVars := bodyScope.MaxVars
	return func(t *vm.Thread) values.Func {
		return &vm.EvalFunc{
			Outer:     t.Frame,
			FrameSize: maxVars,
			Code:      code,
		}
	}
}

// compileAssign compiles an assignment operation without the full generality of an assignCompiler.
// See assignCompiler for a description of the arguments.
func (cc *Compiler) compileAssign(pos token.Pos, b *context.Block, lt vm.Type, rs []*Expr, errOp, errPosName string) func(vm.Value, *vm.Thread) {
	panic("NOT IMPLEMENTED")
}

// Type check the RHS of an assignment, returning a new assignCompiler and indicating if the type check succeeded.
// This always returns an assignCompiler with rmt set, but if type checking fails, slots in the MultiType may be nil.
// If rs contains nil's, type checking will fail and these expressions given a nil type.
func (cc *Compiler) checkAssign(pos token.Pos, rs []*Expr, errOp, errPosName string) (*assignCompiler, bool) {
	c := &assignCompiler{
		Compiler:   cc,
		pos:        pos,
		rs:         rs,
		errOp:      errOp,
		errPosName: errPosName,
	}

	// Is this an unpack?
	if len(rs) == 1 && rs[0] != nil {
		if rmt, isUnpack := rs[0].ExprType.(*types.MultiType); isUnpack {
			c.rmt = rmt
			c.isUnpack = true
			return c, true
		}
	}

	// Create MultiType for RHS and check that all RHS expressions are single-valued.
	rts := make([]vm.Type, len(rs))
	ok := true
	for i, r := range rs {
		if r == nil {
			ok = false
			continue
		}

		if _, isMT := r.ExprType.(*types.MultiType); isMT {
			r.error("multi-valued expression not allowed in %s", errOp)
			ok = false
			continue
		}

		rts[i] = r.ExprType
	}

	c.rmt = types.NewMultiType(rts)
	return c, ok
}
