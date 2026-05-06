.PHONY: build build-windows build-linux build-macos build-macos-arm clean test lint build-all all

BINARY_NAME=auraspeed
VERSION?=dev
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildTime=$(BUILD_TIME)"

build:
	go build $(LDFLAGS) -o $(BINARY_NAME) ./cmd/main.go

build-windows:
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME).exe ./cmd/main.go

build-linux:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-linux-amd64 ./cmd/main.go

build-macos:
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-darwin-amd64 ./cmd/main.go

build-macos-arm:
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BINARY_NAME)-darwin-arm64 ./cmd/main.go

build-all: build-windows build-linux build-macos build-macos-arm

all: build-all

clean:
	@if exist $(BINARY_NAME) del /F /Q $(BINARY_NAME)
	@if exist $(BINARY_NAME).exe del /F /Q $(BINARY_NAME).exe
	@for %%f in ($(BINARY_NAME)-*) do @if exist %%f del /F /Q %%f

test:
	go test ./... -v -cover

lint:
	golangci-lint run ./...
