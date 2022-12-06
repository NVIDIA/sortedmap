# Copyright (c) 2015-2021, NVIDIA CORPORATION.
# SPDX-License-Identifier: Apache-2.0

all: fmt build test

.PHONY: all bench build clean cover fmt test

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

test:
	go test .
