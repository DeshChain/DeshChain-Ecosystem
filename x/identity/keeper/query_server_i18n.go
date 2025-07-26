package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/deshchain/deshchain/x/identity/types"
)

// SupportedLanguages returns all supported languages
func (k Keeper) SupportedLanguages(goCtx context.Context, req *types.QuerySupportedLanguagesRequest) (*types.QuerySupportedLanguagesResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get supported languages from the i18n keeper
	languages := k.i18nKeeper.GetSupportedLanguages(ctx)

	return types.NewQuerySupportedLanguagesResponse(languages), nil
}

// UserLanguage returns user's language preference
func (k Keeper) UserLanguage(goCtx context.Context, req *types.QueryUserLanguageRequest) (*types.QueryUserLanguageResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidRequest, "empty request")
	}

	if err := req.Validate(); err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidRequest, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get user's language preference
	langCode := k.i18nKeeper.getUserLanguagePreference(ctx, req.UserDid)
	if langCode == "" {
		// Return default language if no preference is set
		config := k.i18nKeeper.GetLocalizationConfig(ctx)
		langCode = config.DefaultLanguage
	}

	// Get language info
	langInfo := types.GetLanguageInfo(langCode)

	return types.NewQueryUserLanguageResponse(req.UserDid, langCode, langInfo), nil
}

// LocalizedMessage returns a localized message
func (k Keeper) LocalizedMessage(goCtx context.Context, req *types.QueryLocalizedMessageRequest) (*types.QueryLocalizedMessageResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidRequest, "empty request")
	}

	if err := req.Validate(); err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidRequest, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the localized message
	langCode := types.LanguageCode(req.LanguageCode)
	message := k.i18nKeeper.GetLocalizedMessage(ctx, req.Key, langCode)
	
	found := message != ""
	if !found {
		// Try to get from built-in messages
		builtInMessages := types.GetBuiltInMessages()
		if builtInMsg, exists := builtInMessages[req.Key]; exists {
			message = builtInMsg.Text.GetText(langCode)
			found = message != ""
		}
	}

	return types.NewQueryLocalizedMessageResponse(req.Key, req.LanguageCode, message, found), nil
}

// RegionalLanguages returns languages for a specific region
func (k Keeper) RegionalLanguages(goCtx context.Context, req *types.QueryRegionalLanguagesRequest) (*types.QueryRegionalLanguagesResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidRequest, "empty request")
	}

	if err := req.Validate(); err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidRequest, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get regional languages
	languages := k.i18nKeeper.GetRegionalLanguages(ctx, req.Region)

	return types.NewQueryRegionalLanguagesResponse(req.Region, languages), nil
}

// LanguageInfo returns detailed information about a language
func (k Keeper) LanguageInfo(goCtx context.Context, req *types.QueryLanguageInfoRequest) (*types.QueryLanguageInfoResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidRequest, "empty request")
	}

	if err := req.Validate(); err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidRequest, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get language info
	langCode := types.LanguageCode(req.LanguageCode)
	langInfo := types.GetLanguageInfo(langCode)

	if langInfo == nil {
		return nil, sdkerrors.Wrapf(types.ErrUnsupportedLanguage, "language not found: %s", req.LanguageCode)
	}

	return types.NewQueryLanguageInfoResponse(langInfo), nil
}

// CulturalQuote returns a cultural quote in user's language
func (k Keeper) CulturalQuote(goCtx context.Context, req *types.QueryCulturalQuoteRequest) (*types.QueryCulturalQuoteResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidRequest, "empty request")
	}

	if err := req.Validate(); err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidRequest, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get cultural quote
	quote := k.i18nKeeper.GetCulturalQuote(ctx, req.UserDid, req.Category)
	
	// Get user's language preference
	langCode := k.i18nKeeper.getUserLanguagePreference(ctx, req.UserDid)
	if langCode == "" {
		config := k.i18nKeeper.GetLocalizationConfig(ctx)
		langCode = config.DefaultLanguage
	}

	// For now, we don't extract author from quotes, but this could be enhanced
	author := ""

	return types.NewQueryCulturalQuoteResponse(quote, req.Category, langCode, author), nil
}

// CulturalGreeting returns a culturally appropriate greeting
func (k Keeper) CulturalGreeting(goCtx context.Context, req *types.QueryCulturalGreetingRequest) (*types.QueryCulturalGreetingResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidRequest, "empty request")
	}

	if err := req.Validate(); err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidRequest, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get cultural greeting
	greeting := k.i18nKeeper.GenerateCulturalGreeting(ctx, req.UserDid, req.TimeOfDay, req.Festival)
	
	// Get user's language preference
	langCode := k.i18nKeeper.getUserLanguagePreference(ctx, req.UserDid)
	if langCode == "" {
		config := k.i18nKeeper.GetLocalizationConfig(ctx)
		langCode = config.DefaultLanguage
	}

	// Determine if this is a cultural greeting (non-English or festival-specific)
	cultural := langCode != types.LanguageEnglish || req.Festival != ""

	return types.NewQueryCulturalGreetingResponse(greeting, langCode, req.TimeOfDay, req.Festival, cultural), nil
}

// LocalizationConfig returns the current localization configuration
func (k Keeper) LocalizationConfig(goCtx context.Context, req *types.QueryLocalizationConfigRequest) (*types.QueryLocalizationConfigResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get localization configuration
	config := k.i18nKeeper.GetLocalizationConfig(ctx)

	return types.NewQueryLocalizationConfigResponse(config), nil
}

// CustomMessages returns custom messages, optionally filtered by category
func (k Keeper) CustomMessages(goCtx context.Context, req *types.QueryCustomMessagesRequest) (*types.QueryCustomMessagesResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidRequest, "empty request")
	}

	if err := req.Validate(); err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidRequest, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get all custom messages
	allMessages, err := k.i18nKeeper.GetAllCustomMessages(ctx)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to get custom messages")
	}

	// Filter by category if specified
	var filteredMessages []*types.IdentityMessage
	if req.Category != "" {
		for _, msg := range allMessages {
			if msg.Category == req.Category {
				filteredMessages = append(filteredMessages, msg)
			}
		}
	} else {
		filteredMessages = allMessages
	}

	return types.NewQueryCustomMessagesResponse(filteredMessages, req.Category), nil
}

// LocalizationStats returns localization statistics
func (k Keeper) LocalizationStats(goCtx context.Context, req *types.QueryLocalizationStatsRequest) (*types.QueryLocalizationStatsResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidRequest, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get localization statistics
	stats := k.i18nKeeper.GetLocalizationStats(ctx)

	// Convert to proper types
	totalLangs := stats["total_supported_languages"].(int)
	totalMsgs := stats["total_messages"].(int)
	totalCustomMsgs := stats["custom_messages"].(int)
	mostUsed := stats["most_used_language"].(types.LanguageCode)
	coverage := stats["coverage_percentage"].(float64)

	response := types.NewQueryLocalizationStatsResponse(totalLangs, totalMsgs, totalCustomMsgs, mostUsed, coverage)

	// Add language distribution (mock data for now)
	response.LanguageDistribution = map[string]int{
		"hi": 35, // Hindi
		"en": 25, // English
		"bn": 10, // Bengali
		"ta": 8,  // Tamil
		"te": 6,  // Telugu
		"mr": 5,  // Marathi
		"gu": 4,  // Gujarati
		"others": 7,
	}

	// Add regional statistics (mock data for now)
	response.RegionalStats = map[string]interface{}{
		"north_india": map[string]interface{}{
			"primary_languages": []string{"hi", "ur", "pa"},
			"usage_percentage": 40.0,
		},
		"south_india": map[string]interface{}{
			"primary_languages": []string{"ta", "te", "kn", "ml"},
			"usage_percentage": 30.0,
		},
		"west_india": map[string]interface{}{
			"primary_languages": []string{"mr", "gu"},
			"usage_percentage": 15.0,
		},
		"east_india": map[string]interface{}{
			"primary_languages": []string{"bn", "or", "as"},
			"usage_percentage": 15.0,
		},
	}

	return response, nil
}