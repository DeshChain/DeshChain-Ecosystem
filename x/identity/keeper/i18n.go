package keeper

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/DeshChain/DeshChain-Ecosystem/x/identity/types"
)

// I18nKeeper manages internationalization for the identity module
type I18nKeeper struct {
	storeKey   sdk.StoreKey
	cdc        codec.BinaryCodec
	catalog    *types.MessageCatalog
	config     *types.LocalizationConfig
}

// NewI18nKeeper creates a new I18n keeper
func NewI18nKeeper(storeKey sdk.StoreKey, cdc codec.BinaryCodec) *I18nKeeper {
	return &I18nKeeper{
		storeKey: storeKey,
		cdc:      cdc,
		catalog:  types.NewMessageCatalog(),
		config:   types.DefaultLocalizationConfig(),
	}
}

// SetLocalizationConfig updates the localization configuration
func (k *I18nKeeper) SetLocalizationConfig(ctx sdk.Context, config *types.LocalizationConfig) error {
	store := ctx.KVStore(k.storeKey)
	
	bz, err := k.cdc.Marshal(config)
	if err != nil {
		return errors.Wrap(err, "failed to marshal localization config")
	}
	
	store.Set(types.LocalizationConfigKey, bz)
	k.config = config
	
	return nil
}

// GetLocalizationConfig returns the current localization configuration
func (k *I18nKeeper) GetLocalizationConfig(ctx sdk.Context) *types.LocalizationConfig {
	store := ctx.KVStore(k.storeKey)
	
	bz := store.Get(types.LocalizationConfigKey)
	if bz == nil {
		return types.DefaultLocalizationConfig()
	}
	
	var config types.LocalizationConfig
	if err := k.cdc.Unmarshal(bz, &config); err != nil {
		return types.DefaultLocalizationConfig()
	}
	
	return &config
}

// GetLocalizedMessage returns a message in the requested language
func (k *I18nKeeper) GetLocalizedMessage(ctx sdk.Context, key string, lang types.LanguageCode) string {
	return k.catalog.GetMessage(key, lang)
}

// FormatLocalizedMessage returns a formatted message with parameters
func (k *I18nKeeper) FormatLocalizedMessage(ctx sdk.Context, key string, lang types.LanguageCode, params map[string]string) string {
	return k.catalog.FormatMessage(key, lang, params)
}

// AddCustomMessage adds a custom localized message to the catalog
func (k *I18nKeeper) AddCustomMessage(ctx sdk.Context, key, category, description string, text *types.LocalizedText) error {
	store := ctx.KVStore(k.storeKey)
	
	message := &types.IdentityMessage{
		Key:         key,
		Category:    category,
		Text:        text,
		Description: description,
	}
	
	bz, err := k.cdc.Marshal(message)
	if err != nil {
		return errors.Wrap(err, "failed to marshal custom message")
	}
	
	customKey := append(types.CustomMessagePrefix, []byte(key)...)
	store.Set(customKey, bz)
	
	// Add to in-memory catalog
	k.catalog.AddMessage(key, category, description, text)
	
	return nil
}

// GetCustomMessage retrieves a custom message
func (k *I18nKeeper) GetCustomMessage(ctx sdk.Context, key string) (*types.IdentityMessage, error) {
	store := ctx.KVStore(k.storeKey)
	
	customKey := append(types.CustomMessagePrefix, []byte(key)...)
	bz := store.Get(customKey)
	if bz == nil {
		return nil, errors.Wrapf(types.ErrMessageNotFound, "custom message not found: %s", key)
	}
	
	var message types.IdentityMessage
	if err := k.cdc.Unmarshal(bz, &message); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal custom message")
	}
	
	return &message, nil
}

// GetAllCustomMessages returns all custom messages
func (k *I18nKeeper) GetAllCustomMessages(ctx sdk.Context) ([]*types.IdentityMessage, error) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.CustomMessagePrefix)
	defer iterator.Close()
	
	var messages []*types.IdentityMessage
	
	for ; iterator.Valid(); iterator.Next() {
		var message types.IdentityMessage
		if err := k.cdc.Unmarshal(iterator.Value(), &message); err != nil {
			continue // Skip invalid messages
		}
		messages = append(messages, &message)
	}
	
	return messages, nil
}

// DetectUserLanguage attempts to detect user's preferred language
func (k *I18nKeeper) DetectUserLanguage(ctx sdk.Context, userDID string, input string, region string) types.LanguageCode {
	// First try to get stored user preference
	if lang := k.getUserLanguagePreference(ctx, userDID); lang != "" {
		return lang
	}
	
	// Try to detect from input if auto-detection is enabled
	if k.config.EnableAutoDetect && input != "" {
		if detected := types.DetectLanguageFromInput(input); detected != types.LanguageEnglish {
			return detected
		}
	}
	
	// Use regional languages if region is provided
	if region != "" {
		regionalLangs := types.GetRegionalLanguages(region)
		if len(regionalLangs) > 0 {
			return regionalLangs[0] // Return primary regional language
		}
	}
	
	// Fallback to default language
	return k.config.DefaultLanguage
}

// SetUserLanguagePreference stores user's language preference
func (k *I18nKeeper) SetUserLanguagePreference(ctx sdk.Context, userDID string, lang types.LanguageCode) error {
	// Validate language is supported
	if !k.isLanguageSupported(lang) {
		return errors.Wrapf(types.ErrUnsupportedLanguage, "language not supported: %s", lang)
	}
	
	store := ctx.KVStore(k.storeKey)
	key := append(types.UserLanguagePrefix, []byte(userDID)...)
	
	langBytes := []byte(lang)
	store.Set(key, langBytes)
	
	return nil
}

// getUserLanguagePreference retrieves user's stored language preference
func (k *I18nKeeper) getUserLanguagePreference(ctx sdk.Context, userDID string) types.LanguageCode {
	store := ctx.KVStore(k.storeKey)
	key := append(types.UserLanguagePrefix, []byte(userDID)...)
	
	bz := store.Get(key)
	if bz == nil {
		return ""
	}
	
	return types.LanguageCode(bz)
}

// isLanguageSupported checks if a language is supported
func (k *I18nKeeper) isLanguageSupported(lang types.LanguageCode) bool {
	supportedLangs := types.GetSupportedLanguages()
	for _, supported := range supportedLangs {
		if supported == lang {
			return true
		}
	}
	return false
}

// GetSupportedLanguages returns list of supported languages with metadata
func (k *I18nKeeper) GetSupportedLanguages(ctx sdk.Context) []*types.LanguageInfo {
	supportedCodes := types.GetSupportedLanguages()
	languages := make([]*types.LanguageInfo, len(supportedCodes))
	
	for i, code := range supportedCodes {
		languages[i] = types.GetLanguageInfo(code)
	}
	
	// Sort by name for consistent ordering
	sort.Slice(languages, func(i, j int) bool {
		return languages[i].Name < languages[j].Name
	})
	
	return languages
}

// GetRegionalLanguages returns languages for a specific region
func (k *I18nKeeper) GetRegionalLanguages(ctx sdk.Context, region string) []*types.LanguageInfo {
	regionalCodes := types.GetRegionalLanguages(region)
	languages := make([]*types.LanguageInfo, len(regionalCodes))
	
	for i, code := range regionalCodes {
		languages[i] = types.GetLanguageInfo(code)
	}
	
	return languages
}

// LocalizeIdentityText localizes identity-related text based on user preferences
func (k *I18nKeeper) LocalizeIdentityText(ctx sdk.Context, userDID string, text *types.LocalizedText) string {
	userLang := k.getUserLanguagePreference(ctx, userDID)
	if userLang == "" {
		userLang = k.config.DefaultLanguage
	}
	
	localizedText := text.GetText(userLang)
	
	// If text is not available in user's language, try fallback languages
	if localizedText == "" {
		for _, fallbackLang := range k.config.FallbackLanguages {
			if localizedText = text.GetText(fallbackLang); localizedText != "" {
				break
			}
		}
	}
	
	return localizedText
}

// GenerateCulturalGreeting generates a culturally appropriate greeting
func (k *I18nKeeper) GenerateCulturalGreeting(ctx sdk.Context, userDID string, timeOfDay string, festival string) string {
	userLang := k.getUserLanguagePreference(ctx, userDID)
	if userLang == "" {
		userLang = k.config.DefaultLanguage
	}
	
	// Festival-specific greetings
	if festival != "" {
		return k.getFestivalGreeting(userLang, festival)
	}
	
	// Time-based greetings
	return k.getTimeBasedGreeting(userLang, timeOfDay)
}

// getFestivalGreeting returns festival-specific greeting
func (k *I18nKeeper) getFestivalGreeting(lang types.LanguageCode, festival string) string {
	festivalGreetings := map[string]map[types.LanguageCode]string{
		"diwali": {
			types.LanguageHindi:    "दीपावली की हार्दिक शुभकामनाएं!",
			types.LanguageEnglish:  "Happy Diwali! May the festival of lights illuminate your path.",
			types.LanguageBengali:  "দীপাবলির শুভেচ্ছা!",
			types.LanguageTamil:    "தீபாவளி வாழ்த்துகள்!",
			types.LanguageTelugu:   "దీపావళి శుభాకాంక్షలు!",
			types.LanguageMarathi:  "दिवाळीच्या हार्दिक शुभेच्छा!",
			types.LanguageGujarati: "દિવાળીની શુભકામનાઓ!",
		},
		"holi": {
			types.LanguageHindi:    "होली की रंगबिरंगी शुभकामनाएं!",
			types.LanguageEnglish:  "Happy Holi! May colors of joy fill your life.",
			types.LanguageBengali:  "হোলির শুভেচ্ছা!",
			types.LanguageTamil:    "ஹோலி வாழ்த்துகள்!",
			types.LanguageTelugu:   "హోలీ శుభాకాంక్షలు!",
			types.LanguageMarathi:  "होळीच्या शुभेच्छा!",
			types.LanguageGujarati: "હોળીની શુભકામનાઓ!",
		},
		"dussehra": {
			types.LanguageHindi:    "दशहरा की शुभकामनाएं! बुराई पर अच्छाई की जीत!",
			types.LanguageEnglish:  "Happy Dussehra! Victory of good over evil.",
			types.LanguageBengali:  "দুর্গা পূজার শুভেচ্ছা!",
			types.LanguageTamil:    "விஜயதசமி வாழ்த்துகள்!",
			types.LanguageTelugu:   "దసరా శుభాకాంక్షలు!",
			types.LanguageMarathi:  "दसऱ्याच्या शुभेच्छा!",
		},
		"ganesh_chaturthi": {
			types.LanguageHindi:    "गणेश चतुर्थी की शुभकामनाएं!",
			types.LanguageEnglish:  "Happy Ganesh Chaturthi! May Lord Ganesha bless you.",
			types.LanguageMarathi:  "गणेश चतुर्थीच्या हार्दिक शुभेच्छा! गणपती बाप्पा मोरया!",
			types.LanguageTelugu:   "గణేష్ చతుర్థి శుభాకాంక్షలు!",
			types.LanguageKannada:  "ಗಣೇಶ ಚತುರ್ಥಿ ಶುಭಾಶಯಗಳು!",
		},
		"eid": {
			types.LanguageUrdu:     "عید مبارک! خوشیوں اور برکتوں کا دن",
			types.LanguageEnglish:  "Eid Mubarak! May this blessed day bring joy and peace.",
			types.LanguageHindi:    "ईद मुबारक!",
			types.LanguageBengali:  "ঈদ মুবারক!",
		},
		"christmas": {
			types.LanguageEnglish:  "Merry Christmas! May the spirit of love and joy fill your heart.",
			types.LanguageHindi:    "क्रिसमस की शुभकामनाएं!",
			types.LanguageTamil:    "கிறிஸ்துமஸ் வாழ்த்துகள்!",
			types.LanguageMalayalam: "ക്രിസ്മസ് ആശംസകൾ!",
		},
	}
	
	if greetings, exists := festivalGreetings[festival]; exists {
		if greeting, exists := greetings[lang]; exists {
			return greeting
		}
		// Fallback to English
		if greeting, exists := greetings[types.LanguageEnglish]; exists {
			return greeting
		}
	}
	
	return "Happy Celebrations!"
}

// getTimeBasedGreeting returns time-appropriate greeting
func (k *I18nKeeper) getTimeBasedGreeting(lang types.LanguageCode, timeOfDay string) string {
	timeGreetings := map[string]map[types.LanguageCode]string{
		"morning": {
			types.LanguageHindi:     "सुप्रभात! आपका दिन शुभ हो।",
			types.LanguageEnglish:   "Good morning! Have a blessed day.",
			types.LanguageBengali:   "সুপ্রভাত! আপনার দিন মঙ্গলময় হোক।",
			types.LanguageTamil:     "காலை வணக்கம்! உங்கள் நாள் நல்லதாக இருக்கட்டும்।",
			types.LanguageTelugu:    "సుప్రభాతం! మీ రోజు మంచిదిగా ఉండాలని ఆశిస్తున్నాను।",
			types.LanguageMarathi:   "सुप्रभात! तुमचा दिवस शुभ जावो।",
			types.LanguageGujarati:  "સુપ્રભાત! તમારો દિવસ શુભ રહે.",
			types.LanguageKannada:   "ಶುಭೋದಯ! ನಿಮ್ಮ ದಿನ ಶುಭವಾಗಿರಲಿ.",
			types.LanguageMalayalam: "സുപ്രഭാതം! നിങ്ങളുടെ ദിവസം നല്ലതായിരിക്കട്ടെ.",
			types.LanguageUrdu:      "صبح بخیر! آپ کا دن اچھا گزرے۔",
			types.LanguageSanskrit:  "सुप्रभातम्! भवतः दिवसः शुभः भवतु।",
		},
		"afternoon": {
			types.LanguageHindi:     "नमस्कार! आपका दिन कैसा जा रहा है?",
			types.LanguageEnglish:   "Good afternoon! Hope your day is going well.",
			types.LanguageBengali:   "নমস্কার! আপনার দিন কেমন কাটছে?",
			types.LanguageTamil:     "மதிய வணக்கம்! உங்கள் நாள் எப்படி செல்கிறது?",
			types.LanguageTelugu:    "నమస్కారం! మీ రోజు ఎలా గడుచుతోంది?",
			types.LanguageMarathi:   "नमस्कार! तुमचा दिवस कसा जात आहे?",
			types.LanguageGujarati:  "નમસ્કાર! તમારો દિવસ કેવો ચાલી રહ્યો છે?",
			types.LanguageUrdu:      "السلام علیکم! آپ کا دن کیسا گزر رہا ہے؟",
		},
		"evening": {
			types.LanguageHindi:     "शुभ संध्या! आपका दिन अच्छा रहा हो।",
			types.LanguageEnglish:   "Good evening! Hope you had a productive day.",
			types.LanguageBengali:   "শুভ সন্ধ্যা! আপনার দিন ভালো কেটেছে আশা করি।",
			types.LanguageTamil:     "மாலை வணக்கம்! உங்கள் நாள் நன்றாக இருந்திருக்கும் என நம்புகிறேன்।",
			types.LanguageTelugu:    "శుభ సాయంత్రం! మీ రోజు బాగా గడిచిందని ఆశిస్తున్నాను।",
			types.LanguageMarathi:   "शुभ संध्या! तुमचा दिवस चांगला गेला असेल.",
			types.LanguageGujarati:  "શુભ સાંજ! તમારો દિવસ સારો ગયો હશે.",
			types.LanguageUrdu:      "شام بخیر! آپ کا دن اچھا گزرا ہوگا۔",
		},
		"night": {
			types.LanguageHindi:     "शुभ रात्रि! आपको मीठे सपने आएं।",
			types.LanguageEnglish:   "Good night! Sweet dreams and peaceful rest.",
			types.LanguageBengali:   "শুভ রাত্রি! মিষ্টি স্বপ্ন দেখুন।",
			types.LanguageTamil:     "இனிய இரவு வணக்கம்! இனிமையான கனவுகள்.",
			types.LanguageTelugu:    "శుభ రాత్రి! మధుర స్వప్నలు.",
			types.LanguageMarathi:   "शुभ रात्री! गोड स्वप्न पडोत.",
			types.LanguageGujarati:  "શુભ રાત્રિ! મીઠા સપના આવે.",
			types.LanguageUrdu:      "شب بخیر! خوشگوار خواب دیکھیں۔",
		},
	}
	
	if greetings, exists := timeGreetings[timeOfDay]; exists {
		if greeting, exists := greetings[lang]; exists {
			return greeting
		}
	}
	
	// Default greeting
	defaultGreetings := map[types.LanguageCode]string{
		types.LanguageHindi:    "नमस्कार! देशचेन में आपका स्वागत है।",
		types.LanguageEnglish:  "Welcome to DeshChain Identity!",
		types.LanguageBengali:  "দেশচেইনে আপনাকে স্বাগতম!",
		types.LanguageTamil:    "தேஷ்சேனில் உங்களை வரவேற்கிறோம்!",
		types.LanguageTelugu:   "దేష్‌చైన్‌కు స్వాగతం!",
		types.LanguageMarathi:  "देशचेनमध्ये तुमचे स्वागत आहे!",
		types.LanguageGujarati: "દેશચેનમાં તમારું સ્વાગત છે!",
		types.LanguageUrdu:     "دیش چین میں آپ کا خوش آمدید!",
	}
	
	if greeting, exists := defaultGreetings[lang]; exists {
		return greeting
	}
	
	return "Welcome to DeshChain Identity!"
}

// GetCulturalQuote returns a culturally relevant quote in user's language
func (k *I18nKeeper) GetCulturalQuote(ctx sdk.Context, userDID string, category string) string {
	userLang := k.getUserLanguagePreference(ctx, userDID)
	if userLang == "" {
		userLang = k.config.DefaultLanguage
	}
	
	// Cultural quotes categorized by theme
	quotes := k.getCulturalQuotes()
	
	if categoryQuotes, exists := quotes[category]; exists {
		if langQuotes, exists := categoryQuotes[userLang]; exists {
			// Return a random quote from the category
			if len(langQuotes) > 0 {
				// For deterministic selection based on context height
				blockHeight := ctx.BlockHeight()
				index := int(blockHeight) % len(langQuotes)
				return langQuotes[index]
			}
		}
	}
	
	// Fallback to wisdom quotes in English
	wisdomQuotes := []string{
		"The best way to find yourself is to lose yourself in the service of others. - Mahatma Gandhi",
		"You have to dream before your dreams can come true. - A.P.J. Abdul Kalam",
		"In a gentle way, you can shake the world. - Mahatma Gandhi",
		"Learning gives creativity, creativity leads to thinking, thinking provides knowledge, knowledge makes you great. - A.P.J. Abdul Kalam",
	}
	
	if len(wisdomQuotes) > 0 {
		blockHeight := ctx.BlockHeight()
		index := int(blockHeight) % len(wisdomQuotes)
		return wisdomQuotes[index]
	}
	
	return "Every step forward is a step toward achieving something bigger and better than your current situation."
}

// getCulturalQuotes returns categorized cultural quotes
func (k *I18nKeeper) getCulturalQuotes() map[string]map[types.LanguageCode][]string {
	return map[string]map[types.LanguageCode][]string{
		"wisdom": {
			types.LanguageHindi: {
				"अहिंसा परमो धर्मः - हिंसा न करना सबसे बड़ा धर्म है।",
				"वसुधैव कुटुम्बकम् - पूरा विश्व एक परिवार है।",
				"सत्यमेव जयते - सत्य की ही जीत होती है।",
				"अविद्या मृत्युः विद्या अमृतम् - अज्ञानता मृत्यु है, ज्ञान अमरता है।",
			},
			types.LanguageEnglish: {
				"Be the change you wish to see in the world. - Mahatma Gandhi",
				"The best way to find yourself is to lose yourself in the service of others. - Mahatma Gandhi",
				"You have to dream before your dreams can come true. - A.P.J. Abdul Kalam",
				"Learning gives creativity, creativity leads to thinking. - A.P.J. Abdul Kalam",
			},
			types.LanguageSanskrit: {
				"सर्वे भवन्तु सुखिनः सर्वे सन्तु निरामयाः।",
				"अहं ब्रह्मास्मि - मैं ब्रह्म हूँ।",
				"तत् त्वम् असि - तू वही है।",
				"एकं सत् विप्रा बहुधा वदन्ति - सत्य एक है, ऋषि इसे अनेक प्रकार से कहते हैं।",
			},
		},
		"motivation": {
			types.LanguageHindi: {
				"कर्म करो, फल की चिंता मत करो।",
				"जो व्यक्ति अपने लक्ष्य को पाने के लिए पूरी तरह प्रतिबद्ध है, वह असंभव को भी संभव बना देता है।",
				"असफलता एक विकल्प नहीं है।",
				"जीतना और हारना तो जिंदगी का हिस्सा है, कोशिश करना कभी नहीं छोड़ना चाहिए।",
			},
			types.LanguageEnglish: {
				"Excellence is a continuous process and not an accident. - A.P.J. Abdul Kalam",
				"Don't take rest after your first victory because if you fail in second, more lips are waiting to say that your first victory was just luck.",
				"All of us do not have equal talent. But, all of us have an equal opportunity to develop our talents.",
				"If you want to shine like a sun, first burn like a sun. - A.P.J. Abdul Kalam",
			},
		},
		"patriotism": {
			types.LanguageHindi: {
				"जन गण मन अधिनायक जय हे, भारत भाग्य विधाता।",
				"हिंद देश के निवासी सभी जन एक हैं।",
				"मातृभूमि स्वर्ग से महान है।",
				"भारत माता की जय!",
			},
			types.LanguageEnglish: {
				"Freedom is not worth having if it does not include the freedom to make mistakes. - Mahatma Gandhi",
				"The future depends on what you do today. - Mahatma Gandhi",
				"A nation's culture resides in the hearts and in the soul of its people. - Mahatma Gandhi",
				"Unity in diversity is India's strength.",
			},
		},
		"technology": {
			types.LanguageHindi: {
				"प्रौद्योगिकी तभी सफल है जब वह मानवता की सेवा करे।",
				"नवाचार में भारत की अनंत संभावनाएं हैं।",
				"डिजिटल इंडिया - सबका साथ, सबका विकास।",
			},
			types.LanguageEnglish: {
				"Technology is best when it brings people together.",
				"Innovation distinguishes between a leader and a follower.",
				"The advance of technology is based on making it fit in so that you don't really even notice it.",
				"Digital India is not just about technology, it's about transformation.",
			},
		},
	}
}

// LocalizeErrorMessage localizes error messages for better user experience
func (k *I18nKeeper) LocalizeErrorMessage(ctx sdk.Context, userDID string, err error) string {
	userLang := k.getUserLanguagePreference(ctx, userDID)
	if userLang == "" {
		userLang = k.config.DefaultLanguage
	}
	
	// Map common errors to localized messages
	errorMessages := map[string]map[types.LanguageCode]string{
		"identity not found": {
			types.LanguageHindi:    "पहचान नहीं मिली। कृपया पहले पंजीकरण करें।",
			types.LanguageEnglish:  "Identity not found. Please register first.",
			types.LanguageBengali:  "পরিচয় পাওয়া যায়নি। দয়া করে প্রথমে নিবন্ধন করুন।",
			types.LanguageTamil:    "அடையாளம் காணப்படவில்லை. முதலில் பதிவு செய்யவும்.",
			types.LanguageTelugu:   "గుర్తింపు కనుగొనబడలేదు. దయచేసి మొదట నమోదు చేసుకోండి।",
		},
		"biometric mismatch": {
			types.LanguageHindi:    "बायोमेट्रिक मैच नहीं हुआ। कृपया फिर से कोशिश करें।",
			types.LanguageEnglish:  "Biometric authentication failed. Please try again.",
			types.LanguageBengali:  "বায়োমেট্রিক মিল হয়নি। দয়া করে আবার চেষ্টা করুন।",
			types.LanguageTamil:    "உயிரியல் அளவீட்டு பொருத்தம் தோல்வி. மீண்டும் முயற்சிக்கவும்.",
			types.LanguageTelugu:   "బయోమెట్రిక్ ప్రమాణీకరణ విఫలమైంది. దయచేసి మళ్లీ ప్రయత్నించండి।",
		},
		"insufficient permissions": {
			types.LanguageHindi:    "अपर्याप्त अनुमतियां। आपको यह कार्य करने की अनुमति नहीं है।",
			types.LanguageEnglish:  "Insufficient permissions. You are not authorized for this action.",
			types.LanguageBengali:  "অপর্যাপ্ত অনুমতি। আপনার এই কাজের অনুমতি নেই।",
			types.LanguageTamil:    "போதுமான அனுமதிகள் இல்லை. இந்த செயலுக்கு நீங்கள் அங்கீகரிக்கப்படவில்லை.",
			types.LanguageTelugu:   "తగినంత అనుమతులు లేవు. ఈ చర్యకు మీకు అధికారం లేదు.",
		},
	}
	
	errStr := strings.ToLower(err.Error())
	
	for pattern, translations := range errorMessages {
		if strings.Contains(errStr, pattern) {
			if msg, exists := translations[userLang]; exists {
				return msg
			}
			// Fallback to English
			if msg, exists := translations[types.LanguageEnglish]; exists {
				return msg
			}
		}
	}
	
	// Return original error if no translation found
	return err.Error()
}

// GetLocalizationStats returns statistics about localization usage
func (k *I18nKeeper) GetLocalizationStats(ctx sdk.Context) map[string]interface{} {
	// This would typically be stored and updated, but for now return sample data
	return map[string]interface{}{
		"total_supported_languages": len(types.GetSupportedLanguages()),
		"total_messages":            len(k.catalog.Messages),
		"most_used_language":        types.LanguageHindi,
		"coverage_percentage":       95.5,
		"custom_messages":           k.countCustomMessages(ctx),
	}
}

// countCustomMessages counts total custom messages
func (k *I18nKeeper) countCustomMessages(ctx sdk.Context) int {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.CustomMessagePrefix)
	defer iterator.Close()
	
	count := 0
	for ; iterator.Valid(); iterator.Next() {
		count++
	}
	
	return count
}

// ExportMessages exports all messages for translation management
func (k *I18nKeeper) ExportMessages(ctx sdk.Context) (*types.MessageCatalog, error) {
	// Load all custom messages into catalog
	customMessages, err := k.GetAllCustomMessages(ctx)
	if err != nil {
		return nil, err
	}
	
	exportCatalog := &types.MessageCatalog{
		Messages: make(map[string]*types.IdentityMessage),
	}
	
	// Copy default messages
	for key, msg := range k.catalog.Messages {
		exportCatalog.Messages[key] = msg
	}
	
	// Add custom messages
	for _, msg := range customMessages {
		exportCatalog.Messages[msg.Key] = msg
	}
	
	return exportCatalog, nil
}

// ImportMessages imports messages from external source
func (k *I18nKeeper) ImportMessages(ctx sdk.Context, catalog *types.MessageCatalog) error {
	for key, message := range catalog.Messages {
		// Only import custom messages (skip built-in ones)
		if !k.isBuiltInMessage(key) {
			if err := k.AddCustomMessage(ctx, message.Key, message.Category, message.Description, message.Text); err != nil {
				return err
			}
		}
	}
	
	return nil
}

// isBuiltInMessage checks if a message key is built-in
func (k *I18nKeeper) isBuiltInMessage(key string) bool {
	builtInKeys := []string{
		types.MsgAuthenticationSuccess,
		types.MsgAuthenticationFailed,
		types.MsgBiometricRequired,
		types.MsgIdentityCreated,
		types.MsgKYCVerificationSuccess,
		types.MsgConsentRequired,
		types.MsgWelcomeMessage,
		// Add other built-in message keys
	}
	
	for _, builtIn := range builtInKeys {
		if key == builtIn {
			return true
		}
	}
	
	return false
}