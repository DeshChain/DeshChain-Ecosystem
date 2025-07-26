package types

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

// VerifiableCredential represents a W3C compliant Verifiable Credential
type VerifiableCredential struct {
	Context           []string               `json:"@context"`
	ID                string                 `json:"id"`
	Type              []string               `json:"type"`
	Issuer            string                 `json:"issuer"`
	IssuanceDate      time.Time              `json:"issuanceDate"`
	ExpirationDate    *time.Time             `json:"expirationDate,omitempty"`
	CredentialSubject interface{}            `json:"credentialSubject"`
	CredentialStatus  *CredentialStatus      `json:"credentialStatus,omitempty"`
	CredentialSchema  *CredentialSchema      `json:"credentialSchema,omitempty"`
	Evidence          []Evidence             `json:"evidence,omitempty"`
	TermsOfUse        []TermsOfUse           `json:"termsOfUse,omitempty"`
	RefreshService    *RefreshService        `json:"refreshService,omitempty"`
	Proof             interface{}            `json:"proof,omitempty"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
}

// CredentialStatus represents the revocation status of a credential
type CredentialStatus struct {
	ID                   string `json:"id"`
	Type                 string `json:"type"`
	RevocationListIndex  string `json:"revocationListIndex,omitempty"`
	RevocationListCredential string `json:"revocationListCredential,omitempty"`
	StatusPurpose        string `json:"statusPurpose,omitempty"`
}

// CredentialSchema defines the schema for credential validation
type CredentialSchema struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// Evidence represents evidence supporting the credential claims
type Evidence struct {
	ID               string                 `json:"id,omitempty"`
	Type             []string               `json:"type"`
	Verifier         string                 `json:"verifier"`
	EvidenceDocument string                 `json:"evidenceDocument"`
	SubjectPresence  string                 `json:"subjectPresence"`
	DocumentPresence string                 `json:"documentPresence"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// TermsOfUse represents the terms under which the credential can be used
type TermsOfUse struct {
	ID         string    `json:"id,omitempty"`
	Type       string    `json:"type"`
	Profile    string    `json:"profile,omitempty"`
	Obligation []string  `json:"obligation,omitempty"`
	Assigner   string    `json:"assigner,omitempty"`
	Assignee   string    `json:"assignee,omitempty"`
	ValidFrom  time.Time `json:"validFrom,omitempty"`
	ValidUntil time.Time `json:"validUntil,omitempty"`
}

// RefreshService defines how to refresh the credential
type RefreshService struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// VerifiablePresentation represents a collection of credentials presented by a holder
type VerifiablePresentation struct {
	Context              []string               `json:"@context"`
	ID                   string                 `json:"id"`
	Type                 []string               `json:"type"`
	Holder               string                 `json:"holder"`
	VerifiableCredential []interface{}          `json:"verifiableCredential"`
	Proof                interface{}            `json:"proof,omitempty"`
	Challenge            string                 `json:"challenge,omitempty"`
	Domain               string                 `json:"domain,omitempty"`
	Metadata             map[string]interface{} `json:"metadata,omitempty"`
}

// CredentialSubject represents the subject of a credential
type CredentialSubject struct {
	ID       string                 `json:"id"`
	Type     []string               `json:"type,omitempty"`
	Claims   map[string]interface{} `json:"claims"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// KYCCredentialSubject represents KYC-specific credential subject
type KYCCredentialSubject struct {
	ID              string    `json:"id"`
	Type            []string  `json:"type"`
	GivenName       string    `json:"givenName"`
	FamilyName      string    `json:"familyName"`
	DateOfBirth     string    `json:"dateOfBirth"`
	Nationality     string    `json:"nationality"`
	DocumentType    string    `json:"documentType"`
	DocumentNumber  string    `json:"documentNumber"`
	DocumentCountry string    `json:"documentCountry"`
	Address         Address   `json:"address"`
	KYCLevel        string    `json:"kycLevel"`
	RiskScore       float64   `json:"riskScore"`
	VerifiedAt      time.Time `json:"verifiedAt"`
	AMLStatus       string    `json:"amlStatus,omitempty"`
}

// Address represents a physical address
type Address struct {
	StreetAddress   string `json:"streetAddress"`
	AddressLocality string `json:"addressLocality"`
	AddressRegion   string `json:"addressRegion"`
	PostalCode      string `json:"postalCode"`
	AddressCountry  string `json:"addressCountry"`
}

// BiometricCredentialSubject represents biometric credential subject
type BiometricCredentialSubject struct {
	ID               string                 `json:"id"`
	Type             []string               `json:"type"`
	BiometricType    string                 `json:"biometricType"`
	TemplateFormat   string                 `json:"templateFormat"`
	TemplateHash     string                 `json:"templateHash"`
	QualityScore     float64                `json:"qualityScore"`
	LivenessScore    float64                `json:"livenessScore"`
	DeviceInfo       map[string]interface{} `json:"deviceInfo"`
	EnrollmentDate   time.Time              `json:"enrollmentDate"`
	ConsentGiven     bool                   `json:"consentGiven"`
	ConsentTimestamp time.Time              `json:"consentTimestamp"`
}

// CredentialRequest represents a request for credential issuance
type CredentialRequest struct {
	ID                string                 `json:"id"`
	Type              []string               `json:"type"`
	Holder            string                 `json:"holder"`
	CredentialSchema  *CredentialSchema      `json:"credentialSchema"`
	RequestedClaims   []string               `json:"requestedClaims"`
	Purpose           string                 `json:"purpose"`
	Challenge         string                 `json:"challenge"`
	Domain            string                 `json:"domain"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
}

// SelectiveDisclosure allows revealing only specific claims
type SelectiveDisclosure struct {
	CredentialID    string   `json:"credentialId"`
	RevealedClaims  []string `json:"revealedClaims"`
	BlindedClaims   []string `json:"blindedClaims"`
	ProofType       string   `json:"proofType"`
	ProofValue      string   `json:"proofValue"`
}

// ValidateVerifiableCredential performs basic validation on a VC
func ValidateVerifiableCredential(vc *VerifiableCredential) error {
	if vc.ID == "" {
		return fmt.Errorf("verifiable credential must have an ID")
	}
	
	if len(vc.Context) == 0 {
		return fmt.Errorf("verifiable credential must have at least one context")
	}
	
	if len(vc.Type) == 0 {
		return fmt.Errorf("verifiable credential must have at least one type")
	}
	
	if vc.Issuer == "" {
		return fmt.Errorf("verifiable credential must have an issuer")
	}
	
	if vc.CredentialSubject == nil {
		return fmt.Errorf("verifiable credential must have a credential subject")
	}
	
	// Check if credential is expired
	if vc.ExpirationDate != nil && vc.ExpirationDate.Before(time.Now()) {
		return fmt.Errorf("verifiable credential has expired")
	}
	
	return nil
}

// ToJSON converts VC to JSON
func (vc *VerifiableCredential) ToJSON() ([]byte, error) {
	return json.MarshalIndent(vc, "", "  ")
}

// FromJSON creates VC from JSON
func (vc *VerifiableCredential) FromJSON(data []byte) error {
	return json.Unmarshal(data, vc)
}

// GetHash returns the hash of the credential
func (vc *VerifiableCredential) GetHash() (string, error) {
	// Remove proof before hashing
	vcCopy := *vc
	vcCopy.Proof = nil
	
	data, err := json.Marshal(vcCopy)
	if err != nil {
		return "", err
	}
	
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}

// IsExpired checks if the credential is expired
func (vc *VerifiableCredential) IsExpired() bool {
	if vc.ExpirationDate == nil {
		return false
	}
	return vc.ExpirationDate.Before(time.Now())
}

// HasType checks if the credential has a specific type
func (vc *VerifiableCredential) HasType(credType string) bool {
	for _, t := range vc.Type {
		if t == credType {
			return true
		}
	}
	return false
}

// GetSubjectID returns the ID of the credential subject
func (vc *VerifiableCredential) GetSubjectID() (string, error) {
	switch subject := vc.CredentialSubject.(type) {
	case map[string]interface{}:
		if id, ok := subject["id"].(string); ok {
			return id, nil
		}
	case CredentialSubject:
		return subject.ID, nil
	case KYCCredentialSubject:
		return subject.ID, nil
	case BiometricCredentialSubject:
		return subject.ID, nil
	}
	return "", fmt.Errorf("unable to extract subject ID from credential")
}

// ValidateVerifiablePresentation performs basic validation on a VP
func ValidateVerifiablePresentation(vp *VerifiablePresentation) error {
	if vp.ID == "" {
		return fmt.Errorf("verifiable presentation must have an ID")
	}
	
	if len(vp.Context) == 0 {
		return fmt.Errorf("verifiable presentation must have at least one context")
	}
	
	if len(vp.Type) == 0 {
		return fmt.Errorf("verifiable presentation must have at least one type")
	}
	
	if vp.Holder == "" {
		return fmt.Errorf("verifiable presentation must have a holder")
	}
	
	if len(vp.VerifiableCredential) == 0 {
		return fmt.Errorf("verifiable presentation must have at least one credential")
	}
	
	return nil
}

// Credential type constants
const (
	CredentialTypeVerifiable = "VerifiableCredential"
	CredentialTypeKYC        = "KYCCredential"
	CredentialTypeBiometric  = "BiometricCredential"
	CredentialTypeAddress    = "AddressCredential"
	CredentialTypeIdentity   = "IdentityCredential"
	CredentialTypeAge        = "AgeCredential"
	CredentialTypeIncome     = "IncomeCredential"
	CredentialTypeEducation  = "EducationCredential"
)

// Standard contexts
var (
	ContextW3CCredentials = "https://www.w3.org/2018/credentials/v1"
	ContextDeshChain      = "https://deshchain.com/contexts/identity/v1"
	ContextIndiaStack     = "https://indiastack.org/contexts/v1"
)