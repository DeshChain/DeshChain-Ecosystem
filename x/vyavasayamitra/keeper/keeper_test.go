package keeper_test

import (
	"fmt"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/deshchain/deshchain/x/vyavasayamitra/keeper"
	"github.com/deshchain/deshchain/x/vyavasayamitra/types"
)

type VyavasayaMitraTestSuite struct {
	suite.Suite
	ctx    sdk.Context
	keeper keeper.Keeper
}

func (suite *VyavasayaMitraTestSuite) SetupTest() {
	// Mock setup - in production this would include proper module setup
	suite.ctx = sdk.Context{}
	suite.keeper = keeper.Keeper{}
}

func TestVyavasayaMitraTestSuite(t *testing.T) {
	suite.Run(t, new(VyavasayaMitraTestSuite))
}

func (suite *VyavasayaMitraTestSuite) TestSetAndGetBusinessLoan() {
	// Create a test business loan
	loan := types.BusinessLoan{
		ID:           "test-business-loan-1",
		BusinessID:   "business123",
		Amount:       sdk.NewCoin("namo", sdk.NewInt(500000)),
		InterestRate: "0.10", // 10%
		Duration:     24,     // 24 months
		Purpose:      "Working capital",
		LoanType:     "term_loan",
		Status:       "active",
		DisbursedAt:  time.Now(),
		RepaymentSchedule: []types.RepaymentInstallment{
			{
				DueDate: time.Now().AddDate(0, 1, 0),
				Amount:  sdk.NewCoin("namo", sdk.NewInt(25000)),
				Status:  "pending",
			},
		},
		Collateral: []types.Collateral{
			{
				Type:  "property",
				Value: sdk.NewCoin("namo", sdk.NewInt(1000000)),
			},
		},
	}

	// Test setting loan
	suite.keeper.SetBusinessLoan(suite.ctx, loan)

	// Test getting loan
	retrievedLoan, found := suite.keeper.GetBusinessLoan(suite.ctx, "test-business-loan-1")
	suite.Require().True(found)
	suite.Require().Equal(loan.ID, retrievedLoan.ID)
	suite.Require().Equal(loan.BusinessID, retrievedLoan.BusinessID)
	suite.Require().Equal(loan.Amount, retrievedLoan.Amount)
}

func (suite *VyavasayaMitraTestSuite) TestSetAndGetBusinessProfile() {
	// Create test business profile
	profile := types.BusinessProfile{
		ID:                "business123",
		BusinessName:      "Ram Enterprises",
		BusinessType:      "retail",
		Industry:          "textiles",
		RegistrationNumber: "REG123456789",
		Pincode:           "400001",
		YearsInBusiness:   5,
		AnnualRevenue:     sdk.NewCoin("namo", sdk.NewInt(2000000)),
		EmployeeCount:     25,
		VerificationStatus: "verified",
		KYCDocuments:      []string{"gst_certificate", "bank_statements", "income_tax_returns"},
		CreditScore:       720,
		ActiveLoans:       2,
		TotalLoanHistory:  sdk.NewCoin("namo", sdk.NewInt(1500000)),
		DhanPataAddress:   "deshchain1business123",
		FinancialDocuments: []types.FinancialDocument{
			{
				Type:      "bank_statement",
				Period:    "last_6_months",
				AvgBalance: sdk.NewCoin("namo", sdk.NewInt(100000)),
			},
		},
	}

	// Test setting business profile
	suite.keeper.SetBusinessProfile(suite.ctx, profile)

	// Test getting business profile
	retrievedProfile, found := suite.keeper.GetBusinessProfile(suite.ctx, "business123")
	suite.Require().True(found)
	suite.Require().Equal(profile.ID, retrievedProfile.ID)
	suite.Require().Equal(profile.BusinessName, retrievedProfile.BusinessName)
	suite.Require().Equal(profile.CreditScore, retrievedProfile.CreditScore)
}

func (suite *VyavasayaMitraTestSuite) TestCalculateInterestRate() {
	// Create test application
	application := types.LoanApplication{
		BusinessID:      "business123",
		RequestedAmount: sdk.NewCoin("namo", sdk.NewInt(500000)),
		Purpose:         "Equipment purchase",
		Duration:        18,
		LoanType:        "term_loan",
	}

	// Create business profile with good credit score
	profile := types.BusinessProfile{
		ID:              "business123",
		CreditScore:     750,
		YearsInBusiness: 5,
		AnnualRevenue:   sdk.NewCoin("namo", sdk.NewInt(2000000)),
		BusinessType:    "manufacturing",
	}
	suite.keeper.SetBusinessProfile(suite.ctx, profile)

	// Test interest rate calculation
	rate := suite.keeper.CalculateInterestRate(suite.ctx, application)
	
	// Should be between min and max rates (8-12%)
	minRate := sdk.NewDecWithPrec(8, 2)  // 8%
	maxRate := sdk.NewDecWithPrec(12, 2) // 12%
	
	suite.Require().True(rate.GTE(minRate))
	suite.Require().True(rate.LTE(maxRate))
}

func (suite *VyavasayaMitraTestSuite) TestCheckEligibility() {
	// Create test application
	application := types.LoanApplication{
		BusinessID:      "business123",
		RequestedAmount: sdk.NewCoin("namo", sdk.NewInt(300000)),
		Purpose:         "Working capital",
		Duration:        12,
		LoanType:        "working_capital",
		DhanPataAddress: "deshchain1business123",
	}

	// Create eligible business profile
	profile := types.BusinessProfile{
		ID:                 "business123",
		VerificationStatus: "verified",
		CreditScore:        700,
		YearsInBusiness:    3,
		AnnualRevenue:      sdk.NewCoin("namo", sdk.NewInt(1500000)),
		ActiveLoans:        2, // Less than max of 5
	}
	suite.keeper.SetBusinessProfile(suite.ctx, profile)

	// Test eligibility check
	eligible, reason := suite.keeper.CheckEligibility(suite.ctx, application)
	suite.Require().True(eligible)
	suite.Require().Empty(reason)

	// Test with ineligible business (too new)
	profile.YearsInBusiness = 0 // Too new
	suite.keeper.SetBusinessProfile(suite.ctx, profile)
	
	eligible, reason = suite.keeper.CheckEligibility(suite.ctx, application)
	suite.Require().False(eligible)
	suite.Require().Contains(reason, "Business must be operational for at least 1 year")
}

func (suite *VyavasayaMitraTestSuite) TestIsEligibleForCreditLine() {
	// Test eligible business
	eligibleBusiness := types.BusinessProfile{
		CreditScore:     750,
		YearsInBusiness: 3,
		AnnualRevenue:   sdk.NewCoin("namo", sdk.NewInt(1500000)), // 15 lakhs
	}
	
	annualRevenue := sdk.NewCoin("namo", sdk.NewInt(1500000))
	eligible := suite.keeper.IsEligibleForCreditLine(suite.ctx, "business123", annualRevenue)
	suite.Require().True(eligible)

	// Test ineligible business (low revenue)
	lowRevenue := sdk.NewCoin("namo", sdk.NewInt(500000)) // 5 lakhs
	eligible = suite.keeper.IsEligibleForCreditLine(suite.ctx, "business123", lowRevenue)
	suite.Require().False(eligible)
}

func (suite *VyavasayaMitraTestSuite) TestCalculateCreditUtilization() {
	// Create test active loans
	activeLoans := []types.BusinessLoan{
		{
			Amount:       sdk.NewCoin("namo", sdk.NewInt(100000)),
			OutstandingAmount: sdk.NewCoin("namo", sdk.NewInt(80000)),
		},
		{
			Amount:       sdk.NewCoin("namo", sdk.NewInt(200000)),
			OutstandingAmount: sdk.NewCoin("namo", sdk.NewInt(150000)),
		},
	}

	// Test credit utilization calculation
	utilization := suite.keeper.CalculateCreditUtilization(activeLoans)
	
	// Total outstanding: 80k + 150k = 230k
	// Total credit: 100k + 200k = 300k
	// Utilization: 230k/300k = 76.67%
	expectedUtilization := sdk.NewDecWithPrec(7667, 4) // 76.67%
	suite.Require().True(utilization.Equal(expectedUtilization))
}

func (suite *VyavasayaMitraTestSuite) TestCalculateLoanStatistics() {
	// Create test loans
	loans := []types.BusinessLoan{
		{
			ID:           "loan1",
			Amount:       sdk.NewCoin("namo", sdk.NewInt(500000)),
			InterestRate: "0.10",
			Status:       "active",
			LoanType:     "term_loan",
		},
		{
			ID:           "loan2",
			Amount:       sdk.NewCoin("namo", sdk.NewInt(300000)),
			InterestRate: "0.12",
			Status:       "completed",
			LoanType:     "working_capital",
		},
		{
			ID:           "loan3",
			Amount:       sdk.NewCoin("namo", sdk.NewInt(400000)),
			InterestRate: "0.09",
			Status:       "defaulted",
			LoanType:     "equipment_loan",
		},
	}

	// Set loans in store (mocked)
	for _, loan := range loans {
		suite.keeper.SetBusinessLoan(suite.ctx, loan)
	}

	// Calculate statistics
	stats := suite.keeper.CalculateLoanStatistics(suite.ctx)
	
	// Verify statistics
	expectedTotal := sdk.NewInt(1200000) // 500k + 300k + 400k
	suite.Require().Equal(expectedTotal, stats.TotalLoansDisbursed)
	suite.Require().Equal(uint64(1), stats.ActiveLoans)
	suite.Require().Equal(uint64(1), stats.DefaultedLoans)
	
	// Average interest rate should be around 10.33%
	expectedAvgRate := sdk.NewDecWithPrec(1033, 4) // 10.33%
	suite.Require().True(stats.AverageInterestRate.Equal(expectedAvgRate))
}

func (suite *VyavasayaMitraTestSuite) TestInvoiceFinancing() {
	// Create test invoice
	invoice := types.Invoice{
		ID:          "invoice123",
		BusinessID:  "business123",
		Amount:      sdk.NewCoin("namo", sdk.NewInt(100000)),
		DueDate:     time.Now().AddDate(0, 2, 0), // 2 months from now
		Status:      "verified",
		BuyerRating: 8, // Good buyer rating
	}

	// Test invoice verification
	isVerified := suite.keeper.VerifyInvoice(suite.ctx, invoice)
	suite.Require().True(isVerified)

	// Test financing eligibility
	eligible := suite.keeper.IsEligibleForInvoiceFinancing(suite.ctx, invoice)
	suite.Require().True(eligible)

	// Test financing amount calculation (typically 80-90% of invoice value)
	financingAmount := suite.keeper.CalculateInvoiceFinancingAmount(suite.ctx, invoice)
	expectedMin := sdk.NewInt(80000) // 80% of 100k
	expectedMax := sdk.NewInt(90000) // 90% of 100k
	
	suite.Require().True(financingAmount.GTE(expectedMin))
	suite.Require().True(financingAmount.LTE(expectedMax))
}

// Integration tests
func TestVyavasayaMitraIntegration(t *testing.T) {
	suite := new(VyavasayaMitraTestSuite)
	suite.SetupTest()

	// Test complete business loan application flow
	t.Run("Complete Business Loan Application Flow", func(t *testing.T) {
		// 1. Create business profile
		profile := types.BusinessProfile{
			ID:                 "integration_business",
			BusinessName:       "Integration Test Business",
			BusinessType:       "manufacturing",
			Industry:           "textiles",
			Pincode:            "400001",
			YearsInBusiness:    4,
			AnnualRevenue:      sdk.NewCoin("namo", sdk.NewInt(3000000)),
			EmployeeCount:      50,
			VerificationStatus: "verified",
			CreditScore:        760,
			ActiveLoans:        1,
			DhanPataAddress:    "deshchain1integrationbusiness",
		}
		suite.keeper.SetBusinessProfile(suite.ctx, profile)

		// 2. Create loan application
		application := types.LoanApplication{
			BusinessID:      "integration_business",
			RequestedAmount: sdk.NewCoin("namo", sdk.NewInt(800000)),
			Purpose:         "Equipment upgrade",
			Duration:        36,
			LoanType:        "equipment_loan",
			DhanPataAddress: "deshchain1integrationbusiness",
		}

		// 3. Check eligibility
		eligible, reason := suite.keeper.CheckEligibility(suite.ctx, application)
		require.True(t, eligible, "Business should be eligible: %s", reason)

		// 4. Calculate interest rate
		rate := suite.keeper.CalculateInterestRate(suite.ctx, application)
		require.True(t, rate.GTE(sdk.NewDecWithPrec(8, 2)))
		require.True(t, rate.LTE(sdk.NewDecWithPrec(12, 2)))

		// 5. Create and store loan
		loan := types.BusinessLoan{
			ID:           "integration_business_loan_1",
			BusinessID:   application.BusinessID,
			Amount:       application.RequestedAmount,
			InterestRate: rate.String(),
			Duration:     application.Duration,
			Purpose:      application.Purpose,
			LoanType:     application.LoanType,
			Status:       "active",
			DisbursedAt:  time.Now(),
		}
		suite.keeper.SetBusinessLoan(suite.ctx, loan)

		// 6. Verify loan was stored correctly
		retrievedLoan, found := suite.keeper.GetBusinessLoan(suite.ctx, "integration_business_loan_1")
		require.True(t, found)
		require.Equal(t, loan.BusinessID, retrievedLoan.BusinessID)
	})

	// Test credit line functionality
	t.Run("Credit Line Functionality", func(t *testing.T) {
		// Create eligible business
		profile := types.BusinessProfile{
			ID:              "credit_line_business",
			CreditScore:     780,
			YearsInBusiness: 5,
			AnnualRevenue:   sdk.NewCoin("namo", sdk.NewInt(5000000)), // 50 lakhs
		}
		suite.keeper.SetBusinessProfile(suite.ctx, profile)

		// Test credit line eligibility
		annualRevenue := sdk.NewCoin("namo", sdk.NewInt(5000000))
		eligible := suite.keeper.IsEligibleForCreditLine(suite.ctx, "credit_line_business", annualRevenue)
		require.True(t, eligible)

		// Create credit line
		creditLine := types.CreditLine{
			ID:         "credit_line_1",
			BusinessID: "credit_line_business",
			Limit:      sdk.NewCoin("namo", sdk.NewInt(1000000)), // 10 lakhs limit
			Available:  sdk.NewCoin("namo", sdk.NewInt(1000000)),
			Used:       sdk.NewCoin("namo", sdk.NewInt(0)),
			InterestRate: "0.115", // 11.5%
			Status:     "active",
		}
		suite.keeper.SetCreditLine(suite.ctx, creditLine)

		// Verify credit line
		retrievedCreditLine, found := suite.keeper.GetCreditLine(suite.ctx, "credit_line_1")
		require.True(t, found)
		require.Equal(t, creditLine.Limit, retrievedCreditLine.Limit)
	})
}

// Benchmark tests
func BenchmarkSetBusinessLoan(b *testing.B) {
	suite := new(VyavasayaMitraTestSuite)
	suite.SetupTest()

	loan := types.BusinessLoan{
		ID:         "bench_loan",
		BusinessID: "bench_business",
		Amount:     sdk.NewCoin("namo", sdk.NewInt(500000)),
		Status:     "active",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		loan.ID = fmt.Sprintf("bench_loan_%d", i)
		suite.keeper.SetBusinessLoan(suite.ctx, loan)
	}
}

func BenchmarkCalculateInterestRate(b *testing.B) {
	suite := new(VyavasayaMitraTestSuite)
	suite.SetupTest()

	// Setup test data
	profile := types.BusinessProfile{
		ID:              "bench_business",
		CreditScore:     750,
		YearsInBusiness: 5,
		AnnualRevenue:   sdk.NewCoin("namo", sdk.NewInt(2000000)),
	}
	suite.keeper.SetBusinessProfile(suite.ctx, profile)

	application := types.LoanApplication{
		BusinessID:      "bench_business",
		RequestedAmount: sdk.NewCoin("namo", sdk.NewInt(500000)),
		Duration:        24,
		LoanType:        "term_loan",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		suite.keeper.CalculateInterestRate(suite.ctx, application)
	}
}