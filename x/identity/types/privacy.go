package types

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// ZKProof represents a zero-knowledge proof
type ZKProof struct {
	ID              string            `json:"id"`
	Type            ZKProofType       `json:"type"`
	ProofSystem     string            `json:"proof_system"`
	Statement       string            `json:"statement"`
	Commitment      string            `json:"commitment"`
	Challenge       string            `json:"challenge"`
	Response        string            `json:"response"`
	PublicInputs    []string          `json:"public_inputs,omitempty"`
	VerifierKey     string            `json:"verifier_key"`
	ProofData       []byte            `json:"proof_data"`
	CreatedAt       time.Time         `json:"created_at"`
	ExpiresAt       *time.Time        `json:"expires_at,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
}

// ZKProofType represents different types of zero-knowledge proofs
type ZKProofType int32

const (
	ZKProofType_AGE_RANGE        ZKProofType = 0
	ZKProofType_BALANCE_RANGE    ZKProofType = 1
	ZKProofType_MEMBERSHIP       ZKProofType = 2
	ZKProofType_NON_MEMBERSHIP   ZKProofType = 3
	ZKProofType_CREDENTIAL_OWNERSHIP ZKProofType = 4
	ZKProofType_SELECTIVE_DISCLOSURE ZKProofType = 5
)

// String returns the string representation
func (zpt ZKProofType) String() string {
	switch zpt {
	case ZKProofType_AGE_RANGE:
		return "AGE_RANGE"
	case ZKProofType_BALANCE_RANGE:
		return "BALANCE_RANGE"
	case ZKProofType_MEMBERSHIP:
		return "MEMBERSHIP"
	case ZKProofType_NON_MEMBERSHIP:
		return "NON_MEMBERSHIP"
	case ZKProofType_CREDENTIAL_OWNERSHIP:
		return "CREDENTIAL_OWNERSHIP"
	case ZKProofType_SELECTIVE_DISCLOSURE:
		return "SELECTIVE_DISCLOSURE"
	default:
		return "UNKNOWN"
	}
}

// PrivacyPreservingCredential represents a privacy-enhanced credential
type PrivacyPreservingCredential struct {
	ID                    string                   `json:"id"`
	OriginalCredentialID  string                   `json:"original_credential_id"`
	HolderCommitment      string                   `json:"holder_commitment"`
	IssuerSignature       string                   `json:"issuer_signature"`
	AttributeCommitments  map[string]string        `json:"attribute_commitments"`
	BlindedAttributes     []string                 `json:"blinded_attributes"`
	RevocationHandle      string                   `json:"revocation_handle"`
	AccumulatorValue      string                   `json:"accumulator_value"`
	CreatedAt             time.Time                `json:"created_at"`
	Metadata              map[string]string        `json:"metadata,omitempty"`
}

// AnonymousCredential represents a fully anonymous credential
type AnonymousCredential struct {
	ID                string            `json:"id"`
	Type              string            `json:"type"`
	BlindSignature    string            `json:"blind_signature"`
	Nullifier         string            `json:"nullifier"`
	MerkleRoot        string            `json:"merkle_root"`
	MerkleProof       []string          `json:"merkle_proof"`
	ValidityProof     *ZKProof          `json:"validity_proof"`
	IssuanceTimestamp time.Time         `json:"issuance_timestamp"`
	Metadata          map[string]string `json:"metadata,omitempty"`
}

// PrivacyRequest represents a request for privacy-preserving operations
type PrivacyRequest struct {
	ID              string                 `json:"id"`
	Type            PrivacyRequestType     `json:"type"`
	Requester       string                 `json:"requester"`
	Subject         string                 `json:"subject"`
	RequestedClaims []string               `json:"requested_claims"`
	ProofRequirements []ProofRequirement   `json:"proof_requirements"`
	Challenge       string                 `json:"challenge"`
	ValidUntil      time.Time              `json:"valid_until"`
	Metadata        map[string]string      `json:"metadata,omitempty"`
}

// PrivacyRequestType represents different types of privacy requests
type PrivacyRequestType int32

const (
	PrivacyRequestType_SELECTIVE_DISCLOSURE PrivacyRequestType = 0
	PrivacyRequestType_RANGE_PROOF         PrivacyRequestType = 1
	PrivacyRequestType_SET_MEMBERSHIP      PrivacyRequestType = 2
	PrivacyRequestType_PREDICATE_PROOF     PrivacyRequestType = 3
)

// String returns the string representation
func (prt PrivacyRequestType) String() string {
	switch prt {
	case PrivacyRequestType_SELECTIVE_DISCLOSURE:
		return "SELECTIVE_DISCLOSURE"
	case PrivacyRequestType_RANGE_PROOF:
		return "RANGE_PROOF"
	case PrivacyRequestType_SET_MEMBERSHIP:
		return "SET_MEMBERSHIP"
	case PrivacyRequestType_PREDICATE_PROOF:
		return "PREDICATE_PROOF"
	default:
		return "UNKNOWN"
	}
}

// ProofRequirement defines what needs to be proven
type ProofRequirement struct {
	ID              string            `json:"id"`
	Type            string            `json:"type"`
	Attribute       string            `json:"attribute"`
	Condition       string            `json:"condition"`
	Value           interface{}       `json:"value,omitempty"`
	Range           *Range            `json:"range,omitempty"`
	Set             []string          `json:"set,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
}

// Range represents a numeric range for range proofs
type Range struct {
	Min *int64 `json:"min,omitempty"`
	Max *int64 `json:"max,omitempty"`
}

// PrivacyResponse represents a response to a privacy request
type PrivacyResponse struct {
	ID              string                   `json:"id"`
	RequestID       string                   `json:"request_id"`
	Proofs          []ZKProof                `json:"proofs"`
	RevealedClaims  map[string]interface{}   `json:"revealed_claims,omitempty"`
	Commitments     map[string]string        `json:"commitments"`
	Timestamp       time.Time                `json:"timestamp"`
	Metadata        map[string]string        `json:"metadata,omitempty"`
}

// DerivedCredential represents a credential derived from another with privacy
type DerivedCredential struct {
	ID                   string            `json:"id"`
	SourceCredentialID   string            `json:"source_credential_id"`
	DerivedClaims        []string          `json:"derived_claims"`
	TransformationProof  *ZKProof          `json:"transformation_proof"`
	LinkabilityToken     string            `json:"linkability_token"`
	ValidityPeriod       time.Duration     `json:"validity_period"`
	CreatedAt            time.Time         `json:"created_at"`
	Metadata             map[string]string `json:"metadata,omitempty"`
}

// PrivacySettings represents user privacy preferences
type PrivacySettings struct {
	UserAddress              string            `json:"user_address"`
	DefaultDisclosureLevel   DisclosureLevel   `json:"default_disclosure_level"`
	AllowAnonymousUsage      bool              `json:"allow_anonymous_usage"`
	RequireExplicitConsent   bool              `json:"require_explicit_consent"`
	DataMinimization         bool              `json:"data_minimization"`
	AutoDeleteAfterDays      int32             `json:"auto_delete_after_days"`
	AllowDerivedCredentials  bool              `json:"allow_derived_credentials"`
	PreferredProofSystems    []string          `json:"preferred_proof_systems"`
	BlacklistedVerifiers     []string          `json:"blacklisted_verifiers"`
	UpdatedAt                time.Time         `json:"updated_at"`
	Metadata                 map[string]string `json:"metadata,omitempty"`
}

// DisclosureLevel represents different levels of information disclosure
type DisclosureLevel int32

const (
	DisclosureLevel_MINIMAL     DisclosureLevel = 0
	DisclosureLevel_STANDARD    DisclosureLevel = 1
	DisclosureLevel_ENHANCED    DisclosureLevel = 2
	DisclosureLevel_FULL        DisclosureLevel = 3
)

// String returns the string representation
func (dl DisclosureLevel) String() string {
	switch dl {
	case DisclosureLevel_MINIMAL:
		return "MINIMAL"
	case DisclosureLevel_STANDARD:
		return "STANDARD"
	case DisclosureLevel_ENHANCED:
		return "ENHANCED"
	case DisclosureLevel_FULL:
		return "FULL"
	default:
		return "UNKNOWN"
	}
}

// VerifyZKProof performs basic validation on a ZK proof
func VerifyZKProof(proof *ZKProof) error {
	if proof.ID == "" {
		return fmt.Errorf("proof must have an ID")
	}
	
	if proof.ProofSystem == "" {
		return fmt.Errorf("proof must specify proof system")
	}
	
	if proof.Statement == "" {
		return fmt.Errorf("proof must have a statement")
	}
	
	if len(proof.ProofData) == 0 {
		return fmt.Errorf("proof must have proof data")
	}
	
	// Check if proof is expired
	if proof.ExpiresAt != nil && proof.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("proof has expired")
	}
	
	return nil
}

// GenerateNullifier generates a unique nullifier for double-spend prevention
func GenerateNullifier(credentialID string, usage string) string {
	data := fmt.Sprintf("%s:%s:%d", credentialID, usage, time.Now().Unix())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// CreateCommitment creates a cryptographic commitment to a value
func CreateCommitment(value interface{}, randomness string) string {
	data := fmt.Sprintf("%v:%s", value, randomness)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// Privacy-related constants
const (
	// Proof systems
	ProofSystemGroth16      = "groth16"
	ProofSystemPLONK        = "plonk"
	ProofSystemBulletproofs = "bulletproofs"
	ProofSystemSTARK        = "stark"
	
	// Privacy events
	EventTypeZKProofCreated     = "zk_proof_created"
	EventTypeZKProofVerified    = "zk_proof_verified"
	EventTypePrivacyRequest     = "privacy_request"
	EventTypePrivacyResponse    = "privacy_response"
	EventTypeAnonymousAction    = "anonymous_action"
	
	// Privacy attributes
	AttributeKeyProofType       = "proof_type"
	AttributeKeyProofSystem     = "proof_system"
	AttributeKeyDisclosureLevel = "disclosure_level"
	AttributeKeyAnonymitySet    = "anonymity_set"
)