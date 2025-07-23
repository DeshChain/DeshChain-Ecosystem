package types

// Revenue module event types
const (
	EventTypeRevenueRecorded    = "revenue_recorded"
	EventTypeRevenueDistributed = "revenue_distributed"
	EventTypeYieldCalculated    = "yield_calculated"
	EventTypeMetricsUpdated     = "metrics_updated"
	
	AttributeKeyStreamID        = "stream_id"
	AttributeKeyModule          = "module"
	AttributeKeyAmount          = "amount"
	AttributeKeyDistributionID  = "distribution_id"
	AttributeKeyTotalRevenue    = "total_revenue"
	AttributeKeyCharityAmount   = "charity_amount"
	AttributeKeyCharityPercent  = "charity_percent"
	AttributeKeyYieldRate       = "yield_rate"
	AttributeKeyPerformanceScore = "performance_score"
)