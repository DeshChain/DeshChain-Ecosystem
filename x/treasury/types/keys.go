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
)

const (
	// ModuleName defines the module name
	ModuleName = "treasury"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_treasury"

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName
)

// KeyPrefix returns the store key prefix
func KeyPrefix(p string) []byte {
	return []byte(p)
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

// Storage keys for community proposal system
var (
	CommunityProposalSystemKey      = collections.NewPrefix(400)
	ProposalPhaseConfigKey          = collections.NewPrefix(401)
	ProposalTypeConfigKey           = collections.NewPrefix(402)
	GovernanceParamsKey             = collections.NewPrefix(403)
	FounderControlsKey              = collections.NewPrefix(404)
	CommunityPowersKey              = collections.NewPrefix(405)
	TransitionScheduleKey           = collections.NewPrefix(406)
	EmergencyConfigKey              = collections.NewPrefix(407)
	ProposalStatsKey                = collections.NewPrefix(408)
	PhaseTransitionKey              = collections.NewPrefix(409)
)

// Module account names
const (
	// Community Fund module accounts
	CommunityFundModuleName       = "community_fund"
	CommunityFundPoolName         = "community_fund_pool"
	CommunityFundEscrowName       = "community_fund_escrow"
	CommunityFundRewardsName      = "community_fund_rewards"
	CommunityFundReserveName      = "community_fund_reserve"
	CommunityFundAuditName        = "community_fund_audit"
	CommunityFundGovernanceName   = "community_fund_governance"
	CommunityFundTransparencyName = "community_fund_transparency"
	
	// Development Fund module accounts
	DevelopmentFundModuleName     = "development_fund"
	DevelopmentFundPoolName       = "development_fund_pool"
	DevelopmentFundEscrowName     = "development_fund_escrow"
	DevelopmentFundEmergencyName  = "development_fund_emergency"
	DevelopmentFundIncentiveName  = "development_fund_incentive"
	DevelopmentFundAuditName      = "development_fund_audit"
	DevelopmentFundQualityName    = "development_fund_quality"
	DevelopmentFundReviewName     = "development_fund_review"
	
	// Multi-signature governance module accounts
	MultiSigGovernanceModuleName = "multisig_governance"
	MultiSigEscrowName           = "multisig_escrow"
	MultiSigBondName             = "multisig_bond"
	MultiSigRewardsName          = "multisig_rewards"
	MultiSigPenaltyName          = "multisig_penalty"
	MultiSigAuditModuleName      = "multisig_audit"
	MultiSigComplianceName       = "multisig_compliance"
	MultiSigTransparencyModuleName = "multisig_transparency"
	
	// Proposal system module accounts
	ProposalSystemModuleName    = "proposal_system"
	ProposalEscrowName          = "proposal_escrow"
	ProposalRewardsName         = "proposal_rewards"
	ProposalPenaltyName         = "proposal_penalty"
	ProposalEmergencyName       = "proposal_emergency"
	ProposalTransitionName      = "proposal_transition"
)

// Event types for treasury module
const (
	// Community Fund events
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
	
	// Development Fund events
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
	DevEventTypeAuditFinished      = "audit_finished"
	
	// Multi-signature governance events
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
	MultiSigEventTypeAuditRequested  = "audit_requested"
	MultiSigEventTypeComplianceCheck = "compliance_check"
	MultiSigEventTypeTransparencyUpdate = "transparency_update"
	
	// Proposal system events
	EventTypeSystemInitialized      = "system_initialized"
	EventTypePhaseTransition        = "phase_transition"
	EventTypeProposalCreated        = "proposal_created"
	EventTypeProposalApproved       = "proposal_approved"
	EventTypeProposalRejected       = "proposal_rejected"
	EventTypeProposalVetoed         = "proposal_vetoed"
	EventTypeEmergencyActivated     = "emergency_activated"
	EventTypeEmergencyDeactivated   = "emergency_deactivated"
	EventTypePowerTransferred       = "power_transferred"
	EventTypeCheckpointReached      = "checkpoint_reached"
	EventTypeRollbackInitiated      = "rollback_initiated"
	EventTypeGovernanceUpdated      = "governance_updated"
)

// Query endpoints supported by the treasury querier
const (
	QueryCommunityFundBalance      = "community-fund-balance"
	QueryDevelopmentFundBalance    = "development-fund-balance"
	QueryCommunityProposal         = "community-proposal"
	QueryDevelopmentProposal       = "development-proposal"
	QueryMultiSigGovernance        = "multisig-governance"
	QueryProposalSystem            = "proposal-system"
	QueryDashboard                 = "dashboard"
	QueryTransactions              = "transactions"
	QueryGovernancePhase           = "governance-phase"
	QuerySigners                   = "signers"
	QueryProposalHistory           = "proposal-history"
	QueryTransparencyReport        = "transparency-report"
	QueryAuditHistory              = "audit-history"
	QueryMetrics                   = "metrics"
)

// Attribute keys used in treasury events
const (
	AttributeKeyProposalID      = "proposal_id"
	AttributeKeyProposer        = "proposer"
	AttributeKeyCategory        = "category"
	AttributeKeyAmount          = "amount"
	AttributeKeyRecipient       = "recipient"
	AttributeKeyMilestone       = "milestone"
	AttributeKeyPhase           = "phase"
	AttributeKeyStatus          = "status"
	AttributeKeyTransactionID   = "transaction_id"
	AttributeKeySigner          = "signer"
	AttributeKeyThreshold       = "threshold"
	AttributeKeyFromPhase       = "from_phase"
	AttributeKeyToPhase         = "to_phase"
	AttributeKeyTransitionDate  = "transition_date"
	AttributeKeyEmergencyType   = "emergency_type"
	AttributeKeyExecutor        = "executor"
	AttributeKeyReason          = "reason"
)

// Proposal status constants
const (
	StatusDraft        = "draft"
	StatusPending      = "pending"
	StatusActive       = "active"
	StatusPassed       = "passed"
	StatusRejected     = "rejected"
	StatusExecuting    = "executing"
	StatusCompleted    = "completed"
	StatusCancelled    = "cancelled"
	StatusFailed       = "failed"
	StatusExpired      = "expired"
)

// Transaction types
const (
	TxTypeDeposit       = "deposit"
	TxTypeWithdrawal    = "withdrawal"
	TxTypeAllocation    = "allocation"
	TxTypeTransfer      = "transfer"
	TxTypeRefund        = "refund"
	TxTypePenalty       = "penalty"
	TxTypeReward        = "reward"
	TxTypeEmergency     = "emergency"
)

// Default values
const (
	DefaultQuorumPercentage  = 33  // 33%
	DefaultPassingPercentage = 51  // 51%
	DefaultVotingPeriodDays  = 7   // 7 days
	DefaultReviewPeriodDays  = 14  // 14 days
	DefaultMultiSigThreshold = 5   // 5 out of 9
	DefaultTransparencyScore = 10  // Maximum score
)

// Fund allocation amounts (in NAMO)
const (
	CommunityFundAllocation  = 214294149  // 15% of total supply
	DevelopmentFundAllocation = 214294149 // 15% of total supply
	TotalTreasuryAllocation  = 428588298  // 30% of total supply
)