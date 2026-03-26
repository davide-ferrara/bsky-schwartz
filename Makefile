.PHONY: run build clean test

BIN_DIR := bin
PROJECT_DIR := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

run:
	cd $(PROJECT_DIR) && go run ./cmd/server

build:
	go build -o $(BIN_DIR)/bsky-schwartz ./cmd/server

clean:
	rm -rf $(BIN_DIR)
	rm -f feed_*.txt

test:
	go test -v ./pkg/schwartz/...
