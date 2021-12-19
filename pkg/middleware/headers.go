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

package middleware

import (
	"arisu.land/tsubaki/pkg"
	"fmt"
	"net/http"
	"os"
)

func Headers(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		headers := w.Header()

		var node = ""
		if os.Getenv("DEDI_NODE") != "" {
			dediNode := os.Getenv("DEDI_NODE")
			node = "/" + dediNode
		}

		headers.Set("Access-Control-Allow-Methods", "GET,POST")
		headers.Set("Access-Control-Allow-Origin", "*")
		headers.Set("X-Powered-By", fmt.Sprintf("Arisu/Tsubaki (+https://github.com/arisuland/tsubaki; v%s)", pkg.Version))
		headers.Set("Server", fmt.Sprintf("Noelware%s", node))

		next.ServeHTTP(w, req)
	})
}
