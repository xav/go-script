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

package tests

import (
	"testing"
)

func TestCompile_DeclImplicitVar_Int(t *testing.T) {
	code := `
package main

var a = 2
`
	err := compile(t, code)
	Ok(t, err)
}

func TestCompile_DeclImplicitVar_Float(t *testing.T) {
	code := `
package main

var a = 2
`
	err := compile(t, code)
	Ok(t, err)
}

func TestCompile_DeclImplicitVar_Bool(t *testing.T) {
	code := `
package main

var a = true
`
	err := compile(t, code)
	Ok(t, err)
}

func TestCompile_DeclImplicitVar_LitString(t *testing.T) {
	code := `
package main

var a = "abcd"
`
	err := compile(t, code)
	Ok(t, err)
}

func TestCompile_DeclImplicitVar_LitIntSlice(t *testing.T) {
	code := `
package main

var a = []int{2, 3, 5, 7, 11, 13}
`
	err := compile(t, code)
	Ok(t, err)
}

func TestCompile_DeclImplicitVar_LitIntArray(t *testing.T) {
	code := `
package main

var a = [6]int{2, 3, 5, 7, 11, 13}
`
	err := compile(t, code)
	Ok(t, err)
}

func TestCompile_DeclImplicitVar_LitStringIntMap(t *testing.T) {
	code := `
package main

var a = map[string]int{
	"foo": 1,
	"bar": 2,
}
`
	err := compile(t, code)
	Ok(t, err)
}

func TestCompile_DeclImplicitVar_LitStruct(t *testing.T) {
	code := `
package main

type Vertex struct {
	X, Y int
}

var a = Vertex{1, 2}
`
	err := compile(t, code)
	Ok(t, err)
}
