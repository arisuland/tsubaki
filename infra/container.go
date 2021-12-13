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

package infra

import (
	"arisu.land/tsubaki/kafka"
	"arisu.land/tsubaki/managers"
	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
	"context"
	"os"
)

// Container is the main core part of Tsubaki. This is a containerized
// struct of all the core components Tsubaki needs to use.
type Container struct {
	// Prometheus is the metrics manager that is used to measure metrics.
	Prometheus managers.Prometheus

	// Snowflake is the snowflake generator to use to generate unique IDs.
	Snowflake *managers.SnowflakeManager

	// Kafka returns the producer we need, this can be `nil` if the config
	// option is not defined.
	Kafka *kafka.Producer
}

// NewContainer creates a new Container instance.
func NewContainer() (*Container, error) {
	logger := slog.Make(sloghuman.Sink(os.Stdout))
	logger.Info(context.Background(), "Initializing container...")

	// Create Prometheus instance
	prom := managers.NewPrometheus()

	// Create the snowflake manager
	snowflake, err := managers.NewSnowflakeManager();
	if err != nil {
		return nil, err
	}

	return &Container{
		Prometheus: prom,
		Snowflake:  snowflake,
	}, nil
}
