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

package vm

import "math/big"

// BoundedType represents a numeric type with boundaries on its value.
type BoundedType interface {
	Type
	MinVal() *big.Rat // MinVal returns the smallest value of this type.
	MaxVal() *big.Rat // MaxVal returns the largest value of this type.
}
