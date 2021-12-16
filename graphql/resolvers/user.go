// ☔ Arisu: Translation made with simplicity, yet robust.
// Copyright (C) 2020-2021 Noelware
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
	"arisu.land/tsubaki/graphql/types"
	"context"
)

// User retrieves a user from the database based off its ID.
func (r *Resolver) User(ctx context.Context, args struct{ ID string }) (*types.User, error) {
	user, err := r.Users.GetUser(ctx, args.ID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	return types.FromUserModel(user), nil
}

func (r *Resolver) CreateUser(
	ctx context.Context,
	args struct {
	Email    string
	Password string
	Username string
},
) (*types.User, error) {
	user, err := r.Users.CreateUser(ctx, args.Email, args.Password, args.Username)
	if err != nil {
		return nil, err
	}

	return types.FromUserModel(user), nil
}
