package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// BankKeeper defines the expected bank keeper interface
type BankKeeper interface {
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	HasBalance(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coin) bool
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
}

// AccountKeeper defines the expected account keeper interface
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	SetAccount(ctx sdk.Context, acc authtypes.AccountI)
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, moduleName string) authtypes.ModuleAccountI
	SetModuleAccount(ctx sdk.Context, macc authtypes.ModuleAccountI)
	HasDhanPataID(ctx sdk.Context, address, dhanPataID string) bool
}

// CulturalKeeper defines the expected cultural keeper interface
type CulturalKeeper interface {
	GetRandomCulturalQuote(ctx sdk.Context) string
	GetActiveFestival(ctx sdk.Context) (string, bool)
	IsFestivalActive(ctx sdk.Context, festivalName string) bool
	GetFestivalBonus(ctx sdk.Context, festivalName string) sdk.Dec
}

// DhanPataKeeper defines the expected DhanPata keeper interface
type DhanPataKeeper interface {
	IsAddressVerified(ctx sdk.Context, address string) bool
	GetDhanPataID(ctx sdk.Context, address string) (string, bool)
	ResolveDhanPataAddress(ctx sdk.Context, dhanPataID string) (sdk.AccAddress, bool)
}

// ParamSubspace defines the expected Subspace interface
type ParamSubspace interface {
	Get(ctx sdk.Context, key []byte, ptr interface{})
	GetParamSet(ctx sdk.Context, ps ParamSet)
	SetParamSet(ctx sdk.Context, ps ParamSet)
	HasKeyTable() bool
	WithKeyTable(table KeyTable) ParamSubspace
}

// ParamSet defines an interface that uses x/params module for parameter management
type ParamSet interface {
	ParamSetPairs() ParamSetPairs
}

// KeyTable defines an interface for parameter key table
type KeyTable interface {
	RegisterType(key []byte, ty interface{}) KeyTable
	RegisterParamSet(ps ParamSet) KeyTable
}

// ParamSetPairs defines the params set pairs type
type ParamSetPairs []ParamSetPair

// ParamSetPair is a key-value pair for module parameters
type ParamSetPair struct {
	Key         []byte
	Value       interface{}
	ValidatorFn func(value interface{}) error
}