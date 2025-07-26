package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/deshchain/x/kisaanmitra/types"
)

// GetQueryIdentityCmd returns the identity query commands for KisaanMitra module
func GetQueryIdentityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "identity",
		Short:                      "Querying commands for KisaanMitra identity integration",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdQueryFarmerIdentity(),
		CmdQueryFarmerCredentials(),
		CmdQueryLandOwnership(),
	)

	return cmd
}

// CmdQueryFarmerIdentity implements the farmer identity query command
func CmdQueryFarmerIdentity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "farmer [farmer-address]",
		Short: "Query farmer identity information",
		Long: `Query comprehensive identity information for a farmer including:
- DID and identity status
- KYC verification status and level
- Land records and total area
- Registered crops and insurance status
- Active loans and village information

Example:
$ deshchaind query kisaanmitra identity farmer cosmos1abc...`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			farmerAddress := args[0]

			// Validate farmer address
			_, err = sdk.AccAddressFromBech32(farmerAddress)
			if err != nil {
				return fmt.Errorf("invalid farmer address: %w", err)
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryFarmerIdentityRequest{
				FarmerAddress: farmerAddress,
			}

			res, err := queryClient.FarmerIdentity(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryFarmerCredentials implements the farmer credentials query command
func CmdQueryFarmerCredentials() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "credentials [farmer-address]",
		Short: "Query farmer verifiable credentials",
		Long: `Query all verifiable credentials associated with a farmer including:
- Farmer credentials with crop information
- Land record credentials
- Agricultural loan credentials
- Insurance credentials

Example:
$ deshchaind query kisaanmitra identity credentials cosmos1abc...`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			farmerAddress := args[0]

			// Validate farmer address
			_, err = sdk.AccAddressFromBech32(farmerAddress)
			if err != nil {
				return fmt.Errorf("invalid farmer address: %w", err)
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryFarmerCredentialsRequest{
				FarmerAddress: farmerAddress,
			}

			res, err := queryClient.FarmerCredentials(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryLandOwnership implements the land ownership verification query command
func CmdQueryLandOwnership() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "land-ownership [farmer-address] [required-area]",
		Short: "Verify farmer's land ownership and sufficiency",
		Long: `Verify if a farmer has sufficient land ownership for loan eligibility.
This command checks total land area from verified credentials against required area.

Example:
$ deshchaind query kisaanmitra identity land-ownership cosmos1abc... 2.5`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			farmerAddress := args[0]
			requiredArea := args[1]

			// Validate farmer address
			_, err = sdk.AccAddressFromBech32(farmerAddress)
			if err != nil {
				return fmt.Errorf("invalid farmer address: %w", err)
			}

			// Validate required area
			_, err = sdk.NewDecFromStr(requiredArea)
			if err != nil {
				return fmt.Errorf("invalid required area: %w", err)
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryLandOwnershipRequest{
				FarmerAddress: farmerAddress,
				RequiredArea:  requiredArea,
			}

			res, err := queryClient.LandOwnership(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}