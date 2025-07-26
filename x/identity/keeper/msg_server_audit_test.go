package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/DeshChain/DeshChain-Ecosystem/x/identity/types"
)

type MsgServerAuditTestSuite struct {
	KeeperTestSuite
}

func TestMsgServerAuditTestSuite(t *testing.T) {
	suite.Run(t, new(MsgServerAuditTestSuite))
}

func (suite *MsgServerAuditTestSuite) TestLogAuditEvent() {
	ctx := suite.ctx
	msgServer := suite.msgServer

	// Test successful audit event logging
	req := &types.MsgLogAuditEvent{
		Authority:   suite.authority,
		EventType:   types.AuditEventType_IDENTITY_CREATED,
		Actor:       "desh1actor123",
		Subject:     "desh1subject456",
		Resource:    "identity_789",
		Action:      "create_identity",
		Outcome:     types.AuditOutcome_SUCCESS,
		Severity:    types.AuditSeverity_MEDIUM,
		Description: "Test identity creation",
		TechnicalDetails: map[string]interface{}{
			"method": "POST",
			"endpoint": "/api/identity/create",
		},
		ComplianceFlags: []types.ComplianceFlag{
			{
				Regulation:  "GDPR",
				Requirement: "Article 6 - Lawful basis",
				Status:      types.ComplianceStatus_COMPLIANT,
				Notes:       "User provided explicit consent",
			},
		},
		Metadata: map[string]interface{}{
			"source": "web_interface",
			"version": "1.0.0",
		},
		IPAddress:    "192.168.1.100",
		UserAgent:    "Mozilla/5.0 Test Browser",
		SessionID:    "session_12345",
		ModuleSource: "identity",
	}

	res, err := msgServer.LogAuditEvent(ctx, req)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().NotEmpty(res.EventId)

	// Verify event was logged
	event, found := suite.keeper.GetAuditEvent(suite.ctx, res.EventId)
	suite.Require().True(found)
	suite.Require().Equal(req.EventType, event.EventType)
	suite.Require().Equal(req.Actor, event.Actor)
	suite.Require().Equal(req.Subject, event.Subject)
	suite.Require().Equal(req.Resource, event.Resource)
	suite.Require().Equal(req.Action, event.Action)
	suite.Require().Equal(req.Outcome, event.Outcome)
	suite.Require().Equal(req.Severity, event.Severity)
	suite.Require().Equal(req.Description, event.Description)
	suite.Require().Equal(req.IPAddress, event.IPAddress)
	suite.Require().Equal(req.UserAgent, event.UserAgent)
	suite.Require().Equal(req.SessionID, event.SessionID)
	suite.Require().Equal(req.ModuleSource, event.ModuleSource)
}

func (suite *MsgServerAuditTestSuite) TestLogAuditEventUnauthorized() {
	ctx := suite.ctx
	msgServer := suite.msgServer

	// Test with unauthorized authority
	req := &types.MsgLogAuditEvent{
		Authority:   "desh1unauthorized",
		EventType:   types.AuditEventType_IDENTITY_CREATED,
		Actor:       "desh1actor123",
		Subject:     "desh1subject456",
		Resource:    "identity_789",
		Action:      "create_identity",
		Outcome:     types.AuditOutcome_SUCCESS,
		Severity:    types.AuditSeverity_MEDIUM,
		Description: "Test identity creation",
	}

	_, err := msgServer.LogAuditEvent(ctx, req)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "insufficient permissions")
}

func (suite *MsgServerAuditTestSuite) TestGenerateComplianceReport() {
	ctx := suite.ctx
	msgServer := suite.msgServer

	// First, create some audit events to include in the report
	events := []*types.AuditEvent{
		{
			EventID:      "report_event_001",
			Timestamp:    time.Now().Add(-24 * time.Hour),
			EventType:    types.AuditEventType_IDENTITY_CREATED,
			Actor:        "desh1actor1",
			Subject:      "desh1subject1",
			Resource:     "identity_1",
			Action:       "create",
			Outcome:      types.AuditOutcome_SUCCESS,
			Severity:     types.AuditSeverity_LOW,
			Description:  "Identity creation for report",
			ModuleSource: "identity",
			ChainHeight:  1000,
			ComplianceFlags: []types.ComplianceFlag{
				{
					Regulation:  "GDPR",
					Requirement: "Article 6 - Lawful basis",
					Status:      types.ComplianceStatus_COMPLIANT,
				},
			},
		},
		{
			EventID:      "report_event_002",
			Timestamp:    time.Now().Add(-12 * time.Hour),
			EventType:    types.AuditEventType_COMPLIANCE_VIOLATION,
			Actor:        "desh1actor2",
			Subject:      "desh1subject2",
			Resource:     "sensitive_data",
			Action:       "unauthorized_access",
			Outcome:      types.AuditOutcome_FAILURE,
			Severity:     types.AuditSeverity_CRITICAL,
			Description:  "Compliance violation for report",
			ModuleSource: "identity",
			ChainHeight:  1001,
			ComplianceFlags: []types.ComplianceFlag{
				{
					Regulation:  "GDPR",
					Requirement: "Article 32 - Security",
					Status:      types.ComplianceStatus_NON_COMPLIANT,
				},
			},
		},
	}

	// Store events
	for _, event := range events {
		err := suite.keeper.LogAuditEventToStore(suite.ctx, event)
		suite.Require().NoError(err)
	}

	// Generate compliance report
	req := &types.MsgGenerateComplianceReport{
		Authority:  suite.authority,
		ReportType: types.ComplianceReportType_GDPR_COMPLIANCE,
		TimeRange: types.AuditTimeRange{
			StartTime: time.Now().Add(-48 * time.Hour),
			EndTime:   time.Now(),
		},
		Scope: types.ComplianceScope{
			IncludeModules: []string{"identity"},
		},
		Regulations: []string{"GDPR"},
	}

	res, err := msgServer.GenerateComplianceReport(ctx, req)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().NotEmpty(res.ReportId)
	suite.Require().True(res.ComplianceScore >= 0)
	suite.Require().True(res.ComplianceScore <= 100)
	suite.Require().True(res.TotalFindings >= 0)

	// Verify report was stored
	report, found := suite.keeper.GetComplianceReport(suite.ctx, res.ReportId)
	suite.Require().True(found)
	suite.Require().Equal(req.ReportType, report.ReportType)
	suite.Require().Equal(req.TimeRange, report.TimeRange)
	suite.Require().Equal(req.Regulations, report.Regulations)
}

func (suite *MsgServerAuditTestSuite) TestProcessDataSubjectRequest() {
	ctx := suite.ctx
	msgServer := suite.msgServer

	// Test access request
	accessRequest := &types.DataSubjectRequest{
		RequestType:   types.DataSubjectRequestType_ACCESS,
		DataSubject:   "desh1subject123",
		RequestedBy:   "desh1subject123",
		Description:   "Request access to all personal data",
		RequestDetails: map[string]interface{}{
			"data_categories": []string{"personal", "financial"},
			"format":          "JSON",
		},
		Regulation: "GDPR",
		Priority:   types.DataSubjectRequestPriority_STANDARD,
	}

	req := &types.MsgProcessDataSubjectRequest{
		Authority: suite.authority,
		Request:   accessRequest,
	}

	res, err := msgServer.ProcessDataSubjectRequest(ctx, req)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().NotEmpty(res.RequestId)
	suite.Require().Equal(types.DataSubjectRequestStatus_RECEIVED, res.Status)
	suite.Require().True(res.DueDate.After(time.Now()))

	// Verify request was stored
	storedRequest, found := suite.keeper.GetDataSubjectRequest(suite.ctx, res.RequestId)
	suite.Require().True(found)
	suite.Require().Equal(accessRequest.RequestType, storedRequest.RequestType)
	suite.Require().Equal(accessRequest.DataSubject, storedRequest.DataSubject)
	suite.Require().Equal(accessRequest.Description, storedRequest.Description)
	suite.Require().Equal(accessRequest.Regulation, storedRequest.Regulation)

	// Test erasure request
	erasureRequest := &types.DataSubjectRequest{
		RequestType:   types.DataSubjectRequestType_ERASURE,
		DataSubject:   "desh1subject456",
		RequestedBy:   "desh1subject456",
		Description:   "Request deletion of all personal data",
		RequestDetails: map[string]interface{}{
			"reason": "withdrawing_consent",
		},
		Regulation: "GDPR",
		Priority:   types.DataSubjectRequestPriority_URGENT,
	}

	erasureReq := &types.MsgProcessDataSubjectRequest{
		Authority: suite.authority,
		Request:   erasureRequest,
	}

	erasureRes, err := msgServer.ProcessDataSubjectRequest(ctx, erasureReq)
	suite.Require().NoError(err)
	suite.Require().NotNil(erasureRes)
	suite.Require().NotEmpty(erasureRes.RequestId)
	suite.Require().Equal(types.DataSubjectRequestStatus_RECEIVED, erasureRes.Status)
}

func (suite *MsgServerAuditTestSuite) TestProcessDataSubjectRequestSelfService() {
	ctx := suite.ctx
	msgServer := suite.msgServer

	// Test that data subjects can make requests about themselves
	dataSubject := "desh1selfservice"
	selfRequest := &types.DataSubjectRequest{
		RequestType:   types.DataSubjectRequestType_PORTABILITY,
		DataSubject:   dataSubject,
		RequestedBy:   dataSubject, // Same as data subject
		Description:   "Self-service data portability request",
		Regulation:    "GDPR",
		Priority:      types.DataSubjectRequestPriority_STANDARD,
	}

	req := &types.MsgProcessDataSubjectRequest{
		Authority: dataSubject, // Data subject making their own request
		Request:   selfRequest,
	}

	res, err := msgServer.ProcessDataSubjectRequest(ctx, req)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().NotEmpty(res.RequestId)
}

func (suite *MsgServerAuditTestSuite) TestCreatePrivacyImpactAssessment() {
	ctx := suite.ctx
	msgServer := suite.msgServer

	assessment := &types.PrivacyImpactAssessment{
		Title:              "KYC Data Processing Assessment",
		Description:        "Assessment of KYC data processing activities",
		ProcessingActivity: "kyc_verification",
		DataController:     "DeshChain Foundation",
		DataProcessor:      "Third Party KYC Provider",
		LegalBasis:         []string{"Article 6(1)(c) - Legal obligation"},
		DataCategories:     []string{"identity", "financial", "biometric"},
		DataSubjects:       []string{"customers", "users"},
		ProcessingPurposes: []string{"identity_verification", "compliance"},
		RetentionPeriod:    "7 years",
		TechnicalMeasures:  []string{"encryption", "access_controls", "audit_logging"},
		OrganizationalMeasures: []string{"staff_training", "data_protection_policies", "incident_response"},
		ThirdPartySharing: []types.ThirdPartySharing{
			{
				PartyName:        "Regulatory Authority",
				PartyType:        "controller",
				Purpose:          "compliance_reporting",
				DataCategories:   []string{"identity", "financial"},
				LegalBasis:       "Legal obligation",
				Safeguards:       []string{"encrypted_transmission", "access_controls"},
				ContractualTerms: "Data sharing agreement",
				Country:          "India",
				AdequacyDecision: true,
				StartDate:        time.Now(),
			},
		},
	}

	req := &types.MsgCreatePrivacyImpactAssessment{
		Authority:  suite.authority,
		Assessment: assessment,
	}

	res, err := msgServer.CreatePrivacyImpactAssessment(ctx, req)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().NotEmpty(res.AssessmentId)
	suite.Require().True(res.RiskLevel >= types.RiskLevel_VERY_LOW)
	suite.Require().True(res.RiskLevel <= types.RiskLevel_VERY_HIGH)
	suite.Require().True(res.ReviewDate.After(time.Now()))

	// Verify assessment was stored
	storedAssessment, found := suite.keeper.GetPrivacyImpactAssessment(suite.ctx, res.AssessmentId)
	suite.Require().True(found)
	suite.Require().Equal(assessment.Title, storedAssessment.Title)
	suite.Require().Equal(assessment.ProcessingActivity, storedAssessment.ProcessingActivity)
	suite.Require().Equal(assessment.DataController, storedAssessment.DataController)
	suite.Require().Equal(assessment.DataProcessor, storedAssessment.DataProcessor)
	suite.Require().Equal(types.PIAStatus_DRAFT, storedAssessment.Status)
	suite.Require().Equal(suite.authority, storedAssessment.CreatedBy)

	// Verify risk assessment was performed
	suite.Require().NotEmpty(storedAssessment.RiskAssessment.IdentifiedRisks)
	suite.Require().True(len(storedAssessment.Mitigation) > 0)
}

func (suite *MsgServerAuditTestSuite) TestUpdateAuditSettings() {
	ctx := suite.ctx
	msgServer := suite.msgServer

	newSettings := types.AuditSettings{
		RetentionPeriodDays:         1825, // 5 years
		AutoPurgeEnabled:           false,
		ComplianceMonitoringEnabled: true,
		RealTimeAlertsEnabled:      false,
		AnonymizeExpiredData:       true,
	}

	req := &types.MsgUpdateAuditSettings{
		Authority: suite.authority,
		Settings:  newSettings,
	}

	res, err := msgServer.UpdateAuditSettings(ctx, req)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().True(res.Success)

	// Verify settings were updated
	retrievedSettings, found := suite.keeper.GetAuditSettings(suite.ctx)
	suite.Require().True(found)
	suite.Require().Equal(newSettings.RetentionPeriodDays, retrievedSettings.RetentionPeriodDays)
	suite.Require().Equal(newSettings.AutoPurgeEnabled, retrievedSettings.AutoPurgeEnabled)
	suite.Require().Equal(newSettings.ComplianceMonitoringEnabled, retrievedSettings.ComplianceMonitoringEnabled)
	suite.Require().Equal(newSettings.RealTimeAlertsEnabled, retrievedSettings.RealTimeAlertsEnabled)
	suite.Require().Equal(newSettings.AnonymizeExpiredData, retrievedSettings.AnonymizeExpiredData)
}

func (suite *MsgServerAuditTestSuite) TestUpdateAuditSettingsUnauthorized() {
	ctx := suite.ctx
	msgServer := suite.msgServer

	newSettings := types.AuditSettings{
		RetentionPeriodDays:         365,
		AutoPurgeEnabled:           true,
		ComplianceMonitoringEnabled: false,
		RealTimeAlertsEnabled:      false,
		AnonymizeExpiredData:       false,
	}

	req := &types.MsgUpdateAuditSettings{
		Authority: "desh1unauthorized",
		Settings:  newSettings,
	}

	_, err := msgServer.UpdateAuditSettings(ctx, req)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "insufficient permissions")
}

func (suite *MsgServerAuditTestSuite) TestComplianceViolationLogging() {
	ctx := suite.ctx
	msgServer := suite.msgServer

	// Log a compliance violation event
	req := &types.MsgLogAuditEvent{
		Authority:   suite.authority,
		EventType:   types.AuditEventType_COMPLIANCE_VIOLATION,
		Actor:       "desh1malicious",
		Subject:     "desh1victim",
		Resource:    "personal_data",
		Action:      "unauthorized_processing",
		Outcome:     types.AuditOutcome_FAILURE,
		Severity:    types.AuditSeverity_CRITICAL,
		Description: "Attempted unauthorized processing of personal data",
		ComplianceFlags: []types.ComplianceFlag{
			{
				Regulation:  "GDPR",
				Requirement: "Article 6 - Lawful basis",
				Status:      types.ComplianceStatus_NON_COMPLIANT,
				Notes:       "No legal basis for processing identified",
			},
		},
		ModuleSource: "identity",
	}

	res, err := msgServer.LogAuditEvent(ctx, req)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().NotEmpty(res.EventId)

	// Verify the violation was logged
	event, found := suite.keeper.GetAuditEvent(suite.ctx, res.EventId)
	suite.Require().True(found)
	suite.Require().Equal(types.AuditEventType_COMPLIANCE_VIOLATION, event.EventType)
	suite.Require().Equal(types.AuditSeverity_CRITICAL, event.Severity)
	suite.Require().Len(event.ComplianceFlags, 1)
	suite.Require().Equal("GDPR", event.ComplianceFlags[0].Regulation)
	suite.Require().Equal(types.ComplianceStatus_NON_COMPLIANT, event.ComplianceFlags[0].Status)
}

func (suite *MsgServerAuditTestSuite) TestAuditEventWithTechnicalDetails() {
	ctx := suite.ctx
	msgServer := suite.msgServer

	technicalDetails := map[string]interface{}{
		"http_method":     "POST",
		"endpoint":        "/api/identity/create",
		"response_code":   201,
		"response_time":   "150ms",
		"request_size":    1024,
		"response_size":   512,
		"user_agent":      "DeshChain Mobile App v1.2.3",
		"client_version":  "1.2.3",
		"api_version":     "v2",
		"correlation_id":  "corr_12345",
	}

	req := &types.MsgLogAuditEvent{
		Authority:        suite.authority,
		EventType:        types.AuditEventType_IDENTITY_CREATED,
		Actor:            "desh1actor123",
		Subject:          "desh1subject456",
		Resource:         "identity_789",
		Action:           "create_identity",
		Outcome:          types.AuditOutcome_SUCCESS,
		Severity:         types.AuditSeverity_MEDIUM,
		Description:      "Identity creation with technical details",
		TechnicalDetails: technicalDetails,
		IPAddress:        "10.0.1.100",
		UserAgent:        "DeshChain Mobile App v1.2.3",
		SessionID:        "session_abcdef",
		ModuleSource:     "identity",
	}

	res, err := msgServer.LogAuditEvent(ctx, req)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)

	// Verify technical details were stored
	event, found := suite.keeper.GetAuditEvent(suite.ctx, res.EventId)
	suite.Require().True(found)
	suite.Require().Equal(technicalDetails, event.TechnicalDetails)
	suite.Require().Equal("POST", event.TechnicalDetails["http_method"])
	suite.Require().Equal(float64(201), event.TechnicalDetails["response_code"])
	suite.Require().Equal("150ms", event.TechnicalDetails["response_time"])
}

func (suite *MsgServerAuditTestSuite) TestMultipleComplianceFlags() {
	ctx := suite.ctx
	msgServer := suite.msgServer

	complianceFlags := []types.ComplianceFlag{
		{
			Regulation:  "GDPR",
			Requirement: "Article 6 - Lawful basis",
			Status:      types.ComplianceStatus_COMPLIANT,
			Notes:       "User provided explicit consent",
			Metadata: map[string]interface{}{
				"consent_timestamp": time.Now().Format(time.RFC3339),
				"consent_method":    "web_form",
			},
		},
		{
			Regulation:  "GDPR",
			Requirement: "Article 13 - Information to be provided",
			Status:      types.ComplianceStatus_COMPLIANT,
			Notes:       "Privacy notice displayed and acknowledged",
			Metadata: map[string]interface{}{
				"privacy_notice_version": "2.1",
				"acknowledgment_method":  "checkbox",
			},
		},
		{
			Regulation:  "DPDP",
			Requirement: "Section 7 - Notice",
			Status:      types.ComplianceStatus_PENDING,
			Notes:       "DPDP notice compliance pending verification",
		},
	}

	req := &types.MsgLogAuditEvent{
		Authority:       suite.authority,
		EventType:       types.AuditEventType_CONSENT_GIVEN,
		Actor:           "desh1user123",
		Subject:         "desh1user123",
		Resource:        "consent_record_456",
		Action:          "provide_consent",
		Outcome:         types.AuditOutcome_SUCCESS,
		Severity:        types.AuditSeverity_MEDIUM,
		Description:     "User provided consent for data processing",
		ComplianceFlags: complianceFlags,
		ModuleSource:    "identity",
	}

	res, err := msgServer.LogAuditEvent(ctx, req)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)

	// Verify multiple compliance flags were stored
	event, found := suite.keeper.GetAuditEvent(suite.ctx, res.EventId)
	suite.Require().True(found)
	suite.Require().Len(event.ComplianceFlags, 3)

	// Check GDPR flags
	gdprFlags := 0
	for _, flag := range event.ComplianceFlags {
		if flag.Regulation == "GDPR" {
			gdprFlags++
			suite.Require().Equal(types.ComplianceStatus_COMPLIANT, flag.Status)
		}
	}
	suite.Require().Equal(2, gdprFlags)

	// Check DPDP flag
	dpdpFound := false
	for _, flag := range event.ComplianceFlags {
		if flag.Regulation == "DPDP" {
			dpdpFound = true
			suite.Require().Equal(types.ComplianceStatus_PENDING, flag.Status)
			break
		}
	}
	suite.Require().True(dpdpFound)
}