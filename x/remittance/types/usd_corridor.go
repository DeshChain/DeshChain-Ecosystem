package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CorridorOptimization represents the optimization analysis for a remittance corridor
type CorridorOptimization struct {
	SourceCurrency       string        `json:"source_currency"`
	CorridorCurrency     string        `json:"corridor_currency"`
	DestinationCurrency  string        `json:"destination_currency"`
	SourceAmount         sdk.Coin      `json:"source_amount"`
	CorridorAmount       sdk.Coin      `json:"corridor_amount"`
	DestinationAmount    sdk.Coin      `json:"destination_amount"`
	SourceToCorridorRate sdk.Dec       `json:"source_to_corridor_rate"`
	CorridorToDestRate   sdk.Dec       `json:"corridor_to_dest_rate"`
	TotalCost            sdk.Coin      `json:"total_cost"`
	TraditionalCost      sdk.Coin      `json:"traditional_cost"`
	Savings              sdk.Coin      `json:"savings"`
	ProcessingTime       time.Duration `json:"processing_time"`
	EfficiencyScore      sdk.Dec       `json:"efficiency_score"`
	Route                []string      `json:"route"`
}

// CorridorStats represents statistics for a specific remittance corridor
type CorridorStats struct {
	Corridor              string        `json:"corridor"`
	TotalVolume           sdk.Coin      `json:"total_volume"`
	TotalTransactions     int64         `json:"total_transactions"`
	AverageAmount         sdk.Coin      `json:"average_amount"`
	TotalSavings          sdk.Coin      `json:"total_savings"`
	AverageProcessingTime time.Duration `json:"average_processing_time"`
	TraditionalCost       sdk.Coin      `json:"traditional_cost"`
	DeshChainCost         sdk.Coin      `json:"deshchain_cost"`
	CostSavingPercent     sdk.Dec       `json:"cost_saving_percent"`
}

// Enhanced Transfer key prefix
var EnhancedTransferKeyPrefix = []byte{0x20}

// GetEnhancedTransferKey returns the store key for enhanced transfer data
func GetEnhancedTransferKey(transferID string) []byte {
	return append(EnhancedTransferKeyPrefix, []byte(transferID)...)
}

// MultiCurrencyRemittanceRequest represents a request for multi-currency remittance
type MultiCurrencyRemittanceRequest struct {
	SenderId            string   `json:"sender_id"`
	RecipientId         string   `json:"recipient_id"`
	SourceCurrency      string   `json:"source_currency"`
	SourceAmount        sdk.Coin `json:"source_amount"`
	DestinationCurrency string   `json:"destination_currency"`
	SewaMitraId         string   `json:"sewa_mitra_id,omitempty"`
	Notes               string   `json:"notes,omitempty"`
}

// MultiCurrencyRemittanceResponse represents the response for multi-currency remittance
type MultiCurrencyRemittanceResponse struct {
	TransferId           string            `json:"transfer_id"`
	SourceAmount         sdk.Coin          `json:"source_amount"`
	DestinationAmount    sdk.Coin          `json:"destination_amount"`
	RoutingCurrency      string            `json:"routing_currency"`
	OptimalRoute         []string          `json:"optimal_route"`
	ExchangeRates        map[string]sdk.Dec `json:"exchange_rates"`
	TotalFees            sdk.Coin          `json:"total_fees"`
	TotalSavings         sdk.Coin          `json:"total_savings"`
	ProcessingTime       time.Duration     `json:"processing_time"`
	EfficiencyScore      sdk.Dec           `json:"efficiency_score"`
	Status               string            `json:"status"`
}

// CorridorPerformanceMetrics represents performance metrics for corridors
type CorridorPerformanceMetrics struct {
	Corridor             string            `json:"corridor"`
	Volume24h            sdk.Coin          `json:"volume_24h"`
	Volume7d             sdk.Coin          `json:"volume_7d"`
	Volume30d            sdk.Coin          `json:"volume_30d"`
	Transactions24h      int64             `json:"transactions_24h"`
	Transactions7d       int64             `json:"transactions_7d"`
	Transactions30d      int64             `json:"transactions_30d"`
	AverageAmount        sdk.Coin          `json:"average_amount"`
	MedianAmount         sdk.Coin          `json:"median_amount"`
	AverageProcessingTime time.Duration    `json:"average_processing_time"`
	SuccessRate          sdk.Dec           `json:"success_rate"`
	CostEfficiency       sdk.Dec           `json:"cost_efficiency"`
	PopularRoutes        []string          `json:"popular_routes"`
	PeakHours            []int             `json:"peak_hours"`
}

// SewaMitraUSDCapability represents USD handling capabilities of Sewa Mitra agents
type SewaMitraUSDCapability struct {
	AgentId              string            `json:"agent_id"`
	SupportedCurrencies  []string          `json:"supported_currencies"`
	USDLimits            SewaMitraLimits   `json:"usd_limits"`
	ExchangeRateMarkup   sdk.Dec           `json:"exchange_rate_markup"`
	ProcessingTime       time.Duration     `json:"processing_time"`
	ServiceFee           sdk.Coin          `json:"service_fee"`
	AvailabilityHours    string            `json:"availability_hours"`
	Location             SewaMitraLocation `json:"location"`
	Rating               sdk.Dec           `json:"rating"`
	CompletedTransfers   int64             `json:"completed_transfers"`
	USDExperience        bool              `json:"usd_experience"`
}

// SewaMitraLimits represents transaction limits for Sewa Mitra agents
type SewaMitraLimits struct {
	MinAmount      sdk.Coin `json:"min_amount"`
	MaxAmount      sdk.Coin `json:"max_amount"`
	DailyLimit     sdk.Coin `json:"daily_limit"`
	MonthlyLimit   sdk.Coin `json:"monthly_limit"`
}

// SewaMitraLocation represents the location of a Sewa Mitra agent
type SewaMitraLocation struct {
	Country    string  `json:"country"`
	State      string  `json:"state"`
	City       string  `json:"city"`
	Address    string  `json:"address"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	Timezone   string  `json:"timezone"`
}

// RemittanceAnalytics represents comprehensive analytics for remittance operations
type RemittanceAnalytics struct {
	TotalVolume           sdk.Coin                         `json:"total_volume"`
	TotalTransactions     int64                            `json:"total_transactions"`
	TotalSavings          sdk.Coin                         `json:"total_savings"`
	AverageProcessingTime time.Duration                    `json:"average_processing_time"`
	CorridorBreakdown     map[string]CorridorStats         `json:"corridor_breakdown"`
	CurrencyBreakdown     map[string]sdk.Coin              `json:"currency_breakdown"`
	MonthlyTrends         []MonthlyRemittanceStats         `json:"monthly_trends"`
	TopCorridors          []CorridorStats                  `json:"top_corridors"`
	SewaMitraStats        SewaMitraAggregateStats          `json:"sewa_mitra_stats"`
}

// MonthlyRemittanceStats represents monthly statistics
type MonthlyRemittanceStats struct {
	Month              string   `json:"month"`
	Volume             sdk.Coin `json:"volume"`
	Transactions       int64    `json:"transactions"`
	Savings            sdk.Coin `json:"savings"`
	NewCorridors       int      `json:"new_corridors"`
	ActiveSewaMitras   int      `json:"active_sewa_mitras"`
}

// SewaMitraAggregateStats represents aggregate statistics for Sewa Mitra network
type SewaMitraAggregateStats struct {
	TotalAgents         int64     `json:"total_agents"`
	ActiveAgents        int64     `json:"active_agents"`
	USDCapableAgents    int64     `json:"usd_capable_agents"`
	TotalVolume         sdk.Coin  `json:"total_volume"`
	TotalTransactions   int64     `json:"total_transactions"`
	AverageRating       sdk.Dec   `json:"average_rating"`
	CoverageCountries   []string  `json:"coverage_countries"`
	AverageProcessingTime time.Duration `json:"average_processing_time"`
}