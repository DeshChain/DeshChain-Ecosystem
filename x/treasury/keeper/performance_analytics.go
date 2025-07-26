package keeper

import (
	"fmt"
	"sort"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/treasury/types"
)

// PerformanceAnalytics handles treasury performance reporting and analytics
type PerformanceAnalytics struct {
	keeper Keeper
}

// NewPerformanceAnalytics creates a new performance analytics engine
func NewPerformanceAnalytics(keeper Keeper) *PerformanceAnalytics {
	return &PerformanceAnalytics{
		keeper: keeper,
	}
}

// TreasuryPerformanceReport represents comprehensive treasury performance analysis
type TreasuryPerformanceReport struct {
	ReportID            string                          `json:"report_id"`
	ReportPeriod        types.TimeRange                 `json:"report_period"`
	OverallPerformance  types.OverallPerformance        `json:"overall_performance"`
	PoolPerformance     []types.PoolPerformanceMetrics  `json:"pool_performance"`
	RevenueAnalysis     types.RevenueAnalysis           `json:"revenue_analysis"`
	AllocationAnalysis  types.AllocationAnalysis        `json:"allocation_analysis"`
	RebalanceAnalysis   types.RebalanceAnalysis         `json:"rebalance_analysis"`
	RiskMetrics         types.TreasuryRiskMetrics       `json:"risk_metrics"`
	EfficiencyMetrics   types.EfficiencyMetrics         `json:"efficiency_metrics"`
	BenchmarkComparison types.BenchmarkComparison       `json:"benchmark_comparison"`
	TrendAnalysis       types.TreasuryTrendAnalysis     `json:"trend_analysis"`
	Recommendations     []types.PerformanceRecommendation `json:"recommendations"`
	GeneratedAt         time.Time                       `json:"generated_at"`
}

// GeneratePerformanceReport generates comprehensive treasury performance report
func (pa *PerformanceAnalytics) GeneratePerformanceReport(ctx sdk.Context, timeRange types.TimeRange) (*TreasuryPerformanceReport, error) {
	report := &TreasuryPerformanceReport{
		ReportID:     pa.generateReportID(ctx),
		ReportPeriod: timeRange,
		GeneratedAt:  ctx.BlockTime(),
	}

	// Get all treasury pools
	pools := pa.keeper.GetAllTreasuryPools(ctx)

	// Generate overall performance metrics
	report.OverallPerformance = pa.calculateOverallPerformance(ctx, pools, timeRange)

	// Generate individual pool performance metrics
	for _, pool := range pools {
		poolMetrics := pa.calculatePoolPerformance(ctx, pool, timeRange)
		report.PoolPerformance = append(report.PoolPerformance, poolMetrics)
	}

	// Analyze revenue streams and sources
	report.RevenueAnalysis = pa.analyzeRevenueStreams(ctx, timeRange)

	// Analyze allocation efficiency
	report.AllocationAnalysis = pa.analyzeAllocationEfficiency(ctx, pools, timeRange)

	// Analyze rebalancing performance
	report.RebalanceAnalysis = pa.analyzeRebalancePerformance(ctx, timeRange)

	// Calculate risk metrics
	report.RiskMetrics = pa.calculateRiskMetrics(ctx, pools, timeRange)

	// Calculate efficiency metrics
	report.EfficiencyMetrics = pa.calculateEfficiencyMetrics(ctx, pools, timeRange)

	// Compare against benchmarks
	report.BenchmarkComparison = pa.compareToBenchmarks(ctx, report)

	// Analyze trends
	report.TrendAnalysis = pa.analyzeTrends(ctx, pools, timeRange)

	// Generate recommendations
	report.Recommendations = pa.generateRecommendations(ctx, report)

	// Store performance report
	pa.keeper.SetTreasuryPerformanceReport(ctx, *report)

	// Emit performance report event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTreasuryPerformanceReport,
			sdk.NewAttribute(types.AttributeKeyReportID, report.ReportID),
			sdk.NewAttribute(types.AttributeKeyReportPeriod, fmt.Sprintf("%v-%v", timeRange.StartTime.Unix(), timeRange.EndTime.Unix())),
		),
	)

	return report, nil
}

// calculateOverallPerformance calculates overall treasury performance metrics
func (pa *PerformanceAnalytics) calculateOverallPerformance(ctx sdk.Context, pools []TreasuryPool, timeRange types.TimeRange) types.OverallPerformance {
	performance := types.OverallPerformance{
		TimeRange: timeRange,
	}

	// Calculate total treasury value at start and end of period
	startValue := pa.getTreasuryValueAtTime(ctx, timeRange.StartTime)
	endValue := pa.getTreasuryValueAtTime(ctx, timeRange.EndTime)
	
	performance.StartingValue = startValue
	performance.EndingValue = endValue

	// Calculate absolute and percentage growth
	performance.AbsoluteGrowth = endValue.Sub(startValue)
	if startValue.IsPositive() {
		growth := performance.AbsoluteGrowth.AmountOf("namo").ToDec().Quo(startValue.AmountOf("namo").ToDec())
		performance.PercentageGrowth = growth
	}

	// Calculate total revenue received during period
	totalRevenue := pa.getTotalRevenueInPeriod(ctx, timeRange)
	performance.TotalRevenue = totalRevenue

	// Calculate total expenses during period
	totalExpenses := pa.getTotalExpensesInPeriod(ctx, timeRange)
	performance.TotalExpenses = totalExpenses

	// Calculate net income
	performance.NetIncome = totalRevenue.Sub(totalExpenses)

	// Calculate revenue efficiency (net income / total revenue)
	if totalRevenue.IsPositive() {
		efficiency := performance.NetIncome.AmountOf("namo").ToDec().Quo(totalRevenue.AmountOf("namo").ToDec())
		performance.RevenueEfficiency = efficiency
	}

	// Calculate rebalancing frequency and cost
	rebalances := pa.keeper.GetRebalancesInPeriod(ctx, timeRange.StartTime, timeRange.EndTime)
	performance.RebalanceCount = int64(len(rebalances))
	
	var totalRebalanceCost sdk.Coins
	for _, rebalance := range rebalances {
		totalRebalanceCost = totalRebalanceCost.Add(pa.calculateRebalanceCost(rebalance)...)
	}
	performance.RebalanceCost = totalRebalanceCost

	// Calculate allocation adherence score
	performance.AllocationAdherence = pa.calculateAllocationAdherence(ctx, pools)

	// Calculate volatility metrics
	performance.Volatility = pa.calculateTreasuryVolatility(ctx, timeRange)

	// Calculate Sharpe ratio equivalent for treasury
	performance.RiskAdjustedReturn = pa.calculateRiskAdjustedReturn(performance.PercentageGrowth, performance.Volatility)

	return performance
}

// calculatePoolPerformance calculates performance metrics for individual pool
func (pa *PerformanceAnalytics) calculatePoolPerformance(ctx sdk.Context, pool TreasuryPool, timeRange types.TimeRange) types.PoolPerformanceMetrics {
	metrics := types.PoolPerformanceMetrics{
		PoolID:    pool.PoolID,
		PoolName:  pool.PoolName,
		PoolType:  pool.PoolType,
		TimeRange: timeRange,
	}

	// Get pool balance history
	history := pa.keeper.GetPoolBalanceHistory(ctx, pool.PoolID, timeRange.StartTime, timeRange.EndTime)
	
	if len(history) > 0 {
		metrics.StartingBalance = history[0].Balance
		metrics.EndingBalance = history[len(history)-1].Balance
		
		// Calculate absolute and percentage change
		metrics.AbsoluteChange = metrics.EndingBalance.Sub(metrics.StartingBalance)
		if metrics.StartingBalance.IsPositive() {
			change := metrics.AbsoluteChange.AmountOf("namo").ToDec().Quo(metrics.StartingBalance.AmountOf("namo").ToDec())
			metrics.PercentageChange = change
		}
	}

	// Calculate pool-specific metrics
	metrics.TotalInflows = pa.getPoolInflowsInPeriod(ctx, pool.PoolID, timeRange)
	metrics.TotalOutflows = pa.getPoolOutflowsInPeriod(ctx, pool.PoolID, timeRange)
	metrics.NetFlow = metrics.TotalInflows.Sub(metrics.TotalOutflows)

	// Calculate utilization rate
	metrics.UtilizationRate = pa.calculatePoolUtilization(ctx, pool, timeRange)

	// Calculate allocation target adherence
	metrics.AllocationAdherence = pa.calculatePoolAllocationAdherence(ctx, pool, timeRange)

	// Calculate rebalancing impact on this pool
	metrics.RebalanceImpact = pa.calculatePoolRebalanceImpact(ctx, pool.PoolID, timeRange)

	// Calculate efficiency metrics
	metrics.CostEfficiency = pa.calculatePoolCostEfficiency(ctx, pool, timeRange)
	metrics.PurposeAlignment = pa.calculatePurposeAlignment(ctx, pool, timeRange)

	// Calculate risk metrics for pool
	metrics.VolatilityScore = pa.calculatePoolVolatility(ctx, pool.PoolID, timeRange)
	metrics.LiquidityRisk = pa.calculatePoolLiquidityRisk(ctx, pool)

	// Performance ranking among similar pools
	metrics.PerformanceRanking = pa.calculatePoolRanking(ctx, pool, timeRange)

	return metrics
}

// analyzeRevenueStreams analyzes different revenue sources and their performance
func (pa *PerformanceAnalytics) analyzeRevenueStreams(ctx sdk.Context, timeRange types.TimeRange) types.RevenueAnalysis {
	analysis := types.RevenueAnalysis{
		TimeRange: timeRange,
	}

	// Get all revenue transactions in period
	revenueTransactions := pa.keeper.GetRevenueTransactionsInPeriod(ctx, timeRange.StartTime, timeRange.EndTime)

	// Categorize revenue by source
	revenueBySource := make(map[string]sdk.Coins)
	for _, tx := range revenueTransactions {
		if existing, found := revenueBySource[tx.Source]; found {
			revenueBySource[tx.Source] = existing.Add(tx.Amount...)
		} else {
			revenueBySource[tx.Source] = tx.Amount
		}
	}

	// Calculate total revenue
	var totalRevenue sdk.Coins
	for _, amount := range revenueBySource {
		totalRevenue = totalRevenue.Add(amount...)
	}
	analysis.TotalRevenue = totalRevenue

	// Create revenue source breakdown
	for source, amount := range revenueBySource {
		percentage := amount.AmountOf("namo").ToDec().Quo(totalRevenue.AmountOf("namo").ToDec())
		
		sourceAnalysis := types.RevenueSourceAnalysis{
			Source:     source,
			Amount:     amount,
			Percentage: percentage,
			Growth:     pa.calculateRevenueSourceGrowth(ctx, source, timeRange),
			Stability:  pa.calculateRevenueSourceStability(ctx, source, timeRange),
		}
		
		analysis.RevenueBySource = append(analysis.RevenueBySource, sourceAnalysis)
	}

	// Sort by amount (descending)
	sort.Slice(analysis.RevenueBySource, func(i, j int) bool {
		return analysis.RevenueBySource[i].Amount.AmountOf("namo").GT(analysis.RevenueBySource[j].Amount.AmountOf("namo"))
	})

	// Calculate revenue diversity score
	analysis.DiversityScore = pa.calculateRevenueDiversityScore(analysis.RevenueBySource)

	// Calculate revenue predictability
	analysis.PredictabilityScore = pa.calculateRevenuePredictability(ctx, timeRange)

	// Identify growth opportunities
	analysis.GrowthOpportunities = pa.identifyRevenueGrowthOpportunities(ctx, analysis)

	return analysis
}

// analyzeAllocationEfficiency analyzes how efficiently funds are allocated
func (pa *PerformanceAnalytics) analyzeAllocationEfficiency(ctx sdk.Context, pools []TreasuryPool, timeRange types.TimeRange) types.AllocationAnalysis {
	analysis := types.AllocationAnalysis{
		TimeRange: timeRange,
	}

	// Calculate target vs actual allocations
	totalValue := pa.calculateTotalTreasuryValue(pools)
	
	var totalDeviation sdk.Dec
	for _, pool := range pools {
		currentPercentage := pool.Balance.AmountOf("namo").ToDec().Quo(totalValue.AmountOf("namo").ToDec())
		targetPercentage := pool.Allocation.TargetPercentage
		deviation := currentPercentage.Sub(targetPercentage).Abs()
		
		poolAllocation := types.PoolAllocationAnalysis{
			PoolID:           pool.PoolID,
			TargetAllocation: targetPercentage,
			ActualAllocation: currentPercentage,
			Deviation:        deviation,
			PerformanceScore: pa.calculateAllocationPerformanceScore(ctx, pool, timeRange),
		}
		
		analysis.PoolAllocations = append(analysis.PoolAllocations, poolAllocation)
		totalDeviation = totalDeviation.Add(deviation)
	}

	// Calculate overall allocation efficiency
	if len(pools) > 0 {
		averageDeviation := totalDeviation.QuoInt64(int64(len(pools)))
		analysis.AllocationEfficiency = sdk.OneDec().Sub(averageDeviation)
		if analysis.AllocationEfficiency.LT(sdk.ZeroDec()) {
			analysis.AllocationEfficiency = sdk.ZeroDec()
		}
	}

	// Calculate rebalancing effectiveness
	analysis.RebalancingEffectiveness = pa.calculateRebalancingEffectiveness(ctx, timeRange)

	// Identify misallocated pools
	analysis.MisallocatedPools = pa.identifyMisallocatedPools(analysis.PoolAllocations)

	// Calculate allocation stability
	analysis.AllocationStability = pa.calculateAllocationStability(ctx, pools, timeRange)

	return analysis
}

// calculateRiskMetrics calculates comprehensive risk metrics for treasury
func (pa *PerformanceAnalytics) calculateRiskMetrics(ctx sdk.Context, pools []TreasuryPool, timeRange types.TimeRange) types.TreasuryRiskMetrics {
	metrics := types.TreasuryRiskMetrics{
		TimeRange: timeRange,
	}

	// Calculate concentration risk
	metrics.ConcentrationRisk = pa.calculateConcentrationRisk(pools)

	// Calculate liquidity risk
	metrics.LiquidityRisk = pa.calculateOverallLiquidityRisk(ctx, pools)

	// Calculate operational risk
	metrics.OperationalRisk = pa.calculateOperationalRisk(ctx, timeRange)

	// Calculate market risk
	metrics.MarketRisk = pa.calculateMarketRisk(ctx, timeRange)

	// Calculate governance risk
	metrics.GovernanceRisk = pa.calculateGovernanceRisk(ctx, timeRange)

	// Calculate Value at Risk (VaR) equivalent
	metrics.ValueAtRisk = pa.calculateTreasuryVaR(ctx, pools, timeRange)

	// Calculate stress test scenarios
	metrics.StressTestResults = pa.performStressTests(ctx, pools)

	// Calculate overall risk score
	metrics.OverallRiskScore = pa.calculateOverallRiskScore(metrics)

	// Risk level classification
	metrics.RiskLevel = pa.classifyRiskLevel(metrics.OverallRiskScore)

	return metrics
}

// generateRecommendations generates actionable recommendations based on performance analysis
func (pa *PerformanceAnalytics) generateRecommendations(ctx sdk.Context, report *TreasuryPerformanceReport) []types.PerformanceRecommendation {
	var recommendations []types.PerformanceRecommendation

	// Analyze allocation efficiency recommendations
	if report.AllocationAnalysis.AllocationEfficiency.LT(sdk.NewDecWithPrec(9, 1)) { // < 90% efficiency
		recommendation := types.PerformanceRecommendation{
			Type:        "ALLOCATION_OPTIMIZATION",
			Priority:    "HIGH",
			Title:       "Improve Allocation Efficiency",
			Description: "Current allocation efficiency is below optimal levels",
			Action:      "Implement more frequent rebalancing or adjust target allocations",
			ExpectedImpact: "5-10% improvement in allocation efficiency",
			Timeline:    "1-2 weeks",
		}
		recommendations = append(recommendations, recommendation)
	}

	// Analyze rebalancing frequency recommendations
	if report.OverallPerformance.RebalanceCount > 20 { // Too frequent rebalancing
		recommendation := types.PerformanceRecommendation{
			Type:        "REBALANCING_OPTIMIZATION",
			Priority:    "MEDIUM",
			Title:       "Reduce Rebalancing Frequency",
			Description: "Excessive rebalancing may be increasing costs",
			Action:      "Increase rebalancing thresholds or extend minimum gaps",
			ExpectedImpact: "10-15% reduction in rebalancing costs",
			Timeline:    "Immediate",
		}
		recommendations = append(recommendations, recommendation)
	}

	// Analyze risk recommendations
	if report.RiskMetrics.OverallRiskScore.GT(sdk.NewDecWithPrec(7, 1)) { // > 70% risk score
		recommendation := types.PerformanceRecommendation{
			Type:        "RISK_MANAGEMENT",
			Priority:    "HIGH",
			Title:       "Reduce Treasury Risk Exposure",
			Description: "Overall risk score indicates elevated risk levels",
			Action:      "Increase reserve pool allocation and implement additional safeguards",
			ExpectedImpact: "20-30% reduction in risk exposure",
			Timeline:    "2-4 weeks",
		}
		recommendations = append(recommendations, recommendation)
	}

	// Analyze revenue diversification recommendations
	if report.RevenueAnalysis.DiversityScore.LT(sdk.NewDecWithPrec(6, 1)) { // < 60% diversity
		recommendation := types.PerformanceRecommendation{
			Type:        "REVENUE_DIVERSIFICATION",
			Priority:    "MEDIUM",
			Title:       "Diversify Revenue Sources",
			Description: "Revenue concentration risk identified",
			Action:      "Develop new revenue streams and reduce dependency on top sources",
			ExpectedImpact: "Improved revenue stability and growth potential",
			Timeline:    "3-6 months",
		}
		recommendations = append(recommendations, recommendation)
	}

	// Analyze pool-specific recommendations
	for _, poolMetrics := range report.PoolPerformance {
		if poolMetrics.UtilizationRate.LT(sdk.NewDecWithPrec(3, 1)) { // < 30% utilization
			recommendation := types.PerformanceRecommendation{
				Type:        "POOL_OPTIMIZATION",
				Priority:    "LOW",
				Title:       fmt.Sprintf("Optimize %s Pool Utilization", poolMetrics.PoolName),
				Description: fmt.Sprintf("Pool %s has low utilization rate", poolMetrics.PoolID),
				Action:      "Consider reducing allocation or finding new use cases",
				ExpectedImpact: "Better capital efficiency",
				Timeline:    "1-3 months",
			}
			recommendations = append(recommendations, recommendation)
		}
	}

	// Sort recommendations by priority
	sort.Slice(recommendations, func(i, j int) bool {
		priorityOrder := map[string]int{"HIGH": 1, "MEDIUM": 2, "LOW": 3}
		return priorityOrder[recommendations[i].Priority] < priorityOrder[recommendations[j].Priority]
	})

	return recommendations
}

// Helper utility functions
func (pa *PerformanceAnalytics) generateReportID(ctx sdk.Context) string {
	timestamp := ctx.BlockTime().Unix()
	return fmt.Sprintf("PERF_RPT_%d_%d", ctx.BlockHeight(), timestamp)
}

func (pa *PerformanceAnalytics) calculateTotalTreasuryValue(pools []TreasuryPool) sdk.Coins {
	var total sdk.Coins
	for _, pool := range pools {
		total = total.Add(pool.Balance...)
	}
	return total
}

func (pa *PerformanceAnalytics) calculateRevenueDiversityScore(revenueBySource []types.RevenueSourceAnalysis) sdk.Dec {
	if len(revenueBySource) <= 1 {
		return sdk.ZeroDec()
	}

	// Calculate Herfindahl-Hirschman Index equivalent
	var hhi sdk.Dec
	for _, source := range revenueBySource {
		hhi = hhi.Add(source.Percentage.Mul(source.Percentage))
	}

	// Convert to diversity score (inverse of HHI)
	diversityScore := sdk.OneDec().Sub(hhi)
	if diversityScore.LT(sdk.ZeroDec()) {
		diversityScore = sdk.ZeroDec()
	}

	return diversityScore
}

func (pa *PerformanceAnalytics) classifyRiskLevel(riskScore sdk.Dec) string {
	if riskScore.LT(sdk.NewDecWithPrec(3, 1)) { // < 30%
		return "LOW"
	} else if riskScore.LT(sdk.NewDecWithPrec(6, 1)) { // < 60%
		return "MEDIUM"
	} else if riskScore.LT(sdk.NewDecWithPrec(8, 1)) { // < 80%
		return "HIGH"
	} else {
		return "CRITICAL"
	}
}

// Additional helper methods would include all calculation functions
// referenced in the analysis methods above