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

package vm

type Frame struct {
	Outer *Frame
	Vars  []Value
}

// Get the Value at the specified slot in the specified parent level frame.
func (f *Frame) Get(level int, index int) Value {
	for ; level > 0; level-- {
		f = f.Outer
	}
	return f.Vars[index]
}

// NewChild creates a new child frame with the specified number of variable slots.
func (f *Frame) NewChild(numVars int) *Frame {
	return &Frame{
		Outer: f,
		Vars:  make([]Value, numVars),
	}
}
