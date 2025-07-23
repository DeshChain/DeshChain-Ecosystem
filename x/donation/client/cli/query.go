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
	"context"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"deshchain/x/donation/types"
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
		CmdQueryNGOWallet(),
		CmdQueryNGOWallets(),
		CmdQueryActiveNGOs(),
		CmdQueryDonationRecord(),
		CmdQueryDonationsByDonor(),
		CmdQueryDonationsByNGO(),
		CmdQueryStatistics(),
		CmdQueryCampaign(),
		CmdQueryActiveCampaigns(),
		CmdQueryTransparencyScore(),
	)

	return cmd
}

// CmdQueryParams returns the params query command
func CmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the current donation module parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

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

// CmdQueryNGOWallet returns the NGO wallet query command
func CmdQueryNGOWallet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ngo [ngo-id]",
		Short: "Query an NGO wallet by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			ngoId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid NGO ID: %w", err)
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.NGOWallet(context.Background(), &types.QueryNGOWalletRequest{
				NgoWalletId: ngoId,
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

// CmdQueryNGOWallets returns the NGO wallets list query command
func CmdQueryNGOWallets() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ngos",
		Short: "Query all NGO wallets",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.NGOWallets(context.Background(), &types.QueryNGOWalletsRequest{
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "ngos")

	return cmd
}

// CmdQueryActiveNGOs returns the active NGOs query command
func CmdQueryActiveNGOs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "active-ngos",
		Short: "Query all active and verified NGOs",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.ActiveNGOs(context.Background(), &types.QueryActiveNGOsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// CmdQueryDonationRecord returns the donation record query command
func CmdQueryDonationRecord() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "donation [donation-id]",
		Short: "Query a donation record by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			donationId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid donation ID: %w", err)
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.DonationRecord(context.Background(), &types.QueryDonationRecordRequest{
				DonationId: donationId,
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

// CmdQueryDonationsByDonor returns the donations by donor query command
func CmdQueryDonationsByDonor() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "donor-donations [donor-address]",
		Short: "Query all donations by a donor",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.DonationsByDonor(context.Background(), &types.QueryDonationsByDonorRequest{
				Donor:      args[0],
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "donations")

	return cmd
}

// CmdQueryDonationsByNGO returns the donations by NGO query command
func CmdQueryDonationsByNGO() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ngo-donations [ngo-id]",
		Short: "Query all donations received by an NGO",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			ngoId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid NGO ID: %w", err)
			}

			queryClient := types.NewQueryClient(clientCtx)
			
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.DonationsByNGO(context.Background(), &types.QueryDonationsByNGORequest{
				NgoWalletId: ngoId,
				Pagination:  pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "donations")

	return cmd
}

// CmdQueryStatistics returns the statistics query command
func CmdQueryStatistics() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "statistics",
		Short: "Query donation module statistics",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Statistics(context.Background(), &types.QueryStatisticsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// CmdQueryCampaign returns the campaign query command
func CmdQueryCampaign() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "campaign [campaign-id]",
		Short: "Query a campaign by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			campaignId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid campaign ID: %w", err)
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Campaign(context.Background(), &types.QueryCampaignRequest{
				CampaignId: campaignId,
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

// CmdQueryActiveCampaigns returns the active campaigns query command
func CmdQueryActiveCampaigns() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "active-campaigns",
		Short: "Query all active campaigns",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.ActiveCampaigns(context.Background(), &types.QueryActiveCampaignsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// CmdQueryTransparencyScore returns the transparency score query command
func CmdQueryTransparencyScore() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transparency-score [ngo-id]",
		Short: "Query transparency score for an NGO",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			ngoId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid NGO ID: %w", err)
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.TransparencyScore(context.Background(), &types.QueryTransparencyScoreRequest{
				NgoWalletId: ngoId,
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