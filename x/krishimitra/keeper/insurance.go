package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/namo/x/krishimitra/types"
)

// InsuranceEngine handles crop insurance and weather derivatives
type InsuranceEngine struct {
	keeper Keeper
}

// NewInsuranceEngine creates a new insurance engine
func NewInsuranceEngine(keeper Keeper) *InsuranceEngine {
	return &InsuranceEngine{
		keeper: keeper,
	}
}

// CropInsurancePolicy represents a comprehensive crop insurance policy
type CropInsurancePolicy struct {
	PolicyID          string                    `json:"policy_id"`
	FarmerID          string                    `json:"farmer_id"`
	LoanID            string                    `json:"loan_id,omitempty"`
	CropType          string                    `json:"crop_type"`
	CropVariety       string                    `json:"crop_variety"`
	SeasonType        string                    `json:"season_type"` // KHARIF, RABI, ZAID
	SowingArea        sdk.Dec                   `json:"sowing_area"`
	SumInsured        sdk.Coin                  `json:"sum_insured"`
	PremiumAmount     sdk.Coin                  `json:"premium_amount"`
	CoverageType      string                    `json:"coverage_type"` // COMPREHENSIVE, WEATHER_BASED, YIELD_BASED
	WeatherTriggers   []types.WeatherTrigger    `json:"weather_triggers"`
	YieldTriggers     []types.YieldTrigger      `json:"yield_triggers"`
	PolicyPeriod      types.PolicyPeriod        `json:"policy_period"`
	Status            types.PolicyStatus        `json:"status"`
	ClaimsHistory     []types.InsuranceClaim    `json:"claims_history"`
	PremiumsPaid      []types.PremiumPayment    `json:"premiums_paid"`
	CreatedAt         time.Time                 `json:"created_at"`
	UpdatedAt         time.Time                 `json:"updated_at"`
}

// WeatherDerivative represents weather-based financial instruments
type WeatherDerivative struct {
	DerivativeID      string                     `json:"derivative_id"`
	HolderID          string                     `json:"holder_id"`
	WeatherIndex      string                     `json:"weather_index"` // RAINFALL, TEMPERATURE, HUMIDITY
	ReferenceStation  string                     `json:"reference_station"`
	StrikeValue       sdk.Dec                    `json:"strike_value"`
	PayoutStructure   types.PayoutStructure      `json:"payout_structure"`
	Premium           sdk.Coin                   `json:"premium"`
	MaxPayout         sdk.Coin                   `json:"max_payout"`
	ContractPeriod    types.DerivativePeriod     `json:"contract_period"`
	Status            types.DerivativeStatus     `json:"status"`
	WeatherData       []types.WeatherDataPoint   `json:"weather_data"`
	PayoutHistory     []types.DerivativePayout   `json:"payout_history"`
	CreatedAt         time.Time                  `json:"created_at"`
}

// CreateCropInsurancePolicy creates a new crop insurance policy
func (ie *InsuranceEngine) CreateCropInsurancePolicy(ctx sdk.Context, request *types.InsurancePolicyRequest) (*CropInsurancePolicy, error) {
	// Validate farmer eligibility
	farmerProfile, found := ie.keeper.GetFarmerProfile(ctx, request.FarmerID)
	if !found {
		return nil, fmt.Errorf("farmer profile not found: %s", request.FarmerID)
	}

	// Validate crop and area
	if request.SowingArea.IsZero() || request.SowingArea.GT(farmerProfile.TotalLandArea) {
		return nil, fmt.Errorf("invalid sowing area: %s", request.SowingArea.String())
	}

	// Generate policy ID
	policyID := ie.generatePolicyID(ctx, request.FarmerID, request.CropType)

	// Calculate insurance parameters
	sumInsured := ie.calculateSumInsured(ctx, request)
	premiumAmount := ie.calculatePremiumAmount(ctx, request, sumInsured)
	
	// Create weather triggers
	weatherTriggers := ie.createWeatherTriggers(ctx, request)
	
	// Create yield triggers
	yieldTriggers := ie.createYieldTriggers(ctx, request)

	policy := &CropInsurancePolicy{
		PolicyID:        policyID,
		FarmerID:        request.FarmerID,
		LoanID:          request.LoanID,
		CropType:        request.CropType,
		CropVariety:     request.CropVariety,
		SeasonType:      request.SeasonType,
		SowingArea:      request.SowingArea,
		SumInsured:      sumInsured,
		PremiumAmount:   premiumAmount,
		CoverageType:    request.CoverageType,
		WeatherTriggers: weatherTriggers,
		YieldTriggers:   yieldTriggers,
		PolicyPeriod: types.PolicyPeriod{
			StartDate: request.SowingDate,
			EndDate:   request.HarvestDate,
		},
		Status:        types.PolicyStatusActive,
		ClaimsHistory: []types.InsuranceClaim{},
		PremiumsPaid:  []types.PremiumPayment{},
		CreatedAt:     ctx.BlockTime(),
		UpdatedAt:     ctx.BlockTime(),
	}

	// Store policy
	ie.keeper.SetInsurancePolicy(ctx, *policy)

	// Emit policy creation event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeInsurancePolicyCreated,
			sdk.NewAttribute(types.AttributeKeyPolicyID, policyID),
			sdk.NewAttribute(types.AttributeKeyFarmerID, request.FarmerID),
			sdk.NewAttribute(types.AttributeKeyCropType, request.CropType),
			sdk.NewAttribute(types.AttributeKeySumInsured, sumInsured.String()),
			sdk.NewAttribute(types.AttributeKeyPremium, premiumAmount.String()),
		),
	)

	return policy, nil
}

// CreateWeatherDerivative creates a weather-based derivative contract
func (ie *InsuranceEngine) CreateWeatherDerivative(ctx sdk.Context, request *types.WeatherDerivativeRequest) (*WeatherDerivative, error) {
	// Generate derivative ID
	derivativeID := ie.generateDerivativeID(ctx, request.HolderID, request.WeatherIndex)

	// Calculate premium and payout structure
	premium := ie.calculateDerivativePremium(ctx, request)
	payoutStructure := ie.createPayoutStructure(ctx, request)

	derivative := &WeatherDerivative{
		DerivativeID:     derivativeID,
		HolderID:         request.HolderID,
		WeatherIndex:     request.WeatherIndex,
		ReferenceStation: request.ReferenceStation,
		StrikeValue:      request.StrikeValue,
		PayoutStructure:  payoutStructure,
		Premium:          premium,
		MaxPayout:        request.MaxPayout,
		ContractPeriod: types.DerivativePeriod{
			StartDate: request.StartDate,
			EndDate:   request.EndDate,
		},
		Status:        types.DerivativeStatusActive,
		WeatherData:   []types.WeatherDataPoint{},
		PayoutHistory: []types.DerivativePayout{},
		CreatedAt:     ctx.BlockTime(),
	}

	// Store derivative
	ie.keeper.SetWeatherDerivative(ctx, *derivative)

	// Emit derivative creation event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeWeatherDerivativeCreated,
			sdk.NewAttribute(types.AttributeKeyDerivativeID, derivativeID),
			sdk.NewAttribute(types.AttributeKeyHolderID, request.HolderID),
			sdk.NewAttribute(types.AttributeKeyWeatherIndex, request.WeatherIndex),
			sdk.NewAttribute(types.AttributeKeyPremium, premium.String()),
		),
	)

	return derivative, nil
}

// ProcessWeatherTriggers checks weather conditions and triggers payouts
func (ie *InsuranceEngine) ProcessWeatherTriggers(ctx sdk.Context) error {
	// Get all active insurance policies
	activePolicies := ie.keeper.GetActivePolicies(ctx)

	for _, policy := range activePolicies {
		// Get current weather data for policy location
		weatherData, err := ie.keeper.GetCurrentWeatherData(ctx, policy.FarmerID)
		if err != nil {
			continue // Skip if weather data unavailable
		}

		// Check each weather trigger
		for _, trigger := range policy.WeatherTriggers {
			triggered, payoutAmount := ie.evaluateWeatherTrigger(weatherData, trigger)
			
			if triggered {
				// Create insurance claim
				claim := types.InsuranceClaim{
					ClaimID:       ie.generateClaimID(ctx, policy.PolicyID),
					PolicyID:      policy.PolicyID,
					ClaimType:     "WEATHER_TRIGGER",
					TriggerType:   trigger.TriggerType,
					ClaimAmount:   payoutAmount,
					WeatherData:   weatherData,
					Status:        types.ClaimStatusApproved, // Auto-approved for parametric triggers
					ClaimDate:     ctx.BlockTime(),
					ProcessedDate: ctx.BlockTime(),
				}

				// Process payout
				err := ie.processClaim(ctx, policy.PolicyID, claim)
				if err != nil {
					// Log error but continue processing other policies
					ie.keeper.Logger(ctx).Error("Failed to process weather claim", "policy_id", policy.PolicyID, "error", err)
					continue
				}

				// Emit claim processed event
				ctx.EventManager().EmitEvent(
					sdk.NewEvent(
						types.EventTypeClaimProcessed,
						sdk.NewAttribute(types.AttributeKeyClaimID, claim.ClaimID),
						sdk.NewAttribute(types.AttributeKeyPolicyID, policy.PolicyID),
						sdk.NewAttribute(types.AttributeKeyClaimAmount, payoutAmount.String()),
						sdk.NewAttribute(types.AttributeKeyTriggerType, trigger.TriggerType),
					),
				)
			}
		}
	}

	return nil
}

// ProcessWeatherDerivatives evaluates and settles weather derivatives
func (ie *InsuranceEngine) ProcessWeatherDerivatives(ctx sdk.Context) error {
	activeDerivatives := ie.keeper.GetActiveWeatherDerivatives(ctx)

	for _, derivative := range activeDerivatives {
		// Check if contract period has ended
		if ctx.BlockTime().After(derivative.ContractPeriod.EndDate) {
			// Calculate final settlement
			payout := ie.calculateDerivativePayout(ctx, derivative)
			
			if payout.IsPositive() {
				// Process payout
				derivativePayout := types.DerivativePayout{
					PayoutID:      ie.generatePayoutID(ctx, derivative.DerivativeID),
					DerivativeID:  derivative.DerivativeID,
					PayoutAmount:  payout,
					TriggerValue:  ie.getFinalWeatherValue(ctx, derivative),
					PayoutDate:    ctx.BlockTime(),
				}

				err := ie.processDerivativePayout(ctx, derivative.DerivativeID, derivativePayout)
				if err != nil {
					ie.keeper.Logger(ctx).Error("Failed to process derivative payout", "derivative_id", derivative.DerivativeID, "error", err)
					continue
				}

				// Update derivative status
				derivative.Status = types.DerivativeStatusSettled
				ie.keeper.SetWeatherDerivative(ctx, derivative)

				// Emit settlement event
				ctx.EventManager().EmitEvent(
					sdk.NewEvent(
						types.EventTypeDerivativeSettled,
						sdk.NewAttribute(types.AttributeKeyDerivativeID, derivative.DerivativeID),
						sdk.NewAttribute(types.AttributeKeyPayoutAmount, payout.String()),
					),
				)
			}
		}
	}

	return nil
}

// calculateSumInsured calculates the sum insured based on crop and area
func (ie *InsuranceEngine) calculateSumInsured(ctx sdk.Context, request *types.InsurancePolicyRequest) sdk.Coin {
	params := ie.keeper.GetParams(ctx)
	
	// Get crop-specific insurance rates
	cropInfo, found := ie.keeper.GetCropInsuranceInfo(ctx, request.CropType)
	if !found {
		// Use default rates
		baseValue := request.SowingArea.Mul(params.DefaultInsurancePerAcre)
		return sdk.NewCoin("dinr", baseValue.TruncateInt())
	}

	// Calculate based on expected yield and market price
	expectedYield := request.SowingArea.Mul(cropInfo.AverageYieldPerAcre)
	expectedValue := expectedYield.Mul(cropInfo.AverageMarketPrice)
	
	// Apply coverage percentage (typically 80-90% of expected value)
	coverageRatio := sdk.NewDecWithPrec(85, 2) // 85% coverage
	sumInsured := expectedValue.Mul(coverageRatio).TruncateInt()

	return sdk.NewCoin("dinr", sumInsured)
}

// calculatePremiumAmount calculates the insurance premium
func (ie *InsuranceEngine) calculatePremiumAmount(ctx sdk.Context, request *types.InsurancePolicyRequest, sumInsured sdk.Coin) sdk.Coin {
	params := ie.keeper.GetParams(ctx)
	
	// Base premium rate (percentage of sum insured)
	basePremiumRate := params.BasePremiumRate // e.g., 5%

	// Risk adjustments
	riskMultiplier := ie.calculateRiskMultiplier(ctx, request)
	
	// Calculate premium
	premiumAmount := sumInsured.Amount.ToDec().
		Mul(basePremiumRate).
		Mul(riskMultiplier).
		TruncateInt()

	// Government subsidy (typically 50% of premium for small farmers)
	farmerProfile, _ := ie.keeper.GetFarmerProfile(ctx, request.FarmerID)
	if farmerProfile.TotalLandArea.LT(sdk.NewDec(5)) { // Small farmer (< 5 acres)
		subsidyRate := params.SmallFarmerSubsidy // e.g., 50%
		subsidy := premiumAmount.ToDec().Mul(subsidyRate).TruncateInt()
		premiumAmount = premiumAmount.Sub(subsidy)
	}

	return sdk.NewCoin(sumInsured.Denom, premiumAmount)
}

// calculateRiskMultiplier calculates risk-based premium adjustment
func (ie *InsuranceEngine) calculateRiskMultiplier(ctx sdk.Context, request *types.InsurancePolicyRequest) sdk.Dec {
	multiplier := sdk.OneDec() // Base multiplier

	// Weather risk factor
	weatherRisk, found := ie.keeper.GetWeatherRiskProfile(ctx, request.FarmerID)
	if found {
		if weatherRisk.DroughtFrequency >= 5 { // High drought risk
			multiplier = multiplier.Add(sdk.NewDecWithPrec(5, 1)) // +50%
		} else if weatherRisk.DroughtFrequency >= 3 { // Medium drought risk
			multiplier = multiplier.Add(sdk.NewDecWithPrec(25, 2)) // +25%
		}

		if weatherRisk.FloodRisk == "HIGH" {
			multiplier = multiplier.Add(sdk.NewDecWithPrec(3, 1)) // +30%
		} else if weatherRisk.FloodRisk == "MEDIUM" {
			multiplier = multiplier.Add(sdk.NewDecWithPrec(15, 2)) // +15%
		}
	}

	// Crop-specific risk
	cropRisk := ie.getCropRiskFactor(request.CropType)
	multiplier = multiplier.Add(cropRisk)

	// Farmer's claim history
	claimHistory := ie.keeper.GetFarmerClaimHistory(ctx, request.FarmerID)
	if len(claimHistory) > 0 {
		recentClaims := 0
		for _, claim := range claimHistory {
			if claim.ClaimDate.After(ctx.BlockTime().AddDate(-3, 0, 0)) { // Last 3 years
				recentClaims++
			}
		}
		if recentClaims >= 2 {
			multiplier = multiplier.Add(sdk.NewDecWithPrec(2, 1)) // +20% for frequent claimants
		}
	}

	// Cap the multiplier
	if multiplier.GT(sdk.NewDecWithPrec(25, 1)) { // Max 2.5x
		multiplier = sdk.NewDecWithPrec(25, 1)
	}

	return multiplier
}

// createWeatherTriggers creates parametric weather triggers
func (ie *InsuranceEngine) createWeatherTriggers(ctx sdk.Context, request *types.InsurancePolicyRequest) []types.WeatherTrigger {
	triggers := []types.WeatherTrigger{}

	// Rainfall triggers (most common for Indian agriculture)
	rainfallTrigger := types.WeatherTrigger{
		TriggerType:    "RAINFALL_DEFICIT",
		WeatherMetric:  "CUMULATIVE_RAINFALL",
		ThresholdValue: ie.calculateRainfallThreshold(ctx, request.CropType, request.SeasonType),
		ComparisonType: "LESS_THAN",
		PayoutStructure: types.TriggerPayoutStructure{
			FixedPayout:    false,
			BaseAmount:     sdk.NewCoin("dinr", sdk.NewInt(10000)), // Base payout
			ScalingFactor:  sdk.NewDecWithPrec(2, 0), // 2x multiplier for extreme deficits
		},
		MeasurementPeriod: types.MeasurementPeriod{
			StartOffset: 0,   // From sowing
			Duration:    90,  // Days
		},
	}
	triggers = append(triggers, rainfallTrigger)

	// Temperature triggers (for heat-sensitive crops)
	if ie.isHeatSensitiveCrop(request.CropType) {
		tempTrigger := types.WeatherTrigger{
			TriggerType:    "EXTREME_TEMPERATURE",
			WeatherMetric:  "MAX_TEMPERATURE",
			ThresholdValue: sdk.NewDec(42), // 42°C threshold
			ComparisonType: "GREATER_THAN",
			PayoutStructure: types.TriggerPayoutStructure{
				FixedPayout:   true,
				BaseAmount:    sdk.NewCoin("dinr", sdk.NewInt(5000)),
				ScalingFactor: sdk.OneDec(),
			},
			MeasurementPeriod: types.MeasurementPeriod{
				StartOffset: 30, // After 30 days of sowing
				Duration:    60, // Critical growth period
			},
		}
		triggers = append(triggers, tempTrigger)
	}

	// Drought triggers
	droughtTrigger := types.WeatherTrigger{
		TriggerType:    "PROLONGED_DROUGHT",
		WeatherMetric:  "CONSECUTIVE_DRY_DAYS",
		ThresholdValue: sdk.NewDec(21), // 21 consecutive dry days
		ComparisonType: "GREATER_THAN",
		PayoutStructure: types.TriggerPayoutStructure{
			FixedPayout:   false,
			BaseAmount:    sdk.NewCoin("dinr", sdk.NewInt(15000)),
			ScalingFactor: sdk.NewDecWithPrec(15, 1), // 1.5x for severe drought
		},
		MeasurementPeriod: types.MeasurementPeriod{
			StartOffset: 0,
			Duration:    120,
		},
	}
	triggers = append(triggers, droughtTrigger)

	return triggers
}

// createYieldTriggers creates yield-based triggers
func (ie *InsuranceEngine) createYieldTriggers(ctx sdk.Context, request *types.InsurancePolicyRequest) []types.YieldTrigger {
	triggers := []types.YieldTrigger{}

	// Get historical yield data
	averageYield := ie.keeper.GetAverageYield(ctx, request.CropType, request.FarmerID)
	
	// Yield loss trigger (significant yield reduction)
	yieldTrigger := types.YieldTrigger{
		TriggerType:      "YIELD_LOSS",
		YieldMetric:      "ACTUAL_YIELD_PER_ACRE",
		ThresholdYield:   averageYield.Mul(sdk.NewDecWithPrec(7, 1)), // 70% of average yield
		ComparisonType:   "LESS_THAN",
		PayoutStructure: types.YieldPayoutStructure{
			PayoutType:     "PROPORTIONAL",
			MaxPayoutRate:  sdk.NewDecWithPrec(8, 1), // 80% of sum insured
			BaseThreshold:  sdk.NewDecWithPrec(5, 1), // 50% of average yield
		},
	}
	triggers = append(triggers, yieldTrigger)

	// Total crop failure trigger
	failureTrigger := types.YieldTrigger{
		TriggerType:     "TOTAL_CROP_FAILURE",
		YieldMetric:     "ACTUAL_YIELD_PER_ACRE",
		ThresholdYield:  averageYield.Mul(sdk.NewDecWithPrec(2, 1)), // 20% of average yield
		ComparisonType:  "LESS_THAN",
		PayoutStructure: types.YieldPayoutStructure{
			PayoutType:    "FIXED",
			MaxPayoutRate: sdk.OneDec(), // 100% of sum insured
		},
	}
	triggers = append(triggers, failureTrigger)

	return triggers
}

// Helper functions for specific calculations
func (ie *InsuranceEngine) calculateRainfallThreshold(ctx sdk.Context, cropType, seasonType string) sdk.Dec {
	// Crop and season-specific rainfall requirements (in mm)
	requirements := map[string]map[string]sdk.Dec{
		"RICE": {
			"KHARIF": sdk.NewDec(1000), // 1000mm for kharif rice
			"RABI":   sdk.NewDec(800),  // 800mm for rabi rice
		},
		"WHEAT": {
			"RABI": sdk.NewDec(600), // 600mm for wheat
		},
		"COTTON": {
			"KHARIF": sdk.NewDec(800), // 800mm for cotton
		},
		"SUGARCANE": {
			"ANNUAL": sdk.NewDec(1500), // 1500mm for sugarcane
		},
	}

	if cropReq, found := requirements[cropType]; found {
		if threshold, found := cropReq[seasonType]; found {
			return threshold.Mul(sdk.NewDecWithPrec(8, 1)) // 80% of requirement as threshold
		}
	}

	// Default threshold
	return sdk.NewDec(500) // 500mm default
}

func (ie *InsuranceEngine) isHeatSensitiveCrop(cropType string) bool {
	heatSensitive := []string{"WHEAT", "MUSTARD", "POTATO", "TOMATO"}
	for _, crop := range heatSensitive {
		if crop == cropType {
			return true
		}
	}
	return false
}

func (ie *InsuranceEngine) getCropRiskFactor(cropType string) sdk.Dec {
	// Risk factors based on historical crop performance
	riskFactors := map[string]sdk.Dec{
		"RICE":       sdk.NewDecWithPrec(1, 2),  // 1% (low risk)
		"WHEAT":      sdk.NewDecWithPrec(15, 3), // 1.5% (low-medium risk)
		"COTTON":     sdk.NewDecWithPrec(3, 2),  // 3% (medium risk)
		"SUGARCANE":  sdk.NewDecWithPrec(2, 2),  // 2% (low-medium risk)
		"VEGETABLES": sdk.NewDecWithPrec(5, 2),  // 5% (high risk)
	}

	if factor, found := riskFactors[cropType]; found {
		return factor
	}
	return sdk.NewDecWithPrec(25, 3) // 2.5% default
}

// evaluateWeatherTrigger checks if weather conditions trigger a payout
func (ie *InsuranceEngine) evaluateWeatherTrigger(weatherData types.WeatherData, trigger types.WeatherTrigger) (bool, sdk.Coin) {
	var actualValue sdk.Dec
	
	// Extract relevant weather metric
	switch trigger.WeatherMetric {
	case "CUMULATIVE_RAINFALL":
		actualValue = weatherData.CumulativeRainfall
	case "MAX_TEMPERATURE":
		actualValue = weatherData.MaxTemperature
	case "CONSECUTIVE_DRY_DAYS":
		actualValue = sdk.NewDec(weatherData.ConsecutiveDryDays)
	default:
		return false, sdk.Coin{}
	}

	// Check if trigger condition is met
	triggered := false
	switch trigger.ComparisonType {
	case "LESS_THAN":
		triggered = actualValue.LT(trigger.ThresholdValue)
	case "GREATER_THAN":
		triggered = actualValue.GT(trigger.ThresholdValue)
	case "EQUAL":
		triggered = actualValue.Equal(trigger.ThresholdValue)
	}

	if !triggered {
		return false, sdk.Coin{}
	}

	// Calculate payout amount
	var payoutAmount sdk.Int
	if trigger.PayoutStructure.FixedPayout {
		payoutAmount = trigger.PayoutStructure.BaseAmount.Amount
	} else {
		// Scale payout based on severity
		deviation := trigger.ThresholdValue.Sub(actualValue).Abs()
		scalingMultiplier := sdk.OneDec().Add(deviation.Mul(trigger.PayoutStructure.ScalingFactor))
		payoutAmount = trigger.PayoutStructure.BaseAmount.Amount.ToDec().Mul(scalingMultiplier).TruncateInt()
	}

	return true, sdk.NewCoin(trigger.PayoutStructure.BaseAmount.Denom, payoutAmount)
}

// Additional helper functions would be implemented for:
// - processClaim
// - calculateDerivativePayout
// - processDerivativePayout
// - generatePolicyID, generateClaimID, etc.
// - createPayoutStructure
// - calculateDerivativePremium

// generatePolicyID generates a unique policy ID
func (ie *InsuranceEngine) generatePolicyID(ctx sdk.Context, farmerID, cropType string) string {
	timestamp := ctx.BlockTime().Unix()
	return fmt.Sprintf("POL-%s-%s-%d", farmerID[:8], cropType[:4], timestamp)
}

// generateClaimID generates a unique claim ID
func (ie *InsuranceEngine) generateClaimID(ctx sdk.Context, policyID string) string {
	timestamp := ctx.BlockTime().Unix()
	return fmt.Sprintf("CLM-%s-%d", policyID[4:12], timestamp)
}

// processClaim processes an insurance claim and transfers funds
func (ie *InsuranceEngine) processClaim(ctx sdk.Context, policyID string, claim types.InsuranceClaim) error {
	// Transfer claim amount from insurance pool to farmer
	policy, found := ie.keeper.GetInsurancePolicy(ctx, policyID)
	if !found {
		return fmt.Errorf("policy not found: %s", policyID)
	}

	farmerAddr, err := sdk.AccAddressFromBech32(policy.FarmerID)
	if err != nil {
		return fmt.Errorf("invalid farmer address: %s", policy.FarmerID)
	}

	// Transfer from insurance module to farmer
	err = ie.keeper.bankKeeper.SendCoinsFromModuleToAccount(
		ctx,
		types.InsurancePoolName,
		farmerAddr,
		sdk.NewCoins(claim.ClaimAmount),
	)
	if err != nil {
		return fmt.Errorf("failed to transfer claim amount: %w", err)
	}

	// Update policy with claim
	policy.ClaimsHistory = append(policy.ClaimsHistory, claim)
	ie.keeper.SetInsurancePolicy(ctx, policy)

	return nil
}