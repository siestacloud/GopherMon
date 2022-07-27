.PHONY: build
build:
	go build -v ./cmd/agent
	go build -v ./cmd/server

test:
	go test -v ./cmd/agent
	go test -v ./cmd/server
	go test -v ./internal/agent
	go test -v ./internal/server.DEFAULF_GOAL := build

