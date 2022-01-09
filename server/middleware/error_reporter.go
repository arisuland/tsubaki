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
	"context"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func ErrorReporter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// if sentry wasn't found, let's just recover panic errors.
		if pkg.GlobalContainer.Sentry == nil {
			defer func() {
				if err := recover(); err != nil {
					logrus.Fatalf("Received panic on route '%s %s': %v", req.Method, req.URL.Path, err)
				}
			}()

			next.ServeHTTP(w, req)
			return
		}

		ctx := req.Context()
		hub := sentry.GetHubFromContext(ctx)
		if hub == nil {
			hub = sentry.CurrentHub().Clone()
			ctx = sentry.SetHubOnContext(ctx, hub)
		}

		span := sentry.StartSpan(ctx, "tsubaki.server",
			sentry.TransactionName(fmt.Sprintf("request %s %s", req.Method, req.URL.Path)),
			sentry.ContinueFromRequest(req),
		)

		defer span.Finish()
		req = req.WithContext(span.Context())
		hub.Scope().SetRequest(req)
		defer func() {
			if err := recover(); err != nil {
				logrus.Fatalf("Received panic on route '%s %s': %v", req.Method, req.URL.Path, err)

				eventId := hub.RecoverWithContext(context.WithValue(req.Context(), sentry.RequestContextKey, req), err)
				if eventId != nil {
					hub.Flush(1 * time.Second)
				}
			}
		}()

		next.ServeHTTP(w, req)
	})
}
