package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/deshchain/deshchain/x/shikshamitra/types"
)

// GetIdentityTxCmd returns the identity transaction commands for ShikshaMitra
func GetIdentityTxCmd() *cobra.Command {
	identityTxCmd := &cobra.Command{
		Use:                        "identity",
		Short:                      "Identity-related transactions for ShikshaMitra",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	identityTxCmd.AddCommand(
		CmdMigrateStudentsToIdentity(),
		CmdCreateStudentCredential(),
	)

	return identityTxCmd
}

// CmdMigrateStudentsToIdentity migrates existing students to identity system
func CmdMigrateStudentsToIdentity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate-students",
		Short: "Migrate existing students to identity system",
		Long: `Migrate all existing students to the new identity system.
This command creates identity records and credentials for all existing students.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgMigrateStudentsToIdentity(clientCtx.GetFromAddress().String())
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdCreateStudentCredential creates a credential for a student
func CmdCreateStudentCredential() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-student-credential [student-id]",
		Short: "Create education credential for a student",
		Long: `Create an education profile credential for a student in the identity system.
This command creates verifiable credentials for student profiles.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			studentID := args[0]

			msg := types.NewMsgCreateStudentCredential(
				clientCtx.GetFromAddress().String(),
				studentID,
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