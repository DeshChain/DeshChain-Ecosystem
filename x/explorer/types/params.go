package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Default parameter values
const (
	DefaultEnableRealTimeUpdates = true
	DefaultEnableCulturalQuotes  = true
	DefaultEnablePatriotismScores = true
	DefaultEnableHolderRankings  = true
	DefaultEnableDonationTracking = true
	DefaultEnableAnalytics       = true
	DefaultEnableNotifications   = true
)

// Parameter keys
var (
	KeyEnableRealTimeUpdates = []byte("EnableRealTimeUpdates")
	KeyEnableCulturalQuotes  = []byte("EnableCulturalQuotes")
	KeyEnablePatriotismScores = []byte("EnablePatriotismScores")
	KeyEnableHolderRankings  = []byte("EnableHolderRankings")
	KeyEnableDonationTracking = []byte("EnableDonationTracking")
	KeyMaxSearchResults      = []byte("MaxSearchResults")
	KeyCacheDuration         = []byte("CacheDuration")
	KeyIndexingBatchSize     = []byte("IndexingBatchSize")
	KeyEnableAnalytics       = []byte("EnableAnalytics")
	KeyEnableNotifications   = []byte("EnableNotifications")
)

// DefaultParams returns default parameters
func DefaultParams() ExplorerParams {
	return ExplorerParams{
		EnableRealTimeUpdates: DefaultEnableRealTimeUpdates,
		EnableCulturalQuotes:  DefaultEnableCulturalQuotes,
		EnablePatriotismScores: DefaultEnablePatriotismScores,
		EnableHolderRankings:  DefaultEnableHolderRankings,
		EnableDonationTracking: DefaultEnableDonationTracking,
		MaxSearchResults:      DefaultMaxSearchResults,
		CacheDuration:         DefaultCacheDuration,
		IndexingBatchSize:     DefaultIndexingBatchSize,
		EnableAnalytics:       DefaultEnableAnalytics,
		EnableNotifications:   DefaultEnableNotifications,
	}
}

// ParamKeyTable returns the parameter key table
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&ExplorerParams{})
}

// ParamSetPairs implements the ParamSet interface
func (p *ExplorerParams) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyEnableRealTimeUpdates, &p.EnableRealTimeUpdates, validateBool),
		paramtypes.NewParamSetPair(KeyEnableCulturalQuotes, &p.EnableCulturalQuotes, validateBool),
		paramtypes.NewParamSetPair(KeyEnablePatriotismScores, &p.EnablePatriotismScores, validateBool),
		paramtypes.NewParamSetPair(KeyEnableHolderRankings, &p.EnableHolderRankings, validateBool),
		paramtypes.NewParamSetPair(KeyEnableDonationTracking, &p.EnableDonationTracking, validateBool),
		paramtypes.NewParamSetPair(KeyMaxSearchResults, &p.MaxSearchResults, validateUint32),
		paramtypes.NewParamSetPair(KeyCacheDuration, &p.CacheDuration, validateUint64),
		paramtypes.NewParamSetPair(KeyIndexingBatchSize, &p.IndexingBatchSize, validateUint32),
		paramtypes.NewParamSetPair(KeyEnableAnalytics, &p.EnableAnalytics, validateBool),
		paramtypes.NewParamSetPair(KeyEnableNotifications, &p.EnableNotifications, validateBool),
	}
}

// Validate validates the parameters
func (p ExplorerParams) Validate() error {
	if err := validateBool(p.EnableRealTimeUpdates); err != nil {
		return fmt.Errorf("invalid EnableRealTimeUpdates: %w", err)
	}

	if err := validateBool(p.EnableCulturalQuotes); err != nil {
		return fmt.Errorf("invalid EnableCulturalQuotes: %w", err)
	}

	if err := validateBool(p.EnablePatriotismScores); err != nil {
		return fmt.Errorf("invalid EnablePatriotismScores: %w", err)
	}

	if err := validateBool(p.EnableHolderRankings); err != nil {
		return fmt.Errorf("invalid EnableHolderRankings: %w", err)
	}

	if err := validateBool(p.EnableDonationTracking); err != nil {
		return fmt.Errorf("invalid EnableDonationTracking: %w", err)
	}

	if err := validateUint32(p.MaxSearchResults); err != nil {
		return fmt.Errorf("invalid MaxSearchResults: %w", err)
	}

	if err := validateUint64(p.CacheDuration); err != nil {
		return fmt.Errorf("invalid CacheDuration: %w", err)
	}

	if err := validateUint32(p.IndexingBatchSize); err != nil {
		return fmt.Errorf("invalid IndexingBatchSize: %w", err)
	}

	if err := validateBool(p.EnableAnalytics); err != nil {
		return fmt.Errorf("invalid EnableAnalytics: %w", err)
	}

	if err := validateBool(p.EnableNotifications); err != nil {
		return fmt.Errorf("invalid EnableNotifications: %w", err)
	}

	// Additional validation
	if p.MaxSearchResults == 0 {
		return fmt.Errorf("max search results must be greater than 0")
	}

	if p.MaxSearchResults > MaxSearchResults {
		return fmt.Errorf("max search results cannot exceed %d", MaxSearchResults)
	}

	if p.CacheDuration == 0 {
		return fmt.Errorf("cache duration must be greater than 0")
	}

	if p.IndexingBatchSize == 0 {
		return fmt.Errorf("indexing batch size must be greater than 0")
	}

	if p.IndexingBatchSize > MaxIndexingBatchSize {
		return fmt.Errorf("indexing batch size cannot exceed %d", MaxIndexingBatchSize)
	}

	return nil
}

// String returns a string representation of the parameters
func (p ExplorerParams) String() string {
	return fmt.Sprintf(`Explorer Params:
  EnableRealTimeUpdates: %v
  EnableCulturalQuotes: %v
  EnablePatriotismScores: %v
  EnableHolderRankings: %v
  EnableDonationTracking: %v
  MaxSearchResults: %d
  CacheDuration: %d
  IndexingBatchSize: %d
  EnableAnalytics: %v
  EnableNotifications: %v`,
		p.EnableRealTimeUpdates,
		p.EnableCulturalQuotes,
		p.EnablePatriotismScores,
		p.EnableHolderRankings,
		p.EnableDonationTracking,
		p.MaxSearchResults,
		p.CacheDuration,
		p.IndexingBatchSize,
		p.EnableAnalytics,
		p.EnableNotifications,
	)
}

// Validation functions

func validateBool(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateUint32(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v == 0 {
		return fmt.Errorf("parameter must be greater than 0")
	}
	return nil
}

func validateUint64(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v == 0 {
		return fmt.Errorf("parameter must be greater than 0")
	}
	return nil
}

// Additional constants for validation
const (
	MaxIndexingBatchSize = 10000
	MinCacheDuration     = 1     // 1 second
	MaxCacheDuration     = 86400 // 24 hours
	MinIndexingBatchSize = 1
	MinSearchResults     = 1
)