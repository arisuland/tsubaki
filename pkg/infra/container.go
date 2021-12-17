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
	"arisu.land/tsubaki/pkg/kafka"
	"arisu.land/tsubaki/pkg/managers"
	"arisu.land/tsubaki/pkg/storage"
	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
	"context"
	"errors"
	"golang.org/x/xerrors"
	"os"
)

var log = slog.Make(sloghuman.Sink(os.Stdout))

// Container is the main core part of Tsubaki. This is a containerized
// struct of all the core components Tsubaki needs to use.
type Container struct {
	// Prometheus is the metrics manager that is used to measure metrics.
	Prometheus managers.Prometheus

	// Database is the database connection using Prisma.
	Database managers.Prisma

	// Snowflake is the snowflake generator to use to generate unique IDs.
	Snowflake *managers.SnowflakeManager

	// Storage is the base storage manager to use
	Storage storage.BaseStorageProvider

	// Sentry is the sentry manager for handling error reports to Sentry
	Sentry managers.SentryManager

	// Config is the configuration used to configure this Tsubaki instance
	Config *managers.Config

	// Redis the Redis manager to use for caching user sessions.
	Redis *managers.RedisManager

	// Kafka returns the producer we need, this can be `nil` if the config
	// option is not defined.
	Kafka *kafka.Producer
}

// NewContainer creates a new Container instance.
func NewContainer() (*Container, error) {
	log.Info(context.Background(), "Initializing container...")

	// Load configuration
	config, err := managers.NewConfig()
	if err != nil {
		return nil, err
	}

	// Create Prometheus instance
	prom := managers.NewPrometheus()
	prom.Register()

	// Create the Prisma client and connect
	prisma := managers.NewPrisma()
	if err = prisma.Connect(); err != nil {
		return nil, err
	}

	// Set user count
	users, err := prisma.Client.User.FindMany().Exec(context.TODO())
	if err != nil {
		return nil, err
	}

	managers.UsersCountMetric.Set(float64(len(users)))

	// Create the Redis connection
	redis := managers.NewRedisClient(config.Redis)
	if err := redis.Connect(); err != nil {
		return nil, err
	}

	// Create the snowflake manager
	snowflake, err := managers.NewSnowflakeManager()
	if err != nil {
		return nil, err
	}

	// Create the storage providers
	var storageProvider storage.BaseStorageProvider
	if config.Storage.Filesystem != nil {
		storageProvider := storage.NewFilesystemStorageProvider(*config.Storage.Filesystem)
		if err := storageProvider.Init(); err != nil {
			return nil, err
		}
	} else if config.Storage.S3 != nil {
		storageProvider := storage.NewS3StorageProvider(config.Storage.S3)
		if err := storageProvider.Init(); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("missing storage provider to use")
	}

	// Create the Kafka producer (if config is enabled)
	var producer *kafka.Producer
	if config.Kafka != nil {
		producer = kafka.NewProducer(*config.Kafka)
	}

	// Create Sentry client
	sentry, err := managers.NewSentryManager(config)
	if err != nil {
		log.Error(context.Background(), "Unable to initialize Sentry client, will be noop.", slog.F("error", xerrors.Errorf("%v", err)))
	}

	return &Container{
		Prometheus: prom,
		Snowflake:  snowflake,
		Database:   prisma,
		Storage:    storageProvider,
		Sentry:     sentry,
		Config:     config,
		Redis:      redis,
		Kafka:      producer,
	}, nil
}

// Close closes off all components and destroys data.
func (c *Container) Close() error {
	// Close off Redis
	log.Warn(context.Background(), "Closing Redis connection...")
	if err := c.Redis.Connection.Close(); err != nil {
		return err
	}

	// Close off Prisma
	log.Warn(context.Background(), "Closing PostgreSQL connection...")
	if err := c.Database.Close(); err != nil {
		return err
	}

	// Close off Kafka (if we are connected)
	if c.Kafka != nil {
		log.Warn(context.Background(), "Closing off Kafka producer...")
		if err := c.Kafka.Writer.Close(); err != nil {
			return err
		}
	}

	log.Warn(context.Background(), "Services have been closed.")
	return nil
}
