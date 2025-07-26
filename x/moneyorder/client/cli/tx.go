/*
Copyright 2024 DeshChain Foundation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/DeshChain/DeshChain-Ecosystem/x/moneyorder/types"
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
		CmdCreateMoneyOrder(),
		CmdCreateFixedRatePool(),
		CmdCreateVillagePool(),
		CmdSwap(),
		CmdJoinVillagePool(),
	)

	return cmd
}

// CmdCreateMoneyOrder returns a CLI command for creating a money order
func CmdCreateMoneyOrder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send [receiver-upi] [amount] [note]",
		Short: "Send money via Money Order (UPI-style)",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Send money to another user using their UPI address or wallet address.

Example:
$ %s tx moneyorder send ramesh@deshchain 100namo "For groceries"
$ %s tx moneyorder send desh1abc... 500namo "Monthly support"`,
				"deshchaind", "deshchaind",
			),
		),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			receiverUPI := args[0]
			amount, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}
			note := args[2]

			msg := types.NewMsgCreateMoneyOrder(
				clientCtx.GetFromAddress(),
				receiverUPI,
				amount,
				note,
				"instant", // Default to instant
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdCreateFixedRatePool returns a CLI command for creating a fixed rate pool
func CmdCreateFixedRatePool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-fixed-pool [token0] [token1] [rate] [initial-liquidity]",
		Short: "Create a fixed rate exchange pool",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			token0 := args[0]
			token1 := args[1]
			rate, err := sdk.NewDecFromStr(args[2])
			if err != nil {
				return err
			}
			liquidity, err := sdk.ParseCoinsNormalized(args[3])
			if err != nil {
				return err
			}

			description, _ := cmd.Flags().GetString("description")
			regions, _ := cmd.Flags().GetStringSlice("regions")

			msg := &types.MsgCreateFixedRatePool{
				Creator:          clientCtx.GetFromAddress().String(),
				Token0Denom:      token0,
				Token1Denom:      token1,
				ExchangeRate:     rate,
				InitialLiquidity: liquidity,
				Description:      description,
				SupportedRegions: regions,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("description", "", "Pool description")
	cmd.Flags().StringSlice("regions", []string{}, "Supported postal codes")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdCreateVillagePool returns a CLI command for creating a village pool
func CmdCreateVillagePool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-village-pool [village-name] [postal-code] [state] [district] [liquidity]",
		Short: "Create a community-managed village pool",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			villageName := args[0]
			postalCode := args[1]
			stateCode := args[2]
			districtCode := args[3]
			liquidity, err := sdk.ParseCoinsNormalized(args[4])
			if err != nil {
				return err
			}

			validators, _ := cmd.Flags().GetStringSlice("validators")

			msg := &types.MsgCreateVillagePool{
				PanchayatHead:    clientCtx.GetFromAddress().String(),
				VillageName:      villageName,
				PostalCode:       postalCode,
				StateCode:        stateCode,
				DistrictCode:     districtCode,
				InitialLiquidity: liquidity,
				LocalValidators:  validators,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().StringSlice("validators", []string{}, "Local validator addresses")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdSwap returns a CLI command for swapping tokens
func CmdSwap() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swap [pool-id] [token-in] [token-out-denom]",
		Short: "Swap tokens in a pool",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			poolId, err := ParseUint64(args[0])
			if err != nil {
				return err
			}

			tokenIn, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			tokenOutDenom := args[2]
			minOut, _ := cmd.Flags().GetString("min-out")
			minOutAmount := sdk.ZeroInt()
			if minOut != "" {
				minOutAmount, err = sdk.NewIntFromString(minOut)
				if err != nil {
					return err
				}
			}

			msg := &types.MsgSwapExactAmountIn{
				Sender:        clientCtx.GetFromAddress().String(),
				PoolId:        poolId,
				TokenIn:       tokenIn,
				TokenOutDenom: tokenOutDenom,
				TokenOutMin:   minOutAmount,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("min-out", "0", "Minimum output amount")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdJoinVillagePool returns a CLI command for joining a village pool
func CmdJoinVillagePool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "join-village [pool-id] [local-name] [mobile] [deposit]",
		Short: "Join a village pool as a member",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			poolId, err := ParseUint64(args[0])
			if err != nil {
				return err
			}

			localName := args[1]
			mobile := args[2]
			deposit, err := sdk.ParseCoinsNormalized(args[3])
			if err != nil {
				return err
			}

			msg := &types.MsgJoinVillagePool{
				Member:         clientCtx.GetFromAddress().String(),
				PoolId:         poolId,
				LocalName:      localName,
				MobileNumber:   mobile,
				InitialDeposit: deposit,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// ParseUint64 parses a string to uint64
func ParseUint64(s string) (uint64, error) {
	var n uint64
	_, err := fmt.Sscanf(s, "%d", &n)
	return n, err
}