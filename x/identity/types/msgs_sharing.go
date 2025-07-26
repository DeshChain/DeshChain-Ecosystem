package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Cross-Module Identity Sharing Messages

const (
	TypeMsgCreateShareRequest      = "create_share_request"
	TypeMsgApproveShareRequest     = "approve_share_request"
	TypeMsgDenyShareRequest        = "deny_share_request"
	TypeMsgCreateSharingAgreement  = "create_sharing_agreement"
	TypeMsgUpdateSharingAgreement  = "update_sharing_agreement"
	TypeMsgRevokeSharingAgreement  = "revoke_sharing_agreement"
	TypeMsgCreateAccessPolicy      = "create_access_policy"
	TypeMsgUpdateAccessPolicy      = "update_access_policy"
	TypeMsgDeleteAccessPolicy      = "delete_access_policy"
)

// MsgCreateShareRequest creates a new identity sharing request
type MsgCreateShareRequest struct {
	Authority       string        `json:"authority"`
	RequesterModule string        `json:"requester_module"`
	ProviderModule  string        `json:"provider_module"`
	HolderDID       string        `json:"holder_did"`
	RequestedData   []DataRequest `json:"requested_data"`
	Purpose         string        `json:"purpose"`
	Justification   string        `json:"justification"`
	TTLHours        int64         `json:"ttl_hours"` // TTL in hours
}

// NewMsgCreateShareRequest creates a new MsgCreateShareRequest
func NewMsgCreateShareRequest(
	authority string,
	requesterModule string,
	providerModule string,
	holderDID string,
	requestedData []DataRequest,
	purpose string,
	justification string,
	ttlHours int64,
) *MsgCreateShareRequest {
	return &MsgCreateShareRequest{
		Authority:       authority,
		RequesterModule: requesterModule,
		ProviderModule:  providerModule,
		HolderDID:       holderDID,
		RequestedData:   requestedData,
		Purpose:         purpose,
		Justification:   justification,
		TTLHours:        ttlHours,
	}
}

// Route returns the message route
func (msg *MsgCreateShareRequest) Route() string { return RouterKey }

// Type returns the message type
func (msg *MsgCreateShareRequest) Type() string { return TypeMsgCreateShareRequest }

// GetSigners returns the signers
func (msg *MsgCreateShareRequest) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

// GetSignBytes returns the bytes to sign
func (msg *MsgCreateShareRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// ValidateBasic performs basic validation
func (msg *MsgCreateShareRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address: %v", err)
	}

	if msg.RequesterModule == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "requester module cannot be empty")
	}
	if msg.ProviderModule == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "provider module cannot be empty")
	}
	if msg.HolderDID == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "holder DID cannot be empty")
	}
	if len(msg.RequestedData) == 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "at least one data request is required")
	}
	if msg.Purpose == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "purpose cannot be empty")
	}
	if msg.TTLHours <= 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "TTL must be positive")
	}

	// Validate individual data requests
	for i, dataReq := range msg.RequestedData {
		if err := dataReq.ValidateBasic(); err != nil {
			return sdkerrors.Wrapf(err, "invalid data request at index %d", i)
		}
	}

	return nil
}

// MsgApproveShareRequest approves a pending share request
type MsgApproveShareRequest struct {
	Authority   string `json:"authority"`
	RequestID   string `json:"request_id"`
	AccessToken string `json:"access_token,omitempty"`
}

// NewMsgApproveShareRequest creates a new MsgApproveShareRequest
func NewMsgApproveShareRequest(authority, requestID, accessToken string) *MsgApproveShareRequest {
	return &MsgApproveShareRequest{
		Authority:   authority,
		RequestID:   requestID,
		AccessToken: accessToken,
	}
}

// Route returns the message route
func (msg *MsgApproveShareRequest) Route() string { return RouterKey }

// Type returns the message type
func (msg *MsgApproveShareRequest) Type() string { return TypeMsgApproveShareRequest }

// GetSigners returns the signers
func (msg *MsgApproveShareRequest) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

// GetSignBytes returns the bytes to sign
func (msg *MsgApproveShareRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// ValidateBasic performs basic validation
func (msg *MsgApproveShareRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address: %v", err)
	}

	if msg.RequestID == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "request ID cannot be empty")
	}

	return nil
}

// MsgDenyShareRequest denies a pending share request
type MsgDenyShareRequest struct {
	Authority     string `json:"authority"`
	RequestID     string `json:"request_id"`
	DenialReason  string `json:"denial_reason"`
}

// NewMsgDenyShareRequest creates a new MsgDenyShareRequest
func NewMsgDenyShareRequest(authority, requestID, denialReason string) *MsgDenyShareRequest {
	return &MsgDenyShareRequest{
		Authority:    authority,
		RequestID:    requestID,
		DenialReason: denialReason,
	}
}

// Route returns the message route
func (msg *MsgDenyShareRequest) Route() string { return RouterKey }

// Type returns the message type
func (msg *MsgDenyShareRequest) Type() string { return TypeMsgDenyShareRequest }

// GetSigners returns the signers
func (msg *MsgDenyShareRequest) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

// GetSignBytes returns the bytes to sign
func (msg *MsgDenyShareRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// ValidateBasic performs basic validation
func (msg *MsgDenyShareRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address: %v", err)
	}

	if msg.RequestID == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "request ID cannot be empty")
	}
	if msg.DenialReason == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "denial reason cannot be empty")
	}

	return nil
}

// MsgCreateSharingAgreement creates a standing agreement between modules
type MsgCreateSharingAgreement struct {
	Authority        string   `json:"authority"`
	RequesterModule  string   `json:"requester_module"`
	ProviderModule   string   `json:"provider_module"`
	AllowedDataTypes []string `json:"allowed_data_types"`
	Purposes         []string `json:"purposes"`
	TrustLevel       string   `json:"trust_level"`
	AutoApprove      bool     `json:"auto_approve"`
	MaxTTLHours      int64    `json:"max_ttl_hours"`
	ValidityDays     int64    `json:"validity_days"`
}

// NewMsgCreateSharingAgreement creates a new MsgCreateSharingAgreement
func NewMsgCreateSharingAgreement(
	authority string,
	requesterModule string,
	providerModule string,
	allowedDataTypes []string,
	purposes []string,
	trustLevel string,
	autoApprove bool,
	maxTTLHours int64,
	validityDays int64,
) *MsgCreateSharingAgreement {
	return &MsgCreateSharingAgreement{
		Authority:        authority,
		RequesterModule:  requesterModule,
		ProviderModule:   providerModule,
		AllowedDataTypes: allowedDataTypes,
		Purposes:         purposes,
		TrustLevel:       trustLevel,
		AutoApprove:      autoApprove,
		MaxTTLHours:      maxTTLHours,
		ValidityDays:     validityDays,
	}
}

// Route returns the message route
func (msg *MsgCreateSharingAgreement) Route() string { return RouterKey }

// Type returns the message type
func (msg *MsgCreateSharingAgreement) Type() string { return TypeMsgCreateSharingAgreement }

// GetSigners returns the signers
func (msg *MsgCreateSharingAgreement) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

// GetSignBytes returns the bytes to sign
func (msg *MsgCreateSharingAgreement) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// ValidateBasic performs basic validation
func (msg *MsgCreateSharingAgreement) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address: %v", err)
	}

	if msg.RequesterModule == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "requester module cannot be empty")
	}
	if msg.ProviderModule == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "provider module cannot be empty")
	}
	if len(msg.AllowedDataTypes) == 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "at least one allowed data type is required")
	}
	if len(msg.Purposes) == 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "at least one purpose is required")
	}
	if msg.MaxTTLHours <= 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "max TTL must be positive")
	}
	if msg.ValidityDays <= 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "validity period must be positive")
	}

	return nil
}

// MsgCreateAccessPolicy creates an access policy for identity sharing
type MsgCreateAccessPolicy struct {
	Authority              string            `json:"authority"`
	HolderDID              string            `json:"holder_did"`
	AllowedModules         []string          `json:"allowed_modules,omitempty"`
	DeniedModules          []string          `json:"denied_modules,omitempty"`
	DataRestrictions       map[string][]string `json:"data_restrictions,omitempty"`
	PurposeRestrictions    []string          `json:"purpose_restrictions,omitempty"`
	TimeRestrictions       TimeRestriction   `json:"time_restrictions,omitempty"`
	GeographicRestrictions []string          `json:"geographic_restrictions,omitempty"`
	MaxSharesPerDay        int               `json:"max_shares_per_day"`
	RequireExplicitConsent bool              `json:"require_explicit_consent"`
}

// NewMsgCreateAccessPolicy creates a new MsgCreateAccessPolicy
func NewMsgCreateAccessPolicy(
	authority string,
	holderDID string,
	allowedModules []string,
	deniedModules []string,
	dataRestrictions map[string][]string,
	purposeRestrictions []string,
	timeRestrictions TimeRestriction,
	geographicRestrictions []string,
	maxSharesPerDay int,
	requireExplicitConsent bool,
) *MsgCreateAccessPolicy {
	return &MsgCreateAccessPolicy{
		Authority:              authority,
		HolderDID:              holderDID,
		AllowedModules:         allowedModules,
		DeniedModules:          deniedModules,
		DataRestrictions:       dataRestrictions,
		PurposeRestrictions:    purposeRestrictions,
		TimeRestrictions:       timeRestrictions,
		GeographicRestrictions: geographicRestrictions,
		MaxSharesPerDay:        maxSharesPerDay,
		RequireExplicitConsent: requireExplicitConsent,
	}
}

// Route returns the message route
func (msg *MsgCreateAccessPolicy) Route() string { return RouterKey }

// Type returns the message type
func (msg *MsgCreateAccessPolicy) Type() string { return TypeMsgCreateAccessPolicy }

// GetSigners returns the signers
func (msg *MsgCreateAccessPolicy) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

// GetSignBytes returns the bytes to sign
func (msg *MsgCreateAccessPolicy) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// ValidateBasic performs basic validation
func (msg *MsgCreateAccessPolicy) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address: %v", err)
	}

	if msg.HolderDID == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "holder DID cannot be empty")
	}
	if msg.MaxSharesPerDay < 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "max shares per day cannot be negative")
	}

	return nil
}

// Response types

// MsgCreateShareRequestResponse is the response for MsgCreateShareRequest
type MsgCreateShareRequestResponse struct {
	RequestID string `json:"request_id"`
}

// MsgApproveShareRequestResponse is the response for MsgApproveShareRequest
type MsgApproveShareRequestResponse struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
}

// MsgDenyShareRequestResponse is the response for MsgDenyShareRequest
type MsgDenyShareRequestResponse struct {
	DeniedAt time.Time `json:"denied_at"`
}

// MsgCreateSharingAgreementResponse is the response for MsgCreateSharingAgreement
type MsgCreateSharingAgreementResponse struct {
	AgreementID string `json:"agreement_id"`
}

// MsgCreateAccessPolicyResponse is the response for MsgCreateAccessPolicy
type MsgCreateAccessPolicyResponse struct {
	PolicyID string `json:"policy_id"`
}
