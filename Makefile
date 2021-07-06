PROJECT_NAME := "fsEngine"
PKG := "github.com/fanap-infra/fsEngine"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

.PHONY: all dep build clean test coverage coverhtml lint-test unit-test

all: build

lint-test:dep
	@golangci-lint run

fmt-test:
	@golangci-lint run -p format

fmt:
	@go fmt ./...

test: export EXEC_MODE = TEST

test: unit-test race-test

unit-test: dep
	@go test -count=1 -p=1 -short ./...

race-test: dep
	@go test -race -count=1 -p=1 -short ./...

# msan: dep
# 	@go test -msan -short ${PKG_LIST}

# coverage:
# 	./tools/coverage.sh;

# coverhtml:
# 	./tools/coverage.sh html;

dep:
	@go mod init; \
	go mod tidy; \
	go mod download; \
	go mod verify;