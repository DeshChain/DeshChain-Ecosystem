package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

// Parameter store keys
var (
	KeyLendingBasicThreshold     = []byte("LendingBasicThreshold")
	KeyLendingMediumThreshold    = []byte("LendingMediumThreshold")
	KeyLendingFullThreshold      = []byte("LendingFullThreshold")
	KeyMinimumReserveRatio       = []byte("MinimumReserveRatio")
	KeyEmergencyReserveRatio     = []byte("EmergencyReserveRatio")
	KeyLoanLossProvisionRatio    = []byte("LoanLossProvisionRatio")
	KeyProcessingFeeRate         = []byte("ProcessingFeeRate")
	KeyProcessingFeeCap          = []byte("ProcessingFeeCap")
	KeyEarlySettlementRate       = []byte("EarlySettlementRate")
)

// REVOLUTIONARY FINANCIALLY SOUND PARAMETERS - Based on Deep Risk Analysis
const (
	// CONSERVATIVE CAPITAL REQUIREMENTS - Prevents Systemic Risk
	DefaultLendingBasicThreshold  = "100000000000000" // ₹100 Cr in micro-NAMO (10x safer)
	DefaultLendingMediumThreshold = "250000000000000" // ₹250 Cr in micro-NAMO (10x safer)
	DefaultLendingFullThreshold   = "500000000000000" // ₹500 Cr in micro-NAMO (10x safer)
	
	// ENHANCED RISK RESERVES - Protects Against Defaults & Liquidity Crises
	DefaultMinimumReserveRatio    = "0.50"            // 50% minimum reserve (Bank-grade safety)
	DefaultEmergencyReserveRatio  = "0.15"            // 15% emergency reserve + 10% default provision
	DefaultLoanLossProvisionRatio = "0.10"            // 10% loan loss provision (Industry standard)
	
	// REVOLUTIONARY FEE STRUCTURE - Platform Sustainability
	DefaultProcessingFeeRate      = "0.01"            // 1% processing fee (ultra-competitive)
	DefaultProcessingFeeCap       = "2500000000"      // ₹2500 maximum cap in micro-NAMO
	DefaultEarlySettlementRate    = "0.005"           // 0.5% early settlement fee (borrower-friendly)
)

// LiquidityParams defines the parameters for the liquidity manager module with enhanced risk management
type LiquidityParams struct {
	LendingBasicThreshold     sdk.Int `json:"lending_basic_threshold" yaml:"lending_basic_threshold"`
	LendingMediumThreshold    sdk.Int `json:"lending_medium_threshold" yaml:"lending_medium_threshold"`
	LendingFullThreshold      sdk.Int `json:"lending_full_threshold" yaml:"lending_full_threshold"`
	MinimumReserveRatio       sdk.Dec `json:"minimum_reserve_ratio" yaml:"minimum_reserve_ratio"`
	EmergencyReserveRatio     sdk.Dec `json:"emergency_reserve_ratio" yaml:"emergency_reserve_ratio"`
	LoanLossProvisionRatio    sdk.Dec `json:"loan_loss_provision_ratio" yaml:"loan_loss_provision_ratio"`
	
	// REVOLUTIONARY FEE STRUCTURE
	ProcessingFeeRate         sdk.Dec `json:"processing_fee_rate" yaml:"processing_fee_rate"`
	ProcessingFeeCap          sdk.Int `json:"processing_fee_cap" yaml:"processing_fee_cap"`
	EarlySettlementRate       sdk.Dec `json:"early_settlement_rate" yaml:"early_settlement_rate"`
}

// NewLiquidityParams creates a new LiquidityParams instance with enhanced risk management and fee structure
func NewLiquidityParams(
	basicThreshold, mediumThreshold, fullThreshold sdk.Int,
	minReserve, emergencyReserve, loanLossProvision sdk.Dec,
	processingFeeRate sdk.Dec, processingFeeCap sdk.Int, earlySettlementRate sdk.Dec,
) LiquidityParams {
	return LiquidityParams{
		LendingBasicThreshold:     basicThreshold,
		LendingMediumThreshold:    mediumThreshold,
		LendingFullThreshold:      fullThreshold,
		MinimumReserveRatio:       minReserve,
		EmergencyReserveRatio:     emergencyReserve,
		LoanLossProvisionRatio:    loanLossProvision,
		ProcessingFeeRate:         processingFeeRate,
		ProcessingFeeCap:          processingFeeCap,
		EarlySettlementRate:       earlySettlementRate,
	}
}

// DefaultLiquidityParams returns a default set of parameters with enhanced risk management and revolutionary fee structure
func DefaultLiquidityParams() LiquidityParams {
	basicThreshold, _ := sdk.NewIntFromString(DefaultLendingBasicThreshold)
	mediumThreshold, _ := sdk.NewIntFromString(DefaultLendingMediumThreshold)
	fullThreshold, _ := sdk.NewIntFromString(DefaultLendingFullThreshold)
	minReserve, _ := sdk.NewDecFromStr(DefaultMinimumReserveRatio)
	emergencyReserve, _ := sdk.NewDecFromStr(DefaultEmergencyReserveRatio)
	loanLossProvision, _ := sdk.NewDecFromStr(DefaultLoanLossProvisionRatio)
	processingFeeRate, _ := sdk.NewDecFromStr(DefaultProcessingFeeRate)
	processingFeeCap, _ := sdk.NewIntFromString(DefaultProcessingFeeCap)
	earlySettlementRate, _ := sdk.NewDecFromStr(DefaultEarlySettlementRate)

	return NewLiquidityParams(
		basicThreshold,
		mediumThreshold,
		fullThreshold,
		minReserve,
		emergencyReserve,
		loanLossProvision,
		processingFeeRate,
		processingFeeCap,
		earlySettlementRate,
	)
}

// ParamKeyTable returns the parameter key table for use with the sdk.Params module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&LiquidityParams{})
}

// ParamSetPairs returns the parameter set pairs with enhanced risk management and fee structure
func (p *LiquidityParams) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyLendingBasicThreshold, &p.LendingBasicThreshold, validateThreshold),
		paramtypes.NewParamSetPair(KeyLendingMediumThreshold, &p.LendingMediumThreshold, validateThreshold),
		paramtypes.NewParamSetPair(KeyLendingFullThreshold, &p.LendingFullThreshold, validateThreshold),
		paramtypes.NewParamSetPair(KeyMinimumReserveRatio, &p.MinimumReserveRatio, validateRatio),
		paramtypes.NewParamSetPair(KeyEmergencyReserveRatio, &p.EmergencyReserveRatio, validateRatio),
		paramtypes.NewParamSetPair(KeyLoanLossProvisionRatio, &p.LoanLossProvisionRatio, validateRatio),
		paramtypes.NewParamSetPair(KeyProcessingFeeRate, &p.ProcessingFeeRate, validateFeeRate),
		paramtypes.NewParamSetPair(KeyProcessingFeeCap, &p.ProcessingFeeCap, validateFeeCap),
		paramtypes.NewParamSetPair(KeyEarlySettlementRate, &p.EarlySettlementRate, validateFeeRate),
	}
}

// Validate validates the set of params with enhanced risk management and fee structure checks
func (p LiquidityParams) Validate() error {
	if err := validateThreshold(p.LendingBasicThreshold); err != nil {
		return err
	}
	if err := validateThreshold(p.LendingMediumThreshold); err != nil {
		return err
	}
	if err := validateThreshold(p.LendingFullThreshold); err != nil {
		return err
	}
	if err := validateRatio(p.MinimumReserveRatio); err != nil {
		return err
	}
	if err := validateRatio(p.EmergencyReserveRatio); err != nil {
		return err
	}
	if err := validateRatio(p.LoanLossProvisionRatio); err != nil {
		return err
	}
	if err := validateFeeRate(p.ProcessingFeeRate); err != nil {
		return err
	}
	if err := validateFeeCap(p.ProcessingFeeCap); err != nil {
		return err
	}
	if err := validateFeeRate(p.EarlySettlementRate); err != nil {
		return err
	}

	// Validate threshold ordering
	if p.LendingMediumThreshold.LTE(p.LendingBasicThreshold) {
		return fmt.Errorf("medium threshold must be greater than basic threshold")
	}
	if p.LendingFullThreshold.LTE(p.LendingMediumThreshold) {
		return fmt.Errorf("full threshold must be greater than medium threshold")
	}

	// REVOLUTIONARY: Validate financial soundness
	totalReserveRatio := p.MinimumReserveRatio.Add(p.EmergencyReserveRatio).Add(p.LoanLossProvisionRatio)
	if totalReserveRatio.GTE(sdk.OneDec()) {
		return fmt.Errorf("total reserve ratios cannot exceed 100%: %s", totalReserveRatio.String())
	}

	// Ensure minimum financial safety
	if p.MinimumReserveRatio.LT(sdk.NewDecWithPrec(30, 2)) {
		return fmt.Errorf("minimum reserve ratio must be at least 30%")
	}

	// REVOLUTIONARY: Validate fee structure competitiveness
	if p.ProcessingFeeRate.GT(sdk.NewDecWithPrec(3, 2)) {
		return fmt.Errorf("processing fee rate cannot exceed 3% (industry competitive limit)")
	}
	if p.EarlySettlementRate.GT(sdk.NewDecWithPrec(2, 2)) {
		return fmt.Errorf("early settlement rate cannot exceed 2% (borrower protection)")
	}

	return nil
}

// String implements the Stringer interface
func (p LiquidityParams) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

func validateThreshold(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("threshold cannot be negative: %s", v)
	}

	return nil
}

func validateRatio(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("ratio cannot be negative: %s", v)
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("ratio cannot be greater than 1: %s", v)
	}

	return nil
}

// validateFeeRate validates fee rates with borrower protection limits
func validateFeeRate(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("fee rate cannot be negative: %s", v)
	}

	// REVOLUTIONARY: Cap at 5% for borrower protection (vs industry 3-8%)
	maxFeeRate := sdk.NewDecWithPrec(5, 2) // 5%
	if v.GT(maxFeeRate) {
		return fmt.Errorf("fee rate cannot exceed 5%% (borrower protection): %s", v)
	}

	return nil
}

// validateFeeCap validates processing fee caps
func validateFeeCap(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("fee cap cannot be negative: %s", v)
	}

	// REVOLUTIONARY: Maximum cap of ₹10,000 for borrower protection
	maxFeeCap := sdk.NewInt(10000000000) // ₹10K in micro-NAMO
	if v.GT(maxFeeCap) {
		return fmt.Errorf("fee cap cannot exceed ₹10,000 (borrower protection): %s", v)
	}

	return nil
}