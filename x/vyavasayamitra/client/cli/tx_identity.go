package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/deshchain/x/vyavasayamitra/types"
)

// GetIdentityTxCmd returns the identity transaction commands for VyavasayaMitra
func GetIdentityTxCmd() *cobra.Command {
	identityTxCmd := &cobra.Command{
		Use:                        "identity",
		Short:                      "Identity-related transactions for VyavasayaMitra",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	identityTxCmd.AddCommand(
		CmdMigrateBusinessesToIdentity(),
		CmdCreateBusinessCredential(),
	)

	return identityTxCmd
}

// CmdMigrateBusinessesToIdentity migrates existing businesses to identity system
func CmdMigrateBusinessesToIdentity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate-businesses",
		Short: "Migrate existing businesses to identity system",
		Long: `Migrate all existing businesses to the new identity system.
This command creates identity records and credentials for all existing businesses.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgMigrateBusinessesToIdentity(clientCtx.GetFromAddress().String())
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdCreateBusinessCredential creates a credential for a business
func CmdCreateBusinessCredential() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-business-credential [business-id]",
		Short: "Create business credential for a business profile",
		Long: `Create a business profile credential for a business in the identity system.
This command creates verifiable credentials for business profiles.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			businessID := args[0]

			msg := types.NewMsgCreateBusinessCredential(
				clientCtx.GetFromAddress().String(),
				businessID,
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