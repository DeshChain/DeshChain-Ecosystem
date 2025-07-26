package types

// Built-in message keys for the identity system
const (
	// Authentication messages
	MsgAuthenticationSuccess = "auth.success"
	MsgAuthenticationFailed  = "auth.failed"
	MsgBiometricRequired     = "auth.biometric_required"
	MsgBiometricSuccess      = "auth.biometric_success"
	MsgBiometricFailed       = "auth.biometric_failed"
	MsgPasswordRequired      = "auth.password_required"
	MsgPasswordSuccess       = "auth.password_success"
	MsgPasswordFailed        = "auth.password_failed"
	MsgMFARequired           = "auth.mfa_required"
	MsgMFASuccess            = "auth.mfa_success"
	MsgMFAFailed             = "auth.mfa_failed"

	// Identity management messages
	MsgIdentityCreated       = "identity.created"
	MsgIdentityUpdated       = "identity.updated"
	MsgIdentityDeactivated   = "identity.deactivated"
	MsgIdentityNotFound      = "identity.not_found"
	MsgIdentityInactive      = "identity.inactive"
	MsgIdentityRevoked       = "identity.revoked"
	MsgIdentityRestored      = "identity.restored"
	MsgIdentityVerified      = "identity.verified"
	MsgIdentityVerificationFailed = "identity.verification_failed"

	// DID messages
	MsgDIDCreated            = "did.created"
	MsgDIDUpdated            = "did.updated"
	MsgDIDDeactivated        = "did.deactivated"
	MsgDIDNotFound           = "did.not_found"
	MsgDIDInvalid            = "did.invalid"
	MsgDIDDocumentInvalid    = "did.document_invalid"
	MsgDIDMethodNotSupported = "did.method_not_supported"

	// KYC messages
	MsgKYCRequired           = "kyc.required"
	MsgKYCInProgress         = "kyc.in_progress"
	MsgKYCCompleted          = "kyc.completed"
	MsgKYCVerificationSuccess = "kyc.verification_success"
	MsgKYCVerificationFailed = "kyc.verification_failed"
	MsgKYCExpired            = "kyc.expired"
	MsgKYCLevelInsufficient  = "kyc.level_insufficient"
	MsgKYCDataMismatch       = "kyc.data_mismatch"
	MsgKYCDocumentRequired   = "kyc.document_required"
	MsgKYCDocumentUploaded   = "kyc.document_uploaded"
	MsgKYCDocumentVerified   = "kyc.document_verified"
	MsgKYCDocumentRejected   = "kyc.document_rejected"

	// Credential messages
	MsgCredentialIssued      = "credential.issued"
	MsgCredentialVerified    = "credential.verified"
	MsgCredentialRevoked     = "credential.revoked"
	MsgCredentialExpired     = "credential.expired"
	MsgCredentialNotFound    = "credential.not_found"
	MsgCredentialInvalid     = "credential.invalid"
	MsgCredentialPresented   = "credential.presented"
	MsgCredentialAccepted    = "credential.accepted"
	MsgCredentialRejected    = "credential.rejected"

	// Consent messages
	MsgConsentRequired       = "consent.required"
	MsgConsentGiven          = "consent.given"
	MsgConsentWithdrawn      = "consent.withdrawn"
	MsgConsentExpired        = "consent.expired"
	MsgConsentNotFound       = "consent.not_found"
	MsgConsentInvalid        = "consent.invalid"
	MsgConsentUpdated        = "consent.updated"

	// Privacy messages
	MsgPrivacyEnabled        = "privacy.enabled"
	MsgPrivacyDisabled       = "privacy.disabled"
	MsgZKProofGenerated      = "privacy.zk_proof_generated"
	MsgZKProofVerified       = "privacy.zk_proof_verified"
	MsgZKProofFailed         = "privacy.zk_proof_failed"
	MsgAnonymousMode         = "privacy.anonymous_mode"
	MsgSelectiveDisclosure   = "privacy.selective_disclosure"

	// Recovery messages
	MsgRecoveryInitiated     = "recovery.initiated"
	MsgRecoveryCompleted     = "recovery.completed"
	MsgRecoveryFailed        = "recovery.failed"
	MsgRecoveryMethodSet     = "recovery.method_set"
	MsgRecoveryMethodUpdated = "recovery.method_updated"
	MsgRecoveryCodeGenerated = "recovery.code_generated"
	MsgGuardianAdded         = "recovery.guardian_added"
	MsgGuardianRemoved       = "recovery.guardian_removed"
	MsgGuardianApprovalRequired = "recovery.guardian_approval_required"

	// India Stack messages
	MsgAadhaarLinked         = "india_stack.aadhaar_linked"
	MsgAadhaarVerified       = "india_stack.aadhaar_verified"
	MsgAadhaarLinkFailed     = "india_stack.aadhaar_link_failed"
	MsgDigiLockerConnected   = "india_stack.digilocker_connected"
	MsgDigiLockerDocumentFetched = "india_stack.digilocker_document_fetched"
	MsgUPILinked             = "india_stack.upi_linked"
	MsgDEPAConsentCreated    = "india_stack.depa_consent_created"
	MsgAAConsentGranted      = "india_stack.aa_consent_granted"
	MsgPanchayatKYCCompleted = "india_stack.panchayat_kyc_completed"

	// General system messages
	MsgWelcomeMessage        = "system.welcome"
	MsgGoodbye              = "system.goodbye"
	MsgOperationSuccess      = "system.operation_success"
	MsgOperationFailed       = "system.operation_failed"
	MsgMaintenanceMode       = "system.maintenance_mode"
	MsgServiceUnavailable    = "system.service_unavailable"
	MsgRateLimitExceeded     = "system.rate_limit_exceeded"
	MsgInvalidRequest        = "system.invalid_request"
	MsgUnauthorized          = "system.unauthorized"
	MsgForbidden             = "system.forbidden"
	MsgNotFound              = "system.not_found"
	MsgInternalError         = "system.internal_error"

	// Internationalization system messages
	MsgLanguageChanged       = "i18n.language_changed"
	MsgLanguageNotSupported  = "i18n.language_not_supported"
	MsgTranslationMissing    = "i18n.translation_missing"
	MsgLocalizationUpdated   = "i18n.localization_updated"
	MsgCustomMessageAdded    = "i18n.custom_message_added"
	MsgMessagesImported      = "i18n.messages_imported"

	// Cultural and festival messages
	MsgFestivalGreeting      = "culture.festival_greeting"
	MsgTimeBasedGreeting     = "culture.time_based_greeting"
	MsgCulturalQuote         = "culture.quote"
	MsgRegionalWelcome       = "culture.regional_welcome"
	MsgPatrioticMessage      = "culture.patriotic_message"

	// Notification messages
	MsgNotificationSent      = "notification.sent"
	MsgNotificationDelivered = "notification.delivered"
	MsgNotificationFailed    = "notification.failed"
	MsgEmailSent             = "notification.email_sent"
	MsgSMSSent               = "notification.sms_sent"
	MsgPushNotificationSent  = "notification.push_sent"
)

// Default built-in messages in multiple languages
func GetBuiltInMessages() map[string]*IdentityMessage {
	return map[string]*IdentityMessage{
		MsgWelcomeMessage: {
			Key:      MsgWelcomeMessage,
			Category: "system",
			Description: "Welcome message for new users",
			Text: &LocalizedText{
				DefaultLang: LanguageEnglish,
				Translations: map[LanguageCode]string{
					LanguageEnglish:  "Welcome to DeshChain Identity! Your secure, privacy-preserving digital identity platform.",
					LanguageHindi:    "देशचेन आइडेंटिटी में आपका स्वागत है! आपका सुरक्षित, गोपनीयता-संरक्षित डिजिटल पहचान प्लेटफॉर्म।",
					LanguageBengali:  "দেশচেইন আইডেন্টিটিতে স্বাগতম! আপনার নিরাপদ, গোপনীয়তা-সংরক্ষণকারী ডিজিটাল পরিচয় প্ল্যাটফর্ম।",
					LanguageTamil:    "தேஷ்சேன் அடையாளத்திற்கு வரவேற்கிறோம்! உங்கள் பாதுகாப்பான, தனியுரிமை-பாதுகாக்கும் டிஜிட்டல் அடையாள தளம்।",
					LanguageTelugu:   "దేష్‌చైన్ ఐడెంటిటీకి స్వాగతం! మీ సురక్షితమైన, గోప్యత-రక్షిత డిజిటల్ గుర్తింపు వేదిక।",
					LanguageMarathi:  "देशचेन आयडेंटिटीमध्ये तुमचे स्वागत आहे! तुमचे सुरक्षित, गोपनीयता-संरक्षित डिजिटल ओळख प्लॅटफॉर्म।",
					LanguageGujarati: "દેશચેન આઇડેન્ટિટીમાં તમારું સ્વાગત છે! તમારું સુરક્ષિત, ગોપનીયતા-સંરક્ષિત ડિજિટલ ઓળખ પ્લેટફોર્મ।",
					LanguageUrdu:     "دیش چین آئیڈنٹٹی میں آپ کا خوش آمدید! آپ کا محفوظ، رازداری کو برقرار رکھنے والا ڈیجیٹل شناختی پلیٹ فارم۔",
				},
			},
		},
		MsgAuthenticationSuccess: {
			Key:      MsgAuthenticationSuccess,
			Category: "authentication",
			Description: "Authentication successful message",
			Text: &LocalizedText{
				DefaultLang: LanguageEnglish,
				Translations: map[LanguageCode]string{
					LanguageEnglish:  "Authentication successful! Welcome back.",
					LanguageHindi:    "प्रमाणीकरण सफल! आपका स्वागत है।",
					LanguageBengali:  "প্রমাণীকরণ সফল! স্বাগতম।",
					LanguageTamil:    "அங்கீகாரம் வெற்றிகரமாக! மீண்டும் வரவேற்கிறோம்।",
					LanguageTelugu:   "ప్రమాణీకరణ విజయవంతం! మీకు స్వాగతం।",
					LanguageMarathi:  "प्रमाणीकरण यशस्वी! तुमचे स्वागत आहे।",
					LanguageGujarati: "પ્રમાણીકરણ સફળ! તમારું સ્વાગત છે।",
					LanguageUrdu:     "تصدیق کامیاب! خوش آمدید۔",
				},
			},
		},
		MsgKYCCompleted: {
			Key:      MsgKYCCompleted,
			Category: "kyc",
			Description: "KYC completion message",
			Text: &LocalizedText{
				DefaultLang: LanguageEnglish,
				Translations: map[LanguageCode]string{
					LanguageEnglish:  "KYC verification completed successfully! Your identity is now verified.",
					LanguageHindi:    "केवाईसी सत्यापन सफलतापूर्वक पूरा हुआ! आपकी पहचान अब सत्यापित है।",
					LanguageBengali:  "কেওয়াইসি যাচাইকরণ সফলভাবে সম্পন্ন! আপনার পরিচয় এখন যাচাই করা হয়েছে।",
					LanguageTamil:    "KYC சரிபார்ப்பு வெற்றிகரமாக முடிந்தது! உங்கள் அடையாளம் இப்போது சரிபார்க்கப்பட்டது।",
					LanguageTelugu:   "KYC ధృవీకరణ విజయవంతంగా పూర్తయింది! మీ గుర్తింపు ఇప్పుడు ధృవీకరించబడింది।",
					LanguageMarathi:  "केवायसी पडताळणी यशस्वीरित्या पूर्ण झाली! तुमची ओळख आता पडताळली गेली आहे।",
					LanguageGujarati: "KYC ચકાસણી સફળતાપૂર્વક પૂર્ણ! તમારી ઓળખ હવે ચકાસાયેલ છે।",
					LanguageUrdu:     "KYC تصدیق کامیابی سے مکمل! آپ کی شناخت اب تصدیق شدہ ہے۔",
				},
			},
		},
		MsgCredentialIssued: {
			Key:      MsgCredentialIssued,
			Category: "credentials",
			Description: "Credential issued message",
			Text: &LocalizedText{
				DefaultLang: LanguageEnglish,
				Translations: map[LanguageCode]string{
					LanguageEnglish:  "New credential has been issued to your identity wallet.",
					LanguageHindi:    "आपके पहचान वॉलेट में नया क्रेडेंशियल जारी किया गया है।",
					LanguageBengali:  "আপনার পরিচয় ওয়ালেটে নতুন শংসাপত্র জারি করা হয়েছে।",
					LanguageTamil:    "உங்கள் அடையாள வால்லெட்டில் புதிய சான்றிதழ் வழங்கப்பட்டுள்ளது।",
					LanguageTelugu:   "మీ గుర్తింపు వాలెట్‌కు కొత్త క్రెడెన్షియల్ జారీ చేయబడింది।",
					LanguageMarathi:  "तुमच्या ओळख वॉलेटमध्ये नवीन क्रेडेंशियल जारी केले गेले आहे।",
					LanguageGujarati: "તમારા ઓળખ વૉલેટમાં નવું ક્રેડેન્શિયલ જારી કરવામાં આવ્યું છે।",
					LanguageUrdu:     "آپ کے شناختی والیٹ میں نیا کریڈنشل جاری کیا گیا ہے۔",
				},
			},
		},
		MsgConsentRequired: {
			Key:      MsgConsentRequired,
			Category: "consent",
			Description: "Consent required message",
			Text: &LocalizedText{
				DefaultLang: LanguageEnglish,
				Translations: map[LanguageCode]string{
					LanguageEnglish:  "Your consent is required to proceed with this operation.",
					LanguageHindi:    "इस ऑपरेशन को आगे बढ़ाने के लिए आपकी सहमति आवश्यक है।",
					LanguageBengali:  "এই অপারেশনটি এগিয়ে নিতে আপনার সম্মতি প্রয়োজন।",
					LanguageTamil:    "இந்த செயல்பாட்டைத் தொடர உங்கள் ஒப்புதல் தேவை।",
					LanguageTelugu:   "ఈ ఆపరేషన్‌ను కొనసాగించడానికి మీ సమ్మతి అవసరం।",
					LanguageMarathi:  "हे ऑपरेशन पुढे चालू ठेवण्यासाठी तुमची संमती आवश्यक आहे।",
					LanguageGujarati: "આ ઓપરેશન આગળ વધારવા માટે તમારી સંમતિ જરૂરી છે।",
					LanguageUrdu:     "اس آپریشن کو آگے بڑھانے کے لیے آپ کی رضامندی درکار ہے۔",
				},
			},
		},
		MsgOperationSuccess: {
			Key:      MsgOperationSuccess,
			Category: "system",
			Description: "Generic operation success message",
			Text: &LocalizedText{
				DefaultLang: LanguageEnglish,
				Translations: map[LanguageCode]string{
					LanguageEnglish:  "Operation completed successfully!",
					LanguageHindi:    "ऑपरेशन सफलतापूर्वक पूरा हुआ!",
					LanguageBengali:  "অপারেশন সফলভাবে সম্পন্ন!",
					LanguageTamil:    "செயல்பாடு வெற்றிகரமாக முடிந்தது!",
					LanguageTelugu:   "ఆపరేషన్ విజయవంతంగా పూర్తయింది!",
					LanguageMarathi:  "ऑपरेशन यशस्वीरित्या पूर्ण झाले!",
					LanguageGujarati: "ઓપરેશન સફળતાપૂર્વક પૂર્ણ!",
					LanguageUrdu:     "آپریشن کامیابی سے مکمل!",
				},
			},
		},
		MsgInvalidRequest: {
			Key:      MsgInvalidRequest,
			Category: "error",
			Description: "Invalid request error message",
			Text: &LocalizedText{
				DefaultLang: LanguageEnglish,
				Translations: map[LanguageCode]string{
					LanguageEnglish:  "Invalid request. Please check your input and try again.",
					LanguageHindi:    "अमान्य अनुरोध। कृपया अपना इनपुट जांचें और फिर से कोशिश करें।",
					LanguageBengali:  "অবৈধ অনুরোধ। দয়া করে আপনার ইনপুট পরীক্ষা করুন এবং আবার চেষ্টা করুন।",
					LanguageTamil:    "தவறான கோரிக்கை. உங்கள் உள்ளீட்டைச் சரிபார்த்து மீண்டும் முயற்சிக்கவும்।",
					LanguageTelugu:   "చెల్లని అభ్యర్థన. దయచేసి మీ ఇన్‌పుట్‌ను తనిఖీ చేసి మళ్లీ ప్రయత్నించండి।",
					LanguageMarathi:  "अवैध विनंती. कृपया तुमचे इनपुट तपासा आणि पुन्हा प्रयत्न करा।",
					LanguageGujarati: "અમાન્ય વિનંતી. કૃપા કરીને તમારું ઇનપુટ તપાસો અને ફરીથી પ્રયાસ કરો।",
					LanguageUrdu:     "غلط درخواست۔ براہ کرم اپنا ان پٹ چیک کریں اور دوبارہ کوشش کریں۔",
				},
			},
		},
		MsgLanguageChanged: {
			Key:      MsgLanguageChanged,
			Category: "i18n",
			Description: "Language preference changed message",
			Text: &LocalizedText{
				DefaultLang: LanguageEnglish,
				Translations: map[LanguageCode]string{
					LanguageEnglish:  "Language preference updated successfully!",
					LanguageHindi:    "भाषा वरीयता सफलतापूर्वक अपडेट की गई!",
					LanguageBengali:  "ভাষার পছন্দ সফলভাবে আপডেট করা হয়েছে!",
					LanguageTamil:    "மொழி விருப்பம் வெற்றிகரமாக புதுப்பிக்கப்பட்டது!",
					LanguageTelugu:   "భాష ప్రాధాన్యత విజయవంతంగా అప్‌డేట్ చేయబడింది!",
					LanguageMarathi:  "भाषा प्राधान्य यशस्वीरित्या अपडेट केले!",
					LanguageGujarati: "ભાષા પસંદગી સફળતાપૂર્વક અપડેટ!",
					LanguageUrdu:     "زبان کی ترجیح کامیابی سے اپڈیٹ!",
				},
			},
		},
	}
}

// Quote categories for cultural content
const (
	QuoteCategoryWisdom     = "wisdom"
	QuoteCategoryMotivation = "motivation"
	QuoteCategoryPatriotism = "patriotism"
	QuoteCategoryTechnology = "technology"
	QuoteCategoryGeneral    = "general"
)

// Time of day constants for greetings
const (
	TimeOfDayMorning   = "morning"
	TimeOfDayAfternoon = "afternoon"
	TimeOfDayEvening   = "evening"
	TimeOfDayNight     = "night"
)

// Festival constants for cultural greetings
const (
	FestivalDiwali         = "diwali"
	FestivalHoli           = "holi"
	FestivalDussehra       = "dussehra"
	FestivalGaneshChaturthi = "ganesh_chaturthi"
	FestivalEid            = "eid"
	FestivalChristmas      = "christmas"
	FestivalNewYear        = "new_year"
	FestivalRepublicDay    = "republic_day"
	FestivalIndependenceDay = "independence_day"
)

// Message categories
const (
	CategorySystem         = "system"
	CategoryAuthentication = "authentication"
	CategoryIdentity       = "identity"
	CategoryDID            = "did"
	CategoryKYC            = "kyc"
	CategoryCredentials    = "credentials"
	CategoryConsent        = "consent"
	CategoryPrivacy        = "privacy"
	CategoryRecovery       = "recovery"
	CategoryIndiaStack     = "india_stack"
	CategoryCulture        = "culture"
	CategoryNotification   = "notification"
	CategoryI18n           = "i18n"
	CategoryError          = "error"
)

// GetMessageCategories returns all available message categories
func GetMessageCategories() []string {
	return []string{
		CategorySystem,
		CategoryAuthentication,
		CategoryIdentity,
		CategoryDID,
		CategoryKYC,
		CategoryCredentials,
		CategoryConsent,
		CategoryPrivacy,
		CategoryRecovery,
		CategoryIndiaStack,
		CategoryCulture,
		CategoryNotification,
		CategoryI18n,
		CategoryError,
	}
}

// GetQuoteCategories returns all available quote categories
func GetQuoteCategories() []string {
	return []string{
		QuoteCategoryWisdom,
		QuoteCategoryMotivation,
		QuoteCategoryPatriotism,
		QuoteCategoryTechnology,
		QuoteCategoryGeneral,
	}
}

// GetTimeOfDayOptions returns all available time of day options
func GetTimeOfDayOptions() []string {
	return []string{
		TimeOfDayMorning,
		TimeOfDayAfternoon,
		TimeOfDayEvening,
		TimeOfDayNight,
	}
}

// GetFestivalOptions returns all available festival options
func GetFestivalOptions() []string {
	return []string{
		FestivalDiwali,
		FestivalHoli,
		FestivalDussehra,
		FestivalGaneshChaturthi,
		FestivalEid,
		FestivalChristmas,
		FestivalNewYear,
		FestivalRepublicDay,
		FestivalIndependenceDay,
	}
}