package types

import (
	"cosmossdk.io/collections"
)

const (
	// ModuleName defines the module name
	ModuleName = "grampension"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for gram pension
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_grampension"
)

// KVStore keys
var (
	// ParamsKey is the key for parameters
	ParamsKey = collections.NewPrefix(0)

	// PensionSchemePrefix is the prefix for pension schemes
	PensionSchemePrefix = collections.NewPrefix(1)

	// ParticipantPrefix is the prefix for participants
	ParticipantPrefix = collections.NewPrefix(2)

	// ContributionPrefix is the prefix for contributions
	ContributionPrefix = collections.NewPrefix(3)

	// PaymentRecordPrefix is the prefix for payment records
	PaymentRecordPrefix = collections.NewPrefix(4)

	// WithdrawalRequestPrefix is the prefix for withdrawal requests
	WithdrawalRequestPrefix = collections.NewPrefix(5)

	// ReferralRewardPrefix is the prefix for referral rewards
	ReferralRewardPrefix = collections.NewPrefix(6)

	// KYCStatusPrefix is the prefix for KYC status
	KYCStatusPrefix = collections.NewPrefix(7)

	// PensionStatsPrefix is the prefix for pension statistics
	PensionStatsPrefix = collections.NewPrefix(8)

	// LoyaltyProgramPrefix is the prefix for loyalty program data
	LoyaltyProgramPrefix = collections.NewPrefix(9)

	// PerformanceMetricsPrefix is the prefix for performance metrics
	PerformanceMetricsPrefix = collections.NewPrefix(10)

	// MaturityProjectionPrefix is the prefix for maturity projections
	MaturityProjectionPrefix = collections.NewPrefix(11)

	// SustainabilityReportPrefix is the prefix for sustainability reports
	SustainabilityReportPrefix = collections.NewPrefix(12)

	// FundUtilizationPrefix is the prefix for fund utilization data
	FundUtilizationPrefix = collections.NewPrefix(13)

	// RiskAssessmentPrefix is the prefix for risk assessments
	RiskAssessmentPrefix = collections.NewPrefix(14)

	// CulturalEngagementPrefix is the prefix for cultural engagement data
	CulturalEngagementPrefix = collections.NewPrefix(15)

	// PatriotismScorePrefix is the prefix for patriotism scores
	PatriotismScorePrefix = collections.NewPrefix(16)

	// SchemeStatisticsPrefix is the prefix for scheme statistics
	SchemeStatisticsPrefix = collections.NewPrefix(17)

	// ParticipantIndexPrefix is the prefix for participant indexes
	ParticipantIndexPrefix = collections.NewPrefix(18)

	// ContributionIndexPrefix is the prefix for contribution indexes
	ContributionIndexPrefix = collections.NewPrefix(19)

	// MaturityIndexPrefix is the prefix for maturity indexes
	MaturityIndexPrefix = collections.NewPrefix(20)
)

// Secondary index prefixes
var (
	// ParticipantByAddressPrefix indexes participants by address
	ParticipantByAddressPrefix = collections.NewPrefix(100)

	// ParticipantBySchemePrefix indexes participants by scheme
	ParticipantBySchemePrefix = collections.NewPrefix(101)

	// ParticipantByStatusPrefix indexes participants by status
	ParticipantByStatusPrefix = collections.NewPrefix(102)

	// ParticipantByKYCStatusPrefix indexes participants by KYC status
	ParticipantByKYCStatusPrefix = collections.NewPrefix(103)

	// ParticipantByEnrollmentDatePrefix indexes participants by enrollment date
	ParticipantByEnrollmentDatePrefix = collections.NewPrefix(104)

	// ParticipantByMaturityDatePrefix indexes participants by maturity date
	ParticipantByMaturityDatePrefix = collections.NewPrefix(105)

	// ContributionByParticipantPrefix indexes contributions by participant
	ContributionByParticipantPrefix = collections.NewPrefix(106)

	// ContributionByDatePrefix indexes contributions by date
	ContributionByDatePrefix = collections.NewPrefix(107)

	// ContributionByStatusPrefix indexes contributions by status
	ContributionByStatusPrefix = collections.NewPrefix(108)

	// WithdrawalByParticipantPrefix indexes withdrawals by participant
	WithdrawalByParticipantPrefix = collections.NewPrefix(109)

	// WithdrawalByStatusPrefix indexes withdrawals by status
	WithdrawalByStatusPrefix = collections.NewPrefix(110)

	// WithdrawalByDatePrefix indexes withdrawals by date
	WithdrawalByDatePrefix = collections.NewPrefix(111)

	// ReferralByReferrerPrefix indexes referrals by referrer
	ReferralByReferrerPrefix = collections.NewPrefix(112)

	// ReferralByParticipantPrefix indexes referrals by participant
	ReferralByParticipantPrefix = collections.NewPrefix(113)

	// SchemeByStatusPrefix indexes schemes by status
	SchemeByStatusPrefix = collections.NewPrefix(114)

	// SchemeByCreatedDatePrefix indexes schemes by created date
	SchemeByCreatedDatePrefix = collections.NewPrefix(115)
)

// Event types
const (
	EventTypeSchemeCreated           = "scheme_created"
	EventTypeSchemeUpdated           = "scheme_updated"
	EventTypeSchemeActivated         = "scheme_activated"
	EventTypeSchemeDeactivated       = "scheme_deactivated"
	EventTypeParticipantEnrolled     = "participant_enrolled"
	EventTypeParticipantUpdated      = "participant_updated"
	EventTypeParticipantSuspended    = "participant_suspended"
	EventTypeParticipantReinstated   = "participant_reinstated"
	EventTypeParticipantMatured      = "participant_matured"
	EventTypeParticipantWithdrawn    = "participant_withdrawn"
	EventTypeContributionMade        = "contribution_made"
	EventTypeContributionProcessed   = "contribution_processed"
	EventTypeContributionFailed      = "contribution_failed"
	EventTypeContributionRefunded    = "contribution_refunded"
	EventTypePaymentOverdue          = "payment_overdue"
	EventTypePaymentPenalty          = "payment_penalty"
	EventTypePaymentBonus            = "payment_bonus"
	EventTypeMaturityPayout          = "maturity_payout"
	EventTypeMaturityCalculated      = "maturity_calculated"
	EventTypeMaturityProcessed       = "maturity_processed"
	EventTypeWithdrawalRequested     = "withdrawal_requested"
	EventTypeWithdrawalProcessed     = "withdrawal_processed"
	EventTypeWithdrawalCancelled     = "withdrawal_cancelled"
	EventTypeWithdrawalApproved      = "withdrawal_approved"
	EventTypeWithdrawalRejected      = "withdrawal_rejected"
	EventTypeReferralRewardProcessed = "referral_reward_processed"
	EventTypeReferralMilestone       = "referral_milestone"
	EventTypeKYCStatusUpdated        = "kyc_status_updated"
	EventTypeKYCDocumentSubmitted    = "kyc_document_submitted"
	EventTypeKYCVerificationComplete = "kyc_verification_complete"
	EventTypePerformanceUpdated     = "performance_updated"
	EventTypePerformanceMilestone   = "performance_milestone"
	EventTypeLoyaltyPointsEarned     = "loyalty_points_earned"
	EventTypeLoyaltyRewardClaimed    = "loyalty_reward_claimed"
	EventTypeLoyaltyTierUpgraded     = "loyalty_tier_upgraded"
	EventTypeCulturalEngagement      = "cultural_engagement"
	EventTypePatriotismScoreUpdated  = "patriotism_score_updated"
	EventTypeFundAllocation          = "fund_allocation"
	EventTypeFundUtilization         = "fund_utilization"
	EventTypeRiskAssessment          = "risk_assessment"
	EventTypeSustainabilityReport    = "sustainability_report"
	EventTypeSchemeAudit             = "scheme_audit"
	EventTypeComplianceCheck         = "compliance_check"
	EventTypeEmergencyAction         = "emergency_action"
	EventTypeSystemMaintenance       = "system_maintenance"
)

// Attribute keys
const (
	AttributeKeySchemeID             = "scheme_id"
	AttributeKeySchemeName           = "scheme_name"
	AttributeKeySchemeStatus         = "scheme_status"
	AttributeKeyParticipantID        = "participant_id"
	AttributeKeyParticipantAddress   = "participant_address"
	AttributeKeyParticipantStatus    = "participant_status"
	AttributeKeyContributionAmount   = "contribution_amount"
	AttributeKeyContributionMonth    = "contribution_month"
	AttributeKeyContributionStatus   = "contribution_status"
	AttributeKeyMaturityAmount       = "maturity_amount"
	AttributeKeyMaturityDate         = "maturity_date"
	AttributeKeyMaturityBonus        = "maturity_bonus"
	AttributeKeyWithdrawalAmount     = "withdrawal_amount"
	AttributeKeyWithdrawalPenalty    = "withdrawal_penalty"
	AttributeKeyWithdrawalStatus     = "withdrawal_status"
	AttributeKeyReferralReward       = "referral_reward"
	AttributeKeyReferralBonus        = "referral_bonus"
	AttributeKeyReferrerAddress      = "referrer_address"
	AttributeKeyReferredParticipant  = "referred_participant"
	AttributeKeyKYCStatus            = "kyc_status"
	AttributeKeyKYCLevel             = "kyc_level"
	AttributeKeyKYCDocuments         = "kyc_documents"
	AttributeKeyPerformanceScore     = "performance_score"
	AttributeKeyPerformanceRank      = "performance_rank"
	AttributeKeyPerformanceCategory  = "performance_category"
	AttributeKeyLoyaltyPoints        = "loyalty_points"
	AttributeKeyLoyaltyTier          = "loyalty_tier"
	AttributeKeyLoyaltyReward        = "loyalty_reward"
	AttributeKeyCulturalQuoteID      = "cultural_quote_id"
	AttributeKeyCulturalQuoteText    = "cultural_quote_text"
	AttributeKeyCulturalEngagement   = "cultural_engagement"
	AttributeKeyPatriotismScore      = "patriotism_score"
	AttributeKeyPatriotismRank       = "patriotism_rank"
	AttributeKeyTransactionHash      = "transaction_hash"
	AttributeKeyTransactionFee       = "transaction_fee"
	AttributeKeyReceiptHash          = "receipt_hash"
	AttributeKeyBlockHeight          = "block_height"
	AttributeKeyTimestamp            = "timestamp"
	AttributeKeyExchangeRate         = "exchange_rate"
	AttributeKeyPaymentMethod        = "payment_method"
	AttributeKeyPaymentStatus        = "payment_status"
	AttributeKeyOnTimeBonus          = "on_time_bonus"
	AttributeKeyLatePenalty          = "late_penalty"
	AttributeKeyGracePeriod          = "grace_period"
	AttributeKeyNotificationSent     = "notification_sent"
	AttributeKeyNotificationType     = "notification_type"
	AttributeKeyRiskScore            = "risk_score"
	AttributeKeyRiskLevel            = "risk_level"
	AttributeKeyRiskFactors          = "risk_factors"
	AttributeKeyFundAllocation       = "fund_allocation"
	AttributeKeyFundUtilization      = "fund_utilization"
	AttributeKeyFundBalance          = "fund_balance"
	AttributeKeySustainabilityScore  = "sustainability_score"
	AttributeKeySustainabilityStatus = "sustainability_status"
	AttributeKeyAuditResult          = "audit_result"
	AttributeKeyAuditScore           = "audit_score"
	AttributeKeyComplianceStatus     = "compliance_status"
	AttributeKeyComplianceLevel      = "compliance_level"
	AttributeKeyEmergencyType        = "emergency_type"
	AttributeKeyEmergencyAction      = "emergency_action"
	AttributeKeyMaintenanceType      = "maintenance_type"
	AttributeKeyMaintenanceDuration  = "maintenance_duration"
)

// Default values
const (
	DefaultMinAge                    = 18
	DefaultMaxAge                    = 65
	DefaultContributionPeriod        = 12
	DefaultGracePeriodDays           = 7
	DefaultEarlyWithdrawalPenalty    = "0.10" // 10%
	DefaultLatePaymentPenalty        = "0.02" // 2%
	DefaultReferralRewardPercentage  = "0.01" // 1%
	DefaultOnTimeBonusPercentage     = "0.005" // 0.5%
	DefaultPerformanceBonusThreshold = 80
	DefaultMaxParticipants           = 1000000
	DefaultSchemeAdminFee            = "0.001" // 0.1%
	DefaultMinimumContribution       = "1000000000" // 1000 NAMO (assuming 6 decimals)
	DefaultMaximumContribution       = "10000000000" // 10000 NAMO
	DefaultMaturityBonusPercentage   = "0.50" // 50%
	DefaultSustainabilityThreshold   = "0.80" // 80%
	DefaultRiskThreshold             = "0.70" // 70%
)

// Status values
const (
	StatusActive     = "active"
	StatusInactive   = "inactive"
	StatusSuspended  = "suspended"
	StatusMatured    = "matured"
	StatusWithdrawn  = "withdrawn"
	StatusDefaulted  = "defaulted"
	StatusPending    = "pending"
	StatusProcessing = "processing"
	StatusCompleted  = "completed"
	StatusFailed     = "failed"
	StatusCancelled  = "cancelled"
	StatusApproved   = "approved"
	StatusRejected   = "rejected"
)

// Performance categories
const (
	PerformanceCategoryExcellent = "excellent"
	PerformanceCategoryGood      = "good"
	PerformanceCategoryAverage   = "average"
	PerformanceCategoryPoor      = "poor"
	PerformanceCategoryCritical  = "critical"
)

// Loyalty tiers
const (
	LoyaltyTierBronze   = "bronze"
	LoyaltyTierSilver   = "silver"
	LoyaltyTierGold     = "gold"
	LoyaltyTierPlatinum = "platinum"
	LoyaltyTierDiamond  = "diamond"
)

// Risk levels
const (
	RiskLevelLow      = "low"
	RiskLevelMedium   = "medium"
	RiskLevelHigh     = "high"
	RiskLevelCritical = "critical"
)

// Notification types
const (
	NotificationTypeContribution    = "contribution"
	NotificationTypeMaturity        = "maturity"
	NotificationTypeWithdrawal      = "withdrawal"
	NotificationTypeKYC             = "kyc"
	NotificationTypePerformance     = "performance"
	NotificationTypeLoyalty         = "loyalty"
	NotificationTypeReferral        = "referral"
	NotificationTypeEmergency       = "emergency"
	NotificationTypeMaintenance     = "maintenance"
	NotificationTypePromotion       = "promotion"
	NotificationTypeEducational     = "educational"
	NotificationTypeCompliance      = "compliance"
	NotificationTypeAudit           = "audit"
	NotificationTypeSystemUpdate    = "system_update"
	NotificationTypeSecurityAlert   = "security_alert"
)

// Payment methods
const (
	PaymentMethodBankTransfer  = "bank_transfer"
	PaymentMethodUPI           = "upi"
	PaymentMethodWallet        = "wallet"
	PaymentMethodDebitCard     = "debit_card"
	PaymentMethodCreditCard    = "credit_card"
	PaymentMethodNetBanking    = "net_banking"
	PaymentMethodCrypto        = "crypto"
	PaymentMethodAutoDebit     = "auto_debit"
	PaymentMethodQR            = "qr"
	PaymentMethodDeshPay       = "deshpay"
)

// KYC levels
const (
	KYCLevelBasic      = "basic"
	KYCLevelIntermediate = "intermediate"
	KYCLevelAdvanced   = "advanced"
	KYCLevelPremium    = "premium"
)

// Compliance levels
const (
	ComplianceLevelMinimal = "minimal"
	ComplianceLevelStandard = "standard"
	ComplianceLevelEnhanced = "enhanced"
	ComplianceLevelPremium  = "premium"
)

// Audit frequencies
const (
	AuditFrequencyDaily    = "daily"
	AuditFrequencyWeekly   = "weekly"
	AuditFrequencyMonthly  = "monthly"
	AuditFrequencyQuarterly = "quarterly"
	AuditFrequencyYearly   = "yearly"
)

// Report periods
const (
	ReportPeriodDaily     = "daily"
	ReportPeriodWeekly    = "weekly"
	ReportPeriodMonthly   = "monthly"
	ReportPeriodQuarterly = "quarterly"
	ReportPeriodYearly    = "yearly"
)

// Cultural engagement types
const (
	CulturalEngagementQuote      = "quote"
	CulturalEngagementHistory    = "history"
	CulturalEngagementTradition  = "tradition"
	CulturalEngagementFestival   = "festival"
	CulturalEngagementLanguage   = "language"
	CulturalEngagementArt        = "art"
	CulturalEngagementMusic      = "music"
	CulturalEngagementLiterature = "literature"
	CulturalEngagementPhilosophy = "philosophy"
	CulturalEngagementReligion   = "religion"
)

// Patriotism score categories
const (
	PatriotismCategoryContribution = "contribution"
	PatriotismCategoryDonation     = "donation"
	PatriotismCategoryVolunteering = "volunteering"
	PatriotismCategoryEducation    = "education"
	PatriotismCategoryEnvironment  = "environment"
	PatriotismCategoryHealthcare   = "healthcare"
	PatriotismCategoryDisaster     = "disaster"
	PatriotismCategoryDefense      = "defense"
	PatriotismCategoryInnovation   = "innovation"
	PatriotismCategoryRural        = "rural"
)

// Emergency action types
const (
	EmergencyActionSuspension     = "suspension"
	EmergencyActionEvacuation     = "evacuation"
	EmergencyActionFreeze         = "freeze"
	EmergencyActionReverse        = "reverse"
	EmergencyActionNotification   = "notification"
	EmergencyActionMaintenance    = "maintenance"
	EmergencyActionBackup         = "backup"
	EmergencyActionRestore        = "restore"
	EmergencyActionQuarantine     = "quarantine"
	EmergencyActionEscalation     = "escalation"
)

// Maintenance types
const (
	MaintenanceTypeScheduled   = "scheduled"
	MaintenanceTypeEmergency   = "emergency"
	MaintenanceTypeUpgrade     = "upgrade"
	MaintenanceTypeSecurity    = "security"
	MaintenanceTypePerformance = "performance"
	MaintenanceTypeCompliance  = "compliance"
	MaintenanceTypeDatabase    = "database"
	MaintenanceTypeNetwork     = "network"
	MaintenanceTypeBlockchain  = "blockchain"
	MaintenanceTypeApplication = "application"
)

// Query limits
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
	MinNameLength        = 2
	MaxNameLength        = 100
	MinAddressLength     = 10
	MaxAddressLength     = 200
	MinPhoneLength       = 10
	MaxPhoneLength       = 15
	MinEmailLength       = 5
	MaxEmailLength       = 100
	MinPANLength         = 10
	MaxPANLength         = 10
	MinAadhaarLength     = 12
	MaxAadhaarLength     = 12
	MinAccountLength     = 10
	MaxAccountLength     = 20
	MinIFSCLength        = 11
	MaxIFSCLength        = 11
	MinPincodeLength     = 6
	MaxPincodeLength     = 6
)

// Calculation constants
const (
	PercentageBase       = 100
	BasisPointBase       = 10000
	DecimalPrecision     = 6
	CalculationPrecision = 18
)

// Feature flags
const (
	FeatureAutoRenewal      = "auto_renewal"
	FeatureNotifications    = "notifications"
	FeatureReferralRewards  = "referral_rewards"
	FeatureLoyaltyProgram   = "loyalty_program"
	FeatureCulturalContent  = "cultural_content"
	FeaturePerformanceBonus = "performance_bonus"
	FeatureRiskAssessment   = "risk_assessment"
	FeatureAdvancedAnalytics = "advanced_analytics"
	FeatureAuditTrail       = "audit_trail"
	FeatureEmergencyMode    = "emergency_mode"
)

// Configuration keys
const (
	ConfigKeyMinContribution     = "min_contribution"
	ConfigKeyMaxContribution     = "max_contribution"
	ConfigKeyContributionPeriod  = "contribution_period"
	ConfigKeyMaturityBonus       = "maturity_bonus"
	ConfigKeyGracePeriod         = "grace_period"
	ConfigKeyPenaltyRate         = "penalty_rate"
	ConfigKeyBonusRate           = "bonus_rate"
	ConfigKeyReferralRate        = "referral_rate"
	ConfigKeyMaxParticipants     = "max_participants"
	ConfigKeyMinAge              = "min_age"
	ConfigKeyMaxAge              = "max_age"
	ConfigKeyKYCRequired         = "kyc_required"
	ConfigKeyAutoRenewal         = "auto_renewal"
	ConfigKeyNotificationEnabled = "notification_enabled"
	ConfigKeyAuditFrequency      = "audit_frequency"
	ConfigKeyRiskThreshold       = "risk_threshold"
	ConfigKeySustainabilityRatio = "sustainability_ratio"
)

// Error codes
const (
	ErrorCodeInvalidScheme         = 1001
	ErrorCodeInvalidParticipant    = 1002
	ErrorCodeInvalidContribution   = 1003
	ErrorCodeInvalidPayment        = 1004
	ErrorCodeInvalidWithdrawal     = 1005
	ErrorCodeInvalidKYC            = 1006
	ErrorCodeInvalidPerformance    = 1007
	ErrorCodeInvalidLoyalty        = 1008
	ErrorCodeInvalidReferral       = 1009
	ErrorCodeInvalidRisk           = 1010
	ErrorCodeInvalidFund           = 1011
	ErrorCodeInvalidSustainability = 1012
	ErrorCodeInvalidAudit          = 1013
	ErrorCodeInvalidCompliance     = 1014
	ErrorCodeInvalidEmergency      = 1015
	ErrorCodeInvalidMaintenance    = 1016
	ErrorCodeInvalidConfiguration  = 1017
	ErrorCodeInvalidFeature        = 1018
	ErrorCodeInvalidPermission     = 1019
	ErrorCodeInvalidState          = 1020
)

// Success codes
const (
	SuccessCodeSchemeCreated      = 2001
	SuccessCodeParticipantEnrolled = 2002
	SuccessCodeContributionMade   = 2003
	SuccessCodeMaturityProcessed  = 2004
	SuccessCodeWithdrawalApproved = 2005
	SuccessCodeKYCVerified        = 2006
	SuccessCodePerformanceUpdated = 2007
	SuccessCodeLoyaltyEarned      = 2008
	SuccessCodeReferralRewarded   = 2009
	SuccessCodeRiskAssessed       = 2010
	SuccessCodeFundAllocated      = 2011
	SuccessCodeAuditCompleted     = 2012
)

// Version constants
const (
	ModuleVersion = "1.0.0"
	SchemaVersion = "1.0.0"
	APIVersion    = "v1"
)

// Network constants
const (
	MainnetChainID = "deshchain-1"
	TestnetChainID = "deshchain-testnet-1"
	DevnetChainID  = "deshchain-devnet-1"
)

// Gas constants
const (
	GasCreateScheme       = 100000
	GasEnrollParticipant  = 80000
	GasContribute         = 60000
	GasProcessMaturity    = 120000
	GasProcessWithdrawal  = 100000
	GasUpdateKYC          = 50000
	GasUpdatePerformance  = 40000
	GasProcessReferral    = 60000
	GasUpdateLoyalty      = 40000
	GasRiskAssessment     = 80000
	GasAuditCompliance    = 100000
	GasEmergencyAction    = 150000
	GasMaintenanceAction  = 80000
)