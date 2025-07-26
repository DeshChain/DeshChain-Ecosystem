package types

const (
	// ModuleName defines the module name
	ModuleName = "identity"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName
)

// Key prefixes for store
var (
	// Identity storage
	IdentityPrefix          = []byte{0x01}
	IdentityIndexPrefix     = []byte{0x02}
	
	// DID storage
	DIDDocumentPrefix       = []byte{0x10}
	DIDIndexPrefix          = []byte{0x11}
	DIDVersionPrefix        = []byte{0x12}
	
	// Verifiable Credentials storage
	CredentialPrefix        = []byte{0x20}
	CredentialIndexPrefix   = []byte{0x21}
	CredentialStatusPrefix  = []byte{0x22}
	CredentialSchemaPrefix  = []byte{0x23}
	
	// Revocation storage
	RevocationListPrefix    = []byte{0x30}
	RevocationIndexPrefix   = []byte{0x31}
	
	// Privacy storage
	ZKProofPrefix           = []byte{0x40}
	NullifierPrefix         = []byte{0x41}
	AnonymousCredPrefix     = []byte{0x42}
	PrivacySettingsPrefix   = []byte{0x43}
	
	// India Stack storage
	AadhaarCredPrefix       = []byte{0x50}
	DigiLockerDocPrefix     = []byte{0x51}
	UPIIdentityPrefix       = []byte{0x52}
	DEPAConsentPrefix       = []byte{0x53}
	AAConsentPrefix         = []byte{0x54}
	VillagePanchayatPrefix  = []byte{0x55}
	JANDhanAccountPrefix    = []byte{0x56}
	IndiaStackIntegPrefix   = []byte{0x57}
	
	// Consent storage
	ConsentRecordPrefix     = []byte{0x60}
	ConsentIndexPrefix      = []byte{0x61}
	
	// Recovery storage
	RecoveryMethodPrefix    = []byte{0x70}
	RecoveryRequestPrefix   = []byte{0x71}
	
	// Service registry
	ServiceRegistryPrefix   = []byte{0x80}
	ServiceIndexPrefix      = []byte{0x81}
	
	// Issuer registry
	IssuerRegistryPrefix    = []byte{0x90}
	IssuerIndexPrefix       = []byte{0x91}
	
	// Analytics storage
	AnalyticsPrefix         = []byte{0xA0}
	UsageMetricsPrefix      = []byte{0xA1}
	
	// Configuration storage
	ModuleParamsPrefix      = []byte{0xB0}
	GenesisStatePrefix      = []byte{0xB1}
	
	// Internationalization storage
	LocalizationConfigKey   = []byte{0xE0}
	CustomMessagePrefix     = []byte{0xE1}
	UserLanguagePrefix      = []byte{0xE2}
	MessageCatalogPrefix    = []byte{0xE3}
)

// GetIdentityKey returns the key for an identity
func GetIdentityKey(address string) []byte {
	return append(IdentityPrefix, []byte(address)...)
}

// GetDIDDocumentKey returns the key for a DID document
func GetDIDDocumentKey(did string) []byte {
	return append(DIDDocumentPrefix, []byte(did)...)
}

// GetDIDVersionKey returns the key for a specific DID document version
func GetDIDVersionKey(did string, version string) []byte {
	return append(append(DIDVersionPrefix, []byte(did)...), []byte(version)...)
}

// GetCredentialKey returns the key for a verifiable credential
func GetCredentialKey(credentialID string) []byte {
	return append(CredentialPrefix, []byte(credentialID)...)
}

// GetCredentialStatusKey returns the key for credential status
func GetCredentialStatusKey(credentialID string) []byte {
	return append(CredentialStatusPrefix, []byte(credentialID)...)
}

// GetCredentialSchemaKey returns the key for credential schema
func GetCredentialSchemaKey(schemaID string) []byte {
	return append(CredentialSchemaPrefix, []byte(schemaID)...)
}

// GetRevocationListKey returns the key for revocation list
func GetRevocationListKey(issuer string) []byte {
	return append(RevocationListPrefix, []byte(issuer)...)
}

// GetZKProofKey returns the key for a ZK proof
func GetZKProofKey(proofID string) []byte {
	return append(ZKProofPrefix, []byte(proofID)...)
}

// GetNullifierKey returns the key for a nullifier
func GetNullifierKey(nullifier string) []byte {
	return append(NullifierPrefix, []byte(nullifier)...)
}

// GetPrivacySettingsKey returns the key for privacy settings
func GetPrivacySettingsKey(address string) []byte {
	return append(PrivacySettingsPrefix, []byte(address)...)
}

// GetAadhaarCredentialKey returns the key for Aadhaar credential
func GetAadhaarCredentialKey(credID string) []byte {
	return append(AadhaarCredPrefix, []byte(credID)...)
}

// GetDigiLockerDocumentKey returns the key for DigiLocker document
func GetDigiLockerDocumentKey(docID string) []byte {
	return append(DigiLockerDocPrefix, []byte(docID)...)
}

// GetUPIIdentityKey returns the key for UPI identity
func GetUPIIdentityKey(identityID string) []byte {
	return append(UPIIdentityPrefix, []byte(identityID)...)
}

// GetDEPAConsentKey returns the key for DEPA consent
func GetDEPAConsentKey(consentID string) []byte {
	return append(DEPAConsentPrefix, []byte(consentID)...)
}

// GetAAConsentKey returns the key for AA consent
func GetAAConsentKey(consentID string) []byte {
	return append(AAConsentPrefix, []byte(consentID)...)
}

// GetVillagePanchayatKYCKey returns the key for village panchayat KYC
func GetVillagePanchayatKYCKey(kycID string) []byte {
	return append(VillagePanchayatPrefix, []byte(kycID)...)
}

// GetJANDhanAccountKey returns the key for JAN Dhan account
func GetJANDhanAccountKey(accountID string) []byte {
	return append(JANDhanAccountPrefix, []byte(accountID)...)
}

// GetIndiaStackIntegrationKey returns the key for India Stack integration
func GetIndiaStackIntegrationKey(address string) []byte {
	return append(IndiaStackIntegPrefix, []byte(address)...)
}

// GetConsentRecordKey returns the key for consent record
func GetConsentRecordKey(consentID string) []byte {
	return append(ConsentRecordPrefix, []byte(consentID)...)
}

// GetRecoveryMethodKey returns the key for recovery method
func GetRecoveryMethodKey(address string, methodType string) []byte {
	return append(append(RecoveryMethodPrefix, []byte(address)...), []byte(methodType)...)
}

// GetServiceRegistryKey returns the key for service registry
func GetServiceRegistryKey(serviceID string) []byte {
	return append(ServiceRegistryPrefix, []byte(serviceID)...)
}

// GetIssuerRegistryKey returns the key for issuer registry
func GetIssuerRegistryKey(issuerDID string) []byte {
	return append(IssuerRegistryPrefix, []byte(issuerDID)...)
}

// GetUsageMetricsKey returns the key for usage metrics
func GetUsageMetricsKey(address string, metricType string) []byte {
	return append(append(UsageMetricsPrefix, []byte(address)...), []byte(metricType)...)
}

// Index keys for efficient queries

// GetIdentityByDIDIndexKey returns the index key for identity by DID lookup
func GetIdentityByDIDIndexKey(did string) []byte {
	return append(append(IdentityIndexPrefix, []byte("did:")...), []byte(did)...)
}

// GetCredentialByHolderIndexKey returns the index key for credentials by holder
func GetCredentialByHolderIndexKey(holder string, credentialID string) []byte {
	return append(append(append(CredentialIndexPrefix, []byte("holder:")...), []byte(holder)...), []byte(credentialID)...)
}

// GetCredentialByIssuerIndexKey returns the index key for credentials by issuer
func GetCredentialByIssuerIndexKey(issuer string, credentialID string) []byte {
	return append(append(append(CredentialIndexPrefix, []byte("issuer:")...), []byte(issuer)...), []byte(credentialID)...)
}

// GetCredentialByTypeIndexKey returns the index key for credentials by type
func GetCredentialByTypeIndexKey(credType string, credentialID string) []byte {
	return append(append(append(CredentialIndexPrefix, []byte("type:")...), []byte(credType)...), []byte(credentialID)...)
}

// GetConsentByTypeIndexKey returns the index key for consents by type
func GetConsentByTypeIndexKey(address string, consentType string) []byte {
	return append(append(append(ConsentIndexPrefix, []byte(address)...), []byte(":")...), []byte(consentType)...)
}

// GetServiceByTypeIndexKey returns the index key for services by type
func GetServiceByTypeIndexKey(serviceType string, serviceID string) []byte {
	return append(append(append(ServiceIndexPrefix, []byte("type:")...), []byte(serviceType)...), []byte(serviceID)...)
}

// Sharing Protocol Keys
var (
	ShareRequestKeyPrefix       = []byte{0xC0}
	ShareResponseKeyPrefix      = []byte{0xC1}
	SharingAgreementKeyPrefix   = []byte{0xC2}
	AccessPolicyKeyPrefix       = []byte{0xC3}
	AccessPolicyByHolderPrefix  = []byte{0xC4}
	ShareAuditLogKeyPrefix      = []byte{0xC5}
	ShareAuditLogByHolderPrefix = []byte{0xC6}
	ConsentByHolderKeyPrefix    = []byte{0xC7}
	BiometricTemplateKeyPrefix  = []byte{0xC8}
	PresentationKeyPrefix       = []byte{0xC9}
)

// ShareRequestKey returns the key for a share request
func ShareRequestKey(requestID string) []byte {
	return append(ShareRequestKeyPrefix, []byte(requestID)...)
}

// ShareResponseKey returns the key for a share response
func ShareResponseKey(requestID string) []byte {
	return append(ShareResponseKeyPrefix, []byte(requestID)...)
}

// SharingAgreementKey returns the key for a sharing agreement
func SharingAgreementKey(agreementID string) []byte {
	return append(SharingAgreementKeyPrefix, []byte(agreementID)...)
}

// AccessPolicyKey returns the key for an access policy
func AccessPolicyKey(policyID string) []byte {
	return append(AccessPolicyKeyPrefix, []byte(policyID)...)
}

// AccessPolicyByHolderKey returns the key for access policy by holder index
func AccessPolicyByHolderKey(holderDID, policyID string) []byte {
	return append(AccessPolicyByHolderPrefix, []byte(holderDID+"/"+policyID)...)
}

// AccessPolicyByHolderPrefix returns the prefix for access policies by holder
func AccessPolicyByHolderPrefix(holderDID string) []byte {
	return append(AccessPolicyByHolderPrefix, []byte(holderDID+"/")...)
}

// ShareAuditLogKey returns the key for a share audit log
func ShareAuditLogKey(logID string) []byte {
	return append(ShareAuditLogKeyPrefix, []byte(logID)...)
}

// ShareAuditLogByHolderKey returns the key for audit log by holder index
func ShareAuditLogByHolderKey(holderDID, logID string) []byte {
	return append(ShareAuditLogByHolderPrefix, []byte(holderDID+"/"+logID)...)
}

// ShareAuditLogByHolderPrefix returns the prefix for audit logs by holder
func ShareAuditLogByHolderPrefix(holderDID string) []byte {
	return append(ShareAuditLogByHolderPrefix, []byte(holderDID+"/")...)
}

// ConsentByHolderKey returns the key for consent by holder index
func ConsentByHolderKey(holderDID, consentID string) []byte {
	return append(ConsentByHolderKeyPrefix, []byte(holderDID+"/"+consentID)...)
}

// ConsentByHolderPrefix returns the prefix for consent entries by holder
func ConsentByHolderPrefix(holderDID string) []byte {
	return append(ConsentByHolderKeyPrefix, []byte(holderDID+"/")...)
}

// BiometricTemplateKey returns the key for a biometric template
func BiometricTemplateKey(holderDID, templateType string) []byte {
	return append(BiometricTemplateKeyPrefix, []byte(holderDID+"/"+templateType)...)
}

// PresentationKey returns the key for a presentation
func PresentationKey(presentationID string) []byte {
	return append(PresentationKeyPrefix, []byte(presentationID)...)
}

// Recovery System Keys
var (
	IdentityBackupKeyPrefix        = []byte{0xD0}
	RecoveryRequestKeyPrefix       = []byte{0xD1}
	SocialRecoveryGuardianKeyPrefix = []byte{0xD2}
	GuardianVoteKeyPrefix          = []byte{0xD3}
	BackupVerificationResultKeyPrefix = []byte{0xD4}
	BackupByHolderKeyPrefix        = []byte{0xD5}
	RecoveryRequestByHolderKeyPrefix = []byte{0xD6}
	GuardianByHolderKeyPrefix      = []byte{0xD7}
	DisasterRecoveryConfigKeyPrefix = []byte{0xD8}
)

// IdentityBackupKey returns the key for an identity backup
func IdentityBackupKey(backupID string) []byte {
	return append(IdentityBackupKeyPrefix, []byte(backupID)...)
}

// RecoveryRequestKey returns the key for a recovery request
func RecoveryRequestKey(requestID string) []byte {
	return append(RecoveryRequestKeyPrefix, []byte(requestID)...)
}

// SocialRecoveryGuardianKey returns the key for a social recovery guardian
func SocialRecoveryGuardianKey(guardianID string) []byte {
	return append(SocialRecoveryGuardianKeyPrefix, []byte(guardianID)...)
}

// GuardianVoteKey returns the key for a guardian vote
func GuardianVoteKey(voteID string) []byte {
	return append(GuardianVoteKeyPrefix, []byte(voteID)...)
}

// BackupVerificationResultKey returns the key for a backup verification result
func BackupVerificationResultKey(verificationID string) []byte {
	return append(BackupVerificationResultKeyPrefix, []byte(verificationID)...)
}

// BackupByHolderKey returns the key for backup by holder index
func BackupByHolderKey(holderDID, backupID string) []byte {
	return append(BackupByHolderKeyPrefix, []byte(holderDID+"/"+backupID)...)
}

// BackupByHolderPrefix returns the prefix for backups by holder
func BackupByHolderPrefix(holderDID string) []byte {
	return append(BackupByHolderKeyPrefix, []byte(holderDID+"/")...)
}

// RecoveryRequestByHolderKey returns the key for recovery request by holder index
func RecoveryRequestByHolderKey(holderDID, requestID string) []byte {
	return append(RecoveryRequestByHolderKeyPrefix, []byte(holderDID+"/"+requestID)...)
}

// RecoveryRequestByHolderPrefix returns the prefix for recovery requests by holder
func RecoveryRequestByHolderPrefix(holderDID string) []byte {
	return append(RecoveryRequestByHolderKeyPrefix, []byte(holderDID+"/")...)
}

// GuardianByHolderKey returns the key for guardian by holder index
func GuardianByHolderKey(holderDID, guardianID string) []byte {
	return append(GuardianByHolderKeyPrefix, []byte(holderDID+"/"+guardianID)...)
}

// GuardianByHolderPrefix returns the prefix for guardians by holder
func GuardianByHolderPrefix(holderDID string) []byte {
	return append(GuardianByHolderKeyPrefix, []byte(holderDID+"/")...)
}

// DisasterRecoveryConfigKey returns the key for disaster recovery config
func DisasterRecoveryConfigKey(holderDID string) []byte {
	return append(DisasterRecoveryConfigKeyPrefix, []byte(holderDID)...)
}

// Internationalization helper functions

// GetCustomMessageKey returns the key for a custom message
func GetCustomMessageKey(messageKey string) []byte {
	return append(CustomMessagePrefix, []byte(messageKey)...)
}

// GetUserLanguageKey returns the key for user language preference
func GetUserLanguageKey(userDID string) []byte {
	return append(UserLanguagePrefix, []byte(userDID)...)
}

// GetMessageCatalogKey returns the key for message catalog
func GetMessageCatalogKey(catalogID string) []byte {
	return append(MessageCatalogPrefix, []byte(catalogID)...)
}