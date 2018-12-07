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

package types

import (
	"reflect"

	"github.com/xav/go-script/vm"
)

// Type array maps are used to memoize composite types.
type typeArrayMap map[uintptr]*typeArrayMapEntry

type typeArrayMapEntry struct {
	key  []vm.Type
	v    interface{}
	next *typeArrayMapEntry
}

func newTypeArrayMap() typeArrayMap {
	return make(map[uintptr]*typeArrayMapEntry)
}

func hashTypeArray(key []vm.Type) uintptr {
	hash := uintptr(0)
	for _, t := range key {
		hash = hash * 33
		if t == nil {
			continue
		}
		addr := reflect.ValueOf(t).Pointer()
		hash ^= addr
	}
	return hash
}

func (m typeArrayMap) Get(key []vm.Type) interface{} {
	ent, ok := m[hashTypeArray(key)]
	if !ok {
		return nil
	}

nextEnt:
	for ; ent != nil; ent = ent.next {
		if len(key) != len(ent.key) {
			continue
		}
		for i := 0; i < len(key); i++ {
			if key[i] != ent.key[i] {
				continue nextEnt
			}
		}
		// Found it
		return ent.v
	}

	return nil
}

func (m typeArrayMap) Put(key []vm.Type, v interface{}) interface{} {
	hash := hashTypeArray(key)
	ent := m[hash]

	new := &typeArrayMapEntry{key, v, ent}
	m[hash] = new
	return v
}
