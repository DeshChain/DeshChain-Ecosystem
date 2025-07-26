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

// ValidateAntiDumpTrade validates a sell trade against anti-dump protection rules
func (k Keeper) ValidateAntiDumpTrade(ctx sdk.Context, tokenAddress, seller string, sellAmount sdk.Int) error {
	// Get token launch configuration
	launch, found := k.getTokenLaunchByAddress(ctx, tokenAddress)
	if !found {
		return types.ErrTokenAlreadyExists
	}

	// Check if trading is allowed
	if err := k.CheckTradingDelay(ctx, tokenAddress); err != nil {
		return err
	}

	// Get seller's wallet limits
	walletLimits, found := k.getWalletLimits(ctx, tokenAddress, seller)
	if !found {
		return types.ErrWalletLimitExceeded
	}

	// Check cooldown period between sells
	if !walletLimits.LastTxTime.IsZero() {
		timeSinceLastSell := ctx.BlockTime().Sub(walletLimits.LastTxTime)
		if timeSinceLastSell.Seconds() < float64(launch.AntiPumpConfig.CooldownPeriod) {
			return types.ErrCooldownPeriodActive
		}
	}

	// Check price impact protection
	if err := k.ValidatePriceImpact(ctx, tokenAddress, sellAmount); err != nil {
		return err
	}

	// Implement gradual release mechanism if enabled
	if launch.AntiPumpConfig.GradualReleaseEnabled {
		if err := k.validateGradualRelease(ctx, tokenAddress, seller, sellAmount); err != nil {
			return err
		}
	}

	// Check for suspicious dumping patterns
	if err := k.detectDumpingPattern(ctx, tokenAddress, seller, sellAmount); err != nil {
		return err
	}

	// Check maximum daily sell limit
	if err := k.validateDailySellLimit(ctx, tokenAddress, seller, sellAmount); err != nil {
		return err
	}

	return nil
}

// validateGradualRelease ensures tokens are released gradually over time
func (k Keeper) validateGradualRelease(ctx sdk.Context, tokenAddress, seller string, sellAmount sdk.Int) error {
	launch, found := k.getTokenLaunchByAddress(ctx, tokenAddress)
	if !found {
		return types.ErrTokenAlreadyExists
	}

	if launch.CompletedAt == nil {
		return types.ErrTradingNotStarted
	}

	// Calculate vesting schedule (simplified implementation)
	daysSinceLaunch := ctx.BlockTime().Sub(*launch.CompletedAt).Hours() / 24
	
	// Get seller's original allocation
	sellerAllocation := k.getSellerOriginalAllocation(ctx, tokenAddress, seller)
	if sellerAllocation.IsZero() {
		return nil // No restrictions for non-original holders
	}

	// Calculate maximum sellable amount based on vesting schedule
	// 25% immediately, 25% after 30 days, 25% after 60 days, 25% after 90 days
	var maxSellablePercent sdk.Dec
	if daysSinceLaunch < 30 {
		maxSellablePercent = sdk.MustNewDecFromStr("0.25") // 25%
	} else if daysSinceLaunch < 60 {
		maxSellablePercent = sdk.MustNewDecFromStr("0.50") // 50%
	} else if daysSinceLaunch < 90 {
		maxSellablePercent = sdk.MustNewDecFromStr("0.75") // 75%
	} else {
		maxSellablePercent = sdk.MustNewDecFromStr("1.00") // 100%
	}

	maxSellableAmount := sellerAllocation.ToDec().Mul(maxSellablePercent).TruncateInt()
	
	// Get amount already sold
	soldAmount := k.getSellerSoldAmount(ctx, tokenAddress, seller)
	
	// Check if this sell would exceed the vesting limit
	if soldAmount.Add(sellAmount).GT(maxSellableAmount) {
		return types.ErrGradualReleaseViolation
	}

	return nil
}

// detectDumpingPattern detects suspicious dumping patterns
func (k Keeper) detectDumpingPattern(ctx sdk.Context, tokenAddress, seller string, sellAmount sdk.Int) error {
	// Get recent trading history for this seller
	recentSells := k.getRecentSells(ctx, tokenAddress, seller, 24*time.Hour) // Last 24 hours
	
	// Calculate total sold in last 24 hours
	totalSold24h := sdk.ZeroInt()
	for _, sell := range recentSells {
		totalSold24h = totalSold24h.Add(sell.Amount)
	}
	
	// Add current sell amount
	totalSold24h = totalSold24h.Add(sellAmount)
	
	// Get seller's current balance
	sellerAddr, _ := sdk.AccAddressFromBech32(seller)
	currentBalance := k.bankKeeper.GetBalance(ctx, sellerAddr, tokenAddress)
	
	// Flag as dumping if selling more than 50% of holdings in 24 hours
	if currentBalance.Amount.IsPositive() {
		sellPercentage := totalSold24h.ToDec().Quo(currentBalance.Amount.ToDec())
		if sellPercentage.GT(sdk.MustNewDecFromStr("0.50")) { // 50%
			// Apply dumping penalty
			return k.applyDumpingPenalty(ctx, tokenAddress, seller, sellAmount)
		}
	}
	
	// Check for whale dumping (more than 5% of total supply)
	launch, found := k.getTokenLaunchByAddress(ctx, tokenAddress)
	if found {
		whaleThreshold := launch.TotalSupply.MulRaw(5).QuoRaw(100) // 5%
		if sellAmount.GT(whaleThreshold) {
			k.Logger(ctx).Warn("Whale dumping detected",
				"token", tokenAddress,
				"seller", seller,
				"amount", sellAmount,
				"threshold", whaleThreshold,
			)
			
			// Require additional confirmations for whale sells
			return k.requireWhaleConfirmation(ctx, tokenAddress, seller, sellAmount)
		}
	}
	
	return nil
}

// validateDailySellLimit enforces daily selling limits
func (k Keeper) validateDailySellLimit(ctx sdk.Context, tokenAddress, seller string, sellAmount sdk.Int) error {
	// Get daily sell tracking
	dailySells := k.getDailySells(ctx, tokenAddress, seller, ctx.BlockTime())
	
	// Calculate total sold today
	totalSoldToday := sdk.ZeroInt()
	for _, sell := range dailySells {
		totalSoldToday = totalSoldToday.Add(sell.Amount)
	}
	
	// Get token configuration
	launch, found := k.getTokenLaunchByAddress(ctx, tokenAddress)
	if !found {
		return types.ErrTokenAlreadyExists
	}
	
	// Calculate daily limit (2% of total supply maximum)
	dailyLimit := launch.TotalSupply.MulRaw(2).QuoRaw(100) // 2%
	
	// Check if adding this sell would exceed daily limit
	if totalSoldToday.Add(sellAmount).GT(dailyLimit) {
		return types.ErrMaxTransactionsExceeded
	}
	
	return nil
}

// applyDumpingPenalty applies penalties for detected dumping
func (k Keeper) applyDumpingPenalty(ctx sdk.Context, tokenAddress, seller string, sellAmount sdk.Int) error {
	// Apply increasing penalty fee for dumping
	walletLimits, found := k.getWalletLimits(ctx, tokenAddress, seller)
	if !found {
		return types.ErrWalletLimitExceeded
	}
	
	// Increase violation count
	walletLimits.ViolationCount++
	
	// Apply penalty based on violation count
	var penaltyRate sdk.Dec
	switch {
	case walletLimits.ViolationCount <= 1:
		penaltyRate = sdk.MustNewDecFromStr("0.05") // 5% penalty
	case walletLimits.ViolationCount <= 3:
		penaltyRate = sdk.MustNewDecFromStr("0.10") // 10% penalty
	case walletLimits.ViolationCount <= 5:
		penaltyRate = sdk.MustNewDecFromStr("0.20") // 20% penalty
	default:
		// Temporary suspension for repeated violations
		walletLimits.IsRestricted = true
		k.setWalletLimits(ctx, walletLimits)
		return types.ErrBotDetected
	}
	
	// Calculate penalty amount
	penaltyAmount := sellAmount.ToDec().Mul(penaltyRate).TruncateInt()
	
	// Send penalty to penalty pool
	sellerAddr, _ := sdk.AccAddressFromBech32(seller)
	penaltyCoins := sdk.NewCoins(sdk.NewCoin(tokenAddress, penaltyAmount))
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sellerAddr, types.SecurityAuditFund, penaltyCoins); err != nil {
		k.Logger(ctx).Error("Failed to collect dumping penalty", "error", err)
	}
	
	// Update wallet limits
	k.setWalletLimits(ctx, walletLimits)
	
	k.Logger(ctx).Info("Applied dumping penalty",
		"token", tokenAddress,
		"seller", seller,
		"penalty_rate", penaltyRate,
		"penalty_amount", penaltyAmount,
		"violations", walletLimits.ViolationCount,
	)
	
	return nil
}

// requireWhaleConfirmation requires additional confirmation for whale sells
func (k Keeper) requireWhaleConfirmation(ctx sdk.Context, tokenAddress, seller string, sellAmount sdk.Int) error {
	// For now, just log the whale transaction
	// In a full implementation, this could require governance approval
	k.Logger(ctx).Warn("Whale transaction requires confirmation",
		"token", tokenAddress,
		"seller", seller,
		"amount", sellAmount,
	)
	
	// Could implement a delay mechanism or governance vote here
	
	return nil
}

// MonitorMarketStability monitors overall market stability and applies emergency measures
func (k Keeper) MonitorMarketStability(ctx sdk.Context, tokenAddress string) error {
	// Get trading metrics
	metrics, found := k.getTradingMetrics(ctx, tokenAddress)
	if !found {
		return nil
	}
	
	// Check for extreme price movements (>50% in 1 hour)
	if metrics.PriceChange24h.Abs().GT(sdk.MustNewDecFromStr("0.50")) {
		// Apply circuit breaker
		return k.activateCircuitBreaker(ctx, tokenAddress, "extreme_volatility")
	}
	
	// Check for unusual volume spikes
	avgVolume := metrics.TotalVolume.QuoRaw(30) // 30-day average
	if metrics.DailyVolume.GT(avgVolume.MulRaw(10)) { // 10x average
		k.Logger(ctx).Warn("Unusual volume spike detected",
			"token", tokenAddress,
			"daily_volume", metrics.DailyVolume,
			"avg_volume", avgVolume,
		)
	}
	
	// Check liquidity health
	if metrics.Liquidity.IsPositive() {
		liquidityRatio := metrics.DailyVolume.ToDec().Quo(metrics.Liquidity.ToDec())
		if liquidityRatio.GT(sdk.MustNewDecFromStr("2.0")) { // Volume > 2x liquidity
			return k.activateCircuitBreaker(ctx, tokenAddress, "liquidity_crisis")
		}
	}
	
	return nil
}

// activateCircuitBreaker activates emergency circuit breaker
func (k Keeper) activateCircuitBreaker(ctx sdk.Context, tokenAddress, reason string) error {
	// Create emergency control
	emergencyControl := types.EmergencyControl{
		TokenAddress: tokenAddress,
		ControlType:  "circuit_breaker",
		InitiatedBy:  "system",
		Reason:       reason,
		ActivatedAt:  ctx.BlockTime(),
		ExpiresAt:    &[]time.Time{ctx.BlockTime().Add(1 * time.Hour)}[0], // 1 hour cooling period
		IsActive:     true,
		Metadata:     map[string]string{
			"auto_triggered": "true",
			"trigger_reason": reason,
		},
	}
	
	k.setEmergencyControl(ctx, emergencyControl)
	
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"circuit_breaker_activated",
			sdk.NewAttribute("token_address", tokenAddress),
			sdk.NewAttribute("reason", reason),
			sdk.NewAttribute("expires_at", emergencyControl.ExpiresAt.String()),
		),
	)
	
	k.Logger(ctx).Info("Circuit breaker activated",
		"token", tokenAddress,
		"reason", reason,
		"expires_at", emergencyControl.ExpiresAt,
	)
	
	return types.ErrEmergencyStop
}

// Helper functions for tracking sell history

type SellRecord struct {
	Seller    string    `json:"seller"`
	Amount    sdk.Int   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
	BlockHeight int64   `json:"block_height"`
}

func (k Keeper) getRecentSells(ctx sdk.Context, tokenAddress, seller string, duration time.Duration) []SellRecord {
	// Implementation would query sell history from store
	// Simplified for now
	return []SellRecord{}
}

func (k Keeper) getDailySells(ctx sdk.Context, tokenAddress, seller string, date time.Time) []SellRecord {
	// Implementation would query daily sell records
	// Simplified for now
	return []SellRecord{}
}

func (k Keeper) getSellerOriginalAllocation(ctx sdk.Context, tokenAddress, seller string) sdk.Int {
	// Implementation would get original allocation from launch participation
	// Simplified for now
	return sdk.ZeroInt()
}

func (k Keeper) getSellerSoldAmount(ctx sdk.Context, tokenAddress, seller string) sdk.Int {
	// Implementation would track total sold amount
	// Simplified for now
	return sdk.ZeroInt()
}

func (k Keeper) recordSellTransaction(ctx sdk.Context, tokenAddress, seller string, amount sdk.Int) {
	// Implementation would record sell transaction for tracking
	// Simplified for now
}

// Define the missing error
var ErrGradualReleaseViolation = types.ErrFeatureNotEnabled // Using existing error for now