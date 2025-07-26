package keeper

import (
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/dinr/types"
	"github.com/gogo/protobuf/proto"
)

// GetUserPosition returns a user's DINR position
func (k Keeper) GetUserPosition(ctx sdk.Context, address string) (types.UserPosition, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.UserPositionPrefix)
	
	bz := store.Get([]byte(address))
	if bz == nil {
		return types.UserPosition{}, false
	}
	
	var position types.UserPosition
	k.cdc.MustUnmarshal(bz, &position)
	return position, true
}

// SetUserPosition sets a user's DINR position
func (k Keeper) SetUserPosition(ctx sdk.Context, position types.UserPosition) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.UserPositionPrefix)
	bz := k.cdc.MustMarshal(&position)
	store.Set([]byte(position.Address), bz)
}

// RemoveUserPosition removes a user's DINR position
func (k Keeper) RemoveUserPosition(ctx sdk.Context, address string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.UserPositionPrefix)
	store.Delete([]byte(address))
}

// IterateAllUserPositions iterates over all user positions
func (k Keeper) IterateAllUserPositions(ctx sdk.Context, cb func(position types.UserPosition) bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.UserPositionPrefix)
	
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var position types.UserPosition
		k.cdc.MustUnmarshal(iterator.Value(), &position)
		
		if cb(position) {
			break
		}
	}
}

// GetCollateralAsset returns information about a collateral asset
func (k Keeper) GetCollateralAsset(ctx sdk.Context, denom string) (types.CollateralAsset, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CollateralAssetPrefix)
	
	bz := store.Get([]byte(denom))
	if bz == nil {
		return types.CollateralAsset{}, false
	}
	
	var asset types.CollateralAsset
	k.cdc.MustUnmarshal(bz, &asset)
	return asset, true
}

// SetCollateralAsset sets information about a collateral asset
func (k Keeper) SetCollateralAsset(ctx sdk.Context, asset types.CollateralAsset) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CollateralAssetPrefix)
	bz := k.cdc.MustMarshal(&asset)
	store.Set([]byte(asset.Denom), bz)
}

// GetAllCollateralAssets returns all collateral assets
func (k Keeper) GetAllCollateralAssets(ctx sdk.Context) []types.CollateralAsset {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CollateralAssetPrefix)
	
	var assets []types.CollateralAsset
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var asset types.CollateralAsset
		k.cdc.MustUnmarshal(iterator.Value(), &asset)
		assets = append(assets, asset)
	}
	
	return assets
}

// GetStabilityData returns the current stability data
func (k Keeper) GetStabilityData(ctx sdk.Context) types.StabilityData {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.StabilityDataKey)
	
	if bz == nil {
		// Return default stability data
		return types.StabilityData{
			CurrentPrice:          "1.00",
			TargetPrice:           "1.00",
			PriceDeviation:        0,
			TotalSupply:           sdk.NewCoin(types.DINRDenom, sdk.ZeroInt()),
			TotalCollateralValue:  sdk.NewCoin("inr", sdk.ZeroInt()),
			GlobalCollateralRatio: 0,
			LastUpdate:            ctx.BlockTime(),
		}
	}
	
	var data types.StabilityData
	k.cdc.MustUnmarshal(bz, &data)
	return data
}

// SetStabilityData sets the stability data
func (k Keeper) SetStabilityData(ctx sdk.Context, data types.StabilityData) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&data)
	store.Set(types.StabilityDataKey, bz)
}

// GetInsuranceFund returns the insurance fund data
func (k Keeper) GetInsuranceFund(ctx sdk.Context) types.InsuranceFund {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.InsuranceFundKey)
	
	if bz == nil {
		// Return default insurance fund
		return types.InsuranceFund{
			Balance:       sdk.NewCoin(types.DINRDenom, sdk.ZeroInt()),
			TargetRatio:   types.DefaultInsuranceFundRatio,
			Assets:        sdk.NewCoins(),
			LastRebalance: ctx.BlockTime(),
		}
	}
	
	var fund types.InsuranceFund
	k.cdc.MustUnmarshal(bz, &fund)
	return fund
}

// SetInsuranceFund sets the insurance fund data
func (k Keeper) SetInsuranceFund(ctx sdk.Context, fund types.InsuranceFund) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&fund)
	store.Set(types.InsuranceFundKey, bz)
}

// GetAllYieldStrategies returns all yield strategies
func (k Keeper) GetAllYieldStrategies(ctx sdk.Context) []types.YieldStrategy {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.YieldStrategyPrefix)
	
	var strategies []types.YieldStrategy
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var strategy types.YieldStrategy
		k.cdc.MustUnmarshal(iterator.Value(), &strategy)
		strategies = append(strategies, strategy)
	}
	
	return strategies
}

// SetYieldStrategy sets a yield strategy
func (k Keeper) SetYieldStrategy(ctx sdk.Context, strategy types.YieldStrategy) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.YieldStrategyPrefix)
	bz := k.cdc.MustMarshal(&strategy)
	store.Set([]byte(strategy.Id), bz)
}

// GetLastYieldProcessingTime returns the last time yield was processed
func (k Keeper) GetLastYieldProcessingTime(ctx sdk.Context) time.Time {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.LastYieldProcessingTimeKey)
	
	if bz == nil {
		return ctx.BlockTime().Add(-time.Hour) // Default to 1 hour ago
	}
	
	var timestamp time.Time
	k.cdc.MustUnmarshal(bz, &timestamp)
	return timestamp
}

// SetLastYieldProcessingTime sets the last yield processing time
func (k Keeper) SetLastYieldProcessingTime(ctx sdk.Context, timestamp time.Time) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&timestamp)
	store.Set(types.LastYieldProcessingTimeKey, bz)
}