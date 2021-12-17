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

package managers

import (
	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"os"
	"time"
)

// Prometheus is the exporter for collecting metrics from Tsubaki.
type Prometheus struct {
	// logger is a private method of the slog.Logger to use.
	logger slog.Logger
}

var (
	RequestMetric = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "tsubaki_api_requests",
		Help: "How many requests have been processed, partitioned by HTTP method and URI path.",
	}, []string{"method", "path"})

	RequestLatencyMetric = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "tsubaki_request_latency",
		Help: "Returns the latency of all requests, excluding POST /graphql",
	}, []string{"method", "path"})

	GQLLatencyMetric = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "tsubaki_gql_latency",
		Help: "Returns the latency of every GraphQL execution partitioned by its operation type.",
	}, []string{"operation"})

	UsersCountMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "tsubaki_users_count",
		Help: "Returns how many registered users are in the database.",
	})
)

// NewPrometheus creates a new singleton of a Prometheus instance.
func NewPrometheus() Prometheus {
	return Prometheus{
		logger: slog.Make(sloghuman.Sink(os.Stdout)),
	}
}

// Register registers all counters and histograms
func (prom Prometheus) Register() {
	prom.logger.Info(context.Background(), "Registering metrics...")
	prometheus.MustRegister(RequestLatencyMetric, RequestMetric, GQLLatencyMetric, UsersCountMetric)

	prom.logger.Info(context.Background(), "Registered all metrics!")
}

// Middleware is the chi middleware to observe request latency
func (prom Prometheus) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, req.ProtoMajor)
		next.ServeHTTP(ww, req)

		RequestMetric.WithLabelValues(req.Method, req.URL.Path)
		RequestLatencyMetric.WithLabelValues(req.Method, req.URL.Path).Observe(float64(time.Since(start).Nanoseconds() / 1000000))
	})
}
