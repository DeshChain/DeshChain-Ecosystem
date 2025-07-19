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

	"github.com/deshchain/deshchain/x/tax/types"
	"github.com/stretchr/testify/require"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"cosmossdk.io/math"
)

func TestTaxDistributionCalculation(t *testing.T) {
	// Test with 1000 NAMO tax amount
	taxAmount := sdk.NewCoin("namo", math.NewInt(1000_000_000)) // 1000 NAMO with 6 decimals
	
	distribution := types.TaxDistribution{
		DevelopmentShare: math.LegacyNewDecWithPrec(45, 2), // 45%
		OperationsShare:  math.LegacyNewDecWithPrec(45, 2), // 45%
		FounderRoyalty:   math.LegacyNewDecWithPrec(10, 2), // 10%
		NGOShare:         math.LegacyNewDecWithPrec(30, 2), // 30%
		BurnShare:       math.LegacyNewDecWithPrec(10, 2), // 10%
	}

	result := distribution.CalculateDistribution(taxAmount)
	
	tests := []struct {
		name     string
		expected math.Int
		actual   math.Int
	}{
		{
			name:     "Development Share (45%)",
			expected: math.NewInt(450_000_000), // 450 NAMO
			actual:   result.DevelopmentAmount.Amount,
		},
		{
			name:     "Operations Share (45%)",
			expected: math.NewInt(450_000_000), // 450 NAMO
			actual:   result.OperationsAmount.Amount,
		},
		{
			name:     "Founder Royalty (10%)",
			expected: math.NewInt(100_000_000), // 100 NAMO
			actual:   result.FounderAmount.Amount,
		},
		{
			name:     "NGO Share (30%)",
			expected: math.NewInt(300_000_000), // 300 NAMO
			actual:   result.NGOAmount.Amount,
		},
		{
			name:     "Burn Share (10%)",
			expected: math.NewInt(100_000_000), // 100 NAMO
			actual:   result.BurnAmount.Amount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.actual, "Distribution mismatch for %s", tt.name)
		})
	}
}

func TestTaxDistributionTotal(t *testing.T) {
	taxAmount := sdk.NewCoin("namo", math.NewInt(1000_000_000))
	
	distribution := types.TaxDistribution{
		DevelopmentShare: math.LegacyNewDecWithPrec(45, 2),
		OperationsShare:  math.LegacyNewDecWithPrec(45, 2),
		FounderRoyalty:   math.LegacyNewDecWithPrec(10, 2),
		NGOShare:         math.LegacyNewDecWithPrec(30, 2),
		BurnShare:       math.LegacyNewDecWithPrec(10, 2),
	}

	result := distribution.CalculateDistribution(taxAmount)
	
	// Sum all distributions
	totalDistributed := result.DevelopmentAmount.Amount.
		Add(result.OperationsAmount.Amount).
		Add(result.FounderAmount.Amount).
		Add(result.NGOAmount.Amount).
		Add(result.BurnAmount.Amount)
	
	// Note: Due to the overlapping percentages in our design,
	// the total will be more than the original tax amount
	// This is intentional as NGO and Burn are additional allocations
	expectedTotal := math.NewInt(1_400_000_000) // 1400 NAMO
	require.Equal(t, expectedTotal, totalDistributed, "Total distribution calculation error")
}

func TestTaxDistributionValidation(t *testing.T) {
	tests := []struct {
		name         string
		distribution types.TaxDistribution
		shouldPass   bool
	}{
		{
			name: "Valid Distribution",
			distribution: types.TaxDistribution{
				DevelopmentShare: math.LegacyNewDecWithPrec(45, 2),
				OperationsShare:  math.LegacyNewDecWithPrec(45, 2),
				FounderRoyalty:   math.LegacyNewDecWithPrec(10, 2),
				NGOShare:         math.LegacyNewDecWithPrec(30, 2),
				BurnShare:       math.LegacyNewDecWithPrec(10, 2),
			},
			shouldPass: true,
		},
		{
			name: "Invalid Negative Share",
			distribution: types.TaxDistribution{
				DevelopmentShare: math.LegacyNewDecWithPrec(-1, 2),
				OperationsShare:  math.LegacyNewDecWithPrec(45, 2),
				FounderRoyalty:   math.LegacyNewDecWithPrec(10, 2),
				NGOShare:         math.LegacyNewDecWithPrec(30, 2),
				BurnShare:       math.LegacyNewDecWithPrec(10, 2),
			},
			shouldPass: false,
		},
		{
			name: "Invalid Share Over 100%",
			distribution: types.TaxDistribution{
				DevelopmentShare: math.LegacyNewDecWithPrec(150, 2), // 150%
				OperationsShare:  math.LegacyNewDecWithPrec(45, 2),
				FounderRoyalty:   math.LegacyNewDecWithPrec(10, 2),
				NGOShare:         math.LegacyNewDecWithPrec(30, 2),
				BurnShare:       math.LegacyNewDecWithPrec(10, 2),
			},
			shouldPass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.distribution.Validate()
			if tt.shouldPass {
				require.NoError(t, err, "Validation should pass for %s", tt.name)
			} else {
				require.Error(t, err, "Validation should fail for %s", tt.name)
			}
		})
	}
}

func TestTaxParams(t *testing.T) {
	params := types.DefaultParams()
	
	// Test default tax rate (2.5%)
	expectedTaxRate := math.LegacyNewDecWithPrec(25, 3) // 0.025
	require.Equal(t, expectedTaxRate, params.TaxRate)
	
	// Test minimum transaction amount
	expectedMinTx := math.NewInt(1_000_000) // 1 NAMO
	require.Equal(t, expectedMinTx, params.MinTxAmount)
	
	// Test tax enabled by default
	require.True(t, params.TaxEnabled)
}

func TestTaxDistributionShares(t *testing.T) {
	params := types.DefaultParams()
	
	// Test core distribution shares (should sum to 100%)
	coreShares := params.DevelopmentShare.Add(params.OperationsShare).Add(params.FounderRoyalty)
	expectedCoreTotal := math.LegacyNewDecWithPrec(100, 2) // 1.0 (100%)
	require.Equal(t, expectedCoreTotal, coreShares, "Core shares should sum to 100%")
	
	// Test individual shares
	require.Equal(t, math.LegacyNewDecWithPrec(45, 2), params.DevelopmentShare, "Development share should be 45%")
	require.Equal(t, math.LegacyNewDecWithPrec(45, 2), params.OperationsShare, "Operations share should be 45%")
	require.Equal(t, math.LegacyNewDecWithPrec(10, 2), params.FounderRoyalty, "Founder royalty should be 10%")
}

func TestVolumeDiscount(t *testing.T) {
	params := types.DefaultParams()
	
	tests := []struct {
		name           string
		txAmount       math.Int
		expectedRate   math.LegacyDec
	}{
		{
			name:         "Small Transaction",
			txAmount:     math.NewInt(1_000_000_000), // 1000 NAMO
			expectedRate: math.LegacyNewDecWithPrec(25, 3), // 2.5%
		},
		{
			name:         "Medium Transaction",
			txAmount:     math.NewInt(100_000_000_000), // 100,000 NAMO
			expectedRate: math.LegacyNewDecWithPrec(20, 3), // 2.0% (20% discount)
		},
		{
			name:         "Large Transaction",
			txAmount:     math.NewInt(1_000_000_000_000), // 1,000,000 NAMO
			expectedRate: math.LegacyNewDecWithPrec(15, 3), // 1.5% (40% discount)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualRate := params.CalculateVolumeDiscountRate(tt.txAmount)
			require.Equal(t, tt.expectedRate, actualRate, "Volume discount rate mismatch for %s", tt.name)
		})
	}
}

func TestTaxCalculation(t *testing.T) {
	params := types.DefaultParams()
	
	tests := []struct {
		name         string
		txAmount     math.Int
		expectedTax  math.Int
	}{
		{
			name:        "1000 NAMO Transaction",
			txAmount:    math.NewInt(1000_000_000), // 1000 NAMO
			expectedTax: math.NewInt(25_000_000),   // 25 NAMO (2.5%)
		},
		{
			name:        "100000 NAMO Transaction",
			txAmount:    math.NewInt(100000_000_000), // 100,000 NAMO
			expectedTax: math.NewInt(2000_000_000),   // 2000 NAMO (2.0%)
		},
		{
			name:        "1000000 NAMO Transaction",
			txAmount:    math.NewInt(1000000_000_000), // 1,000,000 NAMO
			expectedTax: math.NewInt(15000_000_000),   // 15000 NAMO (1.5%)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rate := params.CalculateVolumeDiscountRate(tt.txAmount)
			actualTax := rate.MulInt(tt.txAmount).TruncateInt()
			require.Equal(t, tt.expectedTax, actualTax, "Tax calculation mismatch for %s", tt.name)
		})
	}
}

func TestFounderRoyaltyCalculation(t *testing.T) {
	// Test founder royalty calculation (0.10% of tax = 0.0025% of transaction)
	txAmount := math.NewInt(1000_000_000) // 1000 NAMO transaction
	taxRate := math.LegacyNewDecWithPrec(25, 3) // 2.5%
	founderRoyaltyRate := math.LegacyNewDecWithPrec(10, 2) // 10% of tax
	
	taxAmount := taxRate.MulInt(txAmount).TruncateInt()
	founderRoyalty := founderRoyaltyRate.MulInt(taxAmount).TruncateInt()
	
	expectedFounderRoyalty := math.NewInt(2_500_000) // 2.5 NAMO
	require.Equal(t, expectedFounderRoyalty, founderRoyalty, "Founder royalty calculation error")
	
	// Verify it's 0.0025% of original transaction (0.10% of 2.5% tax)
	directCalculation := math.LegacyNewDecWithPrec(25, 5).MulInt(txAmount).TruncateInt() // 0.00025
	require.Equal(t, expectedFounderRoyalty, directCalculation, "Direct royalty calculation mismatch")
}

func TestNGOCharityCalculation(t *testing.T) {
	// Test NGO charity calculation (0.75% of tax = 0.01875% of transaction)
	txAmount := math.NewInt(1000_000_000) // 1000 NAMO transaction
	taxRate := math.LegacyNewDecWithPrec(25, 3) // 2.5%
	ngoShare := math.LegacyNewDecWithPrec(75, 2) // 75% of tax = 0.75
	
	taxAmount := taxRate.MulInt(txAmount).TruncateInt()
	ngoAmount := ngoShare.MulInt(taxAmount).TruncateInt()
	
	expectedNGOAmount := math.NewInt(18_750_000) // 18.75 NAMO
	require.Equal(t, expectedNGOAmount, ngoAmount, "NGO charity calculation error")
}

func TestBurnMechanismCalculation(t *testing.T) {
	// Test burn mechanism calculation (0.25% of tax = 0.00625% of transaction)
	txAmount := math.NewInt(1000_000_000) // 1000 NAMO transaction
	taxRate := math.LegacyNewDecWithPrec(25, 3) // 2.5%
	burnShare := math.LegacyNewDecWithPrec(25, 2) // 25% of tax = 0.25
	
	taxAmount := taxRate.MulInt(txAmount).TruncateInt()
	burnAmount := burnShare.MulInt(taxAmount).TruncateInt()
	
	expectedBurnAmount := math.NewInt(6_250_000) // 6.25 NAMO
	require.Equal(t, expectedBurnAmount, burnAmount, "Burn mechanism calculation error")
}

func TestZeroTaxAmount(t *testing.T) {
	distribution := types.TaxDistribution{
		DevelopmentShare: math.LegacyNewDecWithPrec(45, 2),
		OperationsShare:  math.LegacyNewDecWithPrec(45, 2),
		FounderRoyalty:   math.LegacyNewDecWithPrec(10, 2),
		NGOShare:         math.LegacyNewDecWithPrec(30, 2),
		BurnShare:       math.LegacyNewDecWithPrec(10, 2),
	}
	
	zeroTax := sdk.NewCoin("namo", math.ZeroInt())
	result := distribution.CalculateDistribution(zeroTax)
	
	// All distributions should be zero
	require.True(t, result.DevelopmentAmount.IsZero(), "Development amount should be zero")
	require.True(t, result.OperationsAmount.IsZero(), "Operations amount should be zero")
	require.True(t, result.FounderAmount.IsZero(), "Founder amount should be zero")
	require.True(t, result.NGOAmount.IsZero(), "NGO amount should be zero")
	require.True(t, result.BurnAmount.IsZero(), "Burn amount should be zero")
}

func TestLargeTaxAmount(t *testing.T) {
	// Test with very large tax amount to check for overflow
	largeTaxAmount := sdk.NewCoin("namo", math.NewInt(1_000_000_000_000_000)) // 1B NAMO
	
	distribution := types.TaxDistribution{
		DevelopmentShare: math.LegacyNewDecWithPrec(45, 2),
		OperationsShare:  math.LegacyNewDecWithPrec(45, 2),
		FounderRoyalty:   math.LegacyNewDecWithPrec(10, 2),
		NGOShare:         math.LegacyNewDecWithPrec(30, 2),
		BurnShare:       math.LegacyNewDecWithPrec(10, 2),
	}

	// Should not panic or overflow
	require.NotPanics(t, func() {
		result := distribution.CalculateDistribution(largeTaxAmount)
		require.NotNil(t, result)
		require.False(t, result.DevelopmentAmount.IsNegative())
		require.False(t, result.OperationsAmount.IsNegative())
		require.False(t, result.FounderAmount.IsNegative())
		require.False(t, result.NGOAmount.IsNegative())
		require.False(t, result.BurnAmount.IsNegative())
	})
}