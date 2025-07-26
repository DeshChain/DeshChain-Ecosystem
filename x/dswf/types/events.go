package types

// DSWF module event types
const (
	EventTypeAllocationProposed   = "allocation_proposed"
	EventTypeAllocationApproved   = "allocation_approved"
	EventTypeDisbursementExecuted = "disbursement_executed"
	EventTypeStrategyUpdated      = "strategy_updated"
	EventTypePortfolioRebalanced  = "portfolio_rebalanced"
	EventTypeMetricsSubmitted     = "metrics_submitted"
	EventTypeMonthlyReport        = "monthly_report_generated"
	
	AttributeKeyAllocationID      = "allocation_id"
	AttributeKeyProposer          = "proposer"
	AttributeKeyApprover          = "approver"
	AttributeKeyAmount            = "amount"
	AttributeKeyCategory          = "category"
	AttributeKeyRecipient         = "recipient"
	AttributeKeyStatus            = "status"
	AttributeKeyApproved          = "approved"
	AttributeKeyDisbursementIndex = "disbursement_index"
	AttributeKeyAuthority         = "authority"
	AttributeKeySubmitter         = "submitter"
	AttributeKeyPeriod            = "period"
)