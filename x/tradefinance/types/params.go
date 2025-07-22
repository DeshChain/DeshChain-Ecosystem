package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// Parameter keys
var (
	KeyModuleEnabled                    = []byte("ModuleEnabled")
	KeyFees                            = []byte("Fees")
	KeyMinLcAmount                     = []byte("MinLcAmount")
	KeyMaxLcDurationDays               = []byte("MaxLcDurationDays")
	KeyDocumentSubmissionDeadlineHours = []byte("DocumentSubmissionDeadlineHours")
	KeySupportedCurrencies             = []byte("SupportedCurrencies")
	KeySupportedDocumentTypes          = []byte("SupportedDocumentTypes")
	KeyCollateralRatio                 = []byte("CollateralRatio")
	KeyAutoReleaseEnabled              = []byte("AutoReleaseEnabled")
	KeyDisputeResolutionPeriodHours    = []byte("DisputeResolutionPeriodHours")
)

// Default parameter values
const (
	DefaultMinLcAmount                     = uint64(1000)     // 1000 DINR minimum
	DefaultMaxLcDurationDays               = uint64(365)      // 1 year maximum
	DefaultDocumentSubmissionDeadlineHours = uint64(72)       // 3 days
	DefaultCollateralRatio                 = uint64(11000)    // 110% in basis points
	DefaultDisputeResolutionPeriodHours    = uint64(168)      // 7 days
)

// ParamKeyTable the param key table for the trade finance module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	defaultFees := TradeFinanceFees{
		LcIssuanceFee:           50,   // 0.5% in basis points
		LcAmendmentFee:          25,   // 0.25% in basis points
		DocumentVerificationFee: 100,  // 100 DINR flat fee per document
		SwiftMessageFee:         50,   // 50 DINR flat fee per SWIFT message
		InsuranceProcessingFee:  10,   // 0.1% in basis points
		EarlyPaymentDiscount:    50,   // 0.5% discount for early payment
		LatePaymentPenalty:      100,  // 1% penalty for late payment
	}

	defaultCurrencies := []string{"dinr", "usd", "eur", "inr", "gbp"}
	defaultDocuments := []string{
		"commercial_invoice",
		"bill_of_lading",
		"packing_list",
		"certificate_of_origin",
		"insurance_certificate",
		"inspection_certificate",
		"weight_certificate",
		"health_certificate",
	}

	return Params{
		ModuleEnabled:                    true,
		Fees:                            defaultFees,
		MinLcAmount:                     DefaultMinLcAmount,
		MaxLcDurationDays:               DefaultMaxLcDurationDays,
		DocumentSubmissionDeadlineHours: DefaultDocumentSubmissionDeadlineHours,
		SupportedCurrencies:             defaultCurrencies,
		SupportedDocumentTypes:          defaultDocuments,
		CollateralRatio:                 DefaultCollateralRatio,
		AutoReleaseEnabled:              true,
		DisputeResolutionPeriodHours:    DefaultDisputeResolutionPeriodHours,
	}
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyModuleEnabled, &p.ModuleEnabled, validateBool),
		paramtypes.NewParamSetPair(KeyFees, &p.Fees, validateFees),
		paramtypes.NewParamSetPair(KeyMinLcAmount, &p.MinLcAmount, validateMinLcAmount),
		paramtypes.NewParamSetPair(KeyMaxLcDurationDays, &p.MaxLcDurationDays, validateMaxLcDurationDays),
		paramtypes.NewParamSetPair(KeyDocumentSubmissionDeadlineHours, &p.DocumentSubmissionDeadlineHours, validateDocumentDeadline),
		paramtypes.NewParamSetPair(KeySupportedCurrencies, &p.SupportedCurrencies, validateSupportedCurrencies),
		paramtypes.NewParamSetPair(KeySupportedDocumentTypes, &p.SupportedDocumentTypes, validateSupportedDocumentTypes),
		paramtypes.NewParamSetPair(KeyCollateralRatio, &p.CollateralRatio, validateCollateralRatio),
		paramtypes.NewParamSetPair(KeyAutoReleaseEnabled, &p.AutoReleaseEnabled, validateBool),
		paramtypes.NewParamSetPair(KeyDisputeResolutionPeriodHours, &p.DisputeResolutionPeriodHours, validateDisputePeriod),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateBool(p.ModuleEnabled); err != nil {
		return err
	}
	if err := validateFees(p.Fees); err != nil {
		return err
	}
	if err := validateMinLcAmount(p.MinLcAmount); err != nil {
		return err
	}
	if err := validateMaxLcDurationDays(p.MaxLcDurationDays); err != nil {
		return err
	}
	if err := validateDocumentDeadline(p.DocumentSubmissionDeadlineHours); err != nil {
		return err
	}
	if err := validateSupportedCurrencies(p.SupportedCurrencies); err != nil {
		return err
	}
	if err := validateSupportedDocumentTypes(p.SupportedDocumentTypes); err != nil {
		return err
	}
	if err := validateCollateralRatio(p.CollateralRatio); err != nil {
		return err
	}
	if err := validateBool(p.AutoReleaseEnabled); err != nil {
		return err
	}
	if err := validateDisputePeriod(p.DisputeResolutionPeriodHours); err != nil {
		return err
	}

	return nil
}

// String implements the Stringer interface
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// Validation functions

func validateBool(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateFees(i interface{}) error {
	fees, ok := i.(TradeFinanceFees)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if fees.LcIssuanceFee > 1000 { // Max 10%
		return fmt.Errorf("LC issuance fee cannot exceed 10%%")
	}
	if fees.LcAmendmentFee > 500 { // Max 5%
		return fmt.Errorf("LC amendment fee cannot exceed 5%%")
	}
	if fees.DocumentVerificationFee > 10000 { // Max 10000 DINR
		return fmt.Errorf("document verification fee cannot exceed 10000 DINR")
	}

	return nil
}

func validateMinLcAmount(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("minimum LC amount must be greater than 0")
	}

	return nil
}

func validateMaxLcDurationDays(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("maximum LC duration must be greater than 0")
	}
	if v > 730 { // Max 2 years
		return fmt.Errorf("maximum LC duration cannot exceed 730 days")
	}

	return nil
}

func validateDocumentDeadline(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("document submission deadline must be greater than 0")
	}
	if v > 720 { // Max 30 days
		return fmt.Errorf("document submission deadline cannot exceed 720 hours")
	}

	return nil
}

func validateSupportedCurrencies(i interface{}) error {
	currencies, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if len(currencies) == 0 {
		return fmt.Errorf("at least one supported currency must be specified")
	}

	return nil
}

func validateSupportedDocumentTypes(i interface{}) error {
	docTypes, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if len(docTypes) == 0 {
		return fmt.Errorf("at least one supported document type must be specified")
	}

	return nil
}

func validateCollateralRatio(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 10000 { // Minimum 100%
		return fmt.Errorf("collateral ratio cannot be less than 100%%")
	}
	if v > 20000 { // Maximum 200%
		return fmt.Errorf("collateral ratio cannot exceed 200%%")
	}

	return nil
}

func validateDisputePeriod(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("dispute resolution period must be greater than 0")
	}
	if v > 720 { // Max 30 days
		return fmt.Errorf("dispute resolution period cannot exceed 720 hours")
	}

	return nil
}