package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/DeshChain/DeshChain-Ecosystem/x/identity/types"
)

// OfflineKeeper manages offline verification functionality
type OfflineKeeper struct {
	storeKey sdk.StoreKey
	cdc      codec.BinaryCodec
	config   *types.OfflineVerificationConfig
}

// NewOfflineKeeper creates a new offline verification keeper
func NewOfflineKeeper(storeKey sdk.StoreKey, cdc codec.BinaryCodec) *OfflineKeeper {
	return &OfflineKeeper{
		storeKey: storeKey,
		cdc:      cdc,
		config:   types.DefaultOfflineVerificationConfig(),
	}
}

// PrepareOfflineVerificationData prepares offline verification data for a given identity
func (k *OfflineKeeper) PrepareOfflineVerificationData(ctx sdk.Context, did string, format types.OfflineCredentialFormat, expirationHours uint32) (*types.OfflineVerificationData, error) {
	// Get the identity
	identity, err := k.getIdentity(ctx, did)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get identity")
	}

	// Create offline verification data
	ovd := types.NewOfflineVerificationData(
		identity.DID,
		identity.IdentityHash,
		identity.PublicKey,
		identity.Controller, // Using controller as issuer DID
		identity.VerificationLevel,
	)

	// Set expiration
	if expirationHours > 0 {
		ovd.ExpiresAt = time.Now().Add(time.Duration(expirationHours) * time.Hour)
	}

	// Set format
	ovd.Format = format

	// Add identity proof
	ovd.IdentityProof = k.createIdentityProof(ctx, identity)

	// Add KYC proof if available
	if kycProof := k.createKYCProof(ctx, identity); kycProof != nil {
		ovd.KYCProof = kycProof
	}

	// Add offline-compatible credentials
	offlineCredentials, err := k.getOfflineCredentials(ctx, did)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get offline credentials")
	}
	ovd.Credentials = offlineCredentials

	// Add revocation data
	revocationData, err := k.getRevocationData(ctx, did)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get revocation data")
	}
	ovd.RevocationData = revocationData

	// Add biometric templates if available
	biometricTemplates, err := k.getEncryptedBiometricTemplates(ctx, did)
	if err == nil {
		ovd.BiometricTemplates = biometricTemplates
	}

	// Add emergency contacts
	emergencyContacts, err := k.getEmergencyContacts(ctx, did)
	if err == nil {
		ovd.EmergencyContacts = emergencyContacts
	}

	// Add localization data
	localizationData := k.getOfflineLocalizationData(ctx)
	ovd.LocalizedData = localizationData

	// Compress if requested
	if k.config.EnableCompression {
		ovd.Compressed = true
		ovd.CompressionLevel = k.config.CompressionLevel
	}

	// Compute data hash
	dataHash, err := ovd.ComputeDataHash()
	if err != nil {
		return nil, errors.Wrap(err, "failed to compute data hash")
	}
	ovd.DataHash = dataHash

	// Create signature
	signature, err := k.signOfflineData(ctx, ovd, identity.PrivateKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign offline data")
	}
	ovd.Signature = signature

	return ovd, nil
}

// VerifyOfflineIdentity verifies an identity using offline verification data
func (k *OfflineKeeper) VerifyOfflineIdentity(ctx sdk.Context, request *types.OfflineVerificationRequest, offlineData *types.OfflineVerificationData) (*types.OfflineVerificationResult, error) {
	result := &types.OfflineVerificationResult{
		Success:          false,
		DID:              request.DID,
		VerificationMode: types.OfflineModeFull,
		Timestamp:        time.Now(),
		Errors:           make([]string, 0),
		Warnings:         make([]string, 0),
	}

	// Validate offline data
	if err := offlineData.Validate(); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Invalid offline data: %v", err))
		return result, nil
	}

	// Check if DID matches
	if request.DID != offlineData.DID {
		result.Errors = append(result.Errors, "DID mismatch")
		return result, nil
	}

	// Check expiration
	if offlineData.IsExpired() {
		result.Errors = append(result.Errors, "Offline verification data has expired")
		result.VerificationMode = types.OfflineModeEmergency
		
		// In emergency mode, use lower thresholds
		if !k.config.EmergencyModeEnabled {
			return result, nil
		}
	}

	// Verify data integrity
	if valid, err := k.verifyDataIntegrity(ctx, offlineData); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to verify data integrity: %v", err))
		return result, nil
	} else if !valid {
		result.Errors = append(result.Errors, "Data integrity verification failed")
		return result, nil
	}

	// Verify identity proof
	if identityValid, confidence := k.verifyIdentityProof(ctx, offlineData.IdentityProof, request.Challenge); identityValid {
		result.IdentityVerified = true
		result.IdentityConfidence = confidence
		result.VerifiedTypes = append(result.VerifiedTypes, "identity")
	} else {
		result.Errors = append(result.Errors, "Identity proof verification failed")
	}

	// Verify KYC proof if available and required
	if offlineData.KYCProof != nil && k.requiresKYC(request.RequiredLevel) {
		if kycValid, _ := k.verifyKYCProof(ctx, offlineData.KYCProof); kycValid {
			result.KYCVerified = true
			result.VerifiedTypes = append(result.VerifiedTypes, "kyc")
		} else {
			result.Errors = append(result.Errors, "KYC proof verification failed")
		}
	}

	// Verify biometric data if provided
	if len(request.BiometricData) > 0 && len(offlineData.BiometricTemplates) > 0 {
		if biometricValid, confidence := k.verifyBiometricData(ctx, request.BiometricData, offlineData.BiometricTemplates); biometricValid {
			result.BiometricVerified = true
			result.BiometricConfidence = confidence
			result.VerifiedTypes = append(result.VerifiedTypes, "biometric")
		} else {
			result.Warnings = append(result.Warnings, "Biometric verification failed")
		}
	}

	// Verify credentials if required
	if k.requiresCredentials(request.RequiredTypes) {
		if credentialsValid := k.verifyOfflineCredentials(ctx, offlineData.Credentials, request.RequiredTypes); credentialsValid {
			result.CredentialsValid = true
			result.VerifiedTypes = append(result.VerifiedTypes, "credentials")
		} else {
			result.Errors = append(result.Errors, "Required credentials verification failed")
		}
	}

	// Check revocation status
	if !k.isRevoked(ctx, offlineData.RevocationData, request.DID) {
		result.NotRevoked = true
	} else {
		result.Errors = append(result.Errors, "Identity or credentials have been revoked")
	}

	// Calculate overall confidence
	result.OverallConfidence = k.calculateOverallConfidence(result)

	// Determine verification level achieved
	result.VerificationLevel = k.determineVerificationLevel(result)

	// Check if verification meets requirements
	requiredConfidence := k.config.RequiredConfidence
	if result.VerificationMode == types.OfflineModeEmergency {
		requiredConfidence = k.config.EmergencyThreshold
	}

	if result.OverallConfidence >= requiredConfidence && 
	   result.VerificationLevel >= request.RequiredLevel &&
	   result.IdentityVerified && 
	   result.NotRevoked {
		result.Success = true
	}

	// Set cache age if available
	if !offlineData.IssuedAt.IsZero() {
		cacheAge := time.Since(offlineData.IssuedAt)
		result.CacheAge = cacheAge
	}

	return result, nil
}

// CreateOfflineBackup creates a backup package for offline verification
func (k *OfflineKeeper) CreateOfflineBackup(ctx sdk.Context, did string, includePrivateData bool) (*types.OfflineVerificationData, error) {
	// Prepare comprehensive offline data
	ovd, err := k.PrepareOfflineVerificationData(ctx, did, types.FormatSelfContained, 720) // 30 days
	if err != nil {
		return nil, err
	}

	// Add additional backup data if private data is included
	if includePrivateData {
		// Add recovery methods
		recoveryMethods, err := k.getRecoveryMethods(ctx, did)
		if err == nil {
			// Store encrypted recovery methods
			ovd.LocalizedData.Messages["recovery_methods"] = k.encryptRecoveryMethods(recoveryMethods)
		}

		// Add additional emergency contacts
		emergencyContacts, err := k.getAllEmergencyContacts(ctx, did)
		if err == nil {
			ovd.EmergencyContacts = emergencyContacts
		}
	}

	return ovd, nil
}

// Helper methods

// getIdentity retrieves identity information (simplified)
func (k *OfflineKeeper) getIdentity(ctx sdk.Context, did string) (*types.Identity, error) {
	store := ctx.KVStore(k.storeKey)
	
	// Get identity by DID
	indexKey := types.GetIdentityByDIDIndexKey(did)
	bz := store.Get(indexKey)
	if bz == nil {
		return nil, types.ErrIdentityNotFound
	}

	// Get full identity
	identityKey := types.GetIdentityKey(string(bz))
	identityBz := store.Get(identityKey)
	if identityBz == nil {
		return nil, types.ErrIdentityNotFound
	}

	var identity types.Identity
	if err := k.cdc.Unmarshal(identityBz, &identity); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal identity")
	}

	return &identity, nil
}

// createIdentityProof creates a cryptographic proof for identity
func (k *OfflineKeeper) createIdentityProof(ctx sdk.Context, identity *types.Identity) *types.CryptographicProof {
	return types.NewCryptographicProof(
		"Ed25519Signature2020",
		identity.IdentityHash, // Simplified - should be actual signature
		"assertionMethod",
		identity.PublicKey,
	)
}

// createKYCProof creates a KYC proof if available
func (k *OfflineKeeper) createKYCProof(ctx sdk.Context, identity *types.Identity) *types.CryptographicProof {
	if identity.KYCLevel < 2 {
		return nil
	}

	return types.NewCryptographicProof(
		"KYCProof2020",
		fmt.Sprintf("kyc_level_%d", identity.KYCLevel),
		"assertionMethod",
		identity.PublicKey,
	)
}

// getOfflineCredentials retrieves credentials suitable for offline use
func (k *OfflineKeeper) getOfflineCredentials(ctx sdk.Context, did string) ([]*types.OfflineCredential, error) {
	// This is a simplified implementation
	// In practice, you would query the credential store and filter for offline-compatible credentials
	
	credentials := []*types.OfflineCredential{
		types.NewOfflineCredential(
			"offline_identity_cred_001",
			[]string{"VerifiableCredential", "IdentityCredential"},
			did,
			map[string]interface{}{
				"id": did,
				"type": "identity",
			},
		),
	}

	return credentials, nil
}

// getRevocationData retrieves revocation information for offline checking
func (k *OfflineKeeper) getRevocationData(ctx sdk.Context, did string) (*types.RevocationData, error) {
	// This would typically fetch from a revocation registry
	return &types.RevocationData{
		RevocationListURL:  "https://deshchain.org/revocation/list",
		RevocationListHash: "0x1234567890abcdef",
		LastUpdated:        time.Now().Add(-1 * time.Hour),
		ValidUntil:         time.Now().Add(24 * time.Hour),
		RevokedCredentials: []string{}, // Empty means none revoked
		RevokedIdentities:  []string{},
		MerkleRoot:         "0xabcdef1234567890",
	}, nil
}

// getEncryptedBiometricTemplates retrieves encrypted biometric templates
func (k *OfflineKeeper) getEncryptedBiometricTemplates(ctx sdk.Context, did string) (map[string]*types.EncryptedBiometric, error) {
	templates := make(map[string]*types.EncryptedBiometric)
	
	// This would fetch from biometric storage
	templates["face"] = &types.EncryptedBiometric{
		BiometricType:     types.BiometricTypeFace,
		EncryptedTemplate: "encrypted_face_template_data",
		EncryptionMethod:  "AES256-GCM",
		TemplateHash:      "face_template_hash",
		Quality:           0.95,
		CreatedAt:         time.Now().Add(-24 * time.Hour),
		ExpiresAt:         time.Now().Add(30 * 24 * time.Hour),
	}

	return templates, nil
}

// getEmergencyContacts retrieves emergency contacts
func (k *OfflineKeeper) getEmergencyContacts(ctx sdk.Context, did string) ([]*types.EmergencyContact, error) {
	contacts := []*types.EmergencyContact{
		{
			Name:            "Emergency Contact 1",
			Relationship:    "family",
			ContactMethod:   "phone:+91-9876543210",
			VerificationKey: "emergency_contact_1_key",
			Priority:        1,
		},
	}

	return contacts, nil
}

// getOfflineLocalizationData retrieves localization data for offline use
func (k *OfflineKeeper) getOfflineLocalizationData(ctx sdk.Context) *types.OfflineLocalizationData {
	return &types.OfflineLocalizationData{
		DefaultLanguage: types.LanguageEnglish,
		Messages: map[string]map[types.LanguageCode]string{
			"verification_success": {
				types.LanguageEnglish: "Identity verification successful",
				types.LanguageHindi:   "पहचान सत्यापन सफल",
			},
			"verification_failed": {
				types.LanguageEnglish: "Identity verification failed",
				types.LanguageHindi:   "पहचान सत्यापन असफल",
			},
		},
		ErrorMessages: map[string]map[types.LanguageCode]string{
			"expired_data": {
				types.LanguageEnglish: "Offline verification data has expired",
				types.LanguageHindi:   "ऑफ़लाइन सत्यापन डेटा समाप्त हो गया है",
			},
		},
	}
}

// signOfflineData signs the offline verification data
func (k *OfflineKeeper) signOfflineData(ctx sdk.Context, ovd *types.OfflineVerificationData, privateKey string) (string, error) {
	// This is a simplified implementation
	// In practice, you would use the actual private key to sign the data hash
	dataToSign := fmt.Sprintf("%s:%s", ovd.DataHash, ovd.DID)
	hash := sha256.Sum256([]byte(dataToSign))
	return fmt.Sprintf("signature_%x", hash), nil
}

// verifyDataIntegrity verifies the integrity of offline verification data
func (k *OfflineKeeper) verifyDataIntegrity(ctx sdk.Context, ovd *types.OfflineVerificationData) (bool, error) {
	// Recompute hash and compare
	expectedHash, err := ovd.ComputeDataHash()
	if err != nil {
		return false, err
	}

	return expectedHash == ovd.DataHash, nil
}

// verifyIdentityProof verifies the identity proof
func (k *OfflineKeeper) verifyIdentityProof(ctx sdk.Context, proof *types.CryptographicProof, challenge string) (bool, float64) {
	// This is a simplified implementation
	// In practice, you would verify the cryptographic signature
	if proof == nil {
		return false, 0.0
	}

	// Basic validation
	if proof.ProofValue == "" || proof.VerificationMethod == "" {
		return false, 0.0
	}

	// Return high confidence for valid proof
	return true, 0.95
}

// verifyKYCProof verifies the KYC proof
func (k *OfflineKeeper) verifyKYCProof(ctx sdk.Context, proof *types.CryptographicProof) (bool, float64) {
	// Simplified KYC proof verification
	if proof == nil {
		return false, 0.0
	}

	return proof.ProofType == "KYCProof2020", 0.90
}

// verifyBiometricData verifies biometric data against templates
func (k *OfflineKeeper) verifyBiometricData(ctx sdk.Context, biometricData map[string]string, templates map[string]*types.EncryptedBiometric) (bool, float64) {
	// This is a simplified implementation
	// In practice, you would decrypt templates and perform biometric matching
	
	for biometricType, data := range biometricData {
		if template, exists := templates[biometricType]; exists {
			// Simplified matching - in practice, use biometric matching algorithms
			if len(data) > 0 && template.Quality > 0.8 {
				return true, template.Quality
			}
		}
	}

	return false, 0.0
}

// verifyOfflineCredentials verifies offline credentials
func (k *OfflineKeeper) verifyOfflineCredentials(ctx sdk.Context, credentials []*types.OfflineCredential, requiredTypes []string) bool {
	if len(requiredTypes) == 0 {
		return true
	}

	requiredTypeMap := make(map[string]bool)
	for _, t := range requiredTypes {
		requiredTypeMap[t] = true
	}

	for _, cred := range credentials {
		for _, credType := range cred.Type {
			if requiredTypeMap[credType] {
				delete(requiredTypeMap, credType)
			}
		}
	}

	return len(requiredTypeMap) == 0
}

// isRevoked checks if an identity or credential is revoked
func (k *OfflineKeeper) isRevoked(ctx sdk.Context, revocationData *types.RevocationData, did string) bool {
	if revocationData == nil {
		return false
	}

	return revocationData.IsRevoked(did)
}

// Helper methods for calculations

// requiresKYC checks if KYC is required for the verification level
func (k *OfflineKeeper) requiresKYC(level uint32) bool {
	return level >= 2
}

// requiresCredentials checks if specific credentials are required
func (k *OfflineKeeper) requiresCredentials(requiredTypes []string) bool {
	return len(requiredTypes) > 0
}

// calculateOverallConfidence calculates the overall confidence score
func (k *OfflineKeeper) calculateOverallConfidence(result *types.OfflineVerificationResult) float64 {
	var totalWeight, weightedScore float64

	// Identity verification (weight: 0.4)
	if result.IdentityVerified {
		weightedScore += 0.4 * result.IdentityConfidence
	}
	totalWeight += 0.4

	// Biometric verification (weight: 0.3)
	if result.BiometricVerified {
		weightedScore += 0.3 * result.BiometricConfidence
	} else if result.BiometricConfidence > 0 {
		totalWeight += 0.3
	}

	// KYC verification (weight: 0.2)
	if result.KYCVerified {
		weightedScore += 0.2 * 0.9 // Fixed confidence for KYC
	}
	totalWeight += 0.2

	// Credentials and revocation (weight: 0.1)
	if result.CredentialsValid && result.NotRevoked {
		weightedScore += 0.1 * 1.0
	}
	totalWeight += 0.1

	if totalWeight == 0 {
		return 0.0
	}

	return weightedScore / totalWeight
}

// determineVerificationLevel determines the achieved verification level
func (k *OfflineKeeper) determineVerificationLevel(result *types.OfflineVerificationResult) uint32 {
	if !result.IdentityVerified {
		return 0
	}

	level := uint32(1) // Basic identity

	if result.KYCVerified {
		level = 2 // KYC verified
	}

	if result.BiometricVerified {
		level = 3 // Biometric verified
	}

	if result.CredentialsValid && result.NotRevoked {
		level = 4 // Full verification
	}

	return level
}

// Additional helper methods

// getRecoveryMethods retrieves recovery methods for backup
func (k *OfflineKeeper) getRecoveryMethods(ctx sdk.Context, did string) ([]string, error) {
	// This would fetch from recovery storage
	return []string{"social_recovery", "seed_phrase"}, nil
}

// getAllEmergencyContacts retrieves all emergency contacts
func (k *OfflineKeeper) getAllEmergencyContacts(ctx sdk.Context, did string) ([]*types.EmergencyContact, error) {
	// Extended list for backup purposes
	contacts, err := k.getEmergencyContacts(ctx, did)
	if err != nil {
		return nil, err
	}

	// Add additional contacts for backup
	additionalContacts := []*types.EmergencyContact{
		{
			Name:            "Backup Contact",
			Relationship:    "friend",
			ContactMethod:   "email:backup@example.com",
			VerificationKey: "backup_contact_key",
			Priority:        2,
		},
	}

	return append(contacts, additionalContacts...), nil
}

// encryptRecoveryMethods encrypts recovery methods for storage
func (k *OfflineKeeper) encryptRecoveryMethods(methods []string) map[types.LanguageCode]string {
	// This is a simplified implementation
	// In practice, you would encrypt the recovery methods
	data, _ := json.Marshal(methods)
	
	return map[types.LanguageCode]string{
		types.LanguageEnglish: string(data),
	}
}