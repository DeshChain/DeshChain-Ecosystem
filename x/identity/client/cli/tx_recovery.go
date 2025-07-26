package cli

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/DeshChain/DeshChain-Ecosystem/x/identity/types"
)

// Transaction commands for identity backup and recovery system

// CmdCreateIdentityBackup creates an identity backup command
func CmdCreateIdentityBackup() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-backup [holder-did] [encryption-key] [retention-days]",
		Short: "Create a complete backup of an identity",
		Long: `Create a complete backup of an identity with specified recovery methods.
		
Examples:
deshchaind tx identity create-backup did:desh:user123 <base64-encryption-key> 365 \
  --recovery-methods "mnemonic,biometric,social" \
  --from mykey`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			holderDID := args[0]
			encryptionKey := args[1]
			retentionDays, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid retention days: %w", err)
			}

			// Validate encryption key is base64
			if _, err := base64.StdEncoding.DecodeString(encryptionKey); err != nil {
				return fmt.Errorf("encryption key must be base64 encoded: %w", err)
			}

			// Parse recovery methods
			recoveryMethodsStr, _ := cmd.Flags().GetString("recovery-methods")
			var recoveryMethods []types.RecoveryMethod
			if recoveryMethodsStr != "" {
				methodTypes := strings.Split(recoveryMethodsStr, ",")
				for i, methodType := range methodTypes {
					method := types.RecoveryMethod{
						MethodID:           fmt.Sprintf("method_%d", i+1),
						MethodType:         parseRecoveryMethodType(strings.TrimSpace(methodType)),
						MethodName:         strings.TrimSpace(methodType),
						Configuration:      make(map[string]interface{}),
						TrustLevel:         "medium",
						RequiredConfidence: 80,
						Enabled:            true,
						CreatedAt:          time.Now(),
						UsageCount:         0,
						Metadata:           make(map[string]interface{}),
					}
					recoveryMethods = append(recoveryMethods, method)
				}
			}

			msg := &types.MsgCreateIdentityBackup{
				Authority:       clientCtx.GetFromAddress().String(),
				HolderDID:       holderDID,
				RecoveryMethods: recoveryMethods,
				EncryptionKey:   encryptionKey,
				RetentionDays:   retentionDays,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("recovery-methods", "", "Comma-separated list of recovery methods (mnemonic,biometric,social,etc.)")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// CmdInitiateRecovery initiates identity recovery
func CmdInitiateRecovery() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "initiate-recovery [holder-did] [backup-id] [reason]",
		Short: "Initiate an identity recovery process",
		Long: `Initiate an identity recovery process using a backup.

Examples:
deshchaind tx identity initiate-recovery did:desh:user123 backup_abc123 "Lost device" \
  --from mykey`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgInitiateRecovery{
				Authority: clientCtx.GetFromAddress().String(),
				HolderDID: args[0],
				BackupID:  args[1],
				Reason:    args[2],
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdSubmitRecoveryProof submits recovery proof
func CmdSubmitRecoveryProof() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-recovery-proof [request-id] [method-id] [proof-data]",
		Short: "Submit proof for a recovery method",
		Long: `Submit proof for a recovery method during identity recovery.

Examples:
deshchaind tx identity submit-recovery-proof recovery_xyz123 method_1 <base64-proof-data> \
  --verification-data '{"key":"value"}' \
  --from mykey`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			requestID := args[0]
			methodID := args[1]
			proofData := args[2]

			// Validate proof data is base64
			if _, err := base64.StdEncoding.DecodeString(proofData); err != nil {
				return fmt.Errorf("proof data must be base64 encoded: %w", err)
			}

			// Parse verification data if provided
			verificationData := make(map[string]interface{})
			verificationDataStr, _ := cmd.Flags().GetString("verification-data")
			if verificationDataStr != "" {
				// Simple key=value parsing for CLI
				pairs := strings.Split(verificationDataStr, ",")
				for _, pair := range pairs {
					if kv := strings.Split(pair, "="); len(kv) == 2 {
						verificationData[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
					}
				}
			}

			msg := &types.MsgSubmitRecoveryProof{
				Authority:        clientCtx.GetFromAddress().String(),
				RequestID:        requestID,
				MethodID:         methodID,
				ProofData:        proofData,
				VerificationData: verificationData,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("verification-data", "", "Verification data as key=value pairs separated by commas")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// CmdExecuteRecovery executes recovery
func CmdExecuteRecovery() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "execute-recovery [request-id] [new-controller-address] [decryption-key]",
		Short: "Execute the recovery process",
		Long: `Execute the recovery process and restore identity with new controller.

Examples:
deshchaind tx identity execute-recovery recovery_xyz123 desh1abc... <base64-decryption-key> \
  --from mykey`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			requestID := args[0]
			newControllerAddress := args[1]
			decryptionKey := args[2]

			// Validate decryption key is base64
			if _, err := base64.StdEncoding.DecodeString(decryptionKey); err != nil {
				return fmt.Errorf("decryption key must be base64 encoded: %w", err)
			}

			msg := &types.MsgExecuteRecovery{
				Authority:            clientCtx.GetFromAddress().String(),
				RequestID:            requestID,
				NewControllerAddress: newControllerAddress,
				DecryptionKey:        decryptionKey,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdAddSocialRecoveryGuardian adds a social recovery guardian
func CmdAddSocialRecoveryGuardian() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-guardian [holder-did] [guardian-did] [guardian-address] [guardian-name] [weight]",
		Short: "Add a guardian for social recovery",
		Long: `Add a trusted guardian for social recovery.

Examples:
deshchaind tx identity add-guardian did:desh:user123 did:desh:guardian456 desh1guardian... "John Doe" 10 \
  --contact-info "john@example.com" \
  --public-key "guardian-public-key" \
  --from mykey`,
		Args: cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			holderDID := args[0]
			guardianDID := args[1]
			guardianAddress := args[2]
			guardianName := args[3]
			weight, err := strconv.Atoi(args[4])
			if err != nil {
				return fmt.Errorf("invalid weight: %w", err)
			}

			contactInfo, _ := cmd.Flags().GetString("contact-info")
			publicKey, _ := cmd.Flags().GetString("public-key")

			msg := &types.MsgAddSocialRecoveryGuardian{
				Authority:       clientCtx.GetFromAddress().String(),
				HolderDID:       holderDID,
				GuardianDID:     guardianDID,
				GuardianAddress: guardianAddress,
				GuardianName:    guardianName,
				Weight:          weight,
				ContactInfo:     contactInfo,
				PublicKey:       publicKey,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("contact-info", "", "Guardian contact information")
	cmd.Flags().String("public-key", "", "Guardian public key")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// CmdSubmitGuardianVote submits a guardian vote
func CmdSubmitGuardianVote() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-guardian-vote [request-id] [vote] [signature]",
		Short: "Submit a guardian vote for recovery",
		Long: `Submit a guardian vote for a recovery request.

Examples:
deshchaind tx identity submit-guardian-vote recovery_xyz123 approve <signature> \
  --reason "Guardian approval for legitimate recovery" \
  --from guardian-key`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			requestID := args[0]
			voteStr := args[1]
			signature := args[2]

			var vote types.VoteType
			switch strings.ToLower(voteStr) {
			case "approve":
				vote = types.VoteType_APPROVE
			case "reject":
				vote = types.VoteType_REJECT
			case "abstain":
				vote = types.VoteType_ABSTAIN
			default:
				return fmt.Errorf("invalid vote type: %s (must be approve, reject, or abstain)", voteStr)
			}

			reason, _ := cmd.Flags().GetString("reason")

			msg := &types.MsgSubmitGuardianVote{
				Authority: clientCtx.GetFromAddress().String(),
				RequestID: requestID,
				Vote:      vote,
				Reason:    reason,
				Signature: signature,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("reason", "", "Reason for the vote")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// CmdVerifyBackupIntegrity verifies backup integrity
func CmdVerifyBackupIntegrity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify-backup [backup-id]",
		Short: "Verify the integrity of a backup",
		Long: `Verify the integrity and recoverability of an identity backup.

Examples:
deshchaind tx identity verify-backup backup_abc123 \
  --verification-key <base64-key> \
  --from mykey`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			backupID := args[0]
			verificationKey, _ := cmd.Flags().GetString("verification-key")

			// Validate verification key if provided
			if verificationKey != "" {
				if _, err := base64.StdEncoding.DecodeString(verificationKey); err != nil {
					return fmt.Errorf("verification key must be base64 encoded: %w", err)
				}
			}

			msg := &types.MsgVerifyBackupIntegrity{
				Authority:       clientCtx.GetFromAddress().String(),
				BackupID:        backupID,
				VerificationKey: verificationKey,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("verification-key", "", "Optional verification key (base64 encoded)")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// Helper functions

// parseRecoveryMethodType converts string to RecoveryMethodType
func parseRecoveryMethodType(methodType string) types.RecoveryMethodType {
	switch strings.ToLower(methodType) {
	case "mnemonic", "mnemonic_phrase":
		return types.RecoveryMethodType_MNEMONIC_PHRASE
	case "social", "social_recovery":
		return types.RecoveryMethodType_SOCIAL_RECOVERY
	case "guardian", "guardian_multisig":
		return types.RecoveryMethodType_GUARDIAN_MULTISIG
	case "biometric", "biometric_backup":
		return types.RecoveryMethodType_BIOMETRIC_BACKUP
	case "hardware", "hardware_key":
		return types.RecoveryMethodType_HARDWARE_KEY
	case "backup_codes":
		return types.RecoveryMethodType_BACKUP_CODES
	case "email", "email_verification":
		return types.RecoveryMethodType_EMAIL_VERIFICATION
	case "sms", "sms_verification":
		return types.RecoveryMethodType_SMS_VERIFICATION
	case "identity_provider":
		return types.RecoveryMethodType_IDENTITY_PROVIDER
	case "institutional", "institutional_recovery":
		return types.RecoveryMethodType_INSTITUTIONAL_RECOVERY
	case "zkp", "zero_knowledge_proof":
		return types.RecoveryMethodType_ZERO_KNOWLEDGE_PROOF
	default:
		return types.RecoveryMethodType_MNEMONIC_PHRASE // Default fallback
	}
}