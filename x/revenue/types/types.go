package types

import (
	"time"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RevenueStream represents a source of revenue
type RevenueStream struct {
	StreamID    string    `json:"stream_id"`
	ModuleName  string    `json:"module_name"`
	StreamType  string    `json:"stream_type"`
	Amount      sdk.Coins `json:"amount"`
	Timestamp   time.Time `json:"timestamp"`
	Description string    `json:"description"`
}

// PerformanceMetrics contains platform performance data
type PerformanceMetrics struct {
	Timestamp             time.Time `json:"timestamp"`
	PlatformRevenue       sdk.Dec   `json:"platform_revenue"`
	PlatformExpenses      sdk.Dec   `json:"platform_expenses"`
	TradingVolume         sdk.Dec   `json:"trading_volume"`
	PreviousTradingVolume sdk.Dec   `json:"previous_trading_volume"`
	LendingVolume         sdk.Dec   `json:"lending_volume"`
	DefaultRate           sdk.Dec   `json:"default_rate"`
	DUSDRevenue           sdk.Dec   `json:"dusd_revenue"`
	DUSDVolume            sdk.Dec   `json:"dusd_volume"`
	ActiveUsers           uint64    `json:"active_users"`
	TransactionCount      uint64    `json:"transaction_count"`
}

// RevenueDistribution records how revenue was distributed
type RevenueDistribution struct {
	DistributionID   string    `json:"distribution_id"`
	Timestamp        time.Time `json:"timestamp"`
	TotalRevenue     sdk.Coins `json:"total_revenue"`
	CharityAmount    sdk.Coins `json:"charity_amount"`
	CharityPercent   sdk.Dec   `json:"charity_percent"`
	OperationsAmount sdk.Coins `json:"operations_amount"`
	ReservesAmount   sdk.Coins `json:"reserves_amount"`
	YieldAmount      sdk.Coins `json:"yield_amount"`
}

// ModuleRevenue tracks revenue by module
type ModuleRevenue struct {
	ModuleName       string    `json:"module_name"`
	Period           string    `json:"period"` // daily, weekly, monthly
	Revenue          sdk.Coins `json:"revenue"`
	TransactionCount uint64    `json:"transaction_count"`
	UniqueUsers      uint64    `json:"unique_users"`
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`
}

// YieldCalculation records yield distribution calculations
type YieldCalculation struct {
	CalculationID    string    `json:"calculation_id"`
	Timestamp        time.Time `json:"timestamp"`
	YieldRate        sdk.Dec   `json:"yield_rate"`
	PerformanceScore sdk.Dec   `json:"performance_score"`
	TotalSupply      sdk.Int   `json:"total_supply"`
	YieldAmount      sdk.Int   `json:"yield_amount"`
	Distributed      bool      `json:"distributed"`
}

// PlatformStatistics aggregates platform-wide statistics
type PlatformStatistics struct {
	LastUpdated           time.Time `json:"last_updated"`
	TotalRevenue          sdk.Coins `json:"total_revenue"`
	TotalCharityDistributed sdk.Coins `json:"total_charity_distributed"`
	TotalYieldDistributed sdk.Coins `json:"total_yield_distributed"`
	AverageYieldRate      sdk.Dec   `json:"average_yield_rate"`
	TotalUsers            uint64    `json:"total_users"`
	TotalTransactions     uint64    `json:"total_transactions"`
	PlatformUptime        sdk.Dec   `json:"platform_uptime"`
}