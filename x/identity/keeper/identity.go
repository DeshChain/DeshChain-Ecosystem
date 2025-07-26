package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/DeshChain/DeshChain-Ecosystem/x/identity/types"
)

// SetIdentity stores an identity
func (k Keeper) SetIdentity(ctx sdk.Context, identity types.Identity) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.IdentityPrefix)
	b := k.cdc.MustMarshal(&identity)
	store.Set([]byte(identity.Address), b)
}

// GetIdentity retrieves an identity by address
func (k Keeper) GetIdentity(ctx sdk.Context, address string) (types.Identity, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.IdentityPrefix)
	b := store.Get([]byte(address))
	if b == nil {
		return types.Identity{}, false
	}
	
	var identity types.Identity
	k.cdc.MustUnmarshal(b, &identity)
	return identity, true
}

// HasIdentity checks if an identity exists
func (k Keeper) HasIdentity(ctx sdk.Context, address string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.IdentityPrefix)
	return store.Has([]byte(address))
}

// DeleteIdentity removes an identity
func (k Keeper) DeleteIdentity(ctx sdk.Context, address string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.IdentityPrefix)
	store.Delete([]byte(address))
}

// IterateIdentities iterates over all identities
func (k Keeper) IterateIdentities(ctx sdk.Context, cb func(identity types.Identity) bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.IdentityPrefix)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var identity types.Identity
		k.cdc.MustUnmarshal(iterator.Value(), &identity)
		if cb(identity) {
			break
		}
	}
}

// GetAllIdentities returns all identities
func (k Keeper) GetAllIdentities(ctx sdk.Context) []types.Identity {
	var identities []types.Identity
	k.IterateIdentities(ctx, func(identity types.Identity) bool {
		identities = append(identities, identity)
		return false
	})
	return identities
}

// SetIdentityDIDIndex creates an index from DID to address
func (k Keeper) SetIdentityDIDIndex(ctx sdk.Context, did string, address string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.IdentityIndexPrefix)
	store.Set(types.GetIdentityByDIDIndexKey(did), []byte(address))
}

// GetIdentityByDID retrieves an identity by DID
func (k Keeper) GetIdentityByDID(ctx sdk.Context, did string) (types.Identity, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.IdentityIndexPrefix)
	addressBz := store.Get(types.GetIdentityByDIDIndexKey(did))
	if addressBz == nil {
		return types.Identity{}, false
	}
	
	return k.GetIdentity(ctx, string(addressBz))
}

// DeleteIdentityDIDIndex removes the DID to address index
func (k Keeper) DeleteIdentityDIDIndex(ctx sdk.Context, did string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.IdentityIndexPrefix)
	store.Delete(types.GetIdentityByDIDIndexKey(did))
}

// UpdateIdentityStatus updates the status of an identity
func (k Keeper) UpdateIdentityStatus(ctx sdk.Context, address string, status types.IdentityStatus) error {
	identity, found := k.GetIdentity(ctx, address)
	if !found {
		return types.ErrIdentityNotFound
	}
	
	identity.Status = status
	identity.UpdatedAt = ctx.BlockTime()
	k.SetIdentity(ctx, identity)
	
	return nil
}

// UpdateIdentityKYCStatus updates the KYC status of an identity
func (k Keeper) UpdateIdentityKYCStatus(ctx sdk.Context, address string, kycStatus types.KYCStatus) error {
	identity, found := k.GetIdentity(ctx, address)
	if !found {
		return types.ErrIdentityNotFound
	}
	
	identity.KYCStatus = kycStatus
	identity.UpdatedAt = ctx.BlockTime()
	k.SetIdentity(ctx, identity)
	
	return nil
}

// UpdateIdentityBiometricStatus updates the biometric status of an identity
func (k Keeper) UpdateIdentityBiometricStatus(ctx sdk.Context, address string, biometricStatus types.BiometricStatus) error {
	identity, found := k.GetIdentity(ctx, address)
	if !found {
		return types.ErrIdentityNotFound
	}
	
	identity.BiometricStatus = biometricStatus
	identity.UpdatedAt = ctx.BlockTime()
	k.SetIdentity(ctx, identity)
	
	return nil
}

// AddIdentityCredential adds a credential ID to an identity
func (k Keeper) AddIdentityCredential(ctx sdk.Context, address string, credentialID string) error {
	identity, found := k.GetIdentity(ctx, address)
	if !found {
		return types.ErrIdentityNotFound
	}
	
	// Check if credential already exists
	for _, cid := range identity.Credentials {
		if cid == credentialID {
			return nil // Already exists
		}
	}
	
	// Check max credentials limit
	maxCreds := k.MaxCredentialsPerIdentity(ctx)
	if uint32(len(identity.Credentials)) >= maxCreds {
		return types.ErrInvalidRequest
	}
	
	identity.Credentials = append(identity.Credentials, credentialID)
	identity.UpdatedAt = ctx.BlockTime()
	k.SetIdentity(ctx, identity)
	
	return nil
}

// RemoveIdentityCredential removes a credential ID from an identity
func (k Keeper) RemoveIdentityCredential(ctx sdk.Context, address string, credentialID string) error {
	identity, found := k.GetIdentity(ctx, address)
	if !found {
		return types.ErrIdentityNotFound
	}
	
	// Find and remove credential
	newCredentials := []string{}
	found = false
	for _, cid := range identity.Credentials {
		if cid != credentialID {
			newCredentials = append(newCredentials, cid)
		} else {
			found = true
		}
	}
	
	if !found {
		return types.ErrCredentialNotFound
	}
	
	identity.Credentials = newCredentials
	identity.UpdatedAt = ctx.BlockTime()
	k.SetIdentity(ctx, identity)
	
	return nil
}

// AddIdentityConsent adds a consent record to an identity
func (k Keeper) AddIdentityConsent(ctx sdk.Context, address string, consent types.ConsentRecord) error {
	identity, found := k.GetIdentity(ctx, address)
	if !found {
		return types.ErrIdentityNotFound
	}
	
	identity.Consents = append(identity.Consents, consent)
	identity.UpdatedAt = ctx.BlockTime()
	k.SetIdentity(ctx, identity)
	
	return nil
}

// WithdrawIdentityConsent withdraws a consent
func (k Keeper) WithdrawIdentityConsent(ctx sdk.Context, address string, consentID string) error {
	identity, found := k.GetIdentity(ctx, address)
	if !found {
		return types.ErrIdentityNotFound
	}
	
	consentFound := false
	for i, consent := range identity.Consents {
		if consent.ID == consentID {
			withdrawnAt := ctx.BlockTime()
			identity.Consents[i].Given = false
			identity.Consents[i].WithdrawnAt = &withdrawnAt
			consentFound = true
			break
		}
	}
	
	if !consentFound {
		return types.ErrInvalidRequest
	}
	
	identity.UpdatedAt = ctx.BlockTime()
	k.SetIdentity(ctx, identity)
	
	return nil
}

// AddRecoveryMethod adds a recovery method to an identity
func (k Keeper) AddRecoveryMethod(ctx sdk.Context, address string, method types.RecoveryMethod) error {
	identity, found := k.GetIdentity(ctx, address)
	if !found {
		return types.ErrIdentityNotFound
	}
	
	// Check max recovery methods
	maxMethods := k.MaxRecoveryMethods(ctx)
	if uint32(len(identity.RecoveryMethods)) >= maxMethods {
		return types.ErrInvalidRequest
	}
	
	identity.RecoveryMethods = append(identity.RecoveryMethods, method)
	identity.UpdatedAt = ctx.BlockTime()
	k.SetIdentity(ctx, identity)
	
	return nil
}

// RemoveRecoveryMethod removes a recovery method from an identity
func (k Keeper) RemoveRecoveryMethod(ctx sdk.Context, address string, methodType types.RecoveryType) error {
	identity, found := k.GetIdentity(ctx, address)
	if !found {
		return types.ErrIdentityNotFound
	}
	
	newMethods := []types.RecoveryMethod{}
	found = false
	for _, method := range identity.RecoveryMethods {
		if method.Type != methodType {
			newMethods = append(newMethods, method)
		} else {
			found = true
		}
	}
	
	if !found {
		return types.ErrRecoveryMethodNotSet
	}
	
	identity.RecoveryMethods = newMethods
	identity.UpdatedAt = ctx.BlockTime()
	k.SetIdentity(ctx, identity)
	
	return nil
}

// ValidateIdentityCreation validates if an identity can be created
func (k Keeper) ValidateIdentityCreation(ctx sdk.Context, address string) error {
	// Check if identity already exists
	if k.HasIdentity(ctx, address) {
		return types.ErrIdentityAlreadyExists
	}
	
	// Check if account exists
	acc := k.accountKeeper.GetAccount(ctx, sdk.MustAccAddressFromBech32(address))
	if acc == nil {
		return types.ErrInvalidAddress
	}
	
	return nil
}