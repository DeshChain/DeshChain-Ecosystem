/*
Copyright 2024 DeshChain Foundation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package types

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// BollywoodDialogues - Epic movie dialogues adapted for coding
var BollywoodDialogues = []HumorQuote{
	// Amitabh Bachchan Style
	{
		QuoteId:             1,
		QuoteType:           QuoteType_QUOTE_TYPE_BOLLYWOOD_DIALOGUE,
		Text:                "Rishtey mein toh hum tumhare DEBUGGER lagte hain, naam hai BUG BUSTER!",
		EnglishTranslation:  "In relationships, I'm your DEBUGGER, name is BUG BUSTER!",
		Source:              "Kabhi Khushi Kabhie Bug",
		Character:           "Amitabh Debugchan",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_BUG_FIXES},
		IsFamilyFriendly:    true,
	},
	{
		QuoteId:             2,
		QuoteType:           QuoteType_QUOTE_TYPE_BOLLYWOOD_DIALOGUE,
		Text:                "Jahan commit hoti hai wahin se coding shuru hoti hai",
		EnglishTranslation:  "Where there's a commit, that's where coding begins",
		Source:              "Agnipath",
		Character:           "Big B",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_COMMITS},
		IsFamilyFriendly:    true,
	},

	// Shah Rukh Khan Style
	{
		QuoteId:             3,
		QuoteType:           QuoteType_QUOTE_TYPE_BOLLYWOOD_DIALOGUE,
		Text:                "Bade bade deshon mein aisi choti choti bugs hoti rehti hai",
		EnglishTranslation:  "In big countries, such small bugs keep happening",
		Source:              "DDLJ - Dilwale Debugging Le Jayenge",
		Character:           "SRK",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_BUG_FIXES},
		IsFamilyFriendly:    true,
	},
	{
		QuoteId:             4,
		QuoteType:           QuoteType_QUOTE_TYPE_BOLLYWOOD_DIALOGUE,
		Text:                "Feature nahi banana... developer ban jaana hai",
		EnglishTranslation:  "Don't make features... become the developer",
		Source:              "Feature Naam Hai Mera",
		Character:           "Feature Khan",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_FEATURES},
		IsFamilyFriendly:    true,
	},
	{
		QuoteId:             5,
		QuoteType:           QuoteType_QUOTE_TYPE_BOLLYWOOD_DIALOGUE,
		Text:                "Main GitHub pe 70 repo maintain karta hun... 70!",
		EnglishTranslation:  "I maintain 70 repos on GitHub... 70!",
		Source:              "Chak De Code",
		Character:           "Coach SRK",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_SPECIAL},
		IsFamilyFriendly:    true,
	},

	// Salman Khan Style
	{
		QuoteId:             6,
		QuoteType:           QuoteType_QUOTE_TYPE_BOLLYWOOD_DIALOGUE,
		Text:                "Ek baar jo maine commitment kar di, phir exception bhi nahi sunta",
		EnglishTranslation:  "Once I make a commitment, I don't even listen to exceptions",
		Source:              "Wanted: Bug Free Code",
		Character:           "Salman Bhai",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_COMMITS},
		IsFamilyFriendly:    true,
	},
	{
		QuoteId:             7,
		QuoteType:           QuoteType_QUOTE_TYPE_BOLLYWOOD_DIALOGUE,
		Text:                "Code fast karo, lekin dil se karo",
		EnglishTranslation:  "Code fast, but code from the heart",
		Source:              "Dabangg Developer",
		Character:           "Speed Sultan",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_PERFORMANCE},
		IsFamilyFriendly:    true,
	},

	// Aamir Khan Style
	{
		QuoteId:             8,
		QuoteType:           QuoteType_QUOTE_TYPE_BOLLYWOOD_DIALOGUE,
		Text:                "All izz well... bas console.log() kar do!",
		EnglishTranslation:  "All is well... just do console.log()!",
		Source:              "3 Idiots of Programming",
		Character:           "Rancho",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_BUG_FIXES},
		IsFamilyFriendly:    true,
	},
	{
		QuoteId:             9,
		QuoteType:           QuoteType_QUOTE_TYPE_BOLLYWOOD_DIALOGUE,
		Text:                "Success ke peeche mat bhago, excellence ka peecha karo... merge request khud aayega",
		EnglishTranslation:  "Don't run after success, chase excellence... merge request will come automatically",
		Source:              "3 Idiots",
		Character:           "Rancho Developer",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_FEATURES},
		IsFamilyFriendly:    true,
	},

	// Akshay Kumar Style
	{
		QuoteId:             10,
		QuoteType:           QuoteType_QUOTE_TYPE_BOLLYWOOD_DIALOGUE,
		Text:                "Khiladi commit karta hai, quit nahi",
		EnglishTranslation:  "A player commits, doesn't quit",
		Source:              "Khiladi 420 Commits",
		Character:           "Commit Kumar",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_STREAK},
		IsFamilyFriendly:    true,
	},

	// Rajinikanth Style
	{
		QuoteId:             11,
		QuoteType:           QuoteType_QUOTE_TYPE_BOLLYWOOD_DIALOGUE,
		Text:                "Documentation likhta hun, style mein!",
		EnglishTranslation:  "I write documentation, in style!",
		Source:              "Robot 2.0",
		Character:           "Documentation Rajni",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_DOCUMENTATION},
		IsFamilyFriendly:    true,
	},
	{
		QuoteId:             12,
		QuoteType:           QuoteType_QUOTE_TYPE_BOLLYWOOD_DIALOGUE,
		Text:                "Mind It! Code review pass ho gaya!",
		EnglishTranslation:  "Mind It! Code review passed!",
		Source:              "Sivaji The Coder",
		Character:           "Thalaiva",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_FEATURES},
		IsFamilyFriendly:    true,
	},

	// Hrithik Roshan Style
	{
		QuoteId:             13,
		QuoteType:           QuoteType_QUOTE_TYPE_BOLLYWOOD_DIALOGUE,
		Text:                "Main udta hun, crawl karta hun, code karta hun... main Krrish hun!",
		EnglishTranslation:  "I fly, I crawl, I code... I am Krrish!",
		Source:              "Krrish 3.0",
		Character:           "Super Coder",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_SPECIAL},
		IsFamilyFriendly:    true,
	},

	// Ranveer Singh Style
	{
		QuoteId:             14,
		QuoteType:           QuoteType_QUOTE_TYPE_BOLLYWOOD_DIALOGUE,
		Text:                "Apna time aayega... production deploy hoga!",
		EnglishTranslation:  "My time will come... production will be deployed!",
		Source:              "Gully Code",
		Character:           "MC Developer",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_FEATURES},
		IsFamilyFriendly:    true,
	},

	// Nawazuddin Style
	{
		QuoteId:             15,
		QuoteType:           QuoteType_QUOTE_TYPE_BOLLYWOOD_DIALOGUE,
		Text:                "Sacred bugs ko fix karna mushkil hi nahi, namumkin hai",
		EnglishTranslation:  "Fixing sacred bugs is not difficult, it's impossible",
		Source:              "Sacred Bugs",
		Character:           "Ganesh Gaitonde",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_BUG_FIXES},
		IsFamilyFriendly:    true,
	},

	// Sunny Deol Style
	{
		QuoteId:             16,
		QuoteType:           QuoteType_QUOTE_TYPE_BOLLYWOOD_DIALOGUE,
		Text:                "Yeh keyboard meri machine gun hai!",
		EnglishTranslation:  "This keyboard is my machine gun!",
		Source:              "Border: Line of Code",
		Character:           "Major Debug Singh",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_PERFORMANCE},
		IsFamilyFriendly:    true,
	},

	// Sanjay Dutt Style
	{
		QuoteId:             17,
		QuoteType:           QuoteType_QUOTE_TYPE_BOLLYWOOD_DIALOGUE,
		Text:                "Bhai, bug fix kar na... nahi toh jadoo ki jhappi de dunga!",
		EnglishTranslation:  "Brother, fix the bug... or I'll give you a magic hug!",
		Source:              "Munna Bhai MBBS (Master of Bug Busting Science)",
		Character:           "Munna Bhai",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_BUG_FIXES},
		IsFamilyFriendly:    true,
	},

	// Ajay Devgn Style
	{
		QuoteId:             18,
		QuoteType:           QuoteType_QUOTE_TYPE_BOLLYWOOD_DIALOGUE,
		Text:                "Aata majhi code-ka!",
		EnglishTranslation:  "Now comes my code!",
		Source:              "Singham Returns (with Features)",
		Character:           "Bajirao Singham",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_FEATURES},
		IsFamilyFriendly:    true,
	},

	// Paresh Rawal Style
	{
		QuoteId:             19,
		QuoteType:           QuoteType_QUOTE_TYPE_BOLLYWOOD_DIALOGUE,
		Text:                "Yeh bug kya hai? Yeh bug kya hai?!",
		EnglishTranslation:  "What is this bug? What is this bug?!",
		Source:              "Hera Pheri 2.0",
		Character:           "Baburao Ganpatrao Debugger",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_BUG_FIXES},
		IsFamilyFriendly:    true,
	},

	// Govinda Style
	{
		QuoteId:             20,
		QuoteType:           QuoteType_QUOTE_TYPE_BOLLYWOOD_DIALOGUE,
		Text:                "Code Number 1!",
		EnglishTranslation:  "Code Number 1!",
		Source:              "Coolie No. 1 Developer",
		Character:           "Raju Coder",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_COMMITS},
		IsFamilyFriendly:    true,
	},
}

// CricketCommentary - IPL style commentary for achievements
var CricketCommentary = []HumorQuote{
	{
		QuoteId:             100,
		QuoteType:           QuoteType_QUOTE_TYPE_CRICKET_COMMENTARY,
		Text:                "AND HE HAS HIT IT! Bug fix gone for a SIX! What a shot!",
		EnglishTranslation:  "",
		Source:              "IPL Commentary",
		Character:           "Harsha Bhogle",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_BUG_FIXES},
		IsFamilyFriendly:    true,
	},
	{
		QuoteId:             101,
		QuoteType:           QuoteType_QUOTE_TYPE_CRICKET_COMMENTARY,
		Text:                "Dhoni finishes off in style! Production deployed successfully!",
		EnglishTranslation:  "",
		Source:              "World Cup Final",
		Character:           "Ravi Shastri",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_FEATURES},
		IsFamilyFriendly:    true,
	},
	{
		QuoteId:             102,
		QuoteType:           QuoteType_QUOTE_TYPE_CRICKET_COMMENTARY,
		Text:                "In the air... and TAKEN! Critical bug caught by the developer!",
		EnglishTranslation:  "",
		Source:              "IPL",
		Character:           "Danny Morrison",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_BUG_FIXES},
		IsFamilyFriendly:    true,
	},
	{
		QuoteId:             103,
		QuoteType:           QuoteType_QUOTE_TYPE_CRICKET_COMMENTARY,
		Text:                "BUMRAH-STYLE DEBUGGING! Yorker to the bug! Clean bowled!",
		EnglishTranslation:  "",
		Source:              "IPL",
		Character:           "Akash Chopra",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_PERFORMANCE},
		IsFamilyFriendly:    true,
	},
	{
		QuoteId:             104,
		QuoteType:           QuoteType_QUOTE_TYPE_CRICKET_COMMENTARY,
		Text:                "Virat Kohli-esque consistency! 100 commits in a row!",
		EnglishTranslation:  "",
		Source:              "Test Match",
		Character:           "Sunil Gavaskar",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_STREAK},
		IsFamilyFriendly:    true,
	},
	{
		QuoteId:             105,
		QuoteType:           QuoteType_QUOTE_TYPE_CRICKET_COMMENTARY,
		Text:                "HELICOPTER SHOT! Feature deployed with MS Dhoni style!",
		EnglishTranslation:  "",
		Source:              "IPL Final",
		Character:           "Sanjay Manjrekar",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_FEATURES},
		IsFamilyFriendly:    true,
	},
	{
		QuoteId:             106,
		QuoteType:           QuoteType_QUOTE_TYPE_CRICKET_COMMENTARY,
		Text:                "Gayle storm in the repository! Commits raining like sixes!",
		EnglishTranslation:  "",
		Source:              "IPL",
		Character:           "Ian Bishop",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_COMMITS},
		IsFamilyFriendly:    true,
	},
	{
		QuoteId:             107,
		QuoteType:           QuoteType_QUOTE_TYPE_CRICKET_COMMENTARY,
		Text:                "Remember the name! This developer is on fire!",
		EnglishTranslation:  "",
		Source:              "IPL",
		Character:           "Ian Bishop",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_SPECIAL},
		IsFamilyFriendly:    true,
	},
	{
		QuoteId:             108,
		QuoteType:           QuoteType_QUOTE_TYPE_CRICKET_COMMENTARY,
		Text:                "ABD 360-degree coding! Features from every angle!",
		EnglishTranslation:  "",
		Source:              "IPL",
		Character:           "Kevin Pietersen",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_FEATURES},
		IsFamilyFriendly:    true,
	},
	{
		QuoteId:             109,
		QuoteType:           QuoteType_QUOTE_TYPE_CRICKET_COMMENTARY,
		Text:                "Sachin Tendulkar of coding! Century of pull requests!",
		EnglishTranslation:  "",
		Source:              "World Cup",
		Character:           "Harsha Bhogle",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_SPECIAL},
		IsFamilyFriendly:    true,
	},
}

// SouthIndianSuperstarQuotes - Rajini, Allu Arjun, etc.
var SouthIndianSuperstarQuotes = []HumorQuote{
	{
		QuoteId:             200,
		QuoteType:           QuoteType_QUOTE_TYPE_SOUTH_INDIAN_SUPERSTAR,
		Text:                "Enna Rascala! Bug-a fix panniten da!",
		EnglishTranslation:  "Hey Rascal! I fixed the bug!",
		Source:              "Robot",
		Character:           "Rajinikanth",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_BUG_FIXES},
		IsFamilyFriendly:    true,
	},
	{
		QuoteId:             201,
		QuoteType:           QuoteType_QUOTE_TYPE_SOUTH_INDIAN_SUPERSTAR,
		Text:                "Pushpa ka code... jhukega nahi!",
		EnglishTranslation:  "Pushpa's code... won't bow down!",
		Source:              "Pushpa",
		Character:           "Allu Arjun",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_FEATURES},
		IsFamilyFriendly:    true,
	},
	{
		QuoteId:             202,
		QuoteType:           QuoteType_QUOTE_TYPE_SOUTH_INDIAN_SUPERSTAR,
		Text:                "KGF Chapter Code - Mining bugs like gold!",
		EnglishTranslation:  "",
		Source:              "KGF",
		Character:           "Rocky Bhai",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_BUG_FIXES},
		IsFamilyFriendly:    true,
	},
	{
		QuoteId:             203,
		QuoteType:           QuoteType_QUOTE_TYPE_SOUTH_INDIAN_SUPERSTAR,
		Text:                "RRR - Rapid Revolutionary Refactoring!",
		EnglishTranslation:  "",
		Source:              "RRR",
		Character:           "Ram & Bheem",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_PERFORMANCE},
		IsFamilyFriendly:    true,
	},
	{
		QuoteId:             204,
		QuoteType:           QuoteType_QUOTE_TYPE_SOUTH_INDIAN_SUPERSTAR,
		Text:                "Vikram's commitment: Debug until dawn!",
		EnglishTranslation:  "",
		Source:              "Vikram",
		Character:           "Kamal Haasan",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_STREAK},
		IsFamilyFriendly:    true,
	},
	{
		QuoteId:             205,
		QuoteType:           QuoteType_QUOTE_TYPE_SOUTH_INDIAN_SUPERSTAR,
		Text:                "Master of Code! Thalapathy's teaching!",
		EnglishTranslation:  "",
		Source:              "Master",
		Character:           "Vijay",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_DOCUMENTATION},
		IsFamilyFriendly:    true,
	},
}

// ComedyPunchlines - Munna Bhai, Golmaal style
var ComedyPunchlines = []HumorQuote{
	{
		QuoteId:             300,
		QuoteType:           QuoteType_QUOTE_TYPE_COMEDY_PUNCHLINE,
		Text:                "Subah se deploy kar raha hun, PR merge hi nahi ho raha!",
		EnglishTranslation:  "Been deploying since morning, PR just won't merge!",
		Source:              "Hera Pheri",
		Character:           "Raju",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_FEATURES},
		IsFamilyFriendly:    true,
	},
	{
		QuoteId:             301,
		QuoteType:           QuoteType_QUOTE_TYPE_COMEDY_PUNCHLINE,
		Text:                "Mere ko bug dikhai deta hai, log nahi!",
		EnglishTranslation:  "I can see bugs, not people!",
		Source:              "Munna Bhai MBBS",
		Character:           "Circuit",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_BUG_FIXES},
		IsFamilyFriendly:    true,
	},
	{
		QuoteId:             302,
		QuoteType:           QuoteType_QUOTE_TYPE_COMEDY_PUNCHLINE,
		Text:                "Aye Circuit, documentation likh na!",
		EnglishTranslation:  "Hey Circuit, write the documentation!",
		Source:              "Munna Bhai MBBS",
		Character:           "Munna",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_DOCUMENTATION},
		IsFamilyFriendly:    true,
	},
	{
		QuoteId:             303,
		QuoteType:           QuoteType_QUOTE_TYPE_COMEDY_PUNCHLINE,
		Text:                "Code mein GOLMAAL hai bhai, sab GOLMAAL hai!",
		EnglishTranslation:  "There's chaos in the code, everything is chaos!",
		Source:              "Golmaal",
		Character:           "Gopal",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_BUG_FIXES},
		IsFamilyFriendly:    true,
	},
	{
		QuoteId:             304,
		QuoteType:           QuoteType_QUOTE_TYPE_COMEDY_PUNCHLINE,
		Text:                "Yeh Baburao ka style hai! Async/Await!",
		EnglishTranslation:  "This is Baburao's style! Async/Await!",
		Source:              "Hera Pheri",
		Character:           "Baburao",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_PERFORMANCE},
		IsFamilyFriendly:    true,
	},
}

// CodingParodies - Programming specific humor
var CodingParodies = []HumorQuote{
	{
		QuoteId:             400,
		QuoteType:           QuoteType_QUOTE_TYPE_CODING_PARODY,
		Text:                "Main apni favorite IDE hun!",
		EnglishTranslation:  "I am my favorite IDE!",
		Source:              "Original",
		Character:           "Poo (K3G)",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_COMMITS},
		IsFamilyFriendly:    true,
	},
	{
		QuoteId:             401,
		QuoteType:           QuoteType_QUOTE_TYPE_CODING_PARODY,
		Text:                "Code karta hu main, debug kya karun?",
		EnglishTranslation:  "I code, what's debugging to me?",
		Source:              "Original",
		Character:           "Developer",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_BUG_FIXES},
		IsFamilyFriendly:    true,
	},
	{
		QuoteId:             402,
		QuoteType:           QuoteType_QUOTE_TYPE_CODING_PARODY,
		Text:                "Mere pass CODE hai, LOGIC hai, GITHUB hai... tumhare pass kya hai?",
		EnglishTranslation:  "I have CODE, LOGIC, GITHUB... what do you have?",
		Source:              "Deewar",
		Character:           "Developer Bachchan",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_SPECIAL},
		IsFamilyFriendly:    true,
	},
	{
		QuoteId:             403,
		QuoteType:           QuoteType_QUOTE_TYPE_CODING_PARODY,
		Text:                "How's the code? FIRST CLASS!",
		EnglishTranslation:  "",
		Source:              "Original",
		Character:           "Rajni Style",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_FEATURES},
		IsFamilyFriendly:    true,
	},
	{
		QuoteId:             404,
		QuoteType:           QuoteType_QUOTE_TYPE_CODING_PARODY,
		Text:                "Error 404: Motivation not found... Just kidding! Keep coding!",
		EnglishTranslation:  "",
		Source:              "Developer Humor",
		Character:           "Console",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_STREAK},
		IsFamilyFriendly:    true,
	},
}

// MotivationalFilmy - Inspiring but funny
var MotivationalFilmy = []HumorQuote{
	{
		QuoteId:             500,
		QuoteType:           QuoteType_QUOTE_TYPE_MOTIVATIONAL_FILMY,
		Text:                "Picture abhi baaki hai mere dost! Keep pushing!",
		EnglishTranslation:  "The movie isn't over yet my friend! Keep pushing!",
		Source:              "Om Shanti Om",
		Character:           "Om Kapoor",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_STREAK},
		IsFamilyFriendly:    true,
	},
	{
		QuoteId:             501,
		QuoteType:           QuoteType_QUOTE_TYPE_MOTIVATIONAL_FILMY,
		Text:                "Kuch kuch commit hota hai... tum nahi samjhoge!",
		EnglishTranslation:  "Something something happens with commits... you won't understand!",
		Source:              "Kuch Kuch Hota Hai",
		Character:           "Rahul",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_COMMITS},
		IsFamilyFriendly:    true,
	},
	{
		QuoteId:             502,
		QuoteType:           QuoteType_QUOTE_TYPE_MOTIVATIONAL_FILMY,
		Text:                "Ja developer ja... jee le apni zindagi!",
		EnglishTranslation:  "Go developer go... live your life!",
		Source:              "DDLJ",
		Character:           "Simran's Dad",
		SuitableFor:         []AchievementCategory{AchievementCategory_ACHIEVEMENT_CATEGORY_SPECIAL},
		IsFamilyFriendly:    true,
	},
}

// HumorEngine manages quote selection and rotation
type HumorEngine struct {
	allQuotes       []HumorQuote
	usedQuotes      map[uint64]time.Time
	refreshInterval time.Duration
}

// NewHumorEngine creates a new humor engine
func NewHumorEngine() *HumorEngine {
	// Combine all quote collections
	allQuotes := make([]HumorQuote, 0)
	allQuotes = append(allQuotes, BollywoodDialogues...)
	allQuotes = append(allQuotes, CricketCommentary...)
	allQuotes = append(allQuotes, SouthIndianSuperstarQuotes...)
	allQuotes = append(allQuotes, ComedyPunchlines...)
	allQuotes = append(allQuotes, CodingParodies...)
	allQuotes = append(allQuotes, MotivationalFilmy...)

	return &HumorEngine{
		allQuotes:       allQuotes,
		usedQuotes:      make(map[uint64]time.Time),
		refreshInterval: 24 * time.Hour,
	}
}

// GetQuoteForAchievement returns a suitable quote for an achievement
func (h *HumorEngine) GetQuoteForAchievement(category AchievementCategory, preference HumorPreference) *HumorQuote {
	// Clean old used quotes
	h.cleanUsedQuotes()

	// Filter quotes by category and preference
	suitableQuotes := h.filterQuotes(category, preference)

	// Remove recently used quotes
	availableQuotes := make([]HumorQuote, 0)
	for _, quote := range suitableQuotes {
		if _, used := h.usedQuotes[quote.QuoteId]; !used {
			availableQuotes = append(availableQuotes, quote)
		}
	}

	// If all quotes used, reset and use any suitable quote
	if len(availableQuotes) == 0 {
		availableQuotes = suitableQuotes
	}

	// Select random quote
	if len(availableQuotes) > 0 {
		selected := availableQuotes[rand.Intn(len(availableQuotes))]
		h.usedQuotes[selected.QuoteId] = time.Now()
		return &selected
	}

	// Fallback quote
	return &HumorQuote{
		Text:               "Waah developer ji, waah!",
		EnglishTranslation: "Wow developer, wow!",
		Source:             "Original",
		Character:          "Everyone",
		IsFamilyFriendly:   true,
	}
}

// filterQuotes filters quotes based on category and preference
func (h *HumorEngine) filterQuotes(category AchievementCategory, preference HumorPreference) []HumorQuote {
	filtered := make([]HumorQuote, 0)

	for _, quote := range h.allQuotes {
		// Check category match
		categoryMatch := false
		for _, cat := range quote.SuitableFor {
			if cat == category || category == AchievementCategory_ACHIEVEMENT_CATEGORY_SPECIAL {
				categoryMatch = true
				break
			}
		}

		if !categoryMatch {
			continue
		}

		// Check preference match
		switch preference {
		case HumorPreference_HUMOR_PREFERENCE_BOLLYWOOD:
			if quote.QuoteType == QuoteType_QUOTE_TYPE_BOLLYWOOD_DIALOGUE {
				filtered = append(filtered, quote)
			}
		case HumorPreference_HUMOR_PREFERENCE_CRICKET:
			if quote.QuoteType == QuoteType_QUOTE_TYPE_CRICKET_COMMENTARY {
				filtered = append(filtered, quote)
			}
		case HumorPreference_HUMOR_PREFERENCE_SOUTH_INDIAN:
			if quote.QuoteType == QuoteType_QUOTE_TYPE_SOUTH_INDIAN_SUPERSTAR {
				filtered = append(filtered, quote)
			}
		case HumorPreference_HUMOR_PREFERENCE_COMEDY:
			if quote.QuoteType == QuoteType_QUOTE_TYPE_COMEDY_PUNCHLINE {
				filtered = append(filtered, quote)
			}
		case HumorPreference_HUMOR_PREFERENCE_MIXED:
			filtered = append(filtered, quote)
		default:
			filtered = append(filtered, quote)
		}
	}

	return filtered
}

// cleanUsedQuotes removes quotes older than refresh interval
func (h *HumorEngine) cleanUsedQuotes() {
	now := time.Now()
	for quoteId, usedTime := range h.usedQuotes {
		if now.Sub(usedTime) > h.refreshInterval {
			delete(h.usedQuotes, quoteId)
		}
	}
}

// GenerateStreakQuote creates a custom quote for streak achievements
func (h *HumorEngine) GenerateStreakQuote(days int) string {
	templates := []string{
		"🔥 %d din ki aag! Singham Returns daily!",
		"📿 %d days of devotion! Ram bhakt in GitHub!",
		"🏏 %d days batting! Kohli consistency!",
		"💪 %d din continuous! Dangal in repository!",
		"🚀 %d days non-stop! Pushpa never stops!",
	}

	return fmt.Sprintf(templates[rand.Intn(len(templates))], days)
}

// GenerateEarningsQuote creates earning announcement
func (h *HumorEngine) GenerateEarningsQuote(amount string, username string) string {
	templates := []string{
		"🎬 BREAKING: @%s ne kiya ₹%s ka BLOCKBUSTER commit! Housefull hai GitHub! 🍿",
		"🏏 SIXER! @%s ne bug fix karke kamaye ₹%s! Dhoni bhi proud! 🚁",
		"💃 @%s ka code nachaya blockchain ko! ₹%s ki barish! Nach Baliye! 🕺",
		"🦸 @%s: 'Main hu Hero!' - Feature banaya, ₹%s kamaya! 💪",
		"🎭 @%s ka ₹%s collection! Box office pe dhoom! 💥",
		"⚡ @%s ne maara ₹%s ka chauka! Crowd goes wild! 🏏",
		"🌟 @%s ki kamai ₹%s! Paisa hi paisa hoga! 💰",
	}

	return fmt.Sprintf(templates[rand.Intn(len(templates))], username, amount, username, amount)
}

// GetRandomMoviePosterTitle generates movie poster titles
func GetRandomMoviePosterTitle(achievementName string) string {
	templates := []string{
		"The %s Chronicles",
		"%s: The Developer Rises",
		"Mission %s: Code Protocol",
		"%s Returns",
		"The Legend of %s",
		"%s: Endgame",
		"The %s Strikes Back",
		"%s: No Way Home",
		"%s Forever",
		"The Amazing %s",
	}

	return fmt.Sprintf(templates[rand.Intn(len(templates))], achievementName)
}

// GetRandomTagline generates movie taglines
func GetRandomTagline() string {
	taglines := []string{
		"Ek commit se shuruat... production tak ka safar",
		"Where bugs fear to compile",
		"One developer. One mission. Zero bugs.",
		"Code mein hai dum!",
		"The repository will never be the same",
		"Based on true Git events",
		"From localhost to the cloud",
		"Every line tells a story",
		"Debugging ka badshah",
		"Feature banane ka formula",
	}

	return taglines[rand.Intn(len(taglines))]
}

// GetCriticReview generates fake movie reviews
func GetCriticReview() string {
	reviews := []string{
		"⭐⭐⭐⭐⭐ 'Brilliant performance! Oscar-worthy debugging!' - The GitHub Times",
		"⭐⭐⭐⭐⭐ 'Edge-of-your-seat deployment!' - Code Chronicle",
		"⭐⭐⭐⭐⭐ 'A masterpiece of software engineering!' - Dev Weekly",
		"⭐⭐⭐⭐⭐ 'Blockbuster hit! Must-watch commits!' - Repository Review",
		"⭐⭐⭐⭐⭐ 'Paisa vasool code!' - Bollywood Bugs",
		"⭐⭐⭐⭐⭐ 'Housefull in all terminals!' - Stack Overflow Daily",
		"⭐⭐⭐⭐⭐ 'National Award level programming!' - Indian Express IT",
	}

	return reviews[rand.Intn(len(reviews))]
}

// GenerateIPLCommentary creates live coding commentary
func GenerateIPLCommentary(action string, username string) string {
	switch action {
	case "commit":
		return fmt.Sprintf("AND @%s COMMITS! Straight to the repository! What timing!", username)
	case "bug_fix":
		return fmt.Sprintf("GONE! @%s sends the bug out of the park! Maximum debugging!", username)
	case "feature":
		return fmt.Sprintf("@%s at the keyboard, and... IT'S A FEATURE! The crowd goes wild!", username)
	case "performance":
		return fmt.Sprintf("LIGHTNING FAST! @%s optimizes like Bumrah's yorker! Unplayable!", username)
	default:
		return fmt.Sprintf("@%s is on fire! What a performance in the repository!", username)
	}
}

// GetHashtags generates trending hashtags
func GetHashtags(achievementType string) []string {
	base := []string{"#DeshChain", "#CodingComedy", "#GitHubGolmaal"}

	specific := map[string][]string{
		"bug_fix":     {"#BugBuster", "#DebugDangal", "#NoMoreBugs"},
		"feature":     {"#FeatureFilm", "#NewRelease", "#Blockbuster"},
		"streak":      {"#StreakSultan", "#CommitKing", "#DailyDose"},
		"performance": {"#SpeedSultan", "#OptimizationOscar", "#FastAndCurious"},
		"docs":        {"#DocumentaryDrama", "#ReadMeRockstar", "#DocsKaBaap"},
	}

	if specific[achievementType] != nil {
		return append(base, specific[achievementType]...)
	}

	return base
}

// IsViralWorthy checks if content has viral potential
func IsViralWorthy(content string) bool {
	viralKeywords := []string{
		"BREAKING", "BLOCKBUSTER", "SIXER", "₹", "100", "1000",
		"RECORD", "FIRST", "CENTURY", "THALAIVA", "SULTAN",
		"LEGENDARY", "EPIC", "MIND BLOWN", "🔥", "💯",
	}

	upperContent := strings.ToUpper(content)
	matchCount := 0

	for _, keyword := range viralKeywords {
		if strings.Contains(upperContent, keyword) {
			matchCount++
		}
	}

	return matchCount >= 3
}