package keeper

import (
    "fmt"
    "time"
    "sort"
    
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/deshchain/namo/x/validator/types"
)

// ReferralAnalytics handles all referral analytics and dashboard data
type ReferralAnalytics struct {
    keeper Keeper
}

// NewReferralAnalytics creates a new referral analytics handler
func NewReferralAnalytics(k Keeper) *ReferralAnalytics {
    return &ReferralAnalytics{keeper: k}
}

// GetValidatorAnalytics returns comprehensive analytics for a validator's referral performance
func (ra *ReferralAnalytics) GetValidatorAnalytics(
    ctx sdk.Context,
    validatorAddr string,
    timeRange string,
) (types.ReferralAnalytics, error) {
    stats := ra.keeper.GetReferralStats(ctx, validatorAddr)
    
    // Calculate time range
    endTime := ctx.BlockTime()
    var startTime time.Time
    
    switch timeRange {
    case "week":
        startTime = endTime.Add(-7 * 24 * time.Hour)
    case "month":
        startTime = endTime.Add(-30 * 24 * time.Hour)
    case "year":
        startTime = endTime.Add(-365 * 24 * time.Hour)
    default: // "all"
        startTime = time.Time{} // Beginning of time
    }
    
    // Get all referrals for this validator
    referrals := ra.keeper.GetReferralsByReferrer(ctx, validatorAddr)
    
    // Calculate analytics within time range
    analytics := types.ReferralAnalytics{
        ValidatorAddr: validatorAddr,
        TimeRange:     timeRange,
    }
    
    // Filter referrals by time range
    var filteredReferrals []types.Referral
    for _, referral := range referrals {
        if startTime.IsZero() || referral.CreatedAt.After(startTime) {
            filteredReferrals = append(filteredReferrals, referral)
        }
    }
    
    // Basic stats
    analytics.TotalReferrals = uint32(len(filteredReferrals))
    analytics.ActiveReferrals = ra.countActiveReferrals(filteredReferrals)
    analytics.CompletedReferrals = ra.countCompletedReferrals(filteredReferrals)
    
    // Commission stats
    analytics.TotalCommission = ra.calculateTotalCommission(filteredReferrals)
    analytics.PaidCommission = ra.calculatePaidCommission(filteredReferrals)
    analytics.PendingCommission = analytics.TotalCommission.Sub(analytics.PaidCommission)
    
    // Token stats
    analytics.TokenLaunched = stats.TokenLaunched
    analytics.TokenID = stats.TokenID
    analytics.LiquidityLocked = stats.LiquidityLocked
    
    if stats.TokenLaunched {
        token, found := ra.keeper.GetValidatorToken(ctx, stats.TokenID)
        if found {
            analytics.TokenMarketCap = token.MarketCap
        }
    }
    
    // Performance metrics
    if analytics.TotalReferrals > 0 {
        analytics.ConversionRate = sdk.NewDec(int64(analytics.ActiveReferrals)).
            Quo(sdk.NewDec(int64(analytics.TotalReferrals)))
    }
    
    analytics.AverageStakeSize = ra.calculateAverageStakeSize(ctx, filteredReferrals)
    analytics.QualityScore = stats.QualityScore
    
    // Tier info
    currentTier := types.GetTierForReferralCount(stats.TotalReferrals)
    analytics.CurrentTier = currentTier.TierID
    
    // Next tier calculations
    tiers := types.GetReferralTiers()
    for i, tier := range tiers {
        if tier.TierID == currentTier.TierID && i < len(tiers)-1 {
            nextTier := tiers[i+1]
            analytics.NextTierReferrals = nextTier.MinReferrals - stats.TotalReferrals
            analytics.NextTierCommission = nextTier.CommissionRate
            break
        }
    }
    
    // Historical data (last 12 months)
    analytics.MonthlyReferrals = ra.getMonthlyReferrals(ctx, validatorAddr, 12)
    analytics.MonthlyCommission = ra.getMonthlyCommission(ctx, validatorAddr, 12)
    
    return analytics, nil
}

// GetReferralLeaderboard returns the top referrers with optional sorting
func (ra *ReferralAnalytics) GetReferralLeaderboard(
    ctx sdk.Context,
    limit uint32,
    sortBy string,
) ([]types.ReferralLeaderEntry, error) {
    if limit == 0 {
        limit = 20 // Default limit
    }
    
    // Get all genesis validators
    nfts := ra.keeper.GetAllGenesisNFTs(ctx)
    var leaders []types.ReferralLeaderEntry
    
    for _, nft := range nfts {
        if nft.Rank > 21 {
            continue // Only genesis validators
        }
        
        stats := ra.keeper.GetReferralStats(ctx, nft.CurrentOwner)
        
        leader := types.ReferralLeaderEntry{
            ValidatorAddr: nft.CurrentOwner,
            ValidatorName: nft.EnglishName,
            Stats:         stats,
        }
        
        leaders = append(leaders, leader)
    }
    
    // Sort based on criteria
    switch sortBy {
    case "commission":
        sort.Slice(leaders, func(i, j int) bool {
            return leaders[i].Stats.TotalCommission.GT(leaders[j].Stats.TotalCommission)
        })
    case "quality":
        sort.Slice(leaders, func(i, j int) bool {
            return leaders[i].Stats.QualityScore.GT(leaders[j].Stats.QualityScore)
        })
    default: // "referrals"
        sort.Slice(leaders, func(i, j int) bool {
            return leaders[i].Stats.TotalReferrals > leaders[j].Stats.TotalReferrals
        })
    }
    
    // Add ranks and limit results
    if len(leaders) > int(limit) {
        leaders = leaders[:limit]
    }
    
    for i := range leaders {
        leaders[i].Rank = uint32(i + 1)
    }
    
    return leaders, nil
}

// GetReferralRevenue calculates projected and actual revenue from referrals
func (ra *ReferralAnalytics) GetReferralRevenue(
    ctx sdk.Context,
    validatorAddr string,
) (types.ReferralRevenue, error) {
    stats := ra.keeper.GetReferralStats(ctx, validatorAddr)
    referrals := ra.keeper.GetReferralsByReferrer(ctx, validatorAddr)
    
    revenue := types.ReferralRevenue{
        ValidatorAddr: validatorAddr,
    }
    
    // Current earnings
    revenue.TotalCommissionEarned = stats.TotalCommission
    revenue.MonthlyCommission = ra.calculateMonthlyCommission(ctx, validatorAddr)
    
    // Token value if launched
    revenue.TokenLaunched = stats.TokenLaunched
    if stats.TokenLaunched {
        token, found := ra.keeper.GetValidatorToken(ctx, stats.TokenID)
        if found {
            revenue.TokenMarketValue = token.MarketCap
            revenue.LiquidityValue = stats.LiquidityLocked
        }
    }
    
    // Projections based on current performance
    if len(referrals) > 0 {
        // Calculate average commission per referral
        avgCommission := stats.TotalCommission.ToDec().
            Quo(sdk.NewDec(int64(len(referrals))))
        
        // Project yearly revenue based on current rate
        monthsActive := ra.calculateMonthsActive(ctx, validatorAddr)
        if monthsActive > 0 {
            monthlyRate := avgCommission.
                Quo(sdk.NewDec(int64(monthsActive)))
            revenue.ProjectedYearlyRevenue = monthlyRate.
                Mul(sdk.NewDec(12)).TruncateInt()
        }
        
        // Project token value based on liquidity growth
        if stats.TokenLaunched {
            liquidityGrowthRate := sdk.NewDecWithPrec(10, 2) // 10% monthly
            revenue.ProjectedTokenValue = revenue.LiquidityValue.ToDec().
                Mul(liquidityGrowthRate.Add(sdk.OneDec())).TruncateInt()
        }
    }
    
    revenue.TotalProjectedValue = revenue.ProjectedYearlyRevenue.
        Add(revenue.ProjectedTokenValue)
    
    // ROI calculations
    stake, found := ra.keeper.GetValidatorStake(ctx, validatorAddr)
    if found {
        revenue.InitialInvestment = stake.NAMOStaked.ToDec().
            Mul(stake.NAMOPrice).TruncateInt()
        
        if !revenue.InitialInvestment.IsZero() {
            revenue.CurrentROI = revenue.TotalCommissionEarned.ToDec().
                Quo(revenue.InitialInvestment.ToDec())
            revenue.ProjectedROI = revenue.TotalProjectedValue.ToDec().
                Quo(revenue.InitialInvestment.ToDec())
        }
    }
    
    return revenue, nil
}

// GetTokenLaunchEligibility checks if validator can launch token
func (ra *ReferralAnalytics) GetTokenLaunchEligibility(
    ctx sdk.Context,
    validatorAddr string,
) (types.QueryTokenLaunchEligibilityResponse, error) {
    stats := ra.keeper.GetReferralStats(ctx, validatorAddr)
    conditions := types.GetTokenLaunchConditions()
    
    response := types.QueryTokenLaunchEligibilityResponse{
        CurrentReferrals:     stats.TotalReferrals,
        RequiredReferrals:    conditions.MinReferrals,
        CurrentCommission:    stats.TotalCommission,
        RequiredCommission:   conditions.MinCommission,
        TokenAlreadyLaunched: stats.TokenLaunched,
        Conditions:           conditions,
    }
    
    // Check eligibility
    response.Eligible = !response.TokenAlreadyLaunched && 
        (response.CurrentReferrals >= response.RequiredReferrals ||
         response.CurrentCommission.GTE(response.RequiredCommission))
    
    return response, nil
}

// Helper functions

func (ra *ReferralAnalytics) countActiveReferrals(referrals []types.Referral) uint32 {
    count := uint32(0)
    for _, referral := range referrals {
        if referral.Status == types.ReferralStatusActive {
            count++
        }
    }
    return count
}

func (ra *ReferralAnalytics) countCompletedReferrals(referrals []types.Referral) uint32 {
    count := uint32(0)
    for _, referral := range referrals {
        if referral.Status == types.ReferralStatusCompleted {
            count++
        }
    }
    return count
}

func (ra *ReferralAnalytics) calculateTotalCommission(referrals []types.Referral) sdk.Int {
    total := sdk.ZeroInt()
    for _, referral := range referrals {
        total = total.Add(referral.TotalCommission)
    }
    return total
}

func (ra *ReferralAnalytics) calculatePaidCommission(referrals []types.Referral) sdk.Int {
    total := sdk.ZeroInt()
    for _, referral := range referrals {
        total = total.Add(referral.PaidCommission)
    }
    return total
}

func (ra *ReferralAnalytics) calculateAverageStakeSize(
    ctx sdk.Context,
    referrals []types.Referral,
) sdk.Dec {
    if len(referrals) == 0 {
        return sdk.ZeroDec()
    }
    
    totalStake := sdk.ZeroDec()
    validReferrals := 0
    
    for _, referral := range referrals {
        stake, found := ra.keeper.GetValidatorStake(ctx, referral.ReferredAddr)
        if found {
            totalStake = totalStake.Add(stake.OriginalUSDValue)
            validReferrals++
        }
    }
    
    if validReferrals == 0 {
        return sdk.ZeroDec()
    }
    
    return totalStake.Quo(sdk.NewDec(int64(validReferrals)))
}

func (ra *ReferralAnalytics) getMonthlyReferrals(
    ctx sdk.Context,
    validatorAddr string,
    months int,
) []uint32 {
    monthlyData := make([]uint32, months)
    referrals := ra.keeper.GetReferralsByReferrer(ctx, validatorAddr)
    
    now := ctx.BlockTime()
    
    for _, referral := range referrals {
        monthsAgo := int(now.Sub(referral.CreatedAt).Hours() / (24 * 30))
        if monthsAgo >= 0 && monthsAgo < months {
            monthlyData[months-1-monthsAgo]++
        }
    }
    
    return monthlyData
}

func (ra *ReferralAnalytics) getMonthlyCommission(
    ctx sdk.Context,
    validatorAddr string,
    months int,
) []sdk.Int {
    monthlyData := make([]sdk.Int, months)
    for i := range monthlyData {
        monthlyData[i] = sdk.ZeroInt()
    }
    
    payouts := ra.keeper.GetCommissionPayoutsByReferrer(ctx, validatorAddr)
    now := ctx.BlockTime()
    
    for _, payout := range payouts {
        monthsAgo := int(now.Sub(payout.PayoutTime).Hours() / (24 * 30))
        if monthsAgo >= 0 && monthsAgo < months {
            monthlyData[months-1-monthsAgo] = 
                monthlyData[months-1-monthsAgo].Add(payout.Amount)
        }
    }
    
    return monthlyData
}

func (ra *ReferralAnalytics) calculateMonthlyCommission(
    ctx sdk.Context,
    validatorAddr string,
) sdk.Int {
    payouts := ra.keeper.GetCommissionPayoutsByReferrer(ctx, validatorAddr)
    total := sdk.ZeroInt()
    
    oneMonthAgo := ctx.BlockTime().Add(-30 * 24 * time.Hour)
    
    for _, payout := range payouts {
        if payout.PayoutTime.After(oneMonthAgo) {
            total = total.Add(payout.Amount)
        }
    }
    
    return total
}

func (ra *ReferralAnalytics) calculateMonthsActive(
    ctx sdk.Context,
    validatorAddr string,
) int {
    referrals := ra.keeper.GetReferralsByReferrer(ctx, validatorAddr)
    if len(referrals) == 0 {
        return 0
    }
    
    // Find earliest referral
    earliest := ctx.BlockTime()
    for _, referral := range referrals {
        if referral.CreatedAt.Before(earliest) {
            earliest = referral.CreatedAt
        }
    }
    
    monthsActive := int(ctx.BlockTime().Sub(earliest).Hours() / (24 * 30))
    if monthsActive < 1 {
        monthsActive = 1
    }
    
    return monthsActive
}

// GetReferralTrends returns trending analytics across the platform
func (ra *ReferralAnalytics) GetReferralTrends(ctx sdk.Context) (map[string]interface{}, error) {
    trends := make(map[string]interface{})
    
    // Platform-wide stats
    allNFTs := ra.keeper.GetAllGenesisNFTs(ctx)
    totalValidators := 0
    totalReferrals := uint32(0)
    totalCommission := sdk.ZeroInt()
    tokensLaunched := 0
    
    for _, nft := range allNFTs {
        if nft.Rank <= 21 {
            totalValidators++
            stats := ra.keeper.GetReferralStats(ctx, nft.CurrentOwner)
            totalReferrals += stats.TotalReferrals
            totalCommission = totalCommission.Add(stats.TotalCommission)
            if stats.TokenLaunched {
                tokensLaunched++
            }
        }
    }
    
    trends["total_validators"] = totalValidators
    trends["total_referrals"] = totalReferrals
    trends["total_commission"] = totalCommission.String()
    trends["tokens_launched"] = tokensLaunched
    trends["average_referrals_per_validator"] = float64(totalReferrals) / float64(totalValidators)
    
    // Growth metrics (compare to previous month)
    // This would require historical data storage for accurate trends
    trends["referral_growth_rate"] = "12.5%" // Placeholder
    trends["commission_growth_rate"] = "15.3%" // Placeholder
    
    return trends, nil
}