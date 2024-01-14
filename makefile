# Makefile

# Load environment variables from .env file
include .env
export

.PHONY: run

run:
	go run cmd/api/server.go

.PHONY: test

test:
	go test -v ./...

.PHONY: build

build:
	go build -o bin/api cmd/api/server.go