#!/usr/bin/make -f

PACKAGES_SIMTEST=$(shell go list ./... | grep '/simulation')
VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')
LEDGER_ENABLED ?= true
SDK_PACK := $(shell go list -m github.com/cosmos/cosmos-sdk | sed  's/ /\@/g')
BUILDDIR ?= $(CURDIR)/build
DOCKER := $(shell which docker)
BINDIR ?= $(GOPATH)/bin
SIMAPP = ./app

# process build tags
build_tags = netgo
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif

ifeq (cleveldb,$(findstring cleveldb,$(DESHCHAIN_BUILD_OPTIONS)))
  build_tags += gcc cleveldb
endif

ifeq (rocksdb,$(findstring rocksdb,$(DESHCHAIN_BUILD_OPTIONS)))
  build_tags += gcc rocksdb
endif

build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# process linker flags
ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=deshchain \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=deshchaind \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)"

ifeq (cleveldb,$(findstring cleveldb,$(DESHCHAIN_BUILD_OPTIONS)))
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb
endif
ifeq (rocksdb,$(findstring rocksdb,$(DESHCHAIN_BUILD_OPTIONS)))
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=rocksdb
endif
ifeq (,$(findstring nostrip,$(DESHCHAIN_BUILD_OPTIONS)))
  ldflags += -w -s
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'

# check for nostrip option
ifeq (,$(findstring nostrip,$(DESHCHAIN_BUILD_OPTIONS)))
  BUILD_FLAGS += -trimpath
endif

###############################################################################
###                                  Build                                  ###
###############################################################################

all: install lint test

BUILD_TARGETS := build install

build: BUILD_ARGS=-o $(BUILDDIR)/

$(BUILD_TARGETS): go.sum $(BUILDDIR)/
	go $@ -mod=readonly $(BUILD_FLAGS) $(BUILD_ARGS) ./cmd/deshchaind

$(BUILDDIR)/:
	mkdir -p $(BUILDDIR)/

.PHONY: build install

###############################################################################
###                                Testing                                   ###
###############################################################################

test: test-unit
test-all: test-unit test-integration test-e2e

test-unit:
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' ./...

test-integration:
	@echo "Running integration tests..."
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' -timeout 30m ./tests/integration/...

test-e2e:
	@echo "Running e2e tests..."
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' -timeout 30m ./tests/e2e/...

test-load:
	@echo "Running load tests..."
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' -timeout 30m ./tests/load/...

.PHONY: test test-all test-unit test-integration test-e2e test-load

###############################################################################
###                                Linting                                  ###
###############################################################################

format:
	@echo "Formatting Go files..."
	@go fmt ./...

lint:
	@echo "Running linter..."
	@golangci-lint run --timeout=10m

lint-fix:
	@echo "Running linter with fix..."
	@golangci-lint run --fix --timeout=10m

.PHONY: format lint lint-fix

###############################################################################
###                                Docker                                   ###
###############################################################################

build-docker:
	@echo "Building Docker image..."
	@$(DOCKER) build -t deshchain:local .

build-docker-testnet:
	@echo "Building Docker image for testnet..."
	@$(DOCKER) build -f Dockerfile.testnet -t deshchain:testnet .

.PHONY: build-docker build-docker-testnet

###############################################################################
###                               Development                               ###
###############################################################################

init-testnet:
	@echo "Initializing testnet..."
	@./scripts/init-testnet.sh

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILDDIR)

.PHONY: init-testnet clean

###############################################################################
###                                 Cosmos                                  ###
###############################################################################

go.sum: go.mod
	@echo "Ensuring dependencies have not been modified..."
	@go mod verify

mod-tidy:
	@echo "Tidying go.mod..."
	@go mod tidy

.PHONY: go.sum mod-tidy