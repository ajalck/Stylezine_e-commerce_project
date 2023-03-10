SHELL := /bin/bash

.PHONY: all build test deps deps-cleancache

GOCMD=go
BUILD_DIR=build
BINARY_DIR=$(BUILD_DIR)/bin
CODE_COVERAGE=code-coverage

all: test build

${BINARY_DIR}:
	mkdir -p $(BINARY_DIR)

build: ${BINARY_DIR} ## Compile the code, build Executable File
	$(GOCMD) build -o $(BINARY_DIR) -v ./cmd/api

run: ## Start application
	$(GOCMD) run ./

test: ## Run tests
	$(GOCMD) test ./... -v -cover

test-coverage: ## Run tests and generate coverage file
	$(GOCMD) test ./... -coverprofile=$(CODE_COVERAGE).out
	$(GOCMD) tool cover -html=$(CODE_COVERAGE).out

deps: ## Install dependencies
	# go get $(go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
	$(GOCMD) get -u -t -d -v ./...
	$(GOCMD) mod tidy
	$(GOCMD) mod vendor

deps-cleancache: ## Clear cache in Go module
	$(GOCMD) clean -modcache

wire: ## Generate wire_gen.go
	cd pkg/di && wire

# swag: ## Generate swagger docs
# 	swag init -g pkg/api/middleware/auth.handler.go -o cmd/api/docs

# swag:## Genarate swagger docs
# 	cd cmd/api && swag init --parseDependency --parseInternal --parseDepth 1 -md ./documentation -o ./docs

swag: ## Generate swagger2 docs
	swag init -g pkg/api/handler/auth.handler.go --parseDependency -o ./docs

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

spath:
	cd cmd/api && PATH="$GOPATH/bin:$PATH" && export GOPATH="$HOME/go" && PATH="$GOPATH/bin:$PATH"