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

package script

import "github.com/xav/go-script/vm"

// Runnable represents code runnable by the vm
type Runnable interface {
	Type() vm.Type          // The type of the value Run returns, or nil if Run returns nil.
	Run() (vm.Value, error) // Run runs the code; if the code is a single expression with a value, it returns the value; otherwise it returns nil.
}
