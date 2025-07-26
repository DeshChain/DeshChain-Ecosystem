package keeper

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/identity/types"
)

// BiometricIntegration provides backward-compatible biometric integration for moneyorder module
type BiometricIntegration struct {
	keeper *Keeper
}

// NewBiometricIntegration creates a new biometric integration adapter
func NewBiometricIntegration(k *Keeper) *BiometricIntegration {
	return &BiometricIntegration{keeper: k}
}

// RegisterBiometricCredential creates a biometric credential for a user
func (bi *BiometricIntegration) RegisterBiometricCredential(
	ctx sdk.Context,
	userAddress string,
	biometricType string,
	templateHash string,
	deviceID string,
) (string, error) {
	// Get or create user identity
	did := fmt.Sprintf("did:desh:%s", userAddress)
	identity, exists := bi.keeper.GetIdentity(ctx, did)
	if !exists {
		// Create new identity
		identity = types.Identity{
			Did:        did,
			Controller: userAddress,
			Status:     types.IdentityStatus_ACTIVE,
			CreatedAt:  ctx.BlockTime(),
			UpdatedAt:  ctx.BlockTime(),
			Metadata: map[string]string{
				"source": "moneyorder",
				"type":   "user",
			},
		}
		bi.keeper.SetIdentity(ctx, identity)
	}

	// Generate biometric credential ID
	credID := fmt.Sprintf("vc:biometric:%s:%s:%d", userAddress, biometricType, ctx.BlockTime().Unix())

	// Create biometric data hash
	dataToHash := fmt.Sprintf("%s:%s:%s:%s", userAddress, biometricType, templateHash, deviceID)
	hash := sha256.Sum256([]byte(dataToHash))
	biometricHash := hex.EncodeToString(hash[:])

	// Create biometric credential
	expiryDate := ctx.BlockTime().AddDate(2, 0, 0) // 2 years expiry
	credential := types.VerifiableCredential{
		Context: []string{
			"https://www.w3.org/2018/credentials/v1",
			"https://deshchain.com/contexts/biometric/v1",
		},
		ID:   credID,
		Type: []string{"VerifiableCredential", "BiometricCredential"},
		Issuer: "did:desh:biometric-issuer",
		IssuanceDate: ctx.BlockTime(),
		ExpirationDate: &expiryDate,
		CredentialSubject: map[string]interface{}{
			"id":             did,
			"biometric_type": biometricType,
			"template_hash":  biometricHash,
			"device_id":      deviceID,
			"registered_at":  ctx.BlockTime().Format(time.RFC3339),
		},
		Proof: &types.Proof{
			Type:               "Ed25519Signature2020",
			Created:            ctx.BlockTime(),
			VerificationMethod: "did:desh:biometric-issuer#key-1",
			ProofPurpose:       "assertionMethod",
			ProofValue:         "mock-biometric-signature",
		},
	}

	// Store credential
	bi.keeper.SetCredential(ctx, credential)
	bi.keeper.AddCredentialToSubject(ctx, did, credID)

	// Create biometric authentication method in DID document
	bi.addBiometricAuthMethod(ctx, did, biometricType, biometricHash)

	return credID, nil
}

// AuthenticateBiometric verifies biometric authentication using identity credentials
func (bi *BiometricIntegration) AuthenticateBiometric(
	ctx sdk.Context,
	userAddress string,
	biometricType string,
	templateHash string,
	minScore float64,
) (*BiometricAuthResult, error) {
	did := fmt.Sprintf("did:desh:%s", userAddress)
	
	// Get identity
	identity, exists := bi.keeper.GetIdentity(ctx, did)
	if !exists {
		return &BiometricAuthResult{
			Success: false,
			Error:   "Identity not found",
		}, nil
	}

	// Check if identity is active
	if identity.Status != types.IdentityStatus_ACTIVE {
		return &BiometricAuthResult{
			Success: false,
			Error:   "Identity is not active",
		}, nil
	}

	// Find biometric credentials
	credentials := bi.keeper.GetCredentialsBySubject(ctx, did)
	for _, credID := range credentials {
		cred, found := bi.keeper.GetCredential(ctx, credID)
		if !found {
			continue
		}

		// Check if it's a biometric credential of the right type
		if !bi.isBiometricCredential(cred, biometricType) {
			continue
		}

		// Verify credential is valid
		if cred.Status != nil && cred.Status.Type == "revoked" {
			continue
		}

		// Check expiry
		if cred.ExpirationDate != nil && cred.ExpirationDate.Before(ctx.BlockTime()) {
			continue
		}

		// Extract stored template hash
		if subject, ok := cred.CredentialSubject.(map[string]interface{}); ok {
			if storedHash, ok := subject["template_hash"].(string); ok {
				// Simulate biometric matching (in production, use actual matching)
				score := bi.calculateBiometricScore(templateHash, storedHash)
				
				if score >= minScore {
					// Update last authentication time
					bi.updateLastAuth(ctx, credID)
					
					return &BiometricAuthResult{
						Success:      true,
						Score:        score,
						CredentialID: credID,
					}, nil
				}
			}
		}
	}

	return &BiometricAuthResult{
		Success: false,
		Error:   "No matching biometric credential found",
	}, nil
}

// GetBiometricStatus returns the biometric registration status for a user
func (bi *BiometricIntegration) GetBiometricStatus(
	ctx sdk.Context,
	userAddress string,
) (*BiometricStatus, error) {
	did := fmt.Sprintf("did:desh:%s", userAddress)
	
	status := &BiometricStatus{
		UserAddress:          userAddress,
		RegisteredBiometrics: []RegisteredBiometric{},
	}

	// Get all credentials for the user
	credentials := bi.keeper.GetCredentialsBySubject(ctx, did)
	for _, credID := range credentials {
		cred, found := bi.keeper.GetCredential(ctx, credID)
		if !found {
			continue
		}

		// Check if it's a biometric credential
		isBiometric := false
		for _, credType := range cred.Type {
			if credType == "BiometricCredential" {
				isBiometric = true
				break
			}
		}

		if !isBiometric {
			continue
		}

		// Extract biometric info
		if subject, ok := cred.CredentialSubject.(map[string]interface{}); ok {
			biometricType := ""
			deviceID := ""
			
			if bt, ok := subject["biometric_type"].(string); ok {
				biometricType = bt
			}
			if dev, ok := subject["device_id"].(string); ok {
				deviceID = dev
			}

			regBiometric := RegisteredBiometric{
				CredentialID:  credID,
				BiometricType: biometricType,
				DeviceID:      deviceID,
				RegisteredAt:  cred.IssuanceDate,
				IsActive:      cred.Status == nil || cred.Status.Type != "revoked",
				ExpiresAt:     cred.ExpirationDate,
			}

			status.RegisteredBiometrics = append(status.RegisteredBiometrics, regBiometric)
		}
	}

	status.HasBiometrics = len(status.RegisteredBiometrics) > 0
	return status, nil
}

// DisableBiometric revokes a biometric credential
func (bi *BiometricIntegration) DisableBiometric(
	ctx sdk.Context,
	credentialID string,
	reason string,
) error {
	cred, found := bi.keeper.GetCredential(ctx, credentialID)
	if !found {
		return fmt.Errorf("credential not found: %s", credentialID)
	}

	// Update credential status to revoked
	return bi.keeper.UpdateCredentialStatus(ctx, credentialID, &types.CredentialStatus{
		Type:   "revoked",
		Reason: reason,
	})
}

// MigrateMoneyOrderBiometrics migrates existing moneyorder biometrics to identity module
func (bi *BiometricIntegration) MigrateMoneyOrderBiometrics(
	ctx sdk.Context,
	biometrics []MoneyOrderBiometric,
) (int, error) {
	migrated := 0

	for _, bio := range biometrics {
		// Convert biometric type
		biometricType := bi.convertBiometricType(bio.BiometricType)
		
		// Register biometric credential
		if _, err := bi.RegisterBiometricCredential(
			ctx,
			bio.UserAddress,
			biometricType,
			bio.TemplateHash,
			bio.DeviceID,
		); err != nil {
			bi.keeper.Logger(ctx).Error("Failed to migrate biometric",
				"user", bio.UserAddress,
				"error", err)
			continue
		}

		migrated++
	}

	return migrated, nil
}

// Helper functions

func (bi *BiometricIntegration) addBiometricAuthMethod(
	ctx sdk.Context,
	did string,
	biometricType string,
	biometricHash string,
) error {
	// Get DID document
	didDoc, found := bi.keeper.GetDIDDocument(ctx, did)
	if !found {
		// Create new DID document
		didDoc = types.DIDDocument{
			Context:    []string{"https://www.w3.org/ns/did/v1"},
			ID:         did,
			Controller: did,
			Created:    ctx.BlockTime(),
			Updated:    ctx.BlockTime(),
		}
	}

	// Add biometric authentication method
	authMethod := types.VerificationMethod{
		ID:         fmt.Sprintf("%s#biometric-%s", did, biometricType),
		Type:       "BiometricAuthentication2023",
		Controller: did,
		PublicKeyMultibase: biometricHash,
	}

	didDoc.VerificationMethod = append(didDoc.VerificationMethod, authMethod)
	didDoc.Authentication = append(didDoc.Authentication, authMethod.ID)
	didDoc.Updated = ctx.BlockTime()

	// Store updated DID document
	return bi.keeper.SetDIDDocument(ctx, didDoc)
}

func (bi *BiometricIntegration) isBiometricCredential(cred types.VerifiableCredential, biometricType string) bool {
	// Check if it's a biometric credential
	isBiometric := false
	for _, credType := range cred.Type {
		if credType == "BiometricCredential" {
			isBiometric = true
			break
		}
	}

	if !isBiometric {
		return false
	}

	// Check biometric type if specified
	if biometricType != "" {
		if subject, ok := cred.CredentialSubject.(map[string]interface{}); ok {
			if bt, ok := subject["biometric_type"].(string); ok {
				return bt == biometricType
			}
		}
		return false
	}

	return true
}

func (bi *BiometricIntegration) calculateBiometricScore(templateHash, storedHash string) float64 {
	// Mock biometric matching score
	// In production, this would use actual biometric matching algorithms
	if templateHash == storedHash {
		return 1.0
	}
	
	// Simulate partial match based on hash similarity
	matches := 0
	minLen := len(templateHash)
	if len(storedHash) < minLen {
		minLen = len(storedHash)
	}

	for i := 0; i < minLen && i < 10; i++ {
		if templateHash[i] == storedHash[i] {
			matches++
		}
	}

	return float64(matches) / 10.0
}

func (bi *BiometricIntegration) updateLastAuth(ctx sdk.Context, credentialID string) {
	// In production, update last authentication timestamp
	// This could be stored as credential metadata
}

func (bi *BiometricIntegration) convertBiometricType(moType int32) string {
	switch moType {
	case 0:
		return "FINGERPRINT"
	case 1:
		return "FACE"
	case 2:
		return "IRIS"
	case 3:
		return "VOICE"
	case 4:
		return "PALM"
	default:
		return "UNKNOWN"
	}
}

// Backward compatibility types
type BiometricAuthResult struct {
	Success      bool
	Error        string
	Score        float64
	CredentialID string
}

type BiometricStatus struct {
	UserAddress          string
	HasBiometrics        bool
	RegisteredBiometrics []RegisteredBiometric
}

type RegisteredBiometric struct {
	CredentialID  string
	BiometricType string
	DeviceID      string
	RegisteredAt  time.Time
	IsActive      bool
	ExpiresAt     *time.Time
}

type MoneyOrderBiometric struct {
	UserAddress   string
	BiometricType int32
	TemplateHash  string
	DeviceID      string
	RegisteredAt  time.Time
}