package types

import (
	"time"
)

// Identity Governance Framework Types

// IdentityGovernancePolicy represents a governance policy for identity management
type IdentityGovernancePolicy struct {
	PolicyID          string                    `json:"policy_id"`
	Name              string                    `json:"name"`
	Description       string                    `json:"description"`
	PolicyType        GovernancePolicyType      `json:"policy_type"`
	Status            PolicyStatus              `json:"status"`
	Version           string                    `json:"version"`
	EffectiveDate     time.Time                 `json:"effective_date"`
	ExpirationDate    *time.Time                `json:"expiration_date,omitempty"`
	Scope             GovernanceScope           `json:"scope"`
	Rules             []GovernanceRule          `json:"rules"`
	Enforcement       EnforcementLevel          `json:"enforcement"`
	Exceptions        []PolicyException         `json:"exceptions,omitempty"`
	ApprovalChain     []ApprovalStep            `json:"approval_chain"`
	MonitoringConfig  MonitoringConfiguration   `json:"monitoring_config"`
	ComplianceRefs    []ComplianceReference     `json:"compliance_references"`
	CreatedBy         string                    `json:"created_by"`
	CreatedAt         time.Time                 `json:"created_at"`
	UpdatedAt         time.Time                 `json:"updated_at"`
	ApprovedAt        *time.Time                `json:"approved_at,omitempty"`
	ApprovedBy        string                    `json:"approved_by,omitempty"`
	LastReviewDate    *time.Time                `json:"last_review_date,omitempty"`
	NextReviewDate    time.Time                 `json:"next_review_date"`
	ReviewCycle       time.Duration             `json:"review_cycle"`
	Tags              []string                  `json:"tags,omitempty"`
	Metadata          map[string]interface{}    `json:"metadata,omitempty"`
}

// GovernancePolicyType categorizes different types of governance policies
type GovernancePolicyType int32

const (
	GovernancePolicyType_ACCESS_CONTROL        GovernancePolicyType = 0  // Access control policies
	GovernancePolicyType_DATA_CLASSIFICATION   GovernancePolicyType = 1  // Data classification and handling
	GovernancePolicyType_PRIVACY_PROTECTION    GovernancePolicyType = 2  // Privacy protection policies
	GovernancePolicyType_CREDENTIAL_LIFECYCLE  GovernancePolicyType = 3  // Credential lifecycle management
	GovernancePolicyType_IDENTITY_VERIFICATION GovernancePolicyType = 4  // Identity verification requirements
	GovernancePolicyType_AUDIT_COMPLIANCE      GovernancePolicyType = 5  // Audit and compliance policies
	GovernancePolicyType_RISK_MANAGEMENT       GovernancePolicyType = 6  // Risk management policies
	GovernancePolicyType_FEDERATION_TRUST      GovernancePolicyType = 7  // Federation and trust policies
	GovernancePolicyType_CONSENT_MANAGEMENT    GovernancePolicyType = 8  // Consent management policies
	GovernancePolicyType_RETENTION_DISPOSAL    GovernancePolicyType = 9  // Data retention and disposal
	GovernancePolicyType_INCIDENT_RESPONSE     GovernancePolicyType = 10 // Incident response policies
	GovernancePolicyType_CUSTOM                GovernancePolicyType = 11 // Custom policies
)

// PolicyStatus indicates the current status of a governance policy
type PolicyStatus int32

const (
	PolicyStatus_DRAFT      PolicyStatus = 0
	PolicyStatus_REVIEW     PolicyStatus = 1
	PolicyStatus_APPROVED   PolicyStatus = 2
	PolicyStatus_ACTIVE     PolicyStatus = 3
	PolicyStatus_SUSPENDED  PolicyStatus = 4
	PolicyStatus_DEPRECATED PolicyStatus = 5
	PolicyStatus_ARCHIVED   PolicyStatus = 6
)

// GovernanceScope defines the scope of a governance policy
type GovernanceScope struct {
	Modules          []string              `json:"modules,omitempty"`          // Which modules this applies to
	IdentityTypes    []string              `json:"identity_types,omitempty"`    // Which identity types
	CredentialTypes  []string              `json:"credential_types,omitempty"`  // Which credential types
	UserRoles        []string              `json:"user_roles,omitempty"`        // Which user roles
	GeographicRegions []string             `json:"geographic_regions,omitempty"` // Geographic scope
	Organizations    []string              `json:"organizations,omitempty"`     // Organizational scope
	DataCategories   []string              `json:"data_categories,omitempty"`   // Data category scope
	RiskLevels       []RiskLevel           `json:"risk_levels,omitempty"`       // Risk level scope
	TrustLevels      []TrustLevel          `json:"trust_levels,omitempty"`      // Trust level scope
	Conditions       []ScopeCondition      `json:"conditions,omitempty"`        // Additional conditions
	Exclusions       []ScopeExclusion      `json:"exclusions,omitempty"`        // Explicit exclusions
}

// GovernanceRule represents a specific rule within a policy
type GovernanceRule struct {
	RuleID          string                    `json:"rule_id"`
	Name            string                    `json:"name"`
	Description     string                    `json:"description"`
	RuleType        GovernanceRuleType        `json:"rule_type"`
	Condition       RuleCondition             `json:"condition"`
	Action          RuleAction                `json:"action"`
	Priority        int32                     `json:"priority"`
	Enabled         bool                      `json:"enabled"`
	Parameters      map[string]interface{}    `json:"parameters,omitempty"`
	Schedule        *RuleSchedule             `json:"schedule,omitempty"`
	Dependencies    []string                  `json:"dependencies,omitempty"`    // Other rule IDs
	Conflicts       []string                  `json:"conflicts,omitempty"`       // Conflicting rule IDs
	ErrorHandling   ErrorHandlingPolicy       `json:"error_handling"`
	AuditLevel      AuditLevel                `json:"audit_level"`
	CreatedAt       time.Time                 `json:"created_at"`
	UpdatedAt       time.Time                 `json:"updated_at"`
	LastExecuted    *time.Time                `json:"last_executed,omitempty"`
	ExecutionCount  int64                     `json:"execution_count"`
	SuccessCount    int64                     `json:"success_count"`
	FailureCount    int64                     `json:"failure_count"`
	AverageExecutionTime time.Duration        `json:"average_execution_time"`
}

// GovernanceRuleType categorizes different types of governance rules
type GovernanceRuleType int32

const (
	GovernanceRuleType_AUTHENTICATION      GovernanceRuleType = 0  // Authentication rules
	GovernanceRuleType_AUTHORIZATION       GovernanceRuleType = 1  // Authorization rules
	GovernanceRuleType_VALIDATION          GovernanceRuleType = 2  // Validation rules
	GovernanceRuleType_ENCRYPTION          GovernanceRuleType = 3  // Encryption rules
	GovernanceRuleType_AUDIT_LOGGING       GovernanceRuleType = 4  // Audit logging rules
	GovernanceRuleType_DATA_MASKING        GovernanceRuleType = 5  // Data masking rules
	GovernanceRuleType_RATE_LIMITING       GovernanceRuleType = 6  // Rate limiting rules
	GovernanceRuleType_ANOMALY_DETECTION   GovernanceRuleType = 7  // Anomaly detection rules
	GovernanceRuleType_LIFECYCLE_MANAGEMENT GovernanceRuleType = 8  // Lifecycle management
	GovernanceRuleType_COMPLIANCE_CHECK    GovernanceRuleType = 9  // Compliance checking
	GovernanceRuleType_NOTIFICATION        GovernanceRuleType = 10 // Notification rules
	GovernanceRuleType_REMEDIATION         GovernanceRuleType = 11 // Remediation rules
	GovernanceRuleType_CUSTOM              GovernanceRuleType = 12 // Custom rules
)

// RuleCondition defines when a rule should be triggered
type RuleCondition struct {
	Expression    string                    `json:"expression"`     // Boolean expression
	EventTypes    []string                  `json:"event_types,omitempty"`
	TimeWindows   []TimeWindow              `json:"time_windows,omitempty"`
	Attributes    map[string]interface{}    `json:"attributes,omitempty"`
	Thresholds    map[string]float64        `json:"thresholds,omitempty"`
	Patterns      []ConditionPattern        `json:"patterns,omitempty"`
	CustomLogic   string                    `json:"custom_logic,omitempty"`
}

// RuleAction defines what action to take when a rule is triggered
type RuleAction struct {
	ActionType    GovernanceActionType      `json:"action_type"`
	Parameters    map[string]interface{}    `json:"parameters,omitempty"`
	TargetFields  []string                  `json:"target_fields,omitempty"`
	Notifications []NotificationConfig      `json:"notifications,omitempty"`
	Escalations   []EscalationConfig        `json:"escalations,omitempty"`
	Remediation   *RemediationConfig        `json:"remediation,omitempty"`
	CustomAction  string                    `json:"custom_action,omitempty"`
}

// GovernanceActionType categorizes different types of governance actions
type GovernanceActionType int32

const (
	GovernanceActionType_ALLOW              GovernanceActionType = 0  // Allow the operation
	GovernanceActionType_DENY               GovernanceActionType = 1  // Deny the operation
	GovernanceActionType_REQUIRE_APPROVAL   GovernanceActionType = 2  // Require approval
	GovernanceActionType_LOG_WARNING        GovernanceActionType = 3  // Log a warning
	GovernanceActionType_LOG_ERROR          GovernanceActionType = 4  // Log an error
	GovernanceActionType_SEND_NOTIFICATION  GovernanceActionType = 5  // Send notification
	GovernanceActionType_QUARANTINE         GovernanceActionType = 6  // Quarantine data/identity
	GovernanceActionType_MASK_DATA          GovernanceActionType = 7  // Mask sensitive data
	GovernanceActionType_ENCRYPT_DATA       GovernanceActionType = 8  // Encrypt data
	GovernanceActionType_RATE_LIMIT         GovernanceActionType = 9  // Apply rate limiting
	GovernanceActionType_SUSPEND_IDENTITY   GovernanceActionType = 10 // Suspend identity
	GovernanceActionType_REVOKE_CREDENTIAL  GovernanceActionType = 11 // Revoke credential
	GovernanceActionType_TRIGGER_AUDIT      GovernanceActionType = 12 // Trigger audit
	GovernanceActionType_INITIATE_WORKFLOW  GovernanceActionType = 13 // Initiate workflow
	GovernanceActionType_CUSTOM             GovernanceActionType = 14 // Custom action
)

// EnforcementLevel defines how strictly a policy is enforced
type EnforcementLevel int32

const (
	EnforcementLevel_ADVISORY    EnforcementLevel = 0 // Advisory only, no enforcement
	EnforcementLevel_WARNING     EnforcementLevel = 1 // Warning with logging
	EnforcementLevel_STRICT      EnforcementLevel = 2 // Strict enforcement
	EnforcementLevel_CRITICAL    EnforcementLevel = 3 // Critical enforcement with blocking
)

// PolicyException represents an exception to a governance policy
type PolicyException struct {
	ExceptionID     string                    `json:"exception_id"`
	Name            string                    `json:"name"`
	Description     string                    `json:"description"`
	Justification   string                    `json:"justification"`
	Scope           GovernanceScope           `json:"scope"`
	EffectiveDate   time.Time                 `json:"effective_date"`
	ExpirationDate  time.Time                 `json:"expiration_date"`
	ApprovedBy      string                    `json:"approved_by"`
	ApprovedAt      time.Time                 `json:"approved_at"`
	ReviewRequired  bool                      `json:"review_required"`
	ReviewDate      time.Time                 `json:"review_date"`
	UsageCount      int64                     `json:"usage_count"`
	MaxUsage        int64                     `json:"max_usage,omitempty"`
	Conditions      []ExceptionCondition      `json:"conditions,omitempty"`
	Metadata        map[string]interface{}    `json:"metadata,omitempty"`
}

// ApprovalStep represents a step in the approval workflow
type ApprovalStep struct {
	StepID          string                    `json:"step_id"`
	Name            string                    `json:"name"`
	Description     string                    `json:"description"`
	RequiredRole    string                    `json:"required_role"`
	ApproverGroup   []string                  `json:"approver_group,omitempty"`
	MinApprovals    int32                     `json:"min_approvals"`
	MaxApprovals    int32                     `json:"max_approvals,omitempty"`
	TimeoutDuration time.Duration             `json:"timeout_duration"`
	EscalationRules []EscalationRule          `json:"escalation_rules,omitempty"`
	ParallelSteps   []string                  `json:"parallel_steps,omitempty"`
	Conditions      []ApprovalCondition       `json:"conditions,omitempty"`
	OnApprove       *ApprovalAction           `json:"on_approve,omitempty"`
	OnReject        *ApprovalAction           `json:"on_reject,omitempty"`
	OnTimeout       *ApprovalAction           `json:"on_timeout,omitempty"`
}

// IdentityGovernanceRole represents a governance role with specific permissions
type IdentityGovernanceRole struct {
	RoleID          string                    `json:"role_id"`
	Name            string                    `json:"name"`
	Description     string                    `json:"description"`
	RoleType        GovernanceRoleType        `json:"role_type"`
	Level           GovernanceLevel           `json:"level"`
	Permissions     []GovernancePermission    `json:"permissions"`
	Responsibilities []string                 `json:"responsibilities"`
	Prerequisites   []string                  `json:"prerequisites"`
	Delegation      DelegationConfig          `json:"delegation"`
	RotationPolicy  RotationPolicy            `json:"rotation_policy"`
	AccessControls  RoleAccessControls        `json:"access_controls"`
	AuditRequirements AuditRequirements       `json:"audit_requirements"`
	CreatedAt       time.Time                 `json:"created_at"`
	UpdatedAt       time.Time                 `json:"updated_at"`
	EffectiveDate   time.Time                 `json:"effective_date"`
	ExpirationDate  *time.Time                `json:"expiration_date,omitempty"`
	Status          RoleStatus                `json:"status"`
	Metadata        map[string]interface{}    `json:"metadata,omitempty"`
}

// GovernanceRoleType categorizes different types of governance roles
type GovernanceRoleType int32

const (
	GovernanceRoleType_IDENTITY_ADMIN        GovernanceRoleType = 0  // Identity administrator
	GovernanceRoleType_POLICY_MANAGER        GovernanceRoleType = 1  // Policy manager
	GovernanceRoleType_COMPLIANCE_OFFICER    GovernanceRoleType = 2  // Compliance officer
	GovernanceRoleType_SECURITY_OFFICER      GovernanceRoleType = 3  // Security officer
	GovernanceRoleType_AUDIT_MANAGER         GovernanceRoleType = 4  // Audit manager
	GovernanceRoleType_DATA_STEWARD          GovernanceRoleType = 5  // Data steward
	GovernanceRoleType_PRIVACY_OFFICER       GovernanceRoleType = 6  // Privacy officer
	GovernanceRoleType_RISK_MANAGER          GovernanceRoleType = 7  // Risk manager
	GovernanceRoleType_WORKFLOW_APPROVER     GovernanceRoleType = 8  // Workflow approver
	GovernanceRoleType_TECHNICAL_ADMIN       GovernanceRoleType = 9  // Technical administrator
	GovernanceRoleType_BUSINESS_OWNER        GovernanceRoleType = 10 // Business owner
	GovernanceRoleType_EXTERNAL_AUDITOR      GovernanceRoleType = 11 // External auditor
	GovernanceRoleType_CUSTOM                GovernanceRoleType = 12 // Custom role
)

// GovernanceLevel indicates the level of governance authority
type GovernanceLevel int32

const (
	GovernanceLevel_OPERATIONAL   GovernanceLevel = 0 // Operational level
	GovernanceLevel_TACTICAL      GovernanceLevel = 1 // Tactical level
	GovernanceLevel_STRATEGIC     GovernanceLevel = 2 // Strategic level
	GovernanceLevel_EXECUTIVE     GovernanceLevel = 3 // Executive level
)

// GovernancePermission represents a specific permission in the governance system
type GovernancePermission struct {
	PermissionID    string                    `json:"permission_id"`
	Name            string                    `json:"name"`
	Description     string                    `json:"description"`
	Resource        string                    `json:"resource"`
	Actions         []string                  `json:"actions"`
	Conditions      []PermissionCondition     `json:"conditions,omitempty"`
	Constraints     []PermissionConstraint    `json:"constraints,omitempty"`
	EffectiveDate   time.Time                 `json:"effective_date"`
	ExpirationDate  *time.Time                `json:"expiration_date,omitempty"`
	GrantedBy       string                    `json:"granted_by"`
	GrantedAt       time.Time                 `json:"granted_at"`
	LastUsed        *time.Time                `json:"last_used,omitempty"`
	UsageCount      int64                     `json:"usage_count"`
	Metadata        map[string]interface{}    `json:"metadata,omitempty"`
}

// GovernanceWorkflow represents a governance workflow
type GovernanceWorkflow struct {
	WorkflowID      string                    `json:"workflow_id"`
	Name            string                    `json:"name"`
	Description     string                    `json:"description"`
	WorkflowType    GovernanceWorkflowType    `json:"workflow_type"`
	Version         string                    `json:"version"`
	Status          WorkflowStatus            `json:"status"`
	TriggerEvents   []WorkflowTrigger         `json:"trigger_events"`
	Steps           []WorkflowStep            `json:"steps"`
	Variables       map[string]interface{}    `json:"variables,omitempty"`
	TimeoutDuration time.Duration             `json:"timeout_duration"`
	RetryPolicy     WorkflowRetryPolicy       `json:"retry_policy"`
	ErrorHandling   WorkflowErrorHandling     `json:"error_handling"`
	Notifications   []WorkflowNotification    `json:"notifications,omitempty"`
	Metrics         WorkflowMetrics           `json:"metrics"`
	CreatedBy       string                    `json:"created_by"`
	CreatedAt       time.Time                 `json:"created_at"`
	UpdatedAt       time.Time                 `json:"updated_at"`
	LastExecuted    *time.Time                `json:"last_executed,omitempty"`
	ExecutionCount  int64                     `json:"execution_count"`
	SuccessRate     float64                   `json:"success_rate"`
	Metadata        map[string]interface{}    `json:"metadata,omitempty"`
}

// GovernanceWorkflowType categorizes different types of governance workflows
type GovernanceWorkflowType int32

const (
	GovernanceWorkflowType_APPROVAL           GovernanceWorkflowType = 0  // Approval workflows
	GovernanceWorkflowType_COMPLIANCE_CHECK   GovernanceWorkflowType = 1  // Compliance checking
	GovernanceWorkflowType_INCIDENT_RESPONSE  GovernanceWorkflowType = 2  // Incident response
	GovernanceWorkflowType_POLICY_ENFORCEMENT GovernanceWorkflowType = 3  // Policy enforcement
	GovernanceWorkflowType_RISK_ASSESSMENT    GovernanceWorkflowType = 4  // Risk assessment
	GovernanceWorkflowType_AUDIT_PROCESS      GovernanceWorkflowType = 5  // Audit processes
	GovernanceWorkflowType_LIFECYCLE_MGMT     GovernanceWorkflowType = 6  // Lifecycle management
	GovernanceWorkflowType_EXCEPTION_HANDLING GovernanceWorkflowType = 7  // Exception handling
	GovernanceWorkflowType_REMEDIATION        GovernanceWorkflowType = 8  // Remediation workflows
	GovernanceWorkflowType_CUSTOM             GovernanceWorkflowType = 9  // Custom workflows
)

// WorkflowStatus indicates the current status of a workflow
type WorkflowStatus int32

const (
	WorkflowStatus_DRAFT      WorkflowStatus = 0
	WorkflowStatus_ACTIVE     WorkflowStatus = 1
	WorkflowStatus_SUSPENDED  WorkflowStatus = 2
	WorkflowStatus_DEPRECATED WorkflowStatus = 3
	WorkflowStatus_ARCHIVED   WorkflowStatus = 4
)

// RoleStatus indicates the current status of a governance role
type RoleStatus int32

const (
	RoleStatus_ACTIVE     RoleStatus = 0
	RoleStatus_INACTIVE   RoleStatus = 1
	RoleStatus_SUSPENDED  RoleStatus = 2
	RoleStatus_DEPRECATED RoleStatus = 3
)

// GovernanceDecision represents a governance decision
type GovernanceDecision struct {
	DecisionID      string                    `json:"decision_id"`
	WorkflowID      string                    `json:"workflow_id,omitempty"`
	PolicyID        string                    `json:"policy_id,omitempty"`
	DecisionType    GovernanceDecisionType    `json:"decision_type"`
	Context         DecisionContext           `json:"context"`
	Decision        DecisionOutcome           `json:"decision"`
	Rationale       string                    `json:"rationale"`
	Evidence        []DecisionEvidence        `json:"evidence,omitempty"`
	DecisionMaker   string                    `json:"decision_maker"`
	DecisionDate    time.Time                 `json:"decision_date"`
	EffectiveDate   time.Time                 `json:"effective_date"`
	ReviewDate      *time.Time                `json:"review_date,omitempty"`
	Impact          DecisionImpact            `json:"impact"`
	Precedent       bool                      `json:"precedent"`        // Sets precedent for future decisions
	AppealAllowed   bool                      `json:"appeal_allowed"`
	AppealDeadline  *time.Time                `json:"appeal_deadline,omitempty"`
	Status          DecisionStatus            `json:"status"`
	AuditTrail      []DecisionAuditEntry      `json:"audit_trail"`
	Metadata        map[string]interface{}    `json:"metadata,omitempty"`
}

// GovernanceDecisionType categorizes different types of governance decisions
type GovernanceDecisionType int32

const (
	GovernanceDecisionType_POLICY_APPROVAL     GovernanceDecisionType = 0
	GovernanceDecisionType_EXCEPTION_APPROVAL  GovernanceDecisionType = 1
	GovernanceDecisionType_ACCESS_GRANT        GovernanceDecisionType = 2
	GovernanceDecisionType_RISK_ACCEPTANCE     GovernanceDecisionType = 3
	GovernanceDecisionType_COMPLIANCE_WAIVER   GovernanceDecisionType = 4
	GovernanceDecisionType_INCIDENT_RESOLUTION GovernanceDecisionType = 5
	GovernanceDecisionType_AUDIT_FINDING       GovernanceDecisionType = 6
	GovernanceDecisionType_REMEDIATION_PLAN    GovernanceDecisionType = 7
	GovernanceDecisionType_ROLE_ASSIGNMENT     GovernanceDecisionType = 8
	GovernanceDecisionType_DELEGATION          GovernanceDecisionType = 9
	GovernanceDecisionType_CUSTOM              GovernanceDecisionType = 10
)

// DecisionOutcome represents the outcome of a governance decision
type DecisionOutcome int32

const (
	DecisionOutcome_APPROVED              DecisionOutcome = 0
	DecisionOutcome_REJECTED              DecisionOutcome = 1
	DecisionOutcome_CONDITIONALLY_APPROVED DecisionOutcome = 2
	DecisionOutcome_DEFERRED              DecisionOutcome = 3
	DecisionOutcome_ESCALATED             DecisionOutcome = 4
	DecisionOutcome_WITHDRAWN             DecisionOutcome = 5
)

// DecisionStatus indicates the current status of a decision
type DecisionStatus int32

const (
	DecisionStatus_PENDING    DecisionStatus = 0
	DecisionStatus_FINAL      DecisionStatus = 1
	DecisionStatus_APPEALED   DecisionStatus = 2
	DecisionStatus_OVERTURNED DecisionStatus = 3
	DecisionStatus_EXPIRED    DecisionStatus = 4
)

// Supporting types for complex structures

type ScopeCondition struct {
	Field     string      `json:"field"`
	Operator  string      `json:"operator"`
	Value     interface{} `json:"value"`
	Logic     string      `json:"logic,omitempty"` // AND, OR
}

type ScopeExclusion struct {
	Type        string      `json:"type"`
	Identifier  string      `json:"identifier"`
	Reason      string      `json:"reason"`
	EffectiveDate time.Time `json:"effective_date"`
	ExpirationDate *time.Time `json:"expiration_date,omitempty"`
}

type RuleSchedule struct {
	ScheduleType ScheduleType `json:"schedule_type"`
	CronExpression string     `json:"cron_expression,omitempty"`
	Interval     time.Duration `json:"interval,omitempty"`
	StartTime    time.Time    `json:"start_time"`
	EndTime      *time.Time   `json:"end_time,omitempty"`
	TimeZone     string       `json:"timezone"`
}

type ScheduleType int32

const (
	ScheduleType_IMMEDIATE  ScheduleType = 0
	ScheduleType_SCHEDULED  ScheduleType = 1
	ScheduleType_PERIODIC   ScheduleType = 2
	ScheduleType_EVENT_DRIVEN ScheduleType = 3
)

type TimeWindow struct {
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	TimeZone string    `json:"timezone"`
	Recurring bool     `json:"recurring"`
	Pattern   string   `json:"pattern,omitempty"` // For recurring windows
}

type ConditionPattern struct {
	Name        string `json:"name"`
	Pattern     string `json:"pattern"`
	PatternType string `json:"pattern_type"` // regex, glob, etc.
	CaseSensitive bool `json:"case_sensitive"`
}

type NotificationConfig struct {
	Channel     string                 `json:"channel"`    // email, slack, webhook, etc.
	Recipients  []string               `json:"recipients"`
	Template    string                 `json:"template"`
	Urgency     NotificationUrgency    `json:"urgency"`
	Conditions  []NotificationCondition `json:"conditions,omitempty"`
}

type NotificationUrgency int32

const (
	NotificationUrgency_LOW       NotificationUrgency = 0
	NotificationUrgency_NORMAL    NotificationUrgency = 1
	NotificationUrgency_HIGH      NotificationUrgency = 2
	NotificationUrgency_CRITICAL  NotificationUrgency = 3
)

type NotificationCondition struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

type EscalationConfig struct {
	Level       int32         `json:"level"`
	Trigger     string        `json:"trigger"`
	Recipients  []string      `json:"recipients"`
	Timeout     time.Duration `json:"timeout"`
	Actions     []string      `json:"actions"`
}

type EscalationRule struct {
	TriggerAfter time.Duration `json:"trigger_after"`
	Escalate     []string      `json:"escalate_to"`
	Actions      []string      `json:"actions"`
}

type RemediationConfig struct {
	RemediationType RemediationType        `json:"remediation_type"`
	AutoRemediate   bool                   `json:"auto_remediate"`
	Steps           []RemediationStep      `json:"steps"`
	Timeout         time.Duration          `json:"timeout"`
	Rollback        *RollbackConfig        `json:"rollback,omitempty"`
}

type RemediationType int32

const (
	RemediationType_AUTOMATIC RemediationType = 0
	RemediationType_MANUAL    RemediationType = 1
	RemediationType_HYBRID    RemediationType = 2
)

type RemediationStep struct {
	StepID      string                 `json:"step_id"`
	Name        string                 `json:"name"`
	Action      string                 `json:"action"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
	Timeout     time.Duration          `json:"timeout"`
	OnSuccess   string                 `json:"on_success,omitempty"`
	OnFailure   string                 `json:"on_failure,omitempty"`
}

type RollbackConfig struct {
	Enabled       bool          `json:"enabled"`
	AutoRollback  bool          `json:"auto_rollback"`
	Timeout       time.Duration `json:"timeout"`
	Steps         []RemediationStep `json:"steps"`
}

type ExceptionCondition struct {
	Field     string      `json:"field"`
	Operator  string      `json:"operator"`
	Value     interface{} `json:"value"`
	Required  bool        `json:"required"`
}

type ApprovalCondition struct {
	Field     string      `json:"field"`
	Operator  string      `json:"operator"`
	Value     interface{} `json:"value"`
	Logic     string      `json:"logic,omitempty"`
}

type ApprovalAction struct {
	ActionType  string                 `json:"action_type"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
	NextStep    string                 `json:"next_step,omitempty"`
}

type DelegationConfig struct {
	Allowed           bool          `json:"allowed"`
	MaxDepth          int32         `json:"max_depth"`
	AllowedDelegates  []string      `json:"allowed_delegates,omitempty"`
	RequiresApproval  bool          `json:"requires_approval"`
	ExpirationPolicy  time.Duration `json:"expiration_policy"`
	AuditRequired     bool          `json:"audit_required"`
}

type RotationPolicy struct {
	Enabled         bool          `json:"enabled"`
	RotationPeriod  time.Duration `json:"rotation_period"`
	OverlapPeriod   time.Duration `json:"overlap_period"`
	RequiresHandoff bool          `json:"requires_handoff"`
	KnowledgeTransfer bool        `json:"knowledge_transfer"`
}

type RoleAccessControls struct {
	IPRestrictions    []string      `json:"ip_restrictions,omitempty"`
	TimeRestrictions  []TimeWindow  `json:"time_restrictions,omitempty"`
	LocationRestrictions []string   `json:"location_restrictions,omitempty"`
	DeviceRestrictions []string     `json:"device_restrictions,omitempty"`
	RequiresMFA       bool          `json:"requires_mfa"`
	SessionTimeout    time.Duration `json:"session_timeout"`
	ConcurrentSessions int32        `json:"concurrent_sessions"`
}

type PermissionCondition struct {
	Context   string      `json:"context"`
	Field     string      `json:"field"`
	Operator  string      `json:"operator"`
	Value     interface{} `json:"value"`
	Required  bool        `json:"required"`
}

type PermissionConstraint struct {
	Type        string                 `json:"type"`
	Parameters  map[string]interface{} `json:"parameters"`
	ErrorMessage string                `json:"error_message"`
}

type WorkflowTrigger struct {
	TriggerID   string                 `json:"trigger_id"`
	EventType   string                 `json:"event_type"`
	Conditions  []TriggerCondition     `json:"conditions,omitempty"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
}

type TriggerCondition struct {
	Field     string      `json:"field"`
	Operator  string      `json:"operator"`
	Value     interface{} `json:"value"`
}

type WorkflowStep struct {
	StepID        string                 `json:"step_id"`
	Name          string                 `json:"name"`
	StepType      WorkflowStepType       `json:"step_type"`
	Action        string                 `json:"action"`
	Parameters    map[string]interface{} `json:"parameters,omitempty"`
	Timeout       time.Duration          `json:"timeout"`
	RetryCount    int32                  `json:"retry_count"`
	Dependencies  []string               `json:"dependencies,omitempty"`
	OnSuccess     string                 `json:"on_success,omitempty"`
	OnFailure     string                 `json:"on_failure,omitempty"`
	OnTimeout     string                 `json:"on_timeout,omitempty"`
}

type WorkflowStepType int32

const (
	WorkflowStepType_APPROVAL       WorkflowStepType = 0
	WorkflowStepType_VALIDATION     WorkflowStepType = 1
	WorkflowStepType_NOTIFICATION   WorkflowStepType = 2
	WorkflowStepType_COMPUTATION    WorkflowStepType = 3
	WorkflowStepType_INTEGRATION    WorkflowStepType = 4
	WorkflowStepType_DECISION       WorkflowStepType = 5
	WorkflowStepType_CUSTOM         WorkflowStepType = 6
)

type WorkflowRetryPolicy struct {
	MaxRetries    int32         `json:"max_retries"`
	RetryDelay    time.Duration `json:"retry_delay"`
	BackoffFactor float64       `json:"backoff_factor"`
	MaxDelay      time.Duration `json:"max_delay"`
}

type WorkflowErrorHandling struct {
	OnError      WorkflowErrorAction `json:"on_error"`
	Rollback     bool                `json:"rollback"`
	Notifications []string           `json:"notifications,omitempty"`
	Escalations   []string           `json:"escalations,omitempty"`
}

type WorkflowErrorAction int32

const (
	WorkflowErrorAction_FAIL      WorkflowErrorAction = 0
	WorkflowErrorAction_RETRY     WorkflowErrorAction = 1
	WorkflowErrorAction_SKIP      WorkflowErrorAction = 2
	WorkflowErrorAction_ESCALATE  WorkflowErrorAction = 3
)

type WorkflowNotification struct {
	Event       string   `json:"event"`
	Recipients  []string `json:"recipients"`
	Template    string   `json:"template"`
	Channel     string   `json:"channel"`
}

type WorkflowMetrics struct {
	AverageExecutionTime time.Duration `json:"average_execution_time"`
	SuccessRate         float64       `json:"success_rate"`
	ErrorRate           float64       `json:"error_rate"`
	ThroughputPerHour   float64       `json:"throughput_per_hour"`
	PeakExecutionTime   time.Duration `json:"peak_execution_time"`
	LastUpdated         time.Time     `json:"last_updated"`
}

type DecisionContext struct {
	RequestID     string                 `json:"request_id"`
	RequestType   string                 `json:"request_type"`
	Requestor     string                 `json:"requestor"`
	ResourceType  string                 `json:"resource_type"`
	ResourceID    string                 `json:"resource_id"`
	Environment   string                 `json:"environment"`
	Risk          DecisionRiskContext    `json:"risk"`
	Compliance    DecisionComplianceContext `json:"compliance"`
	Business      DecisionBusinessContext `json:"business"`
	Technical     DecisionTechnicalContext `json:"technical"`
	Stakeholders  []string               `json:"stakeholders"`
	Deadline      *time.Time             `json:"deadline,omitempty"`
	Urgency       DecisionUrgency        `json:"urgency"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

type DecisionUrgency int32

const (
	DecisionUrgency_LOW       DecisionUrgency = 0
	DecisionUrgency_NORMAL    DecisionUrgency = 1
	DecisionUrgency_HIGH      DecisionUrgency = 2
	DecisionUrgency_CRITICAL  DecisionUrgency = 3
)

type DecisionRiskContext struct {
	RiskLevel     RiskLevel              `json:"risk_level"`
	RiskFactors   []string               `json:"risk_factors"`
	Mitigations   []string               `json:"mitigations"`
	Assessment    map[string]interface{} `json:"assessment,omitempty"`
}

type DecisionComplianceContext struct {
	Regulations   []string               `json:"regulations"`
	Requirements  []string               `json:"requirements"`
	Violations    []string               `json:"violations,omitempty"`
	Attestations  []string               `json:"attestations,omitempty"`
}

type DecisionBusinessContext struct {
	Impact        string                 `json:"impact"`
	Justification string                 `json:"justification"`
	CostBenefit   map[string]interface{} `json:"cost_benefit,omitempty"`
	Timeline      string                 `json:"timeline"`
}

type DecisionTechnicalContext struct {
	Implementation map[string]interface{} `json:"implementation,omitempty"`
	Dependencies   []string               `json:"dependencies,omitempty"`
	Constraints    []string               `json:"constraints,omitempty"`
	Resources      map[string]interface{} `json:"resources,omitempty"`
}

type DecisionEvidence struct {
	EvidenceID    string                 `json:"evidence_id"`
	Type          string                 `json:"type"`
	Source        string                 `json:"source"`
	Description   string                 `json:"description"`
	Content       map[string]interface{} `json:"content"`
	Reliability   EvidenceReliability    `json:"reliability"`
	CollectedAt   time.Time              `json:"collected_at"`
	CollectedBy   string                 `json:"collected_by"`
}

type EvidenceReliability int32

const (
	EvidenceReliability_LOW    EvidenceReliability = 0
	EvidenceReliability_MEDIUM EvidenceReliability = 1
	EvidenceReliability_HIGH   EvidenceReliability = 2
)

type DecisionImpact struct {
	Scope         string                 `json:"scope"`
	Severity      ImpactSeverity         `json:"severity"`
	Duration      string                 `json:"duration"`
	Reversible    bool                   `json:"reversible"`
	AffectedParties []string             `json:"affected_parties"`
	BusinessImpact map[string]interface{} `json:"business_impact,omitempty"`
	TechnicalImpact map[string]interface{} `json:"technical_impact,omitempty"`
	ComplianceImpact map[string]interface{} `json:"compliance_impact,omitempty"`
}

type ImpactSeverity int32

const (
	ImpactSeverity_MINIMAL     ImpactSeverity = 0
	ImpactSeverity_MINOR       ImpactSeverity = 1
	ImpactSeverity_MODERATE    ImpactSeverity = 2
	ImpactSeverity_MAJOR       ImpactSeverity = 3
	ImpactSeverity_CRITICAL    ImpactSeverity = 4
)

type DecisionAuditEntry struct {
	Timestamp   time.Time              `json:"timestamp"`
	Actor       string                 `json:"actor"`
	Action      string                 `json:"action"`
	Details     map[string]interface{} `json:"details,omitempty"`
	PreviousState map[string]interface{} `json:"previous_state,omitempty"`
	NewState    map[string]interface{} `json:"new_state,omitempty"`
}

type MonitoringConfiguration struct {
	Enabled        bool                   `json:"enabled"`
	MetricsCollection bool                `json:"metrics_collection"`
	AlertThresholds map[string]float64    `json:"alert_thresholds,omitempty"`
	ReportingSchedule ReportingSchedule   `json:"reporting_schedule"`
	Dashboards     []DashboardConfig      `json:"dashboards,omitempty"`
	Integrations   []MonitoringIntegration `json:"integrations,omitempty"`
}

type ReportingSchedule struct {
	Frequency   string    `json:"frequency"` // daily, weekly, monthly
	Time        string    `json:"time"`      // Time of day for reports
	Recipients  []string  `json:"recipients"`
	Format      string    `json:"format"`    // pdf, html, json
}

type DashboardConfig struct {
	DashboardID  string   `json:"dashboard_id"`
	Name         string   `json:"name"`
	Widgets      []string `json:"widgets"`
	RefreshRate  string   `json:"refresh_rate"`
	AccessRoles  []string `json:"access_roles"`
}

type MonitoringIntegration struct {
	Type         string                 `json:"type"`         // prometheus, grafana, etc.
	Endpoint     string                 `json:"endpoint"`
	Credentials  map[string]string      `json:"credentials,omitempty"`
	Configuration map[string]interface{} `json:"configuration,omitempty"`
}

type ComplianceReference struct {
	Regulation   string `json:"regulation"`  // GDPR, DPDP, etc.
	Article      string `json:"article"`     // Specific article/section
	Description  string `json:"description"`
	RequiredBy   string `json:"required_by"` // Date or event
	Verified     bool   `json:"verified"`
	VerifiedAt   *time.Time `json:"verified_at,omitempty"`
	VerifiedBy   string `json:"verified_by,omitempty"`
}

// Helper methods

// IsActive returns true if the policy is active
func (p *IdentityGovernancePolicy) IsActive() bool {
	now := time.Now()
	return p.Status == PolicyStatus_ACTIVE &&
		now.After(p.EffectiveDate) &&
		(p.ExpirationDate == nil || now.Before(*p.ExpirationDate))
}

// IsExpired returns true if the policy has expired
func (p *IdentityGovernancePolicy) IsExpired() bool {
	return p.ExpirationDate != nil && time.Now().After(*p.ExpirationDate)
}

// RequiresReview returns true if the policy requires review
func (p *IdentityGovernancePolicy) RequiresReview() bool {
	return p.NextReviewDate.Before(time.Now())
}

// String methods for better logging

func (pt GovernancePolicyType) String() string {
	switch pt {
	case GovernancePolicyType_ACCESS_CONTROL:
		return "access_control"
	case GovernancePolicyType_DATA_CLASSIFICATION:
		return "data_classification"
	case GovernancePolicyType_PRIVACY_PROTECTION:
		return "privacy_protection"
	case GovernancePolicyType_CREDENTIAL_LIFECYCLE:
		return "credential_lifecycle"
	case GovernancePolicyType_IDENTITY_VERIFICATION:
		return "identity_verification"
	case GovernancePolicyType_AUDIT_COMPLIANCE:
		return "audit_compliance"
	case GovernancePolicyType_RISK_MANAGEMENT:
		return "risk_management"
	case GovernancePolicyType_FEDERATION_TRUST:
		return "federation_trust"
	case GovernancePolicyType_CONSENT_MANAGEMENT:
		return "consent_management"
	case GovernancePolicyType_RETENTION_DISPOSAL:
		return "retention_disposal"
	case GovernancePolicyType_INCIDENT_RESPONSE:
		return "incident_response"
	case GovernancePolicyType_CUSTOM:
		return "custom"
	default:
		return "unknown"
	}
}

func (rt GovernanceRuleType) String() string {
	switch rt {
	case GovernanceRuleType_AUTHENTICATION:
		return "authentication"
	case GovernanceRuleType_AUTHORIZATION:
		return "authorization"
	case GovernanceRuleType_VALIDATION:
		return "validation"
	case GovernanceRuleType_ENCRYPTION:
		return "encryption"
	case GovernanceRuleType_AUDIT_LOGGING:
		return "audit_logging"
	case GovernanceRuleType_DATA_MASKING:
		return "data_masking"
	case GovernanceRuleType_RATE_LIMITING:
		return "rate_limiting"
	case GovernanceRuleType_ANOMALY_DETECTION:
		return "anomaly_detection"
	case GovernanceRuleType_LIFECYCLE_MANAGEMENT:
		return "lifecycle_management"
	case GovernanceRuleType_COMPLIANCE_CHECK:
		return "compliance_check"
	case GovernanceRuleType_NOTIFICATION:
		return "notification"
	case GovernanceRuleType_REMEDIATION:
		return "remediation"
	case GovernanceRuleType_CUSTOM:
		return "custom"
	default:
		return "unknown"
	}
}

func (at GovernanceActionType) String() string {
	switch at {
	case GovernanceActionType_ALLOW:
		return "allow"
	case GovernanceActionType_DENY:
		return "deny"
	case GovernanceActionType_REQUIRE_APPROVAL:
		return "require_approval"
	case GovernanceActionType_LOG_WARNING:
		return "log_warning"
	case GovernanceActionType_LOG_ERROR:
		return "log_error"
	case GovernanceActionType_SEND_NOTIFICATION:
		return "send_notification"
	case GovernanceActionType_QUARANTINE:
		return "quarantine"
	case GovernanceActionType_MASK_DATA:
		return "mask_data"
	case GovernanceActionType_ENCRYPT_DATA:
		return "encrypt_data"
	case GovernanceActionType_RATE_LIMIT:
		return "rate_limit"
	case GovernanceActionType_SUSPEND_IDENTITY:
		return "suspend_identity"
	case GovernanceActionType_REVOKE_CREDENTIAL:
		return "revoke_credential"
	case GovernanceActionType_TRIGGER_AUDIT:
		return "trigger_audit"
	case GovernanceActionType_INITIATE_WORKFLOW:
		return "initiate_workflow"
	case GovernanceActionType_CUSTOM:
		return "custom"
	default:
		return "unknown"
	}
}

// Error definitions for governance
var (
	ErrGovernancePolicyNotFound     = sdkerrors.Register(ModuleName, 7001, "governance policy not found")
	ErrGovernanceRoleNotFound       = sdkerrors.Register(ModuleName, 7002, "governance role not found")
	ErrGovernanceWorkflowNotFound   = sdkerrors.Register(ModuleName, 7003, "governance workflow not found")
	ErrGovernanceDecisionNotFound   = sdkerrors.Register(ModuleName, 7004, "governance decision not found")
	ErrInsufficientGovernanceRights = sdkerrors.Register(ModuleName, 7005, "insufficient governance rights")
	ErrPolicyValidationFailed       = sdkerrors.Register(ModuleName, 7006, "policy validation failed")
	ErrWorkflowExecutionFailed      = sdkerrors.Register(ModuleName, 7007, "workflow execution failed")
	ErrApprovalRequired             = sdkerrors.Register(ModuleName, 7008, "approval required for operation")
	ErrGovernanceRuleViolation      = sdkerrors.Register(ModuleName, 7009, "governance rule violation")
	ErrPolicyConflict               = sdkerrors.Register(ModuleName, 7010, "policy conflict detected")
	ErrInvalidGovernanceConfiguration = sdkerrors.Register(ModuleName, 7011, "invalid governance configuration")
	ErrDecisionAppealExpired        = sdkerrors.Register(ModuleName, 7012, "decision appeal period has expired")
	ErrWorkflowTimeout              = sdkerrors.Register(ModuleName, 7013, "workflow execution timeout")
	ErrRoleAssignmentFailed         = sdkerrors.Register(ModuleName, 7014, "role assignment failed")
	ErrDelegationNotAllowed         = sdkerrors.Register(ModuleName, 7015, "delegation not allowed for this role")
)