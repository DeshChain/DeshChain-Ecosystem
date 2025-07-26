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
	"github.com/DeshChain/DeshChain-Ecosystem/x/namo/types"
)

// RegisterInvariants registers all NAMO module invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "token-supply", TokenSupplyInvariant(k))
	ir.RegisterRoute(types.ModuleName, "vesting-total", VestingTotalInvariant(k))
	ir.RegisterRoute(types.ModuleName, "allocation-sum", AllocationSumInvariant(k))
}

// TokenSupplyInvariant checks that the total token supply matches the sum of all allocations
func TokenSupplyInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		supply, found := k.GetTokenSupply(ctx)
		if !found {
			return sdk.FormatInvariant(
				types.ModuleName, "token-supply",
				"token supply not found in state",
			), true
		}

		totalSupply, ok := sdk.NewIntFromString(supply.TotalSupply)
		if !ok {
			return sdk.FormatInvariant(
				types.ModuleName, "token-supply",
				"invalid total supply format",
			), true
		}

		// Calculate sum of all allocations
		allocations := []string{
			supply.FounderAllocation,
			supply.TeamAllocation,
			supply.CommunityAllocation,
			supply.DevelopmentAllocation,
			supply.LiquidityAllocation,
			supply.PublicSaleAllocation,
			supply.DAOTreasuryAllocation,
			supply.CoFounderAllocation,
			supply.OperationsAllocation,
			supply.AngelAllocation,
		}

		allocationSum := sdk.ZeroInt()
		for _, allocation := range allocations {
			amount, ok := sdk.NewIntFromString(allocation)
			if !ok {
				return sdk.FormatInvariant(
					types.ModuleName, "token-supply",
					fmt.Sprintf("invalid allocation format: %s", allocation),
				), true
			}
			allocationSum = allocationSum.Add(amount)
		}

		if !totalSupply.Equal(allocationSum) {
			return sdk.FormatInvariant(
				types.ModuleName, "token-supply",
				fmt.Sprintf("total supply (%s) does not equal sum of allocations (%s)", totalSupply, allocationSum),
			), true
		}

		return sdk.FormatInvariant(types.ModuleName, "token-supply", "token supply is valid"), false
	}
}

// VestingTotalInvariant checks that the total vesting amounts are consistent
func VestingTotalInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		schedules := k.GetAllVestingSchedules(ctx)
		
		totalVesting := sdk.ZeroInt()
		totalClaimed := sdk.ZeroInt()

		for _, schedule := range schedules {
			vestingAmount, ok := sdk.NewIntFromString(schedule.TotalAmount)
			if !ok {
				return sdk.FormatInvariant(
					types.ModuleName, "vesting-total",
					fmt.Sprintf("invalid vesting amount for %s: %s", schedule.Beneficiary, schedule.TotalAmount),
				), true
			}

			claimedAmount, ok := sdk.NewIntFromString(schedule.ClaimedAmount)
			if !ok {
				return sdk.FormatInvariant(
					types.ModuleName, "vesting-total",
					fmt.Sprintf("invalid claimed amount for %s: %s", schedule.Beneficiary, schedule.ClaimedAmount),
				), true
			}

			if claimedAmount.GT(vestingAmount) {
				return sdk.FormatInvariant(
					types.ModuleName, "vesting-total",
					fmt.Sprintf("claimed amount (%s) exceeds vesting amount (%s) for %s", 
						claimedAmount, vestingAmount, schedule.Beneficiary),
				), true
			}

			totalVesting = totalVesting.Add(vestingAmount)
			totalClaimed = totalClaimed.Add(claimedAmount)
		}

		// Check that vesting pool has enough tokens
		vestingPoolAddr := k.accountKeeper.GetModuleAddress(types.VestingPoolName)
		if vestingPoolAddr != nil {
			poolBalance := k.bankKeeper.GetBalance(ctx, vestingPoolAddr, types.TokenDenom)
			expectedBalance := totalVesting.Sub(totalClaimed)
			
			if !poolBalance.Amount.Equal(expectedBalance) {
				return sdk.FormatInvariant(
					types.ModuleName, "vesting-total",
					fmt.Sprintf("vesting pool balance (%s) does not match expected balance (%s)", 
						poolBalance.Amount, expectedBalance),
				), true
			}
		}

		return sdk.FormatInvariant(types.ModuleName, "vesting-total", "vesting totals are consistent"), false
	}
}

// AllocationSumInvariant checks that all allocation percentages sum to 100%
func AllocationSumInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		params := k.GetParams(ctx)
		
		allocations := []string{
			params.FounderAllocation,
			params.TeamAllocation,
			params.CommunityAllocation,
			params.DevelopmentAllocation,
			params.LiquidityAllocation,
			params.PublicSaleAllocation,
			params.DAOTreasuryAllocation,
			params.CoFounderAllocation,
			params.OperationsAllocation,
			params.AngelAllocation,
		}

		totalPercentage := sdk.ZeroDec()
		for i, allocation := range allocations {
			percentage, err := sdk.NewDecFromStr(allocation)
			if err != nil {
				return sdk.FormatInvariant(
					types.ModuleName, "allocation-sum",
					fmt.Sprintf("invalid allocation percentage at index %d: %s", i, allocation),
				), true
			}
			totalPercentage = totalPercentage.Add(percentage)
		}

		expectedTotal := sdk.NewDec(100) // 100%
		if !totalPercentage.Equal(expectedTotal) {
			return sdk.FormatInvariant(
				types.ModuleName, "allocation-sum",
				fmt.Sprintf("allocation percentages sum to %s, expected 100.0", totalPercentage),
			), true
		}

		return sdk.FormatInvariant(types.ModuleName, "allocation-sum", "allocation percentages sum to 100%"), false
	}
}

// TokenBalanceInvariant checks that token balances across all accounts are consistent
func TokenBalanceInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		supply, found := k.GetTokenSupply(ctx)
		if !found {
			return sdk.FormatInvariant(
				types.ModuleName, "token-balance",
				"token supply not found",
			), true
		}

		totalSupply, ok := sdk.NewIntFromString(supply.TotalSupply)
		if !ok {
			return sdk.FormatInvariant(
				types.ModuleName, "token-balance",
				"invalid total supply format",
			), true
		}

		// Calculate total tokens in circulation
		circulatingSupply := k.GetCirculatingSupply(ctx)
		
		// Get tokens in vesting pool
		vestingPoolAddr := k.accountKeeper.GetModuleAddress(types.VestingPoolName)
		vestingTokens := sdk.ZeroInt()
		if vestingPoolAddr != nil {
			vestingBalance := k.bankKeeper.GetBalance(ctx, vestingPoolAddr, types.TokenDenom)
			vestingTokens = vestingBalance.Amount
		}

		// Get burned tokens
		burnedTokens := k.GetTotalBurnedTokens(ctx)

		// Total should equal circulating + vesting + burned
		calculatedTotal := circulatingSupply.Add(vestingTokens).Add(burnedTokens)
		
		if !totalSupply.Equal(calculatedTotal) {
			return sdk.FormatInvariant(
				types.ModuleName, "token-balance",
				fmt.Sprintf("total supply (%s) does not equal calculated total (%s = %s circulating + %s vesting + %s burned)", 
					totalSupply, calculatedTotal, circulatingSupply, vestingTokens, burnedTokens),
			), true
		}

		return sdk.FormatInvariant(types.ModuleName, "token-balance", "token balances are consistent"), false
	}
}