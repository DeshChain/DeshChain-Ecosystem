package rest

import (
	"github.com/gorilla/mux"
	
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

// RegisterRoutes registers gamification-related REST handlers to a router
func RegisterRoutes(clientCtx client.Context, r *mux.Router) {
	registerQueryRoutes(clientCtx, r)
	registerTxHandlers(clientCtx, r)
	registerSocialMediaRoutes(clientCtx, r)
}

const (
	RestAddress = "address"
	RestUsername = "username"
	RestAchievementID = "achievement-id"
	RestBattleID = "battle-id"
	RestLeaderboardType = "leaderboard-type"
	RestPlatform = "platform"
	RestLevel = "level"
)