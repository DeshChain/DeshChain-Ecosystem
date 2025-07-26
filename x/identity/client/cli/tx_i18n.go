package cli

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/DeshChain/DeshChain-Ecosystem/x/identity/types"
)

// GetTxI18nCmd returns the transaction commands for internationalization
func GetTxI18nCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "i18n",
		Short:                      "Internationalization transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdSetUserLanguage(),
		CmdAddCustomMessage(),
		CmdUpdateLocalizationConfig(),
		CmdImportMessages(),
	)

	return cmd
}

// CmdSetUserLanguage sets user's preferred language
func CmdSetUserLanguage() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-language [user-did] [language-code]",
		Short: "Set user's preferred language",
		Long: `Set the preferred language for a user's identity interactions.

Example:
$ deshchaind tx identity i18n set-language did:desh:user123 hi --from mykey

Supported language codes:
- hi: Hindi
- en: English  
- bn: Bengali
- ta: Tamil
- te: Telugu
- mr: Marathi
- gu: Gujarati
- kn: Kannada
- ml: Malayalam
- ur: Urdu
- pa: Punjabi
- or: Odia
- as: Assamese
- sa: Sanskrit
- ne: Nepali
- gom: Konkani
- mni: Manipuri
- brx: Bodo
- sat: Santali
- ks: Kashmiri
- mai: Maithili
- doi: Dogri
- sd: Sindhi
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			userDID := args[0]
			languageCode := types.LanguageCode(args[1])

			// Validate language code
			supportedLanguages := types.GetSupportedLanguages()
			isSupported := false
			for _, lang := range supportedLanguages {
				if lang == languageCode {
					isSupported = true
					break
				}
			}

			if !isSupported {
				return fmt.Errorf("unsupported language code: %s", languageCode)
			}

			msg := &types.MsgSetUserLanguage{
				UserDid:      userDID,
				LanguageCode: string(languageCode),
				Signer:       clientCtx.GetFromAddress().String(),
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdAddCustomMessage adds a custom localized message
func CmdAddCustomMessage() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-message [key] [category] [translations-json]",
		Short: "Add a custom localized message",
		Long: `Add a custom message with translations in multiple languages.

The translations should be provided as a JSON object with language codes as keys.

Example:
$ deshchaind tx identity i18n add-message "custom.greeting" "greetings" \
  '{"en":"Hello","hi":"नमस्ते","ta":"வணக்கம்"}' --from mykey
`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			key := args[0]
			category := args[1]
			translationsJSON := args[2]

			// Parse translations JSON
			var translations map[string]string
			if err := json.Unmarshal([]byte(translationsJSON), &translations); err != nil {
				return fmt.Errorf("invalid translations JSON: %w", err)
			}

			// Convert to LocalizedText
			if len(translations) == 0 {
				return fmt.Errorf("at least one translation is required")
			}

			// Use first language as default
			var defaultLang types.LanguageCode
			for langCode := range translations {
				defaultLang = types.LanguageCode(langCode)
				break
			}

			localizedText := types.NewLocalizedText(defaultLang, translations[string(defaultLang)])
			for langCode, text := range translations {
				if types.LanguageCode(langCode) != defaultLang {
					localizedText.AddTranslation(types.LanguageCode(langCode), text)
				}
			}

			msg := &types.MsgAddCustomMessage{
				Key:         key,
				Category:    category,
				Text:        localizedText,
				Description: fmt.Sprintf("Custom message: %s", key),
				Signer:      clientCtx.GetFromAddress().String(),
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdUpdateLocalizationConfig updates the localization configuration
func CmdUpdateLocalizationConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-config [config-json]",
		Short: "Update localization configuration",
		Long: `Update the localization configuration for the identity module.

The configuration should be provided as a JSON object.

Example:
$ deshchaind tx identity i18n update-config \
  '{"default_language":"hi","fallback_languages":["hi","en"],"region":"india","enable_auto_detect":true,"enable_rtl_support":true}' \
  --from admin
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			configJSON := args[0]

			// Parse configuration JSON
			var config types.LocalizationConfig
			if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
				return fmt.Errorf("invalid configuration JSON: %w", err)
			}

			msg := &types.MsgUpdateLocalizationConfig{
				Config: &config,
				Signer: clientCtx.GetFromAddress().String(),
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// CmdImportMessages imports messages from a JSON file
func CmdImportMessages() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import-messages [catalog-json]",
		Short: "Import messages from a message catalog",
		Long: `Import multiple messages from a message catalog JSON.

The catalog should be a JSON object containing message definitions.

Example:
$ deshchaind tx identity i18n import-messages \
  '{"messages":{"msg1":{"key":"msg1","category":"test","text":{"translations":{"en":"Hello","hi":"नमस्ते"}}}}}' \
  --from admin
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			catalogJSON := args[0]

			// Parse catalog JSON
			var catalog types.MessageCatalog
			if err := json.Unmarshal([]byte(catalogJSON), &catalog); err != nil {
				return fmt.Errorf("invalid catalog JSON: %w", err)
			}

			msg := &types.MsgImportMessages{
				Catalog: &catalog,
				Signer:  clientCtx.GetFromAddress().String(),
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}