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

// GetLevelConfigs returns all level configurations
func GetLevelConfigs() []LevelConfig {
	configs := []LevelConfig{}
	
	// Generate level configs from 1 to 100
	for level := uint32(1); level <= 100; level++ {
		config := generateLevelConfig(level)
		configs = append(configs, config)
	}
	
	return configs
}

// generateLevelConfig creates configuration for a specific level
func generateLevelConfig(level uint32) LevelConfig {
	// Base XP requirement with exponential growth
	baseXP := uint64(100)
	xpRequired := baseXP * uint64(level) * uint64(level)
	
	// Calculate bonus reward
	bonusAmount := int64(level * 100)
	if level%10 == 0 {
		bonusAmount *= 5 // 5x bonus for milestone levels
	}
	
	// Get level title and quote
	title, quote := getLevelTitleAndQuote(level)
	
	return LevelConfig{
		Level:       level,
		XpRequired:  xpRequired,
		Title:       title,
		UnlockQuote: quote,
		BonusReward: sdk.NewCoin("namo", sdk.NewInt(bonusAmount)),
	}
}

// getLevelTitleAndQuote returns title and quote for each level
func getLevelTitleAndQuote(level uint32) (string, string) {
	switch {
	case level == 1:
		return "Coding Ka Naya Sitara", "Shubharambh! Welcome to DeshChain family!"
	case level < 5:
		return "Junior Developer", "Seekhte raho, badhte raho!"
	case level < 10:
		return "Code Warrior", "Warrior ban gaye ho, ab hero banna hai!"
	case level == 10:
		return "Rising Star Developer", "10 ka dum! Rising star of DeshChain!"
	case level < 20:
		return "Experienced Coder", "Experience ke saath confidence aata hai!"
	case level == 20:
		return "Senior Developer", "Senior ban gaye! Junior ki help karo!"
	case level == 25:
		return "Code Guru", "Quarter century! Guru level unlocked!"
	case level < 30:
		return "Tech Lead Material", "Leadership qualities dikh rahi hai!"
	case level == 30:
		return "Tech Lead", "30 ka power! Team ko lead karo!"
	case level < 40:
		return "Principal Developer", "Principal position pe pahunch gaye!"
	case level == 40:
		return "Code Architect", "40 fantastic! Architecture master!"
	case level == 50:
		return "Half Century Hero", "50 ka half century! Standing ovation!"
	case level < 60:
		return "Distinguished Engineer", "Distinguished performance!"
	case level == 60:
		return "Engineering Maestro", "60 ka experience! Maestro level!"
	case level < 70:
		return "Tech Visionary", "Vision clear hai, path clear hai!"
	case level == 70:
		return "Code Shahenshah", "70 ka Shahenshah! Korona tumhara!"
	case level == 75:
		return "Diamond Developer", "75 ka diamond jubilee! Chamak raha hai!"
	case level < 80:
		return "Elite Engineer", "Elite club mein entry!"
	case level == 80:
		return "Tech Tycoon", "80 ka tycoon! Business of code!"
	case level < 90:
		return "Code Legend", "Legend in the making!"
	case level == 90:
		return "Tech Titan", "90 ka Titan! Almost at the top!"
	case level < 100:
		return "Supreme Developer", "Supreme court of coding!"
	case level == 100:
		return "Blockchain Ka Thalaiva", "Century complete! Thalaiva of DeshChain!"
	default:
		return "Master Developer", "Master of the code universe!"
	}
}

// CalculateXPRequired calculates total XP needed for a level
func CalculateXPRequired(level uint32) uint64 {
	if level <= 0 || level > 100 {
		return 0
	}
	
	configs := GetLevelConfigs()
	return configs[level-1].XpRequired
}

// GetCurrentLevel calculates current level from XP
func GetCurrentLevel(totalXP uint64) uint32 {
	configs := GetLevelConfigs()
	
	for i := len(configs) - 1; i >= 0; i-- {
		if totalXP >= configs[i].XpRequired {
			return configs[i].Level
		}
	}
	
	return 1 // Default to level 1
}

// GetXPProgressToNextLevel calculates progress to next level
func GetXPProgressToNextLevel(totalXP uint64) (currentXP, requiredXP uint64, percentage float64) {
	currentLevel := GetCurrentLevel(totalXP)
	
	if currentLevel >= 100 {
		return totalXP, totalXP, 100.0
	}
	
	currentLevelXP := CalculateXPRequired(currentLevel)
	nextLevelXP := CalculateXPRequired(currentLevel + 1)
	
	currentXP = totalXP - currentLevelXP
	requiredXP = nextLevelXP - currentLevelXP
	
	if requiredXP > 0 {
		percentage = float64(currentXP) / float64(requiredXP) * 100
	}
	
	return
}

// GetLevelUpMessage creates level up announcement
func GetLevelUpMessage(profile *DeveloperProfile, newLevel uint32) string {
	config := GetLevelConfigByLevel(newLevel)
	if config == nil {
		return ""
	}
	
	avatar := GetAvatarByType(profile.ActiveAvatar)
	avatarName := "Developer"
	if avatar != nil {
		avatarName = avatar.Name
	}
	
	// Special messages for milestone levels
	specialMessage := ""
	switch newLevel {
	case 10:
		specialMessage = "\n\n🎬 SPECIAL SCENE: Entry music plays! Slow motion walk!"
	case 25:
		specialMessage = "\n\n🎪 INTERVAL: Samosa break! You've earned it!"
	case 50:
		specialMessage = "\n\n🏏 HALF CENTURY: Stadium mein standing ovation!"
	case 75:
		specialMessage = "\n\n💎 DIAMOND JUBILEE: Heera hai tu heera!"
	case 100:
		specialMessage = "\n\n👑 THALAIVA MOMENT: Fireworks! Dhol! Party!"
	}
	
	return fmt.Sprintf(`
🎊 LEVEL UP! 🎊
━━━━━━━━━━━━━━━━━━━━━━━━
👤 Developer: %s
🎭 Avatar: %s
⬆️ New Level: %d
🏷️ New Title: %s
💰 Bonus Reward: %s NAMO

💬 %s%s

Aage badhte raho! Next level awaits! 🚀
━━━━━━━━━━━━━━━━━━━━━━━━`,
		profile.GithubUsername,
		avatarName,
		newLevel,
		config.Title,
		config.BonusReward.Amount,
		config.UnlockQuote,
		specialMessage,
	)
}

// GetLevelConfigByLevel returns config for specific level
func GetLevelConfigByLevel(level uint32) *LevelConfig {
	configs := GetLevelConfigs()
	if level > 0 && level <= uint32(len(configs)) {
		return &configs[level-1]
	}
	return nil
}

// CalculateXPFromAction calculates XP based on action and params
func CalculateXPFromAction(params *GamificationParams, action string, avatar AvatarType) uint64 {
	var baseXP uint64
	
	switch action {
	case "commit":
		baseXP = params.BaseXpPerCommit
	case "bug_fix":
		baseXP = params.BaseXpPerBugFix
	case "feature":
		baseXP = params.BaseXpPerFeature
	case "documentation":
		baseXP = params.BaseXpPerDoc
	default:
		baseXP = 50 // Default XP
	}
	
	// Apply avatar bonus
	avatarData := GetAvatarByType(avatar)
	if avatarData != nil && isAvatarSpecialty(avatar, action) {
		multiplier, _ := sdk.NewDecFromStr(avatarData.RewardMultiplier)
		baseXP = uint64(multiplier.MulInt64(int64(baseXP)).TruncateInt64())
	}
	
	return baseXP
}

// isAvatarSpecialty checks if action matches avatar specialty
func isAvatarSpecialty(avatar AvatarType, action string) bool {
	switch avatar {
	case AvatarType_AVATAR_TYPE_BUG_BUSTER_BAHUBALI:
		return action == "bug_fix"
	case AvatarType_AVATAR_TYPE_FEATURE_KHAN:
		return action == "feature"
	case AvatarType_AVATAR_TYPE_DOCUMENTATION_RAJNI:
		return action == "documentation"
	case AvatarType_AVATAR_TYPE_SPEED_SULTAN:
		return action == "performance"
	case AvatarType_AVATAR_TYPE_COMMIT_KUMAR:
		return action == "commit"
	default:
		return false
	}
}

// GetLeaderboardTitle returns special title based on rank
func GetLeaderboardTitle(rank uint64) string {
	switch rank {
	case 1:
		return "👑 Coding Ka Badshah"
	case 2:
		return "🥈 Code Ka Nawab"
	case 3:
		return "🥉 Programming Ka Sultan"
	case 4, 5:
		return "⭐ Top 5 Elite"
	case 6, 7, 8, 9, 10:
		return "🌟 Top 10 Master"
	default:
		if rank <= 25 {
			return "💫 Top 25 Expert"
		} else if rank <= 50 {
			return "✨ Top 50 Pro"
		} else if rank <= 100 {
			return "🔥 Top 100 Warrior"
		}
		return "💪 Active Developer"
	}
}

// CalculateStreakBonus calculates bonus XP for streaks
func CalculateStreakBonus(streakDays uint32, baseXP uint64, multiplierStr string) uint64 {
	if streakDays <= 0 {
		return baseXP
	}
	
	// Parse multiplier
	multiplier, err := sdk.NewDecFromStr(multiplierStr)
	if err != nil {
		multiplier = sdk.NewDecWithPrec(11, 1) // Default 1.1x
	}
	
	// Calculate streak bonus
	var bonus sdk.Dec
	switch {
	case streakDays >= 365:
		bonus = sdk.NewDec(5) // 5x for year-long streak
	case streakDays >= 100:
		bonus = sdk.NewDec(3) // 3x for 100+ days
	case streakDays >= 30:
		bonus = sdk.NewDec(2) // 2x for 30+ days
	case streakDays >= 7:
		bonus = multiplier // Use param multiplier for 7+ days
	default:
		bonus = sdk.NewDec(1) // No bonus
	}
	
	return uint64(bonus.MulInt64(int64(baseXP)).TruncateInt64())
}

// GetRankChangeMessage creates rank change announcement
func GetRankChangeMessage(username string, oldRank, newRank uint64, direction string) string {
	emoji := "📈"
	action := "climbed"
	if direction == "down" {
		emoji = "📉"
		action = "dropped"
	}
	
	oldTitle := GetLeaderboardTitle(oldRank)
	newTitle := GetLeaderboardTitle(newRank)
	
	// Special messages for top ranks
	specialMsg := ""
	if newRank == 1 && direction == "up" {
		specialMsg = "\n\n👑 NEW KING OF THE LEADERBOARD! 👑\nSabka game bajega!"
	} else if newRank <= 3 && direction == "up" {
		specialMsg = fmt.Sprintf("\n\n🏆 PODIUM FINISH! Rank #%d achieved!", newRank)
	} else if newRank <= 10 && oldRank > 10 && direction == "up" {
		specialMsg = "\n\n⭐ TOP 10 ENTRY! Elite club mein aapka swagat hai!"
	}
	
	return fmt.Sprintf(`
%s RANK CHANGE ALERT %s
━━━━━━━━━━━━━━━━━━━━━━━━
👤 Developer: %s
%s Previous Rank: #%d (%s)
%s New Rank: #%d (%s)

You %s %d positions!%s

Keep pushing for the top! 🚀
━━━━━━━━━━━━━━━━━━━━━━━━`,
		emoji, emoji,
		username,
		"📊", oldRank, oldTitle,
		"🎯", newRank, newTitle,
		action, abs(int64(newRank)-int64(oldRank)),
		specialMsg,
	)
}

// Helper function for absolute value
func abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}

// GetMilestoneRewards returns special rewards for milestones
func GetMilestoneRewards(level uint32) sdk.Coins {
	rewards := sdk.NewCoins()
	
	// Special milestone rewards
	switch level {
	case 10:
		rewards = rewards.Add(sdk.NewCoin("namo", sdk.NewInt(5000)))
	case 25:
		rewards = rewards.Add(sdk.NewCoin("namo", sdk.NewInt(15000)))
	case 50:
		rewards = rewards.Add(sdk.NewCoin("namo", sdk.NewInt(50000)))
	case 75:
		rewards = rewards.Add(sdk.NewCoin("namo", sdk.NewInt(100000)))
	case 100:
		rewards = rewards.Add(sdk.NewCoin("namo", sdk.NewInt(500000)))
	}
	
	return rewards
}