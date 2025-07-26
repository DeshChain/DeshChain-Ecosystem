package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"github.com/DeshChain/DeshChain-Ecosystem/x/shikshamitra/types"
)

// GetIdentityQueryCmd returns the identity query commands for ShikshaMitra
func GetIdentityQueryCmd() *cobra.Command {
	identityQueryCmd := &cobra.Command{
		Use:                        "identity",
		Short:                      "Query identity-related information for ShikshaMitra",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	identityQueryCmd.AddCommand(
		CmdQueryStudentIdentity(),
		CmdQueryCoApplicantIdentity(),
		CmdQueryEducationCredentials(),
	)

	return identityQueryCmd
}

// CmdQueryStudentIdentity queries a student's identity status
func CmdQueryStudentIdentity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "student-identity [student-address]",
		Short: "Query student identity and verification status",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryStudentIdentityRequest{
				StudentAddress: args[0],
			}

			res, err := queryClient.StudentIdentity(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryCoApplicantIdentity queries a co-applicant's identity status
func CmdQueryCoApplicantIdentity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "co-applicant-identity [co-applicant-address] [student-address]",
		Short: "Query co-applicant identity and verification status",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryCoApplicantIdentityRequest{
				CoApplicantAddress: args[0],
				StudentAddress:     args[1],
			}

			res, err := queryClient.CoApplicantIdentity(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryEducationCredentials queries education credentials for a student
func CmdQueryEducationCredentials() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "education-credentials [student-address]",
		Short: "Query all education credentials for a student",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryEducationCredentialsRequest{
				StudentAddress: args[0],
			}

			res, err := queryClient.EducationCredentials(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}