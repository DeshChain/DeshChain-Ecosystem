package types

import (
	"sync"
	"time"
)

// Identity Caching System Types and Interfaces

// CacheKey represents a cache key with metadata
type CacheKey struct {
	Type     string // identity, credential, did_document, etc.
	Key      string // actual key (address, DID, credential ID, etc.)
	Version  int64  // version for cache invalidation
	Hash     string // hash of the data for integrity
}

// CacheEntry represents a cached item with metadata
type CacheEntry struct {
	Key         CacheKey               `json:"key"`
	Data        interface{}            `json:"data"`
	ExpiresAt   time.Time              `json:"expires_at"`
	AccessCount int64                  `json:"access_count"`
	LastAccess  time.Time              `json:"last_access"`
	Size        int64                  `json:"size"` // size in bytes
	Tags        []string               `json:"tags"` // for tag-based invalidation
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// CacheStats represents cache performance statistics
type CacheStats struct {
	TotalRequests   int64   `json:"total_requests"`
	CacheHits       int64   `json:"cache_hits"`
	CacheMisses     int64   `json:"cache_misses"`
	HitRatio        float64 `json:"hit_ratio"`
	TotalEntries    int64   `json:"total_entries"`
	TotalSize       int64   `json:"total_size"`
	EvictionCount   int64   `json:"eviction_count"`
	AverageLoadTime float64 `json:"average_load_time"` // in milliseconds
	LastResetTime   time.Time `json:"last_reset_time"`
}

// CacheConfig represents cache configuration
type CacheConfig struct {
	MaxSize        int64         `json:"max_size"`         // maximum cache size in bytes
	MaxEntries     int64         `json:"max_entries"`      // maximum number of entries
	DefaultTTL     time.Duration `json:"default_ttl"`      // default time to live
	CleanupInterval time.Duration `json:"cleanup_interval"` // cleanup interval
	EnableMetrics  bool          `json:"enable_metrics"`   // enable performance metrics
	EnableTags     bool          `json:"enable_tags"`      // enable tag-based invalidation
	
	// Type-specific TTLs
	IdentityTTL    time.Duration `json:"identity_ttl"`
	CredentialTTL  time.Duration `json:"credential_ttl"`
	DIDDocumentTTL time.Duration `json:"did_document_ttl"`
	ConsentTTL     time.Duration `json:"consent_ttl"`
	ZKProofTTL     time.Duration `json:"zk_proof_ttl"`
	
	// Performance tuning
	PreloadIdentities bool `json:"preload_identities"` // preload active identities
	WarmupOnStart     bool `json:"warmup_on_start"`    // warm up cache on startup
	AsyncEviction     bool `json:"async_eviction"`     // use async eviction
}

// EvictionPolicy represents cache eviction policies
type EvictionPolicy int32

const (
	EvictionPolicy_LRU     EvictionPolicy = 0 // Least Recently Used
	EvictionPolicy_LFU     EvictionPolicy = 1 // Least Frequently Used
	EvictionPolicy_TTL     EvictionPolicy = 2 // Time To Live based
	EvictionPolicy_SIZE    EvictionPolicy = 3 // Size based
	EvictionPolicy_HYBRID  EvictionPolicy = 4 // Hybrid (LRU + TTL + Size)
)

// CacheInterface defines the contract for identity cache implementations
type CacheInterface interface {
	// Basic operations
	Get(key CacheKey) (*CacheEntry, bool)
	Set(key CacheKey, data interface{}, ttl time.Duration) error
	Delete(key CacheKey) error
	Exists(key CacheKey) bool
	
	// Batch operations
	GetMultiple(keys []CacheKey) map[string]*CacheEntry
	SetMultiple(entries map[CacheKey]interface{}, ttl time.Duration) error
	DeleteMultiple(keys []CacheKey) error
	
	// Cache management
	Clear() error
	Size() int64
	EntryCount() int64
	Stats() CacheStats
	
	// Advanced operations
	GetByTag(tag string) []*CacheEntry
	InvalidateByTag(tag string) error
	GetByPattern(pattern string) []*CacheEntry
	Refresh(key CacheKey) error
	
	// Lifecycle
	Start() error
	Stop() error
	IsRunning() bool
}

// IdentityCacheEntry represents a cached identity with enriched metadata
type IdentityCacheEntry struct {
	Identity        *Identity              `json:"identity"`
	DIDDocument     *DIDDocument           `json:"did_document,omitempty"`
	Credentials     []*Credential          `json:"credentials,omitempty"`
	ActiveConsents  []*ConsentRecord       `json:"active_consents,omitempty"`
	LastActivity    time.Time              `json:"last_activity"`
	AccessFrequency int64                  `json:"access_frequency"`
	RelatedEntities []string               `json:"related_entities"` // DIDs of related entities
	ComputedFields  map[string]interface{} `json:"computed_fields"`   // pre-computed expensive fields
}

// CredentialCacheEntry represents a cached credential with verification status
type CredentialCacheEntry struct {
	Credential        *Credential    `json:"credential"`
	VerificationStatus string        `json:"verification_status"`
	LastVerified      time.Time      `json:"last_verified"`
	IssuerReputation  float64        `json:"issuer_reputation"`
	UsageCount        int64          `json:"usage_count"`
	RelatedProofs     []string       `json:"related_proofs"` // related ZK proof IDs
}

// DIDCacheEntry represents a cached DID document with resolution metadata
type DIDCacheEntry struct {
	DIDDocument      *DIDDocument   `json:"did_document"`
	ResolutionTime   time.Time      `json:"resolution_time"`
	MethodMetadata   map[string]interface{} `json:"method_metadata"`
	ControllerChain  []string       `json:"controller_chain"` // chain of controllers
	ServiceEndpoints []*Service     `json:"service_endpoints"`
}

// CacheWarmupStrategy represents different warmup strategies
type CacheWarmupStrategy struct {
	PreloadActiveIdentities  bool     `json:"preload_active_identities"`
	PreloadRecentCredentials bool     `json:"preload_recent_credentials"`
	PreloadPopularDIDs       bool     `json:"preload_popular_dids"`
	PreloadByTags           []string  `json:"preload_by_tags"`
	MaxPreloadEntries       int64     `json:"max_preload_entries"`
	PreloadBatchSize        int32     `json:"preload_batch_size"`
}

// CacheMetrics represents detailed cache performance metrics
type CacheMetrics struct {
	// Hit/Miss metrics
	IdentityHits    int64 `json:"identity_hits"`
	IdentityMisses  int64 `json:"identity_misses"`
	CredentialHits  int64 `json:"credential_hits"`
	CredentialMisses int64 `json:"credential_misses"`
	DIDHits         int64 `json:"did_hits"`
	DIDMisses       int64 `json:"did_misses"`
	
	// Performance metrics
	AverageGetTime    float64 `json:"average_get_time"`    // milliseconds
	AverageSetTime    float64 `json:"average_set_time"`    // milliseconds
	AverageDeleteTime float64 `json:"average_delete_time"` // milliseconds
	
	// Memory metrics
	MemoryUsage      int64 `json:"memory_usage"`       // bytes
	MaxMemoryUsage   int64 `json:"max_memory_usage"`   // bytes
	EntryCount       int64 `json:"entry_count"`
	MaxEntryCount    int64 `json:"max_entry_count"`
	
	// Eviction metrics
	LRUEvictions     int64 `json:"lru_evictions"`
	TTLEvictions     int64 `json:"ttl_evictions"`
	SizeEvictions    int64 `json:"size_evictions"`
	ManualEvictions  int64 `json:"manual_evictions"`
	
	// Tag operations
	TagInvalidations int64 `json:"tag_invalidations"`
	TaggedEntries    int64 `json:"tagged_entries"`
	
	// Error metrics
	GetErrors        int64 `json:"get_errors"`
	SetErrors        int64 `json:"set_errors"`
	DeleteErrors     int64 `json:"delete_errors"`
	
	// Timestamp
	LastUpdated      time.Time `json:"last_updated"`
}

// CacheEvent represents cache events for monitoring
type CacheEvent struct {
	Type        string                 `json:"type"`        // hit, miss, set, delete, evict, etc.
	Key         CacheKey               `json:"key"`
	Timestamp   time.Time              `json:"timestamp"`
	Duration    time.Duration          `json:"duration"`    // operation duration
	Size        int64                  `json:"size"`        // data size
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// CacheEventType constants
const (
	CacheEventHit          = "hit"
	CacheEventMiss         = "miss"
	CacheEventSet          = "set"
	CacheEventDelete       = "delete"
	CacheEventEvict        = "evict"
	CacheEventExpire       = "expire"
	CacheEventInvalidate   = "invalidate"
	CacheEventClear        = "clear"
	CacheEventWarmup       = "warmup"
	CacheEventError        = "error"
)

// CacheObserver interface for cache event monitoring
type CacheObserver interface {
	OnCacheEvent(event CacheEvent)
}

// CacheIndexer provides indexing capabilities for efficient lookups
type CacheIndexer interface {
	IndexByDID(did string, entry *CacheEntry) error
	IndexByType(entryType string, entry *CacheEntry) error
	IndexByTag(tag string, entry *CacheEntry) error
	
	GetByDID(did string) []*CacheEntry
	GetByType(entryType string) []*CacheEntry
	GetByTag(tag string) []*CacheEntry
	
	RemoveFromIndex(key CacheKey) error
	ClearIndex() error
}

// Helper functions

// NewCacheKey creates a new cache key
func NewCacheKey(keyType, key string) CacheKey {
	return CacheKey{
		Type:    keyType,
		Key:     key,
		Version: 1,
	}
}

// NewCacheEntry creates a new cache entry
func NewCacheEntry(key CacheKey, data interface{}, ttl time.Duration) *CacheEntry {
	return &CacheEntry{
		Key:        key,
		Data:       data,
		ExpiresAt:  time.Now().Add(ttl),
		LastAccess: time.Now(),
		Metadata:   make(map[string]interface{}),
	}
}

// IsExpired checks if a cache entry is expired
func (ce *CacheEntry) IsExpired() bool {
	return time.Now().After(ce.ExpiresAt)
}

// Touch updates the last access time and increments access count
func (ce *CacheEntry) Touch() {
	ce.LastAccess = time.Now()
	ce.AccessCount++
}

// AddTag adds a tag to the cache entry
func (ce *CacheEntry) AddTag(tag string) {
	for _, t := range ce.Tags {
		if t == tag {
			return // tag already exists
		}
	}
	ce.Tags = append(ce.Tags, tag)
}

// HasTag checks if the entry has a specific tag
func (ce *CacheEntry) HasTag(tag string) bool {
	for _, t := range ce.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

// GetHitRatio calculates the cache hit ratio
func (cs *CacheStats) GetHitRatio() float64 {
	total := cs.CacheHits + cs.CacheMisses
	if total == 0 {
		return 0.0
	}
	return float64(cs.CacheHits) / float64(total)
}

// GetEfficiency calculates cache efficiency (hits per eviction)
func (cs *CacheStats) GetEfficiency() float64 {
	if cs.EvictionCount == 0 {
		return float64(cs.CacheHits)
	}
	return float64(cs.CacheHits) / float64(cs.EvictionCount)
}

// String returns the eviction policy as a string
func (ep EvictionPolicy) String() string {
	switch ep {
	case EvictionPolicy_LRU:
		return "LRU"
	case EvictionPolicy_LFU:
		return "LFU"
	case EvictionPolicy_TTL:
		return "TTL"
	case EvictionPolicy_SIZE:
		return "SIZE"
	case EvictionPolicy_HYBRID:
		return "HYBRID"
	default:
		return "UNKNOWN"
	}
}

// Default cache configurations
var (
	DefaultCacheConfig = CacheConfig{
		MaxSize:         100 * 1024 * 1024, // 100MB
		MaxEntries:      10000,
		DefaultTTL:      30 * time.Minute,
		CleanupInterval: 5 * time.Minute,
		EnableMetrics:   true,
		EnableTags:      true,
		IdentityTTL:     1 * time.Hour,
		CredentialTTL:   30 * time.Minute,
		DIDDocumentTTL:  2 * time.Hour,
		ConsentTTL:      15 * time.Minute,
		ZKProofTTL:      10 * time.Minute,
		PreloadIdentities: true,
		WarmupOnStart:     true,
		AsyncEviction:     true,
	}
	
	HighPerformanceCacheConfig = CacheConfig{
		MaxSize:         500 * 1024 * 1024, // 500MB
		MaxEntries:      50000,
		DefaultTTL:      1 * time.Hour,
		CleanupInterval: 2 * time.Minute,
		EnableMetrics:   true,
		EnableTags:      true,
		IdentityTTL:     4 * time.Hour,
		CredentialTTL:   2 * time.Hour,
		DIDDocumentTTL:  8 * time.Hour,
		ConsentTTL:      1 * time.Hour,
		ZKProofTTL:      30 * time.Minute,
		PreloadIdentities: true,
		WarmupOnStart:     true,
		AsyncEviction:     true,
	}
	
	LowMemoryCacheConfig = CacheConfig{
		MaxSize:         10 * 1024 * 1024, // 10MB
		MaxEntries:      1000,
		DefaultTTL:      10 * time.Minute,
		CleanupInterval: 1 * time.Minute,
		EnableMetrics:   false,
		EnableTags:      false,
		IdentityTTL:     15 * time.Minute,
		CredentialTTL:   10 * time.Minute,
		DIDDocumentTTL:  20 * time.Minute,
		ConsentTTL:      5 * time.Minute,
		ZKProofTTL:      5 * time.Minute,
		PreloadIdentities: false,
		WarmupOnStart:     false,
		AsyncEviction:     false,
	}
)

// Cache key type constants
const (
	CacheTypeIdentity      = "identity"
	CacheTypeCredential    = "credential"
	CacheTypeDIDDocument   = "did_document"
	CacheTypeConsent       = "consent"
	CacheTypeZKProof       = "zk_proof"
	CacheTypeBiometric     = "biometric"
	CacheTypeShareRequest  = "share_request"
	CacheTypeShareResponse = "share_response"
	CacheTypeAccessPolicy = "access_policy"
	CacheTypeRecoveryMethod = "recovery_method"
)