package types

import (
	"fmt"
	"sort"
	"time"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// LeaderboardManager manages different leaderboards
type LeaderboardManager struct {
	GlobalLeaderboard    []LeaderboardEntry
	WeeklyLeaderboard    []LeaderboardEntry
	MonthlyLeaderboard   []LeaderboardEntry
	CategoryLeaderboards map[string][]LeaderboardEntry // By achievement category
}

// NewLeaderboardManager creates a new leaderboard manager
func NewLeaderboardManager() *LeaderboardManager {
	return &LeaderboardManager{
		GlobalLeaderboard:    []LeaderboardEntry{},
		WeeklyLeaderboard:    []LeaderboardEntry{},
		MonthlyLeaderboard:   []LeaderboardEntry{},
		CategoryLeaderboards: make(map[string][]LeaderboardEntry),
	}
}

// UpdateLeaderboards updates all leaderboards with new data
func (lm *LeaderboardManager) UpdateLeaderboards(profiles []DeveloperProfile) {
	// Update global leaderboard
	lm.GlobalLeaderboard = lm.createLeaderboard(profiles, "global")
	
	// Update weekly leaderboard (based on recent activity)
	lm.WeeklyLeaderboard = lm.createLeaderboard(profiles, "weekly")
	
	// Update monthly leaderboard
	lm.MonthlyLeaderboard = lm.createLeaderboard(profiles, "monthly")
	
	// Update category leaderboards
	categories := []string{"commits", "bugs", "features", "docs", "performance", "streak"}
	for _, category := range categories {
		lm.CategoryLeaderboards[category] = lm.createLeaderboard(profiles, category)
	}
}

// createLeaderboard creates a leaderboard based on type
func (lm *LeaderboardManager) createLeaderboard(profiles []DeveloperProfile, leaderboardType string) []LeaderboardEntry {
	entries := []LeaderboardEntry{}
	
	for _, profile := range profiles {
		entry := LeaderboardEntry{
			DeveloperAddress: profile.Address,
			GithubUsername:   profile.GithubUsername,
			Level:            profile.Level,
			TotalXp:          profile.ExperiencePoints,
			TotalEarnings:    profile.TotalEarnings,
			ActiveAvatar:     profile.ActiveAvatar,
			AchievementsCount: uint32(len(profile.AchievementsUnlocked)),
			CurrentStreak:    profile.CurrentStreakDays,
		}
		
		// Set special title based on performance
		entry.SpecialTitle = lm.getSpecialTitle(&profile, leaderboardType)
		
		entries = append(entries, entry)
	}
	
	// Sort based on leaderboard type
	lm.sortLeaderboard(entries, leaderboardType)
	
	// Assign ranks
	for i := range entries {
		entries[i].Rank = uint64(i + 1)
	}
	
	return entries
}

// sortLeaderboard sorts entries based on type
func (lm *LeaderboardManager) sortLeaderboard(entries []LeaderboardEntry, leaderboardType string) {
	switch leaderboardType {
	case "global":
		// Sort by total XP
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].TotalXp > entries[j].TotalXp
		})
	case "weekly", "monthly":
		// Sort by recent earnings (would need timestamp data in real implementation)
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].TotalEarnings.Amount.GT(entries[j].TotalEarnings.Amount)
		})
	case "commits":
		// Sort by total commits (would need this field)
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].TotalXp > entries[j].TotalXp // Placeholder
		})
	case "streak":
		// Sort by current streak
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].CurrentStreak > entries[j].CurrentStreak
		})
	default:
		// Default to XP
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].TotalXp > entries[j].TotalXp
		})
	}
}

// getSpecialTitle assigns special titles based on performance
func (lm *LeaderboardManager) getSpecialTitle(profile *DeveloperProfile, leaderboardType string) string {
	switch leaderboardType {
	case "streak":
		if profile.CurrentStreakDays >= 365 {
			return "🔥 Streak Samrat"
		} else if profile.CurrentStreakDays >= 100 {
			return "🎯 Century Striker"
		} else if profile.CurrentStreakDays >= 30 {
			return "⚡ Streak Sultan"
		}
	case "commits":
		if profile.TotalCommits >= 1000 {
			return "💫 Commit Maharaja"
		} else if profile.TotalCommits >= 500 {
			return "🌟 Commit Commander"
		}
	}
	
	// Avatar-based titles
	avatarTitle := GetAvatarRankTitle(profile.ActiveAvatar, profile.Level)
	return avatarTitle
}

// GetTopPerformers returns top N performers from a leaderboard
func GetTopPerformers(leaderboard []LeaderboardEntry, limit int) []LeaderboardEntry {
	if len(leaderboard) < limit {
		return leaderboard
	}
	return leaderboard[:limit]
}

// GetLeaderboardAroundUser returns entries around a specific user
func GetLeaderboardAroundUser(leaderboard []LeaderboardEntry, userAddress string, range_ int) []LeaderboardEntry {
	userIndex := -1
	for i, entry := range leaderboard {
		if entry.DeveloperAddress == userAddress {
			userIndex = i
			break
		}
	}
	
	if userIndex == -1 {
		return []LeaderboardEntry{}
	}
	
	start := userIndex - range_
	if start < 0 {
		start = 0
	}
	
	end := userIndex + range_ + 1
	if end > len(leaderboard) {
		end = len(leaderboard)
	}
	
	return leaderboard[start:end]
}

// CreateLeaderboardDisplay creates formatted leaderboard display
func CreateLeaderboardDisplay(entries []LeaderboardEntry, displayType string) string {
	if len(entries) == 0 {
		return "No entries in leaderboard yet!"
	}
	
	var title string
	switch displayType {
	case "global":
		title = "🏆 GLOBAL LEADERBOARD - DeshChain Champions 🏆"
	case "weekly":
		title = "📅 WEEKLY WARRIORS - This Week's Heroes 📅"
	case "monthly":
		title = "📆 MONTHLY MASTERS - Month's Top Coders 📆"
	case "streak":
		title = "🔥 STREAK SULTANS - Consistency Kings 🔥"
	default:
		title = "🎯 LEADERBOARD 🎯"
	}
	
	display := fmt.Sprintf(`
%s
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

`, title)
	
	for _, entry := range entries {
		rankEmoji := getRankEmoji(entry.Rank)
		avatar := GetAvatarByType(entry.ActiveAvatar)
		avatarName := "Developer"
		if avatar != nil {
			avatarName = avatar.Name
		}
		
		display += fmt.Sprintf(`%s #%d | @%s
   🎭 %s | Level %d | %s XP
   💰 %s NAMO | 🏆 %d Achievements
   🔥 %d Day Streak | %s
   ─────────────────────────────────────
`,
			rankEmoji, entry.Rank, entry.GithubUsername,
			avatarName, entry.Level, formatNumber(entry.TotalXp),
			entry.TotalEarnings.Amount, entry.AchievementsCount,
			entry.CurrentStreak, entry.SpecialTitle,
		)
	}
	
	return display
}

// getRankEmoji returns emoji based on rank
func getRankEmoji(rank uint64) string {
	switch rank {
	case 1:
		return "👑"
	case 2:
		return "🥈"
	case 3:
		return "🥉"
	default:
		if rank <= 10 {
			return "⭐"
		} else if rank <= 25 {
			return "✨"
		} else if rank <= 50 {
			return "🌟"
		} else if rank <= 100 {
			return "💫"
		}
		return "🎯"
	}
}

// formatNumber formats large numbers with commas
func formatNumber(num uint64) string {
	str := fmt.Sprintf("%d", num)
	n := len(str)
	if n <= 3 {
		return str
	}
	
	// Add commas for readability
	result := ""
	for i, ch := range str {
		if i > 0 && (n-i)%3 == 0 {
			result += ","
		}
		result += string(ch)
	}
	
	return result
}

// CreateBattleLeaderboard creates IPL-style team battle leaderboard
func CreateBattleLeaderboard(battles []TeamBattle) string {
	// Sort battles by total score
	sort.Slice(battles, func(i, j int) bool {
		totalI := battles[i].Team1Score + battles[i].Team2Score
		totalJ := battles[j].Team1Score + battles[j].Team2Score
		return totalI > totalJ
	})
	
	display := `
🏏 IPL-STYLE TEAM BATTLE STANDINGS 🏏
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

`
	
	for i, battle := range battles {
		winner := battle.Team1Name
		winnerScore := battle.Team1Score
		loser := battle.Team2Name
		loserScore := battle.Team2Score
		
		if battle.Team2Score > battle.Team1Score {
			winner = battle.Team2Name
			winnerScore = battle.Team2Score
			loser = battle.Team1Name
			loserScore = battle.Team1Score
		}
		
		display += fmt.Sprintf(`Match %d: %s
🏆 Winner: %s (%d points)
🥈 Runner-up: %s (%d points)
💰 Prize Pool: %s NAMO
⏱️ Duration: %s
─────────────────────────────────────

`,
			i+1, battle.BattleType,
			winner, winnerScore,
			loser, loserScore,
			battle.PrizePool.Amount,
			formatDuration(battle.EndTime.Sub(battle.StartTime)),
		)
	}
	
	return display
}

// formatDuration formats duration in readable format
func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

// GetLeaderboardStats returns statistics about the leaderboard
func GetLeaderboardStats(entries []LeaderboardEntry) map[string]interface{} {
	if len(entries) == 0 {
		return map[string]interface{}{}
	}
	
	stats := make(map[string]interface{})
	
	// Calculate averages
	totalXP := uint64(0)
	totalEarnings := sdk.ZeroInt()
	totalAchievements := uint32(0)
	totalStreak := uint32(0)
	
	// Avatar distribution
	avatarCount := make(map[AvatarType]int)
	
	for _, entry := range entries {
		totalXP += entry.TotalXp
		totalEarnings = totalEarnings.Add(entry.TotalEarnings.Amount)
		totalAchievements += entry.AchievementsCount
		totalStreak += entry.CurrentStreak
		avatarCount[entry.ActiveAvatar]++
	}
	
	count := len(entries)
	stats["total_developers"] = count
	stats["average_xp"] = totalXP / uint64(count)
	stats["average_earnings"] = totalEarnings.Quo(sdk.NewInt(int64(count)))
	stats["average_achievements"] = totalAchievements / uint32(count)
	stats["average_streak"] = totalStreak / uint32(count)
	stats["avatar_distribution"] = avatarCount
	
	// Top performer stats
	if len(entries) > 0 {
		stats["top_xp"] = entries[0].TotalXp
		stats["top_level"] = entries[0].Level
		stats["top_earnings"] = entries[0].TotalEarnings
	}
	
	return stats
}

// GenerateLeaderboardUpdateMessage creates update message
func GenerateLeaderboardUpdateMessage(oldRank, newRank uint64, username string, leaderboardType string) string {
	if oldRank == 0 {
		// New entry
		return fmt.Sprintf(`
🎊 NEW ENTRY ALERT! 🎊
@%s has entered the %s leaderboard at rank #%d!
Welcome to the competition! May the best coder win! 🚀`,
			username, leaderboardType, newRank)
	}
	
	if newRank < oldRank {
		// Moved up
		improvement := oldRank - newRank
		return fmt.Sprintf(`
📈 RANK UP! @%s climbed %d positions!
Old Rank: #%d → New Rank: #%d
%s leaderboard mein tarakki! Keep climbing! 🧗`,
			username, improvement, oldRank, newRank, leaderboardType)
	} else if newRank > oldRank {
		// Moved down
		decline := newRank - oldRank
		return fmt.Sprintf(`
📉 Rank changed: @%s moved down %d positions
Old Rank: #%d → New Rank: #%d
Competition is fierce! Time to code harder! 💪`,
			username, decline, oldRank, newRank)
	}
	
	// Same rank
	return fmt.Sprintf(`
➡️ @%s maintained rank #%d on %s leaderboard!
Consistency is key! Keep it up! 👍`,
		username, newRank, leaderboardType)
}