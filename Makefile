.PHONY: build run test clean lint

APP_NAME := taskbyte
BUILD_DIR := bin

build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/taskbyte

run: build
	@./$(BUILD_DIR)/$(APP_NAME)

test:
	go test -v -cover ./...

test-short:
	go test -short ./...

clean:
	@rm -rf $(BUILD_DIR)
	@echo "Cleaned."

lint:
	golangci-lint run ./...
