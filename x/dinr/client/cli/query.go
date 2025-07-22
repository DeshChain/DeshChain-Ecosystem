package cli

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/deshchain/deshchain/x/dinr/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdQueryParams(),
		CmdQueryPosition(),
		CmdQueryAllPositions(),
		CmdQueryCollateralAsset(),
		CmdQueryAllCollateralAssets(),
		CmdQueryStability(),
		CmdQueryInsuranceFund(),
		CmdEstimateMint(),
		CmdEstimateBurn(),
		CmdQueryLiquidatable(),
	)

	return cmd
}

// CmdQueryParams queries the module parameters
func CmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the current DINR module parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(context.Background(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// CmdQueryPosition queries a user's DINR position
func CmdQueryPosition() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "position [address]",
		Short: "Query a user's DINR position",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.UserPosition(context.Background(), &types.QueryUserPositionRequest{
				Address: args[0],
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

// CmdQueryAllPositions queries all DINR positions
func CmdQueryAllPositions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "positions",
		Short: "Query all DINR positions",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.AllPositions(context.Background(), &types.QueryAllPositionsRequest{
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "positions")

	return cmd
}

// CmdQueryCollateralAsset queries information about a specific collateral type
func CmdQueryCollateralAsset() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collateral [denom]",
		Short: "Query information about a specific collateral asset",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.CollateralAsset(context.Background(), &types.QueryCollateralAssetRequest{
				Denom: args[0],
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

// CmdQueryAllCollateralAssets queries all supported collateral types
func CmdQueryAllCollateralAssets() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all-collateral",
		Short: "Query all supported collateral assets",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.AllCollateralAssets(context.Background(), &types.QueryAllCollateralAssetsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// CmdQueryStability queries current stability metrics
func CmdQueryStability() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stability",
		Short: "Query current DINR stability metrics",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.StabilityInfo(context.Background(), &types.QueryStabilityInfoRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// CmdQueryInsuranceFund queries the insurance fund status
func CmdQueryInsuranceFund() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "insurance",
		Short: "Query insurance fund status",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.InsuranceFund(context.Background(), &types.QueryInsuranceFundRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// CmdEstimateMint estimates the result of a mint operation
func CmdEstimateMint() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "estimate-mint [minter] [collateral-denom] [collateral-amount] [dinr-amount]",
		Short: "Estimate the result of a mint operation",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.EstimateMint(context.Background(), &types.QueryEstimateMintRequest{
				Minter:            args[0],
				CollateralDenom:   args[1],
				CollateralAmount:  args[2],
				DinrAmount:        args[3],
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

// CmdEstimateBurn estimates the result of a burn operation
func CmdEstimateBurn() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "estimate-burn [burner] [dinr-amount] [collateral-denom]",
		Short: "Estimate the result of a burn operation",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.EstimateBurn(context.Background(), &types.QueryEstimateBurnRequest{
				Burner:          args[0],
				DinrAmount:      args[1],
				CollateralDenom: args[2],
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

// CmdQueryLiquidatable queries positions eligible for liquidation
func CmdQueryLiquidatable() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "liquidatable",
		Short: "Query positions eligible for liquidation",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.LiquidatablePositions(context.Background(), &types.QueryLiquidatablePositionsRequest{
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "liquidatable positions")

	return cmd
}