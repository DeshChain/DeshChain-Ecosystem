package types

import (
	"time"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// PerformanceMetrics contains platform performance data for yield calculation
type PerformanceMetrics struct {
	PlatformRevenue       sdk.Dec `json:"platform_revenue"`
	PlatformExpenses      sdk.Dec `json:"platform_expenses"`
	TradingVolume         sdk.Dec `json:"trading_volume"`
	PreviousTradingVolume sdk.Dec `json:"previous_trading_volume"`
	LendingVolume         sdk.Dec `json:"lending_volume"`
	DefaultRate           sdk.Dec `json:"default_rate"`
	DUSDRevenue           sdk.Dec `json:"dusd_revenue"`
	DUSDVolume            sdk.Dec `json:"dusd_volume"`
	Timestamp             time.Time `json:"timestamp"`
}

// YieldDistribution records a yield distribution event
type YieldDistribution struct {
	Timestamp        time.Time `json:"timestamp"`
	YieldRate        sdk.Dec   `json:"yield_rate"`
	TotalDistributed sdk.Int   `json:"total_distributed"`
	TotalSupply      sdk.Int   `json:"total_supply"`
	PerformanceScore sdk.Dec   `json:"performance_score"`
}

// Add yield fields to FeeStructure
type FeeStructure struct {
	MintFee            uint64 `json:"mint_fee"`
	MintFeeCap         string `json:"mint_fee_cap"`
	BurnFee            uint64 `json:"burn_fee"`
	BurnFeeCap         string `json:"burn_fee_cap"`
	LiquidationPenalty uint64 `json:"liquidation_penalty"`
	StabilityFee       uint64 `json:"stability_fee"`
	YieldRateMin       uint64 `json:"yield_rate_min"` // 0% for performance-based
	YieldRateMax       uint64 `json:"yield_rate_max"` // 8% maximum
}