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
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateLaunch     = "create_launch"
	TypeMsgParticipateLaunch = "participate_launch"
	TypeMsgCompleteLaunch   = "complete_launch"
	TypeMsgCancelLaunch     = "cancel_launch"
	TypeMsgClaimTokens      = "claim_tokens"
	TypeMsgInitiateCommunityVeto = "initiate_community_veto"
	TypeMsgVoteCommunityVeto = "vote_community_veto"
	TypeMsgClaimCreatorReward = "claim_creator_reward"
	TypeMsgUpdateAntiPumpConfig = "update_anti_pump_config"
	TypeMsgEmergencyStop    = "emergency_stop"
)

var (
	_ sdk.Msg = &MsgCreateLaunch{}
	_ sdk.Msg = &MsgParticipateLaunch{}
	_ sdk.Msg = &MsgCompleteLaunch{}
	_ sdk.Msg = &MsgCancelLaunch{}
	_ sdk.Msg = &MsgClaimTokens{}
	_ sdk.Msg = &MsgInitiateCommunityVeto{}
	_ sdk.Msg = &MsgVoteCommunityVeto{}
	_ sdk.Msg = &MsgClaimCreatorReward{}
	_ sdk.Msg = &MsgUpdateAntiPumpConfig{}
	_ sdk.Msg = &MsgEmergencyStop{}
)

// MsgCreateLaunch represents a message to create a new token launch
type MsgCreateLaunch struct {
	Creator          string         `json:"creator" yaml:"creator"`
	TokenName        string         `json:"token_name" yaml:"token_name"`
	TokenSymbol      string         `json:"token_symbol" yaml:"token_symbol"`
	TokenDescription string         `json:"token_description" yaml:"token_description"`
	TotalSupply      sdk.Int        `json:"total_supply" yaml:"total_supply"`
	Decimals         uint32         `json:"decimals" yaml:"decimals"`
	LaunchType       string         `json:"launch_type" yaml:"launch_type"`
	TargetAmount     sdk.Int        `json:"target_amount" yaml:"target_amount"`
	MinContribution  sdk.Int        `json:"min_contribution" yaml:"min_contribution"`
	MaxContribution  sdk.Int        `json:"max_contribution" yaml:"max_contribution"`
	StartTime        time.Time      `json:"start_time" yaml:"start_time"`
	EndTime          time.Time      `json:"end_time" yaml:"end_time"`
	TradingDelay     int64          `json:"trading_delay" yaml:"trading_delay"`
	AntiPumpConfig   AntiPumpConfig `json:"anti_pump_config" yaml:"anti_pump_config"`
	CreatorPincode   string         `json:"creator_pincode" yaml:"creator_pincode"`
	CulturalQuote    string         `json:"cultural_quote" yaml:"cultural_quote"`
	Whitelist        []string       `json:"whitelist,omitempty" yaml:"whitelist,omitempty"`
}

func NewMsgCreateLaunch(
	creator, tokenName, tokenSymbol, tokenDescription string,
	totalSupply, targetAmount, minContribution, maxContribution sdk.Int,
	decimals uint32, launchType string,
	startTime, endTime time.Time, tradingDelay int64,
	antiPumpConfig AntiPumpConfig,
	creatorPincode, culturalQuote string,
	whitelist []string,
) *MsgCreateLaunch {
	return &MsgCreateLaunch{
		Creator:          creator,
		TokenName:        tokenName,
		TokenSymbol:      tokenSymbol,
		TokenDescription: tokenDescription,
		TotalSupply:      totalSupply,
		Decimals:         decimals,
		LaunchType:       launchType,
		TargetAmount:     targetAmount,
		MinContribution:  minContribution,
		MaxContribution:  maxContribution,
		StartTime:        startTime,
		EndTime:          endTime,
		TradingDelay:     tradingDelay,
		AntiPumpConfig:   antiPumpConfig,
		CreatorPincode:   creatorPincode,
		CulturalQuote:    culturalQuote,
		Whitelist:        whitelist,
	}
}

func (msg *MsgCreateLaunch) Route() string {
	return RouterKey
}

func (msg *MsgCreateLaunch) Type() string {
	return TypeMsgCreateLaunch
}

func (msg *MsgCreateLaunch) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateLaunch) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateLaunch) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.TokenName == "" {
		return ErrInvalidTokenName
	}

	if msg.TokenSymbol == "" {
		return ErrInvalidTokenSymbol
	}

	if msg.TotalSupply.IsZero() || msg.TotalSupply.IsNegative() {
		return ErrInvalidTotalSupply
	}

	if msg.TargetAmount.IsZero() || msg.TargetAmount.IsNegative() {
		return ErrInvalidTargetAmount
	}

	if msg.MinContribution.IsNegative() || msg.MaxContribution.IsNegative() {
		return ErrInvalidContribution
	}

	if msg.MinContribution.GT(msg.MaxContribution) {
		return ErrInvalidContribution
	}

	if msg.EndTime.Before(msg.StartTime) {
		return ErrInvalidConfiguration
	}

	// Validate anti-pump config
	if err := ValidateAntiPumpConfig(msg.AntiPumpConfig); err != nil {
		return err
	}

	// Validate PIN code
	if len(msg.CreatorPincode) != MaxPincodeLength {
		return ErrInvalidPincode
	}

	// Validate cultural quote length
	if len(msg.CulturalQuote) > MaxCulturalQuoteLength {
		return ErrInvalidCulturalQuote
	}

	// Validate launch type
	validLaunchTypes := []string{LaunchTypeFair, LaunchTypeStealth, LaunchTypeWhitelist, LaunchTypeAuction, LaunchTypePrivate}
	validType := false
	for _, validLT := range validLaunchTypes {
		if msg.LaunchType == validLT {
			validType = true
			break
		}
	}
	if !validType {
		return ErrInvalidConfiguration
	}

	return nil
}

// MsgParticipateLaunch represents a message to participate in a token launch
type MsgParticipateLaunch struct {
	Participant string  `json:"participant" yaml:"participant"`
	LaunchID    string  `json:"launch_id" yaml:"launch_id"`
	Amount      sdk.Int `json:"amount" yaml:"amount"`
}

func NewMsgParticipateLaunch(participant, launchID string, amount sdk.Int) *MsgParticipateLaunch {
	return &MsgParticipateLaunch{
		Participant: participant,
		LaunchID:    launchID,
		Amount:      amount,
	}
}

func (msg *MsgParticipateLaunch) Route() string {
	return RouterKey
}

func (msg *MsgParticipateLaunch) Type() string {
	return TypeMsgParticipateLaunch
}

func (msg *MsgParticipateLaunch) GetSigners() []sdk.AccAddress {
	participant, err := sdk.AccAddressFromBech32(msg.Participant)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{participant}
}

func (msg *MsgParticipateLaunch) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgParticipateLaunch) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Participant)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid participant address (%s)", err)
	}

	if msg.LaunchID == "" {
		return ErrLaunchNotFound
	}

	if msg.Amount.IsZero() || msg.Amount.IsNegative() {
		return ErrInvalidContribution
	}

	return nil
}

// MsgCompleteLaunch represents a message to complete a token launch
type MsgCompleteLaunch struct {
	Creator  string `json:"creator" yaml:"creator"`
	LaunchID string `json:"launch_id" yaml:"launch_id"`
}

func NewMsgCompleteLaunch(creator, launchID string) *MsgCompleteLaunch {
	return &MsgCompleteLaunch{
		Creator:  creator,
		LaunchID: launchID,
	}
}

func (msg *MsgCompleteLaunch) Route() string {
	return RouterKey
}

func (msg *MsgCompleteLaunch) Type() string {
	return TypeMsgCompleteLaunch
}

func (msg *MsgCompleteLaunch) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCompleteLaunch) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCompleteLaunch) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.LaunchID == "" {
		return ErrLaunchNotFound
	}

	return nil
}

// MsgCancelLaunch represents a message to cancel a token launch
type MsgCancelLaunch struct {
	Creator  string `json:"creator" yaml:"creator"`
	LaunchID string `json:"launch_id" yaml:"launch_id"`
	Reason   string `json:"reason" yaml:"reason"`
}

func NewMsgCancelLaunch(creator, launchID, reason string) *MsgCancelLaunch {
	return &MsgCancelLaunch{
		Creator:  creator,
		LaunchID: launchID,
		Reason:   reason,
	}
}

func (msg *MsgCancelLaunch) Route() string {
	return RouterKey
}

func (msg *MsgCancelLaunch) Type() string {
	return TypeMsgCancelLaunch
}

func (msg *MsgCancelLaunch) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCancelLaunch) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCancelLaunch) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.LaunchID == "" {
		return ErrLaunchNotFound
	}

	return nil
}

// MsgClaimTokens represents a message to claim allocated tokens
type MsgClaimTokens struct {
	Participant string `json:"participant" yaml:"participant"`
	LaunchID    string `json:"launch_id" yaml:"launch_id"`
}

func NewMsgClaimTokens(participant, launchID string) *MsgClaimTokens {
	return &MsgClaimTokens{
		Participant: participant,
		LaunchID:    launchID,
	}
}

func (msg *MsgClaimTokens) Route() string {
	return RouterKey
}

func (msg *MsgClaimTokens) Type() string {
	return TypeMsgClaimTokens
}

func (msg *MsgClaimTokens) GetSigners() []sdk.AccAddress {
	participant, err := sdk.AccAddressFromBech32(msg.Participant)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{participant}
}

func (msg *MsgClaimTokens) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgClaimTokens) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Participant)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid participant address (%s)", err)
	}

	if msg.LaunchID == "" {
		return ErrLaunchNotFound
	}

	return nil
}

// MsgInitiateCommunityVeto represents a message to initiate community veto
type MsgInitiateCommunityVeto struct {
	Initiator string `json:"initiator" yaml:"initiator"`
	LaunchID  string `json:"launch_id" yaml:"launch_id"`
	Reason    string `json:"reason" yaml:"reason"`
}

func NewMsgInitiateCommunityVeto(initiator, launchID, reason string) *MsgInitiateCommunityVeto {
	return &MsgInitiateCommunityVeto{
		Initiator: initiator,
		LaunchID:  launchID,
		Reason:    reason,
	}
}

func (msg *MsgInitiateCommunityVeto) Route() string {
	return RouterKey
}

func (msg *MsgInitiateCommunityVeto) Type() string {
	return TypeMsgInitiateCommunityVeto
}

func (msg *MsgInitiateCommunityVeto) GetSigners() []sdk.AccAddress {
	initiator, err := sdk.AccAddressFromBech32(msg.Initiator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{initiator}
}

func (msg *MsgInitiateCommunityVeto) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgInitiateCommunityVeto) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Initiator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid initiator address (%s)", err)
	}

	if msg.LaunchID == "" {
		return ErrLaunchNotFound
	}

	if msg.Reason == "" {
		return ErrInvalidConfiguration
	}

	return nil
}

// MsgVoteCommunityVeto represents a message to vote on community veto
type MsgVoteCommunityVeto struct {
	Voter    string `json:"voter" yaml:"voter"`
	LaunchID string `json:"launch_id" yaml:"launch_id"`
	Vote     bool   `json:"vote" yaml:"vote"`
}

func NewMsgVoteCommunityVeto(voter, launchID string, vote bool) *MsgVoteCommunityVeto {
	return &MsgVoteCommunityVeto{
		Voter:    voter,
		LaunchID: launchID,
		Vote:     vote,
	}
}

func (msg *MsgVoteCommunityVeto) Route() string {
	return RouterKey
}

func (msg *MsgVoteCommunityVeto) Type() string {
	return TypeMsgVoteCommunityVeto
}

func (msg *MsgVoteCommunityVeto) GetSigners() []sdk.AccAddress {
	voter, err := sdk.AccAddressFromBech32(msg.Voter)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{voter}
}

func (msg *MsgVoteCommunityVeto) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgVoteCommunityVeto) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Voter)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid voter address (%s)", err)
	}

	if msg.LaunchID == "" {
		return ErrLaunchNotFound
	}

	return nil
}

// MsgClaimCreatorReward represents a message to claim creator rewards
type MsgClaimCreatorReward struct {
	Creator      string `json:"creator" yaml:"creator"`
	TokenAddress string `json:"token_address" yaml:"token_address"`
}

func NewMsgClaimCreatorReward(creator, tokenAddress string) *MsgClaimCreatorReward {
	return &MsgClaimCreatorReward{
		Creator:      creator,
		TokenAddress: tokenAddress,
	}
}

func (msg *MsgClaimCreatorReward) Route() string {
	return RouterKey
}

func (msg *MsgClaimCreatorReward) Type() string {
	return TypeMsgClaimCreatorReward
}

func (msg *MsgClaimCreatorReward) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgClaimCreatorReward) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgClaimCreatorReward) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.TokenAddress == "" {
		return ErrTokenAlreadyExists
	}

	return nil
}

// MsgUpdateAntiPumpConfig represents a message to update anti-pump configuration
type MsgUpdateAntiPumpConfig struct {
	Creator        string         `json:"creator" yaml:"creator"`
	TokenAddress   string         `json:"token_address" yaml:"token_address"`
	AntiPumpConfig AntiPumpConfig `json:"anti_pump_config" yaml:"anti_pump_config"`
}

func NewMsgUpdateAntiPumpConfig(creator, tokenAddress string, config AntiPumpConfig) *MsgUpdateAntiPumpConfig {
	return &MsgUpdateAntiPumpConfig{
		Creator:        creator,
		TokenAddress:   tokenAddress,
		AntiPumpConfig: config,
	}
}

func (msg *MsgUpdateAntiPumpConfig) Route() string {
	return RouterKey
}

func (msg *MsgUpdateAntiPumpConfig) Type() string {
	return TypeMsgUpdateAntiPumpConfig
}

func (msg *MsgUpdateAntiPumpConfig) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateAntiPumpConfig) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateAntiPumpConfig) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.TokenAddress == "" {
		return ErrTokenAlreadyExists
	}

	if err := ValidateAntiPumpConfig(msg.AntiPumpConfig); err != nil {
		return err
	}

	return nil
}

// MsgEmergencyStop represents a message to trigger emergency stop
type MsgEmergencyStop struct {
	Authority    string `json:"authority" yaml:"authority"`
	TokenAddress string `json:"token_address" yaml:"token_address"`
	Reason       string `json:"reason" yaml:"reason"`
}

func NewMsgEmergencyStop(authority, tokenAddress, reason string) *MsgEmergencyStop {
	return &MsgEmergencyStop{
		Authority:    authority,
		TokenAddress: tokenAddress,
		Reason:       reason,
	}
}

func (msg *MsgEmergencyStop) Route() string {
	return RouterKey
}

func (msg *MsgEmergencyStop) Type() string {
	return TypeMsgEmergencyStop
}

func (msg *MsgEmergencyStop) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgEmergencyStop) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgEmergencyStop) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}

	if msg.TokenAddress == "" {
		return ErrTokenAlreadyExists
	}

	if msg.Reason == "" {
		return ErrInvalidConfiguration
	}

	return nil
}