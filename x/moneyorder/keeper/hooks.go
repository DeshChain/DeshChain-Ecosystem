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

// Hooks wrapper struct for Money Order keeper
type Hooks struct {
	k Keeper
}

var _ types.MoneyOrderHooks = Hooks{}
var _ GramPensionHooks = Hooks{}

// Return the wrapper struct
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// AfterSurakshaContribution - called after a pension contribution is made
// This hook integrates with the Gram Pension module to add liquidity
func (h Hooks) AfterSurakshaContribution(
	ctx sdk.Context,
	pensionAccountId string,
	contributor sdk.AccAddress,
	contribution sdk.Coin,
	villagePostalCode string,
) error {
	// Find the village pool for this postal code
	villagePool, found := h.k.GetVillagePoolByPostalCode(ctx, villagePostalCode)
	if !found {
		// If no village pool exists, skip liquidity provision
		// In production, might create a default pool
		return nil
	}

	// Check if unified pool exists for this village
	unifiedPool, found := h.k.GetUnifiedPoolByVillageId(ctx, villagePool.PoolId)
	if !found {
		// Create unified pool if it doesn't exist
		initialLiquidity := sdk.NewCoins(contribution)
		pool, err := h.k.CreateUnifiedLiquidityPool(ctx, villagePool.PoolId, initialLiquidity)
		if err != nil {
			return err
		}
		unifiedPool = *pool
	}

	// Add pension contribution to unified pool
	return h.k.AddSurakshaContribution(
		ctx,
		unifiedPool.PoolId,
		contribution,
		pensionAccountId,
	)
}

// AfterSurakshaMaturity - called when pension reaches maturity
// This hook handles liquidity withdrawal and return calculation
func (h Hooks) AfterSurakshaMaturity(
	ctx sdk.Context,
	pensionAccountId string,
	beneficiary sdk.AccAddress,
	maturityAmount sdk.Coin,
) error {
	// Find all liquidity positions for this pension account
	positions := h.k.GetSurakshaLiquidityPositions(ctx, pensionAccountId)
	
	totalReturns := sdk.NewCoins()
	
	for _, position := range positions {
		// Process maturity for each position
		if err := h.k.processPensionLiquidityMaturity(ctx, position); err != nil {
			ctx.Logger().Error("failed to process pension liquidity maturity",
				"pension_account", pensionAccountId,
				"error", err)
			continue
		}
		
		// Accumulate returns
		totalReturns = totalReturns.Add(position.RewardsEarned)
	}

	// The actual transfer of maturity amount + returns would be handled by the pension module
	// This hook just ensures liquidity is properly released
	
	// Emit event for tracking
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSurakshaMaturityProcessed,
			sdk.NewAttribute(types.AttributeKeyPensionAccount, pensionAccountId),
			sdk.NewAttribute("beneficiary", beneficiary.String()),
			sdk.NewAttribute(types.AttributeKeyRewards, totalReturns.String()),
		),
	)

	return nil
}

// BeforeAgriLoanApproval - called before agricultural loan approval
// This hook checks liquidity availability in unified pool
func (h Hooks) BeforeAgriLoanApproval(
	ctx sdk.Context,
	borrower sdk.AccAddress,
	loanAmount sdk.Coin,
	villagePostalCode string,
) error {
	// Find village pool
	villagePool, found := h.k.GetVillagePoolByPostalCode(ctx, villagePostalCode)
	if !found {
		return types.ErrVillagePoolNotFound
	}

	// Get unified pool
	unifiedPool, found := h.k.GetUnifiedPoolByVillageId(ctx, villagePool.PoolId)
	if !found {
		return types.ErrPoolNotFound
	}

	// Check if sufficient lending capacity exists
	availableForLending := unifiedPool.AgriLendingPool.AmountOf(loanAmount.Denom)
	if availableForLending.LT(loanAmount.Amount) {
		return types.ErrInsufficientLiquidity
	}

	// Additional checks could include:
	// - Borrower's credit history
	// - Total exposure limits
	// - Seasonal constraints

	return nil
}

// AfterAgriLoanDisbursement - called after agricultural loan is disbursed
// This hook updates the unified pool liquidity
func (h Hooks) AfterAgriLoanDisbursement(
	ctx sdk.Context,
	loanId string,
	borrower sdk.AccAddress,
	loanAmount sdk.Coin,
	loanType string,
	duration uint32,
	villagePostalCode string,
) error {
	// Find village pool
	villagePool, found := h.k.GetVillagePoolByPostalCode(ctx, villagePostalCode)
	if !found {
		return types.ErrVillagePoolNotFound
	}

	// Get unified pool
	unifiedPool, found := h.k.GetUnifiedPoolByVillageId(ctx, villagePool.PoolId)
	if !found {
		return types.ErrPoolNotFound
	}

	// Process the loan through unified pool
	return h.k.ProcessAgriLoan(
		ctx,
		unifiedPool.PoolId,
		loanAmount,
		borrower,
		loanType,
		duration,
	)
}

// AfterAgriLoanRepayment - called after agricultural loan repayment
// This hook processes repayment and distributes profits
func (h Hooks) AfterAgriLoanRepayment(
	ctx sdk.Context,
	loanId string,
	repaymentAmount sdk.Coin,
) error {
	// Convert string loanId to uint64 (in production, handle properly)
	loanIdNum, err := sdk.ParseUint(loanId)
	if err != nil {
		return err
	}

	// Process repayment through unified pool
	return h.k.ProcessLoanRepayment(ctx, loanIdNum.Uint64(), repaymentAmount)
}

// AfterTradingFeeCollection - called after trading fees are collected
// This hook distributes fees to the unified pool for pension returns
func (h Hooks) AfterTradingFeeCollection(
	ctx sdk.Context,
	poolId uint64,
	tradingFees sdk.Coins,
) error {
	// Check if this is a village pool
	villagePool, found := h.k.GetVillagePool(ctx, poolId)
	if !found {
		// Not a village pool, skip unified pool distribution
		return nil
	}

	// Get unified pool for this village
	unifiedPool, found := h.k.GetUnifiedPoolByVillageId(ctx, villagePool.PoolId)
	if !found {
		// No unified pool, use standard fee distribution
		return nil
	}

	// Record DEX revenue in unified pool
	return h.k.RecordDexRevenue(ctx, unifiedPool.PoolId, tradingFees)
}

// BeforePoolCreation - called before creating a new pool
// This hook can enforce village-specific rules
func (h Hooks) BeforePoolCreation(
	ctx sdk.Context,
	creator sdk.AccAddress,
	poolType string,
) error {
	// Check if creator is authorized (e.g., panchayat head for village pools)
	if poolType == types.PoolTypeVillage {
		// Additional validation for village pool creation
		// Could check if creator is registered panchayat head
	}

	return nil
}

// AfterPoolCreation - called after a pool is created
// This hook can trigger additional setup
func (h Hooks) AfterPoolCreation(
	ctx sdk.Context,
	poolId uint64,
	poolType string,
	creator sdk.AccAddress,
) error {
	if poolType == types.PoolTypeVillage {
		// Automatically create unified liquidity pool for village
		pool, found := h.k.GetVillagePool(ctx, poolId)
		if found && pool.Verified {
			// Initialize with minimal liquidity
			initialLiquidity := sdk.NewCoins(sdk.NewCoin("unamo", sdk.NewInt(1000000)))
			_, err := h.k.CreateUnifiedLiquidityPool(ctx, poolId, initialLiquidity)
			return err
		}
	}

	return nil
}

// MonthlyRevenueDistribution - called monthly to distribute revenues
// This hook calculates and allocates returns for pension accounts
func (h Hooks) MonthlyRevenueDistribution(ctx sdk.Context) error {
	// Iterate through all unified pools
	h.k.IterateUnifiedPools(ctx, func(pool UnifiedLiquidityPool) bool {
		// Calculate pension returns from combined revenue
		pensionReturns, err := h.k.CalculatePensionReturns(ctx, pool.PoolId)
		if err != nil {
			ctx.Logger().Error("failed to calculate pension returns",
				"pool_id", pool.PoolId,
				"error", err)
			return false
		}

		// Distribute returns to pension accounts
		// This would integrate with the pension module
		
		// Reset monthly revenue counters
		pool.MonthlyDexRevenue = sdk.NewCoins()
		pool.MonthlyLendingRevenue = sdk.NewCoins()
		h.k.SetUnifiedLiquidityPool(ctx, pool)

		// Emit distribution event
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeMonthlyDistribution,
				sdk.NewAttribute(types.AttributeKeyUnifiedPoolId, fmt.Sprintf("%d", pool.PoolId)),
				sdk.NewAttribute("pension_returns", pensionReturns.String()),
				sdk.NewAttribute("total_revenue", pool.TotalRevenue.String()),
			),
		)

		return false // continue iteration
	})

	return nil
}

// Helper functions

func (k Keeper) GetVillagePoolByPostalCode(ctx sdk.Context, postalCode string) (types.VillagePool, bool) {
	var foundPool types.VillagePool
	found := false

	k.IterateVillagePools(ctx, func(pool types.VillagePool) bool {
		if pool.PostalCode == postalCode {
			foundPool = pool
			found = true
			return true // stop iteration
		}
		return false
	})

	return foundPool, found
}

func (k Keeper) GetUnifiedPoolByVillageId(ctx sdk.Context, villagePoolId uint64) (UnifiedLiquidityPool, bool) {
	var foundPool UnifiedLiquidityPool
	found := false

	k.IterateUnifiedPools(ctx, func(pool UnifiedLiquidityPool) bool {
		if pool.VillagePoolId == villagePoolId {
			foundPool = pool
			found = true
			return true
		}
		return false
	})

	return foundPool, found
}

func (k Keeper) GetSurakshaLiquidityPositions(ctx sdk.Context, pensionAccountId string) []PensionLiquidity {
	var positions []PensionLiquidity

	k.IteratePensionLiquidity(ctx, func(pl PensionLiquidity) bool {
		if pl.PensionAccountId == pensionAccountId {
			positions = append(positions, pl)
		}
		return false
	})

	return positions
}

func (k Keeper) IterateUnifiedPools(ctx sdk.Context, cb func(UnifiedLiquidityPool) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixUnifiedPool)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var pool UnifiedLiquidityPool
		k.cdc.MustUnmarshal(iterator.Value(), &pool)
		if cb(pool) {
			break
		}
	}
}