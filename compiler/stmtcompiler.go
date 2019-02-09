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
	"go/importer"
	"go/token"
	"math/big"
	"strconv"

	gotypes "go/types"

	"github.com/xav/go-script/builtins"
	"github.com/xav/go-script/context"
	"github.com/xav/go-script/types"
	"github.com/xav/go-script/values"
	"github.com/xav/go-script/vm"
)

const (
	returnPC = ^uint(0)
	badPC    = ^uint(1)
)

// stmtCompiler is used for the compilation of individual statements.
type stmtCompiler struct {
	*BlockCompiler           // The BlockCompiler for the block enclosing this statement.
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
		sc.compileAssignStmt(s)

	case *ast.BlockStmt:
		sc.compileBlockStmt(s)

	case *ast.BranchStmt:
		sc.compileBranchStmt(s)

	case *ast.CaseClause:
		// sc.error("case clause outside switch")

	case *ast.CommClause:
		notImplemented = true

	case *ast.DeclStmt:
		sc.compileDeclStmt(s)

	case *ast.DeferStmt:
		notImplemented = true

	case *ast.EmptyStmt:
		// Do nothing.

	case *ast.ExprStmt:
		sc.compileExprStmt(s)

	case *ast.ForStmt:
		sc.compileForStmt(s)

	case *ast.GoStmt:
		notImplemented = true

	case *ast.IfStmt:
		sc.compileIfStmt(s)

	case *ast.IncDecStmt:
		sc.compileIncDecStmt(s)

	case *ast.LabeledStmt:
		sc.compileLabeledStmt(s)

	case *ast.RangeStmt:
		sc.compileRangeStmt(s)

	case *ast.ReturnStmt:
		sc.compileReturnStmt(s)

	case *ast.SelectStmt:
		sc.compileSelectStmt(s)

	case *ast.SendStmt:
		notImplemented = true

	case *ast.SwitchStmt:
		sc.compileSwitchStmt(s)

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

// Statements //////////////////////////////////////////////////////////////////

func (sc *stmtCompiler) compileAssignStmt(stmt *ast.AssignStmt) {
	switch stmt.Tok {
	case token.ASSIGN, token.DEFINE:
		sc.doAssign(stmt.Lhs, stmt.Rhs, stmt.Tok, nil)

	default:
		sc.doAssignOp(stmt)
	}
}

func (sc *stmtCompiler) compileBlockStmt(stmt *ast.BlockStmt) {
	bc := sc.enterChild()
	bc.compileStmts(stmt)
	bc.exit()
}

func (sc *stmtCompiler) compileBranchStmt(stmt *ast.BranchStmt) {
	var pc *uint

	switch stmt.Tok {
	case token.BREAK:
		l := sc.findLexicalLabel(stmt.Label, func(l *Label) bool {
			return l.breakPC != nil
		}, "break", "for loop, switch, or select")
		if l == nil {
			return
		}
		pc = l.breakPC

	case token.CONTINUE:
		l := sc.findLexicalLabel(stmt.Label, func(l *Label) bool { return l.continuePC != nil }, "continue", "for loop")
		if l == nil {
			return
		}
		pc = l.continuePC

	case token.GOTO:
		l, ok := sc.Labels[stmt.Label.Name]
		if !ok {
			pc := badPC
			l = &Label{name: stmt.Label.Name, desc: "unresolved label", gotoPC: &pc, used: stmt.Pos()}
			sc.Labels[l.name] = l
		}

		pc = l.gotoPC
		sc.Flow.putGoto(stmt.Pos(), l.name, sc.Block)

	case token.FALLTHROUGH:
		sc.error("fallthrough outside switch")
		return

	default:
		logger.Panic().Msgf("Unexpected branch token %v", stmt.Tok)
	}

	sc.Flow.putBranching(false, pc)
	sc.Push(func(v *vm.Thread) { v.PC = *pc })
}

func (sc *stmtCompiler) compileDeclStmt(stmt *ast.DeclStmt) {
	switch decl := stmt.Decl.(type) {
	case *ast.BadDecl:
		sc.SilentErrors++ // Error already reported by parser.
		return

	case *ast.FuncDecl:
		sc.compileFuncDecl(decl)

	case *ast.GenDecl:
		sc.compileGenDecl(decl)

	default:
		logger.Panic().
			Str("type", fmt.Sprintf("%T", stmt.Decl)).
			Msg("Unexpected Decl type")
	}
}

func (sc *stmtCompiler) compileExprStmt(stmt *ast.ExprStmt) {
	bc := sc.enterChild()
	defer bc.exit()

	e := sc.CompileExpr(bc.Block, false, stmt.X)
	if e == nil {
		return
	}

	if e.Exec == nil {
		sc.error("%s cannot be used as expression statement", e.desc)
		return
	}

	sc.Push(e.Exec)
}

func (sc *stmtCompiler) compileForStmt(stmt *ast.ForStmt) {
	// Wrap the entire for in a block.
	bc := sc.enterChild()
	defer bc.exit()

	// Compile init statement, if any
	if stmt.Init != nil {
		bc.CompileStmt(stmt.Init)
	}

	bodyPC := badPC
	postPC := badPC
	checkPC := badPC
	endPC := badPC

	// Jump to condition check.
	// We generate slightly less code by placing the condition check after the body.
	sc.Flow.putBranching(false, &checkPC)
	sc.Push(func(v *vm.Thread) {
		v.PC = checkPC
	})

	// Compile body
	bodyPC = sc.NextPC()
	body := bc.enterChild()
	if sc.stmtLabel != nil {
		body.Label = sc.stmtLabel
	} else {
		body.Label = &Label{resolved: stmt.Pos()}
	}
	body.Label.desc = "for loop"
	body.Label.breakPC = &endPC
	body.Label.continuePC = &postPC
	body.compileStmts(stmt.Body)
	body.exit()

	// Compile post, if any
	postPC = sc.NextPC()
	if stmt.Post != nil {
		// TODO: Does the parser disallow short declarations in stmt.Post?
		bc.CompileStmt(stmt.Post)
	}

	// Compile condition check, if any
	checkPC = sc.NextPC()
	if stmt.Cond == nil {
		// If the condition is absent, it is equivalent to true.
		sc.Flow.putBranching(false, &bodyPC)
		sc.Push(func(v *vm.Thread) {
			v.PC = bodyPC
		})
	} else {
		e := bc.CompileExpr(bc.Block, false, stmt.Cond)
		switch {
		case e == nil:
			// Error reported by compileExpr
		case !e.ExprType.IsBoolean():
			sc.error("'for' condition must be boolean\n\t%v", e.ExprType)
		default:
			eval := e.asBool()
			sc.Flow.putBranching(true, &bodyPC)
			sc.Push(func(t *vm.Thread) {
				if eval(t) {
					t.PC = bodyPC
				}
			})
		}
	}

	endPC = sc.NextPC()
}

func (sc *stmtCompiler) compileIfStmt(stmt *ast.IfStmt) {
	// "The scope of any variables declared by [the init] statement extends to the end of the "if" statement
	// and the variables are initialized once before the statement is entered.""
	//
	// What this really means is that there's an implicit scope wrapping every if, for, and switch statement.
	// This is subtly different from what it actually says when there's a non-block else clause,
	// because that else claus has to execute in a scope that is *not* the surrounding scope.
	bc := sc.enterChild()
	defer bc.exit()

	// Compile init statement, if any
	if stmt.Init != nil {
		bc.CompileStmt(stmt.Init)
	}

	elsePC := badPC
	endPC := badPC

	// Compile condition, if any. If there is no condition, we fall through to the body.
	if stmt.Cond != nil {
		e := bc.CompileExpr(bc.Block, false, stmt.Cond)
		switch {
		case e == nil:
			// Error reported by compileExpr
		case !e.ExprType.IsBoolean():
			e.error("'if' condition must be boolean\n\t%v", e.ExprType)
		default:
			eval := e.asBool()
			sc.Flow.putBranching(true, &elsePC)
			sc.Push(func(t *vm.Thread) {
				if !eval(t) {
					t.PC = elsePC
				}
			})
		}
	}

	// Compile body
	body := bc.enterChild()
	body.compileStmts(stmt.Body)
	body.exit()

	// Compile else
	if stmt.Else != nil {
		// Skip over else if we executed the body
		sc.Flow.putBranching(false, &endPC)
		sc.Push(func(v *vm.Thread) { v.PC = endPC })
		elsePC = sc.NextPC()
		bc.CompileStmt(stmt.Else)
	} else {
		elsePC = sc.NextPC()
	}
	endPC = sc.NextPC()
}

func (sc *stmtCompiler) compileIncDecStmt(stmt *ast.IncDecStmt) {
	// Create temporary block for extractEffect
	bc := sc.enterChild()
	defer bc.exit()

	l := sc.CompileExpr(bc.Block, false, stmt.X)
	if l == nil {
		return
	}

	if l.evalAddr == nil {
		l.error("cannot assign to %s", l.desc)
		return
	}
	if !(l.ExprType.IsInteger() || l.ExprType.IsFloat()) {
		l.errorOpType(stmt.Tok, l.ExprType)
		return
	}

	var op token.Token
	var desc string
	switch stmt.Tok {
	case token.INC:
		op = token.ADD
		desc = "increment statement"
	case token.DEC:
		op = token.SUB
		desc = "decrement statement"
	default:
		logger.Panic().Msgf("Unexpected IncDec token %v", stmt.Tok)
	}

	effect, l := l.extractEffect(bc.Block, desc)

	one := l.newExpr(IdealIntType, "constant")
	one.pos = stmt.Pos()
	one.eval = func() *big.Int { return big.NewInt(1) }

	binop := l.compileBinaryExpr(op, l, one)
	if binop == nil {
		return
	}

	assign := sc.compileAssign(stmt.Pos(), bc.Block, l.ExprType, []*Expr{binop}, "", "")
	if assign == nil {
		logger.Panic().Msgf("compileAssign type check failed")
	}

	lf := l.evalAddr
	sc.Push(func(v *vm.Thread) {
		effect(v)
		assign(lf(v), v)
	})
}

func (sc *stmtCompiler) compileLabeledStmt(stmt *ast.LabeledStmt) {
	// Define label
	l, ok := sc.Labels[stmt.Label.Name]
	if ok {
		if l.resolved.IsValid() {
			sc.error("label %s redeclared in this block\n\tprevious declaration at %s", stmt.Label.Name, sc.FSet.Position(l.resolved))
		}
	} else {
		pc := badPC
		l = &Label{name: stmt.Label.Name, gotoPC: &pc}
		sc.Labels[l.name] = l
	}
	l.desc = "regular label"
	l.resolved = stmt.Pos()

	// Set goto PC
	*l.gotoPC = sc.NextPC()

	// Define flow entry so we can check for jumps over declarations.
	sc.Flow.putLabel(l.name, sc.Block)

	// Compile the statement.
	sc = &stmtCompiler{sc.BlockCompiler, stmt.Stmt.Pos(), l}
	sc.compile(stmt.Stmt)
}

func (sc *stmtCompiler) genRangeKeyInit(xt vm.Type, r *Expr) func(kv vm.Value, t *vm.Thread) {
	switch t := xt.(type) {
	case *types.ArrayType:
		return func(lv vm.Value, t *vm.Thread) {
			lv.(values.IntValue).Set(t, 0)
		}
	case *types.SliceType:
		return func(lv vm.Value, t *vm.Thread) {
			lv.(values.IntValue).Set(t, 0)
		}
	case *types.StringType:
		return func(lv vm.Value, t *vm.Thread) {
			lv.(values.IntValue).Set(t, 0)
		}
	case *types.MapType:
		return func(lv vm.Value, t *vm.Thread) {
			lv.(values.IntValue).Set(t, 0)
		}
	case *types.NamedType:
		sc.error("Named types are not implemented")
		return nil
	case *types.ChanType:
		sc.error("Channels are not implemented")
		return nil
	case *types.PtrType:
		return sc.genRangeKeyInit(t.Elem, r)
	default:
		sc.error("cannot range over expression %T", t)
		return nil
	}
}

func (sc *stmtCompiler) getRangeExprType(t vm.Type, p token.Pos) (vm.Type, vm.Type) {
	switch t := t.(type) {
	case *types.ArrayType:
		return builtins.IntType, t.Elem.Lit()
	case *types.SliceType:
		return builtins.IntType, t.Elem.Lit()
	case *types.StringType:
		return builtins.IntType, builtins.Uint8Type
	case *types.MapType:
		return t.Key.Lit(), t.Elem.Lit()
	case *types.NamedType:
		return t.Lit(), nil
	case *types.ChanType:
		sc.errorAt(p, "Channels are not implemented")
		return t, nil
	case *types.PtrType:
		return sc.getRangeExprType(t.Elem, p)
	default:
		sc.errorAt(p, "cannot range over expression %T", t)
		return nil, nil
	}
}

func (sc *stmtCompiler) compileRangeStmt(stmt *ast.RangeStmt) {
	// Wrap the entire for in a block.
	bc := sc.enterChild()
	defer bc.exit()

	// Compile range expression
	rx := bc.CompileExpr(bc.Block, false, stmt.X)
	if rx == nil {
		return
	}

	rx = rx.derefArray()

	var (
		kt vm.Type
		vt vm.Type
		// intIndex bool
	)

	// Type check object
	switch rt := rx.ExprType.Lit().(type) {
	case *types.ArrayType:
		kt = builtins.IntType
		vt = rt.Elem
	case *types.SliceType:
		kt = builtins.IntType
		vt = rt.Elem
	case *types.StringType:
		kt = builtins.IntType
		vt = builtins.Uint8Type
	case *types.MapType:
		kt = rt.Key
		vt = rt.Elem
	case *types.ChanType:
		sc.errorAt(rx.pos, "Channels are not implemented")
		return
	default:
		sc.errorAt(rx.pos, "cannot range over expression %T", rt)
		return
	}

	//TODO: DEBUG
	print(kt)
	print(vt)

	// bodyPC := badPC
	// iteratePC := badPC
	// endPC := badPC

	// ac, ok := sc.checkAssign(sc.pos, rx, "assignment", "value")
	// c := &assignCompiler{
	// 	Compiler:   sc.Compiler,
	// 	pos:        rx.pos,
	// 	rs:         rx,
	// 	errOp:      errOp,
	// 	errPosName: errPosName,
	// }

	// Compile init statement
	// sc.Push(func(t *vm.Thread) {
	// 	temp := lmt.Zero().(values.MultiV)
	// 	assign(temp, t)
	// 	// Copy to destination
	// 	for i := 0; i < n; i++ {
	// 		// TODO: Need to evaluate LHS before RHS
	// 		lfs[i](t).Assign(t, temp[i])
	// 	}
	// })
}

func (sc *stmtCompiler) compileRangeStmtX(stmt *ast.RangeStmt) {
	// Wrap the entire for in a block.
	bc := sc.enterChild()
	defer bc.exit()

	// Compile range expression
	ex := bc.CompileExpr(bc.Block, false, stmt.X)
	if ex == nil {
		return
	}

	// Infer k and v types.
	kt, vt := sc.getRangeExprType(ex.ExprType.Lit(), ex.pos)
	if kt == nil {
		return
	}
	print(kt, vt)

	// Compile init statement ////////////
	//TODO: Check if it's a definition or assignment

	k, ok := stmt.Key.(*ast.Ident)
	if !ok {
		sc.errorAt(stmt.Key.Pos(), "key must be an identifier in range expression")
		return
	}

	var v *ast.Ident
	if stmt.Value != nil {
		v, ok = stmt.Value.(*ast.Ident)
		if !ok {
			sc.errorAt(stmt.Value.Pos(), "value must be an identifier in range expression")
			return
		}
	}

	if v != nil {
		// lf := ls[0].evalAddr
		sc.Push(func(t *vm.Thread) {
			// assign(lf(t), t)
		})
	} else {
		lfs := make([]func(*vm.Thread) vm.Value, 2)
		// for i, l := range ls {
		// 	lfs[i] = l.evalAddr
		// }

		sc.Push(func(t *vm.Thread) {
			dest := make([]vm.Value, 2)
			for i, lf := range lfs {
				dest[i] = lf(t)
			}
			// assign(values.MultiV(dest), t)
		})
	}

	print(k)
	print(v)

	bodyPC := badPC
	iteratePC := badPC
	endPC := badPC

	// Compile body
	bodyPC = sc.NextPC()
	body := bc.enterChild()
	if sc.stmtLabel != nil {
		body.Label = sc.stmtLabel
	} else {
		body.Label = &Label{resolved: stmt.Pos()}
	}
	body.Label.desc = "for range loop"
	body.Label.breakPC = &endPC
	body.Label.continuePC = &iteratePC
	body.compileStmts(stmt.Body)
	body.exit()

	// Compile iteration
	iteratePC = sc.NextPC()
	//TODO: compile iteration
	sc.Push(func(t *vm.Thread) {
		t.PC = bodyPC
	})

	endPC = sc.NextPC()
}

func (sc *stmtCompiler) compileReturnStmt(stmt *ast.ReturnStmt) {
	if sc.FnType == nil {
		sc.error("cannot return at the top level")
		return
	}

	if len(stmt.Results) == 0 && (len(sc.FnType.Out) == 0 || sc.OutVarsNamed) {
		// Simple case. Simply exit from the function.
		sc.Flow.putTerm()
		sc.Push(func(v *vm.Thread) { v.PC = returnPC })
		return
	}

	bc := sc.enterChild()
	defer bc.exit()

	// Compile expressions
	bad := false
	rs := make([]*Expr, len(stmt.Results))
	for i, re := range stmt.Results {
		rs[i] = sc.CompileExpr(bc.Block, false, re)
		if rs[i] == nil {
			bad = true
		}
	}
	if bad {
		return
	}

	// Create assigner

	// If the expression list in the "return" statement is a single call to a multi-valued function,
	// the values returned from the called function will be returned from this one.
	assign := sc.compileAssign(stmt.Pos(), bc.Block, types.NewMultiType(sc.FnType.Out), rs, "return", "value")

	// "The result types of the current function and the called function must match."
	// Match is fuzzy. It should say that they must be assignment compatible.

	// Compile
	start := len(sc.FnType.In)
	nout := len(sc.FnType.Out)
	sc.Flow.putTerm()
	sc.Push(func(t *vm.Thread) {
		assign(values.MultiV(t.Frame.Vars[start:start+nout]), t)
		t.PC = returnPC
	})
}

func (sc *stmtCompiler) compileSelectStmt(stmt *ast.SelectStmt) {
	panic("Select statement not implemented")
}

// Declarations ////////////////////////////////////////////////////////////////

func (sc *stmtCompiler) compileSwitchStmt(stmt *ast.SwitchStmt) {
	// Create implicit scope around switch
	bc := sc.enterChild()
	defer bc.exit()

	// Compile init statement, if any
	if stmt.Init != nil {
		bc.CompileStmt(stmt.Init)
	}

	// Compile condition, if any, and extract its effects
	var cond *Expr
	condbc := bc.enterChild()
	if stmt.Tag != nil {
		e := condbc.CompileExpr(condbc.Block, false, stmt.Tag)
		if e != nil {
			var effect func(*vm.Thread)
			effect, cond = e.extractEffect(condbc.Block, "switch")
			sc.Push(effect)
		}
	}

	// Count cases
	ncases := 0
	hasDefault := false
	for _, c := range stmt.Body.List {
		clause, ok := c.(*ast.CaseClause)
		if !ok {
			sc.errorAt(clause.Pos(), "switch statement must contain case clauses")
			continue
		}
		if clause.List == nil {
			if hasDefault {
				sc.errorAt(clause.Pos(), "switch statement contains more than one default case")
			}
			hasDefault = true
		} else {
			ncases += len(clause.List)
		}
	}

	// Compile case expressions
	cases := make([]func(*vm.Thread) bool, ncases)
	i := 0
	for _, c := range stmt.Body.List {
		clause, ok := c.(*ast.CaseClause)
		if !ok {
			continue
		}
		for _, v := range clause.List {
			e := condbc.CompileExpr(condbc.Block, false, v)
			switch {
			case e == nil:
				// Error reported by compileExpr
			case cond == nil && !e.ExprType.IsBoolean():
				sc.errorAt(v.Pos(), "'case' condition must be boolean")
			case cond == nil:
				cases[i] = e.asBool()
			case cond != nil:
				// Create comparison
				// TODO: This produces bad error messages
				compare := e.compileBinaryExpr(token.EQL, cond, e)
				if compare != nil {
					cases[i] = compare.asBool()
				}
			}
			i++
		}
	}

	// Emit condition
	casePCs := make([]*uint, ncases+1)
	endPC := badPC

	sc.Flow.put(false, false, casePCs)
	sc.Push(func(t *vm.Thread) {
		for i, c := range cases {
			if c(t) {
				t.PC = *casePCs[i]
				return
			}
		}
		t.PC = *casePCs[ncases]
	})
	condbc.exit()

	// Compile cases
	i = 0
	for _, c := range stmt.Body.List {
		clause, ok := c.(*ast.CaseClause)
		if !ok {
			continue
		}

		// Save jump PC's
		pc := sc.NextPC()
		if clause.List != nil {
			for _ = range clause.List {
				casePCs[i] = &pc
				i++
			}
		} else {
			// Default clause
			casePCs[ncases] = &pc
		}

		// Compile body
		fall := false
		for j, s1 := range clause.Body {
			if br, ok := s1.(*ast.BranchStmt); ok && br.Tok == token.FALLTHROUGH {
				// println("Found fallthrough");
				// It may be used only as the final non-empty statement in a case or default clause in an expression "switch" statement.
				for _, s2 := range clause.Body[j+1:] {
					// XXX(Spec) 6g also considers empty blocks to be empty statements.
					if _, ok := s2.(*ast.EmptyStmt); !ok {
						sc.errorAt(s1.Pos(), "fallthrough statement must be final statement in case")
						break
					}
				}
				fall = true
			} else {
				bc.CompileStmt(s1)
			}
		}
		// Jump out of switch, unless there was a fallthrough
		if !fall {
			sc.Flow.putBranching(false, &endPC)
			sc.Push(func(v *vm.Thread) { v.PC = endPC })
		}
	}

	// Get end PC
	endPC = sc.NextPC()
	if !hasDefault {
		casePCs[ncases] = &endPC
	}
}

func (sc *stmtCompiler) compileFuncDecl(fd *ast.FuncDecl) {
	if !sc.Block.Global {
		logger.Panic().Msg("FuncDecl at statement level")
	}

	if fd.Body == nil {
		sc.errorAt(fd.Name.Pos(), "function %s declared without body.", fd.Name.Name)
		return
	}

	ftd := sc.compileFuncType(sc.Block, fd.Type)
	if ftd == nil {
		return
	}

	// Declare and initialize v before compiling func so that body can refer to itself.
	c, prev := sc.Block.DefineConst(fd.Name.Name, sc.pos, ftd.Type, ftd.Type.Zero())
	if prev != nil {
		pos := prev.Pos()
		if pos.IsValid() {
			sc.errorAt(fd.Name.Pos(), "identifier %s redeclared in this block\n\tprevious declaration at %s", fd.Name.Name, sc.FSet.Position(pos))
		} else {
			sc.errorAt(fd.Name.Pos(), "identifier %s redeclared in this block", fd.Name.Name)
		}
	}

	// TODO: If fd.Body is nil, we might be dealing with an assembler function.
	fn := sc.compileFunc(sc.Block, ftd, fd.Body)
	if c == nil || fn == nil {
		// when the compilation failed, remove the func identifier from the block definitions.
		sc.Block.Undefine(fd.Name.Name)
		return
	}

	var zeroThread vm.Thread
	c.Value.(values.FuncValue).Set(nil, fn(&zeroThread))
}

func (sc *stmtCompiler) compileGenDecl(gd *ast.GenDecl) {
	switch gd.Tok {
	case token.CONST:
		sc.compileConstDecl(gd)
	case token.IMPORT:
		sc.compileImportDecl(gd)
	case token.TYPE:
		sc.compileTypeDecl(gd)
	case token.VAR:
		sc.compileVarDecl(gd)
	default:
		logger.Panic().
			Str("token", fmt.Sprintf("%c", gd.Tok)).
			Msg("Unexpected GenDecl token: %v")
	}
}

func (sc *stmtCompiler) compileConstDecl(decl *ast.GenDecl) {
	for _, spec := range decl.Specs {
		spec := spec.(*ast.ValueSpec)
		if spec.Values == nil {
			// Declaration without assignment, the parser should have caught that.
			logger.Panic().Msg("missing constant value")
		}

		lhs := make([]ast.Expr, len(spec.Names))
		for i, n := range spec.Names {
			lhs[i] = n
		}

		sc.doConstAssign(lhs, spec.Values, decl.Tok, spec.Type)
	}
}

func (sc *stmtCompiler) compileImportDecl(decl *ast.GenDecl) {
	if !sc.Block.Global {
		logger.Panic().Msg("import at statement level")
	}

	for _, spec := range decl.Specs {
		spec := spec.(*ast.ImportSpec)
		path, _ := strconv.Unquote(spec.Path.Value)
		id := path
		if spec.Name != nil {
			id = spec.Name.Name
		}

		pkg, err := srcImporter(defaultImporter, path)
		if err != nil {
			sc.errorAt(spec.Pos(), "could not import package [%s]: %v", path, err)
			continue
		}

		if spec.Name != nil {
			sc.definePkg(spec.Name, id, path)
		} else {
			id = pkg.Name()
			sc.definePkg(spec.Path, id, path)
		}
	}
}

func (sc *stmtCompiler) compileTypeDecl(decl *ast.GenDecl) bool {
	name := decl.Specs[0].(*ast.TypeSpec).Name.Name
	_, level, dup := sc.Block.Lookup(name)
	if dup != nil && level == 0 {
		sc.error("type %s redeclared in this block, previous declaration at %s", name, sc.FSet.Position(dup.Pos()))
	}

	ok := true
	for _, spec := range decl.Specs {
		spec := spec.(*ast.TypeSpec)

		// Create incomplete type
		nt := sc.Block.DefineType(spec.Name.Name, spec.Name.Pos(), nil)
		if nt != nil {
			nt.(*types.NamedType).Incomplete = true
		}

		// Compile type
		tc := &typeCompiler{
			Compiler:  sc.BlockCompiler.Compiler,
			block:     sc.Block,
			lateCheck: noLateCheck,
		}
		t := tc.compileType(spec.Type, false)
		if t == nil {
			// Create a placeholder type
			ok = false
		}

		// Fill incomplete type
		nt.(*types.NamedType).Complete(t)

		// Perform late type checking with complete type
		if !tc.lateCheck() {
			ok = false
			if nt != nil {
				nt.(*types.NamedType).Def = nil // Make the type a placeholder
			}
		}
	}

	if !ok {
		sc.Block.Undefine(name)
	}

	return ok
}

// Statement generation helpers ////////////////////////////////////////////////

func (sc *stmtCompiler) compileVarDecl(decl *ast.GenDecl) error {
	for _, spec := range decl.Specs {
		spec := spec.(*ast.ValueSpec)
		if spec.Values == nil { // Declaration without assignment
			if spec.Type == nil {
				// The parser should have caught that
				logger.Panic().Msg("Type and Values nil in variable declaration")
			}

			t := sc.compileType(sc.Block, spec.Type)
			if t != nil {
				// If type compilation succeeded, define placeholders
				for _, n := range spec.Names {
					sc.defineVar(n, t)
				}
			}
		} else { // Declaration with assignment
			lhs := make([]ast.Expr, len(spec.Names))
			for i, n := range spec.Names {
				lhs[i] = n
			}

			sc.doAssign(lhs, spec.Values, decl.Tok, spec.Type)
		}
	}

	return nil
}

func (sc *stmtCompiler) defineVar(ident *ast.Ident, t vm.Type) *context.Variable {
	v, prev := sc.Block.DefineVar(ident.Name, ident.Pos(), t)
	if prev != nil {
		if prev.Pos().IsValid() {
			sc.errorAt(ident.Pos(), "variable %s redeclared in this block\n\tprevious declaration at %s", ident.Name, sc.FSet.Position(prev.Pos()))
		} else {
			sc.errorAt(ident.Pos(), "variable %s redeclared in this block", ident.Name)
		}
		return nil
	}

	// Initialize the variable
	index := v.Index
	if v.Index >= 0 {
		sc.Push(func(v *vm.Thread) {
			v.Frame.Vars[index] = t.Zero()
		})
	}
	return v
}

////////////////////////////////////////////////////////////////////////////////

func (sc *stmtCompiler) definePkg(ident ast.Node, id, path string) *context.PkgIdent {
	v, prev := sc.Block.DefinePackage(id, path, ident.Pos())
	if prev != nil {
		sc.errorAt(ident.Pos(), "%s redeclared as imported package name\n\tprevious declaration at %s", id, sc.FSet.Position(prev.Pos()))
		return nil
	}
	return v
}

func (sc *stmtCompiler) doConstAssign(lhs []ast.Expr, rhs []ast.Expr, tok token.Token, declTypeExpr ast.Expr) {
	nerr := sc.NumError()

	//////////////////////////////////////
	// Compile right side first so we have the types when compiling the left side and
	// so we don't see definitions made on the left side.

	rs := make([]*Expr, len(rhs))
	for i, re := range rhs {
		if !sc.checkConstExpr(re) {
			sc.errorAt(sc.pos, "const initializer %s is not a constant", "TODO:Expr.String()")
			continue
		}
		rs[i] = sc.CompileExpr(sc.Block, false, re)
	}

	//TODO: Check if we really need a assignCompiler
	ac, _ := sc.checkAssign(sc.pos, rs, "const declaration", "value")

	// If this is a definition and the LHS is larger than the RHS, we won't be able to produce
	// the usual error message because we can't begin to infer the types of the LHS.
	if tok == token.CONST && len(lhs) > len(ac.rmt.Elems) {
		sc.error("not enough values for definition")
	}

	//////////////////////////////////////
	// Compile left side

	// Compile left type if there is one
	var declType vm.Type
	if declTypeExpr != nil {
		declType = sc.compileType(sc.Block, declTypeExpr)
	}

	// If there have been errors, our arrays are full of nil's so get out of here now.
	if nerr != sc.NumError() {
		return
	}

	for i, le := range lhs {
		// Get the declaration  identifier and its type
		var lt vm.Type
		ident := le.(*ast.Ident)

		// Compute the identifier's type from the RHS type.
		switch {
		case declTypeExpr != nil:
			// We have a declaration type, use it.
			// If declType is nil, we gave an error when we compiled it.
			lt = declType

		case i >= len(ac.rmt.Elems):
			lt = nil // We already gave the "not enough" error above.
		case ac.rmt.Elems[i] == nil:
			lt = nil // We already gave the error when we compiled the RHS.

		default:
			lt = ac.rmt.Elems[i]
		}

		//////////////////////////////////////
		// define the identifier

		var v vm.Value
		switch lt.Lit().(type) {
		case *types.IdealFloatType:
			v = &values.IdealFloatV{V: rs[i].asIdealFloat()()}
		case *types.IdealIntType:
			v = &values.IdealIntV{V: rs[i].asIdealInt()()}
		case *types.BoolType:
			b := values.BoolV(rs[i].asIdealBool()())
			v = &b

		default:
			sc.errorAt(sc.pos, "const initializer is not a constant")
			v = nil
		}
		_, prev := sc.Block.DefineConst(ident.Name, sc.pos, lt, v)
		if prev != nil {
			pos := prev.Pos()
			if pos.IsValid() {
				sc.errorAt(sc.pos, "identifier %s redeclared in this block\n\tprevious declaration at %s", ident.Name, sc.FSet.Position(pos))
			} else {
				sc.errorAt(sc.pos, "identifier %s redeclared in this block", ident.Name)
			}
		}
	}
}

////////////////////////////////////////////////////////////////////////////////

func (sc *stmtCompiler) checkConstExpr(x ast.Expr) bool {
	switch x := x.(type) {
	case *ast.BasicLit:
		return true

	case *ast.BinaryExpr:
		// Check that both operands are constants
		return sc.checkConstExpr(x.X) && sc.checkConstExpr(x.Y)
	case *ast.ParenExpr:
		// Check that the parenthesized expression is constant
		return sc.checkConstExpr(x.X)
	case *ast.UnaryExpr:
		if x.Op == token.ARROW {
			return false
		}
		return sc.checkConstExpr(x.X)

	case *ast.Ident:
		_, _, v := sc.Block.Lookup(x.Name)
		c, ok := v.(*context.Constant)
		if !ok {
			return false
		}

		switch c.Type.(type) {
		case *types.FuncType:
			return false
		default:
			return true
		}

	default:
		// The rest of the expression types are incompatible with constants
		return false
	}
}

func (sc *stmtCompiler) doAssign(lhs []ast.Expr, rhs []ast.Expr, tok token.Token, declTypeExpr ast.Expr) {
	nerr := sc.NumError()

	//////////////////////////////////////
	// Compile right side first so we have the types when compiling the left side and
	// so we don't see definitions made on the left side.

	rs := make([]*Expr, len(rhs))
	for i, re := range rhs {
		rs[i] = sc.CompileExpr(sc.Block, false, re)
	}

	// Check if compileExpr failed on any expr
	for _, r := range rs {
		if r == nil {
			return
		}
	}

	errOp := "assignment"
	if tok == token.DEFINE || tok == token.VAR {
		errOp = "declaration"
	}

	ac, ok := sc.checkAssign(sc.pos, rs, errOp, "value")
	ac.allowMapForms(len(lhs))

	// If this is a definition and the LHS is larger than the RHS, we won't be able to produce
	// the usual error message because we can't begin to infer the types of the LHS.
	if (tok == token.DEFINE || tok == token.VAR) && len(lhs) > len(ac.rmt.Elems) {
		sc.error("not enough values for definition")
	}

	//////////////////////////////////////
	// Compile left side

	// Compile left type if there is one
	var declType vm.Type
	if declTypeExpr != nil {
		declType = sc.compileType(sc.Block, declTypeExpr)
	}

	ls := make([]*Expr, len(lhs))
	nDefs := 0
	for i, le := range lhs {
		// If this is a definition, get the identifier and its type
		var ident *ast.Ident
		var lt vm.Type
		switch tok {
		case token.DEFINE:
			// Check that it's an identifier
			ident, ok = le.(*ast.Ident)
			if !ok {
				sc.errorAt(le.Pos(), "left side of := must be a name")
				nDefs++ // Suppress new definitions errors
				continue
			}

			// Is this simply an assignment?
			if _, ok := sc.Block.Defs[ident.Name]; ok {
				ident = nil
				break
			}
			nDefs++

		case token.VAR:
			ident = le.(*ast.Ident)
		}

		// If it's a definition, get or infer its type.
		if ident != nil {
			// Compute the identifier's type from the RHS type.
			// We use the computed MultiType so we don't have to worry about unpacking.
			switch {
			case declTypeExpr != nil:
				// We have a declaration type, use it.
				// If declType is nil, we gave an error when we compiled it.
				lt = declType

			case i >= len(ac.rmt.Elems):
				lt = nil // We already gave the "not enough" error above.
			case ac.rmt.Elems[i] == nil:
				lt = nil // We already gave the error when we compiled the RHS.

			case ac.rmt.Elems[i].IsIdeal():
				// If the type is absent and the corresponding expression is a constant expression of
				// ideal integer or ideal float type, the type of the declared variable is int or float respectively.
				switch {
				case ac.rmt.Elems[i].IsInteger():
					lt = builtins.IntType
				case ac.rmt.Elems[i].IsFloat():
					lt = builtins.Float64Type
				default:
					logger.Panic().
						Str("type", fmt.Sprintf("%v", rs[i].ExprType)).
						Msg("unexpected ideal type")
				}

			default:
				lt = ac.rmt.Elems[i]
			}

			// define the identifier
			if sc.defineVar(ident, lt) == nil {
				continue
			}
		}

		// Compile LHS
		ls[i] = sc.CompileExpr(sc.Block, false, le)
		if ls[i] == nil {
			continue
		}

		if ls[i].evalMapValue != nil {
			// Map indexes are not generally addressable, but they are assignable.
			sub := ls[i]
			ls[i] = ls[i].newExpr(sub.ExprType, sub.desc)
			ls[i].evalMapValue = sub.evalMapValue
			mvf := sub.evalMapValue
			et := sub.ExprType
			ls[i].evalAddr = func(t *vm.Thread) vm.Value {
				m, k := mvf(t)
				e := m.Elem(t, k)
				if e == nil {
					e = et.Zero()
					m.SetElem(t, k, e)
				}
				return e
			}
		} else if ls[i].evalAddr == nil {
			ls[i].error("cannot assign to %s", ls[i].desc)
			continue
		}
	}

	// A short variable declaration may redeclare variables provided they were originally
	// declared in the same block with the same type, and at least one of the variables is new.
	if tok == token.DEFINE && nDefs == 0 {
		sc.error("at least one new variable must be declared")
		return
	}

	// If there have been errors, our arrays are full of nil's so get out of here now.
	if nerr != sc.NumError() {
		return
	}

	// Check for 'a[x] = r, ok'
	if len(ls) == 1 && len(rs) == 2 && ls[0].evalMapValue != nil {
		sc.error("a[x] = r, ok form not implemented")
		return
	}

	//////////////////////////////////////
	// Create assigner

	var lt vm.Type
	n := len(lhs)
	if n == 1 {
		lt = ls[0].ExprType
	} else {
		lts := make([]vm.Type, len(ls))
		for i, l := range ls {
			if l != nil {
				lts[i] = l.ExprType
			}
		}
		lt = types.NewMultiType(lts)
	}

	bc := sc.enterChild()
	defer bc.exit()
	assign := ac.compile(bc.Block, lt)
	if assign == nil {
		return
	}

	//////////////////////////////////////
	// Compile

	switch {
	case n == 1:
		// Don't need temporaries and can avoid []Value.
		lf := ls[0].evalAddr
		sc.Push(func(t *vm.Thread) {
			assign(lf(t), t)
		})

	case tok == token.VAR || (tok == token.DEFINE && nDefs == n):
		// Don't need temporaries
		lfs := make([]func(*vm.Thread) vm.Value, n)
		for i, l := range ls {
			lfs[i] = l.evalAddr
		}

		sc.Push(func(t *vm.Thread) {
			dest := make([]vm.Value, n)
			for i, lf := range lfs {
				dest[i] = lf(t)
			}
			assign(values.MultiV(dest), t)
		})

	default:
		// Need temporaries
		lmt := lt.(*types.MultiType)
		lfs := make([]func(*vm.Thread) vm.Value, n)
		for i, l := range ls {
			lfs[i] = l.evalAddr
		}

		sc.Push(func(t *vm.Thread) {
			temp := lmt.Zero().(values.MultiV)
			assign(temp, t)
			// Copy to destination
			for i := 0; i < n; i++ {
				// TODO: Need to evaluate LHS before RHS
				lfs[i](t).Assign(t, temp[i])
			}
		})
	}
}

var assignOpToOp = map[token.Token]token.Token{
	token.ADD_ASSIGN: token.ADD,
	token.SUB_ASSIGN: token.SUB,
	token.MUL_ASSIGN: token.MUL,
	token.QUO_ASSIGN: token.QUO,
	token.REM_ASSIGN: token.REM,

	token.AND_ASSIGN:     token.AND,
	token.OR_ASSIGN:      token.OR,
	token.XOR_ASSIGN:     token.XOR,
	token.SHL_ASSIGN:     token.SHL,
	token.SHR_ASSIGN:     token.SHR,
	token.AND_NOT_ASSIGN: token.AND_NOT,
}

func (sc *stmtCompiler) doAssignOp(stmt *ast.AssignStmt) {
	if len(stmt.Lhs) != 1 || len(stmt.Rhs) != 1 {
		sc.error("tuple assignment cannot be combined with an arithmetic operation")
		return
	}

	// Create temporary block for extractEffect
	bc := sc.enterChild()
	defer bc.exit()

	l := sc.CompileExpr(bc.Block, false, stmt.Lhs[0])
	r := sc.CompileExpr(bc.Block, false, stmt.Rhs[0])
	if l == nil || r == nil {
		return
	}

	if l.evalAddr == nil {
		l.error("cannot assign to %s", l.desc)
		return
	}

	effect, l := l.extractEffect(bc.Block, "operator-assignment")

	binop := r.compileBinaryExpr(assignOpToOp[stmt.Tok], l, r)
	if binop == nil {
		return
	}

	assign := sc.compileAssign(stmt.Pos(), bc.Block, l.ExprType, []*Expr{binop}, "assignment", "value")
	if assign == nil {
		logger.Panic().Msgf("compileAssign type check failed")
	}

	lf := l.evalAddr
	sc.Push(func(t *vm.Thread) {
		effect(t)
		assign(lf(t), t)
	})
}

func (sc *stmtCompiler) findLexicalLabel(name *ast.Ident, pred func(*Label) bool, errOp, errCtx string) *Label {
	bc := sc.BlockCompiler
	for ; bc != nil; bc = bc.Parent {
		if bc.Label == nil {
			continue
		}
		l := bc.Label
		if name == nil && pred(l) {
			return l
		}
		if name != nil && l.name == name.Name {
			if !pred(l) {
				sc.error("cannot %s to %s %s", errOp, l.desc, l.name)
				return nil
			}
			return l
		}
	}
	if name == nil {
		sc.error("%s outside %s", errOp, errCtx)
	} else {
		sc.error("%s label %s not defined", errOp, name.Name)
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////

// default importer of packages
var defaultImporter = importer.Default()

// srcImporter implements the ast.Importer signature.
func srcImporter(typesImporter gotypes.Importer, path string) (pkg *gotypes.Package, err error) {
	return typesImporter.Import(path)
}
