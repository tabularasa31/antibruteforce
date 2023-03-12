#!make
include .env.example

BIN_APP="./bin/antibruteforce"

# ==============================================================================
# Main
run: build
	$(BIN_APP)

build:
	go build -v -o $(BIN_APP) ./cmd/antibruteforce/

test:
	go test -cover ./...

integration-test:
	go test -tags integration ./tests/integration/...

generate:
	 go generate ./...

# ==============================================================================
# Tools commands
install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.50.1

lint: install-lint-deps ### check by golangci linter
	echo "Starting linters"
	golangci-lint run

# ==============================================================================
# Docker compose commands
up:
	echo "Starting docker environment"
	docker-compose -f docker-compose.yml up --build

down:
	docker-compose -f docker-compose.yml down

# ==============================================================================
.PHONY: build run test lint up
