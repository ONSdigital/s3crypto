SHELL=bash

test:
	go test -count=1 -mod=mod -race -cover ./...

.PHONY: test

audit:
	go list -json -m all | nancy sleuth
.PHONY: audit

build:
	go build -mod=mod ./...
.PHONY: build

.PHONY: lint
lint:
	exit

