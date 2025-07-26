package types

import (
	"time"
)

// Identity Federation Types for external system integration

// FederatedIdentityProvider represents an external identity provider
type FederatedIdentityProvider struct {
	ProviderID          string                 `json:"provider_id"`
	Name                string                 `json:"name"`
	Type                FederationProviderType `json:"type"`
	Status              ProviderStatus         `json:"status"`
	Configuration       ProviderConfiguration  `json:"configuration"`
	TrustLevel          TrustLevel             `json:"trust_level"`
	SupportedProtocols  []FederationProtocol   `json:"supported_protocols"`
	CredentialMapping   CredentialMappingRules `json:"credential_mapping"`
	SecuritySettings    SecuritySettings       `json:"security_settings"`
	ComplianceSettings  ComplianceSettings     `json:"compliance_settings"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
	LastSyncAt         *time.Time             `json:"last_sync_at,omitempty"`
	SyncStatus         SyncStatus             `json:"sync_status"`
	Metadata           map[string]interface{} `json:"metadata,omitempty"`
}

// FederationProviderType categorizes different types of identity providers
type FederationProviderType int32

const (
	FederationProviderType_GOVERNMENT_ID     FederationProviderType = 0 // Aadhaar, DigiLocker, etc.
	FederationProviderType_FINANCIAL         FederationProviderType = 1 // Banks, UPI providers
	FederationProviderType_TELECOM           FederationProviderType = 2 // Mobile operators
	FederationProviderType_EDUCATIONAL       FederationProviderType = 3 // Universities, certification bodies
	FederationProviderType_HEALTHCARE        FederationProviderType = 4 // Hospitals, health authorities
	FederationProviderType_ENTERPRISE        FederationProviderType = 5 // Corporate identity systems
	FederationProviderType_BLOCKCHAIN        FederationProviderType = 6 // Other blockchain networks
	FederationProviderType_SOCIAL            FederationProviderType = 7 // Social media platforms
	FederationProviderType_BIOMETRIC         FederationProviderType = 8 // Biometric service providers
	FederationProviderType_CUSTOM            FederationProviderType = 9 // Custom integrations
)

// ProviderStatus indicates the current status of a federated provider
type ProviderStatus int32

const (
	ProviderStatus_ACTIVE      ProviderStatus = 0
	ProviderStatus_INACTIVE    ProviderStatus = 1
	ProviderStatus_SUSPENDED   ProviderStatus = 2
	ProviderStatus_PENDING     ProviderStatus = 3
	ProviderStatus_ERROR       ProviderStatus = 4
	ProviderStatus_MAINTENANCE ProviderStatus = 5
)

// TrustLevel indicates the level of trust for a federated provider
type TrustLevel int32

const (
	TrustLevel_LOW      TrustLevel = 0
	TrustLevel_MEDIUM   TrustLevel = 1
	TrustLevel_HIGH     TrustLevel = 2
	TrustLevel_CRITICAL TrustLevel = 3
)

// FederationProtocol represents supported federation protocols
type FederationProtocol int32

const (
	FederationProtocol_OIDC           FederationProtocol = 0 // OpenID Connect
	FederationProtocol_SAML           FederationProtocol = 1 // SAML 2.0
	FederationProtocol_OAUTH2         FederationProtocol = 2 // OAuth 2.0
	FederationProtocol_DID_COMM       FederationProtocol = 3 // DID Communication
	FederationProtocol_HYPERLEDGER    FederationProtocol = 4 // Hyperledger Aries
	FederationProtocol_CUSTOM_API     FederationProtocol = 5 // Custom REST API
	FederationProtocol_BLOCKCHAIN_RPC FederationProtocol = 6 // Blockchain RPC
	FederationProtocol_INDIA_STACK    FederationProtocol = 7 // India Stack APIs
)

// SyncStatus tracks synchronization status with external providers
type SyncStatus int32

const (
	SyncStatus_SYNCED     SyncStatus = 0
	SyncStatus_PENDING    SyncStatus = 1
	SyncStatus_FAILED     SyncStatus = 2
	SyncStatus_PARTIAL    SyncStatus = 3
	SyncStatus_DISABLED   SyncStatus = 4
)

// ProviderConfiguration contains provider-specific configuration
type ProviderConfiguration struct {
	Endpoint            string                 `json:"endpoint"`
	ClientID            string                 `json:"client_id"`
	ClientSecret        string                 `json:"client_secret,omitempty"` // Encrypted
	Issuer              string                 `json:"issuer,omitempty"`
	WellKnownURL        string                 `json:"well_known_url,omitempty"`
	CertificateURL      string                 `json:"certificate_url,omitempty"`
	TokenEndpoint       string                 `json:"token_endpoint,omitempty"`
	UserInfoEndpoint    string                 `json:"userinfo_endpoint,omitempty"`
	AuthorizationURL    string                 `json:"authorization_url,omitempty"`
	RevocationEndpoint  string                 `json:"revocation_endpoint,omitempty"`
	Scopes              []string               `json:"scopes,omitempty"`
	RedirectURIs        []string               `json:"redirect_uris,omitempty"`
	CustomHeaders       map[string]string      `json:"custom_headers,omitempty"`
	RateLimits          RateLimitConfig        `json:"rate_limits"`
	Timeout             time.Duration          `json:"timeout"`
	RetryPolicy         RetryPolicy            `json:"retry_policy"`
	CustomParameters    map[string]interface{} `json:"custom_parameters,omitempty"`
}

// RateLimitConfig defines rate limiting for external provider calls
type RateLimitConfig struct {
	RequestsPerSecond int64         `json:"requests_per_second"`
	BurstSize         int64         `json:"burst_size"`
	CooldownPeriod    time.Duration `json:"cooldown_period"`
}

// RetryPolicy defines retry behavior for failed requests
type RetryPolicy struct {
	MaxRetries    int           `json:"max_retries"`
	InitialDelay  time.Duration `json:"initial_delay"`
	MaxDelay      time.Duration `json:"max_delay"`
	BackoffFactor float64       `json:"backoff_factor"`
}

// CredentialMappingRules defines how external credentials map to DeshChain credentials
type CredentialMappingRules struct {
	MappingID      string                    `json:"mapping_id"`
	SourceSchema   string                    `json:"source_schema"`
	TargetSchema   string                    `json:"target_schema"`
	FieldMappings  []FieldMapping            `json:"field_mappings"`
	Transformations []DataTransformation     `json:"transformations"`
	ValidationRules []ValidationRule         `json:"validation_rules"`
	DefaultValues  map[string]interface{}    `json:"default_values,omitempty"`
	ConditionalLogic []ConditionalMapping    `json:"conditional_logic,omitempty"`
}

// FieldMapping maps individual fields between external and internal credentials
type FieldMapping struct {
	SourceField   string                 `json:"source_field"`
	TargetField   string                 `json:"target_field"`
	Required      bool                   `json:"required"`
	DataType      string                 `json:"data_type"`
	Format        string                 `json:"format,omitempty"`
	Transformation string                `json:"transformation,omitempty"`
	ValidationRule string                `json:"validation_rule,omitempty"`
	DefaultValue  interface{}            `json:"default_value,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// DataTransformation defines data transformation rules
type DataTransformation struct {
	TransformationID   string                 `json:"transformation_id"`
	Type               TransformationType     `json:"type"`
	SourceFields       []string               `json:"source_fields"`
	TargetField        string                 `json:"target_field"`
	Function           string                 `json:"function"`
	Parameters         map[string]interface{} `json:"parameters,omitempty"`
	ErrorHandling      ErrorHandlingPolicy    `json:"error_handling"`
}

// TransformationType categorizes different types of data transformations
type TransformationType int32

const (
	TransformationType_FORMAT_CONVERSION TransformationType = 0 // Date format, number format, etc.
	TransformationType_DATA_ENRICHMENT   TransformationType = 1 // Add computed fields
	TransformationType_DATA_NORMALIZATION TransformationType = 2 // Normalize to standard format
	TransformationType_CONCATENATION     TransformationType = 3 // Combine multiple fields
	TransformationType_EXTRACTION        TransformationType = 4 // Extract part of a field
	TransformationType_LOOKUP            TransformationType = 5 // Lookup from external source
	TransformationType_CALCULATION       TransformationType = 6 // Mathematical calculations
	TransformationType_CUSTOM            TransformationType = 7 // Custom transformation logic
)

// ErrorHandlingPolicy defines how to handle transformation errors
type ErrorHandlingPolicy int32

const (
	ErrorHandlingPolicy_FAIL      ErrorHandlingPolicy = 0 // Fail the entire operation
	ErrorHandlingPolicy_SKIP      ErrorHandlingPolicy = 1 // Skip this transformation
	ErrorHandlingPolicy_DEFAULT   ErrorHandlingPolicy = 2 // Use default value
	ErrorHandlingPolicy_LOG_WARN  ErrorHandlingPolicy = 3 // Log warning and continue
)

// ValidationRule defines validation logic for mapped data
type ValidationRule struct {
	RuleID      string                 `json:"rule_id"`
	Type        ValidationType         `json:"type"`
	Field       string                 `json:"field"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
	ErrorMessage string                `json:"error_message"`
	Severity    ValidationSeverity     `json:"severity"`
}

// ValidationType categorizes different types of validation rules
type ValidationType int32

const (
	ValidationType_REQUIRED     ValidationType = 0
	ValidationType_FORMAT       ValidationType = 1
	ValidationType_RANGE        ValidationType = 2
	ValidationType_LENGTH       ValidationType = 3
	ValidationType_PATTERN      ValidationType = 4
	ValidationType_ENUM         ValidationType = 5
	ValidationType_CROSS_FIELD  ValidationType = 6
	ValidationType_EXTERNAL     ValidationType = 7
	ValidationType_CUSTOM       ValidationType = 8
)

// ValidationSeverity indicates the severity of validation failures
type ValidationSeverity int32

const (
	ValidationSeverity_ERROR   ValidationSeverity = 0
	ValidationSeverity_WARNING ValidationSeverity = 1
	ValidationSeverity_INFO    ValidationSeverity = 2
)

// ConditionalMapping defines conditional mapping logic
type ConditionalMapping struct {
	ConditionID   string                 `json:"condition_id"`
	Condition     string                 `json:"condition"`      // Expression to evaluate
	ThenMapping   []FieldMapping         `json:"then_mapping"`   // Apply if condition is true
	ElseMapping   []FieldMapping         `json:"else_mapping"`   // Apply if condition is false
	Parameters    map[string]interface{} `json:"parameters,omitempty"`
}

// SecuritySettings defines security requirements for federation
type SecuritySettings struct {
	RequireEncryption     bool                `json:"require_encryption"`
	RequireMutualTLS      bool                `json:"require_mutual_tls"`
	RequireTokenSigning   bool                `json:"require_token_signing"`
	AllowedSigningAlgos   []string            `json:"allowed_signing_algorithms"`
	CertificateValidation CertValidationLevel `json:"certificate_validation"`
	TokenLifetime         time.Duration       `json:"token_lifetime"`
	RefreshTokenLifetime  time.Duration       `json:"refresh_token_lifetime"`
	AllowedIPRanges       []string            `json:"allowed_ip_ranges,omitempty"`
	RequireUserConsent    bool                `json:"require_user_consent"`
	MinTrustLevel         TrustLevel          `json:"min_trust_level"`
	EncryptionSettings    EncryptionSettings  `json:"encryption_settings"`
}

// CertValidationLevel defines certificate validation requirements
type CertValidationLevel int32

const (
	CertValidationLevel_NONE     CertValidationLevel = 0
	CertValidationLevel_BASIC    CertValidationLevel = 1
	CertValidationLevel_FULL     CertValidationLevel = 2
	CertValidationLevel_STRICT   CertValidationLevel = 3
)

// EncryptionSettings defines encryption requirements
type EncryptionSettings struct {
	Algorithm       string `json:"algorithm"`
	KeySize         int    `json:"key_size"`
	RequireKeyRotation bool `json:"require_key_rotation"`
	KeyRotationPeriod time.Duration `json:"key_rotation_period"`
}

// ComplianceSettings defines compliance requirements for federation
type ComplianceSettings struct {
	RequiredCompliance   []string              `json:"required_compliance"`   // GDPR, DPDP, etc.
	DataResidencyRules   []DataResidencyRule   `json:"data_residency_rules"`
	AuditRequirements    AuditRequirements     `json:"audit_requirements"`
	RetentionPolicies    []RetentionPolicy     `json:"retention_policies"`
	ConsentRequirements  ConsentRequirements   `json:"consent_requirements"`
	PrivacySettings      PrivacySettings       `json:"privacy_settings"`
}

// DataResidencyRule defines where data can be stored and processed
type DataResidencyRule struct {
	RuleID          string   `json:"rule_id"`
	DataCategory    string   `json:"data_category"`
	AllowedCountries []string `json:"allowed_countries"`
	AllowedRegions   []string `json:"allowed_regions,omitempty"`
	ProhibitedCountries []string `json:"prohibited_countries,omitempty"`
	RequireEncryption bool    `json:"require_encryption"`
}

// AuditRequirements defines audit logging requirements
type AuditRequirements struct {
	RequireAuditLog    bool          `json:"require_audit_log"`
	AuditLevel         AuditLevel    `json:"audit_level"`
	RetentionPeriod    time.Duration `json:"retention_period"`
	RequireIntegrity   bool          `json:"require_integrity"`
	RequireNonRepudiation bool       `json:"require_non_repudiation"`
}

// AuditLevel defines the level of audit logging required
type AuditLevel int32

const (
	AuditLevel_BASIC       AuditLevel = 0 // Basic operations
	AuditLevel_DETAILED    AuditLevel = 1 // Detailed operations
	AuditLevel_FULL        AuditLevel = 2 // All operations
	AuditLevel_COMPLIANCE  AuditLevel = 3 // Compliance-focused
)

// RetentionPolicy defines data retention requirements
type RetentionPolicy struct {
	PolicyID       string        `json:"policy_id"`
	DataCategory   string        `json:"data_category"`
	RetentionPeriod time.Duration `json:"retention_period"`
	AfterExpiry    ExpiryAction  `json:"after_expiry"`
	Exceptions     []string      `json:"exceptions,omitempty"`
}

// ExpiryAction defines what happens to data after retention period
type ExpiryAction int32

const (
	ExpiryAction_DELETE     ExpiryAction = 0
	ExpiryAction_ANONYMIZE  ExpiryAction = 1
	ExpiryAction_ARCHIVE    ExpiryAction = 2
	ExpiryAction_TRANSFER   ExpiryAction = 3
)

// ConsentRequirements defines consent management requirements
type ConsentRequirements struct {
	RequireExplicitConsent bool                `json:"require_explicit_consent"`
	ConsentPurposes       []string            `json:"consent_purposes"`
	ConsentLifetime       time.Duration       `json:"consent_lifetime"`
	AllowConsentWithdrawal bool               `json:"allow_consent_withdrawal"`
	ConsentGranularity    ConsentGranularity  `json:"consent_granularity"`
	ConsentRecordKeeping  bool                `json:"consent_record_keeping"`
}

// ConsentGranularity defines the granularity of consent
type ConsentGranularity int32

const (
	ConsentGranularity_ALL_OR_NOTHING ConsentGranularity = 0
	ConsentGranularity_PURPOSE_BASED  ConsentGranularity = 1
	ConsentGranularity_DATA_CATEGORY  ConsentGranularity = 2
	ConsentGranularity_FIELD_LEVEL    ConsentGranularity = 3
)

// PrivacySettings defines privacy protection settings
type PrivacySettings struct {
	MinimizeDataCollection   bool                    `json:"minimize_data_collection"`
	RequirePurposeLimitation bool                    `json:"require_purpose_limitation"`
	RequireDataMinimization  bool                    `json:"require_data_minimization"`
	RequireAccuracyControl   bool                    `json:"require_accuracy_control"`
	RequireStorageLimitation bool                    `json:"require_storage_limitation"`
	RequireIntegrityControl  bool                    `json:"require_integrity_control"`
	RequireConfidentiality   bool                    `json:"require_confidentiality"`
	AllowedProcessingBases   []ProcessingLegalBasis  `json:"allowed_processing_bases"`
	PrivacyByDesign         bool                    `json:"privacy_by_design"`
	PrivacyByDefault        bool                    `json:"privacy_by_default"`
}

// ProcessingLegalBasis defines legal bases for data processing
type ProcessingLegalBasis int32

const (
	ProcessingLegalBasis_CONSENT           ProcessingLegalBasis = 0
	ProcessingLegalBasis_CONTRACT          ProcessingLegalBasis = 1
	ProcessingLegalBasis_LEGAL_OBLIGATION  ProcessingLegalBasis = 2
	ProcessingLegalBasis_VITAL_INTERESTS   ProcessingLegalBasis = 3
	ProcessingLegalBasis_PUBLIC_TASK       ProcessingLegalBasis = 4
	ProcessingLegalBasis_LEGITIMATE_INTEREST ProcessingLegalBasis = 5
)

// FederatedCredential represents a credential from an external provider
type FederatedCredential struct {
	CredentialID        string                 `json:"credential_id"`
	ProviderID          string                 `json:"provider_id"`
	ExternalID          string                 `json:"external_id"`
	CredentialType      string                 `json:"credential_type"`
	Subject             string                 `json:"subject"`
	Issuer              string                 `json:"issuer"`
	IssuanceDate        time.Time              `json:"issuance_date"`
	ExpirationDate      *time.Time             `json:"expiration_date,omitempty"`
	RawCredential       string                 `json:"raw_credential"`
	MappedCredential    string                 `json:"mapped_credential"`
	VerificationStatus  VerificationStatus     `json:"verification_status"`
	LastVerified        *time.Time             `json:"last_verified,omitempty"`
	TrustScore          float64                `json:"trust_score"`
	ValidationResults   []ValidationResult     `json:"validation_results"`
	SyncStatus          SyncStatus             `json:"sync_status"`
	Metadata            map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt           time.Time              `json:"created_at"`
	UpdatedAt           time.Time              `json:"updated_at"`
}

// VerificationStatus indicates the verification status of a federated credential
type VerificationStatus int32

const (
	VerificationStatus_UNVERIFIED VerificationStatus = 0
	VerificationStatus_VERIFIED   VerificationStatus = 1
	VerificationStatus_EXPIRED    VerificationStatus = 2
	VerificationStatus_REVOKED    VerificationStatus = 3
	VerificationStatus_SUSPENDED  VerificationStatus = 4
	VerificationStatus_ERROR      VerificationStatus = 5
)

// ValidationResult contains the result of credential validation
type ValidationResult struct {
	RuleID      string             `json:"rule_id"`
	Passed      bool               `json:"passed"`
	Message     string             `json:"message"`
	Severity    ValidationSeverity `json:"severity"`
	Timestamp   time.Time          `json:"timestamp"`
	Details     map[string]interface{} `json:"details,omitempty"`
}

// FederationSession represents an active federation session
type FederationSession struct {
	SessionID       string                 `json:"session_id"`
	ProviderID      string                 `json:"provider_id"`
	UserDID         string                 `json:"user_did"`
	ExternalUserID  string                 `json:"external_user_id"`
	Status          SessionStatus          `json:"status"`
	Protocol        FederationProtocol     `json:"protocol"`
	AccessToken     string                 `json:"access_token,omitempty"`
	RefreshToken    string                 `json:"refresh_token,omitempty"`
	IDToken         string                 `json:"id_token,omitempty"`
	TokenType       string                 `json:"token_type"`
	ExpiresAt       time.Time              `json:"expires_at"`
	Scopes          []string               `json:"scopes"`
	Claims          map[string]interface{} `json:"claims"`
	CreatedAt       time.Time              `json:"created_at"`
	LastAccessedAt  time.Time              `json:"last_accessed_at"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// SessionStatus indicates the status of a federation session
type SessionStatus int32

const (
	SessionStatus_ACTIVE    SessionStatus = 0
	SessionStatus_EXPIRED   SessionStatus = 1
	SessionStatus_REVOKED   SessionStatus = 2
	SessionStatus_SUSPENDED SessionStatus = 3
	SessionStatus_ERROR     SessionStatus = 4
)

// FederationEvent represents events in the federation system
type FederationEvent struct {
	EventID     string                 `json:"event_id"`
	EventType   FederationEventType    `json:"event_type"`
	ProviderID  string                 `json:"provider_id"`
	UserDID     string                 `json:"user_did,omitempty"`
	SessionID   string                 `json:"session_id,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
	Status      FederationEventStatus  `json:"status"`
	Message     string                 `json:"message"`
	ErrorCode   string                 `json:"error_code,omitempty"`
	Details     map[string]interface{} `json:"details,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// FederationEventType categorizes federation events
type FederationEventType int32

const (
	FederationEventType_PROVIDER_REGISTERED   FederationEventType = 0
	FederationEventType_PROVIDER_UPDATED      FederationEventType = 1
	FederationEventType_PROVIDER_SUSPENDED    FederationEventType = 2
	FederationEventType_PROVIDER_ACTIVATED    FederationEventType = 3
	FederationEventType_CREDENTIAL_IMPORTED   FederationEventType = 4
	FederationEventType_CREDENTIAL_VERIFIED   FederationEventType = 5
	FederationEventType_CREDENTIAL_EXPIRED    FederationEventType = 6
	FederationEventType_CREDENTIAL_REVOKED    FederationEventType = 7
	FederationEventType_SESSION_CREATED       FederationEventType = 8
	FederationEventType_SESSION_RENEWED       FederationEventType = 9
	FederationEventType_SESSION_EXPIRED       FederationEventType = 10
	FederationEventType_SESSION_TERMINATED    FederationEventType = 11
	FederationEventType_SYNC_STARTED          FederationEventType = 12
	FederationEventType_SYNC_COMPLETED        FederationEventType = 13
	FederationEventType_SYNC_FAILED           FederationEventType = 14
	FederationEventType_MAPPING_UPDATED       FederationEventType = 15
	FederationEventType_TRUST_SCORE_UPDATED   FederationEventType = 16
	FederationEventType_COMPLIANCE_VIOLATION  FederationEventType = 17
	FederationEventType_SECURITY_INCIDENT     FederationEventType = 18
)

// FederationEventStatus indicates the status of federation events
type FederationEventStatus int32

const (
	FederationEventStatus_SUCCESS  FederationEventStatus = 0
	FederationEventStatus_FAILURE  FederationEventStatus = 1
	FederationEventStatus_WARNING  FederationEventStatus = 2
	FederationEventStatus_INFO     FederationEventStatus = 3
)

// FederationMetrics contains metrics about federation operations
type FederationMetrics struct {
	ProviderID              string    `json:"provider_id"`
	TotalCredentials        int64     `json:"total_credentials"`
	VerifiedCredentials     int64     `json:"verified_credentials"`
	ExpiredCredentials      int64     `json:"expired_credentials"`
	RevokedCredentials      int64     `json:"revoked_credentials"`
	ActiveSessions          int64     `json:"active_sessions"`
	TotalSessions           int64     `json:"total_sessions"`
	SuccessfulSyncs         int64     `json:"successful_syncs"`
	FailedSyncs             int64     `json:"failed_syncs"`
	AverageTrustScore       float64   `json:"average_trust_score"`
	LastSyncTimestamp       time.Time `json:"last_sync_timestamp"`
	AverageResponseTime     time.Duration `json:"average_response_time"`
	ErrorRate               float64   `json:"error_rate"`
	ComplianceScore         float64   `json:"compliance_score"`
	SecurityIncidents       int64     `json:"security_incidents"`
	DataVolumeProcessed     int64     `json:"data_volume_processed"`
	CostPerTransaction      float64   `json:"cost_per_transaction"`
}

// Helper methods

// IsActive returns true if the provider is active
func (p *FederatedIdentityProvider) IsActive() bool {
	return p.Status == ProviderStatus_ACTIVE
}

// IsHighTrust returns true if the provider has high or critical trust level
func (p *FederatedIdentityProvider) IsHighTrust() bool {
	return p.TrustLevel >= TrustLevel_HIGH
}

// SupportsProtocol returns true if the provider supports the given protocol
func (p *FederatedIdentityProvider) SupportsProtocol(protocol FederationProtocol) bool {
	for _, supported := range p.SupportedProtocols {
		if supported == protocol {
			return true
		}
	}
	return false
}

// IsValid returns true if the federated credential is valid
func (c *FederatedCredential) IsValid() bool {
	if c.VerificationStatus != VerificationStatus_VERIFIED {
		return false
	}
	
	if c.ExpirationDate != nil && time.Now().After(*c.ExpirationDate) {
		return false
	}
	
	return true
}

// IsExpired returns true if the credential has expired
func (c *FederatedCredential) IsExpired() bool {
	return c.ExpirationDate != nil && time.Now().After(*c.ExpirationDate)
}

// GetTrustScore returns the trust score of the credential
func (c *FederatedCredential) GetTrustScore() float64 {
	return c.TrustScore
}

// IsActive returns true if the session is active and not expired
func (s *FederationSession) IsActive() bool {
	return s.Status == SessionStatus_ACTIVE && time.Now().Before(s.ExpiresAt)
}

// IsExpired returns true if the session has expired
func (s *FederationSession) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// String methods for better logging

func (pt FederationProviderType) String() string {
	switch pt {
	case FederationProviderType_GOVERNMENT_ID:
		return "government_id"
	case FederationProviderType_FINANCIAL:
		return "financial"
	case FederationProviderType_TELECOM:
		return "telecom"
	case FederationProviderType_EDUCATIONAL:
		return "educational"
	case FederationProviderType_HEALTHCARE:
		return "healthcare"
	case FederationProviderType_ENTERPRISE:
		return "enterprise"
	case FederationProviderType_BLOCKCHAIN:
		return "blockchain"
	case FederationProviderType_SOCIAL:
		return "social"
	case FederationProviderType_BIOMETRIC:
		return "biometric"
	case FederationProviderType_CUSTOM:
		return "custom"
	default:
		return "unknown"
	}
}

func (fp FederationProtocol) String() string {
	switch fp {
	case FederationProtocol_OIDC:
		return "oidc"
	case FederationProtocol_SAML:
		return "saml"
	case FederationProtocol_OAUTH2:
		return "oauth2"
	case FederationProtocol_DID_COMM:
		return "did_comm"
	case FederationProtocol_HYPERLEDGER:
		return "hyperledger"
	case FederationProtocol_CUSTOM_API:
		return "custom_api"
	case FederationProtocol_BLOCKCHAIN_RPC:
		return "blockchain_rpc"
	case FederationProtocol_INDIA_STACK:
		return "india_stack"
	default:
		return "unknown"
	}
}

func (vs VerificationStatus) String() string {
	switch vs {
	case VerificationStatus_UNVERIFIED:
		return "unverified"
	case VerificationStatus_VERIFIED:
		return "verified"
	case VerificationStatus_EXPIRED:
		return "expired"
	case VerificationStatus_REVOKED:
		return "revoked"
	case VerificationStatus_SUSPENDED:
		return "suspended"
	case VerificationStatus_ERROR:
		return "error"
	default:
		return "unknown"
	}
}

// Error definitions for federation
var (
	ErrFederatedProviderNotFound     = sdkerrors.Register(ModuleName, 6001, "federated identity provider not found")
	ErrFederatedCredentialNotFound   = sdkerrors.Register(ModuleName, 6002, "federated credential not found")
	ErrFederationSessionNotFound     = sdkerrors.Register(ModuleName, 6003, "federation session not found")
	ErrInvalidFederationProtocol     = sdkerrors.Register(ModuleName, 6004, "invalid federation protocol")
	ErrProviderNotActive             = sdkerrors.Register(ModuleName, 6005, "federated provider is not active")
	ErrInsufficientTrustLevel        = sdkerrors.Register(ModuleName, 6006, "insufficient trust level for operation")
	ErrCredentialMappingFailed       = sdkerrors.Register(ModuleName, 6007, "credential mapping failed")
	ErrValidationFailed              = sdkerrors.Register(ModuleName, 6008, "credential validation failed")
	ErrSessionExpired                = sdkerrors.Register(ModuleName, 6009, "federation session has expired")
	ErrComplianceViolation           = sdkerrors.Register(ModuleName, 6010, "federation compliance violation")
	ErrSyncFailed                    = sdkerrors.Register(ModuleName, 6011, "federation synchronization failed")
	ErrRateLimitExceeded            = sdkerrors.Register(ModuleName, 6012, "federation rate limit exceeded")
	ErrEncryptionRequired           = sdkerrors.Register(ModuleName, 6013, "encryption required for federation")
	ErrCertificateValidationFailed  = sdkerrors.Register(ModuleName, 6014, "certificate validation failed")
	ErrDataResidencyViolation       = sdkerrors.Register(ModuleName, 6015, "data residency rule violation")
)