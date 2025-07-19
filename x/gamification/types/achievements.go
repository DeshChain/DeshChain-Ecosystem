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
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetDefaultAchievements returns all default achievements
func GetDefaultAchievements() []Achievement {
	return []Achievement{
		// Commit Achievements
		{
			AchievementId:    1,
			Name:             "Pehla Kadam",
			Description:      "Make your first commit to DeshChain",
			Category:         AchievementCategory_ACHIEVEMENT_CATEGORY_COMMITS,
			RequiredValue:    1,
			RewardAmount:     sdk.NewCoin("namo", sdk.NewInt(100)),
			XpReward:         50,
			BadgeImageUrl:    "/badges/pehla_kadam.png",
			UnlockQuote:      "Shubh aarambh! Picture abhi baaki hai mere dost!",
			RarityLevel:      RarityLevel_RARITY_LEVEL_COMMON,
			IsSecret:         false,
		},
		{
			AchievementId:    2,
			Name:             "Dilwale Developer Le Jayenge",
			Description:      "Reach 100 commits - True love for code!",
			Category:         AchievementCategory_ACHIEVEMENT_CATEGORY_COMMITS,
			RequiredValue:    100,
			RewardAmount:     sdk.NewCoin("namo", sdk.NewInt(5000)),
			XpReward:         500,
			BadgeImageUrl:    "/badges/ddlj_developer.png",
			UnlockQuote:      "Ja Simran ja, commit kar le!",
			RarityLevel:      RarityLevel_RARITY_LEVEL_EPIC,
			IsSecret:         false,
		},
		{
			AchievementId:    3,
			Name:             "Commit Kumar 420",
			Description:      "Achieve 420 commits - Khiladi level unlocked!",
			Category:         AchievementCategory_ACHIEVEMENT_CATEGORY_COMMITS,
			RequiredValue:    420,
			RewardAmount:     sdk.NewCoin("namo", sdk.NewInt(42000)),
			XpReward:         4200,
			BadgeImageUrl:    "/badges/commit_kumar_420.png",
			UnlockQuote:      "Khiladi 420 ban gaye! Sabka game bajayega!",
			RarityLevel:      RarityLevel_RARITY_LEVEL_LEGENDARY,
			IsSecret:         false,
		},
		{
			AchievementId:    4,
			Name:             "Code Ka Sachin",
			Description:      "Reach 1000 commits - Century of centuries!",
			Category:         AchievementCategory_ACHIEVEMENT_CATEGORY_COMMITS,
			RequiredValue:    1000,
			RewardAmount:     sdk.NewCoin("namo", sdk.NewInt(100000)),
			XpReward:         10000,
			BadgeImageUrl:    "/badges/code_sachin.png",
			UnlockQuote:      "Aila! Sachin... Sachin... SACHIN!",
			RarityLevel:      RarityLevel_RARITY_LEVEL_MYTHIC,
			IsSecret:         false,
		},

		// Bug Fix Achievements
		{
			AchievementId:    5,
			Name:             "Bug Ka Dushman",
			Description:      "Fix your first bug",
			Category:         AchievementCategory_ACHIEVEMENT_CATEGORY_BUG_FIXES,
			RequiredValue:    1,
			RewardAmount:     sdk.NewCoin("namo", sdk.NewInt(200)),
			XpReward:         100,
			BadgeImageUrl:    "/badges/bug_dushman.png",
			UnlockQuote:      "Bug ko maarne ke liye commitment chahiye!",
			RarityLevel:      RarityLevel_RARITY_LEVEL_COMMON,
			IsSecret:         false,
		},
		{
			AchievementId:    6,
			Name:             "Bahubali Bug Buster",
			Description:      "Fix 50 bugs - Become the Bahubali of debugging!",
			Category:         AchievementCategory_ACHIEVEMENT_CATEGORY_BUG_FIXES,
			RequiredValue:    50,
			RewardAmount:     sdk.NewCoin("namo", sdk.NewInt(10000)),
			XpReward:         1000,
			BadgeImageUrl:    "/badges/bahubali_bug_buster.png",
			UnlockQuote:      "Helicopter debugging activated! Katappa ne bugs ko kyun maara?",
			RarityLevel:      RarityLevel_RARITY_LEVEL_EPIC,
			IsSecret:         false,
		},
		{
			AchievementId:    7,
			Name:             "Thala Bug Terminator",
			Description:      "Fix 100 bugs - Captain Cool of debugging!",
			Category:         AchievementCategory_ACHIEVEMENT_CATEGORY_BUG_FIXES,
			RequiredValue:    100,
			RewardAmount:     sdk.NewCoin("namo", sdk.NewInt(25000)),
			XpReward:         2500,
			BadgeImageUrl:    "/badges/thala_terminator.png",
			UnlockQuote:      "Thala for a reason! Bugs ka baap!",
			RarityLevel:      RarityLevel_RARITY_LEVEL_LEGENDARY,
			IsSecret:         false,
		},

		// Feature Achievements
		{
			AchievementId:    8,
			Name:             "Feature Ka Raja",
			Description:      "Ship your first feature",
			Category:         AchievementCategory_ACHIEVEMENT_CATEGORY_FEATURES,
			RequiredValue:    1,
			RewardAmount:     sdk.NewCoin("namo", sdk.NewInt(500)),
			XpReward:         200,
			BadgeImageUrl:    "/badges/feature_raja.png",
			UnlockQuote:      "Naam hai feature, kaam hai user khush!",
			RarityLevel:      RarityLevel_RARITY_LEVEL_COMMON,
			IsSecret:         false,
		},
		{
			AchievementId:    9,
			Name:             "King Khan of Features",
			Description:      "Ship 25 features - Romance with code!",
			Category:         AchievementCategory_ACHIEVEMENT_CATEGORY_FEATURES,
			RequiredValue:    25,
			RewardAmount:     sdk.NewCoin("namo", sdk.NewInt(15000)),
			XpReward:         1500,
			BadgeImageUrl:    "/badges/king_khan_features.png",
			UnlockQuote:      "Don't underestimate the power of a common developer!",
			RarityLevel:      RarityLevel_RARITY_LEVEL_EPIC,
			IsSecret:         false,
		},
		{
			AchievementId:    10,
			Name:             "Badshah of Features",
			Description:      "Ship 50 features - Feature factory!",
			Category:         AchievementCategory_ACHIEVEMENT_CATEGORY_FEATURES,
			RequiredValue:    50,
			RewardAmount:     sdk.NewCoin("namo", sdk.NewInt(50000)),
			XpReward:         5000,
			BadgeImageUrl:    "/badges/badshah_features.png",
			UnlockQuote:      "Main hoon feature ka badshah! Users ki jaan!",
			RarityLevel:      RarityLevel_RARITY_LEVEL_LEGENDARY,
			IsSecret:         false,
		},

		// Documentation Achievements
		{
			AchievementId:    11,
			Name:             "Documentation Shuru Kiya",
			Description:      "Write your first documentation",
			Category:         AchievementCategory_ACHIEVEMENT_CATEGORY_DOCUMENTATION,
			RequiredValue:    1,
			RewardAmount:     sdk.NewCoin("namo", sdk.NewInt(300)),
			XpReward:         150,
			BadgeImageUrl:    "/badges/doc_shuru.png",
			UnlockQuote:      "Docs likhne ka style dekho!",
			RarityLevel:      RarityLevel_RARITY_LEVEL_COMMON,
			IsSecret:         false,
		},
		{
			AchievementId:    12,
			Name:             "Rajini Style Documentation",
			Description:      "Write 20 comprehensive docs - Mind it!",
			Category:         AchievementCategory_ACHIEVEMENT_CATEGORY_DOCUMENTATION,
			RequiredValue:    20,
			RewardAmount:     sdk.NewCoin("namo", sdk.NewInt(8000)),
			XpReward:         800,
			BadgeImageUrl:    "/badges/rajini_docs.png",
			UnlockQuote:      "En documentation-ah neenga edhirthu paakka mudiyaadhu!",
			RarityLevel:      RarityLevel_RARITY_LEVEL_RARE,
			IsSecret:         false,
		},
		{
			AchievementId:    13,
			Name:             "Thalaiva of Documentation",
			Description:      "Write 50 docs - Documentation superstar!",
			Category:         AchievementCategory_ACHIEVEMENT_CATEGORY_DOCUMENTATION,
			RequiredValue:    50,
			RewardAmount:     sdk.NewCoin("namo", sdk.NewInt(25000)),
			XpReward:         2500,
			BadgeImageUrl:    "/badges/thalaiva_docs.png",
			UnlockQuote:      "Style-u style-u documentation style-u!",
			RarityLevel:      RarityLevel_RARITY_LEVEL_EPIC,
			IsSecret:         false,
		},

		// Performance Achievements
		{
			AchievementId:    14,
			Name:             "Speed Ka Deewana",
			Description:      "Make your first performance improvement",
			Category:         AchievementCategory_ACHIEVEMENT_CATEGORY_PERFORMANCE,
			RequiredValue:    1,
			RewardAmount:     sdk.NewCoin("namo", sdk.NewInt(1000)),
			XpReward:         400,
			BadgeImageUrl:    "/badges/speed_deewana.png",
			UnlockQuote:      "Being fast is being human!",
			RarityLevel:      RarityLevel_RARITY_LEVEL_COMMON,
			IsSecret:         false,
		},
		{
			AchievementId:    15,
			Name:             "Dabangg Optimizer",
			Description:      "Make 10 major performance improvements",
			Category:         AchievementCategory_ACHIEVEMENT_CATEGORY_PERFORMANCE,
			RequiredValue:    10,
			RewardAmount:     sdk.NewCoin("namo", sdk.NewInt(20000)),
			XpReward:         2000,
			BadgeImageUrl:    "/badges/dabangg_optimizer.png",
			UnlockQuote:      "Swagat nahi karoge performance ka?",
			RarityLevel:      RarityLevel_RARITY_LEVEL_EPIC,
			IsSecret:         false,
		},
		{
			AchievementId:    16,
			Name:             "Sultan of Speed",
			Description:      "Make 25 performance improvements - Speed Sultan!",
			Category:         AchievementCategory_ACHIEVEMENT_CATEGORY_PERFORMANCE,
			RequiredValue:    25,
			RewardAmount:     sdk.NewCoin("namo", sdk.NewInt(75000)),
			XpReward:         7500,
			BadgeImageUrl:    "/badges/sultan_speed.png",
			UnlockQuote:      "Sultan ki performance dekhi? Ab code ki dekhlo!",
			RarityLevel:      RarityLevel_RARITY_LEVEL_LEGENDARY,
			IsSecret:         false,
		},

		// Streak Achievements
		{
			AchievementId:    17,
			Name:             "Saat Din Saat Commit",
			Description:      "Maintain a 7-day commit streak",
			Category:         AchievementCategory_ACHIEVEMENT_CATEGORY_STREAK,
			RequiredValue:    7,
			RewardAmount:     sdk.NewCoin("namo", sdk.NewInt(700)),
			XpReward:         350,
			BadgeImageUrl:    "/badges/saat_din.png",
			UnlockQuote:      "Saat din mein code double!",
			RarityLevel:      RarityLevel_RARITY_LEVEL_COMMON,
			IsSecret:         false,
		},
		{
			AchievementId:    18,
			Name:             "30 Din Ka Tapasya",
			Description:      "Maintain a 30-day commit streak - Coding Rishi!",
			Category:         AchievementCategory_ACHIEVEMENT_CATEGORY_STREAK,
			RequiredValue:    30,
			RewardAmount:     sdk.NewCoin("namo", sdk.NewInt(5000)),
			XpReward:         1500,
			BadgeImageUrl:    "/badges/30_din_tapasya.png",
			UnlockQuote:      "30 din mein paisa double? Nahi, code triple!",
			RarityLevel:      RarityLevel_RARITY_LEVEL_RARE,
			IsSecret:         false,
		},
		{
			AchievementId:    19,
			Name:             "Century Streak",
			Description:      "Maintain a 100-day commit streak - Kohli consistency!",
			Category:         AchievementCategory_ACHIEVEMENT_CATEGORY_STREAK,
			RequiredValue:    100,
			RewardAmount:     sdk.NewCoin("namo", sdk.NewInt(25000)),
			XpReward:         5000,
			BadgeImageUrl:    "/badges/century_streak.png",
			UnlockQuote:      "Chase master! 100 days of pure consistency!",
			RarityLevel:      RarityLevel_RARITY_LEVEL_EPIC,
			IsSecret:         false,
		},
		{
			AchievementId:    20,
			Name:             "365 Din Coding Tapasya",
			Description:      "Maintain a 365-day streak - Legendary dedication!",
			Category:         AchievementCategory_ACHIEVEMENT_CATEGORY_STREAK,
			RequiredValue:    365,
			RewardAmount:     sdk.NewCoin("namo", sdk.NewInt(365000)),
			XpReward:         36500,
			BadgeImageUrl:    "/badges/365_tapasya.png",
			UnlockQuote:      "Ek saal bina break! Thalaiva level dedication!",
			RarityLevel:      RarityLevel_RARITY_LEVEL_MYTHIC,
			IsSecret:         false,
		},

		// Social Achievements
		{
			AchievementId:    21,
			Name:             "Social Butterfly Developer",
			Description:      "Get 10 likes on your achievement post",
			Category:         AchievementCategory_ACHIEVEMENT_CATEGORY_SOCIAL,
			RequiredValue:    10,
			RewardAmount:     sdk.NewCoin("namo", sdk.NewInt(500)),
			XpReward:         250,
			BadgeImageUrl:    "/badges/social_butterfly.png",
			UnlockQuote:      "Log kya kahenge? Wah kya code hai!",
			RarityLevel:      RarityLevel_RARITY_LEVEL_COMMON,
			IsSecret:         false,
		},
		{
			AchievementId:    22,
			Name:             "Viral Coder",
			Description:      "Get 100 shares on your achievement post",
			Category:         AchievementCategory_ACHIEVEMENT_CATEGORY_SOCIAL,
			RequiredValue:    100,
			RewardAmount:     sdk.NewCoin("namo", sdk.NewInt(10000)),
			XpReward:         2000,
			BadgeImageUrl:    "/badges/viral_coder.png",
			UnlockQuote:      "Uncle ne WhatsApp pe forward kiya!",
			RarityLevel:      RarityLevel_RARITY_LEVEL_EPIC,
			IsSecret:         false,
		},
		{
			AchievementId:    23,
			Name:             "Influencer Developer",
			Description:      "Get 1000 total engagement on posts",
			Category:         AchievementCategory_ACHIEVEMENT_CATEGORY_SOCIAL,
			RequiredValue:    1000,
			RewardAmount:     sdk.NewCoin("namo", sdk.NewInt(50000)),
			XpReward:         10000,
			BadgeImageUrl:    "/badges/influencer_dev.png",
			UnlockQuote:      "Instagram reels bhi bana lo ab!",
			RarityLevel:      RarityLevel_RARITY_LEVEL_LEGENDARY,
			IsSecret:         false,
		},

		// Special/Secret Achievements
		{
			AchievementId:    24,
			Name:             "3 Idiots Special",
			Description:      "Fix a bug at 3 AM - All izz well!",
			Category:         AchievementCategory_ACHIEVEMENT_CATEGORY_SPECIAL,
			RequiredValue:    1,
			RewardAmount:     sdk.NewCoin("namo", sdk.NewInt(3000)),
			XpReward:         1500,
			BadgeImageUrl:    "/badges/3_idiots.png",
			UnlockQuote:      "Aal izz well! Success ke peeche mat bhago, excellence ka peecha karo!",
			RarityLevel:      RarityLevel_RARITY_LEVEL_RARE,
			IsSecret:         true,
		},
		{
			AchievementId:    25,
			Name:             "Munna Bhai Debugger",
			Description:      "Help 10 other developers fix their bugs",
			Category:         AchievementCategory_ACHIEVEMENT_CATEGORY_SPECIAL,
			RequiredValue:    10,
			RewardAmount:     sdk.NewCoin("namo", sdk.NewInt(15000)),
			XpReward:         3000,
			BadgeImageUrl:    "/badges/munna_bhai.png",
			UnlockQuote:      "Jadoo ki jhappi for debugging! Tension nahi leneka!",
			RarityLevel:      RarityLevel_RARITY_LEVEL_EPIC,
			IsSecret:         false,
		},
		{
			AchievementId:    26,
			Name:             "Kabhi Khushi Kabhi Bug",
			Description:      "Fix a bug and introduce a feature in same PR",
			Category:         AchievementCategory_ACHIEVEMENT_CATEGORY_SPECIAL,
			RequiredValue:    1,
			RewardAmount:     sdk.NewCoin("namo", sdk.NewInt(5000)),
			XpReward:         1000,
			BadgeImageUrl:    "/badges/k3g.png",
			UnlockQuote:      "It's all about loving your code!",
			RarityLevel:      RarityLevel_RARITY_LEVEL_RARE,
			IsSecret:         true,
		},
		{
			AchievementId:    27,
			Name:             "Code Mein Twist",
			Description:      "Refactor code that's older than 6 months",
			Category:         AchievementCategory_ACHIEVEMENT_CATEGORY_SPECIAL,
			RequiredValue:    1,
			RewardAmount:     sdk.NewCoin("namo", sdk.NewInt(8000)),
			XpReward:         1600,
			BadgeImageUrl:    "/badges/code_twist.png",
			UnlockQuote:      "Picture abhi baaki hai mere dost! Legacy code transformed!",
			RarityLevel:      RarityLevel_RARITY_LEVEL_RARE,
			IsSecret:         false,
		},
		{
			AchievementId:    28,
			Name:             "Desi Developer Dhamaka",
			Description:      "Complete 5 achievements in one day",
			Category:         AchievementCategory_ACHIEVEMENT_CATEGORY_SPECIAL,
			RequiredValue:    5,
			RewardAmount:     sdk.NewCoin("namo", sdk.NewInt(25000)),
			XpReward:         5000,
			BadgeImageUrl:    "/badges/desi_dhamaka.png",
			UnlockQuote:      "Ek din mein paanch! Ye hai asli entertainment!",
			RarityLevel:      RarityLevel_RARITY_LEVEL_LEGENDARY,
			IsSecret:         true,
		},
		{
			AchievementId:    29,
			Name:             "IPL Auction Star",
			Description:      "Get recruited by 3 teams for battles",
			Category:         AchievementCategory_ACHIEVEMENT_CATEGORY_SPECIAL,
			RequiredValue:    3,
			RewardAmount:     sdk.NewCoin("namo", sdk.NewInt(30000)),
			XpReward:         6000,
			BadgeImageUrl:    "/badges/ipl_star.png",
			UnlockQuote:      "Sold to highest bidder! Crore mein khelo!",
			RarityLevel:      RarityLevel_RARITY_LEVEL_EPIC,
			IsSecret:         false,
		},
		{
			AchievementId:    30,
			Name:             "Blockchain Ka Thalaiva",
			Description:      "Reach Level 100 - Ultimate achievement!",
			Category:         AchievementCategory_ACHIEVEMENT_CATEGORY_SPECIAL,
			RequiredValue:    100,
			RewardAmount:     sdk.NewCoin("namo", sdk.NewInt(1000000)),
			XpReward:         100000,
			BadgeImageUrl:    "/badges/blockchain_thalaiva.png",
			UnlockQuote:      "Kabaali da! Blockchain ka asli don!",
			RarityLevel:      RarityLevel_RARITY_LEVEL_MYTHIC,
			IsSecret:         false,
		},
	}
}

// GetAchievementByID returns achievement by ID
func GetAchievementByID(id uint64) *Achievement {
	achievements := GetDefaultAchievements()
	for _, achievement := range achievements {
		if achievement.AchievementId == id {
			return &achievement
		}
	}
	return nil
}

// GetAchievementsByCategory returns achievements for a specific category
func GetAchievementsByCategory(category AchievementCategory) []Achievement {
	var categoryAchievements []Achievement
	achievements := GetDefaultAchievements()
	
	for _, achievement := range achievements {
		if achievement.Category == category {
			categoryAchievements = append(categoryAchievements, achievement)
		}
	}
	
	return categoryAchievements
}

// GetAchievementsByRarity returns achievements for a specific rarity
func GetAchievementsByRarity(rarity RarityLevel) []Achievement {
	var rarityAchievements []Achievement
	achievements := GetDefaultAchievements()
	
	for _, achievement := range achievements {
		if achievement.RarityLevel == rarity {
			rarityAchievements = append(rarityAchievements, achievement)
		}
	}
	
	return rarityAchievements
}

// GetSecretAchievements returns all secret achievements
func GetSecretAchievements() []Achievement {
	var secretAchievements []Achievement
	achievements := GetDefaultAchievements()
	
	for _, achievement := range achievements {
		if achievement.IsSecret {
			secretAchievements = append(secretAchievements, achievement)
		}
	}
	
	return secretAchievements
}

// CalculateAchievementProgress calculates progress percentage
func CalculateAchievementProgress(currentValue, requiredValue uint64) float64 {
	if requiredValue == 0 {
		return 100.0
	}
	
	progress := float64(currentValue) / float64(requiredValue) * 100
	if progress > 100 {
		progress = 100
	}
	
	return progress
}

// GetNextAchievementInCategory returns the next achievement to unlock
func GetNextAchievementInCategory(category AchievementCategory, currentValue uint64) *Achievement {
	categoryAchievements := GetAchievementsByCategory(category)
	
	for _, achievement := range categoryAchievements {
		if currentValue < achievement.RequiredValue {
			return &achievement
		}
	}
	
	return nil
}

// GetAchievementCompletionBonus calculates bonus based on rarity
func GetAchievementCompletionBonus(rarity RarityLevel) sdk.Dec {
	switch rarity {
	case RarityLevel_RARITY_LEVEL_COMMON:
		return sdk.NewDec(1) // 1x multiplier
	case RarityLevel_RARITY_LEVEL_RARE:
		return sdk.NewDecWithPrec(15, 1) // 1.5x multiplier
	case RarityLevel_RARITY_LEVEL_EPIC:
		return sdk.NewDec(2) // 2x multiplier
	case RarityLevel_RARITY_LEVEL_LEGENDARY:
		return sdk.NewDec(3) // 3x multiplier
	case RarityLevel_RARITY_LEVEL_MYTHIC:
		return sdk.NewDec(5) // 5x multiplier
	default:
		return sdk.NewDec(1)
	}
}

// GetAchievementUnlockMessage creates unlock announcement
func GetAchievementUnlockMessage(achievement *Achievement, username string) string {
	return fmt.Sprintf(`
🎊 ACHIEVEMENT UNLOCKED! 🎊
━━━━━━━━━━━━━━━━━━━━━━━━
👤 Developer: %s
🏆 Achievement: %s
⭐ Rarity: %s
💰 Reward: %s NAMO
📈 XP Gained: %d

💬 %s

Keep coding, keep achieving! 🚀
━━━━━━━━━━━━━━━━━━━━━━━━`,
		username,
		achievement.Name,
		GetRarityDisplayName(achievement.RarityLevel),
		achievement.RewardAmount.Amount,
		achievement.XpReward,
		achievement.UnlockQuote,
	)
}

// GetRarityDisplayName returns display name for rarity
func GetRarityDisplayName(rarity RarityLevel) string {
	switch rarity {
	case RarityLevel_RARITY_LEVEL_COMMON:
		return "⚪ Common"
	case RarityLevel_RARITY_LEVEL_RARE:
		return "🔵 Rare"
	case RarityLevel_RARITY_LEVEL_EPIC:
		return "🟣 Epic"
	case RarityLevel_RARITY_LEVEL_LEGENDARY:
		return "🟡 Legendary"
	case RarityLevel_RARITY_LEVEL_MYTHIC:
		return "🔴 Mythic"
	default:
		return "Unknown"
	}
}

// CheckMultipleAchievements checks if user unlocked multiple achievements
func CheckMultipleAchievements(profile *DeveloperProfile, actionType string, value uint64) []uint64 {
	var unlockedAchievements []uint64
	
	// Map action types to achievement categories
	category := getAchievementCategory(actionType)
	if category == AchievementCategory_ACHIEVEMENT_CATEGORY_UNSPECIFIED {
		return unlockedAchievements
	}
	
	// Get current value based on action type
	currentValue := getCurrentValue(profile, actionType, value)
	
	// Check all achievements in this category
	achievements := GetAchievementsByCategory(category)
	for _, achievement := range achievements {
		// Skip if already unlocked
		if isAchievementUnlocked(profile.AchievementsUnlocked, achievement.AchievementId) {
			continue
		}
		
		// Check if requirement is met
		if currentValue >= achievement.RequiredValue {
			unlockedAchievements = append(unlockedAchievements, achievement.AchievementId)
		}
	}
	
	// Check special achievements
	specialAchievements := checkSpecialAchievements(profile, actionType, value)
	unlockedAchievements = append(unlockedAchievements, specialAchievements...)
	
	return unlockedAchievements
}

// Helper functions
func getAchievementCategory(actionType string) AchievementCategory {
	switch actionType {
	case "commit":
		return AchievementCategory_ACHIEVEMENT_CATEGORY_COMMITS
	case "bug_fix":
		return AchievementCategory_ACHIEVEMENT_CATEGORY_BUG_FIXES
	case "feature":
		return AchievementCategory_ACHIEVEMENT_CATEGORY_FEATURES
	case "documentation":
		return AchievementCategory_ACHIEVEMENT_CATEGORY_DOCUMENTATION
	case "performance":
		return AchievementCategory_ACHIEVEMENT_CATEGORY_PERFORMANCE
	case "streak":
		return AchievementCategory_ACHIEVEMENT_CATEGORY_STREAK
	case "social":
		return AchievementCategory_ACHIEVEMENT_CATEGORY_SOCIAL
	default:
		return AchievementCategory_ACHIEVEMENT_CATEGORY_UNSPECIFIED
	}
}

func getCurrentValue(profile *DeveloperProfile, actionType string, additionalValue uint64) uint64 {
	switch actionType {
	case "commit":
		return profile.TotalCommits + additionalValue
	case "bug_fix":
		return profile.TotalBugsFi xed + additionalValue
	case "feature":
		return profile.TotalFeaturesShipped + additionalValue
	case "documentation":
		return profile.TotalDocsWritten + additionalValue
	case "performance":
		return profile.PerformanceImprovements + additionalValue
	case "streak":
		return uint64(profile.CurrentStreakDays)
	default:
		return 0
	}
}

func isAchievementUnlocked(unlockedList []uint64, achievementId uint64) bool {
	for _, id := range unlockedList {
		if id == achievementId {
			return true
		}
	}
	return false
}

func checkSpecialAchievements(profile *DeveloperProfile, actionType string, value uint64) []uint64 {
	var specialUnlocks []uint64
	
	// Check for "Desi Developer Dhamaka" - 5 achievements in one day
	if len(profile.AchievementsUnlocked) > 0 && len(profile.AchievementsUnlocked)%5 == 0 {
		if !isAchievementUnlocked(profile.AchievementsUnlocked, 28) {
			specialUnlocks = append(specialUnlocks, 28)
		}
	}
	
	// Check for level 100 achievement
	if profile.Level >= 100 && !isAchievementUnlocked(profile.AchievementsUnlocked, 30) {
		specialUnlocks = append(specialUnlocks, 30)
	}
	
	return specialUnlocks
}