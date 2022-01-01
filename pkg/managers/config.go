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

package managers

import (
	"arisu.land/tsubaki/pkg/kafka"
	"arisu.land/tsubaki/pkg/storage"
	"arisu.land/tsubaki/pkg/util"
	"errors"
	"flag"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var (
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

	// Sets the environment of the project. If this is set to `development` (which is the default),
	// it will have a GraphQL playground available for you to test the API!
	//
	// Default: "development" | Env Variable: TSUBAKI_GO_ENV or GO_ENV
	Environment string `yaml:"environment"`

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

	// SecretKeyBase is the string for hashing JWT tokens. If this is not populated,
	// you can run the `tsubaki generate` command OR have the app create it for you
	// but it will not write to the config file to loosen I/O computation.
	//
	// If you're using Kubernetes, it is best assured that you should use this as a Secret
	// object and use the environment variable to load it in.
	//
	// Default: "<auto-gen>" | Env Variable: TSUBAKI_SECRET_KEY_BASE
	SecretKeyBase string `yaml:"secret_key_base"`

	// If this Tsubaki instance should be invite only, in which,
	// if Registrations is disabled, this is also disabled.
	//
	// If Registrations is enabled, the administrators will
	// allow you to join the instance or be denied via email.
	//
	// Default: false | Env Variable: TSUBAKI_INVITE_ONLY
	InviteOnly bool `yaml:"invite_only"`

	// Returns the site-name to embed on the navbar of your Fubuki instance.
	//
	// Default: "Arisu" | Env Variable: TSUBAKI_SITE_NAME
	SiteName string `yaml:"site_name"`

	// Returns the site icon to use when displayed on the website.
	//
	// Default: https://cdn.arisu.land/lotus.png | Env Variable: TSUBAKI_SITE_ICON
	SiteIcon string `yaml:"site_icon"`

	// Returns a number on how many retries before panicking on an unavailable port
	// to run Tsubaki on. To disable this, use `-1` as the value.
	//
	// Default: 5 | Min: -1 | Max: 15 | Env Variable: TSUBAKI_PORT_RETRY
	PortRetryLimit int `yaml:"port_retry"`

	// Uses a host URI when launching the HTTP server. If you wish to keep Tsubaki
	// running internally, you can use `127.0.0.1` instead of the default `0.0.0.0`.
	//
	// Default: "0.0.0.0" | Env Variable: TSUBAKI_HOST or HOST
	Host string `yaml:"host"`

	// Uses a different port other than 2809. If the port is taken, it will
	// try to find an available one round-robin style and use that one that
	// isn't taken. You can set how many tries using the `port_retry` config
	// variable.
	//
	// Range: 80-65535 on root; 1024-65535 on non-root
	//
	// Default: 28093 | Env Variable: TSUBAKI_PORT or PORT
	Port int `yaml:"port"`

	// Returns the configuration for using the filesystem, S3,
	// or Google Cloud Storage to backup your projects.
	//
	// Default: FilesystemConfig | Prefix: TSUBAKI_STORAGE

	// Configuration to setup a Kafka producer to use for message queues.
	// This is required if you're running the GitHub bot.
	//
	// Default: nil | Prefix: TSUBAKI_KAFKA
	Kafka *kafka.Config `yaml:"kafka"`

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
	logrus.Info("Loading configuration...")

	if configFile == nil {
		// Check if the root directory has the file (most likely)
		_, err := os.Stat("./config.yml")
		if !os.IsNotExist(err) {
			logrus.Info("Found configuration in root directory, loading from that...")

			owo := "./config.yml"
			configFile = &owo
		}

		logrus.Warn("Couldn't find `-c` flag or from current directory, loading from environment variables...")
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

	// Validate if `environment` exists
	if config.Environment == "" {
		return nil, errors.New("missing `environment` scope in config.yml")
	}

	// Check if it's not "development" or "production"
	validEnvs := []string{"development", "production"}
	var valid = false

	for _, key := range validEnvs {
		if key == config.Environment {
			valid = true
		}
	}

	if !valid {
		return nil, errors.New(fmt.Sprintf("unknown go env: %s", config.Environment))
	}

	// Check for site icon / site name
	if config.SiteName == "" {
		config.SiteName = "Arisu"
	}

	if config.SiteIcon == "" {
		config.SiteIcon = "https://cdn.arisu.land/lotus.png"
	}

	// Check if the environment variable exists for secret key base
	if os.Getenv("TSUBAKI_SECRET_KEY_BASE") != "" && config.SecretKeyBase == "" {
		config.SecretKeyBase = os.Getenv("TSUBAKI_SECRET_KEY_BASE")
	}

	// Check if it is empty
	if config.SecretKeyBase == "" {
		randomHash := util.GenerateHash(32)
		if randomHash == "" {
			return nil, errors.New("unable to generate secret key base :<")
		}

		logrus.Warnf(
			"It is recommended to store your generated key for JWT authentication. After a restart, JWTs will fail to be verified.\n%s",
			randomHash,
		)

		config.SecretKeyBase = randomHash
	}

	return &config, nil
}

func loadFromEnvironment() (*Config, error) {
	logrus.Info("Loading configuration from environment variables...")

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
	envSentryDsn := os.Getenv("TSUBAKI_SENTRY_DSN")

	// If the sentry dsn is specified, check if it is a valid one
	var sentryDsn *string
	if envSentryDsn != "" {
		dsn, err := sentry.NewDsn(envSentryDsn)
		if err != nil {
			return nil, err
		}

		dsnUri := dsn.String()
		sentryDsn = &dsnUri
	} else {
		sentryDsn = nil
	}

	telemetryEnabled := convertToBool(os.Getenv("TSUBAKI_ENABLE_TELEMETRY"), false)
	if telemetryEnabled {
		logrus.Warn("You have enabled telemetry reports! You have been warned... :ghost: (https://docs.arisu.land/telemetry)")
	}

	portRetryLimit, err := strconv.Atoi(os.Getenv("TSUBAKI_PORT_RETRY_LIMIT"))
	if err != nil {
		return nil, err
	}

	// this is ugly send help
	var port int
	var portEnv string
	if portEnv = os.Getenv("TSUBAKI_PORT"); portEnv == "" {
		if portEnv = os.Getenv("PORT"); portEnv == "" {
			portEnv = "none"
		}
	}

	if portEnv == "none" {
		port = 28093
	} else {
		portInt, err := convertToInt(portEnv, 28093)
		if err != nil {
			return nil, err
		}

		port = portInt
	}

	// this still looks ugly :sob:
	var host string
	var hostEnv string
	if hostEnv = os.Getenv("TSUBAKI_HOST"); hostEnv == "" {
		if hostEnv = os.Getenv("HOST"); hostEnv == "" {
			host = "0.0.0.0"
		}
	}

	if hostEnv == "none" {
		host = "0.0.0.0"
	} else {
		host = fallbackEnvString(hostEnv, "0.0.0.0")
	}

	// Now we need to get Redis, Kafka, and Storage env handlers WAAAAAA
	kafkaConfig := getKafkaConfig()
	redisConfig, err := getRedisConfig()

	if err != nil {
		return nil, err
	}

	storageConfig, err := getStorageConfig()
	if err != nil {
		return nil, err
	}

	logrus.Info("Loaded configuration from environment variables!")
	return &Config{
		TelemetryEnabled: telemetryEnabled,
		PortRetryLimit:   portRetryLimit,
		Registrations:    convertToBool(os.Getenv("TSUBAKI_REGISTRATIONS"), true),
		SentryDsn:        sentryDsn,
		SiteName:         fallbackEnvString(os.Getenv("TSUBAKI_SITE_NAME"), "Arisu"),
		SiteIcon:         fallbackEnvString(os.Getenv("TSUBAKI_SITE_ICON"), "https://cdn.arisu.land/lotus.png"),
		Storage:          storageConfig,
		Kafka:            kafkaConfig,
		Redis:            *redisConfig,
		Port:             port,
		Host:             host,
	}, nil
}

// convertToBool basically converts the envString provided into a boolean.
func convertToBool(envString string, fallback bool) bool {
	var value bool
	if envString == "" {
		value = fallback
	} else {
		if envString == "yes" || envString == "true" {
			value = true
		}

		if envString == "no" || envString == "false" {
			value = false
		}
	}

	return value
}

func convertToInt(envString string, fallback int) (int, error) {
	if envString == "" {
		return fallback, nil
	} else {
		value, err := strconv.Atoi(envString)
		if err != nil {
			return fallback, NotIntError
		}

		return value, nil
	}
}

func fallbackEnvString(envString string, fallback string) string {
	if envString == "" {
		return fallback
	} else {
		return envString
	}
}

func getKafkaConfig() *kafka.Config {
	// Check if keys exist
	enabled := os.Getenv("TSUBAKI_KAFKA_BROKERS") == "" || os.Getenv("TSUBAKI_KAFKA_TOPIC") == ""
	if !enabled {
		return nil
	}

	return &kafka.Config{
		AutoCreateTopics: convertToBool(os.Getenv("TSUBAKI_KAFKA_AUTO_CREATE_TOPICS"), true),
		Brokers:          strings.Split(os.Getenv("TSUBAKI_KAFKA_BROKERS"), ","),
		Topic:            fallbackEnvString(os.Getenv("TSUBAKI_KAFKA_TOPIC"), "arisu:tsubaki"),
	}
}

func getRedisConfig() (*RedisConfig, error) {
	var password *string
	if os.Getenv("TSUBAKI_REDIS_PASSWORD") != "" {
		p := os.Getenv("TSUBAKI_REDIS_PASSWORD")
		password = &p
	}

	var masterName *string
	if os.Getenv("TSUBAKI_REDIS_SENTINEL_MASTER") != "" {
		m := os.Getenv("TSUBAKI_REDIS_SENTINEL_MASTER")
		masterName = &m
	}

	dbIndex, err := convertToInt(os.Getenv("TSUBAKI_REDIS_DATABASE_INDEX"), 5)
	if err != nil {
		return nil, err
	}

	redisPort, err := convertToInt(os.Getenv("TSUBAKI_REDIS_PORT"), 6379)
	if err != nil {
		return nil, err
	}

	return &RedisConfig{
		Sentinels:  strings.Split(os.Getenv("TSUBAKI_REDIS_SENTINEL_SERVERS"), ","),
		Password:   password,
		MasterName: masterName,
		DbIndex:    dbIndex,
		Host:       fallbackEnvString(os.Getenv("TSUBAKI_REDIS_HOST"), "localhost"),
		Port:       redisPort,
	}, nil
}

func getStorageConfig() (storage.Config, error) {
	// Check if the directory env variable exists
	if os.Getenv("TSUBAKI_STORAGE_FS_DIRECTORY") != "" {
		return storage.Config{
			Filesystem: &storage.FilesystemStorageConfig{
				Directory: fallbackEnvString(os.Getenv("TSUBAKI_STORAGE_FS_DIRECTORY"), "./.arisu"),
			},
		}, nil
	}

	// Check if any S3 keys exist
	// aka this is going to be messy real quick
	if checkIfS3EnvExists() {
		secretKey := fallbackEnvString(os.Getenv("TSUBAKI_STORAGE_S3_SECRET_KEY"), "")
		accessKey := fallbackEnvString(os.Getenv("TSUBAKI_STORAGE_S3_ACCESS_KEY"), "")
		endpoint := fallbackEnvString(os.Getenv("TSUBAKI_STORAGE_S3_ENDPOINT"), "")

		return storage.Config{
			S3: &storage.S3StorageConfig{
				SecretKey: &secretKey,
				AccessKey: &accessKey,
				Provider:  storage.FromProvider(fallbackEnvString(os.Getenv("TSUBAKI_STORAGE_S3_SECRET_KEY"), "")),
				Endpoint:  &endpoint,
				Region:    os.Getenv("TSUBAKI_STORAGE_S3_REGION"),
				Bucket:    fallbackEnvString(os.Getenv("TSUBAKI_STORAGE_S3_BUCKET"), "tsubaki"),
			},
		}, nil
	}

	return storage.Config{}, errors.New("missing storage provider, read more here: https://docs.arisu.land/selfhosting/storage")
}

func checkIfS3EnvExists() bool {
	prefix := "TSUBAKI_STORAGE_S3_"
	provider := os.Getenv(prefix + "PROVIDER")
	region := os.Getenv(prefix + "REGION")

	return provider != "" && region != ""
}
