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
	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
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
	if cfg, err := loadConfig(path); err != nil {
		return nil, err
	} else {
		return cfg, err
	}
}

func TestConfigFromPath(path string) error {
	if err := testConfig(path); err != nil {
		return err
	} else {
		return nil
	}
}

func checkIfMissing(key string, value bool) bool {
	if !value {
		if os.Getenv("TSUBAKI_"+key) != "" {
			return false
		}

		return true
	}

	return false
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

func getKafkaConfig() *KafkaConfig {
	// Check if keys exist
	enabled := os.Getenv("TSUBAKI_KAFKA_BROKERS") == "" || os.Getenv("TSUBAKI_KAFKA_TOPIC") == ""
	if !enabled {
		return nil
	}

	return &KafkaConfig{
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

	var sentinels []string
	if os.Getenv("TSUBAKI_REDIS_SENTINEL_SERVERS") != "" {
		sentinels = strings.Split(os.Getenv("TSUBAKI_REDIS_SENTINEL_SERVERS"), ",")
	}

	return &RedisConfig{
		Sentinels:  &sentinels,
		Password:   password,
		MasterName: masterName,
		DbIndex:    dbIndex,
		Host:       fallbackEnvString(os.Getenv("TSUBAKI_REDIS_HOST"), "localhost"),
		Port:       redisPort,
	}, nil
}

func getStorageConfig() (StorageConfig, error) {
	// Check if the directory env variable exists
	if os.Getenv("TSUBAKI_STORAGE_FS_DIRECTORY") != "" {
		return StorageConfig{
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

		return StorageConfig{
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

	return StorageConfig{}, errors.New("missing storage provider, read more here: https://docs.arisu.land/selfhosting/storage")
}

func getElasticSearchConfig() *ElasticsearchConfig {
	// Check if the required env variable exists
	if os.Getenv("TSUBAKI_ELASTIC_ENDPOINTS") != "" {
		var password *string
		var username *string

		nodes := strings.Split(os.Getenv("TSUBAKI_ELASTIC_ENDPOINTS"), ";")

		if os.Getenv("TSUBAKI_ELASTIC_PASSWORD") != "" {
			raw := os.Getenv("TSUBAKI_ELASTIC_PASSWORD")
			password = &raw
		}

		if os.Getenv("TSUBAKI_ELASTIC_USERNAME") != "" {
			raw := os.Getenv("TSUBAKI_ELASTIC_USERNAME")
			username = &raw
		}

		return &ElasticsearchConfig{
			Password: password,
			Username: username,
			Hosts:    nodes,
		}
	}

	logrus.Warnf("Missing ElasticSearch configuration. This is not required unless you need the `search` GraphQL query.")
	return nil
}

func checkIfS3EnvExists() bool {
	prefix := "TSUBAKI_STORAGE_S3_"
	provider := os.Getenv(prefix + "PROVIDER")
	region := os.Getenv(prefix + "REGION")

	return provider != "" && region != ""
}

func testConfig(path string) error {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var config Config
	err = yaml.Unmarshal(contents, &config)
	if err != nil {
		return nil
	}

	if checkIfMissing("GO_ENV", config.Environment != "") {
		return MissingEnvironmentError
	}

	validEnvs := []string{"development", "production"}
	var valid = false

	for _, key := range validEnvs {
		if key == config.Environment.String() {
			valid = true
		}
	}

	if !valid {
		return InvalidEnvironmentError
	}

	return nil
}

func loadConfig(path string) (*Config, error) {
	logrus.Debug("Loading configuration...")
	if path != "" {
		// Check if it is in the root directory
		_, err := os.Stat("./config.yml")
		if !os.IsNotExist(err) {
			logrus.Debug("Found configuration in root, loading from that...")
			path = "./config.yml"
		} else {
			logrus.Warn("Unable to find configuration in root directory or with `-c` flag. Loading from environment variables...")
			c, err := LoadFromEnv()
			if err != nil {
				return nil, err
			}

			// Check if it is still empty
			if c.SecretKeyBase == "" {
				hash := util.GenerateHash(32)
				if hash == "" {
					return nil, errors.New("unable to generate hash")
				}

				logrus.Warnf(
					"it is recommended to store this generated key for JWT authentication. After a restart, all JWTs will fail: %s",
					hash,
				)

				c.SecretKeyBase = hash
			}

			return c, nil
		}
	}

	logrus.Debugf("Found configuration in path %s! Now loading...", path)
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(contents, &config)
	if err != nil {
		return nil, err
	}

	if checkIfMissing("GO_ENV", config.Environment != "") {
		return nil, MissingEnvironmentError
	}

	validEnvs := []string{"development", "production"}
	var valid = false

	for _, key := range validEnvs {
		if key == config.Environment.String() {
			valid = true
		}
	}

	if !valid {
		return nil, InvalidEnvironmentError
	}

	// If the secret key base environment variable exists,
	// let's set it in the config.
	if os.Getenv("TSUBAKI_SECRET_KEY_BASE") != "" && config.SecretKeyBase == "" {
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
	logrus.Debug("Now loading configuration from environment variables...")

	// Load from .env if the file exists
	if _, err := os.Stat("./.env"); !os.IsNotExist(err) {
		err := godotenv.Load(".env")
		if err != nil {
			return nil, err
		}
	}

	envDsn := os.Getenv("TSUBAKI_SENTRY_DSN")

	var sentryDsn *string
	if envDsn != "" {
		logrus.Debugf("Now verifying DSN %s...", envDsn)
		dsn, err := sentry.NewDsn(envDsn)

		if err != nil {
			return nil, err
		}

		uri := dsn.String()
		sentryDsn = &uri
	}

	telemetry := convertToBool(os.Getenv("TSUBAKI_TELEMETRY"), false)
	if telemetry {
		logrus.Warn("You have enabled telemetry reports! You have been warned... :ghost:")
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
	elasticCfg := getElasticSearchConfig()
	redisConfig, err := getRedisConfig()

	if err != nil {
		return nil, err
	}

	storageConfig, err := getStorageConfig()
	if err != nil {
		return nil, err
	}

	logrus.Debug("Successfully loaded configuration from environment variables!")
	return &Config{
		SecretKeyBase: os.Getenv("TSUBAKI_SECRET_KEY_HASH"),
		ElasticSearch: elasticCfg,
		Registrations: convertToBool(os.Getenv("TSUBAKI_REGISTRATIONS"), true),
		InviteOnly:    convertToBool(os.Getenv("TSUBAKI_INVITE_ONLY"), false),
		Telemetry:     telemetry,
		SentryDSN:     sentryDsn,
		Storage:       storageConfig,
		Kafka:         kafkaConfig,
		Redis:         *redisConfig,
		Port:          &port,
		Host:          &host,
	}, nil
}
