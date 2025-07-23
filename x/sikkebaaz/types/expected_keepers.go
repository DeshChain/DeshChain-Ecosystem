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

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

// AccountKeeper defines the expected account keeper
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	SetAccount(ctx sdk.Context, acc types.AccountI)
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, moduleName string) types.ModuleAccountI
}

// BankKeeper defines the expected bank keeper
type BankKeeper interface {
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
}

// CulturalKeeper defines the expected cultural keeper interface
type CulturalKeeper interface {
	GetRandomCulturalQuote(ctx sdk.Context) string
	GetPatriotismScore(ctx sdk.Context, text string) int
	GetFestivalBonus(ctx sdk.Context, timestamp int64) int
}

// TreasuryKeeper defines the expected treasury keeper interface (NAMO keeper)
type TreasuryKeeper interface {
	TransferToTreasury(ctx sdk.Context, from sdk.AccAddress, amount sdk.Coins) error
	GetTreasuryBalance(ctx sdk.Context) sdk.Coins
}

// RevenueKeeper defines the expected revenue keeper interface
type RevenueKeeper interface {
	CollectLaunchFee(ctx sdk.Context, moduleName string, creator sdk.AccAddress, fee sdk.Coins, tokenSymbol string) error
	CollectTradingFee(ctx sdk.Context, moduleName string, trader sdk.AccAddress, fee sdk.Coins, pair string) error
	RecordRevenue(ctx sdk.Context, source, category string, amount sdk.Coins, description string) error
}

// MoneyOrderKeeper defines the expected money order keeper interface for DEX integration
type MoneyOrderKeeper interface {
	CreateAMMPool(ctx sdk.Context, creator sdk.AccAddress, poolAssets interface{}, swapFee sdk.Dec) (uint64, error)
}