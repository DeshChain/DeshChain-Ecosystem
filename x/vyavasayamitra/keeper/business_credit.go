package keeper

import (
	"fmt"
	"sort"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/namo/x/vyavasayamitra/types"
)

// BusinessCreditAnalyzer handles comprehensive business credit analysis
type BusinessCreditAnalyzer struct {
	keeper Keeper
}

// NewBusinessCreditAnalyzer creates a new business credit analyzer
func NewBusinessCreditAnalyzer(keeper Keeper) *BusinessCreditAnalyzer {
	return &BusinessCreditAnalyzer{
		keeper: keeper,
	}
}

// BusinessCreditProfile represents comprehensive business credit assessment
type BusinessCreditProfile struct {
	BusinessID              string                        `json:"business_id"`
	BusinessInfo            types.BusinessInfo            `json:"business_info"`
	FinancialMetrics        types.FinancialMetrics        `json:"financial_metrics"`
	CashFlowAnalysis        types.CashFlowAnalysis        `json:"cash_flow_analysis"`
	IndustryAnalysis        types.IndustryAnalysis        `json:"industry_analysis"`
	MarketPosition          types.MarketPosition          `json:"market_position"`
	OperationalMetrics      types.OperationalMetrics      `json:"operational_metrics"`
	ComplianceProfile       types.ComplianceProfile       `json:"compliance_profile"`
	DigitalFootprint        types.DigitalFootprint        `json:"digital_footprint"`
	SupplyChainAnalysis     types.SupplyChainAnalysis     `json:"supply_chain_analysis"`
	TechnologyAdoption      types.TechnologyAdoption      `json:"technology_adoption"`
	CreditScore             int64                         `json:"credit_score"`
	RiskCategory            string                        `json:"risk_category"`
	MaxLoanEligibility      sdk.Coin                      `json:"max_loan_eligibility"`
	RecommendedRate         sdk.Dec                       `json:"recommended_rate"`
	CreditLineEligibility   sdk.Coin                      `json:"credit_line_eligibility"`
	SpecialPrograms         []string                      `json:"special_programs"`
	RedFlags                []string                      `json:"red_flags"`
	Recommendations         []string                      `json:"recommendations"`
	LastAssessmentDate      time.Time                     `json:"last_assessment_date"`
}

// AnalyzeBusinessCredit performs comprehensive business credit analysis
func (bca *BusinessCreditAnalyzer) AnalyzeBusinessCredit(ctx sdk.Context, businessID string) (*BusinessCreditProfile, error) {
	// Get business information
	businessInfo, found := bca.keeper.GetBusinessInfo(ctx, businessID)
	if !found {
		return nil, fmt.Errorf("business information not found for ID: %s", businessID)
	}

	profile := &BusinessCreditProfile{
		BusinessID:   businessID,
		BusinessInfo: businessInfo,
	}

	// Perform comprehensive analysis
	profile.FinancialMetrics = bca.analyzeFinancialMetrics(ctx, businessID)
	profile.CashFlowAnalysis = bca.analyzeCashFlow(ctx, businessID)
	profile.IndustryAnalysis = bca.analyzeIndustryFactors(ctx, businessID)
	profile.MarketPosition = bca.analyzeMarketPosition(ctx, businessID)
	profile.OperationalMetrics = bca.analyzeOperationalMetrics(ctx, businessID)
	profile.ComplianceProfile = bca.analyzeCompliance(ctx, businessID)
	profile.DigitalFootprint = bca.analyzeDigitalPresence(ctx, businessID)
	profile.SupplyChainAnalysis = bca.analyzeSupplyChain(ctx, businessID)
	profile.TechnologyAdoption = bca.analyzeTechnologyAdoption(ctx, businessID)

	// Calculate composite credit score
	profile.CreditScore = bca.calculateCompositeCreditScore(profile)

	// Determine risk category
	profile.RiskCategory = bca.determineRiskCategory(profile.CreditScore)

	// Calculate loan eligibility
	profile.MaxLoanEligibility = bca.calculateMaxLoanEligibility(ctx, profile)

	// Calculate recommended interest rate
	profile.RecommendedRate = bca.calculateRecommendedRate(ctx, profile)

	// Calculate credit line eligibility
	profile.CreditLineEligibility = bca.calculateCreditLineEligibility(ctx, profile)

	// Identify special programs
	profile.SpecialPrograms = bca.identifySpecialPrograms(ctx, profile)

	// Identify red flags
	profile.RedFlags = bca.identifyRedFlags(profile)

	// Generate recommendations
	profile.Recommendations = bca.generateRecommendations(profile)

	profile.LastAssessmentDate = ctx.BlockTime()

	return profile, nil
}

// analyzeFinancialMetrics evaluates financial health and performance
func (bca *BusinessCreditAnalyzer) analyzeFinancialMetrics(ctx sdk.Context, businessID string) types.FinancialMetrics {
	financialData := bca.keeper.GetFinancialData(ctx, businessID)
	
	metrics := types.FinancialMetrics{
		BusinessID: businessID,
	}

	if len(financialData.YearlyRevenues) > 0 {
		// Revenue analysis
		metrics.AnnualRevenue = financialData.YearlyRevenues[len(financialData.YearlyRevenues)-1].Amount
		metrics.RevenueGrowthRate = bca.calculateRevenueGrowthRate(financialData.YearlyRevenues)
		metrics.RevenueStability = bca.calculateRevenueStability(financialData.MonthlyRevenues)

		// Profitability analysis
		if len(financialData.YearlyProfits) > 0 {
			latestProfit := financialData.YearlyProfits[len(financialData.YearlyProfits)-1].Amount
			metrics.ProfitMargin = latestProfit.ToDec().Quo(metrics.AnnualRevenue.Amount.ToDec())
			metrics.ProfitGrowthRate = bca.calculateProfitGrowthRate(financialData.YearlyProfits)
		}

		// Liquidity ratios
		if !financialData.CurrentAssets.IsZero() && !financialData.CurrentLiabilities.IsZero() {
			metrics.CurrentRatio = financialData.CurrentAssets.Amount.ToDec().Quo(financialData.CurrentLiabilities.Amount.ToDec())
			
			// Quick ratio (more conservative)
			quickAssets := financialData.CurrentAssets.Amount.Sub(financialData.Inventory.Amount)
			if quickAssets.IsPositive() {
				metrics.QuickRatio = quickAssets.ToDec().Quo(financialData.CurrentLiabilities.Amount.ToDec())
			}
		}

		// Leverage ratios
		totalDebt := financialData.TotalLiabilities.Amount
		totalEquity := financialData.TotalAssets.Amount.Sub(totalDebt)
		if !totalEquity.IsZero() {
			metrics.DebtToEquityRatio = totalDebt.ToDec().Quo(totalEquity.ToDec())
		}

		// Interest coverage ratio
		if !financialData.InterestExpense.IsZero() && len(financialData.YearlyProfits) > 0 {
			ebit := financialData.YearlyProfits[len(financialData.YearlyProfits)-1].Amount.Add(financialData.InterestExpense.Amount)
			metrics.InterestCoverageRatio = ebit.ToDec().Quo(financialData.InterestExpense.Amount.ToDec())
		}

		// Working capital
		metrics.WorkingCapital = financialData.CurrentAssets.Sub(financialData.CurrentLiabilities)

		// Asset turnover
		if !financialData.TotalAssets.IsZero() {
			metrics.AssetTurnoverRatio = metrics.AnnualRevenue.Amount.ToDec().Quo(financialData.TotalAssets.Amount.ToDec())
		}
	}

	return metrics
}

// analyzeCashFlow evaluates cash flow patterns and health
func (bca *BusinessCreditAnalyzer) analyzeCashFlow(ctx sdk.Context, businessID string) types.CashFlowAnalysis {
	cashFlowData := bca.keeper.GetCashFlowData(ctx, businessID)
	
	analysis := types.CashFlowAnalysis{
		BusinessID: businessID,
	}

	if len(cashFlowData.MonthlyCashFlows) > 0 {
		// Calculate operating cash flow
		analysis.OperatingCashFlow = bca.calculateOperatingCashFlow(cashFlowData.MonthlyCashFlows)
		
		// Cash flow stability
		analysis.CashFlowStability = bca.calculateCashFlowStability(cashFlowData.MonthlyCashFlows)
		
		// Seasonal patterns
		analysis.SeasonalityIndex = bca.calculateSeasonalityIndex(cashFlowData.MonthlyCashFlows)
		
		// Cash conversion cycle
		if !cashFlowData.InventoryTurnover.IsZero() && !cashFlowData.ReceivablesTurnover.IsZero() && !cashFlowData.PayablesTurnover.IsZero() {
			daysInventory := sdk.NewDec(365).Quo(cashFlowData.InventoryTurnover)
			daysReceivables := sdk.NewDec(365).Quo(cashFlowData.ReceivablesTurnover)
			daysPayables := sdk.NewDec(365).Quo(cashFlowData.PayablesTurnover)
			analysis.CashConversionCycle = daysInventory.Add(daysReceivables).Sub(daysPayables)
		}

		// Free cash flow
		if !cashFlowData.CapitalExpenditures.IsZero() {
			analysis.FreeCashFlow = analysis.OperatingCashFlow.Sub(cashFlowData.CapitalExpenditures)
		}

		// Cash burnrate for startups
		if len(cashFlowData.MonthlyCashFlows) >= 6 {
			recentFlows := cashFlowData.MonthlyCashFlows[len(cashFlowData.MonthlyCashFlows)-6:]
			totalBurn := sdk.ZeroInt()
			for _, flow := range recentFlows {
				if flow.Amount.IsNegative() {
					totalBurn = totalBurn.Add(flow.Amount.Abs())
				}
			}
			analysis.MonthlyBurnRate = totalBurn.QuoRaw(6)
		}

		// Runway calculation
		if !analysis.MonthlyBurnRate.IsZero() && !cashFlowData.CashReserves.IsZero() {
			analysis.CashRunwayMonths = cashFlowData.CashReserves.Amount.Quo(analysis.MonthlyBurnRate).Int64()
		}
	}

	return analysis
}

// analyzeIndustryFactors evaluates industry-specific factors
func (bca *BusinessCreditAnalyzer) analyzeIndustryFactors(ctx sdk.Context, businessID string) types.IndustryAnalysis {
	businessInfo, _ := bca.keeper.GetBusinessInfo(ctx, businessID)
	
	analysis := types.IndustryAnalysis{
		BusinessID:   businessID,
		IndustryType: businessInfo.Industry,
	}

	// Get industry benchmarks
	benchmarks := bca.keeper.GetIndustryBenchmarks(ctx, businessInfo.Industry)
	if benchmarks != nil {
		analysis.IndustryGrowthRate = benchmarks.AverageGrowthRate
		analysis.IndustryProfitMargin = benchmarks.AverageProfitMargin
		analysis.IndustryRiskLevel = benchmarks.RiskLevel
		analysis.CyclicalityFactor = benchmarks.CyclicalityIndex
		analysis.RegulatoryComplexity = benchmarks.RegulatoryComplexity
	}

	// Industry-specific risk factors
	analysis.IndustryRiskFactors = bca.getIndustryRiskFactors(businessInfo.Industry)
	
	// Competitive landscape
	analysis.CompetitionLevel = bca.assessCompetitionLevel(ctx, businessInfo.Industry, businessInfo.Location)
	
	// Market outlook
	analysis.MarketOutlook = bca.getMarketOutlook(ctx, businessInfo.Industry)

	return analysis
}

// analyzeMarketPosition evaluates business's position in the market
func (bca *BusinessCreditAnalyzer) analyzeMarketPosition(ctx sdk.Context, businessID string) types.MarketPosition {
	businessInfo, _ := bca.keeper.GetBusinessInfo(ctx, businessID)
	
	position := types.MarketPosition{
		BusinessID: businessID,
	}

	// Market share analysis
	if !businessInfo.EstimatedMarketShare.IsZero() {
		position.MarketShare = businessInfo.EstimatedMarketShare
	}

	// Customer concentration
	customerData := bca.keeper.GetCustomerData(ctx, businessID)
	if len(customerData.TopCustomers) > 0 {
		position.CustomerConcentration = bca.calculateCustomerConcentration(customerData.TopCustomers)
		position.CustomerStickiness = bca.calculateCustomerStickiness(customerData.CustomerRetentionData)
	}

	// Competitive advantages
	position.CompetitiveAdvantages = bca.identifyCompetitiveAdvantages(ctx, businessID)
	
	// Supplier relationships
	supplierData := bca.keeper.GetSupplierData(ctx, businessID)
	if len(supplierData.KeySuppliers) > 0 {
		position.SupplierConcentration = bca.calculateSupplierConcentration(supplierData.KeySuppliers)
		position.SupplierRelationshipQuality = bca.assessSupplierRelationships(supplierData)
	}

	// Brand strength
	position.BrandStrength = bca.assessBrandStrength(ctx, businessID)
	
	// Geographic diversification
	position.GeographicDiversification = bca.calculateGeographicDiversification(businessInfo.OperatingLocations)

	return position
}

// analyzeOperationalMetrics evaluates operational efficiency and capabilities
func (bca *BusinessCreditAnalyzer) analyzeOperationalMetrics(ctx sdk.Context, businessID string) types.OperationalMetrics {
	operationalData := bca.keeper.GetOperationalData(ctx, businessID)
	
	metrics := types.OperationalMetrics{
		BusinessID: businessID,
	}

	// Productivity metrics
	if !operationalData.NumberOfEmployees.IsZero() && len(operationalData.MonthlyProduction) > 0 {
		latestProduction := operationalData.MonthlyProduction[len(operationalData.MonthlyProduction)-1]
		metrics.ProductivityPerEmployee = latestProduction.Amount.ToDec().Quo(operationalData.NumberOfEmployees.ToDec())
	}

	// Quality metrics
	metrics.QualityScore = operationalData.QualityScore
	metrics.CustomerSatisfactionScore = operationalData.CustomerSatisfactionScore
	
	// Efficiency metrics
	if !operationalData.TotalCosts.IsZero() && !operationalData.TotalRevenue.IsZero() {
		metrics.OperationalEfficiency = operationalData.TotalRevenue.Amount.ToDec().Quo(operationalData.TotalCosts.Amount.ToDec())
	}

	// Innovation metrics
	metrics.InnovationIndex = bca.calculateInnovationIndex(ctx, businessID)
	
	// Inventory management
	if len(operationalData.InventoryLevels) > 0 {
		metrics.InventoryTurnover = bca.calculateInventoryTurnover(operationalData.InventoryLevels)
	}

	// Employee metrics
	metrics.EmployeeRetentionRate = operationalData.EmployeeRetentionRate
	metrics.EmployeeProductivityGrowth = bca.calculateEmployeeProductivityGrowth(operationalData.ProductivityHistory)

	return metrics
}

// analyzeCompliance evaluates regulatory and compliance status
func (bca *BusinessCreditAnalyzer) analyzeCompliance(ctx sdk.Context, businessID string) types.ComplianceProfile {
	complianceData := bca.keeper.GetComplianceData(ctx, businessID)
	
	profile := types.ComplianceProfile{
		BusinessID: businessID,
	}

	// Regulatory compliance
	profile.TaxCompliance = complianceData.TaxComplianceStatus
	profile.LaborCompliance = complianceData.LaborComplianceStatus
	profile.EnvironmentalCompliance = complianceData.EnvironmentalComplianceStatus
	profile.IndustryCompliance = complianceData.IndustrySpecificCompliance

	// License and permits
	profile.LicenseStatus = complianceData.LicenseStatus
	profile.PermitStatus = complianceData.PermitStatus

	// Audit history
	profile.AuditHistory = complianceData.AuditHistory
	profile.ComplianceScore = bca.calculateComplianceScore(complianceData)

	// Legal issues
	profile.PendingLitigation = complianceData.PendingLitigation
	profile.RegulatoryPenalties = complianceData.RegulatoryPenalties

	return profile
}

// analyzeDigitalPresence evaluates digital footprint and online presence
func (bca *BusinessCreditAnalyzer) analyzeDigitalPresence(ctx sdk.Context, businessID string) types.DigitalFootprint {
	digitalData := bca.keeper.GetDigitalData(ctx, businessID)
	
	footprint := types.DigitalFootprint{
		BusinessID: businessID,
	}

	// Online presence
	footprint.WebsiteQuality = digitalData.WebsiteQualityScore
	footprint.SocialMediaPresence = digitalData.SocialMediaScore
	footprint.OnlineReviews = digitalData.OnlineReviewScore

	// E-commerce metrics
	if !digitalData.OnlineRevenue.IsZero() && !digitalData.TotalRevenue.IsZero() {
		footprint.EcommerceAdoption = digitalData.OnlineRevenue.Amount.ToDec().Quo(digitalData.TotalRevenue.Amount.ToDec())
	}

	// Digital marketing
	footprint.DigitalMarketingEffectiveness = digitalData.DigitalMarketingScore
	
	// Technology adoption
	footprint.TechnologyAdoptionScore = digitalData.TechnologyAdoptionScore
	
	// Cybersecurity
	footprint.CybersecurityScore = digitalData.CybersecurityScore

	return footprint
}

// analyzeSupplyChain evaluates supply chain resilience and management
func (bca *BusinessCreditAnalyzer) analyzeSupplyChain(ctx sdk.Context, businessID string) types.SupplyChainAnalysis {
	supplyChainData := bca.keeper.GetSupplyChainData(ctx, businessID)
	
	analysis := types.SupplyChainAnalysis{
		BusinessID: businessID,
	}

	// Supplier diversity
	if len(supplyChainData.Suppliers) > 0 {
		analysis.SupplierDiversification = bca.calculateSupplierDiversification(supplyChainData.Suppliers)
		analysis.SupplierReliability = bca.calculateSupplierReliability(supplyChainData.SupplierPerformance)
	}

	// Supply chain resilience
	analysis.SupplyChainResilience = bca.assessSupplyChainResilience(supplyChainData)
	
	// Logistics efficiency
	analysis.LogisticsEfficiency = supplyChainData.LogisticsScore
	
	// Inventory management
	analysis.InventoryOptimization = supplyChainData.InventoryOptimizationScore

	return analysis
}

// analyzeTechnologyAdoption evaluates technology integration and digital transformation
func (bca *BusinessCreditAnalyzer) analyzeTechnologyAdoption(ctx sdk.Context, businessID string) types.TechnologyAdoption {
	techData := bca.keeper.GetTechnologyData(ctx, businessID)
	
	adoption := types.TechnologyAdoption{
		BusinessID: businessID,
	}

	// Core technology adoption
	adoption.ERPAdoption = techData.HasERP
	adoption.CRMAdoption = techData.HasCRM
	adoption.CloudAdoption = techData.CloudAdoptionLevel

	// Digital transformation
	adoption.DigitalTransformationStage = techData.DigitalMaturityLevel
	adoption.AutomationLevel = techData.AutomationScore

	// Innovation metrics
	adoption.TechInvestmentRatio = techData.TechInvestmentPercentage
	adoption.DigitalSkillsLevel = techData.EmployeeDigitalSkillsScore

	return adoption
}

// calculateCompositeCreditScore calculates overall credit score
func (bca *BusinessCreditAnalyzer) calculateCompositeCreditScore(profile *BusinessCreditProfile) int64 {
	// Weighted scoring system
	weights := map[string]sdk.Dec{
		"financial":    sdk.NewDecWithPrec(30, 2), // 30%
		"cashflow":     sdk.NewDecWithPrec(20, 2), // 20%
		"industry":     sdk.NewDecWithPrec(10, 2), // 10%
		"market":       sdk.NewDecWithPrec(15, 2), // 15%
		"operational":  sdk.NewDecWithPrec(10, 2), // 10%
		"compliance":   sdk.NewDecWithPrec(10, 2), // 10%
		"digital":      sdk.NewDecWithPrec(5, 2),  // 5%
	}

	// Score each component (0-850 scale)
	financialScore := bca.scoreFinancialMetrics(profile.FinancialMetrics)
	cashflowScore := bca.scoreCashFlowAnalysis(profile.CashFlowAnalysis)
	industryScore := bca.scoreIndustryAnalysis(profile.IndustryAnalysis)
	marketScore := bca.scoreMarketPosition(profile.MarketPosition)
	operationalScore := bca.scoreOperationalMetrics(profile.OperationalMetrics)
	complianceScore := bca.scoreComplianceProfile(profile.ComplianceProfile)
	digitalScore := bca.scoreDigitalFootprint(profile.DigitalFootprint)

	// Calculate weighted average
	totalScore := sdk.ZeroDec()
	totalScore = totalScore.Add(financialScore.Mul(weights["financial"]))
	totalScore = totalScore.Add(cashflowScore.Mul(weights["cashflow"]))
	totalScore = totalScore.Add(industryScore.Mul(weights["industry"]))
	totalScore = totalScore.Add(marketScore.Mul(weights["market"]))
	totalScore = totalScore.Add(operationalScore.Mul(weights["operational"]))
	totalScore = totalScore.Add(complianceScore.Mul(weights["compliance"]))
	totalScore = totalScore.Add(digitalScore.Mul(weights["digital"]))

	return totalScore.TruncateInt64()
}

// determineRiskCategory categorizes businesses based on credit score
func (bca *BusinessCreditAnalyzer) determineRiskCategory(score int64) string {
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

// calculateMaxLoanEligibility determines maximum loan amount
func (bca *BusinessCreditAnalyzer) calculateMaxLoanEligibility(ctx sdk.Context, profile *BusinessCreditProfile) sdk.Coin {
	params := bca.keeper.GetParams(ctx)
	
	// Base eligibility from revenue
	baseEligibility := profile.FinancialMetrics.AnnualRevenue.Amount.Mul(params.MaxLoanToRevenueRatio).Quo(sdk.NewInt(100))

	// Risk-based multiplier
	var riskMultiplier sdk.Dec
	switch profile.RiskCategory {
	case "EXCELLENT":
		riskMultiplier = sdk.NewDecWithPrec(15, 1) // 1.5x
	case "GOOD":
		riskMultiplier = sdk.NewDecWithPrec(12, 1) // 1.2x
	case "FAIR":
		riskMultiplier = sdk.NewDecWithPrec(10, 1) // 1.0x
	case "POOR":
		riskMultiplier = sdk.NewDecWithPrec(7, 1)  // 0.7x
	case "HIGH_RISK":
		riskMultiplier = sdk.NewDecWithPrec(4, 1)  // 0.4x
	}

	maxAmount := baseEligibility.ToDec().Mul(riskMultiplier).TruncateInt()

	// Apply absolute caps
	if maxAmount.GT(params.MaxLoanAmount.Amount) {
		maxAmount = params.MaxLoanAmount.Amount
	}
	if maxAmount.LT(params.MinLoanAmount.Amount) {
		maxAmount = params.MinLoanAmount.Amount
	}

	return sdk.NewCoin(params.MaxLoanAmount.Denom, maxAmount)
}

// Helper functions for scoring individual components
func (bca *BusinessCreditAnalyzer) scoreFinancialMetrics(metrics types.FinancialMetrics) sdk.Dec {
	score := sdk.NewDec(500) // Base score

	// Revenue growth
	if metrics.RevenueGrowthRate.GT(sdk.NewDecWithPrec(2, 1)) { // >20%
		score = score.Add(sdk.NewDec(100))
	} else if metrics.RevenueGrowthRate.GT(sdk.NewDecWithPrec(1, 1)) { // >10%
		score = score.Add(sdk.NewDec(50))
	} else if metrics.RevenueGrowthRate.IsNegative() {
		score = score.Sub(sdk.NewDec(100))
	}

	// Profitability
	if metrics.ProfitMargin.GT(sdk.NewDecWithPrec(15, 2)) { // >15%
		score = score.Add(sdk.NewDec(100))
	} else if metrics.ProfitMargin.GT(sdk.NewDecWithPrec(10, 2)) { // >10%
		score = score.Add(sdk.NewDec(50))
	} else if metrics.ProfitMargin.IsNegative() {
		score = score.Sub(sdk.NewDec(150))
	}

	// Liquidity
	if metrics.CurrentRatio.GT(sdk.NewDecWithPrec(15, 1)) { // >1.5
		score = score.Add(sdk.NewDec(75))
	} else if metrics.CurrentRatio.LT(sdk.OneDec()) { // <1.0
		score = score.Sub(sdk.NewDec(100))
	}

	// Leverage
	if metrics.DebtToEquityRatio.LT(sdk.NewDecWithPrec(5, 1)) { // <0.5
		score = score.Add(sdk.NewDec(50))
	} else if metrics.DebtToEquityRatio.GT(sdk.NewDec(2)) { // >2.0
		score = score.Sub(sdk.NewDec(100))
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
// - scoreCashFlowAnalysis
// - scoreIndustryAnalysis
// - scoreMarketPosition
// - scoreOperationalMetrics
// - scoreComplianceProfile
// - scoreDigitalFootprint

// calculateRevenueGrowthRate calculates year-over-year revenue growth
func (bca *BusinessCreditAnalyzer) calculateRevenueGrowthRate(revenues []sdk.Coin) sdk.Dec {
	if len(revenues) < 2 {
		return sdk.ZeroDec()
	}

	latest := revenues[len(revenues)-1].Amount
	previous := revenues[len(revenues)-2].Amount

	if previous.IsZero() {
		return sdk.ZeroDec()
	}

	growth := latest.Sub(previous).ToDec().Quo(previous.ToDec())
	return growth
}

// Additional helper methods would include all the calculation functions
// referenced in the analysis methods above