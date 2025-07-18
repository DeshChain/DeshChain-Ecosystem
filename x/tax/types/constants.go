package types

// Tax rate constants
const (
	// Base tax rate limits
	MinTaxRateValue = 0.0
	MaxTaxRateValue = 1.0
	
	// Tax amount limits
	MinTaxAmount = 0
	MaxTaxAmount = 1000000 // ₹10 lakh equivalent
	
	// Discount rate limits
	MinDiscountRate = 0.0
	MaxDiscountRate = 1.0
	
	// Time constants
	SecondsInHour  = 3600
	SecondsInDay   = 86400
	SecondsInWeek  = 604800
	SecondsInMonth = 2592000 // 30 days
	SecondsInYear  = 31536000 // 365 days
)

// System modes
const (
	SystemModeNormal       = "normal"
	SystemModeEmergency    = "emergency"
	SystemModeMaintenance  = "maintenance"
	SystemModeUpgrade      = "upgrade"
	SystemModeAudit        = "audit"
	SystemModeCompliance   = "compliance"
	SystemModeOptimization = "optimization"
	SystemModeEducation    = "education"
	SystemModeReporting    = "reporting"
	SystemModeForecasting  = "forecasting"
)

// Volume thresholds for tax discounts
const (
	VolumeThreshold1K    = 1000
	VolumeThreshold10K   = 10000
	VolumeThreshold50K   = 50000
	VolumeThreshold100K  = 100000
	VolumeThreshold500K  = 500000
	VolumeThreshold1M    = 1000000
	VolumeThreshold10M   = 10000000
)

// Volume discount rates (as percentages)
const (
	VolumeDiscount1K   = "0.024"  // 2.4% for 1K+ transactions
	VolumeDiscount10K  = "0.023"  // 2.3% for 10K+ transactions
	VolumeDiscount50K  = "0.022"  // 2.2% for 50K+ transactions
	VolumeDiscount100K = "0.021"  // 2.1% for 100K+ transactions
	VolumeDiscount500K = "0.020"  // 2.0% for 500K+ transactions
	VolumeDiscount1M   = "0.019"  // 1.9% for 1M+ transactions
	VolumeDiscount10M  = "0.018"  // 1.8% for 10M+ transactions
)

// Tax bracket limits (in rupees equivalent)
const (
	TaxBracket1Limit = 40000    // ₹40,000
	TaxBracket2Limit = 400000   // ₹4,00,000
	TaxBracket3Limit = 10000000 // ₹1,00,00,000 (1 crore)
)

// Tax bracket rates and caps
const (
	TaxBracket1Rate = "0.025" // 2.5% for first bracket
	TaxBracket2Cap  = "1000"  // ₹1,000 cap for second bracket
	TaxBracket3Cap  = "1000"  // ₹1,000 flat for third bracket
)

// Event types for tax module
const (
	EventTypeTaxCollected     = "tax_collected"
	EventTypeTaxDistributed   = "tax_distributed"
	EventTypeVolumeDiscount   = "volume_discount_applied"
	EventTypePatriotismBonus  = "patriotism_bonus_applied"
	EventTypeCulturalBonus    = "cultural_bonus_applied"
	EventTypeDonationExemption = "donation_exemption_applied"
	EventTypeTokenBurn        = "token_burn"
	EventTypeFounderRoyalty   = "founder_royalty_distributed"
	EventTypeNGODonation      = "ngo_donation_distributed"
	EventTypeRevenueSharing   = "revenue_sharing_executed"
)

// Attribute keys for events
const (
	AttributeKeyTaxAmount        = "tax_amount"
	AttributeKeyTransactionAmount = "transaction_amount"
	AttributeKeyTaxRate          = "tax_rate"
	AttributeKeyDiscountRate     = "discount_rate"
	AttributeKeyFinalTaxRate     = "final_tax_rate"
	AttributeKeyRecipient        = "recipient"
	AttributeKeyPayer            = "payer"
	AttributeKeyPoolName         = "pool_name"
	AttributeKeyDistribution     = "distribution"
	AttributeKeyBurnAmount       = "burn_amount"
	AttributeKeyRoyaltyAmount    = "royalty_amount"
	AttributeKeyDonationAmount   = "donation_amount"
	AttributeKeyRevenueType      = "revenue_type"
)

// Module errors
const (
	ErrInvalidTaxRate        = "invalid tax rate"
	ErrInvalidTaxAmount      = "invalid tax amount"
	ErrInsufficientFunds     = "insufficient funds for tax"
	ErrInvalidDistribution   = "invalid tax distribution"
	ErrInvalidBeneficiary    = "invalid beneficiary address"
	ErrSystemNotOperational  = "system not operational"
	ErrTaxCalculationFailed  = "tax calculation failed"
	ErrDistributionFailed    = "tax distribution failed"
)