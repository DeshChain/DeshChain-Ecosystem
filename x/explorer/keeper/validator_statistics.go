package keeper

import (
	"fmt"
	"sort"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/explorer/types"
)

// ValidatorStatisticsEngine handles validator performance and statistics
type ValidatorStatisticsEngine struct {
	keeper Keeper
}

// NewValidatorStatisticsEngine creates a new validator statistics engine
func NewValidatorStatisticsEngine(keeper Keeper) *ValidatorStatisticsEngine {
	return &ValidatorStatisticsEngine{
		keeper: keeper,
	}
}

// ValidatorPerformanceReport represents comprehensive validator performance data
type ValidatorPerformanceReport struct {
	ValidatorAddress     string                        `json:"validator_address"`
	ValidatorInfo        types.ValidatorInfo           `json:"validator_info"`
	PerformanceMetrics   types.ValidatorPerformance    `json:"performance_metrics"`
	BlockProduction      types.BlockProductionStats    `json:"block_production"`
	CommissionAnalysis   types.CommissionAnalysis      `json:"commission_analysis"`
	DelegationAnalysis   types.DelegationAnalysis      `json:"delegation_analysis"`
	UptimeAnalysis       types.UptimeAnalysis          `json:"uptime_analysis"`
	SlashingHistory      types.SlashingHistory         `json:"slashing_history"`
	RewardDistribution   types.RewardDistribution      `json:"reward_distribution"`
	NetworkContribution  types.NetworkContribution     `json:"network_contribution"`
	PerformanceRanking   types.PerformanceRanking      `json:"performance_ranking"`
	HistoricalTrends     types.ValidatorTrends         `json:"historical_trends"`
	ReputationScore      sdk.Dec                       `json:"reputation_score"`
	RecommendationLevel  string                        `json:"recommendation_level"`
	GeneratedAt          time.Time                     `json:"generated_at"`
}

// GenerateValidatorPerformanceReport generates comprehensive validator performance report
func (vse *ValidatorStatisticsEngine) GenerateValidatorPerformanceReport(ctx sdk.Context, validatorAddr string, timeRange types.TimeRange) (*ValidatorPerformanceReport, error) {
	// Get validator information
	validator, found := vse.keeper.GetValidator(ctx, validatorAddr)
	if !found {
		return nil, fmt.Errorf("validator not found: %s", validatorAddr)
	}

	report := &ValidatorPerformanceReport{
		ValidatorAddress: validatorAddr,
		ValidatorInfo:    vse.extractValidatorInfo(ctx, validator),
		GeneratedAt:      ctx.BlockTime(),
	}

	// Generate performance metrics
	report.PerformanceMetrics = vse.generatePerformanceMetrics(ctx, validatorAddr, timeRange)

	// Generate block production statistics
	report.BlockProduction = vse.generateBlockProductionStats(ctx, validatorAddr, timeRange)

	// Generate commission analysis
	report.CommissionAnalysis = vse.generateCommissionAnalysis(ctx, validatorAddr, timeRange)

	// Generate delegation analysis
	report.DelegationAnalysis = vse.generateDelegationAnalysis(ctx, validatorAddr, timeRange)

	// Generate uptime analysis
	report.UptimeAnalysis = vse.generateUptimeAnalysis(ctx, validatorAddr, timeRange)

	// Generate slashing history
	report.SlashingHistory = vse.generateSlashingHistory(ctx, validatorAddr, timeRange)

	// Generate reward distribution analysis
	report.RewardDistribution = vse.generateRewardDistribution(ctx, validatorAddr, timeRange)

	// Generate network contribution analysis
	report.NetworkContribution = vse.generateNetworkContribution(ctx, validatorAddr, timeRange)

	// Generate performance ranking
	report.PerformanceRanking = vse.generatePerformanceRanking(ctx, validatorAddr, timeRange)

	// Generate historical trends
	report.HistoricalTrends = vse.generateValidatorTrends(ctx, validatorAddr, timeRange)

	// Calculate reputation score
	report.ReputationScore = vse.calculateReputationScore(report)

	// Determine recommendation level
	report.RecommendationLevel = vse.determineRecommendationLevel(report.ReputationScore)

	// Store validator report
	vse.keeper.SetValidatorPerformanceReport(ctx, *report)

	// Emit report generation event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeValidatorReportGenerated,
			sdk.NewAttribute(types.AttributeKeyValidatorAddress, validatorAddr),
			sdk.NewAttribute(types.AttributeKeyReputationScore, report.ReputationScore.String()),
			sdk.NewAttribute(types.AttributeKeyRecommendationLevel, report.RecommendationLevel),
		),
	)

	return report, nil
}

// generatePerformanceMetrics generates core performance metrics for validator
func (vse *ValidatorStatisticsEngine) generatePerformanceMetrics(ctx sdk.Context, validatorAddr string, timeRange types.TimeRange) types.ValidatorPerformance {
	metrics := types.ValidatorPerformance{
		ValidatorAddress: validatorAddr,
		TimeRange:        timeRange,
	}

	// Get blocks proposed by validator in time range
	blocks := vse.keeper.GetBlocksByProposer(ctx, validatorAddr, timeRange.StartTime, timeRange.EndTime)
	metrics.BlocksProposed = int64(len(blocks))

	// Get total blocks in time range
	totalBlocks := vse.keeper.GetBlocksInTimeRange(ctx, timeRange.StartTime, timeRange.EndTime)
	metrics.TotalBlocks = int64(len(totalBlocks))

	// Calculate block production rate
	if metrics.TotalBlocks > 0 {
		metrics.BlockProductionRate = sdk.NewDec(metrics.BlocksProposed).Quo(sdk.NewDec(metrics.TotalBlocks))
	}

	// Calculate missed blocks
	expectedBlocks := vse.calculateExpectedBlocks(ctx, validatorAddr, timeRange)
	metrics.MissedBlocks = expectedBlocks - metrics.BlocksProposed
	if metrics.MissedBlocks < 0 {
		metrics.MissedBlocks = 0
	}

	// Calculate uptime percentage
	if expectedBlocks > 0 {
		metrics.UptimePercentage = sdk.NewDec(metrics.BlocksProposed).Quo(sdk.NewDec(expectedBlocks))
	}

	// Get signing information
	signingInfo := vse.keeper.GetValidatorSigningInfo(ctx, validatorAddr)
	if signingInfo != nil {
		metrics.StartHeight = signingInfo.StartHeight
		metrics.IndexOffset = signingInfo.IndexOffset
		metrics.JailedUntil = signingInfo.JailedUntil
		metrics.Tombstoned = signingInfo.Tombstoned
		metrics.MissedBlocksCounter = signingInfo.MissedBlocksCounter
	}

	// Calculate performance score
	metrics.PerformanceScore = vse.calculatePerformanceScore(metrics)

	return metrics
}

// generateBlockProductionStats generates detailed block production statistics
func (vse *ValidatorStatisticsEngine) generateBlockProductionStats(ctx sdk.Context, validatorAddr string, timeRange types.TimeRange) types.BlockProductionStats {
	stats := types.BlockProductionStats{
		ValidatorAddress: validatorAddr,
		TimeRange:        timeRange,
	}

	// Get blocks produced by validator
	blocks := vse.keeper.GetBlocksByProposer(ctx, validatorAddr, timeRange.StartTime, timeRange.EndTime)
	stats.TotalBlocksProduced = int64(len(blocks))

	if len(blocks) == 0 {
		return stats
	}

	// Analyze block production patterns
	var blockTimes []time.Duration
	var totalTxs int64
	var totalGasUsed uint64

	for i := 1; i < len(blocks); i++ {
		timeDiff := blocks[i].Time.Sub(blocks[i-1].Time)
		blockTimes = append(blockTimes, timeDiff)
	}

	for _, block := range blocks {
		totalTxs += int64(len(block.Transactions))
		for _, tx := range block.Transactions {
			totalGasUsed += tx.GasUsed
		}
	}

	stats.AverageTransactionsPerBlock = sdk.NewDec(totalTxs).QuoInt64(int64(len(blocks)))
	stats.AverageGasPerBlock = totalGasUsed / uint64(len(blocks))

	// Block timing analysis
	if len(blockTimes) > 0 {
		sort.Slice(blockTimes, func(i, j int) bool {
			return blockTimes[i] < blockTimes[j]
		})

		var totalTime time.Duration
		for _, t := range blockTimes {
			totalTime += t
		}
		stats.AverageBlockTime = totalTime / time.Duration(len(blockTimes))

		// Fastest and slowest blocks
		stats.FastestBlock = blockTimes[0]
		stats.SlowestBlock = blockTimes[len(blockTimes)-1]

		// Block time consistency
		stats.BlockTimeConsistency = vse.calculateBlockTimeConsistency(blockTimes)
	}

	// Daily production pattern
	stats.DailyProductionPattern = vse.calculateDailyProductionPattern(blocks)

	// Production efficiency
	expectedBlocks := vse.calculateExpectedBlocks(ctx, validatorAddr, timeRange)
	if expectedBlocks > 0 {
		stats.ProductionEfficiency = sdk.NewDec(stats.TotalBlocksProduced).Quo(sdk.NewDec(expectedBlocks))
	}

	return stats
}

// generateCommissionAnalysis generates commission rate analysis
func (vse *ValidatorStatisticsEngine) generateCommissionAnalysis(ctx sdk.Context, validatorAddr string, timeRange types.TimeRange) types.CommissionAnalysis {
	analysis := types.CommissionAnalysis{
		ValidatorAddress: validatorAddr,
		TimeRange:        timeRange,
	}

	// Get current validator info
	validator, found := vse.keeper.GetValidator(ctx, validatorAddr)
	if !found {
		return analysis
	}

	analysis.CurrentCommissionRate = validator.Commission.Rate
	analysis.MaxCommissionRate = validator.Commission.MaxRate
	analysis.MaxCommissionChangeRate = validator.Commission.MaxChangeRate

	// Get commission history
	commissionHistory := vse.keeper.GetValidatorCommissionHistory(ctx, validatorAddr, timeRange.StartTime, timeRange.EndTime)
	analysis.CommissionHistory = commissionHistory

	// Calculate commission statistics
	if len(commissionHistory) > 0 {
		var totalCommission sdk.Dec
		var minRate, maxRate sdk.Dec
		minRate = commissionHistory[0].Rate
		maxRate = commissionHistory[0].Rate

		for _, record := range commissionHistory {
			totalCommission = totalCommission.Add(record.Rate)
			if record.Rate.LT(minRate) {
				minRate = record.Rate
			}
			if record.Rate.GT(maxRate) {
				maxRate = record.Rate
			}
		}

		analysis.AverageCommissionRate = totalCommission.QuoInt64(int64(len(commissionHistory)))
		analysis.MinCommissionRate = minRate
		analysis.MaxCommissionRate = maxRate

		// Commission changes
		analysis.CommissionChanges = int64(len(commissionHistory) - 1)
		if analysis.CommissionChanges > 0 {
			totalChange := commissionHistory[len(commissionHistory)-1].Rate.Sub(commissionHistory[0].Rate)
			analysis.TotalCommissionChange = totalChange
		}
	}

	// Commission competitiveness
	allValidators := vse.keeper.GetAllValidators(ctx)
	analysis.CommissionPercentile = vse.calculateCommissionPercentile(validator.Commission.Rate, allValidators)

	// Estimated commission earnings
	analysis.EstimatedEarnings = vse.calculateEstimatedCommissionEarnings(ctx, validatorAddr, timeRange)

	return analysis
}

// generateDelegationAnalysis generates delegation analysis for validator
func (vse *ValidatorStatisticsEngine) generateDelegationAnalysis(ctx sdk.Context, validatorAddr string, timeRange types.TimeRange) types.DelegationAnalysis {
	analysis := types.DelegationAnalysis{
		ValidatorAddress: validatorAddr,
		TimeRange:        timeRange,
	}

	// Get current validator info
	validator, found := vse.keeper.GetValidator(ctx, validatorAddr)
	if !found {
		return analysis
	}

	analysis.TotalDelegations = validator.Tokens
	analysis.SelfDelegation = vse.keeper.GetValidatorSelfDelegation(ctx, validatorAddr)
	analysis.DelegatorShares = validator.DelegatorShares

	// Calculate delegation metrics
	if analysis.TotalDelegations.IsPositive() {
		analysis.SelfDelegationRatio = analysis.SelfDelegation.ToDec().Quo(analysis.TotalDelegations.ToDec())
	}

	// Get delegation history
	delegationHistory := vse.keeper.GetValidatorDelegationHistory(ctx, validatorAddr, timeRange.StartTime, timeRange.EndTime)
	analysis.DelegationHistory = delegationHistory

	// Analyze delegation changes
	var newDelegations, redelegations, undelegations int64
	var newDelegationAmount, redelegationAmount, undelegationAmount sdk.Int

	newDelegationAmount = sdk.ZeroInt()
	redelegationAmount = sdk.ZeroInt()
	undelegationAmount = sdk.ZeroInt()

	for _, record := range delegationHistory {
		switch record.Type {
		case "DELEGATE":
			newDelegations++
			newDelegationAmount = newDelegationAmount.Add(record.Amount.Amount)
		case "REDELEGATE":
			redelegations++
			redelegationAmount = redelegationAmount.Add(record.Amount.Amount)
		case "UNDELEGATE":
			undelegations++
			undelegationAmount = undelegationAmount.Add(record.Amount.Amount)
		}
	}

	analysis.NewDelegations = newDelegations
	analysis.Redelegations = redelegations
	analysis.Undelegations = undelegations
	analysis.NewDelegationAmount = sdk.NewCoin("namo", newDelegationAmount)
	analysis.RedelegationAmount = sdk.NewCoin("namo", redelegationAmount)
	analysis.UndelegationAmount = sdk.NewCoin("namo", undelegationAmount)

	// Net delegation change
	netChange := newDelegationAmount.Add(redelegationAmount).Sub(undelegationAmount)
	analysis.NetDelegationChange = sdk.NewCoin("namo", netChange)

	// Delegation stability
	analysis.DelegationStability = vse.calculateDelegationStability(delegationHistory)

	// Delegator count
	analysis.TotalDelegators = vse.keeper.GetValidatorDelegatorCount(ctx, validatorAddr)

	// Top delegators
	analysis.TopDelegators = vse.keeper.GetValidatorTopDelegators(ctx, validatorAddr, 10)

	// Delegation ranking among all validators
	allValidators := vse.keeper.GetAllValidators(ctx)
	analysis.DelegationRanking = vse.calculateDelegationRanking(validator, allValidators)

	return analysis
}

// generateUptimeAnalysis generates uptime analysis for validator
func (vse *ValidatorStatisticsEngine) generateUptimeAnalysis(ctx sdk.Context, validatorAddr string, timeRange types.TimeRange) types.UptimeAnalysis {
	analysis := types.UptimeAnalysis{
		ValidatorAddress: validatorAddr,
		TimeRange:        timeRange,
	}

	// Get signing information
	signingInfo := vse.keeper.GetValidatorSigningInfo(ctx, validatorAddr)
	if signingInfo != nil {
		analysis.StartHeight = signingInfo.StartHeight
		analysis.MissedBlocksCounter = signingInfo.MissedBlocksCounter
		analysis.IndexOffset = signingInfo.IndexOffset
		analysis.JailedUntil = signingInfo.JailedUntil
		analysis.Tombstoned = signingInfo.Tombstoned
	}

	// Calculate uptime metrics
	blocks := vse.keeper.GetBlocksInTimeRange(ctx, timeRange.StartTime, timeRange.EndTime)
	validatorBlocks := vse.keeper.GetBlocksByProposer(ctx, validatorAddr, timeRange.StartTime, timeRange.EndTime)

	analysis.TotalBlocks = int64(len(blocks))
	analysis.SignedBlocks = int64(len(validatorBlocks))
	analysis.MissedBlocks = analysis.TotalBlocks - analysis.SignedBlocks

	if analysis.TotalBlocks > 0 {
		analysis.UptimePercentage = sdk.NewDec(analysis.SignedBlocks).Quo(sdk.NewDec(analysis.TotalBlocks))
	}

	// Uptime periods
	analysis.UptimePeriods = vse.calculateUptimePeriods(ctx, validatorAddr, timeRange)

	// Downtime analysis
	analysis.DowntimePeriods = vse.calculateDowntimePeriods(ctx, validatorAddr, timeRange)

	// Uptime reliability score
	analysis.ReliabilityScore = vse.calculateReliabilityScore(analysis)

	// Uptime ranking
	allValidators := vse.keeper.GetAllValidators(ctx)
	analysis.UptimeRanking = vse.calculateUptimeRanking(ctx, validatorAddr, allValidators, timeRange)

	return analysis
}

// generateSlashingHistory generates slashing history for validator
func (vse *ValidatorStatisticsEngine) generateSlashingHistory(ctx sdk.Context, validatorAddr string, timeRange types.TimeRange) types.SlashingHistory {
	history := types.SlashingHistory{
		ValidatorAddress: validatorAddr,
		TimeRange:        timeRange,
	}

	// Get slashing events
	slashingEvents := vse.keeper.GetValidatorSlashingEvents(ctx, validatorAddr, timeRange.StartTime, timeRange.EndTime)
	history.SlashingEvents = slashingEvents

	// Analyze slashing events
	var totalSlashed sdk.Int = sdk.ZeroInt()
	var downtimeSlashes, doubleSignSlashes int64

	for _, event := range slashingEvents {
		totalSlashed = totalSlashed.Add(event.SlashedAmount.Amount)
		
		switch event.Reason {
		case "DOWNTIME":
			downtimeSlashes++
		case "DOUBLE_SIGN":
			doubleSignSlashes++
		}
	}

	history.TotalSlashingEvents = int64(len(slashingEvents))
	history.TotalSlashedAmount = sdk.NewCoin("namo", totalSlashed)
	history.DowntimeSlashes = downtimeSlashes
	history.DoubleSignSlashes = doubleSignSlashes

	// Slashing impact
	validator, found := vse.keeper.GetValidator(ctx, validatorAddr)
	if found && validator.Tokens.IsPositive() {
		history.SlashingImpact = totalSlashed.ToDec().Quo(validator.Tokens.ToDec())
	}

	// Time since last slash
	if len(slashingEvents) > 0 {
		lastSlash := slashingEvents[len(slashingEvents)-1].Timestamp
		history.TimeSinceLastSlash = ctx.BlockTime().Sub(lastSlash)
	}

	return history
}

// generateRewardDistribution generates reward distribution analysis
func (vse *ValidatorStatisticsEngine) generateRewardDistribution(ctx sdk.Context, validatorAddr string, timeRange types.TimeRange) types.RewardDistribution {
	distribution := types.RewardDistribution{
		ValidatorAddress: validatorAddr,
		TimeRange:        timeRange,
	}

	// Get reward records
	rewards := vse.keeper.GetValidatorRewards(ctx, validatorAddr, timeRange.StartTime, timeRange.EndTime)
	
	var totalRewards, commissionRewards, delegatorRewards sdk.Int
	totalRewards = sdk.ZeroInt()
	commissionRewards = sdk.ZeroInt()
	delegatorRewards = sdk.ZeroInt()

	for _, reward := range rewards {
		totalRewards = totalRewards.Add(reward.TotalReward.Amount)
		commissionRewards = commissionRewards.Add(reward.CommissionReward.Amount)
		delegatorRewards = delegatorRewards.Add(reward.DelegatorReward.Amount)
	}

	distribution.TotalRewards = sdk.NewCoin("namo", totalRewards)
	distribution.CommissionRewards = sdk.NewCoin("namo", commissionRewards)
	distribution.DelegatorRewards = sdk.NewCoin("namo", delegatorRewards)

	// Average rewards per block
	if len(rewards) > 0 {
		distribution.AverageRewardPerBlock = sdk.NewCoin("namo", totalRewards.QuoRaw(int64(len(rewards))))
	}

	// Reward rate
	validator, found := vse.keeper.GetValidator(ctx, validatorAddr)
	if found && validator.Tokens.IsPositive() {
		duration := timeRange.EndTime.Sub(timeRange.StartTime)
		if duration > 0 {
			annualizedRewards := totalRewards.ToDec().Mul(sdk.NewDec(365)).QuoInt64(int64(duration.Hours()/24))
			distribution.AnnualizedRewardRate = annualizedRewards.Quo(validator.Tokens.ToDec())
		}
	}

	// Reward consistency
	distribution.RewardConsistency = vse.calculateRewardConsistency(rewards)

	return distribution
}

// calculateReputationScore calculates overall reputation score for validator
func (vse *ValidatorStatisticsEngine) calculateReputationScore(report *ValidatorPerformanceReport) sdk.Dec {
	weights := map[string]sdk.Dec{
		"uptime":           sdk.NewDecWithPrec(30, 2), // 30%
		"block_production": sdk.NewDecWithPrec(25, 2), // 25%
		"delegation":       sdk.NewDecWithPrec(15, 2), // 15%
		"commission":       sdk.NewDecWithPrec(10, 2), // 10%
		"slashing":         sdk.NewDecWithPrec(10, 2), // 10%
		"consistency":      sdk.NewDecWithPrec(10, 2), // 10%
	}

	// Score each component (0-100 scale)
	uptimeScore := report.UptimeAnalysis.UptimePercentage.Mul(sdk.NewDec(100))
	blockProductionScore := report.BlockProduction.ProductionEfficiency.Mul(sdk.NewDec(100))
	
	// Delegation score (based on ranking)
	delegationScore := sdk.NewDec(100).Sub(sdk.NewDec(report.DelegationAnalysis.DelegationRanking).Mul(sdk.NewDec(5)))
	if delegationScore.LT(sdk.ZeroDec()) {
		delegationScore = sdk.ZeroDec()
	}

	// Commission score (lower commission = higher score, but not too low)
	commissionScore := vse.calculateCommissionScore(report.CommissionAnalysis.CurrentCommissionRate)

	// Slashing score (penalize slashing events)
	slashingScore := sdk.NewDec(100)
	if report.SlashingHistory.TotalSlashingEvents > 0 {
		penalty := sdk.NewDec(report.SlashingHistory.TotalSlashingEvents).Mul(sdk.NewDec(10))
		slashingScore = slashingScore.Sub(penalty)
		if slashingScore.LT(sdk.ZeroDec()) {
			slashingScore = sdk.ZeroDec()
		}
	}

	// Consistency score
	consistencyScore := report.BlockProduction.BlockTimeConsistency.Mul(sdk.NewDec(100))

	// Calculate weighted average
	totalScore := sdk.ZeroDec()
	totalScore = totalScore.Add(uptimeScore.Mul(weights["uptime"]))
	totalScore = totalScore.Add(blockProductionScore.Mul(weights["block_production"]))
	totalScore = totalScore.Add(delegationScore.Mul(weights["delegation"]))
	totalScore = totalScore.Add(commissionScore.Mul(weights["commission"]))
	totalScore = totalScore.Add(slashingScore.Mul(weights["slashing"]))
	totalScore = totalScore.Add(consistencyScore.Mul(weights["consistency"]))

	return totalScore
}

// determineRecommendationLevel determines recommendation level based on reputation score
func (vse *ValidatorStatisticsEngine) determineRecommendationLevel(score sdk.Dec) string {
	if score.GTE(sdk.NewDec(90)) {
		return "EXCELLENT"
	} else if score.GTE(sdk.NewDec(80)) {
		return "VERY_GOOD"
	} else if score.GTE(sdk.NewDec(70)) {
		return "GOOD"
	} else if score.GTE(sdk.NewDec(60)) {
		return "FAIR"
	} else if score.GTE(sdk.NewDec(50)) {
		return "POOR"
	} else {
		return "NOT_RECOMMENDED"
	}
}

// Helper calculation functions
func (vse *ValidatorStatisticsEngine) calculateExpectedBlocks(ctx sdk.Context, validatorAddr string, timeRange types.TimeRange) int64 {
	// Calculate expected blocks based on validator voting power and time range
	validator, found := vse.keeper.GetValidator(ctx, validatorAddr)
	if !found {
		return 0
	}

	totalValidators := vse.keeper.GetValidatorCount(ctx)
	if totalValidators == 0 {
		return 0
	}

	totalBlocks := vse.keeper.GetBlocksInTimeRange(ctx, timeRange.StartTime, timeRange.EndTime)
	expectedRatio := sdk.OneDec().QuoInt64(int64(totalValidators))
	
	return expectedRatio.MulInt64(int64(len(totalBlocks))).TruncateInt64()
}

func (vse *ValidatorStatisticsEngine) calculatePerformanceScore(metrics types.ValidatorPerformance) sdk.Dec {
	// Simplified performance score calculation
	score := sdk.NewDec(100)

	// Penalize for missed blocks
	if metrics.TotalBlocks > 0 {
		missedRatio := sdk.NewDec(metrics.MissedBlocks).Quo(sdk.NewDec(metrics.TotalBlocks))
		penalty := missedRatio.Mul(sdk.NewDec(50)) // Up to 50 point penalty
		score = score.Sub(penalty)
	}

	// Bonus for high uptime
	if metrics.UptimePercentage.GTE(sdk.NewDecWithPrec(99, 2)) {
		score = score.Add(sdk.NewDec(10))
	}

	if score.LT(sdk.ZeroDec()) {
		score = sdk.ZeroDec()
	}

	return score
}

func (vse *ValidatorStatisticsEngine) calculateCommissionScore(commissionRate sdk.Dec) sdk.Dec {
	// Optimal commission rate is around 5-10%
	optimal := sdk.NewDecWithPrec(75, 3) // 7.5%
	
	if commissionRate.Equal(optimal) {
		return sdk.NewDec(100)
	}

	// Calculate distance from optimal
	distance := commissionRate.Sub(optimal).Abs()
	
	// Score decreases as distance increases
	score := sdk.NewDec(100).Sub(distance.Mul(sdk.NewDec(200)))
	
	if score.LT(sdk.NewDec(20)) {
		score = sdk.NewDec(20) // Minimum score
	}
	
	return score
}

// Additional helper methods would include all calculation functions
// referenced in the generation methods above