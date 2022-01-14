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

package resolvers

import (
	"context"
	"errors"
	"fmt"

	"arisu.land/tsubaki/graphql/types/result"
	"arisu.land/tsubaki/pkg/sessions"
	"arisu.land/tsubaki/prisma/db"
	"arisu.land/tsubaki/util"
)

func (r *Resolver) Login(
	ctx context.Context,
	usernameOrEmail string,
	password string,
) (*result.LoginResult, error) {
	user, err := r.Container.Prisma.User.FindUnique(
		db.User.Username.Equals(usernameOrEmail),
	).Exec(context.TODO())

	if err != nil {
		// Check if the username wasn't found, so let's check the email!
		if errors.Is(err, db.ErrNotFound) {
			user, err = r.Container.Prisma.User.FindUnique(
				db.User.Email.Equals(usernameOrEmail),
			).Exec(context.TODO())

			if err != nil {
				if errors.Is(err, db.ErrNotFound) {
					return &result.LoginResult{
						Token:   "",
						Success: false,
						Errors: []result.Error{
							{
								Message: fmt.Sprintf("Username or email %s doesn't exist.", usernameOrEmail),
								Code:    -1,
							},
						},
					}, nil
				}

				return &result.LoginResult{
					Token:   "",
					Success: false,
					Errors: []result.Error{
						{
							Message: fmt.Sprintf("Unable to retrieve user: %v", err),
							Code:    -1,
						},
					},
				}, nil
			}
		} else {
			return &result.LoginResult{
				Token:   "",
				Success: false,
				Errors: []result.Error{
					{
						Message: fmt.Sprintf("Unable to retrieve user: %v", err),
						Code:    -1,
					},
				},
			}, nil
		}
	}

	// Check if the password is correct
	match, err := util.VerifyPassword(password, user.Password)
	if err != nil {
		return &result.LoginResult{
			Token:   "",
			Success: false,
			Errors: []result.Error{
				{
					Message: "Unable to decode password.",
					Code:    -1,
				},
			},
		}, nil
	}

	if !match {
		return &result.LoginResult{
			Token:   "",
			Success: false,
			Errors: []result.Error{
				{
					Message: "Invalid password.",
					Code:    int32(1002),
				},
			},
		}, nil
	}

	// Create the session
	sess := sessions.Sessions.New(user.ID)
	if sess == nil {
		return &result.LoginResult{
			Token:   "",
			Success: false,
			Errors: []result.Error{
				{
					Message: "Unable to generate session for user.",
					Code:    -1,
				},
			},
		}, nil
	}

	return &result.LoginResult{
		Success: true,
		Errors:  []result.Error{},
		Token:   sess.Token,
	}, nil
}

func (r *Resolver) Logout(ctx context.Context) (bool, error) {
	id := ""
	uid := ctx.Value("user_id")
	if uid != nil {
		id = uid.(string)
	}

	if id == "" {
		return false, errors.New("missing session or bearer token")
	}

	sessions.Sessions.Delete(id)
	return true, nil
}
