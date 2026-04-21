.PHONY: build test release-build clean fmt

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
RELEASE_MANIFEST_URL ?= https://github.com/reorc/apimux-cli/releases/latest/download/latest.json
GORELEASER_CMD ?= go run github.com/goreleaser/goreleaser/v2@v2.12.7
LDFLAGS := -X github.com/reorc/apimux-cli/internal/buildinfo.Version=$(VERSION) \
	-X github.com/reorc/apimux-cli/internal/buildinfo.Commit=$(COMMIT) \
	-X github.com/reorc/apimux-cli/internal/buildinfo.BuildDate=$(BUILD_DATE) \
	-X github.com/reorc/apimux-cli/internal/buildinfo.ReleaseManifestURL=$(RELEASE_MANIFEST_URL)

build:
	mkdir -p dist
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o dist/apimux ./cmd/apimux

test:
	go test ./... -count=1

release-build:
	GORELEASER_CMD='$(GORELEASER_CMD)' ./scripts/build-apimux-cli.sh

fmt:
	gofmt -w ./cmd ./internal

clean:
	rm -rf dist
