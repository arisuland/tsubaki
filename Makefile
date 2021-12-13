VERSION=$(shell cat version.json | jq .version | tr -d '"')
GIT_COMMIT=$(shell git rev-parse --short HEAD)
BUILD_DATE=$(shell date +"%D - %r")

# Usage: `make build`
build:
	@echo "Building Tsubaki..."
	go build -ldflags "-s -w -X main.version=${VERSION} -X main.commitHash=${GIT_COMMIT}" -o ./build/tsubaki
	@echo "Successfully built Tsubaki! Use './build/tsubaki -c config.yml' to run!"

# Usage: `make build.docker`
build.docker:
	@echo "Building Tsubaki Docker image..."
	docker build . -t "arisuland/tsubaki:latest" --no-cache --build-arg VERSION=${VERSION} --build-arg COMMIT_HASH=${GIT_COMMIT} --build-arg BUILD_DATE=${BUILD_DATE}
	docker build . -t "arisuland/tsubaki:${VERSION}" --no-cache
	@echo "Done building images for latest and ${VERSION} tags!"

# Usage: `make fmt`
fmt:
	go fmt

# Usage: `make db.migrate NAME=<string>`
db.migrate:
	@echo "Migrating development database..."
	go run github.com/prisma/prisma-client-go migrate dev --name=$(NAME)

# Usage: `make codegen`
codegen:
	go build cmd/codegen -ldflags "-s -w" -o ./build/codegen
	./build/codegen

# Usage: `make docgen`
docgen:
	go build cmd/docgen -ldflags "-s -w" -o ./build/docgen
	./build/docgen
