package keeper

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"time"

	"cosmossdk.io/core/store"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/deshchain/deshchain/x/tradefinance/types"
)

// ConcentrationRiskManager handles concentration risk assessment and management
type ConcentrationRiskManager struct {
	keeper         *Keeper
	riskLimits     ConcentrationLimits
	riskMetrics    RiskMetrics
	earlyWarning   EarlyWarningSystem
}

// NewConcentrationRiskManager creates a new concentration risk manager
func NewConcentrationRiskManager(k *Keeper) *ConcentrationRiskManager {
	return &ConcentrationRiskManager{
		keeper:       k,
		riskLimits:   initializeConcentrationLimits(),
		riskMetrics:  initializeRiskMetrics(),
		earlyWarning: initializeEarlyWarningSystem(),
	}
}

// Core concentration risk structures

type ConcentrationRiskReport struct {
	ReportID               string                      `json:"report_id"`
	ReportDate             time.Time                   `json:"report_date"`
	InstitutionID          string                      `json:"institution_id"`
	CreditConcentrations   CreditConcentrationAnalysis `json:"credit_concentrations"`
	SectorConcentrations   SectorConcentrationAnalysis `json:"sector_concentrations"`
	GeographicConcentrations GeographicConcentrationAnalysis `json:"geographic_concentrations"`
	CounterpartyConcentrations CounterpartyConcentrationAnalysis `json:"counterparty_concentrations"`
	ProductConcentrations  ProductConcentrationAnalysis `json:"product_concentrations"`
	CollateralConcentrations CollateralConcentrationAnalysis `json:"collateral_concentrations"`
	MaturityConcentrations MaturityConcentrationAnalysis `json:"maturity_concentrations"`
	LargeExposures         LargeExposureAnalysis      `json:"large_exposures"`
	ConnectedLending       ConnectedLendingAnalysis   `json:"connected_lending"`
	RiskMetrics            ConcentrationRiskMetrics   `json:"risk_metrics"`
	ComplianceStatus       ConcentrationCompliance    `json:"compliance_status"`
	MitigationStrategies   []MitigationStrategy       `json:"mitigation_strategies"`
	StressTesting          ConcentrationStressTest    `json:"stress_testing"`
	Recommendations        []RiskRecommendation       `json:"recommendations"`
	Metadata               map[string]interface{}     `json:"metadata"`
}

type CreditConcentrationAnalysis struct {
	SingleBorrowerLimit    LimitAnalysis              `json:"single_borrower_limit"`
	GroupBorrowerLimit     LimitAnalysis              `json:"group_borrower_limit"`
	Top20Exposures         []ExposureDetail           `json:"top_20_exposures"`
	ConcentrationRatio     sdk.Dec                    `json:"concentration_ratio"`
	HerfindahlIndex        sdk.Dec                    `json:"herfindahl_index"`
	GiniCoefficient        sdk.Dec                    `json:"gini_coefficient"`
	LorenzCurve            []LorenzPoint              `json:"lorenz_curve"`
	RiskDistribution       RiskGradeDistribution      `json:"risk_distribution"`
	NPAConcentration       NPAConcentrationAnalysis   `json:"npa_concentration"`
}

type SectorConcentrationAnalysis struct {
	SectorExposures        []SectorExposure           `json:"sector_exposures"`
	SectorLimits           []SectorLimit              `json:"sector_limits"`
	HighRiskSectors        []string                   `json:"high_risk_sectors"`
	SectorCorrelations     []SectorCorrelation        `json:"sector_correlations"`
	CyclicalityAnalysis    CyclicalityMetrics         `json:"cyclicality_analysis"`
	IndustryRiskScores     map[string]RiskScore       `json:"industry_risk_scores"`
}

type GeographicConcentrationAnalysis struct {
	CountryExposures       []CountryExposure          `json:"country_exposures"`
	RegionalDistribution   []RegionalExposure         `json:"regional_distribution"`
	EmergingMarketExposure sdk.Coin                   `json:"emerging_market_exposure"`
	CrossBorderExposure    sdk.Coin                   `json:"cross_border_exposure"`
	PoliticalRiskExposure  []PoliticalRiskExposure    `json:"political_risk_exposure"`
	CurrencyConcentration  []CurrencyExposure         `json:"currency_concentration"`
}

type CounterpartyConcentrationAnalysis struct {
	TopCounterparties      []CounterpartyExposure     `json:"top_counterparties"`
	InterBankExposures     sdk.Coin                   `json:"inter_bank_exposures"`
	FinancialInstitutions  sdk.Coin                   `json:"financial_institutions"`
	CorporateGroups        []CorporateGroupExposure   `json:"corporate_groups"`
	RelatedPartyExposures  sdk.Coin                   `json:"related_party_exposures"`
	SystemicCounterparties []SystemicCounterparty     `json:"systemic_counterparties"`
}

type LargeExposureAnalysis struct {
	LargeExposureLimit     sdk.Dec                    `json:"large_exposure_limit"`
	NumberOfLargeExposures int                        `json:"number_of_large_exposures"`
	TotalLargeExposures    sdk.Coin                   `json:"total_large_exposures"`
	LargeExposureDetails   []LargeExposureDetail      `json:"large_exposure_details"`
	BreachAnalysis         []LimitBreach              `json:"breach_analysis"`
	TrendAnalysis          LargeExposureTrend         `json:"trend_analysis"`
}

// Core calculation functions

// AnalyzeConcentrationRisk performs comprehensive concentration risk analysis
func (crm *ConcentrationRiskManager) AnalyzeConcentrationRisk(ctx context.Context, institutionID string) (*ConcentrationRiskReport, error) {
	report := &ConcentrationRiskReport{
		ReportID:      fmt.Sprintf("CONC_RISK_%s_%d", institutionID, time.Now().Unix()),
		ReportDate:    time.Now(),
		InstitutionID: institutionID,
		Metadata:      make(map[string]interface{}),
	}

	// Get institution data
	institutionData, err := crm.getInstitutionData(ctx, institutionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get institution data: %w", err)
	}

	// Analyze credit concentrations
	creditConc, err := crm.analyzeCreditConcentrations(ctx, institutionData)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze credit concentrations: %w", err)
	}
	report.CreditConcentrations = *creditConc

	// Analyze sector concentrations
	sectorConc, err := crm.analyzeSectorConcentrations(ctx, institutionData)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze sector concentrations: %w", err)
	}
	report.SectorConcentrations = *sectorConc

	// Analyze geographic concentrations
	geoConc, err := crm.analyzeGeographicConcentrations(ctx, institutionData)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze geographic concentrations: %w", err)
	}
	report.GeographicConcentrations = *geoConc

	// Analyze counterparty concentrations
	counterpartyConc, err := crm.analyzeCounterpartyConcentrations(ctx, institutionData)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze counterparty concentrations: %w", err)
	}
	report.CounterpartyConcentrations = *counterpartyConc

	// Analyze product concentrations
	productConc, err := crm.analyzeProductConcentrations(ctx, institutionData)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze product concentrations: %w", err)
	}
	report.ProductConcentrations = *productConc

	// Analyze collateral concentrations
	collateralConc, err := crm.analyzeCollateralConcentrations(ctx, institutionData)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze collateral concentrations: %w", err)
	}
	report.CollateralConcentrations = *collateralConc

	// Analyze maturity concentrations
	maturityConc, err := crm.analyzeMaturityConcentrations(ctx, institutionData)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze maturity concentrations: %w", err)
	}
	report.MaturityConcentrations = *maturityConc

	// Analyze large exposures
	largeExp, err := crm.analyzeLargeExposures(ctx, institutionData)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze large exposures: %w", err)
	}
	report.LargeExposures = *largeExp

	// Analyze connected lending
	connectedLending, err := crm.analyzeConnectedLending(ctx, institutionData)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze connected lending: %w", err)
	}
	report.ConnectedLending = *connectedLending

	// Calculate risk metrics
	riskMetrics := crm.calculateRiskMetrics(report)
	report.RiskMetrics = *riskMetrics

	// Check compliance
	compliance := crm.checkCompliance(report)
	report.ComplianceStatus = *compliance

	// Run stress tests
	stressTest := crm.runConcentrationStressTest(ctx, institutionData, report)
	report.StressTesting = *stressTest

	// Generate mitigation strategies
	mitigations := crm.generateMitigationStrategies(report)
	report.MitigationStrategies = mitigations

	// Generate recommendations
	recommendations := crm.generateRecommendations(report)
	report.Recommendations = recommendations

	// Store report
	if err := crm.storeConcentrationRiskReport(ctx, *report); err != nil {
		return report, fmt.Errorf("failed to store report: %w", err)
	}

	// Emit event
	crm.emitConcentrationRiskEvent(ctx, report)

	return report, nil
}

// Credit concentration analysis

func (crm *ConcentrationRiskManager) analyzeCreditConcentrations(ctx context.Context, data InstitutionData) (*CreditConcentrationAnalysis, error) {
	analysis := &CreditConcentrationAnalysis{}

	// Get all credit exposures
	exposures := crm.getCreditExposures(data)
	
	// Sort by exposure amount
	sort.Slice(exposures, func(i, j int) bool {
		return exposures[i].Amount.IsGT(exposures[j].Amount)
	})

	// Single borrower limit analysis
	capitalBase := data.CapitalBase
	singleBorrowerLimit := capitalBase.Amount.Mul(sdk.NewInt(int64(crm.riskLimits.SingleBorrowerLimit.MustFloat64() * 100))).Quo(sdk.NewInt(100))
	
	breaches := []LimitBreachDetail{}
	for _, exp := range exposures {
		if exp.Amount.Amount.GT(singleBorrowerLimit) {
			breaches = append(breaches, LimitBreachDetail{
				ExposureID:   exp.ID,
				BorrowerName: exp.BorrowerName,
				ExposureAmount: exp.Amount,
				LimitAmount:   sdk.NewCoin(exp.Amount.Denom, singleBorrowerLimit),
				ExcessAmount:  sdk.NewCoin(exp.Amount.Denom, exp.Amount.Amount.Sub(singleBorrowerLimit)),
				BreachPercentage: sdk.NewDecFromInt(exp.Amount.Amount.Sub(singleBorrowerLimit)).Quo(sdk.NewDecFromInt(singleBorrowerLimit)),
			})
		}
	}

	analysis.SingleBorrowerLimit = LimitAnalysis{
		LimitType:        "SINGLE_BORROWER",
		LimitPercentage:  crm.riskLimits.SingleBorrowerLimit,
		LimitAmount:      sdk.NewCoin(capitalBase.Denom, singleBorrowerLimit),
		CurrentExposure:  crm.getMaxSingleExposure(exposures),
		UtilizationRatio: crm.calculateUtilizationRatio(crm.getMaxSingleExposure(exposures), sdk.NewCoin(capitalBase.Denom, singleBorrowerLimit)),
		Breaches:         breaches,
		IsCompliant:      len(breaches) == 0,
	}

	// Group borrower limit analysis
	groupExposures := crm.aggregateGroupExposures(exposures)
	groupBorrowerLimit := capitalBase.Amount.Mul(sdk.NewInt(int64(crm.riskLimits.GroupBorrowerLimit.MustFloat64() * 100))).Quo(sdk.NewInt(100))
	
	groupBreaches := []LimitBreachDetail{}
	for _, grpExp := range groupExposures {
		if grpExp.TotalExposure.Amount.GT(groupBorrowerLimit) {
			groupBreaches = append(groupBreaches, LimitBreachDetail{
				ExposureID:   grpExp.GroupID,
				BorrowerName: grpExp.GroupName,
				ExposureAmount: grpExp.TotalExposure,
				LimitAmount:   sdk.NewCoin(grpExp.TotalExposure.Denom, groupBorrowerLimit),
				ExcessAmount:  sdk.NewCoin(grpExp.TotalExposure.Denom, grpExp.TotalExposure.Amount.Sub(groupBorrowerLimit)),
				BreachPercentage: sdk.NewDecFromInt(grpExp.TotalExposure.Amount.Sub(groupBorrowerLimit)).Quo(sdk.NewDecFromInt(groupBorrowerLimit)),
			})
		}
	}

	analysis.GroupBorrowerLimit = LimitAnalysis{
		LimitType:        "GROUP_BORROWER",
		LimitPercentage:  crm.riskLimits.GroupBorrowerLimit,
		LimitAmount:      sdk.NewCoin(capitalBase.Denom, groupBorrowerLimit),
		CurrentExposure:  crm.getMaxGroupExposure(groupExposures),
		UtilizationRatio: crm.calculateUtilizationRatio(crm.getMaxGroupExposure(groupExposures), sdk.NewCoin(capitalBase.Denom, groupBorrowerLimit)),
		Breaches:         groupBreaches,
		IsCompliant:      len(groupBreaches) == 0,
	}

	// Top 20 exposures
	top20Count := min(20, len(exposures))
	for i := 0; i < top20Count; i++ {
		analysis.Top20Exposures = append(analysis.Top20Exposures, ExposureDetail{
			Rank:             i + 1,
			ExposureID:       exposures[i].ID,
			BorrowerName:     exposures[i].BorrowerName,
			ExposureAmount:   exposures[i].Amount,
			PercentOfCapital: sdk.NewDecFromInt(exposures[i].Amount.Amount).Quo(sdk.NewDecFromInt(capitalBase.Amount)),
			RiskGrade:        exposures[i].RiskGrade,
			Sector:           exposures[i].Sector,
			Geography:        exposures[i].Country,
			MaturityDate:     exposures[i].MaturityDate,
		})
	}

	// Calculate concentration metrics
	analysis.ConcentrationRatio = crm.calculateConcentrationRatio(exposures, capitalBase)
	analysis.HerfindahlIndex = crm.calculateHerfindahlIndex(exposures, data.TotalExposure)
	analysis.GiniCoefficient = crm.calculateGiniCoefficient(exposures)
	analysis.LorenzCurve = crm.calculateLorenzCurve(exposures)

	// Risk distribution analysis
	analysis.RiskDistribution = crm.analyzeRiskDistribution(exposures)

	// NPA concentration analysis
	analysis.NPAConcentration = crm.analyzeNPAConcentration(exposures)

	return analysis, nil
}

// Sector concentration analysis

func (crm *ConcentrationRiskManager) analyzeSectorConcentrations(ctx context.Context, data InstitutionData) (*SectorConcentrationAnalysis, error) {
	analysis := &SectorConcentrationAnalysis{
		SectorExposures:     []SectorExposure{},
		SectorLimits:        []SectorLimit{},
		HighRiskSectors:     []string{},
		SectorCorrelations:  []SectorCorrelation{},
		IndustryRiskScores:  make(map[string]RiskScore),
	}

	// Aggregate exposures by sector
	sectorMap := make(map[string]*SectorExposure)
	for _, exp := range data.CreditExposures {
		if sector, exists := sectorMap[exp.Sector]; exists {
			sector.TotalExposure = sector.TotalExposure.Add(exp.Amount)
			sector.NumberOfBorrowers++
			sector.Exposures = append(sector.Exposures, exp)
		} else {
			sectorMap[exp.Sector] = &SectorExposure{
				SectorCode:        exp.Sector,
				SectorName:        crm.getSectorName(exp.Sector),
				TotalExposure:     exp.Amount,
				PercentOfTotal:    sdk.ZeroDec(),
				NumberOfBorrowers: 1,
				AverageRiskGrade:  exp.RiskGrade,
				NPAAmount:         sdk.NewCoin(exp.Amount.Denom, sdk.ZeroInt()),
				Exposures:         []CreditExposure{exp},
			}
		}
	}

	// Calculate percentages and analyze each sector
	totalExposure := data.TotalExposure
	for sectorCode, sectorExp := range sectorMap {
		sectorExp.PercentOfTotal = sdk.NewDecFromInt(sectorExp.TotalExposure.Amount).Quo(sdk.NewDecFromInt(totalExposure.Amount))
		
		// Calculate average risk grade
		sectorExp.AverageRiskGrade = crm.calculateAverageRiskGrade(sectorExp.Exposures)
		
		// Calculate sector NPA
		sectorExp.NPAAmount = crm.calculateSectorNPA(sectorExp.Exposures)
		sectorExp.NPAPercentage = sdk.NewDecFromInt(sectorExp.NPAAmount.Amount).Quo(sdk.NewDecFromInt(sectorExp.TotalExposure.Amount))
		
		// Check sector limits
		sectorLimit := crm.getSectorLimit(sectorCode)
		sectorExp.LimitUtilization = sectorExp.PercentOfTotal.Quo(sectorLimit.MaxExposure)
		sectorExp.IsHighRisk = crm.isHighRiskSector(sectorCode)
		
		if sectorExp.IsHighRisk {
			analysis.HighRiskSectors = append(analysis.HighRiskSectors, sectorCode)
		}
		
		analysis.SectorExposures = append(analysis.SectorExposures, *sectorExp)
		
		// Add sector limit info
		analysis.SectorLimits = append(analysis.SectorLimits, SectorLimit{
			SectorCode:       sectorCode,
			MaxExposure:      sectorLimit.MaxExposure,
			CurrentExposure:  sectorExp.PercentOfTotal,
			AvailableLimit:   sdk.MaxDec(sdk.ZeroDec(), sectorLimit.MaxExposure.Sub(sectorExp.PercentOfTotal)),
			IsBreached:       sectorExp.PercentOfTotal.GT(sectorLimit.MaxExposure),
		})
		
		// Calculate industry risk score
		analysis.IndustryRiskScores[sectorCode] = crm.calculateIndustryRiskScore(sectorExp)
	}

	// Sort sectors by exposure
	sort.Slice(analysis.SectorExposures, func(i, j int) bool {
		return analysis.SectorExposures[i].TotalExposure.IsGT(analysis.SectorExposures[j].TotalExposure)
	})

	// Calculate sector correlations
	analysis.SectorCorrelations = crm.calculateSectorCorrelations(analysis.SectorExposures)

	// Analyze cyclicality
	analysis.CyclicalityAnalysis = crm.analyzeCyclicality(analysis.SectorExposures)

	return analysis, nil
}

// Geographic concentration analysis

func (crm *ConcentrationRiskManager) analyzeGeographicConcentrations(ctx context.Context, data InstitutionData) (*GeographicConcentrationAnalysis, error) {
	analysis := &GeographicConcentrationAnalysis{
		CountryExposures:      []CountryExposure{},
		RegionalDistribution:  []RegionalExposure{},
		PoliticalRiskExposure: []PoliticalRiskExposure{},
		CurrencyConcentration: []CurrencyExposure{},
	}

	// Aggregate by country
	countryMap := make(map[string]*CountryExposure)
	currencyMap := make(map[string]sdk.Coin)
	
	for _, exp := range data.CreditExposures {
		// Country aggregation
		if country, exists := countryMap[exp.Country]; exists {
			country.TotalExposure = country.TotalExposure.Add(exp.Amount)
			country.NumberOfBorrowers++
		} else {
			countryMap[exp.Country] = &CountryExposure{
				CountryCode:       exp.Country,
				CountryName:       crm.getCountryName(exp.Country),
				TotalExposure:     exp.Amount,
				NumberOfBorrowers: 1,
				CountryRating:     crm.getCountryRating(exp.Country),
				PoliticalRiskScore: crm.getPoliticalRiskScore(exp.Country),
			}
		}
		
		// Currency aggregation
		if amt, exists := currencyMap[exp.Currency]; exists {
			currencyMap[exp.Currency] = amt.Add(exp.Amount)
		} else {
			currencyMap[exp.Currency] = exp.Amount
		}
		
		// Track emerging market exposure
		if crm.isEmergingMarket(exp.Country) {
			analysis.EmergingMarketExposure = analysis.EmergingMarketExposure.Add(exp.Amount)
		}
		
		// Track cross-border exposure
		if exp.Country != data.HomeCountry {
			analysis.CrossBorderExposure = analysis.CrossBorderExposure.Add(exp.Amount)
		}
	}

	// Process country exposures
	for _, country := range countryMap {
		country.PercentOfTotal = sdk.NewDecFromInt(country.TotalExposure.Amount).Quo(sdk.NewDecFromInt(data.TotalExposure.Amount))
		country.ConcentrationRisk = crm.assessCountryConcentrationRisk(country)
		analysis.CountryExposures = append(analysis.CountryExposures, *country)
		
		// Add political risk exposure if high risk
		if country.PoliticalRiskScore > 7 {
			analysis.PoliticalRiskExposure = append(analysis.PoliticalRiskExposure, PoliticalRiskExposure{
				CountryCode:    country.CountryCode,
				ExposureAmount: country.TotalExposure,
				RiskScore:      country.PoliticalRiskScore,
				RiskFactors:    crm.getPoliticalRiskFactors(country.CountryCode),
			})
		}
	}

	// Regional distribution
	analysis.RegionalDistribution = crm.aggregateByRegion(countryMap)

	// Currency concentration
	for currency, amount := range currencyMap {
		analysis.CurrencyConcentration = append(analysis.CurrencyConcentration, CurrencyExposure{
			Currency:       currency,
			ExposureAmount: amount,
			PercentOfTotal: sdk.NewDecFromInt(amount.Amount).Quo(sdk.NewDecFromInt(data.TotalExposure.Amount)),
			ExchangeRisk:   crm.assessExchangeRisk(currency, data.HomeCurrency),
		})
	}

	// Sort by exposure
	sort.Slice(analysis.CountryExposures, func(i, j int) bool {
		return analysis.CountryExposures[i].TotalExposure.IsGT(analysis.CountryExposures[j].TotalExposure)
	})

	return analysis, nil
}

// Large exposure analysis

func (crm *ConcentrationRiskManager) analyzeLargeExposures(ctx context.Context, data InstitutionData) (*LargeExposureAnalysis, error) {
	analysis := &LargeExposureAnalysis{
		LargeExposureLimit:     crm.riskLimits.LargeExposureThreshold,
		NumberOfLargeExposures: 0,
		TotalLargeExposures:    sdk.NewCoin(data.CapitalBase.Denom, sdk.ZeroInt()),
		LargeExposureDetails:   []LargeExposureDetail{},
		BreachAnalysis:         []LimitBreach{},
	}

	// Large exposure threshold
	largeExpThreshold := data.CapitalBase.Amount.Mul(sdk.NewInt(int64(crm.riskLimits.LargeExposureThreshold.MustFloat64() * 100))).Quo(sdk.NewInt(100))
	
	// Aggregate all exposures by counterparty (including groups)
	counterpartyExposures := crm.aggregateAllExposures(data)
	
	for _, exp := range counterpartyExposures {
		if exp.TotalExposure.Amount.GTE(largeExpThreshold) {
			analysis.NumberOfLargeExposures++
			analysis.TotalLargeExposures = analysis.TotalLargeExposures.Add(exp.TotalExposure)
			
			percentOfCapital := sdk.NewDecFromInt(exp.TotalExposure.Amount).Quo(sdk.NewDecFromInt(data.CapitalBase.Amount))
			
			detail := LargeExposureDetail{
				ExposureID:          exp.ID,
				CounterpartyName:    exp.Name,
				ExposureAmount:      exp.TotalExposure,
				PercentOfCapital:    percentOfCapital,
				ExposureType:        exp.Type,
				CreditRating:        exp.Rating,
				CollateralCoverage:  exp.CollateralCoverage,
				NetExposure:         crm.calculateNetExposure(exp),
				ExposureComponents:  exp.Components,
				LastReviewDate:      exp.LastReview,
				NextReviewDate:      exp.NextReview,
			}
			
			// Check for limit breach (25% of capital for non-bank, 100% for interbank)
			var limitPercent sdk.Dec
			if exp.Type == "INTERBANK" {
				limitPercent = sdk.NewDecWithPrec(100, 2) // 100%
			} else {
				limitPercent = sdk.NewDecWithPrec(25, 2) // 25%
			}
			
			if percentOfCapital.GT(limitPercent) {
				analysis.BreachAnalysis = append(analysis.BreachAnalysis, LimitBreach{
					BreachID:        fmt.Sprintf("LARGE_EXP_%s_%d", exp.ID, time.Now().Unix()),
					LimitType:       "LARGE_EXPOSURE",
					CounterpartyID:  exp.ID,
					CounterpartyName: exp.Name,
					LimitAmount:     limitPercent,
					ActualAmount:    percentOfCapital,
					ExcessAmount:    percentOfCapital.Sub(limitPercent),
					BreachDate:      time.Now(),
					Status:          "ACTIVE",
				})
				detail.IsBreaching = true
			}
			
			analysis.LargeExposureDetails = append(analysis.LargeExposureDetails, detail)
		}
	}

	// Trend analysis
	analysis.TrendAnalysis = crm.analyzeLargeExposureTrend(ctx, data.InstitutionID)

	return analysis, nil
}

// Risk metrics calculation

func (crm *ConcentrationRiskManager) calculateRiskMetrics(report *ConcentrationRiskReport) *ConcentrationRiskMetrics {
	metrics := &ConcentrationRiskMetrics{
		CalculationDate: time.Now(),
	}

	// Overall concentration score (0-100, higher is riskier)
	creditScore := crm.scoreCreditConcentration(&report.CreditConcentrations)
	sectorScore := crm.scoreSectorConcentration(&report.SectorConcentrations)
	geoScore := crm.scoreGeographicConcentration(&report.GeographicConcentrations)
	largeExpScore := crm.scoreLargeExposures(&report.LargeExposures)
	
	metrics.OverallConcentrationScore = (creditScore*0.3 + sectorScore*0.25 + geoScore*0.25 + largeExpScore*0.2)
	
	// Risk level classification
	if metrics.OverallConcentrationScore >= 80 {
		metrics.RiskLevel = "CRITICAL"
	} else if metrics.OverallConcentrationScore >= 60 {
		metrics.RiskLevel = "HIGH"
	} else if metrics.OverallConcentrationScore >= 40 {
		metrics.RiskLevel = "MEDIUM"
	} else {
		metrics.RiskLevel = "LOW"
	}

	// Key risk indicators
	metrics.KeyRiskIndicators = []KeyRiskIndicator{
		{
			IndicatorName:  "Herfindahl Index",
			CurrentValue:   report.CreditConcentrations.HerfindahlIndex.MustFloat64(),
			ThresholdValue: 0.15,
			Status:         crm.getKRIStatus(report.CreditConcentrations.HerfindahlIndex.MustFloat64(), 0.15, true),
		},
		{
			IndicatorName:  "Top 20 Concentration",
			CurrentValue:   crm.calculateTop20Concentration(report.CreditConcentrations.Top20Exposures),
			ThresholdValue: 0.40,
			Status:         crm.getKRIStatus(crm.calculateTop20Concentration(report.CreditConcentrations.Top20Exposures), 0.40, true),
		},
		{
			IndicatorName:  "Single Sector Maximum",
			CurrentValue:   crm.getMaxSectorConcentration(report.SectorConcentrations.SectorExposures),
			ThresholdValue: 0.25,
			Status:         crm.getKRIStatus(crm.getMaxSectorConcentration(report.SectorConcentrations.SectorExposures), 0.25, true),
		},
		{
			IndicatorName:  "Large Exposure Count",
			CurrentValue:   float64(report.LargeExposures.NumberOfLargeExposures),
			ThresholdValue: 20,
			Status:         crm.getKRIStatus(float64(report.LargeExposures.NumberOfLargeExposures), 20, true),
		},
	}

	// Diversification metrics
	metrics.DiversificationIndex = crm.calculateDiversificationIndex(report)
	metrics.EffectiveNumberOfExposures = crm.calculateEffectiveNumberOfExposures(report.CreditConcentrations.HerfindahlIndex)

	// Correlation analysis
	metrics.PortfolioCorrelation = crm.estimatePortfolioCorrelation(report)
	metrics.SystemicRiskContribution = crm.calculateSystemicRiskContribution(report)

	return metrics
}

// Stress testing

func (crm *ConcentrationRiskManager) runConcentrationStressTest(ctx context.Context, data InstitutionData, report *ConcentrationRiskReport) *ConcentrationStressTest {
	stressTest := &ConcentrationStressTest{
		TestDate:  time.Now(),
		Scenarios: []ConcentrationStressScenario{},
	}

	// Scenario 1: Single large borrower default
	scenario1 := crm.stressTestLargeBorrowerDefault(data, report)
	stressTest.Scenarios = append(stressTest.Scenarios, scenario1)

	// Scenario 2: Sector downturn
	scenario2 := crm.stressTestSectorDownturn(data, report)
	stressTest.Scenarios = append(stressTest.Scenarios, scenario2)

	// Scenario 3: Geographic crisis
	scenario3 := crm.stressTestGeographicCrisis(data, report)
	stressTest.Scenarios = append(stressTest.Scenarios, scenario3)

	// Scenario 4: Correlated defaults
	scenario4 := crm.stressTestCorrelatedDefaults(data, report)
	stressTest.Scenarios = append(stressTest.Scenarios, scenario4)

	// Calculate aggregate impact
	stressTest.AggregateImpact = crm.calculateAggregateStressImpact(stressTest.Scenarios)

	return stressTest
}

func (crm *ConcentrationRiskManager) stressTestLargeBorrowerDefault(data InstitutionData, report *ConcentrationRiskReport) ConcentrationStressScenario {
	scenario := ConcentrationStressScenario{
		ScenarioID:   "LARGE_BORROWER_DEFAULT",
		ScenarioName: "Top 5 Borrower Default",
		Description:  "Simultaneous default of top 5 borrowers",
		Parameters: map[string]interface{}{
			"default_lgd": 0.45, // 45% loss given default
			"provision_coverage": 0.70,
		},
	}

	// Calculate impact
	totalLoss := sdk.NewCoin(data.CapitalBase.Denom, sdk.ZeroInt())
	for i := 0; i < min(5, len(report.CreditConcentrations.Top20Exposures)); i++ {
		exposure := report.CreditConcentrations.Top20Exposures[i]
		lgd := sdk.NewDecWithPrec(45, 2) // 45%
		loss := sdk.NewCoin(exposure.ExposureAmount.Denom, 
			exposure.ExposureAmount.Amount.Mul(sdk.NewInt(int64(lgd.MustFloat64()*100))).Quo(sdk.NewInt(100)))
		totalLoss = totalLoss.Add(loss)
	}

	scenario.ExpectedLoss = totalLoss
	scenario.CapitalImpact = sdk.NewDecFromInt(totalLoss.Amount).Quo(sdk.NewDecFromInt(data.CapitalBase.Amount))
	scenario.RWAImpact = sdk.NewDecWithPrec(15, 2) // 15% increase in RWA
	
	// Post-stress metrics
	scenario.PostStressCapitalRatio = crm.calculatePostStressCapitalRatio(data, scenario)
	scenario.PostStressConcentration = crm.calculatePostStressConcentration(data, scenario)

	return scenario
}

// Mitigation strategies

func (crm *ConcentrationRiskManager) generateMitigationStrategies(report *ConcentrationRiskReport) []MitigationStrategy {
	strategies := []MitigationStrategy{}

	// Credit concentration mitigation
	if report.CreditConcentrations.ConcentrationRatio.GT(sdk.NewDecWithPrec(50, 2)) {
		strategies = append(strategies, MitigationStrategy{
			StrategyID:   "CREDIT_DIVERSIFICATION",
			StrategyType: "DIVERSIFICATION",
			Description:  "Increase portfolio diversification by targeting smaller exposures",
			Actions: []MitigationAction{
				{
					ActionType:   "EXPOSURE_LIMIT",
					Description:  "Reduce single borrower limit to 20% of capital",
					Timeline:     "6 months",
					ExpectedImpact: "Reduce concentration ratio by 10%",
				},
				{
					ActionType:   "NEW_SECTORS",
					Description:  "Target 3 new sectors for lending",
					Timeline:     "12 months",
					ExpectedImpact: "Improve sector diversification",
				},
			},
			Priority:     "HIGH",
			CostEstimate: sdk.NewCoin("usd", sdk.NewInt(500000)),
		})
	}

	// Sector concentration mitigation
	for _, sector := range report.SectorConcentrations.SectorExposures {
		if sector.PercentOfTotal.GT(sdk.NewDecWithPrec(25, 2)) {
			strategies = append(strategies, MitigationStrategy{
				StrategyID:   fmt.Sprintf("SECTOR_LIMIT_%s", sector.SectorCode),
				StrategyType: "LIMIT_MANAGEMENT",
				Description:  fmt.Sprintf("Reduce %s sector exposure below 25%%", sector.SectorName),
				Actions: []MitigationAction{
					{
						ActionType:   "EXPOSURE_REDUCTION",
						Description:  "No new lending to sector until below limit",
						Timeline:     "Immediate",
						ExpectedImpact: fmt.Sprintf("Reduce %s exposure by %.1f%%", sector.SectorName, sector.PercentOfTotal.Sub(sdk.NewDecWithPrec(25, 2)).MustFloat64()*100),
					},
				},
				Priority: "HIGH",
			})
		}
	}

	// Large exposure mitigation
	if len(report.LargeExposures.BreachAnalysis) > 0 {
		strategies = append(strategies, MitigationStrategy{
			StrategyID:   "LARGE_EXPOSURE_REDUCTION",
			StrategyType: "EXPOSURE_REDUCTION",
			Description:  "Reduce breaching large exposures to regulatory limits",
			Actions: []MitigationAction{
				{
					ActionType:   "SYNDICATION",
					Description:  "Syndicate excess exposures to other banks",
					Timeline:     "3 months",
					ExpectedImpact: "Achieve regulatory compliance",
				},
				{
					ActionType:   "COLLATERAL_ENHANCEMENT",
					Description:  "Obtain additional collateral to reduce net exposure",
					Timeline:     "2 months",
					ExpectedImpact: "Reduce net exposure by 30%",
				},
			},
			Priority:     "CRITICAL",
			RegulatoryRequirement: true,
		})
	}

	// Risk transfer strategies
	if report.RiskMetrics.OverallConcentrationScore > 60 {
		strategies = append(strategies, MitigationStrategy{
			StrategyID:   "RISK_TRANSFER",
			StrategyType: "RISK_TRANSFER",
			Description:  "Transfer concentration risk through derivatives and insurance",
			Actions: []MitigationAction{
				{
					ActionType:   "CREDIT_DERIVATIVES",
					Description:  "Purchase credit default swaps on top 10 exposures",
					Timeline:     "1 month",
					ExpectedImpact: "Transfer 40% of large exposure risk",
				},
				{
					ActionType:   "PORTFOLIO_INSURANCE",
					Description:  "Obtain portfolio credit insurance",
					Timeline:     "3 months",
					ExpectedImpact: "Cap maximum portfolio loss at 5% of capital",
				},
			},
			Priority:     "MEDIUM",
			CostEstimate: sdk.NewCoin("usd", sdk.NewInt(2000000)),
		})
	}

	return strategies
}

// Helper structures and methods

type InstitutionData struct {
	InstitutionID    string
	CapitalBase      sdk.Coin
	TotalExposure    sdk.Coin
	CreditExposures  []CreditExposure
	TotalAssets      sdk.Coin
	HomeCountry      string
	HomeCurrency     string
}

type CreditExposure struct {
	ID             string
	BorrowerName   string
	BorrowerGroup  string
	Amount         sdk.Coin
	Currency       string
	RiskGrade      string
	Sector         string
	Country        string
	MaturityDate   time.Time
	CollateralValue sdk.Coin
	NPAStatus      bool
}

type ExposureDetail struct {
	Rank             int       `json:"rank"`
	ExposureID       string    `json:"exposure_id"`
	BorrowerName     string    `json:"borrower_name"`
	ExposureAmount   sdk.Coin  `json:"exposure_amount"`
	PercentOfCapital sdk.Dec   `json:"percent_of_capital"`
	RiskGrade        string    `json:"risk_grade"`
	Sector           string    `json:"sector"`
	Geography        string    `json:"geography"`
	MaturityDate     time.Time `json:"maturity_date"`
}

type LimitAnalysis struct {
	LimitType        string              `json:"limit_type"`
	LimitPercentage  sdk.Dec             `json:"limit_percentage"`
	LimitAmount      sdk.Coin            `json:"limit_amount"`
	CurrentExposure  sdk.Coin            `json:"current_exposure"`
	UtilizationRatio sdk.Dec             `json:"utilization_ratio"`
	Breaches         []LimitBreachDetail `json:"breaches"`
	IsCompliant      bool                `json:"is_compliant"`
}

type LimitBreachDetail struct {
	ExposureID       string   `json:"exposure_id"`
	BorrowerName     string   `json:"borrower_name"`
	ExposureAmount   sdk.Coin `json:"exposure_amount"`
	LimitAmount      sdk.Coin `json:"limit_amount"`
	ExcessAmount     sdk.Coin `json:"excess_amount"`
	BreachPercentage sdk.Dec  `json:"breach_percentage"`
}

type GroupExposure struct {
	GroupID       string   `json:"group_id"`
	GroupName     string   `json:"group_name"`
	TotalExposure sdk.Coin `json:"total_exposure"`
	Members       []string `json:"members"`
}

type LorenzPoint struct {
	CumulativeCount    sdk.Dec `json:"cumulative_count"`
	CumulativeExposure sdk.Dec `json:"cumulative_exposure"`
}

type RiskGradeDistribution struct {
	AAA     sdk.Dec `json:"aaa"`
	AA      sdk.Dec `json:"aa"`
	A       sdk.Dec `json:"a"`
	BBB     sdk.Dec `json:"bbb"`
	BB      sdk.Dec `json:"bb"`
	B       sdk.Dec `json:"b"`
	CCC     sdk.Dec `json:"ccc"`
	Default sdk.Dec `json:"default"`
}

type NPAConcentrationAnalysis struct {
	TotalNPAAmount      sdk.Coin           `json:"total_npa_amount"`
	NPAPercentage       sdk.Dec            `json:"npa_percentage"`
	Top10NPAExposures   []NPAExposure      `json:"top_10_npa_exposures"`
	SectorNPADistribution []SectorNPA       `json:"sector_npa_distribution"`
	AgeingAnalysis      NPAAgeingAnalysis  `json:"ageing_analysis"`
}

type NPAExposure struct {
	ExposureID     string    `json:"exposure_id"`
	BorrowerName   string    `json:"borrower_name"`
	NPAAmount      sdk.Coin  `json:"npa_amount"`
	DaysPastDue    int       `json:"days_past_due"`
	ProvisionHeld  sdk.Coin  `json:"provision_held"`
	ExpectedLoss   sdk.Coin  `json:"expected_loss"`
}

type SectorNPA struct {
	SectorCode    string   `json:"sector_code"`
	NPAAmount     sdk.Coin `json:"npa_amount"`
	NPAPercentage sdk.Dec  `json:"npa_percentage"`
}

type NPAAgeingAnalysis struct {
	Below90Days    sdk.Coin `json:"below_90_days"`
	Days90To180    sdk.Coin `json:"days_90_to_180"`
	Days180To365   sdk.Coin `json:"days_180_to_365"`
	Above365Days   sdk.Coin `json:"above_365_days"`
}

type SectorExposure struct {
	SectorCode        string           `json:"sector_code"`
	SectorName        string           `json:"sector_name"`
	TotalExposure     sdk.Coin         `json:"total_exposure"`
	PercentOfTotal    sdk.Dec          `json:"percent_of_total"`
	NumberOfBorrowers int              `json:"number_of_borrowers"`
	AverageRiskGrade  string           `json:"average_risk_grade"`
	NPAAmount         sdk.Coin         `json:"npa_amount"`
	NPAPercentage     sdk.Dec          `json:"npa_percentage"`
	LimitUtilization  sdk.Dec          `json:"limit_utilization"`
	IsHighRisk        bool             `json:"is_high_risk"`
	Exposures         []CreditExposure `json:"-"`
}

type SectorLimit struct {
	SectorCode      string  `json:"sector_code"`
	MaxExposure     sdk.Dec `json:"max_exposure"`
	CurrentExposure sdk.Dec `json:"current_exposure"`
	AvailableLimit  sdk.Dec `json:"available_limit"`
	IsBreached      bool    `json:"is_breached"`
}

type SectorCorrelation struct {
	Sector1     string  `json:"sector1"`
	Sector2     string  `json:"sector2"`
	Correlation sdk.Dec `json:"correlation"`
	Significance string `json:"significance"`
}

type CyclicalityMetrics struct {
	CyclicalExposure    sdk.Coin `json:"cyclical_exposure"`
	NonCyclicalExposure sdk.Coin `json:"non_cyclical_exposure"`
	CyclicalPercentage  sdk.Dec  `json:"cyclical_percentage"`
	CyclicalSectors     []string `json:"cyclical_sectors"`
}

type RiskScore struct {
	Score       int    `json:"score"`
	Rating      string `json:"rating"`
	Factors     []string `json:"factors"`
}

type CountryExposure struct {
	CountryCode        string   `json:"country_code"`
	CountryName        string   `json:"country_name"`
	TotalExposure      sdk.Coin `json:"total_exposure"`
	PercentOfTotal     sdk.Dec  `json:"percent_of_total"`
	NumberOfBorrowers  int      `json:"number_of_borrowers"`
	CountryRating      string   `json:"country_rating"`
	PoliticalRiskScore int      `json:"political_risk_score"`
	ConcentrationRisk  string   `json:"concentration_risk"`
}

type RegionalExposure struct {
	Region         string   `json:"region"`
	TotalExposure  sdk.Coin `json:"total_exposure"`
	PercentOfTotal sdk.Dec  `json:"percent_of_total"`
	Countries      []string `json:"countries"`
}

type PoliticalRiskExposure struct {
	CountryCode    string   `json:"country_code"`
	ExposureAmount sdk.Coin `json:"exposure_amount"`
	RiskScore      int      `json:"risk_score"`
	RiskFactors    []string `json:"risk_factors"`
}

type CurrencyExposure struct {
	Currency       string   `json:"currency"`
	ExposureAmount sdk.Coin `json:"exposure_amount"`
	PercentOfTotal sdk.Dec  `json:"percent_of_total"`
	ExchangeRisk   string   `json:"exchange_risk"`
}

type CounterpartyExposure struct {
	CounterpartyID   string   `json:"counterparty_id"`
	CounterpartyName string   `json:"counterparty_name"`
	TotalExposure    sdk.Coin `json:"total_exposure"`
	ExposureType     string   `json:"exposure_type"`
	CreditRating     string   `json:"credit_rating"`
	SystemicImportance bool   `json:"systemic_importance"`
}

type CorporateGroupExposure struct {
	GroupID          string   `json:"group_id"`
	GroupName        string   `json:"group_name"`
	TotalExposure    sdk.Coin `json:"total_exposure"`
	NumberOfEntities int      `json:"number_of_entities"`
	GroupRating      string   `json:"group_rating"`
}

type SystemicCounterparty struct {
	CounterpartyID   string  `json:"counterparty_id"`
	CounterpartyName string  `json:"counterparty_name"`
	SystemicScore    sdk.Dec `json:"systemic_score"`
	Interconnectedness string `json:"interconnectedness"`
}

type ProductConcentrationAnalysis struct {
	ProductExposures     []ProductExposure  `json:"product_exposures"`
	ProductRiskProfile   ProductRiskProfile `json:"product_risk_profile"`
	ProductCorrelations  []ProductCorrelation `json:"product_correlations"`
}

type ProductExposure struct {
	ProductType     string   `json:"product_type"`
	TotalExposure   sdk.Coin `json:"total_exposure"`
	PercentOfTotal  sdk.Dec  `json:"percent_of_total"`
	AverageMaturity string   `json:"average_maturity"`
	RiskWeight      sdk.Dec  `json:"risk_weight"`
}

type ProductRiskProfile struct {
	HighRiskProducts   []string `json:"high_risk_products"`
	LowRiskProducts    []string `json:"low_risk_products"`
	ComplexProducts    []string `json:"complex_products"`
	RiskConcentration  sdk.Dec  `json:"risk_concentration"`
}

type ProductCorrelation struct {
	Product1    string  `json:"product1"`
	Product2    string  `json:"product2"`
	Correlation sdk.Dec `json:"correlation"`
}

type CollateralConcentrationAnalysis struct {
	CollateralTypes      []CollateralType     `json:"collateral_types"`
	CollateralQuality    CollateralQuality    `json:"collateral_quality"`
	GeographicDistribution []CollateralLocation `json:"geographic_distribution"`
	ValuationRisk        ValuationRisk        `json:"valuation_risk"`
}

type CollateralType struct {
	Type           string   `json:"type"`
	TotalValue     sdk.Coin `json:"total_value"`
	PercentOfTotal sdk.Dec  `json:"percent_of_total"`
	AverageHaircut sdk.Dec  `json:"average_haircut"`
	Liquidity      string   `json:"liquidity"`
}

type CollateralQuality struct {
	Grade1Collateral sdk.Coin `json:"grade_1_collateral"`
	Grade2Collateral sdk.Coin `json:"grade_2_collateral"`
	Grade3Collateral sdk.Coin `json:"grade_3_collateral"`
	QualityScore     sdk.Dec  `json:"quality_score"`
}

type CollateralLocation struct {
	Location       string   `json:"location"`
	CollateralValue sdk.Coin `json:"collateral_value"`
	ConcentrationRisk string `json:"concentration_risk"`
}

type ValuationRisk struct {
	LastValuationDate   time.Time `json:"last_valuation_date"`
	StaleValuations     int       `json:"stale_valuations"`
	ValuationVolatility sdk.Dec   `json:"valuation_volatility"`
	StressedValue       sdk.Coin  `json:"stressed_value"`
}

type MaturityConcentrationAnalysis struct {
	MaturityBuckets     []MaturityBucket    `json:"maturity_buckets"`
	AverageMaturity     string              `json:"average_maturity"`
	MaturityMismatch    MaturityMismatch    `json:"maturity_mismatch"`
	RefinancingRisk     RefinancingRisk     `json:"refinancing_risk"`
}

type MaturityBucket struct {
	Bucket         string   `json:"bucket"`
	ExposureAmount sdk.Coin `json:"exposure_amount"`
	PercentOfTotal sdk.Dec  `json:"percent_of_total"`
	NumberOfLoans  int      `json:"number_of_loans"`
}

type MaturityMismatch struct {
	ShortTermAssets      sdk.Coin `json:"short_term_assets"`
	ShortTermLiabilities sdk.Coin `json:"short_term_liabilities"`
	MismatchRatio        sdk.Dec  `json:"mismatch_ratio"`
	LiquidityRisk        string   `json:"liquidity_risk"`
}

type RefinancingRisk struct {
	Next12Months    sdk.Coin `json:"next_12_months"`
	Next24Months    sdk.Coin `json:"next_24_months"`
	ConcentrationDates []time.Time `json:"concentration_dates"`
	RiskLevel       string   `json:"risk_level"`
}

type LargeExposureDetail struct {
	ExposureID         string                `json:"exposure_id"`
	CounterpartyName   string                `json:"counterparty_name"`
	ExposureAmount     sdk.Coin              `json:"exposure_amount"`
	PercentOfCapital   sdk.Dec               `json:"percent_of_capital"`
	ExposureType       string                `json:"exposure_type"`
	CreditRating       string                `json:"credit_rating"`
	CollateralCoverage sdk.Dec               `json:"collateral_coverage"`
	NetExposure        sdk.Coin              `json:"net_exposure"`
	ExposureComponents []ExposureComponent   `json:"exposure_components"`
	LastReviewDate     time.Time             `json:"last_review_date"`
	NextReviewDate     time.Time             `json:"next_review_date"`
	IsBreaching        bool                  `json:"is_breaching"`
}

type ExposureComponent struct {
	ComponentType   string   `json:"component_type"`
	Amount          sdk.Coin `json:"amount"`
	RiskWeight      sdk.Dec  `json:"risk_weight"`
}

type LimitBreach struct {
	BreachID         string    `json:"breach_id"`
	LimitType        string    `json:"limit_type"`
	CounterpartyID   string    `json:"counterparty_id"`
	CounterpartyName string    `json:"counterparty_name"`
	LimitAmount      sdk.Dec   `json:"limit_amount"`
	ActualAmount     sdk.Dec   `json:"actual_amount"`
	ExcessAmount     sdk.Dec   `json:"excess_amount"`
	BreachDate       time.Time `json:"breach_date"`
	Status           string    `json:"status"`
	RemediationPlan  string    `json:"remediation_plan"`
}

type LargeExposureTrend struct {
	PeriodStart            time.Time `json:"period_start"`
	PeriodEnd              time.Time `json:"period_end"`
	AverageLargeExposures  int       `json:"average_large_exposures"`
	MaxLargeExposures      int       `json:"max_large_exposures"`
	TrendDirection         string    `json:"trend_direction"`
	GrowthRate             sdk.Dec   `json:"growth_rate"`
}

type ConnectedLendingAnalysis struct {
	ConnectedExposures    []ConnectedExposure  `json:"connected_exposures"`
	TotalConnectedLending sdk.Coin             `json:"total_connected_lending"`
	PercentOfCapital      sdk.Dec              `json:"percent_of_capital"`
	ComplianceStatus      string               `json:"compliance_status"`
	RelatedPartyGroups    []RelatedPartyGroup  `json:"related_party_groups"`
}

type ConnectedExposure struct {
	ExposureID       string   `json:"exposure_id"`
	ConnectedEntity  string   `json:"connected_entity"`
	ConnectionType   string   `json:"connection_type"`
	ExposureAmount   sdk.Coin `json:"exposure_amount"`
	ApprovalLevel    string   `json:"approval_level"`
	LastReviewDate   time.Time `json:"last_review_date"`
}

type RelatedPartyGroup struct {
	GroupID         string   `json:"group_id"`
	GroupName       string   `json:"group_name"`
	TotalExposure   sdk.Coin `json:"total_exposure"`
	NumberOfEntities int     `json:"number_of_entities"`
	RelationshipType string  `json:"relationship_type"`
}

type ConcentrationRiskMetrics struct {
	OverallConcentrationScore   float64             `json:"overall_concentration_score"`
	RiskLevel                   string              `json:"risk_level"`
	KeyRiskIndicators           []KeyRiskIndicator  `json:"key_risk_indicators"`
	DiversificationIndex        sdk.Dec             `json:"diversification_index"`
	EffectiveNumberOfExposures  float64             `json:"effective_number_of_exposures"`
	PortfolioCorrelation        sdk.Dec             `json:"portfolio_correlation"`
	SystemicRiskContribution    sdk.Dec             `json:"systemic_risk_contribution"`
	CalculationDate             time.Time           `json:"calculation_date"`
}

type KeyRiskIndicator struct {
	IndicatorName  string  `json:"indicator_name"`
	CurrentValue   float64 `json:"current_value"`
	ThresholdValue float64 `json:"threshold_value"`
	Status         string  `json:"status"`
	TrendDirection string  `json:"trend_direction"`
}

type ConcentrationCompliance struct {
	IsCompliant            bool                    `json:"is_compliant"`
	ComplianceDate         time.Time               `json:"compliance_date"`
	ViolationCount         int                     `json:"violation_count"`
	ComplianceViolations   []ComplianceViolation   `json:"compliance_violations"`
	RegulatoryRequirements []RegulatoryRequirement `json:"regulatory_requirements"`
	NextReviewDate         time.Time               `json:"next_review_date"`
}

type ComplianceViolation struct {
	ViolationID     string    `json:"violation_id"`
	RequirementID   string    `json:"requirement_id"`
	ViolationType   string    `json:"violation_type"`
	Description     string    `json:"description"`
	Severity        string    `json:"severity"`
	DateIdentified  time.Time `json:"date_identified"`
	RemediationPlan string    `json:"remediation_plan"`
	Status          string    `json:"status"`
}

type RegulatoryRequirement struct {
	RequirementID   string  `json:"requirement_id"`
	Description     string  `json:"description"`
	LimitType       string  `json:"limit_type"`
	LimitValue      sdk.Dec `json:"limit_value"`
	CurrentValue    sdk.Dec `json:"current_value"`
	ComplianceStatus string `json:"compliance_status"`
}

type ConcentrationStressTest struct {
	TestDate        time.Time                      `json:"test_date"`
	Scenarios       []ConcentrationStressScenario  `json:"scenarios"`
	AggregateImpact AggregateStressImpact          `json:"aggregate_impact"`
}

type ConcentrationStressScenario struct {
	ScenarioID              string                 `json:"scenario_id"`
	ScenarioName            string                 `json:"scenario_name"`
	Description             string                 `json:"description"`
	Parameters              map[string]interface{} `json:"parameters"`
	ExpectedLoss            sdk.Coin               `json:"expected_loss"`
	CapitalImpact           sdk.Dec                `json:"capital_impact"`
	RWAImpact               sdk.Dec                `json:"rwa_impact"`
	PostStressCapitalRatio  sdk.Dec                `json:"post_stress_capital_ratio"`
	PostStressConcentration sdk.Dec                `json:"post_stress_concentration"`
}

type AggregateStressImpact struct {
	TotalExpectedLoss   sdk.Coin `json:"total_expected_loss"`
	MaxCapitalImpact    sdk.Dec  `json:"max_capital_impact"`
	WorstCaseScenario   string   `json:"worst_case_scenario"`
	RequiredCapitalBuffer sdk.Coin `json:"required_capital_buffer"`
}

type MitigationStrategy struct {
	StrategyID            string              `json:"strategy_id"`
	StrategyType          string              `json:"strategy_type"`
	Description           string              `json:"description"`
	Actions               []MitigationAction  `json:"actions"`
	Priority              string              `json:"priority"`
	Timeline              string              `json:"timeline"`
	CostEstimate          sdk.Coin            `json:"cost_estimate,omitempty"`
	ExpectedBenefit       string              `json:"expected_benefit"`
	ResponsibleParty      string              `json:"responsible_party"`
	Status                string              `json:"status"`
	RegulatoryRequirement bool                `json:"regulatory_requirement"`
}

type MitigationAction struct {
	ActionType     string `json:"action_type"`
	Description    string `json:"description"`
	Timeline       string `json:"timeline"`
	ExpectedImpact string `json:"expected_impact"`
	Status         string `json:"status"`
}

type RiskRecommendation struct {
	RecommendationID string `json:"recommendation_id"`
	Category         string `json:"category"`
	Priority         string `json:"priority"`
	Description      string `json:"description"`
	Rationale        string `json:"rationale"`
	Implementation   string `json:"implementation"`
	ExpectedOutcome  string `json:"expected_outcome"`
}

// Limits and configuration

type ConcentrationLimits struct {
	SingleBorrowerLimit    sdk.Dec
	GroupBorrowerLimit     sdk.Dec
	LargeExposureThreshold sdk.Dec
	SectorLimits           map[string]sdk.Dec
	CountryLimits          map[string]sdk.Dec
}

type RiskMetrics struct {
	ConcentrationThresholds map[string]float64
	WarningLevels          map[string]float64
	CriticalLevels         map[string]float64
}

type EarlyWarningSystem struct {
	Triggers      []EarlyWarningTrigger
	AlertLevels   []AlertLevel
	EscalationPath []EscalationStep
}

type EarlyWarningTrigger struct {
	TriggerID    string  `json:"trigger_id"`
	MetricName   string  `json:"metric_name"`
	ThresholdValue float64 `json:"threshold_value"`
	CurrentValue float64 `json:"current_value"`
	IsTriggered  bool    `json:"is_triggered"`
}

type AlertLevel struct {
	Level       string `json:"level"`
	Description string `json:"description"`
	Actions     []string `json:"actions"`
}

type EscalationStep struct {
	Level        string `json:"level"`
	Responsible  string `json:"responsible"`
	TimeFrame    string `json:"time_frame"`
}

// Initialize functions

func initializeConcentrationLimits() ConcentrationLimits {
	return ConcentrationLimits{
		SingleBorrowerLimit:    sdk.NewDecWithPrec(25, 2), // 25% of capital
		GroupBorrowerLimit:     sdk.NewDecWithPrec(40, 2), // 40% of capital
		LargeExposureThreshold: sdk.NewDecWithPrec(10, 2), // 10% of capital
		SectorLimits: map[string]sdk.Dec{
			"REAL_ESTATE":     sdk.NewDecWithPrec(20, 2),
			"INFRASTRUCTURE":  sdk.NewDecWithPrec(30, 2),
			"MANUFACTURING":   sdk.NewDecWithPrec(25, 2),
			"SERVICES":        sdk.NewDecWithPrec(30, 2),
			"AGRICULTURE":     sdk.NewDecWithPrec(15, 2),
			"RETAIL":          sdk.NewDecWithPrec(20, 2),
		},
		CountryLimits: map[string]sdk.Dec{
			"DOMESTIC":        sdk.NewDecWithPrec(80, 2),
			"CROSS_BORDER":    sdk.NewDecWithPrec(20, 2),
		},
	}
}

func initializeRiskMetrics() RiskMetrics {
	return RiskMetrics{
		ConcentrationThresholds: map[string]float64{
			"HERFINDAHL_INDEX":     0.15,
			"TOP_20_CONCENTRATION": 0.40,
			"GINI_COEFFICIENT":     0.60,
		},
		WarningLevels: map[string]float64{
			"SINGLE_BORROWER":  0.20,
			"SECTOR_EXPOSURE":  0.20,
			"COUNTRY_EXPOSURE": 0.15,
		},
		CriticalLevels: map[string]float64{
			"SINGLE_BORROWER":  0.25,
			"SECTOR_EXPOSURE":  0.30,
			"COUNTRY_EXPOSURE": 0.25,
		},
	}
}

func initializeEarlyWarningSystem() EarlyWarningSystem {
	return EarlyWarningSystem{
		Triggers: []EarlyWarningTrigger{
			{
				TriggerID:      "CONCENTRATION_INCREASE",
				MetricName:     "Herfindahl Index Change",
				ThresholdValue: 0.02, // 2% increase
			},
			{
				TriggerID:      "LARGE_EXPOSURE_GROWTH",
				MetricName:     "Large Exposure Count",
				ThresholdValue: 5, // More than 5 new large exposures
			},
		},
		AlertLevels: []AlertLevel{
			{
				Level:       "YELLOW",
				Description: "Elevated concentration risk",
				Actions:     []string{"Monitor closely", "Review limits"},
			},
			{
				Level:       "ORANGE",
				Description: "High concentration risk",
				Actions:     []string{"Implement controls", "Report to board"},
			},
			{
				Level:       "RED",
				Description: "Critical concentration risk",
				Actions:     []string{"Immediate action", "Regulatory notification"},
			},
		},
	}
}

// Helper method implementations (stubs for compilation)

func (crm *ConcentrationRiskManager) getInstitutionData(ctx context.Context, institutionID string) (InstitutionData, error) {
	// Stub implementation
	return InstitutionData{
		InstitutionID: institutionID,
		CapitalBase:   sdk.NewCoin("usd", sdk.NewInt(1000000000)),
		TotalExposure: sdk.NewCoin("usd", sdk.NewInt(8000000000)),
		TotalAssets:   sdk.NewCoin("usd", sdk.NewInt(10000000000)),
		HomeCountry:   "IN",
		HomeCurrency:  "INR",
	}, nil
}

func (crm *ConcentrationRiskManager) getCreditExposures(data InstitutionData) []CreditExposure {
	// Stub - would fetch from database
	return data.CreditExposures
}

func (crm *ConcentrationRiskManager) getMaxSingleExposure(exposures []CreditExposure) sdk.Coin {
	if len(exposures) == 0 {
		return sdk.NewCoin("usd", sdk.ZeroInt())
	}
	return exposures[0].Amount
}

func (crm *ConcentrationRiskManager) calculateUtilizationRatio(current, limit sdk.Coin) sdk.Dec {
	if limit.IsZero() {
		return sdk.ZeroDec()
	}
	return sdk.NewDecFromInt(current.Amount).Quo(sdk.NewDecFromInt(limit.Amount))
}

func (crm *ConcentrationRiskManager) aggregateGroupExposures(exposures []CreditExposure) []GroupExposure {
	groupMap := make(map[string]*GroupExposure)
	
	for _, exp := range exposures {
		if exp.BorrowerGroup != "" {
			if group, exists := groupMap[exp.BorrowerGroup]; exists {
				group.TotalExposure = group.TotalExposure.Add(exp.Amount)
				group.Members = append(group.Members, exp.BorrowerName)
			} else {
				groupMap[exp.BorrowerGroup] = &GroupExposure{
					GroupID:       exp.BorrowerGroup,
					GroupName:     exp.BorrowerGroup,
					TotalExposure: exp.Amount,
					Members:       []string{exp.BorrowerName},
				}
			}
		}
	}
	
	groups := []GroupExposure{}
	for _, group := range groupMap {
		groups = append(groups, *group)
	}
	return groups
}

func (crm *ConcentrationRiskManager) getMaxGroupExposure(groups []GroupExposure) sdk.Coin {
	if len(groups) == 0 {
		return sdk.NewCoin("usd", sdk.ZeroInt())
	}
	
	max := groups[0].TotalExposure
	for _, group := range groups {
		if group.TotalExposure.IsGT(max) {
			max = group.TotalExposure
		}
	}
	return max
}

func (crm *ConcentrationRiskManager) calculateConcentrationRatio(exposures []CreditExposure, capitalBase sdk.Coin) sdk.Dec {
	if len(exposures) < 5 {
		return sdk.ZeroDec()
	}
	
	top5Total := sdk.NewCoin(capitalBase.Denom, sdk.ZeroInt())
	for i := 0; i < 5; i++ {
		top5Total = top5Total.Add(exposures[i].Amount)
	}
	
	return sdk.NewDecFromInt(top5Total.Amount).Quo(sdk.NewDecFromInt(capitalBase.Amount))
}

func (crm *ConcentrationRiskManager) calculateHerfindahlIndex(exposures []CreditExposure, totalExposure sdk.Coin) sdk.Dec {
	hhi := sdk.ZeroDec()
	
	for _, exp := range exposures {
		share := sdk.NewDecFromInt(exp.Amount.Amount).Quo(sdk.NewDecFromInt(totalExposure.Amount))
		hhi = hhi.Add(share.Mul(share))
	}
	
	return hhi
}

func (crm *ConcentrationRiskManager) calculateGiniCoefficient(exposures []CreditExposure) sdk.Dec {
	// Simplified Gini coefficient calculation
	n := len(exposures)
	if n == 0 {
		return sdk.ZeroDec()
	}
	
	// Sort exposures
	sort.Slice(exposures, func(i, j int) bool {
		return exposures[i].Amount.IsLT(exposures[j].Amount)
	})
	
	sumOfProducts := sdk.ZeroDec()
	sumOfExposures := sdk.ZeroDec()
	
	for i, exp := range exposures {
		sumOfProducts = sumOfProducts.Add(
			sdk.NewDec(int64(2*i - n + 1)).Mul(sdk.NewDecFromInt(exp.Amount.Amount)),
		)
		sumOfExposures = sumOfExposures.Add(sdk.NewDecFromInt(exp.Amount.Amount))
	}
	
	if sumOfExposures.IsZero() {
		return sdk.ZeroDec()
	}
	
	gini := sumOfProducts.Quo(sdk.NewDec(int64(n)).Mul(sumOfExposures))
	return gini
}

func (crm *ConcentrationRiskManager) calculateLorenzCurve(exposures []CreditExposure) []LorenzPoint {
	// Sort exposures ascending
	sort.Slice(exposures, func(i, j int) bool {
		return exposures[i].Amount.IsLT(exposures[j].Amount)
	})
	
	points := []LorenzPoint{}
	cumulativeExposure := sdk.ZeroDec()
	totalExposure := sdk.ZeroDec()
	
	for _, exp := range exposures {
		totalExposure = totalExposure.Add(sdk.NewDecFromInt(exp.Amount.Amount))
	}
	
	for i, exp := range exposures {
		cumulativeExposure = cumulativeExposure.Add(sdk.NewDecFromInt(exp.Amount.Amount))
		points = append(points, LorenzPoint{
			CumulativeCount:    sdk.NewDec(int64(i + 1)).Quo(sdk.NewDec(int64(len(exposures)))),
			CumulativeExposure: cumulativeExposure.Quo(totalExposure),
		})
	}
	
	return points
}

func (crm *ConcentrationRiskManager) analyzeRiskDistribution(exposures []CreditExposure) RiskGradeDistribution {
	dist := RiskGradeDistribution{}
	total := sdk.ZeroDec()
	
	gradeAmounts := make(map[string]sdk.Dec)
	for _, exp := range exposures {
		amt := sdk.NewDecFromInt(exp.Amount.Amount)
		gradeAmounts[exp.RiskGrade] = gradeAmounts[exp.RiskGrade].Add(amt)
		total = total.Add(amt)
	}
	
	if total.IsPositive() {
		dist.AAA = gradeAmounts["AAA"].Quo(total)
		dist.AA = gradeAmounts["AA"].Quo(total)
		dist.A = gradeAmounts["A"].Quo(total)
		dist.BBB = gradeAmounts["BBB"].Quo(total)
		dist.BB = gradeAmounts["BB"].Quo(total)
		dist.B = gradeAmounts["B"].Quo(total)
		dist.CCC = gradeAmounts["CCC"].Quo(total)
		dist.Default = gradeAmounts["D"].Quo(total)
	}
	
	return dist
}

func (crm *ConcentrationRiskManager) analyzeNPAConcentration(exposures []CreditExposure) NPAConcentrationAnalysis {
	analysis := NPAConcentrationAnalysis{
		TotalNPAAmount:    sdk.NewCoin("usd", sdk.ZeroInt()),
		Top10NPAExposures: []NPAExposure{},
	}
	
	npaExposures := []NPAExposure{}
	totalExposure := sdk.NewCoin("usd", sdk.ZeroInt())
	
	for _, exp := range exposures {
		totalExposure = totalExposure.Add(exp.Amount)
		if exp.NPAStatus {
			analysis.TotalNPAAmount = analysis.TotalNPAAmount.Add(exp.Amount)
			npaExposures = append(npaExposures, NPAExposure{
				ExposureID:   exp.ID,
				BorrowerName: exp.BorrowerName,
				NPAAmount:    exp.Amount,
				DaysPastDue:  90, // Stub
			})
		}
	}
	
	// Sort and take top 10
	sort.Slice(npaExposures, func(i, j int) bool {
		return npaExposures[i].NPAAmount.IsGT(npaExposures[j].NPAAmount)
	})
	
	for i := 0; i < min(10, len(npaExposures)); i++ {
		analysis.Top10NPAExposures = append(analysis.Top10NPAExposures, npaExposures[i])
	}
	
	if totalExposure.IsPositive() {
		analysis.NPAPercentage = sdk.NewDecFromInt(analysis.TotalNPAAmount.Amount).Quo(sdk.NewDecFromInt(totalExposure.Amount))
	}
	
	return analysis
}

func (crm *ConcentrationRiskManager) getSectorName(sectorCode string) string {
	sectorNames := map[string]string{
		"RE": "Real Estate",
		"INFRA": "Infrastructure",
		"MFG": "Manufacturing",
		"SVC": "Services",
		"AGRI": "Agriculture",
		"RTL": "Retail",
	}
	
	if name, exists := sectorNames[sectorCode]; exists {
		return name
	}
	return sectorCode
}

func (crm *ConcentrationRiskManager) isHighRiskSector(sectorCode string) bool {
	highRiskSectors := []string{"RE", "INFRA", "AVIATION", "HOSPITALITY"}
	for _, sector := range highRiskSectors {
		if sector == sectorCode {
			return true
		}
	}
	return false
}

func (crm *ConcentrationRiskManager) calculateSectorNPA(exposures []CreditExposure) sdk.Coin {
	npa := sdk.NewCoin("usd", sdk.ZeroInt())
	for _, exp := range exposures {
		if exp.NPAStatus {
			npa = npa.Add(exp.Amount)
		}
	}
	return npa
}

func (crm *ConcentrationRiskManager) calculateAverageRiskGrade(exposures []CreditExposure) string {
	if len(exposures) == 0 {
		return "NA"
	}
	
	gradeValues := map[string]int{
		"AAA": 1, "AA": 2, "A": 3, "BBB": 4, "BB": 5, "B": 6, "CCC": 7, "D": 8,
	}
	
	total := 0
	for _, exp := range exposures {
		if val, exists := gradeValues[exp.RiskGrade]; exists {
			total += val
		}
	}
	
	avg := total / len(exposures)
	
	for grade, val := range gradeValues {
		if val == avg {
			return grade
		}
	}
	
	return "BBB" // Default
}

func (crm *ConcentrationRiskManager) getSectorLimit(sectorCode string) SectorLimit {
	if limit, exists := crm.riskLimits.SectorLimits[sectorCode]; exists {
		return SectorLimit{
			SectorCode:  sectorCode,
			MaxExposure: limit,
		}
	}
	
	// Default limit
	return SectorLimit{
		SectorCode:  sectorCode,
		MaxExposure: sdk.NewDecWithPrec(25, 2), // 25%
	}
}

func (crm *ConcentrationRiskManager) calculateIndustryRiskScore(sectorExp *SectorExposure) RiskScore {
	score := 50 // Base score
	
	// Adjust based on NPA
	if sectorExp.NPAPercentage.GT(sdk.NewDecWithPrec(5, 2)) {
		score += 20
	}
	
	// Adjust based on concentration
	if sectorExp.PercentOfTotal.GT(sdk.NewDecWithPrec(20, 2)) {
		score += 15
	}
	
	// High risk sector
	if sectorExp.IsHighRisk {
		score += 15
	}
	
	rating := "LOW"
	if score >= 80 {
		rating = "HIGH"
	} else if score >= 60 {
		rating = "MEDIUM"
	}
	
	return RiskScore{
		Score:  score,
		Rating: rating,
	}
}

func (crm *ConcentrationRiskManager) calculateSectorCorrelations(sectors []SectorExposure) []SectorCorrelation {
	// Simplified correlation calculation
	correlations := []SectorCorrelation{
		{
			Sector1:      "RE",
			Sector2:      "INFRA",
			Correlation:  sdk.NewDecWithPrec(70, 2), // 0.70
			Significance: "HIGH",
		},
		{
			Sector1:      "MFG",
			Sector2:      "SVC",
			Correlation:  sdk.NewDecWithPrec(40, 2), // 0.40
			Significance: "MEDIUM",
		},
	}
	return correlations
}

func (crm *ConcentrationRiskManager) analyzeCyclicality(sectors []SectorExposure) CyclicalityMetrics {
	cyclicalSectors := []string{"RE", "INFRA", "MFG", "AUTO", "AVIATION", "HOSPITALITY"}
	
	metrics := CyclicalityMetrics{
		CyclicalExposure:    sdk.NewCoin("usd", sdk.ZeroInt()),
		NonCyclicalExposure: sdk.NewCoin("usd", sdk.ZeroInt()),
		CyclicalSectors:     []string{},
	}
	
	for _, sector := range sectors {
		isCyclical := false
		for _, cyclical := range cyclicalSectors {
			if sector.SectorCode == cyclical {
				isCyclical = true
				metrics.CyclicalSectors = append(metrics.CyclicalSectors, sector.SectorCode)
				break
			}
		}
		
		if isCyclical {
			metrics.CyclicalExposure = metrics.CyclicalExposure.Add(sector.TotalExposure)
		} else {
			metrics.NonCyclicalExposure = metrics.NonCyclicalExposure.Add(sector.TotalExposure)
		}
	}
	
	total := metrics.CyclicalExposure.Add(metrics.NonCyclicalExposure)
	if total.IsPositive() {
		metrics.CyclicalPercentage = sdk.NewDecFromInt(metrics.CyclicalExposure.Amount).Quo(sdk.NewDecFromInt(total.Amount))
	}
	
	return metrics
}

func (crm *ConcentrationRiskManager) getCountryName(countryCode string) string {
	// Stub implementation
	countryNames := map[string]string{
		"IN": "India",
		"US": "United States",
		"GB": "United Kingdom",
		"CN": "China",
		"JP": "Japan",
	}
	
	if name, exists := countryNames[countryCode]; exists {
		return name
	}
	return countryCode
}

func (crm *ConcentrationRiskManager) getCountryRating(countryCode string) string {
	// Stub implementation
	return "BBB"
}

func (crm *ConcentrationRiskManager) getPoliticalRiskScore(countryCode string) int {
	// Stub implementation - return score 1-10
	riskScores := map[string]int{
		"IN": 3,
		"US": 2,
		"GB": 2,
		"CN": 5,
		"AF": 9,
		"SY": 10,
	}
	
	if score, exists := riskScores[countryCode]; exists {
		return score
	}
	return 5
}

func (crm *ConcentrationRiskManager) isEmergingMarket(countryCode string) bool {
	emergingMarkets := []string{"IN", "CN", "BR", "RU", "ZA", "MX", "ID", "TH", "MY"}
	for _, em := range emergingMarkets {
		if countryCode == em {
			return true
		}
	}
	return false
}

func (crm *ConcentrationRiskManager) assessCountryConcentrationRisk(country *CountryExposure) string {
	if country.PercentOfTotal.GT(sdk.NewDecWithPrec(15, 2)) || country.PoliticalRiskScore > 7 {
		return "HIGH"
	} else if country.PercentOfTotal.GT(sdk.NewDecWithPrec(10, 2)) || country.PoliticalRiskScore > 5 {
		return "MEDIUM"
	}
	return "LOW"
}

func (crm *ConcentrationRiskManager) getPoliticalRiskFactors(countryCode string) []string {
	// Stub implementation
	return []string{"Political instability", "Regulatory changes", "Currency controls"}
}

func (crm *ConcentrationRiskManager) aggregateByRegion(countryMap map[string]*CountryExposure) []RegionalExposure {
	regionMap := map[string][]string{
		"ASIA":   {"IN", "CN", "JP", "KR", "SG"},
		"EUROPE": {"GB", "DE", "FR", "IT", "ES"},
		"AMERICAS": {"US", "CA", "BR", "MX", "AR"},
		"AFRICA": {"ZA", "NG", "EG", "KE", "GH"},
	}
	
	regional := []RegionalExposure{}
	
	for region, countries := range regionMap {
		exposure := sdk.NewCoin("usd", sdk.ZeroInt())
		regionCountries := []string{}
		
		for _, countryCode := range countries {
			if country, exists := countryMap[countryCode]; exists {
				exposure = exposure.Add(country.TotalExposure)
				regionCountries = append(regionCountries, countryCode)
			}
		}
		
		if exposure.IsPositive() {
			regional = append(regional, RegionalExposure{
				Region:        region,
				TotalExposure: exposure,
				Countries:     regionCountries,
			})
		}
	}
	
	return regional
}

func (crm *ConcentrationRiskManager) assessExchangeRisk(currency, homeCurrency string) string {
	if currency == homeCurrency {
		return "NONE"
	}
	
	// Stable currencies
	stableCurrencies := []string{"USD", "EUR", "GBP", "JPY", "CHF"}
	for _, stable := range stableCurrencies {
		if currency == stable {
			return "LOW"
		}
	}
	
	return "MEDIUM"
}

func (crm *ConcentrationRiskManager) analyzeProductConcentrations(ctx context.Context, data InstitutionData) (*ProductConcentrationAnalysis, error) {
	// Stub implementation
	return &ProductConcentrationAnalysis{}, nil
}

func (crm *ConcentrationRiskManager) analyzeCollateralConcentrations(ctx context.Context, data InstitutionData) (*CollateralConcentrationAnalysis, error) {
	// Stub implementation
	return &CollateralConcentrationAnalysis{}, nil
}

func (crm *ConcentrationRiskManager) analyzeMaturityConcentrations(ctx context.Context, data InstitutionData) (*MaturityConcentrationAnalysis, error) {
	// Stub implementation
	return &MaturityConcentrationAnalysis{}, nil
}

func (crm *ConcentrationRiskManager) analyzeCounterpartyConcentrations(ctx context.Context, data InstitutionData) (*CounterpartyConcentrationAnalysis, error) {
	// Stub implementation
	return &CounterpartyConcentrationAnalysis{}, nil
}

func (crm *ConcentrationRiskManager) analyzeConnectedLending(ctx context.Context, data InstitutionData) (*ConnectedLendingAnalysis, error) {
	// Stub implementation
	return &ConnectedLendingAnalysis{
		TotalConnectedLending: sdk.NewCoin("usd", sdk.NewInt(50000000)),
		PercentOfCapital:     sdk.NewDecWithPrec(5, 2),
		ComplianceStatus:     "COMPLIANT",
	}, nil
}

type CounterpartyAggregation struct {
	ID                 string
	Name               string
	Type               string
	TotalExposure      sdk.Coin
	Rating             string
	CollateralCoverage sdk.Dec
	Components         []ExposureComponent
	LastReview         time.Time
	NextReview         time.Time
}

func (crm *ConcentrationRiskManager) aggregateAllExposures(data InstitutionData) []CounterpartyAggregation {
	// Stub implementation
	return []CounterpartyAggregation{}
}

func (crm *ConcentrationRiskManager) calculateNetExposure(exp CounterpartyAggregation) sdk.Coin {
	// Net exposure after collateral
	collateralValue := exp.TotalExposure.Amount.Mul(sdk.NewInt(int64(exp.CollateralCoverage.MustFloat64()*100))).Quo(sdk.NewInt(100))
	netExposure := exp.TotalExposure.Amount.Sub(collateralValue)
	return sdk.NewCoin(exp.TotalExposure.Denom, netExposure)
}

func (crm *ConcentrationRiskManager) analyzeLargeExposureTrend(ctx context.Context, institutionID string) LargeExposureTrend {
	// Stub implementation
	return LargeExposureTrend{
		PeriodStart:           time.Now().AddDate(0, -12, 0),
		PeriodEnd:             time.Now(),
		AverageLargeExposures: 15,
		MaxLargeExposures:     18,
		TrendDirection:        "INCREASING",
		GrowthRate:            sdk.NewDecWithPrec(12, 2), // 12%
	}
}

func (crm *ConcentrationRiskManager) scoreCreditConcentration(analysis *CreditConcentrationAnalysis) float64 {
	score := 0.0
	
	// Herfindahl Index scoring
	hhi := analysis.HerfindahlIndex.MustFloat64()
	if hhi > 0.20 {
		score += 30.0
	} else if hhi > 0.15 {
		score += 20.0
	} else if hhi > 0.10 {
		score += 10.0
	}
	
	// Concentration ratio scoring
	cr := analysis.ConcentrationRatio.MustFloat64()
	if cr > 0.50 {
		score += 25.0
	} else if cr > 0.40 {
		score += 15.0
	} else if cr > 0.30 {
		score += 10.0
	}
	
	// Limit breaches
	if !analysis.SingleBorrowerLimit.IsCompliant {
		score += 20.0
	}
	if !analysis.GroupBorrowerLimit.IsCompliant {
		score += 15.0
	}
	
	// NPA concentration
	if analysis.NPAConcentration.NPAPercentage.GT(sdk.NewDecWithPrec(5, 2)) {
		score += 10.0
	}
	
	return score
}

func (crm *ConcentrationRiskManager) scoreSectorConcentration(analysis *SectorConcentrationAnalysis) float64 {
	score := 0.0
	
	// Maximum sector concentration
	maxConc := crm.getMaxSectorConcentration(analysis.SectorExposures)
	if maxConc > 0.30 {
		score += 30.0
	} else if maxConc > 0.25 {
		score += 20.0
	} else if maxConc > 0.20 {
		score += 10.0
	}
	
	// High risk sector exposure
	highRiskExposure := 0.0
	for _, sector := range analysis.SectorExposures {
		if sector.IsHighRisk {
			highRiskExposure += sector.PercentOfTotal.MustFloat64()
		}
	}
	if highRiskExposure > 0.40 {
		score += 25.0
	} else if highRiskExposure > 0.30 {
		score += 15.0
	}
	
	// Sector limit breaches
	breachCount := 0
	for _, limit := range analysis.SectorLimits {
		if limit.IsBreached {
			breachCount++
		}
	}
	score += float64(breachCount) * 5.0
	
	// Cyclicality
	if analysis.CyclicalityAnalysis.CyclicalPercentage.GT(sdk.NewDecWithPrec(60, 2)) {
		score += 15.0
	}
	
	return score
}

func (crm *ConcentrationRiskManager) scoreGeographicConcentration(analysis *GeographicConcentrationAnalysis) float64 {
	score := 0.0
	
	// Cross-border exposure
	totalExposure := analysis.CrossBorderExposure.Add(analysis.EmergingMarketExposure)
	if len(analysis.CountryExposures) > 0 {
		crossBorderPercent := sdk.NewDecFromInt(analysis.CrossBorderExposure.Amount).Quo(sdk.NewDecFromInt(totalExposure.Amount))
		if crossBorderPercent.GT(sdk.NewDecWithPrec(30, 2)) {
			score += 20.0
		} else if crossBorderPercent.GT(sdk.NewDecWithPrec(20, 2)) {
			score += 10.0
		}
	}
	
	// Emerging market exposure
	if len(analysis.CountryExposures) > 0 {
		emergingPercent := sdk.NewDecFromInt(analysis.EmergingMarketExposure.Amount).Quo(sdk.NewDecFromInt(totalExposure.Amount))
		if emergingPercent.GT(sdk.NewDecWithPrec(20, 2)) {
			score += 15.0
		}
	}
	
	// Political risk exposure
	for _, risk := range analysis.PoliticalRiskExposure {
		if risk.RiskScore >= 8 {
			score += 10.0
		} else if risk.RiskScore >= 6 {
			score += 5.0
		}
	}
	
	// Currency concentration
	nonHomeCurrency := 0
	for _, curr := range analysis.CurrencyConcentration {
		if curr.ExchangeRisk != "NONE" {
			nonHomeCurrency++
		}
	}
	if nonHomeCurrency > 5 {
		score += 15.0
	}
	
	return score
}

func (crm *ConcentrationRiskManager) scoreLargeExposures(analysis *LargeExposureAnalysis) float64 {
	score := 0.0
	
	// Number of large exposures
	if analysis.NumberOfLargeExposures > 20 {
		score += 25.0
	} else if analysis.NumberOfLargeExposures > 15 {
		score += 15.0
	} else if analysis.NumberOfLargeExposures > 10 {
		score += 10.0
	}
	
	// Breaches
	score += float64(len(analysis.BreachAnalysis)) * 10.0
	
	// Trend
	if analysis.TrendAnalysis.TrendDirection == "INCREASING" && analysis.TrendAnalysis.GrowthRate.GT(sdk.NewDecWithPrec(10, 2)) {
		score += 15.0
	}
	
	return score
}

func (crm *ConcentrationRiskManager) getKRIStatus(current, threshold float64, higherIsBad bool) string {
	if higherIsBad {
		if current > threshold*1.2 {
			return "CRITICAL"
		} else if current > threshold {
			return "WARNING"
		}
	} else {
		if current < threshold*0.8 {
			return "CRITICAL"
		} else if current < threshold {
			return "WARNING"
		}
	}
	return "NORMAL"
}

func (crm *ConcentrationRiskManager) calculateTop20Concentration(top20 []ExposureDetail) float64 {
	total := sdk.ZeroDec()
	for _, exp := range top20 {
		total = total.Add(exp.PercentOfCapital)
	}
	return total.MustFloat64()
}

func (crm *ConcentrationRiskManager) getMaxSectorConcentration(sectors []SectorExposure) float64 {
	max := 0.0
	for _, sector := range sectors {
		if sector.PercentOfTotal.MustFloat64() > max {
			max = sector.PercentOfTotal.MustFloat64()
		}
	}
	return max
}

func (crm *ConcentrationRiskManager) calculateDiversificationIndex(report *ConcentrationRiskReport) sdk.Dec {
	// Simpson's diversity index
	return sdk.OneDec().Sub(report.CreditConcentrations.HerfindahlIndex)
}

func (crm *ConcentrationRiskManager) calculateEffectiveNumberOfExposures(hhi sdk.Dec) float64 {
	if hhi.IsZero() {
		return 0
	}
	return 1.0 / hhi.MustFloat64()
}

func (crm *ConcentrationRiskManager) estimatePortfolioCorrelation(report *ConcentrationRiskReport) sdk.Dec {
	// Simplified portfolio correlation estimate
	avgSectorCorr := sdk.NewDecWithPrec(30, 2) // 0.30
	sectorConcentration := crm.getMaxSectorConcentration(report.SectorConcentrations.SectorExposures)
	
	// Higher concentration implies higher effective correlation
	return avgSectorCorr.Mul(sdk.NewDecFromBigInt(sdk.NewInt(int64(sectorConcentration * 100)).BigInt()))
}

func (crm *ConcentrationRiskManager) calculateSystemicRiskContribution(report *ConcentrationRiskReport) sdk.Dec {
	// Simplified systemic risk calculation
	largeExpPercent := float64(report.LargeExposures.NumberOfLargeExposures) / 100.0
	return sdk.NewDecWithPrec(int64(largeExpPercent*100), 2)
}

func (crm *ConcentrationRiskManager) calculatePostStressCapitalRatio(data InstitutionData, scenario ConcentrationStressScenario) sdk.Dec {
	stressedCapital := data.CapitalBase.Amount.Sub(scenario.ExpectedLoss.Amount)
	stressedRWA := data.TotalAssets.Amount.Mul(sdk.NewInt(int64((1 + scenario.RWAImpact.MustFloat64()) * 100))).Quo(sdk.NewInt(100))
	
	if stressedRWA.IsZero() {
		return sdk.ZeroDec()
	}
	
	return sdk.NewDecFromInt(stressedCapital).Quo(sdk.NewDecFromInt(stressedRWA))
}

func (crm *ConcentrationRiskManager) calculatePostStressConcentration(data InstitutionData, scenario ConcentrationStressScenario) sdk.Dec {
	// Stub implementation
	return sdk.NewDecWithPrec(35, 2) // 35%
}

func (crm *ConcentrationRiskManager) stressTestSectorDownturn(data InstitutionData, report *ConcentrationRiskReport) ConcentrationStressScenario {
	// Stub implementation
	return ConcentrationStressScenario{
		ScenarioID:   "SECTOR_DOWNTURN",
		ScenarioName: "Real Estate Sector Crisis",
		Description:  "50% default in real estate sector",
		ExpectedLoss: sdk.NewCoin("usd", sdk.NewInt(500000000)),
		CapitalImpact: sdk.NewDecWithPrec(5, 2), // 5%
		RWAImpact:    sdk.NewDecWithPrec(20, 2), // 20%
	}
}

func (crm *ConcentrationRiskManager) stressTestGeographicCrisis(data InstitutionData, report *ConcentrationRiskReport) ConcentrationStressScenario {
	// Stub implementation
	return ConcentrationStressScenario{
		ScenarioID:   "GEOGRAPHIC_CRISIS",
		ScenarioName: "Emerging Market Crisis",
		Description:  "Currency crisis in emerging markets",
		ExpectedLoss: sdk.NewCoin("usd", sdk.NewInt(300000000)),
		CapitalImpact: sdk.NewDecWithPrec(3, 2), // 3%
		RWAImpact:    sdk.NewDecWithPrec(15, 2), // 15%
	}
}

func (crm *ConcentrationRiskManager) stressTestCorrelatedDefaults(data InstitutionData, report *ConcentrationRiskReport) ConcentrationStressScenario {
	// Stub implementation
	return ConcentrationStressScenario{
		ScenarioID:   "CORRELATED_DEFAULTS",
		ScenarioName: "Systemic Default Event",
		Description:  "Correlated defaults across sectors",
		ExpectedLoss: sdk.NewCoin("usd", sdk.NewInt(800000000)),
		CapitalImpact: sdk.NewDecWithPrec(8, 2), // 8%
		RWAImpact:    sdk.NewDecWithPrec(25, 2), // 25%
	}
}

func (crm *ConcentrationRiskManager) calculateAggregateStressImpact(scenarios []ConcentrationStressScenario) AggregateStressImpact {
	impact := AggregateStressImpact{
		TotalExpectedLoss: sdk.NewCoin("usd", sdk.ZeroInt()),
		MaxCapitalImpact:  sdk.ZeroDec(),
	}
	
	for _, scenario := range scenarios {
		impact.TotalExpectedLoss = impact.TotalExpectedLoss.Add(scenario.ExpectedLoss)
		if scenario.CapitalImpact.GT(impact.MaxCapitalImpact) {
			impact.MaxCapitalImpact = scenario.CapitalImpact
			impact.WorstCaseScenario = scenario.ScenarioName
		}
	}
	
	// Required buffer is max capital impact + 2% cushion
	impact.RequiredCapitalBuffer = sdk.NewCoin("usd", 
		sdk.NewInt(int64(impact.MaxCapitalImpact.Add(sdk.NewDecWithPrec(2, 2)).MustFloat64() * 1000000000)))
	
	return impact
}

func (crm *ConcentrationRiskManager) checkCompliance(report *ConcentrationRiskReport) *ConcentrationCompliance {
	compliance := &ConcentrationCompliance{
		IsCompliant:      true,
		ComplianceDate:   time.Now(),
		ViolationCount:   0,
		ComplianceViolations: []ComplianceViolation{},
		NextReviewDate:   time.Now().AddDate(0, 3, 0), // Quarterly review
	}
	
	// Check single borrower limit
	if !report.CreditConcentrations.SingleBorrowerLimit.IsCompliant {
		compliance.IsCompliant = false
		compliance.ViolationCount++
		for _, breach := range report.CreditConcentrations.SingleBorrowerLimit.Breaches {
			compliance.ComplianceViolations = append(compliance.ComplianceViolations, ComplianceViolation{
				ViolationID:    fmt.Sprintf("SBL_%s", breach.ExposureID),
				RequirementID:  "REG_SBL_001",
				ViolationType:  "SINGLE_BORROWER_LIMIT",
				Description:    fmt.Sprintf("Single borrower exposure exceeds 25%% limit: %s", breach.BorrowerName),
				Severity:       "HIGH",
				DateIdentified: time.Now(),
				Status:         "ACTIVE",
			})
		}
	}
	
	// Check large exposure breaches
	for _, breach := range report.LargeExposures.BreachAnalysis {
		compliance.IsCompliant = false
		compliance.ViolationCount++
		compliance.ComplianceViolations = append(compliance.ComplianceViolations, ComplianceViolation{
			ViolationID:    breach.BreachID,
			RequirementID:  "REG_LE_001",
			ViolationType:  "LARGE_EXPOSURE",
			Description:    fmt.Sprintf("Large exposure breach: %s exceeds limit", breach.CounterpartyName),
			Severity:       "CRITICAL",
			DateIdentified: breach.BreachDate,
			Status:         breach.Status,
		})
	}
	
	// Add regulatory requirements
	compliance.RegulatoryRequirements = []RegulatoryRequirement{
		{
			RequirementID:    "REG_SBL_001",
			Description:      "Single borrower exposure limit",
			LimitType:        "PERCENTAGE_OF_CAPITAL",
			LimitValue:       sdk.NewDecWithPrec(25, 2),
			CurrentValue:     report.CreditConcentrations.SingleBorrowerLimit.UtilizationRatio,
			ComplianceStatus: crm.getComplianceStatus(report.CreditConcentrations.SingleBorrowerLimit.IsCompliant),
		},
		{
			RequirementID:    "REG_GBL_001",
			Description:      "Group borrower exposure limit",
			LimitType:        "PERCENTAGE_OF_CAPITAL",
			LimitValue:       sdk.NewDecWithPrec(40, 2),
			CurrentValue:     report.CreditConcentrations.GroupBorrowerLimit.UtilizationRatio,
			ComplianceStatus: crm.getComplianceStatus(report.CreditConcentrations.GroupBorrowerLimit.IsCompliant),
		},
	}
	
	return compliance
}

func (crm *ConcentrationRiskManager) getComplianceStatus(isCompliant bool) string {
	if isCompliant {
		return "COMPLIANT"
	}
	return "NON_COMPLIANT"
}

func (crm *ConcentrationRiskManager) generateRecommendations(report *ConcentrationRiskReport) []RiskRecommendation {
	recommendations := []RiskRecommendation{}
	
	// Credit concentration recommendations
	if report.CreditConcentrations.HerfindahlIndex.GT(sdk.NewDecWithPrec(15, 2)) {
		recommendations = append(recommendations, RiskRecommendation{
			RecommendationID: "REC_CREDIT_001",
			Category:         "CREDIT_CONCENTRATION",
			Priority:         "HIGH",
			Description:      "Reduce portfolio concentration through diversification",
			Rationale:        fmt.Sprintf("Herfindahl Index of %.2f indicates high concentration", report.CreditConcentrations.HerfindahlIndex.MustFloat64()),
			Implementation:   "Target smaller ticket sizes and new customer segments",
			ExpectedOutcome:  "Reduce HHI below 0.15 within 12 months",
		})
	}
	
	// Sector recommendations
	for _, sector := range report.SectorConcentrations.SectorExposures {
		if sector.PercentOfTotal.GT(sdk.NewDecWithPrec(25, 2)) && sector.IsHighRisk {
			recommendations = append(recommendations, RiskRecommendation{
				RecommendationID: fmt.Sprintf("REC_SECTOR_%s", sector.SectorCode),
				Category:         "SECTOR_CONCENTRATION",
				Priority:         "HIGH",
				Description:      fmt.Sprintf("Reduce exposure to %s sector", sector.SectorName),
				Rationale:        fmt.Sprintf("High-risk sector exposure at %.1f%% exceeds prudent limits", sector.PercentOfTotal.MustFloat64()*100),
				Implementation:   "Implement sector cap and redirect lending to other sectors",
				ExpectedOutcome:  "Reduce sector exposure below 20% within 6 months",
			})
		}
	}
	
	// Geographic recommendations
	if report.GeographicConcentrations.CrossBorderExposure.IsPositive() {
		crossBorderPercent := sdk.NewDecFromInt(report.GeographicConcentrations.CrossBorderExposure.Amount).
			Quo(sdk.NewDecFromInt(report.GeographicConcentrations.CrossBorderExposure.Amount))
		
		if crossBorderPercent.GT(sdk.NewDecWithPrec(25, 2)) {
			recommendations = append(recommendations, RiskRecommendation{
				RecommendationID: "REC_GEO_001",
				Category:         "GEOGRAPHIC_CONCENTRATION",
				Priority:         "MEDIUM",
				Description:      "Diversify geographic exposure",
				Rationale:        "High cross-border concentration increases currency and political risk",
				Implementation:   "Focus on domestic opportunities and hedge foreign exposures",
				ExpectedOutcome:  "Reduce cross-border exposure below 20%",
			})
		}
	}
	
	// Large exposure recommendations
	if report.LargeExposures.NumberOfLargeExposures > 15 {
		recommendations = append(recommendations, RiskRecommendation{
			RecommendationID: "REC_LARGE_001",
			Category:         "LARGE_EXPOSURES",
			Priority:         "HIGH",
			Description:      "Reduce number of large exposures through syndication",
			Rationale:        fmt.Sprintf("%d large exposures create significant concentration risk", report.LargeExposures.NumberOfLargeExposures),
			Implementation:   "Syndicate excess portions of large exposures",
			ExpectedOutcome:  "Reduce large exposure count below 15",
		})
	}
	
	// Risk mitigation recommendations
	if report.RiskMetrics.OverallConcentrationScore > 70 {
		recommendations = append(recommendations, RiskRecommendation{
			RecommendationID: "REC_RISK_001",
			Category:         "RISK_MITIGATION",
			Priority:         "CRITICAL",
			Description:      "Implement comprehensive risk mitigation program",
			Rationale:        fmt.Sprintf("Overall concentration score of %.0f indicates critical risk levels", report.RiskMetrics.OverallConcentrationScore),
			Implementation:   "Execute all HIGH and CRITICAL priority mitigation strategies",
			ExpectedOutcome:  "Reduce overall concentration score below 60 within 6 months",
		})
	}
	
	return recommendations
}

func (crm *ConcentrationRiskManager) storeConcentrationRiskReport(ctx context.Context, report ConcentrationRiskReport) error {
	store := crm.keeper.storeService.OpenKVStore(ctx)
	key := []byte(fmt.Sprintf("concentration_risk_report_%s", report.ReportID))
	bz, err := json.Marshal(report)
	if err != nil {
		return err
	}
	store.Set(key, bz)
	return nil
}

func (crm *ConcentrationRiskManager) emitConcentrationRiskEvent(ctx context.Context, report *ConcentrationRiskReport) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"concentration_risk_analyzed",
			sdk.NewAttribute("report_id", report.ReportID),
			sdk.NewAttribute("institution_id", report.InstitutionID),
			sdk.NewAttribute("risk_level", report.RiskMetrics.RiskLevel),
			sdk.NewAttribute("concentration_score", fmt.Sprintf("%.1f", report.RiskMetrics.OverallConcentrationScore)),
			sdk.NewAttribute("is_compliant", fmt.Sprintf("%v", report.ComplianceStatus.IsCompliant)),
			sdk.NewAttribute("violation_count", fmt.Sprintf("%d", report.ComplianceStatus.ViolationCount)),
		),
	)
}

// Utility functions

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Public API methods

func (crm *ConcentrationRiskManager) GetConcentrationRiskReport(ctx context.Context, reportID string) (*ConcentrationRiskReport, error) {
	store := crm.keeper.storeService.OpenKVStore(ctx)
	key := []byte(fmt.Sprintf("concentration_risk_report_%s", reportID))
	bz := store.Get(key)
	if bz == nil {
		return nil, fmt.Errorf("report not found: %s", reportID)
	}
	
	var report ConcentrationRiskReport
	if err := json.Unmarshal(bz, &report); err != nil {
		return nil, fmt.Errorf("failed to unmarshal report: %w", err)
	}
	
	return &report, nil
}

func (crm *ConcentrationRiskManager) GetInstitutionReports(ctx context.Context, institutionID string) ([]ConcentrationRiskReport, error) {
	// Implementation would query all reports for an institution
	return []ConcentrationRiskReport{}, nil
}