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
	KeyPrefixVolume               = []byte{0x23}
	
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
	KeyPrefixPensionContribution  = []byte{0x1B}
	KeyPrefixAgriLoan             = []byte{0x1C}
	KeyNextLoanId                 = []byte{0x1D}
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

// ParsePoolIdFromKey extracts pool ID from a store key
func ParsePoolIdFromKey(key []byte, prefix []byte) uint64 {
	if len(key) < len(prefix)+8 {
		return 0
	}
	return sdk.BigEndianToUint64(key[len(prefix):])
}