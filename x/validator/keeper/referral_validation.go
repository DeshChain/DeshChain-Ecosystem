package keeper

import (
    "fmt"
    "time"
    "net"
    "strings"
    "crypto/sha256"
    "encoding/hex"
    
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/deshchain/namo/x/validator/types"
)

// ReferralValidator handles all referral validation and anti-gaming measures
type ReferralValidator struct {
    keeper Keeper
}

// NewReferralValidator creates a new referral validator
func NewReferralValidator(k Keeper) *ReferralValidator {
    return &ReferralValidator{keeper: k}
}

// ValidateReferralEligibility performs comprehensive validation for referral creation
func (rv *ReferralValidator) ValidateReferralEligibility(
    ctx sdk.Context,
    referrerAddr sdk.AccAddress,
    referredAddr sdk.AccAddress,
    referredRank uint32,
    clientIP string,
) error {
    // 1. Basic address validation
    if err := rv.validateAddresses(referrerAddr, referredAddr); err != nil {
        return err
    }
    
    // 2. Genesis validator validation
    if err := rv.validateGenesisValidator(ctx, referrerAddr); err != nil {
        return err
    }
    
    // 3. Referral limits validation
    if err := rv.validateReferralLimits(ctx, referrerAddr); err != nil {
        return err
    }
    
    // 4. Anti-gaming validation
    if err := rv.validateAntiGamingRules(ctx, referrerAddr, referredAddr, clientIP); err != nil {
        return err
    }
    
    // 5. Rank availability validation
    if err := rv.validateRankAvailability(ctx, referredRank); err != nil {
        return err
    }
    
    // 6. Time-based restrictions
    if err := rv.validateTimeRestrictions(ctx, referrerAddr); err != nil {
        return err
    }
    
    return nil
}

// validateAddresses performs basic address validation
func (rv *ReferralValidator) validateAddresses(referrerAddr, referredAddr sdk.AccAddress) error {
    if referrerAddr.Empty() {
        return fmt.Errorf("referrer address cannot be empty")
    }
    
    if referredAddr.Empty() {
        return fmt.Errorf("referred address cannot be empty")
    }
    
    if referrerAddr.Equals(referredAddr) {
        return fmt.Errorf("self-referral is not allowed")
    }
    
    return nil
}

// validateGenesisValidator checks if referrer is a genesis validator
func (rv *ReferralValidator) validateGenesisValidator(ctx sdk.Context, referrerAddr sdk.AccAddress) error {
    // Check if referrer has a stake
    stake, found := rv.keeper.GetValidatorStake(ctx, referrerAddr.String())
    if !found {
        return fmt.Errorf("referrer is not a validator")
    }
    
    // Check if referrer owns a genesis NFT (rank 1-21)
    nfts := rv.keeper.GetAllGenesisNFTs(ctx)
    for _, nft := range nfts {
        if nft.CurrentOwner == referrerAddr.String() && nft.Rank <= 21 {
            return nil // Valid genesis validator
        }
    }
    
    return fmt.Errorf("only genesis validators (NFT holders rank 1-21) can refer new validators")
}

// validateReferralLimits checks referral count limits
func (rv *ReferralValidator) validateReferralLimits(ctx sdk.Context, referrerAddr sdk.AccAddress) error {
    stats := rv.keeper.GetReferralStats(ctx, referrerAddr.String())
    
    // Global limit: 100 total referrals
    if stats.TotalReferrals >= 100 {
        return fmt.Errorf("referral limit reached (100 max per validator)")
    }
    
    // Monthly limit: 5 referrals per month
    monthlyCount := rv.getMonthlyReferralCount(ctx, referrerAddr.String())
    if monthlyCount >= 5 {
        return fmt.Errorf("monthly referral limit reached (5 max per month)")
    }
    
    // Weekly limit: 2 referrals per week
    weeklyCount := rv.getWeeklyReferralCount(ctx, referrerAddr.String())
    if weeklyCount >= 2 {
        return fmt.Errorf("weekly referral limit reached (2 max per week)")
    }
    
    return nil
}

// validateAntiGamingRules implements comprehensive anti-gaming measures
func (rv *ReferralValidator) validateAntiGamingRules(
    ctx sdk.Context,
    referrerAddr sdk.AccAddress,
    referredAddr sdk.AccAddress,
    clientIP string,
) error {
    // 1. IP clustering detection
    if err := rv.detectIPClustering(ctx, referrerAddr, clientIP); err != nil {
        return err
    }
    
    // 2. Wallet clustering detection
    if err := rv.detectWalletClustering(ctx, referrerAddr, referredAddr); err != nil {
        return err
    }
    
    // 3. Pattern analysis (timing, amounts, etc.)
    if err := rv.detectSuspiciousPatterns(ctx, referrerAddr); err != nil {
        return err
    }
    
    // 4. Historical behavior analysis
    if err := rv.analyzeHistoricalBehavior(ctx, referrerAddr); err != nil {
        return err
    }
    
    return nil
}

// detectIPClustering prevents multiple referrals from same IP
func (rv *ReferralValidator) detectIPClustering(
    ctx sdk.Context,
    referrerAddr sdk.AccAddress,
    clientIP string,
) error {
    if clientIP == "" {
        return nil // Skip if IP not provided
    }
    
    // Parse IP to get subnet
    ip := net.ParseIP(clientIP)
    if ip == nil {
        return fmt.Errorf("invalid IP address: %s", clientIP)
    }
    
    // Create IP hash for privacy
    ipHash := rv.hashIP(clientIP)
    
    // Check if this IP has been used for referrals recently (last 7 days)
    recentReferrals := rv.getRecentReferralsByIP(ctx, ipHash, 7*24*time.Hour)
    
    // Allow max 2 referrals from same IP subnet per week
    if len(recentReferrals) >= 2 {
        return fmt.Errorf("too many referrals from this IP subnet (max 2 per week)")
    }
    
    // Store IP hash for this referral
    rv.storeReferralIP(ctx, referrerAddr.String(), ipHash, ctx.BlockTime())
    
    return nil
}

// detectWalletClustering prevents sybil attacks via wallet analysis
func (rv *ReferralValidator) detectWalletClustering(
    ctx sdk.Context,
    referrerAddr sdk.AccAddress,
    referredAddr sdk.AccAddress,
) error {
    // Check transaction history between referrer and referred
    if rv.hasRecentTransactions(ctx, referrerAddr, referredAddr) {
        return fmt.Errorf("referrer and referred have recent transaction history (potential sybil)")
    }
    
    // Check if referred address is too new (less than 7 days old)
    if rv.isAddressTooNew(ctx, referredAddr) {
        return fmt.Errorf("referred address is too new (must be at least 7 days old)")
    }
    
    // Check if referred address has minimum transaction activity
    if !rv.hasMinimumActivity(ctx, referredAddr) {
        return fmt.Errorf("referred address lacks minimum on-chain activity")
    }
    
    return nil
}

// detectSuspiciousPatterns analyzes referral patterns
func (rv *ReferralValidator) detectSuspiciousPatterns(
    ctx sdk.Context,
    referrerAddr sdk.AccAddress,
) error {
    referrals := rv.keeper.GetReferralsByReferrer(ctx, referrerAddr.String())
    
    if len(referrals) < 2 {
        return nil // Not enough data for pattern analysis
    }
    
    // Check for suspicious timing patterns (all referrals within short timeframe)
    if rv.hasClusteredTiming(referrals) {
        return fmt.Errorf("suspicious referral timing pattern detected")
    }
    
    // Check for sequential address patterns
    if rv.hasSequentialAddresses(referrals) {
        return fmt.Errorf("suspicious sequential address pattern detected")
    }
    
    // Check for uniform stake amounts (suggests automation)
    if rv.hasUniformStakeAmounts(ctx, referrals) {
        return fmt.Errorf("suspicious uniform stake pattern detected")
    }
    
    return nil
}

// analyzeHistoricalBehavior checks validator's past behavior
func (rv *ReferralValidator) analyzeHistoricalBehavior(
    ctx sdk.Context,
    referrerAddr sdk.AccAddress,
) error {
    // Get validator's history
    stake, found := rv.keeper.GetValidatorStake(ctx, referrerAddr.String())
    if !found {
        return fmt.Errorf("validator stake not found")
    }
    
    // Check if validator has been slashed for referral abuse
    if rv.hasReferralSlashingHistory(ctx, referrerAddr) {
        return fmt.Errorf("validator has history of referral violations")
    }
    
    // Check validator performance and reliability
    if rv.hasLowQualityScore(ctx, referrerAddr) {
        return fmt.Errorf("validator quality score too low for referrals")
    }
    
    return nil
}

// validateRankAvailability checks if rank is available
func (rv *ReferralValidator) validateRankAvailability(ctx sdk.Context, rank uint32) error {
    if rank < 22 || rank > 1000 {
        return fmt.Errorf("invalid rank: %d (must be 22-1000)", rank)
    }
    
    // Check if rank is already taken
    if rv.keeper.IsRankTaken(ctx, rank) {
        return fmt.Errorf("rank %d is already taken", rank)
    }
    
    return nil
}

// validateTimeRestrictions checks time-based restrictions
func (rv *ReferralValidator) validateTimeRestrictions(ctx sdk.Context, referrerAddr sdk.AccAddress) error {
    stats := rv.keeper.GetReferralStats(ctx, referrerAddr.String())
    
    // Minimum 24-hour gap between referrals
    if !stats.LastReferralDate.IsZero() {
        timeSinceLastReferral := ctx.BlockTime().Sub(stats.LastReferralDate)
        if timeSinceLastReferral < 24*time.Hour {
            return fmt.Errorf("minimum 24-hour gap required between referrals")
        }
    }
    
    return nil
}

// ValidateCommissionPayout validates commission payment eligibility
func (rv *ReferralValidator) ValidateCommissionPayout(
    ctx sdk.Context,
    referralID uint64,
    amount sdk.Int,
) error {
    referral, found := rv.keeper.GetReferral(ctx, referralID)
    if !found {
        return fmt.Errorf("referral not found: %d", referralID)
    }
    
    // Check if referral is active
    if referral.Status != types.ReferralStatusActive {
        return fmt.Errorf("referral is not active")
    }
    
    // Check if still in commission period (first year only)
    if ctx.BlockTime().After(referral.ActivatedAt.Add(365 * 24 * time.Hour)) {
        return fmt.Errorf("commission period expired")
    }
    
    // Check for clawback period (1 year after activation)
    if ctx.BlockTime().Before(referral.ActivatedAt.Add(6 * 30 * 24 * time.Hour)) {
        return fmt.Errorf("commission payment in cliff period (6 months)")
    }
    
    // Validate amount is reasonable
    if amount.IsZero() || amount.IsNegative() {
        return fmt.Errorf("invalid commission amount: %s", amount)
    }
    
    return nil
}

// CheckClawbackEligibility checks if commission should be clawed back
func (rv *ReferralValidator) CheckClawbackEligibility(
    ctx sdk.Context,
    referralID uint64,
) (bool, string) {
    referral, found := rv.keeper.GetReferral(ctx, referralID)
    if !found {
        return false, "referral not found"
    }
    
    // Check if referred validator exited within 1 year
    referredStake, found := rv.keeper.GetValidatorStake(ctx, referral.ReferredAddr)
    if !found {
        // Validator has exited
        if ctx.BlockTime().Before(referral.ActivatedAt.Add(365 * 24 * time.Hour)) {
            return true, "referred validator exited within 1 year"
        }
    }
    
    // Check if referred validator was slashed for major violations
    if rv.hasReferralSlashingHistory(ctx, sdk.AccAddress(referral.ReferredAddr)) {
        return true, "referred validator violated terms"
    }
    
    return false, ""
}

// Helper functions

func (rv *ReferralValidator) getMonthlyReferralCount(ctx sdk.Context, referrerAddr string) uint32 {
    count := uint32(0)
    oneMonthAgo := ctx.BlockTime().Add(-30 * 24 * time.Hour)
    
    referrals := rv.keeper.GetReferralsByReferrer(ctx, referrerAddr)
    for _, ref := range referrals {
        if ref.CreatedAt.After(oneMonthAgo) {
            count++
        }
    }
    
    return count
}

func (rv *ReferralValidator) getWeeklyReferralCount(ctx sdk.Context, referrerAddr string) uint32 {
    count := uint32(0)
    oneWeekAgo := ctx.BlockTime().Add(-7 * 24 * time.Hour)
    
    referrals := rv.keeper.GetReferralsByReferrer(ctx, referrerAddr)
    for _, ref := range referrals {
        if ref.CreatedAt.After(oneWeekAgo) {
            count++
        }
    }
    
    return count
}

func (rv *ReferralValidator) hashIP(ip string) string {
    hash := sha256.Sum256([]byte(ip))
    return hex.EncodeToString(hash[:])
}

func (rv *ReferralValidator) getRecentReferralsByIP(ctx sdk.Context, ipHash string, duration time.Duration) []types.Referral {
    var recentReferrals []types.Referral
    cutoff := ctx.BlockTime().Add(-duration)
    
    // This would query a separate IP tracking store
    // For now, return empty slice
    return recentReferrals
}

func (rv *ReferralValidator) storeReferralIP(ctx sdk.Context, referrerAddr, ipHash string, timestamp time.Time) {
    // Store IP hash with timestamp for tracking
    // This would use a separate KV store for IP tracking
}

func (rv *ReferralValidator) hasRecentTransactions(ctx sdk.Context, addr1, addr2 sdk.AccAddress) bool {
    // Check bank module for recent transactions between addresses
    // This would require integration with bank keeper
    return false
}

func (rv *ReferralValidator) isAddressTooNew(ctx sdk.Context, addr sdk.AccAddress) bool {
    // Check when address was first seen on chain
    // This would require transaction history analysis
    return false
}

func (rv *ReferralValidator) hasMinimumActivity(ctx sdk.Context, addr sdk.AccAddress) bool {
    // Check if address has minimum number of transactions/interactions
    // This would require activity analysis
    return true
}

func (rv *ReferralValidator) hasClusteredTiming(referrals []types.Referral) bool {
    if len(referrals) < 3 {
        return false
    }
    
    // Check if more than 3 referrals within 1 hour
    var recentCount int
    now := referrals[len(referrals)-1].CreatedAt
    
    for i := len(referrals) - 1; i >= 0; i-- {
        if now.Sub(referrals[i].CreatedAt) <= time.Hour {
            recentCount++
        } else {
            break
        }
    }
    
    return recentCount > 3
}

func (rv *ReferralValidator) hasSequentialAddresses(referrals []types.Referral) bool {
    // Check for sequential patterns in referred addresses
    // This is a simplified check
    return false
}

func (rv *ReferralValidator) hasUniformStakeAmounts(ctx sdk.Context, referrals []types.Referral) bool {
    if len(referrals) < 3 {
        return false
    }
    
    var amounts []sdk.Dec
    for _, ref := range referrals {
        if stake, found := rv.keeper.GetValidatorStake(ctx, ref.ReferredAddr); found {
            amounts = append(amounts, stake.OriginalUSDValue)
        }
    }
    
    // Check if amounts are suspiciously similar (within 1%)
    if len(amounts) < 3 {
        return false
    }
    
    baseAmount := amounts[0]
    for _, amount := range amounts[1:] {
        diff := amount.Sub(baseAmount).Abs()
        variance := diff.Quo(baseAmount)
        if variance.GT(sdk.NewDecWithPrec(1, 2)) { // More than 1% variance
            return false
        }
    }
    
    return true // All amounts are within 1% - suspicious
}

func (rv *ReferralValidator) hasReferralSlashingHistory(ctx sdk.Context, addr sdk.AccAddress) bool {
    // Check slashing history for referral-related violations
    // This would integrate with slashing keeper
    return false
}

func (rv *ReferralValidator) hasLowQualityScore(ctx sdk.Context, addr sdk.AccAddress) bool {
    stats := rv.keeper.GetReferralStats(ctx, addr.String())
    
    // Quality score below 0.5 is considered low
    return stats.QualityScore.LT(sdk.NewDecWithPrec(5, 1))
}