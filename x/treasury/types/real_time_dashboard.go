/*
Copyright 2024 DeshChain Foundation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

// RealTimeDashboard represents the real-time transparency dashboard
type RealTimeDashboard struct {
	LastUpdated       time.Time            `json:"last_updated"`
	UpdateFrequency   time.Duration        `json:"update_frequency"`
	CommunityFund     CommunityFundMetrics `json:"community_fund"`
	DevelopmentFund   DevelopmentFundMetrics `json:"development_fund"`
	OverallMetrics    OverallFundMetrics   `json:"overall_metrics"`
	TransparencyScore uint8                `json:"transparency_score"`
	ComplianceScore   uint8                `json:"compliance_score"`
	Status            DashboardStatus      `json:"status"`
	Alerts            []Alert              `json:"alerts"`
	Notifications     []Notification       `json:"notifications"`
}

// DashboardStatus represents the status of the dashboard
type DashboardStatus string

const (
	DashboardStatusOnline      DashboardStatus = "online"
	DashboardStatusOffline     DashboardStatus = "offline"
	DashboardStatusMaintenance DashboardStatus = "maintenance"
	DashboardStatusError       DashboardStatus = "error"
	DashboardStatusUpdating    DashboardStatus = "updating"
)

// CommunityFundMetrics represents real-time community fund metrics
type CommunityFundMetrics struct {
	TotalAllocation       sdk.Coin             `json:"total_allocation"`       // 214,294,149 NAMO
	CurrentBalance        sdk.Coin             `json:"current_balance"`
	AllocatedAmount       sdk.Coin             `json:"allocated_amount"`
	SpentAmount           sdk.Coin             `json:"spent_amount"`
	RemainingAmount       sdk.Coin             `json:"remaining_amount"`
	PendingAmount         sdk.Coin             `json:"pending_amount"`
	ReservedAmount        sdk.Coin             `json:"reserved_amount"`
	ActiveProposals       uint64               `json:"active_proposals"`
	CompletedProposals    uint64               `json:"completed_proposals"`
	RejectedProposals     uint64               `json:"rejected_proposals"`
	TotalProposals        uint64               `json:"total_proposals"`
	CategoryBreakdown     []CategoryMetrics    `json:"category_breakdown"`
	TopRecipients         []RecipientMetrics   `json:"top_recipients"`
	RecentTransactions    []TransactionMetrics `json:"recent_transactions"`
	MonthlySpending       []MonthlyMetrics     `json:"monthly_spending"`
	QuarterlyReport       QuarterlyReport      `json:"quarterly_report"`
	AnnualReport          AnnualReport         `json:"annual_report"`
	TransparencyMetrics   TransparencyMetrics  `json:"transparency_metrics"`
	GovernanceMetrics     GovernanceMetrics    `json:"governance_metrics"`
	ImpactMetrics         ImpactMetrics        `json:"impact_metrics"`
	PerformanceMetrics    PerformanceMetrics   `json:"performance_metrics"`
	PredictiveAnalytics   PredictiveAnalytics  `json:"predictive_analytics"`
	LastUpdated           time.Time            `json:"last_updated"`
}

// DevelopmentFundMetrics represents real-time development fund metrics
type DevelopmentFundMetrics struct {
	TotalAllocation       sdk.Coin                `json:"total_allocation"`       // 214,294,149 NAMO
	CurrentBalance        sdk.Coin                `json:"current_balance"`
	AllocatedAmount       sdk.Coin                `json:"allocated_amount"`
	SpentAmount           sdk.Coin                `json:"spent_amount"`
	RemainingAmount       sdk.Coin                `json:"remaining_amount"`
	PendingAmount         sdk.Coin                `json:"pending_amount"`
	EmergencyReserve      sdk.Coin                `json:"emergency_reserve"`
	ActiveProjects        uint64                  `json:"active_projects"`
	CompletedProjects     uint64                  `json:"completed_projects"`
	CancelledProjects     uint64                  `json:"cancelled_projects"`
	TotalProjects         uint64                  `json:"total_projects"`
	CategoryBreakdown     []DevCategoryMetrics    `json:"category_breakdown"`
	TopDevelopers         []DeveloperMetrics      `json:"top_developers"`
	RecentTransactions    []DevTransactionMetrics `json:"recent_transactions"`
	MonthlySpending       []MonthlyMetrics        `json:"monthly_spending"`
	QuarterlyReport       QuarterlyReport         `json:"quarterly_report"`
	AnnualReport          AnnualReport            `json:"annual_report"`
	QualityMetrics        QualityMetrics          `json:"quality_metrics"`
	SecurityMetrics       SecurityMetrics         `json:"security_metrics"`
	PerformanceMetrics    PerformanceMetrics      `json:"performance_metrics"`
	InnovationMetrics     InnovationMetrics       `json:"innovation_metrics"`
	DeliveryMetrics       DeliveryMetrics         `json:"delivery_metrics"`
	ROIMetrics            ROIMetrics              `json:"roi_metrics"`
	LastUpdated           time.Time               `json:"last_updated"`
}

// OverallFundMetrics represents overall fund metrics
type OverallFundMetrics struct {
	TotalFunds            sdk.Coin              `json:"total_funds"`            // 428,588,298 NAMO (both funds)
	TotalAllocated        sdk.Coin              `json:"total_allocated"`
	TotalSpent            sdk.Coin              `json:"total_spent"`
	TotalRemaining        sdk.Coin              `json:"total_remaining"`
	AllocationEfficiency  math.LegacyDec        `json:"allocation_efficiency"`
	SpendingEfficiency    math.LegacyDec        `json:"spending_efficiency"`
	OverallTransparency   uint8                 `json:"overall_transparency"`
	OverallCompliance     uint8                 `json:"overall_compliance"`
	CommunityEngagement   uint8                 `json:"community_engagement"`
	StakeholderSatisfaction uint8               `json:"stakeholder_satisfaction"`
	SustainabilityScore   uint8                 `json:"sustainability_score"`
	InnovationIndex       uint8                 `json:"innovation_index"`
	RiskScore             uint8                 `json:"risk_score"`
	HealthScore           uint8                 `json:"health_score"`
	TrendAnalysis         TrendAnalysis         `json:"trend_analysis"`
	BenchmarkComparison   BenchmarkComparison   `json:"benchmark_comparison"`
	FutureProjections     FutureProjections     `json:"future_projections"`
	LastUpdated           time.Time             `json:"last_updated"`
}

// CategoryMetrics represents metrics for a specific category
type CategoryMetrics struct {
	Category          ProposalCategory `json:"category"`
	TotalAllocated    sdk.Coin         `json:"total_allocated"`
	TotalSpent        sdk.Coin         `json:"total_spent"`
	Remaining         sdk.Coin         `json:"remaining"`
	Utilization       math.LegacyDec   `json:"utilization"`
	Proposals         uint64           `json:"proposals"`
	SuccessfulProposals uint64         `json:"successful_proposals"`
	SuccessRate       math.LegacyDec   `json:"success_rate"`
	AverageAmount     sdk.Coin         `json:"average_amount"`
	AverageTime       time.Duration    `json:"average_time"`
	ImpactScore       uint8            `json:"impact_score"`
	EfficiencyScore   uint8            `json:"efficiency_score"`
	SatisfactionScore uint8            `json:"satisfaction_score"`
	Trend             TrendDirection   `json:"trend"`
	LastUpdated       time.Time        `json:"last_updated"`
}

// DevCategoryMetrics represents metrics for development categories
type DevCategoryMetrics struct {
	Category          DevelopmentCategory `json:"category"`
	TotalAllocated    sdk.Coin            `json:"total_allocated"`
	TotalSpent        sdk.Coin            `json:"total_spent"`
	Remaining         sdk.Coin            `json:"remaining"`
	Utilization       math.LegacyDec      `json:"utilization"`
	Projects          uint64              `json:"projects"`
	CompletedProjects uint64              `json:"completed_projects"`
	CompletionRate    math.LegacyDec      `json:"completion_rate"`
	AverageAmount     sdk.Coin            `json:"average_amount"`
	AverageTime       time.Duration       `json:"average_time"`
	QualityScore      uint8               `json:"quality_score"`
	SecurityScore     uint8               `json:"security_score"`
	InnovationScore   uint8               `json:"innovation_score"`
	ROIScore          uint8               `json:"roi_score"`
	Trend             TrendDirection      `json:"trend"`
	LastUpdated       time.Time           `json:"last_updated"`
}

// RecipientMetrics represents metrics for fund recipients
type RecipientMetrics struct {
	Address           sdk.AccAddress `json:"address"`
	TotalReceived     sdk.Coin       `json:"total_received"`
	TotalProposals    uint64         `json:"total_proposals"`
	SuccessfulProposals uint64       `json:"successful_proposals"`
	SuccessRate       math.LegacyDec `json:"success_rate"`
	AverageAmount     sdk.Coin       `json:"average_amount"`
	Categories        []string       `json:"categories"`
	ReputationScore   uint8          `json:"reputation_score"`
	ComplianceScore   uint8          `json:"compliance_score"`
	ImpactScore       uint8          `json:"impact_score"`
	LastActivity      time.Time      `json:"last_activity"`
	Trend             TrendDirection `json:"trend"`
	Risk              uint8          `json:"risk"`
	Verified          bool           `json:"verified"`
}

// DeveloperMetrics represents metrics for developers
type DeveloperMetrics struct {
	Address           sdk.AccAddress `json:"address"`
	TotalEarned       sdk.Coin       `json:"total_earned"`
	TotalProjects     uint64         `json:"total_projects"`
	CompletedProjects uint64         `json:"completed_projects"`
	CompletionRate    math.LegacyDec `json:"completion_rate"`
	AverageAmount     sdk.Coin       `json:"average_amount"`
	Specializations   []string       `json:"specializations"`
	QualityScore      uint8          `json:"quality_score"`
	SecurityScore     uint8          `json:"security_score"`
	ReliabilityScore  uint8          `json:"reliability_score"`
	InnovationScore   uint8          `json:"innovation_score"`
	LastActivity      time.Time      `json:"last_activity"`
	Trend             TrendDirection `json:"trend"`
	Risk              uint8          `json:"risk"`
	Verified          bool           `json:"verified"`
}

// TransactionMetrics represents transaction metrics
type TransactionMetrics struct {
	TxID        string                 `json:"tx_id"`
	From        sdk.AccAddress         `json:"from"`
	To          sdk.AccAddress         `json:"to"`
	Amount      sdk.Coin               `json:"amount"`
	Type        TransactionType        `json:"type"`
	Category    string                 `json:"category"`
	Status      TransactionStatus      `json:"status"`
	Timestamp   time.Time              `json:"timestamp"`
	BlockHeight int64                  `json:"block_height"`
	Confirmed   bool                   `json:"confirmed"`
	Audited     bool                   `json:"audited"`
	Public      bool                   `json:"public"`
	Impact      uint8                  `json:"impact"`
	Risk        uint8                  `json:"risk"`
}

// DevTransactionMetrics represents development transaction metrics
type DevTransactionMetrics struct {
	TxID        string                 `json:"tx_id"`
	From        sdk.AccAddress         `json:"from"`
	To          sdk.AccAddress         `json:"to"`
	Amount      sdk.Coin               `json:"amount"`
	Type        DevelopmentTxType      `json:"type"`
	Category    DevelopmentCategory    `json:"category"`
	Status      DevelopmentTxStatus    `json:"status"`
	Timestamp   time.Time              `json:"timestamp"`
	BlockHeight int64                  `json:"block_height"`
	ProjectID   uint64                 `json:"project_id"`
	Phase       string                 `json:"phase"`
	Milestone   string                 `json:"milestone"`
	Quality     uint8                  `json:"quality"`
	Security    uint8                  `json:"security"`
	Confirmed   bool                   `json:"confirmed"`
	Reviewed    bool                   `json:"reviewed"`
	Audited     bool                   `json:"audited"`
}

// MonthlyMetrics represents monthly spending metrics
type MonthlyMetrics struct {
	Year         int                    `json:"year"`
	Month        int                    `json:"month"`
	TotalSpent   sdk.Coin               `json:"total_spent"`
	Proposals    uint64                 `json:"proposals"`
	Categories   []CategoryMetrics      `json:"categories"`
	Recipients   []RecipientMetrics     `json:"recipients"`
	Efficiency   math.LegacyDec         `json:"efficiency"`
	Impact       uint8                  `json:"impact"`
	Satisfaction uint8                  `json:"satisfaction"`
	Trend        TrendDirection         `json:"trend"`
	Comparison   ComparisonMetrics      `json:"comparison"`
}

// QuarterlyReport represents quarterly report
type QuarterlyReport struct {
	Year            int                 `json:"year"`
	Quarter         int                 `json:"quarter"`
	TotalSpent      sdk.Coin            `json:"total_spent"`
	TotalProposals  uint64              `json:"total_proposals"`
	Categories      []CategoryMetrics   `json:"categories"`
	TopRecipients   []RecipientMetrics  `json:"top_recipients"`
	Achievements    []Achievement       `json:"achievements"`
	Challenges      []Challenge         `json:"challenges"`
	Improvements    []Improvement       `json:"improvements"`
	NextQuarter     []Objective         `json:"next_quarter"`
	HealthScore     uint8               `json:"health_score"`
	EfficiencyScore uint8               `json:"efficiency_score"`
	ImpactScore     uint8               `json:"impact_score"`
	ReportDate      time.Time           `json:"report_date"`
	Published       bool                `json:"published"`
}

// AnnualReport represents annual report
type AnnualReport struct {
	Year            int                 `json:"year"`
	TotalSpent      sdk.Coin            `json:"total_spent"`
	TotalProposals  uint64              `json:"total_proposals"`
	Categories      []CategoryMetrics   `json:"categories"`
	TopRecipients   []RecipientMetrics  `json:"top_recipients"`
	Achievements    []Achievement       `json:"achievements"`
	Challenges      []Challenge         `json:"challenges"`
	Improvements    []Improvement       `json:"improvements"`
	NextYear        []Objective         `json:"next_year"`
	HealthScore     uint8               `json:"health_score"`
	EfficiencyScore uint8               `json:"efficiency_score"`
	ImpactScore     uint8               `json:"impact_score"`
	OverallRating   uint8               `json:"overall_rating"`
	ReportDate      time.Time           `json:"report_date"`
	Published       bool                `json:"published"`
	Audited         bool                `json:"audited"`
}

// TrendDirection represents trend direction
type TrendDirection string

const (
	TrendUp      TrendDirection = "up"
	TrendDown    TrendDirection = "down"
	TrendStable  TrendDirection = "stable"
	TrendVolatile TrendDirection = "volatile"
)

// Alert represents an alert
type Alert struct {
	ID          uint64    `json:"id"`
	Type        string    `json:"type"`
	Severity    string    `json:"severity"`
	Title       string    `json:"title"`
	Message     string    `json:"message"`
	Category    string    `json:"category"`
	Timestamp   time.Time `json:"timestamp"`
	Acknowledged bool     `json:"acknowledged"`
	Resolved    bool      `json:"resolved"`
	Action      string    `json:"action"`
	Responsible string    `json:"responsible"`
	Deadline    time.Time `json:"deadline"`
}

// Notification represents a notification
type Notification struct {
	ID          uint64    `json:"id"`
	Type        string    `json:"type"`
	Title       string    `json:"title"`
	Message     string    `json:"message"`
	Category    string    `json:"category"`
	Timestamp   time.Time `json:"timestamp"`
	Read        bool      `json:"read"`
	Important   bool      `json:"important"`
	Action      string    `json:"action"`
	Link        string    `json:"link"`
	Expiry      time.Time `json:"expiry"`
}

// TransparencyMetrics represents transparency metrics
type TransparencyMetrics struct {
	OverallScore      uint8          `json:"overall_score"`
	PublicDisclosure  uint8          `json:"public_disclosure"`
	DocumentAvailability uint8       `json:"document_availability"`
	RealTimeUpdates   uint8          `json:"real_time_updates"`
	CommunityInput    uint8          `json:"community_input"`
	AuditFrequency    uint8          `json:"audit_frequency"`
	ResponseTime      time.Duration  `json:"response_time"`
	OpenMeetings      uint8          `json:"open_meetings"`
	Accountability    uint8          `json:"accountability"`
	Accessibility     uint8          `json:"accessibility"`
	Improvement       TrendDirection `json:"improvement"`
	LastAudit         time.Time      `json:"last_audit"`
	NextAudit         time.Time      `json:"next_audit"`
}

// GovernanceMetrics represents governance metrics
type GovernanceMetrics struct {
	ParticipationRate     math.LegacyDec `json:"participation_rate"`
	VotingTurnout         math.LegacyDec `json:"voting_turnout"`
	ProposalSuccessRate   math.LegacyDec `json:"proposal_success_rate"`
	AverageSigningTime    time.Duration  `json:"average_signing_time"`
	ConsensusRate         math.LegacyDec `json:"consensus_rate"`
	DisputeRate           math.LegacyDec `json:"dispute_rate"`
	GovernanceHealth      uint8          `json:"governance_health"`
	DecisionQuality       uint8          `json:"decision_quality"`
	ProcessEfficiency     uint8          `json:"process_efficiency"`
	StakeholderSatisfaction uint8        `json:"stakeholder_satisfaction"`
	ActiveSigners         uint64         `json:"active_signers"`
	InactiveSigners       uint64         `json:"inactive_signers"`
	SignerTurnover        math.LegacyDec `json:"signer_turnover"`
	LastElection          time.Time      `json:"last_election"`
	NextElection          time.Time      `json:"next_election"`
}

// QualityMetrics represents quality metrics
type QualityMetrics struct {
	OverallScore      uint8          `json:"overall_score"`
	CodeQuality       uint8          `json:"code_quality"`
	Documentation     uint8          `json:"documentation"`
	Testing           uint8          `json:"testing"`
	Security          uint8          `json:"security"`
	Performance       uint8          `json:"performance"`
	Maintainability   uint8          `json:"maintainability"`
	Usability         uint8          `json:"usability"`
	Reliability       uint8          `json:"reliability"`
	Scalability       uint8          `json:"scalability"`
	Improvement       TrendDirection `json:"improvement"`
	LastAssessment    time.Time      `json:"last_assessment"`
	NextAssessment    time.Time      `json:"next_assessment"`
}

// SecurityMetrics represents security metrics
type SecurityMetrics struct {
	SecurityScore     uint8          `json:"security_score"`
	VulnerabilityCount uint64        `json:"vulnerability_count"`
	CriticalVulns     uint64         `json:"critical_vulnerabilities"`
	HighRiskVulns     uint64         `json:"high_risk_vulnerabilities"`
	MediumRiskVulns   uint64         `json:"medium_risk_vulnerabilities"`
	LowRiskVulns      uint64         `json:"low_risk_vulnerabilities"`
	ResolvedVulns     uint64         `json:"resolved_vulnerabilities"`
	PendingVulns      uint64         `json:"pending_vulnerabilities"`
	AuditCoverage     uint8          `json:"audit_coverage"`
	ComplianceScore   uint8          `json:"compliance_score"`
	IncidentCount     uint64         `json:"incident_count"`
	LastAudit         time.Time      `json:"last_audit"`
	NextAudit         time.Time      `json:"next_audit"`
	Trend             TrendDirection `json:"trend"`
}

// InnovationMetrics represents innovation metrics
type InnovationMetrics struct {
	InnovationScore   uint8          `json:"innovation_score"`
	NewFeatures       uint64         `json:"new_features"`
	ExperimentalProjects uint64      `json:"experimental_projects"`
	ResearchProjects  uint64         `json:"research_projects"`
	Patents           uint64         `json:"patents"`
	Publications      uint64         `json:"publications"`
	Collaborations    uint64         `json:"collaborations"`
	TechAdoption      uint8          `json:"tech_adoption"`
	DisruptiveIdeas   uint64         `json:"disruptive_ideas"`
	TimeToMarket      time.Duration  `json:"time_to_market"`
	Trend             TrendDirection `json:"trend"`
	LastReview        time.Time      `json:"last_review"`
	NextReview        time.Time      `json:"next_review"`
}

// DeliveryMetrics represents delivery metrics
type DeliveryMetrics struct {
	OnTimeDelivery    math.LegacyDec `json:"on_time_delivery"`
	BudgetAdherence   math.LegacyDec `json:"budget_adherence"`
	ScopeCreep        math.LegacyDec `json:"scope_creep"`
	DefectRate        math.LegacyDec `json:"defect_rate"`
	CustomerSatisfaction uint8       `json:"customer_satisfaction"`
	TeamProductivity  uint8          `json:"team_productivity"`
	ResourceUtilization uint8        `json:"resource_utilization"`
	ProcessEfficiency uint8          `json:"process_efficiency"`
	QualityScore      uint8          `json:"quality_score"`
	RiskMitigation    uint8          `json:"risk_mitigation"`
	AverageDeliveryTime time.Duration `json:"average_delivery_time"`
	Trend             TrendDirection `json:"trend"`
}

// ROIMetrics represents ROI metrics
type ROIMetrics struct {
	ROI               math.LegacyDec `json:"roi"`
	NPV               math.LegacyDec `json:"npv"`
	IRR               math.LegacyDec `json:"irr"`
	PaybackPeriod     time.Duration  `json:"payback_period"`
	TotalCost         sdk.Coin       `json:"total_cost"`
	TotalBenefit      sdk.Coin       `json:"total_benefit"`
	CostPerFeature    sdk.Coin       `json:"cost_per_feature"`
	ValueCreated      sdk.Coin       `json:"value_created"`
	EfficiencyGains   math.LegacyDec `json:"efficiency_gains"`
	UserAdoption      uint8          `json:"user_adoption"`
	MarketImpact      uint8          `json:"market_impact"`
	StrategicValue    uint8          `json:"strategic_value"`
	LastCalculation   time.Time      `json:"last_calculation"`
	Trend             TrendDirection `json:"trend"`
}

// TrendAnalysis represents trend analysis
type TrendAnalysis struct {
	ShortTermTrend    TrendDirection `json:"short_term_trend"`
	MediumTermTrend   TrendDirection `json:"medium_term_trend"`
	LongTermTrend     TrendDirection `json:"long_term_trend"`
	SeasonalPattern   string         `json:"seasonal_pattern"`
	GrowthRate        math.LegacyDec `json:"growth_rate"`
	Volatility        uint8          `json:"volatility"`
	Correlation       math.LegacyDec `json:"correlation"`
	Momentum          uint8          `json:"momentum"`
	Forecast          string         `json:"forecast"`
	Confidence        uint8          `json:"confidence"`
	LastAnalysis      time.Time      `json:"last_analysis"`
	NextAnalysis      time.Time      `json:"next_analysis"`
}

// BenchmarkComparison represents benchmark comparison
type BenchmarkComparison struct {
	IndustryBenchmark  math.LegacyDec `json:"industry_benchmark"`
	PeerComparison     math.LegacyDec `json:"peer_comparison"`
	BestPractice       math.LegacyDec `json:"best_practice"`
	PerformanceGap     math.LegacyDec `json:"performance_gap"`
	Ranking            uint64         `json:"ranking"`
	TotalCompared      uint64         `json:"total_compared"`
	StrongPoints       []string       `json:"strong_points"`
	WeakPoints         []string       `json:"weak_points"`
	Recommendations    []string       `json:"recommendations"`
	ImprovementPlan    string         `json:"improvement_plan"`
	LastComparison     time.Time      `json:"last_comparison"`
	NextComparison     time.Time      `json:"next_comparison"`
}

// FutureProjections represents future projections
type FutureProjections struct {
	NextMonth         ProjectionData `json:"next_month"`
	NextQuarter       ProjectionData `json:"next_quarter"`
	NextYear          ProjectionData `json:"next_year"`
	FiveYear          ProjectionData `json:"five_year"`
	Assumptions       []string       `json:"assumptions"`
	Scenarios         []Scenario     `json:"scenarios"`
	RiskFactors       []string       `json:"risk_factors"`
	Opportunities     []string       `json:"opportunities"`
	Confidence        uint8          `json:"confidence"`
	LastProjection    time.Time      `json:"last_projection"`
	NextProjection    time.Time      `json:"next_projection"`
}

// ProjectionData represents projection data
type ProjectionData struct {
	EstimatedSpending sdk.Coin       `json:"estimated_spending"`
	EstimatedProjects uint64         `json:"estimated_projects"`
	EstimatedImpact   uint8          `json:"estimated_impact"`
	Confidence        uint8          `json:"confidence"`
	Range             Range          `json:"range"`
}

// Range represents a range of values
type Range struct {
	Min sdk.Coin `json:"min"`
	Max sdk.Coin `json:"max"`
}

// Scenario represents a scenario
type Scenario struct {
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Probability  uint8      `json:"probability"`
	Impact       uint8      `json:"impact"`
	Spending     sdk.Coin   `json:"spending"`
	Projects     uint64     `json:"projects"`
	Outcomes     []string   `json:"outcomes"`
	Mitigation   string     `json:"mitigation"`
}

// ComparisonMetrics represents comparison metrics
type ComparisonMetrics struct {
	PreviousMonth   math.LegacyDec `json:"previous_month"`
	PreviousQuarter math.LegacyDec `json:"previous_quarter"`
	PreviousYear    math.LegacyDec `json:"previous_year"`
	Average         math.LegacyDec `json:"average"`
	Median          math.LegacyDec `json:"median"`
	BestMonth       math.LegacyDec `json:"best_month"`
	WorstMonth      math.LegacyDec `json:"worst_month"`
	Trend           TrendDirection `json:"trend"`
}

// Achievement represents an achievement
type Achievement struct {
	ID          uint64    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Impact      uint8     `json:"impact"`
	Date        time.Time `json:"date"`
	Verified    bool      `json:"verified"`
	Public      bool      `json:"public"`
}

// Challenge represents a challenge
type Challenge struct {
	ID          uint64    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Severity    uint8     `json:"severity"`
	Date        time.Time `json:"date"`
	Resolved    bool      `json:"resolved"`
	Resolution  string    `json:"resolution"`
	Lessons     []string  `json:"lessons"`
}

// Improvement represents an improvement
type Improvement struct {
	ID          uint64    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Impact      uint8     `json:"impact"`
	Date        time.Time `json:"date"`
	Implemented bool      `json:"implemented"`
	Results     string    `json:"results"`
	Metrics     []string  `json:"metrics"`
}

// Objective represents an objective
type Objective struct {
	ID          uint64    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Priority    uint8     `json:"priority"`
	Deadline    time.Time `json:"deadline"`
	Assigned    string    `json:"assigned"`
	Progress    uint8     `json:"progress"`
	KPIs        []string  `json:"kpis"`
}

// PredictiveAnalytics represents predictive analytics
type PredictiveAnalytics struct {
	FundingPrediction    FundingPrediction    `json:"funding_prediction"`
	SuccessPrediction    SuccessPrediction    `json:"success_prediction"`
	RiskPrediction       RiskPrediction       `json:"risk_prediction"`
	TrendPrediction      TrendPrediction      `json:"trend_prediction"`
	OptimizationSuggestions []OptimizationSuggestion `json:"optimization_suggestions"`
	AnomalyDetection     AnomalyDetection     `json:"anomaly_detection"`
	SeasonalAdjustments  SeasonalAdjustments  `json:"seasonal_adjustments"`
	ModelAccuracy        uint8                `json:"model_accuracy"`
	LastTraining         time.Time            `json:"last_training"`
	NextTraining         time.Time            `json:"next_training"`
}

// FundingPrediction represents funding prediction
type FundingPrediction struct {
	NextMonthFunding    sdk.Coin  `json:"next_month_funding"`
	NextQuarterFunding  sdk.Coin  `json:"next_quarter_funding"`
	NextYearFunding     sdk.Coin  `json:"next_year_funding"`
	OptimalAllocation   []string  `json:"optimal_allocation"`
	Confidence          uint8     `json:"confidence"`
	LastPrediction      time.Time `json:"last_prediction"`
}

// SuccessPrediction represents success prediction
type SuccessPrediction struct {
	SuccessRate         math.LegacyDec `json:"success_rate"`
	HighRiskProposals   []uint64       `json:"high_risk_proposals"`
	HighPotentialProposals []uint64    `json:"high_potential_proposals"`
	KeySuccessFactors   []string       `json:"key_success_factors"`
	RiskFactors         []string       `json:"risk_factors"`
	Confidence          uint8          `json:"confidence"`
	LastPrediction      time.Time      `json:"last_prediction"`
}

// RiskPrediction represents risk prediction
type RiskPrediction struct {
	OverallRisk         uint8     `json:"overall_risk"`
	FinancialRisk       uint8     `json:"financial_risk"`
	OperationalRisk     uint8     `json:"operational_risk"`
	ReputationalRisk    uint8     `json:"reputational_risk"`
	ComplianceRisk      uint8     `json:"compliance_risk"`
	TechnicalRisk       uint8     `json:"technical_risk"`
	MarketRisk          uint8     `json:"market_risk"`
	RiskTrend           TrendDirection `json:"risk_trend"`
	MitigationStrategies []string  `json:"mitigation_strategies"`
	Confidence          uint8     `json:"confidence"`
	LastPrediction      time.Time `json:"last_prediction"`
}

// TrendPrediction represents trend prediction
type TrendPrediction struct {
	ShortTermTrend      TrendDirection `json:"short_term_trend"`
	MediumTermTrend     TrendDirection `json:"medium_term_trend"`
	LongTermTrend       TrendDirection `json:"long_term_trend"`
	TrendStrength       uint8          `json:"trend_strength"`
	TrendReliability    uint8          `json:"trend_reliability"`
	InfluencingFactors  []string       `json:"influencing_factors"`
	TrendDrivers        []string       `json:"trend_drivers"`
	Confidence          uint8          `json:"confidence"`
	LastPrediction      time.Time      `json:"last_prediction"`
}

// OptimizationSuggestion represents optimization suggestion
type OptimizationSuggestion struct {
	ID              uint64    `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Category        string    `json:"category"`
	Impact          uint8     `json:"impact"`
	Effort          uint8     `json:"effort"`
	Priority        uint8     `json:"priority"`
	EstimatedSavings sdk.Coin  `json:"estimated_savings"`
	EstimatedGains   sdk.Coin  `json:"estimated_gains"`
	Timeline        time.Duration `json:"timeline"`
	Dependencies    []string  `json:"dependencies"`
	Risks           []string  `json:"risks"`
	Benefits        []string  `json:"benefits"`
	Confidence      uint8     `json:"confidence"`
	Generated       time.Time `json:"generated"`
}

// AnomalyDetection represents anomaly detection
type AnomalyDetection struct {
	AnomaliesDetected   uint64         `json:"anomalies_detected"`
	CriticalAnomalies   uint64         `json:"critical_anomalies"`
	RecentAnomalies     []Anomaly      `json:"recent_anomalies"`
	AnomalyTrend        TrendDirection `json:"anomaly_trend"`
	FalsePositiveRate   math.LegacyDec `json:"false_positive_rate"`
	DetectionAccuracy   uint8          `json:"detection_accuracy"`
	LastDetection       time.Time      `json:"last_detection"`
	NextScan            time.Time      `json:"next_scan"`
}

// Anomaly represents an anomaly
type Anomaly struct {
	ID          uint64    `json:"id"`
	Type        string    `json:"type"`
	Severity    string    `json:"severity"`
	Description string    `json:"description"`
	Value       string    `json:"value"`
	Expected    string    `json:"expected"`
	Deviation   string    `json:"deviation"`
	Confidence  uint8     `json:"confidence"`
	Detected    time.Time `json:"detected"`
	Resolved    bool      `json:"resolved"`
	Resolution  string    `json:"resolution"`
}

// SeasonalAdjustments represents seasonal adjustments
type SeasonalAdjustments struct {
	SeasonalFactors     []SeasonalFactor `json:"seasonal_factors"`
	AdjustedForecast    []AdjustedForecast `json:"adjusted_forecast"`
	SeasonalityStrength uint8            `json:"seasonality_strength"`
	PatternRecognition  uint8            `json:"pattern_recognition"`
	LastAdjustment      time.Time        `json:"last_adjustment"`
	NextAdjustment      time.Time        `json:"next_adjustment"`
}

// SeasonalFactor represents seasonal factor
type SeasonalFactor struct {
	Period     string         `json:"period"`
	Factor     math.LegacyDec `json:"factor"`
	Confidence uint8          `json:"confidence"`
	Historical uint8          `json:"historical"`
}

// AdjustedForecast represents adjusted forecast
type AdjustedForecast struct {
	Period          string         `json:"period"`
	RawForecast     sdk.Coin       `json:"raw_forecast"`
	AdjustedForecast sdk.Coin      `json:"adjusted_forecast"`
	Adjustment      math.LegacyDec `json:"adjustment"`
	Confidence      uint8          `json:"confidence"`
}

// Real-time update intervals
const (
	DashboardUpdateInterval     = time.Minute * 5  // 5 minutes
	MetricsUpdateInterval       = time.Minute * 1  // 1 minute
	AlertCheckInterval          = time.Second * 30 // 30 seconds
	TransactionUpdateInterval   = time.Second * 10 // 10 seconds
	PredictionUpdateInterval    = time.Hour * 1    // 1 hour
	ReportGenerationInterval    = time.Hour * 24   // 24 hours
)