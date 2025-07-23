package types

const (
	// ModuleName defines the module name
	ModuleName = "revenue"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_revenue"

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName
)

// Store key prefixes
var (
	// Revenue stream records
	RevenueStreamKey = []byte{0x01}
	
	// Performance metrics
	PerformanceMetricsKey = []byte{0x02}
	
	// Revenue distribution records
	RevenueDistributionKey = []byte{0x03}
	
	// Module revenue tracking
	ModuleRevenueKey = []byte{0x04}
	
	// Yield calculation history
	YieldHistoryKey = []byte{0x05}
	
	// Platform statistics
	PlatformStatsKey = []byte{0x06}
)