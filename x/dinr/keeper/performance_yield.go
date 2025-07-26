package keeper

import (
	"fmt"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/dinr/types"
)

// PerformanceYieldCalculator calculates dynamic yield based on platform performance
type PerformanceYieldCalculator struct {
	keeper *Keeper
}

// NewPerformanceYieldCalculator creates a new performance yield calculator
func NewPerformanceYieldCalculator(k *Keeper) *PerformanceYieldCalculator {
	return &PerformanceYieldCalculator{
		keeper: k,
	}
}

// CalculateCurrentYield calculates the current yield rate based on performance metrics
func (pyc *PerformanceYieldCalculator) CalculateCurrentYield(ctx sdk.Context) (sdk.Dec, error) {
	// Get performance metrics
	metrics := pyc.GetPerformanceMetrics(ctx)
	
	// Calculate weighted yield
	// Platform profitability: 40% weight
	profitabilityScore := pyc.calculateProfitabilityScore(ctx, metrics.PlatformRevenue, metrics.PlatformExpenses)
	
	// Trading volume growth: 20% weight
	volumeGrowthScore := pyc.calculateVolumeGrowthScore(ctx, metrics.TradingVolume, metrics.PreviousTradingVolume)
	
	// Lending book performance: 20% weight
	lendingScore := pyc.calculateLendingScore(ctx, metrics.LendingVolume, metrics.DefaultRate)
	
	// DUSD revenue: 20% weight
	dusdScore := pyc.calculateDUSDScore(ctx, metrics.DUSDRevenue, metrics.DUSDVolume)
	
	// Calculate weighted average (0-8% range)
	weightedScore := profitabilityScore.MulInt64(40).
		Add(volumeGrowthScore.MulInt64(20)).
		Add(lendingScore.MulInt64(20)).
		Add(dusdScore.MulInt64(20)).
		QuoInt64(100)
	
	// Convert to yield percentage (0-8%)
	maxYield := sdk.NewDecWithPrec(8, 2) // 8%
	currentYield := weightedScore.Mul(maxYield)
	
	return currentYield, nil
}

// GetPerformanceMetrics retrieves current platform performance metrics
func (pyc *PerformanceYieldCalculator) GetPerformanceMetrics(ctx sdk.Context) types.PerformanceMetrics {
	// This would integrate with various modules to get real metrics
	// For now, return placeholder metrics
	return types.PerformanceMetrics{
		PlatformRevenue:       sdk.NewDec(1000000),   // ₹10 lakh revenue
		PlatformExpenses:      sdk.NewDec(800000),    // ₹8 lakh expenses
		TradingVolume:         sdk.NewDec(50000000),  // ₹5 Cr trading volume
		PreviousTradingVolume: sdk.NewDec(40000000),  // ₹4 Cr previous volume
		LendingVolume:         sdk.NewDec(10000000),  // ₹1 Cr lending
		DefaultRate:           sdk.NewDecWithPrec(2, 2), // 2% default rate
		DUSDRevenue:           sdk.NewDec(500000),    // ₹5 lakh DUSD revenue
		DUSDVolume:            sdk.NewDec(100000000), // ₹10 Cr DUSD volume
	}
}

// calculateProfitabilityScore scores platform profitability (0-1)
func (pyc *PerformanceYieldCalculator) calculateProfitabilityScore(ctx sdk.Context, revenue, expenses sdk.Dec) sdk.Dec {
	if revenue.IsZero() {
		return sdk.ZeroDec()
	}
	
	// Calculate profit margin
	profit := revenue.Sub(expenses)
	profitMargin := profit.Quo(revenue)
	
	// Score based on profit margin
	// 20%+ margin = 1.0 score
	// 10% margin = 0.5 score
	// 0% margin = 0.0 score
	// Negative = 0.0 score
	if profitMargin.IsNegative() {
		return sdk.ZeroDec()
	}
	
	targetMargin := sdk.NewDecWithPrec(20, 2) // 20%
	score := profitMargin.Quo(targetMargin)
	
	if score.GT(sdk.OneDec()) {
		score = sdk.OneDec()
	}
	
	return score
}

// calculateVolumeGrowthScore scores trading volume growth (0-1)
func (pyc *PerformanceYieldCalculator) calculateVolumeGrowthScore(ctx sdk.Context, currentVolume, previousVolume sdk.Dec) sdk.Dec {
	if previousVolume.IsZero() {
		return sdk.NewDecWithPrec(50, 2) // 0.5 default for new platforms
	}
	
	// Calculate growth rate
	growth := currentVolume.Sub(previousVolume).Quo(previousVolume)
	
	// Score based on growth
	// 50%+ growth = 1.0 score
	// 25% growth = 0.5 score
	// 0% growth = 0.0 score
	// Negative growth = 0.0 score
	if growth.IsNegative() {
		return sdk.ZeroDec()
	}
	
	targetGrowth := sdk.NewDecWithPrec(50, 2) // 50%
	score := growth.Quo(targetGrowth)
	
	if score.GT(sdk.OneDec()) {
		score = sdk.OneDec()
	}
	
	return score
}

// calculateLendingScore scores lending book performance (0-1)
func (pyc *PerformanceYieldCalculator) calculateLendingScore(ctx sdk.Context, lendingVolume sdk.Dec, defaultRate sdk.Dec) sdk.Dec {
	if lendingVolume.IsZero() {
		return sdk.ZeroDec()
	}
	
	// Score based on low default rate
	// 0% default = 1.0 score
	// 2% default = 0.5 score
	// 5%+ default = 0.0 score
	maxAcceptableDefault := sdk.NewDecWithPrec(5, 2) // 5%
	
	if defaultRate.GTE(maxAcceptableDefault) {
		return sdk.ZeroDec()
	}
	
	score := sdk.OneDec().Sub(defaultRate.Quo(maxAcceptableDefault))
	return score
}

// calculateDUSDScore scores DUSD performance (0-1)
func (pyc *PerformanceYieldCalculator) calculateDUSDScore(ctx sdk.Context, dusdRevenue, dusdVolume sdk.Dec) sdk.Dec {
	if dusdVolume.IsZero() {
		return sdk.ZeroDec()
	}
	
	// Calculate revenue efficiency
	revenueEfficiency := dusdRevenue.Quo(dusdVolume)
	
	// Score based on efficiency
	// 1%+ efficiency = 1.0 score
	// 0.5% efficiency = 0.5 score
	// 0% efficiency = 0.0 score
	targetEfficiency := sdk.NewDecWithPrec(1, 2) // 1%
	score := revenueEfficiency.Quo(targetEfficiency)
	
	if score.GT(sdk.OneDec()) {
		score = sdk.OneDec()
	}
	
	return score
}

// DistributePerformanceYield distributes yield to DINR holders based on performance
func (pyc *PerformanceYieldCalculator) DistributePerformanceYield(ctx sdk.Context) error {
	// Calculate current yield rate
	yieldRate, err := pyc.CalculateCurrentYield(ctx)
	if err != nil {
		return err
	}
	
	// Get total DINR supply
	totalSupply := pyc.keeper.GetTotalSupply(ctx)
	if totalSupply.IsZero() {
		return nil // No DINR to distribute to
	}
	
	// Calculate total yield amount
	annualYield := sdk.NewDecFromInt(totalSupply.Amount).Mul(yieldRate)
	// Convert to quarterly distribution
	quarterlyYield := annualYield.QuoInt64(4)
	
	// Store yield distribution record
	distribution := types.YieldDistribution{
		Timestamp:        ctx.BlockTime(),
		YieldRate:        yieldRate,
		TotalDistributed: quarterlyYield.TruncateInt(),
		TotalSupply:      totalSupply.Amount,
	}
	
	pyc.keeper.SetYieldDistribution(ctx, distribution)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeYieldDistributed,
			sdk.NewAttribute("yield_rate", yieldRate.String()),
			sdk.NewAttribute("total_distributed", quarterlyYield.String()),
			sdk.NewAttribute("timestamp", ctx.BlockTime().String()),
		),
	)
	
	pyc.keeper.Logger(ctx).Info("Performance yield distributed",
		"yield_rate", yieldRate.String(),
		"total_distributed", quarterlyYield.String(),
	)
	
	return nil
}

// GetHistoricalYield returns historical yield rates
func (pyc *PerformanceYieldCalculator) GetHistoricalYield(ctx sdk.Context, periods int) []types.YieldDistribution {
	// Retrieve last N yield distributions
	distributions := []types.YieldDistribution{}
	
	// This would query the store for historical distributions
	// For now, return empty slice
	return distributions
}

// ValidateYieldParameters ensures yield parameters are within acceptable ranges
func (pyc *PerformanceYieldCalculator) ValidateYieldParameters(ctx sdk.Context) error {
	params := pyc.keeper.GetParams(ctx)
	
	// Ensure yield is performance-based (no guaranteed minimum)
	if params.Fees.YieldRateMin > 0 {
		return fmt.Errorf("yield must be performance-based with 0% minimum")
	}
	
	// Ensure maximum yield is reasonable
	if params.Fees.YieldRateMax > 800 { // 8% in basis points
		return fmt.Errorf("maximum yield rate cannot exceed 8%")
	}
	
	return nil
}