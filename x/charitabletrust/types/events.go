package types

// CharitableTrust module event types
const (
	EventTypeProposalCreated      = "proposal_created"
	EventTypeProposalVoted        = "proposal_voted"
	EventTypeProposalApproved     = "proposal_approved"
	EventTypeProposalRejected     = "proposal_rejected"
	EventTypeAllocationExecuted   = "allocation_executed"
	EventTypeImpactReportSubmitted = "impact_report_submitted"
	EventTypeImpactReportVerified = "impact_report_verified"
	EventTypeFraudAlertCreated   = "fraud_alert_created"
	EventTypeFraudInvestigated   = "fraud_investigated"
	EventTypeTrusteesUpdated     = "trustees_updated"
	EventTypeFundsDistributed    = "funds_distributed"
	
	AttributeKeyProposalID       = "proposal_id"
	AttributeKeyAllocationID     = "allocation_id"
	AttributeKeyReportID         = "report_id"
	AttributeKeyAlertID          = "alert_id"
	AttributeKeyProposer         = "proposer"
	AttributeKeyVoter            = "voter"
	AttributeKeyVote             = "vote"
	AttributeKeyAmount           = "amount"
	AttributeKeyCategory         = "category"
	AttributeKeyOrganizationID   = "organization_id"
	AttributeKeyOrganizationName = "organization_name"
	AttributeKeyStatus           = "status"
	AttributeKeySeverity         = "severity"
	AttributeKeyVerifier         = "verifier"
	AttributeKeyInvestigator     = "investigator"
	AttributeKeyBeneficiaries    = "beneficiaries"
	AttributeKeyImpactScore      = "impact_score"
)