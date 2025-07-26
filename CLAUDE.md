# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

DeshChain is a comprehensive blockchain platform built on Cosmos SDK that combines:
- Cultural preservation through blockchain technology
- Revolutionary payment system (DeshPay) 
- **Gram Pension Scheme**: Revolutionary blockchain pension system with minimum 8% guaranteed returns, up to 50% based on DeshChain platform performance
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

## Memories

### Development Insights
- The project is an ambitious blend of cultural preservation and blockchain technology
- Key focus on creating a technology platform that serves social and cultural purposes
- Emphasis on transparency, accessibility, and social responsibility in blockchain development
- Strong focus on financial inclusion with culturally-relevant products
- Anti-fraud and investor protection are core priorities
- Traditional Indian financial concepts reimagined for blockchain (money orders, pension schemes)
- Community-driven governance with cultural values integration

### Major Implementation Milestones
- **Phase 1 Foundation (Complete)**: Cosmos SDK fork setup, NAMO token implementation, Cultural heritage system, Basic donation tracking, Blockchain explorer, Smart tax system foundation
- **Founder Protection Framework (July 18, 2025)**: Comprehensive governance protection with immutable 10% allocation, 0.10% tax royalty, 5% platform royalty, all perpetual and inheritable
- **Gamification Module (July 18, 2025)**: Bollywood-style developer achievement system with Bug Buster Bahubali, Feature Khan, social media integration, and movie poster generation
- **Revenue & Royalty System (July 18, 2025)**: Dual-stream revenue model with automated distribution, inheritance mechanism, and comprehensive tracking
- **Social Impact Revolution (July 18, 2025)**: Enhanced NGO donation system with 40% of all fees going to charity (0.75% from tax + 10% from platform revenues)
- **Batua Wallet Implementation (July 17, 2025)**: Complete native Flutter wallet with NAMO token pre-integration, Gram Pension Scheme access, Krishi Mitra preview, cultural UI components, HD wallet support, and secure storage
- **KYC Research & Strategy (July 17, 2025)**: Comprehensive research on open source KYC solutions, with Hyperledger Aries + Indy recommended as the primary blockchain-native identity solution for pension and agriculture finance
- **Sustainability Analysis (July 18, 2025)**: Comprehensive platform and founder sustainability projections showing ₹100,000+ Cr lifetime revenue potential and generational wealth creation
- **Production-Ready Identity Module (July 26, 2025)**: Complete W3C DID-compliant identity system with Verifiable Credentials, Zero-Knowledge Proofs, India Stack integration (Aadhaar, DigiLocker, UPI), biometric authentication, consent management, and backward-compatible integration adapters for existing KYC/biometric systems

### Technical Architecture Achievements
- **Multi-Module Ecosystem**: Governance, Explorer, Revenue, Royalty, Gamification, Tax, Cultural, Donation, NAMO, Identity, DUSD modules
- **Proto Definitions**: Complete gRPC and REST API coverage for all modules
- **Advanced Security**: Multi-signature wallets, encrypted storage, immutable parameters, W3C DID-based identity
- **Cultural Integration**: 10,000+ quotes, 22 language support, festival themes, patriotism scoring
- **Developer Experience**: Bollywood-style gamification, achievement system, social media integration
- **Identity Infrastructure**: W3C DID compliance, Verifiable Credentials, Zero-Knowledge Proofs, India Stack integration, biometric authentication
- **Backward Compatibility**: Seamless integration adapters for existing KYC and biometric systems

### Community-First Approach
- **Reduced founder allocation**: From 20% to 10% for community trust
- **Increased social impact**: 40% of all fees to charity vs industry standard 0-5%
- **Enhanced liquidity**: 20% allocation for market stability
- **Community rewards**: 15% allocation for user engagement
- **Transparent governance**: All founder actions visible on-chain with community oversight

### Innovation Highlights
- **First blockchain with 40% charity allocation**: Revolutionary social impact model
- **Perpetual inheritable royalties**: Unique in crypto space
- **Cultural blockchain**: Only platform integrating Indian heritage with finance
- **Bollywood gamification**: Entertainment industry themes in developer experience
- **Village panchayat KYC**: Traditional governance meets blockchain verification
- **Guaranteed pension returns**: Minimum 8% guaranteed, up to 50% returns based on DeshChain platform performance
- **Agricultural finance**: 6-9% interest rates vs 12-18% traditional banks

### Long-term Vision
- **50-year commitment**: Building sustainable technology for generations
- **Cultural preservation**: Digitizing and protecting Indian heritage through blockchain
- **Financial empowerment**: Creating opportunities for all Indians through accessible DeFi
- **Technological leadership**: Establishing India as a blockchain superpower
- **Global influence**: Spreading Indian values through technology worldwide

A community-first memory that emphasizes holistic technology development with deep cultural roots, showcasing how blockchain can be a tool for social transformation, not just financial innovation. The platform combines the best of traditional Indian values with cutting-edge blockchain technology to create a sustainable ecosystem that benefits all stakeholders while preserving cultural heritage for future generations.

- Added memory to reflect the project's commitment to adding depth and context to the development process