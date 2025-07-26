package keeper

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"sync"
	"time"

	"github.com/namo/x/identity/types"
)

// IdentityCache implements a high-performance in-memory cache for identity data
type IdentityCache struct {
	mu              sync.RWMutex
	entries         map[string]*types.CacheEntry
	config          types.CacheConfig
	stats           types.CacheStats
	metrics         types.CacheMetrics
	observers       []types.CacheObserver
	indexer         *CacheIndexer
	cleanupTicker   *time.Ticker
	running         bool
	evictionPolicy  types.EvictionPolicy
	
	// Performance tracking
	operationTimes  map[string][]float64 // operation -> duration list
	lastCleanup     time.Time
}

// CacheIndexer provides fast lookups by different criteria
type CacheIndexer struct {
	mu          sync.RWMutex
	didIndex    map[string][]string // DID -> cache keys
	typeIndex   map[string][]string // type -> cache keys
	tagIndex    map[string][]string // tag -> cache keys
}

// NewIdentityCache creates a new identity cache instance
func NewIdentityCache(config types.CacheConfig) *IdentityCache {
	cache := &IdentityCache{
		entries:        make(map[string]*types.CacheEntry),
		config:         config,
		stats:          types.CacheStats{LastResetTime: time.Now()},
		metrics:        types.CacheMetrics{LastUpdated: time.Now()},
		observers:      make([]types.CacheObserver, 0),
		indexer:        NewCacheIndexer(),
		evictionPolicy: types.EvictionPolicy_HYBRID,
		operationTimes: make(map[string][]float64),
		lastCleanup:    time.Now(),
	}
	
	return cache
}

// NewCacheIndexer creates a new cache indexer
func NewCacheIndexer() *CacheIndexer {
	return &CacheIndexer{
		didIndex:  make(map[string][]string),
		typeIndex: make(map[string][]string),
		tagIndex:  make(map[string][]string),
	}
}

// Start initializes and starts the cache
func (c *IdentityCache) Start() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	if c.running {
		return fmt.Errorf("cache is already running")
	}
	
	// Start cleanup routine
	c.cleanupTicker = time.NewTicker(c.config.CleanupInterval)
	go c.cleanupRoutine()
	
	// Warmup cache if configured
	if c.config.WarmupOnStart {
		go c.warmupCache()
	}
	
	c.running = true
	c.emitEvent(types.CacheEvent{
		Type:      "start",
		Timestamp: time.Now(),
	})
	
	return nil
}

// Stop gracefully stops the cache
func (c *IdentityCache) Stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	if !c.running {
		return fmt.Errorf("cache is not running")
	}
	
	if c.cleanupTicker != nil {
		c.cleanupTicker.Stop()
	}
	
	c.running = false
	c.emitEvent(types.CacheEvent{
		Type:      "stop",
		Timestamp: time.Now(),
	})
	
	return nil
}

// IsRunning returns whether the cache is running
func (c *IdentityCache) IsRunning() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.running
}

// Get retrieves an entry from the cache
func (c *IdentityCache) Get(key types.CacheKey) (*types.CacheEntry, bool) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		c.trackOperation("get", duration)
	}()
	
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	keyStr := c.keyToString(key)
	entry, exists := c.entries[keyStr]
	
	if !exists {
		c.recordMiss(key.Type)
		c.emitEvent(types.CacheEvent{
			Type:      types.CacheEventMiss,
			Key:       key,
			Timestamp: time.Now(),
			Duration:  time.Since(start),
		})
		return nil, false
	}
	
	// Check if entry is expired
	if entry.IsExpired() {
		c.recordMiss(key.Type)
		// Remove expired entry asynchronously if configured
		if c.config.AsyncEviction {
			go c.deleteExpired(keyStr)
		} else {
			delete(c.entries, keyStr)
			c.indexer.RemoveFromIndex(key)
		}
		c.emitEvent(types.CacheEvent{
			Type:      types.CacheEventExpire,
			Key:       key,
			Timestamp: time.Now(),
			Duration:  time.Since(start),
		})
		return nil, false
	}
	
	// Update access statistics
	entry.Touch()
	c.recordHit(key.Type)
	
	c.emitEvent(types.CacheEvent{
		Type:      types.CacheEventHit,
		Key:       key,
		Timestamp: time.Now(),
		Duration:  time.Since(start),
		Size:      entry.Size,
	})
	
	return entry, true
}

// Set stores an entry in the cache
func (c *IdentityCache) Set(key types.CacheKey, data interface{}, ttl time.Duration) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		c.trackOperation("set", duration)
	}()
	
	c.mu.Lock()
	defer c.mu.Unlock()
	
	// Create cache entry
	entry := types.NewCacheEntry(key, data, ttl)
	entry.Size = c.calculateSize(data)
	
	// Add type-specific tags
	c.addTypeSpecificTags(entry, key.Type)
	
	keyStr := c.keyToString(key)
	
	// Check if we need to evict entries
	if c.needsEviction(entry.Size) {
		if err := c.evictEntries(entry.Size); err != nil {
			return fmt.Errorf("failed to evict entries: %w", err)
		}
	}
	
	// Store the entry
	c.entries[keyStr] = entry
	
	// Update indices
	c.indexer.IndexByType(key.Type, entry)
	if key.Type == types.CacheTypeIdentity || key.Type == types.CacheTypeDIDDocument {
		c.indexer.IndexByDID(key.Key, entry)
	}
	for _, tag := range entry.Tags {
		c.indexer.IndexByTag(tag, entry)
	}
	
	// Update statistics
	c.stats.TotalEntries++
	c.stats.TotalSize += entry.Size
	
	c.emitEvent(types.CacheEvent{
		Type:      types.CacheEventSet,
		Key:       key,
		Timestamp: time.Now(),
		Duration:  time.Since(start),
		Size:      entry.Size,
	})
	
	return nil
}

// Delete removes an entry from the cache
func (c *IdentityCache) Delete(key types.CacheKey) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		c.trackOperation("delete", duration)
	}()
	
	c.mu.Lock()
	defer c.mu.Unlock()
	
	keyStr := c.keyToString(key)
	entry, exists := c.entries[keyStr]
	
	if !exists {
		return fmt.Errorf("key not found: %s", keyStr)
	}
	
	// Remove from cache and indices
	delete(c.entries, keyStr)
	c.indexer.RemoveFromIndex(key)
	
	// Update statistics
	c.stats.TotalEntries--
	c.stats.TotalSize -= entry.Size
	
	c.emitEvent(types.CacheEvent{
		Type:      types.CacheEventDelete,
		Key:       key,
		Timestamp: time.Now(),
		Duration:  time.Since(start),
		Size:      entry.Size,
	})
	
	return nil
}

// Exists checks if a key exists in the cache
func (c *IdentityCache) Exists(key types.CacheKey) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	keyStr := c.keyToString(key)
	entry, exists := c.entries[keyStr]
	
	if !exists {
		return false
	}
	
	return !entry.IsExpired()
}

// GetMultiple retrieves multiple entries from the cache
func (c *IdentityCache) GetMultiple(keys []types.CacheKey) map[string]*types.CacheEntry {
	result := make(map[string]*types.CacheEntry)
	
	for _, key := range keys {
		if entry, found := c.Get(key); found {
			result[c.keyToString(key)] = entry
		}
	}
	
	return result
}

// SetMultiple stores multiple entries in the cache
func (c *IdentityCache) SetMultiple(entries map[types.CacheKey]interface{}, ttl time.Duration) error {
	for key, data := range entries {
		if err := c.Set(key, data, ttl); err != nil {
			return err
		}
	}
	return nil
}

// DeleteMultiple removes multiple entries from the cache
func (c *IdentityCache) DeleteMultiple(keys []types.CacheKey) error {
	for _, key := range keys {
		if err := c.Delete(key); err != nil {
			// Continue with other deletions even if one fails
			continue
		}
	}
	return nil
}

// Clear removes all entries from the cache
func (c *IdentityCache) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	entryCount := len(c.entries)
	c.entries = make(map[string]*types.CacheEntry)
	c.indexer.ClearIndex()
	
	c.stats.TotalEntries = 0
	c.stats.TotalSize = 0
	
	c.emitEvent(types.CacheEvent{
		Type:      types.CacheEventClear,
		Timestamp: time.Now(),
		Metadata:  map[string]interface{}{"cleared_entries": entryCount},
	})
	
	return nil
}

// Size returns the total size of cached data in bytes
func (c *IdentityCache) Size() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.stats.TotalSize
}

// EntryCount returns the number of entries in the cache
func (c *IdentityCache) EntryCount() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.stats.TotalEntries
}

// Stats returns cache statistics
func (c *IdentityCache) Stats() types.CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	stats := c.stats
	stats.HitRatio = stats.GetHitRatio()
	return stats
}

// GetByTag retrieves all entries with a specific tag
func (c *IdentityCache) GetByTag(tag string) []*types.CacheEntry {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	var result []*types.CacheEntry
	keys := c.indexer.GetByTag(tag)
	
	for _, keyStr := range keys {
		if entry, exists := c.entries[keyStr]; exists && !entry.IsExpired() {
			result = append(result, entry)
		}
	}
	
	return result
}

// InvalidateByTag removes all entries with a specific tag
func (c *IdentityCache) InvalidateByTag(tag string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	keys := c.indexer.GetByTag(tag)
	invalidatedCount := 0
	
	for _, keyStr := range keys {
		if entry, exists := c.entries[keyStr]; exists {
			delete(c.entries, keyStr)
			c.stats.TotalEntries--
			c.stats.TotalSize -= entry.Size
			invalidatedCount++
		}
	}
	
	// Remove from tag index
	c.indexer.tagIndex[tag] = nil
	
	c.emitEvent(types.CacheEvent{
		Type:      types.CacheEventInvalidate,
		Timestamp: time.Now(),
		Metadata:  map[string]interface{}{
			"tag": tag,
			"invalidated_count": invalidatedCount,
		},
	})
	
	return nil
}

// GetByPattern retrieves entries matching a key pattern
func (c *IdentityCache) GetByPattern(pattern string) []*types.CacheEntry {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil
	}
	
	var result []*types.CacheEntry
	for keyStr, entry := range c.entries {
		if regex.MatchString(keyStr) && !entry.IsExpired() {
			result = append(result, entry)
		}
	}
	
	return result
}

// Refresh reloads an entry (placeholder for external data source integration)
func (c *IdentityCache) Refresh(key types.CacheKey) error {
	// This would typically reload from the underlying data source
	// For now, we just remove the entry to force a reload
	return c.Delete(key)
}

// GetMetrics returns detailed cache metrics
func (c *IdentityCache) GetMetrics() types.CacheMetrics {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	metrics := c.metrics
	metrics.EntryCount = c.stats.TotalEntries
	metrics.MemoryUsage = c.stats.TotalSize
	metrics.LastUpdated = time.Now()
	
	// Calculate average operation times
	if times, exists := c.operationTimes["get"]; exists && len(times) > 0 {
		metrics.AverageGetTime = c.calculateAverage(times)
	}
	if times, exists := c.operationTimes["set"]; exists && len(times) > 0 {
		metrics.AverageSetTime = c.calculateAverage(times)
	}
	if times, exists := c.operationTimes["delete"]; exists && len(times) > 0 {
		metrics.AverageDeleteTime = c.calculateAverage(times)
	}
	
	return metrics
}

// AddObserver adds a cache event observer
func (c *IdentityCache) AddObserver(observer types.CacheObserver) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.observers = append(c.observers, observer)
}

// RemoveObserver removes a cache event observer
func (c *IdentityCache) RemoveObserver(observer types.CacheObserver) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	for i, obs := range c.observers {
		if obs == observer {
			c.observers = append(c.observers[:i], c.observers[i+1:]...)
			break
		}
	}
}

// Helper methods

func (c *IdentityCache) keyToString(key types.CacheKey) string {
	return fmt.Sprintf("%s:%s:%d", key.Type, key.Key, key.Version)
}

func (c *IdentityCache) calculateSize(data interface{}) int64 {
	// Simple size calculation using JSON marshaling
	// In production, you might want a more sophisticated approach
	if jsonData, err := json.Marshal(data); err == nil {
		return int64(len(jsonData))
	}
	return 0
}

func (c *IdentityCache) addTypeSpecificTags(entry *types.CacheEntry, entryType string) {
	entry.AddTag(entryType)
	
	// Add additional tags based on data type
	switch entryType {
	case types.CacheTypeIdentity:
		entry.AddTag("user_data")
	case types.CacheTypeCredential:
		entry.AddTag("verification")
	case types.CacheTypeDIDDocument:
		entry.AddTag("resolution")
	case types.CacheTypeConsent:
		entry.AddTag("privacy")
	case types.CacheTypeZKProof:
		entry.AddTag("privacy")
		entry.AddTag("proof")
	}
}

func (c *IdentityCache) needsEviction(newEntrySize int64) bool {
	return (c.stats.TotalSize+newEntrySize > c.config.MaxSize) ||
		   (c.stats.TotalEntries >= c.config.MaxEntries)
}

func (c *IdentityCache) evictEntries(sizeNeeded int64) error {
	var evicted int64
	var evictedCount int64
	
	// Get entries sorted by eviction criteria
	candidates := c.getEvictionCandidates()
	
	for _, keyStr := range candidates {
		if entry, exists := c.entries[keyStr]; exists {
			delete(c.entries, keyStr)
			
			evicted += entry.Size
			evictedCount++
			c.stats.EvictionCount++
			
			c.emitEvent(types.CacheEvent{
				Type:      types.CacheEventEvict,
				Timestamp: time.Now(),
				Size:      entry.Size,
			})
			
			// Stop if we've freed enough space
			if evicted >= sizeNeeded && c.stats.TotalEntries < c.config.MaxEntries {
				break
			}
		}
	}
	
	c.stats.TotalEntries -= evictedCount
	c.stats.TotalSize -= evicted
	
	return nil
}

func (c *IdentityCache) getEvictionCandidates() []string {
	type entryInfo struct {
		key    string
		score  float64 // lower score = higher eviction priority
	}
	
	var candidates []entryInfo
	now := time.Now()
	
	for keyStr, entry := range c.entries {
		score := float64(0)
		
		switch c.evictionPolicy {
		case types.EvictionPolicy_LRU:
			score = float64(now.Sub(entry.LastAccess).Seconds())
		case types.EvictionPolicy_LFU:
			score = 1.0 / float64(entry.AccessCount+1)
		case types.EvictionPolicy_TTL:
			score = float64(entry.ExpiresAt.Sub(now).Seconds())
		case types.EvictionPolicy_SIZE:
			score = float64(entry.Size)
		case types.EvictionPolicy_HYBRID:
			// Weighted combination of LRU, LFU, and size
			lruScore := float64(now.Sub(entry.LastAccess).Seconds()) * 0.4
			lfuScore := (1.0 / float64(entry.AccessCount+1)) * 0.3
			sizeScore := float64(entry.Size) * 0.3
			score = lruScore + lfuScore + sizeScore
		}
		
		candidates = append(candidates, entryInfo{key: keyStr, score: score})
	}
	
	// Sort by score (ascending for most eviction policies)
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].score < candidates[j].score
	})
	
	result := make([]string, len(candidates))
	for i, candidate := range candidates {
		result[i] = candidate.key
	}
	
	return result
}

func (c *IdentityCache) recordHit(entryType string) {
	c.stats.TotalRequests++
	c.stats.CacheHits++
	
	// Update type-specific metrics
	switch entryType {
	case types.CacheTypeIdentity:
		c.metrics.IdentityHits++
	case types.CacheTypeCredential:
		c.metrics.CredentialHits++
	case types.CacheTypeDIDDocument:
		c.metrics.DIDHits++
	}
}

func (c *IdentityCache) recordMiss(entryType string) {
	c.stats.TotalRequests++
	c.stats.CacheMisses++
	
	// Update type-specific metrics
	switch entryType {
	case types.CacheTypeIdentity:
		c.metrics.IdentityMisses++
	case types.CacheTypeCredential:
		c.metrics.CredentialMisses++
	case types.CacheTypeDIDDocument:
		c.metrics.DIDMisses++
	}
}

func (c *IdentityCache) trackOperation(operation string, duration time.Duration) {
	if !c.config.EnableMetrics {
		return
	}
	
	durationMs := float64(duration.Nanoseconds()) / 1e6 // Convert to milliseconds
	
	if times, exists := c.operationTimes[operation]; exists {
		// Keep only last 1000 measurements
		if len(times) >= 1000 {
			times = times[1:]
		}
		c.operationTimes[operation] = append(times, durationMs)
	} else {
		c.operationTimes[operation] = []float64{durationMs}
	}
}

func (c *IdentityCache) calculateAverage(times []float64) float64 {
	if len(times) == 0 {
		return 0
	}
	
	var sum float64
	for _, time := range times {
		sum += time
	}
	return sum / float64(len(times))
}

func (c *IdentityCache) emitEvent(event types.CacheEvent) {
	if !c.config.EnableMetrics {
		return
	}
	
	for _, observer := range c.observers {
		go observer.OnCacheEvent(event) // Async to prevent blocking
	}
}

func (c *IdentityCache) deleteExpired(keyStr string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	if entry, exists := c.entries[keyStr]; exists && entry.IsExpired() {
		delete(c.entries, keyStr)
		c.stats.TotalEntries--
		c.stats.TotalSize -= entry.Size
	}
}

func (c *IdentityCache) cleanupRoutine() {
	for range c.cleanupTicker.C {
		c.cleanup()
	}
}

func (c *IdentityCache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	now := time.Now()
	var expiredKeys []string
	var expiredSize int64
	
	for keyStr, entry := range c.entries {
		if entry.IsExpired() {
			expiredKeys = append(expiredKeys, keyStr)
			expiredSize += entry.Size
		}
	}
	
	// Remove expired entries
	for _, keyStr := range expiredKeys {
		delete(c.entries, keyStr)
	}
	
	c.stats.TotalEntries -= int64(len(expiredKeys))
	c.stats.TotalSize -= expiredSize
	c.lastCleanup = now
	
	if len(expiredKeys) > 0 {
		c.emitEvent(types.CacheEvent{
			Type:      "cleanup",
			Timestamp: now,
			Metadata: map[string]interface{}{
				"expired_count": len(expiredKeys),
				"freed_size":    expiredSize,
			},
		})
	}
}

func (c *IdentityCache) warmupCache() {
	// This would typically load frequently accessed data
	// Implementation depends on integration with the keeper
	c.emitEvent(types.CacheEvent{
		Type:      types.CacheEventWarmup,
		Timestamp: time.Now(),
	})
}

// CacheIndexer implementation

func (ci *CacheIndexer) IndexByDID(did string, entry *types.CacheEntry) error {
	ci.mu.Lock()
	defer ci.mu.Unlock()
	
	keyStr := fmt.Sprintf("%s:%s:%d", entry.Key.Type, entry.Key.Key, entry.Key.Version)
	
	if keys, exists := ci.didIndex[did]; exists {
		ci.didIndex[did] = append(keys, keyStr)
	} else {
		ci.didIndex[did] = []string{keyStr}
	}
	
	return nil
}

func (ci *CacheIndexer) IndexByType(entryType string, entry *types.CacheEntry) error {
	ci.mu.Lock()
	defer ci.mu.Unlock()
	
	keyStr := fmt.Sprintf("%s:%s:%d", entry.Key.Type, entry.Key.Key, entry.Key.Version)
	
	if keys, exists := ci.typeIndex[entryType]; exists {
		ci.typeIndex[entryType] = append(keys, keyStr)
	} else {
		ci.typeIndex[entryType] = []string{keyStr}
	}
	
	return nil
}

func (ci *CacheIndexer) IndexByTag(tag string, entry *types.CacheEntry) error {
	ci.mu.Lock()
	defer ci.mu.Unlock()
	
	keyStr := fmt.Sprintf("%s:%s:%d", entry.Key.Type, entry.Key.Key, entry.Key.Version)
	
	if keys, exists := ci.tagIndex[tag]; exists {
		ci.tagIndex[tag] = append(keys, keyStr)
	} else {
		ci.tagIndex[tag] = []string{keyStr}
	}
	
	return nil
}

func (ci *CacheIndexer) GetByDID(did string) []string {
	ci.mu.RLock()
	defer ci.mu.RUnlock()
	
	if keys, exists := ci.didIndex[did]; exists {
		return keys
	}
	return nil
}

func (ci *CacheIndexer) GetByType(entryType string) []string {
	ci.mu.RLock()
	defer ci.mu.RUnlock()
	
	if keys, exists := ci.typeIndex[entryType]; exists {
		return keys
	}
	return nil
}

func (ci *CacheIndexer) GetByTag(tag string) []string {
	ci.mu.RLock()
	defer ci.mu.RUnlock()
	
	if keys, exists := ci.tagIndex[tag]; exists {
		return keys
	}
	return nil
}

func (ci *CacheIndexer) RemoveFromIndex(key types.CacheKey) error {
	ci.mu.Lock()
	defer ci.mu.Unlock()
	
	keyStr := fmt.Sprintf("%s:%s:%d", key.Type, key.Key, key.Version)
	
	// Remove from all indices
	for did, keys := range ci.didIndex {
		ci.didIndex[did] = removeFromSlice(keys, keyStr)
	}
	
	for entryType, keys := range ci.typeIndex {
		ci.typeIndex[entryType] = removeFromSlice(keys, keyStr)
	}
	
	for tag, keys := range ci.tagIndex {
		ci.tagIndex[tag] = removeFromSlice(keys, keyStr)
	}
	
	return nil
}

func (ci *CacheIndexer) ClearIndex() error {
	ci.mu.Lock()
	defer ci.mu.Unlock()
	
	ci.didIndex = make(map[string][]string)
	ci.typeIndex = make(map[string][]string)
	ci.tagIndex = make(map[string][]string)
	
	return nil
}

// Helper function to remove a string from a slice
func removeFromSlice(slice []string, item string) []string {
	for i, v := range slice {
		if v == item {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}