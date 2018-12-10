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
	"go/token"

	"github.com/xav/go-script/builtins"
	"github.com/xav/go-script/context"
	"github.com/xav/go-script/types"
	"github.com/xav/go-script/values"
	"github.com/xav/go-script/vm"
)

// assignCompiler compiles assignment operations.
// Anything other than short declarations should use the compileAssign wrapper.
//
// There are three valid types of assignment:
// 1) T = T
//    Assigning a single expression with single-valued type to a
//    single-valued type.
// 2) MT = T, T, ...
//    Assigning multiple expressions with single-valued types to a
//    multi-valued type.
// 3) MT = MT
//    Assigning a single expression with multi-valued type to a
//    multi-valued type.
type assignCompiler struct {
	*Compiler
	pos         token.Pos
	rs          []*Expr          // The RHS expressions.  This may include nil's for expressions that failed to compile.
	rmt         *types.MultiType // The (possibly unary) MultiType of the RHS.
	isUnpack    bool             // Whether this is an unpack assignment (case 3).
	allowMap    bool             // Whether map special assignment forms are allowed.
	isMapUnpack bool             // Whether this is a "r, ok = a[x]" assignment.
	errOp       string           // The operation name to use in error messages, such as "assignment" or "function call".
	errPosName  string           // The name to use for positions in error messages, such as "argument".
}

// compile type checks and compiles an assignment operation, returning a function that expects
// an l-value and the frame in which to evaluate the RHS expressions.
// The l-value must have exactly the type given by lt.
// Returns nil if type checking fails.
func (ac *assignCompiler) compile(b *context.Block, lt vm.Type) func(vm.Value, *vm.Thread) {
	lmt, isMT := lt.(*types.MultiType)
	rmt, isUnpack := ac.rmt, ac.isUnpack

	// Create unary MultiType for single LHS
	if !isMT {
		lmt = types.NewMultiType([]vm.Type{lt})
	}

	// Check that the assignment count matches
	lc := len(lmt.Elems)
	rc := len(rmt.Elems)
	if lc != rc {
		msg := "not enough"
		pos := ac.pos
		if rc > lc {
			msg = "too many"
			if lc > 0 {
				pos = ac.rs[lc-1].pos
			}
		}
		ac.errorAt(pos, "%s %ss for %s\n\t%s\n\t%s", msg, ac.errPosName, ac.errOp, lt, rmt)
		return nil
	}

	bad := false
	var effect func(*vm.Thread)

	// If this is an unpack, create a temporary to store the multi-value and replace
	// the RHS with expressions to pull out values from the temporary.
	// Technically, this is only necessary when we need to perform assignment conversions.
	if isUnpack {
		temp := b.DefineTemp(ac.rmt)
		tempIdx := temp.Index

		if ac.isMapUnpack {
			rf := ac.rs[0].evalMapValue
			vt := ac.rmt.Elems[0]
			effect = func(t *vm.Thread) {
				m, k := rf(t)
				v := m.Elem(t, k)
				found := values.BoolV(true)
				if v == nil {
					found = values.BoolV(false)
					v = vt.Zero()
				}
				t.Frame.Vars[tempIdx] = values.MultiV([]vm.Value{v, &found})
			}
		} else {
			rf := ac.rs[0].asMulti()
			effect = func(t *vm.Thread) {
				t.Frame.Vars[tempIdx] = values.MultiV(rf(t))
			}
		}
	}

	//////////////////////////////////////
	// At this point, len(ac.rs) == len(ac.rmt) and we've reduced any unpacking to multi-assignment.

	// Values of any type may always be assigned to variables of compatible static type.
	for i, lt := range lmt.Elems {
		rt := rmt.Elems[i]

		// When an ideal is (used in an expression) assigned to a variable or typed constant,
		// the destination must be able to represent the assigned value.
		if rt.IsIdeal() {
			ac.rs[i] = ac.rs[i].resolveIdeal(lmt.Elems[i])
			if ac.rs[i] == nil {
				bad = true
				continue
			}
			rt = ac.rs[i].ExprType
		}

		// A pointer p to an array can be assigned to a slice variable v with compatible
		// element type if the type of p or v is unnamed.
		if rpt, ok := rt.Lit().(*types.PtrType); ok {
			if at, ok := rpt.Elem.Lit().(*types.ArrayType); ok {
				if lst, ok := lt.Lit().(*types.SliceType); ok {
					if lst.Elem.Compat(at.Elem, false) && (rt.Lit() == vm.Type(rt) || lt.Lit() == vm.Type(lt)) {
						rf := ac.rs[i].asPtr()
						ac.rs[i] = ac.rs[i].newExpr(lt, ac.rs[i].desc)
						len := at.Len
						ac.rs[i].eval = func(t *vm.Thread) values.Slice {
							return values.Slice{
								Base: rf(t).(values.ArrayValue),
								Len:  len,
								Cap:  len,
							}
						}
					}
				}
			}
		}

		if !lt.Compat(rt, false) {
			if len(ac.rs) == 1 {
				ac.rs[0].error("illegal operand types for %s\n\t%v\n\t%v", ac.errOp, lt, rt)
			} else {
				ac.rs[i].error("illegal operand types in %s %d of %s\n\t%v\n\t%v", ac.errPosName, i+1, ac.errOp, lt, rt)
			}
			bad = true
		}
	}

	if bad {
		return nil
	}

	//////////////////////////////////////
	// Compile

	if !isMT { // Case 1 (T = T)
		return genAssign(lt, ac.rs[0])
	}

	// Case 2 (MT = T, T, ...) or 3 (MT = MT)
	as := make([]func(lv vm.Value, t *vm.Thread), len(ac.rs))
	for i, r := range ac.rs {
		as[i] = genAssign(lmt.Elems[i], r)
	}

	return func(lv vm.Value, t *vm.Thread) {
		if effect != nil {
			effect(t)
		}
		lmv := lv.(values.MultiV)
		for i, a := range as {
			a(lmv[i], t)
		}
	}
}

func (ac *assignCompiler) allowMapForms(nls int) {
	ac.allowMap = true

	// Update unpacking info if this is 'r, ok = a[x]'
	if nls == 2 && len(ac.rs) == 1 && ac.rs[0] != nil && ac.rs[0].evalMapValue != nil {
		ac.isUnpack = true
		ac.isMapUnpack = true
		ac.rmt = types.NewMultiType([]vm.Type{
			ac.rs[0].ExprType,
			builtins.BoolType,
		})
	}
}
