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

package testing

import (
	"testing"
)

func TestCompile_Builtins_Append(t *testing.T) {
	code := `
package main
func main() {
	a := []int{1,2,3}
	a = append(a, 4, 5, 6)
}`
	testCompilation(t, code)
}

func TestCompile_Builtins_Cap(t *testing.T) {
	code := `
package main
func main() {
	a := []int{1,2,3}
	b := cap(a)
}`
	testCompilation(t, code)
}

func TestCompile_Builtins_Close(t *testing.T) { t.FailNow() }

func TestCompile_Builtins_Complex(t *testing.T) {
	code := `
package main
func main() {
	a = complex(100,8)
}`
	testCompilation(t, code)
}

func TestCompile_Builtins_Copy(t *testing.T) {
	code := `
package main
func main() {
	a := []int{0, 1, 2}
	b := copy(a, a[1:])
}`
	testCompilation(t, code)
}

func TestCompile_Builtins_Delete(t *testing.T) {
	code := `
package main
func main() {
	a := map[string]string {
		"foo": "bar",
	}
	delete(a, "foo")
}`
	testCompilation(t, code)
}

func TestCompile_Builtins_Imag(t *testing.T) {
	code := `
package main
func foo(c complex64) {
	a := imag(c)
}`
	testCompilation(t, code)
}

func TestCompile_Builtins_Len(t *testing.T) {
	code := `
package main
func main() {
	a := len("foo")
}`
	testCompilation(t, code)
}

func TestCompile_Builtins_Make(t *testing.T) {
	code := `
package main
func main() {
	a := make(map[string]int)
}`
	testCompilation(t, code)
}

func TestCompile_Builtins_New(t *testing.T) {
	code := `
package main
func main() {
	a := new(int)
}`
	testCompilation(t, code)
}

func TestCompile_Builtins_Panic(t *testing.T) {
	code := `
package main
func main() {
	panic("foo")
}`
	testCompilation(t, code)
}

func TestCompile_Builtins_Print(t *testing.T) {
	code := `
package main
func main() {
	print("foo")
}`
	testCompilation(t, code)
}

func TestCompile_Builtins_Println(t *testing.T) {
	code := `
package main
func main() {
	println("foo")
}`
	testCompilation(t, code)
}

func TestCompile_Builtins_Real(t *testing.T) {
	code := `
package main
func foo(c complex64) {
	a := real(c)
}`
	testCompilation(t, code)
}

func TestCompile_Builtins_Recover(t *testing.T) {
	code := `
package main
func main() {
	a := recover()
}`
	testCompilation(t, code)
}
