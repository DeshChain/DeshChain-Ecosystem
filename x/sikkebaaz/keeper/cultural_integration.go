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

package keeper

import (
	"fmt"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/namo/x/sikkebaaz/types"
)

// ValidateCulturalRequirements validates cultural aspects of a token launch
func (k Keeper) ValidateCulturalRequirements(ctx sdk.Context, launch *types.TokenLaunch) error {
	// Validate PIN code format (Indian postal codes)
	if err := k.validatePincode(launch.CreatorPincode); err != nil {
		return err
	}

	// Validate cultural quote
	if err := k.validateCulturalQuote(ctx, launch.CulturalQuote); err != nil {
		return err
	}

	// Calculate and validate patriotism score
	patriotismScore := k.calculatePatriotismScore(ctx, launch)
	if patriotismScore < types.MinPatriotismScore {
		return types.ErrInsufficientPatriotismScore
	}

	launch.PatriotismScore = patriotismScore
	return nil
}

// validatePincode validates Indian PIN code format and existence
func (k Keeper) validatePincode(pincode string) error {
	// Check length
	if len(pincode) != 6 {
		return types.ErrInvalidPincode
	}

	// Check if all digits
	if _, err := strconv.Atoi(pincode); err != nil {
		return types.ErrInvalidPincode
	}

	// Check valid PIN code ranges (simplified validation)
	firstDigit, _ := strconv.Atoi(string(pincode[0]))
	if firstDigit < 1 || firstDigit > 8 {
		return types.ErrInvalidPincode
	}

	// Additional validation could include checking against actual PIN code database
	// For now, we accept well-formed PIN codes

	return nil
}

// validateCulturalQuote validates cultural quote content and appropriateness
func (k Keeper) validateCulturalQuote(ctx sdk.Context, quote string) error {
	if len(quote) == 0 {
		return types.ErrInvalidCulturalQuote
	}

	if len(quote) > types.MaxCulturalQuoteLength {
		return types.ErrInvalidCulturalQuote
	}

	// Check if quote contains appropriate cultural content
	if !k.culturalKeeper.IsAppropriateCulturalContent(ctx, quote) {
		return types.ErrInvalidCulturalQuote
	}

	return nil
}

// calculatePatriotismScore calculates patriotism score based on various factors
func (k Keeper) calculatePatriotismScore(ctx sdk.Context, launch *types.TokenLaunch) int64 {
	score := int64(0)

	// Base score for launching on Indian blockchain
	score += 10

	// Score for using cultural quote
	if launch.CulturalQuote != "" {
		score += 15
		
		// Bonus for Sanskrit/Hindi content
		if k.culturalKeeper.ContainsIndianLanguage(ctx, launch.CulturalQuote) {
			score += 10
		}
	}

	// Score for token name containing Indian words
	if k.containsIndianWords(launch.TokenName) {
		score += 10
	}

	// Score for meaningful description
	if k.containsIndianWords(launch.TokenDescription) {
		score += 5
	}

	// Bonus during national festivals
	if k.culturalKeeper.IsNationalFestival(ctx) {
		score += 20
	}

	// Score based on creator's history (if available)
	creatorHistory := k.getCreatorCulturalHistory(ctx, launch.Creator)
	score += creatorHistory * 2

	// Regional bonus for underrepresented areas
	regionalBonus := k.getRegionalBonus(ctx, launch.CreatorPincode)
	score += regionalBonus

	return score
}

// containsIndianWords checks if text contains Indian/Sanskrit words
func (k Keeper) containsIndianWords(text string) bool {
	indianWords := []string{
		"desh", "bharat", "india", "swadesh", "rashtra", "lok", "jana", "krishi",
		"gram", "nagar", "seva", "shakti", "yogi", "guru", "dharma", "karma",
		"yoga", "ayurveda", "veda", "mantra", "om", "namaste", "ji", "bhai",
		"didi", "mata", "pita", "jai", "hind", "vandemataram", "satyameva",
	}

	textLower := strings.ToLower(text)
	for _, word := range indianWords {
		if strings.Contains(textLower, word) {
			return true
		}
	}

	return false
}

// getCreatorCulturalHistory gets creator's cultural participation history
func (k Keeper) getCreatorCulturalHistory(ctx sdk.Context, creator string) int64 {
	// Check previous launches by creator
	launches := k.getCreatorLaunches(ctx, creator)
	culturalLaunches := int64(0)

	for _, launchID := range launches {
		launch, found := k.GetTokenLaunch(ctx, launchID)
		if found && launch.PatriotismScore >= types.MinPatriotismScore {
			culturalLaunches++
		}
	}

	return culturalLaunches
}

// getRegionalBonus calculates bonus for underrepresented regions
func (k Keeper) getRegionalBonus(ctx sdk.Context, pincode string) int64 {
	// Map first digit to regions
	regionMap := map[string]string{
		"1": "Delhi/North",
		"2": "Punjab/Haryana", 
		"3": "Rajasthan",
		"4": "Gujarat/Maharashtra",
		"5": "South India",
		"6": "East India",
		"7": "Northeast",
		"8": "Bihar/Eastern UP",
	}

	firstDigit := string(pincode[0])
	region := regionMap[firstDigit]

	// Get launch count per region
	regionLaunches := k.getRegionLaunchCount(ctx, region)
	
	// Bonus for underrepresented regions (inverse relationship)
	if regionLaunches < 10 {
		return 15
	} else if regionLaunches < 50 {
		return 10
	} else if regionLaunches < 100 {
		return 5
	}

	return 0
}

// ApplyFestivalBonuses applies festival-specific bonuses and features
func (k Keeper) ApplyFestivalBonuses(ctx sdk.Context, launch *types.TokenLaunch) {
	if !k.culturalKeeper.IsActiveFestival(ctx) {
		return
	}

	currentFestival := k.culturalKeeper.GetCurrentFestival(ctx)
	bonusRate := sdk.MustNewDecFromStr(types.FestivalBonusRate)

	// Apply different bonuses based on festival type
	switch currentFestival {
	case "Diwali":
		// 15% bonus for prosperity festival
		bonusRate = bonusRate.Add(sdk.MustNewDecFromStr("0.05"))
		launch.CulturalQuote = k.getDiwaliQuote(ctx)
		
	case "Holi":
		// 12% bonus for color festival  
		bonusRate = bonusRate.Add(sdk.MustNewDecFromStr("0.02"))
		launch.CulturalQuote = k.getHoliQuote(ctx)
		
	case "Independence Day", "Republic Day":
		// 20% bonus for national festivals
		bonusRate = bonusRate.Add(sdk.MustNewDecFromStr("0.10"))
		launch.CulturalQuote = k.getNationalQuote(ctx)
		launch.PatriotismScore += 25
		
	case "Dussehra":
		// 10% bonus for victory festival
		launch.CulturalQuote = k.getDussehraQuote(ctx)
		
	case "Diwali", "Eid", "Christmas", "Guru Nanak Jayanti":
		// Multi-religious harmony bonus
		launch.PatriotismScore += 15
	}

	// Apply festival bonus to target amount
	bonusAmount := launch.TargetAmount.ToDec().Mul(bonusRate).TruncateInt()
	launch.TargetAmount = launch.TargetAmount.Add(bonusAmount)
	launch.FestivalBonus = true

	// Create festival bonus record
	festivalBonus := types.FestivalBonus{
		LaunchID:      launch.LaunchID,
		FestivalName:  currentFestival,
		BonusRate:     bonusRate,
		BonusAmount:   bonusAmount,
		AppliedAt:     ctx.BlockTime(),
		CulturalQuote: launch.CulturalQuote,
	}

	k.setFestivalBonus(ctx, festivalBonus)

	k.Logger(ctx).Info("Applied festival bonus",
		"launch_id", launch.LaunchID,
		"festival", currentFestival,
		"bonus_rate", bonusRate,
		"bonus_amount", bonusAmount,
		"new_target", launch.TargetAmount,
	)
}

// GetLocalCommunityLaunches gets launches from the same PIN code area
func (k Keeper) GetLocalCommunityLaunches(ctx sdk.Context, pincode string) []types.TokenLaunch {
	launches := k.getPincodeLaunches(ctx, pincode)
	result := make([]types.TokenLaunch, 0, len(launches))

	for _, launchID := range launches {
		launch, found := k.GetTokenLaunch(ctx, launchID)
		if found {
			result = append(result, launch)
		}
	}

	return result
}

// CalculateLocalCharityAllocation calculates charity for local NGOs
func (k Keeper) CalculateLocalCharityAllocation(ctx sdk.Context, launch *types.TokenLaunch) sdk.Int {
	baseAllocation := types.CalculateCharityAllocation(launch.RaisedAmount)
	
	// Bonus for high patriotism score
	if launch.PatriotismScore >= 80 {
		bonusAllocation := baseAllocation.MulRaw(50).QuoRaw(100) // 50% bonus
		return baseAllocation.Add(bonusAllocation)
	} else if launch.PatriotismScore >= 70 {
		bonusAllocation := baseAllocation.MulRaw(25).QuoRaw(100) // 25% bonus
		return baseAllocation.Add(bonusAllocation)
	}

	return baseAllocation
}

// ValidateSeasonalRestrictions checks seasonal restrictions
func (k Keeper) ValidateSeasonalRestrictions(ctx sdk.Context, launch *types.TokenLaunch) error {
	// No launches during mourning periods
	if k.culturalKeeper.IsMourningPeriod(ctx) {
		return types.ErrFeatureNotEnabled
	}

	// Encourage launches during auspicious times
	if k.culturalKeeper.IsAuspiciousTime(ctx) {
		launch.PatriotismScore += 10
	}

	return nil
}

// Helper functions for festival quotes

func (k Keeper) getDiwaliQuote(ctx sdk.Context) string {
	quotes := []string{
		"दीपावली की शुभकामनाएं! धन और समृद्धि आपके साथ हो।",
		"May the light of Diwali illuminate your path to prosperity",
		"दीपों का त्योहार आपके जीवन में खुशियां लाए",
		"Light up your investments this Diwali!",
	}
	return k.culturalKeeper.GetRandomQuote(ctx, quotes)
}

func (k Keeper) getHoliQuote(ctx sdk.Context) string {
	quotes := []string{
		"होली के रंग आपके जीवन में नई शुरुआत लाएं",
		"Colors of prosperity this Holi!",
		"Celebrate diversity in your portfolio",
		"होली है! रंग बरसे, पैसा भी बरसे!",
	}
	return k.culturalKeeper.GetRandomQuote(ctx, quotes)
}

func (k Keeper) getNationalQuote(ctx sdk.Context) string {
	quotes := []string{
		"जय हिन्द! स्वतंत्रता के साथ आर्थिक आजादी",
		"Jai Hind! Economic freedom through blockchain",
		"वंदे मातरम्! देश के विकास में योगदान दें",
		"सत्यमेव जयते - Truth in every transaction",
		"आत्मनिर्भर भारत की दिशा में एक कदम",
	}
	return k.culturalKeeper.GetRandomQuote(ctx, quotes)
}

func (k Keeper) getDussehraQuote(ctx sdk.Context) string {
	quotes := []string{
		"दशहरे पर बुराई पर अच्छाई की जीत",
		"Victory of good over evil in markets too!",
		"असत्य पर सत्य की विजय",
		"Dussehra - Triumph over market manipulations",
	}
	return k.culturalKeeper.GetRandomQuote(ctx, quotes)
}

// Helper functions

func (k Keeper) getCreatorLaunches(ctx sdk.Context, creator string) []string {
	// This would return launch IDs for the creator
	// Implementation depends on indexing structure
	return []string{} // Simplified
}

func (k Keeper) getRegionLaunchCount(ctx sdk.Context, region string) int64 {
	// Count launches per region
	// Implementation would iterate through launches and count by region
	return 0 // Simplified
}

func (k Keeper) getPincodeLaunches(ctx sdk.Context, pincode string) []string {
	// This would return launch IDs for the pincode
	// Implementation depends on indexing structure
	return []string{} // Simplified
}