package types

// Tax module event types
const (
	EventTypeTaxCollected = "tax_collected"
	EventTypeTaxDistributed = "tax_distributed"
	EventTypeTaxRefunded = "tax_refunded"
	EventTypeTaxExempted = "tax_exempted"
	EventTypeSustainableFeeUpdate = "sustainable_fee_update"
	
	AttributeKeyAmount = "amount"
	AttributeKeyFrom = "from"
	AttributeKeyTo = "to"
	AttributeKeyTaxRate = "tax_rate"
	AttributeKeyRecipient = "recipient"
	AttributeKeyReason = "reason"
)