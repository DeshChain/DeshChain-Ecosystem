package types

import (
	"fmt"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// Event types
	EventTypeCreateProfile        = "create_profile"
	EventTypeUpdateProfile        = "update_profile"
	EventTypeSelectAvatar         = "select_avatar"
	EventTypeAchievementUnlocked  = "achievement_unlocked"
	EventTypeLevelUp              = "level_up"
	EventTypeStreakUpdate         = "streak_update"
	EventTypeSocialShare          = "social_share"
	EventTypeTeamBattleStart      = "team_battle_start"
	EventTypeTeamBattleEnd        = "team_battle_end"
	EventTypeDailyChallengeComplete = "daily_challenge_complete"
	EventTypeViralPost            = "viral_post"
	EventTypeLeaderboardUpdate    = "leaderboard_update"
	
	// Attribute keys
	AttributeKeyAddress          = "address"
	AttributeKeyUsername         = "username"
	AttributeKeyAvatar           = "avatar"
	AttributeKeyAchievementID    = "achievement_id"
	AttributeKeyAchievementName  = "achievement_name"
	AttributeKeyRarity           = "rarity"
	AttributeKeyReward           = "reward"
	AttributeKeyXP               = "xp"
	AttributeKeyOldLevel         = "old_level"
	AttributeKeyNewLevel         = "new_level"
	AttributeKeyLevelTitle       = "level_title"
	AttributeKeyStreakDays       = "streak_days"
	AttributeKeyStreakMaintained = "streak_maintained"
	AttributeKeyPlatform         = "platform"
	AttributeKeyPostID           = "post_id"
	AttributeKeyEngagement       = "engagement"
	AttributeKeyBattleID         = "battle_id"
	AttributeKeyBattleType       = "battle_type"
	AttributeKeyTeamName         = "team_name"
	AttributeKeyWinner           = "winner"
	AttributeKeyChallengeType    = "challenge_type"
	AttributeKeyOldRank          = "old_rank"
	AttributeKeyNewRank          = "new_rank"
	AttributeKeyLeaderboardType  = "leaderboard_type"
	
	// Module name for events
	ModuleName = "gamification"
)

// EmitAchievementUnlockedEvent emits an achievement unlocked event
func EmitAchievementUnlockedEvent(ctx sdk.Context, profile *DeveloperProfile, achievement *Achievement) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			EventTypeAchievementUnlocked,
			sdk.NewAttribute(AttributeKeyAddress, profile.Address),
			sdk.NewAttribute(AttributeKeyUsername, profile.GithubUsername),
			sdk.NewAttribute(AttributeKeyAchievementID, fmt.Sprintf("%d", achievement.AchievementId)),
			sdk.NewAttribute(AttributeKeyAchievementName, achievement.Name),
			sdk.NewAttribute(AttributeKeyRarity, achievement.RarityLevel.String()),
			sdk.NewAttribute(AttributeKeyReward, achievement.RewardAmount.String()),
			sdk.NewAttribute(AttributeKeyXP, fmt.Sprintf("%d", achievement.XpReward)),
		),
	)
}

// EmitLevelUpEvent emits a level up event
func EmitLevelUpEvent(ctx sdk.Context, profile *DeveloperProfile, oldLevel, newLevel uint32, levelTitle string) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			EventTypeLevelUp,
			sdk.NewAttribute(AttributeKeyAddress, profile.Address),
			sdk.NewAttribute(AttributeKeyUsername, profile.GithubUsername),
			sdk.NewAttribute(AttributeKeyOldLevel, fmt.Sprintf("%d", oldLevel)),
			sdk.NewAttribute(AttributeKeyNewLevel, fmt.Sprintf("%d", newLevel)),
			sdk.NewAttribute(AttributeKeyLevelTitle, levelTitle),
		),
	)
}

// EmitStreakUpdateEvent emits a streak update event
func EmitStreakUpdateEvent(ctx sdk.Context, profile *DeveloperProfile, maintained bool) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			EventTypeStreakUpdate,
			sdk.NewAttribute(AttributeKeyAddress, profile.Address),
			sdk.NewAttribute(AttributeKeyUsername, profile.GithubUsername),
			sdk.NewAttribute(AttributeKeyStreakDays, fmt.Sprintf("%d", profile.CurrentStreakDays)),
			sdk.NewAttribute(AttributeKeyStreakMaintained, fmt.Sprintf("%t", maintained)),
		),
	)
}

// EmitSocialShareEvent emits a social media share event
func EmitSocialShareEvent(ctx sdk.Context, post *SocialMediaPost) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			EventTypeSocialShare,
			sdk.NewAttribute(AttributeKeyAddress, post.DeveloperAddress),
			sdk.NewAttribute(AttributeKeyPostID, fmt.Sprintf("%d", post.PostId)),
			sdk.NewAttribute(AttributeKeyPlatform, post.Platform.String()),
			sdk.NewAttribute(AttributeKeyAchievementID, fmt.Sprintf("%d", post.AchievementId)),
		),
	)
}

// EmitViralPostEvent emits a viral post event
func EmitViralPostEvent(ctx sdk.Context, post *SocialMediaPost, engagementScore uint32) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			EventTypeViralPost,
			sdk.NewAttribute(AttributeKeyAddress, post.DeveloperAddress),
			sdk.NewAttribute(AttributeKeyPostID, fmt.Sprintf("%d", post.PostId)),
			sdk.NewAttribute(AttributeKeyPlatform, post.Platform.String()),
			sdk.NewAttribute(AttributeKeyEngagement, fmt.Sprintf("%d", engagementScore)),
		),
	)
}

// EmitTeamBattleStartEvent emits a team battle start event
func EmitTeamBattleStartEvent(ctx sdk.Context, battle *TeamBattle) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			EventTypeTeamBattleStart,
			sdk.NewAttribute(AttributeKeyBattleID, fmt.Sprintf("%d", battle.BattleId)),
			sdk.NewAttribute(AttributeKeyBattleType, battle.BattleType),
			sdk.NewAttribute(sdk.AttributeKeyAmount, battle.PrizePool.String()),
		),
	)
}

// EmitTeamBattleEndEvent emits a team battle end event
func EmitTeamBattleEndEvent(ctx sdk.Context, battle *TeamBattle, winner string) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			EventTypeTeamBattleEnd,
			sdk.NewAttribute(AttributeKeyBattleID, fmt.Sprintf("%d", battle.BattleId)),
			sdk.NewAttribute(AttributeKeyWinner, winner),
			sdk.NewAttribute(sdk.AttributeKeyAmount, battle.PrizePool.String()),
		),
	)
}

// EmitDailyChallengeCompleteEvent emits a daily challenge completion event
func EmitDailyChallengeCompleteEvent(ctx sdk.Context, address string, challenge *DailyChallenge) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			EventTypeDailyChallengeComplete,
			sdk.NewAttribute(AttributeKeyAddress, address),
			sdk.NewAttribute(AttributeKeyChallengeType, challenge.ChallengeType),
			sdk.NewAttribute(AttributeKeyReward, challenge.Reward.String()),
		),
	)
}

// EmitLeaderboardUpdateEvent emits a leaderboard rank change event
func EmitLeaderboardUpdateEvent(ctx sdk.Context, address string, leaderboardType string, oldRank, newRank uint64) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			EventTypeLeaderboardUpdate,
			sdk.NewAttribute(AttributeKeyAddress, address),
			sdk.NewAttribute(AttributeKeyLeaderboardType, leaderboardType),
			sdk.NewAttribute(AttributeKeyOldRank, fmt.Sprintf("%d", oldRank)),
			sdk.NewAttribute(AttributeKeyNewRank, fmt.Sprintf("%d", newRank)),
		),
	)
}

// EmitProfileCreatedEvent emits a profile created event
func EmitProfileCreatedEvent(ctx sdk.Context, profile *DeveloperProfile) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			EventTypeCreateProfile,
			sdk.NewAttribute(AttributeKeyAddress, profile.Address),
			sdk.NewAttribute(AttributeKeyUsername, profile.GithubUsername),
			sdk.NewAttribute(AttributeKeyAvatar, profile.ActiveAvatar.String()),
		),
	)
}

// EmitAvatarSelectedEvent emits an avatar selection event
func EmitAvatarSelectedEvent(ctx sdk.Context, address string, avatar AvatarType) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			EventTypeSelectAvatar,
			sdk.NewAttribute(AttributeKeyAddress, address),
			sdk.NewAttribute(AttributeKeyAvatar, avatar.String()),
		),
	)
}