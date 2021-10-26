SHALL=/bin/bash

export CGO_ENABLED=0
export DSN=psql://gouser:gopassword@localhost:5432/gotest

up:
	@ echo "-> running docker"
	cd docker && docker-compose up -d && CGO_ENABLED=1 go test -race ../cmd...
.PHONEL: up

down:
	@ echo "-> down docker"
	@ cd docker && docker-compose down
.PHONEL: up

default: build
.PHONY: default

build:
	@ echo "-> build binary ..."
	@ go build -ldflags "-X main.HashCommit=`git rev-parse HEAD` -X main.BuildStamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'`" -o ./calendar ./cmd
.PHONY: build

test:
	@ echo "-> running tests ..."
	@ CGO_ENABLED=1 go test  -race ./cmd...
.PHONY: test

lint:
	@ echo "-> running linters ..."
	@ golangci-lint run --fast ./cmd...
.PHONY: lint