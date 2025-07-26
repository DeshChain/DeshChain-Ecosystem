package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/namo/x/identity/types"
)

// CachedKeeper wraps the identity keeper with caching capabilities
type CachedKeeper struct {
	*Keeper
	cache types.CacheInterface
}

// NewCachedKeeper creates a new cached keeper instance
func NewCachedKeeper(keeper *Keeper, cacheConfig types.CacheConfig) *CachedKeeper {
	cache := NewIdentityCache(cacheConfig)
	
	ck := &CachedKeeper{
		Keeper: keeper,
		cache:  cache,
	}
	
	// Start the cache
	if err := cache.Start(); err != nil {
		panic(fmt.Sprintf("failed to start cache: %v", err))
	}
	
	return ck
}

// GetCache returns the cache instance
func (ck *CachedKeeper) GetCache() types.CacheInterface {
	return ck.cache
}

// Cached Identity Operations

// GetIdentity retrieves an identity with caching
func (ck *CachedKeeper) GetIdentityCached(ctx sdk.Context, address string) (types.Identity, bool) {
	// Try cache first
	cacheKey := types.NewCacheKey(types.CacheTypeIdentity, address)
	if entry, found := ck.cache.Get(cacheKey); found {
		if identity, ok := entry.Data.(types.Identity); ok {
			return identity, true
		}
	}
	
	// Cache miss - load from store
	identity, found := ck.Keeper.GetIdentity(ctx, address)
	if found {
		// Cache the identity with type-specific TTL
		ttl := ck.getIdentityCacheConfig().IdentityTTL
		if err := ck.cache.Set(cacheKey, identity, ttl); err != nil {
			// Log error but don't fail the operation
			ctx.Logger().Error("Failed to cache identity", "address", address, "error", err)
		}
	}
	
	return identity, found
}

// SetIdentity stores an identity and updates cache
func (ck *CachedKeeper) SetIdentityCached(ctx sdk.Context, identity types.Identity) {
	// Update store first
	ck.Keeper.SetIdentity(ctx, identity)
	
	// Update cache
	cacheKey := types.NewCacheKey(types.CacheTypeIdentity, identity.Address)
	ttl := ck.getIdentityCacheConfig().IdentityTTL
	if err := ck.cache.Set(cacheKey, identity, ttl); err != nil {
		ctx.Logger().Error("Failed to cache identity after set", "address", identity.Address, "error", err)
	}
	
	// Create enriched cache entry for frequent lookups
	enrichedEntry := ck.createEnrichedIdentityEntry(ctx, identity)
	enrichedKey := types.NewCacheKey("enriched_identity", identity.Address)
	if err := ck.cache.Set(enrichedKey, enrichedEntry, ttl); err != nil {
		ctx.Logger().Error("Failed to cache enriched identity", "address", identity.Address, "error", err)
	}
}

// DeleteIdentity removes an identity and invalidates cache
func (ck *CachedKeeper) DeleteIdentityCached(ctx sdk.Context, address string) {
	// Remove from store
	ck.Keeper.DeleteIdentity(ctx, address)
	
	// Remove from cache
	cacheKey := types.NewCacheKey(types.CacheTypeIdentity, address)
	if err := ck.cache.Delete(cacheKey); err != nil {
		ctx.Logger().Error("Failed to delete identity from cache", "address", address, "error", err)
	}
	
	// Remove enriched entry
	enrichedKey := types.NewCacheKey("enriched_identity", address)
	ck.cache.Delete(enrichedKey)
	
	// Invalidate related cache entries
	ck.invalidateRelatedEntries(ctx, address)
}

// Cached Credential Operations

// GetCredential retrieves a credential with caching
func (ck *CachedKeeper) GetCredentialCached(ctx sdk.Context, credentialID string) (types.Credential, bool) {
	// Try cache first
	cacheKey := types.NewCacheKey(types.CacheTypeCredential, credentialID)
	if entry, found := ck.cache.Get(cacheKey); found {
		if credential, ok := entry.Data.(types.Credential); ok {
			return credential, true
		}
	}
	
	// Cache miss - load from store
	credential, found := ck.Keeper.GetCredential(ctx, credentialID)
	if found {
		// Cache the credential
		ttl := ck.getIdentityCacheConfig().CredentialTTL
		if err := ck.cache.Set(cacheKey, credential, ttl); err != nil {
			ctx.Logger().Error("Failed to cache credential", "id", credentialID, "error", err)
		}
	}
	
	return credential, found
}

// SetCredential stores a credential and updates cache
func (ck *CachedKeeper) SetCredentialCached(ctx sdk.Context, credential types.Credential) {
	// Update store first
	ck.Keeper.SetCredential(ctx, credential)
	
	// Update cache
	cacheKey := types.NewCacheKey(types.CacheTypeCredential, credential.ID)
	ttl := ck.getIdentityCacheConfig().CredentialTTL
	if err := ck.cache.Set(cacheKey, credential, ttl); err != nil {
		ctx.Logger().Error("Failed to cache credential after set", "id", credential.ID, "error", err)
	}
	
	// Invalidate related entries by holder
	ck.invalidateByTag(fmt.Sprintf("holder:%s", credential.Holder))
}

// GetCredentialsByHolder retrieves credentials by holder with caching
func (ck *CachedKeeper) GetCredentialsByHolderCached(ctx sdk.Context, holder string) []types.Credential {
	// Try to get from cache using tag
	tag := fmt.Sprintf("holder:%s", holder)
	entries := ck.cache.GetByTag(tag)
	
	if len(entries) > 0 {
		var credentials []types.Credential
		for _, entry := range entries {
			if credential, ok := entry.Data.(types.Credential); ok {
				credentials = append(credentials, credential)
			}
		}
		if len(credentials) > 0 {
			return credentials
		}
	}
	
	// Cache miss - load from store
	credentials := ck.Keeper.GetCredentialsByHolder(ctx, holder)
	
	// Cache individual credentials with holder tag
	for _, credential := range credentials {
		cacheKey := types.NewCacheKey(types.CacheTypeCredential, credential.ID)
		ttl := ck.getIdentityCacheConfig().CredentialTTL
		
		// Create cache entry with holder tag
		entry := types.NewCacheEntry(cacheKey, credential, ttl)
		entry.AddTag(tag)
		
		if err := ck.cache.Set(cacheKey, credential, ttl); err != nil {
			ctx.Logger().Error("Failed to cache credential for holder", "holder", holder, "credentialID", credential.ID, "error", err)
		}
	}
	
	return credentials
}

// Cached DID Operations

// GetDIDDocument retrieves a DID document with caching
func (ck *CachedKeeper) GetDIDDocumentCached(ctx sdk.Context, did string) (types.DIDDocument, bool) {
	// Try cache first
	cacheKey := types.NewCacheKey(types.CacheTypeDIDDocument, did)
	if entry, found := ck.cache.Get(cacheKey); found {
		if didDoc, ok := entry.Data.(types.DIDDocument); ok {
			return didDoc, true
		}
	}
	
	// Cache miss - load from store
	didDoc, found := ck.Keeper.GetDIDDocument(ctx, did)
	if found {
		// Cache the DID document
		ttl := ck.getIdentityCacheConfig().DIDDocumentTTL
		if err := ck.cache.Set(cacheKey, didDoc, ttl); err != nil {
			ctx.Logger().Error("Failed to cache DID document", "did", did, "error", err)
		}
	}
	
	return didDoc, found
}

// SetDIDDocument stores a DID document and updates cache
func (ck *CachedKeeper) SetDIDDocumentCached(ctx sdk.Context, didDoc types.DIDDocument) {
	// Update store first
	ck.Keeper.SetDIDDocument(ctx, didDoc)
	
	// Update cache
	cacheKey := types.NewCacheKey(types.CacheTypeDIDDocument, didDoc.ID)
	ttl := ck.getIdentityCacheConfig().DIDDocumentTTL
	if err := ck.cache.Set(cacheKey, didDoc, ttl); err != nil {
		ctx.Logger().Error("Failed to cache DID document after set", "did", didDoc.ID, "error", err)
	}
}

// Cached Consent Operations

// GetConsent retrieves a consent record with caching
func (ck *CachedKeeper) GetConsentCached(ctx sdk.Context, consentID string) (types.ConsentRecord, bool) {
	// Try cache first
	cacheKey := types.NewCacheKey(types.CacheTypeConsent, consentID)
	if entry, found := ck.cache.Get(cacheKey); found {
		if consent, ok := entry.Data.(types.ConsentRecord); ok {
			return consent, true
		}
	}
	
	// Cache miss - load from store
	consent, found := ck.Keeper.GetConsentRecord(ctx, consentID)
	if found {
		// Cache the consent
		ttl := ck.getIdentityCacheConfig().ConsentTTL
		if err := ck.cache.Set(cacheKey, consent, ttl); err != nil {
			ctx.Logger().Error("Failed to cache consent", "id", consentID, "error", err)
		}
	}
	
	return consent, found
}

// Cached ZK Proof Operations

// GetZKProof retrieves a ZK proof with caching
func (ck *CachedKeeper) GetZKProofCached(ctx sdk.Context, proofID string) (types.ZKProof, bool) {
	// Try cache first
	cacheKey := types.NewCacheKey(types.CacheTypeZKProof, proofID)
	if entry, found := ck.cache.Get(cacheKey); found {
		if proof, ok := entry.Data.(types.ZKProof); ok {
			return proof, true
		}
	}
	
	// Cache miss - load from store
	proof, found := ck.Keeper.GetZKProof(ctx, proofID)
	if found {
		// Cache the proof
		ttl := ck.getIdentityCacheConfig().ZKProofTTL
		if err := ck.cache.Set(cacheKey, proof, ttl); err != nil {
			ctx.Logger().Error("Failed to cache ZK proof", "id", proofID, "error", err)
		}
	}
	
	return proof, found
}

// Cached Share Operations

// GetShareRequest retrieves a share request with caching
func (ck *CachedKeeper) GetShareRequestCached(ctx sdk.Context, requestID string) (types.IdentityShareRequest, bool) {
	// Try cache first
	cacheKey := types.NewCacheKey(types.CacheTypeShareRequest, requestID)
	if entry, found := ck.cache.Get(cacheKey); found {
		if request, ok := entry.Data.(types.IdentityShareRequest); ok {
			return request, true
		}
	}
	
	// Cache miss - load from store
	request, found := ck.Keeper.GetShareRequest(ctx, requestID)
	if found {
		// Cache with shorter TTL for share requests
		ttl := 15 * time.Minute
		if err := ck.cache.Set(cacheKey, request, ttl); err != nil {
			ctx.Logger().Error("Failed to cache share request", "id", requestID, "error", err)
		}
	}
	
	return request, found
}

// Cache Management Operations

// InvalidateUserCache invalidates all cache entries for a specific user
func (ck *CachedKeeper) InvalidateUserCache(ctx sdk.Context, address string) {
	// Get user's DID first
	if identity, found := ck.Keeper.GetIdentity(ctx, address); found {
		// Invalidate by DID tag
		ck.invalidateByTag(fmt.Sprintf("did:%s", identity.DID))
	}
	
	// Invalidate by address tag
	ck.invalidateByTag(fmt.Sprintf("address:%s", address))
	
	// Invalidate specific entries
	identityKey := types.NewCacheKey(types.CacheTypeIdentity, address)
	ck.cache.Delete(identityKey)
	
	enrichedKey := types.NewCacheKey("enriched_identity", address)
	ck.cache.Delete(enrichedKey)
}

// RefreshCache performs a full cache refresh
func (ck *CachedKeeper) RefreshCache(ctx sdk.Context) error {
	// Clear the cache
	if err := ck.cache.Clear(); err != nil {
		return err
	}
	
	// Preload active identities if configured
	config := ck.getIdentityCacheConfig()
	if config.PreloadIdentities {
		ck.preloadActiveIdentities(ctx)
	}
	
	return nil
}

// GetCacheStats returns cache performance statistics
func (ck *CachedKeeper) GetCacheStats() types.CacheStats {
	return ck.cache.Stats()
}

// GetCacheMetrics returns detailed cache metrics
func (ck *CachedKeeper) GetCacheMetrics() types.CacheMetrics {
	if identityCache, ok := ck.cache.(*IdentityCache); ok {
		return identityCache.GetMetrics()
	}
	return types.CacheMetrics{}
}

// Helper Methods

func (ck *CachedKeeper) getIdentityCacheConfig() types.CacheConfig {
	if identityCache, ok := ck.cache.(*IdentityCache); ok {
		return identityCache.config
	}
	return types.DefaultCacheConfig
}

func (ck *CachedKeeper) createEnrichedIdentityEntry(ctx sdk.Context, identity types.Identity) types.IdentityCacheEntry {
	// Get DID document
	var didDoc *types.DIDDocument
	if doc, found := ck.Keeper.GetDIDDocument(ctx, identity.DID); found {
		didDoc = &doc
	}
	
	// Get active credentials
	credentials := ck.Keeper.GetCredentialsByHolder(ctx, identity.Address)
	var credentialPtrs []*types.Credential
	for i := range credentials {
		credentialPtrs = append(credentialPtrs, &credentials[i])
	}
	
	// Get active consents
	consents := ck.Keeper.GetConsentsByHolder(ctx, identity.Address)
	var consentPtrs []*types.ConsentRecord
	for i := range consents {
		consentPtrs = append(consentPtrs, &consents[i])
	}
	
	// Calculate computed fields
	computedFields := map[string]interface{}{
		"credential_count": len(credentials),
		"consent_count":    len(consents),
		"kyc_level":        ck.calculateKYCLevel(credentials),
		"last_updated":     time.Now(),
	}
	
	return types.IdentityCacheEntry{
		Identity:        &identity,
		DIDDocument:     didDoc,
		Credentials:     credentialPtrs,
		ActiveConsents:  consentPtrs,
		LastActivity:    identity.LastActivityAt,
		AccessFrequency: 1,
		ComputedFields:  computedFields,
	}
}

func (ck *CachedKeeper) calculateKYCLevel(credentials []types.Credential) string {
	// Simple KYC level calculation based on credential types
	hasBasicKYC := false
	hasAdvancedKYC := false
	
	for _, cred := range credentials {
		for _, credType := range cred.Type {
			switch credType {
			case "BasicKYCCredential", "AadhaarCredential":
				hasBasicKYC = true
			case "AdvancedKYCCredential", "BiometricCredential":
				hasAdvancedKYC = true
			}
		}
	}
	
	if hasAdvancedKYC {
		return "advanced"
	} else if hasBasicKYC {
		return "basic"
	}
	return "none"
}

func (ck *CachedKeeper) invalidateRelatedEntries(ctx sdk.Context, address string) {
	// Get identity to find DID
	if identity, found := ck.Keeper.GetIdentity(ctx, address); found {
		// Invalidate by DID
		ck.invalidateByTag(fmt.Sprintf("did:%s", identity.DID))
	}
	
	// Invalidate credentials by holder
	ck.invalidateByTag(fmt.Sprintf("holder:%s", address))
	
	// Invalidate consents by holder
	ck.invalidateByTag(fmt.Sprintf("consent_holder:%s", address))
}

func (ck *CachedKeeper) invalidateByTag(tag string) {
	if err := ck.cache.InvalidateByTag(tag); err != nil {
		// Log error but don't fail
		// ctx.Logger().Error("Failed to invalidate cache by tag", "tag", tag, "error", err)
	}
}

func (ck *CachedKeeper) preloadActiveIdentities(ctx sdk.Context) {
	// This would typically load the most frequently accessed identities
	// Implementation depends on usage patterns and available indices
	
	// For now, this is a placeholder
	// In production, you might:
	// 1. Load identities with recent activity
	// 2. Load identities with high credential counts
	// 3. Load identities frequently used in cross-module sharing
}

// Background Cache Warming
func (ck *CachedKeeper) StartCacheWarming(ctx sdk.Context) {
	// Start a background goroutine for cache warming
	go ck.cacheWarmingRoutine(ctx)
}

func (ck *CachedKeeper) cacheWarmingRoutine(ctx sdk.Context) {
	ticker := time.NewTicker(30 * time.Minute) // Warm every 30 minutes
	defer ticker.Stop()
	
	for range ticker.C {
		ck.preloadActiveIdentities(ctx)
	}
}