.PHONY: build install

BIN_DIR ?= $(shell go env GOPATH)/bin

build install:
	mkdir -p "$(BIN_DIR)"
	go build -o "$(BIN_DIR)/lsproj" .
