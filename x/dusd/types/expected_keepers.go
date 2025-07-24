package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// BankKeeper defines the expected interface for the bank module
type BankKeeper interface {
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetSupply(ctx sdk.Context, denom string) sdk.Coin
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
}

// AccountKeeper defines the expected interface for the account module
type AccountKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
}

// OracleKeeper defines the expected interface for the oracle module
type OracleKeeper interface {
	GetPrice(ctx sdk.Context, fromDenom, toDenom string) (PriceData, error)
}

// TreasuryKeeper defines the expected interface for the treasury module
type TreasuryKeeper interface {
	AddRevenue(ctx sdk.Context, source string, amount sdk.Coins) error
}

// TaxKeeper defines expected tax keeper for NAMO fee collection
type TaxKeeper interface {
	GetNAMOSwapRouter() NAMOSwapRouter
	DistributePlatformRevenue(ctx sdk.Context, revenueSource string, revenue sdk.Coin) error
}

// RevenueKeeper defines expected revenue keeper
type RevenueKeeper interface {
	RecordRevenue(ctx sdk.Context, source string, amount sdk.Coin) error
}

// NAMOSwapRouter interface for NAMO swapping
type NAMOSwapRouter interface {
	SwapForNAMOFee(ctx sdk.Context, userAddr sdk.AccAddress, feeAmount sdk.Coin, userToken sdk.Coin, inclusive bool) (sdk.Coin, error)
}

// PriceData represents oracle price data
type PriceData struct {
	Price     sdk.Dec
	Timestamp time.Time
}