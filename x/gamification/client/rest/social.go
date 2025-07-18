package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	
	"github.com/deshchain/deshchain/x/gamification/types"
)

func registerSocialMediaRoutes(clientCtx client.Context, r *mux.Router) {
	// Share achievement endpoint
	r.HandleFunc(
		"/gamification/achievement/{achievement-id}/share",
		shareAchievementHandler(clientCtx),
	).Methods("POST")
	
	// Generate share links
	r.HandleFunc(
		"/gamification/achievement/{achievement-id}/share-links",
		getShareLinksHandler(clientCtx),
	).Methods("GET")
	
	// Generate achievement poster
	r.HandleFunc(
		"/gamification/achievement/{achievement-id}/poster",
		generatePosterHandler(clientCtx),
	).Methods("GET")
	
	// Get viral posts
	r.HandleFunc(
		"/gamification/viral-posts",
		getViralPostsHandler(clientCtx),
	).Methods("GET")
	
	// Generate daily motivation
	r.HandleFunc(
		"/gamification/daily-motivation/{address}",
		getDailyMotivationHandler(clientCtx),
	).Methods("GET")
	
	// Get random humor quote
	r.HandleFunc(
		"/gamification/random-quote",
		getRandomQuoteHandler(clientCtx),
	).Methods("GET")
	
	// Track share event
	r.HandleFunc(
		"/gamification/track-share",
		trackShareEventHandler(clientCtx),
	).Methods("POST")
}

// ShareAchievementReq defines the request for sharing an achievement
type ShareAchievementReq struct {
	BaseReq       rest.BaseReq `json:"base_req"`
	Creator       string       `json:"creator"`
	AchievementId uint64       `json:"achievement_id"`
	Platform      string       `json:"platform"`
	Content       string       `json:"content"`
}

func shareAchievementHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		achievementIDStr := vars[RestAchievementID]
		
		achievementID, err := strconv.ParseUint(achievementIDStr, 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid achievement ID")
			return
		}
		
		var req ShareAchievementReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}
		
		req.AchievementId = achievementID
		
		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}
		
		// Parse platform
		platform := types.SocialPlatform_SOCIAL_PLATFORM_TWITTER // Default
		switch req.Platform {
		case "twitter":
			platform = types.SocialPlatform_SOCIAL_PLATFORM_TWITTER
		case "discord":
			platform = types.SocialPlatform_SOCIAL_PLATFORM_DISCORD
		case "telegram":
			platform = types.SocialPlatform_SOCIAL_PLATFORM_TELEGRAM
		case "instagram":
			platform = types.SocialPlatform_SOCIAL_PLATFORM_INSTAGRAM
		case "whatsapp":
			platform = types.SocialPlatform_SOCIAL_PLATFORM_WHATSAPP
		case "linkedin":
			platform = types.SocialPlatform_SOCIAL_PLATFORM_LINKEDIN
		}
		
		msg := &types.MsgShareAchievement{
			Creator:       req.Creator,
			AchievementId: req.AchievementId,
			Platform:      platform,
			Content:       req.Content,
		}
		
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		
		// In a real implementation, this would submit the transaction
		// For now, we'll return a mock response
		response := map[string]interface{}{
			"success": true,
			"post_id": 12345,
			"share_url": "https://deshchain.io/share/abc123",
			"message": "Achievement shared successfully!",
		}
		
		rest.PostProcessResponse(w, clientCtx, response)
	}
}

func getShareLinksHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		achievementIDStr := vars[RestAchievementID]
		
		achievementID, err := strconv.ParseUint(achievementIDStr, 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid achievement ID")
			return
		}
		
		// Get user address from query params
		userAddress := r.URL.Query().Get("address")
		if userAddress == "" {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "address parameter required")
			return
		}
		
		// Generate share links (mock implementation)
		shareLinks := map[string]string{
			"twitter": "https://twitter.com/intent/tweet?text=Just%20unlocked%20achievement%20on%20DeshChain!",
			"linkedin": "https://www.linkedin.com/sharing/share-offsite/?url=https://deshchain.io/achievement/" + achievementIDStr,
			"whatsapp": "https://wa.me/?text=Check%20out%20my%20DeshChain%20achievement!",
			"telegram": "https://t.me/share/url?url=https://deshchain.io/achievement/" + achievementIDStr,
			"facebook": "https://www.facebook.com/sharer/sharer.php?u=https://deshchain.io/achievement/" + achievementIDStr,
		}
		
		response := map[string]interface{}{
			"achievement_id": achievementID,
			"share_links": shareLinks,
			"share_code": "ABC123",
			"qr_code": map[string]interface{}{
				"url": "https://deshchain.io/share/ABC123",
				"image": "data:image/png;base64,iVBORw0KGgoAAAANS...", // Mock QR code
			},
		}
		
		rest.PostProcessResponse(w, clientCtx, response)
	}
}

func generatePosterHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		achievementIDStr := vars[RestAchievementID]
		
		achievementID, err := strconv.ParseUint(achievementIDStr, 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid achievement ID")
			return
		}
		
		// Get user address from query params
		userAddress := r.URL.Query().Get("address")
		if userAddress == "" {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "address parameter required")
			return
		}
		
		// Generate poster data (mock implementation)
		posterData := map[string]interface{}{
			"template": "Classic Bollywood",
			"theme": "Bollywood Gold",
			"aspectRatio": "2:3",
			"elements": []map[string]interface{}{
				{
					"type": "text",
					"content": "Dilwale Developer Le Jayenge",
					"position": map[string]int{
						"x": 50, "y": 30, "width": 700, "height": 100,
					},
					"style": map[string]interface{}{
						"fontFamily": "Bebas Neue",
						"fontSize": 72,
						"color": "#FFD700",
					},
				},
			},
			"effects": map[string]interface{}{
				"particles": true,
				"glow": true,
				"animation": "shimmer",
			},
			"metadata": map[string]interface{}{
				"achievement_id": achievementID,
				"user_address": userAddress,
				"created_at": "2024-01-15T10:00:00Z",
			},
		}
		
		rest.PostProcessResponse(w, clientCtx, posterData)
	}
}

func getViralPostsHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get query parameters
		limit := r.URL.Query().Get("limit")
		if limit == "" {
			limit = "10"
		}
		
		// Mock viral posts
		viralPosts := []map[string]interface{}{
			{
				"post_id": 1001,
				"developer": "@code_khan",
				"achievement": "Commit Kumar 420",
				"platform": "twitter",
				"content": "🎊 Just unlocked 'Commit Kumar 420' on @DeshChain! Khiladi 420 ban gaye!",
				"engagement": map[string]int{
					"likes": 5420,
					"shares": 420,
					"comments": 69,
					"views": 42000,
				},
				"virality_score": 15000,
				"posted_at": "2024-01-15T09:00:00Z",
			},
			{
				"post_id": 1002,
				"developer": "@bug_bahubali",
				"achievement": "Thala Bug Terminator",
				"platform": "linkedin",
				"content": "Proud to achieve 'Thala Bug Terminator' status on DeshChain!",
				"engagement": map[string]int{
					"likes": 3000,
					"shares": 200,
					"comments": 150,
					"views": 25000,
				},
				"virality_score": 12000,
				"posted_at": "2024-01-14T15:30:00Z",
			},
		}
		
		response := map[string]interface{}{
			"viral_posts": viralPosts,
			"total": len(viralPosts),
			"limit": limit,
		}
		
		rest.PostProcessResponse(w, clientCtx, response)
	}
}

func getDailyMotivationHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		address := vars[RestAddress]
		
		// Generate daily motivation (mock)
		motivation := map[string]interface{}{
			"greeting": "🌅 Good morning, developer!",
			"stats": map[string]interface{}{
				"level": 42,
				"next_level": 43,
				"progress": 65.5,
				"current_streak": 15,
				"total_earnings": "125,000 NAMO",
			},
			"next_achievement": map[string]interface{}{
				"name": "30 Din Ka Tapasya",
				"description": "Maintain a 30-day commit streak",
				"progress": "15/30 days",
			},
			"quote": "Code karo, duniya badlo!",
			"tip": "Complete 2 more bug fixes to unlock a special achievement!",
		}
		
		response := map[string]interface{}{
			"address": address,
			"motivation": motivation,
			"generated_at": "2024-01-15T06:00:00Z",
		}
		
		rest.PostProcessResponse(w, clientCtx, response)
	}
}

func getRandomQuoteHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get query parameters
		category := r.URL.Query().Get("category")
		quoteType := r.URL.Query().Get("type")
		
		// Return random quote (mock)
		quote := map[string]interface{}{
			"quote_id": 42,
			"text": "Bade bade deshon mein aisi choti choti bugs hoti rehti hai",
			"translation": "In big big countries, such small small bugs keep happening",
			"source": "DDLJ - Dilwale Debugging Le Jayenge",
			"character": "SRK",
			"type": "bollywood_dialogue",
			"category": "bug_fixes",
			"viral_score": 95,
			"is_family_friendly": true,
		}
		
		response := map[string]interface{}{
			"quote": quote,
			"filters": map[string]string{
				"category": category,
				"type": quoteType,
			},
		}
		
		rest.PostProcessResponse(w, clientCtx, response)
	}
}

// TrackShareEventReq defines the request for tracking a share event
type TrackShareEventReq struct {
	PostID   uint64 `json:"post_id"`
	Platform string `json:"platform"`
	Success  bool   `json:"success"`
	Error    string `json:"error,omitempty"`
}

func trackShareEventHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req TrackShareEventReq
		
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		
		// Track the share event (mock implementation)
		response := map[string]interface{}{
			"tracked": true,
			"event_id": 98765,
			"post_id": req.PostID,
			"platform": req.Platform,
			"tracked_at": "2024-01-15T10:30:00Z",
		}
		
		rest.PostProcessResponse(w, clientCtx, response)
	}
}