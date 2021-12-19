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

package graphql

import (
	"arisu.land/tsubaki/graphql/resolvers"
	"arisu.land/tsubaki/pkg/infra"
	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
	"context"
	"encoding/json"
	"github.com/graph-gophers/graphql-go"
	"io/ioutil"
	"net/http"
	"os"
)

// Manager is the main manager for executing GraphQL queries/mutations/subscriptions.
type Manager struct {
	// Container is the container initialized from main.go
	Container *infra.Container

	// Schema is the GraphQL schema generated from the codegen binary.
	Schema *graphql.Schema

	// Logger returns the logger for this Manager.
	Logger slog.Logger
}

// RequestBody is the request body when requesting from `POST /graphql`
type RequestBody struct {
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
	Query         string                 `json:"query"`
}

// NewGraphQLManager creates a new *Manager instance.
func NewGraphQLManager(container *infra.Container) *Manager {
	return &Manager{
		Container: container,
		Logger:    slog.Make(sloghuman.Sink(os.Stdout)),
		Schema:    nil,
	}
}

// GenerateSchema returns a error if it cannot properly generate the schema
// to be executable. Run `./build/codegen/schema` to generate the `schema.gql` file.
func (m *Manager) GenerateSchema() error {
	m.Logger.Info(context.Background(), "Generating GraphQL schema...")
	contents, err := ioutil.ReadFile("./schema.gql")
	if err != nil {
		return err
	}

	opts := []graphql.SchemaOpt{graphql.UseFieldResolvers()}
	content := string(contents)
	schema := graphql.MustParseSchema(content, resolvers.NewResolver(m.Container), opts...)

	m.Logger.Info(context.Background(), "Generated successfully. :3")
	m.Schema = schema
	return nil
}

// ServeHTTP is the middleware to host the /graphql endpoint for this Manager.
func (m *Manager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var params RequestBody
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		m.Logger.Error(r.Context(), "Unable to decode ")

		http.Error(w, err.Error(), 500)
		return
	}

	result := m.Schema.Exec(r.Context(), params.Query, params.OperationName, params.Variables)
	data, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	_, _ = w.Write(data)
}
