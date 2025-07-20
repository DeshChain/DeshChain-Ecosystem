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

package keeper

import (
	"fmt"
	"math"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/deshchain/x/moneyorder/types"
)

// User Statistics and Trust Score Management

// SetUserStats stores user P2P statistics
func (k Keeper) SetUserStats(ctx sdk.Context, stats *types.UserP2PStats) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetUserStatsKey(stats.Address)
	value := k.cdc.MustMarshal(stats)
	store.Set(key, value)
	
	// Update trust score
	k.UpdateUserTrustScore(ctx, stats)
}

// GetUserStats retrieves user P2P statistics
func (k Keeper) GetUserStats(ctx sdk.Context, address string) *types.UserP2PStats {
	store := ctx.KVStore(k.storeKey)
	key := types.GetUserStatsKey(address)
	value := store.Get(key)
	if value == nil {
		return nil
	}
	
	var stats types.UserP2PStats
	k.cdc.MustUnmarshal(value, &stats)
	return &stats
}

// GetOrCreateUserStats gets existing stats or creates new ones
func (k Keeper) GetOrCreateUserStats(ctx sdk.Context, address string) *types.UserP2PStats {
	stats := k.GetUserStats(ctx, address)
	if stats == nil {
		stats = &types.UserP2PStats{
			Address:          address,
			TotalTrades:      0,
			SuccessfulTrades: 0,
			CancelledTrades:  0,
			DisputedTrades:   0,
			DisputesWon:      0,
			DisputesLost:     0,
			TotalVolume:      sdk.NewCoin(types.DefaultDenom, sdk.ZeroInt()),
			AverageTradeSize: sdk.NewCoin(types.DefaultDenom, sdk.ZeroInt()),
			TrustScore:       50, // Default trust score
			CreatedAt:        ctx.BlockTime(),
			LastTradeAt:      time.Time{},
			PreferredPaymentMethods: []string{},
			TradeResponseTime: 0,
			CompletionRate:   0,
		}
		k.SetUserStats(ctx, stats)
	}
	return stats
}

// UpdateUserTrustScore calculates and updates user trust score
func (k Keeper) UpdateUserTrustScore(ctx sdk.Context, stats *types.UserP2PStats) {
	if stats.TotalTrades == 0 {
		stats.TrustScore = 50 // Default for new users
		return
	}
	
	score := float64(50) // Base score
	
	// 1. Success Rate (0-30 points)
	successRate := float64(stats.SuccessfulTrades) / float64(stats.TotalTrades)
	score += successRate * 30
	
	// 2. Volume Score (0-20 points)
	volumeScore := k.calculateVolumeScore(stats.TotalVolume)
	score += volumeScore * 20
	
	// 3. Account Age Score (0-10 points)
	ageScore := k.calculateAccountAgeScore(ctx.BlockTime(), stats.CreatedAt)
	score += ageScore * 10
	
	// 4. Dispute Score (-20 to +10 points)
	disputeScore := k.calculateDisputeScore(stats)
	score += disputeScore
	
	// 5. Activity Score (0-10 points)
	activityScore := k.calculateActivityScore(ctx.BlockTime(), stats.LastTradeAt, stats.TotalTrades)
	score += activityScore
	
	// 6. Response Time Score (0-5 points)
	if stats.TradeResponseTime > 0 {
		responseScore := k.calculateResponseTimeScore(stats.TradeResponseTime)
		score += responseScore * 5
	}
	
	// 7. KYC Bonus (+5 points)
	if k.IsKYCVerified(stats.Address) {
		score += 5
	}
	
	// Cap between 0-100
	if score < 0 {
		score = 0
	} else if score > 100 {
		score = 100
	}
	
	stats.TrustScore = int32(score)
	
	// Update completion rate
	if stats.TotalTrades > 0 {
		stats.CompletionRate = float64(stats.SuccessfulTrades) / float64(stats.TotalTrades)
	}
}

// GetUserTrustScore returns the current trust score for a user
func (k Keeper) GetUserTrustScore(address string) int32 {
	// Note: This is called from the matching engine, so we need to avoid
	// circular dependency by not accessing ctx here
	// In production, would cache trust scores
	return 75 // Default trust score
}

// GetUserTrustScoreWithContext returns trust score with context
func (k Keeper) GetUserTrustScoreWithContext(ctx sdk.Context, address string) int32 {
	stats := k.GetUserStats(ctx, address)
	if stats == nil {
		return 50 // Default for new users
	}
	return stats.TrustScore
}

// RecordTradeStart records the start of a trade
func (k Keeper) RecordTradeStart(ctx sdk.Context, address string) {
	stats := k.GetOrCreateUserStats(ctx, address)
	stats.TotalTrades++
	stats.LastTradeStartTime = ctx.BlockTime()
	k.SetUserStats(ctx, stats)
}

// RecordTradeCompletion records successful trade completion
func (k Keeper) RecordTradeCompletion(ctx sdk.Context, address string, volume sdk.Coin, responseTime time.Duration) {
	stats := k.GetOrCreateUserStats(ctx, address)
	stats.SuccessfulTrades++
	stats.TotalVolume = stats.TotalVolume.Add(volume)
	stats.LastTradeAt = ctx.BlockTime()
	
	// Update average trade size
	if stats.SuccessfulTrades > 0 {
		avgAmount := stats.TotalVolume.Amount.Quo(sdk.NewInt(int64(stats.SuccessfulTrades)))
		stats.AverageTradeSize = sdk.NewCoin(types.DefaultDenom, avgAmount)
	}
	
	// Update average response time
	if stats.TradeResponseTime == 0 {
		stats.TradeResponseTime = responseTime
	} else {
		// Weighted average
		stats.TradeResponseTime = (stats.TradeResponseTime*time.Duration(stats.SuccessfulTrades-1) + responseTime) / time.Duration(stats.SuccessfulTrades)
	}
	
	k.SetUserStats(ctx, stats)
}

// RecordTradeCancellation records a cancelled trade
func (k Keeper) RecordTradeCancellation(ctx sdk.Context, address string, reason string) {
	stats := k.GetOrCreateUserStats(ctx, address)
	stats.CancelledTrades++
	
	// Track cancellation reasons
	if stats.CancellationReasons == nil {
		stats.CancellationReasons = make(map[string]int32)
	}
	stats.CancellationReasons[reason]++
	
	k.SetUserStats(ctx, stats)
}

// RecordDispute records a trade dispute
func (k Keeper) RecordDispute(ctx sdk.Context, address string, won bool) {
	stats := k.GetOrCreateUserStats(ctx, address)
	stats.DisputedTrades++
	
	if won {
		stats.DisputesWon++
	} else {
		stats.DisputesLost++
	}
	
	k.SetUserStats(ctx, stats)
}

// UpdatePreferredPaymentMethods updates user's preferred payment methods
func (k Keeper) UpdatePreferredPaymentMethods(ctx sdk.Context, address string, methods []string) {
	stats := k.GetOrCreateUserStats(ctx, address)
	stats.PreferredPaymentMethods = methods
	k.SetUserStats(ctx, stats)
}

// Helper functions for trust score calculation

func (k Keeper) calculateVolumeScore(volume sdk.Coin) float64 {
	// Logarithmic scale based on volume
	volumeNAMO := float64(volume.Amount.Int64()) / 1000000 // Convert to NAMO
	
	if volumeNAMO <= 0 {
		return 0
	}
	
	// Log scale with cap at 1M NAMO
	score := math.Log10(volumeNAMO) / math.Log10(1000000) // Max score at 1M NAMO
	if score > 1 {
		score = 1
	}
	
	return score
}

func (k Keeper) calculateAccountAgeScore(currentTime, createdAt time.Time) float64 {
	age := currentTime.Sub(createdAt)
	days := age.Hours() / 24
	
	// Max score at 365 days
	score := days / 365
	if score > 1 {
		score = 1
	}
	
	return score
}

func (k Keeper) calculateDisputeScore(stats *types.UserP2PStats) float64 {
	if stats.DisputedTrades == 0 {
		return 0 // No disputes, neutral score
	}
	
	// Dispute rate penalty
	disputeRate := float64(stats.DisputedTrades) / float64(stats.TotalTrades)
	penalty := disputeRate * 20 // Max -20 points for high dispute rate
	
	// Win rate bonus
	if stats.DisputedTrades > 0 {
		winRate := float64(stats.DisputesWon) / float64(stats.DisputedTrades)
		bonus := winRate * 10 // Max +10 points for winning disputes
		return bonus - penalty
	}
	
	return -penalty
}

func (k Keeper) calculateActivityScore(currentTime, lastTradeTime time.Time, totalTrades int32) float64 {
	if totalTrades == 0 {
		return 0
	}
	
	// Recency score
	daysSinceLastTrade := currentTime.Sub(lastTradeTime).Hours() / 24
	recencyScore := 0.0
	
	if daysSinceLastTrade <= 7 {
		recencyScore = 1.0
	} else if daysSinceLastTrade <= 30 {
		recencyScore = 0.7
	} else if daysSinceLastTrade <= 90 {
		recencyScore = 0.5
	} else if daysSinceLastTrade <= 180 {
		recencyScore = 0.3
	}
	
	// Frequency score
	frequencyScore := math.Min(float64(totalTrades)/100, 1.0) // Cap at 100 trades
	
	// Combined score
	return (recencyScore + frequencyScore) / 2
}

func (k Keeper) calculateResponseTimeScore(avgResponseTime time.Duration) float64 {
	minutes := avgResponseTime.Minutes()
	
	if minutes <= 5 {
		return 1.0 // Excellent
	} else if minutes <= 15 {
		return 0.8 // Good
	} else if minutes <= 30 {
		return 0.6 // Average
	} else if minutes <= 60 {
		return 0.4 // Below average
	} else {
		return 0.2 // Poor
	}
}

// GetUserReputation returns a human-readable reputation level
func (k Keeper) GetUserReputation(ctx sdk.Context, address string) string {
	trustScore := k.GetUserTrustScoreWithContext(ctx, address)
	
	switch {
	case trustScore >= 90:
		return "Diamond"
	case trustScore >= 80:
		return "Platinum"
	case trustScore >= 70:
		return "Gold"
	case trustScore >= 60:
		return "Silver"
	case trustScore >= 50:
		return "Bronze"
	default:
		return "New User"
	}
}

// GetTrustScoreBenefits returns benefits based on trust score
func (k Keeper) GetTrustScoreBenefits(trustScore int32) types.TrustScoreBenefits {
	benefits := types.TrustScoreBenefits{
		TrustScore: trustScore,
	}
	
	// Fee discounts
	if trustScore >= 90 {
		benefits.FeeDiscount = "0.15" // 0.15% discount
		benefits.PriorityMatching = true
		benefits.IncreasedLimits = true
		benefits.FastDispute = true
		benefits.Badge = "Diamond Trader"
	} else if trustScore >= 80 {
		benefits.FeeDiscount = "0.10"
		benefits.PriorityMatching = true
		benefits.IncreasedLimits = true
		benefits.Badge = "Platinum Trader"
	} else if trustScore >= 70 {
		benefits.FeeDiscount = "0.05"
		benefits.PriorityMatching = true
		benefits.Badge = "Gold Trader"
	} else if trustScore >= 60 {
		benefits.FeeDiscount = "0.02"
		benefits.Badge = "Silver Trader"
	} else if trustScore >= 50 {
		benefits.Badge = "Bronze Trader"
	} else {
		benefits.Badge = "New Trader"
	}
	
	return benefits
}

// GetUserTradeHistory returns recent trades for a user
func (k Keeper) GetUserTradeHistory(ctx sdk.Context, address string, limit int) []*types.TradeHistory {
	trades := k.GetP2PTradesByUser(ctx, address)
	
	// Sort by timestamp (most recent first)
	// In production, would implement proper sorting
	
	var history []*types.TradeHistory
	for i, trade := range trades {
		if i >= limit {
			break
		}
		
		// Determine user role in trade
		role := "buyer"
		counterparty := trade.Seller
		if trade.Seller == address {
			role = "seller"
			counterparty = trade.Buyer
		}
		
		history = append(history, &types.TradeHistory{
			TradeId:      trade.TradeId,
			Role:         role,
			Counterparty: counterparty,
			Amount:       trade.NamoAmount,
			FiatAmount:   trade.FiatAmount,
			Status:       trade.Status.String(),
			CreatedAt:    trade.CreatedAt,
			CompletedAt:  trade.CreatedAt, // Would track actual completion time
		})
	}
	
	return history
}

// BanUser bans a user from P2P trading
func (k Keeper) BanUser(ctx sdk.Context, address string, reason string, duration time.Duration) error {
	stats := k.GetOrCreateUserStats(ctx, address)
	stats.Banned = true
	stats.BanReason = reason
	stats.BannedUntil = ctx.BlockTime().Add(duration)
	k.SetUserStats(ctx, stats)
	
	// Cancel all active orders
	orders := k.GetP2POrdersByUser(ctx, address)
	for _, order := range orders {
		if order.Status == types.P2POrderStatus_P2P_STATUS_ACTIVE {
			order.Status = types.P2POrderStatus_P2P_STATUS_CANCELLED
			k.SetP2POrder(ctx, order)
		}
	}
	
	return nil
}

// IsUserBanned checks if a user is banned
func (k Keeper) IsUserBanned(ctx sdk.Context, address string) bool {
	stats := k.GetUserStats(ctx, address)
	if stats == nil || !stats.Banned {
		return false
	}
	
	// Check if ban has expired
	if ctx.BlockTime().After(stats.BannedUntil) {
		stats.Banned = false
		stats.BanReason = ""
		k.SetUserStats(ctx, stats)
		return false
	}
	
	return true
}

// GetP2PMarketStats returns overall P2P market statistics
func (k Keeper) GetP2PMarketStats(ctx sdk.Context) *types.P2PMarketStats {
	stats := &types.P2PMarketStats{
		TotalUsers:       0,
		ActiveUsers:      0,
		TotalVolume:      sdk.NewCoin(types.DefaultDenom, sdk.ZeroInt()),
		TotalTrades:      0,
		ActiveOrders:     0,
		AverageTrustScore: 0,
		TopTradingPairs:  make(map[string]int64),
		VolumeByState:    make(map[string]sdk.Int),
	}
	
	// Aggregate user stats
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.UserStatsPrefix)
	defer iterator.Close()
	
	totalTrustScore := int64(0)
	for ; iterator.Valid(); iterator.Next() {
		var userStats types.UserP2PStats
		k.cdc.MustUnmarshal(iterator.Value(), &userStats)
		
		stats.TotalUsers++
		if userStats.LastTradeAt.After(ctx.BlockTime().Add(-30 * 24 * time.Hour)) {
			stats.ActiveUsers++
		}
		
		stats.TotalVolume = stats.TotalVolume.Add(userStats.TotalVolume)
		stats.TotalTrades += int64(userStats.SuccessfulTrades)
		totalTrustScore += int64(userStats.TrustScore)
	}
	
	// Calculate average trust score
	if stats.TotalUsers > 0 {
		stats.AverageTrustScore = float64(totalTrustScore) / float64(stats.TotalUsers)
	}
	
	// Count active orders
	activeOrders := k.GetActiveP2POrders(ctx)
	stats.ActiveOrders = int32(len(activeOrders))
	
	// Get volume by state (simplified)
	for _, order := range activeOrders {
		if _, exists := stats.VolumeByState[order.State]; !exists {
			stats.VolumeByState[order.State] = sdk.ZeroInt()
		}
		stats.VolumeByState[order.State] = stats.VolumeByState[order.State].Add(order.Amount.Amount)
	}
	
	return stats
}