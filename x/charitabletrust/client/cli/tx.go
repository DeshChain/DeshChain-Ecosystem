package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/DeshChain/DeshChain-Ecosystem/x/charitabletrust/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdCreateAllocationProposal(),
		CmdVoteOnProposal(),
		CmdExecuteAllocation(),
		CmdSubmitImpactReport(),
		CmdVerifyImpactReport(),
		CmdReportFraud(),
		CmdInvestigateFraud(),
	)

	return cmd
}

// CmdCreateAllocationProposal returns the create allocation proposal command
func CmdCreateAllocationProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-proposal [title] [description] [total-amount] [justification] [expected-impact]",
		Short: "Create a new allocation proposal for charitable distributions",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			title := args[0]
			description := args[1]
			totalAmount, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}
			justification := args[3]
			expectedImpact := args[4]

			// Parse allocations from flags
			allocationsStr, _ := cmd.Flags().GetString("allocations")
			var allocations []types.ProposedAllocation
			
			if allocationsStr != "" {
				allocationParts := strings.Split(allocationsStr, ";")
				for _, part := range allocationParts {
					fields := strings.Split(part, ",")
					if len(fields) != 5 {
						return fmt.Errorf("invalid allocation format: %s", part)
					}
					
					orgID, err := strconv.ParseUint(fields[0], 10, 64)
					if err != nil {
						return err
					}
					
					amount, err := sdk.ParseCoinNormalized(fields[2])
					if err != nil {
						return err
					}
					
					allocations = append(allocations, types.ProposedAllocation{
						CharitableOrgWalletId: orgID,
						OrganizationName:      fields[1],
						Amount:                amount,
						Purpose:               fields[3],
						Category:              fields[4],
					})
				}
			}

			documentsStr, _ := cmd.Flags().GetString("documents")
			documents := []string{}
			if documentsStr != "" {
				documents = strings.Split(documentsStr, ",")
			}

			msg := &types.MsgCreateAllocationProposal{
				Proposer:       clientCtx.GetFromAddress().String(),
				Title:          title,
				Description:    description,
				TotalAmount:    totalAmount,
				Allocations:    allocations,
				Justification:  justification,
				ExpectedImpact: expectedImpact,
				Documents:      documents,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().String("allocations", "", "Semicolon-separated allocations (orgID,orgName,amount,purpose,category)")
	cmd.Flags().String("documents", "", "Comma-separated list of document URLs")
	cmd.MarkFlagRequired("allocations")

	return cmd
}

// CmdVoteOnProposal returns the vote on proposal command
func CmdVoteOnProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote-proposal [proposal-id] [vote] [reason]",
		Short: "Vote on an allocation proposal",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			vote := args[1]
			reason := args[2]

			msg := &types.MsgVoteOnProposal{
				ProposalId: proposalID,
				Voter:      clientCtx.GetFromAddress().String(),
				Vote:       vote,
				Reason:     reason,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdExecuteAllocation returns the execute allocation command
func CmdExecuteAllocation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "execute-allocation [proposal-id]",
		Short: "Execute an approved allocation proposal",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := &types.MsgExecuteAllocation{
				ProposalId: proposalID,
				Executor:   clientCtx.GetFromAddress().String(),
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdSubmitImpactReport returns the submit impact report command
func CmdSubmitImpactReport() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-impact-report [allocation-id] [period] [beneficiaries-reached] [funds-utilized]",
		Short: "Submit an impact report for an allocation",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			allocationID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			period := args[1]
			beneficiariesReached, err := strconv.ParseInt(args[2], 10, 32)
			if err != nil {
				return err
			}
			fundsUtilized, err := sdk.ParseCoinNormalized(args[3])
			if err != nil {
				return err
			}

			// Parse metrics from flags
			metricsStr, _ := cmd.Flags().GetString("metrics")
			var metrics []types.ImpactMetric
			
			if metricsStr != "" {
				metricParts := strings.Split(metricsStr, ";")
				for _, part := range metricParts {
					fields := strings.Split(part, ",")
					if len(fields) != 4 {
						return fmt.Errorf("invalid metric format: %s", part)
					}
					
					target, err := strconv.ParseInt(fields[1], 10, 32)
					if err != nil {
						return err
					}
					
					achieved, err := strconv.ParseInt(fields[2], 10, 32)
					if err != nil {
						return err
					}
					
					percentage, err := sdk.NewDecFromStr(fields[3])
					if err != nil {
						return err
					}
					
					metrics = append(metrics, types.ImpactMetric{
						MetricName:            fields[0],
						TargetValue:           int32(target),
						AchievedValue:         int32(achieved),
						AchievementPercentage: percentage,
					})
				}
			}

			documentsStr, _ := cmd.Flags().GetString("documents")
			documents := []string{}
			if documentsStr != "" {
				documents = strings.Split(documentsStr, ",")
			}

			mediaStr, _ := cmd.Flags().GetString("media")
			media := []string{}
			if mediaStr != "" {
				media = strings.Split(mediaStr, ",")
			}

			challenges, _ := cmd.Flags().GetString("challenges")

			msg := &types.MsgSubmitImpactReport{
				AllocationId:         allocationID,
				Submitter:            clientCtx.GetFromAddress().String(),
				Period:               period,
				BeneficiariesReached: int32(beneficiariesReached),
				FundsUtilized:        fundsUtilized,
				Metrics:              metrics,
				Documents:            documents,
				Media:                media,
				Challenges:           challenges,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().String("metrics", "", "Semicolon-separated metrics (name,target,achieved,percentage)")
	cmd.Flags().String("documents", "", "Comma-separated list of document URLs")
	cmd.Flags().String("media", "", "Comma-separated list of media URLs")
	cmd.Flags().String("challenges", "", "Description of challenges faced")

	return cmd
}

// CmdVerifyImpactReport returns the verify impact report command
func CmdVerifyImpactReport() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify-impact-report [report-id] [verified] [notes]",
		Short: "Verify an impact report",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			reportID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			verified, err := strconv.ParseBool(args[1])
			if err != nil {
				return err
			}
			notes := args[2]

			siteVisit, _ := cmd.Flags().GetBool("site-visit")
			financialAudit, _ := cmd.Flags().GetBool("financial-audit")

			msg := &types.MsgVerifyImpactReport{
				ReportId:                reportID,
				Verifier:                clientCtx.GetFromAddress().String(),
				Verified:                verified,
				Notes:                   notes,
				SiteVisitConducted:      siteVisit,
				FinancialAuditConducted: financialAudit,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().Bool("site-visit", false, "Whether a site visit was conducted")
	cmd.Flags().Bool("financial-audit", false, "Whether a financial audit was conducted")
	return cmd
}

// CmdReportFraud returns the report fraud command
func CmdReportFraud() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "report-fraud [allocation-id] [alert-type] [severity] [description]",
		Short: "Report suspected fraud or misuse of funds",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			allocationID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			alertType := args[1]
			severity := args[2]
			description := args[3]

			evidenceStr, _ := cmd.Flags().GetString("evidence")
			evidence := []string{}
			if evidenceStr != "" {
				evidence = strings.Split(evidenceStr, ",")
			}

			msg := &types.MsgReportFraud{
				AllocationId: allocationID,
				Reporter:     clientCtx.GetFromAddress().String(),
				AlertType:    alertType,
				Severity:     severity,
				Description:  description,
				Evidence:     evidence,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().String("evidence", "", "Comma-separated list of evidence URLs")
	return cmd
}

// CmdInvestigateFraud returns the investigate fraud command
func CmdInvestigateFraud() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "investigate-fraud [alert-id] [findings] [recommendation] [investigation-complete]",
		Short: "Submit investigation findings for a fraud alert",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			alertID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			findings := args[1]
			recommendation := args[2]
			complete, err := strconv.ParseBool(args[3])
			if err != nil {
				return err
			}

			report, _ := cmd.Flags().GetString("report")

			msg := &types.MsgInvestigateFraud{
				AlertId:               alertID,
				Investigator:          clientCtx.GetFromAddress().String(),
				Findings:              findings,
				Recommendation:        recommendation,
				Report:                report,
				InvestigationComplete: complete,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().String("report", "", "URL to detailed investigation report")
	return cmd
}