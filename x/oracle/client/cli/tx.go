package cli

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/DeshChain/DeshChain-Ecosystem/x/oracle/types"
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

	cmd.AddCommand(CmdSubmitPrice())
	cmd.AddCommand(CmdSubmitExchangeRate())
	cmd.AddCommand(CmdRegisterOracleValidator())
	cmd.AddCommand(CmdUpdateOracleValidator())

	return cmd
}

func CmdSubmitPrice() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-price [symbol] [price] [source]",
		Short: "Submit a price for a symbol",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argSymbol := args[0]
			argPrice := args[1]
			argSource := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			price, err := sdk.NewDecFromStr(argPrice)
			if err != nil {
				return fmt.Errorf("invalid price: %w", err)
			}

			msg := types.NewMsgSubmitPrice(
				clientCtx.GetFromAddress().String(),
				argSymbol,
				price,
				argSource,
				time.Now(),
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

func CmdSubmitExchangeRate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-exchange-rate [base] [target] [rate] [source]",
		Short: "Submit an exchange rate between two currencies",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argBase := args[0]
			argTarget := args[1]
			argRate := args[2]
			argSource := args[3]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			rate, err := sdk.NewDecFromStr(argRate)
			if err != nil {
				return fmt.Errorf("invalid rate: %w", err)
			}

			msg := types.NewMsgSubmitExchangeRate(
				clientCtx.GetFromAddress().String(),
				argBase,
				argTarget,
				rate,
				argSource,
				time.Now(),
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

func CmdRegisterOracleValidator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-oracle-validator [validator] [power] [description]",
		Short: "Register a new oracle validator",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argValidator := args[0]
			argPower := args[1]
			argDescription := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			power, err := strconv.ParseUint(argPower, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid power: %w", err)
			}

			msg := types.NewMsgRegisterOracleValidator(
				clientCtx.GetFromAddress().String(),
				argValidator,
				power,
				argDescription,
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

func CmdUpdateOracleValidator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-oracle-validator [validator] [power] [active]",
		Short: "Update an oracle validator",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argValidator := args[0]
			argPower := args[1]
			argActive := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			power, err := strconv.ParseUint(argPower, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid power: %w", err)
			}

			active, err := strconv.ParseBool(argActive)
			if err != nil {
				return fmt.Errorf("invalid active status: %w", err)
			}

			msg := types.NewMsgUpdateOracleValidator(
				clientCtx.GetFromAddress().String(),
				argValidator,
				power,
				active,
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