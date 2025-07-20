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
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// Parameter keys
var (
	KeyEnableQuotes            = []byte("EnableQuotes")
	KeyQuoteSelectionAlgorithm = []byte("QuoteSelectionAlgorithm")
	KeyAmountRanges            = []byte("AmountRanges")
	KeySeasonalQuotesEnabled   = []byte("SeasonalQuotesEnabled")
	KeyUserPreferencesEnabled  = []byte("UserPreferencesEnabled")
	KeyMinQuoteLength          = []byte("MinQuoteLength")
	KeyMaxQuoteLength          = []byte("MaxQuoteLength")
	KeyDefaultLanguage         = []byte("DefaultLanguage")
	KeyAvailableLanguages      = []byte("AvailableLanguages")
	KeyQuoteRefreshInterval    = []byte("QuoteRefreshInterval")
	KeyMaxQuotesPerUser        = []byte("MaxQuotesPerUser")
	KeyEnableNFTCreation       = []byte("EnableNFTCreation")
	KeyNFTMintingFee           = []byte("NFTMintingFee")
	KeyNFTRoyaltyPercentage    = []byte("NFTRoyaltyPercentage")
	KeyEnableStatistics        = []byte("EnableStatistics")
	KeyStatisticsRetentionDays = []byte("StatisticsRetentionDays")
)

// Default parameter values
const (
	DefaultEnableQuotes            = true
	DefaultQuoteSelectionAlgorithm = SelectionAlgorithmAmount
	DefaultSeasonalQuotesEnabled   = true
	DefaultUserPreferencesEnabled  = true
	DefaultMinQuoteLength          = uint32(10)
	DefaultMaxQuoteLength          = uint32(500)
	DefaultDefaultLanguage         = LanguageEnglish
	DefaultQuoteRefreshInterval    = int64(86400) // 24 hours
	DefaultMaxQuotesPerUser        = uint32(100)
	DefaultEnableNFTCreation       = true
	DefaultNFTRoyaltyPercentage    = "2.5"
	DefaultEnableStatistics        = true
	DefaultStatisticsRetentionDays = uint32(90)
)

// DefaultAmountRanges returns default amount ranges for quote selection
func DefaultAmountRanges() []AmountRange {
	return []AmountRange{
		{
			MinAmount:  "0",
			MaxAmount:  "1000",
			Categories: []string{CategoryMotivation, CategoryWisdom},
		},
		{
			MinAmount:  "1001",
			MaxAmount:  "10000",
			Categories: []string{CategoryPhilosophy, CategoryEducation, CategoryLeadership},
		},
		{
			MinAmount:  "10001",
			MaxAmount:  "100000",
			Categories: []string{CategoryPatriotism, CategoryLeadership, CategorySpirituality},
		},
		{
			MinAmount:  "100001",
			MaxAmount:  "999999999999999999",
			Categories: []string{CategoryPatriotism, CategoryWisdom, CategoryTruth},
		},
	}
}

// DefaultAvailableLanguages returns default available languages
func DefaultAvailableLanguages() []string {
	return []string{
		LanguageEnglish,
		LanguageHindi,
		LanguageSanskrit,
		LanguageTamil,
		LanguageTelugu,
		LanguageBengali,
		LanguageMarathi,
		LanguageGujarati,
		LanguagePunjabi,
		LanguageKannada,
		LanguageMalayalam,
		LanguageOdia,
		LanguageAssamese,
		LanguageUrdu,
	}
}

// DefaultNFTMintingFee returns default NFT minting fee
func DefaultNFTMintingFee() sdk.Coin {
	return sdk.NewCoin("namo", sdk.NewInt(100000000)) // 100 NAMO
}

// ParamKeyTable returns the key table for cultural module parameters
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(
	enableQuotes bool,
	quoteSelectionAlgorithm string,
	amountRanges []AmountRange,
	seasonalQuotesEnabled bool,
	userPreferencesEnabled bool,
	minQuoteLength uint32,
	maxQuoteLength uint32,
	defaultLanguage string,
	availableLanguages []string,
	quoteRefreshInterval int64,
	maxQuotesPerUser uint32,
	enableNFTCreation bool,
	nftMintingFee sdk.Coin,
	nftRoyaltyPercentage string,
	enableStatistics bool,
	statisticsRetentionDays uint32,
) Params {
	return Params{
		EnableQuotes:            enableQuotes,
		QuoteSelectionAlgorithm: quoteSelectionAlgorithm,
		AmountRanges:            amountRanges,
		SeasonalQuotesEnabled:   seasonalQuotesEnabled,
		UserPreferencesEnabled:  userPreferencesEnabled,
		MinQuoteLength:          minQuoteLength,
		MaxQuoteLength:          maxQuoteLength,
		DefaultLanguage:         defaultLanguage,
		AvailableLanguages:      availableLanguages,
		QuoteRefreshInterval:    quoteRefreshInterval,
		MaxQuotesPerUser:        maxQuotesPerUser,
		EnableNftCreation:       enableNFTCreation,
		NftMintingFee:           nftMintingFee,
		NftRoyaltyPercentage:    nftRoyaltyPercentage,
		EnableStatistics:        enableStatistics,
		StatisticsRetentionDays: statisticsRetentionDays,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultEnableQuotes,
		DefaultQuoteSelectionAlgorithm,
		DefaultAmountRanges(),
		DefaultSeasonalQuotesEnabled,
		DefaultUserPreferencesEnabled,
		DefaultMinQuoteLength,
		DefaultMaxQuoteLength,
		DefaultDefaultLanguage,
		DefaultAvailableLanguages(),
		DefaultQuoteRefreshInterval,
		DefaultMaxQuotesPerUser,
		DefaultEnableNFTCreation,
		DefaultNFTMintingFee(),
		DefaultNFTRoyaltyPercentage,
		DefaultEnableStatistics,
		DefaultStatisticsRetentionDays,
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyEnableQuotes, &p.EnableQuotes, validateEnableQuotes),
		paramtypes.NewParamSetPair(KeyQuoteSelectionAlgorithm, &p.QuoteSelectionAlgorithm, validateQuoteSelectionAlgorithm),
		paramtypes.NewParamSetPair(KeyAmountRanges, &p.AmountRanges, validateAmountRanges),
		paramtypes.NewParamSetPair(KeySeasonalQuotesEnabled, &p.SeasonalQuotesEnabled, validateSeasonalQuotesEnabled),
		paramtypes.NewParamSetPair(KeyUserPreferencesEnabled, &p.UserPreferencesEnabled, validateUserPreferencesEnabled),
		paramtypes.NewParamSetPair(KeyMinQuoteLength, &p.MinQuoteLength, validateMinQuoteLength),
		paramtypes.NewParamSetPair(KeyMaxQuoteLength, &p.MaxQuoteLength, validateMaxQuoteLength),
		paramtypes.NewParamSetPair(KeyDefaultLanguage, &p.DefaultLanguage, validateDefaultLanguage),
		paramtypes.NewParamSetPair(KeyAvailableLanguages, &p.AvailableLanguages, validateAvailableLanguages),
		paramtypes.NewParamSetPair(KeyQuoteRefreshInterval, &p.QuoteRefreshInterval, validateQuoteRefreshInterval),
		paramtypes.NewParamSetPair(KeyMaxQuotesPerUser, &p.MaxQuotesPerUser, validateMaxQuotesPerUser),
		paramtypes.NewParamSetPair(KeyEnableNFTCreation, &p.EnableNftCreation, validateEnableNFTCreation),
		paramtypes.NewParamSetPair(KeyNFTMintingFee, &p.NftMintingFee, validateNFTMintingFee),
		paramtypes.NewParamSetPair(KeyNFTRoyaltyPercentage, &p.NftRoyaltyPercentage, validateNFTRoyaltyPercentage),
		paramtypes.NewParamSetPair(KeyEnableStatistics, &p.EnableStatistics, validateEnableStatistics),
		paramtypes.NewParamSetPair(KeyStatisticsRetentionDays, &p.StatisticsRetentionDays, validateStatisticsRetentionDays),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateEnableQuotes(p.EnableQuotes); err != nil {
		return err
	}
	if err := validateQuoteSelectionAlgorithm(p.QuoteSelectionAlgorithm); err != nil {
		return err
	}
	if err := validateAmountRanges(p.AmountRanges); err != nil {
		return err
	}
	if err := validateMinQuoteLength(p.MinQuoteLength); err != nil {
		return err
	}
	if err := validateMaxQuoteLength(p.MaxQuoteLength); err != nil {
		return err
	}
	if err := validateDefaultLanguage(p.DefaultLanguage); err != nil {
		return err
	}
	if err := validateAvailableLanguages(p.AvailableLanguages); err != nil {
		return err
	}
	if err := validateQuoteRefreshInterval(p.QuoteRefreshInterval); err != nil {
		return err
	}
	if err := validateMaxQuotesPerUser(p.MaxQuotesPerUser); err != nil {
		return err
	}
	if err := validateNFTMintingFee(p.NftMintingFee); err != nil {
		return err
	}
	if err := validateNFTRoyaltyPercentage(p.NftRoyaltyPercentage); err != nil {
		return err
	}
	if err := validateStatisticsRetentionDays(p.StatisticsRetentionDays); err != nil {
		return err
	}
	
	// Validate min < max quote length
	if p.MinQuoteLength > p.MaxQuoteLength {
		return fmt.Errorf("min quote length cannot be greater than max quote length")
	}
	
	return nil
}

// String implements the Stringer interface
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// Validation functions

func validateEnableQuotes(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateQuoteSelectionAlgorithm(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	switch v {
	case SelectionAlgorithmRandom,
		SelectionAlgorithmAmount,
		SelectionAlgorithmSeasonal,
		SelectionAlgorithmPersonalized:
		return nil
	default:
		return fmt.Errorf("invalid quote selection algorithm: %s", v)
	}
}

func validateAmountRanges(i interface{}) error {
	ranges, ok := i.([]AmountRange)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	for _, r := range ranges {
		minAmount, err := sdk.NewIntFromString(r.MinAmount)
		if err != nil {
			return fmt.Errorf("invalid min amount: %s", r.MinAmount)
		}
		maxAmount, err := sdk.NewIntFromString(r.MaxAmount)
		if err != nil {
			return fmt.Errorf("invalid max amount: %s", r.MaxAmount)
		}
		if minAmount.GT(maxAmount) {
			return fmt.Errorf("min amount cannot be greater than max amount")
		}
		if len(r.Categories) == 0 {
			return fmt.Errorf("amount range must have at least one category")
		}
	}
	
	return nil
}

func validateSeasonalQuotesEnabled(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateUserPreferencesEnabled(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateMinQuoteLength(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v == 0 {
		return fmt.Errorf("min quote length must be greater than 0")
	}
	return nil
}

func validateMaxQuoteLength(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v == 0 {
		return fmt.Errorf("max quote length must be greater than 0")
	}
	if v > 10000 {
		return fmt.Errorf("max quote length cannot exceed 10000")
	}
	return nil
}

func validateDefaultLanguage(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if len(v) == 0 {
		return fmt.Errorf("default language cannot be empty")
	}
	return nil
}

func validateAvailableLanguages(i interface{}) error {
	v, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if len(v) == 0 {
		return fmt.Errorf("available languages cannot be empty")
	}
	return nil
}

func validateQuoteRefreshInterval(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v <= 0 {
		return fmt.Errorf("quote refresh interval must be positive")
	}
	return nil
}

func validateMaxQuotesPerUser(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v == 0 {
		return fmt.Errorf("max quotes per user must be greater than 0")
	}
	return nil
}

func validateEnableNFTCreation(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateNFTMintingFee(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if !v.IsValid() {
		return fmt.Errorf("invalid NFT minting fee")
	}
	if v.IsNegative() {
		return fmt.Errorf("NFT minting fee cannot be negative")
	}
	return nil
}

func validateNFTRoyaltyPercentage(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	dec, err := sdk.NewDecFromStr(v)
	if err != nil {
		return fmt.Errorf("invalid royalty percentage: %s", v)
	}
	if dec.IsNegative() || dec.GT(sdk.NewDec(100)) {
		return fmt.Errorf("royalty percentage must be between 0 and 100")
	}
	return nil
}

func validateEnableStatistics(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateStatisticsRetentionDays(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v == 0 {
		return fmt.Errorf("statistics retention days must be greater than 0")
	}
	return nil
}