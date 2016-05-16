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
BUILD_TIME=`date +%FT%T%z`

NAME?=nilstgmd/$(BINARY)
TAG?=$(NAME):$(VERSION)

LDFLAGS=-ldflags "-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME}"

# Docker Compose based builds
.PHONY: all
all: clean build run

.PHONY: clean
clean:
	@docker-compose stop && docker-compose rm --all --force

.PHONY: build
build:
	@env GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY)
	@docker build -t $(TAG) --rm=true $(SOURCEDIR)

.PHONY: run
run:
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
