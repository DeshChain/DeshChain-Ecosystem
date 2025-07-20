package types

const (
	// ModuleName defines the module name
	ModuleName = "krishimitra"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_krishimitra"

	// Interest rate constants
	MinInterestRate = "0.06" // 6%
	MaxInterestRate = "0.09" // 9%

	// Loan limits
	MinLoanAmount = "10000"    // ₹10,000
	MaxLoanAmount = "10000000" // ₹1 Crore

	// Collateral ratio
	MinCollateralRatio = "1.2" // 120%

	// Festival bonus reduction
	FestivalInterestReduction = "0.01" // 1% reduction during festivals
)

// Key prefixes
var (
	LoanKeyPrefix          = []byte{0x01}
	ApplicationKeyPrefix   = []byte{0x02}
	CollateralKeyPrefix    = []byte{0x03}
	RepaymentKeyPrefix     = []byte{0x04}
	EligibilityKeyPrefix   = []byte{0x05}
	FestivalPeriodPrefix   = []byte{0x06}
	PINCodeEligiblePrefix  = []byte{0x07}
)

// GetLoanKey returns the store key for a loan
func GetLoanKey(loanID string) []byte {
	return append(LoanKeyPrefix, []byte(loanID)...)
}

// GetApplicationKey returns the store key for an application
func GetApplicationKey(applicationID string) []byte {
	return append(ApplicationKeyPrefix, []byte(applicationID)...)
}

// GetCollateralKey returns the store key for collateral
func GetCollateralKey(loanID string) []byte {
	return append(CollateralKeyPrefix, []byte(loanID)...)
}

// GetPINCodeKey returns the store key for PIN code eligibility
func GetPINCodeKey(pinCode string) []byte {
	return append(PINCodeEligiblePrefix, []byte(pinCode)...)
}