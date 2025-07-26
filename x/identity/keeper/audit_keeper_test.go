package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/namo/x/identity/keeper"
	"github.com/namo/x/identity/types"
)

type AuditKeeperTestSuite struct {
	KeeperTestSuite
}

func TestAuditKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(AuditKeeperTestSuite))
}

func (suite *AuditKeeperTestSuite) TestLogAuditEvent() {
	ctx := suite.ctx
	keeper := suite.keeper

	// Test successful audit event logging
	event := &types.AuditEvent{
		EventID:      "test_event_001",
		Timestamp:    time.Now(),
		EventType:    types.AuditEventType_IDENTITY_CREATED,
		Actor:        "desh1actor123",
		Subject:      "desh1subject456",
		Resource:     "identity_789",
		Action:       "create_identity",
		Outcome:      types.AuditOutcome_SUCCESS,
		Severity:     types.AuditSeverity_MEDIUM,
		Description:  "Test identity creation",
		ModuleSource: "identity",
		ChainHeight:  1000,
		Metadata: map[string]interface{}{
			"test_key": "test_value",
		},
	}

	err := keeper.LogAuditEventToStore(ctx, event)
	suite.Require().NoError(err)

	// Verify event was stored
	retrievedEvent, found := keeper.GetAuditEvent(ctx, event.EventID)
	suite.Require().True(found)
	suite.Require().Equal(event.EventID, retrievedEvent.EventID)
	suite.Require().Equal(event.EventType, retrievedEvent.EventType)
	suite.Require().Equal(event.Actor, retrievedEvent.Actor)
	suite.Require().Equal(event.Subject, retrievedEvent.Subject)
	suite.Require().Equal(event.Resource, retrievedEvent.Resource)
	suite.Require().Equal(event.Action, retrievedEvent.Action)
	suite.Require().Equal(event.Outcome, retrievedEvent.Outcome)
	suite.Require().Equal(event.Severity, retrievedEvent.Severity)
	suite.Require().Equal(event.Description, retrievedEvent.Description)
}

func (suite *AuditKeeperTestSuite) TestQueryAuditEvents() {
	ctx := suite.ctx
	keeper := suite.keeper

	// Create multiple audit events
	events := []*types.AuditEvent{
		{
			EventID:      "event_001",
			Timestamp:    time.Now().Add(-2 * time.Hour),
			EventType:    types.AuditEventType_IDENTITY_CREATED,
			Actor:        "desh1actor1",
			Subject:      "desh1subject1",
			Resource:     "identity_1",
			Action:       "create",
			Outcome:      types.AuditOutcome_SUCCESS,
			Severity:     types.AuditSeverity_LOW,
			Description:  "First test event",
			ModuleSource: "identity",
			ChainHeight:  1000,
		},
		{
			EventID:      "event_002",
			Timestamp:    time.Now().Add(-1 * time.Hour),
			EventType:    types.AuditEventType_CREDENTIAL_ISSUED,
			Actor:        "desh1actor2",
			Subject:      "desh1subject2",
			Resource:     "credential_2",
			Action:       "issue",
			Outcome:      types.AuditOutcome_SUCCESS,
			Severity:     types.AuditSeverity_HIGH,
			Description:  "Second test event",
			ModuleSource: "identity",
			ChainHeight:  1001,
		},
		{
			EventID:      "event_003",
			Timestamp:    time.Now(),
			EventType:    types.AuditEventType_COMPLIANCE_VIOLATION,
			Actor:        "desh1actor3",
			Subject:      "desh1subject3",
			Resource:     "violation_3",
			Action:       "violation",
			Outcome:      types.AuditOutcome_FAILURE,
			Severity:     types.AuditSeverity_CRITICAL,
			Description:  "Third test event",
			ModuleSource: "identity",
			ChainHeight:  1002,
		},
	}

	// Store all events
	for _, event := range events {
		err := keeper.LogAuditEventToStore(ctx, event)
		suite.Require().NoError(err)
	}

	// Test query all events
	allEvents, err := keeper.QueryAuditEvents(ctx, &types.QueryAuditEventsRequest{})
	suite.Require().NoError(err)
	suite.Require().Len(allEvents.Events, 3)

	// Test query by event type
	identityEvents, err := keeper.QueryAuditEvents(ctx, &types.QueryAuditEventsRequest{
		EventType: &types.AuditEventType_IDENTITY_CREATED,
	})
	suite.Require().NoError(err)
	suite.Require().Len(identityEvents.Events, 1)
	suite.Require().Equal("event_001", identityEvents.Events[0].EventID)

	// Test query by actor
	actorEvents, err := keeper.QueryAuditEvents(ctx, &types.QueryAuditEventsRequest{
		Actor: "desh1actor2",
	})
	suite.Require().NoError(err)
	suite.Require().Len(actorEvents.Events, 1)
	suite.Require().Equal("event_002", actorEvents.Events[0].EventID)

	// Test query by severity
	criticalEvents, err := keeper.QueryAuditEvents(ctx, &types.QueryAuditEventsRequest{
		Severity: &types.AuditSeverity_CRITICAL,
	})
	suite.Require().NoError(err)
	suite.Require().Len(criticalEvents.Events, 1)
	suite.Require().Equal("event_003", criticalEvents.Events[0].EventID)

	// Test query by time range
	startTime := time.Now().Add(-90 * time.Minute)
	endTime := time.Now().Add(-30 * time.Minute)
	timeRangeEvents, err := keeper.QueryAuditEvents(ctx, &types.QueryAuditEventsRequest{
		StartTime: &startTime,
		EndTime:   &endTime,
	})
	suite.Require().NoError(err)
	suite.Require().Len(timeRangeEvents.Events, 1)
	suite.Require().Equal("event_002", timeRangeEvents.Events[0].EventID)
}

func (suite *AuditKeeperTestSuite) TestGenerateComplianceReport() {
	ctx := suite.ctx
	keeper := suite.keeper

	// Create some audit events first
	events := []*types.AuditEvent{
		{
			EventID:      "comp_event_001",
			Timestamp:    time.Now().Add(-24 * time.Hour),
			EventType:    types.AuditEventType_IDENTITY_CREATED,
			Actor:        "desh1actor1",
			Subject:      "desh1subject1",
			Resource:     "identity_1",
			Action:       "create",
			Outcome:      types.AuditOutcome_SUCCESS,
			Severity:     types.AuditSeverity_LOW,
			Description:  "Compliant identity creation",
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
			EventID:      "comp_event_002",
			Timestamp:    time.Now().Add(-12 * time.Hour),
			EventType:    types.AuditEventType_DATA_ACCESS_DENIED,
			Actor:        "desh1actor2",
			Subject:      "desh1subject2",
			Resource:     "sensitive_data",
			Action:       "access_attempt",
			Outcome:      types.AuditOutcome_DENIED,
			Severity:     types.AuditSeverity_MEDIUM,
			Description:  "Unauthorized access attempt",
			ModuleSource: "identity",
			ChainHeight:  1001,
			ComplianceFlags: []types.ComplianceFlag{
				{
					Regulation:  "GDPR",
					Requirement: "Article 32 - Security",
					Status:      types.ComplianceStatus_COMPLIANT,
				},
			},
		},
		{
			EventID:      "comp_event_003",
			Timestamp:    time.Now().Add(-6 * time.Hour),
			EventType:    types.AuditEventType_COMPLIANCE_VIOLATION,
			Actor:        "desh1actor3",
			Subject:      "desh1subject3",
			Resource:     "personal_data",
			Action:       "unauthorized_processing",
			Outcome:      types.AuditOutcome_FAILURE,
			Severity:     types.AuditSeverity_CRITICAL,
			Description:  "Unauthorized data processing",
			ModuleSource: "identity",
			ChainHeight:  1002,
			ComplianceFlags: []types.ComplianceFlag{
				{
					Regulation:  "GDPR",
					Requirement: "Article 6 - Lawful basis",
					Status:      types.ComplianceStatus_NON_COMPLIANT,
				},
			},
		},
	}

	// Store events
	for _, event := range events {
		err := keeper.LogAuditEventToStore(ctx, event)
		suite.Require().NoError(err)
	}

	// Generate compliance report
	timeRange := types.AuditTimeRange{
		StartTime: time.Now().Add(-48 * time.Hour),
		EndTime:   time.Now(),
	}

	scope := types.ComplianceScope{
		IncludeModules: []string{"identity"},
	}

	regulations := []string{"GDPR"}

	report, err := keeper.GenerateComplianceReportInternal(
		ctx,
		types.ComplianceReportType_GDPR_COMPLIANCE,
		timeRange,
		scope,
		regulations,
	)
	suite.Require().NoError(err)
	suite.Require().NotNil(report)

	// Verify report contents
	suite.Require().Equal(types.ComplianceReportType_GDPR_COMPLIANCE, report.ReportType)
	suite.Require().Equal(timeRange, report.TimeRange)
	suite.Require().Equal(int64(3), report.TotalEvents)
	suite.Require().Equal(int64(3), report.DataSubjects) // 3 unique subjects
	suite.Require().True(report.ComplianceScore > 0)
	suite.Require().True(report.ComplianceScore <= 100)

	// Verify summary
	suite.Require().Equal(int64(2), report.Summary.TotalCompliantEvents)
	suite.Require().Equal(int64(1), report.Summary.TotalNonCompliantEvents)
	suite.Require().True(report.Summary.CompliancePercentage > 60)
	suite.Require().Equal(int64(1), report.Summary.CriticalFindings)

	// Verify findings
	suite.Require().True(len(report.Findings) > 0)
	violationFinding := false
	for _, finding := range report.Findings {
		if finding.Regulation == "GDPR" && finding.Status == types.FindingStatus_OPEN {
			violationFinding = true
			break
		}
	}
	suite.Require().True(violationFinding)

	// Store the report
	err = keeper.StoreComplianceReport(ctx, report)
	suite.Require().NoError(err)

	// Verify report was stored
	retrievedReport, found := keeper.GetComplianceReport(ctx, report.ReportID)
	suite.Require().True(found)
	suite.Require().Equal(report.ReportID, retrievedReport.ReportID)
	suite.Require().Equal(report.ComplianceScore, retrievedReport.ComplianceScore)
}

func (suite *AuditKeeperTestSuite) TestDataSubjectRequest() {
	ctx := suite.ctx
	keeper := suite.keeper

	// Create a data subject request
	request := &types.DataSubjectRequest{
		RequestID:     "dsr_test_001",
		RequestType:   types.DataSubjectRequestType_ACCESS,
		RequestStatus: types.DataSubjectRequestStatus_RECEIVED,
		DataSubject:   "desh1subject123",
		RequestedBy:   "desh1subject123",
		RequestDate:   time.Now(),
		DueDate:       time.Now().Add(30 * 24 * time.Hour), // 30 days
		Description:   "Request access to all personal data",
		RequestDetails: map[string]interface{}{
			"data_categories": []string{"personal", "financial"},
			"format":          "JSON",
		},
		Regulation: "GDPR",
		Priority:   types.DataSubjectRequestPriority_STANDARD,
	}

	// Process the request
	response, err := keeper.ProcessDataSubjectRequestInternal(ctx, request)
	suite.Require().NoError(err)
	suite.Require().NotNil(response)

	// Store the request
	err = keeper.StoreDataSubjectRequest(ctx, request)
	suite.Require().NoError(err)

	// Verify request was stored
	retrievedRequest, found := keeper.GetDataSubjectRequest(ctx, request.RequestID)
	suite.Require().True(found)
	suite.Require().Equal(request.RequestID, retrievedRequest.RequestID)
	suite.Require().Equal(request.RequestType, retrievedRequest.RequestType)
	suite.Require().Equal(request.DataSubject, retrievedRequest.DataSubject)
	suite.Require().Equal(request.Regulation, retrievedRequest.Regulation)

	// Test query data subject requests
	requests, err := keeper.QueryDataSubjectRequests(ctx, &types.QueryDataSubjectRequestsRequest{
		DataSubject: "desh1subject123",
	})
	suite.Require().NoError(err)
	suite.Require().Len(requests.Requests, 1)
	suite.Require().Equal("dsr_test_001", requests.Requests[0].RequestID)

	// Test query by request type
	accessRequests, err := keeper.QueryDataSubjectRequests(ctx, &types.QueryDataSubjectRequestsRequest{
		RequestType: &types.DataSubjectRequestType_ACCESS,
	})
	suite.Require().NoError(err)
	suite.Require().Len(accessRequests.Requests, 1)

	// Test query by regulation
	gdprRequests, err := keeper.QueryDataSubjectRequests(ctx, &types.QueryDataSubjectRequestsRequest{
		Regulation: "GDPR",
	})
	suite.Require().NoError(err)
	suite.Require().Len(gdprRequests.Requests, 1)
}

func (suite *AuditKeeperTestSuite) TestPrivacyImpactAssessment() {
	ctx := suite.ctx
	keeper := suite.keeper

	// Create a Privacy Impact Assessment
	assessment := &types.PrivacyImpactAssessment{
		AssessmentID:       "pia_test_001",
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
		TechnicalMeasures:  []string{"encryption", "access_controls"},
		OrganizationalMeasures: []string{"staff_training", "data_protection_policies"},
		ReviewDate: time.Now().AddDate(1, 0, 0),
		Status:     types.PIAStatus_DRAFT,
		CreatedBy:  "desh1admin123",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Perform risk assessment
	riskAssessment, err := keeper.PerformPrivacyRiskAssessment(ctx, assessment)
	suite.Require().NoError(err)
	suite.Require().NotNil(riskAssessment)
	assessment.RiskAssessment = *riskAssessment

	// Store the assessment
	err = keeper.StorePrivacyImpactAssessment(ctx, assessment)
	suite.Require().NoError(err)

	// Verify assessment was stored
	retrievedAssessment, found := keeper.GetPrivacyImpactAssessment(ctx, assessment.AssessmentID)
	suite.Require().True(found)
	suite.Require().Equal(assessment.AssessmentID, retrievedAssessment.AssessmentID)
	suite.Require().Equal(assessment.Title, retrievedAssessment.Title)
	suite.Require().Equal(assessment.ProcessingActivity, retrievedAssessment.ProcessingActivity)
	suite.Require().Equal(assessment.DataController, retrievedAssessment.DataController)

	// Test query PIAs
	assessments, err := keeper.QueryPrivacyImpactAssessments(ctx, &types.QueryPrivacyImpactAssessmentsRequest{
		DataController: "DeshChain Foundation",
	})
	suite.Require().NoError(err)
	suite.Require().Len(assessments.Assessments, 1)
	suite.Require().Equal("pia_test_001", assessments.Assessments[0].AssessmentID)

	// Test query by status
	draftAssessments, err := keeper.QueryPrivacyImpactAssessments(ctx, &types.QueryPrivacyImpactAssessmentsRequest{
		Status: &types.PIAStatus_DRAFT,
	})
	suite.Require().NoError(err)
	suite.Require().Len(draftAssessments.Assessments, 1)
}

func (suite *AuditKeeperTestSuite) TestAuditSettings() {
	ctx := suite.ctx
	keeper := suite.keeper

	// Test default settings
	defaultSettings, found := keeper.GetAuditSettings(ctx)
	if !found {
		// Set default settings if not found
		defaultSettings = types.AuditSettings{
			RetentionPeriodDays:         2555, // 7 years
			AutoPurgeEnabled:           true,
			ComplianceMonitoringEnabled: true,
			RealTimeAlertsEnabled:      true,
			AnonymizeExpiredData:       false,
		}
		err := keeper.UpdateAuditSettings(ctx, defaultSettings)
		suite.Require().NoError(err)
	}

	// Update audit settings
	newSettings := types.AuditSettings{
		RetentionPeriodDays:         1825, // 5 years
		AutoPurgeEnabled:           false,
		ComplianceMonitoringEnabled: true,
		RealTimeAlertsEnabled:      false,
		AnonymizeExpiredData:       true,
	}

	err := keeper.UpdateAuditSettings(ctx, newSettings)
	suite.Require().NoError(err)

	// Verify settings were updated
	retrievedSettings, found := keeper.GetAuditSettings(ctx)
	suite.Require().True(found)
	suite.Require().Equal(newSettings.RetentionPeriodDays, retrievedSettings.RetentionPeriodDays)
	suite.Require().Equal(newSettings.AutoPurgeEnabled, retrievedSettings.AutoPurgeEnabled)
	suite.Require().Equal(newSettings.ComplianceMonitoringEnabled, retrievedSettings.ComplianceMonitoringEnabled)
	suite.Require().Equal(newSettings.RealTimeAlertsEnabled, retrievedSettings.RealTimeAlertsEnabled)
	suite.Require().Equal(newSettings.AnonymizeExpiredData, retrievedSettings.AnonymizeExpiredData)
}

func (suite *AuditKeeperTestSuite) TestComplianceViolationHandling() {
	ctx := suite.ctx
	keeper := suite.keeper

	// Create a compliance violation event
	violationEvent := &types.AuditEvent{
		EventID:      "violation_001",
		Timestamp:    time.Now(),
		EventType:    types.AuditEventType_COMPLIANCE_VIOLATION,
		Actor:        "desh1malicious",
		Subject:      "desh1victim",
		Resource:     "sensitive_data",
		Action:       "unauthorized_access",
		Outcome:      types.AuditOutcome_FAILURE,
		Severity:     types.AuditSeverity_CRITICAL,
		Description:  "Attempted unauthorized access to sensitive data",
		ModuleSource: "identity",
		ChainHeight:  1000,
		ComplianceFlags: []types.ComplianceFlag{
			{
				Regulation:  "GDPR",
				Requirement: "Article 32 - Security of processing",
				Status:      types.ComplianceStatus_NON_COMPLIANT,
				Notes:       "Security breach detected",
			},
		},
	}

	// Handle compliance violation
	err := keeper.HandleComplianceViolation(ctx, violationEvent)
	suite.Require().NoError(err)

	// Verify that additional audit events were created for the violation handling
	// This would typically include notification events, escalation events, etc.
	allEvents, err := keeper.QueryAuditEvents(ctx, &types.QueryAuditEventsRequest{
		EventType: &types.AuditEventType_COMPLIANCE_VIOLATION,
	})
	suite.Require().NoError(err)
	suite.Require().True(len(allEvents.Events) >= 1)
}

func (suite *AuditKeeperTestSuite) TestAuditEventRetention() {
	ctx := suite.ctx
	keeper := suite.keeper

	// Set short retention period for testing
	settings := types.AuditSettings{
		RetentionPeriodDays:         1, // 1 day
		AutoPurgeEnabled:           true,
		ComplianceMonitoringEnabled: true,
		RealTimeAlertsEnabled:      true,
		AnonymizeExpiredData:       false,
	}
	err := keeper.UpdateAuditSettings(ctx, settings)
	suite.Require().NoError(err)

	// Create an old audit event
	oldEvent := &types.AuditEvent{
		EventID:      "old_event_001",
		Timestamp:    time.Now().Add(-48 * time.Hour), // 2 days ago
		EventType:    types.AuditEventType_IDENTITY_CREATED,
		Actor:        "desh1actor",
		Subject:      "desh1subject",
		Resource:     "identity_old",
		Action:       "create",
		Outcome:      types.AuditOutcome_SUCCESS,
		Severity:     types.AuditSeverity_LOW,
		Description:  "Old event for retention testing",
		ModuleSource: "identity",
		ChainHeight:  900,
	}

	err = keeper.LogAuditEventToStore(ctx, oldEvent)
	suite.Require().NoError(err)

	// Create a recent audit event
	recentEvent := &types.AuditEvent{
		EventID:      "recent_event_001",
		Timestamp:    time.Now().Add(-12 * time.Hour), // 12 hours ago
		EventType:    types.AuditEventType_IDENTITY_UPDATED,
		Actor:        "desh1actor",
		Subject:      "desh1subject",
		Resource:     "identity_recent",
		Action:       "update",
		Outcome:      types.AuditOutcome_SUCCESS,
		Severity:     types.AuditSeverity_LOW,
		Description:  "Recent event for retention testing",
		ModuleSource: "identity",
		ChainHeight:  1000,
	}

	err = keeper.LogAuditEventToStore(ctx, recentEvent)
	suite.Require().NoError(err)

	// Run retention cleanup
	deletedCount, err := keeper.CleanupExpiredAuditEvents(ctx)
	suite.Require().NoError(err)
	suite.Require().True(deletedCount > 0)

	// Verify old event was deleted
	_, found := keeper.GetAuditEvent(ctx, oldEvent.EventID)
	suite.Require().False(found)

	// Verify recent event still exists
	_, found = keeper.GetAuditEvent(ctx, recentEvent.EventID)
	suite.Require().True(found)
}

func (suite *AuditKeeperTestSuite) TestComplianceStatistics() {
	ctx := suite.ctx
	keeper := suite.keeper

	// Create various audit events for statistics
	events := []*types.AuditEvent{
		{
			EventID:      "stat_event_001",
			Timestamp:    time.Now().Add(-23 * time.Hour),
			EventType:    types.AuditEventType_IDENTITY_CREATED,
			Outcome:      types.AuditOutcome_SUCCESS,
			Severity:     types.AuditSeverity_LOW,
			ModuleSource: "identity",
			ComplianceFlags: []types.ComplianceFlag{
				{Regulation: "GDPR", Status: types.ComplianceStatus_COMPLIANT},
			},
		},
		{
			EventID:      "stat_event_002",
			Timestamp:    time.Now().Add(-22 * time.Hour),
			EventType:    types.AuditEventType_CREDENTIAL_ISSUED,
			Outcome:      types.AuditOutcome_SUCCESS,
			Severity:     types.AuditSeverity_MEDIUM,
			ModuleSource: "identity",
			ComplianceFlags: []types.ComplianceFlag{
				{Regulation: "GDPR", Status: types.ComplianceStatus_COMPLIANT},
			},
		},
		{
			EventID:      "stat_event_003",
			Timestamp:    time.Now().Add(-21 * time.Hour),
			EventType:    types.AuditEventType_COMPLIANCE_VIOLATION,
			Outcome:      types.AuditOutcome_FAILURE,
			Severity:     types.AuditSeverity_CRITICAL,
			ModuleSource: "identity",
			ComplianceFlags: []types.ComplianceFlag{
				{Regulation: "GDPR", Status: types.ComplianceStatus_NON_COMPLIANT},
			},
		},
	}

	// Store events
	for _, event := range events {
		err := keeper.LogAuditEventToStore(ctx, event)
		suite.Require().NoError(err)
	}

	// Query compliance statistics
	stats, err := keeper.QueryComplianceStatistics(ctx, &types.QueryComplianceStatisticsRequest{
		TimeWindowDays: 1, // Last 24 hours
	})
	suite.Require().NoError(err)
	suite.Require().NotNil(stats)

	// Verify statistics
	suite.Require().True(stats.TotalEvents >= 3)
	suite.Require().True(stats.CompliantEvents >= 2)
	suite.Require().True(stats.NonCompliantEvents >= 1)
	suite.Require().True(stats.CriticalViolations >= 1)
	suite.Require().True(stats.CompliancePercentage > 0)
	suite.Require().True(stats.CompliancePercentage <= 100)
}

// Benchmark tests

func BenchmarkLogAuditEvent(b *testing.B) {
	suite := setupTestSuite()
	defer suite.TearDownSuite()

	ctx := suite.ctx
	keeper := suite.keeper

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		event := &types.AuditEvent{
			EventID:      fmt.Sprintf("bench_event_%d", i),
			Timestamp:    time.Now(),
			EventType:    types.AuditEventType_IDENTITY_CREATED,
			Actor:        "desh1benchactor",
			Subject:      "desh1benchsubject",
			Resource:     fmt.Sprintf("resource_%d", i),
			Action:       "benchmark_action",
			Outcome:      types.AuditOutcome_SUCCESS,
			Severity:     types.AuditSeverity_LOW,
			Description:  "Benchmark audit event",
			ModuleSource: "identity",
			ChainHeight:  int64(1000 + i),
		}

		err := keeper.LogAuditEventToStore(ctx, event)
		require.NoError(b, err)
	}
}

func BenchmarkQueryAuditEvents(b *testing.B) {
	suite := setupTestSuite()
	defer suite.TearDownSuite()

	ctx := suite.ctx
	keeper := suite.keeper

	// Pre-populate with 1000 events
	for i := 0; i < 1000; i++ {
		event := &types.AuditEvent{
			EventID:      fmt.Sprintf("query_bench_event_%d", i),
			Timestamp:    time.Now().Add(-time.Duration(i) * time.Minute),
			EventType:    types.AuditEventType_IDENTITY_CREATED,
			Actor:        fmt.Sprintf("desh1actor%d", i%10),
			Subject:      fmt.Sprintf("desh1subject%d", i%5),
			Resource:     fmt.Sprintf("resource_%d", i),
			Action:       "create",
			Outcome:      types.AuditOutcome_SUCCESS,
			Severity:     types.AuditSeverity_LOW,
			Description:  "Benchmark query event",
			ModuleSource: "identity",
			ChainHeight:  int64(1000 + i),
		}

		err := keeper.LogAuditEventToStore(ctx, event)
		require.NoError(b, err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := keeper.QueryAuditEvents(ctx, &types.QueryAuditEventsRequest{
			Actor: "desh1actor0",
		})
		require.NoError(b, err)
	}
}

// Helper function to setup test suite for benchmarks
func setupTestSuite() *AuditKeeperTestSuite {
	suite := &AuditKeeperTestSuite{}
	suite.SetupTest()
	return suite
}