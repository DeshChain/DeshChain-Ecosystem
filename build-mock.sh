#!/bin/bash
# Build mock binary for testing
mkdir -p build
export PATH=/usr/local/go/bin:$PATH
go build -o build/deshchaind ./cmd/deshchaind/mock.go