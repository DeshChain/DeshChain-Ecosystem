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
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"deshchain/x/donation/types"
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
		CmdDonate(),
		CmdUpdateParams(),
	)

	return cmd
}

// CmdDonate returns the donate command
func CmdDonate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "donate [ngo-wallet-id] [amount] [purpose]",
		Short: "Make a donation to an NGO",
		Long: `Make a donation to a registered NGO wallet. 

Examples:
$ deshchaind tx donation donate 1 1000000namo "Education support" --from mykey
$ deshchaind tx donation donate 2 5000000namo "Disaster relief" --anonymous --from mykey
$ deshchaind tx donation donate 3 10000000namo "Medical aid" --campaign-id 5 --from mykey`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			ngoWalletId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid NGO wallet ID: %w", err)
			}

			amount, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return fmt.Errorf("invalid amount: %w", err)
			}

			purpose := args[2]

			isAnonymous, err := cmd.Flags().GetBool("anonymous")
			if err != nil {
				return err
			}

			campaignId, err := cmd.Flags().GetUint64("campaign-id")
			if err != nil {
				return err
			}

			msg := &types.MsgDonate{
				Donor:        clientCtx.GetFromAddress().String(),
				NgoWalletId:  ngoWalletId,
				Amount:       amount,
				Purpose:      purpose,
				IsAnonymous:  isAnonymous,
				CampaignId:   campaignId,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().Bool("anonymous", false, "Make donation anonymous")
	cmd.Flags().Uint64("campaign-id", 0, "Campaign ID if donating to a specific campaign")

	return cmd
}

// CmdUpdateParams returns the update params command
func CmdUpdateParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-params [enabled]",
		Short: "Update donation module parameters",
		Long: `Update the donation module parameters. Only the governance module account can execute this.

Example:
$ deshchaind tx donation update-params true --from gov`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			enabled, err := strconv.ParseBool(args[0])
			if err != nil {
				return fmt.Errorf("invalid enabled value: %w", err)
			}

			msg := &types.MsgUpdateParams{
				Authority: clientCtx.GetFromAddress().String(),
				Params: types.Params{
					Enabled: enabled,
				},
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}