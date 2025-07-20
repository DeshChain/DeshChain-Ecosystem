package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BulkOrderStatus represents the status of a bulk order
type BulkOrderStatus int32

const (
	BulkOrderStatus_PENDING    BulkOrderStatus = 0
	BulkOrderStatus_PROCESSING BulkOrderStatus = 1
	BulkOrderStatus_COMPLETED  BulkOrderStatus = 2
	BulkOrderStatus_CANCELLED  BulkOrderStatus = 3
	BulkOrderStatus_FAILED     BulkOrderStatus = 4
)

// String returns the string representation of BulkOrderStatus
func (bos BulkOrderStatus) String() string {
	switch bos {
	case BulkOrderStatus_PENDING:
		return "PENDING"
	case BulkOrderStatus_PROCESSING:
		return "PROCESSING"
	case BulkOrderStatus_COMPLETED:
		return "COMPLETED"
	case BulkOrderStatus_CANCELLED:
		return "CANCELLED"
	case BulkOrderStatus_FAILED:
		return "FAILED"
	default:
		return "UNKNOWN"
	}
}

// BulkOrderPriority represents the priority level of bulk orders
type BulkOrderPriority int32

const (
	BulkOrderPriority_LOW    BulkOrderPriority = 0
	BulkOrderPriority_NORMAL BulkOrderPriority = 1
	BulkOrderPriority_HIGH   BulkOrderPriority = 2
	BulkOrderPriority_URGENT BulkOrderPriority = 3
)

// String returns the string representation of BulkOrderPriority
func (bop BulkOrderPriority) String() string {
	switch bop {
	case BulkOrderPriority_LOW:
		return "LOW"
	case BulkOrderPriority_NORMAL:
		return "NORMAL"
	case BulkOrderPriority_HIGH:
		return "HIGH"
	case BulkOrderPriority_URGENT:
		return "URGENT"
	default:
		return "UNKNOWN"
	}
}

// BulkOrder represents a collection of money orders processed together
type BulkOrder struct {
	ID                 string                     `json:"id"`
	BusinessAddress    string                     `json:"business_address"`
	TotalOrders        int64                      `json:"total_orders"`
	TotalAmount        sdk.Int                    `json:"total_amount"`
	Status             BulkOrderStatus            `json:"status"`
	CreatedAt          time.Time                  `json:"created_at"`
	CompletedAt        time.Time                  `json:"completed_at,omitempty"`
	CancelledAt        time.Time                  `json:"cancelled_at,omitempty"`
	CancellationReason string                     `json:"cancellation_reason,omitempty"`
	Metadata           BulkOrderMetadata          `json:"metadata"`
	ProcessingStats    BulkOrderProcessingStats   `json:"processing_stats"`
	FeesCharged        sdk.Int                    `json:"fees_charged"`
	Discount           sdk.Int                    `json:"discount,omitempty"`
}

// BulkOrderMetadata contains additional information about the bulk order
type BulkOrderMetadata struct {
	Description    string            `json:"description,omitempty"`
	Reference      string            `json:"reference,omitempty"`
	CustomerRef    string            `json:"customer_ref,omitempty"`
	Department     string            `json:"department,omitempty"`
	ProjectCode    string            `json:"project_code,omitempty"`
	Tags           []string          `json:"tags,omitempty"`
	CustomFields   map[string]string `json:"custom_fields,omitempty"`
	NotifyEmail    string            `json:"notify_email,omitempty"`
	WebhookURL     string            `json:"webhook_url,omitempty"`
	ScheduledTime  time.Time         `json:"scheduled_time,omitempty"`
	ExpiryTime     time.Time         `json:"expiry_time,omitempty"`
}

// BulkOrderProcessingStats tracks processing statistics
type BulkOrderProcessingStats struct {
	TotalOrders        int64         `json:"total_orders"`
	ProcessedOrders    int64         `json:"processed_orders"`
	SuccessfulOrders   int64         `json:"successful_orders"`
	FailedOrders       int64         `json:"failed_orders"`
	StartTime          time.Time     `json:"start_time"`
	EndTime            time.Time     `json:"end_time,omitempty"`
	ProcessingDuration time.Duration `json:"processing_duration,omitempty"`
	AverageOrderTime   time.Duration `json:"average_order_time,omitempty"`
	ErrorRate          float64       `json:"error_rate"`
}

// BulkOrderItem represents a single order within a bulk order
type BulkOrderItem struct {
	BulkOrderID      string              `json:"bulk_order_id"`
	OrderID          string              `json:"order_id"`
	Index            int                 `json:"index"`
	RecipientAddress string              `json:"recipient_address"`
	Amount           sdk.Int             `json:"amount"`
	Status           MoneyOrderStatus    `json:"status"`
	CreatedAt        time.Time           `json:"created_at"`
	ProcessedAt      time.Time           `json:"processed_at,omitempty"`
	CompletedAt      time.Time           `json:"completed_at,omitempty"`
	FailedAt         time.Time           `json:"failed_at,omitempty"`
	ErrorMessage     string              `json:"error_message,omitempty"`
	Memo             string              `json:"memo,omitempty"`
	Priority         BulkOrderPriority   `json:"priority"`
	RetryCount       int                 `json:"retry_count"`
	LastRetryAt      time.Time           `json:"last_retry_at,omitempty"`
}

// ValidatedBulkOrderItem represents a validated bulk order item
type ValidatedBulkOrderItem struct {
	Index            int               `json:"index"`
	OriginalOrder    BulkOrderItem     `json:"original_order"`
	Amount           sdk.Int           `json:"amount"`
	RecipientAddress string            `json:"recipient_address"`
	SenderAddress    string            `json:"sender_address"`
	Memo             string            `json:"memo"`
	Priority         BulkOrderPriority `json:"priority"`
}

// BulkOrderResult represents the result of processing a bulk order
type BulkOrderResult struct {
	BulkOrderID      string         `json:"bulk_order_id"`
	TotalOrders      int64          `json:"total_orders"`
	SuccessfulOrders []OrderResult  `json:"successful_orders"`
	FailedOrders     []OrderFailure `json:"failed_orders"`
	ProcessingTime   time.Duration  `json:"processing_time"`
	TotalFees        sdk.Int        `json:"total_fees"`
	NetAmount        sdk.Int        `json:"net_amount"`
}

// OrderResult represents a successfully processed order
type OrderResult struct {
	Index            int               `json:"index"`
	OrderID          string            `json:"order_id"`
	RecipientAddress string            `json:"recipient_address"`
	Amount           sdk.Int           `json:"amount"`
	Status           MoneyOrderStatus  `json:"status"`
	CreatedAt        time.Time         `json:"created_at"`
	ProcessingTime   time.Duration     `json:"processing_time"`
	TransactionHash  string            `json:"transaction_hash,omitempty"`
}

// OrderFailure represents a failed order
type OrderFailure struct {
	Index     int     `json:"index"`
	OrderID   string  `json:"order_id,omitempty"`
	Error     string  `json:"error"`
	Amount    sdk.Int `json:"amount"`
	ErrorCode string  `json:"error_code,omitempty"`
	Retryable bool    `json:"retryable"`
}

// BatchProcessingResult represents the result of processing a batch
type BatchProcessingResult struct {
	SuccessfulOrders []OrderResult  `json:"successful_orders"`
	FailedOrders     []OrderFailure `json:"failed_orders"`
}

// BulkOrderTemplate represents a template for creating bulk orders
type BulkOrderTemplate struct {
	Name        string                `json:"name"`
	Description string                `json:"description"`
	BusinessID  string                `json:"business_id"`
	Orders      []BulkOrderItem       `json:"orders"`
	Metadata    BulkOrderMetadata     `json:"metadata"`
	Settings    BulkOrderSettings     `json:"settings"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
	IsActive    bool                  `json:"is_active"`
}

// BulkOrderSettings contains settings for bulk order processing
type BulkOrderSettings struct {
	BatchSize           int           `json:"batch_size"`
	MaxRetries          int           `json:"max_retries"`
	RetryDelay          time.Duration `json:"retry_delay"`
	StopOnFirstFailure  bool          `json:"stop_on_first_failure"`
	ValidateRecipients  bool          `json:"validate_recipients"`
	AllowDuplicates     bool          `json:"allow_duplicates"`
	RequireConfirmation bool          `json:"require_confirmation"`
	NotifyOnCompletion  bool          `json:"notify_on_completion"`
	NotifyOnFailure     bool          `json:"notify_on_failure"`
}

// BulkOrderValidationResult represents the result of validating a bulk order
type BulkOrderValidationResult struct {
	IsValid       bool        `json:"is_valid"`
	TotalOrders   int64       `json:"total_orders"`
	ValidOrders   int64       `json:"valid_orders"`
	InvalidOrders int64       `json:"invalid_orders"`
	TotalAmount   sdk.Int     `json:"total_amount"`
	Warnings      []string    `json:"warnings"`
	Errors        []string    `json:"errors"`
	Duplicates    []int       `json:"duplicates,omitempty"`
	Summary       ValidationSummary `json:"summary"`
}

// ValidationSummary provides a summary of validation results
type ValidationSummary struct {
	AddressValidation  ValidationCategoryResult `json:"address_validation"`
	AmountValidation   ValidationCategoryResult `json:"amount_validation"`
	FormatValidation   ValidationCategoryResult `json:"format_validation"`
	BusinessRules      ValidationCategoryResult `json:"business_rules"`
	LimitChecks        ValidationCategoryResult `json:"limit_checks"`
}

// ValidationCategoryResult represents validation results for a category
type ValidationCategoryResult struct {
	Passed  int64    `json:"passed"`
	Failed  int64    `json:"failed"`
	Warnings int64   `json:"warnings"`
	Details []string `json:"details,omitempty"`
}

// BulkOrderStatusResponse provides detailed status information
type BulkOrderStatusResponse struct {
	BulkOrder BulkOrder           `json:"bulk_order"`
	Orders    []BulkOrderItem     `json:"orders"`
	Summary   BulkOrderSummary    `json:"summary"`
	Timeline  []StatusUpdateEvent `json:"timeline"`
	Metrics   BulkOrderMetrics    `json:"metrics"`
}

// BulkOrderSummary provides a summary of the bulk order
type BulkOrderSummary struct {
	CompletionRate    float64   `json:"completion_rate"`
	SuccessRate       float64   `json:"success_rate"`
	AverageAmount     sdk.Int   `json:"average_amount"`
	LargestAmount     sdk.Int   `json:"largest_amount"`
	SmallestAmount    sdk.Int   `json:"smallest_amount"`
	ProcessingSpeed   float64   `json:"processing_speed"` // Orders per minute
	EstimatedCompletion time.Time `json:"estimated_completion,omitempty"`
	StatusBreakdown   map[string]int64 `json:"status_breakdown"`
}

// StatusUpdateEvent represents a status update in the timeline
type StatusUpdateEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	Details   string    `json:"details,omitempty"`
	OrdersAffected int64 `json:"orders_affected,omitempty"`
}

// BulkOrderMetrics provides detailed metrics
type BulkOrderMetrics struct {
	ThroughputStats    ThroughputStats    `json:"throughput_stats"`
	ErrorAnalysis      ErrorAnalysis      `json:"error_analysis"`
	PerformanceMetrics PerformanceMetrics `json:"performance_metrics"`
	CostAnalysis       CostAnalysis       `json:"cost_analysis"`
}

// ThroughputStats provides throughput statistics
type ThroughputStats struct {
	OrdersPerSecond    float64 `json:"orders_per_second"`
	OrdersPerMinute    float64 `json:"orders_per_minute"`
	PeakThroughput     float64 `json:"peak_throughput"`
	AverageThroughput  float64 `json:"average_throughput"`
	ProcessingEfficiency float64 `json:"processing_efficiency"`
}

// ErrorAnalysis provides error analysis
type ErrorAnalysis struct {
	TotalErrors        int64                    `json:"total_errors"`
	ErrorRate          float64                  `json:"error_rate"`
	ErrorsByType       map[string]int64         `json:"errors_by_type"`
	ErrorsByCategory   map[string]int64         `json:"errors_by_category"`
	RecoverableErrors  int64                    `json:"recoverable_errors"`
	FatalErrors        int64                    `json:"fatal_errors"`
	CommonErrors       []ErrorFrequency         `json:"common_errors"`
}

// ErrorFrequency represents error frequency data
type ErrorFrequency struct {
	Error       string `json:"error"`
	Count       int64  `json:"count"`
	Percentage  float64 `json:"percentage"`
	FirstSeen   time.Time `json:"first_seen"`
	LastSeen    time.Time `json:"last_seen"`
}

// PerformanceMetrics provides performance metrics
type PerformanceMetrics struct {
	AverageProcessingTime time.Duration `json:"average_processing_time"`
	MedianProcessingTime  time.Duration `json:"median_processing_time"`
	P95ProcessingTime     time.Duration `json:"p95_processing_time"`
	P99ProcessingTime     time.Duration `json:"p99_processing_time"`
	TotalProcessingTime   time.Duration `json:"total_processing_time"`
	MemoryUsage           int64         `json:"memory_usage"`
	CPUUsage              float64       `json:"cpu_usage"`
}

// CostAnalysis provides cost analysis
type CostAnalysis struct {
	TotalFees          sdk.Int `json:"total_fees"`
	AverageFeePerOrder sdk.Int `json:"average_fee_per_order"`
	DiscountApplied    sdk.Int `json:"discount_applied"`
	NetCost            sdk.Int `json:"net_cost"`
	CostSavings        sdk.Int `json:"cost_savings"`
	EfficiencyGain     float64 `json:"efficiency_gain"`
}

// BusinessAccount represents a business account with bulk order capabilities
type BusinessAccount struct {
	Address                string                    `json:"address"`
	BusinessName           string                    `json:"business_name"`
	BusinessType           string                    `json:"business_type"`
	RegistrationNumber     string                    `json:"registration_number"`
	ContactEmail           string                    `json:"contact_email"`
	ContactPhone           string                    `json:"contact_phone"`
	IsActive               bool                      `json:"is_active"`
	BulkOrdersEnabled      bool                      `json:"bulk_orders_enabled"`
	DailyLimit             sdk.Int                   `json:"daily_limit"`
	MonthlyLimit           sdk.Int                   `json:"monthly_limit"`
	MaxBulkOrderSize       int64                     `json:"max_bulk_order_size"`
	PreferredBatchSize     int64                     `json:"preferred_batch_size"`
	VerificationLevel      string                    `json:"verification_level"`
	KYCStatus              string                    `json:"kyc_status"`
	ComplianceStatus       string                    `json:"compliance_status"`
	CreatedAt              time.Time                 `json:"created_at"`
	LastBulkOrderAt        time.Time                 `json:"last_bulk_order_at,omitempty"`
	TotalBulkOrders        int64                     `json:"total_bulk_orders"`
	TotalOrdersProcessed   int64                     `json:"total_orders_processed"`
	TotalAmountProcessed   sdk.Int                   `json:"total_amount_processed"`
	SuccessRate            float64                   `json:"success_rate"`
	AverageProcessingTime  time.Duration             `json:"average_processing_time"`
	PreferredSettings      BulkOrderSettings         `json:"preferred_settings"`
	ApiCredentials         ApiCredentials            `json:"api_credentials,omitempty"`
	WebhookConfig          WebhookConfig             `json:"webhook_config,omitempty"`
	ComplianceSettings     ComplianceSettings        `json:"compliance_settings"`
}

// ApiCredentials for programmatic access
type ApiCredentials struct {
	ApiKey        string    `json:"api_key"`
	SecretKey     string    `json:"secret_key,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	ExpiresAt     time.Time `json:"expires_at,omitempty"`
	LastUsedAt    time.Time `json:"last_used_at,omitempty"`
	IsActive      bool      `json:"is_active"`
	Permissions   []string  `json:"permissions"`
	RateLimit     int64     `json:"rate_limit"`
	UsageCount    int64     `json:"usage_count"`
}

// WebhookConfig for notifications
type WebhookConfig struct {
	URL             string            `json:"url"`
	Secret          string            `json:"secret,omitempty"`
	Events          []string          `json:"events"`
	Headers         map[string]string `json:"headers,omitempty"`
	IsActive        bool              `json:"is_active"`
	RetryCount      int               `json:"retry_count"`
	TimeoutSeconds  int               `json:"timeout_seconds"`
	LastTriggered   time.Time         `json:"last_triggered,omitempty"`
	SuccessCount    int64             `json:"success_count"`
	FailureCount    int64             `json:"failure_count"`
}

// ComplianceSettings for regulatory compliance
type ComplianceSettings struct {
	RequireRecipientKYC    bool              `json:"require_recipient_kyc"`
	RequireSenderKYC       bool              `json:"require_sender_kyc"`
	MaxAmountWithoutKYC    sdk.Int           `json:"max_amount_without_kyc"`
	BlockedCountries       []string          `json:"blocked_countries"`
	RequiredDocuments      []string          `json:"required_documents"`
	MonitoringLevel        string            `json:"monitoring_level"`
	ReportingRequirements  []string          `json:"reporting_requirements"`
	AuditTrailRequired     bool              `json:"audit_trail_required"`
	DataRetentionPeriod    time.Duration     `json:"data_retention_period"`
	EncryptionRequired     bool              `json:"encryption_required"`
	AutoFlagThresholds     map[string]sdk.Int `json:"auto_flag_thresholds"`
}

// Key prefixes for bulk order storage
var (
	BulkOrderPrefix        = []byte{0x30}
	BulkOrderItemPrefix    = []byte{0x31}
	BulkOrderTemplatePrefix = []byte{0x32}
	BusinessAccountPrefix  = []byte{0x33}
	BulkOrderStatsPrefix   = []byte{0x34}
)

// Event types for bulk orders
const (
	EventTypeBulkOrderCreated    = "bulk_order_created"
	EventTypeBulkOrderProgress   = "bulk_order_progress"
	EventTypeBulkOrderCompleted  = "bulk_order_completed"
	EventTypeBulkOrderCancelled  = "bulk_order_cancelled"
	EventTypeBulkOrderFailed     = "bulk_order_failed"
	EventTypeBusinessRegistered  = "business_registered"
)

// Attribute keys for bulk order events
const (
	AttributeKeyBulkOrderID      = "bulk_order_id"
	AttributeKeyBusinessAddress  = "business_address"
	AttributeKeyTotalOrders      = "total_orders"
	AttributeKeyProcessedOrders  = "processed_orders"
	AttributeKeySuccessfulOrders = "successful_orders"
	AttributeKeyFailedOrders     = "failed_orders"
	AttributeKeyProcessingTime   = "processing_time"
	AttributeKeyReason           = "reason"
)