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
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/deshchain/deshchain/x/moneyorder/types"
)

// AMMPool represents a basic AMM pool interface
type AMMPool interface {
	GetId() uint64
	GetTokens() []string
	GetBalance(denom string) sdk.Int
	SwapExactAmountIn(ctx sdk.Context, tokenIn sdk.Coin, tokenOutDenom string, minAmountOut sdk.Int) (tokenOut sdk.Coin, err error)
	SwapExactAmountOut(ctx sdk.Context, tokenInDenom string, maxAmountIn sdk.Int, tokenOut sdk.Coin) (tokenIn sdk.Coin, err error)
	JoinPool(ctx sdk.Context, tokensIn sdk.Coins, shareOutAmount sdk.Int) (shareOut sdk.Int, err error)
	ExitPool(ctx sdk.Context, shareInAmount sdk.Int, minAmountsOut sdk.Coins) (tokensOut sdk.Coins, err error)
}

// CreateAMMPool creates a new AMM pool with constant product formula
func (k Keeper) CreateAMMPool(
	ctx sdk.Context,
	creator sdk.AccAddress,
	poolAssets []types.PoolAsset,
	swapFee sdk.Dec,
) (uint64, error) {
	// Validate pool assets
	if len(poolAssets) != 2 {
		return 0, sdkerrors.Wrap(types.ErrInvalidPoolAssets, "AMM pools must have exactly 2 assets")
	}

	// Validate swap fee
	params := k.GetParams(ctx)
	if swapFee.LT(sdk.ZeroDec()) || swapFee.GT(params.TradingFeeRate.Mul(sdk.NewDec(2))) {
		return 0, sdkerrors.Wrap(types.ErrInvalidSwapFee, "swap fee out of allowed range")
	}

	// Get next pool ID
	poolId := k.GetNextPoolId(ctx)
	
	// Calculate initial pool shares based on geometric mean
	shareAmount := sdk.ZeroInt()
	if len(poolAssets) == 2 {
		// shareAmount = sqrt(amount0 * amount1)
		product := poolAssets[0].Amount.Mul(poolAssets[1].Amount)
		shareAmount = product.ApproxSqrt()
	}

	// Create pool
	pool := types.AMMPoolInfo{
		PoolId:          poolId,
		PoolAssets:      poolAssets,
		TotalShares:     shareAmount,
		SwapFee:         swapFee,
		Creator:         creator.String(),
		Active:          true,
		CulturalPair:    k.isCulturalPair(poolAssets[0].Denom, poolAssets[1].Denom),
		FestivalBonus:   params.EnableFestivalBonuses,
		VillagePriority: false,
	}

	// Save pool
	k.SetAMMPool(ctx, pool)
	k.SetNextPoolId(ctx, poolId+1)

	// Transfer assets from creator to module
	for _, asset := range poolAssets {
		coin := sdk.NewCoin(asset.Denom, asset.Amount)
		if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, creator, types.ModuleName, sdk.NewCoins(coin)); err != nil {
			return 0, err
		}
	}

	// Mint and send pool shares to creator
	poolShareCoin := sdk.NewCoin(fmt.Sprintf("pool%d", poolId), shareAmount)
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(poolShareCoin)); err != nil {
		return 0, err
	}
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creator, sdk.NewCoins(poolShareCoin)); err != nil {
		return 0, err
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreatePool,
			sdk.NewAttribute(types.AttributeKeyPoolId, fmt.Sprintf("%d", poolId)),
			sdk.NewAttribute(types.AttributeKeyCreator, creator.String()),
			sdk.NewAttribute(types.AttributeKeyPoolType, "amm"),
		),
	)

	return poolId, nil
}

// SwapExactAmountIn performs a swap with exact input amount
func (k Keeper) SwapExactAmountIn(
	ctx sdk.Context,
	sender sdk.AccAddress,
	poolId uint64,
	tokenIn sdk.Coin,
	tokenOutDenom string,
	minAmountOut sdk.Int,
) (tokenOut sdk.Coin, err error) {
	pool, found := k.GetAMMPool(ctx, poolId)
	if !found {
		return tokenOut, types.ErrPoolNotFound
	}

	if !pool.Active {
		return tokenOut, types.ErrPoolNotActive
	}

	// Find token indices
	tokenInIndex := -1
	tokenOutIndex := -1
	for i, asset := range pool.PoolAssets {
		if asset.Denom == tokenIn.Denom {
			tokenInIndex = i
		}
		if asset.Denom == tokenOutDenom {
			tokenOutIndex = i
		}
	}

	if tokenInIndex == -1 || tokenOutIndex == -1 {
		return tokenOut, types.ErrInvalidTokenPair
	}

	// Calculate output amount using constant product formula
	// (x + dx) * (y - dy) = x * y
	// dy = y * dx / (x + dx)
	
	reserveIn := pool.PoolAssets[tokenInIndex].Amount
	reserveOut := pool.PoolAssets[tokenOutIndex].Amount
	
	// Apply swap fee
	params := k.GetParams(ctx)
	effectiveFee := pool.SwapFee
	
	// Apply cultural discount if applicable
	if pool.CulturalPair && k.IsInFestivalPeriod(ctx) && params.EnableFestivalBonuses {
		effectiveFee = effectiveFee.Mul(sdk.OneDec().Sub(params.FestivalDiscount))
	}
	
	amountInWithFee := tokenIn.Amount.ToDec().Mul(sdk.OneDec().Sub(effectiveFee))
	
	// Calculate output amount
	numerator := amountInWithFee.Mul(reserveOut.ToDec())
	denominator := reserveIn.ToDec().Add(amountInWithFee)
	amountOut := numerator.Quo(denominator).TruncateInt()
	
	if amountOut.LT(minAmountOut) {
		return tokenOut, sdkerrors.Wrap(types.ErrInsufficientOutput, "output amount less than minimum")
	}
	
	tokenOut = sdk.NewCoin(tokenOutDenom, amountOut)
	
	// Update pool reserves
	pool.PoolAssets[tokenInIndex].Amount = pool.PoolAssets[tokenInIndex].Amount.Add(tokenIn.Amount)
	pool.PoolAssets[tokenOutIndex].Amount = pool.PoolAssets[tokenOutIndex].Amount.Sub(amountOut)
	
	// Transfer tokens
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(tokenIn)); err != nil {
		return tokenOut, err
	}
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, sdk.NewCoins(tokenOut)); err != nil {
		return tokenOut, err
	}
	
	// Collect and distribute trading fees
	feeAmount := tokenIn.Amount.ToDec().Mul(effectiveFee).TruncateInt()
	if feeAmount.GT(sdk.ZeroInt()) && k.revenueKeeper != nil {
		feeCoin := sdk.NewCoin(tokenIn.Denom, feeAmount)
		feeCoins := sdk.NewCoins(feeCoin)
		
		// Collect trading fee using revenue keeper
		pair := tokenIn.Denom + "-" + tokenOutDenom
		if err := k.revenueKeeper.CollectTradingFee(ctx, types.ModuleName, sender, feeCoins, pair); err != nil {
			k.Logger(ctx).Error("Failed to collect trading fee", "error", err, "sender", sender, "fee", feeAmount)
			// Continue even if fee collection fails
		}
	}
	
	// Update pool
	k.SetAMMPool(ctx, pool)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSwap,
			sdk.NewAttribute(types.AttributeKeyPoolId, fmt.Sprintf("%d", poolId)),
			sdk.NewAttribute(types.AttributeKeySender, sender.String()),
			sdk.NewAttribute(types.AttributeKeyTokenIn, tokenIn.String()),
			sdk.NewAttribute(types.AttributeKeyTokenOut, tokenOut.String()),
		),
	)
	
	return tokenOut, nil
}

// JoinPool adds liquidity to an AMM pool
func (k Keeper) JoinPool(
	ctx sdk.Context,
	sender sdk.AccAddress,
	poolId uint64,
	tokensIn sdk.Coins,
	shareOutAmount sdk.Int,
) (shareOut sdk.Int, err error) {
	pool, found := k.GetAMMPool(ctx, poolId)
	if !found {
		return sdk.ZeroInt(), types.ErrPoolNotFound
	}

	if !pool.Active {
		return sdk.ZeroInt(), types.ErrPoolNotActive
	}

	// Calculate share amount based on token ratios
	if pool.TotalShares.IsZero() {
		// Initial liquidity
		if len(tokensIn) != 2 {
			return sdk.ZeroInt(), types.ErrInvalidPoolAssets
		}
		shareOut = tokensIn[0].Amount.Mul(tokensIn[1].Amount).ApproxSqrt()
	} else {
		// Calculate proportional shares
		ratio0 := tokensIn[0].Amount.ToDec().Quo(pool.PoolAssets[0].Amount.ToDec())
		ratio1 := tokensIn[1].Amount.ToDec().Quo(pool.PoolAssets[1].Amount.ToDec())
		
		if !ratio0.Equal(ratio1) {
			return sdk.ZeroInt(), sdkerrors.Wrap(types.ErrInvalidPoolAssets, "token ratios must match pool ratios")
		}
		
		shareOut = pool.TotalShares.ToDec().Mul(ratio0).TruncateInt()
	}
	
	if shareOutAmount.GT(sdk.ZeroInt()) && shareOut.LT(shareOutAmount) {
		return sdk.ZeroInt(), sdkerrors.Wrap(types.ErrInsufficientShares, "calculated shares less than requested")
	}
	
	// Transfer tokens to pool
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, tokensIn); err != nil {
		return sdk.ZeroInt(), err
	}
	
	// Update pool reserves
	for i, tokenIn := range tokensIn {
		pool.PoolAssets[i].Amount = pool.PoolAssets[i].Amount.Add(tokenIn.Amount)
	}
	pool.TotalShares = pool.TotalShares.Add(shareOut)
	
	// Mint and send pool shares
	poolShareCoin := sdk.NewCoin(fmt.Sprintf("pool%d", poolId), shareOut)
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(poolShareCoin)); err != nil {
		return sdk.ZeroInt(), err
	}
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, sdk.NewCoins(poolShareCoin)); err != nil {
		return sdk.ZeroInt(), err
	}
	
	// Update pool
	k.SetAMMPool(ctx, pool)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeJoinPool,
			sdk.NewAttribute(types.AttributeKeyPoolId, fmt.Sprintf("%d", poolId)),
			sdk.NewAttribute(types.AttributeKeySender, sender.String()),
			sdk.NewAttribute(types.AttributeKeySharesOut, shareOut.String()),
		),
	)
	
	return shareOut, nil
}

// ExitPool removes liquidity from an AMM pool
func (k Keeper) ExitPool(
	ctx sdk.Context,
	sender sdk.AccAddress,
	poolId uint64,
	shareInAmount sdk.Int,
	minAmountsOut sdk.Coins,
) (tokensOut sdk.Coins, err error) {
	pool, found := k.GetAMMPool(ctx, poolId)
	if !found {
		return nil, types.ErrPoolNotFound
	}

	// Calculate proportional token amounts
	shareRatio := shareInAmount.ToDec().Quo(pool.TotalShares.ToDec())
	
	tokensOut = sdk.NewCoins()
	for _, asset := range pool.PoolAssets {
		amountOut := asset.Amount.ToDec().Mul(shareRatio).TruncateInt()
		tokenOut := sdk.NewCoin(asset.Denom, amountOut)
		tokensOut = tokensOut.Add(tokenOut)
		
		// Check minimum amounts
		minAmount := minAmountsOut.AmountOf(asset.Denom)
		if amountOut.LT(minAmount) {
			return nil, sdkerrors.Wrap(types.ErrInsufficientOutput, "output amount less than minimum")
		}
	}
	
	// Burn pool shares
	poolShareCoin := sdk.NewCoin(fmt.Sprintf("pool%d", poolId), shareInAmount)
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(poolShareCoin)); err != nil {
		return nil, err
	}
	if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(poolShareCoin)); err != nil {
		return nil, err
	}
	
	// Update pool reserves
	for i, tokenOut := range tokensOut {
		pool.PoolAssets[i].Amount = pool.PoolAssets[i].Amount.Sub(tokenOut.Amount)
	}
	pool.TotalShares = pool.TotalShares.Sub(shareInAmount)
	
	// Transfer tokens to sender
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, tokensOut); err != nil {
		return nil, err
	}
	
	// Update pool
	k.SetAMMPool(ctx, pool)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeExitPool,
			sdk.NewAttribute(types.AttributeKeyPoolId, fmt.Sprintf("%d", poolId)),
			sdk.NewAttribute(types.AttributeKeySender, sender.String()),
			sdk.NewAttribute(types.AttributeKeySharesIn, shareInAmount.String()),
		),
	)
	
	return tokensOut, nil
}

// Helper functions

func (k Keeper) isCulturalPair(denom0, denom1 string) bool {
	culturalDenoms := map[string]bool{
		"unamo":     true,
		"uinr":      true,
		"ucultural": true,
		"uheritage": true,
	}
	
	return culturalDenoms[denom0] || culturalDenoms[denom1]
}

func (k Keeper) GetAMMPool(ctx sdk.Context, poolId uint64) (types.AMMPoolInfo, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetAMMPoolKey(poolId)
	
	bz := store.Get(key)
	if bz == nil {
		return types.AMMPoolInfo{}, false
	}
	
	var pool types.AMMPoolInfo
	k.cdc.MustUnmarshal(bz, &pool)
	return pool, true
}

func (k Keeper) SetAMMPool(ctx sdk.Context, pool types.AMMPoolInfo) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetAMMPoolKey(pool.PoolId)
	bz := k.cdc.MustMarshal(&pool)
	store.Set(key, bz)
}

func (k Keeper) GetNextPoolId(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyNextPoolId)
	if bz == nil {
		return 1
	}
	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) SetNextPoolId(ctx sdk.Context, poolId uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyNextPoolId, sdk.Uint64ToBigEndian(poolId))
}