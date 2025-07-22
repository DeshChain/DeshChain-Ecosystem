package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/deshchain/deshchain/x/dinr/types"
)

var (
	DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())
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
		CmdMintDINR(),
		CmdBurnDINR(),
		CmdDepositCollateral(),
		CmdWithdrawCollateral(),
		CmdLiquidate(),
	)

	return cmd
}

// CmdMintDINR returns a CLI command to mint DINR tokens
func CmdMintDINR() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint [collateral-amount] [dinr-amount]",
		Short: "Mint DINR by depositing collateral",
		Long: `Mint DINR stablecoins by depositing collateral. The collateral must meet the minimum ratio requirement.
Example:
$ deshchaind tx dinr mint 1000usdt 750dinr --from mykey`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			collateral, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			dinrToMint, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgMintDINR(
				clientCtx.GetFromAddress().String(),
				collateral,
				dinrToMint,
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

// CmdBurnDINR returns a CLI command to burn DINR tokens
func CmdBurnDINR() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn [dinr-amount] [collateral-denom]",
		Short: "Burn DINR and retrieve collateral",
		Long: `Burn DINR stablecoins and retrieve the specified collateral.
Example:
$ deshchaind tx dinr burn 500dinr usdt --from mykey`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			dinrToBurn, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			collateralDenom := args[1]

			msg := types.NewMsgBurnDINR(
				clientCtx.GetFromAddress().String(),
				dinrToBurn,
				collateralDenom,
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

// CmdDepositCollateral returns a CLI command to deposit additional collateral
func CmdDepositCollateral() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit-collateral [collateral-amount]",
		Short: "Deposit additional collateral to improve position health",
		Long: `Deposit additional collateral to an existing position to improve the collateral ratio.
Example:
$ deshchaind tx dinr deposit-collateral 500usdt --from mykey`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			collateral, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgDepositCollateral(
				clientCtx.GetFromAddress().String(),
				collateral,
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

// CmdWithdrawCollateral returns a CLI command to withdraw excess collateral
func CmdWithdrawCollateral() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-collateral [collateral-amount]",
		Short: "Withdraw excess collateral from position",
		Long: `Withdraw excess collateral from your position while maintaining the minimum ratio.
Example:
$ deshchaind tx dinr withdraw-collateral 200usdt --from mykey`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			collateral, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgWithdrawCollateral(
				clientCtx.GetFromAddress().String(),
				collateral,
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

// CmdLiquidate returns a CLI command to liquidate an undercollateralized position
func CmdLiquidate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "liquidate [user-address] [dinr-amount]",
		Short: "Liquidate an undercollateralized position",
		Long: `Liquidate an undercollateralized position by repaying DINR debt and receiving collateral with a discount.
Example:
$ deshchaind tx dinr liquidate desh1abc... 1000dinr --from mykey`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			userAddress := args[0]
			dinrToCover, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgLiquidate(
				clientCtx.GetFromAddress().String(),
				userAddress,
				dinrToCover,
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