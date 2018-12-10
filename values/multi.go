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

import "github.com/xav/go-script/vm"

type MultiV []vm.Value

func (v MultiV) String() string {
	res := "("
	for i, v := range v {
		if i > 0 {
			res += ", "
		}
		res += v.String()
	}
	return res + ")"
}

func (v MultiV) Assign(t *vm.Thread, o vm.Value) {
	omv := o.(MultiV)
	for i := range v {
		v[i].Assign(t, omv[i])
	}
}
