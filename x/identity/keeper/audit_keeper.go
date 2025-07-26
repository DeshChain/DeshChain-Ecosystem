package keeper

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/DeshChain/DeshChain-Ecosystem/x/identity/types"
)

// Identity Audit and Compliance Keeper Implementation

// LogAuditEvent creates and stores an audit event
func (k Keeper) LogAuditEvent(
	ctx sdk.Context,
	eventType types.AuditEventType,
	actor string,
	subject string,
	resource string,
	action string,
	outcome types.AuditOutcome,
	description string,
	metadata map[string]interface{},
) (*types.AuditEvent, error) {
	// Generate event ID
	eventID, err := k.generateAuditEventID()
	if err != nil {
		return nil, err
	}

	// Determine severity based on event type and outcome
	severity := k.calculateEventSeverity(eventType, outcome)

	// Extract technical details from context
	technicalDetails := map[string]interface{}{
		"block_height": ctx.BlockHeight(),
		"block_time":   ctx.BlockTime().Format(time.RFC3339),
		"gas_used":     ctx.GasMeter().GasConsumed(),
	}

	// Add transaction hash if available
	txHash := ""
	if txBytes := ctx.TxBytes(); txBytes != nil {
		txHash = fmt.Sprintf("%x", txBytes)
	}

	// Create audit event
	auditEvent := &types.AuditEvent{
		EventID:          eventID,
		Timestamp:        ctx.BlockTime(),
		EventType:        eventType,
		Actor:            actor,
		Subject:          subject,
		Resource:         resource,
		Action:           action,
		Outcome:          outcome,
		Severity:         severity,
		Description:      description,
		TechnicalDetails: technicalDetails,
		ComplianceFlags:  k.generateComplianceFlags(ctx, eventType, actor, subject),
		Metadata:         metadata,
		ModuleSource:     "identity",
		ChainHeight:      ctx.BlockHeight(),
		TxHash:           txHash,
	}

	// Store audit event
	k.SetAuditEvent(ctx, *auditEvent)

	// Index for efficient queries
	k.SetAuditEventByActor(ctx, actor, eventID)
	k.SetAuditEventBySubject(ctx, subject, eventID)
	k.SetAuditEventByType(ctx, eventType, eventID)
	k.SetAuditEventByDate(ctx, ctx.BlockTime(), eventID)

	// Check for compliance violations
	if violation := k.checkComplianceViolation(ctx, auditEvent); violation != nil {
		k.handleComplianceViolation(ctx, violation)
	}

	// Emit blockchain event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"audit_event_logged",
			sdk.NewAttribute("event_id", eventID),
			sdk.NewAttribute("event_type", eventType.String()),
			sdk.NewAttribute("actor", actor),
			sdk.NewAttribute("subject", subject),
			sdk.NewAttribute("severity", severity.String()),
		),
	)

	return auditEvent, nil
}

// GetAuditEvent retrieves an audit event by ID
func (k Keeper) GetAuditEvent(ctx sdk.Context, eventID string) (types.AuditEvent, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.AuditEventKey(eventID))
	if bz == nil {
		return types.AuditEvent{}, false
	}

	var event types.AuditEvent
	k.cdc.MustUnmarshal(bz, &event)
	return event, true
}

// GetAuditEventsByActor retrieves audit events for a specific actor
func (k Keeper) GetAuditEventsByActor(ctx sdk.Context, actor string, startTime, endTime time.Time) []types.AuditEvent {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.AuditEventByActorPrefix(actor))
	defer iterator.Close()

	var events []types.AuditEvent
	for ; iterator.Valid(); iterator.Next() {
		eventID := string(iterator.Value())
		if event, found := k.GetAuditEvent(ctx, eventID); found {
			if event.Timestamp.After(startTime) && event.Timestamp.Before(endTime) {
				events = append(events, event)
			}
		}
	}

	return events
}

// GetAuditEventsBySubject retrieves audit events for a specific subject
func (k Keeper) GetAuditEventsBySubject(ctx sdk.Context, subject string, startTime, endTime time.Time) []types.AuditEvent {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.AuditEventBySubjectPrefix(subject))
	defer iterator.Close()

	var events []types.AuditEvent
	for ; iterator.Valid(); iterator.Next() {
		eventID := string(iterator.Value())
		if event, found := k.GetAuditEvent(ctx, eventID); found {
			if event.Timestamp.After(startTime) && event.Timestamp.Before(endTime) {
				events = append(events, event)
			}
		}
	}

	return events
}

// GetAuditEventsByType retrieves audit events by type
func (k Keeper) GetAuditEventsByType(ctx sdk.Context, eventType types.AuditEventType, startTime, endTime time.Time) []types.AuditEvent {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.AuditEventByTypePrefix(eventType))
	defer iterator.Close()

	var events []types.AuditEvent
	for ; iterator.Valid(); iterator.Next() {
		eventID := string(iterator.Value())
		if event, found := k.GetAuditEvent(ctx, eventID); found {
			if event.Timestamp.After(startTime) && event.Timestamp.Before(endTime) {
				events = append(events, event)
			}
		}
	}

	return events
}

// GenerateComplianceReport creates a comprehensive compliance report
func (k Keeper) GenerateComplianceReport(
	ctx sdk.Context,
	reportType types.ComplianceReportType,
	timeRange types.AuditTimeRange,
	scope types.ComplianceScope,
	regulations []string,
	generatedBy string,
) (*types.ComplianceReport, error) {
	// Generate report ID
	reportID, err := k.generateComplianceReportID()
	if err != nil {
		return nil, err
	}

	// Collect audit events within the time range and scope
	events := k.collectAuditEventsForReport(ctx, timeRange, scope)

	// Analyze events for compliance
	summary := k.analyzeComplianceEvents(events, regulations)
	findings := k.identifyComplianceFindings(events, regulations)
	recommendations := k.generateComplianceRecommendations(findings)

	// Calculate compliance score
	complianceScore := k.calculateComplianceScore(summary, findings)
	riskLevel := k.assessComplianceRisk(findings, complianceScore)

	// Count unique data subjects
	dataSubjects := k.countUniqueDataSubjects(events)

	// Create compliance report
	report := &types.ComplianceReport{
		ReportID:        reportID,
		GeneratedAt:     ctx.BlockTime(),
		GeneratedBy:     generatedBy,
		ReportType:      reportType,
		TimeRange:       timeRange,
		Scope:           scope,
		Regulations:     regulations,
		Summary:         summary,
		Findings:        findings,
		Recommendations: recommendations,
		DataSubjects:    dataSubjects,
		TotalEvents:     int64(len(events)),
		ComplianceScore: complianceScore,
		RiskLevel:       riskLevel,
		NextReviewDate:  ctx.BlockTime().Add(90 * 24 * time.Hour), // 90 days
		Metadata:        make(map[string]interface{}),
	}

	// Store compliance report
	k.SetComplianceReport(ctx, *report)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"compliance_report_generated",
			sdk.NewAttribute("report_id", reportID),
			sdk.NewAttribute("report_type", fmt.Sprintf("%d", reportType)),
			sdk.NewAttribute("compliance_score", fmt.Sprintf("%.2f", complianceScore)),
			sdk.NewAttribute("risk_level", fmt.Sprintf("%d", riskLevel)),
		),
	)

	return report, nil
}

// CreateDataSubjectRequest handles data subject requests (GDPR Article 12-22)
func (k Keeper) CreateDataSubjectRequest(
	ctx sdk.Context,
	requestType types.DataSubjectRequestType,
	dataSubject string,
	requestedBy string,
	description string,
	requestDetails map[string]interface{},
	regulation string,
) (*types.DataSubjectRequest, error) {
	// Generate request ID
	requestID, err := k.generateDataSubjectRequestID()
	if err != nil {
		return nil, err
	}

	// Calculate due date based on regulation
	dueDate := k.calculateDataSubjectRequestDueDate(requestType, regulation, ctx.BlockTime())

	// Determine priority
	priority := k.determineDataSubjectRequestPriority(requestType)

	// Create request
	request := &types.DataSubjectRequest{
		RequestID:       requestID,
		RequestType:     requestType,
		RequestStatus:   types.DataSubjectRequestStatus_RECEIVED,
		DataSubject:     dataSubject,
		RequestedBy:     requestedBy,
		RequestDate:     ctx.BlockTime(),
		DueDate:         dueDate,
		Description:     description,
		RequestDetails:  requestDetails,
		LegalBasis:      k.determineLegalBasis(requestType),
		Regulation:      regulation,
		Priority:        priority,
		Metadata:        make(map[string]interface{}),
	}

	// Store request
	k.SetDataSubjectRequest(ctx, *request)

	// Index for efficient queries
	k.SetDataSubjectRequestBySubject(ctx, dataSubject, requestID)
	k.SetDataSubjectRequestByType(ctx, requestType, requestID)
	k.SetDataSubjectRequestByStatus(ctx, request.RequestStatus, requestID)

	// Log audit event
	k.LogAuditEvent(
		ctx,
		types.AuditEventType_DATA_ACCESS_REQUESTED,
		requestedBy,
		dataSubject,
		requestID,
		fmt.Sprintf("data_subject_request_%s", requestType),
		types.AuditOutcome_SUCCESS,
		fmt.Sprintf("Data subject request created: %s", description),
		map[string]interface{}{
			"request_type": requestType,
			"regulation":   regulation,
			"due_date":     dueDate.Format(time.RFC3339),
		},
	)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"data_subject_request_created",
			sdk.NewAttribute("request_id", requestID),
			sdk.NewAttribute("request_type", fmt.Sprintf("%d", requestType)),
			sdk.NewAttribute("data_subject", dataSubject),
			sdk.NewAttribute("due_date", dueDate.Format(time.RFC3339)),
		),
	)

	return request, nil
}

// ProcessDataSubjectRequest processes a data subject request
func (k Keeper) ProcessDataSubjectRequest(
	ctx sdk.Context,
	requestID string,
	processedBy string,
	processingNotes string,
) error {
	// Get request
	request, found := k.GetDataSubjectRequest(ctx, requestID)
	if !found {
		return types.ErrDataSubjectRequestNotFound
	}

	// Check if request is expired
	if request.IsExpired() {
		return types.ErrDataSubjectRequestExpired
	}

	// Update request status
	request.RequestStatus = types.DataSubjectRequestStatus_IN_PROGRESS
	processedDate := ctx.BlockTime()
	request.ProcessedDate = &processedDate
	request.ProcessingNotes = processingNotes
	request.AssignedTo = processedBy

	// Process based on request type
	responseData, err := k.executeDataSubjectRequest(ctx, &request)
	if err != nil {
		request.RequestStatus = types.DataSubjectRequestStatus_REJECTED
		k.SetDataSubjectRequest(ctx, request)
		return err
	}

	// Create response
	if responseData != nil {
		request.ResponseData = responseData
	}

	// Update status
	request.RequestStatus = types.DataSubjectRequestStatus_COMPLETED
	completedDate := ctx.BlockTime()
	request.CompletedDate = &completedDate

	// Store updated request
	k.SetDataSubjectRequest(ctx, request)

	// Log audit event
	k.LogAuditEvent(
		ctx,
		types.AuditEventType_DATA_ACCESS_REQUESTED,
		processedBy,
		request.DataSubject,
		requestID,
		fmt.Sprintf("data_subject_request_processed_%s", request.RequestType),
		types.AuditOutcome_SUCCESS,
		fmt.Sprintf("Data subject request processed: %s", request.Description),
		map[string]interface{}{
			"request_type":     request.RequestType,
			"processing_time":  request.GetProcessingTime().String(),
			"processing_notes": processingNotes,
		},
	)

	return nil
}

// CreatePrivacyImpactAssessment creates a DPIA
func (k Keeper) CreatePrivacyImpactAssessment(
	ctx sdk.Context,
	title string,
	description string,
	processingActivity string,
	dataController string,
	createdBy string,
) (*types.PrivacyImpactAssessment, error) {
	// Generate assessment ID
	assessmentID, err := k.generatePrivacyImpactAssessmentID()
	if err != nil {
		return nil, err
	}

	// Create PIA
	pia := &types.PrivacyImpactAssessment{
		AssessmentID:       assessmentID,
		Title:              title,
		Description:        description,
		ProcessingActivity: processingActivity,
		DataController:     dataController,
		Status:             types.PIAStatus_DRAFT,
		CreatedBy:          createdBy,
		CreatedAt:          ctx.BlockTime(),
		UpdatedAt:          ctx.BlockTime(),
		ReviewDate:         ctx.BlockTime().Add(365 * 24 * time.Hour), // Annual review
		Metadata:           make(map[string]interface{}),
	}

	// Store PIA
	k.SetPrivacyImpactAssessment(ctx, *pia)

	// Log audit event
	k.LogAuditEvent(
		ctx,
		types.AuditEventType_ADMIN_ACTION,
		createdBy,
		dataController,
		assessmentID,
		"privacy_impact_assessment_created",
		types.AuditOutcome_SUCCESS,
		fmt.Sprintf("Privacy Impact Assessment created: %s", title),
		map[string]interface{}{
			"processing_activity": processingActivity,
		},
	)

	return pia, nil
}

// Storage functions

// SetAuditEvent stores an audit event
func (k Keeper) SetAuditEvent(ctx sdk.Context, event types.AuditEvent) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&event)
	store.Set(types.AuditEventKey(event.EventID), bz)
}

// SetComplianceReport stores a compliance report
func (k Keeper) SetComplianceReport(ctx sdk.Context, report types.ComplianceReport) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&report)
	store.Set(types.ComplianceReportKey(report.ReportID), bz)
}

// GetComplianceReport retrieves a compliance report
func (k Keeper) GetComplianceReport(ctx sdk.Context, reportID string) (types.ComplianceReport, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ComplianceReportKey(reportID))
	if bz == nil {
		return types.ComplianceReport{}, false
	}

	var report types.ComplianceReport
	k.cdc.MustUnmarshal(bz, &report)
	return report, true
}

// SetDataSubjectRequest stores a data subject request
func (k Keeper) SetDataSubjectRequest(ctx sdk.Context, request types.DataSubjectRequest) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&request)
	store.Set(types.DataSubjectRequestKey(request.RequestID), bz)
}

// GetDataSubjectRequest retrieves a data subject request
func (k Keeper) GetDataSubjectRequest(ctx sdk.Context, requestID string) (types.DataSubjectRequest, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.DataSubjectRequestKey(requestID))
	if bz == nil {
		return types.DataSubjectRequest{}, false
	}

	var request types.DataSubjectRequest
	k.cdc.MustUnmarshal(bz, &request)
	return request, true
}

// SetPrivacyImpactAssessment stores a privacy impact assessment
func (k Keeper) SetPrivacyImpactAssessment(ctx sdk.Context, pia types.PrivacyImpactAssessment) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&pia)
	store.Set(types.PrivacyImpactAssessmentKey(pia.AssessmentID), bz)
}

// GetPrivacyImpactAssessment retrieves a privacy impact assessment
func (k Keeper) GetPrivacyImpactAssessment(ctx sdk.Context, assessmentID string) (types.PrivacyImpactAssessment, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.PrivacyImpactAssessmentKey(assessmentID))
	if bz == nil {
		return types.PrivacyImpactAssessment{}, false
	}

	var pia types.PrivacyImpactAssessment
	k.cdc.MustUnmarshal(bz, &pia)
	return pia, true
}

// Index functions

// SetAuditEventByActor creates an index for audit events by actor
func (k Keeper) SetAuditEventByActor(ctx sdk.Context, actor, eventID string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.AuditEventByActorKey(actor, eventID), []byte(eventID))
}

// SetAuditEventBySubject creates an index for audit events by subject
func (k Keeper) SetAuditEventBySubject(ctx sdk.Context, subject, eventID string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.AuditEventBySubjectKey(subject, eventID), []byte(eventID))
}

// SetAuditEventByType creates an index for audit events by type
func (k Keeper) SetAuditEventByType(ctx sdk.Context, eventType types.AuditEventType, eventID string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.AuditEventByTypeKey(eventType, eventID), []byte(eventID))
}

// SetAuditEventByDate creates an index for audit events by date
func (k Keeper) SetAuditEventByDate(ctx sdk.Context, date time.Time, eventID string) {
	store := ctx.KVStore(k.storeKey)
	dateKey := date.Format("2006-01-02")
	store.Set(types.AuditEventByDateKey(dateKey, eventID), []byte(eventID))
}

// SetDataSubjectRequestBySubject creates an index for requests by subject
func (k Keeper) SetDataSubjectRequestBySubject(ctx sdk.Context, subject, requestID string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.DataSubjectRequestBySubjectKey(subject, requestID), []byte(requestID))
}

// SetDataSubjectRequestByType creates an index for requests by type
func (k Keeper) SetDataSubjectRequestByType(ctx sdk.Context, requestType types.DataSubjectRequestType, requestID string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.DataSubjectRequestByTypeKey(requestType, requestID), []byte(requestID))
}

// SetDataSubjectRequestByStatus creates an index for requests by status
func (k Keeper) SetDataSubjectRequestByStatus(ctx sdk.Context, status types.DataSubjectRequestStatus, requestID string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.DataSubjectRequestByStatusKey(status, requestID), []byte(requestID))
}

// Helper functions

// generateAuditEventID generates a unique audit event ID
func (k Keeper) generateAuditEventID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "audit_" + hex.EncodeToString(bytes), nil
}

// generateComplianceReportID generates a unique compliance report ID
func (k Keeper) generateComplianceReportID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "report_" + hex.EncodeToString(bytes), nil
}

// generateDataSubjectRequestID generates a unique data subject request ID
func (k Keeper) generateDataSubjectRequestID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "dsr_" + hex.EncodeToString(bytes), nil
}

// generatePrivacyImpactAssessmentID generates a unique PIA ID
func (k Keeper) generatePrivacyImpactAssessmentID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "pia_" + hex.EncodeToString(bytes), nil
}

// calculateEventSeverity determines the severity of an audit event
func (k Keeper) calculateEventSeverity(eventType types.AuditEventType, outcome types.AuditOutcome) types.AuditSeverity {
	// Critical events
	if eventType == types.AuditEventType_COMPLIANCE_VIOLATION ||
		eventType == types.AuditEventType_SUSPICIOUS_ACTIVITY {
		return types.AuditSeverity_CRITICAL
	}

	// High severity for failed security-sensitive operations
	if outcome == types.AuditOutcome_FAILURE || outcome == types.AuditOutcome_DENIED {
		switch eventType {
		case types.AuditEventType_IDENTITY_ACCESSED,
			types.AuditEventType_CREDENTIAL_ACCESSED,
			types.AuditEventType_DATA_ACCESS_REQUESTED:
			return types.AuditSeverity_HIGH
		}
	}

	// Medium severity for data operations
	switch eventType {
	case types.AuditEventType_DATA_SHARED,
		types.AuditEventType_CONSENT_WITHDRAWN,
		types.AuditEventType_CREDENTIAL_REVOKED:
		return types.AuditSeverity_MEDIUM
	}

	// Default to low severity
	return types.AuditSeverity_LOW
}

// generateComplianceFlags creates compliance flags for an event
func (k Keeper) generateComplianceFlags(ctx sdk.Context, eventType types.AuditEventType, actor, subject string) []types.ComplianceFlag {
	var flags []types.ComplianceFlag

	// GDPR compliance flags
	switch eventType {
	case types.AuditEventType_CONSENT_GIVEN:
		flags = append(flags, types.ComplianceFlag{
			Regulation:  "GDPR",
			Requirement: "Article 7 - Consent",
			Status:      types.ComplianceStatus_COMPLIANT,
			Notes:       "Explicit consent recorded",
		})
	case types.AuditEventType_CONSENT_WITHDRAWN:
		flags = append(flags, types.ComplianceFlag{
			Regulation:  "GDPR",
			Requirement: "Article 7(3) - Withdrawal of consent",
			Status:      types.ComplianceStatus_COMPLIANT,
			Notes:       "Consent withdrawal processed",
		})
	case types.AuditEventType_DATA_SHARED:
		flags = append(flags, types.ComplianceFlag{
			Regulation:  "GDPR",
			Requirement: "Article 6 - Lawfulness of processing",
			Status:      types.ComplianceStatus_PENDING,
			Notes:       "Requires legal basis verification",
		})
	}

	// DPDP Act compliance flags
	if strings.Contains(subject, "desh") { // Indian users
		switch eventType {
		case types.AuditEventType_DATA_SHARED:
			flags = append(flags, types.ComplianceFlag{
				Regulation:  "DPDP",
				Requirement: "Section 8 - Purpose limitation",
				Status:      types.ComplianceStatus_PENDING,
				Notes:       "Purpose validation required",
			})
		}
	}

	return flags
}

// checkComplianceViolation checks if an event represents a compliance violation
func (k Keeper) checkComplianceViolation(ctx sdk.Context, event *types.AuditEvent) *types.ComplianceFinding {
	// Check for suspicious patterns
	if event.EventType == types.AuditEventType_SUSPICIOUS_ACTIVITY {
		return &types.ComplianceFinding{
			Regulation:  "GDPR",
			Requirement: "Article 32 - Security of processing",
			Severity:    types.AuditSeverity_HIGH,
			Status:      types.FindingStatus_OPEN,
			Description: "Suspicious activity detected",
			Evidence:    []string{event.EventID},
		}
	}

	// Check for failed access attempts
	if event.Outcome == types.AuditOutcome_DENIED &&
		event.EventType == types.AuditEventType_DATA_ACCESS_REQUESTED {
		// Check if there are multiple failed attempts
		recentEvents := k.GetAuditEventsByActor(ctx, event.Actor,
			ctx.BlockTime().Add(-1*time.Hour), ctx.BlockTime())
		
		failedAttempts := 0
		for _, recentEvent := range recentEvents {
			if recentEvent.Outcome == types.AuditOutcome_DENIED {
				failedAttempts++
			}
		}

		if failedAttempts > 5 { // Threshold for suspicious activity
			return &types.ComplianceFinding{
				Regulation:  "GDPR",
				Requirement: "Article 32 - Security of processing",
				Severity:    types.AuditSeverity_MEDIUM,
				Status:      types.FindingStatus_OPEN,
				Description: fmt.Sprintf("Multiple failed access attempts (%d) by actor %s", failedAttempts, event.Actor),
				Evidence:    []string{event.EventID},
			}
		}
	}

	return nil
}

// handleComplianceViolation processes a compliance violation
func (k Keeper) handleComplianceViolation(ctx sdk.Context, violation *types.ComplianceFinding) {
	// Log the violation as a critical audit event
	k.LogAuditEvent(
		ctx,
		types.AuditEventType_COMPLIANCE_VIOLATION,
		"system",
		"compliance_monitor",
		violation.FindingID,
		"compliance_violation_detected",
		types.AuditOutcome_SUCCESS,
		violation.Description,
		map[string]interface{}{
			"regulation":   violation.Regulation,
			"requirement":  violation.Requirement,
			"severity":     violation.Severity.String(),
			"evidence":     violation.Evidence,
		},
	)
}

// collectAuditEventsForReport collects audit events for a compliance report
func (k Keeper) collectAuditEventsForReport(ctx sdk.Context, timeRange types.AuditTimeRange, scope types.ComplianceScope) []types.AuditEvent {
	var events []types.AuditEvent

	// If specific event types are included, query by type
	if len(scope.IncludeEventTypes) > 0 {
		for _, eventType := range scope.IncludeEventTypes {
			typeEvents := k.GetAuditEventsByType(ctx, eventType, timeRange.StartTime, timeRange.EndTime)
			events = append(events, typeEvents...)
		}
	} else {
		// Query all events in time range and filter later
		store := ctx.KVStore(k.storeKey)
		iterator := sdk.KVStorePrefixIterator(store, types.AuditEventPrefix)
		defer iterator.Close()

		for ; iterator.Valid(); iterator.Next() {
			var event types.AuditEvent
			k.cdc.MustUnmarshal(iterator.Value(), &event)
			
			if event.Timestamp.After(timeRange.StartTime) && event.Timestamp.Before(timeRange.EndTime) {
				events = append(events, event)
			}
		}
	}

	// Apply scope filters
	return k.filterEventsByScope(events, scope)
}

// filterEventsByScope applies scope filters to events
func (k Keeper) filterEventsByScope(events []types.AuditEvent, scope types.ComplianceScope) []types.AuditEvent {
	var filtered []types.AuditEvent

	for _, event := range events {
		// Apply module filters
		if len(scope.IncludeModules) > 0 {
			if !contains(scope.IncludeModules, event.ModuleSource) {
				continue
			}
		}
		if len(scope.ExcludeModules) > 0 {
			if contains(scope.ExcludeModules, event.ModuleSource) {
				continue
			}
		}

		// Apply event type filters
		if len(scope.ExcludeEventTypes) > 0 {
			if containsEventType(scope.ExcludeEventTypes, event.EventType) {
				continue
			}
		}

		// Apply data subject filters
		if len(scope.DataSubjects) > 0 {
			if !contains(scope.DataSubjects, event.Subject) {
				continue
			}
		}

		filtered = append(filtered, event)
	}

	return filtered
}

// analyzeComplianceEvents analyzes events for compliance summary
func (k Keeper) analyzeComplianceEvents(events []types.AuditEvent, regulations []string) types.ComplianceReportSummary {
	summary := types.ComplianceReportSummary{}

	for _, event := range events {
		for _, flag := range event.ComplianceFlags {
			if contains(regulations, flag.Regulation) {
				switch flag.Status {
				case types.ComplianceStatus_COMPLIANT:
					summary.TotalCompliantEvents++
				case types.ComplianceStatus_NON_COMPLIANT:
					summary.TotalNonCompliantEvents++
				}
			}
		}

		// Count specific event types
		switch event.EventType {
		case types.AuditEventType_CONSENT_WITHDRAWN:
			summary.ConsentWithdrawals++
		case types.AuditEventType_DATA_ACCESS_REQUESTED:
			summary.DataAccessRequests++
		case types.AuditEventType_IDENTITY_DELETED:
			summary.DataDeletionRequests++
		case types.AuditEventType_DATA_SHARED:
			summary.DataProcessingActivities++
		}
	}

	total := summary.TotalCompliantEvents + summary.TotalNonCompliantEvents
	if total > 0 {
		summary.CompliancePercentage = float64(summary.TotalCompliantEvents) / float64(total) * 100
	}

	return summary
}

// identifyComplianceFindings identifies compliance issues
func (k Keeper) identifyComplianceFindings(events []types.AuditEvent, regulations []string) []types.ComplianceFinding {
	var findings []types.ComplianceFinding

	// Track patterns that might indicate compliance issues
	actorEventCounts := make(map[string]int)
	
	for _, event := range events {
		actorEventCounts[event.Actor]++

		// Check for non-compliant flags
		for _, flag := range event.ComplianceFlags {
			if flag.Status == types.ComplianceStatus_NON_COMPLIANT {
				finding := types.ComplianceFinding{
					FindingID:   fmt.Sprintf("finding_%s_%s", event.EventID, flag.Regulation),
					Regulation:  flag.Regulation,
					Requirement: flag.Requirement,
					Severity:    event.Severity,
					Status:      types.FindingStatus_OPEN,
					Description: fmt.Sprintf("Compliance violation: %s", flag.Notes),
					Evidence:    []string{event.EventID},
					CreatedAt:   event.Timestamp,
					UpdatedAt:   event.Timestamp,
				}
				findings = append(findings, finding)
			}
		}
	}

	// Check for suspicious patterns
	for actor, count := range actorEventCounts {
		if count > 100 { // Threshold for high activity
			finding := types.ComplianceFinding{
				FindingID:   fmt.Sprintf("finding_high_activity_%s", actor),
				Regulation:  "GDPR",
				Requirement: "Article 32 - Security of processing",
				Severity:    types.AuditSeverity_MEDIUM,
				Status:      types.FindingStatus_OPEN,
				Description: fmt.Sprintf("High activity detected for actor %s (%d events)", actor, count),
				Evidence:    []string{},
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			findings = append(findings, finding)
		}
	}

	return findings
}

// generateComplianceRecommendations generates recommendations based on findings
func (k Keeper) generateComplianceRecommendations(findings []types.ComplianceFinding) []types.ComplianceRecommendation {
	var recommendations []types.ComplianceRecommendation

	// Group findings by regulation
	findingsByRegulation := make(map[string][]types.ComplianceFinding)
	for _, finding := range findings {
		findingsByRegulation[finding.Regulation] = append(findingsByRegulation[finding.Regulation], finding)
	}

	// Generate recommendations for each regulation
	for regulation, regulationFindings := range findingsByRegulation {
		if len(regulationFindings) > 5 {
			recommendation := types.ComplianceRecommendation{
				RecommendationID: fmt.Sprintf("rec_%s_training", regulation),
				Title:           fmt.Sprintf("Enhanced %s Training", regulation),
				Description:     fmt.Sprintf("Multiple compliance issues detected for %s. Consider enhanced training for staff.", regulation),
				Priority:        types.RecommendationPriority_HIGH,
				Category:        types.RecommendationCategory_TRAINING,
				EstimatedEffort: "2-4 weeks",
				ExpectedBenefit: "Reduced compliance violations",
				Implementation:  "Develop and deliver comprehensive training program",
			}
			recommendations = append(recommendations, recommendation)
		}
	}

	return recommendations
}

// calculateComplianceScore calculates overall compliance score
func (k Keeper) calculateComplianceScore(summary types.ComplianceReportSummary, findings []types.ComplianceFinding) float64 {
	baseScore := summary.CompliancePercentage

	// Deduct points for findings based on severity
	deductions := 0.0
	for _, finding := range findings {
		switch finding.Severity {
		case types.AuditSeverity_CRITICAL:
			deductions += 10.0
		case types.AuditSeverity_HIGH:
			deductions += 5.0
		case types.AuditSeverity_MEDIUM:
			deductions += 2.0
		case types.AuditSeverity_LOW:
			deductions += 0.5
		}
	}

	score := baseScore - deductions
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score
}

// assessComplianceRisk assesses overall compliance risk level
func (k Keeper) assessComplianceRisk(findings []types.ComplianceFinding, complianceScore float64) types.ComplianceRiskLevel {
	if complianceScore < 50 {
		return types.ComplianceRiskLevel_CRITICAL
	}

	criticalFindings := 0
	for _, finding := range findings {
		if finding.Severity == types.AuditSeverity_CRITICAL {
			criticalFindings++
		}
	}

	if criticalFindings > 0 {
		return types.ComplianceRiskLevel_HIGH
	}

	if complianceScore < 75 {
		return types.ComplianceRiskLevel_MEDIUM
	}

	return types.ComplianceRiskLevel_LOW
}

// countUniqueDataSubjects counts unique data subjects in events
func (k Keeper) countUniqueDataSubjects(events []types.AuditEvent) int64 {
	subjects := make(map[string]bool)
	for _, event := range events {
		if event.Subject != "" {
			subjects[event.Subject] = true
		}
	}
	return int64(len(subjects))
}

// calculateDataSubjectRequestDueDate calculates due date based on regulation
func (k Keeper) calculateDataSubjectRequestDueDate(requestType types.DataSubjectRequestType, regulation string, requestDate time.Time) time.Time {
	switch regulation {
	case "GDPR":
		return requestDate.Add(30 * 24 * time.Hour) // 30 days
	case "DPDP":
		return requestDate.Add(30 * 24 * time.Hour) // 30 days
	case "CCPA":
		return requestDate.Add(45 * 24 * time.Hour) // 45 days
	default:
		return requestDate.Add(30 * 24 * time.Hour) // Default 30 days
	}
}

// determineDataSubjectRequestPriority determines priority of data subject request
func (k Keeper) determineDataSubjectRequestPriority(requestType types.DataSubjectRequestType) types.DataSubjectRequestPriority {
	switch requestType {
	case types.DataSubjectRequestType_ERASURE,
		types.DataSubjectRequestType_RESTRICT,
		types.DataSubjectRequestType_OBJECT:
		return types.DataSubjectRequestPriority_URGENT
	case types.DataSubjectRequestType_COMPLAINT:
		return types.DataSubjectRequestPriority_CRITICAL
	default:
		return types.DataSubjectRequestPriority_STANDARD
	}
}

// determineLegalBasis determines legal basis for processing
func (k Keeper) determineLegalBasis(requestType types.DataSubjectRequestType) string {
	switch requestType {
	case types.DataSubjectRequestType_ACCESS:
		return "GDPR Article 15 - Right of access"
	case types.DataSubjectRequestType_RECTIFICATION:
		return "GDPR Article 16 - Right to rectification"
	case types.DataSubjectRequestType_ERASURE:
		return "GDPR Article 17 - Right to erasure"
	case types.DataSubjectRequestType_RESTRICT:
		return "GDPR Article 18 - Right to restriction"
	case types.DataSubjectRequestType_PORTABILITY:
		return "GDPR Article 20 - Right to portability"
	case types.DataSubjectRequestType_OBJECT:
		return "GDPR Article 21 - Right to object"
	case types.DataSubjectRequestType_WITHDRAW_CONSENT:
		return "GDPR Article 7(3) - Withdrawal of consent"
	default:
		return "Data subject rights"
	}
}

// executeDataSubjectRequest executes the actual data subject request
func (k Keeper) executeDataSubjectRequest(ctx sdk.Context, request *types.DataSubjectRequest) (*types.DataSubjectResponse, error) {
	switch request.RequestType {
	case types.DataSubjectRequestType_ACCESS:
		return k.executeDataAccessRequest(ctx, request)
	case types.DataSubjectRequestType_ERASURE:
		return k.executeDataErasureRequest(ctx, request)
	case types.DataSubjectRequestType_PORTABILITY:
		return k.executeDataPortabilityRequest(ctx, request)
	case types.DataSubjectRequestType_RECTIFICATION:
		return k.executeDataRectificationRequest(ctx, request)
	default:
		return nil, fmt.Errorf("unsupported request type: %s", request.RequestType)
	}
}

// executeDataAccessRequest provides data access to the subject
func (k Keeper) executeDataAccessRequest(ctx sdk.Context, request *types.DataSubjectRequest) (*types.DataSubjectResponse, error) {
	// Collect all data for the subject
	identity, found := k.GetIdentity(ctx, request.DataSubject)
	if !found {
		return nil, types.ErrIdentityNotFound
	}

	credentials := k.GetCredentialsByHolder(ctx, request.DataSubject)
	consents := k.GetConsentsByHolder(ctx, request.DataSubject)
	auditEvents := k.GetAuditEventsBySubject(ctx, request.DataSubject, 
		time.Now().Add(-365*24*time.Hour), time.Now())

	responseData := map[string]interface{}{
		"identity":     identity,
		"credentials":  credentials,
		"consents":     consents,
		"audit_events": auditEvents,
		"exported_at":  ctx.BlockTime(),
	}

	response := &types.DataSubjectResponse{
		ResponseID:     fmt.Sprintf("resp_%s", request.RequestID),
		ResponseType:   "data_access",
		ResponseData:   responseData,
		ExportFormat:   "JSON",
		DeliveryMethod: "secure_download",
		DownloadCount:  0,
	}

	return response, nil
}

// executeDataErasureRequest erases data for the subject (right to be forgotten)
func (k Keeper) executeDataErasureRequest(ctx sdk.Context, request *types.DataSubjectRequest) (*types.DataSubjectResponse, error) {
	// Note: This is a simplified implementation
	// In production, you need to consider legal obligations to retain data
	
	// Mark identity as deleted but keep audit trail
	identity, found := k.GetIdentity(ctx, request.DataSubject)
	if found {
		identity.Status = types.IdentityStatus_DELETED
		k.SetIdentity(ctx, identity)
	}

	// Revoke all credentials
	credentials := k.GetCredentialsByHolder(ctx, request.DataSubject)
	for _, credential := range credentials {
		k.RevokeCredential(ctx, credential.ID, "Data erasure request")
	}

	// Withdraw all consents
	consents := k.GetConsentsByHolder(ctx, request.DataSubject)
	for _, consent := range consents {
		k.WithdrawConsent(ctx, request.DataSubject, consent.ID, "Data erasure request")
	}

	response := &types.DataSubjectResponse{
		ResponseID:     fmt.Sprintf("resp_%s", request.RequestID),
		ResponseType:   "data_erasure",
		ResponseData:   map[string]interface{}{"erased_at": ctx.BlockTime()},
		ExportFormat:   "JSON",
		DeliveryMethod: "notification",
	}

	return response, nil
}

// executeDataPortabilityRequest provides data in portable format
func (k Keeper) executeDataPortabilityRequest(ctx sdk.Context, request *types.DataSubjectRequest) (*types.DataSubjectResponse, error) {
	// Similar to access request but in standardized portable format
	identity, found := k.GetIdentity(ctx, request.DataSubject)
	if !found {
		return nil, types.ErrIdentityNotFound
	}

	credentials := k.GetCredentialsByHolder(ctx, request.DataSubject)

	// Create portable format (simplified)
	portableData := map[string]interface{}{
		"personal_data": map[string]interface{}{
			"identity":    identity,
			"credentials": credentials,
		},
		"format_version": "1.0",
		"exported_at":    ctx.BlockTime(),
	}

	response := &types.DataSubjectResponse{
		ResponseID:     fmt.Sprintf("resp_%s", request.RequestID),
		ResponseType:   "data_portability",
		ResponseData:   portableData,
		ExportFormat:   "JSON",
		DeliveryMethod: "secure_download",
		DownloadCount:  0,
	}

	return response, nil
}

// executeDataRectificationRequest corrects inaccurate data
func (k Keeper) executeDataRectificationRequest(ctx sdk.Context, request *types.DataSubjectRequest) (*types.DataSubjectResponse, error) {
	// This would need to implement specific correction logic based on request details
	// For now, return a placeholder response
	
	response := &types.DataSubjectResponse{
		ResponseID:     fmt.Sprintf("resp_%s", request.RequestID),
		ResponseType:   "data_rectification",
		ResponseData:   map[string]interface{}{"rectified_at": ctx.BlockTime()},
		ExportFormat:   "JSON",
		DeliveryMethod: "notification",
	}

	return response, nil
}

// Helper utility functions

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func containsEventType(slice []types.AuditEventType, item types.AuditEventType) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}