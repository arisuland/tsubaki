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

package main

import (
	"context"
	"os"

	"arisu.land/tsubaki/internal"
	"arisu.land/tsubaki/prisma/db"
	"arisu.land/tsubaki/util"
	"github.com/bwmarrin/snowflake"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&internal.Formatter{})
	logrus.SetLevel(logrus.DebugLevel)

	if _, err := os.Stat("./.env"); !os.IsNotExist(err) {
		err := godotenv.Load("./.env")
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	logrus.Info("🌱 Seeding database...")

	prisma := db.NewClient()
	if err := prisma.Connect(); err != nil {
		logrus.Fatalf("Unable to connect to Prisma: %v", err)
	}

	snowflake, err := snowflake.NewNode(1)
	if err != nil {
		logrus.Fatalf("Unable to create snowflake node: %v", err)
	}

	logrus.Info("Connected to Prisma!")
	id := snowflake.Generate().String()

	pass, err := util.GeneratePassword("admin")
	if err != nil {
		logrus.Fatalf("Unable to generate admin password: %v", err)
	}

	_, err = prisma.User.CreateOne(
		db.User.Username.Set("admin"),
		db.User.Password.Set(pass),
		db.User.Email.Set("admin@arisu.land"),
		db.User.ID.Set(id),
	).Exec(context.TODO())

	if err != nil {
		logrus.Fatalf("Unable to generate admin user: %v", err)
	}

	logrus.Info("🌱 Seeded database successfully!")
}
