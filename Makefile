SHALL=/bin/bash

export CGO_ENABLED=0
export DSN=psql://gouser:gopassword@localhost:5432/gotest

up:
	@ echo "-> running docker"
	@ cd docker && docker-compose up -d \
 	@ go test -ncolor
.PHONEL: up

down:
	@ echo "-> down docker"
	@ cd docker && docker-compose down
.PHONEL: up

default: build
.PHONY: default

build:
	@ echo "-> build binary ..."
	@ go build -ldflags "-X main.HashCommit=`git rev-parse HEAD` -X main.BuildStamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'`" -o ./calendar .
.PHONY: build

test:
	@ echo "-> running tests ..."
	@ CGO_ENABLED=1 go test -race ./...
.PHONY: test

lint:
	@ echo "-> running linters ..."
	@ golangci-lint run --enable-all ./...
.PHONY: lint