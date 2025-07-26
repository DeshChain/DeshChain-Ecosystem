package types

import (
	"fmt"
)

// GenesisState defines the identity module's genesis state
type GenesisState struct {
	Params               Params                       `json:"params"`
	Identities           []Identity                   `json:"identities"`
	DIDDocuments         []DIDDocument                `json:"did_documents"`
	Credentials          []VerifiableCredential       `json:"credentials"`
	CredentialSchemas    []CredentialSchema           `json:"credential_schemas"`
	RevocationLists      []RevocationList             `json:"revocation_lists"`
	PrivacySettings      []PrivacySettings            `json:"privacy_settings"`
	IndiaStackIntegrations []IndiaStackIntegration    `json:"india_stack_integrations"`
	ServiceRegistry      []ServiceRegistryEntry       `json:"service_registry"`
	IssuerRegistry       []IssuerRegistryEntry        `json:"issuer_registry"`
}

// DefaultGenesisState returns the default genesis state
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:               DefaultParams(),
		Identities:           []Identity{},
		DIDDocuments:         []DIDDocument{},
		Credentials:          []VerifiableCredential{},
		CredentialSchemas:    DefaultCredentialSchemas(),
		RevocationLists:      []RevocationList{},
		PrivacySettings:      []PrivacySettings{},
		IndiaStackIntegrations: []IndiaStackIntegration{},
		ServiceRegistry:      DefaultServiceRegistry(),
		IssuerRegistry:       DefaultIssuerRegistry(),
	}
}

// Validate performs basic genesis state validation
func (gs GenesisState) Validate() error {
	// Validate params
	if err := gs.Params.Validate(); err != nil {
		return fmt.Errorf("invalid params: %w", err)
	}
	
	// Validate identities
	identityMap := make(map[string]bool)
	didMap := make(map[string]bool)
	
	for _, identity := range gs.Identities {
		if err := identity.Validate(); err != nil {
			return fmt.Errorf("invalid identity %s: %w", identity.Address, err)
		}
		
		if identityMap[identity.Address] {
			return fmt.Errorf("duplicate identity address: %s", identity.Address)
		}
		identityMap[identity.Address] = true
		
		if didMap[identity.DID] {
			return fmt.Errorf("duplicate DID: %s", identity.DID)
		}
		didMap[identity.DID] = true
	}
	
	// Validate DID documents
	for _, did := range gs.DIDDocuments {
		if err := ValidateDIDDocument(&did); err != nil {
			return fmt.Errorf("invalid DID document %s: %w", did.ID, err)
		}
		
		if !didMap[did.ID] {
			return fmt.Errorf("DID document %s has no corresponding identity", did.ID)
		}
	}
	
	// Validate credentials
	credentialMap := make(map[string]bool)
	for _, cred := range gs.Credentials {
		if err := ValidateVerifiableCredential(&cred); err != nil {
			return fmt.Errorf("invalid credential %s: %w", cred.ID, err)
		}
		
		if credentialMap[cred.ID] {
			return fmt.Errorf("duplicate credential ID: %s", cred.ID)
		}
		credentialMap[cred.ID] = true
	}
	
	// Validate credential schemas
	schemaMap := make(map[string]bool)
	for _, schema := range gs.CredentialSchemas {
		if schema.ID == "" {
			return fmt.Errorf("credential schema must have an ID")
		}
		
		if schemaMap[schema.ID] {
			return fmt.Errorf("duplicate credential schema ID: %s", schema.ID)
		}
		schemaMap[schema.ID] = true
	}
	
	// Validate privacy settings
	privacyMap := make(map[string]bool)
	for _, settings := range gs.PrivacySettings {
		if settings.UserAddress == "" {
			return fmt.Errorf("privacy settings must have a user address")
		}
		
		if privacyMap[settings.UserAddress] {
			return fmt.Errorf("duplicate privacy settings for address: %s", settings.UserAddress)
		}
		privacyMap[settings.UserAddress] = true
		
		if !identityMap[settings.UserAddress] {
			return fmt.Errorf("privacy settings for non-existent identity: %s", settings.UserAddress)
		}
	}
	
	// Validate India Stack integrations
	indiaStackMap := make(map[string]bool)
	for _, integration := range gs.IndiaStackIntegrations {
		if integration.UserAddress == "" {
			return fmt.Errorf("India Stack integration must have a user address")
		}
		
		if indiaStackMap[integration.UserAddress] {
			return fmt.Errorf("duplicate India Stack integration for address: %s", integration.UserAddress)
		}
		indiaStackMap[integration.UserAddress] = true
		
		if !identityMap[integration.UserAddress] {
			return fmt.Errorf("India Stack integration for non-existent identity: %s", integration.UserAddress)
		}
	}
	
	// Validate service registry
	serviceMap := make(map[string]bool)
	for _, service := range gs.ServiceRegistry {
		if err := service.Validate(); err != nil {
			return fmt.Errorf("invalid service registry entry %s: %w", service.ID, err)
		}
		
		if serviceMap[service.ID] {
			return fmt.Errorf("duplicate service ID: %s", service.ID)
		}
		serviceMap[service.ID] = true
	}
	
	// Validate issuer registry
	issuerMap := make(map[string]bool)
	for _, issuer := range gs.IssuerRegistry {
		if err := issuer.Validate(); err != nil {
			return fmt.Errorf("invalid issuer registry entry %s: %w", issuer.DID, err)
		}
		
		if issuerMap[issuer.DID] {
			return fmt.Errorf("duplicate issuer DID: %s", issuer.DID)
		}
		issuerMap[issuer.DID] = true
	}
	
	return nil
}

// RevocationList represents a credential revocation list
type RevocationList struct {
	ID                string   `json:"id"`
	Issuer            string   `json:"issuer"`
	RevokedCredentials []string `json:"revoked_credentials"`
	LastUpdated       string   `json:"last_updated"`
}

// ServiceRegistryEntry represents a registered service
type ServiceRegistryEntry struct {
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	Type              string            `json:"type"`
	DID               string            `json:"did"`
	Endpoint          string            `json:"endpoint"`
	IsActive          bool              `json:"is_active"`
	AllowedOperations []string          `json:"allowed_operations"`
	Metadata          map[string]string `json:"metadata,omitempty"`
}

// Validate validates a service registry entry
func (s *ServiceRegistryEntry) Validate() error {
	if s.ID == "" {
		return fmt.Errorf("service must have an ID")
	}
	if s.Name == "" {
		return fmt.Errorf("service must have a name")
	}
	if s.Type == "" {
		return fmt.Errorf("service must have a type")
	}
	if s.DID == "" {
		return fmt.Errorf("service must have a DID")
	}
	if s.Endpoint == "" {
		return fmt.Errorf("service must have an endpoint")
	}
	return nil
}

// IssuerRegistryEntry represents a registered credential issuer
type IssuerRegistryEntry struct {
	DID                  string            `json:"did"`
	Name                 string            `json:"name"`
	Type                 string            `json:"type"`
	IsActive             bool              `json:"is_active"`
	AllowedCredentialTypes []string        `json:"allowed_credential_types"`
	TrustLevel           int32             `json:"trust_level"`
	Metadata             map[string]string `json:"metadata,omitempty"`
}

// Validate validates an issuer registry entry
func (i *IssuerRegistryEntry) Validate() error {
	if i.DID == "" {
		return fmt.Errorf("issuer must have a DID")
	}
	if i.Name == "" {
		return fmt.Errorf("issuer must have a name")
	}
	if i.Type == "" {
		return fmt.Errorf("issuer must have a type")
	}
	if i.TrustLevel < 0 || i.TrustLevel > 100 {
		return fmt.Errorf("issuer trust level must be between 0 and 100")
	}
	return nil
}

// DefaultCredentialSchemas returns default credential schemas
func DefaultCredentialSchemas() []CredentialSchema {
	return []CredentialSchema{
		{
			ID:   "https://deshchain.com/schemas/identity/v1/KYCCredential",
			Type: "JsonSchemaValidator2018",
		},
		{
			ID:   "https://deshchain.com/schemas/identity/v1/BiometricCredential",
			Type: "JsonSchemaValidator2018",
		},
		{
			ID:   "https://deshchain.com/schemas/identity/v1/AddressCredential",
			Type: "JsonSchemaValidator2018",
		},
		{
			ID:   "https://deshchain.com/schemas/identity/v1/AgeCredential",
			Type: "JsonSchemaValidator2018",
		},
		{
			ID:   "https://deshchain.com/schemas/identity/v1/IncomeCredential",
			Type: "JsonSchemaValidator2018",
		},
	}
}

// DefaultServiceRegistry returns default service registry entries
func DefaultServiceRegistry() []ServiceRegistryEntry {
	return []ServiceRegistryEntry{
		{
			ID:       "kyc-service-1",
			Name:     "DeshChain KYC Service",
			Type:     "KYCVerification",
			DID:      "did:desh:kyc-service",
			Endpoint: "https://kyc.deshchain.com/api/v1",
			IsActive: true,
			AllowedOperations: []string{
				"verify_identity",
				"issue_kyc_credential",
				"check_aml_status",
			},
		},
		{
			ID:       "biometric-service-1",
			Name:     "DeshChain Biometric Service",
			Type:     "BiometricEnrollment",
			DID:      "did:desh:biometric-service",
			Endpoint: "https://biometric.deshchain.com/api/v1",
			IsActive: true,
			AllowedOperations: []string{
				"enroll_biometric",
				"verify_biometric",
				"issue_biometric_credential",
			},
		},
	}
}

// DefaultIssuerRegistry returns default issuer registry entries
func DefaultIssuerRegistry() []IssuerRegistryEntry {
	return []IssuerRegistryEntry{
		{
			DID:      "did:desh:kyc-issuer",
			Name:     "DeshChain KYC Issuer",
			Type:     "KYCIssuer",
			IsActive: true,
			AllowedCredentialTypes: []string{
				"KYCCredential",
				"AMLCredential",
				"AddressCredential",
			},
			TrustLevel: 90,
		},
		{
			DID:      "did:desh:biometric-issuer",
			Name:     "DeshChain Biometric Issuer",
			Type:     "BiometricIssuer",
			IsActive: true,
			AllowedCredentialTypes: []string{
				"BiometricCredential",
			},
			TrustLevel: 85,
		},
		{
			DID:      "did:desh:government-issuer",
			Name:     "Government Document Issuer",
			Type:     "GovernmentIssuer",
			IsActive: true,
			AllowedCredentialTypes: []string{
				"AadhaarCredential",
				"PANCredential",
				"VoterIDCredential",
			},
			TrustLevel: 100,
		},
	}
}