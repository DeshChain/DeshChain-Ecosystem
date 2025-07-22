package types

const (
	// ModuleName defines the module name
	ModuleName = "dinr"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_dinr"

	// DINRDenom is the denomination for the DINR stablecoin
	DINRDenom = "dinr"
)

// Key prefixes for store
var (
	ParamsKey                  = []byte{0x01}
	UserPositionPrefix         = []byte{0x02}
	CollateralAssetPrefix      = []byte{0x03}
	StabilityDataKey           = []byte{0x04}
	InsuranceFundKey           = []byte{0x05}
	YieldStrategyPrefix        = []byte{0x06}
	TotalDINRMintedKey         = []byte{0x07}
	TotalFeesCollectedKey      = []byte{0x08}
	LastYieldProcessingTimeKey = []byte{0x09}
)

// GetUserPositionKey returns the store key for a user position
func GetUserPositionKey(address string) []byte {
	return append(UserPositionPrefix, []byte(address)...)
}

// GetCollateralAssetKey returns the store key for a collateral asset
func GetCollateralAssetKey(denom string) []byte {
	return append(CollateralAssetPrefix, []byte(denom)...)
}

// GetYieldStrategyKey returns the store key for a yield strategy
func GetYieldStrategyKey(id string) []byte {
	return append(YieldStrategyPrefix, []byte(id)...)
}