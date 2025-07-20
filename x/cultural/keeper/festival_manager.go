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
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/deshchain/x/cultural/types"
)

// FestivalInfo represents comprehensive festival information
type FestivalInfo struct {
	ID                  string    `json:"id"`
	Name                string    `json:"name"`
	Description         string    `json:"description"`
	Date                time.Time `json:"date"`
	EndDate             *time.Time `json:"end_date,omitempty"`
	Type                string    `json:"type"`
	Region              string    `json:"region"`
	BonusRate           float64   `json:"bonus_rate"`
	IsActive            bool      `json:"is_active"`
	TraditionalGreeting string    `json:"traditional_greeting"`
	CulturalTheme       string    `json:"cultural_theme"`
	Colors              FestivalColors `json:"colors"`
	DaysRemaining       int       `json:"days_remaining"`
	Significance        string    `json:"significance"`
}

type FestivalColors struct {
	Primary   string `json:"primary"`
	Secondary string `json:"secondary"`
	Accent    string `json:"accent"`
}

// IsActiveFestival checks if there is any active festival
func (k Keeper) IsActiveFestival(ctx sdk.Context) bool {
	festivals := k.GetActiveFestivals(ctx)
	return len(festivals) > 0
}

// IsNationalFestival checks if current festival is a national festival
func (k Keeper) IsNationalFestival(ctx sdk.Context) bool {
	festivals := k.GetActiveFestivals(ctx)
	for _, festival := range festivals {
		if festival.Type == "national" {
			return true
		}
	}
	return false
}

// GetCurrentFestival returns the primary active festival
func (k Keeper) GetCurrentFestival(ctx sdk.Context) string {
	festivals := k.GetActiveFestivals(ctx)
	if len(festivals) > 0 {
		return festivals[0].Name
	}
	return ""
}

// GetActiveFestivals returns all currently active festivals
func (k Keeper) GetActiveFestivals(ctx sdk.Context) []FestivalInfo {
	currentTime := ctx.BlockTime()
	allFestivals := k.getAllFestivals()
	
	var activeFestivals []FestivalInfo
	for _, festival := range allFestivals {
		if k.isFestivalActive(festival, currentTime) {
			festival.IsActive = true
			festival.DaysRemaining = k.calculateDaysRemaining(festival, currentTime)
			activeFestivals = append(activeFestivals, festival)
		}
	}
	
	return activeFestivals
}

// GetUpcomingFestivals returns festivals happening in the next 30 days
func (k Keeper) GetUpcomingFestivals(ctx sdk.Context) []FestivalInfo {
	currentTime := ctx.BlockTime()
	allFestivals := k.getAllFestivals()
	
	var upcomingFestivals []FestivalInfo
	for _, festival := range allFestivals {
		daysUntil := int(festival.Date.Sub(currentTime).Hours() / 24)
		if daysUntil > 0 && daysUntil <= 30 {
			festival.IsActive = false
			festival.DaysRemaining = daysUntil
			upcomingFestivals = append(upcomingFestivals, festival)
		}
	}
	
	return upcomingFestivals
}

// GetFestivalByID returns a specific festival by ID
func (k Keeper) GetFestivalByID(ctx sdk.Context, festivalID string) (FestivalInfo, bool) {
	allFestivals := k.getAllFestivals()
	for _, festival := range allFestivals {
		if festival.ID == festivalID {
			currentTime := ctx.BlockTime()
			festival.IsActive = k.isFestivalActive(festival, currentTime)
			festival.DaysRemaining = k.calculateDaysRemaining(festival, currentTime)
			return festival, true
		}
	}
	return FestivalInfo{}, false
}

// GetFestivalBonusRate returns the bonus rate for the current festival
func (k Keeper) GetFestivalBonusRate(ctx sdk.Context) sdk.Dec {
	festivals := k.GetActiveFestivals(ctx)
	if len(festivals) == 0 {
		return sdk.ZeroDec()
	}
	
	// Return highest bonus rate if multiple festivals are active
	maxBonus := 0.0
	for _, festival := range festivals {
		if festival.BonusRate > maxBonus {
			maxBonus = festival.BonusRate
		}
	}
	
	return sdk.MustNewDecFromStr(sdk.NewDecWithPrec(int64(maxBonus*100), 2).String())
}

// GetRandomQuote returns a random quote from a list with cultural context
func (k Keeper) GetRandomQuote(ctx sdk.Context, quotes []string) string {
	if len(quotes) == 0 {
		return ""
	}
	
	// Use block time for deterministic randomness
	seed := ctx.BlockTime().Unix() + ctx.BlockHeight()
	index := seed % int64(len(quotes))
	
	return quotes[index]
}

// IsAppropriateCulturalContent validates if content is culturally appropriate
func (k Keeper) IsAppropriateCulturalContent(ctx sdk.Context, content string) bool {
	// Basic content validation
	if len(content) == 0 {
		return false
	}
	
	// Check for inappropriate words (simplified implementation)
	inappropriateWords := []string{
		"hate", "violence", "discrimination", "offensive"} // Add more as needed
	
	contentLower := strings.ToLower(content)
	for _, word := range inappropriateWords {
		if strings.Contains(contentLower, word) {
			return false
		}
	}
	
	return true
}

// ContainsIndianLanguage checks if text contains Indian language content
func (k Keeper) ContainsIndianLanguage(ctx sdk.Context, text string) bool {
	// Check for Devanagari script (Hindi, Sanskrit, etc.)
	for _, char := range text {
		if char >= 0x0900 && char <= 0x097F {
			return true
		}
	}
	
	// Check for other Indian scripts (simplified)
	indianWords := []string{
		"नमस्ते", "धन्यवाद", "स्वागत", "जय", "हिन्द",
		"வணக்கம்", "நன்றி", "स्वदेश", "भारत", "राष्ट्र",
	}
	
	for _, word := range indianWords {
		if strings.Contains(text, word) {
			return true
		}
	}
	
	return false
}

// IsMourningPeriod checks if current time is during a mourning period
func (k Keeper) IsMourningPeriod(ctx sdk.Context) bool {
	// Check for national mourning days or tragic anniversaries
	currentTime := ctx.BlockTime()
	
	// Gandhi's assassination anniversary
	if currentTime.Month() == time.January && currentTime.Day() == 30 {
		return true
	}
	
	// 26/11 Mumbai attacks anniversary
	if currentTime.Month() == time.November && currentTime.Day() == 26 {
		return true
	}
	
	return false
}

// IsAuspiciousTime checks if current time is considered auspicious
func (k Keeper) IsAuspiciousTime(ctx sdk.Context) bool {
	currentTime := ctx.BlockTime()
	
	// Auspicious times based on Hindu calendar (simplified)
	// Early morning (4-6 AM) and evening (6-8 PM) are generally auspicious
	hour := currentTime.Hour()
	if (hour >= 4 && hour <= 6) || (hour >= 18 && hour <= 20) {
		return true
	}
	
	// Full moon days are auspicious
	// This is a simplified check - real implementation would use lunar calendar
	if currentTime.Day() == 15 {
		return true
	}
	
	return false
}

// getAllFestivals returns the complete list of festivals
func (k Keeper) getAllFestivals() []FestivalInfo {
	return []FestivalInfo{
		{
			ID:          "diwali",
			Name:        "Diwali",
			Description: "Festival of Lights celebrating prosperity and good fortune",
			Date:        k.getDiwaliDate(time.Now().Year()),
			Type:        "national",
			Region:      "all_india",
			BonusRate:   0.15, // 15% bonus
			TraditionalGreeting: "दीपावली की शुभकामनाएं! (Deepavali ki Shubhkamnayein!)",
			CulturalTheme: "lights_prosperity",
			Colors: FestivalColors{
				Primary:   "#FF6B35",
				Secondary: "#F7931E",
				Accent:    "#FFD700",
			},
			Significance: "Celebrates the victory of light over darkness and good over evil",
		},
		{
			ID:          "holi",
			Name:        "Holi",
			Description: "Festival of Colors celebrating spring and new beginnings",
			Date:        k.getHoliDate(time.Now().Year()),
			Type:        "national",
			Region:      "all_india",
			BonusRate:   0.12, // 12% bonus
			TraditionalGreeting: "होली की शुभकामनाएं! (Holi ki Shubhkamnayein!)",
			CulturalTheme: "colors_spring",
			Colors: FestivalColors{
				Primary:   "#E91E63",
				Secondary: "#9C27B0",
				Accent:    "#FF9800",
			},
			Significance: "Celebrates the triumph of good over evil and the arrival of spring",
		},
		{
			ID:          "independence_day",
			Name:        "Independence Day",
			Description: "Celebrating India's independence from British rule",
			Date:        time.Date(time.Now().Year(), time.August, 15, 0, 0, 0, 0, time.UTC),
			Type:        "national",
			Region:      "all_india",
			BonusRate:   0.20, // 20% bonus for national day
			TraditionalGreeting: "स्वतंत्रता दिवस की हार्दिक शुभकामनाएं! (Swatantrata Divas ki Hardik Shubhkamnayein!)",
			CulturalTheme: "patriotism_freedom",
			Colors: FestivalColors{
				Primary:   "#FF9933",
				Secondary: "#FFFFFF",
				Accent:    "#138808",
			},
			Significance: "Commemorates India's freedom struggle and celebrates national sovereignty",
		},
		{
			ID:          "republic_day",
			Name:        "Republic Day",
			Description: "Celebrating the adoption of the Indian Constitution",
			Date:        time.Date(time.Now().Year(), time.January, 26, 0, 0, 0, 0, time.UTC),
			Type:        "national",
			Region:      "all_india",
			BonusRate:   0.20, // 20% bonus for national day
			TraditionalGreeting: "गणतंत्र दिवस की हार्दिक शुभकामनाएं! (Gantantra Divas ki Hardik Shubhkamnayein!)",
			CulturalTheme: "constitution_democracy",
			Colors: FestivalColors{
				Primary:   "#FF9933",
				Secondary: "#FFFFFF", 
				Accent:    "#138808",
			},
			Significance: "Honors the Constitution of India and celebrates democratic values",
		},
		{
			ID:          "dussehra",
			Name:        "Dussehra",
			Description: "Victory of good over evil, celebrating Lord Rama's triumph",
			Date:        k.getDussehraDate(time.Now().Year()),
			Type:        "national",
			Region:      "all_india",
			BonusRate:   0.10, // 10% bonus
			TraditionalGreeting: "दशहरा की शुभकामनाएं! (Dussehra ki Shubhkamnayein!)",
			CulturalTheme: "victory_righteousness",
			Colors: FestivalColors{
				Primary:   "#DC143C",
				Secondary: "#FFD700",
				Accent:    "#FF4500",
			},
			Significance: "Symbolizes the victory of righteousness over evil",
		},
		{
			ID:          "eid_ul_fitr",
			Name:        "Eid ul-Fitr",
			Description: "Festival marking the end of Ramadan",
			Date:        k.getEidDate(time.Now().Year()),
			Type:        "religious",
			Region:      "all_india",
			BonusRate:   0.12, // 12% bonus
			TraditionalGreeting: "ईद मुबारक! (Eid Mubarak!)",
			CulturalTheme: "harmony_brotherhood",
			Colors: FestivalColors{
				Primary:   "#00A86B",
				Secondary: "#FFFFFF",
				Accent:    "#FFD700",
			},
			Significance: "Celebrates the completion of fasting and spiritual reflection",
		},
		{
			ID:          "christmas",
			Name:        "Christmas",
			Description: "Celebrating the birth of Jesus Christ",
			Date:        time.Date(time.Now().Year(), time.December, 25, 0, 0, 0, 0, time.UTC),
			Type:        "religious",
			Region:      "all_india",
			BonusRate:   0.10, // 10% bonus
			TraditionalGreeting: "क्रिसमस की शुभकामनाएं! (Christmas ki Shubhkamnayein!)",
			CulturalTheme: "peace_joy",
			Colors: FestivalColors{
				Primary:   "#DC143C",
				Secondary: "#228B22",
				Accent:    "#FFD700",
			},
			Significance: "Celebrates love, peace, and goodwill among all people",
		},
		{
			ID:          "guru_nanak_jayanti",
			Name:        "Guru Nanak Jayanti",
			Description: "Birthday of Guru Nanak, founder of Sikhism",
			Date:        k.getGuruNanakJayantiDate(time.Now().Year()),
			Type:        "religious",
			Region:      "all_india",
			BonusRate:   0.12, // 12% bonus
			TraditionalGreeting: "गुरु नानक जयंती की शुभकामनाएं! (Guru Nanak Jayanti ki Shubhkamnayein!)",
			CulturalTheme: "wisdom_devotion",
			Colors: FestivalColors{
				Primary:   "#FF9933",
				Secondary: "#FFFFFF",
				Accent:    "#138808",
			},
			Significance: "Honors the teachings of Guru Nanak and Sikh values",
		},
	}
}

// Helper functions to calculate festival dates (simplified implementations)
func (k Keeper) getDiwaliDate(year int) time.Time {
	// Diwali dates vary each year based on lunar calendar
	// This is a simplified approximation
	diwaliDates := map[int]time.Time{
		2024: time.Date(2024, time.November, 1, 0, 0, 0, 0, time.UTC),
		2025: time.Date(2025, time.October, 20, 0, 0, 0, 0, time.UTC),
		2026: time.Date(2026, time.November, 8, 0, 0, 0, 0, time.UTC),
	}
	
	if date, exists := diwaliDates[year]; exists {
		return date
	}
	
	// Default approximation if year not in map
	return time.Date(year, time.October, 25, 0, 0, 0, 0, time.UTC)
}

func (k Keeper) getHoliDate(year int) time.Time {
	// Holi dates vary each year based on lunar calendar
	holiDates := map[int]time.Time{
		2024: time.Date(2024, time.March, 25, 0, 0, 0, 0, time.UTC),
		2025: time.Date(2025, time.March, 14, 0, 0, 0, 0, time.UTC),
		2026: time.Date(2026, time.March, 3, 0, 0, 0, 0, time.UTC),
	}
	
	if date, exists := holiDates[year]; exists {
		return date
	}
	
	return time.Date(year, time.March, 15, 0, 0, 0, 0, time.UTC)
}

func (k Keeper) getDussehraDate(year int) time.Time {
	// Dussehra typically occurs 20 days before Diwali
	diwaliDate := k.getDiwaliDate(year)
	return diwaliDate.AddDate(0, 0, -20)
}

func (k Keeper) getEidDate(year int) time.Time {
	// Eid dates vary based on lunar calendar
	eidDates := map[int]time.Time{
		2024: time.Date(2024, time.April, 10, 0, 0, 0, 0, time.UTC),
		2025: time.Date(2025, time.March, 30, 0, 0, 0, 0, time.UTC),
		2026: time.Date(2026, time.March, 20, 0, 0, 0, 0, time.UTC),
	}
	
	if date, exists := eidDates[year]; exists {
		return date
	}
	
	return time.Date(year, time.April, 15, 0, 0, 0, 0, time.UTC)
}

func (k Keeper) getGuruNanakJayantiDate(year int) time.Time {
	// Guru Nanak Jayanti dates vary based on lunar calendar
	dates := map[int]time.Time{
		2024: time.Date(2024, time.November, 15, 0, 0, 0, 0, time.UTC),
		2025: time.Date(2025, time.November, 5, 0, 0, 0, 0, time.UTC),
		2026: time.Date(2026, time.November, 24, 0, 0, 0, 0, time.UTC),
	}
	
	if date, exists := dates[year]; exists {
		return date
	}
	
	return time.Date(year, time.November, 15, 0, 0, 0, 0, time.UTC)
}

// isFestivalActive checks if a festival is currently active
func (k Keeper) isFestivalActive(festival FestivalInfo, currentTime time.Time) bool {
	// Most festivals are active for their specific day
	if festival.EndDate != nil {
		return currentTime.After(festival.Date) && currentTime.Before(*festival.EndDate)
	}
	
	// For single-day festivals, check if it's the same day
	return currentTime.Year() == festival.Date.Year() &&
		currentTime.Month() == festival.Date.Month() &&
		currentTime.Day() == festival.Date.Day()
}

// calculateDaysRemaining calculates days remaining for a festival
func (k Keeper) calculateDaysRemaining(festival FestivalInfo, currentTime time.Time) int {
	if festival.EndDate != nil {
		return int(festival.EndDate.Sub(currentTime).Hours() / 24)
	}
	
	// Calculate days until festival starts
	daysUntil := int(festival.Date.Sub(currentTime).Hours() / 24)
	if daysUntil < 0 {
		return 0 // Festival is today or past
	}
	
	return daysUntil
}