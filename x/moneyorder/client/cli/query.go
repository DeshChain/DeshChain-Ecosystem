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

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/deshchain/deshchain/x/moneyorder/types"
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

	cmd.AddCommand(
		CmdQueryParams(),
		CmdQueryReceipt(),
		CmdQueryPool(),
		CmdQueryUserOrders(),
	)

	return cmd
}

// CmdQueryParams returns the command to query module parameters
func CmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the current money order parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			// Query would be implemented with proper gRPC query client
			fmt.Println("Query params not yet implemented")
			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryReceipt returns the command to query a money order receipt
func CmdQueryReceipt() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "receipt [order-id-or-reference]",
		Short: "Query a money order receipt by order ID or reference number",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			orderRef := args[0]
			// Query would be implemented with proper gRPC query client
			fmt.Printf("Query receipt %s not yet implemented\n", orderRef)
			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryPool returns the command to query a pool
func CmdQueryPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pool [pool-id]",
		Short: "Query a pool by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			poolId, err := ParseUint64(args[0])
			if err != nil {
				return err
			}

			// Query would be implemented with proper gRPC query client
			fmt.Printf("Query pool %d not yet implemented\n", poolId)
			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryUserOrders returns the command to query user's orders
func CmdQueryUserOrders() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "orders [address]",
		Short: "Query all money orders for a user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			userAddr := args[0]
			// Query would be implemented with proper gRPC query client
			fmt.Printf("Query orders for %s not yet implemented\n", userAddr)
			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}