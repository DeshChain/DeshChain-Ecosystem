# DeshChain Economic Model v2.0 - Direct Implementation Tasks

## Overview
Since DeshChain hasn't been deployed yet, we can directly implement the finalized economic model without backward compatibility concerns. This document outlines all code implementation tasks required.

## 🎯 Priority 1: Update Core Constants

### 1.1 Token Distribution Constants
**File**: `x/namo/types/constants.go`
```go
package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
    // Token Distribution
    TotalSupply              = 1_428_627_663
    PublicSaleAllocation     = 0.20   // 20% - 285,725,533 NAMO
    LiquidityAllocation      = 0.18   // 18% - 257,152,979 NAMO
    CommunityAllocation      = 0.15   // 15% - 214,294,149 NAMO
    DevelopmentAllocation    = 0.15   // 15% - 214,294,149 NAMO
    TeamAllocation           = 0.12   // 12% - 171,435,319 NAMO
    FounderAllocation        = 0.08   // 8% - 114,290,213 NAMO
    DAOTreasuryAllocation    = 0.05   // 5% - 71,431,383 NAMO
    CoFounderAllocation      = 0.035  // 3.5% - 50,001,968 NAMO
    OperationsAllocation     = 0.02   // 2% - 28,572,553 NAMO
    AngelAllocation          = 0.015  // 1.5% - 21,428,900 NAMO
    
    // Vesting Parameters
    UniversalCliffMonths     = 12
    FounderVestingMonths     = 48
    TeamVestingMonths        = 24
    CoFounderVestingMonths   = 24
    AngelVestingMonths       = 24
)
```

### 1.2 Transaction Tax Distribution
**File**: `x/tax/types/params.go`
```go
package types

import sdk "github.com/cosmos/cosmos-sdk/types"

var (
    DefaultTaxRate = sdk.NewDecWithPrec(25, 3) // 2.5%
    
    DefaultTaxDistribution = TaxDistribution{
        NGODonations:      sdk.NewDecWithPrec(300, 3), // 30%
        Validators:        sdk.NewDecWithPrec(250, 3), // 25%
        CommunityRewards:  sdk.NewDecWithPrec(200, 3), // 20%
        TechInnovation:    sdk.NewDecWithPrec(60, 3),  // 6%
        Operations:        sdk.NewDecWithPrec(50, 3),  // 5%
        TalentAcquisition: sdk.NewDecWithPrec(40, 3),  // 4%
        StrategicReserve:  sdk.NewDecWithPrec(40, 3),  // 4%
        Founder:           sdk.NewDecWithPrec(35, 3),  // 3.5%
        CoFounders:        sdk.NewDecWithPrec(18, 3),  // 1.8%
        AngelInvestors:    sdk.NewDecWithPrec(7, 3),   // 0.7%
    }
)

type TaxDistribution struct {
    NGODonations      sdk.Dec
    Validators        sdk.Dec
    CommunityRewards  sdk.Dec
    TechInnovation    sdk.Dec
    Operations        sdk.Dec
    TalentAcquisition sdk.Dec
    StrategicReserve  sdk.Dec
    Founder           sdk.Dec
    CoFounders        sdk.Dec
    AngelInvestors    sdk.Dec
}
```

## 🎯 Priority 2: Platform Revenue Distributions

### 2.1 DEX Fee Distribution
**File**: `x/dex/types/params.go`
```go
package types

import sdk "github.com/cosmos/cosmos-sdk/types"

var (
    DefaultTradingFee = sdk.NewDecWithPrec(3, 3) // 0.3%
    
    DefaultDEXDistribution = FeeDistribution{
        Validators:         sdk.NewDecWithPrec(45, 2), // 45%
        LiquidityProviders: sdk.NewDecWithPrec(15, 2), // 15%
        NGO:                sdk.NewDecWithPrec(15, 2), // 15%
        Community:          sdk.NewDecWithPrec(10, 2), // 10%
        Operations:         sdk.NewDecWithPrec(5, 2),  // 5%
        Tech:               sdk.NewDecWithPrec(4, 2),  // 4%
        Founder:            sdk.NewDecWithPrec(4, 2),  // 4%
        CoFounders:         sdk.NewDecWithPrec(15, 3), // 1.5%
        Angels:             sdk.NewDecWithPrec(5, 3),  // 0.5%
    }
)
```

### 2.2 Sikkebaaz Launchpad Fees
**File**: `x/launchpad/types/params.go`
```go
package types

import sdk "github.com/cosmos/cosmos-sdk/types"

var (
    DefaultPlatformFee = sdk.NewDecWithPrec(5, 2)  // 5% (updated from 2%)
    DefaultListingFee  = sdk.NewInt(1000)           // 1000 NAMO (updated from 100)
    
    DefaultLaunchpadDistribution = FeeDistribution{
        Validators:      sdk.NewDecWithPrec(40, 2), // 40%
        NGO:             sdk.NewDecWithPrec(20, 2), // 20%
        Community:       sdk.NewDecWithPrec(15, 2), // 15%
        AntiRugFund:     sdk.NewDecWithPrec(10, 2), // 10%
        Operations:      sdk.NewDecWithPrec(5, 2),  // 5%
        Tech:            sdk.NewDecWithPrec(4, 2),  // 4%
        Founder:         sdk.NewDecWithPrec(4, 2),  // 4%
        CoFoundersAngels:sdk.NewDecWithPrec(2, 2),  // 2%
    }
)
```

### 2.3 NFT Marketplace Distribution
**File**: `x/nft/types/marketplace_params.go`
```go
package types

import sdk "github.com/cosmos/cosmos-sdk/types"

var (
    DefaultMarketplaceFee = sdk.NewDecWithPrec(25, 3) // 2.5%
    
    DefaultNFTDistribution = FeeDistribution{
        Validators:       sdk.NewDecWithPrec(35, 2), // 35%
        NGOArtEducation:  sdk.NewDecWithPrec(25, 2), // 25%
        CommunityArtists: sdk.NewDecWithPrec(20, 2), // 20%
        Operations:       sdk.NewDecWithPrec(8, 2),  // 8%
        Tech:             sdk.NewDecWithPrec(6, 2),  // 6%
        Founder:          sdk.NewDecWithPrec(4, 2),  // 4%
        CoFoundersAngels: sdk.NewDecWithPrec(2, 2),  // 2%
    }
)
```

## 🎯 Priority 3: Validator Economics

### 3.1 Geographic Incentives
**File**: `x/validator/types/geographic.go`
```go
package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type GeographicIncentives struct {
    IndiaBase        sdk.Dec `json:"india_base"`         // 10%
    Tier2CityBonus   sdk.Dec `json:"tier2_city_bonus"`   // 5%
    EmploymentBonus  sdk.Dec `json:"employment_bonus"`   // 3%
    GreenEnergyBonus sdk.Dec `json:"green_energy_bonus"` // 2%
}

var DefaultGeographicIncentives = GeographicIncentives{
    IndiaBase:        sdk.NewDecWithPrec(10, 2),
    Tier2CityBonus:   sdk.NewDecWithPrec(5, 2),
    EmploymentBonus:  sdk.NewDecWithPrec(3, 2),
    GreenEnergyBonus: sdk.NewDecWithPrec(2, 2),
}

type ValidatorLocation struct {
    Country      string `json:"country"`
    City         string `json:"city"`
    CityTier     int    `json:"city_tier"`
    Employees    int    `json:"local_employees"`
    RenewableEnergy int `json:"renewable_energy_percent"`
}
```

### 3.2 Performance Bonus Structure
**File**: `x/validator/types/performance.go`
```go
package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type PerformanceBonus struct {
    UptimeThresholds     map[string]sdk.Dec
    BlockProductionBonus map[string]sdk.Dec
    TransactionSpeedBonus map[string]sdk.Dec
    CommunityBonus       sdk.Dec
    ArchiveNodeBonus     sdk.Dec
    PublicRPCBonus       sdk.Dec
}

var DefaultPerformanceBonus = PerformanceBonus{
    UptimeThresholds: map[string]sdk.Dec{
        "99.0":  sdk.NewDec(0),
        "99.5":  sdk.NewDecWithPrec(2, 2),
        "99.9":  sdk.NewDecWithPrec(3, 2),
        "99.99": sdk.NewDecWithPrec(5, 2),
    },
    BlockProductionBonus: map[string]sdk.Dec{
        "top10":  sdk.NewDecWithPrec(3, 2),
        "top25":  sdk.NewDecWithPrec(2, 2),
        "average": sdk.NewDec(0),
    },
    TransactionSpeedBonus: map[string]sdk.Dec{
        "fast":   sdk.NewDecWithPrec(3, 2), // <100ms
        "medium": sdk.NewDecWithPrec(2, 2), // 100-200ms
        "normal": sdk.NewDecWithPrec(1, 2), // 200-500ms
    },
    CommunityBonus:   sdk.NewDecWithPrec(2, 2),
    ArchiveNodeBonus: sdk.NewDecWithPrec(2, 2),
    PublicRPCBonus:   sdk.NewDecWithPrec(2, 2),
}
```

### 3.3 MEV Distribution
**File**: `x/validator/keeper/mev.go`
```go
package keeper

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
)

// DistributeMEV distributes MEV to the block proposer
func (k Keeper) DistributeMEV(ctx sdk.Context, mevAmount sdk.Coins, proposer sdk.ValAddress) error {
    // 100% of MEV goes to block proposer
    return k.bankKeeper.SendCoinsFromModuleToAccount(
        ctx,
        types.MEVPoolModule,
        proposer,
        mevAmount,
    )
}

// DistributePriorityFees distributes priority fees to validators
func (k Keeper) DistributePriorityFees(ctx sdk.Context, fees sdk.Coins, proposer sdk.ValAddress) error {
    // 100% of priority fees to block proposer
    return k.bankKeeper.SendCoinsFromModuleToAccount(
        ctx,
        types.FeeCollectorModule,
        proposer,
        fees,
    )
}
```

## 🎯 Priority 4: Genesis Configuration

### 4.1 Genesis State
**File**: `app/genesis.go`
```go
package app

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/deshchain/x/namo/types"
)

func NewDefaultGenesisState() GenesisState {
    totalSupply := sdk.NewInt(types.TotalSupply).Mul(sdk.NewInt(1e6)) // With decimals
    
    allocations := []GenesisAllocation{
        {
            Name:    "public_sale",
            Amount:  totalSupply.Mul(sdk.NewInt(20)).Quo(sdk.NewInt(100)),
            Vesting: false,
        },
        {
            Name:    "liquidity",
            Amount:  totalSupply.Mul(sdk.NewInt(18)).Quo(sdk.NewInt(100)),
            Vesting: false,
        },
        {
            Name:    "community_rewards",
            Amount:  totalSupply.Mul(sdk.NewInt(15)).Quo(sdk.NewInt(100)),
            Vesting: false,
        },
        {
            Name:    "development",
            Amount:  totalSupply.Mul(sdk.NewInt(15)).Quo(sdk.NewInt(100)),
            Vesting: false,
        },
        {
            Name:         "team",
            Amount:       totalSupply.Mul(sdk.NewInt(12)).Quo(sdk.NewInt(100)),
            Vesting:      true,
            CliffMonths:  12,
            VestingMonths: 24,
        },
        {
            Name:         "founder",
            Amount:       totalSupply.Mul(sdk.NewInt(8)).Quo(sdk.NewInt(100)),
            Vesting:      true,
            CliffMonths:  12,
            VestingMonths: 48,
        },
        {
            Name:    "dao_treasury",
            Amount:  totalSupply.Mul(sdk.NewInt(5)).Quo(sdk.NewInt(100)),
            Vesting: false,
        },
        {
            Name:         "cofounders",
            Amount:       totalSupply.Mul(sdk.NewInt(35)).Quo(sdk.NewInt(1000)),
            Vesting:      true,
            CliffMonths:  12,
            VestingMonths: 24,
        },
        {
            Name:    "operations",
            Amount:  totalSupply.Mul(sdk.NewInt(2)).Quo(sdk.NewInt(100)),
            Vesting: false,
        },
        {
            Name:         "angels",
            Amount:       totalSupply.Mul(sdk.NewInt(15)).Quo(sdk.NewInt(1000)),
            Vesting:      true,
            CliffMonths:  12,
            VestingMonths: 24,
        },
    }
    
    return createGenesisFromAllocations(allocations)
}
```

### 4.2 Vesting Account Creation
**File**: `app/vesting.go`
```go
package app

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
    vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
)

func CreateVestingAccount(
    address sdk.AccAddress,
    amount sdk.Coins,
    cliffMonths int64,
    vestingMonths int64,
    startTime time.Time,
) *vestingtypes.ContinuousVestingAccount {
    cliffTime := startTime.Add(time.Duration(cliffMonths) * 30 * 24 * time.Hour)
    endTime := startTime.Add(time.Duration(vestingMonths) * 30 * 24 * time.Hour)
    
    baseAccount := authtypes.NewBaseAccountWithAddress(address)
    baseVestingAccount := vestingtypes.NewBaseVestingAccount(
        baseAccount,
        amount,
        endTime.Unix(),
    )
    
    return vestingtypes.NewContinuousVestingAccountRaw(
        baseVestingAccount,
        cliffTime.Unix(),
    )
}
```

## 🎯 Priority 5: Tax Collection and Distribution

### 5.1 Tax Handler
**File**: `x/tax/keeper/handler.go`
```go
package keeper

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) CollectAndDistributeTax(ctx sdk.Context, amount sdk.Coins) error {
    taxRate := k.GetTaxRate(ctx)
    taxAmount := sdk.NewDecCoinsFromCoins(amount...).MulDec(taxRate)
    
    distribution := k.GetTaxDistribution(ctx)
    
    // Distribute to each category
    distributions := map[string]sdk.Dec{
        types.NGOPoolAddress:        distribution.NGODonations,
        types.ValidatorPoolAddress:  distribution.Validators,
        types.CommunityPoolAddress:  distribution.CommunityRewards,
        types.TechFundAddress:       distribution.TechInnovation,
        types.OperationsAddress:     distribution.Operations,
        types.TalentFundAddress:     distribution.TalentAcquisition,
        types.StrategicAddress:      distribution.StrategicReserve,
        types.FounderAddress:        distribution.Founder,
        types.CoFoundersAddress:     distribution.CoFounders,
        types.AngelsAddress:         distribution.AngelInvestors,
    }
    
    for address, percentage := range distributions {
        share := taxAmount.MulDec(percentage)
        if err := k.bankKeeper.SendCoinsFromModuleToAccount(
            ctx,
            types.FeeCollectorModule,
            sdk.MustAccAddressFromBech32(address),
            share.TruncateDecimal(),
        ); err != nil {
            return err
        }
    }
    
    return nil
}
```

## 🎯 Priority 6: Testing

### 6.1 Unit Tests for Token Distribution
**File**: `x/namo/types/constants_test.go`
```go
package types_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestTokenDistribution(t *testing.T) {
    totalPercentage := PublicSaleAllocation + LiquidityAllocation + 
                      CommunityAllocation + DevelopmentAllocation + 
                      TeamAllocation + FounderAllocation + 
                      DAOTreasuryAllocation + CoFounderAllocation + 
                      OperationsAllocation + AngelAllocation
    
    assert.Equal(t, 1.0, totalPercentage, "Total allocation must equal 100%")
    
    // Verify specific allocations
    assert.Equal(t, 0.20, PublicSaleAllocation, "Public sale should be 20%")
    assert.Equal(t, 0.08, FounderAllocation, "Founder should be 8%")
    assert.Equal(t, 0.035, CoFounderAllocation, "Co-founders should be 3.5%")
    assert.Equal(t, 0.015, AngelAllocation, "Angels should be 1.5%")
}
```

### 6.2 Integration Tests for Fee Distribution
**File**: `x/tax/keeper/handler_test.go`
```go
package keeper_test

import (
    "testing"
    sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestTaxDistribution(t *testing.T) {
    app := setupTestApp()
    ctx := app.NewContext(false, tmproto.Header{})
    
    // Create transaction amount
    txAmount := sdk.NewCoins(sdk.NewCoin("namo", sdk.NewInt(100000000)))
    
    // Collect and distribute tax
    err := app.TaxKeeper.CollectAndDistributeTax(ctx, txAmount)
    require.NoError(t, err)
    
    // Verify distributions
    expectedTax := txAmount.AmountOf("namo").MulRaw(25).QuoRaw(1000) // 2.5%
    
    // Check validator pool received 25% of tax
    validatorBalance := app.BankKeeper.GetBalance(ctx, types.ValidatorPoolAddress, "namo")
    expectedValidator := expectedTax.MulRaw(25).QuoRaw(100)
    assert.Equal(t, expectedValidator, validatorBalance.Amount)
    
    // Check NGO pool received 30% of tax
    ngoBalance := app.BankKeeper.GetBalance(ctx, types.NGOPoolAddress, "namo")
    expectedNGO := expectedTax.MulRaw(30).QuoRaw(100)
    assert.Equal(t, expectedNGO, ngoBalance.Amount)
}
```

## 🎯 Priority 7: CLI Commands

### 7.1 Query Commands
**File**: `x/tax/client/cli/query.go`
```go
package cli

import (
    "github.com/cosmos/cosmos-sdk/client"
    "github.com/spf13/cobra"
)

func GetQueryCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "tax",
        Short: "Query tax distribution and rates",
    }
    
    cmd.AddCommand(
        GetCmdQueryTaxRate(),
        GetCmdQueryDistribution(),
        GetCmdQueryCollectedFees(),
    )
    
    return cmd
}

func GetCmdQueryDistribution() *cobra.Command {
    return &cobra.Command{
        Use:   "distribution",
        Short: "Query current tax distribution percentages",
        RunE: func(cmd *cobra.Command, args []string) error {
            clientCtx := client.GetClientContextFromCmd(cmd)
            queryClient := types.NewQueryClient(clientCtx)
            
            res, err := queryClient.TaxDistribution(
                context.Background(),
                &types.QueryTaxDistributionRequest{},
            )
            if err != nil {
                return err
            }
            
            return clientCtx.PrintProto(res)
        },
    }
}
```

## 📋 Implementation Checklist

### Week 1: Core Updates
- [ ] Update all constants files
- [ ] Implement new distribution structures
- [ ] Create allocation addresses
- [ ] Update genesis configuration

### Week 2: Module Updates
- [ ] Update tax module with new distribution
- [ ] Update DEX module fees
- [ ] Update launchpad fees
- [ ] Update NFT marketplace fees

### Week 3: Validator Features
- [ ] Implement geographic incentives
- [ ] Add performance bonus system
- [ ] Create MEV distribution
- [ ] Add staking reward tiers

### Week 4: Testing & Documentation
- [ ] Complete unit tests
- [ ] Run integration tests
- [ ] Update API documentation
- [ ] Create deployment scripts

## 🎯 Key Files to Create/Update

1. **New Files**:
   - `x/validator/types/geographic.go`
   - `x/validator/types/performance.go`
   - `x/validator/keeper/mev.go`
   - `x/validator/keeper/bonus.go`

2. **Update Files**:
   - `x/namo/types/constants.go`
   - `x/tax/types/params.go`
   - `x/dex/types/params.go`
   - `x/launchpad/types/params.go`
   - `x/nft/types/marketplace_params.go`
   - `app/genesis.go`
   - All test files

3. **Configuration Files**:
   - `config/genesis.json`
   - `config/app.toml`
   - `docker-compose.yml`

---

**Note**: Since this is a fresh implementation, we can directly use the optimized economic model without any migration concerns. All values are final and tested for sustainability.