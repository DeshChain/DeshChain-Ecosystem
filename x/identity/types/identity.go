package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Identity represents a complete identity record in DeshChain
type Identity struct {
	Address          string              `json:"address"`           // DeshChain address
	DID              string              `json:"did"`               // Decentralized Identifier
	DIDDocument      *DIDDocument        `json:"did_document"`      // Full DID Document
	Credentials      []string            `json:"credentials"`       // List of credential IDs
	Status           IdentityStatus      `json:"status"`            // Active, suspended, revoked
	KYCStatus        KYCStatus           `json:"kyc_status"`        // KYC verification status
	BiometricStatus  BiometricStatus     `json:"biometric_status"`  // Biometric enrollment status
	RecoveryMethods  []RecoveryMethod    `json:"recovery_methods"`  // Account recovery options
	Consents         []ConsentRecord     `json:"consents"`          // Privacy consents
	CreatedAt        time.Time           `json:"created_at"`
	UpdatedAt        time.Time           `json:"updated_at"`
	LastActivityAt   time.Time           `json:"last_activity_at"`
	Metadata         map[string]string   `json:"metadata,omitempty"`
}

// IdentityStatus represents the status of an identity
type IdentityStatus int32

const (
	IdentityStatus_PENDING    IdentityStatus = 0
	IdentityStatus_ACTIVE     IdentityStatus = 1
	IdentityStatus_SUSPENDED  IdentityStatus = 2
	IdentityStatus_REVOKED    IdentityStatus = 3
)

// String returns the string representation
func (is IdentityStatus) String() string {
	switch is {
	case IdentityStatus_PENDING:
		return "PENDING"
	case IdentityStatus_ACTIVE:
		return "ACTIVE"
	case IdentityStatus_SUSPENDED:
		return "SUSPENDED"
	case IdentityStatus_REVOKED:
		return "REVOKED"
	default:
		return "UNKNOWN"
	}
}

// KYCStatus represents KYC verification status
type KYCStatus struct {
	Level            KYCLevel          `json:"level"`
	Status           VerificationStatus `json:"status"`
	VerifiedAt       time.Time         `json:"verified_at"`
	ExpiresAt        time.Time         `json:"expires_at"`
	Verifier         string            `json:"verifier"`
	CredentialID     string            `json:"credential_id"`
	RiskScore        float64           `json:"risk_score"`
	ComplianceFlags  []string          `json:"compliance_flags,omitempty"`
}

// KYCLevel represents different levels of KYC verification
type KYCLevel int32

const (
	KYCLevel_NONE     KYCLevel = 0
	KYCLevel_BASIC    KYCLevel = 1
	KYCLevel_ENHANCED KYCLevel = 2
	KYCLevel_FULL     KYCLevel = 3
)

// String returns the string representation
func (kl KYCLevel) String() string {
	switch kl {
	case KYCLevel_NONE:
		return "NONE"
	case KYCLevel_BASIC:
		return "BASIC"
	case KYCLevel_ENHANCED:
		return "ENHANCED"
	case KYCLevel_FULL:
		return "FULL"
	default:
		return "UNKNOWN"
	}
}

// BiometricStatus represents biometric enrollment status
type BiometricStatus struct {
	Enrolled         bool              `json:"enrolled"`
	BiometricTypes   []string          `json:"biometric_types"`
	LastEnrollment   time.Time         `json:"last_enrollment"`
	CredentialIDs    []string          `json:"credential_ids"`
	SecurityLevel    string            `json:"security_level"`
	DeviceCount      int32             `json:"device_count"`
}

// VerificationStatus represents the status of a verification
type VerificationStatus int32

const (
	VerificationStatus_NOT_VERIFIED VerificationStatus = 0
	VerificationStatus_PENDING      VerificationStatus = 1
	VerificationStatus_VERIFIED     VerificationStatus = 2
	VerificationStatus_FAILED       VerificationStatus = 3
	VerificationStatus_EXPIRED      VerificationStatus = 4
)

// String returns the string representation
func (vs VerificationStatus) String() string {
	switch vs {
	case VerificationStatus_NOT_VERIFIED:
		return "NOT_VERIFIED"
	case VerificationStatus_PENDING:
		return "PENDING"
	case VerificationStatus_VERIFIED:
		return "VERIFIED"
	case VerificationStatus_FAILED:
		return "FAILED"
	case VerificationStatus_EXPIRED:
		return "EXPIRED"
	default:
		return "UNKNOWN"
	}
}

// RecoveryMethod represents a method for account recovery
type RecoveryMethod struct {
	Type        RecoveryType `json:"type"`
	Value       string       `json:"value"`        // Encrypted/hashed value
	AddedAt     time.Time    `json:"added_at"`
	LastUsed    time.Time    `json:"last_used,omitempty"`
	IsActive    bool         `json:"is_active"`
}

// RecoveryType represents different recovery methods
type RecoveryType int32

const (
	RecoveryType_EMAIL       RecoveryType = 0
	RecoveryType_PHONE       RecoveryType = 1
	RecoveryType_GUARDIAN    RecoveryType = 2
	RecoveryType_SOCIAL      RecoveryType = 3
	RecoveryType_MNEMONIC    RecoveryType = 4
)

// String returns the string representation
func (rt RecoveryType) String() string {
	switch rt {
	case RecoveryType_EMAIL:
		return "EMAIL"
	case RecoveryType_PHONE:
		return "PHONE"
	case RecoveryType_GUARDIAN:
		return "GUARDIAN"
	case RecoveryType_SOCIAL:
		return "SOCIAL"
	case RecoveryType_MNEMONIC:
		return "MNEMONIC"
	default:
		return "UNKNOWN"
	}
}

// ConsentRecord represents a privacy consent
type ConsentRecord struct {
	ID              string            `json:"id"`
	Type            ConsentType       `json:"type"`
	Purpose         string            `json:"purpose"`
	DataController  string            `json:"data_controller"`
	DataCategories  []string          `json:"data_categories"`
	ProcessingTypes []string          `json:"processing_types"`
	Given           bool              `json:"given"`
	GivenAt         time.Time         `json:"given_at"`
	ExpiresAt       *time.Time        `json:"expires_at,omitempty"`
	WithdrawnAt     *time.Time        `json:"withdrawn_at,omitempty"`
	Version         string            `json:"version"`
}

// ConsentType represents different types of consent
type ConsentType int32

const (
	ConsentType_DATA_COLLECTION ConsentType = 0
	ConsentType_DATA_SHARING    ConsentType = 1
	ConsentType_MARKETING       ConsentType = 2
	ConsentType_ANALYTICS       ConsentType = 3
	ConsentType_BIOMETRIC       ConsentType = 4
	ConsentType_KYC             ConsentType = 5
)

// String returns the string representation
func (ct ConsentType) String() string {
	switch ct {
	case ConsentType_DATA_COLLECTION:
		return "DATA_COLLECTION"
	case ConsentType_DATA_SHARING:
		return "DATA_SHARING"
	case ConsentType_MARKETING:
		return "MARKETING"
	case ConsentType_ANALYTICS:
		return "ANALYTICS"
	case ConsentType_BIOMETRIC:
		return "BIOMETRIC"
	case ConsentType_KYC:
		return "KYC"
	default:
		return "UNKNOWN"
	}
}

// IdentityRequest represents a request to create or update identity
type IdentityRequest struct {
	Address         string                 `json:"address"`
	PublicKey       string                 `json:"public_key"`
	DIDMethod       string                 `json:"did_method"`
	ServiceEndpoints []Service             `json:"service_endpoints,omitempty"`
	RecoveryMethods []RecoveryMethod      `json:"recovery_methods"`
	InitialConsents []ConsentRecord       `json:"initial_consents"`
	Metadata        map[string]string     `json:"metadata,omitempty"`
}

// CredentialIssuanceRequest represents a request to issue a credential
type CredentialIssuanceRequest struct {
	Issuer           string                 `json:"issuer"`
	Holder           string                 `json:"holder"`
	CredentialType   string                 `json:"credential_type"`
	Claims           map[string]interface{} `json:"claims"`
	ExpirationDays   int32                  `json:"expiration_days"`
	Evidence         []Evidence             `json:"evidence,omitempty"`
	RequireConsent   bool                   `json:"require_consent"`
	Metadata         map[string]string      `json:"metadata,omitempty"`
}

// IdentityEventType defines event types for identity module
const (
	EventTypeIdentityCreated      = "identity_created"
	EventTypeIdentityUpdated      = "identity_updated"
	EventTypeIdentityRevoked      = "identity_revoked"
	EventTypeDIDRegistered        = "did_registered"
	EventTypeDIDUpdated           = "did_updated"
	EventTypeCredentialIssued     = "credential_issued"
	EventTypeCredentialRevoked    = "credential_revoked"
	EventTypeCredentialPresented  = "credential_presented"
	EventTypeKYCCompleted         = "kyc_completed"
	EventTypeBiometricEnrolled    = "biometric_enrolled"
	EventTypeConsentGiven         = "consent_given"
	EventTypeConsentWithdrawn     = "consent_withdrawn"
	EventTypeRecoveryInitiated    = "recovery_initiated"
	EventTypeRecoveryCompleted    = "recovery_completed"
)

// AttributeKey defines attribute keys for events
const (
	AttributeKeyAddress       = "address"
	AttributeKeyDID           = "did"
	AttributeKeyCredentialID  = "credential_id"
	AttributeKeyCredentialType = "credential_type"
	AttributeKeyIssuer        = "issuer"
	AttributeKeyHolder        = "holder"
	AttributeKeyKYCLevel      = "kyc_level"
	AttributeKeyBiometricType = "biometric_type"
	AttributeKeyConsentType   = "consent_type"
	AttributeKeyRecoveryType  = "recovery_type"
	AttributeKeyStatus        = "status"
)

// Validate performs validation on Identity
func (i *Identity) Validate() error {
	if i.Address == "" {
		return fmt.Errorf("identity must have an address")
	}
	
	if _, err := sdk.AccAddressFromBech32(i.Address); err != nil {
		return fmt.Errorf("invalid address: %w", err)
	}
	
	if i.DID == "" {
		return fmt.Errorf("identity must have a DID")
	}
	
	if i.DIDDocument == nil {
		return fmt.Errorf("identity must have a DID document")
	}
	
	if err := ValidateDIDDocument(i.DIDDocument); err != nil {
		return fmt.Errorf("invalid DID document: %w", err)
	}
	
	return nil
}

// IsActive returns true if identity is active
func (i *Identity) IsActive() bool {
	return i.Status == IdentityStatus_ACTIVE
}

// HasKYC returns true if identity has completed KYC
func (i *Identity) HasKYC(minLevel KYCLevel) bool {
	return i.KYCStatus.Status == VerificationStatus_VERIFIED && 
		i.KYCStatus.Level >= minLevel &&
		i.KYCStatus.ExpiresAt.After(time.Now())
}

// HasBiometrics returns true if identity has biometric enrollment
func (i *Identity) HasBiometrics() bool {
	return i.BiometricStatus.Enrolled && len(i.BiometricStatus.BiometricTypes) > 0
}

// GetActiveConsents returns all active consents
func (i *Identity) GetActiveConsents() []ConsentRecord {
	var active []ConsentRecord
	for _, consent := range i.Consents {
		if consent.Given && consent.WithdrawnAt == nil {
			if consent.ExpiresAt == nil || consent.ExpiresAt.After(time.Now()) {
				active = append(active, consent)
			}
		}
	}
	return active
}

// HasConsent checks if identity has given consent for a specific type
func (i *Identity) HasConsent(consentType ConsentType) bool {
	for _, consent := range i.GetActiveConsents() {
		if consent.Type == consentType {
			return true
		}
	}
	return false
}

// RecoveryRequest represents an account recovery request
type RecoveryRequest struct {
	ID              string            `json:"id"`
	TargetAddress   string            `json:"target_address"`
	InitiatorAddress string           `json:"initiator_address"`
	RecoveryType    RecoveryType      `json:"recovery_type"`
	RecoveryData    string            `json:"recovery_data"`
	Status          string            `json:"status"` // pending, ready, completed, rejected
	InitiatedAt     time.Time         `json:"initiated_at"`
	CompletedAt     *time.Time        `json:"completed_at,omitempty"`
	NewPublicKey    string            `json:"new_public_key,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
}