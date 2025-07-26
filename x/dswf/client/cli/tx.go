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

	"github.com/DeshChain/DeshChain-Ecosystem/x/dswf/types"
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

	cmd.AddCommand(
		CmdProposeAllocation(),
		CmdApproveAllocation(),
		CmdExecuteDisbursement(),
		CmdUpdatePortfolio(),
		CmdSubmitMonthlyReport(),
		CmdRecordReturns(),
	)

	return cmd
}

// CmdProposeAllocation returns the propose allocation command
func CmdProposeAllocation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "propose-allocation [purpose] [amount] [category] [recipient] [justification] [expected-impact] [expected-returns] [risk-assessment]",
		Short: "Propose a new allocation from the DSWF",
		Args:  cobra.ExactArgs(8),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			purpose := args[0]
			amount, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}
			category := args[2]
			recipient := args[3]
			justification := args[4]
			expectedImpact := args[5]
			expectedReturns := args[6]
			riskAssessment := args[7]

			proposersStr, _ := cmd.Flags().GetString("proposers")
			proposers := strings.Split(proposersStr, ",")

			msg := &types.MsgProposeAllocation{
				Proposers:       proposers,
				Purpose:         purpose,
				Amount:          amount,
				Category:        category,
				Recipient:       recipient,
				Justification:   justification,
				ExpectedImpact:  expectedImpact,
				ExpectedReturns: expectedReturns,
				RiskAssessment:  riskAssessment,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().String("proposers", "", "Comma-separated list of proposer addresses")
	cmd.MarkFlagRequired("proposers")

	return cmd
}

// CmdApproveAllocation returns the approve allocation command
func CmdApproveAllocation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approve-allocation [allocation-id] [decision] [comments]",
		Short: "Approve or reject an allocation proposal",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			allocationID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			decision := args[1]
			comments := args[2]

			approversStr, _ := cmd.Flags().GetString("approvers")
			approvers := strings.Split(approversStr, ",")

			msg := &types.MsgApproveAllocation{
				AllocationId: allocationID,
				Approvers:    approvers,
				Decision:     decision,
				Comments:     comments,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().String("approvers", "", "Comma-separated list of approver addresses")
	cmd.MarkFlagRequired("approvers")

	return cmd
}

// CmdExecuteDisbursement returns the execute disbursement command
func CmdExecuteDisbursement() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "execute-disbursement [allocation-id] [disbursement-index] [verification-notes]",
		Short: "Execute a scheduled disbursement",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			allocationID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			disbursementIndex, err := strconv.ParseUint(args[1], 10, 32)
			if err != nil {
				return err
			}
			verificationNotes := args[2]

			msg := &types.MsgExecuteDisbursement{
				AllocationId:      allocationID,
				DisbursementIndex: uint32(disbursementIndex),
				Executor:          clientCtx.GetFromAddress().String(),
				VerificationNotes: verificationNotes,
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

// CmdUpdatePortfolio returns the update portfolio command
func CmdUpdatePortfolio() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-portfolio [total-value] [total-returns] [annual-return-rate] [risk-score] [rebalance-reason]",
		Short: "Update the DSWF investment portfolio",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			totalValue, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}
			totalReturns, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}
			annualReturnRate := args[2]
			riskScore, err := strconv.ParseInt(args[3], 10, 32)
			if err != nil {
				return err
			}
			rebalanceReason := args[4]

			msg := &types.MsgUpdatePortfolio{
				Authority:        clientCtx.GetFromAddress().String(),
				TotalValue:       totalValue,
				TotalReturns:     totalReturns,
				AnnualReturnRate: annualReturnRate,
				RiskScore:        int32(riskScore),
				RebalanceReason:  rebalanceReason,
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

// CmdSubmitMonthlyReport returns the submit monthly report command
func CmdSubmitMonthlyReport() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-monthly-report [period] [opening-balance] [closing-balance] [total-inflows] [total-outflows] [total-allocated] [total-disbursed] [total-returns] [average-return-rate]",
		Short: "Submit a monthly report for the DSWF",
		Args:  cobra.ExactArgs(9),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			period := args[0]
			openingBalance, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}
			closingBalance, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}
			totalInflows, err := sdk.ParseCoinNormalized(args[3])
			if err != nil {
				return err
			}
			totalOutflows, err := sdk.ParseCoinNormalized(args[4])
			if err != nil {
				return err
			}
			totalAllocated, err := sdk.ParseCoinNormalized(args[5])
			if err != nil {
				return err
			}
			totalDisbursed, err := sdk.ParseCoinNormalized(args[6])
			if err != nil {
				return err
			}
			totalReturns, err := sdk.ParseCoinNormalized(args[7])
			if err != nil {
				return err
			}
			averageReturnRate := args[8]

			msg := &types.MsgSubmitMonthlyReport{
				Reporter:          clientCtx.GetFromAddress().String(),
				Period:            period,
				OpeningBalance:    openingBalance,
				ClosingBalance:    closingBalance,
				TotalInflows:      totalInflows,
				TotalOutflows:     totalOutflows,
				TotalAllocated:    totalAllocated,
				TotalDisbursed:    totalDisbursed,
				TotalReturns:      totalReturns,
				AverageReturnRate: averageReturnRate,
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

// CmdRecordReturns returns the record returns command
func CmdRecordReturns() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "record-returns [allocation-id] [actual-returns] [period]",
		Short: "Record actual returns for an allocation",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			allocationID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			actualReturns, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}
			period, err := strconv.ParseInt(args[2], 10, 32)
			if err != nil {
				return err
			}

			msg := &types.MsgRecordReturns{
				AllocationId:  allocationID,
				ActualReturns: actualReturns,
				Period:        int32(period),
				Reporter:      clientCtx.GetFromAddress().String(),
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