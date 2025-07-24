package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TaxConfig represents the tax module configuration
type TaxConfig struct {
	Enabled                  bool     `json:"enabled"`
	ExemptMessages          []string `json:"exempt_messages"`
	ExemptAddresses         []string `json:"exempt_addresses"`
	OptimizationEnabled     bool     `json:"optimization_enabled"`
	DefaultOptimizationMode string   `json:"default_optimization_mode"`
}

// NewDefaultTaxConfig creates a default tax configuration
func NewDefaultTaxConfig() *TaxConfig {
	return &TaxConfig{
		Enabled:                  true,
		ExemptMessages:          []string{"/cosmos.gov.v1beta1.MsgVote", "/cosmos.gov.v1beta1.MsgDeposit"},
		ExemptAddresses:         []string{},
		OptimizationEnabled:     true,
		DefaultOptimizationMode: "automatic",
	}
}

// UserTaxProfile represents a user's tax profile for optimization
type UserTaxProfile struct {
	Address                 string   `json:"address"`
	PatriotismScore        int32    `json:"patriotism_score"`
	CulturalEngagementScore int32    `json:"cultural_engagement_score"`
	TransactionVolume      sdk.Coin `json:"transaction_volume"`
	TotalTaxPaid           sdk.Coin `json:"total_tax_paid"`
	OptimizationPreference string   `json:"optimization_preference"`
	CreatedAt              int64    `json:"created_at"`
	UpdatedAt              int64    `json:"updated_at"`
}

// TaxTransaction represents a tax transaction record
type TaxTransaction struct {
	TransactionID       string   `json:"transaction_id"`
	UserAddress        string   `json:"user_address"`
	TransactionAmount  sdk.Coin `json:"transaction_amount"`
	TaxAmount          sdk.Coin `json:"tax_amount"`
	EffectiveRate      string   `json:"effective_rate"`
	AppliedDiscounts   []string `json:"applied_discounts"`
	OptimizationApplied bool     `json:"optimization_applied"`
	DonationFlag       bool     `json:"donation_flag"`
	Timestamp          int64    `json:"timestamp"`
}

// Compliance levels
const (
	ComplianceLevelBasic    = "BASIC"
	ComplianceLevelStandard = "STANDARD"
	ComplianceLevelAdvanced = "ADVANCED"
)