// ☔ Arisu: Translation made with simplicity, yet robust.
// Copyright (C) 2020-2022 Noelware
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package acl

import (
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/sirupsen/logrus"
)

// FormatVersion is the format version to use. Defaults to Version1.
type FormatVersion int

var (
	// Version1 is referred to the first version of the format version.
	// Following versions will have a way to migrate from version to version.
	Version1 FormatVersion = 1
)

// Int will convert the FormatVersion to its original value.
func (v FormatVersion) Int() int {
	switch v {
	case Version1:
		return 1

	default:
		return -1
	}
}

// String stringifies a FormatVersion type.
func (v FormatVersion) String() string {
	switch v {
	case Version1:
		return "version 1"

	default:
		return "unknown version"
	}
}

// Permission is a type to represent an ACL permission.
type Permission string

var (
	// REPO_UPDATE allows a AclMember to update this project's metadata.
	REPO_UPDATE Permission = "REPO_UPDATE"

	// WRITE allows a AclMember to write to this project.
	WRITE Permission = "WRITE"
)

// Role is the role a Member can inherit permissions from.
type Role struct {
	// List of Permission objects that a Member can inherit from to grant actions.
	Allow []Permission `hcl:"allow"`

	// List of Permission objects that a Member can inherit from to prohibit actions.
	Deny []Permission `hcl:"deny"`
}

// Member is a member object that is structured to grant or deny permissions.
type Member struct {
	// List of Permission objects that a Member can inherit from to grant actions.
	Allow []Permission `hcl:"allow"`

	// List of Permission objects that a Member can inherit from to prohibit actions.
	Deny []Permission `hcl:"deny"`
}

// Object is the HCL structure for the project's ACL.
type Object struct {
	// FormatVersion refers to the object's format version.
	FormatVersion FormatVersion `hcl:"formatVersion"`

	// Members is a list of Member objects to grant or deny permissions.
	Members map[string]Member `hcl:"members"`
}

// DecodeFromSource decodes the contents and converts it to an ACL Object.
func DecodeFromSource(contents string) (*Object, error) {
	var object Object
	err := hclsimple.Decode("permissions.hcl", []byte(contents), nil, &object)
	if err != nil {
		logrus.Fatalf("Unable to load project ACL: %v", err)
		return nil, err
	}

	return &object, nil
}
