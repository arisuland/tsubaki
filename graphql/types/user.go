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

import (
	"time"

	"arisu.land/tsubaki/prisma/db"
)

// User is a registered account, a invited account, or a created account
// to represents someone. This holds external metadata about them but
// doesn't leak sensitive information like emails or passwords.
type User struct {
	// Returns the User's description. This can be `null` by default
	// if not edited.
	Description *string `json:"description"`

	// Returns a ISO-8601 formatted string of when the user
	// has updated their account.
	UpdatedAt string `json:"updated_at"`

	// Returns a ISO-8601 formatted string of when the user
	// has created their account.
	CreatedAt string `json:"created_at"`

	// The user's unique username.
	Username string `json:"username"`

	// If the user's account is disabled.
	Disabled bool `json:"disabled"`

	// Returns a list of projects belonging to this user.
	Projects []Project `json:"projects"`

	// Returns the user's flags. The flags are as followed:
	//
	//  - 1 << 0: Admin - They are an administrator of this instance. - You can obtain this flag by being an administrator.
	//  - 1 << 1: Employee - This person is a Noelware employee. - You cannot obtain this flag.
	//  - 1 << 2: Cutie - ... - You cannot obtain this flag.
	//  - 1 << 3: Founder - This person is the founder of Arisu! You cannot obtain this flag.
	Flags int32 `json:"flags"`

	// The user's display name. This can be `null` if not edited.
	Name *string `json:"name"`

	// The user's ID.
	ID string `json:"id"`
}

// SelfUser is the user metadata with the session metadata
type SelfUser struct {
	SessionExpiresIn string    `json:"session_expires_in"`
	SessionType      string    `json:"session_type"`
	Description      *string   `json:"description"`
	UpdatedAt        string    `json:"updated_at"`
	CreatedAt        string    `json:"created_at"`
	Username         string    `json:"username"`
	Disabled         bool      `json:"disabled"`
	Projects         []Project `json:"projects"`
	Flags            int32     `json:"flags"`
	Name             *string   `json:"name"`
	ID               string    `json:"id"`
}

// FromDbModel is a function to convert the db.UserModel into a typed User.
func FromDbModel(user *db.UserModel) *User {
	var name *string
	var desc *string

	description, ok := user.Description()
	if !ok {
		desc = nil
	} else {
		desc = &description
	}

	n, ok := user.Name()
	if !ok {
		name = nil
	} else {
		name = &n
	}

	projects := make([]Project, 0)
	if user.RelationsUser.Projects != nil {
		for _, proj := range user.RelationsUser.Projects {
			projects = append(projects, FromProjectModel(&proj))
		}
	}

	return &User{
		Description: desc,
		UpdatedAt:   user.UpdatedAt.Format(time.RFC3339),
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
		Username:    user.Username,
		Disabled:    user.Disabled,
		Projects:    projects,
		Flags:       int32(user.Flags),
		Name:        name,
		ID:          user.ID,
	}
}
