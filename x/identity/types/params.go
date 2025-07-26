package types

import (
	"fmt"
	"time"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys
var (
	KeyMaxDIDDocumentSize        = []byte("MaxDIDDocumentSize")
	KeyMaxCredentialSize         = []byte("MaxCredentialSize")
	KeyMaxProofSize              = []byte("MaxProofSize")
	KeyCredentialExpiryDays      = []byte("CredentialExpiryDays")
	KeyKYCExpiryDays             = []byte("KYCExpiryDays")
	KeyBiometricExpiryDays       = []byte("BiometricExpiryDays")
	KeyMaxRecoveryMethods        = []byte("MaxRecoveryMethods")
	KeyMaxCredentialsPerIdentity = []byte("MaxCredentialsPerIdentity")
	KeyMinAnonymitySetSize       = []byte("MinAnonymitySetSize")
	KeyProofExpiryMinutes        = []byte("ProofExpiryMinutes")
	KeyMaxConsentDurationDays    = []byte("MaxConsentDurationDays")
	KeyEnableAnonymousCredentials = []byte("EnableAnonymousCredentials")
	KeyEnableZKProofs            = []byte("EnableZKProofs")
	KeyEnableIndiaStack          = []byte("EnableIndiaStack")
	KeyRequireKYCForHighValue    = []byte("RequireKYCForHighValue")
	KeyHighValueThreshold        = []byte("HighValueThreshold")
	KeyMaxFailedAuthAttempts     = []byte("MaxFailedAuthAttempts")
	KeyAuthLockoutDurationHours  = []byte("AuthLockoutDurationHours")
	KeySupportedDIDMethods       = []byte("SupportedDIDMethods")
	KeySupportedProofSystems     = []byte("SupportedProofSystems")
	KeyTrustedIssuers            = []byte("TrustedIssuers")
)

// Params defines the parameters for the identity module
type Params struct {
	// Size limits
	MaxDIDDocumentSize        uint64 `json:"max_did_document_size"`
	MaxCredentialSize         uint64 `json:"max_credential_size"`
	MaxProofSize              uint64 `json:"max_proof_size"`
	
	// Expiry settings
	CredentialExpiryDays      uint32 `json:"credential_expiry_days"`
	KYCExpiryDays             uint32 `json:"kyc_expiry_days"`
	BiometricExpiryDays       uint32 `json:"biometric_expiry_days"`
	ProofExpiryMinutes        uint32 `json:"proof_expiry_minutes"`
	MaxConsentDurationDays    uint32 `json:"max_consent_duration_days"`
	
	// Limits
	MaxRecoveryMethods        uint32 `json:"max_recovery_methods"`
	MaxCredentialsPerIdentity uint32 `json:"max_credentials_per_identity"`
	MinAnonymitySetSize       uint32 `json:"min_anonymity_set_size"`
	
	// Feature flags
	EnableAnonymousCredentials bool   `json:"enable_anonymous_credentials"`
	EnableZKProofs            bool   `json:"enable_zk_proofs"`
	EnableIndiaStack          bool   `json:"enable_india_stack"`
	
	// Security settings
	RequireKYCForHighValue    bool   `json:"require_kyc_for_high_value"`
	HighValueThreshold        int64  `json:"high_value_threshold"`
	MaxFailedAuthAttempts     uint32 `json:"max_failed_auth_attempts"`
	AuthLockoutDurationHours  uint32 `json:"auth_lockout_duration_hours"`
	
	// Supported methods
	SupportedDIDMethods       []string `json:"supported_did_methods"`
	SupportedProofSystems     []string `json:"supported_proof_systems"`
	TrustedIssuers            []string `json:"trusted_issuers"`
}

// NewParams creates a new Params instance
func NewParams(
	maxDIDDocumentSize uint64,
	maxCredentialSize uint64,
	maxProofSize uint64,
	credentialExpiryDays uint32,
	kycExpiryDays uint32,
	biometricExpiryDays uint32,
	maxRecoveryMethods uint32,
	maxCredentialsPerIdentity uint32,
	minAnonymitySetSize uint32,
	proofExpiryMinutes uint32,
	maxConsentDurationDays uint32,
	enableAnonymousCredentials bool,
	enableZKProofs bool,
	enableIndiaStack bool,
	requireKYCForHighValue bool,
	highValueThreshold int64,
	maxFailedAuthAttempts uint32,
	authLockoutDurationHours uint32,
	supportedDIDMethods []string,
	supportedProofSystems []string,
	trustedIssuers []string,
) Params {
	return Params{
		MaxDIDDocumentSize:        maxDIDDocumentSize,
		MaxCredentialSize:         maxCredentialSize,
		MaxProofSize:              maxProofSize,
		CredentialExpiryDays:      credentialExpiryDays,
		KYCExpiryDays:             kycExpiryDays,
		BiometricExpiryDays:       biometricExpiryDays,
		MaxRecoveryMethods:        maxRecoveryMethods,
		MaxCredentialsPerIdentity: maxCredentialsPerIdentity,
		MinAnonymitySetSize:       minAnonymitySetSize,
		ProofExpiryMinutes:        proofExpiryMinutes,
		MaxConsentDurationDays:    maxConsentDurationDays,
		EnableAnonymousCredentials: enableAnonymousCredentials,
		EnableZKProofs:            enableZKProofs,
		EnableIndiaStack:          enableIndiaStack,
		RequireKYCForHighValue:    requireKYCForHighValue,
		HighValueThreshold:        highValueThreshold,
		MaxFailedAuthAttempts:     maxFailedAuthAttempts,
		AuthLockoutDurationHours:  authLockoutDurationHours,
		SupportedDIDMethods:       supportedDIDMethods,
		SupportedProofSystems:     supportedProofSystems,
		TrustedIssuers:            trustedIssuers,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return Params{
		// Size limits
		MaxDIDDocumentSize:        65536,  // 64 KB
		MaxCredentialSize:         32768,  // 32 KB
		MaxProofSize:              16384,  // 16 KB
		
		// Expiry settings
		CredentialExpiryDays:      365,    // 1 year
		KYCExpiryDays:             180,    // 6 months
		BiometricExpiryDays:       730,    // 2 years
		ProofExpiryMinutes:        60,     // 1 hour
		MaxConsentDurationDays:    365,    // 1 year
		
		// Limits
		MaxRecoveryMethods:        5,
		MaxCredentialsPerIdentity: 100,
		MinAnonymitySetSize:       10,
		
		// Feature flags
		EnableAnonymousCredentials: true,
		EnableZKProofs:            true,
		EnableIndiaStack:          true,
		
		// Security settings
		RequireKYCForHighValue:    true,
		HighValueThreshold:        100000000000, // 100,000 NAMO (in unamo)
		MaxFailedAuthAttempts:     5,
		AuthLockoutDurationHours:  24,
		
		// Supported methods
		SupportedDIDMethods: []string{
			"desh",
			"web",
		},
		SupportedProofSystems: []string{
			ProofSystemGroth16,
			ProofSystemPLONK,
			ProofSystemBulletproofs,
		},
		TrustedIssuers: []string{
			"did:desh:kyc-issuer",
			"did:desh:biometric-issuer",
			"did:desh:government-issuer",
		},
	}
}

// ParamKeyTable returns the key table for identity module parameters
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs returns the parameter set pairs
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMaxDIDDocumentSize, &p.MaxDIDDocumentSize, validateMaxDIDDocumentSize),
		paramtypes.NewParamSetPair(KeyMaxCredentialSize, &p.MaxCredentialSize, validateMaxCredentialSize),
		paramtypes.NewParamSetPair(KeyMaxProofSize, &p.MaxProofSize, validateMaxProofSize),
		paramtypes.NewParamSetPair(KeyCredentialExpiryDays, &p.CredentialExpiryDays, validateCredentialExpiryDays),
		paramtypes.NewParamSetPair(KeyKYCExpiryDays, &p.KYCExpiryDays, validateKYCExpiryDays),
		paramtypes.NewParamSetPair(KeyBiometricExpiryDays, &p.BiometricExpiryDays, validateBiometricExpiryDays),
		paramtypes.NewParamSetPair(KeyMaxRecoveryMethods, &p.MaxRecoveryMethods, validateMaxRecoveryMethods),
		paramtypes.NewParamSetPair(KeyMaxCredentialsPerIdentity, &p.MaxCredentialsPerIdentity, validateMaxCredentialsPerIdentity),
		paramtypes.NewParamSetPair(KeyMinAnonymitySetSize, &p.MinAnonymitySetSize, validateMinAnonymitySetSize),
		paramtypes.NewParamSetPair(KeyProofExpiryMinutes, &p.ProofExpiryMinutes, validateProofExpiryMinutes),
		paramtypes.NewParamSetPair(KeyMaxConsentDurationDays, &p.MaxConsentDurationDays, validateMaxConsentDurationDays),
		paramtypes.NewParamSetPair(KeyEnableAnonymousCredentials, &p.EnableAnonymousCredentials, validateBool),
		paramtypes.NewParamSetPair(KeyEnableZKProofs, &p.EnableZKProofs, validateBool),
		paramtypes.NewParamSetPair(KeyEnableIndiaStack, &p.EnableIndiaStack, validateBool),
		paramtypes.NewParamSetPair(KeyRequireKYCForHighValue, &p.RequireKYCForHighValue, validateBool),
		paramtypes.NewParamSetPair(KeyHighValueThreshold, &p.HighValueThreshold, validateHighValueThreshold),
		paramtypes.NewParamSetPair(KeyMaxFailedAuthAttempts, &p.MaxFailedAuthAttempts, validateMaxFailedAuthAttempts),
		paramtypes.NewParamSetPair(KeyAuthLockoutDurationHours, &p.AuthLockoutDurationHours, validateAuthLockoutDurationHours),
		paramtypes.NewParamSetPair(KeySupportedDIDMethods, &p.SupportedDIDMethods, validateSupportedDIDMethods),
		paramtypes.NewParamSetPair(KeySupportedProofSystems, &p.SupportedProofSystems, validateSupportedProofSystems),
		paramtypes.NewParamSetPair(KeyTrustedIssuers, &p.TrustedIssuers, validateTrustedIssuers),
	}
}

// Validate validates the parameter set
func (p Params) Validate() error {
	if err := validateMaxDIDDocumentSize(p.MaxDIDDocumentSize); err != nil {
		return err
	}
	if err := validateMaxCredentialSize(p.MaxCredentialSize); err != nil {
		return err
	}
	if err := validateMaxProofSize(p.MaxProofSize); err != nil {
		return err
	}
	if err := validateCredentialExpiryDays(p.CredentialExpiryDays); err != nil {
		return err
	}
	if err := validateKYCExpiryDays(p.KYCExpiryDays); err != nil {
		return err
	}
	if err := validateBiometricExpiryDays(p.BiometricExpiryDays); err != nil {
		return err
	}
	if err := validateMaxRecoveryMethods(p.MaxRecoveryMethods); err != nil {
		return err
	}
	if err := validateMaxCredentialsPerIdentity(p.MaxCredentialsPerIdentity); err != nil {
		return err
	}
	if err := validateMinAnonymitySetSize(p.MinAnonymitySetSize); err != nil {
		return err
	}
	if err := validateProofExpiryMinutes(p.ProofExpiryMinutes); err != nil {
		return err
	}
	if err := validateMaxConsentDurationDays(p.MaxConsentDurationDays); err != nil {
		return err
	}
	if err := validateHighValueThreshold(p.HighValueThreshold); err != nil {
		return err
	}
	if err := validateMaxFailedAuthAttempts(p.MaxFailedAuthAttempts); err != nil {
		return err
	}
	if err := validateAuthLockoutDurationHours(p.AuthLockoutDurationHours); err != nil {
		return err
	}
	if err := validateSupportedDIDMethods(p.SupportedDIDMethods); err != nil {
		return err
	}
	if err := validateSupportedProofSystems(p.SupportedProofSystems); err != nil {
		return err
	}
	if err := validateTrustedIssuers(p.TrustedIssuers); err != nil {
		return err
	}
	
	return nil
}

// Validation functions

func validateMaxDIDDocumentSize(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	if v < 1024 || v > 1048576 { // 1KB to 1MB
		return fmt.Errorf("max DID document size must be between 1KB and 1MB")
	}
	
	return nil
}

func validateMaxCredentialSize(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	if v < 512 || v > 524288 { // 512B to 512KB
		return fmt.Errorf("max credential size must be between 512B and 512KB")
	}
	
	return nil
}

func validateMaxProofSize(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	if v < 256 || v > 262144 { // 256B to 256KB
		return fmt.Errorf("max proof size must be between 256B and 256KB")
	}
	
	return nil
}

func validateCredentialExpiryDays(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	if v < 30 || v > 3650 { // 30 days to 10 years
		return fmt.Errorf("credential expiry must be between 30 days and 10 years")
	}
	
	return nil
}

func validateKYCExpiryDays(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	if v < 30 || v > 730 { // 30 days to 2 years
		return fmt.Errorf("KYC expiry must be between 30 days and 2 years")
	}
	
	return nil
}

func validateBiometricExpiryDays(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	if v < 180 || v > 3650 { // 6 months to 10 years
		return fmt.Errorf("biometric expiry must be between 6 months and 10 years")
	}
	
	return nil
}

func validateMaxRecoveryMethods(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	if v < 1 || v > 10 {
		return fmt.Errorf("max recovery methods must be between 1 and 10")
	}
	
	return nil
}

func validateMaxCredentialsPerIdentity(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	if v < 10 || v > 1000 {
		return fmt.Errorf("max credentials per identity must be between 10 and 1000")
	}
	
	return nil
}

func validateMinAnonymitySetSize(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	if v < 5 || v > 100 {
		return fmt.Errorf("min anonymity set size must be between 5 and 100")
	}
	
	return nil
}

func validateProofExpiryMinutes(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	if v < 5 || v > 1440 { // 5 minutes to 24 hours
		return fmt.Errorf("proof expiry must be between 5 minutes and 24 hours")
	}
	
	return nil
}

func validateMaxConsentDurationDays(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	if v < 1 || v > 3650 { // 1 day to 10 years
		return fmt.Errorf("max consent duration must be between 1 day and 10 years")
	}
	
	return nil
}

func validateBool(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateHighValueThreshold(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	if v < 0 {
		return fmt.Errorf("high value threshold must be non-negative")
	}
	
	return nil
}

func validateMaxFailedAuthAttempts(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	if v < 3 || v > 10 {
		return fmt.Errorf("max failed auth attempts must be between 3 and 10")
	}
	
	return nil
}

func validateAuthLockoutDurationHours(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	if v < 1 || v > 168 { // 1 hour to 1 week
		return fmt.Errorf("auth lockout duration must be between 1 hour and 1 week")
	}
	
	return nil
}

func validateSupportedDIDMethods(i interface{}) error {
	v, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	if len(v) == 0 {
		return fmt.Errorf("at least one DID method must be supported")
	}
	
	// Check for duplicates
	methodMap := make(map[string]bool)
	for _, method := range v {
		if method == "" {
			return fmt.Errorf("DID method cannot be empty")
		}
		if methodMap[method] {
			return fmt.Errorf("duplicate DID method: %s", method)
		}
		methodMap[method] = true
	}
	
	return nil
}

func validateSupportedProofSystems(i interface{}) error {
	v, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	if len(v) == 0 {
		return fmt.Errorf("at least one proof system must be supported")
	}
	
	// Check for duplicates
	systemMap := make(map[string]bool)
	for _, system := range v {
		if system == "" {
			return fmt.Errorf("proof system cannot be empty")
		}
		if systemMap[system] {
			return fmt.Errorf("duplicate proof system: %s", system)
		}
		systemMap[system] = true
	}
	
	return nil
}

func validateTrustedIssuers(i interface{}) error {
	v, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	// Trusted issuers can be empty (trust no one by default)
	
	// Check for duplicates
	issuerMap := make(map[string]bool)
	for _, issuer := range v {
		if issuer == "" {
			return fmt.Errorf("trusted issuer cannot be empty")
		}
		// Basic DID format check
		if len(issuer) < 7 || issuer[:4] != "did:" {
			return fmt.Errorf("invalid DID format for trusted issuer: %s", issuer)
		}
		if issuerMap[issuer] {
			return fmt.Errorf("duplicate trusted issuer: %s", issuer)
		}
		issuerMap[issuer] = true
	}
	
	return nil
}

// GetCredentialExpiry returns the credential expiry duration
func (p Params) GetCredentialExpiry() time.Duration {
	return time.Duration(p.CredentialExpiryDays) * 24 * time.Hour
}

// GetKYCExpiry returns the KYC expiry duration
func (p Params) GetKYCExpiry() time.Duration {
	return time.Duration(p.KYCExpiryDays) * 24 * time.Hour
}

// GetBiometricExpiry returns the biometric expiry duration
func (p Params) GetBiometricExpiry() time.Duration {
	return time.Duration(p.BiometricExpiryDays) * 24 * time.Hour
}

// GetProofExpiry returns the proof expiry duration
func (p Params) GetProofExpiry() time.Duration {
	return time.Duration(p.ProofExpiryMinutes) * time.Minute
}

// GetMaxConsentDuration returns the max consent duration
func (p Params) GetMaxConsentDuration() time.Duration {
	return time.Duration(p.MaxConsentDurationDays) * 24 * time.Hour
}

// GetAuthLockoutDuration returns the auth lockout duration
func (p Params) GetAuthLockoutDuration() time.Duration {
	return time.Duration(p.AuthLockoutDurationHours) * time.Hour
}

// IsDIDMethodSupported checks if a DID method is supported
func (p Params) IsDIDMethodSupported(method string) bool {
	for _, m := range p.SupportedDIDMethods {
		if m == method {
			return true
		}
	}
	return false
}

// IsProofSystemSupported checks if a proof system is supported
func (p Params) IsProofSystemSupported(system string) bool {
	for _, s := range p.SupportedProofSystems {
		if s == system {
			return true
		}
	}
	return false
}

// IsTrustedIssuer checks if an issuer is trusted
func (p Params) IsTrustedIssuer(issuer string) bool {
	for _, i := range p.TrustedIssuers {
		if i == issuer {
			return true
		}
	}
	return false
}