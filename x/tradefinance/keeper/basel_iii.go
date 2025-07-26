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

	"github.com/DeshChain/DeshChain-Ecosystem/x/tradefinance/types"
)

// Basel III Capital Management System
type BaselIIICapitalManager struct {
	keeper                *Keeper
	capitalRequirements   CapitalRequirements
	riskWeights          RiskWeightConfig
	buffers              CapitalBuffers
	stressTestScenarios  []StressTestScenario
}

// NewBaselIIICapitalManager creates a new Basel III capital manager
func NewBaselIIICapitalManager(k *Keeper) *BaselIIICapitalManager {
	return &BaselIIICapitalManager{
		keeper:              k,
		capitalRequirements: initializeCapitalRequirements(),
		riskWeights:        initializeRiskWeights(),
		buffers:            initializeCapitalBuffers(),
		stressTestScenarios: initializeStressTestScenarios(),
	}
}

// Core Basel III structures

type CapitalAdequacyReport struct {
	ReportID               string                     `json:"report_id"`
	ReportDate             time.Time                  `json:"report_date"`
	BankID                 string                     `json:"bank_id"`
	CapitalComponents      CapitalComponents          `json:"capital_components"`
	RiskWeightedAssets     RiskWeightedAssets         `json:"risk_weighted_assets"`
	CapitalRatios          CapitalRatios              `json:"capital_ratios"`
	BufferRequirements     BufferRequirements         `json:"buffer_requirements"`
	LeverageRatio          LeverageRatioCalculation   `json:"leverage_ratio"`
	LiquidityMetrics       LiquidityMetrics           `json:"liquidity_metrics"`
	ComplianceStatus       ComplianceStatus           `json:"compliance_status"`
	StressTestResults      []StressTestResult         `json:"stress_test_results"`
	CapitalPlan            CapitalPlan                `json:"capital_plan"`
	RegulatoryAdjustments  []RegulatoryAdjustment     `json:"regulatory_adjustments"`
	Metadata               map[string]interface{}     `json:"metadata"`
}

type CapitalComponents struct {
	// Tier 1 Capital
	CommonEquityTier1      CommonEquityTier1Capital   `json:"common_equity_tier1"`
	AdditionalTier1        AdditionalTier1Capital     `json:"additional_tier1"`
	TotalTier1             sdk.Coin                   `json:"total_tier1"`
	
	// Tier 2 Capital
	Tier2Capital           Tier2Capital               `json:"tier2_capital"`
	
	// Total Capital
	TotalCapital           sdk.Coin                   `json:"total_capital"`
	
	// Deductions
	RegulatoryDeductions   []CapitalDeduction         `json:"regulatory_deductions"`
	NetCapital             sdk.Coin                   `json:"net_capital"`
}

type CommonEquityTier1Capital struct {
	PaidUpCapital          sdk.Coin                   `json:"paid_up_capital"`
	SharePremium           sdk.Coin                   `json:"share_premium"`
	RetainedEarnings       sdk.Coin                   `json:"retained_earnings"`
	AccumulatedOCI         sdk.Coin                   `json:"accumulated_oci"`
	StatutoryReserves      sdk.Coin                   `json:"statutory_reserves"`
	MinorityInterest       sdk.Coin                   `json:"minority_interest"`
	RegulatoryAdjustments  []CapitalAdjustment        `json:"regulatory_adjustments"`
	Total                  sdk.Coin                   `json:"total"`
	Percentage             sdk.Dec                    `json:"percentage_of_rwa"`
}

type AdditionalTier1Capital struct {
	PerpetualBonds         sdk.Coin                   `json:"perpetual_bonds"`
	PreferenceShares       sdk.Coin                   `json:"preference_shares"`
	ContingentConvertible  sdk.Coin                   `json:"contingent_convertible"`
	MinorityInterest       sdk.Coin                   `json:"minority_interest"`
	RegulatoryAdjustments  []CapitalAdjustment        `json:"regulatory_adjustments"`
	Total                  sdk.Coin                   `json:"total"`
}

type Tier2Capital struct {
	SubordinatedDebt       sdk.Coin                   `json:"subordinated_debt"`
	RevaluationReserves    sdk.Coin                   `json:"revaluation_reserves"`
	GeneralProvisions      sdk.Coin                   `json:"general_provisions"`
	HybridInstruments      sdk.Coin                   `json:"hybrid_instruments"`
	MinorityInterest       sdk.Coin                   `json:"minority_interest"`
	RegulatoryAdjustments  []CapitalAdjustment        `json:"regulatory_adjustments"`
	Total                  sdk.Coin                   `json:"total"`
}

type RiskWeightedAssets struct {
	CreditRiskRWA          CreditRiskRWA              `json:"credit_risk_rwa"`
	MarketRiskRWA          MarketRiskRWA              `json:"market_risk_rwa"`
	OperationalRiskRWA     OperationalRiskRWA         `json:"operational_risk_rwa"`
	CVARiskRWA             CVARiskRWA                 `json:"cva_risk_rwa"`
	TotalRWA               sdk.Coin                   `json:"total_rwa"`
	RWADensity             sdk.Dec                    `json:"rwa_density"`
}

type CreditRiskRWA struct {
	OnBalanceSheet         []AssetRiskWeight          `json:"on_balance_sheet"`
	OffBalanceSheet        []AssetRiskWeight          `json:"off_balance_sheet"`
	Derivatives            []DerivativeExposure       `json:"derivatives"`
	SecuritiesFinancing    []SFTExposure              `json:"securities_financing"`
	CentralCounterparties  []CCPExposure              `json:"central_counterparties"`
	TotalCreditRWA         sdk.Coin                   `json:"total_credit_rwa"`
	StandardizedApproach   bool                       `json:"standardized_approach"`
	IRBApproach            IRBApproachDetails         `json:"irb_approach,omitempty"`
}

type AssetRiskWeight struct {
	AssetClass             AssetClass                 `json:"asset_class"`
	ExposureAmount         sdk.Coin                   `json:"exposure_amount"`
	RiskWeight             sdk.Dec                    `json:"risk_weight"`
	RiskWeightedAmount     sdk.Coin                   `json:"risk_weighted_amount"`
	CreditRating           string                     `json:"credit_rating,omitempty"`
	Maturity               time.Duration              `json:"maturity,omitempty"`
	Collateral             CollateralDetails          `json:"collateral,omitempty"`
	CCF                    sdk.Dec                    `json:"credit_conversion_factor,omitempty"`
}

type MarketRiskRWA struct {
	InterestRateRisk       sdk.Coin                   `json:"interest_rate_risk"`
	EquityPositionRisk     sdk.Coin                   `json:"equity_position_risk"`
	ForeignExchangeRisk    sdk.Coin                   `json:"foreign_exchange_risk"`
	CommodityRisk          sdk.Coin                   `json:"commodity_risk"`
	OptionsRisk            sdk.Coin                   `json:"options_risk"`
	TotalMarketRWA         sdk.Coin                   `json:"total_market_rwa"`
	InternalModelsApproach bool                       `json:"internal_models_approach"`
	VaRModel               VaRModelDetails            `json:"var_model,omitempty"`
}

type OperationalRiskRWA struct {
	BasicIndicatorApproach bool                       `json:"basic_indicator_approach"`
	StandardizedApproach   bool                       `json:"standardized_approach"`
	AdvancedMeasurement    bool                       `json:"advanced_measurement"`
	GrossIncome            []YearlyGrossIncome        `json:"gross_income"`
	BusinessLineRWA        []BusinessLineRWA          `json:"business_line_rwa,omitempty"`
	TotalOperationalRWA    sdk.Coin                   `json:"total_operational_rwa"`
}

type CapitalRatios struct {
	CET1Ratio              sdk.Dec                    `json:"cet1_ratio"`
	Tier1Ratio             sdk.Dec                    `json:"tier1_ratio"`
	TotalCapitalRatio      sdk.Dec                    `json:"total_capital_ratio"`
	LeverageRatio          sdk.Dec                    `json:"leverage_ratio"`
	RequiredCET1           sdk.Dec                    `json:"required_cet1"`
	RequiredTier1          sdk.Dec                    `json:"required_tier1"`
	RequiredTotal          sdk.Dec                    `json:"required_total"`
	CET1Surplus            sdk.Coin                   `json:"cet1_surplus"`
	Tier1Surplus           sdk.Coin                   `json:"tier1_surplus"`
	TotalSurplus           sdk.Coin                   `json:"total_surplus"`
}

type BufferRequirements struct {
	CapitalConservation    BufferRequirement          `json:"capital_conservation_buffer"`
	CounterCyclical        BufferRequirement          `json:"counter_cyclical_buffer"`
	SystemicRisk           BufferRequirement          `json:"systemic_risk_buffer"`
	GSIB                   BufferRequirement          `json:"gsib_buffer"`
	DSIB                   BufferRequirement          `json:"dsib_buffer"`
	TotalBufferRequirement sdk.Dec                    `json:"total_buffer_requirement"`
	AvailableBuffer        sdk.Dec                    `json:"available_buffer"`
	BufferBreach           bool                       `json:"buffer_breach"`
	RestrictionsApplied    []BufferRestriction        `json:"restrictions_applied"`
}

type LiquidityMetrics struct {
	LCR                    LiquidityCoverageRatio     `json:"liquidity_coverage_ratio"`
	NSFR                   NetStableFundingRatio      `json:"net_stable_funding_ratio"`
	LiquidityBuffers       []LiquidityBuffer          `json:"liquidity_buffers"`
	StressedOutflows       sdk.Coin                   `json:"stressed_outflows"`
	ContingencyFunding     ContingencyFundingPlan     `json:"contingency_funding"`
}

// Calculation functions

// CalculateCapitalAdequacy performs complete Basel III capital adequacy calculation
func (bcm *BaselIIICapitalManager) CalculateCapitalAdequacy(ctx context.Context, bankID string) (*CapitalAdequacyReport, error) {
	report := &CapitalAdequacyReport{
		ReportID:   fmt.Sprintf("CAR_%s_%d", bankID, time.Now().Unix()),
		ReportDate: time.Now(),
		BankID:     bankID,
		Metadata:   make(map[string]interface{}),
	}

	// Step 1: Calculate capital components
	capitalComponents, err := bcm.calculateCapitalComponents(ctx, bankID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate capital components: %w", err)
	}
	report.CapitalComponents = *capitalComponents

	// Step 2: Calculate risk-weighted assets
	rwa, err := bcm.calculateRiskWeightedAssets(ctx, bankID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate RWA: %w", err)
	}
	report.RiskWeightedAssets = *rwa

	// Step 3: Calculate capital ratios
	capitalRatios := bcm.calculateCapitalRatios(capitalComponents, rwa)
	report.CapitalRatios = *capitalRatios

	// Step 4: Calculate buffer requirements
	bufferReqs := bcm.calculateBufferRequirements(ctx, bankID, capitalRatios)
	report.BufferRequirements = *bufferReqs

	// Step 5: Calculate leverage ratio
	leverageRatio, err := bcm.calculateLeverageRatio(ctx, bankID, capitalComponents)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate leverage ratio: %w", err)
	}
	report.LeverageRatio = *leverageRatio

	// Step 6: Calculate liquidity metrics
	liquidityMetrics, err := bcm.calculateLiquidityMetrics(ctx, bankID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate liquidity metrics: %w", err)
	}
	report.LiquidityMetrics = *liquidityMetrics

	// Step 7: Determine compliance status
	complianceStatus := bcm.determineComplianceStatus(capitalRatios, bufferReqs, leverageRatio, liquidityMetrics)
	report.ComplianceStatus = *complianceStatus

	// Step 8: Run stress tests
	stressResults := bcm.runStressTests(ctx, bankID, capitalComponents, rwa)
	report.StressTestResults = stressResults

	// Step 9: Generate capital plan
	capitalPlan := bcm.generateCapitalPlan(capitalRatios, bufferReqs, stressResults)
	report.CapitalPlan = *capitalPlan

	// Store report
	if err := bcm.storeCapitalAdequacyReport(ctx, *report); err != nil {
		return report, fmt.Errorf("failed to store report: %w", err)
	}

	// Emit event
	bcm.emitCapitalAdequacyEvent(ctx, report)

	return report, nil
}

// Capital component calculations

func (bcm *BaselIIICapitalManager) calculateCapitalComponents(ctx context.Context, bankID string) (*CapitalComponents, error) {
	// Fetch bank's capital data
	bankData, err := bcm.getBankCapitalData(ctx, bankID)
	if err != nil {
		return nil, err
	}

	// Calculate CET1
	cet1 := bcm.calculateCET1Capital(bankData)
	
	// Calculate Additional Tier 1
	at1 := bcm.calculateAdditionalTier1(bankData)
	
	// Calculate Tier 2
	tier2 := bcm.calculateTier2Capital(bankData)
	
	// Apply regulatory deductions
	deductions := bcm.calculateRegulatoryDeductions(bankData)
	
	// Calculate totals
	totalTier1 := cet1.Total.Add(at1.Total)
	totalCapital := totalTier1.Add(tier2.Total)
	
	// Apply deductions
	netCapital := totalCapital
	for _, deduction := range deductions {
		netCapital = netCapital.Sub(deduction.Amount)
	}

	return &CapitalComponents{
		CommonEquityTier1:    cet1,
		AdditionalTier1:      at1,
		TotalTier1:          totalTier1,
		Tier2Capital:        tier2,
		TotalCapital:        totalCapital,
		RegulatoryDeductions: deductions,
		NetCapital:          netCapital,
	}, nil
}

func (bcm *BaselIIICapitalManager) calculateCET1Capital(bankData BankCapitalData) CommonEquityTier1Capital {
	cet1 := CommonEquityTier1Capital{
		PaidUpCapital:     bankData.PaidUpCapital,
		SharePremium:      bankData.SharePremium,
		RetainedEarnings:  bankData.RetainedEarnings,
		AccumulatedOCI:    bankData.AccumulatedOCI,
		StatutoryReserves: bankData.StatutoryReserves,
		MinorityInterest:  bcm.calculateMinorityInterestCET1(bankData),
	}

	// Apply regulatory adjustments
	adjustments := bcm.getCET1RegulatoryAdjustments(bankData)
	cet1.RegulatoryAdjustments = adjustments

	// Calculate total
	total := cet1.PaidUpCapital.
		Add(cet1.SharePremium).
		Add(cet1.RetainedEarnings).
		Add(cet1.AccumulatedOCI).
		Add(cet1.StatutoryReserves).
		Add(cet1.MinorityInterest)

	// Apply adjustments
	for _, adj := range adjustments {
		if adj.IsDeduction {
			total = total.Sub(adj.Amount)
		} else {
			total = total.Add(adj.Amount)
		}
	}

	cet1.Total = total
	return cet1
}

// Risk-weighted assets calculations

func (bcm *BaselIIICapitalManager) calculateRiskWeightedAssets(ctx context.Context, bankID string) (*RiskWeightedAssets, error) {
	// Calculate credit risk RWA
	creditRWA, err := bcm.calculateCreditRiskRWA(ctx, bankID)
	if err != nil {
		return nil, fmt.Errorf("credit risk RWA calculation failed: %w", err)
	}

	// Calculate market risk RWA
	marketRWA, err := bcm.calculateMarketRiskRWA(ctx, bankID)
	if err != nil {
		return nil, fmt.Errorf("market risk RWA calculation failed: %w", err)
	}

	// Calculate operational risk RWA
	operationalRWA, err := bcm.calculateOperationalRiskRWA(ctx, bankID)
	if err != nil {
		return nil, fmt.Errorf("operational risk RWA calculation failed: %w", err)
	}

	// Calculate CVA risk RWA
	cvaRWA, err := bcm.calculateCVARiskRWA(ctx, bankID)
	if err != nil {
		return nil, fmt.Errorf("CVA risk RWA calculation failed: %w", err)
	}

	// Calculate total RWA
	totalRWA := creditRWA.TotalCreditRWA.
		Add(marketRWA.TotalMarketRWA).
		Add(operationalRWA.TotalOperationalRWA).
		Add(cvaRWA.TotalCVARWA)

	// Calculate RWA density
	totalAssets, _ := bcm.getTotalAssets(ctx, bankID)
	rwaDensity := sdk.NewDecFromInt(totalRWA.Amount).Quo(sdk.NewDecFromInt(totalAssets.Amount))

	return &RiskWeightedAssets{
		CreditRiskRWA:      *creditRWA,
		MarketRiskRWA:      *marketRWA,
		OperationalRiskRWA: *operationalRWA,
		CVARiskRWA:         *cvaRWA,
		TotalRWA:           totalRWA,
		RWADensity:         rwaDensity,
	}, nil
}

func (bcm *BaselIIICapitalManager) calculateCreditRiskRWA(ctx context.Context, bankID string) (*CreditRiskRWA, error) {
	creditRWA := &CreditRiskRWA{
		OnBalanceSheet:  []AssetRiskWeight{},
		OffBalanceSheet: []AssetRiskWeight{},
		Derivatives:     []DerivativeExposure{},
		StandardizedApproach: true,
	}

	// Get bank's credit exposures
	exposures, err := bcm.getBankCreditExposures(ctx, bankID)
	if err != nil {
		return nil, err
	}

	// Calculate on-balance sheet RWA
	for _, exposure := range exposures.OnBalanceSheet {
		riskWeight := bcm.getRiskWeight(exposure)
		rwa := AssetRiskWeight{
			AssetClass:         exposure.AssetClass,
			ExposureAmount:     exposure.Amount,
			RiskWeight:         riskWeight,
			RiskWeightedAmount: sdk.NewCoin(exposure.Amount.Denom, exposure.Amount.Amount.Mul(sdk.NewInt(int64(riskWeight.MustFloat64()*100))).Quo(sdk.NewInt(100))),
			CreditRating:       exposure.Rating,
			Maturity:          exposure.Maturity,
		}
		creditRWA.OnBalanceSheet = append(creditRWA.OnBalanceSheet, rwa)
	}

	// Calculate off-balance sheet RWA with CCF
	for _, exposure := range exposures.OffBalanceSheet {
		ccf := bcm.getCreditConversionFactor(exposure)
		adjustedExposure := sdk.NewCoin(
			exposure.Amount.Denom,
			exposure.Amount.Amount.Mul(sdk.NewInt(int64(ccf.MustFloat64()*100))).Quo(sdk.NewInt(100)),
		)
		riskWeight := bcm.getRiskWeight(exposure)
		
		rwa := AssetRiskWeight{
			AssetClass:         exposure.AssetClass,
			ExposureAmount:     exposure.Amount,
			RiskWeight:         riskWeight,
			RiskWeightedAmount: sdk.NewCoin(adjustedExposure.Denom, adjustedExposure.Amount.Mul(sdk.NewInt(int64(riskWeight.MustFloat64()*100))).Quo(sdk.NewInt(100))),
			CCF:               ccf,
		}
		creditRWA.OffBalanceSheet = append(creditRWA.OffBalanceSheet, rwa)
	}

	// Calculate total credit RWA
	totalCreditRWA := sdk.NewCoin(bcm.getBaseCurrency(), sdk.ZeroInt())
	for _, rwa := range creditRWA.OnBalanceSheet {
		totalCreditRWA = totalCreditRWA.Add(rwa.RiskWeightedAmount)
	}
	for _, rwa := range creditRWA.OffBalanceSheet {
		totalCreditRWA = totalCreditRWA.Add(rwa.RiskWeightedAmount)
	}

	creditRWA.TotalCreditRWA = totalCreditRWA
	return creditRWA, nil
}

// Capital ratio calculations

func (bcm *BaselIIICapitalManager) calculateCapitalRatios(capital *CapitalComponents, rwa *RiskWeightedAssets) *CapitalRatios {
	// Calculate ratios
	cet1Ratio := sdk.NewDecFromInt(capital.CommonEquityTier1.Total.Amount).Quo(sdk.NewDecFromInt(rwa.TotalRWA.Amount))
	tier1Ratio := sdk.NewDecFromInt(capital.TotalTier1.Amount).Quo(sdk.NewDecFromInt(rwa.TotalRWA.Amount))
	totalRatio := sdk.NewDecFromInt(capital.TotalCapital.Amount).Quo(sdk.NewDecFromInt(rwa.TotalRWA.Amount))

	// Get regulatory minimums
	requiredCET1 := bcm.capitalRequirements.MinimumCET1Ratio
	requiredTier1 := bcm.capitalRequirements.MinimumTier1Ratio
	requiredTotal := bcm.capitalRequirements.MinimumTotalRatio

	// Calculate surplus/deficit
	cet1Surplus := capital.CommonEquityTier1.Total.Sub(
		sdk.NewCoin(capital.CommonEquityTier1.Total.Denom, 
			rwa.TotalRWA.Amount.Mul(sdk.NewInt(int64(requiredCET1.MustFloat64()*10000))).Quo(sdk.NewInt(10000)),
		),
	)
	
	tier1Surplus := capital.TotalTier1.Sub(
		sdk.NewCoin(capital.TotalTier1.Denom,
			rwa.TotalRWA.Amount.Mul(sdk.NewInt(int64(requiredTier1.MustFloat64()*10000))).Quo(sdk.NewInt(10000)),
		),
	)
	
	totalSurplus := capital.TotalCapital.Sub(
		sdk.NewCoin(capital.TotalCapital.Denom,
			rwa.TotalRWA.Amount.Mul(sdk.NewInt(int64(requiredTotal.MustFloat64()*10000))).Quo(sdk.NewInt(10000)),
		),
	)

	return &CapitalRatios{
		CET1Ratio:         cet1Ratio,
		Tier1Ratio:        tier1Ratio,
		TotalCapitalRatio: totalRatio,
		RequiredCET1:      requiredCET1,
		RequiredTier1:     requiredTier1,
		RequiredTotal:     requiredTotal,
		CET1Surplus:       cet1Surplus,
		Tier1Surplus:      tier1Surplus,
		TotalSurplus:      totalSurplus,
	}
}

// Buffer calculations

func (bcm *BaselIIICapitalManager) calculateBufferRequirements(ctx context.Context, bankID string, ratios *CapitalRatios) *BufferRequirements {
	bufferReqs := &BufferRequirements{}

	// Capital conservation buffer (2.5%)
	bufferReqs.CapitalConservation = BufferRequirement{
		BufferType:    "CAPITAL_CONSERVATION",
		RequiredRatio: sdk.NewDecWithPrec(25, 3), // 2.5%
		ActualRatio:   sdk.MaxDec(sdk.ZeroDec(), ratios.CET1Ratio.Sub(bcm.capitalRequirements.MinimumCET1Ratio)),
		IsMet:         ratios.CET1Ratio.GTE(bcm.capitalRequirements.MinimumCET1Ratio.Add(sdk.NewDecWithPrec(25, 3))),
	}

	// Counter-cyclical buffer (0-2.5% based on credit growth)
	ccyBuffer := bcm.calculateCounterCyclicalBuffer(ctx, bankID)
	bufferReqs.CounterCyclical = ccyBuffer

	// Systemic risk buffer (if applicable)
	if bcm.isSystemicallyImportant(ctx, bankID) {
		bufferReqs.SystemicRisk = BufferRequirement{
			BufferType:    "SYSTEMIC_RISK",
			RequiredRatio: sdk.NewDecWithPrec(10, 3), // 1.0%
			ActualRatio:   sdk.MaxDec(sdk.ZeroDec(), ratios.CET1Ratio.Sub(bcm.capitalRequirements.MinimumCET1Ratio.Add(sdk.NewDecWithPrec(25, 3)))),
			IsMet:         ratios.CET1Ratio.GTE(bcm.capitalRequirements.MinimumCET1Ratio.Add(sdk.NewDecWithPrec(35, 3))),
		}
	}

	// DSIB buffer (for domestic systemically important banks)
	if bcm.isDSIB(ctx, bankID) {
		dsibScore := bcm.calculateDSIBScore(ctx, bankID)
		dsibBuffer := bcm.getDSIBBufferRequirement(dsibScore)
		bufferReqs.DSIB = BufferRequirement{
			BufferType:    "DSIB",
			RequiredRatio: dsibBuffer,
			ActualRatio:   sdk.MaxDec(sdk.ZeroDec(), ratios.CET1Ratio.Sub(bcm.capitalRequirements.MinimumCET1Ratio.Add(sdk.NewDecWithPrec(25, 3)))),
			IsMet:         ratios.CET1Ratio.GTE(bcm.capitalRequirements.MinimumCET1Ratio.Add(sdk.NewDecWithPrec(25, 3)).Add(dsibBuffer)),
		}
	}

	// Calculate total buffer requirement
	totalBuffer := bufferReqs.CapitalConservation.RequiredRatio
	if bufferReqs.CounterCyclical.RequiredRatio.IsPositive() {
		totalBuffer = totalBuffer.Add(bufferReqs.CounterCyclical.RequiredRatio)
	}
	if bufferReqs.SystemicRisk.RequiredRatio.IsPositive() {
		totalBuffer = totalBuffer.Add(bufferReqs.SystemicRisk.RequiredRatio)
	}
	if bufferReqs.DSIB.RequiredRatio.IsPositive() {
		totalBuffer = totalBuffer.Add(bufferReqs.DSIB.RequiredRatio)
	}
	bufferReqs.TotalBufferRequirement = totalBuffer

	// Calculate available buffer
	bufferReqs.AvailableBuffer = sdk.MaxDec(sdk.ZeroDec(), ratios.CET1Ratio.Sub(bcm.capitalRequirements.MinimumCET1Ratio))

	// Check for buffer breach
	bufferReqs.BufferBreach = bufferReqs.AvailableBuffer.LT(bufferReqs.TotalBufferRequirement)

	// Apply restrictions if buffer breached
	if bufferReqs.BufferBreach {
		bufferReqs.RestrictionsApplied = bcm.getBufferBreachRestrictions(bufferReqs.AvailableBuffer, bufferReqs.TotalBufferRequirement)
	}

	return bufferReqs
}

// Leverage ratio calculation

func (bcm *BaselIIICapitalManager) calculateLeverageRatio(ctx context.Context, bankID string, capital *CapitalComponents) (*LeverageRatioCalculation, error) {
	// Get total exposure measure
	exposureMeasure, err := bcm.calculateLeverageExposure(ctx, bankID)
	if err != nil {
		return nil, err
	}

	// Calculate ratio
	leverageRatio := sdk.NewDecFromInt(capital.TotalTier1.Amount).Quo(sdk.NewDecFromInt(exposureMeasure.TotalExposure.Amount))

	// Get regulatory minimum (3%)
	minimumRatio := sdk.NewDecWithPrec(3, 2)

	return &LeverageRatioCalculation{
		Tier1Capital:        capital.TotalTier1,
		ExposureMeasure:     *exposureMeasure,
		LeverageRatio:       leverageRatio,
		MinimumRequirement:  minimumRatio,
		Buffer:             leverageRatio.Sub(minimumRatio),
		IsCompliant:        leverageRatio.GTE(minimumRatio),
	}, nil
}

// Liquidity calculations

func (bcm *BaselIIICapitalManager) calculateLiquidityMetrics(ctx context.Context, bankID string) (*LiquidityMetrics, error) {
	metrics := &LiquidityMetrics{}

	// Calculate LCR
	lcr, err := bcm.calculateLCR(ctx, bankID)
	if err != nil {
		return nil, err
	}
	metrics.LCR = *lcr

	// Calculate NSFR
	nsfr, err := bcm.calculateNSFR(ctx, bankID)
	if err != nil {
		return nil, err
	}
	metrics.NSFR = *nsfr

	// Get liquidity buffers
	buffers, err := bcm.getLiquidityBuffers(ctx, bankID)
	if err != nil {
		return nil, err
	}
	metrics.LiquidityBuffers = buffers

	// Calculate stressed outflows
	stressedOutflows := bcm.calculateStressedOutflows(ctx, bankID)
	metrics.StressedOutflows = stressedOutflows

	// Get contingency funding plan
	cfp := bcm.getContingencyFundingPlan(ctx, bankID)
	metrics.ContingencyFunding = cfp

	return metrics, nil
}

func (bcm *BaselIIICapitalManager) calculateLCR(ctx context.Context, bankID string) (*LiquidityCoverageRatio, error) {
	// Get high-quality liquid assets
	hqla, err := bcm.getHQLA(ctx, bankID)
	if err != nil {
		return nil, err
	}

	// Calculate net cash outflows
	netOutflows, err := bcm.calculateNetCashOutflows(ctx, bankID)
	if err != nil {
		return nil, err
	}

	// Calculate LCR
	lcr := sdk.NewDecFromInt(hqla.Total.Amount).Quo(sdk.NewDecFromInt(netOutflows.Amount))

	return &LiquidityCoverageRatio{
		HQLA:               *hqla,
		TotalNetOutflows:   netOutflows,
		LCR:               lcr,
		MinimumRequirement: sdk.NewDecWithPrec(100, 2), // 100%
		IsCompliant:       lcr.GTE(sdk.NewDecWithPrec(100, 2)),
		ObservationPeriod: 30 * 24 * time.Hour,
	}, nil
}

// Stress testing

func (bcm *BaselIIICapitalManager) runStressTests(ctx context.Context, bankID string, capital *CapitalComponents, rwa *RiskWeightedAssets) []StressTestResult {
	results := []StressTestResult{}

	for _, scenario := range bcm.stressTestScenarios {
		result := bcm.runStressScenario(ctx, bankID, capital, rwa, scenario)
		results = append(results, result)
	}

	return results
}

func (bcm *BaselIIICapitalManager) runStressScenario(ctx context.Context, bankID string, capital *CapitalComponents, rwa *RiskWeightedAssets, scenario StressTestScenario) StressTestResult {
	result := StressTestResult{
		ScenarioID:   scenario.ID,
		ScenarioName: scenario.Name,
		RunDate:      time.Now(),
	}

	// Apply stress to capital
	stressedCapital := bcm.applyCapitalStress(capital, scenario.CapitalImpacts)
	result.StressedCapital = stressedCapital

	// Apply stress to RWA
	stressedRWA := bcm.applyRWAStress(rwa, scenario.RWAImpacts)
	result.StressedRWA = stressedRWA

	// Recalculate ratios under stress
	stressedRatios := bcm.calculateCapitalRatios(&stressedCapital, &stressedRWA)
	result.StressedRatios = *stressedRatios

	// Determine if bank passes stress test
	result.PassesTest = stressedRatios.CET1Ratio.GTE(scenario.MinimumCET1) &&
		stressedRatios.Tier1Ratio.GTE(scenario.MinimumTier1) &&
		stressedRatios.TotalCapitalRatio.GTE(scenario.MinimumTotal)

	// Calculate capital shortfall if any
	if !result.PassesTest {
		shortfall := bcm.calculateCapitalShortfall(stressedCapital, stressedRWA, scenario)
		result.CapitalShortfall = shortfall
	}

	return result
}

// Compliance determination

func (bcm *BaselIIICapitalManager) determineComplianceStatus(ratios *CapitalRatios, buffers *BufferRequirements, leverage *LeverageRatioCalculation, liquidity *LiquidityMetrics) *ComplianceStatus {
	status := &ComplianceStatus{
		IsCompliant:      true,
		ComplianceDate:   time.Now(),
		ComplianceItems:  []ComplianceItem{},
		RequiredActions:  []RequiredAction{},
	}

	// Check capital ratios
	if ratios.CET1Ratio.LT(ratios.RequiredCET1) {
		status.IsCompliant = false
		status.ComplianceItems = append(status.ComplianceItems, ComplianceItem{
			Requirement: "Minimum CET1 Ratio",
			Required:    ratios.RequiredCET1.String(),
			Actual:      ratios.CET1Ratio.String(),
			Status:      "NON_COMPLIANT",
		})
		status.RequiredActions = append(status.RequiredActions, RequiredAction{
			ActionType:  "CAPITAL_RAISE",
			Description: "Raise additional CET1 capital",
			Amount:      ratios.CET1Surplus.Neg(),
			Deadline:    time.Now().AddDate(0, 3, 0),
		})
	}

	// Check buffers
	if buffers.BufferBreach {
		status.IsCompliant = false
		status.ComplianceItems = append(status.ComplianceItems, ComplianceItem{
			Requirement: "Capital Buffers",
			Required:    buffers.TotalBufferRequirement.String(),
			Actual:      buffers.AvailableBuffer.String(),
			Status:      "BUFFER_BREACH",
		})
	}

	// Check leverage ratio
	if !leverage.IsCompliant {
		status.IsCompliant = false
		status.ComplianceItems = append(status.ComplianceItems, ComplianceItem{
			Requirement: "Leverage Ratio",
			Required:    leverage.MinimumRequirement.String(),
			Actual:      leverage.LeverageRatio.String(),
			Status:      "NON_COMPLIANT",
		})
	}

	// Check liquidity
	if !liquidity.LCR.IsCompliant {
		status.IsCompliant = false
		status.ComplianceItems = append(status.ComplianceItems, ComplianceItem{
			Requirement: "Liquidity Coverage Ratio",
			Required:    liquidity.LCR.MinimumRequirement.String(),
			Actual:      liquidity.LCR.LCR.String(),
			Status:      "NON_COMPLIANT",
		})
	}

	return status
}

// Capital planning

func (bcm *BaselIIICapitalManager) generateCapitalPlan(ratios *CapitalRatios, buffers *BufferRequirements, stressResults []StressTestResult) *CapitalPlan {
	plan := &CapitalPlan{
		PlanningHorizon: 3 * 365 * 24 * time.Hour, // 3 years
		TargetRatios: TargetRatios{
			CET1:  ratios.RequiredCET1.Add(buffers.TotalBufferRequirement).Add(sdk.NewDecWithPrec(10, 3)), // +1% management buffer
			Tier1: ratios.RequiredTier1.Add(buffers.TotalBufferRequirement).Add(sdk.NewDecWithPrec(10, 3)),
			Total: ratios.RequiredTotal.Add(buffers.TotalBufferRequirement).Add(sdk.NewDecWithPrec(10, 3)),
		},
		PlannedActions: []PlannedCapitalAction{},
	}

	// Determine capital needs based on stress tests
	maxShortfall := sdk.NewCoin(bcm.getBaseCurrency(), sdk.ZeroInt())
	for _, result := range stressResults {
		if result.CapitalShortfall.IsPositive() && result.CapitalShortfall.IsGT(maxShortfall) {
			maxShortfall = result.CapitalShortfall
		}
	}

	// Plan capital actions if needed
	if maxShortfall.IsPositive() || ratios.CET1Ratio.LT(plan.TargetRatios.CET1) {
		// Plan capital raise
		plan.PlannedActions = append(plan.PlannedActions, PlannedCapitalAction{
			ActionType:   "EQUITY_RAISE",
			Amount:       maxShortfall.Add(sdk.NewCoin(maxShortfall.Denom, maxShortfall.Amount.Quo(sdk.NewInt(10)))), // +10% buffer
			Timeline:     time.Now().AddDate(0, 6, 0),
			Impact:       "Increase CET1 capital",
			Status:       "PLANNED",
		})
	}

	// Plan dividend policy
	if ratios.CET1Ratio.GT(plan.TargetRatios.CET1) {
		// Can pay dividends
		plan.DividendPolicy = "NORMAL"
		plan.MaxDividendPayout = sdk.NewDecWithPrec(40, 2) // 40% payout ratio
	} else {
		// Restrict dividends
		plan.DividendPolicy = "RESTRICTED"
		plan.MaxDividendPayout = sdk.ZeroDec()
	}

	return plan
}

// Helper functions and data structures

type BankCapitalData struct {
	BankID              string
	PaidUpCapital       sdk.Coin
	SharePremium        sdk.Coin
	RetainedEarnings    sdk.Coin
	AccumulatedOCI      sdk.Coin
	StatutoryReserves   sdk.Coin
	SubordinatedDebt    sdk.Coin
	RevaluationReserves sdk.Coin
	GeneralProvisions   sdk.Coin
	MinorityInterests   []MinorityInterest
	Deductions          []Deduction
	TotalAssets         sdk.Coin
	TotalLiabilities    sdk.Coin
}

type MinorityInterest struct {
	Entity     string   `json:"entity"`
	Amount     sdk.Coin `json:"amount"`
	Percentage sdk.Dec  `json:"percentage"`
}

type Deduction struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Amount      sdk.Coin `json:"amount"`
}

type CreditExposures struct {
	OnBalanceSheet  []CreditExposure
	OffBalanceSheet []CreditExposure
}

type CreditExposure struct {
	ExposureID  string
	AssetClass  AssetClass
	Amount      sdk.Coin
	Rating      string
	Maturity    time.Duration
	Collateral  bool
	Guarantor   string
	Country     string
}

type AssetClass int

const (
	ASSET_CLASS_SOVEREIGN AssetClass = iota
	ASSET_CLASS_BANK
	ASSET_CLASS_CORPORATE
	ASSET_CLASS_RETAIL_MORTGAGE
	ASSET_CLASS_RETAIL_REVOLVING
	ASSET_CLASS_RETAIL_OTHER
	ASSET_CLASS_COMMERCIAL_REAL_ESTATE
	ASSET_CLASS_EQUITY
	ASSET_CLASS_OTHER
)

type CapitalRequirements struct {
	MinimumCET1Ratio  sdk.Dec
	MinimumTier1Ratio sdk.Dec
	MinimumTotalRatio sdk.Dec
}

type RiskWeightConfig struct {
	SovereignWeights  map[string]sdk.Dec
	BankWeights       map[string]sdk.Dec
	CorporateWeights  map[string]sdk.Dec
	RetailWeights     map[string]sdk.Dec
	MortgageWeights   map[string]sdk.Dec
}

type CapitalBuffers struct {
	ConservationBuffer    sdk.Dec
	CounterCyclicalMax    sdk.Dec
	SystemicRiskBuffer    sdk.Dec
	DSIBBufferBuckets     []sdk.Dec
}

type StressTestScenario struct {
	ID             string
	Name           string
	Description    string
	Severity       string
	CapitalImpacts CapitalStressImpacts
	RWAImpacts     RWAStressImpacts
	MinimumCET1    sdk.Dec
	MinimumTier1   sdk.Dec
	MinimumTotal   sdk.Dec
}

// Additional supporting types

type CapitalAdjustment struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Amount      sdk.Coin `json:"amount"`
	IsDeduction bool     `json:"is_deduction"`
}

type CapitalDeduction struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Amount      sdk.Coin `json:"amount"`
	TierApplied string   `json:"tier_applied"`
}

type CollateralDetails struct {
	Type      string   `json:"type"`
	Value     sdk.Coin `json:"value"`
	Haircut   sdk.Dec  `json:"haircut"`
	NetValue  sdk.Coin `json:"net_value"`
}

type DerivativeExposure struct {
	ContractType string   `json:"contract_type"`
	Notional     sdk.Coin `json:"notional"`
	MTM          sdk.Coin `json:"mtm"`
	AddOn        sdk.Coin `json:"add_on"`
	NetExposure  sdk.Coin `json:"net_exposure"`
}

type SFTExposure struct {
	TransactionType string   `json:"transaction_type"`
	Collateral      sdk.Coin `json:"collateral"`
	Exposure        sdk.Coin `json:"exposure"`
}

type CCPExposure struct {
	CCPName         string   `json:"ccp_name"`
	ExposureAmount  sdk.Coin `json:"exposure_amount"`
	DefaultFund     sdk.Coin `json:"default_fund"`
}

type IRBApproachDetails struct {
	ApproachType string              `json:"approach_type"`
	PD          map[string]sdk.Dec  `json:"probability_of_default"`
	LGD         map[string]sdk.Dec  `json:"loss_given_default"`
	EAD         map[string]sdk.Coin `json:"exposure_at_default"`
}

type VaRModelDetails struct {
	ModelType       string  `json:"model_type"`
	ConfidenceLevel sdk.Dec `json:"confidence_level"`
	HoldingPeriod   int     `json:"holding_period"`
	LookbackPeriod  int     `json:"lookback_period"`
}

type YearlyGrossIncome struct {
	Year   int      `json:"year"`
	Income sdk.Coin `json:"income"`
}

type BusinessLineRWA struct {
	BusinessLine string   `json:"business_line"`
	Beta         sdk.Dec  `json:"beta"`
	GrossIncome  sdk.Coin `json:"gross_income"`
	RWA          sdk.Coin `json:"rwa"`
}

type CVARiskRWA struct {
	StandardizedCVA sdk.Coin `json:"standardized_cva"`
	AdvancedCVA     sdk.Coin `json:"advanced_cva,omitempty"`
	TotalCVARWA     sdk.Coin `json:"total_cva_rwa"`
}

type BufferRequirement struct {
	BufferType    string  `json:"buffer_type"`
	RequiredRatio sdk.Dec `json:"required_ratio"`
	ActualRatio   sdk.Dec `json:"actual_ratio"`
	IsMet         bool    `json:"is_met"`
}

type BufferRestriction struct {
	RestrictionType     string  `json:"restriction_type"`
	MaxDistribution     sdk.Dec `json:"max_distribution"`
	EffectiveDate       time.Time `json:"effective_date"`
}

type LeverageRatioCalculation struct {
	Tier1Capital       sdk.Coin       `json:"tier1_capital"`
	ExposureMeasure    ExposureMeasure `json:"exposure_measure"`
	LeverageRatio      sdk.Dec         `json:"leverage_ratio"`
	MinimumRequirement sdk.Dec         `json:"minimum_requirement"`
	Buffer             sdk.Dec         `json:"buffer"`
	IsCompliant        bool            `json:"is_compliant"`
}

type ExposureMeasure struct {
	OnBalanceSheet     sdk.Coin `json:"on_balance_sheet"`
	Derivatives        sdk.Coin `json:"derivatives"`
	SFT                sdk.Coin `json:"sft"`
	OffBalanceSheet    sdk.Coin `json:"off_balance_sheet"`
	OtherExposures     sdk.Coin `json:"other_exposures"`
	TotalExposure      sdk.Coin `json:"total_exposure"`
}

type LiquidityCoverageRatio struct {
	HQLA               HQLAComposition `json:"hqla"`
	TotalNetOutflows   sdk.Coin        `json:"total_net_outflows"`
	LCR                sdk.Dec         `json:"lcr"`
	MinimumRequirement sdk.Dec         `json:"minimum_requirement"`
	IsCompliant        bool            `json:"is_compliant"`
	ObservationPeriod  time.Duration   `json:"observation_period"`
}

type HQLAComposition struct {
	Level1     sdk.Coin `json:"level1"`
	Level2A    sdk.Coin `json:"level2a"`
	Level2B    sdk.Coin `json:"level2b"`
	Total      sdk.Coin `json:"total"`
}

type NetStableFundingRatio struct {
	AvailableStableFunding sdk.Coin `json:"available_stable_funding"`
	RequiredStableFunding  sdk.Coin `json:"required_stable_funding"`
	NSFR                   sdk.Dec  `json:"nsfr"`
	MinimumRequirement     sdk.Dec  `json:"minimum_requirement"`
	IsCompliant            bool     `json:"is_compliant"`
}

type LiquidityBuffer struct {
	BufferType  string   `json:"buffer_type"`
	Amount      sdk.Coin `json:"amount"`
	Quality     string   `json:"quality"`
	Availability string  `json:"availability"`
}

type ContingencyFundingPlan struct {
	PlanID          string                `json:"plan_id"`
	TriggerLevels   []LiquidityTrigger    `json:"trigger_levels"`
	FundingSources  []ContingencySource   `json:"funding_sources"`
	ActionPlan      []ContingencyAction   `json:"action_plan"`
	LastUpdated     time.Time             `json:"last_updated"`
	NextReview      time.Time             `json:"next_review"`
}

type LiquidityTrigger struct {
	Level       string  `json:"level"`
	Metric      string  `json:"metric"`
	Threshold   sdk.Dec `json:"threshold"`
	Action      string  `json:"action"`
}

type ContingencySource struct {
	SourceType   string   `json:"source_type"`
	Amount       sdk.Coin `json:"amount"`
	Availability string   `json:"availability"`
	Cost         sdk.Dec  `json:"cost"`
}

type ContingencyAction struct {
	ActionType   string        `json:"action_type"`
	Description  string        `json:"description"`
	Timeline     time.Duration `json:"timeline"`
	Responsible  string        `json:"responsible"`
}

type ComplianceStatus struct {
	IsCompliant      bool              `json:"is_compliant"`
	ComplianceDate   time.Time         `json:"compliance_date"`
	ComplianceItems  []ComplianceItem  `json:"compliance_items"`
	RequiredActions  []RequiredAction  `json:"required_actions"`
	NextReviewDate   time.Time         `json:"next_review_date"`
}

type ComplianceItem struct {
	Requirement string `json:"requirement"`
	Required    string `json:"required"`
	Actual      string `json:"actual"`
	Status      string `json:"status"`
}

type RequiredAction struct {
	ActionType   string    `json:"action_type"`
	Description  string    `json:"description"`
	Amount       sdk.Coin  `json:"amount,omitempty"`
	Deadline     time.Time `json:"deadline"`
	Responsible  string    `json:"responsible"`
	Status       string    `json:"status"`
}

type StressTestResult struct {
	ScenarioID       string          `json:"scenario_id"`
	ScenarioName     string          `json:"scenario_name"`
	RunDate          time.Time       `json:"run_date"`
	StressedCapital  CapitalComponents `json:"stressed_capital"`
	StressedRWA      RiskWeightedAssets `json:"stressed_rwa"`
	StressedRatios   CapitalRatios   `json:"stressed_ratios"`
	PassesTest       bool            `json:"passes_test"`
	CapitalShortfall sdk.Coin        `json:"capital_shortfall,omitempty"`
}

type CapitalPlan struct {
	PlanningHorizon   time.Duration          `json:"planning_horizon"`
	TargetRatios      TargetRatios           `json:"target_ratios"`
	PlannedActions    []PlannedCapitalAction `json:"planned_actions"`
	DividendPolicy    string                 `json:"dividend_policy"`
	MaxDividendPayout sdk.Dec                `json:"max_dividend_payout"`
	ReviewSchedule    []time.Time            `json:"review_schedule"`
}

type TargetRatios struct {
	CET1  sdk.Dec `json:"cet1"`
	Tier1 sdk.Dec `json:"tier1"`
	Total sdk.Dec `json:"total"`
}

type PlannedCapitalAction struct {
	ActionType  string    `json:"action_type"`
	Amount      sdk.Coin  `json:"amount"`
	Timeline    time.Time `json:"timeline"`
	Impact      string    `json:"impact"`
	Status      string    `json:"status"`
}

type RegulatoryAdjustment struct {
	AdjustmentType string   `json:"adjustment_type"`
	Description    string   `json:"description"`
	Amount         sdk.Coin `json:"amount"`
	TierAffected   string   `json:"tier_affected"`
}

type CapitalStressImpacts struct {
	RetainedEarningsImpact sdk.Dec `json:"retained_earnings_impact"`
	OCIImpact              sdk.Dec `json:"oci_impact"`
	ProvisioningImpact     sdk.Dec `json:"provisioning_impact"`
}

type RWAStressImpacts struct {
	CreditRWAMultiplier     sdk.Dec `json:"credit_rwa_multiplier"`
	MarketRWAMultiplier     sdk.Dec `json:"market_rwa_multiplier"`
	OperationalRWAMultiplier sdk.Dec `json:"operational_rwa_multiplier"`
}

// Initialize functions

func initializeCapitalRequirements() CapitalRequirements {
	return CapitalRequirements{
		MinimumCET1Ratio:  sdk.NewDecWithPrec(45, 3),  // 4.5%
		MinimumTier1Ratio: sdk.NewDecWithPrec(60, 3),  // 6.0%
		MinimumTotalRatio: sdk.NewDecWithPrec(80, 3),  // 8.0%
	}
}

func initializeRiskWeights() RiskWeightConfig {
	return RiskWeightConfig{
		SovereignWeights: map[string]sdk.Dec{
			"AAA": sdk.ZeroDec(),
			"AA":  sdk.ZeroDec(),
			"A":   sdk.NewDecWithPrec(20, 2),
			"BBB": sdk.NewDecWithPrec(50, 2),
			"BB":  sdk.NewDecWithPrec(100, 2),
			"B":   sdk.NewDecWithPrec(100, 2),
			"CCC": sdk.NewDecWithPrec(150, 2),
		},
		BankWeights: map[string]sdk.Dec{
			"AAA": sdk.NewDecWithPrec(20, 2),
			"AA":  sdk.NewDecWithPrec(20, 2),
			"A":   sdk.NewDecWithPrec(50, 2),
			"BBB": sdk.NewDecWithPrec(50, 2),
			"BB":  sdk.NewDecWithPrec(100, 2),
			"B":   sdk.NewDecWithPrec(100, 2),
			"CCC": sdk.NewDecWithPrec(150, 2),
		},
		CorporateWeights: map[string]sdk.Dec{
			"AAA": sdk.NewDecWithPrec(20, 2),
			"AA":  sdk.NewDecWithPrec(50, 2),
			"A":   sdk.NewDecWithPrec(100, 2),
			"BBB": sdk.NewDecWithPrec(100, 2),
			"BB":  sdk.NewDecWithPrec(150, 2),
			"B":   sdk.NewDecWithPrec(150, 2),
			"CCC": sdk.NewDecWithPrec(150, 2),
		},
	}
}

func initializeCapitalBuffers() CapitalBuffers {
	return CapitalBuffers{
		ConservationBuffer: sdk.NewDecWithPrec(25, 3),  // 2.5%
		CounterCyclicalMax: sdk.NewDecWithPrec(25, 3),  // 2.5%
		SystemicRiskBuffer: sdk.NewDecWithPrec(10, 3),  // 1.0%
		DSIBBufferBuckets: []sdk.Dec{
			sdk.NewDecWithPrec(10, 3),  // 1.0%
			sdk.NewDecWithPrec(15, 3),  // 1.5%
			sdk.NewDecWithPrec(20, 3),  // 2.0%
			sdk.NewDecWithPrec(25, 3),  // 2.5%
		},
	}
}

func initializeStressTestScenarios() []StressTestScenario {
	return []StressTestScenario{
		{
			ID:          "BASELINE_ADVERSE",
			Name:        "Baseline Adverse Scenario",
			Description: "Moderate economic downturn",
			Severity:    "MODERATE",
			CapitalImpacts: CapitalStressImpacts{
				RetainedEarningsImpact: sdk.NewDecWithPrec(-15, 2), // -15%
				OCIImpact:              sdk.NewDecWithPrec(-10, 2), // -10%
				ProvisioningImpact:     sdk.NewDecWithPrec(20, 2),  // +20%
			},
			RWAImpacts: RWAStressImpacts{
				CreditRWAMultiplier:      sdk.NewDecWithPrec(115, 2), // 1.15x
				MarketRWAMultiplier:      sdk.NewDecWithPrec(120, 2), // 1.20x
				OperationalRWAMultiplier: sdk.NewDecWithPrec(110, 2), // 1.10x
			},
			MinimumCET1:  sdk.NewDecWithPrec(40, 3), // 4.0%
			MinimumTier1: sdk.NewDecWithPrec(55, 3), // 5.5%
			MinimumTotal: sdk.NewDecWithPrec(75, 3), // 7.5%
		},
		{
			ID:          "SEVERELY_ADVERSE",
			Name:        "Severely Adverse Scenario",
			Description: "Severe economic crisis",
			Severity:    "SEVERE",
			CapitalImpacts: CapitalStressImpacts{
				RetainedEarningsImpact: sdk.NewDecWithPrec(-30, 2), // -30%
				OCIImpact:              sdk.NewDecWithPrec(-20, 2), // -20%
				ProvisioningImpact:     sdk.NewDecWithPrec(50, 2),  // +50%
			},
			RWAImpacts: RWAStressImpacts{
				CreditRWAMultiplier:      sdk.NewDecWithPrec(130, 2), // 1.30x
				MarketRWAMultiplier:      sdk.NewDecWithPrec(150, 2), // 1.50x
				OperationalRWAMultiplier: sdk.NewDecWithPrec(125, 2), // 1.25x
			},
			MinimumCET1:  sdk.NewDecWithPrec(35, 3), // 3.5%
			MinimumTier1: sdk.NewDecWithPrec(45, 3), // 4.5%
			MinimumTotal: sdk.NewDecWithPrec(60, 3), // 6.0%
		},
	}
}

// Helper method stubs for compilation

func (bcm *BaselIIICapitalManager) getBankCapitalData(ctx context.Context, bankID string) (BankCapitalData, error) {
	// Stub implementation
	return BankCapitalData{
		BankID:           bankID,
		PaidUpCapital:    sdk.NewCoin("usd", sdk.NewInt(1000000000)),
		RetainedEarnings: sdk.NewCoin("usd", sdk.NewInt(500000000)),
		TotalAssets:      sdk.NewCoin("usd", sdk.NewInt(10000000000)),
	}, nil
}

func (bcm *BaselIIICapitalManager) calculateMinorityInterestCET1(data BankCapitalData) sdk.Coin {
	return sdk.NewCoin("usd", sdk.NewInt(50000000))
}

func (bcm *BaselIIICapitalManager) getCET1RegulatoryAdjustments(data BankCapitalData) []CapitalAdjustment {
	return []CapitalAdjustment{}
}

func (bcm *BaselIIICapitalManager) calculateAdditionalTier1(data BankCapitalData) AdditionalTier1Capital {
	return AdditionalTier1Capital{
		Total: sdk.NewCoin("usd", sdk.NewInt(200000000)),
	}
}

func (bcm *BaselIIICapitalManager) calculateTier2Capital(data BankCapitalData) Tier2Capital {
	return Tier2Capital{
		SubordinatedDebt: sdk.NewCoin("usd", sdk.NewInt(300000000)),
		Total:           sdk.NewCoin("usd", sdk.NewInt(400000000)),
	}
}

func (bcm *BaselIIICapitalManager) calculateRegulatoryDeductions(data BankCapitalData) []CapitalDeduction {
	return []CapitalDeduction{}
}

func (bcm *BaselIIICapitalManager) getBankCreditExposures(ctx context.Context, bankID string) (CreditExposures, error) {
	return CreditExposures{}, nil
}

func (bcm *BaselIIICapitalManager) getRiskWeight(exposure CreditExposure) sdk.Dec {
	// Simplified risk weight assignment
	switch exposure.AssetClass {
	case ASSET_CLASS_SOVEREIGN:
		return sdk.ZeroDec()
	case ASSET_CLASS_BANK:
		return sdk.NewDecWithPrec(20, 2)
	case ASSET_CLASS_RETAIL_MORTGAGE:
		return sdk.NewDecWithPrec(35, 2)
	case ASSET_CLASS_CORPORATE:
		return sdk.NewDecWithPrec(100, 2)
	default:
		return sdk.NewDecWithPrec(100, 2)
	}
}

func (bcm *BaselIIICapitalManager) getCreditConversionFactor(exposure CreditExposure) sdk.Dec {
	// Simplified CCF
	return sdk.NewDecWithPrec(50, 2) // 50%
}

func (bcm *BaselIIICapitalManager) getBaseCurrency() string {
	return "usd"
}

func (bcm *BaselIIICapitalManager) calculateMarketRiskRWA(ctx context.Context, bankID string) (*MarketRiskRWA, error) {
	return &MarketRiskRWA{
		TotalMarketRWA: sdk.NewCoin("usd", sdk.NewInt(500000000)),
	}, nil
}

func (bcm *BaselIIICapitalManager) calculateOperationalRiskRWA(ctx context.Context, bankID string) (*OperationalRiskRWA, error) {
	return &OperationalRiskRWA{
		TotalOperationalRWA: sdk.NewCoin("usd", sdk.NewInt(800000000)),
		BasicIndicatorApproach: true,
	}, nil
}

func (bcm *BaselIIICapitalManager) calculateCVARiskRWA(ctx context.Context, bankID string) (*CVARiskRWA, error) {
	return &CVARiskRWA{
		TotalCVARWA: sdk.NewCoin("usd", sdk.NewInt(100000000)),
	}, nil
}

func (bcm *BaselIIICapitalManager) getTotalAssets(ctx context.Context, bankID string) (sdk.Coin, error) {
	return sdk.NewCoin("usd", sdk.NewInt(10000000000)), nil
}

func (bcm *BaselIIICapitalManager) calculateCounterCyclicalBuffer(ctx context.Context, bankID string) BufferRequirement {
	return BufferRequirement{
		BufferType:    "COUNTER_CYCLICAL",
		RequiredRatio: sdk.NewDecWithPrec(5, 3), // 0.5%
		ActualRatio:   sdk.NewDecWithPrec(5, 3),
		IsMet:         true,
	}
}

func (bcm *BaselIIICapitalManager) isSystemicallyImportant(ctx context.Context, bankID string) bool {
	return false
}

func (bcm *BaselIIICapitalManager) isDSIB(ctx context.Context, bankID string) bool {
	return true
}

func (bcm *BaselIIICapitalManager) calculateDSIBScore(ctx context.Context, bankID string) int {
	return 150 // Basis points
}

func (bcm *BaselIIICapitalManager) getDSIBBufferRequirement(score int) sdk.Dec {
	if score >= 200 {
		return sdk.NewDecWithPrec(25, 3) // 2.5%
	} else if score >= 150 {
		return sdk.NewDecWithPrec(20, 3) // 2.0%
	} else if score >= 100 {
		return sdk.NewDecWithPrec(15, 3) // 1.5%
	} else if score >= 50 {
		return sdk.NewDecWithPrec(10, 3) // 1.0%
	}
	return sdk.ZeroDec()
}

func (bcm *BaselIIICapitalManager) getBufferBreachRestrictions(available, required sdk.Dec) []BufferRestriction {
	restrictions := []BufferRestriction{}
	
	shortfall := required.Sub(available)
	if shortfall.IsPositive() {
		// Calculate distribution restrictions based on shortfall
		var maxPayout sdk.Dec
		if shortfall.LTE(sdk.NewDecWithPrec(625, 4)) { // 0.625%
			maxPayout = sdk.NewDecWithPrec(60, 2) // 60%
		} else if shortfall.LTE(sdk.NewDecWithPrec(125, 3)) { // 1.25%
			maxPayout = sdk.NewDecWithPrec(40, 2) // 40%
		} else if shortfall.LTE(sdk.NewDecWithPrec(1875, 4)) { // 1.875%
			maxPayout = sdk.NewDecWithPrec(20, 2) // 20%
		} else {
			maxPayout = sdk.ZeroDec() // 0%
		}

		restrictions = append(restrictions, BufferRestriction{
			RestrictionType: "DIVIDEND_RESTRICTION",
			MaxDistribution: maxPayout,
			EffectiveDate:   time.Now(),
		})
	}

	return restrictions
}

func (bcm *BaselIIICapitalManager) calculateLeverageExposure(ctx context.Context, bankID string) (*ExposureMeasure, error) {
	return &ExposureMeasure{
		OnBalanceSheet:  sdk.NewCoin("usd", sdk.NewInt(8000000000)),
		Derivatives:     sdk.NewCoin("usd", sdk.NewInt(500000000)),
		SFT:            sdk.NewCoin("usd", sdk.NewInt(300000000)),
		OffBalanceSheet: sdk.NewCoin("usd", sdk.NewInt(1200000000)),
		TotalExposure:   sdk.NewCoin("usd", sdk.NewInt(10000000000)),
	}, nil
}

func (bcm *BaselIIICapitalManager) getHQLA(ctx context.Context, bankID string) (*HQLAComposition, error) {
	return &HQLAComposition{
		Level1:  sdk.NewCoin("usd", sdk.NewInt(1000000000)),
		Level2A: sdk.NewCoin("usd", sdk.NewInt(200000000)),
		Level2B: sdk.NewCoin("usd", sdk.NewInt(100000000)),
		Total:   sdk.NewCoin("usd", sdk.NewInt(1300000000)),
	}, nil
}

func (bcm *BaselIIICapitalManager) calculateNetCashOutflows(ctx context.Context, bankID string) (sdk.Coin, error) {
	return sdk.NewCoin("usd", sdk.NewInt(1200000000)), nil
}

func (bcm *BaselIIICapitalManager) calculateNSFR(ctx context.Context, bankID string) (*NetStableFundingRatio, error) {
	return &NetStableFundingRatio{
		AvailableStableFunding: sdk.NewCoin("usd", sdk.NewInt(5000000000)),
		RequiredStableFunding:  sdk.NewCoin("usd", sdk.NewInt(4500000000)),
		NSFR:                   sdk.NewDecWithPrec(111, 2), // 111%
		MinimumRequirement:     sdk.NewDecWithPrec(100, 2),
		IsCompliant:           true,
	}, nil
}

func (bcm *BaselIIICapitalManager) getLiquidityBuffers(ctx context.Context, bankID string) ([]LiquidityBuffer, error) {
	return []LiquidityBuffer{
		{
			BufferType:   "CENTRAL_BANK_RESERVES",
			Amount:       sdk.NewCoin("usd", sdk.NewInt(500000000)),
			Quality:      "LEVEL_1",
			Availability: "IMMEDIATE",
		},
	}, nil
}

func (bcm *BaselIIICapitalManager) calculateStressedOutflows(ctx context.Context, bankID string) sdk.Coin {
	return sdk.NewCoin("usd", sdk.NewInt(2000000000))
}

func (bcm *BaselIIICapitalManager) getContingencyFundingPlan(ctx context.Context, bankID string) ContingencyFundingPlan {
	return ContingencyFundingPlan{
		PlanID:      "CFP_2025",
		LastUpdated: time.Now(),
		NextReview:  time.Now().AddDate(0, 6, 0),
	}
}

func (bcm *BaselIIICapitalManager) applyCapitalStress(capital *CapitalComponents, impacts CapitalStressImpacts) CapitalComponents {
	stressedCapital := *capital
	
	// Apply stress to retained earnings
	stressedCapital.CommonEquityTier1.RetainedEarnings = sdk.NewCoin(
		capital.CommonEquityTier1.RetainedEarnings.Denom,
		capital.CommonEquityTier1.RetainedEarnings.Amount.Mul(
			sdk.NewInt(100).Add(sdk.NewInt(int64(impacts.RetainedEarningsImpact.MustFloat64()*100))),
		).Quo(sdk.NewInt(100)),
	)
	
	return stressedCapital
}

func (bcm *BaselIIICapitalManager) applyRWAStress(rwa *RiskWeightedAssets, impacts RWAStressImpacts) RiskWeightedAssets {
	stressedRWA := *rwa
	
	// Apply multipliers
	stressedRWA.CreditRiskRWA.TotalCreditRWA = sdk.NewCoin(
		rwa.CreditRiskRWA.TotalCreditRWA.Denom,
		rwa.CreditRiskRWA.TotalCreditRWA.Amount.Mul(sdk.NewInt(int64(impacts.CreditRWAMultiplier.MustFloat64()*100))).Quo(sdk.NewInt(100)),
	)
	
	return stressedRWA
}

func (bcm *BaselIIICapitalManager) calculateCapitalShortfall(capital CapitalComponents, rwa RiskWeightedAssets, scenario StressTestScenario) sdk.Coin {
	requiredCapital := rwa.TotalRWA.Amount.Mul(sdk.NewInt(int64(scenario.MinimumCET1.MustFloat64()*10000))).Quo(sdk.NewInt(10000))
	
	if capital.CommonEquityTier1.Total.Amount.LT(requiredCapital) {
		return sdk.NewCoin(capital.CommonEquityTier1.Total.Denom, requiredCapital.Sub(capital.CommonEquityTier1.Total.Amount))
	}
	
	return sdk.NewCoin(bcm.getBaseCurrency(), sdk.ZeroInt())
}

func (bcm *BaselIIICapitalManager) storeCapitalAdequacyReport(ctx context.Context, report CapitalAdequacyReport) error {
	store := bcm.keeper.storeService.OpenKVStore(ctx)
	key := []byte(fmt.Sprintf("basel_report_%s", report.ReportID))
	bz, err := json.Marshal(report)
	if err != nil {
		return err
	}
	store.Set(key, bz)
	return nil
}

func (bcm *BaselIIICapitalManager) emitCapitalAdequacyEvent(ctx context.Context, report *CapitalAdequacyReport) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"basel_iii_capital_adequacy_calculated",
			sdk.NewAttribute("report_id", report.ReportID),
			sdk.NewAttribute("bank_id", report.BankID),
			sdk.NewAttribute("cet1_ratio", report.CapitalRatios.CET1Ratio.String()),
			sdk.NewAttribute("tier1_ratio", report.CapitalRatios.Tier1Ratio.String()),
			sdk.NewAttribute("total_ratio", report.CapitalRatios.TotalCapitalRatio.String()),
			sdk.NewAttribute("is_compliant", fmt.Sprintf("%v", report.ComplianceStatus.IsCompliant)),
		),
	)
}

// Public API methods

func (bcm *BaselIIICapitalManager) GetCapitalAdequacyReport(ctx context.Context, reportID string) (*CapitalAdequacyReport, error) {
	store := bcm.keeper.storeService.OpenKVStore(ctx)
	key := []byte(fmt.Sprintf("basel_report_%s", reportID))
	bz := store.Get(key)
	if bz == nil {
		return nil, fmt.Errorf("report not found: %s", reportID)
	}
	
	var report CapitalAdequacyReport
	if err := json.Unmarshal(bz, &report); err != nil {
		return nil, fmt.Errorf("failed to unmarshal report: %w", err)
	}
	
	return &report, nil
}

func (bcm *BaselIIICapitalManager) GetBankReports(ctx context.Context, bankID string) ([]CapitalAdequacyReport, error) {
	// Implementation would query all reports for a bank
	return []CapitalAdequacyReport{}, nil
}