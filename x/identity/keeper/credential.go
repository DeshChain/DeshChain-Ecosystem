package keeper

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/deshchain/x/identity/types"
)

// SetCredential stores a verifiable credential
func (k Keeper) SetCredential(ctx sdk.Context, credential types.VerifiableCredential) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CredentialPrefix)
	b := k.cdc.MustMarshal(&credential)
	store.Set([]byte(credential.ID), b)
}

// GetCredential retrieves a verifiable credential
func (k Keeper) GetCredential(ctx sdk.Context, credentialID string) (types.VerifiableCredential, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CredentialPrefix)
	b := store.Get([]byte(credentialID))
	if b == nil {
		return types.VerifiableCredential{}, false
	}
	
	var credential types.VerifiableCredential
	k.cdc.MustUnmarshal(b, &credential)
	return credential, true
}

// HasCredential checks if a credential exists
func (k Keeper) HasCredential(ctx sdk.Context, credentialID string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CredentialPrefix)
	return store.Has([]byte(credentialID))
}

// DeleteCredential removes a credential
func (k Keeper) DeleteCredential(ctx sdk.Context, credentialID string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CredentialPrefix)
	store.Delete([]byte(credentialID))
}

// SetCredentialHolderIndex creates an index for credentials by holder
func (k Keeper) SetCredentialHolderIndex(ctx sdk.Context, holder string, credentialID string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CredentialIndexPrefix)
	store.Set(types.GetCredentialByHolderIndexKey(holder, credentialID), []byte{1})
}

// RemoveCredentialHolderIndex removes a credential from holder index
func (k Keeper) RemoveCredentialHolderIndex(ctx sdk.Context, holder string, credentialID string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CredentialIndexPrefix)
	store.Delete(types.GetCredentialByHolderIndexKey(holder, credentialID))
}

// SetCredentialIssuerIndex creates an index for credentials by issuer
func (k Keeper) SetCredentialIssuerIndex(ctx sdk.Context, issuer string, credentialID string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CredentialIndexPrefix)
	store.Set(types.GetCredentialByIssuerIndexKey(issuer, credentialID), []byte{1})
}

// RemoveCredentialIssuerIndex removes a credential from issuer index
func (k Keeper) RemoveCredentialIssuerIndex(ctx sdk.Context, issuer string, credentialID string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CredentialIndexPrefix)
	store.Delete(types.GetCredentialByIssuerIndexKey(issuer, credentialID))
}

// SetCredentialTypeIndex creates an index for credentials by type
func (k Keeper) SetCredentialTypeIndex(ctx sdk.Context, credType string, credentialID string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CredentialIndexPrefix)
	store.Set(types.GetCredentialByTypeIndexKey(credType, credentialID), []byte{1})
}

// RemoveCredentialTypeIndex removes a credential from type index
func (k Keeper) RemoveCredentialTypeIndex(ctx sdk.Context, credType string, credentialID string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CredentialIndexPrefix)
	store.Delete(types.GetCredentialByTypeIndexKey(credType, credentialID))
}

// GetCredentialsByHolder returns all credentials for a holder
func (k Keeper) GetCredentialsByHolder(ctx sdk.Context, holder string) []types.VerifiableCredential {
	var credentials []types.VerifiableCredential
	
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CredentialIndexPrefix)
	iterator := sdk.KVStorePrefixIterator(store, append([]byte("holder:"), []byte(holder)...))
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		// Extract credential ID from the key
		key := string(iterator.Key())
		credentialID := key[len("holder:"+holder):]
		
		if credential, found := k.GetCredential(ctx, credentialID); found {
			credentials = append(credentials, credential)
		}
	}
	
	return credentials
}

// GetCredentialsByIssuer returns all credentials issued by an issuer
func (k Keeper) GetCredentialsByIssuer(ctx sdk.Context, issuer string) []types.VerifiableCredential {
	var credentials []types.VerifiableCredential
	
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CredentialIndexPrefix)
	iterator := sdk.KVStorePrefixIterator(store, append([]byte("issuer:"), []byte(issuer)...))
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		// Extract credential ID from the key
		key := string(iterator.Key())
		credentialID := key[len("issuer:"+issuer):]
		
		if credential, found := k.GetCredential(ctx, credentialID); found {
			credentials = append(credentials, credential)
		}
	}
	
	return credentials
}

// GetCredentialsByType returns all credentials of a specific type
func (k Keeper) GetCredentialsByType(ctx sdk.Context, credType string) []types.VerifiableCredential {
	var credentials []types.VerifiableCredential
	
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CredentialIndexPrefix)
	iterator := sdk.KVStorePrefixIterator(store, append([]byte("type:"), []byte(credType)...))
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		// Extract credential ID from the key
		key := string(iterator.Key())
		credentialID := key[len("type:"+credType):]
		
		if credential, found := k.GetCredential(ctx, credentialID); found {
			credentials = append(credentials, credential)
		}
	}
	
	return credentials
}

// IterateCredentials iterates over all credentials
func (k Keeper) IterateCredentials(ctx sdk.Context, cb func(credential types.VerifiableCredential) bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CredentialPrefix)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var credential types.VerifiableCredential
		k.cdc.MustUnmarshal(iterator.Value(), &credential)
		if cb(credential) {
			break
		}
	}
}

// IssueCredential issues a new verifiable credential
func (k Keeper) IssueCredential(
	ctx sdk.Context,
	issuer string,
	holder string,
	credentialType string,
	claims map[string]interface{},
	expirationDays int32,
) (*types.VerifiableCredential, error) {
	// Validate issuer
	params := k.GetParams(ctx)
	if !params.IsTrustedIssuer(issuer) {
		// Check if issuer is registered
		if !k.IsRegisteredIssuer(ctx, issuer) {
			return nil, types.ErrInvalidIssuer
		}
	}
	
	// Validate holder has identity
	if !k.HasIdentity(ctx, holder) {
		return nil, types.ErrIdentityNotFound
	}
	
	// Generate credential ID
	credentialID := fmt.Sprintf("urn:uuid:%s", sdk.NewRand().Str(32))
	
	// Calculate expiration
	var expirationDate *time.Time
	if expirationDays > 0 {
		expDate := ctx.BlockTime().AddDate(0, 0, int(expirationDays))
		expirationDate = &expDate
	} else {
		// Use default expiration
		expDate := ctx.BlockTime().Add(params.GetCredentialExpiry())
		expirationDate = &expDate
	}
	
	// Create credential subject
	credentialSubject := types.CredentialSubject{
		ID:     holder,
		Claims: claims,
	}
	
	// Create credential
	credential := types.VerifiableCredential{
		Context: []string{
			types.ContextW3CCredentials,
			types.ContextDeshChain,
		},
		ID:   credentialID,
		Type: []string{types.CredentialTypeVerifiable, credentialType},
		Issuer:           issuer,
		IssuanceDate:     ctx.BlockTime(),
		ExpirationDate:   expirationDate,
		CredentialSubject: credentialSubject,
	}
	
	// Validate credential size
	if err := k.ValidateCredentialSize(ctx, &credential); err != nil {
		return nil, err
	}
	
	// Store credential
	k.SetCredential(ctx, credential)
	
	// Create indexes
	k.SetCredentialHolderIndex(ctx, holder, credentialID)
	k.SetCredentialIssuerIndex(ctx, issuer, credentialID)
	k.SetCredentialTypeIndex(ctx, credentialType, credentialID)
	
	// Add to holder's identity
	if err := k.AddIdentityCredential(ctx, holder, credentialID); err != nil {
		return nil, err
	}
	
	return &credential, nil
}

// RevokeCredential adds a credential to the revocation list
func (k Keeper) RevokeCredential(ctx sdk.Context, credentialID string, reason string) error {
	credential, found := k.GetCredential(ctx, credentialID)
	if !found {
		return types.ErrCredentialNotFound
	}
	
	// Add to revocation list
	k.AddToRevocationList(ctx, credential.Issuer, credentialID)
	
	// Update credential status
	if credential.CredentialStatus == nil {
		credential.CredentialStatus = &types.CredentialStatus{}
	}
	credential.CredentialStatus.Type = "RevocationList2020Status"
	credential.CredentialStatus.RevocationListCredential = credential.Issuer
	
	// Update metadata
	if credential.Metadata == nil {
		credential.Metadata = make(map[string]interface{})
	}
	credential.Metadata["revoked"] = true
	credential.Metadata["revokedAt"] = ctx.BlockTime().Format(time.RFC3339)
	credential.Metadata["revocationReason"] = reason
	
	k.SetCredential(ctx, credential)
	
	return nil
}

// VerifyCredential performs basic verification of a credential
func (k Keeper) VerifyCredential(ctx sdk.Context, credentialID string) error {
	credential, found := k.GetCredential(ctx, credentialID)
	if !found {
		return types.ErrCredentialNotFound
	}
	
	// Check if expired
	if credential.IsExpired() {
		return types.ErrCredentialExpired
	}
	
	// Check if revoked
	if k.IsCredentialRevoked(ctx, credentialID) {
		return types.ErrCredentialRevoked
	}
	
	// Validate credential structure
	if err := types.ValidateVerifiableCredential(&credential); err != nil {
		return types.ErrInvalidCredential
	}
	
	// Verify issuer is trusted or registered
	params := k.GetParams(ctx)
	if !params.IsTrustedIssuer(credential.Issuer) {
		if !k.IsRegisteredIssuer(ctx, credential.Issuer) {
			return types.ErrInvalidIssuer
		}
	}
	
	return nil
}

// IsCredentialRevoked checks if a credential is revoked
func (k Keeper) IsCredentialRevoked(ctx sdk.Context, credentialID string) bool {
	credential, found := k.GetCredential(ctx, credentialID)
	if !found {
		return false
	}
	
	// Check metadata for revocation
	if credential.Metadata != nil {
		if revoked, ok := credential.Metadata["revoked"].(bool); ok && revoked {
			return true
		}
	}
	
	// Check revocation list
	return k.IsInRevocationList(ctx, credential.Issuer, credentialID)
}

// SetCredentialSchema stores a credential schema
func (k Keeper) SetCredentialSchema(ctx sdk.Context, schema types.CredentialSchema) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CredentialSchemaPrefix)
	b := k.cdc.MustMarshal(&schema)
	store.Set([]byte(schema.ID), b)
}

// GetCredentialSchema retrieves a credential schema
func (k Keeper) GetCredentialSchema(ctx sdk.Context, schemaID string) (types.CredentialSchema, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CredentialSchemaPrefix)
	b := store.Get([]byte(schemaID))
	if b == nil {
		return types.CredentialSchema{}, false
	}
	
	var schema types.CredentialSchema
	k.cdc.MustUnmarshal(b, &schema)
	return schema, true
}

// IterateCredentialSchemas iterates over all credential schemas
func (k Keeper) IterateCredentialSchemas(ctx sdk.Context, cb func(schema types.CredentialSchema) bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CredentialSchemaPrefix)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var schema types.CredentialSchema
		k.cdc.MustUnmarshal(iterator.Value(), &schema)
		if cb(schema) {
			break
		}
	}
}

// ValidateCredentialSize validates the size of a credential
func (k Keeper) ValidateCredentialSize(ctx sdk.Context, credential *types.VerifiableCredential) error {
	bz, err := credential.ToJSON()
	if err != nil {
		return err
	}
	
	maxSize := k.MaxCredentialSize(ctx)
	if uint64(len(bz)) > maxSize {
		return types.ErrDataTooLarge
	}
	
	return nil
}

// PresentCredential creates a verifiable presentation of credentials
func (k Keeper) PresentCredential(
	ctx sdk.Context,
	holder string,
	credentialIDs []string,
	verifier string,
	challenge string,
	domain string,
) (*types.VerifiablePresentation, error) {
	// Validate holder
	if !k.HasIdentity(ctx, holder) {
		return nil, types.ErrIdentityNotFound
	}
	
	// Validate verifier
	if !k.HasIdentity(ctx, verifier) {
		return nil, types.ErrInvalidAddress
	}
	
	// Collect credentials
	var credentials []interface{}
	for _, credID := range credentialIDs {
		credential, found := k.GetCredential(ctx, credID)
		if !found {
			return nil, types.ErrCredentialNotFound
		}
		
		// Verify holder owns the credential
		subjectID, err := credential.GetSubjectID()
		if err != nil || subjectID != holder {
			return nil, types.ErrUnauthorized
		}
		
		// Verify credential is valid
		if err := k.VerifyCredential(ctx, credID); err != nil {
			return nil, err
		}
		
		credentials = append(credentials, credential)
	}
	
	// Create presentation
	presentation := types.VerifiablePresentation{
		Context: []string{
			types.ContextW3CCredentials,
			types.ContextDeshChain,
		},
		ID:                   fmt.Sprintf("urn:uuid:%s", sdk.NewRand().Str(32)),
		Type:                 []string{"VerifiablePresentation"},
		Holder:               holder,
		VerifiableCredential: credentials,
		Challenge:            challenge,
		Domain:               domain,
		Metadata: map[string]interface{}{
			"verifier":    verifier,
			"presentedAt": ctx.BlockTime().Format(time.RFC3339),
		},
	}
	
	return &presentation, nil
}