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

// GetHashtags returns relevant hashtags for achievement categories
func GetHashtags(category string) []string {
	baseHashtags := []string{"#DeshChain", "#ProudIndianDeveloper", "#CodingWithCulture"}
	
	categoryHashtags := map[string][]string{
		"COMMITS":      {"#CommitKumar", "#DailyCommit", "#ConsistencyIsKey"},
		"BUG_FIXES":    {"#BugBusterBahubali", "#DebuggingNinja", "#BugFree"},
		"FEATURES":     {"#FeatureKhan", "#NewFeature", "#Innovation"},
		"DOCUMENTATION": {"#DocumentationRajni", "#DocsStyle", "#ClearDocs"},
		"PERFORMANCE":  {"#SpeedSultan", "#PerformanceMatters", "#FastCode"},
		"STREAK":       {"#StreakMaster", "#NeverMissADay", "#Dedication"},
		"SOCIAL":       {"#ViralCoder", "#CommunityFirst", "#ShareTheCode"},
		"SPECIAL":      {"#SpecialAchievement", "#EliteClub", "#RareFind"},
	}
	
	if specific, ok := categoryHashtags[category]; ok {
		return append(baseHashtags, specific...)
	}
	
	return baseHashtags
}

// IsViralWorthy checks if content has viral potential
func IsViralWorthy(content string) bool {
	viralKeywords := []string{
		"thalaiva", "bahubali", "sultan", "khan", "rajini",
		"bollywood", "cricket", "ipl", "dhoni", "kohli",
		"100", "420", "achievement", "legendary", "mythic",
		"first", "record", "breaking", "epic", "amazing",
	}
	
	contentLower := strings.ToLower(content)
	viralScore := 0
	
	for _, keyword := range viralKeywords {
		if strings.Contains(contentLower, keyword) {
			viralScore++
		}
	}
	
	// Has emojis
	if strings.ContainsAny(content, "🎊🏆💰🔥⭐👑🎯💯") {
		viralScore += 2
	}
	
	// Has exclamation marks
	if strings.Count(content, "!") >= 2 {
		viralScore++
	}
	
	return viralScore >= 3
}

// GetRandomMoviePosterTitle generates movie-style achievement titles
func GetRandomMoviePosterTitle(achievementName string) string {
	templates := []string{
		"The %s Returns",
		"%s: The Untold Story",
		"%s Ka Khiladi",
		"Bade %s",
		"%s No. 1",
		"Main Hoon %s",
		"%s: The Legend",
		"%s Zindabad",
		"Super %s",
		"%s: Breaking Records",
	}
	
	template := templates[rand.Intn(len(templates))]
	return fmt.Sprintf(template, achievementName)
}

// GetRandomTagline generates movie taglines
func GetRandomTagline() string {
	taglines := []string{
		"Code. Commit. Conquer.",
		"Where bugs fear to tread",
		"One developer. Infinite possibilities.",
		"Breaking bugs, not hearts",
		"The code must go on",
		"Debugging with style",
		"Features that touch hearts",
		"Speed that thrills",
		"Documentation that speaks",
		"Commitment personified",
		"Code mein hai dum",
		"Bug ka baap aaya",
		"Feature ka badshah",
		"Performance ka sultan",
		"Documentation ka don",
	}
	
	return taglines[rand.Intn(len(taglines))]
}

// GetCriticReview generates funny critic reviews
func GetCriticReview() string {
	reviews := []string{
		"⭐⭐⭐⭐⭐ 'Blockbuster code!' - Code Times",
		"⭐⭐⭐⭐⭐ 'Paisa vasool performance!' - Debug Daily",
		"⭐⭐⭐⭐⭐ 'Must-watch commits!' - Feature Films",
		"⭐⭐⭐⭐⭐ 'Housefull repository!' - Git Gazette",
		"⭐⭐⭐⭐⭐ 'Super hit debugging!' - Bug Bulletin",
		"⭐⭐⭐⭐⭐ '100 crore club code!' - Performance Post",
		"⭐⭐⭐⭐⭐ 'Family entertainer functions!' - Dev Digest",
		"⭐⭐⭐⭐⭐ 'Whistle-worthy features!' - Stack Stories",
		"⭐⭐⭐⭐⭐ 'Single screen sensation!' - Commit Chronicles",
		"⭐⭐⭐⭐⭐ 'Mass masala code!' - Terminal Tribune",
	}
	
	return reviews[rand.Intn(len(reviews))]
}

// GenerateTeamName creates IPL-style team names
func GenerateTeamName() string {
	cities := []string{
		"Mumbai", "Delhi", "Bangalore", "Chennai", "Kolkata",
		"Hyderabad", "Pune", "Jaipur", "Ahmedabad", "Lucknow",
	}
	
	suffixes := []string{
		"Coders", "Developers", "Debuggers", "Hackers", "Engineers",
		"Programmers", "Scripters", "Builders", "Architects", "Ninjas",
	}
	
	city := cities[rand.Intn(len(cities))]
	suffix := suffixes[rand.Intn(len(suffixes))]
	
	return fmt.Sprintf("%s %s", city, suffix)
}

// CalculateEngagementScore calculates overall engagement
func CalculateEngagementScore(metrics *EngagementMetrics) uint32 {
	if metrics == nil {
		return 0
	}
	
	// Weighted scoring
	score := metrics.Likes*1 + 
		metrics.Shares*3 + 
		metrics.Comments*2 + 
		metrics.Views/10 +
		metrics.Clicks/5
	
	if metrics.IsViral {
		score *= 2
	}
	
	return score
}

// GetTimeBasedGreeting returns greeting based on time
func GetTimeBasedGreeting() string {
	hour := time.Now().Hour()
	
	switch {
	case hour < 6:
		return "🌙 Late night coding session!"
	case hour < 12:
		return "🌅 Good morning, developer!"
	case hour < 16:
		return "☀️ Good afternoon, keep coding!"
	case hour < 20:
		return "🌆 Good evening, time to debug!"
	default:
		return "🌃 Night shift coding!"
	}
}

// GetMotivationalQuote returns coding motivation
func GetMotivationalQuote() string {
	quotes := []string{
		"Code karo, duniya badlo!",
		"Bug today, feature tomorrow!",
		"Commit karke dikha!",
		"Documentation is love!",
		"Performance matters, always!",
		"Keep calm and code on!",
		"Debugging is an art!",
		"Features make users happy!",
		"Clean code, happy life!",
		"Test first, code later!",
	}
	
	return quotes[rand.Intn(len(quotes))]
}

// GenerateWelcomeMessage creates welcome message for new developers
func GenerateWelcomeMessage(username string) string {
	return fmt.Sprintf(`
🎊 WELCOME TO DESHCHAIN FAMILY! 🎊
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Namaste @%s! 🙏

You've joined India's first cultural blockchain!
Where code meets culture, and developers become legends!

🎯 Your Journey Begins:
• Choose your avatar (5 Bollywood-style characters!)
• Start earning NAMO tokens
• Unlock achievements with desi tadka
• Compete in IPL-style coding battles
• Share your wins with viral-worthy posts

💡 Pro Tips:
• Commit daily for streak bonuses
• Fix bugs to become Bahubali
• Ship features like King Khan
• Document with Rajini style
• Optimize for Sultan speed

Remember: "Picture abhi shuru hui hai!"

Type '/help gamification' to start your journey!
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━`,
		username,
	)
}

// GenerateDailyMotivation creates daily motivation message
func GenerateDailyMotivation(profile *DeveloperProfile) string {
	greeting := GetTimeBasedGreeting()
	quote := GetMotivationalQuote()
	
	// Get progress info
	nextLevel := profile.Level + 1
	xpCurrent, xpRequired, percentage := GetXPProgressToNextLevel(profile.ExperiencePoints)
	
	// Get next achievement hint
	var nextAchievementHint string
	categories := []AchievementCategory{
		AchievementCategory_ACHIEVEMENT_CATEGORY_COMMITS,
		AchievementCategory_ACHIEVEMENT_CATEGORY_BUG_FIXES,
		AchievementCategory_ACHIEVEMENT_CATEGORY_FEATURES,
	}
	
	for _, category := range categories {
		if next := GetNextAchievementInCategory(category, 0); next != nil {
			nextAchievementHint = fmt.Sprintf("🎯 Next: %s (%s)", next.Name, next.Description)
			break
		}
	}
	
	return fmt.Sprintf(`
%s @%s!

📊 Daily Stats:
• Level %d → %d (%.1f%% complete)
• Current Streak: %d days
• Total Earnings: %s NAMO

%s

💭 "%s"

Keep pushing forward! Every commit counts! 🚀`,
		greeting, profile.GithubUsername,
		profile.Level, nextLevel, percentage,
		profile.CurrentStreakDays,
		profile.TotalEarnings.Amount,
		nextAchievementHint,
		quote,
	)
}

// ValidateAchievementUnlock validates if achievement can be unlocked
func ValidateAchievementUnlock(profile *DeveloperProfile, achievementId uint64) error {
	achievement := GetAchievementByID(achievementId)
	if achievement == nil {
		return fmt.Errorf("achievement not found")
	}
	
	// Check if already unlocked
	for _, unlockedId := range profile.AchievementsUnlocked {
		if unlockedId == achievementId {
			return fmt.Errorf("achievement already unlocked")
		}
	}
	
	// Check prerequisites
	for _, prereqId := range achievement.PrerequisiteAchievements {
		found := false
		for _, unlockedId := range profile.AchievementsUnlocked {
			if unlockedId == prereqId {
				found = true
				break
			}
		}
		if !found {
			prereq := GetAchievementByID(prereqId)
			if prereq != nil {
				return fmt.Errorf("prerequisite achievement '%s' not unlocked", prereq.Name)
			}
		}
	}
	
	return nil
}

// GenerateShareableContent creates content optimized for sharing
func GenerateShareableContent(achievement *Achievement, profile *DeveloperProfile) map[string]string {
	content := make(map[string]string)
	
	// Twitter (short and punchy)
	content["twitter"] = fmt.Sprintf(
		"🎊 Just unlocked '%s' on @DeshChain! %s #ProudIndianDeveloper",
		achievement.Name,
		achievement.UnlockQuote,
	)
	
	// LinkedIn (professional)
	content["linkedin"] = fmt.Sprintf(
		"Excited to share that I've achieved '%s' on DeshChain! %s This achievement represents my commitment to continuous learning and contribution to India's blockchain ecosystem. #DeshChain #BlockchainIndia",
		achievement.Name,
		achievement.Description,
	)
	
	// WhatsApp (uncle-friendly)
	content["whatsapp"] = fmt.Sprintf(
		"*Achievement Unlocked!* 🎊\n\n*%s*\n%s\n\n_%s_\n\nProud to be part of India's blockchain revolution! 🇮🇳",
		achievement.Name,
		achievement.Description,
		achievement.UnlockQuote,
	)
	
	return content
}