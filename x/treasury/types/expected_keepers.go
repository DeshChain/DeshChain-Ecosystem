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
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// AccountKeeper defines the expected interface for the Account module.
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	SetAccount(ctx sdk.Context, acc types.AccountI)
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, moduleName string) types.ModuleAccountI
	SetModuleAccount(ctx sdk.Context, macc types.ModuleAccountI)
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	LockedCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetSupply(ctx sdk.Context, denom string) sdk.Coin
}

// StakingKeeper defines the expected interface for the Staking module.
type StakingKeeper interface {
	GetValidator(ctx sdk.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, found bool)
	GetBondedValidatorsByPower(ctx sdk.Context) []stakingtypes.Validator
	GetAllValidators(ctx sdk.Context) (validators []stakingtypes.Validator)
	GetValidatorByConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) (validator stakingtypes.Validator, found bool)
	TotalBondedTokens(ctx sdk.Context) sdk.Int
	StakingTokenSupply(ctx sdk.Context) sdk.Int
	GetDelegation(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (delegation stakingtypes.Delegation, found bool)
	GetAllDelegatorDelegations(ctx sdk.Context, delegator sdk.AccAddress) []stakingtypes.Delegation
	GetValidatorDelegations(ctx sdk.Context, valAddr sdk.ValAddress) []stakingtypes.Delegation
}

// DistributionKeeper defines the expected interface for the Distribution module.
type DistributionKeeper interface {
	GetFeePool(ctx sdk.Context) (feePool sdk.DecCoins)
	SetFeePool(ctx sdk.Context, feePool sdk.DecCoins)
	GetCommunityTax(ctx sdk.Context) sdk.Dec
	FundCommunityPool(ctx sdk.Context, amount sdk.Coins, sender sdk.AccAddress) error
	AllocateTokensToValidator(ctx sdk.Context, val stakingtypes.ValidatorI, tokens sdk.DecCoins)
}

// GovernanceKeeper defines the expected interface for the Governance module.
type GovernanceKeeper interface {
	GetProposal(ctx sdk.Context, proposalID uint64) (proposal interface{}, found bool)
	GetProposals(ctx sdk.Context) []interface{}
	SubmitProposal(ctx sdk.Context, content interface{}) (interface{}, error)
	GetVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress) (interface{}, bool)
	GetVotes(ctx sdk.Context, proposalID uint64) []interface{}
	AddVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress, options interface{}, metadata string) error
}

// SlashingKeeper defines the expected interface for the Slashing module.
type SlashingKeeper interface {
	GetValidatorSigningInfo(ctx sdk.Context, consAddr sdk.ConsAddress) (info interface{}, found bool)
	JailUntil(ctx sdk.Context, consAddr sdk.ConsAddress, jailTime int64)
	Slash(ctx sdk.Context, consAddr sdk.ConsAddress, fraction sdk.Dec, power, distributionHeight int64)
	IsTombstoned(ctx sdk.Context, consAddr sdk.ConsAddress) bool
}

// ParamsKeeper defines the expected interface for the Params module.
type ParamsKeeper interface {
	GetSubspace(moduleName string) interface{}
}