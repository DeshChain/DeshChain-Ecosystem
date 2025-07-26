#!/bin/bash
set -e

echo "Fixing Go dependencies..."

# Update module path
go mod edit -module github.com/deshchain/deshchain

# Add missing dependencies
go get cosmossdk.io/store@v1.1.0
go get cosmossdk.io/x/evidence@v0.1.0
go get cosmossdk.io/x/feegrant@v0.1.0
go get cosmossdk.io/x/upgrade@v0.1.1
go get cosmossdk.io/collections@v0.4.0
go get cosmossdk.io/errors@v1.0.1
go get cosmossdk.io/core@v0.11.0
go get github.com/cometbft/cometbft@v0.38.7
go get github.com/cometbft/cometbft-db@v0.9.1
go get github.com/cosmos/cosmos-sdk@v0.50.7
go get github.com/cosmos/gogoproto@v1.4.12
go get github.com/cosmos/ibc-go/v8@v8.2.0
go get github.com/prometheus/client_golang@v1.19.0
go get github.com/gogo/protobuf@v1.3.2

# Create missing directories
mkdir -p x/nft/keeper
mkdir -p x/sikkebaaz
mkdir -p x/shikshamitra/client/rest
mkdir -p x/vyavasayamitra/client/rest

# Create placeholder files
touch x/nft/keeper/keeper.go
touch x/sikkebaaz/module.go
touch x/shikshamitra/client/rest/rest.go
touch x/vyavasayamitra/client/rest/rest.go

# Run go mod tidy
go mod tidy