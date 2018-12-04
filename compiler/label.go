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

package compiler

import "go/token"

type Label struct {
	name       string
	desc       string
	gotoPC     *uint     // The PC goto statements should jump to, or nil if this label cannot be goto'd (such as an anonymous for loop label).
	breakPC    *uint     // The PC break statements should jump to, or nil if a break statement is invalid.
	continuePC *uint     // The PC continue statements should jump to, or nil if a continue statement is invalid.
	resolved   token.Pos // The position where this label was resolved.  If it has not been resolved yet, an invalid position.
	used       token.Pos // The position where this label was first jumped to.
}
