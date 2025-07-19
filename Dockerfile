# Copyright 2024 DeshChain Foundation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# DeshChain Node Dockerfile
# Multi-stage build for optimized production image

# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache \
    git \
    build-base \
    linux-headers \
    ca-certificates

# Set working directory
WORKDIR /src

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w -X github.com/cosmos/cosmos-sdk/version.Name=deshchain \
    -X github.com/cosmos/cosmos-sdk/version.AppName=deshchaind \
    -X github.com/cosmos/cosmos-sdk/version.Version=$(git describe --tags 2>/dev/null || echo 'development') \
    -X github.com/cosmos/cosmos-sdk/version.Commit=$(git rev-parse HEAD 2>/dev/null || echo 'unknown')" \
    -o /bin/deshchaind ./cmd/deshchaind

# Production stage
FROM alpine:3.18

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    jq \
    curl \
    bash \
    vim \
    tzdata \
    su-exec

# Create deshchain user
RUN addgroup -g 1000 deshchain && \
    adduser -D -u 1000 -G deshchain deshchain

# Copy binary from builder
COPY --from=builder /bin/deshchaind /usr/local/bin/deshchaind

# Create directories
RUN mkdir -p /home/deshchain/.deshchain/config \
    /home/deshchain/.deshchain/data \
    /home/deshchain/.deshchain/cosmovisor/genesis/bin \
    /home/deshchain/.deshchain/cosmovisor/upgrades

# Set ownership
RUN chown -R deshchain:deshchain /home/deshchain

# Copy entrypoint script
COPY docker/entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

# Expose ports
EXPOSE 26656 26657 26658 1317 9090 9091

# Set environment variables
ENV DAEMON_NAME=deshchaind
ENV DAEMON_HOME=/home/deshchain/.deshchain
ENV DAEMON_RESTART_AFTER_UPGRADE=true
ENV DAEMON_ALLOW_DOWNLOAD_BINARIES=false
ENV UNSAFE_SKIP_BACKUP=false

# Health check
HEALTHCHECK --interval=30s --timeout=10s --retries=3 \
    CMD curl -f http://localhost:26657/health || exit 1

# Set working directory
WORKDIR /home/deshchain

# Use entrypoint script
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]

# Default command
CMD ["deshchaind", "start"]