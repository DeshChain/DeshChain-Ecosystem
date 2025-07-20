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
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/namo/x/sikkebaaz/types"
)

// UpdateTradingMetrics updates trading metrics for a token
func (k Keeper) UpdateTradingMetrics(ctx sdk.Context, tokenAddress string, tradeVolume sdk.Int, tradeCount uint64, priceChange sdk.Dec) error {
	// Get existing metrics
	metrics, found := k.getTradingMetrics(ctx, tokenAddress)
	if !found {
		// Initialize new metrics
		metrics = types.TradingMetrics{
			TokenAddress:   tokenAddress,
			TotalVolume:    sdk.ZeroInt(),
			DailyVolume:    sdk.ZeroInt(),
			TotalTrades:    0,
			DailyTrades:    0,
			UniqueTraders:  0,
			CurrentPrice:   sdk.ZeroDec(),
			PriceChange24h: sdk.ZeroDec(),
			MarketCap:      sdk.ZeroInt(),
			Liquidity:      sdk.ZeroInt(),
			LastUpdated:    ctx.BlockTime(),
		}
	}

	// Reset daily counters if it's a new day
	if k.isNewDay(metrics.LastUpdated, ctx.BlockTime()) {
		metrics.DailyVolume = sdk.ZeroInt()
		metrics.DailyTrades = 0
	}

	// Update metrics
	metrics.TotalVolume = metrics.TotalVolume.Add(tradeVolume)
	metrics.DailyVolume = metrics.DailyVolume.Add(tradeVolume)
	metrics.TotalTrades += tradeCount
	metrics.DailyTrades += tradeCount
	metrics.PriceChange24h = priceChange
	metrics.LastUpdated = ctx.BlockTime()

	// Update market cap based on current price and total supply
	launch, found := k.getTokenLaunchByAddress(ctx, tokenAddress)
	if found {
		metrics.MarketCap = launch.TotalSupply.ToDec().Mul(metrics.CurrentPrice).TruncateInt()
	}

	// Store updated metrics
	k.setTradingMetrics(ctx, metrics)

	// Update creator rewards
	k.updateCreatorRewards(ctx, tokenAddress, tradeVolume)

	return nil
}

// ProcessTrade processes a trade and updates all relevant metrics
func (k Keeper) ProcessTrade(ctx sdk.Context, tokenAddress, trader string, volume sdk.Int, price sdk.Dec, isBuy bool) error {
	// Update trading metrics
	priceChange := k.calculatePriceChange(ctx, tokenAddress, price)
	if err := k.UpdateTradingMetrics(ctx, tokenAddress, volume, 1, priceChange); err != nil {
		return err
	}

	// Update trader metrics
	k.updateTraderMetrics(ctx, tokenAddress, trader, volume, isBuy)

	// Update liquidity metrics
	k.updateLiquidityMetrics(ctx, tokenAddress, volume, isBuy)

	// Check for trading milestones
	k.checkTradingMilestones(ctx, tokenAddress)

	// Emit trading event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"token_trade",
			sdk.NewAttribute("token_address", tokenAddress),
			sdk.NewAttribute("trader", trader),
			sdk.NewAttribute("volume", volume.String()),
			sdk.NewAttribute("price", price.String()),
			sdk.NewAttribute("type", func() string {
				if isBuy {
					return "buy"
				}
				return "sell"
			}()),
		),
	)

	return nil
}

// updateCreatorRewards updates creator trading rewards
func (k Keeper) updateCreatorRewards(ctx sdk.Context, tokenAddress string, tradeVolume sdk.Int) {
	// Get creator reward info
	launch, found := k.getTokenLaunchByAddress(ctx, tokenAddress)
	if !found {
		return
	}

	reward, found := k.getCreatorReward(ctx, launch.Creator, tokenAddress)
	if !found {
		// Initialize creator reward
		reward = types.CreatorReward{
			Creator:           launch.Creator,
			TokenAddress:      tokenAddress,
			RewardRate:        sdk.MustNewDecFromStr(types.CreatorTradingReward), // 2%
			AccumulatedReward: sdk.ZeroInt(),
			LastClaimedAt:     ctx.BlockTime(),
			TotalClaimed:      sdk.ZeroInt(),
			IsActive:          true,
		}
	}

	if !reward.IsActive {
		return
	}

	// Calculate trading reward (2% of volume)
	tradingReward := tradeVolume.ToDec().Mul(reward.RewardRate).TruncateInt()
	
	// Add to accumulated rewards
	reward.AccumulatedReward = reward.AccumulatedReward.Add(tradingReward)

	k.setCreatorReward(ctx, reward)

	// Emit creator reward event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreatorReward,
			sdk.NewAttribute("action", "accumulated"),
			sdk.NewAttribute("creator", reward.Creator),
			sdk.NewAttribute("token_address", tokenAddress),
			sdk.NewAttribute("reward_amount", tradingReward.String()),
			sdk.NewAttribute("total_accumulated", reward.AccumulatedReward.String()),
		),
	)
}

// updateTraderMetrics updates individual trader metrics
func (k Keeper) updateTraderMetrics(ctx sdk.Context, tokenAddress, trader string, volume sdk.Int, isBuy bool) {
	// This would track individual trader statistics
	// Implementation would maintain trader-specific metrics
	
	// For now, we'll update unique trader count in trading metrics
	metrics, found := k.getTradingMetrics(ctx, tokenAddress)
	if found {
		// Check if this is a new trader (simplified implementation)
		if k.isNewTrader(ctx, tokenAddress, trader) {
			metrics.UniqueTraders++
			k.setTradingMetrics(ctx, metrics)
		}
	}
}

// updateLiquidityMetrics updates liquidity-related metrics
func (k Keeper) updateLiquidityMetrics(ctx sdk.Context, tokenAddress string, volume sdk.Int, isBuy bool) {
	metrics, found := k.getTradingMetrics(ctx, tokenAddress)
	if !found {
		return
	}

	// Calculate liquidity impact
	// Large trades relative to current liquidity indicate lower liquidity
	if metrics.Liquidity.IsPositive() {
		liquidityImpact := volume.ToDec().Quo(metrics.Liquidity.ToDec())
		
		// If impact is high (>10%), adjust liquidity down
		if liquidityImpact.GT(sdk.MustNewDecFromStr("0.1")) {
			adjustmentFactor := sdk.OneDec().Sub(liquidityImpact.Mul(sdk.MustNewDecFromStr("0.5")))
			if adjustmentFactor.IsPositive() {
				metrics.Liquidity = metrics.Liquidity.ToDec().Mul(adjustmentFactor).TruncateInt()
			}
		}
	} else {
		// Initialize liquidity estimate based on trade volume
		metrics.Liquidity = volume.MulRaw(10) // Estimate 10x trade volume as liquidity
	}

	k.setTradingMetrics(ctx, metrics)
}

// checkTradingMilestones checks and awards trading milestone achievements
func (k Keeper) checkTradingMilestones(ctx sdk.Context, tokenAddress string) {
	metrics, found := k.getTradingMetrics(ctx, tokenAddress)
	if !found {
		return
	}

	// Check volume milestones
	volumeMilestones := []sdk.Int{
		sdk.NewInt(1000000).MulRaw(1e12),   // 1M NAMO
		sdk.NewInt(10000000).MulRaw(1e12),  // 10M NAMO
		sdk.NewInt(100000000).MulRaw(1e12), // 100M NAMO
	}

	for i, milestone := range volumeMilestones {
		if metrics.TotalVolume.GTE(milestone) && !k.hasMilestoneAchieved(ctx, tokenAddress, fmt.Sprintf("volume_%d", i)) {
			k.awardMilestoneReward(ctx, tokenAddress, fmt.Sprintf("volume_%d", i), milestone)
		}
	}

	// Check trade count milestones
	tradeMilestones := []uint64{1000, 10000, 100000}
	for i, milestone := range tradeMilestones {
		if metrics.TotalTrades >= milestone && !k.hasMilestoneAchieved(ctx, tokenAddress, fmt.Sprintf("trades_%d", i)) {
			k.awardMilestoneReward(ctx, tokenAddress, fmt.Sprintf("trades_%d", i), sdk.NewInt(int64(milestone)))
		}
	}

	// Check unique trader milestones
	traderMilestones := []uint64{100, 1000, 10000}
	for i, milestone := range traderMilestones {
		if metrics.UniqueTraders >= milestone && !k.hasMilestoneAchieved(ctx, tokenAddress, fmt.Sprintf("traders_%d", i)) {
			k.awardMilestoneReward(ctx, tokenAddress, fmt.Sprintf("traders_%d", i), sdk.NewInt(int64(milestone)))
		}
	}
}

// awardMilestoneReward awards milestone achievement rewards
func (k Keeper) awardMilestoneReward(ctx sdk.Context, tokenAddress, milestoneType string, milestoneValue sdk.Int) {
	// Get token launch for creator info
	launch, found := k.getTokenLaunchByAddress(ctx, tokenAddress)
	if !found {
		return
	}

	// Calculate milestone reward (0.1% of milestone value)
	rewardAmount := milestoneValue.MulRaw(1).QuoRaw(1000) // 0.1%

	// Award to creator
	creatorAddr, err := sdk.AccAddressFromBech32(launch.Creator)
	if err != nil {
		return
	}

	rewardCoins := sdk.NewCoins(sdk.NewCoin(types.DefaultDenom, rewardAmount))
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.CreatorRewardsPool, creatorAddr, rewardCoins); err != nil {
		k.Logger(ctx).Error("Failed to award milestone reward",
			"creator", launch.Creator,
			"milestone", milestoneType,
			"reward", rewardAmount,
			"error", err,
		)
		return
	}

	// Mark milestone as achieved
	k.setMilestoneAchieved(ctx, tokenAddress, milestoneType)

	// Emit milestone event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"trading_milestone",
			sdk.NewAttribute("token_address", tokenAddress),
			sdk.NewAttribute("creator", launch.Creator),
			sdk.NewAttribute("milestone_type", milestoneType),
			sdk.NewAttribute("milestone_value", milestoneValue.String()),
			sdk.NewAttribute("reward_amount", rewardAmount.String()),
		),
	)

	k.Logger(ctx).Info("Trading milestone achieved",
		"token_address", tokenAddress,
		"creator", launch.Creator,
		"milestone", milestoneType,
		"value", milestoneValue,
		"reward", rewardAmount,
	)
}

// GetTopTradingTokens returns top trading tokens by volume
func (k Keeper) GetTopTradingTokens(ctx sdk.Context, limit int) []types.TradingMetrics {
	var allMetrics []types.TradingMetrics
	
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixTradingMetrics)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var metrics types.TradingMetrics
		k.cdc.MustUnmarshal(iterator.Value(), &metrics)
		allMetrics = append(allMetrics, metrics)
	}

	// Sort by total volume (simplified bubble sort)
	for i := 0; i < len(allMetrics)-1; i++ {
		for j := 0; j < len(allMetrics)-i-1; j++ {
			if allMetrics[j].TotalVolume.LT(allMetrics[j+1].TotalVolume) {
				allMetrics[j], allMetrics[j+1] = allMetrics[j+1], allMetrics[j]
			}
		}
	}

	// Return top results
	if limit > len(allMetrics) {
		limit = len(allMetrics)
	}

	return allMetrics[:limit]
}

// GetTradingMetricsByToken returns trading metrics for a specific token
func (k Keeper) GetTradingMetricsByToken(ctx sdk.Context, tokenAddress string) (types.TradingMetrics, bool) {
	return k.getTradingMetrics(ctx, tokenAddress)
}

// Helper functions

func (k Keeper) isNewDay(lastUpdate, currentTime time.Time) bool {
	return lastUpdate.Day() != currentTime.Day() || 
		   lastUpdate.Month() != currentTime.Month() || 
		   lastUpdate.Year() != currentTime.Year()
}

func (k Keeper) calculatePriceChange(ctx sdk.Context, tokenAddress string, currentPrice sdk.Dec) sdk.Dec {
	metrics, found := k.getTradingMetrics(ctx, tokenAddress)
	if !found || metrics.CurrentPrice.IsZero() {
		return sdk.ZeroDec()
	}

	// Calculate percentage change
	priceChange := currentPrice.Sub(metrics.CurrentPrice).Quo(metrics.CurrentPrice)
	return priceChange
}

func (k Keeper) isNewTrader(ctx sdk.Context, tokenAddress, trader string) bool {
	// This would check if trader has traded this token before
	// Simplified implementation returns true for now
	return true
}

func (k Keeper) hasMilestoneAchieved(ctx sdk.Context, tokenAddress, milestoneType string) bool {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte("milestone_"), []byte(tokenAddress+"_"+milestoneType)...)
	return store.Has(key)
}

func (k Keeper) setMilestoneAchieved(ctx sdk.Context, tokenAddress, milestoneType string) {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte("milestone_"), []byte(tokenAddress+"_"+milestoneType)...)
	store.Set(key, []byte("achieved"))
}

// ClaimCreatorRewards allows creators to claim accumulated trading rewards
func (k Keeper) ClaimCreatorRewards(ctx sdk.Context, creator, tokenAddress string) error {
	reward, found := k.getCreatorReward(ctx, creator, tokenAddress)
	if !found {
		return types.ErrCreatorRewardNotFound
	}

	if reward.AccumulatedReward.IsZero() {
		return types.ErrNoRewardsToClaim
	}

	if !reward.IsActive {
		return types.ErrCreatorRewardInactive
	}

	// Transfer rewards to creator
	creatorAddr, err := sdk.AccAddressFromBech32(creator)
	if err != nil {
		return err
	}

	rewardCoins := sdk.NewCoins(sdk.NewCoin(types.DefaultDenom, reward.AccumulatedReward))
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.CreatorRewardsPool, creatorAddr, rewardCoins); err != nil {
		return err
	}

	// Update reward record
	claimedAmount := reward.AccumulatedReward
	reward.TotalClaimed = reward.TotalClaimed.Add(claimedAmount)
	reward.AccumulatedReward = sdk.ZeroInt()
	reward.LastClaimedAt = ctx.BlockTime()

	k.setCreatorReward(ctx, reward)

	// Emit claim event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreatorReward,
			sdk.NewAttribute("action", "claimed"),
			sdk.NewAttribute("creator", creator),
			sdk.NewAttribute("token_address", tokenAddress),
			sdk.NewAttribute("claimed_amount", claimedAmount.String()),
			sdk.NewAttribute("total_claimed", reward.TotalClaimed.String()),
		),
	)

	return nil
}

// GetCreatorRewardInfo returns creator reward information
func (k Keeper) GetCreatorRewardInfo(ctx sdk.Context, creator, tokenAddress string) (types.CreatorReward, bool) {
	return k.getCreatorReward(ctx, creator, tokenAddress)
}

// GetAllCreatorRewards returns all creator rewards
func (k Keeper) GetAllCreatorRewards(ctx sdk.Context) []types.CreatorReward {
	var rewards []types.CreatorReward
	
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixCreatorRewards)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var reward types.CreatorReward
		k.cdc.MustUnmarshal(iterator.Value(), &reward)
		rewards = append(rewards, reward)
	}

	return rewards
}

func (k Keeper) getCreatorReward(ctx sdk.Context, creator, tokenAddress string) (types.CreatorReward, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetCreatorRewardsKey(creator, tokenAddress))
	if bz == nil {
		return types.CreatorReward{}, false
	}

	var reward types.CreatorReward
	k.cdc.MustUnmarshal(bz, &reward)
	return reward, true
}