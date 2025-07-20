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
	ErrInvalidSchemeID            = sdkerrors.Register(ModuleName, 1001, "invalid scheme id")
	ErrInvalidSchemeName          = sdkerrors.Register(ModuleName, 1002, "invalid scheme name")
	ErrInvalidContribution        = sdkerrors.Register(ModuleName, 1003, "invalid contribution amount")
	ErrInvalidContributionPeriod  = sdkerrors.Register(ModuleName, 1004, "invalid contribution period")
	ErrInvalidMaturityBonus       = sdkerrors.Register(ModuleName, 1005, "invalid maturity bonus")
	ErrInvalidAgeRange            = sdkerrors.Register(ModuleName, 1006, "invalid age range")
	ErrInvalidLiquidityUtilization = sdkerrors.Register(ModuleName, 1007, "invalid liquidity utilization rate")
	
	ErrInvalidParticipantID       = sdkerrors.Register(ModuleName, 1008, "invalid participant id")
	ErrInvalidAddress             = sdkerrors.Register(ModuleName, 1009, "invalid address")
	ErrInvalidAge                 = sdkerrors.Register(ModuleName, 1010, "invalid age")
	
	ErrSchemeNotFound             = sdkerrors.Register(ModuleName, 1011, "scheme not found")
	ErrParticipantNotFound        = sdkerrors.Register(ModuleName, 1012, "participant not found")
	ErrSchemeInactive             = sdkerrors.Register(ModuleName, 1013, "scheme is inactive")
	ErrSchemeAtCapacity           = sdkerrors.Register(ModuleName, 1014, "scheme is at maximum capacity")
	ErrAgeNotEligible             = sdkerrors.Register(ModuleName, 1015, "age not eligible for scheme")
	ErrAlreadyEnrolled            = sdkerrors.Register(ModuleName, 1016, "already enrolled in scheme")
	ErrNotEnrolled                = sdkerrors.Register(ModuleName, 1017, "not enrolled in scheme")
	
	ErrInsufficientBalance        = sdkerrors.Register(ModuleName, 1018, "insufficient balance")
	ErrContributionAlreadyMade    = sdkerrors.Register(ModuleName, 1019, "contribution already made for this month")
	ErrContributionPeriodExceeded = sdkerrors.Register(ModuleName, 1020, "contribution period exceeded")
	ErrInvalidContributionMonth   = sdkerrors.Register(ModuleName, 1021, "invalid contribution month")
	ErrMissedTooManyPayments      = sdkerrors.Register(ModuleName, 1022, "missed too many payments")
	
	ErrNotMatured                 = sdkerrors.Register(ModuleName, 1023, "pension not matured yet")
	ErrAlreadyMatured             = sdkerrors.Register(ModuleName, 1024, "pension already matured")
	ErrAlreadyProcessed           = sdkerrors.Register(ModuleName, 1025, "already processed")
	ErrCannotWithdrawEarly        = sdkerrors.Register(ModuleName, 1026, "cannot withdraw early")
	ErrWithdrawalPending          = sdkerrors.Register(ModuleName, 1027, "withdrawal already pending")
	
	ErrKYCNotVerified             = sdkerrors.Register(ModuleName, 1028, "KYC not verified")
	ErrKYCRequired                = sdkerrors.Register(ModuleName, 1029, "KYC verification required")
	ErrInvalidKYCStatus           = sdkerrors.Register(ModuleName, 1030, "invalid KYC status")
	
	ErrInvalidReferrer            = sdkerrors.Register(ModuleName, 1031, "invalid referrer")
	ErrSelfReferral               = sdkerrors.Register(ModuleName, 1032, "self referral not allowed")
	ErrReferrerNotEligible        = sdkerrors.Register(ModuleName, 1033, "referrer not eligible")
	
	ErrInvalidWithdrawalAmount    = sdkerrors.Register(ModuleName, 1034, "invalid withdrawal amount")
	ErrInvalidWithdrawalReason    = sdkerrors.Register(ModuleName, 1035, "invalid withdrawal reason")
	ErrWithdrawalNotAllowed       = sdkerrors.Register(ModuleName, 1036, "withdrawal not allowed")
	
	ErrInvalidLiquidityPool       = sdkerrors.Register(ModuleName, 1037, "invalid liquidity pool")
	ErrLiquidityProvisionFailed   = sdkerrors.Register(ModuleName, 1038, "liquidity provision failed")
	ErrLiquidityReturnFailed      = sdkerrors.Register(ModuleName, 1039, "liquidity return failed")
	
	ErrInvalidVillageCode         = sdkerrors.Register(ModuleName, 1040, "invalid village postal code")
	ErrVillagePoolNotFound        = sdkerrors.Register(ModuleName, 1041, "village pool not found")
	
	ErrUnauthorized               = sdkerrors.Register(ModuleName, 1042, "unauthorized")
	ErrInvalidParams              = sdkerrors.Register(ModuleName, 1043, "invalid parameters")
	ErrInternalError              = sdkerrors.Register(ModuleName, 1044, "internal error")
)