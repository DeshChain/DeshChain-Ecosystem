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
	"time"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NotificationType defines types of notifications
type NotificationType string

const (
	NotificationTypeAchievementUnlocked NotificationType = "achievement_unlocked"
	NotificationTypeLevelUp             NotificationType = "level_up"
	NotificationTypeStreakMaintained    NotificationType = "streak_maintained"
	NotificationTypeStreakBroken        NotificationType = "streak_broken"
	NotificationTypeSocialShare         NotificationType = "social_share"
	NotificationTypeTeamBattleStart     NotificationType = "team_battle_start"
	NotificationTypeTeamBattleEnd       NotificationType = "team_battle_end"
	NotificationTypeDailyChallenge      NotificationType = "daily_challenge"
	NotificationTypeViralPost           NotificationType = "viral_post"
	NotificationTypeLeaderboardChange   NotificationType = "leaderboard_change"
	NotificationTypeMilestone           NotificationType = "milestone"
)

// NotificationPriority defines notification priority levels
type NotificationPriority string

const (
	NotificationPriorityHigh   NotificationPriority = "high"
	NotificationPriorityMedium NotificationPriority = "medium"
	NotificationPriorityLow    NotificationPriority = "low"
)

// Notification represents a gamification notification
type Notification struct {
	ID           string               `json:"id"`
	Type         NotificationType     `json:"type"`
	Priority     NotificationPriority `json:"priority"`
	Title        string               `json:"title"`
	Message      string               `json:"message"`
	Timestamp    time.Time            `json:"timestamp"`
	UserAddress  string               `json:"user_address"`
	Data         map[string]interface{} `json:"data"`
	Read         bool                 `json:"read"`
	ActionURL    string               `json:"action_url,omitempty"`
	ImageURL     string               `json:"image_url,omitempty"`
	ExpiresAt    *time.Time           `json:"expires_at,omitempty"`
}

// NotificationManager handles notification creation and delivery
type NotificationManager struct {
	notifications []Notification
}

// NewNotificationManager creates a new notification manager
func NewNotificationManager() *NotificationManager {
	return &NotificationManager{
		notifications: []Notification{},
	}
}

// CreateAchievementNotification creates notification for achievement unlock
func (nm *NotificationManager) CreateAchievementNotification(
	profile *DeveloperProfile,
	achievement *Achievement,
) Notification {
	
	title := fmt.Sprintf("🎊 Achievement Unlocked: %s!", achievement.Name)
	message := fmt.Sprintf(
		"%s\n\n💰 Reward: %s NAMO\n⭐ XP: %d\n\n%s",
		achievement.Description,
		achievement.RewardAmount.Amount,
		achievement.XpReward,
		achievement.UnlockQuote,
	)
	
	notification := Notification{
		ID:          fmt.Sprintf("achievement_%d_%d", achievement.AchievementId, time.Now().Unix()),
		Type:        NotificationTypeAchievementUnlocked,
		Priority:    getPriorityForRarity(achievement.RarityLevel),
		Title:       title,
		Message:     message,
		Timestamp:   time.Now(),
		UserAddress: profile.Address,
		Data: map[string]interface{}{
			"achievement_id": achievement.AchievementId,
			"achievement_name": achievement.Name,
			"rarity": achievement.RarityLevel.String(),
			"reward": achievement.RewardAmount,
			"xp": achievement.XpReward,
			"badge_url": achievement.BadgeImageUrl,
		},
		ActionURL: fmt.Sprintf("/achievement/%d", achievement.AchievementId),
		ImageURL: achievement.BadgeImageUrl,
	}
	
	return notification
}

// CreateLevelUpNotification creates notification for level up
func (nm *NotificationManager) CreateLevelUpNotification(
	profile *DeveloperProfile,
	oldLevel, newLevel uint32,
	bonusReward sdk.Coin,
) Notification {
	
	levelConfig := GetLevelConfigByLevel(newLevel)
	
	title := fmt.Sprintf("🎯 Level Up! Welcome to Level %d!", newLevel)
	message := fmt.Sprintf(
		"Congratulations %s!\n\nYou've reached Level %d: %s\n\n💰 Bonus Reward: %s NAMO\n\n%s",
		profile.GithubUsername,
		newLevel,
		levelConfig.Title,
		bonusReward.Amount,
		levelConfig.UnlockQuote,
	)
	
	// Special celebration for milestone levels
	priority := NotificationPriorityMedium
	if newLevel%10 == 0 || newLevel == 25 || newLevel == 50 || newLevel == 75 || newLevel == 100 {
		priority = NotificationPriorityHigh
	}
	
	notification := Notification{
		ID:          fmt.Sprintf("levelup_%d_%d", newLevel, time.Now().Unix()),
		Type:        NotificationTypeLevelUp,
		Priority:    priority,
		Title:       title,
		Message:     message,
		Timestamp:   time.Now(),
		UserAddress: profile.Address,
		Data: map[string]interface{}{
			"old_level": oldLevel,
			"new_level": newLevel,
			"level_title": levelConfig.Title,
			"bonus_reward": bonusReward,
		},
		ActionURL: fmt.Sprintf("/profile/%s", profile.Address),
	}
	
	return notification
}

// CreateStreakNotification creates notification for streak events
func (nm *NotificationManager) CreateStreakNotification(
	profile *DeveloperProfile,
	maintained bool,
	days uint32,
) Notification {
	
	var title, message string
	notifType := NotificationTypeStreakMaintained
	priority := NotificationPriorityLow
	
	if maintained {
		title = fmt.Sprintf("🔥 %d Day Streak Maintained!", days)
		message = GetStreakMessage(days, profile.GithubUsername)
		
		// Higher priority for milestone streaks
		if days == 7 || days == 30 || days == 100 || days == 365 {
			priority = NotificationPriorityHigh
		}
	} else {
		notifType = NotificationTypeStreakBroken
		title = "😔 Streak Broken!"
		message = fmt.Sprintf(
			"Your %d day streak has ended. Don't worry %s, start a new one today!\n\n💪 Remember: Every expert was once a beginner!",
			days,
			profile.GithubUsername,
		)
		priority = NotificationPriorityMedium
	}
	
	notification := Notification{
		ID:          fmt.Sprintf("streak_%s_%d", profile.Address, time.Now().Unix()),
		Type:        notifType,
		Priority:    priority,
		Title:       title,
		Message:     message,
		Timestamp:   time.Now(),
		UserAddress: profile.Address,
		Data: map[string]interface{}{
			"streak_days": days,
			"maintained": maintained,
		},
	}
	
	return notification
}

// CreateViralPostNotification creates notification when post goes viral
func (nm *NotificationManager) CreateViralPostNotification(
	profile *DeveloperProfile,
	post *SocialMediaPost,
) Notification {
	
	title := "🚀 Your Post Went Viral!"
	message := fmt.Sprintf(
		"Your achievement post is trending!\n\n👍 %d Likes\n🔄 %d Shares\n💬 %d Comments\n👁️ %d Views\n\nYou're now a DeshChain influencer!",
		post.EngagementMetrics.Likes,
		post.EngagementMetrics.Shares,
		post.EngagementMetrics.Comments,
		post.EngagementMetrics.Views,
	)
	
	notification := Notification{
		ID:          fmt.Sprintf("viral_%d", post.PostId),
		Type:        NotificationTypeViralPost,
		Priority:    NotificationPriorityHigh,
		Title:       title,
		Message:     message,
		Timestamp:   time.Now(),
		UserAddress: profile.Address,
		Data: map[string]interface{}{
			"post_id": post.PostId,
			"platform": post.Platform.String(),
			"engagement": post.EngagementMetrics,
		},
		ActionURL: fmt.Sprintf("/post/%d", post.PostId),
	}
	
	return notification
}

// CreateTeamBattleNotification creates notification for team battles
func (nm *NotificationManager) CreateTeamBattleNotification(
	userAddress string,
	battle *TeamBattle,
	isStart bool,
) Notification {
	
	var title, message string
	notifType := NotificationTypeTeamBattleStart
	
	if isStart {
		title = fmt.Sprintf("🏏 Team Battle Started: %s!", battle.BattleType)
		message = fmt.Sprintf(
			"%s vs %s\n\n💰 Prize Pool: %s NAMO\n⏱️ Duration: %s\n\nMay the best team win!",
			battle.Team1Name,
			battle.Team2Name,
			battle.PrizePool.Amount,
			formatDuration(time.Duration(battle.EndTime.Sub(battle.StartTime).Seconds()) * time.Second),
		)
	} else {
		notifType = NotificationTypeTeamBattleEnd
		winner := battle.Team1Name
		if battle.Team2Score > battle.Team1Score {
			winner = battle.Team2Name
		}
		
		title = fmt.Sprintf("🏆 Team Battle Ended: %s Wins!", winner)
		message = fmt.Sprintf(
			"Final Score:\n%s: %d\n%s: %d\n\n🎊 Congratulations to %s!\n💰 Prize: %s NAMO",
			battle.Team1Name, battle.Team1Score,
			battle.Team2Name, battle.Team2Score,
			winner,
			battle.PrizePool.Amount,
		)
	}
	
	notification := Notification{
		ID:          fmt.Sprintf("battle_%d_%d", battle.BattleId, time.Now().Unix()),
		Type:        notifType,
		Priority:    NotificationPriorityHigh,
		Title:       title,
		Message:     message,
		Timestamp:   time.Now(),
		UserAddress: userAddress,
		Data: map[string]interface{}{
			"battle_id": battle.BattleId,
			"battle_type": battle.BattleType,
			"is_start": isStart,
		},
		ActionURL: fmt.Sprintf("/battle/%d", battle.BattleId),
	}
	
	return notification
}

// CreateDailyChallengeNotification creates notification for daily challenge
func (nm *NotificationManager) CreateDailyChallengeNotification(
	userAddress string,
	challenge *DailyChallenge,
) Notification {
	
	title := "🎯 New Daily Challenge Available!"
	message := fmt.Sprintf(
		"%s\n\n📋 Goal: %d\n💰 Reward: %s NAMO\n\n%s",
		challenge.Description,
		challenge.TargetValue,
		challenge.Reward.Amount,
		challenge.ThemeQuote,
	)
	
	// Set expiry for daily challenge notifications
	tomorrow := time.Now().Add(24 * time.Hour)
	
	notification := Notification{
		ID:          fmt.Sprintf("daily_%s", challenge.ChallengeDate.Format("2006-01-02")),
		Type:        NotificationTypeDailyChallenge,
		Priority:    NotificationPriorityMedium,
		Title:       title,
		Message:     message,
		Timestamp:   time.Now(),
		UserAddress: userAddress,
		Data: map[string]interface{}{
			"challenge_type": challenge.ChallengeType,
			"target": challenge.TargetValue,
			"reward": challenge.Reward,
		},
		ActionURL: "/daily-challenge",
		ExpiresAt: &tomorrow,
	}
	
	return notification
}

// CreateLeaderboardChangeNotification creates notification for rank changes
func (nm *NotificationManager) CreateLeaderboardChangeNotification(
	profile *DeveloperProfile,
	oldRank, newRank uint64,
	leaderboardType string,
) Notification {
	
	improved := newRank < oldRank
	emoji := "📈"
	if !improved {
		emoji = "📉"
	}
	
	title := fmt.Sprintf("%s Leaderboard Rank Changed!", emoji)
	message := GetRankChangeMessage(profile.GithubUsername, oldRank, newRank, "up")
	if !improved {
		message = GetRankChangeMessage(profile.GithubUsername, oldRank, newRank, "down")
	}
	
	priority := NotificationPriorityLow
	// Higher priority for top 10 changes
	if newRank <= 10 || oldRank <= 10 {
		priority = NotificationPriorityMedium
	}
	// Highest priority for podium changes
	if newRank <= 3 || oldRank <= 3 {
		priority = NotificationPriorityHigh
	}
	
	notification := Notification{
		ID:          fmt.Sprintf("rank_%s_%d", profile.Address, time.Now().Unix()),
		Type:        NotificationTypeLeaderboardChange,
		Priority:    priority,
		Title:       title,
		Message:     message,
		Timestamp:   time.Now(),
		UserAddress: profile.Address,
		Data: map[string]interface{}{
			"old_rank": oldRank,
			"new_rank": newRank,
			"leaderboard_type": leaderboardType,
			"improved": improved,
		},
		ActionURL: fmt.Sprintf("/leaderboard/%s", leaderboardType),
	}
	
	return notification
}

// Helper functions

func getPriorityForRarity(rarity RarityLevel) NotificationPriority {
	switch rarity {
	case RarityLevel_RARITY_LEVEL_MYTHIC:
		return NotificationPriorityHigh
	case RarityLevel_RARITY_LEVEL_LEGENDARY:
		return NotificationPriorityHigh
	case RarityLevel_RARITY_LEVEL_EPIC:
		return NotificationPriorityMedium
	case RarityLevel_RARITY_LEVEL_RARE:
		return NotificationPriorityMedium
	default:
		return NotificationPriorityLow
	}
}

// GetStreakMessage returns appropriate message for streak days
func GetStreakMessage(days uint32, username string) string {
	switch days {
	case 7:
		return fmt.Sprintf("🔥 SAAT DIN SAAT COMMIT! %s ne lagayi aag! Keep the fire burning!", username)
	case 30:
		return fmt.Sprintf("📿 30 DIN KA TAPASYA! %s = Modern day Coding Rishi! Respect! 🙏", username)
	case 100:
		return fmt.Sprintf("💯 CENTURY MUBARAK! %s = Sachin of Coding! Standing ovation! 👏", username)
	case 365:
		return fmt.Sprintf("🎊 SAAL BHAR CODING! %s completed 365 days! Legendary dedication! 🙌", username)
	default:
		return fmt.Sprintf("🔥 %d Day Streak! %s is on fire! Keep coding, keep growing!", days, username)
	}
}

// NotificationPreferences stores user notification preferences
type NotificationPreferences struct {
	UserAddress             string `json:"user_address"`
	AchievementNotifs       bool   `json:"achievement_notifs"`
	LevelUpNotifs           bool   `json:"level_up_notifs"`
	StreakNotifs            bool   `json:"streak_notifs"`
	SocialShareNotifs       bool   `json:"social_share_notifs"`
	TeamBattleNotifs        bool   `json:"team_battle_notifs"`
	DailyChallengeNotifs    bool   `json:"daily_challenge_notifs"`
	ViralPostNotifs         bool   `json:"viral_post_notifs"`
	LeaderboardChangeNotifs bool   `json:"leaderboard_change_notifs"`
	EmailNotifications      bool   `json:"email_notifications"`
	PushNotifications       bool   `json:"push_notifications"`
	SMSNotifications        bool   `json:"sms_notifications"`
	QuietHoursStart         string `json:"quiet_hours_start"` // e.g., "22:00"
	QuietHoursEnd           string `json:"quiet_hours_end"`   // e.g., "08:00"
}

// ShouldSendNotification checks if notification should be sent based on preferences
func ShouldSendNotification(prefs *NotificationPreferences, notifType NotificationType) bool {
	switch notifType {
	case NotificationTypeAchievementUnlocked:
		return prefs.AchievementNotifs
	case NotificationTypeLevelUp:
		return prefs.LevelUpNotifs
	case NotificationTypeStreakMaintained, NotificationTypeStreakBroken:
		return prefs.StreakNotifs
	case NotificationTypeSocialShare:
		return prefs.SocialShareNotifs
	case NotificationTypeTeamBattleStart, NotificationTypeTeamBattleEnd:
		return prefs.TeamBattleNotifs
	case NotificationTypeDailyChallenge:
		return prefs.DailyChallengeNotifs
	case NotificationTypeViralPost:
		return prefs.ViralPostNotifs
	case NotificationTypeLeaderboardChange:
		return prefs.LeaderboardChangeNotifs
	default:
		return true
	}
}