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

package managers

import (
	"arisu.land/tsubaki/prisma/db"
	"context"
)

// Prisma is a struct that manages the Prisma client
type Prisma struct {
	// Client is the db.PrismaClient available. This is `nil` until you call
	// the Connect method.
	Client *db.PrismaClient
}

// NewPrisma creates a new Prisma instance.
func NewPrisma() Prisma {
	return Prisma{
		Client: db.NewClient(),
	}
}

// Connect is a method to connect this Prisma instance to the world!
func (p Prisma) Connect() error {
	log.Info(context.Background(), "Connecting to PostgreSQL...")

	err := p.Client.Connect()
	if err != nil {
		return err
	}

	log.Info(context.Background(), "Connected successfully!")
	return nil
}

// Close is a method to close this Prisma instance.
func (p Prisma) Close() error {
	log.Warn(context.Background(), "Closing off connection...")

	if err := p.Client.Disconnect(); err != nil {
		return err
	}

	log.Warn(context.Background(), "Closed off the connection!")
	return nil
}