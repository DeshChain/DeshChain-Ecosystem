package types

// Oracle module event types and attribute keys
const (
	// Event types
	EventTypePriceSubmission      = "price_submission"
	EventTypePriceUpdate          = "price_update"
	EventTypeExchangeRateUpdate   = "exchange_rate_update"
	EventTypeValidatorRegistered  = "oracle_validator_registered"
	EventTypeValidatorUpdated     = "oracle_validator_updated"
	EventTypeValidatorSlashed     = "oracle_validator_slashed"
	EventTypePriceAggregation     = "price_aggregation"
	EventTypePriceDeviation       = "price_deviation"

	// Attribute keys
	AttributeKeyValidator           = "validator"
	AttributeKeySymbol              = "symbol"
	AttributeKeyPrice               = "price"
	AttributeKeyOldPrice            = "old_price"
	AttributeKeyNewPrice            = "new_price"
	AttributeKeySource              = "source"
	AttributeKeyTimestamp           = "timestamp"
	AttributeKeyBase                = "base"
	AttributeKeyTarget              = "target"
	AttributeKeyRate                = "rate"
	AttributeKeyPower               = "power"
	AttributeKeyActive              = "active"
	AttributeKeySlashAmount         = "slash_amount"
	AttributeKeyReason              = "reason"
	AttributeKeyBlockHeight         = "block_height"
	AttributeKeyValidatorCount      = "validator_count"
	AttributeKeyMedianPrice         = "median_price"
	AttributeKeyMeanPrice           = "mean_price"
	AttributeKeyDeviation           = "deviation"
	AttributeKeyMaxDeviation        = "max_deviation"
	AttributeKeyDeviationExceeded   = "deviation_exceeded"
	AttributeKeyWindowStart         = "window_start"
	AttributeKeyWindowEnd           = "window_end"

	// Attribute values
	AttributeValueCategory = ModuleName
)