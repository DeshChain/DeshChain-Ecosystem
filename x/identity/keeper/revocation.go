package keeper

import (
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/deshchain/x/identity/types"
)

// SetRevocationList stores a revocation list
func (k Keeper) SetRevocationList(ctx sdk.Context, revList types.RevocationList) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RevocationListPrefix)
	b := k.cdc.MustMarshal(&revList)
	store.Set([]byte(revList.Issuer), b)
}

// GetRevocationList retrieves a revocation list for an issuer
func (k Keeper) GetRevocationList(ctx sdk.Context, issuer string) (types.RevocationList, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RevocationListPrefix)
	b := store.Get([]byte(issuer))
	if b == nil {
		return types.RevocationList{}, false
	}
	
	var revList types.RevocationList
	k.cdc.MustUnmarshal(b, &revList)
	return revList, true
}

// AddToRevocationList adds a credential to the revocation list
func (k Keeper) AddToRevocationList(ctx sdk.Context, issuer string, credentialID string) {
	revList, found := k.GetRevocationList(ctx, issuer)
	if !found {
		revList = types.RevocationList{
			ID:                 issuer + "-revocation-list",
			Issuer:             issuer,
			RevokedCredentials: []string{},
			LastUpdated:        ctx.BlockTime().Format(time.RFC3339),
		}
	}
	
	// Check if already revoked
	for _, revokedID := range revList.RevokedCredentials {
		if revokedID == credentialID {
			return // Already revoked
		}
	}
	
	// Add to revoked list
	revList.RevokedCredentials = append(revList.RevokedCredentials, credentialID)
	revList.LastUpdated = ctx.BlockTime().Format(time.RFC3339)
	
	k.SetRevocationList(ctx, revList)
}

// RemoveFromRevocationList removes a credential from the revocation list
func (k Keeper) RemoveFromRevocationList(ctx sdk.Context, issuer string, credentialID string) {
	revList, found := k.GetRevocationList(ctx, issuer)
	if !found {
		return
	}
	
	// Remove from list
	newRevoked := []string{}
	for _, revokedID := range revList.RevokedCredentials {
		if revokedID != credentialID {
			newRevoked = append(newRevoked, revokedID)
		}
	}
	
	revList.RevokedCredentials = newRevoked
	revList.LastUpdated = ctx.BlockTime().Format(time.RFC3339)
	
	k.SetRevocationList(ctx, revList)
}

// IsInRevocationList checks if a credential is in the revocation list
func (k Keeper) IsInRevocationList(ctx sdk.Context, issuer string, credentialID string) bool {
	revList, found := k.GetRevocationList(ctx, issuer)
	if !found {
		return false
	}
	
	for _, revokedID := range revList.RevokedCredentials {
		if revokedID == credentialID {
			return true
		}
	}
	
	return false
}

// IterateRevocationLists iterates over all revocation lists
func (k Keeper) IterateRevocationLists(ctx sdk.Context, cb func(revList types.RevocationList) bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RevocationListPrefix)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var revList types.RevocationList
		k.cdc.MustUnmarshal(iterator.Value(), &revList)
		if cb(revList) {
			break
		}
	}
}

// GetRevokedCredentialCount returns the number of revoked credentials for an issuer
func (k Keeper) GetRevokedCredentialCount(ctx sdk.Context, issuer string) int {
	revList, found := k.GetRevocationList(ctx, issuer)
	if !found {
		return 0
	}
	
	return len(revList.RevokedCredentials)
}

// CleanupRevocationList removes non-existent credentials from revocation list
func (k Keeper) CleanupRevocationList(ctx sdk.Context, issuer string) {
	revList, found := k.GetRevocationList(ctx, issuer)
	if !found {
		return
	}
	
	// Check each revoked credential still exists
	validRevoked := []string{}
	for _, credentialID := range revList.RevokedCredentials {
		if k.HasCredential(ctx, credentialID) {
			validRevoked = append(validRevoked, credentialID)
		}
	}
	
	if len(validRevoked) != len(revList.RevokedCredentials) {
		revList.RevokedCredentials = validRevoked
		revList.LastUpdated = ctx.BlockTime().Format(time.RFC3339)
		k.SetRevocationList(ctx, revList)
	}
}

// GenerateRevocationListCredential creates a verifiable credential for the revocation list
func (k Keeper) GenerateRevocationListCredential(ctx sdk.Context, issuer string) (*types.VerifiableCredential, error) {
	revList, found := k.GetRevocationList(ctx, issuer)
	if !found {
		return nil, types.ErrInvalidRequest
	}
	
	// Create revocation list as credential subject
	credentialSubject := map[string]interface{}{
		"id":                 revList.ID,
		"type":               "RevocationList2020",
		"issuer":             revList.Issuer,
		"revokedCredentials": revList.RevokedCredentials,
		"lastUpdated":        revList.LastUpdated,
	}
	
	// Create the revocation list credential
	credential := types.VerifiableCredential{
		Context: []string{
			types.ContextW3CCredentials,
			"https://w3id.org/vc-status-list-2020/v1",
		},
		ID:                revList.ID,
		Type:              []string{"VerifiableCredential", "StatusList2020Credential"},
		Issuer:            issuer,
		IssuanceDate:      ctx.BlockTime(),
		CredentialSubject: credentialSubject,
	}
	
	return &credential, nil
}

// BatchRevokeCredentials revokes multiple credentials at once
func (k Keeper) BatchRevokeCredentials(ctx sdk.Context, issuer string, credentialIDs []string, reason string) error {
	// Validate issuer
	if !k.IsRegisteredIssuer(ctx, issuer) {
		params := k.GetParams(ctx)
		if !params.IsTrustedIssuer(issuer) {
			return types.ErrInvalidIssuer
		}
	}
	
	// Revoke each credential
	for _, credentialID := range credentialIDs {
		if err := k.RevokeCredential(ctx, credentialID, reason); err != nil {
			// Log error but continue with other revocations
			k.Logger(ctx).Error("failed to revoke credential", "credential_id", credentialID, "error", err)
		}
	}
	
	return nil
}