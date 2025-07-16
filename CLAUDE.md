# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

DeshChain is a comprehensive blockchain platform built on Cosmos SDK that combines:
- Cultural preservation through blockchain technology
- Revolutionary payment system (DeshPay) 
- Transparent donation tracking
- NFT-based rewards and gamification
- Privacy-preserving transactions
- Dynamic tax system with volume-based reduction
- Decentralized governance

## Technical Architecture

### Core Stack
- **Blockchain**: Cosmos SDK (Go)
- **Consensus**: Tendermint
- **Frontend**: React/TypeScript + PWA
- **Backend**: Go, Node.js/TypeScript
- **Database**: PostgreSQL, Redis
- **Storage**: IPFS (cultural content, NFT metadata)
- **Privacy**: zk-SNARKs implementation
- **Mobile**: React Native / PWA

### Development Commands

```bash
# Blockchain Development
go mod init deshchain
go mod tidy
go build ./cmd/deshchaind
go test ./...

# Frontend Development
npm install
npm start
npm run build
npm run test

# Testing
make test-unit
make test-integration
make test-e2e
make test-load

# Deployment
make build-docker
make deploy-testnet
make deploy-mainnet
```

[... rest of the existing content remains unchanged ...]

## Memories

### Development Insights
- The project is an ambitious blend of cultural preservation and blockchain technology
- Key focus on creating a technology platform that serves social and cultural purposes
- Emphasis on transparency, accessibility, and social responsibility in blockchain development