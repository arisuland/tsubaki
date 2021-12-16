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
	"arisu.land/tsubaki/managers"
	"arisu.land/tsubaki/prisma/db"
	"arisu.land/tsubaki/util"
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
