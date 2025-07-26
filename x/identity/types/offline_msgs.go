package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Message constants for offline verification
const (
	TypeMsgPrepareOfflineVerification = "prepare_offline_verification"
	TypeMsgCreateOfflineBackup        = "create_offline_backup"
	TypeMsgUpdateOfflineConfig        = "update_offline_config"
	TypeMsgRegisterOfflineDevice      = "register_offline_device"
	TypeMsgRevokeOfflineAccess        = "revoke_offline_access"
)

var (
	_ sdk.Msg = &MsgPrepareOfflineVerification{}
	_ sdk.Msg = &MsgCreateOfflineBackup{}
	_ sdk.Msg = &MsgUpdateOfflineConfig{}
	_ sdk.Msg = &MsgRegisterOfflineDevice{}
	_ sdk.Msg = &MsgRevokeOfflineAccess{}
)

// MsgPrepareOfflineVerification prepares offline verification data for an identity
type MsgPrepareOfflineVerification struct {
	DID              string                     `json:"did" yaml:"did"`
	Format           OfflineCredentialFormat    `json:"format" yaml:"format"`
	ExpirationHours  uint32                     `json:"expiration_hours" yaml:"expiration_hours"`
	IncludeBiometric bool                       `json:"include_biometric" yaml:"include_biometric"`
	IncludeCredentials []string                 `json:"include_credentials" yaml:"include_credentials"`
	RequiredLevel    uint32                     `json:"required_level" yaml:"required_level"`
	DeviceID         string                     `json:"device_id,omitempty" yaml:"device_id,omitempty"`
	Signer           string                     `json:"signer" yaml:"signer"`
}

// NewMsgPrepareOfflineVerification creates a new MsgPrepareOfflineVerification instance
func NewMsgPrepareOfflineVerification(did string, format OfflineCredentialFormat, expirationHours uint32, signer string) *MsgPrepareOfflineVerification {
	return &MsgPrepareOfflineVerification{
		DID:             did,
		Format:          format,
		ExpirationHours: expirationHours,
		RequiredLevel:   2, // Default to KYC level
		Signer:          signer,
	}
}

// Route returns the module route
func (msg *MsgPrepareOfflineVerification) Route() string {
	return RouterKey
}

// Type returns the message type
func (msg *MsgPrepareOfflineVerification) Type() string {
	return TypeMsgPrepareOfflineVerification
}

// GetSigners returns the signers
func (msg *MsgPrepareOfflineVerification) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

// GetSignBytes returns the sign bytes
func (msg *MsgPrepareOfflineVerification) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validates the message
func (msg *MsgPrepareOfflineVerification) ValidateBasic() error {
	if err := ValidateUserDID(msg.DID); err != nil {
		return sdkerrors.Wrap(ErrInvalidRequest, err.Error())
	}

	if msg.Format == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "format cannot be empty")
	}

	validFormats := []OfflineCredentialFormat{
		FormatSelfContained,
		FormatCompressed,
		FormatQRCode,
		FormatNFC,
		FormatPrintable,
	}

	isValidFormat := false
	for _, format := range validFormats {
		if msg.Format == format {
			isValidFormat = true
			break
		}
	}

	if !isValidFormat {
		return sdkerrors.Wrapf(ErrInvalidRequest, "invalid format: %s", msg.Format)
	}

	if msg.ExpirationHours == 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "expiration hours must be greater than 0")
	}

	if msg.ExpirationHours > 8760 { // Max 1 year
		return sdkerrors.Wrap(ErrInvalidRequest, "expiration hours cannot exceed 8760 (1 year)")
	}

	if msg.RequiredLevel < 1 || msg.RequiredLevel > 5 {
		return sdkerrors.Wrap(ErrInvalidRequest, "required level must be between 1 and 5")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid signer address: %v", err)
	}

	return nil
}

// MsgCreateOfflineBackup creates a comprehensive offline backup
type MsgCreateOfflineBackup struct {
	DID                string `json:"did" yaml:"did"`
	IncludePrivateData bool   `json:"include_private_data" yaml:"include_private_data"`
	BackupPassword     string `json:"backup_password,omitempty" yaml:"backup_password,omitempty"`
	BackupLocation     string `json:"backup_location,omitempty" yaml:"backup_location,omitempty"`
	Signer             string `json:"signer" yaml:"signer"`
}

// NewMsgCreateOfflineBackup creates a new MsgCreateOfflineBackup instance
func NewMsgCreateOfflineBackup(did string, includePrivateData bool, signer string) *MsgCreateOfflineBackup {
	return &MsgCreateOfflineBackup{
		DID:                did,
		IncludePrivateData: includePrivateData,
		Signer:             signer,
	}
}

// Route returns the module route
func (msg *MsgCreateOfflineBackup) Route() string {
	return RouterKey
}

// Type returns the message type
func (msg *MsgCreateOfflineBackup) Type() string {
	return TypeMsgCreateOfflineBackup
}

// GetSigners returns the signers
func (msg *MsgCreateOfflineBackup) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

// GetSignBytes returns the sign bytes
func (msg *MsgCreateOfflineBackup) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validates the message
func (msg *MsgCreateOfflineBackup) ValidateBasic() error {
	if err := ValidateUserDID(msg.DID); err != nil {
		return sdkerrors.Wrap(ErrInvalidRequest, err.Error())
	}

	if msg.IncludePrivateData && msg.BackupPassword == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "backup password required when including private data")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid signer address: %v", err)
	}

	return nil
}

// MsgUpdateOfflineConfig updates offline verification configuration
type MsgUpdateOfflineConfig struct {
	Config *OfflineVerificationConfig `json:"config" yaml:"config"`
	Signer string                     `json:"signer" yaml:"signer"`
}

// NewMsgUpdateOfflineConfig creates a new MsgUpdateOfflineConfig instance
func NewMsgUpdateOfflineConfig(config *OfflineVerificationConfig, signer string) *MsgUpdateOfflineConfig {
	return &MsgUpdateOfflineConfig{
		Config: config,
		Signer: signer,
	}
}

// Route returns the module route
func (msg *MsgUpdateOfflineConfig) Route() string {
	return RouterKey
}

// Type returns the message type
func (msg *MsgUpdateOfflineConfig) Type() string {
	return TypeMsgUpdateOfflineConfig
}

// GetSigners returns the signers
func (msg *MsgUpdateOfflineConfig) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

// GetSignBytes returns the sign bytes
func (msg *MsgUpdateOfflineConfig) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validates the message
func (msg *MsgUpdateOfflineConfig) ValidateBasic() error {
	if msg.Config == nil {
		return sdkerrors.Wrap(ErrInvalidRequest, "config cannot be nil")
	}

	// Validate configuration values
	if msg.Config.MaxOfflineDuration <= 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "max offline duration must be positive")
	}

	if msg.Config.RequiredConfidence < 0 || msg.Config.RequiredConfidence > 1 {
		return sdkerrors.Wrap(ErrInvalidRequest, "required confidence must be between 0 and 1")
	}

	if msg.Config.BiometricThreshold < 0 || msg.Config.BiometricThreshold > 1 {
		return sdkerrors.Wrap(ErrInvalidRequest, "biometric threshold must be between 0 and 1")
	}

	if msg.Config.MaxCacheSize == 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "max cache size must be positive")
	}

	if msg.Config.CacheExpirationPeriod <= 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "cache expiration period must be positive")
	}

	if msg.Config.EmergencyThreshold < 0 || msg.Config.EmergencyThreshold > 1 {
		return sdkerrors.Wrap(ErrInvalidRequest, "emergency threshold must be between 0 and 1")
	}

	if msg.Config.DefaultLanguage == "" {
		return sdkerrors.Wrap(ErrInvalidLanguageCode, "default language cannot be empty")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid signer address: %v", err)
	}

	return nil
}

// MsgRegisterOfflineDevice registers a device for offline verification
type MsgRegisterOfflineDevice struct {
	DID                string                  `json:"did" yaml:"did"`
	DeviceID           string                  `json:"device_id" yaml:"device_id"`
	DeviceName         string                  `json:"device_name" yaml:"device_name"`
	DeviceType         string                  `json:"device_type" yaml:"device_type"`
	PublicKey          string                  `json:"public_key" yaml:"public_key"`
	Capabilities       []string                `json:"capabilities" yaml:"capabilities"`
	SecurityLevel      uint32                  `json:"security_level" yaml:"security_level"`
	MaxOfflineDuration uint32                  `json:"max_offline_duration" yaml:"max_offline_duration"` // hours
	Signer             string                  `json:"signer" yaml:"signer"`
}

// NewMsgRegisterOfflineDevice creates a new MsgRegisterOfflineDevice instance
func NewMsgRegisterOfflineDevice(did, deviceID, deviceName, deviceType, publicKey string, signer string) *MsgRegisterOfflineDevice {
	return &MsgRegisterOfflineDevice{
		DID:                did,
		DeviceID:           deviceID,
		DeviceName:         deviceName,
		DeviceType:         deviceType,
		PublicKey:          publicKey,
		Capabilities:       []string{"identity_verification"},
		SecurityLevel:      2,
		MaxOfflineDuration: 24, // 24 hours default
		Signer:             signer,
	}
}

// Route returns the module route
func (msg *MsgRegisterOfflineDevice) Route() string {
	return RouterKey
}

// Type returns the message type
func (msg *MsgRegisterOfflineDevice) Type() string {
	return TypeMsgRegisterOfflineDevice
}

// GetSigners returns the signers
func (msg *MsgRegisterOfflineDevice) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

// GetSignBytes returns the sign bytes
func (msg *MsgRegisterOfflineDevice) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validates the message
func (msg *MsgRegisterOfflineDevice) ValidateBasic() error {
	if err := ValidateUserDID(msg.DID); err != nil {
		return sdkerrors.Wrap(ErrInvalidRequest, err.Error())
	}

	if msg.DeviceID == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "device ID cannot be empty")
	}

	if msg.DeviceName == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "device name cannot be empty")
	}

	if msg.DeviceType == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "device type cannot be empty")
	}

	validDeviceTypes := []string{"mobile", "tablet", "laptop", "desktop", "iot", "smartcard", "hardware_token"}
	isValidType := false
	for _, validType := range validDeviceTypes {
		if msg.DeviceType == validType {
			isValidType = true
			break
		}
	}

	if !isValidType {
		return sdkerrors.Wrapf(ErrInvalidRequest, "invalid device type: %s", msg.DeviceType)
	}

	if msg.PublicKey == "" {
		return sdkerrors.Wrap(ErrInvalidPublicKey, "public key cannot be empty")
	}

	if len(msg.Capabilities) == 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "capabilities cannot be empty")
	}

	validCapabilities := []string{
		"identity_verification",
		"biometric_capture",
		"credential_storage",
		"offline_signing",
		"secure_element",
	}

	for _, capability := range msg.Capabilities {
		isValidCapability := false
		for _, validCap := range validCapabilities {
			if capability == validCap {
				isValidCapability = true
				break
			}
		}
		if !isValidCapability {
			return sdkerrors.Wrapf(ErrInvalidRequest, "invalid capability: %s", capability)
		}
	}

	if msg.SecurityLevel < 1 || msg.SecurityLevel > 5 {
		return sdkerrors.Wrap(ErrInvalidRequest, "security level must be between 1 and 5")
	}

	if msg.MaxOfflineDuration == 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "max offline duration must be greater than 0")
	}

	if msg.MaxOfflineDuration > 8760 { // Max 1 year
		return sdkerrors.Wrap(ErrInvalidRequest, "max offline duration cannot exceed 8760 hours (1 year)")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid signer address: %v", err)
	}

	return nil
}

// MsgRevokeOfflineAccess revokes offline access for a device or all devices
type MsgRevokeOfflineAccess struct {
	DID      string `json:"did" yaml:"did"`
	DeviceID string `json:"device_id,omitempty" yaml:"device_id,omitempty"` // Empty means revoke all devices
	Reason   string `json:"reason,omitempty" yaml:"reason,omitempty"`
	Signer   string `json:"signer" yaml:"signer"`
}

// NewMsgRevokeOfflineAccess creates a new MsgRevokeOfflineAccess instance
func NewMsgRevokeOfflineAccess(did, deviceID, reason, signer string) *MsgRevokeOfflineAccess {
	return &MsgRevokeOfflineAccess{
		DID:      did,
		DeviceID: deviceID,
		Reason:   reason,
		Signer:   signer,
	}
}

// Route returns the module route
func (msg *MsgRevokeOfflineAccess) Route() string {
	return RouterKey
}

// Type returns the message type
func (msg *MsgRevokeOfflineAccess) Type() string {
	return TypeMsgRevokeOfflineAccess
}

// GetSigners returns the signers
func (msg *MsgRevokeOfflineAccess) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

// GetSignBytes returns the sign bytes
func (msg *MsgRevokeOfflineAccess) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validates the message
func (msg *MsgRevokeOfflineAccess) ValidateBasic() error {
	if err := ValidateUserDID(msg.DID); err != nil {
		return sdkerrors.Wrap(ErrInvalidRequest, err.Error())
	}

	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid signer address: %v", err)
	}

	return nil
}

// Helper functions for offline verification

// ValidateOfflineFormat validates offline credential format
func ValidateOfflineFormat(format OfflineCredentialFormat) error {
	validFormats := []OfflineCredentialFormat{
		FormatSelfContained,
		FormatCompressed,
		FormatQRCode,
		FormatNFC,
		FormatPrintable,
	}

	for _, validFormat := range validFormats {
		if format == validFormat {
			return nil
		}
	}

	return fmt.Errorf("invalid offline credential format: %s", format)
}

// ValidateOfflineMode validates offline verification mode
func ValidateOfflineMode(mode OfflineVerificationMode) error {
	validModes := []OfflineVerificationMode{
		OfflineModeFull,
		OfflineModePartial,
		OfflineModeMinimal,
		OfflineModeEmergency,
	}

	for _, validMode := range validModes {
		if mode == validMode {
			return nil
		}
	}

	return fmt.Errorf("invalid offline verification mode: %s", mode)
}

// GetSupportedOfflineFormats returns all supported offline credential formats
func GetSupportedOfflineFormats() []OfflineCredentialFormat {
	return []OfflineCredentialFormat{
		FormatSelfContained,
		FormatCompressed,
		FormatQRCode,
		FormatNFC,
		FormatPrintable,
	}
}

// GetSupportedOfflineModes returns all supported offline verification modes
func GetSupportedOfflineModes() []OfflineVerificationMode {
	return []OfflineVerificationMode{
		OfflineModeFull,
		OfflineModePartial,
		OfflineModeMinimal,
		OfflineModeEmergency,
	}
}

// GetSupportedDeviceTypes returns all supported device types for offline verification
func GetSupportedDeviceTypes() []string {
	return []string{
		"mobile",
		"tablet",
		"laptop",
		"desktop",
		"iot",
		"smartcard",
		"hardware_token",
	}
}

// GetSupportedDeviceCapabilities returns all supported device capabilities
func GetSupportedDeviceCapabilities() []string {
	return []string{
		"identity_verification",
		"biometric_capture",
		"credential_storage",
		"offline_signing",
		"secure_element",
	}
}