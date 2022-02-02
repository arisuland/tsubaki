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
	"arisu.land/tsubaki/internal"
	"arisu.land/tsubaki/pkg"
	"arisu.land/tsubaki/server/middleware"
	"arisu.land/tsubaki/util"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type versionResponse struct {
	CommitSHA string `json:"commit_sha"`
	BuildDate string `json:"build_date"`
	Version   string `json:"version"`
}

func NewVersionRouter(container *pkg.Container) chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.BasicAuth(container.Config))
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		util.WriteJson(w, 200, versionResponse{
			CommitSHA: internal.CommitSHA,
			BuildDate: internal.BuildDate,
			Version:   internal.Version,
		})
	})

	return router
}
