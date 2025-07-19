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
	"fmt"
	"strconv"
	"strings"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys
var (
	KeyBaseTaxRate             = []byte("BaseTaxRate")
	KeyMaxTaxAmount            = []byte("MaxTaxAmount")
	KeyMinTaxAmount            = []byte("MinTaxAmount")
	KeyVolumeDiscountEnabled   = []byte("VolumeDiscountEnabled")
	KeyPatriotismDiscountEnabled = []byte("PatriotismDiscountEnabled")
	KeyCulturalDiscountEnabled = []byte("CulturalDiscountEnabled")
	KeyDonationExemptionEnabled = []byte("DonationExemptionEnabled")
	KeyOptimizationEnabled     = []byte("OptimizationEnabled")
	KeyProgressiveTaxEnabled   = []byte("ProgressiveTaxEnabled")
	KeyComplianceRequired      = []byte("ComplianceRequired")
	KeyAuditTrailEnabled       = []byte("AuditTrailEnabled")
	KeyEducationEnabled        = []byte("EducationEnabled")
	KeyForecastingEnabled      = []byte("ForecastingEnabled")
	KeyReportingEnabled        = []byte("ReportingEnabled")
	KeyTransparencyEnabled     = []byte("TransparencyEnabled")
	KeyTaxCapResetPeriod       = []byte("TaxCapResetPeriod")
	KeyGracePeriodDays         = []byte("GracePeriodDays")
	KeyMaxVolumeDiscount       = []byte("MaxVolumeDiscount")
	KeyEarlyPaymentDiscount    = []byte("EarlyPaymentDiscount")
	KeyLatePaymentPenalty      = []byte("LatePaymentPenalty")
	KeyPatriotismDiscountRate  = []byte("PatriotismDiscountRate")
	KeyCulturalBonusRate       = []byte("CulturalBonusRate")
	KeyTaxCalculationPrecision = []byte("TaxCalculationPrecision")
	KeyVolumeThresholds        = []byte("VolumeThresholds")
	KeyTaxBrackets             = []byte("TaxBrackets")
	KeySystemMode              = []byte("SystemMode")
	KeyEmergencyMode           = []byte("EmergencyMode")
	KeyMaintenanceMode         = []byte("MaintenanceMode")
	KeyUpgradeMode             = []byte("UpgradeMode")
	KeyAuditMode               = []byte("AuditMode")
	KeyDebugMode               = []byte("DebugMode")
	KeyTestMode                = []byte("TestMode")
	KeyPerformanceMode         = []byte("PerformanceMode")
)

// Default parameter values
var (
	DefaultBaseTaxRate             = "0.025"    // 2.5%
	DefaultMaxTaxAmount            = "1000"     // ₹1,000 equivalent
	DefaultMinTaxAmount            = "0"        // ₹0
	DefaultVolumeDiscountEnabled   = true
	DefaultPatriotismDiscountEnabled = true
	DefaultCulturalDiscountEnabled = true
	DefaultDonationExemptionEnabled = true
	DefaultOptimizationEnabled     = true
	DefaultProgressiveTaxEnabled   = true
	DefaultComplianceRequired      = true
	DefaultAuditTrailEnabled       = true
	DefaultEducationEnabled        = true
	DefaultForecastingEnabled      = true
	DefaultReportingEnabled        = true
	DefaultTransparencyEnabled     = true
	DefaultTaxCapResetPeriod       = int64(86400) // 24 hours
	DefaultGracePeriodDays         = int64(7)     // 7 days
	DefaultMaxVolumeDiscount       = "0.9"        // 90% max discount
	DefaultEarlyPaymentDiscount    = "0.01"       // 1%
	DefaultLatePaymentPenalty      = "0.02"       // 2%
	DefaultPatriotismDiscountRate  = "0.005"      // 0.5%
	DefaultCulturalBonusRate       = "0.002"      // 0.2%
	DefaultTaxCalculationPrecision = int64(6)     // 6 decimal places
	DefaultSystemMode              = SystemModeNormal
	DefaultEmergencyMode           = false
	DefaultMaintenanceMode         = false
	DefaultUpgradeMode             = false
	DefaultAuditMode               = false
	DefaultDebugMode               = false
	DefaultTestMode                = false
	DefaultPerformanceMode         = false
)

// DefaultVolumeThresholds returns default volume thresholds
func DefaultVolumeThresholds() []*VolumeThreshold {
	return []*VolumeThreshold{
		{
			TransactionCount: VolumeThreshold1K,
			TaxRate:          VolumeDiscount1K,
			Description:      "1K transactions daily",
		},
		{
			TransactionCount: VolumeThreshold10K,
			TaxRate:          VolumeDiscount10K,
			Description:      "10K transactions daily",
		},
		{
			TransactionCount: VolumeThreshold50K,
			TaxRate:          VolumeDiscount50K,
			Description:      "50K transactions daily",
		},
		{
			TransactionCount: VolumeThreshold100K,
			TaxRate:          VolumeDiscount100K,
			Description:      "100K transactions daily",
		},
		{
			TransactionCount: VolumeThreshold500K,
			TaxRate:          VolumeDiscount500K,
			Description:      "500K transactions daily",
		},
		{
			TransactionCount: VolumeThreshold1M,
			TaxRate:          VolumeDiscount1M,
			Description:      "1M transactions daily",
		},
		{
			TransactionCount: VolumeThreshold10M,
			TaxRate:          VolumeDiscount10M,
			Description:      "10M+ transactions daily",
		},
	}
}

// DefaultTaxBrackets returns default tax brackets
func DefaultTaxBrackets() []*TaxBracket {
	return []*TaxBracket{
		{
			AmountLimit: TaxBracket1Limit,
			TaxRate:     TaxBracket1Rate,
			TaxCap:      "", // No cap for first bracket
			Description: "₹0 - ₹40,000: Full percentage",
		},
		{
			AmountLimit: TaxBracket2Limit,
			TaxRate:     TaxBracket1Rate,
			TaxCap:      TaxBracket2Cap,
			Description: "₹40,001 - ₹4,00,000: ₹1,000 cap",
		},
		{
			AmountLimit: TaxBracket3Limit,
			TaxRate:     TaxBracket1Rate,
			TaxCap:      TaxBracket3Cap,
			Description: "Above ₹4,00,000: Flat ₹1,000",
		},
	}
}

// TaxParams defines the parameters for the tax module
type TaxParams struct {
	// Basic tax configuration
	BaseTaxRate             string `protobuf:"bytes,1,opt,name=base_tax_rate,json=baseTaxRate,proto3" json:"base_tax_rate,omitempty"`
	MaxTaxAmount            string `protobuf:"bytes,2,opt,name=max_tax_amount,json=maxTaxAmount,proto3" json:"max_tax_amount,omitempty"`
	MinTaxAmount            string `protobuf:"bytes,3,opt,name=min_tax_amount,json=minTaxAmount,proto3" json:"min_tax_amount,omitempty"`
	
	// Feature flags
	VolumeDiscountEnabled   bool `protobuf:"varint,4,opt,name=volume_discount_enabled,json=volumeDiscountEnabled,proto3" json:"volume_discount_enabled,omitempty"`
	PatriotismDiscountEnabled bool `protobuf:"varint,5,opt,name=patriotism_discount_enabled,json=patriotismDiscountEnabled,proto3" json:"patriotism_discount_enabled,omitempty"`
	CulturalDiscountEnabled bool `protobuf:"varint,6,opt,name=cultural_discount_enabled,json=culturalDiscountEnabled,proto3" json:"cultural_discount_enabled,omitempty"`
	DonationExemptionEnabled bool `protobuf:"varint,7,opt,name=donation_exemption_enabled,json=donationExemptionEnabled,proto3" json:"donation_exemption_enabled,omitempty"`
	OptimizationEnabled     bool `protobuf:"varint,8,opt,name=optimization_enabled,json=optimizationEnabled,proto3" json:"optimization_enabled,omitempty"`
	ProgressiveTaxEnabled   bool `protobuf:"varint,9,opt,name=progressive_tax_enabled,json=progressiveTaxEnabled,proto3" json:"progressive_tax_enabled,omitempty"`
	ComplianceRequired      bool `protobuf:"varint,10,opt,name=compliance_required,json=complianceRequired,proto3" json:"compliance_required,omitempty"`
	AuditTrailEnabled       bool `protobuf:"varint,11,opt,name=audit_trail_enabled,json=auditTrailEnabled,proto3" json:"audit_trail_enabled,omitempty"`
	EducationEnabled        bool `protobuf:"varint,12,opt,name=education_enabled,json=educationEnabled,proto3" json:"education_enabled,omitempty"`
	ForecastingEnabled      bool `protobuf:"varint,13,opt,name=forecasting_enabled,json=forecastingEnabled,proto3" json:"forecasting_enabled,omitempty"`
	ReportingEnabled        bool `protobuf:"varint,14,opt,name=reporting_enabled,json=reportingEnabled,proto3" json:"reporting_enabled,omitempty"`
	TransparencyEnabled     bool `protobuf:"varint,15,opt,name=transparency_enabled,json=transparencyEnabled,proto3" json:"transparency_enabled,omitempty"`
	
	// Timing parameters
	TaxCapResetPeriod       int64 `protobuf:"varint,16,opt,name=tax_cap_reset_period,json=taxCapResetPeriod,proto3" json:"tax_cap_reset_period,omitempty"`
	GracePeriodDays         int64 `protobuf:"varint,17,opt,name=grace_period_days,json=gracePeriodDays,proto3" json:"grace_period_days,omitempty"`
	
	// Discount and penalty rates
	MaxVolumeDiscount       string `protobuf:"bytes,18,opt,name=max_volume_discount,json=maxVolumeDiscount,proto3" json:"max_volume_discount,omitempty"`
	EarlyPaymentDiscount    string `protobuf:"bytes,19,opt,name=early_payment_discount,json=earlyPaymentDiscount,proto3" json:"early_payment_discount,omitempty"`
	LatePaymentPenalty      string `protobuf:"bytes,20,opt,name=late_payment_penalty,json=latePaymentPenalty,proto3" json:"late_payment_penalty,omitempty"`
	PatriotismDiscountRate  string `protobuf:"bytes,21,opt,name=patriotism_discount_rate,json=patriotismDiscountRate,proto3" json:"patriotism_discount_rate,omitempty"`
	CulturalBonusRate       string `protobuf:"bytes,22,opt,name=cultural_bonus_rate,json=culturalBonusRate,proto3" json:"cultural_bonus_rate,omitempty"`
	
	// Technical parameters
	TaxCalculationPrecision int64 `protobuf:"varint,23,opt,name=tax_calculation_precision,json=taxCalculationPrecision,proto3" json:"tax_calculation_precision,omitempty"`
	
	// Complex parameters
	VolumeThresholds        []*VolumeThreshold `protobuf:"bytes,24,rep,name=volume_thresholds,json=volumeThresholds,proto3" json:"volume_thresholds,omitempty"`
	TaxBrackets             []*TaxBracket      `protobuf:"bytes,25,rep,name=tax_brackets,json=taxBrackets,proto3" json:"tax_brackets,omitempty"`
	
	// System modes
	SystemMode              string `protobuf:"bytes,26,opt,name=system_mode,json=systemMode,proto3" json:"system_mode,omitempty"`
	EmergencyMode           bool   `protobuf:"varint,27,opt,name=emergency_mode,json=emergencyMode,proto3" json:"emergency_mode,omitempty"`
	MaintenanceMode         bool   `protobuf:"varint,28,opt,name=maintenance_mode,json=maintenanceMode,proto3" json:"maintenance_mode,omitempty"`
	UpgradeMode             bool   `protobuf:"varint,29,opt,name=upgrade_mode,json=upgradeMode,proto3" json:"upgrade_mode,omitempty"`
	AuditMode               bool   `protobuf:"varint,30,opt,name=audit_mode,json=auditMode,proto3" json:"audit_mode,omitempty"`
	DebugMode               bool   `protobuf:"varint,31,opt,name=debug_mode,json=debugMode,proto3" json:"debug_mode,omitempty"`
	TestMode                bool   `protobuf:"varint,32,opt,name=test_mode,json=testMode,proto3" json:"test_mode,omitempty"`
	PerformanceMode         bool   `protobuf:"varint,33,opt,name=performance_mode,json=performanceMode,proto3" json:"performance_mode,omitempty"`
}

// NewTaxParams creates a new TaxParams instance
func NewTaxParams(
	baseTaxRate,
	maxTaxAmount,
	minTaxAmount string,
	volumeDiscountEnabled,
	patriotismDiscountEnabled,
	culturalDiscountEnabled,
	donationExemptionEnabled,
	optimizationEnabled,
	progressiveTaxEnabled,
	complianceRequired,
	auditTrailEnabled,
	educationEnabled,
	forecastingEnabled,
	reportingEnabled,
	transparencyEnabled bool,
	taxCapResetPeriod,
	gracePeriodDays int64,
	maxVolumeDiscount,
	earlyPaymentDiscount,
	latePaymentPenalty,
	patriotismDiscountRate,
	culturalBonusRate string,
	taxCalculationPrecision int64,
	volumeThresholds []*VolumeThreshold,
	taxBrackets []*TaxBracket,
	systemMode string,
	emergencyMode,
	maintenanceMode,
	upgradeMode,
	auditMode,
	debugMode,
	testMode,
	performanceMode bool,
) TaxParams {
	return TaxParams{
		BaseTaxRate:             baseTaxRate,
		MaxTaxAmount:            maxTaxAmount,
		MinTaxAmount:            minTaxAmount,
		VolumeDiscountEnabled:   volumeDiscountEnabled,
		PatriotismDiscountEnabled: patriotismDiscountEnabled,
		CulturalDiscountEnabled: culturalDiscountEnabled,
		DonationExemptionEnabled: donationExemptionEnabled,
		OptimizationEnabled:     optimizationEnabled,
		ProgressiveTaxEnabled:   progressiveTaxEnabled,
		ComplianceRequired:      complianceRequired,
		AuditTrailEnabled:       auditTrailEnabled,
		EducationEnabled:        educationEnabled,
		ForecastingEnabled:      forecastingEnabled,
		ReportingEnabled:        reportingEnabled,
		TransparencyEnabled:     transparencyEnabled,
		TaxCapResetPeriod:       taxCapResetPeriod,
		GracePeriodDays:         gracePeriodDays,
		MaxVolumeDiscount:       maxVolumeDiscount,
		EarlyPaymentDiscount:    earlyPaymentDiscount,
		LatePaymentPenalty:      latePaymentPenalty,
		PatriotismDiscountRate:  patriotismDiscountRate,
		CulturalBonusRate:       culturalBonusRate,
		TaxCalculationPrecision: taxCalculationPrecision,
		VolumeThresholds:        volumeThresholds,
		TaxBrackets:             taxBrackets,
		SystemMode:              systemMode,
		EmergencyMode:           emergencyMode,
		MaintenanceMode:         maintenanceMode,
		UpgradeMode:             upgradeMode,
		AuditMode:               auditMode,
		DebugMode:               debugMode,
		TestMode:                testMode,
		PerformanceMode:         performanceMode,
	}
}

// DefaultParams returns default tax parameters
func DefaultParams() TaxParams {
	return NewTaxParams(
		DefaultBaseTaxRate,
		DefaultMaxTaxAmount,
		DefaultMinTaxAmount,
		DefaultVolumeDiscountEnabled,
		DefaultPatriotismDiscountEnabled,
		DefaultCulturalDiscountEnabled,
		DefaultDonationExemptionEnabled,
		DefaultOptimizationEnabled,
		DefaultProgressiveTaxEnabled,
		DefaultComplianceRequired,
		DefaultAuditTrailEnabled,
		DefaultEducationEnabled,
		DefaultForecastingEnabled,
		DefaultReportingEnabled,
		DefaultTransparencyEnabled,
		DefaultTaxCapResetPeriod,
		DefaultGracePeriodDays,
		DefaultMaxVolumeDiscount,
		DefaultEarlyPaymentDiscount,
		DefaultLatePaymentPenalty,
		DefaultPatriotismDiscountRate,
		DefaultCulturalBonusRate,
		DefaultTaxCalculationPrecision,
		DefaultVolumeThresholds(),
		DefaultTaxBrackets(),
		DefaultSystemMode,
		DefaultEmergencyMode,
		DefaultMaintenanceMode,
		DefaultUpgradeMode,
		DefaultAuditMode,
		DefaultDebugMode,
		DefaultTestMode,
		DefaultPerformanceMode,
	)
}

// ParamKeyTable returns the parameter key table for use with the sdk.Params
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&TaxParams{})
}

// ParamSetPairs returns the parameter set pairs
func (p *TaxParams) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyBaseTaxRate, &p.BaseTaxRate, validateBaseTaxRate),
		paramtypes.NewParamSetPair(KeyMaxTaxAmount, &p.MaxTaxAmount, validateMaxTaxAmount),
		paramtypes.NewParamSetPair(KeyMinTaxAmount, &p.MinTaxAmount, validateMinTaxAmount),
		paramtypes.NewParamSetPair(KeyVolumeDiscountEnabled, &p.VolumeDiscountEnabled, validateBool),
		paramtypes.NewParamSetPair(KeyPatriotismDiscountEnabled, &p.PatriotismDiscountEnabled, validateBool),
		paramtypes.NewParamSetPair(KeyCulturalDiscountEnabled, &p.CulturalDiscountEnabled, validateBool),
		paramtypes.NewParamSetPair(KeyDonationExemptionEnabled, &p.DonationExemptionEnabled, validateBool),
		paramtypes.NewParamSetPair(KeyOptimizationEnabled, &p.OptimizationEnabled, validateBool),
		paramtypes.NewParamSetPair(KeyProgressiveTaxEnabled, &p.ProgressiveTaxEnabled, validateBool),
		paramtypes.NewParamSetPair(KeyComplianceRequired, &p.ComplianceRequired, validateBool),
		paramtypes.NewParamSetPair(KeyAuditTrailEnabled, &p.AuditTrailEnabled, validateBool),
		paramtypes.NewParamSetPair(KeyEducationEnabled, &p.EducationEnabled, validateBool),
		paramtypes.NewParamSetPair(KeyForecastingEnabled, &p.ForecastingEnabled, validateBool),
		paramtypes.NewParamSetPair(KeyReportingEnabled, &p.ReportingEnabled, validateBool),
		paramtypes.NewParamSetPair(KeyTransparencyEnabled, &p.TransparencyEnabled, validateBool),
		paramtypes.NewParamSetPair(KeyTaxCapResetPeriod, &p.TaxCapResetPeriod, validateTaxCapResetPeriod),
		paramtypes.NewParamSetPair(KeyGracePeriodDays, &p.GracePeriodDays, validateGracePeriodDays),
		paramtypes.NewParamSetPair(KeyMaxVolumeDiscount, &p.MaxVolumeDiscount, validateMaxVolumeDiscount),
		paramtypes.NewParamSetPair(KeyEarlyPaymentDiscount, &p.EarlyPaymentDiscount, validateEarlyPaymentDiscount),
		paramtypes.NewParamSetPair(KeyLatePaymentPenalty, &p.LatePaymentPenalty, validateLatePaymentPenalty),
		paramtypes.NewParamSetPair(KeyPatriotismDiscountRate, &p.PatriotismDiscountRate, validatePatriotismDiscountRate),
		paramtypes.NewParamSetPair(KeyCulturalBonusRate, &p.CulturalBonusRate, validateCulturalBonusRate),
		paramtypes.NewParamSetPair(KeyTaxCalculationPrecision, &p.TaxCalculationPrecision, validateTaxCalculationPrecision),
		paramtypes.NewParamSetPair(KeyVolumeThresholds, &p.VolumeThresholds, validateVolumeThresholds),
		paramtypes.NewParamSetPair(KeyTaxBrackets, &p.TaxBrackets, validateTaxBrackets),
		paramtypes.NewParamSetPair(KeySystemMode, &p.SystemMode, validateSystemMode),
		paramtypes.NewParamSetPair(KeyEmergencyMode, &p.EmergencyMode, validateBool),
		paramtypes.NewParamSetPair(KeyMaintenanceMode, &p.MaintenanceMode, validateBool),
		paramtypes.NewParamSetPair(KeyUpgradeMode, &p.UpgradeMode, validateBool),
		paramtypes.NewParamSetPair(KeyAuditMode, &p.AuditMode, validateBool),
		paramtypes.NewParamSetPair(KeyDebugMode, &p.DebugMode, validateBool),
		paramtypes.NewParamSetPair(KeyTestMode, &p.TestMode, validateBool),
		paramtypes.NewParamSetPair(KeyPerformanceMode, &p.PerformanceMode, validateBool),
	}
}

// Validate validates the tax parameters
func (p TaxParams) Validate() error {
	if err := validateBaseTaxRate(p.BaseTaxRate); err != nil {
		return err
	}
	if err := validateMaxTaxAmount(p.MaxTaxAmount); err != nil {
		return err
	}
	if err := validateMinTaxAmount(p.MinTaxAmount); err != nil {
		return err
	}
	if err := validateTaxCapResetPeriod(p.TaxCapResetPeriod); err != nil {
		return err
	}
	if err := validateGracePeriodDays(p.GracePeriodDays); err != nil {
		return err
	}
	if err := validateMaxVolumeDiscount(p.MaxVolumeDiscount); err != nil {
		return err
	}
	if err := validateEarlyPaymentDiscount(p.EarlyPaymentDiscount); err != nil {
		return err
	}
	if err := validateLatePaymentPenalty(p.LatePaymentPenalty); err != nil {
		return err
	}
	if err := validatePatriotismDiscountRate(p.PatriotismDiscountRate); err != nil {
		return err
	}
	if err := validateCulturalBonusRate(p.CulturalBonusRate); err != nil {
		return err
	}
	if err := validateTaxCalculationPrecision(p.TaxCalculationPrecision); err != nil {
		return err
	}
	if err := validateVolumeThresholds(p.VolumeThresholds); err != nil {
		return err
	}
	if err := validateTaxBrackets(p.TaxBrackets); err != nil {
		return err
	}
	if err := validateSystemMode(p.SystemMode); err != nil {
		return err
	}
	return nil
}

// String returns a human-readable string representation of the parameters
func (p TaxParams) String() string {
	var sb strings.Builder
	sb.WriteString("Tax Parameters:\n")
	sb.WriteString(fmt.Sprintf("  Base Tax Rate: %s\n", p.BaseTaxRate))
	sb.WriteString(fmt.Sprintf("  Max Tax Amount: %s\n", p.MaxTaxAmount))
	sb.WriteString(fmt.Sprintf("  Min Tax Amount: %s\n", p.MinTaxAmount))
	sb.WriteString(fmt.Sprintf("  Volume Discount Enabled: %t\n", p.VolumeDiscountEnabled))
	sb.WriteString(fmt.Sprintf("  Patriotism Discount Enabled: %t\n", p.PatriotismDiscountEnabled))
	sb.WriteString(fmt.Sprintf("  Cultural Discount Enabled: %t\n", p.CulturalDiscountEnabled))
	sb.WriteString(fmt.Sprintf("  Donation Exemption Enabled: %t\n", p.DonationExemptionEnabled))
	sb.WriteString(fmt.Sprintf("  Optimization Enabled: %t\n", p.OptimizationEnabled))
	sb.WriteString(fmt.Sprintf("  Progressive Tax Enabled: %t\n", p.ProgressiveTaxEnabled))
	sb.WriteString(fmt.Sprintf("  Compliance Required: %t\n", p.ComplianceRequired))
	sb.WriteString(fmt.Sprintf("  Audit Trail Enabled: %t\n", p.AuditTrailEnabled))
	sb.WriteString(fmt.Sprintf("  Education Enabled: %t\n", p.EducationEnabled))
	sb.WriteString(fmt.Sprintf("  Forecasting Enabled: %t\n", p.ForecastingEnabled))
	sb.WriteString(fmt.Sprintf("  Reporting Enabled: %t\n", p.ReportingEnabled))
	sb.WriteString(fmt.Sprintf("  Transparency Enabled: %t\n", p.TransparencyEnabled))
	sb.WriteString(fmt.Sprintf("  Tax Cap Reset Period: %d\n", p.TaxCapResetPeriod))
	sb.WriteString(fmt.Sprintf("  Grace Period Days: %d\n", p.GracePeriodDays))
	sb.WriteString(fmt.Sprintf("  Max Volume Discount: %s\n", p.MaxVolumeDiscount))
	sb.WriteString(fmt.Sprintf("  Early Payment Discount: %s\n", p.EarlyPaymentDiscount))
	sb.WriteString(fmt.Sprintf("  Late Payment Penalty: %s\n", p.LatePaymentPenalty))
	sb.WriteString(fmt.Sprintf("  Patriotism Discount Rate: %s\n", p.PatriotismDiscountRate))
	sb.WriteString(fmt.Sprintf("  Cultural Bonus Rate: %s\n", p.CulturalBonusRate))
	sb.WriteString(fmt.Sprintf("  Tax Calculation Precision: %d\n", p.TaxCalculationPrecision))
	sb.WriteString(fmt.Sprintf("  Volume Thresholds Count: %d\n", len(p.VolumeThresholds)))
	sb.WriteString(fmt.Sprintf("  Tax Brackets Count: %d\n", len(p.TaxBrackets)))
	sb.WriteString(fmt.Sprintf("  System Mode: %s\n", p.SystemMode))
	sb.WriteString(fmt.Sprintf("  Emergency Mode: %t\n", p.EmergencyMode))
	sb.WriteString(fmt.Sprintf("  Maintenance Mode: %t\n", p.MaintenanceMode))
	sb.WriteString(fmt.Sprintf("  Upgrade Mode: %t\n", p.UpgradeMode))
	sb.WriteString(fmt.Sprintf("  Audit Mode: %t\n", p.AuditMode))
	sb.WriteString(fmt.Sprintf("  Debug Mode: %t\n", p.DebugMode))
	sb.WriteString(fmt.Sprintf("  Test Mode: %t\n", p.TestMode))
	sb.WriteString(fmt.Sprintf("  Performance Mode: %t\n", p.PerformanceMode))
	return sb.String()
}

// Validation functions

func validateBaseTaxRate(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	rate, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return fmt.Errorf("invalid base tax rate: %s", v)
	}
	
	if rate < MinTaxRateValue || rate > MaxTaxRateValue {
		return fmt.Errorf("base tax rate must be between %f and %f", MinTaxRateValue, MaxTaxRateValue)
	}
	
	return nil
}

func validateMaxTaxAmount(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	amount, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid max tax amount: %s", v)
	}
	
	if amount < MinTaxAmount || amount > MaxTaxAmount {
		return fmt.Errorf("max tax amount must be between %d and %d", MinTaxAmount, MaxTaxAmount)
	}
	
	return nil
}

func validateMinTaxAmount(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	amount, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid min tax amount: %s", v)
	}
	
	if amount < 0 {
		return fmt.Errorf("min tax amount cannot be negative")
	}
	
	return nil
}

func validateTaxCapResetPeriod(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	if v < SecondsInHour || v > SecondsInWeek {
		return fmt.Errorf("tax cap reset period must be between %d and %d seconds", SecondsInHour, SecondsInWeek)
	}
	
	return nil
}

func validateGracePeriodDays(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	if v < 0 || v > 30 {
		return fmt.Errorf("grace period days must be between 0 and 30")
	}
	
	return nil
}

func validateMaxVolumeDiscount(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	discount, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return fmt.Errorf("invalid max volume discount: %s", v)
	}
	
	if discount < MinDiscountRate || discount > MaxDiscountRate {
		return fmt.Errorf("max volume discount must be between %f and %f", MinDiscountRate, MaxDiscountRate)
	}
	
	return nil
}

func validateEarlyPaymentDiscount(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	discount, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return fmt.Errorf("invalid early payment discount: %s", v)
	}
	
	if discount < MinDiscountRate || discount > MaxDiscountRate {
		return fmt.Errorf("early payment discount must be between %f and %f", MinDiscountRate, MaxDiscountRate)
	}
	
	return nil
}

func validateLatePaymentPenalty(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	penalty, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return fmt.Errorf("invalid late payment penalty: %s", v)
	}
	
	if penalty < 0 || penalty > 1.0 {
		return fmt.Errorf("late payment penalty must be between 0 and 1.0")
	}
	
	return nil
}

func validatePatriotismDiscountRate(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	rate, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return fmt.Errorf("invalid patriotism discount rate: %s", v)
	}
	
	if rate < MinDiscountRate || rate > MaxDiscountRate {
		return fmt.Errorf("patriotism discount rate must be between %f and %f", MinDiscountRate, MaxDiscountRate)
	}
	
	return nil
}

func validateCulturalBonusRate(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	rate, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return fmt.Errorf("invalid cultural bonus rate: %s", v)
	}
	
	if rate < MinDiscountRate || rate > MaxDiscountRate {
		return fmt.Errorf("cultural bonus rate must be between %f and %f", MinDiscountRate, MaxDiscountRate)
	}
	
	return nil
}

func validateTaxCalculationPrecision(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	if v < 0 || v > 18 {
		return fmt.Errorf("tax calculation precision must be between 0 and 18")
	}
	
	return nil
}

func validateVolumeThresholds(i interface{}) error {
	v, ok := i.([]*VolumeThreshold)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	if len(v) == 0 {
		return fmt.Errorf("volume thresholds cannot be empty")
	}
	
	for _, threshold := range v {
		if threshold.TransactionCount == 0 {
			return fmt.Errorf("volume threshold transaction count cannot be zero")
		}
		
		rate, err := strconv.ParseFloat(threshold.TaxRate, 64)
		if err != nil {
			return fmt.Errorf("invalid volume threshold tax rate: %s", threshold.TaxRate)
		}
		
		if rate < MinTaxRateValue || rate > MaxTaxRateValue {
			return fmt.Errorf("volume threshold tax rate must be between %f and %f", MinTaxRateValue, MaxTaxRateValue)
		}
	}
	
	return nil
}

func validateTaxBrackets(i interface{}) error {
	v, ok := i.([]*TaxBracket)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	if len(v) == 0 {
		return fmt.Errorf("tax brackets cannot be empty")
	}
	
	for _, bracket := range v {
		if bracket.AmountLimit <= 0 {
			return fmt.Errorf("tax bracket amount limit must be positive")
		}
		
		rate, err := strconv.ParseFloat(bracket.TaxRate, 64)
		if err != nil {
			return fmt.Errorf("invalid tax bracket tax rate: %s", bracket.TaxRate)
		}
		
		if rate < MinTaxRateValue || rate > MaxTaxRateValue {
			return fmt.Errorf("tax bracket tax rate must be between %f and %f", MinTaxRateValue, MaxTaxRateValue)
		}
		
		if bracket.TaxCap != "" {
			cap, err := strconv.ParseFloat(bracket.TaxCap, 64)
			if err != nil {
				return fmt.Errorf("invalid tax bracket tax cap: %s", bracket.TaxCap)
			}
			
			if cap < 0 {
				return fmt.Errorf("tax bracket tax cap cannot be negative")
			}
		}
	}
	
	return nil
}

func validateSystemMode(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	validModes := []string{
		SystemModeNormal,
		SystemModeEmergency,
		SystemModeMaintenance,
		SystemModeUpgrade,
		SystemModeAudit,
		SystemModeCompliance,
		SystemModeOptimization,
		SystemModeEducation,
		SystemModeReporting,
		SystemModeForecasting,
	}
	
	for _, mode := range validModes {
		if v == mode {
			return nil
		}
	}
	
	return fmt.Errorf("invalid system mode: %s", v)
}

func validateBool(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

// GetBaseTaxRateAsDec returns the base tax rate as a decimal
func (p TaxParams) GetBaseTaxRateAsDec() sdk.Dec {
	rate, _ := sdk.NewDecFromStr(p.BaseTaxRate)
	return rate
}

// GetMaxTaxAmountAsCoin returns the max tax amount as a coin
func (p TaxParams) GetMaxTaxAmountAsCoin(denom string) sdk.Coin {
	amount, _ := strconv.ParseInt(p.MaxTaxAmount, 10, 64)
	return sdk.NewCoin(denom, math.NewInt(amount))
}

// GetMinTaxAmountAsCoin returns the min tax amount as a coin
func (p TaxParams) GetMinTaxAmountAsCoin(denom string) sdk.Coin {
	amount, _ := strconv.ParseInt(p.MinTaxAmount, 10, 64)
	return sdk.NewCoin(denom, math.NewInt(amount))
}

// GetMaxVolumeDiscountAsDec returns the max volume discount as a decimal
func (p TaxParams) GetMaxVolumeDiscountAsDec() sdk.Dec {
	discount, _ := sdk.NewDecFromStr(p.MaxVolumeDiscount)
	return discount
}

// GetEarlyPaymentDiscountAsDec returns the early payment discount as a decimal
func (p TaxParams) GetEarlyPaymentDiscountAsDec() sdk.Dec {
	discount, _ := sdk.NewDecFromStr(p.EarlyPaymentDiscount)
	return discount
}

// GetLatePaymentPenaltyAsDec returns the late payment penalty as a decimal
func (p TaxParams) GetLatePaymentPenaltyAsDec() sdk.Dec {
	penalty, _ := sdk.NewDecFromStr(p.LatePaymentPenalty)
	return penalty
}

// GetPatriotismDiscountRateAsDec returns the patriotism discount rate as a decimal
func (p TaxParams) GetPatriotismDiscountRateAsDec() sdk.Dec {
	rate, _ := sdk.NewDecFromStr(p.PatriotismDiscountRate)
	return rate
}

// GetCulturalBonusRateAsDec returns the cultural bonus rate as a decimal
func (p TaxParams) GetCulturalBonusRateAsDec() sdk.Dec {
	rate, _ := sdk.NewDecFromStr(p.CulturalBonusRate)
	return rate
}

// IsInEmergencyMode returns true if the system is in emergency mode
func (p TaxParams) IsInEmergencyMode() bool {
	return p.EmergencyMode || p.SystemMode == SystemModeEmergency
}

// IsInMaintenanceMode returns true if the system is in maintenance mode
func (p TaxParams) IsInMaintenanceMode() bool {
	return p.MaintenanceMode || p.SystemMode == SystemModeMaintenance
}

// IsInUpgradeMode returns true if the system is in upgrade mode
func (p TaxParams) IsInUpgradeMode() bool {
	return p.UpgradeMode || p.SystemMode == SystemModeUpgrade
}

// IsInAuditMode returns true if the system is in audit mode
func (p TaxParams) IsInAuditMode() bool {
	return p.AuditMode || p.SystemMode == SystemModeAudit
}

// IsInDebugMode returns true if the system is in debug mode
func (p TaxParams) IsInDebugMode() bool {
	return p.DebugMode
}

// IsInTestMode returns true if the system is in test mode
func (p TaxParams) IsInTestMode() bool {
	return p.TestMode
}

// IsInPerformanceMode returns true if the system is in performance mode
func (p TaxParams) IsInPerformanceMode() bool {
	return p.PerformanceMode
}

// IsOperational returns true if the system is operational (not in emergency or maintenance mode)
func (p TaxParams) IsOperational() bool {
	return !p.IsInEmergencyMode() && !p.IsInMaintenanceMode() && !p.IsInUpgradeMode()
}