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

import "github.com/xav/go-script/context"

var universe = newUniverse()

// Universe is a special scope for the universe for all the worlds.
type Universe struct {
	*context.Scope
	Pkgs map[string]*context.Scope // a lookup-table for easy retrieval of packages by their "path"
}

// The universal scope
func newUniverse() *Universe {
	scope := &Universe{
		Scope: new(context.Scope),
		Pkgs:  make(map[string]*context.Scope),
	}

	scope.Block = &context.Block{
		Scope:  scope.Scope,
		Offset: 0,
		Global: true,
		Defs:   make(map[string]context.Def),
	}

	return scope
}
