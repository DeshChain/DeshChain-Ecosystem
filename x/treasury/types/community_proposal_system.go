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

// CommunityProposalSystem represents the phased community proposal system
type CommunityProposalSystem struct {
	ID                    uint64                      `json:"id"`
	Name                  string                      `json:"name"`
	Description           string                      `json:"description"`
	LaunchDate            time.Time                   `json:"launch_date"`
	ActivationDate        time.Time                   `json:"activation_date"`       // 3 years after launch
	CurrentPhase          ProposalSystemPhase         `json:"current_phase"`
	PhaseHistory          []PhaseTransition           `json:"phase_history"`
	ProposalTypes         []ProposalTypeConfig        `json:"proposal_types"`
	GovernanceParameters  PhaseGovernanceParams       `json:"governance_parameters"`
	FounderControls       FounderControlConfig        `json:"founder_controls"`
	CommunityPowers       CommunityPowerConfig        `json:"community_powers"`
	TransitionSchedule    TransitionSchedule          `json:"transition_schedule"`
	EmergencyPowers       EmergencyPowerConfig        `json:"emergency_powers"`
	Stats                 ProposalSystemStats         `json:"stats"`
	Status                ProposalSystemStatus        `json:"status"`
}

// ProposalSystemPhase represents the current phase of the system
type ProposalSystemPhase string

const (
	PhaseFounderDriven      ProposalSystemPhase = "founder_driven"      // Years 0-3
	PhaseTransitional       ProposalSystemPhase = "transitional"        // Year 3-4
	PhaseCommunityProposal  ProposalSystemPhase = "community_proposal"  // Year 4-7
	PhaseFullGovernance     ProposalSystemPhase = "full_governance"     // Year 7+
	PhaseEmergency          ProposalSystemPhase = "emergency"           // Special phase
)

// ProposalSystemStatus represents the status of the system
type ProposalSystemStatus string

const (
	SystemStatusActive      ProposalSystemStatus = "active"
	SystemStatusMaintenance ProposalSystemStatus = "maintenance"
	SystemStatusEmergency   ProposalSystemStatus = "emergency"
	SystemStatusPaused      ProposalSystemStatus = "paused"
	SystemStatusUpgrading   ProposalSystemStatus = "upgrading"
)

// PhaseTransition represents a transition between phases
type PhaseTransition struct {
	From              ProposalSystemPhase `json:"from"`
	To                ProposalSystemPhase `json:"to"`
	TransitionDate    time.Time           `json:"transition_date"`
	Reason            string              `json:"reason"`
	ApprovedBy        sdk.AccAddress      `json:"approved_by"`
	CommunityConsent  bool                `json:"community_consent"`
	TransitionMetrics TransitionMetrics   `json:"transition_metrics"`
}

// TransitionMetrics represents metrics during transition
type TransitionMetrics struct {
	ProposalsActive      uint64         `json:"proposals_active"`
	ProposalsCompleted   uint64         `json:"proposals_completed"`
	CommunityReadiness   uint8          `json:"community_readiness"`
	SystemHealth         uint8          `json:"system_health"`
	TransitionRisk       uint8          `json:"transition_risk"`
	StakeholderSupport   math.LegacyDec `json:"stakeholder_support"`
}

// ProposalTypeConfig represents configuration for proposal types
type ProposalTypeConfig struct {
	Type                ProposalType              `json:"type"`
	Name                string                    `json:"name"`
	Description         string                    `json:"description"`
	EnabledInPhase      []ProposalSystemPhase     `json:"enabled_in_phase"`
	RequiredApprovers   []ApproverConfig          `json:"required_approvers"`
	ThresholdOverrides  map[string]math.LegacyDec `json:"threshold_overrides"`
	Limits              ProposalLimits            `json:"limits"`
	ReviewRequirements  ReviewRequirements        `json:"review_requirements"`
	Priority            uint8                     `json:"priority"`
}

// ApproverConfig represents approver configuration
type ApproverConfig struct {
	Role              string  `json:"role"`
	Required          bool    `json:"required"`
	VetoEnabled       bool    `json:"veto_enabled"`
	WeightMultiplier  float64 `json:"weight_multiplier"`
	TimeLimit         time.Duration `json:"time_limit"`
}

// ProposalLimits represents limits for proposals
type ProposalLimits struct {
	MinAmount           sdk.Coin       `json:"min_amount"`
	MaxAmount           sdk.Coin       `json:"max_amount"`
	MaxPerPeriod        uint64         `json:"max_per_period"`
	Period              time.Duration  `json:"period"`
	CoolingPeriod       time.Duration  `json:"cooling_period"`
	MaxActiveProposals  uint64         `json:"max_active_proposals"`
	RequiredStake       sdk.Coin       `json:"required_stake"`
	RequiredReputation  uint8          `json:"required_reputation"`
}

// ReviewRequirements represents review requirements
type ReviewRequirements struct {
	TechnicalReview     bool           `json:"technical_review"`
	FinancialReview     bool           `json:"financial_review"`
	SecurityReview      bool           `json:"security_review"`
	CommunityReview     bool           `json:"community_review"`
	AuditRequired       bool           `json:"audit_required"`
	MinReviewPeriod     time.Duration  `json:"min_review_period"`
	MaxReviewPeriod     time.Duration  `json:"max_review_period"`
	RequiredReviewScore uint8          `json:"required_review_score"`
}

// PhaseGovernanceParams represents governance parameters for each phase
type PhaseGovernanceParams struct {
	CurrentPhase        ProposalSystemPhase       `json:"current_phase"`
	FounderDriven       FounderDrivenParams       `json:"founder_driven"`       // Years 0-3
	Transitional        TransitionalParams        `json:"transitional"`         // Year 3-4
	CommunityProposal   CommunityProposalParams   `json:"community_proposal"`   // Year 4-7
	FullGovernance      FullGovernanceParams      `json:"full_governance"`      // Year 7+
	Emergency           EmergencyParams           `json:"emergency"`            // Special
}

// FounderDrivenParams represents parameters for founder-driven phase (Years 0-3)
type FounderDrivenParams struct {
	FounderAllocationPower  uint8              `json:"founder_allocation_power"`  // 100%
	CommunityVisibility     uint8              `json:"community_visibility"`      // 100%
	CommunityFeedback       bool               `json:"community_feedback"`        // true
	CommunityVoting         bool               `json:"community_voting"`          // false
	FounderVetoRequired     bool               `json:"founder_veto_required"`     // false (founder decides)
	ProposalTypes           []ProposalType     `json:"proposal_types"`
	MaxProposalSize         sdk.Coin           `json:"max_proposal_size"`
	ReviewPeriod            time.Duration      `json:"review_period"`
	ExecutionSpeed          ExecutionSpeed     `json:"execution_speed"`
	TransparencyLevel       uint8              `json:"transparency_level"`
	AuditFrequency          time.Duration      `json:"audit_frequency"`
	ReportingRequirements   []string           `json:"reporting_requirements"`
	AllowedCategories       []string           `json:"allowed_categories"`
	RestrictedActions       []string           `json:"restricted_actions"`
}

// TransitionalParams represents parameters for transitional phase (Year 3-4)
type TransitionalParams struct {
	FounderAllocationPower  uint8              `json:"founder_allocation_power"`  // 70%
	CommunityProposalPower  uint8              `json:"community_proposal_power"`  // 30%
	CommunityVotingEnabled  bool               `json:"community_voting_enabled"`  // true (advisory)
	FounderVetoEnabled      bool               `json:"founder_veto_enabled"`      // true
	ProposalTypes           []ProposalType     `json:"proposal_types"`
	CommunityProposalLimit  sdk.Coin           `json:"community_proposal_limit"`
	RequiredCommunityStake  sdk.Coin           `json:"required_community_stake"`
	VotingPeriod            time.Duration      `json:"voting_period"`
	QuorumRequirement       math.LegacyDec     `json:"quorum_requirement"`
	PassingThreshold        math.LegacyDec     `json:"passing_threshold"`
	TrainingPrograms        []TrainingProgram  `json:"training_programs"`
	PilotProjects           []PilotProject     `json:"pilot_projects"`
	TransitionMetrics       []string           `json:"transition_metrics"`
	ReadinessAssessment     ReadinessAssessment `json:"readiness_assessment"`
}

// CommunityProposalParams represents parameters for community proposal phase (Year 4-7)
type CommunityProposalParams struct {
	FounderAllocationPower  uint8              `json:"founder_allocation_power"`  // 30%
	CommunityProposalPower  uint8              `json:"community_proposal_power"`  // 70%
	CommunityVotingEnabled  bool               `json:"community_voting_enabled"`  // true
	FounderVetoEnabled      bool               `json:"founder_veto_enabled"`      // true (limited)
	FounderVetoThreshold    sdk.Coin           `json:"founder_veto_threshold"`    // Only for large proposals
	ProposalTypes           []ProposalType     `json:"proposal_types"`
	MaxProposalSize         sdk.Coin           `json:"max_proposal_size"`
	RequiredStake           sdk.Coin           `json:"required_stake"`
	VotingPeriod            time.Duration      `json:"voting_period"`
	ReviewPeriod            time.Duration      `json:"review_period"`
	QuorumRequirement       math.LegacyDec     `json:"quorum_requirement"`
	PassingThreshold        math.LegacyDec     `json:"passing_threshold"`
	ProposalDeposit         sdk.Coin           `json:"proposal_deposit"`
	SlashingConditions      []SlashingCondition `json:"slashing_conditions"`
	ReputationRequirement   uint8              `json:"reputation_requirement"`
	CommunityCouncil        CommunityCouncil   `json:"community_council"`
	AppealProcess           AppealProcess      `json:"appeal_process"`
}

// FullGovernanceParams represents parameters for full governance phase (Year 7+)
type FullGovernanceParams struct {
	FounderAllocationPower  uint8              `json:"founder_allocation_power"`  // 10%
	CommunityProposalPower  uint8              `json:"community_proposal_power"`  // 90%
	CommunityVotingEnabled  bool               `json:"community_voting_enabled"`  // true
	FounderVetoEnabled      bool               `json:"founder_veto_enabled"`      // true (emergency only)
	FounderEmergencyPowers  []EmergencyPower   `json:"founder_emergency_powers"`
	ProposalTypes           []ProposalType     `json:"proposal_types"`
	MaxProposalSize         sdk.Coin           `json:"max_proposal_size"`
	RequiredStake           sdk.Coin           `json:"required_stake"`
	VotingPeriod            time.Duration      `json:"voting_period"`
	ReviewPeriod            time.Duration      `json:"review_period"`
	QuorumRequirement       math.LegacyDec     `json:"quorum_requirement"`
	PassingThreshold        math.LegacyDec     `json:"passing_threshold"`
	ConstitutionalRules     []ConstitutionalRule `json:"constitutional_rules"`
	GovernanceEvolution     GovernanceEvolution  `json:"governance_evolution"`
	DecentralizationScore   uint8              `json:"decentralization_score"`
	AutomationLevel         uint8              `json:"automation_level"`
}

// EmergencyParams represents parameters for emergency phase
type EmergencyParams struct {
	TriggerConditions       []EmergencyTrigger `json:"trigger_conditions"`
	FounderPowers           []EmergencyPower   `json:"founder_powers"`
	SecurityCouncilPowers   []EmergencyPower   `json:"security_council_powers"`
	Duration                time.Duration      `json:"duration"`
	ExtensionProcess        ExtensionProcess   `json:"extension_process"`
	CommunityNotification   bool               `json:"community_notification"`
	TransparencyRequired    bool               `json:"transparency_required"`
	PostEmergencyReview     bool               `json:"post_emergency_review"`
	AutomaticReversion      bool               `json:"automatic_reversion"`
	ReversionDate           time.Time          `json:"reversion_date"`
}

// ExecutionSpeed represents execution speed
type ExecutionSpeed string

const (
	ExecutionSpeedImmediate ExecutionSpeed = "immediate"
	ExecutionSpeedFast      ExecutionSpeed = "fast"
	ExecutionSpeedNormal    ExecutionSpeed = "normal"
	ExecutionSpeedSlow      ExecutionSpeed = "slow"
	ExecutionSpeedManual    ExecutionSpeed = "manual"
)

// TrainingProgram represents a training program
type TrainingProgram struct {
	ID              uint64    `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Duration        time.Duration `json:"duration"`
	Participants    uint64    `json:"participants"`
	CompletionRate  uint8     `json:"completion_rate"`
	Effectiveness   uint8     `json:"effectiveness"`
	NextSession     time.Time `json:"next_session"`
}

// PilotProject represents a pilot project
type PilotProject struct {
	ID              uint64    `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Budget          sdk.Coin  `json:"budget"`
	Duration        time.Duration `json:"duration"`
	Participants    uint64    `json:"participants"`
	SuccessMetrics  []string  `json:"success_metrics"`
	Status          string    `json:"status"`
	Results         string    `json:"results"`
	LessonsLearned  []string  `json:"lessons_learned"`
}

// ReadinessAssessment represents readiness assessment
type ReadinessAssessment struct {
	OverallScore        uint8     `json:"overall_score"`
	CommunityReadiness  uint8     `json:"community_readiness"`
	TechnicalReadiness  uint8     `json:"technical_readiness"`
	ProcessReadiness    uint8     `json:"process_readiness"`
	RiskReadiness       uint8     `json:"risk_readiness"`
	Recommendations     []string  `json:"recommendations"`
	BlockingIssues      []string  `json:"blocking_issues"`
	Timeline            string    `json:"timeline"`
	NextAssessment      time.Time `json:"next_assessment"`
}

// SlashingCondition represents slashing condition
type SlashingCondition struct {
	Condition       string   `json:"condition"`
	Penalty         sdk.Coin `json:"penalty"`
	Duration        time.Duration `json:"duration"`
	RepeatOffense   bool     `json:"repeat_offense"`
	AppealAllowed   bool     `json:"appeal_allowed"`
}

// CommunityCouncil represents community council
type CommunityCouncil struct {
	Members         []CouncilMember `json:"members"`
	ElectionProcess ElectionProcess `json:"election_process"`
	Term            time.Duration   `json:"term"`
	Powers          []string        `json:"powers"`
	Limitations     []string        `json:"limitations"`
	MeetingSchedule string          `json:"meeting_schedule"`
	Transparency    uint8           `json:"transparency"`
}

// CouncilMember represents a council member
type CouncilMember struct {
	Address         sdk.AccAddress `json:"address"`
	Role            string         `json:"role"`
	ElectedDate     time.Time      `json:"elected_date"`
	TermEnd         time.Time      `json:"term_end"`
	Votes           uint64         `json:"votes"`
	Reputation      uint8          `json:"reputation"`
	Attendance      uint8          `json:"attendance"`
	ProposalsLed    uint64         `json:"proposals_led"`
}

// ElectionProcess represents election process
type ElectionProcess struct {
	Type            string        `json:"type"`
	Frequency       time.Duration `json:"frequency"`
	NominationPeriod time.Duration `json:"nomination_period"`
	VotingPeriod    time.Duration `json:"voting_period"`
	Requirements    []string      `json:"requirements"`
	VotingMethod    string        `json:"voting_method"`
}

// AppealProcess represents appeal process
type AppealProcess struct {
	Enabled         bool          `json:"enabled"`
	AppealPeriod    time.Duration `json:"appeal_period"`
	ReviewBody      string        `json:"review_body"`
	RequiredStake   sdk.Coin      `json:"required_stake"`
	SuccessRate     uint8         `json:"success_rate"`
	AverageTime     time.Duration `json:"average_time"`
}

// EmergencyTrigger represents emergency trigger
type EmergencyTrigger struct {
	Type            string `json:"type"`
	Description     string `json:"description"`
	Threshold       string `json:"threshold"`
	AutoTrigger     bool   `json:"auto_trigger"`
	RequiredVotes   uint8  `json:"required_votes"`
	Duration        time.Duration `json:"duration"`
}

// EmergencyPower represents emergency power
type EmergencyPower struct {
	Power           string        `json:"power"`
	Description     string        `json:"description"`
	Limitations     []string      `json:"limitations"`
	Duration        time.Duration `json:"duration"`
	RequiresConsent bool          `json:"requires_consent"`
	Reversible      bool          `json:"reversible"`
}

// ExtensionProcess represents extension process
type ExtensionProcess struct {
	MaxExtensions   uint8         `json:"max_extensions"`
	ExtensionPeriod time.Duration `json:"extension_period"`
	RequiredVotes   uint8         `json:"required_votes"`
	ReviewRequired  bool          `json:"review_required"`
	Justification   bool          `json:"justification_required"`
}

// ConstitutionalRule represents constitutional rule
type ConstitutionalRule struct {
	ID              uint64   `json:"id"`
	Rule            string   `json:"rule"`
	Description     string   `json:"description"`
	Immutable       bool     `json:"immutable"`
	AmendmentProcess string  `json:"amendment_process"`
	Violations      uint64   `json:"violations"`
	Enforced        bool     `json:"enforced"`
}

// GovernanceEvolution represents governance evolution
type GovernanceEvolution struct {
	CurrentModel    string        `json:"current_model"`
	NextEvolution   string        `json:"next_evolution"`
	Timeline        time.Duration `json:"timeline"`
	Requirements    []string      `json:"requirements"`
	Benefits        []string      `json:"benefits"`
	Risks           []string      `json:"risks"`
	CommunitySupport uint8        `json:"community_support"`
}

// FounderControlConfig represents founder control configuration
type FounderControlConfig struct {
	CurrentPowers       []FounderPower      `json:"current_powers"`
	PowerTransition     PowerTransition     `json:"power_transition"`
	VetoRights          VetoRightConfig     `json:"veto_rights"`
	EmergencyOverride   bool                `json:"emergency_override"`
	SuccessionPlan      SuccessionPlan      `json:"succession_plan"`
	AccountabilityRules []AccountabilityRule `json:"accountability_rules"`
}

// FounderPower represents founder power
type FounderPower struct {
	Power           string              `json:"power"`
	Description     string              `json:"description"`
	Scope           string              `json:"scope"`
	ValidUntil      time.Time           `json:"valid_until"`
	TransitionTo    string              `json:"transition_to"`
	Conditions      []string            `json:"conditions"`
	UsageCount      uint64              `json:"usage_count"`
	LastUsed        time.Time           `json:"last_used"`
}

// PowerTransition represents power transition
type PowerTransition struct {
	Schedule        []TransitionMilestone `json:"schedule"`
	CurrentMilestone uint64               `json:"current_milestone"`
	Progress        uint8                 `json:"progress"`
	Blockers        []string              `json:"blockers"`
	NextTransition  time.Time             `json:"next_transition"`
}

// TransitionMilestone represents transition milestone
type TransitionMilestone struct {
	ID              uint64    `json:"id"`
	Name            string    `json:"name"`
	Date            time.Time `json:"date"`
	PowersRetained  []string  `json:"powers_retained"`
	PowersTransferred []string `json:"powers_transferred"`
	Requirements    []string  `json:"requirements"`
	Status          string    `json:"status"`
	Completed       bool      `json:"completed"`
}

// VetoRightConfig represents veto right configuration
type VetoRightConfig struct {
	Enabled         bool                        `json:"enabled"`
	Scope           []VetoScope                 `json:"scope"`
	Thresholds      map[string]sdk.Coin         `json:"thresholds"`
	TimeLimit       time.Duration               `json:"time_limit"`
	Override        VetoOverride                `json:"override"`
	UsageStats      VetoUsageStats              `json:"usage_stats"`
}

// VetoScope represents veto scope
type VetoScope struct {
	Category        string   `json:"category"`
	Enabled         bool     `json:"enabled"`
	Conditions      []string `json:"conditions"`
	ValidUntil      time.Time `json:"valid_until"`
}

// VetoOverride represents veto override
type VetoOverride struct {
	Possible        bool           `json:"possible"`
	RequiredVotes   math.LegacyDec `json:"required_votes"`
	CoolingPeriod   time.Duration  `json:"cooling_period"`
	AppealProcess   bool           `json:"appeal_process"`
}

// VetoUsageStats represents veto usage statistics
type VetoUsageStats struct {
	TotalVetos      uint64         `json:"total_vetos"`
	Overridden      uint64         `json:"overridden"`
	Categories      map[string]uint64 `json:"categories"`
	AverageResponse time.Duration  `json:"average_response"`
	Trend           TrendDirection `json:"trend"`
}

// SuccessionPlan represents succession plan
type SuccessionPlan struct {
	Enabled         bool              `json:"enabled"`
	Triggers        []SuccessionTrigger `json:"triggers"`
	Successors      []Successor        `json:"successors"`
	TransitionPeriod time.Duration     `json:"transition_period"`
	Requirements    []string           `json:"requirements"`
	LastUpdated     time.Time          `json:"last_updated"`
}

// SuccessionTrigger represents succession trigger
type SuccessionTrigger struct {
	Type            string `json:"type"`
	Description     string `json:"description"`
	AutoTrigger     bool   `json:"auto_trigger"`
	Verification    string `json:"verification"`
}

// Successor represents a successor
type Successor struct {
	Address         sdk.AccAddress `json:"address"`
	Priority        uint8          `json:"priority"`
	Qualifications  []string       `json:"qualifications"`
	Approved        bool           `json:"approved"`
	ApprovalDate    time.Time      `json:"approval_date"`
	TrainingStatus  string         `json:"training_status"`
}

// AccountabilityRule represents accountability rule
type AccountabilityRule struct {
	Rule            string   `json:"rule"`
	Description     string   `json:"description"`
	Enforcement     string   `json:"enforcement"`
	Penalties       []string `json:"penalties"`
	Violations      uint64   `json:"violations"`
	Active          bool     `json:"active"`
}

// CommunityPowerConfig represents community power configuration
type CommunityPowerConfig struct {
	CurrentPowers       []CommunityPower    `json:"current_powers"`
	PowerGrowth         PowerGrowthSchedule `json:"power_growth"`
	ParticipationRules  ParticipationRules  `json:"participation_rules"`
	DecisionMaking      DecisionMakingRules `json:"decision_making"`
	Safeguards          []Safeguard         `json:"safeguards"`
}

// CommunityPower represents community power
type CommunityPower struct {
	Power           string    `json:"power"`
	Description     string    `json:"description"`
	ActivationDate  time.Time `json:"activation_date"`
	Requirements    []string  `json:"requirements"`
	Limitations     []string  `json:"limitations"`
	UsageCount      uint64    `json:"usage_count"`
	Effectiveness   uint8     `json:"effectiveness"`
}

// PowerGrowthSchedule represents power growth schedule
type PowerGrowthSchedule struct {
	CurrentLevel    uint8                  `json:"current_level"`
	TargetLevel     uint8                  `json:"target_level"`
	GrowthRate      math.LegacyDec         `json:"growth_rate"`
	Milestones      []GrowthMilestone      `json:"milestones"`
	NextMilestone   time.Time              `json:"next_milestone"`
	Accelerators    []string               `json:"accelerators"`
	Impediments     []string               `json:"impediments"`
}

// GrowthMilestone represents growth milestone
type GrowthMilestone struct {
	Level           uint8     `json:"level"`
	Date            time.Time `json:"date"`
	PowersGranted   []string  `json:"powers_granted"`
	Requirements    []string  `json:"requirements"`
	Achieved        bool      `json:"achieved"`
}

// ParticipationRules represents participation rules
type ParticipationRules struct {
	MinStake        sdk.Coin       `json:"min_stake"`
	MinReputation   uint8          `json:"min_reputation"`
	MinActivity     uint8          `json:"min_activity"`
	VotingPower     VotingPowerCalc `json:"voting_power"`
	Delegation      bool           `json:"delegation"`
	ProxyVoting     bool           `json:"proxy_voting"`
	Incentives      []Incentive    `json:"incentives"`
}

// VotingPowerCalc represents voting power calculation
type VotingPowerCalc struct {
	BaseFormula     string         `json:"base_formula"`
	Multipliers     []Multiplier   `json:"multipliers"`
	Caps            VotingCaps     `json:"caps"`
	Distribution    PowerDistribution `json:"distribution"`
}

// Multiplier represents a multiplier
type Multiplier struct {
	Type            string         `json:"type"`
	Value           math.LegacyDec `json:"value"`
	Conditions      []string       `json:"conditions"`
	MaxMultiplier   math.LegacyDec `json:"max_multiplier"`
}

// VotingCaps represents voting caps
type VotingCaps struct {
	MaxPerAddress   math.LegacyDec `json:"max_per_address"`
	MaxPerProposal  math.LegacyDec `json:"max_per_proposal"`
	MaxDelegation   math.LegacyDec `json:"max_delegation"`
}

// PowerDistribution represents power distribution
type PowerDistribution struct {
	Top10Percent    math.LegacyDec `json:"top_10_percent"`
	Top50Percent    math.LegacyDec `json:"top_50_percent"`
	GiniCoefficient math.LegacyDec `json:"gini_coefficient"`
	Decentralization uint8         `json:"decentralization"`
}

// Incentive represents an incentive
type Incentive struct {
	Type            string   `json:"type"`
	Description     string   `json:"description"`
	Reward          sdk.Coin `json:"reward"`
	Requirements    []string `json:"requirements"`
	Active          bool     `json:"active"`
}

// DecisionMakingRules represents decision making rules
type DecisionMakingRules struct {
	ConsensusModel  string              `json:"consensus_model"`
	VotingMethods   []VotingMethod      `json:"voting_methods"`
	QuorumRules     QuorumRules         `json:"quorum_rules"`
	TimeConstraints TimeConstraints     `json:"time_constraints"`
	DisputeResolution DisputeResolution `json:"dispute_resolution"`
}

// VotingMethod represents voting method
type VotingMethod struct {
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	UseCase         []string `json:"use_case"`
	Algorithm       string   `json:"algorithm"`
	Transparency    uint8    `json:"transparency"`
}

// QuorumRules represents quorum rules
type QuorumRules struct {
	BaseQuorum      math.LegacyDec            `json:"base_quorum"`
	DynamicQuorum   bool                      `json:"dynamic_quorum"`
	CategoryQuorum  map[string]math.LegacyDec `json:"category_quorum"`
	QuorumDecay     bool                      `json:"quorum_decay"`
	MinParticipation math.LegacyDec           `json:"min_participation"`
}

// TimeConstraints represents time constraints
type TimeConstraints struct {
	MinVotingPeriod time.Duration             `json:"min_voting_period"`
	MaxVotingPeriod time.Duration             `json:"max_voting_period"`
	ReviewPeriod    time.Duration             `json:"review_period"`
	ExecutionDelay  time.Duration             `json:"execution_delay"`
	FastTrack       FastTrackRules            `json:"fast_track"`
}

// FastTrackRules represents fast track rules
type FastTrackRules struct {
	Enabled         bool           `json:"enabled"`
	Conditions      []string       `json:"conditions"`
	ReducedPeriod   time.Duration  `json:"reduced_period"`
	RequiredSupport math.LegacyDec `json:"required_support"`
}

// DisputeResolution represents dispute resolution
type DisputeResolution struct {
	Mechanism       string        `json:"mechanism"`
	Arbitrators     []Arbitrator  `json:"arbitrators"`
	Process         []ProcessStep `json:"process"`
	Timeline        time.Duration `json:"timeline"`
	Binding         bool          `json:"binding"`
	AppealProcess   bool          `json:"appeal_process"`
}

// Arbitrator represents an arbitrator
type Arbitrator struct {
	Address         sdk.AccAddress `json:"address"`
	Expertise       []string       `json:"expertise"`
	Cases           uint64         `json:"cases"`
	SuccessRate     uint8          `json:"success_rate"`
	Active          bool           `json:"active"`
}

// ProcessStep represents a process step
type ProcessStep struct {
	Step            uint8         `json:"step"`
	Name            string        `json:"name"`
	Duration        time.Duration `json:"duration"`
	Requirements    []string      `json:"requirements"`
	Outcome         string        `json:"outcome"`
}

// Safeguard represents a safeguard
type Safeguard struct {
	Type            string   `json:"type"`
	Description     string   `json:"description"`
	Trigger         string   `json:"trigger"`
	Action          string   `json:"action"`
	AutoActivate    bool     `json:"auto_activate"`
	Override        bool     `json:"override_possible"`
	LastActivated   time.Time `json:"last_activated"`
	Effectiveness   uint8    `json:"effectiveness"`
}

// TransitionSchedule represents transition schedule
type TransitionSchedule struct {
	Phases          []PhaseSchedule `json:"phases"`
	CurrentPhase    uint8           `json:"current_phase"`
	NextTransition  time.Time       `json:"next_transition"`
	Triggers        []TransitionTrigger `json:"triggers"`
	Checkpoints     []Checkpoint    `json:"checkpoints"`
	RollbackPlan    RollbackPlan    `json:"rollback_plan"`
}

// PhaseSchedule represents phase schedule
type PhaseSchedule struct {
	Phase           ProposalSystemPhase `json:"phase"`
	StartDate       time.Time           `json:"start_date"`
	EndDate         time.Time           `json:"end_date"`
	Duration        time.Duration       `json:"duration"`
	Objectives      []string            `json:"objectives"`
	SuccessCriteria []string            `json:"success_criteria"`
	Status          string              `json:"status"`
}

// TransitionTrigger represents transition trigger
type TransitionTrigger struct {
	Type            string `json:"type"`
	Condition       string `json:"condition"`
	AutoTrigger     bool   `json:"auto_trigger"`
	RequiredVotes   uint8  `json:"required_votes"`
	LastChecked     time.Time `json:"last_checked"`
}

// Checkpoint represents a checkpoint
type Checkpoint struct {
	ID              uint64    `json:"id"`
	Name            string    `json:"name"`
	Date            time.Time `json:"date"`
	Criteria        []string  `json:"criteria"`
	Passed          bool      `json:"passed"`
	Score           uint8     `json:"score"`
	Feedback        string    `json:"feedback"`
}

// RollbackPlan represents rollback plan
type RollbackPlan struct {
	Enabled         bool              `json:"enabled"`
	Triggers        []RollbackTrigger `json:"triggers"`
	Process         []string          `json:"process"`
	Timeline        time.Duration     `json:"timeline"`
	Authority       string            `json:"authority"`
	LastRollback    time.Time         `json:"last_rollback"`
}

// RollbackTrigger represents rollback trigger
type RollbackTrigger struct {
	Type            string `json:"type"`
	Threshold       string `json:"threshold"`
	Duration        time.Duration `json:"duration"`
	AutoTrigger     bool   `json:"auto_trigger"`
}

// EmergencyPowerConfig represents emergency power configuration
type EmergencyPowerConfig struct {
	Enabled         bool                    `json:"enabled"`
	Powers          []EmergencyPower        `json:"powers"`
	Activation      EmergencyActivation     `json:"activation"`
	Duration        time.Duration           `json:"duration"`
	Oversight       EmergencyOversight      `json:"oversight"`
	PostAction      PostEmergencyAction     `json:"post_action"`
}

// EmergencyActivation represents emergency activation
type EmergencyActivation struct {
	Triggers        []string          `json:"triggers"`
	RequiredVotes   uint8             `json:"required_votes"`
	FastTrack       bool              `json:"fast_track"`
	Notification    NotificationRules `json:"notification"`
	Documentation   bool              `json:"documentation_required"`
}

// NotificationRules represents notification rules
type NotificationRules struct {
	Immediate       []string      `json:"immediate"`
	Within24Hours   []string      `json:"within_24_hours"`
	Weekly          []string      `json:"weekly"`
	Channels        []string      `json:"channels"`
	Language        []string      `json:"languages"`
}

// EmergencyOversight represents emergency oversight
type EmergencyOversight struct {
	OversightBody   string        `json:"oversight_body"`
	ReportingFreq   time.Duration `json:"reporting_frequency"`
	Transparency    uint8         `json:"transparency"`
	Limitations     []string      `json:"limitations"`
	Accountability  []string      `json:"accountability"`
}

// PostEmergencyAction represents post-emergency action
type PostEmergencyAction struct {
	Review          bool          `json:"review_required"`
	Report          bool          `json:"report_required"`
	Timeline        time.Duration `json:"timeline"`
	Compensation    bool          `json:"compensation_possible"`
	PolicyUpdate    bool          `json:"policy_update_required"`
}

// ProposalSystemStats represents proposal system statistics
type ProposalSystemStats struct {
	TotalProposals      uint64         `json:"total_proposals"`
	ActiveProposals     uint64         `json:"active_proposals"`
	PassedProposals     uint64         `json:"passed_proposals"`
	RejectedProposals   uint64         `json:"rejected_proposals"`
	VetoedProposals     uint64         `json:"vetoed_proposals"`
	TotalFundsRequested sdk.Coin       `json:"total_funds_requested"`
	TotalFundsApproved  sdk.Coin       `json:"total_funds_approved"`
	AverageApprovalTime time.Duration  `json:"average_approval_time"`
	ParticipationRate   math.LegacyDec `json:"participation_rate"`
	SuccessRate         math.LegacyDec `json:"success_rate"`
	PhaseStats          map[string]PhaseStats `json:"phase_stats"`
	LastUpdated         time.Time      `json:"last_updated"`
}

// PhaseStats represents phase statistics
type PhaseStats struct {
	Proposals       uint64         `json:"proposals"`
	FundsAllocated  sdk.Coin       `json:"funds_allocated"`
	SuccessRate     math.LegacyDec `json:"success_rate"`
	Participation   math.LegacyDec `json:"participation"`
	Efficiency      uint8          `json:"efficiency"`
	Satisfaction    uint8          `json:"satisfaction"`
}

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

// Module account names for community proposal system
const (
	ProposalSystemModuleName    = "proposal_system"
	ProposalEscrowName          = "proposal_escrow"
	ProposalRewardsName         = "proposal_rewards"
	ProposalPenaltyName         = "proposal_penalty"
	ProposalEmergencyName       = "proposal_emergency"
	ProposalTransitionName      = "proposal_transition"
)

// Event types for community proposal system
const (
	EventTypeSystemInitialized      = "system_initialized"
	EventTypePhaseTransition        = "phase_transition"
	EventTypeProposalSubmitted      = "proposal_submitted"
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

// Phase durations
const (
	FounderDrivenDuration     = time.Hour * 24 * 365 * 3  // 3 years
	TransitionalDuration      = time.Hour * 24 * 365      // 1 year
	CommunityProposalDuration = time.Hour * 24 * 365 * 3  // 3 years
	// Full governance starts after year 7 and continues indefinitely
)

// Default phase configurations
var (
	DefaultFounderDrivenParams = FounderDrivenParams{
		FounderAllocationPower: 100,
		CommunityVisibility:    100,
		CommunityFeedback:      true,
		CommunityVoting:        false,
		FounderVetoRequired:    false,
		MaxProposalSize:        sdk.NewCoin("namo", math.NewInt(10000000)),
		ReviewPeriod:           time.Hour * 24 * 3,
		ExecutionSpeed:         ExecutionSpeedFast,
		TransparencyLevel:      10,
		AuditFrequency:         time.Hour * 24 * 90,
	}

	DefaultTransitionalParams = TransitionalParams{
		FounderAllocationPower: 70,
		CommunityProposalPower: 30,
		CommunityVotingEnabled: true,
		FounderVetoEnabled:     true,
		CommunityProposalLimit: sdk.NewCoin("namo", math.NewInt(1000000)),
		RequiredCommunityStake: sdk.NewCoin("namo", math.NewInt(10000)),
		VotingPeriod:           time.Hour * 24 * 7,
		QuorumRequirement:      math.LegacyNewDecWithPrec(20, 2),
		PassingThreshold:       math.LegacyNewDecWithPrec(51, 2),
	}

	DefaultCommunityProposalParams = CommunityProposalParams{
		FounderAllocationPower: 30,
		CommunityProposalPower: 70,
		CommunityVotingEnabled: true,
		FounderVetoEnabled:     true,
		FounderVetoThreshold:   sdk.NewCoin("namo", math.NewInt(5000000)),
		MaxProposalSize:        sdk.NewCoin("namo", math.NewInt(5000000)),
		RequiredStake:          sdk.NewCoin("namo", math.NewInt(5000)),
		VotingPeriod:           time.Hour * 24 * 14,
		ReviewPeriod:           time.Hour * 24 * 7,
		QuorumRequirement:      math.LegacyNewDecWithPrec(33, 2),
		PassingThreshold:       math.LegacyNewDecWithPrec(60, 2),
		ProposalDeposit:        sdk.NewCoin("namo", math.NewInt(1000)),
		ReputationRequirement:  6,
	}

	DefaultFullGovernanceParams = FullGovernanceParams{
		FounderAllocationPower: 10,
		CommunityProposalPower: 90,
		CommunityVotingEnabled: true,
		FounderVetoEnabled:     true,
		MaxProposalSize:        sdk.NewCoin("namo", math.NewInt(10000000)),
		RequiredStake:          sdk.NewCoin("namo", math.NewInt(1000)),
		VotingPeriod:           time.Hour * 24 * 21,
		ReviewPeriod:           time.Hour * 24 * 14,
		QuorumRequirement:      math.LegacyNewDecWithPrec(40, 2),
		PassingThreshold:       math.LegacyNewDecWithPrec(67, 2),
		DecentralizationScore:  8,
		AutomationLevel:        7,
	}
)