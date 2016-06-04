# Boilerplate code to setup make
MAKEFLAGS += --warn-undefined-variables
SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := all

# Watch all source files
SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=graphql-server

VERSION := $(shell cat version)
BUILD_TIME := $(shell date +%FT%T%z)
BUILDER := $(shell echo "`git config user.name` <`git config user.email`>")
GOVERSION := $(shell go version)

NAME ?= nilstgmd/$(BINARY)
TAG ?= $(NAME):$(VERSION)

LDFLAGS := -X 'main.version=$(VERSION)' \
		-X 'main.buildTime=$(BUILD_TIME)' \
		-X 'main.builder=$(BUILDER)' \
		-X 'main.goversion=$(GOVERSION)'

# Docker Compose based builds
.PHONY: all
all: clean test lint build run

.PHONY: clean
clean:
	@echo "Cleanup old containers ..."
	@docker-compose stop && docker-compose rm --all --force

.PHONY: test
test:
	@echo "Starting tests ..."
	@go test $$(go list ./... | grep -v /vendor/ | grep -v /cmd/)

.PHONY: lint
lint:
	@echo "Execute go vet and golint ..."
	@go vet -v $$(go list ./... | grep -v /vendor/ | grep -v /cmd/)
	@golint ./...

.PHONY: build
build:
	@echo "Starting build ..."
	@env GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BINARY)
	@docker build -t $(TAG) --rm=true $(SOURCEDIR)

.PHONY: run
run:
	@echo "Starting containers ..."
	@docker-compose up -d

# Stand-alone Docker based builds
.PHONY: $(BINARY)
$(BINARY):
	@docker run -d -p 8080:8080 --name $(BINARY) --link mongo:mongo --link cassandra:cassandra $(TAG)

.PHONY: mongo
mongo:
	@docker run --name mongo -d mongo

.PHONY: cassandra
cassandra:
	@docker run --name cassandra -d cassandra:3.5
