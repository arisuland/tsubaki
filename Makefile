# ☔ Arisu: Translation made with simplicity, yet robust.
# Copyright (C) 2020-2022 Noelware
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <https://www.gnu.org/licenses/>.

JQ := $(shell command -v jq 2>/dev/null)
ifndef JQ
	$(error "`jq` is missing. please install jq!")
endif

VERSION    := $(shell cat version.json | jq .version | tr -d '"')
COMMIT_SHA := $(shell git rev-parse --short HEAD)
BUILD_DATE := $(shell go run ./cmd/build-date/main.go)
GIT_TAG    ?= $(shell git describe --tags --match "v[0-9]*")

GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

ifeq ($(GOOS), linux)
	TARGET_OS ?= linux
else ifeq ($(GOOS),darwin)
	TARGET_OS ?= darwin
else ifeq ($(GOOS),windows)
	TARGET_OS ?= windows
else
	$(error System $(GOOS) is not supported at this time)
endif

EXTENSION :=
ifeq ($(TARGET_OS),windows)
	EXTENSION := .exe
endif

# Usage: `make deps`
deps:
	@echo Updating dependency tree...
	go mod tidy
	go mod download
	@echo Updated dependency tree successfully.

# Usage: `make build`
build:
	@echo Now building Tsubaki for $(GOOS)!
	go build -ldflags "-s -w -X arisu.land/tsubaki/internal.Version=${VERSION} arisu.land/tsubaki/internal.CommitSHA=${COMMIT_SHA} \"arisu.land/tsubaki/internal.BuildDate=${BUILD_DATE}\"" -o ./bin/tsubaki$(EXTENSION)
	@echo Successfully built the binary. Use './bin/tsubaki$(EXTENSION) -c config.yml' to run!

# Usage: `make clean`
clean:
	@echo Now cleaning project..
	rm -rf bin/ .profile/
	go clean
	@echo Done!

# Usage: `make fmt`
fmt:
	@echo Formatting project...
	go fmt
	@echo Formatted!

# Usage: `make db.migrate NAME=...`
db.migrate:
	@echo Migrating database for development...
	go run github.com/prisma/prisma-client-go migrate dev --name=$(NAME)

# Usage: `make db.fmt`
db.fmt:
	@echo Formatting Prisma schema...
	go run github.com/prisma/prisma-client-go format

# Usage: `make db.generate`
db.generate:
	@echo Generating Prisma artifacts...
	go run github.com/prisma/prisma-client-go generate

# Usage: `make docgen`
docgen:
	@echo Now building documentation schema...
	go run ./cmd/docgen/main.go
