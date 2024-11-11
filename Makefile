BINARY_NAME=veranad
MAIN_PATH=./cmd/veranad/main.go
VERSION := $(shell git describe --tags)
COMMIT := $(shell git log -1 --format='%H')
GOBIN = $(shell go env GOPATH)/bin

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=verana \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=verana \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT)

BUILD_FLAGS := -ldflags '$(ldflags)'

.PHONY: all build install clean test test-verbose test-coverage

all: install

build:
	@echo "Building Veranad binary..."
	@go build $(BUILD_FLAGS) -o $(GOBIN)/$(BINARY_NAME) $(MAIN_PATH)

install: build
	@echo "Veranad binary installed at: $(GOBIN)/$(BINARY_NAME)"

clean:
	@echo "Removing Veranad binary..."
	@rm -f $(GOBIN)/$(BINARY_NAME)

test:
	@echo "Running tests..."
	@go test ./x/trustregistry/keeper/... ./x/diddirectory/keeper/...

test-verbose:
	@echo "Running tests with verbose output..."
	@go test -v ./x/trustregistry/keeper/... ./x/diddirectory/keeper/...

test-coverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./x/trustregistry/keeper/... ./x/diddirectory/keeper/...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"