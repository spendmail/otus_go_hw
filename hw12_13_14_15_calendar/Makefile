BIN := "./bin/calendar"
BIN_SCHEDULER := "./bin/calendar_scheduler"
BIN_SENDER := "./bin/calendar_sender"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(BIN_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/calendar_scheduler
	go build -v -o $(BIN_SENDER) -ldflags "$(LDFLAGS)" ./cmd/calendar_sender

run: build
	$(BIN) -config ./configs/calendar.toml

run-scheduler: build
	$(BIN_SCHEDULER) -config ./configs/calendar_scheduler.toml

run-sender: build
	$(BIN_SENDER) -config ./configs/calendar_sender.toml

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
	go test -race ./internal/... ./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.37.0

lint: install-lint-deps
	golangci-lint run ./...

migrate:
	(which goose > /dev/null) || go get -u github.com/pressly/goose/cmd/goose
	goose -dir migrations postgres "user=calendar password=calendar dbname=calendar sslmode=disable" up

generate:
	go generate ./...

up:
	docker-compose up -d

upl:
	docker-compose up

down:
	docker-compose down

rebuild:
	docker-compose up --build

bdd-clean:
	docker-compose -f docker-compose.test.yaml down \
    --rmi local \
		--volumes \
		--remove-orphans \
		--timeout 60; \
  	docker-compose rm -f

bdd:
	set -e ;\
	docker-compose -f docker-compose.test.yaml up --build -d ;\
	test_status_code=0 ;\
	docker-compose -f docker-compose.test.yaml run integration_tests go test -v || test_status_code=$$? ;\
	docker-compose -f docker-compose.test.yaml down ;\
	exit $$test_status_code ;

.PHONY: build run build-img run-img version test lint
