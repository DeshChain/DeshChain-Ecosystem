package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/DeshChain/DeshChain-Ecosystem/x/tradefinance/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdQueryParams())
	cmd.AddCommand(CmdShowParty())
	cmd.AddCommand(CmdListParties())
	cmd.AddCommand(CmdShowLc())
	cmd.AddCommand(CmdListLcs())
	cmd.AddCommand(CmdShowDocument())
	cmd.AddCommand(CmdListDocuments())
	cmd.AddCommand(CmdShowInsurancePolicy())
	cmd.AddCommand(CmdListInsurancePolicies())
	cmd.AddCommand(CmdShowShipment())
	cmd.AddCommand(CmdListShipments())
	cmd.AddCommand(CmdShowPaymentInstruction())
	cmd.AddCommand(CmdListPaymentInstructions())
	cmd.AddCommand(CmdQueryStats())
	cmd.AddCommand(CmdQueryLcsByParty())
	cmd.AddCommand(CmdQueryLcsByStatus())
	cmd.AddCommand(CmdQueryTradeVolume())

	return cmd
}

func CmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "shows the parameters of the module",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryParamsRequest{}

			res, err := queryClient.Params(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdShowParty() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "party [party-id]",
		Short: "shows a trade party",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetPartyRequest{
				PartyId: args[0],
			}

			res, err := queryClient.Party(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdListParties() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "parties",
		Short: "list all trade parties",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllPartiesRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.PartiesAll(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)

	return cmd
}

func CmdShowLc() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lc [lc-id]",
		Short: "shows a Letter of Credit",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetLcRequest{
				LcId: args[0],
			}

			res, err := queryClient.Lc(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdListLcs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lcs",
		Short: "list all Letters of Credit",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllLcsRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.LcsAll(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)

	return cmd
}

func CmdShowDocument() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "document [document-id]",
		Short: "shows a trade document",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetDocumentRequest{
				DocumentId: args[0],
			}

			res, err := queryClient.Document(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdListDocuments() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "documents [lc-id]",
		Short: "list all documents for an LC",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllDocumentsRequest{
				LcId:       args[0],
				Pagination: pageReq,
			}

			res, err := queryClient.DocumentsAll(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)

	return cmd
}

func CmdShowInsurancePolicy() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "insurance [policy-id]",
		Short: "shows an insurance policy",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetInsurancePolicyRequest{
				PolicyId: args[0],
			}

			res, err := queryClient.InsurancePolicy(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdListInsurancePolicies() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "insurances [lc-id]",
		Short: "list all insurance policies for an LC",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllInsurancePoliciesRequest{
				LcId:       args[0],
				Pagination: pageReq,
			}

			res, err := queryClient.InsurancePoliciesAll(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)

	return cmd
}

func CmdShowShipment() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "shipment [tracking-id]",
		Short: "shows shipment tracking information",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetShipmentRequest{
				TrackingId: args[0],
			}

			res, err := queryClient.Shipment(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdListShipments() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "shipments [lc-id]",
		Short: "list all shipments for an LC",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllShipmentsRequest{
				LcId:       args[0],
				Pagination: pageReq,
			}

			res, err := queryClient.ShipmentsAll(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)

	return cmd
}

func CmdShowPaymentInstruction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "payment [instruction-id]",
		Short: "shows a payment instruction",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetPaymentInstructionRequest{
				InstructionId: args[0],
			}

			res, err := queryClient.PaymentInstruction(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdListPaymentInstructions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "payments [lc-id]",
		Short: "list all payment instructions for an LC",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllPaymentInstructionsRequest{
				LcId:       args[0],
				Pagination: pageReq,
			}

			res, err := queryClient.PaymentInstructionsAll(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)

	return cmd
}

func CmdQueryStats() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats",
		Short: "shows trade finance statistics",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryStatsRequest{}

			res, err := queryClient.Stats(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// Advanced query commands

func CmdQueryLcsByParty() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lcs-by-party [party-id] [role]",
		Short: "list LCs by party role (applicant, beneficiary, issuing_bank)",
		Args:  cobra.ExactArgs(2),
		Example: fmt.Sprintf(`%s query tradefinance lcs-by-party PARTY001 applicant`, version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryLcsByPartyRequest{
				PartyId:    args[0],
				Role:       args[1],
				Pagination: pageReq,
			}

			res, err := queryClient.LcsByParty(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)

	return cmd
}

func CmdQueryLcsByStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lcs-by-status [status]",
		Short: "list LCs by status (draft, issued, accepted, documents_presented, paid, cancelled)",
		Args:  cobra.ExactArgs(1),
		Example: fmt.Sprintf(`%s query tradefinance lcs-by-status issued`, version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryLcsByStatusRequest{
				Status:     args[0],
				Pagination: pageReq,
			}

			res, err := queryClient.LcsByStatus(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)

	return cmd
}

func CmdQueryTradeVolume() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "trade-volume [days]",
		Short: "shows trade volume statistics for specified days",
		Args:  cobra.ExactArgs(1),
		Example: fmt.Sprintf(`%s query tradefinance trade-volume 30`, version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			days, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryTradeVolumeRequest{
				Days: uint32(days),
			}

			res, err := queryClient.TradeVolume(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}