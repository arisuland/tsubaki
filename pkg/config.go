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

package pkg

import (
	"arisu.land/tsubaki/pkg/storage"
	"arisu.land/tsubaki/util"
	"errors"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

// Environment is a type to determine the current Tsubaki environment.
type Environment string

var (
	Development Environment = "development"
	Production  Environment = "production"
)

func (e Environment) String() string {
	switch {
	case e == Development:
		return "development"

	case e == Production:
		return "production"

	default:
		return ""
	}
}

var (
	// NotIntError is a error if the a non-integer was provided.
	NotIntError = errors.New("provided value was not a valid integer")

	// MissingEnvironmentError is a error if the `environment` config option is not
	// available in the config.yml file or by environment variables.
	MissingEnvironmentError = errors.New("missing `environment` config option")

	// InvalidEnvironmentError is a error if the `environment` config option
	// is not a valid environment.
	InvalidEnvironmentError = errors.New("environment option was not `development` or `production`")
)

// Config is the structure of the `config.yml` file. Any environment variable
// can be set and overridden in this document.
type Config struct {
	// If registrations are enabled on the server. If not, the `createUser` GraphQL mutation
	// is disabled and administrations must make your account on the administration dashboard.
	//
	// Default: true | Environment Variable: TSUBAKI_REGISTRATIONS
	Registrations bool `yaml:"registrations"`

	// Returns the current environment for the instance. If this is set to `development`,
	// it will log debug messages and has a GraphQL playground to play around with the API.
	//
	// Default: "development" | Environment Variable: TSUBAKI_GO_ENV or GO_ENV
	Environment Environment `yaml:"environment"`

	// DSN to link up Sentry error handling with Tsubaki. This will output
	// GraphQL, request, and database errors towards Sentry.
	//
	// Default: `nil` | Environment Variable: TSUBAKI_SENTRY_DSN
	SentryDSN *string `yaml:"sentry_dsn,omitempty"`

	// If this instance should be able to send out Telemetry events to
	// our main instance (telemetry.arisu.land) to view what users usually do,
	// what size of users, etc. This will enhance the scale on how much Tsubaki
	// is used to determine what to focus on. This is defaulted to `false` due
	// to people being most likely afraid.
	//
	// Note: This does not send out your system information, more of general information,
	// like if it's running on Docker or Kubernetes, memory / cpu usage values (as in percentage),
	// and how many users are linked into this instance.
	//
	// Default: false | Environment Variable: TSUBAKI_TELEMETRY_ENABLED
	Telemetry bool `yaml:"telemetry"`

	// SecretKeyBase is the string to for hashing JWT tokens. If this is not populated,
	// Tsubaki will automatically generate one for you and will remind you to store this
	// so you can be able to verify old JWT tokens.
	//
	// If you're using our Helm Chart with Kubernetes. This is automatically generated
	// and stored in a Secret object, so you don't need to worry about it.
	//
	// You can also use the `tsubaki generate` command to generate one and stores it
	// in your configuration file.
	//
	// Default: "<auto generated>" | Environment Variable: TSUBAKI_SECRET_KEY_BASE
	SecretKeyBase string `yaml:"secret_key_base"`

	// If this instance should be invite only. Which, you are
	// required to get a email from the server admins that you
	// can join this instance and create projects and what-not.
	//
	// Default: false | Environment Variable: TSUBAKI_INVITE_ONLY
	InviteOnly bool `yaml:"invite_only"`

	// Sets the current username to use for non-JWT requests, example
	// would be a request to `127.0.0.1/` since that is non-authenticated.
	Username *string `yaml:"username"`

	// Sets the current password to use for non-JWT request, example
	// would be a request to `127.0.0.1/` since that is non-authenticated.
	Password *string `yaml:"password"`

	// Uses a host URI when launching the HTTP server. If you wish to keep Tsubaki
	// running internally, you can use `127.0.0.1` instead of the default `0.0.0.0`.
	//
	// Default: "0.0.0.0" | Env Variable: TSUBAKI_HOST or HOST
	Host *string `yaml:"host,omitempty"`

	// Uses a different port other than 2809. If the port is taken, it will
	// try to find an available one round-robin style and use that one that
	// isn't taken. You can set how many tries using the `port_retry` config
	// variable.
	//
	// Range: 80-65535 on root; 1024-65535 on non-root
	//
	// Default: 28093 | Env Variable: TSUBAKI_PORT or PORT
	Port *int `yaml:"port,omitempty"`

	// Returns the configuration for using the filesystem, S3,
	// or Google Cloud Storage to backup your projects.
	//
	// Default:
	//   fs:
	//    directory: <cwd>/.arisu
	//
	// Prefix: TSUBAKI_STORAGE_*
	Storage StorageConfig `yaml:"storage"`

	// Configuration to setup a Kafka producer to use for message queues.
	// This is required if you're running the GitHub bot.
	//
	// Default: nil | Prefix: TSUBAKI_KAFKA
	Kafka *KafkaConfig `yaml:"kafka,omitempty"`

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

	// Returns the configuration to use ElasticSearch as the search
	// engine for the `search` GraphQL query.
	//
	// Default: nil | Prefix: TSUBAKI_ELASTIC_*
	ElasticSearch *ElasticsearchConfig `yaml:"elasticsearch,omitempty"`
}

// KafkaConfig is the configuration options for configuring the Kafka Producer.
// This is required if you're running the GitHub bot.
//
// Default: nil | Prefix: TSUBAKI_KAFKA
type KafkaConfig struct {
	// If the producer should create the topic or not.
	//
	// Default: true | Variable: TSUBAKI_KAFKA_AUTO_CREATE_TOPICS
	AutoCreateTopics bool `yaml:"auto_create_topics"`

	// A list of brokers to connect to. This is a List of `host:port` strings.
	//
	// Default: []string{"localhost:9092"} | Variable: TSUBAKI_KAFKA_BROKERS
	Brokers []string `yaml:"brokers"`

	// Returns the topic to use when sending messages towards the GitHub bot (vice-versa).
	//
	// Warning: This must be the same topic you set from the GitHub bot configuration
	// or it will not receive messages!
	//
	// Default: "arisu:tsubaki" | Variable: TSUBAKI_KAFKA_TOPIC
	Topic string `yaml:"topic"`
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
	Sentinels *[]string `yaml:"sentinels,omitempty"`

	// If `requirepass` is set on your Redis server config, this property will authenticate
	// Tsubaki once the connection is being dealt with.
	//
	// Default: nil | Variable: TSUBAKI_REDIS_PASSWORD
	Password *string `yaml:"password,omitempty"`

	// Returns the master name for connecting to any Redis sentinel servers.
	//
	// Default: nil | Variable: TSUBAKI_REDIS_SENTINEL_MASTER
	MasterName *string `yaml:"master,omitempty"`

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

// StorageConfig is the configuration for storing projects.
//
// Prefix: TSUBAKI_STORAGE
type StorageConfig struct {
	// Configures using the filesystem to host your projects, once the data
	// is removed, Arisu will fix it but cannot restore your projects.
	// If you're using Docker or Kubernetes, it is best assured that you
	// must create a volume so Arisu can interact with it.
	//
	// Aliases: `fs` | Prefix: TSUBAKI_STORAGE_FS_*
	Filesystem *storage.FilesystemStorageConfig `yaml:"fs,omitempty"`

	// Configures using Amazon S3 to host your projects.
	// This is a recommended option to store your projects. :3
	S3 *storage.S3StorageConfig `yaml:"s3,omitempty"`
}

// ElasticsearchConfig is the configuration for using Elasticsearch for
// holding search metadata. This is used to power the `search` GraphQL query.
//
// If the Elasticsearch service isn't available the `search` query will just
// return a error.
type ElasticsearchConfig struct {
	// Password is the password to use to authenticate.
	Password *string `yaml:"password,omitempty"`

	// Username is the username if you have authentication enabled.
	Username *string `yaml:"username,omitempty"`

	// Hosts is the URI endpoints to match your ElasticSearch cluster.
	Hosts []string `yaml:"hosts"`
}

func NewConfig(path string) (*Config, error) {
	logrus.Debugf("Now loading configuration in path %s...", path)
	return loadConfigFromPath(path)
}

func TestConfigFromPath(path string) error {
	return testConfig(path)
}

//////////////////////////////////////////////////////////
///////////////// UTILITY FUNCTIONS :D //////////////////
////////////////////////////////////////////////////////

// checkIfVariableExists returns a bool if `value` is false, checks from the
// system environment variables and see if it exists. Otherwise, if `value` is true,
// then it returns false (since this actually "exists" but h)
func checkIfVariableExists(key string, value bool) bool {
	// If we can't find it in the config.toml file, maybe it is
	// an environment variable?
	if !value {
		if _, ok := os.LookupEnv("TSUBAKI_" + key); !ok {
			return false
		}

		return true
	}

	return false
}

// convertToBool converts a "booleanish" [value] string into the actual raw
// value it should be.
//
// Notes:
//   - It panic's if the [value] was not the following regex:
//        - /(yes|true)$/
//        - /(no|false)$/
func convertToBool(value string, fallback bool) bool {
	if value == "" {
		return fallback
	} else {
		regex := regexp.MustCompile("(yes|true)$")
		if regex.MatchString(value) {
			return true
		}

		noRegex := regexp.MustCompile("(no|false)$")
		if noRegex.MatchString(value) {
			return false
		}

		panic(fmt.Errorf("value %s was not `yes`, `no, `true`, or `false` :(", value))
	}
}

// convertToInt converts the "integer" [value] to an actual int.
func convertToInt(value string, fallback int) (int, error) {
	if value == "" {
		return fallback, nil
	} else {
		v, err := strconv.Atoi(value)
		if err != nil {
			return fallback, NotIntError
		}

		return v, nil
	}
}

// fallbackToString check if [value] is empty, falls back to the [fallback] parameter.
func fallbackToString(value string, fallback string) string {
	if value == "" {
		return fallback
	} else {
		return value
	}
}

func getKafkaConfigFromEnv() *KafkaConfig {
	logrus.Debug("Finding Kafka configuration from system environment variables...")
	disabled := os.Getenv("TSUBAKI_KAFKA_BROKERS") == ""
	if disabled {
		logrus.Debug("Kafka producer will not be initialized due to `TSUBAKI_KAFKA_BROKERS` not existing.")
		return nil
	}

	logrus.Debug("Found Kafka configuration from system environment variables!")
	return &KafkaConfig{
		AutoCreateTopics: convertToBool(os.Getenv("TSUBAKI_KAFKA_AUTO_CREATE_TOPICS"), true),
		Brokers:          strings.Split(os.Getenv("TSUBAKI_KAFKA_BROKERS"), ","),
		Topic:            fallbackToString(os.Getenv("TSUBAKI_KAFKA_TOPIC"), "tsubaki"),
	}
}

func getRedisConfigFromEnv() (*RedisConfig, error) {
	logrus.Debug("Now finding Redis configuration from environment variables...")
	var password *string
	if value, ok := os.LookupEnv("TSUBAKI_REDIS_PASSWORD"); ok {
		password = &value
	}

	// Check if we have a sentinel connection
	// If we need to connect with Redis Sentinel, then
	// `TSUBAKI_REDIS_SENTINEL_SERVERS` and `TSUBAKI_REDIS_SENTINEL_MASTER` to exist.
	if servers, ok := os.LookupEnv("TSUBAKI_REDIS_SENTINEL_SERVERS"); ok {
		logrus.Debug("Found `TSUBAKI_REDIS_SENTINEL_SERVERS` environment variable!")
		if master, ok := os.LookupEnv("TSUBAKI_REDIS_SENTINEL_MASTER"); ok {
			logrus.Debug("We can now create a Sentinel connection since both `TSUBAKI_REDIS_SENTINEL_SERVERS` and `TSUBAKI_REDIS_SENTINEL_MASTER` exist!")
			serverList := strings.Split(servers, ",")

			index, err := convertToInt(os.Getenv("TSUBAKI_REDIS_DB_INDEX"), 8)
			if err != nil {
				return nil, err
			}

			return &RedisConfig{
				Sentinels:  &serverList,
				MasterName: &master,
				DbIndex:    index,

				// we can use these as defaults since
				// the sentinel connection won't use them.
				Host: "localhost",
				Port: 6379,
			}, nil
		}
	}

	logrus.Debug("Determined that we are using a Redis Standalone connection!")

	index, err := convertToInt(os.Getenv("TSUBAKI_REDIS_DB_INDEX"), 8)
	if err != nil {
		return nil, err
	}

	redisPort, err := convertToInt(os.Getenv("TSUBAKI_REDIS_DB_INDEX"), 6379)
	if err != nil {
		return nil, err
	}

	return &RedisConfig{
		Sentinels:  nil,
		Password:   password,
		MasterName: nil,
		DbIndex:    index,
		Host:       fallbackToString(os.Getenv("TSUBAKI_REDIS_HOST"), "localhost"),
		Port:       redisPort,
	}, nil
}

func getStorageConfigFromEnv() StorageConfig {
	if provider, ok := os.LookupEnv("TSUBAKI_STORAGE_PROVIDER"); ok {
		logrus.Debugf("storage: We found a provider with %s!", provider)
		switch provider {
		case "filesystem":
		case "fs":
			{
				logrus.Debug("storage: using filesystem configuration...")
				logrus.Warn("It is not recommended to use this under production, we recommend using the S3 storage provider to have a backup -- https://docs.arisu.land/self-hosting/warnings")

				if actualPath, ok := os.LookupEnv("TSUBAKI_STORAGE_FILESYSTEM_DIRECTORY"); ok {
					logrus.Debugf("Found path %s to use to init provider.", actualPath)
					return StorageConfig{
						Filesystem: &storage.FilesystemStorageConfig{
							Directory: actualPath,
						},
					}
				}

				path := ""
				if runtime.GOOS == "linux" {
					path = "/etc/arisu/tsubaki/storage"
				} else if runtime.GOOS == "windows" {
					appdata := os.Getenv("APPDATA")
					path = appdata + "\\Tsubaki\\Storage"
				} else {
					// TODO: macos - try to figure out what to use
					panic(fmt.Errorf("unsupported runtime: %s", runtime.GOOS))
				}

				logrus.Warn("It is recommended to set your own path if using the filesystem configuration with `TSUBAKI_STORAGE_FILESYSTEM_DIRECTORY` environment variable with a valid path, so you don't run into permission errors!")
				return StorageConfig{
					Filesystem: &storage.FilesystemStorageConfig{
						Directory: path,
					},
				}
			}

		case "s3":
			{
				logrus.Debug("storage: using s3 storage configuration...")

				secretKey := fallbackToString(os.Getenv("TSUBAKI_STORAGE_S3_SECRET_KEY"), "")
				accessKey := fallbackToString(os.Getenv("TSUBAKI_STORAGE_S3_ACCESS_KEY"), "")
				endpoint := fallbackToString(os.Getenv("TSUBAKI_STORAGE_S3_ENDPOINT"), "")

				if secretKey == "" || accessKey == "" {
					panic(fmt.Errorf("missing secret and access keys for s3 authentication"))
				}

				provider := storage.FromProvider(fallbackToString(os.Getenv("TSUBAKI_STORAGE_S3_PROVIDER"), "amazon"))
				if provider == storage.Empty {
					panic(fmt.Errorf("unable to determine s3 provider with %s :(", fallbackToString(os.Getenv("TSUBAKI_STORAGE_S3_PROVIDER"), "amazon")))
				}

				return StorageConfig{
					S3: &storage.S3StorageConfig{
						SecretKey: &secretKey,
						AccessKey: &accessKey,
						Provider:  storage.FromProvider(fallbackToString(os.Getenv("TSUBAKI_STORAGE_S3_PROVIDER"), "amazon")),
						Endpoint:  &endpoint,
						Region:    fallbackToString(os.Getenv("TSUBAKI_STORAGE_S3_REGION"), "us-east1"),
						Bucket:    fallbackToString(os.Getenv("TSUBAKI_STORAGE_S3_BUCKET"), "tsubaki"),
					},
				}
			}
		}
	} else {
		logrus.Warnf("Missing `TSUBAKI_STORAGE_PROVIDER` system environment variable!")
		if actualPath, ok := os.LookupEnv("TSUBAKI_STORAGE_FILESYSTEM_DIRECTORY"); ok {
			logrus.Debugf("Found path %s to use to init provider.", actualPath)
			return StorageConfig{
				Filesystem: &storage.FilesystemStorageConfig{
					Directory: actualPath,
				},
			}
		}

		path := ""
		if runtime.GOOS == "linux" {
			path = "/etc/arisu/tsubaki/storage"
		} else if runtime.GOOS == "windows" {
			appdata := os.Getenv("APPDATA")
			path = appdata + "\\Tsubaki\\Storage"
		} else {
			// TODO: macos - try to figure out what to use
			panic(fmt.Errorf("unsupported runtime: %s", runtime.GOOS))
		}

		logrus.Warn("It is recommended to set your own path if using the filesystem configuration with `TSUBAKI_STORAGE_FILESYSTEM_DIRECTORY` environment variable with a valid path, so you don't run into permission errors!")
		return StorageConfig{
			Filesystem: &storage.FilesystemStorageConfig{
				Directory: path,
			},
		}
	}

	panic("we should never end up here")
}

func getElasticsearchConfigFromEnv() *ElasticsearchConfig {
	logrus.Debug("Now loading Elasticsearch configuration from system environment variables...")

	// Check if it is enabled
	endpoints, enabled := os.LookupEnv("TSUBAKI_ELASTIC_ENDPOINTS")
	if !enabled {
		logrus.Warn("Missing Elasticsearch configuration! You will lose the `search` GraphQL query!")
		return nil
	}

	nodes := strings.Split(endpoints, ",")
	var password *string
	var username *string

	if value, ok := os.LookupEnv("TSUBAKI_ELASTIC_PASSWORD"); ok {
		password = &value
	}

	if value, ok := os.LookupEnv("TSUBAKI_ELASTIC_USERNAME"); ok {
		username = &value
	}

	return &ElasticsearchConfig{
		Password: password,
		Username: username,
		Hosts:    nodes,
	}
}

///////////////////////////////////////////////////////////
/////////////// ✨ LOADER FUNCTIONS :D ✨ /////////////////
//////////////////////////////////////////////////////////

func testConfig(path string) error {
	// Check if we can read the file
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var config Config
	err = yaml.Unmarshal(contents, &config)
	if err != nil {
		return err
	}

	if checkIfVariableExists("GO_ENV", config.Environment != "") {
		return MissingEnvironmentError
	}

	valid := []string{"development", "production"}
	isValid := false
	for _, key := range valid {
		if key == config.Environment.String() {
			isValid = true
		}
	}

	if !isValid {
		return InvalidEnvironmentError
	}

	return nil
}

func loadConfigFromPath(path string) (*Config, error) {
	if path != "" {
		// Check if it is in the root directory
		if _, err := os.Stat("./config.yml"); !os.IsNotExist(err) {
			logrus.Debugf("We will not use config in path '%s' since we found a config.yml file in the root directory.", path)
			path = "./config.yml"
		} else {
			logrus.Warnf("Unable to find the configuration path in %s (if it was loaded with `-c` or not), opting to use environment variables...", path)
			c, err := LoadFromEnv()
			if err != nil {
				return nil, err
			}

			// Check if we can find the secret key base
			if c.SecretKeyBase == "" {
				hash := util.GenerateHash(32)
				if hash == "" {
					panic("we should never be here (location=generating hash)")
				}

				c.SecretKeyBase = hash
				logrus.Warnf(
					"I was unable to find the `TSUBAKI_SECRET_KEY_BASE` environment variable, so I created it myself: %s",
					hash,
				)

				return c, nil
			}
		}
	}

	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(contents, &config)
	if err != nil {
		return nil, err
	}

	if checkIfVariableExists("GO_ENV", config.Environment != "") {
		return nil, MissingEnvironmentError
	}

	valid := []string{"development", "production"}
	isValid := false
	for _, key := range valid {
		if key == config.Environment.String() {
			isValid = true
		}
	}

	if !isValid {
		return nil, InvalidEnvironmentError
	}

	// Usually, some users can load this using an environment variable
	// for security reasons, so we can load that from here!
	if config.SecretKeyBase == "" && os.Getenv("TSUBAKI_SECRET_KEY_BASE") != "" {
		config.SecretKeyBase = os.Getenv("TSUBAKI_SECRET_KEY_BASE")
	}

	// Check if it is still empty
	if config.SecretKeyBase == "" {
		hash := util.GenerateHash(32)
		if hash == "" {
			return nil, errors.New("unable to generate hash")
		}

		logrus.Warnf(
			"it is recommended to store this generated key for JWT authentication. After a restart, all JWTs will fail: %s",
			hash,
		)

		config.SecretKeyBase = hash
	}

	return &config, nil
}

func LoadFromEnv() (*Config, error) {
	logrus.Debug("Now loading configuration from system environment variables...")

	// Check if .env exists in the root directory, if it does
	// let's load it!
	if _, err := os.Stat("./.env"); !os.IsNotExist(err) {
		err := godotenv.Load(".env", "./.env")
		if err != nil {
			return nil, err
		}
	}

	// Check if the Sentry DSN exists here
	var actualDsn *string
	if dsn, ok := os.LookupEnv("TSUBAKI_SENTRY_DSN"); ok {
		logrus.Debug("Found Sentry DSN %s!", dsn)
		d, err := sentry.NewDsn(dsn)
		if err != nil {
			return nil, err
		}

		uri := d.String()
		actualDsn = &uri
	}

	// Check if we can send telemetry events to our main servers.
	// If this is enabled, we will send telemetry events to Noelware.
	// By default, we do not enable it, but if you want to,
	// then you can.
	//
	// Curious? You can read up in the documentation:
	// https://docs.noelware.org/telemetry/services#arisu
	telemetryEnabled := convertToBool(os.Getenv("TSUBAKI_TELEMETRY_ENABLED"), false)
	if telemetryEnabled {
		logrus.Warn("Looks like you enabled Telemetry on this Tsubaki instance (probably by accident; https://docs.noelware.org/telemetry/services#arisu)")
	}

	port := 28093
	host := "0.0.0.0"

	if h, ok := os.LookupEnv("TSUBAKI_HOST"); !ok {
		if o, k := os.LookupEnv("HOST"); k {
			host = o
		}
	} else {
		host = h
	}

	if p, ok := os.LookupEnv("TSUBAKI_PORT"); !ok {
		if o, k := os.LookupEnv("PORT"); k {
			h, err := convertToInt(o, 28093)
			if err != nil {
				return nil, err
			}

			port = h
		}
	} else {
		h, err := convertToInt(p, 28093)
		if err != nil {
			return nil, err
		}

		port = h
	}

	kafkaConfig := getKafkaConfigFromEnv()
	elasticConfig := getElasticsearchConfigFromEnv()
	redisConfig, err := getRedisConfigFromEnv()

	if err != nil {
		return nil, err
	}

	storageConfig := getStorageConfigFromEnv()

	// check if we can enable basic auth
	var password *string
	var username *string

	if user, ok := os.LookupEnv("TSUBAKI_AUTH_USERNAME"); ok {
		if pass, k := os.LookupEnv("TSUBAKI_AUTH_PASSWORD"); k {
			logrus.Debug("Basic authentication on non-authenticated routes is enabled!")
			password = &pass
			username = &user
		} else {
			return nil, errors.New("you are required to have `TSUBAKI_AUTH_PASSWORD` also with `TSUBAKI_AUTH_USERNAME`")
		}
	}

	logrus.Debug("Loaded configuration from system environment variables!")
	return &Config{
		SecretKeyBase: os.Getenv("TSUBAKI_SECRET_KEY_BASE"),
		ElasticSearch: elasticConfig,
		Registrations: convertToBool(os.Getenv("TSUBAKI_REGISTRATIONS_ENABLED"), true),
		InviteOnly:    convertToBool(os.Getenv("TSUBAKI_INVITE_ONLY"), false),
		Telemetry:     telemetryEnabled,
		SentryDSN:     actualDsn,
		Username:      username,
		Password:      password,
		Storage:       storageConfig,
		Kafka:         kafkaConfig,
		Redis:         *redisConfig,
		Port:          &port,
		Host:          &host,
	}, nil
}
