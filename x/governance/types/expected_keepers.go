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
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

// AccountKeeper defines the expected account keeper
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
}

// BankKeeper defines the expected bank keeper
type BankKeeper interface {
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
}

// GovKeeper defines the expected governance keeper interface
type GovKeeper interface {
	GetProposal(ctx sdk.Context, proposalID uint64) (govv1.Proposal, bool)
	SetProposal(ctx sdk.Context, proposal govv1.Proposal)
	GetProposals(ctx sdk.Context) (proposals govv1.Proposals)
	GetVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress) (vote govv1.Vote, found bool)
	GetVotes(ctx sdk.Context, proposalID uint64) (votes govv1.Votes)
	GetDeposit(ctx sdk.Context, proposalID uint64, depositorAddr sdk.AccAddress) (deposit govv1.Deposit, found bool)
	GetDeposits(ctx sdk.Context, proposalID uint64) (deposits govv1.Deposits)
	GetTallyResult(ctx sdk.Context, proposal govv1.Proposal) (tallyResult govv1.TallyResult)
}

// StakingKeeper defines the expected staking keeper
type StakingKeeper interface {
	TotalBondedTokens(ctx sdk.Context) sdk.Int
	IterateBondedValidatorsByPower(ctx sdk.Context, fn func(index int64, validator StakingValidatorI) (stop bool))
}

// StakingValidatorI defines the expected validator interface
type StakingValidatorI interface {
	GetOperator() sdk.ValAddress
	GetTokens() sdk.Int
	GetBondedTokens() sdk.Int
	GetDelegatorShares() sdk.Dec
	IsBonded() bool
}