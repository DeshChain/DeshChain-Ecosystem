package types

import (
	"cosmossdk.io/collections"
)

// GetBorrowerKey returns the key for a borrower
func GetBorrowerKey(borrowerID string) []byte {
	return append(BorrowerPrefix.Bytes(), []byte(borrowerID)...)
}

// GetLoanKey returns the key for a loan
func GetLoanKey(loanID string) []byte {
	return append(LoanPrefix.Bytes(), []byte(loanID)...)
}

// GetLoanSchemeKey returns the key for a loan scheme
func GetLoanSchemeKey(schemeID string) []byte {
	return append(LoanSchemePrefix.Bytes(), []byte(schemeID)...)
}

// GetLoanApplicationKey returns the key for a loan application
func GetLoanApplicationKey(applicationID string) []byte {
	return append(LoanApplicationPrefix.Bytes(), []byte(applicationID)...)
}

// GetRepaymentKey returns the key for a repayment
func GetRepaymentKey(repaymentID string) []byte {
	return append(RepaymentPrefix.Bytes(), []byte(repaymentID)...)
}

// GetCollateralKey returns the key for collateral
func GetCollateralKey(collateralID string) []byte {
	return append(CollateralPrefix.Bytes(), []byte(collateralID)...)
}

// GetVillagePoolKey returns the key for a village pool
func GetVillagePoolKey(poolID string) []byte {
	return append(VillagePoolPrefix.Bytes(), []byte(poolID)...)
}

// GetCreditHistoryKey returns the key for credit history
func GetCreditHistoryKey(borrowerID string) []byte {
	return append(CreditHistoryPrefix.Bytes(), []byte(borrowerID)...)
}

// GetRiskAssessmentKey returns the key for risk assessment
func GetRiskAssessmentKey(borrowerID string) []byte {
	return append(RiskAssessmentPrefix.Bytes(), []byte(borrowerID)...)
}

// GetCropInsuranceKey returns the key for crop insurance
func GetCropInsuranceKey(insuranceID string) []byte {
	return append(CropInsurancePrefix.Bytes(), []byte(insuranceID)...)
}

// GetWeatherDataKey returns the key for weather data
func GetWeatherDataKey(locationID string) []byte {
	return append(WeatherDataPrefix.Bytes(), []byte(locationID)...)
}

// GetCropCycleKey returns the key for crop cycle
func GetCropCycleKey(cycleID string) []byte {
	return append(CropCyclePrefix.Bytes(), []byte(cycleID)...)
}

// GetMarketPriceKey returns the key for market price
func GetMarketPriceKey(cropType string) []byte {
	return append(MarketPricePrefix.Bytes(), []byte(cropType)...)
}

// GetCommunityValidationKey returns the key for community validation
func GetCommunityValidationKey(validationID string) []byte {
	return append(CommunityValidationPrefix.Bytes(), []byte(validationID)...)
}

// GetLiquidityPoolKey returns the key for liquidity pool
func GetLiquidityPoolKey(poolID string) []byte {
	return append(LiquidityPoolPrefix.Bytes(), []byte(poolID)...)
}

// GetInterestRateKey returns the key for interest rate
func GetInterestRateKey(rateID string) []byte {
	return append(InterestRatePrefix.Bytes(), []byte(rateID)...)
}

// GetLoanPerformanceKey returns the key for loan performance
func GetLoanPerformanceKey(loanID string) []byte {
	return append(LoanPerformancePrefix.Bytes(), []byte(loanID)...)
}

// GetVillageCoordinatorKey returns the key for village coordinator
func GetVillageCoordinatorKey(coordinatorID string) []byte {
	return append(VillageCoordinatorPrefix.Bytes(), []byte(coordinatorID)...)
}