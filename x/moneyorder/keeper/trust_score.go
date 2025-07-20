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
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/deshchain/x/moneyorder/types"
)

// TrustScore represents a user's trust score
type TrustScore struct {
	Address              string
	Score                int32
	TotalTransactions    int64
	SuccessfulTrades     int64
	DisputesWon          int32
	DisputesLost         int32
	AverageResponseTime  int32 // minutes
	LastUpdated          int64
}

// GetUserTrustScore retrieves user's trust score
func (k Keeper) GetUserTrustScore(address string) int32 {
	// Simplified implementation - in production would fetch from store
	return 75 // Default trust score
}

// CalculateTrustScore calculates trust score based on various factors
func (k Keeper) CalculateTrustScore(ctx sdk.Context, address string) int32 {
	userStats := k.GetUserStats(ctx, address)
	if userStats == nil {
		return 50 // New user default score
	}
	
	score := float64(50) // Base score
	
	// Transaction success rate (0-30 points)
	if userStats.TotalTransactions > 0 {
		successRate := float64(userStats.SuccessfulTransactions) / float64(userStats.TotalTransactions)
		score += successRate * 30
	}
	
	// Volume score (0-20 points)
	volumeScore := math.Min(float64(userStats.TotalVolume.Amount.Int64())/1000000000, 20)
	score += volumeScore
	
	// Dispute score (-20 to +10 points)
	if userStats.DisputesInvolved > 0 {
		disputeRate := float64(userStats.DisputesLost) / float64(userStats.DisputesInvolved)
		score -= disputeRate * 20
		score += float64(userStats.DisputesWon) * 2
	}
	
	// Activity score (0-10 points)
	daysSinceLastActive := ctx.BlockTime().Sub(userStats.LastActiveTime).Hours() / 24
	if daysSinceLastActive < 7 {
		score += 10
	} else if daysSinceLastActive < 30 {
		score += 5
	}
	
	// Response time score (0-10 points)
	if userStats.AverageResponseTime < 300 { // 5 minutes
		score += 10
	} else if userStats.AverageResponseTime < 900 { // 15 minutes
		score += 5
	}
	
	// Cap between 0-100
	if score < 0 {
		score = 0
	} else if score > 100 {
		score = 100
	}
	
	return int32(score)
}

// UpdateUserStats updates user statistics after a transaction
func (k Keeper) updateUserStats(ctx sdk.Context, buyer, seller string, amount sdk.Coin) {
	// Update buyer stats
	buyerStats := k.GetUserStats(ctx, buyer)
	if buyerStats == nil {
		buyerStats = &types.UserStats{
			Address:         buyer,
			TotalVolume:     sdk.NewCoin(types.DefaultDenom, sdk.ZeroInt()),
			LastActiveTime:  ctx.BlockTime(),
		}
	}
	buyerStats.TotalTransactions++
	buyerStats.SuccessfulTransactions++
	buyerStats.TotalVolume = buyerStats.TotalVolume.Add(amount)
	buyerStats.LastActiveTime = ctx.BlockTime()
	k.SetUserStats(ctx, buyerStats)
	
	// Update seller stats
	sellerStats := k.GetUserStats(ctx, seller)
	if sellerStats == nil {
		sellerStats = &types.UserStats{
			Address:         seller,
			TotalVolume:     sdk.NewCoin(types.DefaultDenom, sdk.ZeroInt()),
			LastActiveTime:  ctx.BlockTime(),
		}
	}
	sellerStats.TotalTransactions++
	sellerStats.SuccessfulTransactions++
	sellerStats.TotalVolume = sellerStats.TotalVolume.Add(amount)
	sellerStats.LastActiveTime = ctx.BlockTime()
	k.SetUserStats(ctx, sellerStats)
}

// IncreaseTrustScore increases user's trust score
func (k Keeper) increaseTrustScore(ctx sdk.Context, address string, points int32) {
	stats := k.GetUserStats(ctx, address)
	if stats == nil {
		return
	}
	
	currentScore := k.CalculateTrustScore(ctx, address)
	newScore := currentScore + points
	if newScore > 100 {
		newScore = 100
	}
	
	// Store trust score adjustment
	k.SetTrustScoreAdjustment(ctx, address, newScore-currentScore)
}

// DecreaseTrustScore decreases user's trust score
func (k Keeper) decreaseTrustScore(ctx sdk.Context, address string, points int32) {
	stats := k.GetUserStats(ctx, address)
	if stats == nil {
		return
	}
	
	currentScore := k.CalculateTrustScore(ctx, address)
	newScore := currentScore - points
	if newScore < 0 {
		newScore = 0
	}
	
	// Store trust score adjustment
	k.SetTrustScoreAdjustment(ctx, address, newScore-currentScore)
}

// GetTrustLevel returns user's trust level based on score
func (k Keeper) GetTrustLevel(score int32) string {
	switch {
	case score >= 90:
		return "Diamond"
	case score >= 75:
		return "Gold"
	case score >= 60:
		return "Silver"
	case score >= 40:
		return "Bronze"
	default:
		return "New"
	}
}

// GetTrustBadge returns trust badge emoji based on level
func (k Keeper) GetTrustBadge(level string) string {
	switch level {
	case "Diamond":
		return "💎"
	case "Gold":
		return "🥇"
	case "Silver":
		return "🥈"
	case "Bronze":
		return "🥉"
	default:
		return "🆕"
	}
}

// CalculateTrustDiscount returns fee discount based on trust score
func (k Keeper) CalculateTrustDiscount(score int32) sdk.Dec {
	// Higher trust score = lower fees
	// 90+ = 50% discount
	// 75+ = 30% discount
	// 60+ = 20% discount
	// 40+ = 10% discount
	// <40 = 0% discount
	
	switch {
	case score >= 90:
		return sdk.NewDecWithPrec(50, 2) // 0.50
	case score >= 75:
		return sdk.NewDecWithPrec(30, 2) // 0.30
	case score >= 60:
		return sdk.NewDecWithPrec(20, 2) // 0.20
	case score >= 40:
		return sdk.NewDecWithPrec(10, 2) // 0.10
	default:
		return sdk.ZeroDec()
	}
}