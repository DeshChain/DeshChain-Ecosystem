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

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SocialMediaGenerator creates viral-worthy posts
type SocialMediaGenerator struct {
	humorEngine *HumorEngine
}

// NewSocialMediaGenerator creates a new generator
func NewSocialMediaGenerator() *SocialMediaGenerator {
	return &SocialMediaGenerator{
		humorEngine: NewHumorEngine(),
	}
}

// GenerateAchievementPost creates a platform-specific post
func (g *SocialMediaGenerator) GenerateAchievementPost(
	profile *DeveloperProfile,
	achievement *Achievement,
	platform SocialPlatform,
) *SocialMediaPost {
	// Get suitable quote
	quote := g.humorEngine.GetQuoteForAchievement(achievement.Category, profile.HumorPreference)

	// Generate content based on platform
	var content string
	var hashtags []string
	var imageUrl string

	switch platform {
	case SocialPlatform_SOCIAL_PLATFORM_TWITTER:
		content = g.generateTwitterPost(profile, achievement, quote)
		hashtags = GetHashtags(achievement.Category.String())
	case SocialPlatform_SOCIAL_PLATFORM_DISCORD:
		content = g.generateDiscordPost(profile, achievement, quote)
		hashtags = GetHashtags(achievement.Category.String())
	case SocialPlatform_SOCIAL_PLATFORM_WHATSAPP:
		content = g.generateWhatsAppPost(profile, achievement, quote)
	case SocialPlatform_SOCIAL_PLATFORM_INSTAGRAM:
		content = g.generateInstagramPost(profile, achievement, quote)
		hashtags = append(GetHashtags(achievement.Category.String()), "#InstaCode", "#DevLife")
		imageUrl = g.generatePosterUrl(profile, achievement)
	case SocialPlatform_SOCIAL_PLATFORM_LINKEDIN:
		content = g.generateLinkedInPost(profile, achievement, quote)
		hashtags = []string{"#DeshChain", "#Blockchain", "#Achievement", "#TechIndia"}
	}

	// Create post object
	post := &SocialMediaPost{
		PostId:           uint64(time.Now().UnixNano()),
		DeveloperAddress: profile.Address,
		AchievementId:    achievement.AchievementId,
		Platform:         platform,
		Content:          content,
		ImageUrl:         imageUrl,
		Hashtags:         hashtags,
		PostTime:         time.Now(),
		EngagementMetrics: &EngagementMetrics{
			IsViral: IsViralWorthy(content),
		},
	}

	if quote != nil {
		post.QuoteIds = []uint64{quote.QuoteId}
	}

	return post
}

// generateTwitterPost creates Twitter-style post (280 chars)
func (g *SocialMediaGenerator) generateTwitterPost(profile *DeveloperProfile, achievement *Achievement, quote *HumorQuote) string {
	// Short and punchy for Twitter
	templates := []string{
		"🎬 %s unlocked: %s!\n%s\n💰 %s earned!",
		"🏏 SIXER! %s achieved %s!\n%s\n🎯 %s NAMO!",
		"🌟 %s ne kar diya kamaal!\n✅ %s\n%s",
		"💥 BREAKING: %s rocks %s!\n%s\n₹%s collection!",
	}

	quoteText := ""
	if quote != nil {
		quoteText = quote.Text
		if len(quoteText) > 100 {
			quoteText = quoteText[:97] + "..."
		}
	}

	template := templates[rand.Intn(len(templates))]
	content := fmt.Sprintf(template,
		profile.GithubUsername,
		achievement.Name,
		quoteText,
		achievement.RewardAmount.Amount,
	)

	// Ensure within Twitter limit
	if len(content) > 250 { // Leave room for hashtags
		content = content[:247] + "..."
	}

	return content
}

// generateDiscordPost creates Discord-style post (detailed)
func (g *SocialMediaGenerator) generateDiscordPost(profile *DeveloperProfile, achievement *Achievement, quote *HumorQuote) string {
	avatar := GetAvatarByType(profile.ActiveAvatar)
	
	content := fmt.Sprintf(`🎊 **ACHIEVEMENT UNLOCKED** 🎊
━━━━━━━━━━━━━━━━━━━━━━━━━━━
👤 **Developer**: %s
🎭 **Avatar**: %s
🏆 **Achievement**: %s
📝 **Description**: %s
💰 **Reward**: %s NAMO
⭐ **Rarity**: %s
━━━━━━━━━━━━━━━━━━━━━━━━━━━

💬 **%s says**: "%s"

🎬 **Bollywood Moment**: %s

🔥 **Stats Update**:
├─ Level: %d
├─ Total Earnings: %s NAMO
├─ Streak: %d days
└─ Rank: #%d

Keep coding, keep achieving! 🚀`,
		profile.GithubUsername,
		avatar.Name,
		achievement.Name,
		achievement.Description,
		achievement.RewardAmount.Amount,
		achievement.RarityLevel.String(),
		avatar.Character,
		avatar.SignatureQuote,
		quote.Text,
		profile.Level,
		profile.TotalEarnings.Amount,
		profile.CurrentStreakDays,
		profile.LeaderboardRank,
	)

	return content
}

// generateWhatsAppPost creates WhatsApp forward-friendly post
func (g *SocialMediaGenerator) generateWhatsAppPost(profile *DeveloperProfile, achievement *Achievement, quote *HumorQuote) string {
	// Uncle-friendly format with lots of emojis
	content := fmt.Sprintf(`*🌟 DeshChain Achievement Alert 🌟*

देखिये जी! 👏

*%s* ने कमाल कर दिया! 

✅ Achievement: *%s*
💰 Reward: *₹%s*
🎯 Level: *%d*

_"%s"_

👉 *Proud Indian Developer!* 🇮🇳
👉 Making India Digital! 💻
👉 Jai Hind! 🙏

*Share with other developers!*
_Forwarded many times_`,
		profile.GithubUsername,
		achievement.Name,
		formatIndianCurrency(achievement.RewardAmount.Amount),
		profile.Level,
		quote.Text,
	)

	return content
}

// generateInstagramPost creates Instagram-style post
func (g *SocialMediaGenerator) generateInstagramPost(profile *DeveloperProfile, achievement *Achievement, quote *HumorQuote) string {
	avatar := GetAvatarByType(profile.ActiveAvatar)
	
	content := fmt.Sprintf(`🎬 NEW ACHIEVEMENT DROP! 🎬

%s as %s in...
"%s"

%s

🏆 Achievement Unlocked: %s
💰 Box Office Collection: ₹%s
🎭 Character: %s
⭐ Rating: %s/5

Director: @DeshChain
Producer: @%s
Music: GitHub Actions 🎵

Now Streaming on Blockchain! 🚀

Double tap if you're proud of Indian developers! ❤️`,
		profile.GithubUsername,
		avatar.Name,
		GetRandomMoviePosterTitle(achievement.Name),
		quote.Text,
		achievement.Name,
		formatIndianCurrency(achievement.RewardAmount.Amount),
		avatar.Title,
		"5", // Always 5 stars for our developers!
		profile.GithubUsername,
	)

	return content
}

// generateLinkedInPost creates professional yet fun LinkedIn post
func (g *SocialMediaGenerator) generateLinkedInPost(profile *DeveloperProfile, achievement *Achievement, quote *HumorQuote) string {
	content := fmt.Sprintf(`🎯 Proud Achievement Moment at DeshChain! 

I'm thrilled to share that I've unlocked "%s" on DeshChain, India's first cultural blockchain platform.

📊 Achievement Details:
• Challenge: %s
• Reward: %s NAMO tokens
• Impact: Contributing to India's digital transformation

💭 As they say in Bollywood: "%s"

This achievement represents not just personal growth, but also my commitment to building a digitally empowered India through blockchain technology.

Special thanks to the DeshChain community for creating a platform where coding meets culture, and where every contribution matters.

#ProudIndianDeveloper #DeshChain #BlockchainIndia #CodingWithCulture #DigitalIndia`,
		achievement.Name,
		achievement.Description,
		achievement.RewardAmount.Amount,
		quote.Text,
	)

	return content
}

// GenerateStreakPost creates streak-specific posts
func (g *SocialMediaGenerator) GenerateStreakPost(profile *DeveloperProfile, days uint32) string {
	milestones := map[uint32]string{
		7:   "🔥 SAAT DIN SAAT COMMIT! %s ne lagayi aag! Singham Returns daily! 🦁",
		30:  "📿 30 DIN KA TAPASYA! %s = Modern day Coding Rishi! 🙏",
		50:  "🏏 HALF CENTURY! %s scores 50-day streak! Kohli consistency! 🏏",
		100: "💯 CENTURY MUBARAK! %s = Sachin of Coding! Standing ovation! 👏",
		365: "🎊 SAAL BHAR CODING! %s completed 365 days! Thalaiva! 🙌",
		420: "🎯 420 ACHIEVED! %s unlocked Khiladi 420! Legend! 🏆",
	}

	// Check for milestone
	for milestone, template := range milestones {
		if days == milestone {
			return fmt.Sprintf(template, profile.GithubUsername)
		}
	}

	// Generic streak message
	return g.humorEngine.GenerateStreakQuote(int(days))
}

// GenerateEarningsPost creates earning announcement
func (g *SocialMediaGenerator) GenerateEarningsPost(profile *DeveloperProfile, amount sdk.Coin) string {
	return g.humorEngine.GenerateEarningsQuote(
		formatIndianCurrency(amount.Amount),
		profile.GithubUsername,
	)
}

// GenerateLevelUpPost creates level up announcement
func (g *SocialMediaGenerator) GenerateLevelUpPost(profile *DeveloperProfile, newLevel uint32) string {
	levelPosts := map[uint32]string{
		10:  "🎭 %s ab JUNIOR se HERO! Picture abhi baaki hai mere dost! 🎬",
		25:  "⭐ %s promoted to SUPERSTAR coder! Taaliyan! 👏",
		50:  "👑 %s is now SHAHENSHAH of DeshChain! Korona unka, takht unka! 🏰",
		75:  "🔥 %s achieved LEGENDARY status! Naam hai risk, kaam hai disk! 💾",
		100: "🔱 %s achieved THALAIVA status! Mind = Blown! 🤯",
	}

	if post, exists := levelPosts[newLevel]; exists {
		return fmt.Sprintf(post, profile.GithubUsername)
	}

	// Generic level up
	return fmt.Sprintf("🎯 Level %d! %s climbing the coding ladder! Next stop: Bollywood! 🎬",
		newLevel, profile.GithubUsername)
}

// GenerateTeamBattleUpdate creates IPL-style commentary
func (g *SocialMediaGenerator) GenerateTeamBattleUpdate(battle *TeamBattle, event string) string {
	switch event {
	case "start":
		return fmt.Sprintf(`🏏 MATCH START! 🏏
%s vs %s
🎯 Battle: %s
💰 Prize Pool: %s NAMO
May the best code win! Let's go!`,
			battle.Team1Name, battle.Team2Name,
			battle.BattleType, battle.PrizePool.Amount)

	case "score_update":
		lead := ""
		if battle.Team1Score > battle.Team2Score {
			lead = fmt.Sprintf("%s leading by %d!", battle.Team1Name, battle.Team1Score-battle.Team2Score)
		} else if battle.Team2Score > battle.Team1Score {
			lead = fmt.Sprintf("%s leading by %d!", battle.Team2Name, battle.Team2Score-battle.Team1Score)
		} else {
			lead = "Neck and neck competition!"
		}

		return fmt.Sprintf(`🔥 SCORE UPDATE 🔥
%s: %d
%s: %d
%s
Time remaining: Calculate yourself! 😄`,
			battle.Team1Name, battle.Team1Score,
			battle.Team2Name, battle.Team2Score,
			lead)

	case "finish":
		winner := battle.Team1Name
		winnerScore := battle.Team1Score
		if battle.Team2Score > battle.Team1Score {
			winner = battle.Team2Name
			winnerScore = battle.Team2Score
		}

		return fmt.Sprintf(`🏆 MATCH FINISHED! 🏆
WINNER: %s 🎊
Final Score: %d
Prize Won: %s NAMO!

What a match! Incredible performance! 👏`,
			winner, winnerScore, battle.PrizePool.Amount)

	default:
		return "🏏 Epic battle in progress! Stay tuned!"
	}
}

// GenerateDailyChallenge creates movie-themed daily challenges
func (g *SocialMediaGenerator) GenerateDailyChallenge() *DailyChallenge {
	challenges := []struct {
		Type        string
		Description string
		Target      uint64
		Reward      int64
		Quote       string
	}{
		{
			Type:        "bug_hunt",
			Description: "Sholay Bug Hunt: 'Kitne bug the?' Fix them all!",
			Target:      5,
			Reward:      1000,
			Quote:       "Jo darr gaya, samjho marr gaya!",
		},
		{
			Type:        "feature_friday",
			Description: "DDLJ Feature Challenge: Create features that live forever!",
			Target:      3,
			Reward:      2000,
			Quote:       "Bade bade repos mein aisi choti choti features hoti rehti hai",
		},
		{
			Type:        "speed_run",
			Description: "Dhoom Speed Challenge: Optimize like there's no tomorrow!",
			Target:      50, // 50% improvement
			Reward:      1500,
			Quote:       "Speed thrills but kills... bugs!",
		},
		{
			Type:        "doc_day",
			Description: "3 Idiots Documentation: All izz well documented!",
			Target:      10,
			Reward:      800,
			Quote:       "Success ke peeche mat bhago, documentation ka peecha karo",
		},
		{
			Type:        "streak_starter",
			Description: "Dangal Commit Challenge: Commit daily like a wrestler!",
			Target:      7,
			Reward:      500,
			Quote:       "Gold medal code chahiye toh mehnat karni padegi",
		},
	}

	selected := challenges[rand.Intn(len(challenges))]

	return &DailyChallenge{
		ChallengeDate: time.Now(),
		ChallengeType: selected.Type,
		Description:   selected.Description,
		TargetValue:   selected.Target,
		Reward:        sdk.NewCoin("namo", sdk.NewInt(selected.Reward)),
		ThemeQuote:    selected.Quote,
		Participants:  []string{},
		Winners:       []string{},
	}
}

// generatePosterUrl creates movie poster style image URL
func (g *SocialMediaGenerator) generatePosterUrl(profile *DeveloperProfile, achievement *Achievement) string {
	// In production, this would call an image generation service
	// For now, return a template URL
	return fmt.Sprintf("/posters/%s_%d_%d.jpg",
		profile.GithubUsername,
		achievement.AchievementId,
		time.Now().Unix(),
	)
}

// formatIndianCurrency formats amount in Indian style (lakhs, crores)
func formatIndianCurrency(amount sdk.Int) string {
	amountStr := amount.String()
	
	// Convert to float for calculation
	val, _ := sdk.NewDecFromStr(amountStr)
	
	if val.GTE(sdk.NewDec(10000000)) { // >= 1 crore
		crores := val.Quo(sdk.NewDec(10000000))
		return fmt.Sprintf("%.2f Cr", crores.MustFloat64())
	} else if val.GTE(sdk.NewDec(100000)) { // >= 1 lakh
		lakhs := val.Quo(sdk.NewDec(100000))
		return fmt.Sprintf("%.2f L", lakhs.MustFloat64())
	} else if val.GTE(sdk.NewDec(1000)) { // >= 1000
		thousands := val.Quo(sdk.NewDec(1000))
		return fmt.Sprintf("%.1fK", thousands.MustFloat64())
	}
	
	return amountStr
}

// GenerateViralityReport checks post performance
func (g *SocialMediaGenerator) GenerateViralityReport(post *SocialMediaPost) string {
	metrics := post.EngagementMetrics
	
	viralityScore := metrics.Likes*3 + metrics.Shares*5 + metrics.Comments*2 + metrics.Views
	
	status := "🌱 Growing"
	if viralityScore > 10000 {
		status = "🔥 VIRAL!"
		metrics.IsViral = true
	} else if viralityScore > 5000 {
		status = "⚡ Hot!"
	} else if viralityScore > 1000 {
		status = "📈 Trending"
	}
	
	return fmt.Sprintf(`📊 Post Performance Report:
Status: %s
👍 Likes: %d
🔄 Shares: %d
💬 Comments: %d
👁️ Views: %d
📈 Virality Score: %d

%s`,
		status,
		metrics.Likes,
		metrics.Shares,
		metrics.Comments,
		metrics.Views,
		viralityScore,
		getViralityTip(viralityScore),
	)
}

// getViralityTip provides tips to increase virality
func getViralityTip(score uint32) string {
	if score < 100 {
		return "💡 Tip: Add more Bollywood references for instant virality!"
	} else if score < 1000 {
		return "💡 Tip: Tag your team members for more engagement!"
	} else if score < 5000 {
		return "💡 Tip: You're doing great! Add a cricket reference to go viral!"
	}
	return "🎊 Congratulations! You've mastered the art of viral posts!"
}

// GenerateMoviePosterCard creates detailed movie poster data
func GenerateMoviePosterCard(profile *DeveloperProfile, achievement *Achievement) *AchievementCard {
	avatar := GetAvatarByType(profile.ActiveAvatar)
	
	return &AchievementCard{
		CardId:           uint64(time.Now().UnixNano()),
		DeveloperAddress: profile.Address,
		AchievementId:    achievement.AchievementId,
		MovieTitle:       GetRandomMoviePosterTitle(achievement.Name),
		TagLine:          GetRandomTagline(),
		StarringText:     fmt.Sprintf("Starring: @%s as %s", profile.GithubUsername, avatar.Name),
		BoxOfficeText:    fmt.Sprintf("Box Office: ₹%s NAMO Collected", formatIndianCurrency(achievement.RewardAmount.Amount)),
		DirectorText:     "Directed by: GitHub Actions | Produced by: DeshChain Studios",
		ReleaseDate:      time.Now(),
		CriticReviews: []string{
			GetCriticReview(),
			GetCriticReview(),
			GetCriticReview(),
		},
		IsNft:      false, // Can be minted later
		ShareCount: 0,
	}
}