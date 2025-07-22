package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*OracleParams)(nil)

var (
	KeyMaxPriceDeviation  = []byte("MaxPriceDeviation")
	KeyMinValidators      = []byte("MinValidators")
	KeyAggregationWindow  = []byte("AggregationWindow")
	KeyHeartbeatInterval  = []byte("HeartbeatInterval")
	KeySlashFraction      = []byte("SlashFraction")
)

// ParamKeyTable returns the parameter key table for the oracle module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&OracleParams{})
}

// NewOracleParams creates a new OracleParams instance
func NewOracleParams(
	maxPriceDeviation sdk.Dec,
	minValidators uint64,
	aggregationWindow uint64,
	heartbeatInterval uint64,
	slashFraction sdk.Dec,
) OracleParams {
	return OracleParams{
		MaxPriceDeviation: maxPriceDeviation,
		MinValidators:     minValidators,
		AggregationWindow: aggregationWindow,
		HeartbeatInterval: heartbeatInterval,
		SlashFraction:     slashFraction,
	}
}

// DefaultOracleParams returns the default oracle parameters
func DefaultOracleParams() OracleParams {
	return NewOracleParams(
		sdk.NewDecWithPrec(5, 2),  // 5% max price deviation
		3,                         // minimum 3 validators required
		5,                         // 5 block aggregation window
		10,                        // heartbeat every 10 blocks
		sdk.NewDecWithPrec(1, 4),  // 0.01% slash fraction
	)
}

// ParamSetPairs implements the ParamSet interface and returns the key/value pairs
func (p *OracleParams) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMaxPriceDeviation, &p.MaxPriceDeviation, validateMaxPriceDeviation),
		paramtypes.NewParamSetPair(KeyMinValidators, &p.MinValidators, validateMinValidators),
		paramtypes.NewParamSetPair(KeyAggregationWindow, &p.AggregationWindow, validateAggregationWindow),
		paramtypes.NewParamSetPair(KeyHeartbeatInterval, &p.HeartbeatInterval, validateHeartbeatInterval),
		paramtypes.NewParamSetPair(KeySlashFraction, &p.SlashFraction, validateSlashFraction),
	}
}

// Validate validates the oracle parameters
func (p OracleParams) Validate() error {
	if err := validateMaxPriceDeviation(p.MaxPriceDeviation); err != nil {
		return err
	}
	if err := validateMinValidators(p.MinValidators); err != nil {
		return err
	}
	if err := validateAggregationWindow(p.AggregationWindow); err != nil {
		return err
	}
	if err := validateHeartbeatInterval(p.HeartbeatInterval); err != nil {
		return err
	}
	if err := validateSlashFraction(p.SlashFraction); err != nil {
		return err
	}
	return nil
}

// String implements the Stringer interface
func (p OracleParams) String() string {
	return fmt.Sprintf(`Oracle Params:
  MaxPriceDeviation: %s
  MinValidators: %d
  AggregationWindow: %d
  HeartbeatInterval: %d
  SlashFraction: %s`,
		p.MaxPriceDeviation, p.MinValidators, p.AggregationWindow, p.HeartbeatInterval, p.SlashFraction)
}

func validateMaxPriceDeviation(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("max price deviation cannot be nil")
	}

	if v.IsNegative() {
		return fmt.Errorf("max price deviation cannot be negative: %s", v)
	}

	if v.GT(sdk.NewDecWithPrec(50, 2)) { // 50%
		return fmt.Errorf("max price deviation cannot be greater than 50%%: %s", v)
	}

	return nil
}

func validateMinValidators(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("min validators cannot be zero")
	}

	if v > 100 {
		return fmt.Errorf("min validators cannot be greater than 100: %d", v)
	}

	return nil
}

func validateAggregationWindow(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("aggregation window cannot be zero")
	}

	if v > 100 {
		return fmt.Errorf("aggregation window cannot be greater than 100 blocks: %d", v)
	}

	return nil
}

func validateHeartbeatInterval(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("heartbeat interval cannot be zero")
	}

	if v > 1000 {
		return fmt.Errorf("heartbeat interval cannot be greater than 1000 blocks: %d", v)
	}

	return nil
}

func validateSlashFraction(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("slash fraction cannot be nil")
	}

	if v.IsNegative() {
		return fmt.Errorf("slash fraction cannot be negative: %s", v)
	}

	if v.GT(sdk.NewDecWithPrec(10, 2)) { // 10%
		return fmt.Errorf("slash fraction cannot be greater than 10%%: %s", v)
	}

	return nil
}