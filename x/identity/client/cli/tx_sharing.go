package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/DeshChain/DeshChain-Ecosystem/x/identity/types"
)

// GetTxSharingCmd returns the transaction commands for identity sharing
func GetTxSharingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "sharing",
		Short:                      "Identity sharing transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdCreateShareRequest(),
		CmdApproveShareRequest(),
		CmdDenyShareRequest(),
		CmdCreateSharingAgreement(),
		CmdCreateAccessPolicy(),
	)

	return cmd
}

// CmdCreateShareRequest implements the create share request command
func CmdCreateShareRequest() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-request [requester-module] [provider-module] [holder-did] [purpose] [justification] [ttl-hours] [requested-data]",
		Short: "Create a new identity sharing request",
		Long: `Create a new identity sharing request between modules.

The requested-data parameter should be a JSON array of data requests with format:
[{"credential_type":"KYCCredential","attributes":["name","age"],"minimum_trust_level":"high","required":true}]

Example:
$ deshchaind tx identity sharing create-request tradefinance gramsuraksha did:desh:user123 "Loan eligibility verification" "Required for trade finance loan" 24 '[{"credential_type":"KYCCredential","attributes":["name","age"],"required":true}]' --from mykey`,
		Args: cobra.ExactArgs(7),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			requesterModule := args[0]
			providerModule := args[1]
			holderDID := args[2]
			purpose := args[3]
			justification := args[4]
			ttlHours, err := strconv.ParseInt(args[5], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid TTL hours: %w", err)
			}

			// Parse requested data JSON
			var requestedData []types.DataRequest
			if err := clientCtx.Codec.UnmarshalJSON([]byte(args[6]), &requestedData); err != nil {
				return fmt.Errorf("invalid requested data JSON: %w", err)
			}

			authority := clientCtx.GetFromAddress()

			msg := types.NewMsgCreateShareRequest(
				authority.String(),
				requesterModule,
				providerModule,
				holderDID,
				requestedData,
				purpose,
				justification,
				ttlHours,
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

// CmdApproveShareRequest implements the approve share request command
func CmdApproveShareRequest() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approve [request-id] [access-token]",
		Short: "Approve a pending identity sharing request",
		Long: `Approve a pending identity sharing request.

The access-token parameter is optional. If not provided, a random token will be generated.

Example:
$ deshchaind tx identity sharing approve share_req_1234567890abcdef token_abc123 --from mykey`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			requestID := args[0]
			accessToken := ""
			if len(args) > 1 {
				accessToken = args[1]
			}

			authority := clientCtx.GetFromAddress()

			msg := types.NewMsgApproveShareRequest(
				authority.String(),
				requestID,
				accessToken,
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

// CmdDenyShareRequest implements the deny share request command
func CmdDenyShareRequest() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deny [request-id] [denial-reason]",
		Short: "Deny a pending identity sharing request",
		Long: `Deny a pending identity sharing request with a reason.

Example:
$ deshchaind tx identity sharing deny share_req_1234567890abcdef "Insufficient KYC level" --from mykey`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			requestID := args[0]
			denialReason := args[1]

			authority := clientCtx.GetFromAddress()

			msg := types.NewMsgDenyShareRequest(
				authority.String(),
				requestID,
				denialReason,
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

// CmdCreateSharingAgreement implements the create sharing agreement command
func CmdCreateSharingAgreement() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-agreement [requester-module] [provider-module] [trust-level] [max-ttl-hours] [validity-days] [allowed-data-types] [purposes] [auto-approve]",
		Short: "Create a standing agreement between modules",
		Long: `Create a standing agreement between modules for automatic sharing.

The allowed-data-types and purposes parameters should be comma-separated lists.
The auto-approve parameter should be true or false.

Example:
$ deshchaind tx identity sharing create-agreement tradefinance gramsuraksha high 24 365 "KYCCredential,BiometricCredential" "loan_verification,risk_assessment" true --from mykey`,
		Args: cobra.ExactArgs(8),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			requesterModule := args[0]
			providerModule := args[1]
			trustLevel := args[2]
			maxTTLHours, err := strconv.ParseInt(args[3], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid max TTL hours: %w", err)
			}
			validityDays, err := strconv.ParseInt(args[4], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid validity days: %w", err)
			}

			allowedDataTypes := strings.Split(args[5], ",")
			for i, dataType := range allowedDataTypes {
				allowedDataTypes[i] = strings.TrimSpace(dataType)
			}

			purposes := strings.Split(args[6], ",")
			for i, purpose := range purposes {
				purposes[i] = strings.TrimSpace(purpose)
			}

			autoApprove, err := strconv.ParseBool(args[7])
			if err != nil {
				return fmt.Errorf("invalid auto-approve value: %w", err)
			}

			authority := clientCtx.GetFromAddress()

			msg := types.NewMsgCreateSharingAgreement(
				authority.String(),
				requesterModule,
				providerModule,
				allowedDataTypes,
				purposes,
				trustLevel,
				autoApprove,
				maxTTLHours,
				validityDays,
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

// CmdCreateAccessPolicy implements the create access policy command
func CmdCreateAccessPolicy() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-policy [holder-did] [max-shares-per-day] [require-explicit-consent]",
		Short: "Create an access policy for identity sharing",
		Long: `Create an access policy to control how identity data can be shared.

Optional flags:
--allowed-modules: Comma-separated list of allowed modules
--denied-modules: Comma-separated list of denied modules
--purpose-restrictions: Comma-separated list of allowed purposes
--data-restrictions: JSON map of credential types to allowed attributes

Example:
$ deshchaind tx identity sharing create-policy did:desh:user123 5 true --allowed-modules "tradefinance,gramsuraksha" --purpose-restrictions "loan_verification,kyc_check" --from mykey`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			holderDID := args[0]
			maxSharesPerDay, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid max shares per day: %w", err)
			}
			requireExplicitConsent, err := strconv.ParseBool(args[2])
			if err != nil {
				return fmt.Errorf("invalid require explicit consent value: %w", err)
			}

			// Parse optional flags
			allowedModulesStr, _ := cmd.Flags().GetString("allowed-modules")
			var allowedModules []string
			if allowedModulesStr != "" {
				allowedModules = strings.Split(allowedModulesStr, ",")
				for i, module := range allowedModules {
					allowedModules[i] = strings.TrimSpace(module)
				}
			}

			deniedModulesStr, _ := cmd.Flags().GetString("denied-modules")
			var deniedModules []string
			if deniedModulesStr != "" {
				deniedModules = strings.Split(deniedModulesStr, ",")
				for i, module := range deniedModules {
					deniedModules[i] = strings.TrimSpace(module)
				}
			}

			purposeRestrictionsStr, _ := cmd.Flags().GetString("purpose-restrictions")
			var purposeRestrictions []string
			if purposeRestrictionsStr != "" {
				purposeRestrictions = strings.Split(purposeRestrictionsStr, ",")
				for i, purpose := range purposeRestrictions {
					purposeRestrictions[i] = strings.TrimSpace(purpose)
				}
			}

			dataRestrictionsStr, _ := cmd.Flags().GetString("data-restrictions")
			var dataRestrictions map[string][]string
			if dataRestrictionsStr != "" {
				if err := clientCtx.Codec.UnmarshalJSON([]byte(dataRestrictionsStr), &dataRestrictions); err != nil {
					return fmt.Errorf("invalid data restrictions JSON: %w", err)
				}
			}

			// Create time restrictions (empty for now)
			timeRestrictions := types.TimeRestriction{}

			authority := clientCtx.GetFromAddress()

			msg := types.NewMsgCreateAccessPolicy(
				authority.String(),
				holderDID,
				allowedModules,
				deniedModules,
				dataRestrictions,
				purposeRestrictions,
				timeRestrictions,
				[]string{}, // geographic restrictions
				maxSharesPerDay,
				requireExplicitConsent,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("allowed-modules", "", "Comma-separated list of allowed modules")
	cmd.Flags().String("denied-modules", "", "Comma-separated list of denied modules")
	cmd.Flags().String("purpose-restrictions", "", "Comma-separated list of allowed purposes")
	cmd.Flags().String("data-restrictions", "", "JSON map of credential types to allowed attributes")

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
