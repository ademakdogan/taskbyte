.PHONY: build run test clean lint install fmt cover

APP_NAME := taskbyte
BUILD_DIR := bin
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
LDFLAGS := -ldflags "-X github.com/adem/taskbyte/internal/app.Version=$(VERSION) -X github.com/adem/taskbyte/internal/app.BuildTime=$(BUILD_TIME)"

build:
	@echo "Building $(APP_NAME) $(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) ./cmd/taskbyte

run: build
	@./$(BUILD_DIR)/$(APP_NAME)

install:
	go install $(LDFLAGS) ./cmd/taskbyte

test:
	go test -v -cover ./...

test-short:
	go test -short ./...

cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

fmt:
	gofmt -w .

clean:
	@rm -rf $(BUILD_DIR) coverage.out coverage.html
	@echo "Cleaned."

lint:
	golangci-lint run ./...
