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

package rest

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	
	"github.com/deshchain/deshchain/x/gamification/types"
)

func registerTxHandlers(clientCtx client.Context, r *mux.Router) {
	r.HandleFunc("/gamification/profile", createProfileHandler(clientCtx)).Methods("POST")
	r.HandleFunc("/gamification/profile", updateProfileHandler(clientCtx)).Methods("PUT")
	r.HandleFunc("/gamification/avatar", selectAvatarHandler(clientCtx)).Methods("POST")
	r.HandleFunc("/gamification/achievement/claim", claimAchievementHandler(clientCtx)).Methods("POST")
	r.HandleFunc("/gamification/action", recordActionHandler(clientCtx)).Methods("POST")
	r.HandleFunc("/gamification/team-battle", createTeamBattleHandler(clientCtx)).Methods("POST")
	r.HandleFunc("/gamification/team-battle/join", joinTeamBattleHandler(clientCtx)).Methods("POST")
	r.HandleFunc("/gamification/daily-challenge/complete", completeDailyChallengeHandler(clientCtx)).Methods("POST")
}

// CreateProfileReq defines the request for creating a profile
type CreateProfileReq struct {
	BaseReq          rest.BaseReq    `json:"base_req"`
	Creator          string          `json:"creator"`
	GithubUsername   string          `json:"github_username"`
	PreferredAvatar  string          `json:"preferred_avatar"`
	PreferredLanguage string          `json:"preferred_language"`
	Region           string          `json:"region"`
	HumorPreference  string          `json:"humor_preference"`
}

func createProfileHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateProfileReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// Parse avatar type
		avatarType := types.AvatarType_AVATAR_TYPE_BUG_BUSTER_BAHUBALI // Default
		// TODO: Parse from string

		// Parse humor preference
		humorPref := types.HumorPreference_HUMOR_PREFERENCE_MIXED // Default
		// TODO: Parse from string

		msg := types.NewMsgCreateProfile(
			req.Creator,
			req.GithubUsername,
			avatarType,
			req.PreferredLanguage,
			req.Region,
			humorPref,
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

// UpdateProfileReq defines the request for updating a profile
type UpdateProfileReq struct {
	BaseReq          rest.BaseReq    `json:"base_req"`
	Creator          string          `json:"creator"`
	PreferredLanguage string          `json:"preferred_language"`
	Region           string          `json:"region"`
	HumorPreference  string          `json:"humor_preference"`
	SocialHandles    *types.SocialMediaHandles `json:"social_handles"`
}

func updateProfileHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UpdateProfileReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// Parse humor preference
		humorPref := types.HumorPreference_HUMOR_PREFERENCE_MIXED // Default
		// TODO: Parse from string

		msg := &types.MsgUpdateProfile{
			Creator:          req.Creator,
			PreferredLanguage: req.PreferredLanguage,
			Region:           req.Region,
			HumorPreference:  humorPref,
			SocialHandles:    req.SocialHandles,
		}
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

// SelectAvatarReq defines the request for selecting an avatar
type SelectAvatarReq struct {
	BaseReq    rest.BaseReq `json:"base_req"`
	Creator    string       `json:"creator"`
	AvatarType string       `json:"avatar_type"`
}

func selectAvatarHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SelectAvatarReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// Parse avatar type
		avatarType := types.AvatarType_AVATAR_TYPE_BUG_BUSTER_BAHUBALI // Default
		// TODO: Parse from string

		msg := &types.MsgSelectAvatar{
			Creator:    req.Creator,
			AvatarType: avatarType,
		}
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

// ClaimAchievementReq defines the request for claiming an achievement
type ClaimAchievementReq struct {
	BaseReq       rest.BaseReq `json:"base_req"`
	Creator       string       `json:"creator"`
	AchievementId uint64       `json:"achievement_id"`
}

func claimAchievementHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ClaimAchievementReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		msg := &types.MsgClaimAchievement{
			Creator:       req.Creator,
			AchievementId: req.AchievementId,
		}
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

// RecordActionReq defines the request for recording an action
type RecordActionReq struct {
	BaseReq    rest.BaseReq `json:"base_req"`
	Creator    string       `json:"creator"`
	ActionType string       `json:"action_type"`
	Value      uint64       `json:"value"`
	Metadata   string       `json:"metadata"`
}

func recordActionHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RecordActionReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		msg := &types.MsgRecordAction{
			Creator:    req.Creator,
			ActionType: req.ActionType,
			Value:      req.Value,
			Metadata:   req.Metadata,
		}
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

// CreateTeamBattleReq defines the request for creating a team battle
type CreateTeamBattleReq struct {
	BaseReq    rest.BaseReq `json:"base_req"`
	Creator    string       `json:"creator"`
	BattleType string       `json:"battle_type"`
	Duration   int64        `json:"duration"`
	PrizePool  sdk.Coin     `json:"prize_pool"`
	Team1Name  string       `json:"team_1_name"`
	Team2Name  string       `json:"team_2_name"`
}

func createTeamBattleHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateTeamBattleReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		msg := &types.MsgCreateTeamBattle{
			Creator:    req.Creator,
			BattleType: req.BattleType,
			Duration:   req.Duration,
			PrizePool:  req.PrizePool,
			Team1Name:  req.Team1Name,
			Team2Name:  req.Team2Name,
		}
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

// JoinTeamBattleReq defines the request for joining a team battle
type JoinTeamBattleReq struct {
	BaseReq  rest.BaseReq `json:"base_req"`
	Creator  string       `json:"creator"`
	BattleId uint64       `json:"battle_id"`
	TeamName string       `json:"team_name"`
}

func joinTeamBattleHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req JoinTeamBattleReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		msg := &types.MsgJoinTeamBattle{
			Creator:  req.Creator,
			BattleId: req.BattleId,
			TeamName: req.TeamName,
		}
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

// CompleteDailyChallengeReq defines the request for completing a daily challenge
type CompleteDailyChallengeReq struct {
	BaseReq       rest.BaseReq `json:"base_req"`
	Creator       string       `json:"creator"`
	ChallengeDate string       `json:"challenge_date"`
	ProofData     string       `json:"proof_data"`
}

func completeDailyChallengeHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CompleteDailyChallengeReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		msg := &types.MsgCompleteDailyChallenge{
			Creator:       req.Creator,
			ChallengeDate: req.ChallengeDate,
			ProofData:     req.ProofData,
		}
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}