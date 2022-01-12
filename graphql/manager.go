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

package graphql

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"arisu.land/tsubaki/graphql/resolvers"
	"github.com/graph-gophers/graphql-go"
	"github.com/sirupsen/logrus"
)

// Manager is the main manager for executing GraphQL queries/mutations/subscriptions.
type Manager struct {
	// Schema is the GraphQL schema generated from the codegen binary.
	Schema *graphql.Schema
}

// RequestBody is the request body when requesting from `POST /graphql`
type RequestBody struct {
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
	Query         string                 `json:"query"`
}

func NewGraphQLManager() (*Manager, error) {
	logrus.Info("Building GraphQL schema...")
	contents, err := ioutil.ReadFile("./schema.gql")
	if err != nil {
		return nil, err
	}

	opts := []graphql.SchemaOpt{graphql.UseFieldResolvers()}
	schema := graphql.MustParseSchema(string(contents), resolvers.NewResolver(), opts...)

	logrus.Info("Successfully built schema.")
	return &Manager{
		Schema: schema,
	}, nil
}

// ServeHTTP is the middleware to host the /graphql endpoint for this Manager.
func (m *Manager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var params RequestBody
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		logrus.Errorf("Unable to unmarshal payload:\n%v", err)
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
