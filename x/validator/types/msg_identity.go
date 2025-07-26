package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgMigrateValidatorsToIdentity    = "migrate_validators_to_identity"
	TypeMsgCreateValidatorCredential      = "create_validator_credential"
	TypeMsgCreateNFTBindingCredential     = "create_nft_binding_credential"
	TypeMsgCreateReferralCredential       = "create_referral_credential"
	TypeMsgCreateTokenLaunchCredential    = "create_token_launch_credential"
	TypeMsgCreateComplianceCredential     = "create_compliance_credential"
)

var (
	_ sdk.Msg = &MsgMigrateValidatorsToIdentity{}
	_ sdk.Msg = &MsgCreateValidatorCredential{}
	_ sdk.Msg = &MsgCreateNFTBindingCredential{}
	_ sdk.Msg = &MsgCreateReferralCredential{}
	_ sdk.Msg = &MsgCreateTokenLaunchCredential{}
	_ sdk.Msg = &MsgCreateComplianceCredential{}
)

// MsgMigrateValidatorsToIdentity migrates validators to identity system
type MsgMigrateValidatorsToIdentity struct {
	Authority string `json:"authority"`
}

func NewMsgMigrateValidatorsToIdentity(authority string) *MsgMigrateValidatorsToIdentity {
	return &MsgMigrateValidatorsToIdentity{
		Authority: authority,
	}
}

func (msg *MsgMigrateValidatorsToIdentity) Route() string {
	return RouterKey
}

func (msg *MsgMigrateValidatorsToIdentity) Type() string {
	return TypeMsgMigrateValidatorsToIdentity
}

func (msg *MsgMigrateValidatorsToIdentity) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgMigrateValidatorsToIdentity) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgMigrateValidatorsToIdentity) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	return nil
}

// MsgCreateValidatorCredential creates a validator credential
type MsgCreateValidatorCredential struct {
	Authority         string `json:"authority"`
	ValidatorAddress  string `json:"validator_address"`
	OperatorAddress   string `json:"operator_address"`
	ValidatorRank     uint32 `json:"validator_rank"`
	StakeAmount       string `json:"stake_amount"`
}

func NewMsgCreateValidatorCredential(authority, validatorAddress, operatorAddress string, validatorRank uint32, stakeAmount string) *MsgCreateValidatorCredential {
	return &MsgCreateValidatorCredential{
		Authority:        authority,
		ValidatorAddress: validatorAddress,
		OperatorAddress:  operatorAddress,
		ValidatorRank:    validatorRank,
		StakeAmount:      stakeAmount,
	}
}

func (msg *MsgCreateValidatorCredential) Route() string {
	return RouterKey
}

func (msg *MsgCreateValidatorCredential) Type() string {
	return TypeMsgCreateValidatorCredential
}

func (msg *MsgCreateValidatorCredential) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgCreateValidatorCredential) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateValidatorCredential) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	
	_, err = sdk.AccAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid validator address (%s)", err)
	}
	
	if msg.OperatorAddress == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "operator address cannot be empty")
	}
	
	if msg.ValidatorRank == 0 || msg.ValidatorRank > 21 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "validator rank must be between 1 and 21")
	}
	
	if msg.StakeAmount == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "stake amount cannot be empty")
	}
	
	// Validate stake amount is a valid integer
	_, ok := sdk.NewIntFromString(msg.StakeAmount)
	if !ok {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid stake amount: %s", msg.StakeAmount)
	}
	
	return nil
}

// MsgCreateNFTBindingCredential creates an NFT binding credential
type MsgCreateNFTBindingCredential struct {
	Authority         string `json:"authority"`
	ValidatorAddress  string `json:"validator_address"`
	NFTTokenID        uint64 `json:"nft_token_id"`
}

func NewMsgCreateNFTBindingCredential(authority, validatorAddress string, nftTokenID uint64) *MsgCreateNFTBindingCredential {
	return &MsgCreateNFTBindingCredential{
		Authority:        authority,
		ValidatorAddress: validatorAddress,
		NFTTokenID:       nftTokenID,
	}
}

func (msg *MsgCreateNFTBindingCredential) Route() string {
	return RouterKey
}

func (msg *MsgCreateNFTBindingCredential) Type() string {
	return TypeMsgCreateNFTBindingCredential
}

func (msg *MsgCreateNFTBindingCredential) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgCreateNFTBindingCredential) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateNFTBindingCredential) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	
	_, err = sdk.AccAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid validator address (%s)", err)
	}
	
	if msg.NFTTokenID == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "NFT token ID cannot be zero")
	}
	
	return nil
}

// MsgCreateReferralCredential creates a referral credential
type MsgCreateReferralCredential struct {
	Authority        string `json:"authority"`
	ReferrerAddress  string `json:"referrer_address"`
	ReferredAddress  string `json:"referred_address"`
	ReferralID       uint64 `json:"referral_id"`
}

func NewMsgCreateReferralCredential(authority, referrerAddress, referredAddress string, referralID uint64) *MsgCreateReferralCredential {
	return &MsgCreateReferralCredential{
		Authority:       authority,
		ReferrerAddress: referrerAddress,
		ReferredAddress: referredAddress,
		ReferralID:      referralID,
	}
}

func (msg *MsgCreateReferralCredential) Route() string {
	return RouterKey
}

func (msg *MsgCreateReferralCredential) Type() string {
	return TypeMsgCreateReferralCredential
}

func (msg *MsgCreateReferralCredential) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgCreateReferralCredential) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateReferralCredential) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	
	_, err = sdk.AccAddressFromBech32(msg.ReferrerAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid referrer address (%s)", err)
	}
	
	_, err = sdk.AccAddressFromBech32(msg.ReferredAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid referred address (%s)", err)
	}
	
	if msg.ReferralID == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "referral ID cannot be zero")
	}
	
	return nil
}

// MsgCreateTokenLaunchCredential creates a token launch credential
type MsgCreateTokenLaunchCredential struct {
	Authority         string `json:"authority"`
	ValidatorAddress  string `json:"validator_address"`
	TokenID           uint64 `json:"token_id"`
}

func NewMsgCreateTokenLaunchCredential(authority, validatorAddress string, tokenID uint64) *MsgCreateTokenLaunchCredential {
	return &MsgCreateTokenLaunchCredential{
		Authority:        authority,
		ValidatorAddress: validatorAddress,
		TokenID:          tokenID,
	}
}

func (msg *MsgCreateTokenLaunchCredential) Route() string {
	return RouterKey
}

func (msg *MsgCreateTokenLaunchCredential) Type() string {
	return TypeMsgCreateTokenLaunchCredential
}

func (msg *MsgCreateTokenLaunchCredential) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgCreateTokenLaunchCredential) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateTokenLaunchCredential) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	
	_, err = sdk.AccAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid validator address (%s)", err)
	}
	
	if msg.TokenID == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "token ID cannot be zero")
	}
	
	return nil
}

// MsgCreateComplianceCredential creates a compliance credential
type MsgCreateComplianceCredential struct {
	Authority           string   `json:"authority"`
	ValidatorAddress    string   `json:"validator_address"`
	ComplianceLevel     string   `json:"compliance_level"`
	Jurisdictions       []string `json:"jurisdictions"`
	RequiredDocuments   []string `json:"required_documents"`
	VerifiedDocuments   []string `json:"verified_documents"`
}

func NewMsgCreateComplianceCredential(authority, validatorAddress, complianceLevel string, jurisdictions, requiredDocuments, verifiedDocuments []string) *MsgCreateComplianceCredential {
	return &MsgCreateComplianceCredential{
		Authority:         authority,
		ValidatorAddress:  validatorAddress,
		ComplianceLevel:   complianceLevel,
		Jurisdictions:     jurisdictions,
		RequiredDocuments: requiredDocuments,
		VerifiedDocuments: verifiedDocuments,
	}
}

func (msg *MsgCreateComplianceCredential) Route() string {
	return RouterKey
}

func (msg *MsgCreateComplianceCredential) Type() string {
	return TypeMsgCreateComplianceCredential
}

func (msg *MsgCreateComplianceCredential) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgCreateComplianceCredential) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateComplianceCredential) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	
	_, err = sdk.AccAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid validator address (%s)", err)
	}
	
	if msg.ComplianceLevel == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "compliance level cannot be empty")
	}
	
	if len(msg.Jurisdictions) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "at least one jurisdiction must be specified")
	}
	
	return nil
}

// Query types for identity integration

// QueryValidatorIdentityRequest is the request for validator identity
type QueryValidatorIdentityRequest struct {
	ValidatorAddress string `json:"validator_address"`
}

// QueryValidatorIdentityResponse is the response for validator identity
type QueryValidatorIdentityResponse struct {
	Address            string `json:"address"`
	HasIdentity        bool   `json:"has_identity"`
	DID                string `json:"did,omitempty"`
	IsKYCVerified      bool   `json:"is_kyc_verified"`
	KYCLevel           string `json:"kyc_level"`
	ValidatorRank      uint32 `json:"validator_rank"`
	StakeVerified      bool   `json:"stake_verified"`
	NFTBound           bool   `json:"nft_bound"`
	ReferralCredential bool   `json:"referral_credential"`
	TokenLaunched      bool   `json:"token_launched"`
	ComplianceStatus   string `json:"compliance_status"`
	JurisdictionCode   string `json:"jurisdiction_code"`
	GeographicRegion   string `json:"geographic_region"`
}

// QueryValidatorComplianceRequest is the request for validator compliance
type QueryValidatorComplianceRequest struct {
	ValidatorAddress      string   `json:"validator_address"`
	RequiredJurisdictions []string `json:"required_jurisdictions"`
}

// QueryValidatorComplianceResponse is the response for validator compliance
type QueryValidatorComplianceResponse struct {
	Address            string    `json:"address"`
	IsCompliant        bool      `json:"is_compliant"`
	ComplianceLevel    string    `json:"compliance_level"`
	JurisdictionCodes  []string  `json:"jurisdiction_codes"`
	RequiredDocuments  []string  `json:"required_documents"`
	VerifiedDocuments  []string  `json:"verified_documents"`
	AMLVerified        bool      `json:"aml_verified"`
	SanctionsChecked   bool      `json:"sanctions_checked"`
	KYBCompleted       bool      `json:"kyb_completed"`
	ComplianceExpiry   string    `json:"compliance_expiry"`
}

// QueryValidatorCredentialsRequest is the request for validator credentials
type QueryValidatorCredentialsRequest struct {
	ValidatorAddress string `json:"validator_address"`
}

// QueryValidatorCredentialsResponse is the response for validator credentials
type QueryValidatorCredentialsResponse struct {
	Address     string                     `json:"address"`
	DID         string                     `json:"did"`
	Credentials []ValidatorCredentialInfo  `json:"credentials"`
}

// ValidatorCredentialInfo contains validator credential information
type ValidatorCredentialInfo struct {
	CredentialID     string   `json:"credential_id"`
	Type             []string `json:"type"`
	ValidatorRank    uint32   `json:"validator_rank,omitempty"`
	NFTTokenID       uint64   `json:"nft_token_id,omitempty"`
	ReferralID       uint64   `json:"referral_id,omitempty"`
	TokenID          uint64   `json:"token_id,omitempty"`
	ComplianceLevel  string   `json:"compliance_level,omitempty"`
	Status           string   `json:"status"`
	IsValid          bool     `json:"is_valid"`
}