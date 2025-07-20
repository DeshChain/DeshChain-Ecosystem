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
	TypeMsgRegisterDhanPataAddress        = "register_dhanpata_address"
	TypeMsgCreateKshetraCoin             = "create_kshetra_coin"
	TypeMsgRegisterEnhancedMitra         = "register_enhanced_mitra"
	TypeMsgProcessMoneyOrderWithDhanPata = "process_money_order_with_dhanpata"
	TypeMsgUpdateDhanPataMetadata        = "update_dhanpata_metadata"
)

var (
	_ sdk.Msg = &MsgRegisterDhanPataAddress{}
	_ sdk.Msg = &MsgCreateKshetraCoin{}
	_ sdk.Msg = &MsgRegisterEnhancedMitra{}
	_ sdk.Msg = &MsgProcessMoneyOrderWithDhanPata{}
	_ sdk.Msg = &MsgUpdateDhanPataMetadata{}
)

// MsgRegisterDhanPataAddress defines a message to register a DhanPata virtual address
type MsgRegisterDhanPataAddress struct {
	Sender          string   `json:"sender" yaml:"sender"`
	Name            string   `json:"name" yaml:"name"`                       // username@dhan
	AddressType     string   `json:"address_type" yaml:"address_type"`       // personal, business, service
	DisplayName     string   `json:"display_name" yaml:"display_name"`
	Description     string   `json:"description" yaml:"description"`
	ProfileImageUrl string   `json:"profile_image_url" yaml:"profile_image_url"`
	Tags            []string `json:"tags" yaml:"tags"`
}

// NewMsgRegisterDhanPataAddress creates a new MsgRegisterDhanPataAddress instance
func NewMsgRegisterDhanPataAddress(
	sender sdk.AccAddress,
	name, addressType, displayName, description, profileImageUrl string,
	tags []string,
) *MsgRegisterDhanPataAddress {
	return &MsgRegisterDhanPataAddress{
		Sender:          sender.String(),
		Name:            name,
		AddressType:     addressType,
		DisplayName:     displayName,
		Description:     description,
		ProfileImageUrl: profileImageUrl,
		Tags:            tags,
	}
}

func (msg *MsgRegisterDhanPataAddress) Route() string { return RouterKey }
func (msg *MsgRegisterDhanPataAddress) Type() string  { return TypeMsgRegisterDhanPataAddress }

func (msg *MsgRegisterDhanPataAddress) GetSigners() []sdk.AccAddress {
	sender, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{sender}
}

func (msg *MsgRegisterDhanPataAddress) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRegisterDhanPataAddress) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address: %s", err)
	}

	if err := ValidateDhanPataName(msg.Name); err != nil {
		return err
	}

	if msg.AddressType != "personal" && msg.AddressType != "business" && msg.AddressType != "service" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid address type")
	}

	return nil
}

// MsgCreateKshetraCoin defines a message to create a pincode-based memecoin
type MsgCreateKshetraCoin struct {
	Creator        string   `json:"creator" yaml:"creator"`
	Pincode        string   `json:"pincode" yaml:"pincode"`
	CoinName       string   `json:"coin_name" yaml:"coin_name"`
	CoinPrefix     string   `json:"coin_prefix" yaml:"coin_prefix"`         // Will be combined with pincode
	TotalSupply    sdk.Int  `json:"total_supply" yaml:"total_supply"`
	NgoBeneficiary string   `json:"ngo_beneficiary" yaml:"ngo_beneficiary"`
	Description    string   `json:"description" yaml:"description"`
	LocalLandmarks []string `json:"local_landmarks" yaml:"local_landmarks"`
}

// NewMsgCreateKshetraCoin creates a new MsgCreateKshetraCoin instance
func NewMsgCreateKshetraCoin(
	creator sdk.AccAddress,
	pincode, coinName, coinPrefix string,
	totalSupply sdk.Int,
	ngoBeneficiary, description string,
	localLandmarks []string,
) *MsgCreateKshetraCoin {
	return &MsgCreateKshetraCoin{
		Creator:        creator.String(),
		Pincode:        pincode,
		CoinName:       coinName,
		CoinPrefix:     coinPrefix,
		TotalSupply:    totalSupply,
		NgoBeneficiary: ngoBeneficiary,
		Description:    description,
		LocalLandmarks: localLandmarks,
	}
}

func (msg *MsgCreateKshetraCoin) Route() string { return RouterKey }
func (msg *MsgCreateKshetraCoin) Type() string  { return TypeMsgCreateKshetraCoin }

func (msg *MsgCreateKshetraCoin) GetSigners() []sdk.AccAddress {
	creator, _ := sdk.AccAddressFromBech32(msg.Creator)
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateKshetraCoin) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateKshetraCoin) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address: %s", err)
	}

	if err := ValidatePincode(msg.Pincode); err != nil {
		return err
	}

	if msg.CoinName == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "coin name cannot be empty")
	}

	if msg.TotalSupply.IsZero() || msg.TotalSupply.IsNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "total supply must be positive")
	}

	return nil
}

// PaymentMethodMsg represents a payment method in messages
type PaymentMethodMsg struct {
	Type        string `json:"type" yaml:"type"`
	Provider    string `json:"provider" yaml:"provider"`
	Identifier  string `json:"identifier" yaml:"identifier"`
	IsPreferred bool   `json:"is_preferred" yaml:"is_preferred"`
}

// MsgRegisterEnhancedMitra defines a message to register an enhanced mitra
type MsgRegisterEnhancedMitra struct {
	Sender           string              `json:"sender" yaml:"sender"`
	MitraId          string              `json:"mitra_id" yaml:"mitra_id"`
	DhanpataName     string              `json:"dhanpata_name" yaml:"dhanpata_name"`
	MitraType        string              `json:"mitra_type" yaml:"mitra_type"`
	TrustScore       int64               `json:"trust_score" yaml:"trust_score"`
	Specializations  []string            `json:"specializations" yaml:"specializations"`
	OperatingRegions []string            `json:"operating_regions" yaml:"operating_regions"`
	PaymentMethods   []PaymentMethodMsg  `json:"payment_methods" yaml:"payment_methods"`
}

// NewMsgRegisterEnhancedMitra creates a new MsgRegisterEnhancedMitra instance
func NewMsgRegisterEnhancedMitra(
	sender sdk.AccAddress,
	mitraId, dhanpataName, mitraType string,
	trustScore int64,
	specializations, operatingRegions []string,
	paymentMethods []PaymentMethodMsg,
) *MsgRegisterEnhancedMitra {
	return &MsgRegisterEnhancedMitra{
		Sender:           sender.String(),
		MitraId:          mitraId,
		DhanpataName:     dhanpataName,
		MitraType:        mitraType,
		TrustScore:       trustScore,
		Specializations:  specializations,
		OperatingRegions: operatingRegions,
		PaymentMethods:   paymentMethods,
	}
}

func (msg *MsgRegisterEnhancedMitra) Route() string { return RouterKey }
func (msg *MsgRegisterEnhancedMitra) Type() string  { return TypeMsgRegisterEnhancedMitra }

func (msg *MsgRegisterEnhancedMitra) GetSigners() []sdk.AccAddress {
	sender, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{sender}
}

func (msg *MsgRegisterEnhancedMitra) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRegisterEnhancedMitra) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address: %s", err)
	}

	if msg.MitraType != MitraTypeIndividual && msg.MitraType != MitraTypeBusiness && msg.MitraType != MitraTypeGlobal {
		return ErrInvalidMitraType
	}

	if msg.TrustScore < 0 || msg.TrustScore > 100 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "trust score must be between 0 and 100")
	}

	return nil
}

// MsgProcessMoneyOrderWithDhanPata defines a message to process money order with DhanPata
type MsgProcessMoneyOrderWithDhanPata struct {
	Sender           string   `json:"sender" yaml:"sender"`
	ReceiverDhanpata string   `json:"receiver_dhanpata" yaml:"receiver_dhanpata"`
	Amount           sdk.Coin `json:"amount" yaml:"amount"`
	Note             string   `json:"note" yaml:"note"`
}

// NewMsgProcessMoneyOrderWithDhanPata creates a new MsgProcessMoneyOrderWithDhanPata instance
func NewMsgProcessMoneyOrderWithDhanPata(
	sender sdk.AccAddress,
	receiverDhanpata string,
	amount sdk.Coin,
	note string,
) *MsgProcessMoneyOrderWithDhanPata {
	return &MsgProcessMoneyOrderWithDhanPata{
		Sender:           sender.String(),
		ReceiverDhanpata: receiverDhanpata,
		Amount:           amount,
		Note:             note,
	}
}

func (msg *MsgProcessMoneyOrderWithDhanPata) Route() string { return RouterKey }
func (msg *MsgProcessMoneyOrderWithDhanPata) Type() string  { return TypeMsgProcessMoneyOrderWithDhanPata }

func (msg *MsgProcessMoneyOrderWithDhanPata) GetSigners() []sdk.AccAddress {
	sender, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{sender}
}

func (msg *MsgProcessMoneyOrderWithDhanPata) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgProcessMoneyOrderWithDhanPata) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address: %s", err)
	}

	if err := ValidateDhanPataName(msg.ReceiverDhanpata); err != nil {
		return err
	}

	if !msg.Amount.IsValid() || msg.Amount.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid amount")
	}

	return nil
}

// MsgUpdateDhanPataMetadata defines a message to update DhanPata metadata
type MsgUpdateDhanPataMetadata struct {
	Owner           string   `json:"owner" yaml:"owner"`
	Name            string   `json:"name" yaml:"name"`
	DisplayName     string   `json:"display_name" yaml:"display_name"`
	Description     string   `json:"description" yaml:"description"`
	ProfileImageUrl string   `json:"profile_image_url" yaml:"profile_image_url"`
	Tags            []string `json:"tags" yaml:"tags"`
}

// NewMsgUpdateDhanPataMetadata creates a new MsgUpdateDhanPataMetadata instance
func NewMsgUpdateDhanPataMetadata(
	owner sdk.AccAddress,
	name, displayName, description, profileImageUrl string,
	tags []string,
) *MsgUpdateDhanPataMetadata {
	return &MsgUpdateDhanPataMetadata{
		Owner:           owner.String(),
		Name:            name,
		DisplayName:     displayName,
		Description:     description,
		ProfileImageUrl: profileImageUrl,
		Tags:            tags,
	}
}

func (msg *MsgUpdateDhanPataMetadata) Route() string { return RouterKey }
func (msg *MsgUpdateDhanPataMetadata) Type() string  { return TypeMsgUpdateDhanPataMetadata }

func (msg *MsgUpdateDhanPataMetadata) GetSigners() []sdk.AccAddress {
	owner, _ := sdk.AccAddressFromBech32(msg.Owner)
	return []sdk.AccAddress{owner}
}

func (msg *MsgUpdateDhanPataMetadata) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateDhanPataMetadata) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address: %s", err)
	}

	if err := ValidateDhanPataName(msg.Name); err != nil {
		return err
	}

	return nil
}

// Response types

// MsgRegisterDhanPataAddressResponse defines the response for registering DhanPata address
type MsgRegisterDhanPataAddressResponse struct {
	Success      bool   `json:"success" yaml:"success"`
	DhanpataName string `json:"dhanpata_name" yaml:"dhanpata_name"`
}

// MsgCreateKshetraCoinResponse defines the response for creating Kshetra coin
type MsgCreateKshetraCoinResponse struct {
	Success    bool   `json:"success" yaml:"success"`
	CoinSymbol string `json:"coin_symbol" yaml:"coin_symbol"`
	Pincode    string `json:"pincode" yaml:"pincode"`
}

// MsgRegisterEnhancedMitraResponse defines the response for registering enhanced mitra
type MsgRegisterEnhancedMitraResponse struct {
	Success bool   `json:"success" yaml:"success"`
	MitraId string `json:"mitra_id" yaml:"mitra_id"`
}

// MsgProcessMoneyOrderWithDhanPataResponse defines the response for processing money order
type MsgProcessMoneyOrderWithDhanPataResponse struct {
	Success          bool     `json:"success" yaml:"success"`
	ReceiverDhanpata string   `json:"receiver_dhanpata" yaml:"receiver_dhanpata"`
	ProcessedAmount  sdk.Coin `json:"processed_amount" yaml:"processed_amount"`
	Fee              sdk.Coin `json:"fee" yaml:"fee"`
}

// MsgUpdateDhanPataMetadataResponse defines the response for updating DhanPata metadata
type MsgUpdateDhanPataMetadataResponse struct {
	Success bool `json:"success" yaml:"success"`
}