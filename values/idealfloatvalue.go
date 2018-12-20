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
	"math/big"

	"github.com/xav/go-script/vm"
)

// TODO: Remove IdealFloatValue since ideals are not l-values.

type IdealFloatValue interface {
	vm.Value
	Get() *big.Rat
}

type IdealFloatV struct {
	V *big.Rat
}

func (v *IdealFloatV) String() string {
	return v.V.FloatString(6)
}

func (v *IdealFloatV) Assign(t *vm.Thread, o vm.Value) {
	v.V = o.(IdealFloatValue).Get()
}

func (v *IdealFloatV) Get() *big.Rat {
	return v.V
}
