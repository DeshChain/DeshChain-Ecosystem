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
	"context"
	"fmt"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/deshchain/x/moneyorder/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// CreateMoneyOrder handles creation of a new money order
func (k msgServer) CreateMoneyOrder(
	goCtx context.Context,
	msg *types.MsgCreateMoneyOrder,
) (*types.MsgCreateMoneyOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Parse sender address
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	
	// Create the money order
	orderId, referenceNumber, err := k.Keeper.CreateMoneyOrder(
		ctx,
		sender,
		msg.ReceiverUPI,
		msg.Amount,
		msg.Note,
		msg.OrderType,
		msg.ScheduledTime,
	)
	if err != nil {
		return nil, err
	}
	
	return &types.MsgCreateMoneyOrderResponse{
		OrderId:         orderId,
		ReferenceNumber: referenceNumber,
	}, nil
}

// CreateFixedRatePool handles creation of a fixed rate pool
func (k msgServer) CreateFixedRatePool(
	goCtx context.Context,
	msg *types.MsgCreateFixedRatePool,
) (*types.MsgCreateFixedRatePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Parse creator address
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}
	
	// Create the pool
	poolId, err := k.Keeper.CreateFixedRatePool(
		ctx,
		creator,
		msg.Token0Denom,
		msg.Token1Denom,
		msg.ExchangeRate,
		msg.InitialLiquidity,
		msg.Description,
		msg.SupportedRegions,
	)
	if err != nil {
		return nil, err
	}
	
	return &types.MsgCreateFixedRatePoolResponse{
		PoolId: poolId,
	}, nil
}

// CreateVillagePool handles creation of a village pool
func (k msgServer) CreateVillagePool(
	goCtx context.Context,
	msg *types.MsgCreateVillagePool,
) (*types.MsgCreateVillagePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Parse panchayat head address
	panchayatHead, err := sdk.AccAddressFromBech32(msg.PanchayatHead)
	if err != nil {
		return nil, err
	}
	
	// Parse validator addresses
	var validators []sdk.ValAddress
	for _, valStr := range msg.LocalValidators {
		val, err := sdk.ValAddressFromBech32(valStr)
		if err != nil {
			return nil, fmt.Errorf("invalid validator address %s: %w", valStr, err)
		}
		validators = append(validators, val)
	}
	
	// Create the pool
	poolId, err := k.Keeper.CreateVillagePool(
		ctx,
		panchayatHead,
		msg.VillageName,
		msg.PostalCode,
		msg.StateCode,
		msg.DistrictCode,
		msg.InitialLiquidity,
		validators,
	)
	if err != nil {
		return nil, err
	}
	
	return &types.MsgCreateVillagePoolResponse{
		PoolId: poolId,
	}, nil
}

// AddLiquidity handles adding liquidity to a pool
func (k msgServer) AddLiquidity(
	goCtx context.Context,
	msg *types.MsgAddLiquidity,
) (*types.MsgAddLiquidityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Parse depositor address
	depositor, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		return nil, err
	}
	
	// Determine pool type and add liquidity
	// First check if it's a village pool
	villagePool, found := k.GetVillagePoolById(ctx, msg.PoolId)
	if found {
		err := k.AddVillagePoolLiquidity(
			ctx,
			msg.PoolId,
			depositor,
			msg.TokenAmounts,
		)
		if err != nil {
			return nil, err
		}
		
		// Calculate shares (simplified)
		shares := sdk.NewInt(0)
		for _, coin := range msg.TokenAmounts {
			shares = shares.Add(coin.Amount)
		}
		
		return &types.MsgAddLiquidityResponse{
			SharesOut: shares,
		}, nil
	}
	
	// Check if it's a fixed rate pool
	fixedPool, found := k.GetFixedRatePool(ctx, msg.PoolId)
	if found {
		// Fixed rate pools don't have traditional liquidity adding
		// This would be handled differently in production
		return nil, fmt.Errorf("fixed rate pools don't support direct liquidity addition")
	}
	
	// If neither type found
	return nil, types.ErrPoolNotFound
}

// RemoveLiquidity handles removing liquidity from a pool
func (k msgServer) RemoveLiquidity(
	goCtx context.Context,
	msg *types.MsgRemoveLiquidity,
) (*types.MsgRemoveLiquidityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Parse withdrawer address
	withdrawer, err := sdk.AccAddressFromBech32(msg.Withdrawer)
	if err != nil {
		return nil, err
	}
	
	// Implementation would handle liquidity removal
	// For now, return placeholder
	return &types.MsgRemoveLiquidityResponse{
		TokensOut: sdk.NewCoins(),
	}, fmt.Errorf("liquidity removal not yet implemented")
}

// SwapExactAmountIn handles swap with exact input amount
func (k msgServer) SwapExactAmountIn(
	goCtx context.Context,
	msg *types.MsgSwapExactAmountIn,
) (*types.MsgSwapExactAmountInResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Parse sender address
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	
	// Check if it's a fixed rate pool
	_, found := k.GetFixedRatePool(ctx, msg.PoolId)
	if found {
		// Execute fixed rate swap
		tokenOut, err := k.ExecuteFixedRateSwap(
			ctx,
			msg.PoolId,
			sender,
			msg.TokenIn,
			msg.TokenOutDenom,
			msg.TokenOutMin,
		)
		if err != nil {
			return nil, err
		}
		
		return &types.MsgSwapExactAmountInResponse{
			TokenOut: tokenOut,
		}, nil
	}
	
	// Check other pool types (AMM, etc.) - to be implemented
	return nil, types.ErrPoolNotFound
}

// SwapExactAmountOut handles swap with exact output amount
func (k msgServer) SwapExactAmountOut(
	goCtx context.Context,
	msg *types.MsgSwapExactAmountOut,
) (*types.MsgSwapExactAmountOutResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Parse sender address
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	
	// Implementation would handle exact output swaps
	// For now, return placeholder
	return &types.MsgSwapExactAmountOutResponse{
		TokenIn: sdk.Coin{},
	}, fmt.Errorf("exact amount out swaps not yet implemented")
}

// JoinVillagePool handles joining a village pool
func (k msgServer) JoinVillagePool(
	goCtx context.Context,
	msg *types.MsgJoinVillagePool,
) (*types.MsgJoinVillagePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Parse member address
	member, err := sdk.AccAddressFromBech32(msg.Member)
	if err != nil {
		return nil, err
	}
	
	// Join the village pool
	err = k.Keeper.JoinVillagePool(
		ctx,
		member,
		msg.PoolId,
		msg.InitialDeposit,
		msg.LocalName,
		msg.MobileNumber,
	)
	if err != nil {
		return nil, err
	}
	
	// Generate member ID
	memberId := fmt.Sprintf("MEMBER-%d-%s", msg.PoolId, member.String()[:8])
	
	return &types.MsgJoinVillagePoolResponse{
		MemberId: memberId,
	}, nil
}

// ClaimRewards handles claiming rewards from a pool
func (k msgServer) ClaimRewards(
	goCtx context.Context,
	msg *types.MsgClaimRewards,
) (*types.MsgClaimRewardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Parse claimer address
	claimer, err := sdk.AccAddressFromBech32(msg.Claimer)
	if err != nil {
		return nil, err
	}
	
	// Check if it's a village pool
	_, found := k.GetVillagePoolById(ctx, msg.PoolId)
	if found {
		rewards, err := k.ClaimVillagePoolRewards(
			ctx,
			msg.PoolId,
			claimer,
		)
		if err != nil {
			return nil, err
		}
		
		return &types.MsgClaimRewardsResponse{
			RewardsClaimed: rewards,
		}, nil
	}
	
	// Check other pool types
	return nil, types.ErrPoolNotFound
}

// UpdatePoolParams handles updating pool parameters (governance only)
func (k msgServer) UpdatePoolParams(
	goCtx context.Context,
	msg *types.MsgUpdatePoolParams,
) (*types.MsgUpdatePoolParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Parse authority address
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return nil, err
	}
	
	// Check if authority is governance module
	// In production, this would check against the governance module address
	// For now, simplified check
	if authority.String() != "desh1gov..." { // Placeholder
		return nil, types.ErrUnauthorized
	}
	
	// Update pool parameters based on pool type
	_, found := k.GetFixedRatePool(ctx, msg.PoolId)
	if found {
		err := k.UpdateFixedRatePoolParams(
			ctx,
			msg.PoolId,
			msg.BaseFee,
			msg.Active,
		)
		if err != nil {
			return nil, err
		}
		
		return &types.MsgUpdatePoolParamsResponse{
			Success: true,
		}, nil
	}
	
	return nil, types.ErrPoolNotFound
}