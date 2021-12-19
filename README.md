<div align='center'>
  <h2>☔ Tsubaki</h2>
  <div align='center'>
    <a href="https://arisu.land"><strong>Website</strong></a>  •  <a href="https://arisu.land/discord"><strong>Discord</strong></a>   •   <a href="https://github.com/auguwu/Arisu/discussions"><strong>Discussions</strong></a>
  </div>
  <br />
  <blockquote>Translation made with simplicity, yet robust. Made with 💖 using <a href='https://typescriptlang.org'><strong>TypeScript</strong></a>, <a href='https://reactjs.org'><strong>React</strong></a> with <a href='https://nextjs.org'><strong>Next.js</strong></a>.</blockquote>
</div>

Tsubaki is the **backend** portion of Arisu. This has been moved from [Arisu/app](https://github.com/auguwu/Arisu/tree/master/app) because of two reasons:

- The JavaScript port will not work in scalability, so this is a re-worked version in a good language for managing HTTP services.
- The more complex Arisu gets, I don't think that keep maintaining the code will be matured in the long run.

The frontend repository and the GitHub bot will reside in the [main repository](https://github.com/auguwu/Arisu).

## Features
- :octocat: **Open Source and Free** — Arisu is 100% free and open source, so you can host your own instance if you please!
- ✨ **Monorepos** — Create multiple subprojects into a single repository without having to maintain multiple repositories.
- ⚡ **Robust** — Arisu makes your workflow way easier and manageable without any high latency.

## Tech Stack
### Backend
- [**PostgreSQL**](https://postgresql.org) (with [Prisma](https://prisma.io))
- [**Kafka**](https://kafka.apache.org) (optional)
- [**Redis**](https://redis.io)
- [**chi**](https://github.com/go-chi/chi)
- [**Go**](https://golang.org)

### APIs
- [**GraphQL**](https://graphql.org)

### Frontend
- [**TailwindCSS**](https://tailwindcss.com)
- [**React.js**](https://reactjs.org)
- [**Next.js**](https://nextjs.org)

## Projects
Arisu is split into multiple projects to reduce workloads.

- 🌌 [telemetry-server](https://github.com/arisuland/telemetry-server) - **Telemetry server to track error reports and data usage between services**
- :octocat: [github-bot](https://github.com/auguwu/Arisu/tree/github-bot) - **GitHub bot to sync your translation from both parties.**
- 💝 [fubuki](https://github.com/auguwu/Arisu/tree/web) - **Website portion of Arisu, made with React + Next.js**
- 🐳 [docs](https://github.com/arisuland/docs) - **Documentation site for Arisu.**
- ⛴ [cli](https://github.com/arisuland/cli) - **CLI to connect your Arisu project with your machine or any CI service.**

## Installation
You have four options to run Tsubaki:

- under [Codespaces](https://github.com/features/codespaces) with a [development container](./.devcontainer) under **Visual Studio Code**
- Installing it on **Kubernetes** with our [Helm chart](https://charts.arisu.land)
- [Cloning the repository and running it](#installation-locally)
- Under a [Docker container](#installation-docker)

### Prerequisites
Before running your own instance of Arisu, you are required to have:

- [**PostgreSQL**](https://postgresql.org)
- [**Meilisearch**](https://meilisearch.com)
- [**Redis**](https://redis.io)
- [**Go**](https://golang.org)
- 1GB of RAM higher on your system
- 2 CPU cores or higher on your system

Any optional software to optimize the experience of Arisu:

- [**Docker**](https://docker.io) - A containerization tool to run Tsubaki. Can be used with our [docker compose](./docker-compose.yml) file.
- [**Kafka**](https://kafka.apache.org) - Used for messaging queues to and from the [GitHub bot](https://github.arisu.land)

#### Installation: Docker
Before we get started, you will need [Docker](https://docker.io) required to be running, optionally Docker Desktop for Mac or Windows.

You are allowed to use a SemVer version from our [Releases](https://github.com/arisuland/tsubaki/releases) in the Docker image.

```shell
# 1. Pull our Docker image to your machine
$ docker pull arisuland/tsubaki:latest # Replace `:latest` with the version, i.e, `:1.0.0`

# 2. Create a network to run the application under localhost
$ docker network create <name> --bridge app-tier

# 3. Run the image with the network attached
$ docker run -d -p 8787:8787 --name tsubaki \
  -v <path to config.yml>:/opt/arisu/tsubaki/config.yml \
  -e TSUBAKI_PORT=8787 \
  arisuland/tsubaki:latest -c /opt/arisu/tsubaki/config.yml
```

#### Installation: Locally
You are required to have [**Git**](https://git-scm.com) on your machine before contributing / continuing.

```shell
# 1. Pull the repository down to your machine
$ git pull https://github.com/arisuland/tsubaki

# 2. Change the directory to `tsubaki` and run `go get`
$ cd tsubaki && go get

# 3. Build the binary with GNU Make
# If you're using Windows, install Make with Chocolatey: `choco install make`
$ make build

# 4. Run the binary
$ ./build/tsubaki -c config.yml # Add `.exe` if running on Windows
```

## Configuration
Tsubaki's configuration file can be set in any directory with the **.yml** extension, you can also use environment variables
prefixed with `TSUBAKI_` to override them from `config.yml` or set them if it doesn't exist in `config.yml`.

```yml
# If registrations should be enabled on the server. If not, you are
# required to create a user on the administration dashboard.
#
# Type: Boolean
# Default: true
# Variable: TSUBAKI_REGISTRATIONS
registrations: Boolean

# DSN URI to link up Sentry to Tsubaki. This will output GraphQL, request, and database
# errors towards Sentry.
#
# Type: String?
# Default: nil
# Variable: TSUBAKI_SENTRY_DSN
sentry_dsn: String?

# Returns the site-name to embed on the navbar of your Fubuki instance.
#
# Type: String?
# Default: Arisu
# Variable: TSUBAKI_SITE_NAME
site_name: String?

# Returns the site icon to use when displayed on the website.
#
# Type: String?
# Default: https://cdn.arisu.land/lotus.png
# Variable: TSUBAKI_SITE_ICON
site_icon: String?

# Returns a number on how many retries before panicking on an unavailable port
# to run Tsubaki on. To disable this, use `-1` as the value.
#
# Type: Int
# Default: 5
# Max: 15
# Min: -1
# Variable: TSUBAKI_PORT_RETRY
port_retry: Int

# Uses a host URI when launching the HTTP server. If you wish to keep Tsubaki
# running internally, you can use `127.0.0.1` instead of the default `0.0.0.0`.
#
# Type: String
# Default: "0.0.0.0"
# Variable: TSUBAKI_HOST or HOST
host: String

# Uses a different port other than 2809. If the port is taken, it will
# try to find an available one round-robin style and use that one that
# isn't taken. You can set how many tries using the `port_retry` config
# variable.
#
# Type: Int
# Default: 28093
# Range: 80-65535 on root; 1024-65535 on non-root
port: Int

# Returns the configuration for using the filesystem, S3,
# or Google Cloud Storage to backup your projects.
#
# Type: FilesystemConfig
# Default:
#   fs:
#     directory: <cwd>/.arisu
# Prefix: TSUBAKI_STORAGE
storage:
  # Configures using the filesystem to host your projects, once the data
  # is removed, Arisu will fix it but cannot restore your projects.
  # If you're using Docker or Kubernetes, it is best assured that you
  # must create a volume so Arisu can interact with it.
  #
  # Type: FileSystemStorageConfig
  # Aliases: `fs`
  # Variable: TSUBAKI_STORAGE_FS_*
  filesystem:
    # Returns the directory to store your projects in. Arisu will attempt
    # to create a `arisu.lock` file to show that the directory can be used.
    #
    # Type: String
    # Variable: TSUBAKI_STORAGE_FS_DIRECTORY
    # Default: <cwd>/.arisu
    directory: String

  # Configures using S3 to host your projects, once the bucket is gone,
  # Arisu will attempt to create the bucket but your data will be lost.
  #
  # Type: S3StorageConfig
  s3:
    # Returns the provider to use when authenticating to S3. Arisu supports
    # Amazon S3, Wasabi, or using Minio. By default, it will attempt to use
    # Amazon S3.
    #
    # Type: S3Provider?
    # Variable: TSUBAKI_STORAGE_S3_PROVIDER
    # Default: S3Provider.AMAZON
    provider: S3Provider

    # Returns the bucket to use when storing files. If this bucket
    # doesn't exist, Arisu will attempt to create the bucket.
    # By default, Arisu will use `arisu` as the default bucket name
    # if this is not set.
    #
    # Type: String
    # Variable: TSUBAKI_STORAGE_S3_BUCKET
    # Default: "arisu"
    bucket: String

    # Returns the access key for authenticating to S3. If this isn't provided,
    # it will attempt to look for your credentials stored in `~/.aws`. This is a
    # recommended variable to set if using the S3 provider.
    #
    # Type: String
    # Variable: TSUBAKI_STORAGE_S3_ACCESS_KEY
    # Default: "access_key" key in ~/.aws/tsubaki_config
    access_key: String

    # Returns the secret key for authenticating to S3. If this isn't provided,
    # it will attempt to look for your credentials stored in `~/.aws`. This is a
    # recommended variable to set if using the S3 provider.
    #
    # Type: String
    # Variable: TSUBAKI_STORAGE_S3_SECRET_KEY
    # Default: "access_key" key in ~/.aws/tsubaki_config
    secret_key: String

    # Returns the region to host your bucket, this is dependant on if you
    # created the bucket without running Tsubaki. This is required to set to
    # so no errors will occur while authenticating to S3.
    #
    # Type: String
    # Variable: TSUBAKI_STORAGE_S3_REGION
    # Default: "us-east1"
    region: String

# Configuration to setup a Kafka producer to use for messaging queues.
# This is required if you're running the GitHub bot.
#
# Type: KafkaConfig
# Default: nil
# Prefix: TSUBAKI_KAFKA
kafka:
  # If the producer should auto create the topic for you.
  #
  # Type: Boolean
  # Variable: TSUBAKI_KAFKA_AUTO_CREATE_TOPICS
  # Default: true
  auto_create_topics: Boolean

  # A list of brokers to connect to. This returns a List of `host:port` strings. If you're
  # using an environment variable to set this, split it with `,` so it can be registered properly!
  #
  # Type: List<String>
  # Variable: TSUBAKI_KAFKA_BROKERS
  # Default: localhost:9092
  brokers: List<String>

  # Returns the topic to send messages towards the GitHub bot to Tsubaki.
  # Warning: This must be the same as the one you set on the GitHub bot configuration
  # or it will not receive messages.
  #
  # Type: String
  # Variable: TSUBAKI_KAFKA_TOPIC
  # Default: `arisu:tsubaki`
  topic: String

# Returns the configuration for using Redis to cache user sessions.
# Sentinel and Standalone are supported.
#
# Type: RedisConfig
# Default:
#   redis:
#     host: localhost
#     port: 6379
# Prefix: TSUBAKI_REDIS
redis:
  # Sets a list of sentinel servers to use when using Redis Sentinel
  # instead of Redis Standalone. The `master` key is required if this
  # is defined. This returns a List of `host:port` strings. If you're
  # using an environment variable to set this, split it with `,` so it can be registered properly!
  #
  # Type: List<String>?
  # Variable: TSUBAKI_REDIS_SENTINEL_SERVERS
  # Default: nil
  sentinels: List<String>?

  # If `requirepass` is set on your Redis server, this property will authenticate
  # Tsubaki once the connection is being dealt with.
  #
  # Type: String?
  # Variable: TSUBAKI_REDIS_PASSWORD
  # Default: nil
  password: String?

  # Returns the master name for connecting to any redis sentinel servers.
  #
  # Type: String?
  # Variable: TSUBAKI_REDIS_SENTINEL_MASTER
  # Default: nil
  master: String?

  # Returns the database index to use so you don't clutter your server
  # with Tsubaki-set configs.
  #
  # Type: Int
  # Variable: TSUBAKI_REDIS_DATABASE_INDEX
  # Default: 5
  index: Int

  # Returns the host for connecting to Redis.
  #
  # Type: String
  # Default: "localhost"
  # Variable: TSUBAKI_REDIS_HOST
  host: String

  # Returns the port to use when connecting to Redis.
  #
  # Type: Int
  # Default: 6379
  # Range: 80-65535 on root; 1024-65535 on non-root
  # Variable: TSUBAKI_REDIS_PORT
  port: Int
```

## Frequently Asked Questions
### Q. Why is Meilisearch a required service to run?
> Since it's for indexing project and user names to use for the search bar, this will be under the [search](https://docs.arisu.land/graphql/queries#search)
query which is ***not authenticated*** but API requests have a ratelimiter!

### Q. Why Kafka over ...
> Because Kafka will be more stabilized when using multiple Kafka brokers rather than <insert service here>. :3

### Q. Why Go over ...?
> Personally, I really like Go for creating backend applications. I could've used Rust but there was no libraries for my use case
> that I could find (e.g. having Sentinel support for any Redis libraries)

## License
**Tsubaki** is released under the **GPL-3.0** License by Noelware. If you wish to view the whole license, read the [LICENSE](/LICENSE) file.
