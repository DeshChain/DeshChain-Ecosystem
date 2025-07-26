package cli

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/DeshChain/DeshChain-Ecosystem/x/identity/types"
)

// GetTxOfflineCmd returns the transaction commands for offline verification
func GetTxOfflineCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "offline",
		Short:                      "Offline verification transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdPrepareOfflineVerification(),
		CmdCreateOfflineBackup(),
		CmdUpdateOfflineConfig(),
		CmdRegisterOfflineDevice(),
		CmdRevokeOfflineAccess(),
	)

	return cmd
}

// CmdPrepareOfflineVerification prepares offline verification data
func CmdPrepareOfflineVerification() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prepare-verification [did] [format] [expiration-hours]",
		Short: "Prepare offline verification data for an identity",
		Long: `Prepare offline verification data that can be used for identity verification without network connectivity.

Supported formats:
- self_contained: Complete data package (default)
- compressed: Compressed format for space efficiency  
- qr_code: QR code encodable format
- nfc: NFC tag compatible format
- printable: Human-readable printable format

Example:
$ deshchaind tx identity offline prepare-verification did:desh:user123 qr_code 24 --from mykey
$ deshchaind tx identity offline prepare-verification did:desh:user456 self_contained 72 --include-biometric --include-credentials="IdentityCredential,EducationCredential" --from mykey
`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			did := args[0]
			format := types.OfflineCredentialFormat(args[1])
			expirationHours, err := strconv.ParseUint(args[2], 10, 32)
			if err != nil {
				return fmt.Errorf("invalid expiration hours: %v", err)
			}

			// Get flags
			includeBiometric, _ := cmd.Flags().GetBool("include-biometric")
			includeCredentialsStr, _ := cmd.Flags().GetString("include-credentials")
			requiredLevel, _ := cmd.Flags().GetUint32("required-level")
			deviceID, _ := cmd.Flags().GetString("device-id")

			// Parse credentials list
			var includeCredentials []string
			if includeCredentialsStr != "" {
				includeCredentials = strings.Split(includeCredentialsStr, ",")
				for i, cred := range includeCredentials {
					includeCredentials[i] = strings.TrimSpace(cred)
				}
			}

			msg := &types.MsgPrepareOfflineVerification{
				DID:                did,
				Format:             format,
				ExpirationHours:    uint32(expirationHours),
				IncludeBiometric:   includeBiometric,
				IncludeCredentials: includeCredentials,
				RequiredLevel:      requiredLevel,
				DeviceID:           deviceID,
				Signer:             clientCtx.GetFromAddress().String(),
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().Bool("include-biometric", false, "Include biometric templates for offline verification")
	cmd.Flags().String("include-credentials", "", "Comma-separated list of credential types to include")
	cmd.Flags().Uint32("required-level", 2, "Required verification level (1-5)")
	cmd.Flags().String("device-id", "", "Target device ID for offline verification")

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdCreateOfflineBackup creates an offline backup package
func CmdCreateOfflineBackup() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-backup [did]",
		Short: "Create a comprehensive offline backup for an identity",
		Long: `Create a comprehensive offline backup package that includes identity data, credentials, 
recovery methods, and emergency contacts for offline access and disaster recovery.

SECURITY WARNING: Backups with private data should be encrypted and stored securely.

Example:
$ deshchaind tx identity offline create-backup did:desh:user123 --from mykey
$ deshchaind tx identity offline create-backup did:desh:user456 --include-private-data --backup-password="strong_password_123" --backup-location="/secure/backup/path" --from mykey
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			did := args[0]

			// Get flags
			includePrivateData, _ := cmd.Flags().GetBool("include-private-data")
			backupPassword, _ := cmd.Flags().GetString("backup-password")
			backupLocation, _ := cmd.Flags().GetString("backup-location")

			msg := &types.MsgCreateOfflineBackup{
				DID:                did,
				IncludePrivateData: includePrivateData,
				BackupPassword:     backupPassword,
				BackupLocation:     backupLocation,
				Signer:             clientCtx.GetFromAddress().String(),
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().Bool("include-private-data", false, "Include private data such as recovery methods (requires password)")
	cmd.Flags().String("backup-password", "", "Password for encrypting private data in backup")
	cmd.Flags().String("backup-location", "", "Optional backup storage location or identifier")

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdUpdateOfflineConfig updates offline verification configuration
func CmdUpdateOfflineConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-config",
		Short: "Update offline verification configuration",
		Long: `Update the global offline verification configuration parameters.

This command requires admin privileges and updates system-wide settings for offline verification.

Example:
$ deshchaind tx identity offline update-config --max-offline-hours=48 --required-confidence=0.9 --emergency-mode=true --from admin
`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Get flags
			maxOfflineHours, _ := cmd.Flags().GetUint32("max-offline-hours")
			requiredConfidence, _ := cmd.Flags().GetFloat64("required-confidence")
			biometricThreshold, _ := cmd.Flags().GetFloat64("biometric-threshold")
			maxCacheMB, _ := cmd.Flags().GetUint64("max-cache-mb")
			cacheExpirationHours, _ := cmd.Flags().GetUint32("cache-expiration-hours")
			enableCompression, _ := cmd.Flags().GetBool("enable-compression")
			compressionLevel, _ := cmd.Flags().GetUint32("compression-level")
			emergencyMode, _ := cmd.Flags().GetBool("emergency-mode")
			emergencyThreshold, _ := cmd.Flags().GetFloat64("emergency-threshold")
			defaultLanguage, _ := cmd.Flags().GetString("default-language")
			supportedRegions, _ := cmd.Flags().GetStringSlice("supported-regions")

			// Create config with current values or defaults
			config := types.DefaultOfflineVerificationConfig()

			// Update config with provided values
			if cmd.Flags().Changed("max-offline-hours") {
				config.MaxOfflineDuration = time.Duration(maxOfflineHours) * time.Hour
			}
			if cmd.Flags().Changed("required-confidence") {
				config.RequiredConfidence = requiredConfidence
			}
			if cmd.Flags().Changed("biometric-threshold") {
				config.BiometricThreshold = biometricThreshold
			}
			if cmd.Flags().Changed("max-cache-mb") {
				config.MaxCacheSize = maxCacheMB * 1024 * 1024 // Convert MB to bytes
			}
			if cmd.Flags().Changed("cache-expiration-hours") {
				config.CacheExpirationPeriod = time.Duration(cacheExpirationHours) * time.Hour
			}
			if cmd.Flags().Changed("enable-compression") {
				config.EnableCompression = enableCompression
			}
			if cmd.Flags().Changed("compression-level") {
				config.CompressionLevel = compressionLevel
			}
			if cmd.Flags().Changed("emergency-mode") {
				config.EmergencyModeEnabled = emergencyMode
			}
			if cmd.Flags().Changed("emergency-threshold") {
				config.EmergencyThreshold = emergencyThreshold
			}
			if cmd.Flags().Changed("default-language") {
				config.DefaultLanguage = types.LanguageCode(defaultLanguage)
			}
			if cmd.Flags().Changed("supported-regions") {
				config.SupportedRegions = supportedRegions
			}

			msg := &types.MsgUpdateOfflineConfig{
				Config: config,
				Signer: clientCtx.GetFromAddress().String(),
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().Uint32("max-offline-hours", 24, "Maximum hours for offline verification validity")
	cmd.Flags().Float64("required-confidence", 0.85, "Required confidence score for verification (0.0-1.0)")
	cmd.Flags().Float64("biometric-threshold", 0.95, "Biometric matching threshold (0.0-1.0)")
	cmd.Flags().Uint64("max-cache-mb", 100, "Maximum cache size in megabytes")
	cmd.Flags().Uint32("cache-expiration-hours", 168, "Cache expiration period in hours (default: 7 days)")
	cmd.Flags().Bool("enable-compression", true, "Enable data compression for offline packages")
	cmd.Flags().Uint32("compression-level", 6, "Compression level (1-9, higher = better compression)")
	cmd.Flags().Bool("emergency-mode", true, "Enable emergency mode with reduced security")
	cmd.Flags().Float64("emergency-threshold", 0.70, "Emergency mode confidence threshold (0.0-1.0)")
	cmd.Flags().String("default-language", "en", "Default language for offline messages")
	cmd.Flags().StringSlice("supported-regions", []string{"india", "global"}, "Supported regions for offline verification")

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdRegisterOfflineDevice registers a device for offline verification
func CmdRegisterOfflineDevice() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-device [did] [device-id] [device-name] [device-type] [public-key]",
		Short: "Register a device for offline verification",
		Long: `Register a device that can perform offline identity verification.

Supported device types:
- mobile: Mobile phone or tablet
- tablet: Tablet device
- laptop: Laptop computer
- desktop: Desktop computer  
- iot: IoT device
- smartcard: Smart card reader
- hardware_token: Hardware security token

Supported capabilities:
- identity_verification: Basic identity verification
- biometric_capture: Biometric data capture
- credential_storage: Secure credential storage
- offline_signing: Offline transaction signing
- secure_element: Hardware secure element

Example:
$ deshchaind tx identity offline register-device did:desh:user123 "device_001" "John's iPhone" mobile "pubkey_abc123" --from mykey
$ deshchaind tx identity offline register-device did:desh:user456 "laptop_001" "Work Laptop" laptop "pubkey_def456" --capabilities="identity_verification,biometric_capture" --security-level=3 --max-offline-hours=48 --from mykey
`,
		Args: cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			did := args[0]
			deviceID := args[1]
			deviceName := args[2]
			deviceType := args[3]
			publicKey := args[4]

			// Get flags
			capabilitiesStr, _ := cmd.Flags().GetString("capabilities")
			securityLevel, _ := cmd.Flags().GetUint32("security-level")
			maxOfflineHours, _ := cmd.Flags().GetUint32("max-offline-hours")

			// Parse capabilities
			var capabilities []string
			if capabilitiesStr != "" {
				capabilities = strings.Split(capabilitiesStr, ",")
				for i, cap := range capabilities {
					capabilities[i] = strings.TrimSpace(cap)
				}
			} else {
				capabilities = []string{"identity_verification"} // Default capability
			}

			msg := &types.MsgRegisterOfflineDevice{
				DID:                did,
				DeviceID:           deviceID,
				DeviceName:         deviceName,
				DeviceType:         deviceType,
				PublicKey:          publicKey,
				Capabilities:       capabilities,
				SecurityLevel:      securityLevel,
				MaxOfflineDuration: maxOfflineHours,
				Signer:             clientCtx.GetFromAddress().String(),
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("capabilities", "identity_verification", "Comma-separated list of device capabilities")
	cmd.Flags().Uint32("security-level", 2, "Device security level (1-5)")
	cmd.Flags().Uint32("max-offline-hours", 24, "Maximum offline duration in hours")

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdRevokeOfflineAccess revokes offline access for devices
func CmdRevokeOfflineAccess() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke-access [did]",
		Short: "Revoke offline verification access",
		Long: `Revoke offline verification access for a specific device or all devices associated with an identity.

This is useful when:
- A device is lost or stolen
- A device is compromised
- Updating security policies
- Decommissioning devices

Example:
$ deshchaind tx identity offline revoke-access did:desh:user123 --device-id="device_001" --reason="Device lost" --from mykey
$ deshchaind tx identity offline revoke-access did:desh:user456 --reason="Security policy update" --from mykey  # Revokes all devices
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			did := args[0]

			// Get flags
			deviceID, _ := cmd.Flags().GetString("device-id")
			reason, _ := cmd.Flags().GetString("reason")

			msg := &types.MsgRevokeOfflineAccess{
				DID:      did,
				DeviceID: deviceID,
				Reason:   reason,
				Signer:   clientCtx.GetFromAddress().String(),
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String("device-id", "", "Specific device ID to revoke (empty = revoke all devices)")
	cmd.Flags().String("reason", "", "Reason for revocation")

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}