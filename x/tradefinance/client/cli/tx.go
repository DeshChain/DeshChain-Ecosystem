package cli

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/deshchain/deshchain/x/tradefinance/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdRegisterParty())
	cmd.AddCommand(CmdIssueLc())
	cmd.AddCommand(CmdAcceptLc())
	cmd.AddCommand(CmdSubmitDocuments())
	cmd.AddCommand(CmdVerifyDocument())
	cmd.AddCommand(CmdRequestPayment())
	cmd.AddCommand(CmdMakePayment())
	cmd.AddCommand(CmdAmendLc())
	cmd.AddCommand(CmdCancelLc())
	cmd.AddCommand(CmdCreateInsurancePolicy())
	cmd.AddCommand(CmdUpdateShipment())

	return cmd
}

func CmdRegisterParty() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-party [party-type] [name] [registration-number] [country] [industry]",
		Short: "Register a new trade party (bank, exporter, importer, etc.)",
		Args:  cobra.ExactArgs(5),
		Example: fmt.Sprintf(`%s tx tradefinance register-party bank "State Bank of India" "REG001" "India" "Banking"`, version.AppName),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.MsgRegisterParty{
				Creator:            clientCtx.GetFromAddress().String(),
				PartyType:          args[0],
				Name:               args[1],
				RegistrationNumber: args[2],
				Country:            args[3],
				Industry:           args[4],
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdIssueLc() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue-lc [applicant-id] [beneficiary-id] [lc-number] [amount] [collateral] [expiry-days] [description]",
		Short: "Issue a new Letter of Credit",
		Args:  cobra.ExactArgs(7),
		Example: fmt.Sprintf(`%s tx tradefinance issue-lc PARTY001 PARTY002 LC2024001 1000000dinr 1200000dinr 90 "Cotton export LC"`, version.AppName),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinNormalized(args[3])
			if err != nil {
				return err
			}

			collateral, err := sdk.ParseCoinNormalized(args[4])
			if err != nil {
				return err
			}

			expiryDays, err := strconv.ParseInt(args[5], 10, 64)
			if err != nil {
				return err
			}

			expiryDate := time.Now().AddDate(0, 0, int(expiryDays))

			msg := types.MsgIssueLc{
				IssuingBank:   clientCtx.GetFromAddress().String(),
				ApplicantId:   args[0],
				BeneficiaryId: args[1],
				LcNumber:      args[2],
				Amount:        amount,
				Collateral:    collateral,
				ExpiryDate:    expiryDate,
				Description:   args[6],
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdAcceptLc() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "accept-lc [lc-id]",
		Short: "Accept a Letter of Credit as beneficiary",
		Args:  cobra.ExactArgs(1),
		Example: fmt.Sprintf(`%s tx tradefinance accept-lc LC001`, version.AppName),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.MsgAcceptLc{
				Beneficiary: clientCtx.GetFromAddress().String(),
				LcId:        args[0],
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdSubmitDocuments() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-documents [lc-id] [document-types] [ipfs-hashes]",
		Short: "Submit trade documents for LC",
		Args:  cobra.ExactArgs(3),
		Long: `Submit trade documents for a Letter of Credit.
Document types and IPFS hashes should be comma-separated lists in the same order.`,
		Example: fmt.Sprintf(`%s tx tradefinance submit-documents LC001 "invoice,bill_of_lading,packing_list" "QmHash1,QmHash2,QmHash3"`, version.AppName),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			docTypes := strings.Split(args[1], ",")
			ipfsHashes := strings.Split(args[2], ",")

			if len(docTypes) != len(ipfsHashes) {
				return fmt.Errorf("number of document types must match number of IPFS hashes")
			}

			msg := types.MsgSubmitDocuments{
				Submitter:     clientCtx.GetFromAddress().String(),
				LcId:          args[0],
				DocumentTypes: docTypes,
				IpfsHashes:    ipfsHashes,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdVerifyDocument() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify-document [document-id] [verification-result] [comments]",
		Short: "Verify a trade document",
		Args:  cobra.ExactArgs(3),
		Example: fmt.Sprintf(`%s tx tradefinance verify-document DOC001 approved "All details verified"`, version.AppName),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.MsgVerifyDocument{
				Verifier:           clientCtx.GetFromAddress().String(),
				DocumentId:         args[0],
				VerificationResult: args[1],
				Comments:           args[2],
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdRequestPayment() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "request-payment [lc-id] [amount] [justification]",
		Short: "Request payment under LC",
		Args:  cobra.ExactArgs(3),
		Example: fmt.Sprintf(`%s tx tradefinance request-payment LC001 1000000dinr "All documents verified and compliant"`, version.AppName),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.MsgRequestPayment{
				Requestor:     clientCtx.GetFromAddress().String(),
				LcId:          args[0],
				Amount:        amount,
				Justification: args[2],
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdMakePayment() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "make-payment [instruction-id]",
		Short: "Process payment under LC",
		Args:  cobra.ExactArgs(1),
		Example: fmt.Sprintf(`%s tx tradefinance make-payment INST001`, version.AppName),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.MsgMakePayment{
				Payer:         clientCtx.GetFromAddress().String(),
				InstructionId: args[0],
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdAmendLc() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "amend-lc [lc-id] [amendment-type] [new-value] [justification]",
		Short: "Amend a Letter of Credit",
		Args:  cobra.ExactArgs(4),
		Long: `Amend a Letter of Credit. Amendment types include:
- amount: Change LC amount
- expiry: Change expiry date
- terms: Modify terms and conditions
- documents: Change required documents`,
		Example: fmt.Sprintf(`%s tx tradefinance amend-lc LC001 amount 1200000dinr "Increased order quantity"`, version.AppName),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.MsgAmendLc{
				Requester:      clientCtx.GetFromAddress().String(),
				LcId:           args[0],
				AmendmentType:  args[1],
				NewValue:       args[2],
				Justification:  args[3],
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdCancelLc() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-lc [lc-id] [reason]",
		Short: "Cancel a Letter of Credit",
		Args:  cobra.ExactArgs(2),
		Example: fmt.Sprintf(`%s tx tradefinance cancel-lc LC001 "Order cancelled by buyer"`, version.AppName),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.MsgCancelLc{
				Requester: clientCtx.GetFromAddress().String(),
				LcId:      args[0],
				Reason:    args[1],
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdCreateInsurancePolicy() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-insurance [lc-id] [coverage-amount] [premium] [policy-type] [terms]",
		Short: "Create insurance policy for LC",
		Args:  cobra.ExactArgs(5),
		Example: fmt.Sprintf(`%s tx tradefinance create-insurance LC001 1000000dinr 25000dinr marine "Full marine coverage"`, version.AppName),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			coverageAmount, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			premium, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}

			msg := types.MsgCreateInsurancePolicy{
				Creator:        clientCtx.GetFromAddress().String(),
				LcId:           args[0],
				CoverageAmount: coverageAmount,
				Premium:        premium,
				PolicyType:     args[3],
				Terms:          args[4],
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdUpdateShipment() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-shipment [lc-id] [tracking-id] [status] [location] [notes]",
		Short: "Update shipment tracking information",
		Args:  cobra.ExactArgs(5),
		Example: fmt.Sprintf(`%s tx tradefinance update-shipment LC001 TRACK001 in_transit "Mumbai Port" "Departed on schedule"`, version.AppName),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.MsgUpdateShipment{
				Updater:    clientCtx.GetFromAddress().String(),
				LcId:       args[0],
				TrackingId: args[1],
				Status:     args[2],
				Location:   args[3],
				Notes:      args[4],
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}