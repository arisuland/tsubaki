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
	"arisu.land/tsubaki/graphql/types/result"
	"context"
)

func (r *Resolver) Login(ctx context.Context, args struct {
	UsernameOrEmail string
	Password        string
}) result.LoginResult {
	return r.login.Login(ctx, args.UsernameOrEmail, args.Password)
}

func (r *Resolver) Logout(ctx context.Context) bool {
	id := ""
	uid := ctx.Value("user_id")
	if uid != nil {
		id = uid.(string)
	}

	if id == "" {
		return false
	}

	return r.login.Logout(id)
}
