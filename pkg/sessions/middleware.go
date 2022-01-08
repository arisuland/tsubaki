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

import "net/http"

func (m SessionManager) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

	})
}

/*
func (m Manager) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// if there is no authorization header, let's just skip it.
		if req.Header.Get("Authorization") == "" {
			next.ServeHTTP(w, req)
			return
		}

		auth := req.Header.Get("Authorization")
		if strings.HasPrefix(auth, "Bearer") {
			// TODO: access tokens
			next.ServeHTTP(w, req)
		} else if strings.HasPrefix(auth, "Session") {
			// remove any spaces -> trim "Session" off the string
			token := strings.Trim(strings.TrimPrefix(req.Header.Get("Authorization"), "Session"), " ")
			decoded, err := DecodeToken(token)
			if err != nil {
				logrus.Errorf("Unable to decode token '%s': %v", token, err)
				w.WriteHeader(406)
				_ = json.NewEncoder(w).Encode(&response{
					Message: fmt.Sprintf("Invalid token: %s", token),
				})

				return
			}

			// validate token
			validated, err := ValidateToken(token)
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
			uid, ok := decoded["userId"].(string)
			if !ok {
				w.WriteHeader(500)
				_ = json.NewEncoder(w).Encode(&response{
					Message: "Unable to cast `uid` ~> string.",
				})

				return
			}

			// add it to the request context
			user, err := m.db.Client.User.FindUnique(db.User.ID.Equals(uid)).Exec(context.TODO())
			if err != nil {
				w.WriteHeader(500)
				_ = json.NewEncoder(w).Encode(&response{
					Message: "Unable to retrieve user from database.",
				})

				return
			}

			// blep!
			ctx := context.WithValue(req.Context(), "user_id", user.ID)
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
*/
