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

func TestRange_LitString(t *testing.T) {
	code := `
package main

func main() {
	for i, v := range "abcd" {
		print(i, "/", v, "\n")
	}
}`
	err := compile(t, code)
	Ok(t, err)
}

func TestRange_CompositeLitArray(t *testing.T) {
	code := `
package main

func main() {
	for i, v := range [4]int{1, 2, 3, 4} {
		print(i, "/", v, "\n")
	}
}`
	err := compile(t, code)
	Ok(t, err)
}

func TestRange_CompositeLitSlice(t *testing.T) {
	code := `
package main

func main() {
	for i, v := range []int{1, 2, 3, 4} {
		print(i, "/", v, "\n")
	}
}`
	err := compile(t, code)
	Ok(t, err)
}

func TestRange_CompositeLitMap(t *testing.T) {
	code := `
package main

func main() {
	for i, v := range map[string]int{"a": 1, "b": 2, "c": 3} {
		print(i, "/", v, "\n")
	}
}
`
	err := compile(t, code)
	Ok(t, err)
}

func TestRange_IdentString(t *testing.T) {
	code := `
package main

func main() {
	s := "abcd"
	for i, v := range s {
		print(i, "/", v, "\n")
	}
}`
	err := compile(t, code)
	Ok(t, err)
}

func TestRange_IdentArray(t *testing.T) {
	code := `
package main

func main() {
	n := [4]int{1, 2, 3, 4}
	for i, v := range n {
		print(i, "/", v, "\n")
	}
}`
	err := compile(t, code)
	Ok(t, err)
}

func TestRange_IdentArrayPointer(t *testing.T) {
	code := `
package main

func main() {
	numbers := [4]int{1, 2, 3, 4}
	for i, v := range &numbers {
		print(i, "/", v, "\n")
	}
}`
	err := compile(t, code)
	Ok(t, err)
}

func TestRange_IdentArrayPointerValue(t *testing.T) {
	code := `
package main

func main() {
	numbers := [4]int{1, 2, 3, 4}
	np := &numbers
	for i, v := range *np {
		print(i, "/", v, "\n")
	}
}`
	err := compile(t, code)
	Ok(t, err)
}

func TestRange_IdentSlice(t *testing.T) {
	code := `
package main

func main() {
	numbers := []int{1, 2, 3, 4}
	for i, v := range numbers {
		print(i, "/", v, "\n")
	}
}`
	err := compile(t, code)
	Ok(t, err)
}

func TestRange_IdentMap(t *testing.T) {
	code := `
package main

func main() {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	for i, v := range m {
		print(i, "/", v, "\n")
	}
}
`
	err := compile(t, code)
	Ok(t, err)
}

func TestRange_Channel(t *testing.T) {
	code := `
package main

func main() {
    queue := make(chan string, 2)
    queue <- "one"
    queue <- "two"
    close(queue)

    for elem := range queue {
        fmt.Println(elem)
		print(elem, "\n")
    }
}
`
	err := compile(t, code)
	Ok(t, err)
}

func TestRange_RangeFailOnIntExpr(t *testing.T) {
	code := `
package main

func main() {
	for i, v := range 12 {
		print(i, "/", v, "\n")
	}
}`
	err := compile(t, code)
	ScannerErrEquals(t, "5:20: cannot range over 12 (type INT)", err)
}
