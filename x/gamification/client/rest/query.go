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
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	
	"github.com/DeshChain/DeshChain-Ecosystem/x/gamification/types"
)

func registerQueryRoutes(clientCtx client.Context, r *mux.Router) {
	// Profile queries
	r.HandleFunc(
		"/gamification/profile/{address}",
		queryProfileHandlerFn(clientCtx),
	).Methods("GET")
	
	r.HandleFunc(
		"/gamification/profile/username/{username}",
		queryProfileByUsernameHandlerFn(clientCtx),
	).Methods("GET")
	
	// Achievement queries
	r.HandleFunc(
		"/gamification/achievements",
		queryAchievementsHandlerFn(clientCtx),
	).Methods("GET")
	
	r.HandleFunc(
		"/gamification/achievement/{achievement-id}",
		queryAchievementHandlerFn(clientCtx),
	).Methods("GET")
	
	// Leaderboard queries
	r.HandleFunc(
		"/gamification/leaderboard/{leaderboard-type}",
		queryLeaderboardHandlerFn(clientCtx),
	).Methods("GET")
	
	// Daily challenge query
	r.HandleFunc(
		"/gamification/daily-challenge",
		queryDailyChallengeHandlerFn(clientCtx),
	).Methods("GET")
	
	// Team battles queries
	r.HandleFunc(
		"/gamification/team-battles",
		queryTeamBattlesHandlerFn(clientCtx),
	).Methods("GET")
	
	r.HandleFunc(
		"/gamification/team-battle/{battle-id}",
		queryTeamBattleHandlerFn(clientCtx),
	).Methods("GET")
	
	// Social posts query
	r.HandleFunc(
		"/gamification/social-posts",
		querySocialPostsHandlerFn(clientCtx),
	).Methods("GET")
	
	// Humor quotes query
	r.HandleFunc(
		"/gamification/humor-quotes",
		queryHumorQuotesHandlerFn(clientCtx),
	).Methods("GET")
	
	// Level config query
	r.HandleFunc(
		"/gamification/level-config/{level}",
		queryLevelConfigHandlerFn(clientCtx),
	).Methods("GET")
	
	// Params query
	r.HandleFunc(
		"/gamification/params",
		queryParamsHandlerFn(clientCtx),
	).Methods("GET")
}

func queryProfileHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		address := vars[RestAddress]
		
		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}
		
		params := types.QueryProfileRequest{
			Address: address,
		}
		
		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		
		res, height, err := clientCtx.QueryWithData(
			fmt.Sprintf("custom/%s/profile/%s", types.QuerierRoute, address),
			bz,
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		
		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func queryProfileByUsernameHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars[RestUsername]
		
		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}
		
		params := types.QueryProfileByUsernameRequest{
			Username: username,
		}
		
		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		
		res, height, err := clientCtx.QueryWithData(
			fmt.Sprintf("custom/%s/profile/username/%s", types.QuerierRoute, username),
			bz,
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		
		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func queryAchievementsHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}
		
		// Parse query parameters
		category := r.URL.Query().Get("category")
		rarity := r.URL.Query().Get("rarity")
		
		params := types.QueryAchievementsRequest{}
		
		// TODO: Parse category and rarity enums
		
		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		
		res, height, err := clientCtx.QueryWithData(
			fmt.Sprintf("custom/%s/achievements", types.QuerierRoute),
			bz,
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		
		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func queryAchievementHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		achievementID := vars[RestAchievementID]
		
		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}
		
		// Convert achievement ID to uint64
		var id uint64
		fmt.Sscanf(achievementID, "%d", &id)
		
		params := types.QueryAchievementRequest{
			AchievementId: id,
		}
		
		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		
		res, height, err := clientCtx.QueryWithData(
			fmt.Sprintf("custom/%s/achievement/%s", types.QuerierRoute, achievementID),
			bz,
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		
		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func queryLeaderboardHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		leaderboardType := vars[RestLeaderboardType]
		
		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}
		
		params := types.QueryLeaderboardRequest{
			LeaderboardType: leaderboardType,
		}
		
		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		
		res, height, err := clientCtx.QueryWithData(
			fmt.Sprintf("custom/%s/leaderboard/%s", types.QuerierRoute, leaderboardType),
			bz,
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		
		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func queryDailyChallengeHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}
		
		date := r.URL.Query().Get("date")
		
		params := types.QueryDailyChallengeRequest{
			Date: date,
		}
		
		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		
		res, height, err := clientCtx.QueryWithData(
			fmt.Sprintf("custom/%s/daily-challenge", types.QuerierRoute),
			bz,
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		
		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func queryTeamBattlesHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}
		
		activeOnly := r.URL.Query().Get("active_only") == "true"
		
		params := types.QueryTeamBattlesRequest{
			ActiveOnly: activeOnly,
		}
		
		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		
		res, height, err := clientCtx.QueryWithData(
			fmt.Sprintf("custom/%s/team-battles", types.QuerierRoute),
			bz,
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		
		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func queryTeamBattleHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		battleID := vars[RestBattleID]
		
		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}
		
		// Convert battle ID to uint64
		var id uint64
		fmt.Sscanf(battleID, "%d", &id)
		
		params := types.QueryTeamBattleRequest{
			BattleId: id,
		}
		
		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		
		res, height, err := clientCtx.QueryWithData(
			fmt.Sprintf("custom/%s/team-battle/%s", types.QuerierRoute, battleID),
			bz,
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		
		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func querySocialPostsHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}
		
		developerAddress := r.URL.Query().Get("developer_address")
		platform := r.URL.Query().Get("platform")
		viralOnly := r.URL.Query().Get("viral_only") == "true"
		
		params := types.QuerySocialPostsRequest{
			DeveloperAddress: developerAddress,
			ViralOnly: viralOnly,
		}
		
		// TODO: Parse platform enum
		
		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		
		res, height, err := clientCtx.QueryWithData(
			fmt.Sprintf("custom/%s/social-posts", types.QuerierRoute),
			bz,
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		
		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func queryHumorQuotesHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}
		
		quoteType := r.URL.Query().Get("quote_type")
		category := r.URL.Query().Get("category")
		familyFriendlyOnly := r.URL.Query().Get("family_friendly_only") == "true"
		
		params := types.QueryHumorQuotesRequest{
			FamilyFriendlyOnly: familyFriendlyOnly,
		}
		
		// TODO: Parse quote type and category enums
		
		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		
		res, height, err := clientCtx.QueryWithData(
			fmt.Sprintf("custom/%s/humor-quotes", types.QuerierRoute),
			bz,
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		
		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func queryLevelConfigHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		level := vars[RestLevel]
		
		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}
		
		// Convert level to uint32
		var lvl uint32
		fmt.Sscanf(level, "%d", &lvl)
		
		params := types.QueryLevelConfigRequest{
			Level: lvl,
		}
		
		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		
		res, height, err := clientCtx.QueryWithData(
			fmt.Sprintf("custom/%s/level-config/%s", types.QuerierRoute, level),
			bz,
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		
		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func queryParamsHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}
		
		res, height, err := clientCtx.QueryWithData(
			fmt.Sprintf("custom/%s/params", types.QuerierRoute),
			nil,
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		
		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}