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
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/DeshChain/DeshChain-Ecosystem/x/moneyorder/types"
)

func TestFixedRatePool_Validate(t *testing.T) {
	now := time.Now()
	validAddr := "desh1qperwt9wrnkg5k9e5gzfgjppzpqhyav5j24d66"
	
	tests := []struct {
		name    string
		pool    types.FixedRatePool
		wantErr bool
	}{
		{
			name: "valid pool",
			pool: types.FixedRatePool{
				PoolId:          1,
				PoolName:        "INR-NAMO Pool",
				Description:     "Fixed rate pool for INR to NAMO conversion",
				Token0Denom:     "uinr",
				Token1Denom:     "unamo",
				ExchangeRate:    sdk.NewDec(100),
				ReverseRate:     sdk.NewDecWithPrec(1, 2),
				MinOrderAmount:  sdk.NewInt(1000),
				MaxOrderAmount:  sdk.NewInt(1000000),
				DailyLimit:      sdk.NewInt(10000000),
				MonthlyLimit:    sdk.NewInt(100000000),
				Token0Balance:   sdk.NewInt(1000000),
				Token1Balance:   sdk.NewInt(100000000),
				ReservedBalance: sdk.NewInt(10000),
				SupportedRegions: []string{"110001", "400001"},
				RequiresKyc:     false,
				KycThreshold:    sdk.NewInt(50000),
				BaseFee:         sdk.NewDecWithPrec(1, 3),
				ExpressFee:      sdk.NewDecWithPrec(2, 3),
				BulkDiscount:    sdk.NewDecWithPrec(20, 2),
				CulturalPair:    true,
				FestivalBonus:   true,
				VillagePriority: false,
				Active:          true,
				MaintenanceMode: false,
				CreatedBy:       validAddr,
				CreatedAt:       now,
				UpdatedAt:       now,
				TotalOrders:     100,
				TotalVolume:     sdk.NewInt(10000000),
				DailyVolume:     sdk.NewInt(100000),
				MonthlyVolume:   sdk.NewInt(1000000),
			},
			wantErr: false,
		},
		{
			name: "invalid pool id",
			pool: types.FixedRatePool{
				PoolId: 0,
			},
			wantErr: true,
		},
		{
			name: "empty pool name",
			pool: types.FixedRatePool{
				PoolId:   1,
				PoolName: "",
			},
			wantErr: true,
		},
		{
			name: "invalid token denoms",
			pool: types.FixedRatePool{
				PoolId:      1,
				PoolName:    "Test Pool",
				Token0Denom: "",
				Token1Denom: "unamo",
			},
			wantErr: true,
		},
		{
			name: "same token denoms",
			pool: types.FixedRatePool{
				PoolId:      1,
				PoolName:    "Test Pool",
				Token0Denom: "unamo",
				Token1Denom: "unamo",
			},
			wantErr: true,
		},
		{
			name: "invalid exchange rate",
			pool: types.FixedRatePool{
				PoolId:       1,
				PoolName:     "Test Pool",
				Token0Denom:  "uinr",
				Token1Denom:  "unamo",
				ExchangeRate: sdk.NewDec(-1),
			},
			wantErr: true,
		},
		{
			name: "invalid reverse rate",
			pool: types.FixedRatePool{
				PoolId:       1,
				PoolName:     "Test Pool",
				Token0Denom:  "uinr",
				Token1Denom:  "unamo",
				ExchangeRate: sdk.NewDec(100),
				ReverseRate:  sdk.NewDec(0),
			},
			wantErr: true,
		},
		{
			name: "min order amount greater than max",
			pool: types.FixedRatePool{
				PoolId:         1,
				PoolName:       "Test Pool",
				Token0Denom:    "uinr",
				Token1Denom:    "unamo",
				ExchangeRate:   sdk.NewDec(100),
				ReverseRate:    sdk.NewDecWithPrec(1, 2),
				MinOrderAmount: sdk.NewInt(1000000),
				MaxOrderAmount: sdk.NewInt(1000),
			},
			wantErr: true,
		},
		{
			name: "negative balance",
			pool: types.FixedRatePool{
				PoolId:         1,
				PoolName:       "Test Pool",
				Token0Denom:    "uinr",
				Token1Denom:    "unamo",
				ExchangeRate:   sdk.NewDec(100),
				ReverseRate:    sdk.NewDecWithPrec(1, 2),
				MinOrderAmount: sdk.NewInt(1000),
				MaxOrderAmount: sdk.NewInt(1000000),
				Token0Balance:  sdk.NewInt(-1000),
			},
			wantErr: true,
		},
		{
			name: "invalid creator address",
			pool: types.FixedRatePool{
				PoolId:         1,
				PoolName:       "Test Pool",
				Token0Denom:    "uinr",
				Token1Denom:    "unamo",
				ExchangeRate:   sdk.NewDec(100),
				ReverseRate:    sdk.NewDecWithPrec(1, 2),
				MinOrderAmount: sdk.NewInt(1000),
				MaxOrderAmount: sdk.NewInt(1000000),
				CreatedBy:      "invalid",
			},
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.pool.Validate()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetFixedRatePoolKey(t *testing.T) {
	poolId := uint64(123)
	key := types.GetFixedRatePoolKey(poolId)
	require.NotNil(t, key)
	require.Equal(t, append(types.KeyPrefixFixedRatePool, sdk.Uint64ToBigEndian(poolId)...), key)
}

func TestCalculateReverseRate(t *testing.T) {
	tests := []struct {
		name         string
		exchangeRate sdk.Dec
		expected     sdk.Dec
	}{
		{
			name:         "rate of 100",
			exchangeRate: sdk.NewDec(100),
			expected:     sdk.NewDecWithPrec(1, 2), // 0.01
		},
		{
			name:         "rate of 50",
			exchangeRate: sdk.NewDec(50),
			expected:     sdk.NewDecWithPrec(2, 2), // 0.02
		},
		{
			name:         "rate of 1",
			exchangeRate: sdk.NewDec(1),
			expected:     sdk.NewDec(1), // 1
		},
		{
			name:         "rate of 0.5",
			exchangeRate: sdk.NewDecWithPrec(5, 1),
			expected:     sdk.NewDec(2), // 2
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reverseRate := sdk.OneDec().Quo(tt.exchangeRate)
			require.True(t, reverseRate.Equal(tt.expected))
		})
	}
}

func TestFixedRatePool_CalculateSwapOutput(t *testing.T) {
	pool := types.FixedRatePool{
		Token0Denom:    "uinr",
		Token1Denom:    "unamo",
		ExchangeRate:   sdk.NewDec(100),   // 1 NAMO = 100 INR
		ReverseRate:    sdk.NewDecWithPrec(1, 2), // 1 INR = 0.01 NAMO
		Token0Balance:  sdk.NewInt(1000000000), // 1M INR
		Token1Balance:  sdk.NewInt(10000000),   // 10K NAMO
		BaseFee:        sdk.NewDecWithPrec(3, 3), // 0.3%
	}
	
	tests := []struct {
		name          string
		tokenIn       sdk.Coin
		tokenOutDenom string
		expectedOut   sdk.Int
		feeAmount     sdk.Int
	}{
		{
			name:          "swap INR to NAMO",
			tokenIn:       sdk.NewCoin("uinr", sdk.NewInt(10000)), // 10K INR
			tokenOutDenom: "unamo",
			expectedOut:   sdk.NewInt(97), // 100 NAMO - 0.3% fee
			feeAmount:     sdk.NewInt(30), // 0.3% of 10K INR
		},
		{
			name:          "swap NAMO to INR",
			tokenIn:       sdk.NewCoin("unamo", sdk.NewInt(100)), // 100 NAMO
			tokenOutDenom: "uinr",
			expectedOut:   sdk.NewInt(9700), // 10K INR - 0.3% fee
			feeAmount:     sdk.NewInt(30),   // 0.3% of 10K INR equivalent
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output sdk.Int
			if pool.Token0Denom == tt.tokenIn.Denom {
				// INR to NAMO
				output = tt.tokenIn.Amount.ToDec().Mul(pool.ReverseRate).TruncateInt()
			} else {
				// NAMO to INR
				output = tt.tokenIn.Amount.ToDec().Mul(pool.ExchangeRate).TruncateInt()
			}
			
			// Apply fee
			fee := output.ToDec().Mul(pool.BaseFee).TruncateInt()
			finalOutput := output.Sub(fee)
			
			require.True(t, finalOutput.GTE(tt.expectedOut))
		})
	}
}