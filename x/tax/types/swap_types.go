package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SwapRequest represents a request to swap tokens for NAMO fee payment
type SwapRequest struct {
	UserAddr  sdk.AccAddress `json:"user_addr"`
	FromCoin  sdk.Coin       `json:"from_coin"`
	ToCoin    string         `json:"to_coin"`
	Inclusive bool           `json:"inclusive"`
}

// SwapResult represents the result of a token swap
type SwapResult struct {
	UserAddr   sdk.AccAddress `json:"user_addr"`
	InputCoin  sdk.Coin       `json:"input_coin"`
	OutputCoin sdk.Coin       `json:"output_coin"`
	SwapRate   sdk.Dec        `json:"swap_rate"`
	Success    bool           `json:"success"`
	Error      string         `json:"error,omitempty"`
}

// FeeOption represents user's choice for fee payment
type FeeOption struct {
	Inclusive bool   `json:"inclusive"` // true: deduct from amount, false: add on top
	AutoSwap  bool   `json:"auto_swap"` // true: auto-swap to NAMO
	PayDenom  string `json:"pay_denom"` // denomination user wants to pay in
}