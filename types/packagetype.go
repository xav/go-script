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
	"github.com/xav/go-script/values"
	"github.com/xav/go-script/vm"
)

var packageTypes = make(map[string]*PackageType)

type PackageType struct {
	commonType
	Elems []PackageField
}

type PackageField struct {
	Name string
	Type vm.Type
}

// NewPackageType creates a new package import type.
func NewPackageType(path string, fields []PackageField) *PackageType {
	t, ok := packageTypes[path]
	if !ok {
		t = &PackageType{commonType{}, fields}
		packageTypes[path] = t
	}
	return t
}

// Type interface //////////////////////////////////////////////////////////////

// Compat returns whether this type is compatible with another type.
func (p *PackageType) Compat(o vm.Type, conv bool) bool {
	return false
}

// Lit returns this type's literal.
func (p *PackageType) Lit() vm.Type {
	return p
}

// Zero returns a new zero value of this type.
func (p *PackageType) Zero() vm.Value {
	return &values.PackageV{
		Name:   "",
		Idents: make([]vm.Value, 0),
	}
}

// String returns the string representation of this type.
func (p *PackageType) String() string {
	return "<package>"
}
