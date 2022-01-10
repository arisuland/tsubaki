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
	"arisu.land/tsubaki/util"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func newStorageRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		util.WriteJson(w, 200, struct {
			Message string `json:"message"`
		}{
			Message: "Welcome to the not-so secret internal API. This is not documented, beware!",
		})
	})

	return r
}
