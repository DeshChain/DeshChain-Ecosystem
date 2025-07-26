package types

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Offline verification structures and types for DeshChain Identity

// OfflineVerificationMode defines the mode of offline verification
type OfflineVerificationMode string

const (
	// Verification modes
	OfflineModeFull       OfflineVerificationMode = "full"        // Complete offline verification
	OfflineModePartial    OfflineVerificationMode = "partial"     // Partial verification with cached data
	OfflineModeMinimal    OfflineVerificationMode = "minimal"     // Basic identity checks only
	OfflineModeEmergency  OfflineVerificationMode = "emergency"   // Emergency mode with reduced security
)

// OfflineCredentialFormat defines the format for offline credentials
type OfflineCredentialFormat string

const (
	// Credential formats
	FormatSelfContained   OfflineCredentialFormat = "self_contained"   // All data embedded
	FormatCompressed      OfflineCredentialFormat = "compressed"       // Compressed format
	FormatQRCode          OfflineCredentialFormat = "qr_code"          // QR code encodable
	FormatNFC             OfflineCredentialFormat = "nfc"              // NFC tag format
	FormatPrintable       OfflineCredentialFormat = "printable"       // Human-readable printable
)

// OfflineVerificationData contains all data needed for offline verification
type OfflineVerificationData struct {
	// Core identity information
	DID               string                     `json:"did" yaml:"did"`
	IdentityHash      string                     `json:"identity_hash" yaml:"identity_hash"`
	PublicKey         string                     `json:"public_key" yaml:"public_key"`
	
	// Verification metadata
	IssuedAt          time.Time                  `json:"issued_at" yaml:"issued_at"`
	ExpiresAt         time.Time                  `json:"expires_at" yaml:"expires_at"`
	IssuerDID         string                     `json:"issuer_did" yaml:"issuer_did"`
	VerificationLevel uint32                     `json:"verification_level" yaml:"verification_level"`
	
	// Cryptographic proofs
	IdentityProof     *CryptographicProof        `json:"identity_proof" yaml:"identity_proof"`
	KYCProof          *CryptographicProof        `json:"kyc_proof,omitempty" yaml:"kyc_proof,omitempty"`
	BiometricProof    *CryptographicProof        `json:"biometric_proof,omitempty" yaml:"biometric_proof,omitempty"`
	
	// Cached credentials for offline use
	Credentials       []*OfflineCredential       `json:"credentials" yaml:"credentials"`
	
	// Revocation information
	RevocationData    *RevocationData            `json:"revocation_data" yaml:"revocation_data"`
	
	// Biometric templates (encrypted)
	BiometricTemplates map[string]*EncryptedBiometric `json:"biometric_templates,omitempty" yaml:"biometric_templates,omitempty"`
	
	// Emergency contacts and recovery
	EmergencyContacts  []*EmergencyContact        `json:"emergency_contacts,omitempty" yaml:"emergency_contacts,omitempty"`
	
	// Localization for offline use
	LocalizedData     *OfflineLocalizationData   `json:"localized_data,omitempty" yaml:"localized_data,omitempty"`
	
	// Format and compression information
	Format            OfflineCredentialFormat    `json:"format" yaml:"format"`
	Compressed        bool                       `json:"compressed" yaml:"compressed"`
	CompressionLevel  uint32                     `json:"compression_level,omitempty" yaml:"compression_level,omitempty"`
	
	// Integrity verification
	DataHash          string                     `json:"data_hash" yaml:"data_hash"`
	Signature         string                     `json:"signature" yaml:"signature"`
}

// CryptographicProof contains cryptographic proof for offline verification
type CryptographicProof struct {
	ProofType         string                     `json:"proof_type" yaml:"proof_type"`
	ProofValue        string                     `json:"proof_value" yaml:"proof_value"`
	ProofPurpose      string                     `json:"proof_purpose" yaml:"proof_purpose"`
	Created           time.Time                  `json:"created" yaml:"created"`
	VerificationMethod string                    `json:"verification_method" yaml:"verification_method"`
	Challenge         string                     `json:"challenge,omitempty" yaml:"challenge,omitempty"`
	Domain            string                     `json:"domain,omitempty" yaml:"domain,omitempty"`
	Nonce             string                     `json:"nonce,omitempty" yaml:"nonce,omitempty"`
}

// OfflineCredential represents a credential optimized for offline use
type OfflineCredential struct {
	ID                string                     `json:"id" yaml:"id"`
	Type              []string                   `json:"type" yaml:"type"`
	Issuer            string                     `json:"issuer" yaml:"issuer"`
	IssuanceDate      time.Time                  `json:"issuance_date" yaml:"issuance_date"`
	ExpirationDate    *time.Time                 `json:"expiration_date,omitempty" yaml:"expiration_date,omitempty"`
	CredentialSubject map[string]interface{}     `json:"credential_subject" yaml:"credential_subject"`
	Proof             *CryptographicProof        `json:"proof" yaml:"proof"`
	
	// Offline-specific fields
	OfflineUsable     bool                       `json:"offline_usable" yaml:"offline_usable"`
	CompressionLevel  uint32                     `json:"compression_level,omitempty" yaml:"compression_level,omitempty"`
	Priority          uint32                     `json:"priority" yaml:"priority"` // Higher priority credentials checked first
}

// RevocationData contains information for checking credential revocation offline
type RevocationData struct {
	RevocationListURL    string                  `json:"revocation_list_url" yaml:"revocation_list_url"`
	RevocationListHash   string                  `json:"revocation_list_hash" yaml:"revocation_list_hash"`
	LastUpdated          time.Time               `json:"last_updated" yaml:"last_updated"`
	ValidUntil           time.Time               `json:"valid_until" yaml:"valid_until"`
	
	// Cached revocation entries for offline use
	RevokedCredentials   []string                `json:"revoked_credentials" yaml:"revoked_credentials"`
	RevokedIdentities    []string                `json:"revoked_identities" yaml:"revoked_identities"`
	
	// Merkle tree root for efficient verification
	MerkleRoot           string                  `json:"merkle_root" yaml:"merkle_root"`
	MerkleProofs         map[string]string       `json:"merkle_proofs,omitempty" yaml:"merkle_proofs,omitempty"`
}

// EncryptedBiometric contains encrypted biometric template for offline matching
type EncryptedBiometric struct {
	BiometricType     BiometricType              `json:"biometric_type" yaml:"biometric_type"`
	EncryptedTemplate string                     `json:"encrypted_template" yaml:"encrypted_template"`
	EncryptionMethod  string                     `json:"encryption_method" yaml:"encryption_method"`
	TemplateHash      string                     `json:"template_hash" yaml:"template_hash"`
	Quality           float64                    `json:"quality" yaml:"quality"`
	CreatedAt         time.Time                  `json:"created_at" yaml:"created_at"`
	ExpiresAt         time.Time                  `json:"expires_at" yaml:"expires_at"`
}

// EmergencyContact contains emergency contact information for offline scenarios
type EmergencyContact struct {
	Name              string                     `json:"name" yaml:"name"`
	Relationship      string                     `json:"relationship" yaml:"relationship"`
	ContactMethod     string                     `json:"contact_method" yaml:"contact_method"`
	VerificationKey   string                     `json:"verification_key" yaml:"verification_key"`
	Priority          uint32                     `json:"priority" yaml:"priority"`
}

// OfflineLocalizationData contains localized strings for offline use
type OfflineLocalizationData struct {
	DefaultLanguage   LanguageCode               `json:"default_language" yaml:"default_language"`
	Messages          map[string]map[LanguageCode]string `json:"messages" yaml:"messages"`
	ErrorMessages     map[string]map[LanguageCode]string `json:"error_messages" yaml:"error_messages"`
}

// OfflineVerificationRequest represents a request for offline verification
type OfflineVerificationRequest struct {
	DID               string                     `json:"did" yaml:"did"`
	Challenge         string                     `json:"challenge" yaml:"challenge"`
	RequiredLevel     uint32                     `json:"required_level" yaml:"required_level"`
	RequiredTypes     []string                   `json:"required_types" yaml:"required_types"`
	BiometricData     map[string]string          `json:"biometric_data,omitempty" yaml:"biometric_data,omitempty"`
	Context           string                     `json:"context,omitempty" yaml:"context,omitempty"`
	Timestamp         time.Time                  `json:"timestamp" yaml:"timestamp"`
	Nonce             string                     `json:"nonce" yaml:"nonce"`
}

// OfflineVerificationResult represents the result of offline verification
type OfflineVerificationResult struct {
	Success           bool                       `json:"success" yaml:"success"`
	DID               string                     `json:"did" yaml:"did"`
	VerificationLevel uint32                     `json:"verification_level" yaml:"verification_level"`
	VerifiedTypes     []string                   `json:"verified_types" yaml:"verified_types"`
	Timestamp         time.Time                  `json:"timestamp" yaml:"timestamp"`
	
	// Verification details
	IdentityVerified  bool                       `json:"identity_verified" yaml:"identity_verified"`
	KYCVerified       bool                       `json:"kyc_verified" yaml:"kyc_verified"`
	BiometricVerified bool                       `json:"biometric_verified" yaml:"biometric_verified"`
	CredentialsValid  bool                       `json:"credentials_valid" yaml:"credentials_valid"`
	NotRevoked        bool                       `json:"not_revoked" yaml:"not_revoked"`
	
	// Error information
	Errors            []string                   `json:"errors,omitempty" yaml:"errors,omitempty"`
	Warnings          []string                   `json:"warnings,omitempty" yaml:"warnings,omitempty"`
	
	// Confidence scores
	IdentityConfidence  float64                  `json:"identity_confidence" yaml:"identity_confidence"`
	BiometricConfidence float64                  `json:"biometric_confidence,omitempty" yaml:"biometric_confidence,omitempty"`
	OverallConfidence   float64                  `json:"overall_confidence" yaml:"overall_confidence"`
	
	// Additional metadata
	VerificationMode  OfflineVerificationMode    `json:"verification_mode" yaml:"verification_mode"`
	CacheAge          time.Duration              `json:"cache_age,omitempty" yaml:"cache_age,omitempty"`
	LastOnlineSync    *time.Time                 `json:"last_online_sync,omitempty" yaml:"last_online_sync,omitempty"`
}

// OfflineVerificationConfig contains configuration for offline verification
type OfflineVerificationConfig struct {
	// Security settings
	MaxOfflineDuration    time.Duration            `json:"max_offline_duration" yaml:"max_offline_duration"`
	RequiredConfidence    float64                  `json:"required_confidence" yaml:"required_confidence"`
	BiometricThreshold    float64                  `json:"biometric_threshold" yaml:"biometric_threshold"`
	
	// Cache settings
	MaxCacheSize          uint64                   `json:"max_cache_size" yaml:"max_cache_size"`
	CacheExpirationPeriod time.Duration            `json:"cache_expiration_period" yaml:"cache_expiration_period"`
	
	// Compression settings
	EnableCompression     bool                     `json:"enable_compression" yaml:"enable_compression"`
	CompressionLevel      uint32                   `json:"compression_level" yaml:"compression_level"`
	
	// Fallback settings
	EmergencyModeEnabled  bool                     `json:"emergency_mode_enabled" yaml:"emergency_mode_enabled"`
	EmergencyThreshold    float64                  `json:"emergency_threshold" yaml:"emergency_threshold"`
	
	// Regional settings
	SupportedRegions      []string                 `json:"supported_regions" yaml:"supported_regions"`
	DefaultLanguage       LanguageCode             `json:"default_language" yaml:"default_language"`
}

// Constructor functions

// NewOfflineVerificationData creates a new offline verification data structure
func NewOfflineVerificationData(did, identityHash, publicKey, issuerDID string, level uint32) *OfflineVerificationData {
	now := time.Now()
	return &OfflineVerificationData{
		DID:               did,
		IdentityHash:      identityHash,
		PublicKey:         publicKey,
		IssuedAt:          now,
		ExpiresAt:         now.Add(24 * time.Hour), // Default 24 hour expiration
		IssuerDID:         issuerDID,
		VerificationLevel: level,
		Credentials:       make([]*OfflineCredential, 0),
		BiometricTemplates: make(map[string]*EncryptedBiometric),
		EmergencyContacts: make([]*EmergencyContact, 0),
		Format:            FormatSelfContained,
		Compressed:        false,
	}
}

// NewCryptographicProof creates a new cryptographic proof
func NewCryptographicProof(proofType, proofValue, proofPurpose, verificationMethod string) *CryptographicProof {
	return &CryptographicProof{
		ProofType:          proofType,
		ProofValue:         proofValue,
		ProofPurpose:       proofPurpose,
		Created:            time.Now(),
		VerificationMethod: verificationMethod,
	}
}

// NewOfflineCredential creates a new offline credential
func NewOfflineCredential(id string, credType []string, issuer string, subject map[string]interface{}) *OfflineCredential {
	return &OfflineCredential{
		ID:                id,
		Type:              credType,
		Issuer:            issuer,
		IssuanceDate:      time.Now(),
		CredentialSubject: subject,
		OfflineUsable:     true,
		Priority:          1,
	}
}

// NewOfflineVerificationRequest creates a new offline verification request
func NewOfflineVerificationRequest(did, challenge string, requiredLevel uint32, requiredTypes []string) *OfflineVerificationRequest {
	return &OfflineVerificationRequest{
		DID:           did,
		Challenge:     challenge,
		RequiredLevel: requiredLevel,
		RequiredTypes: requiredTypes,
		Timestamp:     time.Now(),
		Nonce:         generateNonce(),
	}
}

// Validation methods

// Validate validates the offline verification data
func (ovd *OfflineVerificationData) Validate() error {
	if ovd.DID == "" {
		return fmt.Errorf("DID cannot be empty")
	}
	
	if ovd.IdentityHash == "" {
		return fmt.Errorf("identity hash cannot be empty")
	}
	
	if ovd.PublicKey == "" {
		return fmt.Errorf("public key cannot be empty")
	}
	
	if ovd.IssuedAt.IsZero() {
		return fmt.Errorf("issued at time cannot be zero")
	}
	
	if ovd.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("offline verification data has expired")
	}
	
	if ovd.VerificationLevel < 1 || ovd.VerificationLevel > 5 {
		return fmt.Errorf("verification level must be between 1 and 5")
	}
	
	// Validate cryptographic proofs
	if ovd.IdentityProof != nil {
		if err := ovd.IdentityProof.Validate(); err != nil {
			return fmt.Errorf("invalid identity proof: %w", err)
		}
	}
	
	// Validate credentials
	for i, cred := range ovd.Credentials {
		if err := cred.Validate(); err != nil {
			return fmt.Errorf("invalid credential at index %d: %w", i, err)
		}
	}
	
	// Validate revocation data
	if ovd.RevocationData != nil {
		if err := ovd.RevocationData.Validate(); err != nil {
			return fmt.Errorf("invalid revocation data: %w", err)
		}
	}
	
	return nil
}

// Validate validates a cryptographic proof
func (cp *CryptographicProof) Validate() error {
	if cp.ProofType == "" {
		return fmt.Errorf("proof type cannot be empty")
	}
	
	if cp.ProofValue == "" {
		return fmt.Errorf("proof value cannot be empty")
	}
	
	if cp.VerificationMethod == "" {
		return fmt.Errorf("verification method cannot be empty")
	}
	
	if cp.Created.IsZero() {
		return fmt.Errorf("created time cannot be zero")
	}
	
	return nil
}

// Validate validates an offline credential
func (oc *OfflineCredential) Validate() error {
	if oc.ID == "" {
		return fmt.Errorf("credential ID cannot be empty")
	}
	
	if len(oc.Type) == 0 {
		return fmt.Errorf("credential type cannot be empty")
	}
	
	if oc.Issuer == "" {
		return fmt.Errorf("issuer cannot be empty")
	}
	
	if oc.IssuanceDate.IsZero() {
		return fmt.Errorf("issuance date cannot be zero")
	}
	
	if oc.ExpirationDate != nil && oc.ExpirationDate.Before(time.Now()) {
		return fmt.Errorf("credential has expired")
	}
	
	if oc.CredentialSubject == nil {
		return fmt.Errorf("credential subject cannot be nil")
	}
	
	if oc.Proof != nil {
		if err := oc.Proof.Validate(); err != nil {
			return fmt.Errorf("invalid proof: %w", err)
		}
	}
	
	return nil
}

// Validate validates revocation data
func (rd *RevocationData) Validate() error {
	if rd.RevocationListURL == "" {
		return fmt.Errorf("revocation list URL cannot be empty")
	}
	
	if rd.RevocationListHash == "" {
		return fmt.Errorf("revocation list hash cannot be empty")
	}
	
	if rd.LastUpdated.IsZero() {
		return fmt.Errorf("last updated time cannot be zero")
	}
	
	if rd.ValidUntil.Before(time.Now()) {
		return fmt.Errorf("revocation data has expired")
	}
	
	return nil
}

// Utility methods

// ComputeDataHash computes the hash of the offline verification data
func (ovd *OfflineVerificationData) ComputeDataHash() (string, error) {
	// Create a copy without the hash and signature fields
	copy := *ovd
	copy.DataHash = ""
	copy.Signature = ""
	
	data, err := json.Marshal(copy)
	if err != nil {
		return "", fmt.Errorf("failed to marshal data: %w", err)
	}
	
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash), nil
}

// IsExpired checks if the offline verification data has expired
func (ovd *OfflineVerificationData) IsExpired() bool {
	return time.Now().After(ovd.ExpiresAt)
}

// GetCredentialByType returns credentials of a specific type
func (ovd *OfflineVerificationData) GetCredentialByType(credType string) []*OfflineCredential {
	var result []*OfflineCredential
	
	for _, cred := range ovd.Credentials {
		for _, t := range cred.Type {
			if t == credType {
				result = append(result, cred)
				break
			}
		}
	}
	
	return result
}

// IsRevoked checks if a credential or identity is revoked
func (rd *RevocationData) IsRevoked(id string) bool {
	for _, revokedID := range rd.RevokedCredentials {
		if revokedID == id {
			return true
		}
	}
	
	for _, revokedID := range rd.RevokedIdentities {
		if revokedID == id {
			return true
		}
	}
	
	return false
}

// generateNonce generates a random nonce for verification requests
func generateNonce() string {
	// This is a simplified implementation
	// In production, use cryptographically secure random generation
	hash := sha256.Sum256([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))
	return fmt.Sprintf("%x", hash)[:16]
}

// DefaultOfflineVerificationConfig returns default configuration for offline verification
func DefaultOfflineVerificationConfig() *OfflineVerificationConfig {
	return &OfflineVerificationConfig{
		MaxOfflineDuration:    24 * time.Hour,
		RequiredConfidence:    0.85,
		BiometricThreshold:    0.95,
		MaxCacheSize:          100 * 1024 * 1024, // 100MB
		CacheExpirationPeriod: 7 * 24 * time.Hour, // 7 days
		EnableCompression:     true,
		CompressionLevel:      6,
		EmergencyModeEnabled:  true,
		EmergencyThreshold:    0.70,
		SupportedRegions:      []string{"india", "global"},
		DefaultLanguage:       LanguageEnglish,
	}
}

// Codec registration

// RegisterCodec registers the offline verification types with the codec
func RegisterOfflineCodec(cdc codec.LegacyAmino) {
	cdc.RegisterConcrete(&OfflineVerificationData{}, "identity/OfflineVerificationData", nil)
	cdc.RegisterConcrete(&CryptographicProof{}, "identity/CryptographicProof", nil)
	cdc.RegisterConcrete(&OfflineCredential{}, "identity/OfflineCredential", nil)
	cdc.RegisterConcrete(&RevocationData{}, "identity/RevocationData", nil)
	cdc.RegisterConcrete(&EncryptedBiometric{}, "identity/EncryptedBiometric", nil)
	cdc.RegisterConcrete(&OfflineVerificationRequest{}, "identity/OfflineVerificationRequest", nil)
	cdc.RegisterConcrete(&OfflineVerificationResult{}, "identity/OfflineVerificationResult", nil)
	cdc.RegisterConcrete(&OfflineVerificationConfig{}, "identity/OfflineVerificationConfig", nil)
}