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

import (
	"github.com/xav/go-script/types"
)

// FuncCompiler captures information used throughout the compilation of a single function body.
type FuncCompiler struct {
	*Compiler                      // The Compiler for the package enclosing this function.
	*CodeBuf                       //
	FnType       *types.FuncType   //
	OutVarsNamed bool              // Whether the out variables are named. This affects what kinds of return statements are legal.
	Flow         *FlowBuf          //
	Labels       map[string]*Label //
}

// CheckLabels checks that labels were resolved and that all jumps obey scoping rules.
// Reports an error and set fc.err if any check fails.
func (fc *FuncCompiler) CheckLabels() {
	nerr := fc.NumError()
	for _, l := range fc.Labels {
		if !l.resolved.IsValid() {
			fc.errorAt(l.used, "label %s not defined", l.name)
		}
	}
	if nerr != fc.NumError() {
		// Don't check scopes if we have unresolved labels
		return
	}

	// Executing a "goto" statement must not cause any variables to come
	// into scope that were not already in scope at the point of the goto.
	fc.Flow.gotosObeyScopes(fc.Compiler)
}
