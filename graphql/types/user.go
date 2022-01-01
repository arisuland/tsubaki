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
	"arisu.land/tsubaki/prisma/db"
	"time"
)

// User is the user metadata from the database.
type User struct {
	Description *string   `json:"description"`
	UpdatedAt   string    `json:"updated_at"`
	CreatedAt   string    `json:"created_at"`
	Username    string    `json:"username"`
	Disabled    bool      `json:"disabled"`
	Projects    []Project `json:"projects"`
	Flags       int32     `json:"flags"`
	Name        *string   `json:"name"`
	ID          string    `json:"id"`
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

// FromUserModel returns a new User entity based off the db result.
func FromUserModel(model *db.UserModel) *User {
	var name *string
	var desc *string

	description, ok := model.Description()
	if !ok {
		desc = nil
	} else {
		desc = &description
	}

	n, ok := model.Name()
	if !ok {
		name = nil
	} else {
		name = &n
	}

	projects := make([]Project, 0)
	if model.RelationsUser.Projects != nil {
		for _, proj := range model.RelationsUser.Projects {
			projects = append(projects, FromProjectModel(&proj))
		}
	}

	return &User{
		Description: desc,
		UpdatedAt:   model.UpdatedAt.Format(time.RFC3339),
		CreatedAt:   model.CreatedAt.Format(time.RFC3339),
		Username:    model.Username,
		Disabled:    model.Disabled,
		Projects:    projects,
		Flags:       int32(model.Flags),
		Name:        name,
		ID:          model.ID,
	}
}
