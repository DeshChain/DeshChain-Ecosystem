package types

import (
	"time"
)

// Query request and response types for offline verification

// QueryOfflineConfigRequest is the request for offline verification configuration
type QueryOfflineConfigRequest struct{}

// QueryOfflineConfigResponse is the response for offline verification configuration
type QueryOfflineConfigResponse struct {
	Config *OfflineVerificationConfig `json:"config" yaml:"config"`
}

// QueryOfflineDevicesRequest is the request for offline devices
type QueryOfflineDevicesRequest struct {
	Did      string `json:"did" yaml:"did"`
	DeviceId string `json:"device_id,omitempty" yaml:"device_id,omitempty"`
}

// QueryOfflineDevicesResponse is the response for offline devices
type QueryOfflineDevicesResponse struct {
	Devices []*OfflineDevice `json:"devices" yaml:"devices"`
	Total   int              `json:"total" yaml:"total"`
}

// QueryOfflineFormatsRequest is the request for supported offline formats
type QueryOfflineFormatsRequest struct{}

// QueryOfflineFormatsResponse is the response for supported offline formats
type QueryOfflineFormatsResponse struct {
	Formats []*OfflineFormatInfo `json:"formats" yaml:"formats"`
}

// QueryOfflineModesRequest is the request for supported offline modes
type QueryOfflineModesRequest struct{}

// QueryOfflineModesResponse is the response for supported offline modes
type QueryOfflineModesResponse struct {
	Modes []*OfflineModeInfo `json:"modes" yaml:"modes"`
}

// QueryVerifyOfflineRequest is the request for offline verification
type QueryVerifyOfflineRequest struct {
	OfflineData *OfflineVerificationData    `json:"offline_data" yaml:"offline_data"`
	Request     *OfflineVerificationRequest `json:"request" yaml:"request"`
}

// QueryVerifyOfflineResponse is the response for offline verification
type QueryVerifyOfflineResponse struct {
	Result *OfflineVerificationResult `json:"result" yaml:"result"`
}

// QueryValidateOfflinePackageRequest is the request for package validation
type QueryValidateOfflinePackageRequest struct {
	PackageData []byte `json:"package_data" yaml:"package_data"`
	Format      string `json:"format" yaml:"format"`
}

// QueryValidateOfflinePackageResponse is the response for package validation
type QueryValidateOfflinePackageResponse struct {
	Valid      bool                         `json:"valid" yaml:"valid"`
	Errors     []string                     `json:"errors,omitempty" yaml:"errors,omitempty"`
	Warnings   []string                     `json:"warnings,omitempty" yaml:"warnings,omitempty"`
	PackageInfo *OfflinePackageInfo         `json:"package_info,omitempty" yaml:"package_info,omitempty"`
}

// Supporting types for offline verification queries

// OfflineDevice represents a registered offline device
type OfflineDevice struct {
	DID                string                  `json:"did" yaml:"did"`
	DeviceID           string                  `json:"device_id" yaml:"device_id"`
	DeviceName         string                  `json:"device_name" yaml:"device_name"`
	DeviceType         string                  `json:"device_type" yaml:"device_type"`
	PublicKey          string                  `json:"public_key" yaml:"public_key"`
	Capabilities       []string                `json:"capabilities" yaml:"capabilities"`
	SecurityLevel      uint32                  `json:"security_level" yaml:"security_level"`
	MaxOfflineDuration uint32                  `json:"max_offline_duration" yaml:"max_offline_duration"`
	RegisteredAt       time.Time               `json:"registered_at" yaml:"registered_at"`
	LastUsed           *time.Time              `json:"last_used,omitempty" yaml:"last_used,omitempty"`
	Status             OfflineDeviceStatus     `json:"status" yaml:"status"`
	UsageCount         uint64                  `json:"usage_count" yaml:"usage_count"`
}

// OfflineDeviceStatus represents the status of an offline device
type OfflineDeviceStatus string

const (
	DeviceStatusActive    OfflineDeviceStatus = "active"
	DeviceStatusInactive  OfflineDeviceStatus = "inactive"
	DeviceStatusRevoked   OfflineDeviceStatus = "revoked"
	DeviceStatusSuspended OfflineDeviceStatus = "suspended"
)

// OfflineFormatInfo provides information about an offline format
type OfflineFormatInfo struct {
	Format      OfflineCredentialFormat `json:"format" yaml:"format"`
	Name        string                  `json:"name" yaml:"name"`
	Description string                  `json:"description" yaml:"description"`
	MaxSize     uint64                  `json:"max_size" yaml:"max_size"`
	Compressed  bool                    `json:"compressed" yaml:"compressed"`
	Portable    bool                    `json:"portable" yaml:"portable"`
	HumanReadable bool                  `json:"human_readable" yaml:"human_readable"`
	UseCases    []string                `json:"use_cases" yaml:"use_cases"`
}

// OfflineModeInfo provides information about an offline verification mode
type OfflineModeInfo struct {
	Mode            OfflineVerificationMode `json:"mode" yaml:"mode"`
	Name            string                  `json:"name" yaml:"name"`
	Description     string                  `json:"description" yaml:"description"`
	SecurityLevel   uint32                  `json:"security_level" yaml:"security_level"`
	RequiredData    []string                `json:"required_data" yaml:"required_data"`
	OptionalData    []string                `json:"optional_data" yaml:"optional_data"`
	MaxDuration     time.Duration           `json:"max_duration" yaml:"max_duration"`
	UseCases        []string                `json:"use_cases" yaml:"use_cases"`
}

// OfflinePackageInfo provides information about an offline package
type OfflinePackageInfo struct {
	DID            string                  `json:"did" yaml:"did"`
	Format         OfflineCredentialFormat `json:"format" yaml:"format"`
	Size           uint64                  `json:"size" yaml:"size"`
	Compressed     bool                    `json:"compressed" yaml:"compressed"`
	CreatedAt      time.Time               `json:"created_at" yaml:"created_at"`
	ExpiresAt      time.Time               `json:"expires_at" yaml:"expires_at"`
	IsExpired      bool                    `json:"is_expired" yaml:"is_expired"`
	CredentialCount int                    `json:"credential_count" yaml:"credential_count"`
	HasBiometrics  bool                    `json:"has_biometrics" yaml:"has_biometrics"`
	HasRecovery    bool                    `json:"has_recovery" yaml:"has_recovery"`
	DataIntegrity  bool                    `json:"data_integrity" yaml:"data_integrity"`
	SignatureValid bool                    `json:"signature_valid" yaml:"signature_valid"`
}

// Validation functions for query types

// Validate validates QueryOfflineDevicesRequest
func (q *QueryOfflineDevicesRequest) Validate() error {
	return ValidateUserDID(q.Did)
}

// Validate validates QueryVerifyOfflineRequest
func (q *QueryVerifyOfflineRequest) Validate() error {
	if q.OfflineData == nil {
		return ErrInvalidRequest
	}

	if q.Request == nil {
		return ErrInvalidRequest
	}

	if err := q.OfflineData.Validate(); err != nil {
		return err
	}

	return nil
}

// Validate validates QueryValidateOfflinePackageRequest
func (q *QueryValidateOfflinePackageRequest) Validate() error {
	if len(q.PackageData) == 0 {
		return ErrInvalidRequest
	}

	if q.Format != "" {
		if err := ValidateOfflineFormat(OfflineCredentialFormat(q.Format)); err != nil {
			return err
		}
	}

	return nil
}

// Constructor functions for query responses

// NewQueryOfflineConfigResponse creates a new offline config response
func NewQueryOfflineConfigResponse(config *OfflineVerificationConfig) *QueryOfflineConfigResponse {
	return &QueryOfflineConfigResponse{
		Config: config,
	}
}

// NewQueryOfflineDevicesResponse creates a new offline devices response
func NewQueryOfflineDevicesResponse(devices []*OfflineDevice) *QueryOfflineDevicesResponse {
	return &QueryOfflineDevicesResponse{
		Devices: devices,
		Total:   len(devices),
	}
}

// NewQueryOfflineFormatsResponse creates a new offline formats response
func NewQueryOfflineFormatsResponse() *QueryOfflineFormatsResponse {
	formats := []*OfflineFormatInfo{
		{
			Format:        FormatSelfContained,
			Name:          "Self-Contained",
			Description:   "Complete offline verification package with all required data",
			MaxSize:       10 * 1024 * 1024, // 10MB
			Compressed:    false,
			Portable:      true,
			HumanReadable: false,
			UseCases:      []string{"comprehensive_verification", "backup", "disaster_recovery"},
		},
		{
			Format:        FormatCompressed,
			Name:          "Compressed",
			Description:   "Space-efficient compressed format for bandwidth-limited scenarios",
			MaxSize:       2 * 1024 * 1024, // 2MB
			Compressed:    true,
			Portable:      true,
			HumanReadable: false,
			UseCases:      []string{"mobile_apps", "low_bandwidth", "storage_optimization"},
		},
		{
			Format:        FormatQRCode,
			Name:          "QR Code",
			Description:   "QR code compatible format for visual scanning",
			MaxSize:       4 * 1024, // 4KB
			Compressed:    true,
			Portable:      true,
			HumanReadable: true,
			UseCases:      []string{"visual_verification", "mobile_scanning", "paper_backup"},
		},
		{
			Format:        FormatNFC,
			Name:          "NFC",
			Description:   "NFC tag compatible format for contactless verification",
			MaxSize:       8 * 1024, // 8KB
			Compressed:    true,
			Portable:      true,
			HumanReadable: false,
			UseCases:      []string{"contactless_verification", "smart_cards", "iot_devices"},
		},
		{
			Format:        FormatPrintable,
			Name:          "Printable",
			Description:   "Human-readable printable format for paper-based verification",
			MaxSize:       100 * 1024, // 100KB
			Compressed:    false,
			Portable:      false,
			HumanReadable: true,
			UseCases:      []string{"paper_documents", "manual_verification", "emergency_backup"},
		},
	}

	return &QueryOfflineFormatsResponse{
		Formats: formats,
	}
}

// NewQueryOfflineModesResponse creates a new offline modes response
func NewQueryOfflineModesResponse() *QueryOfflineModesResponse {
	modes := []*OfflineModeInfo{
		{
			Mode:          OfflineModeFull,
			Name:          "Full Verification",
			Description:   "Complete offline verification with all security checks",
			SecurityLevel: 5,
			RequiredData:  []string{"identity_proof", "kyc_proof", "credentials"},
			OptionalData:  []string{"biometric_templates", "emergency_contacts"},
			MaxDuration:   24 * time.Hour,
			UseCases:      []string{"high_security", "financial_transactions", "government_services"},
		},
		{
			Mode:          OfflineModePartial,
			Name:          "Partial Verification",
			Description:   "Partial verification with cached data and reduced security",
			SecurityLevel: 3,
			RequiredData:  []string{"identity_proof"},
			OptionalData:  []string{"kyc_proof", "credentials", "biometric_templates"},
			MaxDuration:   7 * 24 * time.Hour,
			UseCases:      []string{"cached_verification", "mobile_apps", "routine_access"},
		},
		{
			Mode:          OfflineModeMinimal,
			Name:          "Minimal Verification",
			Description:   "Basic identity verification with minimal data requirements",
			SecurityLevel: 2,
			RequiredData:  []string{"identity_proof"},
			OptionalData:  []string{},
			MaxDuration:   30 * 24 * time.Hour,
			UseCases:      []string{"basic_identification", "public_services", "information_access"},
		},
		{
			Mode:          OfflineModeEmergency,
			Name:          "Emergency Mode",
			Description:   "Emergency verification with relaxed security for critical situations",
			SecurityLevel: 1,
			RequiredData:  []string{},
			OptionalData:  []string{"identity_proof", "emergency_contacts"},
			MaxDuration:   1 * time.Hour,
			UseCases:      []string{"emergency_access", "disaster_recovery", "critical_services"},
		},
	}

	return &QueryOfflineModesResponse{
		Modes: modes,
	}
}

// NewQueryVerifyOfflineResponse creates a new verify offline response
func NewQueryVerifyOfflineResponse(result *OfflineVerificationResult) *QueryVerifyOfflineResponse {
	return &QueryVerifyOfflineResponse{
		Result: result,
	}
}

// NewQueryValidateOfflinePackageResponse creates a new validate package response
func NewQueryValidateOfflinePackageResponse(valid bool, errors, warnings []string, packageInfo *OfflinePackageInfo) *QueryValidateOfflinePackageResponse {
	return &QueryValidateOfflinePackageResponse{
		Valid:       valid,
		Errors:      errors,
		Warnings:    warnings,
		PackageInfo: packageInfo,
	}
}

// Helper functions for offline device management

// NewOfflineDevice creates a new offline device
func NewOfflineDevice(did, deviceID, deviceName, deviceType, publicKey string, capabilities []string, securityLevel, maxOfflineDuration uint32) *OfflineDevice {
	return &OfflineDevice{
		DID:                did,
		DeviceID:           deviceID,
		DeviceName:         deviceName,
		DeviceType:         deviceType,
		PublicKey:          publicKey,
		Capabilities:       capabilities,
		SecurityLevel:      securityLevel,
		MaxOfflineDuration: maxOfflineDuration,
		RegisteredAt:       time.Now(),
		Status:             DeviceStatusActive,
		UsageCount:         0,
	}
}

// IsActive checks if the device is active
func (d *OfflineDevice) IsActive() bool {
	return d.Status == DeviceStatusActive
}

// CanPerformVerification checks if the device can perform a specific type of verification
func (d *OfflineDevice) CanPerformVerification(verificationType string) bool {
	for _, capability := range d.Capabilities {
		if capability == verificationType || capability == "identity_verification" {
			return true
		}
	}
	return false
}

// UpdateLastUsed updates the last used timestamp and increments usage count
func (d *OfflineDevice) UpdateLastUsed() {
	now := time.Now()
	d.LastUsed = &now
	d.UsageCount++
}

// GetDeviceCapabilityDescription returns human-readable descriptions for device capabilities
func GetDeviceCapabilityDescription(capability string) string {
	descriptions := map[string]string{
		"identity_verification": "Basic identity verification and authentication",
		"biometric_capture":     "Capture and process biometric data",
		"credential_storage":    "Secure storage of verifiable credentials",
		"offline_signing":       "Cryptographic signing without network connectivity",
		"secure_element":        "Hardware-based secure key storage and operations",
	}

	if desc, exists := descriptions[capability]; exists {
		return desc
	}

	return "Unknown capability"
}

// GetDeviceTypeDescription returns human-readable descriptions for device types
func GetDeviceTypeDescription(deviceType string) string {
	descriptions := map[string]string{
		"mobile":         "Mobile phone or smartphone",
		"tablet":         "Tablet device",
		"laptop":         "Laptop computer",
		"desktop":        "Desktop computer",
		"iot":            "Internet of Things device",
		"smartcard":      "Smart card or card reader",
		"hardware_token": "Dedicated hardware security token",
	}

	if desc, exists := descriptions[deviceType]; exists {
		return desc
	}

	return "Unknown device type"
}

// Validation functions for offline device management

// ValidateOfflineDevice validates an offline device
func ValidateOfflineDevice(device *OfflineDevice) error {
	if device.DID == "" {
		return ErrInvalidRequest
	}

	if device.DeviceID == "" {
		return ErrInvalidRequest
	}

	if device.DeviceName == "" {
		return ErrInvalidRequest
	}

	if device.DeviceType == "" {
		return ErrInvalidRequest
	}

	if device.PublicKey == "" {
		return ErrInvalidPublicKey
	}

	if len(device.Capabilities) == 0 {
		return ErrInvalidRequest
	}

	if device.SecurityLevel < 1 || device.SecurityLevel > 5 {
		return ErrInvalidRequest
	}

	if device.MaxOfflineDuration == 0 {
		return ErrInvalidRequest
	}

	return nil
}

// GetSupportedDeviceStatuses returns all supported device statuses
func GetSupportedDeviceStatuses() []OfflineDeviceStatus {
	return []OfflineDeviceStatus{
		DeviceStatusActive,
		DeviceStatusInactive,
		DeviceStatusRevoked,
		DeviceStatusSuspended,
	}
}