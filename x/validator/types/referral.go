package types

import (
    "time"
    sdk "github.com/cosmos/cosmos-sdk/types"
)

// ReferralStatus represents the status of a referral
type ReferralStatus string

const (
    ReferralStatusPending   ReferralStatus = "pending"
    ReferralStatusActive    ReferralStatus = "active"
    ReferralStatusCompleted ReferralStatus = "completed"
    ReferralStatusCancelled ReferralStatus = "cancelled"
    ReferralStatusClawback  ReferralStatus = "clawback"
)

// Referral represents a validator referral
type Referral struct {
    ReferralID       uint64         `json:"referral_id"`
    ReferrerAddr     string         `json:"referrer_address"`     // Genesis validator who referred
    ReferredAddr     string         `json:"referred_address"`     // New validator
    ReferredRank     uint32         `json:"referred_rank"`        // Rank of new validator
    Status           ReferralStatus `json:"status"`
    CreatedAt        time.Time      `json:"created_at"`
    ActivatedAt      time.Time      `json:"activated_at"`
    CommissionRate   sdk.Dec        `json:"commission_rate"`      // 10-20% based on tier
    TotalCommission  sdk.Int        `json:"total_commission"`     // Total earned
    PaidCommission   sdk.Int        `json:"paid_commission"`      // Already paid
    LiquidityLocked  sdk.Int        `json:"liquidity_locked"`     // Locked in DEX
    ClawbackPeriod   time.Time      `json:"clawback_period"`      // 1 year from activation
    ClawbackAmount   sdk.Int        `json:"clawback_amount"`      // Amount clawed back
    ClawbackReason   string         `json:"clawback_reason"`      // Reason for clawback
}

// ReferralTier defines commission tiers based on referral count
type ReferralTier struct {
    TierID          uint32  `json:"tier_id"`
    MinReferrals    uint32  `json:"min_referrals"`
    MaxReferrals    uint32  `json:"max_referrals"`
    CommissionRate  sdk.Dec `json:"commission_rate"`
    TokenBonus      sdk.Int `json:"token_bonus"`
    BadgeNFT        string  `json:"badge_nft"`
}

// ValidatorToken represents auto-launched token for validators
type ValidatorToken struct {
    TokenID          uint64    `json:"token_id"`
    ValidatorAddr    string    `json:"validator_address"`
    TokenName        string    `json:"token_name"`        // "[NFT Name] Coin"
    TokenSymbol      string    `json:"token_symbol"`      // Auto-generated
    TotalSupply      sdk.Int   `json:"total_supply"`      // 1 billion
    Decimals         uint32    `json:"decimals"`          // 6
    LogoURI          string    `json:"logo_uri"`          // NFT image
    
    // Distribution
    ValidatorAllocation sdk.Int `json:"validator_allocation"` // 40%
    LiquidityAllocation sdk.Int `json:"liquidity_allocation"` // 30%
    AirdropAllocation   sdk.Int `json:"airdrop_allocation"`   // 15%
    DevelopmentAllocation sdk.Int `json:"development_allocation"` // 10%
    InitialLiquidity    sdk.Int `json:"initial_liquidity"`    // 5%
    
    // Launch conditions
    LaunchedAt       time.Time `json:"launched_at"`
    LaunchTrigger    string    `json:"launch_trigger"` // "referrals" or "commission" or "manual"
    ReferralCount    uint32    `json:"referral_count"`
    CommissionEarned sdk.Int   `json:"commission_earned"`
    
    // DEX Integration
    LiquidityPoolID  string    `json:"liquidity_pool_id"`
    TradingPairID    string    `json:"trading_pair_id"`
    CurrentPrice     sdk.Dec   `json:"current_price"`
    MarketCap        sdk.Int   `json:"market_cap"`
    
    // Anti-dump parameters
    MaxWalletPercent    sdk.Dec `json:"max_wallet_percent"`     // 2%
    MaxTxPercent        sdk.Dec `json:"max_tx_percent"`         // 0.5%
    SellTaxPercent      sdk.Dec `json:"sell_tax_percent"`       // 5%
    BuyTaxPercent       sdk.Dec `json:"buy_tax_percent"`        // 2%
    CooldownSeconds     uint64  `json:"cooldown_seconds"`       // 3600 (1 hour)
}

// ReferralStats tracks validator referral performance
type ReferralStats struct {
    ValidatorAddr       string    `json:"validator_address"`
    TotalReferrals      uint32    `json:"total_referrals"`
    ActiveReferrals     uint32    `json:"active_referrals"`
    TotalCommission     sdk.Int   `json:"total_commission"`
    CurrentTier         uint32    `json:"current_tier"`
    TokenLaunched       bool      `json:"token_launched"`
    TokenID             uint64    `json:"token_id"`
    LiquidityLocked     sdk.Int   `json:"liquidity_locked"`
    LastReferralDate    time.Time `json:"last_referral_date"`
    QualityScore        sdk.Dec   `json:"quality_score"` // Based on referred validator performance
}

// CommissionPayout represents a referral commission payment
type CommissionPayout struct {
    PayoutID        uint64    `json:"payout_id"`
    ReferralID      uint64    `json:"referral_id"`
    ReferrerAddr    string    `json:"referrer_address"`
    Amount          sdk.Int   `json:"amount"`           // In NAMO value
    TokenAmount     sdk.Int   `json:"token_amount"`     // Validator tokens created
    LiquidityAdded  sdk.Int   `json:"liquidity_added"`  // NAMO added to pool
    PayoutTime      time.Time `json:"payout_time"`
    BlockHeight     int64     `json:"block_height"`
}

// GetReferralTiers returns the referral commission tiers
func GetReferralTiers() []ReferralTier {
    return []ReferralTier{
        {
            TierID:         1,
            MinReferrals:   0,
            MaxReferrals:   10,
            CommissionRate: sdk.NewDecWithPrec(10, 2), // 10%
            TokenBonus:     sdk.ZeroInt(),
            BadgeNFT:       "",
        },
        {
            TierID:         2,
            MinReferrals:   11,
            MaxReferrals:   25,
            CommissionRate: sdk.NewDecWithPrec(12, 2), // 12%
            TokenBonus:     sdk.NewInt(1000000000), // 1,000 tokens
            BadgeNFT:       "bronze_recruiter",
        },
        {
            TierID:         3,
            MinReferrals:   26,
            MaxReferrals:   50,
            CommissionRate: sdk.NewDecWithPrec(15, 2), // 15%
            TokenBonus:     sdk.NewInt(5000000000), // 5,000 tokens
            BadgeNFT:       "silver_recruiter",
        },
        {
            TierID:         4,
            MinReferrals:   51,
            MaxReferrals:   100,
            CommissionRate: sdk.NewDecWithPrec(20, 2), // 20%
            TokenBonus:     sdk.NewInt(10000000000), // 10,000 tokens
            BadgeNFT:       "gold_recruiter",
        },
    }
}

// GetTierForReferralCount returns the tier based on referral count
func GetTierForReferralCount(count uint32) ReferralTier {
    tiers := GetReferralTiers()
    for _, tier := range tiers {
        if count >= tier.MinReferrals && count <= tier.MaxReferrals {
            return tier
        }
    }
    // Return highest tier if count exceeds all
    return tiers[len(tiers)-1]
}

// TokenLaunchConditions defines when a validator token can be launched
type TokenLaunchConditions struct {
    MinReferrals      uint32  `json:"min_referrals"`       // 5
    MinCommission     sdk.Int `json:"min_commission"`      // ₹50 lakhs
    MinValidatorAge   int64   `json:"min_validator_age"`   // 1 year in seconds
    RequireActiveNFT  bool    `json:"require_active_nft"`  // Must own genesis NFT
}

// GetTokenLaunchConditions returns the conditions for token launch
func GetTokenLaunchConditions() TokenLaunchConditions {
    return TokenLaunchConditions{
        MinReferrals:     5,
        MinCommission:    sdk.NewInt(5000000000000), // ₹50 lakhs in NAMO (assuming ₹0.01)
        MinValidatorAge:  365 * 24 * 60 * 60,         // 1 year
        RequireActiveNFT: true,
    }
}

// CalculateTokenDistribution calculates token allocations
func CalculateTokenDistribution(totalSupply sdk.Int) ValidatorTokenDistribution {
    return ValidatorTokenDistribution{
        ValidatorAllocation:   totalSupply.MulRaw(40).QuoRaw(100),
        LiquidityAllocation:   totalSupply.MulRaw(30).QuoRaw(100),
        AirdropAllocation:     totalSupply.MulRaw(15).QuoRaw(100),
        DevelopmentAllocation: totalSupply.MulRaw(10).QuoRaw(100),
        InitialLiquidity:      totalSupply.MulRaw(5).QuoRaw(100),
    }
}

// ValidatorTokenDistribution holds token allocation amounts
type ValidatorTokenDistribution struct {
    ValidatorAllocation   sdk.Int
    LiquidityAllocation   sdk.Int
    AirdropAllocation     sdk.Int
    DevelopmentAllocation sdk.Int
    InitialLiquidity      sdk.Int
}