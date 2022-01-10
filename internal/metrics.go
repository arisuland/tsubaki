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

package internal

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

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

// RegisterMetrics registers all counters and histograms
func RegisterMetrics() {
	logrus.Debug("Creating metrics...")
	prometheus.MustRegister(RequestLatencyMetric, RequestMetric, GQLLatencyMetric, UsersCountMetric)

	logrus.Debug("Metrics have been established.")
}
