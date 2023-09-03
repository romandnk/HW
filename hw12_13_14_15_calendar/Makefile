include ./configs/calendar.env

BIN_calendar := "./bin/calendar"
BIN_scheduler := "./bin/scheduler"
BIN_sender := "./bin/sender"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

up:
	docker compose -f ./deployments/docker-compose.yaml up --build -d

down:
	docker compose -f ./deployments/docker-compose.yaml down && \
	docker volume rm calendar_log calendar_pgdata

build-calendar:
	go build -v -o $(BIN_calendar) -ldflags "$(LDFLAGS)" ./cmd/calendar

build-calendar-scheduler:
	go build -o $(BIN_scheduler) ./cmd/calendar_scheduler

build-calendar-sender:
	go build -v -o $(BIN_sender) ./cmd/calendar_sender

build: build-calendar build-calendar-scheduler build-calendar-sender

run-calendar: build-calendar
	$(BIN_calendar) -config ./configs/calendar_config.toml

run-calendar-scheduler: build-calendar-scheduler
	$(BIN_scheduler) -config ./configs/scheduler_config.toml

run-calendar-sender: build-calendar-sender
	$(BIN_sender) -config ./configs/sender_config.toml

clean:
	rm -f $(BIN_calendar) $(BIN_scheduler) $(BIN_sender)

generate:
	go generate ./api/.

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

integration-tests:
	docker compose -f deployments/docker-compose-test.yaml up -d --build && \
	docker wait calendar_test-integration_test-1 && \
	docker compose -f deployments/docker-compose-test.yaml down && \
    docker volume rm calendar_test_pgdata_test

integration-tests-down:
	docker compose -f deployments/docker-compose-test.yaml down && \
    docker volume rm calendar_test_pgdata_test

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.53.3

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint