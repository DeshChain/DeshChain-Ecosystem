package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// AccountKeeper defines the expected account keeper
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, name string) authtypes.ModuleAccountI
}

// BankKeeper defines the expected bank keeper
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
}

// TaxKeeper defines the expected tax keeper interface
type TaxKeeper interface {
	GetTransactionVolume(ctx sdk.Context) sdk.Dec
}

// DINRKeeper defines the expected DINR keeper interface
type DINRKeeper interface {
	GetTotalSupply(ctx sdk.Context) sdk.Coin
	GetRevenue(ctx sdk.Context) sdk.Dec
}

// DUSDKeeper defines the expected DUSD keeper interface
type DUSDKeeper interface {
	GetRevenue(ctx sdk.Context) sdk.Dec
	GetVolume(ctx sdk.Context) sdk.Dec
}

// LendingKeeper defines the expected lending keeper interface
type LendingKeeper interface {
	GetTotalLendingVolume(ctx sdk.Context) sdk.Dec
	GetDefaultRate(ctx sdk.Context) sdk.Dec
}

// TradeKeeper defines the expected trade finance keeper interface
type TradeKeeper interface {
	GetTradeVolume(ctx sdk.Context) sdk.Dec
}

// RemittanceKeeper defines the expected remittance keeper interface
type RemittanceKeeper interface {
	GetRemittanceVolume(ctx sdk.Context) sdk.Dec
}

// GovernanceKeeper defines the expected governance keeper interface
type GovernanceKeeper interface {
	GetCurrentCharityPercentage(ctx sdk.Context) sdk.Dec
}