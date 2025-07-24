/*
Copyright 2024 DeshChain Foundation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package types

import (
	"math/big"
	"strconv"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TaxCalculator handles all tax calculations for DeshChain
type TaxCalculator struct {
	config           *TaxConfig
	volumeData       *VolumeData
	patriotismScore  int32
	culturalScore    int32
	donationFlag     bool
	optimizationMode bool
}

// NewTaxCalculator creates a new tax calculator instance
func NewTaxCalculator(config *TaxConfig) *TaxCalculator {
	return &TaxCalculator{
		config:           config,
		volumeData:       &VolumeData{},
		patriotismScore:  0,
		culturalScore:    0,
		donationFlag:     false,
		optimizationMode: true,
	}
}

// SetVolumeData sets the volume data for tax calculations
func (tc *TaxCalculator) SetVolumeData(volumeData *VolumeData) {
	tc.volumeData = volumeData
}

// SetPatriotismScore sets the patriotism score for tax calculations
func (tc *TaxCalculator) SetPatriotismScore(score int32) {
	tc.patriotismScore = score
}

// SetCulturalScore sets the cultural score for tax calculations
func (tc *TaxCalculator) SetCulturalScore(score int32) {
	tc.culturalScore = score
}

// SetDonationFlag sets whether this transaction is a donation
func (tc *TaxCalculator) SetDonationFlag(isDonation bool) {
	tc.donationFlag = isDonation
}

// SetOptimizationMode enables or disables tax optimization
func (tc *TaxCalculator) SetOptimizationMode(enabled bool) {
	tc.optimizationMode = enabled
}

// CalculateTax calculates the total tax for a transaction
func (tc *TaxCalculator) CalculateTax(transactionAmount sdk.Coin, messageType string) (*TaxCalculationResult, error) {
	// If it's a donation, no tax is applied
	if tc.donationFlag {
		return &TaxCalculationResult{
			OriginalAmount:    transactionAmount,
			TaxAmount:         sdk.NewCoin(transactionAmount.Denom, math.ZeroInt()),
			EffectiveRate:     "0.0",
			AppliedDiscounts:  []string{"donation_exemption"},
			TaxBreakdown:      tc.createTaxBreakdown(transactionAmount, math.ZeroInt()),
			Optimizations:     []string{"donation_exemption"},
			SavingsAmount:     tc.calculateBaseTax(transactionAmount),
			OptimizationScore: "100",
		}, nil
	}

	// Calculate base tax
	baseTax := tc.calculateBaseTax(transactionAmount)
	
	// Apply volume-based discounts
	volumeDiscountRate := tc.calculateVolumeDiscount()
	volumeDiscountAmount := tc.applyDiscount(baseTax, volumeDiscountRate)
	
	// Apply patriotism-based discounts
	patriotismDiscountRate := tc.calculatePatriotismDiscount()
	patriotismDiscountAmount := tc.applyDiscount(baseTax, patriotismDiscountRate)
	
	// Apply cultural engagement discounts
	culturalDiscountRate := tc.calculateCulturalDiscount()
	culturalDiscountAmount := tc.applyDiscount(baseTax, culturalDiscountRate)
	
	// Calculate total discounts
	totalDiscountAmount := volumeDiscountAmount.Add(patriotismDiscountAmount).Add(culturalDiscountAmount)
	
	// Calculate final tax amount
	finalTaxAmount := baseTax.Sub(totalDiscountAmount)
	
	// Apply tax cap
	cappedTaxAmount := tc.applyTaxCap(finalTaxAmount, transactionAmount)
	
	// Calculate effective tax rate
	effectiveRate := tc.calculateEffectiveRate(cappedTaxAmount, transactionAmount)
	
	// Prepare applied discounts list
	appliedDiscounts := tc.buildAppliedDiscountsList(volumeDiscountRate, patriotismDiscountRate, culturalDiscountRate)
	
	// Calculate savings
	savingsAmount := baseTax.Sub(cappedTaxAmount)
	
	// Calculate optimization score
	optimizationScore := tc.calculateOptimizationScore(baseTax, cappedTaxAmount)
	
	// Build optimizations list
	optimizations := tc.buildOptimizationsList()
	
	result := &TaxCalculationResult{
		OriginalAmount:    transactionAmount,
		TaxAmount:         cappedTaxAmount,
		EffectiveRate:     effectiveRate,
		AppliedDiscounts:  appliedDiscounts,
		TaxBreakdown:      tc.createTaxBreakdown(transactionAmount, cappedTaxAmount),
		Optimizations:     optimizations,
		SavingsAmount:     savingsAmount,
		OptimizationScore: optimizationScore,
	}
	
	return result, nil
}

// calculateBaseTax calculates the base tax amount with progressive structure
func (tc *TaxCalculator) calculateBaseTax(amount sdk.Coin) sdk.Coin {
	// Convert amount to INR equivalent (assuming 1 NAMO = 1 INR with 6 decimals)
	amountINR := amount.Amount.Int64() / 1000000
	
	// Progressive tax structure
	var taxAmount int64
	
	switch {
	case amountINR < 100:
		// < ₹100: FREE (0%)
		taxAmount = 0
	case amountINR >= 100 && amountINR < 500:
		// ₹100 - ₹500: ₹0.01 fixed
		taxAmount = 10000 // 0.01 NAMO in micro units
	case amountINR >= 500 && amountINR < 1000:
		// ₹500 - ₹1,000: ₹0.05 fixed
		taxAmount = 50000 // 0.05 NAMO in micro units
	case amountINR >= 1000 && amountINR < 10000:
		// ₹1,000 - ₹10,000: 0.25%
		taxAmount = amount.Amount.Int64() * 25 / 10000
	case amountINR >= 10000 && amountINR < 100000:
		// ₹10,000 - ₹1 lakh: 0.5%
		taxAmount = amount.Amount.Int64() * 50 / 10000
	case amountINR >= 100000 && amountINR < 1000000:
		// ₹1 lakh - ₹10 lakh: 0.3%
		taxAmount = amount.Amount.Int64() * 30 / 10000
	default:
		// > ₹10 lakh: 0.2%
		taxAmount = amount.Amount.Int64() * 20 / 10000
	}
	
	// Apply ₹1,000 cap for large transactions
	maxTaxINR := int64(1000) * 1000000 // ₹1,000 in micro units
	if taxAmount > maxTaxINR {
		taxAmount = maxTaxINR
	}
	
	return sdk.NewCoin("namo", math.NewInt(taxAmount))
}

// calculateVolumeDiscount calculates the volume-based discount rate
func (tc *TaxCalculator) calculateVolumeDiscount() float64 {
	if tc.volumeData == nil || tc.volumeData.DailyTransactionCount == 0 {
		return 0.0
	}
	
	dailyCount := tc.volumeData.DailyTransactionCount
	
	// Apply volume-based tax reduction
	switch {
	case dailyCount >= VolumeThreshold10M:
		rate, _ := strconv.ParseFloat(VolumeDiscount10M, 64)
		return rate
	case dailyCount >= VolumeThreshold1M:
		rate, _ := strconv.ParseFloat(VolumeDiscount1M, 64)
		return rate
	case dailyCount >= VolumeThreshold500K:
		rate, _ := strconv.ParseFloat(VolumeDiscount500K, 64)
		return rate
	case dailyCount >= VolumeThreshold100K:
		rate, _ := strconv.ParseFloat(VolumeDiscount100K, 64)
		return rate
	case dailyCount >= VolumeThreshold50K:
		rate, _ := strconv.ParseFloat(VolumeDiscount50K, 64)
		return rate
	case dailyCount >= VolumeThreshold10K:
		rate, _ := strconv.ParseFloat(VolumeDiscount10K, 64)
		return rate
	default:
		return 0.0
	}
}

// calculatePatriotismDiscount calculates the patriotism-based discount rate
func (tc *TaxCalculator) calculatePatriotismDiscount() float64 {
	if tc.patriotismScore <= 0 {
		return 0.0
	}
	
	// Calculate patriotism discount: 0.5% per 100 patriotism score
	patriotismDiscountRate, _ := strconv.ParseFloat(DefaultPatriotismDiscountRate, 64)
	scoreMultiplier := float64(tc.patriotismScore) / float64(PatriotismScoreMultiplier)
	
	return patriotismDiscountRate * scoreMultiplier
}

// calculateCulturalDiscount calculates the cultural engagement discount rate
func (tc *TaxCalculator) calculateCulturalDiscount() float64 {
	if tc.culturalScore <= 0 {
		return 0.0
	}
	
	// Calculate cultural discount: 0.2% per 100 cultural score
	culturalDiscountRate, _ := strconv.ParseFloat(DefaultCulturalBonusRate, 64)
	scoreMultiplier := float64(tc.culturalScore) / float64(CulturalScoreMultiplier)
	
	return culturalDiscountRate * scoreMultiplier
}

// applyDiscount applies a discount rate to a tax amount
func (tc *TaxCalculator) applyDiscount(taxAmount sdk.Coin, discountRate float64) sdk.Coin {
	if discountRate <= 0 {
		return sdk.NewCoin(taxAmount.Denom, math.ZeroInt())
	}
	
	// Convert to big.Int for precise calculation
	taxBigInt := new(big.Int).Set(taxAmount.Amount.BigInt())
	
	// Calculate discount amount
	discountRateInt := big.NewInt(int64(discountRate * float64(PercentageMultiplier) * 100))
	discountAmount := new(big.Int).Mul(taxBigInt, discountRateInt)
	discountAmount = discountAmount.Div(discountAmount, big.NewInt(PercentageMultiplier*100))
	
	return sdk.NewCoin(taxAmount.Denom, math.NewIntFromBigInt(discountAmount))
}

// applyTaxCap ensures tax doesn't exceed ₹1,000 cap (already handled in calculateBaseTax)
func (tc *TaxCalculator) applyTaxCap(taxAmount sdk.Coin, transactionAmount sdk.Coin) sdk.Coin {
	// Tax cap is already applied in calculateBaseTax method
	// This method is kept for backward compatibility
	return taxAmount
}

// calculateEffectiveRate calculates the effective tax rate
func (tc *TaxCalculator) calculateEffectiveRate(taxAmount sdk.Coin, transactionAmount sdk.Coin) string {
	if transactionAmount.Amount.IsZero() {
		return "0.0"
	}
	
	// Calculate effective rate: (tax / transaction) * 100
	taxBigInt := new(big.Int).Set(taxAmount.Amount.BigInt())
	transactionBigInt := new(big.Int).Set(transactionAmount.Amount.BigInt())
	
	// Multiply by 10000 for precision (4 decimal places)
	effectiveRateBigInt := new(big.Int).Mul(taxBigInt, big.NewInt(10000))
	effectiveRateBigInt = effectiveRateBigInt.Div(effectiveRateBigInt, transactionBigInt)
	
	effectiveRate := float64(effectiveRateBigInt.Int64()) / 100.0
	
	return strconv.FormatFloat(effectiveRate, 'f', 4, 64)
}

// buildAppliedDiscountsList builds a list of applied discounts
func (tc *TaxCalculator) buildAppliedDiscountsList(volumeDiscount, patriotismDiscount, culturalDiscount float64) []string {
	var discounts []string
	
	if volumeDiscount > 0 {
		discounts = append(discounts, "volume_discount")
	}
	
	if patriotismDiscount > 0 {
		discounts = append(discounts, "patriotism_discount")
	}
	
	if culturalDiscount > 0 {
		discounts = append(discounts, "cultural_discount")
	}
	
	return discounts
}

// buildOptimizationsList builds a list of applied optimizations
func (tc *TaxCalculator) buildOptimizationsList() []string {
	var optimizations []string
	
	if tc.optimizationMode {
		optimizations = append(optimizations, "automatic_optimization")
	}
	
	if tc.volumeData != nil && tc.volumeData.DailyTransactionCount > VolumeThreshold1K {
		optimizations = append(optimizations, "volume_optimization")
	}
	
	if tc.patriotismScore > 0 {
		optimizations = append(optimizations, "patriotism_optimization")
	}
	
	if tc.culturalScore > 0 {
		optimizations = append(optimizations, "cultural_optimization")
	}
	
	return optimizations
}

// calculateOptimizationScore calculates the optimization score
func (tc *TaxCalculator) calculateOptimizationScore(baseTax, finalTax sdk.Coin) string {
	if baseTax.Amount.IsZero() {
		return "0"
	}
	
	// Calculate savings percentage
	savingsAmount := baseTax.Amount.Sub(finalTax.Amount)
	savingsPercentage := new(big.Int).Mul(savingsAmount.BigInt(), big.NewInt(100))
	savingsPercentage = savingsPercentage.Div(savingsPercentage, baseTax.Amount.BigInt())
	
	return savingsPercentage.String()
}

// createTaxBreakdown creates a detailed tax breakdown
func (tc *TaxCalculator) createTaxBreakdown(transactionAmount sdk.Coin, taxAmount sdk.Coin) *TaxBreakdown {
	baseTax := tc.calculateBaseTax(transactionAmount)
	
	volumeDiscount := tc.applyDiscount(baseTax, tc.calculateVolumeDiscount())
	patriotismDiscount := tc.applyDiscount(baseTax, tc.calculatePatriotismDiscount())
	culturalDiscount := tc.applyDiscount(baseTax, tc.calculateCulturalDiscount())
	
	totalDiscounts := volumeDiscount.Add(patriotismDiscount).Add(culturalDiscount)
	
	return &TaxBreakdown{
		BaseTaxAmount:        baseTax,
		VolumeDiscountAmount: volumeDiscount,
		PatriotismDiscountAmount: patriotismDiscount,
		CulturalDiscountAmount: culturalDiscount,
		TotalDiscountAmount:  totalDiscounts,
		FinalTaxAmount:       taxAmount,
		EffectiveTaxRate:     tc.calculateEffectiveRate(taxAmount, transactionAmount),
	}
}

// TaxCalculationResult represents the result of a tax calculation
type TaxCalculationResult struct {
	OriginalAmount    sdk.Coin
	TaxAmount         sdk.Coin
	EffectiveRate     string
	AppliedDiscounts  []string
	TaxBreakdown      *TaxBreakdown
	Optimizations     []string
	SavingsAmount     sdk.Coin
	OptimizationScore string
}

// TaxBreakdown provides detailed breakdown of tax calculation
type TaxBreakdown struct {
	BaseTaxAmount            sdk.Coin
	VolumeDiscountAmount     sdk.Coin
	PatriotismDiscountAmount sdk.Coin
	CulturalDiscountAmount   sdk.Coin
	TotalDiscountAmount      sdk.Coin
	FinalTaxAmount           sdk.Coin
	EffectiveTaxRate         string
}

// VolumeData represents volume-based data for tax calculations
type VolumeData struct {
	DailyTransactionCount   uint64
	WeeklyTransactionCount  uint64
	MonthlyTransactionCount uint64
	DailyVolume             sdk.Coin
	WeeklyVolume            sdk.Coin
	MonthlyVolume           sdk.Coin
	AverageTransactionSize  sdk.Coin
	PeakTransactionHour     int32
	VolumeGrowthRate        string
	VolumeStabilityScore    string
}

// TaxOptimizer provides advanced tax optimization strategies
type TaxOptimizer struct {
	calculator *TaxCalculator
}

// NewTaxOptimizer creates a new tax optimizer
func NewTaxOptimizer(calculator *TaxCalculator) *TaxOptimizer {
	return &TaxOptimizer{
		calculator: calculator,
	}
}

// OptimizeTaxStrategy provides optimization recommendations
func (to *TaxOptimizer) OptimizeTaxStrategy(userProfile *UserTaxProfile, transactionAmount sdk.Coin) *TaxOptimizationResult {
	result := &TaxOptimizationResult{
		UserProfile:       userProfile,
		TransactionAmount: transactionAmount,
		Recommendations:   []string{},
		PotentialSavings:  sdk.NewCoin(transactionAmount.Denom, math.ZeroInt()),
		OptimizationScore: "0",
	}
	
	// Analyze current tax situation
	currentTax, _ := to.calculator.CalculateTax(transactionAmount, "transfer")
	
	// Provide recommendations based on profile
	if userProfile.PatriotismScore < 5000 {
		result.Recommendations = append(result.Recommendations, "Increase patriotism score through donations and cultural engagement")
	}
	
	if userProfile.CulturalEngagementScore < 3000 {
		result.Recommendations = append(result.Recommendations, "Participate in cultural activities to earn cultural engagement bonus")
	}
	
	if userProfile.TransactionVolume.Amount.LT(math.NewInt(1000000)) {
		result.Recommendations = append(result.Recommendations, "Consider consolidating transactions to benefit from volume discounts")
	}
	
	// Calculate potential savings
	if len(result.Recommendations) > 0 {
		// Estimate potential savings (simplified calculation)
		potentialSavings := currentTax.TaxAmount.Amount.Quo(math.NewInt(4)) // 25% potential savings
		result.PotentialSavings = sdk.NewCoin(transactionAmount.Denom, potentialSavings)
		result.OptimizationScore = "75"
	}
	
	return result
}

// TaxOptimizationResult represents the result of tax optimization analysis
type TaxOptimizationResult struct {
	UserProfile       *UserTaxProfile
	TransactionAmount sdk.Coin
	Recommendations   []string
	PotentialSavings  sdk.Coin
	OptimizationScore string
}

// TaxComplianceChecker checks tax compliance
type TaxComplianceChecker struct {
	config *TaxConfig
}

// NewTaxComplianceChecker creates a new compliance checker
func NewTaxComplianceChecker(config *TaxConfig) *TaxComplianceChecker {
	return &TaxComplianceChecker{
		config: config,
	}
}

// CheckCompliance checks if a transaction is compliant
func (tcc *TaxComplianceChecker) CheckCompliance(transaction *TaxTransaction) *ComplianceResult {
	result := &ComplianceResult{
		IsCompliant:     true,
		ComplianceLevel: ComplianceLevelStandard,
		Violations:      []string{},
		Recommendations: []string{},
		RiskScore:       "LOW",
	}
	
	// Check for compliance violations
	if transaction.TaxAmount.Amount.IsZero() && !transaction.DonationFlag {
		result.Violations = append(result.Violations, "Zero tax amount for non-donation transaction")
		result.IsCompliant = false
	}
	
	// Check for suspicious patterns
	if transaction.OptimizationApplied && len(transaction.AppliedDiscounts) > 3 {
		result.Recommendations = append(result.Recommendations, "High number of discounts applied - monitor for abuse")
		result.RiskScore = "MEDIUM"
	}
	
	// Set compliance level based on violations
	if len(result.Violations) > 0 {
		result.ComplianceLevel = ComplianceLevelBasic
	}
	
	return result
}

// ComplianceResult represents the result of compliance checking
type ComplianceResult struct {
	IsCompliant     bool
	ComplianceLevel string
	Violations      []string
	Recommendations []string
	RiskScore       string
}

// TaxReporter generates tax reports
type TaxReporter struct {
	calculator *TaxCalculator
}

// NewTaxReporter creates a new tax reporter
func NewTaxReporter(calculator *TaxCalculator) *TaxReporter {
	return &TaxReporter{
		calculator: calculator,
	}
}

// GenerateReport generates a tax report
func (tr *TaxReporter) GenerateReport(reportType string, startTime, endTime int64) *TaxReport {
	return &TaxReport{
		ReportType:        reportType,
		StartTime:         startTime,
		EndTime:           endTime,
		TotalTransactions: 0,
		TotalTaxCollected: sdk.NewCoin("namo", math.ZeroInt()),
		AverageEffectiveRate: "0.0",
		TopOptimizations:  []string{},
		ComplianceScore:   "100",
		GeneratedAt:       startTime, // In production, use current time
	}
}

// TaxReport represents a tax report
type TaxReport struct {
	ReportType           string
	StartTime            int64
	EndTime              int64
	TotalTransactions    uint64
	TotalTaxCollected    sdk.Coin
	AverageEffectiveRate string
	TopOptimizations     []string
	ComplianceScore      string
	GeneratedAt          int64
}

// TaxForecaster provides tax forecasting capabilities
type TaxForecaster struct {
	calculator *TaxCalculator
}

// NewTaxForecaster creates a new tax forecaster
func NewTaxForecaster(calculator *TaxCalculator) *TaxForecaster {
	return &TaxForecaster{
		calculator: calculator,
	}
}

// ForecastTax forecasts future tax revenue
func (tf *TaxForecaster) ForecastTax(forecastPeriod string, currentVolume sdk.Coin) *TaxForecast {
	return &TaxForecast{
		ForecastPeriod:    forecastPeriod,
		CurrentVolume:     currentVolume,
		ProjectedVolume:   currentVolume, // Simplified - use growth models in production
		ProjectedTax:      sdk.NewCoin("namo", math.ZeroInt()),
		ConfidenceLevel:   "HIGH",
		GrowthRate:        "0.0",
		SeasonalFactors:   []string{},
		RiskFactors:       []string{},
		ForecastedAt:      0, // In production, use current time
	}
}

// TaxForecast represents a tax forecast
type TaxForecast struct {
	ForecastPeriod  string
	CurrentVolume   sdk.Coin
	ProjectedVolume sdk.Coin
	ProjectedTax    sdk.Coin
	ConfidenceLevel string
	GrowthRate      string
	SeasonalFactors []string
	RiskFactors     []string
	ForecastedAt    int64
}

// TaxEducator provides tax education content
type TaxEducator struct {
	config *TaxConfig
}

// NewTaxEducator creates a new tax educator
func NewTaxEducator(config *TaxConfig) *TaxEducator {
	return &TaxEducator{
		config: config,
	}
}

// GetEducationContent provides educational content about taxes
func (te *TaxEducator) GetEducationContent(contentType string) *TaxEducationContent {
	return &TaxEducationContent{
		ContentType:   contentType,
		Title:         "Understanding DeshChain Tax System",
		Description:   "Learn how DeshChain's innovative tax system works",
		Content:       "DeshChain implements a progressive tax system with volume-based discounts...",
		Examples:      []string{},
		References:    []string{},
		QuizQuestions: []string{},
		CreatedAt:     0, // In production, use current time
		UpdatedAt:     0, // In production, use current time
	}
}

// TaxEducationContent represents educational content about taxes
type TaxEducationContent struct {
	ContentType   string
	Title         string
	Description   string
	Content       string
	Examples      []string
	References    []string
	QuizQuestions []string
	CreatedAt     int64
	UpdatedAt     int64
}