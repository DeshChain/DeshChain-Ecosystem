package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgMigrateToIdentity         = "migrate_to_identity"
	TypeMsgCreateParticipantCredential = "create_participant_credential"
)

var (
	_ sdk.Msg = &MsgMigrateToIdentity{}
	_ sdk.Msg = &MsgCreateParticipantCredential{}
)

// MsgMigrateToIdentity migrates participants to identity system
type MsgMigrateToIdentity struct {
	Authority string `json:"authority"`
}

func NewMsgMigrateToIdentity(authority string) *MsgMigrateToIdentity {
	return &MsgMigrateToIdentity{
		Authority: authority,
	}
}

func (msg *MsgMigrateToIdentity) Route() string {
	return RouterKey
}

func (msg *MsgMigrateToIdentity) Type() string {
	return TypeMsgMigrateToIdentity
}

func (msg *MsgMigrateToIdentity) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgMigrateToIdentity) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgMigrateToIdentity) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	return nil
}

// MsgCreateParticipantCredential creates a credential for a participant
type MsgCreateParticipantCredential struct {
	Authority     string `json:"authority"`
	ParticipantId string `json:"participant_id"`
	SchemeId      string `json:"scheme_id"`
}

func NewMsgCreateParticipantCredential(authority, participantId, schemeId string) *MsgCreateParticipantCredential {
	return &MsgCreateParticipantCredential{
		Authority:     authority,
		ParticipantId: participantId,
		SchemeId:      schemeId,
	}
}

func (msg *MsgCreateParticipantCredential) Route() string {
	return RouterKey
}

func (msg *MsgCreateParticipantCredential) Type() string {
	return TypeMsgCreateParticipantCredential
}

func (msg *MsgCreateParticipantCredential) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgCreateParticipantCredential) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateParticipantCredential) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	
	if msg.ParticipantId == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "participant ID cannot be empty")
	}
	
	if msg.SchemeId == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "scheme ID cannot be empty")
	}
	
	return nil
}

// Query types for identity integration

// QueryParticipantIdentityRequest is the request for participant identity
type QueryParticipantIdentityRequest struct {
	ParticipantAddress string `json:"participant_address"`
}

// QueryParticipantIdentityResponse is the response for participant identity
type QueryParticipantIdentityResponse struct {
	Address                 string                      `json:"address"`
	HasTraditionalKYC       bool                        `json:"has_traditional_kyc"`
	HasIdentity             bool                        `json:"has_identity"`
	DID                     string                      `json:"did,omitempty"`
	IsKYCVerified           bool                        `json:"is_kyc_verified"`
	KYCLevel                string                      `json:"kyc_level"`
	EnrolledSchemes         []string                    `json:"enrolled_schemes"`
	GramSurakshaCredentials []GramSurakshaCredentialInfo `json:"gramsuraksha_credentials"`
}

// GramSurakshaCredentialInfo contains credential information
type GramSurakshaCredentialInfo struct {
	CredentialID string `json:"credential_id"`
	SchemeID     string `json:"scheme_id"`
	Status       string `json:"status"`
	IsValid      bool   `json:"is_valid"`
}

// QueryVerifyAgeRequest is the request for age verification
type QueryVerifyAgeRequest struct {
	Address string `json:"address"`
	MinAge  int32  `json:"min_age"`
	MaxAge  int32  `json:"max_age"`
}

// QueryVerifyAgeResponse is the response for age verification
type QueryVerifyAgeResponse struct {
	Verified     bool   `json:"verified"`
	Method       string `json:"method"` // "traditional" or "zkp"
	ErrorMessage string `json:"error_message,omitempty"`
}