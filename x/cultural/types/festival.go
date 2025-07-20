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
	"time"
)

// Festival represents a cultural or religious festival
type Festival struct {
	Id                  uint64            `json:"id" yaml:"id"`
	Name                string            `json:"name" yaml:"name"`
	Description         string            `json:"description" yaml:"description"`
	Date                time.Time         `json:"date" yaml:"date"`
	EndDate             *time.Time        `json:"end_date,omitempty" yaml:"end_date,omitempty"`
	Type                string            `json:"type" yaml:"type"`                     // national, religious, regional, seasonal
	Region              string            `json:"region" yaml:"region"`                 // all_india, north, south, east, west, specific_state
	BonusRate           float64           `json:"bonus_rate" yaml:"bonus_rate"`         // Bonus percentage (0.10 = 10%)
	TraditionalGreeting string            `json:"traditional_greeting" yaml:"traditional_greeting"`
	CulturalTheme       string            `json:"cultural_theme" yaml:"cultural_theme"`
	Colors              FestivalColors    `json:"colors" yaml:"colors"`
	Significance        string            `json:"significance" yaml:"significance"`
	Tags                []string          `json:"tags" yaml:"tags"`
	CreatedAt           int64             `json:"created_at" yaml:"created_at"`
	UpdatedAt           int64             `json:"updated_at" yaml:"updated_at"`
	IsVerified          bool              `json:"is_verified" yaml:"is_verified"`
	UsageCount          uint64            `json:"usage_count" yaml:"usage_count"`
	Metadata            map[string]string `json:"metadata" yaml:"metadata"`
}

// FestivalColors represents the color scheme for a festival
type FestivalColors struct {
	Primary   string `json:"primary" yaml:"primary"`     // Main festival color
	Secondary string `json:"secondary" yaml:"secondary"` // Secondary accent color
	Accent    string `json:"accent" yaml:"accent"`       // Highlight/accent color
}

// FestivalBonus represents bonus applied during festivals
type FestivalBonus struct {
	FestivalId    string    `json:"festival_id" yaml:"festival_id"`
	UserId        string    `json:"user_id" yaml:"user_id"`
	TransactionId string    `json:"transaction_id" yaml:"transaction_id"`
	BonusAmount   string    `json:"bonus_amount" yaml:"bonus_amount"`
	BonusRate     float64   `json:"bonus_rate" yaml:"bonus_rate"`
	AppliedAt     time.Time `json:"applied_at" yaml:"applied_at"`
	IsActive      bool      `json:"is_active" yaml:"is_active"`
}

// FestivalStatus represents the current status of festivals
type FestivalStatus struct {
	ActiveFestivals   []string          `json:"active_festivals" yaml:"active_festivals"`
	UpcomingFestivals []string          `json:"upcoming_festivals" yaml:"upcoming_festivals"`
	CurrentBonusRate  float64           `json:"current_bonus_rate" yaml:"current_bonus_rate"`
	NextFestival      *FestivalInfo     `json:"next_festival,omitempty" yaml:"next_festival,omitempty"`
	LastUpdated       time.Time         `json:"last_updated" yaml:"last_updated"`
	Metadata          map[string]string `json:"metadata" yaml:"metadata"`
}

// FestivalInfo represents comprehensive festival information for API responses
type FestivalInfo struct {
	Id                  string         `json:"id" yaml:"id"`
	Name                string         `json:"name" yaml:"name"`
	Description         string         `json:"description" yaml:"description"`
	Date                time.Time      `json:"date" yaml:"date"`
	EndDate             *time.Time     `json:"end_date,omitempty" yaml:"end_date,omitempty"`
	Type                string         `json:"type" yaml:"type"`
	Region              string         `json:"region" yaml:"region"`
	BonusRate           float64        `json:"bonus_rate" yaml:"bonus_rate"`
	IsActive            bool           `json:"is_active" yaml:"is_active"`
	TraditionalGreeting string         `json:"traditional_greeting" yaml:"traditional_greeting"`
	CulturalTheme       string         `json:"cultural_theme" yaml:"cultural_theme"`
	Colors              FestivalColors `json:"colors" yaml:"colors"`
	DaysRemaining       int            `json:"days_remaining" yaml:"days_remaining"`
	Significance        string         `json:"significance" yaml:"significance"`
	Tags                []string       `json:"tags" yaml:"tags"`
	IsVerified          bool           `json:"is_verified" yaml:"is_verified"`
}

// FestivalPreferences represents user preferences for festival notifications
type FestivalPreferences struct {
	UserId                string   `json:"user_id" yaml:"user_id"`
	EnableNotifications   bool     `json:"enable_notifications" yaml:"enable_notifications"`
	PreferredFestivals    []string `json:"preferred_festivals" yaml:"preferred_festivals"`
	PreferredRegions      []string `json:"preferred_regions" yaml:"preferred_regions"`
	EnableBonuses         bool     `json:"enable_bonuses" yaml:"enable_bonuses"`
	LanguagePreference    string   `json:"language_preference" yaml:"language_preference"`
	NotificationThreshold int      `json:"notification_threshold" yaml:"notification_threshold"` // Days before festival
	CreatedAt             int64    `json:"created_at" yaml:"created_at"`
	UpdatedAt             int64    `json:"updated_at" yaml:"updated_at"`
}

// Festival type constants
const (
	FestivalTypeNational  = "national"
	FestivalTypeReligious = "religious"
	FestivalTypeRegional  = "regional"
	FestivalTypeSeasonal  = "seasonal"
	FestivalTypeCultural  = "cultural"
)

// Festival region constants
const (
	FestivalRegionAllIndia = "all_india"
	FestivalRegionNorth    = "north"
	FestivalRegionSouth    = "south"
	FestivalRegionEast     = "east"
	FestivalRegionWest     = "west"
	FestivalRegionCentral  = "central"
	FestivalRegionNortheast = "northeast"
)

// Cultural theme constants
const (
	ThemeLightsProsperity      = "lights_prosperity"
	ThemeColorsSpring          = "colors_spring"
	ThemePatriotismFreedom     = "patriotism_freedom"
	ThemeConstitutionDemocracy = "constitution_democracy"
	ThemeVictoryRighteousness  = "victory_righteousness"
	ThemeHarmonyBrotherhood    = "harmony_brotherhood"
	ThemePeaceJoy              = "peace_joy"
	ThemeWisdomDevotion        = "wisdom_devotion"
	ThemeHarvestGratitude      = "harvest_gratitude"
	ThemeKnowledgeWisdom       = "knowledge_wisdom"
)

// ValidateFestival validates festival data
func ValidateFestival(festival Festival) error {
	if festival.Name == "" {
		return ErrInvalidData
	}
	
	if festival.BonusRate < 0 || festival.BonusRate > 1 {
		return ErrInvalidData
	}
	
	if festival.Type == "" {
		return ErrInvalidData
	}
	
	if festival.Region == "" {
		return ErrInvalidData
	}
	
	return nil
}

// IsNationalFestival checks if festival is national level
func (f Festival) IsNationalFestival() bool {
	return f.Type == FestivalTypeNational
}

// IsReligiousFestival checks if festival is religious
func (f Festival) IsReligiousFestival() bool {
	return f.Type == FestivalTypeReligious
}

// IsRegionalFestival checks if festival is regional
func (f Festival) IsRegionalFestival() bool {
	return f.Type == FestivalTypeRegional
}

// GetBonusRateDecimal returns bonus rate as decimal percentage
func (f Festival) GetBonusRateDecimal() float64 {
	return f.BonusRate * 100
}

// IsActiveOn checks if festival is active on given date
func (f Festival) IsActiveOn(date time.Time) bool {
	if f.EndDate != nil {
		return date.After(f.Date) && date.Before(*f.EndDate)
	}
	
	// Single day festival
	return date.Year() == f.Date.Year() &&
		date.Month() == f.Date.Month() &&
		date.Day() == f.Date.Day()
}

// DaysUntil calculates days until festival starts
func (f Festival) DaysUntil(currentDate time.Time) int {
	duration := f.Date.Sub(currentDate)
	days := int(duration.Hours() / 24)
	
	if days < 0 {
		return 0
	}
	
	return days
}

// DaysRemaining calculates days remaining for festival
func (f Festival) DaysRemaining(currentDate time.Time) int {
	if f.EndDate != nil {
		duration := f.EndDate.Sub(currentDate)
		days := int(duration.Hours() / 24)
		return days
	}
	
	// For single-day festivals, check if it's today
	if f.IsActiveOn(currentDate) {
		return 0
	}
	
	return f.DaysUntil(currentDate)
}

// ToFestivalInfo converts Festival to FestivalInfo with runtime data
func (f Festival) ToFestivalInfo(currentDate time.Time) FestivalInfo {
	return FestivalInfo{
		Id:                  string(rune(f.Id)), // Convert uint64 to string
		Name:                f.Name,
		Description:         f.Description,
		Date:                f.Date,
		EndDate:             f.EndDate,
		Type:                f.Type,
		Region:              f.Region,
		BonusRate:           f.BonusRate,
		IsActive:            f.IsActiveOn(currentDate),
		TraditionalGreeting: f.TraditionalGreeting,
		CulturalTheme:       f.CulturalTheme,
		Colors:              f.Colors,
		DaysRemaining:       f.DaysRemaining(currentDate),
		Significance:        f.Significance,
		Tags:                f.Tags,
		IsVerified:          f.IsVerified,
	}
}

// GetFestivalIcon returns appropriate icon for festival theme
func (f Festival) GetFestivalIcon() string {
	switch f.CulturalTheme {
	case ThemeLightsProsperity:
		return "celebration" // Diwali - lights and celebration
	case ThemeColorsSpring:
		return "palette" // Holi - colors and art
	case ThemePatriotismFreedom, ThemeConstitutionDemocracy:
		return "flag" // National festivals - patriotic
	case ThemeVictoryRighteousness:
		return "military_tech" // Dussehra - victory and valor
	case ThemeHarmonyBrotherhood:
		return "diversity_1" // Unity festivals - harmony
	case ThemePeaceJoy:
		return "favorite" // Christmas, peaceful festivals
	case ThemeWisdomDevotion:
		return "auto_stories" // Guru festivals - wisdom
	case ThemeHarvestGratitude:
		return "agriculture" // Harvest festivals
	case ThemeKnowledgeWisdom:
		return "school" // Educational festivals
	default:
		return "celebration" // Default celebration icon
	}
}