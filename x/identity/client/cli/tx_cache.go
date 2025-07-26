package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/namo/x/identity/types"
)

// Transaction commands for identity cache management

// CmdClearCache clears the identity cache
func CmdClearCache() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clear-cache",
		Short: "Clear the identity cache",
		Long: `Clear all entries from the identity cache. This operation requires admin privileges.

Examples:
deshchaind tx identity clear-cache --from admin-key`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgClearCache{
				Authority: clientCtx.GetFromAddress().String(),
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdRefreshCache refreshes the identity cache
func CmdRefreshCache() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "refresh-cache",
		Short: "Refresh the identity cache",
		Long: `Refresh the identity cache by clearing and reloading frequently accessed data.

Examples:
deshchaind tx identity refresh-cache --from admin-key
deshchaind tx identity refresh-cache --preload-identities --from admin-key`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			preloadIdentities, _ := cmd.Flags().GetBool("preload-identities")
			preloadCredentials, _ := cmd.Flags().GetBool("preload-credentials")
			preloadDIDs, _ := cmd.Flags().GetBool("preload-dids")

			msg := &types.MsgRefreshCache{
				Authority:          clientCtx.GetFromAddress().String(),
				PreloadIdentities:  preloadIdentities,
				PreloadCredentials: preloadCredentials,
				PreloadDIDs:        preloadDIDs,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().Bool("preload-identities", true, "Preload active identities after refresh")
	cmd.Flags().Bool("preload-credentials", false, "Preload recent credentials after refresh")
	cmd.Flags().Bool("preload-dids", false, "Preload popular DID documents after refresh")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// CmdUpdateCacheConfig updates cache configuration
func CmdUpdateCacheConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-cache-config",
		Short: "Update cache configuration",
		Long: `Update the identity cache configuration settings.

Examples:
deshchaind tx identity update-cache-config \
  --max-size 104857600 \
  --max-entries 10000 \
  --default-ttl 1800 \
  --from admin-key`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Parse configuration flags
			maxSize, _ := cmd.Flags().GetInt64("max-size")
			maxEntries, _ := cmd.Flags().GetInt64("max-entries")
			defaultTTL, _ := cmd.Flags().GetInt64("default-ttl")
			cleanupInterval, _ := cmd.Flags().GetInt64("cleanup-interval")
			enableMetrics, _ := cmd.Flags().GetBool("enable-metrics")
			enableTags, _ := cmd.Flags().GetBool("enable-tags")

			// Parse type-specific TTLs
			identityTTL, _ := cmd.Flags().GetInt64("identity-ttl")
			credentialTTL, _ := cmd.Flags().GetInt64("credential-ttl")
			didDocumentTTL, _ := cmd.Flags().GetInt64("did-document-ttl")
			consentTTL, _ := cmd.Flags().GetInt64("consent-ttl")
			zkProofTTL, _ := cmd.Flags().GetInt64("zk-proof-ttl")

			config := types.CacheConfig{
				MaxSize:         maxSize,
				MaxEntries:      maxEntries,
				DefaultTTL:      defaultTTL,
				CleanupInterval: cleanupInterval,
				EnableMetrics:   enableMetrics,
				EnableTags:      enableTags,
				IdentityTTL:     identityTTL,
				CredentialTTL:   credentialTTL,
				DIDDocumentTTL:  didDocumentTTL,
				ConsentTTL:      consentTTL,
				ZKProofTTL:      zkProofTTL,
			}

			msg := &types.MsgUpdateCacheConfig{
				Authority: clientCtx.GetFromAddress().String(),
				Config:    config,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	// Cache size and capacity
	cmd.Flags().Int64("max-size", 0, "Maximum cache size in bytes")
	cmd.Flags().Int64("max-entries", 0, "Maximum number of cache entries")
	cmd.Flags().Int64("default-ttl", 0, "Default TTL in seconds")
	cmd.Flags().Int64("cleanup-interval", 0, "Cleanup interval in seconds")

	// Feature flags
	cmd.Flags().Bool("enable-metrics", true, "Enable cache metrics collection")
	cmd.Flags().Bool("enable-tags", true, "Enable tag-based cache operations")

	// Type-specific TTLs
	cmd.Flags().Int64("identity-ttl", 0, "Identity cache TTL in seconds")
	cmd.Flags().Int64("credential-ttl", 0, "Credential cache TTL in seconds")
	cmd.Flags().Int64("did-document-ttl", 0, "DID document cache TTL in seconds")
	cmd.Flags().Int64("consent-ttl", 0, "Consent cache TTL in seconds")
	cmd.Flags().Int64("zk-proof-ttl", 0, "ZK proof cache TTL in seconds")

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdInvalidateCache invalidates specific cache entries
func CmdInvalidateCache() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "invalidate-cache [type] [keys]",
		Short: "Invalidate specific cache entries",
		Long: `Invalidate specific cache entries by type and keys, or by tag.

Examples:
deshchaind tx identity invalidate-cache identity desh1abc...,desh1def... --from admin-key
deshchaind tx identity invalidate-cache --by-tag user_data --from admin-key
deshchaind tx identity invalidate-cache --by-pattern "did:desh:user*" --from admin-key`,
		Args: cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			byTag, _ := cmd.Flags().GetString("by-tag")
			byPattern, _ := cmd.Flags().GetString("by-pattern")

			var msg *types.MsgInvalidateCache

			if byTag != "" {
				// Invalidate by tag
				msg = &types.MsgInvalidateCache{
					Authority:        clientCtx.GetFromAddress().String(),
					InvalidationType: types.InvalidationType_BY_TAG,
					Tag:              byTag,
				}
			} else if byPattern != "" {
				// Invalidate by pattern
				msg = &types.MsgInvalidateCache{
					Authority:        clientCtx.GetFromAddress().String(),
					InvalidationType: types.InvalidationType_BY_PATTERN,
					Pattern:          byPattern,
				}
			} else if len(args) >= 2 {
				// Invalidate by type and keys
				entryType := args[0]
				keys := strings.Split(args[1], ",")

				msg = &types.MsgInvalidateCache{
					Authority:        clientCtx.GetFromAddress().String(),
					InvalidationType: types.InvalidationType_BY_KEYS,
					Type:             entryType,
					Keys:             keys,
				}
			} else {
				return fmt.Errorf("must specify either type/keys, --by-tag, or --by-pattern")
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("by-tag", "", "Invalidate all entries with this tag")
	cmd.Flags().String("by-pattern", "", "Invalidate all entries matching this pattern")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// CmdWarmupCache preloads data into cache
func CmdWarmupCache() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "warmup-cache",
		Short: "Warm up the cache with frequently accessed data",
		Long: `Preload frequently accessed data into the cache to improve performance.

Examples:
deshchaind tx identity warmup-cache --identities 100 --credentials 500 --from admin-key
deshchaind tx identity warmup-cache --by-tags "user_data,verification" --from admin-key`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			maxIdentities, _ := cmd.Flags().GetInt32("identities")
			maxCredentials, _ := cmd.Flags().GetInt32("credentials")
			maxDIDs, _ := cmd.Flags().GetInt32("dids")
			tagsStr, _ := cmd.Flags().GetString("by-tags")

			var tags []string
			if tagsStr != "" {
				tags = strings.Split(tagsStr, ",")
				for i, tag := range tags {
					tags[i] = strings.TrimSpace(tag)
				}
			}

			strategy := types.CacheWarmupStrategy{
				PreloadActiveIdentities:  maxIdentities > 0,
				PreloadRecentCredentials: maxCredentials > 0,
				PreloadPopularDIDs:       maxDIDs > 0,
				PreloadByTags:           tags,
				MaxPreloadEntries:       int64(maxIdentities + maxCredentials + maxDIDs),
				PreloadBatchSize:        100,
			}

			msg := &types.MsgWarmupCache{
				Authority: clientCtx.GetFromAddress().String(),
				Strategy:  strategy,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().Int32("identities", 0, "Number of identities to preload")
	cmd.Flags().Int32("credentials", 0, "Number of credentials to preload")
	cmd.Flags().Int32("dids", 0, "Number of DID documents to preload")
	cmd.Flags().String("by-tags", "", "Comma-separated list of tags to preload")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// CmdSetCacheExpiry sets custom expiry for cache entries
func CmdSetCacheExpiry() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-cache-expiry [type] [key] [ttl-seconds]",
		Short: "Set custom expiry time for cache entries",
		Long: `Set a custom expiry time for specific cache entries.

Examples:
deshchaind tx identity set-cache-expiry identity desh1abc... 3600 --from admin-key
deshchaind tx identity set-cache-expiry credential cred_123 1800 --from admin-key`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			entryType := args[0]
			key := args[1]
			ttlSeconds, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid TTL seconds: %w", err)
			}

			msg := &types.MsgSetCacheExpiry{
				Authority:  clientCtx.GetFromAddress().String(),
				Type:       entryType,
				Key:        key,
				TTLSeconds: ttlSeconds,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdOptimizeCache triggers cache optimization
func CmdOptimizeCache() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "optimize-cache",
		Short: "Optimize cache performance",
		Long: `Trigger cache optimization including cleanup of expired entries, defragmentation, and performance tuning.

Examples:
deshchaind tx identity optimize-cache --from admin-key
deshchaind tx identity optimize-cache --aggressive --from admin-key`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			aggressive, _ := cmd.Flags().GetBool("aggressive")
			cleanupExpired, _ := cmd.Flags().GetBool("cleanup-expired")
			defragment, _ := cmd.Flags().GetBool("defragment")
			rebalance, _ := cmd.Flags().GetBool("rebalance")

			msg := &types.MsgOptimizeCache{
				Authority:      clientCtx.GetFromAddress().String(),
				Aggressive:     aggressive,
				CleanupExpired: cleanupExpired,
				Defragment:     defragment,
				Rebalance:      rebalance,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().Bool("aggressive", false, "Perform aggressive optimization")
	cmd.Flags().Bool("cleanup-expired", true, "Clean up expired entries")
	cmd.Flags().Bool("defragment", false, "Defragment cache memory")
	cmd.Flags().Bool("rebalance", false, "Rebalance cache distribution")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}