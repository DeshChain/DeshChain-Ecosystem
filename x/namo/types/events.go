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