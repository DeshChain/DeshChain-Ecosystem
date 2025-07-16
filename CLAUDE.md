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

## Key Features Documentation

### Gram Pension Scheme
- **Location**: `/root/namo/proto/deshchain/grampension/v1/`
- **Purpose**: Revolutionary blockchain pension system with guaranteed 50% returns
- **Implementation**: Complete proto definitions with 12 transaction handlers and 21 query endpoints
- **Financial Model**: 80.6% profit margin with comprehensive sustainability analysis
- **Features**: KYC integration, referral rewards, loyalty programs, risk management
- **Cultural Integration**: Patriotism scoring, cultural engagement bonuses

### Sikkebaaz - Desi MemeCoin Launchpad  
- **Location**: `/root/namo/sikkebaaz/`
- **Purpose**: India's first safety-first memecoin launchpad with anti-pump & dump protection
- **Safety Features**: 72-hour community review, mandatory KYC, liquidity locks, progressive releases
- **Cultural Themes**: Bollywood categories, festival launches, regional customization
- **Protection**: Investor protection fund, guaranteed liquidity, reputation scoring
- **Community**: Traditional joint family tokenomics, cultural meme contests

### Money Order DEX
- **Location**: `/root/namo/money-order/`
- **Purpose**: Culturally-rooted DEX inspired by traditional Indian money orders
- **Emotional Connection**: Recreates trust and reliability of traditional money orders
- **Features**: Multi-chain support, cultural themes, festival-based trading events
- **Languages**: 22 Indian languages support with regional customization
- **Trading**: AMM, order books, derivatives, yield farming with cultural integration

### Kisaan Mitra - Agricultural Lending Platform
- **Location**: `/root/namo/kisaan-mitra/`
- **Purpose**: India's first community-backed agricultural lending platform for farmers
- **Interest Rates**: 6-9% compared to 12-18% from commercial banks
- **Protection**: Triple-layer fraud protection with village-level verification
- **Community**: Village panchayat verification system with peer monitoring
- **Features**: Seasonal lending, crop-specific loans, regional customization
- **Impact**: 1,00,000+ farmers served, ₹500 crores disbursed, 95% repayment rate
- **Cultural Integration**: Festival-based lending, traditional wisdom integration

### Smart Tax System Foundation
- **Location**: `/root/namo/x/tax/types/`
- **Implementation**: Complete tax calculation engine with progressive rates
- **Features**: Volume-based discounts, patriotism bonuses, cultural engagement rewards
- **Protection**: ₹1,000 daily cap, donation exemptions, optimization algorithms
- **Components**: TaxCalculator, TaxOptimizer, ComplianceChecker, TaxReporter

## Module Structure

### Core Modules
- **NAMO Token**: `/root/namo/x/namo/` - Native token with cultural tokenomics
- **Cultural Heritage**: `/root/namo/x/cultural/` - 10,000+ quotes and cultural content
- **Donation Tracking**: `/root/namo/x/donation/` - Transparent multi-signature NGO wallets
- **Blockchain Explorer**: `/root/namo/x/explorer/` - Real-time indexing and search
- **Tax System**: `/root/namo/x/tax/` - Progressive tax with volume discounts
- **Gram Pension**: `/root/namo/x/grampension/` - Blockchain pension platform

### Proto Definitions
- All modules have comprehensive proto definitions in `/root/namo/proto/deshchain/`
- Message types for transactions, queries, and responses
- gRPC service definitions for all operations
- Extensive error handling and validation

## Development Guidelines

### Cultural Integration
- Every financial transaction includes cultural quotes from Indian leaders
- Festival-based features and timing (Diwali, Holi, Dussehra, etc.)
- Regional customization for different Indian states
- Traditional values integrated with modern blockchain technology

### Safety & Security
- Multi-layer security with formal verification
- Anti-pump & dump mechanisms in all financial products
- Investor protection funds and insurance coverage
- 24/7 monitoring and fraud detection

### User Experience
- Mobile-first design with PWA capabilities
- 22 Indian languages support
- Voice commands in local languages
- Simplified onboarding with cultural context

## Pending Development Tasks

### Founder Allocation Updates Required
- **Task 32**: Update founder's token allocation to 20% of total supply (285,725,533 tokens)
- **Task 33**: Add 15% founder royalty share to all revenue streams in tokenomics
- **Task 34**: Implement 60-month vesting period for founder tokens with proper schedule
- **Task 35**: Update README.md with new founder allocation and revenue structure
- **Task 36**: Update CLAUDE.md with founder allocation technical documentation
- **Task 37**: Modify tax distribution to include 15% founder royalty share
- **Task 38**: Update all revenue stream documentation with founder's 15% share
- **Task 39**: Create founder vesting smart contract specification
- **Task 40**: Update tokenomics section with new distribution percentages
- **Task 41**: Re-estimate founder's wealth projections based on new allocations

### Proposed Token Distribution Changes
```
Current Distribution:
- 25% Public Sale (357,156,916 tokens)
- 15% Liquidity (214,294,149 tokens)
- 20% Team (285,725,533 tokens) - 24-month vesting
- 15% DeshChain Development (214,294,149 tokens)
- 10% Community Rewards (142,862,766 tokens)
- 5% DAO Treasury (71,431,383 tokens)
- 10% Initial Burn (142,862,766 tokens)

Proposed Distribution:
- 25% Public Sale (357,156,916 tokens)
- 15% Liquidity (214,294,149 tokens)
- 20% Founder (285,725,533 tokens) - 60-month vesting
- 15% Team (214,294,149 tokens) - 24-month vesting
- 10% DeshChain Development (142,862,766 tokens)
- 10% Community Rewards (142,862,766 tokens)
- 5% DAO Treasury (71,431,383 tokens)
- 0% Initial Burn (redirected to founder)
```

### Revenue Stream Updates Required
All revenue streams need 15% founder royalty:
- Transaction Tax: 2.5% → 15% to founder
- Privacy Fees: ₹50-150 → 15% to founder
- DEX Trading Fees: 0.3% → 15% to founder
- NFT Marketplace: 2.5% → 15% to founder
- Sikkebaaz Launchpad: 100 NAMO + 2% → 15% to founder
- Gram Pension: 80.6% margin → 15% to founder
- Kisaan Mitra: 6-9% APR → 15% to founder
- All other revenue streams → 15% to founder

### Vesting Schedule Specification
- **Total Founder Tokens**: 285,725,533 NAMO (20% of supply)
- **Vesting Period**: 60 months (5 years)
- **Cliff Period**: 12 months (no tokens released in first year)
- **Release Schedule**: Linear release over 48 months after cliff
- **Monthly Release**: 5,952,615 tokens after cliff period

## Memories

### Development Insights
- The project is an ambitious blend of cultural preservation and blockchain technology
- Key focus on creating a technology platform that serves social and cultural purposes
- Emphasis on transparency, accessibility, and social responsibility in blockchain development
- Strong focus on financial inclusion with culturally-relevant products
- Anti-fraud and investor protection are core priorities
- Traditional Indian financial concepts reimagined for blockchain (money orders, pension schemes)
- Community-driven governance with cultural values integration