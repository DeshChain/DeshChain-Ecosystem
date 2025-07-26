package cli

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/namo/x/identity/types"
)

// Transaction commands for identity audit and compliance

// CmdLogAuditEvent logs an audit event
func CmdLogAuditEvent() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "log-audit-event [event-type] [actor] [subject] [resource] [action] [outcome] [description]",
		Short: "Log an audit event",
		Long: `Log an audit event to the identity audit trail.

Examples:
deshchaind tx identity log-audit-event identity_created desh1abc... desh1def... identity_123 create_identity success "Created new identity" --from admin-key
deshchaind tx identity log-audit-event credential_issued desh1abc... desh1def... cred_123 issue_credential success "Issued KYC credential" --severity high --from admin-key`,
		Args: cobra.ExactArgs(7),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Parse event type
			eventTypeStr := args[0]
			eventType, err := parseAuditEventType(eventTypeStr)
			if err != nil {
				return fmt.Errorf("invalid event type: %w", err)
			}

			// Parse outcome
			outcomeStr := args[5]
			outcome, err := parseAuditOutcome(outcomeStr)
			if err != nil {
				return fmt.Errorf("invalid outcome: %w", err)
			}

			// Parse flags
			severityStr, _ := cmd.Flags().GetString("severity")
			severity, err := parseAuditSeverity(severityStr)
			if err != nil {
				return fmt.Errorf("invalid severity: %w", err)
			}

			ipAddress, _ := cmd.Flags().GetString("ip-address")
			userAgent, _ := cmd.Flags().GetString("user-agent")
			sessionID, _ := cmd.Flags().GetString("session-id")
			moduleSource, _ := cmd.Flags().GetString("module-source")
			metadataStr, _ := cmd.Flags().GetString("metadata")
			technicalDetailsStr, _ := cmd.Flags().GetString("technical-details")
			complianceFlagsStr, _ := cmd.Flags().GetString("compliance-flags")

			// Parse metadata
			var metadata map[string]interface{}
			if metadataStr != "" {
				if err := json.Unmarshal([]byte(metadataStr), &metadata); err != nil {
					return fmt.Errorf("invalid metadata JSON: %w", err)
				}
			}

			// Parse technical details
			var technicalDetails map[string]interface{}
			if technicalDetailsStr != "" {
				if err := json.Unmarshal([]byte(technicalDetailsStr), &technicalDetails); err != nil {
					return fmt.Errorf("invalid technical details JSON: %w", err)
				}
			}

			// Parse compliance flags
			var complianceFlags []types.ComplianceFlag
			if complianceFlagsStr != "" {
				if err := json.Unmarshal([]byte(complianceFlagsStr), &complianceFlags); err != nil {
					return fmt.Errorf("invalid compliance flags JSON: %w", err)
				}
			}

			msg := &types.MsgLogAuditEvent{
				Authority:        clientCtx.GetFromAddress().String(),
				EventType:        eventType,
				Actor:            args[1],
				Subject:          args[2],
				Resource:         args[3],
				Action:           args[4],
				Outcome:          outcome,
				Severity:         severity,
				Description:      args[6],
				TechnicalDetails: technicalDetails,
				ComplianceFlags:  complianceFlags,
				Metadata:         metadata,
				IPAddress:        ipAddress,
				UserAgent:        userAgent,
				SessionID:        sessionID,
				ModuleSource:     moduleSource,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("severity", "medium", "Event severity (low, medium, high, critical)")
	cmd.Flags().String("ip-address", "", "IP address of the actor")
	cmd.Flags().String("user-agent", "", "User agent string")
	cmd.Flags().String("session-id", "", "Session ID")
	cmd.Flags().String("module-source", "identity", "Source module")
	cmd.Flags().String("metadata", "", "Additional metadata as JSON")
	cmd.Flags().String("technical-details", "", "Technical details as JSON")
	cmd.Flags().String("compliance-flags", "", "Compliance flags as JSON array")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// CmdGenerateComplianceReport generates a compliance report
func CmdGenerateComplianceReport() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-compliance-report [report-type] [start-time] [end-time]",
		Short: "Generate a compliance report",
		Long: `Generate a comprehensive compliance report for audit and regulatory purposes.

Examples:
deshchaind tx identity generate-compliance-report gdpr_compliance 2024-01-01T00:00:00Z 2024-12-31T23:59:59Z --from admin-key
deshchaind tx identity generate-compliance-report security_audit 2024-06-01T00:00:00Z 2024-06-30T23:59:59Z --regulations "GDPR,DPDP" --from admin-key`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Parse report type
			reportTypeStr := args[0]
			reportType, err := parseComplianceReportType(reportTypeStr)
			if err != nil {
				return fmt.Errorf("invalid report type: %w", err)
			}

			// Parse time range
			startTime, err := time.Parse(time.RFC3339, args[1])
			if err != nil {
				return fmt.Errorf("invalid start time: %w", err)
			}

			endTime, err := time.Parse(time.RFC3339, args[2])
			if err != nil {
				return fmt.Errorf("invalid end time: %w", err)
			}

			timeRange := types.AuditTimeRange{
				StartTime: startTime,
				EndTime:   endTime,
			}

			// Parse flags
			regulationsStr, _ := cmd.Flags().GetString("regulations")
			includeModulesStr, _ := cmd.Flags().GetString("include-modules")
			excludeModulesStr, _ := cmd.Flags().GetString("exclude-modules")

			var regulations []string
			if regulationsStr != "" {
				regulations = strings.Split(regulationsStr, ",")
				for i, reg := range regulations {
					regulations[i] = strings.TrimSpace(reg)
				}
			}

			var includeModules []string
			if includeModulesStr != "" {
				includeModules = strings.Split(includeModulesStr, ",")
				for i, mod := range includeModules {
					includeModules[i] = strings.TrimSpace(mod)
				}
			}

			var excludeModules []string
			if excludeModulesStr != "" {
				excludeModules = strings.Split(excludeModulesStr, ",")
				for i, mod := range excludeModules {
					excludeModules[i] = strings.TrimSpace(mod)
				}
			}

			scope := types.ComplianceScope{
				IncludeModules: includeModules,
				ExcludeModules: excludeModules,
			}

			msg := &types.MsgGenerateComplianceReport{
				Authority:   clientCtx.GetFromAddress().String(),
				ReportType:  reportType,
				TimeRange:   timeRange,
				Scope:       scope,
				Regulations: regulations,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("regulations", "", "Comma-separated list of regulations (GDPR, DPDP, CCPA)")
	cmd.Flags().String("include-modules", "", "Comma-separated list of modules to include")
	cmd.Flags().String("exclude-modules", "", "Comma-separated list of modules to exclude")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// CmdProcessDataSubjectRequest processes a data subject request
func CmdProcessDataSubjectRequest() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "process-data-subject-request [request-type] [data-subject] [description]",
		Short: "Process a data subject request (GDPR/DPDP)",
		Long: `Process a data subject request for access, rectification, erasure, or other rights.

Examples:
deshchaind tx identity process-data-subject-request access desh1abc... "Request access to all personal data" --regulation GDPR --from admin-key
deshchaind tx identity process-data-subject-request erasure desh1def... "Request deletion of all personal data" --regulation DPDP --priority urgent --from admin-key`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Parse request type
			requestTypeStr := args[0]
			requestType, err := parseDataSubjectRequestType(requestTypeStr)
			if err != nil {
				return fmt.Errorf("invalid request type: %w", err)
			}

			// Parse flags
			regulation, _ := cmd.Flags().GetString("regulation")
			priorityStr, _ := cmd.Flags().GetString("priority")
			requestedBy, _ := cmd.Flags().GetString("requested-by")
			legalBasis, _ := cmd.Flags().GetString("legal-basis")
			requestDetailsStr, _ := cmd.Flags().GetString("request-details")

			priority, err := parseDataSubjectRequestPriority(priorityStr)
			if err != nil {
				return fmt.Errorf("invalid priority: %w", err)
			}

			// Parse request details
			var requestDetails map[string]interface{}
			if requestDetailsStr != "" {
				if err := json.Unmarshal([]byte(requestDetailsStr), &requestDetails); err != nil {
					return fmt.Errorf("invalid request details JSON: %w", err)
				}
			}

			// If requested-by is not specified, use data subject
			if requestedBy == "" {
				requestedBy = args[1]
			}

			request := &types.DataSubjectRequest{
				RequestType:    requestType,
				RequestStatus:  types.DataSubjectRequestStatus_RECEIVED,
				DataSubject:    args[1],
				RequestedBy:    requestedBy,
				Description:    args[2],
				RequestDetails: requestDetails,
				LegalBasis:     legalBasis,
				Regulation:     regulation,
				Priority:       priority,
			}

			msg := &types.MsgProcessDataSubjectRequest{
				Authority: clientCtx.GetFromAddress().String(),
				Request:   request,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("regulation", "GDPR", "Applicable regulation (GDPR, DPDP, CCPA)")
	cmd.Flags().String("priority", "standard", "Request priority (standard, urgent, critical)")
	cmd.Flags().String("requested-by", "", "Address of the person making the request (if different from data subject)")
	cmd.Flags().String("legal-basis", "", "Legal basis for the request")
	cmd.Flags().String("request-details", "", "Additional request details as JSON")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// CmdCreatePrivacyImpactAssessment creates a Privacy Impact Assessment
func CmdCreatePrivacyImpactAssessment() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-privacy-impact-assessment [title] [processing-activity] [data-controller]",
		Short: "Create a Privacy Impact Assessment (DPIA)",
		Long: `Create a Privacy Impact Assessment for data processing activities.

Examples:
deshchaind tx identity create-privacy-impact-assessment "KYC Data Processing" "kyc_verification" "DeshChain Foundation" --from admin-key
deshchaind tx identity create-privacy-impact-assessment "Biometric Authentication" "biometric_auth" "DeshChain Foundation" --data-categories "biometric,identity" --from admin-key`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Parse flags
			description, _ := cmd.Flags().GetString("description")
			dataProcessor, _ := cmd.Flags().GetString("data-processor")
			legalBasisStr, _ := cmd.Flags().GetString("legal-basis")
			dataCategoriesStr, _ := cmd.Flags().GetString("data-categories")
			dataSubjectsStr, _ := cmd.Flags().GetString("data-subjects")
			processingPurposesStr, _ := cmd.Flags().GetString("processing-purposes")
			retentionPeriod, _ := cmd.Flags().GetString("retention-period")
			technicalMeasuresStr, _ := cmd.Flags().GetString("technical-measures")
			organizationalMeasuresStr, _ := cmd.Flags().GetString("organizational-measures")

			var legalBasis []string
			if legalBasisStr != "" {
				legalBasis = strings.Split(legalBasisStr, ",")
				for i, basis := range legalBasis {
					legalBasis[i] = strings.TrimSpace(basis)
				}
			}

			var dataCategories []string
			if dataCategoriesStr != "" {
				dataCategories = strings.Split(dataCategoriesStr, ",")
				for i, cat := range dataCategories {
					dataCategories[i] = strings.TrimSpace(cat)
				}
			}

			var dataSubjects []string
			if dataSubjectsStr != "" {
				dataSubjects = strings.Split(dataSubjectsStr, ",")
				for i, sub := range dataSubjects {
					dataSubjects[i] = strings.TrimSpace(sub)
				}
			}

			var processingPurposes []string
			if processingPurposesStr != "" {
				processingPurposes = strings.Split(processingPurposesStr, ",")
				for i, purpose := range processingPurposes {
					processingPurposes[i] = strings.TrimSpace(purpose)
				}
			}

			var technicalMeasures []string
			if technicalMeasuresStr != "" {
				technicalMeasures = strings.Split(technicalMeasuresStr, ",")
				for i, measure := range technicalMeasures {
					technicalMeasures[i] = strings.TrimSpace(measure)
				}
			}

			var organizationalMeasures []string
			if organizationalMeasuresStr != "" {
				organizationalMeasures = strings.Split(organizationalMeasuresStr, ",")
				for i, measure := range organizationalMeasures {
					organizationalMeasures[i] = strings.TrimSpace(measure)
				}
			}

			assessment := &types.PrivacyImpactAssessment{
				Title:                  args[0],
				Description:            description,
				ProcessingActivity:     args[1],
				DataController:         args[2],
				DataProcessor:          dataProcessor,
				LegalBasis:             legalBasis,
				DataCategories:         dataCategories,
				DataSubjects:           dataSubjects,
				ProcessingPurposes:     processingPurposes,
				RetentionPeriod:        retentionPeriod,
				TechnicalMeasures:      technicalMeasures,
				OrganizationalMeasures: organizationalMeasures,
				ReviewDate:             time.Now().AddDate(1, 0, 0), // Review in 1 year
			}

			msg := &types.MsgCreatePrivacyImpactAssessment{
				Authority:  clientCtx.GetFromAddress().String(),
				Assessment: assessment,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("description", "", "Detailed description of the assessment")
	cmd.Flags().String("data-processor", "", "Data processor (if different from controller)")
	cmd.Flags().String("legal-basis", "", "Comma-separated list of legal basis")
	cmd.Flags().String("data-categories", "", "Comma-separated list of data categories")
	cmd.Flags().String("data-subjects", "", "Comma-separated list of data subject types")
	cmd.Flags().String("processing-purposes", "", "Comma-separated list of processing purposes")
	cmd.Flags().String("retention-period", "", "Data retention period")
	cmd.Flags().String("technical-measures", "", "Comma-separated list of technical measures")
	cmd.Flags().String("organizational-measures", "", "Comma-separated list of organizational measures")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// CmdUpdateAuditSettings updates audit and compliance settings
func CmdUpdateAuditSettings() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-audit-settings",
		Short: "Update audit and compliance settings",
		Long: `Update the configuration settings for audit and compliance monitoring.

Examples:
deshchaind tx identity update-audit-settings --retention-period 2555 --auto-purge --compliance-monitoring --from admin-key`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			retentionPeriod, _ := cmd.Flags().GetInt64("retention-period")
			autoPurge, _ := cmd.Flags().GetBool("auto-purge")
			complianceMonitoring, _ := cmd.Flags().GetBool("compliance-monitoring")
			realTimeAlerts, _ := cmd.Flags().GetBool("real-time-alerts")
			anonymizeExpiredData, _ := cmd.Flags().GetBool("anonymize-expired")

			settings := types.AuditSettings{
				RetentionPeriodDays:         retentionPeriod,
				AutoPurgeEnabled:           autoPurge,
				ComplianceMonitoringEnabled: complianceMonitoring,
				RealTimeAlertsEnabled:      realTimeAlerts,
				AnonymizeExpiredData:       anonymizeExpiredData,
			}

			msg := &types.MsgUpdateAuditSettings{
				Authority: clientCtx.GetFromAddress().String(),
				Settings:  settings,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().Int64("retention-period", 2555, "Audit data retention period in days (7 years default)")
	cmd.Flags().Bool("auto-purge", true, "Enable automatic purging of expired audit data")
	cmd.Flags().Bool("compliance-monitoring", true, "Enable real-time compliance monitoring")
	cmd.Flags().Bool("real-time-alerts", true, "Enable real-time alerts for violations")
	cmd.Flags().Bool("anonymize-expired", false, "Anonymize data instead of deleting when expired")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// Helper functions for parsing enums

func parseAuditEventType(s string) (types.AuditEventType, error) {
	switch strings.ToLower(s) {
	case "identity_created":
		return types.AuditEventType_IDENTITY_CREATED, nil
	case "identity_updated":
		return types.AuditEventType_IDENTITY_UPDATED, nil
	case "identity_deleted":
		return types.AuditEventType_IDENTITY_DELETED, nil
	case "identity_accessed":
		return types.AuditEventType_IDENTITY_ACCESSED, nil
	case "credential_issued":
		return types.AuditEventType_CREDENTIAL_ISSUED, nil
	case "credential_verified":
		return types.AuditEventType_CREDENTIAL_VERIFIED, nil
	case "credential_revoked":
		return types.AuditEventType_CREDENTIAL_REVOKED, nil
	case "credential_accessed":
		return types.AuditEventType_CREDENTIAL_ACCESSED, nil
	case "consent_given":
		return types.AuditEventType_CONSENT_GIVEN, nil
	case "consent_withdrawn":
		return types.AuditEventType_CONSENT_WITHDRAWN, nil
	case "consent_accessed":
		return types.AuditEventType_CONSENT_ACCESSED, nil
	case "data_shared":
		return types.AuditEventType_DATA_SHARED, nil
	case "data_access_requested":
		return types.AuditEventType_DATA_ACCESS_REQUESTED, nil
	case "data_access_denied":
		return types.AuditEventType_DATA_ACCESS_DENIED, nil
	case "privacy_settings_changed":
		return types.AuditEventType_PRIVACY_SETTINGS_CHANGED, nil
	case "did_created":
		return types.AuditEventType_DID_CREATED, nil
	case "did_updated":
		return types.AuditEventType_DID_UPDATED, nil
	case "did_deactivated":
		return types.AuditEventType_DID_DEACTIVATED, nil
	case "biometric_enrolled":
		return types.AuditEventType_BIOMETRIC_ENROLLED, nil
	case "biometric_verified":
		return types.AuditEventType_BIOMETRIC_VERIFIED, nil
	case "kyc_initiated":
		return types.AuditEventType_KYC_INITIATED, nil
	case "kyc_completed":
		return types.AuditEventType_KYC_COMPLETED, nil
	case "recovery_initiated":
		return types.AuditEventType_RECOVERY_INITIATED, nil
	case "recovery_completed":
		return types.AuditEventType_RECOVERY_COMPLETED, nil
	case "suspicious_activity":
		return types.AuditEventType_SUSPICIOUS_ACTIVITY, nil
	case "compliance_violation":
		return types.AuditEventType_COMPLIANCE_VIOLATION, nil
	case "system_error":
		return types.AuditEventType_SYSTEM_ERROR, nil
	case "admin_action":
		return types.AuditEventType_ADMIN_ACTION, nil
	case "export_request":
		return types.AuditEventType_EXPORT_REQUEST, nil
	case "deletion_request":
		return types.AuditEventType_DELETION_REQUEST, nil
	default:
		return 0, fmt.Errorf("unknown audit event type: %s", s)
	}
}

func parseAuditOutcome(s string) (types.AuditOutcome, error) {
	switch strings.ToLower(s) {
	case "success":
		return types.AuditOutcome_SUCCESS, nil
	case "failure":
		return types.AuditOutcome_FAILURE, nil
	case "partial_success":
		return types.AuditOutcome_PARTIAL_SUCCESS, nil
	case "denied":
		return types.AuditOutcome_DENIED, nil
	case "error":
		return types.AuditOutcome_ERROR, nil
	case "timeout":
		return types.AuditOutcome_TIMEOUT, nil
	case "cancelled":
		return types.AuditOutcome_CANCELLED, nil
	default:
		return 0, fmt.Errorf("unknown audit outcome: %s", s)
	}
}

func parseAuditSeverity(s string) (types.AuditSeverity, error) {
	switch strings.ToLower(s) {
	case "low":
		return types.AuditSeverity_LOW, nil
	case "medium":
		return types.AuditSeverity_MEDIUM, nil
	case "high":
		return types.AuditSeverity_HIGH, nil
	case "critical":
		return types.AuditSeverity_CRITICAL, nil
	default:
		return types.AuditSeverity_MEDIUM, nil // Default to medium
	}
}

func parseComplianceReportType(s string) (types.ComplianceReportType, error) {
	switch strings.ToLower(s) {
	case "gdpr_compliance":
		return types.ComplianceReportType_GDPR_COMPLIANCE, nil
	case "dpdp_compliance":
		return types.ComplianceReportType_DPDP_COMPLIANCE, nil
	case "ccpa_compliance":
		return types.ComplianceReportType_CCPA_COMPLIANCE, nil
	case "general_audit":
		return types.ComplianceReportType_GENERAL_AUDIT, nil
	case "security_audit":
		return types.ComplianceReportType_SECURITY_AUDIT, nil
	case "data_breach_report":
		return types.ComplianceReportType_DATA_BREACH_REPORT, nil
	case "periodic_review":
		return types.ComplianceReportType_PERIODIC_REVIEW, nil
	case "incident_report":
		return types.ComplianceReportType_INCIDENT_REPORT, nil
	case "risk_assessment":
		return types.ComplianceReportType_RISK_ASSESSMENT, nil
	default:
		return 0, fmt.Errorf("unknown compliance report type: %s", s)
	}
}

func parseDataSubjectRequestType(s string) (types.DataSubjectRequestType, error) {
	switch strings.ToLower(s) {
	case "access":
		return types.DataSubjectRequestType_ACCESS, nil
	case "rectification":
		return types.DataSubjectRequestType_RECTIFICATION, nil
	case "erasure":
		return types.DataSubjectRequestType_ERASURE, nil
	case "restrict":
		return types.DataSubjectRequestType_RESTRICT, nil
	case "portability":
		return types.DataSubjectRequestType_PORTABILITY, nil
	case "object":
		return types.DataSubjectRequestType_OBJECT, nil
	case "withdraw_consent":
		return types.DataSubjectRequestType_WITHDRAW_CONSENT, nil
	case "complaint":
		return types.DataSubjectRequestType_COMPLAINT, nil
	case "information":
		return types.DataSubjectRequestType_INFORMATION, nil
	case "stop_processing":
		return types.DataSubjectRequestType_STOP_PROCESSING, nil
	default:
		return 0, fmt.Errorf("unknown data subject request type: %s", s)
	}
}

func parseDataSubjectRequestPriority(s string) (types.DataSubjectRequestPriority, error) {
	switch strings.ToLower(s) {
	case "standard":
		return types.DataSubjectRequestPriority_STANDARD, nil
	case "urgent":
		return types.DataSubjectRequestPriority_URGENT, nil
	case "critical":
		return types.DataSubjectRequestPriority_CRITICAL, nil
	default:
		return types.DataSubjectRequestPriority_STANDARD, nil // Default to standard
	}
}