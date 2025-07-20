package keeper_test

import (
	"fmt"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/deshchain/deshchain/x/krishimitra/keeper"
	"github.com/deshchain/deshchain/x/krishimitra/types"
)

type KrishiMitraTestSuite struct {
	suite.Suite
	ctx    sdk.Context
	keeper keeper.Keeper
}

func (suite *KrishiMitraTestSuite) SetupTest() {
	// Mock setup - in production this would include proper module setup
	suite.ctx = sdk.Context{}
	suite.keeper = keeper.Keeper{}
}

func TestKrishiMitraTestSuite(t *testing.T) {
	suite.Run(t, new(KrishiMitraTestSuite))
}

func (suite *KrishiMitraTestSuite) TestSetAndGetLoan() {
	// Create a test loan
	loan := types.Loan{
		ID:          "test-loan-1",
		FarmerID:    "farmer123",
		Amount:      sdk.NewCoin("namo", sdk.NewInt(100000)),
		InterestRate: "0.07", // 7%
		Duration:    12, // 12 months
		Purpose:     "Seed purchase",
		Status:      "active",
		DisbursedAt: time.Now(),
		RepaymentSchedule: []types.RepaymentInstallment{
			{
				DueDate: time.Now().AddDate(0, 1, 0),
				Amount:  sdk.NewCoin("namo", sdk.NewInt(9000)),
				Status:  "pending",
			},
		},
	}

	// Test setting loan
	suite.keeper.SetLoan(suite.ctx, loan)

	// Test getting loan
	retrievedLoan, found := suite.keeper.GetLoan(suite.ctx, "test-loan-1")
	suite.Require().True(found)
	suite.Require().Equal(loan.ID, retrievedLoan.ID)
	suite.Require().Equal(loan.FarmerID, retrievedLoan.FarmerID)
	suite.Require().Equal(loan.Amount, retrievedLoan.Amount)
}

func (suite *KrishiMitraTestSuite) TestSetAndGetFarmerProfile() {
	// Create test farmer profile
	profile := types.FarmerProfile{
		ID:               "farmer123",
		Name:             "Ram Singh",
		Pincode:          "110001",
		LandHolding:      sdk.NewDec(2), // 2 acres
		CropTypes:        []string{"wheat", "rice"},
		VerificationStatus: "verified",
		KYCDocuments:     []string{"aadhaar", "land_records"},
		CreditScore:      750,
		ActiveLoans:      1,
		TotalLoanHistory: sdk.NewCoin("namo", sdk.NewInt(500000)),
		DhanPataAddress:  "deshchain1farmer123",
	}

	// Test setting farmer profile
	suite.keeper.SetFarmerProfile(suite.ctx, profile)

	// Test getting farmer profile
	retrievedProfile, found := suite.keeper.GetFarmerProfile(suite.ctx, "farmer123")
	suite.Require().True(found)
	suite.Require().Equal(profile.ID, retrievedProfile.ID)
	suite.Require().Equal(profile.Name, retrievedProfile.Name)
	suite.Require().Equal(profile.CreditScore, retrievedProfile.CreditScore)
}

func (suite *KrishiMitraTestSuite) TestCalculateInterestRate() {
	// Create test application
	application := types.LoanApplication{
		FarmerID:        "farmer123",
		RequestedAmount: sdk.NewCoin("namo", sdk.NewInt(100000)),
		Purpose:         "Seed purchase",
		Duration:        12,
		CropType:        "wheat",
		LandSize:        sdk.NewDec(2),
		SeasonType:      "kharif",
	}

	// Create farmer profile with good credit score
	profile := types.FarmerProfile{
		ID:          "farmer123",
		CreditScore: 800,
		LandHolding: sdk.NewDec(2),
		Pincode:     "500001", // Rural pincode
	}
	suite.keeper.SetFarmerProfile(suite.ctx, profile)

	// Test interest rate calculation
	rate := suite.keeper.CalculateInterestRate(suite.ctx, application)
	
	// Should be between min and max rates (6-9%)
	minRate := sdk.NewDecWithPrec(6, 2)  // 6%
	maxRate := sdk.NewDecWithPrec(9, 2)  // 9%
	
	suite.Require().True(rate.GTE(minRate))
	suite.Require().True(rate.LTE(maxRate))
}

func (suite *KrishiMitraTestSuite) TestCheckEligibility() {
	// Create test application
	application := types.LoanApplication{
		FarmerID:        "farmer123",
		RequestedAmount: sdk.NewCoin("namo", sdk.NewInt(50000)),
		Purpose:         "Fertilizer purchase",
		Duration:        6,
		CropType:        "rice",
		LandSize:        sdk.NewDec(1),
		SeasonType:      "kharif",
		DhanPataAddress: "deshchain1farmer123",
	}

	// Create eligible farmer profile
	profile := types.FarmerProfile{
		ID:                 "farmer123",
		VerificationStatus: "verified",
		CreditScore:        700,
		LandHolding:        sdk.NewDec(1),
		ActiveLoans:        1, // Less than max of 3
	}
	suite.keeper.SetFarmerProfile(suite.ctx, profile)

	// Test eligibility check
	eligible, reason := suite.keeper.CheckEligibility(suite.ctx, application)
	suite.Require().True(eligible)
	suite.Require().Empty(reason)

	// Test with ineligible farmer (too many active loans)
	profile.ActiveLoans = 4 // Exceeds limit
	suite.keeper.SetFarmerProfile(suite.ctx, profile)
	
	eligible, reason = suite.keeper.CheckEligibility(suite.ctx, application)
	suite.Require().False(eligible)
	suite.Require().Contains(reason, "Maximum active loans limit reached")
}

func (suite *KrishiMitraTestSuite) TestGetCropRecommendations() {
	// Test weather data
	weather := types.WeatherInfo{
		Temperature: sdk.NewDec(28), // 28°C
		Humidity:    sdk.NewDec(70), // 70%
		Rainfall:    sdk.NewDec(50), // 50mm
		Season:      "kharif",
	}

	// Test crop recommendations
	recommendations := suite.keeper.GetCropRecommendations(suite.ctx, "110001", weather)
	
	// Should include suitable crops for the temperature range
	suite.Require().NotEmpty(recommendations)
	suite.Require().Contains(recommendations, "Rice")
}

func (suite *KrishiMitraTestSuite) TestIsEligibleForSubsidy() {
	// Test small farmer (should be eligible)
	smallFarmer := types.FarmerProfile{
		LandHolding: sdk.NewDecWithPrec(5, 1), // 0.5 acres
		CreditScore: 650,
	}
	
	eligible := suite.keeper.IsEligibleForSubsidy(suite.ctx, smallFarmer, "seed_subsidy")
	suite.Require().True(eligible)

	// Test large farmer (should not be eligible for certain subsidies)
	largeFarmer := types.FarmerProfile{
		LandHolding: sdk.NewDec(10), // 10 acres
		CreditScore: 800,
	}
	
	eligible = suite.keeper.IsEligibleForSubsidy(suite.ctx, largeFarmer, "seed_subsidy")
	suite.Require().False(eligible) // Large farmers not eligible for small farmer subsidies
}

func (suite *KrishiMitraTestSuite) TestCalculateLoanStatistics() {
	// Create test loans
	loans := []types.Loan{
		{
			ID:           "loan1",
			Amount:       sdk.NewCoin("namo", sdk.NewInt(100000)),
			InterestRate: "0.07",
			Status:       "active",
		},
		{
			ID:           "loan2",
			Amount:       sdk.NewCoin("namo", sdk.NewInt(200000)),
			InterestRate: "0.08",
			Status:       "completed",
		},
		{
			ID:           "loan3",
			Amount:       sdk.NewCoin("namo", sdk.NewInt(150000)),
			InterestRate: "0.06",
			Status:       "defaulted",
		},
	}

	// Set loans in store
	for _, loan := range loans {
		suite.keeper.SetLoan(suite.ctx, loan)
	}

	// Calculate statistics
	stats := suite.keeper.CalculateLoanStatistics(suite.ctx)
	
	// Verify statistics
	expectedTotal := sdk.NewInt(450000) // 100k + 200k + 150k
	suite.Require().Equal(expectedTotal, stats.TotalLoansDisbursed)
	suite.Require().Equal(uint64(1), stats.ActiveLoans)
	suite.Require().Equal(uint64(1), stats.DefaultedLoans)
	
	// Average interest rate should be around 7%
	expectedAvgRate := sdk.NewDecWithPrec(7, 2) // 7%
	suite.Require().True(stats.AverageInterestRate.Equal(expectedAvgRate))
}

// Integration tests
func TestKrishiMitraIntegration(t *testing.T) {
	suite := new(KrishiMitraTestSuite)
	suite.SetupTest()

	// Test complete loan application flow
	t.Run("Complete Loan Application Flow", func(t *testing.T) {
		// 1. Create farmer profile
		profile := types.FarmerProfile{
			ID:                 "integration_farmer",
			Name:               "Integration Test Farmer",
			Pincode:            "500001",
			LandHolding:        sdk.NewDec(3),
			CropTypes:          []string{"wheat", "cotton"},
			VerificationStatus: "verified",
			CreditScore:        750,
			ActiveLoans:        0,
			DhanPataAddress:    "deshchain1integrationfarmer",
		}
		suite.keeper.SetFarmerProfile(suite.ctx, profile)

		// 2. Create loan application
		application := types.LoanApplication{
			FarmerID:        "integration_farmer",
			RequestedAmount: sdk.NewCoin("namo", sdk.NewInt(75000)),
			Purpose:         "Equipment purchase",
			Duration:        18,
			CropType:        "wheat",
			LandSize:        sdk.NewDec(3),
			SeasonType:      "rabi",
			DhanPataAddress: "deshchain1integrationfarmer",
		}

		// 3. Check eligibility
		eligible, reason := suite.keeper.CheckEligibility(suite.ctx, application)
		require.True(t, eligible, "Farmer should be eligible: %s", reason)

		// 4. Calculate interest rate
		rate := suite.keeper.CalculateInterestRate(suite.ctx, application)
		require.True(t, rate.GTE(sdk.NewDecWithPrec(6, 2)))
		require.True(t, rate.LTE(sdk.NewDecWithPrec(9, 2)))

		// 5. Create and store loan
		loan := types.Loan{
			ID:           "integration_loan_1",
			FarmerID:     application.FarmerID,
			Amount:       application.RequestedAmount,
			InterestRate: rate.String(),
			Duration:     application.Duration,
			Purpose:      application.Purpose,
			Status:       "active",
			DisbursedAt:  time.Now(),
		}
		suite.keeper.SetLoan(suite.ctx, loan)

		// 6. Verify loan was stored correctly
		retrievedLoan, found := suite.keeper.GetLoan(suite.ctx, "integration_loan_1")
		require.True(t, found)
		require.Equal(t, loan.FarmerID, retrievedLoan.FarmerID)
	})
}

// Benchmark tests
func BenchmarkSetLoan(b *testing.B) {
	suite := new(KrishiMitraTestSuite)
	suite.SetupTest()

	loan := types.Loan{
		ID:       "bench_loan",
		FarmerID: "bench_farmer",
		Amount:   sdk.NewCoin("namo", sdk.NewInt(100000)),
		Status:   "active",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		loan.ID = fmt.Sprintf("bench_loan_%d", i)
		suite.keeper.SetLoan(suite.ctx, loan)
	}
}

func BenchmarkGetLoan(b *testing.B) {
	suite := new(KrishiMitraTestSuite)
	suite.SetupTest()

	// Setup test loan
	loan := types.Loan{
		ID:       "bench_loan",
		FarmerID: "bench_farmer",
		Amount:   sdk.NewCoin("namo", sdk.NewInt(100000)),
		Status:   "active",
	}
	suite.keeper.SetLoan(suite.ctx, loan)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		suite.keeper.GetLoan(suite.ctx, "bench_loan")
	}
}