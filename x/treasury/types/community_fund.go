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
	"github.com/cosmos/cosmos-sdk/types/address"
	"time"
)

// CommunityFundProposal represents a proposal for using community funds
type CommunityFundProposal struct {
	ProposalID      uint64           `json:"proposal_id"`
	Proposer        sdk.AccAddress   `json:"proposer"`
	Title           string           `json:"title"`
	Description     string           `json:"description"`
	Category        ProposalCategory `json:"category"`
	RequestedAmount sdk.Coin         `json:"requested_amount"`
	Recipients      []Recipient      `json:"recipients"`
	Milestones      []Milestone      `json:"milestones"`
	Status          ProposalStatus   `json:"status"`
	VotingPeriod    time.Duration    `json:"voting_period"`
	SubmissionTime  time.Time        `json:"submission_time"`
	VotingEndTime   time.Time        `json:"voting_end_time"`
	ExecutionTime   *time.Time       `json:"execution_time,omitempty"`
	VotesFor        math.Int         `json:"votes_for"`
	VotesAgainst    math.Int         `json:"votes_against"`
	VotesAbstain    math.Int         `json:"votes_abstain"`
	TotalVotes      math.Int         `json:"total_votes"`
	QuorumReached   bool             `json:"quorum_reached"`
	Passed          bool             `json:"passed"`
	AuditRequired   bool             `json:"audit_required"`
	TransparencyScore uint8          `json:"transparency_score"`
}

// ProposalCategory defines the types of community fund proposals
type ProposalCategory string

const (
	CategoryCommunityRewards    ProposalCategory = "community_rewards"
	CategoryDeveloperIncentives ProposalCategory = "developer_incentives"
	CategoryEducationPrograms   ProposalCategory = "education_programs"
	CategoryMarketingCampaigns  ProposalCategory = "marketing_campaigns"
	CategoryEvents              ProposalCategory = "events"
	CategoryInfrastructure      ProposalCategory = "infrastructure"
	CategoryPartnerships        ProposalCategory = "partnerships"
	CategorySocialImpact        ProposalCategory = "social_impact"
	CategoryResearch            ProposalCategory = "research"
	CategoryEmergencyFund       ProposalCategory = "emergency_fund"
)

// ProposalStatus defines the status of a proposal
type ProposalStatus string

const (
	StatusPending     ProposalStatus = "pending"
	StatusActive      ProposalStatus = "active"
	StatusPassed      ProposalStatus = "passed"
	StatusRejected    ProposalStatus = "rejected"
	StatusExecuting   ProposalStatus = "executing"
	StatusCompleted   ProposalStatus = "completed"
	StatusCancelled   ProposalStatus = "cancelled"
	StatusAuditFailed ProposalStatus = "audit_failed"
)

// Recipient represents a fund recipient
type Recipient struct {
	Address     sdk.AccAddress `json:"address"`
	Amount      sdk.Coin       `json:"amount"`
	Description string         `json:"description"`
	KYCVerified bool           `json:"kyc_verified"`
	Reputation  uint8          `json:"reputation"`
}

// Milestone represents a project milestone
type Milestone struct {
	ID          uint64     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Amount      sdk.Coin   `json:"amount"`
	DueDate     time.Time  `json:"due_date"`
	Completed   bool       `json:"completed"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Evidence    string     `json:"evidence"`
	Approved    bool       `json:"approved"`
}

// CommunityFundBalance tracks the current balance and allocations
type CommunityFundBalance struct {
	TotalBalance      sdk.Coin `json:"total_balance"`
	AllocatedAmount   sdk.Coin `json:"allocated_amount"`
	AvailableAmount   sdk.Coin `json:"available_amount"`
	ReservedAmount    sdk.Coin `json:"reserved_amount"`
	PendingAmount     sdk.Coin `json:"pending_amount"`
	LastUpdateHeight  int64    `json:"last_update_height"`
	LastUpdateTime    time.Time `json:"last_update_time"`
}

// CommunityFundTransaction represents a fund transaction
type CommunityFundTransaction struct {
	TxID          string              `json:"tx_id"`
	ProposalID    uint64              `json:"proposal_id"`
	From          sdk.AccAddress      `json:"from"`
	To            sdk.AccAddress      `json:"to"`
	Amount        sdk.Coin            `json:"amount"`
	Type          TransactionType     `json:"type"`
	Category      ProposalCategory    `json:"category"`
	Description   string              `json:"description"`
	Timestamp     time.Time           `json:"timestamp"`
	BlockHeight   int64               `json:"block_height"`
	Status        TransactionStatus   `json:"status"`
	AuditTrail    []AuditEntry        `json:"audit_trail"`
	Verified      bool                `json:"verified"`
}

// TransactionType defines the type of fund transaction
type TransactionType string

const (
	TxTypeDeposit       TransactionType = "deposit"
	TxTypeWithdrawal    TransactionType = "withdrawal"
	TxTypeAllocation    TransactionType = "allocation"
	TxTypeRefund        TransactionType = "refund"
	TxTypeReward        TransactionType = "reward"
	TxTypeTransfer      TransactionType = "transfer"
	TxTypeBurn          TransactionType = "burn"
	TxTypeStake         TransactionType = "stake"
	TxTypeUnstake       TransactionType = "unstake"
	TxTypeSlash         TransactionType = "slash"
)

// TransactionStatus defines the status of a transaction
type TransactionStatus string

const (
	TxStatusPending   TransactionStatus = "pending"
	TxStatusConfirmed TransactionStatus = "confirmed"
	TxStatusFailed    TransactionStatus = "failed"
	TxStatusCancelled TransactionStatus = "cancelled"
	TxStatusAudited   TransactionStatus = "audited"
)

// AuditEntry represents an audit log entry
type AuditEntry struct {
	Timestamp   time.Time  `json:"timestamp"`
	Action      string     `json:"action"`
	Actor       sdk.AccAddress `json:"actor"`
	Description string     `json:"description"`
	Evidence    string     `json:"evidence"`
	Approved    bool       `json:"approved"`
}

// CommunityGovernance represents governance parameters
type CommunityGovernance struct {
	QuorumPercentage      math.LegacyDec `json:"quorum_percentage"`
	PassingPercentage     math.LegacyDec `json:"passing_percentage"`
	VotingPeriod          time.Duration  `json:"voting_period"`
	MinDeposit            sdk.Coin       `json:"min_deposit"`
	MaxProposalSize       sdk.Coin       `json:"max_proposal_size"`
	RequiredStake         sdk.Coin       `json:"required_stake"`
	AuditThreshold        sdk.Coin       `json:"audit_threshold"`
	TransparencyRequired  bool           `json:"transparency_required"`
	CommunityApprovalReq  bool           `json:"community_approval_required"`
	MultiSigRequired      bool           `json:"multi_sig_required"`
	MinReputation         uint8          `json:"min_reputation"`
}

// TransparencyReport represents a transparency report
type TransparencyReport struct {
	ReportID        uint64            `json:"report_id"`
	ReportingPeriod time.Duration     `json:"reporting_period"`
	StartDate       time.Time         `json:"start_date"`
	EndDate         time.Time         `json:"end_date"`
	TotalFunds      sdk.Coin          `json:"total_funds"`
	AllocatedFunds  sdk.Coin          `json:"allocated_funds"`
	SpentFunds      sdk.Coin          `json:"spent_funds"`
	RemainingFunds  sdk.Coin          `json:"remaining_funds"`
	Categories      []CategoryReport  `json:"categories"`
	TopRecipients   []RecipientReport `json:"top_recipients"`
	Milestones      []MilestoneReport `json:"milestones"`
	ImpactMetrics   ImpactMetrics     `json:"impact_metrics"`
	AuditStatus     AuditStatus       `json:"audit_status"`
	PublicFeedback  []FeedbackEntry   `json:"public_feedback"`
	NextReportDate  time.Time         `json:"next_report_date"`
}

// CategoryReport represents spending by category
type CategoryReport struct {
	Category      ProposalCategory `json:"category"`
	TotalSpent    sdk.Coin         `json:"total_spent"`
	Proposals     uint64           `json:"proposals"`
	Success       uint64           `json:"success"`
	Pending       uint64           `json:"pending"`
	Failed        uint64           `json:"failed"`
	Impact        string           `json:"impact"`
	Efficiency    math.LegacyDec   `json:"efficiency"`
}

// RecipientReport represents top recipients
type RecipientReport struct {
	Address         sdk.AccAddress `json:"address"`
	TotalReceived   sdk.Coin       `json:"total_received"`
	Proposals       uint64         `json:"proposals"`
	CompletedProj   uint64         `json:"completed_projects"`
	SuccessRate     math.LegacyDec `json:"success_rate"`
	Reputation      uint8          `json:"reputation"`
	LastActivity    time.Time      `json:"last_activity"`
}

// MilestoneReport represents milestone completion
type MilestoneReport struct {
	Category          ProposalCategory `json:"category"`
	TotalMilestones   uint64           `json:"total_milestones"`
	CompletedMilestones uint64         `json:"completed_milestones"`
	OnTimeMilestones  uint64           `json:"on_time_milestones"`
	OverdueMilestones uint64           `json:"overdue_milestones"`
	CompletionRate    math.LegacyDec   `json:"completion_rate"`
}

// ImpactMetrics represents impact measurements
type ImpactMetrics struct {
	DevelopersSupported   uint64  `json:"developers_supported"`
	ProjectsCompleted     uint64  `json:"projects_completed"`
	CommunityGrowth       uint64  `json:"community_growth"`
	EducationPrograms     uint64  `json:"education_programs"`
	PartnershipsFormed    uint64  `json:"partnerships_formed"`
	SocialImpactScore     uint8   `json:"social_impact_score"`
	CommunityEngagement   uint8   `json:"community_engagement"`
	InnovationIndex       uint8   `json:"innovation_index"`
	SustainabilityScore   uint8   `json:"sustainability_score"`
}

// AuditStatus represents audit information
type AuditStatus struct {
	LastAuditDate     time.Time         `json:"last_audit_date"`
	NextAuditDate     time.Time         `json:"next_audit_date"`
	AuditFirm         string            `json:"audit_firm"`
	AuditScore        uint8             `json:"audit_score"`
	Recommendations   []string          `json:"recommendations"`
	ComplianceStatus  ComplianceStatus  `json:"compliance_status"`
	Issues            []AuditIssue      `json:"issues"`
	Resolved          uint64            `json:"resolved"`
	Pending           uint64            `json:"pending"`
}

// ComplianceStatus represents compliance status
type ComplianceStatus string

const (
	ComplianceGreen  ComplianceStatus = "green"
	ComplianceYellow ComplianceStatus = "yellow"
	ComplianceRed    ComplianceStatus = "red"
)

// AuditIssue represents an audit issue
type AuditIssue struct {
	ID          uint64    `json:"id"`
	Severity    string    `json:"severity"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Found       time.Time `json:"found"`
	Resolved    *time.Time `json:"resolved,omitempty"`
	Status      string    `json:"status"`
}

// FeedbackEntry represents community feedback
type FeedbackEntry struct {
	ID          uint64         `json:"id"`
	Submitter   sdk.AccAddress `json:"submitter"`
	Content     string         `json:"content"`
	Rating      uint8          `json:"rating"`
	Category    string         `json:"category"`
	Timestamp   time.Time      `json:"timestamp"`
	Verified    bool           `json:"verified"`
	Response    string         `json:"response,omitempty"`
	Addressed   bool           `json:"addressed"`
}

// Storage keys for community fund
var (
	CommunityFundProposalKey        = collections.NewPrefix(100)
	CommunityFundBalanceKey         = collections.NewPrefix(101)
	CommunityFundTransactionKey     = collections.NewPrefix(102)
	CommunityFundGovernanceKey      = collections.NewPrefix(103)
	CommunityFundTransparencyKey    = collections.NewPrefix(104)
	CommunityFundVoteKey            = collections.NewPrefix(105)
	CommunityFundAuditKey           = collections.NewPrefix(106)
	CommunityFundFeedbackKey        = collections.NewPrefix(107)
	CommunityFundStatsKey           = collections.NewPrefix(108)
	CommunityFundReputationKey      = collections.NewPrefix(109)
)

// Module account names for community fund
const (
	CommunityFundModuleName       = "community_fund"
	CommunityFundPoolName         = "community_fund_pool"
	CommunityFundEscrowName       = "community_fund_escrow"
	CommunityFundRewardsName      = "community_fund_rewards"
	CommunityFundReserveName      = "community_fund_reserve"
	CommunityFundAuditName        = "community_fund_audit"
	CommunityFundGovernanceName   = "community_fund_governance"
	CommunityFundTransparencyName = "community_fund_transparency"
)

// Event types for community fund
const (
	EventTypeProposalSubmitted     = "proposal_submitted"
	EventTypeProposalVoted         = "proposal_voted"
	EventTypeProposalPassed        = "proposal_passed"
	EventTypeProposalRejected      = "proposal_rejected"
	EventTypeProposalExecuted      = "proposal_executed"
	EventTypeFundAllocated         = "fund_allocated"
	EventTypeFundWithdrawn         = "fund_withdrawn"
	EventTypeFundTransferred       = "fund_transferred"
	EventTypeMilestoneCompleted    = "milestone_completed"
	EventTypeTransparencyReport    = "transparency_report"
	EventTypeAuditCompleted        = "audit_completed"
	EventTypeFeedbackSubmitted     = "feedback_submitted"
	EventTypeReputationUpdated     = "reputation_updated"
)

// Default governance parameters
var DefaultCommunityGovernance = CommunityGovernance{
	QuorumPercentage:      math.LegacyNewDecWithPrec(33, 2), // 33%
	PassingPercentage:     math.LegacyNewDecWithPrec(51, 2), // 51%
	VotingPeriod:          time.Hour * 24 * 7,               // 7 days
	MinDeposit:            sdk.NewCoin("namo", math.NewInt(1000)),
	MaxProposalSize:       sdk.NewCoin("namo", math.NewInt(10000000)), // 10M NAMO
	RequiredStake:         sdk.NewCoin("namo", math.NewInt(10000)),
	AuditThreshold:        sdk.NewCoin("namo", math.NewInt(1000000)), // 1M NAMO
	TransparencyRequired:  true,
	CommunityApprovalReq:  true,
	MultiSigRequired:      true,
	MinReputation:         7, // Out of 10
}

// CommunityFundAllocation represents the 15% allocation
const (
	CommunityFundPercentage = 15 // 15% of total supply
	CommunityFundAllocation = 214294149 // 214,294,149 NAMO tokens
)

// Community Fund Categories with allocation limits
var CategoryLimits = map[ProposalCategory]math.LegacyDec{
	CategoryCommunityRewards:    math.LegacyNewDecWithPrec(30, 2), // 30%
	CategoryDeveloperIncentives: math.LegacyNewDecWithPrec(25, 2), // 25%
	CategoryEducationPrograms:   math.LegacyNewDecWithPrec(15, 2), // 15%
	CategoryMarketingCampaigns:  math.LegacyNewDecWithPrec(10, 2), // 10%
	CategoryEvents:              math.LegacyNewDecWithPrec(5, 2),  // 5%
	CategoryInfrastructure:      math.LegacyNewDecWithPrec(5, 2),  // 5%
	CategoryPartnerships:        math.LegacyNewDecWithPrec(3, 2),  // 3%
	CategorySocialImpact:        math.LegacyNewDecWithPrec(3, 2),  // 3%
	CategoryResearch:            math.LegacyNewDecWithPrec(2, 2),  // 2%
	CategoryEmergencyFund:       math.LegacyNewDecWithPrec(2, 2),  // 2%
}

// Multi-signature configuration
const (
	MultiSigThreshold = 5 // 5 out of 9 signatures required
	MultiSigSigners   = 9 // 9 total signers
)

// Multi-signature signer roles
const (
	SignerRoleFounder       = "founder"
	SignerRoleCommunityLead = "community_lead"
	SignerRoleTechnicalLead = "technical_lead"
	SignerRoleFinancialLead = "financial_lead"
	SignerRoleAuditor       = "auditor"
	SignerRoleValidator     = "validator"
	SignerRoleDeveloper     = "developer"
	SignerRoleAdviser       = "adviser"
	SignerRoleRepresentative = "representative"
)

// Transparency requirements
const (
	TransparencyScoreThreshold = 8 // Out of 10
	MonthlyReportRequired     = true
	QuarterlyAuditRequired    = true
	PublicFeedbackRequired    = true
	RealTimeTrackingRequired  = true
)