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

package sessions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"arisu.land/tsubaki/pkg"
	"github.com/sirupsen/logrus"
)

func (m SessionManager) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// If there is no Authorization header, skip it
		if req.Header.Get("Authorization") == "" {
			next.ServeHTTP(w, req)
			return
		}

		auth := req.Header.Get("Authorization")

		// If we have "Basic," we can skip this (since some routes can
		// have this if `config.username` and `config.password` are
		// enabled)
		if strings.HasPrefix(auth, "Basic") {
			next.ServeHTTP(w, req)
			return
		}

		// access token! Let's check if it is valid!
		if strings.HasPrefix(auth, "Bearer") {
			// Coming soon!
			next.ServeHTTP(w, req)
		} else if strings.HasPrefix(auth, "Session") {
			// remove any spaces -> trim "Session" off the auth string
			token := strings.Trim(strings.TrimPrefix(auth, "Session"), " ")
			decoded, err := pkg.DecodeToken(token)

			if err != nil {
				logrus.Errorf("Unable to decode session token: %v", err)
				w.WriteHeader(400)
				_ = json.NewEncoder(w).Encode(&errorResponse{
					Message: fmt.Sprintf("Invalid token: %s", token),
				})

				return
			}

			// validate token
			validated, err := pkg.ValidateToken(token)
			if err != nil {
				logrus.Errorf("Unable to validate token: %v", err)
				w.WriteHeader(400)
				_ = json.NewEncoder(w).Encode(&errorResponse{
					Message: fmt.Sprintf("Unable to validate token %s", token),
				})

				return
			}

			if !validated {
				logrus.Errorf("Unable to validate session token")
				w.WriteHeader(400)
				_ = json.NewEncoder(w).Encode(&errorResponse{
					Message: fmt.Sprintf("Unable to validate token %s", token),
				})

				return
			}

			// get user id from MapClaims
			uid, ok := decoded["user_id"].(string)
			if !ok {
				w.WriteHeader(500)
				_ = json.NewEncoder(w).Encode(&errorResponse{
					Message: "Unable to cast `user_id` ~> string.",
				})

				return
			}

			ctx := context.WithValue(req.Context(), "userId", uid)
			req = req.WithContext(ctx)
			next.ServeHTTP(w, req)
		} else {
			w.WriteHeader(406)
			_ = json.NewEncoder(w).Encode(&errorResponse{
				Message: "Missing `Bearer` or `Session` prefix.",
			})
		}
	})
}
