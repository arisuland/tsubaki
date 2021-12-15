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

package routers

import (
	"arisu.land/tsubaki/graphql"
	"arisu.land/tsubaki/infra"
	"arisu.land/tsubaki/middleware"
	"arisu.land/tsubaki/util"
	"fmt"
	"github.com/go-chi/chi/v5"
	"html/template"
	"net/http"
)

func NewGraphQLRouter(container *infra.Container, gql *graphql.Manager) chi.Router {
	router := chi.NewRouter()

	router.Use(middleware.LogMiddleware)
	router.Post("/", gql.ServeHTTP)
	router.Get("/", func(w http.ResponseWriter, req *http.Request) {
		if container.Config.Environment == "development" {
			t := template.New("graphql-playground")
			t, err := t.Parse(util.PlaygroundTemplate)
			if err != nil {
				http.Error(w, err.Error(), 500)
			}

			data := util.PlaygroundTemplateData{
				Endpoint: fmt.Sprintf("http://localhost:%d/graphql", container.Config.Port),
			}

			if err := t.ExecuteTemplate(w, "index", data); err != nil {
				http.Error(w, err.Error(), 500)
			}

			return
		}

		util.WriteJson(w, 405, struct{ Message string }{
			Message: "You can only use the GraphQL API via POST /graphql only.",
		})
	})

	return router
}
