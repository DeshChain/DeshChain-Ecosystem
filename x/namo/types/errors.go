package types

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/namo module sentinel errors
var (
	ErrInvalidTokenDenom       = sdkerrors.Register(ModuleName, 1, "invalid token denomination")
	ErrInvalidAmount           = sdkerrors.Register(ModuleName, 2, "invalid amount")
	ErrInsufficientBalance     = sdkerrors.Register(ModuleName, 3, "insufficient balance")
	ErrBurningDisabled         = sdkerrors.Register(ModuleName, 4, "token burning is disabled")
	ErrVestingDisabled         = sdkerrors.Register(ModuleName, 5, "vesting is disabled")
	ErrInvalidVestingSchedule  = sdkerrors.Register(ModuleName, 6, "invalid vesting schedule")
	ErrVestingScheduleNotFound = sdkerrors.Register(ModuleName, 7, "vesting schedule not found")
	ErrNoClaimableTokens       = sdkerrors.Register(ModuleName, 8, "no claimable tokens available")
	ErrInvalidRecipient        = sdkerrors.Register(ModuleName, 9, "invalid recipient address")
	ErrInvalidAuthority        = sdkerrors.Register(ModuleName, 10, "invalid authority")
	ErrMinBurnAmountNotMet     = sdkerrors.Register(ModuleName, 11, "amount is below minimum burn amount")
	ErrVestingPeriodInvalid    = sdkerrors.Register(ModuleName, 12, "vesting period is invalid")
	ErrCliffPeriodInvalid      = sdkerrors.Register(ModuleName, 13, "cliff period is invalid")
	ErrScheduleAlreadyExists   = sdkerrors.Register(ModuleName, 14, "vesting schedule already exists")
	ErrUnauthorized            = sdkerrors.Register(ModuleName, 15, "unauthorized")
	ErrInvalidParams           = sdkerrors.Register(ModuleName, 16, "invalid parameters")
	ErrTokenSupplyExhausted    = sdkerrors.Register(ModuleName, 17, "token supply exhausted")
	ErrAllocationExceeded      = sdkerrors.Register(ModuleName, 18, "allocation exceeded")
	ErrInvalidDistributionType = sdkerrors.Register(ModuleName, 19, "invalid distribution type")
	ErrDistributionNotAllowed  = sdkerrors.Register(ModuleName, 20, "distribution not allowed")
)