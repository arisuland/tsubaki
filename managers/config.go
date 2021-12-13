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
	"arisu.land/tsubaki/kafka"
	"arisu.land/tsubaki/storage"
	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
	"context"
	"errors"
	"flag"
	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

var (
	log        = slog.Make(sloghuman.Sink(os.Stdout))
	configFile = flag.String("c", "./config.yml", "the config file path")

	NotIntError       = errors.New("provided value was not a integer")
	InvalidUriError   = errors.New("invalid uri structure")
	InvalidRangeError = errors.New("invalid range provided")
)

// Config is the struct for the actual `config.yml` file. Any environment variable
// set can be overridden if the value doesn't exist
type Config struct {
	// If registrations should be enabled on the server. If not, the administrators
	// of this Tsubaki instance is required to create your user on the /admin endpoint.
	//
	// Default: true | Env Variable: TSUBAKI_REGISTRATIONS
	Registrations bool `yaml:"registrations"`

	// DSN URI to link up Sentry to Tsubaki. This will output GraphQL, request, and database
	// errors towards Sentry.
	//
	// Default: nil | Env Variable: TSUBAKI_SENTRY_DSN
	SentryDsn *string `yaml:"sentry_dsn"`

	// If Tsubaki should send telemetry events to Arisu, read up on what we do
	// before enabling: https://docs.arisu.land/telemetry
	//
	// Default: false | Env Variable: TSUBAKI_ENABLE_TELEMETRY
	TelemetryEnabled bool `yaml:"telemetry"`

	// Returns the site-name to embed on the navbar of your Fubuki instance.
	//
	// Default: "Arisu" | Env Variable: TSUBAKI_SITE_NAME
	SiteName *string `yaml:"site_name"`

	// Returns the site icon to use when displayed on the website.
	//
	// Default: https://cdn.arisu.land/lotus.png | Env Variable: TSUBAKI_SITE_ICON
	SiteIcon *string `yaml:"site_icon"`

	// Returns a number on how many retries before panicking on an unavailable port
	// to run Tsubaki on. To disable this, use `-1` as the value.
	//
	// Default: 5 | Min: -1 | Max: 15 | Env Variable: TSUBAKI_PORT_RETRY
	PortRetryLimit int8 `yaml:"port_retry"`

	// Uses a host URI when launching the HTTP server. If you wish to keep Tsubaki
	// running internally, you can use `127.0.0.1` instead of the default `0.0.0.0`.
	//
	// Default: "0.0.0.0" | Env Variable: TSUBAKI_HOST or HOST
	Host *string `yaml:"host"`

	// Uses a different port other than 2809. If the port is taken, it will
	// try to find an available one round-robin style and use that one that
	// isn't taken. You can set how many tries using the `port_retry` config
	// variable.
	//
	// Range: 80-65535 on root; 1024-65535 on non-root
	//
	// Default: 28093 | Env Variable: TSUBAKI_PORT or PORT
	Port int16 `yaml:"port"`

	// Returns the configuration for using the filesystem, S3,
	// or Google Cloud Storage to backup your projects.
	//
	// Default: FilesystemConfig | Prefix: TSUBAKI_STORAGE

	// Configuration to setup a Kafka producer to use for message queues.
	// This is required if you're running the GitHub bot.
	//
	// Default: nil | Prefix: TSUBAKI_KAFKA
	Kafka kafka.Config `yaml:"kafka"`

	// Returns the configuration for using the filesystem, S3,
	// or Google Cloud Storage to backup your projects.
	//
	// Default:
	//   fs:
	//    directory: <cwd>/.arisu
	//
	// Prefix: TSUBAKI_STORAGE_*
	Storage storage.Config `yaml:"storage"`

	// Returns the configuration for using Redis to cache user sessions.
	// Sentinel and Standalone are supported.
	//
	// Default:
	//  redis:
	//   host: localhost
	//   port: 6379
	//
	// Prefix: TSUBAKI_REDIS_*
	Redis RedisConfig `yaml:"redis"`
}

// NewConfig loads the configuration with a tuple of (*Config, error).
func NewConfig() (*Config, error) {
	if cfg, err := loadConfig(); err != nil {
		return nil, err
	} else {
		return cfg, nil
	}
}

func loadConfig() (*Config, error) {
	log.Info(context.Background(), "Loading configuration...")
	flag.Parse()

	if configFile == nil {
		log.Warn(context.Background(), "Missing config file flag, checking from environment variables...")
		return loadFromEnvironment()
	}

	contents, err := ioutil.ReadFile(*configFile)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(contents, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func loadFromEnvironment() (*Config, error) {
	log.Info(context.Background(), "Loading configuration from environment variables...")

	// Load from .env if there is any
	if _, err := os.Stat("./.env"); !os.IsNotExist(err) {
		err := godotenv.Load(".env")
		if err != nil {
			// it's fine to panic since there might be something wrong
			// with the .env file.
			return nil, err
		}
	}

	// this is going to get messy real quick...
	envRegistrations := os.Getenv("TSUBAKI_REGISTRATIONS")
	envSentryDsn := os.Getenv("TSUBAKI_SENTRY_DSN")
	//envTelemetryEnabled := os.Getenv("TSUBAKI_ENABLE_TELEMETRY")
	//envSiteName := os.Getenv("TSUBAKI_SITE_NAME")
	//envSiteIcon := os.Getenv("TSUBAKI_SITE_ICON")
	//envPortRetryLimit := os.Getenv("TSUBAKI_PORT_RETRY")
	//envHost := os.Getenv("TSUBAKI_HOST")
	//envPort := os.Getenv("TSUBAKI_PORT")

	// Check if we can cast `registrations` into a bool
	var registrations bool
	if envRegistrations == "" {
		// Registrations are enabled by default if it is not provided
		registrations = true
	} else {
		// We can do "yes"/"no" style booleans.
		if envRegistrations == "yes" || envRegistrations == "true" {
			registrations = true
		}

		if envRegistrations == "no" || envRegistrations == "false" {
			registrations = false
		}
	}

	// If the sentry dsn is specified, check if it is a valid one
	var sentryDsn *string
	if envSentryDsn != "" {
		dsn, err := sentry.NewDsn(envSentryDsn);
		if err != nil {
			return nil, err
		}

		dsnUri := dsn.String()
		sentryDsn = &dsnUri
	} else {
		sentryDsn = nil
	}

	return &Config{
		Registrations: registrations,
		SentryDsn:     sentryDsn,
	}, nil
}
