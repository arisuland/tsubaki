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
	"arisu.land/tsubaki/pkg"
	"arisu.land/tsubaki/pkg/util"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type versionResponse struct {
	CommitHash string `json:"commit_hash"`
	Version    string `json:"version"`
}

func NewVersionRouter() chi.Router {
	router := chi.NewRouter()
	router.Get("/", func(w http.ResponseWriter, req *http.Request) {
		util.WriteJson(w, 200, versionResponse{
			CommitHash: pkg.CommitHash,
			Version:    pkg.Version,
		})
	})

	return router
}
