package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Default parameter values
const (
	DefaultEnabled = true
	DefaultTaxRate = "0.025" // 2.5%
)

// Parameter store keys
var (
	KeyEnabled = []byte("Enabled")
	KeyTaxRate = []byte("TaxRate")
)

// Params defines the parameters for the tax module
type Params struct {
	Enabled bool   `json:"enabled" yaml:"enabled"`
	TaxRate sdk.Dec `json:"tax_rate" yaml:"tax_rate"`
}

// NewParams creates a new Params instance
func NewParams(enabled bool, taxRate sdk.Dec) Params {
	return Params{
		Enabled: enabled,
		TaxRate: taxRate,
	}
}

// DefaultParams returns default parameters
func DefaultParams() Params {
	return Params{
		Enabled: DefaultEnabled,
		TaxRate: sdk.MustNewDecFromStr(DefaultTaxRate),
	}
}

// ParamKeyTable returns the parameter key table
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs returns the parameter set pairs
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyEnabled, &p.Enabled, validateEnabled),
		paramtypes.NewParamSetPair(KeyTaxRate, &p.TaxRate, validateTaxRate),
	}
}

// Validate validates the parameters
func (p Params) Validate() error {
	if err := validateEnabled(p.Enabled); err != nil {
		return err
	}
	if err := validateTaxRate(p.TaxRate); err != nil {
		return err
	}
	return nil
}

func validateEnabled(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return ErrInvalidType
	}
	return nil
}

func validateTaxRate(i interface{}) error {
	rate, ok := i.(sdk.Dec)
	if !ok {
		return ErrInvalidType
	}
	
	if rate.IsNegative() || rate.GT(sdk.OneDec()) {
		return ErrInvalidTaxRate
	}
	
	return nil
}