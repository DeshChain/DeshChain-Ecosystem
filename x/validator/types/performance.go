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
	"time"
)

// PerformanceMetrics defines the performance criteria for validator bonuses
type PerformanceMetrics struct {
	UptimeThresholds      map[string]sdk.Dec `json:"uptime_thresholds" yaml:"uptime_thresholds"`
	BlockProductionBonus  map[string]sdk.Dec `json:"block_production_bonus" yaml:"block_production_bonus"`
	TransactionSpeedBonus map[string]sdk.Dec `json:"transaction_speed_bonus" yaml:"transaction_speed_bonus"`
	CommunityBonus        sdk.Dec            `json:"community_bonus" yaml:"community_bonus"`
	ArchiveNodeBonus      sdk.Dec            `json:"archive_node_bonus" yaml:"archive_node_bonus"`
	PublicRPCBonus        sdk.Dec            `json:"public_rpc_bonus" yaml:"public_rpc_bonus"`
}

// DefaultPerformanceMetrics returns the default performance bonus structure
func DefaultPerformanceMetrics() PerformanceMetrics {
	return PerformanceMetrics{
		UptimeThresholds: map[string]sdk.Dec{
			"99.0":  sdk.NewDec(0),                    // No bonus for 99.0%
			"99.5":  sdk.NewDecWithPrec(2, 2),         // 2% bonus for 99.5%
			"99.9":  sdk.NewDecWithPrec(3, 2),         // 3% bonus for 99.9%
			"99.99": sdk.NewDecWithPrec(5, 2),         // 5% bonus for 99.99%
			"100":   sdk.NewDecWithPrec(5, 2),         // 5% bonus for 100%
		},
		BlockProductionBonus: map[string]sdk.Dec{
			"top10":  sdk.NewDecWithPrec(3, 2),        // 3% for top 10% efficiency
			"top25":  sdk.NewDecWithPrec(2, 2),        // 2% for top 25% efficiency
			"average": sdk.NewDec(0),                  // No bonus for average
		},
		TransactionSpeedBonus: map[string]sdk.Dec{
			"fast":   sdk.NewDecWithPrec(3, 2),        // 3% for <100ms average
			"medium": sdk.NewDecWithPrec(2, 2),        // 2% for 100-200ms
			"normal": sdk.NewDecWithPrec(1, 2),        // 1% for 200-500ms
			"slow":   sdk.NewDec(0),                   // No bonus for >500ms
		},
		CommunityBonus:   sdk.NewDecWithPrec(2, 2),    // 2% for active community participation
		ArchiveNodeBonus: sdk.NewDecWithPrec(2, 2),    // 2% for maintaining full archive
		PublicRPCBonus:   sdk.NewDecWithPrec(2, 2),    // 2% for providing free public RPC
	}
}

// ValidatorPerformanceData represents a validator's performance metrics
type ValidatorPerformanceData struct {
	ValidatorAddress     string    `json:"validator_address" yaml:"validator_address"`
	
	// Uptime metrics
	UptimePercentage     sdk.Dec   `json:"uptime_percentage" yaml:"uptime_percentage"`
	TotalBlocks          int64     `json:"total_blocks" yaml:"total_blocks"`
	MissedBlocks         int64     `json:"missed_blocks" yaml:"missed_blocks"`
	
	// Block production metrics
	BlocksProduced       int64     `json:"blocks_produced" yaml:"blocks_produced"`
	AverageBlockTime     int64     `json:"average_block_time" yaml:"average_block_time"` // in milliseconds
	EfficiencyRank       int       `json:"efficiency_rank" yaml:"efficiency_rank"`        // Rank among all validators
	TotalValidators      int       `json:"total_validators" yaml:"total_validators"`
	
	// Transaction processing metrics
	TransactionsProcessed int64    `json:"transactions_processed" yaml:"transactions_processed"`
	AverageResponseTime   int64    `json:"average_response_time" yaml:"average_response_time"` // in milliseconds
	
	// Service provision
	ProvidesArchiveNode   bool     `json:"provides_archive_node" yaml:"provides_archive_node"`
	ProvidesPublicRPC     bool     `json:"provides_public_rpc" yaml:"provides_public_rpc"`
	RPCEndpoint           string   `json:"rpc_endpoint" yaml:"rpc_endpoint"`
	
	// Community participation
	CommunityContributions []CommunityContribution `json:"community_contributions" yaml:"community_contributions"`
	
	// Performance period
	PeriodStart          time.Time `json:"period_start" yaml:"period_start"`
	PeriodEnd            time.Time `json:"period_end" yaml:"period_end"`
	LastUpdated          time.Time `json:"last_updated" yaml:"last_updated"`
}

// CommunityContribution represents different ways validators can contribute to the community
type CommunityContribution struct {
	Type        string    `json:"type" yaml:"type"`               // documentation, tools, support, education
	Description string    `json:"description" yaml:"description"`
	Verified    bool      `json:"verified" yaml:"verified"`
	Points      int       `json:"points" yaml:"points"`           // Contribution score
	Date        time.Time `json:"date" yaml:"date"`
}

// Community contribution types
const (
	ContributionTypeDocumentation = "documentation"
	ContributionTypeTools         = "tools"
	ContributionTypeSupport       = "support"
	ContributionTypeEducation     = "education"
	ContributionTypeBugReports    = "bug_reports"
	ContributionTypeNetworkHealth = "network_health"
)

// CalculatePerformanceMultiplier calculates the total performance bonus multiplier
func (pm PerformanceMetrics) CalculatePerformanceMultiplier(data ValidatorPerformanceData) sdk.Dec {
	multiplier := sdk.OneDec() // Start with 1.0 (no bonus)
	
	// Uptime bonus
	uptimeBonus := pm.getUptimeBonus(data.UptimePercentage)
	multiplier = multiplier.Add(uptimeBonus)
	
	// Block production efficiency bonus
	efficiencyBonus := pm.getEfficiencyBonus(data.EfficiencyRank, data.TotalValidators)
	multiplier = multiplier.Add(efficiencyBonus)
	
	// Transaction speed bonus
	speedBonus := pm.getSpeedBonus(data.AverageResponseTime)
	multiplier = multiplier.Add(speedBonus)
	
	// Archive node bonus
	if data.ProvidesArchiveNode {
		multiplier = multiplier.Add(pm.ArchiveNodeBonus)
	}
	
	// Public RPC bonus
	if data.ProvidesPublicRPC {
		multiplier = multiplier.Add(pm.PublicRPCBonus)
	}
	
	// Community contribution bonus
	if pm.hasCommunityContributions(data.CommunityContributions) {
		multiplier = multiplier.Add(pm.CommunityBonus)
	}
	
	return multiplier
}

// getUptimeBonus calculates uptime-based bonus
func (pm PerformanceMetrics) getUptimeBonus(uptimePercentage sdk.Dec) sdk.Dec {
	// Convert to percentage for comparison
	uptime := uptimePercentage.Mul(sdk.NewDec(100))
	
	switch {
	case uptime.GTE(sdk.NewDec(100)):
		return pm.UptimeThresholds["100"]
	case uptime.GTE(sdk.NewDecWithPrec(9999, 2)):
		return pm.UptimeThresholds["99.99"]
	case uptime.GTE(sdk.NewDecWithPrec(999, 1)):
		return pm.UptimeThresholds["99.9"]
	case uptime.GTE(sdk.NewDecWithPrec(995, 1)):
		return pm.UptimeThresholds["99.5"]
	case uptime.GTE(sdk.NewDecWithPrec(99, 0)):
		return pm.UptimeThresholds["99.0"]
	default:
		return sdk.NewDec(0)
	}
}

// getEfficiencyBonus calculates block production efficiency bonus
func (pm PerformanceMetrics) getEfficiencyBonus(rank int, totalValidators int) sdk.Dec {
	if totalValidators == 0 {
		return sdk.NewDec(0)
	}
	
	percentile := float64(rank) / float64(totalValidators)
	
	switch {
	case percentile <= 0.10: // Top 10%
		return pm.BlockProductionBonus["top10"]
	case percentile <= 0.25: // Top 25%
		return pm.BlockProductionBonus["top25"]
	default:
		return pm.BlockProductionBonus["average"]
	}
}

// getSpeedBonus calculates transaction processing speed bonus
func (pm PerformanceMetrics) getSpeedBonus(averageResponseTime int64) sdk.Dec {
	switch {
	case averageResponseTime < 100: // <100ms
		return pm.TransactionSpeedBonus["fast"]
	case averageResponseTime < 200: // 100-200ms
		return pm.TransactionSpeedBonus["medium"]
	case averageResponseTime < 500: // 200-500ms
		return pm.TransactionSpeedBonus["normal"]
	default: // >500ms
		return pm.TransactionSpeedBonus["slow"]
	}
}

// hasCommunityContributions checks if validator has meaningful community contributions
func (pm PerformanceMetrics) hasCommunityContributions(contributions []CommunityContribution) bool {
	totalPoints := 0
	for _, contrib := range contributions {
		if contrib.Verified {
			totalPoints += contrib.Points
		}
	}
	// Require at least 10 points worth of verified contributions
	return totalPoints >= 10
}

// GetMaxPossiblePerformanceMultiplier returns the maximum achievable performance multiplier
func (pm PerformanceMetrics) GetMaxPossiblePerformanceMultiplier() sdk.Dec {
	return sdk.OneDec().
		Add(pm.UptimeThresholds["100"]).
		Add(pm.BlockProductionBonus["top10"]).
		Add(pm.TransactionSpeedBonus["fast"]).
		Add(pm.CommunityBonus).
		Add(pm.ArchiveNodeBonus).
		Add(pm.PublicRPCBonus)
}

// Validate validates the performance metrics structure
func (pm PerformanceMetrics) Validate() error {
	// Validate uptime thresholds
	for threshold, bonus := range pm.UptimeThresholds {
		if bonus.IsNegative() {
			return fmt.Errorf("uptime bonus for %s cannot be negative: %s", threshold, bonus)
		}
		if bonus.GT(sdk.NewDecWithPrec(20, 2)) { // Max 20% for any single bonus
			return fmt.Errorf("uptime bonus for %s cannot exceed 20%%: %s", threshold, bonus.Mul(sdk.NewDec(100)))
		}
	}
	
	// Validate block production bonuses
	for level, bonus := range pm.BlockProductionBonus {
		if bonus.IsNegative() {
			return fmt.Errorf("block production bonus for %s cannot be negative: %s", level, bonus)
		}
		if bonus.GT(sdk.NewDecWithPrec(20, 2)) {
			return fmt.Errorf("block production bonus for %s cannot exceed 20%%: %s", level, bonus.Mul(sdk.NewDec(100)))
		}
	}
	
	// Validate transaction speed bonuses
	for speed, bonus := range pm.TransactionSpeedBonus {
		if bonus.IsNegative() {
			return fmt.Errorf("transaction speed bonus for %s cannot be negative: %s", speed, bonus)
		}
		if bonus.GT(sdk.NewDecWithPrec(20, 2)) {
			return fmt.Errorf("transaction speed bonus for %s cannot exceed 20%%: %s", speed, bonus.Mul(sdk.NewDec(100)))
		}
	}
	
	// Validate individual bonuses
	bonuses := []struct {
		name  string
		value sdk.Dec
	}{
		{"community_bonus", pm.CommunityBonus},
		{"archive_node_bonus", pm.ArchiveNodeBonus},
		{"public_rpc_bonus", pm.PublicRPCBonus},
	}
	
	for _, bonus := range bonuses {
		if bonus.value.IsNegative() {
			return fmt.Errorf("%s cannot be negative: %s", bonus.name, bonus.value)
		}
		if bonus.value.GT(sdk.NewDecWithPrec(20, 2)) {
			return fmt.Errorf("%s cannot exceed 20%%: %s", bonus.name, bonus.value.Mul(sdk.NewDec(100)))
		}
	}
	
	// Check total possible bonus doesn't exceed 50%
	maxTotal := pm.GetMaxPossiblePerformanceMultiplier().Sub(sdk.OneDec())
	if maxTotal.GT(sdk.NewDecWithPrec(50, 2)) {
		return fmt.Errorf("maximum total performance bonus cannot exceed 50%%: %s", maxTotal.Mul(sdk.NewDec(100)))
	}
	
	return nil
}

// Validate validates the validator performance data
func (vpd ValidatorPerformanceData) Validate() error {
	if vpd.ValidatorAddress == "" {
		return fmt.Errorf("validator address cannot be empty")
	}
	
	if vpd.UptimePercentage.IsNegative() || vpd.UptimePercentage.GT(sdk.OneDec()) {
		return fmt.Errorf("uptime percentage must be between 0 and 1")
	}
	
	if vpd.TotalBlocks < 0 {
		return fmt.Errorf("total blocks cannot be negative")
	}
	
	if vpd.MissedBlocks < 0 {
		return fmt.Errorf("missed blocks cannot be negative")
	}
	
	if vpd.MissedBlocks > vpd.TotalBlocks {
		return fmt.Errorf("missed blocks cannot exceed total blocks")
	}
	
	if vpd.AverageResponseTime < 0 {
		return fmt.Errorf("average response time cannot be negative")
	}
	
	if vpd.EfficiencyRank < 0 {
		return fmt.Errorf("efficiency rank cannot be negative")
	}
	
	if vpd.TotalValidators < 0 {
		return fmt.Errorf("total validators cannot be negative")
	}
	
	if vpd.EfficiencyRank > vpd.TotalValidators {
		return fmt.Errorf("efficiency rank cannot exceed total validators")
	}
	
	return nil
}

// CalculateUptimePercentage calculates uptime percentage from block data
func CalculateUptimePercentage(totalBlocks, missedBlocks int64) sdk.Dec {
	if totalBlocks == 0 {
		return sdk.ZeroDec()
	}
	
	successfulBlocks := totalBlocks - missedBlocks
	return sdk.NewDec(successfulBlocks).Quo(sdk.NewDec(totalBlocks))
}

// GetCommunityContributionPoints returns standard points for different contribution types
func GetCommunityContributionPoints(contributionType string) int {
	switch contributionType {
	case ContributionTypeDocumentation:
		return 5  // 5 points per documentation contribution
	case ContributionTypeTools:
		return 10 // 10 points per tool contribution
	case ContributionTypeSupport:
		return 3  // 3 points per support activity
	case ContributionTypeEducation:
		return 8  // 8 points per educational content
	case ContributionTypeBugReports:
		return 5  // 5 points per verified bug report
	case ContributionTypeNetworkHealth:
		return 15 // 15 points per network health improvement
	default:
		return 1  // 1 point for unrecognized contributions
	}
}

// PerformanceReportingPeriod defines the period for performance evaluation
type PerformanceReportingPeriod struct {
	StartTime time.Time `json:"start_time" yaml:"start_time"`
	EndTime   time.Time `json:"end_time" yaml:"end_time"`
	Type      string    `json:"type" yaml:"type"` // daily, weekly, monthly, quarterly
}

// Standard reporting periods
const (
	PeriodTypeDaily     = "daily"
	PeriodTypeWeekly    = "weekly"
	PeriodTypeMonthly   = "monthly"
	PeriodTypeQuarterly = "quarterly"
)

// GetStandardReportingPeriod returns a standard reporting period
func GetStandardReportingPeriod(periodType string, referenceTime time.Time) PerformanceReportingPeriod {
	switch periodType {
	case PeriodTypeDaily:
		start := time.Date(referenceTime.Year(), referenceTime.Month(), referenceTime.Day(), 0, 0, 0, 0, referenceTime.Location())
		end := start.Add(24 * time.Hour)
		return PerformanceReportingPeriod{StartTime: start, EndTime: end, Type: PeriodTypeDaily}
		
	case PeriodTypeWeekly:
		// Start from Monday
		days := int(referenceTime.Weekday())
		if days == 0 { // Sunday
			days = 7
		}
		start := referenceTime.AddDate(0, 0, -days+1)
		start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
		end := start.Add(7 * 24 * time.Hour)
		return PerformanceReportingPeriod{StartTime: start, EndTime: end, Type: PeriodTypeWeekly}
		
	case PeriodTypeMonthly:
		start := time.Date(referenceTime.Year(), referenceTime.Month(), 1, 0, 0, 0, 0, referenceTime.Location())
		end := start.AddDate(0, 1, 0)
		return PerformanceReportingPeriod{StartTime: start, EndTime: end, Type: PeriodTypeMonthly}
		
	case PeriodTypeQuarterly:
		quarter := (int(referenceTime.Month()) - 1) / 3
		startMonth := time.Month(quarter*3 + 1)
		start := time.Date(referenceTime.Year(), startMonth, 1, 0, 0, 0, 0, referenceTime.Location())
		end := start.AddDate(0, 3, 0)
		return PerformanceReportingPeriod{StartTime: start, EndTime: end, Type: PeriodTypeQuarterly}
		
	default:
		// Default to daily
		return GetStandardReportingPeriod(PeriodTypeDaily, referenceTime)
	}
}