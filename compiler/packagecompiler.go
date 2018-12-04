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
	"fmt"
	"go/scanner"
	"go/token"
)

// PackageCompiler captures information used throughout a package compilation.
type PackageCompiler struct {
	FSet         *token.FileSet
	Errors       *scanner.ErrorList
	NumErrors    int
	SilentErrors int
}

// Reports a compilation error att the specified position
func (pc *PackageCompiler) errorAt(pos token.Pos, format string, args ...interface{}) {
	pc.Errors.Add(pc.FSet.Position(pos), fmt.Sprintf(format, args...))
	pc.NumErrors++
}

// NumError returns the total number of errors detected yet
func (pc *PackageCompiler) NumError() int {
	return pc.NumErrors + pc.SilentErrors
}
