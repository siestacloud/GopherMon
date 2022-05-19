.PHONY: build
build:
	go build -v ./cmd/agent
	go build -v ./cmd/server

.DEFAULF_GOAL := build