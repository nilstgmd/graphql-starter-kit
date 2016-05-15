# Watch all source files
SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=graphql-server

VERSION := $(shell cat version)
BUILD_TIME=`date +%FT%T%z`

NAME?=nilstgmd/$(BINARY)
TAG?=$(NAME):$(VERSION)

LDFLAGS=-ldflags "-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME}"

.DEFAULT_GOAL: $(BINARY)
$(BINARY): $(SOURCES)
	@go build $(LDFLAGS) -o $(BINARY) -v

.PHONY: clean
clean:
	@if [ -f ${BINARY} ] ; then rm ${BINARY}; go clean ./...; fi

.PHONY: test
test:
	@go test ./...

.PHONY: docker
docker: dockerclean dockerbuild dockerrun

.PHONY: dockerclean
dockerclean:
	@docker ps -a -q --filter ancestor=$(TAG) --format="{{.ID}}" | xargs docker stop | xargs docker rm

.PHONY: dockerbuild
dockerbuild:
	@env GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY)
	@docker build -t $(TAG) --rm=true $(SOURCEDIR)

.PHONY: dockerrun
dockerrun:
	@docker run -d -p 8080:8080 --name $(BINARY) --link mongo:mongo --link cassandra:cassandra $(TAG)

.PHONY: mongo
mongo:
	@docker run --name mongo -d mongo

.PHONY: cassandra
cassandra:
	@docker run --name cassandra -d cassandra:3.5
