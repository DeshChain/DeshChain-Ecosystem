package types

import (
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Default parameter values
const (
	DefaultEnabled = true
)

// Parameter store keys
var (
	KeyEnabled = []byte("Enabled")
)

// Params defines the parameters for the donation module
type Params struct {
	Enabled bool `json:"enabled" yaml:"enabled"`
}

// NewParams creates a new Params instance
func NewParams(enabled bool) Params {
	return Params{
		Enabled: enabled,
	}
}

// DefaultParams returns default parameters
func DefaultParams() Params {
	return Params{
		Enabled: DefaultEnabled,
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
	}
}

// Validate validates the parameters
func (p Params) Validate() error {
	if err := validateEnabled(p.Enabled); err != nil {
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