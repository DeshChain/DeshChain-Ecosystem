package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DUSD-specific treasury types

// ReserveAsset represents a USD reserve asset in the treasury
type ReserveAsset struct {
	AssetID      string     `json:"asset_id"`
	AssetType    string     `json:"asset_type"`    // STABLECOIN, GOVERNMENT_BOND, CORPORATE_BOND, CASH
	Currency     string     `json:"currency"`      // USD, USDC, USDT, etc.
	Value        sdk.Coin   `json:"value"`         // Current value
	Liquidity    sdk.Dec    `json:"liquidity"`     // Liquidity score 0-1
	Risk         sdk.Dec    `json:"risk"`          // Risk score 0-1
	Yield        sdk.Dec    `json:"yield"`         // Annual yield percentage
	MaturityDate *MaturityDate `json:"maturity_date,omitempty"` // For bonds
}

// MaturityDate represents the maturity date of a financial instrument
type MaturityDate struct {
	Date time.Time `json:"date"`
}

// USDRebalanceStrategy defines rebalancing strategy for USD reserves
type USDRebalanceStrategy struct {
	Strategy           string                 `json:"strategy"`            // CONSERVATIVE, AGGRESSIVE, OPERATIONAL, CROSS_CURRENCY
	MaxDeviation       sdk.Dec               `json:"max_deviation"`       // Maximum allowed deviation
	RebalanceFrequency time.Duration         `json:"rebalance_frequency"` // How often to rebalance
	AssetWeights       map[string]sdk.Dec    `json:"asset_weights"`       // Target weights for each asset type
	RiskLimits         map[string]sdk.Dec    `json:"risk_limits"`         // Risk limits for different categories
}

// RebalanceAction represents a rebalancing action to be executed
type RebalanceAction struct {
	Currency        string    `json:"currency"`
	CurrentExposure sdk.Dec   `json:"current_exposure"`
	TargetExposure  sdk.Dec   `json:"target_exposure"`
	Deviation       sdk.Dec   `json:"deviation"`
	Action          string    `json:"action"`          // INCREASE, REDUCE, MAINTAIN
	Amount          sdk.Coin  `json:"amount"`          // Amount to buy/sell
	Urgency         string    `json:"urgency"`         // LOW, MEDIUM, HIGH, CRITICAL
	Timestamp       time.Time `json:"timestamp"`
}

// DUSDPoolKey prefixes
var (
	DUSDPoolKeyPrefix = []byte{0x30}
)

// GetDUSDPoolKey returns the store key for DUSD pool data
func GetDUSDPoolKey(poolID string) []byte {
	return append(DUSDPoolKeyPrefix, []byte(poolID)...)
}

// DUSDReserveStats represents statistics for DUSD reserves
type DUSDReserveStats struct {
	TotalUSDReserves       sdk.Coin                   `json:"total_usd_reserves"`
	TotalDUSDSupply        sdk.Coin                   `json:"total_dusd_supply"`
	CollateralizationRatio sdk.Dec                    `json:"collateralization_ratio"`
	ReserveBreakdown       map[string]sdk.Coin        `json:"reserve_breakdown"`
	AssetAllocation        map[string]sdk.Dec         `json:"asset_allocation"`
	RiskMetrics            DUSDRiskMetrics            `json:"risk_metrics"`
	PerformanceMetrics     DUSDPerformanceMetrics     `json:"performance_metrics"`
	RebalanceHistory       []DUSDRebalanceRecord      `json:"rebalance_history"`
	LastUpdated            time.Time                  `json:"last_updated"`
}

// DUSDRiskMetrics represents risk metrics for DUSD reserves
type DUSDRiskMetrics struct {
	VaR95              sdk.Dec   `json:"var_95"`               // Value at Risk (95% confidence)
	VaR99              sdk.Dec   `json:"var_99"`               // Value at Risk (99% confidence)
	MaxDrawdown        sdk.Dec   `json:"max_drawdown"`         // Maximum historical drawdown
	Volatility         sdk.Dec   `json:"volatility"`           // Annualized volatility
	SharpeRatio        sdk.Dec   `json:"sharpe_ratio"`         // Risk-adjusted returns
	ConcentrationRisk  sdk.Dec   `json:"concentration_risk"`   // Concentration risk score
	LiquidityRisk      sdk.Dec   `json:"liquidity_risk"`       // Liquidity risk score
	CounterpartyRisk   sdk.Dec   `json:"counterparty_risk"`    // Counterparty risk score
	OverallRiskScore   sdk.Dec   `json:"overall_risk_score"`   // Overall risk score 0-100
	LastUpdated        time.Time `json:"last_updated"`
}

// DUSDPerformanceMetrics represents performance metrics for DUSD reserves
type DUSDPerformanceMetrics struct {
	TotalReturn        sdk.Dec   `json:"total_return"`         // Total return since inception
	AnnualizedReturn   sdk.Dec   `json:"annualized_return"`    // Annualized return
	MonthlyReturns     []sdk.Dec `json:"monthly_returns"`      // Last 12 months
	BenchmarkReturn    sdk.Dec   `json:"benchmark_return"`     // Benchmark return
	Alpha              sdk.Dec   `json:"alpha"`                // Alpha vs benchmark
	Beta               sdk.Dec   `json:"beta"`                 // Beta vs benchmark
	TrackingError      sdk.Dec   `json:"tracking_error"`       // Tracking error vs benchmark
	InformationRatio   sdk.Dec   `json:"information_ratio"`    // Information ratio
	YieldGenerated     sdk.Coin  `json:"yield_generated"`      // Total yield generated
	FeesIncurred       sdk.Coin  `json:"fees_incurred"`        // Total fees incurred
	NetPerformance     sdk.Dec   `json:"net_performance"`      // Net performance after fees
	LastUpdated        time.Time `json:"last_updated"`
}

// DUSDRebalanceRecord represents a historical rebalancing record
type DUSDRebalanceRecord struct {
	RebalanceID        string                    `json:"rebalance_id"`
	Timestamp          time.Time                 `json:"timestamp"`
	Trigger            string                    `json:"trigger"`           // DEVIATION, SCHEDULED, MANUAL, EMERGENCY
	PreRebalanceState  DUSDPortfolioState        `json:"pre_rebalance_state"`
	PostRebalanceState DUSDPortfolioState        `json:"post_rebalance_state"`
	Actions            []RebalanceAction         `json:"actions"`
	TotalCost          sdk.Coin                  `json:"total_cost"`
	ImpactOnPerformance sdk.Dec                  `json:"impact_on_performance"`
	Success            bool                      `json:"success"`
	ErrorMessage       string                    `json:"error_message,omitempty"`
}

// DUSDPortfolioState represents the state of the DUSD portfolio at a point in time
type DUSDPortfolioState struct {
	Timestamp          time.Time                 `json:"timestamp"`
	TotalValue         sdk.Coin                  `json:"total_value"`
	AssetBreakdown     map[string]sdk.Coin       `json:"asset_breakdown"`
	WeightBreakdown    map[string]sdk.Dec        `json:"weight_breakdown"`
	CurrencyExposure   map[string]sdk.Dec        `json:"currency_exposure"`
	RiskMetrics        DUSDRiskMetrics           `json:"risk_metrics"`
	LiquidityProfile   map[string]sdk.Dec        `json:"liquidity_profile"`
}

// CrossCurrencyExposure represents exposure across multiple currencies
type CrossCurrencyExposure struct {
	Currency           string    `json:"currency"`
	Exposure           sdk.Dec   `json:"exposure"`            // Percentage exposure
	ExposureAmount     sdk.Coin  `json:"exposure_amount"`     // Absolute exposure
	TargetExposure     sdk.Dec   `json:"target_exposure"`     // Target percentage
	Deviation          sdk.Dec   `json:"deviation"`           // Current vs target deviation
	HedgeRatio         sdk.Dec   `json:"hedge_ratio"`         // Hedge ratio for this currency
	LastRebalanced     time.Time `json:"last_rebalanced"`
}

// DUSDStabilityMetrics represents stability metrics for DUSD
type DUSDStabilityMetrics struct {
	PriceDivergence    sdk.Dec                   `json:"price_divergence"`     // Divergence from $1.00 target
	PriceVolatility    sdk.Dec                   `json:"price_volatility"`     // Price volatility over time
	RedemptionPressure sdk.Dec                   `json:"redemption_pressure"`  // Pressure from redemptions
	MintingDemand      sdk.Dec                   `json:"minting_demand"`       // Demand for new minting
	LiquidityDepth     sdk.Coin                  `json:"liquidity_depth"`      // Available liquidity
	ArbitrageGaps      []ArbitrageGap            `json:"arbitrage_gaps"`       // Arbitrage opportunities
	StabilityScore     sdk.Dec                   `json:"stability_score"`      // Overall stability score 0-100
	LastUpdated        time.Time                 `json:"last_updated"`
}

// ArbitrageGap represents an arbitrage opportunity
type ArbitrageGap struct {
	Source         string    `json:"source"`           // Source exchange/protocol
	Target         string    `json:"target"`           // Target exchange/protocol
	PriceGap       sdk.Dec   `json:"price_gap"`        // Price difference percentage
	PotentialProfit sdk.Coin `json:"potential_profit"` // Potential profit
	LiquiditySize  sdk.Coin  `json:"liquidity_size"`   // Available liquidity
	RiskLevel      string    `json:"risk_level"`       // LOW, MEDIUM, HIGH
	Timestamp      time.Time `json:"timestamp"`
}

// DUSDCollateralRequirement represents collateral requirements for DUSD
type DUSDCollateralRequirement struct {
	DUSDAmount             sdk.Coin  `json:"dusd_amount"`
	RequiredUSDCollateral  sdk.Coin  `json:"required_usd_collateral"`
	CurrentUSDCollateral   sdk.Coin  `json:"current_usd_collateral"`
	CollateralRatio        sdk.Dec   `json:"collateral_ratio"`
	MinCollateralRatio     sdk.Dec   `json:"min_collateral_ratio"`
	TargetCollateralRatio  sdk.Dec   `json:"target_collateral_ratio"`
	ExcessCollateral       sdk.Coin  `json:"excess_collateral"`
	DeficitCollateral      sdk.Coin  `json:"deficit_collateral"`
	HealthScore            sdk.Dec   `json:"health_score"`          // 0-100 health score
	ActionRequired         string    `json:"action_required"`       // NONE, ADD_COLLATERAL, REDUCE_EXPOSURE
	LastUpdated            time.Time `json:"last_updated"`
}

// USDTreasuryOperations represents USD treasury operations
type USDTreasuryOperations struct {
	OperationID     string                    `json:"operation_id"`
	OperationType   string                    `json:"operation_type"`   // MINT_DUSD, BURN_DUSD, REBALANCE, YIELD_COLLECTION
	Status          string                    `json:"status"`           // PENDING, IN_PROGRESS, COMPLETED, FAILED
	InitiatedBy     string                    `json:"initiated_by"`     // User or system component
	Amount          sdk.Coin                  `json:"amount"`
	USDAmount       sdk.Coin                  `json:"usd_amount"`       // USD equivalent
	Parameters      map[string]string         `json:"parameters"`       // Operation-specific parameters
	PreState        DUSDPortfolioState        `json:"pre_state"`
	PostState       DUSDPortfolioState        `json:"post_state"`
	TransactionHash string                    `json:"transaction_hash"`
	GasUsed         uint64                    `json:"gas_used"`
	Fees            sdk.Coin                  `json:"fees"`
	StartTime       time.Time                 `json:"start_time"`
	EndTime         time.Time                 `json:"end_time"`
	ErrorMessage    string                    `json:"error_message,omitempty"`
}