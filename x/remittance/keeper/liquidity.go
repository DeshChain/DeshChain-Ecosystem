package keeper

import (
	"context"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/remittance/types"
)

// ========================= Liquidity Pool Management =========================

// SetLiquidityPool stores a liquidity pool
func (k Keeper) SetLiquidityPool(ctx context.Context, pool types.LiquidityPool) error {
	store := k.GetStore(ctx)
	key := types.LiquidityPoolKey(pool.Id)
	bz := k.cdc.MustMarshal(&pool)
	store.Set(key, bz)
	return nil
}

// GetLiquidityPool retrieves a liquidity pool by ID
func (k Keeper) GetLiquidityPool(ctx context.Context, poolID string) (types.LiquidityPool, error) {
	store := k.GetStore(ctx)
	key := types.LiquidityPoolKey(poolID)
	bz := store.Get(key)
	if bz == nil {
		return types.LiquidityPool{}, types.ErrPoolNotFound
	}

	var pool types.LiquidityPool
	k.cdc.MustUnmarshal(bz, &pool)
	return pool, nil
}

// HasLiquidityPool checks if a pool exists
func (k Keeper) HasLiquidityPool(ctx context.Context, poolID string) bool {
	store := k.GetStore(ctx)
	key := types.LiquidityPoolKey(poolID)
	return store.Has(key)
}

// DeleteLiquidityPool removes a liquidity pool
func (k Keeper) DeleteLiquidityPool(ctx context.Context, poolID string) error {
	store := k.GetStore(ctx)
	key := types.LiquidityPoolKey(poolID)
	store.Delete(key)
	return nil
}

// GetAllLiquidityPools returns all liquidity pools
func (k Keeper) GetAllLiquidityPools(ctx context.Context) ([]types.LiquidityPool, error) {
	store := k.GetStore(ctx)
	iterator := store.Iterator(types.LiquidityPoolKeyPrefix, nil)
	defer iterator.Close()

	var pools []types.LiquidityPool
	for ; iterator.Valid(); iterator.Next() {
		var pool types.LiquidityPool
		k.cdc.MustUnmarshal(iterator.Value(), &pool)
		pools = append(pools, pool)
	}

	return pools, nil
}

// CreateLiquidityPool creates a new liquidity pool
func (k Keeper) CreateLiquidityPool(
	ctx context.Context,
	baseCurrency, quoteCurrency string,
	baseAmount, quoteAmount sdk.Int,
) (string, error) {
	// Validate currencies
	if baseCurrency == quoteCurrency {
		return "", fmt.Errorf("base and quote currencies cannot be the same")
	}

	// Check if pool already exists
	poolID := k.generatePoolID(baseCurrency, quoteCurrency)
	if k.HasLiquidityPool(ctx, poolID) {
		return "", types.ErrPoolExists
	}

	// Create the pool
	pool := types.LiquidityPool{
		Id:            poolID,
		BaseCurrency:  baseCurrency,
		QuoteCurrency: quoteCurrency,
		BaseReserves:  baseAmount,
		QuoteReserves: quoteAmount,
		TotalShares:   sdk.ZeroInt(),
		IsActive:      true,
		Providers:     []types.LiquidityProvider{},
	}

	// Calculate initial LP tokens
	// For initial deposit, LP tokens = sqrt(baseAmount * quoteAmount)
	lpTokens := sdk.NewIntFromBigInt(baseAmount.Mul(quoteAmount).BigInt())
	pool.TotalShares = lpTokens

	return poolID, k.SetLiquidityPool(ctx, pool)
}

// AddLiquidity adds liquidity to an existing pool
func (k Keeper) AddLiquidity(
	ctx context.Context,
	poolID string,
	provider string,
	baseAmount, quoteAmount sdk.Int,
	minLPTokens sdk.Int,
) (sdk.Int, error) {
	// Get the pool
	pool, err := k.GetLiquidityPool(ctx, poolID)
	if err != nil {
		return sdk.ZeroInt(), err
	}

	if !pool.IsActive {
		return sdk.ZeroInt(), types.ErrPoolInactive
	}

	// Validate amounts
	if baseAmount.IsZero() || quoteAmount.IsZero() {
		return sdk.ZeroInt(), types.ErrInvalidLiquidityAmount
	}

	// Calculate optimal amounts to maintain ratio
	optimalQuoteAmount := baseAmount.Mul(pool.QuoteReserves).Quo(pool.BaseReserves)
	optimalBaseAmount := quoteAmount.Mul(pool.BaseReserves).Quo(pool.QuoteReserves)

	var actualBaseAmount, actualQuoteAmount sdk.Int

	if optimalQuoteAmount.LTE(quoteAmount) {
		// Use base amount as reference
		actualBaseAmount = baseAmount
		actualQuoteAmount = optimalQuoteAmount
	} else {
		// Use quote amount as reference
		actualBaseAmount = optimalBaseAmount
		actualQuoteAmount = quoteAmount
	}

	// Calculate LP tokens to mint
	// LP tokens = min(baseAmount * totalShares / baseReserves, quoteAmount * totalShares / quoteReserves)
	var lpTokens sdk.Int
	if pool.TotalShares.IsZero() {
		// First deposit
		lpTokens = sdk.NewIntFromBigInt(actualBaseAmount.Mul(actualQuoteAmount).BigInt())
	} else {
		lpTokensFromBase := actualBaseAmount.Mul(pool.TotalShares).Quo(pool.BaseReserves)
		lpTokensFromQuote := actualQuoteAmount.Mul(pool.TotalShares).Quo(pool.QuoteReserves)
		
		if lpTokensFromBase.LT(lpTokensFromQuote) {
			lpTokens = lpTokensFromBase
		} else {
			lpTokens = lpTokensFromQuote
		}
	}

	// Check minimum LP tokens requirement
	if lpTokens.LT(minLPTokens) {
		return sdk.ZeroInt(), fmt.Errorf("insufficient LP tokens: got %s, minimum %s", lpTokens, minLPTokens)
	}

	// Update pool reserves and shares
	pool.BaseReserves = pool.BaseReserves.Add(actualBaseAmount)
	pool.QuoteReserves = pool.QuoteReserves.Add(actualQuoteAmount)
	pool.TotalShares = pool.TotalShares.Add(lpTokens)

	// Update or add liquidity provider
	k.updateLiquidityProvider(ctx, &pool, provider, lpTokens, true)

	// Save updated pool
	if err := k.SetLiquidityPool(ctx, pool); err != nil {
		return sdk.ZeroInt(), err
	}

	return lpTokens, nil
}

// RemoveLiquidity removes liquidity from a pool
func (k Keeper) RemoveLiquidity(
	ctx context.Context,
	poolID string,
	provider string,
	lpTokens sdk.Int,
	minBaseAmount, minQuoteAmount sdk.Int,
) (sdk.Int, sdk.Int, error) {
	// Get the pool
	pool, err := k.GetLiquidityPool(ctx, poolID)
	if err != nil {
		return sdk.ZeroInt(), sdk.ZeroInt(), err
	}

	// Validate LP tokens
	if lpTokens.IsZero() || lpTokens.IsNegative() {
		return sdk.ZeroInt(), sdk.ZeroInt(), types.ErrInvalidLPTokens
	}

	// Check provider has enough LP tokens
	providerLP := k.getLiquidityProviderTokens(ctx, pool, provider)
	if providerLP.LT(lpTokens) {
		return sdk.ZeroInt(), sdk.ZeroInt(), fmt.Errorf("insufficient LP tokens: have %s, want %s", providerLP, lpTokens)
	}

	// Calculate amounts to return
	baseAmount := lpTokens.Mul(pool.BaseReserves).Quo(pool.TotalShares)
	quoteAmount := lpTokens.Mul(pool.QuoteReserves).Quo(pool.TotalShares)

	// Check minimum amounts
	if baseAmount.LT(minBaseAmount) {
		return sdk.ZeroInt(), sdk.ZeroInt(), fmt.Errorf("insufficient base amount: got %s, minimum %s", baseAmount, minBaseAmount)
	}
	if quoteAmount.LT(minQuoteAmount) {
		return sdk.ZeroInt(), sdk.ZeroInt(), fmt.Errorf("insufficient quote amount: got %s, minimum %s", quoteAmount, minQuoteAmount)
	}

	// Update pool reserves and shares
	pool.BaseReserves = pool.BaseReserves.Sub(baseAmount)
	pool.QuoteReserves = pool.QuoteReserves.Sub(quoteAmount)
	pool.TotalShares = pool.TotalShares.Sub(lpTokens)

	// Update liquidity provider
	k.updateLiquidityProvider(ctx, &pool, provider, lpTokens, false)

	// Save updated pool
	if err := k.SetLiquidityPool(ctx, pool); err != nil {
		return sdk.ZeroInt(), sdk.ZeroInt(), err
	}

	return baseAmount, quoteAmount, nil
}

// GetExchangeRate calculates the exchange rate between two currencies using liquidity pools
func (k Keeper) GetExchangeRate(ctx context.Context, baseCurrency, quoteCurrency string) (sdk.Dec, error) {
	// Try direct pool
	poolID := k.generatePoolID(baseCurrency, quoteCurrency)
	if pool, err := k.GetLiquidityPool(ctx, poolID); err == nil && pool.IsActive {
		if pool.BaseReserves.IsZero() || pool.QuoteReserves.IsZero() {
			return sdk.ZeroDec(), types.ErrInsufficientLiquidity
		}
		return sdk.NewDecFromInt(pool.QuoteReserves).Quo(sdk.NewDecFromInt(pool.BaseReserves)), nil
	}

	// Try reverse pool
	reversePoolID := k.generatePoolID(quoteCurrency, baseCurrency)
	if pool, err := k.GetLiquidityPool(ctx, reversePoolID); err == nil && pool.IsActive {
		if pool.BaseReserves.IsZero() || pool.QuoteReserves.IsZero() {
			return sdk.ZeroDec(), types.ErrInsufficientLiquidity
		}
		return sdk.NewDecFromInt(pool.BaseReserves).Quo(sdk.NewDecFromInt(pool.QuoteReserves)), nil
	}

	// TODO: Implement indirect routing through intermediate currencies (e.g., USD)
	return sdk.ZeroDec(), types.ErrExchangeRateNotFound
}

// SwapCurrency performs a currency swap using liquidity pools
func (k Keeper) SwapCurrency(
	ctx context.Context,
	poolID string,
	inputAmount sdk.Int,
	inputCurrency string,
	minOutputAmount sdk.Int,
) (sdk.Int, error) {
	// Get the pool
	pool, err := k.GetLiquidityPool(ctx, poolID)
	if err != nil {
		return sdk.ZeroInt(), err
	}

	if !pool.IsActive {
		return sdk.ZeroInt(), types.ErrPoolInactive
	}

	// Determine input and output reserves
	var inputReserves, outputReserves sdk.Int
	if inputCurrency == pool.BaseCurrency {
		inputReserves = pool.BaseReserves
		outputReserves = pool.QuoteReserves
	} else if inputCurrency == pool.QuoteCurrency {
		inputReserves = pool.QuoteReserves
		outputReserves = pool.BaseReserves
	} else {
		return sdk.ZeroInt(), fmt.Errorf("currency %s not supported in pool %s", inputCurrency, poolID)
	}

	// Calculate output amount using constant product formula (x * y = k)
	// outputAmount = outputReserves * inputAmount / (inputReserves + inputAmount)
	outputAmount := outputReserves.Mul(inputAmount).Quo(inputReserves.Add(inputAmount))

	// Check minimum output amount
	if outputAmount.LT(minOutputAmount) {
		return sdk.ZeroInt(), types.ErrSlippageExceeded
	}

	// Update reserves
	if inputCurrency == pool.BaseCurrency {
		pool.BaseReserves = pool.BaseReserves.Add(inputAmount)
		pool.QuoteReserves = pool.QuoteReserves.Sub(outputAmount)
	} else {
		pool.QuoteReserves = pool.QuoteReserves.Add(inputAmount)
		pool.BaseReserves = pool.BaseReserves.Sub(outputAmount)
	}

	// Save updated pool
	if err := k.SetLiquidityPool(ctx, pool); err != nil {
		return sdk.ZeroInt(), err
	}

	return outputAmount, nil
}

// CheckLiquidityAvailability checks if sufficient liquidity is available for a swap
func (k Keeper) CheckLiquidityAvailability(
	ctx context.Context,
	amount sdk.Int,
	fromCurrency, toCurrency string,
) bool {
	// Get exchange rate to estimate required liquidity
	rate, err := k.GetExchangeRate(ctx, fromCurrency, toCurrency)
	if err != nil {
		return false
	}

	// Estimate output amount
	estimatedOutput := sdk.NewDecFromInt(amount).Mul(rate).TruncateInt()

	// Check if pools have sufficient reserves
	poolID := k.generatePoolID(fromCurrency, toCurrency)
	if pool, err := k.GetLiquidityPool(ctx, poolID); err == nil && pool.IsActive {
		if fromCurrency == pool.BaseCurrency {
			return pool.QuoteReserves.GT(estimatedOutput)
		} else {
			return pool.BaseReserves.GT(estimatedOutput)
		}
	}

	// Check reverse pool
	reversePoolID := k.generatePoolID(toCurrency, fromCurrency)
	if pool, err := k.GetLiquidityPool(ctx, reversePoolID); err == nil && pool.IsActive {
		if toCurrency == pool.BaseCurrency {
			return pool.BaseReserves.GT(estimatedOutput)
		} else {
			return pool.QuoteReserves.GT(estimatedOutput)
		}
	}

	return false
}

// Helper functions

// generatePoolID generates a consistent pool ID for a currency pair
func (k Keeper) generatePoolID(baseCurrency, quoteCurrency string) string {
	// Always use alphabetical order to ensure consistency
	currencies := []string{baseCurrency, quoteCurrency}
	if currencies[0] > currencies[1] {
		currencies[0], currencies[1] = currencies[1], currencies[0]
	}
	return fmt.Sprintf("%s-%s", currencies[0], currencies[1])
}

// updateLiquidityProvider updates or adds a liquidity provider in the pool
func (k Keeper) updateLiquidityProvider(
	ctx context.Context,
	pool *types.LiquidityPool,
	provider string,
	lpTokens sdk.Int,
	isAdd bool,
) {
	// Find existing provider
	for i, lp := range pool.Providers {
		if lp.Address == provider {
			if isAdd {
				pool.Providers[i].LpTokens = pool.Providers[i].LpTokens.Add(lpTokens)
			} else {
				pool.Providers[i].LpTokens = pool.Providers[i].LpTokens.Sub(lpTokens)
			}
			return
		}
	}

	// Add new provider
	if isAdd {
		newProvider := types.LiquidityProvider{
			Address:  provider,
			LpTokens: lpTokens,
		}
		pool.Providers = append(pool.Providers, newProvider)
	}
}

// getLiquidityProviderTokens gets the LP tokens for a specific provider
func (k Keeper) getLiquidityProviderTokens(
	ctx context.Context,
	pool types.LiquidityPool,
	provider string,
) sdk.Int {
	for _, lp := range pool.Providers {
		if lp.Address == provider {
			return lp.LpTokens
		}
	}
	return sdk.ZeroInt()
}

// GetPoolStatistics returns statistics for a liquidity pool
func (k Keeper) GetPoolStatistics(ctx context.Context, poolID string) (types.PoolStatistics, error) {
	pool, err := k.GetLiquidityPool(ctx, poolID)
	if err != nil {
		return types.PoolStatistics{}, err
	}

	// Calculate TVL in base currency terms
	tvl := pool.BaseReserves.Mul(sdk.NewInt(2)) // Assuming equal value on both sides

	stats := types.PoolStatistics{
		PoolId:           pool.Id,
		BaseCurrency:     pool.BaseCurrency,
		QuoteCurrency:    pool.QuoteCurrency,
		BaseReserves:     pool.BaseReserves,
		QuoteReserves:    pool.QuoteReserves,
		TotalShares:      pool.TotalShares,
		Tvl:              tvl,
		ProviderCount:    uint64(len(pool.Providers)),
		IsActive:         pool.IsActive,
	}

	// Calculate current exchange rate
	if !pool.BaseReserves.IsZero() && !pool.QuoteReserves.IsZero() {
		stats.CurrentRate = sdk.NewDecFromInt(pool.QuoteReserves).Quo(sdk.NewDecFromInt(pool.BaseReserves))
	}

	return stats, nil
}

// SetLiquidityProvider stores a liquidity provider record
func (k Keeper) SetLiquidityProvider(ctx context.Context, poolID string, provider types.LiquidityProvider) error {
	store := k.GetStore(ctx)
	key := types.LiquidityProviderKey(poolID, provider.Address)
	bz := k.cdc.MustMarshal(&provider)
	store.Set(key, bz)
	return nil
}

// GetLiquidityProvider retrieves a liquidity provider record
func (k Keeper) GetLiquidityProvider(ctx context.Context, poolID string, providerAddress string) (types.LiquidityProvider, error) {
	store := k.GetStore(ctx)
	key := types.LiquidityProviderKey(poolID, providerAddress)
	bz := store.Get(key)
	if bz == nil {
		return types.LiquidityProvider{}, types.ErrLiquidityProviderNotFound
	}

	var provider types.LiquidityProvider
	k.cdc.MustUnmarshal(bz, &provider)
	return provider, nil
}