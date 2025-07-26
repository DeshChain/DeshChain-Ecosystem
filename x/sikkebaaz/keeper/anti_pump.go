/*
Copyright 2024 DeshChain Foundation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package keeper

import (
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/sikkebaaz/types"
)

// ValidateAntiPumpTrade validates a trade against anti-pump and dump protection rules
func (k Keeper) ValidateAntiPumpTrade(ctx sdk.Context, tokenAddress, trader string, amount sdk.Int, isSource bool) error {
	// Get wallet limits
	walletLimits, found := k.getWalletLimits(ctx, tokenAddress, trader)
	if !found {
		// Initialize wallet limits for new trader
		launch, found := k.getTokenLaunchByAddress(ctx, tokenAddress)
		if !found {
			return types.ErrTokenAlreadyExists
		}
		
		currentLimit := launch.GetCurrentWalletLimit(ctx.BlockTime())
		maxAmount := launch.TotalSupply.MulRaw(int64(currentLimit)).QuoRaw(10000)
		
		walletLimits = types.WalletLimits{
			TokenAddress:   tokenAddress,
			WalletAddress:  trader,
			MaxAmount:      maxAmount,
			CurrentAmount:  sdk.ZeroInt(),
			LastTxTime:     time.Time{},
			LastTxBlock:    0,
			ViolationCount: 0,
			IsRestricted:   false,
		}
	}

	// Check if wallet is restricted
	if walletLimits.IsRestricted {
		return types.ErrWalletLimitExceeded
	}

	// Check block-based cooldown for bot protection
	blocksSinceLastTx := ctx.BlockHeight() - walletLimits.LastTxBlock
	if blocksSinceLastTx < int64(types.MinBlocksBetweenTx) {
		walletLimits.ViolationCount++
		if walletLimits.ViolationCount >= 3 {
			walletLimits.IsRestricted = true
			k.setWalletLimits(ctx, walletLimits)
			return types.ErrBotDetected
		}
	}

	// For buy transactions (isSource = false), check wallet limits
	if !isSource {
		newAmount := walletLimits.CurrentAmount.Add(amount)
		if newAmount.GT(walletLimits.MaxAmount) {
			return types.ErrWalletLimitExceeded
		}
		walletLimits.CurrentAmount = newAmount
	} else {
		// For sell transactions, subtract from current amount
		walletLimits.CurrentAmount = walletLimits.CurrentAmount.Sub(amount)
		if walletLimits.CurrentAmount.IsNegative() {
			walletLimits.CurrentAmount = sdk.ZeroInt()
		}
	}

	// Check cooldown period for sell transactions
	if isSource {
		launch, found := k.getTokenLaunchByAddress(ctx, tokenAddress)
		if found {
			timeSinceLastTx := ctx.BlockTime().Sub(walletLimits.LastTxTime)
			if timeSinceLastTx.Seconds() < float64(launch.AntiPumpConfig.CooldownPeriod) {
				return types.ErrCooldownPeriodActive
			}
		}
	}

	// Update wallet limits
	walletLimits.LastTxTime = ctx.BlockTime()
	walletLimits.LastTxBlock = ctx.BlockHeight()
	k.setWalletLimits(ctx, walletLimits)

	// Update trading metrics
	k.updateTradingMetrics(ctx, tokenAddress, amount, isSource)

	return nil
}

// ValidatePriceImpact checks if the trade would cause excessive price impact
func (k Keeper) ValidatePriceImpact(ctx sdk.Context, tokenAddress string, tradeAmount sdk.Int) error {
	// Get current liquidity
	metrics, found := k.getTradingMetrics(ctx, tokenAddress)
	if !found {
		return nil // Allow trades for new tokens
	}

	if metrics.Liquidity.IsZero() {
		return types.ErrNoLiquidity
	}

	// Calculate price impact: tradeAmount / liquidity
	priceImpact := tradeAmount.ToDec().Quo(metrics.Liquidity.ToDec())
	
	// Get launch config for max price impact
	launch, found := k.getTokenLaunchByAddress(ctx, tokenAddress)
	if found && priceImpact.GT(launch.AntiPumpConfig.MaxPriceImpact) {
		return types.ErrPriceImpactTooHigh
	}

	return nil
}

// CheckTradingDelay verifies if trading is allowed for a token
func (k Keeper) CheckTradingDelay(ctx sdk.Context, tokenAddress string) error {
	launch, found := k.getTokenLaunchByAddress(ctx, tokenAddress)
	if !found {
		return types.ErrTokenAlreadyExists
	}

	if launch.CompletedAt == nil {
		return types.ErrTradingNotStarted
	}

	// Check if trading delay has passed
	tradingStartTime := launch.CompletedAt.Add(time.Duration(launch.TradingDelay) * time.Second)
	if ctx.BlockTime().Before(tradingStartTime) {
		return types.ErrTradingNotStarted
	}

	return nil
}

// IsLiquidityLocked checks if liquidity is currently locked
func (k Keeper) IsLiquidityLocked(ctx sdk.Context, tokenAddress string) bool {
	lock, found := k.getLiquidityLock(ctx, tokenAddress)
	if !found {
		return false
	}

	return !lock.IsWithdrawn && ctx.BlockTime().Before(lock.UnlockDate)
}

// UpdateWalletLimitsAfter24h updates wallet limits after 24 hours
func (k Keeper) UpdateWalletLimitsAfter24h(ctx sdk.Context, tokenAddress string) {
	launch, found := k.getTokenLaunchByAddress(ctx, tokenAddress)
	if !found {
		return
	}

	if launch.CompletedAt == nil {
		return
	}

	// Check if 24 hours have passed
	hoursSinceLaunch := ctx.BlockTime().Sub(*launch.CompletedAt).Hours()
	if hoursSinceLaunch < 24 {
		return
	}

	// Update all wallet limits for this token
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixWalletLimits)
	iterator := sdk.KVStorePrefixIterator(store, []byte(tokenAddress))
	defer iterator.Close()

	newMaxPercent := launch.AntiPumpConfig.MaxWalletPercentAfter
	newMaxAmount := launch.TotalSupply.MulRaw(int64(newMaxPercent)).QuoRaw(10000)

	for ; iterator.Valid(); iterator.Next() {
		var limits types.WalletLimits
		k.cdc.MustUnmarshal(iterator.Value(), &limits)

		if limits.TokenAddress == tokenAddress {
			limits.MaxAmount = newMaxAmount
			k.setWalletLimits(ctx, limits)
		}
	}
}

// DetectSuspiciousActivity detects and handles suspicious trading patterns
func (k Keeper) DetectSuspiciousActivity(ctx sdk.Context, tokenAddress, trader string) error {
	// Get trading metrics
	metrics, found := k.getTradingMetrics(ctx, tokenAddress)
	if !found {
		return nil
	}

	// Check for unusual trading volume
	avgDailyVolume := metrics.TotalVolume.QuoRaw(int64(ctx.BlockTime().Sub(time.Time{}).Hours() / 24))
	if metrics.DailyVolume.GT(avgDailyVolume.MulRaw(5)) { // 5x average
		k.Logger(ctx).Warn("Suspicious trading volume detected", 
			"token", tokenAddress, 
			"daily_volume", metrics.DailyVolume,
			"average", avgDailyVolume,
		)
	}

	// Check for rapid-fire transactions
	walletLimits, found := k.getWalletLimits(ctx, tokenAddress, trader)
	if found && walletLimits.ViolationCount >= 5 {
		// Temporary restriction
		walletLimits.IsRestricted = true
		k.setWalletLimits(ctx, walletLimits)
		
		k.Logger(ctx).Warn("Trader temporarily restricted for suspicious activity",
			"token", tokenAddress,
			"trader", trader,
			"violations", walletLimits.ViolationCount,
		)
		
		return types.ErrBotDetected
	}

	return nil
}

// updateTradingMetrics updates trading statistics
func (k Keeper) updateTradingMetrics(ctx sdk.Context, tokenAddress string, amount sdk.Int, isSource bool) {
	metrics, found := k.getTradingMetrics(ctx, tokenAddress)
	if !found {
		// Initialize new metrics
		metrics = types.TradingMetrics{
			TokenAddress:   tokenAddress,
			TotalVolume:    sdk.ZeroInt(),
			DailyVolume:    sdk.ZeroInt(),
			TotalTrades:    0,
			DailyTrades:    0,
			UniqueTraders:  0,
			CurrentPrice:   sdk.ZeroDec(),
			PriceChange24h: sdk.ZeroDec(),
			MarketCap:      sdk.ZeroInt(),
			Liquidity:      sdk.ZeroInt(),
			LastUpdated:    ctx.BlockTime(),
		}
	}

	// Update volumes
	metrics.TotalVolume = metrics.TotalVolume.Add(amount)
	metrics.TotalTrades++

	// Reset daily counters if it's a new day
	if ctx.BlockTime().Day() != metrics.LastUpdated.Day() {
		metrics.DailyVolume = sdk.ZeroInt()
		metrics.DailyTrades = 0
	}

	metrics.DailyVolume = metrics.DailyVolume.Add(amount)
	metrics.DailyTrades++
	metrics.LastUpdated = ctx.BlockTime()

	k.setTradingMetrics(ctx, metrics)

	// Emit trading event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"token_trade",
			sdk.NewAttribute("token_address", tokenAddress),
			sdk.NewAttribute("amount", amount.String()),
			sdk.NewAttribute("is_source", string(rune(map[bool]int{true: 1, false: 0}[isSource]))),
		),
	)
}

// Helper functions

func (k Keeper) getWalletLimits(ctx sdk.Context, tokenAddress, wallet string) (types.WalletLimits, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetWalletLimitsKey(tokenAddress, wallet))
	if bz == nil {
		return types.WalletLimits{}, false
	}

	var limits types.WalletLimits
	k.cdc.MustUnmarshal(bz, &limits)
	return limits, true
}

func (k Keeper) getLiquidityLock(ctx sdk.Context, tokenAddress string) (types.LiquidityLock, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetLiquidityLockKey(tokenAddress))
	if bz == nil {
		return types.LiquidityLock{}, false
	}

	var lock types.LiquidityLock
	k.cdc.MustUnmarshal(bz, &lock)
	return lock, true
}

func (k Keeper) getTradingMetrics(ctx sdk.Context, tokenAddress string) (types.TradingMetrics, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixTradingMetrics)
	bz := store.Get([]byte(tokenAddress))
	if bz == nil {
		return types.TradingMetrics{}, false
	}

	var metrics types.TradingMetrics
	k.cdc.MustUnmarshal(bz, &metrics)
	return metrics, true
}

func (k Keeper) getTokenLaunchByAddress(ctx sdk.Context, tokenAddress string) (types.TokenLaunch, bool) {
	// This is a simplified implementation - in practice, you'd need an index
	// from token address to launch ID
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixTokenLaunch)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var launch types.TokenLaunch
		k.cdc.MustUnmarshal(iterator.Value(), &launch)
		if launch.TokenSymbol == tokenAddress {
			return launch, true
		}
	}

	return types.TokenLaunch{}, false
}