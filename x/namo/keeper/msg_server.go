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

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/namo/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// BurnTokens handles token burning requests
func (k msgServer) BurnTokens(goCtx context.Context, msg *types.MsgBurnTokens) (*types.MsgBurnTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate message
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	// Parse sender address
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	// Parse amount
	amount, ok := sdk.NewIntFromString(msg.Amount)
	if !ok {
		return nil, types.ErrInvalidTokenAmount
	}

	// Validate token operation
	if err := k.ValidateTokenOperation(ctx, "burn", amount); err != nil {
		return nil, err
	}

	// Burn tokens
	if err := k.BurnTokens(ctx, senderAddr, amount); err != nil {
		return nil, err
	}

	return &types.MsgBurnTokensResponse{}, nil
}

// ClaimVestedTokens handles vested token claiming requests
func (k msgServer) ClaimVestedTokens(goCtx context.Context, msg *types.MsgClaimVestedTokens) (*types.MsgClaimVestedTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate message
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	// Claim vested tokens
	if err := k.ClaimVestedTokens(ctx, msg.Beneficiary); err != nil {
		return nil, err
	}

	// Calculate claimed amount for response
	claimedAmount, err := k.CalculateVestedTokens(ctx, msg.Beneficiary)
	if err != nil {
		claimedAmount = sdk.ZeroInt()
	}

	return &types.MsgClaimVestedTokensResponse{
		ClaimedAmount: claimedAmount.String(),
	}, nil
}

// CreateVestingSchedule handles vesting schedule creation requests
func (k msgServer) CreateVestingSchedule(goCtx context.Context, msg *types.MsgCreateVestingSchedule) (*types.MsgCreateVestingScheduleResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate message
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	// Check if vesting schedule already exists for beneficiary
	if _, found := k.GetVestingSchedule(ctx, msg.Beneficiary); found {
		return nil, types.ErrVestingScheduleExists
	}

	// Validate vesting parameters
	params := k.GetParams(ctx)
	if !params.VestingEnabled {
		return nil, types.ErrVestingDisabled
	}

	// Validate vesting period
	vestingPeriod := msg.EndTime - msg.CliffTime
	if vestingPeriod < params.MinVestingPeriod {
		return nil, types.ErrInvalidVestingPeriod
	}
	if vestingPeriod > params.MaxVestingPeriod {
		return nil, types.ErrInvalidVestingPeriod
	}

	// Create vesting schedule
	schedule := types.VestingSchedule{
		Beneficiary:     msg.Beneficiary,
		TotalAmount:     msg.TotalAmount,
		ClaimedAmount:   "0",
		CliffTime:       msg.CliffTime,
		EndTime:         msg.EndTime,
		VestingCategory: msg.VestingCategory,
		CreatedAt:       ctx.BlockTime().Unix(),
	}

	// Validate schedule
	if err := types.ValidateVestingSchedule(schedule); err != nil {
		return nil, err
	}

	// Set vesting schedule
	k.SetVestingSchedule(ctx, schedule)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateVestingSchedule,
			sdk.NewAttribute(types.AttributeKeyBeneficiary, msg.Beneficiary),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.TotalAmount),
			sdk.NewAttribute(types.AttributeKeyVestingCategory, msg.VestingCategory),
		),
	)

	return &types.MsgCreateVestingScheduleResponse{}, nil
}

// UpdateParams handles parameter update requests
func (k msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate message
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	// TODO: Implement proper authority validation
	// For now, we'll skip authority validation, but in production this should check
	// that the sender has the appropriate governance authority

	// Validate params
	if err := msg.Params.Validate(); err != nil {
		return nil, err
	}

	// Set new params
	k.SetParams(ctx, msg.Params)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpdateParams,
			sdk.NewAttribute(types.AttributeKeyAuthority, msg.Authority),
		),
	)

	return &types.MsgUpdateParamsResponse{}, nil
}