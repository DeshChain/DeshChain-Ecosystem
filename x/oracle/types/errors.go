package types

import (
	sdkerrors "cosmossdk.io/errors"
)

// DONTCOVER

// x/oracle module sentinel errors
var (
	ErrInvalidSigner              = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrInvalidSymbol              = sdkerrors.Register(ModuleName, 1101, "invalid asset symbol")
	ErrInvalidPrice               = sdkerrors.Register(ModuleName, 1102, "invalid price value")
	ErrInvalidValidator           = sdkerrors.Register(ModuleName, 1103, "invalid validator address")
	ErrValidatorNotFound          = sdkerrors.Register(ModuleName, 1104, "oracle validator not found")
	ErrValidatorNotActive         = sdkerrors.Register(ModuleName, 1105, "oracle validator is not active")
	ErrUnauthorizedValidator      = sdkerrors.Register(ModuleName, 1106, "validator not authorized for oracle submissions")
	ErrPriceNotFound              = sdkerrors.Register(ModuleName, 1107, "price data not found for symbol")
	ErrExchangeRateNotFound       = sdkerrors.Register(ModuleName, 1108, "exchange rate not found")
	ErrInvalidExchangeRate        = sdkerrors.Register(ModuleName, 1109, "invalid exchange rate")
	ErrInvalidCurrency            = sdkerrors.Register(ModuleName, 1110, "invalid currency code")
	ErrInsufficientValidators     = sdkerrors.Register(ModuleName, 1111, "insufficient number of validators for price consensus")
	ErrPriceDeviationTooHigh      = sdkerrors.Register(ModuleName, 1112, "price deviation exceeds maximum allowed")
	ErrStalePrice                 = sdkerrors.Register(ModuleName, 1113, "price data is too old")
	ErrInvalidMinSources          = sdkerrors.Register(ModuleName, 1121, "invalid minimum sources")
	ErrInvalidMaxDeviation        = sdkerrors.Register(ModuleName, 1122, "invalid maximum deviation")
	ErrInvalidUpdateInterval      = sdkerrors.Register(ModuleName, 1123, "invalid update interval")
	ErrInvalidPriceExpiryTime     = sdkerrors.Register(ModuleName, 1124, "invalid price expiry time")
	ErrNoEnabledSymbols           = sdkerrors.Register(ModuleName, 1125, "no enabled symbols")
	ErrInvalidSourceName          = sdkerrors.Register(ModuleName, 1126, "invalid source name")
	ErrInvalidSourceEndpoint      = sdkerrors.Register(ModuleName, 1127, "invalid source endpoint")
	ErrInvalidSourceWeight        = sdkerrors.Register(ModuleName, 1128, "invalid source weight")
	ErrInvalidSourceTimeout       = sdkerrors.Register(ModuleName, 1129, "invalid source timeout")
	ErrPriceExpired               = sdkerrors.Register(ModuleName, 1130, "price expired")
	ErrInsufficientSources        = sdkerrors.Register(ModuleName, 1131, "insufficient oracle sources")
	ErrOracleNotHealthy           = sdkerrors.Register(ModuleName, 1132, "oracle not healthy")
	ErrDuplicateSubmission        = sdkerrors.Register(ModuleName, 1114, "duplicate price submission within aggregation window")
	ErrInvalidTimestamp           = sdkerrors.Register(ModuleName, 1115, "invalid timestamp")
	ErrInvalidSource              = sdkerrors.Register(ModuleName, 1116, "invalid price source")
	ErrValidatorSlashed           = sdkerrors.Register(ModuleName, 1117, "validator has been slashed for oracle misbehavior")
	ErrAggregationWindowNotReady  = sdkerrors.Register(ModuleName, 1118, "aggregation window is not ready for finalization")
	ErrInvalidPower               = sdkerrors.Register(ModuleName, 1119, "invalid validator power")
	ErrValidatorAlreadyExists     = sdkerrors.Register(ModuleName, 1120, "oracle validator already exists")
)