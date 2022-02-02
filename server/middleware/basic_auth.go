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

package middleware

import (
	"arisu.land/tsubaki/pkg"
	"arisu.land/tsubaki/util"
	"crypto/subtle"
	"fmt"
	"net/http"
)

func BasicAuth(config *pkg.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if config.Username != nil && config.Password != nil {
				user, pass, ok := req.BasicAuth()
				if !ok {
					w.Header().Add("WWW-Authenticate", `Basic realm="Arisu/Tsubaki"`)
					util.WriteJson(w, http.StatusUnauthorized, &response{
						Message: "Unable to obtain credentials (server has basic authentication enabled)",
					})

					return
				}

				if user != *config.Username {
					w.Header().Add("WWW-Authenticate", `Basic realm="Arisu/Tsubaki"`)
					util.WriteJson(w, http.StatusForbidden, &response{
						Message: fmt.Sprintf("unknown username: %s", user),
					})

					return
				}

				if subtle.ConstantTimeCompare([]byte(*config.Password), []byte(pass)) != 1 {
					w.Header().Add("WWW-Authenticate", `Basic realm="Arisu/Tsubaki"`)
					util.WriteJson(w, http.StatusUnauthorized, &response{
						Message: "invalid credentials provided",
					})

					return
				}

				next.ServeHTTP(w, req)
			} else {
				next.ServeHTTP(w, req)
			}
		})

		return fn
	}
}
