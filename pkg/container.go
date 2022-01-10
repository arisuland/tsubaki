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
	"arisu.land/tsubaki/internal"
	"arisu.land/tsubaki/pkg/storage"
	"arisu.land/tsubaki/prisma/db"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bwmarrin/snowflake"
	es "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/getsentry/sentry-go"
	"github.com/go-redis/redis/v8"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
	"time"
)

// GlobalContainer represents the global Container instance
// that is constructed using the NewContainer function.
var GlobalContainer *Container = nil

// Container is a object that holds all the dependencies for every part
// of Tsubaki's lifecycle.
type Container struct {
	ElasticSearch *es.Client
	Snowflake     *snowflake.Node
	Storage       storage.BaseStorageProvider
	Prisma        *db.PrismaClient
	Sentry        *sentry.Client
	Config        *Config
	Redis         *redis.Client
	Kafka         *kafka.Writer
}

func NewContainer(path string) error {
	if GlobalContainer != nil {
		panic("tried to create a new global container, but one already exists")
	}

	snowflake.Epoch = int64(1641020400000)
	logrus.Debug("Creating container...")

	config, err := NewConfig(path)
	if err != nil {
		return err
	}

	// Register Prometheus objects
	internal.RegisterMetrics()

	// create snowflake creator
	node, err := snowflake.NewNode(0)
	if err != nil {
		return err
	}

	// Create the Prisma client
	logrus.Debug("Connecting to PostgreSQL...")
	prisma := db.NewClient()
	if err := prisma.Connect(); err != nil {
		return err
	}

	logrus.Debug("Connected to PostgreSQL successfully!")
	users, err := prisma.User.FindMany().Exec(context.TODO())
	if err != nil {
		return err
	}

	internal.UsersCountMetric.Set(float64(len(users)))

	// Create Redis client
	logrus.Debug("Now connecting to Redis...")

	password := ""
	if config.Redis.Password != nil {
		password = *config.Redis.Password
	}

	var re *redis.Client
	if len(config.Redis.Sentinels) > 0 {
		masterName := ""
		if config.Redis.MasterName == nil {
			return errors.New("config option 'redis.master_name' needs to be defined to use a sentinel connection")
		} else {
			masterName = *config.Redis.MasterName
		}

		re = redis.NewFailoverClient(&redis.FailoverOptions{
			SentinelAddrs: config.Redis.Sentinels,
			MasterName:    masterName,
			Password:      password,
			DB:            config.Redis.DbIndex,
			DialTimeout:   10 * time.Second,
			ReadTimeout:   15 * time.Second,
			WriteTimeout:  15 * time.Second,
		})
	} else {
		re = redis.NewClient(&redis.Options{
			Password:     password,
			Addr:         fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
			DB:           config.Redis.DbIndex,
			DialTimeout:  10 * time.Second,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
		})
	}

	logrus.Debug("Created Redis client, checking connection...")
	if err := re.Ping(context.TODO()).Err(); err != nil {
		return err
	}

	logrus.Debug("Connected to Redis!")

	// create storage provider
	logrus.Info("Now creating storage provider...")
	var provider storage.BaseStorageProvider
	if config.Storage.Filesystem != nil {
		provider = storage.NewFilesystemStorageProvider(*config.Storage.Filesystem)
		if err := provider.Init(); err != nil {
			return err
		}
	} else if config.Storage.S3 != nil {
		provider = storage.NewS3StorageProvider(config.Storage.S3)
		if err := provider.Init(); err != nil {
			return err
		}
	} else {
		return errors.New("missing storage provider to use")
	}

	logrus.Info("Created storage provider!")
	var writer *kafka.Writer

	if config.Kafka != nil {
		logrus.Debug("Creating Kafka producer...")
		writer = &kafka.Writer{
			Addr:     kafka.TCP(config.Kafka.Brokers...),
			Topic:    config.Kafka.Topic,
			Balancer: &kafka.LeastBytes{},
		}
	}

	var sc *sentry.Client
	if config.SentryDSN != nil {
		hostName, err := os.Hostname()
		if err != nil {
			hostName = "localhost"
		}

		client, err := sentry.NewClient(sentry.ClientOptions{
			Dsn:              *config.SentryDSN,
			AttachStacktrace: true,
			SampleRate:       1.0,
			ServerName:       fmt.Sprintf("arisu.tsubaki v%s @ %s", internal.Version, hostName),
		})

		if err != nil {
			return err
		}

		sc = client
	}

	var e *es.Client
	if config.ElasticSearch != nil {
		logrus.Debug("Connecting to ElasticSearch...")
		if config.ElasticSearch.Username == nil || config.ElasticSearch.Password == nil {
			logrus.Warn("It is recommended to keep your ElasticSearch instance secured~")
		}

		cfg := es.Config{
			Addresses:            config.ElasticSearch.Hosts,
			DiscoverNodesOnStart: true,
		}

		if config.ElasticSearch.Username != nil {
			cfg.Username = *config.ElasticSearch.Username
		}

		if config.ElasticSearch.Password != nil {
			cfg.Password = *config.ElasticSearch.Password
		}

		client, err := es.NewClient(cfg)
		if err != nil {
			return err
		}

		e = client

		// check if we can query
		res, err := client.Info()
		if err != nil {
			return err
		}

		// unmarshal it from `res`
		defer func() {
			_ = res.Body.Close()
		}()

		var data map[string]interface{}
		err = json.NewDecoder(res.Body).Decode(&data)
		if err != nil {
			return err
		}

		serverVersion := data["version"].(map[string]interface{})["number"].(string)
		logrus.Debugf("Server: %s | Client: %s", serverVersion, es.Version)

		// Now we have to index documents, in a separate goroutine
		_ = indexDocuments(prisma, client)
	}

	GlobalContainer = &Container{
		ElasticSearch: e,
		Snowflake:     node,
		Storage:       provider,
		Prisma:        prisma,
		Sentry:        sc,
		Config:        config,
		Redis:         re,
		Kafka:         writer,
	}

	return nil
}

func (c *Container) Close() error {
	s := time.Now()
	if err := c.Redis.Close(); err != nil {
		return err
	}

	logrus.Warnf("Disconnected from Redis in %s!", time.Since(s).String())
	s = time.Now()
	if err := c.Prisma.Disconnect(); err != nil {
		return err
	}

	logrus.Warnf("Disconnected from PostgreSQL in %s!", time.Since(s).String())
	if c.Kafka != nil {
		logrus.Warn("Closing off Kafka broker...")
		if err := c.Kafka.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (c *Container) DbPing() int64 {
	t := time.Now()
	if _, err := c.Prisma.User.FindMany().Exec(context.TODO()); err != nil {
		return -1
	} else {
		return time.Since(t).Milliseconds()
	}
}

func (c *Container) RedisPing() int64 {
	t := time.Now()
	if err := c.Redis.Ping(context.TODO()).Err(); err != nil {
		return -1
	} else {
		return time.Since(t).Milliseconds()
	}
}

func indexDocuments(prisma *db.PrismaClient, client *es.Client) error {
	logrus.Debug("Now indexing documents...")
	wg := sync.WaitGroup{}

	// query all projects
	t := time.Now()
	projects, err := prisma.Project.FindMany().Exec(context.TODO())
	if err != nil {
		logrus.Fatalf("Unable to retrieve projects: %v", err)
		return err
	}

	logrus.Debugf("Took %s to grab all projects.", time.Since(t).String())

	t = time.Now()
	// query all users
	users, err := prisma.User.FindMany().Exec(context.TODO())
	if err != nil {
		logrus.Fatalf("Unable to retrieve users: %v", err)
		return err
	}

	logrus.Debugf("Took %s to grab all users.", time.Since(t).String())

	// Now, we actually index them!
	for i, project := range projects {
		logrus.Debugf("Now indexing project %s...", project.ID)

		wg.Add(1)
		go func(i int, project db.ProjectModel) {
			t := time.Now()

			// unmarshal it from json
			data, err := json.Marshal(&project.InnerProject)
			if err != nil {
				logrus.Fatalf("Unable to unmarshal project %s: %v", project.ID, err)
				return
			}

			// Setup the request payload
			req := esapi.IndexRequest{
				Index:      "arisu:tsubaki",
				DocumentID: project.ID,
				Body:       bytes.NewReader(data),
				Refresh:    "true",
			}

			// Perform the request!!!
			res, err := req.Do(context.TODO(), client)
			if err != nil {
				logrus.Fatalf("Unable to get response from ElasticSearch: %v", err)
				return
			}

			defer func() {
				_ = res.Body.Close()
			}()

			logrus.Debug("Took %s to send out a request to server", time.Since(t).String())
			if res.IsError() {
				logrus.Fatalf("Unable to index project %s [%s]", project.ID, res.Status())
			} else {
				var data map[string]interface{}
				if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
					logrus.Fatalf("Unable to parse request body: %v", err)
				} else {
					logrus.Debugf("Indexed project %s successfully! version=%d", project.ID, int(data["_version"].(float64)))
				}
			}
		}(i, project)
	}

	// now we index the users also!! :D
	// Now, we actually index them!
	for i, user := range users {
		logrus.Debugf("Now indexing user %s...", user.ID)

		wg.Add(1)
		go func(i int, user db.UserModel) {
			t := time.Now()

			// unmarshal it from json
			u := FromDbUserModel(user)
			data, err := json.Marshal(&u)
			if err != nil {
				logrus.Fatalf("Unable to unmarshal project %s: %v", user.ID, err)
				return
			}

			// Setup the request payload
			req := esapi.IndexRequest{
				Index:      "arisu:tsubaki",
				DocumentID: user.ID,
				Body:       bytes.NewReader(data),
				Refresh:    "true",
			}

			// Perform the request!!!
			res, err := req.Do(context.TODO(), client)
			if err != nil {
				logrus.Fatalf("Unable to get response from ElasticSearch: %v", err)
				return
			}

			defer func() {
				_ = res.Body.Close()
			}()

			logrus.Debug("Took %s to send out a request to server", time.Since(t).String())
			if res.IsError() {
				logrus.Fatalf("Unable to index project %s [%s]", user.ID, res.Status())
			} else {
				var data map[string]interface{}
				if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
					logrus.Fatalf("Unable to parse request body: %v", err)
				} else {
					logrus.Debugf("Indexed user %s successfully! version=%d", user.ID, int(data["_version"].(float64)))
				}
			}
		}(i, user)
	}

	return nil
}
