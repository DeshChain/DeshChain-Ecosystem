package keeper

import (
    "fmt"
    "time"
    
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/deshchain/namo/x/validator/types"
)

// StakingManager handles all validator staking operations
type StakingManager struct {
    keeper Keeper
}

// NewStakingManager creates a new staking manager
func NewStakingManager(k Keeper) *StakingManager {
    return &StakingManager{keeper: k}
}

// OnboardValidator handles the initial staking for a new validator
func (sm *StakingManager) OnboardValidator(
    ctx sdk.Context,
    validatorAddr sdk.AccAddress,
    rank uint32,
    namoPrice sdk.Dec,
) error {
    // Get stake requirement for this rank
    requirements := types.GetValidatorStakeRequirements()
    var requirement *types.ValidatorStakeRequirement
    
    for _, req := range requirements {
        if req.ValidatorRank == rank {
            requirement = &req
            break
        }
    }
    
    if requirement == nil {
        return fmt.Errorf("no stake requirement found for rank %d", rank)
    }
    
    // Calculate NAMO tokens needed
    namoRequired := types.CalculateNAMORequired(requirement.RequiredUSD, namoPrice)
    
    // Get tier information
    tier, found := types.GetTierForRank(rank)
    if !found {
        return fmt.Errorf("no tier found for rank %d", rank)
    }
    
    // Calculate allocations
    performanceBond := namoRequired.ToDec().Mul(tier.PerformanceBondPct).TruncateInt()
    insuranceContribution := namoRequired.ToDec().MulInt64(2).QuoInt64(100).TruncateInt() // 2%
    vestableAmount := namoRequired.Sub(performanceBond)
    
    // Check validator has sufficient balance
    balance := sm.keeper.bankKeeper.GetBalance(ctx, validatorAddr, "namo")
    if balance.Amount.LT(namoRequired) {
        return fmt.Errorf("insufficient balance: have %s, need %s NAMO", 
            balance.Amount, namoRequired)
    }
    
    // Lock the tokens
    if err := sm.keeper.bankKeeper.SendCoinsFromAccountToModule(
        ctx,
        validatorAddr,
        types.ModuleName,
        sdk.NewCoins(sdk.NewCoin("namo", namoRequired)),
    ); err != nil {
        return fmt.Errorf("failed to lock tokens: %w", err)
    }
    
    // Create stake record
    stake := types.ValidatorStake{
        ValidatorAddr:      validatorAddr.String(),
        OriginalUSDValue:   requirement.RequiredUSD,
        NAMOTokensStaked:   namoRequired,
        NAMOPriceAtStake:   namoPrice,
        StakeTimestamp:     ctx.BlockTime(),
        LockEndTime:        ctx.BlockTime().Add(time.Duration(tier.LockPeriodMonths) * 30 * 24 * time.Hour),
        VestingEndTime:     ctx.BlockTime().Add(time.Duration(tier.LockPeriodMonths + tier.VestingMonths) * 30 * 24 * time.Hour),
        PerformanceBond:    performanceBond,
        VestableAmount:     vestableAmount,
        UnlockedAmount:     sdk.ZeroInt(),
        InsuranceContribution: insuranceContribution,
        Tier:               tier.TierID,
        SlashingHistory:    []types.SlashingEvent{},
        LastWithdrawal:     ctx.BlockTime(),
    }
    
    // Store the stake
    sm.keeper.SetValidatorStake(ctx, stake)
    
    // Transfer insurance contribution to insurance pool
    if err := sm.transferToInsurancePool(ctx, insuranceContribution); err != nil {
        return fmt.Errorf("failed to fund insurance pool: %w", err)
    }
    
    // Mint NFT for genesis validators
    if rank <= 21 {
        if _, err := sm.keeper.MintGenesisNFT(ctx, validatorAddr.String(), rank); err != nil {
            return fmt.Errorf("failed to mint genesis NFT: %w", err)
        }
    }
    
    // Emit event
    ctx.EventManager().EmitEvent(
        sdk.NewEvent(
            types.EventTypeValidatorStaked,
            sdk.NewAttribute(types.AttributeKeyValidator, validatorAddr.String()),
            sdk.NewAttribute(types.AttributeKeyRank, fmt.Sprintf("%d", rank)),
            sdk.NewAttribute(types.AttributeKeyStakeAmount, namoRequired.String()),
            sdk.NewAttribute(types.AttributeKeyUSDValue, requirement.RequiredUSD.String()),
            sdk.NewAttribute(types.AttributeKeyNAMOPrice, namoPrice.String()),
            sdk.NewAttribute(types.AttributeKeyTier, fmt.Sprintf("%d", tier.TierID)),
        ),
    )
    
    sm.keeper.Logger(ctx).Info("Validator onboarded",
        "validator", validatorAddr,
        "rank", rank,
        "usd_value", requirement.RequiredUSD,
        "namo_tokens", namoRequired,
        "namo_price", namoPrice,
        "tier", tier.TierID)
    
    return nil
}

// ProcessVesting handles the gradual unlocking of tokens
func (sm *StakingManager) ProcessVesting(ctx sdk.Context, validatorAddr sdk.AccAddress) error {
    stake, found := sm.keeper.GetValidatorStake(ctx, validatorAddr.String())
    if !found {
        return fmt.Errorf("no stake found for validator %s", validatorAddr)
    }
    
    currentTime := ctx.BlockTime()
    
    // Check if lock period has passed
    if currentTime.Before(stake.LockEndTime) {
        return fmt.Errorf("tokens still locked until %s", stake.LockEndTime)
    }
    
    // Check if already fully vested
    if stake.UnlockedAmount.GTE(stake.VestableAmount) {
        return fmt.Errorf("all vestable tokens already unlocked")
    }
    
    // Calculate vesting progress
    vestingDuration := stake.VestingEndTime.Sub(stake.LockEndTime)
    timeSinceLock := currentTime.Sub(stake.LockEndTime)
    
    var unlockedAmount sdk.Int
    
    if currentTime.After(stake.VestingEndTime) {
        // Fully vested
        unlockedAmount = stake.VestableAmount
    } else {
        // Calculate proportional unlock
        progress := sdk.NewDec(timeSinceLock.Nanoseconds()).Quo(
            sdk.NewDec(vestingDuration.Nanoseconds()))
        unlockedAmount = stake.VestableAmount.ToDec().Mul(progress).TruncateInt()
    }
    
    // Calculate newly unlocked amount
    newlyUnlocked := unlockedAmount.Sub(stake.UnlockedAmount)
    
    if newlyUnlocked.IsZero() {
        return fmt.Errorf("no new tokens to unlock")
    }
    
    // Update stake record
    stake.UnlockedAmount = unlockedAmount
    sm.keeper.SetValidatorStake(ctx, stake)
    
    // Emit event
    ctx.EventManager().EmitEvent(
        sdk.NewEvent(
            types.EventTypeTokensVested,
            sdk.NewAttribute(types.AttributeKeyValidator, validatorAddr.String()),
            sdk.NewAttribute(types.AttributeKeyUnlockedAmount, newlyUnlocked.String()),
            sdk.NewAttribute(types.AttributeKeyTotalUnlocked, unlockedAmount.String()),
        ),
    )
    
    return nil
}

// WithdrawUnlocked allows validators to withdraw unlocked tokens
func (sm *StakingManager) WithdrawUnlocked(
    ctx sdk.Context,
    validatorAddr sdk.AccAddress,
    amount sdk.Int,
) error {
    stake, found := sm.keeper.GetValidatorStake(ctx, validatorAddr.String())
    if !found {
        return fmt.Errorf("no stake found for validator %s", validatorAddr)
    }
    
    // Check available balance
    availableForWithdrawal := stake.UnlockedAmount
    if amount.GT(availableForWithdrawal) {
        return fmt.Errorf("insufficient unlocked balance: have %s, requested %s",
            availableForWithdrawal, amount)
    }
    
    // Apply daily sell limit
    tier, _ := types.GetTierForRank(stake.Tier)
    dailyLimit := stake.VestableAmount.ToDec().Mul(tier.DailySellLimitPct).TruncateInt()
    
    // Check if within daily limit
    if ctx.BlockTime().Sub(stake.LastWithdrawal).Hours() < 24 {
        if amount.GT(dailyLimit) {
            return fmt.Errorf("exceeds daily sell limit of %s NAMO", dailyLimit)
        }
    }
    
    // Send tokens back to validator
    if err := sm.keeper.bankKeeper.SendCoinsFromModuleToAccount(
        ctx,
        types.ModuleName,
        validatorAddr,
        sdk.NewCoins(sdk.NewCoin("namo", amount)),
    ); err != nil {
        return fmt.Errorf("failed to withdraw tokens: %w", err)
    }
    
    // Update stake record
    stake.UnlockedAmount = stake.UnlockedAmount.Sub(amount)
    stake.NAMOTokensStaked = stake.NAMOTokensStaked.Sub(amount)
    stake.LastWithdrawal = ctx.BlockTime()
    sm.keeper.SetValidatorStake(ctx, stake)
    
    // Emit event
    ctx.EventManager().EmitEvent(
        sdk.NewEvent(
            types.EventTypeTokensWithdrawn,
            sdk.NewAttribute(types.AttributeKeyValidator, validatorAddr.String()),
            sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
            sdk.NewAttribute(types.AttributeKeyRemainingStake, stake.NAMOTokensStaked.String()),
        ),
    )
    
    return nil
}

// ReleasePerformanceBond releases the performance bond after 3 years
func (sm *StakingManager) ReleasePerformanceBond(
    ctx sdk.Context,
    validatorAddr sdk.AccAddress,
) error {
    stake, found := sm.keeper.GetValidatorStake(ctx, validatorAddr.String())
    if !found {
        return fmt.Errorf("no stake found for validator %s", validatorAddr)
    }
    
    // Check if 3 years have passed
    threeYears := stake.StakeTimestamp.Add(3 * 365 * 24 * time.Hour)
    if ctx.BlockTime().Before(threeYears) {
        return fmt.Errorf("performance bond locked until %s", threeYears)
    }
    
    // Check for slashing history
    if len(stake.SlashingHistory) > 0 {
        return fmt.Errorf("cannot release performance bond due to slashing history")
    }
    
    // Release the bond
    if err := sm.keeper.bankKeeper.SendCoinsFromModuleToAccount(
        ctx,
        types.ModuleName,
        validatorAddr,
        sdk.NewCoins(sdk.NewCoin("namo", stake.PerformanceBond)),
    ); err != nil {
        return fmt.Errorf("failed to release performance bond: %w", err)
    }
    
    // Update stake record
    stake.PerformanceBond = sdk.ZeroInt()
    sm.keeper.SetValidatorStake(ctx, stake)
    
    // Emit event
    ctx.EventManager().EmitEvent(
        sdk.NewEvent(
            types.EventTypePerformanceBondReleased,
            sdk.NewAttribute(types.AttributeKeyValidator, validatorAddr.String()),
            sdk.NewAttribute(types.AttributeKeyAmount, stake.PerformanceBond.String()),
        ),
    )
    
    return nil
}

// transferToInsurancePool transfers tokens to the insurance pool
func (sm *StakingManager) transferToInsurancePool(ctx sdk.Context, amount sdk.Int) error {
    // Implementation depends on insurance pool module
    // For now, we'll store in a separate module account
    insuranceAddr := sm.keeper.GetInsurancePoolAddress(ctx)
    
    return sm.keeper.bankKeeper.SendCoinsFromModuleToAccount(
        ctx,
        types.ModuleName,
        insuranceAddr,
        sdk.NewCoins(sdk.NewCoin("namo", amount)),
    )
}