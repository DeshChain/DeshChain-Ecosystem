package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/deshchain/x/remittance/types"
)

// GetQueryIdentityCmd returns the identity query commands for Remittance module
func GetQueryIdentityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "identity",
		Short:                      "Querying commands for Remittance identity integration",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdQuerySenderIdentity(),
		CmdQueryRecipientIdentity(),
		CmdQuerySewaMitraIdentity(),
		CmdQueryRemittanceCompliance(),
	)

	return cmd
}

// CmdQuerySenderIdentity implements the sender identity query command
func CmdQuerySenderIdentity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sender [sender-address] [transfer-amount] [recipient-country]",
		Short: "Query sender identity and compliance information",
		Long: `Query comprehensive identity information for a remittance sender including:
- DID and identity status
- KYC verification status and level
- AML and sanctions screening status
- Source of funds verification
- Risk level assessment
- Transfer limits (max, daily, monthly)
- Country restrictions
- Compliance expiry

Example:
$ deshchaind query remittance identity sender cosmos1abc... 1000usd IN`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			senderAddress := args[0]
			transferAmount := args[1]
			recipientCountry := args[2]

			// Validate sender address
			_, err = sdk.AccAddressFromBech32(senderAddress)
			if err != nil {
				return fmt.Errorf("invalid sender address: %w", err)
			}

			// Validate transfer amount
			_, err = sdk.ParseCoinNormalized(transferAmount)
			if err != nil {
				return fmt.Errorf("invalid transfer amount: %w", err)
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QuerySenderIdentityRequest{
				SenderAddress:    senderAddress,
				TransferAmount:   transferAmount,
				RecipientCountry: recipientCountry,
			}

			res, err := queryClient.SenderIdentity(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryRecipientIdentity implements the recipient identity query command
func CmdQueryRecipientIdentity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "recipient [recipient-address] [sender-country]",
		Short: "Query recipient identity and verification information",
		Long: `Query comprehensive identity information for a remittance recipient including:
- DID and identity status
- KYC verification status and level
- Verification documents
- Purpose of funds
- Beneficiary type (individual, business, etc.)
- Country of residence
- Allowed sender countries

Example:
$ deshchaind query remittance identity recipient cosmos1abc... US`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			recipientAddress := args[0]
			senderCountry := args[1]

			// Validate recipient address
			_, err = sdk.AccAddressFromBech32(recipientAddress)
			if err != nil {
				return fmt.Errorf("invalid recipient address: %w", err)
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryRecipientIdentityRequest{
				RecipientAddress: recipientAddress,
				SenderCountry:    senderCountry,
			}

			res, err := queryClient.RecipientIdentity(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQuerySewaMitraIdentity implements the Sewa Mitra identity query command
func CmdQuerySewaMitraIdentity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sewa-mitra [agent-address] [agent-id]",
		Short: "Query Sewa Mitra agent identity and capabilities",
		Long: `Query comprehensive identity information for a Sewa Mitra agent including:
- DID and identity status
- KYC verification status and level
- Business license validation
- Compliance verification status
- AML compliance status
- Certifications and qualifications
- Service areas coverage
- Supported currencies
- Transaction limits
- Background verification
- Insurance coverage
- Performance rating

Example:
$ deshchaind query remittance identity sewa-mitra cosmos1abc... SMA-000001`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			agentAddress := args[0]
			agentID := args[1]

			// Validate agent address
			_, err = sdk.AccAddressFromBech32(agentAddress)
			if err != nil {
				return fmt.Errorf("invalid agent address: %w", err)
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QuerySewaMitraIdentityRequest{
				AgentAddress: agentAddress,
				AgentID:      agentID,
			}

			res, err := queryClient.SewaMitraIdentity(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryRemittanceCompliance implements the remittance compliance query command
func CmdQueryRemittanceCompliance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compliance [sender-address] [transfer-amount] [source-currency] [dest-currency] [sender-country] [recipient-country] --recipient-address [recipient-address] --agent-address [agent-address]",
		Short: "Query comprehensive remittance compliance check",
		Long: `Query comprehensive compliance verification for a remittance transaction including:
- Sender compliance verification
- Recipient compliance verification (if recipient address provided)
- Agent compliance verification (if agent address provided)
- Corridor allowability check
- Amount limits verification
- Sanctions screening results
- AML compliance status
- Regulatory compliance assessment
- Overall compliance score
- Required actions for compliance
- Warning messages

Example:
$ deshchaind query remittance identity compliance cosmos1sender... 1000usd USD INR US IN --recipient-address cosmos1recipient... --agent-address cosmos1agent...`,
		Args: cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			senderAddress := args[0]
			transferAmount := args[1]
			sourceCurrency := args[2]
			destCurrency := args[3]
			senderCountry := args[4]
			recipientCountry := args[5]

			// Get optional flags
			recipientAddress, _ := cmd.Flags().GetString("recipient-address")
			agentAddress, _ := cmd.Flags().GetString("agent-address")

			// Validate sender address
			_, err = sdk.AccAddressFromBech32(senderAddress)
			if err != nil {
				return fmt.Errorf("invalid sender address: %w", err)
			}

			// Validate transfer amount
			_, err = sdk.ParseCoinNormalized(transferAmount)
			if err != nil {
				return fmt.Errorf("invalid transfer amount: %w", err)
			}

			// Validate optional addresses
			if recipientAddress != "" {
				_, err = sdk.AccAddressFromBech32(recipientAddress)
				if err != nil {
					return fmt.Errorf("invalid recipient address: %w", err)
				}
			}

			if agentAddress != "" {
				_, err = sdk.AccAddressFromBech32(agentAddress)
				if err != nil {
					return fmt.Errorf("invalid agent address: %w", err)
				}
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryRemittanceComplianceRequest{
				SenderAddress:    senderAddress,
				RecipientAddress: recipientAddress,
				AgentAddress:     agentAddress,
				TransferAmount:   transferAmount,
				SourceCurrency:   sourceCurrency,
				DestCurrency:     destCurrency,
				SenderCountry:    senderCountry,
				RecipientCountry: recipientCountry,
			}

			res, err := queryClient.RemittanceCompliance(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().String("recipient-address", "", "Optional recipient address for compliance check")
	cmd.Flags().String("agent-address", "", "Optional Sewa Mitra agent address for compliance check")

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}