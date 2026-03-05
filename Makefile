# Project variables
PROJECT_NAME := broker-moc

# Build variables
.DEFAULT_GOAL = test
BUILD_DIR := build
TOOLS_DIR := $(shell pwd)/tools
TOOLS_BIN_DIR := ${TOOLS_DIR}/bin
DEV_GOARCH := $(shell go env GOARCH)
DEV_GOOS := $(shell go env GOOS)
EXE =
ifeq ($(DEV_GOOS),windows)
EXE = .exe
endif


## tools: Install required tooling.
.PHONY: tools
tools:
	@echo "==> Installing required tooling..."
	@cd $(TOOLS_DIR) && GOBIN=$(TOOLS_BIN_DIR) go install github.com/git-chglog/git-chglog/cmd/git-chglog
	@cd $(TOOLS_DIR) && GOBIN=$(TOOLS_BIN_DIR) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint

# clean: Delete the build directory.
.PHONY: clean
clean:
	@echo "==> Removing '$(BUILD_DIR)' directory..."
	@rm -rf $(BUILD_DIR)

## lint: Lint code with golangci-lint.
.PHONY: lint
lint:
	@echo "==> Linting code with 'golangci-lint'..."
	@$(TOOLS_BIN_DIR)/golangci-lint run

## test: Run all unit tests.
.PHONY: test
test:
	@echo "==> Running unit tests..."
	@mkdir -p $(BUILD_DIR)
	@go test -count=1 -v -cover -coverprofile=$(BUILD_DIR)/coverage.out -parallel=4 ./...

## build: Build CLI binary for default local system's operating system and architecture.
.PHONY: build
build:
	@echo "==> Building CLI binary..."
	@echo "    running go build for GOOS=$(DEV_GOOS) GOARCH=$(DEV_GOARCH)"
	@go build -o $(BUILD_DIR)/$(PROJECT_NAME)$(EXE) cmd/broker-moc/main.go

help: Makefile
	@echo "Usage: make <command>"
	@echo ""
	@echo "Commands:"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
