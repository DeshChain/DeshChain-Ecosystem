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

package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/namo/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.GetEnableTokenOperations(ctx),
		k.GetTokenDenom(ctx),
		k.GetMaxSupply(ctx),
		k.GetInitialSupply(ctx),
		k.GetMintable(ctx),
		k.GetBurnable(ctx),
		k.GetVestingEnabled(ctx),
		k.GetMinVestingPeriod(ctx),
		k.GetMaxVestingPeriod(ctx),
		k.GetBurnRatio(ctx),
		k.GetDistributionEnabled(ctx),
		k.GetMaxDistributionEvents(ctx),
		k.GetFounderAllocation(ctx),
		k.GetTeamAllocation(ctx),
		k.GetCommunityAllocation(ctx),
		k.GetDevelopmentAllocation(ctx),
		k.GetLiquidityAllocation(ctx),
		k.GetPublicSaleAllocation(ctx),
		k.GetDAOTreasuryAllocation(ctx),
		k.GetCoFounderAllocation(ctx),
		k.GetOperationsAllocation(ctx),
		k.GetAngelAllocation(ctx),
	)
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// GetEnableTokenOperations returns whether token operations are enabled
func (k Keeper) GetEnableTokenOperations(ctx sdk.Context) bool {
	var res bool
	k.paramstore.Get(ctx, types.KeyEnableTokenOperations, &res)
	return res
}

// GetTokenDenom returns the token denomination
func (k Keeper) GetTokenDenom(ctx sdk.Context) string {
	var res string
	k.paramstore.Get(ctx, types.KeyTokenDenom, &res)
	return res
}

// GetMaxSupply returns the maximum token supply
func (k Keeper) GetMaxSupply(ctx sdk.Context) string {
	var res string
	k.paramstore.Get(ctx, types.KeyMaxSupply, &res)
	return res
}

// GetInitialSupply returns the initial token supply
func (k Keeper) GetInitialSupply(ctx sdk.Context) string {
	var res string
	k.paramstore.Get(ctx, types.KeyInitialSupply, &res)
	return res
}

// GetMintable returns whether tokens are mintable
func (k Keeper) GetMintable(ctx sdk.Context) bool {
	var res bool
	k.paramstore.Get(ctx, types.KeyMintable, &res)
	return res
}

// GetBurnable returns whether tokens are burnable
func (k Keeper) GetBurnable(ctx sdk.Context) bool {
	var res bool
	k.paramstore.Get(ctx, types.KeyBurnable, &res)
	return res
}

// GetVestingEnabled returns whether vesting is enabled
func (k Keeper) GetVestingEnabled(ctx sdk.Context) bool {
	var res bool
	k.paramstore.Get(ctx, types.KeyVestingEnabled, &res)
	return res
}

// GetMinVestingPeriod returns the minimum vesting period in seconds
func (k Keeper) GetMinVestingPeriod(ctx sdk.Context) int64 {
	var res int64
	k.paramstore.Get(ctx, types.KeyMinVestingPeriod, &res)
	return res
}

// GetMaxVestingPeriod returns the maximum vesting period in seconds
func (k Keeper) GetMaxVestingPeriod(ctx sdk.Context) int64 {
	var res int64
	k.paramstore.Get(ctx, types.KeyMaxVestingPeriod, &res)
	return res
}

// GetBurnRatio returns the burn ratio for token burning
func (k Keeper) GetBurnRatio(ctx sdk.Context) string {
	var res string
	k.paramstore.Get(ctx, types.KeyBurnRatio, &res)
	return res
}

// GetDistributionEnabled returns whether token distribution is enabled
func (k Keeper) GetDistributionEnabled(ctx sdk.Context) bool {
	var res bool
	k.paramstore.Get(ctx, types.KeyDistributionEnabled, &res)
	return res
}

// GetMaxDistributionEvents returns the maximum number of distribution events
func (k Keeper) GetMaxDistributionEvents(ctx sdk.Context) uint64 {
	var res uint64
	k.paramstore.Get(ctx, types.KeyMaxDistributionEvents, &res)
	return res
}

// GetFounderAllocation returns the founder allocation percentage
func (k Keeper) GetFounderAllocation(ctx sdk.Context) string {
	var res string
	k.paramstore.Get(ctx, types.KeyFounderAllocation, &res)
	return res
}

// GetTeamAllocation returns the team allocation percentage
func (k Keeper) GetTeamAllocation(ctx sdk.Context) string {
	var res string
	k.paramstore.Get(ctx, types.KeyTeamAllocation, &res)
	return res
}

// GetCommunityAllocation returns the community allocation percentage
func (k Keeper) GetCommunityAllocation(ctx sdk.Context) string {
	var res string
	k.paramstore.Get(ctx, types.KeyCommunityAllocation, &res)
	return res
}

// GetDevelopmentAllocation returns the development allocation percentage
func (k Keeper) GetDevelopmentAllocation(ctx sdk.Context) string {
	var res string
	k.paramstore.Get(ctx, types.KeyDevelopmentAllocation, &res)
	return res
}

// GetLiquidityAllocation returns the liquidity allocation percentage
func (k Keeper) GetLiquidityAllocation(ctx sdk.Context) string {
	var res string
	k.paramstore.Get(ctx, types.KeyLiquidityAllocation, &res)
	return res
}

// GetPublicSaleAllocation returns the public sale allocation percentage
func (k Keeper) GetPublicSaleAllocation(ctx sdk.Context) string {
	var res string
	k.paramstore.Get(ctx, types.KeyPublicSaleAllocation, &res)
	return res
}

// GetDAOTreasuryAllocation returns the DAO treasury allocation percentage
func (k Keeper) GetDAOTreasuryAllocation(ctx sdk.Context) string {
	var res string
	k.paramstore.Get(ctx, types.KeyDAOTreasuryAllocation, &res)
	return res
}

// GetCoFounderAllocation returns the co-founder allocation percentage
func (k Keeper) GetCoFounderAllocation(ctx sdk.Context) string {
	var res string
	k.paramstore.Get(ctx, types.KeyCoFounderAllocation, &res)
	return res
}

// GetOperationsAllocation returns the operations allocation percentage
func (k Keeper) GetOperationsAllocation(ctx sdk.Context) string {
	var res string
	k.paramstore.Get(ctx, types.KeyOperationsAllocation, &res)
	return res
}

// GetAngelAllocation returns the angel allocation percentage
func (k Keeper) GetAngelAllocation(ctx sdk.Context) string {
	var res string
	k.paramstore.Get(ctx, types.KeyAngelAllocation, &res)
	return res
}