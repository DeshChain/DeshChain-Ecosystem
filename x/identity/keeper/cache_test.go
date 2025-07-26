package keeper_test

import (
	"fmt"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/DeshChain/DeshChain-Ecosystem/x/identity/keeper"
	"github.com/DeshChain/DeshChain-Ecosystem/x/identity/types"
	"github.com/DeshChain/DeshChain-Ecosystem/testutil"
)

type CacheTestSuite struct {
	suite.Suite

	ctx           sdk.Context
	keeper        *keeper.Keeper
	cachedKeeper  *keeper.CachedKeeper
	cache         types.CacheInterface
	addrs         []sdk.AccAddress
}

func (suite *CacheTestSuite) SetupTest() {
	suite.ctx, suite.keeper = testutil.IdentityKeeperTestSetup(suite.T())
	suite.addrs = testutil.CreateIncrementalAccounts(10)
	
	// Create cached keeper with test configuration
	cacheConfig := types.CacheConfig{
		MaxSize:         1024 * 1024, // 1MB
		MaxEntries:      1000,
		DefaultTTL:      5 * time.Minute,
		CleanupInterval: 30 * time.Second,
		EnableMetrics:   true,
		EnableTags:      true,
		IdentityTTL:     10 * time.Minute,
		CredentialTTL:   5 * time.Minute,
		DIDDocumentTTL:  15 * time.Minute,
		ConsentTTL:      3 * time.Minute,
		ZKProofTTL:      2 * time.Minute,
		PreloadIdentities: false,
		WarmupOnStart:     false,
		AsyncEviction:     false,
	}
	
	suite.cachedKeeper = keeper.NewCachedKeeper(suite.keeper, cacheConfig)
	suite.cache = suite.cachedKeeper.GetCache()
}

func (suite *CacheTestSuite) TearDownTest() {
	if suite.cache != nil {
		suite.cache.Stop()
	}
}

func TestCacheTestSuite(t *testing.T) {
	suite.Run(t, new(CacheTestSuite))
}

// TestBasicCacheOperations tests basic cache functionality
func (suite *CacheTestSuite) TestBasicCacheOperations() {
	cache := suite.cache
	
	// Test Set and Get
	key := types.NewCacheKey("test", "key1")
	data := "test data"
	ttl := 1 * time.Minute
	
	err := cache.Set(key, data, ttl)
	suite.NoError(err)
	
	entry, found := cache.Get(key)
	suite.True(found)
	suite.NotNil(entry)
	suite.Equal(data, entry.Data)
	
	// Test Exists
	suite.True(cache.Exists(key))
	
	// Test Delete
	err = cache.Delete(key)
	suite.NoError(err)
	
	_, found = cache.Get(key)
	suite.False(found)
	suite.False(cache.Exists(key))
}

// TestCacheExpiration tests TTL and expiration
func (suite *CacheTestSuite) TestCacheExpiration() {
	cache := suite.cache
	
	// Set entry with very short TTL
	key := types.NewCacheKey("test", "expiry_test")
	data := "expiring data"
	ttl := 100 * time.Millisecond
	
	err := cache.Set(key, data, ttl)
	suite.NoError(err)
	
	// Should be available immediately
	_, found := cache.Get(key)
	suite.True(found)
	
	// Wait for expiration
	time.Sleep(150 * time.Millisecond)
	
	// Should be expired now
	_, found = cache.Get(key)
	suite.False(found)
}

// TestCacheStats tests cache statistics
func (suite *CacheTestSuite) TestCacheStats() {
	cache := suite.cache
	
	// Get initial stats
	initialStats := cache.Stats()
	
	// Add some entries
	for i := 0; i < 5; i++ {
		key := types.NewCacheKey("test", fmt.Sprintf("key%d", i))
		cache.Set(key, fmt.Sprintf("data%d", i), 1*time.Minute)
	}
	
	// Get entries (hits)
	for i := 0; i < 3; i++ {
		key := types.NewCacheKey("test", fmt.Sprintf("key%d", i))
		cache.Get(key)
	}
	
	// Try to get non-existent entries (misses)
	for i := 5; i < 8; i++ {
		key := types.NewCacheKey("test", fmt.Sprintf("key%d", i))
		cache.Get(key)
	}
	
	// Check stats
	stats := cache.Stats()
	suite.Equal(initialStats.TotalEntries+5, stats.TotalEntries)
	suite.Equal(initialStats.CacheHits+3, stats.CacheHits)
	suite.Equal(initialStats.CacheMisses+3, stats.CacheMisses)
	
	hitRatio := stats.GetHitRatio()
	suite.True(hitRatio > 0 && hitRatio <= 1)
}

// TestTagOperations tests tag-based operations
func (suite *CacheTestSuite) TestTagOperations() {
	cache := suite.cache
	
	// Create entries with tags
	key1 := types.NewCacheKey("test", "tagged1")
	entry1 := types.NewCacheEntry(key1, "data1", 1*time.Minute)
	entry1.AddTag("tag1")
	entry1.AddTag("common")
	cache.Set(key1, entry1.Data, 1*time.Minute)
	
	key2 := types.NewCacheKey("test", "tagged2")
	entry2 := types.NewCacheEntry(key2, "data2", 1*time.Minute)
	entry2.AddTag("tag2")
	entry2.AddTag("common")
	cache.Set(key2, entry2.Data, 1*time.Minute)
	
	key3 := types.NewCacheKey("test", "untagged")
	cache.Set(key3, "data3", 1*time.Minute)
	
	// Test GetByTag
	commonEntries := cache.GetByTag("common")
	suite.Len(commonEntries, 2)
	
	tag1Entries := cache.GetByTag("tag1")
	suite.Len(tag1Entries, 1)
	
	// Test InvalidateByTag
	err := cache.InvalidateByTag("common")
	suite.NoError(err)
	
	// Tagged entries should be gone
	_, found1 := cache.Get(key1)
	_, found2 := cache.Get(key2)
	_, found3 := cache.Get(key3)
	
	suite.False(found1)
	suite.False(found2)
	suite.True(found3) // Untagged entry should remain
}

// TestMultipleOperations tests batch operations
func (suite *CacheTestSuite) TestMultipleOperations() {
	cache := suite.cache
	
	// Test SetMultiple
	entries := make(map[types.CacheKey]interface{})
	keys := make([]types.CacheKey, 3)
	
	for i := 0; i < 3; i++ {
		key := types.NewCacheKey("test", fmt.Sprintf("multi%d", i))
		keys[i] = key
		entries[key] = fmt.Sprintf("data%d", i)
	}
	
	err := cache.SetMultiple(entries, 1*time.Minute)
	suite.NoError(err)
	
	// Test GetMultiple
	result := cache.GetMultiple(keys)
	suite.Len(result, 3)
	
	// Test DeleteMultiple
	err = cache.DeleteMultiple(keys)
	suite.NoError(err)
	
	// Entries should be gone
	for _, key := range keys {
		_, found := cache.Get(key)
		suite.False(found)
	}
}

// TestCachedIdentityOperations tests cached identity operations
func (suite *CacheTestSuite) TestCachedIdentityOperations() {
	ctx := suite.ctx
	ck := suite.cachedKeeper
	address := suite.addrs[0].String()
	
	// Create an identity
	identity := types.Identity{
		Address:        address,
		DID:            fmt.Sprintf("did:desh:test_%s", address),
		Status:         types.IdentityStatus_ACTIVE,
		CreatedAt:      ctx.BlockTime(),
		UpdatedAt:      ctx.BlockTime(),
		LastActivityAt: ctx.BlockTime(),
	}
	
	// Set identity (should cache it)
	ck.SetIdentityCached(ctx, identity)
	
	// Get identity - should hit cache
	retrieved, found := ck.GetIdentityCached(ctx, address)
	suite.True(found)
	suite.Equal(identity.Address, retrieved.Address)
	suite.Equal(identity.DID, retrieved.DID)
	
	// Verify cache hit
	stats := ck.GetCacheStats()
	suite.True(stats.CacheHits > 0)
	
	// Update identity
	identity.Status = types.IdentityStatus_SUSPENDED
	ck.SetIdentityCached(ctx, identity)
	
	// Get updated identity
	updated, found := ck.GetIdentityCached(ctx, address)
	suite.True(found)
	suite.Equal(types.IdentityStatus_SUSPENDED, updated.Status)
	
	// Delete identity
	ck.DeleteIdentityCached(ctx, address)
	
	// Should not be found in cache
	_, found = ck.GetIdentityCached(ctx, address)
	suite.False(found)
}

// TestCachedCredentialOperations tests cached credential operations
func (suite *CacheTestSuite) TestCachedCredentialOperations() {
	ctx := suite.ctx
	ck := suite.cachedKeeper
	holder := suite.addrs[0].String()
	
	// Create identity first
	identity := types.Identity{
		Address:        holder,
		DID:            fmt.Sprintf("did:desh:test_%s", holder),
		Status:         types.IdentityStatus_ACTIVE,
		CreatedAt:      ctx.BlockTime(),
		UpdatedAt:      ctx.BlockTime(),
		LastActivityAt: ctx.BlockTime(),
	}
	ck.SetIdentityCached(ctx, identity)
	
	// Issue a credential
	issuer := suite.addrs[1].String()
	credentialTypes := []string{"VerifiableCredential", "TestCredential"}
	claims := map[string]interface{}{
		"name": "Test User",
		"age":  25,
	}
	
	credential, err := ck.IssueCredential(ctx, issuer, holder, credentialTypes, claims, 365)
	suite.NoError(err)
	
	// Cache the credential
	ck.SetCredentialCached(ctx, *credential)
	
	// Get credential - should hit cache
	retrieved, found := ck.GetCredentialCached(ctx, credential.ID)
	suite.True(found)
	suite.Equal(credential.ID, retrieved.ID)
	suite.Equal(credential.Holder, retrieved.Holder)
	
	// Get credentials by holder - should hit cache
	credentials := ck.GetCredentialsByHolderCached(ctx, holder)
	suite.Len(credentials, 1)
	suite.Equal(credential.ID, credentials[0].ID)
}

// TestCachedDIDOperations tests cached DID operations
func (suite *CacheTestSuite) TestCachedDIDOperations() {
	ctx := suite.ctx
	ck := suite.cachedKeeper
	did := "did:desh:test123"
	
	// Create DID document
	didDoc := types.DIDDocument{
		Context:    []string{"https://www.w3.org/ns/did/v1"},
		ID:         did,
		Controller: did,
		Created:    ctx.BlockTime(),
		Updated:    ctx.BlockTime(),
	}
	
	// Set DID document (should cache it)
	ck.SetDIDDocumentCached(ctx, didDoc)
	
	// Get DID document - should hit cache
	retrieved, found := ck.GetDIDDocumentCached(ctx, did)
	suite.True(found)
	suite.Equal(didDoc.ID, retrieved.ID)
	suite.Equal(didDoc.Controller, retrieved.Controller)
	
	// Verify cache hit
	stats := ck.GetCacheStats()
	suite.True(stats.CacheHits > 0)
}

// TestCacheInvalidation tests cache invalidation
func (suite *CacheTestSuite) TestCacheInvalidation() {
	ctx := suite.ctx
	ck := suite.cachedKeeper
	address := suite.addrs[0].String()
	
	// Create and cache identity
	identity := types.Identity{
		Address:        address,
		DID:            fmt.Sprintf("did:desh:test_%s", address),
		Status:         types.IdentityStatus_ACTIVE,
		CreatedAt:      ctx.BlockTime(),
		UpdatedAt:      ctx.BlockTime(),
		LastActivityAt: ctx.BlockTime(),
	}
	ck.SetIdentityCached(ctx, identity)
	
	// Create and cache credential
	credential, err := ck.IssueCredential(ctx, suite.addrs[1].String(), address, []string{"TestCredential"}, map[string]interface{}{"test": "data"}, 365)
	suite.NoError(err)
	ck.SetCredentialCached(ctx, *credential)
	
	// Verify both are cached
	_, found1 := ck.GetIdentityCached(ctx, address)
	_, found2 := ck.GetCredentialCached(ctx, credential.ID)
	suite.True(found1)
	suite.True(found2)
	
	// Invalidate user cache
	ck.InvalidateUserCache(ctx, address)
	
	// Identity should be invalidated
	_, found1 = ck.GetIdentityCached(ctx, address)
	suite.False(found1)
	
	// Note: In a full implementation, related credentials would also be invalidated
}

// TestCachePerformance tests cache performance improvements
func (suite *CacheTestSuite) TestCachePerformance() {
	ctx := suite.ctx
	ck := suite.cachedKeeper
	address := suite.addrs[0].String()
	
	// Create identity
	identity := types.Identity{
		Address:        address,
		DID:            fmt.Sprintf("did:desh:perf_%s", address),
		Status:         types.IdentityStatus_ACTIVE,
		CreatedAt:      ctx.BlockTime(),
		UpdatedAt:      ctx.BlockTime(),
		LastActivityAt: ctx.BlockTime(),
	}
	
	// Measure cold access time (cache miss)
	start := time.Now()
	ck.SetIdentityCached(ctx, identity)
	coldTime := time.Since(start)
	
	// Measure warm access time (cache hit)
	start = time.Now()
	_, found := ck.GetIdentityCached(ctx, address)
	warmTime := time.Since(start)
	
	suite.True(found)
	suite.True(warmTime < coldTime, "Cache hit should be faster than cache miss")
	
	// Multiple access should maintain performance
	for i := 0; i < 10; i++ {
		start = time.Now()
		_, found = ck.GetIdentityCached(ctx, address)
		accessTime := time.Since(start)
		
		suite.True(found)
		suite.True(accessTime < coldTime, "Repeated access should remain fast")
	}
}

// TestCacheMetrics tests detailed cache metrics
func (suite *CacheTestSuite) TestCacheMetrics() {
	ctx := suite.ctx
	ck := suite.cachedKeeper
	
	// Perform various operations
	for i := 0; i < 5; i++ {
		address := suite.addrs[i].String()
		identity := types.Identity{
			Address:        address,
			DID:            fmt.Sprintf("did:desh:metrics_%d", i),
			Status:         types.IdentityStatus_ACTIVE,
			CreatedAt:      ctx.BlockTime(),
			UpdatedAt:      ctx.BlockTime(),
			LastActivityAt: ctx.BlockTime(),
		}
		ck.SetIdentityCached(ctx, identity)
	}
	
	// Access some entries (hits)
	for i := 0; i < 3; i++ {
		address := suite.addrs[i].String()
		_, found := ck.GetIdentityCached(ctx, address)
		suite.True(found)
	}
	
	// Try to access non-existent entries (misses)
	for i := 5; i < 8; i++ {
		address := fmt.Sprintf("non-existent-%d", i)
		_, found := ck.GetIdentityCached(ctx, address)
		suite.False(found)
	}
	
	// Get metrics
	metrics := ck.GetCacheMetrics()
	
	suite.True(metrics.IdentityHits >= 3)
	suite.True(metrics.IdentityMisses >= 3)
	suite.True(metrics.EntryCount >= 5)
	suite.True(metrics.MemoryUsage > 0)
	suite.True(metrics.AverageGetTime >= 0)
	suite.True(metrics.AverageSetTime >= 0)
}

// TestCacheConcurrency tests cache thread safety
func (suite *CacheTestSuite) TestCacheConcurrency() {
	cache := suite.cache
	
	// Number of concurrent operations
	numGoroutines := 10
	numOperations := 100
	
	// Channel to signal completion
	done := make(chan bool, numGoroutines)
	
	// Start concurrent goroutines
	for g := 0; g < numGoroutines; g++ {
		go func(goroutineID int) {
			defer func() { done <- true }()
			
			for i := 0; i < numOperations; i++ {
				key := types.NewCacheKey("concurrent", fmt.Sprintf("g%d_key%d", goroutineID, i))
				data := fmt.Sprintf("g%d_data%d", goroutineID, i)
				
				// Set
				err := cache.Set(key, data, 1*time.Minute)
				suite.NoError(err)
				
				// Get
				entry, found := cache.Get(key)
				suite.True(found)
				suite.Equal(data, entry.Data)
				
				// Delete
				err = cache.Delete(key)
				suite.NoError(err)
			}
		}(g)
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		select {
		case <-done:
			// Goroutine completed
		case <-time.After(30 * time.Second):
			suite.Fail("Timeout waiting for concurrent operations")
		}
	}
	
	// Verify cache integrity
	stats := cache.Stats()
	suite.True(stats.TotalRequests > 0)
}

// TestCacheEviction tests cache eviction policies
func (suite *CacheTestSuite) TestCacheEviction() {
	// Create a cache with very small limits for testing eviction
	smallCacheConfig := types.CacheConfig{
		MaxSize:       1024,  // 1KB
		MaxEntries:    5,     // Only 5 entries
		DefaultTTL:    1 * time.Minute,
		EnableMetrics: true,
	}
	
	smallCache := keeper.NewIdentityCache(smallCacheConfig)
	err := smallCache.Start()
	suite.NoError(err)
	defer smallCache.Stop()
	
	// Fill cache beyond capacity
	for i := 0; i < 10; i++ {
		key := types.NewCacheKey("test", fmt.Sprintf("eviction_key%d", i))
		data := fmt.Sprintf("eviction_data_%d_with_some_extra_content", i)
		
		err := smallCache.Set(key, data, 1*time.Minute)
		suite.NoError(err)
	}
	
	// Check that eviction occurred
	stats := smallCache.Stats()
	suite.True(stats.TotalEntries <= 5, "Cache should not exceed max entries")
	suite.True(stats.EvictionCount > 0, "Some evictions should have occurred")
	
	// Verify some entries were evicted
	evictedCount := 0
	for i := 0; i < 10; i++ {
		key := types.NewCacheKey("test", fmt.Sprintf("eviction_key%d", i))
		if !smallCache.Exists(key) {
			evictedCount++
		}
	}
	
	suite.True(evictedCount > 0, "Some entries should have been evicted")
}

// TestCacheCleanup tests automatic cleanup of expired entries
func (suite *CacheTestSuite) TestCacheCleanup() {
	cache := suite.cache
	
	// Add entries with very short TTL
	for i := 0; i < 5; i++ {
		key := types.NewCacheKey("cleanup", fmt.Sprintf("key%d", i))
		data := fmt.Sprintf("data%d", i)
		ttl := 100 * time.Millisecond
		
		err := cache.Set(key, data, ttl)
		suite.NoError(err)
	}
	
	// All entries should exist initially
	suite.Equal(int64(5), cache.EntryCount())
	
	// Wait for expiration
	time.Sleep(150 * time.Millisecond)
	
	// Trigger cleanup by accessing the cache
	// (In practice, cleanup runs automatically)
	for i := 0; i < 5; i++ {
		key := types.NewCacheKey("cleanup", fmt.Sprintf("key%d", i))
		cache.Get(key) // This should trigger cleanup of expired entries
	}
	
	// Give some time for cleanup to complete
	time.Sleep(100 * time.Millisecond)
	
	// Verify entries were cleaned up
	for i := 0; i < 5; i++ {
		key := types.NewCacheKey("cleanup", fmt.Sprintf("key%d", i))
		suite.False(cache.Exists(key))
	}
}