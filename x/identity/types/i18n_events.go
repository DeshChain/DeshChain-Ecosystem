package types

// Event types for internationalization operations
const (
	// Event types
	EventTypeUserLanguageSet          = "user_language_set"
	EventTypeCustomMessageAdded       = "custom_message_added"
	EventTypeLocalizationConfigUpdated = "localization_config_updated"
	EventTypeMessagesImported         = "messages_imported"
	EventTypeLanguageDetected         = "language_detected"
	EventTypeCulturalGreetingGenerated = "cultural_greeting_generated"
	EventTypeCulturalQuoteGenerated   = "cultural_quote_generated"
	EventTypeTranslationRequested     = "translation_requested"
	EventTypeLocalizationStatsUpdated = "localization_stats_updated"

	// Attribute keys
	AttributeKeyUserDID           = "user_did"
	AttributeKeyLanguageCode      = "language_code"
	AttributeKeyMessageKey        = "message_key"
	AttributeKeyCategory          = "category"
	AttributeKeySigner            = "signer"
	AttributeKeyDefaultLanguage   = "default_language"
	AttributeKeyRegion            = "region"
	AttributeKeyMessageCount      = "message_count"
	AttributeKeyTimeOfDay         = "time_of_day"
	AttributeKeyFestival          = "festival"
	AttributeKeyQuoteCategory     = "quote_category"
	AttributeKeyDetectedLanguage  = "detected_language"
	AttributeKeyTranslationKey    = "translation_key"
	AttributeKeySourceLanguage    = "source_language"
	AttributeKeyTargetLanguage    = "target_language"
	AttributeKeySuccess           = "success"
	AttributeKeyErrorMessage      = "error_message"
	AttributeKeyGreetingType      = "greeting_type"
	AttributeKeyIsCultural        = "is_cultural"
	AttributeKeyIsTimeBasedGreeting = "is_time_based"
	AttributeKeyIsFestivalGreeting = "is_festival_greeting"

	// Module attribute values
	AttributeValueCategory        = ModuleName
)

// Event helper functions

// NewUserLanguageSetEvent creates a new user language set event
func NewUserLanguageSetEvent(userDID, languageCode, signer string) map[string]string {
	return map[string]string{
		AttributeKeyUserDID:      userDID,
		AttributeKeyLanguageCode: languageCode,
		AttributeKeySigner:       signer,
	}
}

// NewCustomMessageAddedEvent creates a new custom message added event
func NewCustomMessageAddedEvent(messageKey, category, signer string) map[string]string {
	return map[string]string{
		AttributeKeyMessageKey: messageKey,
		AttributeKeyCategory:   category,
		AttributeKeySigner:     signer,
	}
}

// NewLocalizationConfigUpdatedEvent creates a new localization config updated event
func NewLocalizationConfigUpdatedEvent(defaultLanguage, region, signer string) map[string]string {
	return map[string]string{
		AttributeKeyDefaultLanguage: defaultLanguage,
		AttributeKeyRegion:          region,
		AttributeKeySigner:          signer,
	}
}

// NewMessagesImportedEvent creates a new messages imported event
func NewMessagesImportedEvent(messageCount, signer string) map[string]string {
	return map[string]string{
		AttributeKeyMessageCount: messageCount,
		AttributeKeySigner:       signer,
	}
}

// NewLanguageDetectedEvent creates a new language detected event
func NewLanguageDetectedEvent(userDID, detectedLanguage, sourceText string) map[string]string {
	return map[string]string{
		AttributeKeyUserDID:          userDID,
		AttributeKeyDetectedLanguage: detectedLanguage,
	}
}

// NewCulturalGreetingGeneratedEvent creates a new cultural greeting generated event
func NewCulturalGreetingGeneratedEvent(userDID, languageCode, timeOfDay, festival string, isCultural bool) map[string]string {
	greetingType := "standard"
	if festival != "" {
		greetingType = "festival"
	} else if timeOfDay != "" {
		greetingType = "time_based"
	}

	event := map[string]string{
		AttributeKeyUserDID:              userDID,
		AttributeKeyLanguageCode:         languageCode,
		AttributeKeyGreetingType:         greetingType,
		AttributeKeyIsCultural:           boolToString(isCultural),
	}

	if timeOfDay != "" {
		event[AttributeKeyTimeOfDay] = timeOfDay
		event[AttributeKeyIsTimeBasedGreeting] = "true"
	}

	if festival != "" {
		event[AttributeKeyFestival] = festival
		event[AttributeKeyIsFestivalGreeting] = "true"
	}

	return event
}

// NewCulturalQuoteGeneratedEvent creates a new cultural quote generated event
func NewCulturalQuoteGeneratedEvent(userDID, languageCode, category string) map[string]string {
	return map[string]string{
		AttributeKeyUserDID:       userDID,
		AttributeKeyLanguageCode:  languageCode,
		AttributeKeyQuoteCategory: category,
	}
}

// NewTranslationRequestedEvent creates a new translation requested event
func NewTranslationRequestedEvent(translationKey, sourceLanguage, targetLanguage string, success bool, errorMessage string) map[string]string {
	event := map[string]string{
		AttributeKeyTranslationKey: translationKey,
		AttributeKeySourceLanguage: sourceLanguage,
		AttributeKeyTargetLanguage: targetLanguage,
		AttributeKeySuccess:        boolToString(success),
	}

	if errorMessage != "" {
		event[AttributeKeyErrorMessage] = errorMessage
	}

	return event
}

// NewLocalizationStatsUpdatedEvent creates a new localization stats updated event
func NewLocalizationStatsUpdatedEvent(totalLanguages, totalMessages, customMessages int) map[string]string {
	return map[string]string{
		"total_languages":    intToString(totalLanguages),
		"total_messages":     intToString(totalMessages),
		"custom_messages":    intToString(customMessages),
	}
}

// Helper functions for event attributes

// boolToString converts bool to string
func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// intToString converts int to string
func intToString(i int) string {
	return string(rune(i + 48)) // Simple conversion for small numbers
}

// Event validation functions

// ValidateUserLanguageSetEvent validates user language set event attributes
func ValidateUserLanguageSetEvent(attrs map[string]string) error {
	if attrs[AttributeKeyUserDID] == "" {
		return ErrInvalidRequest
	}
	if attrs[AttributeKeyLanguageCode] == "" {
		return ErrInvalidLanguageCode
	}
	if attrs[AttributeKeySigner] == "" {
		return ErrInvalidAddress
	}
	return nil
}

// ValidateCustomMessageAddedEvent validates custom message added event attributes
func ValidateCustomMessageAddedEvent(attrs map[string]string) error {
	if attrs[AttributeKeyMessageKey] == "" {
		return ErrInvalidRequest
	}
	if attrs[AttributeKeyCategory] == "" {
		return ErrInvalidRequest
	}
	if attrs[AttributeKeySigner] == "" {
		return ErrInvalidAddress
	}
	return nil
}

// ValidateLocalizationConfigUpdatedEvent validates localization config updated event attributes
func ValidateLocalizationConfigUpdatedEvent(attrs map[string]string) error {
	if attrs[AttributeKeyDefaultLanguage] == "" {
		return ErrInvalidLanguageCode
	}
	if attrs[AttributeKeySigner] == "" {
		return ErrInvalidAddress
	}
	return nil
}

// ValidateMessagesImportedEvent validates messages imported event attributes
func ValidateMessagesImportedEvent(attrs map[string]string) error {
	if attrs[AttributeKeyMessageCount] == "" {
		return ErrInvalidRequest
	}
	if attrs[AttributeKeySigner] == "" {
		return ErrInvalidAddress
	}
	return nil
}

// Event type constants for easier reference
var (
	// All internationalization event types
	I18nEventTypes = []string{
		EventTypeUserLanguageSet,
		EventTypeCustomMessageAdded,
		EventTypeLocalizationConfigUpdated,
		EventTypeMessagesImported,
		EventTypeLanguageDetected,
		EventTypeCulturalGreetingGenerated,
		EventTypeCulturalQuoteGenerated,
		EventTypeTranslationRequested,
		EventTypeLocalizationStatsUpdated,
	}

	// All internationalization attribute keys
	I18nAttributeKeys = []string{
		AttributeKeyUserDID,
		AttributeKeyLanguageCode,
		AttributeKeyMessageKey,
		AttributeKeyCategory,
		AttributeKeySigner,
		AttributeKeyDefaultLanguage,
		AttributeKeyRegion,
		AttributeKeyMessageCount,
		AttributeKeyTimeOfDay,
		AttributeKeyFestival,
		AttributeKeyQuoteCategory,
		AttributeKeyDetectedLanguage,
		AttributeKeyTranslationKey,
		AttributeKeySourceLanguage,
		AttributeKeyTargetLanguage,
		AttributeKeySuccess,
		AttributeKeyErrorMessage,
		AttributeKeyGreetingType,
		AttributeKeyIsCultural,
		AttributeKeyIsTimeBasedGreeting,
		AttributeKeyIsFestivalGreeting,
	}
)

// GetI18nEventTypes returns all internationalization event types
func GetI18nEventTypes() []string {
	return I18nEventTypes
}

// GetI18nAttributeKeys returns all internationalization attribute keys
func GetI18nAttributeKeys() []string {
	return I18nAttributeKeys
}

// IsI18nEvent checks if an event type is an internationalization event
func IsI18nEvent(eventType string) bool {
	for _, i18nEventType := range I18nEventTypes {
		if eventType == i18nEventType {
			return true
		}
	}
	return false
}

// IsI18nAttributeKey checks if an attribute key is an internationalization attribute
func IsI18nAttributeKey(attributeKey string) bool {
	for _, i18nAttributeKey := range I18nAttributeKeys {
		if attributeKey == i18nAttributeKey {
			return true
		}
	}
	return false
}