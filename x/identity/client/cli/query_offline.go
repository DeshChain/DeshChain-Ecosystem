package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/deshchain/deshchain/x/identity/types"
)

// GetQueryOfflineCmd returns the query commands for offline verification
func GetQueryOfflineCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "offline",
		Short:                      "Querying commands for offline verification",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdQueryOfflineConfig(),
		CmdQueryOfflineDevices(),
		CmdQueryOfflineFormats(),
		CmdQueryOfflineModes(),
		CmdVerifyOfflineData(),
		CmdValidateOfflinePackage(),
	)

	return cmd
}

// CmdQueryOfflineConfig queries the offline verification configuration
func CmdQueryOfflineConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Query offline verification configuration",
		Long: `Query the current system-wide offline verification configuration including security thresholds,
cache settings, compression options, and supported regions.

Example:
$ deshchaind query identity offline config
`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.OfflineConfig(context.Background(), &types.QueryOfflineConfigRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryOfflineDevices queries registered offline devices for an identity
func CmdQueryOfflineDevices() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "devices [did]",
		Short: "Query registered offline devices for an identity",
		Long: `Query all devices registered for offline verification for a specific identity.
Shows device information including capabilities, security levels, and registration status.

Example:
$ deshchaind query identity offline devices did:desh:user123
$ deshchaind query identity offline devices did:desh:user456 --device-id="device_001"
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			did := args[0]
			deviceID, _ := cmd.Flags().GetString("device-id")

			res, err := queryClient.OfflineDevices(context.Background(), &types.QueryOfflineDevicesRequest{
				Did:      did,
				DeviceId: deviceID,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().String("device-id", "", "Query specific device by ID")
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryOfflineFormats queries supported offline credential formats
func CmdQueryOfflineFormats() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "formats",
		Short: "Query supported offline credential formats",
		Long: `Query all supported offline credential formats with their descriptions and use cases.

Example:
$ deshchaind query identity offline formats
`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.OfflineFormats(context.Background(), &types.QueryOfflineFormatsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryOfflineModes queries supported offline verification modes
func CmdQueryOfflineModes() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "modes",
		Short: "Query supported offline verification modes",
		Long: `Query all supported offline verification modes with their security levels and use cases.

Example:
$ deshchaind query identity offline modes
`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.OfflineModes(context.Background(), &types.QueryOfflineModesRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdVerifyOfflineData performs offline verification of identity data
func CmdVerifyOfflineData() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify [offline-data-file] [verification-request-file]",
		Short: "Verify identity using offline verification data",
		Long: `Perform identity verification using offline verification data and a verification request.
This simulates the offline verification process that would occur without network connectivity.

The offline-data-file should contain the offline verification data structure.
The verification-request-file should contain the verification request parameters.

Example:
$ deshchaind query identity offline verify offline_data.json verification_request.json
$ deshchaind query identity offline verify offline_data.json verification_request.json --output=json
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			offlineDataFile := args[0]
			verificationRequestFile := args[1]

			// Read offline data file
			offlineDataBytes, err := ioutil.ReadFile(offlineDataFile)
			if err != nil {
				return fmt.Errorf("failed to read offline data file: %w", err)
			}

			var offlineData types.OfflineVerificationData
			if err := json.Unmarshal(offlineDataBytes, &offlineData); err != nil {
				return fmt.Errorf("failed to parse offline data: %w", err)
			}

			// Read verification request file
			requestBytes, err := ioutil.ReadFile(verificationRequestFile)
			if err != nil {
				return fmt.Errorf("failed to read verification request file: %w", err)
			}

			var request types.OfflineVerificationRequest
			if err := json.Unmarshal(requestBytes, &request); err != nil {
				return fmt.Errorf("failed to parse verification request: %w", err)
			}

			// Perform verification
			res, err := queryClient.VerifyOffline(context.Background(), &types.QueryVerifyOfflineRequest{
				OfflineData: &offlineData,
				Request:     &request,
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

// CmdValidateOfflinePackage validates an offline verification package
func CmdValidateOfflinePackage() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate [offline-package-file]",
		Short: "Validate an offline verification package",
		Long: `Validate the integrity and structure of an offline verification package.
Checks data integrity, expiration, signatures, and format compliance.

Example:
$ deshchaind query identity offline validate offline_package.json
$ deshchaind query identity offline validate compressed_package.bin --format=compressed
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			packageFile := args[0]
			format, _ := cmd.Flags().GetString("format")

			// Read package file
			packageBytes, err := ioutil.ReadFile(packageFile)
			if err != nil {
				return fmt.Errorf("failed to read package file: %w", err)
			}

			res, err := queryClient.ValidateOfflinePackage(context.Background(), &types.QueryValidateOfflinePackageRequest{
				PackageData: packageBytes,
				Format:      format,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().String("format", "self_contained", "Package format (self_contained, compressed, qr_code, nfc, printable)")
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// Helper commands for offline verification testing and debugging

// CmdCreateSampleOfflineData creates sample offline data for testing
func CmdCreateSampleOfflineData() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-sample [did] [output-file]",
		Short: "Create sample offline verification data for testing",
		Long: `Create sample offline verification data for testing and development purposes.
This is useful for testing offline verification flows without real identity data.

WARNING: This creates sample data only and should not be used in production.

Example:
$ deshchaind query identity offline create-sample did:desh:test123 sample_offline_data.json
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			did := args[0]
			outputFile := args[1]

			// Create sample offline verification data
			sampleData := createSampleOfflineData(did)

			// Marshal to JSON
			jsonData, err := json.MarshalIndent(sampleData, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal sample data: %w", err)
			}

			// Write to file
			if err := ioutil.WriteFile(outputFile, jsonData, 0644); err != nil {
				return fmt.Errorf("failed to write sample data file: %w", err)
			}

			fmt.Printf("Sample offline verification data created: %s\n", outputFile)
			return nil
		},
	}

	return cmd
}

// CmdCreateSampleVerificationRequest creates sample verification request for testing
func CmdCreateSampleVerificationRequest() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-sample-request [did] [output-file]",
		Short: "Create sample verification request for testing",
		Long: `Create sample verification request for testing offline verification flows.

Example:
$ deshchaind query identity offline create-sample-request did:desh:test123 sample_request.json
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			did := args[0]
			outputFile := args[1]

			requiredLevel, _ := cmd.Flags().GetUint32("required-level")
			requiredTypes, _ := cmd.Flags().GetStringSlice("required-types")

			// Create sample verification request
			sampleRequest := createSampleVerificationRequest(did, requiredLevel, requiredTypes)

			// Marshal to JSON
			jsonData, err := json.MarshalIndent(sampleRequest, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal sample request: %w", err)
			}

			// Write to file
			if err := ioutil.WriteFile(outputFile, jsonData, 0644); err != nil {
				return fmt.Errorf("failed to write sample request file: %w", err)
			}

			fmt.Printf("Sample verification request created: %s\n", outputFile)
			return nil
		},
	}

	cmd.Flags().Uint32("required-level", 2, "Required verification level")
	cmd.Flags().StringSlice("required-types", []string{"IdentityCredential"}, "Required credential types")

	return cmd
}

// Helper functions for creating sample data

func createSampleOfflineData(did string) *types.OfflineVerificationData {
	sampleData := types.NewOfflineVerificationData(
		did,
		"sample_identity_hash_123456",
		"sample_public_key_abcdef",
		"did:desh:issuer",
		2,
	)

	// Add sample identity proof
	sampleData.IdentityProof = types.NewCryptographicProof(
		"Ed25519Signature2020",
		"sample_identity_proof_value",
		"assertionMethod",
		"sample_verification_method",
	)

	// Add sample KYC proof
	sampleData.KYCProof = types.NewCryptographicProof(
		"KYCProof2020",
		"sample_kyc_proof_value",
		"assertionMethod",
		"sample_kyc_verification_method",
	)

	// Add sample credential
	sampleCredential := types.NewOfflineCredential(
		"sample_credential_001",
		[]string{"VerifiableCredential", "IdentityCredential"},
		did,
		map[string]interface{}{
			"id":   did,
			"name": "Sample User",
			"type": "identity",
		},
	)
	sampleData.Credentials = []*types.OfflineCredential{sampleCredential}

	// Add sample revocation data
	sampleData.RevocationData = &types.RevocationData{
		RevocationListURL:  "https://example.com/revocation",
		RevocationListHash: "sample_revocation_hash",
		LastUpdated:        sampleData.IssuedAt,
		ValidUntil:         sampleData.ExpiresAt,
		RevokedCredentials: []string{},
		RevokedIdentities:  []string{},
		MerkleRoot:         "sample_merkle_root",
	}

	// Add sample localization data
	sampleData.LocalizedData = &types.OfflineLocalizationData{
		DefaultLanguage: types.LanguageEnglish,
		Messages: map[string]map[types.LanguageCode]string{
			"verification_success": {
				types.LanguageEnglish: "Verification successful",
				types.LanguageHindi:   "सत्यापन सफल",
			},
		},
		ErrorMessages: map[string]map[types.LanguageCode]string{
			"verification_failed": {
				types.LanguageEnglish: "Verification failed",
				types.LanguageHindi:   "सत्यापन असफल",
			},
		},
	}

	// Compute and set data hash
	if hash, err := sampleData.ComputeDataHash(); err == nil {
		sampleData.DataHash = hash
		sampleData.Signature = "sample_signature_" + hash[:16]
	}

	return sampleData
}

func createSampleVerificationRequest(did string, requiredLevel uint32, requiredTypes []string) *types.OfflineVerificationRequest {
	return types.NewOfflineVerificationRequest(
		did,
		"sample_challenge_123456",
		requiredLevel,
		requiredTypes,
	)
}