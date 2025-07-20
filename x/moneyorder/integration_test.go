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

package moneyorder_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/deshchain/deshchain/x/moneyorder/keeper"
	"github.com/deshchain/deshchain/x/moneyorder/types"
)

type IntegrationTestSuite struct {
	suite.Suite
	keeper     keeper.Keeper
	ctx        sdk.Context
	testAddrs  []sdk.AccAddress
}

func (s *IntegrationTestSuite) SetupSuite() {
	// Test setup would be implemented here
	// For now, this is a placeholder structure
}

func (s *IntegrationTestSuite) TestUnifiedLiquidityPoolIntegration() {
	s.Run("Pension Contribution to Unified Pool", func() {
		// Test pension contribution triggering liquidity addition
		pensionAccountID := "TEST-PENSION-001"
		contributor := s.testAddrs[0]
		contribution := sdk.NewCoin("unamo", sdk.NewInt(1000000)) // 1000 NAMO
		villageCode := "110001"

		// Create village pool first
		villagePool := types.VillagePool{
			PoolId:      1,
			VillageName: "Test Village",
			PostalCode:  villageCode,
			Verified:    true,
			Status:      "active",
		}
		s.keeper.SetVillagePool(s.ctx, villagePool)

		// Test hook integration
		hooks := s.keeper.Hooks()
		err := hooks.AfterSurakshaContribution(s.ctx, pensionAccountID, contributor, contribution, villageCode)
		s.Require().NoError(err)

		// Verify unified pool was created and liquidity added
		unifiedPool, found := s.keeper.GetUnifiedPoolByVillageId(s.ctx, 1)
		s.Require().True(found)
		s.Require().True(unifiedPool.TotalLiquidity.IsPositive())

		// Verify allocation percentages
		s.Require().Equal(uint64(1), unifiedPool.VillagePoolId)
		s.Require().True(unifiedPool.SurakshaReserve.IsPositive())
		s.Require().True(unifiedPool.DexLiquidity.IsPositive())
		s.Require().True(unifiedPool.AgriLendingPool.IsPositive())
	})

	s.Run("DEX Trading with Unified Pool Liquidity", func() {
		// Test Money Order transaction using unified pool liquidity
		poolId := uint64(1)
		tokenIn := sdk.NewCoin("unamo", sdk.NewInt(100000))
		tokenOut := sdk.NewCoin("inr", sdk.NewInt(7500)) // 1 NAMO = 0.075 INR

		// Get initial pool state
		initialPool, found := s.keeper.GetUnifiedPoolByVillageId(s.ctx, poolId)
		s.Require().True(found)
		initialDexLiquidity := initialPool.DexLiquidity

		// Simulate DEX transaction
		hooks := s.keeper.Hooks()
		err := hooks.AfterSwap(s.ctx, poolId, tokenIn, tokenOut)
		s.Require().NoError(err)

		// Verify liquidity was used appropriately
		updatedPool, found := s.keeper.GetUnifiedPoolByVillageId(s.ctx, poolId)
		s.Require().True(found)

		// Pool should track the swap
		s.Require().NotEqual(initialDexLiquidity, updatedPool.DexLiquidity)
	})

	s.Run("Agricultural Loan from Unified Pool", func() {
		// Test agricultural lending using unified pool
		poolId := uint64(1)
		loanAmount := sdk.NewCoin("unamo", sdk.NewInt(500000)) // 500 NAMO loan
		borrower := s.testAddrs[1]
		loanType := "crop_input"
		duration := uint32(6) // 6 months

		// Get initial pool state
		initialPool, found := s.keeper.GetUnifiedPoolByVillageId(s.ctx, poolId)
		s.Require().True(found)
		initialLendingPool := initialPool.AgriLendingPool

		// Test loan processing
		err := s.keeper.ProcessAgriLoan(s.ctx, poolId, loanAmount, borrower, loanType, duration)
		s.Require().NoError(err)

		// Verify lending pool was reduced
		updatedPool, found := s.keeper.GetUnifiedPoolByVillageId(s.ctx, poolId)
		s.Require().True(found)
		s.Require().True(updatedPool.AgriLendingPool.AmountOf("unamo").LT(initialLendingPool.AmountOf("unamo")))

		// Verify loan was recorded
		s.Require().True(updatedPool.ActiveLoans > 0)
		s.Require().True(updatedPool.TotalLoansValue.IsPositive())
	})

	s.Run("Loan Repayment to Unified Pool", func() {
		// Test loan repayment returning funds to pool
		loanId := uint64(1)
		repaymentAmount := sdk.NewCoin("unamo", sdk.NewInt(550000)) // 500k + 50k interest

		// Get initial pool state
		poolId := uint64(1)
		initialPool, found := s.keeper.GetUnifiedPoolByVillageId(s.ctx, poolId)
		s.Require().True(found)
		initialLendingPool := initialPool.AgriLendingPool

		// Process repayment
		err := s.keeper.ProcessLoanRepayment(s.ctx, loanId, repaymentAmount)
		s.Require().NoError(err)

		// Verify pool received repayment + interest
		updatedPool, found := s.keeper.GetUnifiedPoolByVillageId(s.ctx, poolId)
		s.Require().True(found)
		s.Require().True(updatedPool.AgriLendingPool.AmountOf("unamo").GT(initialLendingPool.AmountOf("unamo")))

		// Verify interest was recorded as revenue
		s.Require().True(updatedPool.MonthlyLendingRevenue.IsPositive())
	})

	s.Run("Pension Maturity from Unified Pool", func() {
		// Test pension maturity processing
		pensionAccountID := "TEST-PENSION-001"
		beneficiary := s.testAddrs[0]
		maturityAmount := sdk.NewCoin("unamo", sdk.NewInt(1500000)) // 1.5M NAMO (50% return)

		// Get initial pool state
		poolId := uint64(1)
		initialPool, found := s.keeper.GetUnifiedPoolByVillageId(s.ctx, poolId)
		s.Require().True(found)
		initialReserve := initialPool.SurakshaReserve

		// Process maturity
		hooks := s.keeper.Hooks()
		err := hooks.AfterSurakshaMaturity(s.ctx, pensionAccountID, beneficiary, maturityAmount)
		s.Require().NoError(err)

		// Verify pension reserve was used
		updatedPool, found := s.keeper.GetUnifiedPoolByVillageId(s.ctx, poolId)
		s.Require().True(found)
		s.Require().True(updatedPool.SurakshaReserve.AmountOf("unamo").LT(initialReserve.AmountOf("unamo")))
	})

	s.Run("Monthly Revenue Distribution", func() {
		// Test monthly revenue distribution across the pool
		hooks := s.keeper.Hooks()
		err := hooks.MonthlyRevenueDistribution(s.ctx)
		s.Require().NoError(err)

		// Verify revenue was distributed
		poolId := uint64(1)
		pool, found := s.keeper.GetUnifiedPoolByVillageId(s.ctx, poolId)
		s.Require().True(found)

		// Check that monthly revenue counters were reset
		s.Require().True(pool.MonthlyDexRevenue.IsZero())
		s.Require().True(pool.MonthlyLendingRevenue.IsZero())

		// Verify cumulative revenue tracking
		s.Require().True(pool.TotalRevenue.IsPositive())
	})
}

func (s *IntegrationTestSuite) TestVillagePoolManagement() {
	s.Run("Create Village Pool", func() {
		villagePool := types.VillagePool{
			PoolId:      2,
			VillageName: "Integration Test Village",
			PostalCode:  "110002",
			Verified:    true,
			Status:      "active",
		}

		err := s.keeper.CreateVillagePool(s.ctx, villagePool)
		s.Require().NoError(err)

		// Verify pool was created
		retrievedPool, found := s.keeper.GetVillagePool(s.ctx, 2)
		s.Require().True(found)
		s.Require().Equal(villagePool.VillageName, retrievedPool.VillageName)
		s.Require().Equal(villagePool.PostalCode, retrievedPool.PostalCode)
	})

	s.Run("Join Village Pool", func() {
		poolId := uint64(2)
		member := s.testAddrs[2]

		err := s.keeper.JoinVillagePool(s.ctx, poolId, member)
		s.Require().NoError(err)

		// Verify member was added
		pool, found := s.keeper.GetVillagePool(s.ctx, poolId)
		s.Require().True(found)
		s.Require().Contains(pool.Members, member.String())
	})

	s.Run("Village Pool Performance Metrics", func() {
		poolId := uint64(2)

		// Update performance metrics
		err := s.keeper.UpdateVillagePoolPerformance(s.ctx, poolId)
		s.Require().NoError(err)

		// Verify metrics were calculated
		pool, found := s.keeper.GetVillagePool(s.ctx, poolId)
		s.Require().True(found)
		s.Require().GreaterOrEqual(pool.TrustScore, uint32(0))
	})
}

func (s *IntegrationTestSuite) TestFixedRateExchange() {
	s.Run("Create Fixed Rate Pool", func() {
		fixedPool := types.FixedRatePool{
			PoolId:       3,
			TokenA:       "unamo",
			TokenB:       "inr",
			ExchangeRate: sdk.NewDecWithPrec(75, 3), // 0.075 INR per NAMO
			CreatedBy:    s.testAddrs[0],
			Status:       "active",
			CreatedAt:    time.Now(),
		}

		err := s.keeper.CreateFixedRatePool(s.ctx, fixedPool)
		s.Require().NoError(err)

		// Verify pool was created
		retrievedPool, found := s.keeper.GetFixedRatePool(s.ctx, 3)
		s.Require().True(found)
		s.Require().Equal(fixedPool.TokenA, retrievedPool.TokenA)
		s.Require().Equal(fixedPool.TokenB, retrievedPool.TokenB)
		s.Require().True(fixedPool.ExchangeRate.Equal(retrievedPool.ExchangeRate))
	})

	s.Run("Execute Fixed Rate Swap", func() {
		poolId := uint64(3)
		amountIn := sdk.NewCoin("unamo", sdk.NewInt(100000)) // 100 NAMO
		trader := s.testAddrs[1]

		// Execute swap
		amountOut, err := s.keeper.SwapFixedRate(s.ctx, poolId, trader, amountIn)
		s.Require().NoError(err)
		s.Require().True(amountOut.IsPositive())

		// Verify correct exchange rate was applied
		expectedOut := sdk.NewInt(7500) // 100 * 0.075 * 1000000 (micro units)
		s.Require().Equal(expectedOut, amountOut.Amount)
		s.Require().Equal("inr", amountOut.Denom)
	})
}

func (s *IntegrationTestSuite) TestCulturalIntegration() {
	s.Run("Festival Bonus Application", func() {
		// Test festival bonus during special periods
		amount := sdk.NewCoin("unamo", sdk.NewInt(100000))
		festivalPeriod := "diwali"

		bonus, err := s.keeper.ApplyFestivalBonus(s.ctx, amount, festivalPeriod)
		s.Require().NoError(err)
		s.Require().True(bonus.IsPositive())

		// Verify bonus is appropriate (e.g., 5% during Diwali)
		expectedBonus := amount.Amount.MulRaw(5).QuoRaw(100) // 5%
		s.Require().Equal(expectedBonus, bonus.Amount)
	})

	s.Run("Cultural Quote Integration", func() {
		// Test cultural quote selection for transactions
		transactionType := "money_order"
		amount := sdk.NewCoin("unamo", sdk.NewInt(50000))

		quote, category, err := s.keeper.GetContextualQuote(s.ctx, transactionType, amount)
		s.Require().NoError(err)
		s.Require().NotEmpty(quote)
		s.Require().NotEmpty(category)

		// Verify quote is appropriate for the context
		s.Require().Contains([]string{"finance", "prosperity", "community", "trust"}, category)
	})
}

func (s *IntegrationTestSuite) TestMoneyOrderReceipts() {
	s.Run("Generate Money Order Receipt", func() {
		orderId := "MO-TEST-001"
		sender := s.testAddrs[0]
		receiver := s.testAddrs[1]
		amount := sdk.NewCoin("unamo", sdk.NewInt(75000))

		receipt, err := s.keeper.GenerateReceipt(s.ctx, orderId, sender, receiver, amount)
		s.Require().NoError(err)
		s.Require().NotEmpty(receipt.ReceiptId)
		s.Require().Equal(orderId, receipt.OrderId)
		s.Require().Equal(sender, receipt.Sender)
		s.Require().Equal(receiver, receipt.Receiver)
		s.Require().Equal(amount, receipt.Amount)
	})

	s.Run("Track Receipt Status", func() {
		receiptId := "RCP-TEST-001"

		// Update receipt status
		err := s.keeper.UpdateReceiptStatus(s.ctx, receiptId, "completed")
		s.Require().NoError(err)

		// Verify status was updated
		receipt, found := s.keeper.GetReceipt(s.ctx, receiptId)
		s.Require().True(found)
		s.Require().Equal("completed", receipt.Status)
	})
}

func (s *IntegrationTestSuite) TestErrorHandling() {
	s.Run("Invalid Pool Access", func() {
		// Test accessing non-existent pool
		_, found := s.keeper.GetVillagePool(s.ctx, 999)
		s.Require().False(found)

		// Test creating pool with invalid data
		invalidPool := types.VillagePool{
			PoolId:      0, // Invalid ID
			VillageName: "",
			PostalCode:  "invalid",
			Verified:    false,
			Status:      "unknown",
		}

		err := s.keeper.CreateVillagePool(s.ctx, invalidPool)
		s.Require().Error(err)
	})

	s.Run("Insufficient Liquidity Handling", func() {
		// Test loan request exceeding available liquidity
		poolId := uint64(1)
		excessiveLoanAmount := sdk.NewCoin("unamo", sdk.NewInt(999999999)) // Very large amount
		borrower := s.testAddrs[3]

		err := s.keeper.ProcessAgriLoan(s.ctx, poolId, excessiveLoanAmount, borrower, "equipment", 12)
		s.Require().Error(err)
		s.Require().Contains(err.Error(), "insufficient")
	})

	s.Run("Invalid Hook Integration", func() {
		// Test hook calls with invalid data
		hooks := s.keeper.Hooks()

		// Invalid pension contribution
		err := hooks.AfterSurakshaContribution(s.ctx, "", sdk.AccAddress{}, sdk.Coin{}, "")
		s.Require().Error(err)

		// Invalid maturity processing
		err = hooks.AfterSurakshaMaturity(s.ctx, "", sdk.AccAddress{}, sdk.Coin{})
		s.Require().Error(err)
	})
}

func (s *IntegrationTestSuite) TestPerformanceMetrics() {
	s.Run("Pool Performance Calculation", func() {
		poolId := uint64(1)

		// Calculate performance metrics
		performance, err := s.keeper.CalculatePoolPerformance(s.ctx, poolId)
		s.Require().NoError(err)
		s.Require().GreaterOrEqual(performance.SuccessRate, sdk.ZeroDec())
		s.Require().LessOrEqual(performance.SuccessRate, sdk.OneDec())
	})

	s.Run("Revenue Tracking", func() {
		poolId := uint64(1)

		// Get revenue metrics
		revenue, err := s.keeper.GetPoolRevenue(s.ctx, poolId)
		s.Require().NoError(err)
		s.Require().GreaterOrEqual(revenue.TotalRevenue.Amount.Int64(), int64(0))
	})
}

func (s *IntegrationTestSuite) TestSeasonalFinanceIntegration() {
	s.Run("Festival Season Finance", func() {
		// Test Diwali season jewelry business financing
		poolId := uint64(1)
		loanAmount := sdk.NewCoin("unamo", sdk.NewInt(2000000)) // 2M NAMO for jewelry business
		borrower := s.testAddrs[0]
		businessType := "jewelry"
		season := "diwali"

		// Apply seasonal premium rates
		err := s.keeper.ProcessSeasonalLoan(s.ctx, poolId, loanAmount, borrower, businessType, season)
		s.Require().NoError(err)

		// Verify premium rate was applied (18-22% during festival season)
		loan, found := s.keeper.GetSeasonalLoan(s.ctx, borrower.String(), businessType)
		s.Require().True(found)
		s.Require().True(loan.InterestRate.GTE(sdk.NewDecWithPrec(18, 2))) // >= 18%
		s.Require().True(loan.InterestRate.LTE(sdk.NewDecWithPrec(22, 2))) // <= 22%
	})

	s.Run("Wedding Season Working Capital", func() {
		// Test catering business working capital during wedding season
		poolId := uint64(1)
		workingCapital := sdk.NewCoin("unamo", sdk.NewInt(1500000)) // 1.5M NAMO
		business := s.testAddrs[1]
		businessType := "catering"
		season := "wedding"

		err := s.keeper.ProcessSeasonalWorkingCapital(s.ctx, poolId, workingCapital, business, businessType, season)
		s.Require().NoError(err)

		// Verify seasonal bonus was applied
		capital, found := s.keeper.GetWorkingCapital(s.ctx, business.String())
		s.Require().True(found)
		s.Require().True(capital.SeasonalBonus.IsPositive())
	})
}

func (s *IntegrationTestSuite) TestCulturalBusinessSupport() {
	s.Run("Traditional Handicraft Financing", func() {
		// Test special rates for traditional handicraft businesses
		poolId := uint64(1)
		loanAmount := sdk.NewCoin("unamo", sdk.NewInt(500000)) // 500K NAMO
		artisan := s.testAddrs[2]
		craftType := "handloom_weaving"

		err := s.keeper.ProcessTraditionalCraftLoan(s.ctx, poolId, loanAmount, artisan, craftType)
		s.Require().NoError(err)

		// Verify cultural preservation discount was applied (5-6% rate)
		loan, found := s.keeper.GetCraftLoan(s.ctx, artisan.String())
		s.Require().True(found)
		s.Require().True(loan.InterestRate.LTE(sdk.NewDecWithPrec(6, 2))) // <= 6%
		s.Require().True(loan.CulturalPreservationBonus.IsPositive())
	})

	s.Run("Women Enterprise Support", func() {
		// Test women entrepreneurship support with rate discounts
		poolId := uint64(1)
		loanAmount := sdk.NewCoin("unamo", sdk.NewInt(750000)) // 750K NAMO
		womenEntrepreneur := s.testAddrs[3]
		businessType := "home_textile"

		err := s.keeper.ProcessWomenEnterpriseLoan(s.ctx, poolId, loanAmount, womenEntrepreneur, businessType)
		s.Require().NoError(err)

		// Verify women entrepreneurship discount (-1% rate)
		loan, found := s.keeper.GetWomenEnterpriseLoan(s.ctx, womenEntrepreneur.String())
		s.Require().True(found)
		s.Require().True(loan.WomenEntrepreneurDiscount.Equal(sdk.NewDecWithPrec(1, 2))) // 1% discount
	})
}

func (s *IntegrationTestSuite) TestAdvancedFinancialInstruments() {
	s.Run("Revenue Based Financing (RBF)", func() {
		// Test RBF for high-growth tech startup
		businessId := "TECH-STARTUP-001"
		owner := s.testAddrs[0]
		investmentAmount := sdk.NewCoin("unamo", sdk.NewInt(2000000)) // 2M NAMO investment
		revenuePercentage := sdk.NewDecWithPrec(5, 2) // 5% of monthly revenue
		capMultiple := sdk.NewDecWithPrec(25, 1) // 2.5x cap

		err := s.keeper.CreateRBFInvestment(s.ctx, businessId, owner, investmentAmount, revenuePercentage, capMultiple)
		s.Require().NoError(err)

		// Verify RBF was created with correct terms
		rbf, found := s.keeper.GetRBFInvestment(s.ctx, businessId)
		s.Require().True(found)
		s.Require().Equal(investmentAmount, rbf.InvestmentAmount)
		s.Require().Equal(revenuePercentage, rbf.RevenuePercentage)
		s.Require().Equal(capMultiple, rbf.CapMultiple)
	})

	s.Run("Supply Chain Financing", func() {
		// Test invoice financing for B2B traders
		invoiceId := "INV-001"
		supplier := s.testAddrs[1]
		buyer := s.testAddrs[2]
		invoiceAmount := sdk.NewCoin("unamo", sdk.NewInt(1000000)) // 1M NAMO invoice
		advancePercentage := sdk.NewDecWithPrec(80, 2) // 80% advance

		err := s.keeper.CreateSupplyChainFinancing(s.ctx, invoiceId, supplier, buyer, invoiceAmount, advancePercentage)
		s.Require().NoError(err)

		// Verify advance was calculated and disbursed correctly
		scf, found := s.keeper.GetSupplyChainFinancing(s.ctx, invoiceId)
		s.Require().True(found)
		expectedAdvance := invoiceAmount.Amount.ToDec().Mul(advancePercentage).TruncateInt()
		s.Require().Equal(expectedAdvance, scf.AdvanceAmount.Amount)
	})

	s.Run("Digital Gold Lending", func() {
		// Test gold-backed micro-lending
		goldLoanId := "GOLD-001"
		borrower := s.testAddrs[3]
		goldWeight := sdk.NewDecWithPrec(100, 0) // 100 grams
		goldPurity := uint32(22) // 22K gold
		ltvRatio := sdk.NewDecWithPrec(80, 2) // 80% LTV

		err := s.keeper.CreateDigitalGoldLoan(s.ctx, goldLoanId, borrower, goldWeight, goldPurity, ltvRatio)
		s.Require().NoError(err)

		// Verify loan amount calculated based on gold value and LTV
		goldLoan, found := s.keeper.GetDigitalGoldLoan(s.ctx, goldLoanId)
		s.Require().True(found)
		s.Require().Equal(goldWeight, goldLoan.GoldWeight)
		s.Require().Equal(goldPurity, goldLoan.GoldPurity)
		s.Require().True(goldLoan.LoanAmount.IsPositive())
	})
}

func (s *IntegrationTestSuite) TestCrossModuleIntegration() {
	s.Run("Urban Pension to SME Loan Pipeline", func() {
		// Test complete flow from urban pension contribution to SME loan funding
		contributorAddr := s.testAddrs[0]
		monthlyContribution := sdk.NewCoin("unamo", sdk.NewInt(2500000)) // 2500 NAMO
		cityCode := "DEL"

		// Create urban pension scheme
		scheme, err := s.keeper.CreateUrbanSurakshaScheme(s.ctx, contributorAddr, cityCode, sdk.AccAddress{})
		s.Require().NoError(err)

		// Process monthly contribution
		err = s.keeper.ProcessUrbanSurakshaContribution(s.ctx, scheme.SchemeID, contributorAddr, monthlyContribution)
		s.Require().NoError(err)

		// Verify liquidity was added to SME pool (35% allocation)
		urbanPool, found := s.keeper.GetUrbanUnifiedPool(s.ctx, scheme.UrbanPoolID)
		s.Require().True(found)
		s.Require().True(urbanPool.EducationLoanPool.IsPositive())

		// Apply for SME loan from the same pool
		smeApplicant := s.testAddrs[1]
		loanAmount := sdk.NewCoin("unamo", sdk.NewInt(500000)) // 500K NAMO
		businessType := "manufacturing"

		err = s.keeper.ProcessSMELoanFromUrbanPool(s.ctx, scheme.UrbanPoolID, smeApplicant, loanAmount, businessType)
		s.Require().NoError(err)

		// Verify loan was disbursed and pool liquidity reduced
		updatedPool, found := s.keeper.GetUrbanUnifiedPool(s.ctx, scheme.UrbanPoolID)
		s.Require().True(found)
		s.Require().True(updatedPool.EducationLoanPool.AmountOf("unamo").LT(urbanPool.EducationLoanPool.AmountOf("unamo")))
	})

	s.Run("Referral Reward Distribution", func() {
		// Test referral reward system across pension and SME signups
		referrer := s.testAddrs[0]
		referee := s.testAddrs[1]
		rewardType := "pension"

		// Process referral for pension signup
		err := s.keeper.ProcessReferralReward(s.ctx, referrer, referee, rewardType)
		s.Require().NoError(err)

		// Verify referral reward was paid
		referralCount := s.keeper.GetReferralCount(s.ctx, referrer)
		s.Require().Equal(uint32(1), referralCount)

		// Test milestone multiplier for multiple referrals
		for i := 0; i < 4; i++ {
			newReferee := s.testAddrs[i+2]
			err = s.keeper.ProcessReferralReward(s.ctx, referrer, newReferee, rewardType)
			s.Require().NoError(err)
		}

		// Verify milestone bonus was applied (5+ referrals = 1.2x multiplier)
		finalCount := s.keeper.GetReferralCount(s.ctx, referrer)
		s.Require().Equal(uint32(5), finalCount)

		reward := s.keeper.GetLatestReferralReward(s.ctx, referrer)
		s.Require().True(reward.MilestoneBonus.IsPositive())
	})
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}