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
	TypeMsgVetoProposal                    = "veto_proposal"
	TypeMsgApproveFounderConsentProposal   = "approve_founder_consent_proposal"
	TypeMsgUpdateProtectedParameter        = "update_protected_parameter"
)

var (
	_ sdk.Msg = &MsgVetoProposal{}
	_ sdk.Msg = &MsgApproveFounderConsentProposal{}
	_ sdk.Msg = &MsgUpdateProtectedParameter{}
)

// MsgVetoProposal defines a message to veto a proposal
type MsgVetoProposal struct {
	Authority  string `json:"authority" yaml:"authority"`
	ProposalId uint64 `json:"proposal_id" yaml:"proposal_id"`
	Reason     string `json:"reason" yaml:"reason"`
}

func NewMsgVetoProposal(authority string, proposalId uint64, reason string) *MsgVetoProposal {
	return &MsgVetoProposal{
		Authority:  authority,
		ProposalId: proposalId,
		Reason:     reason,
	}
}

func (msg *MsgVetoProposal) Route() string {
	return RouterKey
}

func (msg *MsgVetoProposal) Type() string {
	return TypeMsgVetoProposal
}

func (msg *MsgVetoProposal) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgVetoProposal) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgVetoProposal) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}

	if msg.ProposalId == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "proposal id cannot be 0")
	}

	if msg.Reason == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "reason cannot be empty")
	}

	return nil
}

// MsgApproveFounderConsentProposal defines a message to approve a proposal requiring founder consent
type MsgApproveFounderConsentProposal struct {
	Authority  string `json:"authority" yaml:"authority"`
	ProposalId uint64 `json:"proposal_id" yaml:"proposal_id"`
}

func NewMsgApproveFounderConsentProposal(authority string, proposalId uint64) *MsgApproveFounderConsentProposal {
	return &MsgApproveFounderConsentProposal{
		Authority:  authority,
		ProposalId: proposalId,
	}
}

func (msg *MsgApproveFounderConsentProposal) Route() string {
	return RouterKey
}

func (msg *MsgApproveFounderConsentProposal) Type() string {
	return TypeMsgApproveFounderConsentProposal
}

func (msg *MsgApproveFounderConsentProposal) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgApproveFounderConsentProposal) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgApproveFounderConsentProposal) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}

	if msg.ProposalId == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "proposal id cannot be 0")
	}

	return nil
}

// MsgUpdateProtectedParameter defines a message to update a protected parameter
type MsgUpdateProtectedParameter struct {
	Authority     string         `json:"authority" yaml:"authority"`
	Name          string         `json:"name" yaml:"name"`
	NewProtection ProtectionType `json:"new_protection" yaml:"new_protection"`
}

func NewMsgUpdateProtectedParameter(authority string, name string, newProtection ProtectionType) *MsgUpdateProtectedParameter {
	return &MsgUpdateProtectedParameter{
		Authority:     authority,
		Name:          name,
		NewProtection: newProtection,
	}
}

func (msg *MsgUpdateProtectedParameter) Route() string {
	return RouterKey
}

func (msg *MsgUpdateProtectedParameter) Type() string {
	return TypeMsgUpdateProtectedParameter
}

func (msg *MsgUpdateProtectedParameter) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgUpdateProtectedParameter) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateProtectedParameter) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}

	if msg.Name == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "parameter name cannot be empty")
	}

	if msg.NewProtection < ProtectionType_NONE || msg.NewProtection > ProtectionType_SUPERMAJORITY {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid protection type")
	}

	return nil
}

// Response types
type MsgVetoProposalResponse struct{}
type MsgApproveFounderConsentProposalResponse struct{}
type MsgUpdateProtectedParameterResponse struct{}