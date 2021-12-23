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
	"arisu.land/tsubaki/pkg/infra"
	"context"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type ErrorHandlerMiddleware struct {
	container *infra.Container
}

func NewErrorHandler(container *infra.Container) ErrorHandlerMiddleware {
	return ErrorHandlerMiddleware{container}
}

func (middle ErrorHandlerMiddleware) Serve(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Do nothing and serve the next middleware in the tree.
		if middle.container.Sentry.Client == nil {
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
			sentry.TransactionName(fmt.Sprintf("req %s %s", req.Method, req.URL.Path)),
			sentry.ContinueFromRequest(req),
		)

		defer span.Finish()
		req = req.WithContext(span.Context())
		hub.Scope().SetRequest(req)
		defer middle.recover(hub, req)

		next.ServeHTTP(w, req)
	})
}

func (middle ErrorHandlerMiddleware) recover(hub *sentry.Hub, req *http.Request) {
	if err := recover(); err != nil {
		eventId := hub.RecoverWithContext(context.WithValue(req.Context(), sentry.RequestContextKey, req), err)
		if eventId != nil {
			hub.Flush(1 * time.Second)
		}

		// print the panic if there was one
		logrus.Fatalf("Received panic on %s %s:\n%v", req.Method, req.URL.Path, err)
	}
}
