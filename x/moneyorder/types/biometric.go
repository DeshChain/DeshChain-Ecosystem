package types

import (
	"time"
)

// BiometricType represents the type of biometric authentication
type BiometricType int32

const (
	BiometricType_FINGERPRINT BiometricType = 0
	BiometricType_FACE        BiometricType = 1
	BiometricType_IRIS        BiometricType = 2
	BiometricType_VOICE       BiometricType = 3
	BiometricType_PALM        BiometricType = 4
)

// String returns the string representation of BiometricType
func (bt BiometricType) String() string {
	switch bt {
	case BiometricType_FINGERPRINT:
		return "FINGERPRINT"
	case BiometricType_FACE:
		return "FACE"
	case BiometricType_IRIS:
		return "IRIS"
	case BiometricType_VOICE:
		return "VOICE"
	case BiometricType_PALM:
		return "PALM"
	default:
		return "UNKNOWN"
	}
}

// SecurityLevel represents the security level based on biometric configuration
type SecurityLevel int32

const (
	SecurityLevel_BASIC    SecurityLevel = 0
	SecurityLevel_ENHANCED SecurityLevel = 1
	SecurityLevel_PREMIUM  SecurityLevel = 2
)

// String returns the string representation of SecurityLevel
func (sl SecurityLevel) String() string {
	switch sl {
	case SecurityLevel_BASIC:
		return "BASIC"
	case SecurityLevel_ENHANCED:
		return "ENHANCED"
	case SecurityLevel_PREMIUM:
		return "PREMIUM"
	default:
		return "UNKNOWN"
	}
}

// BiometricRegistration represents a user's biometric registration
type BiometricRegistration struct {
	BiometricId     string        `json:"biometric_id"`
	UserAddress     string        `json:"user_address"`
	BiometricType   BiometricType `json:"biometric_type"`
	TemplateHash    string        `json:"template_hash"`
	DeviceId        string        `json:"device_id"`
	RegisteredAt    time.Time     `json:"registered_at"`
	IsActive        bool          `json:"is_active"`
	FailCount       int32         `json:"fail_count"`
	LastUsed        time.Time     `json:"last_used"`
	DisabledAt      time.Time     `json:"disabled_at,omitempty"`
	DisabledReason  string        `json:"disabled_reason,omitempty"`
}

// BiometricAuthResult represents the result of a biometric authentication
type BiometricAuthResult struct {
	Success      bool    `json:"success"`
	ErrorMessage string  `json:"error_message,omitempty"`
	AuthScore    float64 `json:"auth_score"`
	BiometricId  string  `json:"biometric_id,omitempty"`
}

// BiometricAuthAttempt represents an authentication attempt for audit purposes
type BiometricAuthAttempt struct {
	UserAddress string    `json:"user_address"`
	BiometricId string    `json:"biometric_id"`
	Success     bool      `json:"success"`
	AuthScore   float64   `json:"auth_score"`
	AttemptTime time.Time `json:"attempt_time"`
	ClientIP    string    `json:"client_ip,omitempty"`
}

// BiometricSecurityConfig represents system-wide biometric security configuration
type BiometricSecurityConfig struct {
	MinAuthScore          float64 `json:"min_auth_score"`
	MaxFailAttempts       int32   `json:"max_fail_attempts"`
	LockoutDurationHours  int32   `json:"lockout_duration_hours"`
	RequiredForHighValue  bool    `json:"required_for_high_value"`
	HighValueThreshold    int64   `json:"high_value_threshold"`
	MultiFactorRequired   bool    `json:"multi_factor_required"`
	AllowedBiometricTypes []BiometricType `json:"allowed_biometric_types"`
}

// DefaultBiometricSecurityConfig returns the default biometric security configuration
func DefaultBiometricSecurityConfig() BiometricSecurityConfig {
	return BiometricSecurityConfig{
		MinAuthScore:         0.85,
		MaxFailAttempts:      5,
		LockoutDurationHours: 24,
		RequiredForHighValue: true,
		HighValueThreshold:   100000, // 1 Lakh NAMO
		MultiFactorRequired:  false,
		AllowedBiometricTypes: []BiometricType{
			BiometricType_FINGERPRINT,
			BiometricType_FACE,
			BiometricType_IRIS,
		},
	}
}

// Key prefixes for biometric data storage
var (
	BiometricRegistrationPrefix = []byte{0x10}
	BiometricAuthAttemptPrefix  = []byte{0x11}
	BiometricConfigPrefix       = []byte{0x12}
)

// Event types for biometric authentication
const (
	EventTypeBiometricRegistered   = "biometric_registered"
	EventTypeBiometricAuthSuccess  = "biometric_auth_success"
	EventTypeBiometricAuthFailed   = "biometric_auth_failed"
	EventTypeBiometricDisabled     = "biometric_disabled"
	EventTypeMoneyOrderBiometricAuth = "money_order_biometric_auth"
)

// Attribute keys for biometric events
const (
	AttributeKeyBiometricType = "biometric_type"
	AttributeKeyBiometricID   = "biometric_id"
	AttributeKeyAuthScore     = "auth_score"
	AttributeKeyFailCount     = "fail_count"
	AttributeKeyDeviceID      = "device_id"
	AttributeKeyReason        = "reason"
)