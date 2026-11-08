.PHONY: test

install: install-dependencies

install-dependencies:
	go mod download

test: unit-tests integration-tests

unit-tests:
	@echo "=== Starting Unit Tests ==="
	go test -v $$(go list ./... | grep -v /tests/integration)

integration-tests:
	@echo "=== Starting Integration Tests ==="
	go test -v ./tests/integration/...
