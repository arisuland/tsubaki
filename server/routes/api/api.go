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

package api

import (
	"arisu.land/tsubaki/internal/controllers"
	"arisu.land/tsubaki/util"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func NewApiV1Router(controller controllers.Controller) chi.Router {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		util.WriteJson(w, 200, struct {
			Message        string `json:"message"`
			DocsUri        string `json:"docs_url"`
			DefaultVersion string `json:"default_version"`
			CurrentVersion string `json:"current_version"`
		}{
			Message:        "hello world!",
			DocsUri:        "https://docs.arisu.land",
			DefaultVersion: "v1",
			CurrentVersion: "v1",
		})
	})

	r.Mount("/users", newUserApiRouter(controller))
	r.Mount("/admin", newAdminRouter())
	r.Mount("/login", newLoginApiRouter())
	r.Mount("/search", newSearchApiRouter())
	r.Mount("/storage", newStorageRouter())
	r.Mount("/projects", newProjectsApiRouter())

	return r
}
