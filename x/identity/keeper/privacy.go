package keeper

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/DeshChain/DeshChain-Ecosystem/x/identity/types"
)

// SetZKProof stores a zero-knowledge proof
func (k Keeper) SetZKProof(ctx sdk.Context, proof types.ZKProof) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ZKProofPrefix)
	b := k.cdc.MustMarshal(&proof)
	store.Set([]byte(proof.ID), b)
}

// GetZKProof retrieves a zero-knowledge proof
func (k Keeper) GetZKProof(ctx sdk.Context, proofID string) (types.ZKProof, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ZKProofPrefix)
	b := store.Get([]byte(proofID))
	if b == nil {
		return types.ZKProof{}, false
	}
	
	var proof types.ZKProof
	k.cdc.MustUnmarshal(b, &proof)
	return proof, true
}

// DeleteZKProof removes a zero-knowledge proof
func (k Keeper) DeleteZKProof(ctx sdk.Context, proofID string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ZKProofPrefix)
	store.Delete([]byte(proofID))
}

// IterateZKProofs iterates over all ZK proofs
func (k Keeper) IterateZKProofs(ctx sdk.Context, cb func(proof types.ZKProof) bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ZKProofPrefix)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var proof types.ZKProof
		k.cdc.MustUnmarshal(iterator.Value(), &proof)
		if cb(proof) {
			break
		}
	}
}

// SetNullifier stores a nullifier to prevent double-spending
func (k Keeper) SetNullifier(ctx sdk.Context, nullifier string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NullifierPrefix)
	store.Set([]byte(nullifier), []byte{1})
}

// HasNullifier checks if a nullifier exists
func (k Keeper) HasNullifier(ctx sdk.Context, nullifier string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NullifierPrefix)
	return store.Has([]byte(nullifier))
}

// SetPrivacySettings stores privacy settings for a user
func (k Keeper) SetPrivacySettings(ctx sdk.Context, settings types.PrivacySettings) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrivacySettingsPrefix)
	b := k.cdc.MustMarshal(&settings)
	store.Set([]byte(settings.UserAddress), b)
}

// GetPrivacySettings retrieves privacy settings for a user
func (k Keeper) GetPrivacySettings(ctx sdk.Context, address string) (types.PrivacySettings, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrivacySettingsPrefix)
	b := store.Get([]byte(address))
	if b == nil {
		// Return default settings if not found
		return types.PrivacySettings{
			UserAddress:              address,
			DefaultDisclosureLevel:   types.DisclosureLevel_STANDARD,
			AllowAnonymousUsage:      false,
			RequireExplicitConsent:   true,
			DataMinimization:         true,
			AutoDeleteAfterDays:      365,
			AllowDerivedCredentials:  true,
			PreferredProofSystems:    []string{types.ProofSystemGroth16},
			BlacklistedVerifiers:     []string{},
			UpdatedAt:                ctx.BlockTime(),
		}, false
	}
	
	var settings types.PrivacySettings
	k.cdc.MustUnmarshal(b, &settings)
	return settings, true
}

// IteratePrivacySettings iterates over all privacy settings
func (k Keeper) IteratePrivacySettings(ctx sdk.Context, cb func(settings types.PrivacySettings) bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PrivacySettingsPrefix)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var settings types.PrivacySettings
		k.cdc.MustUnmarshal(iterator.Value(), &settings)
		if cb(settings) {
			break
		}
	}
}

// CreateZKProof creates and stores a new zero-knowledge proof
func (k Keeper) CreateZKProof(
	ctx sdk.Context,
	creator string,
	proofType types.ZKProofType,
	statement string,
	proofData []byte,
	expiryMinutes int32,
) (*types.ZKProof, error) {
	// Check if ZK proofs are enabled
	if !k.EnableZKProofs(ctx) {
		return nil, types.ErrInvalidRequest
	}
	
	// Validate creator has identity
	if !k.HasIdentity(ctx, creator) {
		return nil, types.ErrIdentityNotFound
	}
	
	// Check proof size
	maxSize := k.MaxProofSize(ctx)
	if uint64(len(proofData)) > maxSize {
		return nil, types.ErrDataTooLarge
	}
	
	// Generate proof ID
	proofID := fmt.Sprintf("proof:%s:%d", sdk.NewRand().Str(16), ctx.BlockHeight())
	
	// Calculate expiry
	var expiresAt *time.Time
	if expiryMinutes > 0 {
		expiry := ctx.BlockTime().Add(time.Duration(expiryMinutes) * time.Minute)
		expiresAt = &expiry
	} else {
		// Use default expiry
		expiry := ctx.BlockTime().Add(time.Duration(k.ProofExpiryMinutes(ctx)) * time.Minute)
		expiresAt = &expiry
	}
	
	// Create proof
	proof := types.ZKProof{
		ID:          proofID,
		Type:        proofType,
		ProofSystem: types.ProofSystemGroth16, // Default, could be parameterized
		Statement:   statement,
		ProofData:   proofData,
		CreatedAt:   ctx.BlockTime(),
		ExpiresAt:   expiresAt,
		Metadata: map[string]string{
			"creator":     creator,
			"blockHeight": fmt.Sprintf("%d", ctx.BlockHeight()),
		},
	}
	
	// Store proof
	k.SetZKProof(ctx, proof)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeZKProofCreated,
			sdk.NewAttribute(types.AttributeKeyProofType, proofType.String()),
			sdk.NewAttribute("proof_id", proofID),
			sdk.NewAttribute("creator", creator),
		),
	)
	
	return &proof, nil
}

// VerifyZKProof verifies a zero-knowledge proof
func (k Keeper) VerifyZKProof(ctx sdk.Context, proofID string, verifier string) error {
	// Get proof
	proof, found := k.GetZKProof(ctx, proofID)
	if !found {
		return types.ErrInvalidZKProof
	}
	
	// Check if expired
	if proof.ExpiresAt != nil && proof.ExpiresAt.Before(ctx.BlockTime()) {
		return types.ErrProofExpired
	}
	
	// Basic validation
	if err := types.VerifyZKProof(&proof); err != nil {
		return err
	}
	
	// Check if proof system is supported
	params := k.GetParams(ctx)
	if !params.IsProofSystemSupported(proof.ProofSystem) {
		return types.ErrInvalidZKProof
	}
	
	// Check privacy settings of proof creator
	if creatorAddr, ok := proof.Metadata["creator"]; ok {
		settings, _ := k.GetPrivacySettings(ctx, creatorAddr)
		
		// Check if verifier is blacklisted
		for _, blacklisted := range settings.BlacklistedVerifiers {
			if blacklisted == verifier {
				return types.ErrPrivacyViolation
			}
		}
	}
	
	// TODO: Actual cryptographic verification would go here
	// This would involve calling the appropriate proof verification function
	// based on the proof system (Groth16, PLONK, etc.)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeZKProofVerified,
			sdk.NewAttribute(types.AttributeKeyProofType, proof.Type.String()),
			sdk.NewAttribute("proof_id", proofID),
			sdk.NewAttribute("verifier", verifier),
		),
	)
	
	return nil
}

// CreateAnonymousCredential creates an anonymous credential
func (k Keeper) CreateAnonymousCredential(
	ctx sdk.Context,
	originalCredentialID string,
	blindingFactor string,
) (*types.AnonymousCredential, error) {
	// Check if anonymous credentials are enabled
	if !k.EnableAnonymousCredentials(ctx) {
		return nil, types.ErrInvalidRequest
	}
	
	// Get original credential
	credential, found := k.GetCredential(ctx, originalCredentialID)
	if !found {
		return nil, types.ErrCredentialNotFound
	}
	
	// Verify credential is valid
	if err := k.VerifyCredential(ctx, originalCredentialID); err != nil {
		return nil, err
	}
	
	// Generate nullifier to prevent double-spending
	nullifier := types.GenerateNullifier(originalCredentialID, "anonymous")
	
	// Check if nullifier already used
	if k.HasNullifier(ctx, nullifier) {
		return nil, types.ErrNullifierUsed
	}
	
	// Create anonymous credential
	anonCred := types.AnonymousCredential{
		ID:                fmt.Sprintf("anon:%s", sdk.NewRand().Str(32)),
		Type:              credential.Type[0], // Use primary type
		BlindSignature:    blindingFactor,     // In real implementation, this would be cryptographically generated
		Nullifier:         nullifier,
		IssuanceTimestamp: ctx.BlockTime(),
		Metadata: map[string]string{
			"originalType": credential.Type[0],
		},
	}
	
	// Store nullifier
	k.SetNullifier(ctx, nullifier)
	
	// Store anonymous credential
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.AnonymousCredPrefix)
	b := k.cdc.MustMarshal(&anonCred)
	store.Set([]byte(anonCred.ID), b)
	
	return &anonCred, nil
}

// CreateSelectiveDisclosureProof creates a proof for selective disclosure
func (k Keeper) CreateSelectiveDisclosureProof(
	ctx sdk.Context,
	credentialID string,
	revealedClaims []string,
	holder string,
) (*types.SelectiveDisclosure, error) {
	// Get credential
	credential, found := k.GetCredential(ctx, credentialID)
	if !found {
		return nil, types.ErrCredentialNotFound
	}
	
	// Verify holder owns the credential
	subjectID, err := credential.GetSubjectID()
	if err != nil || subjectID != holder {
		return nil, types.ErrUnauthorized
	}
	
	// Get all claims
	var allClaims []string
	var blindedClaims []string
	
	switch subject := credential.CredentialSubject.(type) {
	case map[string]interface{}:
		for claim := range subject {
			allClaims = append(allClaims, claim)
		}
	case types.CredentialSubject:
		for claim := range subject.Claims {
			allClaims = append(allClaims, claim)
		}
	}
	
	// Determine blinded claims
	for _, claim := range allClaims {
		revealed := false
		for _, revealedClaim := range revealedClaims {
			if claim == revealedClaim {
				revealed = true
				break
			}
		}
		if !revealed {
			blindedClaims = append(blindedClaims, claim)
		}
	}
	
	// Create selective disclosure
	disclosure := types.SelectiveDisclosure{
		CredentialID:   credentialID,
		RevealedClaims: revealedClaims,
		BlindedClaims:  blindedClaims,
		ProofType:      "BBS+",
		ProofValue:     "mock-proof-value", // In real implementation, this would be cryptographically generated
	}
	
	return &disclosure, nil
}

// CheckAnonymitySet verifies the anonymity set size meets requirements
func (k Keeper) CheckAnonymitySet(ctx sdk.Context, credentialType string) error {
	// Get all credentials of this type
	credentials := k.GetCredentialsByType(ctx, credentialType)
	
	// Check anonymity set size
	minSize := k.MinAnonymitySetSize(ctx)
	if uint32(len(credentials)) < minSize {
		return types.ErrAnonymitySetTooSmall
	}
	
	return nil
}

// UpdatePrivacyMetrics updates privacy-related metrics
func (k Keeper) UpdatePrivacyMetrics(ctx sdk.Context, operation string) {
	// This would update metrics for monitoring privacy operations
	// For now, just emit an event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"privacy_operation",
			sdk.NewAttribute("operation", operation),
			sdk.NewAttribute("timestamp", ctx.BlockTime().Format(time.RFC3339)),
		),
	)
}