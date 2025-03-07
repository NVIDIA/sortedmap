# Copyright (c) 2015-2021, NVIDIA CORPORATION.
# SPDX-License-Identifier: Apache-2.0

all: fmt build test

.PHONY: all bench build clean cover fmt lint-update lint test

bench:
	go test -bench .

build:
	go build .

clean:
	go clean -i .

cover:
	go test -cover .

fmt:
	go fmt .

lint-update:
	rm -f $(GOPATH)/bin/golangci-lint
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin latest

lint:
	golangci-lint run --config .golangci.yml .

test:
	go test .
