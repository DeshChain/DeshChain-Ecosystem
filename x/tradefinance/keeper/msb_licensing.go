package keeper

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MSBLicensingSystem manages Money Service Business licensing
type MSBLicensingSystem struct {
	keeper                   Keeper
	applicationManager       *MSBApplicationManager
	licenseManager           *MSBLicenseManager
	complianceMonitor        *MSBComplianceMonitor
	renewalSystem            *LicenseRenewalSystem
	reportingEngine          *MSBReportingEngine
	bondManager              *SuretyBondManager
	mu                       sync.RWMutex
}

// MSBApplicationManager handles license applications
type MSBApplicationManager struct {
	applications             map[string]*MSBApplication
	workflowEngine           *ApplicationWorkflowEngine
	documentVerifier         *DocumentVerificationSystem
	backgroundChecker        *BackgroundCheckSystem
	financialAnalyzer        *FinancialAnalysisEngine
	riskAssessor             *ApplicationRiskAssessor
}

// MSBApplication represents a license application
type MSBApplication struct {
	ApplicationID            string
	ApplicantInfo            ApplicantInformation
	BusinessInfo             BusinessInformation
	FinancialInfo            FinancialInformation
	ComplianceProgram        MSBComplianceProgram
	OwnershipStructure       OwnershipStructure
	OperationalPlan          OperationalPlan
	RiskManagement           RiskManagementPlan
	Status                   ApplicationStatus
	SubmittedAt              time.Time
	LastUpdated              time.Time
	AssignedOfficer          string
	ReviewNotes              []ReviewNote
	RequiredDocuments        []RequiredDocument
	SubmittedDocuments       []SubmittedDocument
	BackgroundCheckResults   *BackgroundCheckReport
	FinancialAnalysisResults *FinancialAnalysisReport
	RiskAssessmentResults    *RiskAssessmentReport
	Deficiencies             []Deficiency
	ApprovalConditions       []ApprovalCondition
}

// MSBLicenseManager manages active licenses
type MSBLicenseManager struct {
	licenses                 map[string]*MSBLicense
	stateRegistry            *StateRegistryManager
	federalRegistry          *FederalRegistryManager
	activityTracker          *LicensedActivityTracker
	restrictionEnforcer      *RestrictionEnforcer
	publicRegistry           *PublicLicenseRegistry
}

// MSBLicense represents an active MSB license
type MSBLicense struct {
	LicenseID                string
	LicenseNumber            string
	LicenseType              MSBLicenseType
	IssuingState             string
	IssuingAuthority         string
	LicenseeName             string
	DBA                      []string // Doing Business As names
	IssueDate                time.Time
	ExpiryDate               time.Time
	Status                   MSBLicenseStatus
	AuthorizedActivities     []AuthorizedActivity
	Restrictions             []LicenseRestriction
	Locations                []LicensedLocation
	AgentLocations           []AgentLocation
	SuretyBond               *SuretyBondInfo
	NetWorthRequirement      sdk.Int
	CurrentNetWorth          sdk.Int
	ComplianceOfficer        ComplianceOfficerInfo
	AMLProgram               *AMLProgramInfo
	LastExamination          *ExaminationRecord
	EnforcementActions       []EnforcementAction
	Amendments               []LicenseAmendment
}

// MSBComplianceMonitor ensures ongoing compliance
type MSBComplianceMonitor struct {
	complianceChecks         map[string]ComplianceCheck
	monitoringSchedule       *MonitoringSchedule
	violationDetector        *ViolationDetector
	remediationTracker       *RemediationTracker
	examPreparation          *ExaminationPreparationSystem
	selfAssessmentTool       *SelfAssessmentTool
}

// LicenseRenewalSystem manages license renewals
type LicenseRenewalSystem struct {
	renewalQueue             *RenewalQueue
	renewalProcessor         *RenewalProcessor
	earlyWarningSystem       *RenewalWarningSystem
	autoRenewalManager       *AutoRenewalManager
	renewalDocumentTracker   *RenewalDocumentTracker
}

// MSBReportingEngine handles regulatory reporting
type MSBReportingEngine struct {
	reportGenerators         map[string]MSBReportGenerator
	filingSystem             *RegulatoryFilingSystem
	scheduleTracker          *ReportingScheduleTracker
	dataAggregator           *ReportDataAggregator
	validationEngine         *ReportValidationEngine
}

// SuretyBondManager manages surety bonds
type SuretyBondManager struct {
	bonds                    map[string]*SuretyBond
	bondProviders            map[string]*BondProvider
	claimProcessor           *BondClaimProcessor
	bondCalculator           *BondAmountCalculator
	renewalManager           *BondRenewalManager
}

// Types and enums
type ApplicationStatus int
type MSBLicenseType int
type MSBLicenseStatus int
type AuthorizedActivityType int
type ComplianceCheckType int
type ReportType int
type BondStatus int

const (
	// Application Status
	ApplicationDraft ApplicationStatus = iota
	ApplicationSubmitted
	ApplicationUnderReview
	ApplicationPendingDocuments
	ApplicationPendingBackgroundCheck
	ApplicationPendingFinancialReview
	ApplicationApproved
	ApplicationDenied
	ApplicationWithdrawn
	
	// MSB License Types
	MoneyTransmitter MSBLicenseType = iota
	CheckCasher
	CurrencyExchanger
	CheckSeller
	PrepaidAccessProvider
	
	// License Status
	LicenseActive MSBLicenseStatus = iota
	LicenseExpired
	LicenseSuspended
	LicenseRevoked
	LicenseInactive
	
	// Compliance Check Types
	NetWorthCheck ComplianceCheckType = iota
	BondVerification
	AMLProgramReview
	LocationCompliance
	ActivityCompliance
)

// Core MSB licensing methods

// SubmitMSBApplication submits a new MSB license application
func (k Keeper) SubmitMSBApplication(ctx context.Context, appRequest MSBApplicationRequest) (*MSBApplication, error) {
	mls := k.getMSBLicensingSystem()
	
	// Validate application request
	if err := mls.validateApplicationRequest(appRequest); err != nil {
		return nil, fmt.Errorf("invalid application request: %w", err)
	}
	
	// Create application
	application := &MSBApplication{
		ApplicationID:      generateID("MSBAPP"),
		ApplicantInfo:      appRequest.ApplicantInfo,
		BusinessInfo:       appRequest.BusinessInfo,
		FinancialInfo:      appRequest.FinancialInfo,
		ComplianceProgram:  appRequest.ComplianceProgram,
		OwnershipStructure: appRequest.OwnershipStructure,
		OperationalPlan:    appRequest.OperationalPlan,
		RiskManagement:     appRequest.RiskManagement,
		Status:             ApplicationSubmitted,
		SubmittedAt:        time.Now(),
		LastUpdated:        time.Now(),
		RequiredDocuments:  mls.getRequiredDocuments(appRequest.LicenseType, appRequest.States),
	}
	
	// Start application workflow
	workflow := mls.applicationManager.workflowEngine.createWorkflow(application)
	
	// Initial document verification
	docResults := mls.applicationManager.documentVerifier.verifyInitialDocuments(appRequest.Documents)
	application.SubmittedDocuments = docResults.VerifiedDocuments
	
	if len(docResults.MissingDocuments) > 0 {
		application.Status = ApplicationPendingDocuments
		application.Deficiencies = append(application.Deficiencies, Deficiency{
			Type:        DocumentDeficiency,
			Description: fmt.Sprintf("Missing documents: %v", docResults.MissingDocuments),
			DateAdded:   time.Now(),
		})
	}
	
	// Store application
	mls.applicationManager.applications[application.ApplicationID] = application
	if err := k.storeMSBApplication(ctx, application); err != nil {
		return nil, fmt.Errorf("failed to store application: %w", err)
	}
	
	// Start background checks if initial documents are complete
	if application.Status == ApplicationSubmitted {
		go mls.startBackgroundChecks(ctx, application)
	}
	
	// Send confirmation
	mls.sendApplicationConfirmation(application)
	
	return application, nil
}

// ProcessMSBApplication processes an MSB application
func (k Keeper) ProcessMSBApplication(ctx context.Context, applicationID string, action ApplicationAction) (*ApplicationResult, error) {
	mls := k.getMSBLicensingSystem()
	
	// Get application
	application, exists := mls.applicationManager.applications[applicationID]
	if !exists {
		return nil, fmt.Errorf("application not found")
	}
	
	// Validate action
	if err := mls.validateApplicationAction(application, action); err != nil {
		return nil, fmt.Errorf("invalid action: %w", err)
	}
	
	result := &ApplicationResult{
		ApplicationID: applicationID,
		ActionTaken:   action,
		Timestamp:     time.Now(),
	}
	
	switch action.Type {
	case ReviewAction:
		result = mls.processReview(ctx, application, action)
		
	case ApprovalAction:
		// Final compliance checks
		complianceResult := mls.performFinalComplianceChecks(application)
		if !complianceResult.Passed {
			result.Success = false
			result.Reason = "Failed final compliance checks"
			result.Details = complianceResult.Issues
			break
		}
		
		// Create license
		license, err := mls.createMSBLicense(ctx, application)
		if err != nil {
			result.Success = false
			result.Reason = err.Error()
			break
		}
		
		application.Status = ApplicationApproved
		result.Success = true
		result.LicenseID = license.LicenseID
		result.LicenseNumber = license.LicenseNumber
		
	case DenialAction:
		application.Status = ApplicationDenied
		result.Success = true
		result.Reason = action.Reason
		
		// Record denial reasons
		mls.recordDenialReasons(application, action)
		
	case RequestInfoAction:
		// Add deficiencies
		for _, deficiency := range action.Deficiencies {
			application.Deficiencies = append(application.Deficiencies, deficiency)
		}
		application.Status = ApplicationPendingDocuments
		result.Success = true
		
		// Notify applicant
		mls.notifyDeficiencies(application)
	}
	
	// Update application
	application.LastUpdated = time.Now()
	if err := k.updateMSBApplication(ctx, application); err != nil {
		return nil, fmt.Errorf("failed to update application: %w", err)
	}
	
	// Log action
	mls.logApplicationAction(application, action, result)
	
	return result, nil
}

// License management methods

func (mls *MSBLicensingSystem) createMSBLicense(ctx context.Context, application *MSBApplication) (*MSBLicense, error) {
	// Generate license number
	licenseNumber := mls.generateLicenseNumber(application.BusinessInfo.State, application.BusinessInfo.LicenseType)
	
	// Calculate surety bond requirement
	bondAmount := mls.bondManager.bondCalculator.calculateRequiredBond(
		application.BusinessInfo.State,
		application.BusinessInfo.LicenseType,
		application.FinancialInfo.ProjectedVolume,
	)
	
	// Create license
	license := &MSBLicense{
		LicenseID:            generateID("MSBL"),
		LicenseNumber:        licenseNumber,
		LicenseType:          application.BusinessInfo.LicenseType,
		IssuingState:         application.BusinessInfo.State,
		IssuingAuthority:     mls.getIssuingAuthority(application.BusinessInfo.State),
		LicenseeName:         application.ApplicantInfo.LegalName,
		DBA:                  application.ApplicantInfo.DBANames,
		IssueDate:            time.Now(),
		ExpiryDate:           mls.calculateExpiryDate(application.BusinessInfo.State),
		Status:               LicenseActive,
		AuthorizedActivities: mls.getAuthorizedActivities(application),
		Restrictions:         mls.getInitialRestrictions(application),
		Locations:            application.OperationalPlan.Locations,
		AgentLocations:       application.OperationalPlan.AgentLocations,
		NetWorthRequirement:  mls.getNetWorthRequirement(application),
		CurrentNetWorth:      application.FinancialInfo.NetWorth,
		ComplianceOfficer:    application.ComplianceProgram.ComplianceOfficer,
		AMLProgram: &AMLProgramInfo{
			ProgramID:        generateID("AML"),
			LastReview:       time.Now(),
			NextReview:       time.Now().Add(365 * 24 * time.Hour),
			RiskAssessment:   application.RiskManagement.AMLRiskAssessment,
			TrainingProgram:  application.ComplianceProgram.TrainingProgram,
		},
	}
	
	// Register with state
	if err := mls.licenseManager.stateRegistry.registerLicense(license); err != nil {
		return nil, fmt.Errorf("state registration failed: %w", err)
	}
	
	// Register with FinCEN (federal)
	if err := mls.licenseManager.federalRegistry.registerMSB(license); err != nil {
		return nil, fmt.Errorf("federal registration failed: %w", err)
	}
	
	// Store license
	mls.licenseManager.licenses[license.LicenseID] = license
	
	// Set up compliance monitoring
	mls.complianceMonitor.setupMonitoring(license)
	
	// Add to public registry
	mls.licenseManager.publicRegistry.addLicense(license)
	
	return license, nil
}

// Compliance monitoring methods

func (cm *MSBComplianceMonitor) performComplianceCheck(ctx context.Context, licenseID string) (*ComplianceCheckResult, error) {
	license, exists := cm.getLicense(licenseID)
	if !exists {
		return nil, fmt.Errorf("license not found")
	}
	
	result := &ComplianceCheckResult{
		CheckID:       generateID("COMP"),
		LicenseID:     licenseID,
		CheckDate:     time.Now(),
		OverallStatus: CompliantStatus,
		Checks:        []IndividualCheck{},
	}
	
	// Net worth check
	netWorthCheck := cm.checkNetWorth(license)
	result.Checks = append(result.Checks, netWorthCheck)
	if !netWorthCheck.Passed {
		result.OverallStatus = NonCompliantStatus
		result.RequiredActions = append(result.RequiredActions, ComplianceAction{
			Type:        CapitalInfusion,
			Description: fmt.Sprintf("Increase net worth by %s", license.NetWorthRequirement.Sub(license.CurrentNetWorth)),
			Deadline:    time.Now().Add(30 * 24 * time.Hour),
		})
	}
	
	// Surety bond verification
	bondCheck := cm.checkSuretyBond(license)
	result.Checks = append(result.Checks, bondCheck)
	if !bondCheck.Passed {
		result.OverallStatus = NonCompliantStatus
		result.RequiredActions = append(result.RequiredActions, ComplianceAction{
			Type:        BondRenewal,
			Description: "Renew or increase surety bond",
			Deadline:    time.Now().Add(15 * 24 * time.Hour),
		})
	}
	
	// AML program review
	amlCheck := cm.checkAMLProgram(license)
	result.Checks = append(result.Checks, amlCheck)
	if !amlCheck.Passed {
		result.OverallStatus = min(result.OverallStatus, ConditionallyCompliantStatus)
		result.RequiredActions = append(result.RequiredActions, amlCheck.RequiredActions...)
	}
	
	// Location compliance
	locationCheck := cm.checkLocationCompliance(license)
	result.Checks = append(result.Checks, locationCheck)
	
	// Activity compliance
	activityCheck := cm.checkActivityCompliance(license)
	result.Checks = append(result.Checks, activityCheck)
	
	// Store result
	cm.storeComplianceCheck(result)
	
	// Trigger remediation if needed
	if result.OverallStatus != CompliantStatus {
		cm.remediationTracker.createRemediationPlan(license, result)
	}
	
	return result, nil
}

// Renewal system methods

func (lrs *LicenseRenewalSystem) processRenewal(ctx context.Context, licenseID string) (*RenewalResult, error) {
	// Get license
	license := lrs.getLicense(licenseID)
	if license == nil {
		return nil, fmt.Errorf("license not found")
	}
	
	// Check renewal eligibility
	eligibility := lrs.checkRenewalEligibility(license)
	if !eligibility.Eligible {
		return &RenewalResult{
			LicenseID: licenseID,
			Success:   false,
			Reason:    eligibility.Reason,
		}, nil
	}
	
	// Create renewal application
	renewal := &RenewalApplication{
		RenewalID:     generateID("REN"),
		LicenseID:     licenseID,
		SubmittedAt:   time.Now(),
		CurrentExpiry: license.ExpiryDate,
		Status:        RenewalPending,
	}
	
	// Collect renewal requirements
	requirements := lrs.getRenewalRequirements(license)
	renewal.Requirements = requirements
	
	// Process renewal documents
	if err := lrs.renewalDocumentTracker.processDocuments(renewal); err != nil {
		renewal.Status = RenewalPendingDocuments
		return &RenewalResult{
			LicenseID:  licenseID,
			Success:    false,
			Reason:     "Missing renewal documents",
			RenewalID:  renewal.RenewalID,
		}, nil
	}
	
	// Update license
	license.ExpiryDate = license.ExpiryDate.Add(365 * 24 * time.Hour) // 1 year extension
	renewal.NewExpiry = license.ExpiryDate
	renewal.Status = RenewalApproved
	renewal.ApprovedAt = timePtr(time.Now())
	
	// Update registries
	lrs.updateRegistries(license)
	
	// Store renewal
	lrs.storeRenewal(renewal)
	
	return &RenewalResult{
		LicenseID:  licenseID,
		Success:    true,
		RenewalID:  renewal.RenewalID,
		NewExpiry:  license.ExpiryDate,
	}, nil
}

// Reporting methods

func (mre *MSBReportingEngine) generateMSBReport(ctx context.Context, reportRequest MSBReportRequest) (*MSBReport, error) {
	// Get report generator
	generator, exists := mre.reportGenerators[reportRequest.ReportType]
	if !exists {
		return nil, fmt.Errorf("unknown report type: %s", reportRequest.ReportType)
	}
	
	// Aggregate data
	data, err := mre.dataAggregator.aggregateData(reportRequest)
	if err != nil {
		return nil, fmt.Errorf("data aggregation failed: %w", err)
	}
	
	// Generate report
	report := &MSBReport{
		ReportID:    generateID("MSBREP"),
		ReportType:  reportRequest.ReportType,
		PeriodStart: reportRequest.PeriodStart,
		PeriodEnd:   reportRequest.PeriodEnd,
		GeneratedAt: time.Now(),
		LicenseID:   reportRequest.LicenseID,
	}
	
	// Generate content
	content, err := generator.Generate(data)
	if err != nil {
		return nil, fmt.Errorf("report generation failed: %w", err)
	}
	report.Content = content
	
	// Validate report
	validation := mre.validationEngine.validateReport(report)
	if !validation.Valid {
		report.ValidationErrors = validation.Errors
		return report, fmt.Errorf("report validation failed")
	}
	
	// File report if requested
	if reportRequest.AutoFile {
		filing, err := mre.filingSystem.fileReport(report)
		if err != nil {
			return report, fmt.Errorf("report filing failed: %w", err)
		}
		report.FilingReference = filing.Reference
		report.FiledAt = timePtr(filing.FiledAt)
	}
	
	return report, nil
}

// Helper types

type MSBApplicationRequest struct {
	ApplicantInfo      ApplicantInformation
	BusinessInfo       BusinessInformation
	FinancialInfo      FinancialInformation
	ComplianceProgram  MSBComplianceProgram
	OwnershipStructure OwnershipStructure
	OperationalPlan    OperationalPlan
	RiskManagement     RiskManagementPlan
	LicenseType        MSBLicenseType
	States             []string
	Documents          []Document
}

type ApplicantInformation struct {
	LegalName          string
	DBANames           []string
	TaxID              string
	FormationType      string
	DateOfFormation    time.Time
	StateOfFormation   string
	PrincipalAddress   Address
	MailingAddress     Address
	Phone              string
	Email              string
	Website            string
}

type BusinessInformation struct {
	LicenseType           MSBLicenseType
	State                 string
	BusinessActivities    []string
	CustomerTypes         []string
	TransactionTypes      []string
	MonthlyVolume         sdk.Int
	AverageTransactionSize sdk.Int
	MaxTransactionSize    sdk.Int
	Countries             []string
}

type FinancialInformation struct {
	NetWorth              sdk.Int
	LiquidAssets          sdk.Int
	TotalAssets           sdk.Int
	TotalLiabilities      sdk.Int
	AnnualRevenue         sdk.Int
	ProjectedVolume       sdk.Int
	BankAccounts          []BankAccount
	FundingSources        []FundingSource
	FinancialStatements   []FinancialStatement
	AuditorInfo           *AuditorInformation
}

type MSBComplianceProgram struct {
	ComplianceOfficer     ComplianceOfficerInfo
	AMLPolicy             Document
	RiskAssessment        Document
	CustomerIDProgram     Document
	SuspiciousActivityProc Document
	TrainingProgram       TrainingProgramInfo
	IndependentReview     *IndependentReviewInfo
	VendorManagement      *VendorManagementInfo
}

type OwnershipStructure struct {
	Owners              []OwnerInfo
	Directors           []DirectorInfo
	Officers            []OfficerInfo
	BeneficialOwners    []BeneficialOwnerInfo
	CorporateStructure  Document
	OwnershipChart      Document
}

type AuthorizedActivity struct {
	ActivityType        AuthorizedActivityType
	Description         string
	VolumeLimit         *sdk.Int
	TransactionLimit    *sdk.Int
	GeographicScope     []string
	Restrictions        []string
}

type LicenseRestriction struct {
	RestrictionID   string
	Type            string
	Description     string
	EffectiveDate   time.Time
	ExpiryDate      *time.Time
	Conditions      []string
}

type ComplianceCheckResult struct {
	CheckID          string
	LicenseID        string
	CheckDate        time.Time
	OverallStatus    ComplianceStatus
	Checks           []IndividualCheck
	RequiredActions  []ComplianceAction
	NextCheckDate    time.Time
}

type IndividualCheck struct {
	CheckType       ComplianceCheckType
	CheckName       string
	Passed          bool
	Details         string
	Evidence        []string
	RequiredActions []ComplianceAction
}

type ComplianceAction struct {
	Type        ActionType
	Description string
	Deadline    time.Time
	Priority    Priority
}

type RenewalApplication struct {
	RenewalID     string
	LicenseID     string
	SubmittedAt   time.Time
	CurrentExpiry time.Time
	NewExpiry     time.Time
	Status        RenewalStatus
	Requirements  []RenewalRequirement
	Documents     []Document
	ApprovedAt    *time.Time
	ApprovedBy    string
}

type MSBReport struct {
	ReportID         string
	ReportType       string
	LicenseID        string
	PeriodStart      time.Time
	PeriodEnd        time.Time
	GeneratedAt      time.Time
	Content          []byte
	Format           string
	ValidationErrors []string
	FilingReference  string
	FiledAt          *time.Time
}

// Enums for status types
type ComplianceStatus int
type ActionType int
type Priority int
type RenewalStatus int

const (
	CompliantStatus ComplianceStatus = iota
	ConditionallyCompliantStatus
	NonCompliantStatus
	
	CapitalInfusion ActionType = iota
	BondRenewal
	PolicyUpdate
	TrainingRequired
	
	HighPriority Priority = iota
	MediumPriority
	LowPriority
	
	RenewalPending RenewalStatus = iota
	RenewalPendingDocuments
	RenewalApproved
	RenewalDenied
)

// Utility functions

func (mls *MSBLicensingSystem) generateLicenseNumber(state string, licenseType MSBLicenseType) string {
	prefix := mls.getStatePrefix(state)
	typeCode := mls.getLicenseTypeCode(licenseType)
	timestamp := time.Now().Unix()
	
	return fmt.Sprintf("%s-%s-%d", prefix, typeCode, timestamp)
}

func (mls *MSBLicensingSystem) calculateExpiryDate(state string) time.Time {
	// Most states have 1-year licenses
	duration := 365 * 24 * time.Hour
	
	// Some states have different durations
	switch state {
	case "CA":
		duration = 2 * 365 * 24 * time.Hour // 2 years
	case "NY":
		duration = 2 * 365 * 24 * time.Hour // 2 years
	}
	
	return time.Now().Add(duration)
}

func timePtr(t time.Time) *time.Time {
	return &t
}