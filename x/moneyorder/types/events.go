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

// Money Order module event types
const (
	// Module events
	EventTypeMoneyOrder       = "money_order"
	EventTypeCreatePool       = "create_pool"
	EventTypeSwap             = "swap"
	EventTypeAddLiquidity     = "add_liquidity"
	EventTypeRemoveLiquidity  = "remove_liquidity"
	EventTypeJoinPool         = "join_pool"
	EventTypeExitPool         = "exit_pool"
	EventTypeJoinVillagePool  = "join_village_pool"
	EventTypeClaimRewards     = "claim_rewards"
	EventTypeUpdatePool       = "update_pool"
	
	// Money Order specific events
	EventTypeMoneyOrderCreated    = "money_order_created"
	EventTypeMoneyOrderCompleted  = "money_order_completed"
	EventTypeMoneyOrderCancelled  = "money_order_cancelled"
	EventTypeMoneyOrderExpired    = "money_order_expired"
	EventTypeReceiptGenerated     = "receipt_generated"
	
	// Pool events
	EventTypeFixedRatePoolCreated = "fixed_rate_pool_created"
	EventTypeVillagePoolCreated   = "village_pool_created"
	EventTypeAMMPoolCreated       = "amm_pool_created"
	EventTypePoolActivated        = "pool_activated"
	EventTypePoolDeactivated      = "pool_deactivated"
	
	// Village pool events
	EventTypeVillageMemberAdded   = "village_member_added"
	EventTypeVillageAchievement   = "village_achievement"
	EventTypeVillageVerified      = "village_verified"
	EventTypeTrustScoreUpdated    = "trust_score_updated"
	
	// Cultural events
	EventTypeFestivalBonusApplied = "festival_bonus_applied"
	EventTypeCulturalDiscount     = "cultural_discount_applied"
	EventTypeFestivalPeriodUpdated = "festival_period_updated"
	
	// Common attributes
	AttributeKeySender            = "sender"
	AttributeKeyReceiver          = "receiver"
	AttributeKeyAmount            = "amount"
	AttributeKeyPoolId            = "pool_id"
	AttributeKeyOrderId           = "order_id"
	AttributeKeyReferenceNumber   = "reference_number"
	
	// Money Order attributes
	AttributeKeyOrderType         = "order_type"
	AttributeKeyReceiverUPI       = "receiver_upi"
	AttributeKeyNote              = "note"
	AttributeKeyPostalCodeFrom    = "postal_code_from"
	AttributeKeyPostalCodeTo      = "postal_code_to"
	AttributeKeyExchangeRate      = "exchange_rate"
	AttributeKeyFees              = "fees"
	AttributeKeyStatus            = "status"
	
	// Pool attributes
	AttributeKeyPoolType          = "pool_type"
	AttributeKeyToken0            = "token0"
	AttributeKeyToken1            = "token1"
	AttributeKeyLiquidity         = "liquidity"
	AttributeKeyShares            = "shares"
	AttributeKeyTokenIn           = "token_in"
	AttributeKeyTokenOut          = "token_out"
	AttributeKeySharesIn          = "shares_in"
	AttributeKeySharesOut         = "shares_out"
	AttributeKeyCreator           = "creator"
	
	// Village pool attributes
	AttributeKeyVillageName       = "village_name"
	AttributeKeyPostalCode        = "postal_code"
	AttributeKeyPanchayatHead     = "panchayat_head"
	AttributeKeyMemberCount       = "member_count"
	AttributeKeyTrustScore        = "trust_score"
	AttributeKeyAchievement       = "achievement"
	
	// Cultural attributes
	AttributeKeyFestival          = "festival"
	AttributeKeyDiscount          = "discount"
	AttributeKeyLanguage          = "language"
	AttributeKeyCulturalQuote     = "cultural_quote"
	
	// Fee distribution attributes
	AttributeKeyValidatorFee      = "validator_fee"
	AttributeKeyPlatformFee       = "platform_fee"
	AttributeKeyLPFee             = "lp_fee"
	AttributeKeyDevelopmentFee    = "development_fee"
	AttributeKeyNGOFee            = "ngo_fee"
	AttributeKeyFounderFee        = "founder_fee"
	
	// Error attributes
	AttributeKeyError             = "error"
	AttributeKeyReason            = "reason"
	
	// Success attributes
	AttributeKeySuccess           = "success"
	AttributeKeyQRCode            = "qr_code"
	AttributeKeyTrackingURL       = "tracking_url"
)