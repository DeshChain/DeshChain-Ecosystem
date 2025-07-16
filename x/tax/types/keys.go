package types

import (
	"cosmossdk.io/collections"
)

const (
	// ModuleName defines the module name
	ModuleName = "tax"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for tax
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_tax"
)

// KVStore keys
var (
	// ParamsKey is the key for parameters
	ParamsKey = collections.NewPrefix(0)

	// TaxConfigPrefix is the prefix for tax configuration
	TaxConfigPrefix = collections.NewPrefix(1)

	// UserTaxProfilePrefix is the prefix for user tax profiles
	UserTaxProfilePrefix = collections.NewPrefix(2)

	// TaxTransactionPrefix is the prefix for tax transactions
	TaxTransactionPrefix = collections.NewPrefix(3)

	// TaxIncentivePrefix is the prefix for tax incentives
	TaxIncentivePrefix = collections.NewPrefix(4)

	// TaxHolidayPrefix is the prefix for tax holidays
	TaxHolidayPrefix = collections.NewPrefix(5)

	// TaxDiscountPrefix is the prefix for tax discounts
	TaxDiscountPrefix = collections.NewPrefix(6)

	// TaxRefundPrefix is the prefix for tax refunds
	TaxRefundPrefix = collections.NewPrefix(7)

	// TaxStatisticsPrefix is the prefix for tax statistics
	TaxStatisticsPrefix = collections.NewPrefix(8)

	// TaxOptimizationPrefix is the prefix for tax optimization records
	TaxOptimizationPrefix = collections.NewPrefix(9)

	// TaxCompliancePrefix is the prefix for tax compliance records
	TaxCompliancePrefix = collections.NewPrefix(10)

	// TaxAuditPrefix is the prefix for tax audit records
	TaxAuditPrefix = collections.NewPrefix(11)

	// TaxForecastingPrefix is the prefix for tax forecasting data
	TaxForecastingPrefix = collections.NewPrefix(12)

	// TaxReportPrefix is the prefix for tax reports
	TaxReportPrefix = collections.NewPrefix(13)

	// TaxEducationPrefix is the prefix for tax education content
	TaxEducationPrefix = collections.NewPrefix(14)

	// TaxNewsPrefix is the prefix for tax news and updates
	TaxNewsPrefix = collections.NewPrefix(15)
)

// Secondary index prefixes
var (
	// TaxTransactionByUserPrefix indexes tax transactions by user
	TaxTransactionByUserPrefix = collections.NewPrefix(100)

	// TaxTransactionByDatePrefix indexes tax transactions by date
	TaxTransactionByDatePrefix = collections.NewPrefix(101)

	// TaxTransactionByAmountPrefix indexes tax transactions by amount
	TaxTransactionByAmountPrefix = collections.NewPrefix(102)

	// TaxIncentiveByTypePrefix indexes tax incentives by type
	TaxIncentiveByTypePrefix = collections.NewPrefix(103)

	// TaxIncentiveByDatePrefix indexes tax incentives by date
	TaxIncentiveByDatePrefix = collections.NewPrefix(104)

	// TaxHolidayByDatePrefix indexes tax holidays by date
	TaxHolidayByDatePrefix = collections.NewPrefix(105)

	// TaxDiscountByUserPrefix indexes tax discounts by user
	TaxDiscountByUserPrefix = collections.NewPrefix(106)

	// TaxRefundByUserPrefix indexes tax refunds by user
	TaxRefundByUserPrefix = collections.NewPrefix(107)

	// TaxRefundByStatusPrefix indexes tax refunds by status
	TaxRefundByStatusPrefix = collections.NewPrefix(108)
)

// Default tax configuration values
const (
	DefaultBaseTaxRate             = "0.025"    // 2.5%
	DefaultMaxTaxAmountINR         = "1000"     // ₹1,000
	DefaultMinTaxAmountINR         = "0"        // ₹0
	DefaultVolumeDiscountThreshold = "1000000"  // 1M daily transactions
	DefaultMaxVolumeDiscount       = "0.9"      // 90% max discount
	DefaultEarlyPaymentDiscount    = "0.01"     // 1% early payment discount
	DefaultLatePaymentPenalty      = "0.02"     // 2% late payment penalty
	DefaultTaxCapResetPeriod       = "86400"    // 24 hours (1 day)
	DefaultGracePeriodDays         = "7"        // 7 days grace period
	DefaultPatriotismDiscountRate  = "0.005"    // 0.5% per 100 patriotism score
	DefaultCulturalBonusRate       = "0.002"    // 0.2% cultural engagement bonus
	DefaultDonationTaxExemption    = "1.0"      // 100% exemption for donations
	DefaultOptimizationEnabled     = true
	DefaultProgressiveTaxEnabled   = true
	DefaultComplianceRequired      = true
	DefaultAuditTrailEnabled       = true
	DefaultEducationEnabled        = true
	DefaultForecastingEnabled      = true
	DefaultReportingEnabled        = true
	DefaultTransparencyEnabled     = true
)

// Tax calculation constants
const (
	TaxCalculationPrecision = 6
	MinTaxableAmount       = 1000000 // 1 NAMO (assuming 6 decimals)
	MaxTaxableAmount       = 1000000000000 // 1M NAMO
	TaxRateDecimalPlaces   = 6
	PercentageMultiplier   = 100
	BasisPointMultiplier   = 10000
)

// Volume-based tax reduction thresholds
const (
	VolumeThreshold1K     = 1000
	VolumeThreshold10K    = 10000
	VolumeThreshold50K    = 50000
	VolumeThreshold100K   = 100000
	VolumeThreshold500K   = 500000
	VolumeThreshold1M     = 1000000
	VolumeThreshold10M    = 10000000

	VolumeDiscount1K      = "0.025"  // 2.5%
	VolumeDiscount10K     = "0.0225" // 2.25%
	VolumeDiscount50K     = "0.020"  // 2.0%
	VolumeDiscount100K    = "0.015"  // 1.5%
	VolumeDiscount500K    = "0.010"  // 1.0%
	VolumeDiscount1M      = "0.005"  // 0.5%
	VolumeDiscount10M     = "0.0025" // 0.25%
)

// Progressive tax brackets (INR equivalent)
const (
	TaxBracket1Limit   = 40000     // ₹40,000
	TaxBracket2Limit   = 400000    // ₹4,00,000
	TaxBracket3Limit   = 4000000   // ₹40,00,000

	TaxBracket1Rate    = "0.025"   // 2.5%
	TaxBracket2Cap     = "1000"    // ₹1,000 cap
	TaxBracket3Cap     = "1000"    // ₹1,000 cap (flat)
)

// Tax period constants
const (
	TaxPeriodDaily     = "daily"
	TaxPeriodWeekly    = "weekly"
	TaxPeriodMonthly   = "monthly"
	TaxPeriodQuarterly = "quarterly"
	TaxPeriodYearly    = "yearly"
)

// Tax types
const (
	TaxTypeTransaction    = "transaction"
	TaxTypeTransfer      = "transfer"
	TaxTypeTrading       = "trading"
	TaxTypeStaking       = "staking"
	TaxTypeNFT           = "nft"
	TaxTypeDeFi          = "defi"
	TaxTypePrivacy       = "privacy"
	TaxTypeGovernance    = "governance"
	TaxTypeLaunchpad     = "launchpad"
	TaxTypeInsurance     = "insurance"
	TaxTypeOracle        = "oracle"
	TaxTypeBridge        = "bridge"
)

// Tax status values
const (
	TaxStatusPending    = "pending"
	TaxStatusCalculated = "calculated"
	TaxStatusPaid       = "paid"
	TaxStatusOverdue    = "overdue"
	TaxStatusWaived     = "waived"
	TaxStatusRefunded   = "refunded"
	TaxStatusDisputed   = "disputed"
	TaxStatusAudited    = "audited"
	TaxStatusCompliant  = "compliant"
	TaxStatusNonCompliant = "non_compliant"
)

// Tax optimization types
const (
	OptimizationTypeVolume      = "volume"
	OptimizationTypePatriotism  = "patriotism"
	OptimizationTypeCultural    = "cultural"
	OptimizationTypeDonation    = "donation"
	OptimizationTypeEarlyPayment = "early_payment"
	OptimizationTypeLoyalty     = "loyalty"
	OptimizationTypeStaking     = "staking"
	OptimizationTypeGovernance  = "governance"
	OptimizationTypeEducation   = "education"
	OptimizationTypeReferral    = "referral"
)

// Tax incentive types
const (
	IncentiveTypePatriotism     = "patriotism"
	IncentiveTypeCultural       = "cultural"
	IncentiveTypeDonation       = "donation"
	IncentiveTypeEducation      = "education"
	IncentiveTypeEnvironment    = "environment"
	IncentiveTypeHealthcare     = "healthcare"
	IncentiveTypeRural          = "rural"
	IncentiveTypeStartup        = "startup"
	IncentiveTypeInnovation     = "innovation"
	IncentiveTypeDefense        = "defense"
	IncentiveTypeDisaster       = "disaster"
	IncentiveTypeInfrastructure = "infrastructure"
	IncentiveTypeAgriculture    = "agriculture"
	IncentiveTypeManufacturing  = "manufacturing"
	IncentiveTypeServices       = "services"
	IncentiveTypeExport         = "export"
	IncentiveTypeImport         = "import"
	IncentiveTypeForeign        = "foreign"
	IncentiveTypeDigital        = "digital"
	IncentiveTypeGreen          = "green"
)

// Tax holiday types
const (
	HolidayTypeFestival       = "festival"
	HolidayTypeNational       = "national"
	HolidayTypeEmergency      = "emergency"
	HolidayTypeEconomic       = "economic"
	HolidayTypePromotion      = "promotion"
	HolidayTypeOnboarding     = "onboarding"
	HolidayTypeMilestone      = "milestone"
	HolidayTypeSpecial        = "special"
	HolidayTypeGovernment     = "government"
	HolidayTypeCompliance     = "compliance"
	HolidayTypeUpgrade        = "upgrade"
	HolidayTypeMaintenance    = "maintenance"
	HolidayTypePartnership    = "partnership"
	HolidayTypeCharity        = "charity"
	HolidayTypeEducational    = "educational"
)

// Tax discount types
const (
	DiscountTypeVolume        = "volume"
	DiscountTypePatriotism    = "patriotism"
	DiscountTypeCultural      = "cultural"
	DiscountTypeLoyalty       = "loyalty"
	DiscountTypeReferral      = "referral"
	DiscountTypeEarlyPayment  = "early_payment"
	DiscountTypeStaking       = "staking"
	DiscountTypeGovernance    = "governance"
	DiscountTypeEducation     = "education"
	DiscountTypeFirst         = "first_time"
	DiscountTypeBulk          = "bulk"
	DiscountTypePromotional   = "promotional"
	DiscountTypePartnership   = "partnership"
	DiscountTypeCharity       = "charity"
	DiscountTypeEmergency     = "emergency"
)

// Tax refund types
const (
	RefundTypeOverpayment   = "overpayment"
	RefundTypeError         = "error"
	RefundTypeDispute       = "dispute"
	RefundTypePolicy        = "policy"
	RefundTypeGovernment    = "government"
	RefundTypeEmergency     = "emergency"
	RefundTypeCompliance    = "compliance"
	RefundTypeAudit         = "audit"
	RefundTypeCorrection    = "correction"
	RefundTypeGoodwill      = "goodwill"
	RefundTypeCancellation  = "cancellation"
	RefundTypeFailure       = "failure"
	RefundTypeReversal      = "reversal"
	RefundTypeAdjustment    = "adjustment"
	RefundTypeCompensation  = "compensation"
)

// Tax compliance levels
const (
	ComplianceLevelBasic    = "basic"
	ComplianceLevelStandard = "standard"
	ComplianceLevelAdvanced = "advanced"
	ComplianceLevelPremium  = "premium"
	ComplianceLevelExpert   = "expert"
)

// Tax audit types
const (
	AuditTypeRoutine    = "routine"
	AuditTypeRandom     = "random"
	AuditTypeTargeted   = "targeted"
	AuditTypeCompliance = "compliance"
	AuditTypeRisk       = "risk"
	AuditTypeComplaint  = "complaint"
	AuditTypeGovernment = "government"
	AuditTypeInternal   = "internal"
	AuditTypeExternal   = "external"
	AuditTypeForensic   = "forensic"
)

// Tax report types
const (
	ReportTypeDaily        = "daily"
	ReportTypeWeekly       = "weekly"
	ReportTypeMonthly      = "monthly"
	ReportTypeQuarterly    = "quarterly"
	ReportTypeYearly       = "yearly"
	ReportTypeCustom       = "custom"
	ReportTypeCompliance   = "compliance"
	ReportTypeAudit        = "audit"
	ReportTypeStatistical  = "statistical"
	ReportTypeForecasting  = "forecasting"
	ReportTypeAnalytical   = "analytical"
	ReportTypeComparative  = "comparative"
	ReportTypeTrend        = "trend"
	ReportTypePerformance  = "performance"
	ReportTypeRisk         = "risk"
)

// Tax education types
const (
	EducationTypeBasic       = "basic"
	EducationTypeIntermediate = "intermediate"
	EducationTypeAdvanced    = "advanced"
	EducationTypeSpecialized = "specialized"
	EducationTypeCompliance  = "compliance"
	EducationTypeOptimization = "optimization"
	EducationTypePlanning    = "planning"
	EducationTypeStrategy    = "strategy"
	EducationTypeRegulatory  = "regulatory"
	EducationTypePolicy      = "policy"
	EducationTypeInternational = "international"
	EducationTypeDigital     = "digital"
	EducationTypeCryptocurrency = "cryptocurrency"
	EducationTypeBlockchain  = "blockchain"
	EducationTypeDeFi        = "defi"
)

// Event types
const (
	EventTypeTaxCalculated         = "tax_calculated"
	EventTypeTaxPaid              = "tax_paid"
	EventTypeTaxOptimized         = "tax_optimized"
	EventTypeTaxRefunded          = "tax_refunded"
	EventTypeTaxDiscountApplied   = "tax_discount_applied"
	EventTypeTaxIncentiveApplied  = "tax_incentive_applied"
	EventTypeTaxHolidayActivated  = "tax_holiday_activated"
	EventTypeTaxConfigUpdated     = "tax_config_updated"
	EventTypeTaxProfileUpdated    = "tax_profile_updated"
	EventTypeTaxComplianceChecked = "tax_compliance_checked"
	EventTypeTaxAuditCompleted    = "tax_audit_completed"
	EventTypeTaxReportGenerated   = "tax_report_generated"
	EventTypeTaxEducationCompleted = "tax_education_completed"
	EventTypeTaxForecastUpdated   = "tax_forecast_updated"
	EventTypeTaxStatisticsUpdated = "tax_statistics_updated"
	EventTypeTaxPenaltyApplied    = "tax_penalty_applied"
	EventTypeTaxGraceGranted      = "tax_grace_granted"
	EventTypeTaxCapReached        = "tax_cap_reached"
	EventTypeTaxExemptionGranted  = "tax_exemption_granted"
	EventTypeTaxThresholdChanged  = "tax_threshold_changed"
	EventTypeTaxRateChanged       = "tax_rate_changed"
)

// Attribute keys for events
const (
	AttributeKeyTaxAmount      = "tax_amount"
	AttributeKeyTaxRate        = "tax_rate"
	AttributeKeyTaxType        = "tax_type"
	AttributeKeyTaxStatus      = "tax_status"
	AttributeKeyTaxPeriod      = "tax_period"
	AttributeKeyTaxMethod      = "tax_method"
	AttributeKeyTaxOptimization = "tax_optimization"
	AttributeKeyTaxSavings     = "tax_savings"
	AttributeKeyTaxDiscount    = "tax_discount"
	AttributeKeyTaxIncentive   = "tax_incentive"
	AttributeKeyTaxHoliday     = "tax_holiday"
	AttributeKeyTaxRefund      = "tax_refund"
	AttributeKeyTaxPenalty     = "tax_penalty"
	AttributeKeyTaxGrace       = "tax_grace"
	AttributeKeyTaxCap         = "tax_cap"
	AttributeKeyTaxExemption   = "tax_exemption"
	AttributeKeyTransactionHash = "transaction_hash"
	AttributeKeyBlockHeight    = "block_height"
	AttributeKeyTimestamp      = "timestamp"
	AttributeKeyUserAddress    = "user_address"
	AttributeKeyRecipient      = "recipient"
	AttributeKeyVolumeData     = "volume_data"
	AttributeKeyPatriotismScore = "patriotism_score"
	AttributeKeyCulturalScore  = "cultural_score"
	AttributeKeyDonationFlag   = "donation_flag"
	AttributeKeyOptimizationType = "optimization_type"
	AttributeKeyComplianceLevel = "compliance_level"
	AttributeKeyAuditResult    = "audit_result"
	AttributeKeyReportType     = "report_type"
	AttributeKeyEducationType  = "education_type"
	AttributeKeyForecastPeriod = "forecast_period"
	AttributeKeyRiskLevel      = "risk_level"
	AttributeKeyConfidenceLevel = "confidence_level"
)

// Query limits and pagination
const (
	DefaultQueryLimit = 100
	MaxQueryLimit     = 1000
	DefaultPageSize   = 20
	MaxPageSize       = 100
)

// Time constants (in seconds)
const (
	SecondsInMinute = 60
	SecondsInHour   = 3600
	SecondsInDay    = 86400
	SecondsInWeek   = 604800
	SecondsInMonth  = 2592000  // 30 days
	SecondsInYear   = 31536000 // 365 days
)

// Validation constants
const (
	MinTaxRateValue = 0.0
	MaxTaxRateValue = 1.0
	MinTaxAmount    = 0
	MaxTaxAmount    = 1000000000000 // 1 trillion
	MinDiscountRate = 0.0
	MaxDiscountRate = 1.0
	MinIncentiveRate = 0.0
	MaxIncentiveRate = 1.0
	MinRefundAmount = 0
	MaxRefundAmount = 1000000000000 // 1 trillion
)

// Cultural and patriotism scoring
const (
	MaxPatriotismScore      = 10000
	MaxCulturalScore        = 10000
	PatriotismScoreMultiplier = 100
	CulturalScoreMultiplier = 100
	DonationScoreBonus      = 500
	VolunteerScoreBonus     = 300
	EducationScoreBonus     = 200
	EnvironmentScoreBonus   = 150
	HealthcareScoreBonus    = 100
)

// System configuration
const (
	SystemModeNormal     = "normal"
	SystemModeEmergency  = "emergency"
	SystemModeMaintenance = "maintenance"
	SystemModeUpgrade    = "upgrade"
	SystemModeAudit      = "audit"
	SystemModeCompliance = "compliance"
	SystemModeOptimization = "optimization"
	SystemModeEducation  = "education"
	SystemModeReporting  = "reporting"
	SystemModeForecasting = "forecasting"
)

// Error codes
const (
	ErrCodeInvalidTaxRate      = 3001
	ErrCodeInvalidTaxAmount    = 3002
	ErrCodeInvalidUser         = 3003
	ErrCodeInvalidTransaction  = 3004
	ErrCodeInvalidOptimization = 3005
	ErrCodeInvalidDiscount     = 3006
	ErrCodeInvalidIncentive    = 3007
	ErrCodeInvalidRefund       = 3008
	ErrCodeInvalidCompliance   = 3009
	ErrCodeInvalidAudit        = 3010
	ErrCodeInvalidReport       = 3011
	ErrCodeInvalidEducation    = 3012
	ErrCodeInvalidForecast     = 3013
	ErrCodeInvalidConfig       = 3014
	ErrCodeInvalidProfile      = 3015
	ErrCodeInvalidHoliday      = 3016
	ErrCodeInvalidStatistics   = 3017
	ErrCodeInvalidCalculation  = 3018
	ErrCodeInvalidPeriod       = 3019
	ErrCodeInvalidThreshold    = 3020
)

// Success codes
const (
	SuccessCodeTaxCalculated   = 4001
	SuccessCodeTaxPaid         = 4002
	SuccessCodeTaxOptimized    = 4003
	SuccessCodeTaxRefunded     = 4004
	SuccessCodeDiscountApplied = 4005
	SuccessCodeIncentiveApplied = 4006
	SuccessCodeHolidayActivated = 4007
	SuccessCodeConfigUpdated   = 4008
	SuccessCodeProfileUpdated  = 4009
	SuccessCodeComplianceChecked = 4010
	SuccessCodeAuditCompleted  = 4011
	SuccessCodeReportGenerated = 4012
)

// Gas costs for tax operations
const (
	GasTaxCalculation   = 50000
	GasTaxPayment      = 60000
	GasTaxOptimization = 80000
	GasTaxRefund       = 70000
	GasTaxDiscount     = 40000
	GasTaxIncentive    = 50000
	GasTaxHoliday      = 60000
	GasTaxConfig       = 100000
	GasTaxProfile      = 40000
	GasTaxCompliance   = 80000
	GasTaxAudit        = 100000
	GasTaxReport       = 60000
	GasTaxEducation    = 30000
	GasTaxForecast     = 70000
	GasTaxStatistics   = 50000
)

// Version information
const (
	TaxModuleVersion = "1.0.0"
	TaxAPIVersion    = "v1"
	TaxSchemaVersion = "1.0.0"
)

// Network-specific constants
const (
	MainnetTaxMultiplier = 1.0
	TestnetTaxMultiplier = 0.1
	DevnetTaxMultiplier  = 0.01
)