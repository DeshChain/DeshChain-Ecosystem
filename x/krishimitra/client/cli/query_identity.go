package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"github.com/DeshChain/DeshChain-Ecosystem/x/krishimitra/types"
)

// GetIdentityQueryCmd returns the identity query commands for KrishiMitra
func GetIdentityQueryCmd() *cobra.Command {
	identityQueryCmd := &cobra.Command{
		Use:                        "identity",
		Short:                      "Query identity-related information for KrishiMitra",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	identityQueryCmd.AddCommand(
		CmdQueryFarmerIdentity(),
		CmdQueryFarmerCredentials(),
		CmdQueryLandOwnership(),
	)

	return identityQueryCmd
}

// CmdQueryFarmerIdentity queries a farmer's identity status
func CmdQueryFarmerIdentity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "farmer-identity [farmer-address]",
		Short: "Query farmer identity and verification status",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryFarmerIdentityRequest{
				FarmerAddress: args[0],
			}

			res, err := queryClient.FarmerIdentity(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryFarmerCredentials queries farmer credentials
func CmdQueryFarmerCredentials() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "farmer-credentials [farmer-address]",
		Short: "Query all farmer credentials",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryFarmerCredentialsRequest{
				FarmerAddress: args[0],
			}

			res, err := queryClient.FarmerCredentials(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryLandOwnership queries land ownership verification
func CmdQueryLandOwnership() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "land-ownership [farmer-address] [required-area]",
		Short: "Query land ownership verification for a farmer",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryLandOwnershipRequest{
				FarmerAddress: args[0],
				RequiredArea:  args[1],
			}

			res, err := queryClient.LandOwnership(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}