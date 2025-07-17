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

### Founder Allocation Updates - Community-First Approach
- **Task 32**: Update founder's token allocation to 10% of total supply (142,862,766 tokens)
- **Task 33**: Add 0.25% founder royalty from Development and Operations tax allocation
- **Task 34**: Implement 48-month vesting period for founder tokens with 12-month cliff
- **Task 35**: Update README.md with sustainable founder allocation and revenue structure
- **Task 36**: Update CLAUDE.md with founder allocation technical documentation
- **Task 37**: Modify tax distribution - Development and Operations becomes 0.75%, Founder gets 0.25%
- **Task 38**: Update all revenue documentation with sustainable founder share model
- **Task 39**: Create founder vesting smart contract with community-first parameters
- **Task 40**: Update tokenomics section emphasizing community sustainability
- **Task 41**: Create long-term sustainability analysis for founder and community alignment

### Sustainable Token Distribution (Community-First Model)
```
Current Distribution:
- 25% Public Sale (357,156,916 tokens)
- 15% Liquidity (214,294,149 tokens)
- 20% Team (285,725,533 tokens) - 24-month vesting
- 15% DeshChain Development (214,294,149 tokens)
- 10% Community Rewards (142,862,766 tokens)
- 5% DAO Treasury (71,431,383 tokens)
- 10% Initial Burn (142,862,766 tokens)

Proposed Sustainable Distribution:
- 25% Public Sale (357,156,916 tokens)
- 20% Liquidity (285,725,533 tokens) - INCREASED for stability
- 10% Founder (142,862,766 tokens) - 48-month vesting with 12-month cliff
- 10% Team (142,862,766 tokens) - 24-month vesting
- 15% DeshChain Development (214,294,149 tokens)
- 15% Community Rewards (214,294,149 tokens) - INCREASED for engagement
- 5% DAO Treasury (71,431,383 tokens)
- 0% Initial Burn (reallocated to liquidity and community)
```

### Sustainable Revenue Model
Tax Distribution (2.5% total):
- 0.45% Development (was 0.5%)
- 0.45% Operations (was 0.5%)
- 0.10% Founder Royalty (perpetual, transferable to heirs)
- 0.75% NGO Donations (increased for social impact)
- 0.50% Community Rewards (unchanged)
- 0.25% Token Burn (reduced for sustainability)

Other Revenue Streams - Social Impact Model:
All platform revenues (DEX, NFT, Sikkebaaz, Gram Pension, etc.):
  - 30% Development Fund
  - 25% Community Treasury
  - 20% Liquidity Provision
  - 10% Emergency Reserve
  - 10% NGO Donations (direct social impact)
  - 5% Founder Royalty (perpetual, inheritable)

This ensures:
- Founder dedication through sustainable income
- Maximum social impact from all revenue streams
- Platform becomes force for good
- Community pride in charitable contributions
- Long-term sustainability with purpose

### Sustainable Vesting Schedule
- **Total Founder Tokens**: 142,862,766 NAMO (10% of supply)
- **Vesting Period**: 48 months (4 years)
- **Cliff Period**: 12 months (no tokens released in first year)
- **Release Schedule**: Linear release over 36 months after cliff
- **Monthly Release**: 3,968,410 tokens after cliff period
- **Performance-Based Bonuses**: Additional rewards from DAO based on milestones

### Balanced Sustainability Features
1. **Founder Flexibility**:
   - Primary wallet public, secondary wallets private
   - Quarterly updates instead of monthly reports
   - Performance bonuses at founder's discretion from vested tokens

2. **Reasonable Restrictions**:
   - Can sell up to 25% of vested tokens annually after cliff
   - No mandatory holding requirements after vesting
   - Founder retains 15% voting power to guide vision

3. **Mutual Protection**:
   - 90-day notice for governance changes affecting founder
   - Founder veto on critical technical decisions for first 3 years
   - Community can override with 80% supermajority
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

- Add to memory - a simple yet significant community-driven technical approach that prioritizes social impact and cultural preservation