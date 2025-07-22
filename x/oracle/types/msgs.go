package types

import (
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "cosmossdk.io/errors"
)

const TypeMsgSubmitPrice = "submit_price"
const TypeMsgSubmitExchangeRate = "submit_exchange_rate"
const TypeMsgRegisterOracleValidator = "register_oracle_validator"
const TypeMsgUpdateOracleValidator = "update_oracle_validator"
const TypeMsgUpdateParams = "update_params"

var _ sdk.Msg = &MsgSubmitPrice{}
var _ sdk.Msg = &MsgSubmitExchangeRate{}
var _ sdk.Msg = &MsgRegisterOracleValidator{}
var _ sdk.Msg = &MsgUpdateOracleValidator{}
var _ sdk.Msg = &MsgUpdateParams{}

// NewMsgSubmitPrice creates a new MsgSubmitPrice instance
func NewMsgSubmitPrice(validator, symbol string, price sdk.Dec, source string, timestamp time.Time) *MsgSubmitPrice {
	return &MsgSubmitPrice{
		Validator: validator,
		Symbol:    symbol,
		Price:     price,
		Source:    source,
		Timestamp: timestamp,
	}
}

// Route returns the name of the module
func (msg MsgSubmitPrice) Route() string { return RouterKey }

// Type returns the action type
func (msg MsgSubmitPrice) Type() string { return TypeMsgSubmitPrice }

// GetSigners returns the expected signers
func (msg *MsgSubmitPrice) GetSigners() []sdk.AccAddress {
	validator, err := sdk.AccAddressFromBech32(msg.Validator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{validator}
}

// GetSignBytes encodes the message for signing
func (msg *MsgSubmitPrice) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// ValidateBasic runs stateless checks on the message
func (msg MsgSubmitPrice) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Validator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid validator address (%s)", err)
	}

	if strings.TrimSpace(msg.Symbol) == "" {
		return sdkerrors.Wrap(ErrInvalidSymbol, "symbol cannot be empty")
	}

	if msg.Price.IsNil() || msg.Price.IsNegative() || msg.Price.IsZero() {
		return sdkerrors.Wrap(ErrInvalidPrice, "price must be positive")
	}

	if strings.TrimSpace(msg.Source) == "" {
		return sdkerrors.Wrap(ErrInvalidSource, "source cannot be empty")
	}

	if msg.Timestamp.IsZero() {
		return sdkerrors.Wrap(ErrInvalidTimestamp, "timestamp cannot be zero")
	}

	// Check if timestamp is not too far in the future (max 5 minutes)
	if msg.Timestamp.After(time.Now().Add(5 * time.Minute)) {
		return sdkerrors.Wrap(ErrInvalidTimestamp, "timestamp cannot be too far in the future")
	}

	// Check if timestamp is not too old (max 1 hour)
	if msg.Timestamp.Before(time.Now().Add(-1 * time.Hour)) {
		return sdkerrors.Wrap(ErrInvalidTimestamp, "timestamp cannot be too old")
	}

	return nil
}

// NewMsgSubmitExchangeRate creates a new MsgSubmitExchangeRate instance
func NewMsgSubmitExchangeRate(validator, base, target string, rate sdk.Dec, source string, timestamp time.Time) *MsgSubmitExchangeRate {
	return &MsgSubmitExchangeRate{
		Validator: validator,
		Base:      base,
		Target:    target,
		Rate:      rate,
		Source:    source,
		Timestamp: timestamp,
	}
}

// Route returns the name of the module
func (msg MsgSubmitExchangeRate) Route() string { return RouterKey }

// Type returns the action type
func (msg MsgSubmitExchangeRate) Type() string { return TypeMsgSubmitExchangeRate }

// GetSigners returns the expected signers
func (msg *MsgSubmitExchangeRate) GetSigners() []sdk.AccAddress {
	validator, err := sdk.AccAddressFromBech32(msg.Validator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{validator}
}

// GetSignBytes encodes the message for signing
func (msg *MsgSubmitExchangeRate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// ValidateBasic runs stateless checks on the message
func (msg MsgSubmitExchangeRate) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Validator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid validator address (%s)", err)
	}

	if strings.TrimSpace(msg.Base) == "" {
		return sdkerrors.Wrap(ErrInvalidCurrency, "base currency cannot be empty")
	}

	if strings.TrimSpace(msg.Target) == "" {
		return sdkerrors.Wrap(ErrInvalidCurrency, "target currency cannot be empty")
	}

	if msg.Base == msg.Target {
		return sdkerrors.Wrap(ErrInvalidCurrency, "base and target currencies cannot be the same")
	}

	if msg.Rate.IsNil() || msg.Rate.IsNegative() || msg.Rate.IsZero() {
		return sdkerrors.Wrap(ErrInvalidExchangeRate, "exchange rate must be positive")
	}

	if strings.TrimSpace(msg.Source) == "" {
		return sdkerrors.Wrap(ErrInvalidSource, "source cannot be empty")
	}

	if msg.Timestamp.IsZero() {
		return sdkerrors.Wrap(ErrInvalidTimestamp, "timestamp cannot be zero")
	}

	return nil
}

// NewMsgRegisterOracleValidator creates a new MsgRegisterOracleValidator instance
func NewMsgRegisterOracleValidator(creator, validator string, power uint64, description string) *MsgRegisterOracleValidator {
	return &MsgRegisterOracleValidator{
		Creator:     creator,
		Validator:   validator,
		Power:       power,
		Description: description,
	}
}

// Route returns the name of the module
func (msg MsgRegisterOracleValidator) Route() string { return RouterKey }

// Type returns the action type
func (msg MsgRegisterOracleValidator) Type() string { return TypeMsgRegisterOracleValidator }

// GetSigners returns the expected signers
func (msg *MsgRegisterOracleValidator) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes encodes the message for signing
func (msg *MsgRegisterOracleValidator) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// ValidateBasic runs stateless checks on the message
func (msg MsgRegisterOracleValidator) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Validator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid validator address (%s)", err)
	}

	if msg.Power == 0 {
		return sdkerrors.Wrap(ErrInvalidPower, "power cannot be zero")
	}

	return nil
}

// NewMsgUpdateOracleValidator creates a new MsgUpdateOracleValidator instance
func NewMsgUpdateOracleValidator(creator, validator string, power uint64, active bool) *MsgUpdateOracleValidator {
	return &MsgUpdateOracleValidator{
		Creator:   creator,
		Validator: validator,
		Power:     power,
		Active:    active,
	}
}

// Route returns the name of the module
func (msg MsgUpdateOracleValidator) Route() string { return RouterKey }

// Type returns the action type
func (msg MsgUpdateOracleValidator) Type() string { return TypeMsgUpdateOracleValidator }

// GetSigners returns the expected signers
func (msg *MsgUpdateOracleValidator) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes encodes the message for signing
func (msg *MsgUpdateOracleValidator) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// ValidateBasic runs stateless checks on the message
func (msg MsgUpdateOracleValidator) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Validator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid validator address (%s)", err)
	}

	return nil
}

// NewMsgUpdateParams creates a new MsgUpdateParams instance
func NewMsgUpdateParams(authority string, params OracleParams) *MsgUpdateParams {
	return &MsgUpdateParams{
		Authority: authority,
		Params:    params,
	}
}

// Route returns the name of the module
func (msg MsgUpdateParams) Route() string { return RouterKey }

// Type returns the action type
func (msg MsgUpdateParams) Type() string { return TypeMsgUpdateParams }

// GetSigners returns the expected signers
func (msg *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

// GetSignBytes encodes the message for signing
func (msg *MsgUpdateParams) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// ValidateBasic runs stateless checks on the message
func (msg MsgUpdateParams) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}

	return msg.Params.Validate()
}