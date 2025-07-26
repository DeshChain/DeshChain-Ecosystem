package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"github.com/DeshChain/DeshChain-Ecosystem/x/vyavasayamitra/types"
)

// GetIdentityQueryCmd returns the identity query commands for VyavasayaMitra
func GetIdentityQueryCmd() *cobra.Command {
	identityQueryCmd := &cobra.Command{
		Use:                        "identity",
		Short:                      "Query identity-related information for VyavasayaMitra",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	identityQueryCmd.AddCommand(
		CmdQueryBusinessIdentity(),
		CmdQueryBusinessCredentials(),
		CmdQueryBusinessCompliance(),
	)

	return identityQueryCmd
}

// CmdQueryBusinessIdentity queries a business's identity status
func CmdQueryBusinessIdentity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "business-identity [business-address]",
		Short: "Query business identity and verification status",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryBusinessIdentityRequest{
				BusinessAddress: args[0],
			}

			res, err := queryClient.BusinessIdentity(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryBusinessCredentials queries business credentials
func CmdQueryBusinessCredentials() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "business-credentials [business-address]",
		Short: "Query all business credentials for a business",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryBusinessCredentialsRequest{
				BusinessAddress: args[0],
			}

			res, err := queryClient.BusinessCredentials(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryBusinessCompliance queries business compliance status
func CmdQueryBusinessCompliance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "business-compliance [business-address] [business-type]",
		Short: "Query business compliance requirements and status",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryBusinessComplianceRequest{
				BusinessAddress: args[0],
				BusinessType:    args[1],
			}

			res, err := queryClient.BusinessCompliance(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}