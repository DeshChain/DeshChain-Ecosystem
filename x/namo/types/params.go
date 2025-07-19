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

// Parameter store keys
var (
	KeyTokenDenom      = []byte("TokenDenom")
	KeyEnableVesting   = []byte("EnableVesting")
	KeyEnableBurning   = []byte("EnableBurning")
	KeyMinBurnAmount   = []byte("MinBurnAmount")
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(
	tokenDenom string,
	enableVesting bool,
	enableBurning bool,
	minBurnAmount sdk.Int,
) Params {
	return Params{
		TokenDenom:    tokenDenom,
		EnableVesting: enableVesting,
		EnableBurning: enableBurning,
		MinBurnAmount: minBurnAmount,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultTokenDenom,
		true,  // enable vesting
		true,  // enable burning
		sdk.NewInt(1000000), // minimum burn amount (1 NAMO with 6 decimals)
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyTokenDenom, &p.TokenDenom, validateTokenDenom),
		paramtypes.NewParamSetPair(KeyEnableVesting, &p.EnableVesting, validateEnableVesting),
		paramtypes.NewParamSetPair(KeyEnableBurning, &p.EnableBurning, validateEnableBurning),
		paramtypes.NewParamSetPair(KeyMinBurnAmount, &p.MinBurnAmount, validateMinBurnAmount),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateTokenDenom(p.TokenDenom); err != nil {
		return err
	}
	if err := validateEnableVesting(p.EnableVesting); err != nil {
		return err
	}
	if err := validateEnableBurning(p.EnableBurning); err != nil {
		return err
	}
	if err := validateMinBurnAmount(p.MinBurnAmount); err != nil {
		return err
	}
	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// validateTokenDenom validates the token denomination
func validateTokenDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if len(v) == 0 {
		return fmt.Errorf("token denom cannot be empty")
	}

	return sdk.ValidateDenom(v)
}

// validateEnableVesting validates the enable vesting parameter
func validateEnableVesting(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

// validateEnableBurning validates the enable burning parameter
func validateEnableBurning(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

// validateMinBurnAmount validates the minimum burn amount parameter
func validateMinBurnAmount(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("minimum burn amount cannot be nil")
	}

	if v.IsNegative() {
		return fmt.Errorf("minimum burn amount cannot be negative: %s", v)
	}

	return nil
}