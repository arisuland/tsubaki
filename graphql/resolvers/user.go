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
	"arisu.land/tsubaki/graphql/types"
	"arisu.land/tsubaki/graphql/types/update"
	"arisu.land/tsubaki/pkg/sessions"
	"arisu.land/tsubaki/prisma/db"
	"context"
	"errors"
	"time"
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

func (r *Resolver) UpdateUser(
	ctx context.Context,
	args struct {
		Args update.UserArgs
	},
) (bool, error) {
	id := ""
	uid := ctx.Value("user_id")
	if uid != nil {
		id = uid.(string)
	}

	if id == "" {
		return false, errors.New("missing session or bearer token")
	}

	// get session from context
	session := sessions.Sessions.Get(id)
	if session == nil {
		return false, nil
	}

	err := r.Users.UpdateUser(ctx, session.User.ID, args.Args)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *Resolver) DeleteUser(ctx context.Context) bool {
	id := ""
	uid := ctx.Value("user_id")
	if uid != nil {
		id = uid.(string)
	}

	if id == "" {
		return false
	}

	// get session from context
	session := sessions.Sessions.Get(id)
	if session == nil {
		return false
	}

	err := r.Users.DeleteUser(ctx, session.User.ID)
	if err != nil {
		return false
	}

	return true
}

func (r *Resolver) ReenableUser(ctx context.Context) bool {
	id := ""
	uid := ctx.Value("user_id")
	if uid != nil {
		id = uid.(string)
	}

	if id == "" {
		return false
	}

	// get session from context
	session := sessions.Sessions.Get(id)
	if session == nil {
		return false
	}

	err := r.Users.ReenableUser(ctx, session.User.ID)
	if err != nil {
		return false
	}

	return true
}

func (r *Resolver) DisableUser(ctx context.Context) bool {
	id := ""
	uid := ctx.Value("user_id")
	if uid != nil {
		id = uid.(string)
	}

	if id == "" {
		return false
	}

	// get session from context
	session := sessions.Sessions.Get(id)
	if session == nil {
		return false
	}

	err := r.Users.DisableUser(ctx, session.User.ID)
	if err != nil {
		return false
	}

	return true
}

func (r *Resolver) Me(ctx context.Context) (*types.SelfUser, error) {
	id := ""
	uid := ctx.Value("user_id")
	if uid != nil {
		id = uid.(string)
	}

	if id == "" {
		return nil, nil
	}

	// get session from context
	session := sessions.Sessions.Get(id)
	if session == nil {
		return nil, nil
	}

	// owo
	user, err := r.Users.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	return toSelfUserModel(user, session), nil
}

func toSelfUserModel(model *db.UserModel, session *sessions.Session) *types.SelfUser {
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

	projects := make([]types.Project, 0)
	if model.RelationsUser.Projects != nil {
		for _, proj := range model.RelationsUser.Projects {
			projects = append(projects, types.FromProjectModel(&proj))
		}
	}

	return &types.SelfUser{
		SessionExpiresIn: session.ExpiresIn.Format(time.RFC3339),
		SessionType:      session.Type.String(),
		Description:      desc,
		UpdatedAt:        model.UpdatedAt.Format(time.RFC3339),
		CreatedAt:        model.CreatedAt.Format(time.RFC3339),
		Username:         model.Username,
		Disabled:         model.Disabled,
		Projects:         projects,
		Flags:            int32(model.Flags),
		Name:             name,
		ID:               model.ID,
	}
}
