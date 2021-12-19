VERSION=$(shell cat version.json | jq .version | tr -d '"')
GIT_COMMIT=$(shell git rev-parse --short HEAD)
BUILD_DATE=$(shell go run ./cmd/build-date/main.go)

HOME_OS ?= $(shell go env GOOS)
ifeq ($(HOME_OS),linux)
	TARGET_OS ?= linux
else ifeq ($(HOME_OS),darwin)
	TARGET_OS ?= darwin
else ifeq ($(HOME_OS),windows)
	TARGET_OS ?= windows
else
	$(error System $(LOCAL_OS) is not supported.)
endif

ifeq ($(HOME_OS),windows)
	TARGET_FILE ?= ./build/tsubaki.exe
	TARGET_FILE_DOCGEN ?= ./main
else
	TARGET_FILE ?= ./build/tsubaki
	TARGET_FILE_CODEGEN ?= ./main
endif

# Usage: `make build`
build:
	@echo Building Tsubaki...
	go build -ldflags "-s -w -X main.version=${VERSION} -X main.commitHash=${GIT_COMMIT} -X \"main.buildDate=${BUILD_DATE}\"" -o $(TARGET_FILE)
	@echo Successfully built Tsubaki! Use '$(TARGET_FILE) -c config.yml' to run!

# Usage: `make build.docker`
build.docker:
	@echo Building Tsubaki Docker image...
	docker build . -t "arisuland/tsubaki:latest" --no-cache --build-arg VERSION=${VERSION} --build-arg COMMIT_HASH=${GIT_COMMIT} --build-arg BUILD_DATE=${BUILD_DATE}
	docker build . -t "arisuland/tsubaki:${VERSION}" --no-cache
	@echo Done building images for latest and ${VERSION} tags!

clean:
	@echo Cleaning build/ directories
	go clean
	rm -rf build/
	@echo Done!

# Usage: `make fmt`
fmt:
	go fmt

# Usage: `make db.migrate NAME=<string>`
db.migrate:
	@echo Migrating development database...
	go run github.com/prisma/prisma-client-go migrate dev --name=$(NAME)

# Usage: `make db.fmt`
db.fmt:
	@echo Formatting .prisma files...
	go run github.com/prisma/prisma-client-go format

# Usage: `make db.generate`
db.generate:
	@echo Generating Prisma artifacts...
	go run github.com/prisma/prisma-client-go generate

# Usage: `make docgen`
docgen:
	go build cmd/docgen/main.go
	$(TARGET_FILE_DOCGEN)
