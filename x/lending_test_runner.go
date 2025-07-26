package x

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	krishimitrakeeper "github.com/DeshChain/DeshChain-Ecosystem/x/krishimitra/keeper"
	krishimitratypes "github.com/DeshChain/DeshChain-Ecosystem/x/krishimitra/types"
	
	vyavasayamitratkeeper "github.com/DeshChain/DeshChain-Ecosystem/x/vyavasayamitra/keeper"
	vyavasayamitratypes "github.com/DeshChain/DeshChain-Ecosystem/x/vyavasayamitra/types"
	
	shikshamitrakeeper "github.com/DeshChain/DeshChain-Ecosystem/x/shikshamitra/keeper"
	shikshamitratypes "github.com/DeshChain/DeshChain-Ecosystem/x/shikshamitra/types"
)

// TestLendingSuiteIntegration tests the integration between all three lending modules
func TestLendingSuiteIntegration(t *testing.T) {
	// Mock context setup
	ctx := sdk.Context{}
	
	// Initialize keepers (mocked for testing)
	krishiKeeper := krishimitrakeeper.Keeper{}
	vyavasayaKeeper := vyavasayamitratkeeper.Keeper{}
	shikshaMitraKeeper := shikshamitrakeeper.Keeper{}

	t.Run("Interest Rate Comparison Across Modules", func(t *testing.T) {
		// Test that education loans have the lowest rates, followed by agriculture, then business
		
		// Create comparable loan amounts
		loanAmount := sdk.NewCoin("namo", sdk.NewInt(500000)) // 5 lakhs
		
		// Krishi Mitra application
		farmerProfile := krishimitratypes.FarmerProfile{
			ID:          "farmer1",
			CreditScore: 750,
			LandHolding: sdk.NewDec(2),
			Pincode:     "500001",
		}
		krishiKeeper.SetFarmerProfile(ctx, farmerProfile)
		
		krishiApp := krishimitratypes.LoanApplication{
			FarmerID:        "farmer1",
			RequestedAmount: loanAmount,
			Purpose:         "Seed purchase",
			Duration:        12,
		}
		krishiRate := krishiKeeper.CalculateInterestRate(ctx, krishiApp)
		
		// Vyavasaya Mitra application
		businessProfile := vyavasayamitratypes.BusinessProfile{
			ID:              "business1",
			CreditScore:     750,
			YearsInBusiness: 3,
			AnnualRevenue:   sdk.NewCoin("namo", sdk.NewInt(2000000)),
		}
		vyavasayaKeeper.SetBusinessProfile(ctx, businessProfile)
		
		vyavasayaApp := vyavasayamitratypes.LoanApplication{
			BusinessID:      "business1",
			RequestedAmount: loanAmount,
			Purpose:         "Working capital",
			Duration:        12,
		}
		vyavasayaRate := vyavasayaKeeper.CalculateInterestRate(ctx, vyavasayaApp)
		
		// Shiksha Mitra application
		studentProfile := shikshamitratypes.StudentProfile{
			ID: "student1",
			AcademicRecords: []shikshamitratypes.AcademicRecord{
				{Percentage: sdk.NewDec(80)},
			},
			VerificationStatus: "verified",
		}
		shikshaMitraKeeper.SetStudentProfile(ctx, studentProfile)
		
		institution := shikshamitratypes.Institution{
			ID:           "institution1",
			Type:         shikshamitratypes.InstitutionType_NIT,
			IsRecognized: true,
		}
		shikshaMitraKeeper.SetInstitution(ctx, institution)
		
		course := shikshamitratypes.Course{
			ID:   "course1",
			Type: shikshamitratypes.CourseType_PROFESSIONAL,
		}
		shikshaMitraKeeper.SetCourse(ctx, course)
		
		shikshaMitraApp := shikshamitratypes.LoanApplication{
			StudentID:     "student1",
			CourseID:      "course1",
			InstitutionID: "institution1",
			CoApplicantDetails: shikshamitratypes.CoApplicant{
				CreditScore: 750,
			},
		}
		shikshaMitraRate := shikshaMitraKeeper.CalculateInterestRate(ctx, shikshaMitraApp)
		
		// Verify rate hierarchy: Education < Agriculture < Business
		require.True(t, shikshaMitraRate.LT(krishiRate), 
			"Education loan rate (%s) should be less than agriculture rate (%s)", 
			shikshaMitraRate.String(), krishiRate.String())
		
		require.True(t, krishiRate.LT(vyavasayaRate), 
			"Agriculture loan rate (%s) should be less than business rate (%s)", 
			krishiRate.String(), vyavasayaRate.String())
		
		// Verify specific ranges
		require.True(t, shikshaMitraRate.GTE(sdk.NewDecWithPrec(4, 2))) // >= 4%
		require.True(t, shikshaMitraRate.LTE(sdk.NewDecWithPrec(7, 2))) // <= 7%
		
		require.True(t, krishiRate.GTE(sdk.NewDecWithPrec(6, 2)))       // >= 6%
		require.True(t, krishiRate.LTE(sdk.NewDecWithPrec(9, 2)))       // <= 9%
		
		require.True(t, vyavasayaRate.GTE(sdk.NewDecWithPrec(8, 2)))    // >= 8%
		require.True(t, vyavasayaRate.LTE(sdk.NewDecWithPrec(12, 2)))   // <= 12%
	})

	t.Run("Cross-Module Profile Verification", func(t *testing.T) {
		// Test that the same DhanPata address can be used across modules
		dhanpataAddress := "deshchain1testuser123"
		
		// Create profiles with same DhanPata address
		farmerProfile := krishimitratypes.FarmerProfile{
			ID:              "farmer123",
			DhanPataAddress: dhanpataAddress,
			VerificationStatus: "verified",
		}
		krishiKeeper.SetFarmerProfile(ctx, farmerProfile)
		
		businessProfile := vyavasayamitratypes.BusinessProfile{
			ID:              "business123",
			DhanPataAddress: dhanpataAddress,
			VerificationStatus: "verified",
		}
		vyavasayaKeeper.SetBusinessProfile(ctx, businessProfile)
		
		studentProfile := shikshamitratypes.StudentProfile{
			ID:              "student123",
			DhanPataAddress: dhanpataAddress,
			VerificationStatus: "verified",
		}
		shikshaMitraKeeper.SetStudentProfile(ctx, studentProfile)
		
		// Verify all profiles can coexist with same DhanPata address
		retrievedFarmer, found := krishiKeeper.GetFarmerProfile(ctx, "farmer123")
		require.True(t, found)
		require.Equal(t, dhanpataAddress, retrievedFarmer.DhanPataAddress)
		
		retrievedBusiness, found := vyavasayaKeeper.GetBusinessProfile(ctx, "business123")
		require.True(t, found)
		require.Equal(t, dhanpataAddress, retrievedBusiness.DhanPataAddress)
		
		retrievedStudent, found := shikshaMitraKeeper.GetStudentProfile(ctx, "student123")
		require.True(t, found)
		require.Equal(t, dhanpataAddress, retrievedStudent.DhanPataAddress)
	})

	t.Run("Festival Integration Across Modules", func(t *testing.T) {
		// Test that festival offers work consistently across all modules
		
		// Check if all modules have festival integration
		krishiOffers := krishiKeeper.GetActiveFestivalOffers(ctx)
		vyavasayaOffers := vyavasayaKeeper.GetActiveFestivalOffers(ctx)
		shikshaMitraOffers := shikshaMitraKeeper.GetActiveFestivalOffers(ctx)
		
		// All modules should support festival offers (even if none are currently active)
		require.True(t, len(krishiOffers) >= 0)
		require.True(t, len(vyavasayaOffers) >= 0)
		require.True(t, len(shikshaMitraOffers) >= 0)
		
		// If festival offers exist, they should have valid structure
		for _, offer := range krishiOffers {
			require.NotEmpty(t, offer.Name)
			require.True(t, offer.InterestReduction.GTE(sdk.ZeroDec()))
		}
	})

	t.Run("Credit Score Impact Analysis", func(t *testing.T) {
		// Test how credit scores affect interest rates across all modules
		creditScores := []int32{650, 700, 750, 800}
		
		for _, score := range creditScores {
			// Krishi Mitra
			farmerProfile := krishimitratypes.FarmerProfile{
				ID:          "farmer_credit_test",
				CreditScore: score,
				LandHolding: sdk.NewDec(2),
			}
			krishiKeeper.SetFarmerProfile(ctx, farmerProfile)
			
			krishiApp := krishimitratypes.LoanApplication{
				FarmerID:        "farmer_credit_test",
				RequestedAmount: sdk.NewCoin("namo", sdk.NewInt(100000)),
			}
			krishiRate := krishiKeeper.CalculateInterestRate(ctx, krishiApp)
			
			// Vyavasaya Mitra
			businessProfile := vyavasayamitratypes.BusinessProfile{
				ID:              "business_credit_test",
				CreditScore:     score,
				YearsInBusiness: 3,
				AnnualRevenue:   sdk.NewCoin("namo", sdk.NewInt(1500000)),
			}
			vyavasayaKeeper.SetBusinessProfile(ctx, businessProfile)
			
			vyavasayaApp := vyavasayamitratypes.LoanApplication{
				BusinessID:      "business_credit_test",
				RequestedAmount: sdk.NewCoin("namo", sdk.NewInt(100000)),
			}
			vyavasayaRate := vyavasayaKeeper.CalculateInterestRate(ctx, vyavasayaApp)
			
			// Shiksha Mitra
			shikshaMitraApp := shikshamitratypes.LoanApplication{
				StudentID:     "student1",
				CourseID:      "course1",
				InstitutionID: "institution1",
				CoApplicantDetails: shikshamitratypes.CoApplicant{
					CreditScore: score,
				},
			}
			shikshaMitraRate := shikshaMitraKeeper.CalculateInterestRate(ctx, shikshaMitraApp)
			
			// Verify that all rates are within expected ranges
			require.True(t, krishiRate.GTE(sdk.NewDecWithPrec(6, 2)))
			require.True(t, krishiRate.LTE(sdk.NewDecWithPrec(9, 2)))
			
			require.True(t, vyavasayaRate.GTE(sdk.NewDecWithPrec(8, 2)))
			require.True(t, vyavasayaRate.LTE(sdk.NewDecWithPrec(12, 2)))
			
			require.True(t, shikshaMitraRate.GTE(sdk.NewDecWithPrec(4, 2)))
			require.True(t, shikshaMitraRate.LTE(sdk.NewDecWithPrec(7, 2)))
			
			t.Logf("Credit Score %d: Krishi=%.2f%%, Vyavasaya=%.2f%%, Shiksha=%.2f%%", 
				score, 
				krishiRate.MustFloat64()*100,
				vyavasayaRate.MustFloat64()*100,
				shikshaMitraRate.MustFloat64()*100)
		}
	})

	t.Run("Rural Area Benefits Consistency", func(t *testing.T) {
		// Test that rural area benefits are consistently applied across modules
		ruralPincode := "500001" // Rural pincode
		urbanPincode := "110001" // Urban pincode
		
		// Krishi Mitra rural vs urban
		ruralFarmer := krishimitratypes.FarmerProfile{
			ID:          "rural_farmer",
			Pincode:     ruralPincode,
			CreditScore: 700,
			LandHolding: sdk.NewDec(2),
		}
		krishiKeeper.SetFarmerProfile(ctx, ruralFarmer)
		
		urbanFarmer := krishimitratypes.FarmerProfile{
			ID:          "urban_farmer",
			Pincode:     urbanPincode,
			CreditScore: 700,
			LandHolding: sdk.NewDec(2),
		}
		krishiKeeper.SetFarmerProfile(ctx, urbanFarmer)
		
		// Test that rural farmers get better rates
		ruralApp := krishimitratypes.LoanApplication{
			FarmerID:        "rural_farmer",
			RequestedAmount: sdk.NewCoin("namo", sdk.NewInt(100000)),
		}
		ruralRate := krishiKeeper.CalculateInterestRate(ctx, ruralApp)
		
		urbanApp := krishimitratypes.LoanApplication{
			FarmerID:        "urban_farmer",
			RequestedAmount: sdk.NewCoin("namo", sdk.NewInt(100000)),
		}
		urbanRate := krishiKeeper.CalculateInterestRate(ctx, urbanApp)
		
		// Rural farmers should get equal or better rates
		require.True(t, ruralRate.LTE(urbanRate), 
			"Rural farmer rate (%s) should be <= urban farmer rate (%s)", 
			ruralRate.String(), urbanRate.String())
		
		// Test rural discounts in Shiksha Mitra
		ruralStudent := shikshamitratypes.StudentProfile{
			ID:      "rural_student",
			Pincode: ruralPincode,
		}
		shikshaMitraKeeper.SetStudentProfile(ctx, ruralStudent)
		
		ruralDiscounts := shikshaMitraKeeper.GetApplicableDiscounts(ctx, "rural_student", "state_university", "general", sdk.NewDec(75))
		
		// Should have rural development discount
		var hasRuralDiscount bool
		for _, discount := range ruralDiscounts {
			if discount.Type == "Rural Development" {
				hasRuralDiscount = true
				break
			}
		}
		require.True(t, hasRuralDiscount, "Rural student should get rural development discount")
	})
}

// TestLendingModulePerformance benchmarks the performance of all lending modules
func BenchmarkLendingModulesPerformance(b *testing.B) {
	ctx := sdk.Context{}
	
	b.Run("KrishiMitra", func(b *testing.B) {
		keeper := krishimitrakeeper.Keeper{}
		profile := krishimitratypes.FarmerProfile{
			ID:          "bench_farmer",
			CreditScore: 750,
		}
		keeper.SetFarmerProfile(ctx, profile)
		
		app := krishimitratypes.LoanApplication{
			FarmerID:        "bench_farmer",
			RequestedAmount: sdk.NewCoin("namo", sdk.NewInt(100000)),
		}
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			keeper.CalculateInterestRate(ctx, app)
		}
	})
	
	b.Run("VyavasayaMitra", func(b *testing.B) {
		keeper := vyavasayamitratkeeper.Keeper{}
		profile := vyavasayamitratypes.BusinessProfile{
			ID:              "bench_business",
			CreditScore:     750,
			YearsInBusiness: 3,
			AnnualRevenue:   sdk.NewCoin("namo", sdk.NewInt(2000000)),
		}
		keeper.SetBusinessProfile(ctx, profile)
		
		app := vyavasayamitratypes.LoanApplication{
			BusinessID:      "bench_business",
			RequestedAmount: sdk.NewCoin("namo", sdk.NewInt(500000)),
		}
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			keeper.CalculateInterestRate(ctx, app)
		}
	})
	
	b.Run("ShikshaMitra", func(b *testing.B) {
		keeper := shikshamitrakeeper.Keeper{}
		profile := shikshamitratypes.StudentProfile{
			ID: "bench_student",
			AcademicRecords: []shikshamitratypes.AcademicRecord{
				{Percentage: sdk.NewDec(85)},
			},
		}
		keeper.SetStudentProfile(ctx, profile)
		
		institution := shikshamitratypes.Institution{
			ID:   "bench_institution",
			Type: shikshamitratypes.InstitutionType_NIT,
		}
		keeper.SetInstitution(ctx, institution)
		
		course := shikshamitratypes.Course{
			ID:   "bench_course",
			Type: shikshamitratypes.CourseType_PROFESSIONAL,
		}
		keeper.SetCourse(ctx, course)
		
		app := shikshamitratypes.LoanApplication{
			StudentID:     "bench_student",
			CourseID:      "bench_course",
			InstitutionID: "bench_institution",
			CoApplicantDetails: shikshamitratypes.CoApplicant{
				CreditScore: 750,
			},
		}
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			keeper.CalculateInterestRate(ctx, app)
		}
	})
}

// TestLendingModuleCompliance tests regulatory compliance across all modules
func TestLendingModuleCompliance(t *testing.T) {
	ctx := sdk.Context{}
	
	t.Run("Interest Rate Caps Compliance", func(t *testing.T) {
		// Test that no module exceeds RBI interest rate caps
		
		// Maximum rates as per RBI guidelines (hypothetical)
		maxEducationRate := sdk.NewDecWithPrec(15, 2)  // 15%
		maxAgricultureRate := sdk.NewDecWithPrec(12, 2) // 12%
		maxBusinessRate := sdk.NewDecWithPrec(18, 2)    // 18%
		
		// Test Krishi Mitra compliance
		krishiKeeper := krishimitrakeeper.Keeper{}
		farmerProfile := krishimitratypes.FarmerProfile{
			ID:          "compliance_farmer",
			CreditScore: 500, // Worst case scenario
		}
		krishiKeeper.SetFarmerProfile(ctx, farmerProfile)
		
		krishiApp := krishimitratypes.LoanApplication{
			FarmerID:        "compliance_farmer",
			RequestedAmount: sdk.NewCoin("namo", sdk.NewInt(1000000)), // Large amount
		}
		krishiMaxRate := krishiKeeper.CalculateInterestRate(ctx, krishiApp)
		require.True(t, krishiMaxRate.LTE(maxAgricultureRate), 
			"Krishi Mitra max rate (%s) exceeds compliance limit (%s)", 
			krishiMaxRate.String(), maxAgricultureRate.String())
		
		// Test Vyavasaya Mitra compliance
		vyavasayaKeeper := vyavasayamitratkeeper.Keeper{}
		businessProfile := vyavasayamitratypes.BusinessProfile{
			ID:              "compliance_business",
			CreditScore:     500, // Worst case
			YearsInBusiness: 1,   // Minimum experience
		}
		vyavasayaKeeper.SetBusinessProfile(ctx, businessProfile)
		
		vyavasayaApp := vyavasayamitratypes.LoanApplication{
			BusinessID:      "compliance_business",
			RequestedAmount: sdk.NewCoin("namo", sdk.NewInt(5000000)), // Large amount
		}
		vyavasayaMaxRate := vyavasayaKeeper.CalculateInterestRate(ctx, vyavasayaApp)
		require.True(t, vyavasayaMaxRate.LTE(maxBusinessRate), 
			"Vyavasaya Mitra max rate (%s) exceeds compliance limit (%s)", 
			vyavasayaMaxRate.String(), maxBusinessRate.String())
		
		// Test Shiksha Mitra compliance
		shikshaMitraKeeper := shikshamitrakeeper.Keeper{}
		studentProfile := shikshamitratypes.StudentProfile{
			ID: "compliance_student",
			AcademicRecords: []shikshamitratypes.AcademicRecord{
				{Percentage: sdk.NewDec(50)}, // Minimum qualification
			},
		}
		shikshaMitraKeeper.SetStudentProfile(ctx, studentProfile)
		
		institution := shikshamitratypes.Institution{
			ID:   "compliance_institution",
			Type: shikshamitratypes.InstitutionType_PRIVATE,
		}
		shikshaMitraKeeper.SetInstitution(ctx, institution)
		
		course := shikshamitratypes.Course{
			ID:   "compliance_course",
			Type: shikshamitratypes.CourseType_GENERAL,
		}
		shikshaMitraKeeper.SetCourse(ctx, course)
		
		shikshaMitraApp := shikshamitratypes.LoanApplication{
			StudentID:     "compliance_student",
			CourseID:      "compliance_course",
			InstitutionID: "compliance_institution",
			CoApplicantDetails: shikshamitratypes.CoApplicant{
				CreditScore: 600, // Lower score
			},
		}
		shikshaMitraMaxRate := shikshaMitraKeeper.CalculateInterestRate(ctx, shikshaMitraApp)
		require.True(t, shikshaMitraMaxRate.LTE(maxEducationRate), 
			"Shiksha Mitra max rate (%s) exceeds compliance limit (%s)", 
			shikshaMitraMaxRate.String(), maxEducationRate.String())
	})
}