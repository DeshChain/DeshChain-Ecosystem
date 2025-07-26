package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgMigrateRemittanceToIdentity    = "migrate_remittance_to_identity"
	TypeMsgCreateSenderCredential         = "create_sender_credential"
	TypeMsgCreateSewaMitraCredential      = "create_sewa_mitra_credential"
	TypeMsgCreateRecipientCredential      = "create_recipient_credential"
	TypeMsgCreateTransferCredential       = "create_transfer_credential"
)

var (
	_ sdk.Msg = &MsgMigrateRemittanceToIdentity{}
	_ sdk.Msg = &MsgCreateSenderCredential{}
	_ sdk.Msg = &MsgCreateSewaMitraCredential{}
	_ sdk.Msg = &MsgCreateRecipientCredential{}
	_ sdk.Msg = &MsgCreateTransferCredential{}
)

// MsgMigrateRemittanceToIdentity migrates remittance data to identity system
type MsgMigrateRemittanceToIdentity struct {
	Authority string `json:"authority"`
}

func NewMsgMigrateRemittanceToIdentity(authority string) *MsgMigrateRemittanceToIdentity {
	return &MsgMigrateRemittanceToIdentity{
		Authority: authority,
	}
}

func (msg *MsgMigrateRemittanceToIdentity) Route() string {
	return RouterKey
}

func (msg *MsgMigrateRemittanceToIdentity) Type() string {
	return TypeMsgMigrateRemittanceToIdentity
}

func (msg *MsgMigrateRemittanceToIdentity) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgMigrateRemittanceToIdentity) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgMigrateRemittanceToIdentity) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	return nil
}

// MsgCreateSenderCredential creates a remittance sender credential
type MsgCreateSenderCredential struct {
	Authority           string   `json:"authority"`
	SenderAddress       string   `json:"sender_address"`
	KycLevel            string   `json:"kyc_level"`
	SourceOfFunds       string   `json:"source_of_funds"`
	RiskLevel           string   `json:"risk_level"`
	MaxTransferLimit    string   `json:"max_transfer_limit"`
	DailyLimit          string   `json:"daily_limit"`
	MonthlyLimit        string   `json:"monthly_limit"`
	CountryRestrictions []string `json:"country_restrictions"`
}

func NewMsgCreateSenderCredential(authority, senderAddress, kycLevel, sourceOfFunds, riskLevel, maxTransferLimit, dailyLimit, monthlyLimit string, countryRestrictions []string) *MsgCreateSenderCredential {
	return &MsgCreateSenderCredential{
		Authority:           authority,
		SenderAddress:       senderAddress,
		KycLevel:            kycLevel,
		SourceOfFunds:       sourceOfFunds,
		RiskLevel:           riskLevel,
		MaxTransferLimit:    maxTransferLimit,
		DailyLimit:          dailyLimit,
		MonthlyLimit:        monthlyLimit,
		CountryRestrictions: countryRestrictions,
	}
}

func (msg *MsgCreateSenderCredential) Route() string {
	return RouterKey
}

func (msg *MsgCreateSenderCredential) Type() string {
	return TypeMsgCreateSenderCredential
}

func (msg *MsgCreateSenderCredential) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgCreateSenderCredential) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateSenderCredential) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	
	_, err = sdk.AccAddressFromBech32(msg.SenderAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	
	if msg.KycLevel == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "KYC level cannot be empty")
	}
	
	if msg.SourceOfFunds == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "source of funds cannot be empty")
	}
	
	if msg.RiskLevel == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "risk level cannot be empty")
	}
	
	// Validate coin amounts
	if msg.MaxTransferLimit != "" {
		_, err = sdk.ParseCoinNormalized(msg.MaxTransferLimit)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid max transfer limit: %s", err)
		}
	}
	
	if msg.DailyLimit != "" {
		_, err = sdk.ParseCoinNormalized(msg.DailyLimit)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid daily limit: %s", err)
		}
	}
	
	if msg.MonthlyLimit != "" {
		_, err = sdk.ParseCoinNormalized(msg.MonthlyLimit)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid monthly limit: %s", err)
		}
	}
	
	return nil
}

// MsgCreateSewaMitraCredential creates a Sewa Mitra agent credential
type MsgCreateSewaMitraCredential struct {
	Authority           string   `json:"authority"`
	AgentAddress        string   `json:"agent_address"`
	AgentID             string   `json:"agent_id"`
	BusinessName        string   `json:"business_name"`
	ServiceAreas        []string `json:"service_areas"`
	SupportedCurrencies []string `json:"supported_currencies"`
	MaxTransactionLimit string   `json:"max_transaction_limit"`
	Certifications      []string `json:"certifications"`
}

func NewMsgCreateSewaMitraCredential(authority, agentAddress, agentID, businessName string, serviceAreas, supportedCurrencies []string, maxTransactionLimit string, certifications []string) *MsgCreateSewaMitraCredential {
	return &MsgCreateSewaMitraCredential{
		Authority:           authority,
		AgentAddress:        agentAddress,
		AgentID:             agentID,
		BusinessName:        businessName,
		ServiceAreas:        serviceAreas,
		SupportedCurrencies: supportedCurrencies,
		MaxTransactionLimit: maxTransactionLimit,
		Certifications:      certifications,
	}
}

func (msg *MsgCreateSewaMitraCredential) Route() string {
	return RouterKey
}

func (msg *MsgCreateSewaMitraCredential) Type() string {
	return TypeMsgCreateSewaMitraCredential
}

func (msg *MsgCreateSewaMitraCredential) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgCreateSewaMitraCredential) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateSewaMitraCredential) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	
	_, err = sdk.AccAddressFromBech32(msg.AgentAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid agent address (%s)", err)
	}
	
	if msg.AgentID == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "agent ID cannot be empty")
	}
	
	if msg.BusinessName == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "business name cannot be empty")
	}
	
	if len(msg.ServiceAreas) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "at least one service area must be specified")
	}
	
	if len(msg.SupportedCurrencies) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "at least one supported currency must be specified")
	}
	
	if msg.MaxTransactionLimit != "" {
		_, err = sdk.ParseCoinNormalized(msg.MaxTransactionLimit)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid max transaction limit: %s", err)
		}
	}
	
	return nil
}

// MsgCreateRecipientCredential creates a remittance recipient credential
type MsgCreateRecipientCredential struct {
	Authority         string   `json:"authority"`
	RecipientAddress  string   `json:"recipient_address"`
	KycLevel          string   `json:"kyc_level"`
	VerificationDocs  []string `json:"verification_docs"`
	PurposeOfFunds    string   `json:"purpose_of_funds"`
	BeneficiaryType   string   `json:"beneficiary_type"`
	Country           string   `json:"country"`
	CanReceiveFrom    []string `json:"can_receive_from"`
}

func NewMsgCreateRecipientCredential(authority, recipientAddress, kycLevel string, verificationDocs []string, purposeOfFunds, beneficiaryType, country string, canReceiveFrom []string) *MsgCreateRecipientCredential {
	return &MsgCreateRecipientCredential{
		Authority:        authority,
		RecipientAddress: recipientAddress,
		KycLevel:         kycLevel,
		VerificationDocs: verificationDocs,
		PurposeOfFunds:   purposeOfFunds,
		BeneficiaryType:  beneficiaryType,
		Country:          country,
		CanReceiveFrom:   canReceiveFrom,
	}
}

func (msg *MsgCreateRecipientCredential) Route() string {
	return RouterKey
}

func (msg *MsgCreateRecipientCredential) Type() string {
	return TypeMsgCreateRecipientCredential
}

func (msg *MsgCreateRecipientCredential) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgCreateRecipientCredential) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateRecipientCredential) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	
	_, err = sdk.AccAddressFromBech32(msg.RecipientAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid recipient address (%s)", err)
	}
	
	if msg.KycLevel == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "KYC level cannot be empty")
	}
	
	if msg.Country == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "country cannot be empty")
	}
	
	return nil
}

// MsgCreateTransferCredential creates a transfer completion credential
type MsgCreateTransferCredential struct {
	Authority  string `json:"authority"`
	TransferID string `json:"transfer_id"`
}

func NewMsgCreateTransferCredential(authority, transferID string) *MsgCreateTransferCredential {
	return &MsgCreateTransferCredential{
		Authority:  authority,
		TransferID: transferID,
	}
}

func (msg *MsgCreateTransferCredential) Route() string {
	return RouterKey
}

func (msg *MsgCreateTransferCredential) Type() string {
	return TypeMsgCreateTransferCredential
}

func (msg *MsgCreateTransferCredential) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgCreateTransferCredential) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateTransferCredential) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	
	if msg.TransferID == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "transfer ID cannot be empty")
	}
	
	return nil
}

// Query types for identity integration

// QuerySenderIdentityRequest is the request for sender identity verification
type QuerySenderIdentityRequest struct {
	SenderAddress    string `json:"sender_address"`
	TransferAmount   string `json:"transfer_amount"`
	RecipientCountry string `json:"recipient_country"`
}

// QuerySenderIdentityResponse is the response for sender identity verification
type QuerySenderIdentityResponse struct {
	Address             string   `json:"address"`
	HasIdentity         bool     `json:"has_identity"`
	DID                 string   `json:"did,omitempty"`
	IsKYCVerified       bool     `json:"is_kyc_verified"`
	KYCLevel            string   `json:"kyc_level"`
	AMLVerified         bool     `json:"aml_verified"`
	SanctionsChecked    bool     `json:"sanctions_checked"`
	SourceOfFunds       string   `json:"source_of_funds"`
	RiskLevel           string   `json:"risk_level"`
	MaxTransferLimit    string   `json:"max_transfer_limit"`
	DailyLimit          string   `json:"daily_limit"`
	MonthlyLimit        string   `json:"monthly_limit"`
	CountryRestrictions []string `json:"country_restrictions"`
	ComplianceExpiry    string   `json:"compliance_expiry"`
}

// QueryRecipientIdentityRequest is the request for recipient identity verification
type QueryRecipientIdentityRequest struct {
	RecipientAddress string `json:"recipient_address"`
	SenderCountry    string `json:"sender_country"`
}

// QueryRecipientIdentityResponse is the response for recipient identity verification
type QueryRecipientIdentityResponse struct {
	Address           string   `json:"address"`
	HasIdentity       bool     `json:"has_identity"`
	DID               string   `json:"did,omitempty"`
	IsKYCVerified     bool     `json:"is_kyc_verified"`
	KYCLevel          string   `json:"kyc_level"`
	VerificationDocs  []string `json:"verification_docs"`
	PurposeOfFunds    string   `json:"purpose_of_funds"`
	BeneficiaryType   string   `json:"beneficiary_type"`
	Country           string   `json:"country"`
	CanReceiveFrom    []string `json:"can_receive_from"`
}

// QuerySewaMitraIdentityRequest is the request for Sewa Mitra agent identity verification
type QuerySewaMitraIdentityRequest struct {
	AgentAddress string `json:"agent_address"`
	AgentID      string `json:"agent_id"`
}

// QuerySewaMitraIdentityResponse is the response for Sewa Mitra agent identity verification
type QuerySewaMitraIdentityResponse struct {
	Address              string   `json:"address"`
	HasIdentity          bool     `json:"has_identity"`
	DID                  string   `json:"did,omitempty"`
	IsKYCVerified        bool     `json:"is_kyc_verified"`
	KYCLevel             string   `json:"kyc_level"`
	BusinessLicenseValid bool     `json:"business_license_valid"`
	ComplianceVerified   bool     `json:"compliance_verified"`
	AMLCompliant         bool     `json:"aml_compliant"`
	CertificationsValid  []string `json:"certifications_valid"`
	ServiceAreas         []string `json:"service_areas"`
	SupportedCurrencies  []string `json:"supported_currencies"`
	MaxTransactionLimit  string   `json:"max_transaction_limit"`
	BackgroundVerified   bool     `json:"background_verified"`
	InsuranceCovered     bool     `json:"insurance_covered"`
	Rating               string   `json:"rating"`
}

// QueryRemittanceComplianceRequest is the request for comprehensive compliance check
type QueryRemittanceComplianceRequest struct {
	SenderAddress     string `json:"sender_address"`
	RecipientAddress  string `json:"recipient_address,omitempty"`
	AgentAddress      string `json:"agent_address,omitempty"`
	TransferAmount    string `json:"transfer_amount"`
	SourceCurrency    string `json:"source_currency"`
	DestCurrency      string `json:"dest_currency"`
	SenderCountry     string `json:"sender_country"`
	RecipientCountry  string `json:"recipient_country"`
}

// QueryRemittanceComplianceResponse is the response for comprehensive compliance check
type QueryRemittanceComplianceResponse struct {
	SenderCompliant     bool     `json:"sender_compliant"`
	RecipientCompliant  bool     `json:"recipient_compliant"`
	AgentCompliant      bool     `json:"agent_compliant"`
	CorridorAllowed     bool     `json:"corridor_allowed"`
	AmountWithinLimits  bool     `json:"amount_within_limits"`
	SanctionsCleared    bool     `json:"sanctions_cleared"`
	AMLCompliant        bool     `json:"aml_compliant"`
	RegulatoryCompliant bool     `json:"regulatory_compliant"`
	ComplianceScore     string   `json:"compliance_score"`
	RequiredActions     []string `json:"required_actions"`
	Warnings            []string `json:"warnings"`
}