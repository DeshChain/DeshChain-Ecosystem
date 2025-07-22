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

import (
	"time"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "moneyorder"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_moneyorder"
	
	// DefaultDenom is the default denomination for the module
	DefaultDenom = "namo"
	
	// EscrowModuleName for escrow account
	EscrowModuleName = "moneyorder_escrow"
	
	// SevaMitraSecurityPool for seva mitra deposits
	SevaMitraSecurityPool = "sevamitra_security"
)

// Store key prefixes
var (
	// Pools
	KeyPrefixAMMPool              = []byte{0x01}
	KeyPrefixFixedRatePool        = []byte{0x02}
	KeyPrefixVillagePool          = []byte{0x03}
	KeyPrefixConcentratedPool     = []byte{0x04}
	
	// Orders and Receipts
	KeyPrefixMoneyOrder           = []byte{0x10}
	KeyPrefixReceipt              = []byte{0x11}
	KeyPrefixOrderByUser          = []byte{0x12}
	KeyPrefixOrderByPostalCode    = []byte{0x13}
	
	// Trading
	KeyPrefixTradingPair          = []byte{0x20}
	KeyPrefixLiquidity            = []byte{0x21}
	KeyPrefixPosition             = []byte{0x22}
	KeyPrefixTick                 = []byte{0x23}
	KeyPrefixVolume               = []byte{0x24}
	
	// Cultural Features
	KeyPrefixFestivalBonus        = []byte{0x30}
	KeyPrefixCulturalPair         = []byte{0x31}
	KeyPrefixVillageIncentive     = []byte{0x32}
	
	// Params and State
	KeyPrefixParams               = []byte{0x40}
	KeyPrefixSequence             = []byte{0x41}
	KeyPrefixFeeDistribution      = []byte{0x42}
	
	// Additional keys
	KeyPrefixUserReceipt          = []byte{0x10}
	KeyPrefixUPIAddress           = []byte{0x11}
	KeyPrefixUserDailyLimit       = []byte{0x12}
	KeyPrefixScheduledOrder       = []byte{0x13}
	KeyPrefixCulturalQuote        = []byte{0x14}
	KeyPrefixFestivalPeriod       = []byte{0x15}
	KeyNextPoolId                 = []byte{0x16}
	KeyPrefixVillagePoolMember    = []byte{0x17}
	KeyPrefixPensionLiquidity     = []byte{0x18}
	KeyPrefixUnifiedPool          = []byte{0x19}
	KeyNextUnifiedPoolId          = []byte{0x1A}
	KeyPrefixSurakshaContribution  = []byte{0x1B}
	KeyPrefixAgriLoan             = []byte{0x1C}
	KeyNextLoanId                 = []byte{0x1D}
	KeyDynamicPensionRate         = []byte{0x1E}
	
	// P2P and Escrow prefixes
	P2POrderPrefix                = []byte{0x50}
	P2PTradePrefix                = []byte{0x51}
	AgentPrefix                   = []byte{0x52}
	EscrowPrefix                  = []byte{0x53}
	DisputePrefix                 = []byte{0x54}
	UserStatsPrefix               = []byte{0x55}
	TrustScorePrefix              = []byte{0x56}
	EscrowOrderIndexPrefix        = []byte{0x57}
	EscrowTradeIndexPrefix        = []byte{0x58}
	DisputeEscrowIndexPrefix      = []byte{0x59}
	EscrowExpiryQueuePrefix       = []byte{0x5A}
	RefundQueuePrefix             = []byte{0x5B}
	P2POrderTypeIndexPrefix       = []byte{0x5C}
	AgentDistrictIndexPrefix      = []byte{0x5D}
	UserStatsPrefix               = []byte{0x5E}
)

// Module account names
const (
	// Fee collection accounts
	MoneyOrderFeeCollector      = "moneyorder_fees"
	ValidatorPoolName           = "moneyorder_validator_pool"
	PlatformPoolName            = "moneyorder_platform_pool"
	LiquidityProviderPoolName   = "moneyorder_lp_pool"
	DevelopmentPoolName         = "moneyorder_dev_pool"
	NGODonationPoolName         = "moneyorder_ngo_pool"
	FounderRoyaltyPoolName      = "moneyorder_founder_pool"
	
	// Special purpose accounts
	VillagePoolReserve          = "village_pool_reserve"
	FestivalBonusPool           = "festival_bonus_pool"
	KYCEscrowAccount            = "kyc_escrow_account"
	FeeCollectorName            = "fee_collector"
)

// Order status constants
const (
	OrderStatusPending    = "pending"
	OrderStatusProcessing = "processing"
	OrderStatusCompleted  = "completed"
	OrderStatusCancelled  = "cancelled"
	OrderStatusRefunded   = "refunded"
	OrderStatusExpired    = "expired"
)

// Pool type constants
const (
	PoolTypeAMM           = "amm"
	PoolTypeFixedRate     = "fixed_rate"
	PoolTypeConcentrated  = "concentrated"
	PoolTypeVillage       = "village"
)

// Cultural constants
const (
	// Festival periods (can be updated via governance)
	FestivalDiwali        = "diwali"
	FestivalHoli          = "holi"
	FestivalDussehra      = "dussehra"
	FestivalIndependence  = "independence_day"
	FestivalRepublic      = "republic_day"
	
	// Language codes
	LanguageHindi         = "hi"
	LanguageEnglish       = "en"
	LanguageBengali       = "bn"
	LanguageTelugu        = "te"
	LanguageMarathi       = "mr"
	LanguageTamil         = "ta"
	LanguageGujarati      = "gu"
	LanguageUrdu          = "ur"
)

// Money Order limits
var (
	// Traditional money order limits (in base NAMO units)
	MinMoneyOrderAmount     = sdk.NewInt(10_000_000)      // 10 NAMO
	MaxMoneyOrderAmount     = sdk.NewInt(50_000_000_000)  // 50,000 NAMO
	MaxDailyUserLimit       = sdk.NewInt(100_000_000_000) // 100,000 NAMO
	
	// KYC thresholds
	KYCRequiredThreshold    = sdk.NewInt(10_000_000_000)  // 10,000 NAMO
	
	// Village pool minimums
	MinVillagePoolLiquidity = sdk.NewInt(1_000_000_000)   // 1,000 NAMO
)

// Fee constants (matching x/dex but with Money Order context)
const (
	// Base trading fee (0.3%)
	DefaultTradingFeeRate = "0.003"
	
	// Money Order specific fees
	MoneyOrderBaseFee     = "0.001"  // 0.1% for simple transfers
	ExpressDeliveryFee    = "0.002"  // 0.2% for instant
	BulkOrderDiscount     = "0.20"   // 20% discount for bulk
	
	// Cultural discounts
	VillagePoolDiscount   = "0.30"   // 30% discount
	FestivalDiscount      = "0.25"   // 25% discount
	SeniorCitizenDiscount = "0.15"   // 15% discount
)

// GetAMMPoolKey returns the store key for an AMM pool
func GetAMMPoolKey(poolId uint64) []byte {
	return append(KeyPrefixAMMPool, sdk.Uint64ToBigEndian(poolId)...)
}

// GetFixedRatePoolKey returns the store key for a fixed rate pool
func GetFixedRatePoolKey(poolId uint64) []byte {
	return append(KeyPrefixFixedRatePool, sdk.Uint64ToBigEndian(poolId)...)
}

// GetVillagePoolKey returns the store key for a village pool
func GetVillagePoolKey(postalCode string) []byte {
	return append(KeyPrefixVillagePool, []byte(postalCode)...)
}

// GetMoneyOrderKey returns the store key for a money order
func GetMoneyOrderKey(orderId string) []byte {
	return append(KeyPrefixMoneyOrder, []byte(orderId)...)
}

// GetReceiptKey returns the store key for a receipt
func GetReceiptKey(receiptId string) []byte {
	return append(KeyPrefixReceipt, []byte(receiptId)...)
}

// GetUserOrdersKey returns the store key for user's orders
func GetUserOrdersKey(userAddr sdk.AccAddress) []byte {
	return append(KeyPrefixOrderByUser, userAddr.Bytes()...)
}

// GetPostalCodeOrdersKey returns the store key for postal code orders
func GetPostalCodeOrdersKey(postalCode string) []byte {
	return append(KeyPrefixOrderByPostalCode, []byte(postalCode)...)
}

// GetConcentratedPoolKey returns the store key for a concentrated liquidity pool
func GetConcentratedPoolKey(poolId uint64) []byte {
	return append(KeyPrefixConcentratedPool, sdk.Uint64ToBigEndian(poolId)...)
}

// GetPositionKey returns the store key for a liquidity position
func GetPositionKey(positionId uint64) []byte {
	return append(KeyPrefixPosition, sdk.Uint64ToBigEndian(positionId)...)
}

// GetTickKey returns the store key for a tick
func GetTickKey(poolId uint64, tickIndex int64) []byte {
	key := append(KeyPrefixTick, sdk.Uint64ToBigEndian(poolId)...)
	return append(key, sdk.Int64ToBigEndian(tickIndex)...)
}

// KeyConcentratedPool returns key for concentrated liquidity pool
func KeyConcentratedPool(poolId uint64) []byte {
	return GetConcentratedPoolKey(poolId)
}

// KeyPosition returns key for liquidity position
func KeyPosition(positionId uint64) []byte {
	return GetPositionKey(positionId)
}

// KeyTick returns key for price tick
func KeyTick(poolId uint64, tickIndex int64) []byte {
	return GetTickKey(poolId, tickIndex)
}

// ParsePoolIdFromKey extracts pool ID from a store key
func ParsePoolIdFromKey(key []byte, prefix []byte) uint64 {
	if len(key) < len(prefix)+8 {
		return 0
	}
	return sdk.BigEndianToUint64(key[len(prefix):])
}

// GetP2POrderKey returns the store key for a P2P order
func GetP2POrderKey(orderId string) []byte {
	return append(P2POrderPrefix, []byte(orderId)...)
}

// GetP2PTradeKey returns the store key for a P2P trade
func GetP2PTradeKey(tradeId string) []byte {
	return append(P2PTradePrefix, []byte(tradeId)...)
}

// GetAgentKey returns the store key for an agent
func GetAgentKey(agentId string) []byte {
	return append(AgentPrefix, []byte(agentId)...)
}

// GetEscrowKey returns the store key for an escrow
func GetEscrowKey(escrowId string) []byte {
	return append(EscrowPrefix, []byte(escrowId)...)
}

// GetDisputeKey returns the store key for a dispute
func GetDisputeKey(disputeId string) []byte {
	return append(DisputePrefix, []byte(disputeId)...)
}

// GetUserStatsKey returns the store key for user stats
func GetUserStatsKey(address string) []byte {
	return append(UserStatsPrefix, []byte(address)...)
}

// GetEscrowOrderIndexKey returns the store key for escrow order index
func GetEscrowOrderIndexKey(orderId string) []byte {
	return append(EscrowOrderIndexPrefix, []byte(orderId)...)
}

// GetEscrowTradeIndexKey returns the store key for escrow trade index
func GetEscrowTradeIndexKey(tradeId string) []byte {
	return append(EscrowTradeIndexPrefix, []byte(tradeId)...)
}

// GetDisputeEscrowIndexKey returns the store key for dispute escrow index
func GetDisputeEscrowIndexKey(escrowId string) []byte {
	return append(DisputeEscrowIndexPrefix, []byte(escrowId)...)
}

// GetEscrowExpiryQueueKey returns the store key for escrow expiry queue
func GetEscrowExpiryQueueKey(expiresAt time.Time, escrowId string) []byte {
	timeBytes := sdk.FormatTimeBytes(expiresAt)
	key := append(EscrowExpiryQueuePrefix, timeBytes...)
	return append(key, []byte(escrowId)...)
}

// ParseTimeFromBytes parses time from bytes
func ParseTimeFromBytes(bz []byte) time.Time {
	return sdk.ParseTimeBytes(bz)
}

// GetSevaMitraKey returns the store key for a seva mitra
func GetSevaMitraKey(mitraId string) []byte {
	return append(AgentPrefix, []byte(mitraId)...)
}

// GetSevaMitraAddressIndexKey returns the store key for seva mitra address index
func GetSevaMitraAddressIndexKey(address string) []byte {
	return append([]byte("mitra-addr-"), []byte(address)...)
}

// GetSevaMitraDistrictIndexKey returns the store key for seva mitra district index
func GetSevaMitraDistrictIndexKey(district, mitraId string) []byte {
	key := append(AgentDistrictIndexPrefix, []byte(district)...)
	return append(key, []byte(mitraId)...)
}

// GetSevaMitraDistrictIndexPrefix returns the prefix for seva mitra district index
func GetSevaMitraDistrictIndexPrefix(district string) []byte {
	return append(AgentDistrictIndexPrefix, []byte(district)...)
}

// GetSevaMitraDailyVolumeKey returns the store key for seva mitra daily volume
func GetSevaMitraDailyVolumeKey(mitraId, date string) []byte {
	key := append([]byte("mitra-daily-vol-"), []byte(mitraId)...)
	return append(key, []byte(date)...)
}

// GetSevaMitraRatingKey returns the store key for seva mitra rating
func GetSevaMitraRatingKey(mitraId string, rater sdk.AccAddress) []byte {
	key := append([]byte("mitra-rating-"), []byte(mitraId)...)
	return append(key, rater.Bytes()...)
}

// GetSevaMitraRatingsPrefix returns the prefix for seva mitra ratings
func GetSevaMitraRatingsPrefix(mitraId string) []byte {
	return append([]byte("mitra-rating-"), []byte(mitraId)...)
}

// GetKYCQueueKey returns the store key for KYC queue
func GetKYCQueueKey(scheduledAt time.Time, mitraId string) []byte {
	timeBytes := sdk.FormatTimeBytes(scheduledAt)
	key := append([]byte("kyc-queue-"), timeBytes...)
	return append(key, []byte(mitraId)...)
}

// Note: Using AgentPrefix and AgentDistrictIndexPrefix as SevaMitraPrefix 
// and SevaMitraDistrictIndexPrefix for backward compatibility
var SevaMitraPrefix = AgentPrefix
var SevaMitraDistrictIndexPrefix = AgentDistrictIndexPrefix

// Event types
const (
	EventTypeMoneyOrder          = "money_order"
	EventTypeP2POrderCreated     = "p2p_order_created"
	EventTypeP2PTradeMatched     = "p2p_trade_matched"
	EventTypeP2POrderRefunded    = "p2p_order_refunded"
	EventTypeEscrowRefunded      = "escrow_refunded"
	EventTypeSevaMitraRegistered = "seva_mitra_registered"
	EventTypeSevaMitraKYCCompleted = "seva_mitra_kyc_completed"
	EventTypeSevaMitraStatusUpdated = "seva_mitra_status_updated"
	
	// Attribute keys
	AttributeKeyOrderID    = "order_id"
	AttributeKeyFees       = "fees"
	AttributeKeySender     = "sender"
	AttributeKeyAmount     = "amount"
	AttributeKeyPostalCode = "postal_code"
)