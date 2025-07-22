package keeper

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GDPRPrivacyComplianceSystem manages GDPR and privacy compliance
type GDPRPrivacyComplianceSystem struct {
	keeper                   Keeper
	consentManager           *ConsentManagementSystem
	dataProcessor            *PersonalDataProcessor
	rightsManager            *DataSubjectRightsManager
	privacyEngine            *PrivacyByDesignEngine
	dataProtectionOfficer    *DataProtectionOfficer
	breachNotificationSystem *BreachNotificationSystem
	impactAssessment         *PrivacyImpactAssessment
	auditLogger              *PrivacyAuditLogger
	mu                       sync.RWMutex
}

// ConsentManagementSystem handles consent collection and management
type ConsentManagementSystem struct {
	consents              map[string]*ConsentRecord
	consentTemplates      map[string]*ConsentTemplate
	consentVersioning     *ConsentVersionManager
	withdrawalProcessor   *ConsentWithdrawalProcessor
	preferenceManager     *PrivacyPreferenceManager
	cookieManager         *CookieConsentManager
	granularityEngine     *ConsentGranularityEngine
}

// ConsentRecord represents a user's consent
type ConsentRecord struct {
	ConsentID           string
	UserID              string
	ConsentType         ConsentType
	Purposes            []ProcessingPurpose
	DataCategories      []DataCategory
	GrantedAt           time.Time
	ExpiresAt           *time.Time
	WithdrawnAt         *time.Time
	ConsentMethod       ConsentMethod
	ConsentVersion      string
	IPAddress           string
	UserAgent           string
	LegalBasis          LegalBasis
	ThirdPartySharing   []ThirdPartyConsent
	SpecialCategories   []SpecialCategoryConsent
	ChildConsent        *ChildConsentInfo
	ConsentProof        ConsentProof
	Preferences         map[string]bool
	UpdateHistory       []ConsentUpdate
}

// PersonalDataProcessor handles personal data processing
type PersonalDataProcessor struct {
	dataInventory        *PersonalDataInventory
	processingRegister   *ProcessingActivityRegister
	dataMinimizer        *DataMinimizationEngine
	pseudonymizer        *PseudonymizationService
	anonymizer           *AnonymizationService
	encryptionService    *DataEncryptionService
	retentionManager     *DataRetentionManager
	dataCleaner          *DataCleaningService
}

// PersonalDataInventory tracks all personal data
type PersonalDataInventory struct {
	DataAssets           map[string]*DataAsset
	DataFlows            map[string]*DataFlow
	ProcessingActivities map[string]*ProcessingActivity
	DataMapping          *DataMappingSystem
	ClassificationEngine *DataClassificationEngine
	SensitivityLabeler   *SensitivityLabeler
}

// DataAsset represents a collection of personal data
type DataAsset struct {
	AssetID             string
	Name                string
	Description         string
	DataCategories      []DataCategory
	DataSubjects        []DataSubjectType
	Sources             []DataSource
	RetentionPeriod     time.Duration
	LegalBasis          LegalBasis
	SecurityMeasures    []SecurityMeasure
	AccessControls      []AccessControl
	ThirdPartyAccess    []ThirdPartyAccess
	CrossBorderTransfer []CrossBorderTransfer
	LastUpdated         time.Time
	DataController      string
	DataProcessors      []string
}

// DataSubjectRightsManager handles GDPR rights requests
type DataSubjectRightsManager struct {
	requestQueue         *RightsRequestQueue
	accessProvider       *DataAccessProvider
	rectificationEngine  *DataRectificationEngine
	erasureProcessor     *DataErasureProcessor
	portabilityExporter  *DataPortabilityExporter
	restrictionManager   *ProcessingRestrictionManager
	objectionHandler     *ProcessingObjectionHandler
	automatedDecisionMgr *AutomatedDecisionManager
}

// RightsRequest represents a data subject rights request
type RightsRequest struct {
	RequestID           string
	UserID              string
	RequestType         RightsRequestType
	RequestedAt         time.Time
	Description         string
	VerificationStatus  VerificationStatus
	ProcessingStatus    ProcessingStatus
	AssignedTo          string
	Deadline            time.Time
	CompletedAt         *time.Time
	ResponseMethod      ResponseMethod
	RequestData         map[string]interface{}
	ProcessingNotes     []ProcessingNote
	Attachments         []Document
	Response            *RightsResponse
}

// PrivacyByDesignEngine implements privacy by design principles
type PrivacyByDesignEngine struct {
	designPrinciples     *PrivacyDesignPrinciples
	privacyPatterns      *PrivacyPatternLibrary
	defaultSettings      *PrivacyDefaultSettings
	architectureReview   *PrivacyArchitectureReview
	codeAnalyzer         *PrivacyCodeAnalyzer
	configValidator      *PrivacyConfigValidator
}

// DataProtectionOfficer manages DPO responsibilities
type DataProtectionOfficer struct {
	dpoInfo              *DPOInformation
	advisoryService      *PrivacyAdvisoryService
	trainingProgram      *PrivacyTrainingProgram
	complianceMonitor    *ComplianceMonitoringService
	stakeholderManager   *StakeholderRelationshipManager
	documentationMgr     *ComplianceDocumentationManager
}

// BreachNotificationSystem handles data breach incidents
type BreachNotificationSystem struct {
	breachDetector       *BreachDetectionEngine
	incidentManager      *BreachIncidentManager
	notificationEngine   *NotificationDispatcher
	impactAnalyzer       *BreachImpactAnalyzer
	remediationTracker   *BreachRemediationTracker
	regulatoryReporter   *RegulatoryBreachReporter
}

// DataBreach represents a data breach incident
type DataBreach struct {
	BreachID            string
	DetectedAt          time.Time
	BreachType          BreachType
	Severity            BreachSeverity
	AffectedDataTypes   []DataCategory
	AffectedUsers       []string
	EstimatedUserCount  int
	BreachSource        string
	AttackVector        string
	DataCompromised     []DataCompromiseDetail
	RiskAssessment      *BreachRiskAssessment
	NotificationStatus  NotificationStatus
	RemediationActions  []RemediationAction
	RegulatoryFilings   []RegulatoryFiling
	LessonsLearned      []string
}

// Types and enums
type ConsentType int
type ProcessingPurpose int
type DataCategory int
type ConsentMethod int
type LegalBasis int
type DataSubjectType int
type RightsRequestType int
type VerificationStatus int
type ProcessingStatus int
type ResponseMethod int
type BreachType int
type BreachSeverity int
type NotificationStatus int

const (
	// Consent Types
	ServiceConsent ConsentType = iota
	MarketingConsent
	AnalyticsConsent
	ThirdPartyConsent
	CookieConsent
	
	// Legal Basis
	ConsentBasis LegalBasis = iota
	ContractBasis
	LegalObligationBasis
	VitalInterestsBasis
	PublicTaskBasis
	LegitimateInterestsBasis
	
	// Rights Request Types
	AccessRequest RightsRequestType = iota
	RectificationRequest
	ErasureRequest
	RestrictProcessingRequest
	DataPortabilityRequest
	ObjectionRequest
	AutomatedDecisionRequest
	
	// Breach Severity
	LowSeverity BreachSeverity = iota
	MediumSeverity
	HighSeverity
	CriticalSeverity
)

// Core GDPR compliance methods

// ProcessConsentRequest processes a consent request
func (k Keeper) ProcessConsentRequest(ctx context.Context, request ConsentRequest) (*ConsentRecord, error) {
	gpcs := k.getGDPRPrivacyComplianceSystem()
	
	// Validate consent request
	if err := gpcs.validateConsentRequest(request); err != nil {
		return nil, fmt.Errorf("invalid consent request: %w", err)
	}
	
	// Check if user is a child (requires parental consent)
	isChild, parentalConsent := gpcs.checkChildStatus(request.UserID)
	if isChild && parentalConsent == nil {
		return nil, fmt.Errorf("parental consent required for users under 16")
	}
	
	// Create consent record
	consent := &ConsentRecord{
		ConsentID:         generateID("CONSENT"),
		UserID:            request.UserID,
		ConsentType:       request.ConsentType,
		Purposes:          request.Purposes,
		DataCategories:    request.DataCategories,
		GrantedAt:         time.Now(),
		ConsentMethod:     request.Method,
		ConsentVersion:    gpcs.consentManager.getCurrentVersion(request.ConsentType),
		IPAddress:         request.IPAddress,
		UserAgent:         request.UserAgent,
		LegalBasis:        request.LegalBasis,
		ThirdPartySharing: request.ThirdPartySharing,
		Preferences:       request.Preferences,
	}
	
	// Handle child consent if applicable
	if isChild {
		consent.ChildConsent = &ChildConsentInfo{
			ChildAge:          parentalConsent.ChildAge,
			ParentID:          parentalConsent.ParentID,
			VerificationMethod: parentalConsent.VerificationMethod,
			VerifiedAt:        parentalConsent.VerifiedAt,
		}
	}
	
	// Generate consent proof
	consent.ConsentProof = gpcs.generateConsentProof(consent)
	
	// Store consent
	gpcs.consentManager.consents[consent.ConsentID] = consent
	if err := k.storeConsent(ctx, consent); err != nil {
		return nil, fmt.Errorf("failed to store consent: %w", err)
	}
	
	// Update privacy preferences
	gpcs.consentManager.preferenceManager.updatePreferences(request.UserID, consent.Preferences)
	
	// Log consent event
	gpcs.auditLogger.logConsentEvent(ConsentGrantedEvent, consent)
	
	return consent, nil
}

// ProcessDataSubjectRequest handles GDPR rights requests
func (k Keeper) ProcessDataSubjectRequest(ctx context.Context, request DataSubjectRequest) (*RightsResponse, error) {
	gpcs := k.getGDPRPrivacyComplianceSystem()
	
	// Verify identity
	verified, err := gpcs.verifyDataSubjectIdentity(request)
	if err != nil || !verified {
		return nil, fmt.Errorf("identity verification failed: %w", err)
	}
	
	// Create rights request
	rightsRequest := &RightsRequest{
		RequestID:          generateID("DSR"),
		UserID:             request.UserID,
		RequestType:        request.RequestType,
		RequestedAt:        time.Now(),
		Description:        request.Description,
		VerificationStatus: VerificationCompleted,
		ProcessingStatus:   ProcessingInProgress,
		Deadline:           gpcs.calculateDeadline(request.RequestType),
		ResponseMethod:     request.PreferredResponseMethod,
		RequestData:        request.AdditionalData,
	}
	
	// Add to processing queue
	gpcs.rightsManager.requestQueue.enqueue(rightsRequest)
	
	// Process based on request type
	var response *RightsResponse
	
	switch request.RequestType {
	case AccessRequest:
		response, err = gpcs.processAccessRequest(ctx, rightsRequest)
		
	case RectificationRequest:
		response, err = gpcs.processRectificationRequest(ctx, rightsRequest)
		
	case ErasureRequest:
		response, err = gpcs.processErasureRequest(ctx, rightsRequest)
		
	case RestrictProcessingRequest:
		response, err = gpcs.processRestrictionRequest(ctx, rightsRequest)
		
	case DataPortabilityRequest:
		response, err = gpcs.processPortabilityRequest(ctx, rightsRequest)
		
	case ObjectionRequest:
		response, err = gpcs.processObjectionRequest(ctx, rightsRequest)
		
	case AutomatedDecisionRequest:
		response, err = gpcs.processAutomatedDecisionRequest(ctx, rightsRequest)
	}
	
	if err != nil {
		rightsRequest.ProcessingStatus = ProcessingFailed
		return nil, fmt.Errorf("failed to process request: %w", err)
	}
	
	// Update request status
	rightsRequest.ProcessingStatus = ProcessingCompleted
	rightsRequest.CompletedAt = timePtr(time.Now())
	rightsRequest.Response = response
	
	// Store request and response
	if err := k.storeRightsRequest(ctx, rightsRequest); err != nil {
		return nil, fmt.Errorf("failed to store request: %w", err)
	}
	
	// Log the request
	gpcs.auditLogger.logRightsRequest(rightsRequest)
	
	// Send response to user
	gpcs.sendResponseToUser(rightsRequest, response)
	
	return response, nil
}

// Data processing methods

func (gpcs *GDPRPrivacyComplianceSystem) processAccessRequest(ctx context.Context, request *RightsRequest) (*RightsResponse, error) {
	// Collect all personal data
	personalData, err := gpcs.dataProcessor.collectPersonalData(request.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to collect personal data: %w", err)
	}
	
	// Format data for export
	exportData := gpcs.rightsManager.accessProvider.formatDataForExport(personalData)
	
	// Include processing information
	processingInfo := gpcs.dataProcessor.getProcessingInformation(request.UserID)
	exportData.ProcessingActivities = processingInfo
	
	// Include consent history
	consentHistory := gpcs.consentManager.getConsentHistory(request.UserID)
	exportData.ConsentHistory = consentHistory
	
	response := &RightsResponse{
		ResponseID:   generateID("RESP"),
		RequestID:    request.RequestID,
		ResponseType: AccessResponse,
		Data:         exportData,
		GeneratedAt:  time.Now(),
		ValidUntil:   time.Now().Add(30 * 24 * time.Hour),
	}
	
	return response, nil
}

func (gpcs *GDPRPrivacyComplianceSystem) processErasureRequest(ctx context.Context, request *RightsRequest) (*RightsResponse, error) {
	// Check if erasure can be performed
	erasureCheck := gpcs.rightsManager.erasureProcessor.checkErasureEligibility(request.UserID)
	if !erasureCheck.CanErase {
		return &RightsResponse{
			ResponseID:   generateID("RESP"),
			RequestID:    request.RequestID,
			ResponseType: ErasureResponse,
			Success:      false,
			Reason:       erasureCheck.Reason,
			LegalBasis:   erasureCheck.RetentionBasis,
		}, nil
	}
	
	// Perform erasure
	erasureResult := gpcs.rightsManager.erasureProcessor.performErasure(ctx, request.UserID, erasureCheck.ErasableData)
	
	// Notify third parties
	if len(erasureResult.ThirdPartiesToNotify) > 0 {
		gpcs.notifyThirdPartiesOfErasure(request.UserID, erasureResult.ThirdPartiesToNotify)
	}
	
	response := &RightsResponse{
		ResponseID:    generateID("RESP"),
		RequestID:     request.RequestID,
		ResponseType:  ErasureResponse,
		Success:       erasureResult.Success,
		ErasedData:    erasureResult.ErasedCategories,
		RetainedData:  erasureResult.RetainedCategories,
		GeneratedAt:   time.Now(),
	}
	
	return response, nil
}

// Data breach handling

func (k Keeper) ReportDataBreach(ctx context.Context, breachReport BreachReport) (*DataBreach, error) {
	gpcs := k.getGDPRPrivacyComplianceSystem()
	
	// Validate breach report
	if err := gpcs.validateBreachReport(breachReport); err != nil {
		return nil, fmt.Errorf("invalid breach report: %w", err)
	}
	
	// Create breach incident
	breach := &DataBreach{
		BreachID:           generateID("BREACH"),
		DetectedAt:         breachReport.DetectedAt,
		BreachType:         breachReport.BreachType,
		Severity:           gpcs.assessBreachSeverity(breachReport),
		AffectedDataTypes:  breachReport.AffectedDataTypes,
		BreachSource:       breachReport.Source,
		AttackVector:       breachReport.AttackVector,
		NotificationStatus: NotificationPending,
	}
	
	// Analyze breach impact
	impactAnalysis := gpcs.breachNotificationSystem.impactAnalyzer.analyzeImpact(breachReport)
	breach.RiskAssessment = impactAnalysis.RiskAssessment
	breach.AffectedUsers = impactAnalysis.AffectedUsers
	breach.EstimatedUserCount = len(impactAnalysis.AffectedUsers)
	breach.DataCompromised = impactAnalysis.CompromisedData
	
	// Determine notification requirements
	notificationReq := gpcs.determineNotificationRequirements(breach)
	
	// 72-hour regulatory notification if high risk
	if notificationReq.RequiresRegulatoryNotification {
		deadline := breach.DetectedAt.Add(72 * time.Hour)
		gpcs.scheduleRegulatoryNotification(breach, deadline)
	}
	
	// User notification if required
	if notificationReq.RequiresUserNotification {
		gpcs.scheduleUserNotifications(breach, notificationReq.NotificationPriority)
	}
	
	// Start remediation
	remediationPlan := gpcs.breachNotificationSystem.createRemediationPlan(breach)
	breach.RemediationActions = remediationPlan.Actions
	
	// Store breach record
	if err := k.storeDataBreach(ctx, breach); err != nil {
		return nil, fmt.Errorf("failed to store breach: %w", err)
	}
	
	// Log breach event
	gpcs.auditLogger.logBreachEvent(breach)
	
	// Notify DPO
	gpcs.dataProtectionOfficer.notifyOfBreach(breach)
	
	return breach, nil
}

// Privacy by design implementation

func (pde *PrivacyByDesignEngine) enforcePrivacyByDesign(operation DataOperation) (*PrivacyAssessment, error) {
	assessment := &PrivacyAssessment{
		OperationID:  operation.ID,
		AssessedAt:   time.Now(),
		Principles:   []PrincipleAssessment{},
		OverallScore: 100,
	}
	
	// Check proactive not reactive
	proactiveCheck := pde.checkProactiveMeasures(operation)
	assessment.Principles = append(assessment.Principles, proactiveCheck)
	assessment.OverallScore = min(assessment.OverallScore, proactiveCheck.Score)
	
	// Check privacy as default
	defaultCheck := pde.checkPrivacyDefaults(operation)
	assessment.Principles = append(assessment.Principles, defaultCheck)
	assessment.OverallScore = min(assessment.OverallScore, defaultCheck.Score)
	
	// Check full functionality
	functionalityCheck := pde.checkFullFunctionality(operation)
	assessment.Principles = append(assessment.Principles, functionalityCheck)
	assessment.OverallScore = min(assessment.OverallScore, functionalityCheck.Score)
	
	// Check end-to-end security
	securityCheck := pde.checkEndToEndSecurity(operation)
	assessment.Principles = append(assessment.Principles, securityCheck)
	assessment.OverallScore = min(assessment.OverallScore, securityCheck.Score)
	
	// Check visibility and transparency
	transparencyCheck := pde.checkTransparency(operation)
	assessment.Principles = append(assessment.Principles, transparencyCheck)
	assessment.OverallScore = min(assessment.OverallScore, transparencyCheck.Score)
	
	// Check respect for user privacy
	respectCheck := pde.checkUserPrivacyRespect(operation)
	assessment.Principles = append(assessment.Principles, respectCheck)
	assessment.OverallScore = min(assessment.OverallScore, respectCheck.Score)
	
	// Check privacy embedded into design
	embeddedCheck := pde.checkPrivacyEmbedded(operation)
	assessment.Principles = append(assessment.Principles, embeddedCheck)
	assessment.OverallScore = min(assessment.OverallScore, embeddedCheck.Score)
	
	// Require minimum score
	if assessment.OverallScore < 70 {
		return assessment, fmt.Errorf("operation does not meet privacy by design standards: score %d", assessment.OverallScore)
	}
	
	return assessment, nil
}

// Data minimization and retention

func (dme *DataMinimizationEngine) minimizeDataCollection(dataRequest DataCollectionRequest) (*MinimizedDataSet, error) {
	minimized := &MinimizedDataSet{
		RequestID:        dataRequest.ID,
		OriginalFields:   dataRequest.RequestedFields,
		MinimizedFields:  []DataField{},
		RemovedFields:    []DataField{},
		Justifications:   map[string]string{},
	}
	
	// Analyze each requested field
	for _, field := range dataRequest.RequestedFields {
		necessity := dme.assessFieldNecessity(field, dataRequest.Purpose)
		
		if necessity.IsNecessary {
			// Check if we can collect less precise data
			if alternative := dme.findLessIntrusiveAlternative(field); alternative != nil {
				minimized.MinimizedFields = append(minimized.MinimizedFields, *alternative)
				minimized.Justifications[field.Name] = fmt.Sprintf("Replaced with %s for privacy", alternative.Name)
			} else {
				minimized.MinimizedFields = append(minimized.MinimizedFields, field)
				minimized.Justifications[field.Name] = necessity.Justification
			}
		} else {
			minimized.RemovedFields = append(minimized.RemovedFields, field)
			minimized.Justifications[field.Name] = "Not necessary for stated purpose"
		}
	}
	
	return minimized, nil
}

// Encryption and pseudonymization

func (es *DataEncryptionService) encryptPersonalData(data []byte, userID string) ([]byte, error) {
	// Generate unique key for user
	key := es.deriveKeyForUser(userID)
	
	// Create cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}
	
	// GCM mode for authenticated encryption
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}
	
	// Generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}
	
	// Encrypt data
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	
	return ciphertext, nil
}

func (ps *PseudonymizationService) pseudonymizeData(data PersonalData) (*PseudonymizedData, error) {
	// Generate pseudonym
	pseudonym := ps.generatePseudonym(data.UserID)
	
	// Create mapping
	mapping := &PseudonymMapping{
		Original:  data.UserID,
		Pseudonym: pseudonym,
		CreatedAt: time.Now(),
		ValidUntil: time.Now().Add(ps.getRetentionPeriod(data.Category)),
	}
	
	// Store mapping securely
	ps.storePseudonymMapping(mapping)
	
	// Replace identifiers
	pseudonymized := &PseudonymizedData{
		Pseudonym:   pseudonym,
		Data:        ps.replaceIdentifiers(data.Data, mapping),
		Category:    data.Category,
		ProcessedAt: time.Now(),
	}
	
	return pseudonymized, nil
}

// Helper types

type ConsentRequest struct {
	UserID            string
	ConsentType       ConsentType
	Purposes          []ProcessingPurpose
	DataCategories    []DataCategory
	Method            ConsentMethod
	IPAddress         string
	UserAgent         string
	LegalBasis        LegalBasis
	ThirdPartySharing []ThirdPartyConsent
	Preferences       map[string]bool
}

type DataSubjectRequest struct {
	UserID                  string
	RequestType             RightsRequestType
	Description             string
	PreferredResponseMethod ResponseMethod
	AdditionalData          map[string]interface{}
	VerificationDocuments   []Document
}

type RightsResponse struct {
	ResponseID    string
	RequestID     string
	ResponseType  ResponseType
	Success       bool
	Data          interface{}
	ErasedData    []DataCategory
	RetainedData  []DataCategory
	Reason        string
	LegalBasis    []LegalBasis
	GeneratedAt   time.Time
	ValidUntil    time.Time
}

type BreachReport struct {
	DetectedAt        time.Time
	BreachType        BreachType
	Source            string
	AttackVector      string
	AffectedDataTypes []DataCategory
	InitialAssessment string
	ReportedBy        string
}

type ConsentProof struct {
	Timestamp      time.Time
	Hash           string
	Signature      string
	ProofData      map[string]interface{}
}

type ChildConsentInfo struct {
	ChildAge           int
	ParentID           string
	VerificationMethod string
	VerifiedAt         time.Time
}

type ThirdPartyConsent struct {
	ThirdPartyID   string
	ThirdPartyName string
	Purposes       []string
	DataShared     []DataCategory
	Consented      bool
}

type SpecialCategoryConsent struct {
	Category    string
	Consented   bool
	ExplicitConsent bool
}

type PersonalData struct {
	UserID   string
	Category DataCategory
	Data     map[string]interface{}
}

type PseudonymizedData struct {
	Pseudonym   string
	Data        map[string]interface{}
	Category    DataCategory
	ProcessedAt time.Time
}

type PseudonymMapping struct {
	Original   string
	Pseudonym  string
	CreatedAt  time.Time
	ValidUntil time.Time
}

type PrivacyAssessment struct {
	OperationID  string
	AssessedAt   time.Time
	Principles   []PrincipleAssessment
	OverallScore int
	Recommendations []string
}

type PrincipleAssessment struct {
	Principle   string
	Score       int
	Met         bool
	Findings    []string
}

type DataOperation struct {
	ID          string
	Type        string
	Purpose     string
	DataTypes   []DataCategory
	Processing  []ProcessingActivity
}

type MinimizedDataSet struct {
	RequestID       string
	OriginalFields  []DataField
	MinimizedFields []DataField
	RemovedFields   []DataField
	Justifications  map[string]string
}

type DataField struct {
	Name         string
	Type         string
	Required     bool
	Sensitivity  SensitivityLevel
}

type DataCollectionRequest struct {
	ID              string
	Purpose         ProcessingPurpose
	RequestedFields []DataField
	LegalBasis      LegalBasis
}

// Enums
type ResponseType int
type SensitivityLevel int

const (
	AccessResponse ResponseType = iota
	RectificationResponse
	ErasureResponse
	RestrictionResponse
	PortabilityResponse
	ObjectionResponse
	
	PublicData SensitivityLevel = iota
	InternalData
	ConfidentialData
	RestrictedData
	HighlyRestrictedData
)

// Utility functions

func (gpcs *GDPRPrivacyComplianceSystem) calculateDeadline(requestType RightsRequestType) time.Time {
	// GDPR requires response within 1 month, with possible 2-month extension
	baseDeadline := 30 * 24 * time.Hour
	
	// Urgent requests get priority
	if requestType == ErasureRequest || requestType == RestrictProcessingRequest {
		baseDeadline = 7 * 24 * time.Hour
	}
	
	return time.Now().Add(baseDeadline)
}

func (gpcs *GDPRPrivacyComplianceSystem) assessBreachSeverity(report BreachReport) BreachSeverity {
	severity := LowSeverity
	
	// Check data types
	for _, dataType := range report.AffectedDataTypes {
		if isSpecialCategory(dataType) {
			severity = max(severity, HighSeverity)
		} else if isSensitive(dataType) {
			severity = max(severity, MediumSeverity)
		}
	}
	
	// Check attack vector
	if report.AttackVector == "external_malicious" {
		severity = max(severity, HighSeverity)
	}
	
	// Check if encryption was compromised
	if strings.Contains(report.InitialAssessment, "encryption_compromised") {
		severity = CriticalSeverity
	}
	
	return severity
}

func (gpcs *GDPRPrivacyComplianceSystem) generateConsentProof(consent *ConsentRecord) ConsentProof {
	// Create proof data
	proofData := map[string]interface{}{
		"consent_id":    consent.ConsentID,
		"user_id":       consent.UserID,
		"purposes":      consent.Purposes,
		"granted_at":    consent.GrantedAt,
		"version":       consent.ConsentVersion,
	}
	
	// Generate hash
	jsonData, _ := json.Marshal(proofData)
	hash := sha256.Sum256(jsonData)
	
	// Create proof
	return ConsentProof{
		Timestamp: time.Now(),
		Hash:      base64.StdEncoding.EncodeToString(hash[:]),
		Signature: gpcs.signData(hash[:]),
		ProofData: proofData,
	}
}

func isSpecialCategory(category DataCategory) bool {
	// Special categories under GDPR Article 9
	specialCategories := []DataCategory{
		HealthData,
		BiometricData,
		GeneticData,
		ReligiousBeliefs,
		PoliticalOpinions,
		TradeUnionMembership,
		SexualOrientation,
		RacialEthnicOrigin,
	}
	
	for _, special := range specialCategories {
		if category == special {
			return true
		}
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b BreachSeverity) BreachSeverity {
	if a > b {
		return a
	}
	return b
}