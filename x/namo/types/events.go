package types

// NAMO module event types
const (
	EventTypeBurnTokens            = "burn_tokens"
	EventTypeCreateVestingSchedule = "create_vesting_schedule"
	EventTypeClaimVestedTokens     = "claim_vested_tokens"
	EventTypeUpdateParams          = "update_params"
	EventTypeDistributeTokens      = "distribute_tokens"
	EventTypeInitialDistribution   = "initial_distribution"
)

// NAMO module event attribute keys
const (
	AttributeKeyRecipient       = "recipient"
	AttributeKeySender          = "sender"
	AttributeKeyAmount          = "amount"
	AttributeKeyBurnedAmount    = "burned_amount"
	AttributeKeyClaimedAmount   = "claimed_amount"
	AttributeKeyScheduleID      = "schedule_id"
	AttributeKeyVestingPeriod   = "vesting_period"
	AttributeKeyCliffPeriod     = "cliff_period"
	AttributeKeyTotalAmount     = "total_amount"
	AttributeKeyEventType       = "event_type"
	AttributeKeyDistributionType = "distribution_type"
	AttributeKeyAuthority       = "authority"
	AttributeKeyTokenDenom      = "token_denom"
	AttributeKeyEnableVesting   = "enable_vesting"
	AttributeKeyEnableBurning   = "enable_burning"
	AttributeKeyMinBurnAmount   = "min_burn_amount"
)

// NAMO module event attribute values
const (
	AttributeValueCategory = ModuleName
)