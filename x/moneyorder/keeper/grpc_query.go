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

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/DeshChain/DeshChain-Ecosystem/x/moneyorder/types"
)

var _ types.QueryServer = Keeper{}

// Params queries the parameters of the Money Order module
func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}

// FixedRatePool queries a fixed rate pool by ID
func (k Keeper) FixedRatePool(c context.Context, req *types.QueryFixedRatePoolRequest) (*types.QueryFixedRatePoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	pool, found := k.GetFixedRatePool(ctx, req.PoolId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "fixed rate pool %d not found", req.PoolId)
	}

	return &types.QueryFixedRatePoolResponse{Pool: pool}, nil
}

// FixedRatePools queries all fixed rate pools with pagination
func (k Keeper) FixedRatePools(c context.Context, req *types.QueryFixedRatePoolsRequest) (*types.QueryFixedRatePoolsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var pools []types.FixedRatePool
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	poolStore := prefix.NewStore(store, types.KeyPrefixFixedRatePool)

	pageRes, err := query.Paginate(poolStore, req.Pagination, func(key []byte, value []byte) error {
		var pool types.FixedRatePool
		if err := k.cdc.Unmarshal(value, &pool); err != nil {
			return err
		}
		pools = append(pools, pool)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryFixedRatePoolsResponse{
		Pools:      pools,
		Pagination: pageRes,
	}, nil
}

// VillagePool queries a village pool by ID
func (k Keeper) VillagePool(c context.Context, req *types.QueryVillagePoolRequest) (*types.QueryVillagePoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	pool, found := k.GetVillagePool(ctx, req.PoolId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "village pool %d not found", req.PoolId)
	}

	return &types.QueryVillagePoolResponse{Pool: pool}, nil
}

// VillagePools queries all village pools with pagination
func (k Keeper) VillagePools(c context.Context, req *types.QueryVillagePoolsRequest) (*types.QueryVillagePoolsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var pools []types.VillagePool
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	poolStore := prefix.NewStore(store, types.KeyPrefixVillagePool)

	pageRes, err := query.Paginate(poolStore, req.Pagination, func(key []byte, value []byte) error {
		var pool types.VillagePool
		if err := k.cdc.Unmarshal(value, &pool); err != nil {
			return err
		}
		pools = append(pools, pool)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryVillagePoolsResponse{
		Pools:      pools,
		Pagination: pageRes,
	}, nil
}

// MoneyOrderReceipt queries a money order receipt by ID
func (k Keeper) MoneyOrderReceipt(c context.Context, req *types.QueryMoneyOrderReceiptRequest) (*types.QueryMoneyOrderReceiptResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	receipt, found := k.GetMoneyOrderReceipt(ctx, req.OrderId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "money order receipt %s not found", req.OrderId)
	}

	return &types.QueryMoneyOrderReceiptResponse{Receipt: receipt}, nil
}

// UserReceipts queries all receipts for a user
func (k Keeper) UserReceipts(c context.Context, req *types.QueryUserReceiptsRequest) (*types.QueryUserReceiptsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	receipts := k.GetUserReceipts(ctx, req.Address)

	// Apply pagination manually since receipts are stored differently
	start, end := 0, len(receipts)
	if req.Pagination != nil {
		if req.Pagination.Offset > uint64(len(receipts)) {
			receipts = []types.MoneyOrderReceipt{}
		} else {
			start = int(req.Pagination.Offset)
			limit := int(req.Pagination.Limit)
			if limit == 0 {
				limit = 100 // Default limit
			}
			if start+limit < len(receipts) {
				end = start + limit
			}
			receipts = receipts[start:end]
		}
	}

	return &types.QueryUserReceiptsResponse{
		Receipts: receipts,
		Pagination: &query.PageResponse{
			Total: uint64(len(receipts)),
		},
	}, nil
}

// EstimateSwap estimates the output for a swap
func (k Keeper) EstimateSwap(c context.Context, req *types.QueryEstimateSwapRequest) (*types.QueryEstimateSwapResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	
	// Parse the amount
	tokenIn, err := sdk.ParseCoinNormalized(req.TokenIn)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Find a suitable pool
	pool, found := k.findBestPool(ctx, tokenIn.Denom, req.TokenOutDenom)
	if !found {
		return nil, status.Error(codes.NotFound, "no suitable pool found for this token pair")
	}

	// Calculate output amount
	outputAmount, fee := k.calculateSwapOutput(ctx, pool, tokenIn)
	
	return &types.QueryEstimateSwapResponse{
		TokenOut: sdk.NewCoin(req.TokenOutDenom, outputAmount),
		Fee:      sdk.NewCoin(tokenIn.Denom, fee),
		PoolId:   pool.PoolId,
		Rate:     pool.ExchangeRate,
	}, nil
}

// UserUPIAddress queries the UPI address for a user
func (k Keeper) UserUPIAddress(c context.Context, req *types.QueryUserUPIAddressRequest) (*types.QueryUserUPIAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	upiAddress, found := k.GetUPIAddress(ctx, req.Address)
	if !found {
		return nil, status.Errorf(codes.NotFound, "UPI address not found for %s", req.Address)
	}

	return &types.QueryUserUPIAddressResponse{UpiAddress: upiAddress}, nil
}

// AddressFromUPI queries the blockchain address from UPI address
func (k Keeper) AddressFromUPI(c context.Context, req *types.QueryAddressFromUPIRequest) (*types.QueryAddressFromUPIResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	address, found := k.GetAddressFromUPI(ctx, req.UpiAddress)
	if !found {
		return nil, status.Errorf(codes.NotFound, "address not found for UPI %s", req.UpiAddress)
	}

	return &types.QueryAddressFromUPIResponse{Address: address}, nil
}

// VillagePoolMembers queries members of a village pool
func (k Keeper) VillagePoolMembers(c context.Context, req *types.QueryVillagePoolMembersRequest) (*types.QueryVillagePoolMembersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	members := k.GetVillagePoolMembers(ctx, req.PoolId)

	// Apply pagination
	start, end := 0, len(members)
	if req.Pagination != nil {
		if req.Pagination.Offset > uint64(len(members)) {
			members = []types.VillagePoolMember{}
		} else {
			start = int(req.Pagination.Offset)
			limit := int(req.Pagination.Limit)
			if limit == 0 {
				limit = 100
			}
			if start+limit < len(members) {
				end = start + limit
			}
			members = members[start:end]
		}
	}

	return &types.QueryVillagePoolMembersResponse{
		Members: members,
		Pagination: &query.PageResponse{
			Total: uint64(len(members)),
		},
	}, nil
}

// Helper function to find the best pool for a token pair
func (k Keeper) findBestPool(ctx sdk.Context, tokenIn, tokenOut string) (types.FixedRatePool, bool) {
	var bestPool types.FixedRatePool
	found := false
	
	k.IterateFixedRatePools(ctx, func(pool types.FixedRatePool) bool {
		if pool.Active && !pool.MaintenanceMode {
			if (pool.Token0Denom == tokenIn && pool.Token1Denom == tokenOut) ||
			   (pool.Token1Denom == tokenIn && pool.Token0Denom == tokenOut) {
				bestPool = pool
				found = true
				return true // stop iteration
			}
		}
		return false
	})
	
	return bestPool, found
}

// Helper function to calculate swap output
func (k Keeper) calculateSwapOutput(ctx sdk.Context, pool types.FixedRatePool, tokenIn sdk.Coin) (sdk.Int, sdk.Int) {
	params := k.GetParams(ctx)
	
	// Calculate base output
	var outputAmount sdk.Int
	if pool.Token0Denom == tokenIn.Denom {
		outputAmount = tokenIn.Amount.ToDec().Mul(pool.ExchangeRate).TruncateInt()
	} else {
		outputAmount = tokenIn.Amount.ToDec().Mul(pool.ReverseRate).TruncateInt()
	}
	
	// Calculate fee
	fee := tokenIn.Amount.ToDec().Mul(params.TradingFeeRate).TruncateInt()
	
	// Apply cultural discount if applicable
	if k.IsInFestivalPeriod(ctx) && params.EnableFestivalBonuses {
		discount := fee.ToDec().Mul(params.FestivalDiscount).TruncateInt()
		fee = fee.Sub(discount)
	}
	
	// Deduct fee from output
	outputAmount = outputAmount.Sub(fee)
	
	return outputAmount, fee
}