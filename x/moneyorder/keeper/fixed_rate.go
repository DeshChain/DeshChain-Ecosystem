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
	"fmt"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/moneyorder/types"
)

// CreateFixedRatePool creates a new fixed rate pool
func (k Keeper) CreateFixedRatePool(
	ctx sdk.Context,
	creator sdk.AccAddress,
	token0Denom string,
	token1Denom string,
	exchangeRate sdk.Dec,
	initialLiquidity sdk.Coins,
	description string,
	supportedRegions []string,
) (uint64, error) {
	params := k.GetParams(ctx)
	if !params.EnableFixedRatePools {
		return 0, types.ErrPoolInactive
	}
	
	// Generate new pool ID
	poolId := k.GetNextPoolId(ctx)
	
	// Create the pool
	pool := types.NewFixedRatePool(
		poolId,
		token0Denom,
		token1Denom,
		exchangeRate,
		creator,
	)
	
	// Set additional fields
	pool.Description = description
	pool.SupportedRegions = supportedRegions
	
	// Validate the pool
	if err := pool.ValidatePool(); err != nil {
		return 0, err
	}
	
	// Check initial liquidity
	token0Amount := initialLiquidity.AmountOf(token0Denom)
	token1Amount := initialLiquidity.AmountOf(token1Denom)
	
	if token0Amount.IsZero() || token1Amount.IsZero() {
		return 0, types.ErrInsufficientLiquidity
	}
	
	// Transfer initial liquidity from creator
	if err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx, creator, types.ModuleName, initialLiquidity,
	); err != nil {
		return 0, err
	}
	
	// Update pool balances
	pool.Token0Balance = token0Amount
	pool.Token1Balance = token1Amount
	
	// Store the pool
	k.SetFixedRatePool(ctx, pool)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeFixedRatePoolCreated,
			sdk.NewAttribute(types.AttributeKeyPoolId, fmt.Sprintf("%d", poolId)),
			sdk.NewAttribute(types.AttributeKeyToken0, token0Denom),
			sdk.NewAttribute(types.AttributeKeyToken1, token1Denom),
			sdk.NewAttribute(types.AttributeKeyExchangeRate, exchangeRate.String()),
			sdk.NewAttribute(types.AttributeKeySender, creator.String()),
		),
	)
	
	// Call hooks
	k.AfterPoolCreated(ctx, poolId)
	
	return poolId, nil
}

// SetFixedRatePool stores a fixed rate pool
func (k Keeper) SetFixedRatePool(ctx sdk.Context, pool *types.FixedRatePool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetFixedRatePoolKey(pool.PoolId)
	bz := k.cdc.MustMarshal(pool)
	store.Set(key, bz)
}

// GetFixedRatePool retrieves a fixed rate pool by ID
func (k Keeper) GetFixedRatePool(ctx sdk.Context, poolId uint64) (*types.FixedRatePool, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetFixedRatePoolKey(poolId)
	bz := store.Get(key)
	
	if bz == nil {
		return nil, false
	}
	
	var pool types.FixedRatePool
	k.cdc.MustUnmarshal(bz, &pool)
	return &pool, true
}

// ExecuteFixedRateSwap executes a swap in a fixed rate pool
func (k Keeper) ExecuteFixedRateSwap(
	ctx sdk.Context,
	poolId uint64,
	sender sdk.AccAddress,
	tokenIn sdk.Coin,
	tokenOutDenom string,
	minTokenOut sdk.Int,
) (sdk.Coin, error) {
	// Get the pool
	pool, found := k.GetFixedRatePool(ctx, poolId)
	if !found {
		return sdk.Coin{}, types.ErrPoolNotFound
	}
	
	// Check if pool is active
	if !pool.Active {
		return sdk.Coin{}, types.ErrPoolInactive
	}
	
	if pool.MaintenanceMode {
		return sdk.Coin{}, types.ErrPoolMaintenance
	}
	
	// Determine swap direction
	isForward := tokenIn.Denom == pool.Token0Denom && tokenOutDenom == pool.Token1Denom
	isReverse := tokenIn.Denom == pool.Token1Denom && tokenOutDenom == pool.Token0Denom
	
	if !isForward && !isReverse {
		return sdk.Coin{}, fmt.Errorf("invalid token pair for pool")
	}
	
	// Check daily limit for sender
	if err := k.CheckDailyLimit(ctx, sender, tokenIn.Amount); err != nil {
		return sdk.Coin{}, err
	}
	
	// Check KYC if required
	if pool.RequiresKYC && tokenIn.Amount.GTE(pool.KYCThreshold) {
		if err := k.ValidateKYC(ctx, sender); err != nil {
			return sdk.Coin{}, err
		}
	}
	
	// Calculate output amount and fees
	outputAmount, feeAmount, err := pool.CalculateOutput(tokenIn.Amount, isForward)
	if err != nil {
		return sdk.Coin{}, err
	}
	
	// Check minimum output
	if outputAmount.LT(minTokenOut) {
		return sdk.Coin{}, types.ErrSlippageExceeded
	}
	
	// Create output coin
	tokenOut := sdk.NewCoin(tokenOutDenom, outputAmount)
	
	// Check if pool can fulfill the order
	if !pool.CanFulfillOrder(outputAmount, tokenOutDenom) {
		return sdk.Coin{}, types.ErrInsufficientPoolFunds
	}
	
	// Transfer tokens from sender to pool
	if err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx, sender, types.ModuleName, sdk.NewCoins(tokenIn),
	); err != nil {
		return sdk.Coin{}, err
	}
	
	// Transfer output tokens from pool to sender
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx, types.ModuleName, sender, sdk.NewCoins(tokenOut),
	); err != nil {
		return sdk.Coin{}, err
	}
	
	// Collect fees
	feeCoin := sdk.NewCoin(tokenIn.Denom, feeAmount)
	if err := k.CollectFees(ctx, feeCoin, sender); err != nil {
		return sdk.Coin{}, err
	}
	
	// Update pool balances
	if isForward {
		pool.Token0Balance = pool.Token0Balance.Add(tokenIn.Amount)
		pool.Token1Balance = pool.Token1Balance.Sub(outputAmount)
	} else {
		pool.Token1Balance = pool.Token1Balance.Add(tokenIn.Amount)
		pool.Token0Balance = pool.Token0Balance.Sub(outputAmount)
	}
	
	// Update pool statistics
	pool.UpdateVolume(tokenIn.Amount)
	
	// Save updated pool
	k.SetFixedRatePool(ctx, pool)
	
	// Create money order receipt
	orderId := k.GetNextOrderId(ctx)
	receipt := types.NewMoneyOrderReceipt(
		orderId,
		sender,
		sender, // For swaps, sender is also receiver
		tokenIn,
		fmt.Sprintf("Fixed rate swap: %s to %s", tokenIn.Denom, tokenOutDenom),
	)
	receipt.ExchangeRate = pool.ExchangeRate
	receipt.Fees = feeCoin
	receipt.Status = types.OrderStatusCompleted
	receipt.CompletedAt = ctx.BlockTime()
	
	k.SetMoneyOrderReceipt(ctx, receipt)
	
	// Emit events
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSwap,
			sdk.NewAttribute(types.AttributeKeyPoolId, fmt.Sprintf("%d", poolId)),
			sdk.NewAttribute(types.AttributeKeySender, sender.String()),
			sdk.NewAttribute(types.AttributeKeyTokenIn, tokenIn.String()),
			sdk.NewAttribute(types.AttributeKeyTokenOut, tokenOut.String()),
			sdk.NewAttribute(types.AttributeKeyFees, feeCoin.String()),
		),
	)
	
	// Call hooks
	k.AfterSwap(ctx, poolId, tokenIn, tokenOut)
	
	return tokenOut, nil
}

// UpdateFixedRatePoolParams updates pool parameters (governance only)
func (k Keeper) UpdateFixedRatePoolParams(
	ctx sdk.Context,
	poolId uint64,
	baseFee sdk.Dec,
	active bool,
) error {
	pool, found := k.GetFixedRatePool(ctx, poolId)
	if !found {
		return types.ErrPoolNotFound
	}
	
	// Update parameters
	if !baseFee.IsNil() && !baseFee.IsNegative() {
		pool.BaseFee = baseFee
	}
	
	pool.Active = active
	pool.UpdatedAt = ctx.BlockTime()
	
	// Save updated pool
	k.SetFixedRatePool(ctx, pool)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpdatePool,
			sdk.NewAttribute(types.AttributeKeyPoolId, fmt.Sprintf("%d", poolId)),
			sdk.NewAttribute(types.AttributeKeyStatus, fmt.Sprintf("%v", active)),
		),
	)
	
	return nil
}

// GetAllFixedRatePools returns all fixed rate pools
func (k Keeper) GetAllFixedRatePools(ctx sdk.Context) []*types.FixedRatePool {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixFixedRatePool)
	defer iterator.Close()
	
	var pools []*types.FixedRatePool
	for ; iterator.Valid(); iterator.Next() {
		var pool types.FixedRatePool
		k.cdc.MustUnmarshal(iterator.Value(), &pool)
		pools = append(pools, &pool)
	}
	
	return pools
}

// GetFixedRatePoolsByRegion returns pools supporting a specific region
func (k Keeper) GetFixedRatePoolsByRegion(ctx sdk.Context, postalCode string) []*types.FixedRatePool {
	allPools := k.GetAllFixedRatePools(ctx)
	var regionalPools []*types.FixedRatePool
	
	for _, pool := range allPools {
		if pool.Active {
			// Check if pool supports all regions or specific postal code
			if len(pool.SupportedRegions) == 0 {
				regionalPools = append(regionalPools, pool)
			} else {
				for _, region := range pool.SupportedRegions {
					if region == postalCode || region == postalCode[:3] { // Support partial postal codes
						regionalPools = append(regionalPools, pool)
						break
					}
				}
			}
		}
	}
	
	return regionalPools
}