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

import (
	"fmt"
)

// SliceError is used to abort a thread when invalid indices are used for
// a slice manipulation of a Slice, Array or String.
type SliceError struct {
	Lo, Hi, Cap int64
}

func (e SliceError) Error() string {
	return fmt.Sprintf("slice [%d:%d]; cap %d", e.Lo, e.Hi, e.Cap)
}
