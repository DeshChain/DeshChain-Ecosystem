package types

const (
	// ModuleName defines the module name
	ModuleName = "liquiditymanager"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_liquiditymanager"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

// KVStore keys
var (
	LiquidityStatusKey      = []byte{0x01}
	DailyLendingUsedKey     = []byte{0x02}
	VillageMemberPrefix     = []byte{0x03}
	UrbanMemberPrefix       = []byte{0x04}
	StakedNAMOPrefix        = []byte{0x05}
	LockedCollateralPrefix  = []byte{0x06}
	CollateralLoanPrefix    = []byte{0x07}
	LendingStatsPrefix      = []byte{0x08}
)

// GetVillageMemberKey returns the key for village member storage
func GetVillageMemberKey(user string) []byte {
	return append(VillageMemberPrefix, []byte(user)...)
}

// GetUrbanMemberKey returns the key for urban member storage
func GetUrbanMemberKey(user string) []byte {
	return append(UrbanMemberPrefix, []byte(user)...)
}

// GetStakedNAMOKey returns the key for staked NAMO storage
func GetStakedNAMOKey(user string) []byte {
	return append(StakedNAMOPrefix, []byte(user)...)
}

// GetLockedCollateralKey returns the key for locked collateral storage
func GetLockedCollateralKey(user string) []byte {
	return append(LockedCollateralPrefix, []byte(user)...)
}

// GetCollateralLoanKey returns the key for collateral loan storage
func GetCollateralLoanKey(user string) []byte {
	return append(CollateralLoanPrefix, []byte(user)...)
}