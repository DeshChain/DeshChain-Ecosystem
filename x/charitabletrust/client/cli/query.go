package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/DeshChain/DeshChain-Ecosystem/x/charitabletrust/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdQueryTrustFundBalance(),
		CmdQueryAllocation(),
		CmdQueryAllocations(),
		CmdQueryAllocationProposal(),
		CmdQueryAllocationProposals(),
		CmdQueryImpactReport(),
		CmdQueryImpactReports(),
		CmdQueryFraudAlert(),
		CmdQueryFraudAlerts(),
		CmdQueryTrustGovernance(),
		CmdQueryAllocationsByOrganization(),
		CmdQueryParams(),
	)

	return cmd
}

// CmdQueryTrustFundBalance returns the trust fund balance query command
func CmdQueryTrustFundBalance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "trust-fund-balance",
		Short: "Query the current trust fund balance",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.TrustFundBalance(cmd.Context(), &types.QueryTrustFundBalanceRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryAllocation returns the allocation query command
func CmdQueryAllocation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "allocation [allocation-id]",
		Short: "Query a specific allocation by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			allocationID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			res, err := queryClient.Allocation(cmd.Context(), &types.QueryAllocationRequest{
				AllocationId: allocationID,
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

// CmdQueryAllocations returns the allocations query command
func CmdQueryAllocations() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "allocations",
		Short: "Query all allocations with optional filters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			status, _ := cmd.Flags().GetString("status")
			category, _ := cmd.Flags().GetString("category")

			res, err := queryClient.Allocations(cmd.Context(), &types.QueryAllocationsRequest{
				Status:     status,
				Category:   category,
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "allocations")
	cmd.Flags().String("status", "", "Filter by allocation status")
	cmd.Flags().String("category", "", "Filter by allocation category")
	return cmd
}

// CmdQueryAllocationProposal returns the allocation proposal query command
func CmdQueryAllocationProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposal [proposal-id]",
		Short: "Query a specific allocation proposal by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			res, err := queryClient.AllocationProposal(cmd.Context(), &types.QueryAllocationProposalRequest{
				ProposalId: proposalID,
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

// CmdQueryAllocationProposals returns the allocation proposals query command
func CmdQueryAllocationProposals() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposals",
		Short: "Query all allocation proposals with optional filters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			status, _ := cmd.Flags().GetString("status")

			res, err := queryClient.AllocationProposals(cmd.Context(), &types.QueryAllocationProposalsRequest{
				Status:     status,
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "proposals")
	cmd.Flags().String("status", "", "Filter by proposal status (pending, approved, rejected, executed)")
	return cmd
}

// CmdQueryImpactReport returns the impact report query command
func CmdQueryImpactReport() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "impact-report [report-id]",
		Short: "Query a specific impact report by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			reportID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			res, err := queryClient.ImpactReport(cmd.Context(), &types.QueryImpactReportRequest{
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

// CmdQueryImpactReports returns the impact reports query command
func CmdQueryImpactReports() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "impact-reports [allocation-id]",
		Short: "Query impact reports for an allocation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			allocationID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			res, err := queryClient.ImpactReports(cmd.Context(), &types.QueryImpactReportsRequest{
				AllocationId: allocationID,
				Pagination:   pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "impact-reports")
	return cmd
}

// CmdQueryFraudAlert returns the fraud alert query command
func CmdQueryFraudAlert() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fraud-alert [alert-id]",
		Short: "Query a specific fraud alert by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			alertID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			res, err := queryClient.FraudAlert(cmd.Context(), &types.QueryFraudAlertRequest{
				AlertId: alertID,
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

// CmdQueryFraudAlerts returns the fraud alerts query command
func CmdQueryFraudAlerts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fraud-alerts",
		Short: "Query all fraud alerts with optional filters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			status, _ := cmd.Flags().GetString("status")
			severity, _ := cmd.Flags().GetString("severity")

			res, err := queryClient.FraudAlerts(cmd.Context(), &types.QueryFraudAlertsRequest{
				Status:     status,
				Severity:   severity,
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "fraud-alerts")
	cmd.Flags().String("status", "", "Filter by alert status")
	cmd.Flags().String("severity", "", "Filter by severity (low, medium, high, critical)")
	return cmd
}

// CmdQueryTrustGovernance returns the trust governance query command
func CmdQueryTrustGovernance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "trust-governance",
		Short: "Query the trust governance configuration",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.TrustGovernance(cmd.Context(), &types.QueryTrustGovernanceRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryAllocationsByOrganization returns the allocations by organization query command
func CmdQueryAllocationsByOrganization() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "allocations-by-org [org-wallet-id]",
		Short: "Query allocations for a specific organization",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			orgWalletID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			status, _ := cmd.Flags().GetString("status")

			res, err := queryClient.AllocationsByOrganization(cmd.Context(), &types.QueryAllocationsByOrganizationRequest{
				OrgWalletId: orgWalletID,
				Status:      status,
				Pagination:  pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "allocations-by-org")
	cmd.Flags().String("status", "", "Filter by allocation status")
	return cmd
}

// CmdQueryParams returns the params query command
func CmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the module parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(cmd.Context(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}