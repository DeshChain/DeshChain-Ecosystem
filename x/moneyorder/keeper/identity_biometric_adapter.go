package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	identitykeeper "github.com/DeshChain/DeshChain-Ecosystem/x/identity/keeper"
	"github.com/DeshChain/DeshChain-Ecosystem/x/moneyorder/types"
)

// IdentityBiometricAdapter provides backward-compatible biometric authentication using identity module
type IdentityBiometricAdapter struct {
	keeper            *Keeper
	identityKeeper    identitykeeper.Keeper
	biometricIntegration *identitykeeper.BiometricIntegration
}

// NewIdentityBiometricAdapter creates a new adapter for biometric integration
func NewIdentityBiometricAdapter(k *Keeper, ik identitykeeper.Keeper) *IdentityBiometricAdapter {
	return &IdentityBiometricAdapter{
		keeper:               k,
		identityKeeper:       ik,
		biometricIntegration: identitykeeper.NewBiometricIntegration(&ik),
	}
}

// RegisterBiometricWithIdentity registers biometric using both old and new systems
func (iba *IdentityBiometricAdapter) RegisterBiometricWithIdentity(
	ctx sdk.Context,
	userAddress string,
	biometricType types.BiometricType,
	templateHash string,
	deviceID string,
) error {
	// First, register using traditional biometric system
	biometricMgr := NewBiometricAuthManager(iba.keeper)
	err := biometricMgr.RegisterBiometric(ctx, userAddress, biometricType, templateHash, deviceID)
	if err != nil {
		return err
	}

	// Then, create biometric credential in identity module
	biometricTypeStr := biometricType.String()
	credID, err := iba.biometricIntegration.RegisterBiometricCredential(
		ctx, userAddress, biometricTypeStr, templateHash, deviceID,
	)
	if err != nil {
		// Log error but don't fail - maintain backward compatibility
		iba.keeper.Logger(ctx).Error("Failed to create biometric credential", 
			"user", userAddress, "error", err)
	} else {
		// Store credential ID reference
		iba.storeBiometricCredentialRef(ctx, userAddress, biometricType, credID)
	}

	return nil
}

// AuthenticateWithIdentity performs biometric authentication using both systems
func (iba *IdentityBiometricAdapter) AuthenticateWithIdentity(
	ctx sdk.Context,
	userAddress string,
	biometricType types.BiometricType,
	templateHash string,
	deviceID string,
) (*EnhancedBiometricAuthResult, error) {
	result := &EnhancedBiometricAuthResult{
		UserAddress:    userAddress,
		BiometricType:  biometricType.String(),
		Timestamp:      ctx.BlockTime(),
	}

	// Try traditional authentication first
	biometricMgr := NewBiometricAuthManager(iba.keeper)
	tradResult, err := biometricMgr.AuthenticateBiometric(ctx, userAddress, biometricType, templateHash, deviceID)
	if err == nil && tradResult != nil {
		result.TraditionalAuthSuccess = tradResult.Success
		result.TraditionalAuthScore = tradResult.AuthScore
	}

	// Get security config for minimum score
	config := iba.getSecurityConfig(ctx)
	minScore := config.MinAuthScore

	// Try identity-based authentication
	identityResult, err := iba.biometricIntegration.AuthenticateBiometric(
		ctx, userAddress, biometricType.String(), templateHash, minScore,
	)
	if err == nil && identityResult != nil {
		result.IdentityAuthSuccess = identityResult.Success
		result.IdentityAuthScore = identityResult.Score
		result.CredentialID = identityResult.CredentialID
		
		if identityResult.Success {
			result.DID = fmt.Sprintf("did:desh:%s", userAddress)
		}
	}

	// Determine overall success (either system succeeds)
	result.Success = result.TraditionalAuthSuccess || result.IdentityAuthSuccess
	
	// Use the higher score
	if result.TraditionalAuthScore > result.IdentityAuthScore {
		result.FinalScore = result.TraditionalAuthScore
	} else {
		result.FinalScore = result.IdentityAuthScore
	}

	// Emit enhanced event
	iba.emitEnhancedAuthEvent(ctx, result)

	return result, nil
}

// GetBiometricStatusWithIdentity returns comprehensive biometric status
func (iba *IdentityBiometricAdapter) GetBiometricStatusWithIdentity(
	ctx sdk.Context,
	userAddress string,
) (*ComprehensiveBiometricStatus, error) {
	status := &ComprehensiveBiometricStatus{
		UserAddress: userAddress,
		Traditional: TraditionalBiometricInfo{
			Registrations: []types.BiometricRegistration{},
		},
		Identity: IdentityBiometricInfo{
			Credentials: []BiometricCredentialInfo{},
		},
	}

	// Get traditional biometric registrations
	tradRegistrations := iba.getTraditionalBiometricRegistrations(ctx, userAddress)
	status.Traditional.Registrations = tradRegistrations
	status.Traditional.HasBiometrics = len(tradRegistrations) > 0

	// Get identity-based biometric status
	identityStatus, err := iba.biometricIntegration.GetBiometricStatus(ctx, userAddress)
	if err == nil && identityStatus != nil {
		status.Identity.HasBiometrics = identityStatus.HasBiometrics
		status.Identity.DID = fmt.Sprintf("did:desh:%s", userAddress)
		
		// Convert registered biometrics to credential info
		for _, regBio := range identityStatus.RegisteredBiometrics {
			credInfo := BiometricCredentialInfo{
				CredentialID:  regBio.CredentialID,
				BiometricType: regBio.BiometricType,
				DeviceID:      regBio.DeviceID,
				IsActive:      regBio.IsActive,
				RegisteredAt:  regBio.RegisteredAt,
				ExpiresAt:     regBio.ExpiresAt,
			}
			status.Identity.Credentials = append(status.Identity.Credentials, credInfo)
		}
	}

	// Overall status
	status.HasAnyBiometric = status.Traditional.HasBiometrics || status.Identity.HasBiometrics
	status.PreferredSystem = iba.determinePreferredSystem(status)

	return status, nil
}

// MigrateBiometricsToIdentity migrates existing biometrics to identity module
func (iba *IdentityBiometricAdapter) MigrateBiometricsToIdentity(ctx sdk.Context) (int, error) {
	// Get all biometric registrations
	registrations := iba.getAllBiometricRegistrations(ctx)
	
	// Convert to integration format
	moneyOrderBiometrics := make([]identitykeeper.MoneyOrderBiometric, len(registrations))
	for i, reg := range registrations {
		moneyOrderBiometrics[i] = identitykeeper.MoneyOrderBiometric{
			UserAddress:   reg.UserAddress,
			BiometricType: int32(reg.BiometricType),
			TemplateHash:  reg.TemplateHash,
			DeviceID:      reg.DeviceId,
			RegisteredAt:  reg.RegisteredAt,
		}
	}

	// Perform migration
	return iba.biometricIntegration.MigrateMoneyOrderBiometrics(ctx, moneyOrderBiometrics)
}

// Enhanced biometric operations

// CreateBiometricPresentation creates a verifiable presentation for biometric auth
func (iba *IdentityBiometricAdapter) CreateBiometricPresentation(
	ctx sdk.Context,
	userAddress string,
	purpose string,
) (*BiometricPresentation, error) {
	did := fmt.Sprintf("did:desh:%s", userAddress)
	
	// Get biometric credentials
	status, err := iba.biometricIntegration.GetBiometricStatus(ctx, userAddress)
	if err != nil {
		return nil, err
	}

	if !status.HasBiometrics {
		return nil, sdkerrors.Wrap(types.ErrNoBiometric, "no biometric credentials found")
	}

	presentation := &BiometricPresentation{
		ID:        fmt.Sprintf("vp:biometric:%s:%d", userAddress, ctx.BlockTime().Unix()),
		Holder:    did,
		Purpose:   purpose,
		CreatedAt: ctx.BlockTime(),
		Biometrics: []BiometricClaim{},
	}

	// Add biometric claims
	for _, bio := range status.RegisteredBiometrics {
		if bio.IsActive {
			claim := BiometricClaim{
				CredentialID:  bio.CredentialID,
				BiometricType: bio.BiometricType,
				DeviceID:      bio.DeviceID,
				ValidUntil:    bio.ExpiresAt,
			}
			presentation.Biometrics = append(presentation.Biometrics, claim)
		}
	}

	return presentation, nil
}

// VerifyBiometricForHighValue performs enhanced verification for high-value transactions
func (iba *IdentityBiometricAdapter) VerifyBiometricForHighValue(
	ctx sdk.Context,
	userAddress string,
	amount sdk.Coin,
	biometricData BiometricAuthData,
) (bool, error) {
	// Check if amount requires biometric
	config := iba.getSecurityConfig(ctx)
	if !config.RequiredForHighValue || amount.Amount.LT(sdk.NewInt(config.HighValueThreshold)) {
		return true, nil // Not required for this amount
	}

	// Perform enhanced authentication
	result, err := iba.AuthenticateWithIdentity(
		ctx,
		userAddress,
		biometricData.Type,
		biometricData.TemplateHash,
		biometricData.DeviceID,
	)
	if err != nil {
		return false, err
	}

	// For high-value, require higher score and preferably both systems
	requiredScore := config.MinAuthScore * 1.1 // 10% higher for high-value
	if config.MultiFactorRequired {
		// Require both systems to succeed
		return result.TraditionalAuthSuccess && result.IdentityAuthSuccess && 
			   result.FinalScore >= requiredScore, nil
	}

	return result.Success && result.FinalScore >= requiredScore, nil
}

// Helper methods

func (iba *IdentityBiometricAdapter) getSecurityConfig(ctx sdk.Context) types.BiometricSecurityConfig {
	// Get from params or use default
	return types.DefaultBiometricSecurityConfig()
}

func (iba *IdentityBiometricAdapter) storeBiometricCredentialRef(
	ctx sdk.Context,
	userAddress string,
	biometricType types.BiometricType,
	credentialID string,
) {
	// Store mapping between traditional biometric and identity credential
	store := ctx.KVStore(iba.keeper.storeKey)
	key := []byte(fmt.Sprintf("biometric_cred_ref:%s:%d", userAddress, biometricType))
	store.Set(key, []byte(credentialID))
}

func (iba *IdentityBiometricAdapter) getTraditionalBiometricRegistrations(
	ctx sdk.Context,
	userAddress string,
) []types.BiometricRegistration {
	// Implementation to get all biometric registrations for a user
	// This is a placeholder - actual implementation would query store
	return []types.BiometricRegistration{}
}

func (iba *IdentityBiometricAdapter) getAllBiometricRegistrations(ctx sdk.Context) []types.BiometricRegistration {
	// Implementation to get all biometric registrations
	// This is a placeholder - actual implementation would iterate through store
	return []types.BiometricRegistration{}
}

func (iba *IdentityBiometricAdapter) determinePreferredSystem(status *ComprehensiveBiometricStatus) string {
	// Prefer identity system if available and has more registrations
	if status.Identity.HasBiometrics && len(status.Identity.Credentials) >= len(status.Traditional.Registrations) {
		return "identity"
	}
	if status.Traditional.HasBiometrics {
		return "traditional"
	}
	return "none"
}

func (iba *IdentityBiometricAdapter) emitEnhancedAuthEvent(ctx sdk.Context, result *EnhancedBiometricAuthResult) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"enhanced_biometric_auth",
			sdk.NewAttribute("user", result.UserAddress),
			sdk.NewAttribute("success", fmt.Sprintf("%t", result.Success)),
			sdk.NewAttribute("traditional_success", fmt.Sprintf("%t", result.TraditionalAuthSuccess)),
			sdk.NewAttribute("identity_success", fmt.Sprintf("%t", result.IdentityAuthSuccess)),
			sdk.NewAttribute("final_score", fmt.Sprintf("%.2f", result.FinalScore)),
			sdk.NewAttribute("has_did", fmt.Sprintf("%t", result.DID != "")),
		),
	)
}

// New types for enhanced biometric features

type EnhancedBiometricAuthResult struct {
	UserAddress            string
	BiometricType          string
	Success                bool
	TraditionalAuthSuccess bool
	TraditionalAuthScore   float64
	IdentityAuthSuccess    bool
	IdentityAuthScore      float64
	FinalScore             float64
	DID                    string
	CredentialID           string
	Timestamp              time.Time
}

type ComprehensiveBiometricStatus struct {
	UserAddress     string
	HasAnyBiometric bool
	PreferredSystem string
	Traditional     TraditionalBiometricInfo
	Identity        IdentityBiometricInfo
}

type TraditionalBiometricInfo struct {
	HasBiometrics bool
	Registrations []types.BiometricRegistration
}

type IdentityBiometricInfo struct {
	HasBiometrics bool
	DID           string
	Credentials   []BiometricCredentialInfo
}

type BiometricCredentialInfo struct {
	CredentialID  string
	BiometricType string
	DeviceID      string
	IsActive      bool
	RegisteredAt  time.Time
	ExpiresAt     *time.Time
}

type BiometricPresentation struct {
	ID         string
	Holder     string
	Purpose    string
	CreatedAt  time.Time
	Biometrics []BiometricClaim
}

type BiometricClaim struct {
	CredentialID  string
	BiometricType string
	DeviceID      string
	ValidUntil    *time.Time
}

type BiometricAuthData struct {
	Type         types.BiometricType
	TemplateHash string
	DeviceID     string
}