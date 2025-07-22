package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgMintDINR           = "mint_dinr"
	TypeMsgBurnDINR           = "burn_dinr"
	TypeMsgDepositCollateral  = "deposit_collateral"
	TypeMsgWithdrawCollateral = "withdraw_collateral"
	TypeMsgLiquidate          = "liquidate"
	TypeMsgUpdateParams       = "update_params"
)

var (
	_ sdk.Msg = &MsgMintDINR{}
	_ sdk.Msg = &MsgBurnDINR{}
	_ sdk.Msg = &MsgDepositCollateral{}
	_ sdk.Msg = &MsgWithdrawCollateral{}
	_ sdk.Msg = &MsgLiquidate{}
	_ sdk.Msg = &MsgUpdateParams{}
)

// NewMsgMintDINR creates a new MsgMintDINR instance
func NewMsgMintDINR(minter string, collateral sdk.Coin, dinrToMint sdk.Coin) *MsgMintDINR {
	return &MsgMintDINR{
		Minter:     minter,
		Collateral: collateral,
		DinrToMint: dinrToMint,
	}
}

// Route implements sdk.Msg
func (msg *MsgMintDINR) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (msg *MsgMintDINR) Type() string {
	return TypeMsgMintDINR
}

// GetSigners implements sdk.Msg
func (msg *MsgMintDINR) GetSigners() []sdk.AccAddress {
	minter, err := sdk.AccAddressFromBech32(msg.Minter)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{minter}
}

// GetSignBytes implements sdk.Msg
func (msg *MsgMintDINR) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements sdk.Msg
func (msg *MsgMintDINR) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Minter)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid minter address (%s)", err)
	}

	if !msg.Collateral.IsValid() || msg.Collateral.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid collateral amount")
	}

	if !msg.DinrToMint.IsValid() || msg.DinrToMint.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid DINR mint amount")
	}

	if msg.DinrToMint.Denom != DINRDenom {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "mint denom must be %s", DINRDenom)
	}

	return nil
}

// NewMsgBurnDINR creates a new MsgBurnDINR instance
func NewMsgBurnDINR(burner string, dinrToBurn sdk.Coin, collateralDenom string) *MsgBurnDINR {
	return &MsgBurnDINR{
		Burner:          burner,
		DinrToBurn:      dinrToBurn,
		CollateralDenom: collateralDenom,
	}
}

// Route implements sdk.Msg
func (msg *MsgBurnDINR) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (msg *MsgBurnDINR) Type() string {
	return TypeMsgBurnDINR
}

// GetSigners implements sdk.Msg
func (msg *MsgBurnDINR) GetSigners() []sdk.AccAddress {
	burner, err := sdk.AccAddressFromBech32(msg.Burner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{burner}
}

// GetSignBytes implements sdk.Msg
func (msg *MsgBurnDINR) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements sdk.Msg
func (msg *MsgBurnDINR) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Burner)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid burner address (%s)", err)
	}

	if !msg.DinrToBurn.IsValid() || msg.DinrToBurn.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid DINR burn amount")
	}

	if msg.DinrToBurn.Denom != DINRDenom {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "burn denom must be %s", DINRDenom)
	}

	if msg.CollateralDenom == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "collateral denom cannot be empty")
	}

	return nil
}

// NewMsgDepositCollateral creates a new MsgDepositCollateral instance
func NewMsgDepositCollateral(depositor string, collateral sdk.Coin) *MsgDepositCollateral {
	return &MsgDepositCollateral{
		Depositor:  depositor,
		Collateral: collateral,
	}
}

// Route implements sdk.Msg
func (msg *MsgDepositCollateral) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (msg *MsgDepositCollateral) Type() string {
	return TypeMsgDepositCollateral
}

// GetSigners implements sdk.Msg
func (msg *MsgDepositCollateral) GetSigners() []sdk.AccAddress {
	depositor, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{depositor}
}

// GetSignBytes implements sdk.Msg
func (msg *MsgDepositCollateral) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements sdk.Msg
func (msg *MsgDepositCollateral) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid depositor address (%s)", err)
	}

	if !msg.Collateral.IsValid() || msg.Collateral.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid collateral amount")
	}

	return nil
}

// NewMsgWithdrawCollateral creates a new MsgWithdrawCollateral instance
func NewMsgWithdrawCollateral(withdrawer string, collateral sdk.Coin) *MsgWithdrawCollateral {
	return &MsgWithdrawCollateral{
		Withdrawer: withdrawer,
		Collateral: collateral,
	}
}

// Route implements sdk.Msg
func (msg *MsgWithdrawCollateral) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (msg *MsgWithdrawCollateral) Type() string {
	return TypeMsgWithdrawCollateral
}

// GetSigners implements sdk.Msg
func (msg *MsgWithdrawCollateral) GetSigners() []sdk.AccAddress {
	withdrawer, err := sdk.AccAddressFromBech32(msg.Withdrawer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{withdrawer}
}

// GetSignBytes implements sdk.Msg
func (msg *MsgWithdrawCollateral) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements sdk.Msg
func (msg *MsgWithdrawCollateral) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Withdrawer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid withdrawer address (%s)", err)
	}

	if !msg.Collateral.IsValid() || msg.Collateral.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid collateral amount")
	}

	return nil
}

// NewMsgLiquidate creates a new MsgLiquidate instance
func NewMsgLiquidate(liquidator string, user string, dinrToCover sdk.Coin) *MsgLiquidate {
	return &MsgLiquidate{
		Liquidator:  liquidator,
		User:        user,
		DinrToCover: dinrToCover,
	}
}

// Route implements sdk.Msg
func (msg *MsgLiquidate) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (msg *MsgLiquidate) Type() string {
	return TypeMsgLiquidate
}

// GetSigners implements sdk.Msg
func (msg *MsgLiquidate) GetSigners() []sdk.AccAddress {
	liquidator, err := sdk.AccAddressFromBech32(msg.Liquidator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{liquidator}
}

// GetSignBytes implements sdk.Msg
func (msg *MsgLiquidate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements sdk.Msg
func (msg *MsgLiquidate) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Liquidator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid liquidator address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.User)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid user address (%s)", err)
	}

	if !msg.DinrToCover.IsValid() || msg.DinrToCover.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid DINR cover amount")
	}

	if msg.DinrToCover.Denom != DINRDenom {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "cover denom must be %s", DINRDenom)
	}

	return nil
}

// NewMsgUpdateParams creates a new MsgUpdateParams instance
func NewMsgUpdateParams(authority string, params Params) *MsgUpdateParams {
	return &MsgUpdateParams{
		Authority: authority,
		Params:    params,
	}
}

// Route implements sdk.Msg
func (msg *MsgUpdateParams) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (msg *MsgUpdateParams) Type() string {
	return TypeMsgUpdateParams
}

// GetSigners implements sdk.Msg
func (msg *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

// GetSignBytes implements sdk.Msg
func (msg *MsgUpdateParams) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements sdk.Msg
func (msg *MsgUpdateParams) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}

	return msg.Params.Validate()
}