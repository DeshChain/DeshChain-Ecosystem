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

// Money Order module errors
var (
	// General errors
	ErrInvalidPoolType        = sdkerrors.Register(ModuleName, 1100, "invalid pool type")
	ErrPoolNotFound           = sdkerrors.Register(ModuleName, 1101, "pool not found")
	ErrInsufficientLiquidity  = sdkerrors.Register(ModuleName, 1102, "insufficient liquidity")
	ErrInvalidSwapAmount      = sdkerrors.Register(ModuleName, 1103, "invalid swap amount")
	ErrSlippageExceeded       = sdkerrors.Register(ModuleName, 1104, "slippage tolerance exceeded")
	
	// Money Order specific errors
	ErrInvalidMoneyOrder      = sdkerrors.Register(ModuleName, 1200, "invalid money order")
	ErrOrderAmountTooLow      = sdkerrors.Register(ModuleName, 1201, "order amount below minimum")
	ErrOrderAmountTooHigh     = sdkerrors.Register(ModuleName, 1202, "order amount exceeds maximum")
	ErrDailyLimitExceeded     = sdkerrors.Register(ModuleName, 1203, "daily limit exceeded")
	ErrMonthlyLimitExceeded   = sdkerrors.Register(ModuleName, 1204, "monthly limit exceeded")
	ErrInvalidReceiverUPI     = sdkerrors.Register(ModuleName, 1205, "invalid receiver UPI address")
	ErrOrderExpired           = sdkerrors.Register(ModuleName, 1206, "money order has expired")
	ErrInvalidPostalCode      = sdkerrors.Register(ModuleName, 1207, "invalid or unsupported postal code")
	
	// KYC errors
	ErrKYCRequired            = sdkerrors.Register(ModuleName, 1300, "KYC verification required")
	ErrKYCNotCompleted        = sdkerrors.Register(ModuleName, 1301, "KYC verification not completed")
	ErrKYCExpired             = sdkerrors.Register(ModuleName, 1302, "KYC verification has expired")
	ErrKYCRejected            = sdkerrors.Register(ModuleName, 1303, "KYC verification was rejected")
	
	// Fixed rate pool errors
	ErrInvalidExchangeRate    = sdkerrors.Register(ModuleName, 1400, "invalid exchange rate")
	ErrPoolInactive           = sdkerrors.Register(ModuleName, 1401, "pool is inactive")
	ErrPoolMaintenance        = sdkerrors.Register(ModuleName, 1402, "pool is under maintenance")
	ErrInsufficientPoolFunds  = sdkerrors.Register(ModuleName, 1403, "insufficient funds in pool")
	ErrUnsupportedRegion      = sdkerrors.Register(ModuleName, 1404, "region not supported by pool")
	
	// Village pool errors
	ErrVillagePoolNotFound    = sdkerrors.Register(ModuleName, 1500, "village pool not found")
	ErrNotVillageMember       = sdkerrors.Register(ModuleName, 1501, "not a village pool member")
	ErrInsufficientVillageSignatures = sdkerrors.Register(ModuleName, 1502, "insufficient village validator signatures")
	ErrVillagePoolInactive    = sdkerrors.Register(ModuleName, 1503, "village pool is inactive")
	ErrVillageNotVerified     = sdkerrors.Register(ModuleName, 1504, "village pool not government verified")
	ErrMembershipLimitReached = sdkerrors.Register(ModuleName, 1505, "village pool membership limit reached")
	
	// Liquidity errors
	ErrInsufficientShares     = sdkerrors.Register(ModuleName, 1600, "insufficient liquidity shares")
	ErrMinSharesNotMet        = sdkerrors.Register(ModuleName, 1601, "minimum shares requirement not met")
	ErrImbalancedLiquidity    = sdkerrors.Register(ModuleName, 1602, "liquidity amounts are imbalanced")
	ErrLiquidityLocked        = sdkerrors.Register(ModuleName, 1603, "liquidity is locked")
	
	// Cultural feature errors
	ErrCulturalFeaturesDisabled = sdkerrors.Register(ModuleName, 1700, "cultural features are disabled")
	ErrInvalidFestivalPeriod    = sdkerrors.Register(ModuleName, 1701, "not in festival period")
	ErrCulturalTokenNotSupported = sdkerrors.Register(ModuleName, 1702, "cultural token not supported")
	
	// Permission errors
	ErrUnauthorized           = sdkerrors.Register(ModuleName, 1800, "unauthorized")
	ErrNotPanchayatHead       = sdkerrors.Register(ModuleName, 1801, "not panchayat head")
	ErrNotLocalValidator      = sdkerrors.Register(ModuleName, 1802, "not a local validator")
	
	// Parameter errors
	ErrInvalidParams          = sdkerrors.Register(ModuleName, 1900, "invalid parameters")
	ErrFeeTooHigh             = sdkerrors.Register(ModuleName, 1901, "fee exceeds maximum allowed")
	ErrDiscountTooHigh        = sdkerrors.Register(ModuleName, 1902, "discount exceeds maximum allowed")
	
	// AMM Pool errors
	ErrInvalidPoolAssets      = sdkerrors.Register(ModuleName, 2000, "invalid pool assets")
	ErrInvalidSwapFee         = sdkerrors.Register(ModuleName, 2001, "invalid swap fee")
	ErrInsufficientOutput     = sdkerrors.Register(ModuleName, 2002, "insufficient output amount")
	ErrInvalidTokenPair       = sdkerrors.Register(ModuleName, 2003, "invalid token pair")
	ErrPoolNotActive          = sdkerrors.Register(ModuleName, 2004, "pool not active")
	ErrInvalidPoolId          = sdkerrors.Register(ModuleName, 2005, "invalid pool id")
	ErrInvalidShares          = sdkerrors.Register(ModuleName, 2006, "invalid shares amount")
	
	// General errors
	ErrNotFound               = sdkerrors.Register(ModuleName, 2100, "not found")
	
	// Biometric authentication errors
	ErrInvalidBiometric       = sdkerrors.Register(ModuleName, 2200, "invalid biometric data")
	ErrBiometricNotFound      = sdkerrors.Register(ModuleName, 2201, "biometric registration not found")
	ErrBiometricAuthFailed    = sdkerrors.Register(ModuleName, 2202, "biometric authentication failed")
	ErrInvalidDevice          = sdkerrors.Register(ModuleName, 2203, "invalid device")
	ErrBiometricLocked        = sdkerrors.Register(ModuleName, 2204, "biometric authentication locked")
	ErrInvalidAuthScore       = sdkerrors.Register(ModuleName, 2205, "invalid authentication score")
	ErrTooManyFailedAttempts  = sdkerrors.Register(ModuleName, 2206, "too many failed authentication attempts")
	
	// Order status errors
	ErrOrderNotFound          = sdkerrors.Register(ModuleName, 2300, "money order not found")
	ErrInvalidOrderStatus     = sdkerrors.Register(ModuleName, 2301, "invalid order status")
	ErrOrderAlreadyExecuted   = sdkerrors.Register(ModuleName, 2302, "money order already executed")
	ErrInvalidRecipient       = sdkerrors.Register(ModuleName, 2303, "invalid recipient")
	ErrInvalidSender          = sdkerrors.Register(ModuleName, 2304, "invalid sender")
	ErrInvalidAddress         = sdkerrors.Register(ModuleName, 2305, "invalid address")
	ErrInvalidAmount          = sdkerrors.Register(ModuleName, 2306, "invalid amount")
	ErrInsufficientFunds      = sdkerrors.Register(ModuleName, 2307, "insufficient funds")
	
	// User management errors
	ErrUserNotFound           = sdkerrors.Register(ModuleName, 2400, "user not found")
	ErrUserNotActive          = sdkerrors.Register(ModuleName, 2401, "user not active")
	ErrInvalidUserStatus      = sdkerrors.Register(ModuleName, 2402, "invalid user status")
	ErrUserBlacklisted        = sdkerrors.Register(ModuleName, 2403, "user is blacklisted")
	ErrInvalidTrustScore      = sdkerrors.Register(ModuleName, 2404, "invalid trust score")
)