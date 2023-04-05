OS_ARCH=darwin_arm64
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
TEST?=$$(go list ./... | grep -v 'vendor')

.DEFAULT_GOAL := help

.PHONY: help run test fmt test_with_profile lint

default: run

help:
	@echo "Please use 'make <target>' where <target> is one of"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z\._-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run:
	go run cmd/server/main.go

lint: ## Lint soruce code
	@echo "==> Checking source code against linters..."
	@GOGC=30 golangci-lint run .
	goimports -w .

fmt: ## Format code
	gofmt -w $(GOFMT_FILES)

test:
	go test $(TEST) -cover -v  $(TESTARGS) || exit 1

test_with_profile:
	go test $(TEST) -coverprofile=coverage.out -v  $(TESTARGS)
	go tool cover -html=coverage.out
