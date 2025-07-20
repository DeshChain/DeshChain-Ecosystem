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
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Module errors
var (
	ErrInvalidLoanScheme         = sdkerrors.Register(ModuleName, 2001, "invalid loan scheme")
	ErrInvalidBorrower           = sdkerrors.Register(ModuleName, 2002, "invalid borrower information")
	ErrInvalidLoanAmount         = sdkerrors.Register(ModuleName, 2003, "invalid loan amount")
	ErrInvalidInterestRate       = sdkerrors.Register(ModuleName, 2004, "invalid interest rate")
	ErrInvalidLoanTerm           = sdkerrors.Register(ModuleName, 2005, "invalid loan term")
	ErrInvalidCollateral         = sdkerrors.Register(ModuleName, 2006, "invalid collateral")
	ErrInvalidCreditScore        = sdkerrors.Register(ModuleName, 2007, "invalid credit score")
	ErrInsufficientLiquidity     = sdkerrors.Register(ModuleName, 2008, "insufficient liquidity")
	ErrLoanNotFound              = sdkerrors.Register(ModuleName, 2009, "loan not found")
	ErrBorrowerNotFound          = sdkerrors.Register(ModuleName, 2010, "borrower not found")
	ErrUnauthorizedAccess        = sdkerrors.Register(ModuleName, 2011, "unauthorized access")
	ErrValidationFailed          = sdkerrors.Register(ModuleName, 2012, "validation failed")
	ErrRiskTooHigh               = sdkerrors.Register(ModuleName, 2013, "risk level too high")
	ErrLoanAlreadyExists         = sdkerrors.Register(ModuleName, 2014, "loan already exists")
	ErrRepaymentFailed           = sdkerrors.Register(ModuleName, 2015, "repayment failed")
	ErrCollateralInsufficient    = sdkerrors.Register(ModuleName, 2016, "collateral insufficient")
	ErrWeatherClaim              = sdkerrors.Register(ModuleName, 2017, "weather insurance claim")
	ErrMarketVolatility          = sdkerrors.Register(ModuleName, 2018, "market price volatility")
	ErrCommunityRejection        = sdkerrors.Register(ModuleName, 2019, "community validation rejected")
	ErrComplianceViolation       = sdkerrors.Register(ModuleName, 2020, "compliance violation")
	
	// Borrower related errors
	ErrBorrowerAlreadyRegistered = sdkerrors.Register(ModuleName, 2021, "borrower already registered")
	ErrBorrowerInactive          = sdkerrors.Register(ModuleName, 2022, "borrower is inactive")
	ErrBorrowerKYCPending        = sdkerrors.Register(ModuleName, 2023, "borrower KYC pending")
	ErrBorrowerBlacklisted       = sdkerrors.Register(ModuleName, 2024, "borrower is blacklisted")
	ErrBorrowerMaxLoansReached   = sdkerrors.Register(ModuleName, 2025, "borrower has reached maximum loans")
	
	// Application related errors
	ErrApplicationNotFound       = sdkerrors.Register(ModuleName, 2026, "loan application not found")
	ErrApplicationAlreadyProcessed = sdkerrors.Register(ModuleName, 2027, "application already processed")
	ErrApplicationExpired        = sdkerrors.Register(ModuleName, 2028, "application has expired")
	ErrInvalidApplicationStatus  = sdkerrors.Register(ModuleName, 2029, "invalid application status")
	ErrApplicationUnderReview    = sdkerrors.Register(ModuleName, 2030, "application is under review")
	
	// Loan related errors
	ErrLoanNotActive             = sdkerrors.Register(ModuleName, 2031, "loan is not active")
	ErrLoanAlreadyRepaid         = sdkerrors.Register(ModuleName, 2032, "loan already fully repaid")
	ErrLoanDefaulted             = sdkerrors.Register(ModuleName, 2033, "loan is in default")
	ErrLoanMatured               = sdkerrors.Register(ModuleName, 2034, "loan has matured")
	ErrInvalidRepaymentAmount    = sdkerrors.Register(ModuleName, 2035, "invalid repayment amount")
	ErrEarlyRepaymentNotAllowed  = sdkerrors.Register(ModuleName, 2036, "early repayment not allowed")
	ErrPartialRepaymentNotAllowed = sdkerrors.Register(ModuleName, 2037, "partial repayment not allowed")
	
	// Village pool related errors
	ErrVillagePoolNotFound       = sdkerrors.Register(ModuleName, 2038, "village pool not found")
	ErrVillagePoolInactive       = sdkerrors.Register(ModuleName, 2039, "village pool is inactive")
	ErrInsufficientPoolBalance   = sdkerrors.Register(ModuleName, 2040, "insufficient pool balance")
	ErrPoolExposureLimitReached  = sdkerrors.Register(ModuleName, 2041, "pool exposure limit reached")
	ErrInvalidCoordinator        = sdkerrors.Register(ModuleName, 2042, "invalid village coordinator")
	
	// Community validation errors
	ErrInsufficientValidators    = sdkerrors.Register(ModuleName, 2043, "insufficient community validators")
	ErrValidatorNotEligible      = sdkerrors.Register(ModuleName, 2044, "validator not eligible")
	ErrValidationTimeout         = sdkerrors.Register(ModuleName, 2045, "validation timeout")
	ErrConflictingValidations    = sdkerrors.Register(ModuleName, 2046, "conflicting validation results")
	ErrSelfValidationNotAllowed  = sdkerrors.Register(ModuleName, 2047, "self validation not allowed")
	
	// Risk assessment errors
	ErrRiskAssessmentFailed      = sdkerrors.Register(ModuleName, 2048, "risk assessment failed")
	ErrInvalidRiskFactors        = sdkerrors.Register(ModuleName, 2049, "invalid risk factors")
	ErrRiskModelNotFound         = sdkerrors.Register(ModuleName, 2050, "risk model not found")
	ErrHighRiskBorrower          = sdkerrors.Register(ModuleName, 2051, "borrower is high risk")
	
	// Insurance related errors
	ErrInsuranceRequired         = sdkerrors.Register(ModuleName, 2052, "insurance is required")
	ErrInsuranceNotFound         = sdkerrors.Register(ModuleName, 2053, "insurance policy not found")
	ErrInsuranceExpired          = sdkerrors.Register(ModuleName, 2054, "insurance policy expired")
	ErrInsuranceClaimRejected    = sdkerrors.Register(ModuleName, 2055, "insurance claim rejected")
	ErrInvalidInsurancePremium   = sdkerrors.Register(ModuleName, 2056, "invalid insurance premium")
	
	// Weather and crop errors
	ErrInvalidCropType           = sdkerrors.Register(ModuleName, 2057, "invalid crop type")
	ErrInvalidCropSeason         = sdkerrors.Register(ModuleName, 2058, "invalid crop season")
	ErrWeatherDataUnavailable    = sdkerrors.Register(ModuleName, 2059, "weather data unavailable")
	ErrCropFailure               = sdkerrors.Register(ModuleName, 2060, "crop failure detected")
	ErrHarvestDelayed            = sdkerrors.Register(ModuleName, 2061, "harvest delayed")
	
	// Market related errors
	ErrMarketDataUnavailable     = sdkerrors.Register(ModuleName, 2062, "market data unavailable")
	ErrPriceVolatilityHigh       = sdkerrors.Register(ModuleName, 2063, "price volatility too high")
	ErrMarketConditionsUnfavorable = sdkerrors.Register(ModuleName, 2064, "market conditions unfavorable")
	ErrInvalidMarketPrice        = sdkerrors.Register(ModuleName, 2065, "invalid market price")
	
	// Collateral related errors
	ErrCollateralNotFound        = sdkerrors.Register(ModuleName, 2066, "collateral not found")
	ErrCollateralAlreadyPledged  = sdkerrors.Register(ModuleName, 2067, "collateral already pledged")
	ErrCollateralValueInsufficient = sdkerrors.Register(ModuleName, 2068, "collateral value insufficient")
	ErrCollateralLiquidationFailed = sdkerrors.Register(ModuleName, 2069, "collateral liquidation failed")
	ErrInvalidCollateralType     = sdkerrors.Register(ModuleName, 2070, "invalid collateral type")
	
	// Liquidity and financial errors
	ErrLiquidityPoolNotFound     = sdkerrors.Register(ModuleName, 2071, "liquidity pool not found")
	ErrInsufficientFunds         = sdkerrors.Register(ModuleName, 2072, "insufficient funds")
	ErrInvalidTransactionAmount  = sdkerrors.Register(ModuleName, 2073, "invalid transaction amount")
	ErrTransactionFailed         = sdkerrors.Register(ModuleName, 2074, "transaction failed")
	ErrFundsLocked               = sdkerrors.Register(ModuleName, 2075, "funds are locked")
	
	// Scheme related errors
	ErrSchemeNotFound            = sdkerrors.Register(ModuleName, 2076, "loan scheme not found")
	ErrSchemeInactive            = sdkerrors.Register(ModuleName, 2077, "loan scheme is inactive")
	ErrSchemeExpired             = sdkerrors.Register(ModuleName, 2078, "loan scheme has expired")
	ErrBorrowerNotEligible       = sdkerrors.Register(ModuleName, 2079, "borrower not eligible for scheme")
	ErrAmountExceedsLimit        = sdkerrors.Register(ModuleName, 2080, "amount exceeds scheme limit")
	
	// Document and verification errors
	ErrInvalidDocuments          = sdkerrors.Register(ModuleName, 2081, "invalid documents")
	ErrDocumentVerificationFailed = sdkerrors.Register(ModuleName, 2082, "document verification failed")
	ErrMissingRequiredDocuments  = sdkerrors.Register(ModuleName, 2083, "missing required documents")
	ErrDocumentExpired           = sdkerrors.Register(ModuleName, 2084, "document has expired")
	
	// System and operational errors
	ErrSystemMaintenanceMode     = sdkerrors.Register(ModuleName, 2085, "system in maintenance mode")
	ErrRateLimitExceeded         = sdkerrors.Register(ModuleName, 2086, "rate limit exceeded")
	ErrServiceUnavailable        = sdkerrors.Register(ModuleName, 2087, "service temporarily unavailable")
	ErrInternalError             = sdkerrors.Register(ModuleName, 2088, "internal system error")
	ErrConfigurationError        = sdkerrors.Register(ModuleName, 2089, "configuration error")
	
	// Governance and permission errors
	ErrUnauthorizedOperation     = sdkerrors.Register(ModuleName, 2090, "unauthorized operation")
	ErrInsufficientPermissions   = sdkerrors.Register(ModuleName, 2091, "insufficient permissions")
	ErrInvalidAuthority          = sdkerrors.Register(ModuleName, 2092, "invalid authority")
	ErrOperationNotAllowed       = sdkerrors.Register(ModuleName, 2093, "operation not allowed")
	ErrGovernanceRestriction     = sdkerrors.Register(ModuleName, 2094, "governance restriction applies")
	
	// Integration errors
	ErrMoneyOrderIntegrationFailed = sdkerrors.Register(ModuleName, 2095, "money order integration failed")
	ErrPensionIntegrationFailed    = sdkerrors.Register(ModuleName, 2096, "pension integration failed")
	ErrBankingIntegrationFailed    = sdkerrors.Register(ModuleName, 2097, "banking integration failed")
	ErrExternalAPIFailure          = sdkerrors.Register(ModuleName, 2098, "external API failure")
	ErrDataSyncFailed              = sdkerrors.Register(ModuleName, 2099, "data synchronization failed")
	
	// Performance and limit errors
	ErrTooManyRequests           = sdkerrors.Register(ModuleName, 2100, "too many requests")
	ErrQueryLimitExceeded        = sdkerrors.Register(ModuleName, 2101, "query limit exceeded")
	ErrTimeoutError              = sdkerrors.Register(ModuleName, 2102, "operation timeout")
	ErrResourceLimitReached      = sdkerrors.Register(ModuleName, 2103, "resource limit reached")
)