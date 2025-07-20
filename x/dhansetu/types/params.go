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

// Parameter keys
var (
	KeyDhanSetuEnabled             = []byte("DhanSetuEnabled")
	KeyDhanPataRegistrationFee     = []byte("DhanPataRegistrationFee")
	KeyKshetraCoinCreationFee      = []byte("KshetraCoinCreationFee")
	KeyMitraRegistrationFee        = []byte("MitraRegistrationFee")
	KeyMaxDhanPataLength           = []byte("MaxDhanPataLength")
	KeyMinTrustScoreForMitra       = []byte("MinTrustScoreForMitra")
	KeyDhanSetuFeeRate             = []byte("DhanSetuFeeRate")
	KeyCrossModuleFeeSharing       = []byte("CrossModuleFeeSharing")
	KeyKYCRequiredAmount           = []byte("KYCRequiredAmount")
	KeyMitraCooldownPeriod         = []byte("MitraCooldownPeriod")
)

// Default parameter values
const (
	DefaultDhanSetuEnabled         = true
	DefaultDhanPataRegistrationFee = "100"      // 100 NAMO
	DefaultKshetraCoinCreationFee  = "1000"     // 1000 NAMO
	DefaultMitraRegistrationFee    = "500"      // 500 NAMO
	DefaultMaxDhanPataLength       = uint64(32)
	DefaultMinTrustScoreForMitra   = int64(50)
	DefaultDhanSetuFeeRate         = "0.005"    // 0.5%
	DefaultCrossModuleFeeSharing   = true
	DefaultKYCRequiredAmount       = "10000000000000" // 10,000 NAMO
	DefaultMitraCooldownPeriod     = int64(3600)     // 1 hour in seconds
)

// Params defines the parameters for the DhanSetu module
type Params struct {
	DhanSetuEnabled         bool    `json:"dhansetu_enabled" yaml:"dhansetu_enabled"`
	DhanPataRegistrationFee sdk.Int `json:"dhanpata_registration_fee" yaml:"dhanpata_registration_fee"`
	KshetraCoinCreationFee  sdk.Int `json:"kshetra_coin_creation_fee" yaml:"kshetra_coin_creation_fee"`
	MitraRegistrationFee    sdk.Int `json:"mitra_registration_fee" yaml:"mitra_registration_fee"`
	MaxDhanPataLength       uint64  `json:"max_dhanpata_length" yaml:"max_dhanpata_length"`
	MinTrustScoreForMitra   int64   `json:"min_trust_score_for_mitra" yaml:"min_trust_score_for_mitra"`
	DhanSetuFeeRate         sdk.Dec `json:"dhansetu_fee_rate" yaml:"dhansetu_fee_rate"`
	CrossModuleFeeSharing   bool    `json:"cross_module_fee_sharing" yaml:"cross_module_fee_sharing"`
	KYCRequiredAmount       sdk.Int `json:"kyc_required_amount" yaml:"kyc_required_amount"`
	MitraCooldownPeriod     int64   `json:"mitra_cooldown_period" yaml:"mitra_cooldown_period"`
}

// NewParams creates a new Params instance
func NewParams(
	dhanSetuEnabled bool,
	dhanPataFee, kshetraFee, mitraFee sdk.Int,
	maxLength uint64,
	minTrustScore int64,
	feeRate sdk.Dec,
	feeSharing bool,
	kycAmount sdk.Int,
	cooldownPeriod int64,
) Params {
	return Params{
		DhanSetuEnabled:         dhanSetuEnabled,
		DhanPataRegistrationFee: dhanPataFee,
		KshetraCoinCreationFee:  kshetraFee,
		MitraRegistrationFee:    mitraFee,
		MaxDhanPataLength:       maxLength,
		MinTrustScoreForMitra:   minTrustScore,
		DhanSetuFeeRate:         feeRate,
		CrossModuleFeeSharing:   feeSharing,
		KYCRequiredAmount:       kycAmount,
		MitraCooldownPeriod:     cooldownPeriod,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	dhanPataFee, _ := sdk.NewIntFromString(DefaultDhanPataRegistrationFee)
	kshetraFee, _ := sdk.NewIntFromString(DefaultKshetraCoinCreationFee)
	mitraFee, _ := sdk.NewIntFromString(DefaultMitraRegistrationFee)
	feeRate, _ := sdk.NewDecFromStr(DefaultDhanSetuFeeRate)
	kycAmount, _ := sdk.NewIntFromString(DefaultKYCRequiredAmount)

	return NewParams(
		DefaultDhanSetuEnabled,
		dhanPataFee,
		kshetraFee,
		mitraFee,
		DefaultMaxDhanPataLength,
		DefaultMinTrustScoreForMitra,
		feeRate,
		DefaultCrossModuleFeeSharing,
		kycAmount,
		DefaultMitraCooldownPeriod,
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyDhanSetuEnabled, &p.DhanSetuEnabled, validateBool),
		paramtypes.NewParamSetPair(KeyDhanPataRegistrationFee, &p.DhanPataRegistrationFee, validateInt),
		paramtypes.NewParamSetPair(KeyKshetraCoinCreationFee, &p.KshetraCoinCreationFee, validateInt),
		paramtypes.NewParamSetPair(KeyMitraRegistrationFee, &p.MitraRegistrationFee, validateInt),
		paramtypes.NewParamSetPair(KeyMaxDhanPataLength, &p.MaxDhanPataLength, validateUint64),
		paramtypes.NewParamSetPair(KeyMinTrustScoreForMitra, &p.MinTrustScoreForMitra, validateInt64),
		paramtypes.NewParamSetPair(KeyDhanSetuFeeRate, &p.DhanSetuFeeRate, validateDec),
		paramtypes.NewParamSetPair(KeyCrossModuleFeeSharing, &p.CrossModuleFeeSharing, validateBool),
		paramtypes.NewParamSetPair(KeyKYCRequiredAmount, &p.KYCRequiredAmount, validateInt),
		paramtypes.NewParamSetPair(KeyMitraCooldownPeriod, &p.MitraCooldownPeriod, validateInt64),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateBool(p.DhanSetuEnabled); err != nil {
		return err
	}
	if err := validateInt(p.DhanPataRegistrationFee); err != nil {
		return err
	}
	if err := validateInt(p.KshetraCoinCreationFee); err != nil {
		return err
	}
	if err := validateInt(p.MitraRegistrationFee); err != nil {
		return err
	}
	if err := validateUint64(p.MaxDhanPataLength); err != nil {
		return err
	}
	if err := validateInt64(p.MinTrustScoreForMitra); err != nil {
		return err
	}
	if err := validateDec(p.DhanSetuFeeRate); err != nil {
		return err
	}
	if err := validateBool(p.CrossModuleFeeSharing); err != nil {
		return err
	}
	if err := validateInt(p.KYCRequiredAmount); err != nil {
		return err
	}
	if err := validateInt64(p.MitraCooldownPeriod); err != nil {
		return err
	}
	return nil
}

// String implements the Stringer interface
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// ParamKeyTable for DhanSetu module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// Validation functions

func validateBool(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateInt(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNegative() {
		return fmt.Errorf("invalid parameter value: cannot be negative")
	}
	return nil
}

func validateUint64(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v == 0 {
		return fmt.Errorf("invalid parameter value: cannot be zero")
	}
	return nil
}

func validateInt64(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v < 0 {
		return fmt.Errorf("invalid parameter value: cannot be negative")
	}
	return nil
}

func validateDec(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNegative() {
		return fmt.Errorf("invalid parameter value: cannot be negative")
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("invalid parameter value: cannot be greater than 1.0")
	}
	return nil
}