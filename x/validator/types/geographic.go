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
	"strings"
)

// GeographicIncentives defines the bonus structure for India-based validators
type GeographicIncentives struct {
	IndiaBase        sdk.Dec `json:"india_base" yaml:"india_base"`                 // 10% base bonus for India
	Tier2CityBonus   sdk.Dec `json:"tier2_city_bonus" yaml:"tier2_city_bonus"`     // 5% additional for Tier 2/3 cities
	EmploymentBonus  sdk.Dec `json:"employment_bonus" yaml:"employment_bonus"`     // 3% for hiring 5+ local employees
	GreenEnergyBonus sdk.Dec `json:"green_energy_bonus" yaml:"green_energy_bonus"` // 2% for renewable energy usage
}

// DefaultGeographicIncentives returns the default geographic incentive structure
func DefaultGeographicIncentives() GeographicIncentives {
	return GeographicIncentives{
		IndiaBase:        sdk.NewDecWithPrec(10, 2), // 10%
		Tier2CityBonus:   sdk.NewDecWithPrec(5, 2),  // 5%
		EmploymentBonus:  sdk.NewDecWithPrec(3, 2),  // 3%
		GreenEnergyBonus: sdk.NewDecWithPrec(2, 2),  // 2%
	}
}

// ValidatorLocation represents the geographical and operational details of a validator
type ValidatorLocation struct {
	Country             string `json:"country" yaml:"country"`
	State               string `json:"state" yaml:"state"`
	City                string `json:"city" yaml:"city"`
	CityTier            int    `json:"city_tier" yaml:"city_tier"`                         // 1, 2, or 3
	DataCenterAddress   string `json:"datacenter_address" yaml:"datacenter_address"`
	IPAddress           string `json:"ip_address" yaml:"ip_address"`
	LocalEmployees      int    `json:"local_employees" yaml:"local_employees"`
	RenewableEnergyPerc int    `json:"renewable_energy_percent" yaml:"renewable_energy_percent"`
	Verified            bool   `json:"verified" yaml:"verified"`
	VerificationDate    int64  `json:"verification_date" yaml:"verification_date"`
}

// GeographicMultiplier calculates the total geographic bonus multiplier for a validator
func (gi GeographicIncentives) CalculateMultiplier(location ValidatorLocation) sdk.Dec {
	multiplier := sdk.OneDec() // Start with 1.0 (no bonus)
	
	// Only apply bonuses for India-based validators
	if strings.ToLower(location.Country) != "india" {
		return multiplier
	}
	
	// Base India bonus
	multiplier = multiplier.Add(gi.IndiaBase)
	
	// Tier 2/3 city bonus
	if location.CityTier >= 2 {
		multiplier = multiplier.Add(gi.Tier2CityBonus)
	}
	
	// Local employment bonus
	if location.LocalEmployees >= 5 {
		multiplier = multiplier.Add(gi.EmploymentBonus)
	}
	
	// Green energy bonus
	if location.RenewableEnergyPerc >= 50 {
		multiplier = multiplier.Add(gi.GreenEnergyBonus)
	}
	
	return multiplier
}

// GetMaxPossibleMultiplier returns the maximum possible geographic multiplier
func (gi GeographicIncentives) GetMaxPossibleMultiplier() sdk.Dec {
	return sdk.OneDec().
		Add(gi.IndiaBase).
		Add(gi.Tier2CityBonus).
		Add(gi.EmploymentBonus).
		Add(gi.GreenEnergyBonus)
}

// Validate validates the geographic incentives structure
func (gi GeographicIncentives) Validate() error {
	incentives := []struct {
		name  string
		value sdk.Dec
	}{
		{"india_base", gi.IndiaBase},
		{"tier2_city_bonus", gi.Tier2CityBonus},
		{"employment_bonus", gi.EmploymentBonus},
		{"green_energy_bonus", gi.GreenEnergyBonus},
	}
	
	for _, incentive := range incentives {
		if incentive.value.IsNegative() {
			return fmt.Errorf("%s cannot be negative: %s", incentive.name, incentive.value)
		}
		if incentive.value.GT(sdk.NewDecWithPrec(50, 2)) { // Max 50% for any single incentive
			return fmt.Errorf("%s cannot exceed 50%%: %s", incentive.name, incentive.value.Mul(sdk.NewDec(100)))
		}
	}
	
	// Check total bonus doesn't exceed 100%
	total := gi.IndiaBase.Add(gi.Tier2CityBonus).Add(gi.EmploymentBonus).Add(gi.GreenEnergyBonus)
	if total.GT(sdk.OneDec()) {
		return fmt.Errorf("total geographic incentives cannot exceed 100%%: %s", total.Mul(sdk.NewDec(100)))
	}
	
	return nil
}

// Validate validates the validator location information
func (vl ValidatorLocation) Validate() error {
	if vl.Country == "" {
		return fmt.Errorf("country cannot be empty")
	}
	
	if vl.CityTier < 1 || vl.CityTier > 3 {
		return fmt.Errorf("city tier must be 1, 2, or 3")
	}
	
	if vl.LocalEmployees < 0 {
		return fmt.Errorf("local employees cannot be negative")
	}
	
	if vl.RenewableEnergyPerc < 0 || vl.RenewableEnergyPerc > 100 {
		return fmt.Errorf("renewable energy percentage must be between 0 and 100")
	}
	
	return nil
}

// IsEligibleForIndiaBonus checks if validator is eligible for India-specific bonuses
func (vl ValidatorLocation) IsEligibleForIndiaBonus() bool {
	return strings.ToLower(vl.Country) == "india" && vl.Verified
}

// GetTierDescription returns a human-readable description of the city tier
func (vl ValidatorLocation) GetTierDescription() string {
	switch vl.CityTier {
	case 1:
		return "Tier 1 City (Mumbai, Delhi, Bangalore, etc.)"
	case 2:
		return "Tier 2 City (Pune, Ahmedabad, Hyderabad, etc.)"
	case 3:
		return "Tier 3 City (Smaller cities and towns)"
	default:
		return "Unknown Tier"
	}
}

// GeographicVerification represents the verification process for geographic claims
type GeographicVerification struct {
	ValidatorAddress  string   `json:"validator_address" yaml:"validator_address"`
	LocationClaimed   ValidatorLocation `json:"location_claimed" yaml:"location_claimed"`
	VerificationItems []VerificationItem `json:"verification_items" yaml:"verification_items"`
	Status           string   `json:"status" yaml:"status"` // pending, verified, rejected
	VerifierAddress  string   `json:"verifier_address" yaml:"verifier_address"`
	VerificationDate int64    `json:"verification_date" yaml:"verification_date"`
	ExpiryDate       int64    `json:"expiry_date" yaml:"expiry_date"`
}

// VerificationItem represents a single verification requirement
type VerificationItem struct {
	Type        string `json:"type" yaml:"type"`               // ip_check, physical_audit, utility_bill, etc.
	Description string `json:"description" yaml:"description"`
	Required    bool   `json:"required" yaml:"required"`
	Completed   bool   `json:"completed" yaml:"completed"`
	Evidence    string `json:"evidence" yaml:"evidence"`       // Hash or URL of evidence
	Verifier    string `json:"verifier" yaml:"verifier"`       // Address of verifier
}

// GetRequiredVerificationItems returns the list of required verification items for India bonus
func GetRequiredVerificationItems() []VerificationItem {
	return []VerificationItem{
		{
			Type:        "ip_verification",
			Description: "IP address geolocation verification",
			Required:    true,
			Completed:   false,
		},
		{
			Type:        "utility_bills",
			Description: "Local electricity/internet bills",
			Required:    true,
			Completed:   false,
		},
		{
			Type:        "datacenter_registration",
			Description: "Government registration of data center",
			Required:    true,
			Completed:   false,
		},
		{
			Type:        "physical_audit",
			Description: "Physical verification by local auditor",
			Required:    true,
			Completed:   false,
		},
		{
			Type:        "employment_records",
			Description: "Payroll and tax records for local employees",
			Required:    false, // Only if claiming employment bonus
			Completed:   false,
		},
		{
			Type:        "green_energy_certificate",
			Description: "Renewable energy usage certification",
			Required:    false, // Only if claiming green bonus
			Completed:   false,
		},
	}
}

// VerificationStatus constants
const (
	VerificationStatusPending  = "pending"
	VerificationStatusVerified = "verified"
	VerificationStatusRejected = "rejected"
	VerificationStatusExpired  = "expired"
)

// City tier mappings for major Indian cities
var IndianCityTiers = map[string]int{
	// Tier 1 Cities
	"mumbai":     1,
	"delhi":      1,
	"bangalore":  1,
	"kolkata":    1,
	"chennai":    1,
	"hyderabad":  1,
	"pune":       1,
	"ahmedabad":  1,
	
	// Tier 2 Cities
	"jaipur":     2,
	"lucknow":    2,
	"kanpur":     2,
	"nagpur":     2,
	"indore":     2,
	"thane":      2,
	"bhopal":     2,
	"visakhapatnam": 2,
	"pimpri":     2,
	"patna":      2,
	"vadodara":   2,
	"ghaziabad":  2,
	"ludhiana":   2,
	"agra":       2,
	"nashik":     2,
	"ranchi":     2,
	"faridabad":  2,
	"meerut":     2,
	"rajkot":     2,
	"kalyan":     2,
	"vasai":      2,
	"varanasi":   2,
	"srinagar":   2,
	"aurangabad": 2,
	"dhanbad":    2,
	"amritsar":   2,
	"navi mumbai": 2,
	"allahabad":  2,
	"howrah":     2,
	"gwalior":    2,
	"jabalpur":   2,
	"coimbatore": 2,
	"vijayawada": 2,
	"jodhpur":    2,
	"madurai":    2,
	"raipur":     2,
	"kota":       2,
	"chanddigarh": 2,
	"guwahati":   2,
	"solapur":    2,
	"hubballi":   2,
	"tiruchirappalli": 2,
	"bareilly":   2,
	"mysore":     2,
	"tiruppur":   2,
	"gurgaon":    2,
	"aligarh":    2,
	"jalandhar":  2,
	"bhubaneswar": 2,
	"salem":      2,
	"warangal":   2,
	"guntur":     2,
	"bhiwandi":   2,
	"saharanpur": 2,
	"gorakhpur":  2,
	"bikaner":    2,
	"amravati":   2,
	"noida":      2,
	"jamshedpur": 2,
	"bhilai":     2,
	"cuttack":    2,
	"firozabad":  2,
	"kochi":      2,
	"bhavnagar":  2,
	"dehradun":   2,
	"durgapur":   2,
	"asansol":    2,
	"nanded":     2,
	"kolhapur":   2,
	"ajmer":      2,
	"gulbarga":   2,
	"jamnagar":   2,
	"ujjain":     2,
	"loni":       2,
	"siliguri":   2,
	"jhansi":     2,
	"ulhasnagar": 2,
	"jammu":      2,
	"sangli":     2,
	"mangalore":  2,
	"erode":      2,
	"belgaum":    2,
	"ambattur":   2,
	"tirunelveli": 2,
	"malegaon":   2,
	"gaya":       2,
	"jalgaon":    2,
	"udaipur":    2,
	"maheshtala": 2,
}

// GetCityTier returns the tier of an Indian city
func GetCityTier(city string) int {
	if tier, exists := IndianCityTiers[strings.ToLower(city)]; exists {
		return tier
	}
	// Default to Tier 3 for unlisted cities
	return 3
}

// CalculateAnnualBenefitValue calculates the annual financial benefit of geographic bonuses
func CalculateAnnualBenefitValue(baseEarnings sdk.Dec, multiplier sdk.Dec) sdk.Dec {
	bonus := multiplier.Sub(sdk.OneDec()) // Get the bonus percentage
	return baseEarnings.Mul(bonus)         // Calculate bonus amount
}