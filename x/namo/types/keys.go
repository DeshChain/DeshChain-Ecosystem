package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "namo"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_namo"
)

var (
	// ParamsKey is the key for module parameters
	ParamsKey = collections.NewPrefix(0)

	// TokenSupplyKey is the key for token supply information
	TokenSupplyKey = collections.NewPrefix(1)

	// VestingScheduleKey is the key for vesting schedules
	VestingScheduleKey = collections.NewPrefix(2)

	// VestingScheduleCountKey is the key for vesting schedule counter
	VestingScheduleCountKey = collections.NewPrefix(3)

	// DistributionEventKey is the key for distribution events
	DistributionEventKey = collections.NewPrefix(4)

	// DistributionEventCountKey is the key for distribution event counter
	DistributionEventCountKey = collections.NewPrefix(5)
)

// NAMO Token Configuration Constants
const (
	// DefaultTokenDenom is the default denomination for NAMO tokens
	DefaultTokenDenom = "namo"

	// TotalSupply is the total supply of NAMO tokens: 1,428,627,663
	TotalSupply = 1428627663

	// Supply distribution percentages - Updated for sustainable model
	PublicSalePercent      = 25 // 25% - 357,156,916 tokens
	LiquidityPercent       = 20 // 20% - 285,725,533 tokens (increased for stability)
	FounderPercent         = 10 // 10% - 142,862,766 tokens (48-month vesting with 12-month cliff)
	TeamPercent            = 10 // 10% - 142,862,766 tokens (24-month vesting)
	DevelopmentPercent     = 15 // 15% - 214,294,149 tokens
	CommunityPercent       = 15 // 15% - 214,294,149 tokens (increased for engagement)
	DAOTreasuryPercent     = 5  // 5% - 71,431,383 tokens
	InitialBurnPercent     = 0  // 0% - 0 tokens (reallocated to liquidity and community)

	// Vesting parameters
	FounderVestingMonths    = 48 // 48 months for founder vesting
	FounderCliffMonths      = 12 // 12 months cliff for founder
	TeamVestingMonths       = 24 // 24 months for team vesting
	CommunityDistribMonths  = 60 // 60 months for community rewards
	DefaultCliffMonths      = 6  // 6 months cliff period
)

// Token allocation amounts (in base units)
var (
	PublicSaleAllocation    = TotalSupply * PublicSalePercent / 100      // 357,156,916
	LiquidityAllocation     = TotalSupply * LiquidityPercent / 100       // 285,725,533
	FounderAllocation       = TotalSupply * FounderPercent / 100          // 142,862,766
	TeamAllocation          = TotalSupply * TeamPercent / 100            // 142,862,766
	DevelopmentAllocation   = TotalSupply * DevelopmentPercent / 100     // 214,294,149
	CommunityAllocation     = TotalSupply * CommunityPercent / 100       // 214,294,149
	DAOTreasuryAllocation   = TotalSupply * DAOTreasuryPercent / 100     // 71,431,383
	InitialBurnAllocation   = TotalSupply * InitialBurnPercent / 100     // 0
)

// Module account names
const (
	// PublicSalePoolName is the name of the public sale pool
	PublicSalePoolName = "public_sale_pool"

	// LiquidityPoolName is the name of the liquidity pool
	LiquidityPoolName = "liquidity_pool"

	// FounderPoolName is the name of the founder pool
	FounderPoolName = "founder_pool"

	// TeamPoolName is the name of the team pool
	TeamPoolName = "team_pool"

	// DevelopmentPoolName is the name of the development pool
	DevelopmentPoolName = "development_pool"

	// CommunityPoolName is the name of the community pool
	CommunityPoolName = "community_pool"

	// DAOTreasuryPoolName is the name of the DAO treasury pool
	DAOTreasuryPoolName = "dao_treasury_pool"

	// BurnPoolName is the name of the burn pool
	BurnPoolName = "burn_pool"
)

// Event types for distribution events
const (
	EventTypeInitialDistribution = "initial_distribution"
	EventTypeVestingClaim        = "vesting_claim"
	EventTypeCommunityReward     = "community_reward"
	EventTypeTokenBurn           = "token_burn"
	EventTypeDevelopmentFund     = "development_fund"
	EventTypeDAOTreasury         = "dao_treasury"
	EventTypeLiquidityProvision  = "liquidity_provision"
)