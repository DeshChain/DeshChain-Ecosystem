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
	"fmt"
	"time"
	
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/deshchain/x/gramsuraksha/types"
)

// GetChainPerformanceMetrics retrieves current chain performance metrics
func (k Keeper) GetChainPerformanceMetrics(ctx sdk.Context) types.ChainPerformanceMetrics {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixPerformanceMetrics)
	
	var metrics types.ChainPerformanceMetrics
	bz := store.Get(types.KeyCurrentMetrics)
	if bz != nil {
		k.cdc.MustUnmarshal(bz, &metrics)
	} else {
		// Return default metrics if none exist
		metrics = k.InitializeDefaultMetrics(ctx)
	}
	
	return metrics
}

// SetChainPerformanceMetrics updates the chain performance metrics
func (k Keeper) SetChainPerformanceMetrics(ctx sdk.Context, metrics types.ChainPerformanceMetrics) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixPerformanceMetrics)
	metrics.LastUpdated = ctx.BlockTime().Unix()
	bz := k.cdc.MustMarshal(&metrics)
	store.Set(types.KeyCurrentMetrics, bz)
}

// UpdateChainPerformanceMetrics updates metrics based on current chain state
func (k Keeper) UpdateChainPerformanceMetrics(ctx sdk.Context) {
	metrics := k.GetChainPerformanceMetrics(ctx)
	
	// Update block time (simplified - in production, calculate average)
	metrics.BlockTime = 5000 // 5 seconds average
	
	// Update validator count
	validators := k.stakingKeeper.GetAllValidators(ctx)
	activeCount := 0
	for _, val := range validators {
		if val.IsBonded() {
			activeCount++
		}
	}
	metrics.ActiveValidators = uint32(activeCount)
	
	// Calculate network uptime (simplified - assume 99.9% for now)
	metrics.NetworkUptime = sdk.NewDecWithPrec(999, 3)
	
	// Update pension statistics
	metrics.ActivePensions = k.GetActivePensionCount(ctx)
	
	// Calculate health score
	metrics.HealthScore = k.calculateHealthScore(metrics)
	
	// Calculate risk score
	metrics.RiskScore = k.calculateRiskScore(metrics)
	
	// Calculate performance score
	metrics.PerformanceScore = k.calculatePerformanceScore(metrics)
	
	// Save updated metrics
	k.SetChainPerformanceMetrics(ctx, metrics)
}

// InitializeDefaultMetrics creates default metrics for a new chain
func (k Keeper) InitializeDefaultMetrics(ctx sdk.Context) types.ChainPerformanceMetrics {
	return types.ChainPerformanceMetrics{
		BlockTime:          5000,
		TotalTransactions:  0,
		ActiveValidators:   21,
		NetworkUptime:      sdk.NewDecWithPrec(999, 3), // 99.9%
		TotalRevenue:       sdk.NewCoins(),
		MonthlyRevenue:     sdk.NewCoins(),
		RevenueTrend:       sdk.ZeroDec(),
		TotalWriteoffs:     sdk.NewCoins(),
		MonthlyWriteoffs:   sdk.NewCoins(),
		WriteoffRate:       sdk.ZeroDec(),
		DetectedFrauds:     0,
		CollaborativeFraud: 0,
		FraudLossAmount:    sdk.NewCoins(),
		FraudRate:          sdk.ZeroDec(),
		ActivePensions:     0,
		PensionPayouts:     sdk.NewCoins(),
		PayoutSuccessRate:  sdk.OneDec(), // 100% initially
		HealthScore:        sdk.NewDec(90),
		RiskScore:          sdk.NewDec(10),
		PerformanceScore:   sdk.NewDec(80),
		LastUpdated:        ctx.BlockTime().Unix(),
	}
}

// RecordWriteoff records a loan writeoff and updates metrics
func (k Keeper) RecordWriteoff(ctx sdk.Context, amount sdk.Coin, reason string) {
	metrics := k.GetChainPerformanceMetrics(ctx)
	
	// Update writeoff amounts
	metrics.TotalWriteoffs = metrics.TotalWriteoffs.Add(amount)
	metrics.MonthlyWriteoffs = metrics.MonthlyWriteoffs.Add(amount)
	
	// Calculate writeoff rate (writeoffs / total loans)
	totalLoans := k.GetTotalLoanAmount(ctx)
	if !totalLoans.IsZero() {
		writeoffAmount := metrics.TotalWriteoffs.AmountOf(amount.Denom)
		metrics.WriteoffRate = writeoffAmount.ToDec().Quo(totalLoans.ToDec())
	}
	
	k.SetChainPerformanceMetrics(ctx, metrics)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeWriteoffRecorded,
			sdk.NewAttribute("amount", amount.String()),
			sdk.NewAttribute("reason", reason),
			sdk.NewAttribute("writeoff_rate", metrics.WriteoffRate.String()),
		),
	)
}

// RecordFraud records a fraud case and updates metrics
func (k Keeper) RecordFraud(ctx sdk.Context, fraudType string, lossAmount sdk.Coin, isCollaborative bool) {
	metrics := k.GetChainPerformanceMetrics(ctx)
	
	// Update fraud counts
	metrics.DetectedFrauds++
	if isCollaborative {
		metrics.CollaborativeFraud++
	}
	
	// Update fraud loss
	metrics.FraudLossAmount = metrics.FraudLossAmount.Add(lossAmount)
	
	// Calculate fraud rate
	if metrics.ActivePensions > 0 {
		metrics.FraudRate = sdk.NewDec(int64(metrics.DetectedFrauds)).Quo(sdk.NewDec(int64(metrics.ActivePensions)))
	}
	
	k.SetChainPerformanceMetrics(ctx, metrics)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeFraudDetected,
			sdk.NewAttribute("fraud_type", fraudType),
			sdk.NewAttribute("loss_amount", lossAmount.String()),
			sdk.NewAttribute("is_collaborative", fmt.Sprintf("%t", isCollaborative)),
			sdk.NewAttribute("total_frauds", fmt.Sprintf("%d", metrics.DetectedFrauds)),
		),
	)
}

// calculateHealthScore calculates overall chain health (0-100)
func (k Keeper) calculateHealthScore(metrics types.ChainPerformanceMetrics) sdk.Dec {
	score := sdk.NewDec(100)
	
	// Deduct for low validator count
	if metrics.ActiveValidators < 20 {
		score = score.Sub(sdk.NewDec(20))
	} else if metrics.ActiveValidators < 50 {
		score = score.Sub(sdk.NewDec(10))
	}
	
	// Deduct for network downtime
	downtimePercent := sdk.NewDec(100).Sub(metrics.NetworkUptime)
	score = score.Sub(downtimePercent.Mul(sdk.NewDec(2))) // 2x penalty for downtime
	
	// Deduct for low payout success rate
	if metrics.PayoutSuccessRate.LT(sdk.NewDecWithPrec(95, 2)) {
		penalty := sdk.NewDec(100).Sub(metrics.PayoutSuccessRate.Mul(sdk.NewDec(100)))
		score = score.Sub(penalty)
	}
	
	// Ensure score is within bounds
	if score.LT(sdk.ZeroDec()) {
		return sdk.ZeroDec()
	}
	if score.GT(sdk.NewDec(100)) {
		return sdk.NewDec(100)
	}
	
	return score
}

// calculateRiskScore calculates risk level (0-100, lower is better)
func (k Keeper) calculateRiskScore(metrics types.ChainPerformanceMetrics) sdk.Dec {
	score := sdk.ZeroDec()
	
	// Add risk for writeoffs
	writeoffRisk := metrics.WriteoffRate.Mul(sdk.NewDec(100))
	score = score.Add(writeoffRisk.Mul(sdk.NewDec(3))) // 3x weight for writeoffs
	
	// Add risk for fraud
	fraudRisk := metrics.FraudRate.Mul(sdk.NewDec(100))
	score = score.Add(fraudRisk.Mul(sdk.NewDec(5))) // 5x weight for fraud
	
	// Add heavy penalty for collaborative fraud
	if metrics.CollaborativeFraud > 0 {
		collabPenalty := sdk.NewDec(int64(metrics.CollaborativeFraud * 10))
		score = score.Add(collabPenalty)
	}
	
	// Add risk for revenue decline
	if metrics.RevenueTrend.IsNegative() {
		revenuePenalty := metrics.RevenueTrend.Abs().Mul(sdk.NewDec(50))
		score = score.Add(revenuePenalty)
	}
	
	// Ensure score is within bounds
	if score.GT(sdk.NewDec(100)) {
		return sdk.NewDec(100)
	}
	
	return score
}

// calculatePerformanceScore calculates performance rating (0-100)
func (k Keeper) calculatePerformanceScore(metrics types.ChainPerformanceMetrics) sdk.Dec {
	score := sdk.NewDec(50) // Start at baseline
	
	// Add points for revenue growth
	if metrics.RevenueTrend.IsPositive() {
		revenueBonus := metrics.RevenueTrend.Mul(sdk.NewDec(100))
		if revenueBonus.GT(sdk.NewDec(30)) {
			revenueBonus = sdk.NewDec(30) // Cap at 30 points
		}
		score = score.Add(revenueBonus)
	}
	
	// Add points for low writeoff rate
	if metrics.WriteoffRate.LT(sdk.NewDecWithPrec(1, 2)) { // Less than 1%
		writeoffBonus := sdk.NewDec(20)
		score = score.Add(writeoffBonus)
	} else if metrics.WriteoffRate.LT(sdk.NewDecWithPrec(3, 2)) { // Less than 3%
		writeoffBonus := sdk.NewDec(10)
		score = score.Add(writeoffBonus)
	}
	
	// Add points for high payout success
	if metrics.PayoutSuccessRate.GTE(sdk.NewDecWithPrec(98, 2)) {
		successBonus := sdk.NewDec(20)
		score = score.Add(successBonus)
	} else if metrics.PayoutSuccessRate.GTE(sdk.NewDecWithPrec(95, 2)) {
		successBonus := sdk.NewDec(10)
		score = score.Add(successBonus)
	}
	
	// Deduct points for fraud
	fraudPenalty := sdk.NewDec(int64(metrics.DetectedFrauds)).Quo(sdk.NewDec(100))
	score = score.Sub(fraudPenalty)
	
	// Ensure score is within bounds
	if score.LT(sdk.ZeroDec()) {
		return sdk.ZeroDec()
	}
	if score.GT(sdk.NewDec(100)) {
		return sdk.NewDec(100)
	}
	
	return score
}

// GetDynamicPayoutParams returns the current payout parameters
func (k Keeper) GetDynamicPayoutParams(ctx sdk.Context) types.DynamicPayoutParams {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixPayoutParams)
	
	var params types.DynamicPayoutParams
	bz := store.Get(types.KeyPayoutParams)
	if bz != nil {
		k.cdc.MustUnmarshal(bz, &params)
	} else {
		// Return default params if none exist
		params = types.DefaultDynamicPayoutParams()
		k.SetDynamicPayoutParams(ctx, params)
	}
	
	return params
}

// SetDynamicPayoutParams updates the payout parameters
func (k Keeper) SetDynamicPayoutParams(ctx sdk.Context, params types.DynamicPayoutParams) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixPayoutParams)
	bz := k.cdc.MustMarshal(&params)
	store.Set(types.KeyPayoutParams, bz)
}

// Helper functions

// GetActivePensionCount returns the number of active pension accounts
func (k Keeper) GetActivePensionCount(ctx sdk.Context) uint64 {
	// Implementation depends on how pensions are stored
	count := uint64(0)
	k.IterateAllParticipants(ctx, func(participant types.SurakshaParticipant) bool {
		if participant.Status == types.StatusActive {
			count++
		}
		return false
	})
	return count
}

// GetTotalLoanAmount returns total outstanding loan amount
func (k Keeper) GetTotalLoanAmount(ctx sdk.Context) sdk.Int {
	// This would integrate with the lending module
	// For now, return a placeholder
	return sdk.NewInt(1000000000000) // 1M NAMO
}