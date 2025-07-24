package keeper

import (
    "fmt"
    "time"
    
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/deshchain/namo/x/validator/types"
)

// CircuitBreakerState tracks market conditions
type CircuitBreakerState struct {
    LastPrice          sdk.Dec   `json:"last_price"`
    LastUpdateTime     time.Time `json:"last_update_time"`
    DailyHighPrice     sdk.Dec   `json:"daily_high_price"`
    DailyLowPrice      sdk.Dec   `json:"daily_low_price"`
    CircuitBreakerActive bool    `json:"circuit_breaker_active"`
    BreakerEndTime     time.Time `json:"breaker_end_time"`
    ConsecutiveBreakers uint32   `json:"consecutive_breakers"`
}

// CheckPriceMovement monitors price movements and triggers circuit breakers
func (k Keeper) CheckPriceMovement(ctx sdk.Context) error {
    // Get current NAMO price from oracle
    currentPrice, err := k.oracleKeeper.GetNAMOPriceUSD(ctx)
    if err != nil {
        return fmt.Errorf("failed to get NAMO price: %w", err)
    }
    
    // Get circuit breaker state
    state := k.GetCircuitBreakerState(ctx)
    
    // Check if we're still in a circuit breaker period
    if state.CircuitBreakerActive && ctx.BlockTime().Before(state.BreakerEndTime) {
        return fmt.Errorf("circuit breaker active until %s", state.BreakerEndTime)
    }
    
    // Calculate price change percentage
    var priceChange sdk.Dec
    if !state.LastPrice.IsZero() {
        priceChange = currentPrice.Sub(state.LastPrice).Quo(state.LastPrice).Mul(sdk.NewDec(100))
    }
    
    // Check for circuit breaker triggers
    triggered := false
    var breakerDuration time.Duration
    var reason string
    
    priceDropPercent := priceChange.Neg()
    
    switch {
    case priceDropPercent.GTE(sdk.NewDec(20)):
        // 20% drop: Emergency DAO vote required
        triggered = true
        breakerDuration = 24 * time.Hour
        reason = "20% price drop - emergency measures"
        
    case priceDropPercent.GTE(sdk.NewDec(10)):
        // 10% drop: 1-hour pause
        triggered = true
        breakerDuration = 1 * time.Hour
        reason = "10% price drop"
        
    case priceDropPercent.GTE(sdk.NewDec(5)):
        // 5% drop: 15-minute pause
        triggered = true
        breakerDuration = 15 * time.Minute
        reason = "5% price drop"
    }
    
    if triggered {
        state.CircuitBreakerActive = true
        state.BreakerEndTime = ctx.BlockTime().Add(breakerDuration)
        state.ConsecutiveBreakers++
        
        // Apply progressive measures based on consecutive breakers
        if state.ConsecutiveBreakers >= 3 {
            // Extend all unbonding periods by 7 days
            k.ExtendAllUnbondingPeriods(ctx, 7*24*time.Hour)
            
            // Freeze large transfers
            k.SetLargeTransferFreeze(ctx, true)
        }
        
        // Emit circuit breaker event
        ctx.EventManager().EmitEvent(
            sdk.NewEvent(
                "circuit_breaker_triggered",
                sdk.NewAttribute("reason", reason),
                sdk.NewAttribute("duration", breakerDuration.String()),
                sdk.NewAttribute("price_drop", priceDropPercent.String()),
                sdk.NewAttribute("consecutive_count", fmt.Sprintf("%d", state.ConsecutiveBreakers)),
            ),
        )
        
        k.Logger(ctx).Info("Circuit breaker triggered",
            "reason", reason,
            "duration", breakerDuration,
            "price_drop", priceDropPercent)
    } else {
        // Reset consecutive counter if no breaker for 24 hours
        if ctx.BlockTime().Sub(state.BreakerEndTime).Hours() > 24 {
            state.ConsecutiveBreakers = 0
        }
    }
    
    // Update daily high/low
    if currentPrice.GT(state.DailyHighPrice) || state.DailyHighPrice.IsZero() {
        state.DailyHighPrice = currentPrice
    }
    if currentPrice.LT(state.DailyLowPrice) || state.DailyLowPrice.IsZero() {
        state.DailyLowPrice = currentPrice
    }
    
    // Reset daily values at midnight UTC
    if ctx.BlockTime().Day() != state.LastUpdateTime.Day() {
        state.DailyHighPrice = currentPrice
        state.DailyLowPrice = currentPrice
    }
    
    state.LastPrice = currentPrice
    state.LastUpdateTime = ctx.BlockTime()
    
    k.SetCircuitBreakerState(ctx, state)
    
    return nil
}

// EnforceDailySellLimits checks and enforces daily sell limits for validators
func (k Keeper) EnforceDailySellLimits(
    ctx sdk.Context,
    validatorAddr sdk.AccAddress,
    sellAmount sdk.Int,
) error {
    stake, found := k.GetValidatorStake(ctx, validatorAddr.String())
    if !found {
        return fmt.Errorf("no stake found for validator")
    }
    
    // Get tier for limits
    tier, found := types.GetTierForRank(stake.Tier)
    if !found {
        return fmt.Errorf("invalid tier %d", stake.Tier)
    }
    
    // Calculate daily limit
    dailyLimit := stake.VestableAmount.ToDec().Mul(tier.DailySellLimitPct).TruncateInt()
    
    // Check if within 24-hour window
    if ctx.BlockTime().Sub(stake.LastWithdrawal).Hours() < 24 {
        if sellAmount.GT(dailyLimit) {
            return fmt.Errorf("exceeds daily sell limit of %s NAMO (tier %d: %s%%)", 
                dailyLimit, stake.Tier, tier.DailySellLimitPct.Mul(sdk.NewDec(100)))
        }
    }
    
    // Additional checks during circuit breaker
    cbState := k.GetCircuitBreakerState(ctx)
    if cbState.CircuitBreakerActive {
        // Reduce limits by 50% during circuit breaker
        reducedLimit := dailyLimit.QuoRaw(2)
        if sellAmount.GT(reducedLimit) {
            return fmt.Errorf("during circuit breaker, limit reduced to %s NAMO", reducedLimit)
        }
    }
    
    // Check weekly limit (10%, 5%, 2.5% for tiers)
    weeklyLimitPct := tier.DailySellLimitPct.MulInt64(5) // 5x daily = weekly
    weeklyLimit := stake.VestableAmount.ToDec().Mul(weeklyLimitPct).TruncateInt()
    
    // Get weekly withdrawal total
    weeklyTotal := k.GetWeeklyWithdrawalTotal(ctx, validatorAddr.String())
    if weeklyTotal.Add(sellAmount).GT(weeklyLimit) {
        return fmt.Errorf("exceeds weekly sell limit of %s NAMO", weeklyLimit)
    }
    
    // Check monthly limit (25%, 15%, 10% for tiers)
    monthlyLimitPct := tier.DailySellLimitPct.MulInt64(12) // ~12x daily = monthly
    monthlyLimit := stake.VestableAmount.ToDec().Mul(monthlyLimitPct).TruncateInt()
    
    monthlyTotal := k.GetMonthlyWithdrawalTotal(ctx, validatorAddr.String())
    if monthlyTotal.Add(sellAmount).GT(monthlyLimit) {
        return fmt.Errorf("exceeds monthly sell limit of %s NAMO", monthlyLimit)
    }
    
    return nil
}

// Circuit breaker state management

func (k Keeper) GetCircuitBreakerState(ctx sdk.Context) CircuitBreakerState {
    store := ctx.KVStore(k.storeKey)
    bz := store.Get([]byte("circuit_breaker_state"))
    
    var state CircuitBreakerState
    if bz != nil {
        k.cdc.MustUnmarshal(bz, &state)
    }
    
    return state
}

func (k Keeper) SetCircuitBreakerState(ctx sdk.Context, state CircuitBreakerState) {
    store := ctx.KVStore(k.storeKey)
    bz := k.cdc.MustMarshal(&state)
    store.Set([]byte("circuit_breaker_state"), bz)
}

// Helper functions

func (k Keeper) ExtendAllUnbondingPeriods(ctx sdk.Context, extension time.Duration) {
    stakes := k.GetAllValidatorStakes(ctx)
    
    for _, stake := range stakes {
        stake.LockEndTime = stake.LockEndTime.Add(extension)
        stake.VestingEndTime = stake.VestingEndTime.Add(extension)
        k.SetValidatorStake(ctx, stake)
    }
    
    k.Logger(ctx).Info("Extended all unbonding periods",
        "extension", extension,
        "validators_affected", len(stakes))
}

func (k Keeper) SetLargeTransferFreeze(ctx sdk.Context, frozen bool) {
    store := ctx.KVStore(k.storeKey)
    if frozen {
        store.Set([]byte("large_transfer_freeze"), []byte{1})
    } else {
        store.Delete([]byte("large_transfer_freeze"))
    }
}

func (k Keeper) IsLargeTransferFrozen(ctx sdk.Context) bool {
    store := ctx.KVStore(k.storeKey)
    return store.Has([]byte("large_transfer_freeze"))
}

func (k Keeper) GetWeeklyWithdrawalTotal(ctx sdk.Context, validatorAddr string) sdk.Int {
    // Implementation would track withdrawals over past 7 days
    // For now, return zero
    return sdk.ZeroInt()
}

func (k Keeper) GetMonthlyWithdrawalTotal(ctx sdk.Context, validatorAddr string) sdk.Int {
    // Implementation would track withdrawals over past 30 days
    // For now, return zero
    return sdk.ZeroInt()
}