package types

import (
    "time"
    sdk "github.com/cosmos/cosmos-sdk/types"
)

// ValidatorStakeRequirement defines the stake requirement for a validator
type ValidatorStakeRequirement struct {
    ValidatorRank     uint32    `json:"validator_rank"`
    RequiredUSD       sdk.Dec   `json:"required_usd"`
    NAMOAmount        sdk.Int   `json:"namo_amount"`
    NAMOPriceUSD      sdk.Dec   `json:"namo_price_usd"`
    StakeTime         time.Time `json:"stake_time"`
    LockPeriodMonths  uint32    `json:"lock_period_months"`
    VestingMonths     uint32    `json:"vesting_months"`
    PerformanceBondPct sdk.Dec  `json:"performance_bond_pct"`
}

// StakeTier defines the tier-based staking parameters
type StakeTier struct {
    TierID            uint32  `json:"tier_id"`
    MinRank           uint32  `json:"min_rank"`
    MaxRank           uint32  `json:"max_rank"`
    LockPeriodMonths  uint32  `json:"lock_period_months"`
    VestingMonths     uint32  `json:"vesting_months"`
    PerformanceBondPct sdk.Dec `json:"performance_bond_pct"`
    DailySellLimitPct sdk.Dec `json:"daily_sell_limit_pct"`
    SlashingMultiplier sdk.Dec `json:"slashing_multiplier"`
}

// ValidatorStake represents a validator's staked tokens
type ValidatorStake struct {
    ValidatorAddr      string    `json:"validator_address"`
    OriginalUSDValue   sdk.Dec   `json:"original_usd_value"`
    NAMOTokensStaked   sdk.Int   `json:"namo_tokens_staked"`
    NAMOPriceAtStake   sdk.Dec   `json:"namo_price_at_stake"`
    StakeTimestamp     time.Time `json:"stake_timestamp"`
    
    // Lock and vesting
    LockEndTime        time.Time `json:"lock_end_time"`
    VestingEndTime     time.Time `json:"vesting_end_time"`
    
    // Performance bond (carved from stake)
    PerformanceBond    sdk.Int   `json:"performance_bond"`
    VestableAmount     sdk.Int   `json:"vestable_amount"`
    UnlockedAmount     sdk.Int   `json:"unlocked_amount"`
    
    // Insurance contribution (2% of stake)
    InsuranceContribution sdk.Int `json:"insurance_contribution"`
    
    // NFT binding
    BoundNFTID         uint64    `json:"bound_nft_id"`
    NFTBindingActive   bool      `json:"nft_binding_active"`
    
    // Tracking
    Tier               uint32    `json:"tier"`
    SlashingHistory    []SlashingEvent `json:"slashing_history"`
    LastWithdrawal     time.Time `json:"last_withdrawal"`
}

// SlashingEvent records a slashing incident
type SlashingEvent struct {
    Timestamp     time.Time `json:"timestamp"`
    Reason        string    `json:"reason"`
    SlashedAmount sdk.Int   `json:"slashed_amount"`
    SlashingRate  sdk.Dec   `json:"slashing_rate"`
}

// GetStakeTiers returns the three stake tiers
func GetStakeTiers() []StakeTier {
    return []StakeTier{
        {
            TierID:            1,
            MinRank:           1,
            MaxRank:           10,
            LockPeriodMonths:  6,
            VestingMonths:     18,
            PerformanceBondPct: sdk.NewDecWithPrec(20, 2), // 20%
            DailySellLimitPct: sdk.NewDecWithPrec(2, 2),   // 2%
            SlashingMultiplier: sdk.NewDec(1),             // 1x base rate
        },
        {
            TierID:            2,
            MinRank:           11,
            MaxRank:           20,
            LockPeriodMonths:  9,
            VestingMonths:     24,
            PerformanceBondPct: sdk.NewDecWithPrec(25, 2), // 25%
            DailySellLimitPct: sdk.NewDecWithPrec(1, 2),   // 1%
            SlashingMultiplier: sdk.NewDecWithPrec(15, 1), // 1.5x
        },
        {
            TierID:            3,
            MinRank:           21,
            MaxRank:           21,
            LockPeriodMonths:  12,
            VestingMonths:     36,
            PerformanceBondPct: sdk.NewDecWithPrec(30, 2), // 30%
            DailySellLimitPct: sdk.NewDecWithPrec(5, 3),   // 0.5%
            SlashingMultiplier: sdk.NewDec(2),             // 2x
        },
    }
}

// GetValidatorStakeRequirements returns stake requirements for all 21 validators
func GetValidatorStakeRequirements() []ValidatorStakeRequirement {
    return []ValidatorStakeRequirement{
        // Tier 1: Validators 1-10
        {ValidatorRank: 1, RequiredUSD: sdk.NewDec(200000)},
        {ValidatorRank: 2, RequiredUSD: sdk.NewDec(220000)},
        {ValidatorRank: 3, RequiredUSD: sdk.NewDec(240000)},
        {ValidatorRank: 4, RequiredUSD: sdk.NewDec(260000)},
        {ValidatorRank: 5, RequiredUSD: sdk.NewDec(280000)},
        {ValidatorRank: 6, RequiredUSD: sdk.NewDec(300000)},
        {ValidatorRank: 7, RequiredUSD: sdk.NewDec(320000)},
        {ValidatorRank: 8, RequiredUSD: sdk.NewDec(340000)},
        {ValidatorRank: 9, RequiredUSD: sdk.NewDec(360000)},
        {ValidatorRank: 10, RequiredUSD: sdk.NewDec(380000)},
        
        // Tier 2: Validators 11-20
        {ValidatorRank: 11, RequiredUSD: sdk.NewDec(800000)},
        {ValidatorRank: 12, RequiredUSD: sdk.NewDec(820000)},
        {ValidatorRank: 13, RequiredUSD: sdk.NewDec(840000)},
        {ValidatorRank: 14, RequiredUSD: sdk.NewDec(860000)},
        {ValidatorRank: 15, RequiredUSD: sdk.NewDec(880000)},
        {ValidatorRank: 16, RequiredUSD: sdk.NewDec(900000)},
        {ValidatorRank: 17, RequiredUSD: sdk.NewDec(920000)},
        {ValidatorRank: 18, RequiredUSD: sdk.NewDec(940000)},
        {ValidatorRank: 19, RequiredUSD: sdk.NewDec(960000)},
        {ValidatorRank: 20, RequiredUSD: sdk.NewDec(980000)},
        
        // Tier 3: Validator 21
        {ValidatorRank: 21, RequiredUSD: sdk.NewDec(1500000)},
    }
}

// CalculateNAMORequired calculates NAMO tokens needed based on USD value and current price
func CalculateNAMORequired(usdAmount sdk.Dec, namoPriceUSD sdk.Dec) sdk.Int {
    if namoPriceUSD.IsZero() {
        panic("NAMO price cannot be zero")
    }
    
    // USD amount / NAMO price = NAMO tokens needed
    namoTokens := usdAmount.Quo(namoPriceUSD)
    
    // Convert to integer (with 6 decimal places for NAMO)
    // Multiply by 10^6 to preserve precision
    namoTokensInt := namoTokens.MulInt64(1000000).TruncateInt()
    
    return namoTokensInt
}

// GetTierForRank returns the tier for a given validator rank
func GetTierForRank(rank uint32) (StakeTier, bool) {
    tiers := GetStakeTiers()
    for _, tier := range tiers {
        if rank >= tier.MinRank && rank <= tier.MaxRank {
            return tier, true
        }
    }
    return StakeTier{}, false
}