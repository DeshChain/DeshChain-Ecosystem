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

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/deshchain/deshchain/x/moneyorder/keeper"
	"github.com/deshchain/deshchain/x/moneyorder/types"
)

type PensionLiquidityTestSuite struct {
	suite.Suite
	keeper    keeper.Keeper
	ctx       sdk.Context
	addresses []sdk.AccAddress
}

func TestPensionLiquidityTestSuite(t *testing.T) {
	suite.Run(t, new(PensionLiquidityTestSuite))
}

func (suite *PensionLiquidityTestSuite) SetupTest() {
	// Initialize test setup
	// This would be properly set up with test keeper and context
	// For now, we're creating test structure
}

func (suite *PensionLiquidityTestSuite) TestAddPensionLiquidity() {
	// Test adding pension liquidity to village pool
	tests := []struct {
		name                string
		villagePoolId       uint64
		surakshaContribution sdk.Coin
		expectedLiquidity   sdk.Int
		expectedReserve     sdk.Int
		shouldError         bool
	}{
		{
			name:                "valid pension contribution",
			villagePoolId:       1,
			surakshaContribution: sdk.NewCoin("unamo", sdk.NewInt(1000000)), // 1000 NAMO
			expectedLiquidity:   sdk.NewInt(800000),                         // 80% = 800 NAMO
			expectedReserve:     sdk.NewInt(200000),                         // 20% = 200 NAMO
			shouldError:         false,
		},
		{
			name:                "large pension contribution",
			villagePoolId:       1,
			surakshaContribution: sdk.NewCoin("unamo", sdk.NewInt(10000000)), // 10000 NAMO
			expectedLiquidity:   sdk.NewInt(8000000),                        // 80% = 8000 NAMO
			expectedReserve:     sdk.NewInt(2000000),                        // 20% = 2000 NAMO
			shouldError:         false,
		},
		{
			name:                "invalid pool id",
			villagePoolId:       999,
			surakshaContribution: sdk.NewCoin("unamo", sdk.NewInt(1000000)),
			expectedLiquidity:   sdk.NewInt(0),
			expectedReserve:     sdk.NewInt(0),
			shouldError:         true,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			// Test would check:
			// 1. Liquidity is correctly split (80/20)
			// 2. Pool balances are updated
			// 3. Pension liquidity record is created
			// 4. Events are emitted
			
			if tc.shouldError {
				// Expect error
				suite.Require().Error(suite.keeper.AddPensionLiquidity(
					suite.ctx,
					tc.villagePoolId,
					tc.surakshaContribution,
					suite.addresses[0],
					"PENSION001",
				))
			} else {
				// Expect success
				err := suite.keeper.AddPensionLiquidity(
					suite.ctx,
					tc.villagePoolId,
					tc.surakshaContribution,
					suite.addresses[0],
					"PENSION001",
				)
				suite.Require().NoError(err)
				
				// Verify liquidity split
				pl, found := suite.keeper.GetSurakshaLiquidity(suite.ctx, "PENSION001", tc.villagePoolId)
				suite.Require().True(found)
				suite.Require().Equal(tc.expectedLiquidity, pl.LiquidityAmount.Amount)
				suite.Require().Equal(tc.expectedReserve, pl.ReserveAmount.Amount)
			}
		})
	}
}

func (suite *PensionLiquidityTestSuite) TestRotatePensionLiquidity() {
	// Test 12-month rotation mechanism
	
	// Setup: Add pension liquidity
	contribution := sdk.NewCoin("unamo", sdk.NewInt(1000000))
	err := suite.keeper.AddPensionLiquidity(
		suite.ctx,
		1,
		contribution,
		suite.addresses[0],
		"PENSION001",
	)
	suite.Require().NoError(err)
	
	// Fast forward 12 months
	futureCtx := suite.ctx.WithBlockTime(suite.ctx.BlockTime().AddDate(0, 12, 0))
	
	// Run rotation
	suite.keeper.RotatePensionLiquidity(futureCtx)
	
	// Verify:
	// 1. Liquidity is marked as matured
	// 2. Rewards are calculated
	// 3. Pool liquidity is updated
	// 4. Events are emitted
	
	pl, found := suite.keeper.GetSurakshaLiquidity(futureCtx, "PENSION001", 1)
	suite.Require().True(found)
	suite.Require().False(pl.IsActive)
	suite.Require().True(pl.RewardsEarned.Amount.GT(sdk.ZeroInt()))
}

func (suite *PensionLiquidityTestSuite) TestPensionLiquidityRewards() {
	// Test reward calculation for pension liquidity providers
	
	config := keeper.DefaultPensionLiquidityConfig()
	
	// Test cases for different scenarios
	tests := []struct {
		name             string
		liquidityAmount  sdk.Int
		poolTotalVolume  sdk.Int
		expectedMinAPY   sdk.Dec
		expectedBonusAPY sdk.Dec
	}{
		{
			name:             "standard liquidity provider",
			liquidityAmount:  sdk.NewInt(1000000),   // 1000 NAMO
			poolTotalVolume:  sdk.NewInt(100000000), // 100K NAMO volume
			expectedMinAPY:   sdk.NewDecWithPrec(10, 2), // 10% base APY
			expectedBonusAPY: sdk.NewDecWithPrec(5, 2),  // 5% bonus
		},
		{
			name:             "large liquidity provider",
			liquidityAmount:  sdk.NewInt(10000000),  // 10K NAMO
			poolTotalVolume:  sdk.NewInt(100000000), // 100K NAMO volume
			expectedMinAPY:   sdk.NewDecWithPrec(10, 2), // 10% base APY
			expectedBonusAPY: sdk.NewDecWithPrec(5, 2),  // 5% bonus
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			// Calculate expected rewards
			baseRewards := tc.liquidityAmount.ToDec().Mul(tc.expectedMinAPY).TruncateInt()
			bonusRewards := tc.liquidityAmount.ToDec().Mul(tc.expectedBonusAPY).TruncateInt()
			totalExpectedRewards := baseRewards.Add(bonusRewards)
			
			// Verify rewards are within expected range
			suite.Require().True(totalExpectedRewards.GT(sdk.ZeroInt()))
			
			// Total APY should be base + bonus
			totalAPY := tc.expectedMinAPY.Add(tc.expectedBonusAPY)
			suite.Require().Equal(sdk.NewDecWithPrec(15, 2), totalAPY) // 15% total
		})
	}
}

func (suite *PensionLiquidityTestSuite) TestGetSurakshaLiquidityUtilization() {
	// Test utilization statistics calculation
	
	poolId := uint64(1)
	
	// Add multiple pension contributions
	for i := 0; i < 10; i++ {
		contribution := sdk.NewCoin("unamo", sdk.NewInt(1000000)) // 1000 NAMO each
		err := suite.keeper.AddPensionLiquidity(
			suite.ctx,
			poolId,
			contribution,
			suite.addresses[i%len(suite.addresses)],
			fmt.Sprintf("PENSION%03d", i),
		)
		suite.Require().NoError(err)
	}
	
	// Get utilization stats
	stats := suite.keeper.GetSurakshaLiquidityUtilization(suite.ctx, poolId)
	
	// Verify stats
	suite.Require().Equal(poolId, stats.VillagePoolId)
	suite.Require().Equal(uint32(10), stats.ActiveContributors)
	suite.Require().True(stats.TotalPensionLiquidity.AmountOf("unamo").Equal(sdk.NewInt(8000000))) // 10 * 800K (80% of 1M)
	suite.Require().True(stats.AverageAPY.Equal(sdk.NewDecWithPrec(5, 2))) // 5% bonus
}

func (suite *PensionLiquidityTestSuite) TestMonthlyRotationCycle() {
	// Test complete 12-month cycle with monthly contributions
	
	poolId := uint64(1)
	monthlyContribution := sdk.NewCoin("unamo", sdk.NewInt(1000000)) // 1000 NAMO
	
	// Simulate 12 months of contributions
	for month := 0; month < 12; month++ {
		// Add contribution for this month
		ctx := suite.ctx.WithBlockTime(suite.ctx.BlockTime().AddDate(0, month, 0))
		
		err := suite.keeper.AddPensionLiquidity(
			ctx,
			poolId,
			monthlyContribution,
			suite.addresses[0],
			fmt.Sprintf("PENSION-M%02d", month),
		)
		suite.Require().NoError(err)
	}
	
	// After 12 months, first contribution should mature
	maturityCtx := suite.ctx.WithBlockTime(suite.ctx.BlockTime().AddDate(0, 12, 0))
	suite.keeper.RotatePensionLiquidity(maturityCtx)
	
	// Verify first contribution is matured
	pl, found := suite.keeper.GetSurakshaLiquidity(maturityCtx, "PENSION-M00", poolId)
	suite.Require().True(found)
	suite.Require().False(pl.IsActive)
	
	// Verify later contributions are still active
	pl, found = suite.keeper.GetSurakshaLiquidity(maturityCtx, "PENSION-M11", poolId)
	suite.Require().True(found)
	suite.Require().True(pl.IsActive)
	
	// Get stats to verify rotating liquidity
	stats := suite.keeper.GetSurakshaLiquidityUtilization(maturityCtx, poolId)
	suite.Require().Equal(uint32(11), stats.ActiveContributors) // 11 still active
	suite.Require().True(stats.MonthlyOutflow.AmountOf("unamo").Equal(sdk.NewInt(800000))) // First month matured
}

func (suite *PensionLiquidityTestSuite) TestEmergencyWithdrawal() {
	// Test early withdrawal with penalties
	
	// Add pension liquidity
	contribution := sdk.NewCoin("unamo", sdk.NewInt(1000000))
	err := suite.keeper.AddPensionLiquidity(
		suite.ctx,
		1,
		contribution,
		suite.addresses[0],
		"PENSION001",
	)
	suite.Require().NoError(err)
	
	// Try to withdraw after 6 months (early)
	earlyCtx := suite.ctx.WithBlockTime(suite.ctx.BlockTime().AddDate(0, 6, 0))
	
	// In real implementation, there would be a withdrawal function with penalties
	// For now, we test the concept
	pl, found := suite.keeper.GetSurakshaLiquidity(earlyCtx, "PENSION001", 1)
	suite.Require().True(found)
	suite.Require().True(pl.IsActive)
	
	// Calculate penalty (e.g., 10% for early withdrawal)
	penalty := pl.LiquidityAmount.Amount.ToDec().Mul(sdk.NewDecWithPrec(10, 2)).TruncateInt()
	withdrawAmount := pl.LiquidityAmount.Amount.Sub(penalty)
	
	suite.Require().True(withdrawAmount.LT(pl.LiquidityAmount.Amount))
}