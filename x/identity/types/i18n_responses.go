package types

// Response types for internationalization messages

// MsgSetUserLanguageResponse is the response for setting user language
type MsgSetUserLanguageResponse struct {
	Success bool   `json:"success" yaml:"success"`
	Message string `json:"message" yaml:"message"`
}

// MsgAddCustomMessageResponse is the response for adding custom message
type MsgAddCustomMessageResponse struct {
	Success bool   `json:"success" yaml:"success"`
	Message string `json:"message" yaml:"message"`
}

// MsgUpdateLocalizationConfigResponse is the response for updating localization config
type MsgUpdateLocalizationConfigResponse struct {
	Success bool   `json:"success" yaml:"success"`
	Message string `json:"message" yaml:"message"`
}

// MsgImportMessagesResponse is the response for importing messages
type MsgImportMessagesResponse struct {
	Success       bool   `json:"success" yaml:"success"`
	Message       string `json:"message" yaml:"message"`
	ImportedCount int32  `json:"imported_count" yaml:"imported_count"`
}

// Constructor functions for responses

// NewMsgSetUserLanguageResponse creates a new MsgSetUserLanguageResponse
func NewMsgSetUserLanguageResponse(success bool, message string) *MsgSetUserLanguageResponse {
	return &MsgSetUserLanguageResponse{
		Success: success,
		Message: message,
	}
}

// NewMsgAddCustomMessageResponse creates a new MsgAddCustomMessageResponse
func NewMsgAddCustomMessageResponse(success bool, message string) *MsgAddCustomMessageResponse {
	return &MsgAddCustomMessageResponse{
		Success: success,
		Message: message,
	}
}

// NewMsgUpdateLocalizationConfigResponse creates a new MsgUpdateLocalizationConfigResponse
func NewMsgUpdateLocalizationConfigResponse(success bool, message string) *MsgUpdateLocalizationConfigResponse {
	return &MsgUpdateLocalizationConfigResponse{
		Success: success,
		Message: message,
	}
}

// NewMsgImportMessagesResponse creates a new MsgImportMessagesResponse
func NewMsgImportMessagesResponse(success bool, message string, importedCount int32) *MsgImportMessagesResponse {
	return &MsgImportMessagesResponse{
		Success:       success,
		Message:       message,
		ImportedCount: importedCount,
	}
}