package types

import (
	"fmt"
	"strings"
)

// LanguageCode represents supported language codes
type LanguageCode string

// Supported language codes for Indian languages and international standards
const (
	// Hindi and its variants
	LanguageHindi     LanguageCode = "hi"    // Hindi (Devanagari)
	LanguageHinglish  LanguageCode = "hi-en" // Hinglish (Hindi-English mix)
	
	// Official Indian languages
	LanguageAssamese   LanguageCode = "as"    // Assamese
	LanguageBengali    LanguageCode = "bn"    // Bengali
	LanguageBodo       LanguageCode = "brx"   // Bodo
	LanguageDogri      LanguageCode = "doi"   // Dogri
	LanguageGujarati   LanguageCode = "gu"    // Gujarati
	LanguageKannada    LanguageCode = "kn"    // Kannada
	LanguageKashmiri   LanguageCode = "ks"    // Kashmiri
	LanguageKonkani    LanguageCode = "gom"   // Konkani
	LanguageMaithili   LanguageCode = "mai"   // Maithili
	LanguageMalayalam  LanguageCode = "ml"    // Malayalam
	LanguageManipuri   LanguageCode = "mni"   // Manipuri (Meitei)
	LanguageMarathi    LanguageCode = "mr"    // Marathi
	LanguageNepali     LanguageCode = "ne"    // Nepali
	LanguageOdia       LanguageCode = "or"    // Odia
	LanguagePunjabi    LanguageCode = "pa"    // Punjabi
	LanguageSanskrit   LanguageCode = "sa"    // Sanskrit
	LanguageSantali    LanguageCode = "sat"   // Santali
	LanguageSindhi     LanguageCode = "sd"    // Sindhi
	LanguageTamil      LanguageCode = "ta"    // Tamil
	LanguageTelugu     LanguageCode = "te"    // Telugu
	LanguageUrdu       LanguageCode = "ur"    // Urdu
	
	// International languages
	LanguageEnglish    LanguageCode = "en"    // English
	LanguageArabic     LanguageCode = "ar"    // Arabic
	LanguageChinese    LanguageCode = "zh"    // Chinese
	LanguageSpanish    LanguageCode = "es"    // Spanish
	LanguageFrench     LanguageCode = "fr"    // French
	LanguageRussian    LanguageCode = "ru"    // Russian
	LanguageJapanese   LanguageCode = "ja"    // Japanese
	LanguageKorean     LanguageCode = "ko"    // Korean
	LanguageGerman     LanguageCode = "de"    // German
	LanguagePortuguese LanguageCode = "pt"    // Portuguese
)

// LocalizedText represents text in multiple languages
type LocalizedText struct {
	Translations map[LanguageCode]string `json:"translations"`
	DefaultLang  LanguageCode            `json:"default_language"`
}

// NewLocalizedText creates a new localized text with default language
func NewLocalizedText(defaultLang LanguageCode, defaultText string) *LocalizedText {
	return &LocalizedText{
		Translations: map[LanguageCode]string{
			defaultLang: defaultText,
		},
		DefaultLang: defaultLang,
	}
}

// AddTranslation adds a translation for a specific language
func (lt *LocalizedText) AddTranslation(lang LanguageCode, text string) {
	if lt.Translations == nil {
		lt.Translations = make(map[LanguageCode]string)
	}
	lt.Translations[lang] = text
}

// GetText returns text in the requested language, falling back to default
func (lt *LocalizedText) GetText(lang LanguageCode) string {
	if text, exists := lt.Translations[lang]; exists {
		return text
	}
	
	// Fallback to default language
	if text, exists := lt.Translations[lt.DefaultLang]; exists {
		return text
	}
	
	// Fallback to English if default is not available
	if text, exists := lt.Translations[LanguageEnglish]; exists {
		return text
	}
	
	// Return first available translation
	for _, text := range lt.Translations {
		return text
	}
	
	return ""
}

// HasTranslation checks if translation exists for a language
func (lt *LocalizedText) HasTranslation(lang LanguageCode) bool {
	_, exists := lt.Translations[lang]
	return exists
}

// GetAvailableLanguages returns list of available languages
func (lt *LocalizedText) GetAvailableLanguages() []LanguageCode {
	languages := make([]LanguageCode, 0, len(lt.Translations))
	for lang := range lt.Translations {
		languages = append(languages, lang)
	}
	return languages
}

// LanguageInfo contains metadata about a language
type LanguageInfo struct {
	Code         LanguageCode `json:"code"`
	Name         string       `json:"name"`
	NativeName   string       `json:"native_name"`
	Script       string       `json:"script"`
	Direction    string       `json:"direction"` // ltr or rtl
	Region       string       `json:"region"`
	IsOfficial   bool         `json:"is_official"`   // Official language of India
	IsSupported  bool         `json:"is_supported"`  // Supported in identity flows
}

// GetLanguageInfo returns information about a language
func GetLanguageInfo(code LanguageCode) *LanguageInfo {
	languageMap := map[LanguageCode]*LanguageInfo{
		LanguageHindi: {
			Code: LanguageHindi, Name: "Hindi", NativeName: "हिन्दी",
			Script: "Devanagari", Direction: "ltr", Region: "India", IsOfficial: true, IsSupported: true,
		},
		LanguageEnglish: {
			Code: LanguageEnglish, Name: "English", NativeName: "English",
			Script: "Latin", Direction: "ltr", Region: "Global", IsOfficial: true, IsSupported: true,
		},
		LanguageBengali: {
			Code: LanguageBengali, Name: "Bengali", NativeName: "বাংলা",
			Script: "Bengali", Direction: "ltr", Region: "West Bengal, Bangladesh", IsOfficial: true, IsSupported: true,
		},
		LanguageTelugu: {
			Code: LanguageTelugu, Name: "Telugu", NativeName: "తెలుగు",
			Script: "Telugu", Direction: "ltr", Region: "Andhra Pradesh, Telangana", IsOfficial: true, IsSupported: true,
		},
		LanguageMarathi: {
			Code: LanguageMarathi, Name: "Marathi", NativeName: "मराठी",
			Script: "Devanagari", Direction: "ltr", Region: "Maharashtra", IsOfficial: true, IsSupported: true,
		},
		LanguageTamil: {
			Code: LanguageTamil, Name: "Tamil", NativeName: "தமிழ்",
			Script: "Tamil", Direction: "ltr", Region: "Tamil Nadu, Sri Lanka", IsOfficial: true, IsSupported: true,
		},
		LanguageUrdu: {
			Code: LanguageUrdu, Name: "Urdu", NativeName: "اردو",
			Script: "Arabic", Direction: "rtl", Region: "India, Pakistan", IsOfficial: true, IsSupported: true,
		},
		LanguageGujarati: {
			Code: LanguageGujarati, Name: "Gujarati", NativeName: "ગુજરાતી",
			Script: "Gujarati", Direction: "ltr", Region: "Gujarat", IsOfficial: true, IsSupported: true,
		},
		LanguageKannada: {
			Code: LanguageKannada, Name: "Kannada", NativeName: "ಕನ್ನಡ",
			Script: "Kannada", Direction: "ltr", Region: "Karnataka", IsOfficial: true, IsSupported: true,
		},
		LanguageMalayalam: {
			Code: LanguageMalayalam, Name: "Malayalam", NativeName: "മലയാളം",
			Script: "Malayalam", Direction: "ltr", Region: "Kerala", IsOfficial: true, IsSupported: true,
		},
		LanguageOdia: {
			Code: LanguageOdia, Name: "Odia", NativeName: "ଓଡ଼ିଆ",
			Script: "Odia", Direction: "ltr", Region: "Odisha", IsOfficial: true, IsSupported: true,
		},
		LanguagePunjabi: {
			Code: LanguagePunjabi, Name: "Punjabi", NativeName: "ਪੰਜਾਬੀ",
			Script: "Gurmukhi", Direction: "ltr", Region: "Punjab", IsOfficial: true, IsSupported: true,
		},
		LanguageAssamese: {
			Code: LanguageAssamese, Name: "Assamese", NativeName: "অসমীয়া",
			Script: "Bengali", Direction: "ltr", Region: "Assam", IsOfficial: true, IsSupported: true,
		},
		LanguageSanskrit: {
			Code: LanguageSanskrit, Name: "Sanskrit", NativeName: "संस्कृतम्",
			Script: "Devanagari", Direction: "ltr", Region: "India", IsOfficial: true, IsSupported: true,
		},
		LanguageNepali: {
			Code: LanguageNepali, Name: "Nepali", NativeName: "नेपाली",
			Script: "Devanagari", Direction: "ltr", Region: "Nepal, India", IsOfficial: true, IsSupported: true,
		},
		LanguageKonkani: {
			Code: LanguageKonkani, Name: "Konkani", NativeName: "कोंकणी",
			Script: "Devanagari", Direction: "ltr", Region: "Goa", IsOfficial: true, IsSupported: true,
		},
		LanguageManipuri: {
			Code: LanguageManipuri, Name: "Manipuri", NativeName: "ꯃꯤꯇꯩ ꯂꯣꯟ",
			Script: "Meitei Mayek", Direction: "ltr", Region: "Manipur", IsOfficial: true, IsSupported: true,
		},
		LanguageBodo: {
			Code: LanguageBodo, Name: "Bodo", NativeName: "बड़ो",
			Script: "Devanagari", Direction: "ltr", Region: "Assam", IsOfficial: true, IsSupported: true,
		},
		LanguageSantali: {
			Code: LanguageSantali, Name: "Santali", NativeName: "ᱥᱟᱱᱛᱟᱲᱤ",
			Script: "Ol Chiki", Direction: "ltr", Region: "Jharkhand", IsOfficial: true, IsSupported: true,
		},
		LanguageKashmiri: {
			Code: LanguageKashmiri, Name: "Kashmiri", NativeName: "कॉशुर",
			Script: "Devanagari", Direction: "ltr", Region: "Kashmir", IsOfficial: true, IsSupported: true,
		},
		LanguageMaithili: {
			Code: LanguageMaithili, Name: "Maithili", NativeName: "मैथिली",
			Script: "Devanagari", Direction: "ltr", Region: "Bihar", IsOfficial: true, IsSupported: true,
		},
		LanguageDogri: {
			Code: LanguageDogri, Name: "Dogri", NativeName: "डोगरी",
			Script: "Devanagari", Direction: "ltr", Region: "Jammu", IsOfficial: true, IsSupported: true,
		},
		LanguageSindhi: {
			Code: LanguageSindhi, Name: "Sindhi", NativeName: "سندھی",
			Script: "Arabic", Direction: "rtl", Region: "Sindh", IsOfficial: true, IsSupported: true,
		},
	}
	
	if info, exists := languageMap[code]; exists {
		return info
	}
	
	return &LanguageInfo{
		Code: code, Name: string(code), NativeName: string(code),
		Script: "Unknown", Direction: "ltr", Region: "Unknown", IsOfficial: false, IsSupported: false,
	}
}

// GetSupportedLanguages returns list of all supported languages
func GetSupportedLanguages() []LanguageCode {
	return []LanguageCode{
		LanguageHindi, LanguageEnglish, LanguageBengali, LanguageTelugu,
		LanguageMarathi, LanguageTamil, LanguageUrdu, LanguageGujarati,
		LanguageKannada, LanguageMalayalam, LanguageOdia, LanguagePunjabi,
		LanguageAssamese, LanguageSanskrit, LanguageNepali, LanguageKonkani,
		LanguageManipuri, LanguageBodo, LanguageSantali, LanguageKashmiri,
		LanguageMaithili, LanguageDogri, LanguageSindhi,
	}
}

// GetOfficialIndianLanguages returns list of official Indian languages
func GetOfficialIndianLanguages() []LanguageCode {
	return GetSupportedLanguages() // All supported languages are official
}

// IdentityMessage represents localized identity system messages
type IdentityMessage struct {
	Key         string         `json:"key"`
	Category    string         `json:"category"`
	Text        *LocalizedText `json:"text"`
	Description string         `json:"description"`
}

// Common identity message keys
const (
	// Authentication messages
	MsgAuthenticationSuccess    = "auth.success"
	MsgAuthenticationFailed     = "auth.failed"
	MsgBiometricRequired       = "auth.biometric_required"
	MsgBiometricMismatch       = "auth.biometric_mismatch"
	MsgMultiFactorRequired     = "auth.mfa_required"
	
	// Identity creation messages
	MsgIdentityCreated         = "identity.created"
	MsgIdentityUpdateSuccess   = "identity.updated"
	MsgIdentityNotFound        = "identity.not_found"
	MsgDIDResolutionFailed     = "identity.did_resolution_failed"
	
	// Credential messages
	MsgCredentialIssued        = "credential.issued"
	MsgCredentialVerified      = "credential.verified"
	MsgCredentialExpired       = "credential.expired"
	MsgCredentialRevoked       = "credential.revoked"
	MsgCredentialInvalid       = "credential.invalid"
	
	// KYC messages
	MsgKYCVerificationStarted  = "kyc.verification_started"
	MsgKYCVerificationSuccess  = "kyc.verification_success"
	MsgKYCVerificationFailed   = "kyc.verification_failed"
	MsgKYCDocumentRequired     = "kyc.document_required"
	MsgKYCAadhaarLinked        = "kyc.aadhaar_linked"
	
	// Privacy messages
	MsgConsentRequired         = "privacy.consent_required"
	MsgConsentGranted          = "privacy.consent_granted"
	MsgConsentWithdrawn        = "privacy.consent_withdrawn"
	MsgDataMinimized           = "privacy.data_minimized"
	MsgPrivacyPolicyUpdated    = "privacy.policy_updated"
	
	// Error messages
	MsgInvalidInput            = "error.invalid_input"
	MsgInsufficientPermissions = "error.insufficient_permissions"
	MsgSystemError             = "error.system_error"
	MsgNetworkError            = "error.network_error"
	MsgTimeoutError            = "error.timeout"
	
	// Success messages
	MsgOperationComplete       = "success.operation_complete"
	MsgDataSaved              = "success.data_saved"
	MsgSecurityVerified       = "success.security_verified"
	
	// Cultural messages
	MsgWelcomeMessage         = "cultural.welcome"
	MsgFestivalGreeting       = "cultural.festival_greeting"
	MsgPatriotismQuote        = "cultural.patriotism_quote"
	MsgWisdomQuote            = "cultural.wisdom_quote"
)

// MessageCatalog manages localized messages
type MessageCatalog struct {
	Messages map[string]*IdentityMessage `json:"messages"`
}

// NewMessageCatalog creates a new message catalog
func NewMessageCatalog() *MessageCatalog {
	catalog := &MessageCatalog{
		Messages: make(map[string]*IdentityMessage),
	}
	catalog.initializeDefaultMessages()
	return catalog
}

// GetMessage returns a localized message
func (mc *MessageCatalog) GetMessage(key string, lang LanguageCode) string {
	if msg, exists := mc.Messages[key]; exists {
		return msg.Text.GetText(lang)
	}
	return fmt.Sprintf("[%s]", key) // Return key in brackets if not found
}

// AddMessage adds a new message to the catalog
func (mc *MessageCatalog) AddMessage(key, category, description string, text *LocalizedText) {
	mc.Messages[key] = &IdentityMessage{
		Key:         key,
		Category:    category,
		Text:        text,
		Description: description,
	}
}

// GetMessagesForCategory returns all messages in a category
func (mc *MessageCatalog) GetMessagesForCategory(category string) map[string]*IdentityMessage {
	result := make(map[string]*IdentityMessage)
	for key, msg := range mc.Messages {
		if msg.Category == category {
			result[key] = msg
		}
	}
	return result
}

// initializeDefaultMessages loads default messages in multiple languages
func (mc *MessageCatalog) initializeDefaultMessages() {
	// Authentication success message
	authSuccessText := NewLocalizedText(LanguageEnglish, "Authentication successful")
	authSuccessText.AddTranslation(LanguageHindi, "प्रमाणीकरण सफल")
	authSuccessText.AddTranslation(LanguageBengali, "প্রমাণীকরণ সফল")
	authSuccessText.AddTranslation(LanguageTamil, "அங்கீகரிப்பு வெற்றிகரமானது")
	authSuccessText.AddTranslation(LanguageTelugu, "ధృవీకరణ విజయవంతమైంది")
	authSuccessText.AddTranslation(LanguageMarathi, "प्रमाणीकरण यशस्वी")
	authSuccessText.AddTranslation(LanguageGujarati, "પ્રમાણીકરણ સફળ")
	authSuccessText.AddTranslation(LanguageKannada, "ದೃಢೀಕರಣ ಯಶಸ್ವಿ")
	authSuccessText.AddTranslation(LanguageMalayalam, "ആധികാരികത വിജയിച്ചു")
	authSuccessText.AddTranslation(LanguageUrdu, "تصدیق کامیاب")
	mc.AddMessage(MsgAuthenticationSuccess, "authentication", "User authentication successful", authSuccessText)
	
	// Biometric required message
	biometricText := NewLocalizedText(LanguageEnglish, "Biometric verification required for this operation")
	biometricText.AddTranslation(LanguageHindi, "इस ऑपरेशन के लिए बायोमेट्रिक सत्यापन आवश्यक है")
	biometricText.AddTranslation(LanguageBengali, "এই অপারেশনের জন্য বায়োমেট্রিক যাচাইকরণ প্রয়োজন")
	biometricText.AddTranslation(LanguageTamil, "இந்த செயல்பாட்டிற்கு உயிரியல் அளவீட்டு சரிபார்ப்பு தேவை")
	biometricText.AddTranslation(LanguageTelugu, "ఈ ఆపరేషన్ కోసం బయోమెట్రిక్ ధృవీకరణ అవసరం")
	biometricText.AddTranslation(LanguageMarathi, "या ऑपरेशनसाठी बायोमेट्रिक सत्यापन आवश्यक आहे")
	biometricText.AddTranslation(LanguageGujarati, "આ ઓપરેશન માટે બાયોમેટ્રિક ચકાસણી જરૂરી છે")
	biometricText.AddTranslation(LanguageKannada, "ಈ ಕಾರ್ಯಾಚರಣೆಗೆ ಬಯೋಮೆಟ್ರಿಕ್ ಪರಿಶೀಲನೆ ಅಗತ್ಯವಿದೆ")
	biometricText.AddTranslation(LanguageMalayalam, "ഈ പ്രവർത്തനത്തിന് ബയോമെട്രിക് പരിശോധന ആവശ്യമാണ്")
	biometricText.AddTranslation(LanguageUrdu, "اس آپریشن کے لیے بائیو میٹرک تصدیق ضروری ہے")
	mc.AddMessage(MsgBiometricRequired, "authentication", "Biometric verification required", biometricText)
	
	// Identity created message
	identityCreatedText := NewLocalizedText(LanguageEnglish, "Your digital identity has been created successfully")
	identityCreatedText.AddTranslation(LanguageHindi, "आपकी डिजिटल पहचान सफलतापूर्वक बनाई गई है")
	identityCreatedText.AddTranslation(LanguageBengali, "আপনার ডিজিটাল পরিচয় সফলভাবে তৈরি হয়েছে")
	identityCreatedText.AddTranslation(LanguageTamil, "உங்கள் டிஜிட்டல் அடையாளம் வெற்றிகரமாக உருவாக்கப்பட்டது")
	identityCreatedText.AddTranslation(LanguageTelugu, "మీ డిజిటల్ గుర్తింపు విజయవంతంగా సృష్టించబడింది")
	identityCreatedText.AddTranslation(LanguageMarathi, "तुमची डिजिटल ओळख यशस्वीरित्या तयार झाली आहे")
	identityCreatedText.AddTranslation(LanguageGujarati, "તમારી ડિજિટલ ઓળખ સફળતાપૂર્વક બનાવવામાં આવી છે")
	identityCreatedText.AddTranslation(LanguageKannada, "ನಿಮ್ಮ ಡಿಜಿಟಲ್ ಗುರುತು ಯಶಸ್ವಿಯಾಗಿ ರಚಿಸಲಾಗಿದೆ")
	identityCreatedText.AddTranslation(LanguageMalayalam, "നിങ്ങളുടെ ഡിജിറ്റൽ ഐഡന്റിറ്റി വിജയകരമായി സൃഷ്ടിച്ചു")
	identityCreatedText.AddTranslation(LanguageUrdu, "آپ کی ڈیجیٹل شناخت کامیابی سے بنائی گئی ہے")
	mc.AddMessage(MsgIdentityCreated, "identity", "Digital identity creation success", identityCreatedText)
	
	// KYC verification success message
	kycSuccessText := NewLocalizedText(LanguageEnglish, "KYC verification completed successfully")
	kycSuccessText.AddTranslation(LanguageHindi, "केवाईसी सत्यापन सफलतापूर्वक पूरा हुआ")
	kycSuccessText.AddTranslation(LanguageBengali, "কেওয়াইসি যাচাইকরণ সফলভাবে সম্পন্ন হয়েছে")
	kycSuccessText.AddTranslation(LanguageTamil, "KYC சரிபார்ப்பு வெற்றிகரமாக முடிந்தது")
	kycSuccessText.AddTranslation(LanguageTelugu, "KYC ధృవీకరణ విజయవంతంగా పూర్తయింది")
	kycSuccessText.AddTranslation(LanguageMarathi, "KYC सत्यापन यशस्वीरित्या पूर्ण झाले")
	kycSuccessText.AddTranslation(LanguageGujarati, "KYC ચકાસણી સફળતાપૂર્વક પૂર્ણ થઈ")
	kycSuccessText.AddTranslation(LanguageKannada, "KYC ಪರಿಶೀಲನೆ ಯಶಸ್ವಿಯಾಗಿ ಪೂರ್ಣಗೊಂಡಿತು")
	kycSuccessText.AddTranslation(LanguageMalayalam, "KYC പരിശോധന വിജയകരമായി പൂർത്തിയായി")
	kycSuccessText.AddTranslation(LanguageUrdu, "KYC تصدیق کامیابی سے مکمل ہوئی")
	mc.AddMessage(MsgKYCVerificationSuccess, "kyc", "KYC verification success", kycSuccessText)
	
	// Consent required message
	consentText := NewLocalizedText(LanguageEnglish, "Your consent is required to share this information")
	consentText.AddTranslation(LanguageHindi, "इस जानकारी को साझा करने के लिए आपकी सहमति आवश्यक है")
	consentText.AddTranslation(LanguageBengali, "এই তথ্য শেয়ার করার জন্য আপনার সম্মতি প্রয়োজন")
	consentText.AddTranslation(LanguageTamil, "இந்த தகவலைப் பகிர உங்கள் ஒப்புதல் தேவை")
	consentText.AddTranslation(LanguageTelugu, "ఈ సమాచారాన్ని పంచుకోవడానికి మీ అనుమతి అవసరం")
	consentText.AddTranslation(LanguageMarathi, "ही माहिती सामायिक करण्यासाठी तुमची संमती आवश्यक आहे")
	consentText.AddTranslation(LanguageGujarati, "આ માહિતી શેર કરવા માટે તમારી સંમતિ જરૂરી છે")
	consentText.AddTranslation(LanguageKannada, "ಈ ಮಾಹಿತಿಯನ್ನು ಹಂಚಿಕೊಳ್ಳಲು ನಿಮ್ಮ ಒಪ್ಪಿಗೆ ಅಗತ್ಯವಿದೆ")
	consentText.AddTranslation(LanguageMalayalam, "ഈ വിവരങ്ങൾ പങ്കിടാൻ നിങ്ങളുടെ സമ്മതം ആവശ്യമാണ്")
	consentText.AddTranslation(LanguageUrdu, "اس معلومات کو شیئر کرنے کے لیے آپ کی رضامندی ضروری ہے")
	mc.AddMessage(MsgConsentRequired, "privacy", "Consent required for data sharing", consentText)
	
	// Welcome message with cultural context
	welcomeText := NewLocalizedText(LanguageEnglish, "Welcome to DeshChain Identity - Your Digital Sovereignty Begins Here")
	welcomeText.AddTranslation(LanguageHindi, "देशचेन आइडेंटिटी में आपका स्वागत है - आपकी डिजिटल संप्रभुता यहाँ शुरू होती है")
	welcomeText.AddTranslation(LanguageSanskrit, "देशश्रृंखला पहचान में आपका स्वागतम् - अत्र आपस्य डिजिटल स्वराज्यम् आरभते")
	welcomeText.AddTranslation(LanguageBengali, "দেশচেইন আইডেন্টিটিতে আপনাকে স্বাগতম - আপনার ডিজিটাল সার্বভৌমত্ব এখানে শুরু হয়")
	welcomeText.AddTranslation(LanguageTamil, "தேஷ்சேன் அடையாளத்திற்கு வரவேற்கிறோம் - உங்கள் டிஜிட்டல் இறையாண்மை இங்கே தொடங்குகிறது")
	welcomeText.AddTranslation(LanguageTelugu, "దేష్‌చైన్ ఐడెంటిటీకి స్వాగతం - మీ డిజిటల్ సార్వభౌమత్వం ఇక్కడ మొదలవుతుంది")
	mc.AddMessage(MsgWelcomeMessage, "cultural", "Welcome message with cultural context", welcomeText)
}

// FormatMessage formats a message with parameters
func (mc *MessageCatalog) FormatMessage(key string, lang LanguageCode, params map[string]string) string {
	message := mc.GetMessage(key, lang)
	
	// Replace parameters in the message
	for param, value := range params {
		placeholder := fmt.Sprintf("{%s}", param)
		message = strings.ReplaceAll(message, placeholder, value)
	}
	
	return message
}

// DetectLanguageFromInput attempts to detect language from user input
func DetectLanguageFromInput(input string) LanguageCode {
	// Simple detection based on script/characters
	// This could be enhanced with proper language detection libraries
	
	// Check for Devanagari script (Hindi, Marathi, Sanskrit, Nepali)
	for _, char := range input {
		if char >= 0x0900 && char <= 0x097F {
			return LanguageHindi // Default to Hindi for Devanagari
		}
	}
	
	// Check for Bengali script
	for _, char := range input {
		if char >= 0x0980 && char <= 0x09FF {
			return LanguageBengali
		}
	}
	
	// Check for Tamil script
	for _, char := range input {
		if char >= 0x0B80 && char <= 0x0BFF {
			return LanguageTamil
		}
	}
	
	// Check for Telugu script
	for _, char := range input {
		if char >= 0x0C00 && char <= 0x0C7F {
			return LanguageTelugu
		}
	}
	
	// Check for Gujarati script
	for _, char := range input {
		if char >= 0x0A80 && char <= 0x0AFF {
			return LanguageGujarati
		}
	}
	
	// Check for Kannada script
	for _, char := range input {
		if char >= 0x0C80 && char <= 0x0CFF {
			return LanguageKannada
		}
	}
	
	// Check for Malayalam script
	for _, char := range input {
		if char >= 0x0D00 && char <= 0x0D7F {
			return LanguageMalayalam
		}
	}
	
	// Check for Arabic script (Urdu)
	for _, char := range input {
		if char >= 0x0600 && char <= 0x06FF {
			return LanguageUrdu
		}
	}
	
	// Default to English for Latin script or unknown
	return LanguageEnglish
}

// GetRegionalLanguages returns languages for a specific region/state
func GetRegionalLanguages(region string) []LanguageCode {
	regionMap := map[string][]LanguageCode{
		"maharashtra": {LanguageMarathi, LanguageHindi, LanguageEnglish},
		"gujarat":     {LanguageGujarati, LanguageHindi, LanguageEnglish},
		"punjab":      {LanguagePunjabi, LanguageHindi, LanguageEnglish},
		"karnataka":   {LanguageKannada, LanguageEnglish, LanguageHindi},
		"tamil_nadu":  {LanguageTamil, LanguageEnglish},
		"telangana":   {LanguageTelugu, LanguageHindi, LanguageEnglish},
		"andhra_pradesh": {LanguageTelugu, LanguageHindi, LanguageEnglish},
		"kerala":      {LanguageMalayalam, LanguageEnglish, LanguageHindi},
		"west_bengal": {LanguageBengali, LanguageHindi, LanguageEnglish},
		"odisha":      {LanguageOdia, LanguageHindi, LanguageEnglish},
		"assam":       {LanguageAssamese, LanguageHindi, LanguageEnglish},
		"kashmir":     {LanguageKashmiri, LanguageUrdu, LanguageHindi, LanguageEnglish},
		"goa":         {LanguageKonkani, LanguageMarathi, LanguageEnglish},
		"manipur":     {LanguageManipuri, LanguageEnglish, LanguageHindi},
		"bihar":       {LanguageMaithili, LanguageHindi, LanguageEnglish},
		"jharkhand":   {LanguageSantali, LanguageHindi, LanguageEnglish},
	}
	
	regionKey := strings.ToLower(strings.ReplaceAll(region, " ", "_"))
	if languages, exists := regionMap[regionKey]; exists {
		return languages
	}
	
	// Default to Hindi and English
	return []LanguageCode{LanguageHindi, LanguageEnglish}
}

// LocalizationConfig contains localization settings
type LocalizationConfig struct {
	DefaultLanguage    LanguageCode   `json:"default_language"`
	FallbackLanguages  []LanguageCode `json:"fallback_languages"`
	Region            string         `json:"region"`
	EnableAutoDetect  bool           `json:"enable_auto_detect"`
	EnableRTLSupport  bool           `json:"enable_rtl_support"`
}

// DefaultLocalizationConfig returns default localization configuration
func DefaultLocalizationConfig() *LocalizationConfig {
	return &LocalizationConfig{
		DefaultLanguage:   LanguageEnglish,
		FallbackLanguages: []LanguageCode{LanguageHindi, LanguageEnglish},
		Region:           "india",
		EnableAutoDetect: true,
		EnableRTLSupport: true,
	}
}

// IsRTLLanguage checks if a language is right-to-left
func IsRTLLanguage(lang LanguageCode) bool {
	rtlLanguages := map[LanguageCode]bool{
		LanguageUrdu:   true,
		LanguageArabic: true,
	}
	return rtlLanguages[lang]
}