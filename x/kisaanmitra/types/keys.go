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
	"cosmossdk.io/collections"
)

const (
	// ModuleName defines the module name
	ModuleName = "kisaanmitra"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for kisaan mitra
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_kisaanmitra"
)

// KVStore keys
var (
	// ParamsKey is the key for parameters
	ParamsKey = collections.NewPrefix(0)

	// LoanSchemePrefix is the prefix for loan schemes
	LoanSchemePrefix = collections.NewPrefix(1)

	// BorrowerPrefix is the prefix for borrowers
	BorrowerPrefix = collections.NewPrefix(2)

	// LoanApplicationPrefix is the prefix for loan applications
	LoanApplicationPrefix = collections.NewPrefix(3)

	// LoanPrefix is the prefix for active loans
	LoanPrefix = collections.NewPrefix(4)

	// RepaymentPrefix is the prefix for repayments
	RepaymentPrefix = collections.NewPrefix(5)

	// CollateralPrefix is the prefix for collateral
	CollateralPrefix = collections.NewPrefix(6)

	// VillagePoolPrefix is the prefix for village lending pools
	VillagePoolPrefix = collections.NewPrefix(7)

	// CreditHistoryPrefix is the prefix for credit history
	CreditHistoryPrefix = collections.NewPrefix(8)

	// RiskAssessmentPrefix is the prefix for risk assessments
	RiskAssessmentPrefix = collections.NewPrefix(9)

	// LendingStatsPrefix is the prefix for lending statistics
	LendingStatsPrefix = collections.NewPrefix(10)

	// DefaultPrefix is the prefix for defaulted loans
	DefaultPrefix = collections.NewPrefix(11)

	// CropInsurancePrefix is the prefix for crop insurance
	CropInsurancePrefix = collections.NewPrefix(12)

	// WeatherDataPrefix is the prefix for weather data
	WeatherDataPrefix = collections.NewPrefix(13)

	// CropCyclePrefix is the prefix for crop cycles
	CropCyclePrefix = collections.NewPrefix(14)

	// MarketPricePrefix is the prefix for market prices
	MarketPricePrefix = collections.NewPrefix(15)

	// CommunityValidationPrefix is the prefix for community validations
	CommunityValidationPrefix = collections.NewPrefix(16)

	// LiquidityPoolPrefix is the prefix for liquidity pools
	LiquidityPoolPrefix = collections.NewPrefix(17)

	// InterestRatePrefix is the prefix for interest rates
	InterestRatePrefix = collections.NewPrefix(18)

	// LoanPerformancePrefix is the prefix for loan performance metrics
	LoanPerformancePrefix = collections.NewPrefix(19)

	// VillageCoordinatorPrefix is the prefix for village coordinators
	VillageCoordinatorPrefix = collections.NewPrefix(20)
)

// Secondary index prefixes
var (
	// BorrowerByAddressPrefix indexes borrowers by address
	BorrowerByAddressPrefix = collections.NewPrefix(100)

	// BorrowerByVillagePrefix indexes borrowers by village
	BorrowerByVillagePrefix = collections.NewPrefix(101)

	// BorrowerByCreditScorePrefix indexes borrowers by credit score
	BorrowerByCreditScorePrefix = collections.NewPrefix(102)

	// LoanByBorrowerPrefix indexes loans by borrower
	LoanByBorrowerPrefix = collections.NewPrefix(103)

	// LoanByStatusPrefix indexes loans by status
	LoanByStatusPrefix = collections.NewPrefix(104)

	// LoanByMaturityPrefix indexes loans by maturity date
	LoanByMaturityPrefix = collections.NewPrefix(105)

	// LoanByAmountPrefix indexes loans by amount
	LoanByAmountPrefix = collections.NewPrefix(106)

	// LoanByVillagePrefix indexes loans by village
	LoanByVillagePrefix = collections.NewPrefix(107)

	// RepaymentByLoanPrefix indexes repayments by loan
	RepaymentByLoanPrefix = collections.NewPrefix(108)

	// RepaymentByDatePrefix indexes repayments by date
	RepaymentByDatePrefix = collections.NewPrefix(109)

	// ApplicationByStatusPrefix indexes applications by status
	ApplicationByStatusPrefix = collections.NewPrefix(110)

	// ApplicationByDatePrefix indexes applications by date
	ApplicationByDatePrefix = collections.NewPrefix(111)

	// DefaultByDatePrefix indexes defaults by date
	DefaultByDatePrefix = collections.NewPrefix(112)

	// CropBySeasonPrefix indexes crops by season
	CropBySeasonPrefix = collections.NewPrefix(113)

	// VillageByPerformancePrefix indexes villages by performance
	VillageByPerformancePrefix = collections.NewPrefix(114)

	// LoanByInterestRatePrefix indexes loans by interest rate
	LoanByInterestRatePrefix = collections.NewPrefix(115)
)

// Event types
const (
	EventTypeLoanSchemeCreated      = "loan_scheme_created"
	EventTypeLoanSchemeUpdated      = "loan_scheme_updated"
	EventTypeBorrowerRegistered     = "borrower_registered"
	EventTypeBorrowerUpdated        = "borrower_updated"
	EventTypeLoanApplicationSubmitted = "loan_application_submitted"
	EventTypeLoanApplicationApproved  = "loan_application_approved"
	EventTypeLoanApplicationRejected  = "loan_application_rejected"
	EventTypeLoanDisbursed           = "loan_disbursed"
	EventTypeLoanRepayment           = "loan_repayment"
	EventTypeLoanMatured             = "loan_matured"
	EventTypeLoanDefaulted           = "loan_defaulted"
	EventTypeLoanRescheduled         = "loan_rescheduled"
	EventTypeCollateralDeposited     = "collateral_deposited"
	EventTypeCollateralReleased      = "collateral_released"
	EventTypeCollateralLiquidated    = "collateral_liquidated"
	EventTypeCropInsuranceClaimed    = "crop_insurance_claimed"
	EventTypeWeatherEventTriggered   = "weather_event_triggered"
	EventTypeCropCycleStarted        = "crop_cycle_started"
	EventTypeCropCycleCompleted      = "crop_cycle_completed"
	EventTypeCommunityValidation     = "community_validation"
	EventTypeRiskAssessmentUpdated   = "risk_assessment_updated"
	EventTypeLiquidityProvided       = "liquidity_provided"
	EventTypeLiquidityWithdrawn      = "liquidity_withdrawn"
	EventTypeInterestRateUpdated     = "interest_rate_updated"
	EventTypeVillagePoolCreated      = "village_pool_created"
	EventTypePerformanceMetricsUpdated = "performance_metrics_updated"
)

// Attribute keys
const (
	AttributeKeyLoanSchemeID      = "loan_scheme_id"
	AttributeKeyLoanSchemeName    = "loan_scheme_name"
	AttributeKeyBorrowerID        = "borrower_id"
	AttributeKeyBorrowerAddress   = "borrower_address"
	AttributeKeyLoanApplicationID = "loan_application_id"
	AttributeKeyLoanID            = "loan_id"
	AttributeKeyLoanAmount        = "loan_amount"
	AttributeKeyLoanTerm          = "loan_term"
	AttributeKeyInterestRate      = "interest_rate"
	AttributeKeyRepaymentAmount   = "repayment_amount"
	AttributeKeyDueDate           = "due_date"
	AttributeKeyMaturityDate      = "maturity_date"
	AttributeKeyLoanStatus        = "loan_status"
	AttributeKeyCollateralType    = "collateral_type"
	AttributeKeyCollateralValue   = "collateral_value"
	AttributeKeyCropType          = "crop_type"
	AttributeKeyCropSeason        = "crop_season"
	AttributeKeyVillageCode       = "village_code"
	AttributeKeyVillageName       = "village_name"
	AttributeKeyCoordinatorAddress = "coordinator_address"
	AttributeKeyCreditScore       = "credit_score"
	AttributeKeyRiskScore         = "risk_score"
	AttributeKeyValidationStatus  = "validation_status"
	AttributeKeyValidatorAddress  = "validator_address"
	AttributeKeyLiquidityProvider = "liquidity_provider"
	AttributeKeyLiquidityAmount   = "liquidity_amount"
	AttributeKeyInsuranceAmount   = "insurance_amount"
	AttributeKeyWeatherCondition  = "weather_condition"
	AttributeKeyMarketPrice       = "market_price"
	AttributeKeyPerformanceScore  = "performance_score"
	AttributeKeyDefaultReason     = "default_reason"
	AttributeKeyRescheduleReason  = "reschedule_reason"
)

// Status values
const (
	StatusActive      = "active"
	StatusInactive    = "inactive"
	StatusPending     = "pending"
	StatusApproved    = "approved"
	StatusRejected    = "rejected"
	StatusDisbursed   = "disbursed"
	StatusRepaid      = "repaid"
	StatusDefaulted   = "defaulted"
	StatusRescheduled = "rescheduled"
	StatusCancelled   = "cancelled"
	StatusMatured     = "matured"
	StatusPartial     = "partial"
	StatusOverdue     = "overdue"
	StatusUnderReview = "under_review"
	StatusVerified    = "verified"
	StatusLiquidated  = "liquidated"
)

// Loan types
const (
	LoanTypeCropInput      = "crop_input"     // Seeds, fertilizers, pesticides
	LoanTypeEquipment      = "equipment"      // Tractors, tools, machinery
	LoanTypeEmergency      = "emergency"      // Medical, disaster relief
	LoanTypeHarvest        = "harvest"        // Post-harvest financing
	LoanTypeStorage        = "storage"        // Storage facility financing
	LoanTypeProcessing     = "processing"     // Food processing equipment
	LoanTypeLivestock      = "livestock"      // Animal husbandry
	LoanTypeOrganicFarming = "organic_farming" // Organic certification and inputs
	LoanTypeIrrigation     = "irrigation"     // Water management systems
	LoanTypeMarketing      = "marketing"      // Market linkage and transportation
)

// Crop seasons
const (
	SeasonKharif = "kharif"   // Monsoon season (June-October)
	SeasonRabi   = "rabi"     // Winter season (November-April)
	SeasonZaid   = "zaid"     // Summer season (April-June)
	SeasonPerennial = "perennial" // Year-round crops
)

// Risk levels
const (
	RiskLevelVeryLow  = "very_low"
	RiskLevelLow      = "low"
	RiskLevelMedium   = "medium"
	RiskLevelHigh     = "high"
	RiskLevelVeryHigh = "very_high"
	RiskLevelCritical = "critical"
)

// Validation types
const (
	ValidationType1Community   = "community"
	ValidationTypeCoordinator  = "coordinator"
	ValidationTypeAI          = "ai"
	ValidationTypeGovernment  = "government"
	ValidationTypeCreditBureau = "credit_bureau"
)

// Default values and limits
const (
	DefaultMinLoanAmount        = 1000          // ₹1,000
	DefaultMaxLoanAmount        = 5000000       // ₹50 lakhs
	DefaultMinInterestRate      = "0.06"        // 6%
	DefaultMaxInterestRate      = "0.18"        // 18%
	DefaultMinLoanTerm          = 3             // 3 months
	DefaultMaxLoanTerm          = 60            // 5 years
	DefaultMinCreditScore       = 300
	DefaultMaxCreditScore       = 900
	DefaultGracePeriodDays      = 15
	DefaultLatePaymentPenalty   = "0.02"        // 2%
	DefaultCollateralRatio      = "1.25"        // 125%
	DefaultInsurancePremium     = "0.03"        // 3%
	DefaultProcessingFee        = "0.01"        // 1%
	DefaultCommunityValidators  = 3
	DefaultLiquidityBuffer      = "0.20"        // 20%
	DefaultMaxExposurePerBorrower = "0.05"     // 5% of pool
	DefaultMaxExposurePerVillage  = "0.30"     // 30% of pool
)

// Configuration keys
const (
	ConfigKeyMinLoanAmount       = "min_loan_amount"
	ConfigKeyMaxLoanAmount       = "max_loan_amount"
	ConfigKeyMinInterestRate     = "min_interest_rate"
	ConfigKeyMaxInterestRate     = "max_interest_rate"
	ConfigKeyMinLoanTerm         = "min_loan_term"
	ConfigKeyMaxLoanTerm         = "max_loan_term"
	ConfigKeyGracePeriod         = "grace_period"
	ConfigKeyLatePaymentPenalty  = "late_payment_penalty"
	ConfigKeyCollateralRequired  = "collateral_required"
	ConfigKeyInsuranceRequired   = "insurance_required"
	ConfigKeyCommunityValidation = "community_validation"
	ConfigKeyRiskBasedPricing    = "risk_based_pricing"
	ConfigKeyWeatherIntegration  = "weather_integration"
	ConfigKeyMarketIntegration   = "market_integration"
)

// Error codes
const (
	ErrorCodeInvalidLoanScheme    = 2001
	ErrorCodeInvalidBorrower      = 2002
	ErrorCodeInvalidLoanAmount    = 2003
	ErrorCodeInvalidInterestRate  = 2004
	ErrorCodeInvalidLoanTerm      = 2005
	ErrorCodeInvalidCollateral    = 2006
	ErrorCodeInvalidCreditScore   = 2007
	ErrorCodeInsufficientLiquidity = 2008
	ErrorCodeLoanNotFound         = 2009
	ErrorCodeBorrowerNotFound     = 2010
	ErrorCodeUnauthorizedAccess   = 2011
	ErrorCodeValidationFailed     = 2012
	ErrorCodeRiskTooHigh          = 2013
	ErrorCodeLoanAlreadyExists    = 2014
	ErrorCodeRepaymentFailed      = 2015
	ErrorCodeCollateralInsufficient = 2016
	ErrorCodeWeatherClaim         = 2017
	ErrorCodeMarketVolatility     = 2018
	ErrorCodeCommunityRejection   = 2019
	ErrorCodeComplianceViolation  = 2020
)

// Success codes
const (
	SuccessCodeLoanApproved     = 3001
	SuccessCodeLoanDisbursed    = 3002
	SuccessCodeRepaymentSuccess = 3003
	SuccessCodeLoanMatured      = 3004
	SuccessCodeCollateralReleased = 3005
	SuccessCodeInsuranceClaimed = 3006
	SuccessCodeCommunityValidated = 3007
	SuccessCodeRiskAssessed     = 3008
)

// Time constants (in seconds)
const (
	SecondsInDay   = 86400
	SecondsInWeek  = 604800
	SecondsInMonth = 2592000  // 30 days
	SecondsInYear  = 31536000 // 365 days
)

// Validation constants
const (
	MinNameLength        = 2
	MaxNameLength        = 100
	MinPhoneLength       = 10
	MaxPhoneLength       = 15
	MinPincodeLength     = 6
	MaxPincodeLength     = 6
	MinAadhaarLength     = 12
	MaxAadhaarLength     = 12
	MinPANLength         = 10
	MaxPANLength         = 10
	MinBankAccountLength = 10
	MaxBankAccountLength = 20
)

// Percentage constants
const (
	PercentageBase   = 100
	BasisPointBase   = 10000
	DecimalPrecision = 6
)

// Version constants
const (
	ModuleVersion = "1.0.0"
	SchemaVersion = "1.0.0"
	APIVersion    = "v1"
)

// Gas constants
const (
	GasCreateLoanScheme       = 100000
	GasRegisterBorrower       = 80000
	GasSubmitApplication      = 60000
	GasApproveLoan           = 120000
	GasDisburseLoan          = 100000
	GasRepayLoan             = 80000
	GasLiquidateCollateral   = 150000
	GasUpdateRiskAssessment  = 60000
	GasCommunityValidation   = 40000
	GasWeatherClaim          = 100000
)

// Query limits
const (
	DefaultQueryLimit = 100
	MaxQueryLimit     = 1000
	DefaultPageSize   = 20
	MaxPageSize       = 100
)