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

	// Supply distribution percentages - Final v2.0 Economic Model
	PublicSalePercent      = 20  // 20% - 285,725,533 tokens (reduced for scarcity)
	LiquidityPercent       = 18  // 18% - 257,152,979 tokens 
	CommunityPercent       = 15  // 15% - 214,294,149 tokens
	DevelopmentPercent     = 15  // 15% - 214,294,149 tokens
	TeamPercent            = 12  // 12% - 171,435,319 tokens (24-month vesting with 12-month cliff)
	FounderPercent         = 8   // 8% - 114,290,213 tokens (48-month vesting with 12-month cliff)
	DAOTreasuryPercent     = 5   // 5% - 71,431,383 tokens
	CoFounderPercent       = 3.5 // 3.5% - 50,001,968 tokens (24-month vesting with 12-month cliff)
	OperationsPercent      = 2   // 2% - 28,572,553 tokens
	AngelPercent           = 1.5 // 1.5% - 21,428,900 tokens (24-month vesting with 12-month cliff)

	// Vesting parameters - Universal 12-month cliff
	UniversalCliffMonths    = 12 // 12 months cliff for ALL vested allocations
	FounderVestingMonths    = 48 // 48 months for founder vesting
	TeamVestingMonths       = 24 // 24 months for team vesting
	CoFounderVestingMonths  = 24 // 24 months for co-founder vesting
	AngelVestingMonths      = 24 // 24 months for angel vesting
	CommunityDistribMonths  = 60 // 60 months for community rewards
)

// Token allocation amounts (in base units) - Final v2.0 Model
var (
	PublicSaleAllocation    = TotalSupply * 20 / 100        // 285,725,533
	LiquidityAllocation     = TotalSupply * 18 / 100        // 257,152,979
	CommunityAllocation     = TotalSupply * 15 / 100        // 214,294,149
	DevelopmentAllocation   = TotalSupply * 15 / 100        // 214,294,149
	TeamAllocation          = TotalSupply * 12 / 100        // 171,435,319
	FounderAllocation       = TotalSupply * 8 / 100         // 114,290,213
	DAOTreasuryAllocation   = TotalSupply * 5 / 100         // 71,431,383
	CoFounderAllocation     = TotalSupply * 35 / 1000       // 50,001,968 (3.5%)
	OperationsAllocation    = TotalSupply * 2 / 100         // 28,572,553
	AngelAllocation         = TotalSupply * 15 / 1000       // 21,428,900 (1.5%)
)

// Module account names - Updated for v2.0 model
const (
	// PublicSalePoolName is the name of the public sale pool
	PublicSalePoolName = "public_sale_pool"

	// LiquidityPoolName is the name of the liquidity pool
	LiquidityPoolName = "liquidity_pool"

	// CommunityPoolName is the name of the community pool
	CommunityPoolName = "community_pool"

	// DevelopmentPoolName is the name of the development pool
	DevelopmentPoolName = "development_pool"

	// TeamPoolName is the name of the team pool
	TeamPoolName = "team_pool"

	// FounderPoolName is the name of the founder pool
	FounderPoolName = "founder_pool"

	// DAOTreasuryPoolName is the name of the DAO treasury pool
	DAOTreasuryPoolName = "dao_treasury_pool"

	// CoFounderPoolName is the name of the co-founder pool
	CoFounderPoolName = "cofounder_pool"

	// OperationsPoolName is the name of the operations pool
	OperationsPoolName = "operations_pool"

	// AngelPoolName is the name of the angel investor pool
	AngelPoolName = "angel_pool"
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