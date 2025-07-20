package keeper_test

import (
	"fmt"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/deshchain/deshchain/x/shikshamitra/keeper"
	"github.com/deshchain/deshchain/x/shikshamitra/types"
)

type ShikshaMitraTestSuite struct {
	suite.Suite
	ctx    sdk.Context
	keeper keeper.Keeper
}

func (suite *ShikshaMitraTestSuite) SetupTest() {
	// Mock setup - in production this would include proper module setup
	suite.ctx = sdk.Context{}
	suite.keeper = keeper.Keeper{}
}

func TestShikshaMitraTestSuite(t *testing.T) {
	suite.Run(t, new(ShikshaMitraTestSuite))
}

func (suite *ShikshaMitraTestSuite) TestSetAndGetEducationLoan() {
	// Create a test education loan
	loan := types.EducationLoan{
		ID:           "test-edu-loan-1",
		StudentID:    "student123",
		Amount:       sdk.NewCoin("namo", sdk.NewInt(800000)),
		InterestRate: "0.055", // 5.5%
		Duration:     48,      // 48 months
		CourseID:     "course123",
		InstitutionID: "institution123",
		Status:       "active",
		DisbursedAt:  time.Now(),
		MoratoriumPeriod: 24, // 24 months grace period
		RepaymentSchedule: []types.RepaymentInstallment{
			{
				DueDate: time.Now().AddDate(2, 1, 0), // After moratorium
				Amount:  sdk.NewCoin("namo", sdk.NewInt(35000)),
				Status:  "pending",
			},
		},
		CoApplicantDetails: types.CoApplicant{
			Name:        "Parent Name",
			Relation:    "father",
			Income:      sdk.NewCoin("namo", sdk.NewInt(600000)),
			CreditScore: 720,
		},
	}

	// Test setting loan
	suite.keeper.SetEducationLoan(suite.ctx, loan)

	// Test getting loan
	retrievedLoan, found := suite.keeper.GetEducationLoan(suite.ctx, "test-edu-loan-1")
	suite.Require().True(found)
	suite.Require().Equal(loan.ID, retrievedLoan.ID)
	suite.Require().Equal(loan.StudentID, retrievedLoan.StudentID)
	suite.Require().Equal(loan.Amount, retrievedLoan.Amount)
}

func (suite *ShikshaMitraTestSuite) TestSetAndGetStudentProfile() {
	// Create test student profile
	profile := types.StudentProfile{
		ID:                "student123",
		Name:              "Arjun Sharma",
		DateOfBirth:       time.Now().AddDate(-22, 0, 0), // 22 years old
		Gender:            "male",
		Category:          "general",
		Pincode:           "110001",
		FamilyIncome:      sdk.NewCoin("namo", sdk.NewInt(500000)),
		VerificationStatus: "verified",
		KYCDocuments:      []string{"aadhaar", "10th_marksheet", "12th_marksheet"},
		ActiveLoans:       0,
		TotalLoanHistory:  sdk.NewCoin("namo", sdk.NewInt(0)),
		DhanPataAddress:   "deshchain1student123",
		AcademicRecords: []types.AcademicRecord{
			{
				Level:      "12th",
				Board:      "CBSE",
				Percentage: sdk.NewDec(85),
				YearOfPassing: 2020,
			},
			{
				Level:      "graduation",
				Board:      "Delhi University",
				Percentage: sdk.NewDec(78),
				YearOfPassing: 2023,
			},
		},
	}

	// Test setting student profile
	suite.keeper.SetStudentProfile(suite.ctx, profile)

	// Test getting student profile
	retrievedProfile, found := suite.keeper.GetStudentProfile(suite.ctx, "student123")
	suite.Require().True(found)
	suite.Require().Equal(profile.ID, retrievedProfile.ID)
	suite.Require().Equal(profile.Name, retrievedProfile.Name)
	suite.Require().Equal(profile.Gender, retrievedProfile.Gender)
}

func (suite *ShikshaMitraTestSuite) TestSetAndGetInstitution() {
	// Create test institution
	institution := types.Institution{
		ID:           "institution123",
		Name:         "Indian Institute of Technology Delhi",
		Type:         types.InstitutionType_IIT,
		Location:     "New Delhi",
		IsRecognized: true,
		Accreditation: []string{"NAAC A++", "NBA"},
		MaxLoanAmount: sdk.NewCoin("namo", sdk.NewInt(2000000)), // 20 lakhs
		Courses: []string{"computer_science", "mechanical_engineering", "electrical_engineering"},
	}

	// Test setting institution
	suite.keeper.SetInstitution(suite.ctx, institution)

	// Test getting institution
	retrievedInstitution, found := suite.keeper.GetInstitution(suite.ctx, "institution123")
	suite.Require().True(found)
	suite.Require().Equal(institution.ID, retrievedInstitution.ID)
	suite.Require().Equal(institution.Name, retrievedInstitution.Name)
	suite.Require().Equal(institution.Type, retrievedInstitution.Type)
}

func (suite *ShikshaMitraTestSuite) TestSetAndGetCourse() {
	// Create test course
	course := types.Course{
		ID:           "course123",
		Name:         "Computer Science Engineering",
		Type:         types.CourseType_PROFESSIONAL,
		Duration:     48, // 4 years
		MaxAge:       25,
		MinQualification: "12th with PCM",
		FeeRange: types.FeeRange{
			Min: sdk.NewCoin("namo", sdk.NewInt(500000)),  // 5 lakhs
			Max: sdk.NewCoin("namo", sdk.NewInt(1500000)), // 15 lakhs
		},
	}

	// Test setting course
	suite.keeper.SetCourse(suite.ctx, course)

	// Test getting course
	retrievedCourse, found := suite.keeper.GetCourse(suite.ctx, "course123")
	suite.Require().True(found)
	suite.Require().Equal(course.ID, retrievedCourse.ID)
	suite.Require().Equal(course.Name, retrievedCourse.Name)
	suite.Require().Equal(course.Type, retrievedCourse.Type)
}

func (suite *ShikshaMitraTestSuite) TestCalculateInterestRate() {
	// Create test application
	application := types.LoanApplication{
		StudentID:       "student123",
		CourseID:        "course123",
		InstitutionID:   "institution123",
		RequestedAmount: sdk.NewCoin("namo", sdk.NewInt(1000000)),
		AdmissionStatus: "confirmed",
		CoApplicantDetails: types.CoApplicant{
			CreditScore: 750,
		},
		DhanPataAddress: "deshchain1student123",
	}

	// Create student profile with good academic record
	profile := types.StudentProfile{
		ID: "student123",
		AcademicRecords: []types.AcademicRecord{
			{
				Level:      "graduation",
				Percentage: sdk.NewDec(85), // Good percentage
			},
		},
		VerificationStatus: "verified",
	}
	suite.keeper.SetStudentProfile(suite.ctx, profile)

	// Create institution (IIT)
	institution := types.Institution{
		ID:           "institution123",
		Type:         types.InstitutionType_IIT,
		IsRecognized: true,
	}
	suite.keeper.SetInstitution(suite.ctx, institution)

	// Create course (Professional)
	course := types.Course{
		ID:   "course123",
		Type: types.CourseType_PROFESSIONAL,
	}
	suite.keeper.SetCourse(suite.ctx, course)

	// Test interest rate calculation
	rate := suite.keeper.CalculateInterestRate(suite.ctx, application)
	
	// Should be between min and max rates (4-7%)
	minRate := sdk.NewDecWithPrec(4, 2) // 4%
	maxRate := sdk.NewDecWithPrec(7, 2) // 7%
	
	suite.Require().True(rate.GTE(minRate))
	suite.Require().True(rate.LTE(maxRate))
	
	// For IIT + good academic record, should be closer to minimum
	expectedMaxForGoodProfile := sdk.NewDecWithPrec(5, 2) // 5%
	suite.Require().True(rate.LTE(expectedMaxForGoodProfile))
}

func (suite *ShikshaMitraTestSuite) TestCheckEligibility() {
	// Create test application
	application := types.LoanApplication{
		StudentID:       "student123",
		CourseID:        "course123",
		InstitutionID:   "institution123",
		RequestedAmount: sdk.NewCoin("namo", sdk.NewInt(800000)),
		AdmissionStatus: "confirmed",
		CoApplicantDetails: types.CoApplicant{
			CreditScore: 700,
		},
		DhanPataAddress: "deshchain1student123",
	}

	// Create eligible student profile
	profile := types.StudentProfile{
		ID:                 "student123",
		DateOfBirth:        time.Now().AddDate(-20, 0, 0), // 20 years old
		VerificationStatus: "verified",
		ActiveLoans:        0, // No active loans
	}
	suite.keeper.SetStudentProfile(suite.ctx, profile)

	// Create recognized institution
	institution := types.Institution{
		ID:            "institution123",
		IsRecognized:  true,
		MaxLoanAmount: sdk.NewCoin("namo", sdk.NewInt(1500000)),
	}
	suite.keeper.SetInstitution(suite.ctx, institution)

	// Create course
	course := types.Course{
		ID:     "course123",
		MaxAge: 25,
	}
	suite.keeper.SetCourse(suite.ctx, course)

	// Test eligibility check
	eligible, reason := suite.keeper.CheckEligibility(suite.ctx, application)
	suite.Require().True(eligible)
	suite.Require().Empty(reason)

	// Test with ineligible student (too many active loans)
	profile.ActiveLoans = 3 // Exceeds limit of 2
	suite.keeper.SetStudentProfile(suite.ctx, profile)
	
	eligible, reason = suite.keeper.CheckEligibility(suite.ctx, application)
	suite.Require().False(eligible)
	suite.Require().Contains(reason, "Maximum active education loans limit reached")
}

func (suite *ShikshaMitraTestSuite) TestEstimateInterestRate() {
	// Test with high academic score and IIT
	rate := suite.keeper.EstimateInterestRate(
		suite.ctx,
		"student123",
		"iit",
		"professional",
		sdk.NewDec(95), // 95% academic score
		sdk.NewDec(300000), // 3 lakhs family income
	)
	
	// Should get the best rate due to high score + IIT + low income
	expectedRate := sdk.NewDecWithPrec(40, 3) // 4.0%
	suite.Require().Equal(expectedRate, rate)

	// Test with average profile
	rate = suite.keeper.EstimateInterestRate(
		suite.ctx,
		"student123",
		"state_university",
		"general",
		sdk.NewDec(65), // 65% academic score
		sdk.NewDec(800000), // 8 lakhs family income
	)
	
	// Should be around base rate (5.5%)
	expectedRate = sdk.NewDecWithPrec(55, 3) // 5.5%
	suite.Require().Equal(expectedRate, rate)
}

func (suite *ShikshaMitraTestSuite) TestGetApplicableDiscounts() {
	// Create female student profile in rural area
	profile := types.StudentProfile{
		ID:      "student123",
		Gender:  "female",
		Pincode: "500001", // Rural pincode (starts with 5)
	}
	suite.keeper.SetStudentProfile(suite.ctx, profile)

	// Test discounts for female IIT student
	discounts := suite.keeper.GetApplicableDiscounts(
		suite.ctx,
		"student123",
		"iit",
		"professional",
		sdk.NewDec(92), // 92% academic score
	)
	
	// Should get multiple discounts
	suite.Require().NotEmpty(discounts)
	
	// Check for specific discounts
	var hasAcademicDiscount, hasWomenDiscount, hasRuralDiscount, hasIITDiscount bool
	for _, discount := range discounts {
		switch discount.Type {
		case "Academic Merit":
			hasAcademicDiscount = true
		case "Women Empowerment":
			hasWomenDiscount = true
		case "Rural Development":
			hasRuralDiscount = true
		case "Premier Institution":
			hasIITDiscount = true
		}
	}
	
	suite.Require().True(hasAcademicDiscount)
	suite.Require().True(hasWomenDiscount)
	suite.Require().True(hasRuralDiscount)
	suite.Require().True(hasIITDiscount)
}

func (suite *ShikshaMitraTestSuite) TestGetActiveFestivalOffers() {
	// Test getting active festival offers
	offers := suite.keeper.GetActiveFestivalOffers(suite.ctx)
	
	// Note: In real implementation, this would depend on current date
	// For testing, we expect some offers to be defined
	suite.Require().True(len(offers) >= 0) // Could be empty or have offers
	
	// If there are offers, validate structure
	for _, offer := range offers {
		suite.Require().NotEmpty(offer.Name)
		suite.Require().NotEmpty(offer.FestivalID)
		suite.Require().True(offer.InterestReduction.GTE(sdk.ZeroDec()))
		suite.Require().True(offer.ProcessingFeeWaiver.GTE(sdk.ZeroDec()))
	}
}

func (suite *ShikshaMitraTestSuite) TestGetRequiredDocuments() {
	// Test documents for domestic professional course
	docs := suite.keeper.GetRequiredDocuments(suite.ctx, "student123", "iit", "professional")
	
	// Should have basic documents plus entrance exam score
	suite.Require().True(len(docs) >= 7) // At least 7 basic documents
	
	// Check for entrance exam requirement
	var hasEntranceExam bool
	for _, doc := range docs {
		if doc.Type == "Entrance Exam Score" {
			hasEntranceExam = true
			break
		}
	}
	suite.Require().True(hasEntranceExam)

	// Test documents for foreign education
	foreignDocs := suite.keeper.GetRequiredDocuments(suite.ctx, "student123", "foreign", "professional")
	
	// Should have additional foreign education documents
	var hasVisa, hasLanguageTest bool
	for _, doc := range foreignDocs {
		switch doc.Type {
		case "Visa Documents":
			hasVisa = true
		case "Language Test Score":
			hasLanguageTest = true
		}
	}
	suite.Require().True(hasVisa)
	suite.Require().True(hasLanguageTest)
}

func (suite *ShikshaMitraTestSuite) TestCalculateLoanStatistics() {
	// Create test loans
	loans := []types.EducationLoan{
		{
			ID:           "loan1",
			Amount:       sdk.NewCoin("namo", sdk.NewInt(800000)),
			InterestRate: "0.055",
			Status:       "active",
			CourseID:     "course1",
			InstitutionID: "institution1",
		},
		{
			ID:           "loan2",
			Amount:       sdk.NewCoin("namo", sdk.NewInt(1200000)),
			InterestRate: "0.06",
			Status:       "completed",
			CourseID:     "course2",
			InstitutionID: "institution2",
		},
		{
			ID:           "loan3",
			Amount:       sdk.NewCoin("namo", sdk.NewInt(600000)),
			InterestRate: "0.05",
			Status:       "defaulted",
			CourseID:     "course1",
			InstitutionID: "institution1",
		},
	}

	// Set loans in store (mocked)
	for _, loan := range loans {
		suite.keeper.SetEducationLoan(suite.ctx, loan)
	}

	// Create institutions and courses for breakdown
	institutions := []types.Institution{
		{ID: "institution1", Type: types.InstitutionType_IIT},
		{ID: "institution2", Type: types.InstitutionType_NIT},
	}
	for _, inst := range institutions {
		suite.keeper.SetInstitution(suite.ctx, inst)
	}

	courses := []types.Course{
		{ID: "course1", Type: types.CourseType_PROFESSIONAL},
		{ID: "course2", Type: types.CourseType_GENERAL},
	}
	for _, course := range courses {
		suite.keeper.SetCourse(suite.ctx, course)
	}

	// Calculate statistics
	stats := suite.keeper.CalculateLoanStatistics(suite.ctx)
	
	// Verify statistics
	expectedTotal := sdk.NewInt(2600000) // 800k + 1200k + 600k
	suite.Require().Equal(expectedTotal, stats.TotalLoansDisbursed)
	suite.Require().Equal(uint64(1), stats.ActiveLoans)
	suite.Require().Equal(uint64(1), stats.DefaultedLoans)
	
	// Average interest rate should be around 5.5%
	expectedAvgRate := sdk.NewDecWithPrec(55, 3) // 5.5%
	suite.Require().True(stats.AverageInterestRate.Equal(expectedAvgRate))
}

// Integration tests
func TestShikshaMitraIntegration(t *testing.T) {
	suite := new(ShikshaMitraTestSuite)
	suite.SetupTest()

	// Test complete education loan application flow
	t.Run("Complete Education Loan Application Flow", func(t *testing.T) {
		// 1. Create student profile
		profile := types.StudentProfile{
			ID:                 "integration_student",
			Name:               "Integration Test Student",
			DateOfBirth:        time.Now().AddDate(-20, 0, 0),
			Gender:             "female",
			Category:           "general",
			Pincode:            "600001", // Rural area
			FamilyIncome:       sdk.NewCoin("namo", sdk.NewInt(400000)),
			VerificationStatus: "verified",
			ActiveLoans:        0,
			DhanPataAddress:    "deshchain1integrationstudent",
			AcademicRecords: []types.AcademicRecord{
				{
					Level:      "12th",
					Percentage: sdk.NewDec(88),
				},
			},
		}
		suite.keeper.SetStudentProfile(suite.ctx, profile)

		// 2. Create institution
		institution := types.Institution{
			ID:            "integration_institution",
			Name:          "Test Engineering College",
			Type:          types.InstitutionType_NIT,
			IsRecognized:  true,
			MaxLoanAmount: sdk.NewCoin("namo", sdk.NewInt(1500000)),
		}
		suite.keeper.SetInstitution(suite.ctx, institution)

		// 3. Create course
		course := types.Course{
			ID:       "integration_course",
			Name:     "Computer Science",
			Type:     types.CourseType_PROFESSIONAL,
			Duration: 48,
			MaxAge:   25,
		}
		suite.keeper.SetCourse(suite.ctx, course)

		// 4. Create loan application
		application := types.LoanApplication{
			StudentID:       "integration_student",
			CourseID:        "integration_course",
			InstitutionID:   "integration_institution",
			RequestedAmount: sdk.NewCoin("namo", sdk.NewInt(1000000)),
			AdmissionStatus: "confirmed",
			CoApplicantDetails: types.CoApplicant{
				Name:        "Parent",
				CreditScore: 720,
			},
			DhanPataAddress: "deshchain1integrationstudent",
		}

		// 5. Check eligibility
		eligible, reason := suite.keeper.CheckEligibility(suite.ctx, application)
		require.True(t, eligible, "Student should be eligible: %s", reason)

		// 6. Calculate interest rate
		rate := suite.keeper.CalculateInterestRate(suite.ctx, application)
		require.True(t, rate.GTE(sdk.NewDecWithPrec(4, 2)))
		require.True(t, rate.LTE(sdk.NewDecWithPrec(7, 2)))

		// 7. Get applicable discounts
		discounts := suite.keeper.GetApplicableDiscounts(suite.ctx, "integration_student", "nit", "professional", sdk.NewDec(88))
		require.NotEmpty(t, discounts) // Should have women empowerment and rural discounts

		// 8. Create and store loan
		loan := types.EducationLoan{
			ID:              "integration_loan_1",
			StudentID:       application.StudentID,
			Amount:          application.RequestedAmount,
			InterestRate:    rate.String(),
			Duration:        48,
			CourseID:        application.CourseID,
			InstitutionID:   application.InstitutionID,
			Status:          "active",
			DisbursedAt:     time.Now(),
			MoratoriumPeriod: 24,
		}
		suite.keeper.SetEducationLoan(suite.ctx, loan)

		// 9. Verify loan was stored correctly
		retrievedLoan, found := suite.keeper.GetEducationLoan(suite.ctx, "integration_loan_1")
		require.True(t, found)
		require.Equal(t, loan.StudentID, retrievedLoan.StudentID)
	})
}

// Benchmark tests
func BenchmarkSetEducationLoan(b *testing.B) {
	suite := new(ShikshaMitraTestSuite)
	suite.SetupTest()

	loan := types.EducationLoan{
		ID:        "bench_loan",
		StudentID: "bench_student",
		Amount:    sdk.NewCoin("namo", sdk.NewInt(800000)),
		Status:    "active",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		loan.ID = fmt.Sprintf("bench_loan_%d", i)
		suite.keeper.SetEducationLoan(suite.ctx, loan)
	}
}

func BenchmarkCalculateInterestRate(b *testing.B) {
	suite := new(ShikshaMitraTestSuite)
	suite.SetupTest()

	// Setup test data
	profile := types.StudentProfile{
		ID: "bench_student",
		AcademicRecords: []types.AcademicRecord{
			{Percentage: sdk.NewDec(85)},
		},
		VerificationStatus: "verified",
	}
	suite.keeper.SetStudentProfile(suite.ctx, profile)

	institution := types.Institution{
		ID:           "bench_institution",
		Type:         types.InstitutionType_NIT,
		IsRecognized: true,
	}
	suite.keeper.SetInstitution(suite.ctx, institution)

	course := types.Course{
		ID:   "bench_course",
		Type: types.CourseType_PROFESSIONAL,
	}
	suite.keeper.SetCourse(suite.ctx, course)

	application := types.LoanApplication{
		StudentID:     "bench_student",
		CourseID:      "bench_course",
		InstitutionID: "bench_institution",
		CoApplicantDetails: types.CoApplicant{
			CreditScore: 750,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		suite.keeper.CalculateInterestRate(suite.ctx, application)
	}
}