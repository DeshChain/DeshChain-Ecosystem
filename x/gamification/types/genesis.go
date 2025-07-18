package types

import (
	"fmt"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
		Profiles: []DeveloperProfile{},
		Achievements: GetDefaultAchievements(),
		HumorQuotes: GetDefaultHumorQuotes(),
		SocialPosts: []SocialMediaPost{},
		TeamBattles: []TeamBattle{},
		DailyChallenges: []DailyChallenge{},
		AchievementCards: []AchievementCard{},
		LevelConfigs: GetLevelConfigs(),
		NextAchievementId: uint64(len(GetDefaultAchievements()) + 1),
		NextPostId: 1,
		NextBattleId: 1,
		NextQuoteId: uint64(len(GetDefaultHumorQuotes()) + 1),
		NextCardId: 1,
	}
}

// DefaultParams returns default gamification parameters
func DefaultParams() GamificationParams {
	return GamificationParams{
		EnableGamification: true,
		BaseXpPerCommit: 10,
		BaseXpPerBugFix: 50,
		BaseXpPerFeature: 100,
		BaseXpPerDoc: 30,
		StreakBonusMultiplier: "1.1", // 10% bonus
		ViralThreshold: 1000, // 1000 engagement for viral
		EnableSocialPosting: true,
		MaxDailyAchievements: 10,
		QuoteRefreshHours: 24,
	}
}

// ValidateGenesis validates the genesis state
func ValidateGenesis(data GenesisState) error {
	// Validate params
	if err := validateParams(data.Params); err != nil {
		return err
	}

	// Validate profiles
	profileMap := make(map[string]bool)
	usernameMap := make(map[string]bool)
	for _, profile := range data.Profiles {
		if _, ok := profileMap[profile.Address]; ok {
			return fmt.Errorf("duplicate profile address %s", profile.Address)
		}
		profileMap[profile.Address] = true
		
		if _, ok := usernameMap[profile.GithubUsername]; ok {
			return fmt.Errorf("duplicate github username %s", profile.GithubUsername)
		}
		usernameMap[profile.GithubUsername] = true
		
		if err := host.AddressFromBech32(profile.Address); err != nil {
			return fmt.Errorf("invalid profile address %s: %w", profile.Address, err)
		}
	}

	// Validate achievements
	achievementMap := make(map[uint64]bool)
	for _, achievement := range data.Achievements {
		if _, ok := achievementMap[achievement.AchievementId]; ok {
			return fmt.Errorf("duplicate achievement id %d", achievement.AchievementId)
		}
		achievementMap[achievement.AchievementId] = true
		
		if achievement.AchievementId == 0 {
			return fmt.Errorf("achievement id cannot be zero")
		}
		
		if achievement.Name == "" {
			return fmt.Errorf("achievement name cannot be empty")
		}
		
		if !achievement.RewardAmount.IsValid() || achievement.RewardAmount.IsZero() {
			return fmt.Errorf("invalid achievement reward amount")
		}
	}

	// Validate humor quotes
	quoteMap := make(map[uint64]bool)
	for _, quote := range data.HumorQuotes {
		if _, ok := quoteMap[quote.QuoteId]; ok {
			return fmt.Errorf("duplicate quote id %d", quote.QuoteId)
		}
		quoteMap[quote.QuoteId] = true
		
		if quote.QuoteId == 0 {
			return fmt.Errorf("quote id cannot be zero")
		}
		
		if quote.Text == "" {
			return fmt.Errorf("quote text cannot be empty")
		}
	}

	// Validate social posts
	postMap := make(map[uint64]bool)
	for _, post := range data.SocialPosts {
		if _, ok := postMap[post.PostId]; ok {
			return fmt.Errorf("duplicate post id %d", post.PostId)
		}
		postMap[post.PostId] = true
		
		if post.PostId == 0 {
			return fmt.Errorf("post id cannot be zero")
		}
		
		if err := host.AddressFromBech32(post.DeveloperAddress); err != nil {
			return fmt.Errorf("invalid developer address in post %d: %w", post.PostId, err)
		}
	}

	// Validate team battles
	battleMap := make(map[uint64]bool)
	for _, battle := range data.TeamBattles {
		if _, ok := battleMap[battle.BattleId]; ok {
			return fmt.Errorf("duplicate battle id %d", battle.BattleId)
		}
		battleMap[battle.BattleId] = true
		
		if battle.BattleId == 0 {
			return fmt.Errorf("battle id cannot be zero")
		}
		
		if !battle.PrizePool.IsValid() || battle.PrizePool.IsZero() {
			return fmt.Errorf("invalid battle prize pool")
		}
	}

	// Validate level configs
	levelMap := make(map[uint32]bool)
	for _, config := range data.LevelConfigs {
		if _, ok := levelMap[config.Level]; ok {
			return fmt.Errorf("duplicate level config %d", config.Level)
		}
		levelMap[config.Level] = true
		
		if config.Level == 0 || config.Level > 100 {
			return fmt.Errorf("invalid level %d (must be 1-100)", config.Level)
		}
		
		if !config.BonusReward.IsValid() || config.BonusReward.IsZero() {
			return fmt.Errorf("invalid level bonus reward")
		}
	}

	// Validate next IDs
	if data.NextAchievementId <= uint64(len(data.Achievements)) {
		return fmt.Errorf("next achievement id must be greater than existing achievements")
	}
	
	if data.NextQuoteId <= uint64(len(data.HumorQuotes)) {
		return fmt.Errorf("next quote id must be greater than existing quotes")
	}

	return nil
}

// validateParams validates the module parameters
func validateParams(params GamificationParams) error {
	if params.BaseXpPerCommit == 0 {
		return fmt.Errorf("base XP per commit cannot be zero")
	}
	
	if params.BaseXpPerBugFix == 0 {
		return fmt.Errorf("base XP per bug fix cannot be zero")
	}
	
	if params.BaseXpPerFeature == 0 {
		return fmt.Errorf("base XP per feature cannot be zero")
	}
	
	if params.BaseXpPerDoc == 0 {
		return fmt.Errorf("base XP per doc cannot be zero")
	}
	
	// Validate streak multiplier is a valid decimal
	if _, err := sdk.NewDecFromStr(params.StreakBonusMultiplier); err != nil {
		return fmt.Errorf("invalid streak bonus multiplier: %w", err)
	}
	
	if params.ViralThreshold == 0 {
		return fmt.Errorf("viral threshold cannot be zero")
	}
	
	if params.MaxDailyAchievements == 0 {
		return fmt.Errorf("max daily achievements cannot be zero")
	}
	
	if params.QuoteRefreshHours == 0 {
		return fmt.Errorf("quote refresh hours cannot be zero")
	}
	
	return nil
}

// GetDefaultHumorQuotes returns default humor quotes with IDs
func GetDefaultHumorQuotes() []HumorQuote {
	humorEngine := NewHumorEngine()
	quotes := []HumorQuote{}
	id := uint64(1)
	
	// Add Bollywood dialogues
	for _, dialogue := range BollywoodDialogues {
		quotes = append(quotes, HumorQuote{
			QuoteId:            id,
			QuoteType:          QuoteType_QUOTE_TYPE_BOLLYWOOD_DIALOGUE,
			Text:               dialogue.Text,
			EnglishTranslation: dialogue.Translation,
			Source:             dialogue.Source,
			Character:          dialogue.Character,
			SuitableFor:        dialogue.Categories,
			UsageCount:         0,
			ViralScore:         0,
			IsFamilyFriendly:   true,
			RegionalTags:       dialogue.Tags,
		})
		id++
	}
	
	// Add Cricket commentary
	for _, commentary := range CricketCommentary {
		quotes = append(quotes, HumorQuote{
			QuoteId:            id,
			QuoteType:          QuoteType_QUOTE_TYPE_CRICKET_COMMENTARY,
			Text:               commentary.Text,
			EnglishTranslation: "",
			Source:             commentary.Source,
			Character:          commentary.Character,
			SuitableFor:        commentary.Categories,
			UsageCount:         0,
			ViralScore:         0,
			IsFamilyFriendly:   true,
			RegionalTags:       commentary.Tags,
		})
		id++
	}
	
	// Add South Indian quotes
	for _, quote := range SouthIndianQuotes {
		quotes = append(quotes, HumorQuote{
			QuoteId:            id,
			QuoteType:          QuoteType_QUOTE_TYPE_SOUTH_INDIAN_SUPERSTAR,
			Text:               quote.Text,
			EnglishTranslation: quote.Translation,
			Source:             quote.Source,
			Character:          quote.Character,
			SuitableFor:        quote.Categories,
			UsageCount:         0,
			ViralScore:         0,
			IsFamilyFriendly:   true,
			RegionalTags:       quote.Tags,
		})
		id++
	}
	
	// Add Comedy punchlines
	for _, punchline := range ComedyPunchlines {
		quotes = append(quotes, HumorQuote{
			QuoteId:            id,
			QuoteType:          QuoteType_QUOTE_TYPE_COMEDY_PUNCHLINE,
			Text:               punchline.Text,
			EnglishTranslation: punchline.Translation,
			Source:             punchline.Source,
			Character:          punchline.Character,
			SuitableFor:        punchline.Categories,
			UsageCount:         0,
			ViralScore:         0,
			IsFamilyFriendly:   true,
			RegionalTags:       punchline.Tags,
		})
		id++
	}
	
	return quotes
}