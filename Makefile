.PHONY: build build-windows build-linux build-macos build-macos-arm build-all run clean test lint format install release all

BINARY_NAME=auraspeed
VERSION?=$(shell cat .version 2>/dev/null || echo "dev")
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS=-ldflags "-s -w -X 'auraspeed/cmd/root.Version=$(VERSION)' -X 'auraspeed/cmd/root.Commit=$(COMMIT)' -X 'auraspeed/cmd/root.BuildTime=$(BUILD_TIME)'"

build:
	go build $(LDFLAGS) -o $(BINARY_NAME) ./cmd/main.go

build-windows:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(BINARY_NAME)-windows-amd64.exe ./cmd/main.go

build-linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(BINARY_NAME)-linux-amd64 ./cmd/main.go

build-macos:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(BINARY_NAME)-darwin-amd64 ./cmd/main.go

build-macos-arm:
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(BINARY_NAME)-darwin-arm64 ./cmd/main.go

build-all: build-windows build-linux build-macos build-macos-arm

run: build
	./$(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME).exe $(BINARY_NAME)-*

test:
	go test ./... -v -short -cover

lint:
	golangci-lint run ./...
	go vet ./...

format:
	go fmt ./...

install:
	go install $(LDFLAGS) ./cmd/main.go

release: build-all
	@echo "Release $(VERSION) binaries built in current directory"

all: build-all
