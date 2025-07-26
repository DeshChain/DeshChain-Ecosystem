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
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
}

// StakingKeeper defines the expected staking keeper
type StakingKeeper interface {
	BondDenom(ctx sdk.Context) string
}

// DonationKeeper defines the expected donation keeper interface
type DonationKeeper interface {
	GetCharitableTrustWallet(ctx sdk.Context, walletID uint64) (interface{}, bool)
	IsWalletVerified(ctx sdk.Context, walletID uint64) bool
	IsWalletActive(ctx sdk.Context, walletID uint64) bool
}

// GovKeeper defines the expected governance keeper
type GovKeeper interface {
	GetProposal(ctx sdk.Context, proposalID uint64) (proposal interface{}, found bool)
}

// RevenueKeeper defines the expected revenue keeper for integration
type RevenueKeeper interface {
	RecordRevenueStream(ctx sdk.Context, stream interface{})
	CalculateAndDistributeRevenue(ctx sdk.Context) error
	IsRevenueEnabled(ctx sdk.Context) bool
}