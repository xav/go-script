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

package types

import (
	"fmt"
	"go/token"

	"github.com/xav/go-script/vm"
	"github.com/xav/horus/warden/vm/value"
)

// NamedType represents a user defined type.
//
// If Incomplete is true, Def will be nil.
// If Incomplete is false and Def is still nil, then this is a placeholder type representing an error.
type NamedType struct {
	NTPos      token.Pos
	Name       string
	Def        vm.Type           // Underlying type.
	Incomplete bool              // True while this type is being defined.
	Methods    map[string]Method //
}

type Method struct {
	decl *FuncDecl
	fn   value.Func
}

// Pos returns the position of the definition in the source code
func (t *NamedType) Pos() token.Pos {
	return t.NTPos
}

// Type interface //////////////////////////////////////////////////////////////

// Compat returns whether this type is compatible with another type.
func (t *NamedType) Compat(o vm.Type, conv bool) bool {
	if t2, ok := o.(*NamedType); ok {
		if conv {
			// Two named types are conversion compatible if their literals are conversion compatible.
			return t.Def.Compat(t2.Def, conv)
		}

		// Two named types are compatible if their type names originate in the same type declaration.
		return t == t2
	}

	// A named and an unnamed type are compatible if the respective type literals are compatible.
	return o.Compat(t.Def, conv)
}

// Lit returns the underlying type's literal.
func (t *NamedType) Lit() vm.Type {
	return t.Def.Lit()
}

// IsBoolean returns true if this is a boolean type.
func (t *NamedType) IsBoolean() bool {
	return t.Def.IsBoolean()
}

// IsInteger returns true if this is an integer type.
func (t *NamedType) IsInteger() bool {
	return t.Def.IsInteger()
}

// IsFloat returns true if this is a floating type.
func (t *NamedType) IsFloat() bool {
	return t.Def.IsFloat()
}

// IsIdeal returns true if this represents an ideal value.
func (t *NamedType) IsIdeal() bool {
	return false
}

// Zero returns a new zero value of this type.
func (t *NamedType) Zero() vm.Value {
	return t.Def.Zero()
}

// String returns the string representation of this type.
func (t *NamedType) String() string {
	return t.Name
}

////////////////////////////////////////////////////////////////////////////////

// Complete marks this named type as completed
func (t *NamedType) Complete(def vm.Type) {
	if !t.Incomplete {
		logger.Panic().
			Str("type", fmt.Sprintf("%+v", *t)).
			Msg("cannot complete already completed NamedType")
	}
	// We strip the name from def because multiple levels of naming are useless.
	if ndef, ok := def.(*NamedType); ok {
		def = ndef.Def
	}
	t.Def = def
	t.Incomplete = false
}
