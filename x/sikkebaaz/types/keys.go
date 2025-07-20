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

const (
	// ModuleName defines the module name
	ModuleName = "sikkebaaz"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_sikkebaaz"
	
	// DefaultDenom is the default denomination for Sikkebaaz
	DefaultDenom = "namo"
)

// Store key prefixes for Sikkebaaz launchpad
var (
	// Token Launch Management
	KeyPrefixTokenLaunch        = []byte{0x01} // launch_id -> TokenLaunch
	KeyPrefixCreatorLaunches    = []byte{0x02} // creator -> launch_ids
	KeyPrefixPincodeLaunches    = []byte{0x03} // pincode -> launch_ids
	KeyPrefixActiveLaunches     = []byte{0x04} // active launches
	
	// Trading and Liquidity
	KeyPrefixTokenTradingPair   = []byte{0x10} // trading pairs
	KeyPrefixLiquidityLock      = []byte{0x11} // locked liquidity
	KeyPrefixCreatorRewards     = []byte{0x12} // creator reward tracking
	KeyPrefixTradingMetrics     = []byte{0x13} // volume, fees, etc.
	
	// Anti-Pump Protection
	KeyPrefixWalletLimits       = []byte{0x20} // wallet -> limits
	KeyPrefixTradingRestrictions = []byte{0x21} // token -> restrictions
	KeyPrefixBotDetection       = []byte{0x22} // anti-bot measures
	KeyPrefixLaunchTimestamps   = []byte{0x23} // launch timing
	
	// Community Features
	KeyPrefixCommunityVeto      = []byte{0x30} // community voting
	KeyPrefixLocalNGOAllocation = []byte{0x31} // charity allocations
	KeyPrefixFestivalBonuses    = []byte{0x32} // festival rewards
	KeyPrefixCulturalQuotes     = []byte{0x33} // cultural integration
	
	// Governance and Security
	KeyPrefixMultisigApprovals  = []byte{0x40} // large launch approvals
	KeyPrefixSecurityAudits     = []byte{0x41} // audit results
	KeyPrefixEmergencyControls  = []byte{0x42} // emergency stops
	KeyPrefixComplianceFlags    = []byte{0x43} // regulatory compliance
)

// Launch status constants
const (
	LaunchStatusPending     = "pending"     // Created but not active
	LaunchStatusActive      = "active"      // Currently raising funds
	LaunchStatusSuccessful  = "successful" // Successfully completed
	LaunchStatusFailed      = "failed"     // Failed to reach target
	LaunchStatusCancelled   = "cancelled"  // Cancelled by creator
	LaunchStatusVetoed      = "vetoed"     // Community vetoed
)

// Launch type constants
const (
	LaunchTypeFair      = "fair"      // Fair launch - no pre-sale
	LaunchTypeStealth   = "stealth"   // Stealth launch - surprise
	LaunchTypeWhitelist = "whitelist" // Whitelist only
	LaunchTypeAuction   = "auction"   // Dutch auction style
	LaunchTypePrivate   = "private"   // Private round first
)

// Anti-pump configuration constants
const (
	// Wallet limits (percentage of total supply)
	MaxWalletPercent24h = 500  // 5% for first 24 hours
	MaxWalletPercentAfter = 1000 // 10% after 24 hours
	
	// Trading restrictions
	MinTradingDelay = 300      // 5 minutes before trading starts
	MaxTradingDelay = 86400    // 24 hours maximum delay
	
	// Liquidity lock requirements
	MinLiquidityLockDays = 365 // 1 year minimum
	MaxLiquidityLockDays = 1825 // 5 years maximum
	
	// Bot protection
	MinBlocksBetweenTx = 3     // Minimum blocks between transactions
	MaxGasPrice = 50000000000  // Maximum gas price (50 gwei equivalent)
)

// Fee structure constants
const (
	// Launch fees
	BaseLaunchFee = "1000000000000000" // 1000 NAMO (with 12 decimals)
	VariableFeeRate = "0.05"           // 5% of raised amount
	
	// Creator rewards
	CreatorTradingReward = "0.02"      // 2% of trading volume
	
	// Charity allocation
	LocalNGOAllocation = "0.01"        // 1% to local NGO
	
	// Community veto
	CommunityVetoThreshold = "0.70"    // 70% vote required
)

// Cultural integration constants
const (
	// Festival bonuses
	FestivalBonusRate = "0.10"         // 10% bonus during festivals
	
	// Cultural features
	MaxCulturalQuoteLength = 280       // Tweet-like length
	MinPatriotismScore = 50            // Minimum score for launches
	
	// Pincode features
	MaxPincodeLength = 6               // Indian PIN code length
	MinPincodePopulation = 1000        // Minimum population for memecoin
)

// Module account names for Sikkebaaz
const (
	// Fee collection
	SikkebaazFeeCollector = "sikkebaaz_fees"
	
	// Security and escrow
	LaunchEscrowAccount = "launch_escrow"
	LiquidityLockAccount = "liquidity_lock"
	
	// Rewards and incentives
	CreatorRewardsPool = "creator_rewards"
	CommunityIncentivePool = "community_incentives"
	
	// Charity and social impact
	LocalNGOPool = "local_ngo_pool"
	FestivalBonusPool = "festival_bonus"
	
	// Security and compliance
	SecurityAuditFund = "security_audit"
	EmergencyFund = "emergency_fund"
)

// Event types for Sikkebaaz
const (
	EventTypeTokenLaunched = "token_launched"
	EventTypeLaunchCompleted = "launch_completed"
	EventTypeLaunchFailed = "launch_failed"
	EventTypeCommunityVeto = "community_veto"
	EventTypeCreatorReward = "creator_reward"
	EventTypeLiquidityLocked = "liquidity_locked"
	EventTypeAntiPumpTriggered = "anti_pump_triggered"
	EventTypeFestivalBonus = "festival_bonus"
	
	// Attribute keys
	AttributeKeyLaunchID = "launch_id"
	AttributeKeyTokenSymbol = "token_symbol"
	AttributeKeyCreator = "creator"
	AttributeKeyPincode = "pincode"
	AttributeKeyRaisedAmount = "raised_amount"
	AttributeKeyLaunchType = "launch_type"
	AttributeKeyTargetAmount = "target_amount"
)

// GetTokenLaunchKey returns the store key for a token launch
func GetTokenLaunchKey(launchId string) []byte {
	return append(KeyPrefixTokenLaunch, []byte(launchId)...)
}

// GetCreatorLaunchesKey returns the store key for creator's launches
func GetCreatorLaunchesKey(creator string) []byte {
	return append(KeyPrefixCreatorLaunches, []byte(creator)...)
}

// GetPincodeLaunchesKey returns the store key for pincode launches
func GetPincodeLaunchesKey(pincode string) []byte {
	return append(KeyPrefixPincodeLaunches, []byte(pincode)...)
}

// GetLiquidityLockKey returns the store key for liquidity lock
func GetLiquidityLockKey(tokenAddress string) []byte {
	return append(KeyPrefixLiquidityLock, []byte(tokenAddress)...)
}

// GetCreatorRewardsKey returns the store key for creator rewards
func GetCreatorRewardsKey(creator, tokenAddress string) []byte {
	key := append(KeyPrefixCreatorRewards, []byte(creator)...)
	return append(key, []byte(tokenAddress)...)
}

// GetCommunityVetoKey returns the store key for community veto
func GetCommunityVetoKey(launchId string) []byte {
	return append(KeyPrefixCommunityVeto, []byte(launchId)...)
}

// GetWalletLimitsKey returns the store key for wallet limits
func GetWalletLimitsKey(tokenAddress, wallet string) []byte {
	key := append(KeyPrefixWalletLimits, []byte(tokenAddress)...)
	return append(key, []byte(wallet)...)
}