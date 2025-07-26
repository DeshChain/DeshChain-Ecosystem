package keeper

import (
	"context"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/DeshChain/DeshChain-Ecosystem/x/identity/types"
)

// LogAuditEvent logs an audit event to the audit trail
func (k msgServer) LogAuditEvent(goCtx context.Context, req *types.MsgLogAuditEvent) (*types.MsgLogAuditEventResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Verify authority permissions
	if err := k.ValidateAuthority(req.Authority); err != nil {
		return nil, err
	}

	// Create the audit event
	event := &types.AuditEvent{
		EventID:       k.GenerateEventID(ctx),
		Timestamp:     time.Now(),
		EventType:     req.EventType,
		Actor:         req.Actor,
		Subject:       req.Subject,
		Resource:      req.Resource,
		Action:        req.Action,
		Outcome:       req.Outcome,
		Severity:      req.Severity,
		Description:   req.Description,
		TechnicalDetails: req.TechnicalDetails,
		ComplianceFlags: req.ComplianceFlags,
		Metadata:      req.Metadata,
		IPAddress:     req.IPAddress,
		UserAgent:     req.UserAgent,
		SessionID:     req.SessionID,
		ModuleSource:  req.ModuleSource,
		ChainHeight:   ctx.BlockHeight(),
		TxHash:        fmt.Sprintf("%x", ctx.TxBytes()),
	}

	// Log the audit event
	if err := k.keeper.LogAuditEventToStore(ctx, event); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrComplianceViolation, "failed to log audit event: %v", err)
	}

	// Check for compliance violations
	if req.EventType == types.AuditEventType_COMPLIANCE_VIOLATION {
		if err := k.keeper.HandleComplianceViolation(ctx, event); err != nil {
			k.keeper.Logger(ctx).Error("Failed to handle compliance violation", "error", err)
		}
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAuditEventLogged,
			sdk.NewAttribute(types.AttributeKeyEventID, event.EventID),
			sdk.NewAttribute(types.AttributeKeyEventType, event.EventType.String()),
			sdk.NewAttribute(types.AttributeKeyActor, event.Actor),
			sdk.NewAttribute(types.AttributeKeySubject, event.Subject),
			sdk.NewAttribute(types.AttributeKeySeverity, event.Severity.String()),
		),
	)

	return &types.MsgLogAuditEventResponse{
		EventId: event.EventID,
	}, nil
}

// GenerateComplianceReport generates a comprehensive compliance report
func (k msgServer) GenerateComplianceReport(goCtx context.Context, req *types.MsgGenerateComplianceReport) (*types.MsgGenerateComplianceReportResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Verify authority permissions
	if err := k.ValidateAuthority(req.Authority); err != nil {
		return nil, err
	}

	// Generate the compliance report
	report, err := k.keeper.GenerateComplianceReportInternal(ctx, req.ReportType, req.TimeRange, req.Scope, req.Regulations)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrComplianceReportNotFound, "failed to generate compliance report: %v", err)
	}

	// Store the report
	if err := k.keeper.StoreComplianceReport(ctx, report); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrComplianceViolation, "failed to store compliance report: %v", err)
	}

	// Log audit event for report generation
	auditEvent := &types.AuditEvent{
		EventID:      k.GenerateEventID(ctx),
		Timestamp:    time.Now(),
		EventType:    types.AuditEventType_ADMIN_ACTION,
		Actor:        req.Authority,
		Subject:      "",
		Resource:     fmt.Sprintf("compliance_report:%s", report.ReportID),
		Action:       "generate_compliance_report",
		Outcome:      types.AuditOutcome_SUCCESS,
		Severity:     types.AuditSeverity_MEDIUM,
		Description:  fmt.Sprintf("Generated %s compliance report", report.ReportType),
		ModuleSource: "identity",
		ChainHeight:  ctx.BlockHeight(),
		TxHash:       fmt.Sprintf("%x", ctx.TxBytes()),
		Metadata: map[string]interface{}{
			"report_type":    report.ReportType,
			"time_range":     report.TimeRange,
			"regulations":    report.Regulations,
			"total_events":   report.TotalEvents,
			"data_subjects":  report.DataSubjects,
		},
	}

	k.keeper.LogAuditEventToStore(ctx, auditEvent)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeComplianceReportGenerated,
			sdk.NewAttribute(types.AttributeKeyReportID, report.ReportID),
			sdk.NewAttribute(types.AttributeKeyReportType, fmt.Sprintf("%d", report.ReportType)),
			sdk.NewAttribute(types.AttributeKeyGeneratedBy, report.GeneratedBy),
			sdk.NewAttribute(types.AttributeKeyTotalEvents, fmt.Sprintf("%d", report.TotalEvents)),
		),
	)

	return &types.MsgGenerateComplianceReportResponse{
		ReportId:        report.ReportID,
		ComplianceScore: report.ComplianceScore,
		RiskLevel:       report.RiskLevel,
		TotalFindings:   int64(len(report.Findings)),
	}, nil
}

// ProcessDataSubjectRequest processes a data subject request (GDPR/DPDP)
func (k msgServer) ProcessDataSubjectRequest(goCtx context.Context, req *types.MsgProcessDataSubjectRequest) (*types.MsgProcessDataSubjectRequestResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Verify authority permissions or subject identity
	if err := k.ValidateDataSubjectRequestPermissions(req.Authority, req.Request.DataSubject); err != nil {
		return nil, err
	}

	// Calculate due date based on regulation
	dueDate := k.keeper.CalculateDataSubjectRequestDueDate(req.Request.RequestType, req.Request.Regulation)
	req.Request.DueDate = dueDate
	req.Request.RequestDate = time.Now()

	// Assign request ID
	req.Request.RequestID = k.GenerateDataSubjectRequestID(ctx)

	// Process the request
	response, err := k.keeper.ProcessDataSubjectRequestInternal(ctx, req.Request)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrDataSubjectRequestNotFound, "failed to process data subject request: %v", err)
	}

	// Store the request
	if err := k.keeper.StoreDataSubjectRequest(ctx, req.Request); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrComplianceViolation, "failed to store data subject request: %v", err)
	}

	// Log audit event
	auditEvent := &types.AuditEvent{
		EventID:      k.GenerateEventID(ctx),
		Timestamp:    time.Now(),
		EventType:    k.getAuditEventTypeForDataSubjectRequest(req.Request.RequestType),
		Actor:        req.Authority,
		Subject:      req.Request.DataSubject,
		Resource:     fmt.Sprintf("data_subject_request:%s", req.Request.RequestID),
		Action:       "process_data_subject_request",
		Outcome:      types.AuditOutcome_SUCCESS,
		Severity:     types.AuditSeverity_MEDIUM,
		Description:  fmt.Sprintf("Processed %s data subject request", req.Request.RequestType),
		ModuleSource: "identity",
		ChainHeight:  ctx.BlockHeight(),
		TxHash:       fmt.Sprintf("%x", ctx.TxBytes()),
		Metadata: map[string]interface{}{
			"request_type": req.Request.RequestType,
			"regulation":   req.Request.Regulation,
			"due_date":     req.Request.DueDate,
			"priority":     req.Request.Priority,
		},
	}

	k.keeper.LogAuditEventToStore(ctx, auditEvent)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDataSubjectRequestProcessed,
			sdk.NewAttribute(types.AttributeKeyRequestID, req.Request.RequestID),
			sdk.NewAttribute(types.AttributeKeyRequestType, fmt.Sprintf("%d", req.Request.RequestType)),
			sdk.NewAttribute(types.AttributeKeyDataSubject, req.Request.DataSubject),
			sdk.NewAttribute(types.AttributeKeyRegulation, req.Request.Regulation),
		),
	)

	return &types.MsgProcessDataSubjectRequestResponse{
		RequestId:      req.Request.RequestID,
		Status:         req.Request.RequestStatus,
		DueDate:        req.Request.DueDate,
		ResponseData:   response,
		EstimatedProcessingTime: k.keeper.EstimateProcessingTime(req.Request.RequestType),
	}, nil
}

// CreatePrivacyImpactAssessment creates a new Privacy Impact Assessment
func (k msgServer) CreatePrivacyImpactAssessment(goCtx context.Context, req *types.MsgCreatePrivacyImpactAssessment) (*types.MsgCreatePrivacyImpactAssessmentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Verify authority permissions
	if err := k.ValidateAuthority(req.Authority); err != nil {
		return nil, err
	}

	// Set assessment metadata
	req.Assessment.AssessmentID = k.GeneratePIAID(ctx)
	req.Assessment.CreatedAt = time.Now()
	req.Assessment.UpdatedAt = time.Now()
	req.Assessment.CreatedBy = req.Authority
	req.Assessment.Status = types.PIAStatus_DRAFT

	// Perform initial risk assessment
	riskAssessment, err := k.keeper.PerformPrivacyRiskAssessment(ctx, req.Assessment)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrPrivacyImpactAssessmentRequired, "failed to perform risk assessment: %v", err)
	}
	req.Assessment.RiskAssessment = *riskAssessment

	// Store the PIA
	if err := k.keeper.StorePrivacyImpactAssessment(ctx, req.Assessment); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrComplianceViolation, "failed to store privacy impact assessment: %v", err)
	}

	// Log audit event
	auditEvent := &types.AuditEvent{
		EventID:      k.GenerateEventID(ctx),
		Timestamp:    time.Now(),
		EventType:    types.AuditEventType_ADMIN_ACTION,
		Actor:        req.Authority,
		Subject:      "",
		Resource:     fmt.Sprintf("privacy_impact_assessment:%s", req.Assessment.AssessmentID),
		Action:       "create_privacy_impact_assessment",
		Outcome:      types.AuditOutcome_SUCCESS,
		Severity:     types.AuditSeverity_MEDIUM,
		Description:  fmt.Sprintf("Created Privacy Impact Assessment for %s", req.Assessment.ProcessingActivity),
		ModuleSource: "identity",
		ChainHeight:  ctx.BlockHeight(),
		TxHash:       fmt.Sprintf("%x", ctx.TxBytes()),
		Metadata: map[string]interface{}{
			"processing_activity": req.Assessment.ProcessingActivity,
			"data_controller":     req.Assessment.DataController,
			"overall_risk_level":  req.Assessment.RiskAssessment.OverallRiskLevel,
		},
	}

	k.keeper.LogAuditEventToStore(ctx, auditEvent)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypePrivacyImpactAssessmentCreated,
			sdk.NewAttribute(types.AttributeKeyAssessmentID, req.Assessment.AssessmentID),
			sdk.NewAttribute(types.AttributeKeyProcessingActivity, req.Assessment.ProcessingActivity),
			sdk.NewAttribute(types.AttributeKeyDataController, req.Assessment.DataController),
			sdk.NewAttribute(types.AttributeKeyRiskLevel, fmt.Sprintf("%d", req.Assessment.RiskAssessment.OverallRiskLevel)),
		),
	)

	return &types.MsgCreatePrivacyImpactAssessmentResponse{
		AssessmentId:   req.Assessment.AssessmentID,
		RiskLevel:      req.Assessment.RiskAssessment.OverallRiskLevel,
		RequiredReview: req.Assessment.RiskAssessment.OverallRiskLevel >= types.RiskLevel_HIGH,
		ReviewDate:     req.Assessment.ReviewDate,
	}, nil
}

// UpdateAuditSettings updates audit and compliance settings
func (k msgServer) UpdateAuditSettings(goCtx context.Context, req *types.MsgUpdateAuditSettings) (*types.MsgUpdateAuditSettingsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Verify authority permissions
	if err := k.ValidateAuthority(req.Authority); err != nil {
		return nil, err
	}

	// Update audit settings
	if err := k.keeper.UpdateAuditSettings(ctx, req.Settings); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrComplianceViolation, "failed to update audit settings: %v", err)
	}

	// Log audit event
	auditEvent := &types.AuditEvent{
		EventID:      k.GenerateEventID(ctx),
		Timestamp:    time.Now(),
		EventType:    types.AuditEventType_ADMIN_ACTION,
		Actor:        req.Authority,
		Subject:      "",
		Resource:     "audit_settings",
		Action:       "update_audit_settings",
		Outcome:      types.AuditOutcome_SUCCESS,
		Severity:     types.AuditSeverity_HIGH,
		Description:  "Updated audit and compliance settings",
		ModuleSource: "identity",
		ChainHeight:  ctx.BlockHeight(),
		TxHash:       fmt.Sprintf("%x", ctx.TxBytes()),
		Metadata: map[string]interface{}{
			"retention_period":      req.Settings.RetentionPeriodDays,
			"auto_purge_enabled":    req.Settings.AutoPurgeEnabled,
			"compliance_monitoring": req.Settings.ComplianceMonitoringEnabled,
		},
	}

	k.keeper.LogAuditEventToStore(ctx, auditEvent)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAuditSettingsUpdated,
			sdk.NewAttribute(types.AttributeKeyAuthority, req.Authority),
			sdk.NewAttribute(types.AttributeKeyRetentionPeriod, fmt.Sprintf("%d", req.Settings.RetentionPeriodDays)),
		),
	)

	return &types.MsgUpdateAuditSettingsResponse{
		Success: true,
	}, nil
}

// Helper methods

// ValidateAuthority validates that the authority has permission to perform audit operations
func (k msgServer) ValidateAuthority(authority string) error {
	// Check if authority is the module authority or admin
	if authority != k.keeper.GetAuthority() {
		// Check if authority has admin privileges
		if !k.keeper.HasAdminPrivileges(authority) {
			return sdkerrors.Wrapf(types.ErrInsufficientAuditPermissions, "insufficient permissions for audit operation")
		}
	}
	return nil
}

// ValidateDataSubjectRequestPermissions validates permissions for data subject requests
func (k msgServer) ValidateDataSubjectRequestPermissions(authority, dataSubject string) error {
	// Data subjects can make requests about themselves
	if authority == dataSubject {
		return nil
	}

	// Check if authority is an admin or has data protection officer role
	if !k.keeper.HasDataProtectionOfficerRole(authority) && !k.keeper.HasAdminPrivileges(authority) {
		return sdkerrors.Wrapf(types.ErrInsufficientAuditPermissions, "insufficient permissions for data subject request")
	}

	return nil
}

// GenerateEventID generates a unique audit event ID
func (k msgServer) GenerateEventID(ctx sdk.Context) string {
	return fmt.Sprintf("audit_%d_%x", ctx.BlockHeight(), ctx.TxBytes()[:8])
}

// GenerateDataSubjectRequestID generates a unique data subject request ID
func (k msgServer) GenerateDataSubjectRequestID(ctx sdk.Context) string {
	return fmt.Sprintf("dsr_%d_%x", ctx.BlockHeight(), ctx.TxBytes()[:8])
}

// GeneratePIAID generates a unique Privacy Impact Assessment ID
func (k msgServer) GeneratePIAID(ctx sdk.Context) string {
	return fmt.Sprintf("pia_%d_%x", ctx.BlockHeight(), ctx.TxBytes()[:8])
}

// getAuditEventTypeForDataSubjectRequest maps data subject request types to audit event types
func (k msgServer) getAuditEventTypeForDataSubjectRequest(requestType types.DataSubjectRequestType) types.AuditEventType {
	switch requestType {
	case types.DataSubjectRequestType_ACCESS:
		return types.AuditEventType_DATA_ACCESS_REQUESTED
	case types.DataSubjectRequestType_ERASURE:
		return types.AuditEventType_DELETION_REQUEST
	case types.DataSubjectRequestType_PORTABILITY:
		return types.AuditEventType_EXPORT_REQUEST
	default:
		return types.AuditEventType_DATA_ACCESS_REQUESTED
	}
}