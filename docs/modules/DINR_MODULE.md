# DINR Module Documentation

## Overview

The DINR (Desh INR) module implements an algorithmic stablecoin pegged to the Indian Rupee (INR). It provides a decentralized, collateral-backed stable currency for the DeshChain ecosystem, enabling stable value transactions, remittances, and DeFi applications while maintaining the purchasing power of the Indian Rupee. All fees are paid in NAMO tokens with automatic swapping, featuring a tiered fee structure from 0.5% (< ₹10K) to 0.2% (> ₹10L) with a ₹830 maximum cap.

## Module Architecture

```mermaid
graph TB
    subgraph "DINR Module Components"
        KEEPER[Keeper]
        STABILITY[Stability Engine]
        COLLATERAL[Collateral Manager]
        LIQUIDATION[Liquidation Engine]
        YIELD[Yield Generator]
        ORACLE[Oracle Integration]
    end
    
    subgraph "External Dependencies"
        BANK[Bank Module]
        ORACLE_EXT[Oracle Module]
        NAMO[NAMO Module]
        TAX[Tax Module]
        INSURANCE[Insurance Fund]
    end
    
    KEEPER --> STABILITY
    KEEPER --> COLLATERAL
    KEEPER --> LIQUIDATION
    KEEPER --> YIELD
    
    STABILITY --> ORACLE
    COLLATERAL --> BANK
    LIQUIDATION --> INSURANCE
    ORACLE --> ORACLE_EXT
    
    KEEPER <--> BANK
    KEEPER <--> ORACLE_EXT
    KEEPER <--> TAX
```

## Stablecoin Mechanism

### Algorithmic Peg Maintenance

```mermaid
graph LR
    subgraph "Price Stability Loop"
        MONITOR[Price Monitor] --> DETECT{Deviation?}
        DETECT -->|Above Peg| INCREASE[Increase Supply]
        DETECT -->|Below Peg| DECREASE[Decrease Supply]
        DETECT -->|At Peg| MAINTAIN[Maintain]
        
        INCREASE --> INCENTIVE1[Lower Mint Fees]
        DECREASE --> INCENTIVE2[Higher Burn Rewards]
        
        INCENTIVE1 --> MONITOR
        INCENTIVE2 --> MONITOR
    end
```

### Collateralization Model

```mermaid
pie title "Collateral Tier Distribution"
    "Tier 1 Stables" : 40
    "Tier 2 Crypto" : 35
    "Tier 3 Alts" : 15
    "Insurance Reserve" : 10
```

## Collateral Tiers

### Tier 1: Stablecoins (140% Min Ratio)
| Asset | Max Allocation | Oracle Feed | Liquidation |
|-------|----------------|-------------|-------------|
| USDT | 25% | Chainlink | 130% |
| USDC | 25% | Band Protocol | 130% |
| DAI | 15% | DIA | 130% |
| BUSD | 10% | API3 | 130% |

### Tier 2: Major Cryptocurrencies (150% Min Ratio)
| Asset | Max Allocation | Oracle Feed | Liquidation |
|-------|----------------|-------------|-------------|
| BTC | 20% | Multiple | 135% |
| ETH | 20% | Multiple | 135% |
| BNB | 10% | Multiple | 135% |

### Tier 3: Alternative Assets (170% Min Ratio)
| Asset | Max Allocation | Oracle Feed | Liquidation |
|-------|----------------|-------------|-------------|
| MATIC | 5% | Multiple | 140% |
| SOL | 5% | Multiple | 140% |
| NAMO | 10% | Internal | 140% |

## Core Operations

### 1. Minting DINR

```mermaid
sequenceDiagram
    participant User
    participant DINR Module
    participant Oracle
    participant Bank
    participant Tax
    
    User->>DINR Module: Request Mint DINR
    DINR Module->>Oracle: Get Collateral Price
    Oracle-->>DINR Module: Price Data
    DINR Module->>DINR Module: Calculate Collateral Ratio
    
    alt Ratio >= 150%
        DINR Module->>Bank: Lock Collateral
        DINR Module->>Tax: Calculate Tiered Fees in NAMO
        DINR Module->>Bank: Mint DINR
        DINR Module->>User: Transfer DINR
    else Ratio < 150%
        DINR Module->>User: Reject (Insufficient Collateral)
    end
```

**Process Steps:**
1. User deposits collateral (multi-asset supported)
2. Oracle provides real-time price feeds
3. System calculates collateral value in INR
4. Minimum 150% collateralization enforced
5. Tiered fee applied: 0.5% (< ₹10K) → 0.3% (₹10K-1L) → 0.2% (> ₹1L), capped at ₹830 in NAMO
6. DINR minted and transferred to user

### 2. Burning DINR

```mermaid
sequenceDiagram
    participant User
    participant DINR Module
    participant Oracle
    participant Bank
    
    User->>DINR Module: Request Burn DINR
    DINR Module->>Oracle: Get Current Prices
    Oracle-->>DINR Module: Price Data
    DINR Module->>Bank: Burn DINR from User
    DINR Module->>DINR Module: Calculate Collateral Return
    DINR Module->>Bank: Release Collateral
    DINR Module->>User: Transfer Collateral
```

**Process Steps:**
1. User submits DINR for burning
2. System calculates proportional collateral return
3. Tiered fee applied: 0.5% (< ₹10K) → 0.3% (₹10K-1L) → 0.2% (> ₹1L), capped at ₹830 in NAMO
4. DINR burned from circulation
5. Collateral released to user

### 3. Liquidation Process

```mermaid
graph TB
    subgraph "Liquidation Engine"
        MONITOR[Position Monitor] --> CHECK{Health Check}
        CHECK -->|Ratio < 130%| TRIGGER[Trigger Liquidation]
        CHECK -->|Ratio >= 130%| SAFE[Position Safe]
        
        TRIGGER --> AUCTION[Dutch Auction]
        AUCTION --> LIQUIDATOR[Liquidator Buys]
        LIQUIDATOR --> PENALTY[10% Penalty]
        PENALTY --> INSURANCE[5% to Insurance]
        PENALTY --> PLATFORM[5% to Platform]
    end
```

## Module Parameters

```go
type Params struct {
    // Fee Parameters (all fees paid in NAMO)
    TieredFeeStructure   bool    // true - enables tiered fees
    Tier1Fee             uint64  // 50 (0.5% for < ₹10K)
    Tier2Fee             uint64  // 30 (0.3% for ₹10K-1L)
    Tier3Fee             uint64  // 20 (0.2% for > ₹1L)
    MaxFeeNAMO           string  // "830000000" (₹830 in micro units)
    
    // Collateral Parameters
    MinCollateralRatio   uint64  // 15000 (150%)
    LiquidationThreshold uint64  // 13000 (130%)
    LiquidationPenalty   uint64  // 1000 (10%)
    
    // Stability Parameters
    MaxPriceDeviation    uint64  // 100 (1%)
    RebalanceInterval    int64   // 3600 (1 hour)
    InsuranceFundTarget  uint64  // 500 (5%)
    
    // Yield Parameters
    YieldDeploymentRatio uint64  // 8000 (80%)
    MinYieldThreshold    string  // "1000000" (₹10 lakh)
}
```

## Stability Mechanisms

### 1. Price Oracle System

```mermaid
graph LR
    subgraph "Multi-Oracle Aggregation"
        O1[Chainlink] --> AGG[Aggregator]
        O2[Band Protocol] --> AGG
        O3[DIA] --> AGG
        O4[API3] --> AGG
        
        AGG --> MEDIAN[Median Filter]
        MEDIAN --> OUTLIER[Outlier Detection]
        OUTLIER --> FINAL[Final Price]
    end
```

### 2. Stability Fee Adjustment

```mermaid
graph TD
    PRICE[DINR Price] --> COMPARE{Compare to ₹1}
    
    COMPARE -->|> ₹1.01| HIGH[Above Peg]
    COMPARE -->|< ₹0.99| LOW[Below Peg]
    COMPARE -->|₹0.99-1.01| STABLE[Stable]
    
    HIGH --> ADJ1[Reduce Mint Fee<br/>Increase Burn Reward]
    LOW --> ADJ2[Increase Mint Fee<br/>Reduce Burn Reward]
    STABLE --> MAINTAIN[No Adjustment]
```

### 3. Insurance Fund

- Target: 5% of total DINR supply
- Funded by:
  - 50% of liquidation penalties
  - Excess yield from strategies
  - Emergency minting (governance approved)
- Used for:
  - Black swan event coverage
  - Bad debt absorption
  - Peg defense operations

## Yield Generation Strategy

```mermaid
graph TB
    subgraph "Yield Deployment"
        COLLATERAL[Idle Collateral] --> ASSESS{Risk Assessment}
        
        ASSESS -->|Low Risk| STABLE[Stable Strategies<br/>3-5% APY]
        ASSESS -->|Medium Risk| DEFI[DeFi Protocols<br/>5-10% APY]
        ASSESS -->|Reserve| IDLE[Keep Idle<br/>0% APY]
        
        STABLE --> YIELD[Yield Generated]
        DEFI --> YIELD
        
        YIELD --> DIST{Distribution}
        DIST -->|40%| INSURANCE[Insurance Fund]
        DIST -->|30%| PLATFORM[Platform Revenue]
        DIST -->|30%| HOLDERS[DINR Holders]
    end
```

## Transaction Types

### 1. MsgMintDINR
Mints new DINR tokens against collateral.

```go
type MsgMintDINR struct {
    Minter      string
    Collateral  sdk.Coin
    DinrToMint  sdk.Coin
}
```

### 2. MsgBurnDINR
Burns DINR tokens and returns collateral.

```go
type MsgBurnDINR struct {
    Burner           string
    DinrToBurn       sdk.Coin
    CollateralDenom  string
}
```

### 3. MsgAddCollateral
Adds additional collateral to improve position health.

```go
type MsgAddCollateral struct {
    Depositor       string
    AdditionalColl  sdk.Coin
}
```

### 4. MsgLiquidatePosition
Liquidates an undercollateralized position.

```go
type MsgLiquidatePosition struct {
    Liquidator      string
    TargetAddress   string
    MaxDinrToPay    sdk.Coin
}
```

## Query Endpoints

### 1. QueryParams
Returns current module parameters.

**Request**: `/deshchain/dinr/v1/params`

**Response**:
```json
{
  "params": {
    "tiered_fee_structure": true,
    "tier1_fee": "50",
    "tier2_fee": "30",
    "tier3_fee": "20",
    "max_fee_namo": "830000000",
    "min_collateral_ratio": "15000",
    "liquidation_threshold": "13000"
  }
}
```

### 2. QueryUserPosition
Returns user's DINR position details.

**Request**: `/deshchain/dinr/v1/position/{address}`

**Response**:
```json
{
  "position": {
    "address": "deshchain1...",
    "collateral": [
      {"denom": "usdt", "amount": "100000000"}
    ],
    "dinr_minted": {"denom": "dinr", "amount": "7000000"},
    "collateral_ratio": "15500",
    "health_factor": "1.19",
    "is_liquidatable": false
  }
}
```

### 3. QueryStabilityData
Returns current stability metrics.

**Request**: `/deshchain/dinr/v1/stability`

**Response**:
```json
{
  "stability_data": {
    "current_price": "1.0023",
    "target_price": "1.0000",
    "price_deviation": "23",
    "total_supply": {"denom": "dinr", "amount": "1000000000000"},
    "total_collateral_value": {"denom": "inr", "amount": "1550000000000"},
    "global_collateral_ratio": "15500"
  }
}
```

### 4. QueryCollateralAssets
Returns supported collateral assets.

**Request**: `/deshchain/dinr/v1/collateral-assets`

**Response**:
```json
{
  "assets": [
    {
      "denom": "usdt",
      "tier": "tier1_stable",
      "min_collateral_ratio": "14000",
      "max_allocation": "2500",
      "is_active": true
    }
  ]
}
```

## Events

### 1. DINR Minted Event
```json
{
  "type": "dinr_minted",
  "attributes": [
    {"key": "minter", "value": "{address}"},
    {"key": "collateral", "value": "{amount}"},
    {"key": "dinr_minted", "value": "{amount}"},
    {"key": "collateral_ratio", "value": "{ratio}"},
    {"key": "fee_paid", "value": "{amount}"}
  ]
}
```

### 2. DINR Burned Event
```json
{
  "type": "dinr_burned",
  "attributes": [
    {"key": "burner", "value": "{address}"},
    {"key": "dinr_burned", "value": "{amount}"},
    {"key": "collateral_returned", "value": "{amount}"},
    {"key": "fee_paid", "value": "{amount}"}
  ]
}
```

### 3. Position Liquidated Event
```json
{
  "type": "position_liquidated",
  "attributes": [
    {"key": "liquidator", "value": "{address}"},
    {"key": "liquidated_user", "value": "{address}"},
    {"key": "dinr_repaid", "value": "{amount}"},
    {"key": "collateral_seized", "value": "{amount}"},
    {"key": "penalty", "value": "{amount}"}
  ]
}
```

## Integration with Other Modules

### 1. Oracle Module Integration
- Real-time price feeds for all collateral assets
- INR/USD exchange rate for stability calculations
- Multi-oracle aggregation for reliability

### 2. Tax Module Integration
- Transaction fees subject to platform tax
- Tax distribution to NGOs and validators

### 3. NAMO Module Integration
- NAMO can be used as Tier 3 collateral
- All platform fees exclusively paid in NAMO tokens
- Automatic token swapping for fee collection
- 2% of all fees burned for deflationary pressure

### 4. Money Order Module Integration
- DINR used as primary currency for money orders
- Instant settlement without volatility risk

### 5. Remittance Module Integration
- DINR enables stable cross-border transfers
- No FX risk during transaction processing

## Security Considerations

1. **Oracle Manipulation Protection**
   - Multi-oracle aggregation with median pricing
   - Deviation limits and circuit breakers
   - Time-weighted average pricing (TWAP)

2. **Flash Loan Protection**
   - Minimum lock period for collateral
   - Rate limiting on large operations
   - Progressive fee structure

3. **Liquidation Safety**
   - Grace period notifications
   - Partial liquidation support
   - Liquidator incentive caps

4. **Emergency Mechanisms**
   - Pause functionality (governance only)
   - Emergency collateral ratio adjustment
   - Insurance fund deployment

## Risk Management

### Collateral Risk Matrix

```mermaid
graph TB
    subgraph "Risk Assessment"
        LOW[Low Risk<br/>Stablecoins] --> RATIO1[140% Min Ratio]
        MED[Medium Risk<br/>BTC/ETH] --> RATIO2[150% Min Ratio]
        HIGH[High Risk<br/>Alts] --> RATIO3[170% Min Ratio]
        
        RATIO1 --> MONITOR[Continuous Monitoring]
        RATIO2 --> MONITOR
        RATIO3 --> MONITOR
        
        MONITOR --> ACTION{Action Required?}
        ACTION -->|Yes| LIQUIDATE[Liquidate]
        ACTION -->|No| CONTINUE[Continue]
    end
```

## Best Practices

1. **For Users**
   - Maintain healthy collateral ratios (>150%)
   - Diversify collateral across tiers
   - Monitor position health regularly
   - Use DINR for stable transactions

2. **For Developers**
   - Always check collateral ratios before operations
   - Handle liquidation events gracefully
   - Implement proper error handling for oracle failures
   - Use batch operations for efficiency

3. **For Liquidators**
   - Monitor positions approaching liquidation threshold
   - Prepare sufficient DINR for liquidations
   - Understand gas optimization strategies
   - Calculate profitability including penalties

## Revenue Model

### Fee Structure (All Fees in NAMO)
- **Minting/Burning Fees**:
  - < ₹10K: 0.5%
  - ₹10K - ₹1L: 0.3%
  - > ₹1L: 0.2%
  - Maximum cap: ₹830 in NAMO tokens
- **Liquidation Penalty**: 10% of liquidated amount
- **Yield Generation**: Performance-based 0-8% APY on deployed collateral

### Revenue Distribution
```mermaid
pie title "DINR Revenue Distribution"
    "NGO Charity (28%)" : 28
    "Validators (25%)" : 25
    "Community Rewards (18%)" : 18
    "Development (14%)" : 14
    "Operations (8%)" : 8
    "Founder Royalty (5%)" : 5
    "NAMO Burn (2%)" : 2
```

### Projected Revenue
- Year 1: ₹150 Crore
- Year 2: ₹842 Crore  
- Year 3: ₹2,636 Crore
- Year 5: ₹11,866 Crore

## CLI Commands

### Query Commands
```bash
# Query module parameters
deshchaind query dinr params

# Query user position
deshchaind query dinr position [address]

# Query stability data
deshchaind query dinr stability

# Query collateral assets
deshchaind query dinr collateral-assets

# Query total supply
deshchaind query dinr total-supply
```

### Transaction Commands
```bash
# Mint DINR
deshchaind tx dinr mint [collateral-amount] [dinr-amount] --from [key]

# Burn DINR
deshchaind tx dinr burn [dinr-amount] [collateral-denom] --from [key]

# Add collateral
deshchaind tx dinr add-collateral [amount] --from [key]

# Liquidate position
deshchaind tx dinr liquidate [target-address] [max-dinr] --from [key]
```

## FAQ

**Q: How is DINR different from other stablecoins?**
A: DINR is specifically pegged to INR, uses multi-tier collateral system with tiered fees (0.5% → 0.2%), all fees paid in NAMO with ₹830 cap, generates yield on idle collateral, and dedicates 28% of revenue to charity with 2% burned.

**Q: What happens if DINR loses its peg?**
A: The stability mechanism automatically adjusts fees and incentives. If deviation persists, the insurance fund can be deployed to defend the peg.

**Q: Can I use NAMO tokens as collateral?**
A: Yes, NAMO is accepted as Tier 3 collateral with a 170% minimum collateralization ratio. Additionally, all fees are paid exclusively in NAMO tokens with automatic swapping from your preferred token.

**Q: How are liquidations handled?**
A: Positions below 130% ratio can be liquidated through a Dutch auction mechanism, with a 10% penalty split between insurance fund and platform.

**Q: Is there a minimum amount to mint DINR?**
A: Yes, the minimum minting amount is ₹100 worth of DINR to ensure economic viability. Note that transactions under ₹100 in the broader ecosystem enjoy FREE fees under the progressive tax structure.

---

For more information, see the [Module Overview](../MODULE_OVERVIEW.md) or explore other [DeshChain Modules](../MODULE_OVERVIEW.md#module-categories).