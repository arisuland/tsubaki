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
	"arisu.land/tsubaki/pkg"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type response struct {
	Message string `json:"message"`
}

func (m SessionManager) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// If there is no Authorization header, skip it
		if req.Header.Get("Authorization") == "" {
			next.ServeHTTP(w, req)
			return
		}

		auth := req.Header.Get("Authorization")

		// access token! Let's check if it is valid!
		if strings.HasPrefix(auth, "Bearer") {
			// Coming soon!
			next.ServeHTTP(w, req)
		} else if strings.HasPrefix(auth, "Session") {
			// remove any spaces -> trim "Session" off the auth string
			token := strings.Trim(strings.TrimPrefix(auth, "Session"), " ")
			decoded, err := pkg.DecodeToken(token)

			if err != nil {
				logrus.Errorf("Unable to decode session token '%s': %v", token, err)
				w.WriteHeader(400)
				_ = json.NewEncoder(w).Encode(&response{
					Message: fmt.Sprintf("Invalid token: %s", token),
				})

				return
			}

			// validate token
			validated, err := pkg.ValidateToken(token)
			if err != nil {
				logrus.Errorf("Unable to validate token '%s': %v", token, err)
				w.WriteHeader(400)
				_ = json.NewEncoder(w).Encode(&response{
					Message: fmt.Sprintf("Unable to validate token %s", token),
				})

				return
			}

			if !validated {
				logrus.Errorf("Unable to validate token '%s'", token)
				w.WriteHeader(400)
				_ = json.NewEncoder(w).Encode(&response{
					Message: fmt.Sprintf("Unable to validate token %s", token),
				})

				return
			}

			// get user id from MapClaims
			uid, ok := decoded["user_id"].(string)
			if !ok {
				w.WriteHeader(500)
				_ = json.NewEncoder(w).Encode(&response{
					Message: "Unable to cast `user_id` ~> string.",
				})

				return
			}

			ctx := context.WithValue(req.Context(), "userId", uid)
			req = req.WithContext(ctx)
			next.ServeHTTP(w, req)
		} else {
			w.WriteHeader(406)
			_ = json.NewEncoder(w).Encode(&response{
				Message: "Missing `Bearer` or `Session` prefix.",
			})
		}
	})
}
