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
type Block struct {
	Outer *Block         // The block enclosing this one
	Inner *Block         // The nested block currently being compiled, or nil.
	Scope *Scope         // The Scope containing this block.
	Defs  map[string]Def // The Variables, Constants, and Types defined in this block.

	// The index of the first variable defined in this block.
	// This must be greater than the index of any variable defined in any parent of this block
	// within the same Scope at the time this block is entered.
	Offset int
	// The number of Variables defined in this block.
	NumVars int

	// If Global, do not allocate new vars and consts in the frame;
	// assume that the refs will be compiled in using defs[name].Init.
	Global bool
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
		Offset: b.Offset + b.NumVars,
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
	sub.Offset = 0
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
		logger.Panic().Msg("failed to exit child block before defining variable")
	}

	return &types.Variable{
		VarPos: token.NoPos,
		Index:  b.defineSlot(true),
		Type:   t,
	}
}

// DefineConst creates a new constant definition and allocate it in the current scope,
// or returns the existing definition if the symbol was already defined.
func (b *Block) DefineConst(name string, pos token.Pos, t vm.Type) (*types.Constant, Def) {
	if prev, ok := b.Defs[name]; ok {
		if _, ok := prev.(*types.Constant); ok {
			logger.Panic().
				Str("symbol", name).
				Msgf("symbol redeclaration with different primitives types (%+v -> Constant)", prev)
		}
		return nil, prev
	}

	if b.Inner != nil && b.Inner.Scope == b.Scope {
		logger.Panic().Msg("failed to exit child block before defining constant")
	}

	c := &types.Constant{
		ConstPos: pos,
		Index:    b.defineSlot(false),
		Type:     t,
	}
	b.Defs[name] = c
	return c, nil
}

// defineSlot allocates a slot in the block's scope.
func (b *Block) defineSlot(temp bool) int {
	if b.Inner != nil && b.Inner.Scope == b.Scope {
		logger.Panic().Msg("failed to exit child block before defining symbol")
	}

	index := -1
	if !b.Global || temp {
		index = b.Offset + b.NumVars
		b.NumVars++
		if index >= b.Scope.MaxVars {
			b.Scope.MaxVars = index + 1
		}
	}

	return index
}

// DefineType creates a user defined type.
func (b *Block) DefineType(name string, pos token.Pos, t vm.Type) vm.Type {
	if prev, ok := b.Defs[name]; ok {
		if _, ok := prev.(*types.NamedType); ok {
			logger.Panic().
				Str("symbol", name).
				Msgf("symbol redeclaration with different primitives types (%+v -> NamedType)", prev)
		}
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

//TODO(xav): channels

// Undefine removes a symbole definition from this block.
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
