package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/namo/x/krishimitra/types"
)

// CreditScoringEngine handles comprehensive farmer credit assessment
type CreditScoringEngine struct {
	keeper Keeper
}

// NewCreditScoringEngine creates a new credit scoring engine
func NewCreditScoringEngine(keeper Keeper) *CreditScoringEngine {
	return &CreditScoringEngine{
		keeper: keeper,
	}
}

// FarmerCreditProfile represents comprehensive credit profile
type FarmerCreditProfile struct {
	FarmerID            string                     `json:"farmer_id"`
	BasicInfo           types.FarmerBasicInfo      `json:"basic_info"`
	LandOwnership       types.LandOwnershipInfo    `json:"land_ownership"`
	CropHistory         []types.CropRecord         `json:"crop_history"`
	FinancialHistory    types.FinancialHistory     `json:"financial_history"`
	WeatherRiskProfile  types.WeatherRiskProfile   `json:"weather_risk"`
	MarketAccessProfile types.MarketAccessProfile  `json:"market_access"`
	SocialProfile       types.SocialProfile        `json:"social_profile"`
	CreditScore         int64                      `json:"credit_score"`
	RiskCategory        string                     `json:"risk_category"`
	MaxLoanEligibility  sdk.Coin                   `json:"max_loan_eligibility"`
	RecommendedRate     sdk.Dec                    `json:"recommended_rate"`
}

// CalculateComprehensiveCreditScore calculates a detailed credit score
func (cse *CreditScoringEngine) CalculateComprehensiveCreditScore(ctx sdk.Context, farmerID string) (*FarmerCreditProfile, error) {
	// Get basic farmer information
	basicInfo, found := cse.keeper.GetFarmerBasicInfo(ctx, farmerID)
	if !found {
		return nil, fmt.Errorf("farmer basic information not found for ID: %s", farmerID)
	}

	profile := &FarmerCreditProfile{
		FarmerID:  farmerID,
		BasicInfo: basicInfo,
	}

	// Calculate individual scoring components
	landScore := cse.calculateLandOwnershipScore(ctx, farmerID)
	cropScore := cse.calculateCropHistoryScore(ctx, farmerID)
	financialScore := cse.calculateFinancialScore(ctx, farmerID)
	weatherScore := cse.calculateWeatherRiskScore(ctx, farmerID)
	marketScore := cse.calculateMarketAccessScore(ctx, farmerID)
	socialScore := cse.calculateSocialScore(ctx, farmerID)

	// Weighted credit score calculation
	weights := map[string]sdk.Dec{
		"land":      sdk.NewDecWithPrec(25, 2), // 25%
		"crop":      sdk.NewDecWithPrec(20, 2), // 20%
		"financial": sdk.NewDecWithPrec(25, 2), // 25%
		"weather":   sdk.NewDecWithPrec(10, 2), // 10%
		"market":    sdk.NewDecWithPrec(10, 2), // 10%
		"social":    sdk.NewDecWithPrec(10, 2), // 10%
	}

	totalScore := sdk.ZeroDec()
	totalScore = totalScore.Add(landScore.Mul(weights["land"]))
	totalScore = totalScore.Add(cropScore.Mul(weights["crop"]))
	totalScore = totalScore.Add(financialScore.Mul(weights["financial"]))
	totalScore = totalScore.Add(weatherScore.Mul(weights["weather"]))
	totalScore = totalScore.Add(marketScore.Mul(weights["market"]))
	totalScore = totalScore.Add(socialScore.Mul(weights["social"]))

	profile.CreditScore = totalScore.TruncateInt64()

	// Determine risk category
	profile.RiskCategory = cse.determineRiskCategory(profile.CreditScore)

	// Calculate max loan eligibility
	profile.MaxLoanEligibility = cse.calculateMaxLoanEligibility(ctx, profile)

	// Calculate recommended interest rate
	profile.RecommendedRate = cse.calculateRecommendedRate(ctx, profile)

	return profile, nil
}

// calculateLandOwnershipScore evaluates land-related factors
func (cse *CreditScoringEngine) calculateLandOwnershipScore(ctx sdk.Context, farmerID string) sdk.Dec {
	landInfo, found := cse.keeper.GetLandOwnershipInfo(ctx, farmerID)
	if !found {
		return sdk.NewDec(300) // Low score for no land documentation
	}

	score := sdk.NewDec(500) // Base score

	// Land size factor (more land = higher score, with diminishing returns)
	if landInfo.TotalArea.GT(sdk.NewDec(10)) {
		score = score.Add(sdk.NewDec(150)) // Large farmer bonus
	} else if landInfo.TotalArea.GT(sdk.NewDec(5)) {
		score = score.Add(sdk.NewDec(100)) // Medium farmer bonus
	} else if landInfo.TotalArea.GT(sdk.NewDec(2)) {
		score = score.Add(sdk.NewDec(50)) // Small farmer modest bonus
	}

	// Ownership type factor
	if landInfo.OwnershipType == "OWNED" {
		score = score.Add(sdk.NewDec(100))
	} else if landInfo.OwnershipType == "LEASED_LONG_TERM" {
		score = score.Add(sdk.NewDec(50))
	}

	// Irrigation factor
	if landInfo.IrrigationAvailable {
		score = score.Add(sdk.NewDec(75))
	}

	// Soil quality factor
	switch landInfo.SoilQuality {
	case "EXCELLENT":
		score = score.Add(sdk.NewDec(75))
	case "GOOD":
		score = score.Add(sdk.NewDec(50))
	case "FAIR":
		score = score.Add(sdk.NewDec(25))
	}

	// Cap at 850
	if score.GT(sdk.NewDec(850)) {
		score = sdk.NewDec(850)
	}

	return score
}

// calculateCropHistoryScore evaluates farming track record
func (cse *CreditScoringEngine) calculateCropHistoryScore(ctx sdk.Context, farmerID string) sdk.Dec {
	cropHistory := cse.keeper.GetCropHistory(ctx, farmerID)
	if len(cropHistory) == 0 {
		return sdk.NewDec(400) // No history
	}

	score := sdk.NewDec(500) // Base score
	
	// Experience factor (more seasons = higher score)
	yearsOfExperience := len(cropHistory)
	if yearsOfExperience >= 10 {
		score = score.Add(sdk.NewDec(150))
	} else if yearsOfExperience >= 5 {
		score = score.Add(sdk.NewDec(100))
	} else if yearsOfExperience >= 3 {
		score = score.Add(sdk.NewDec(50))
	}

	// Success rate factor
	successfulSeasons := 0
	totalRevenue := sdk.ZeroDec()
	totalCost := sdk.ZeroDec()

	for _, record := range cropHistory {
		if record.HarvestStatus == "SUCCESSFUL" {
			successfulSeasons++
		}
		if !record.Revenue.IsZero() {
			totalRevenue = totalRevenue.Add(record.Revenue.Amount.ToDec())
		}
		if !record.Cost.IsZero() {
			totalCost = totalCost.Add(record.Cost.Amount.ToDec())
		}
	}

	if yearsOfExperience > 0 {
		successRate := sdk.NewDec(int64(successfulSeasons)).Quo(sdk.NewDec(int64(yearsOfExperience)))
		if successRate.GT(sdk.NewDecWithPrec(8, 1)) { // >80% success
			score = score.Add(sdk.NewDec(100))
		} else if successRate.GT(sdk.NewDecWithPrec(6, 1)) { // >60% success
			score = score.Add(sdk.NewDec(50))
		}
	}

	// Profitability factor
	if !totalCost.IsZero() && totalRevenue.GT(totalCost) {
		profitMargin := totalRevenue.Sub(totalCost).Quo(totalCost)
		if profitMargin.GT(sdk.NewDecWithPrec(3, 1)) { // >30% profit margin
			score = score.Add(sdk.NewDec(75))
		} else if profitMargin.GT(sdk.NewDecWithPrec(15, 2)) { // >15% profit margin
			score = score.Add(sdk.NewDec(50))
		}
	}

	// Crop diversification factor
	uniqueCrops := make(map[string]bool)
	for _, record := range cropHistory {
		uniqueCrops[record.CropType] = true
	}
	if len(uniqueCrops) >= 3 {
		score = score.Add(sdk.NewDec(50)) // Diversification bonus
	}

	// Cap at 850
	if score.GT(sdk.NewDec(850)) {
		score = sdk.NewDec(850)
	}

	return score
}

// calculateFinancialScore evaluates financial history and behavior
func (cse *CreditScoringEngine) calculateFinancialScore(ctx sdk.Context, farmerID string) sdk.Dec {
	financialHistory, found := cse.keeper.GetFinancialHistory(ctx, farmerID)
	if !found {
		return sdk.NewDec(350) // No financial history
	}

	score := sdk.NewDec(500) // Base score

	// Previous loan repayment history
	if financialHistory.TotalLoans > 0 {
		repaymentRate := sdk.NewDec(financialHistory.SuccessfulRepayments).Quo(sdk.NewDec(financialHistory.TotalLoans))
		if repaymentRate.Equal(sdk.OneDec()) { // 100% repayment
			score = score.Add(sdk.NewDec(150))
		} else if repaymentRate.GT(sdk.NewDecWithPrec(9, 1)) { // >90% repayment
			score = score.Add(sdk.NewDec(100))
		} else if repaymentRate.GT(sdk.NewDecWithPrec(8, 1)) { // >80% repayment
			score = score.Add(sdk.NewDec(50))
		} else if repaymentRate.LT(sdk.NewDecWithPrec(5, 1)) { // <50% repayment
			score = score.Sub(sdk.NewDec(100))
		}
	}

	// Credit utilization
	if !financialHistory.CreditLimit.IsZero() && !financialHistory.OutstandingDebt.IsZero() {
		utilization := financialHistory.OutstandingDebt.Amount.ToDec().Quo(financialHistory.CreditLimit.Amount.ToDec())
		if utilization.LT(sdk.NewDecWithPrec(3, 1)) { // <30% utilization
			score = score.Add(sdk.NewDec(50))
		} else if utilization.GT(sdk.NewDecWithPrec(8, 1)) { // >80% utilization
			score = score.Sub(sdk.NewDec(50))
		}
	}

	// Banking relationship
	if financialHistory.BankAccountAge >= 5 { // 5+ years banking
		score = score.Add(sdk.NewDec(50))
	} else if financialHistory.BankAccountAge >= 2 { // 2+ years banking
		score = score.Add(sdk.NewDec(25))
	}

	// Income stability
	if len(financialHistory.MonthlyIncomes) >= 12 {
		incomeVariability := cse.calculateIncomeVariability(financialHistory.MonthlyIncomes)
		if incomeVariability.LT(sdk.NewDecWithPrec(2, 1)) { // <20% variability
			score = score.Add(sdk.NewDec(50))
		} else if incomeVariability.GT(sdk.NewDecWithPrec(5, 1)) { // >50% variability
			score = score.Sub(sdk.NewDec(25))
		}
	}

	// Cap at 850
	if score.GT(sdk.NewDec(850)) {
		score = sdk.NewDec(850)
	}

	return score
}

// calculateWeatherRiskScore evaluates climate and weather-related risks
func (cse *CreditScoringEngine) calculateWeatherRiskScore(ctx sdk.Context, farmerID string) sdk.Dec {
	weatherProfile, found := cse.keeper.GetWeatherRiskProfile(ctx, farmerID)
	if !found {
		return sdk.NewDec(500) // Average score if no data
	}

	score := sdk.NewDec(600) // Base score (lower risk = higher score)

	// Climate zone risk
	switch weatherProfile.ClimateZone {
	case "ARID", "SEMI_ARID":
		score = score.Sub(sdk.NewDec(100)) // Higher risk
	case "TROPICAL", "SUBTROPICAL":
		score = score.Add(sdk.NewDec(50)) // Lower risk
	case "TEMPERATE":
		score = score.Add(sdk.NewDec(75)) // Optimal
	}

	// Drought frequency (last 10 years)
	if weatherProfile.DroughtFrequency >= 5 { // 50%+ drought years
		score = score.Sub(sdk.NewDec(150))
	} else if weatherProfile.DroughtFrequency >= 3 { // 30%+ drought years
		score = score.Sub(sdk.NewDec(75))
	} else if weatherProfile.DroughtFrequency <= 1 { // <=10% drought years
		score = score.Add(sdk.NewDec(50))
	}

	// Flood risk
	if weatherProfile.FloodRisk == "HIGH" {
		score = score.Sub(sdk.NewDec(100))
	} else if weatherProfile.FloodRisk == "MEDIUM" {
		score = score.Sub(sdk.NewDec(50))
	} else if weatherProfile.FloodRisk == "LOW" {
		score = score.Add(sdk.NewDec(25))
	}

	// Rainfall variability
	if weatherProfile.RainfallVariability.GT(sdk.NewDecWithPrec(4, 1)) { // >40% variability
		score = score.Sub(sdk.NewDec(75))
	} else if weatherProfile.RainfallVariability.LT(sdk.NewDecWithPrec(2, 1)) { // <20% variability
		score = score.Add(sdk.NewDec(50))
	}

	// Insurance coverage
	if weatherProfile.HasCropInsurance {
		score = score.Add(sdk.NewDec(100)) // Major boost for insured farmers
	}

	// Cap between 300-800 (weather is inherently risky)
	if score.GT(sdk.NewDec(800)) {
		score = sdk.NewDec(800)
	}
	if score.LT(sdk.NewDec(300)) {
		score = sdk.NewDec(300)
	}

	return score
}

// calculateMarketAccessScore evaluates market connectivity and access
func (cse *CreditScoringEngine) calculateMarketAccessScore(ctx sdk.Context, farmerID string) sdk.Dec {
	marketProfile, found := cse.keeper.GetMarketAccessProfile(ctx, farmerID)
	if !found {
		return sdk.NewDec(450) // Below average if no data
	}

	score := sdk.NewDec(500) // Base score

	// Distance to nearest market
	if marketProfile.DistanceToNearestMarket.LT(sdk.NewDec(10)) { // <10 km
		score = score.Add(sdk.NewDec(100))
	} else if marketProfile.DistanceToNearestMarket.LT(sdk.NewDec(25)) { // <25 km
		score = score.Add(sdk.NewDec(50))
	} else if marketProfile.DistanceToNearestMarket.GT(sdk.NewDec(50)) { // >50 km
		score = score.Sub(sdk.NewDec(50))
	}

	// Transportation access
	if marketProfile.HasOwnTransport {
		score = score.Add(sdk.NewDec(75))
	} else if marketProfile.HasAccessToTransport {
		score = score.Add(sdk.NewDec(50))
	}

	// Storage facilities
	if marketProfile.HasStorageFacility {
		score = score.Add(sdk.NewDec(75))
	}

	// Price realization
	if !marketProfile.AverageMarketPrice.IsZero() && !marketProfile.AverageSalePrice.IsZero() {
		priceRealization := marketProfile.AverageSalePrice.Amount.ToDec().Quo(marketProfile.AverageMarketPrice.Amount.ToDec())
		if priceRealization.GT(sdk.NewDecWithPrec(9, 1)) { // >90% price realization
			score = score.Add(sdk.NewDec(100))
		} else if priceRealization.GT(sdk.NewDecWithPrec(8, 1)) { // >80% price realization
			score = score.Add(sdk.NewDec(50))
		} else if priceRealization.LT(sdk.NewDecWithPrec(7, 1)) { // <70% price realization
			score = score.Sub(sdk.NewDec(50))
		}
	}

	// Market connections
	if marketProfile.NumberOfBuyers >= 5 {
		score = score.Add(sdk.NewDec(50)) // Multiple buyer relationships
	} else if marketProfile.NumberOfBuyers >= 2 {
		score = score.Add(sdk.NewDec(25))
	}

	// Digital market access
	if marketProfile.UsesDigitalPlatforms {
		score = score.Add(sdk.NewDec(50))
	}

	// Cap at 850
	if score.GT(sdk.NewDec(850)) {
		score = sdk.NewDec(850)
	}

	return score
}

// calculateSocialScore evaluates social factors
func (cse *CreditScoringEngine) calculateSocialScore(ctx sdk.Context, farmerID string) sdk.Dec {
	socialProfile, found := cse.keeper.GetSocialProfile(ctx, farmerID)
	if !found {
		return sdk.NewDec(500) // Average score if no data
	}

	score := sdk.NewDec(500) // Base score

	// Education level
	switch socialProfile.EducationLevel {
	case "GRADUATE", "POST_GRADUATE":
		score = score.Add(sdk.NewDec(75))
	case "HIGHER_SECONDARY":
		score = score.Add(sdk.NewDec(50))
	case "SECONDARY":
		score = score.Add(sdk.NewDec(25))
	case "PRIMARY":
		score = score.Add(sdk.NewDec(10))
	}

	// Age factor (experience vs. innovation balance)
	if socialProfile.Age >= 30 && socialProfile.Age <= 50 {
		score = score.Add(sdk.NewDec(50)) // Optimal age range
	} else if socialProfile.Age > 60 {
		score = score.Sub(sdk.NewDec(25)) // Aging farmer risk
	}

	// Family support
	if socialProfile.FamilyMembersInAgriculture >= 2 {
		score = score.Add(sdk.NewDec(50)) // Family farming support
	}

	// Technology adoption
	if socialProfile.SmartphoneAccess {
		score = score.Add(sdk.NewDec(25))
	}
	if socialProfile.InternetAccess {
		score = score.Add(sdk.NewDec(25))
	}

	// Group membership (SHG, FPO, etc.)
	if socialProfile.IsMemberOfFPO {
		score = score.Add(sdk.NewDec(50))
	}
	if socialProfile.IsMemberOfSHG {
		score = score.Add(sdk.NewDec(25))
	}

	// Training and certification
	if socialProfile.HasReceivedTraining {
		score = score.Add(sdk.NewDec(25))
	}
	if socialProfile.HasCertifications {
		score = score.Add(sdk.NewDec(50))
	}

	// Cap at 850
	if score.GT(sdk.NewDec(850)) {
		score = sdk.NewDec(850)
	}

	return score
}

// determineRiskCategory categorizes farmers based on credit score
func (cse *CreditScoringEngine) determineRiskCategory(score int64) string {
	if score >= 750 {
		return "LOW_RISK"
	} else if score >= 650 {
		return "MEDIUM_RISK"
	} else if score >= 550 {
		return "HIGH_RISK"
	} else {
		return "VERY_HIGH_RISK"
	}
}

// calculateMaxLoanEligibility determines maximum loan amount
func (cse *CreditScoringEngine) calculateMaxLoanEligibility(ctx sdk.Context, profile *FarmerCreditProfile) sdk.Coin {
	params := cse.keeper.GetParams(ctx)
	
	// Base eligibility from land value and income
	landInfo, _ := cse.keeper.GetLandOwnershipInfo(ctx, profile.FarmerID)
	baseEligibility := landInfo.TotalArea.Mul(params.LoanPerAcre) // Base on land size

	// Risk-based multiplier
	var riskMultiplier sdk.Dec
	switch profile.RiskCategory {
	case "LOW_RISK":
		riskMultiplier = sdk.NewDecWithPrec(15, 1) // 1.5x
	case "MEDIUM_RISK":
		riskMultiplier = sdk.NewDecWithPrec(12, 1) // 1.2x
	case "HIGH_RISK":
		riskMultiplier = sdk.NewDecWithPrec(8, 1)  // 0.8x
	case "VERY_HIGH_RISK":
		riskMultiplier = sdk.NewDecWithPrec(5, 1)  // 0.5x
	}

	maxAmount := baseEligibility.Mul(riskMultiplier).TruncateInt()

	// Apply absolute caps
	if maxAmount.GT(params.MaxLoanAmount.Amount) {
		maxAmount = params.MaxLoanAmount.Amount
	}
	if maxAmount.LT(params.MinLoanAmount.Amount) {
		maxAmount = params.MinLoanAmount.Amount
	}

	return sdk.NewCoin(params.MaxLoanAmount.Denom, maxAmount)
}

// calculateRecommendedRate determines optimal interest rate
func (cse *CreditScoringEngine) calculateRecommendedRate(ctx sdk.Context, profile *FarmerCreditProfile) sdk.Dec {
	params := cse.keeper.GetParams(ctx)
	baseRate, _ := sdk.NewDecFromStr(params.BaseInterestRate)

	// Risk-based adjustment
	var riskAdjustment sdk.Dec
	switch profile.RiskCategory {
	case "LOW_RISK":
		riskAdjustment = sdk.NewDecWithPrec(-100, 4) // -1%
	case "MEDIUM_RISK":
		riskAdjustment = sdk.NewDecWithPrec(-25, 4)  // -0.25%
	case "HIGH_RISK":
		riskAdjustment = sdk.NewDecWithPrec(50, 4)   // +0.5%
	case "VERY_HIGH_RISK":
		riskAdjustment = sdk.NewDecWithPrec(150, 4)  // +1.5%
	}

	finalRate := baseRate.Add(riskAdjustment)

	// Apply bounds
	minRate, _ := sdk.NewDecFromStr(params.MinInterestRate)
	maxRate, _ := sdk.NewDecFromStr(params.MaxInterestRate)

	if finalRate.LT(minRate) {
		finalRate = minRate
	}
	if finalRate.GT(maxRate) {
		finalRate = maxRate
	}

	return finalRate
}

// calculateIncomeVariability calculates coefficient of variation for income
func (cse *CreditScoringEngine) calculateIncomeVariability(incomes []sdk.Coin) sdk.Dec {
	if len(incomes) == 0 {
		return sdk.OneDec() // 100% variability if no data
	}

	// Calculate mean
	total := sdk.ZeroDec()
	for _, income := range incomes {
		total = total.Add(income.Amount.ToDec())
	}
	mean := total.QuoInt64(int64(len(incomes)))

	if mean.IsZero() {
		return sdk.OneDec()
	}

	// Calculate variance
	variance := sdk.ZeroDec()
	for _, income := range incomes {
		deviation := income.Amount.ToDec().Sub(mean)
		variance = variance.Add(deviation.Mul(deviation))
	}
	variance = variance.QuoInt64(int64(len(incomes)))

	// Calculate coefficient of variation (standard deviation / mean)
	stdDev := variance.ApproxSqrt()
	return stdDev.Quo(mean)
}

// AssessLoanEligibility provides comprehensive loan eligibility assessment
func (cse *CreditScoringEngine) AssessLoanEligibility(ctx sdk.Context, farmerID string, requestedAmount sdk.Coin, cropType string) (*types.LoanEligibilityAssessment, error) {
	// Get comprehensive credit profile
	profile, err := cse.CalculateComprehensiveCreditScore(ctx, farmerID)
	if err != nil {
		return nil, err
	}

	assessment := &types.LoanEligibilityAssessment{
		FarmerID:            farmerID,
		RequestedAmount:     requestedAmount,
		CropType:            cropType,
		CreditScore:         profile.CreditScore,
		RiskCategory:        profile.RiskCategory,
		MaxEligibleAmount:   profile.MaxLoanEligibility,
		RecommendedRate:     profile.RecommendedRate,
		AssessmentTime:      ctx.BlockTime(),
	}

	// Determine eligibility
	if requestedAmount.Amount.LTE(profile.MaxLoanEligibility.Amount) && profile.CreditScore >= 550 {
		assessment.IsEligible = true
		assessment.EligibilityReason = "Meets all eligibility criteria"
	} else {
		assessment.IsEligible = false
		if profile.CreditScore < 550 {
			assessment.EligibilityReason = "Credit score below minimum threshold (550)"
		} else {
			assessment.EligibilityReason = fmt.Sprintf("Requested amount exceeds maximum eligibility of %s", profile.MaxLoanEligibility.String())
		}
	}

	// Add recommendations
	assessment.Recommendations = cse.generateRecommendations(ctx, profile, requestedAmount)

	return assessment, nil
}

// generateRecommendations provides actionable recommendations
func (cse *CreditScoringEngine) generateRecommendations(ctx sdk.Context, profile *FarmerCreditProfile, requestedAmount sdk.Coin) []string {
	recommendations := []string{}

	if profile.CreditScore < 650 {
		recommendations = append(recommendations, "Consider building credit history with smaller loans")
		recommendations = append(recommendations, "Maintain consistent farming records and income documentation")
	}

	if !profile.LandOwnership.IrrigationAvailable {
		recommendations = append(recommendations, "Consider irrigation investment to improve land value and reduce weather risk")
	}

	if len(profile.CropHistory) < 3 {
		recommendations = append(recommendations, "Build farming track record with consistent crop production")
	}

	if !profile.WeatherRiskProfile.HasCropInsurance {
		recommendations = append(recommendations, "Obtain crop insurance to reduce weather-related risks")
	}

	if profile.MarketAccessProfile.DistanceToNearestMarket.GT(sdk.NewDec(25)) {
		recommendations = append(recommendations, "Explore transportation and storage solutions to improve market access")
	}

	if requestedAmount.Amount.GT(profile.MaxLoanEligibility.Amount) {
		recommendations = append(recommendations, fmt.Sprintf("Consider applying for %s initially and building eligibility for larger amounts", profile.MaxLoanEligibility.String()))
	}

	return recommendations
}