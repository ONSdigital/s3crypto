audit:
	go list -json -m all | nancy sleuth
.PHONY: audit

build:
	go build ./...
.PHONY: build

test:
	@echo "No tests available"
.PHONY: build