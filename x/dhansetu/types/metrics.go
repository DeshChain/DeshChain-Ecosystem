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
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MitraPerformanceMetrics represents performance metrics for a mitra
type MitraPerformanceMetrics struct {
	MitraId           string    `json:"mitra_id" yaml:"mitra_id"`
	DhanPataName      string    `json:"dhanpata_name" yaml:"dhanpata_name"`
	TrustScore        int64     `json:"trust_score" yaml:"trust_score"`
	TotalTrades       uint64    `json:"total_trades" yaml:"total_trades"`
	SuccessfulTrades  uint64    `json:"successful_trades" yaml:"successful_trades"`
	SuccessRate       float64   `json:"success_rate" yaml:"success_rate"`
	DailyVolume       sdk.Int   `json:"daily_volume" yaml:"daily_volume"`
	MonthlyVolume     sdk.Int   `json:"monthly_volume" yaml:"monthly_volume"`
	DailyLimit        sdk.Int   `json:"daily_limit" yaml:"daily_limit"`
	MonthlyLimit      sdk.Int   `json:"monthly_limit" yaml:"monthly_limit"`
	DailyUtilization  float64   `json:"daily_utilization" yaml:"daily_utilization"`
	ActiveEscrowCount uint64    `json:"active_escrow_count" yaml:"active_escrow_count"`
	LastActiveAt      time.Time `json:"last_active_at" yaml:"last_active_at"`
}

// DhanSetuEcosystemMetrics represents overall ecosystem metrics
type DhanSetuEcosystemMetrics struct {
	TotalDhanPataAddresses  uint64 `json:"total_dhanpata_addresses" yaml:"total_dhanpata_addresses"`
	TotalKshetraCoins       uint64 `json:"total_kshetra_coins" yaml:"total_kshetra_coins"`
	TotalEnhancedMitras     uint64 `json:"total_enhanced_mitras" yaml:"total_enhanced_mitras"`
	TotalTradeVolume        sdk.Int `json:"total_trade_volume" yaml:"total_trade_volume"`
	TotalFeesDisbursed      sdk.Int `json:"total_fees_disbursed" yaml:"total_fees_disbursed"`
	NGODonationsTotal       sdk.Int `json:"ngo_donations_total" yaml:"ngo_donations_total"`
	FounderRoyaltyTotal     sdk.Int `json:"founder_royalty_total" yaml:"founder_royalty_total"`
	
	// Activity metrics
	DailyActiveUsers        uint64 `json:"daily_active_users" yaml:"daily_active_users"`
	MonthlyActiveUsers      uint64 `json:"monthly_active_users" yaml:"monthly_active_users"`
	AverageTransactionSize  sdk.Dec `json:"average_transaction_size" yaml:"average_transaction_size"`
	
	// Geographic distribution
	TopPincodesByVolume     []PincodeMetrics `json:"top_pincodes_by_volume" yaml:"top_pincodes_by_volume"`
	RegionalDistribution    map[string]uint64 `json:"regional_distribution" yaml:"regional_distribution"`
}

// PincodeMetrics represents metrics for a specific pincode
type PincodeMetrics struct {
	Pincode             string  `json:"pincode" yaml:"pincode"`
	ActiveUsers         uint64  `json:"active_users" yaml:"active_users"`
	TotalVolume         sdk.Int `json:"total_volume" yaml:"total_volume"`
	HasKshetraCoin      bool    `json:"has_kshetra_coin" yaml:"has_kshetra_coin"`
	KshetraCoinSymbol   string  `json:"kshetra_coin_symbol,omitempty" yaml:"kshetra_coin_symbol,omitempty"`
	LocalMitrasCount    uint64  `json:"local_mitras_count" yaml:"local_mitras_count"`
}

// DhanPataAnalytics represents analytics for DhanPata addresses
type DhanPataAnalytics struct {
	Name                string            `json:"name" yaml:"name"`
	AddressType         string            `json:"address_type" yaml:"address_type"`
	TotalTransactions   uint64            `json:"total_transactions" yaml:"total_transactions"`
	TotalVolume         sdk.Int           `json:"total_volume" yaml:"total_volume"`
	AverageTransaction  sdk.Dec           `json:"average_transaction" yaml:"average_transaction"`
	MostUsedProducts    []string          `json:"most_used_products" yaml:"most_used_products"`
	PreferredMitras     []string          `json:"preferred_mitras" yaml:"preferred_mitras"`
	GeographicReach     []string          `json:"geographic_reach" yaml:"geographic_reach"`
	
	// Time-based metrics
	DailyTransactions   map[string]uint64 `json:"daily_transactions" yaml:"daily_transactions"`
	MonthlyVolume       map[string]sdk.Int `json:"monthly_volume" yaml:"monthly_volume"`
	
	// Engagement metrics
	LastActiveDate      time.Time         `json:"last_active_date" yaml:"last_active_date"`
	AccountAge          int64             `json:"account_age" yaml:"account_age"` // Days since creation
	EngagementScore     int64             `json:"engagement_score" yaml:"engagement_score"` // 0-100
}

// KshetraCoinMetrics represents metrics for Kshetra coins
type KshetraCoinMetrics struct {
	Pincode             string    `json:"pincode" yaml:"pincode"`
	CoinName            string    `json:"coin_name" yaml:"coin_name"`
	CoinSymbol          string    `json:"coin_symbol" yaml:"coin_symbol"`
	MarketCap           sdk.Int   `json:"market_cap" yaml:"market_cap"`
	HolderCount         uint64    `json:"holder_count" yaml:"holder_count"`
	DailyVolume         sdk.Int   `json:"daily_volume" yaml:"daily_volume"`
	CommunityFundSize   sdk.Int   `json:"community_fund_size" yaml:"community_fund_size"`
	NGODonationsTotal   sdk.Int   `json:"ngo_donations_total" yaml:"ngo_donations_total"`
	LocalAdoptionRate   float64   `json:"local_adoption_rate" yaml:"local_adoption_rate"`
	CrossPincodeHolders uint64    `json:"cross_pincode_holders" yaml:"cross_pincode_holders"`
	CreatedAt           time.Time `json:"created_at" yaml:"created_at"`
	DaysSinceCreation   int64     `json:"days_since_creation" yaml:"days_since_creation"`
}

// CrossModuleIntegrationMetrics represents metrics for cross-module integration
type CrossModuleIntegrationMetrics struct {
	MoneyOrderIntegration struct {
		TotalOrders         uint64  `json:"total_orders" yaml:"total_orders"`
		DhanPataOrders      uint64  `json:"dhanpata_orders" yaml:"dhanpata_orders"`
		DhanPataAdoption    float64 `json:"dhanpata_adoption" yaml:"dhanpata_adoption"` // Percentage
		AverageFeeGenerated sdk.Dec `json:"average_fee_generated" yaml:"average_fee_generated"`
	} `json:"moneyorder_integration" yaml:"moneyorder_integration"`
	
	CulturalIntegration struct {
		FestivalBonusesClaimed uint64  `json:"festival_bonuses_claimed" yaml:"festival_bonuses_claimed"`
		CulturalQuotesShown    uint64  `json:"cultural_quotes_shown" yaml:"cultural_quotes_shown"`
		RegionalPreferences    map[string]uint64 `json:"regional_preferences" yaml:"regional_preferences"`
	} `json:"cultural_integration" yaml:"cultural_integration"`
	
	NAMOIntegration struct {
		VestingClaims       uint64  `json:"vesting_claims" yaml:"vesting_claims"`
		TokenBurns          uint64  `json:"token_burns" yaml:"token_burns"`
		DhanSetuRewards     sdk.Int `json:"dhansetu_rewards" yaml:"dhansetu_rewards"`
	} `json:"namo_integration" yaml:"namo_integration"`
}

// CalculateEngagementScore calculates engagement score based on activity
func CalculateEngagementScore(analytics DhanPataAnalytics) int64 {
	score := int64(0)
	
	// Base activity score (0-40 points)
	if analytics.TotalTransactions > 100 {
		score += 40
	} else if analytics.TotalTransactions > 50 {
		score += 30
	} else if analytics.TotalTransactions > 10 {
		score += 20
	} else if analytics.TotalTransactions > 0 {
		score += 10
	}
	
	// Volume score (0-20 points)
	// This would be calculated based on volume relative to user type
	score += 15 // Simplified
	
	// Recency score (0-20 points)
	daysSinceActive := time.Since(analytics.LastActiveDate).Hours() / 24
	if daysSinceActive <= 1 {
		score += 20
	} else if daysSinceActive <= 7 {
		score += 15
	} else if daysSinceActive <= 30 {
		score += 10
	} else if daysSinceActive <= 90 {
		score += 5
	}
	
	// Diversification score (0-20 points)
	if len(analytics.MostUsedProducts) >= 3 {
		score += 20
	} else if len(analytics.MostUsedProducts) >= 2 {
		score += 15
	} else if len(analytics.MostUsedProducts) >= 1 {
		score += 10
	}
	
	// Cap at 100
	if score > 100 {
		score = 100
	}
	
	return score
}