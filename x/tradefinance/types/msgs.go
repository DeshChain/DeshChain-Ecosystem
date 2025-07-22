package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgRegisterParty = "register_party"

var _ sdk.Msg = &MsgRegisterParty{}

func (msg *MsgRegisterParty) Route() string {
	return RouterKey
}

func (msg *MsgRegisterParty) Type() string {
	return TypeMsgRegisterParty
}

func (msg *MsgRegisterParty) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRegisterParty) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRegisterParty) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.PartyType == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "party type cannot be empty")
	}

	if msg.Name == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "party name cannot be empty")
	}

	return nil
}

const TypeMsgIssueLc = "issue_lc"

var _ sdk.Msg = &MsgIssueLc{}

func (msg *MsgIssueLc) Route() string {
	return RouterKey
}

func (msg *MsgIssueLc) Type() string {
	return TypeMsgIssueLc
}

func (msg *MsgIssueLc) GetSigners() []sdk.AccAddress {
	issuer, err := sdk.AccAddressFromBech32(msg.IssuingBank)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{issuer}
}

func (msg *MsgIssueLc) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgIssueLc) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.IssuingBank)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid issuing bank address (%s)", err)
	}

	if msg.ApplicantId == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "applicant ID cannot be empty")
	}

	if msg.BeneficiaryId == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "beneficiary ID cannot be empty")
	}

	if !msg.Amount.IsValid() || msg.Amount.Amount.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid LC amount")
	}

	if !msg.Collateral.IsValid() || msg.Collateral.Amount.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid collateral amount")
	}

	return nil
}

const TypeMsgAcceptLc = "accept_lc"

var _ sdk.Msg = &MsgAcceptLc{}

func (msg *MsgAcceptLc) Route() string {
	return RouterKey
}

func (msg *MsgAcceptLc) Type() string {
	return TypeMsgAcceptLc
}

func (msg *MsgAcceptLc) GetSigners() []sdk.AccAddress {
	beneficiary, err := sdk.AccAddressFromBech32(msg.Beneficiary)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{beneficiary}
}

func (msg *MsgAcceptLc) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAcceptLc) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Beneficiary)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid beneficiary address (%s)", err)
	}

	if msg.LcId == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "LC ID cannot be empty")
	}

	return nil
}

// Add similar implementations for other message types...

const TypeMsgUpdateParams = "update_params"

var _ sdk.Msg = &MsgUpdateParams{}

func (msg *MsgUpdateParams) Route() string {
	return RouterKey
}

func (msg *MsgUpdateParams) Type() string {
	return TypeMsgUpdateParams
}

func (msg *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgUpdateParams) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateParams) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}

	return msg.Params.Validate()
}