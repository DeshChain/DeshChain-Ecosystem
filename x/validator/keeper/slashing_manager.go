package keeper

import (
    "fmt"
    "time"
    
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/DeshChain/DeshChain-Ecosystem/x/validator/types"
)

// SlashingReason defines the reason for slashing
type SlashingReason string

const (
    SlashingReasonDowntime          SlashingReason = "downtime"
    SlashingReasonDoubleSign        SlashingReason = "double_sign"
    SlashingReasonMissedBlocks      SlashingReason = "missed_blocks"
    SlashingReasonDumpAttempt       SlashingReason = "dump_attempt"
    SlashingReasonMarketManipulation SlashingReason = "market_manipulation"
    SlashingReasonCollusion         SlashingReason = "collusion"
    SlashingReasonGovernanceAbuse   SlashingReason = "governance_abuse"
)

// GetBaseSlashingRate returns the base slashing rate for a given reason
func GetBaseSlashingRate(reason SlashingReason) sdk.Dec {
    switch reason {
    // Technical violations
    case SlashingReasonDowntime:
        return sdk.NewDecWithPrec(1, 3) // 0.1%
    case SlashingReasonDoubleSign:
        return sdk.NewDecWithPrec(5, 2) // 5%
    case SlashingReasonMissedBlocks:
        return sdk.NewDecWithPrec(1, 4) // 0.01%
    
    // Economic violations
    case SlashingReasonDumpAttempt:
        return sdk.NewDecWithPrec(25, 2) // 25%
    case SlashingReasonMarketManipulation:
        return sdk.NewDecWithPrec(15, 2) // 15%
    case SlashingReasonCollusion:
        return sdk.NewDecWithPrec(30, 2) // 30%
    
    // Governance violations
    case SlashingReasonGovernanceAbuse:
        return sdk.NewDecWithPrec(10, 2) // 10%
    
    default:
        return sdk.ZeroDec()
    }
}

// SlashValidator applies slashing to a validator's stake
func (k Keeper) SlashValidator(
    ctx sdk.Context,
    validatorAddr sdk.AccAddress,
    reason SlashingReason,
    customRate *sdk.Dec, // Optional custom rate
) error {
    stake, found := k.GetValidatorStake(ctx, validatorAddr.String())
    if !found {
        return fmt.Errorf("no stake found for validator %s", validatorAddr)
    }
    
    // Get base slashing rate
    baseRate := GetBaseSlashingRate(reason)
    if customRate != nil {
        baseRate = *customRate
    }
    
    // Apply tier multiplier
    tier, found := types.GetTierForRank(stake.Tier)
    if !found {
        return fmt.Errorf("invalid tier %d", stake.Tier)
    }
    
    finalRate := baseRate.Mul(tier.SlashingMultiplier)
    
    // Calculate slash amount from total stake
    slashAmount := stake.NAMOTokensStaked.ToDec().Mul(finalRate).TruncateInt()
    
    // Ensure we don't slash more than available
    if slashAmount.GT(stake.NAMOTokensStaked) {
        slashAmount = stake.NAMOTokensStaked
    }
    
    // Apply the slash
    stake.NAMOTokensStaked = stake.NAMOTokensStaked.Sub(slashAmount)
    
    // Slash proportionally from vestable amount and performance bond
    vestableRatio := stake.VestableAmount.ToDec().Quo(stake.NAMOTokensStaked.Add(slashAmount).ToDec())
    vestableSlash := slashAmount.ToDec().Mul(vestableRatio).TruncateInt()
    bondSlash := slashAmount.Sub(vestableSlash)
    
    stake.VestableAmount = stake.VestableAmount.Sub(vestableSlash)
    stake.PerformanceBond = stake.PerformanceBond.Sub(bondSlash)
    
    // Record slashing event
    slashingEvent := types.SlashingEvent{
        Timestamp:     ctx.BlockTime(),
        Reason:        string(reason),
        SlashedAmount: slashAmount,
        SlashingRate:  finalRate,
    }
    stake.SlashingHistory = append(stake.SlashingHistory, slashingEvent)
    
    // Update stake
    k.SetValidatorStake(ctx, stake)
    
    // Send slashed tokens to insurance pool (50%) and burn (50%)
    halfSlash := slashAmount.QuoRaw(2)
    
    // To insurance pool
    insuranceAddr := k.GetInsurancePoolAddress(ctx)
    if err := k.bankKeeper.SendCoinsFromModuleToAccount(
        ctx,
        types.ModuleName,
        insuranceAddr,
        sdk.NewCoins(sdk.NewCoin("namo", halfSlash)),
    ); err != nil {
        return fmt.Errorf("failed to send to insurance pool: %w", err)
    }
    
    // Burn the other half
    if err := k.bankKeeper.BurnCoins(
        ctx,
        types.ModuleName,
        sdk.NewCoins(sdk.NewCoin("namo", slashAmount.Sub(halfSlash))),
    ); err != nil {
        return fmt.Errorf("failed to burn slashed tokens: %w", err)
    }
    
    // Emit event
    ctx.EventManager().EmitEvent(
        sdk.NewEvent(
            types.EventTypeValidatorSlashed,
            sdk.NewAttribute(types.AttributeKeyValidator, validatorAddr.String()),
            sdk.NewAttribute(types.AttributeKeySlashingReason, string(reason)),
            sdk.NewAttribute(types.AttributeKeySlashingRate, finalRate.String()),
            sdk.NewAttribute(types.AttributeKeyAmount, slashAmount.String()),
        ),
    )
    
    k.Logger(ctx).Info("Validator slashed",
        "validator", validatorAddr,
        "reason", reason,
        "rate", finalRate,
        "amount", slashAmount)
    
    return nil
}

// DetectDumpAttempt checks for potential dump attempts
func (k Keeper) DetectDumpAttempt(
    ctx sdk.Context,
    validatorAddr sdk.AccAddress,
    withdrawAmount sdk.Int,
) (bool, error) {
    stake, found := k.GetValidatorStake(ctx, validatorAddr.String())
    if !found {
        return false, nil
    }
    
    // Get tier for daily limit
    tier, found := types.GetTierForRank(stake.Tier)
    if !found {
        return false, fmt.Errorf("invalid tier %d", stake.Tier)
    }
    
    // Check if withdrawal exceeds 5x daily limit
    dailyLimit := stake.VestableAmount.ToDec().Mul(tier.DailySellLimitPct).TruncateInt()
    suspiciousThreshold := dailyLimit.MulRaw(5)
    
    if withdrawAmount.GT(suspiciousThreshold) {
        // Check recent withdrawal pattern
        timeSinceLastWithdrawal := ctx.BlockTime().Sub(stake.LastWithdrawal)
        
        // If large withdrawal within 24 hours, flag as dump attempt
        if timeSinceLastWithdrawal.Hours() < 24 {
            return true, nil
        }
    }
    
    // Check if trying to withdraw more than 50% of unlocked amount
    if withdrawAmount.GT(stake.UnlockedAmount.QuoRaw(2)) {
        return true, nil
    }
    
    return false, nil
}

// HandleDowntime processes validator downtime
func (k Keeper) HandleDowntime(
    ctx sdk.Context,
    validatorAddr sdk.AccAddress,
    downtimeHours int64,
) error {
    if downtimeHours < 24 {
        return nil // No slashing for less than 24 hours
    }
    
    // Calculate slashing rate: 0.1% per day after 24 hours
    daysDown := (downtimeHours - 24) / 24
    if daysDown == 0 {
        daysDown = 1
    }
    
    customRate := sdk.NewDecWithPrec(1, 3).MulInt64(daysDown) // 0.1% * days
    
    return k.SlashValidator(ctx, validatorAddr, SlashingReasonDowntime, &customRate)
}

// HandleMissedBlocks processes missed block violations
func (k Keeper) HandleMissedBlocks(
    ctx sdk.Context,
    validatorAddr sdk.AccAddress,
    missedBlocks int64,
) error {
    // Slash 0.01% per 100 missed blocks
    if missedBlocks < 100 {
        return nil
    }
    
    multiplier := missedBlocks / 100
    customRate := sdk.NewDecWithPrec(1, 4).MulInt64(multiplier)
    
    return k.SlashValidator(ctx, validatorAddr, SlashingReasonMissedBlocks, &customRate)
}