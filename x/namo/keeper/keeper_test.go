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

package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/namo/types"
)

func TestDefaultParams(t *testing.T) {
	params := types.DefaultParams()
	require.NotNil(t, params)
	
	// Validate params
	err := params.Validate()
	require.NoError(t, err)
	
	// Check default values
	require.True(t, params.EnableTokenOperations)
	require.Equal(t, types.TokenDenom, params.TokenDenom)
	require.Equal(t, types.DefaultTotalSupply, params.InitialSupply)
	require.True(t, params.VestingEnabled)
	require.True(t, params.Burnable)
}

func TestVestingScheduleValidation(t *testing.T) {
	testCases := []struct {
		name      string
		schedule  types.VestingSchedule
		shouldErr bool
	}{
		{
			name: "valid schedule",
			schedule: types.VestingSchedule{
				Beneficiary:     "desh1abc123...",
				TotalAmount:     "1000000000000", // 1M tokens
				ClaimedAmount:   "0",
				CliffTime:       time.Now().Unix() + 3600*24*365,     // 1 year from now
				EndTime:         time.Now().Unix() + 3600*24*365*2,   // 2 years from now
				VestingCategory: types.VestingCategoryFounder,
				CreatedAt:       time.Now().Unix(),
			},
			shouldErr: false,
		},
		{
			name: "invalid beneficiary",
			schedule: types.VestingSchedule{
				Beneficiary:     "",
				TotalAmount:     "1000000000000",
				ClaimedAmount:   "0",
				CliffTime:       time.Now().Unix() + 3600*24*365,
				EndTime:         time.Now().Unix() + 3600*24*365*2,
				VestingCategory: types.VestingCategoryFounder,
				CreatedAt:       time.Now().Unix(),
			},
			shouldErr: true,
		},
		{
			name: "invalid amount",
			schedule: types.VestingSchedule{
				Beneficiary:     "desh1abc123...",
				TotalAmount:     "0",
				ClaimedAmount:   "0",
				CliffTime:       time.Now().Unix() + 3600*24*365,
				EndTime:         time.Now().Unix() + 3600*24*365*2,
				VestingCategory: types.VestingCategoryFounder,
				CreatedAt:       time.Now().Unix(),
			},
			shouldErr: true,
		},
		{
			name: "invalid time order",
			schedule: types.VestingSchedule{
				Beneficiary:     "desh1abc123...",
				TotalAmount:     "1000000000000",
				ClaimedAmount:   "0",
				CliffTime:       time.Now().Unix() + 3600*24*365*2,   // 2 years from now
				EndTime:         time.Now().Unix() + 3600*24*365,     // 1 year from now
				VestingCategory: types.VestingCategoryFounder,
				CreatedAt:       time.Now().Unix(),
			},
			shouldErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := types.ValidateVestingSchedule(tc.schedule)
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestTokenSupplyValidation(t *testing.T) {
	supply := types.TokenSupply{
		TotalSupply:           types.DefaultTotalSupply,
		FounderAllocation:     types.CalculateAllocation(types.DefaultTotalSupply, types.FounderAllocationPercentage),
		TeamAllocation:        types.CalculateAllocation(types.DefaultTotalSupply, types.TeamAllocationPercentage),
		CommunityAllocation:   types.CalculateAllocation(types.DefaultTotalSupply, types.CommunityAllocationPercentage),
		DevelopmentAllocation: types.CalculateAllocation(types.DefaultTotalSupply, types.DevelopmentAllocationPercentage),
		LiquidityAllocation:   types.CalculateAllocation(types.DefaultTotalSupply, types.LiquidityAllocationPercentage),
		PublicSaleAllocation:  types.CalculateAllocation(types.DefaultTotalSupply, types.PublicSaleAllocationPercentage),
		DAOTreasuryAllocation: types.CalculateAllocation(types.DefaultTotalSupply, types.DAOTreasuryAllocationPercentage),
		CoFounderAllocation:   types.CalculateAllocation(types.DefaultTotalSupply, types.CoFounderAllocationPercentage),
		OperationsAllocation:  types.CalculateAllocation(types.DefaultTotalSupply, types.OperationsAllocationPercentage),
		AngelAllocation:       types.CalculateAllocation(types.DefaultTotalSupply, types.AngelAllocationPercentage),
	}

	// Validate that all allocations sum to total supply
	totalAllocation := sdk.ZeroInt()
	allocationAmounts := []string{
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

	for _, allocation := range allocationAmounts {
		amount, ok := sdk.NewIntFromString(allocation)
		require.True(t, ok, "Invalid allocation amount: %s", allocation)
		totalAllocation = totalAllocation.Add(amount)
	}

	totalSupplyInt, ok := sdk.NewIntFromString(supply.TotalSupply)
	require.True(t, ok)
	require.True(t, totalSupplyInt.Equal(totalAllocation), 
		"Total supply (%s) should equal sum of allocations (%s)", 
		totalSupplyInt, totalAllocation)
}

func TestCalculateAllocation(t *testing.T) {
	totalSupply := "1000000000000000000" // 1B tokens with 12 decimals
	
	testCases := []struct {
		name       string
		percentage string
		expected   string
	}{
		{
			name:       "founder allocation 8%",
			percentage: "8.0",
			expected:   "80000000000000000", // 80M tokens
		},
		{
			name:       "team allocation 12%",
			percentage: "12.0",
			expected:   "120000000000000000", // 120M tokens
		},
		{
			name:       "public sale allocation 20%",
			percentage: "20.0",
			expected:   "200000000000000000", // 200M tokens
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := types.CalculateAllocation(totalSupply, tc.percentage)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestVestingCategories(t *testing.T) {
	// Test that all vesting categories are defined
	categories := []string{
		types.VestingCategoryFounder,
		types.VestingCategoryTeam,
		types.VestingCategoryCommunity,
		types.VestingCategoryDevelopment,
		types.VestingCategoryDAO,
		types.VestingCategoryCoFounder,
		types.VestingCategoryOperations,
		types.VestingCategoryAngel,
	}

	for _, category := range categories {
		require.NotEmpty(t, category)
	}
}

func TestTokenDenomination(t *testing.T) {
	require.Equal(t, "namo", types.TokenDenom)
	require.Equal(t, int32(12), types.TokenDecimals)
}

func TestModuleAccounts(t *testing.T) {
	// Test that module account names are defined
	require.NotEmpty(t, types.ModuleName)
	require.NotEmpty(t, types.VestingPoolName)
	require.NotEmpty(t, types.BurnPoolName)
	
	// Test that module account names are unique
	accounts := []string{
		types.ModuleName,
		types.VestingPoolName,
		types.BurnPoolName,
	}
	
	accountSet := make(map[string]bool)
	for _, account := range accounts {
		require.False(t, accountSet[account], "Duplicate module account name: %s", account)
		accountSet[account] = true
	}
}