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

// Type is the common interface for all supported types
type Type interface {
	// Compat returns whether this type is compatible with another type.
	// If conv is false, this is normal compatibility, where two named types are compatible only if they are the same named type.
	// If conv if true, this is conversion compatibility, where two named types are conversion compatible if their definitions are conversion compatible.
	Compat(t Type, conv bool) bool // TODO: Deal with recursive types

	// Lit returns this type's literal.
	// If this is a named type, this is the unnamed underlying type.
	// Otherwise, this is an identity operation.
	Lit() Type

	// IsBoolean returns true if this is a boolean type.
	IsBoolean() bool

	// IsInteger returns true if this is an integer type.
	IsInteger() bool

	// IsFloat returns true if this is a floating type.
	IsFloat() bool

	// IsIdeal returns true if this represents an ideal value.
	IsIdeal() bool

	// Zero returns a new zero value of this type.
	Zero() Value

	// String returns the string representation of this type.
	String() string
}
