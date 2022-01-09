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
	"arisu.land/tsubaki/pkg/ratelimit"
	"arisu.land/tsubaki/util"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		s := time.Now()
		ww := middleware.NewWrapResponseWriter(w, req.ProtoMajor)
		next.ServeHTTP(ww, req)

		// skip on `/graphql` requests since they can get spammy :<
		if req.URL.Path == "/graphql" {
			return
		}

		code := util.GetStatusCode(ww.Status())
		logrus.Infof("[%s] %s %s (%s) => %d %s (%d bytes) [%s]",
			ratelimit.RealIP(req),
			req.Method,
			req.URL.Path,
			req.Proto,
			ww.Status(),
			code,
			ww.BytesWritten(),
			time.Since(s).String(),
		)

		internal.RequestMetric.WithLabelValues(req.Method, req.URL.Path).Inc()
		internal.RequestLatencyMetric.WithLabelValues(req.Method, req.URL.Path).Observe(float64(time.Since(s).Nanoseconds() / 1000000))
	})
}
