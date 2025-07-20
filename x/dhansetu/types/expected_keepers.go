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
	GetModuleAddress(moduleName string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, moduleName string) types.ModuleAccountI
	SetModuleAccount(ctx sdk.Context, moduleAccount types.ModuleAccountI)
}

// BankKeeper defines the expected bank keeper
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
}

// MoneyOrderKeeper defines the expected money order keeper interface
type MoneyOrderKeeper interface {
	// Core money order functions
	CreateMoneyOrder(ctx sdk.Context, sender sdk.AccAddress, receiverUPI string, amount sdk.Coin, note string) error
	GetMoneyOrder(ctx sdk.Context, orderId string) (interface{}, bool) // Using interface{} to avoid circular imports
	
	// P2P matching functions
	CreateP2POrder(ctx sdk.Context, order interface{}) error
	MatchOrders(ctx sdk.Context, buyOrderId, sellOrderId string) error
	
	// Escrow functions
	CreateEscrow(ctx sdk.Context, orderId string, amount sdk.Coin, expiryTime int64) error
	ReleaseEscrow(ctx sdk.Context, escrowId string) error
	RefundEscrow(ctx sdk.Context, escrowId string) error
	
	// Mitra functions
	RegisterSevaMitra(ctx sdk.Context, mitra interface{}) error
	GetSevaMitra(ctx sdk.Context, mitraId string) (interface{}, bool)
	UpdateMitraStats(ctx sdk.Context, mitraId string, tradeVolume sdk.Int, successful bool) error
}

// CulturalKeeper defines the expected cultural keeper interface
type CulturalKeeper interface {
	GetRandomQuote(ctx sdk.Context, language string) string
	GetFestivalBonus(ctx sdk.Context, festivalName string) sdk.Dec
	IsFestivalPeriod(ctx sdk.Context, festivalName string) bool
	GetPatriotismScore(ctx sdk.Context, address string) int64
}

// NAMOKeeper defines the expected NAMO token keeper interface
type NAMOKeeper interface {
	GetTokenSupply(ctx sdk.Context) interface{} // Using interface{} to avoid circular imports
	BurnTokens(ctx sdk.Context, from sdk.AccAddress, amount sdk.Coin) error
	CreateVestingSchedule(ctx sdk.Context, beneficiary sdk.AccAddress, amount sdk.Int, cliffMonths, vestingMonths int64) error
	ClaimVestedTokens(ctx sdk.Context, beneficiary sdk.AccAddress) (sdk.Coin, error)
	GetVestingSchedule(ctx sdk.Context, beneficiary sdk.AccAddress) (interface{}, bool)
}

// ParamSubspace defines the expected Subspace interface
type ParamSubspace interface {
	Get(ctx sdk.Context, key []byte, ptr interface{})
	Set(ctx sdk.Context, key []byte, param interface{})
	GetParamSet(ctx sdk.Context, ps ParamSet)
	SetParamSet(ctx sdk.Context, ps ParamSet)
}

// ParamSet defines an interface for structs containing parameters
type ParamSet interface {
	ParamSetPairs() ParamSetPairs
}

// ParamSetPair is a key-value pair for module parameters
type ParamSetPair struct {
	Key         []byte
	Value       interface{}
	ValidatorFn func(value interface{}) error
}

// ParamSetPairs represents a list of ParamSetPair
type ParamSetPairs []ParamSetPair