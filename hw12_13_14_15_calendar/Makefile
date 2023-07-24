include ./configs/.env

BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) -config ./configs/calendar_config.toml

migrate:
	CALENDAR_DB_USER=$(CALENDAR_DB_USER) CALENDAR_DB_PASSWORD=$(CALENDAR_DB_PASSWORD) migrate -path ./migrations -database 'postgres://$(CALENDAR_DB_USER):$(CALENDAR_DB_PASSWORD)@localhost:5432/calendar_db?sslmode=disable' up

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.53.3

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint

# local tests
db:
	docker run -d --name pg -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=1234 -e POSTGRES_DB=calendar_db -p 5432:5432 postgres:latest

d:
	docker stop pg && docker rm pg