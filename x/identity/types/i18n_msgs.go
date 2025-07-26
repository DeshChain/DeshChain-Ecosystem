package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Message constants for i18n
const (
	TypeMsgSetUserLanguage          = "set_user_language"
	TypeMsgAddCustomMessage         = "add_custom_message"
	TypeMsgUpdateLocalizationConfig = "update_localization_config"
	TypeMsgImportMessages           = "import_messages"
)

var (
	_ sdk.Msg = &MsgSetUserLanguage{}
	_ sdk.Msg = &MsgAddCustomMessage{}
	_ sdk.Msg = &MsgUpdateLocalizationConfig{}
	_ sdk.Msg = &MsgImportMessages{}
)

// MsgSetUserLanguage sets user's preferred language
type MsgSetUserLanguage struct {
	UserDid      string `json:"user_did" yaml:"user_did"`
	LanguageCode string `json:"language_code" yaml:"language_code"`
	Signer       string `json:"signer" yaml:"signer"`
}

// NewMsgSetUserLanguage creates a new MsgSetUserLanguage instance
func NewMsgSetUserLanguage(userDid, languageCode, signer string) *MsgSetUserLanguage {
	return &MsgSetUserLanguage{
		UserDid:      userDid,
		LanguageCode: languageCode,
		Signer:       signer,
	}
}

// Route returns the module route
func (msg *MsgSetUserLanguage) Route() string {
	return RouterKey
}

// Type returns the message type
func (msg *MsgSetUserLanguage) Type() string {
	return TypeMsgSetUserLanguage
}

// GetSigners returns the signers
func (msg *MsgSetUserLanguage) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

// GetSignBytes returns the sign bytes
func (msg *MsgSetUserLanguage) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validates the message
func (msg *MsgSetUserLanguage) ValidateBasic() error {
	if msg.UserDid == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "user DID cannot be empty")
	}

	if msg.LanguageCode == "" {
		return sdkerrors.Wrap(ErrInvalidLanguageCode, "language code cannot be empty")
	}

	// Validate language code is supported
	supportedLanguages := GetSupportedLanguages()
	isSupported := false
	for _, lang := range supportedLanguages {
		if string(lang) == msg.LanguageCode {
			isSupported = true
			break
		}
	}

	if !isSupported {
		return sdkerrors.Wrapf(ErrUnsupportedLanguage, "language code not supported: %s", msg.LanguageCode)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid signer address: %v", err)
	}

	return nil
}

// MsgAddCustomMessage adds a custom localized message
type MsgAddCustomMessage struct {
	Key         string         `json:"key" yaml:"key"`
	Category    string         `json:"category" yaml:"category"`
	Text        *LocalizedText `json:"text" yaml:"text"`
	Description string         `json:"description" yaml:"description"`
	Signer      string         `json:"signer" yaml:"signer"`
}

// NewMsgAddCustomMessage creates a new MsgAddCustomMessage instance
func NewMsgAddCustomMessage(key, category, description, signer string, text *LocalizedText) *MsgAddCustomMessage {
	return &MsgAddCustomMessage{
		Key:         key,
		Category:    category,
		Text:        text,
		Description: description,
		Signer:      signer,
	}
}

// Route returns the module route
func (msg *MsgAddCustomMessage) Route() string {
	return RouterKey
}

// Type returns the message type
func (msg *MsgAddCustomMessage) Type() string {
	return TypeMsgAddCustomMessage
}

// GetSigners returns the signers
func (msg *MsgAddCustomMessage) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

// GetSignBytes returns the sign bytes
func (msg *MsgAddCustomMessage) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validates the message
func (msg *MsgAddCustomMessage) ValidateBasic() error {
	if msg.Key == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "message key cannot be empty")
	}

	if msg.Category == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "category cannot be empty")
	}

	if msg.Text == nil {
		return sdkerrors.Wrap(ErrInvalidRequest, "localized text cannot be nil")
	}

	if len(msg.Text.Translations) == 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "at least one translation is required")
	}

	// Validate that all language codes in translations are supported
	supportedLanguages := GetSupportedLanguages()
	supportedMap := make(map[string]bool)
	for _, lang := range supportedLanguages {
		supportedMap[string(lang)] = true
	}

	for langCode := range msg.Text.Translations {
		if !supportedMap[langCode] {
			return sdkerrors.Wrapf(ErrUnsupportedLanguage, "unsupported language in translations: %s", langCode)
		}
	}

	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid signer address: %v", err)
	}

	return nil
}

// MsgUpdateLocalizationConfig updates the localization configuration
type MsgUpdateLocalizationConfig struct {
	Config *LocalizationConfig `json:"config" yaml:"config"`
	Signer string              `json:"signer" yaml:"signer"`
}

// NewMsgUpdateLocalizationConfig creates a new MsgUpdateLocalizationConfig instance
func NewMsgUpdateLocalizationConfig(config *LocalizationConfig, signer string) *MsgUpdateLocalizationConfig {
	return &MsgUpdateLocalizationConfig{
		Config: config,
		Signer: signer,
	}
}

// Route returns the module route
func (msg *MsgUpdateLocalizationConfig) Route() string {
	return RouterKey
}

// Type returns the message type
func (msg *MsgUpdateLocalizationConfig) Type() string {
	return TypeMsgUpdateLocalizationConfig
}

// GetSigners returns the signers
func (msg *MsgUpdateLocalizationConfig) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

// GetSignBytes returns the sign bytes
func (msg *MsgUpdateLocalizationConfig) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validates the message
func (msg *MsgUpdateLocalizationConfig) ValidateBasic() error {
	if msg.Config == nil {
		return sdkerrors.Wrap(ErrInvalidRequest, "config cannot be nil")
	}

	if err := msg.Config.Validate(); err != nil {
		return sdkerrors.Wrap(ErrLocalizationConfigError, err.Error())
	}

	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid signer address: %v", err)
	}

	return nil
}

// MsgImportMessages imports messages from a catalog
type MsgImportMessages struct {
	Catalog *MessageCatalog `json:"catalog" yaml:"catalog"`
	Signer  string          `json:"signer" yaml:"signer"`
}

// NewMsgImportMessages creates a new MsgImportMessages instance
func NewMsgImportMessages(catalog *MessageCatalog, signer string) *MsgImportMessages {
	return &MsgImportMessages{
		Catalog: catalog,
		Signer:  signer,
	}
}

// Route returns the module route
func (msg *MsgImportMessages) Route() string {
	return RouterKey
}

// Type returns the message type
func (msg *MsgImportMessages) Type() string {
	return TypeMsgImportMessages
}

// GetSigners returns the signers
func (msg *MsgImportMessages) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

// GetSignBytes returns the sign bytes
func (msg *MsgImportMessages) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validates the message
func (msg *MsgImportMessages) ValidateBasic() error {
	if msg.Catalog == nil {
		return sdkerrors.Wrap(ErrInvalidRequest, "catalog cannot be nil")
	}

	if len(msg.Catalog.Messages) == 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "catalog must contain at least one message")
	}

	// Validate each message in the catalog
	for key, message := range msg.Catalog.Messages {
		if message.Key != key {
			return sdkerrors.Wrapf(ErrInvalidRequest, "message key mismatch: expected %s, got %s", key, message.Key)
		}

		if message.Text == nil {
			return sdkerrors.Wrapf(ErrInvalidRequest, "message text cannot be nil for key: %s", key)
		}

		if len(message.Text.Translations) == 0 {
			return sdkerrors.Wrapf(ErrInvalidRequest, "message must have at least one translation for key: %s", key)
		}
	}

	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid signer address: %v", err)
	}

	return nil
}

// Helper functions for message validation

// ValidateLanguageCode validates if a language code is supported
func ValidateLanguageCode(langCode string) error {
	supportedLanguages := GetSupportedLanguages()
	for _, lang := range supportedLanguages {
		if string(lang) == langCode {
			return nil
		}
	}
	return fmt.Errorf("unsupported language code: %s", langCode)
}

// ValidateUserDID validates DID format
func ValidateUserDID(did string) error {
	if did == "" {
		return fmt.Errorf("DID cannot be empty")
	}

	// Basic DID format validation
	if len(did) < 10 {
		return fmt.Errorf("DID too short: %s", did)
	}

	// Check if it starts with "did:"
	if len(did) < 4 || did[:4] != "did:" {
		return fmt.Errorf("DID must start with 'did:': %s", did)
	}

	return nil
}