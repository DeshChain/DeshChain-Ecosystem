package types

// Query request and response types for internationalization

// QuerySupportedLanguagesRequest is the request for supported languages
type QuerySupportedLanguagesRequest struct{}

// QuerySupportedLanguagesResponse is the response for supported languages
type QuerySupportedLanguagesResponse struct {
	Languages []*LanguageInfo `json:"languages" yaml:"languages"`
}

// QueryUserLanguageRequest is the request for user language preference
type QueryUserLanguageRequest struct {
	UserDid string `json:"user_did" yaml:"user_did"`
}

// QueryUserLanguageResponse is the response for user language preference
type QueryUserLanguageResponse struct {
	UserDid      string       `json:"user_did" yaml:"user_did"`
	LanguageCode LanguageCode `json:"language_code" yaml:"language_code"`
	LanguageInfo *LanguageInfo `json:"language_info" yaml:"language_info"`
}

// QueryLocalizedMessageRequest is the request for a localized message
type QueryLocalizedMessageRequest struct {
	Key          string `json:"key" yaml:"key"`
	LanguageCode string `json:"language_code" yaml:"language_code"`
}

// QueryLocalizedMessageResponse is the response for a localized message
type QueryLocalizedMessageResponse struct {
	Key          string `json:"key" yaml:"key"`
	LanguageCode string `json:"language_code" yaml:"language_code"`
	Message      string `json:"message" yaml:"message"`
	Found        bool   `json:"found" yaml:"found"`
}

// QueryRegionalLanguagesRequest is the request for regional languages
type QueryRegionalLanguagesRequest struct {
	Region string `json:"region" yaml:"region"`
}

// QueryRegionalLanguagesResponse is the response for regional languages
type QueryRegionalLanguagesResponse struct {
	Region    string          `json:"region" yaml:"region"`
	Languages []*LanguageInfo `json:"languages" yaml:"languages"`
}

// QueryLanguageInfoRequest is the request for language information
type QueryLanguageInfoRequest struct {
	LanguageCode string `json:"language_code" yaml:"language_code"`
}

// QueryLanguageInfoResponse is the response for language information
type QueryLanguageInfoResponse struct {
	LanguageInfo *LanguageInfo `json:"language_info" yaml:"language_info"`
}

// QueryCulturalQuoteRequest is the request for a cultural quote
type QueryCulturalQuoteRequest struct {
	UserDid  string `json:"user_did" yaml:"user_did"`
	Category string `json:"category" yaml:"category"`
}

// QueryCulturalQuoteResponse is the response for a cultural quote
type QueryCulturalQuoteResponse struct {
	Quote        string       `json:"quote" yaml:"quote"`
	Category     string       `json:"category" yaml:"category"`
	LanguageCode LanguageCode `json:"language_code" yaml:"language_code"`
	Author       string       `json:"author,omitempty" yaml:"author,omitempty"`
}

// QueryCulturalGreetingRequest is the request for a cultural greeting
type QueryCulturalGreetingRequest struct {
	UserDid   string `json:"user_did" yaml:"user_did"`
	TimeOfDay string `json:"time_of_day" yaml:"time_of_day"`
	Festival  string `json:"festival,omitempty" yaml:"festival,omitempty"`
}

// QueryCulturalGreetingResponse is the response for a cultural greeting
type QueryCulturalGreetingResponse struct {
	Greeting     string       `json:"greeting" yaml:"greeting"`
	LanguageCode LanguageCode `json:"language_code" yaml:"language_code"`
	TimeOfDay    string       `json:"time_of_day" yaml:"time_of_day"`
	Festival     string       `json:"festival,omitempty" yaml:"festival,omitempty"`
	Cultural     bool         `json:"cultural" yaml:"cultural"`
}

// QueryLocalizationConfigRequest is the request for localization configuration
type QueryLocalizationConfigRequest struct{}

// QueryLocalizationConfigResponse is the response for localization configuration
type QueryLocalizationConfigResponse struct {
	Config *LocalizationConfig `json:"config" yaml:"config"`
}

// QueryCustomMessagesRequest is the request for custom messages
type QueryCustomMessagesRequest struct {
	Category string `json:"category,omitempty" yaml:"category,omitempty"`
}

// QueryCustomMessagesResponse is the response for custom messages
type QueryCustomMessagesResponse struct {
	Messages []*IdentityMessage `json:"messages" yaml:"messages"`
	Total    int                `json:"total" yaml:"total"`
	Category string             `json:"category,omitempty" yaml:"category,omitempty"`
}

// QueryLocalizationStatsRequest is the request for localization statistics
type QueryLocalizationStatsRequest struct{}

// QueryLocalizationStatsResponse is the response for localization statistics
type QueryLocalizationStatsResponse struct {
	TotalSupportedLanguages int                    `json:"total_supported_languages" yaml:"total_supported_languages"`
	TotalMessages          int                    `json:"total_messages" yaml:"total_messages"`
	TotalCustomMessages    int                    `json:"total_custom_messages" yaml:"total_custom_messages"`
	MostUsedLanguage       LanguageCode           `json:"most_used_language" yaml:"most_used_language"`
	CoveragePercentage     float64                `json:"coverage_percentage" yaml:"coverage_percentage"`
	LanguageDistribution   map[string]int         `json:"language_distribution" yaml:"language_distribution"`
	RegionalStats          map[string]interface{} `json:"regional_stats" yaml:"regional_stats"`
}

// Validation functions for query types

// Validate validates QueryUserLanguageRequest
func (q *QueryUserLanguageRequest) Validate() error {
	return ValidateUserDID(q.UserDid)
}

// Validate validates QueryLocalizedMessageRequest
func (q *QueryLocalizedMessageRequest) Validate() error {
	if q.Key == "" {
		return ErrInvalidRequest
	}
	return ValidateLanguageCode(q.LanguageCode)
}

// Validate validates QueryRegionalLanguagesRequest
func (q *QueryRegionalLanguagesRequest) Validate() error {
	if q.Region == "" {
		return ErrInvalidRequest
	}
	return nil
}

// Validate validates QueryLanguageInfoRequest
func (q *QueryLanguageInfoRequest) Validate() error {
	return ValidateLanguageCode(q.LanguageCode)
}

// Validate validates QueryCulturalQuoteRequest
func (q *QueryCulturalQuoteRequest) Validate() error {
	if err := ValidateUserDID(q.UserDid); err != nil {
		return err
	}
	
	// Validate category
	validCategories := []string{"wisdom", "motivation", "patriotism", "technology", "general"}
	if q.Category != "" {
		isValid := false
		for _, cat := range validCategories {
			if q.Category == cat {
				isValid = true
				break
			}
		}
		if !isValid {
			return ErrInvalidRequest
		}
	}
	
	return nil
}

// Validate validates QueryCulturalGreetingRequest
func (q *QueryCulturalGreetingRequest) Validate() error {
	if err := ValidateUserDID(q.UserDid); err != nil {
		return err
	}
	
	// Validate time of day
	if q.TimeOfDay != "" {
		validTimes := []string{"morning", "afternoon", "evening", "night"}
		isValid := false
		for _, time := range validTimes {
			if q.TimeOfDay == time {
				isValid = true
				break
			}
		}
		if !isValid {
			return ErrInvalidRequest
		}
	}
	
	return nil
}

// Validate validates QueryCustomMessagesRequest
func (q *QueryCustomMessagesRequest) Validate() error {
	// Category is optional, no validation needed if empty
	return nil
}

// Helper functions for creating query responses

// NewQuerySupportedLanguagesResponse creates a new supported languages response
func NewQuerySupportedLanguagesResponse(languages []*LanguageInfo) *QuerySupportedLanguagesResponse {
	return &QuerySupportedLanguagesResponse{
		Languages: languages,
	}
}

// NewQueryUserLanguageResponse creates a new user language response
func NewQueryUserLanguageResponse(userDid string, langCode LanguageCode, langInfo *LanguageInfo) *QueryUserLanguageResponse {
	return &QueryUserLanguageResponse{
		UserDid:      userDid,
		LanguageCode: langCode,
		LanguageInfo: langInfo,
	}
}

// NewQueryLocalizedMessageResponse creates a new localized message response
func NewQueryLocalizedMessageResponse(key, langCode, message string, found bool) *QueryLocalizedMessageResponse {
	return &QueryLocalizedMessageResponse{
		Key:          key,
		LanguageCode: langCode,
		Message:      message,
		Found:        found,
	}
}

// NewQueryRegionalLanguagesResponse creates a new regional languages response
func NewQueryRegionalLanguagesResponse(region string, languages []*LanguageInfo) *QueryRegionalLanguagesResponse {
	return &QueryRegionalLanguagesResponse{
		Region:    region,
		Languages: languages,
	}
}

// NewQueryLanguageInfoResponse creates a new language info response
func NewQueryLanguageInfoResponse(langInfo *LanguageInfo) *QueryLanguageInfoResponse {
	return &QueryLanguageInfoResponse{
		LanguageInfo: langInfo,
	}
}

// NewQueryCulturalQuoteResponse creates a new cultural quote response
func NewQueryCulturalQuoteResponse(quote, category string, langCode LanguageCode, author string) *QueryCulturalQuoteResponse {
	return &QueryCulturalQuoteResponse{
		Quote:        quote,
		Category:     category,
		LanguageCode: langCode,
		Author:       author,
	}
}

// NewQueryCulturalGreetingResponse creates a new cultural greeting response
func NewQueryCulturalGreetingResponse(greeting string, langCode LanguageCode, timeOfDay, festival string, cultural bool) *QueryCulturalGreetingResponse {
	return &QueryCulturalGreetingResponse{
		Greeting:     greeting,
		LanguageCode: langCode,
		TimeOfDay:    timeOfDay,
		Festival:     festival,
		Cultural:     cultural,
	}
}

// NewQueryLocalizationConfigResponse creates a new localization config response
func NewQueryLocalizationConfigResponse(config *LocalizationConfig) *QueryLocalizationConfigResponse {
	return &QueryLocalizationConfigResponse{
		Config: config,
	}
}

// NewQueryCustomMessagesResponse creates a new custom messages response
func NewQueryCustomMessagesResponse(messages []*IdentityMessage, category string) *QueryCustomMessagesResponse {
	return &QueryCustomMessagesResponse{
		Messages: messages,
		Total:    len(messages),
		Category: category,
	}
}

// NewQueryLocalizationStatsResponse creates a new localization stats response
func NewQueryLocalizationStatsResponse(totalLangs, totalMsgs, totalCustomMsgs int, mostUsed LanguageCode, coverage float64) *QueryLocalizationStatsResponse {
	return &QueryLocalizationStatsResponse{
		TotalSupportedLanguages: totalLangs,
		TotalMessages:          totalMsgs,
		TotalCustomMessages:    totalCustomMsgs,
		MostUsedLanguage:       mostUsed,
		CoveragePercentage:     coverage,
		LanguageDistribution:   make(map[string]int),
		RegionalStats:          make(map[string]interface{}),
	}
}