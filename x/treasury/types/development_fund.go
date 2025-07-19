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
	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

// DevelopmentFundProposal represents a proposal for using development funds
type DevelopmentFundProposal struct {
	ProposalID          uint64              `json:"proposal_id"`
	Proposer            sdk.AccAddress      `json:"proposer"`
	Title               string              `json:"title"`
	Description         string              `json:"description"`
	Category            DevelopmentCategory `json:"category"`
	Priority            PriorityLevel       `json:"priority"`
	RequestedAmount     sdk.Coin            `json:"requested_amount"`
	TechnicalSpecs      TechnicalSpecs      `json:"technical_specs"`
	Timeline            ProjectTimeline     `json:"timeline"`
	Team                []TeamMember        `json:"team"`
	Deliverables        []Deliverable       `json:"deliverables"`
	Status              DevelopmentStatus   `json:"status"`
	SubmissionTime      time.Time           `json:"submission_time"`
	ReviewPeriod        time.Duration       `json:"review_period"`
	ReviewEndTime       time.Time           `json:"review_end_time"`
	ApprovalTime        *time.Time          `json:"approval_time,omitempty"`
	ExecutionStartTime  *time.Time          `json:"execution_start_time,omitempty"`
	CompletionTime      *time.Time          `json:"completion_time,omitempty"`
	TechnicalReview     TechnicalReview     `json:"technical_review"`
	FinancialReview     FinancialReview     `json:"financial_review"`
	SecurityReview      SecurityReview      `json:"security_review"`
	CodeQuality         CodeQuality         `json:"code_quality"`
	TestingCoverage     TestingCoverage     `json:"testing_coverage"`
	Documentation       Documentation       `json:"documentation"`
	Impact              ImpactAssessment    `json:"impact"`
	RiskAssessment      RiskAssessment      `json:"risk_assessment"`
	ApprovalScore       uint8               `json:"approval_score"`
	CommunityFeedback   []CommunityFeedback `json:"community_feedback"`
	TransparencyLevel   uint8               `json:"transparency_level"`
}

// DevelopmentCategory defines the types of development fund proposals
type DevelopmentCategory string

const (
	CategoryCoreBlockchain    DevelopmentCategory = "core_blockchain"
	CategorySmartContracts    DevelopmentCategory = "smart_contracts"
	CategoryUserInterface     DevelopmentCategory = "user_interface"
	CategoryMobileApp         DevelopmentCategory = "mobile_app"
	CategoryWalletDevelopment DevelopmentCategory = "wallet_development"
	CategorySecurityAudits    DevelopmentCategory = "security_audits"
	CategoryPerformanceOpt    DevelopmentCategory = "performance_optimization"
	CategoryTesting           DevelopmentCategory = "testing"
	CategoryDocumentation     DevelopmentCategory = "documentation"
	CategoryDevTools          DevelopmentCategory = "dev_tools"
	CategoryInfrastructure    DevelopmentCategory = "infrastructure"
	CategoryIntegrations      DevelopmentCategory = "integrations"
	CategoryResearch          DevelopmentCategory = "research"
	CategoryPrototyping       DevelopmentCategory = "prototyping"
	CategoryMaintenance       DevelopmentCategory = "maintenance"
	CategoryEmergencyFixes    DevelopmentCategory = "emergency_fixes"
)

// PriorityLevel defines the priority of development proposals
type PriorityLevel string

const (
	PriorityCritical PriorityLevel = "critical"
	PriorityHigh     PriorityLevel = "high"
	PriorityMedium   PriorityLevel = "medium"
	PriorityLow      PriorityLevel = "low"
	PriorityNice     PriorityLevel = "nice_to_have"
)

// DevelopmentStatus defines the status of development proposals
type DevelopmentStatus string

const (
	DevStatusPending        DevelopmentStatus = "pending"
	DevStatusUnderReview    DevelopmentStatus = "under_review"
	DevStatusApproved       DevelopmentStatus = "approved"
	DevStatusRejected       DevelopmentStatus = "rejected"
	DevStatusInProgress     DevelopmentStatus = "in_progress"
	DevStatusTesting        DevelopmentStatus = "testing"
	DevStatusCompleted      DevelopmentStatus = "completed"
	DevStatusDeployed       DevelopmentStatus = "deployed"
	DevStatusMaintenance    DevelopmentStatus = "maintenance"
	DevStatusCancelled      DevelopmentStatus = "cancelled"
	DevStatusOnHold         DevelopmentStatus = "on_hold"
	DevStatusRevisionNeeded DevelopmentStatus = "revision_needed"
)

// TechnicalSpecs represents technical specifications
type TechnicalSpecs struct {
	Languages         []string            `json:"languages"`
	Frameworks        []string            `json:"frameworks"`
	Dependencies      []string            `json:"dependencies"`
	Architecture      string              `json:"architecture"`
	Database          string              `json:"database"`
	APIs              []string            `json:"apis"`
	SecurityFeatures  []string            `json:"security_features"`
	Performance       PerformanceSpecs    `json:"performance"`
	Scalability       ScalabilitySpecs    `json:"scalability"`
	Compatibility     CompatibilitySpecs  `json:"compatibility"`
	Testing           TestingSpecs        `json:"testing"`
	Documentation     DocumentationSpecs  `json:"documentation"`
	Deployment        DeploymentSpecs     `json:"deployment"`
	Monitoring        MonitoringSpecs     `json:"monitoring"`
}

// PerformanceSpecs represents performance requirements
type PerformanceSpecs struct {
	ExpectedTPS       uint64        `json:"expected_tps"`
	LatencyTarget     time.Duration `json:"latency_target"`
	MemoryUsage       string        `json:"memory_usage"`
	CPUUsage          string        `json:"cpu_usage"`
	NetworkBandwidth  string        `json:"network_bandwidth"`
	StorageRequirements string      `json:"storage_requirements"`
}

// ScalabilitySpecs represents scalability requirements
type ScalabilitySpecs struct {
	HorizontalScaling bool   `json:"horizontal_scaling"`
	VerticalScaling   bool   `json:"vertical_scaling"`
	LoadBalancing     bool   `json:"load_balancing"`
	Caching           string `json:"caching"`
	CDN               bool   `json:"cdn"`
	Microservices     bool   `json:"microservices"`
}

// CompatibilitySpecs represents compatibility requirements
type CompatibilitySpecs struct {
	Platforms       []string `json:"platforms"`
	Browsers        []string `json:"browsers"`
	MobileOS        []string `json:"mobile_os"`
	DatabaseSystems []string `json:"database_systems"`
	APIVersions     []string `json:"api_versions"`
	BackwardCompat  bool     `json:"backward_compatibility"`
}

// TestingSpecs represents testing requirements
type TestingSpecs struct {
	UnitTests        bool   `json:"unit_tests"`
	IntegrationTests bool   `json:"integration_tests"`
	E2ETests         bool   `json:"e2e_tests"`
	LoadTests        bool   `json:"load_tests"`
	SecurityTests    bool   `json:"security_tests"`
	CoverageTarget   uint8  `json:"coverage_target"`
	TestFrameworks   []string `json:"test_frameworks"`
}

// DocumentationSpecs represents documentation requirements
type DocumentationSpecs struct {
	APIDocumentation  bool `json:"api_documentation"`
	UserGuides        bool `json:"user_guides"`
	DeveloperGuides   bool `json:"developer_guides"`
	TechnicalSpecs    bool `json:"technical_specs"`
	InstallationGuide bool `json:"installation_guide"`
	TroubleshootingGuide bool `json:"troubleshooting_guide"`
}

// DeploymentSpecs represents deployment requirements
type DeploymentSpecs struct {
	ContainerSupport bool     `json:"container_support"`
	CloudPlatforms   []string `json:"cloud_platforms"`
	CICD             bool     `json:"ci_cd"`
	RollbackStrategy string   `json:"rollback_strategy"`
	BlueGreenDeploy  bool     `json:"blue_green_deploy"`
	MonitoringIntegration bool `json:"monitoring_integration"`
}

// MonitoringSpecs represents monitoring requirements
type MonitoringSpecs struct {
	Metrics     []string `json:"metrics"`
	Logging     bool     `json:"logging"`
	Alerting    bool     `json:"alerting"`
	Dashboards  bool     `json:"dashboards"`
	Tracing     bool     `json:"tracing"`
	HealthChecks bool    `json:"health_checks"`
}

// ProjectTimeline represents project timeline
type ProjectTimeline struct {
	StartDate         time.Time `json:"start_date"`
	EndDate           time.Time `json:"end_date"`
	Duration          time.Duration `json:"duration"`
	Phases            []Phase   `json:"phases"`
	Milestones        []Milestone `json:"milestones"`
	Dependencies      []string  `json:"dependencies"`
	CriticalPath      []string  `json:"critical_path"`
	BufferTime        time.Duration `json:"buffer_time"`
	ReviewCheckpoints []time.Time `json:"review_checkpoints"`
}

// Phase represents a project phase
type Phase struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Status      string    `json:"status"`
	Progress    uint8     `json:"progress"`
	Deliverables []string `json:"deliverables"`
}

// TeamMember represents a team member
type TeamMember struct {
	Address      sdk.AccAddress `json:"address"`
	Role         string         `json:"role"`
	Expertise    []string       `json:"expertise"`
	Allocation   math.LegacyDec `json:"allocation"`
	Compensation sdk.Coin       `json:"compensation"`
	KPIs         []string       `json:"kpis"`
	Reputation   uint8          `json:"reputation"`
	Experience   uint8          `json:"experience"`
	Availability uint8          `json:"availability"`
}

// Deliverable represents a project deliverable
type Deliverable struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status"`
	Progress    uint8     `json:"progress"`
	Quality     uint8     `json:"quality"`
	TestStatus  string    `json:"test_status"`
	ReviewStatus string   `json:"review_status"`
	DeployStatus string   `json:"deploy_status"`
}

// TechnicalReview represents technical review results
type TechnicalReview struct {
	Reviewer        sdk.AccAddress `json:"reviewer"`
	ReviewDate      time.Time      `json:"review_date"`
	Score           uint8          `json:"score"`
	Approved        bool           `json:"approved"`
	Comments        string         `json:"comments"`
	Recommendations []string       `json:"recommendations"`
	TechnicalRisks  []string       `json:"technical_risks"`
	Complexity      uint8          `json:"complexity"`
	Feasibility     uint8          `json:"feasibility"`
	Innovation      uint8          `json:"innovation"`
	Maintainability uint8          `json:"maintainability"`
}

// FinancialReview represents financial review results
type FinancialReview struct {
	Reviewer        sdk.AccAddress `json:"reviewer"`
	ReviewDate      time.Time      `json:"review_date"`
	Score           uint8          `json:"score"`
	Approved        bool           `json:"approved"`
	BudgetAnalysis  string         `json:"budget_analysis"`
	CostBreakdown   []CostItem     `json:"cost_breakdown"`
	ROI             math.LegacyDec `json:"roi"`
	NPV             math.LegacyDec `json:"npv"`
	Payback         time.Duration  `json:"payback"`
	RiskFactors     []string       `json:"risk_factors"`
	Recommendations []string       `json:"recommendations"`
}

// CostItem represents a cost item
type CostItem struct {
	Category    string   `json:"category"`
	Amount      sdk.Coin `json:"amount"`
	Description string   `json:"description"`
	Recurring   bool     `json:"recurring"`
	Frequency   string   `json:"frequency"`
}

// SecurityReview represents security review results
type SecurityReview struct {
	Reviewer           sdk.AccAddress `json:"reviewer"`
	ReviewDate         time.Time      `json:"review_date"`
	Score              uint8          `json:"score"`
	Approved           bool           `json:"approved"`
	SecurityRisks      []SecurityRisk `json:"security_risks"`
	Vulnerabilities    []string       `json:"vulnerabilities"`
	Mitigations        []string       `json:"mitigations"`
	ComplianceStatus   string         `json:"compliance_status"`
	AuditRequired      bool           `json:"audit_required"`
	PenetrationTesting bool           `json:"penetration_testing"`
	Recommendations    []string       `json:"recommendations"`
}

// SecurityRisk represents a security risk
type SecurityRisk struct {
	ID          uint64 `json:"id"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
	Impact      string `json:"impact"`
	Likelihood  string `json:"likelihood"`
	Mitigation  string `json:"mitigation"`
	Status      string `json:"status"`
}

// CodeQuality represents code quality metrics
type CodeQuality struct {
	OverallScore      uint8    `json:"overall_score"`
	Maintainability   uint8    `json:"maintainability"`
	Reliability       uint8    `json:"reliability"`
	Security          uint8    `json:"security"`
	Performance       uint8    `json:"performance"`
	Duplication       uint8    `json:"duplication"`
	Complexity        uint8    `json:"complexity"`
	TechnicalDebt     uint8    `json:"technical_debt"`
	CodeCoverage      uint8    `json:"code_coverage"`
	DocumentationCov  uint8    `json:"documentation_coverage"`
	Standards         uint8    `json:"standards"`
	BestPractices     uint8    `json:"best_practices"`
	Issues            []string `json:"issues"`
	Improvements      []string `json:"improvements"`
}

// TestingCoverage represents testing coverage
type TestingCoverage struct {
	UnitTestCoverage        uint8    `json:"unit_test_coverage"`
	IntegrationTestCoverage uint8    `json:"integration_test_coverage"`
	E2ETestCoverage         uint8    `json:"e2e_test_coverage"`
	OverallCoverage         uint8    `json:"overall_coverage"`
	TestsPassed             uint64   `json:"tests_passed"`
	TestsFailed             uint64   `json:"tests_failed"`
	TestsSkipped            uint64   `json:"tests_skipped"`
	TestExecutionTime       time.Duration `json:"test_execution_time"`
	CriticalPathsCovered    uint8    `json:"critical_paths_covered"`
	EdgeCasesCovered        uint8    `json:"edge_cases_covered"`
	Performance             uint8    `json:"performance"`
	LoadTesting             uint8    `json:"load_testing"`
	SecurityTesting         uint8    `json:"security_testing"`
	Recommendations         []string `json:"recommendations"`
}

// Documentation represents documentation status
type Documentation struct {
	OverallScore      uint8    `json:"overall_score"`
	Completeness      uint8    `json:"completeness"`
	Clarity           uint8    `json:"clarity"`
	Accuracy          uint8    `json:"accuracy"`
	Examples          uint8    `json:"examples"`
	APIDocumentation  uint8    `json:"api_documentation"`
	UserGuides        uint8    `json:"user_guides"`
	DeveloperGuides   uint8    `json:"developer_guides"`
	InstallationGuide uint8    `json:"installation_guide"`
	TroubleshootingGuide uint8 `json:"troubleshooting_guide"`
	CodeComments      uint8    `json:"code_comments"`
	Changelog         uint8    `json:"changelog"`
	Recommendations   []string `json:"recommendations"`
}

// ImpactAssessment represents impact assessment
type ImpactAssessment struct {
	UserImpact        uint8    `json:"user_impact"`
	DeveloperImpact   uint8    `json:"developer_impact"`
	PlatformImpact    uint8    `json:"platform_impact"`
	SecurityImpact    uint8    `json:"security_impact"`
	PerformanceImpact uint8    `json:"performance_impact"`
	ScalabilityImpact uint8    `json:"scalability_impact"`
	MaintenanceImpact uint8    `json:"maintenance_impact"`
	CostImpact        uint8    `json:"cost_impact"`
	RevenueImpact     uint8    `json:"revenue_impact"`
	MarketImpact      uint8    `json:"market_impact"`
	CompetitiveAdv    uint8    `json:"competitive_advantage"`
	Innovation        uint8    `json:"innovation"`
	StrategicValue    uint8    `json:"strategic_value"`
	LongTermValue     uint8    `json:"long_term_value"`
	Justification     string   `json:"justification"`
	Metrics           []string `json:"metrics"`
}

// RiskAssessment represents risk assessment
type RiskAssessment struct {
	OverallRisk       uint8  `json:"overall_risk"`
	TechnicalRisk     uint8  `json:"technical_risk"`
	FinancialRisk     uint8  `json:"financial_risk"`
	SecurityRisk      uint8  `json:"security_risk"`
	TimelineRisk      uint8  `json:"timeline_risk"`
	ResourceRisk      uint8  `json:"resource_risk"`
	MarketRisk        uint8  `json:"market_risk"`
	CompetitiveRisk   uint8  `json:"competitive_risk"`
	RegulatoryRisk    uint8  `json:"regulatory_risk"`
	ReputationRisk    uint8  `json:"reputation_risk"`
	MitigationPlan    string `json:"mitigation_plan"`
	ContingencyPlan   string `json:"contingency_plan"`
	RiskFactors       []string `json:"risk_factors"`
	Assumptions       []string `json:"assumptions"`
	Dependencies      []string `json:"dependencies"`
	Success           uint8    `json:"success_probability"`
}

// CommunityFeedback represents community feedback
type CommunityFeedback struct {
	ID          uint64         `json:"id"`
	Submitter   sdk.AccAddress `json:"submitter"`
	Content     string         `json:"content"`
	Rating      uint8          `json:"rating"`
	Category    string         `json:"category"`
	Timestamp   time.Time      `json:"timestamp"`
	Verified    bool           `json:"verified"`
	Helpful     uint64         `json:"helpful"`
	Response    string         `json:"response,omitempty"`
	Addressed   bool           `json:"addressed"`
}

// DevelopmentFundBalance tracks the development fund balance
type DevelopmentFundBalance struct {
	TotalBalance     sdk.Coin  `json:"total_balance"`
	AllocatedAmount  sdk.Coin  `json:"allocated_amount"`
	AvailableAmount  sdk.Coin  `json:"available_amount"`
	PendingAmount    sdk.Coin  `json:"pending_amount"`
	EmergencyReserve sdk.Coin  `json:"emergency_reserve"`
	LastUpdateHeight int64     `json:"last_update_height"`
	LastUpdateTime   time.Time `json:"last_update_time"`
}

// DevelopmentFundTransaction represents a development fund transaction
type DevelopmentFundTransaction struct {
	TxID          string                 `json:"tx_id"`
	ProposalID    uint64                 `json:"proposal_id"`
	From          sdk.AccAddress         `json:"from"`
	To            sdk.AccAddress         `json:"to"`
	Amount        sdk.Coin               `json:"amount"`
	Type          DevelopmentTxType      `json:"type"`
	Category      DevelopmentCategory    `json:"category"`
	Description   string                 `json:"description"`
	Timestamp     time.Time              `json:"timestamp"`
	BlockHeight   int64                  `json:"block_height"`
	Status        DevelopmentTxStatus    `json:"status"`
	Phase         string                 `json:"phase"`
	Milestone     string                 `json:"milestone"`
	Deliverable   string                 `json:"deliverable"`
	QualityScore  uint8                  `json:"quality_score"`
	Approved      bool                   `json:"approved"`
	Reviewed      bool                   `json:"reviewed"`
	Audited       bool                   `json:"audited"`
}

// DevelopmentTxType defines the type of development transaction
type DevelopmentTxType string

const (
	DevTxTypeAllocation    DevelopmentTxType = "allocation"
	DevTxTypePayment       DevelopmentTxType = "payment"
	DevTxTypeMilestone     DevelopmentTxType = "milestone"
	DevTxTypeBonus         DevelopmentTxType = "bonus"
	DevTxTypeRefund        DevelopmentTxType = "refund"
	DevTxTypeReallocation  DevelopmentTxType = "reallocation"
	DevTxTypeEmergency     DevelopmentTxType = "emergency"
	DevTxTypeIncentive     DevelopmentTxType = "incentive"
	DevTxTypeAuditPayment  DevelopmentTxType = "audit_payment"
	DevTxTypeInfrastructure DevelopmentTxType = "infrastructure"
)

// DevelopmentTxStatus defines the status of development transaction
type DevelopmentTxStatus string

const (
	DevTxStatusPending   DevelopmentTxStatus = "pending"
	DevTxStatusApproved  DevelopmentTxStatus = "approved"
	DevTxStatusExecuted  DevelopmentTxStatus = "executed"
	DevTxStatusCompleted DevelopmentTxStatus = "completed"
	DevTxStatusFailed    DevelopmentTxStatus = "failed"
	DevTxStatusCancelled DevelopmentTxStatus = "cancelled"
	DevTxStatusRefunded  DevelopmentTxStatus = "refunded"
	DevTxStatusAudited   DevelopmentTxStatus = "audited"
)

// Storage keys for development fund
var (
	DevelopmentFundProposalKey     = collections.NewPrefix(200)
	DevelopmentFundBalanceKey      = collections.NewPrefix(201)
	DevelopmentFundTransactionKey  = collections.NewPrefix(202)
	DevelopmentFundGovernanceKey   = collections.NewPrefix(203)
	DevelopmentFundReviewKey       = collections.NewPrefix(204)
	DevelopmentFundProgressKey     = collections.NewPrefix(205)
	DevelopmentFundQualityKey      = collections.NewPrefix(206)
	DevelopmentFundAuditKey        = collections.NewPrefix(207)
	DevelopmentFundFeedbackKey     = collections.NewPrefix(208)
	DevelopmentFundMetricsKey      = collections.NewPrefix(209)
)

// Module account names for development fund
const (
	DevelopmentFundModuleName     = "development_fund"
	DevelopmentFundPoolName       = "development_fund_pool"
	DevelopmentFundEscrowName     = "development_fund_escrow"
	DevelopmentFundEmergencyName  = "development_fund_emergency"
	DevelopmentFundIncentiveName  = "development_fund_incentive"
	DevelopmentFundAuditName      = "development_fund_audit"
	DevelopmentFundQualityName    = "development_fund_quality"
	DevelopmentFundReviewName     = "development_fund_review"
)

// Event types for development fund
const (
	DevEventTypeProposalSubmitted  = "dev_proposal_submitted"
	DevEventTypeProposalApproved   = "dev_proposal_approved"
	DevEventTypeProposalRejected   = "dev_proposal_rejected"
	DevEventTypeProjectStarted     = "project_started"
	DevEventTypePhaseCompleted     = "phase_completed"
	DevEventTypeMilestoneAchieved  = "milestone_achieved"
	DevEventTypeDeliverableSubmitted = "deliverable_submitted"
	DevEventTypeQualityReview      = "quality_review"
	DevEventTypeSecurityAudit      = "security_audit"
	DevEventTypeProjectCompleted   = "project_completed"
	DevEventTypeProjectDeployed    = "project_deployed"
	DevEventTypePaymentReleased    = "payment_released"
	DevEventTypeIncentiveAwarded   = "incentive_awarded"
	DevEventTypeReviewCompleted    = "review_completed"
	DevEventTypeAuditCompleted     = "audit_completed"
)

// Development Fund Allocation - 15% of total supply
const (
	DevelopmentFundPercentage = 15 // 15% of total supply
	DevelopmentFundAllocation = 214294149 // 214,294,149 NAMO tokens
)

// Development Fund Categories with allocation limits
var DevelopmentCategoryLimits = map[DevelopmentCategory]math.LegacyDec{
	CategoryCoreBlockchain:    math.LegacyNewDecWithPrec(25, 2), // 25%
	CategorySmartContracts:    math.LegacyNewDecWithPrec(15, 2), // 15%
	CategoryUserInterface:     math.LegacyNewDecWithPrec(12, 2), // 12%
	CategoryMobileApp:         math.LegacyNewDecWithPrec(10, 2), // 10%
	CategoryWalletDevelopment: math.LegacyNewDecWithPrec(8, 2),  // 8%
	CategorySecurityAudits:    math.LegacyNewDecWithPrec(8, 2),  // 8%
	CategoryPerformanceOpt:    math.LegacyNewDecWithPrec(5, 2),  // 5%
	CategoryTesting:           math.LegacyNewDecWithPrec(5, 2),  // 5%
	CategoryDocumentation:     math.LegacyNewDecWithPrec(3, 2),  // 3%
	CategoryDevTools:          math.LegacyNewDecWithPrec(3, 2),  // 3%
	CategoryInfrastructure:    math.LegacyNewDecWithPrec(2, 2),  // 2%
	CategoryIntegrations:      math.LegacyNewDecWithPrec(2, 2),  // 2%
	CategoryResearch:          math.LegacyNewDecWithPrec(1, 2),  // 1%
	CategoryEmergencyFixes:    math.LegacyNewDecWithPrec(1, 2),  // 1%
}

// Multi-signature configuration for development fund
const (
	DevMultiSigThreshold = 6 // 6 out of 11 signatures required
	DevMultiSigSigners   = 11 // 11 total signers
)

// Development fund governance parameters
const (
	DevMinProposalAmount     = 1000    // 1,000 NAMO minimum
	DevMaxProposalAmount     = 5000000 // 5M NAMO maximum
	DevReviewPeriod          = 14      // 14 days review period
	DevQualityThreshold      = 8       // 8/10 quality score required
	DevTransparencyThreshold = 9       // 9/10 transparency score required
	DevSecurityThreshold     = 9       // 9/10 security score required
	DevTestCoverageThreshold = 90      // 90% test coverage required
	DevDocumentationThreshold = 8      // 8/10 documentation score required
)