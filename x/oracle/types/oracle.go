package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

// OracleSourceStatus represents the status of an oracle data source
type OracleSourceStatus struct {
	Name    string  `json:"name"`
	Healthy bool    `json:"healthy"`
	Weight  sdk.Dec `json:"weight"`
}

// PriceFeed represents a price feed for a specific symbol
type PriceFeed struct {
	Symbol      string    `json:"symbol"`
	Price       sdk.Dec   `json:"price"`
	Timestamp   time.Time `json:"timestamp"`
	BlockHeight int64     `json:"block_height"`
	Sources     []string  `json:"sources"`
}

// OracleParams defines the parameters for the oracle module
type OracleParams struct {
	// MinSources is the minimum number of healthy sources required for price aggregation
	MinSources int `json:"min_sources"`
	
	// MaxDeviation is the maximum allowed price deviation between sources (as decimal percentage)
	MaxDeviation sdk.Dec `json:"max_deviation"`
	
	// UpdateInterval is the minimum time between price updates (in seconds)
	UpdateInterval int64 `json:"update_interval"`
	
	// PriceExpiryTime is how long a price remains valid (in seconds)
	PriceExpiryTime int64 `json:"price_expiry_time"`
	
	// EnabledSymbols is the list of symbols for which price feeds are enabled
	EnabledSymbols []string `json:"enabled_symbols"`
}

// DefaultParams returns default oracle parameters
func DefaultParams() OracleParams {
	return OracleParams{
		MinSources:      2,
		MaxDeviation:    sdk.NewDecWithPrec(5, 2), // 5% maximum deviation
		UpdateInterval:  60,                       // 1 minute
		PriceExpiryTime: 300,                      // 5 minutes
		EnabledSymbols:  []string{"DINR", "BTC", "ETH", "NAMO"},
	}
}

// Validate validates the oracle parameters
func (p OracleParams) Validate() error {
	if p.MinSources < 1 {
		return ErrInvalidMinSources
	}
	
	if p.MaxDeviation.IsNegative() || p.MaxDeviation.GT(sdk.OneDec()) {
		return ErrInvalidMaxDeviation
	}
	
	if p.UpdateInterval <= 0 {
		return ErrInvalidUpdateInterval
	}
	
	if p.PriceExpiryTime <= 0 {
		return ErrInvalidPriceExpiryTime
	}
	
	if len(p.EnabledSymbols) == 0 {
		return ErrNoEnabledSymbols
	}
	
	return nil
}

// OracleData represents stored oracle price data
type OracleData struct {
	Symbol       string    `json:"symbol"`
	Price        sdk.Dec   `json:"price"`
	Confidence   sdk.Dec   `json:"confidence"`
	Timestamp    time.Time `json:"timestamp"`
	BlockHeight  int64     `json:"block_height"`
	SourceCount  int       `json:"source_count"`
	Sources      []string  `json:"sources"`
	LastUpdateBy string    `json:"last_update_by"`
}

// IsExpired checks if the oracle data is expired based on params
func (od *OracleData) IsExpired(params OracleParams) bool {
	expiryTime := od.Timestamp.Add(time.Duration(params.PriceExpiryTime) * time.Second)
	return time.Now().After(expiryTime)
}

// GetAge returns the age of the oracle data in seconds
func (od *OracleData) GetAge() int64 {
	return int64(time.Since(od.Timestamp).Seconds())
}

// PriceUpdate represents a price update event
type PriceUpdate struct {
	Symbol      string    `json:"symbol"`
	OldPrice    sdk.Dec   `json:"old_price"`
	NewPrice    sdk.Dec   `json:"new_price"`
	Change      sdk.Dec   `json:"change"`
	ChangeRate  sdk.Dec   `json:"change_rate"`
	Timestamp   time.Time `json:"timestamp"`
	BlockHeight int64     `json:"block_height"`
	UpdatedBy   string    `json:"updated_by"`
}

// CalculateChange calculates the price change and change rate
func (pu *PriceUpdate) CalculateChange() {
	pu.Change = pu.NewPrice.Sub(pu.OldPrice)
	
	if !pu.OldPrice.IsZero() {
		pu.ChangeRate = pu.Change.Quo(pu.OldPrice)
	} else {
		pu.ChangeRate = sdk.ZeroDec()
	}
}

// OracleMetrics represents oracle performance metrics
type OracleMetrics struct {
	TotalUpdates    int64     `json:"total_updates"`
	FailedUpdates   int64     `json:"failed_updates"`
	LastUpdate      time.Time `json:"last_update"`
	AverageLatency  int64     `json:"average_latency_ms"`
	UptimePercent   sdk.Dec   `json:"uptime_percent"`
	ActiveSources   int       `json:"active_sources"`
	TotalSources    int       `json:"total_sources"`
}

// CalculateSuccessRate calculates the success rate of oracle updates
func (om *OracleMetrics) CalculateSuccessRate() sdk.Dec {
	if om.TotalUpdates == 0 {
		return sdk.ZeroDec()
	}
	
	successfulUpdates := om.TotalUpdates - om.FailedUpdates
	return sdk.NewDec(successfulUpdates).Quo(sdk.NewDec(om.TotalUpdates))
}

// IsHealthy checks if the oracle metrics indicate a healthy system
func (om *OracleMetrics) IsHealthy() bool {
	// Consider healthy if:
	// 1. Success rate > 95%
	// 2. Has recent updates (within last 10 minutes)
	// 3. Has at least 1 active source
	
	successRate := om.CalculateSuccessRate()
	recentUpdate := time.Since(om.LastUpdate) < 10*time.Minute
	hasSources := om.ActiveSources > 0
	
	return successRate.GT(sdk.NewDecWithPrec(95, 2)) && recentUpdate && hasSources
}

// OracleAlert represents an oracle system alert
type OracleAlert struct {
	ID          string              `json:"id"`
	Type        OracleAlertType     `json:"type"`
	Severity    OracleAlertSeverity `json:"severity"`
	Symbol      string              `json:"symbol,omitempty"`
	Message     string              `json:"message"`
	Timestamp   time.Time           `json:"timestamp"`
	BlockHeight int64               `json:"block_height"`
	Resolved    bool                `json:"resolved"`
	ResolvedAt  *time.Time          `json:"resolved_at,omitempty"`
}

// OracleAlertType defines the type of oracle alert
type OracleAlertType string

const (
	AlertTypePriceDeviation  OracleAlertType = "price_deviation"
	AlertTypeSourceDown      OracleAlertType = "source_down"
	AlertTypeDataStale       OracleAlertType = "data_stale"
	AlertTypeUpdateFailed    OracleAlertType = "update_failed"
	AlertTypeSystemHealth    OracleAlertType = "system_health"
)

// OracleAlertSeverity defines the severity level of an alert
type OracleAlertSeverity string

const (
	SeverityLow      OracleAlertSeverity = "low"
	SeverityMedium   OracleAlertSeverity = "medium"
	SeverityHigh     OracleAlertSeverity = "high"
	SeverityCritical OracleAlertSeverity = "critical"
)

// Resolve marks the alert as resolved
func (oa *OracleAlert) Resolve() {
	oa.Resolved = true
	now := time.Now()
	oa.ResolvedAt = &now
}

// GetDuration returns how long the alert has been active
func (oa *OracleAlert) GetDuration() time.Duration {
	if oa.Resolved && oa.ResolvedAt != nil {
		return oa.ResolvedAt.Sub(oa.Timestamp)
	}
	return time.Since(oa.Timestamp)
}

// OracleConfiguration represents the complete oracle configuration
type OracleConfiguration struct {
	Params        OracleParams          `json:"params"`
	SourceConfigs []OracleSourceConfig  `json:"source_configs"`
	EnabledAlerts []OracleAlertType     `json:"enabled_alerts"`
	UpdatedAt     time.Time             `json:"updated_at"`
	UpdatedBy     string                `json:"updated_by"`
}

// OracleSourceConfig represents configuration for an oracle source
type OracleSourceConfig struct {
	Name      string  `json:"name"`
	Endpoint  string  `json:"endpoint"`
	APIKey    string  `json:"api_key,omitempty"`
	Weight    sdk.Dec `json:"weight"`
	Enabled   bool    `json:"enabled"`
	Timeout   int64   `json:"timeout_seconds"`
}

// Validate validates the oracle source configuration
func (osc *OracleSourceConfig) Validate() error {
	if osc.Name == "" {
		return ErrInvalidSourceName
	}
	
	if osc.Endpoint == "" {
		return ErrInvalidSourceEndpoint
	}
	
	if osc.Weight.IsNegative() || osc.Weight.IsZero() {
		return ErrInvalidSourceWeight
	}
	
	if osc.Timeout <= 0 {
		return ErrInvalidSourceTimeout
	}
	
	return nil
}