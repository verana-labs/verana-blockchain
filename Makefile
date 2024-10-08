BINARY_NAME=verana
MAIN_PATH=./cmd/verana/main.go
VERSION := $(shell git describe --tags)
COMMIT := $(shell git log -1 --format='%H')
GOBIN = $(shell go env GOPATH)/bin

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=verana \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=verana \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT)

BUILD_FLAGS := -ldflags '$(ldflags)'

.PHONY: all build install clean

all: install

build:
	@echo "Building Verana binary..."
	@go build $(BUILD_FLAGS) -o $(GOBIN)/$(BINARY_NAME) $(MAIN_PATH)

install: build
	@echo "Verana binary installed at: $(GOBIN)/$(BINARY_NAME)"

clean:
	@echo "Removing Verana binary..."
	@rm -f $(GOBIN)/$(BINARY_NAME)