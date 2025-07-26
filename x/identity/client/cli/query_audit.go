package cli

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/DeshChain/DeshChain-Ecosystem/x/identity/types"
)

// Query commands for identity audit and compliance

// CmdQueryAuditEvents queries audit events
func CmdQueryAuditEvents() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "audit-events",
		Short: "Query audit events",
		Long: `Query audit events with optional filtering by event type, actor, subject, or time range.

Examples:
deshchaind query identity audit-events
deshchaind query identity audit-events --event-type identity_created --limit 10
deshchaind query identity audit-events --actor desh1abc... --start-time 2024-01-01T00:00:00Z`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			// Parse flags
			eventTypeStr, _ := cmd.Flags().GetString("event-type")
			actor, _ := cmd.Flags().GetString("actor")
			subject, _ := cmd.Flags().GetString("subject")
			resource, _ := cmd.Flags().GetString("resource")
			startTimeStr, _ := cmd.Flags().GetString("start-time")
			endTimeStr, _ := cmd.Flags().GetString("end-time")
			severityStr, _ := cmd.Flags().GetString("severity")
			moduleSource, _ := cmd.Flags().GetString("module-source")

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			// Parse event type filter
			var eventTypeFilter *types.AuditEventType
			if eventTypeStr != "" {
				eventType, err := parseAuditEventType(eventTypeStr)
				if err != nil {
					return fmt.Errorf("invalid event type: %w", err)
				}
				eventTypeFilter = &eventType
			}

			// Parse severity filter
			var severityFilter *types.AuditSeverity
			if severityStr != "" {
				severity, err := parseAuditSeverity(severityStr)
				if err != nil {
					return fmt.Errorf("invalid severity: %w", err)
				}
				severityFilter = &severity
			}

			// Parse time range
			var startTime, endTime *time.Time
			if startTimeStr != "" {
				t, err := time.Parse(time.RFC3339, startTimeStr)
				if err != nil {
					return fmt.Errorf("invalid start time: %w", err)
				}
				startTime = &t
			}
			if endTimeStr != "" {
				t, err := time.Parse(time.RFC3339, endTimeStr)
				if err != nil {
					return fmt.Errorf("invalid end time: %w", err)
				}
				endTime = &t
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.AuditEvents(context.Background(), &types.QueryAuditEventsRequest{
				EventType:    eventTypeFilter,
				Actor:        actor,
				Subject:      subject,
				Resource:     resource,
				StartTime:    startTime,
				EndTime:      endTime,
				Severity:     severityFilter,
				ModuleSource: moduleSource,
				Pagination:   pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().String("event-type", "", "Filter by event type")
	cmd.Flags().String("actor", "", "Filter by actor address")
	cmd.Flags().String("subject", "", "Filter by subject address")
	cmd.Flags().String("resource", "", "Filter by resource")
	cmd.Flags().String("start-time", "", "Start time (RFC3339 format)")
	cmd.Flags().String("end-time", "", "End time (RFC3339 format)")
	cmd.Flags().String("severity", "", "Filter by severity (low, medium, high, critical)")
	cmd.Flags().String("module-source", "", "Filter by module source")
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "audit events")

	return cmd
}

// CmdQueryAuditEvent queries a specific audit event
func CmdQueryAuditEvent() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "audit-event [event-id]",
		Short: "Query a specific audit event",
		Long: `Query details of a specific audit event by its ID.

Examples:
deshchaind query identity audit-event audit_12345_abcdef
deshchaind query identity audit-event audit_67890_fedcba`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			eventID := args[0]

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.AuditEvent(context.Background(), &types.QueryAuditEventRequest{
				EventId: eventID,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryComplianceReports queries compliance reports
func CmdQueryComplianceReports() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compliance-reports",
		Short: "Query compliance reports",
		Long: `Query compliance reports with optional filtering by type and time range.

Examples:
deshchaind query identity compliance-reports
deshchaind query identity compliance-reports --report-type gdpr_compliance
deshchaind query identity compliance-reports --generated-by desh1abc...`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			// Parse flags
			reportTypeStr, _ := cmd.Flags().GetString("report-type")
			generatedBy, _ := cmd.Flags().GetString("generated-by")
			startTimeStr, _ := cmd.Flags().GetString("start-time")
			endTimeStr, _ := cmd.Flags().GetString("end-time")

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			// Parse report type filter
			var reportTypeFilter *types.ComplianceReportType
			if reportTypeStr != "" {
				reportType, err := parseComplianceReportType(reportTypeStr)
				if err != nil {
					return fmt.Errorf("invalid report type: %w", err)
				}
				reportTypeFilter = &reportType
			}

			// Parse time range
			var startTime, endTime *time.Time
			if startTimeStr != "" {
				t, err := time.Parse(time.RFC3339, startTimeStr)
				if err != nil {
					return fmt.Errorf("invalid start time: %w", err)
				}
				startTime = &t
			}
			if endTimeStr != "" {
				t, err := time.Parse(time.RFC3339, endTimeStr)
				if err != nil {
					return fmt.Errorf("invalid end time: %w", err)
				}
				endTime = &t
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.ComplianceReports(context.Background(), &types.QueryComplianceReportsRequest{
				ReportType:  reportTypeFilter,
				GeneratedBy: generatedBy,
				StartTime:   startTime,
				EndTime:     endTime,
				Pagination:  pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().String("report-type", "", "Filter by report type")
	cmd.Flags().String("generated-by", "", "Filter by generator address")
	cmd.Flags().String("start-time", "", "Start time (RFC3339 format)")
	cmd.Flags().String("end-time", "", "End time (RFC3339 format)")
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "compliance reports")

	return cmd
}

// CmdQueryComplianceReport queries a specific compliance report
func CmdQueryComplianceReport() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compliance-report [report-id]",
		Short: "Query a specific compliance report",
		Long: `Query details of a specific compliance report by its ID.

Examples:
deshchaind query identity compliance-report report_12345_abcdef
deshchaind query identity compliance-report report_67890_fedcba`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			reportID := args[0]

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.ComplianceReport(context.Background(), &types.QueryComplianceReportRequest{
				ReportId: reportID,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryDataSubjectRequests queries data subject requests
func CmdQueryDataSubjectRequests() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "data-subject-requests",
		Short: "Query data subject requests",
		Long: `Query data subject requests with optional filtering.

Examples:
deshchaind query identity data-subject-requests
deshchaind query identity data-subject-requests --data-subject desh1abc...
deshchaind query identity data-subject-requests --request-type access --status completed`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			// Parse flags
			dataSubject, _ := cmd.Flags().GetString("data-subject")
			requestTypeStr, _ := cmd.Flags().GetString("request-type")
			statusStr, _ := cmd.Flags().GetString("status")
			regulation, _ := cmd.Flags().GetString("regulation")

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			// Parse request type filter
			var requestTypeFilter *types.DataSubjectRequestType
			if requestTypeStr != "" {
				requestType, err := parseDataSubjectRequestType(requestTypeStr)
				if err != nil {
					return fmt.Errorf("invalid request type: %w", err)
				}
				requestTypeFilter = &requestType
			}

			// Parse status filter
			var statusFilter *types.DataSubjectRequestStatus
			if statusStr != "" {
				status, err := parseDataSubjectRequestStatus(statusStr)
				if err != nil {
					return fmt.Errorf("invalid status: %w", err)
				}
				statusFilter = &status
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.DataSubjectRequests(context.Background(), &types.QueryDataSubjectRequestsRequest{
				DataSubject: dataSubject,
				RequestType: requestTypeFilter,
				Status:      statusFilter,
				Regulation:  regulation,
				Pagination:  pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().String("data-subject", "", "Filter by data subject address")
	cmd.Flags().String("request-type", "", "Filter by request type")
	cmd.Flags().String("status", "", "Filter by request status")
	cmd.Flags().String("regulation", "", "Filter by regulation")
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "data subject requests")

	return cmd
}

// CmdQueryDataSubjectRequest queries a specific data subject request
func CmdQueryDataSubjectRequest() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "data-subject-request [request-id]",
		Short: "Query a specific data subject request",
		Long: `Query details of a specific data subject request by its ID.

Examples:
deshchaind query identity data-subject-request dsr_12345_abcdef
deshchaind query identity data-subject-request dsr_67890_fedcba`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			requestID := args[0]

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.DataSubjectRequest(context.Background(), &types.QueryDataSubjectRequestRequest{
				RequestId: requestID,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryPrivacyImpactAssessments queries Privacy Impact Assessments
func CmdQueryPrivacyImpactAssessments() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "privacy-impact-assessments",
		Short: "Query Privacy Impact Assessments",
		Long: `Query Privacy Impact Assessments with optional filtering.

Examples:
deshchaind query identity privacy-impact-assessments
deshchaind query identity privacy-impact-assessments --data-controller "DeshChain Foundation"
deshchaind query identity privacy-impact-assessments --status approved --risk-level high`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			// Parse flags
			dataController, _ := cmd.Flags().GetString("data-controller")
			statusStr, _ := cmd.Flags().GetString("status")
			riskLevelStr, _ := cmd.Flags().GetString("risk-level")

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			// Parse status filter
			var statusFilter *types.PIAStatus
			if statusStr != "" {
				status, err := parsePIAStatus(statusStr)
				if err != nil {
					return fmt.Errorf("invalid status: %w", err)
				}
				statusFilter = &status
			}

			// Parse risk level filter
			var riskLevelFilter *types.PrivacyRiskLevel
			if riskLevelStr != "" {
				riskLevel, err := parsePrivacyRiskLevel(riskLevelStr)
				if err != nil {
					return fmt.Errorf("invalid risk level: %w", err)
				}
				riskLevelFilter = &riskLevel
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.PrivacyImpactAssessments(context.Background(), &types.QueryPrivacyImpactAssessmentsRequest{
				DataController: dataController,
				Status:         statusFilter,
				RiskLevel:      riskLevelFilter,
				Pagination:     pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().String("data-controller", "", "Filter by data controller")
	cmd.Flags().String("status", "", "Filter by assessment status")
	cmd.Flags().String("risk-level", "", "Filter by risk level")
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "privacy impact assessments")

	return cmd
}

// CmdQueryPrivacyImpactAssessment queries a specific Privacy Impact Assessment
func CmdQueryPrivacyImpactAssessment() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "privacy-impact-assessment [assessment-id]",
		Short: "Query a specific Privacy Impact Assessment",
		Long: `Query details of a specific Privacy Impact Assessment by its ID.

Examples:
deshchaind query identity privacy-impact-assessment pia_12345_abcdef
deshchaind query identity privacy-impact-assessment pia_67890_fedcba`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			assessmentID := args[0]

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.PrivacyImpactAssessment(context.Background(), &types.QueryPrivacyImpactAssessmentRequest{
				AssessmentId: assessmentID,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryAuditSettings queries current audit settings
func CmdQueryAuditSettings() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "audit-settings",
		Short: "Query current audit and compliance settings",
		Long: `Query the current configuration settings for audit and compliance monitoring.

Examples:
deshchaind query identity audit-settings`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.AuditSettings(context.Background(), &types.QueryAuditSettingsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryComplianceStatistics queries compliance statistics
func CmdQueryComplianceStatistics() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compliance-statistics [time-window-days]",
		Short: "Query compliance statistics for a time period",
		Long: `Query compliance statistics including violation counts, report generation stats, and trends.

Examples:
deshchaind query identity compliance-statistics 30  # Last 30 days
deshchaind query identity compliance-statistics 365 # Last year`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			timeWindowDays, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid time window days: %w", err)
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.ComplianceStatistics(context.Background(), &types.QueryComplianceStatisticsRequest{
				TimeWindowDays: timeWindowDays,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// Helper functions for parsing additional enums

func parseDataSubjectRequestStatus(s string) (types.DataSubjectRequestStatus, error) {
	switch s {
	case "received":
		return types.DataSubjectRequestStatus_RECEIVED, nil
	case "acknowledged":
		return types.DataSubjectRequestStatus_ACKNOWLEDGED, nil
	case "in_review":
		return types.DataSubjectRequestStatus_IN_REVIEW, nil
	case "in_progress":
		return types.DataSubjectRequestStatus_IN_PROGRESS, nil
	case "completed":
		return types.DataSubjectRequestStatus_COMPLETED, nil
	case "rejected":
		return types.DataSubjectRequestStatus_REJECTED, nil
	case "escalated":
		return types.DataSubjectRequestStatus_ESCALATED, nil
	case "expired":
		return types.DataSubjectRequestStatus_EXPIRED, nil
	default:
		return 0, fmt.Errorf("unknown data subject request status: %s", s)
	}
}

func parsePIAStatus(s string) (types.PIAStatus, error) {
	switch s {
	case "draft":
		return types.PIAStatus_DRAFT, nil
	case "review":
		return types.PIAStatus_REVIEW, nil
	case "approved":
		return types.PIAStatus_APPROVED, nil
	case "rejected":
		return types.PIAStatus_REJECTED, nil
	case "superseded":
		return types.PIAStatus_SUPERSEDED, nil
	default:
		return 0, fmt.Errorf("unknown PIA status: %s", s)
	}
}

func parsePrivacyRiskLevel(s string) (types.PrivacyRiskLevel, error) {
	switch s {
	case "very_low":
		return types.RiskLevel_VERY_LOW, nil
	case "low":
		return types.RiskLevel_LOW, nil
	case "medium":
		return types.RiskLevel_MEDIUM, nil
	case "high":
		return types.RiskLevel_HIGH, nil
	case "very_high":
		return types.RiskLevel_VERY_HIGH, nil
	default:
		return 0, fmt.Errorf("unknown privacy risk level: %s", s)
	}
}