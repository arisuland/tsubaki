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

package main

import (
	"arisu.land/tsubaki/graphql/resolvers"
	"arisu.land/tsubaki/pkg/logging"
	"context"
	"fmt"
	"github.com/graph-gophers/graphql-go"
	"github.com/prometheus/common/log"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

func init() {
	logrus.SetFormatter(&logging.Formatter{})
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {
	logrus.Info("Generating documentation from schema...")
	contents, err := ioutil.ReadFile("./schema.gql")
	if err != nil {
		log.Fatal(context.Background(), "Unable to find schema.gql file. You must be in the root directory of the project.")
		os.Exit(1)
	}

	opts := []graphql.SchemaOpt{graphql.UseFieldResolvers()}

	// It's fine if we have the container as `nil`
	// since we are not making any requests, so it's perfectly fine.
	schema := graphql.MustParseSchema(string(contents), resolvers.NewResolver(nil), opts...)
	logrus.Info("Generated schema! Now creating docs.json file...")

	fmt.Println(schema.ASTSchema())
}
