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

const (
	TypeMsgDonate       = "donate"
	TypeMsgUpdateParams = "update_params"
)

var (
	_ sdk.Msg = &MsgDonate{}
	_ sdk.Msg = &MsgUpdateParams{}
)

// MsgDonate defines a message to make a donation to an NGO
type MsgDonate struct {
	Donor        string    `protobuf:"bytes,1,opt,name=donor,proto3" json:"donor,omitempty"`
	NgoWalletId  uint64    `protobuf:"varint,2,opt,name=ngo_wallet_id,json=ngoWalletId,proto3" json:"ngo_wallet_id,omitempty"`
	Amount       sdk.Coins `protobuf:"bytes,3,rep,name=amount,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"amount"`
	Purpose      string    `protobuf:"bytes,4,opt,name=purpose,proto3" json:"purpose,omitempty"`
	IsAnonymous  bool      `protobuf:"varint,5,opt,name=is_anonymous,json=isAnonymous,proto3" json:"is_anonymous,omitempty"`
	CampaignId   uint64    `protobuf:"varint,6,opt,name=campaign_id,json=campaignId,proto3" json:"campaign_id,omitempty"`
}

// Route returns the name of the module
func (msg MsgDonate) Route() string { return RouterKey }

// Type returns the the action
func (msg MsgDonate) Type() string { return TypeMsgDonate }

// GetSigners returns the expected signers for a MsgDonate message
func (msg *MsgDonate) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Donor)
	return []sdk.AccAddress{addr}
}

// GetSignBytes returns the bytes to sign for a MsgDonate message
func (msg *MsgDonate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic performs basic validation
func (msg *MsgDonate) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Donor)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid donor address: %s", err)
	}

	if msg.NgoWalletId == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("NGO wallet ID cannot be zero")
	}

	if msg.Amount.IsZero() {
		return sdkerrors.ErrInvalidRequest.Wrap("donation amount cannot be zero")
	}

	if !msg.Amount.IsValid() {
		return sdkerrors.ErrInvalidRequest.Wrap("invalid donation amount")
	}

	if len(msg.Purpose) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("donation purpose cannot be empty")
	}

	return nil
}

// MsgDonateResponse defines the response for MsgDonate
type MsgDonateResponse struct {
	DonationId  uint64 `protobuf:"varint,1,opt,name=donation_id,json=donationId,proto3" json:"donation_id,omitempty"`
	ReceiptHash string `protobuf:"bytes,2,opt,name=receipt_hash,json=receiptHash,proto3" json:"receipt_hash,omitempty"`
}

// MsgUpdateParams defines a message to update module parameters
type MsgUpdateParams struct {
	Authority string `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority,omitempty"`
	Params    Params `protobuf:"bytes,2,opt,name=params,proto3" json:"params"`
}

// Route returns the name of the module
func (msg MsgUpdateParams) Route() string { return RouterKey }

// Type returns the the action
func (msg MsgUpdateParams) Type() string { return TypeMsgUpdateParams }

// GetSigners returns the expected signers for a MsgUpdateParams message
func (msg *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

// GetSignBytes returns the bytes to sign for a MsgUpdateParams message
func (msg *MsgUpdateParams) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic performs basic validation
func (msg *MsgUpdateParams) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid authority address: %s", err)
	}

	return msg.Params.Validate()
}

// MsgUpdateParamsResponse defines the response for MsgUpdateParams
type MsgUpdateParamsResponse struct{}