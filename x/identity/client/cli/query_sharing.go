package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/deshchain/deshchain/x/identity/types"
)

// GetQuerySharingCmd returns the query commands for identity sharing
func GetQuerySharingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "sharing",
		Short:                      "Querying commands for identity sharing",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdQueryShareRequest(),
		CmdQueryShareResponse(),
		CmdQuerySharingAgreement(),
		CmdQueryAccessPolicy(),
		CmdQueryShareAuditLogs(),
		CmdQueryModuleCapabilities(),
	)

	return cmd
}

// CmdQueryShareRequest implements the query share request command
func CmdQueryShareRequest() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "request [request-id]",
		Short: "Query a specific identity sharing request",
		Long: `Query details of a specific identity sharing request by its ID.

Example:
$ deshchaind query identity sharing request share_req_1234567890abcdef`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			requestID := args[0]

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryShareRequestRequest{
				RequestID: requestID,
			}

			res, err := queryClient.ShareRequest(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryShareResponse implements the query share response command
func CmdQueryShareResponse() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "response [request-id]",
		Short: "Query the response for a sharing request",
		Long: `Query the response details for a specific identity sharing request.

Example:
$ deshchaind query identity sharing response share_req_1234567890abcdef`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			requestID := args[0]

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryShareResponseRequest{
				RequestID: requestID,
			}

			res, err := queryClient.ShareResponse(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQuerySharingAgreement implements the query sharing agreement command
func CmdQuerySharingAgreement() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agreement [agreement-id]",
		Short: "Query a specific sharing agreement",
		Long: `Query details of a specific sharing agreement between modules.

Example:
$ deshchaind query identity sharing agreement share_agr_1234567890abcdef`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			agreementID := args[0]

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QuerySharingAgreementRequest{
				AgreementID: agreementID,
			}

			res, err := queryClient.SharingAgreement(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryAccessPolicy implements the query access policy command
func CmdQueryAccessPolicy() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "policy [policy-id] [--holder-did holder-did]",
		Short: "Query access policy details",
		Long: `Query access policy details by policy ID or all policies for a holder.

Examples:
$ deshchaind query identity sharing policy access_pol_1234567890abcdef
$ deshchaind query identity sharing policy --holder-did did:desh:user123`,
		Args: cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			holderDID, _ := cmd.Flags().GetString("holder-did")

			if len(args) == 1 {
				// Query specific policy
				policyID := args[0]
				req := &types.QueryAccessPolicyRequest{
					PolicyID: policyID,
				}

				res, err := queryClient.AccessPolicy(context.Background(), req)
				if err != nil {
					return err
				}

				return clientCtx.PrintProto(res)
			} else if holderDID != "" {
				// Query policies by holder
				req := &types.QueryAccessPoliciesByHolderRequest{
					HolderDID: holderDID,
				}

				res, err := queryClient.AccessPoliciesByHolder(context.Background(), req)
				if err != nil {
					return err
				}

				return clientCtx.PrintProto(res)
			} else {
				return fmt.Errorf("must provide either policy-id as argument or --holder-did flag")
			}
		},
	}

	cmd.Flags().String("holder-did", "", "Query all policies for a specific holder DID")
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryShareAuditLogs implements the query share audit logs command
func CmdQueryShareAuditLogs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "audit-logs [--holder-did holder-did] [--request-id request-id] [--limit limit]",
		Short: "Query identity sharing audit logs",
		Long: `Query identity sharing audit logs with optional filters.

Examples:
$ deshchaind query identity sharing audit-logs --holder-did did:desh:user123
$ deshchaind query identity sharing audit-logs --request-id share_req_1234567890abcdef
$ deshchaind query identity sharing audit-logs --limit 50`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			holderDID, _ := cmd.Flags().GetString("holder-did")
			requestID, _ := cmd.Flags().GetString("request-id")
			limit, _ := cmd.Flags().GetUint64("limit")

			if limit == 0 {
				limit = 100 // Default limit
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryShareAuditLogsRequest{
				HolderDID: holderDID,
				RequestID: requestID,
				Limit:     limit,
			}

			res, err := queryClient.ShareAuditLogs(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().String("holder-did", "", "Filter by holder DID")
	cmd.Flags().String("request-id", "", "Filter by request ID")
	cmd.Flags().Uint64("limit", 100, "Maximum number of logs to return")
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryModuleCapabilities implements the query module capabilities command
func CmdQueryModuleCapabilities() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "capabilities [module-name]",
		Short: "Query module capabilities for identity sharing",
		Long: `Query the capabilities of a specific module for identity sharing,
including what data types it can request and provide.

Example:
$ deshchaind query identity sharing capabilities tradefinance`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			moduleName := args[0]

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryModuleCapabilitiesRequest{
				ModuleName: moduleName,
			}

			res, err := queryClient.ModuleCapabilities(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
