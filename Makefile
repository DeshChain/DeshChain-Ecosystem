#!/usr/bin/make -f

PACKAGES_SIMTEST=$(shell go list ./... | grep -v '/simulation')
VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')
LEDGER_ENABLED ?= true
SDK_PACK := $(shell go list -m github.com/cosmos/cosmos-sdk | sed  's/ /\@/g')
BINDIR ?= $(GOPATH)/bin
SIMAPP = ./app

# for dockerized protobuf tools
DOCKER := $(shell which docker)
PROTO_CONTAINER := cosmwasm/proto-builder:0.11.2
HTTPS_GIT := https://github.com/aeye-employed/Desh-Chain-The-Blockchain-of-India.git

export GO111MODULE = on

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

ifeq ($(LINK_STATICALLY),true)
	ldflags += -linkmode=external -extldflags "-Wl,-z,muldefs -static"
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'

all: install

install: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/deshchaind

build:
	go build $(BUILD_FLAGS) -o bin/deshchaind ./cmd/deshchaind

###############################################################################
###                                Protobuf                                 ###
###############################################################################

proto-all: proto-format proto-gen

proto-gen:
	@echo "Generating Protobuf files"
	@$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(PROTO_CONTAINER) \
		sh ./scripts/protocgen.sh

proto-format:
	@$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(PROTO_CONTAINER) \
		find ./ -name "*.proto" -exec clang-format -i {} \;

proto-swagger-gen:
	@./scripts/protoc-swagger-gen.sh

proto-lint:
	@$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(PROTO_CONTAINER) \
		buf lint --error-format=json

proto-check-breaking:
	@$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(PROTO_CONTAINER) \
		buf breaking --against $(HTTPS_GIT)#branch=main

###############################################################################
###                                 Devdoc                                  ###
###############################################################################

DEVDOC_SAVE = docker commit `docker ps -a -n 1 -q` devdoc:local

devdoc-init:
	docker run -it -v "$(CURDIR):/go/src/github.com/deshchain/deshchain" -w "/go/src/github.com/deshchain/deshchain" tendermint/devdoc echo
	# TODO make this safer
	$(call DEVDOC_SAVE)

devdoc:
	docker run -it -v "$(CURDIR):/go/src/github.com/deshchain/deshchain" -w "/go/src/github.com/deshchain/deshchain" devdoc:local bash

devdoc-save:
	# TODO make this safer
	$(call DEVDOC_SAVE)

devdoc-clean:
	docker rmi -f $$(docker images -f "dangling=true" -q)

devdoc-update:
	docker pull tendermint/devdoc

###############################################################################
###                                Localnet                                 ###
###############################################################################

# Run a single node
localnet: build
	./scripts/localnet.sh

###############################################################################
###                               Deployment                                ###
###############################################################################

init-testnet:
	./scripts/init-testnet.sh

init-mainnet:
	./scripts/init-mainnet.sh

start:
	./bin/deshchaind start

start-testnet:
	./bin/deshchaind start --home ~/.deshchain-testnet

start-mainnet:
	./bin/deshchaind start --home ~/.deshchain

deploy-testnet:
	@echo "Deploying to testnet..."
	./scripts/deploy-testnet.sh

deploy-mainnet:
	@echo "Deploying to mainnet..."
	@echo "WARNING: This will deploy to mainnet. Are you sure? [y/N]"
	@read -r response; \
	if [ "$$response" = "y" ]; then \
		./scripts/deploy-mainnet.sh; \
	else \
		echo "Deployment cancelled."; \
	fi

###############################################################################
###                                  Test                                   ###
###############################################################################

test: test-unit test-build

test-unit:
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' -ldflags '$(ldflags)' ./...

test-race:
	@VERSION=$(VERSION) go test -mod=readonly -race -tags='ledger test_ledger_mock' ./...

test-cover:
	@go test -mod=readonly -timeout 30m -race -coverprofile=coverage.txt -covermode=atomic -tags='ledger test_ledger_mock' ./...

test-build: build
	@go test -mod=readonly -p 4 -tags='ledger test_ledger_mock' ./...

test-integration:
	@go test -mod=readonly -tags='ledger test_ledger_mock integration' ./integration_test.go

test-e2e:
	@go test -mod=readonly -tags='ledger test_ledger_mock e2e' ./e2e/...

test-load:
	@echo "Running basic load test..."
	@go run scripts/load-testing/load-test.go -workers=10 -tx-per-worker=100

test-stress:
	@echo "Running stress test suite..."
	@./scripts/load-testing/stress-test.sh

test-benchmark:
	@echo "Running comprehensive benchmark suite..."
	@./scripts/load-testing/benchmark-suite.sh

test-performance:
	@echo "Starting performance monitoring (60s)..."
	@python3 scripts/load-testing/performance-monitor.py --duration=60 --interval=10

# Quick performance validation
test-perf-quick:
	@echo "Running quick performance validation..."
	@go run scripts/load-testing/load-test.go -workers=5 -tx-per-worker=50 -output=quick-perf.json
	@echo "Quick performance test completed. Results in quick-perf.json"

# CI/CD performance gate
test-perf-gate:
	@echo "Running performance gate for CI/CD..."
	@go run scripts/load-testing/load-test.go -workers=20 -tx-per-worker=100 -output=ci-perf.json
	@SUCCESS_RATE=$$(jq -r '.results.success_rate' ci-perf.json); \
	if [ "$$(echo "$$SUCCESS_RATE < 95" | bc -l)" -eq 1 ]; then \
		echo "Performance gate failed: Success rate $$SUCCESS_RATE% < 95%"; \
		exit 1; \
	else \
		echo "Performance gate passed: Success rate $$SUCCESS_RATE%"; \
	fi

test-sim-import-export: runsim
	@echo "Running application import/export simulation. This may take several minutes..."
	@$(BINDIR)/runsim -Jobs=4 -SimAppPkg=$(SIMAPP) -ExitOnFail 50 5 TestAppImportExport

test-sim-multi-seed-short: runsim
	@echo "Running short multi-seed application simulation. This may take awhile!"
	@$(BINDIR)/runsim -Jobs=4 -SimAppPkg=$(SIMAPP) -ExitOnFail 50 5 TestFullAppSimulation

benchmark:
	@go test -mod=readonly -bench=. ./...

###############################################################################
###                                Linting                                  ###
###############################################################################

lint:
	golangci-lint run --out-format=tab

lint-fix:
	golangci-lint run --fix --out-format=tab --issues-exit-code=0

# Security audit
audit:
	@echo "Running security audit..."
	@./scripts/security/audit.sh

# Network security scan
network-scan:
	@echo "Running network security scan..."
	@./scripts/security/network-scan.sh localhost

# Comprehensive security check
security-check: audit network-scan
	@echo "Security checks completed"

# Clean build artifacts
clean:
	rm -rf bin/ build/ artifacts/

# Format code
format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -name "*.pb.go" -not -name "*.pb.gw.go" | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -name "*.pb.go" -not -name "*.pb.gw.go" | xargs goimports -w -local github.com/deshchain/namo

###############################################################################
###                            Optimization                                ###
###############################################################################

# Complete optimization suite
optimize:
	@echo "Running complete optimization analysis..."
	@./scripts/optimization/run-all-optimizations.sh

# Individual optimization tools
optimize-memory:
	@echo "Running memory optimization analysis..."
	@./scripts/optimization/memory-optimizer.sh

optimize-cpu:
	@echo "Running CPU performance profiling..."
	@./scripts/optimization/cpu-profiler.sh

optimize-blockchain:
	@echo "Running blockchain optimization analysis..."
	@python3 scripts/optimization/blockchain-optimizer.py --duration=300

optimize-application:
	@echo "Running application profiling..."
	@go run scripts/optimization/performance-profiler.go -duration=5m

# Quick optimization check
optimize-quick:
	@echo "Running quick optimization check..."
	@ANALYSIS_DURATION=120 ./scripts/optimization/run-all-optimizations.sh

.PHONY: all build install proto-gen proto-format proto-lint test test-unit test-race test-cover test-build test-integration test-e2e test-load test-stress test-benchmark test-performance test-perf-quick test-perf-gate benchmark lint lint-fix audit network-scan security-check clean format optimize optimize-memory optimize-cpu optimize-blockchain optimize-application optimize-quick

###############################################################################
###                                 Devdoc                                  ###
###############################################################################

contract-tests:
	@go test -mod=readonly -v ./x/wasm/...

###############################################################################
###                                Releasing                                ###
###############################################################################

GORELEASER_DEBUG := false

ifdef GITHUB_TOKEN
release:
	docker run \
		--rm \
		--env GITHUB_TOKEN=$(GITHUB_TOKEN) \
		--env GORELEASER_CURRENT_TAG=$(VERSION) \
		--env GORELEASER_DEBUG=$(GORELEASER_DEBUG) \
		--env COSMOVISOR_ENABLED=$(COSMOVISOR_ENABLED) \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v $(CURDIR):/workspace \
		-w /workspace \
		goreleaser/goreleaser-cross:latest \
		release \
		--clean
else
release:
	@echo "Error: GITHUB_TOKEN not defined. Please define it before running 'make release'."
endif

.PHONY: release