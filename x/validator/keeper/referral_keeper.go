package keeper

import (
    "fmt"
    "time"
    "strings"
    
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/DeshChain/DeshChain-Ecosystem/x/validator/types"
)

// ReferralKeeper handles all referral operations
type ReferralKeeper struct {
    keeper Keeper
    sikkebaazIntegration *SikkebaazIntegration
    validator *ReferralValidator
}

// NewReferralKeeper creates a new referral keeper
func NewReferralKeeper(k Keeper) *ReferralKeeper {
    return &ReferralKeeper{
        keeper: k,
        sikkebaazIntegration: NewSikkebaazIntegration(k),
        validator: NewReferralValidator(k),
    }
}

// CreateReferral creates a new validator referral with comprehensive validation
func (rk *ReferralKeeper) CreateReferral(
    ctx sdk.Context,
    referrerAddr sdk.AccAddress,
    referredAddr sdk.AccAddress,
    referredRank uint32,
) (uint64, error) {
    return rk.CreateReferralWithIP(ctx, referrerAddr, referredAddr, referredRank, "")
}

// CreateReferralWithIP creates a new validator referral with IP tracking
func (rk *ReferralKeeper) CreateReferralWithIP(
    ctx sdk.Context,
    referrerAddr sdk.AccAddress,
    referredAddr sdk.AccAddress,
    referredRank uint32,
    clientIP string,
) (uint64, error) {
    // Comprehensive validation using ReferralValidator
    if err := rk.validator.ValidateReferralEligibility(
        ctx, referrerAddr, referredAddr, referredRank, clientIP,
    ); err != nil {
        return 0, fmt.Errorf("referral validation failed: %w", err)
    }
    
    // Check if referred validator already exists
    if _, found := rk.keeper.GetValidatorStake(ctx, referredAddr.String()); found {
        return 0, fmt.Errorf("referred address is already a validator")
    }
    
    // Get referral stats and commission tier
    stats := rk.GetReferralStats(ctx, referrerAddr.String())
    tier := types.GetTierForReferralCount(stats.TotalReferrals)
    
    // Create referral
    referralID := rk.keeper.GetNextReferralID(ctx)
    referral := types.Referral{
        ReferralID:      referralID,
        ReferrerAddr:    referrerAddr.String(),
        ReferredAddr:    referredAddr.String(),
        ReferredRank:    referredRank,
        Status:          types.ReferralStatusPending,
        CreatedAt:       ctx.BlockTime(),
        CommissionRate:  tier.CommissionRate,
        TotalCommission: sdk.ZeroInt(),
        PaidCommission:  sdk.ZeroInt(),
        LiquidityLocked: sdk.ZeroInt(),
        ClawbackPeriod:  ctx.BlockTime().Add(365 * 24 * time.Hour), // 1 year
    }
    
    rk.keeper.SetReferral(ctx, referral)
    rk.keeper.SetNextReferralID(ctx, referralID+1)
    
    // Update referrer stats
    stats.TotalReferrals++
    stats.LastReferralDate = ctx.BlockTime()
    rk.keeper.SetReferralStats(ctx, stats)
    
    // Emit event
    ctx.EventManager().EmitEvent(
        sdk.NewEvent(
            "referral_created",
            sdk.NewAttribute("referral_id", fmt.Sprintf("%d", referralID)),
            sdk.NewAttribute("referrer", referrerAddr.String()),
            sdk.NewAttribute("referred", referredAddr.String()),
            sdk.NewAttribute("commission_rate", tier.CommissionRate.String()),
        ),
    )
    
    return referralID, nil
}

// ActivateReferral activates a referral when the referred validator stakes
func (rk *ReferralKeeper) ActivateReferral(
    ctx sdk.Context,
    referredAddr sdk.AccAddress,
) error {
    // Find pending referral for this address
    referral, found := rk.getPendingReferralByReferred(ctx, referredAddr.String())
    if !found {
        return nil // No referral, normal validator onboarding
    }
    
    // Verify the validator has actually staked
    stake, found := rk.keeper.GetValidatorStake(ctx, referredAddr.String())
    if !found {
        return fmt.Errorf("referred validator has not staked")
    }
    
    // Activate the referral
    referral.Status = types.ReferralStatusActive
    referral.ActivatedAt = ctx.BlockTime()
    referral.ClawbackPeriod = ctx.BlockTime().Add(365 * 24 * time.Hour)
    
    rk.keeper.SetReferral(ctx, referral)
    
    // Update referrer stats
    stats := rk.GetReferralStats(ctx, referral.ReferrerAddr)
    stats.ActiveReferrals++
    stats.CurrentTier = uint32(types.GetTierForReferralCount(stats.TotalReferrals).TierID)
    
    // Calculate quality score based on stake amount
    qualityScore := rk.calculateQualityScore(stake.OriginalUSDValue)
    stats.QualityScore = stats.QualityScore.Add(qualityScore).Quo(sdk.NewDec(int64(stats.ActiveReferrals)))
    
    rk.keeper.SetReferralStats(ctx, stats)
    
    // Check if eligible for token launch
    rk.checkTokenLaunchEligibility(ctx, referral.ReferrerAddr)
    
    // Emit event
    ctx.EventManager().EmitEvent(
        sdk.NewEvent(
            "referral_activated",
            sdk.NewAttribute("referral_id", fmt.Sprintf("%d", referral.ReferralID)),
            sdk.NewAttribute("referrer", referral.ReferrerAddr),
            sdk.NewAttribute("referred", referral.ReferredAddr),
        ),
    )
    
    return nil
}

// ProcessReferralCommission calculates and pays referral commission
func (rk *ReferralKeeper) ProcessReferralCommission(
    ctx sdk.Context,
    referredAddr sdk.AccAddress,
    revenueAmount sdk.Int,
) error {
    // Find active referral
    referral, found := rk.getActiveReferralByReferred(ctx, referredAddr.String())
    if !found {
        return nil // No referral commission
    }
    
    // Calculate commission
    commission := revenueAmount.ToDec().Mul(referral.CommissionRate).TruncateInt()
    
    // Validate commission payout with anti-gaming checks
    if err := rk.validator.ValidateCommissionPayout(ctx, referral.ReferralID, commission); err != nil {
        return fmt.Errorf("commission validation failed: %w", err)
    }
    
    // Check for clawback eligibility
    if shouldClawback, reason := rk.validator.CheckClawbackEligibility(ctx, referral.ReferralID); shouldClawback {
        // Mark for clawback instead of paying
        return rk.processClawback(ctx, referral, reason)
    }
    
    // Check if still in first year
    if ctx.BlockTime().After(referral.ActivatedAt.Add(365 * 24 * time.Hour)) {
        referral.Status = types.ReferralStatusCompleted
        rk.keeper.SetReferral(ctx, referral)
        return nil
    }
    
    // Check if validator has launched token
    stats := rk.GetReferralStats(ctx, referral.ReferrerAddr)
    
    if stats.TokenLaunched {
        // Convert commission to liquidity
        err := rk.convertCommissionToLiquidity(ctx, referral.ReferrerAddr, commission, stats.TokenID)
        if err != nil {
            return fmt.Errorf("failed to convert commission to liquidity: %w", err)
        }
        
        referral.LiquidityLocked = referral.LiquidityLocked.Add(commission)
    }
    
    // Update referral
    referral.TotalCommission = referral.TotalCommission.Add(commission)
    referral.PaidCommission = referral.PaidCommission.Add(commission)
    rk.keeper.SetReferral(ctx, referral)
    
    // Update stats
    stats.TotalCommission = stats.TotalCommission.Add(commission)
    stats.LiquidityLocked = stats.LiquidityLocked.Add(commission)
    rk.keeper.SetReferralStats(ctx, stats)
    
    // Create payout record
    payout := types.CommissionPayout{
        PayoutID:       rk.keeper.GetNextPayoutID(ctx),
        ReferralID:     referral.ReferralID,
        ReferrerAddr:   referral.ReferrerAddr,
        Amount:         commission,
        PayoutTime:     ctx.BlockTime(),
        BlockHeight:    ctx.BlockHeight(),
    }
    
    if stats.TokenLaunched {
        payout.TokenAmount = commission // Will be calculated based on pool ratio
        payout.LiquidityAdded = commission
    }
    
    rk.keeper.SetCommissionPayout(ctx, payout)
    rk.keeper.SetNextPayoutID(ctx, payout.PayoutID+1)
    
    // Emit event
    ctx.EventManager().EmitEvent(
        sdk.NewEvent(
            "referral_commission_paid",
            sdk.NewAttribute("referral_id", fmt.Sprintf("%d", referral.ReferralID)),
            sdk.NewAttribute("amount", commission.String()),
            sdk.NewAttribute("liquidity_locked", fmt.Sprintf("%v", stats.TokenLaunched)),
        ),
    )
    
    return nil
}

// checkTokenLaunchEligibility checks if validator can launch token
func (rk *ReferralKeeper) checkTokenLaunchEligibility(
    ctx sdk.Context,
    validatorAddr string,
) bool {
    stats := rk.GetReferralStats(ctx, validatorAddr)
    if stats.TokenLaunched {
        return false
    }
    
    conditions := types.GetTokenLaunchConditions()
    
    // Check referral count
    if stats.TotalReferrals >= conditions.MinReferrals {
        return rk.launchValidatorToken(ctx, validatorAddr, "referrals")
    }
    
    // Check commission earned
    if stats.TotalCommission.GTE(conditions.MinCommission) {
        return rk.launchValidatorToken(ctx, validatorAddr, "commission")
    }
    
    return false
}

// launchValidatorToken launches token on Sikkebaaz
func (rk *ReferralKeeper) launchValidatorToken(
    ctx sdk.Context,
    validatorAddr string,
    trigger string,
) bool {
    // Get validator's NFT for token details
    var validatorNFT *types.GenesisValidatorNFT
    nfts := rk.keeper.GetAllGenesisNFTs(ctx)
    for _, nft := range nfts {
        if nft.CurrentOwner == validatorAddr {
            validatorNFT = &nft
            break
        }
    }
    
    if validatorNFT == nil {
        return false
    }
    
    // Create token
    tokenID := rk.keeper.GetNextTokenID(ctx)
    totalSupply := sdk.NewInt(1000000000000000) // 1 billion with 6 decimals
    
    token := types.ValidatorToken{
        TokenID:       tokenID,
        ValidatorAddr: validatorAddr,
        TokenName:     fmt.Sprintf("%s Coin", validatorNFT.EnglishName),
        TokenSymbol:   rk.generateTokenSymbol(validatorNFT.EnglishName),
        TotalSupply:   totalSupply,
        Decimals:      6,
        LogoURI:       fmt.Sprintf("/nfts/%d.png", validatorNFT.TokenID),
        
        // Distribution
        ValidatorAllocation:   totalSupply.MulRaw(40).QuoRaw(100),
        LiquidityAllocation:   totalSupply.MulRaw(30).QuoRaw(100),
        AirdropAllocation:     totalSupply.MulRaw(15).QuoRaw(100),
        DevelopmentAllocation: totalSupply.MulRaw(10).QuoRaw(100),
        InitialLiquidity:      totalSupply.MulRaw(5).QuoRaw(100),
        
        // Launch info
        LaunchedAt:    ctx.BlockTime(),
        LaunchTrigger: trigger,
        
        // Anti-dump
        MaxWalletPercent: sdk.NewDecWithPrec(2, 2),   // 2%
        MaxTxPercent:     sdk.NewDecWithPrec(5, 3),   // 0.5%
        SellTaxPercent:   sdk.NewDecWithPrec(5, 2),   // 5%
        BuyTaxPercent:    sdk.NewDecWithPrec(2, 2),   // 2%
        CooldownSeconds:  3600,                        // 1 hour
    }
    
    // Launch on Sikkebaaz (would integrate with Sikkebaaz module)
    err := rk.launchOnSikkebaaz(ctx, token)
    if err != nil {
        return false
    }
    
    rk.keeper.SetValidatorToken(ctx, token)
    rk.keeper.SetNextTokenID(ctx, tokenID+1)
    
    // Update stats
    stats := rk.GetReferralStats(ctx, validatorAddr)
    stats.TokenLaunched = true
    stats.TokenID = tokenID
    rk.keeper.SetReferralStats(ctx, stats)
    
    // Emit event
    ctx.EventManager().EmitEvent(
        sdk.NewEvent(
            "validator_token_launched",
            sdk.NewAttribute("token_id", fmt.Sprintf("%d", tokenID)),
            sdk.NewAttribute("validator", validatorAddr),
            sdk.NewAttribute("token_name", token.TokenName),
            sdk.NewAttribute("token_symbol", token.TokenSymbol),
            sdk.NewAttribute("trigger", trigger),
        ),
    )
    
    return true
}

// Helper functions

func (rk *ReferralKeeper) getMonthlyReferralCount(ctx sdk.Context, referrerAddr string) uint32 {
    count := uint32(0)
    oneMonthAgo := ctx.BlockTime().Add(-30 * 24 * time.Hour)
    
    // Iterate through recent referrals
    referrals := rk.keeper.GetReferralsByReferrer(ctx, referrerAddr)
    for _, ref := range referrals {
        if ref.CreatedAt.After(oneMonthAgo) {
            count++
        }
    }
    
    return count
}

func (rk *ReferralKeeper) getPendingReferralByReferred(ctx sdk.Context, referredAddr string) (types.Referral, bool) {
    referrals := rk.keeper.GetAllReferrals(ctx)
    for _, ref := range referrals {
        if ref.ReferredAddr == referredAddr && ref.Status == types.ReferralStatusPending {
            return ref, true
        }
    }
    return types.Referral{}, false
}

func (rk *ReferralKeeper) getActiveReferralByReferred(ctx sdk.Context, referredAddr string) (types.Referral, bool) {
    referrals := rk.keeper.GetAllReferrals(ctx)
    for _, ref := range referrals {
        if ref.ReferredAddr == referredAddr && ref.Status == types.ReferralStatusActive {
            return ref, true
        }
    }
    return types.Referral{}, false
}

func (rk *ReferralKeeper) calculateQualityScore(stakeUSD sdk.Dec) sdk.Dec {
    // Higher stake = higher quality score
    // $200K = 1.0, $1M = 2.0, $2M = 3.0
    return stakeUSD.QuoInt(sdk.NewInt(200000)).Add(sdk.OneDec())
}

func (rk *ReferralKeeper) generateTokenSymbol(nftName string) string {
    // Generate 3-5 letter symbol from NFT name
    words := strings.Fields(strings.ToUpper(nftName))
    if len(words) == 1 {
        // Take first 4 letters
        if len(words[0]) >= 4 {
            return words[0][:4]
        }
        return words[0]
    }
    
    // Take first letter of each word
    symbol := ""
    for _, word := range words {
        if len(word) > 0 {
            symbol += string(word[0])
        }
    }
    
    if len(symbol) > 5 {
        symbol = symbol[:5]
    }
    
    return symbol
}

func (rk *ReferralKeeper) convertCommissionToLiquidity(
    ctx sdk.Context,
    validatorAddr string,
    commission sdk.Int,
    tokenID uint64,
) error {
    // Platform takes 5% fee
    platformFee := commission.ToDec().Mul(sdk.NewDecWithPrec(5, 2)).TruncateInt()
    netLiquidity := commission.Sub(platformFee)
    
    // Add liquidity to validator token pool on Sikkebaaz
    tokenAmount, err := rk.sikkebaazIntegration.AddLiquidity(ctx, tokenID, netLiquidity)
    if err != nil {
        return fmt.Errorf("failed to add liquidity: %w", err)
    }
    
    // Transfer platform fee to treasury
    treasuryAddr := rk.keeper.accountKeeper.GetModuleAddress(types.ModuleName)
    if err := rk.keeper.bankKeeper.SendCoinsFromModuleToAccount(
        ctx,
        types.ModuleName,
        treasuryAddr,
        sdk.NewCoins(sdk.NewCoin("namo", platformFee)),
    ); err != nil {
        return fmt.Errorf("failed to send platform fee: %w", err)
    }
    
    // Emit event
    ctx.EventManager().EmitEvent(
        sdk.NewEvent(
            "commission_converted_to_liquidity",
            sdk.NewAttribute("validator", validatorAddr),
            sdk.NewAttribute("commission", commission.String()),
            sdk.NewAttribute("net_liquidity", netLiquidity.String()),
            sdk.NewAttribute("token_amount", tokenAmount.String()),
            sdk.NewAttribute("platform_fee", platformFee.String()),
        ),
    )
    
    return nil
}

func (rk *ReferralKeeper) launchOnSikkebaaz(ctx sdk.Context, token types.ValidatorToken) error {
    // Validate token launch parameters
    if err := rk.sikkebaazIntegration.ValidateTokenLaunch(ctx, token); err != nil {
        return fmt.Errorf("token launch validation failed: %w", err)
    }
    
    // Launch on Sikkebaaz platform
    if err := rk.sikkebaazIntegration.LaunchValidatorToken(ctx, token); err != nil {
        return fmt.Errorf("failed to launch token on Sikkebaaz: %w", err)
    }
    
    return nil
}

// GetReferralStats returns referral statistics for a validator
func (rk *ReferralKeeper) GetReferralStats(ctx sdk.Context, validatorAddr string) types.ReferralStats {
    stats, found := rk.keeper.GetReferralStats(ctx, validatorAddr)
    if !found {
        // Initialize new stats
        stats = types.ReferralStats{
            ValidatorAddr:   validatorAddr,
            TotalReferrals:  0,
            ActiveReferrals: 0,
            TotalCommission: sdk.ZeroInt(),
            CurrentTier:     1,
            TokenLaunched:   false,
            LiquidityLocked: sdk.ZeroInt(),
            QualityScore:    sdk.OneDec(),
        }
    }
    return stats
}

// processClawback handles commission clawback for referral violations
func (rk *ReferralKeeper) processClawback(
    ctx sdk.Context,
    referral types.Referral,
    reason string,
) error {
    // Calculate total clawback amount
    clawbackAmount := referral.PaidCommission
    
    if clawbackAmount.IsZero() {
        return nil // Nothing to clawback
    }
    
    // Mark referral as clawed back
    referral.Status = types.ReferralStatusClawedBack
    referral.ClawbackAmount = clawbackAmount
    referral.ClawbackReason = reason
    rk.keeper.SetReferral(ctx, referral)
    
    // Update referrer stats
    stats := rk.GetReferralStats(ctx, referral.ReferrerAddr)
    stats.TotalCommission = stats.TotalCommission.Sub(clawbackAmount)
    stats.LiquidityLocked = stats.LiquidityLocked.Sub(clawbackAmount)
    rk.keeper.SetReferralStats(ctx, stats)
    
    // Emit clawback event
    ctx.EventManager().EmitEvent(
        sdk.NewEvent(
            "referral_commission_clawed_back",
            sdk.NewAttribute("referral_id", fmt.Sprintf("%d", referral.ReferralID)),
            sdk.NewAttribute("amount", clawbackAmount.String()),
            sdk.NewAttribute("reason", reason),
        ),
    )
    
    return nil
}