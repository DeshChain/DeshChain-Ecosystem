package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Event types for tiered fees
const (
	EventTypeTieredFeeCalculated = "tiered_fee_calculated"
	
	// Additional keys for tiered fee structure
	KeyTieredFeeStructure = []byte("tiered_fee_structure")
)

// FeeInfo contains information about calculated fees
type FeeInfo struct {
	Amount     sdk.Int `json:"amount"`
	FeeRate    sdk.Dec `json:"fee_rate"`
	Fee        sdk.Int `json:"fee"`
	TierIndex  int     `json:"tier_index"`
	MinFee     sdk.Int `json:"min_fee"`
	HasCap     bool    `json:"has_cap"`
}