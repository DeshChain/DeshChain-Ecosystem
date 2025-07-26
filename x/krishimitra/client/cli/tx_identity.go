package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/DeshChain/DeshChain-Ecosystem/x/krishimitra/types"
)

// GetIdentityTxCmd returns the identity transaction commands for KrishiMitra
func GetIdentityTxCmd() *cobra.Command {
	identityTxCmd := &cobra.Command{
		Use:                        "identity",
		Short:                      "Identity-related transactions for KrishiMitra",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	identityTxCmd.AddCommand(
		CmdMigrateFarmersToIdentity(),
		CmdCreateFarmerCredential(),
		CmdCreateLandRecordCredential(),
	)

	return identityTxCmd
}

// CmdMigrateFarmersToIdentity migrates existing farmers to identity system
func CmdMigrateFarmersToIdentity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate-farmers",
		Short: "Migrate existing farmers to identity system",
		Long: `Migrate all existing farmers to the new identity system.
This command creates identity records and credentials for all existing farmers.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgMigrateFarmersToIdentity(clientCtx.GetFromAddress().String())
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdCreateFarmerCredential creates a credential for a farmer
func CmdCreateFarmerCredential() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-farmer-credential [farmer-address] [pincode]",
		Short: "Create farmer credential in the identity system",
		Long: `Create a farmer profile credential in the identity system.
This command creates verifiable credentials for farmer profiles.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			farmerAddress := args[0]
			pincode := args[1]

			msg := types.NewMsgCreateFarmerCredential(
				clientCtx.GetFromAddress().String(),
				farmerAddress,
				pincode,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdCreateLandRecordCredential creates a land record credential
func CmdCreateLandRecordCredential() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-land-credential [farmer-address] [khata-number] [total-area]",
		Short: "Create land record credential for a farmer",
		Long: `Create a land ownership credential for a farmer in the identity system.
This command creates verifiable credentials for land records.`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			farmerAddress := args[0]
			khataNumber := args[1]
			totalArea := args[2]

			msg := types.NewMsgCreateLandRecordCredential(
				clientCtx.GetFromAddress().String(),
				farmerAddress,
				khataNumber,
				totalArea,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}