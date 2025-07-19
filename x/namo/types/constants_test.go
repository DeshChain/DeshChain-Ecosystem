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

package types_test

import (
	"testing"

	"github.com/deshchain/deshchain/x/namo/types"
	"github.com/stretchr/testify/require"
	"cosmossdk.io/math"
)

func TestTokenConstants(t *testing.T) {
	tests := []struct {
		name     string
		expected math.Int
		actual   math.Int
	}{
		{
			name:     "Total Supply",
			expected: math.NewInt(1_428_627_660_000_000),
			actual:   types.TotalSupply,
		},
		{
			name:     "Founder Allocation",
			expected: math.NewInt(142_862_766_000_000),
			actual:   types.FounderAllocation,
		},
		{
			name:     "Community Allocation",
			expected: math.NewInt(214_294_149_000_000),
			actual:   types.CommunityAllocation,
		},
		{
			name:     "Development Allocation",
			expected: math.NewInt(214_294_149_000_000),
			actual:   types.DevelopmentAllocation,
		},
		{
			name:     "Liquidity Allocation",
			expected: math.NewInt(285_725_532_000_000),
			actual:   types.LiquidityAllocation,
		},
		{
			name:     "Ecosystem Allocation",
			expected: math.NewInt(142_862_766_000_000),
			actual:   types.EcosystemAllocation,
		},
		{
			name:     "Treasury Allocation",
			expected: math.NewInt(214_294_149_000_000),
			actual:   types.TreasuryAllocation,
		},
		{
			name:     "Validator Allocation",
			expected: math.NewInt(214_294_149_000_000),
			actual:   types.ValidatorAllocation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.actual, "Allocation mismatch for %s", tt.name)
		})
	}
}

func TestAllocationPercentages(t *testing.T) {
	tests := []struct {
		name       string
		allocation math.Int
		percentage uint64
	}{
		{
			name:       "Founder 10%",
			allocation: types.FounderAllocation,
			percentage: 10,
		},
		{
			name:       "Community 15%",
			allocation: types.CommunityAllocation,
			percentage: 15,
		},
		{
			name:       "Development 15%",
			allocation: types.DevelopmentAllocation,
			percentage: 15,
		},
		{
			name:       "Liquidity 20%",
			allocation: types.LiquidityAllocation,
			percentage: 20,
		},
		{
			name:       "Ecosystem 10%",
			allocation: types.EcosystemAllocation,
			percentage: 10,
		},
		{
			name:       "Treasury 15%",
			allocation: types.TreasuryAllocation,
			percentage: 15,
		},
		{
			name:       "Validator 15%",
			allocation: types.ValidatorAllocation,
			percentage: 15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := types.TotalSupply.MulRaw(int64(tt.percentage)).QuoRaw(100)
			require.Equal(t, expected, tt.allocation, "Percentage calculation mismatch for %s", tt.name)
		})
	}
}

func TestTotalAllocationSum(t *testing.T) {
	totalAllocated := types.FounderAllocation.
		Add(types.CommunityAllocation).
		Add(types.DevelopmentAllocation).
		Add(types.LiquidityAllocation).
		Add(types.EcosystemAllocation).
		Add(types.TreasuryAllocation).
		Add(types.ValidatorAllocation)

	require.Equal(t, types.TotalSupply, totalAllocated, "Total allocations must equal total supply")
}

func TestCommunityFocusedAllocations(t *testing.T) {
	// Community-focused allocations: Community + Development + Liquidity + Treasury + Validator
	communityFocused := types.CommunityAllocation.
		Add(types.DevelopmentAllocation).
		Add(types.LiquidityAllocation).
		Add(types.TreasuryAllocation).
		Add(types.ValidatorAllocation)

	// This should be 80% of total supply (excluding Founder 10% and Ecosystem 10%)
	expectedCommunityFocused := types.TotalSupply.MulRaw(80).QuoRaw(100)
	
	require.Equal(t, expectedCommunityFocused, communityFocused, 
		"Community-focused allocations should be 80% of total supply")
}

func TestDenomConstants(t *testing.T) {
	require.Equal(t, "namo", types.DefaultDenom)
	require.Equal(t, "unamo", types.BaseDenom)
	require.Equal(t, uint32(6), types.DenomExponent)
}

func TestVestingPeriods(t *testing.T) {
	tests := []struct {
		name     string
		expected uint64
		actual   uint64
	}{
		{
			name:     "Founder Vesting Period",
			expected: 126_144_000, // 48 months in seconds
			actual:   types.FounderVestingPeriod,
		},
		{
			name:     "Founder Cliff Period",
			expected: 31_536_000, // 12 months in seconds
			actual:   types.FounderCliffPeriod,
		},
		{
			name:     "Ecosystem Vesting Period",
			expected: 94_608_000, // 36 months in seconds
			actual:   types.EcosystemVestingPeriod,
		},
		{
			name:     "Treasury Vesting Period",
			expected: 157_680_000, // 60 months in seconds
			actual:   types.TreasuryVestingPeriod,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.actual, "Vesting period mismatch for %s", tt.name)
		})
	}
}

func TestTokenPrecision(t *testing.T) {
	// Test that our allocations have proper precision (multiples of 10^6)
	allocations := []math.Int{
		types.FounderAllocation,
		types.CommunityAllocation,
		types.DevelopmentAllocation,
		types.LiquidityAllocation,
		types.EcosystemAllocation,
		types.TreasuryAllocation,
		types.ValidatorAllocation,
	}

	for i, allocation := range allocations {
		remainder := allocation.Mod(math.NewInt(1_000_000))
		require.True(t, remainder.IsZero(), 
			"Allocation %d should be divisible by 1,000,000 (6 decimal places)", i)
	}
}

func TestMinimumValidatorStake(t *testing.T) {
	// 10,000 NAMO with 6 decimals
	expected := math.NewInt(10_000_000_000)
	require.Equal(t, expected, types.MinValidatorStake)
}

func TestTransparencyAllocations(t *testing.T) {
	// Verify Community and Development funds for transparency
	require.Equal(t, types.CommunityAllocation, types.DevelopmentAllocation, 
		"Community and Development allocations should be equal for balanced transparency")
	
	// Combined transparency funds should be 30% of total
	transparencyFunds := types.CommunityAllocation.Add(types.DevelopmentAllocation)
	expectedTransparency := types.TotalSupply.MulRaw(30).QuoRaw(100)
	
	require.Equal(t, expectedTransparency, transparencyFunds,
		"Combined Community and Development funds should be 30% of total supply")
}

func TestFounderProtectionLimits(t *testing.T) {
	// Founder allocation should not exceed 10%
	maxFounderAllocation := types.TotalSupply.MulRaw(10).QuoRaw(100)
	require.True(t, types.FounderAllocation.LTE(maxFounderAllocation),
		"Founder allocation should not exceed 10% of total supply")
	
	// Founder allocation should be exactly 10%
	require.Equal(t, maxFounderAllocation, types.FounderAllocation,
		"Founder allocation should be exactly 10% for community trust")
}

func TestRoyaltyConstants(t *testing.T) {
	// Test tax royalty (0.10% = 0.001)
	expectedTaxRoyalty := math.LegacyNewDecWithPrec(1, 3) // 0.001
	require.Equal(t, expectedTaxRoyalty, types.FounderTaxRoyalty)
	
	// Test platform royalty (5% = 0.05)
	expectedPlatformRoyalty := math.LegacyNewDecWithPrec(5, 2) // 0.05
	require.Equal(t, expectedPlatformRoyalty, types.FounderPlatformRoyalty)
}

func TestNGOCharityAllocation(t *testing.T) {
	// Test NGO allocation from tax (0.75% of tax = 0.3% of transaction)
	expectedNGOTaxShare := math.LegacyNewDecWithPrec(75, 2) // 0.75 of tax
	require.Equal(t, expectedNGOTaxShare, types.NGOTaxAllocation)
	
	// Test NGO allocation from platform revenue (10%)
	expectedNGOPlatformShare := math.LegacyNewDecWithPrec(10, 2) // 0.10
	require.Equal(t, expectedNGOPlatformShare, types.NGOPlatformAllocation)
}

func TestBurnMechanism(t *testing.T) {
	// Test burn allocation from tax (0.25% of tax = 0.0625% of transaction)
	expectedBurnAllocation := math.LegacyNewDecWithPrec(25, 2) // 0.25 of tax
	require.Equal(t, expectedBurnAllocation, types.BurnAllocation)
}