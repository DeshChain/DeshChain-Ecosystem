package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "cultural"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_cultural"
)

var (
	// ParamsKey is the key for module parameters
	ParamsKey = collections.NewPrefix(0)

	// QuoteKey is the key for storing quotes
	QuoteKey = collections.NewPrefix(1)

	// QuoteCountKey is the key for quote counter
	QuoteCountKey = collections.NewPrefix(2)

	// HistoricalEventKey is the key for storing historical events
	HistoricalEventKey = collections.NewPrefix(3)

	// HistoricalEventCountKey is the key for historical event counter
	HistoricalEventCountKey = collections.NewPrefix(4)

	// CulturalWisdomKey is the key for storing cultural wisdom
	CulturalWisdomKey = collections.NewPrefix(5)

	// CulturalWisdomCountKey is the key for cultural wisdom counter
	CulturalWisdomCountKey = collections.NewPrefix(6)

	// TransactionQuoteKey is the key for transaction quotes
	TransactionQuoteKey = collections.NewPrefix(7)

	// UserFavoriteKey is the key for user favorites
	UserFavoriteKey = collections.NewPrefix(8)

	// DailyWisdomKey is the key for daily wisdom
	DailyWisdomKey = collections.NewPrefix(9)

	// QuoteStatisticsKey is the key for quote statistics
	QuoteStatisticsKey = collections.NewPrefix(10)

	// ContentModerationKey is the key for content moderation
	ContentModerationKey = collections.NewPrefix(11)
)

// Quote Categories
const (
	CategoryLeadership      = "leadership"
	CategoryPhilosophy      = "philosophy"
	CategorySpirituality    = "spirituality"
	CategoryPatriotism      = "patriotism"
	CategoryScience         = "science"
	CategoryEducation       = "education"
	CategoryUnity           = "unity"
	CategoryPeace           = "peace"
	CategoryWisdom          = "wisdom"
	CategoryMotivation      = "motivation"
	CategoryLife            = "life"
	CategoryTruth           = "truth"
	CategoryDharma          = "dharma"
	CategoryKarma           = "karma"
	CategoryAhimsa          = "ahimsa"
	CategoryService         = "service"
	CategoryKnowledge       = "knowledge"
	CategoryHumanity        = "humanity"
	CategoryNature          = "nature"
	CategoryProgress        = "progress"
)

// Famous Indian Leaders and Authors
const (
	AuthorGandhi         = "Mahatma Gandhi"
	AuthorVivekananda    = "Swami Vivekananda"
	AuthorKalam          = "Dr. APJ Abdul Kalam"
	AuthorTagore         = "Rabindranath Tagore"
	AuthorBuddha         = "Gautama Buddha"
	AuthorNehru          = "Jawaharlal Nehru"
	AuthorSaraswati      = "Dayananda Saraswati"
	AuthorAurobindo      = "Sri Aurobindo"
	AuthorRaman          = "Sir C.V. Raman"
	AuthorBhagat         = "Bhagat Singh"
	AuthorSubhash        = "Subhash Chandra Bose"
	AuthorRoy            = "Ram Mohan Roy"
	AuthorTilak          = "Bal Gangadhar Tilak"
	AuthorAzad           = "Maulana Abul Kalam Azad"
	AuthorNaidu          = "Sarojini Naidu"
	AuthorRoy2           = "Aruna Roy"
	AuthorKiran          = "Kiran Bedi"
	AuthorYunus          = "Muhammad Yunus"
	AuthorSharma         = "Pandit Ravi Shankar"
	AuthorSadhguru       = "Sadhguru Jaggi Vasudev"
)

// Cultural Traditions
const (
	TraditionVedic         = "vedic"
	TraditionBuddhist      = "buddhist"
	TraditionSikh          = "sikh"
	TraditionSufi          = "sufi"
	TraditionTamil         = "tamil"
	TraditionSanskrit      = "sanskrit"
	TraditionHindi         = "hindi"
	TraditionUrdu          = "urdu"
	TraditionBengali       = "bengali"
	TraditionMarathi       = "marathi"
	TraditionGujarati      = "gujarati"
	TraditionTelugu        = "telugu"
	TraditionKannada       = "kannada"
	TraditionMalayalam     = "malayalam"
	TraditionPunjabi       = "punjabi"
	TraditionAssamese      = "assamese"
	TraditionOdiya         = "odiya"
)

// Historical Event Categories
const (
	EventIndependenceMovement = "independence_movement"
	EventGreenRevolution     = "green_revolution"
	EventWhiteRevolution     = "white_revolution"
	EventSpaceAchievements   = "space_achievements"
	EventSportsVictories     = "sports_victories"
	EventCulturalHeritage    = "cultural_heritage"
	EventScientificBreakthroughs = "scientific_breakthroughs"
	EventSocialReform        = "social_reform"
	EventEconomicMilestones  = "economic_milestones"
	EventConstitutionalValues = "constitutional_values"
	EventDefenseAchievements = "defense_achievements"
	EventEducationalProgress = "educational_progress"
	EventHealthcareAdvances  = "healthcare_advances"
	EventTechnologyInnovation = "technology_innovation"
	EventEnvironmentalConservation = "environmental_conservation"
)

// Sacred Texts and Sources
const (
	SourceVedas            = "vedas"
	SourceUpanishads       = "upanishads"
	SourceBhagavadGita     = "bhagavad_gita"
	SourceRamayana         = "ramayana"
	SourceMahabharata      = "mahabharata"
	SourceThirukkural      = "thirukkural"
	SourceGuruGranthSahib  = "guru_granth_sahib"
	SourceTripitaka        = "tripitaka"
	SourceQuran            = "quran"
	SourceConstitution     = "constitution"
	SourcePuranas          = "puranas"
	SourceYogaSutras       = "yoga_sutras"
	SourceArthashastra     = "arthashastra"
	SourceKamasutra        = "kamasutra"
	SourceAyurveda         = "ayurveda"
	SourceSangamLiterature = "sangam_literature"
	SourceBuddhist         = "buddhist_texts"
	SourceJain             = "jain_texts"
	SourceSikh             = "sikh_texts"
	SourceSufi             = "sufi_texts"
)

// Indian Regions
const (
	RegionNorth     = "north"
	RegionSouth     = "south"
	RegionEast      = "east"
	RegionWest      = "west"
	RegionCentral   = "central"
	RegionNortheast = "northeast"
	RegionKashmir   = "kashmir"
	RegionRajasthan = "rajasthan"
	RegionGujarat   = "gujarat"
	RegionPunjab    = "punjab"
	RegionBengal    = "bengal"
	RegionTamil     = "tamil"
	RegionKerala    = "kerala"
	RegionKarnataka = "karnataka"
	RegionTelangana = "telangana"
	RegionMaharashtra = "maharashtra"
	RegionMadhya    = "madhya_pradesh"
	RegionUttar     = "uttar_pradesh"
	RegionBihar     = "bihar"
	RegionOdisha    = "odisha"
	RegionAssam     = "assam"
	RegionManipur   = "manipur"
	RegionNagaland  = "nagaland"
	RegionMizoram   = "mizoram"
	RegionTripura   = "tripura"
	RegionSikkim    = "sikkim"
	RegionArunachal = "arunachal_pradesh"
	RegionMeghalaya = "meghalaya"
	RegionGoa       = "goa"
	RegionHimachal  = "himachal_pradesh"
	RegionUttarakhand = "uttarakhand"
	RegionJharkhand = "jharkhand"
	RegionChhattisgarh = "chhattisgarh"
	RegionDelhi     = "delhi"
)

// Languages
const (
	LanguageHindi     = "hindi"
	LanguageEnglish   = "english"
	LanguageSanskrit  = "sanskrit"
	LanguageTamil     = "tamil"
	LanguageBengali   = "bengali"
	LanguageTelugu    = "telugu"
	LanguageMarathi   = "marathi"
	LanguageGujarati  = "gujarati"
	LanguageKannada   = "kannada"
	LanguageMalayalam = "malayalam"
	LanguagePunjabi   = "punjabi"
	LanguageUrdu      = "urdu"
	LanguageAssamese  = "assamese"
	LanguageOdiya     = "odiya"
	LanguageKashmiri  = "kashmiri"
	LanguageNepali    = "nepali"
	LanguageSindhi    = "sindhi"
	LanguageMaithili  = "maithili"
	LanguageSantali   = "santali"
	LanguageBodo      = "bodo"
)

// Educational Levels
const (
	LevelPrimary    = "primary"
	LevelSecondary  = "secondary"
	LevelHighSchool = "high_school"
	LevelCollege    = "college"
	LevelUniversity = "university"
	LevelResearch   = "research"
	LevelGeneral    = "general"
)

// Quote Selection Algorithms
const (
	AlgorithmRandom         = "random"
	AlgorithmAmountBased    = "amount_based"
	AlgorithmTimeBased      = "time_based"
	AlgorithmRegionBased    = "region_based"
	AlgorithmUserProfile    = "user_profile"
	AlgorithmSeasonal       = "seasonal"
	AlgorithmFestival       = "festival"
	AlgorithmTrending       = "trending"
	AlgorithmPersonalized   = "personalized"
	AlgorithmEducational    = "educational"
	AlgorithmMotivational   = "motivational"
	AlgorithmCultural       = "cultural"
	AlgorithmSpiritual      = "spiritual"
	AlgorithmHistorical     = "historical"
	AlgorithmInspirational  = "inspirational"
)

// Event types for cultural module
const (
	EventTypeQuoteUsed         = "quote_used"
	EventTypeEventShared       = "event_shared"
	EventTypeWisdomDisplayed   = "wisdom_displayed"
	EventTypeNFTCreated        = "nft_created"
	EventTypeUserRating        = "user_rating"
	EventTypeContentAdded     = "content_added"
	EventTypeContentModerated = "content_moderated"
	EventTypeDailyWisdom      = "daily_wisdom"
	EventTypeFavoriteAdded    = "favorite_added"
	EventTypeQuoteShared      = "quote_shared"
)

// Default Values
const (
	DefaultMaxQuotesPerTransaction = 1
	DefaultNFTCreationFee          = 1000000 // 1 NAMO
	DefaultDifficultyLevel         = 5
	DefaultSignificanceThreshold   = 7
	DefaultMaxQuotesInDatabase     = 100000
	DefaultMaxEventsInDatabase     = 10000
	DefaultMaxWisdomInDatabase     = 50000
)