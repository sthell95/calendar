SHALL=/bin/bash

export CGO_ENABLED=0
export DSN=psql://gouser:gopassword@localhost:5432/gotest

up:
	@ echo "-> running docker"
	cd docker && docker-compose up -d && CGO_ENABLED=1 go test -race ../...
	make migrate
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
	@ CGO_ENABLED=1 go test  -race ./...
.PHONY: test

lint:
	@ echo "-> running linters ..."
	@ golangci-lint run --fast ./...
.PHONY: lint

migration create:
	migrate create -ext sql -dir ./db/migrations -seq `date +"%s.%N"` ' * 1000000)/1'
.PHONY: create

migrate:
	migrate -database postgres://gouser:gopassword@localhost:5432/calendar?sslmode=disable -path ./db/migrations/ up
.PHONY: migrate