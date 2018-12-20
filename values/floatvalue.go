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

type FloatValue interface {
	vm.Value
	Get(*vm.Thread) float64
	Set(*vm.Thread, float64)
}

// float32 /////////////////////////////////////////////////////////////////////

type Float32V float32

func (v *Float32V) String() string {
	return fmt.Sprint(*v)
}

func (v *Float32V) Assign(t *vm.Thread, o vm.Value) {
	*v = Float32V(o.(FloatValue).Get(t))
}

func (v *Float32V) Get(*vm.Thread) float64 {
	return float64(*v)
}

func (v *Float32V) Set(t *vm.Thread, x float64) {
	*v = Float32V(x)
}

// float64 /////////////////////////////////////////////////////////////////////

type Float64V float64

func (v *Float64V) String() string {
	return fmt.Sprint(*v)
}

func (v *Float64V) Assign(t *vm.Thread, o vm.Value) {
	*v = Float64V(o.(FloatValue).Get(t))
}

func (v *Float64V) Get(*vm.Thread) float64 {
	return float64(*v)
}

func (v *Float64V) Set(t *vm.Thread, x float64) {
	*v = Float64V(x)
}
