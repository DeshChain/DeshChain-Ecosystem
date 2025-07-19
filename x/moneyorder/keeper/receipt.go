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
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/deshchain/x/moneyorder/types"
)

// CreateMoneyOrder creates a new money order with UPI-style simplicity
func (k Keeper) CreateMoneyOrder(
	ctx sdk.Context,
	sender sdk.AccAddress,
	receiverUPI string,
	amount sdk.Coin,
	note string,
	orderType string,
	scheduledTime int64,
) (string, string, error) {
	params := k.GetParams(ctx)
	
	// Validate amount
	if amount.Amount.LT(params.MinMoneyOrderAmount) {
		return "", "", types.ErrOrderAmountTooLow
	}
	
	if amount.Amount.GT(params.MaxMoneyOrderAmount) {
		return "", "", types.ErrOrderAmountTooHigh
	}
	
	// Check daily limit
	if err := k.CheckDailyLimit(ctx, sender, amount.Amount); err != nil {
		return "", "", err
	}
	
	// Check KYC if required
	if params.IsLargeOrder(amount.Amount) {
		if err := k.ValidateKYC(ctx, sender); err != nil {
			return "", "", err
		}
	}
	
	// Resolve receiver address
	receiver, err := k.ResolveUPIAddress(ctx, receiverUPI)
	if err != nil {
		return "", "", err
	}
	
	// Calculate fees
	applicableDiscounts := k.calculateDiscounts(ctx, sender, amount.Denom)
	fees := params.CalculateTradingFee(amount.Amount, orderType, applicableDiscounts)
	feeCoin := sdk.NewCoin(amount.Denom, fees)
	
	// Check sender has sufficient balance
	totalRequired := amount.Add(feeCoin)
	senderBalance := k.bankKeeper.GetBalance(ctx, sender, amount.Denom)
	if senderBalance.IsLT(totalRequired) {
		return "", "", fmt.Errorf("insufficient balance")
	}
	
	// Generate order ID and reference number
	orderId := k.GetNextOrderId(ctx)
	
	// Create receipt
	receipt := types.NewMoneyOrderReceipt(
		orderId,
		sender,
		receiver,
		amount,
		note,
	)
	
	// Set additional receipt fields
	receipt.TransactionType = orderType
	receipt.Fees = feeCoin
	receipt.TotalAmount = totalRequired
	receipt.Language = "en" // Default, can be customized
	receipt.CulturalQuote = k.GetCulturalQuote(ctx, receipt.Language)
	
	// Add festival greeting if applicable
	if k.IsFestivalPeriod(ctx) {
		receipt.FestivalGreeting = "Happy Festival! Enjoy special discounts."
	}
	
	// Handle different order types
	switch orderType {
	case "instant":
		// Execute immediate transfer
		if err := k.executeMoneyOrderTransfer(ctx, sender, receiver, amount, feeCoin); err != nil {
			return "", "", err
		}
		receipt.Status = types.OrderStatusCompleted
		receipt.CompletedAt = ctx.BlockTime()
		
	case "normal":
		// Queue for processing (simplified - in production would have batch processing)
		if err := k.executeMoneyOrderTransfer(ctx, sender, receiver, amount, feeCoin); err != nil {
			return "", "", err
		}
		receipt.Status = types.OrderStatusCompleted
		receipt.ProcessedAt = ctx.BlockTime()
		receipt.CompletedAt = ctx.BlockTime()
		
	case "scheduled":
		// Store for future execution
		receipt.Status = types.OrderStatusPending
		receipt.ExpiresAt = ctx.BlockTime().Add(24 * 3600 * 1e9) // 24 hours
		// In production, would store scheduled orders separately
		
	default:
		return "", "", fmt.Errorf("invalid order type")
	}
	
	// Store receipt
	k.SetMoneyOrderReceipt(ctx, receipt)
	
	// Store order by user
	k.AddUserOrder(ctx, sender, orderId)
	k.AddUserOrder(ctx, receiver, orderId)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMoneyOrderCreated,
			sdk.NewAttribute(types.AttributeKeyOrderId, orderId),
			sdk.NewAttribute(types.AttributeKeyReferenceNumber, receipt.ReferenceNumber),
			sdk.NewAttribute(types.AttributeKeySender, sender.String()),
			sdk.NewAttribute(types.AttributeKeyReceiver, receiver.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(types.AttributeKeyOrderType, orderType),
			sdk.NewAttribute(types.AttributeKeyFees, feeCoin.String()),
		),
	)
	
	// Call hooks
	if k.hooks != nil {
		k.hooks.AfterMoneyOrderCreated(ctx, orderId, sender, receiver, amount)
	}
	
	return orderId, receipt.ReferenceNumber, nil
}

// executeMoneyOrderTransfer handles the actual transfer of funds
func (k Keeper) executeMoneyOrderTransfer(
	ctx sdk.Context,
	sender sdk.AccAddress,
	receiver sdk.AccAddress,
	amount sdk.Coin,
	fees sdk.Coin,
) error {
	// Transfer amount from sender to receiver
	if err := k.bankKeeper.SendCoins(ctx, sender, receiver, sdk.NewCoins(amount)); err != nil {
		return err
	}
	
	// Collect fees
	if !fees.IsZero() {
		if err := k.CollectFees(ctx, fees, sender); err != nil {
			return err
		}
	}
	
	return nil
}

// SetMoneyOrderReceipt stores a money order receipt
func (k Keeper) SetMoneyOrderReceipt(ctx sdk.Context, receipt *types.MoneyOrderReceipt) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetMoneyOrderKey(receipt.OrderId)
	bz := k.cdc.MustMarshal(receipt)
	store.Set(key, bz)
	
	// Also store by reference number for easy lookup
	refKey := types.GetReceiptKey(receipt.ReferenceNumber)
	store.Set(refKey, []byte(receipt.OrderId))
}

// GetMoneyOrderReceipt retrieves a money order receipt by ID
func (k Keeper) GetMoneyOrderReceipt(ctx sdk.Context, orderId string) (*types.MoneyOrderReceipt, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetMoneyOrderKey(orderId)
	bz := store.Get(key)
	
	if bz == nil {
		return nil, false
	}
	
	var receipt types.MoneyOrderReceipt
	k.cdc.MustUnmarshal(bz, &receipt)
	return &receipt, true
}

// GetReceiptByReference retrieves a receipt by reference number
func (k Keeper) GetReceiptByReference(ctx sdk.Context, referenceNumber string) (*types.MoneyOrderReceipt, bool) {
	store := ctx.KVStore(k.storeKey)
	refKey := types.GetReceiptKey(referenceNumber)
	orderIdBz := store.Get(refKey)
	
	if orderIdBz == nil {
		return nil, false
	}
	
	return k.GetMoneyOrderReceipt(ctx, string(orderIdBz))
}

// AddUserOrder adds an order to a user's order list
func (k Keeper) AddUserOrder(ctx sdk.Context, user sdk.AccAddress, orderId string) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetUserOrdersKey(user)
	
	// Get existing orders
	var orderIds []string
	bz := store.Get(key)
	if bz != nil {
		k.cdc.MustUnmarshal(bz, &orderIds)
	}
	
	// Add new order
	orderIds = append(orderIds, orderId)
	
	// Store updated list
	bz = k.cdc.MustMarshal(&orderIds)
	store.Set(key, bz)
}

// GetUserOrders retrieves all orders for a user
func (k Keeper) GetUserOrders(ctx sdk.Context, user sdk.AccAddress) []string {
	store := ctx.KVStore(k.storeKey)
	key := types.GetUserOrdersKey(user)
	
	var orderIds []string
	bz := store.Get(key)
	if bz != nil {
		k.cdc.MustUnmarshal(bz, &orderIds)
	}
	
	return orderIds
}

// GetUserReceipts retrieves all receipts for a user
func (k Keeper) GetUserReceipts(ctx sdk.Context, user sdk.AccAddress) []*types.MoneyOrderReceipt {
	orderIds := k.GetUserOrders(ctx, user)
	var receipts []*types.MoneyOrderReceipt
	
	for _, orderId := range orderIds {
		receipt, found := k.GetMoneyOrderReceipt(ctx, orderId)
		if found {
			receipts = append(receipts, receipt)
		}
	}
	
	return receipts
}

// calculateDiscounts calculates applicable discounts for a user
func (k Keeper) calculateDiscounts(ctx sdk.Context, user sdk.AccAddress, denom string) sdk.Dec {
	params := k.GetParams(ctx)
	
	// Check if cultural features are enabled
	if !params.EnableCulturalFeatures {
		return sdk.ZeroDec()
	}
	
	isFestival := k.IsFestivalPeriod(ctx)
	isVillagePoolMember := k.IsVillagePoolMember(ctx, user)
	isSeniorCitizen := k.IsSeniorCitizen(ctx, user)
	isCulturalToken := k.IsCulturalToken(denom)
	
	return params.GetCulturalDiscount(
		isFestival,
		isVillagePoolMember,
		isSeniorCitizen,
		isCulturalToken,
	)
}

// IsVillagePoolMember checks if user is a member of any village pool
func (k Keeper) IsVillagePoolMember(ctx sdk.Context, user sdk.AccAddress) bool {
	// Simplified check - in production would query village pool memberships
	store := ctx.KVStore(k.storeKey)
	key := append([]byte("village_member:"), user.Bytes()...)
	return store.Has(key)
}

// IsSeniorCitizen checks if user is a senior citizen (would integrate with KYC)
func (k Keeper) IsSeniorCitizen(ctx sdk.Context, user sdk.AccAddress) bool {
	// Placeholder - would check KYC data
	return false
}

// IsCulturalToken checks if a token is a cultural token
func (k Keeper) IsCulturalToken(denom string) bool {
	culturalTokens := map[string]bool{
		"uheritage": true,
		"uwisdom":   true,
		"ufestival": true,
		"unamo":     true, // NAMO is also a cultural token
	}
	return culturalTokens[denom]
}

// ProcessScheduledOrders processes all scheduled orders that are due
func (k Keeper) ProcessScheduledOrders(ctx sdk.Context) {
	// This would be called in BeginBlock
	// Simplified implementation - in production would have proper queue
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixMoneyOrder)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var receipt types.MoneyOrderReceipt
		k.cdc.MustUnmarshal(iterator.Value(), &receipt)
		
		// Process pending scheduled orders
		if receipt.Status == types.OrderStatusPending && 
		   receipt.TransactionType == "scheduled" &&
		   ctx.BlockTime().After(receipt.CreatedAt) {
			
			// Execute the transfer
			err := k.executeMoneyOrderTransfer(
				ctx,
				receipt.SenderAddress,
				receipt.ReceiverAddress,
				receipt.Amount,
				receipt.Fees,
			)
			
			if err == nil {
				receipt.Status = types.OrderStatusCompleted
				receipt.CompletedAt = ctx.BlockTime()
			} else {
				receipt.Status = types.OrderStatusCancelled
				receipt.StatusMessage = err.Error()
			}
			
			k.SetMoneyOrderReceipt(ctx, &receipt)
		}
	}
}

// GetMoneyOrderStatistics returns statistics for money orders
func (k Keeper) GetMoneyOrderStatistics(ctx sdk.Context) map[string]interface{} {
	stats := make(map[string]interface{})
	
	// Count orders by status
	statusCount := make(map[string]int)
	totalVolume := sdk.ZeroInt()
	
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixMoneyOrder)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var receipt types.MoneyOrderReceipt
		k.cdc.MustUnmarshal(iterator.Value(), &receipt)
		
		statusCount[receipt.Status]++
		if receipt.Status == types.OrderStatusCompleted {
			totalVolume = totalVolume.Add(receipt.Amount.Amount)
		}
	}
	
	stats["order_count_by_status"] = statusCount
	stats["total_volume"] = totalVolume.String()
	stats["total_orders"] = len(statusCount)
	
	return stats
}