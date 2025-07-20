package types

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/liquiditymanager module sentinel errors
var (
	ErrLoanRejected             = sdkerrors.Register(ModuleName, 1100, "loan request rejected")
	ErrCollateralLoanRejected   = sdkerrors.Register(ModuleName, 1101, "collateral loan request rejected")
	ErrCollateralLockFailed     = sdkerrors.Register(ModuleName, 1102, "failed to lock collateral")
	ErrCollateralUnlockFailed   = sdkerrors.Register(ModuleName, 1103, "failed to unlock collateral")
	ErrInsufficientLiquidity    = sdkerrors.Register(ModuleName, 1104, "insufficient liquidity for lending")
	ErrNotPoolMember           = sdkerrors.Register(ModuleName, 1105, "user is not a pool member")
	ErrInsufficientCollateral  = sdkerrors.Register(ModuleName, 1106, "insufficient collateral for loan")
	ErrCollateralAlreadyLocked = sdkerrors.Register(ModuleName, 1107, "collateral is already locked")
	ErrInvalidLoanAmount       = sdkerrors.Register(ModuleName, 1108, "invalid loan amount")
	ErrDailyLimitExceeded      = sdkerrors.Register(ModuleName, 1109, "daily lending limit exceeded")
	ErrLendingNotAvailable     = sdkerrors.Register(ModuleName, 1110, "lending is not available")
	ErrModuleNotAvailable      = sdkerrors.Register(ModuleName, 1111, "lending module not available at current liquidity level")
)