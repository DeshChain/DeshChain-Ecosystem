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
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/deshchain/x/gramsuraksha/types"
)

// VerifierLevel represents different community levels
type VerifierLevel string

const (
	VerifierLevelVillage  VerifierLevel = "village"
	VerifierLevelDistrict VerifierLevel = "district"
	VerifierLevelState    VerifierLevel = "state"
	VerifierLevelNational VerifierLevel = "national"
)

// NotifyVerifiersAboutPayoutChange sends notifications to verifiers at different levels
func (k Keeper) NotifyVerifiersAboutPayoutChange(ctx sdk.Context, oldRate, newRate sdk.Dec) {
	metrics := k.GetChainPerformanceMetrics(ctx)
	
	// Get payout message based on current rate
	generalMessage := types.GetPayoutMessage(newRate, metrics)
	
	// Send notifications to different verifier levels
	k.notifyVillageVerifiers(ctx, oldRate, newRate, generalMessage)
	k.notifyDistrictVerifiers(ctx, oldRate, newRate, metrics)
	k.notifyStateVerifiers(ctx, oldRate, newRate, metrics)
	k.notifyNationalVerifiers(ctx, oldRate, newRate, metrics)
	
	// Emit event for transparency
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypePayoutRateChanged,
			sdk.NewAttribute("old_rate", oldRate.String()),
			sdk.NewAttribute("new_rate", newRate.String()),
			sdk.NewAttribute("health_score", metrics.HealthScore.String()),
			sdk.NewAttribute("risk_score", metrics.RiskScore.String()),
			sdk.NewAttribute("message", generalMessage),
		),
	)
}

// notifyVillageVerifiers sends detailed messages to village-level verifiers
func (k Keeper) notifyVillageVerifiers(ctx sdk.Context, oldRate, newRate sdk.Dec, message string) {
	// Village verifiers get the most detailed information
	detailedMsg := fmt.Sprintf(
		"📢 GRAM SURAKSHA PAYOUT UPDATE\n\n"+
		"Dear Village Verifier,\n\n"+
		"Pension payout rates have been adjusted:\n"+
		"Previous Rate: %s%%\n"+
		"New Rate: %s%%\n\n"+
		"%s\n\n"+
		"Action Items for Village Level:\n"+
		"1. Inform all pension participants about the rate change\n"+
		"2. Explain the reasons for the adjustment\n"+
		"3. Help participants understand chain performance metrics\n"+
		"4. Report any suspicious activities or fraud attempts\n"+
		"5. Encourage timely contributions to improve chain health\n\n"+
		"For questions, contact your District Coordinator.\n"+
		"Jai Hind! 🇮🇳",
		oldRate.Mul(sdk.NewDec(100)).TruncateInt().String(),
		newRate.Mul(sdk.NewDec(100)).TruncateInt().String(),
		message,
	)
	
	// Store message for village verifiers
	k.storeVerifierMessage(ctx, VerifierLevelVillage, detailedMsg)
}

// notifyDistrictVerifiers sends summary messages to district-level verifiers
func (k Keeper) notifyDistrictVerifiers(ctx sdk.Context, oldRate, newRate sdk.Dec, metrics types.ChainPerformanceMetrics) {
	// District verifiers get summary information
	summaryMsg := fmt.Sprintf(
		"📊 DISTRICT PENSION UPDATE\n\n"+
		"Payout Rate Change: %s%% → %s%%\n\n"+
		"District Performance Summary:\n"+
		"• Active Pensions: %d\n"+
		"• Writeoff Rate: %s%%\n"+
		"• Fraud Cases: %d\n"+
		"• Health Score: %s/100\n\n"+
		"Required Actions:\n"+
		"1. Review village-level performance\n"+
		"2. Coordinate fraud prevention measures\n"+
		"3. Monitor collection efficiency\n"+
		"4. Submit monthly performance report\n\n"+
		"Target: Achieve 40%% payout rate by improving metrics.",
		oldRate.Mul(sdk.NewDec(100)).TruncateInt().String(),
		newRate.Mul(sdk.NewDec(100)).TruncateInt().String(),
		metrics.ActivePensions,
		metrics.WriteoffRate.Mul(sdk.NewDec(100)).TruncateInt().String(),
		metrics.DetectedFrauds,
		metrics.HealthScore.TruncateInt().String(),
	)
	
	k.storeVerifierMessage(ctx, VerifierLevelDistrict, summaryMsg)
}

// notifyStateVerifiers sends strategic messages to state-level verifiers
func (k Keeper) notifyStateVerifiers(ctx sdk.Context, oldRate, newRate sdk.Dec, metrics types.ChainPerformanceMetrics) {
	// State verifiers get strategic overview
	strategicMsg := fmt.Sprintf(
		"🏛️ STATE PENSION ADVISORY\n\n"+
		"Gram Suraksha Payout Adjustment\n"+
		"Current Rate: %s%% (Change: %s%%)\n\n"+
		"State-Level Metrics:\n"+
		"• Performance Score: %s/100\n"+
		"• Risk Score: %s/100\n"+
		"• Revenue Trend: %s%%\n\n"+
		"Strategic Priorities:\n"+
		"1. Enhance district coordination\n"+
		"2. Implement fraud detection training\n"+
		"3. Optimize resource allocation\n"+
		"4. Quarterly review with National Committee\n\n"+
		"Excellence Target: 45%% payout rate",
		newRate.Mul(sdk.NewDec(100)).TruncateInt().String(),
		newRate.Sub(oldRate).Mul(sdk.NewDec(100)).TruncateInt().String(),
		metrics.PerformanceScore.TruncateInt().String(),
		metrics.RiskScore.TruncateInt().String(),
		metrics.RevenueTrend.Mul(sdk.NewDec(100)).TruncateInt().String(),
	)
	
	k.storeVerifierMessage(ctx, VerifierLevelState, strategicMsg)
}

// notifyNationalVerifiers sends executive summary to national-level verifiers
func (k Keeper) notifyNationalVerifiers(ctx sdk.Context, oldRate, newRate sdk.Dec, metrics types.ChainPerformanceMetrics) {
	// National verifiers get executive summary
	executiveMsg := fmt.Sprintf(
		"🇮🇳 NATIONAL PENSION BOARD UPDATE\n\n"+
		"Gram Suraksha Performance Report\n"+
		"Payout Rate: %s%%\n\n"+
		"National KPIs:\n"+
		"• Total Participants: %s\n"+
		"• Chain Health: %s/100\n"+
		"• Risk Level: %s/100\n"+
		"• Payout Success: %s%%\n\n"+
		"Board Recommendations:\n"+
		"• %s\n\n"+
		"Next Review: Monthly Board Meeting",
		newRate.Mul(sdk.NewDec(100)).TruncateInt().String(),
		formatLargeNumber(metrics.ActivePensions),
		metrics.HealthScore.TruncateInt().String(),
		metrics.RiskScore.TruncateInt().String(),
		metrics.PayoutSuccessRate.Mul(sdk.NewDec(100)).TruncateInt().String(),
		getNationalRecommendation(newRate),
	)
	
	k.storeVerifierMessage(ctx, VerifierLevelNational, executiveMsg)
}

// storeVerifierMessage stores a message for verifiers at a specific level
func (k Keeper) storeVerifierMessage(ctx sdk.Context, level VerifierLevel, message string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixVerifierMessages)
	
	key := append([]byte(level), ctx.BlockTime().Format("20060102150405")...)
	store.Set(key, []byte(message))
	
	// Also store as latest message for quick access
	latestKey := append([]byte("latest_"), []byte(level)...)
	store.Set(latestKey, []byte(message))
}

// GetLatestVerifierMessage retrieves the latest message for a verifier level
func (k Keeper) GetLatestVerifierMessage(ctx sdk.Context, level VerifierLevel) string {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixVerifierMessages)
	
	latestKey := append([]byte("latest_"), []byte(level)...)
	bz := store.Get(latestKey)
	if bz == nil {
		return "No messages available"
	}
	
	return string(bz)
}

// GetVerifierMessageHistory retrieves message history for a verifier level
func (k Keeper) GetVerifierMessageHistory(ctx sdk.Context, level VerifierLevel, limit uint32) []string {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixVerifierMessages)
	
	messages := []string{}
	iterator := sdk.KVStoreReversePrefixIterator(store, []byte(level))
	defer iterator.Close()
	
	count := uint32(0)
	for ; iterator.Valid() && count < limit; iterator.Next() {
		// Skip latest message keys
		key := string(iterator.Key())
		if len(key) > 7 && key[:7] == "latest_" {
			continue
		}
		
		messages = append(messages, string(iterator.Value()))
		count++
	}
	
	return messages
}

// Helper functions

// formatLargeNumber formats large numbers with appropriate units
func formatLargeNumber(num uint64) string {
	if num >= 10000000 { // 1 crore
		return fmt.Sprintf("%.2f Cr", float64(num)/10000000)
	} else if num >= 100000 { // 1 lakh
		return fmt.Sprintf("%.2f L", float64(num)/100000)
	} else if num >= 1000 {
		return fmt.Sprintf("%.1fK", float64(num)/1000)
	}
	return fmt.Sprintf("%d", num)
}

// getNationalRecommendation provides recommendations based on payout rate
func getNationalRecommendation(rate sdk.Dec) string {
	ratePercent := rate.Mul(sdk.NewDec(100)).TruncateInt().Int64()
	
	switch {
	case ratePercent >= 45:
		return "Maintain current excellence standards. Consider expansion."
	case ratePercent >= 35:
		return "Good performance. Focus on fraud prevention and efficiency."
	case ratePercent >= 25:
		return "Moderate performance. Implement risk management measures."
	case ratePercent >= 15:
		return "Below target. Urgent review of operations required."
	default:
		return "Critical situation. Emergency measures needed immediately."
	}
}

// MonthlyPayoutReview performs monthly review and sends notifications
func (k Keeper) MonthlyPayoutReview(ctx sdk.Context) {
	// Get current metrics
	metrics := k.GetChainPerformanceMetrics(ctx)
	params := k.GetDynamicPayoutParams(ctx)
	
	// Calculate new payout rate
	newRate := types.CalculateDynamicPayout(metrics, params)
	
	// Get previous rate
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixPayoutHistory)
	var oldRate sdk.Dec
	bz := store.Get(types.KeyLastPayoutRate)
	if bz != nil {
		oldRate, _ = sdk.NewDecFromStr(string(bz))
	} else {
		oldRate = params.BaselinePayout
	}
	
	// Only notify if rate changed significantly (>1%)
	rateDiff := newRate.Sub(oldRate).Abs()
	if rateDiff.GT(sdk.NewDecWithPrec(1, 2)) {
		k.NotifyVerifiersAboutPayoutChange(ctx, oldRate, newRate)
	}
	
	// Store new rate as current
	store.Set(types.KeyLastPayoutRate, []byte(newRate.String()))
	
	// Log the review
	ctx.Logger().Info("Monthly payout review completed",
		"old_rate", oldRate.String(),
		"new_rate", newRate.String(),
		"rate_diff", rateDiff.String(),
		"notifications_sent", rateDiff.GT(sdk.NewDecWithPrec(1, 2)),
	)
}