package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/namo/x/validator/types"
)

// GetQueryIdentityCmd returns the identity query commands for Validator module
func GetQueryIdentityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "identity",
		Short:                      "Querying commands for Validator identity integration",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdQueryValidatorIdentity(),
		CmdQueryValidatorCompliance(),
		CmdQueryValidatorCredentials(),
	)

	return cmd
}

// CmdQueryValidatorIdentity implements the validator identity query command
func CmdQueryValidatorIdentity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator [validator-address]",
		Short: "Query validator identity information",
		Long: `Query comprehensive identity information for a validator including:
- DID and identity status
- KYC verification status and level  
- Validator rank and stake verification
- NFT binding status
- Referral and token launch credentials
- Compliance status and jurisdiction

Example:
$ deshchaind query validator identity validator cosmos1abc...`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			validatorAddress := args[0]

			// Validate validator address
			_, err = sdk.AccAddressFromBech32(validatorAddress)
			if err != nil {
				return fmt.Errorf("invalid validator address: %w", err)
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryValidatorIdentityRequest{
				ValidatorAddress: validatorAddress,
			}

			res, err := queryClient.ValidatorIdentity(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryValidatorCompliance implements the validator compliance query command
func CmdQueryValidatorCompliance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compliance [validator-address] [required-jurisdictions]",
		Short: "Query validator compliance verification",
		Long: `Query comprehensive compliance information for a validator including:
- KYB (Know Your Business) completion status
- AML (Anti-Money Laundering) verification
- Sanctions screening results
- Document verification status
- Jurisdiction-specific compliance
- Compliance expiry dates

Arguments:
- required-jurisdictions: comma-separated list of jurisdiction codes to check (e.g., "US,IN,EU")

Example:
$ deshchaind query validator identity compliance cosmos1abc... "US,IN"`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			validatorAddress := args[0]
			jurisdictionsStr := args[1]

			// Validate validator address
			_, err = sdk.AccAddressFromBech32(validatorAddress)
			if err != nil {
				return fmt.Errorf("invalid validator address: %w", err)
			}

			// Parse jurisdictions
			var requiredJurisdictions []string
			if jurisdictionsStr != "" {
				requiredJurisdictions = strings.Split(jurisdictionsStr, ",")
				for i, j := range requiredJurisdictions {
					requiredJurisdictions[i] = strings.TrimSpace(j)
				}
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryValidatorComplianceRequest{
				ValidatorAddress:      validatorAddress,
				RequiredJurisdictions: requiredJurisdictions,
			}

			res, err := queryClient.ValidatorCompliance(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryValidatorCredentials implements the validator credentials query command
func CmdQueryValidatorCredentials() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "credentials [validator-address]",
		Short: "Query validator verifiable credentials",
		Long: `Query all verifiable credentials associated with a validator including:
- Validator credentials with rank and stake information
- NFT binding credentials
- Referral credentials for successful referrals
- Token launch credentials
- Compliance credentials with verification details
- KYC/KYB credentials

Example:
$ deshchaind query validator identity credentials cosmos1abc...`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			validatorAddress := args[0]

			// Validate validator address
			_, err = sdk.AccAddressFromBech32(validatorAddress)
			if err != nil {
				return fmt.Errorf("invalid validator address: %w", err)
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryValidatorCredentialsRequest{
				ValidatorAddress: validatorAddress,
			}

			res, err := queryClient.ValidatorCredentials(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}