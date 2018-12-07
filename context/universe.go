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

// UniverseScope contains the global scope for the whole vm.
type UniverseScope struct {
	*Scope
	Pkgs map[string]*Scope // a lookup-table for easy retrieval of packages by their "path"
}

// The universal scope
func NewUniverse() *UniverseScope {
	scope := &UniverseScope{
		Scope: new(Scope),
		Pkgs:  make(map[string]*Scope),
	}

	scope.Block = &Block{
		Scope:  scope.Scope,
		Global: true,
		Defs:   make(map[string]Def),
	}

	return scope
}
