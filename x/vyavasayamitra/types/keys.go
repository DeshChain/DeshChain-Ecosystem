package types

const (
	// ModuleName defines the module name
	ModuleName = "vyavasayamitra"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_vyavasayamitra"

	// Interest rate constants
	MinInterestRate = "0.08" // 8%
	MaxInterestRate = "0.12" // 12%

	// Loan limits
	MinLoanAmount = "50000"     // ₹50,000
	MaxLoanAmount = "100000000" // ₹10 Crore

	// Collateral ratio
	MinCollateralRatio = "1.5" // 150%

	// Festival bonus reduction
	FestivalInterestReduction = "0.005" // 0.5% reduction during festivals

	// Business categories
	CategoryManufacturing = "manufacturing"
	CategoryRetail        = "retail"
	CategoryServices      = "services"
	CategoryTechnology    = "technology"
	CategoryExportImport  = "export_import"
)

// Key prefixes
var (
	LoanKeyPrefix           = []byte{0x01}
	ApplicationKeyPrefix    = []byte{0x02}
	CollateralKeyPrefix     = []byte{0x03}
	RepaymentKeyPrefix      = []byte{0x04}
	BusinessProfilePrefix   = []byte{0x05}
	CreditLineKeyPrefix     = []byte{0x06}
	InvoiceKeyPrefix        = []byte{0x07}
	PINCodeEligiblePrefix   = []byte{0x08}
	FestivalPeriodPrefix    = []byte{0x09}
	MerchantRatingPrefix    = []byte{0x0A}
)

// GetLoanKey returns the store key for a loan
func GetLoanKey(loanID string) []byte {
	return append(LoanKeyPrefix, []byte(loanID)...)
}

// GetApplicationKey returns the store key for an application
func GetApplicationKey(applicationID string) []byte {
	return append(ApplicationKeyPrefix, []byte(applicationID)...)
}

// GetBusinessProfileKey returns the store key for a business profile
func GetBusinessProfileKey(businessID string) []byte {
	return append(BusinessProfilePrefix, []byte(businessID)...)
}

// GetCreditLineKey returns the store key for a credit line
func GetCreditLineKey(creditLineID string) []byte {
	return append(CreditLineKeyPrefix, []byte(creditLineID)...)
}

// GetInvoiceKey returns the store key for an invoice
func GetInvoiceKey(invoiceID string) []byte {
	return append(InvoiceKeyPrefix, []byte(invoiceID)...)
}

// GetPINCodeKey returns the store key for PIN code eligibility
func GetPINCodeKey(pinCode string) []byte {
	return append(PINCodeEligiblePrefix, []byte(pinCode)...)
}

// GetMerchantRatingKey returns the store key for merchant rating
func GetMerchantRatingKey(merchantID string) []byte {
	return append(MerchantRatingPrefix, []byte(merchantID)...)
}