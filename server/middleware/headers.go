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
	"arisu.land/tsubaki/internal"
	"fmt"
	"net/http"
)

func Headers(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Sets some basic headers, if available
		req.Header.Set("X-Powered-By", fmt.Sprintf("Arisu/Tsubaki v%s", internal.Version))
		req.Header.Set("Cache-Control", "public, max-age=7776000")

		// Add some security headers (for reasons:tm:)
		// TODO: should this be a config variable? (config.server.security_headers)
		req.Header.Set("X-Frame-Options", "deny")
		req.Header.Set("X-Content-Type-Options", "nosniff")
		req.Header.Set("X-XSS-Protection", "1; mode=block")

		next.ServeHTTP(w, req)
	})
}
