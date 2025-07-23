package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/namo/x/shikshaamitra/types"
)

// StudentCreditAnalyzer handles comprehensive student credit analysis
type StudentCreditAnalyzer struct {
	keeper Keeper
}

// NewStudentCreditAnalyzer creates a new student credit analyzer
func NewStudentCreditAnalyzer(keeper Keeper) *StudentCreditAnalyzer {
	return &StudentCreditAnalyzer{
		keeper: keeper,
	}
}

// StudentCreditProfile represents comprehensive student credit assessment
type StudentCreditProfile struct {
	StudentID           string                      `json:"student_id"`
	CoApplicantID       string                      `json:"co_applicant_id,omitempty"`
	PersonalInfo        types.StudentPersonalInfo   `json:"personal_info"`
	AcademicProfile     types.AcademicProfile       `json:"academic_profile"`
	FamilyProfile       types.FamilyProfile         `json:"family_profile"`
	FinancialProfile    types.StudentFinancialProfile `json:"financial_profile"`
	CreditHistory       types.StudentCreditHistory  `json:"credit_history"`
	EmploymentPotential types.EmploymentPotential   `json:"employment_potential"`
	CourseAnalysis      types.CourseAnalysis        `json:"course_analysis"`
	GeographicFactors   types.GeographicFactors     `json:"geographic_factors"`
	SocialFactors       types.SocialFactors         `json:"social_factors"`
	CreditScore         int64                       `json:"credit_score"`
	RiskCategory        string                      `json:"risk_category"`
	MaxLoanEligibility  sdk.Coin                    `json:"max_loan_eligibility"`
	RecommendedRate     sdk.Dec                     `json:"recommended_rate"`
	SpecialPrograms     []string                    `json:"special_programs"`
	RiskFactors         []string                    `json:"risk_factors"`
	Recommendations     []string                    `json:"recommendations"`
	LastAssessmentDate  time.Time                   `json:"last_assessment_date"`
}

// AnalyzeStudentCredit performs comprehensive student credit analysis
func (sca *StudentCreditAnalyzer) AnalyzeStudentCredit(ctx sdk.Context, studentID, coApplicantID string) (*StudentCreditProfile, error) {
	// Get student information
	studentInfo, found := sca.keeper.GetStudentInfo(ctx, studentID)
	if !found {
		return nil, fmt.Errorf("student information not found for ID: %s", studentID)
	}

	profile := &StudentCreditProfile{
		StudentID:     studentID,
		CoApplicantID: coApplicantID,
		PersonalInfo:  studentInfo,
	}

	// Perform comprehensive analysis
	profile.AcademicProfile = sca.analyzeAcademicProfile(ctx, studentID)
	profile.FamilyProfile = sca.analyzeFamilyProfile(ctx, studentID, coApplicantID)
	profile.FinancialProfile = sca.analyzeFinancialProfile(ctx, studentID, coApplicantID)
	profile.CreditHistory = sca.analyzeCreditHistory(ctx, studentID, coApplicantID)
	profile.EmploymentPotential = sca.analyzeEmploymentPotential(ctx, studentID)
	profile.CourseAnalysis = sca.analyzeCourseFactors(ctx, studentID)
	profile.GeographicFactors = sca.analyzeGeographicFactors(ctx, studentID)
	profile.SocialFactors = sca.analyzeSocialFactors(ctx, studentID)

	// Calculate composite credit score
	profile.CreditScore = sca.calculateStudentCreditScore(profile)

	// Determine risk category
	profile.RiskCategory = sca.determineStudentRiskCategory(profile.CreditScore)

	// Calculate loan eligibility
	profile.MaxLoanEligibility = sca.calculateStudentLoanEligibility(ctx, profile)

	// Calculate recommended interest rate
	profile.RecommendedRate = sca.calculateStudentRecommendedRate(ctx, profile)

	// Identify special programs
	profile.SpecialPrograms = sca.identifyStudentSpecialPrograms(ctx, profile)

	// Identify risk factors
	profile.RiskFactors = sca.identifyStudentRiskFactors(profile)

	// Generate recommendations
	profile.Recommendations = sca.generateStudentRecommendations(profile)

	profile.LastAssessmentDate = ctx.BlockTime()

	return profile, nil
}

// analyzeAcademicProfile evaluates academic performance and potential
func (sca *StudentCreditAnalyzer) analyzeAcademicProfile(ctx sdk.Context, studentID string) types.AcademicProfile {
	academicData := sca.keeper.GetStudentAcademicData(ctx, studentID)
	
	profile := types.AcademicProfile{
		StudentID: studentID,
	}

	if academicData != nil {
		// Academic performance metrics
		profile.PreviousEducationScore = academicData.HighSchoolPercentage
		profile.EntranceExamScore = academicData.EntranceExamScore
		profile.AcademicConsistency = sca.calculateAcademicConsistency(academicData.YearlyPerformance)
		
		// Calculate academic ranking percentile
		profile.AcademicRanking = sca.calculateAcademicRanking(ctx, academicData)
		
		// Extracurricular activities
		profile.ExtracurricularScore = sca.calculateExtracurricularScore(academicData.Activities)
		
		// Academic awards and achievements
		profile.AwardsAndHonors = academicData.Awards
		profile.ResearchExperience = academicData.ResearchProjects > 0
		
		// Language proficiency (important for abroad studies)
		profile.LanguageProficiency = academicData.LanguageScores
		
		// Academic referee strength
		profile.RefereeStrength = sca.assessRefereeStrength(academicData.References)
	}

	return profile
}

// analyzeFamilyProfile evaluates family financial background
func (sca *StudentCreditAnalyzer) analyzeFamilyProfile(ctx sdk.Context, studentID, coApplicantID string) types.FamilyProfile {
	familyData := sca.keeper.GetFamilyData(ctx, studentID)
	
	profile := types.FamilyProfile{
		StudentID: studentID,
	}

	if familyData != nil {
		// Income analysis
		profile.MonthlyIncome = familyData.CombinedMonthlyIncome
		profile.IncomeStability = sca.calculateIncomeStability(familyData.IncomeHistory)
		profile.IncomeSource = familyData.PrimaryIncomeSource
		
		// Asset analysis
		profile.TotalAssets = familyData.TotalAssets
		profile.PropertyOwnership = familyData.PropertyValue.GT(sdk.ZeroDec())
		profile.InvestmentPortfolio = familyData.Investments
		
		// Existing liabilities
		profile.ExistingLoans = familyData.ExistingLoans
		profile.DebtToIncomeRatio = sca.calculateDebtToIncomeRatio(familyData.ExistingLoans, familyData.CombinedMonthlyIncome)
		
		// Family size and dependents
		profile.FamilySize = familyData.FamilySize
		profile.NumberOfDependents = familyData.Dependents
		profile.EducationExpenses = familyData.CurrentEducationExpenses
		
		// Co-applicant analysis
		if coApplicantID != "" {
			coApplicantData := sca.keeper.GetCoApplicantData(ctx, coApplicantID)
			if coApplicantData != nil {
				profile.CoApplicantIncome = coApplicantData.MonthlyIncome
				profile.CoApplicantCreditScore = coApplicantData.CreditScore
				profile.CoApplicantEmploymentStability = coApplicantData.EmploymentYears
			}
		}
	}

	return profile
}

// analyzeFinancialProfile evaluates student's own financial standing
func (sca *StudentCreditAnalyzer) analyzeFinancialProfile(ctx sdk.Context, studentID, coApplicantID string) types.StudentFinancialProfile {
	financialData := sca.keeper.GetStudentFinancialData(ctx, studentID)
	
	profile := types.StudentFinancialProfile{
		StudentID: studentID,
	}

	if financialData != nil {
		// Student's own income (if any)
		profile.StudentIncome = financialData.PartTimeIncome
		profile.Scholarships = financialData.Scholarships
		profile.Grants = financialData.Grants
		
		// Savings and financial discipline
		profile.SavingsAmount = financialData.Savings
		profile.FinancialLiteracyScore = sca.assessFinancialLiteracy(ctx, studentID)
		
		// Existing financial commitments
		profile.ExistingEducationLoans = financialData.ExistingLoans
		profile.CreditCardDebt = financialData.CreditCardDebt
		
		// Bank relationship
		profile.BankRelationshipYears = financialData.BankingHistory
		profile.AccountConductScore = sca.calculateAccountConductScore(financialData.BankStatements)
	}

	return profile
}

// analyzeCreditHistory evaluates credit behavior
func (sca *StudentCreditAnalyzer) analyzeCreditHistory(ctx sdk.Context, studentID, coApplicantID string) types.StudentCreditHistory {
	creditData := sca.keeper.GetStudentCreditData(ctx, studentID)
	
	history := types.StudentCreditHistory{
		StudentID: studentID,
	}

	if creditData != nil {
		// Credit score from external agencies
		history.CreditScore = creditData.CIBILScore
		history.CreditLength = creditData.CreditHistoryLength
		
		// Payment behavior
		history.PaymentHistory = creditData.PaymentRecords
		history.DefaultedLoans = creditData.Defaults
		history.LatePayments = creditData.LatePayments
		
		// Credit utilization
		history.CreditUtilization = sca.calculateCreditUtilization(creditData.CreditCards)
		
		// Mix of credit
		history.CreditMix = sca.assessCreditMix(creditData.CreditTypes)
	}

	// Co-applicant credit history
	if coApplicantID != "" {
		coApplicantCredit := sca.keeper.GetCoApplicantCreditData(ctx, coApplicantID)
		if coApplicantCredit != nil {
			history.CoApplicantCreditScore = coApplicantCredit.CreditScore
			history.CoApplicantPaymentHistory = coApplicantCredit.PaymentBehavior
		}
	}

	return history
}

// analyzeEmploymentPotential evaluates future earning potential
func (sca *StudentCreditAnalyzer) analyzeEmploymentPotential(ctx sdk.Context, studentID string) types.EmploymentPotential {
	courseData := sca.keeper.GetStudentCourseData(ctx, studentID)
	
	potential := types.EmploymentPotential{
		StudentID: studentID,
	}

	if courseData != nil {
		// Industry employment statistics
		industryStats := sca.keeper.GetIndustryEmploymentStats(ctx, courseData.Industry)
		if industryStats != nil {
			potential.IndustryGrowthRate = industryStats.GrowthRate
			potential.AverageSalaryRange = industryStats.SalaryRange
			potential.EmploymentRate = industryStats.EmploymentRate
		}
		
		// Institution placement records
		institutionStats := sca.keeper.GetInstitutionPlacementStats(ctx, courseData.InstitutionID)
		if institutionStats != nil {
			potential.InstitutionPlacementRate = institutionStats.PlacementRate
			potential.AverageStartingSalary = institutionStats.AverageSalary
			potential.TopRecruiters = institutionStats.TopRecruiters
		}
		
		// Course-specific factors
		potential.CourseRelevance = sca.assessCourseMarketRelevance(ctx, courseData.CourseType)
		potential.SkillDemand = sca.assessSkillDemand(ctx, courseData.Skills)
		
		// Geographic employment factors
		potential.LocationAdvantage = sca.assessLocationAdvantage(ctx, courseData.StudyLocation)
	}

	return potential
}

// analyzeCourseFactors evaluates course-specific risk and potential
func (sca *StudentCreditAnalyzer) analyzeCourseFactors(ctx sdk.Context, studentID string) types.CourseAnalysis {
	courseData := sca.keeper.GetStudentCourseData(ctx, studentID)
	
	analysis := types.CourseAnalysis{
		StudentID: studentID,
	}

	if courseData != nil {
		// Course ROI analysis
		analysis.ExpectedROI = sca.calculateCourseROI(ctx, courseData)
		analysis.PaybackPeriod = sca.calculatePaybackPeriod(ctx, courseData)
		
		// Course difficulty and completion rates
		analysis.CourseCompletionRate = sca.getCourseCompletionRate(ctx, courseData.CourseType, courseData.InstitutionID)
		analysis.CourseDifficultyLevel = sca.assessCourseDifficulty(courseData.CourseType)
		
		// Market demand and saturation
		analysis.MarketDemand = sca.assessMarketDemand(ctx, courseData.CourseType)
		analysis.MarketSaturation = sca.assessMarketSaturation(ctx, courseData.CourseType)
		
		// Technology and future-proofing
		analysis.TechnologyRelevance = sca.assessTechnologyRelevance(courseData.CourseType)
		analysis.FutureProofing = sca.assessFutureProofing(courseData.Skills)
	}

	return analysis
}

// analyzeGeographicFactors evaluates location-based factors
func (sca *StudentCreditAnalyzer) analyzeGeographicFactors(ctx sdk.Context, studentID string) types.GeographicFactors {
	locationData := sca.keeper.GetStudentLocationData(ctx, studentID)
	
	factors := types.GeographicFactors{
		StudentID: studentID,
	}

	if locationData != nil {
		// Regional economic factors
		factors.RegionalGDPGrowth = sca.getRegionalGDPGrowth(ctx, locationData.Region)
		factors.UnemploymentRate = sca.getRegionalUnemploymentRate(ctx, locationData.Region)
		factors.CostOfLiving = sca.getCostOfLiving(ctx, locationData.City)
		
		// Educational infrastructure
		factors.EducationalInfrastructure = sca.assessEducationalInfrastructure(ctx, locationData.Region)
		factors.IndustryPresence = sca.assessIndustryPresence(ctx, locationData.Region)
		
		// Migration patterns
		factors.OutMigrationRate = sca.getOutMigrationRate(ctx, locationData.Region)
		factors.ReturnMigrationRate = sca.getReturnMigrationRate(ctx, locationData.Region)
	}

	return factors
}

// analyzeSocialFactors evaluates social and demographic factors
func (sca *StudentCreditAnalyzer) analyzeSocialFactors(ctx sdk.Context, studentID string) types.SocialFactors {
	socialData := sca.keeper.GetStudentSocialData(ctx, studentID)
	
	factors := types.SocialFactors{
		StudentID: studentID,
	}

	if socialData != nil {
		// Demographic factors
		factors.Gender = socialData.Gender
		factors.SocialCategory = socialData.Category
		factors.RuralUrbanBackground = socialData.Background
		
		// First generation learner
		factors.FirstGenerationLearner = socialData.FirstGeneration
		
		// Social network and support
		factors.AlumniNetwork = sca.assessAlumniNetworkStrength(ctx, socialData.InstitutionID)
		factors.PeerSupport = sca.assessPeerSupportSystem(ctx, studentID)
		
		// Special categories and quotas
		factors.SpecialCategory = socialData.SpecialCategory
		factors.DisabilityStatus = socialData.DisabilityStatus
	}

	return factors
}

// calculateStudentCreditScore calculates overall credit score for student
func (sca *StudentCreditAnalyzer) calculateStudentCreditScore(profile *StudentCreditProfile) int64 {
	// Weighted scoring system for students
	weights := map[string]sdk.Dec{
		"academic":     sdk.NewDecWithPrec(25, 2), // 25%
		"family":       sdk.NewDecWithPrec(20, 2), // 20%
		"financial":    sdk.NewDecWithPrec(15, 2), // 15%
		"credit":       sdk.NewDecWithPrec(10, 2), // 10%
		"employment":   sdk.NewDecWithPrec(15, 2), // 15%
		"course":       sdk.NewDecWithPrec(10, 2), // 10%
		"social":       sdk.NewDecWithPrec(5, 2),  // 5%
	}

	// Score each component (0-850 scale)
	academicScore := sca.scoreAcademicProfile(profile.AcademicProfile)
	familyScore := sca.scoreFamilyProfile(profile.FamilyProfile)
	financialScore := sca.scoreFinancialProfile(profile.FinancialProfile)
	creditScore := sca.scoreCreditHistory(profile.CreditHistory)
	employmentScore := sca.scoreEmploymentPotential(profile.EmploymentPotential)
	courseScore := sca.scoreCourseAnalysis(profile.CourseAnalysis)
	socialScore := sca.scoreSocialFactors(profile.SocialFactors)

	// Calculate weighted average
	totalScore := sdk.ZeroDec()
	totalScore = totalScore.Add(academicScore.Mul(weights["academic"]))
	totalScore = totalScore.Add(familyScore.Mul(weights["family"]))
	totalScore = totalScore.Add(financialScore.Mul(weights["financial"]))
	totalScore = totalScore.Add(creditScore.Mul(weights["credit"]))
	totalScore = totalScore.Add(employmentScore.Mul(weights["employment"]))
	totalScore = totalScore.Add(courseScore.Mul(weights["course"]))
	totalScore = totalScore.Add(socialScore.Mul(weights["social"]))

	return totalScore.TruncateInt64()
}

// determineStudentRiskCategory categorizes students based on credit score
func (sca *StudentCreditAnalyzer) determineStudentRiskCategory(score int64) string {
	if score >= 750 {
		return "EXCELLENT"
	} else if score >= 700 {
		return "GOOD"
	} else if score >= 650 {
		return "FAIR"
	} else if score >= 600 {
		return "POOR"
	} else {
		return "HIGH_RISK"
	}
}

// calculateStudentLoanEligibility determines maximum loan amount for student
func (sca *StudentCreditAnalyzer) calculateStudentLoanEligibility(ctx sdk.Context, profile *StudentCreditProfile) sdk.Coin {
	params := sca.keeper.GetParams(ctx)
	
	// Base eligibility from family income
	baseEligibility := profile.FamilyProfile.MonthlyIncome.Mul(sdk.NewInt(60)) // 5 years of income
	
	// Risk-based multiplier
	var riskMultiplier sdk.Dec
	switch profile.RiskCategory {
	case "EXCELLENT":
		riskMultiplier = sdk.NewDecWithPrec(20, 1) // 2.0x
	case "GOOD":
		riskMultiplier = sdk.NewDecWithPrec(15, 1) // 1.5x
	case "FAIR":
		riskMultiplier = sdk.NewDecWithPrec(12, 1) // 1.2x
	case "POOR":
		riskMultiplier = sdk.NewDecWithPrec(8, 1)  // 0.8x
	case "HIGH_RISK":
		riskMultiplier = sdk.NewDecWithPrec(5, 1)  // 0.5x
	}

	maxAmount := baseEligibility.Mul(riskMultiplier).TruncateInt()

	// Apply course-specific caps
	courseTypeMultiplier := sca.getCourseTypeMultiplier(profile.CourseAnalysis.CourseType)
	maxAmount = maxAmount.ToDec().Mul(courseTypeMultiplier).TruncateInt()

	// Apply absolute caps
	if maxAmount.GT(params.MaxEducationLoanAmount.Amount) {
		maxAmount = params.MaxEducationLoanAmount.Amount
	}
	if maxAmount.LT(params.MinEducationLoanAmount.Amount) {
		maxAmount = params.MinEducationLoanAmount.Amount
	}

	return sdk.NewCoin(params.MaxEducationLoanAmount.Denom, maxAmount)
}

// Helper scoring functions for individual components
func (sca *StudentCreditAnalyzer) scoreAcademicProfile(profile types.AcademicProfile) sdk.Dec {
	score := sdk.NewDec(500) // Base score

	// Previous education score
	if profile.PreviousEducationScore >= 90 {
		score = score.Add(sdk.NewDec(100))
	} else if profile.PreviousEducationScore >= 80 {
		score = score.Add(sdk.NewDec(75))
	} else if profile.PreviousEducationScore >= 70 {
		score = score.Add(sdk.NewDec(50))
	} else if profile.PreviousEducationScore < 60 {
		score = score.Sub(sdk.NewDec(50))
	}

	// Entrance exam score
	if profile.EntranceExamScore >= 95 {
		score = score.Add(sdk.NewDec(75))
	} else if profile.EntranceExamScore >= 85 {
		score = score.Add(sdk.NewDec(50))
	} else if profile.EntranceExamScore >= 75 {
		score = score.Add(sdk.NewDec(25))
	}

	// Extracurricular activities
	if profile.ExtracurricularScore >= 8 {
		score = score.Add(sdk.NewDec(25))
	}

	// Research experience
	if profile.ResearchExperience {
		score = score.Add(sdk.NewDec(25))
	}

	// Cap at 850
	if score.GT(sdk.NewDec(850)) {
		score = sdk.NewDec(850)
	}
	if score.LT(sdk.NewDec(300)) {
		score = sdk.NewDec(300)
	}

	return score
}

// Additional scoring methods would be implemented for:
// - scoreFamilyProfile
// - scoreFinancialProfile
// - scoreCreditHistory
// - scoreEmploymentPotential
// - scoreCourseAnalysis
// - scoreSocialFactors

// Helper calculation functions
func (sca *StudentCreditAnalyzer) calculateAcademicConsistency(yearlyPerformance []sdk.Dec) sdk.Dec {
	if len(yearlyPerformance) < 2 {
		return sdk.ZeroDec()
	}

	// Calculate variance in performance
	var sum sdk.Dec
	for _, perf := range yearlyPerformance {
		sum = sum.Add(perf)
	}
	average := sum.QuoInt64(int64(len(yearlyPerformance)))

	var variance sdk.Dec
	for _, perf := range yearlyPerformance {
		diff := perf.Sub(average)
		variance = variance.Add(diff.Mul(diff))
	}
	variance = variance.QuoInt64(int64(len(yearlyPerformance)))

	// Convert variance to consistency score (lower variance = higher consistency)
	consistency := sdk.NewDec(100).Sub(variance.QuoInt64(10))
	if consistency.LT(sdk.ZeroDec()) {
		consistency = sdk.ZeroDec()
	}

	return consistency
}

func (sca *StudentCreditAnalyzer) getCourseTypeMultiplier(courseType string) sdk.Dec {
	// Course type multipliers based on employment potential
	multipliers := map[string]sdk.Dec{
		"ENGINEERING":    sdk.NewDecWithPrec(15, 1), // 1.5x
		"MEDICINE":       sdk.NewDecWithPrec(20, 1), // 2.0x
		"MBA":           sdk.NewDecWithPrec(18, 1), // 1.8x
		"LAW":           sdk.NewDecWithPrec(14, 1), // 1.4x
		"COMPUTER_SCIENCE": sdk.NewDecWithPrec(16, 1), // 1.6x
		"ARTS":          sdk.NewDecWithPrec(08, 1), // 0.8x
		"COMMERCE":      sdk.NewDecWithPrec(10, 1), // 1.0x
		"SCIENCE":       sdk.NewDecWithPrec(12, 1), // 1.2x
	}

	if multiplier, found := multipliers[courseType]; found {
		return multiplier
	}
	return sdk.OneDec() // Default 1.0x
}

// Additional helper methods would include all the calculation functions
// referenced in the analysis methods above