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
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Message types
const (
	TypeMsgCreateScheme      = "create_scheme"
	TypeMsgUpdateScheme      = "update_scheme"
	TypeMsgEnrollParticipant = "enroll_participant"
	TypeMsgMakeContribution  = "make_contribution"
	TypeMsgProcessMaturity   = "process_maturity"
	TypeMsgRequestWithdrawal = "request_withdrawal"
	TypeMsgProcessWithdrawal = "process_withdrawal"
	TypeMsgUpdateKYCStatus   = "update_kyc_status"
	TypeMsgClaimReferral     = "claim_referral"
)

// Ensure messages implement sdk.Msg interface
var (
	_ sdk.Msg = &MsgCreateScheme{}
	_ sdk.Msg = &MsgUpdateScheme{}
	_ sdk.Msg = &MsgEnrollParticipant{}
	_ sdk.Msg = &MsgMakeContribution{}
	_ sdk.Msg = &MsgProcessMaturity{}
	_ sdk.Msg = &MsgRequestWithdrawal{}
	_ sdk.Msg = &MsgProcessWithdrawal{}
	_ sdk.Msg = &MsgUpdateKYCStatus{}
	_ sdk.Msg = &MsgClaimReferral{}
)

// MsgCreateScheme creates a new pension scheme
type MsgCreateScheme struct {
	Creator                sdk.AccAddress `json:"creator"`
	SchemeName             string         `json:"scheme_name"`
	Description            string         `json:"description"`
	MonthlyContribution    sdk.Coin       `json:"monthly_contribution"`
	ContributionPeriod     uint32         `json:"contribution_period"`
	MaturityBonus          sdk.Dec        `json:"maturity_bonus"`
	MinAge                 uint32         `json:"min_age"`
	MaxAge                 uint32         `json:"max_age"`
	MaxParticipants        uint64         `json:"max_participants"`
	GracePeriodDays        uint32         `json:"grace_period_days"`
	EarlyWithdrawalPenalty sdk.Dec        `json:"early_withdrawal_penalty"`
	LatePaymentPenalty     sdk.Dec        `json:"late_payment_penalty"`
	ReferralRewardPercent  sdk.Dec        `json:"referral_reward_percent"`
	OnTimeBonusPercent     sdk.Dec        `json:"on_time_bonus_percent"`
	KYCRequired            bool           `json:"kyc_required"`
	LiquidityProvision     bool           `json:"liquidity_provision"`
	LiquidityUtilization   sdk.Dec        `json:"liquidity_utilization"`
}

func NewMsgCreateScheme(creator sdk.AccAddress, schemeName, description string, monthlyContribution sdk.Coin,
	contributionPeriod uint32, maturityBonus sdk.Dec, minAge, maxAge uint32, maxParticipants uint64) *MsgCreateScheme {
	return &MsgCreateScheme{
		Creator:             creator,
		SchemeName:          schemeName,
		Description:         description,
		MonthlyContribution: monthlyContribution,
		ContributionPeriod:  contributionPeriod,
		MaturityBonus:       maturityBonus,
		MinAge:              minAge,
		MaxAge:              maxAge,
		MaxParticipants:     maxParticipants,
	}
}

func (msg *MsgCreateScheme) Route() string { return RouterKey }
func (msg *MsgCreateScheme) Type() string  { return TypeMsgCreateScheme }
func (msg *MsgCreateScheme) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}

func (msg *MsgCreateScheme) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateScheme) ValidateBasic() error {
	if msg.Creator.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "creator cannot be empty")
	}
	if msg.SchemeName == "" {
		return ErrInvalidSchemeName
	}
	if !msg.MonthlyContribution.IsPositive() {
		return ErrInvalidContribution
	}
	if msg.ContributionPeriod == 0 || msg.ContributionPeriod > 120 {
		return ErrInvalidContributionPeriod
	}
	if msg.MaturityBonus.IsNegative() || msg.MaturityBonus.GT(sdk.NewDec(1)) {
		return ErrInvalidMaturityBonus
	}
	if msg.MinAge == 0 || msg.MaxAge == 0 || msg.MinAge > msg.MaxAge {
		return ErrInvalidAgeRange
	}
	if msg.LiquidityUtilization.IsNegative() || msg.LiquidityUtilization.GT(sdk.NewDec(1)) {
		return ErrInvalidLiquidityUtilization
	}
	return nil
}

// MsgUpdateScheme updates an existing pension scheme
type MsgUpdateScheme struct {
	Authority   sdk.AccAddress `json:"authority"`
	SchemeID    string         `json:"scheme_id"`
	Description string         `json:"description,omitempty"`
	Status      string         `json:"status,omitempty"`
}

func (msg *MsgUpdateScheme) Route() string { return RouterKey }
func (msg *MsgUpdateScheme) Type() string  { return TypeMsgUpdateScheme }
func (msg *MsgUpdateScheme) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Authority}
}

func (msg *MsgUpdateScheme) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateScheme) ValidateBasic() error {
	if msg.Authority.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "authority cannot be empty")
	}
	if msg.SchemeID == "" {
		return ErrInvalidSchemeID
	}
	return nil
}

// MsgEnrollParticipant enrolls a participant in a pension scheme
type MsgEnrollParticipant struct {
	Participant       sdk.AccAddress `json:"participant"`
	SchemeID          string         `json:"scheme_id"`
	Name              string         `json:"name"`
	Age               uint32         `json:"age"`
	VillagePostalCode string         `json:"village_postal_code"`
	ReferrerAddress   sdk.AccAddress `json:"referrer_address,omitempty"`
}

func NewMsgEnrollParticipant(participant sdk.AccAddress, schemeID, name string, age uint32, villageCode string) *MsgEnrollParticipant {
	return &MsgEnrollParticipant{
		Participant:       participant,
		SchemeID:          schemeID,
		Name:              name,
		Age:               age,
		VillagePostalCode: villageCode,
	}
}

func (msg *MsgEnrollParticipant) Route() string { return RouterKey }
func (msg *MsgEnrollParticipant) Type() string  { return TypeMsgEnrollParticipant }
func (msg *MsgEnrollParticipant) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Participant}
}

func (msg *MsgEnrollParticipant) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgEnrollParticipant) ValidateBasic() error {
	if msg.Participant.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "participant cannot be empty")
	}
	if msg.SchemeID == "" {
		return ErrInvalidSchemeID
	}
	if msg.Name == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "name cannot be empty")
	}
	if msg.Age == 0 {
		return ErrInvalidAge
	}
	if !msg.ReferrerAddress.Empty() && msg.Participant.Equals(msg.ReferrerAddress) {
		return ErrSelfReferral
	}
	return nil
}

// MsgMakeContribution makes a monthly contribution
type MsgMakeContribution struct {
	Contributor   sdk.AccAddress `json:"contributor"`
	ParticipantID string         `json:"participant_id"`
	Amount        sdk.Coin       `json:"amount"`
	Month         uint32         `json:"month"`
}

func NewMsgMakeContribution(contributor sdk.AccAddress, participantID string, amount sdk.Coin, month uint32) *MsgMakeContribution {
	return &MsgMakeContribution{
		Contributor:   contributor,
		ParticipantID: participantID,
		Amount:        amount,
		Month:         month,
	}
}

func (msg *MsgMakeContribution) Route() string { return RouterKey }
func (msg *MsgMakeContribution) Type() string  { return TypeMsgMakeContribution }
func (msg *MsgMakeContribution) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Contributor}
}

func (msg *MsgMakeContribution) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgMakeContribution) ValidateBasic() error {
	if msg.Contributor.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "contributor cannot be empty")
	}
	if msg.ParticipantID == "" {
		return ErrInvalidParticipantID
	}
	if !msg.Amount.IsPositive() {
		return ErrInvalidContribution
	}
	if msg.Month == 0 {
		return ErrInvalidContributionMonth
	}
	return nil
}

// MsgProcessMaturity processes pension maturity
type MsgProcessMaturity struct {
	Authority     sdk.AccAddress `json:"authority"`
	ParticipantID string         `json:"participant_id"`
}

func (msg *MsgProcessMaturity) Route() string { return RouterKey }
func (msg *MsgProcessMaturity) Type() string  { return TypeMsgProcessMaturity }
func (msg *MsgProcessMaturity) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Authority}
}

func (msg *MsgProcessMaturity) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgProcessMaturity) ValidateBasic() error {
	if msg.Authority.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "authority cannot be empty")
	}
	if msg.ParticipantID == "" {
		return ErrInvalidParticipantID
	}
	return nil
}

// MsgRequestWithdrawal requests early withdrawal
type MsgRequestWithdrawal struct {
	Participant   sdk.AccAddress `json:"participant"`
	ParticipantID string         `json:"participant_id"`
	Reason        string         `json:"reason"`
}

func (msg *MsgRequestWithdrawal) Route() string { return RouterKey }
func (msg *MsgRequestWithdrawal) Type() string  { return TypeMsgRequestWithdrawal }
func (msg *MsgRequestWithdrawal) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Participant}
}

func (msg *MsgRequestWithdrawal) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRequestWithdrawal) ValidateBasic() error {
	if msg.Participant.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "participant cannot be empty")
	}
	if msg.ParticipantID == "" {
		return ErrInvalidParticipantID
	}
	if msg.Reason == "" {
		return ErrInvalidWithdrawalReason
	}
	return nil
}

// MsgProcessWithdrawal processes withdrawal request
type MsgProcessWithdrawal struct {
	Authority    sdk.AccAddress `json:"authority"`
	WithdrawalID string         `json:"withdrawal_id"`
	Approved     bool           `json:"approved"`
}

func (msg *MsgProcessWithdrawal) Route() string { return RouterKey }
func (msg *MsgProcessWithdrawal) Type() string  { return TypeMsgProcessWithdrawal }
func (msg *MsgProcessWithdrawal) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Authority}
}

func (msg *MsgProcessWithdrawal) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgProcessWithdrawal) ValidateBasic() error {
	if msg.Authority.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "authority cannot be empty")
	}
	if msg.WithdrawalID == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "withdrawal id cannot be empty")
	}
	return nil
}

// MsgUpdateKYCStatus updates KYC status
type MsgUpdateKYCStatus struct {
	Authority sdk.AccAddress `json:"authority"`
	Address   sdk.AccAddress `json:"address"`
	KYCStatus string         `json:"kyc_status"`
	KYCLevel  string         `json:"kyc_level"`
}

func (msg *MsgUpdateKYCStatus) Route() string { return RouterKey }
func (msg *MsgUpdateKYCStatus) Type() string  { return TypeMsgUpdateKYCStatus }
func (msg *MsgUpdateKYCStatus) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Authority}
}

func (msg *MsgUpdateKYCStatus) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateKYCStatus) ValidateBasic() error {
	if msg.Authority.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "authority cannot be empty")
	}
	if msg.Address.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "address cannot be empty")
	}
	if msg.KYCStatus == "" {
		return ErrInvalidKYCStatus
	}
	return nil
}

// MsgClaimReferral claims referral rewards
type MsgClaimReferral struct {
	Referrer sdk.AccAddress `json:"referrer"`
}

func (msg *MsgClaimReferral) Route() string { return RouterKey }
func (msg *MsgClaimReferral) Type() string  { return TypeMsgClaimReferral }
func (msg *MsgClaimReferral) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Referrer}
}

func (msg *MsgClaimReferral) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgClaimReferral) ValidateBasic() error {
	if msg.Referrer.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "referrer cannot be empty")
	}
	return nil
}