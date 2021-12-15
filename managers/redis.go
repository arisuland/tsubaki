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
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"time"
)

// RedisManager is a manager to handle Redis connections.
type RedisManager struct {
	// Connection is the Redis client used when using the Connect method.
	Connection *redis.Client

	// config represents the Redis configuration.
	config RedisConfig

	// logger is the slog.Logger for this RedisManager instance.
	logger slog.Logger
}

// RedisConfig represents the configuration for using Redis as a cache.
// Sentinel and Standalone are supported!
//
// Prefix: TSUBAKI_REDIS
type RedisConfig struct {
	// Sets a list of sentinel servers to use when using Redis Sentinel
	// instead of Redis Standalone. The `master` key is required if this
	// is defined. This returns a List of `host:port` strings. If you're
	// using an environment variable to set this, split it with `,` so it can be registered properly!
	//
	// Default: nil | Variable: TSUBAKI_REDIS_SENTINEL_SERVERS
	Sentinels []string `yaml:"sentinels"`

	// If `requirepass` is set on your Redis server config, this property will authenticate
	// Tsubaki once the connection is being dealt with.
	//
	// Default: nil | Variable: TSUBAKI_REDIS_PASSWORD
	Password *string `yaml:"password"`

	// Returns the master name for connecting to any Redis sentinel servers.
	//
	// Default: nil | Variable: TSUBAKI_REDIS_SENTINEL_MASTER
	MasterName *string `yaml:"master"`

	// Returns the database index to use so you don't clutter your server
	// with Tsubaki-set configs.
	//
	// Default: 5 | Variable: TSUBAKI_REDIS_DATABASE_INDEX
	DbIndex int `yaml:"index"`

	// Returns the host for connecting to Redis.
	//
	// Default: "localhost" | Variable: TSUBAKI_REDIS_HOST
	Host string `yaml:"host"`

	// Returns the port to use when connecting to Redis.
	//
	// Default: 6379 | Variable: TSUBAKI_REDIS_PORT
	Port int `yaml:"port"`
}

// NewRedisClient initialises a new RedisManager instance.
func NewRedisClient(config RedisConfig) *RedisManager {
	return &RedisManager{
		Connection: nil,
		config:     config,
		logger:     slog.Make(sloghuman.Sink(os.Stdout)),
	}
}

// Connect creates a new Connection towards Redis.
func (m *RedisManager) Connect() error {
	m.logger.Info(context.TODO(), "Connecting to Redis...")

	var password string
	if m.config.Password == nil {
		password = ""
	} else {
		password = *m.config.Password
	}

	if len(m.config.Sentinels) > 0 {
		var masterName string
		if m.config.MasterName == nil {
			return errors.New("config option 'redis.master_name' needs to be defined to use a sentinel connection")
		} else {
			masterName = *m.config.MasterName
		}

		m.Connection = redis.NewFailoverClient(&redis.FailoverOptions{
			SentinelAddrs: m.config.Sentinels,
			MasterName:    masterName,
			Password:      password,
			DB:            m.config.DbIndex,
			DialTimeout:   10 * time.Second,
			ReadTimeout:   15 * time.Second,
			WriteTimeout:  15 * time.Second,
		})
	} else {
		m.Connection = redis.NewClient(&redis.Options{
			Addr:         fmt.Sprintf("%s:%d", m.config.Host, m.config.Port),
			Password:     password,
			DB:           m.config.DbIndex,
			DialTimeout:  10 * time.Second,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
		})
	}

	if err := m.Connection.Ping(context.TODO()).Err(); err != nil {
		return err
	}

	return nil
}
