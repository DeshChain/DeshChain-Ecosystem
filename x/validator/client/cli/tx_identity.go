package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/namo/x/validator/types"
)

// GetTxIdentityCmd returns the identity transaction commands for Validator module
func GetTxIdentityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "identity",
		Short:                      "Identity transaction subcommands for Validator",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdMigrateValidatorsToIdentity(),
		CmdCreateValidatorCredential(),
		CmdCreateNFTBindingCredential(),
		CmdCreateReferralCredential(),
		CmdCreateTokenLaunchCredential(),
		CmdCreateComplianceCredential(),
	)

	return cmd
}

// CmdMigrateValidatorsToIdentity implements the migrate validators to identity command
func CmdMigrateValidatorsToIdentity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate-validators",
		Short: "Migrate all validators to identity system",
		Long: `Migrate all existing validators in the Validator module to the new identity system.
This command creates identities and basic credentials for all active validators.

Example:
$ deshchaind tx validator identity migrate-validators --from mykey`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			authority := clientCtx.GetFromAddress()

			msg := types.NewMsgMigrateValidatorsToIdentity(authority.String())
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdCreateValidatorCredential implements the create validator credential command
func CmdCreateValidatorCredential() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-validator-credential [validator-address] [operator-address] [validator-rank] [stake-amount]",
		Short: "Create a validator credential",
		Long: `Create a verifiable credential for a validator with their staking information.
This includes validator rank, operator address, stake amount, and tier information.

Example:
$ deshchaind tx validator identity create-validator-credential cosmos1abc... deshvaloper1abc... 1 1000000000000 --from mykey`,
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			authority := clientCtx.GetFromAddress()
			validatorAddress := args[0]
			operatorAddress := args[1]
			validatorRankStr := args[2]
			stakeAmount := args[3]

			// Validate validator address
			_, err = sdk.AccAddressFromBech32(validatorAddress)
			if err != nil {
				return fmt.Errorf("invalid validator address: %w", err)
			}

			// Validate validator rank
			validatorRank, err := strconv.ParseUint(validatorRankStr, 10, 32)
			if err != nil {
				return fmt.Errorf("invalid validator rank: %w", err)
			}
			if validatorRank == 0 || validatorRank > 21 {
				return fmt.Errorf("validator rank must be between 1 and 21")
			}

			// Validate stake amount
			_, ok := sdk.NewIntFromString(stakeAmount)
			if !ok {
				return fmt.Errorf("invalid stake amount: %s", stakeAmount)
			}

			msg := types.NewMsgCreateValidatorCredential(authority.String(), validatorAddress, operatorAddress, uint32(validatorRank), stakeAmount)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdCreateNFTBindingCredential implements the create NFT binding credential command
func CmdCreateNFTBindingCredential() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-nft-binding [validator-address] [nft-token-id]",
		Short: "Create an NFT binding credential for a validator",
		Long: `Create a verifiable credential for a validator's NFT binding.
This credential proves the validator owns and has bound a genesis NFT.

Example:
$ deshchaind tx validator identity create-nft-binding cosmos1abc... 1 --from mykey`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			authority := clientCtx.GetFromAddress()
			validatorAddress := args[0]
			nftTokenIDStr := args[1]

			// Validate validator address
			_, err = sdk.AccAddressFromBech32(validatorAddress)
			if err != nil {
				return fmt.Errorf("invalid validator address: %w", err)
			}

			// Validate NFT token ID
			nftTokenID, err := strconv.ParseUint(nftTokenIDStr, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid NFT token ID: %w", err)
			}

			msg := types.NewMsgCreateNFTBindingCredential(authority.String(), validatorAddress, nftTokenID)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdCreateReferralCredential implements the create referral credential command
func CmdCreateReferralCredential() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-referral-credential [referrer-address] [referred-address] [referral-id]",
		Short: "Create a referral credential",
		Long: `Create a verifiable credential for a validator referral.
This proves that a validator successfully referred another validator.

Example:
$ deshchaind tx validator identity create-referral-credential cosmos1abc... cosmos1def... 1 --from mykey`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			authority := clientCtx.GetFromAddress()
			referrerAddress := args[0]
			referredAddress := args[1]
			referralIDStr := args[2]

			// Validate referrer address
			_, err = sdk.AccAddressFromBech32(referrerAddress)
			if err != nil {
				return fmt.Errorf("invalid referrer address: %w", err)
			}

			// Validate referred address
			_, err = sdk.AccAddressFromBech32(referredAddress)
			if err != nil {
				return fmt.Errorf("invalid referred address: %w", err)
			}

			// Validate referral ID
			referralID, err := strconv.ParseUint(referralIDStr, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid referral ID: %w", err)
			}

			msg := types.NewMsgCreateReferralCredential(authority.String(), referrerAddress, referredAddress, referralID)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdCreateTokenLaunchCredential implements the create token launch credential command
func CmdCreateTokenLaunchCredential() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-token-launch [validator-address] [token-id]",
		Short: "Create a token launch credential",
		Long: `Create a verifiable credential for a validator's token launch.
This proves that a validator has successfully launched their validator token.

Example:
$ deshchaind tx validator identity create-token-launch cosmos1abc... 1 --from mykey`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			authority := clientCtx.GetFromAddress()
			validatorAddress := args[0]
			tokenIDStr := args[1]

			// Validate validator address
			_, err = sdk.AccAddressFromBech32(validatorAddress)
			if err != nil {
				return fmt.Errorf("invalid validator address: %w", err)
			}

			// Validate token ID
			tokenID, err := strconv.ParseUint(tokenIDStr, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid token ID: %w", err)
			}

			msg := types.NewMsgCreateTokenLaunchCredential(authority.String(), validatorAddress, tokenID)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdCreateComplianceCredential implements the create compliance credential command
func CmdCreateComplianceCredential() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-compliance [validator-address] [compliance-level] [jurisdictions] [required-docs] [verified-docs]",
		Short: "Create a compliance credential",
		Long: `Create a verifiable credential for validator compliance verification.
This includes KYB, AML, sanctions screening, and document verification.

Arguments:
- jurisdictions: comma-separated list of jurisdiction codes (e.g., "US,IN,EU")
- required-docs: comma-separated list of required documents (e.g., "passport,license")
- verified-docs: comma-separated list of verified documents (e.g., "passport,license")

Example:
$ deshchaind tx validator identity create-compliance cosmos1abc... enhanced "US,IN" "passport,license" "passport,license" --from mykey`,
		Args: cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			authority := clientCtx.GetFromAddress()
			validatorAddress := args[0]
			complianceLevel := args[1]
			jurisdictionsStr := args[2]
			requiredDocsStr := args[3]
			verifiedDocsStr := args[4]

			// Validate validator address
			_, err = sdk.AccAddressFromBech32(validatorAddress)
			if err != nil {
				return fmt.Errorf("invalid validator address: %w", err)
			}

			// Parse comma-separated lists
			var jurisdictions []string
			if jurisdictionsStr != "" {
				jurisdictions = strings.Split(jurisdictionsStr, ",")
				for i, j := range jurisdictions {
					jurisdictions[i] = strings.TrimSpace(j)
				}
			}

			var requiredDocuments []string
			if requiredDocsStr != "" {
				requiredDocuments = strings.Split(requiredDocsStr, ",")
				for i, d := range requiredDocuments {
					requiredDocuments[i] = strings.TrimSpace(d)
				}
			}

			var verifiedDocuments []string
			if verifiedDocsStr != "" {
				verifiedDocuments = strings.Split(verifiedDocsStr, ",")
				for i, d := range verifiedDocuments {
					verifiedDocuments[i] = strings.TrimSpace(d)
				}
			}

			msg := types.NewMsgCreateComplianceCredential(authority.String(), validatorAddress, complianceLevel, jurisdictions, requiredDocuments, verifiedDocuments)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}