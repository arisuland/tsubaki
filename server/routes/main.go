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

package routes

import (
	"arisu.land/tsubaki/pkg"
	"arisu.land/tsubaki/server/middleware"
	"arisu.land/tsubaki/util"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type mainResponse struct {
	Message string `json:"message"`
	DocsURI string `json:"docs_url"`
}

func NewMainRouter(container *pkg.Container) chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.BasicAuth(container.Config))
	router.Get("/", func(w http.ResponseWriter, req *http.Request) {
		util.WriteJson(w, 200, mainResponse{
			Message: "hello world!",
			DocsURI: "https://docs.arisu.land",
		})
	})

	return router
}
