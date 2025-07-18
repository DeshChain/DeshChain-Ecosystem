package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "royalty"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_royalty"
)

var (
	// ParamsKey is the key for module parameters
	ParamsKey = collections.NewPrefix(0)

	// RoyaltyConfigKey is the key for royalty configuration
	RoyaltyConfigKey = collections.NewPrefix(1)

	// RoyaltyClaimKey is the key for royalty claims
	RoyaltyClaimKey = collections.NewPrefix(2)

	// InheritanceRecordKey is the key for inheritance records
	InheritanceRecordKey = collections.NewPrefix(3)

	// BeneficiaryHistoryKey is the key for beneficiary history
	BeneficiaryHistoryKey = collections.NewPrefix(4)

	// RoyaltyAccumulatorKey is the key for accumulated royalties
	RoyaltyAccumulatorKey = collections.NewPrefix(5)

	// TransactionRoyaltyKey is the key for transaction-based royalties
	TransactionRoyaltyKey = collections.NewPrefix(6)

	// PlatformRoyaltyKey is the key for platform revenue royalties
	PlatformRoyaltyKey = collections.NewPrefix(7)
)

// Module account names
const (
	// RoyaltyCollectorAccount collects all royalties before distribution
	RoyaltyCollectorAccount = "royalty_collector"

	// TransactionRoyaltyAccount holds transaction-based royalties
	TransactionRoyaltyAccount = "transaction_royalty"

	// PlatformRoyaltyAccount holds platform revenue royalties
	PlatformRoyaltyAccount = "platform_royalty"

	// InheritanceEscrowAccount holds funds during inheritance transition
	InheritanceEscrowAccount = "inheritance_escrow"
)