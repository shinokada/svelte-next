BINARY  := svelte-next
VERSION ?= $(shell git describe --tags --always 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X github.com/shinokada/svelte-next/cmd.Version=$(VERSION)"

.PHONY: all build go.sum test test-verbose lint tidy clean release-dry install help

all: build

## build: compile the binary into ./bin/
build: go.sum
	go build $(LDFLAGS) -o bin/$(BINARY) .

## go.sum: generate go.sum if missing
go.sum:
	go mod tidy

## test: run all unit and integration tests
test: go.sum
	go test ./...

## test-verbose: run tests with output
test-verbose:
	go test -v ./...

## lint: run golangci-lint (must be installed)
lint:
	golangci-lint run ./...

## tidy: update go.sum and tidy dependencies
tidy:
	go mod tidy

## clean: remove build artefacts
clean:
	rm -rf bin/ dist/

## release-dry: run goreleaser in snapshot mode (no publish)
release-dry:
	goreleaser release --snapshot --clean

## install: install binary to GOPATH/bin
install:
	go install $(LDFLAGS) .

help:
	@grep -E '^## ' Makefile | sed 's/## //'
