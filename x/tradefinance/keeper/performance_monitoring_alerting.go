package keeper

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/prometheus/client_golang/prometheus"
)

// PerformanceMonitoringSystem manages comprehensive performance monitoring
type PerformanceMonitoringSystem struct {
	keeper                Keeper
	metricsCollector      *MetricsCollector
	performanceAnalyzer   *PerformanceAnalyzer
	alertingEngine        *AlertingEngine
	dashboardManager      *DashboardManager
	anomalyDetector       *PerformanceAnomalyDetector
	capacityPlanner       *CapacityPlanningEngine
	optimizationEngine    *PerformanceOptimizationEngine
	reportGenerator       *PerformanceReportGenerator
	mu                    sync.RWMutex
}

// MetricsCollector collects performance metrics
type MetricsCollector struct {
	systemMetrics         *SystemMetricsCollector
	applicationMetrics    *ApplicationMetricsCollector
	businessMetrics       *BusinessMetricsCollector
	customMetrics         map[string]*CustomMetric
	aggregator            *MetricsAggregator
	timeSeries            *TimeSeriesDatabase
	exporters             []MetricsExporter
}

// SystemMetricsCollector collects system-level metrics
type SystemMetricsCollector struct {
	cpuMonitor            *CPUMonitor
	memoryMonitor         *MemoryMonitor
	diskMonitor           *DiskMonitor
	networkMonitor        *NetworkMonitor
	processMonitor        *ProcessMonitor
	containerMonitor      *ContainerMonitor
}

// ApplicationMetricsCollector collects application metrics
type ApplicationMetricsCollector struct {
	transactionMetrics    *TransactionMetricsCollector
	apiMetrics            *APIMetricsCollector
	databaseMetrics       *DatabaseMetricsCollector
	cacheMetrics          *CacheMetricsCollector
	queueMetrics          *QueueMetricsCollector
	blockchainMetrics     *BlockchainMetricsCollector
}

// AlertingEngine manages alerts
type AlertingEngine struct {
	alertRules            map[string]*AlertRule
	alertProcessor        *AlertProcessor
	notificationChannels  map[string]NotificationChannel
	escalationPolicies    map[string]*EscalationPolicy
	alertHistory          *AlertHistoryManager
	suppressionRules      map[string]*SuppressionRule
	correlationEngine     *AlertCorrelationEngine
}

// AlertRule defines conditions for alerts
type AlertRule struct {
	RuleID                string
	Name                  string
	Description           string
	MetricName            string
	Condition             AlertCondition
	Threshold             float64
	Duration              time.Duration
	Severity              AlertSeverity
	Labels                map[string]string
	Annotations           map[string]string
	NotificationChannels  []string
	EscalationPolicy      string
	Active                bool
	LastEvaluation        time.Time
	LastState             AlertState
	ConsecutiveFailures   int
}

// DashboardManager manages monitoring dashboards
type DashboardManager struct {
	dashboards            map[string]*Dashboard
	widgetLibrary         *WidgetLibrary
	layoutEngine          *LayoutEngine
	refreshScheduler      *RefreshScheduler
	sharingManager        *DashboardSharingManager
	templateEngine        *DashboardTemplateEngine
}

// Dashboard represents a monitoring dashboard
type Dashboard struct {
	DashboardID           string
	Name                  string
	Description           string
	Owner                 string
	Widgets               []Widget
	Layout                DashboardLayout
	RefreshInterval       time.Duration
	TimeRange             TimeRange
	Variables             map[string]DashboardVariable
	Permissions           DashboardPermissions
	Tags                  []string
	LastModified          time.Time
	Version               int
}

// PerformanceAnomalyDetector detects performance anomalies
type PerformanceAnomalyDetector struct {
	baselineCalculator    *BaselineCalculator
	anomalyModels         map[string]AnomalyModel
	detectionEngine       *AnomalyDetectionEngine
	seasonalityAnalyzer   *SeasonalityAnalyzer
	trendAnalyzer         *TrendAnalyzer
	forecastingEngine     *ForecastingEngine
}

// CapacityPlanningEngine handles capacity planning
type CapacityPlanningEngine struct {
	resourceAnalyzer      *ResourceUtilizationAnalyzer
	growthPredictor       *GrowthPredictor
	bottleneckIdentifier  *BottleneckIdentifier
	scalingRecommender    *ScalingRecommender
	costOptimizer         *CostOptimizationEngine
	scenarioPlanner       *ScenarioPlanner
}

// Types and enums
type MetricType int
type AlertSeverity int
type AlertState int
type AlertCondition int
type AnomalyType int
type OptimizationType int
type NotificationPriority int

const (
	// Metric Types
	CounterMetric MetricType = iota
	GaugeMetric
	HistogramMetric
	SummaryMetric
	
	// Alert Severities
	InfoSeverity AlertSeverity = iota
	WarningSeverity
	ErrorSeverity
	CriticalSeverity
	
	// Alert States
	AlertInactive AlertState = iota
	AlertPending
	AlertFiring
	AlertResolved
	
	// Alert Conditions
	GreaterThan AlertCondition = iota
	LessThan
	Equal
	NotEqual
	Contains
	Regex
)

// Core monitoring methods

// CollectMetrics collects all system metrics
func (k Keeper) CollectMetrics(ctx context.Context) (*MetricsSnapshot, error) {
	pms := k.getPerformanceMonitoringSystem()
	
	snapshot := &MetricsSnapshot{
		SnapshotID:   generateID("METRIC"),
		Timestamp:    time.Now(),
		Metrics:      make(map[string]MetricValue),
		SystemHealth: HealthStatus{},
	}
	
	// Collect system metrics
	systemMetrics := pms.collectSystemMetrics()
	for name, value := range systemMetrics {
		snapshot.Metrics[name] = value
	}
	
	// Collect application metrics
	appMetrics := pms.collectApplicationMetrics(ctx)
	for name, value := range appMetrics {
		snapshot.Metrics[name] = value
	}
	
	// Collect business metrics
	businessMetrics := pms.collectBusinessMetrics(ctx)
	for name, value := range businessMetrics {
		snapshot.Metrics[name] = value
	}
	
	// Calculate system health
	snapshot.SystemHealth = pms.calculateSystemHealth(snapshot.Metrics)
	
	// Store metrics in time series database
	if err := pms.metricsCollector.timeSeries.store(snapshot); err != nil {
		return nil, fmt.Errorf("failed to store metrics: %w", err)
	}
	
	// Check for anomalies
	anomalies := pms.anomalyDetector.detectAnomalies(snapshot)
	if len(anomalies) > 0 {
		snapshot.Anomalies = anomalies
		pms.handleAnomalies(anomalies)
	}
	
	// Export metrics to external systems
	for _, exporter := range pms.metricsCollector.exporters {
		go exporter.Export(snapshot)
	}
	
	return snapshot, nil
}

// EvaluateAlerts evaluates all alert rules
func (k Keeper) EvaluateAlerts(ctx context.Context) (*AlertEvaluationResult, error) {
	pms := k.getPerformanceMonitoringSystem()
	
	result := &AlertEvaluationResult{
		EvaluationID:   generateID("ALERTEVAL"),
		Timestamp:      time.Now(),
		RulesEvaluated: 0,
		AlertsFired:    0,
		AlertsResolved: 0,
		Errors:         []error{},
	}
	
	// Get current metrics
	metrics, err := k.CollectMetrics(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to collect metrics: %w", err)
	}
	
	// Evaluate each alert rule
	for _, rule := range pms.alertingEngine.alertRules {
		if !rule.Active {
			continue
		}
		
		result.RulesEvaluated++
		
		// Evaluate rule condition
		evaluation := pms.evaluateAlertRule(rule, metrics)
		
		// Update rule state
		previousState := rule.LastState
		rule.LastEvaluation = time.Now()
		rule.LastState = evaluation.State
		
		// Handle state transitions
		switch {
		case previousState != AlertFiring && evaluation.State == AlertFiring:
			// New alert firing
			alert := pms.createAlert(rule, evaluation, metrics)
			pms.alertingEngine.alertProcessor.processAlert(alert)
			result.AlertsFired++
			
		case previousState == AlertFiring && evaluation.State != AlertFiring:
			// Alert resolved
			pms.resolveAlert(rule, evaluation)
			result.AlertsResolved++
			
		case previousState == AlertFiring && evaluation.State == AlertFiring:
			// Alert still firing - check for re-notification
			pms.checkRenotification(rule, evaluation)
		}
		
		// Store evaluation result
		pms.storeAlertEvaluation(rule, evaluation)
	}
	
	// Correlate alerts
	correlations := pms.alertingEngine.correlationEngine.correlateAlerts()
	result.Correlations = correlations
	
	// Apply suppression rules
	suppressions := pms.applySuppresionRules()
	result.Suppressions = suppressions
	
	return result, nil
}

// System metrics collection

func (pms *PerformanceMonitoringSystem) collectSystemMetrics() map[string]MetricValue {
	metrics := make(map[string]MetricValue)
	
	// CPU metrics
	cpuMetrics := pms.metricsCollector.systemMetrics.cpuMonitor.collect()
	metrics["system.cpu.usage"] = MetricValue{
		Value:     cpuMetrics.Usage,
		Type:      GaugeMetric,
		Unit:      "percent",
		Timestamp: time.Now(),
	}
	metrics["system.cpu.load.1m"] = MetricValue{
		Value:     cpuMetrics.Load1m,
		Type:      GaugeMetric,
		Unit:      "load",
		Timestamp: time.Now(),
	}
	
	// Memory metrics
	memoryMetrics := pms.metricsCollector.systemMetrics.memoryMonitor.collect()
	metrics["system.memory.used"] = MetricValue{
		Value:     float64(memoryMetrics.Used),
		Type:      GaugeMetric,
		Unit:      "bytes",
		Timestamp: time.Now(),
	}
	metrics["system.memory.available"] = MetricValue{
		Value:     float64(memoryMetrics.Available),
		Type:      GaugeMetric,
		Unit:      "bytes",
		Timestamp: time.Now(),
	}
	metrics["system.memory.usage"] = MetricValue{
		Value:     memoryMetrics.UsagePercent,
		Type:      GaugeMetric,
		Unit:      "percent",
		Timestamp: time.Now(),
	}
	
	// Disk metrics
	diskMetrics := pms.metricsCollector.systemMetrics.diskMonitor.collect()
	for _, disk := range diskMetrics {
		prefix := fmt.Sprintf("system.disk.%s", disk.Device)
		metrics[prefix+".used"] = MetricValue{
			Value:     float64(disk.Used),
			Type:      GaugeMetric,
			Unit:      "bytes",
			Timestamp: time.Now(),
		}
		metrics[prefix+".free"] = MetricValue{
			Value:     float64(disk.Free),
			Type:      GaugeMetric,
			Unit:      "bytes",
			Timestamp: time.Now(),
		}
		metrics[prefix+".usage"] = MetricValue{
			Value:     disk.UsagePercent,
			Type:      GaugeMetric,
			Unit:      "percent",
			Timestamp: time.Now(),
		}
		metrics[prefix+".iops.read"] = MetricValue{
			Value:     float64(disk.ReadIOPS),
			Type:      GaugeMetric,
			Unit:      "ops/s",
			Timestamp: time.Now(),
		}
		metrics[prefix+".iops.write"] = MetricValue{
			Value:     float64(disk.WriteIOPS),
			Type:      GaugeMetric,
			Unit:      "ops/s",
			Timestamp: time.Now(),
		}
	}
	
	// Network metrics
	networkMetrics := pms.metricsCollector.systemMetrics.networkMonitor.collect()
	for _, iface := range networkMetrics {
		prefix := fmt.Sprintf("system.network.%s", iface.Interface)
		metrics[prefix+".bytes.sent"] = MetricValue{
			Value:     float64(iface.BytesSent),
			Type:      CounterMetric,
			Unit:      "bytes",
			Timestamp: time.Now(),
		}
		metrics[prefix+".bytes.received"] = MetricValue{
			Value:     float64(iface.BytesReceived),
			Type:      CounterMetric,
			Unit:      "bytes",
			Timestamp: time.Now(),
		}
		metrics[prefix+".packets.sent"] = MetricValue{
			Value:     float64(iface.PacketsSent),
			Type:      CounterMetric,
			Unit:      "packets",
			Timestamp: time.Now(),
		}
		metrics[prefix+".packets.received"] = MetricValue{
			Value:     float64(iface.PacketsReceived),
			Type:      CounterMetric,
			Unit:      "packets",
			Timestamp: time.Now(),
		}
		metrics[prefix+".errors"] = MetricValue{
			Value:     float64(iface.Errors),
			Type:      CounterMetric,
			Unit:      "errors",
			Timestamp: time.Now(),
		}
	}
	
	return metrics
}

// Application metrics collection

func (pms *PerformanceMonitoringSystem) collectApplicationMetrics(ctx context.Context) map[string]MetricValue {
	metrics := make(map[string]MetricValue)
	
	// Transaction metrics
	txMetrics := pms.metricsCollector.applicationMetrics.transactionMetrics.collect()
	metrics["app.transactions.total"] = MetricValue{
		Value:     float64(txMetrics.Total),
		Type:      CounterMetric,
		Unit:      "transactions",
		Timestamp: time.Now(),
	}
	metrics["app.transactions.rate"] = MetricValue{
		Value:     txMetrics.Rate,
		Type:      GaugeMetric,
		Unit:      "tx/s",
		Timestamp: time.Now(),
	}
	metrics["app.transactions.latency.p50"] = MetricValue{
		Value:     txMetrics.LatencyP50.Seconds() * 1000,
		Type:      GaugeMetric,
		Unit:      "ms",
		Timestamp: time.Now(),
	}
	metrics["app.transactions.latency.p95"] = MetricValue{
		Value:     txMetrics.LatencyP95.Seconds() * 1000,
		Type:      GaugeMetric,
		Unit:      "ms",
		Timestamp: time.Now(),
	}
	metrics["app.transactions.latency.p99"] = MetricValue{
		Value:     txMetrics.LatencyP99.Seconds() * 1000,
		Type:      GaugeMetric,
		Unit:      "ms",
		Timestamp: time.Now(),
	}
	
	// API metrics
	apiMetrics := pms.metricsCollector.applicationMetrics.apiMetrics.collect()
	for endpoint, stats := range apiMetrics {
		prefix := fmt.Sprintf("app.api.%s", sanitizeMetricName(endpoint))
		metrics[prefix+".requests"] = MetricValue{
			Value:     float64(stats.Requests),
			Type:      CounterMetric,
			Unit:      "requests",
			Timestamp: time.Now(),
		}
		metrics[prefix+".errors"] = MetricValue{
			Value:     float64(stats.Errors),
			Type:      CounterMetric,
			Unit:      "errors",
			Timestamp: time.Now(),
		}
		metrics[prefix+".latency.avg"] = MetricValue{
			Value:     stats.AvgLatency.Seconds() * 1000,
			Type:      GaugeMetric,
			Unit:      "ms",
			Timestamp: time.Now(),
		}
		metrics[prefix+".success_rate"] = MetricValue{
			Value:     stats.SuccessRate,
			Type:      GaugeMetric,
			Unit:      "percent",
			Timestamp: time.Now(),
		}
	}
	
	// Database metrics
	dbMetrics := pms.metricsCollector.applicationMetrics.databaseMetrics.collect()
	metrics["app.database.connections.active"] = MetricValue{
		Value:     float64(dbMetrics.ActiveConnections),
		Type:      GaugeMetric,
		Unit:      "connections",
		Timestamp: time.Now(),
	}
	metrics["app.database.connections.idle"] = MetricValue{
		Value:     float64(dbMetrics.IdleConnections),
		Type:      GaugeMetric,
		Unit:      "connections",
		Timestamp: time.Now(),
	}
	metrics["app.database.queries.total"] = MetricValue{
		Value:     float64(dbMetrics.TotalQueries),
		Type:      CounterMetric,
		Unit:      "queries",
		Timestamp: time.Now(),
	}
	metrics["app.database.queries.slow"] = MetricValue{
		Value:     float64(dbMetrics.SlowQueries),
		Type:      CounterMetric,
		Unit:      "queries",
		Timestamp: time.Now(),
	}
	metrics["app.database.query.latency.avg"] = MetricValue{
		Value:     dbMetrics.AvgQueryTime.Seconds() * 1000,
		Type:      GaugeMetric,
		Unit:      "ms",
		Timestamp: time.Now(),
	}
	
	// Cache metrics
	cacheMetrics := pms.metricsCollector.applicationMetrics.cacheMetrics.collect()
	metrics["app.cache.hits"] = MetricValue{
		Value:     float64(cacheMetrics.Hits),
		Type:      CounterMetric,
		Unit:      "hits",
		Timestamp: time.Now(),
	}
	metrics["app.cache.misses"] = MetricValue{
		Value:     float64(cacheMetrics.Misses),
		Type:      CounterMetric,
		Unit:      "misses",
		Timestamp: time.Now(),
	}
	metrics["app.cache.hit_rate"] = MetricValue{
		Value:     cacheMetrics.HitRate,
		Type:      GaugeMetric,
		Unit:      "percent",
		Timestamp: time.Now(),
	}
	metrics["app.cache.evictions"] = MetricValue{
		Value:     float64(cacheMetrics.Evictions),
		Type:      CounterMetric,
		Unit:      "evictions",
		Timestamp: time.Now(),
	}
	metrics["app.cache.size"] = MetricValue{
		Value:     float64(cacheMetrics.Size),
		Type:      GaugeMetric,
		Unit:      "bytes",
		Timestamp: time.Now(),
	}
	
	// Blockchain metrics
	bcMetrics := pms.metricsCollector.applicationMetrics.blockchainMetrics.collect()
	metrics["blockchain.height"] = MetricValue{
		Value:     float64(bcMetrics.Height),
		Type:      GaugeMetric,
		Unit:      "blocks",
		Timestamp: time.Now(),
	}
	metrics["blockchain.validators.active"] = MetricValue{
		Value:     float64(bcMetrics.ActiveValidators),
		Type:      GaugeMetric,
		Unit:      "validators",
		Timestamp: time.Now(),
	}
	metrics["blockchain.transactions.pending"] = MetricValue{
		Value:     float64(bcMetrics.PendingTx),
		Type:      GaugeMetric,
		Unit:      "transactions",
		Timestamp: time.Now(),
	}
	metrics["blockchain.block.time"] = MetricValue{
		Value:     bcMetrics.BlockTime.Seconds(),
		Type:      GaugeMetric,
		Unit:      "seconds",
		Timestamp: time.Now(),
	}
	
	return metrics
}

// Anomaly detection

func (pad *PerformanceAnomalyDetector) detectAnomalies(snapshot *MetricsSnapshot) []Anomaly {
	anomalies := []Anomaly{}
	
	for metricName, metricValue := range snapshot.Metrics {
		// Get baseline for metric
		baseline := pad.baselineCalculator.getBaseline(metricName, snapshot.Timestamp)
		if baseline == nil {
			continue
		}
		
		// Check for statistical anomaly
		if anomaly := pad.checkStatisticalAnomaly(metricName, metricValue, baseline); anomaly != nil {
			anomalies = append(anomalies, *anomaly)
		}
		
		// Check for seasonal anomaly
		if anomaly := pad.seasonalityAnalyzer.checkSeasonalAnomaly(metricName, metricValue, snapshot.Timestamp); anomaly != nil {
			anomalies = append(anomalies, *anomaly)
		}
		
		// Check for trend anomaly
		if anomaly := pad.trendAnalyzer.checkTrendAnomaly(metricName, metricValue); anomaly != nil {
			anomalies = append(anomalies, *anomaly)
		}
		
		// Use ML models for complex anomaly detection
		for _, model := range pad.anomalyModels {
			if model.AppliesTo(metricName) {
				if anomaly := model.Detect(metricValue, snapshot); anomaly != nil {
					anomalies = append(anomalies, *anomaly)
				}
			}
		}
	}
	
	// Correlate anomalies
	correlatedAnomalies := pad.correlateAnomalies(anomalies)
	
	return correlatedAnomalies
}

// Performance optimization

func (poe *PerformanceOptimizationEngine) optimizePerformance(metrics map[string]MetricValue) []OptimizationRecommendation {
	recommendations := []OptimizationRecommendation{}
	
	// Check CPU optimization opportunities
	if cpuUsage, exists := metrics["system.cpu.usage"]; exists && cpuUsage.Value > 80 {
		rec := poe.analyzeCPUOptimization(metrics)
		if rec != nil {
			recommendations = append(recommendations, *rec)
		}
	}
	
	// Check memory optimization
	if memUsage, exists := metrics["system.memory.usage"]; exists && memUsage.Value > 85 {
		rec := poe.analyzeMemoryOptimization(metrics)
		if rec != nil {
			recommendations = append(recommendations, *rec)
		}
	}
	
	// Check database optimization
	if slowQueries, exists := metrics["app.database.queries.slow"]; exists && slowQueries.Value > 100 {
		rec := poe.analyzeDatabaseOptimization(metrics)
		if rec != nil {
			recommendations = append(recommendations, *rec)
		}
	}
	
	// Check cache optimization
	if hitRate, exists := metrics["app.cache.hit_rate"]; exists && hitRate.Value < 80 {
		rec := poe.analyzeCacheOptimization(metrics)
		if rec != nil {
			recommendations = append(recommendations, *rec)
		}
	}
	
	// Sort by impact
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].ExpectedImpact > recommendations[j].ExpectedImpact
	})
	
	return recommendations
}

// Helper types

type MetricsSnapshot struct {
	SnapshotID   string
	Timestamp    time.Time
	Metrics      map[string]MetricValue
	SystemHealth HealthStatus
	Anomalies    []Anomaly
}

type MetricValue struct {
	Value     float64
	Type      MetricType
	Unit      string
	Timestamp time.Time
	Labels    map[string]string
}

type AlertEvaluationResult struct {
	EvaluationID   string
	Timestamp      time.Time
	RulesEvaluated int
	AlertsFired    int
	AlertsResolved int
	Correlations   []AlertCorrelation
	Suppressions   []AlertSuppression
	Errors         []error
}

type Anomaly struct {
	AnomalyID    string
	MetricName   string
	Type         AnomalyType
	Severity     float64
	Description  string
	DetectedAt   time.Time
	ExpectedValue float64
	ActualValue  float64
	Confidence   float64
}

type OptimizationRecommendation struct {
	RecommendationID string
	Type             OptimizationType
	Title            string
	Description      string
	ExpectedImpact   float64
	Implementation   string
	Risk             string
	Priority         int
}

type HealthStatus struct {
	Overall      string
	Score        float64
	Components   map[string]ComponentHealth
	LastChecked  time.Time
}

type ComponentHealth struct {
	Name     string
	Status   string
	Score    float64
	Message  string
	Metrics  map[string]float64
}

type Alert struct {
	AlertID       string
	RuleID        string
	RuleName      string
	Severity      AlertSeverity
	State         AlertState
	Value         float64
	Threshold     float64
	Description   string
	Labels        map[string]string
	Annotations   map[string]string
	FiredAt       time.Time
	ResolvedAt    *time.Time
	Fingerprint   string
}

type Widget struct {
	WidgetID     string
	Type         string
	Title        string
	Query        string
	Visualization string
	Options      map[string]interface{}
	Position     WidgetPosition
	Size         WidgetSize
}

type DashboardLayout struct {
	Type    string
	Columns int
	Rows    int
	Widgets []WidgetPosition
}

type BaselineMetric struct {
	MetricName   string
	Mean         float64
	StdDev       float64
	Min          float64
	Max          float64
	Percentiles  map[int]float64
	LastUpdated  time.Time
}

// CPU Monitor implementation
type CPUMonitor struct {
	previousCPU  uint64
	previousTime time.Time
}

type CPUMetrics struct {
	Usage   float64
	Load1m  float64
	Load5m  float64
	Load15m float64
	Cores   int
}

func (cm *CPUMonitor) collect() CPUMetrics {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	
	return CPUMetrics{
		Usage:  calculateCPUUsage(),
		Load1m: getLoadAverage(1),
		Load5m: getLoadAverage(5),
		Load15m: getLoadAverage(15),
		Cores:  runtime.NumCPU(),
	}
}

// Memory Monitor implementation
type MemoryMonitor struct{}

type MemoryMetrics struct {
	Total        uint64
	Used         uint64
	Available    uint64
	UsagePercent float64
	SwapTotal    uint64
	SwapUsed     uint64
}

func (mm *MemoryMonitor) collect() MemoryMetrics {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	
	return MemoryMetrics{
		Used:         stats.Alloc,
		Total:        stats.Sys,
		Available:    stats.Sys - stats.Alloc,
		UsagePercent: float64(stats.Alloc) / float64(stats.Sys) * 100,
	}
}

// Utility functions

func sanitizeMetricName(name string) string {
	// Replace invalid characters with underscores
	sanitized := ""
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '.' {
			sanitized += string(r)
		} else {
			sanitized += "_"
		}
	}
	return sanitized
}

func calculateCPUUsage() float64 {
	// Simplified CPU usage calculation
	// In production, use proper OS-specific APIs
	return math.Min(math.Max(0, math.Sin(float64(time.Now().Unix()))*50+50), 100)
}

func getLoadAverage(minutes int) float64 {
	// Simplified load average
	// In production, read from /proc/loadavg or equivalent
	base := 0.5
	variance := 0.3
	return base + variance*math.Sin(float64(time.Now().Unix())/float64(minutes*60))
}

// Prometheus metrics registration
var (
	transactionCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "deshchain_transactions_total",
			Help: "Total number of transactions processed",
		},
		[]string{"type", "status"},
	)
	
	transactionDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "deshchain_transaction_duration_seconds",
			Help:    "Transaction processing duration",
			Buckets: prometheus.ExponentialBuckets(0.001, 2, 10),
		},
		[]string{"type"},
	)
	
	systemMemoryGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "deshchain_system_memory_bytes",
			Help: "System memory usage in bytes",
		},
	)
	
	apiRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "deshchain_api_requests_total",
			Help: "Total API requests",
		},
		[]string{"endpoint", "method", "status"},
	)
)

func init() {
	// Register Prometheus metrics
	prometheus.MustRegister(transactionCounter)
	prometheus.MustRegister(transactionDuration)
	prometheus.MustRegister(systemMemoryGauge)
	prometheus.MustRegister(apiRequestsTotal)
}