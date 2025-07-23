package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/dinr module sentinel errors
var (
	ErrMintingDisabled          = sdkerrors.Register(ModuleName, 1100, "minting is currently disabled")
	ErrBurningDisabled          = sdkerrors.Register(ModuleName, 1101, "burning is currently disabled")
	ErrInvalidCollateral        = sdkerrors.Register(ModuleName, 1102, "invalid or inactive collateral asset")
	ErrInsufficientCollateral   = sdkerrors.Register(ModuleName, 1103, "insufficient collateral ratio")
	ErrPositionNotFound         = sdkerrors.Register(ModuleName, 1104, "user position not found")
	ErrInsufficientDINR         = sdkerrors.Register(ModuleName, 1105, "insufficient DINR minted")
	ErrInsufficientBalance      = sdkerrors.Register(ModuleName, 1106, "insufficient balance")
	ErrOraclePriceNotAvailable  = sdkerrors.Register(ModuleName, 1107, "oracle price not available")
	ErrBelowMinimumMint         = sdkerrors.Register(ModuleName, 1108, "mint amount below minimum")
	ErrExceedsMaxSupply         = sdkerrors.Register(ModuleName, 1109, "exceeds maximum DINR supply")
	ErrInvalidCollateralRatio   = sdkerrors.Register(ModuleName, 1110, "invalid collateral ratio")
	ErrPositionNotLiquidatable  = sdkerrors.Register(ModuleName, 1111, "position not eligible for liquidation")
	ErrInvalidLiquidationAmount = sdkerrors.Register(ModuleName, 1112, "invalid liquidation amount")
	ErrExceedsWithdrawLimit     = sdkerrors.Register(ModuleName, 1113, "exceeds maximum withdrawable collateral")
	ErrInvalidFeeStructure      = sdkerrors.Register(ModuleName, 1114, "invalid fee structure")
	ErrYieldStrategyNotFound    = sdkerrors.Register(ModuleName, 1115, "yield strategy not found")
	ErrInsuranceFundInsufficient = sdkerrors.Register(ModuleName, 1116, "insurance fund insufficient")
	ErrExcessiveLiquidation     = sdkerrors.Register(ModuleName, 1117, "liquidation amount exceeds debt")
	ErrUnauthorized             = sdkerrors.Register(ModuleName, 1118, "unauthorized")
	ErrDailyMintingLimitExceeded = sdkerrors.Register(ModuleName, 1119, "daily minting limit exceeded")
	ErrDailyBurningLimitExceeded = sdkerrors.Register(ModuleName, 1120, "daily burning limit exceeded")
	ErrInsufficientStabilityPoolBalance = sdkerrors.Register(ModuleName, 1121, "insufficient stability pool balance")
	ErrInvalidTargetPrice        = sdkerrors.Register(ModuleName, 1122, "invalid target price")
	ErrInvalidTolerance          = sdkerrors.Register(ModuleName, 1123, "invalid tolerance")
)