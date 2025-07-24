package types

// Tax module event types
const (
	EventTypeTaxCollected = "tax_collected"
	EventTypeTaxDistributed = "tax_distributed"
	EventTypeTaxRefunded = "tax_refunded"
	EventTypeTaxExempted = "tax_exempted"
	EventTypeSustainableFeeUpdate = "sustainable_fee_update"
	EventTypeNAMOSwap = "namo_swap"
	EventTypeNAMOBurn = "namo_burn"
	
	AttributeKeyAmount = "amount"
	AttributeKeyFrom = "from"
	AttributeKeyTo = "to"
	AttributeKeyTaxRate = "tax_rate"
	AttributeKeyRecipient = "recipient"
	AttributeKeyReason = "reason"
	AttributeKeyUser = "user"
	AttributeKeyFromToken = "from_token"
	AttributeKeyToToken = "to_token"
	AttributeKeySwapRate = "swap_rate"
	AttributeKeyBurnAmount = "burn_amount"
)