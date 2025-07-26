package types

import (
	"time"
)

// Analytics and monitoring types for DeshChain Identity System

// IdentityAnalytics contains comprehensive analytics data for the identity system
type IdentityAnalytics struct {
	// System overview
	TotalIdentities        uint64                  `json:"total_identities" yaml:"total_identities"`
	ActiveIdentities       uint64                  `json:"active_identities" yaml:"active_identities"`
	VerifiedIdentities     uint64                  `json:"verified_identities" yaml:"verified_identities"`
	
	// Verification statistics
	VerificationStats      *VerificationStatistics `json:"verification_stats" yaml:"verification_stats"`
	KYCStats              *KYCStatistics          `json:"kyc_stats" yaml:"kyc_stats"`
	BiometricStats        *BiometricStatistics    `json:"biometric_stats" yaml:"biometric_stats"`
	
	// Credential analytics
	CredentialStats       *CredentialStatistics   `json:"credential_stats" yaml:"credential_stats"`
	
	// Geographic and demographic analytics
	GeographicStats       *GeographicStatistics   `json:"geographic_stats" yaml:"geographic_stats"`
	DemographicStats      *DemographicStatistics  `json:"demographic_stats" yaml:"demographic_stats"`
	
	// Language and localization analytics
	LanguageStats         *LanguageStatistics     `json:"language_stats" yaml:"language_stats"`
	
	// Offline verification analytics
	OfflineStats          *OfflineStatistics      `json:"offline_stats" yaml:"offline_stats"`
	
	// Security and fraud analytics
	SecurityStats         *SecurityStatistics     `json:"security_stats" yaml:"security_stats"`
	FraudDetectionStats   *FraudDetectionStatistics `json:"fraud_detection_stats" yaml:"fraud_detection_stats"`
	
	// Performance metrics
	PerformanceMetrics    *PerformanceMetrics     `json:"performance_metrics" yaml:"performance_metrics"`
	
	// Usage patterns
	UsagePatterns         *UsagePatterns          `json:"usage_patterns" yaml:"usage_patterns"`
	
	// Growth and trends
	GrowthMetrics         *GrowthMetrics          `json:"growth_metrics" yaml:"growth_metrics"`
	
	// Compliance and audit metrics
	ComplianceMetrics     *ComplianceMetrics      `json:"compliance_metrics" yaml:"compliance_metrics"`
	
	// Timestamp information
	GeneratedAt           time.Time               `json:"generated_at" yaml:"generated_at"`
	TimeRange             *TimeRange              `json:"time_range" yaml:"time_range"`
	DataFreshness         time.Duration           `json:"data_freshness" yaml:"data_freshness"`
}

// VerificationStatistics contains verification-related analytics
type VerificationStatistics struct {
	TotalVerifications        uint64                    `json:"total_verifications" yaml:"total_verifications"`
	SuccessfulVerifications   uint64                    `json:"successful_verifications" yaml:"successful_verifications"`
	FailedVerifications       uint64                    `json:"failed_verifications" yaml:"failed_verifications"`
	SuccessRate               float64                   `json:"success_rate" yaml:"success_rate"`
	AverageVerificationTime   time.Duration             `json:"average_verification_time" yaml:"average_verification_time"`
	VerificationsByType       map[string]uint64         `json:"verifications_by_type" yaml:"verifications_by_type"`
	VerificationsByLevel      map[uint32]uint64         `json:"verifications_by_level" yaml:"verifications_by_level"`
	HourlyVerifications       map[int]uint64            `json:"hourly_verifications" yaml:"hourly_verifications"`
	DailyVerifications        map[string]uint64         `json:"daily_verifications" yaml:"daily_verifications"`
	PeakVerificationTimes     []PeakTime                `json:"peak_verification_times" yaml:"peak_verification_times"`
}

// KYCStatistics contains KYC-related analytics
type KYCStatistics struct {
	TotalKYCSubmissions       uint64                    `json:"total_kyc_submissions" yaml:"total_kyc_submissions"`
	PendingKYC                uint64                    `json:"pending_kyc" yaml:"pending_kyc"`
	ApprovedKYC               uint64                    `json:"approved_kyc" yaml:"approved_kyc"`
	RejectedKYC               uint64                    `json:"rejected_kyc" yaml:"rejected_kyc"`
	KYCApprovalRate           float64                   `json:"kyc_approval_rate" yaml:"kyc_approval_rate"`
	AverageKYCProcessingTime  time.Duration             `json:"average_kyc_processing_time" yaml:"average_kyc_processing_time"`
	KYCByLevel                map[uint32]uint64         `json:"kyc_by_level" yaml:"kyc_by_level"`
	KYCByMethod               map[string]uint64         `json:"kyc_by_method" yaml:"kyc_by_method"`
	DocumentTypes             map[string]uint64         `json:"document_types" yaml:"document_types"`
	IndiaStackIntegrations    uint64                    `json:"india_stack_integrations" yaml:"india_stack_integrations"`
	AadhaarVerifications      uint64                    `json:"aadhaar_verifications" yaml:"aadhaar_verifications"`
	DigiLockerConnections     uint64                    `json:"digilocker_connections" yaml:"digilocker_connections"`
	PanchayatKYCVerifications uint64                    `json:"panchayat_kyc_verifications" yaml:"panchayat_kyc_verifications"`
}

// BiometricStatistics contains biometric-related analytics
type BiometricStatistics struct {
	TotalBiometricEnrollments uint64                    `json:"total_biometric_enrollments" yaml:"total_biometric_enrollments"`
	BiometricsByType          map[BiometricType]uint64  `json:"biometrics_by_type" yaml:"biometrics_by_type"`
	BiometricVerifications    uint64                    `json:"biometric_verifications" yaml:"biometric_verifications"`
	BiometricSuccessRate      float64                   `json:"biometric_success_rate" yaml:"biometric_success_rate"`
	AverageMatchScore         float64                   `json:"average_match_score" yaml:"average_match_score"`
	BiometricQualityScores    *QualityScoreDistribution `json:"biometric_quality_scores" yaml:"biometric_quality_scores"`
	FalseAcceptanceRate       float64                   `json:"false_acceptance_rate" yaml:"false_acceptance_rate"`
	FalseRejectionRate        float64                   `json:"false_rejection_rate" yaml:"false_rejection_rate"`
	DeviceCompatibility       map[string]uint64         `json:"device_compatibility" yaml:"device_compatibility"`
}

// CredentialStatistics contains credential-related analytics
type CredentialStatistics struct {
	TotalCredentials          uint64                    `json:"total_credentials" yaml:"total_credentials"`
	ActiveCredentials         uint64                    `json:"active_credentials" yaml:"active_credentials"`
	ExpiredCredentials        uint64                    `json:"expired_credentials" yaml:"expired_credentials"`
	RevokedCredentials        uint64                    `json:"revoked_credentials" yaml:"revoked_credentials"`
	CredentialsByType         map[string]uint64         `json:"credentials_by_type" yaml:"credentials_by_type"`
	CredentialsByIssuer       map[string]uint64         `json:"credentials_by_issuer" yaml:"credentials_by_issuer"`
	CredentialIssuanceRate    float64                   `json:"credential_issuance_rate" yaml:"credential_issuance_rate"`
	CredentialVerifications   uint64                    `json:"credential_verifications" yaml:"credential_verifications"`
	AverageCredentialLifetime time.Duration             `json:"average_credential_lifetime" yaml:"average_credential_lifetime"`
	PresentationRequests      uint64                    `json:"presentation_requests" yaml:"presentation_requests"`
	PresentationSuccess       uint64                    `json:"presentation_success" yaml:"presentation_success"`
}

// GeographicStatistics contains geographic distribution analytics
type GeographicStatistics struct {
	IdentitiesByCountry       map[string]uint64         `json:"identities_by_country" yaml:"identities_by_country"`
	IdentitiesByState         map[string]uint64         `json:"identities_by_state" yaml:"identities_by_state"`
	IdentitiesByCity          map[string]uint64         `json:"identities_by_city" yaml:"identities_by_city"`
	VerificationsByRegion     map[string]uint64         `json:"verifications_by_region" yaml:"verifications_by_region"`
	PopularRegions            []RegionUsage             `json:"popular_regions" yaml:"popular_regions"`
	RuralVsUrbanSplit         *RuralUrbanSplit          `json:"rural_vs_urban_split" yaml:"rural_vs_urban_split"`
	GeographicGrowthTrends    map[string]*GrowthTrend   `json:"geographic_growth_trends" yaml:"geographic_growth_trends"`
}

// DemographicStatistics contains demographic analytics
type DemographicStatistics struct {
	AgeDistribution           map[string]uint64         `json:"age_distribution" yaml:"age_distribution"`
	GenderDistribution        map[string]uint64         `json:"gender_distribution" yaml:"gender_distribution"`
	EducationLevels           map[string]uint64         `json:"education_levels" yaml:"education_levels"`
	OccupationCategories      map[string]uint64         `json:"occupation_categories" yaml:"occupation_categories"`
	IncomeRanges              map[string]uint64         `json:"income_ranges" yaml:"income_ranges"`
	DigitalLiteracyLevels     map[string]uint64         `json:"digital_literacy_levels" yaml:"digital_literacy_levels"`
	DeviceOwnership           map[string]uint64         `json:"device_ownership" yaml:"device_ownership"`
	InternetAccess            *InternetAccessStats      `json:"internet_access" yaml:"internet_access"`
}

// LanguageStatistics contains language and localization analytics
type LanguageStatistics struct {
	LanguageDistribution      map[LanguageCode]uint64   `json:"language_distribution" yaml:"language_distribution"`
	RegionalLanguageUsage     map[string]map[LanguageCode]uint64 `json:"regional_language_usage" yaml:"regional_language_usage"`
	LanguagePreferenceChanges uint64                    `json:"language_preference_changes" yaml:"language_preference_changes"`
	CulturalContentRequests   uint64                    `json:"cultural_content_requests" yaml:"cultural_content_requests"`
	LocalizationCoverage      float64                   `json:"localization_coverage" yaml:"localization_coverage"`
	TranslationRequests       map[string]uint64         `json:"translation_requests" yaml:"translation_requests"`
	FestivalGreetingsServed   map[string]uint64         `json:"festival_greetings_served" yaml:"festival_greetings_served"`
	CulturalQuotesRequested   map[string]uint64         `json:"cultural_quotes_requested" yaml:"cultural_quotes_requested"`
}

// OfflineStatistics contains offline verification analytics
type OfflineStatistics struct {
	OfflineVerificationPackages uint64                  `json:"offline_verification_packages" yaml:"offline_verification_packages"`
	OfflineVerifications        uint64                  `json:"offline_verifications" yaml:"offline_verifications"`
	OfflineSuccessRate          float64                 `json:"offline_success_rate" yaml:"offline_success_rate"`
	PackagesByFormat            map[OfflineCredentialFormat]uint64 `json:"packages_by_format" yaml:"packages_by_format"`
	VerificationsByMode         map[OfflineVerificationMode]uint64 `json:"verifications_by_mode" yaml:"verifications_by_mode"`
	RegisteredDevices           uint64                  `json:"registered_devices" yaml:"registered_devices"`
	ActiveDevices               uint64                  `json:"active_devices" yaml:"active_devices"`
	DevicesByType               map[string]uint64       `json:"devices_by_type" yaml:"devices_by_type"`
	OfflineBackupsCreated       uint64                  `json:"offline_backups_created" yaml:"offline_backups_created"`
	EmergencyModeUsage          uint64                  `json:"emergency_mode_usage" yaml:"emergency_mode_usage"`
	AverageOfflineDuration      time.Duration           `json:"average_offline_duration" yaml:"average_offline_duration"`
}

// SecurityStatistics contains security-related analytics
type SecurityStatistics struct {
	SecurityIncidents           uint64                  `json:"security_incidents" yaml:"security_incidents"`
	BlockedAttacks              uint64                  `json:"blocked_attacks" yaml:"blocked_attacks"`
	SuspiciousActivities        uint64                  `json:"suspicious_activities" yaml:"suspicious_activities"`
	FailedAuthenticationAttempts uint64                 `json:"failed_authentication_attempts" yaml:"failed_authentication_attempts"`
	AccountLockouts             uint64                  `json:"account_lockouts" yaml:"account_lockouts"`
	PasswordResets              uint64                  `json:"password_resets" yaml:"password_resets"`
	TwoFactorAdoption           float64                 `json:"two_factor_adoption" yaml:"two_factor_adoption"`
	BiometricSecurityScore      float64                 `json:"biometric_security_score" yaml:"biometric_security_score"`
	EncryptionCoverage          float64                 `json:"encryption_coverage" yaml:"encryption_coverage"`
	SecurityAudits              uint64                  `json:"security_audits" yaml:"security_audits"`
	ComplianceScore             float64                 `json:"compliance_score" yaml:"compliance_score"`
}

// FraudDetectionStatistics contains fraud detection analytics
type FraudDetectionStatistics struct {
	TotalFraudAttempts          uint64                  `json:"total_fraud_attempts" yaml:"total_fraud_attempts"`
	DetectedFraudAttempts       uint64                  `json:"detected_fraud_attempts" yaml:"detected_fraud_attempts"`
	PreventedFraudAttempts      uint64                  `json:"prevented_fraud_attempts" yaml:"prevented_fraud_attempts"`
	FraudDetectionRate          float64                 `json:"fraud_detection_rate" yaml:"fraud_detection_rate"`
	FalsePositiveRate           float64                 `json:"false_positive_rate" yaml:"false_positive_rate"`
	FraudByType                 map[string]uint64       `json:"fraud_by_type" yaml:"fraud_by_type"`
	FraudByRegion               map[string]uint64       `json:"fraud_by_region" yaml:"fraud_by_region"`
	SyntheticIdentityDetection  uint64                  `json:"synthetic_identity_detection" yaml:"synthetic_identity_detection"`
	DuplicateIdentityDetection  uint64                  `json:"duplicate_identity_detection" yaml:"duplicate_identity_detection"`
	BiometricSpoofingAttempts   uint64                  `json:"biometric_spoofing_attempts" yaml:"biometric_spoofing_attempts"`
	DocumentForgeryDetection    uint64                  `json:"document_forgery_detection" yaml:"document_forgery_detection"`
	MLModelAccuracy             float64                 `json:"ml_model_accuracy" yaml:"ml_model_accuracy"`
}

// PerformanceMetrics contains system performance analytics
type PerformanceMetrics struct {
	AverageResponseTime         time.Duration           `json:"average_response_time" yaml:"average_response_time"`
	ThroughputPerSecond         float64                 `json:"throughput_per_second" yaml:"throughput_per_second"`
	SystemUptime                float64                 `json:"system_uptime" yaml:"system_uptime"`
	ErrorRate                   float64                 `json:"error_rate" yaml:"error_rate"`
	CacheHitRate                float64                 `json:"cache_hit_rate" yaml:"cache_hit_rate"`
	DatabasePerformance         *DatabaseMetrics        `json:"database_performance" yaml:"database_performance"`
	APIPerformance              *APIMetrics             `json:"api_performance" yaml:"api_performance"`
	BlockchainMetrics           *BlockchainMetrics      `json:"blockchain_metrics" yaml:"blockchain_metrics"`
	ResourceUtilization         *ResourceUtilization    `json:"resource_utilization" yaml:"resource_utilization"`
}

// UsagePatterns contains usage pattern analytics
type UsagePatterns struct {
	PeakUsageHours              []int                   `json:"peak_usage_hours" yaml:"peak_usage_hours"`
	WeeklyUsagePattern          map[string]uint64       `json:"weekly_usage_pattern" yaml:"weekly_usage_pattern"`
	MonthlyTrends               map[string]uint64       `json:"monthly_trends" yaml:"monthly_trends"`
	SeasonalPatterns            map[string]uint64       `json:"seasonal_patterns" yaml:"seasonal_patterns"`
	FeatureUsageFrequency       map[string]uint64       `json:"feature_usage_frequency" yaml:"feature_usage_frequency"`
	UserRetentionRates          *RetentionRates         `json:"user_retention_rates" yaml:"user_retention_rates"`
	SessionDuration             *SessionMetrics         `json:"session_duration" yaml:"session_duration"`
	DevicePlatformUsage         map[string]uint64       `json:"device_platform_usage" yaml:"device_platform_usage"`
}

// GrowthMetrics contains growth and trend analytics
type GrowthMetrics struct {
	UserGrowthRate              float64                 `json:"user_growth_rate" yaml:"user_growth_rate"`
	VerificationGrowthRate      float64                 `json:"verification_growth_rate" yaml:"verification_growth_rate"`
	MonthlyActiveUsers          uint64                  `json:"monthly_active_users" yaml:"monthly_active_users"`
	DailyActiveUsers            uint64                  `json:"daily_active_users" yaml:"daily_active_users"`
	NewUserAcquisition          *AcquisitionMetrics     `json:"new_user_acquisition" yaml:"new_user_acquisition"`
	ChurnRate                   float64                 `json:"churn_rate" yaml:"churn_rate"`
	UserEngagementScore         float64                 `json:"user_engagement_score" yaml:"user_engagement_score"`
	FeatureAdoptionRates        map[string]float64      `json:"feature_adoption_rates" yaml:"feature_adoption_rates"`
	GeographicExpansion         *ExpansionMetrics       `json:"geographic_expansion" yaml:"geographic_expansion"`
}

// ComplianceMetrics contains compliance and audit analytics
type ComplianceMetrics struct {
	GDPRCompliance              float64                 `json:"gdpr_compliance" yaml:"gdpr_compliance"`
	DPDPActCompliance           float64                 `json:"dpdp_act_compliance" yaml:"dpdp_act_compliance"`
	CCPACompliance              float64                 `json:"ccpa_compliance" yaml:"ccpa_compliance"`
	SOC2Compliance              float64                 `json:"soc2_compliance" yaml:"soc2_compliance"`
	ISO27001Compliance          float64                 `json:"iso27001_compliance" yaml:"iso27001_compliance"`
	DataRetentionCompliance     float64                 `json:"data_retention_compliance" yaml:"data_retention_compliance"`
	ConsentManagementScore      float64                 `json:"consent_management_score" yaml:"consent_management_score"`
	DataMinimizationScore       float64                 `json:"data_minimization_score" yaml:"data_minimization_score"`
	AuditTrailCompleteness      float64                 `json:"audit_trail_completeness" yaml:"audit_trail_completeness"`
	ComplianceIncidents         uint64                  `json:"compliance_incidents" yaml:"compliance_incidents"`
	DataBreaches                uint64                  `json:"data_breaches" yaml:"data_breaches"`
	PrivacyRequests             uint64                  `json:"privacy_requests" yaml:"privacy_requests"`
}

// Supporting types for analytics

// TimeRange represents a time range for analytics
type TimeRange struct {
	StartTime   time.Time `json:"start_time" yaml:"start_time"`
	EndTime     time.Time `json:"end_time" yaml:"end_time"`
	Duration    time.Duration `json:"duration" yaml:"duration"`
	Granularity string    `json:"granularity" yaml:"granularity"` // hour, day, week, month
}

// PeakTime represents peak usage time information
type PeakTime struct {
	Hour        int     `json:"hour" yaml:"hour"`
	Count       uint64  `json:"count" yaml:"count"`
	Percentage  float64 `json:"percentage" yaml:"percentage"`
}

// QualityScoreDistribution represents distribution of quality scores
type QualityScoreDistribution struct {
	Excellent   uint64  `json:"excellent" yaml:"excellent"`   // 90-100%
	Good        uint64  `json:"good" yaml:"good"`             // 80-89%
	Fair        uint64  `json:"fair" yaml:"fair"`             // 70-79%
	Poor        uint64  `json:"poor" yaml:"poor"`             // <70%
	Average     float64 `json:"average" yaml:"average"`
}

// RegionUsage represents regional usage statistics
type RegionUsage struct {
	Region      string  `json:"region" yaml:"region"`
	Count       uint64  `json:"count" yaml:"count"`
	Percentage  float64 `json:"percentage" yaml:"percentage"`
	GrowthRate  float64 `json:"growth_rate" yaml:"growth_rate"`
}

// RuralUrbanSplit represents rural vs urban distribution
type RuralUrbanSplit struct {
	Rural       uint64  `json:"rural" yaml:"rural"`
	Urban       uint64  `json:"urban" yaml:"urban"`
	RuralPercent float64 `json:"rural_percent" yaml:"rural_percent"`
	UrbanPercent float64 `json:"urban_percent" yaml:"urban_percent"`
}

// GrowthTrend represents growth trend information
type GrowthTrend struct {
	Period      string  `json:"period" yaml:"period"`
	StartValue  uint64  `json:"start_value" yaml:"start_value"`
	EndValue    uint64  `json:"end_value" yaml:"end_value"`
	GrowthRate  float64 `json:"growth_rate" yaml:"growth_rate"`
	Projection  uint64  `json:"projection" yaml:"projection"`
}

// InternetAccessStats represents internet access statistics
type InternetAccessStats struct {
	HighSpeedInternet   uint64  `json:"high_speed_internet" yaml:"high_speed_internet"`
	MobileDataOnly      uint64  `json:"mobile_data_only" yaml:"mobile_data_only"`
	LimitedAccess       uint64  `json:"limited_access" yaml:"limited_access"`
	NoAccess            uint64  `json:"no_access" yaml:"no_access"`
	AverageSpeed        float64 `json:"average_speed" yaml:"average_speed"` // Mbps
}

// DatabaseMetrics represents database performance metrics
type DatabaseMetrics struct {
	QueryResponseTime   time.Duration `json:"query_response_time" yaml:"query_response_time"`
	ConnectionPoolUsage float64       `json:"connection_pool_usage" yaml:"connection_pool_usage"`
	DeadlockCount       uint64        `json:"deadlock_count" yaml:"deadlock_count"`
	SlowQueryCount      uint64        `json:"slow_query_count" yaml:"slow_query_count"`
	CacheEfficiency     float64       `json:"cache_efficiency" yaml:"cache_efficiency"`
}

// APIMetrics represents API performance metrics
type APIMetrics struct {
	RequestsPerSecond   float64       `json:"requests_per_second" yaml:"requests_per_second"`
	AverageLatency      time.Duration `json:"average_latency" yaml:"average_latency"`
	ErrorRate           float64       `json:"error_rate" yaml:"error_rate"`
	TimeoutRate         float64       `json:"timeout_rate" yaml:"timeout_rate"`
	EndpointPerformance map[string]*EndpointMetrics `json:"endpoint_performance" yaml:"endpoint_performance"`
}

// BlockchainMetrics represents blockchain performance metrics
type BlockchainMetrics struct {
	BlockTime           time.Duration `json:"block_time" yaml:"block_time"`
	TransactionsPerBlock float64      `json:"transactions_per_block" yaml:"transactions_per_block"`
	NetworkLatency      time.Duration `json:"network_latency" yaml:"network_latency"`
	ConsensusTime       time.Duration `json:"consensus_time" yaml:"consensus_time"`
	NodeSyncStatus      float64       `json:"node_sync_status" yaml:"node_sync_status"`
}

// ResourceUtilization represents system resource utilization
type ResourceUtilization struct {
	CPUUsage            float64 `json:"cpu_usage" yaml:"cpu_usage"`
	MemoryUsage         float64 `json:"memory_usage" yaml:"memory_usage"`
	DiskUsage           float64 `json:"disk_usage" yaml:"disk_usage"`
	NetworkBandwidth    float64 `json:"network_bandwidth" yaml:"network_bandwidth"`
	StorageIOPS         float64 `json:"storage_iops" yaml:"storage_iops"`
}

// RetentionRates represents user retention analytics
type RetentionRates struct {
	Day1        float64 `json:"day_1" yaml:"day_1"`
	Day7        float64 `json:"day_7" yaml:"day_7"`
	Day30       float64 `json:"day_30" yaml:"day_30"`
	Day90       float64 `json:"day_90" yaml:"day_90"`
	Day365      float64 `json:"day_365" yaml:"day_365"`
}

// SessionMetrics represents session analytics
type SessionMetrics struct {
	AverageSessionDuration time.Duration `json:"average_session_duration" yaml:"average_session_duration"`
	MedianSessionDuration  time.Duration `json:"median_session_duration" yaml:"median_session_duration"`
	BounceRate             float64       `json:"bounce_rate" yaml:"bounce_rate"`
	SessionsPerUser        float64       `json:"sessions_per_user" yaml:"sessions_per_user"`
}

// AcquisitionMetrics represents user acquisition analytics
type AcquisitionMetrics struct {
	OrganicSignups      uint64 `json:"organic_signups" yaml:"organic_signups"`
	ReferralSignups     uint64 `json:"referral_signups" yaml:"referral_signups"`
	PartnerSignups      uint64 `json:"partner_signups" yaml:"partner_signups"`
	GovernmentSignups   uint64 `json:"government_signups" yaml:"government_signups"`
	CostPerAcquisition  float64 `json:"cost_per_acquisition" yaml:"cost_per_acquisition"`
}

// ExpansionMetrics represents geographic expansion analytics
type ExpansionMetrics struct {
	NewRegions          []string          `json:"new_regions" yaml:"new_regions"`
	RegionPenetration   map[string]float64 `json:"region_penetration" yaml:"region_penetration"`
	LocalizationGaps    []string          `json:"localization_gaps" yaml:"localization_gaps"`
	RegulatoryReadiness map[string]float64 `json:"regulatory_readiness" yaml:"regulatory_readiness"`
}

// EndpointMetrics represents individual API endpoint metrics
type EndpointMetrics struct {
	Path            string        `json:"path" yaml:"path"`
	Method          string        `json:"method" yaml:"method"`
	RequestCount    uint64        `json:"request_count" yaml:"request_count"`
	AverageLatency  time.Duration `json:"average_latency" yaml:"average_latency"`
	ErrorCount      uint64        `json:"error_count" yaml:"error_count"`
	ErrorRate       float64       `json:"error_rate" yaml:"error_rate"`
	P95Latency      time.Duration `json:"p95_latency" yaml:"p95_latency"`
	P99Latency      time.Duration `json:"p99_latency" yaml:"p99_latency"`
}

// Constructor functions

// NewIdentityAnalytics creates a new identity analytics instance
func NewIdentityAnalytics() *IdentityAnalytics {
	return &IdentityAnalytics{
		VerificationStats:    &VerificationStatistics{},
		KYCStats:            &KYCStatistics{},
		BiometricStats:      &BiometricStatistics{},
		CredentialStats:     &CredentialStatistics{},
		GeographicStats:     &GeographicStatistics{},
		DemographicStats:    &DemographicStatistics{},
		LanguageStats:       &LanguageStatistics{},
		OfflineStats:        &OfflineStatistics{},
		SecurityStats:       &SecurityStatistics{},
		FraudDetectionStats: &FraudDetectionStatistics{},
		PerformanceMetrics:  &PerformanceMetrics{},
		UsagePatterns:       &UsagePatterns{},
		GrowthMetrics:       &GrowthMetrics{},
		ComplianceMetrics:   &ComplianceMetrics{},
		GeneratedAt:         time.Now(),
	}
}

// NewTimeRange creates a new time range
func NewTimeRange(start, end time.Time, granularity string) *TimeRange {
	return &TimeRange{
		StartTime:   start,
		EndTime:     end,
		Duration:    end.Sub(start),
		Granularity: granularity,
	}
}

// Utility methods

// CalculateSuccessRate calculates success rate from total and successful counts
func CalculateSuccessRate(successful, total uint64) float64 {
	if total == 0 {
		return 0.0
	}
	return float64(successful) / float64(total) * 100.0
}

// CalculateGrowthRate calculates growth rate between two values
func CalculateGrowthRate(current, previous uint64) float64 {
	if previous == 0 {
		return 0.0
	}
	return (float64(current) - float64(previous)) / float64(previous) * 100.0
}

// CalculatePercentage calculates percentage of part from total
func CalculatePercentage(part, total uint64) float64 {
	if total == 0 {
		return 0.0
	}
	return float64(part) / float64(total) * 100.0
}

// IsDataFresh checks if analytics data is fresh (within threshold)
func (ia *IdentityAnalytics) IsDataFresh(threshold time.Duration) bool {
	return time.Since(ia.GeneratedAt) <= threshold
}

// GetOverallHealthScore calculates an overall system health score
func (ia *IdentityAnalytics) GetOverallHealthScore() float64 {
	scores := []float64{
		ia.VerificationStats.SuccessRate,
		ia.KYCStats.KYCApprovalRate,
		ia.BiometricStats.BiometricSuccessRate,
		ia.PerformanceMetrics.SystemUptime,
		ia.SecurityStats.ComplianceScore,
		(100.0 - ia.SecurityStats.FailedAuthenticationAttempts/float64(max(ia.TotalIdentities, 1))*100.0),
	}

	var total float64
	for _, score := range scores {
		total += score
	}

	return total / float64(len(scores))
}

// max returns the maximum of two uint64 values
func max(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

// GetTrendDirection determines trend direction based on growth rate
func GetTrendDirection(growthRate float64) string {
	if growthRate > 5.0 {
		return "📈 Strong Growth"
	} else if growthRate > 0 {
		return "📊 Growth"
	} else if growthRate == 0 {
		return "📊 Stable"
	} else if growthRate > -5.0 {
		return "📉 Decline"
	} else {
		return "📉 Strong Decline"
	}
}

// GetRiskLevel determines risk level based on various metrics
func (ia *IdentityAnalytics) GetRiskLevel() string {
	riskFactors := 0
	
	if ia.SecurityStats.FailedAuthenticationAttempts > 1000 {
		riskFactors++
	}
	if ia.FraudDetectionStats.FraudDetectionRate < 80.0 {
		riskFactors++
	}
	if ia.PerformanceMetrics.ErrorRate > 5.0 {
		riskFactors++
	}
	if ia.ComplianceMetrics.ComplianceIncidents > 5 {
		riskFactors++
	}

	switch riskFactors {
	case 0:
		return "🟢 Low Risk"
	case 1:
		return "🟡 Medium Risk"
	case 2:
		return "🟠 High Risk"
	default:
		return "🔴 Critical Risk"
	}
}