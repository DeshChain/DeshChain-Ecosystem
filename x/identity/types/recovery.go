package types

import (
	"time"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Identity Backup and Recovery System Types

// IdentityBackup represents a complete backup of an identity
type IdentityBackup struct {
	BackupID           string                 `json:"backup_id"`
	HolderDID          string                 `json:"holder_did"`
	BackupVersion      int64                  `json:"backup_version"`
	CreatedAt          time.Time              `json:"created_at"`
	ExpiresAt          time.Time              `json:"expires_at"`
	EncryptionMethod   string                 `json:"encryption_method"`
	BackupData         EncryptedBackupData    `json:"backup_data"`
	RecoveryMethods    []RecoveryMethod       `json:"recovery_methods"`
	IntegrityHash      string                 `json:"integrity_hash"`
	Status             BackupStatus           `json:"status"`
	Metadata           map[string]interface{} `json:"metadata,omitempty"`
}

// EncryptedBackupData contains the encrypted identity data
type EncryptedBackupData struct {
	IdentityData       []byte `json:"identity_data"`        // Encrypted identity
	CredentialsData    []byte `json:"credentials_data"`     // Encrypted credentials
	BiometricData      []byte `json:"biometric_data"`       // Encrypted biometric templates
	ConsentData        []byte `json:"consent_data"`         // Encrypted consent records
	AccessPolicyData   []byte `json:"access_policy_data"`   // Encrypted access policies
	ZKProofData        []byte `json:"zk_proof_data"`        // Encrypted ZK proof data
	EncryptionKeyInfo  []byte `json:"encryption_key_info"`  // Encrypted key derivation info
}

// RecoveryMethod defines different ways to recover an identity
type RecoveryMethod struct {
	MethodID           string                 `json:"method_id"`
	MethodType         RecoveryMethodType     `json:"method_type"`
	MethodName         string                 `json:"method_name"`
	Configuration      map[string]interface{} `json:"configuration"`
	TrustLevel         string                 `json:"trust_level"`
	RequiredConfidence int                    `json:"required_confidence"` // 1-100
	Enabled            bool                   `json:"enabled"`
	CreatedAt          time.Time              `json:"created_at"`
	LastUsed           *time.Time             `json:"last_used,omitempty"`
	UsageCount         int64                  `json:"usage_count"`
	Metadata           map[string]interface{} `json:"metadata,omitempty"`
}

// RecoveryMethodType represents different recovery method types
type RecoveryMethodType int32

const (
	RecoveryMethodType_MNEMONIC_PHRASE         RecoveryMethodType = 0
	RecoveryMethodType_SOCIAL_RECOVERY         RecoveryMethodType = 1
	RecoveryMethodType_GUARDIAN_MULTISIG       RecoveryMethodType = 2
	RecoveryMethodType_BIOMETRIC_BACKUP        RecoveryMethodType = 3
	RecoveryMethodType_HARDWARE_KEY            RecoveryMethodType = 4
	RecoveryMethodType_BACKUP_CODES            RecoveryMethodType = 5
	RecoveryMethodType_EMAIL_VERIFICATION      RecoveryMethodType = 6
	RecoveryMethodType_SMS_VERIFICATION        RecoveryMethodType = 7
	RecoveryMethodType_IDENTITY_PROVIDER       RecoveryMethodType = 8
	RecoveryMethodType_INSTITUTIONAL_RECOVERY  RecoveryMethodType = 9
	RecoveryMethodType_ZERO_KNOWLEDGE_PROOF    RecoveryMethodType = 10
)

// BackupStatus represents the status of a backup
type BackupStatus int32

const (
	BackupStatus_ACTIVE     BackupStatus = 0
	BackupStatus_EXPIRED    BackupStatus = 1
	BackupStatus_REVOKED    BackupStatus = 2
	BackupStatus_CORRUPTED  BackupStatus = 3
	BackupStatus_RECOVERING BackupStatus = 4
)

// RecoveryRequest represents a request to recover an identity
type RecoveryRequest struct {
	RequestID          string                 `json:"request_id"`
	HolderDID          string                 `json:"holder_did"`
	BackupID           string                 `json:"backup_id"`
	RequestedBy        string                 `json:"requested_by"`
	RecoveryMethods    []RecoveryAttempt      `json:"recovery_methods"`
	Reason             string                 `json:"reason"`
	RequestedAt        time.Time              `json:"requested_at"`
	ExpiresAt          time.Time              `json:"expires_at"`
	Status             RecoveryRequestStatus  `json:"status"`
	ConfidenceScore    int                    `json:"confidence_score"`
	RequiredScore      int                    `json:"required_score"`
	AttemptCount       int                    `json:"attempt_count"`
	MaxAttempts        int                    `json:"max_attempts"`
	Metadata           map[string]interface{} `json:"metadata,omitempty"`
}

// RecoveryAttempt represents an attempt to use a recovery method
type RecoveryAttempt struct {
	AttemptID       string                 `json:"attempt_id"`
	MethodID        string                 `json:"method_id"`
	MethodType      RecoveryMethodType     `json:"method_type"`
	ProofData       []byte                 `json:"proof_data"`
	AttemptedAt     time.Time              `json:"attempted_at"`
	Status          AttemptStatus          `json:"status"`
	Confidence      int                    `json:"confidence"`
	ErrorMessage    string                 `json:"error_message,omitempty"`
	VerificationData map[string]interface{} `json:"verification_data,omitempty"`
}

// RecoveryRequestStatus represents the status of a recovery request
type RecoveryRequestStatus int32

const (
	RecoveryRequestStatus_PENDING    RecoveryRequestStatus = 0
	RecoveryRequestStatus_APPROVED   RecoveryRequestStatus = 1
	RecoveryRequestStatus_REJECTED   RecoveryRequestStatus = 2
	RecoveryRequestStatus_EXPIRED    RecoveryRequestStatus = 3
	RecoveryRequestStatus_COMPLETED  RecoveryRequestStatus = 4
	RecoveryRequestStatus_FAILED     RecoveryRequestStatus = 5
)

// AttemptStatus represents the status of a recovery attempt
type AttemptStatus int32

const (
	AttemptStatus_PENDING    AttemptStatus = 0
	AttemptStatus_VERIFIED   AttemptStatus = 1
	AttemptStatus_FAILED     AttemptStatus = 2
	AttemptStatus_EXPIRED    AttemptStatus = 3
)

// SocialRecoveryGuardian represents a trusted guardian for social recovery
type SocialRecoveryGuardian struct {
	GuardianID      string                 `json:"guardian_id"`
	GuardianDID     string                 `json:"guardian_did"`
	GuardianAddress string                 `json:"guardian_address"`
	GuardianName    string                 `json:"guardian_name"`
	TrustLevel      string                 `json:"trust_level"`
	Weight          int                    `json:"weight"`           // Voting weight
	ContactInfo     string                 `json:"contact_info"`     // Encrypted contact info
	PublicKey       string                 `json:"public_key"`       // Guardian's public key
	Status          GuardianStatus         `json:"status"`
	AddedAt         time.Time              `json:"added_at"`
	LastActive      *time.Time             `json:"last_active,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// GuardianStatus represents the status of a guardian
type GuardianStatus int32

const (
	GuardianStatus_ACTIVE     GuardianStatus = 0
	GuardianStatus_INACTIVE   GuardianStatus = 1
	GuardianStatus_SUSPENDED  GuardianStatus = 2
	GuardianStatus_REVOKED    GuardianStatus = 3
)

// GuardianVote represents a guardian's vote on a recovery request
type GuardianVote struct {
	VoteID        string                 `json:"vote_id"`
	RequestID     string                 `json:"request_id"`
	GuardianID    string                 `json:"guardian_id"`
	Vote          VoteType               `json:"vote"`
	Reason        string                 `json:"reason,omitempty"`
	Signature     string                 `json:"signature"`
	VotedAt       time.Time              `json:"voted_at"`
	Weight        int                    `json:"weight"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// VoteType represents the type of vote
type VoteType int32

const (
	VoteType_APPROVE VoteType = 0
	VoteType_REJECT  VoteType = 1
	VoteType_ABSTAIN VoteType = 2
)

// DisasterRecoveryConfig represents disaster recovery configuration
type DisasterRecoveryConfig struct {
	ConfigID                string                 `json:"config_id"`
	HolderDID               string                 `json:"holder_did"`
	BackupFrequency         time.Duration          `json:"backup_frequency"`
	BackupRetention         time.Duration          `json:"backup_retention"`
	AutoBackupEnabled       bool                   `json:"auto_backup_enabled"`
	CrossChainBackup        bool                   `json:"cross_chain_backup"`
	OffChainBackupEnabled   bool                   `json:"off_chain_backup_enabled"`
	RecoveryTimeObjective   time.Duration          `json:"recovery_time_objective"`   // RTO
	RecoveryPointObjective  time.Duration          `json:"recovery_point_objective"`  // RPO
	MinGuardians            int                    `json:"min_guardians"`
	GuardianThreshold       int                    `json:"guardian_threshold"`        // Required votes
	EmergencyContacts       []EmergencyContact     `json:"emergency_contacts"`
	NotificationSettings   NotificationSettings   `json:"notification_settings"`
	CreatedAt               time.Time              `json:"created_at"`
	UpdatedAt               time.Time              `json:"updated_at"`
	Metadata                map[string]interface{} `json:"metadata,omitempty"`
}

// EmergencyContact represents an emergency contact for recovery
type EmergencyContact struct {
	ContactID    string                 `json:"contact_id"`
	Name         string                 `json:"name"`
	ContactType  string                 `json:"contact_type"` // email, sms, etc.
	ContactInfo  string                 `json:"contact_info"` // Encrypted
	Priority     int                    `json:"priority"`
	Verified     bool                   `json:"verified"`
	Enabled      bool                   `json:"enabled"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// NotificationSettings represents notification preferences
type NotificationSettings struct {
	EmailEnabled      bool   `json:"email_enabled"`
	SMSEnabled        bool   `json:"sms_enabled"`
	PushEnabled       bool   `json:"push_enabled"`
	WebhookEnabled    bool   `json:"webhook_enabled"`
	WebhookURL        string `json:"webhook_url,omitempty"`
	NotifyOnBackup    bool   `json:"notify_on_backup"`
	NotifyOnRecovery  bool   `json:"notify_on_recovery"`
	NotifyOnFailure   bool   `json:"notify_on_failure"`
}

// BackupVerificationResult represents the result of backup verification
type BackupVerificationResult struct {
	BackupID           string                 `json:"backup_id"`
	VerificationID     string                 `json:"verification_id"`
	VerifiedAt         time.Time              `json:"verified_at"`
	IntegrityValid     bool                   `json:"integrity_valid"`
	DecryptionValid    bool                   `json:"decryption_valid"`
	DataCompleteness   float64                `json:"data_completeness"` // 0-100%
	RecoverabilityScore int                   `json:"recoverability_score"` // 0-100
	IssuesFound        []string               `json:"issues_found,omitempty"`
	Recommendations    []string               `json:"recommendations,omitempty"`
	Metadata           map[string]interface{} `json:"metadata,omitempty"`
}

// Validation functions

// ValidateBasic performs basic validation of IdentityBackup
func (b *IdentityBackup) ValidateBasic() error {
	if b.BackupID == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "backup ID cannot be empty")
	}
	if b.HolderDID == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "holder DID cannot be empty")
	}
	if b.BackupVersion <= 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "backup version must be positive")
	}
	if b.EncryptionMethod == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "encryption method cannot be empty")
	}
	if len(b.RecoveryMethods) == 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "at least one recovery method is required")
	}
	if b.IntegrityHash == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "integrity hash cannot be empty")
	}

	// Validate individual recovery methods
	for i, method := range b.RecoveryMethods {
		if err := method.ValidateBasic(); err != nil {
			return sdkerrors.Wrapf(err, "invalid recovery method at index %d", i)
		}
	}

	return nil
}

// ValidateBasic performs basic validation of RecoveryMethod
func (rm *RecoveryMethod) ValidateBasic() error {
	if rm.MethodID == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "method ID cannot be empty")
	}
	if rm.MethodName == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "method name cannot be empty")
	}
	if rm.RequiredConfidence < 1 || rm.RequiredConfidence > 100 {
		return sdkerrors.Wrap(ErrInvalidRequest, "required confidence must be between 1 and 100")
	}
	return nil
}

// ValidateBasic performs basic validation of RecoveryRequest
func (rr *RecoveryRequest) ValidateBasic() error {
	if rr.RequestID == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "request ID cannot be empty")
	}
	if rr.HolderDID == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "holder DID cannot be empty")
	}
	if rr.BackupID == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "backup ID cannot be empty")
	}
	if rr.RequestedBy == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "requested by cannot be empty")
	}
	if len(rr.RecoveryMethods) == 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "at least one recovery method is required")
	}
	if rr.RequiredScore < 1 || rr.RequiredScore > 100 {
		return sdkerrors.Wrap(ErrInvalidRequest, "required score must be between 1 and 100")
	}
	if rr.MaxAttempts <= 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "max attempts must be positive")
	}
	return nil
}

// ValidateBasic performs basic validation of SocialRecoveryGuardian
func (g *SocialRecoveryGuardian) ValidateBasic() error {
	if g.GuardianID == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "guardian ID cannot be empty")
	}
	if g.GuardianDID == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "guardian DID cannot be empty")
	}
	if g.GuardianAddress == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "guardian address cannot be empty")
	}
	if g.Weight <= 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "guardian weight must be positive")
	}
	if g.PublicKey == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "guardian public key cannot be empty")
	}
	return nil
}

// Helper functions

// IsExpired checks if a backup is expired
func (b *IdentityBackup) IsExpired() bool {
	return time.Now().After(b.ExpiresAt)
}

// IsValid checks if a backup is valid and usable
func (b *IdentityBackup) IsValid() bool {
	return b.Status == BackupStatus_ACTIVE && !b.IsExpired()
}

// IsExpired checks if a recovery request is expired
func (rr *RecoveryRequest) IsExpired() bool {
	return time.Now().After(rr.ExpiresAt)
}

// CanAttemptRecovery checks if more recovery attempts are allowed
func (rr *RecoveryRequest) CanAttemptRecovery() bool {
	return rr.AttemptCount < rr.MaxAttempts && !rr.IsExpired() && rr.Status == RecoveryRequestStatus_PENDING
}

// HasSufficientConfidence checks if the recovery request has sufficient confidence
func (rr *RecoveryRequest) HasSufficientConfidence() bool {
	return rr.ConfidenceScore >= rr.RequiredScore
}

// IsActive checks if a guardian is active
func (g *SocialRecoveryGuardian) IsActive() bool {
	return g.Status == GuardianStatus_ACTIVE
}

// GetMethodTypeString returns a string representation of the recovery method type
func (t RecoveryMethodType) String() string {
	switch t {
	case RecoveryMethodType_MNEMONIC_PHRASE:
		return "mnemonic_phrase"
	case RecoveryMethodType_SOCIAL_RECOVERY:
		return "social_recovery"
	case RecoveryMethodType_GUARDIAN_MULTISIG:
		return "guardian_multisig"
	case RecoveryMethodType_BIOMETRIC_BACKUP:
		return "biometric_backup"
	case RecoveryMethodType_HARDWARE_KEY:
		return "hardware_key"
	case RecoveryMethodType_BACKUP_CODES:
		return "backup_codes"
	case RecoveryMethodType_EMAIL_VERIFICATION:
		return "email_verification"
	case RecoveryMethodType_SMS_VERIFICATION:
		return "sms_verification"
	case RecoveryMethodType_IDENTITY_PROVIDER:
		return "identity_provider"
	case RecoveryMethodType_INSTITUTIONAL_RECOVERY:
		return "institutional_recovery"
	case RecoveryMethodType_ZERO_KNOWLEDGE_PROOF:
		return "zero_knowledge_proof"
	default:
		return "unknown"
	}
}

// Error definitions
var (
	ErrBackupNotFound         = sdkerrors.Register(ModuleName, 4001, "backup not found")
	ErrBackupExpired          = sdkerrors.Register(ModuleName, 4002, "backup expired")
	ErrBackupCorrupted        = sdkerrors.Register(ModuleName, 4003, "backup corrupted")
	ErrRecoveryRequestNotFound = sdkerrors.Register(ModuleName, 4004, "recovery request not found")
	ErrRecoveryRequestExpired = sdkerrors.Register(ModuleName, 4005, "recovery request expired")
	ErrInsufficientConfidence = sdkerrors.Register(ModuleName, 4006, "insufficient confidence for recovery")
	ErrMaxAttemptsExceeded    = sdkerrors.Register(ModuleName, 4007, "maximum recovery attempts exceeded")
	ErrGuardianNotFound       = sdkerrors.Register(ModuleName, 4008, "guardian not found")
	ErrInvalidRecoveryMethod  = sdkerrors.Register(ModuleName, 4009, "invalid recovery method")
	ErrRecoveryNotAllowed     = sdkerrors.Register(ModuleName, 4010, "recovery not allowed")
)
