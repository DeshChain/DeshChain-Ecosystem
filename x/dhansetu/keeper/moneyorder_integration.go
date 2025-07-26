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
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/dhansetu/types"
)

// MoneyOrderIntegrationHooks implements hooks for Money Order DEX integration
type MoneyOrderIntegrationHooks struct {
	keeper Keeper
}

// NewMoneyOrderIntegrationHooks creates new integration hooks
func (k Keeper) NewMoneyOrderIntegrationHooks() MoneyOrderIntegrationHooks {
	return MoneyOrderIntegrationHooks{keeper: k}
}

// AfterMoneyOrderCreated is called after a money order is created
func (h MoneyOrderIntegrationHooks) AfterMoneyOrderCreated(ctx sdk.Context, orderId, sender, receiver string, amount sdk.Coin) error {
	// Check if receiver is a DhanPata address (contains @dhan)
	if strings.Contains(receiver, "@dhan") {
		// Create bridge mapping for DhanPata integration
		return h.keeper.CreateOrderBridge(ctx, orderId, receiver)
	}
	return nil
}

// AfterP2PTradeMatched is called after P2P trade matching
func (h MoneyOrderIntegrationHooks) AfterP2PTradeMatched(ctx sdk.Context, tradeId, buyOrderId, sellOrderId string, amount sdk.Coin) error {
	// Record trade in DhanSetu unified history
	trade := types.TradeHistoryEntry{
		TradeId:       tradeId,
		TradeType:     "p2p_trade",
		SourceProduct: "moneyorder",
		Amount:        amount,
		Status:        "matched",
		Metadata: map[string]interface{}{
			"buy_order_id":  buyOrderId,
			"sell_order_id": sellOrderId,
		},
		Timestamp: ctx.BlockTime(),
	}

	h.keeper.RecordTradeHistory(ctx, trade)
	return nil
}

// AfterSevaMitraRegistered is called after a SevaMitra is registered
func (h MoneyOrderIntegrationHooks) AfterSevaMitraRegistered(ctx sdk.Context, mitraId, address string) error {
	// Check if mitra has a DhanPata address
	dhanpataName, found := h.keeper.GetDhanPataByAddress(ctx, address)
	if found {
		// Enhance the existing mitra profile with DhanSetu features
		return h.keeper.enhanceSevaMitra(ctx, mitraId, dhanpataName)
	}
	return nil
}

// enhanceSevaMitra upgrades a SevaMitra to an Enhanced Mitra with DhanSetu features
func (k Keeper) enhanceSevaMitra(ctx sdk.Context, mitraId, dhanpataName string) error {
	// Create enhanced mitra profile
	profile := types.EnhancedMitraProfile{
		MitraId:          mitraId,
		DhanPataName:     dhanpataName,
		MitraType:        types.MitraTypeIndividual, // Default to individual
		TrustScore:       50,                        // Base trust score
		DailyVolume:      sdk.ZeroInt(),
		MonthlyVolume:    sdk.ZeroInt(),
		TotalTrades:      0,
		SuccessfulTrades: 0,
		ActiveEscrows:    []string{},
		Specializations:  []string{"money_order", "p2p_trading"},
		OperatingRegions: []string{}, // Will be populated based on location
		PaymentMethods:   []types.PaymentMethod{},
		KYCStatus:        "pending",
		IsActive:         true,
	}

	// Calculate limits based on trust score
	daily, monthly := types.CalculateMitraLimits(profile.MitraType, profile.TrustScore)
	profile.DailyLimit = daily
	profile.MonthlyLimit = monthly

	return k.RegisterEnhancedMitra(ctx, profile)
}

// ResolveDhanPataForMoneyOrder resolves DhanPata address for money order processing
func (k Keeper) ResolveDhanPataForMoneyOrder(ctx sdk.Context, receiverUPI string) (string, error) {
	// Check if it's a DhanPata address
	if strings.HasSuffix(receiverUPI, "@dhan") {
		dhanpataAddr, found := k.GetDhanPataAddress(ctx, receiverUPI)
		if !found {
			return "", types.ErrDhanPataNotFound
		}
		return dhanpataAddr.BlockchainAddr, nil
	}

	// If it's already a blockchain address, return as-is
	return receiverUPI, nil
}

// ProcessDhanSetuMoneyOrder processes a money order with full DhanSetu integration
func (k Keeper) ProcessDhanSetuMoneyOrder(ctx sdk.Context, msg *types.MsgProcessMoneyOrderWithDhanPata) error {
	// Resolve DhanPata to blockchain address
	receiverAddr, err := k.ResolveDhanPataForMoneyOrder(ctx, msg.ReceiverDhanpata)
	if err != nil {
		return err
	}

	// Get sender address
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return err
	}

	// Calculate DhanSetu fee
	params := k.GetParams(ctx)
	feeAmount := params.DhanSetuFeeRate.MulInt(msg.Amount.Amount).TruncateInt()
	fee := sdk.NewCoin(msg.Amount.Denom, feeAmount)

	// Charge DhanSetu fee
	if err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx, senderAddr, types.ModuleName, sdk.NewCoins(fee),
	); err != nil {
		return err
	}

	// Create order ID
	orderId := fmt.Sprintf("dh-%d-%s", ctx.BlockHeight(), ctx.TxBytes())

	// Create money order through money order keeper
	// Note: This would call the actual money order keeper method
	// For now, we'll simulate the core functionality

	// Create bridge mapping
	err = k.CreateOrderBridge(ctx, orderId, msg.ReceiverDhanpata)
	if err != nil {
		return err
	}

	// Record in trade history
	trade := types.TradeHistoryEntry{
		TradeId:       orderId,
		UserDhanPata:  msg.ReceiverDhanpata,
		TradeType:     "dhansetu_money_order",
		SourceProduct: "dhansetu",
		Amount:        msg.Amount,
		Fee:           fee,
		Counterparty:  msg.Sender,
		Status:        "completed",
		Metadata: map[string]interface{}{
			"note":            msg.Note,
			"receiver_addr":   receiverAddr,
			"dhansetu_integration": true,
		},
		Timestamp: ctx.BlockTime(),
	}

	k.RecordTradeHistory(ctx, trade)

	// Distribute fees (40% platform, 40% charity, 20% founder)
	k.distributeDhanSetuFees(ctx, fee)

	// Emit events
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCrossModuleTransfer,
			sdk.NewAttribute(types.AttributeKeySourceModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyTargetModule, "moneyorder"),
			sdk.NewAttribute("sender", msg.Sender),
			sdk.NewAttribute("receiver_dhanpata", msg.ReceiverDhanpata),
			sdk.NewAttribute("receiver_address", receiverAddr),
			sdk.NewAttribute("amount", msg.Amount.String()),
			sdk.NewAttribute("fee", fee.String()),
			sdk.NewAttribute("order_id", orderId),
		),
	)

	return nil
}

// distributeDhanSetuFees distributes fees according to DeshChain Platform Revenue Model
func (k Keeper) distributeDhanSetuFees(ctx sdk.Context, totalFee sdk.Coin) error {
	// Calculate distribution per Platform Revenue Model
	developmentShare := totalFee.Amount.MulRaw(30).QuoRaw(100)    // 30% Development Fund
	communityShare := totalFee.Amount.MulRaw(25).QuoRaw(100)      // 25% Community Treasury
	liquidityShare := totalFee.Amount.MulRaw(20).QuoRaw(100)      // 20% Liquidity Provision
	ngoShare := totalFee.Amount.MulRaw(10).QuoRaw(100)           // 10% NGO Donations
	emergencyShare := totalFee.Amount.MulRaw(10).QuoRaw(100)      // 10% Emergency Reserve
	founderShare := totalFee.Amount.MulRaw(5).QuoRaw(100)        // 5% Founder Royalty

	// Create coins for distribution
	developmentCoin := sdk.NewCoin(totalFee.Denom, developmentShare)
	communityCoin := sdk.NewCoin(totalFee.Denom, communityShare)
	liquidityCoin := sdk.NewCoin(totalFee.Denom, liquidityShare)
	ngoCoin := sdk.NewCoin(totalFee.Denom, ngoShare)
	emergencyCoin := sdk.NewCoin(totalFee.Denom, emergencyShare)
	founderCoin := sdk.NewCoin(totalFee.Denom, founderShare)

	// Send to respective module accounts per established model
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)

	// Development Fund (30% - platform development and features)
	if err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx, moduleAddr, "platform_development_fund", sdk.NewCoins(developmentCoin),
	); err != nil {
		return err
	}

	// Community Treasury (25% - community programs and governance)
	if err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx, moduleAddr, "platform_community_treasury", sdk.NewCoins(communityCoin),
	); err != nil {
		return err
	}

	// Liquidity Provision (20% - market making and DEX liquidity)
	if err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx, moduleAddr, "platform_liquidity_pool", sdk.NewCoins(liquidityCoin),
	); err != nil {
		return err
	}

	// NGO Donations (10% - social impact and charity)
	if err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx, moduleAddr, "ngo_donation_pool", sdk.NewCoins(ngoCoin),
	); err != nil {
		return err
	}

	// Emergency Reserve (10% - risk management and stability)
	if err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx, moduleAddr, "platform_emergency_reserve", sdk.NewCoins(emergencyCoin),
	); err != nil {
		return err
	}

	// Founder Royalty (5% - perpetual compensation per platform revenue model)
	if err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx, moduleAddr, "founder_royalty_pool", sdk.NewCoins(founderCoin),
	); err != nil {
		return err
	}

	return nil
}

// GetDhanPataOrderHistory returns order history for a DhanPata user
func (k Keeper) GetDhanPataOrderHistory(ctx sdk.Context, dhanpataName string) []types.TradeHistoryEntry {
	return k.GetTradeHistory(ctx, dhanpataName)
}

// GetMitraPerformanceMetrics returns performance metrics for a mitra
func (k Keeper) GetMitraPerformanceMetrics(ctx sdk.Context, mitraId string) (*types.MitraPerformanceMetrics, error) {
	profile, found := k.GetEnhancedMitraProfile(ctx, mitraId)
	if !found {
		return nil, types.ErrMitraNotFound
	}

	// Calculate success rate
	successRate := float64(0)
	if profile.TotalTrades > 0 {
		successRate = float64(profile.SuccessfulTrades) / float64(profile.TotalTrades) * 100
	}

	// Calculate daily volume utilization
	dailyUtilization := float64(0)
	if !profile.DailyLimit.IsZero() {
		dailyUtilization = float64(profile.DailyVolume.Int64()) / float64(profile.DailyLimit.Int64()) * 100
	}

	metrics := &types.MitraPerformanceMetrics{
		MitraId:           mitraId,
		DhanPataName:      profile.DhanPataName,
		TrustScore:        profile.TrustScore,
		TotalTrades:       profile.TotalTrades,
		SuccessfulTrades:  profile.SuccessfulTrades,
		SuccessRate:       successRate,
		DailyVolume:       profile.DailyVolume,
		MonthlyVolume:     profile.MonthlyVolume,
		DailyLimit:        profile.DailyLimit,
		MonthlyLimit:      profile.MonthlyLimit,
		DailyUtilization:  dailyUtilization,
		ActiveEscrowCount: uint64(len(profile.ActiveEscrows)),
		LastActiveAt:      profile.LastActiveAt,
	}

	return metrics, nil
}

// UpdateMitraTradeStats updates mitra trading statistics
func (k Keeper) UpdateMitraTradeStats(ctx sdk.Context, mitraId string, tradeVolume sdk.Int, successful bool) error {
	profile, found := k.GetEnhancedMitraProfile(ctx, mitraId)
	if !found {
		return types.ErrMitraNotFound
	}

	// Update statistics
	profile.TotalTrades++
	if successful {
		profile.SuccessfulTrades++
	}

	// Update volume (assuming daily reset happens elsewhere)
	profile.DailyVolume = profile.DailyVolume.Add(tradeVolume)
	profile.MonthlyVolume = profile.MonthlyVolume.Add(tradeVolume)

	// Update last active time
	profile.LastActiveAt = ctx.BlockTime()

	// Recalculate trust score based on performance
	profile.TrustScore = k.calculateUpdatedTrustScore(profile)

	// Recalculate limits based on new trust score
	daily, monthly := types.CalculateMitraLimits(profile.MitraType, profile.TrustScore)
	profile.DailyLimit = daily
	profile.MonthlyLimit = monthly

	// Save updated profile
	k.SetEnhancedMitraProfile(ctx, profile)

	return nil
}

// calculateUpdatedTrustScore calculates updated trust score based on performance
func (k Keeper) calculateUpdatedTrustScore(profile types.EnhancedMitraProfile) int64 {
	baseScore := int64(50) // Base score

	// Success rate bonus (0-30 points)
	successRate := float64(0)
	if profile.TotalTrades > 0 {
		successRate = float64(profile.SuccessfulTrades) / float64(profile.TotalTrades)
	}
	successBonus := int64(successRate * 30)

	// Volume activity bonus (0-10 points)
	volumeBonus := int64(10) // Simplified - would be based on actual volume metrics

	// Trade count bonus (0-10 points)
	tradeBonus := int64(0)
	if profile.TotalTrades > 100 {
		tradeBonus = 10
	} else if profile.TotalTrades > 50 {
		tradeBonus = 5
	}

	newScore := baseScore + successBonus + volumeBonus + tradeBonus

	// Cap at 100
	if newScore > 100 {
		newScore = 100
	}

	return newScore
}