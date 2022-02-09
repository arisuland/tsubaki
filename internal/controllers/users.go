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

package controllers

import (
	"arisu.land/tsubaki/pkg"
	"arisu.land/tsubaki/pkg/result"
	"arisu.land/tsubaki/prisma/db"
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

type UserController struct{}

func newUserController() UserController {
	return UserController{}
}

// User is the underlying User structure that is returned using the Users API.
// Due to some data we store, some of this shouldn't be present!
type User struct {
	// Returns the Gravatar email if `use_gravatar` is set to true. This
	// can be nil if this account has decided to use this account's email
	// address as the email address.
	GravatarEmail *string `json:"gravatar_email"`

	// If this account has opted to use Gravatar for their avatar.
	UseGravatar bool `json:"use_gravatar"`

	// The account's avatar URL from the unified storage. (NOT AVAILABLE)
	AvatarUrl *string `json:"avatar_url"`

	// Returns the account's description, can be `nil`
	// if no description was provided.
	Description *string `json:"description"`

	// Returns a RFC3339 timestamp of when this account's
	// metadata has been updated.
	UpdatedAt string `json:"updated_at"`

	// Returns a RFC3339-compliant timestamp of when this account
	// was registered at.
	CreatedAt string `json:"created_at"`

	// Returns the account's unique username that can be
	// displayed at `<FUBUKI_URL>/@{username}`
	Username string `json:"username"`

	// Returns a boolean if this account is currently
	// disabled by the administrators or not.
	Disabled bool `json:"disabled"`

	// Returns the account public flags that represented
	// their permissions (i.e, admin).
	Flags int `json:"flags"`

	// The user's display name, this can be `nil` if none
	// was set.
	Name *string `json:"name"`

	// Returns this account's ID that can be queried from the API.
	ID string `json:"id"`
}

func fromUserModel(user *db.UserModel) *User {
	return &User{
		GravatarEmail: user.InnerUser.GravatarEmail,
		UseGravatar:   user.UseGravatar,
		AvatarUrl:     user.InnerUser.AvatarURL,
		Description:   user.InnerUser.Description,
		UpdatedAt:     user.InnerUser.UpdatedAt.Format(time.RFC3339),
		CreatedAt:     user.InnerUser.CreatedAt.Format(time.RFC3339),
		Username:      user.InnerUser.Username,
		Disabled:      user.InnerUser.Disabled,
		Flags:         user.InnerUser.Flags,
		Name:          user.InnerUser.Name,
		ID:            user.InnerUser.ID,
	}
}

func (UserController) Get(id string) *result.Result {
	user, err := pkg.GlobalContainer.Prisma.User.FindUnique(
		db.User.ID.Equals(id)).Exec(context.TODO())

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return result.Err(404, "USER_NOT_FOUND", fmt.Sprintf("user with id %s was not found.", id))
		} else {
			logrus.Errorf("Unable to retrieve user %s from the database: %v", id, err)
			return result.Err(500, "UNKNOWN_ERROR", fmt.Sprintf("unknown error while retrieving user %s...", id))
		}
	}

	return result.Ok(fromUserModel(user))
}

func (UserController) Create(username string, password string, email string) *result.Result {
	return result.Ok(struct{}{})
}
