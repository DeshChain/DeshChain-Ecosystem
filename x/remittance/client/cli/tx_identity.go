package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/remittance/types"
)

// GetTxIdentityCmd returns the identity transaction commands for Remittance module
func GetTxIdentityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "identity",
		Short:                      "Identity transaction subcommands for Remittance",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdMigrateRemittanceToIdentity(),
		CmdCreateSenderCredential(),
		CmdCreateSewaMitraCredential(),
		CmdCreateRecipientCredential(),
		CmdCreateTransferCredential(),
	)

	return cmd
}

// CmdMigrateRemittanceToIdentity implements the migrate remittance to identity command
func CmdMigrateRemittanceToIdentity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate-remittance",
		Short: "Migrate all remittance data to identity system",
		Long: `Migrate all existing transfers and Sewa Mitra agents in the Remittance module to the new identity system.
This command creates identities and credentials for all transfers and agents.

Example:
$ deshchaind tx remittance identity migrate-remittance --from mykey`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			authority := clientCtx.GetFromAddress()

			msg := types.NewMsgMigrateRemittanceToIdentity(authority.String())
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdCreateSenderCredential implements the create sender credential command
func CmdCreateSenderCredential() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-sender-credential [sender-address] [kyc-level] [source-of-funds] [risk-level] [max-transfer-limit] [daily-limit] [monthly-limit] [country-restrictions]",
		Short: "Create a remittance sender credential",
		Long: `Create a verifiable credential for a remittance sender with compliance information.
This includes KYC level, source of funds, risk assessment, transfer limits, and country restrictions.

Arguments:
- country-restrictions: comma-separated list of restricted countries (e.g., "US,CN" or "none")

Example:
$ deshchaind tx remittance identity create-sender-credential cosmos1abc... enhanced "employment" low "100000usd" "10000usd" "50000usd" "none" --from mykey`,
		Args: cobra.ExactArgs(8),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			authority := clientCtx.GetFromAddress()
			senderAddress := args[0]
			kycLevel := args[1]
			sourceOfFunds := args[2]
			riskLevel := args[3]
			maxTransferLimit := args[4]
			dailyLimit := args[5]
			monthlyLimit := args[6]
			countryRestrictionsStr := args[7]

			// Validate sender address
			_, err = sdk.AccAddressFromBech32(senderAddress)
			if err != nil {
				return fmt.Errorf("invalid sender address: %w", err)
			}

			// Validate coin amounts
			if maxTransferLimit != "none" {
				_, err = sdk.ParseCoinNormalized(maxTransferLimit)
				if err != nil {
					return fmt.Errorf("invalid max transfer limit: %w", err)
				}
			} else {
				maxTransferLimit = ""
			}

			if dailyLimit != "none" {
				_, err = sdk.ParseCoinNormalized(dailyLimit)
				if err != nil {
					return fmt.Errorf("invalid daily limit: %w", err)
				}
			} else {
				dailyLimit = ""
			}

			if monthlyLimit != "none" {
				_, err = sdk.ParseCoinNormalized(monthlyLimit)
				if err != nil {
					return fmt.Errorf("invalid monthly limit: %w", err)
				}
			} else {
				monthlyLimit = ""
			}

			// Parse country restrictions
			var countryRestrictions []string
			if countryRestrictionsStr != "none" && countryRestrictionsStr != "" {
				countryRestrictions = strings.Split(countryRestrictionsStr, ",")
				for i, country := range countryRestrictions {
					countryRestrictions[i] = strings.TrimSpace(country)
				}
			}

			msg := types.NewMsgCreateSenderCredential(authority.String(), senderAddress, kycLevel, sourceOfFunds, riskLevel, maxTransferLimit, dailyLimit, monthlyLimit, countryRestrictions)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdCreateSewaMitraCredential implements the create Sewa Mitra credential command
func CmdCreateSewaMitraCredential() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-sewa-mitra-credential [agent-address] [agent-id] [business-name] [service-areas] [supported-currencies] [max-transaction-limit] [certifications]",
		Short: "Create a Sewa Mitra agent credential",
		Long: `Create a verifiable credential for a Sewa Mitra agent with service capabilities.
This includes business information, service areas, supported currencies, and certifications.

Arguments:
- service-areas: comma-separated list of service areas (e.g., "Mumbai,Delhi")
- supported-currencies: comma-separated list of currencies (e.g., "USD,INR,EUR")
- certifications: comma-separated list of certifications (e.g., "AML,KYC,Remittance" or "none")

Example:
$ deshchaind tx remittance identity create-sewa-mitra-credential cosmos1abc... SMA-000001 "Mumbai Remittance Services" "Mumbai,Pune" "USD,INR" "50000usd" "AML,KYC" --from mykey`,
		Args: cobra.ExactArgs(7),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			authority := clientCtx.GetFromAddress()
			agentAddress := args[0]
			agentID := args[1]
			businessName := args[2]
			serviceAreasStr := args[3]
			supportedCurrenciesStr := args[4]
			maxTransactionLimit := args[5]
			certificationsStr := args[6]

			// Validate agent address
			_, err = sdk.AccAddressFromBech32(agentAddress)
			if err != nil {
				return fmt.Errorf("invalid agent address: %w", err)
			}

			// Validate max transaction limit
			if maxTransactionLimit != "none" {
				_, err = sdk.ParseCoinNormalized(maxTransactionLimit)
				if err != nil {
					return fmt.Errorf("invalid max transaction limit: %w", err)
				}
			} else {
				maxTransactionLimit = ""
			}

			// Parse service areas
			var serviceAreas []string
			if serviceAreasStr != "" {
				serviceAreas = strings.Split(serviceAreasStr, ",")
				for i, area := range serviceAreas {
					serviceAreas[i] = strings.TrimSpace(area)
				}
			}

			// Parse supported currencies
			var supportedCurrencies []string
			if supportedCurrenciesStr != "" {
				supportedCurrencies = strings.Split(supportedCurrenciesStr, ",")
				for i, currency := range supportedCurrencies {
					supportedCurrencies[i] = strings.TrimSpace(currency)
				}
			}

			// Parse certifications
			var certifications []string
			if certificationsStr != "none" && certificationsStr != "" {
				certifications = strings.Split(certificationsStr, ",")
				for i, cert := range certifications {
					certifications[i] = strings.TrimSpace(cert)
				}
			}

			msg := types.NewMsgCreateSewaMitraCredential(authority.String(), agentAddress, agentID, businessName, serviceAreas, supportedCurrencies, maxTransactionLimit, certifications)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdCreateRecipientCredential implements the create recipient credential command
func CmdCreateRecipientCredential() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-recipient-credential [recipient-address] [kyc-level] [verification-docs] [purpose-of-funds] [beneficiary-type] [country] [can-receive-from]",
		Short: "Create a remittance recipient credential",
		Long: `Create a verifiable credential for a remittance recipient with verification information.
This includes KYC level, verification documents, purpose of funds, and receiving restrictions.

Arguments:
- verification-docs: comma-separated list of documents (e.g., "passport,license")
- can-receive-from: comma-separated list of allowed sender countries (e.g., "US,IN,UK" or "any")

Example:
$ deshchaind tx remittance identity create-recipient-credential cosmos1abc... basic "passport,license" "family_support" "individual" "IN" "US,UK" --from mykey`,
		Args: cobra.ExactArgs(7),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			authority := clientCtx.GetFromAddress()
			recipientAddress := args[0]
			kycLevel := args[1]
			verificationDocsStr := args[2]
			purposeOfFunds := args[3]
			beneficiaryType := args[4]
			country := args[5]
			canReceiveFromStr := args[6]

			// Validate recipient address
			_, err = sdk.AccAddressFromBech32(recipientAddress)
			if err != nil {
				return fmt.Errorf("invalid recipient address: %w", err)
			}

			// Parse verification documents
			var verificationDocs []string
			if verificationDocsStr != "" && verificationDocsStr != "none" {
				verificationDocs = strings.Split(verificationDocsStr, ",")
				for i, doc := range verificationDocs {
					verificationDocs[i] = strings.TrimSpace(doc)
				}
			}

			// Parse can receive from countries
			var canReceiveFrom []string
			if canReceiveFromStr != "any" && canReceiveFromStr != "" {
				canReceiveFrom = strings.Split(canReceiveFromStr, ",")
				for i, country := range canReceiveFrom {
					canReceiveFrom[i] = strings.TrimSpace(country)
				}
			}

			msg := types.NewMsgCreateRecipientCredential(authority.String(), recipientAddress, kycLevel, verificationDocs, purposeOfFunds, beneficiaryType, country, canReceiveFrom)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdCreateTransferCredential implements the create transfer credential command
func CmdCreateTransferCredential() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-transfer-credential [transfer-id]",
		Short: "Create a transfer completion credential",
		Long: `Create a verifiable credential for a completed remittance transfer.
This credential serves as proof of successful transfer completion.

Example:
$ deshchaind tx remittance identity create-transfer-credential RMT-123456 --from mykey`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			authority := clientCtx.GetFromAddress()
			transferID := args[0]

			msg := types.NewMsgCreateTransferCredential(authority.String(), transferID)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}