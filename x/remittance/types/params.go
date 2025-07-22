package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Default parameter values
const (
	// DefaultTransferEnabled defines if remittance transfers are enabled
	DefaultTransferEnabled = true
	
	// DefaultMaxTransferAmount defines the maximum transfer amount
	DefaultMaxTransferAmount = 1000000 // 1M NAMO tokens
	
	// DefaultMinTransferAmount defines the minimum transfer amount
	DefaultMinTransferAmount = 1 // 1 NAMO token
	
	// DefaultTransferFeeRate defines the default transfer fee rate (0.5%)
	DefaultTransferFeeRate = "0.005"
	
	// DefaultMaxDailyTransfers defines maximum transfers per day per user
	DefaultMaxDailyTransfers = 10
	
	// DefaultTransferExpiryHours defines when transfers expire
	DefaultTransferExpiryHours = 24
	
	// DefaultRequireKYC defines if KYC is required
	DefaultRequireKYC = true
	
	// DefaultMinKYCLevel defines minimum KYC level required
	DefaultMinKYCLevel = KYC_LEVEL_BASIC
	
	// DefaultMaxSlippageTolerance defines maximum slippage tolerance (2%)
	DefaultMaxSlippageTolerance = "0.02"
	
	// DefaultLiquidityPoolEnabled defines if liquidity pools are enabled
	DefaultLiquidityPoolEnabled = true
	
	// DefaultMinLiquidityAmount defines minimum liquidity amount
	DefaultMinLiquidityAmount = 1000 // 1000 NAMO tokens
)

// Parameter store keys
var (
	KeyTransferEnabled          = []byte("TransferEnabled")
	KeyMaxTransferAmount        = []byte("MaxTransferAmount")
	KeyMinTransferAmount        = []byte("MinTransferAmount")
	KeyTransferFeeRate          = []byte("TransferFeeRate")
	KeyMaxDailyTransfers        = []byte("MaxDailyTransfers")
	KeyTransferExpiryHours      = []byte("TransferExpiryHours")
	KeyRequireKYC              = []byte("RequireKYC")
	KeyMinKYCLevel             = []byte("MinKYCLevel")
	KeyMaxSlippageTolerance    = []byte("MaxSlippageTolerance")
	KeyLiquidityPoolEnabled    = []byte("LiquidityPoolEnabled")
	KeyMinLiquidityAmount      = []byte("MinLiquidityAmount")
	KeySupportedCurrencies     = []byte("SupportedCurrencies")
	KeySupportedCountries      = []byte("SupportedCountries")
	KeyAuthorizedOperators     = []byte("AuthorizedOperators")
	KeyEmergencyMode          = []byte("EmergencyMode")
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&RemittanceParams{})
}

// DefaultParams returns a default set of parameters
func DefaultParams() RemittanceParams {
	transferFeeRate, _ := sdk.NewDecFromStr(DefaultTransferFeeRate)
	maxSlippageTolerance, _ := sdk.NewDecFromStr(DefaultMaxSlippageTolerance)
	
	return RemittanceParams{
		TransferEnabled:         DefaultTransferEnabled,
		MaxTransferAmount:       sdk.NewInt(DefaultMaxTransferAmount),
		MinTransferAmount:       sdk.NewInt(DefaultMinTransferAmount),
		TransferFeeRate:         transferFeeRate,
		MaxDailyTransfers:       DefaultMaxDailyTransfers,
		TransferExpiryHours:     DefaultTransferExpiryHours,
		RequireKyc:             DefaultRequireKYC,
		MinKycLevel:            DefaultMinKYCLevel,
		MaxSlippageTolerance:   maxSlippageTolerance,
		LiquidityPoolEnabled:   DefaultLiquidityPoolEnabled,
		MinLiquidityAmount:     sdk.NewInt(DefaultMinLiquidityAmount),
		SupportedCurrencies:    []string{"NAMO", "USD", "INR", "EUR", "GBP"},
		SupportedCountries:     []string{"US", "IN", "GB", "CA", "AU"},
		AuthorizedOperators:    []string{},
		EmergencyMode:          false,
	}
}

// ParamSetPairs get the params.ParamSet
func (p *RemittanceParams) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyTransferEnabled, &p.TransferEnabled, validateTransferEnabled),
		paramtypes.NewParamSetPair(KeyMaxTransferAmount, &p.MaxTransferAmount, validateMaxTransferAmount),
		paramtypes.NewParamSetPair(KeyMinTransferAmount, &p.MinTransferAmount, validateMinTransferAmount),
		paramtypes.NewParamSetPair(KeyTransferFeeRate, &p.TransferFeeRate, validateTransferFeeRate),
		paramtypes.NewParamSetPair(KeyMaxDailyTransfers, &p.MaxDailyTransfers, validateMaxDailyTransfers),
		paramtypes.NewParamSetPair(KeyTransferExpiryHours, &p.TransferExpiryHours, validateTransferExpiryHours),
		paramtypes.NewParamSetPair(KeyRequireKYC, &p.RequireKyc, validateRequireKYC),
		paramtypes.NewParamSetPair(KeyMinKYCLevel, &p.MinKycLevel, validateMinKYCLevel),
		paramtypes.NewParamSetPair(KeyMaxSlippageTolerance, &p.MaxSlippageTolerance, validateMaxSlippageTolerance),
		paramtypes.NewParamSetPair(KeyLiquidityPoolEnabled, &p.LiquidityPoolEnabled, validateLiquidityPoolEnabled),
		paramtypes.NewParamSetPair(KeyMinLiquidityAmount, &p.MinLiquidityAmount, validateMinLiquidityAmount),
		paramtypes.NewParamSetPair(KeySupportedCurrencies, &p.SupportedCurrencies, validateSupportedCurrencies),
		paramtypes.NewParamSetPair(KeySupportedCountries, &p.SupportedCountries, validateSupportedCountries),
		paramtypes.NewParamSetPair(KeyAuthorizedOperators, &p.AuthorizedOperators, validateAuthorizedOperators),
		paramtypes.NewParamSetPair(KeyEmergencyMode, &p.EmergencyMode, validateEmergencyMode),
	}
}

// Validate validates the set of params
func (p RemittanceParams) Validate() error {
	if err := validateTransferEnabled(p.TransferEnabled); err != nil {
		return err
	}
	if err := validateMaxTransferAmount(p.MaxTransferAmount); err != nil {
		return err
	}
	if err := validateMinTransferAmount(p.MinTransferAmount); err != nil {
		return err
	}
	if err := validateTransferFeeRate(p.TransferFeeRate); err != nil {
		return err
	}
	if err := validateMaxDailyTransfers(p.MaxDailyTransfers); err != nil {
		return err
	}
	if err := validateTransferExpiryHours(p.TransferExpiryHours); err != nil {
		return err
	}
	if err := validateRequireKYC(p.RequireKyc); err != nil {
		return err
	}
	if err := validateMinKYCLevel(p.MinKycLevel); err != nil {
		return err
	}
	if err := validateMaxSlippageTolerance(p.MaxSlippageTolerance); err != nil {
		return err
	}
	if err := validateLiquidityPoolEnabled(p.LiquidityPoolEnabled); err != nil {
		return err
	}
	if err := validateMinLiquidityAmount(p.MinLiquidityAmount); err != nil {
		return err
	}
	if err := validateSupportedCurrencies(p.SupportedCurrencies); err != nil {
		return err
	}
	if err := validateSupportedCountries(p.SupportedCountries); err != nil {
		return err
	}
	if err := validateAuthorizedOperators(p.AuthorizedOperators); err != nil {
		return err
	}
	if err := validateEmergencyMode(p.EmergencyMode); err != nil {
		return err
	}

	// Cross-validation
	if p.MinTransferAmount.GT(p.MaxTransferAmount) {
		return fmt.Errorf("min transfer amount cannot be greater than max transfer amount")
	}

	return nil
}

// String implements the Stringer interface.
func (p RemittanceParams) String() string {
	return fmt.Sprintf(`Remittance Params:
  Transfer Enabled: %t
  Max Transfer Amount: %s
  Min Transfer Amount: %s
  Transfer Fee Rate: %s
  Max Daily Transfers: %d
  Transfer Expiry Hours: %d
  Require KYC: %t
  Min KYC Level: %s
  Max Slippage Tolerance: %s
  Liquidity Pool Enabled: %t
  Min Liquidity Amount: %s
  Supported Currencies: %v
  Supported Countries: %v
  Authorized Operators: %v
  Emergency Mode: %t`,
		p.TransferEnabled,
		p.MaxTransferAmount.String(),
		p.MinTransferAmount.String(),
		p.TransferFeeRate.String(),
		p.MaxDailyTransfers,
		p.TransferExpiryHours,
		p.RequireKyc,
		p.MinKycLevel.String(),
		p.MaxSlippageTolerance.String(),
		p.LiquidityPoolEnabled,
		p.MinLiquidityAmount.String(),
		p.SupportedCurrencies,
		p.SupportedCountries,
		p.AuthorizedOperators,
		p.EmergencyMode,
	)
}

// Validation functions
func validateTransferEnabled(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateMaxTransferAmount(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNegative() {
		return fmt.Errorf("max transfer amount cannot be negative")
	}
	return nil
}

func validateMinTransferAmount(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNegative() {
		return fmt.Errorf("min transfer amount cannot be negative")
	}
	return nil
}

func validateTransferFeeRate(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNegative() {
		return fmt.Errorf("transfer fee rate cannot be negative")
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("transfer fee rate cannot be greater than 100%%")
	}
	return nil
}

func validateMaxDailyTransfers(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v == 0 {
		return fmt.Errorf("max daily transfers cannot be zero")
	}
	return nil
}

func validateTransferExpiryHours(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v == 0 {
		return fmt.Errorf("transfer expiry hours cannot be zero")
	}
	return nil
}

func validateRequireKYC(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateMinKYCLevel(i interface{}) error {
	v, ok := i.(KYCLevel)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v < KYC_LEVEL_NONE || v > KYC_LEVEL_ENHANCED {
		return fmt.Errorf("invalid KYC level: %d", v)
	}
	return nil
}

func validateMaxSlippageTolerance(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNegative() {
		return fmt.Errorf("max slippage tolerance cannot be negative")
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("max slippage tolerance cannot be greater than 100%%")
	}
	return nil
}

func validateLiquidityPoolEnabled(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateMinLiquidityAmount(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNegative() {
		return fmt.Errorf("min liquidity amount cannot be negative")
	}
	return nil
}

func validateSupportedCurrencies(i interface{}) error {
	v, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if len(v) == 0 {
		return fmt.Errorf("supported currencies cannot be empty")
	}
	for _, currency := range v {
		if len(currency) != 3 && len(currency) != 4 { // ISO 4217 or crypto codes
			return fmt.Errorf("invalid currency code: %s", currency)
		}
	}
	return nil
}

func validateSupportedCountries(i interface{}) error {
	v, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if len(v) == 0 {
		return fmt.Errorf("supported countries cannot be empty")
	}
	for _, country := range v {
		if len(country) != 2 { // ISO 3166-1 alpha-2
			return fmt.Errorf("invalid country code: %s", country)
		}
	}
	return nil
}

func validateAuthorizedOperators(i interface{}) error {
	v, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	for _, operator := range v {
		if _, err := sdk.AccAddressFromBech32(operator); err != nil {
			return fmt.Errorf("invalid operator address: %s", operator)
		}
	}
	return nil
}

func validateEmergencyMode(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}