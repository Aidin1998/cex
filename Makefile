# Go parameters
GOCMD      ?= go
GOBUILD    := $(GOCMD) build
GOCLEAN    := $(GOCMD) clean
GOTEST     := $(GOCMD) test
GOGET      := $(GOCMD) get
GOFMT      := gofumpt
LINTER     := golangci-lint
LINTCONFIG ?= .golangci.yml


# Binary names
USERS_BINARY  	:= cex-users
ACCOUNTS_BINARY := cex-accounts

# Directories
BUILD_DIR	?= ./build/
USERS_DIR 	:= ./cmd/users
ACCOUNTS_DIR	:= ./cmd/accounts

.PHONY: all build clean test fmt lint run deps build-users build-accounts run-users run-accounts

all: test build

build: build-users build-accounts

build-users:
	$(GOBUILD) -o $(USERS_BINARY) -v $(USERS_DIR)

build-accounts:
	$(GOBUILD) -o $(ACCOUNTS_BINARY) -v $(ACCOUNTS_DIR)

clean:
	$(GOCLEAN)
	rm -f $(BUILD_DIR)

test:
	$(GOTEST) -v ./...

fmt:
	$(GOFMT) -l -w .

lint:
	$(LINTER) run --config $(LINTCONFIG)

deps:
	$(GOGET) -v ./...

# Cross compilation
build-linux-users:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(USERS_BINARY)_linux -v $(USERS_DIR)

build-linux-accounts:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(ACCOUNTS_BINARY)_linux -v $(ACCOUNTS_DIR)

build-linux: build-linux-users build-linux-accounts

