# DeshChain Validator Staking & Security System

## Overview

DeshChain implements a revolutionary USD-pegged staking system with comprehensive security mechanisms to protect the network from validator token dumps while ensuring long-term commitment and alignment with platform success.

## Table of Contents

1. [USD-Pegged Staking Model](#usd-pegged-staking-model)
2. [Tiered Validator Structure](#tiered-validator-structure)
3. [Security Mechanisms](#security-mechanisms)
4. [Vesting & Lock Periods](#vesting--lock-periods)
5. [Slashing Framework](#slashing-framework)
6. [Circuit Breakers](#circuit-breakers)
7. [Insurance Pool](#insurance-pool)
8. [NFT-Stake Binding](#nft-stake-binding)
9. [Implementation Guide](#implementation-guide)
10. [FAQ](#faq)

## USD-Pegged Staking Model

### Core Principle

Validators stake based on USD value, not token count. Once staked, the NAMO token amount is fixed forever, regardless of price changes.

**Example**:
- Validator 11 needs $800K stake
- NAMO price at onboarding: $0.10
- Required tokens: 8,000,000 NAMO
- These 8M NAMO remain locked regardless of price fluctuations

### Benefits

1. **Predictable Entry Cost**: Validators know exact USD investment required
2. **No Speculation Incentive**: Can't profit from price appreciation during lock
3. **Fair Competition**: Later validators don't face exponentially higher token requirements
4. **Stability**: Reduces incentive for early dumps

## Tiered Validator Structure

### Investment Requirements

| Tier | Validators | Contract | Stake (USD) | Total | Performance Bond |
|------|------------|----------|-------------|-------|------------------|
| 1 | 1-10 | $100-280K | $200-380K | $300-660K | 20% |
| 2 | 11-20 | $400-580K | $800-980K | $1.2-1.56M | 25% |
| 3 | 21 | $650K | $1.5M | $2.15M | 30% |

### Detailed Breakdown

```
Validator 1:  $100k contract + $200k stake = $300k total
Validator 2:  $120k contract + $220k stake = $340k total
Validator 3:  $140k contract + $240k stake = $380k total
Validator 4:  $160k contract + $260k stake = $420k total
Validator 5:  $180k contract + $280k stake = $460k total
Validator 6:  $200k contract + $300k stake = $500k total
Validator 7:  $220k contract + $320k stake = $540k total
Validator 8:  $240k contract + $340k stake = $580k total
Validator 9:  $260k contract + $360k stake = $620k total
Validator 10: $280k contract + $380k stake = $660k total
--- TIER 2 ---
Validator 11: $400k contract + $800k stake = $1,200k total
Validator 12: $420k contract + $820k stake = $1,240k total
Validator 13: $440k contract + $840k stake = $1,280k total
Validator 14: $460k contract + $860k stake = $1,320k total
Validator 15: $480k contract + $880k stake = $1,360k total
Validator 16: $500k contract + $900k stake = $1,400k total
Validator 17: $520k contract + $920k stake = $1,440k total
Validator 18: $540k contract + $940k stake = $1,480k total
Validator 19: $560k contract + $960k stake = $1,520k total
Validator 20: $580k contract + $980k stake = $1,560k total
--- TIER 3 ---
Validator 21: $650k contract + $1,500k stake = $2,150k total

TOTAL: $8.25M contracts + $11.8M stakes = $20.05M
```

## Security Mechanisms

### 1. Multi-Layer Protection

```
┌─────────────────────────────────────────────────────────┐
│                  VALIDATOR SECURITY LAYERS               │
├─────────────────────────────────────────────────────────┤
│  Layer 1: STAKE LOCKING                                 │
│  ├─ 6/9/12-month minimum lock                          │
│  ├─ No unstaking during active validation              │
│  └─ 30-day notice for emergency exit                   │
│                                                         │
│  Layer 2: VESTING SCHEDULE                             │
│  ├─ 18/24/36-month gradual unlock                      │
│  ├─ Monthly release after initial lock                 │
│  └─ Performance-based acceleration                     │
│                                                         │
│  Layer 3: PERFORMANCE BOND                             │
│  ├─ 20/25/30% permanently locked                       │
│  ├─ Released after 3 years good standing              │
│  └─ Forfeited on malicious behavior                   │
│                                                         │
│  Layer 4: SLASHING CONDITIONS                          │
│  ├─ Technical violations: 0.01-5%                      │
│  ├─ Economic violations: 15-30%                        │
│  └─ Tiered multipliers by stake level                 │
│                                                         │
│  Layer 5: NFT-STAKE BINDING                            │
│  ├─ NFT transfer requires stake transfer              │
│  ├─ New owner inherits all obligations                │
│  └─ Cannot separate NFT from duties                   │
│                                                         │
│  Layer 6: INSURANCE POOL                               │
│  ├─ 2% of stakes fund pool (~$236K)                   │
│  ├─ Covers dumps up to tier limits                    │
│  └─ Community vote for claims                         │
└─────────────────────────────────────────────────────────┘
```

### 2. Stake Allocation

From each validator's stake:
- **78%**: Vestable amount (gradual unlock)
- **20-30%**: Performance bond (3-year lock)
- **2%**: Insurance pool contribution

## Vesting & Lock Periods

### Lock Period (No Access)

| Tier | Lock Period | Validators | Justification |
|------|-------------|------------|---------------|
| 1 | 6 months | 1-10 | Lower stakes, shorter commitment |
| 2 | 9 months | 11-20 | Higher stakes, medium commitment |
| 3 | 12 months | 21 | Highest stake, maximum commitment |

### Vesting Schedule (After Lock)

```
Tier 1 (18-month vesting after 6-month lock):
Month 7:  10% unlocked
Month 8:  15% unlocked
Month 9:  20% unlocked
Month 10-24: 3.67% monthly

Tier 2 (24-month vesting after 9-month lock):
Month 10: 10% unlocked
Month 11: 15% unlocked
Month 12-33: 3.41% monthly

Tier 3 (36-month vesting after 12-month lock):
Month 13: 10% unlocked
Month 14: 12% unlocked
Month 15-48: 2.29% monthly
```

### Performance Bond Release

- **Timeline**: 3 years from stake date
- **Conditions**: 
  - No major slashing events
  - Consistent uptime (>95%)
  - Active governance participation
  - No dump attempts

## Slashing Framework

### Violation Categories

#### Technical Violations

| Violation | Base Rate | Description |
|-----------|-----------|-------------|
| Downtime | 0.1%/day | After 24 hours offline |
| Double Sign | 5% | Signing conflicting blocks |
| Missed Blocks | 0.01% | Per 100 consecutive misses |

#### Economic Violations

| Violation | Base Rate | Description |
|-----------|-----------|-------------|
| Dump Attempt | 25% | Large coordinated sells |
| Market Manipulation | 15% | Price manipulation |
| Collusion | 30% | Coordinated attacks |
| Wash Trading | 20% | Fake volume generation |

### Tier Multipliers

- **Tier 1**: 1.0x base rate
- **Tier 2**: 1.5x base rate
- **Tier 3**: 2.0x base rate

*Higher stakes = higher responsibility = higher penalties*

### Slashing Distribution

- **50%**: To insurance pool
- **50%**: Burned permanently

## Circuit Breakers

### Price-Based Triggers

| Price Drop | Action | Duration | Additional Measures |
|------------|--------|----------|--------------------|
| 5% | Pause | 15 min | Reduced sell limits |
| 10% | Freeze | 1 hour | Large transfer freeze |
| 20% | Emergency | 24 hours | DAO vote required |

### Progressive Measures

```go
if consecutiveBreakers >= 3 {
    // Extend all unbonding by 7 days
    // Freeze transfers > $10,000
    // Reduce daily limits by 50%
    // Activate emergency insurance
}
```

## Daily Sell Limits

### Base Limits by Tier

| Tier | Daily | Weekly | Monthly |
|------|-------|--------|----------|
| 1 | 2% | 10% | 25% |
| 2 | 1% | 5% | 15% |
| 3 | 0.5% | 2.5% | 10% |

### Circuit Breaker Adjustments

- During circuit breaker: 50% reduction
- After 3 consecutive: 75% reduction
- Emergency mode: 90% reduction

## Insurance Pool

### Funding Structure

- **Initial**: 2% from all stakes = ~$236K
- **Ongoing**: 50% of slashing penalties
- **Target**: 5% of total stake value

### Coverage Limits

| Tier | Per Incident | Annual Cap | Deductible |
|------|--------------|------------|------------|
| 1 | $100K | $300K | 10% |
| 2 | $250K | $750K | 10% |
| 3 | $500K | $1.5M | 10% |

### Claim Process

1. **Detection**: Automated monitoring triggers
2. **Filing**: 7-day window to file claim
3. **Investigation**: Community review period
4. **Voting**: 75% approval required
5. **Payout**: Over 6 months if approved

## NFT-Stake Binding

### Binding Rules

1. **Inseparable**: NFT ownership = Validator rights + Stake obligations
2. **Transfer Requirements**:
   - 6-month minimum holding
   - 30-day notice period
   - Buyer accepts all obligations
   - 5% transfer fee to treasury

### Transfer Process

```go
// Pseudo-code for NFT transfer with stake
function transferNFTWithStake(from, to, price) {
    // Verify 6-month holding
    require(currentTime > nft.mintTime + 6 months)
    
    // Check 30-day notice
    require(transferNotice.active)
    
    // Transfer payment
    transfer(to, from, price * 0.95)  // 5% fee
    transfer(to, treasury, price * 0.05)
    
    // Pay 5% royalty to original validator
    transfer(to, originalValidator, price * 0.05)
    
    // Transfer NFT + all stake obligations
    nft.owner = to
    stake.owner = to
    // Vesting continues unchanged
    // Lock periods remain
    // Performance bond stays locked
}
```

## Implementation Guide

### For Validators

#### 1. Pre-Onboarding

```bash
# Check current NAMO price
deshchaind query oracle namo-price-usd

# Calculate required tokens
# Example: Validator 11 at $0.10/NAMO
# $800,000 / $0.10 = 8,000,000 NAMO

# Ensure sufficient balance
deshchaind query bank balances $(deshchaind keys show validator-key -a)
```

#### 2. Onboarding Process

```bash
# Stake tokens with rank
deshchaind tx validator onboard \
  --rank=11 \
  --stake-amount=8000000000000 \
  --from=validator-key \
  --gas=auto \
  --gas-adjustment=1.5

# Verify stake
deshchaind query validator stake-info $(deshchaind keys show validator-key -a)
```

#### 3. Monitoring

```bash
# Check vesting status
deshchaind query validator vesting-status $(deshchaind keys show validator-key -a)

# View available for withdrawal
deshchaind query validator unlocked-amount $(deshchaind keys show validator-key -a)

# Check slashing history
deshchaind query validator slashing-history $(deshchaind keys show validator-key -a)
```

#### 4. Withdrawals

```bash
# Withdraw unlocked tokens
deshchaind tx validator withdraw-unlocked \
  --amount=1000000000000 \
  --from=validator-key

# Check daily limit first
deshchaind query validator daily-limit $(deshchaind keys show validator-key -a)
```

### For NFT Traders

```bash
# Query NFT details
deshchaind query validator nft 1

# Check transfer eligibility
deshchaind query validator nft-transfer-status 1

# Initiate transfer (after 6 months)
deshchaind tx validator transfer-nft \
  --nft-id=1 \
  --to=deshchain1... \
  --price=50000000000000unamo \
  --from=validator-key
```

## Best Practices

### For Validators

1. **Capital Planning**: Ensure 20% buffer above minimum
2. **Price Monitoring**: Track NAMO price before onboarding
3. **Uptime Focus**: Maintain >99.9% to avoid slashing
4. **Governance Participation**: Vote to maintain reputation
5. **Exit Planning**: Plan 3+ years commitment

### For the Network

1. **Gradual Onboarding**: Don't onboard all 21 at once
2. **Price Stability**: Implement measures during onboarding
3. **Communication**: Clear timeline communication
4. **Support System**: Validator help channels
5. **Monitoring**: Real-time dashboards

## Risk Analysis

### Identified Risks

1. **Coordinated Dump** (High Impact, Low Probability)
   - Mitigation: Insurance pool, circuit breakers
   
2. **Price Manipulation** (Medium Impact, Medium Probability)
   - Mitigation: Oracle redundancy, time-weighted pricing
   
3. **Technical Failures** (Low Impact, Medium Probability)
   - Mitigation: Redundant infrastructure, quick recovery

### Stress Test Scenarios

| Scenario | Impact | Response |
|----------|--------|----------|
| 3 validators dump | -15% price | Circuit breaker + insurance |
| 50% price crash | Panic selling | Extended locks + DAO |
| Coordinated attack | Network risk | Emergency shutdown |

## FAQ

### General Questions

**Q: Why USD-pegged staking instead of token amount?**
A: Ensures fair entry regardless of token price, reduces speculation, and creates predictable costs.

**Q: Can I profit from NAMO price increase during staking?**
A: No, your tokens are locked at the entry amount. Profits come from validation rewards, not price appreciation.

**Q: What happens if NAMO price drops significantly?**
A: Your locked token amount remains the same. Circuit breakers and insurance protect the network.

### Technical Questions

**Q: How is the USD price determined at staking?**
A: Through decentralized oracles with multiple price feeds and time-weighted averages.

**Q: Can I partially withdraw vested tokens?**
A: Yes, after the lock period ends, you can withdraw any unlocked amount within daily limits.

**Q: What happens to slashed tokens?**
A: 50% goes to insurance pool, 50% is permanently burned.

### NFT Questions

**Q: Can I sell my validator NFT immediately?**
A: No, there's a 6-month lock period before any transfers.

**Q: Does NFT buyer get my stake back?**
A: The buyer inherits all obligations including lock periods and vesting schedules.

**Q: What's the minimum NFT trading price?**
A: 10,000 NAMO tokens to prevent manipulation.

## Conclusion

DeshChain's validator security system represents a paradigm shift in blockchain staking:

- **USD-pegged entry** ensures fairness
- **Multi-layer security** prevents dumps
- **Insurance protection** safeguards community
- **NFT binding** creates unique value
- **Progressive penalties** align incentives

This comprehensive approach protects the network while rewarding long-term commitment, creating a sustainable ecosystem for generations.

---

*For technical implementation details, see [Validator Module Documentation](./modules/VALIDATOR_MODULE.md)*