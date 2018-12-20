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

	"github.com/xav/go-script/context"
)

type FlowBuf struct {
	cb *CodeBuf

	ents   map[uint]*FlowEnt        // map from PC's to flow entries. Any PC missing from this map is assumed to reach only PC+1.
	gotos  map[token.Pos]*FlowBlock // map from goto positions to information on the block at the point of the goto.
	labels map[string]*FlowBlock    // map from label name to information on the block at the point of the label.

	// labels are tracked by name, since multiple labels at the same PC can have different blocks.
}

type FlowBlock struct {
	target  string         // If this is a goto, the target label.
	block   *context.Block // The inner-most block containing definitions.
	numVars []int          // The numVars from each block leading to the root of the scope, starting at block.
}

type FlowEnt struct {
	cond    bool    // Whether this flow entry is conditional. If true, flow can continue to the next PC.
	term    bool    // True if this will terminate flow (e.g., a return statement). cond must be false and jumps must be nil if this is true.
	jumps   []*uint // PC's that can be reached from this flow entry.
	visited bool    // Whether this flow entry has been visited by reachesEnd.
}

// NewFlowBuf creates a new FlowBuf using the specified CodeBuf
func NewFlowBuf(cb *CodeBuf) *FlowBuf {
	return &FlowBuf{
		cb:     cb,
		ents:   make(map[uint]*FlowEnt),
		gotos:  make(map[token.Pos]*FlowBlock),
		labels: make(map[string]*FlowBlock),
	}
}

// reachesEnd returns true if the end of f's code buffer can be reached from the given program counter.
// Error reporting is the caller's responsibility.
func (f *FlowBuf) reachesEnd(pc uint) bool {
	endPC := f.cb.NextPC()
	if pc > endPC {
		logger.Panic().Msgf("Reached bad PC %d past end PC %d", pc, endPC)
	}

	for ; pc < endPC; pc++ {
		ent, ok := f.ents[pc]
		if !ok {
			continue
		}

		if ent.visited {
			return false
		}
		ent.visited = true

		if ent.term {
			return false
		}

		// If anything can reach the end, we can reach the end from pc.
		for _, j := range ent.jumps {
			if f.reachesEnd(*j) {
				return true
			}
		}

		// If the jump was conditional, we can reach the next
		// PC, so try reaching the end from it.
		if ent.cond {
			continue
		}

		return false
	}

	return true
}

// gotosObeyScopes returns true if no goto statement causes any variables to come
// into scope that were not in scope at the point of the goto.
// Reports any errors using the specified compiler.
func (f *FlowBuf) gotosObeyScopes(pc *Compiler) {
	for pos, src := range f.gotos {
		tgt := f.labels[src.target]

		// The target block must be a parent of this block
		numVars := src.numVars
		b := src.block
		for len(numVars) > 0 && b != tgt.block {
			b = b.Outer
			numVars = numVars[1:]
		}
		if b != tgt.block {
			// We jumped into a deeper block
			pc.errorAt(pos, "goto causes variables to come into scope")
			return
		}

		// There must be no variables in the target block that did not exist at the jump
		tgtNumVars := tgt.numVars
		for i := range numVars {
			if tgtNumVars[i] > numVars[i] {
				pc.errorAt(pos, "goto causes variables to come into scope")
				return
			}
		}
	}
}
