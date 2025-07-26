package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgMigrateBusinessesToIdentity = "migrate_businesses_to_identity"
	TypeMsgCreateBusinessCredential    = "create_business_credential"
)

var (
	_ sdk.Msg = &MsgMigrateBusinessesToIdentity{}
	_ sdk.Msg = &MsgCreateBusinessCredential{}
)

// MsgMigrateBusinessesToIdentity migrates businesses to identity system
type MsgMigrateBusinessesToIdentity struct {
	Authority string `json:"authority"`
}

func NewMsgMigrateBusinessesToIdentity(authority string) *MsgMigrateBusinessesToIdentity {
	return &MsgMigrateBusinessesToIdentity{
		Authority: authority,
	}
}

func (msg *MsgMigrateBusinessesToIdentity) Route() string {
	return RouterKey
}

func (msg *MsgMigrateBusinessesToIdentity) Type() string {
	return TypeMsgMigrateBusinessesToIdentity
}

func (msg *MsgMigrateBusinessesToIdentity) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgMigrateBusinessesToIdentity) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgMigrateBusinessesToIdentity) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	return nil
}

// MsgCreateBusinessCredential creates a credential for a business
type MsgCreateBusinessCredential struct {
	Authority  string `json:"authority"`
	BusinessId string `json:"business_id"`
}

func NewMsgCreateBusinessCredential(authority, businessId string) *MsgCreateBusinessCredential {
	return &MsgCreateBusinessCredential{
		Authority:  authority,
		BusinessId: businessId,
	}
}

func (msg *MsgCreateBusinessCredential) Route() string {
	return RouterKey
}

func (msg *MsgCreateBusinessCredential) Type() string {
	return TypeMsgCreateBusinessCredential
}

func (msg *MsgCreateBusinessCredential) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgCreateBusinessCredential) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateBusinessCredential) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	
	if msg.BusinessId == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "business ID cannot be empty")
	}
	
	return nil
}

// Query types for identity integration

// QueryBusinessIdentityRequest is the request for business identity
type QueryBusinessIdentityRequest struct {
	BusinessAddress string `json:"business_address"`
}

// QueryBusinessIdentityResponse is the response for business identity
type QueryBusinessIdentityResponse struct {
	Address            string   `json:"address"`
	HasIdentity        bool     `json:"has_identity"`
	DID                string   `json:"did,omitempty"`
	IsKYCVerified      bool     `json:"is_kyc_verified"`
	KYCLevel           string   `json:"kyc_level"`
	HasComplianceDocs  bool     `json:"has_compliance_docs"`
	HasFinancialDocs   bool     `json:"has_financial_docs"`
	BusinessType       string   `json:"business_type"`
	GSTNumber          string   `json:"gst_number,omitempty"`
	PANNumber          string   `json:"pan_number,omitempty"`
	VerifiedRevenue    sdk.Dec  `json:"verified_revenue,omitempty"`
	CreditScore        int32    `json:"credit_score"`
	ActiveLoans        []string `json:"active_loans"`
}

// QueryBusinessCredentialsRequest is the request for business credentials
type QueryBusinessCredentialsRequest struct {
	BusinessAddress string `json:"business_address"`
}

// QueryBusinessCredentialsResponse is the response for business credentials
type QueryBusinessCredentialsResponse struct {
	Address         string                    `json:"address"`
	DID             string                    `json:"did"`
	Credentials     []BusinessCredentialInfo  `json:"credentials"`
}

// BusinessCredentialInfo contains business credential information
type BusinessCredentialInfo struct {
	CredentialID string   `json:"credential_id"`
	Type         []string `json:"type"`
	LoanID       string   `json:"loan_id,omitempty"`
	Status       string   `json:"status"`
	BusinessName string   `json:"business_name,omitempty"`
	IsValid      bool     `json:"is_valid"`
}

// QueryBusinessComplianceRequest is the request for business compliance
type QueryBusinessComplianceRequest struct {
	BusinessAddress string `json:"business_address"`
	BusinessType    string `json:"business_type"`
}

// QueryBusinessComplianceResponse is the response for business compliance
type QueryBusinessComplianceResponse struct {
	Address              string   `json:"address"`
	RequiredLicenses     []string `json:"required_licenses"`
	AvailableLicenses    []string `json:"available_licenses"`
	IsCompliant          bool     `json:"is_compliant"`
	MissingLicenses      []string `json:"missing_licenses"`
	ComplianceCredentialID string `json:"compliance_credential_id,omitempty"`
}