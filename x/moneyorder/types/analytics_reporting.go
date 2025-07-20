package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ReportType represents different types of analytics reports
type ReportType int32

const (
	ReportType_UNKNOWN      ReportType = 0
	ReportType_SYSTEM       ReportType = 1
	ReportType_BUSINESS     ReportType = 2
	ReportType_TRANSACTION  ReportType = 3
	ReportType_PERFORMANCE  ReportType = 4
	ReportType_COMPLIANCE   ReportType = 5
	ReportType_SECURITY     ReportType = 6
	ReportType_ANOMALY      ReportType = 7
)

// String returns the string representation of ReportType
func (rt ReportType) String() string {
	switch rt {
	case ReportType_SYSTEM:
		return "SYSTEM"
	case ReportType_BUSINESS:
		return "BUSINESS"
	case ReportType_TRANSACTION:
		return "TRANSACTION"
	case ReportType_PERFORMANCE:
		return "PERFORMANCE"
	case ReportType_COMPLIANCE:
		return "COMPLIANCE"
	case ReportType_SECURITY:
		return "SECURITY"
	case ReportType_ANOMALY:
		return "ANOMALY"
	default:
		return "UNKNOWN"
	}
}

// ExportFormat represents different export formats for reports
type ExportFormat int32

const (
	ExportFormat_JSON  ExportFormat = 0
	ExportFormat_CSV   ExportFormat = 1
	ExportFormat_PDF   ExportFormat = 2
	ExportFormat_EXCEL ExportFormat = 3
)

// AnomalySeverity represents the severity level of detected anomalies
type AnomalySeverity int32

const (
	AnomalySeverity_LOW      AnomalySeverity = 0
	AnomalySeverity_MEDIUM   AnomalySeverity = 1
	AnomalySeverity_HIGH     AnomalySeverity = 2
	AnomalySeverity_CRITICAL AnomalySeverity = 3
)

// SystemAnalyticsReport represents a comprehensive system analytics report
type SystemAnalyticsReport struct {
	ReportID     string                 `json:"report_id"`
	ReportType   ReportType            `json:"report_type"`
	StartDate    time.Time             `json:"start_date"`
	EndDate      time.Time             `json:"end_date"`
	GeneratedAt  time.Time             `json:"generated_at"`
	Summary      SystemSummary         `json:"summary"`
	Metrics      SystemMetrics         `json:"metrics"`
	Trends       TrendAnalysis         `json:"trends"`
	Geography    GeographicAnalysis    `json:"geography"`
	Performance  PerformanceMetrics    `json:"performance"`
	Predictions  PredictiveAnalysis    `json:"predictions"`
}

// BusinessAnalyticsReport represents analytics report for a specific business
type BusinessAnalyticsReport struct {
	ReportID        string                    `json:"report_id"`
	BusinessAddress string                    `json:"business_address"`
	BusinessName    string                    `json:"business_name"`
	StartDate       time.Time                 `json:"start_date"`
	EndDate         time.Time                 `json:"end_date"`
	GeneratedAt     time.Time                 `json:"generated_at"`
	Summary         BusinessSummary           `json:"summary"`
	Transactions    TransactionAnalytics      `json:"transactions"`
	BulkOrders      BulkOrderAnalytics        `json:"bulk_orders"`
	Performance     BusinessPerformance       `json:"performance"`
	Compliance      ComplianceReport          `json:"compliance"`
	Recommendations BusinessRecommendations   `json:"recommendations"`
}

// SystemSummary provides high-level system statistics
type SystemSummary struct {
	TotalTransactions       int64     `json:"total_transactions"`
	TotalVolume             sdk.Int   `json:"total_volume"`
	UniqueUsers             int64     `json:"unique_users"`
	AverageTransactionValue sdk.Dec   `json:"average_transaction_value"`
	GrowthRate              float64   `json:"growth_rate"`
	TopRegions              []string  `json:"top_regions"`
	SuccessRate             float64   `json:"success_rate"`
	AverageProcessingTime   time.Duration `json:"average_processing_time"`
}

// SystemMetrics contains detailed system metrics
type SystemMetrics struct {
	TransactionMetrics TransactionMetrics `json:"transaction_metrics"`
	UserMetrics        UserMetrics        `json:"user_metrics"`
	VolumeMetrics      VolumeMetrics      `json:"volume_metrics"`
	PerformanceMetrics PerformanceMetrics `json:"performance_metrics"`
	ErrorMetrics       ErrorMetrics       `json:"error_metrics"`
	SecurityMetrics    SecurityMetrics    `json:"security_metrics"`
}

// TransactionMetrics provides transaction-related metrics
type TransactionMetrics struct {
	TotalCount            int64              `json:"total_count"`
	SuccessfulCount       int64              `json:"successful_count"`
	FailedCount           int64              `json:"failed_count"`
	PendingCount          int64              `json:"pending_count"`
	AverageValue          sdk.Dec            `json:"average_value"`
	MedianValue           sdk.Dec            `json:"median_value"`
	LargestTransaction    sdk.Int            `json:"largest_transaction"`
	SmallestTransaction   sdk.Int            `json:"smallest_transaction"`
	TransactionsByType    map[string]int64   `json:"transactions_by_type"`
	TransactionsByStatus  map[string]int64   `json:"transactions_by_status"`
	HourlyDistribution    []HourlyData       `json:"hourly_distribution"`
	DailyDistribution     []DailyData        `json:"daily_distribution"`
}

// UserMetrics provides user-related metrics
type UserMetrics struct {
	TotalUsers         int64                `json:"total_users"`
	ActiveUsers        int64                `json:"active_users"`
	NewUsers           int64                `json:"new_users"`
	RetentionRate      float64              `json:"retention_rate"`
	UsersByRegion      map[string]int64     `json:"users_by_region"`
	UsersByType        map[string]int64     `json:"users_by_type"`
	AverageSessionTime time.Duration        `json:"average_session_time"`
	UserEngagement     UserEngagementData   `json:"user_engagement"`
}

// VolumeMetrics provides volume-related metrics
type VolumeMetrics struct {
	TotalVolume       sdk.Int           `json:"total_volume"`
	VolumeByRegion    map[string]sdk.Int `json:"volume_by_region"`
	VolumeByType      map[string]sdk.Int `json:"volume_by_type"`
	VolumeGrowth      float64           `json:"volume_growth"`
	PeakVolume        sdk.Int           `json:"peak_volume"`
	AverageVolume     sdk.Dec           `json:"average_volume"`
	VolumeDistribution []VolumeRangeData `json:"volume_distribution"`
}

// PerformanceMetrics provides system performance metrics
type PerformanceMetrics struct {
	TransactionThroughput ThroughputMetrics    `json:"transaction_throughput"`
	LatencyMetrics        LatencyMetrics       `json:"latency_metrics"`
	ErrorRates            ErrorRateMetrics     `json:"error_rates"`
	SystemUtilization     UtilizationMetrics   `json:"system_utilization"`
	NetworkMetrics        NetworkMetrics       `json:"network_metrics"`
	ResourceUsage         ResourceUsageMetrics `json:"resource_usage"`
}

// ThroughputMetrics provides throughput statistics
type ThroughputMetrics struct {
	TransactionsPerSecond  float64 `json:"transactions_per_second"`
	TransactionsPerMinute  float64 `json:"transactions_per_minute"`
	TransactionsPerHour    float64 `json:"transactions_per_hour"`
	PeakThroughput         float64 `json:"peak_throughput"`
	AverageThroughput      float64 `json:"average_throughput"`
	ThroughputGrowth       float64 `json:"throughput_growth"`
}

// LatencyMetrics provides latency statistics
type LatencyMetrics struct {
	AverageLatency    time.Duration `json:"average_latency"`
	MedianLatency     time.Duration `json:"median_latency"`
	P90Latency        time.Duration `json:"p90_latency"`
	P95Latency        time.Duration `json:"p95_latency"`
	P99Latency        time.Duration `json:"p99_latency"`
	MaxLatency        time.Duration `json:"max_latency"`
	MinLatency        time.Duration `json:"min_latency"`
}

// ErrorRateMetrics provides error rate statistics
type ErrorRateMetrics struct {
	OverallErrorRate    float64              `json:"overall_error_rate"`
	ErrorsByType        map[string]float64   `json:"errors_by_type"`
	ErrorsByEndpoint    map[string]float64   `json:"errors_by_endpoint"`
	CriticalErrors      int64                `json:"critical_errors"`
	RecoverableErrors   int64                `json:"recoverable_errors"`
	ErrorTrend          []ErrorTrendData     `json:"error_trend"`
}

// TrendAnalysis provides trend analysis data
type TrendAnalysis struct {
	VolumeTrend       TrendData      `json:"volume_trend"`
	UserGrowthTrend   TrendData      `json:"user_growth_trend"`
	TransactionTrend  TrendData      `json:"transaction_trend"`
	GeographicTrend   TrendData      `json:"geographic_trend"`
	SeasonalPatterns  []PatternData  `json:"seasonal_patterns"`
	WeeklyPatterns    []PatternData  `json:"weekly_patterns"`
	HourlyPatterns    []PatternData  `json:"hourly_patterns"`
}

// TrendData represents trend information
type TrendData struct {
	Direction    string      `json:"direction"` // "up", "down", "stable"
	Percentage   float64     `json:"percentage"`
	Confidence   float64     `json:"confidence"`
	DataPoints   []DataPoint `json:"data_points"`
	Correlation  float64     `json:"correlation"`
	Seasonality  bool        `json:"seasonality"`
}

// DataPoint represents a single data point in trend analysis
type DataPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
	Label     string    `json:"label,omitempty"`
}

// GeographicAnalysis provides geographic distribution analysis
type GeographicAnalysis struct {
	TransactionsByRegion map[string]TransactionMetrics `json:"transactions_by_region"`
	VolumeByRegion       map[string]sdk.Int           `json:"volume_by_region"`
	UsersByRegion        map[string]int64             `json:"users_by_region"`
	GrowthByRegion       map[string]float64           `json:"growth_by_region"`
	TopCities            []CityData                   `json:"top_cities"`
	RegionalTrends       map[string]TrendData         `json:"regional_trends"`
}

// CityData represents city-specific data
type CityData struct {
	City         string  `json:"city"`
	State        string  `json:"state"`
	Transactions int64   `json:"transactions"`
	Volume       sdk.Int `json:"volume"`
	Users        int64   `json:"users"`
	GrowthRate   float64 `json:"growth_rate"`
}

// PredictiveAnalysis provides predictive analytics
type PredictiveAnalysis struct {
	VolumeForecast       []ForecastData     `json:"volume_forecast"`
	TransactionForecast  []ForecastData     `json:"transaction_forecast"`
	UserGrowthForecast   []ForecastData     `json:"user_growth_forecast"`
	ResourceNeedsForecast []ResourceForecast `json:"resource_needs_forecast"`
	TrendPredictions     []TrendPrediction  `json:"trend_predictions"`
	RiskAssessment       RiskAssessment     `json:"risk_assessment"`
}

// ForecastData represents forecasted data
type ForecastData struct {
	Date       time.Time `json:"date"`
	Value      float64   `json:"value"`
	LowerBound float64   `json:"lower_bound"`
	UpperBound float64   `json:"upper_bound"`
	Confidence float64   `json:"confidence"`
}

// BusinessSummary provides business-specific summary
type BusinessSummary struct {
	TotalTransactions int64         `json:"total_transactions"`
	TotalVolume       sdk.Int       `json:"total_volume"`
	BulkOrdersCount   int64         `json:"bulk_orders_count"`
	AverageOrderSize  sdk.Dec       `json:"average_order_size"`
	SuccessRate       float64       `json:"success_rate"`
	CostSavings       sdk.Int       `json:"cost_savings"`
	ProcessingTime    time.Duration `json:"processing_time"`
	ComplianceScore   float64       `json:"compliance_score"`
}

// TransactionAnalytics provides detailed transaction analytics
type TransactionAnalytics struct {
	Filters            AnalyticsFilters      `json:"filters"`
	GeneratedAt        time.Time             `json:"generated_at"`
	VolumeAnalysis     VolumeAnalysis        `json:"volume_analysis"`
	AmountDistribution AmountDistribution    `json:"amount_distribution"`
	TimePatterns       TimePatterns          `json:"time_patterns"`
	GeographicPatterns GeographicPatterns    `json:"geographic_patterns"`
	UserBehavior       UserBehaviorAnalysis  `json:"user_behavior"`
	AnomalyDetection   AnomalyDetection      `json:"anomaly_detection"`
	Correlations       CorrelationAnalysis   `json:"correlations"`
}

// AnalyticsFilters represents filters for analytics queries
type AnalyticsFilters struct {
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	UserAddress   string    `json:"user_address,omitempty"`
	MinAmount     sdk.Int   `json:"min_amount,omitempty"`
	MaxAmount     sdk.Int   `json:"max_amount,omitempty"`
	Status        string    `json:"status,omitempty"`
	Region        string    `json:"region,omitempty"`
	TransactionType string  `json:"transaction_type,omitempty"`
}

// DashboardMetrics provides real-time dashboard metrics
type DashboardMetrics struct {
	UserAddress       string                `json:"user_address"`
	TimeRange         string                `json:"time_range"`
	GeneratedAt       time.Time             `json:"generated_at"`
	TransactionStats  TransactionStats      `json:"transaction_stats"`
	VolumeMetrics     UserVolumeMetrics     `json:"volume_metrics"`
	TrendData         []TrendDataPoint      `json:"trend_data"`
	TopCounterparties []CounterpartyData    `json:"top_counterparties"`
	ActivityHeatmap   []ActivityData        `json:"activity_heatmap"`
	Alerts            []AlertData           `json:"alerts"`
}

// TransactionStats provides transaction statistics for a user
type TransactionStats struct {
	TotalCount    int64   `json:"total_count"`
	SuccessCount  int64   `json:"success_count"`
	FailedCount   int64   `json:"failed_count"`
	PendingCount  int64   `json:"pending_count"`
	SuccessRate   float64 `json:"success_rate"`
	AverageAmount sdk.Dec `json:"average_amount"`
}

// RealTimeMetrics provides real-time system metrics
type RealTimeMetrics struct {
	Timestamp             time.Time `json:"timestamp"`
	ActiveUsers           int64     `json:"active_users"`
	TransactionsPerSecond float64   `json:"transactions_per_second"`
	TotalValueLocked      sdk.Int   `json:"total_value_locked"`
	NetworkHealth         float64   `json:"network_health"`
	SystemLoad            float64   `json:"system_load"`
	QueueStatus           QueueStatus `json:"queue_status"`
	ErrorRate             float64   `json:"error_rate"`
	ResponseTime          time.Duration `json:"response_time"`
	GasUsage              int64     `json:"gas_usage"`
	MemoryUsage           int64     `json:"memory_usage"`
	AlertsCount           int64     `json:"alerts_count"`
}

// QueueStatus represents the status of processing queues
type QueueStatus struct {
	PendingTransactions int64 `json:"pending_transactions"`
	ProcessingRate      float64 `json:"processing_rate"`
	AverageWaitTime     time.Duration `json:"average_wait_time"`
	QueueHealth         string `json:"queue_health"` // "healthy", "warning", "critical"
}

// AnomalyReport represents an anomaly detection report
type AnomalyReport struct {
	ReportID        string                 `json:"report_id"`
	StartDate       time.Time              `json:"start_date"`
	EndDate         time.Time              `json:"end_date"`
	GeneratedAt     time.Time              `json:"generated_at"`
	Severity        AnomalySeverity        `json:"severity"`
	Anomalies       []DetectedAnomaly      `json:"anomalies"`
	Summary         AnomalySummary         `json:"summary"`
	Impact          AnomalyImpact          `json:"impact"`
	Recommendations []AnomalyRecommendation `json:"recommendations"`
}

// DetectedAnomaly represents a detected anomaly
type DetectedAnomaly struct {
	AnomalyID     string          `json:"anomaly_id"`
	Type          string          `json:"type"`
	Severity      AnomalySeverity `json:"severity"`
	DetectedAt    time.Time       `json:"detected_at"`
	Description   string          `json:"description"`
	AffectedMetric string         `json:"affected_metric"`
	ExpectedValue float64         `json:"expected_value"`
	ActualValue   float64         `json:"actual_value"`
	Deviation     float64         `json:"deviation"`
	Confidence    float64         `json:"confidence"`
	Context       map[string]interface{} `json:"context"`
}

// ExportResult represents the result of a report export
type ExportResult struct {
	ReportID    string        `json:"report_id"`
	Format      ExportFormat  `json:"format"`
	Data        []byte        `json:"data"`
	ContentType string        `json:"content_type"`
	Filename    string        `json:"filename"`
	Size        int64         `json:"size"`
	GeneratedAt time.Time     `json:"generated_at"`
}

// ReportSummary provides a summary of available reports
type ReportSummary struct {
	ReportID    string      `json:"report_id"`
	ReportType  ReportType  `json:"report_type"`
	StartDate   time.Time   `json:"start_date"`
	EndDate     time.Time   `json:"end_date"`
	GeneratedAt time.Time   `json:"generated_at"`
	Status      string      `json:"status"`
	Size        int64       `json:"size,omitempty"`
}

// ReportSchedule represents a scheduled report
type ReportSchedule struct {
	ScheduleID   string                 `json:"schedule_id"`
	ReportType   ReportType            `json:"report_type"`
	Frequency    string                `json:"frequency"` // "daily", "weekly", "monthly"
	Recipients   []string              `json:"recipients"`
	Filters      AnalyticsFilters      `json:"filters"`
	ExportFormat ExportFormat          `json:"export_format"`
	IsActive     bool                  `json:"is_active"`
	CreatedAt    time.Time             `json:"created_at"`
	LastRun      time.Time             `json:"last_run,omitempty"`
	NextRun      time.Time             `json:"next_run,omitempty"`
	Settings     map[string]interface{} `json:"settings,omitempty"`
}

// Additional supporting types for various metrics
type HourlyData struct {
	Hour  int   `json:"hour"`
	Count int64 `json:"count"`
	Value sdk.Int `json:"value,omitempty"`
}

type DailyData struct {
	Date  time.Time `json:"date"`
	Count int64     `json:"count"`
	Value sdk.Int   `json:"value,omitempty"`
}

type UserEngagementData struct {
	DailyActiveUsers   int64   `json:"daily_active_users"`
	WeeklyActiveUsers  int64   `json:"weekly_active_users"`
	MonthlyActiveUsers int64   `json:"monthly_active_users"`
	EngagementScore    float64 `json:"engagement_score"`
}

type VolumeRangeData struct {
	Range       string  `json:"range"`
	Count       int64   `json:"count"`
	Percentage  float64 `json:"percentage"`
	TotalVolume sdk.Int `json:"total_volume"`
}

type UtilizationMetrics struct {
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
	DiskUsage   float64 `json:"disk_usage"`
	NetworkIO   int64   `json:"network_io"`
}

type NetworkMetrics struct {
	Bandwidth       int64   `json:"bandwidth"`
	Latency         time.Duration `json:"latency"`
	PacketLoss      float64 `json:"packet_loss"`
	Connections     int64   `json:"connections"`
	ThroughputMbps  float64 `json:"throughput_mbps"`
}

type ResourceUsageMetrics struct {
	TotalMemory      int64 `json:"total_memory"`
	UsedMemory       int64 `json:"used_memory"`
	TotalStorage     int64 `json:"total_storage"`
	UsedStorage      int64 `json:"used_storage"`
	ActiveConnections int64 `json:"active_connections"`
}

type ErrorTrendData struct {
	Timestamp time.Time `json:"timestamp"`
	ErrorRate float64   `json:"error_rate"`
	Count     int64     `json:"count"`
}

type PatternData struct {
	Pattern     string    `json:"pattern"`
	Occurrences int64     `json:"occurrences"`
	Strength    float64   `json:"strength"`
	Periods     []string  `json:"periods"`
}

type ResourceForecast struct {
	Date           time.Time `json:"date"`
	CPURequirement float64   `json:"cpu_requirement"`
	MemoryRequirement int64  `json:"memory_requirement"`
	StorageRequirement int64 `json:"storage_requirement"`
	Confidence     float64   `json:"confidence"`
}

type TrendPrediction struct {
	Metric      string    `json:"metric"`
	Direction   string    `json:"direction"`
	Magnitude   float64   `json:"magnitude"`
	Probability float64   `json:"probability"`
	Timeframe   string    `json:"timeframe"`
}

type RiskAssessment struct {
	OverallRiskScore float64            `json:"overall_risk_score"`
	RiskFactors      []RiskFactor       `json:"risk_factors"`
	Mitigations      []string           `json:"mitigations"`
	Recommendations  []string           `json:"recommendations"`
}

type RiskFactor struct {
	Factor      string  `json:"factor"`
	Severity    string  `json:"severity"`
	Probability float64 `json:"probability"`
	Impact      float64 `json:"impact"`
	Description string  `json:"description"`
}

// Key prefixes for analytics storage
var (
	AnalyticsReportPrefix   = []byte{0x40}
	BusinessReportPrefix    = []byte{0x41}
	ReportSchedulePrefix    = []byte{0x42}
	MetricsCachePrefix      = []byte{0x43}
	AnomalyDetectionPrefix  = []byte{0x44}
)

// Event types for analytics and reporting
const (
	EventTypeReportGenerated  = "report_generated"
	EventTypeReportScheduled  = "report_scheduled"
	EventTypeAnomalyDetected  = "anomaly_detected"
	EventTypeMetricsUpdated   = "metrics_updated"
)

// Attribute keys for analytics events
const (
	AttributeKeyReportType   = "report_type"
	AttributeKeyScheduleID   = "schedule_id"
	AttributeKeyFrequency    = "frequency"
	AttributeKeyAnomalyType  = "anomaly_type"
	AttributeKeySeverity     = "severity"
	AttributeKeyMetricType   = "metric_type"
)