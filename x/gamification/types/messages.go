package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateProfile          = "create_profile"
	TypeMsgUpdateProfile          = "update_profile"
	TypeMsgSelectAvatar           = "select_avatar"
	TypeMsgClaimAchievement       = "claim_achievement"
	TypeMsgRecordAction           = "record_action"
	TypeMsgShareAchievement       = "share_achievement"
	TypeMsgJoinTeamBattle         = "join_team_battle"
	TypeMsgCreateTeamBattle       = "create_team_battle"
	TypeMsgCompleteDailyChallenge = "complete_daily_challenge"
)

var _ sdk.Msg = &MsgCreateProfile{}

// MsgCreateProfile creates a new developer profile
type MsgCreateProfile struct {
	Creator        string `json:"creator"`
	GithubUsername string `json:"github_username"`
	PreferredAvatar AvatarType `json:"preferred_avatar"`
	PreferredLanguage string `json:"preferred_language"`
	Region         string `json:"region"`
	HumorPreference HumorPreference `json:"humor_preference"`
}

func NewMsgCreateProfile(creator, githubUsername string, avatar AvatarType, language, region string, humor HumorPreference) *MsgCreateProfile {
	return &MsgCreateProfile{
		Creator:        creator,
		GithubUsername: githubUsername,
		PreferredAvatar: avatar,
		PreferredLanguage: language,
		Region:         region,
		HumorPreference: humor,
	}
}

func (msg *MsgCreateProfile) Route() string { return RouterKey }
func (msg *MsgCreateProfile) Type() string  { return TypeMsgCreateProfile }
func (msg *MsgCreateProfile) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateProfile) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateProfile) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.GithubUsername == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "github username cannot be empty")
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateProfile{}

// MsgUpdateProfile updates developer profile
type MsgUpdateProfile struct {
	Creator           string             `json:"creator"`
	PreferredLanguage string             `json:"preferred_language"`
	Region            string             `json:"region"`
	HumorPreference   HumorPreference    `json:"humor_preference"`
	SocialHandles     *SocialMediaHandles `json:"social_handles"`
}

func (msg *MsgUpdateProfile) Route() string { return RouterKey }
func (msg *MsgUpdateProfile) Type() string  { return TypeMsgUpdateProfile }
func (msg *MsgUpdateProfile) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateProfile) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateProfile) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgSelectAvatar{}

// MsgSelectAvatar selects a developer avatar
type MsgSelectAvatar struct {
	Creator    string     `json:"creator"`
	AvatarType AvatarType `json:"avatar_type"`
}

func (msg *MsgSelectAvatar) Route() string { return RouterKey }
func (msg *MsgSelectAvatar) Type() string  { return TypeMsgSelectAvatar }
func (msg *MsgSelectAvatar) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSelectAvatar) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSelectAvatar) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.AvatarType == AvatarType_AVATAR_TYPE_UNSPECIFIED {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "avatar type must be specified")
	}
	return nil
}

var _ sdk.Msg = &MsgClaimAchievement{}

// MsgClaimAchievement claims an achievement
type MsgClaimAchievement struct {
	Creator       string `json:"creator"`
	AchievementId uint64 `json:"achievement_id"`
}

func (msg *MsgClaimAchievement) Route() string { return RouterKey }
func (msg *MsgClaimAchievement) Type() string  { return TypeMsgClaimAchievement }
func (msg *MsgClaimAchievement) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgClaimAchievement) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgClaimAchievement) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.AchievementId == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "achievement id cannot be zero")
	}
	return nil
}

var _ sdk.Msg = &MsgRecordAction{}

// MsgRecordAction records a developer action
type MsgRecordAction struct {
	Creator    string `json:"creator"`
	ActionType string `json:"action_type"`
	Value      uint64 `json:"value"`
	Metadata   string `json:"metadata"`
}

func (msg *MsgRecordAction) Route() string { return RouterKey }
func (msg *MsgRecordAction) Type() string  { return TypeMsgRecordAction }
func (msg *MsgRecordAction) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRecordAction) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRecordAction) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.ActionType == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "action type cannot be empty")
	}
	return nil
}

var _ sdk.Msg = &MsgShareAchievement{}

// MsgShareAchievement shares an achievement on social media
type MsgShareAchievement struct {
	Creator       string         `json:"creator"`
	AchievementId uint64         `json:"achievement_id"`
	Platform      SocialPlatform `json:"platform"`
	Content       string         `json:"content"`
}

func (msg *MsgShareAchievement) Route() string { return RouterKey }
func (msg *MsgShareAchievement) Type() string  { return TypeMsgShareAchievement }
func (msg *MsgShareAchievement) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgShareAchievement) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgShareAchievement) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.AchievementId == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "achievement id cannot be zero")
	}
	if msg.Platform == SocialPlatform_SOCIAL_PLATFORM_UNSPECIFIED {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "platform must be specified")
	}
	return nil
}

var _ sdk.Msg = &MsgJoinTeamBattle{}

// MsgJoinTeamBattle joins a team battle
type MsgJoinTeamBattle struct {
	Creator  string `json:"creator"`
	BattleId uint64 `json:"battle_id"`
	TeamName string `json:"team_name"`
}

func (msg *MsgJoinTeamBattle) Route() string { return RouterKey }
func (msg *MsgJoinTeamBattle) Type() string  { return TypeMsgJoinTeamBattle }
func (msg *MsgJoinTeamBattle) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgJoinTeamBattle) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgJoinTeamBattle) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.BattleId == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "battle id cannot be zero")
	}
	if msg.TeamName == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "team name cannot be empty")
	}
	return nil
}

var _ sdk.Msg = &MsgCreateTeamBattle{}

// MsgCreateTeamBattle creates a new team battle
type MsgCreateTeamBattle struct {
	Creator    string        `json:"creator"`
	BattleType string        `json:"battle_type"`
	Duration   int64         `json:"duration"`
	PrizePool  sdk.Coin      `json:"prize_pool"`
	Team1Name  string        `json:"team_1_name"`
	Team2Name  string        `json:"team_2_name"`
}

func (msg *MsgCreateTeamBattle) Route() string { return RouterKey }
func (msg *MsgCreateTeamBattle) Type() string  { return TypeMsgCreateTeamBattle }
func (msg *MsgCreateTeamBattle) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateTeamBattle) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateTeamBattle) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.BattleType == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "battle type cannot be empty")
	}
	if msg.Duration <= 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "duration must be positive")
	}
	if !msg.PrizePool.IsValid() || msg.PrizePool.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid prize pool")
	}
	if msg.Team1Name == "" || msg.Team2Name == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "team names cannot be empty")
	}
	return nil
}

var _ sdk.Msg = &MsgCompleteDailyChallenge{}

// MsgCompleteDailyChallenge completes a daily challenge
type MsgCompleteDailyChallenge struct {
	Creator       string `json:"creator"`
	ChallengeDate string `json:"challenge_date"`
	ProofData     string `json:"proof_data"`
}

func (msg *MsgCompleteDailyChallenge) Route() string { return RouterKey }
func (msg *MsgCompleteDailyChallenge) Type() string  { return TypeMsgCompleteDailyChallenge }
func (msg *MsgCompleteDailyChallenge) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCompleteDailyChallenge) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCompleteDailyChallenge) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.ChallengeDate == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "challenge date cannot be empty")
	}
	return nil
}