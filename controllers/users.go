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

package controllers

import (
	"arisu.land/tsubaki/graphql/types"
	"arisu.land/tsubaki/graphql/types/update"
	"arisu.land/tsubaki/pkg/managers"
	"arisu.land/tsubaki/pkg/sessions"
	"arisu.land/tsubaki/pkg/util"
	"arisu.land/tsubaki/prisma/db"
	"context"
	"errors"
	"fmt"
	"net/mail"
)

// UserController is a struct that handles fetching, updating, deleting, and creating users.
type UserController struct {
	Snowflake *managers.SnowflakeManager

	// Prisma is the attached Prisma client to use.
	Prisma managers.Prisma
}

// NewUserController creates a new UserController instance.
func NewUserController(prisma managers.Prisma, snowflake *managers.SnowflakeManager) UserController {
	return UserController{
		Snowflake: snowflake,
		Prisma:    prisma,
	}
}

// GetUser returns a db.UserModel instance or `nil` if it was not found. If any errors occur,
// the `error` will be filled in.
func (m UserController) GetUser(ctx context.Context, id string) (*db.UserModel, error) {
	user, err := m.Prisma.Client.User.FindFirst(db.User.ID.Equals(id)).With(db.User.Projects.Fetch()).Exec(ctx)
	if errors.Is(err, db.ErrNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

// CreateUser creates a new user entry in the database.
func (m UserController) CreateUser(
	ctx context.Context,
	email string,
	password string,
	username string,
) (*db.UserModel, error) {
	// Check if user exists in database
	userByName, err := m.Prisma.Client.User.FindUnique(db.User.Username.Equals(username)).Exec(ctx)
	if err != nil && !errors.Is(err, db.ErrNotFound) {
		return nil, err
	}

	if userByName != nil {
		return nil, errors.New(fmt.Sprintf("user with name %s exists.", username))
	}

	// Check if email is valid
	_, err = mail.ParseAddress(email)
	if err != nil {
		return nil, err
	}

	// Check if email exists in db
	userByEmail, err := m.Prisma.Client.User.FindUnique(db.User.Email.Equals(email)).Exec(ctx)
	if err != nil && !errors.Is(err, db.ErrNotFound) {
		return nil, err
	}

	if userByEmail != nil {
		return nil, errors.New(fmt.Sprintf("user with email %s already exists.", email))
	}

	// We're good to go! Let's now
	salt := util.GenerateHash(16)
	if salt == "" {
		return nil, errors.New("unable to generate password salt")
	}

	hash, err := util.GeneratePassword(password)
	if err != nil {
		return nil, err
	}

	id := m.Snowflake.Generate()
	user, err := m.Prisma.Client.User.CreateOne(
		db.User.Username.Set(username),
		db.User.Password.Set(hash),
		db.User.Email.Set(email),
		db.User.ID.Set(id),
		db.User.Projects.Link(),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (c UserController) UpdateUser(
	ctx context.Context,
	uid string,
	args update.UserArgs,
) error {
	_, err := c.Prisma.Client.User.FindUnique(db.User.ID.Equals(uid)).Exec(ctx)
	if err != nil {
		return err
	}

	if args.Description != nil {
		desc := *args.Description
		if len(desc) > 240 {
			code, msg := types.Get(1008)
			return fmt.Errorf("%d: %s", code, msg)
		}

		_, err := c.Prisma.Client.User.FindUnique(db.User.ID.Equals(uid)).Update(db.User.Description.Set(desc)).Exec(ctx)
		if err != nil {
			return err
		}
	}

	if args.Username != nil {
		username := *args.Username
		existing, err := c.Prisma.Client.User.FindUnique(db.User.Username.Equals(username)).Exec(ctx)
		if err != nil {
			return err
		}

		if existing != nil {
			code, msg := types.Get(1004)
			return fmt.Errorf("%d: %s", code, msg)
		}

		_, err = c.Prisma.Client.User.FindUnique(db.User.ID.Equals(uid)).Update(db.User.Username.Set(username)).Exec(ctx)
		if err != nil {
			return err
		}
	}

	if args.Name != nil {
		name := *args.Name
		if len(name) > 60 {
			code, msg := types.Get(1008)
			return fmt.Errorf("%d: %s", code, msg)
		}

		_, err := c.Prisma.Client.User.FindUnique(db.User.ID.Equals(uid)).Update(db.User.Name.Set(name)).Exec(ctx)
		if err != nil {
			return err
		}
	}

	return errors.New("no update query was performed")
}

func (c UserController) DeleteUser(
	ctx context.Context,
	uid string,
) error {
	_, err := c.Prisma.Client.User.FindUnique(db.User.ID.Equals(uid)).Exec(ctx)
	if err != nil {
		return err
	}

	// now we delete
	_, err = c.Prisma.Client.User.FindUnique(db.User.ID.Equals(uid)).Delete().Exec(ctx)
	if err != nil {
		return err
	}

	// delete all sessions
	sessions.Sessions.Delete(uid)
	return nil
}

func (c UserController) DisableUser(
	ctx context.Context,
	uid string,
) error {
	_, err := c.Prisma.Client.User.FindUnique(db.User.ID.Equals(uid)).Exec(ctx)
	if err != nil {
		return err
	}

	// now we delete
	_, err = c.Prisma.Client.User.FindUnique(db.User.ID.Equals(uid)).Update(db.User.Disabled.Set(true)).Exec(ctx)
	if err != nil {
		return err
	}

	// delete all sessions
	sessions.Sessions.Delete(uid)
	return nil
}

func (c UserController) ReenableUser(
	ctx context.Context,
	uid string,
) error {
	_, err := c.Prisma.Client.User.FindUnique(db.User.ID.Equals(uid)).Exec(ctx)
	if err != nil {
		return err
	}

	// now we delete
	_, err = c.Prisma.Client.User.FindUnique(db.User.ID.Equals(uid)).Update(db.User.Disabled.Set(false)).Exec(ctx)
	if err != nil {
		return err
	}

	// delete all sessions
	sessions.Sessions.Delete(uid)
	return nil
}
