package keeper

import (
	"fmt"
	"sort"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/explorer/types"
)

// AnalyticsEngine handles advanced blockchain analytics and reporting
type AnalyticsEngine struct {
	keeper Keeper
}

// NewAnalyticsEngine creates a new analytics engine
func NewAnalyticsEngine(keeper Keeper) *AnalyticsEngine {
	return &AnalyticsEngine{
		keeper: keeper,
	}
}

// ChainAnalytics represents comprehensive chain analytics
type ChainAnalytics struct {
	AnalyticsID       string                    `json:"analytics_id"`
	TimeRange         types.TimeRange           `json:"time_range"`
	NetworkStats      types.NetworkStatistics   `json:"network_stats"`
	TransactionStats  types.TransactionStatistics `json:"transaction_stats"`
	ModuleActivity    []types.ModuleAnalytics   `json:"module_activity"`
	ValidatorStats    types.ValidatorStatistics `json:"validator_stats"`
	TokenEconomics    types.TokenEconomics      `json:"token_economics"`
	GovernanceStats   types.GovernanceStatistics `json:"governance_stats"`
	CulturalMetrics   types.CulturalMetrics     `json:"cultural_metrics"`
	LendingMetrics    types.LendingMetrics      `json:"lending_metrics"`
	PerformanceMetrics types.PerformanceMetrics `json:"performance_metrics"`
	TrendAnalysis     types.TrendAnalysis       `json:"trend_analysis"`
	PredictiveMetrics types.PredictiveMetrics   `json:"predictive_metrics"`
	GeneratedAt       time.Time                 `json:"generated_at"`
}

// GenerateChainAnalytics generates comprehensive chain analytics for time period
func (ae *AnalyticsEngine) GenerateChainAnalytics(ctx sdk.Context, startTime, endTime time.Time) (*ChainAnalytics, error) {
	analytics := &ChainAnalytics{
		AnalyticsID: ae.generateAnalyticsID(ctx),
		TimeRange: types.TimeRange{
			StartTime: startTime,
			EndTime:   endTime,
		},
		GeneratedAt: ctx.BlockTime(),
	}

	// Generate network statistics
	analytics.NetworkStats = ae.generateNetworkStatistics(ctx, startTime, endTime)

	// Generate transaction statistics
	analytics.TransactionStats = ae.generateTransactionStatistics(ctx, startTime, endTime)

	// Generate module activity analysis
	analytics.ModuleActivity = ae.generateModuleAnalytics(ctx, startTime, endTime)

	// Generate validator statistics
	analytics.ValidatorStats = ae.generateValidatorStatistics(ctx, startTime, endTime)

	// Generate token economics analysis
	analytics.TokenEconomics = ae.generateTokenEconomics(ctx, startTime, endTime)

	// Generate governance statistics
	analytics.GovernanceStats = ae.generateGovernanceStatistics(ctx, startTime, endTime)

	// Generate cultural metrics
	analytics.CulturalMetrics = ae.generateCulturalMetrics(ctx, startTime, endTime)

	// Generate lending metrics
	analytics.LendingMetrics = ae.generateLendingMetrics(ctx, startTime, endTime)

	// Generate performance metrics
	analytics.PerformanceMetrics = ae.generatePerformanceMetrics(ctx, startTime, endTime)

	// Generate trend analysis
	analytics.TrendAnalysis = ae.generateTrendAnalysis(ctx, startTime, endTime)

	// Generate predictive metrics
	analytics.PredictiveMetrics = ae.generatePredictiveMetrics(ctx, analytics)

	// Store analytics report
	ae.keeper.SetChainAnalytics(ctx, *analytics)

	// Emit analytics generation event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAnalyticsGenerated,
			sdk.NewAttribute(types.AttributeKeyAnalyticsID, analytics.AnalyticsID),
			sdk.NewAttribute(types.AttributeKeyTimeRange, fmt.Sprintf("%v-%v", startTime.Unix(), endTime.Unix())),
		),
	)

	return analytics, nil
}

// generateNetworkStatistics generates comprehensive network statistics
func (ae *AnalyticsEngine) generateNetworkStatistics(ctx sdk.Context, startTime, endTime time.Time) types.NetworkStatistics {
	stats := types.NetworkStatistics{
		TimeRange: types.TimeRange{StartTime: startTime, EndTime: endTime},
	}

	// Get blocks in time range
	blocks := ae.keeper.GetBlocksInTimeRange(ctx, startTime, endTime)
	stats.TotalBlocks = int64(len(blocks))

	if len(blocks) == 0 {
		return stats
	}

	// Calculate block timing statistics
	var blockTimes []time.Duration
	for i := 1; i < len(blocks); i++ {
		timeDiff := blocks[i].Time.Sub(blocks[i-1].Time)
		blockTimes = append(blockTimes, timeDiff)
	}

	if len(blockTimes) > 0 {
		sort.Slice(blockTimes, func(i, j int) bool {
			return blockTimes[i] < blockTimes[j]
		})

		// Average block time
		var totalTime time.Duration
		for _, t := range blockTimes {
			totalTime += t
		}
		stats.AverageBlockTime = totalTime / time.Duration(len(blockTimes))

		// Median block time
		mid := len(blockTimes) / 2
		if len(blockTimes)%2 == 0 {
			stats.MedianBlockTime = (blockTimes[mid-1] + blockTimes[mid]) / 2
		} else {
			stats.MedianBlockTime = blockTimes[mid]
		}

		// Min/Max block times
		stats.MinBlockTime = blockTimes[0]
		stats.MaxBlockTime = blockTimes[len(blockTimes)-1]
	}

	// Network health indicators
	stats.ChainStability = ae.calculateChainStability(blockTimes)
	stats.NetworkLatency = ae.calculateNetworkLatency(ctx, blocks)
	stats.SyncStatus = ae.getNetworkSyncStatus(ctx)

	// Uptime statistics
	stats.NetworkUptime = ae.calculateNetworkUptime(ctx, startTime, endTime)
	stats.ActiveValidators = ae.getActiveValidatorCount(ctx)

	return stats
}

// generateTransactionStatistics generates comprehensive transaction statistics
func (ae *AnalyticsEngine) generateTransactionStatistics(ctx sdk.Context, startTime, endTime time.Time) types.TransactionStatistics {
	stats := types.TransactionStatistics{
		TimeRange: types.TimeRange{StartTime: startTime, EndTime: endTime},
	}

	// Get transactions in time range
	txs := ae.keeper.GetTransactionsInTimeRange(ctx, startTime, endTime)
	stats.TotalTransactions = int64(len(txs))

	if len(txs) == 0 {
		return stats
	}

	// Success/failure analysis
	var successCount, failureCount int64
	var totalGasUsed, totalGasWanted uint64
	var totalFees sdk.Int = sdk.ZeroInt()
	var gasPrices []sdk.Dec

	for _, tx := range txs {
		if tx.Code == 0 {
			successCount++
		} else {
			failureCount++
		}

		totalGasUsed += tx.GasUsed
		totalGasWanted += tx.GasWanted

		// Calculate fees
		for _, fee := range tx.Fees {
			if fee.Denom == "namo" {
				totalFees = totalFees.Add(fee.Amount)
			}
		}

		// Calculate gas price
		if tx.GasWanted > 0 {
			gasPrice := totalFees.ToDec().QuoInt64(int64(tx.GasWanted))
			gasPrices = append(gasPrices, gasPrice)
		}
	}

	stats.SuccessfulTransactions = successCount
	stats.FailedTransactions = failureCount
	stats.SuccessRate = sdk.NewDec(successCount).QuoInt64(int64(len(txs)))

	// Gas statistics
	stats.TotalGasUsed = totalGasUsed
	stats.TotalGasWanted = totalGasWanted
	stats.AverageGasUsed = totalGasUsed / uint64(len(txs))
	stats.GasEfficiency = sdk.NewDec(int64(totalGasUsed)).QuoInt64(int64(totalGasWanted))

	// Fee statistics
	stats.TotalFees = sdk.NewCoin("namo", totalFees)
	stats.AverageFee = sdk.NewCoin("namo", totalFees.QuoRaw(int64(len(txs))))

	// Gas price statistics
	if len(gasPrices) > 0 {
		sort.Slice(gasPrices, func(i, j int) bool {
			return gasPrices[i].LT(gasPrices[j])
		})

		// Calculate average gas price
		var totalGasPrice sdk.Dec
		for _, price := range gasPrices {
			totalGasPrice = totalGasPrice.Add(price)
		}
		stats.AverageGasPrice = totalGasPrice.QuoInt64(int64(len(gasPrices)))

		// Median gas price
		mid := len(gasPrices) / 2
		if len(gasPrices)%2 == 0 {
			stats.MedianGasPrice = gasPrices[mid-1].Add(gasPrices[mid]).QuoInt64(2)
		} else {
			stats.MedianGasPrice = gasPrices[mid]
		}

		stats.MinGasPrice = gasPrices[0]
		stats.MaxGasPrice = gasPrices[len(gasPrices)-1]
	}

	// Transaction type distribution
	stats.TransactionTypes = ae.analyzeTransactionTypes(txs)

	// Transaction size analysis
	stats.AverageTransactionSize = ae.calculateAverageTransactionSize(txs)

	// TPS calculations
	timeDuration := endTime.Sub(startTime)
	if timeDuration > 0 {
		stats.TransactionsPerSecond = float64(len(txs)) / timeDuration.Seconds()
		stats.PeakTPS = ae.calculatePeakTPS(ctx, txs, timeDuration)
	}

	return stats
}

// generateModuleAnalytics generates activity analytics for all modules
func (ae *AnalyticsEngine) generateModuleAnalytics(ctx sdk.Context, startTime, endTime time.Time) []types.ModuleAnalytics {
	modules := []string{
		"bank", "staking", "gov", "cultural", "namo", "explorer",
		"krishimitra", "vyavasayamitra", "shikshaamitra", "dinr", "oracle",
	}

	var analytics []types.ModuleAnalytics

	for _, moduleName := range modules {
		moduleAnalytics := ae.generateSingleModuleAnalytics(ctx, moduleName, startTime, endTime)
		analytics = append(analytics, moduleAnalytics)
	}

	return analytics
}

// generateSingleModuleAnalytics generates analytics for a specific module
func (ae *AnalyticsEngine) generateSingleModuleAnalytics(ctx sdk.Context, moduleName string, startTime, endTime time.Time) types.ModuleAnalytics {
	analytics := types.ModuleAnalytics{
		ModuleName: moduleName,
		TimeRange:  types.TimeRange{StartTime: startTime, EndTime: endTime},
	}

	// Get module-specific transactions
	txs := ae.keeper.GetModuleTransactions(ctx, moduleName, startTime, endTime)
	analytics.TotalTransactions = int64(len(txs))

	if len(txs) == 0 {
		return analytics
	}

	// Calculate module-specific metrics
	var totalGasUsed uint64
	var totalFees sdk.Int = sdk.ZeroInt()
	var successCount int64

	for _, tx := range txs {
		totalGasUsed += tx.GasUsed
		for _, fee := range tx.Fees {
			if fee.Denom == "namo" {
				totalFees = totalFees.Add(fee.Amount)
			}
		}
		if tx.Code == 0 {
			successCount++
		}
	}

	analytics.TotalGasUsed = totalGasUsed
	analytics.TotalFees = sdk.NewCoin("namo", totalFees)
	analytics.SuccessRate = sdk.NewDec(successCount).QuoInt64(int64(len(txs)))

	// Module-specific features
	switch moduleName {
	case "cultural":
		analytics.CulturalEvents = ae.getCulturalEventCount(ctx, startTime, endTime)
		analytics.FestivalParticipation = ae.getFestivalParticipation(ctx, startTime, endTime)
	case "krishimitra":
		analytics.LoansProcessed = ae.getKrishiLoanCount(ctx, startTime, endTime)
		analytics.TotalLoanAmount = ae.getKrishiTotalLoanAmount(ctx, startTime, endTime)
	case "vyavasayamitra":
		analytics.BusinessLoansProcessed = ae.getBusinessLoanCount(ctx, startTime, endTime)
		analytics.TotalBusinessLoanAmount = ae.getBusinessTotalLoanAmount(ctx, startTime, endTime)
	case "shikshaamitra":
		analytics.EducationLoansProcessed = ae.getEducationLoanCount(ctx, startTime, endTime)
		analytics.TotalEducationLoanAmount = ae.getEducationTotalLoanAmount(ctx, startTime, endTime)
	}

	// Activity trends
	analytics.ActivityTrends = ae.calculateModuleActivityTrends(ctx, moduleName, startTime, endTime)

	return analytics
}

// generateValidatorStatistics generates comprehensive validator statistics
func (ae *AnalyticsEngine) generateValidatorStatistics(ctx sdk.Context, startTime, endTime time.Time) types.ValidatorStatistics {
	stats := types.ValidatorStatistics{
		TimeRange: types.TimeRange{StartTime: startTime, EndTime: endTime},
	}

	// Get all validators
	validators := ae.keeper.GetAllValidators(ctx)
	stats.TotalValidators = int64(len(validators))

	var activeValidators, jailedValidators int64
	var totalDelegations, totalSelfDelegations sdk.Int = sdk.ZeroInt(), sdk.ZeroInt()
	var commissionRates []sdk.Dec

	for _, validator := range validators {
		if validator.IsBonded() {
			activeValidators++
		}
		if validator.IsJailed() {
			jailedValidators++
		}

		// Add delegations
		totalDelegations = totalDelegations.Add(validator.Tokens)
		totalSelfDelegations = totalSelfDelegations.Add(validator.DelegatorShares.TruncateInt())

		// Commission rates
		commissionRates = append(commissionRates, validator.Commission.Rate)
	}

	stats.ActiveValidators = activeValidators
	stats.JailedValidators = jailedValidators
	stats.TotalDelegations = sdk.NewCoin("namo", totalDelegations)
	stats.TotalSelfDelegations = sdk.NewCoin("namo", totalSelfDelegations)

	// Commission statistics
	if len(commissionRates) > 0 {
		sort.Slice(commissionRates, func(i, j int) bool {
			return commissionRates[i].LT(commissionRates[j])
		})

		var totalCommission sdk.Dec
		for _, rate := range commissionRates {
			totalCommission = totalCommission.Add(rate)
		}
		stats.AverageCommission = totalCommission.QuoInt64(int64(len(commissionRates)))

		mid := len(commissionRates) / 2
		if len(commissionRates)%2 == 0 {
			stats.MedianCommission = commissionRates[mid-1].Add(commissionRates[mid]).QuoInt64(2)
		} else {
			stats.MedianCommission = commissionRates[mid]
		}
	}

	// Block production statistics
	blocks := ae.keeper.GetBlocksInTimeRange(ctx, startTime, endTime)
	proposerCounts := make(map[string]int64)
	
	for _, block := range blocks {
		proposerCounts[block.ProposerAddress]++
	}

	stats.BlockProduction = proposerCounts

	// Validator performance metrics
	stats.ValidatorPerformance = ae.calculateValidatorPerformance(ctx, validators, startTime, endTime)

	return stats
}

// generateTokenEconomics generates token economics analysis
func (ae *AnalyticsEngine) generateTokenEconomics(ctx sdk.Context, startTime, endTime time.Time) types.TokenEconomics {
	economics := types.TokenEconomics{
		TimeRange: types.TimeRange{StartTime: startTime, EndTime: endTime},
	}

	// Total supply
	economics.TotalSupply = ae.keeper.GetTotalSupply(ctx, "namo")

	// Circulating supply (total - locked/staked)
	stakedAmount := ae.keeper.GetTotalStakedAmount(ctx)
	economics.CirculatingSupply = economics.TotalSupply.Sub(sdk.NewCoin("namo", stakedAmount))

	// Staking statistics
	economics.TotalStaked = sdk.NewCoin("namo", stakedAmount)
	economics.StakingRatio = stakedAmount.ToDec().Quo(economics.TotalSupply.Amount.ToDec())

	// Transaction volumes
	txs := ae.keeper.GetTransactionsInTimeRange(ctx, startTime, endTime)
	var totalVolume sdk.Int = sdk.ZeroInt()
	
	for _, tx := range txs {
		volume := ae.calculateTransactionVolume(tx)
		totalVolume = totalVolume.Add(volume)
	}
	economics.TransactionVolume = sdk.NewCoin("namo", totalVolume)

	// Fee burn analysis
	var totalFeesBurned sdk.Int = sdk.ZeroInt()
	for _, tx := range txs {
		for _, fee := range tx.Fees {
			if fee.Denom == "namo" {
				// Assuming 30% of fees are burned (from tax system)
				burnAmount := fee.Amount.MulRaw(30).QuoRaw(100)
				totalFeesBurned = totalFeesBurned.Add(burnAmount)
			}
		}
	}
	economics.FeesBurned = sdk.NewCoin("namo", totalFeesBurned)

	// Inflation and deflation
	economics.InflationRate = ae.keeper.GetInflationRate(ctx)
	economics.EffectiveInflation = economics.InflationRate.Sub(totalFeesBurned.ToDec().Quo(economics.TotalSupply.Amount.ToDec()))

	// Yield and rewards
	economics.StakingYield = ae.calculateStakingYield(ctx, startTime, endTime)
	economics.TotalRewardsDistributed = ae.getTotalRewardsDistributed(ctx, startTime, endTime)

	return economics
}

// generateCulturalMetrics generates cultural blockchain metrics
func (ae *AnalyticsEngine) generateCulturalMetrics(ctx sdk.Context, startTime, endTime time.Time) types.CulturalMetrics {
	metrics := types.CulturalMetrics{
		TimeRange: types.TimeRange{StartTime: startTime, EndTime: endTime},
	}

	// Cultural transactions
	culturalTxs := ae.keeper.GetCulturalTransactions(ctx, startTime, endTime)
	metrics.CulturalTransactions = int64(len(culturalTxs))

	// Festival participation
	metrics.FestivalParticipation = ae.getFestivalParticipationMetrics(ctx, startTime, endTime)

	// Language diversity
	metrics.LanguageDistribution = ae.getLanguageDistribution(ctx, startTime, endTime)

	// Regional activity
	metrics.RegionalActivity = ae.getRegionalActivity(ctx, startTime, endTime)

	// Cultural content creation
	metrics.CulturalContentCreated = ae.getCulturalContentCreated(ctx, startTime, endTime)

	// Heritage preservation
	metrics.HeritageItemsPreserved = ae.getHeritageItemsPreserved(ctx, startTime, endTime)

	// Community engagement
	metrics.CommunityEngagement = ae.getCommunityEngagementScore(ctx, startTime, endTime)

	return metrics
}

// generateLendingMetrics generates lending ecosystem metrics
func (ae *AnalyticsEngine) generateLendingMetrics(ctx sdk.Context, startTime, endTime time.Time) types.LendingMetrics {
	metrics := types.LendingMetrics{
		TimeRange: types.TimeRange{StartTime: startTime, EndTime: endTime},
	}

	// Agricultural lending (KrishiMitra)
	metrics.AgricultureLoans = ae.getAgricultureLoanMetrics(ctx, startTime, endTime)

	// Business lending (VyavasayaMitra)
	metrics.BusinessLoans = ae.getBusinessLoanMetrics(ctx, startTime, endTime)

	// Education lending (ShikshaMitra)
	metrics.EducationLoans = ae.getEducationLoanMetrics(ctx, startTime, endTime)

	// Overall lending statistics
	metrics.TotalLoansOriginated = metrics.AgricultureLoans.Count + metrics.BusinessLoans.Count + metrics.EducationLoans.Count
	metrics.TotalLoanAmount = metrics.AgricultureLoans.Amount.Add(metrics.BusinessLoans.Amount).Add(metrics.EducationLoans.Amount)

	// Risk metrics
	metrics.DefaultRate = ae.calculateOverallDefaultRate(ctx, startTime, endTime)
	metrics.AverageInterestRate = ae.calculateAverageInterestRate(ctx, startTime, endTime)

	// Repayment statistics
	metrics.RepaymentPerformance = ae.getRepaymentPerformance(ctx, startTime, endTime)

	return metrics
}

// generateTrendAnalysis generates trend analysis for key metrics
func (ae *AnalyticsEngine) generateTrendAnalysis(ctx sdk.Context, startTime, endTime time.Time) types.TrendAnalysis {
	analysis := types.TrendAnalysis{
		TimeRange: types.TimeRange{StartTime: startTime, EndTime: endTime},
	}

	// Transaction trends
	analysis.TransactionTrends = ae.calculateTransactionTrends(ctx, startTime, endTime)

	// User adoption trends
	analysis.UserAdoptionTrends = ae.calculateUserAdoptionTrends(ctx, startTime, endTime)

	// Module usage trends
	analysis.ModuleUsageTrends = ae.calculateModuleUsageTrends(ctx, startTime, endTime)

	// Economic trends
	analysis.EconomicTrends = ae.calculateEconomicTrends(ctx, startTime, endTime)

	// Growth indicators
	analysis.GrowthIndicators = ae.calculateGrowthIndicators(ctx, startTime, endTime)

	return analysis
}

// generatePredictiveMetrics generates predictive analytics based on current data
func (ae *AnalyticsEngine) generatePredictiveMetrics(ctx sdk.Context, analytics *ChainAnalytics) types.PredictiveMetrics {
	metrics := types.PredictiveMetrics{}

	// Predict transaction volume growth
	metrics.PredictedTransactionGrowth = ae.predictTransactionGrowth(analytics.TransactionStats, analytics.TrendAnalysis)

	// Predict user adoption
	metrics.PredictedUserAdoption = ae.predictUserAdoption(analytics.TrendAnalysis.UserAdoptionTrends)

	// Predict network load
	metrics.PredictedNetworkLoad = ae.predictNetworkLoad(analytics.NetworkStats, analytics.TransactionStats)

	// Predict token economics
	metrics.PredictedTokenMetrics = ae.predictTokenMetrics(analytics.TokenEconomics)

	// Risk predictions
	metrics.RiskPredictions = ae.generateRiskPredictions(analytics)

	// Growth forecasts
	metrics.GrowthForecasts = ae.generateGrowthForecasts(analytics)

	return metrics
}

// Helper functions for calculations
func (ae *AnalyticsEngine) generateAnalyticsID(ctx sdk.Context) string {
	timestamp := ctx.BlockTime().Unix()
	return fmt.Sprintf("ANALYTICS-%d", timestamp)
}

func (ae *AnalyticsEngine) calculateChainStability(blockTimes []time.Duration) sdk.Dec {
	if len(blockTimes) < 2 {
		return sdk.ZeroDec()
	}

	// Calculate coefficient of variation
	var sum, sumSquares float64
	for _, t := range blockTimes {
		seconds := t.Seconds()
		sum += seconds
		sumSquares += seconds * seconds
	}

	mean := sum / float64(len(blockTimes))
	variance := (sumSquares / float64(len(blockTimes))) - (mean * mean)
	stdDev := variance
	if variance > 0 {
		stdDev = variance // Simplified square root
	}

	coefficientOfVariation := stdDev / mean
	stability := 1.0 - coefficientOfVariation

	if stability < 0 {
		stability = 0
	}
	if stability > 1 {
		stability = 1
	}

	return sdk.NewDecWithPrec(int64(stability*100), 2)
}

func (ae *AnalyticsEngine) calculateTransactionVolume(tx types.Transaction) sdk.Int {
	volume := sdk.ZeroInt()
	for _, msg := range tx.Messages {
		if msg.Amount.Amount.IsPositive() {
			volume = volume.Add(msg.Amount.Amount)
		}
	}
	return volume
}

// Additional helper methods would include all calculation functions
// referenced in the generation methods above