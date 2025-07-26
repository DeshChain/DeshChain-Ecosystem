package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateIdentity         = "create_identity"
	TypeMsgUpdateIdentity         = "update_identity"
	TypeMsgRevokeIdentity         = "revoke_identity"
	TypeMsgRegisterDID            = "register_did"
	TypeMsgUpdateDID              = "update_did"
	TypeMsgDeactivateDID          = "deactivate_did"
	TypeMsgIssueCredential        = "issue_credential"
	TypeMsgRevokeCredential       = "revoke_credential"
	TypeMsgPresentCredential      = "present_credential"
	TypeMsgCreateZKProof          = "create_zk_proof"
	TypeMsgVerifyZKProof          = "verify_zk_proof"
	TypeMsgLinkAadhaar            = "link_aadhaar"
	TypeMsgConnectDigiLocker      = "connect_digilocker"
	TypeMsgLinkUPI                = "link_upi"
	TypeMsgGiveConsent            = "give_consent"
	TypeMsgWithdrawConsent        = "withdraw_consent"
	TypeMsgAddRecoveryMethod      = "add_recovery_method"
	TypeMsgInitiateRecovery       = "initiate_recovery"
	TypeMsgCompleteRecovery       = "complete_recovery"
	TypeMsgUpdatePrivacySettings  = "update_privacy_settings"
)

// Ensure messages implement sdk.Msg interface
var (
	_ sdk.Msg = &MsgCreateIdentity{}
	_ sdk.Msg = &MsgUpdateIdentity{}
	_ sdk.Msg = &MsgRevokeIdentity{}
	_ sdk.Msg = &MsgRegisterDID{}
	_ sdk.Msg = &MsgUpdateDID{}
	_ sdk.Msg = &MsgDeactivateDID{}
	_ sdk.Msg = &MsgIssueCredential{}
	_ sdk.Msg = &MsgRevokeCredential{}
	_ sdk.Msg = &MsgPresentCredential{}
	_ sdk.Msg = &MsgCreateZKProof{}
	_ sdk.Msg = &MsgVerifyZKProof{}
	_ sdk.Msg = &MsgLinkAadhaar{}
	_ sdk.Msg = &MsgConnectDigiLocker{}
	_ sdk.Msg = &MsgLinkUPI{}
	_ sdk.Msg = &MsgGiveConsent{}
	_ sdk.Msg = &MsgWithdrawConsent{}
	_ sdk.Msg = &MsgAddRecoveryMethod{}
	_ sdk.Msg = &MsgInitiateRecovery{}
	_ sdk.Msg = &MsgCompleteRecovery{}
	_ sdk.Msg = &MsgUpdatePrivacySettings{}
)

// MsgCreateIdentity creates a new identity
type MsgCreateIdentity struct {
	Creator          string            `json:"creator"`
	PublicKey        string            `json:"public_key"`
	ServiceEndpoints []Service         `json:"service_endpoints,omitempty"`
	RecoveryMethods  []RecoveryMethod  `json:"recovery_methods"`
	InitialConsents  []ConsentRecord   `json:"initial_consents"`
	Metadata         map[string]string `json:"metadata,omitempty"`
}

func NewMsgCreateIdentity(creator, publicKey string) *MsgCreateIdentity {
	return &MsgCreateIdentity{
		Creator:   creator,
		PublicKey: publicKey,
	}
}

func (msg *MsgCreateIdentity) Route() string { return RouterKey }
func (msg *MsgCreateIdentity) Type() string  { return TypeMsgCreateIdentity }
func (msg *MsgCreateIdentity) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateIdentity) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateIdentity) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address: %s", msg.Creator)
	}
	if msg.PublicKey == "" {
		return sdkerrors.Wrap(ErrInvalidPublicKey, "public key cannot be empty")
	}
	return nil
}

// MsgUpdateIdentity updates an existing identity
type MsgUpdateIdentity struct {
	Creator          string                 `json:"creator"`
	ServiceEndpoints []Service              `json:"service_endpoints,omitempty"`
	RecoveryMethods  []RecoveryMethod       `json:"recovery_methods,omitempty"`
	Metadata         map[string]string      `json:"metadata,omitempty"`
}

func NewMsgUpdateIdentity(creator string) *MsgUpdateIdentity {
	return &MsgUpdateIdentity{
		Creator: creator,
	}
}

func (msg *MsgUpdateIdentity) Route() string { return RouterKey }
func (msg *MsgUpdateIdentity) Type() string  { return TypeMsgUpdateIdentity }
func (msg *MsgUpdateIdentity) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateIdentity) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateIdentity) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address: %s", msg.Creator)
	}
	return nil
}

// MsgRevokeIdentity revokes an identity
type MsgRevokeIdentity struct {
	Creator string `json:"creator"`
	Reason  string `json:"reason"`
}

func NewMsgRevokeIdentity(creator, reason string) *MsgRevokeIdentity {
	return &MsgRevokeIdentity{
		Creator: creator,
		Reason:  reason,
	}
}

func (msg *MsgRevokeIdentity) Route() string { return RouterKey }
func (msg *MsgRevokeIdentity) Type() string  { return TypeMsgRevokeIdentity }
func (msg *MsgRevokeIdentity) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRevokeIdentity) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRevokeIdentity) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address: %s", msg.Creator)
	}
	if msg.Reason == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "reason cannot be empty")
	}
	return nil
}

// MsgRegisterDID registers a new DID
type MsgRegisterDID struct {
	Creator            string               `json:"creator"`
	DIDDocument        *DIDDocument         `json:"did_document"`
}

func NewMsgRegisterDID(creator string, didDoc *DIDDocument) *MsgRegisterDID {
	return &MsgRegisterDID{
		Creator:     creator,
		DIDDocument: didDoc,
	}
}

func (msg *MsgRegisterDID) Route() string { return RouterKey }
func (msg *MsgRegisterDID) Type() string  { return TypeMsgRegisterDID }
func (msg *MsgRegisterDID) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRegisterDID) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRegisterDID) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address: %s", msg.Creator)
	}
	if msg.DIDDocument == nil {
		return sdkerrors.Wrap(ErrInvalidDIDDocument, "DID document cannot be nil")
	}
	if err := ValidateDIDDocument(msg.DIDDocument); err != nil {
		return sdkerrors.Wrap(ErrInvalidDIDDocument, err.Error())
	}
	return nil
}

// MsgUpdateDID updates a DID document
type MsgUpdateDID struct {
	Creator               string               `json:"creator"`
	DID                   string               `json:"did"`
	VerificationMethods   []VerificationMethod `json:"verification_methods,omitempty"`
	Services              []Service            `json:"services,omitempty"`
	RemoveVerificationMethods []string         `json:"remove_verification_methods,omitempty"`
	RemoveServices        []string             `json:"remove_services,omitempty"`
}

func NewMsgUpdateDID(creator, did string) *MsgUpdateDID {
	return &MsgUpdateDID{
		Creator: creator,
		DID:     did,
	}
}

func (msg *MsgUpdateDID) Route() string { return RouterKey }
func (msg *MsgUpdateDID) Type() string  { return TypeMsgUpdateDID }
func (msg *MsgUpdateDID) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateDID) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateDID) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address: %s", msg.Creator)
	}
	if msg.DID == "" {
		return sdkerrors.Wrap(ErrInvalidDID, "DID cannot be empty")
	}
	return nil
}

// MsgDeactivateDID deactivates a DID
type MsgDeactivateDID struct {
	Creator string `json:"creator"`
	DID     string `json:"did"`
	Reason  string `json:"reason"`
}

func NewMsgDeactivateDID(creator, did, reason string) *MsgDeactivateDID {
	return &MsgDeactivateDID{
		Creator: creator,
		DID:     did,
		Reason:  reason,
	}
}

func (msg *MsgDeactivateDID) Route() string { return RouterKey }
func (msg *MsgDeactivateDID) Type() string  { return TypeMsgDeactivateDID }
func (msg *MsgDeactivateDID) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeactivateDID) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeactivateDID) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address: %s", msg.Creator)
	}
	if msg.DID == "" {
		return sdkerrors.Wrap(ErrInvalidDID, "DID cannot be empty")
	}
	return nil
}

// MsgIssueCredential issues a new verifiable credential
type MsgIssueCredential struct {
	Issuer           string                 `json:"issuer"`
	Holder           string                 `json:"holder"`
	CredentialType   string                 `json:"credential_type"`
	Claims           map[string]interface{} `json:"claims"`
	ExpirationDays   int32                  `json:"expiration_days"`
	Evidence         []Evidence             `json:"evidence,omitempty"`
	RequireConsent   bool                   `json:"require_consent"`
	Metadata         map[string]string      `json:"metadata,omitempty"`
}

func NewMsgIssueCredential(issuer, holder, credType string) *MsgIssueCredential {
	return &MsgIssueCredential{
		Issuer:         issuer,
		Holder:         holder,
		CredentialType: credType,
	}
}

func (msg *MsgIssueCredential) Route() string { return RouterKey }
func (msg *MsgIssueCredential) Type() string  { return TypeMsgIssueCredential }
func (msg *MsgIssueCredential) GetSigners() []sdk.AccAddress {
	issuer, err := sdk.AccAddressFromBech32(msg.Issuer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{issuer}
}

func (msg *MsgIssueCredential) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgIssueCredential) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Issuer); err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid issuer address: %s", msg.Issuer)
	}
	if _, err := sdk.AccAddressFromBech32(msg.Holder); err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid holder address: %s", msg.Holder)
	}
	if msg.CredentialType == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "credential type cannot be empty")
	}
	if len(msg.Claims) == 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "claims cannot be empty")
	}
	if msg.ExpirationDays < 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "expiration days cannot be negative")
	}
	return nil
}

// MsgRevokeCredential revokes a credential
type MsgRevokeCredential struct {
	Issuer       string `json:"issuer"`
	CredentialID string `json:"credential_id"`
	Reason       string `json:"reason"`
}

func NewMsgRevokeCredential(issuer, credentialID, reason string) *MsgRevokeCredential {
	return &MsgRevokeCredential{
		Issuer:       issuer,
		CredentialID: credentialID,
		Reason:       reason,
	}
}

func (msg *MsgRevokeCredential) Route() string { return RouterKey }
func (msg *MsgRevokeCredential) Type() string  { return TypeMsgRevokeCredential }
func (msg *MsgRevokeCredential) GetSigners() []sdk.AccAddress {
	issuer, err := sdk.AccAddressFromBech32(msg.Issuer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{issuer}
}

func (msg *MsgRevokeCredential) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRevokeCredential) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Issuer); err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid issuer address: %s", msg.Issuer)
	}
	if msg.CredentialID == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "credential ID cannot be empty")
	}
	if msg.Reason == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "reason cannot be empty")
	}
	return nil
}

// MsgPresentCredential presents a credential with selective disclosure
type MsgPresentCredential struct {
	Holder          string              `json:"holder"`
	CredentialIDs   []string            `json:"credential_ids"`
	Verifier        string              `json:"verifier"`
	RevealedClaims  map[string][]string `json:"revealed_claims"` // credentialID -> claims
	Challenge       string              `json:"challenge"`
	Domain          string              `json:"domain"`
}

func NewMsgPresentCredential(holder, verifier string) *MsgPresentCredential {
	return &MsgPresentCredential{
		Holder:   holder,
		Verifier: verifier,
	}
}

func (msg *MsgPresentCredential) Route() string { return RouterKey }
func (msg *MsgPresentCredential) Type() string  { return TypeMsgPresentCredential }
func (msg *MsgPresentCredential) GetSigners() []sdk.AccAddress {
	holder, err := sdk.AccAddressFromBech32(msg.Holder)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{holder}
}

func (msg *MsgPresentCredential) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgPresentCredential) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Holder); err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid holder address: %s", msg.Holder)
	}
	if _, err := sdk.AccAddressFromBech32(msg.Verifier); err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid verifier address: %s", msg.Verifier)
	}
	if len(msg.CredentialIDs) == 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "credential IDs cannot be empty")
	}
	return nil
}

// MsgCreateZKProof creates a zero-knowledge proof
type MsgCreateZKProof struct {
	Creator         string             `json:"creator"`
	ProofType       ZKProofType        `json:"proof_type"`
	Statement       string             `json:"statement"`
	CredentialIDs   []string           `json:"credential_ids"`
	ProofData       []byte             `json:"proof_data"`
	PublicInputs    []string           `json:"public_inputs,omitempty"`
	ExpiryMinutes   int32              `json:"expiry_minutes"`
}

func NewMsgCreateZKProof(creator string, proofType ZKProofType) *MsgCreateZKProof {
	return &MsgCreateZKProof{
		Creator:   creator,
		ProofType: proofType,
	}
}

func (msg *MsgCreateZKProof) Route() string { return RouterKey }
func (msg *MsgCreateZKProof) Type() string  { return TypeMsgCreateZKProof }
func (msg *MsgCreateZKProof) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateZKProof) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateZKProof) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address: %s", msg.Creator)
	}
	if msg.Statement == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "statement cannot be empty")
	}
	if len(msg.ProofData) == 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "proof data cannot be empty")
	}
	return nil
}

// MsgVerifyZKProof verifies a zero-knowledge proof
type MsgVerifyZKProof struct {
	Verifier  string `json:"verifier"`
	ProofID   string `json:"proof_id"`
	Challenge string `json:"challenge,omitempty"`
}

func NewMsgVerifyZKProof(verifier, proofID string) *MsgVerifyZKProof {
	return &MsgVerifyZKProof{
		Verifier: verifier,
		ProofID:  proofID,
	}
}

func (msg *MsgVerifyZKProof) Route() string { return RouterKey }
func (msg *MsgVerifyZKProof) Type() string  { return TypeMsgVerifyZKProof }
func (msg *MsgVerifyZKProof) GetSigners() []sdk.AccAddress {
	verifier, err := sdk.AccAddressFromBech32(msg.Verifier)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{verifier}
}

func (msg *MsgVerifyZKProof) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgVerifyZKProof) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Verifier); err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid verifier address: %s", msg.Verifier)
	}
	if msg.ProofID == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "proof ID cannot be empty")
	}
	return nil
}

// MsgLinkAadhaar links Aadhaar to identity
type MsgLinkAadhaar struct {
	Creator            string `json:"creator"`
	AadhaarHash        string `json:"aadhaar_hash"`
	DemographicHash    string `json:"demographic_hash"`
	BiometricHash      string `json:"biometric_hash"`
	VerificationMethod string `json:"verification_method"`
	ConsentArtefact    string `json:"consent_artefact"`
}

func NewMsgLinkAadhaar(creator string) *MsgLinkAadhaar {
	return &MsgLinkAadhaar{
		Creator: creator,
	}
}

func (msg *MsgLinkAadhaar) Route() string { return RouterKey }
func (msg *MsgLinkAadhaar) Type() string  { return TypeMsgLinkAadhaar }
func (msg *MsgLinkAadhaar) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgLinkAadhaar) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgLinkAadhaar) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address: %s", msg.Creator)
	}
	if msg.AadhaarHash == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "Aadhaar hash cannot be empty")
	}
	if msg.ConsentArtefact == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "consent artefact cannot be empty")
	}
	return nil
}

// MsgConnectDigiLocker connects DigiLocker account
type MsgConnectDigiLocker struct {
	Creator        string   `json:"creator"`
	AuthToken      string   `json:"auth_token"`
	ConsentID      string   `json:"consent_id"`
	DocumentTypes  []string `json:"document_types"`
}

func NewMsgConnectDigiLocker(creator string) *MsgConnectDigiLocker {
	return &MsgConnectDigiLocker{
		Creator: creator,
	}
}

func (msg *MsgConnectDigiLocker) Route() string { return RouterKey }
func (msg *MsgConnectDigiLocker) Type() string  { return TypeMsgConnectDigiLocker }
func (msg *MsgConnectDigiLocker) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgConnectDigiLocker) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgConnectDigiLocker) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address: %s", msg.Creator)
	}
	if msg.AuthToken == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "auth token cannot be empty")
	}
	if msg.ConsentID == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "consent ID cannot be empty")
	}
	return nil
}

// MsgLinkUPI links UPI ID to identity
type MsgLinkUPI struct {
	Creator      string `json:"creator"`
	VPAURI       string `json:"vpa_uri"`
	PSPProvider  string `json:"psp_provider"`
	AuthToken    string `json:"auth_token"`
}

func NewMsgLinkUPI(creator, vpaURI string) *MsgLinkUPI {
	return &MsgLinkUPI{
		Creator: creator,
		VPAURI:  vpaURI,
	}
}

func (msg *MsgLinkUPI) Route() string { return RouterKey }
func (msg *MsgLinkUPI) Type() string  { return TypeMsgLinkUPI }
func (msg *MsgLinkUPI) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgLinkUPI) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgLinkUPI) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address: %s", msg.Creator)
	}
	if msg.VPAURI == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "VPA URI cannot be empty")
	}
	if msg.PSPProvider == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "PSP provider cannot be empty")
	}
	return nil
}

// MsgGiveConsent gives consent for data usage
type MsgGiveConsent struct {
	Creator         string      `json:"creator"`
	ConsentType     ConsentType `json:"consent_type"`
	Purpose         string      `json:"purpose"`
	DataController  string      `json:"data_controller"`
	DataCategories  []string    `json:"data_categories"`
	ProcessingTypes []string    `json:"processing_types"`
	ExpirationDays  int32       `json:"expiration_days"`
}

func NewMsgGiveConsent(creator string, consentType ConsentType) *MsgGiveConsent {
	return &MsgGiveConsent{
		Creator:     creator,
		ConsentType: consentType,
	}
}

func (msg *MsgGiveConsent) Route() string { return RouterKey }
func (msg *MsgGiveConsent) Type() string  { return TypeMsgGiveConsent }
func (msg *MsgGiveConsent) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgGiveConsent) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgGiveConsent) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address: %s", msg.Creator)
	}
	if msg.Purpose == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "purpose cannot be empty")
	}
	if msg.DataController == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "data controller cannot be empty")
	}
	if len(msg.DataCategories) == 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "data categories cannot be empty")
	}
	return nil
}

// MsgWithdrawConsent withdraws consent
type MsgWithdrawConsent struct {
	Creator   string `json:"creator"`
	ConsentID string `json:"consent_id"`
	Reason    string `json:"reason"`
}

func NewMsgWithdrawConsent(creator, consentID, reason string) *MsgWithdrawConsent {
	return &MsgWithdrawConsent{
		Creator:   creator,
		ConsentID: consentID,
		Reason:    reason,
	}
}

func (msg *MsgWithdrawConsent) Route() string { return RouterKey }
func (msg *MsgWithdrawConsent) Type() string  { return TypeMsgWithdrawConsent }
func (msg *MsgWithdrawConsent) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgWithdrawConsent) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgWithdrawConsent) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address: %s", msg.Creator)
	}
	if msg.ConsentID == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "consent ID cannot be empty")
	}
	return nil
}

// MsgAddRecoveryMethod adds a recovery method
type MsgAddRecoveryMethod struct {
	Creator  string       `json:"creator"`
	Type     RecoveryType `json:"type"`
	Value    string       `json:"value"`
}

func NewMsgAddRecoveryMethod(creator string, recoveryType RecoveryType, value string) *MsgAddRecoveryMethod {
	return &MsgAddRecoveryMethod{
		Creator: creator,
		Type:    recoveryType,
		Value:   value,
	}
}

func (msg *MsgAddRecoveryMethod) Route() string { return RouterKey }
func (msg *MsgAddRecoveryMethod) Type() string  { return TypeMsgAddRecoveryMethod }
func (msg *MsgAddRecoveryMethod) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgAddRecoveryMethod) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAddRecoveryMethod) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address: %s", msg.Creator)
	}
	if msg.Value == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "recovery value cannot be empty")
	}
	return nil
}

// MsgInitiateRecovery initiates account recovery
type MsgInitiateRecovery struct {
	Creator         string       `json:"creator"`
	TargetAddress   string       `json:"target_address"`
	RecoveryType    RecoveryType `json:"recovery_type"`
	RecoveryData    string       `json:"recovery_data"`
}

func NewMsgInitiateRecovery(creator, targetAddress string, recoveryType RecoveryType) *MsgInitiateRecovery {
	return &MsgInitiateRecovery{
		Creator:       creator,
		TargetAddress: targetAddress,
		RecoveryType:  recoveryType,
	}
}

func (msg *MsgInitiateRecovery) Route() string { return RouterKey }
func (msg *MsgInitiateRecovery) Type() string  { return TypeMsgInitiateRecovery }
func (msg *MsgInitiateRecovery) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgInitiateRecovery) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgInitiateRecovery) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address: %s", msg.Creator)
	}
	if _, err := sdk.AccAddressFromBech32(msg.TargetAddress); err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid target address: %s", msg.TargetAddress)
	}
	return nil
}

// MsgCompleteRecovery completes account recovery
type MsgCompleteRecovery struct {
	Creator       string `json:"creator"`
	RecoveryID    string `json:"recovery_id"`
	NewPublicKey  string `json:"new_public_key"`
	ProofData     string `json:"proof_data"`
}

func NewMsgCompleteRecovery(creator, recoveryID, newPublicKey string) *MsgCompleteRecovery {
	return &MsgCompleteRecovery{
		Creator:      creator,
		RecoveryID:   recoveryID,
		NewPublicKey: newPublicKey,
	}
}

func (msg *MsgCompleteRecovery) Route() string { return RouterKey }
func (msg *MsgCompleteRecovery) Type() string  { return TypeMsgCompleteRecovery }
func (msg *MsgCompleteRecovery) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCompleteRecovery) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCompleteRecovery) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address: %s", msg.Creator)
	}
	if msg.RecoveryID == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "recovery ID cannot be empty")
	}
	if msg.NewPublicKey == "" {
		return sdkerrors.Wrap(ErrInvalidPublicKey, "new public key cannot be empty")
	}
	return nil
}

// MsgUpdatePrivacySettings updates privacy settings
type MsgUpdatePrivacySettings struct {
	Creator                  string            `json:"creator"`
	DefaultDisclosureLevel   DisclosureLevel   `json:"default_disclosure_level"`
	AllowAnonymousUsage      bool              `json:"allow_anonymous_usage"`
	RequireExplicitConsent   bool              `json:"require_explicit_consent"`
	DataMinimization         bool              `json:"data_minimization"`
	AutoDeleteAfterDays      int32             `json:"auto_delete_after_days"`
	AllowDerivedCredentials  bool              `json:"allow_derived_credentials"`
	PreferredProofSystems    []string          `json:"preferred_proof_systems"`
	BlacklistedVerifiers     []string          `json:"blacklisted_verifiers"`
}

func NewMsgUpdatePrivacySettings(creator string) *MsgUpdatePrivacySettings {
	return &MsgUpdatePrivacySettings{
		Creator: creator,
	}
}

func (msg *MsgUpdatePrivacySettings) Route() string { return RouterKey }
func (msg *MsgUpdatePrivacySettings) Type() string  { return TypeMsgUpdatePrivacySettings }
func (msg *MsgUpdatePrivacySettings) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdatePrivacySettings) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdatePrivacySettings) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address: %s", msg.Creator)
	}
	if msg.AutoDeleteAfterDays < 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "auto delete days cannot be negative")
	}
	return nil
}