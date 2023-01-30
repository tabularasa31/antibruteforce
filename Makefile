# ==============================================================================
# Main
run:
	go run ./cmd/main.go

build:
	go build ./cmd/main.go

test:
	go test -cover ./...


# ==============================================================================
# Tools commands

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps ### check by golangci linter
	echo "Starting linters"
	golangci-lint run

