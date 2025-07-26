package cli

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/namo/x/identity/types"
)

// Query commands for identity backup and recovery system

// CmdQueryIdentityBackup queries an identity backup
func CmdQueryIdentityBackup() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup [backup-id]",
		Short: "Query an identity backup",
		Long: `Query details of an identity backup by its ID.

Examples:
deshchaind query identity backup backup_abc123`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.IdentityBackup(context.Background(), &types.QueryIdentityBackupRequest{
				BackupId: args[0],
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryRecoveryRequest queries a recovery request
func CmdQueryRecoveryRequest() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "recovery-request [request-id]",
		Short: "Query a recovery request",
		Long: `Query details of a recovery request by its ID.

Examples:
deshchaind query identity recovery-request recovery_xyz123`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.RecoveryRequest(context.Background(), &types.QueryRecoveryRequestRequest{
				RequestId: args[0],
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryBackupsByHolder queries backups by holder DID
func CmdQueryBackupsByHolder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backups-by-holder [holder-did]",
		Short: "Query all backups for a holder",
		Long: `Query all identity backups for a specific holder DID.

Examples:
deshchaind query identity backups-by-holder did:desh:user123`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.BackupsByHolder(context.Background(), &types.QueryBackupsByHolderRequest{
				HolderDid:  args[0],
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "backups")
	return cmd
}

// CmdQueryRecoveryRequestsByHolder queries recovery requests by holder DID
func CmdQueryRecoveryRequestsByHolder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "recovery-requests-by-holder [holder-did]",
		Short: "Query all recovery requests for a holder",
		Long: `Query all recovery requests for a specific holder DID.

Examples:
deshchaind query identity recovery-requests-by-holder did:desh:user123`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.RecoveryRequestsByHolder(context.Background(), &types.QueryRecoveryRequestsByHolderRequest{
				HolderDid:  args[0],
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "recovery-requests")
	return cmd
}

// CmdQuerySocialRecoveryGuardian queries a social recovery guardian
func CmdQuerySocialRecoveryGuardian() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "guardian [guardian-id]",
		Short: "Query a social recovery guardian",
		Long: `Query details of a social recovery guardian by its ID.

Examples:
deshchaind query identity guardian guardian_def456`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.SocialRecoveryGuardian(context.Background(), &types.QuerySocialRecoveryGuardianRequest{
				GuardianId: args[0],
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryGuardiansByHolder queries guardians by holder DID
func CmdQueryGuardiansByHolder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "guardians-by-holder [holder-did]",
		Short: "Query all guardians for a holder",
		Long: `Query all social recovery guardians for a specific holder DID.

Examples:
deshchaind query identity guardians-by-holder did:desh:user123`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.GuardiansByHolder(context.Background(), &types.QueryGuardiansByHolderRequest{
				HolderDid:  args[0],
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "guardians")
	return cmd
}

// CmdQueryBackupVerificationResult queries a backup verification result
func CmdQueryBackupVerificationResult() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup-verification [verification-id]",
		Short: "Query a backup verification result",
		Long: `Query the result of a backup integrity verification by its ID.

Examples:
deshchaind query identity backup-verification verify_ghi789`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.BackupVerificationResult(context.Background(), &types.QueryBackupVerificationResultRequest{
				VerificationId: args[0],
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryRecoveryStats queries recovery system statistics
func CmdQueryRecoveryStats() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "recovery-stats",
		Short: "Query recovery system statistics",
		Long: `Query overall statistics for the identity backup and recovery system.

Examples:
deshchaind query identity recovery-stats`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.RecoveryStats(context.Background(), &types.QueryRecoveryStatsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryGuardianVotes queries guardian votes for a recovery request
func CmdQueryGuardianVotes() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "guardian-votes [request-id]",
		Short: "Query guardian votes for a recovery request",
		Long: `Query all guardian votes for a specific recovery request.

Examples:
deshchaind query identity guardian-votes recovery_xyz123`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.GuardianVotes(context.Background(), &types.QueryGuardianVotesRequest{
				RequestId:  args[0],
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "votes")
	return cmd
}

// CmdQueryDisasterRecoveryConfig queries disaster recovery configuration
func CmdQueryDisasterRecoveryConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disaster-recovery-config [holder-did]",
		Short: "Query disaster recovery configuration",
		Long: `Query the disaster recovery configuration for a specific holder DID.

Examples:
deshchaind query identity disaster-recovery-config did:desh:user123`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.DisasterRecoveryConfig(context.Background(), &types.QueryDisasterRecoveryConfigRequest{
				HolderDid: args[0],
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}