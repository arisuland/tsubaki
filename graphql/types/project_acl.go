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

package types

// ACLPermission is the permission type that is tied to a Role or Member.
type ACLPermission string

var (
	// ACL_REPO_UPDATE allows a Member to update the project's metadata.
	ACL_REPO_UPDATE ACLPermission = "REPO_UPDATE"

	// ACL_WRITE allows a Member to write to this project.
	ACL_WRITE ACLPermission = "WRITE"
)

// Role is a ACL role that can be tied to a user to inherit the permissions
// object.
type Role struct {
	// Allowed is a list of permissions available to a Member.
	// This is automatically inherited so the Member's allowed object
	// will be overridden by the Role's allowed object.
	Allowed []ACLPermission `json:"allowed"`

	// Deny is a list of permissions unavailable to a Member.
	// This is automatically inherited so the Member's allowed object
	// will be overridden by the Role's allowed object.
	Deny []ACLPermission `json:"denied"`

	// Name is the Role's name that a Member can inherit the permissions from.
	Name string `json:"name"`
}

// Member is the member that is in the project's ACL
type Member struct {
	// Allowed is the ACLPermission that this Member is allowed to do.
	Allowed []ACLPermission `json:"allowed"`

	// Deny is the ACLPermission that this Member is not allowed to do.
	Deny []ACLPermission `json:"denied"`

	// Roles is a list of Roles that this Member can inherit from.
	Roles []Role `json:"roles"`
}

// ProjectACL is a object that is all the project's ACL
// permissions for a specific user or all users.
type ProjectACL struct {
	// FormatVersion refers to the format version to use.
	// https://docs.arisu.land/projects/acl#format-version
	FormatVersion int32 `json:"format_version"`

	// UpdatedAt refers to a ISO8601 formatted string on when
	// this ACL was updated at. Since this isn't stored in the database,
	// this will refer to the file update date since this is stored
	// under `$project_volume/:user/:project/permissions.hcl`
	UpdatedAt string `json:"updated_at"`

	// CreatedAt refers to a ISO8601 formatted string on when
	// this ACL was created at. Since this isn't stored in the database,
	// this will refer to the file creation date since this is stored
	// under `$project_volume/:user/:project/permissions.hcl`
	CreatedAt string `json:"created_at"`

	// Members refer to the ProjectACL members block.
	Members map[string]Member `json:"members"`
}
