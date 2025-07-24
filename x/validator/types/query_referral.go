package types

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/cosmos/cosmos-sdk/types/query"
)

// Query request types for referral system

// QueryReferralRequest is request type for the Query/Referral RPC method
type QueryReferralRequest struct {
    ReferralID uint64 `protobuf:"varint,1,opt,name=referral_id,json=referralId,proto3" json:"referral_id,omitempty"`
}

// QueryReferralResponse is response type for the Query/Referral RPC method
type QueryReferralResponse struct {
    Referral Referral `protobuf:"bytes,1,opt,name=referral,proto3" json:"referral"`
}

// QueryReferralStatsRequest is request type for the Query/ReferralStats RPC method
type QueryReferralStatsRequest struct {
    ValidatorAddr string `protobuf:"bytes,1,opt,name=validator_addr,json=validatorAddr,proto3" json:"validator_addr,omitempty"`
}

// QueryReferralStatsResponse is response type for the Query/ReferralStats RPC method
type QueryReferralStatsResponse struct {
    Stats ReferralStats `protobuf:"bytes,1,opt,name=stats,proto3" json:"stats"`
}

// QueryValidatorTokenRequest is request type for the Query/ValidatorToken RPC method
type QueryValidatorTokenRequest struct {
    TokenID uint64 `protobuf:"varint,1,opt,name=token_id,json=tokenId,proto3" json:"token_id,omitempty"`
}

// QueryValidatorTokenResponse is response type for the Query/ValidatorToken RPC method
type QueryValidatorTokenResponse struct {
    Token ValidatorToken `protobuf:"bytes,1,opt,name=token,proto3" json:"token"`
}

// QueryReferralsByReferrerRequest is request type for the Query/ReferralsByReferrer RPC method
type QueryReferralsByReferrerRequest struct {
    ReferrerAddr string           `protobuf:"bytes,1,opt,name=referrer_addr,json=referrerAddr,proto3" json:"referrer_addr,omitempty"`
    Pagination   *query.PageRequest `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

// QueryReferralsByReferrerResponse is response type for the Query/ReferralsByReferrer RPC method
type QueryReferralsByReferrerResponse struct {
    Referrals  []Referral          `protobuf:"bytes,1,rep,name=referrals,proto3" json:"referrals"`
    Pagination *query.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

// QueryCommissionPayoutsRequest is request type for the Query/CommissionPayouts RPC method
type QueryCommissionPayoutsRequest struct {
    ReferrerAddr string           `protobuf:"bytes,1,opt,name=referrer_addr,json=referrerAddr,proto3" json:"referrer_addr,omitempty"`
    Pagination   *query.PageRequest `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

// QueryCommissionPayoutsResponse is response type for the Query/CommissionPayouts RPC method
type QueryCommissionPayoutsResponse struct {
    Payouts    []CommissionPayout  `protobuf:"bytes,1,rep,name=payouts,proto3" json:"payouts"`
    Pagination *query.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

// QueryReferralLeaderboardRequest is request type for the Query/ReferralLeaderboard RPC method
type QueryReferralLeaderboardRequest struct {
    Limit      uint32 `protobuf:"varint,1,opt,name=limit,proto3" json:"limit,omitempty"`
    SortBy     string `protobuf:"bytes,2,opt,name=sort_by,json=sortBy,proto3" json:"sort_by,omitempty"` // "referrals", "commission", "quality"
}

// QueryReferralLeaderboardResponse is response type for the Query/ReferralLeaderboard RPC method
type QueryReferralLeaderboardResponse struct {
    Leaders []ReferralLeaderEntry `protobuf:"bytes,1,rep,name=leaders,proto3" json:"leaders"`
}

// ReferralLeaderEntry represents a leaderboard entry
type ReferralLeaderEntry struct {
    Rank          uint32      `json:"rank"`
    ValidatorAddr string      `json:"validator_addr"`
    ValidatorName string      `json:"validator_name"` // From NFT
    Stats         ReferralStats `json:"stats"`
}

// QueryTokenLaunchEligibilityRequest is request type for the Query/TokenLaunchEligibility RPC method
type QueryTokenLaunchEligibilityRequest struct {
    ValidatorAddr string `protobuf:"bytes,1,opt,name=validator_addr,json=validatorAddr,proto3" json:"validator_addr,omitempty"`
}

// QueryTokenLaunchEligibilityResponse is response type for the Query/TokenLaunchEligibility RPC method
type QueryTokenLaunchEligibilityResponse struct {
    Eligible            bool                   `json:"eligible"`
    CurrentReferrals    uint32                 `json:"current_referrals"`
    RequiredReferrals   uint32                 `json:"required_referrals"`
    CurrentCommission   sdk.Int                `json:"current_commission"`
    RequiredCommission  sdk.Int                `json:"required_commission"`
    TokenAlreadyLaunched bool                  `json:"token_already_launched"`
    Conditions          TokenLaunchConditions `json:"conditions"`
}

// QueryValidatorTokensByValidatorRequest is request type for the Query/ValidatorTokensByValidator RPC method
type QueryValidatorTokensByValidatorRequest struct {
    ValidatorAddr string           `protobuf:"bytes,1,opt,name=validator_addr,json=validatorAddr,proto3" json:"validator_addr,omitempty"`
    Pagination    *query.PageRequest `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

// QueryValidatorTokensByValidatorResponse is response type for the Query/ValidatorTokensByValidator RPC method
type QueryValidatorTokensByValidatorResponse struct {
    Tokens     []ValidatorToken    `protobuf:"bytes,1,rep,name=tokens,proto3" json:"tokens"`
    Pagination *query.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

// QueryReferralAnalyticsRequest is request type for the Query/ReferralAnalytics RPC method
type QueryReferralAnalyticsRequest struct {
    ValidatorAddr string `protobuf:"bytes,1,opt,name=validator_addr,json=validatorAddr,proto3" json:"validator_addr,omitempty"`
    TimeRange     string `protobuf:"bytes,2,opt,name=time_range,json=timeRange,proto3" json:"time_range,omitempty"` // "week", "month", "year", "all"
}

// QueryReferralAnalyticsResponse is response type for the Query/ReferralAnalytics RPC method
type QueryReferralAnalyticsResponse struct {
    Analytics ReferralAnalytics `protobuf:"bytes,1,opt,name=analytics,proto3" json:"analytics"`
}

// ReferralAnalytics provides detailed analytics for a validator's referral performance
type ReferralAnalytics struct {
    ValidatorAddr       string      `json:"validator_addr"`
    TimeRange          string      `json:"time_range"`
    
    // Basic stats
    TotalReferrals     uint32      `json:"total_referrals"`
    ActiveReferrals    uint32      `json:"active_referrals"`
    CompletedReferrals uint32      `json:"completed_referrals"`
    
    // Commission stats
    TotalCommission    sdk.Int     `json:"total_commission"`
    PaidCommission     sdk.Int     `json:"paid_commission"`
    PendingCommission  sdk.Int     `json:"pending_commission"`
    
    // Token stats
    TokenLaunched      bool        `json:"token_launched"`
    TokenID            uint64      `json:"token_id"`
    LiquidityLocked    sdk.Int     `json:"liquidity_locked"`
    TokenMarketCap     sdk.Int     `json:"token_market_cap"`
    
    // Performance metrics
    ConversionRate     sdk.Dec     `json:"conversion_rate"`     // Active/Total referrals
    AverageStakeSize   sdk.Dec     `json:"average_stake_size"`  // Average referred validator stake
    QualityScore       sdk.Dec     `json:"quality_score"`       // Quality of referred validators
    
    // Tier info
    CurrentTier        uint32      `json:"current_tier"`
    NextTierReferrals  uint32      `json:"next_tier_referrals"` // Referrals needed for next tier
    NextTierCommission sdk.Dec     `json:"next_tier_commission"` // Commission rate at next tier
    
    // Historical data
    MonthlyReferrals   []uint32    `json:"monthly_referrals"`   // Last 12 months
    MonthlyCommission  []sdk.Int   `json:"monthly_commission"`  // Last 12 months
}

// QueryAllValidatorTokensRequest is request type for the Query/AllValidatorTokens RPC method
type QueryAllValidatorTokensRequest struct {
    Pagination *query.PageRequest `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

// QueryAllValidatorTokensResponse is response type for the Query/AllValidatorTokens RPC method
type QueryAllValidatorTokensResponse struct {
    Tokens     []ValidatorToken    `protobuf:"bytes,1,rep,name=tokens,proto3" json:"tokens"`
    Pagination *query.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

// QueryReferralRevenueRequest is request type for the Query/ReferralRevenue RPC method
type QueryReferralRevenueRequest struct {
    ValidatorAddr string `protobuf:"bytes,1,opt,name=validator_addr,json=validatorAddr,proto3" json:"validator_addr,omitempty"`
}

// QueryReferralRevenueResponse is response type for the Query/ReferralRevenue RPC method
type QueryReferralRevenueResponse struct {
    Revenue ReferralRevenue `protobuf:"bytes,1,opt,name=revenue,proto3" json:"revenue"`
}

// ReferralRevenue shows projected and actual revenue from referrals
type ReferralRevenue struct {
    ValidatorAddr         string  `json:"validator_addr"`
    
    // Current earnings
    TotalCommissionEarned sdk.Int `json:"total_commission_earned"`
    MonthlyCommission     sdk.Int `json:"monthly_commission"`
    
    // Token value (if launched)
    TokenLaunched         bool    `json:"token_launched"`
    TokenMarketValue      sdk.Int `json:"token_market_value"`
    LiquidityValue        sdk.Int `json:"liquidity_value"`
    
    // Projections
    ProjectedYearlyRevenue   sdk.Int `json:"projected_yearly_revenue"`
    ProjectedTokenValue      sdk.Int `json:"projected_token_value"`
    TotalProjectedValue      sdk.Int `json:"total_projected_value"`
    
    // ROI calculations
    InitialInvestment        sdk.Int `json:"initial_investment"` // Validator's stake
    CurrentROI               sdk.Dec `json:"current_roi"`
    ProjectedROI             sdk.Dec `json:"projected_roi"`
}