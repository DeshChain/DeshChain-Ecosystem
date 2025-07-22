package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// Parameter keys
var (
	KeyFees                   = []byte("Fees")
	KeyMinCollateralRatio     = []byte("MinCollateralRatio")
	KeyLiquidationThreshold   = []byte("LiquidationThreshold")
	KeyMaxSupply              = []byte("MaxSupply")
	KeyCollateralAssets       = []byte("CollateralAssets")
	KeyMintingEnabled         = []byte("MintingEnabled")
	KeyBurningEnabled         = []byte("BurningEnabled")
	KeyOracleUpdateFrequency  = []byte("OracleUpdateFrequency")
	KeyInsuranceFundRatio     = []byte("InsuranceFundRatio")
)

// Default parameter values
const (
	DefaultMintFee              = uint64(10)    // 0.1% in basis points
	DefaultMintFeeCap           = "100"         // 100 DINR cap
	DefaultBurnFee              = uint64(10)    // 0.1% in basis points
	DefaultBurnFeeCap           = "100"         // 100 DINR cap
	DefaultLiquidationPenalty   = uint64(1000)  // 10% in basis points
	DefaultStabilityFee         = uint64(200)   // 2% annual in basis points
	DefaultMinCollateralRatio   = uint64(15000) // 150% in basis points
	DefaultLiquidationThreshold = uint64(13000) // 130% in basis points
	DefaultMaxSupply            = uint64(1000000000) // 1 billion DINR
	DefaultOracleUpdateFrequency = uint64(300)  // 5 minutes
	DefaultInsuranceFundRatio   = "0.05"        // 5% of total supply
	DefaultYieldToInsuranceRatio = uint64(2000) // 20% of yield to insurance fund
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(
	fees FeeStructure,
	minCollateralRatio uint64,
	liquidationThreshold uint64,
	maxSupply uint64,
	collateralAssets []CollateralAsset,
	mintingEnabled bool,
	burningEnabled bool,
	oracleUpdateFrequency uint64,
	insuranceFundRatio string,
) Params {
	return Params{
		Fees:                  fees,
		MinCollateralRatio:    minCollateralRatio,
		LiquidationThreshold:  liquidationThreshold,
		MaxSupply:             maxSupply,
		CollateralAssets:      collateralAssets,
		MintingEnabled:        mintingEnabled,
		BurningEnabled:        burningEnabled,
		OracleUpdateFrequency: oracleUpdateFrequency,
		InsuranceFundRatio:    insuranceFundRatio,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	defaultFees := FeeStructure{
		MintFee:            DefaultMintFee,
		MintFeeCap:         DefaultMintFeeCap,
		BurnFee:            DefaultBurnFee,
		BurnFeeCap:         DefaultBurnFeeCap,
		LiquidationPenalty: DefaultLiquidationPenalty,
		StabilityFee:       DefaultStabilityFee,
	}

	// Default collateral assets with NO NAMO to prevent circular dependency
	defaultCollateralAssets := []CollateralAsset{
		{
			Denom:               "usdt",
			Tier:                "tier1_stable",
			MinCollateralRatio:  14000, // 140%
			MaxAllocation:       2500,  // 25%
			IsActive:            true,
			OracleScriptId:      "usdt_inr",
		},
		{
			Denom:               "usdc",
			Tier:                "tier1_stable",
			MinCollateralRatio:  14000, // 140%
			MaxAllocation:       2500,  // 25%
			IsActive:            true,
			OracleScriptId:      "usdc_inr",
		},
		{
			Denom:               "dai",
			Tier:                "tier1_stable",
			MinCollateralRatio:  14000, // 140%
			MaxAllocation:       1000,  // 10%
			IsActive:            true,
			OracleScriptId:      "dai_inr",
		},
		{
			Denom:               "btc",
			Tier:                "tier2_crypto",
			MinCollateralRatio:  15000, // 150%
			MaxAllocation:       2000,  // 20%
			IsActive:            true,
			OracleScriptId:      "btc_inr",
		},
		{
			Denom:               "eth",
			Tier:                "tier2_crypto",
			MinCollateralRatio:  15000, // 150%
			MaxAllocation:       2000,  // 20%
			IsActive:            true,
			OracleScriptId:      "eth_inr",
		},
		{
			Denom:               "bnb",
			Tier:                "tier3_alt",
			MinCollateralRatio:  17000, // 170%
			MaxAllocation:       500,   // 5%
			IsActive:            true,
			OracleScriptId:      "bnb_inr",
		},
		{
			Denom:               "matic",
			Tier:                "tier3_alt",
			MinCollateralRatio:  17000, // 170%
			MaxAllocation:       500,   // 5%
			IsActive:            true,
			OracleScriptId:      "matic_inr",
		},
	}

	return NewParams(
		defaultFees,
		DefaultMinCollateralRatio,
		DefaultLiquidationThreshold,
		DefaultMaxSupply,
		defaultCollateralAssets,
		true,  // minting enabled
		true,  // burning enabled
		DefaultOracleUpdateFrequency,
		DefaultInsuranceFundRatio,
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyFees, &p.Fees, validateFees),
		paramtypes.NewParamSetPair(KeyMinCollateralRatio, &p.MinCollateralRatio, validateMinCollateralRatio),
		paramtypes.NewParamSetPair(KeyLiquidationThreshold, &p.LiquidationThreshold, validateLiquidationThreshold),
		paramtypes.NewParamSetPair(KeyMaxSupply, &p.MaxSupply, validateMaxSupply),
		paramtypes.NewParamSetPair(KeyCollateralAssets, &p.CollateralAssets, validateCollateralAssets),
		paramtypes.NewParamSetPair(KeyMintingEnabled, &p.MintingEnabled, validateBool),
		paramtypes.NewParamSetPair(KeyBurningEnabled, &p.BurningEnabled, validateBool),
		paramtypes.NewParamSetPair(KeyOracleUpdateFrequency, &p.OracleUpdateFrequency, validateOracleUpdateFrequency),
		paramtypes.NewParamSetPair(KeyInsuranceFundRatio, &p.InsuranceFundRatio, validateInsuranceFundRatio),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateFees(p.Fees); err != nil {
		return err
	}
	if err := validateMinCollateralRatio(p.MinCollateralRatio); err != nil {
		return err
	}
	if err := validateLiquidationThreshold(p.LiquidationThreshold); err != nil {
		return err
	}
	if err := validateMaxSupply(p.MaxSupply); err != nil {
		return err
	}
	if err := validateCollateralAssets(p.CollateralAssets); err != nil {
		return err
	}
	if err := validateOracleUpdateFrequency(p.OracleUpdateFrequency); err != nil {
		return err
	}
	if err := validateInsuranceFundRatio(p.InsuranceFundRatio); err != nil {
		return err
	}

	// Additional validation: liquidation threshold must be less than min collateral ratio
	if p.LiquidationThreshold >= p.MinCollateralRatio {
		return fmt.Errorf("liquidation threshold must be less than minimum collateral ratio")
	}

	return nil
}

// String implements the Stringer interface
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// Validation functions
func validateFees(i interface{}) error {
	fees, ok := i.(FeeStructure)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if fees.MintFee > 1000 { // Max 10%
		return fmt.Errorf("mint fee cannot exceed 10%%")
	}
	if fees.BurnFee > 1000 { // Max 10%
		return fmt.Errorf("burn fee cannot exceed 10%%")
	}
	if fees.LiquidationPenalty > 2000 { // Max 20%
		return fmt.Errorf("liquidation penalty cannot exceed 20%%")
	}

	return nil
}

func validateMinCollateralRatio(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 10000 { // Minimum 100%
		return fmt.Errorf("minimum collateral ratio cannot be less than 100%%")
	}
	if v > 30000 { // Maximum 300%
		return fmt.Errorf("minimum collateral ratio cannot exceed 300%%")
	}

	return nil
}

func validateLiquidationThreshold(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 10000 { // Minimum 100%
		return fmt.Errorf("liquidation threshold cannot be less than 100%%")
	}
	if v > 25000 { // Maximum 250%
		return fmt.Errorf("liquidation threshold cannot exceed 250%%")
	}

	return nil
}

func validateMaxSupply(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("max supply must be greater than 0")
	}

	return nil
}

func validateCollateralAssets(i interface{}) error {
	assets, ok := i.([]CollateralAsset)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if len(assets) == 0 {
		return fmt.Errorf("at least one collateral asset must be defined")
	}

	denomMap := make(map[string]bool)
	totalMaxAllocation := uint64(0)

	for _, asset := range assets {
		// Check for duplicate denoms
		if denomMap[asset.Denom] {
			return fmt.Errorf("duplicate collateral asset denom: %s", asset.Denom)
		}
		denomMap[asset.Denom] = true

		// Validate individual asset
		if asset.MinCollateralRatio < 10000 {
			return fmt.Errorf("collateral ratio for %s cannot be less than 100%%", asset.Denom)
		}
		if asset.MaxAllocation > 10000 {
			return fmt.Errorf("max allocation for %s cannot exceed 100%%", asset.Denom)
		}

		totalMaxAllocation += asset.MaxAllocation
	}

	// Total max allocation can exceed 100% as it's a cap per asset
	return nil
}

func validateBool(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateOracleUpdateFrequency(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 60 { // Minimum 1 minute
		return fmt.Errorf("oracle update frequency cannot be less than 60 seconds")
	}
	if v > 3600 { // Maximum 1 hour
		return fmt.Errorf("oracle update frequency cannot exceed 3600 seconds")
	}

	return nil
}

func validateInsuranceFundRatio(i interface{}) error {
	_, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	// Additional decimal validation could be added here
	return nil
}