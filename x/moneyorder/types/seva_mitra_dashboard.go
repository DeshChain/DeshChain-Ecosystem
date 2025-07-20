package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SevaMitraDashboardData represents comprehensive dashboard data for a Seva Mitra
type SevaMitraDashboardData struct {
	MitraInfo        SevaMitra         `json:"mitra_info"`
	Summary          DashboardSummary  `json:"summary"`
	EarningsData     EarningsData      `json:"earnings_data"`
	ServiceRequests  []ServiceRequest  `json:"service_requests"`
	PerformanceStats PerformanceStats  `json:"performance_stats"`
	Analytics        AnalyticsData     `json:"analytics"`
	Notifications    []Notification    `json:"notifications"`
	Rankings         RankingData       `json:"rankings"`
}

// DashboardSummary provides key summary metrics
type DashboardSummary struct {
	TotalServices     int64         `json:"total_services"`
	TodayServices     int64         `json:"today_services"`
	MonthlyServices   int64         `json:"monthly_services"`
	TotalEarnings     sdk.Int       `json:"total_earnings"`
	TodayEarnings     sdk.Int       `json:"today_earnings"`
	MonthlyEarnings   sdk.Int       `json:"monthly_earnings"`
	PendingRequests   int64         `json:"pending_requests"`
	TrustScore        float64       `json:"trust_score"`
	ResponseTime      time.Duration `json:"response_time"`
	CustomerRating    float64       `json:"customer_rating"`
	OnlineStatus      bool          `json:"online_status"`
	LastActiveTime    time.Time     `json:"last_active_time"`
}

// EarningsData provides detailed earnings information
type EarningsData struct {
	TotalEarnings        sdk.Int                    `json:"total_earnings"`
	WeeklyEarnings       []DailyEarnings           `json:"weekly_earnings"`
	MonthlyEarnings      []MonthlyEarnings         `json:"monthly_earnings"`
	ServiceTypeBreakdown map[string]sdk.Int        `json:"service_type_breakdown"`
	RecentTransactions   []EarningsTransaction     `json:"recent_transactions"`
	PendingPayments      []PendingPayment          `json:"pending_payments"`
	PaymentHistory       []PaymentRecord           `json:"payment_history"`
}

// DailyEarnings represents earnings for a specific day
type DailyEarnings struct {
	Date     time.Time `json:"date"`
	Amount   sdk.Int   `json:"amount"`
	Services int64     `json:"services"`
}

// MonthlyEarnings represents earnings for a specific month
type MonthlyEarnings struct {
	Month    string  `json:"month"`
	Year     int     `json:"year"`
	Amount   sdk.Int `json:"amount"`
	Services int64   `json:"services"`
	Growth   float64 `json:"growth"` // Percentage growth from previous month
}

// EarningsTransaction represents a single earnings transaction
type EarningsTransaction struct {
	ID            string    `json:"id"`
	ServiceID     string    `json:"service_id"`
	ServiceType   string    `json:"service_type"`
	Amount        sdk.Int   `json:"amount"`
	Commission    sdk.Int   `json:"commission"`
	NetAmount     sdk.Int   `json:"net_amount"`
	CustomerName  string    `json:"customer_name"`
	Timestamp     time.Time `json:"timestamp"`
	Status        string    `json:"status"`
}

// PendingPayment represents a payment that is pending
type PendingPayment struct {
	ID            string    `json:"id"`
	Amount        sdk.Int   `json:"amount"`
	ServiceType   string    `json:"service_type"`
	CustomerName  string    `json:"customer_name"`
	DueDate       time.Time `json:"due_date"`
	DaysOverdue   int       `json:"days_overdue"`
}

// PaymentRecord represents a historical payment record
type PaymentRecord struct {
	ID          string    `json:"id"`
	Amount      sdk.Int   `json:"amount"`
	PaymentDate time.Time `json:"payment_date"`
	Method      string    `json:"method"`
	Status      string    `json:"status"`
	Reference   string    `json:"reference"`
}

// PerformanceStats provides performance metrics
type PerformanceStats struct {
	CompletionRate       float64              `json:"completion_rate"`       // Percentage of completed services
	AverageResponseTime  time.Duration        `json:"average_response_time"` // Average time to respond to requests
	CustomerSatisfaction float64              `json:"customer_satisfaction"` // Average customer rating
	RepeatedCustomers    int64                `json:"repeated_customers"`    // Number of customers who used service multiple times
	DisputeRate          float64              `json:"dispute_rate"`          // Percentage of services that resulted in disputes
	OnTimeDelivery       float64              `json:"on_time_delivery"`      // Percentage of services completed on time
	MonthlyGrowth        float64              `json:"monthly_growth"`        // Growth rate compared to previous month
	ServiceReliability   float64              `json:"service_reliability"`   // Overall reliability score
	TrustScoreHistory    []TrustScorePoint    `json:"trust_score_history"`   // Historical trust score data
}

// TrustScorePoint represents a trust score at a specific time
type TrustScorePoint struct {
	Date  time.Time `json:"date"`
	Score float64   `json:"score"`
}

// AnalyticsData provides analytical insights
type AnalyticsData struct {
	ServiceTypesChart      []ServiceTypeData      `json:"service_types_chart"`
	HourlyActivity         []HourlyActivityData   `json:"hourly_activity"`
	GeographicDistribution []GeographicData       `json:"geographic_distribution"`
	CustomerDemographics   CustomerDemographics   `json:"customer_demographics"`
	SeasonalTrends         []SeasonalTrendData    `json:"seasonal_trends"`
	CompetitorAnalysis     CompetitorAnalysis     `json:"competitor_analysis"`
	MarketOpportunities    []MarketOpportunity    `json:"market_opportunities"`
}

// ServiceTypeData represents service type distribution
type ServiceTypeData struct {
	ServiceType string  `json:"service_type"`
	Count       int64   `json:"count"`
	Percentage  float64 `json:"percentage"`
	Revenue     sdk.Int `json:"revenue"`
}

// HourlyActivityData represents activity patterns by hour
type HourlyActivityData struct {
	Hour     int   `json:"hour"`
	Requests int64 `json:"requests"`
	Revenue  sdk.Int `json:"revenue"`
}

// GeographicData represents geographic distribution of services
type GeographicData struct {
	Location string  `json:"location"`
	Pincode  string  `json:"pincode"`
	Count    int64   `json:"count"`
	Revenue  sdk.Int `json:"revenue"`
}

// CustomerDemographics provides customer demographic information
type CustomerDemographics struct {
	AgeGroups       []AgeGroupData       `json:"age_groups"`
	GenderSplit     GenderSplitData      `json:"gender_split"`
	IncomeRanges    []IncomeRangeData    `json:"income_ranges"`
	NewVsReturning  NewVsReturningData   `json:"new_vs_returning"`
}

// AgeGroupData represents age group distribution
type AgeGroupData struct {
	AgeRange   string  `json:"age_range"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
}

// GenderSplitData represents gender distribution
type GenderSplitData struct {
	Male        int64   `json:"male"`
	Female      int64   `json:"female"`
	Other       int64   `json:"other"`
	MalePercent float64 `json:"male_percent"`
	FemalePercent float64 `json:"female_percent"`
}

// IncomeRangeData represents income range distribution
type IncomeRangeData struct {
	Range      string  `json:"range"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
}

// NewVsReturningData represents new vs returning customer data
type NewVsReturningData struct {
	NewCustomers       int64   `json:"new_customers"`
	ReturningCustomers int64   `json:"returning_customers"`
	NewPercent         float64 `json:"new_percent"`
	ReturningPercent   float64 `json:"returning_percent"`
}

// SeasonalTrendData represents seasonal trend information
type SeasonalTrendData struct {
	Month    string  `json:"month"`
	Services int64   `json:"services"`
	Revenue  sdk.Int `json:"revenue"`
	Growth   float64 `json:"growth"`
}

// CompetitorAnalysis provides competitor analysis data
type CompetitorAnalysis struct {
	MarketPosition    int     `json:"market_position"`
	TrustScoreRank    int     `json:"trust_score_rank"`
	PricingComparison float64 `json:"pricing_comparison"` // How pricing compares to market average
	ServiceQuality    float64 `json:"service_quality"`    // Relative service quality score
	MarketShare       float64 `json:"market_share"`       // Estimated market share percentage
}

// MarketOpportunity represents a market opportunity
type MarketOpportunity struct {
	OpportunityType string  `json:"opportunity_type"`
	Description     string  `json:"description"`
	PotentialRevenue sdk.Int `json:"potential_revenue"`
	Difficulty      string  `json:"difficulty"` // "Low", "Medium", "High"
	Timeline        string  `json:"timeline"`   // Expected timeline to realize
}

// ServiceRequest represents a service request in the dashboard
type ServiceRequest struct {
	ID                      string                  `json:"id"`
	CustomerAddress         string                  `json:"customer_address"`
	CustomerName            string                  `json:"customer_name"`
	SevaMitraAddress        string                  `json:"seva_mitra_address"`
	ServiceType             string                  `json:"service_type"`
	Amount                  sdk.Int                 `json:"amount"`
	Status                  ServiceRequestStatus    `json:"status"`
	Priority                ServicePriority         `json:"priority"`
	Location                ServiceLocation         `json:"location"`
	CreatedAt               time.Time               `json:"created_at"`
	AcceptedAt              time.Time               `json:"accepted_at,omitempty"`
	CompletedAt             time.Time               `json:"completed_at,omitempty"`
	EstimatedCompletionTime time.Time               `json:"estimated_completion_time,omitempty"`
	ServiceDuration         time.Duration           `json:"service_duration,omitempty"`
	EarningsAmount          sdk.Int                 `json:"earnings_amount"`
	CustomerRating          float64                 `json:"customer_rating,omitempty"`
	CustomerFeedback        string                  `json:"customer_feedback,omitempty"`
	CompletionNote          string                  `json:"completion_note,omitempty"`
}

// ServiceRequestStatus represents the status of a service request
type ServiceRequestStatus int32

const (
	ServiceRequestStatus_PENDING     ServiceRequestStatus = 0
	ServiceRequestStatus_ACCEPTED    ServiceRequestStatus = 1
	ServiceRequestStatus_IN_PROGRESS ServiceRequestStatus = 2
	ServiceRequestStatus_COMPLETED   ServiceRequestStatus = 3
	ServiceRequestStatus_CANCELLED   ServiceRequestStatus = 4
	ServiceRequestStatus_DISPUTED    ServiceRequestStatus = 5
)

// String returns the string representation of ServiceRequestStatus
func (srs ServiceRequestStatus) String() string {
	switch srs {
	case ServiceRequestStatus_PENDING:
		return "PENDING"
	case ServiceRequestStatus_ACCEPTED:
		return "ACCEPTED"
	case ServiceRequestStatus_IN_PROGRESS:
		return "IN_PROGRESS"
	case ServiceRequestStatus_COMPLETED:
		return "COMPLETED"
	case ServiceRequestStatus_CANCELLED:
		return "CANCELLED"
	case ServiceRequestStatus_DISPUTED:
		return "DISPUTED"
	default:
		return "UNKNOWN"
	}
}

// ServicePriority represents the priority level of a service request
type ServicePriority int32

const (
	ServicePriority_LOW      ServicePriority = 0
	ServicePriority_NORMAL   ServicePriority = 1
	ServicePriority_HIGH     ServicePriority = 2
	ServicePriority_URGENT   ServicePriority = 3
)

// String returns the string representation of ServicePriority
func (sp ServicePriority) String() string {
	switch sp {
	case ServicePriority_LOW:
		return "LOW"
	case ServicePriority_NORMAL:
		return "NORMAL"
	case ServicePriority_HIGH:
		return "HIGH"
	case ServicePriority_URGENT:
		return "URGENT"
	default:
		return "UNKNOWN"
	}
}

// ServiceLocation represents the location information for a service
type ServiceLocation struct {
	Address     string  `json:"address"`
	City        string  `json:"city"`
	State       string  `json:"state"`
	Pincode     string  `json:"pincode"`
	Latitude    float64 `json:"latitude,omitempty"`
	Longitude   float64 `json:"longitude,omitempty"`
	Landmark    string  `json:"landmark,omitempty"`
}

// Notification represents a notification for the Seva Mitra
type Notification struct {
	ID        string                    `json:"id"`
	Type      NotificationType          `json:"type"`
	Title     string                    `json:"title"`
	Message   string                    `json:"message"`
	Priority  NotificationPriority      `json:"priority"`
	CreatedAt time.Time                 `json:"created_at"`
	ReadAt    time.Time                 `json:"read_at,omitempty"`
	Data      map[string]interface{}    `json:"data,omitempty"`
}

// NotificationType represents different types of notifications
type NotificationType int32

const (
	NotificationType_URGENT_REQUEST  NotificationType = 0
	NotificationType_RATING_ALERT    NotificationType = 1
	NotificationType_PAYMENT         NotificationType = 2
	NotificationType_ACHIEVEMENT     NotificationType = 3
	NotificationType_SYSTEM          NotificationType = 4
	NotificationType_PROMOTION       NotificationType = 5
)

// NotificationPriority represents notification priority levels
type NotificationPriority int32

const (
	NotificationPriority_LOW    NotificationPriority = 0
	NotificationPriority_MEDIUM NotificationPriority = 1
	NotificationPriority_HIGH   NotificationPriority = 2
)

// RankingData provides ranking information
type RankingData struct {
	LocalRank      int   `json:"local_rank"`       // Rank within local area (same pincode/city)
	RegionalRank   int   `json:"regional_rank"`    // Rank within region (same state)
	NationalRank   int   `json:"national_rank"`    // Rank nationally
	CategoryRank   int   `json:"category_rank"`    // Rank within service category
	TrustScoreRank int   `json:"trust_score_rank"` // Rank based on trust score
	EarningsRank   int   `json:"earnings_rank"`    // Rank based on earnings
	TotalMitras    int64 `json:"total_mitras"`     // Total number of Seva Mitras for context
}

// SevaMitra represents a Seva Mitra profile
type SevaMitra struct {
	Address               string             `json:"address"`
	Name                  string             `json:"name"`
	PhoneNumber           string             `json:"phone_number"`
	Email                 string             `json:"email,omitempty"`
	Services              []string           `json:"services"`
	Location              ServiceLocation    `json:"location"`
	TrustScore            float64            `json:"trust_score"`
	CompletedServices     int64              `json:"completed_services"`
	TotalEarnings         sdk.Int            `json:"total_earnings"`
	JoinedAt              time.Time          `json:"joined_at"`
	IsActive              bool               `json:"is_active"`
	IsOnline              bool               `json:"is_online"`
	OnlineSince           time.Time          `json:"online_since,omitempty"`
	LastActiveTime        time.Time          `json:"last_active_time"`
	KYCLevel              string             `json:"kyc_level"`
	BiometricEnabled      bool               `json:"biometric_enabled"`
	ServiceAvailability   map[string]bool    `json:"service_availability"`
	WorkingHours          WorkingHours       `json:"working_hours"`
	CommissionRates       map[string]float64 `json:"commission_rates"`
	MaxServiceAmount      sdk.Int            `json:"max_service_amount"`
	AverageResponseTime   time.Duration      `json:"average_response_time"`
	CustomerRating        float64            `json:"customer_rating"`
	Languages             []string           `json:"languages"`
}

// WorkingHours represents working hours for a Seva Mitra
type WorkingHours struct {
	Monday    DaySchedule `json:"monday"`
	Tuesday   DaySchedule `json:"tuesday"`
	Wednesday DaySchedule `json:"wednesday"`
	Thursday  DaySchedule `json:"thursday"`
	Friday    DaySchedule `json:"friday"`
	Saturday  DaySchedule `json:"saturday"`
	Sunday    DaySchedule `json:"sunday"`
}

// DaySchedule represents schedule for a single day
type DaySchedule struct {
	IsWorking bool   `json:"is_working"`
	StartTime string `json:"start_time"` // Format: "09:00"
	EndTime   string `json:"end_time"`   // Format: "17:00"
	BreakTime string `json:"break_time,omitempty"` // Format: "12:00-13:00"
}

// Key prefixes for dashboard data storage
var (
	ServiceRequestPrefix       = []byte{0x20}
	EarningsTransactionPrefix  = []byte{0x21}
	NotificationPrefix         = []byte{0x22}
	PerformanceStatsPrefix     = []byte{0x23}
	SevaMitraPrefix           = []byte{0x24}
)

// Event types for Seva Mitra dashboard
const (
	EventTypeSevaMitraStatusUpdate    = "seva_mitra_status_update"
	EventTypeSevaMitraServiceUpdate   = "seva_mitra_service_update"
	EventTypeServiceRequestAccepted   = "service_request_accepted"
	EventTypeServiceRequestCompleted  = "service_request_completed"
	EventTypeEarningsUpdated          = "earnings_updated"
)

// Attribute keys for Seva Mitra events
const (
	AttributeKeySevaMitraAddress = "seva_mitra_address"
	AttributeKeyServiceType      = "service_type"
	AttributeKeyRequestID        = "request_id"
	AttributeKeyEarnings         = "earnings"
	AttributeKeyOnlineStatus     = "online_status"
	AttributeKeyAvailable        = "available"
	AttributeKeyTimestamp        = "timestamp"
)