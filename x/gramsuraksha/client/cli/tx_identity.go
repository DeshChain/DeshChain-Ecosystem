package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/DeshChain/DeshChain-Ecosystem/x/gramsuraksha/types"
	"github.com/spf13/cobra"
)

// GetIdentityTxCmd returns identity-related transaction commands
func GetIdentityTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "identity",
		Short:                      "Identity integration commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		GetMigrateToIdentityCmd(),
		GetVerifyIdentityCmd(),
	)

	return cmd
}

// GetMigrateToIdentityCmd returns the command to migrate participants to identity system
func GetMigrateToIdentityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate existing participants to identity system",
		Long: `Migrate all existing GramSuraksha participants to use the new identity system.
This will create DID-based identities and verifiable credentials for all participants.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgMigrateToIdentity{
				Authority: clientCtx.GetFromAddress().String(),
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetVerifyIdentityCmd returns the command to verify participant identity
func GetVerifyIdentityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify [participant-address]",
		Short: "Verify participant identity status",
		Long: `Verify the identity status of a GramSuraksha participant.
Shows both traditional KYC status and DID-based identity credentials.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryParticipantIdentityRequest{
				ParticipantAddress: args[0],
			}

			res, err := queryClient.ParticipantIdentity(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// Additional commands for identity management

// GetCreateIdentityCredentialCmd creates a credential for a participant
func GetCreateIdentityCredentialCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-credential [participant-id] [scheme-id]",
		Short: "Create identity credential for a participant",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgCreateParticipantCredential{
				Authority:     clientCtx.GetFromAddress().String(),
				ParticipantId: args[0],
				SchemeId:      args[1],
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetVerifyAgeWithZKPCmd verifies age using zero-knowledge proof
func GetVerifyAgeWithZKPCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify-age-zkp [address] [min-age] [max-age]",
		Short: "Verify participant age using zero-knowledge proof",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			// Parse ages
			var minAge, maxAge int32
			if _, err := fmt.Sscanf(args[1], "%d", &minAge); err != nil {
				return fmt.Errorf("invalid min-age: %s", args[1])
			}
			if _, err := fmt.Sscanf(args[2], "%d", &maxAge); err != nil {
				return fmt.Errorf("invalid max-age: %s", args[2])
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryVerifyAgeRequest{
				Address: args[0],
				MinAge:  minAge,
				MaxAge:  maxAge,
			}

			res, err := queryClient.VerifyAge(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}