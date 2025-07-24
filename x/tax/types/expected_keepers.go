package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// BankKeeper defines the expected bank keeper interface
type BankKeeper interface {
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	GetModuleAddress(moduleName string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, moduleName string) authtypes.ModuleAccountI
}

// RevenueKeeper defines the expected revenue keeper interface
type RevenueKeeper interface {
	CollectRevenue(ctx sdk.Context, source string, amount sdk.Coins) error
	RecordRevenue(ctx sdk.Context, source string, amount sdk.Coin) error
	GetTotalRevenue(ctx sdk.Context, source string) sdk.Coin
	GetRevenueByPeriod(ctx sdk.Context, source string, startTime, endTime int64) sdk.Coin
}

// DEXKeeper defines the expected interface for DEX integration
type DEXKeeper interface {
	GetSwapRate(ctx sdk.Context, fromDenom, toDenom string) sdk.Dec
	Swap(ctx sdk.Context, trader sdk.AccAddress, fromCoin sdk.Coin, toDenom string) (sdk.Coin, error)
	GetLiquidity(ctx sdk.Context, denom1, denom2 string) sdk.Dec
}

// OracleKeeper defines the expected interface for price oracle
type OracleKeeper interface {
	GetPrice(ctx sdk.Context, denom string) sdk.Dec
	GetExchangeRate(ctx sdk.Context, fromDenom, toDenom string) sdk.Dec
	IsPriceAvailable(ctx sdk.Context, denom string) bool
}