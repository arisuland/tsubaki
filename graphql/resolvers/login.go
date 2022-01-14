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

import "context"

func (r *Resolver) Login(
	ctx context.Context,
	usernameOrEmail string,
	password string,
) error {
	return nil
}

// func (c LoginController) Login(
// 	ctx context.Context,
// 	usernameOrEmail string,
// 	password string,
// ) result.LoginResult {
// 	user, err := c.Prisma.Client.User.FindUnique(db.User.Username.Equals(usernameOrEmail)).Exec(ctx)
// 	if err != nil {
// 		if errors.Is(err, db.ErrNotFound) {
// 			user, err = c.Prisma.Client.User.FindUnique(db.User.Email.Equals(usernameOrEmail)).Exec(ctx)
// 			if err != nil {
// 				if errors.Is(err, db.ErrNotFound) {
// 					return result.LoginResult{
// 						Success: false,
// 						Errors: []result.Error{
// 							{
// 								Message: fmt.Sprintf("Username or email %s doesn't exist.", usernameOrEmail),
// 								Code:    -1,
// 							},
// 						},

// 						Token: "",
// 					}
// 				}

// 				return result.LoginResult{
// 					Success: false,
// 					Errors: []result.Error{
// 						{
// 							Message: fmt.Sprintf("Unable to retrieve user: %v", err),
// 							Code:    -1,
// 						},
// 					},

// 					Token: "",
// 				}
// 			}
// 		} else {
// 			return result.LoginResult{
// 				Success: false,
// 				Errors: []result.Error{
// 					{
// 						Message: fmt.Sprintf("Unable to retrieve user: %v", err),
// 						Code:    -1,
// 					},
// 				},

// 				Token: "",
// 			}
// 		}
// 	}

// 	match, err := util.VerifyPassword(password, user.Password)
// 	if err != nil {
// 		return result.LoginResult{
// 			Success: false,
// 			Errors: []result.Error{
// 				{
// 					Message: fmt.Sprintf("Unable to decode password: %v", err),
// 					Code:    -1,
// 				},
// 			},

// 			Token: "",
// 		}
// 	}

// 	if !match {
// 		code, msg := types.Get(1002)
// 		return result.LoginResult{
// 			Success: false,
// 			Errors: []result.Error{
// 				{
// 					Message: msg,
// 					Code:    int32(code),
// 				},
// 			},

// 			Token: "",
// 		}
// 	}

// 	// create session
// 	session := sessions.Sessions.New(user.ID)
// 	if session == nil {
// 		return result.LoginResult{
// 			Success: false,
// 			Errors: []result.Error{
// 				{
// 					Message: "Unable to generate session for user.",
// 					Code:    -1,
// 				},
// 			},

// 			Token: "",
// 		}
// 	}

// 	return result.LoginResult{
// 		Success: true,
// 		Errors:  make([]result.Error, 0),
// 		Token:   session.Token,
// 	}
// }

// func (c LoginController) Logout(uid string) bool {
// 	session := sessions.Sessions.Get(uid)
// 	if session == nil {
// 		return false
// 	}

// 	sessions.Sessions.Delete(uid)
// 	return true
// }
