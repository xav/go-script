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

// EvalFunc is used for user-defined functions.
type EvalFunc struct {
	Outer     *Frame
	FrameSize int
	Code      Code
}

func (f *EvalFunc) NewFrame() *Frame {
	return f.Outer.NewChild(f.FrameSize)
}

func (f *EvalFunc) Call(t *Thread) {
	f.Code.Exec(t)
}
