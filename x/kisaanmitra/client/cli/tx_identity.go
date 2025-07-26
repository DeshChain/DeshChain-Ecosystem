package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/kisaanmitra/types"
)

// GetTxIdentityCmd returns the identity transaction commands for KisaanMitra module
func GetTxIdentityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "identity",
		Short:                      "Identity transaction subcommands for KisaanMitra",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdMigrateFarmersToIdentity(),
		CmdCreateFarmerCredential(),
		CmdCreateLandRecordCredential(),
	)

	return cmd
}

// CmdMigrateFarmersToIdentity implements the migrate farmers to identity command
func CmdMigrateFarmersToIdentity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate-farmers",
		Short: "Migrate all farmers to identity system",
		Long: `Migrate all existing farmers in the KisaanMitra module to the new identity system.
This command creates identities and basic credentials for all registered farmers.

Example:
$ deshchaind tx kisaanmitra identity migrate-farmers --from mykey`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			authority := clientCtx.GetFromAddress()

			msg := types.NewMsgMigrateFarmersToIdentity(authority.String())
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdCreateFarmerCredential implements the create farmer credential command
func CmdCreateFarmerCredential() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-farmer-credential [farmer-address] [pincode]",
		Short: "Create a farmer credential for a borrower",
		Long: `Create a verifiable credential for a farmer with their basic information.
This includes farmer details, land information, and registered crops.

Example:
$ deshchaind tx kisaanmitra identity create-farmer-credential cosmos1abc... 560001 --from mykey`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			authority := clientCtx.GetFromAddress()
			farmerAddress := args[0]
			pincode := args[1]

			// Validate farmer address
			_, err = sdk.AccAddressFromBech32(farmerAddress)
			if err != nil {
				return fmt.Errorf("invalid farmer address: %w", err)
			}

			msg := types.NewMsgCreateFarmerCredential(authority.String(), farmerAddress, pincode)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdCreateLandRecordCredential implements the create land record credential command
func CmdCreateLandRecordCredential() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-land-record [farmer-address] [khata-number] [total-area]",
		Short: "Create a land record credential for a farmer",
		Long: `Create a verifiable credential for a farmer's land records.
This includes khata number, total area, and land ownership details.

Example:
$ deshchaind tx kisaanmitra identity create-land-record cosmos1abc... KH123456 2.5 --from mykey`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			authority := clientCtx.GetFromAddress()
			farmerAddress := args[0]
			khataNumber := args[1]
			totalArea := args[2]

			// Validate farmer address
			_, err = sdk.AccAddressFromBech32(farmerAddress)
			if err != nil {
				return fmt.Errorf("invalid farmer address: %w", err)
			}

			// Validate total area
			_, err = sdk.NewDecFromStr(totalArea)
			if err != nil {
				return fmt.Errorf("invalid total area: %w", err)
			}

			msg := types.NewMsgCreateLandRecordCredential(authority.String(), farmerAddress, khataNumber, totalArea)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}