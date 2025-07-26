package types

// Message response types

// MsgCreateIdentityResponse defines the response for creating an identity
type MsgCreateIdentityResponse struct {
	Did string `json:"did"`
}

// MsgUpdateIdentityResponse defines the response for updating an identity
type MsgUpdateIdentityResponse struct{}

// MsgRevokeIdentityResponse defines the response for revoking an identity
type MsgRevokeIdentityResponse struct{}

// MsgRegisterDIDResponse defines the response for registering a DID
type MsgRegisterDIDResponse struct{}

// MsgUpdateDIDResponse defines the response for updating a DID
type MsgUpdateDIDResponse struct{}

// MsgDeactivateDIDResponse defines the response for deactivating a DID
type MsgDeactivateDIDResponse struct{}

// MsgIssueCredentialResponse defines the response for issuing a credential
type MsgIssueCredentialResponse struct {
	CredentialId string `json:"credential_id"`
}

// MsgRevokeCredentialResponse defines the response for revoking a credential
type MsgRevokeCredentialResponse struct{}

// MsgPresentCredentialResponse defines the response for presenting credentials
type MsgPresentCredentialResponse struct {
	PresentationId string `json:"presentation_id"`
}

// MsgCreateZKProofResponse defines the response for creating a ZK proof
type MsgCreateZKProofResponse struct {
	ProofId string `json:"proof_id"`
}

// MsgVerifyZKProofResponse defines the response for verifying a ZK proof
type MsgVerifyZKProofResponse struct {
	Valid bool   `json:"valid"`
	Error string `json:"error,omitempty"`
}

// MsgLinkAadhaarResponse defines the response for linking Aadhaar
type MsgLinkAadhaarResponse struct{}

// MsgConnectDigiLockerResponse defines the response for connecting DigiLocker
type MsgConnectDigiLockerResponse struct {
	DocumentCount int32 `json:"document_count"`
}

// MsgLinkUPIResponse defines the response for linking UPI
type MsgLinkUPIResponse struct{}

// MsgGiveConsentResponse defines the response for giving consent
type MsgGiveConsentResponse struct {
	ConsentId string `json:"consent_id"`
}

// MsgWithdrawConsentResponse defines the response for withdrawing consent
type MsgWithdrawConsentResponse struct{}

// MsgAddRecoveryMethodResponse defines the response for adding a recovery method
type MsgAddRecoveryMethodResponse struct{}

// MsgInitiateRecoveryResponse defines the response for initiating recovery
type MsgInitiateRecoveryResponse struct {
	RecoveryId string `json:"recovery_id"`
}

// MsgCompleteRecoveryResponse defines the response for completing recovery
type MsgCompleteRecoveryResponse struct{}

// MsgUpdatePrivacySettingsResponse defines the response for updating privacy settings
type MsgUpdatePrivacySettingsResponse struct{}