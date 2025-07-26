package cli

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/namo/x/identity/types"
)

// Query commands for identity cache management

// CmdQueryCacheStats queries cache statistics
func CmdQueryCacheStats() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cache-stats",
		Short: "Query identity cache statistics",
		Long: `Query performance statistics for the identity cache system.

Examples:
deshchaind query identity cache-stats`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.CacheStats(context.Background(), &types.QueryCacheStatsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryCacheMetrics queries detailed cache metrics
func CmdQueryCacheMetrics() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cache-metrics",
		Short: "Query detailed cache performance metrics",
		Long: `Query detailed performance metrics for the identity cache system including hit ratios, operation times, and memory usage.

Examples:
deshchaind query identity cache-metrics`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.CacheMetrics(context.Background(), &types.QueryCacheMetricsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryCacheEntry queries a specific cache entry
func CmdQueryCacheEntry() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cache-entry [type] [key]",
		Short: "Query a specific cache entry",
		Long: `Query details of a specific cache entry by type and key.

Examples:
deshchaind query identity cache-entry identity desh1abc...
deshchaind query identity cache-entry credential cred_123
deshchaind query identity cache-entry did_document did:desh:user123`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			entryType := args[0]
			key := args[1]

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.CacheEntry(context.Background(), &types.QueryCacheEntryRequest{
				Type: entryType,
				Key:  key,
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

// CmdQueryCacheByTag queries cache entries by tag
func CmdQueryCacheByTag() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cache-by-tag [tag]",
		Short: "Query cache entries by tag",
		Long: `Query all cache entries that have a specific tag.

Examples:
deshchaind query identity cache-by-tag user_data
deshchaind query identity cache-by-tag verification
deshchaind query identity cache-by-tag did:desh:user123`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			tag := args[0]
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.CacheByTag(context.Background(), &types.QueryCacheByTagRequest{
				Tag:        tag,
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "cache entries")
	return cmd
}

// CmdQueryCacheByType queries cache entries by type
func CmdQueryCacheByType() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cache-by-type [type]",
		Short: "Query cache entries by type",
		Long: `Query all cache entries of a specific type.

Examples:
deshchaind query identity cache-by-type identity
deshchaind query identity cache-by-type credential
deshchaind query identity cache-by-type did_document`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			entryType := args[0]
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.CacheByType(context.Background(), &types.QueryCacheByTypeRequest{
				Type:       entryType,
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "cache entries")
	return cmd
}

// CmdQueryCacheHealth queries cache health status
func CmdQueryCacheHealth() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cache-health",
		Short: "Query cache health and performance status",
		Long: `Query the health status of the identity cache including performance indicators and recommendations.

Examples:
deshchaind query identity cache-health`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.CacheHealth(context.Background(), &types.QueryCacheHealthRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryCacheConfig queries cache configuration
func CmdQueryCacheConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cache-config",
		Short: "Query cache configuration settings",
		Long: `Query the current configuration settings for the identity cache.

Examples:
deshchaind query identity cache-config`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.CacheConfig(context.Background(), &types.QueryCacheConfigRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryCacheSize queries cache size information
func CmdQueryCacheSize() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cache-size",
		Short: "Query cache size and capacity information",
		Long: `Query information about cache size, entry count, and capacity utilization.

Examples:
deshchaind query identity cache-size`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.CacheSize(context.Background(), &types.QueryCacheSizeRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryCacheHitRatio queries cache hit ratio
func CmdQueryCacheHitRatio() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cache-hit-ratio [window-minutes]",
		Short: "Query cache hit ratio for a time window",
		Long: `Query the cache hit ratio for a specific time window in minutes.

Examples:
deshchaind query identity cache-hit-ratio 60    # Last 60 minutes
deshchaind query identity cache-hit-ratio 1440  # Last 24 hours`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			windowMinutes, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid window minutes: %w", err)
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.CacheHitRatio(context.Background(), &types.QueryCacheHitRatioRequest{
				WindowMinutes: windowMinutes,
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

// CmdQueryCacheTopKeys queries most accessed cache keys
func CmdQueryCacheTopKeys() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cache-top-keys [limit]",
		Short: "Query most frequently accessed cache keys",
		Long: `Query the most frequently accessed cache keys to identify hot data.

Examples:
deshchaind query identity cache-top-keys 10   # Top 10 keys
deshchaind query identity cache-top-keys 50   # Top 50 keys`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			limit, err := strconv.ParseInt(args[0], 10, 32)
			if err != nil {
				return fmt.Errorf("invalid limit: %w", err)
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.CacheTopKeys(context.Background(), &types.QueryCacheTopKeysRequest{
				Limit: int32(limit),
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