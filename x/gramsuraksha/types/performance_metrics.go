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

package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ChainPerformanceMetrics tracks performance indicators for dynamic payout calculation
type ChainPerformanceMetrics struct {
	// Chain health indicators
	BlockTime          int64   `json:"block_time"`          // Average block time in ms
	TotalTransactions  uint64  `json:"total_transactions"`  // Total transactions processed
	ActiveValidators   uint32  `json:"active_validators"`   // Number of active validators
	NetworkUptime      sdk.Dec `json:"network_uptime"`      // Network uptime percentage (0-100)
	
	// Financial performance
	TotalRevenue       sdk.Coins `json:"total_revenue"`       // Total revenue generated
	MonthlyRevenue     sdk.Coins `json:"monthly_revenue"`     // Last month's revenue
	RevenueTrend       sdk.Dec   `json:"revenue_trend"`       // Revenue growth/decline rate
	
	// Risk indicators
	TotalWriteoffs     sdk.Coins `json:"total_writeoffs"`     // Total loan writeoffs
	MonthlyWriteoffs   sdk.Coins `json:"monthly_writeoffs"`   // Last month's writeoffs
	WriteoffRate       sdk.Dec   `json:"writeoff_rate"`       // Writeoff percentage
	
	// Fraud detection
	DetectedFrauds     uint64    `json:"detected_frauds"`     // Number of detected fraud cases
	CollaborativeFraud uint64    `json:"collaborative_fraud"` // Detected collaborative fraud cases
	FraudLossAmount    sdk.Coins `json:"fraud_loss_amount"`   // Total loss due to fraud
	FraudRate          sdk.Dec   `json:"fraud_rate"`          // Fraud rate percentage
	
	// Community health
	ActivePensions     uint64    `json:"active_pensions"`     // Active pension accounts
	PensionPayouts     sdk.Coins `json:"pension_payouts"`     // Total pension payouts made
	PayoutSuccessRate  sdk.Dec   `json:"payout_success_rate"` // Successful payout percentage
	
	// Calculated scores
	HealthScore        sdk.Dec   `json:"health_score"`        // Overall chain health (0-100)
	RiskScore          sdk.Dec   `json:"risk_score"`          // Risk level (0-100, lower is better)
	PerformanceScore   sdk.Dec   `json:"performance_score"`   // Performance rating (0-100)
	
	LastUpdated        int64     `json:"last_updated"`        // Unix timestamp of last update
}

// DynamicPayoutParams configures how payouts are calculated based on performance
type DynamicPayoutParams struct {
	// Base configuration
	MinimumPayout      sdk.Dec `json:"minimum_payout"`      // Minimum guaranteed payout (8%)
	MaximumPayout      sdk.Dec `json:"maximum_payout"`      // Maximum possible payout (50%)
	BaselinePayout     sdk.Dec `json:"baseline_payout"`     // Normal conditions payout (30%)
	
	// Performance multipliers
	HealthMultiplier   sdk.Dec `json:"health_multiplier"`   // Impact of health score
	RevenueMultiplier  sdk.Dec `json:"revenue_multiplier"`  // Impact of revenue performance
	
	// Risk deductions
	WriteoffImpact     sdk.Dec `json:"writeoff_impact"`     // Deduction per % of writeoffs
	FraudImpact        sdk.Dec `json:"fraud_impact"`        // Deduction per fraud case
	CollabFraudImpact  sdk.Dec `json:"collab_fraud_impact"` // Heavy penalty for collaborative fraud
	
	// Thresholds
	CriticalRiskLevel  sdk.Dec `json:"critical_risk_level"` // Risk score that triggers minimum payout
	OptimalHealthLevel sdk.Dec `json:"optimal_health_level"` // Health score for maximum payout
}

// DefaultDynamicPayoutParams returns the default parameters
func DefaultDynamicPayoutParams() DynamicPayoutParams {
	return DynamicPayoutParams{
		MinimumPayout:      sdk.NewDecWithPrec(8, 2),   // 8%
		MaximumPayout:      sdk.NewDecWithPrec(50, 2),  // 50%
		BaselinePayout:     sdk.NewDecWithPrec(30, 2),  // 30%
		HealthMultiplier:   sdk.NewDecWithPrec(20, 2),  // 0.2x multiplier
		RevenueMultiplier:  sdk.NewDecWithPrec(15, 2),  // 0.15x multiplier
		WriteoffImpact:     sdk.NewDecWithPrec(5, 2),   // -5% per 1% writeoff
		FraudImpact:        sdk.NewDecWithPrec(2, 2),   // -2% per fraud case (per 1000 accounts)
		CollabFraudImpact:  sdk.NewDecWithPrec(10, 2),  // -10% for collaborative fraud
		CriticalRiskLevel:  sdk.NewDecWithPrec(80, 2),  // 80% risk score
		OptimalHealthLevel: sdk.NewDecWithPrec(90, 2),  // 90% health score
	}
}

// CalculateDynamicPayout calculates the actual payout percentage based on performance
func CalculateDynamicPayout(metrics ChainPerformanceMetrics, params DynamicPayoutParams) sdk.Dec {
	// Start with baseline
	payout := params.BaselinePayout
	
	// Apply health score bonus (up to +20% at 100% health)
	healthBonus := metrics.HealthScore.Sub(sdk.NewDecWithPrec(50, 2)).Mul(params.HealthMultiplier).Quo(sdk.NewDec(100))
	if healthBonus.IsPositive() {
		payout = payout.Add(healthBonus)
	}
	
	// Apply revenue performance bonus (up to +15% for strong growth)
	if metrics.RevenueTrend.IsPositive() {
		revenueBonus := metrics.RevenueTrend.Mul(params.RevenueMultiplier)
		payout = payout.Add(revenueBonus)
	}
	
	// Apply writeoff penalty
	writeoffPenalty := metrics.WriteoffRate.Mul(params.WriteoffImpact)
	payout = payout.Sub(writeoffPenalty)
	
	// Apply fraud penalty (normalized per 1000 accounts)
	fraudsPerThousand := sdk.NewDec(int64(metrics.DetectedFrauds)).Mul(sdk.NewDec(1000)).Quo(sdk.NewDec(int64(metrics.ActivePensions)))
	fraudPenalty := fraudsPerThousand.Mul(params.FraudImpact)
	payout = payout.Sub(fraudPenalty)
	
	// Apply collaborative fraud penalty
	if metrics.CollaborativeFraud > 0 {
		collabPenalty := sdk.NewDec(int64(metrics.CollaborativeFraud)).Mul(params.CollabFraudImpact).Quo(sdk.NewDec(100))
		payout = payout.Sub(collabPenalty)
	}
	
	// Check critical risk level
	if metrics.RiskScore.GTE(params.CriticalRiskLevel) {
		// Force to minimum payout
		return params.MinimumPayout
	}
	
	// Ensure within bounds
	if payout.LT(params.MinimumPayout) {
		return params.MinimumPayout
	}
	if payout.GT(params.MaximumPayout) {
		return params.MaximumPayout
	}
	
	return payout
}

// GetPayoutMessage returns appropriate message for verifiers based on payout level
func GetPayoutMessage(payoutRate sdk.Dec, metrics ChainPerformanceMetrics) string {
	if payoutRate.LTE(sdk.NewDecWithPrec(10, 2)) {
		return fmt.Sprintf(
			"⚠️ CRITICAL: Pension payouts reduced to %s%% due to high risk factors:\n"+
			"- Writeoff Rate: %s%%\n"+
			"- Fraud Cases: %d (Collaborative: %d)\n"+
			"- Risk Score: %s/100\n"+
			"Immediate action required to improve chain health.",
			payoutRate.Mul(sdk.NewDec(100)).TruncateInt().String(),
			metrics.WriteoffRate.Mul(sdk.NewDec(100)).TruncateInt().String(),
			metrics.DetectedFrauds,
			metrics.CollaborativeFraud,
			metrics.RiskScore.TruncateInt().String(),
		)
	} else if payoutRate.LTE(sdk.NewDecWithPrec(25, 2)) {
		return fmt.Sprintf(
			"⚡ WARNING: Pension payouts at %s%% due to performance issues:\n"+
			"- Health Score: %s/100\n"+
			"- Revenue Trend: %s%%\n"+
			"- Writeoffs: %s\n"+
			"Please review risk management practices.",
			payoutRate.Mul(sdk.NewDec(100)).TruncateInt().String(),
			metrics.HealthScore.TruncateInt().String(),
			metrics.RevenueTrend.Mul(sdk.NewDec(100)).TruncateInt().String(),
			metrics.MonthlyWriteoffs.String(),
		)
	} else if payoutRate.GTE(sdk.NewDecWithPrec(45, 2)) {
		return fmt.Sprintf(
			"🎉 EXCELLENT: Pension payouts at %s%% due to outstanding performance:\n"+
			"- Health Score: %s/100\n"+
			"- Revenue Growth: %s%%\n"+
			"- Low Risk Score: %s/100\n"+
			"Keep up the great work!",
			payoutRate.Mul(sdk.NewDec(100)).TruncateInt().String(),
			metrics.HealthScore.TruncateInt().String(),
			metrics.RevenueTrend.Mul(sdk.NewDec(100)).TruncateInt().String(),
			metrics.RiskScore.TruncateInt().String(),
		)
	} else {
		return fmt.Sprintf(
			"✅ NORMAL: Pension payouts at %s%% based on current performance:\n"+
			"- Health Score: %s/100\n"+
			"- Performance Score: %s/100\n"+
			"- Payout Success Rate: %s%%\n"+
			"Chain operating within normal parameters.",
			payoutRate.Mul(sdk.NewDec(100)).TruncateInt().String(),
			metrics.HealthScore.TruncateInt().String(),
			metrics.PerformanceScore.TruncateInt().String(),
			metrics.PayoutSuccessRate.Mul(sdk.NewDec(100)).TruncateInt().String(),
		)
	}
}

// ValidateMetrics ensures metrics are within valid ranges
func (m ChainPerformanceMetrics) Validate() error {
	if m.NetworkUptime.LT(sdk.ZeroDec()) || m.NetworkUptime.GT(sdk.NewDec(100)) {
		return fmt.Errorf("network uptime must be between 0-100")
	}
	if m.WriteoffRate.LT(sdk.ZeroDec()) || m.WriteoffRate.GT(sdk.OneDec()) {
		return fmt.Errorf("writeoff rate must be between 0-1")
	}
	if m.FraudRate.LT(sdk.ZeroDec()) || m.FraudRate.GT(sdk.OneDec()) {
		return fmt.Errorf("fraud rate must be between 0-1")
	}
	if m.HealthScore.LT(sdk.ZeroDec()) || m.HealthScore.GT(sdk.NewDec(100)) {
		return fmt.Errorf("health score must be between 0-100")
	}
	if m.RiskScore.LT(sdk.ZeroDec()) || m.RiskScore.GT(sdk.NewDec(100)) {
		return fmt.Errorf("risk score must be between 0-100")
	}
	if m.PerformanceScore.LT(sdk.ZeroDec()) || m.PerformanceScore.GT(sdk.NewDec(100)) {
		return fmt.Errorf("performance score must be between 0-100")
	}
	return nil
}