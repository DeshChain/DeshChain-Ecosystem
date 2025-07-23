package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/namo/x/vyavasayamitra/types"
)

// Business data management functions for VyavasayaMitra keeper

// GetBusinessCreditAnalyzer returns a new business credit analyzer
func (k Keeper) GetBusinessCreditAnalyzer() *BusinessCreditAnalyzer {
	return NewBusinessCreditAnalyzer(k)
}

// GetBusinessLoanProcessor returns a new business loan processor
func (k Keeper) GetBusinessLoanProcessor() *BusinessLoanProcessor {
	return NewBusinessLoanProcessor(k)
}

// ProcessComprehensiveBusinessLoan processes a complete business loan from application to disbursement
func (k Keeper) ProcessComprehensiveBusinessLoan(ctx sdk.Context, applicationID string) error {
	loanProcessor := k.GetBusinessLoanProcessor()
	
	// Process business loan application
	loan, err := loanProcessor.ProcessBusinessLoanApplication(ctx, applicationID)
	if err != nil {
		return fmt.Errorf("failed to process business loan application: %w", err)
	}
	
	// Auto-disburse if all conditions are met
	if loan.Status == types.BusinessLoanStatusApproved && 
	   (!loan.CollateralRequired || k.isCollateralVerified(ctx, loan.LoanID)) &&
	   (!loan.GuaranteeRequired || k.isGuaranteeVerified(ctx, loan.LoanID)) &&
	   loan.ComplianceChecks.OverallStatus == "PASSED" {
		err = loanProcessor.DisburseBusinessLoan(ctx, loan.LoanID)
		if err != nil {
			k.Logger(ctx).Error("Failed to auto-disburse business loan", "loan_id", loan.LoanID, "error", err)
		}
	}
	
	return nil
}

// Additional keeper methods for business lending data management

// SetProcessedBusinessLoan stores a processed business loan
func (k Keeper) SetProcessedBusinessLoan(ctx sdk.Context, loan ProcessedBusinessLoan) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&loan)
	store.Set(types.GetProcessedBusinessLoanKey(loan.LoanID), bz)
}

// GetProcessedBusinessLoan retrieves a processed business loan
func (k Keeper) GetProcessedBusinessLoan(ctx sdk.Context, loanID string) (ProcessedBusinessLoan, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetProcessedBusinessLoanKey(loanID))
	if bz == nil {
		return ProcessedBusinessLoan{}, false
	}
	var loan ProcessedBusinessLoan
	k.cdc.MustUnmarshal(bz, &loan)
	return loan, true
}

// Data management functions needed by business credit analyzer

// GetBusinessInfo retrieves comprehensive business information
func (k Keeper) GetBusinessInfo(ctx sdk.Context, businessID string) (types.BusinessInfo, bool) {
	profile, found := k.GetBusinessProfile(ctx, businessID)
	if !found {
		return types.BusinessInfo{}, false
	}
	
	// Convert business profile to business info
	businessInfo := types.BusinessInfo{
		BusinessID:       businessID,
		BusinessName:     profile.BusinessName,
		BusinessType:     profile.BusinessType,
		Industry:         profile.Industry,
		Category:         profile.Category,
		YearsInBusiness:  profile.YearsInBusiness,
		Location:         profile.RegisteredAddress,
		ActiveLoans:      profile.ActiveLoans,
		EstablishedDate:  profile.EstablishedDate,
		OperatingLocations: []string{profile.RegisteredAddress},
		EstimatedMarketShare: sdk.ZeroDec(), // Would need market data
	}
	
	return businessInfo, true
}

// GetFinancialData retrieves financial data for credit analysis
func (k Keeper) GetFinancialData(ctx sdk.Context, businessID string) types.FinancialData {
	// Placeholder implementation - in real system this would fetch from financial data store
	return types.FinancialData{
		BusinessID:           businessID,
		YearlyRevenues:       []sdk.Coin{sdk.NewCoin("dinr", sdk.NewInt(1000000))},
		MonthlyRevenues:      []sdk.Coin{sdk.NewCoin("dinr", sdk.NewInt(83333))},
		YearlyProfits:        []sdk.Coin{sdk.NewCoin("dinr", sdk.NewInt(150000))},
		CurrentAssets:        sdk.NewCoin("dinr", sdk.NewInt(500000)),
		CurrentLiabilities:   sdk.NewCoin("dinr", sdk.NewInt(300000)),
		TotalAssets:          sdk.NewCoin("dinr", sdk.NewInt(2000000)),
		TotalLiabilities:     sdk.NewCoin("dinr", sdk.NewInt(800000)),
		Inventory:            sdk.NewCoin("dinr", sdk.NewInt(200000)),
		InterestExpense:      sdk.NewCoin("dinr", sdk.NewInt(50000)),
	}
}

// GetCashFlowData retrieves cash flow data for analysis
func (k Keeper) GetCashFlowData(ctx sdk.Context, businessID string) types.CashFlowData {
	// Placeholder implementation
	return types.CashFlowData{
		BusinessID:           businessID,
		MonthlyCashFlows:     []sdk.Coin{sdk.NewCoin("dinr", sdk.NewInt(50000))},
		InventoryTurnover:    sdk.NewDec(6),
		ReceivablesTurnover:  sdk.NewDec(8),
		PayablesTurnover:     sdk.NewDec(12),
		CapitalExpenditures:  sdk.NewCoin("dinr", sdk.NewInt(100000)),
		CashReserves:         sdk.NewCoin("dinr", sdk.NewInt(300000)),
	}
}

// GetIndustryBenchmarks retrieves industry-specific benchmarks
func (k Keeper) GetIndustryBenchmarks(ctx sdk.Context, industry string) *types.IndustryBenchmarks {
	// Placeholder implementation with default benchmarks
	return &types.IndustryBenchmarks{
		Industry:             industry,
		AverageGrowthRate:    sdk.NewDecWithPrec(15, 2), // 15%
		AverageProfitMargin:  sdk.NewDecWithPrec(12, 2), // 12%
		RiskLevel:           "MEDIUM",
		CyclicalityIndex:    sdk.NewDecWithPrec(5, 1),   // 0.5
		RegulatoryComplexity: "MEDIUM",
	}
}

// GetCustomerData retrieves customer data for analysis
func (k Keeper) GetCustomerData(ctx sdk.Context, businessID string) types.CustomerData {
	// Placeholder implementation
	return types.CustomerData{
		BusinessID:           businessID,
		TopCustomers:         []types.CustomerInfo{},
		CustomerRetentionData: types.RetentionData{},
	}
}

// GetSupplierData retrieves supplier data for analysis  
func (k Keeper) GetSupplierData(ctx sdk.Context, businessID string) types.SupplierData {
	// Placeholder implementation
	return types.SupplierData{
		BusinessID:    businessID,
		KeySuppliers:  []types.SupplierInfo{},
	}
}

// GetOperationalData retrieves operational metrics
func (k Keeper) GetOperationalData(ctx sdk.Context, businessID string) types.OperationalData {
	// Placeholder implementation
	return types.OperationalData{
		BusinessID:                   businessID,
		NumberOfEmployees:            sdk.NewInt(25),
		MonthlyProduction:            []sdk.Coin{sdk.NewCoin("units", sdk.NewInt(1000))},
		QualityScore:                 sdk.NewDec(85),
		CustomerSatisfactionScore:    sdk.NewDec(90),
		TotalCosts:                   sdk.NewCoin("dinr", sdk.NewInt(800000)),
		TotalRevenue:                 sdk.NewCoin("dinr", sdk.NewInt(1000000)),
		InventoryLevels:              []sdk.Coin{sdk.NewCoin("dinr", sdk.NewInt(200000))},
		EmployeeRetentionRate:        sdk.NewDecWithPrec(85, 2),
		ProductivityHistory:          []sdk.Dec{sdk.NewDec(100)},
	}
}

// GetComplianceData retrieves compliance information
func (k Keeper) GetComplianceData(ctx sdk.Context, businessID string) types.ComplianceData {
	profile, _ := k.GetBusinessProfile(ctx, businessID)
	
	return types.ComplianceData{
		BusinessID:                   businessID,
		TaxComplianceStatus:          "COMPLIANT",
		LaborComplianceStatus:        "COMPLIANT",
		EnvironmentalComplianceStatus: "COMPLIANT",
		IndustrySpecificCompliance:   "COMPLIANT",
		LicenseStatus:                "VALID",
		PermitStatus:                 "VALID",
		AuditHistory:                 []types.AuditRecord{},
		PendingLitigation:            []types.LegalCase{},
		RegulatoryPenalties:          []types.Penalty{},
	}
}

// GetDigitalData retrieves digital presence data
func (k Keeper) GetDigitalData(ctx sdk.Context, businessID string) types.DigitalData {
	// Placeholder implementation
	return types.DigitalData{
		BusinessID:              businessID,
		WebsiteQualityScore:     sdk.NewDec(75),
		SocialMediaScore:        sdk.NewDec(60),
		OnlineReviewScore:       sdk.NewDec(85),
		OnlineRevenue:           sdk.NewCoin("dinr", sdk.NewInt(300000)),
		TotalRevenue:            sdk.NewCoin("dinr", sdk.NewInt(1000000)),
		DigitalMarketingScore:   sdk.NewDec(70),
		TechnologyAdoptionScore: sdk.NewDec(80),
		CybersecurityScore:      sdk.NewDec(75),
	}
}

// GetSupplyChainData retrieves supply chain information
func (k Keeper) GetSupplyChainData(ctx sdk.Context, businessID string) types.SupplyChainData {
	// Placeholder implementation
	return types.SupplyChainData{
		BusinessID:                businessID,
		Suppliers:                 []types.SupplierInfo{},
		SupplierPerformance:       []types.PerformanceMetric{},
		LogisticsScore:            sdk.NewDec(80),
		InventoryOptimizationScore: sdk.NewDec(75),
	}
}

// GetTechnologyData retrieves technology adoption data
func (k Keeper) GetTechnologyData(ctx sdk.Context, businessID string) types.TechnologyData {
	// Placeholder implementation
	return types.TechnologyData{
		BusinessID:                businessID,
		HasERP:                    true,
		HasCRM:                    false,
		CloudAdoptionLevel:        sdk.NewDec(60),
		DigitalMaturityLevel:      "INTERMEDIATE",
		AutomationScore:           sdk.NewDec(50),
		TechInvestmentPercentage:  sdk.NewDecWithPrec(5, 2), // 5%
		EmployeeDigitalSkillsScore: sdk.NewDec(70),
	}
}

// Additional helper functions
func (k Keeper) isCollateralVerified(ctx sdk.Context, loanID string) bool {
	// Implementation would check collateral verification status
	return true // Simplified for now
}

func (k Keeper) isGuaranteeVerified(ctx sdk.Context, loanID string) bool {
	// Implementation would check guarantee verification status
	return true // Simplified for now
}

// Additional credit scoring helper functions
func (k Keeper) GetWeatherRiskProfile(ctx sdk.Context, businessID string) (types.WeatherRiskProfile, bool) {
	// Placeholder implementation
	return types.WeatherRiskProfile{
		BusinessID:       businessID,
		DroughtFrequency: 2,
		FloodRisk:        "LOW",
	}, true
}

func (k Keeper) GetFarmerClaimHistory(ctx sdk.Context, businessID string) []types.InsuranceClaim {
	// Placeholder implementation
	return []types.InsuranceClaim{}
}

func (k Keeper) GetCropInsuranceInfo(ctx sdk.Context, cropType string) (types.CropInsuranceInfo, bool) {
	// Placeholder implementation
	return types.CropInsuranceInfo{
		CropType:            cropType,
		AverageYieldPerAcre: sdk.NewDec(25),
		AverageMarketPrice:  sdk.NewDec(2000),
	}, true
}

func (k Keeper) GetAverageYield(ctx sdk.Context, cropType, farmerID string) sdk.Dec {
	// Placeholder implementation
	return sdk.NewDec(25) // 25 quintals per acre
}

func (k Keeper) GetFarmerProfile(ctx sdk.Context, farmerID string) (types.FarmerProfile, bool) {
	// Placeholder implementation for business context
	return types.FarmerProfile{
		FarmerID:       farmerID,
		TotalLandArea:  sdk.NewDec(10), // 10 acres
		IsWomenFarmer:  false,
		TotalLoans:     2,
		DefaultedLoans: 0,
		LandSize:       sdk.NewDec(10),
		ActiveLoans:    1,
	}, true
}

// Credit scoring data structures placeholders
func (k Keeper) calculateComplianceScore(data types.ComplianceData) int64 {
	// Simplified compliance scoring
	score := int64(100)
	if data.TaxComplianceStatus != "COMPLIANT" {
		score -= 20
	}
	if data.LaborComplianceStatus != "COMPLIANT" {
		score -= 15
	}
	if len(data.PendingLitigation) > 0 {
		score -= 25
	}
	return score
}

func (k Keeper) calculateCustomerConcentration(customers []types.CustomerInfo) sdk.Dec {
	// Placeholder implementation
	return sdk.NewDecWithPrec(2, 1) // 20% concentration
}

func (k Keeper) calculateCustomerStickiness(retention types.RetentionData) sdk.Dec {
	// Placeholder implementation
	return sdk.NewDecWithPrec(85, 2) // 85% stickiness
}

func (k Keeper) identifyCompetitiveAdvantages(ctx sdk.Context, businessID string) []string {
	// Placeholder implementation
	return []string{"Strong local presence", "Competitive pricing"}
}

func (k Keeper) calculateSupplierConcentration(suppliers []types.SupplierInfo) sdk.Dec {
	// Placeholder implementation
	return sdk.NewDecWithPrec(3, 1) // 30% concentration
}

func (k Keeper) assessSupplierRelationships(data types.SupplierData) sdk.Dec {
	// Placeholder implementation
	return sdk.NewDec(8) // Good relationships (8/10)
}

func (k Keeper) assessBrandStrength(ctx sdk.Context, businessID string) int64 {
	// Placeholder implementation
	return 7 // Strong brand (7/10)
}

func (k Keeper) calculateGeographicDiversification(locations []string) sdk.Dec {
	// Placeholder implementation
	return sdk.NewDecWithPrec(5, 1) // 50% diversification
}

func (k Keeper) calculateInnovationIndex(ctx sdk.Context, businessID string) sdk.Dec {
	// Placeholder implementation
	return sdk.NewDec(6) // Moderate innovation (6/10)
}

func (k Keeper) calculateInventoryTurnover(levels []sdk.Coin) sdk.Dec {
	// Placeholder implementation
	return sdk.NewDec(6) // 6 times per year
}

func (k Keeper) calculateEmployeeProductivityGrowth(history []sdk.Dec) sdk.Dec {
	// Placeholder implementation
	return sdk.NewDecWithPrec(5, 2) // 5% growth
}

func (k Keeper) calculateSupplierDiversification(suppliers []types.SupplierInfo) sdk.Dec {
	// Placeholder implementation
	return sdk.NewDecWithPrec(7, 1) // 70% diversification
}

func (k Keeper) calculateSupplierReliability(performance []types.PerformanceMetric) sdk.Dec {
	// Placeholder implementation
	return sdk.NewDecWithPrec(85, 2) // 85% reliability
}

func (k Keeper) assessSupplyChainResilience(data types.SupplyChainData) sdk.Dec {
	// Placeholder implementation
	return sdk.NewDec(8) // High resilience (8/10)
}

func (k Keeper) getIndustryRiskFactors(industry string) []string {
	// Industry-specific risk factors
	riskMap := map[string][]string{
		"MANUFACTURING": {"Supply chain disruption", "Regulatory changes"},
		"EXPORT":        {"Currency fluctuation", "Trade policy changes"},
		"TECHNOLOGY":    {"Rapid obsolescence", "Talent shortage"},
		"RETAIL":        {"Market competition", "Consumer preference shifts"},
	}
	
	if risks, found := riskMap[industry]; found {
		return risks
	}
	return []string{"General market risk"}
}

func (k Keeper) assessCompetitionLevel(ctx sdk.Context, industry, location string) string {
	// Simplified competition assessment
	competitiveIndustries := []string{"RETAIL", "TECHNOLOGY", "SERVICES"}
	for _, competitive := range competitiveIndustries {
		if industry == competitive {
			return "HIGH"
		}
	}
	return "MEDIUM"
}

func (k Keeper) getMarketOutlook(ctx sdk.Context, industry string) string {
	// Simplified market outlook
	growingIndustries := []string{"TECHNOLOGY", "RENEWABLE_ENERGY", "HEALTHCARE"}
	for _, growing := range growingIndustries {
		if industry == growing {
			return "POSITIVE"
		}
	}
	return "STABLE"
}