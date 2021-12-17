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
	"arisu.land/tsubaki/pkg/managers"
	"arisu.land/tsubaki/pkg/util"
	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"os"
	"time"
)

var log = slog.Make(sloghuman.Sink(os.Stdout))

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, req.ProtoMajor)
		next.ServeHTTP(ww, req)

		statusCode := util.GetStatusCode(ww.Status())
		log.Info(req.Context(), fmt.Sprintf(
			"[%s] %s %s (%s) => %d %s (%d bytes written) [%s]",
			req.RemoteAddr,
			req.Method,
			req.URL.Path,
			req.Proto,
			ww.Status(),
			statusCode,
			ww.BytesWritten(),
			time.Since(start).String(),
		))

		managers.RequestMetric.WithLabelValues(req.Method, req.URL.Path).Inc()
		managers.RequestLatencyMetric.WithLabelValues(req.Method, req.URL.Path).Observe(float64(time.Since(start).Nanoseconds() / 1000000))
	})
}
