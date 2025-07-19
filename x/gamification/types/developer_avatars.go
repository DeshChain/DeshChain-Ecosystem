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

// GetAllDeveloperAvatars returns all available developer avatars
func GetAllDeveloperAvatars() []DeveloperAvatar {
	return []DeveloperAvatar{
		{
			AvatarType:        AvatarType_AVATAR_TYPE_BUG_BUSTER_BAHUBALI,
			Name:              "Bug Buster Bahubali",
			Title:             "The Legendary Debugger",
			SignatureQuote:    "Bug ko maarne ke liye commitment chahiye, junoon chahiye!",
			SpecialMove:       "Helicopter Debug - Fixes critical bugs with lightning speed",
			RewardMultiplier:  "2",
			IconUrl:           "/avatars/bug_buster_bahubali.png",
			StyleDescription:  "MS Dhoni Style - Cool under pressure, finishes with style",
		},
		{
			AvatarType:        AvatarType_AVATAR_TYPE_FEATURE_KHAN,
			Name:              "Feature Khan",
			Title:             "The King of Features",
			SignatureQuote:    "Feature nahi banana… developer ban jaana hai",
			SpecialMove:       "DDLJ Deploy - Creates features that users never forget",
			RewardMultiplier:  "3",
			IconUrl:           "/avatars/feature_khan.png",
			StyleDescription:  "Shah Rukh Khan Style - Romantic with code, features that touch hearts",
		},
		{
			AvatarType:        AvatarType_AVATAR_TYPE_DOCUMENTATION_RAJNI,
			Name:              "Documentation Rajni",
			Title:             "The Superstar of Docs",
			SignatureQuote:    "Documentation likhta hun, style mein!",
			SpecialMove:       "Mind It! - Documentation so clear, even juniors understand instantly",
			RewardMultiplier:  "2",
			IconUrl:           "/avatars/documentation_rajni.png",
			StyleDescription:  "Rajinikanth Style - Makes the impossible possible, docs with swag",
		},
		{
			AvatarType:        AvatarType_AVATAR_TYPE_SPEED_SULTAN,
			Name:              "Speed Sultan",
			Title:             "The Performance Champion",
			SignatureQuote:    "Code fast karo, lekin dil se karo",
			SpecialMove:       "Dabangg Optimization - Makes code run 10x faster",
			RewardMultiplier:  "4",
			IconUrl:           "/avatars/speed_sultan.png",
			StyleDescription:  "Salman Khan Style - Raw power, breaks performance records",
		},
		{
			AvatarType:        AvatarType_AVATAR_TYPE_COMMIT_KUMAR,
			Name:              "Commit Kumar",
			Title:             "The Consistency King",
			SignatureQuote:    "Khiladi commit karta hai, quit nahi",
			SpecialMove:       "Khiladi 420 - Maintains 420-day commit streaks",
			RewardMultiplier:  "1.5",
			IconUrl:           "/avatars/commit_kumar.png",
			StyleDescription:  "Akshay Kumar Style - Disciplined, punctual, never misses a day",
		},
	}
}

// GetAvatarByType returns a specific avatar
func GetAvatarByType(avatarType AvatarType) *DeveloperAvatar {
	avatars := GetAllDeveloperAvatars()
	for _, avatar := range avatars {
		if avatar.AvatarType == avatarType {
			return &avatar
		}
	}
	return nil
}

// CalculateAvatarBonus calculates the bonus based on avatar and action
func CalculateAvatarBonus(avatar AvatarType, baseReward sdk.Coin, actionType string) sdk.Coin {
	avatarData := GetAvatarByType(avatar)
	if avatarData == nil {
		return baseReward
	}

	// Parse multiplier
	multiplier := sdk.MustNewDecFromStr(avatarData.RewardMultiplier)

	// Check if action matches avatar specialty
	matchesSpecialty := false
	switch avatar {
	case AvatarType_AVATAR_TYPE_BUG_BUSTER_BAHUBALI:
		matchesSpecialty = actionType == "bug_fix"
	case AvatarType_AVATAR_TYPE_FEATURE_KHAN:
		matchesSpecialty = actionType == "feature"
	case AvatarType_AVATAR_TYPE_DOCUMENTATION_RAJNI:
		matchesSpecialty = actionType == "documentation"
	case AvatarType_AVATAR_TYPE_SPEED_SULTAN:
		matchesSpecialty = actionType == "performance"
	case AvatarType_AVATAR_TYPE_COMMIT_KUMAR:
		matchesSpecialty = actionType == "commit"
	}

	if matchesSpecialty {
		// Apply multiplier
		bonusAmount := sdk.NewDecFromInt(baseReward.Amount).Mul(multiplier).TruncateInt()
		return sdk.NewCoin(baseReward.Denom, bonusAmount)
	}

	return baseReward
}

// GetAvatarSpecialMoveDescription returns detailed special move info
func GetAvatarSpecialMoveDescription(avatar AvatarType) string {
	switch avatar {
	case AvatarType_AVATAR_TYPE_BUG_BUSTER_BAHUBALI:
		return `🚁 HELICOPTER DEBUG:
		- Instantly identifies root cause of critical bugs
		- Fixes P0 issues in record time
		- 95% first-time fix success rate
		- Bonus: Extra rewards for fixing bugs others couldn't`

	case AvatarType_AVATAR_TYPE_FEATURE_KHAN:
		return `💕 DDLJ DEPLOY:
		- Creates features with exceptional user experience
		- High user retention and satisfaction
		- Features become instant hits
		- Bonus: 3x rewards for features with >90% user approval`

	case AvatarType_AVATAR_TYPE_DOCUMENTATION_RAJNI:
		return `🎯 MIND IT! DOCUMENTATION:
		- Crystal clear documentation with examples
		- Interactive tutorials and guides
		- Multi-language support
		- Bonus: 2x rewards when docs reduce support tickets`

	case AvatarType_AVATAR_TYPE_SPEED_SULTAN:
		return `⚡ DABANGG OPTIMIZATION:
		- Reduces load time by 10x
		- Memory optimization specialist
		- Database query expert
		- Bonus: 4x rewards for >50% performance improvement`

	case AvatarType_AVATAR_TYPE_COMMIT_KUMAR:
		return `🔥 KHILADI 420:
		- Never breaks commit streaks
		- Consistent high-quality contributions
		- Master of time management
		- Bonus: Exponential rewards for long streaks`

	default:
		return "Special abilities locked!"
	}
}

// GetAvatarProgressionQuotes returns level-up quotes for each avatar
func GetAvatarProgressionQuotes(avatar AvatarType, level uint32) string {
	quotes := map[AvatarType]map[uint32]string{
		AvatarType_AVATAR_TYPE_BUG_BUSTER_BAHUBALI: {
			10: "Debugging ka junoon badh raha hai!",
			25: "Bug Hunter se Bug Terminator ban gaye!",
			50: "Thala for a reason! Bug fixing machine!",
			100: "MS Dhoni level achieved! Captain Cool of debugging!",
		},
		AvatarType_AVATAR_TYPE_FEATURE_KHAN: {
			10: "Romance with code has begun!",
			25: "Features me hai dum!",
			50: "King Khan of features!",
			100: "Badshah of feature development!",
		},
		AvatarType_AVATAR_TYPE_DOCUMENTATION_RAJNI: {
			10: "Style se documentation shuru!",
			25: "Mind it! Docs getting better!",
			50: "Superstar documentation writer!",
			100: "Thalaiva of technical writing!",
		},
		AvatarType_AVATAR_TYPE_SPEED_SULTAN: {
			10: "Speed ka deewana!",
			25: "Performance ka Sultan!",
			50: "Dabangg optimization expert!",
			100: "Bhaijaan of blazing fast code!",
		},
		AvatarType_AVATAR_TYPE_COMMIT_KUMAR: {
			10: "Khiladi mode ON!",
			25: "Consistency ka Kumar!",
			50: "Khiladi of commitments!",
			100: "Boss of daily contributions!",
		},
	}

	if avatarQuotes, ok := quotes[avatar]; ok {
		// Find the appropriate quote for the level
		for lvl := uint32(100); lvl >= 10; lvl -= 15 {
			if level >= lvl {
				if quote, exists := avatarQuotes[lvl]; exists {
					return quote
				}
			}
		}
	}

	return "Level up! Keep coding!"
}

// GetAvatarChallenges returns specific challenges for each avatar type
func GetAvatarChallenges(avatar AvatarType) []string {
	challenges := map[AvatarType][]string{
		AvatarType_AVATAR_TYPE_BUG_BUSTER_BAHUBALI: {
			"Fix 5 critical bugs in one day",
			"Debug without using console.log",
			"Fix a bug that's been open for >30 days",
			"Help 3 teammates debug their issues",
			"Write comprehensive bug report",
		},
		AvatarType_AVATAR_TYPE_FEATURE_KHAN: {
			"Ship 3 features in a week",
			"Get 95% user approval on a feature",
			"Implement a feature requested by >10 users",
			"Create a feature with <5% bug rate",
			"Feature with internationalization",
		},
		AvatarType_AVATAR_TYPE_DOCUMENTATION_RAJNI: {
			"Write docs for 10 functions",
			"Create interactive tutorial",
			"Document entire module",
			"Add examples to all APIs",
			"Translate docs to Hindi",
		},
		AvatarType_AVATAR_TYPE_SPEED_SULTAN: {
			"Reduce page load by 50%",
			"Optimize database queries",
			"Implement caching strategy",
			"Reduce bundle size by 30%",
			"Profile and fix memory leaks",
		},
		AvatarType_AVATAR_TYPE_COMMIT_KUMAR: {
			"7-day commit streak",
			"30-day commit streak",
			"100-day commit streak",
			"Commit at same time daily",
			"Quality commits for 420 days",
		},
	}

	return challenges[avatar]
}

// GenerateAvatarIntroduction creates a Bollywood-style introduction
func GenerateAvatarIntroduction(profile *DeveloperProfile) string {
	avatar := GetAvatarByType(profile.ActiveAvatar)
	if avatar == nil {
		return "A new developer has joined!"
	}

	templates := []string{
		"🎬 Introducing %s as %s! %s",
		"🌟 %s has chosen the path of %s! %s",
		"🎭 Meet %s - The %s! %s",
		"🏆 %s transforms into %s! %s",
		"💫 %s embraces the power of %s! %s",
	}

	return fmt.Sprintf(
		templates[profile.Level%uint32(len(templates))],
		profile.GithubUsername,
		avatar.Name,
		avatar.SignatureQuote,
	)
}

// GetAvatarRankTitle returns special titles based on avatar and level
func GetAvatarRankTitle(avatar AvatarType, level uint32) string {
	titles := map[AvatarType]map[uint32]string{
		AvatarType_AVATAR_TYPE_BUG_BUSTER_BAHUBALI: {
			1:   "Bug Padawan",
			10:  "Bug Hunter",
			25:  "Bug Slayer",
			50:  "Bug Terminator",
			75:  "Bug Annihilator",
			100: "Thala of Debugging",
		},
		AvatarType_AVATAR_TYPE_FEATURE_KHAN: {
			1:   "Feature Beginner",
			10:  "Feature Builder",
			25:  "Feature Artist",
			50:  "Feature Maestro",
			75:  "Feature Wizard",
			100: "Badshah of Features",
		},
		AvatarType_AVATAR_TYPE_DOCUMENTATION_RAJNI: {
			1:   "Doc Writer",
			10:  "Doc Expert",
			25:  "Doc Master",
			50:  "Doc Superstar",
			75:  "Doc Legend",
			100: "Thalaiva of Docs",
		},
		AvatarType_AVATAR_TYPE_SPEED_SULTAN: {
			1:   "Speed Learner",
			10:  "Speed Demon",
			25:  "Speed Master",
			50:  "Speed Sultan",
			75:  "Speed Emperor",
			100: "Bhaijaan of Performance",
		},
		AvatarType_AVATAR_TYPE_COMMIT_KUMAR: {
			1:   "Commit Rookie",
			10:  "Commit Regular",
			25:  "Commit Expert",
			50:  "Commit Machine",
			75:  "Commit Legend",
			100: "Boss of Commits",
		},
	}

	if avatarTitles, ok := titles[avatar]; ok {
		// Find the appropriate title for the level
		for lvl := uint32(100); lvl >= 1; lvl-- {
			if level >= lvl {
				if title, exists := avatarTitles[lvl]; exists {
					return title
				}
			}
		}
	}

	return "Developer"
}

// GetAvatarBattleCry returns battle cries for team competitions
func GetAvatarBattleCry(avatar AvatarType) string {
	cries := map[AvatarType]string{
		AvatarType_AVATAR_TYPE_BUG_BUSTER_BAHUBALI:   "Bugs ka baap aaya! Helicopter debugging activate! 🚁",
		AvatarType_AVATAR_TYPE_FEATURE_KHAN:          "Dilwale feature le jayenge! Code me romance! 💕",
		AvatarType_AVATAR_TYPE_DOCUMENTATION_RAJNI:   "Mind it! Documentation ka style dekho! 🎯",
		AvatarType_AVATAR_TYPE_SPEED_SULTAN:          "Being fast is being human! Dabangg speed! ⚡",
		AvatarType_AVATAR_TYPE_COMMIT_KUMAR:          "Khiladi ready! 420 commits loading! 🔥",
	}

	return cries[avatar]
}

// CalculateTeamSynergy calculates bonus when certain avatars work together
func CalculateTeamSynergy(avatars []AvatarType) float64 {
	synergy := 1.0

	// Check for power combos
	hasBugBuster := false
	hasFeatureKhan := false
	hasDocRajni := false
	hasSpeedSultan := false
	hasCommitKumar := false

	for _, avatar := range avatars {
		switch avatar {
		case AvatarType_AVATAR_TYPE_BUG_BUSTER_BAHUBALI:
			hasBugBuster = true
		case AvatarType_AVATAR_TYPE_FEATURE_KHAN:
			hasFeatureKhan = true
		case AvatarType_AVATAR_TYPE_DOCUMENTATION_RAJNI:
			hasDocRajni = true
		case AvatarType_AVATAR_TYPE_SPEED_SULTAN:
			hasSpeedSultan = true
		case AvatarType_AVATAR_TYPE_COMMIT_KUMAR:
			hasCommitKumar = true
		}
	}

	// Power combos
	if hasBugBuster && hasFeatureKhan {
		synergy += 0.2 // "Bug-free features" combo
	}
	if hasDocRajni && hasFeatureKhan {
		synergy += 0.15 // "Well-documented features" combo
	}
	if hasSpeedSultan && hasBugBuster {
		synergy += 0.25 // "Fast debugging" combo
	}
	if hasCommitKumar && len(avatars) >= 3 {
		synergy += 0.1 // "Consistent team" bonus
	}
	if len(avatars) == 5 && hasBugBuster && hasFeatureKhan && hasDocRajni && hasSpeedSultan && hasCommitKumar {
		synergy += 0.5 // "Avengers assembled" bonus
	}

	return synergy
}