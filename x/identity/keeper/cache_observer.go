package keeper

import (
	"fmt"
	"sync"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/DeshChain/DeshChain-Ecosystem/x/identity/types"
)

// CacheObserverManager manages multiple cache observers
type CacheObserverManager struct {
	mu        sync.RWMutex
	observers []types.CacheObserver
	enabled   bool
}

// NewCacheObserverManager creates a new observer manager
func NewCacheObserverManager() *CacheObserverManager {
	return &CacheObserverManager{
		observers: make([]types.CacheObserver, 0),
		enabled:   true,
	}
}

// AddObserver adds a cache observer
func (com *CacheObserverManager) AddObserver(observer types.CacheObserver) {
	com.mu.Lock()
	defer com.mu.Unlock()
	com.observers = append(com.observers, observer)
}

// RemoveObserver removes a cache observer
func (com *CacheObserverManager) RemoveObserver(observer types.CacheObserver) {
	com.mu.Lock()
	defer com.mu.Unlock()
	
	for i, obs := range com.observers {
		if obs == observer {
			com.observers = append(com.observers[:i], com.observers[i+1:]...)
			break
		}
	}
}

// NotifyObservers notifies all observers of an event
func (com *CacheObserverManager) NotifyObservers(event types.CacheEvent) {
	if !com.enabled {
		return
	}
	
	com.mu.RLock()
	observers := make([]types.CacheObserver, len(com.observers))
	copy(observers, com.observers)
	com.mu.RUnlock()
	
	for _, observer := range observers {
		go observer.OnCacheEvent(event) // Async notification
	}
}

// Enable enables observer notifications
func (com *CacheObserverManager) Enable() {
	com.mu.Lock()
	defer com.mu.Unlock()
	com.enabled = true
}

// Disable disables observer notifications
func (com *CacheObserverManager) Disable() {
	com.mu.Lock()
	defer com.mu.Unlock()
	com.enabled = false
}

// MetricsObserver collects and aggregates cache metrics
type MetricsObserver struct {
	mu                sync.RWMutex
	hitCounters       map[string]int64
	missCounters      map[string]int64
	operationTimes    map[string][]float64
	lastReset         time.Time
	eventHistory      []types.CacheEvent
	maxHistorySize    int
	alertThresholds   MetricsThresholds
	alerts            []CacheAlert
}

// MetricsThresholds defines alert thresholds for cache metrics
type MetricsThresholds struct {
	MinHitRatio        float64 // Alert if hit ratio falls below this
	MaxOperationTime   float64 // Alert if operation time exceeds this (ms)
	MaxMemoryUsage     int64   // Alert if memory usage exceeds this (bytes)
	MaxMissRate        float64 // Alert if miss rate exceeds this (misses per second)
}

// CacheAlert represents a cache performance alert
type CacheAlert struct {
	Type        string    `json:"type"`
	Severity    string    `json:"severity"`
	Message     string    `json:"message"`
	Timestamp   time.Time `json:"timestamp"`
	MetricValue float64   `json:"metric_value"`
	Threshold   float64   `json:"threshold"`
}

// NewMetricsObserver creates a new metrics observer
func NewMetricsObserver() *MetricsObserver {
	return &MetricsObserver{
		hitCounters:     make(map[string]int64),
		missCounters:    make(map[string]int64),
		operationTimes:  make(map[string][]float64),
		lastReset:       time.Now(),
		eventHistory:    make([]types.CacheEvent, 0),
		maxHistorySize:  1000,
		alertThresholds: MetricsThresholds{
			MinHitRatio:      0.7,  // 70%
			MaxOperationTime: 10.0, // 10ms
			MaxMemoryUsage:   100 * 1024 * 1024, // 100MB
			MaxMissRate:      100.0, // 100 misses per second
		},
		alerts: make([]CacheAlert, 0),
	}
}

// OnCacheEvent processes cache events
func (mo *MetricsObserver) OnCacheEvent(event types.CacheEvent) {
	mo.mu.Lock()
	defer mo.mu.Unlock()
	
	// Update counters
	switch event.Type {
	case types.CacheEventHit:
		mo.hitCounters[event.Key.Type]++
		mo.hitCounters["total"]++
	case types.CacheEventMiss:
		mo.missCounters[event.Key.Type]++
		mo.missCounters["total"]++
	}
	
	// Track operation times
	if event.Duration > 0 {
		durationMs := float64(event.Duration.Nanoseconds()) / 1e6
		mo.operationTimes[event.Type] = append(mo.operationTimes[event.Type], durationMs)
		
		// Keep only last 100 measurements per operation type
		if len(mo.operationTimes[event.Type]) > 100 {
			mo.operationTimes[event.Type] = mo.operationTimes[event.Type][1:]
		}
		
		// Check for performance alerts
		mo.checkPerformanceAlerts(event.Type, durationMs)
	}
	
	// Add to event history
	mo.eventHistory = append(mo.eventHistory, event)
	if len(mo.eventHistory) > mo.maxHistorySize {
		mo.eventHistory = mo.eventHistory[1:]
	}
}

// GetHitRatio returns the cache hit ratio for a specific type
func (mo *MetricsObserver) GetHitRatio(entryType string) float64 {
	mo.mu.RLock()
	defer mo.mu.RUnlock()
	
	hits := mo.hitCounters[entryType]
	misses := mo.missCounters[entryType]
	total := hits + misses
	
	if total == 0 {
		return 0.0
	}
	
	return float64(hits) / float64(total)
}

// GetAverageOperationTime returns the average operation time for a type
func (mo *MetricsObserver) GetAverageOperationTime(operationType string) float64 {
	mo.mu.RLock()
	defer mo.mu.RUnlock()
	
	times, exists := mo.operationTimes[operationType]
	if !exists || len(times) == 0 {
		return 0.0
	}
	
	var sum float64
	for _, time := range times {
		sum += time
	}
	
	return sum / float64(len(times))
}

// GetRecentEvents returns recent cache events
func (mo *MetricsObserver) GetRecentEvents(limit int) []types.CacheEvent {
	mo.mu.RLock()
	defer mo.mu.RUnlock()
	
	if limit <= 0 || limit > len(mo.eventHistory) {
		limit = len(mo.eventHistory)
	}
	
	start := len(mo.eventHistory) - limit
	events := make([]types.CacheEvent, limit)
	copy(events, mo.eventHistory[start:])
	
	return events
}

// GetAlerts returns recent cache alerts
func (mo *MetricsObserver) GetAlerts() []CacheAlert {
	mo.mu.RLock()
	defer mo.mu.RUnlock()
	
	alerts := make([]CacheAlert, len(mo.alerts))
	copy(alerts, mo.alerts)
	
	return alerts
}

// Reset resets metrics counters
func (mo *MetricsObserver) Reset() {
	mo.mu.Lock()
	defer mo.mu.Unlock()
	
	mo.hitCounters = make(map[string]int64)
	mo.missCounters = make(map[string]int64)
	mo.operationTimes = make(map[string][]float64)
	mo.lastReset = time.Now()
	mo.eventHistory = make([]types.CacheEvent, 0)
	mo.alerts = make([]CacheAlert, 0)
}

// checkPerformanceAlerts checks for performance issues and creates alerts
func (mo *MetricsObserver) checkPerformanceAlerts(operationType string, durationMs float64) {
	// Check operation time threshold
	if durationMs > mo.alertThresholds.MaxOperationTime {
		alert := CacheAlert{
			Type:        "slow_operation",
			Severity:    "warning",
			Message:     fmt.Sprintf("Slow %s operation: %.2fms (threshold: %.2fms)", operationType, durationMs, mo.alertThresholds.MaxOperationTime),
			Timestamp:   time.Now(),
			MetricValue: durationMs,
			Threshold:   mo.alertThresholds.MaxOperationTime,
		}
		mo.alerts = append(mo.alerts, alert)
		
		// Keep only last 50 alerts
		if len(mo.alerts) > 50 {
			mo.alerts = mo.alerts[1:]
		}
	}
	
	// Check hit ratio
	hitRatio := mo.GetHitRatio("total")
	if hitRatio < mo.alertThresholds.MinHitRatio && mo.hitCounters["total"]+mo.missCounters["total"] > 100 {
		alert := CacheAlert{
			Type:        "low_hit_ratio",
			Severity:    "warning",
			Message:     fmt.Sprintf("Low cache hit ratio: %.2f%% (threshold: %.2f%%)", hitRatio*100, mo.alertThresholds.MinHitRatio*100),
			Timestamp:   time.Now(),
			MetricValue: hitRatio,
			Threshold:   mo.alertThresholds.MinHitRatio,
		}
		mo.alerts = append(mo.alerts, alert)
	}
}

// LoggingObserver logs cache events to the blockchain logger
type LoggingObserver struct {
	ctx     sdk.Context
	enabled bool
	logLevel string
}

// NewLoggingObserver creates a new logging observer
func NewLoggingObserver(ctx sdk.Context) *LoggingObserver {
	return &LoggingObserver{
		ctx:      ctx,
		enabled:  true,
		logLevel: "info",
	}
}

// OnCacheEvent logs cache events
func (lo *LoggingObserver) OnCacheEvent(event types.CacheEvent) {
	if !lo.enabled {
		return
	}
	
	logger := lo.ctx.Logger().With("module", "identity_cache")
	
	switch event.Type {
	case types.CacheEventHit, types.CacheEventMiss:
		if lo.logLevel == "debug" {
			logger.Debug(
				fmt.Sprintf("Cache %s", event.Type),
				"key_type", event.Key.Type,
				"key", event.Key.Key,
				"duration_ms", float64(event.Duration.Nanoseconds())/1e6,
			)
		}
	case types.CacheEventSet:
		logger.Info(
			"Cache entry set",
			"key_type", event.Key.Type,
			"key", event.Key.Key,
			"size_bytes", event.Size,
		)
	case types.CacheEventEvict:
		logger.Info(
			"Cache entry evicted",
			"key_type", event.Key.Type,
			"key", event.Key.Key,
			"size_bytes", event.Size,
		)
	case types.CacheEventClear:
		logger.Info("Cache cleared")
	case types.CacheEventError:
		logger.Error(
			"Cache error",
			"key_type", event.Key.Type,
			"key", event.Key.Key,
			"error", event.Metadata["error"],
		)
	}
}

// SetLogLevel sets the logging level
func (lo *LoggingObserver) SetLogLevel(level string) {
	lo.logLevel = level
}

// Enable enables logging
func (lo *LoggingObserver) Enable() {
	lo.enabled = true
}

// Disable disables logging
func (lo *LoggingObserver) Disable() {
	lo.enabled = false
}

// AuditObserver creates audit trails for cache operations
type AuditObserver struct {
	auditLog []CacheAuditEntry
	mu       sync.RWMutex
	maxSize  int
}

// CacheAuditEntry represents an audit log entry
type CacheAuditEntry struct {
	Timestamp   time.Time              `json:"timestamp"`
	EventType   string                 `json:"event_type"`
	KeyType     string                 `json:"key_type"`
	Key         string                 `json:"key"`
	Size        int64                  `json:"size"`
	Duration    time.Duration          `json:"duration"`
	Metadata    map[string]interface{} `json:"metadata"`
	UserAddress string                 `json:"user_address,omitempty"`
}

// NewAuditObserver creates a new audit observer
func NewAuditObserver(maxSize int) *AuditObserver {
	return &AuditObserver{
		auditLog: make([]CacheAuditEntry, 0),
		maxSize:  maxSize,
	}
}

// OnCacheEvent creates audit entries for cache events
func (ao *AuditObserver) OnCacheEvent(event types.CacheEvent) {
	ao.mu.Lock()
	defer ao.mu.Unlock()
	
	// Only audit significant events
	if !ao.shouldAudit(event.Type) {
		return
	}
	
	entry := CacheAuditEntry{
		Timestamp: event.Timestamp,
		EventType: event.Type,
		KeyType:   event.Key.Type,
		Key:       event.Key.Key,
		Size:      event.Size,
		Duration:  event.Duration,
		Metadata:  event.Metadata,
	}
	
	// Extract user address from metadata if available
	if userAddr, exists := event.Metadata["user_address"]; exists {
		if addr, ok := userAddr.(string); ok {
			entry.UserAddress = addr
		}
	}
	
	ao.auditLog = append(ao.auditLog, entry)
	
	// Maintain size limit
	if len(ao.auditLog) > ao.maxSize {
		ao.auditLog = ao.auditLog[1:]
	}
}

// shouldAudit determines if an event should be audited
func (ao *AuditObserver) shouldAudit(eventType string) bool {
	switch eventType {
	case types.CacheEventSet, types.CacheEventDelete, types.CacheEventEvict, 
		 types.CacheEventClear, types.CacheEventInvalidate:
		return true
	default:
		return false
	}
}

// GetAuditLog returns the audit log entries
func (ao *AuditObserver) GetAuditLog() []CacheAuditEntry {
	ao.mu.RLock()
	defer ao.mu.RUnlock()
	
	log := make([]CacheAuditEntry, len(ao.auditLog))
	copy(log, ao.auditLog)
	
	return log
}

// GetAuditLogForUser returns audit log entries for a specific user
func (ao *AuditObserver) GetAuditLogForUser(userAddress string) []CacheAuditEntry {
	ao.mu.RLock()
	defer ao.mu.RUnlock()
	
	var userLog []CacheAuditEntry
	for _, entry := range ao.auditLog {
		if entry.UserAddress == userAddress {
			userLog = append(userLog, entry)
		}
	}
	
	return userLog
}

// PerformanceMonitor monitors cache performance and provides recommendations
type PerformanceMonitor struct {
	metricsObserver *MetricsObserver
	sampleInterval  time.Duration
	recommendations []string
	mu              sync.RWMutex
}

// NewPerformanceMonitor creates a new performance monitor
func NewPerformanceMonitor(metricsObserver *MetricsObserver) *PerformanceMonitor {
	return &PerformanceMonitor{
		metricsObserver: metricsObserver,
		sampleInterval:  5 * time.Minute,
		recommendations: make([]string, 0),
	}
}

// StartMonitoring starts performance monitoring
func (pm *PerformanceMonitor) StartMonitoring() {
	go pm.monitoringLoop()
}

// monitoringLoop runs the performance monitoring loop
func (pm *PerformanceMonitor) monitoringLoop() {
	ticker := time.NewTicker(pm.sampleInterval)
	defer ticker.Stop()
	
	for range ticker.C {
		pm.analyzePerformance()
	}
}

// analyzePerformance analyzes cache performance and generates recommendations
func (pm *PerformanceMonitor) analyzePerformance() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	pm.recommendations = pm.recommendations[:0] // Clear previous recommendations
	
	// Check hit ratio
	hitRatio := pm.metricsObserver.GetHitRatio("total")
	if hitRatio < 0.5 {
		pm.recommendations = append(pm.recommendations, 
			"Consider increasing cache size or TTL - hit ratio is below 50%")
	}
	
	// Check operation times
	avgGetTime := pm.metricsObserver.GetAverageOperationTime(types.CacheEventHit)
	if avgGetTime > 5.0 {
		pm.recommendations = append(pm.recommendations, 
			"Cache get operations are slow - consider optimizing cache structure")
	}
	
	avgSetTime := pm.metricsObserver.GetAverageOperationTime(types.CacheEventSet)
	if avgSetTime > 10.0 {
		pm.recommendations = append(pm.recommendations, 
			"Cache set operations are slow - consider reducing data size or optimizing serialization")
	}
	
	// Check for memory issues
	alerts := pm.metricsObserver.GetAlerts()
	memoryAlerts := 0
	for _, alert := range alerts {
		if alert.Type == "high_memory_usage" {
			memoryAlerts++
		}
	}
	
	if memoryAlerts > 5 {
		pm.recommendations = append(pm.recommendations, 
			"High memory usage detected - consider reducing cache size or implementing more aggressive eviction")
	}
}

// GetRecommendations returns performance recommendations
func (pm *PerformanceMonitor) GetRecommendations() []string {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	
	recommendations := make([]string, len(pm.recommendations))
	copy(recommendations, pm.recommendations)
	
	return recommendations
}