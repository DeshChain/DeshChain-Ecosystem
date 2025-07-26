package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/DeshChain/DeshChain-Ecosystem/x/dinr/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

// UserPosition returns a specific user's position
func (k Keeper) UserPosition(c context.Context, req *types.QueryUserPositionRequest) (*types.QueryUserPositionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	position, found := k.GetUserPosition(ctx, req.Address)
	if !found {
		return nil, status.Error(codes.NotFound, "position not found")
	}

	return &types.QueryUserPositionResponse{Position: &position}, nil
}

// AllPositions returns all user positions with pagination
func (k Keeper) AllPositions(c context.Context, req *types.QueryAllPositionsRequest) (*types.QueryAllPositionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var positions []types.UserPosition
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	positionStore := prefix.NewStore(store, types.UserPositionPrefix)

	pageRes, err := query.Paginate(positionStore, req.Pagination, func(key []byte, value []byte) error {
		var position types.UserPosition
		if err := k.cdc.Unmarshal(value, &position); err != nil {
			return err
		}

		positions = append(positions, position)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllPositionsResponse{
		Positions:  positions,
		Pagination: pageRes,
	}, nil
}

// CollateralAsset returns information about a specific collateral asset
func (k Keeper) CollateralAsset(c context.Context, req *types.QueryCollateralAssetRequest) (*types.QueryCollateralAssetResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	asset, found := k.GetCollateralAsset(ctx, req.Denom)
	if !found {
		return nil, status.Error(codes.NotFound, "collateral asset not found")
	}

	return &types.QueryCollateralAssetResponse{Asset: &asset}, nil
}

// AllCollateralAssets returns all supported collateral assets
func (k Keeper) AllCollateralAssets(c context.Context, req *types.QueryAllCollateralAssetsRequest) (*types.QueryAllCollateralAssetsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	assets := k.GetAllCollateralAssets(ctx)

	return &types.QueryAllCollateralAssetsResponse{Assets: assets}, nil
}

// StabilityInfo returns current stability metrics
func (k Keeper) StabilityInfo(c context.Context, req *types.QueryStabilityInfoRequest) (*types.QueryStabilityInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	stabilityData := k.GetStabilityData(ctx)

	return &types.QueryStabilityInfoResponse{StabilityData: &stabilityData}, nil
}

// InsuranceFund returns insurance fund information
func (k Keeper) InsuranceFund(c context.Context, req *types.QueryInsuranceFundRequest) (*types.QueryInsuranceFundResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	fund := k.GetInsuranceFund(ctx)

	return &types.QueryInsuranceFundResponse{Fund: &fund}, nil
}

// EstimateMint estimates the result of a mint operation
func (k Keeper) EstimateMint(c context.Context, req *types.QueryEstimateMintRequest) (*types.QueryEstimateMintResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	// Parse amounts
	collateralAmount, err := sdk.ParseCoinNormalized(req.CollateralAmount)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid collateral amount")
	}

	dinrAmount, err := sdk.ParseCoinNormalized(req.DinrAmount)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid DINR amount")
	}

	// Validate collateral is acceptable
	collateralAsset, found := k.GetCollateralAsset(ctx, req.CollateralDenom)
	if !found || !collateralAsset.IsEnabled {
		return nil, status.Error(codes.InvalidArgument, "invalid collateral asset")
	}

	// Get collateral price
	price, err := k.oracleKeeper.GetPrice(ctx, req.CollateralDenom)
	if err != nil {
		return nil, status.Error(codes.Internal, "could not fetch collateral price")
	}

	// Calculate collateral value in INR
	collateralValue := k.calculateCollateralValue(collateralAmount, price)

	// Calculate collateral ratio
	collateralRatio := k.calculateCollateralRatio(collateralValue, dinrAmount.Amount)

	// Get params
	params := k.GetParams(ctx)

	// Check if ratio meets minimum
	canMint := collateralRatio >= uint64(params.MinCollateralRatio)

	// Calculate fees
	mintFee := k.calculateMintingFee(dinrAmount, params.Fees)
	netDINR := dinrAmount.Sub(mintFee)

	return &types.QueryEstimateMintResponse{
		CanMint:         canMint,
		CollateralValue: sdk.NewCoin("inr", collateralValue),
		CollateralRatio: collateralRatio,
		MintFee:         mintFee,
		NetDinrToMint:   netDINR,
	}, nil
}

// EstimateBurn estimates the result of a burn operation
func (k Keeper) EstimateBurn(c context.Context, req *types.QueryEstimateBurnRequest) (*types.QueryEstimateBurnResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	// Parse amount
	dinrAmount, err := sdk.ParseCoinNormalized(req.DinrAmount)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid DINR amount")
	}

	// Get user position
	position, found := k.GetUserPosition(ctx, req.Burner)
	if !found {
		return nil, status.Error(codes.NotFound, "position not found")
	}

	// Check if user has enough DINR minted
	if position.DinrMinted.Amount.LT(dinrAmount.Amount) {
		return nil, status.Error(codes.InvalidArgument, "insufficient DINR minted")
	}

	// Get params
	params := k.GetParams(ctx)

	// Calculate fees
	burnFee := k.calculateBurningFee(dinrAmount, params.Fees)
	totalDINRNeeded := dinrAmount.Add(burnFee)

	// Calculate collateral to return
	collateralToReturn := k.calculateCollateralToReturn(ctx, position, dinrAmount, req.CollateralDenom)

	return &types.QueryEstimateBurnResponse{
		BurnFee:            burnFee,
		TotalDinrNeeded:    totalDINRNeeded,
		CollateralToReturn: collateralToReturn,
	}, nil
}

// LiquidatablePositions returns positions eligible for liquidation
func (k Keeper) LiquidatablePositions(c context.Context, req *types.QueryLiquidatablePositionsRequest) (*types.QueryLiquidatablePositionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var liquidatablePositions []types.UserPosition
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)

	store := ctx.KVStore(k.storeKey)
	positionStore := prefix.NewStore(store, types.UserPositionPrefix)

	pageRes, err := query.Paginate(positionStore, req.Pagination, func(key []byte, value []byte) error {
		var position types.UserPosition
		if err := k.cdc.Unmarshal(value, &position); err != nil {
			return err
		}

		// Check if position is liquidatable
		collateralValue := k.calculateTotalCollateralValue(ctx, position.Collateral)
		collateralRatio := k.calculateCollateralRatio(collateralValue, position.DinrMinted.Amount)

		if collateralRatio < uint64(params.LiquidationThreshold) {
			position.IsLiquidatable = true
			position.CollateralRatio = collateralRatio
			liquidatablePositions = append(liquidatablePositions, position)
		}

		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryLiquidatablePositionsResponse{
		Positions:  liquidatablePositions,
		Pagination: pageRes,
	}, nil
}