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
	"math/rand"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/deshchain/deshchain/x/moneyorder/types"
)

// Festival periods and cultural quotes
var (
	// Festival dates (can be updated via governance)
	festivalPeriods = map[string][]time.Time{
		types.FestivalDiwali: {
			time.Date(2024, time.November, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, time.November, 5, 23, 59, 59, 0, time.UTC),
		},
		types.FestivalHoli: {
			time.Date(2024, time.March, 25, 0, 0, 0, 0, time.UTC),
			time.Date(2024, time.March, 26, 23, 59, 59, 0, time.UTC),
		},
		types.FestivalDussehra: {
			time.Date(2024, time.October, 12, 0, 0, 0, 0, time.UTC),
			time.Date(2024, time.October, 13, 23, 59, 59, 0, time.UTC),
		},
		types.FestivalIndependence: {
			time.Date(2024, time.August, 15, 0, 0, 0, 0, time.UTC),
			time.Date(2024, time.August, 15, 23, 59, 59, 0, time.UTC),
		},
		types.FestivalRepublic: {
			time.Date(2024, time.January, 26, 0, 0, 0, 0, time.UTC),
			time.Date(2024, time.January, 26, 23, 59, 59, 0, time.UTC),
		},
	}

	// Cultural quotes by language
	culturalQuotes = map[string][]string{
		types.LanguageHindi: {
			"वसुधैव कुटुम्बकम् - संपूर्ण विश्व एक परिवार है",
			"अतिथि देवो भव - अतिथि भगवान के समान है",
			"सत्यमेव जयते - सत्य की ही विजय होती है",
			"कर्मण्येवाधिकारस्ते मा फलेषु कदाचन",
			"जननी जन्मभूमिश्च स्वर्गादपि गरीयसी",
		},
		types.LanguageEnglish: {
			"Unity in Diversity - The strength of India",
			"Where the mind is without fear - Tagore",
			"Be the change you wish to see - Gandhi",
			"A nation's culture resides in the hearts of its people",
			"India is the cradle of human race - Mark Twain",
		},
		types.LanguageBengali: {
			"একলা চলো রে - Walk alone if needed",
			"যেখানে দেখিবে ছাই, উড়াইয়া দেখ তাই",
			"মানুষ মানুষের জন্য - Humans for humanity",
			"আমার সোনার বাংলা, আমি তোমায় ভালোবাসি",
			"নিজের চেষ্টায় নিজের ভাগ্য গড়",
		},
		types.LanguageTamil: {
			"யாதும் ஊரே யாவரும் கேளிர் - Every place is home, everyone is kin",
			"கற்றது கைமண் அளவு, கல்லாதது உலகளவு",
			"அன்பே சிவம் - Love is God",
			"பொறுத்தார் பூமி ஆள்வார்",
			"தமிழுக்கும் அமுதென்று பேர்",
		},
	}

	// Festival greetings
	festivalGreetings = map[string]map[string]string{
		types.FestivalDiwali: {
			types.LanguageHindi:    "दीपावली की हार्दिक शुभकामनाएं! 🪔",
			types.LanguageEnglish:  "Happy Diwali! May light triumph over darkness 🪔",
			types.LanguageBengali:  "শুভ দীপাবলি! আলোর উৎসব 🪔",
			types.LanguageTamil:    "தீபாவளி நல்வாழ்த்துக்கள்! 🪔",
			types.LanguageGujarati: "દીપાવલીની હાર્દિક શુભકામના! 🪔",
		},
		types.FestivalHoli: {
			types.LanguageHindi:    "होली की रंगीन शुभकामनाएं! 🎨",
			types.LanguageEnglish:  "Happy Holi! Festival of colors 🎨",
			types.LanguageBengali:  "শুভ হোলি! রঙের উৎসব 🎨",
			types.LanguageTamil:    "ஹோலி வாழ்த்துக்கள்! 🎨",
			types.LanguageGujarati: "હોળીની શુભકામના! 🎨",
		},
	}
)

// IsInFestivalPeriod checks if current time is within any festival period
func (k Keeper) IsInFestivalPeriod(ctx sdk.Context) bool {
	currentTime := ctx.BlockTime()
	
	for _, periods := range festivalPeriods {
		if len(periods) >= 2 && currentTime.After(periods[0]) && currentTime.Before(periods[1]) {
			return true
		}
	}
	
	return false
}

// GetCurrentFestival returns the current festival if any
func (k Keeper) GetCurrentFestival(ctx sdk.Context) (string, bool) {
	currentTime := ctx.BlockTime()
	
	for festival, periods := range festivalPeriods {
		if len(periods) >= 2 && currentTime.After(periods[0]) && currentTime.Before(periods[1]) {
			return festival, true
		}
	}
	
	return "", false
}

// GetCulturalQuote returns a random cultural quote in the specified language
func (k Keeper) GetCulturalQuote(ctx sdk.Context, language string) string {
	quotes, exists := culturalQuotes[language]
	if !exists {
		// Default to English if language not found
		quotes = culturalQuotes[types.LanguageEnglish]
	}
	
	// Use block height as seed for deterministic randomness
	rand.Seed(ctx.BlockHeight())
	return quotes[rand.Intn(len(quotes))]
}

// GetFestivalGreeting returns festival greeting in specified language
func (k Keeper) GetFestivalGreeting(ctx sdk.Context, festival, language string) string {
	greetings, exists := festivalGreetings[festival]
	if !exists {
		return ""
	}
	
	greeting, exists := greetings[language]
	if !exists {
		// Default to English if language not found
		greeting = greetings[types.LanguageEnglish]
	}
	
	return greeting
}

// ApplyFestivalBonus applies festival bonus to fees if applicable
func (k Keeper) ApplyFestivalBonus(ctx sdk.Context, baseFee sdk.Dec) sdk.Dec {
	params := k.GetParams(ctx)
	
	if !params.EnableFestivalBonuses {
		return baseFee
	}
	
	if k.IsInFestivalPeriod(ctx) {
		// Apply festival discount
		discount := baseFee.Mul(params.FestivalDiscount)
		return baseFee.Sub(discount)
	}
	
	return baseFee
}

// ApplyCulturalDiscount applies cultural discount for cultural token pairs
func (k Keeper) ApplyCulturalDiscount(ctx sdk.Context, baseFee sdk.Dec, token0, token1 string) sdk.Dec {
	params := k.GetParams(ctx)
	
	if !params.EnableCulturalFeatures {
		return baseFee
	}
	
	// Check if it's a cultural pair
	culturalTokens := map[string]bool{
		"unamo":     true,
		"uinr":      true,
		"ucultural": true,
		"uheritage": true,
		"utemple":   true,
		"ufestival": true,
	}
	
	if culturalTokens[token0] || culturalTokens[token1] {
		// Apply cultural discount
		discount := baseFee.Mul(params.CulturalTokenDiscount)
		return baseFee.Sub(discount)
	}
	
	return baseFee
}

// ApplyVillagePoolDiscount applies discount for village pool members
func (k Keeper) ApplyVillagePoolDiscount(ctx sdk.Context, baseFee sdk.Dec, userAddr sdk.AccAddress) sdk.Dec {
	params := k.GetParams(ctx)
	
	if !params.EnableVillagePools {
		return baseFee
	}
	
	// Check if user is a member of any active village pool
	isVillageMember := false
	k.IterateVillagePools(ctx, func(pool types.VillagePool) bool {
		if pool.Active && pool.Verified {
			members := k.GetVillagePoolMembers(ctx, pool.PoolId)
			for _, member := range members {
				memberAddr, _ := sdk.AccAddressFromBech32(member.MemberAddress)
				if memberAddr.Equals(userAddr) {
					isVillageMember = true
					return true // stop iteration
				}
			}
		}
		return false
	})
	
	if isVillageMember {
		// Apply village pool discount
		discount := baseFee.Mul(params.VillagePoolDiscount)
		return baseFee.Sub(discount)
	}
	
	return baseFee
}

// ApplySeniorCitizenDiscount applies discount for senior citizens (if KYC verified)
func (k Keeper) ApplySeniorCitizenDiscount(ctx sdk.Context, baseFee sdk.Dec, userAddr sdk.AccAddress) sdk.Dec {
	params := k.GetParams(ctx)
	
	// In a real implementation, this would check KYC data for age verification
	// For now, we'll return the base fee
	// TODO: Integrate with KYC module to verify senior citizen status
	
	return baseFee
}

// GetBestDiscount returns the best applicable discount for a transaction
func (k Keeper) GetBestDiscount(ctx sdk.Context, baseFee sdk.Dec, userAddr sdk.AccAddress, token0, token1 string) sdk.Dec {
	// Calculate all applicable discounts
	festivalFee := k.ApplyFestivalBonus(ctx, baseFee)
	culturalFee := k.ApplyCulturalDiscount(ctx, baseFee, token0, token1)
	villageFee := k.ApplyVillagePoolDiscount(ctx, baseFee, userAddr)
	seniorFee := k.ApplySeniorCitizenDiscount(ctx, baseFee, userAddr)
	
	// Return the lowest fee (best discount)
	minFee := baseFee
	if festivalFee.LT(minFee) {
		minFee = festivalFee
	}
	if culturalFee.LT(minFee) {
		minFee = culturalFee
	}
	if villageFee.LT(minFee) {
		minFee = villageFee
	}
	if seniorFee.LT(minFee) {
		minFee = seniorFee
	}
	
	return minFee
}

// EmitCulturalEvent emits a cultural event with appropriate attributes
func (k Keeper) EmitCulturalEvent(ctx sdk.Context, eventType string, sender sdk.AccAddress, attributes ...sdk.Attribute) {
	// Add cultural context to events
	festival, hasFestival := k.GetCurrentFestival(ctx)
	language := types.LanguageEnglish // Default, would be determined by user preference
	
	culturalAttrs := []sdk.Attribute{
		sdk.NewAttribute(types.AttributeKeyLanguage, language),
		sdk.NewAttribute(types.AttributeKeyCulturalQuote, k.GetCulturalQuote(ctx, language)),
	}
	
	if hasFestival {
		culturalAttrs = append(culturalAttrs,
			sdk.NewAttribute(types.AttributeKeyFestival, festival),
			sdk.NewAttribute(types.AttributeKeyFestivalGreeting, k.GetFestivalGreeting(ctx, festival, language)),
		)
	}
	
	// Combine all attributes
	allAttrs := append(culturalAttrs, attributes...)
	allAttrs = append(allAttrs, sdk.NewAttribute(types.AttributeKeySender, sender.String()))
	
	// Emit the event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(eventType, allAttrs...),
	)
}

// UpdateFestivalPeriods updates festival periods (called by governance)
func (k Keeper) UpdateFestivalPeriods(ctx sdk.Context, festival string, startTime, endTime time.Time) error {
	// Validate festival
	validFestivals := map[string]bool{
		types.FestivalDiwali:       true,
		types.FestivalHoli:         true,
		types.FestivalDussehra:     true,
		types.FestivalIndependence: true,
		types.FestivalRepublic:     true,
	}
	
	if !validFestivals[festival] {
		return fmt.Errorf("invalid festival: %s", festival)
	}
	
	// Update festival periods
	festivalPeriods[festival] = []time.Time{startTime, endTime}
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeFestivalPeriodUpdated,
			sdk.NewAttribute(types.AttributeKeyFestival, festival),
			sdk.NewAttribute("start_time", startTime.String()),
			sdk.NewAttribute("end_time", endTime.String()),
		),
	)
	
	return nil
}