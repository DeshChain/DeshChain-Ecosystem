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
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/moneyorder/types"
)

// ConcentratedLiquidityPool represents a concentrated liquidity pool for efficient trading
type ConcentratedLiquidityPool struct {
	PoolId           uint64         `json:"pool_id"`
	TokenA           string         `json:"token_a"`
	TokenB           string         `json:"token_b"`
	TickSpacing      uint64         `json:"tick_spacing"`      // Minimum tick spacing
	Fee              sdk.Dec        `json:"fee"`               // Pool fee percentage
	CurrentTick      int64          `json:"current_tick"`      // Current price tick
	CurrentSqrtPrice sdk.Dec        `json:"current_sqrt_price"` // Square root of current price
	Liquidity        sdk.Int        `json:"liquidity"`         // Current active liquidity
	Creator          sdk.AccAddress `json:"creator"`
	CreatedAt        int64          `json:"created_at"`
	Status           string         `json:"status"`
	
	// Village-specific features
	VillageRestricted bool           `json:"village_restricted"`
	AllowedVillages   []string       `json:"allowed_villages,omitempty"`
	CulturalBonus     sdk.Dec        `json:"cultural_bonus"` // Festival bonus
	
	// Unified pool integration
	UnifiedPoolId     uint64         `json:"unified_pool_id,omitempty"`
	LiquiditySource   string         `json:"liquidity_source"` // "unified", "dedicated", "mixed"
}

// Position represents a liquidity position in a concentrated liquidity pool
type Position struct {
	PositionId       uint64         `json:"position_id"`
	PoolId           uint64         `json:"pool_id"`
	Owner            sdk.AccAddress `json:"owner"`
	LowerTick        int64          `json:"lower_tick"`
	UpperTick        int64          `json:"upper_tick"`
	Liquidity        sdk.Int        `json:"liquidity"`
	TokenAAmount     sdk.Coin       `json:"token_a_amount"`
	TokenBAmount     sdk.Coin       `json:"token_b_amount"`
	FeesAccumulated  sdk.Coins      `json:"fees_accumulated"`
	CreatedAt        int64          `json:"created_at"`
	LastUpdated      int64          `json:"last_updated"`
	Status           string         `json:"status"`
	
	// Cultural features
	CulturalQuote    string         `json:"cultural_quote,omitempty"`
	PatriotismBonus  sdk.Coin       `json:"patriotism_bonus,omitempty"`
}

// Tick represents a price tick in the concentrated liquidity pool
type Tick struct {
	PoolId              uint64  `json:"pool_id"`
	TickIndex           int64   `json:"tick_index"`
	LiquidityGross      sdk.Int `json:"liquidity_gross"`      // Total liquidity referencing this tick
	LiquidityNet        sdk.Int `json:"liquidity_net"`        // Net liquidity change at this tick
	FeeGrowthOutsideA   sdk.Dec `json:"fee_growth_outside_a"` // Fee growth per unit liquidity on token A
	FeeGrowthOutsideB   sdk.Dec `json:"fee_growth_outside_b"` // Fee growth per unit liquidity on token B
	Initialized         bool    `json:"initialized"`
}

// SwapResult represents the result of a concentrated liquidity swap
type SwapResult struct {
	AmountIn      sdk.Coin `json:"amount_in"`
	AmountOut     sdk.Coin `json:"amount_out"`
	Fee           sdk.Coin `json:"fee"`
	NewTick       int64    `json:"new_tick"`
	NewSqrtPrice  sdk.Dec  `json:"new_sqrt_price"`
	PriceImpact   sdk.Dec  `json:"price_impact"`
}

// CreateConcentratedLiquidityPool creates a new concentrated liquidity pool
func (k Keeper) CreateConcentratedLiquidityPool(
	ctx sdk.Context,
	creator sdk.AccAddress,
	tokenA, tokenB string,
	tickSpacing uint64,
	fee sdk.Dec,
	initialPrice sdk.Dec,
) (*ConcentratedLiquidityPool, error) {
	// Validate inputs
	if tokenA == tokenB {
		return nil, types.ErrInvalidTokenPair
	}
	
	if fee.IsNegative() || fee.GTE(sdk.OneDec()) {
		return nil, types.ErrInvalidFee
	}
	
	if initialPrice.IsNegative() || initialPrice.IsZero() {
		return nil, types.ErrInvalidPrice
	}
	
	// Generate pool ID
	poolId := k.GetNextPoolId(ctx)
	
	// Calculate initial tick and sqrt price
	initialTick := k.priceToTick(initialPrice)
	initialSqrtPrice := k.tickToSqrtPrice(initialTick)
	
	pool := &ConcentratedLiquidityPool{
		PoolId:           poolId,
		TokenA:           tokenA,
		TokenB:           tokenB,
		TickSpacing:      tickSpacing,
		Fee:              fee,
		CurrentTick:      initialTick,
		CurrentSqrtPrice: initialSqrtPrice,
		Liquidity:        sdk.ZeroInt(),
		Creator:          creator,
		CreatedAt:        ctx.BlockTime().Unix(),
		Status:           "active",
		CulturalBonus:    sdk.ZeroDec(),
		LiquiditySource:  "dedicated",
	}
	
	// Store the pool
	k.SetConcentratedLiquidityPool(ctx, *pool)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"concentrated_pool_created",
			sdk.NewAttribute("pool_id", fmt.Sprintf("%d", poolId)),
			sdk.NewAttribute("token_a", tokenA),
			sdk.NewAttribute("token_b", tokenB),
			sdk.NewAttribute("initial_price", initialPrice.String()),
			sdk.NewAttribute("creator", creator.String()),
		),
	)
	
	return pool, nil
}

// AddLiquidity adds liquidity to a concentrated liquidity pool within a price range
func (k Keeper) AddLiquidity(
	ctx sdk.Context,
	poolId uint64,
	provider sdk.AccAddress,
	lowerTick, upperTick int64,
	amountA, amountB sdk.Coin,
) (*Position, error) {
	// Get pool
	pool, found := k.GetConcentratedLiquidityPool(ctx, poolId)
	if !found {
		return nil, types.ErrPoolNotFound
	}
	
	// Validate tick range
	if lowerTick >= upperTick {
		return nil, types.ErrInvalidTickRange
	}
	
	// Calculate liquidity to be provided
	liquidity := k.calculateLiquidityFromAmounts(pool.CurrentSqrtPrice, lowerTick, upperTick, amountA, amountB)
	
	if liquidity.IsZero() {
		return nil, types.ErrInsufficientLiquidity
	}
	
	// Transfer tokens from provider to pool
	if err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx, provider, types.ModuleName, sdk.NewCoins(amountA, amountB),
	); err != nil {
		return nil, err
	}
	
	// Create position
	positionId := k.GetNextPositionId(ctx)
	position := &Position{
		PositionId:      positionId,
		PoolId:          poolId,
		Owner:           provider,
		LowerTick:       lowerTick,
		UpperTick:       upperTick,
		Liquidity:       liquidity,
		TokenAAmount:    amountA,
		TokenBAmount:    amountB,
		FeesAccumulated: sdk.NewCoins(),
		CreatedAt:       ctx.BlockTime().Unix(),
		LastUpdated:     ctx.BlockTime().Unix(),
		Status:          "active",
	}
	
	// Add cultural quote if enabled
	if pool.CulturalBonus.IsPositive() {
		position.CulturalQuote = k.getRandomCulturalQuote(ctx, "liquidity_provision")
	}
	
	// Update pool liquidity if position is in range
	if k.isPositionInRange(pool.CurrentTick, lowerTick, upperTick) {
		pool.Liquidity = pool.Liquidity.Add(liquidity)
	}
	
	// Update ticks
	k.updateTick(ctx, poolId, lowerTick, liquidity, true)
	k.updateTick(ctx, poolId, upperTick, liquidity, false)
	
	// Store updates
	k.SetConcentratedLiquidityPool(ctx, pool)
	k.SetPosition(ctx, *position)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"liquidity_added",
			sdk.NewAttribute("pool_id", fmt.Sprintf("%d", poolId)),
			sdk.NewAttribute("position_id", fmt.Sprintf("%d", positionId)),
			sdk.NewAttribute("provider", provider.String()),
			sdk.NewAttribute("liquidity", liquidity.String()),
			sdk.NewAttribute("lower_tick", fmt.Sprintf("%d", lowerTick)),
			sdk.NewAttribute("upper_tick", fmt.Sprintf("%d", upperTick)),
		),
	)
	
	return position, nil
}

// Swap executes a swap in a concentrated liquidity pool
func (k Keeper) SwapConcentratedLiquidity(
	ctx sdk.Context,
	poolId uint64,
	trader sdk.AccAddress,
	amountIn sdk.Coin,
	minAmountOut sdk.Coin,
	zeroForOne bool, // true if swapping token0 for token1
) (*SwapResult, error) {
	// Get pool
	pool, found := k.GetConcentratedLiquidityPool(ctx, poolId)
	if !found {
		return nil, types.ErrPoolNotFound
	}
	
	// Validate input
	if !amountIn.IsPositive() {
		return nil, types.ErrInvalidAmount
	}
	
	// Calculate swap
	result, err := k.calculateSwap(ctx, pool, amountIn, zeroForOne)
	if err != nil {
		return nil, err
	}
	
	// Check slippage protection
	if result.AmountOut.Amount.LT(minAmountOut.Amount) {
		return nil, types.ErrExcessiveSlippage
	}
	
	// Transfer input token from trader
	if err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx, trader, types.ModuleName, sdk.NewCoins(amountIn),
	); err != nil {
		return nil, err
	}
	
	// Transfer output token to trader
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx, types.ModuleName, trader, sdk.NewCoins(result.AmountOut),
	); err != nil {
		return nil, err
	}
	
	// Update pool state
	pool.CurrentTick = result.NewTick
	pool.CurrentSqrtPrice = result.NewSqrtPrice
	
	// Apply cultural bonus if applicable
	if pool.CulturalBonus.IsPositive() {
		bonus := k.calculateCulturalBonus(ctx, result.AmountOut, pool.CulturalBonus)
		if bonus.IsPositive() {
			bonusCoin := sdk.NewCoin(result.AmountOut.Denom, bonus)
			if err := k.bankKeeper.SendCoinsFromModuleToAccount(
				ctx, types.ModuleName, trader, sdk.NewCoins(bonusCoin),
			); err == nil {
				result.AmountOut = result.AmountOut.Add(bonusCoin)
			}
		}
	}
	
	// Store updated pool
	k.SetConcentratedLiquidityPool(ctx, pool)
	
	// Update fees for positions
	k.updatePositionFees(ctx, poolId, result.Fee)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"concentrated_swap",
			sdk.NewAttribute("pool_id", fmt.Sprintf("%d", poolId)),
			sdk.NewAttribute("trader", trader.String()),
			sdk.NewAttribute("amount_in", result.AmountIn.String()),
			sdk.NewAttribute("amount_out", result.AmountOut.String()),
			sdk.NewAttribute("fee", result.Fee.String()),
			sdk.NewAttribute("new_tick", fmt.Sprintf("%d", result.NewTick)),
		),
	)
	
	return result, nil
}

// RemoveLiquidity removes liquidity from a position
func (k Keeper) RemoveLiquidity(
	ctx sdk.Context,
	positionId uint64,
	liquidityToRemove sdk.Int,
) (sdk.Coin, sdk.Coin, sdk.Coins, error) {
	// Get position
	position, found := k.GetPosition(ctx, positionId)
	if !found {
		return sdk.Coin{}, sdk.Coin{}, sdk.Coins{}, types.ErrPositionNotFound
	}
	
	// Get pool
	pool, found := k.GetConcentratedLiquidityPool(ctx, position.PoolId)
	if !found {
		return sdk.Coin{}, sdk.Coin{}, sdk.Coins{}, types.ErrPoolNotFound
	}
	
	// Validate removal amount
	if liquidityToRemove.GT(position.Liquidity) {
		return sdk.Coin{}, sdk.Coin{}, sdk.Coins{}, types.ErrInsufficientLiquidity
	}
	
	// Calculate token amounts to return
	amountA, amountB := k.calculateAmountsFromLiquidity(
		pool.CurrentSqrtPrice, position.LowerTick, position.UpperTick, liquidityToRemove,
	)
	
	// Calculate fees earned
	feesEarned := k.calculateFeesEarned(ctx, position)
	
	// Update position
	position.Liquidity = position.Liquidity.Sub(liquidityToRemove)
	position.TokenAAmount = position.TokenAAmount.Sub(sdk.NewCoin(position.TokenAAmount.Denom, amountA))
	position.TokenBAmount = position.TokenBAmount.Sub(sdk.NewCoin(position.TokenBAmount.Denom, amountB))
	position.LastUpdated = ctx.BlockTime().Unix()
	
	// If no liquidity left, mark position as closed
	if position.Liquidity.IsZero() {
		position.Status = "closed"
	}
	
	// Update pool liquidity if position is in range
	if k.isPositionInRange(pool.CurrentTick, position.LowerTick, position.UpperTick) {
		pool.Liquidity = pool.Liquidity.Sub(liquidityToRemove)
	}
	
	// Update ticks
	k.updateTick(ctx, position.PoolId, position.LowerTick, liquidityToRemove, false)
	k.updateTick(ctx, position.PoolId, position.UpperTick, liquidityToRemove, true)
	
	// Transfer tokens back to owner
	tokensToReturn := sdk.NewCoins(
		sdk.NewCoin(pool.TokenA, amountA),
		sdk.NewCoin(pool.TokenB, amountB),
	).Add(feesEarned...)
	
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx, types.ModuleName, position.Owner, tokensToReturn,
	); err != nil {
		return sdk.Coin{}, sdk.Coin{}, sdk.Coins{}, err
	}
	
	// Store updates
	k.SetConcentratedLiquidityPool(ctx, pool)
	k.SetPosition(ctx, position)
	
	return sdk.NewCoin(pool.TokenA, amountA), sdk.NewCoin(pool.TokenB, amountB), feesEarned, nil
}

// Helper functions

func (k Keeper) priceToTick(price sdk.Dec) int64 {
	// Simplified price to tick conversion
	// In production, this would use proper logarithmic calculation
	if price.Equal(sdk.OneDec()) {
		return 0
	}
	// Placeholder implementation
	return int64(price.MulInt64(1000).TruncateInt64()) - 1000
}

func (k Keeper) tickToSqrtPrice(tick int64) sdk.Dec {
	// Simplified tick to sqrt price conversion
	// In production, this would use proper exponential calculation
	adjustedTick := sdk.NewDec(tick + 1000).QuoInt64(1000)
	return adjustedTick.ApproxSqrt()
}

func (k Keeper) calculateLiquidityFromAmounts(
	sqrtPrice sdk.Dec, 
	lowerTick, upperTick int64, 
	amountA, amountB sdk.Coin,
) sdk.Int {
	// Simplified liquidity calculation
	// In production, this would use proper concentrated liquidity math
	avgAmount := amountA.Amount.Add(amountB.Amount).QuoRaw(2)
	return avgAmount
}

func (k Keeper) isPositionInRange(currentTick, lowerTick, upperTick int64) bool {
	return currentTick >= lowerTick && currentTick < upperTick
}

func (k Keeper) calculateSwap(
	ctx sdk.Context,
	pool ConcentratedLiquidityPool,
	amountIn sdk.Coin,
	zeroForOne bool,
) (*SwapResult, error) {
	// Simplified swap calculation
	// In production, this would iterate through ticks and calculate precise amounts
	
	// Apply fee
	feeAmount := amountIn.Amount.ToDec().Mul(pool.Fee).TruncateInt()
	amountAfterFee := amountIn.Amount.Sub(feeAmount)
	
	// Simple 1:1 swap for demonstration (would be more complex in production)
	var amountOut sdk.Coin
	var newTick int64
	var newSqrtPrice sdk.Dec
	
	if zeroForOne {
		amountOut = sdk.NewCoin(pool.TokenB, amountAfterFee)
		newTick = pool.CurrentTick - 1
	} else {
		amountOut = sdk.NewCoin(pool.TokenA, amountAfterFee)
		newTick = pool.CurrentTick + 1
	}
	
	newSqrtPrice = k.tickToSqrtPrice(newTick)
	
	return &SwapResult{
		AmountIn:     amountIn,
		AmountOut:    amountOut,
		Fee:          sdk.NewCoin(amountIn.Denom, feeAmount),
		NewTick:      newTick,
		NewSqrtPrice: newSqrtPrice,
		PriceImpact:  sdk.NewDecWithPrec(1, 2), // 1% placeholder
	}, nil
}

func (k Keeper) calculateAmountsFromLiquidity(
	sqrtPrice sdk.Dec,
	lowerTick, upperTick int64,
	liquidity sdk.Int,
) (sdk.Int, sdk.Int) {
	// Simplified amount calculation
	// In production, this would use proper concentrated liquidity formulas
	halfLiquidity := liquidity.QuoRaw(2)
	return halfLiquidity, halfLiquidity
}

func (k Keeper) calculateFeesEarned(ctx sdk.Context, position Position) sdk.Coins {
	// Simplified fee calculation
	// In production, this would calculate fees based on position's share of trading volume
	return position.FeesAccumulated
}

func (k Keeper) updateTick(ctx sdk.Context, poolId uint64, tickIndex int64, liquidityDelta sdk.Int, upper bool) {
	tick, found := k.GetTick(ctx, poolId, tickIndex)
	if !found {
		tick = Tick{
			PoolId:              poolId,
			TickIndex:           tickIndex,
			LiquidityGross:      sdk.ZeroInt(),
			LiquidityNet:        sdk.ZeroInt(),
			FeeGrowthOutsideA:   sdk.ZeroDec(),
			FeeGrowthOutsideB:   sdk.ZeroDec(),
			Initialized:         true,
		}
	}
	
	tick.LiquidityGross = tick.LiquidityGross.Add(liquidityDelta)
	
	if upper {
		tick.LiquidityNet = tick.LiquidityNet.Sub(liquidityDelta)
	} else {
		tick.LiquidityNet = tick.LiquidityNet.Add(liquidityDelta)
	}
	
	k.SetTick(ctx, tick)
}

func (k Keeper) updatePositionFees(ctx sdk.Context, poolId uint64, fee sdk.Coin) {
	// Update fee tracking for all positions in the pool
	// This is a simplified implementation
}

func (k Keeper) calculateCulturalBonus(ctx sdk.Context, amount sdk.Coin, bonusRate sdk.Dec) sdk.Int {
	return amount.Amount.ToDec().Mul(bonusRate).TruncateInt()
}

func (k Keeper) getRandomCulturalQuote(ctx sdk.Context, category string) string {
	// Return a random cultural quote based on category
	quotes := map[string][]string{
		"liquidity_provision": {
			"जल ही जीवन है - Water is life, liquidity is prosperity",
			"सहयोग से समृद्धि - Prosperity through cooperation",
			"एकता में शक्ति - Strength in unity, wealth in sharing",
		},
	}
	
	if categoryQuotes, exists := quotes[category]; exists && len(categoryQuotes) > 0 {
		// Simple selection - in production would use proper randomization
		return categoryQuotes[0]
	}
	
	return "धन्यवाद - Thank you for contributing to our community"
}

// Storage functions
func (k Keeper) SetConcentratedLiquidityPool(ctx sdk.Context, pool ConcentratedLiquidityPool) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&pool)
	store.Set(types.KeyConcentratedPool(pool.PoolId), bz)
}

func (k Keeper) GetConcentratedLiquidityPool(ctx sdk.Context, poolId uint64) (ConcentratedLiquidityPool, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyConcentratedPool(poolId))
	if bz == nil {
		return ConcentratedLiquidityPool{}, false
	}
	
	var pool ConcentratedLiquidityPool
	k.cdc.MustUnmarshal(bz, &pool)
	return pool, true
}

func (k Keeper) SetPosition(ctx sdk.Context, position Position) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&position)
	store.Set(types.KeyPosition(position.PositionId), bz)
}

func (k Keeper) GetPosition(ctx sdk.Context, positionId uint64) (Position, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyPosition(positionId))
	if bz == nil {
		return Position{}, false
	}
	
	var position Position
	k.cdc.MustUnmarshal(bz, &position)
	return position, true
}

func (k Keeper) SetTick(ctx sdk.Context, tick Tick) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&tick)
	store.Set(types.KeyTick(tick.PoolId, tick.TickIndex), bz)
}

func (k Keeper) GetTick(ctx sdk.Context, poolId uint64, tickIndex int64) (Tick, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyTick(poolId, tickIndex))
	if bz == nil {
		return Tick{}, false
	}
	
	var tick Tick
	k.cdc.MustUnmarshal(bz, &tick)
	return tick, true
}

func (k Keeper) GetNextPositionId(ctx sdk.Context) uint64 {
	// Simplified counter - in production would use proper sequence
	return uint64(ctx.BlockTime().Unix())
}