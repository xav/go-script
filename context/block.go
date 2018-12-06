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

package context

import (
	"go/token"

	"github.com/xav/go-script/types"
	"github.com/xav/go-script/vm"
)

// Block represents a definition context in which a name may not be defined more than once.
//
// The variables offset must be greater than the index of any variable  defined in any
// parent of this block within the same Scope at the time this block is entered.
//
// Global blocks assume that the refs will be compiled in using defs[name].Init.
type Block struct {
	Outer  *Block         // The block enclosing this one
	Inner  *Block         // The nested block currently being compiled, or nil.
	Scope  *Scope         // The Scope containing this block.
	Defs   map[string]Def // The Variables, Constants, and Types defined in this block.
	Global bool           // If Global, do not allocate new vars and consts in the frame;

	offset  int // The index of the first variable defined in this block.
	numVars int // The number of variables defined in this block.
}

// EnterChild creates a new block with this one as parent, using the same scope.
func (b *Block) EnterChild() *Block {
	if b.Inner != nil && b.Inner.Scope == b.Scope {
		logger.Panic().Msg("failed to exit child block before entering another child")
	}

	sub := &Block{
		Outer:  b,
		Scope:  b.Scope,
		Defs:   make(map[string]Def),
		offset: b.offset + b.numVars,
	}
	b.Inner = sub
	return sub
}

// EnterChildScope creates a new block with this one as parent, using a new scope.
func (b *Block) EnterChildScope() *Scope {
	if b.Inner != nil && b.Inner.Scope == b.Scope {
		logger.Panic().Msg("failed to Exit child block before entering a child scope")
	}

	sub := b.EnterChild()
	sub.offset = 0
	sub.Scope = &Scope{
		Block:   sub,
		MaxVars: 0,
	}
	return sub.Scope
}

// Exit closes this block in it's parent.
func (b *Block) Exit() {
	if b.Outer == nil {
		logger.Panic().Msg("cannot exit top-level block")
	}

	if b.Outer.Scope == b.Scope {
		if b.Outer.Inner != b {
			logger.Panic().Msg("already exited block")
		}
		if b.Inner != nil && b.Inner.Scope == b.Scope {
			logger.Panic().Msg("exit of parent block without exit of child block")
		}
	}

	b.Outer.Inner = nil
}

// DefineVar creates a new variable definition and allocate it in the current scope,
// or returns the existing definition if the symbol was already defined.
func (b *Block) DefineVar(name string, pos token.Pos, t vm.Type) (*types.Variable, Def) {
	if prev, ok := b.Defs[name]; ok {
		if _, ok := prev.(*types.Variable); ok {
			logger.Panic().
				Str("symbol", name).
				Msgf("symbol redeclaration with different primitives types (%+v -> Variable)", prev)
		}

		return nil, prev
	}

	if b.Inner != nil && b.Inner.Scope == b.Scope {
		logger.Panic().Msg("failed to exit child block before defining variable")
	}

	v := &types.Variable{
		VarPos: pos,
		Index:  b.defineSlot(false),
		Type:   t,
	}
	b.Defs[name] = v
	return v, nil
}

// DefineTemp allocates a anonymous slot in the current scope.
func (b *Block) DefineTemp(t vm.Type) *types.Variable {
	if b.Inner != nil && b.Inner.Scope == b.Scope {
		logger.Panic().Msg("failed to exit child block before defining temp")
	}

	return &types.Variable{
		VarPos: token.NoPos,
		Index:  b.defineSlot(true),
		Type:   t,
	}
}

// defineSlot allocates a var slot in the block's scope.
func (b *Block) defineSlot(temp bool) int {
	index := -1
	if !b.Global || temp {
		index = b.offset + b.numVars
		b.numVars++
		if index >= b.Scope.MaxVars {
			b.Scope.MaxVars = index + 1
		}
	}

	return index
}

// DefineConst creates a new constant definition.
func (b *Block) DefineConst(name string, pos token.Pos, t vm.Type) *types.Constant {
	if _, ok := b.Defs[name]; ok {
		logger.Error().
			Str("symbol", name).
			Msg("constant already declared in this block")
		return nil
	}

	c := &types.Constant{
		ConstPos: pos,
		Type:     t,
	}

	b.Defs[name] = c
	return c
}

// DefineType creates a user defined type.
func (b *Block) DefineType(name string, pos token.Pos, t vm.Type) vm.Type {
	if _, ok := b.Defs[name]; ok {
		logger.Error().
			Str("symbol", name).
			Msg("type already declared in this block")
		return nil
	}

	nt := &types.NamedType{
		NTPos:      pos,
		Name:       name,
		Def:        nil,
		Incomplete: true,
	}

	if t != nil {
		nt.Complete(t)
	}

	b.Defs[name] = nt
	return nt
}

//TODO(xav): DefineChan for channels

// Undefine removes a symbol definition from this block.
func (b *Block) Undefine(name string) {
	delete(b.Defs, name)
}

// Lookup search for the definition of the specified symbol by marching up
// the blocks tree until either the symbol is found, or the scope root is reached.
func (b *Block) Lookup(name string) (bl *Block, level int, def Def) {
	for b != nil {
		if d, ok := b.Defs[name]; ok {
			return b, level, d
		}

		if b.Outer != nil && b.Scope != b.Outer.Scope {
			level++
		}
		b = b.Outer
	}
	return nil, 0, nil
}
