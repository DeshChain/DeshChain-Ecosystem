package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MultiCurrencyStats represents statistics for multi-currency trade finance
type MultiCurrencyStats struct {
	TotalLCs               int64                 `json:"total_lcs"`
	TotalVolume           sdk.Coin             `json:"total_volume"`
	TotalSavings          sdk.Coin             `json:"total_savings"`
	AverageProcessingTime time.Duration        `json:"average_processing_time"`
	CurrencyBreakdown     map[string]int64     `json:"currency_breakdown"`
	StablecoinBreakdown   map[string]int64     `json:"stablecoin_breakdown"`
}

// Enhanced LC key prefix
var EnhancedLCKeyPrefix = []byte{0x10}

// GetEnhancedLCKey returns the store key for enhanced LC data
func GetEnhancedLCKey(lcID string) []byte {
	return append(EnhancedLCKeyPrefix, []byte(lcID)...)
}

// MultiCurrencyLCRequest represents a request for multi-currency LC processing
type MultiCurrencyLCRequest struct {
	LcId               string   `json:"lc_id"`
	OriginalCurrency   string   `json:"original_currency"`
	OriginalAmount     sdk.Coin `json:"original_amount"`
	SettlementCurrency string   `json:"settlement_currency"`
	ApplicantId        string   `json:"applicant_id"`
	BeneficiaryId      string   `json:"beneficiary_id"`
	IssuingBankId      string   `json:"issuing_bank_id"`
}

// MultiCurrencyLCResponse represents the response for multi-currency LC processing
type MultiCurrencyLCResponse struct {
	LcId               string            `json:"lc_id"`
	OriginalAmount     sdk.Coin         `json:"original_amount"`
	SettlementAmount   sdk.Coin         `json:"settlement_amount"`
	LocalAmount        sdk.Coin         `json:"local_amount"`
	ExchangeRates      map[string]sdk.Dec `json:"exchange_rates"`
	TotalFees          sdk.Coin         `json:"total_fees"`
	TotalSavings       sdk.Coin         `json:"total_savings"`
	ProcessingTime     time.Duration    `json:"processing_time"`
	CurrencyRoute      []string         `json:"currency_route"`
	Status             string           `json:"status"`
}

// CurrencyConversionStep represents a step in currency conversion
type CurrencyConversionStep struct {
	FromCurrency string   `json:"from_currency"`
	ToCurrency   string   `json:"to_currency"`
	FromAmount   sdk.Coin `json:"from_amount"`
	ToAmount     sdk.Coin `json:"to_amount"`
	ExchangeRate sdk.Dec  `json:"exchange_rate"`
	Fee          sdk.Coin `json:"fee"`
	Timestamp    time.Time `json:"timestamp"`
}

// TradeFinanceCorridorStats represents statistics for specific trade corridors
type TradeFinanceCorridorStats struct {
	Corridor           string        `json:"corridor"`           // e.g., "USD-INR", "EUR-INR"
	TotalVolume        sdk.Coin      `json:"total_volume"`
	TotalTransactions  int64         `json:"total_transactions"`
	AverageAmount      sdk.Coin      `json:"average_amount"`
	TotalSavings       sdk.Coin      `json:"total_savings"`
	AverageProcessingTime time.Duration `json:"average_processing_time"`
	TraditionalCost    sdk.Coin      `json:"traditional_cost"`
	DeshChainCost      sdk.Coin      `json:"deshchain_cost"`
	CostSavingPercent  sdk.Dec       `json:"cost_saving_percent"`
}

// MultiCurrencyLCStatus represents the status of multi-currency LC processing
type MultiCurrencyLCStatus struct {
	LcId                string                    `json:"lc_id"`
	Status              string                    `json:"status"` // "initiated", "processing", "completed", "failed"
	CurrentStep         string                    `json:"current_step"`
	ConversionSteps     []CurrencyConversionStep  `json:"conversion_steps"`
	StartTime           time.Time                 `json:"start_time"`
	CompletionTime      time.Time                 `json:"completion_time"`
	ErrorMessage        string                    `json:"error_message,omitempty"`
	Metadata            map[string]string         `json:"metadata"`
}