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

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/deshchain/deshchain/x/namo/types"
)

var _ types.QueryServer = Keeper{}

// Params returns the module parameters
func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}

// TokenSupply returns the current token supply
func (k Keeper) TokenSupply(c context.Context, req *types.QueryTokenSupplyRequest) (*types.QueryTokenSupplyResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	supply, found := k.GetTokenSupply(ctx)
	if !found {
		return nil, status.Error(codes.NotFound, "token supply not found")
	}

	return &types.QueryTokenSupplyResponse{TokenSupply: supply}, nil
}

// VestingSchedule returns a vesting schedule by beneficiary
func (k Keeper) VestingSchedule(c context.Context, req *types.QueryVestingScheduleRequest) (*types.QueryVestingScheduleResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	schedule, found := k.GetVestingSchedule(ctx, req.Beneficiary)
	if !found {
		return nil, status.Error(codes.NotFound, "vesting schedule not found")
	}

	return &types.QueryVestingScheduleResponse{VestingSchedule: schedule}, nil
}

// VestingSchedules returns all vesting schedules with pagination
func (k Keeper) VestingSchedules(c context.Context, req *types.QueryVestingSchedulesRequest) (*types.QueryVestingSchedulesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	schedules := k.GetAllVestingSchedules(ctx)

	// Apply pagination
	pageRes, err := query.Paginate(len(schedules), req.Pagination, func(offset, limit int) error {
		end := offset + limit
		if end > len(schedules) {
			end = len(schedules)
		}
		schedules = schedules[offset:end]
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryVestingSchedulesResponse{
		VestingSchedules: schedules,
		Pagination:       pageRes,
	}, nil
}

// VestedTokens returns the amount of vested tokens for a beneficiary
func (k Keeper) VestedTokens(c context.Context, req *types.QueryVestedTokensRequest) (*types.QueryVestedTokensResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	vestedAmount, err := k.CalculateVestedTokens(ctx, req.Beneficiary)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryVestedTokensResponse{
		VestedAmount: vestedAmount.String(),
	}, nil
}

// TokenDistributionEvent returns a token distribution event by ID
func (k Keeper) TokenDistributionEvent(c context.Context, req *types.QueryTokenDistributionEventRequest) (*types.QueryTokenDistributionEventResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	event, found := k.GetTokenDistributionEvent(ctx, req.Id)
	if !found {
		return nil, status.Error(codes.NotFound, "token distribution event not found")
	}

	return &types.QueryTokenDistributionEventResponse{Event: event}, nil
}

// TokenDistributionEvents returns all token distribution events with pagination
func (k Keeper) TokenDistributionEvents(c context.Context, req *types.QueryTokenDistributionEventsRequest) (*types.QueryTokenDistributionEventsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	events := k.GetAllTokenDistributionEvents(ctx)

	// Apply pagination
	pageRes, err := query.Paginate(len(events), req.Pagination, func(offset, limit int) error {
		end := offset + limit
		if end > len(events) {
			end = len(events)
		}
		events = events[offset:end]
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryTokenDistributionEventsResponse{
		Events:     events,
		Pagination: pageRes,
	}, nil
}

// CirculatingSupply returns the circulating token supply
func (k Keeper) CirculatingSupply(c context.Context, req *types.QueryCirculatingSupplyRequest) (*types.QueryCirculatingSupplyResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	circulatingSupply := k.GetCirculatingSupply(ctx)

	return &types.QueryCirculatingSupplyResponse{
		CirculatingSupply: circulatingSupply.String(),
	}, nil
}

// BurnedTokens returns the total amount of burned tokens
func (k Keeper) BurnedTokens(c context.Context, req *types.QueryBurnedTokensRequest) (*types.QueryBurnedTokensResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	burnedTokens := k.GetTotalBurnedTokens(ctx)

	return &types.QueryBurnedTokensResponse{
		BurnedAmount: burnedTokens.String(),
	}, nil
}

// TokenStats returns comprehensive token statistics
func (k Keeper) TokenStats(c context.Context, req *types.QueryTokenStatsRequest) (*types.QueryTokenStatsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	// Get token supply
	supply, found := k.GetTokenSupply(ctx)
	if !found {
		return nil, status.Error(codes.NotFound, "token supply not found")
	}

	// Get circulating supply
	circulatingSupply := k.GetCirculatingSupply(ctx)

	// Get burned tokens
	burnedTokens := k.GetTotalBurnedTokens(ctx)

	// Get vesting pool balance
	vestingPoolAddr := k.accountKeeper.GetModuleAddress(types.VestingPoolName)
	vestingTokens := sdk.ZeroInt()
	if vestingPoolAddr != nil {
		vestingBalance := k.bankKeeper.GetBalance(ctx, vestingPoolAddr, types.TokenDenom)
		vestingTokens = vestingBalance.Amount
	}

	// Count vesting schedules
	schedules := k.GetAllVestingSchedules(ctx)
	totalVestingSchedules := uint64(len(schedules))

	// Count distribution events
	events := k.GetAllTokenDistributionEvents(ctx)
	totalDistributionEvents := uint64(len(events))

	return &types.QueryTokenStatsResponse{
		TotalSupply:               supply.TotalSupply,
		CirculatingSupply:         circulatingSupply.String(),
		BurnedTokens:              burnedTokens.String(),
		VestingTokens:             vestingTokens.String(),
		TotalVestingSchedules:     totalVestingSchedules,
		TotalDistributionEvents:   totalDistributionEvents,
	}, nil
}