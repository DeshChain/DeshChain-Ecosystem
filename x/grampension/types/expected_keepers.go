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
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// AccountKeeper defines the expected account keeper
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	GetModuleAccount(ctx sdk.Context, moduleName string) authtypes.ModuleAccountI
	SetModuleAccount(ctx sdk.Context, macc authtypes.ModuleAccountI)
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	SetAccount(ctx sdk.Context, acc authtypes.AccountI)
}

// BankKeeper defines the expected bank keeper
type BankKeeper interface {
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	BlockedAddr(addr sdk.AccAddress) bool
}

// StakingKeeper defines the expected staking keeper
type StakingKeeper interface {
	BondDenom(ctx sdk.Context) string
}

// MoneyOrderKeeper defines the expected money order keeper for liquidity integration
type MoneyOrderKeeper interface {
	// Hooks for pension events
	GetMoneyOrderHooks() GramPensionHooks
}

// GramPensionHooks defines the hooks that money order module will implement
type GramPensionHooks interface {
	// Called after a pension contribution is made
	AfterPensionContribution(
		ctx sdk.Context,
		pensionAccountId string,
		contributor sdk.AccAddress,
		contribution sdk.Coin,
		villagePostalCode string,
	) error

	// Called when pension reaches maturity
	AfterPensionMaturity(
		ctx sdk.Context,
		pensionAccountId string,
		beneficiary sdk.AccAddress,
		maturityAmount sdk.Coin,
	) error

	// Called monthly for revenue distribution
	MonthlyRevenueDistribution(ctx sdk.Context) error
}

// CulturalKeeper defines the expected cultural keeper
type CulturalKeeper interface {
	GetRandomQuote(ctx sdk.Context, category string) (string, error)
	GetPatriotismScore(ctx sdk.Context, address sdk.AccAddress) (uint64, error)
	UpdatePatriotismScore(ctx sdk.Context, address sdk.AccAddress, points uint64) error
}

// TaxKeeper defines the expected tax keeper
type TaxKeeper interface {
	GetCurrentTaxRate(ctx sdk.Context) sdk.Dec
	CalculateTax(ctx sdk.Context, amount sdk.Coin) (sdk.Coin, error)
}

// DonationKeeper defines the expected donation keeper
type DonationKeeper interface {
	ProcessDonation(ctx sdk.Context, from sdk.AccAddress, amount sdk.Coin, ngoId string) error
	GetNGODetails(ctx sdk.Context, ngoId string) (NGODetails, error)
}

// NGODetails represents basic NGO information
type NGODetails struct {
	ID          string
	Name        string
	Description string
	Address     sdk.AccAddress
	Status      string
}

// KYCKeeper defines the expected KYC keeper
type KYCKeeper interface {
	GetKYCStatus(ctx sdk.Context, address sdk.AccAddress) (string, error)
	VerifyKYC(ctx sdk.Context, address sdk.AccAddress, level string) error
	IsKYCVerified(ctx sdk.Context, address sdk.AccAddress) bool
}

// Placeholder for future integration with agricultural lending
type KisaanMitraKeeper interface {
	// To be implemented when Kisaan Mitra module is created
	GetLoanEligibility(ctx sdk.Context, farmer sdk.AccAddress) (sdk.Coin, error)
	ProcessLoanApplication(ctx sdk.Context, farmer sdk.AccAddress, amount sdk.Coin, duration uint32) error
}