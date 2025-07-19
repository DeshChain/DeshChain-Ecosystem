/*
Copyright 2024 DeshChain Foundation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Message types
const (
	TypeMsgCreateMoneyOrder     = "create_money_order"
	TypeMsgCreateFixedRatePool  = "create_fixed_rate_pool"
	TypeMsgCreateVillagePool    = "create_village_pool"
	TypeMsgAddLiquidity         = "add_liquidity"
	TypeMsgRemoveLiquidity      = "remove_liquidity"
	TypeMsgSwapExactAmountIn    = "swap_exact_amount_in"
	TypeMsgSwapExactAmountOut   = "swap_exact_amount_out"
	TypeMsgJoinVillagePool      = "join_village_pool"
	TypeMsgClaimRewards         = "claim_rewards"
	TypeMsgUpdatePoolParams     = "update_pool_params"
)

var (
	_ sdk.Msg = &MsgCreateMoneyOrder{}
	_ sdk.Msg = &MsgCreateFixedRatePool{}
	_ sdk.Msg = &MsgCreateVillagePool{}
	_ sdk.Msg = &MsgAddLiquidity{}
	_ sdk.Msg = &MsgRemoveLiquidity{}
	_ sdk.Msg = &MsgSwapExactAmountIn{}
	_ sdk.Msg = &MsgSwapExactAmountOut{}
	_ sdk.Msg = &MsgJoinVillagePool{}
	_ sdk.Msg = &MsgClaimRewards{}
	_ sdk.Msg = &MsgUpdatePoolParams{}
)

// MsgCreateMoneyOrder - Simple UPI-style money transfer
type MsgCreateMoneyOrder struct {
	Sender          string   `json:"sender" yaml:"sender"`
	ReceiverUPI     string   `json:"receiver_upi" yaml:"receiver_upi"`     // name@deshchain or address
	Amount          sdk.Coin `json:"amount" yaml:"amount"`
	Note            string   `json:"note" yaml:"note"`                     // Personal message
	OrderType       string   `json:"order_type" yaml:"order_type"`         // "instant", "normal", "scheduled"
	ScheduledTime   int64    `json:"scheduled_time" yaml:"scheduled_time"` // Unix timestamp if scheduled
}

// NewMsgCreateMoneyOrder creates a new money order message
func NewMsgCreateMoneyOrder(
	sender sdk.AccAddress,
	receiverUPI string,
	amount sdk.Coin,
	note string,
	orderType string,
) *MsgCreateMoneyOrder {
	return &MsgCreateMoneyOrder{
		Sender:      sender.String(),
		ReceiverUPI: receiverUPI,
		Amount:      amount,
		Note:        note,
		OrderType:   orderType,
	}
}

func (msg *MsgCreateMoneyOrder) Route() string { return RouterKey }
func (msg *MsgCreateMoneyOrder) Type() string  { return TypeMsgCreateMoneyOrder }

func (msg *MsgCreateMoneyOrder) GetSigners() []sdk.AccAddress {
	sender, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{sender}
}

func (msg *MsgCreateMoneyOrder) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateMoneyOrder) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address: %s", err)
	}
	
	if msg.ReceiverUPI == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "receiver UPI cannot be empty")
	}
	
	if !msg.Amount.IsValid() || msg.Amount.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid amount")
	}
	
	if msg.OrderType != "instant" && msg.OrderType != "normal" && msg.OrderType != "scheduled" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid order type")
	}
	
	if msg.OrderType == "scheduled" && msg.ScheduledTime <= 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "scheduled time must be provided for scheduled orders")
	}
	
	return nil
}

// MsgCreateFixedRatePool - Create a fixed exchange rate pool
type MsgCreateFixedRatePool struct {
	Creator          string   `json:"creator" yaml:"creator"`
	Token0Denom      string   `json:"token0_denom" yaml:"token0_denom"`
	Token1Denom      string   `json:"token1_denom" yaml:"token1_denom"`
	ExchangeRate     sdk.Dec  `json:"exchange_rate" yaml:"exchange_rate"`
	InitialLiquidity sdk.Coins `json:"initial_liquidity" yaml:"initial_liquidity"`
	Description      string   `json:"description" yaml:"description"`
	SupportedRegions []string `json:"supported_regions" yaml:"supported_regions"`
}

func (msg *MsgCreateFixedRatePool) Route() string { return RouterKey }
func (msg *MsgCreateFixedRatePool) Type() string  { return TypeMsgCreateFixedRatePool }

func (msg *MsgCreateFixedRatePool) GetSigners() []sdk.AccAddress {
	creator, _ := sdk.AccAddressFromBech32(msg.Creator)
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateFixedRatePool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateFixedRatePool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address: %s", err)
	}
	
	if msg.Token0Denom == "" || msg.Token1Denom == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "token denoms cannot be empty")
	}
	
	if msg.Token0Denom == msg.Token1Denom {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "token denoms must be different")
	}
	
	if msg.ExchangeRate.IsZero() || msg.ExchangeRate.IsNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid exchange rate")
	}
	
	if !msg.InitialLiquidity.IsValid() || msg.InitialLiquidity.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid initial liquidity")
	}
	
	return nil
}

// MsgCreateVillagePool - Create a community-managed pool
type MsgCreateVillagePool struct {
	PanchayatHead    string   `json:"panchayat_head" yaml:"panchayat_head"`
	VillageName      string   `json:"village_name" yaml:"village_name"`
	PostalCode       string   `json:"postal_code" yaml:"postal_code"`
	StateCode        string   `json:"state_code" yaml:"state_code"`
	DistrictCode     string   `json:"district_code" yaml:"district_code"`
	InitialLiquidity sdk.Coins `json:"initial_liquidity" yaml:"initial_liquidity"`
	LocalValidators  []string `json:"local_validators" yaml:"local_validators"`
}

func (msg *MsgCreateVillagePool) Route() string { return RouterKey }
func (msg *MsgCreateVillagePool) Type() string  { return TypeMsgCreateVillagePool }

func (msg *MsgCreateVillagePool) GetSigners() []sdk.AccAddress {
	head, _ := sdk.AccAddressFromBech32(msg.PanchayatHead)
	return []sdk.AccAddress{head}
}

func (msg *MsgCreateVillagePool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateVillagePool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.PanchayatHead)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid panchayat head address: %s", err)
	}
	
	if msg.VillageName == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "village name cannot be empty")
	}
	
	if len(msg.PostalCode) != 6 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid postal code")
	}
	
	if !msg.InitialLiquidity.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid initial liquidity")
	}
	
	if len(msg.LocalValidators) < 2 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "at least 2 local validators required")
	}
	
	for _, val := range msg.LocalValidators {
		_, err := sdk.ValAddressFromBech32(val)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid validator address: %s", err)
		}
	}
	
	return nil
}

// MsgAddLiquidity - Add liquidity to any pool
type MsgAddLiquidity struct {
	Depositor     string    `json:"depositor" yaml:"depositor"`
	PoolId        uint64    `json:"pool_id" yaml:"pool_id"`
	TokenAmounts  sdk.Coins `json:"token_amounts" yaml:"token_amounts"`
	ShareOutMin   sdk.Int   `json:"share_out_min" yaml:"share_out_min"`
}

func (msg *MsgAddLiquidity) Route() string { return RouterKey }
func (msg *MsgAddLiquidity) Type() string  { return TypeMsgAddLiquidity }

func (msg *MsgAddLiquidity) GetSigners() []sdk.AccAddress {
	depositor, _ := sdk.AccAddressFromBech32(msg.Depositor)
	return []sdk.AccAddress{depositor}
}

func (msg *MsgAddLiquidity) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAddLiquidity) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid depositor address: %s", err)
	}
	
	if msg.PoolId == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid pool ID")
	}
	
	if !msg.TokenAmounts.IsValid() || msg.TokenAmounts.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid token amounts")
	}
	
	if msg.ShareOutMin.IsNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "share out minimum cannot be negative")
	}
	
	return nil
}

// MsgRemoveLiquidity - Remove liquidity from a pool
type MsgRemoveLiquidity struct {
	Withdrawer      string  `json:"withdrawer" yaml:"withdrawer"`
	PoolId          uint64  `json:"pool_id" yaml:"pool_id"`
	ShareAmount     sdk.Int `json:"share_amount" yaml:"share_amount"`
	TokenOutMins    sdk.Coins `json:"token_out_mins" yaml:"token_out_mins"`
}

func (msg *MsgRemoveLiquidity) Route() string { return RouterKey }
func (msg *MsgRemoveLiquidity) Type() string  { return TypeMsgRemoveLiquidity }

func (msg *MsgRemoveLiquidity) GetSigners() []sdk.AccAddress {
	withdrawer, _ := sdk.AccAddressFromBech32(msg.Withdrawer)
	return []sdk.AccAddress{withdrawer}
}

func (msg *MsgRemoveLiquidity) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRemoveLiquidity) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Withdrawer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid withdrawer address: %s", err)
	}
	
	if msg.PoolId == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid pool ID")
	}
	
	if msg.ShareAmount.IsZero() || msg.ShareAmount.IsNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid share amount")
	}
	
	if !msg.TokenOutMins.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid minimum token outputs")
	}
	
	return nil
}

// MsgSwapExactAmountIn - Swap with exact input amount
type MsgSwapExactAmountIn struct {
	Sender       string   `json:"sender" yaml:"sender"`
	PoolId       uint64   `json:"pool_id" yaml:"pool_id"`
	TokenIn      sdk.Coin `json:"token_in" yaml:"token_in"`
	TokenOutDenom string  `json:"token_out_denom" yaml:"token_out_denom"`
	TokenOutMin  sdk.Int  `json:"token_out_min" yaml:"token_out_min"`
}

func (msg *MsgSwapExactAmountIn) Route() string { return RouterKey }
func (msg *MsgSwapExactAmountIn) Type() string  { return TypeMsgSwapExactAmountIn }

func (msg *MsgSwapExactAmountIn) GetSigners() []sdk.AccAddress {
	sender, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{sender}
}

func (msg *MsgSwapExactAmountIn) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSwapExactAmountIn) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address: %s", err)
	}
	
	if msg.PoolId == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid pool ID")
	}
	
	if !msg.TokenIn.IsValid() || msg.TokenIn.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid token in")
	}
	
	if msg.TokenOutDenom == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "token out denom cannot be empty")
	}
	
	if msg.TokenIn.Denom == msg.TokenOutDenom {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "cannot swap same token")
	}
	
	if msg.TokenOutMin.IsNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "token out minimum cannot be negative")
	}
	
	return nil
}

// MsgSwapExactAmountOut - Swap with exact output amount
type MsgSwapExactAmountOut struct {
	Sender       string   `json:"sender" yaml:"sender"`
	PoolId       uint64   `json:"pool_id" yaml:"pool_id"`
	TokenInDenom string   `json:"token_in_denom" yaml:"token_in_denom"`
	TokenInMax   sdk.Int  `json:"token_in_max" yaml:"token_in_max"`
	TokenOut     sdk.Coin `json:"token_out" yaml:"token_out"`
}

func (msg *MsgSwapExactAmountOut) Route() string { return RouterKey }
func (msg *MsgSwapExactAmountOut) Type() string  { return TypeMsgSwapExactAmountOut }

func (msg *MsgSwapExactAmountOut) GetSigners() []sdk.AccAddress {
	sender, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{sender}
}

func (msg *MsgSwapExactAmountOut) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSwapExactAmountOut) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address: %s", err)
	}
	
	if msg.PoolId == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid pool ID")
	}
	
	if msg.TokenInDenom == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "token in denom cannot be empty")
	}
	
	if !msg.TokenOut.IsValid() || msg.TokenOut.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid token out")
	}
	
	if msg.TokenInDenom == msg.TokenOut.Denom {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "cannot swap same token")
	}
	
	if msg.TokenInMax.IsNegative() || msg.TokenInMax.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid token in maximum")
	}
	
	return nil
}

// MsgJoinVillagePool - Join a village pool as a member
type MsgJoinVillagePool struct {
	Member         string   `json:"member" yaml:"member"`
	PoolId         uint64   `json:"pool_id" yaml:"pool_id"`
	InitialDeposit sdk.Coins `json:"initial_deposit" yaml:"initial_deposit"`
	LocalName      string   `json:"local_name" yaml:"local_name"`
	MobileNumber   string   `json:"mobile_number" yaml:"mobile_number"`
}

func (msg *MsgJoinVillagePool) Route() string { return RouterKey }
func (msg *MsgJoinVillagePool) Type() string  { return TypeMsgJoinVillagePool }

func (msg *MsgJoinVillagePool) GetSigners() []sdk.AccAddress {
	member, _ := sdk.AccAddressFromBech32(msg.Member)
	return []sdk.AccAddress{member}
}

func (msg *MsgJoinVillagePool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgJoinVillagePool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Member)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid member address: %s", err)
	}
	
	if msg.PoolId == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid pool ID")
	}
	
	if msg.LocalName == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "local name cannot be empty")
	}
	
	if msg.MobileNumber == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "mobile number cannot be empty")
	}
	
	return nil
}

// MsgClaimRewards - Claim accumulated rewards
type MsgClaimRewards struct {
	Claimer  string `json:"claimer" yaml:"claimer"`
	PoolId   uint64 `json:"pool_id" yaml:"pool_id"`
}

func (msg *MsgClaimRewards) Route() string { return RouterKey }
func (msg *MsgClaimRewards) Type() string  { return TypeMsgClaimRewards }

func (msg *MsgClaimRewards) GetSigners() []sdk.AccAddress {
	claimer, _ := sdk.AccAddressFromBech32(msg.Claimer)
	return []sdk.AccAddress{claimer}
}

func (msg *MsgClaimRewards) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgClaimRewards) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Claimer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid claimer address: %s", err)
	}
	
	if msg.PoolId == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid pool ID")
	}
	
	return nil
}

// MsgUpdatePoolParams - Update pool parameters (governance)
type MsgUpdatePoolParams struct {
	Authority    string  `json:"authority" yaml:"authority"`
	PoolId       uint64  `json:"pool_id" yaml:"pool_id"`
	BaseFee      sdk.Dec `json:"base_fee" yaml:"base_fee"`
	Active       bool    `json:"active" yaml:"active"`
}

func (msg *MsgUpdatePoolParams) Route() string { return RouterKey }
func (msg *MsgUpdatePoolParams) Type() string  { return TypeMsgUpdatePoolParams }

func (msg *MsgUpdatePoolParams) GetSigners() []sdk.AccAddress {
	authority, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{authority}
}

func (msg *MsgUpdatePoolParams) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdatePoolParams) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address: %s", err)
	}
	
	if msg.PoolId == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid pool ID")
	}
	
	if msg.BaseFee.IsNegative() || msg.BaseFee.GT(sdk.NewDecWithPrec(10, 2)) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid base fee")
	}
	
	return nil
}