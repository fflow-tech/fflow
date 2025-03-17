.PHONY: all build test lint

SERVER    := faas
HTTP_PORT := 50031
USERNAME  := $(shell git config user.name)
TARNAME := $(addsuffix .tgz, $(SERVER))

GOPATH ?= $(shell go env GOPATH)
SERVER_BIN_PATH := ../../../../bin
GO ?= go
env ?= dev
BUILD_DATE ?= $(shell date -u)
BUILD_HASH ?= $(shell git rev-parse HEAD)
LDFLAGS += -X "model.BuildDate=$(BUILD_DATE)"
LDFLAGS += -X "model.BuildHash=$(BUILD_HASH)"
MAIN = main.go

all: lint build

build:
	@echo "make build"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -o $(SERVER) -ldflags '$(LDFLAGS)'
	mkdir -p $(SERVER_BIN_PATH)
	cp $(SERVER) $(SERVER_BIN_PATH)

test:
	@echo "make test"
	$(GO) test -v -race -cover -gcflags=all=-l -coverprofile=coverage.data ./...
	$(GO) tool cover -func=coverage.data -o coverage.txt
	rm -rf coverage.data

swag:
	@echo "make swag"
	# 解决swag引用不到pb的问题
	swag init --parseDependency --parseInternal --parseDepth 2 -g main.go
	@echo "http://127.0.0.1:$(HTTP_PORT)/$(SERVER)/swagger/index.html"

generate: #生成mock桩代码
	go generate ./...

lint:
	@echo "make lint"
	golangci-lint run -c .golangci.yml
