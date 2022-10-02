.PHONY: install build all lint run

all: build coverage

install:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.49.0

build:
	go build ./cmd/listener

test:
	go test ./...

coverage:
	go test -v -coverprofile cover.out ./...
	go tool cover -html=cover.out -o cover.html
	open cover.html

coverage-ci:
	go test -v -coverprofile cover.out ./...
	go tool cover -func=cover.out

lint:
	./bin/golangci-lint run --timeout 5m

ci: install lint build test coverage-ci

run: build
	./listener

run-debug: build
	./listener
