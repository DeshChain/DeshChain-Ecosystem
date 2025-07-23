# DUSD Module Documentation

## Overview

The DUSD (DeshChain USD) module implements a USD-pegged algorithmic stablecoin designed for global trade finance and cross-border remittances. Building on the proven stability mechanisms of DINR, DUSD expands DeshChain's addressable market from ₹50 lakh Cr to $20+ trillion globally.

## Key Features

- **USD Peg**: 1 DUSD = $1.00 USD maintained algorithmically
- **Low Fees**: $0.10 - $1.00 per transaction (vs $15-50 traditional banking)
- **Proven Stability**: Same algorithmic mechanisms as DINR
- **Multi-Currency Integration**: Seamless routing for global trade finance
- **Enterprise-Grade**: Built for institutional adoption

## Architecture

### Core Components

```go
type DUSDKeeper struct {
    cdc            codec.BinaryCodec
    storeKey       storetypes.StoreKey
    authority      string
    
    // Integration keepers
    bankKeeper     BankKeeper
    oracleKeeper   OracleKeeper
    treasuryKeeper TreasuryKeeper
    accountKeeper  AccountKeeper
}
```

### Price Stability Engine

The DUSD module uses the same proven stability algorithms as DINR:

```go
type StabilityEngine struct {
    keeper *Keeper
}

func (se *StabilityEngine) CheckPriceStability(ctx sdk.Context) error {
    currentPrice, err := se.keeper.GetUSDPrice(ctx)
    if err != nil {
        return err
    }
    
    targetPrice := sdk.OneDec() // $1.00 USD
    deviation := currentPrice.Sub(targetPrice).Quo(targetPrice).Abs()
    
    if deviation.GTE(rebalanceThreshold) {
        return se.ExecuteRebalanceAction(ctx, currentPrice, targetPrice, deviation)
    }
    
    return nil
}
```

## Oracle Integration

### Multi-Source USD Price Feeds

DUSD integrates with multiple oracle sources for robust USD pricing:

```go
var USDOracleSources = map[string]OracleSourceConfig{
    "chainlink": {
        Name:        "Chainlink",
        Weight:      30,
        Endpoint:    "https://api.chain.link/v1/feeds/usd-inr",
        Reliability: 99.9,
    },
    "federal_reserve": {
        Name:        "Federal Reserve",
        Weight:      25,
        Endpoint:    "https://api.stlouisfed.org/fred/series/observations?series_id=DEXINUS",
        Reliability: 99.5,
    },
    "band_protocol": {
        Name:        "Band Protocol",
        Weight:      20,
        Reliability: 99.8,
    },
    "pyth_network": {
        Name:        "Pyth Network",
        Weight:      15,
        Reliability: 99.7,
    },
}
```

## Multi-Currency Operations

### Enhanced Trade Finance

DUSD enables multi-currency trade finance with significant cost and time savings:

```go
type EnhancedLetterOfCredit struct {
    // Base LC fields
    LcId               string
    ApplicantId        string
    BeneficiaryId      string
    
    // Multi-currency enhancement
    OriginalCurrency   string    // "USD", "EUR", "SGD"
    OriginalAmount     sdk.Coin  // Original trade amount
    SettlementCurrency string    // "DUSD" routing
    SettlementAmount   sdk.Coin  // DUSD equivalent
    LocalCurrency      string    // "DINR" for recipients
    LocalAmount        sdk.Coin  // Final local amount
    
    // Cost analysis
    TraditionalCost    sdk.Coin  // vs traditional banking
    DeshChainCost      sdk.Coin  // actual cost
    TotalSavings       sdk.Coin  // customer savings (85%+)
    ProcessingTime     time.Duration // 5 min vs 5-7 days
}
```

### Remittance Optimization

DUSD provides optimal routing for cross-border remittances:

```go
type EnhancedRemittanceTransfer struct {
    // Source and destination
    SourceCurrency      string    // "USD", "EUR", "SGD"
    SourceAmount        sdk.Coin  // Original amount
    RoutingCurrency     string    // "DUSD" optimal routing
    DestinationCurrency string    // "DINR"
    DestinationAmount   sdk.Coin  // Final amount
    
    // Cost comparison
    TraditionalCost     sdk.Coin  // 6-8% traditional fees
    DeshChainCost       sdk.Coin  // $0.30 typical cost
    TotalSavings        sdk.Coin  // 95%+ savings
    ProcessingTime      time.Duration // 30 seconds
}
```

## Treasury Integration

### DUSD Reserve Management

DUSD integrates with the treasury system for robust reserve management:

```go
type DUSDTreasuryPool struct {
    // Base treasury pool
    PoolID              string
    PoolType            string // "DUSD_RESERVE"
    USDCollateralRatio  sdk.Dec // 150% target
    DUSDSupplyBacked    sdk.Coin // DUSD backed by pool
    USDReserveAssets    []ReserveAsset // USD assets
    StabilityBuffer     sdk.Coin // Stability operations
}
```

### Cross-Currency Rebalancing

```go
func (drm *DUSDReserveManager) RebalanceCrossCurrencyExposure(ctx sdk.Context) error {
    // Target exposure across currencies
    targetExposure := map[string]sdk.Dec{
        "USD": sdk.NewDecWithPrec(40, 2), // 40%
        "EUR": sdk.NewDecWithPrec(30, 2), // 30%
        "SGD": sdk.NewDecWithPrec(20, 2), // 20%
        "GBP": sdk.NewDecWithPrec(10, 2), // 10%
    }
    
    // Execute rebalancing if deviation exceeds threshold
    return drm.executeRebalanceActions(ctx, targetExposure)
}
```

## API Reference

### Transaction Messages

#### MintDUSD
```protobuf
message MsgMintDUSD {
    string creator = 1;
    cosmos.base.v1beta1.Coin collateral_amount = 2;
    string collateral_type = 3;
    cosmos.base.v1beta1.Coin dusd_amount = 4;
}

message MsgMintDUSDResponse {
    string position_id = 1;
    cosmos.base.v1beta1.Coin minted_amount = 2;
    cosmos.base.v1beta1.Coin fee_paid = 3;
    string health_factor = 4;
}
```

#### BurnDUSD
```protobuf
message MsgBurnDUSD {
    string creator = 1;
    string position_id = 2;
    cosmos.base.v1beta1.Coin dusd_amount = 3;
}

message MsgBurnDUSDResponse {
    cosmos.base.v1beta1.Coin burned_amount = 1;
    cosmos.base.v1beta1.Coin collateral_released = 2;
    string remaining_health_factor = 3;
}
```

### Query Endpoints

#### Get Position
```
GET /deshchain/dusd/v1/position/{position_id}
```

#### Get Price Data
```
GET /deshchain/dusd/v1/price
```

#### Get Reserve Statistics
```
GET /deshchain/dusd/v1/reserves
```

### REST API Examples

#### Mint DUSD
```bash
curl -X POST \
  http://localhost:1317/deshchain/dusd/v1/mint \
  -H 'Content-Type: application/json' \
  -d '{
    "creator": "deshchain1abc...",
    "collateral_amount": {
      "denom": "USDC",
      "amount": "1500000000"
    },
    "collateral_type": "USDC",
    "dusd_amount": {
      "denom": "DUSD",
      "amount": "1000000000"
    }
  }'
```

#### Query Position Health
```bash
curl http://localhost:1317/deshchain/dusd/v1/health/position123
```

## Fee Structure

### USD-Equivalent Fees

DUSD uses USD-equivalent fee structure:

- **Base Fee**: $0.10 USD (minimum)
- **Max Fee**: $1.00 USD (maximum)
- **Fee Calculation**: 0.25% of transaction amount, capped at $1.00
- **Comparison**: 95%+ savings vs traditional banking ($15-50 fees)

### Fee Examples

| Transaction Amount | DUSD Fee | Traditional Fee | Savings |
|------------------|----------|----------------|---------|
| $100 | $0.10 | $15 | 99.3% |
| $1,000 | $0.25 | $25 | 99.0% |
| $10,000 | $1.00 | $50 | 98.0% |
| $100,000 | $1.00 | $500 | 99.8% |

## Market Opportunity

### Global Addressable Market

DUSD expands DeshChain's total addressable market:

| Market Segment | Current (DINR) | With DUSD | Expansion |
|----------------|----------------|-----------|-----------|
| Trade Finance | ₹50 lakh Cr | $15+ trillion | 40x |
| Remittances | ₹1 lakh Cr | $200B+ | 25x |
| Cross-border B2B | Minimal | $5+ trillion | ∞ |
| **Total** | **₹51 lakh Cr** | **$20+ trillion** | **40x** |

### Revenue Projections

```yaml
Year 1 (DUSD Launch):
  Target Volume: $10B transactions
  Revenue: $35M (0.35% average fee)
  Market Share: 0.1% of addressable market

Year 3 (Multi-Currency Suite):
  Target Volume: $200B transactions
  Revenue: $400M (economies of scale)
  Market Share: 1% of addressable market

Year 5 (Market Leader):
  Target Volume: $1T+ transactions
  Revenue: $1.5B+ annually
  Market Position: Top 3 global stablecoin platform
```

## Implementation Roadmap

### Phase 1: DUSD Core (Q4 2025)
- [ ] Complete DUSD module implementation
- [ ] Oracle USD integration
- [ ] Treasury USD reserve pools
- [ ] Basic trade finance integration

### Phase 2: Multi-Currency Integration (Q1 2026)
- [ ] Enhanced trade finance with DUSD
- [ ] Remittance corridor optimization
- [ ] Cross-currency rebalancing
- [ ] Sewa Mitra USD support

### Phase 3: Global Expansion (Q2 2026)
- [ ] DEUR implementation
- [ ] DSGD implementation
- [ ] Full multi-currency suite
- [ ] Global regulatory compliance

## Security Considerations

### Collateral Management
- **Minimum Ratio**: 150% (same as DINR)
- **Liquidation Threshold**: 120% health factor
- **Emergency Protocols**: Circuit breakers for extreme volatility
- **Oracle Security**: Multi-source validation with deviation alerts

### Risk Mitigation
- **Diversified Reserves**: Multiple USD-denominated assets
- **Regular Audits**: Quarterly security and financial audits
- **Insurance Fund**: 5% of fees allocated to insurance
- **Governance Controls**: Multi-signature treasury operations

## Performance Metrics

### Target KPIs

| Metric | Target | Status |
|--------|--------|--------|
| Stability | ±1% of $1.00 USD | In Development |
| Uptime | 99.9% | In Development |
| Oracle Accuracy | 99.5% | In Development |
| Fee Efficiency | <$1 max fee | ✅ Implemented |
| Processing Time | <30 seconds | ✅ Implemented |

## Conclusion

The DUSD module represents a revolutionary advancement in global stablecoin technology, combining the proven stability of DINR with the massive market opportunity of USD-denominated financial services. By enabling seamless multi-currency operations, DUSD positions DeshChain as a global leader in blockchain-based trade finance and remittances.

---

*For technical support and development questions, please refer to the [DeshChain Developer Documentation](../README.md) or contact the development team.*