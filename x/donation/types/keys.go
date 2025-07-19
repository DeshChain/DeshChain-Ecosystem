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

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "donation"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_donation"
)

var (
	// ParamsKey is the key for module parameters
	ParamsKey = collections.NewPrefix(0)

	// NGOWalletKey is the key for NGO wallets
	NGOWalletKey = collections.NewPrefix(1)

	// NGOWalletCountKey is the key for NGO wallet counter
	NGOWalletCountKey = collections.NewPrefix(2)

	// DonationRecordKey is the key for donation records
	DonationRecordKey = collections.NewPrefix(3)

	// DonationRecordCountKey is the key for donation record counter
	DonationRecordCountKey = collections.NewPrefix(4)

	// DistributionRecordKey is the key for distribution records
	DistributionRecordKey = collections.NewPrefix(5)

	// DistributionRecordCountKey is the key for distribution record counter
	DistributionRecordCountKey = collections.NewPrefix(6)

	// AuditReportKey is the key for audit reports
	AuditReportKey = collections.NewPrefix(7)

	// AuditReportCountKey is the key for audit report counter
	AuditReportCountKey = collections.NewPrefix(8)

	// BeneficiaryTestimonialKey is the key for beneficiary testimonials
	BeneficiaryTestimonialKey = collections.NewPrefix(9)

	// BeneficiaryTestimonialCountKey is the key for beneficiary testimonial counter
	BeneficiaryTestimonialCountKey = collections.NewPrefix(10)

	// NGOByAddressKey is the key for NGO lookup by address
	NGOByAddressKey = collections.NewPrefix(11)

	// DonationByDonorKey is the key for donations by donor
	DonationByDonorKey = collections.NewPrefix(12)

	// DonationByNGOKey is the key for donations by NGO
	DonationByNGOKey = collections.NewPrefix(13)

	// DistributionByNGOKey is the key for distributions by NGO
	DistributionByNGOKey = collections.NewPrefix(14)

	// AuditByNGOKey is the key for audits by NGO
	AuditByNGOKey = collections.NewPrefix(15)

	// TestimonialByNGOKey is the key for testimonials by NGO
	TestimonialByNGOKey = collections.NewPrefix(16)

	// CampaignKey is the key for campaigns
	CampaignKey = collections.NewPrefix(17)

	// CampaignCountKey is the key for campaign counter
	CampaignCountKey = collections.NewPrefix(18)

	// RecurringDonationKey is the key for recurring donations
	RecurringDonationKey = collections.NewPrefix(19)

	// RecurringDonationCountKey is the key for recurring donation counter
	RecurringDonationCountKey = collections.NewPrefix(20)

	// EmergencyPauseKey is the key for emergency pause status
	EmergencyPauseKey = collections.NewPrefix(21)

	// StatisticsKey is the key for statistics cache
	StatisticsKey = collections.NewPrefix(22)

	// FundFlowKey is the key for fund flow tracking
	FundFlowKey = collections.NewPrefix(23)

	// TransparencyScoreKey is the key for transparency score tracking
	TransparencyScoreKey = collections.NewPrefix(24)

	// VerificationQueueKey is the key for verification queue
	VerificationQueueKey = collections.NewPrefix(25)
)

// Default NGO Categories
const (
	CategoryArmyWelfare         = "army_welfare"
	CategoryWarRelief          = "war_relief"
	CategoryDisabledSoldiers   = "disabled_soldiers"
	CategoryBorderAreaSchools  = "border_area_schools"
	CategoryMartyrsChildren    = "martyrs_children"
	CategoryDisasterRelief     = "disaster_relief"
	CategoryEducation          = "education"
	CategoryHealthcare         = "healthcare"
	CategoryEnvironment        = "environment"
	CategoryWomenEmpowerment   = "women_empowerment"
	CategoryChildWelfare       = "child_welfare"
	CategoryElderCare          = "elder_care"
	CategoryAnimalWelfare      = "animal_welfare"
	CategoryRuralDevelopment   = "rural_development"
	CategoryUrbanDevelopment   = "urban_development"
	CategoryPovertyAlleviation = "poverty_alleviation"
	CategorySkillDevelopment   = "skill_development"
	CategoryCulturalPreservation = "cultural_preservation"
	CategorySportsPromotion    = "sports_promotion"
	CategoryTechnologyAccess   = "technology_access"
)

// Default NGO Wallet Addresses (Multi-signature)
const (
	// ArmyWelfareWalletAddress is the address for Army Welfare Fund
	ArmyWelfareWalletAddress = "desh1army000000000000000000000000000000000000"

	// WarReliefWalletAddress is the address for War Relief Fund
	WarReliefWalletAddress = "desh1war0000000000000000000000000000000000000"

	// DisabledSoldiersWalletAddress is the address for Disabled Soldiers Fund
	DisabledSoldiersWalletAddress = "desh1disabled0000000000000000000000000000000"

	// BorderAreaSchoolsWalletAddress is the address for Border Area Schools Fund
	BorderAreaSchoolsWalletAddress = "desh1border00000000000000000000000000000000000"

	// MartyrsChildrenWalletAddress is the address for Martyrs' Children Fund
	MartyrsChildrenWalletAddress = "desh1martyrs0000000000000000000000000000000000"

	// DisasterReliefWalletAddress is the address for Disaster Relief Fund
	DisasterReliefWalletAddress = "desh1disaster000000000000000000000000000000000"
)

// Distribution Categories
const (
	DistributionCategoryMedical     = "medical"
	DistributionCategoryEducation   = "education"
	DistributionCategoryFood        = "food"
	DistributionCategoryShelter     = "shelter"
	DistributionCategoryEmergency   = "emergency"
	DistributionCategoryInfrastructure = "infrastructure"
	DistributionCategoryEquipment   = "equipment"
	DistributionCategoryTraining    = "training"
	DistributionCategoryResearch    = "research"
	DistributionCategorySupport     = "support"
	DistributionCategoryRehabilitation = "rehabilitation"
	DistributionCategoryPrevention  = "prevention"
	DistributionCategoryAwareness   = "awareness"
	DistributionCategoryCapacity    = "capacity_building"
	DistributionCategoryMaintenance = "maintenance"
)

// Audit Types
const (
	AuditTypeFinancial   = "financial"
	AuditTypeOperational = "operational"
	AuditTypeCompliance  = "compliance"
	AuditTypeImpact      = "impact"
	AuditTypeGovernance  = "governance"
	AuditTypeTransparency = "transparency"
	AuditTypeEfficiency  = "efficiency"
	AuditTypeRisk        = "risk"
	AuditTypeInternal    = "internal"
	AuditTypeExternal    = "external"
	AuditTypeSpecial     = "special"
)

// Verification Status
const (
	VerificationStatusPending   = "pending"
	VerificationStatusApproved  = "approved"
	VerificationStatusRejected  = "rejected"
	VerificationStatusUnderReview = "under_review"
	VerificationStatusRevoked   = "revoked"
	VerificationStatusSuspended = "suspended"
	VerificationStatusExpired   = "expired"
)

// Audit Severity Levels
const (
	SeverityLow      = "low"
	SeverityMedium   = "medium"
	SeverityHigh     = "high"
	SeverityCritical = "critical"
)

// Audit Finding Status
const (
	FindingStatusOpen        = "open"
	FindingStatusInProgress  = "in_progress"
	FindingStatusResolved    = "resolved"
	FindingStatusClosed      = "closed"
	FindingStatusOverdue     = "overdue"
	FindingStatusEscalated   = "escalated"
)

// Multi-signature Requirements
const (
	// DefaultRequiredSignatures is the default number of required signatures
	DefaultRequiredSignatures = 5

	// MaxSigners is the maximum number of signers allowed
	MaxSigners = 9

	// MinSigners is the minimum number of signers required
	MinSigners = 3

	// DefaultSigners is the default number of signers
	DefaultSigners = 9
)

// Transparency Score Ranges
const (
	TransparencyScoreExcellent = 9  // 9-10
	TransparencyScoreGood      = 7  // 7-8
	TransparencyScoreAverage   = 5  // 5-6
	TransparencyScorePoor      = 3  // 3-4
	TransparencyScoreUnacceptable = 1  // 1-2
)

// Impact Metric Types
const (
	MetricTypeBeneficiaries    = "beneficiaries_served"
	MetricTypeFundsUtilized    = "funds_utilized"
	MetricTypeProjectsCompleted = "projects_completed"
	MetricTypeGeographicReach  = "geographic_reach"
	MetricTypeVolunteers       = "volunteers_engaged"
	MetricTypePartnerships     = "partnerships_formed"
	MetricTypeTrainingsSessions = "training_sessions"
	MetricTypeInfrastructureBuilt = "infrastructure_built"
	MetricTypeJobs             = "jobs_created"
	MetricTypeScholarships     = "scholarships_provided"
	MetricTypeHealthcareServices = "healthcare_services"
	MetricTypeEnvironmentalImpact = "environmental_impact"
	MetricTypeTechnologyAccess = "technology_access"
	MetricTypeCommunityDevelopment = "community_development"
	MetricTypeCapacityBuilding = "capacity_building"
)

// Default values
const (
	DefaultMinDonationAmount = 1000000  // 1 NAMO (with 6 decimals)
	DefaultMaxDonationAmount = 10000000000000  // 10 Million NAMO
	DefaultAuditFrequency    = 12       // 12 months
	DefaultTransparencyScore = 5        // 5 out of 10
	DefaultTaxBenefitPercentage = "50"  // 50% tax benefit
	DefaultDonationFeePercentage = "0"  // No fee for donations
	DefaultReceiptGenerationFee = 100000 // 0.1 NAMO
	DefaultDistributionApprovalThreshold = 3 // 3 approvals needed
)

// Government Signers (example addresses)
var DefaultGovernmentSigners = []string{
	"desh1gov1000000000000000000000000000000000000",
	"desh1gov2000000000000000000000000000000000000",
	"desh1gov3000000000000000000000000000000000000",
	"desh1gov4000000000000000000000000000000000000",
	"desh1gov5000000000000000000000000000000000000",
	"desh1gov6000000000000000000000000000000000000",
	"desh1gov7000000000000000000000000000000000000",
	"desh1gov8000000000000000000000000000000000000",
	"desh1gov9000000000000000000000000000000000000",
}

// Default Verification Authorities
var DefaultVerificationAuthorities = []string{
	"desh1verifier1000000000000000000000000000000000",
	"desh1verifier2000000000000000000000000000000000",
	"desh1verifier3000000000000000000000000000000000",
}

// Default Audit Authorities
var DefaultAuditAuthorities = []string{
	"desh1auditor1000000000000000000000000000000000",
	"desh1auditor2000000000000000000000000000000000",
	"desh1auditor3000000000000000000000000000000000",
}

// Default Emergency Pause Authorities
var DefaultEmergencyPauseAuthorities = []string{
	"desh1emergency1000000000000000000000000000000000",
	"desh1emergency2000000000000000000000000000000000",
}

// Event types for donation module
const (
	EventTypeRegisterNGO              = "register_ngo"
	EventTypeVerifyNGO                = "verify_ngo"
	EventTypeUpdateNGO                = "update_ngo"
	EventTypeDonate                   = "donate"
	EventTypeDistributeFunds          = "distribute_funds"
	EventTypeSubmitAuditReport        = "submit_audit_report"
	EventTypeAddBeneficiaryTestimonial = "add_beneficiary_testimonial"
	EventTypeUpdateImpactMetrics      = "update_impact_metrics"
	EventTypeCreateCampaign           = "create_campaign"
	EventTypeEmergencyPause           = "emergency_pause"
	EventTypeUpdateParams             = "update_params"
	EventTypeMultiSigApproval         = "multi_sig_approval"
	EventTypeTransparencyUpdate       = "transparency_update"
	EventTypeFundFlowUpdate           = "fund_flow_update"
	EventTypeAuditScheduled           = "audit_scheduled"
	EventTypeVerificationExpired      = "verification_expired"
	EventTypeRecurringDonation        = "recurring_donation"
	EventTypeCampaignCompleted        = "campaign_completed"
	EventTypeReceiptGenerated         = "receipt_generated"
	EventTypeNFTReceiptCreated        = "nft_receipt_created"
	EventTypeImpactMeasured           = "impact_measured"
	EventTypeTestimonialVerified      = "testimonial_verified"
	EventTypeAuditReminder            = "audit_reminder"
	EventTypeEmergencyPauseLifted     = "emergency_pause_lifted"
)

// Event attribute keys
const (
	AttributeKeyNGOWalletID         = "ngo_wallet_id"
	AttributeKeyNGOName             = "ngo_name"
	AttributeKeyNGOCategory         = "ngo_category"
	AttributeKeyNGOAddress          = "ngo_address"
	AttributeKeyDonor               = "donor"
	AttributeKeyDonationID          = "donation_id"
	AttributeKeyDonationAmount      = "donation_amount"
	AttributeKeyDonationPurpose     = "donation_purpose"
	AttributeKeyDistributionID      = "distribution_id"
	AttributeKeyDistributionAmount  = "distribution_amount"
	AttributeKeyRecipient           = "recipient"
	AttributeKeyBeneficiaryName     = "beneficiary_name"
	AttributeKeyProjectName         = "project_name"
	AttributeKeyAuditReportID       = "audit_report_id"
	AttributeKeyAuditor             = "auditor"
	AttributeKeyAuditType           = "audit_type"
	AttributeKeyOverallRating       = "overall_rating"
	AttributeKeyTransparencyScore   = "transparency_score"
	AttributeKeyTestimonialID       = "testimonial_id"
	AttributeKeyTestimonialRating   = "testimonial_rating"
	AttributeKeyImpactMetric        = "impact_metric"
	AttributeKeyBeneficiaryCount    = "beneficiary_count"
	AttributeKeyCampaignID          = "campaign_id"
	AttributeKeyCampaignName        = "campaign_name"
	AttributeKeyTargetAmount        = "target_amount"
	AttributeKeyVerifier            = "verifier"
	AttributeKeyVerificationStatus  = "verification_status"
	AttributeKeyAuthority           = "authority"
	AttributeKeyPauseReason         = "pause_reason"
	AttributeKeyPauseDuration       = "pause_duration"
	AttributeKeyReceiptHash         = "receipt_hash"
	AttributeKeyNFTReceiptID        = "nft_receipt_id"
	AttributeKeyTaxBenefitAmount    = "tax_benefit_amount"
	AttributeKeyFundFlowType        = "fund_flow_type"
	AttributeKeyFromAddress         = "from_address"
	AttributeKeyToAddress           = "to_address"
	AttributeKeyTransactionHash     = "transaction_hash"
	AttributeKeyBlockHeight         = "block_height"
	AttributeKeyGPSCoordinates      = "gps_coordinates"
	AttributeKeyRegion              = "region"
	AttributeKeyDocumentationHash   = "documentation_hash"
	AttributeKeyPhotosHash          = "photos_hash"
	AttributeKeyVideoHash           = "video_hash"
	AttributeKeyApprovalCount       = "approval_count"
	AttributeKeyRequiredApprovals   = "required_approvals"
	AttributeKeySignerAddress       = "signer_address"
	AttributeKeyMultiSigAddress     = "multi_sig_address"
	AttributeKeyEmergencyLevel      = "emergency_level"
	AttributeKeyAuditDueDate        = "audit_due_date"
	AttributeKeyComplianceScore     = "compliance_score"
	AttributeKeyEfficiencyScore     = "efficiency_score"
	AttributeKeyImpactScore         = "impact_score"
	AttributeKeyUtilizationRate     = "utilization_rate"
	AttributeKeyDonationCount       = "donation_count"
	AttributeKeyDistributionCount   = "distribution_count"
	AttributeKeyAverageTransactionSize = "average_transaction_size"
	AttributeKeyMonthlyGrowthRate   = "monthly_growth_rate"
	AttributeKeyRegionalImpact      = "regional_impact"
	AttributeKeyPartnershipCount    = "partnership_count"
	AttributeKeyVolunteerCount      = "volunteer_count"
	AttributeKeyProjectCount        = "project_count"
	AttributeKeySuccessRate         = "success_rate"
	AttributeKeyResponseTime        = "response_time"
	AttributeKeyBeneficiarySatisfaction = "beneficiary_satisfaction"
	AttributeKeyDonorRetentionRate  = "donor_retention_rate"
	AttributeKeyTransparencyRank    = "transparency_rank"
	AttributeKeyEfficiencyRank      = "efficiency_rank"
	AttributeKeyImpactRank          = "impact_rank"
	AttributeKeyOverallRank         = "overall_rank"
	AttributeKeyCulturalQuoteID     = "cultural_quote_id"
	AttributeKeyMatchingFunds       = "matching_funds"
	AttributeKeyRecurringDonationID = "recurring_donation_id"
	AttributeKeyIsAnonymous         = "is_anonymous"
	AttributeKeyIsPublic            = "is_public"
	AttributeKeyLanguage            = "language"
	AttributeKeyTranslation         = "translation"
	AttributeKeyVerificationMethod  = "verification_method"
	AttributeKeyDocumentType        = "document_type"
	AttributeKeyExpirationDate      = "expiration_date"
	AttributeKeyRenewalDate         = "renewal_date"
	AttributeKeyMaintenanceMode     = "maintenance_mode"
	AttributeKeySystemStatus        = "system_status"
	AttributeKeyDataIntegrity       = "data_integrity"
	AttributeKeyBackupStatus        = "backup_status"
	AttributeKeySecurityLevel       = "security_level"
	AttributeKeyAccessLevel         = "access_level"
	AttributeKeyPermissionLevel     = "permission_level"
	AttributeKeyNotificationSent    = "notification_sent"
	AttributeKeyEmailSent           = "email_sent"
	AttributeKeySMSSent             = "sms_sent"
	AttributeKeyPushNotificationSent = "push_notification_sent"
	AttributeKeyAlertLevel          = "alert_level"
	AttributeKeyEscalationLevel     = "escalation_level"
	AttributeKeyFollowUpRequired    = "follow_up_required"
	AttributeKeyFollowUpDate        = "follow_up_date"
	AttributeKeyCompletionDate      = "completion_date"
	AttributeKeyDeadlineDate        = "deadline_date"
	AttributeKeyPriorityLevel       = "priority_level"
	AttributeKeyUrgencyLevel        = "urgency_level"
	AttributeKeyStakeholderNotified = "stakeholder_notified"
	AttributeKeyPublicNotified      = "public_notified"
	AttributeKeyMediaNotified       = "media_notified"
	AttributeKeyGovernmentNotified  = "government_notified"
	AttributeKeyDonorNotified       = "donor_notified"
	AttributeKeyBeneficiaryNotified = "beneficiary_notified"
	AttributeKeyVolunteerNotified   = "volunteer_notified"
	AttributeKeyPartnerNotified     = "partner_notified"
	AttributeKeyAuditorNotified     = "auditor_notified"
	AttributeKeyVerifierNotified    = "verifier_notified"
	AttributeKeyAuthorityNotified   = "authority_notified"
	AttributeKeyEmergencyContactNotified = "emergency_contact_notified"
	AttributeKeyLegalContactNotified = "legal_contact_notified"
	AttributeKeyTechnicalContactNotified = "technical_contact_notified"
	AttributeKeyMediaContactNotified = "media_contact_notified"
	AttributeKeyPublicRelationsNotified = "public_relations_notified"
	AttributeKeyComplianceOfficerNotified = "compliance_officer_notified"
	AttributeKeyRiskManagerNotified = "risk_manager_notified"
	AttributeKeyQualityAssuranceNotified = "quality_assurance_notified"
	AttributeKeySecurityOfficerNotified = "security_officer_notified"
	AttributeKeyDataProtectionOfficerNotified = "data_protection_officer_notified"
)

// Event attribute values
const (
	AttributeValueCategory = ModuleName
)