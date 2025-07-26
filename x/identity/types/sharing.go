package types

import (
	"time"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Cross-Module Identity Sharing Protocol Types

// IdentityShareRequest represents a request to share identity data between modules
type IdentityShareRequest struct {
	RequestID        string                 `json:"request_id"`
	RequesterModule  string                 `json:"requester_module"`
	ProviderModule   string                 `json:"provider_module"`
	HolderDID        string                 `json:"holder_did"`
	RequestedData    []DataRequest          `json:"requested_data"`
	Purpose          string                 `json:"purpose"`
	Justification    string                 `json:"justification"`
	TTL              time.Duration          `json:"ttl"`
	RequestedAt      time.Time              `json:"requested_at"`
	Status           ShareRequestStatus     `json:"status"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// DataRequest specifies what data is being requested
type DataRequest struct {
	CredentialType   string   `json:"credential_type"`
	Attributes       []string `json:"attributes"`
	MinimumTrustLevel string  `json:"minimum_trust_level"`
	Required         bool     `json:"required"`
}

// ShareRequestStatus represents the status of a sharing request
type ShareRequestStatus int32

const (
	ShareRequestStatus_PENDING  ShareRequestStatus = 0
	ShareRequestStatus_APPROVED ShareRequestStatus = 1
	ShareRequestStatus_DENIED   ShareRequestStatus = 2
	ShareRequestStatus_EXPIRED  ShareRequestStatus = 3
	ShareRequestStatus_REVOKED  ShareRequestStatus = 4
)

// IdentityShareResponse represents the response to an identity sharing request
type IdentityShareResponse struct {
	RequestID       string                 `json:"request_id"`
	Status          ShareRequestStatus     `json:"status"`
	SharedData      []SharedCredential     `json:"shared_data,omitempty"`
	DenialReason    string                 `json:"denial_reason,omitempty"`
	ExpiresAt       time.Time              `json:"expires_at"`
	ResponseAt      time.Time              `json:"response_at"`
	AccessToken     string                 `json:"access_token,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// SharedCredential represents credential data shared between modules
type SharedCredential struct {
	CredentialID    string                 `json:"credential_id"`
	CredentialType  string                 `json:"credential_type"`
	Issuer          string                 `json:"issuer"`
	IssuedAt        time.Time              `json:"issued_at"`
	SharedData      map[string]interface{} `json:"shared_data"`
	TrustLevel      string                 `json:"trust_level"`
	VerificationProof string               `json:"verification_proof,omitempty"`
}

// IdentityShareAgreement represents a standing agreement between modules
type IdentityShareAgreement struct {
	AgreementID      string                 `json:"agreement_id"`
	RequesterModule  string                 `json:"requester_module"`
	ProviderModule   string                 `json:"provider_module"`
	AllowedDataTypes []string               `json:"allowed_data_types"`
	Purposes         []string               `json:"purposes"`
	TrustLevel       string                 `json:"trust_level"`
	AutoApprove      bool                   `json:"auto_approve"`
	MaxTTL           time.Duration          `json:"max_ttl"`
	CreatedAt        time.Time              `json:"created_at"`
	ExpiresAt        time.Time              `json:"expires_at"`
	Status           AgreementStatus        `json:"status"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// AgreementStatus represents the status of a sharing agreement
type AgreementStatus int32

const (
	AgreementStatus_ACTIVE    AgreementStatus = 0
	AgreementStatus_SUSPENDED AgreementStatus = 1
	AgreementStatus_EXPIRED   AgreementStatus = 2
	AgreementStatus_REVOKED   AgreementStatus = 3
)

// IdentityShareAuditLog represents an audit entry for identity sharing
type IdentityShareAuditLog struct {
	LogID           string                 `json:"log_id"`
	RequestID       string                 `json:"request_id"`
	HolderDID       string                 `json:"holder_did"`
	RequesterModule string                 `json:"requester_module"`
	ProviderModule  string                 `json:"provider_module"`
	Action          ShareAuditAction       `json:"action"`
	SharedData      []string               `json:"shared_data"`
	Purpose         string                 `json:"purpose"`
	Timestamp       time.Time              `json:"timestamp"`
	UserAgent       string                 `json:"user_agent,omitempty"`
	IPAddress       string                 `json:"ip_address,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// ShareAuditAction represents the type of sharing action
type ShareAuditAction int32

const (
	ShareAuditAction_REQUEST_CREATED  ShareAuditAction = 0
	ShareAuditAction_REQUEST_APPROVED ShareAuditAction = 1
	ShareAuditAction_REQUEST_DENIED   ShareAuditAction = 2
	ShareAuditAction_DATA_ACCESSED    ShareAuditAction = 3
	ShareAuditAction_REQUEST_EXPIRED  ShareAuditAction = 4
	ShareAuditAction_REQUEST_REVOKED  ShareAuditAction = 5
)

// ModuleCapabilities represents what a module can request and provide
type ModuleCapabilities struct {
	ModuleName       string            `json:"module_name"`
	CanRequest       []string          `json:"can_request"`       // Data types this module can request
	CanProvide       []string          `json:"can_provide"`       // Data types this module can provide
	TrustLevel       string            `json:"trust_level"`       // Module's trust level
	Certifications   []string          `json:"certifications"`   // Module certifications
	SecurityPolicies map[string]string `json:"security_policies"` // Security policies
	ContactInfo      string            `json:"contact_info"`      // Module maintainer contact
	Version          string            `json:"version"`           // Module version
	LastUpdated      time.Time         `json:"last_updated"`
}

// IdentityAccessPolicy defines access control policies for identity sharing
type IdentityAccessPolicy struct {
	PolicyID        string                 `json:"policy_id"`
	HolderDID       string                 `json:"holder_did"`
	AllowedModules  []string               `json:"allowed_modules"`
	DeniedModules   []string               `json:"denied_modules"`
	DataRestrictions map[string][]string   `json:"data_restrictions"` // credential_type -> allowed_attributes
	PurposeRestrictions []string           `json:"purpose_restrictions"`
	TimeRestrictions TimeRestriction       `json:"time_restrictions"`
	GeographicRestrictions []string        `json:"geographic_restrictions"`
	MaxSharesPerDay int                    `json:"max_shares_per_day"`
	RequireExplicitConsent bool            `json:"require_explicit_consent"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// TimeRestriction defines time-based access restrictions
type TimeRestriction struct {
	AllowedHours    []int     `json:"allowed_hours,omitempty"`    // Hours of day (0-23)
	AllowedDays     []int     `json:"allowed_days,omitempty"`     // Days of week (0-6)
	AllowedDateFrom time.Time `json:"allowed_date_from,omitempty"`
	AllowedDateTo   time.Time `json:"allowed_date_to,omitempty"`
	TimeZone        string    `json:"time_zone,omitempty"`
}

// Validation functions

// ValidateBasic performs basic validation of IdentityShareRequest
func (r *IdentityShareRequest) ValidateBasic() error {
	if r.RequestID == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "request ID cannot be empty")
	}
	if r.RequesterModule == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "requester module cannot be empty")
	}
	if r.ProviderModule == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "provider module cannot be empty")
	}
	if r.HolderDID == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "holder DID cannot be empty")
	}
	if len(r.RequestedData) == 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "at least one data request is required")
	}
	if r.Purpose == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "purpose cannot be empty")
	}
	if r.TTL <= 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "TTL must be positive")
	}

	// Validate individual data requests
	for i, dataReq := range r.RequestedData {
		if err := dataReq.ValidateBasic(); err != nil {
			return sdkerrors.Wrapf(err, "invalid data request at index %d", i)
		}
	}

	return nil
}

// ValidateBasic performs basic validation of DataRequest
func (dr *DataRequest) ValidateBasic() error {
	if dr.CredentialType == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "credential type cannot be empty")
	}
	if len(dr.Attributes) == 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "at least one attribute must be specified")
	}
	return nil
}

// ValidateBasic performs basic validation of IdentityShareAgreement
func (a *IdentityShareAgreement) ValidateBasic() error {
	if a.AgreementID == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "agreement ID cannot be empty")
	}
	if a.RequesterModule == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "requester module cannot be empty")
	}
	if a.ProviderModule == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "provider module cannot be empty")
	}
	if len(a.AllowedDataTypes) == 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "at least one allowed data type is required")
	}
	if len(a.Purposes) == 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "at least one purpose is required")
	}
	if a.MaxTTL <= 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "max TTL must be positive")
	}
	return nil
}

// ValidateBasic performs basic validation of IdentityAccessPolicy
func (p *IdentityAccessPolicy) ValidateBasic() error {
	if p.PolicyID == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "policy ID cannot be empty")
	}
	if p.HolderDID == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "holder DID cannot be empty")
	}
	if p.MaxSharesPerDay < 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "max shares per day cannot be negative")
	}
	return nil
}

// Helper functions

// IsExpired checks if a share request is expired
func (r *IdentityShareRequest) IsExpired() bool {
	return time.Now().After(r.RequestedAt.Add(r.TTL))
}

// IsActive checks if a share agreement is active
func (a *IdentityShareAgreement) IsActive() bool {
	return a.Status == AgreementStatus_ACTIVE && time.Now().Before(a.ExpiresAt)
}

// CanAutoApprove checks if a request can be auto-approved based on agreement
func (a *IdentityShareAgreement) CanAutoApprove(request *IdentityShareRequest) bool {
	if !a.AutoApprove || !a.IsActive() {
		return false
	}
	
	// Check if requester and provider match
	if a.RequesterModule != request.RequesterModule || a.ProviderModule != request.ProviderModule {
		return false
	}
	
	// Check if TTL is within allowed range
	if request.TTL > a.MaxTTL {
		return false
	}
	
	// Check if purpose is allowed
	purposeAllowed := false
	for _, allowedPurpose := range a.Purposes {
		if allowedPurpose == request.Purpose {
			purposeAllowed = true
			break
		}
	}
	if !purposeAllowed {
		return false
	}
	
	// Check if all requested data types are allowed
	for _, dataReq := range request.RequestedData {
		dataTypeAllowed := false
		for _, allowedType := range a.AllowedDataTypes {
			if allowedType == dataReq.CredentialType {
				dataTypeAllowed = true
				break
			}
		}
		if !dataTypeAllowed {
			return false
		}
	}
	
	return true
}

// Error definitions
var (
	ErrInvalidShareRequest    = sdkerrors.Register(ModuleName, 3001, "invalid share request")
	ErrShareRequestNotFound   = sdkerrors.Register(ModuleName, 3002, "share request not found")
	ErrShareRequestExpired    = sdkerrors.Register(ModuleName, 3003, "share request expired")
	ErrShareRequestDenied     = sdkerrors.Register(ModuleName, 3004, "share request denied")
	ErrInvalidAgreement       = sdkerrors.Register(ModuleName, 3005, "invalid sharing agreement")
	ErrAgreementNotFound      = sdkerrors.Register(ModuleName, 3006, "sharing agreement not found")
	ErrUnauthorizedModule     = sdkerrors.Register(ModuleName, 3007, "unauthorized module")
	ErrInsufficientTrustLevel = sdkerrors.Register(ModuleName, 3008, "insufficient trust level")
	ErrAccessPolicyViolation  = sdkerrors.Register(ModuleName, 3009, "access policy violation")
	ErrDailyLimitExceeded     = sdkerrors.Register(ModuleName, 3010, "daily sharing limit exceeded")
)
