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

package routers

import (
	"arisu.land/tsubaki/pkg/infra"
	"arisu.land/tsubaki/pkg/util"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type GenericResponse struct {
	Hello    string `json:"hello"`
	DocsUrl  string `json:"docs_url"`
	SiteIcon string `json:"site_icon"`
	SiteName string `json:"site_name"`
}

func NewMainRouter(container *infra.Container) chi.Router {
	router := chi.NewRouter()
	router.Get("/", func(w http.ResponseWriter, req *http.Request) {
		util.WriteJson(w, 200, GenericResponse{
			Hello:    "world",
			DocsUrl:  "https://docs.arisu.land/graphql",
			SiteName: container.Config.SiteName,
			SiteIcon: container.Config.SiteIcon,
		})
	})

	return router
}
