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

import "github.com/xav/go-script/vm"

// Scope is the compile-time analogue of a Frame, which captures some subtree of blocks.
type Scope struct {
	// The root block of this scope.
	*Block
	// The maximum number of variables required at any point in this Scope.
	// This determines the number of slots needed in the Frame created from this Scope at run-time.
	MaxVars int
}

// NewFrame creates a new child frame to the one specified, with the number of variable slots
// needed by the current scope.
func (s *Scope) NewFrame(outerFrame *vm.Frame) *vm.Frame {
	return outerFrame.NewChild(s.MaxVars)
}
