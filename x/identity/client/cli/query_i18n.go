package cli

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/DeshChain/DeshChain-Ecosystem/x/identity/types"
)

// GetQueryI18nCmd returns the query commands for internationalization
func GetQueryI18nCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "i18n",
		Short:                      "Querying commands for internationalization",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdQuerySupportedLanguages(),
		CmdQueryUserLanguage(),
		CmdQueryLocalizedMessage(),
		CmdQueryRegionalLanguages(),
		CmdQueryLanguageInfo(),
		CmdQueryCulturalQuote(),
		CmdQueryGreeting(),
		CmdQueryLocalizationConfig(),
		CmdQueryCustomMessages(),
		CmdQueryLocalizationStats(),
	)

	return cmd
}

// CmdQuerySupportedLanguages queries all supported languages
func CmdQuerySupportedLanguages() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "supported-languages",
		Short: "Query all supported languages",
		Long: `Query all languages supported by the DeshChain identity system.

Example:
$ deshchaind query identity i18n supported-languages
`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.SupportedLanguages(context.Background(), &types.QuerySupportedLanguagesRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryUserLanguage queries user's language preference
func CmdQueryUserLanguage() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user-language [user-did]",
		Short: "Query user's language preference",
		Long: `Query the preferred language setting for a specific user.

Example:
$ deshchaind query identity i18n user-language did:desh:user123
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			userDID := args[0]

			res, err := queryClient.UserLanguage(context.Background(), &types.QueryUserLanguageRequest{
				UserDid: userDID,
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

// CmdQueryLocalizedMessage queries a localized message
func CmdQueryLocalizedMessage() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "message [key] [language-code]",
		Short: "Query a localized message",
		Long: `Query a specific message in the requested language.

Example:
$ deshchaind query identity i18n message "auth.success" "hi"
$ deshchaind query identity i18n message "welcome" "en"
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			key := args[0]
			languageCode := args[1]

			res, err := queryClient.LocalizedMessage(context.Background(), &types.QueryLocalizedMessageRequest{
				Key:          key,
				LanguageCode: languageCode,
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

// CmdQueryRegionalLanguages queries languages for a specific region
func CmdQueryRegionalLanguages() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "regional-languages [region]",
		Short: "Query languages for a specific region",
		Long: `Query the languages commonly used in a specific region or state.

Example:
$ deshchaind query identity i18n regional-languages "maharashtra"
$ deshchaind query identity i18n regional-languages "tamil_nadu"
$ deshchaind query identity i18n regional-languages "gujarat"
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			region := args[0]

			res, err := queryClient.RegionalLanguages(context.Background(), &types.QueryRegionalLanguagesRequest{
				Region: region,
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

// CmdQueryLanguageInfo queries detailed information about a language
func CmdQueryLanguageInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "language-info [language-code]",
		Short: "Query detailed information about a language",
		Long: `Query detailed information about a specific language including script, direction, and regional information.

Example:
$ deshchaind query identity i18n language-info "hi"
$ deshchaind query identity i18n language-info "ta"
$ deshchaind query identity i18n language-info "ur"
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			languageCode := args[0]

			res, err := queryClient.LanguageInfo(context.Background(), &types.QueryLanguageInfoRequest{
				LanguageCode: languageCode,
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

// CmdQueryCulturalQuote queries a cultural quote
func CmdQueryCulturalQuote() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cultural-quote [user-did] [category]",
		Short: "Query a cultural quote in user's language",
		Long: `Query a culturally relevant quote in the user's preferred language.

Categories:
- wisdom: Wisdom quotes from Indian philosophers and leaders
- motivation: Motivational quotes for inspiration
- patriotism: Quotes about patriotism and national pride
- technology: Quotes about technology and innovation

Example:
$ deshchaind query identity i18n cultural-quote did:desh:user123 "wisdom"
$ deshchaind query identity i18n cultural-quote did:desh:user456 "motivation"
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			userDID := args[0]
			category := args[1]

			res, err := queryClient.CulturalQuote(context.Background(), &types.QueryCulturalQuoteRequest{
				UserDid:  userDID,
				Category: category,
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

// CmdQueryGreeting queries a culturally appropriate greeting
func CmdQueryGreeting() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "greeting [user-did] [time-of-day] [festival]",
		Short: "Query a culturally appropriate greeting",
		Long: `Query a greeting message in the user's preferred language, appropriate for the time of day and any ongoing festival.

Time of day options: morning, afternoon, evening, night
Festival options: diwali, holi, dussehra, ganesh_chaturthi, eid, christmas, etc.

Example:
$ deshchaind query identity i18n greeting did:desh:user123 "morning" ""
$ deshchaind query identity i18n greeting did:desh:user456 "evening" "diwali"
`,
		Args: cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			userDID := args[0]
			timeOfDay := args[1]
			festival := ""
			if len(args) > 2 {
				festival = args[2]
			}

			res, err := queryClient.CulturalGreeting(context.Background(), &types.QueryCulturalGreetingRequest{
				UserDid:   userDID,
				TimeOfDay: timeOfDay,
				Festival:  festival,
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

// CmdQueryLocalizationConfig queries the current localization configuration
func CmdQueryLocalizationConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Query localization configuration",
		Long: `Query the current localization configuration for the identity module.

Example:
$ deshchaind query identity i18n config
`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.LocalizationConfig(context.Background(), &types.QueryLocalizationConfigRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// CmdQueryCustomMessages queries custom messages
func CmdQueryCustomMessages() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "custom-messages [category]",
		Short: "Query custom messages",
		Long: `Query custom messages, optionally filtered by category.

Example:
$ deshchaind query identity i18n custom-messages
$ deshchaind query identity i18n custom-messages "greetings"
`,
		Args: cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			category := ""
			if len(args) > 0 {
				category = args[0]
			}

			res, err := queryClient.CustomMessages(context.Background(), &types.QueryCustomMessagesRequest{
				Category: category,
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

// CmdQueryLocalizationStats queries localization statistics
func CmdQueryLocalizationStats() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Query localization statistics",
		Long: `Query statistics about localization usage and coverage.

Example:
$ deshchaind query identity i18n stats
`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.LocalizationStats(context.Background(), &types.QueryLocalizationStatsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// Helper function to display language information in a user-friendly format
func displayLanguageInfo(info *types.LanguageInfo) {
	fmt.Printf("Language Information:\n")
	fmt.Printf("  Code: %s\n", info.Code)
	fmt.Printf("  Name: %s\n", info.Name)
	fmt.Printf("  Native Name: %s\n", info.NativeName)
	fmt.Printf("  Script: %s\n", info.Script)
	fmt.Printf("  Direction: %s\n", info.Direction)
	fmt.Printf("  Region: %s\n", info.Region)
	fmt.Printf("  Official Language: %t\n", info.IsOfficial)
	fmt.Printf("  Supported: %t\n", info.IsSupported)
}