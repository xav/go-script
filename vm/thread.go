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

import "runtime"

type Thread struct {
	abort chan error //
	PC    uint       // Program Counter
	Frame *Frame     // The execution frame of this function. This remains the same throughout a function invocation.
}

// Abort aborts the thread's current computation, causing the innermost Try to return err.
func (t *Thread) Abort(err error) {
	if t.abort == nil {
		panic("abort: " + err.Error())
	}
	t.abort <- err
	runtime.Goexit()
}

// Try executes a computation; if the computation Aborts, Try returns the error passed to abort.
func (t *Thread) Try(f func(t *Thread)) error {
	oc := t.abort
	c := make(chan error)
	t.abort = c

	go func() {
		f(t)
		c <- nil
	}()

	err := <-c
	t.abort = oc

	return err
}
