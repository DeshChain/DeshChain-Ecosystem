package types

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Message types for referral system

// MsgCreateReferral creates a new validator referral
type MsgCreateReferral struct {
    Referrer     string `json:"referrer"`      // Genesis validator address
    Referred     string `json:"referred"`      // New validator address
    ReferredRank uint32 `json:"referred_rank"` // Rank for new validator
}

// NewMsgCreateReferral creates a new MsgCreateReferral instance
func NewMsgCreateReferral(referrer, referred string, rank uint32) *MsgCreateReferral {
    return &MsgCreateReferral{
        Referrer:     referrer,
        Referred:     referred,
        ReferredRank: rank,
    }
}

// Route returns the name of the module
func (msg MsgCreateReferral) Route() string { return RouterKey }

// Type returns the action
func (msg MsgCreateReferral) Type() string { return "create_referral" }

// ValidateBasic runs stateless checks on the message
func (msg MsgCreateReferral) ValidateBasic() error {
    if _, err := sdk.AccAddressFromBech32(msg.Referrer); err != nil {
        return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid referrer address: %v", err)
    }
    if _, err := sdk.AccAddressFromBech32(msg.Referred); err != nil {
        return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid referred address: %v", err)
    }
    if msg.ReferredRank == 0 || msg.ReferredRank > 1000 {
        return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid rank: %d (must be 1-1000)", msg.ReferredRank)
    }
    if msg.Referrer == msg.Referred {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "cannot refer yourself")
    }
    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgCreateReferral) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners defines whose signature is required
func (msg MsgCreateReferral) GetSigners() []sdk.AccAddress {
    addr, _ := sdk.AccAddressFromBech32(msg.Referrer)
    return []sdk.AccAddress{addr}
}

// MsgLaunchValidatorToken launches a validator token
type MsgLaunchValidatorToken struct {
    Validator string `json:"validator"` // Validator address
}

// NewMsgLaunchValidatorToken creates a new MsgLaunchValidatorToken instance
func NewMsgLaunchValidatorToken(validator string) *MsgLaunchValidatorToken {
    return &MsgLaunchValidatorToken{
        Validator: validator,
    }
}

// Route returns the name of the module
func (msg MsgLaunchValidatorToken) Route() string { return RouterKey }

// Type returns the action
func (msg MsgLaunchValidatorToken) Type() string { return "launch_validator_token" }

// ValidateBasic runs stateless checks on the message
func (msg MsgLaunchValidatorToken) ValidateBasic() error {
    if _, err := sdk.AccAddressFromBech32(msg.Validator); err != nil {
        return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid validator address: %v", err)
    }
    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgLaunchValidatorToken) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners defines whose signature is required
func (msg MsgLaunchValidatorToken) GetSigners() []sdk.AccAddress {
    addr, _ := sdk.AccAddressFromBech32(msg.Validator)
    return []sdk.AccAddress{addr}
}

// MsgAirdropTokens airdrops validator tokens to addresses
type MsgAirdropTokens struct {
    Validator   string              `json:"validator"`   // Validator address
    TokenID     uint64              `json:"token_id"`    // Token ID
    Recipients  []AirdropRecipient  `json:"recipients"`  // Recipients list
}

// AirdropRecipient represents an airdrop recipient
type AirdropRecipient struct {
    Address string  `json:"address"`
    Amount  sdk.Int `json:"amount"`
}

// NewMsgAirdropTokens creates a new MsgAirdropTokens instance
func NewMsgAirdropTokens(validator string, tokenID uint64, recipients []AirdropRecipient) *MsgAirdropTokens {
    return &MsgAirdropTokens{
        Validator:  validator,
        TokenID:    tokenID,
        Recipients: recipients,
    }
}

// Route returns the name of the module
func (msg MsgAirdropTokens) Route() string { return RouterKey }

// Type returns the action
func (msg MsgAirdropTokens) Type() string { return "airdrop_tokens" }

// ValidateBasic runs stateless checks on the message
func (msg MsgAirdropTokens) ValidateBasic() error {
    if _, err := sdk.AccAddressFromBech32(msg.Validator); err != nil {
        return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid validator address: %v", err)
    }
    if msg.TokenID == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "token ID cannot be zero")
    }
    if len(msg.Recipients) == 0 {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "recipients list cannot be empty")
    }
    if len(msg.Recipients) > 1000 {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "too many recipients (max 1000)")
    }
    
    for i, recipient := range msg.Recipients {
        if _, err := sdk.AccAddressFromBech32(recipient.Address); err != nil {
            return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid recipient address at index %d: %v", i, err)
        }
        if recipient.Amount.IsZero() || recipient.Amount.IsNegative() {
            return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid amount for recipient %d: %s", i, recipient.Amount)
        }
    }
    
    return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgAirdropTokens) GetSignBytes() []byte {
    return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners defines whose signature is required
func (msg MsgAirdropTokens) GetSigners() []sdk.AccAddress {
    addr, _ := sdk.AccAddressFromBech32(msg.Validator)
    return []sdk.AccAddress{addr}
}

// Response types

// MsgCreateReferralResponse defines the Msg/CreateReferral response type
type MsgCreateReferralResponse struct {
    ReferralID uint64 `json:"referral_id"`
}

// MsgLaunchValidatorTokenResponse defines the Msg/LaunchValidatorToken response type
type MsgLaunchValidatorTokenResponse struct {
    TokenID uint64 `json:"token_id"`
    Success bool   `json:"success"`
}

// MsgAirdropTokensResponse defines the Msg/AirdropTokens response type
type MsgAirdropTokensResponse struct {
    RecipientsCount uint32  `json:"recipients_count"`
    TotalAmount     sdk.Int `json:"total_amount"`
}