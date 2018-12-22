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

package values

import (
	"fmt"

	"github.com/xav/go-script/vm"
)

type StringValue interface {
	vm.Value
	Get(*vm.Thread) string
	Set(*vm.Thread, string)
}

// string //////////////////////////////////////////////////////////////////////

type StringV string

func (v *StringV) String() string {
	return fmt.Sprint(*v)
}

func (v *StringV) Assign(t *vm.Thread, o vm.Value) {
	*v = StringV(o.(StringValue).Get(t))
}

func (v *StringV) Get(*vm.Thread) string {
	return string(*v)
}

func (v *StringV) Set(t *vm.Thread, x string) {
	*v = StringV(x)
}
