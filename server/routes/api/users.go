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
	"arisu.land/tsubaki/pkg/result"
	"arisu.land/tsubaki/util"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type userResponse struct {
	Message string `json:"message"`
}

func newUserApiRouter(controller controllers.Controller) chi.Router {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		util.WriteJson(w, 200, userResponse{
			Message: "Welcome to the Users API! Read more here: https://docs.arisu.land/api/users",
		})
	})

	r.Get("/@me", func(w http.ResponseWriter, req *http.Request) {
		// Check if we have the `Authorization` header
		if req.Header.Get("Authorization") == "" {
			util.WriteJson(w, 401, result.Err(401, "MISSING_AUTH_TOKEN", "Missing user token in request!"))
			return
		}

		// Check if we have the user token available
		uid := req.Context().Value("userId")
		if uid == nil {
			util.WriteJson(w, 403, result.Err(403, "UNKNOWN_USER_ID", "Unable to determine the user ID."))
			return
		}

		res := controller.Users.Get(uid.(string))
		util.WriteJson(w, res.StatusCode, res)
	})

	r.Get("/{id}", func(w http.ResponseWriter, req *http.Request) {
		res := controller.Users.Get(chi.URLParam(req, "id"))
		util.WriteJson(w, res.StatusCode, res)
	})

	//r.Get("/{id}/projects", func(w http.ResponseWriter, req *http.Request) {
	//
	//})

	r.Post("/", func(w http.ResponseWriter, req *http.Request) {
		statusCode, data, err := util.GetJsonBody(req)
		if err != nil {
			util.WriteJson(w, statusCode, result.Err(statusCode, "INVALID_BODY_STRUCTURE", err.Error()))
			return
		}

		username, ok := data["username"].(string)
		if !ok {
			util.WriteJson(w, 406, result.Err(406, "MISSING_USERNAME", "Missing `username` field in body or `username` was not a valid string."))
			return
		}

		email, ok := data["email"].(string)
		if !ok {
			util.WriteJson(w, 406, result.Err(406, "MISSING_USERNAME", "Missing `email` field in body or `email` was not a valid string."))
			return
		}

		password, ok := data["password"].(string)
		if !ok {
			util.WriteJson(w, 406, result.Err(406, "MISSING_USERNAME", "Missing `password` field in body or `password` was not a valid string."))
			return
		}

		res := controller.Users.Create(username, password, email)
		util.WriteJson(w, res.StatusCode, res)
	})

	return r
}
