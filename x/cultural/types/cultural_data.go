package types

import (
	"fmt"
	"math/rand"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultQuotes returns a curated list of default quotes from Indian leaders
func DefaultQuotes() []Quote {
	quotes := []Quote{
		// Mahatma Gandhi Quotes
		{
			Id:                   1,
			Text:                 "Be the change you want to see in the world.",
			Author:               AuthorGandhi,
			Source:               "Various speeches and writings",
			Language:             LanguageEnglish,
			Category:             CategoryLeadership,
			Tags:                 []string{"change", "leadership", "personal development"},
			Region:               RegionGujarat,
			HistoricalPeriod:     "1900-1950",
			DifficultyLevel:      3,
			Translation:          map[string]string{"hindi": "वह परिवर्तन बनो जो तुम दुनिया में देखना चाहते हो।"},
			CreatedAt:            time.Now().Unix(),
			Verified:             true,
			CulturalSignificance: 10,
		},
		{
			Id:                   2,
			Text:                 "An eye for an eye only ends up making the whole world blind.",
			Author:               AuthorGandhi,
			Source:               "Young India, 1924",
			Language:             LanguageEnglish,
			Category:             CategoryPeace,
			Tags:                 []string{"peace", "non-violence", "ahimsa"},
			Region:               RegionGujarat,
			HistoricalPeriod:     "1900-1950",
			DifficultyLevel:      5,
			Translation:          map[string]string{"hindi": "आंख के बदले आंख पूरी दुनिया को अंधा बना देती है।"},
			CreatedAt:            time.Now().Unix(),
			Verified:             true,
			CulturalSignificance: 9,
		},
		{
			Id:                   3,
			Text:                 "Live as if you were to die tomorrow. Learn as if you were to live forever.",
			Author:               AuthorGandhi,
			Source:               "Harijan, 1947",
			Language:             LanguageEnglish,
			Category:             CategoryEducation,
			Tags:                 []string{"learning", "life", "education"},
			Region:               RegionGujarat,
			HistoricalPeriod:     "1900-1950",
			DifficultyLevel:      4,
			Translation:          map[string]string{"hindi": "ऐसे जिओ जैसे कि कल तुम मर जाओगे। ऐसे सीखो जैसे कि तुम हमेशा जीवित रहोगे।"},
			CreatedAt:            time.Now().Unix(),
			Verified:             true,
			CulturalSignificance: 8,
		},

		// Swami Vivekananda Quotes
		{
			Id:                   4,
			Text:                 "Arise, awake, and stop not until the goal is reached.",
			Author:               AuthorVivekananda,
			Source:               "Kathopanishad interpretation",
			Language:             LanguageEnglish,
			Category:             CategoryMotivation,
			Tags:                 []string{"motivation", "determination", "goals"},
			Region:               RegionBengal,
			HistoricalPeriod:     "1860-1900",
			DifficultyLevel:      6,
			Translation:          map[string]string{"hindi": "उठो, जागो और तब तक नहीं रुको जब तक लक्ष्य न मिल जाए।", "sanskrit": "उत्तिष्ठत जाग्रत प्राप्य वरान्निबोधत।"},
			CreatedAt:            time.Now().Unix(),
			Verified:             true,
			CulturalSignificance: 10,
		},
		{
			Id:                   5,
			Text:                 "You cannot believe in God until you believe in yourself.",
			Author:               AuthorVivekananda,
			Source:               "Complete Works of Swami Vivekananda",
			Language:             LanguageEnglish,
			Category:             CategorySpirituality,
			Tags:                 []string{"spirituality", "self-belief", "god"},
			Region:               RegionBengal,
			HistoricalPeriod:     "1860-1900",
			DifficultyLevel:      7,
			Translation:          map[string]string{"hindi": "तुम तब तक भगवान में विश्वास नहीं कर सकते जब तक तुम अपने आप में विश्वास नहीं करते।"},
			CreatedAt:            time.Now().Unix(),
			Verified:             true,
			CulturalSignificance: 9,
		},

		// Dr. APJ Abdul Kalam Quotes
		{
			Id:                   6,
			Text:                 "Dream is not that which you see while sleeping, it is something that does not let you sleep.",
			Author:               AuthorKalam,
			Source:               "Wings of Fire",
			Language:             LanguageEnglish,
			Category:             CategoryMotivation,
			Tags:                 []string{"dreams", "ambition", "success"},
			Region:               RegionTamil,
			HistoricalPeriod:     "1950-2015",
			DifficultyLevel:      5,
			Translation:          map[string]string{"hindi": "सपना वह नहीं है जो तुम सोते वक्त देखते हो, सपना वह है जो तुम्हें सोने नहीं देता।"},
			CreatedAt:            time.Now().Unix(),
			Verified:             true,
			CulturalSignificance: 8,
		},
		{
			Id:                   7,
			Text:                 "If you want to shine like a sun, first burn like a sun.",
			Author:               AuthorKalam,
			Source:               "Ignited Minds",
			Language:             LanguageEnglish,
			Category:             CategoryMotivation,
			Tags:                 []string{"hard work", "excellence", "dedication"},
			Region:               RegionTamil,
			HistoricalPeriod:     "1950-2015",
			DifficultyLevel:      4,
			Translation:          map[string]string{"hindi": "यदि तुम सूरज की तरह चमकना चाहते हो तो पहले सूरज की तरह जलना होगा।"},
			CreatedAt:            time.Now().Unix(),
			Verified:             true,
			CulturalSignificance: 8,
		},

		// Rabindranath Tagore Quotes
		{
			Id:                   8,
			Text:                 "Where the mind is without fear and the head is held high.",
			Author:               AuthorTagore,
			Source:               "Gitanjali",
			Language:             LanguageEnglish,
			Category:             CategoryPatriotism,
			Tags:                 []string{"freedom", "courage", "nation"},
			Region:               RegionBengal,
			HistoricalPeriod:     "1860-1940",
			DifficultyLevel:      6,
			Translation:          map[string]string{"hindi": "जहाँ मन भय से मुक्त है और सिर गर्व से ऊँचा है।"},
			CreatedAt:            time.Now().Unix(),
			Verified:             true,
			CulturalSignificance: 10,
		},
		{
			Id:                   9,
			Text:                 "You can't cross the sea merely by standing and staring at the water.",
			Author:               AuthorTagore,
			Source:               "Stray Birds",
			Language:             LanguageEnglish,
			Category:             CategoryMotivation,
			Tags:                 []string{"action", "achievement", "courage"},
			Region:               RegionBengal,
			HistoricalPeriod:     "1860-1940",
			DifficultyLevel:      5,
			Translation:          map[string]string{"hindi": "केवल खड़े होकर पानी को देखने से आप समुद्र पार नहीं कर सकते।"},
			CreatedAt:            time.Now().Unix(),
			Verified:             true,
			CulturalSignificance: 8,
		},

		// Buddha Quotes
		{
			Id:                   10,
			Text:                 "The mind is everything. What you think you become.",
			Author:               AuthorBuddha,
			Source:               "Dhammapada",
			Language:             LanguageEnglish,
			Category:             CategoryWisdom,
			Tags:                 []string{"mind", "thoughts", "self-development"},
			Region:               RegionBihar,
			HistoricalPeriod:     "500 BCE",
			DifficultyLevel:      7,
			Translation:          map[string]string{"hindi": "मन ही सब कुछ है। जो तुम सोचते हो, वही बन जाते हो।"},
			CreatedAt:            time.Now().Unix(),
			Verified:             true,
			CulturalSignificance: 10,
		},
		{
			Id:                   11,
			Text:                 "Three things cannot be long hidden: the sun, the moon, and the truth.",
			Author:               AuthorBuddha,
			Source:               "Buddhist Teachings",
			Language:             LanguageEnglish,
			Category:             CategoryTruth,
			Tags:                 []string{"truth", "honesty", "wisdom"},
			Region:               RegionBihar,
			HistoricalPeriod:     "500 BCE",
			DifficultyLevel:      5,
			Translation:          map[string]string{"hindi": "तीन चीजें लंबे समय तक छुपाई नहीं जा सकतीं: सूरज, चांद, और सत्य।"},
			CreatedAt:            time.Now().Unix(),
			Verified:             true,
			CulturalSignificance: 9,
		},

		// Jawaharlal Nehru Quotes
		{
			Id:                   12,
			Text:                 "The only way to deal with an unfree world is to become so absolutely free that your very existence is an act of rebellion.",
			Author:               AuthorNehru,
			Source:               "The Discovery of India",
			Language:             LanguageEnglish,
			Category:             CategoryPatriotism,
			Tags:                 []string{"freedom", "rebellion", "independence"},
			Region:               RegionUttar,
			HistoricalPeriod:     "1900-1964",
			DifficultyLevel:      8,
			Translation:          map[string]string{"hindi": "अस्वतंत्र दुनिया से निपटने का एकमात्र तरीका यह है कि बिल्कुल स्वतंत्र बन जाओ कि तुम्हारा अस्तित्व ही विद्रोह का कार्य बन जाए।"},
			CreatedAt:            time.Now().Unix(),
			Verified:             true,
			CulturalSignificance: 8,
		},

		// Bhagat Singh Quotes
		{
			Id:                   13,
			Text:                 "They may kill me, but they cannot kill my ideas. They can crush my body, but they will not be able to crush my spirit.",
			Author:               AuthorBhagat,
			Source:               "Letter from prison",
			Language:             LanguageEnglish,
			Category:             CategoryPatriotism,
			Tags:                 []string{"sacrifice", "courage", "patriotism"},
			Region:               RegionPunjab,
			HistoricalPeriod:     "1900-1931",
			DifficultyLevel:      6,
			Translation:          map[string]string{"hindi": "वे मुझे मार सकते हैं, लेकिन वे मेरे विचारों को नहीं मार सकते। वे मेरे शरीर को कुचल सकते हैं, लेकिन मेरी आत्मा को कुचल नहीं सकेंगे।"},
			CreatedAt:            time.Now().Unix(),
			Verified:             true,
			CulturalSignificance: 10,
		},

		// Subhash Chandra Bose Quotes
		{
			Id:                   14,
			Text:                 "Give me blood and I will give you freedom.",
			Author:               AuthorSubhash,
			Source:               "Independence movement speeches",
			Language:             LanguageEnglish,
			Category:             CategoryPatriotism,
			Tags:                 []string{"freedom", "sacrifice", "independence"},
			Region:               RegionBengal,
			HistoricalPeriod:     "1900-1945",
			DifficultyLevel:      5,
			Translation:          map[string]string{"hindi": "तुम मुझे खून दो और मैं तुम्हें आजादी दूंगा।"},
			CreatedAt:            time.Now().Unix(),
			Verified:             true,
			CulturalSignificance: 10,
		},

		// Sri Aurobindo Quotes
		{
			Id:                   15,
			Text:                 "All life is yoga.",
			Author:               AuthorAurobindo,
			Source:               "The Synthesis of Yoga",
			Language:             LanguageEnglish,
			Category:             CategorySpirituality,
			Tags:                 []string{"yoga", "life", "spirituality"},
			Region:               RegionBengal,
			HistoricalPeriod:     "1872-1950",
			DifficultyLevel:      7,
			Translation:          map[string]string{"hindi": "सारा जीवन योग है।"},
			CreatedAt:            time.Now().Unix(),
			Verified:             true,
			CulturalSignificance: 8,
		},
	}

	return quotes
}

// DefaultHistoricalEvents returns a curated list of significant Indian historical events
func DefaultHistoricalEvents() []HistoricalEvent {
	events := []HistoricalEvent{
		{
			Id:               1,
			Title:            "Independence Day",
			Description:      "India gained independence from British rule on August 15, 1947, marking the end of nearly 200 years of colonial rule.",
			Date:             "August 15, 1947",
			Year:             1947,
			Category:         EventIndependenceMovement,
			Location:         "New Delhi, India",
			KeyFigures:       []string{"Mahatma Gandhi", "Jawaharlal Nehru", "Sardar Patel"},
			Significance:     "Marked the birth of independent India and the beginning of a new era of self-governance.",
			Tags:             []string{"independence", "freedom", "nationalism"},
			CreatedAt:        time.Now().Unix(),
			Verified:         true,
			EducationalLevel: LevelPrimary,
		},
		{
			Id:               2,
			Title:            "Green Revolution",
			Description:      "The Green Revolution transformed Indian agriculture in the 1960s and 1970s, making India self-sufficient in food production.",
			Date:             "1960s-1970s",
			Year:             1960,
			Category:         EventGreenRevolution,
			Location:         "Punjab, Haryana, Western UP",
			KeyFigures:       []string{"M.S. Swaminathan", "Norman Borlaug", "Indira Gandhi"},
			Significance:     "Ended India's dependence on food imports and made the country self-sufficient in wheat and rice production.",
			Tags:             []string{"agriculture", "food security", "technology"},
			CreatedAt:        time.Now().Unix(),
			Verified:         true,
			EducationalLevel: LevelSecondary,
		},
		{
			Id:               3,
			Title:            "Chandrayaan-1 Mission",
			Description:      "India's first lunar mission, launched in 2008, confirmed the presence of water molecules on the moon.",
			Date:             "October 22, 2008",
			Year:             2008,
			Category:         EventSpaceAchievements,
			Location:         "Satish Dhawan Space Centre, Sriharikota",
			KeyFigures:       []string{"A.P.J. Abdul Kalam", "G. Madhavan Nair", "ISRO Team"},
			Significance:     "Established India as a major space-faring nation and contributed to lunar science globally.",
			Tags:             []string{"space", "moon", "science", "technology"},
			CreatedAt:        time.Now().Unix(),
			Verified:         true,
			EducationalLevel: LevelHighSchool,
		},
		{
			Id:               4,
			Title:            "Dandi March",
			Description:      "Gandhi's 240-mile march to Dandi to protest the British salt monopoly, starting the civil disobedience movement.",
			Date:             "March 12 - April 6, 1930",
			Year:             1930,
			Category:         EventIndependenceMovement,
			Location:         "Sabarmati Ashram to Dandi, Gujarat",
			KeyFigures:       []string{"Mahatma Gandhi", "Sarojini Naidu", "Abbas Tyabji"},
			Significance:     "Sparked nationwide civil disobedience and became a symbol of non-violent resistance.",
			Tags:             []string{"civil disobedience", "salt satyagraha", "non-violence"},
			CreatedAt:        time.Now().Unix(),
			Verified:         true,
			EducationalLevel: LevelPrimary,
		},
		{
			Id:               5,
			Title:            "Pokhran Nuclear Tests",
			Description:      "India conducted its first nuclear test in 1974 and became a nuclear power with the 1998 tests.",
			Date:             "May 18, 1974 and May 11, 1998",
			Year:             1998,
			Category:         EventDefenseAchievements,
			Location:         "Pokhran, Rajasthan",
			KeyFigures:       []string{"Atal Bihari Vajpayee", "Raja Ramanna", "A.P.J. Abdul Kalam"},
			Significance:     "Established India as a nuclear power and changed the strategic balance in South Asia.",
			Tags:             []string{"nuclear", "defense", "technology"},
			CreatedAt:        time.Now().Unix(),
			Verified:         true,
			EducationalLevel: LevelCollege,
		},
		{
			Id:               6,
			Title:            "India's First Cricket World Cup Victory",
			Description:      "India won their first Cricket World Cup in 1983, defeating the West Indies in the final.",
			Date:             "June 25, 1983",
			Year:             1983,
			Category:         EventSportsVictories,
			Location:         "Lord's Cricket Ground, London",
			KeyFigures:       []string{"Kapil Dev", "Mohinder Amarnath", "Krishnamachari Srikkanth"},
			Significance:     "Transformed cricket in India and inspired a generation of cricketers.",
			Tags:             []string{"cricket", "sports", "victory", "pride"},
			CreatedAt:        time.Now().Unix(),
			Verified:         true,
			EducationalLevel: LevelGeneral,
		},
		{
			Id:               7,
			Title:            "Constitution Day",
			Description:      "The Indian Constitution was adopted on November 26, 1949, establishing India as a sovereign democratic republic.",
			Date:             "November 26, 1949",
			Year:             1949,
			Category:         EventConstitutionalValues,
			Location:         "New Delhi, India",
			KeyFigures:       []string{"Dr. B.R. Ambedkar", "Rajendra Prasad", "Constituent Assembly"},
			Significance:     "Established the framework for Indian democracy and fundamental rights.",
			Tags:             []string{"constitution", "democracy", "rights", "governance"},
			CreatedAt:        time.Now().Unix(),
			Verified:         true,
			EducationalLevel: LevelSecondary,
		},
		{
			Id:               8,
			Title:            "Mangalyaan Mission",
			Description:      "India's Mars Orbiter Mission made India the first country to reach Mars orbit in its first attempt.",
			Date:             "September 24, 2014",
			Year:             2014,
			Category:         EventSpaceAchievements,
			Location:         "Mars Orbit",
			KeyFigures:       []string{"ISRO Team", "K. Radhakrishnan", "Narendra Modi"},
			Significance:     "Made India the first Asian nation to reach Mars and the first to do so in the first attempt.",
			Tags:             []string{"mars", "space", "achievement", "technology"},
			CreatedAt:        time.Now().Unix(),
			Verified:         true,
			EducationalLevel: LevelHighSchool,
		},
	}

	return events
}

// DefaultCulturalWisdom returns a curated list of cultural wisdom from various Indian traditions
func DefaultCulturalWisdom() []CulturalWisdom {
	wisdom := []CulturalWisdom{
		{
			Id:                    1,
			Text:                  "सर्वे भवन्तु सुखिनः सर्वे सन्तु निरामयाः।",
			Tradition:             TraditionSanskrit,
			Language:              LanguageSanskrit,
			Transliteration:       "Sarve bhavantu sukhinah sarve santu niramayah",
			Translation:           "May all beings be happy, may all beings be free from illness",
			Meaning:               "This is a universal prayer for the well-being of all living beings, expressing the wish for happiness and good health for everyone.",
			Context:               "This is a commonly recited Sanskrit prayer that embodies the principle of universal love and compassion.",
			Scripture:             SourceUpanishads,
			Tags:                  []string{"prayer", "well-being", "universal love", "compassion"},
			CreatedAt:             time.Now().Unix(),
			Verified:              true,
			SpiritualSignificance: 10,
		},
		{
			Id:                    2,
			Text:                  "वसुधैव कुटुम्बकम्।",
			Tradition:             TraditionSanskrit,
			Language:              LanguageSanskrit,
			Transliteration:       "Vasudhaiva Kutumbakam",
			Translation:           "The world is one family",
			Meaning:               "This profound concept teaches that all humanity is one extended family, transcending boundaries of race, religion, and nationality.",
			Context:               "This principle from the Upanishads is fundamental to Indian philosophy and has been adopted as a guiding principle for global harmony.",
			Scripture:             SourceUpanishads,
			VerseReference:        "Maha Upanishad 6.72",
			Tags:                  []string{"unity", "family", "global harmony", "philosophy"},
			CreatedAt:             time.Now().Unix(),
			Verified:              true,
			SpiritualSignificance: 10,
		},
		{
			Id:                    3,
			Text:                  "अहिंसा परमो धर्मः।",
			Tradition:             TraditionSanskrit,
			Language:              LanguageSanskrit,
			Transliteration:       "Ahimsa paramo dharmah",
			Translation:           "Non-violence is the highest virtue",
			Meaning:               "This principle establishes non-violence as the supreme moral law, applicable to thoughts, words, and actions.",
			Context:               "This is a fundamental principle in Jainism, Buddhism, and Hinduism, emphasizing the importance of non-violence in all aspects of life.",
			Scripture:             SourceMahabharata,
			Tags:                  []string{"non-violence", "dharma", "morality", "virtue"},
			CreatedAt:             time.Now().Unix(),
			Verified:              true,
			SpiritualSignificance: 9,
		},
		{
			Id:                    4,
			Text:                  "குறள் 380: அறத்துப் பால் - இல்லறவியல் - வாழ்க்கை துணைநலம்",
			Tradition:             TraditionTamil,
			Language:              LanguageTamil,
			Transliteration:       "Manathukkum uruvam illai arivukkum uruvam illai",
			Translation:           "The mind has no form, knowledge has no form",
			Meaning:               "This teaches that the most powerful things in life - the mind and knowledge - are formless and boundless.",
			Context:               "From Thirukkural, this verse emphasizes the power of the mind and knowledge over physical form.",
			Scripture:             SourceThirukkural,
			VerseReference:        "Kural 380",
			Tags:                  []string{"mind", "knowledge", "formless", "wisdom"},
			CreatedAt:             time.Now().Unix(),
			Verified:              true,
			SpiritualSignificance: 8,
		},
		{
			Id:                    5,
			Text:                  "ਸਤਿ ਨਾਮੁ ਕਰਤਾ ਪੁਰਖੁ ਨਿਰਭਉ ਨਿਰਵੈਰੁ",
			Tradition:             TraditionSikh,
			Language:              LanguagePunjabi,
			Transliteration:       "Sati naam karata purakhu nirbhau nirvair",
			Translation:           "True is the Name, the Creator, without fear, without hatred",
			Meaning:               "This is the opening of the Mool Mantra, describing the nature of the divine as truth, creator, fearless, and without enmity.",
			Context:               "This forms the beginning of the Guru Granth Sahib and encapsulates the fundamental beliefs of Sikhism.",
			Scripture:             SourceGuruGranthSahib,
			VerseReference:        "Mool Mantra",
			Tags:                  []string{"truth", "creator", "fearless", "divine"},
			CreatedAt:             time.Now().Unix(),
			Verified:              true,
			SpiritualSignificance: 10,
		},
		{
			Id:                    6,
			Text:                  "যে জন দিব্যজ্ঞান হীন, সে জন আন্ধ সমান",
			Tradition:             TraditionBengali,
			Language:              LanguageBengali,
			Transliteration:       "Je jon dibyagyan heen, se jon andho soman",
			Translation:           "One who lacks divine knowledge is like a blind person",
			Meaning:               "This emphasizes that spiritual wisdom is essential for truly seeing and understanding life.",
			Context:               "From Bengali spiritual literature, highlighting the importance of divine knowledge for spiritual sight.",
			Scripture:             "Bengali Spiritual Literature",
			Tags:                  []string{"knowledge", "wisdom", "spiritual sight", "divine"},
			CreatedAt:             time.Now().Unix(),
			Verified:              true,
			SpiritualSignificance: 8,
		},
		{
			Id:                    7,
			Text:                  "بسم الله الرحمن الرحيم",
			Tradition:             TraditionSufi,
			Language:              LanguageUrdu,
			Transliteration:       "Bismillah ir-Rahman ir-Rahim",
			Translation:           "In the name of Allah, the Most Gracious, the Most Merciful",
			Meaning:               "This invocation acknowledges the divine source of all actions and seeks blessings from the most compassionate creator.",
			Context:               "This is the opening verse of the Quran, recited by Muslims before beginning any task, emphasizing divine grace and mercy.",
			Scripture:             SourceQuran,
			VerseReference:        "Al-Fatiha 1:1",
			Tags:                  []string{"divine", "mercy", "grace", "blessing"},
			CreatedAt:             time.Now().Unix(),
			Verified:              true,
			SpiritualSignificance: 10,
		},
		{
			Id:                    8,
			Text:                  "योग: कर्मसु कौशलम्",
			Tradition:             TraditionSanskrit,
			Language:              LanguageSanskrit,
			Transliteration:       "Yogah karmasu kaushalam",
			Translation:           "Yoga is skill in action",
			Meaning:               "This defines yoga not just as physical postures but as skillful, mindful action in all aspects of life.",
			Context:               "From the Bhagavad Gita, this verse provides a practical definition of yoga as conscious, skillful living.",
			Scripture:             SourceBhagavadGita,
			VerseReference:        "Bhagavad Gita 2.50",
			Tags:                  []string{"yoga", "skill", "action", "mindfulness"},
			CreatedAt:             time.Now().Unix(),
			Verified:              true,
			SpiritualSignificance: 9,
		},
	}

	return wisdom
}

// GetQuoteByAmountRange returns quotes suitable for different transaction amounts
func GetQuoteByAmountRange(amount sdk.Int) []Quote {
	quotes := DefaultQuotes()
	
	// Simple amount-based filtering
	if amount.LTE(sdk.NewInt(1000)) {
		// Small amounts: motivational and simple quotes
		return filterQuotesByCategory(quotes, []string{CategoryMotivation, CategoryWisdom})
	} else if amount.LTE(sdk.NewInt(10000)) {
		// Medium amounts: philosophical and educational quotes
		return filterQuotesByCategory(quotes, []string{CategoryPhilosophy, CategoryEducation, CategoryLeadership})
	} else if amount.LTE(sdk.NewInt(100000)) {
		// Large amounts: patriotic and significant quotes
		return filterQuotesByCategory(quotes, []string{CategoryPatriotism, CategoryLeadership, CategorySpirituality})
	} else {
		// Very large amounts: most significant quotes
		return filterQuotesBySignificance(quotes, 8)
	}
}

// filterQuotesByCategory filters quotes by category
func filterQuotesByCategory(quotes []Quote, categories []string) []Quote {
	var filtered []Quote
	for _, quote := range quotes {
		for _, category := range categories {
			if quote.Category == category {
				filtered = append(filtered, quote)
				break
			}
		}
	}
	return filtered
}

// filterQuotesBySignificance filters quotes by cultural significance
func filterQuotesBySignificance(quotes []Quote, minSignificance int32) []Quote {
	var filtered []Quote
	for _, quote := range quotes {
		if quote.CulturalSignificance >= minSignificance {
			filtered = append(filtered, quote)
		}
	}
	return filtered
}

// GetRandomQuote returns a random quote from the provided list
func GetRandomQuote(quotes []Quote) Quote {
	if len(quotes) == 0 {
		return Quote{}
	}
	
	rand.Seed(time.Now().UnixNano())
	return quotes[rand.Intn(len(quotes))]
}

// GetSeasonalQuote returns a quote appropriate for the current season
func GetSeasonalQuote(quotes []Quote) Quote {
	// Simple implementation - can be enhanced with seasonal logic
	return GetRandomQuote(quotes)
}

// GetRegionalQuote returns a quote from a specific region
func GetRegionalQuote(quotes []Quote, region string) Quote {
	regionalQuotes := []Quote{}
	for _, quote := range quotes {
		if quote.Region == region {
			regionalQuotes = append(regionalQuotes, quote)
		}
	}
	
	if len(regionalQuotes) == 0 {
		return GetRandomQuote(quotes)
	}
	
	return GetRandomQuote(regionalQuotes)
}

// FormatQuoteForDisplay formats a quote for display in transaction
func FormatQuoteForDisplay(quote Quote, language string) string {
	if translation, exists := quote.Translation[language]; exists {
		return fmt.Sprintf("%s\n- %s", translation, quote.Author)
	}
	return fmt.Sprintf("%s\n- %s", quote.Text, quote.Author)
}

// GetDailyWisdom returns wisdom for daily display
func GetDailyWisdom(date time.Time) CulturalWisdom {
	wisdom := DefaultCulturalWisdom()
	
	// Use date to generate consistent daily wisdom
	dayOfYear := date.YearDay()
	index := dayOfYear % len(wisdom)
	
	return wisdom[index]
}

// ValidateQuote validates a quote for content and accuracy
func ValidateQuote(quote Quote) error {
	if len(quote.Text) == 0 {
		return fmt.Errorf("quote text cannot be empty")
	}
	
	if len(quote.Author) == 0 {
		return fmt.Errorf("quote author cannot be empty")
	}
	
	if quote.DifficultyLevel < 1 || quote.DifficultyLevel > 10 {
		return fmt.Errorf("difficulty level must be between 1 and 10")
	}
	
	if quote.CulturalSignificance < 1 || quote.CulturalSignificance > 10 {
		return fmt.Errorf("cultural significance must be between 1 and 10")
	}
	
	return nil
}

// GetQuoteStatistics returns statistics about quote usage
func GetQuoteStatistics(quotes []Quote) map[string]interface{} {
	stats := make(map[string]interface{})
	
	categoryCount := make(map[string]int)
	authorCount := make(map[string]int)
	totalUsage := uint64(0)
	
	for _, quote := range quotes {
		categoryCount[quote.Category]++
		authorCount[quote.Author]++
		totalUsage += quote.UsageCount
	}
	
	stats["total_quotes"] = len(quotes)
	stats["categories"] = categoryCount
	stats["authors"] = authorCount
	stats["total_usage"] = totalUsage
	
	return stats
}