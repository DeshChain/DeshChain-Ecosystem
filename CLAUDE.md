# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

DeshChain is a comprehensive blockchain platform built on Cosmos SDK that combines:
- Cultural preservation through blockchain technology
- Revolutionary payment system (DeshPay) 
- **Gram Pension Scheme**: Revolutionary blockchain pension system with 50% guaranteed returns
- **Sikkebaaz**: Desi memecoin launchpad with anti-pump & dump protection
- **Money Order DEX**: Culturally-rooted decentralized exchange inspired by traditional money orders
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
- **Native Wallet**: Flutter (Batua Wallet)

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

# Batua Wallet Development
cd batua/mobile
flutter pub get
flutter run
flutter build apk
flutter build ios

# KYC System Development
cd kyc
docker-compose up -d  # Start Hyperledger Aries agents
python -m aries_cloudagent.app --help  # ACA-Py commands
npm run kyc-schema-gen  # Generate credential schemas
npm run kyc-test  # Run KYC integration tests

# Deployment
make build-docker
make deploy-testnet
make deploy-mainnet
```

(Rest of the existing file content remains the same)

## Memories

### Development Insights
- The project is an ambitious blend of cultural preservation and blockchain technology
- Key focus on creating a technology platform that serves social and cultural purposes
- Emphasis on transparency, accessibility, and social responsibility in blockchain development
- Strong focus on financial inclusion with culturally-relevant products
- Anti-fraud and investor protection are core priorities
- Traditional Indian financial concepts reimagined for blockchain (money orders, pension schemes)
- Community-driven governance with cultural values integration

- Add to memory - a simple yet significant community-driven technical approach that prioritizes social impact and cultural preservation
- **Batua Wallet Implementation (July 17, 2025)**: Complete native Flutter wallet with NAMO token pre-integration, Gram Pension Scheme access, Krishi Mitra preview, cultural UI components, HD wallet support, and secure storage
- **KYC Research & Strategy (July 17, 2025)**: Comprehensive research on open source KYC solutions, with Hyperledger Aries + Indy recommended as the primary blockchain-native identity solution for pension and agriculture finance
- A community-first memory that emphasizes holistic technology development with deep cultural roots
- Add to memory - an implementation that showcases how blockchain can be a tool for social transformation, not just financial innovation
- We will continue tomorrow