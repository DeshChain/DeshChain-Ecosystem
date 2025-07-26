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
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/DeshChain/DeshChain-Ecosystem/x/moneyorder/types"
)

// PensionLiquidityConfig defines how pension funds are utilized for liquidity
type PensionLiquidityConfig struct {
	// Percentage of pension contributions that can be used for liquidity (e.g., 80%)
	LiquidityUtilizationRate sdk.Dec
	// Minimum reserve to maintain for pension payouts (e.g., 20%)
	MinimumReserveRate sdk.Dec
	// Lock period for pension liquidity (12 months)
	LockPeriodMonths uint32
	// Interest rate bonus for pension liquidity providers
	PensionLiquidityBonus sdk.Dec
}

// DefaultPensionLiquidityConfig returns default configuration
func DefaultPensionLiquidityConfig() PensionLiquidityConfig {
	return PensionLiquidityConfig{
		LiquidityUtilizationRate: sdk.NewDecWithPrec(80, 2), // 80% can be used for liquidity
		MinimumReserveRate:       sdk.NewDecWithPrec(20, 2), // 20% must be kept as reserve
		LockPeriodMonths:         12,                         // 12-month lock period
		PensionLiquidityBonus:    sdk.NewDecWithPrec(5, 2),  // 5% bonus APY for pension liquidity
	}
}

// AddPensionLiquidity adds pension contributions to village pool liquidity
func (k Keeper) AddPensionLiquidity(
	ctx sdk.Context,
	villagePoolId uint64,
	surakshaContribution sdk.Coin,
	contributorAddr sdk.AccAddress,
	pensionAccountId string,
) error {
	// Get village pool
	pool, found := k.GetVillagePool(ctx, villagePoolId)
	if !found {
		return types.ErrVillagePoolNotFound
	}

	if !pool.Active {
		return types.ErrVillagePoolInactive
	}

	// Get pension liquidity config
	config := DefaultPensionLiquidityConfig()

	// Calculate usable liquidity (80% of contribution)
	usableLiquidity := surakshaContribution.Amount.ToDec().Mul(config.LiquidityUtilizationRate).TruncateInt()
	reserveAmount := surakshaContribution.Amount.Sub(usableLiquidity)

	// Create liquidity coin
	liquidityCoin := sdk.NewCoin(surakshaContribution.Denom, usableLiquidity)
	reserveCoin := sdk.NewCoin(surakshaContribution.Denom, reserveAmount)

	// Add to pool liquidity
	pool.TotalLiquidity = pool.TotalLiquidity.Add(liquidityCoin)
	pool.AvailableLiquidity = pool.AvailableLiquidity.Add(liquidityCoin)
	
	// Track pension liquidity separately for accounting
	k.SetSurakshaLiquidity(ctx, PensionLiquidity{
		PensionAccountId:  pensionAccountId,
		VillagePoolId:     villagePoolId,
		ContributorAddr:   contributorAddr.String(),
		LiquidityAmount:   liquidityCoin,
		ReserveAmount:     reserveCoin,
		ContributionMonth: uint32(ctx.BlockTime().Month()),
		ContributionYear:  uint32(ctx.BlockTime().Year()),
		MaturityTime:      ctx.BlockTime().AddDate(0, int(config.LockPeriodMonths), 0),
		IsActive:          true,
		BonusRate:         config.PensionLiquidityBonus,
	})

	// Update pool
	k.SetVillagePool(ctx, pool)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypePensionLiquidityAdded,
			sdk.NewAttribute(types.AttributeKeyPoolId, fmt.Sprintf("%d", villagePoolId)),
			sdk.NewAttribute(types.AttributeKeyPensionAccount, pensionAccountId),
			sdk.NewAttribute(types.AttributeKeyLiquidity, liquidityCoin.String()),
			sdk.NewAttribute(types.AttributeKeyReserve, reserveCoin.String()),
			sdk.NewAttribute("contributor", contributorAddr.String()),
			sdk.NewAttribute("maturity_months", fmt.Sprintf("%d", config.LockPeriodMonths)),
		),
	)

	return nil
}

// RotatePensionLiquidity handles the 12-month rotation of pension liquidity
func (k Keeper) RotatePensionLiquidity(ctx sdk.Context) {
	// This function is called monthly to handle maturing pension liquidity
	
	// Iterate through all pension liquidity positions
	k.IteratePensionLiquidity(ctx, func(pl PensionLiquidity) bool {
		// Check if liquidity has matured (12 months passed)
		if ctx.BlockTime().After(pl.MaturityTime) && pl.IsActive {
			// Process maturity
			if err := k.processPensionLiquidityMaturity(ctx, pl); err != nil {
				// Log error but continue iteration
				ctx.Logger().Error("failed to process pension liquidity maturity", 
					"pension_account", pl.PensionAccountId, 
					"error", err)
			}
		}
		return false // continue iteration
	})
}

// processPensionLiquidityMaturity handles matured pension liquidity
func (k Keeper) processPensionLiquidityMaturity(ctx sdk.Context, pl PensionLiquidity) error {
	// Get village pool
	pool, found := k.GetVillagePool(ctx, pl.VillagePoolId)
	if !found {
		return types.ErrVillagePoolNotFound
	}

	// Calculate rewards earned from providing liquidity
	// Base rewards from trading fees + bonus rate
	baseRewards := k.calculatePensionLiquidityRewards(ctx, pl)
	bonusRewards := pl.LiquidityAmount.Amount.ToDec().Mul(pl.BonusRate).TruncateInt()
	totalRewards := baseRewards.Add(sdk.NewCoin(pl.LiquidityAmount.Denom, bonusRewards))

	// Remove liquidity from pool
	pool.AvailableLiquidity = pool.AvailableLiquidity.Sub(pl.LiquidityAmount)
	
	// Mark pension liquidity as inactive
	pl.IsActive = false
	pl.MaturedAt = ctx.BlockTime()
	pl.RewardsEarned = totalRewards
	
	// Update records
	k.SetVillagePool(ctx, pool)
	k.SetSurakshaLiquidity(ctx, pl)

	// Transfer liquidity + rewards back to pension account
	// This will be handled by the pension module when paying out maturity
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypePensionLiquidityMatured,
			sdk.NewAttribute(types.AttributeKeyPensionAccount, pl.PensionAccountId),
			sdk.NewAttribute(types.AttributeKeyPoolId, fmt.Sprintf("%d", pl.VillagePoolId)),
			sdk.NewAttribute(types.AttributeKeyLiquidity, pl.LiquidityAmount.String()),
			sdk.NewAttribute(types.AttributeKeyRewards, totalRewards.String()),
			sdk.NewAttribute("maturity_time", pl.MaturityTime.String()),
		),
	)

	return nil
}

// calculatePensionLiquidityRewards calculates rewards earned from providing liquidity
func (k Keeper) calculatePensionLiquidityRewards(ctx sdk.Context, pl PensionLiquidity) sdk.Coin {
	// Get pool's total earnings during the liquidity period
	pool, found := k.GetVillagePool(ctx, pl.VillagePoolId)
	if !found {
		return sdk.NewCoin(pl.LiquidityAmount.Denom, sdk.ZeroInt())
	}

	// Calculate share of pool earnings based on liquidity contribution
	// This is a simplified calculation - in production, track actual earnings
	liquidityShare := pl.LiquidityAmount.Amount.ToDec().Quo(pool.TotalLiquidity.AmountOf(pl.LiquidityAmount.Denom).ToDec())
	
	// Assume pool earned 10% APY from trading fees (this would be tracked accurately)
	annualEarnings := pool.TotalLiquidity.AmountOf(pl.LiquidityAmount.Denom).ToDec().Mul(sdk.NewDecWithPrec(10, 2))
	
	// Pro-rate for actual months (could be less than 12 if early withdrawal)
	monthsActive := uint32(12) // Full term
	if !pl.IsActive && !pl.MaturedAt.IsZero() {
		monthsActive = uint32(pl.MaturedAt.Sub(ctx.BlockTime()).Hours() / 24 / 30)
	}
	
	earnings := annualEarnings.Mul(liquidityShare).Mul(sdk.NewDec(int64(monthsActive))).Quo(sdk.NewDec(12)).TruncateInt()
	
	return sdk.NewCoin(pl.LiquidityAmount.Denom, earnings)
}

// GetSurakshaLiquidityUtilization returns current utilization stats
func (k Keeper) GetSurakshaLiquidityUtilization(ctx sdk.Context, villagePoolId uint64) PensionLiquidityStats {
	stats := PensionLiquidityStats{
		VillagePoolId:        villagePoolId,
		TotalPensionLiquidity: sdk.NewCoins(),
		ActiveContributors:   0,
		MonthlyInflow:        sdk.NewCoins(),
		MonthlyOutflow:       sdk.NewCoins(),
		AverageAPY:           sdk.ZeroDec(),
	}

	// Calculate stats by iterating pension liquidity
	totalAPY := sdk.ZeroDec()
	contributorMap := make(map[string]bool)
	currentMonth := uint32(ctx.BlockTime().Month())
	currentYear := uint32(ctx.BlockTime().Year())

	k.IteratePensionLiquidityByPool(ctx, villagePoolId, func(pl PensionLiquidity) bool {
		if pl.IsActive {
			stats.TotalPensionLiquidity = stats.TotalPensionLiquidity.Add(pl.LiquidityAmount)
			contributorMap[pl.ContributorAddr] = true
			
			// Track monthly flows
			if pl.ContributionMonth == currentMonth && pl.ContributionYear == currentYear {
				stats.MonthlyInflow = stats.MonthlyInflow.Add(pl.LiquidityAmount)
			}
			
			// Add bonus rate to average APY calculation
			totalAPY = totalAPY.Add(pl.BonusRate)
		} else if pl.MaturedAt.Month() == time.Month(currentMonth) && pl.MaturedAt.Year() == int(currentYear) {
			// Track outflows
			stats.MonthlyOutflow = stats.MonthlyOutflow.Add(pl.LiquidityAmount)
		}
		
		return false
	})

	stats.ActiveContributors = uint32(len(contributorMap))
	if stats.ActiveContributors > 0 {
		stats.AverageAPY = totalAPY.Quo(sdk.NewDec(int64(stats.ActiveContributors)))
	}

	return stats
}

// Types for pension liquidity tracking

type PensionLiquidity struct {
	PensionAccountId  string    `json:"pension_account_id"`
	VillagePoolId     uint64    `json:"village_pool_id"`
	ContributorAddr   string    `json:"contributor_addr"`
	LiquidityAmount   sdk.Coin  `json:"liquidity_amount"`
	ReserveAmount     sdk.Coin  `json:"reserve_amount"`
	ContributionMonth uint32    `json:"contribution_month"`
	ContributionYear  uint32    `json:"contribution_year"`
	MaturityTime      time.Time `json:"maturity_time"`
	MaturedAt         time.Time `json:"matured_at,omitempty"`
	IsActive          bool      `json:"is_active"`
	BonusRate         sdk.Dec   `json:"bonus_rate"`
	RewardsEarned     sdk.Coin  `json:"rewards_earned,omitempty"`
}

type PensionLiquidityStats struct {
	VillagePoolId         uint64    `json:"village_pool_id"`
	TotalPensionLiquidity sdk.Coins `json:"total_pension_liquidity"`
	ActiveContributors    uint32    `json:"active_contributors"`
	MonthlyInflow         sdk.Coins `json:"monthly_inflow"`
	MonthlyOutflow        sdk.Coins `json:"monthly_outflow"`
	AverageAPY            sdk.Dec   `json:"average_apy"`
}

// Store functions

func (k Keeper) SetSurakshaLiquidity(ctx sdk.Context, pl PensionLiquidity) {
	store := ctx.KVStore(k.storeKey)
	key := getPensionLiquidityKey(pl.PensionAccountId, pl.VillagePoolId)
	bz := k.cdc.MustMarshal(&pl)
	store.Set(key, bz)
}

func (k Keeper) GetSurakshaLiquidity(ctx sdk.Context, pensionAccountId string, poolId uint64) (PensionLiquidity, bool) {
	store := ctx.KVStore(k.storeKey)
	key := getPensionLiquidityKey(pensionAccountId, poolId)
	bz := store.Get(key)
	if bz == nil {
		return PensionLiquidity{}, false
	}
	
	var pl PensionLiquidity
	k.cdc.MustUnmarshal(bz, &pl)
	return pl, true
}

func (k Keeper) IteratePensionLiquidity(ctx sdk.Context, cb func(PensionLiquidity) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixPensionLiquidity)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var pl PensionLiquidity
		k.cdc.MustUnmarshal(iterator.Value(), &pl)
		if cb(pl) {
			break
		}
	}
}

func (k Keeper) IteratePensionLiquidityByPool(ctx sdk.Context, poolId uint64, cb func(PensionLiquidity) bool) {
	k.IteratePensionLiquidity(ctx, func(pl PensionLiquidity) bool {
		if pl.VillagePoolId == poolId {
			return cb(pl)
		}
		return false
	})
}

func getPensionLiquidityKey(pensionAccountId string, poolId uint64) []byte {
	return append(append(types.KeyPrefixPensionLiquidity, []byte(pensionAccountId)...), sdk.Uint64ToBigEndian(poolId)...)
}