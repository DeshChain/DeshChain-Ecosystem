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

// MultiSigGovernance represents the multi-signature governance system
type MultiSigGovernance struct {
	ID                uint64              `json:"id"`
	Name              string              `json:"name"`
	Description       string              `json:"description"`
	Type              GovernanceType      `json:"type"`
	Threshold         uint8               `json:"threshold"`
	TotalSigners      uint8               `json:"total_signers"`
	Signers           []Signer            `json:"signers"`
	Proposals         []GovernanceProposal `json:"proposals"`
	ActiveProposals   uint64              `json:"active_proposals"`
	CompletedProposals uint64             `json:"completed_proposals"`
	CreatedAt         time.Time           `json:"created_at"`
	UpdatedAt         time.Time           `json:"updated_at"`
	Status            GovernanceStatus    `json:"status"`
	Permissions       []Permission        `json:"permissions"`
	Rules             GovernanceRules     `json:"rules"`
}

// GovernanceType defines the type of governance
type GovernanceType string

const (
	GovernanceTypeCommunityFund   GovernanceType = "community_fund"
	GovernanceTypeDevelopmentFund GovernanceType = "development_fund"
	GovernanceTypeEmergencyFund   GovernanceType = "emergency_fund"
	GovernanceTypeFounderRoyalty  GovernanceType = "founder_royalty"
	GovernanceTypeNGOFund         GovernanceType = "ngo_fund"
	GovernanceTypeGeneral         GovernanceType = "general"
	GovernanceTypeSecurityCouncil GovernanceType = "security_council"
	GovernanceTypeTechnicalCouncil GovernanceType = "technical_council"
)

// GovernanceStatus defines the status of governance
type GovernanceStatus string

const (
	GovernanceStatusActive   GovernanceStatus = "active"
	GovernanceStatusInactive GovernanceStatus = "inactive"
	GovernanceStatusSuspended GovernanceStatus = "suspended"
	GovernanceStatusUpdating GovernanceStatus = "updating"
)

// Signer represents a multi-signature signer
type Signer struct {
	Address         sdk.AccAddress  `json:"address"`
	Role            SignerRole      `json:"role"`
	Weight          uint8           `json:"weight"`
	Status          SignerStatus    `json:"status"`
	Reputation      uint8           `json:"reputation"`
	JoinedAt        time.Time       `json:"joined_at"`
	LastActivity    time.Time       `json:"last_activity"`
	SignedProposals uint64          `json:"signed_proposals"`
	TotalProposals  uint64          `json:"total_proposals"`
	SigningRate     math.LegacyDec  `json:"signing_rate"`
	Expertise       []string        `json:"expertise"`
	Permissions     []Permission    `json:"permissions"`
	PublicKey       string          `json:"public_key"`
	Contact         ContactInfo     `json:"contact"`
	KYCVerified     bool            `json:"kyc_verified"`
	BackgroundCheck bool            `json:"background_check"`
	Insurance       bool            `json:"insurance"`
	Bond            sdk.Coin        `json:"bond"`
}

// SignerRole defines the role of a signer
type SignerRole string

const (
	SignerRoleFounder           SignerRole = "founder"
	SignerRoleCommunityLead     SignerRole = "community_lead"
	SignerRoleTechnicalLead     SignerRole = "technical_lead"
	SignerRoleFinancialLead     SignerRole = "financial_lead"
	SignerRoleSecurityExpert    SignerRole = "security_expert"
	SignerRoleAuditor           SignerRole = "auditor"
	SignerRoleValidator         SignerRole = "validator"
	SignerRoleDeveloper         SignerRole = "developer"
	SignerRoleAdviser           SignerRole = "adviser"
	SignerRoleRepresentative    SignerRole = "representative"
	SignerRoleLegal             SignerRole = "legal"
	SignerRoleCompliance        SignerRole = "compliance"
	SignerRoleRiskManagement    SignerRole = "risk_management"
	SignerRoleProductManager    SignerRole = "product_manager"
	SignerRoleDesignLead        SignerRole = "design_lead"
	SignerRoleQualityAssurance  SignerRole = "quality_assurance"
	SignerRoleDataAnalyst       SignerRole = "data_analyst"
	SignerRoleMarketingLead     SignerRole = "marketing_lead"
	SignerRolePartnershipLead   SignerRole = "partnership_lead"
	SignerRoleGovernanceCoordinator SignerRole = "governance_coordinator"
)

// SignerStatus defines the status of a signer
type SignerStatus string

const (
	SignerStatusActive    SignerStatus = "active"
	SignerStatusInactive  SignerStatus = "inactive"
	SignerStatusSuspended SignerStatus = "suspended"
	SignerStatusPending   SignerStatus = "pending"
	SignerStatusRevoked   SignerStatus = "revoked"
	SignerStatusOnLeave   SignerStatus = "on_leave"
)

// ContactInfo represents contact information
type ContactInfo struct {
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Telegram        string `json:"telegram"`
	Discord         string `json:"discord"`
	LinkedIn        string `json:"linkedin"`
	Twitter         string `json:"twitter"`
	GitHub          string `json:"github"`
	PGPKey          string `json:"pgp_key"`
	TimeZone        string `json:"timezone"`
	PreferredComm   string `json:"preferred_communication"`
	EmergencyContact string `json:"emergency_contact"`
}

// Permission represents a permission
type Permission string

const (
	PermissionPropose       Permission = "propose"
	PermissionVote          Permission = "vote"
	PermissionExecute       Permission = "execute"
	PermissionEmergencyStop Permission = "emergency_stop"
	PermissionAudit         Permission = "audit"
	PermissionReview        Permission = "review"
	PermissionApprove       Permission = "approve"
	PermissionReject        Permission = "reject"
	PermissionModify        Permission = "modify"
	PermissionCancel        Permission = "cancel"
	PermissionDelegate      Permission = "delegate"
	PermissionViewAll       Permission = "view_all"
	PermissionManageSigners Permission = "manage_signers"
	PermissionManageRules   Permission = "manage_rules"
	PermissionFinancialOp   Permission = "financial_operations"
	PermissionTechnicalOp   Permission = "technical_operations"
	PermissionSecurityOp    Permission = "security_operations"
	PermissionEmergencyOp   Permission = "emergency_operations"
	PermissionGovernanceOp  Permission = "governance_operations"
)

// GovernanceProposal represents a governance proposal
type GovernanceProposal struct {
	ID                uint64              `json:"id"`
	Title             string              `json:"title"`
	Description       string              `json:"description"`
	Type              ProposalType        `json:"type"`
	Category          string              `json:"category"`
	Priority          PriorityLevel       `json:"priority"`
	Proposer          sdk.AccAddress      `json:"proposer"`
	RequestedAmount   sdk.Coin            `json:"requested_amount"`
	Recipients        []ProposalRecipient `json:"recipients"`
	Attachments       []Attachment        `json:"attachments"`
	Status            ProposalStatus      `json:"status"`
	SubmissionTime    time.Time           `json:"submission_time"`
	ReviewPeriod      time.Duration       `json:"review_period"`
	VotingPeriod      time.Duration       `json:"voting_period"`
	ExecutionPeriod   time.Duration       `json:"execution_period"`
	ReviewEndTime     time.Time           `json:"review_end_time"`
	VotingEndTime     time.Time           `json:"voting_end_time"`
	ExecutionEndTime  time.Time           `json:"execution_end_time"`
	Signatures        []Signature         `json:"signatures"`
	RequiredSignatures uint8              `json:"required_signatures"`
	CurrentSignatures  uint8              `json:"current_signatures"`
	Approved          bool                `json:"approved"`
	Executed          bool                `json:"executed"`
	ExecutedAt        *time.Time          `json:"executed_at,omitempty"`
	CancelledAt       *time.Time          `json:"cancelled_at,omitempty"`
	Reason            string              `json:"reason,omitempty"`
	Impact            ImpactAssessment    `json:"impact"`
	Risk              RiskAssessment      `json:"risk"`
	Compliance        ComplianceCheck     `json:"compliance"`
	Audit             AuditInfo           `json:"audit"`
	Feedback          []ProposalFeedback  `json:"feedback"`
	Metrics           ProposalMetrics     `json:"metrics"`
	Transparency      TransparencyInfo    `json:"transparency"`
}

// ProposalType defines the type of proposal
type ProposalType string

const (
	ProposalTypeFundAllocation    ProposalType = "fund_allocation"
	ProposalTypeParameterChange   ProposalType = "parameter_change"
	ProposalTypeSignerAddition    ProposalType = "signer_addition"
	ProposalTypeSignerRemoval     ProposalType = "signer_removal"
	ProposalTypeThresholdChange   ProposalType = "threshold_change"
	ProposalTypeRuleChange        ProposalType = "rule_change"
	ProposalTypeEmergencyAction   ProposalType = "emergency_action"
	ProposalTypeAuditRequest      ProposalType = "audit_request"
	ProposalTypeSecurityUpdate    ProposalType = "security_update"
	ProposalTypeSystemUpgrade     ProposalType = "system_upgrade"
	ProposalTypePartnership       ProposalType = "partnership"
	ProposalTypeCompliance        ProposalType = "compliance"
	ProposalTypeGovernanceUpdate  ProposalType = "governance_update"
	ProposalTypeGeneral           ProposalType = "general"
)

// ProposalRecipient represents a proposal recipient
type ProposalRecipient struct {
	Address     sdk.AccAddress `json:"address"`
	Amount      sdk.Coin       `json:"amount"`
	Purpose     string         `json:"purpose"`
	Milestones  []string       `json:"milestones"`
	KYCVerified bool           `json:"kyc_verified"`
	Reputation  uint8          `json:"reputation"`
	Risk        uint8          `json:"risk"`
}

// Attachment represents a proposal attachment
type Attachment struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Size        uint64    `json:"size"`
	Hash        string    `json:"hash"`
	URL         string    `json:"url"`
	Description string    `json:"description"`
	Uploader    sdk.AccAddress `json:"uploader"`
	UploadTime  time.Time `json:"upload_time"`
	Verified    bool      `json:"verified"`
}

// Signature represents a signature
type Signature struct {
	Signer      sdk.AccAddress `json:"signer"`
	Signature   string         `json:"signature"`
	SignedAt    time.Time      `json:"signed_at"`
	Weight      uint8          `json:"weight"`
	Status      SignatureStatus `json:"status"`
	Comments    string         `json:"comments"`
	Conditions  []string       `json:"conditions"`
	Delegated   bool           `json:"delegated"`
	DelegatedBy sdk.AccAddress `json:"delegated_by,omitempty"`
}

// SignatureStatus defines the status of a signature
type SignatureStatus string

const (
	SignatureStatusSigned   SignatureStatus = "signed"
	SignatureStatusPending  SignatureStatus = "pending"
	SignatureStatusRevoked  SignatureStatus = "revoked"
	SignatureStatusExpired  SignatureStatus = "expired"
	SignatureStatusInvalid  SignatureStatus = "invalid"
	SignatureStatusConditional SignatureStatus = "conditional"
)

// ComplianceCheck represents compliance check results
type ComplianceCheck struct {
	Passed            bool      `json:"passed"`
	Score             uint8     `json:"score"`
	Checks            []string  `json:"checks"`
	Issues            []string  `json:"issues"`
	Recommendations   []string  `json:"recommendations"`
	Reviewer          sdk.AccAddress `json:"reviewer"`
	ReviewDate        time.Time `json:"review_date"`
	NextReviewDate    time.Time `json:"next_review_date"`
	ComplianceLevel   string    `json:"compliance_level"`
	RegulatoryStatus  string    `json:"regulatory_status"`
	AuditRequired     bool      `json:"audit_required"`
	DocumentationReq  bool      `json:"documentation_required"`
	ApprovalRequired  bool      `json:"approval_required"`
}

// AuditInfo represents audit information
type AuditInfo struct {
	Required      bool           `json:"required"`
	Type          string         `json:"type"`
	Scope         string         `json:"scope"`
	Auditor       sdk.AccAddress `json:"auditor"`
	ScheduledDate time.Time      `json:"scheduled_date"`
	CompletedDate *time.Time     `json:"completed_date,omitempty"`
	Status        string         `json:"status"`
	Report        string         `json:"report"`
	Score         uint8          `json:"score"`
	Issues        []string       `json:"issues"`
	Recommendations []string     `json:"recommendations"`
	Cost          sdk.Coin       `json:"cost"`
	NextAuditDate time.Time      `json:"next_audit_date"`
}

// ProposalFeedback represents feedback on a proposal
type ProposalFeedback struct {
	ID          uint64         `json:"id"`
	Submitter   sdk.AccAddress `json:"submitter"`
	Type        string         `json:"type"`
	Content     string         `json:"content"`
	Rating      uint8          `json:"rating"`
	Category    string         `json:"category"`
	Timestamp   time.Time      `json:"timestamp"`
	Verified    bool           `json:"verified"`
	Helpful     uint64         `json:"helpful"`
	Response    string         `json:"response,omitempty"`
	Addressed   bool           `json:"addressed"`
	Public      bool           `json:"public"`
	Anonymous   bool           `json:"anonymous"`
}

// ProposalMetrics represents proposal metrics
type ProposalMetrics struct {
	Views            uint64        `json:"views"`
	Comments         uint64        `json:"comments"`
	Shares           uint64        `json:"shares"`
	Likes            uint64        `json:"likes"`
	Dislikes         uint64        `json:"dislikes"`
	EngagementRate   math.LegacyDec `json:"engagement_rate"`
	SentimentScore   uint8         `json:"sentiment_score"`
	CommunitySupport uint8         `json:"community_support"`
	ExpertEndorsement uint8        `json:"expert_endorsement"`
	MediaCoverage    uint8         `json:"media_coverage"`
	SocialImpact     uint8         `json:"social_impact"`
	ViralityScore    uint8         `json:"virality_score"`
}

// TransparencyInfo represents transparency information
type TransparencyInfo struct {
	Level             uint8     `json:"level"`
	PublicDisclosure  bool      `json:"public_disclosure"`
	DocumentsPublic   bool      `json:"documents_public"`
	VotingPublic      bool      `json:"voting_public"`
	FinancialsPublic  bool      `json:"financials_public"`
	ProgressPublic    bool      `json:"progress_public"`
	CommunityInput    bool      `json:"community_input"`
	ExternalReview    bool      `json:"external_review"`
	IndependentAudit  bool      `json:"independent_audit"`
	RealTimeTracking  bool      `json:"real_time_tracking"`
	OpenSource        bool      `json:"open_source"`
	PublicMeetings    bool      `json:"public_meetings"`
	RegularReports    bool      `json:"regular_reports"`
	TransparencyScore uint8     `json:"transparency_score"`
	LastUpdated       time.Time `json:"last_updated"`
}

// GovernanceRules represents governance rules
type GovernanceRules struct {
	MinProposalDeposit     sdk.Coin       `json:"min_proposal_deposit"`
	MaxProposalSize        sdk.Coin       `json:"max_proposal_size"`
	ReviewPeriod           time.Duration  `json:"review_period"`
	VotingPeriod           time.Duration  `json:"voting_period"`
	ExecutionPeriod        time.Duration  `json:"execution_period"`
	QuorumThreshold        uint8          `json:"quorum_threshold"`
	PassingThreshold       uint8          `json:"passing_threshold"`
	VetoThreshold          uint8          `json:"veto_threshold"`
	MaxActiveProposals     uint64         `json:"max_active_proposals"`
	CooldownPeriod         time.Duration  `json:"cooldown_period"`
	SignerTerm             time.Duration  `json:"signer_term"`
	MaxSignerInactivity    time.Duration  `json:"max_signer_inactivity"`
	RequiredExpertise      []string       `json:"required_expertise"`
	MinReputation          uint8          `json:"min_reputation"`
	MaxRiskTolerance       uint8          `json:"max_risk_tolerance"`
	AuditThreshold         sdk.Coin       `json:"audit_threshold"`
	TransparencyRequired   bool           `json:"transparency_required"`
	CommunityInputRequired bool           `json:"community_input_required"`
	ExternalReviewRequired bool           `json:"external_review_required"`
	EmergencyPowers        []Permission   `json:"emergency_powers"`
	RestrictedActions      []string       `json:"restricted_actions"`
	DelegationAllowed      bool           `json:"delegation_allowed"`
	ProxyVotingAllowed     bool           `json:"proxy_voting_allowed"`
	Anonymous              bool           `json:"anonymous_voting"`
	WeightedVoting         bool           `json:"weighted_voting"`
	TimeBasedVoting        bool           `json:"time_based_voting"`
	LocationBasedVoting    bool           `json:"location_based_voting"`
	StakeBasedVoting       bool           `json:"stake_based_voting"`
	ReputationBasedVoting  bool           `json:"reputation_based_voting"`
}

// DashboardMetrics represents dashboard metrics
type DashboardMetrics struct {
	TotalProposals        uint64         `json:"total_proposals"`
	ActiveProposals       uint64         `json:"active_proposals"`
	CompletedProposals    uint64         `json:"completed_proposals"`
	RejectedProposals     uint64         `json:"rejected_proposals"`
	TotalFundsAllocated   sdk.Coin       `json:"total_funds_allocated"`
	TotalFundsSpent       sdk.Coin       `json:"total_funds_spent"`
	TotalFundsRemaining   sdk.Coin       `json:"total_funds_remaining"`
	AverageProposalValue  sdk.Coin       `json:"average_proposal_value"`
	AverageCompletionTime time.Duration  `json:"average_completion_time"`
	SuccessRate           math.LegacyDec `json:"success_rate"`
	CommunityEngagement   uint8          `json:"community_engagement"`
	TransparencyScore     uint8          `json:"transparency_score"`
	ComplianceScore       uint8          `json:"compliance_score"`
	SecurityScore         uint8          `json:"security_score"`
	EfficiencyScore       uint8          `json:"efficiency_score"`
	ImpactScore           uint8          `json:"impact_score"`
	SustainabilityScore   uint8          `json:"sustainability_score"`
	ActiveSigners         uint64         `json:"active_signers"`
	AverageSigningTime    time.Duration  `json:"average_signing_time"`
	SigningRate           math.LegacyDec `json:"signing_rate"`
	LastUpdated           time.Time      `json:"last_updated"`
}

// Storage keys for multi-signature governance
var (
	MultiSigGovernanceKey        = collections.NewPrefix(300)
	MultiSigSignerKey            = collections.NewPrefix(301)
	MultiSigProposalKey          = collections.NewPrefix(302)
	MultiSigSignatureKey         = collections.NewPrefix(303)
	MultiSigRulesKey             = collections.NewPrefix(304)
	MultiSigMetricsKey           = collections.NewPrefix(305)
	MultiSigAuditKey             = collections.NewPrefix(306)
	MultiSigFeedbackKey          = collections.NewPrefix(307)
	MultiSigTransparencyKey      = collections.NewPrefix(308)
	MultiSigComplianceKey        = collections.NewPrefix(309)
)

// Module account names for multi-signature governance
const (
	MultiSigGovernanceModuleName = "multisig_governance"
	MultiSigEscrowName           = "multisig_escrow"
	MultiSigBondName             = "multisig_bond"
	MultiSigRewardsName          = "multisig_rewards"
	MultiSigPenaltyName          = "multisig_penalty"
	MultiSigAuditName            = "multisig_audit"
	MultiSigComplianceName       = "multisig_compliance"
	MultiSigTransparencyName     = "multisig_transparency"
)

// Event types for multi-signature governance
const (
	MultiSigEventTypeSignerAdded     = "signer_added"
	MultiSigEventTypeSignerRemoved   = "signer_removed"
	MultiSigEventTypeSignerUpdated   = "signer_updated"
	MultiSigEventTypeProposalCreated = "proposal_created"
	MultiSigEventTypeProposalSigned  = "proposal_signed"
	MultiSigEventTypeProposalApproved = "proposal_approved"
	MultiSigEventTypeProposalExecuted = "proposal_executed"
	MultiSigEventTypeProposalRejected = "proposal_rejected"
	MultiSigEventTypeProposalCancelled = "proposal_cancelled"
	MultiSigEventTypeThresholdChanged = "threshold_changed"
	MultiSigEventTypeRulesUpdated    = "rules_updated"
	MultiSigEventTypeEmergencyAction = "emergency_action"
	MultiSigEventTypeAuditCompleted  = "audit_completed"
	MultiSigEventTypeComplianceCheck = "compliance_check"
	MultiSigEventTypeTransparencyUpdate = "transparency_update"
)

// Default multi-signature configurations
var (
	DefaultCommunityFundMultiSig = MultiSigGovernance{
		Name:         "Community Fund Governance",
		Description:  "Multi-signature governance for community fund allocation",
		Type:         GovernanceTypeCommunityFund,
		Threshold:    5,
		TotalSigners: 9,
		Status:       GovernanceStatusActive,
		Rules: GovernanceRules{
			MinProposalDeposit:     sdk.NewCoin("namo", math.NewInt(1000)),
			MaxProposalSize:        sdk.NewCoin("namo", math.NewInt(10000000)),
			ReviewPeriod:           time.Hour * 24 * 7,
			VotingPeriod:           time.Hour * 24 * 14,
			ExecutionPeriod:        time.Hour * 24 * 7,
			QuorumThreshold:        60,
			PassingThreshold:       67,
			VetoThreshold:          33,
			MaxActiveProposals:     10,
			CooldownPeriod:         time.Hour * 24,
			SignerTerm:             time.Hour * 24 * 365,
			MaxSignerInactivity:    time.Hour * 24 * 30,
			MinReputation:          7,
			MaxRiskTolerance:       7,
			AuditThreshold:         sdk.NewCoin("namo", math.NewInt(1000000)),
			TransparencyRequired:   true,
			CommunityInputRequired: true,
			ExternalReviewRequired: true,
			WeightedVoting:         true,
			ReputationBasedVoting:  true,
		},
	}

	DefaultDevelopmentFundMultiSig = MultiSigGovernance{
		Name:         "Development Fund Governance",
		Description:  "Multi-signature governance for development fund allocation",
		Type:         GovernanceTypeDevelopmentFund,
		Threshold:    6,
		TotalSigners: 11,
		Status:       GovernanceStatusActive,
		Rules: GovernanceRules{
			MinProposalDeposit:     sdk.NewCoin("namo", math.NewInt(5000)),
			MaxProposalSize:        sdk.NewCoin("namo", math.NewInt(5000000)),
			ReviewPeriod:           time.Hour * 24 * 14,
			VotingPeriod:           time.Hour * 24 * 21,
			ExecutionPeriod:        time.Hour * 24 * 14,
			QuorumThreshold:        70,
			PassingThreshold:       75,
			VetoThreshold:          25,
			MaxActiveProposals:     5,
			CooldownPeriod:         time.Hour * 24 * 3,
			SignerTerm:             time.Hour * 24 * 365 * 2,
			MaxSignerInactivity:    time.Hour * 24 * 21,
			MinReputation:          8,
			MaxRiskTolerance:       6,
			AuditThreshold:         sdk.NewCoin("namo", math.NewInt(2000000)),
			TransparencyRequired:   true,
			CommunityInputRequired: true,
			ExternalReviewRequired: true,
			WeightedVoting:         true,
			ReputationBasedVoting:  true,
		},
	}
)

// Transparency levels
const (
	TransparencyLevelNone       = 0
	TransparencyLevelBasic      = 3
	TransparencyLevelStandard   = 5
	TransparencyLevelHigh       = 7
	TransparencyLevelMaximum    = 10
)

// Compliance levels
const (
	ComplianceLevelBasic      = "basic"
	ComplianceLevelStandard   = "standard"
	ComplianceLevelHigh       = "high"
	ComplianceLevelMaximum    = "maximum"
	ComplianceLevelRegulatory = "regulatory"
)

// Risk levels
const (
	RiskLevelVeryLow  = 1
	RiskLevelLow      = 3
	RiskLevelMedium   = 5
	RiskLevelHigh     = 7
	RiskLevelVeryHigh = 9
	RiskLevelCritical = 10
)

// Reputation levels
const (
	ReputationLevelPoor      = 1
	ReputationLevelFair      = 3
	ReputationLevelGood      = 5
	ReputationLevelVeryGood  = 7
	ReputationLevelExcellent = 9
	ReputationLevelPerfect   = 10
)