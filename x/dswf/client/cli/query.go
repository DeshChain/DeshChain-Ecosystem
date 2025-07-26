package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/deshchain/deshchain/x/dswf/types"
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
		CmdQueryFundStatus(),
		CmdQueryAllocation(),
		CmdQueryAllocations(),
		CmdQueryPortfolio(),
		CmdQueryMonthlyReports(),
		CmdQueryGovernance(),
		CmdQueryAllocationsByCategory(),
		CmdQueryPendingDisbursements(),
		CmdQueryParams(),
	)

	return cmd
}

// CmdQueryFundStatus returns the fund status query command
func CmdQueryFundStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fund-status",
		Short: "Query the current status of the DSWF",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.FundStatus(cmd.Context(), &types.QueryFundStatusRequest{})
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

			res, err := queryClient.Allocations(cmd.Context(), &types.QueryAllocationsRequest{
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
	flags.AddPaginationFlagsToCmd(cmd, "allocations")
	cmd.Flags().String("status", "", "Filter by allocation status (proposed, approved, active, completed, rejected)")
	return cmd
}

// CmdQueryPortfolio returns the portfolio query command
func CmdQueryPortfolio() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "portfolio",
		Short: "Query the current investment portfolio",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Portfolio(cmd.Context(), &types.QueryPortfolioRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryMonthlyReports returns the monthly reports query command
func CmdQueryMonthlyReports() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "monthly-reports",
		Short: "Query monthly reports with optional period filters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			fromPeriod, _ := cmd.Flags().GetString("from")
			toPeriod, _ := cmd.Flags().GetString("to")

			res, err := queryClient.MonthlyReports(cmd.Context(), &types.QueryMonthlyReportsRequest{
				FromPeriod: fromPeriod,
				ToPeriod:   toPeriod,
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "monthly-reports")
	cmd.Flags().String("from", "", "From period (YYYY-MM)")
	cmd.Flags().String("to", "", "To period (YYYY-MM)")
	return cmd
}

// CmdQueryGovernance returns the governance query command
func CmdQueryGovernance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "governance",
		Short: "Query the fund governance configuration",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Governance(cmd.Context(), &types.QueryGovernanceRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryAllocationsByCategory returns the allocations by category query command
func CmdQueryAllocationsByCategory() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "allocations-by-category [category]",
		Short: "Query allocations filtered by category",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			status, _ := cmd.Flags().GetString("status")

			res, err := queryClient.AllocationsByCategory(cmd.Context(), &types.QueryAllocationsByCategoryRequest{
				Category:   args[0],
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
	flags.AddPaginationFlagsToCmd(cmd, "allocations-by-category")
	cmd.Flags().String("status", "", "Filter by allocation status")
	return cmd
}

// CmdQueryPendingDisbursements returns the pending disbursements query command
func CmdQueryPendingDisbursements() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pending-disbursements",
		Short: "Query pending disbursements for the next 30 days",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.PendingDisbursements(cmd.Context(), &types.QueryPendingDisbursementsRequest{
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "pending-disbursements")
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