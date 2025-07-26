package types

import (
	"encoding/json"
	"fmt"
	"time"
)

// DIDMethod represents the DID method for DeshChain
const DIDMethod = "desh"

// DIDDocument represents a W3C compliant DID Document
type DIDDocument struct {
	Context            []string                  `json:"@context"`
	ID                 string                    `json:"id"`
	Controller         string                    `json:"controller,omitempty"`
	VerificationMethod []VerificationMethod      `json:"verificationMethod"`
	Authentication     []interface{}             `json:"authentication,omitempty"`
	AssertionMethod    []interface{}             `json:"assertionMethod,omitempty"`
	KeyAgreement       []interface{}             `json:"keyAgreement,omitempty"`
	Service            []Service                 `json:"service,omitempty"`
	Created            time.Time                 `json:"created"`
	Updated            time.Time                 `json:"updated"`
	Proof              *Proof                    `json:"proof,omitempty"`
	Metadata           map[string]interface{}    `json:"metadata,omitempty"`
}

// VerificationMethod represents a cryptographic public key
type VerificationMethod struct {
	ID                 string    `json:"id"`
	Type               string    `json:"type"`
	Controller         string    `json:"controller"`
	PublicKeyMultibase string    `json:"publicKeyMultibase,omitempty"`
	PublicKeyJwk       *JWK      `json:"publicKeyJwk,omitempty"`
	BlockchainAccount  string    `json:"blockchainAccountId,omitempty"`
	Created            time.Time `json:"created"`
	Revoked            bool      `json:"revoked,omitempty"`
	RevokedAt          time.Time `json:"revokedAt,omitempty"`
}

// JWK represents a JSON Web Key
type JWK struct {
	Kty string `json:"kty"`
	Crv string `json:"crv,omitempty"`
	X   string `json:"x,omitempty"`
	Y   string `json:"y,omitempty"`
	N   string `json:"n,omitempty"`
	E   string `json:"e,omitempty"`
}

// Service represents a service endpoint in the DID Document
type Service struct {
	ID              string                 `json:"id"`
	Type            string                 `json:"type"`
	ServiceEndpoint interface{}            `json:"serviceEndpoint"`
	Privacy         string                 `json:"privacy,omitempty"` // public, private, permissioned
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// Proof represents a cryptographic proof
type Proof struct {
	Type               string    `json:"type"`
	Created            time.Time `json:"created"`
	ProofPurpose       string    `json:"proofPurpose"`
	VerificationMethod string    `json:"verificationMethod"`
	ProofValue         string    `json:"proofValue"`
	Challenge          string    `json:"challenge,omitempty"`
	Domain             string    `json:"domain,omitempty"`
	Nonce              string    `json:"nonce,omitempty"`
}

// DIDResolutionMetadata contains metadata about the DID resolution process
type DIDResolutionMetadata struct {
	ContentType     string    `json:"contentType"`
	ResolutionTime  time.Time `json:"resolutionTime"`
	Error           string    `json:"error,omitempty"`
	ErrorMessage    string    `json:"errorMessage,omitempty"`
	CanonicalID     string    `json:"canonicalId,omitempty"`
	EquivalentID    []string  `json:"equivalentId,omitempty"`
}

// DIDDocumentMetadata contains metadata about the DID Document
type DIDDocumentMetadata struct {
	Created         time.Time `json:"created"`
	Updated         time.Time `json:"updated"`
	Deactivated     bool      `json:"deactivated"`
	NextUpdate      time.Time `json:"nextUpdate,omitempty"`
	VersionID       string    `json:"versionId"`
	NextVersionID   string    `json:"nextVersionId,omitempty"`
	EquivalentID    []string  `json:"equivalentId,omitempty"`
	CanonicalID     string    `json:"canonicalId,omitempty"`
}

// DIDResolutionResult represents the result of DID resolution
type DIDResolutionResult struct {
	DidDocument           *DIDDocument           `json:"didDocument"`
	DidResolutionMetadata *DIDResolutionMetadata `json:"didResolutionMetadata"`
	DidDocumentMetadata   *DIDDocumentMetadata   `json:"didDocumentMetadata"`
}

// DIDState represents the state of a DID
type DIDState int32

const (
	DIDState_ACTIVE      DIDState = 0
	DIDState_DEACTIVATED DIDState = 1
	DIDState_REVOKED     DIDState = 2
)

// String returns the string representation of DIDState
func (ds DIDState) String() string {
	switch ds {
	case DIDState_ACTIVE:
		return "ACTIVE"
	case DIDState_DEACTIVATED:
		return "DEACTIVATED"
	case DIDState_REVOKED:
		return "REVOKED"
	default:
		return "UNKNOWN"
	}
}

// GenerateDID creates a new DID for DeshChain
func GenerateDID(identifier string) string {
	return fmt.Sprintf("did:%s:%s", DIDMethod, identifier)
}

// ParseDID parses a DID string and returns its components
func ParseDID(did string) (method, identifier string, err error) {
	var parts []string
	fmt.Sscanf(did, "did:%s:%s", &method, &identifier)
	if method == "" || identifier == "" {
		return "", "", fmt.Errorf("invalid DID format: %s", did)
	}
	return method, identifier, nil
}

// ValidateDIDDocument performs basic validation on a DID Document
func ValidateDIDDocument(doc *DIDDocument) error {
	if doc.ID == "" {
		return fmt.Errorf("DID Document must have an ID")
	}
	
	if len(doc.Context) == 0 {
		return fmt.Errorf("DID Document must have at least one context")
	}
	
	if len(doc.VerificationMethod) == 0 {
		return fmt.Errorf("DID Document must have at least one verification method")
	}
	
	// Validate DID format
	method, _, err := ParseDID(doc.ID)
	if err != nil {
		return err
	}
	
	if method != DIDMethod {
		return fmt.Errorf("invalid DID method: expected %s, got %s", DIDMethod, method)
	}
	
	return nil
}

// ToJSON converts DID Document to JSON
func (d *DIDDocument) ToJSON() ([]byte, error) {
	return json.MarshalIndent(d, "", "  ")
}

// FromJSON creates DID Document from JSON
func (d *DIDDocument) FromJSON(data []byte) error {
	return json.Unmarshal(data, d)
}

// AddVerificationMethod adds a new verification method to the DID Document
func (d *DIDDocument) AddVerificationMethod(vm VerificationMethod) {
	d.VerificationMethod = append(d.VerificationMethod, vm)
	d.Updated = time.Now()
}

// AddService adds a new service endpoint to the DID Document
func (d *DIDDocument) AddService(svc Service) {
	d.Service = append(d.Service, svc)
	d.Updated = time.Now()
}

// GetVerificationMethod returns a verification method by ID
func (d *DIDDocument) GetVerificationMethod(id string) (*VerificationMethod, error) {
	for _, vm := range d.VerificationMethod {
		if vm.ID == id {
			return &vm, nil
		}
	}
	return nil, fmt.Errorf("verification method not found: %s", id)
}

// GetService returns a service by ID
func (d *DIDDocument) GetService(id string) (*Service, error) {
	for _, svc := range d.Service {
		if svc.ID == id {
			return &svc, nil
		}
	}
	return nil, fmt.Errorf("service not found: %s", id)
}

// IsActive checks if the DID Document is active
func (d *DIDDocument) IsActive() bool {
	if d.Metadata != nil {
		if deactivated, ok := d.Metadata["deactivated"].(bool); ok && deactivated {
			return false
		}
	}
	return true
}