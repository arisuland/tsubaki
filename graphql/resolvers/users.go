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
	"net/mail"
	"time"

	"arisu.land/tsubaki/graphql/types"
	"arisu.land/tsubaki/graphql/types/result"
	"arisu.land/tsubaki/graphql/types/update"
	"arisu.land/tsubaki/pkg/sessions"
	"arisu.land/tsubaki/prisma/db"
	"arisu.land/tsubaki/util"
)

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
	user, err := r.Container.Prisma.User.FindFirst(db.User.ID.Equals(id)).With(db.User.Projects.Fetch()).Exec(ctx)
	if errors.Is(err, db.ErrNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return toSelfUserModel(user, session), nil
}

func (r *Resolver) User(ctx context.Context, args struct {
	ID string
}) (*types.User, error) {
	user, err := r.Container.Prisma.User.FindFirst(db.User.ID.Equals(args.ID)).With(db.User.Projects.Fetch()).Exec(ctx)
	if errors.Is(err, db.ErrNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return types.FromDbModel(user), nil
}

func (r *Resolver) Signup(
	ctx context.Context,
	args struct {
		Email    string
		Password string
		Username string
	},
) (*result.SignupResult, error) {
	// Check if user exists in database
	userByName, err := r.Container.Prisma.User.FindUnique(db.User.Username.Equals(args.Username)).Exec(ctx)
	if err != nil && !errors.Is(err, db.ErrNotFound) {
		return nil, err
	}

	if userByName != nil {
		return nil, errors.New(fmt.Sprintf("user with name %s exists.", args.Username))
	}

	// Check if email is valid
	_, err = mail.ParseAddress(args.Email)
	if err != nil {
		return nil, err
	}

	// Check if email exists in db
	userByEmail, err := r.Container.Prisma.User.FindUnique(db.User.Email.Equals(args.Email)).Exec(ctx)
	if err != nil && !errors.Is(err, db.ErrNotFound) {
		return nil, err
	}

	if userByEmail != nil {
		return nil, errors.New(fmt.Sprintf("user with email %s already exists.", args.Email))
	}

	// We're good to go! Let's now
	salt := util.GenerateHash(16)
	if salt == "" {
		return nil, errors.New("unable to generate password salt")
	}

	hash, err := util.GeneratePassword(args.Password)
	if err != nil {
		return nil, err
	}

	id := r.Container.Snowflake.Generate().String()
	user, err := r.Container.Prisma.User.CreateOne(
		db.User.Username.Set(args.Username),
		db.User.Password.Set(hash),
		db.User.Email.Set(args.Email),
		db.User.ID.Set(id),
		db.User.Projects.Link(),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &result.SignupResult{
		Success: true,
		Errors:  []result.Error{},
		User:    types.FromDbModel(user),
	}, nil
}

func (r *Resolver) DeleteUser(ctx context.Context) result.Result {
	id := ""
	uid := ctx.Value("user_id")
	if uid != nil {
		id = uid.(string)
	}

	if id == "" {
		return result.ErrWithMessage("Missing session token", -1)
	}

	// get session from context
	session := sessions.Sessions.Get(id)
	if session == nil {
		return result.ErrWithMessage("Missing session token", -1)
	}

	_, err := r.Container.Prisma.User.FindUnique(db.User.ID.Equals(id)).Delete().Exec(ctx)
	if err != nil {
		return result.Err([]error{err})
	}

	return result.Ok()
}

func (r *Resolver) UpdateUser(ctx context.Context, args struct {
	Args update.UpdateUserArgs
}) result.Result {
	id := ""
	uid := ctx.Value("user_id")
	if uid != nil {
		id = uid.(string)
	}

	if id == "" {
		return result.ErrWithMessage("Missing session token", -1)
	}

	_, err := r.Container.Prisma.User.FindUnique(db.User.ID.Equals(id)).Exec(ctx)
	if err != nil {
		return result.Err([]error{err})
	}

	hasUpdated := false
	if args.Args.Description != nil {
		desc := *args.Args.Description
		if len(desc) > 240 {
			return result.ErrWithMessage("Length for string was over the maximum amount of characters.", 1008)
		}

		_, err := r.Container.Prisma.User.FindUnique(
			db.User.ID.Equals(id),
		).Update(db.User.Description.Set(desc)).Exec(ctx)

		if err != nil {
			return result.Err([]error{err})
		}

		if !hasUpdated {
			hasUpdated = true
		}
	}

	if args.Args.Username != nil {
		uname := *args.Args.Username

		_, err := r.Container.Prisma.User.FindUnique(db.User.Username.Equals(uname)).Exec(ctx)
		if err == nil {
			return result.ErrWithMessage(fmt.Sprintf("Username %s is taken already.", uname), 1009)
		}

		_, err = r.Container.Prisma.User.FindUnique(
			db.User.ID.Equals(id),
		).Update(db.User.Username.Set(uname)).Exec(ctx)

		if err != nil {
			return result.Err([]error{err})
		}

		if !hasUpdated {
			hasUpdated = true
		}
	}

	if args.Args.Name != nil {
		name := *args.Args.Name

		if len(name) > 60 {
			return result.ErrWithMessage("Length for string was over the maximum amount of characters.", 1008)
		}

		_, err = r.Container.Prisma.User.FindUnique(
			db.User.ID.Equals(id),
		).Update(db.User.Name.Set(name)).Exec(ctx)

		if err != nil {
			return result.Err([]error{err})
		}

		if !hasUpdated {
			hasUpdated = true
		}
	}

	if !hasUpdated {
		return result.ErrWithMessage("Unable to update user: nothing was updated.", 1010)
	}

	return result.Ok()
}
