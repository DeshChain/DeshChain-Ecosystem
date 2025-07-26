package types

import (
	"time"
)

// Identity Audit and Compliance Types

// AuditEvent represents a single audit event in the identity system
type AuditEvent struct {
	EventID       string                 `json:"event_id"`
	Timestamp     time.Time              `json:"timestamp"`
	EventType     AuditEventType         `json:"event_type"`
	Actor         string                 `json:"actor"`          // Address of the actor performing the action
	Subject       string                 `json:"subject"`        // Address/DID of the subject being acted upon
	Resource      string                 `json:"resource"`       // Resource identifier (credential ID, DID, etc.)
	Action        string                 `json:"action"`         // Specific action performed
	Outcome       AuditOutcome           `json:"outcome"`        // Success, failure, etc.
	Severity      AuditSeverity          `json:"severity"`       // Low, medium, high, critical
	Description   string                 `json:"description"`    // Human-readable description
	TechnicalDetails map[string]interface{} `json:"technical_details,omitempty"`
	ComplianceFlags []ComplianceFlag      `json:"compliance_flags,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	IPAddress     string                 `json:"ip_address,omitempty"`
	UserAgent     string                 `json:"user_agent,omitempty"`
	SessionID     string                 `json:"session_id,omitempty"`
	ModuleSource  string                 `json:"module_source"`  // Which module generated this event
	ChainHeight   int64                  `json:"chain_height"`   // Block height when event occurred
	TxHash        string                 `json:"tx_hash,omitempty"` // Transaction hash if applicable
}

// AuditEventType categorizes different types of audit events
type AuditEventType int32

const (
	AuditEventType_IDENTITY_CREATED     AuditEventType = 0
	AuditEventType_IDENTITY_UPDATED     AuditEventType = 1
	AuditEventType_IDENTITY_DELETED     AuditEventType = 2
	AuditEventType_IDENTITY_ACCESSED    AuditEventType = 3
	AuditEventType_CREDENTIAL_ISSUED    AuditEventType = 4
	AuditEventType_CREDENTIAL_VERIFIED  AuditEventType = 5
	AuditEventType_CREDENTIAL_REVOKED   AuditEventType = 6
	AuditEventType_CREDENTIAL_ACCESSED  AuditEventType = 7
	AuditEventType_CONSENT_GIVEN        AuditEventType = 8
	AuditEventType_CONSENT_WITHDRAWN    AuditEventType = 9
	AuditEventType_CONSENT_ACCESSED     AuditEventType = 10
	AuditEventType_DATA_SHARED          AuditEventType = 11
	AuditEventType_DATA_ACCESS_REQUESTED AuditEventType = 12
	AuditEventType_DATA_ACCESS_DENIED   AuditEventType = 13
	AuditEventType_PRIVACY_SETTINGS_CHANGED AuditEventType = 14
	AuditEventType_DID_CREATED          AuditEventType = 15
	AuditEventType_DID_UPDATED          AuditEventType = 16
	AuditEventType_DID_DEACTIVATED      AuditEventType = 17
	AuditEventType_BIOMETRIC_ENROLLED   AuditEventType = 18
	AuditEventType_BIOMETRIC_VERIFIED   AuditEventType = 19
	AuditEventType_KYC_INITIATED        AuditEventType = 20
	AuditEventType_KYC_COMPLETED        AuditEventType = 21
	AuditEventType_RECOVERY_INITIATED   AuditEventType = 22
	AuditEventType_RECOVERY_COMPLETED   AuditEventType = 23
	AuditEventType_SUSPICIOUS_ACTIVITY  AuditEventType = 24
	AuditEventType_COMPLIANCE_VIOLATION AuditEventType = 25
	AuditEventType_SYSTEM_ERROR         AuditEventType = 26
	AuditEventType_ADMIN_ACTION         AuditEventType = 27
	AuditEventType_EXPORT_REQUEST       AuditEventType = 28
	AuditEventType_DELETION_REQUEST     AuditEventType = 29
)

// AuditOutcome represents the result of an audited action
type AuditOutcome int32

const (
	AuditOutcome_SUCCESS           AuditOutcome = 0
	AuditOutcome_FAILURE           AuditOutcome = 1
	AuditOutcome_PARTIAL_SUCCESS   AuditOutcome = 2
	AuditOutcome_DENIED            AuditOutcome = 3
	AuditOutcome_ERROR             AuditOutcome = 4
	AuditOutcome_TIMEOUT           AuditOutcome = 5
	AuditOutcome_CANCELLED         AuditOutcome = 6
)

// AuditSeverity indicates the importance level of an audit event
type AuditSeverity int32

const (
	AuditSeverity_LOW      AuditSeverity = 0
	AuditSeverity_MEDIUM   AuditSeverity = 1
	AuditSeverity_HIGH     AuditSeverity = 2
	AuditSeverity_CRITICAL AuditSeverity = 3
)

// ComplianceFlag indicates specific compliance requirements related to an event
type ComplianceFlag struct {
	Regulation  string                 `json:"regulation"`  // GDPR, DPDP, CCPA, etc.
	Requirement string                 `json:"requirement"` // Specific requirement (e.g., "Article 17 - Right to erasure")
	Status      ComplianceStatus       `json:"status"`      // Compliant, non-compliant, pending
	Notes       string                 `json:"notes,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// ComplianceStatus indicates compliance state
type ComplianceStatus int32

const (
	ComplianceStatus_COMPLIANT     ComplianceStatus = 0
	ComplianceStatus_NON_COMPLIANT ComplianceStatus = 1
	ComplianceStatus_PENDING       ComplianceStatus = 2
	ComplianceStatus_EXEMPT        ComplianceStatus = 3
	ComplianceStatus_UNKNOWN       ComplianceStatus = 4
)

// ComplianceReport represents a comprehensive compliance report
type ComplianceReport struct {
	ReportID          string                    `json:"report_id"`
	GeneratedAt       time.Time                 `json:"generated_at"`
	GeneratedBy       string                    `json:"generated_by"`
	ReportType        ComplianceReportType      `json:"report_type"`
	TimeRange         AuditTimeRange            `json:"time_range"`
	Scope             ComplianceScope           `json:"scope"`
	Regulations       []string                  `json:"regulations"`
	Summary           ComplianceReportSummary   `json:"summary"`
	Findings          []ComplianceFinding       `json:"findings"`
	Recommendations   []ComplianceRecommendation `json:"recommendations"`
	DataSubjects      int64                     `json:"data_subjects"`     // Number of unique data subjects
	TotalEvents       int64                     `json:"total_events"`      // Total audit events analyzed
	ComplianceScore   float64                   `json:"compliance_score"`  // Overall compliance score (0-100)
	RiskLevel         ComplianceRiskLevel       `json:"risk_level"`
	NextReviewDate    time.Time                 `json:"next_review_date"`
	CertificationInfo *CertificationInfo        `json:"certification_info,omitempty"`
	Metadata          map[string]interface{}    `json:"metadata,omitempty"`
}

// ComplianceReportType categorizes different types of compliance reports
type ComplianceReportType int32

const (
	ComplianceReportType_GDPR_COMPLIANCE    ComplianceReportType = 0
	ComplianceReportType_DPDP_COMPLIANCE    ComplianceReportType = 1
	ComplianceReportType_CCPA_COMPLIANCE    ComplianceReportType = 2
	ComplianceReportType_GENERAL_AUDIT      ComplianceReportType = 3
	ComplianceReportType_SECURITY_AUDIT     ComplianceReportType = 4
	ComplianceReportType_DATA_BREACH_REPORT ComplianceReportType = 5
	ComplianceReportType_PERIODIC_REVIEW    ComplianceReportType = 6
	ComplianceReportType_INCIDENT_REPORT    ComplianceReportType = 7
	ComplianceReportType_RISK_ASSESSMENT    ComplianceReportType = 8
)

// AuditTimeRange defines the time period for audit analysis
type AuditTimeRange struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// ComplianceScope defines the scope of compliance analysis
type ComplianceScope struct {
	IncludeModules    []string `json:"include_modules,omitempty"`
	ExcludeModules    []string `json:"exclude_modules,omitempty"`
	IncludeEventTypes []AuditEventType `json:"include_event_types,omitempty"`
	ExcludeEventTypes []AuditEventType `json:"exclude_event_types,omitempty"`
	DataSubjects      []string `json:"data_subjects,omitempty"`      // Specific users to include
	GeographicScope   []string `json:"geographic_scope,omitempty"`   // Countries/regions
	OrganizationalScope []string `json:"organizational_scope,omitempty"` // Departments/business units
}

// ComplianceReportSummary provides high-level statistics
type ComplianceReportSummary struct {
	TotalCompliantEvents    int64   `json:"total_compliant_events"`
	TotalNonCompliantEvents int64   `json:"total_non_compliant_events"`
	CompliancePercentage    float64 `json:"compliance_percentage"`
	CriticalFindings        int64   `json:"critical_findings"`
	HighRiskFindings        int64   `json:"high_risk_findings"`
	MediumRiskFindings      int64   `json:"medium_risk_findings"`
	LowRiskFindings         int64   `json:"low_risk_findings"`
	ResolvedFindings        int64   `json:"resolved_findings"`
	OutstandingFindings     int64   `json:"outstanding_findings"`
	DataProcessingActivities int64  `json:"data_processing_activities"`
	ConsentWithdrawals      int64   `json:"consent_withdrawals"`
	DataAccessRequests      int64   `json:"data_access_requests"`
	DataDeletionRequests    int64   `json:"data_deletion_requests"`
}

// ComplianceFinding represents a specific compliance issue or observation
type ComplianceFinding struct {
	FindingID     string                 `json:"finding_id"`
	Regulation    string                 `json:"regulation"`
	Requirement   string                 `json:"requirement"`
	Severity      AuditSeverity          `json:"severity"`
	Status        FindingStatus          `json:"status"`
	Description   string                 `json:"description"`
	Evidence      []string               `json:"evidence"`         // Event IDs supporting this finding
	Impact        string                 `json:"impact"`
	Remediation   string                 `json:"remediation"`
	DueDate       *time.Time             `json:"due_date,omitempty"`
	AssignedTo    string                 `json:"assigned_to,omitempty"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
	ResolvedAt    *time.Time             `json:"resolved_at,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// FindingStatus indicates the current state of a compliance finding
type FindingStatus int32

const (
	FindingStatus_OPEN        FindingStatus = 0
	FindingStatus_IN_PROGRESS FindingStatus = 1
	FindingStatus_RESOLVED    FindingStatus = 2
	FindingStatus_DISMISSED   FindingStatus = 3
	FindingStatus_ESCALATED   FindingStatus = 4
)

// ComplianceRecommendation provides actionable guidance
type ComplianceRecommendation struct {
	RecommendationID string                 `json:"recommendation_id"`
	Title            string                 `json:"title"`
	Description      string                 `json:"description"`
	Priority         RecommendationPriority `json:"priority"`
	Category         RecommendationCategory `json:"category"`
	EstimatedEffort  string                 `json:"estimated_effort"`
	ExpectedBenefit  string                 `json:"expected_benefit"`
	Implementation   string                 `json:"implementation"`
	RelatedFindings  []string               `json:"related_findings"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// RecommendationPriority indicates urgency of recommendations
type RecommendationPriority int32

const (
	RecommendationPriority_LOW      RecommendationPriority = 0
	RecommendationPriority_MEDIUM   RecommendationPriority = 1
	RecommendationPriority_HIGH     RecommendationPriority = 2
	RecommendationPriority_CRITICAL RecommendationPriority = 3
)

// RecommendationCategory categorizes types of recommendations
type RecommendationCategory int32

const (
	RecommendationCategory_POLICY       RecommendationCategory = 0
	RecommendationCategory_TECHNICAL    RecommendationCategory = 1
	RecommendationCategory_PROCEDURAL   RecommendationCategory = 2
	RecommendationCategory_TRAINING     RecommendationCategory = 3
	RecommendationCategory_MONITORING   RecommendationCategory = 4
	RecommendationCategory_GOVERNANCE   RecommendationCategory = 5
)

// ComplianceRiskLevel indicates overall risk assessment
type ComplianceRiskLevel int32

const (
	ComplianceRiskLevel_LOW      ComplianceRiskLevel = 0
	ComplianceRiskLevel_MEDIUM   ComplianceRiskLevel = 1
	ComplianceRiskLevel_HIGH     ComplianceRiskLevel = 2
	ComplianceRiskLevel_CRITICAL ComplianceRiskLevel = 3
)

// CertificationInfo contains information about compliance certifications
type CertificationInfo struct {
	CertificationType string    `json:"certification_type"` // ISO 27001, SOC 2, etc.
	CertificationBody string    `json:"certification_body"`
	CertificationDate time.Time `json:"certification_date"`
	ExpiryDate        time.Time `json:"expiry_date"`
	CertificateNumber string    `json:"certificate_number"`
	Scope             string    `json:"scope"`
	Status            string    `json:"status"`
}

// DataSubjectRequest represents requests from data subjects (GDPR Article 12-22)
type DataSubjectRequest struct {
	RequestID      string                    `json:"request_id"`
	RequestType    DataSubjectRequestType    `json:"request_type"`
	RequestStatus  DataSubjectRequestStatus  `json:"request_status"`
	DataSubject    string                    `json:"data_subject"`    // Address/DID of requestor
	RequestedBy    string                    `json:"requested_by"`    // May be different if via representative
	RequestDate    time.Time                 `json:"request_date"`
	DueDate        time.Time                 `json:"due_date"`        // Legal deadline (e.g., 30 days for GDPR)
	ProcessedDate  *time.Time                `json:"processed_date,omitempty"`
	CompletedDate  *time.Time                `json:"completed_date,omitempty"`
	Description    string                    `json:"description"`
	RequestDetails map[string]interface{}    `json:"request_details"`
	ProcessingNotes string                   `json:"processing_notes,omitempty"`
	ResponseData   *DataSubjectResponse      `json:"response_data,omitempty"`
	LegalBasis     string                    `json:"legal_basis,omitempty"`
	Regulation     string                    `json:"regulation"`      // GDPR, DPDP, etc.
	Priority       DataSubjectRequestPriority `json:"priority"`
	AssignedTo     string                    `json:"assigned_to,omitempty"`
	Metadata       map[string]interface{}    `json:"metadata,omitempty"`
}

// DataSubjectRequestType categorizes different types of data subject requests
type DataSubjectRequestType int32

const (
	DataSubjectRequestType_ACCESS         DataSubjectRequestType = 0  // Article 15
	DataSubjectRequestType_RECTIFICATION  DataSubjectRequestType = 1  // Article 16
	DataSubjectRequestType_ERASURE        DataSubjectRequestType = 2  // Article 17 (Right to be forgotten)
	DataSubjectRequestType_RESTRICT       DataSubjectRequestType = 3  // Article 18
	DataSubjectRequestType_PORTABILITY    DataSubjectRequestType = 4  // Article 20
	DataSubjectRequestType_OBJECT         DataSubjectRequestType = 5  // Article 21
	DataSubjectRequestType_WITHDRAW_CONSENT DataSubjectRequestType = 6 // Article 7(3)
	DataSubjectRequestType_COMPLAINT      DataSubjectRequestType = 7
	DataSubjectRequestType_INFORMATION    DataSubjectRequestType = 8  // Article 13-14
	DataSubjectRequestType_STOP_PROCESSING DataSubjectRequestType = 9
)

// DataSubjectRequestStatus tracks the processing status of requests
type DataSubjectRequestStatus int32

const (
	DataSubjectRequestStatus_RECEIVED     DataSubjectRequestStatus = 0
	DataSubjectRequestStatus_ACKNOWLEDGED DataSubjectRequestStatus = 1
	DataSubjectRequestStatus_IN_REVIEW    DataSubjectRequestStatus = 2
	DataSubjectRequestStatus_IN_PROGRESS  DataSubjectRequestStatus = 3
	DataSubjectRequestStatus_COMPLETED    DataSubjectRequestStatus = 4
	DataSubjectRequestStatus_REJECTED     DataSubjectRequestStatus = 5
	DataSubjectRequestStatus_ESCALATED    DataSubjectRequestStatus = 6
	DataSubjectRequestStatus_EXPIRED      DataSubjectRequestStatus = 7
)

// DataSubjectRequestPriority indicates urgency of data subject requests
type DataSubjectRequestPriority int32

const (
	DataSubjectRequestPriority_STANDARD DataSubjectRequestPriority = 0
	DataSubjectRequestPriority_URGENT   DataSubjectRequestPriority = 1
	DataSubjectRequestPriority_CRITICAL DataSubjectRequestPriority = 2
)

// DataSubjectResponse contains the response to a data subject request
type DataSubjectResponse struct {
	ResponseID    string                 `json:"response_id"`
	ResponseType  string                 `json:"response_type"`
	ResponseData  map[string]interface{} `json:"response_data,omitempty"`
	ExportFormat  string                 `json:"export_format,omitempty"`  // JSON, CSV, PDF
	DeliveryMethod string                `json:"delivery_method,omitempty"` // Email, secure download, etc.
	DeliveredAt   *time.Time             `json:"delivered_at,omitempty"`
	ExpiresAt     *time.Time             `json:"expires_at,omitempty"`
	DownloadCount int64                  `json:"download_count"`
	Checksum      string                 `json:"checksum,omitempty"`
	EncryptionKey string                 `json:"encryption_key,omitempty"`
}

// PrivacyImpactAssessment represents a Data Protection Impact Assessment (DPIA)
type PrivacyImpactAssessment struct {
	AssessmentID       string                 `json:"assessment_id"`
	Title              string                 `json:"title"`
	Description        string                 `json:"description"`
	ProcessingActivity string                 `json:"processing_activity"`
	DataController     string                 `json:"data_controller"`
	DataProcessor      string                 `json:"data_processor,omitempty"`
	LegalBasis         []string               `json:"legal_basis"`
	DataCategories     []string               `json:"data_categories"`
	DataSubjects       []string               `json:"data_subjects"`
	ProcessingPurposes []string               `json:"processing_purposes"`
	RetentionPeriod    string                 `json:"retention_period"`
	ThirdPartySharing  []ThirdPartySharing    `json:"third_party_sharing,omitempty"`
	TechnicalMeasures  []string               `json:"technical_measures"`
	OrganizationalMeasures []string           `json:"organizational_measures"`
	RiskAssessment     PrivacyRiskAssessment  `json:"risk_assessment"`
	Mitigation         []MitigationMeasure    `json:"mitigation"`
	ReviewDate         time.Time              `json:"review_date"`
	Status             PIAStatus              `json:"status"`
	CreatedBy          string                 `json:"created_by"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
	ApprovedBy         string                 `json:"approved_by,omitempty"`
	ApprovedAt         *time.Time             `json:"approved_at,omitempty"`
	Metadata           map[string]interface{} `json:"metadata,omitempty"`
}

// ThirdPartySharing describes data sharing with third parties
type ThirdPartySharing struct {
	PartyName        string    `json:"party_name"`
	PartyType        string    `json:"party_type"`        // Processor, controller, etc.
	Purpose          string    `json:"purpose"`
	DataCategories   []string  `json:"data_categories"`
	LegalBasis       string    `json:"legal_basis"`
	Safeguards       []string  `json:"safeguards"`
	ContractualTerms string    `json:"contractual_terms"`
	Country          string    `json:"country"`
	AdequacyDecision bool      `json:"adequacy_decision"`
	StartDate        time.Time `json:"start_date"`
	EndDate          *time.Time `json:"end_date,omitempty"`
}

// PrivacyRiskAssessment contains risk analysis for privacy impact
type PrivacyRiskAssessment struct {
	OverallRiskLevel  PrivacyRiskLevel    `json:"overall_risk_level"`
	IdentifiedRisks   []PrivacyRisk       `json:"identified_risks"`
	ResidualRiskLevel PrivacyRiskLevel    `json:"residual_risk_level"`
	ReviewDate        time.Time           `json:"review_date"`
	RiskMatrix        map[string]string   `json:"risk_matrix,omitempty"`
	Assumptions       []string            `json:"assumptions,omitempty"`
}

// PrivacyRisk represents a specific privacy risk
type PrivacyRisk struct {
	RiskID          string           `json:"risk_id"`
	RiskType        string           `json:"risk_type"`
	Description     string           `json:"description"`
	Likelihood      RiskLevel        `json:"likelihood"`
	Impact          RiskLevel        `json:"impact"`
	OverallRisk     RiskLevel        `json:"overall_risk"`
	AffectedRights  []string         `json:"affected_rights"`
	Consequences    []string         `json:"consequences"`
	ExistingControls []string        `json:"existing_controls"`
}

// RiskLevel represents risk severity levels
type RiskLevel int32

const (
	RiskLevel_VERY_LOW  RiskLevel = 0
	RiskLevel_LOW       RiskLevel = 1
	RiskLevel_MEDIUM    RiskLevel = 2
	RiskLevel_HIGH      RiskLevel = 3
	RiskLevel_VERY_HIGH RiskLevel = 4
)

// PrivacyRiskLevel extends RiskLevel for privacy-specific context
type PrivacyRiskLevel = RiskLevel

// MitigationMeasure describes measures to reduce privacy risks
type MitigationMeasure struct {
	MeasureID      string                 `json:"measure_id"`
	Type           MitigationType         `json:"type"`
	Description    string                 `json:"description"`
	Implementation string                 `json:"implementation"`
	Effectiveness  EffectivenessLevel     `json:"effectiveness"`
	Cost           CostLevel              `json:"cost"`
	Timeline       string                 `json:"timeline"`
	ResponsibleParty string               `json:"responsible_party"`
	Status         MitigationStatus       `json:"status"`
	TargetRisks    []string               `json:"target_risks"`
	Metrics        []string               `json:"metrics"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// MitigationType categorizes different types of mitigation measures
type MitigationType int32

const (
	MitigationType_TECHNICAL        MitigationType = 0
	MitigationType_ORGANIZATIONAL   MitigationType = 1
	MitigationType_LEGAL           MitigationType = 2
	MitigationType_PROCEDURAL      MitigationType = 3
	MitigationType_TRAINING        MitigationType = 4
	MitigationType_MONITORING      MitigationType = 5
)

// EffectivenessLevel indicates how effective a mitigation measure is
type EffectivenessLevel int32

const (
	EffectivenessLevel_LOW    EffectivenessLevel = 0
	EffectivenessLevel_MEDIUM EffectivenessLevel = 1
	EffectivenessLevel_HIGH   EffectivenessLevel = 2
	EffectivenessLevel_VERY_HIGH EffectivenessLevel = 3
)

// CostLevel indicates the cost of implementing a mitigation measure
type CostLevel int32

const (
	CostLevel_LOW    CostLevel = 0
	CostLevel_MEDIUM CostLevel = 1
	CostLevel_HIGH   CostLevel = 2
	CostLevel_VERY_HIGH CostLevel = 3
)

// MitigationStatus tracks implementation status of mitigation measures
type MitigationStatus int32

const (
	MitigationStatus_PLANNED      MitigationStatus = 0
	MitigationStatus_IN_PROGRESS  MitigationStatus = 1
	MitigationStatus_IMPLEMENTED  MitigationStatus = 2
	MitigationStatus_VERIFIED     MitigationStatus = 3
	MitigationStatus_INEFFECTIVE  MitigationStatus = 4
)

// PIAStatus represents the status of a Privacy Impact Assessment
type PIAStatus int32

const (
	PIAStatus_DRAFT      PIAStatus = 0
	PIAStatus_REVIEW     PIAStatus = 1
	PIAStatus_APPROVED   PIAStatus = 2
	PIAStatus_REJECTED   PIAStatus = 3
	PIAStatus_SUPERSEDED PIAStatus = 4
)

// Helper methods

// String returns string representation of AuditEventType
func (aet AuditEventType) String() string {
	switch aet {
	case AuditEventType_IDENTITY_CREATED:
		return "identity_created"
	case AuditEventType_IDENTITY_UPDATED:
		return "identity_updated"
	case AuditEventType_IDENTITY_DELETED:
		return "identity_deleted"
	case AuditEventType_CREDENTIAL_ISSUED:
		return "credential_issued"
	case AuditEventType_CREDENTIAL_VERIFIED:
		return "credential_verified"
	case AuditEventType_CREDENTIAL_REVOKED:
		return "credential_revoked"
	case AuditEventType_CONSENT_GIVEN:
		return "consent_given"
	case AuditEventType_CONSENT_WITHDRAWN:
		return "consent_withdrawn"
	case AuditEventType_DATA_SHARED:
		return "data_shared"
	case AuditEventType_DATA_ACCESS_REQUESTED:
		return "data_access_requested"
	case AuditEventType_DATA_ACCESS_DENIED:
		return "data_access_denied"
	case AuditEventType_PRIVACY_SETTINGS_CHANGED:
		return "privacy_settings_changed"
	case AuditEventType_SUSPICIOUS_ACTIVITY:
		return "suspicious_activity"
	case AuditEventType_COMPLIANCE_VIOLATION:
		return "compliance_violation"
	default:
		return "unknown"
	}
}

// IsHighPriority returns true if the event type requires immediate attention
func (aet AuditEventType) IsHighPriority() bool {
	switch aet {
	case AuditEventType_SUSPICIOUS_ACTIVITY, AuditEventType_COMPLIANCE_VIOLATION,
		 AuditEventType_SYSTEM_ERROR, AuditEventType_DATA_ACCESS_DENIED:
		return true
	default:
		return false
	}
}

// RequiresNotification returns true if the event should trigger notifications
func (aet AuditEventType) RequiresNotification() bool {
	switch aet {
	case AuditEventType_CONSENT_WITHDRAWN, AuditEventType_CREDENTIAL_REVOKED,
		 AuditEventType_SUSPICIOUS_ACTIVITY, AuditEventType_COMPLIANCE_VIOLATION:
		return true
	default:
		return false
	}
}

// GetRetentionPeriod returns the recommended retention period for this event type
func (aet AuditEventType) GetRetentionPeriod() time.Duration {
	switch aet {
	case AuditEventType_COMPLIANCE_VIOLATION, AuditEventType_SUSPICIOUS_ACTIVITY:
		return 7 * 365 * 24 * time.Hour // 7 years
	case AuditEventType_CONSENT_GIVEN, AuditEventType_CONSENT_WITHDRAWN:
		return 6 * 365 * 24 * time.Hour // 6 years
	default:
		return 3 * 365 * 24 * time.Hour // 3 years
	}
}

// String returns string representation of AuditOutcome
func (ao AuditOutcome) String() string {
	switch ao {
	case AuditOutcome_SUCCESS:
		return "success"
	case AuditOutcome_FAILURE:
		return "failure"
	case AuditOutcome_PARTIAL_SUCCESS:
		return "partial_success"
	case AuditOutcome_DENIED:
		return "denied"
	case AuditOutcome_ERROR:
		return "error"
	case AuditOutcome_TIMEOUT:
		return "timeout"
	case AuditOutcome_CANCELLED:
		return "cancelled"
	default:
		return "unknown"
	}
}

// String returns string representation of AuditSeverity
func (as AuditSeverity) String() string {
	switch as {
	case AuditSeverity_LOW:
		return "low"
	case AuditSeverity_MEDIUM:
		return "medium"
	case AuditSeverity_HIGH:
		return "high"
	case AuditSeverity_CRITICAL:
		return "critical"
	default:
		return "unknown"
	}
}

// IsExpired returns true if the data subject request is past its due date
func (dsr *DataSubjectRequest) IsExpired() bool {
	return time.Now().After(dsr.DueDate)
}

// IsOverdue returns true if the request is overdue
func (dsr *DataSubjectRequest) IsOverdue() bool {
	return dsr.IsExpired() && dsr.RequestStatus != DataSubjectRequestStatus_COMPLETED
}

// GetProcessingTime returns how long the request has been in processing
func (dsr *DataSubjectRequest) GetProcessingTime() time.Duration {
	if dsr.CompletedDate != nil {
		return dsr.CompletedDate.Sub(dsr.RequestDate)
	}
	return time.Since(dsr.RequestDate)
}

// Error definitions for audit and compliance
var (
	ErrAuditEventNotFound        = sdkerrors.Register(ModuleName, 5001, "audit event not found")
	ErrComplianceReportNotFound  = sdkerrors.Register(ModuleName, 5002, "compliance report not found")
	ErrDataSubjectRequestNotFound = sdkerrors.Register(ModuleName, 5003, "data subject request not found")
	ErrInvalidAuditQuery         = sdkerrors.Register(ModuleName, 5004, "invalid audit query parameters")
	ErrAuditRetentionViolation   = sdkerrors.Register(ModuleName, 5005, "audit retention policy violation")
	ErrComplianceViolation       = sdkerrors.Register(ModuleName, 5006, "compliance violation detected")
	ErrPrivacyImpactAssessmentRequired = sdkerrors.Register(ModuleName, 5007, "privacy impact assessment required")
	ErrDataSubjectRequestExpired = sdkerrors.Register(ModuleName, 5008, "data subject request has expired")
	ErrInsufficientAuditPermissions = sdkerrors.Register(ModuleName, 5009, "insufficient permissions for audit operation")
)