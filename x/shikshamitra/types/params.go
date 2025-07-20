package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys
var (
	KeyMinInterestRate = []byte("MinInterestRate")
	KeyMaxInterestRate = []byte("MaxInterestRate")
	KeyMeritDiscount = []byte("MeritDiscount")
	KeyWomenStudentDiscount = []byte("WomenStudentDiscount")
	KeyReservedCategoryDiscount = []byte("ReservedCategoryDiscount")
	KeyMinLoanDuration = []byte("MinLoanDuration")
	KeyMaxLoanDuration = []byte("MaxLoanDuration")
	KeyMinLoanAmount = []byte("MinLoanAmount")
	KeyMaxLoanAmount = []byte("MaxLoanAmount")
	KeyAuthorizedApprovers = []byte("AuthorizedApprovers")
	KeyAuthorizedDisbursers = []byte("AuthorizedDisbursers")
	KeyEducationVerifiers = []byte("EducationVerifiers")
	KeyCollateralThreshold = []byte("CollateralThreshold")
	KeyMoratoriumPeriod = []byte("MoratoriumPeriod")
	KeyInstitutionRankings = []byte("InstitutionRankings")
	KeyCourseInterestRates = []byte("CourseInterestRates")
)

// DefaultParams returns default module parameters
func DefaultParams() Params {
	return Params{
		MinInterestRate: "0.04", // 4%
		MaxInterestRate: "0.07", // 7%
		MeritDiscount: "0.01", // 1% discount for 80%+ marks
		WomenStudentDiscount: "0.005", // 0.5% discount
		ReservedCategoryDiscount: "0.01", // 1% discount
		MinLoanDuration: 12, // 12 months
		MaxLoanDuration: 180, // 15 years
		MinLoanAmount: sdk.NewCoin("NAMO", sdk.NewInt(25000)), // ₹25,000
		MaxLoanAmount: sdk.NewCoin("NAMO", sdk.NewInt(2000000)), // ₹20 Lakhs
		AuthorizedApprovers: []string{},
		AuthorizedDisbursers: []string{},
		EducationVerifiers: []string{},
		CollateralThreshold: "1000000", // ₹10 Lakhs
		MoratoriumPeriod: 6, // 6 months after course completion
		InstitutionRankings: map[string]string{
			"IIT": "TIER1",
			"IIM": "TIER1",
			"NIT": "TIER2",
			"AIIMS": "TIER1",
		},
		CourseInterestRates: map[string]string{
			"ENGINEERING": "0.05",
			"MEDICAL": "0.05",
			"MANAGEMENT": "0.055",
			"LAW": "0.06",
		},
	}
}

// ParamKeyTable returns the parameter key table for shikshamitra module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs returns the parameter set pairs
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMinInterestRate, &p.MinInterestRate, validateInterestRate),
		paramtypes.NewParamSetPair(KeyMaxInterestRate, &p.MaxInterestRate, validateInterestRate),
		paramtypes.NewParamSetPair(KeyMeritDiscount, &p.MeritDiscount, validateDiscount),
		paramtypes.NewParamSetPair(KeyWomenStudentDiscount, &p.WomenStudentDiscount, validateDiscount),
		paramtypes.NewParamSetPair(KeyReservedCategoryDiscount, &p.ReservedCategoryDiscount, validateDiscount),
		paramtypes.NewParamSetPair(KeyMinLoanDuration, &p.MinLoanDuration, validateDuration),
		paramtypes.NewParamSetPair(KeyMaxLoanDuration, &p.MaxLoanDuration, validateDuration),
		paramtypes.NewParamSetPair(KeyMinLoanAmount, &p.MinLoanAmount, validateLoanAmount),
		paramtypes.NewParamSetPair(KeyMaxLoanAmount, &p.MaxLoanAmount, validateLoanAmount),
		paramtypes.NewParamSetPair(KeyAuthorizedApprovers, &p.AuthorizedApprovers, validateAddressList),
		paramtypes.NewParamSetPair(KeyAuthorizedDisbursers, &p.AuthorizedDisbursers, validateAddressList),
		paramtypes.NewParamSetPair(KeyEducationVerifiers, &p.EducationVerifiers, validateAddressList),
		paramtypes.NewParamSetPair(KeyCollateralThreshold, &p.CollateralThreshold, validateThreshold),
		paramtypes.NewParamSetPair(KeyMoratoriumPeriod, &p.MoratoriumPeriod, validateDuration),
		paramtypes.NewParamSetPair(KeyInstitutionRankings, &p.InstitutionRankings, validateStringMap),
		paramtypes.NewParamSetPair(KeyCourseInterestRates, &p.CourseInterestRates, validateStringMap),
	}
}

// Validate validates all params
func (p Params) Validate() error {
	if err := validateInterestRate(p.MinInterestRate); err != nil {
		return err
	}
	if err := validateInterestRate(p.MaxInterestRate); err != nil {
		return err
	}
	if err := validateDiscount(p.MeritDiscount); err != nil {
		return err
	}
	if err := validateDiscount(p.WomenStudentDiscount); err != nil {
		return err
	}
	if err := validateDiscount(p.ReservedCategoryDiscount); err != nil {
		return err
	}
	if err := validateDuration(p.MinLoanDuration); err != nil {
		return err
	}
	if err := validateDuration(p.MaxLoanDuration); err != nil {
		return err
	}
	if err := validateLoanAmount(p.MinLoanAmount); err != nil {
		return err
	}
	if err := validateLoanAmount(p.MaxLoanAmount); err != nil {
		return err
	}
	if err := validateAddressList(p.AuthorizedApprovers); err != nil {
		return err
	}
	if err := validateAddressList(p.AuthorizedDisbursers); err != nil {
		return err
	}
	if err := validateAddressList(p.EducationVerifiers); err != nil {
		return err
	}
	if err := validateThreshold(p.CollateralThreshold); err != nil {
		return err
	}
	if err := validateDuration(p.MoratoriumPeriod); err != nil {
		return err
	}
	if err := validateStringMap(p.InstitutionRankings); err != nil {
		return err
	}
	if err := validateStringMap(p.CourseInterestRates); err != nil {
		return err
	}
	return nil
}

func validateInterestRate(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	rate, err := sdk.NewDecFromStr(v)
	if err != nil {
		return err
	}

	if rate.IsNegative() || rate.GT(sdk.OneDec()) {
		return fmt.Errorf("interest rate must be between 0 and 1")
	}

	return nil
}

func validateDiscount(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	discount, err := sdk.NewDecFromStr(v)
	if err != nil {
		return err
	}

	if discount.IsNegative() || discount.GT(sdk.NewDecWithPrec(5, 2)) {
		return fmt.Errorf("discount must be between 0 and 0.05 (5%%)")
	}

	return nil
}

func validateDuration(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("duration must be positive")
	}

	return nil
}

func validateLoanAmount(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if !v.IsValid() || v.IsZero() {
		return fmt.Errorf("loan amount must be positive")
	}

	return nil
}

func validateAddressList(i interface{}) error {
	v, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	for _, addr := range v {
		if _, err := sdk.AccAddressFromBech32(addr); err != nil {
			return fmt.Errorf("invalid address %s: %w", addr, err)
		}
	}

	return nil
}

func validateThreshold(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	threshold, err := sdk.NewDecFromStr(v)
	if err != nil {
		return err
	}

	if threshold.IsNegative() {
		return fmt.Errorf("threshold must be non-negative")
	}

	return nil
}

func validateStringMap(i interface{}) error {
	_, ok := i.(map[string]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}