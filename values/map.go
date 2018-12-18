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

type Map interface {
	Len(*vm.Thread) int64
	Elem(t *vm.Thread, key interface{}) vm.Value         // Retrieve an element from the map, returning nil if it does not exist.
	SetElem(t *vm.Thread, key interface{}, val vm.Value) // Set an entry in the map. If val is nil, delete the entry.
	Iter(func(key interface{}, val vm.Value) bool)       // TODO:  Perhaps there should be an iterator interface instead.
}

type MapValue interface {
	vm.Value
	Get(*vm.Thread) Map
	Set(*vm.Thread, Map)
}

// map /////////////////////////////////////////////////////////////////////////

type MapV struct {
	Target Map
}

func (v *MapV) String() string {
	if v.Target == nil {
		return "<nil>"
	}
	res := "map["
	i := 0
	v.Target.Iter(func(key interface{}, val vm.Value) bool {
		if i > 0 {
			res += ", "
		}
		i++
		res += fmt.Sprint(key) + ":" + val.String()
		return true
	})
	return res + "]"
}

func (v *MapV) Assign(t *vm.Thread, o vm.Value) {
	v.Target = o.(MapValue).Get(t)
}

func (v *MapV) Get(*vm.Thread) Map {
	return v.Target
}

func (v *MapV) Set(t *vm.Thread, x Map) {
	v.Target = x
}
