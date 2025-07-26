package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Identity module error codes
var (
	// Identity errors
	ErrIdentityNotFound       = sdkerrors.Register(ModuleName, 100, "identity not found")
	ErrIdentityAlreadyExists  = sdkerrors.Register(ModuleName, 101, "identity already exists")
	ErrIdentityInactive       = sdkerrors.Register(ModuleName, 102, "identity is not active")
	ErrIdentityRevoked        = sdkerrors.Register(ModuleName, 103, "identity has been revoked")
	ErrInvalidIdentity        = sdkerrors.Register(ModuleName, 104, "invalid identity")
	
	// DID errors
	ErrDIDNotFound            = sdkerrors.Register(ModuleName, 200, "DID not found")
	ErrDIDAlreadyExists       = sdkerrors.Register(ModuleName, 201, "DID already exists")
	ErrInvalidDID             = sdkerrors.Register(ModuleName, 202, "invalid DID format")
	ErrDIDDeactivated         = sdkerrors.Register(ModuleName, 203, "DID is deactivated")
	ErrInvalidDIDDocument     = sdkerrors.Register(ModuleName, 204, "invalid DID document")
	ErrDIDMethodNotSupported  = sdkerrors.Register(ModuleName, 205, "DID method not supported")
	
	// Credential errors
	ErrCredentialNotFound     = sdkerrors.Register(ModuleName, 300, "credential not found")
	ErrCredentialExpired      = sdkerrors.Register(ModuleName, 301, "credential has expired")
	ErrCredentialRevoked      = sdkerrors.Register(ModuleName, 302, "credential has been revoked")
	ErrInvalidCredential      = sdkerrors.Register(ModuleName, 303, "invalid credential")
	ErrInvalidProof           = sdkerrors.Register(ModuleName, 304, "invalid proof")
	ErrInvalidIssuer          = sdkerrors.Register(ModuleName, 305, "invalid or unauthorized issuer")
	ErrInvalidHolder          = sdkerrors.Register(ModuleName, 306, "invalid holder")
	ErrCredentialNotVerified  = sdkerrors.Register(ModuleName, 307, "credential verification failed")
	ErrSchemaNotFound         = sdkerrors.Register(ModuleName, 308, "credential schema not found")
	ErrInvalidSchema          = sdkerrors.Register(ModuleName, 309, "invalid credential schema")
	
	// KYC errors
	ErrKYCNotCompleted        = sdkerrors.Register(ModuleName, 400, "KYC not completed")
	ErrKYCExpired             = sdkerrors.Register(ModuleName, 401, "KYC has expired")
	ErrKYCLevelInsufficient   = sdkerrors.Register(ModuleName, 402, "KYC level insufficient")
	ErrKYCVerificationFailed  = sdkerrors.Register(ModuleName, 403, "KYC verification failed")
	ErrKYCDataMismatch        = sdkerrors.Register(ModuleName, 404, "KYC data mismatch")
	
	// Biometric errors
	ErrBiometricNotEnrolled   = sdkerrors.Register(ModuleName, 500, "biometric not enrolled")
	ErrBiometricAuthFailed    = sdkerrors.Register(ModuleName, 501, "biometric authentication failed")
	ErrBiometricExpired       = sdkerrors.Register(ModuleName, 502, "biometric enrollment expired")
	ErrBiometricQualityLow    = sdkerrors.Register(ModuleName, 503, "biometric quality too low")
	ErrBiometricLocked        = sdkerrors.Register(ModuleName, 504, "biometric locked due to failed attempts")
	ErrInvalidBiometricType   = sdkerrors.Register(ModuleName, 505, "invalid biometric type")
	
	// Consent errors
	ErrConsentNotGiven        = sdkerrors.Register(ModuleName, 600, "consent not given")
	ErrConsentExpired         = sdkerrors.Register(ModuleName, 601, "consent has expired")
	ErrConsentWithdrawn       = sdkerrors.Register(ModuleName, 602, "consent has been withdrawn")
	ErrInvalidConsentType     = sdkerrors.Register(ModuleName, 603, "invalid consent type")
	ErrConsentAlreadyExists   = sdkerrors.Register(ModuleName, 604, "consent already exists")
	
	// Privacy errors
	ErrInvalidZKProof         = sdkerrors.Register(ModuleName, 700, "invalid zero-knowledge proof")
	ErrProofExpired           = sdkerrors.Register(ModuleName, 701, "proof has expired")
	ErrNullifierUsed          = sdkerrors.Register(ModuleName, 702, "nullifier already used")
	ErrAnonymitySetTooSmall   = sdkerrors.Register(ModuleName, 703, "anonymity set too small")
	ErrPrivacyViolation       = sdkerrors.Register(ModuleName, 704, "privacy settings violation")
	ErrSelectiveDisclosureFailed = sdkerrors.Register(ModuleName, 705, "selective disclosure failed")
	
	// India Stack errors
	ErrAadhaarNotLinked       = sdkerrors.Register(ModuleName, 800, "Aadhaar not linked")
	ErrAadhaarVerificationFailed = sdkerrors.Register(ModuleName, 801, "Aadhaar verification failed")
	ErrDigiLockerNotConnected = sdkerrors.Register(ModuleName, 802, "DigiLocker not connected")
	ErrDocumentNotFound       = sdkerrors.Register(ModuleName, 803, "document not found in DigiLocker")
	ErrUPINotLinked           = sdkerrors.Register(ModuleName, 804, "UPI not linked")
	ErrDEPAConsentInvalid     = sdkerrors.Register(ModuleName, 805, "DEPA consent invalid")
	ErrAAConsentExpired       = sdkerrors.Register(ModuleName, 806, "Account Aggregator consent expired")
	ErrPanchayatKYCInvalid    = sdkerrors.Register(ModuleName, 807, "Village Panchayat KYC invalid")
	
	// Recovery errors
	ErrRecoveryMethodNotSet   = sdkerrors.Register(ModuleName, 900, "recovery method not set")
	ErrRecoveryFailed         = sdkerrors.Register(ModuleName, 901, "recovery failed")
	ErrInvalidRecoveryMethod  = sdkerrors.Register(ModuleName, 902, "invalid recovery method")
	ErrRecoveryLimitExceeded  = sdkerrors.Register(ModuleName, 903, "recovery attempt limit exceeded")
	ErrGuardianNotFound       = sdkerrors.Register(ModuleName, 904, "guardian not found")
	
	// Authorization errors
	ErrUnauthorized           = sdkerrors.Register(ModuleName, 1000, "unauthorized")
	ErrInsufficientPermission = sdkerrors.Register(ModuleName, 1001, "insufficient permission")
	ErrDelegationNotFound     = sdkerrors.Register(ModuleName, 1002, "delegation not found")
	ErrDelegationExpired      = sdkerrors.Register(ModuleName, 1003, "delegation expired")
	
	// Service errors
	ErrServiceNotFound        = sdkerrors.Register(ModuleName, 1100, "service not found")
	ErrServiceUnavailable     = sdkerrors.Register(ModuleName, 1101, "service unavailable")
	ErrServiceNotAuthorized   = sdkerrors.Register(ModuleName, 1102, "service not authorized")
	ErrInvalidServiceEndpoint = sdkerrors.Register(ModuleName, 1103, "invalid service endpoint")
	
	// General errors
	ErrInvalidRequest         = sdkerrors.Register(ModuleName, 1200, "invalid request")
	ErrInvalidAddress         = sdkerrors.Register(ModuleName, 1201, "invalid address")
	ErrInvalidPublicKey       = sdkerrors.Register(ModuleName, 1202, "invalid public key")
	ErrInvalidSignature       = sdkerrors.Register(ModuleName, 1203, "invalid signature")
	ErrDataTooLarge           = sdkerrors.Register(ModuleName, 1204, "data too large")
	ErrRateLimitExceeded      = sdkerrors.Register(ModuleName, 1205, "rate limit exceeded")
	ErrMaintenanceMode        = sdkerrors.Register(ModuleName, 1206, "identity system in maintenance mode")
	
	// Internationalization errors
	ErrUnsupportedLanguage    = sdkerrors.Register(ModuleName, 1300, "unsupported language")
	ErrMessageNotFound        = sdkerrors.Register(ModuleName, 1301, "message not found")
	ErrInvalidLocalization    = sdkerrors.Register(ModuleName, 1302, "invalid localization")
	ErrTranslationMissing     = sdkerrors.Register(ModuleName, 1303, "translation missing")
	ErrInvalidLanguageCode    = sdkerrors.Register(ModuleName, 1304, "invalid language code")
	ErrLocalizationConfigError = sdkerrors.Register(ModuleName, 1305, "localization config error")
	ErrCulturalContextMissing = sdkerrors.Register(ModuleName, 1306, "cultural context missing")
	ErrRegionalNotSupported   = sdkerrors.Register(ModuleName, 1307, "regional settings not supported")
)