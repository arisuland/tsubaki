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
	"arisu.land/tsubaki/controllers"
	"arisu.land/tsubaki/infra"
)

// Resolver represents the main resolver tree to use in our GraphQL schema.
type Resolver struct {
	// Container is the main container that holds everything in.
	Container *infra.Container

	// Users is the controller for handling user metadata.
	Users controllers.UserController
}

// NewResolver creates a new Resolver instance.
func NewResolver(container *infra.Container) *Resolver {
	return &Resolver{
		Container: container,
		Users:     controllers.NewUserController(container.Database, container.Snowflake),
	}
}
