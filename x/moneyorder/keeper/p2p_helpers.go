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
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/deshchain/x/moneyorder/types"
)

// P2P Helper Functions

// generateRandomString generates a cryptographically secure random string
func (k Keeper) generateRandomString(length int) string {
	bytes := make([]byte, length/2)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)[:length]
}

// generateEscrowID generates a unique escrow ID
func (k Keeper) generateEscrowID(ctx sdk.Context) string {
	return fmt.Sprintf("ESC-%d-%s", ctx.BlockHeight(), k.generateRandomString(8))
}

// GetEscrowByOrderID retrieves escrow by order ID
func (k Keeper) GetEscrowByOrderID(ctx sdk.Context, orderID string) (*types.Escrow, bool) {
	store := ctx.KVStore(k.storeKey)
	indexKey := types.GetEscrowOrderIndexKey(orderID)
	escrowID := store.Get(indexKey)
	if escrowID == nil {
		return nil, false
	}
	
	return k.GetEscrow(ctx, string(escrowID))
}

// SetEscrow stores an escrow record
func (k Keeper) SetEscrow(ctx sdk.Context, escrow *types.Escrow) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetEscrowKey(escrow.EscrowId)
	value := k.cdc.MustMarshal(escrow)
	store.Set(key, value)
	
	// Set indexes
	if escrow.OrderId != "" {
		indexKey := types.GetEscrowOrderIndexKey(escrow.OrderId)
		store.Set(indexKey, []byte(escrow.EscrowId))
	}
	
	if escrow.TradeId != "" {
		indexKey := types.GetEscrowTradeIndexKey(escrow.TradeId)
		store.Set(indexKey, []byte(escrow.EscrowId))
	}
	
	// Add to expiry queue
	if !escrow.ExpiresAt.IsZero() {
		k.AddToEscrowExpiryQueue(ctx, escrow.EscrowId, escrow.ExpiresAt)
	}
}

// GetEscrow retrieves an escrow by ID
func (k Keeper) GetEscrow(ctx sdk.Context, escrowID string) (*types.Escrow, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetEscrowKey(escrowID)
	value := store.Get(key)
	if value == nil {
		return nil, false
	}
	
	var escrow types.Escrow
	k.cdc.MustUnmarshal(value, &escrow)
	return &escrow, true
}

// AddToEscrowExpiryQueue adds escrow to expiry queue
func (k Keeper) AddToEscrowExpiryQueue(ctx sdk.Context, escrowID string, expiresAt time.Time) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetEscrowExpiryQueueKey(expiresAt, escrowID)
	store.Set(key, []byte{1})
}

// getCommonPaymentMethods finds common payment methods between two lists
func (k Keeper) getCommonPaymentMethods(methods1, methods2 []types.PaymentMethod) []types.PaymentMethod {
	var common []types.PaymentMethod
	
	for _, m1 := range methods1 {
		for _, m2 := range methods2 {
			if k.arePaymentMethodsCompatible(m1, m2) {
				common = append(common, m1)
				break
			}
		}
	}
	
	return common
}

// arePaymentMethodsCompatible checks if two payment methods are compatible
func (k Keeper) arePaymentMethodsCompatible(m1, m2 types.PaymentMethod) bool {
	// Same method type required
	if m1.MethodType != m2.MethodType {
		return false
	}
	
	// For UPI, any provider works
	if m1.MethodType == "UPI" {
		return true
	}
	
	// For bank transfers, compatibility is automatic
	if m1.MethodType == "IMPS" || m1.MethodType == "NEFT" || m1.MethodType == "RTGS" {
		return true
	}
	
	// For specific providers, must match
	return m1.Provider == m2.Provider
}

// calculateAmountOverlap calculates the overlap between two order amounts
func (k Keeper) calculateAmountOverlap(order1, order2 *types.P2POrder) float64 {
	// Get ranges for both orders
	min1 := order1.MinAmount
	if min1.IsZero() {
		min1 = order1.Amount
	}
	max1 := order1.MaxAmount
	if max1.IsZero() {
		max1 = order1.Amount
	}
	
	min2 := order2.MinAmount
	if min2.IsZero() {
		min2 = order2.Amount
	}
	max2 := order2.MaxAmount
	if max2.IsZero() {
		max2 = order2.Amount
	}
	
	// Check if ranges overlap
	if min1.IsGT(max2) || min2.IsGT(max1) {
		return 0
	}
	
	// Calculate overlap size
	overlapMin := sdk.MaxInt(min1.Amount, min2.Amount)
	overlapMax := sdk.MinInt(max1.Amount, max2.Amount)
	overlapSize := overlapMax.Sub(overlapMin)
	
	// Calculate average range size
	range1 := max1.Amount.Sub(min1.Amount)
	range2 := max2.Amount.Sub(min2.Amount)
	avgRange := range1.Add(range2).QuoRaw(2)
	
	if avgRange.IsZero() {
		// Both are fixed amounts
		if order1.Amount.Equal(order2.Amount) {
			return 1.0
		}
		return 0
	}
	
	// Return overlap percentage
	return float64(overlapSize.Int64()) / float64(avgRange.Int64())
}

// determineTradeAmount calculates the actual trade amount for matching orders
func (k Keeper) determineTradeAmount(buyOrder, sellOrder *types.P2POrder) sdk.Coin {
	// Get effective ranges
	buyMin := buyOrder.MinAmount
	if buyMin.IsZero() {
		buyMin = buyOrder.Amount
	}
	buyMax := buyOrder.MaxAmount
	if buyMax.IsZero() {
		buyMax = buyOrder.Amount
	}
	
	sellMin := sellOrder.MinAmount
	if sellMin.IsZero() {
		sellMin = sellOrder.Amount
	}
	sellMax := sellOrder.MaxAmount
	if sellMax.IsZero() {
		sellMax = sellOrder.Amount
	}
	
	// Find overlap
	tradeMin := sdk.MaxInt(buyMin.Amount, sellMin.Amount)
	tradeMax := sdk.MinInt(buyMax.Amount, sellMax.Amount)
	
	if tradeMin.GT(tradeMax) {
		return sdk.NewCoin(types.DefaultDenom, sdk.ZeroInt())
	}
	
	// Use the maximum possible amount within overlap
	return sdk.NewCoin(types.DefaultDenom, tradeMax)
}

// calculateFiatAmount calculates fiat amount based on NAMO amount and rate
func (k Keeper) calculateFiatAmount(namoAmount, referenceFiat, referenceNamo sdk.Coin) sdk.Coin {
	if referenceNamo.IsZero() {
		return sdk.NewCoin(referenceFiat.Denom, sdk.ZeroInt())
	}
	
	// Calculate rate: fiat per NAMO
	rate := sdk.NewDecFromInt(referenceFiat.Amount).Quo(sdk.NewDecFromInt(referenceNamo.Amount))
	
	// Calculate fiat amount
	fiatAmount := rate.MulInt(namoAmount.Amount).TruncateInt()
	
	return sdk.NewCoin(referenceFiat.Denom, fiatAmount)
}

// selectPaymentMethod selects the best common payment method
func (k Keeper) selectPaymentMethod(buyerMethods, sellerMethods []types.PaymentMethod) *types.PaymentMethod {
	// Preference order: UPI > IMPS > NEFT > Others
	preferenceOrder := []string{"UPI", "IMPS", "NEFT", "RTGS", "CASH"}
	
	for _, prefMethod := range preferenceOrder {
		for _, bMethod := range buyerMethods {
			if bMethod.MethodType != prefMethod {
				continue
			}
			
			for _, sMethod := range sellerMethods {
				if k.arePaymentMethodsCompatible(bMethod, sMethod) {
					// Return buyer's method (contains their account info)
					return &bMethod
				}
			}
		}
	}
	
	// Fallback to first compatible method
	common := k.getCommonPaymentMethods(buyerMethods, sellerMethods)
	if len(common) > 0 {
		return &common[0]
	}
	
	return nil
}

// ProcessExpiredEscrows processes all expired escrows for automatic refund
func (k Keeper) ProcessExpiredEscrows(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	endKey := append(types.EscrowExpiryQueuePrefix, sdk.FormatTimeBytes(ctx.BlockTime())...)
	iterator := store.Iterator(types.EscrowExpiryQueuePrefix, endKey)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		// Extract escrow ID from key
		key := iterator.Key()
		timeLen := len(types.EscrowExpiryQueuePrefix) + 8 // 8 bytes for time
		if len(key) > timeLen {
			escrowID := string(key[timeLen:])
			
			// Process refund
			if err := k.RefundEscrow(ctx, escrowID, "expired"); err != nil {
				k.Logger(ctx).Error("Failed to refund expired escrow", "escrow_id", escrowID, "error", err)
			}
			
			// Remove from queue
			store.Delete(key)
		}
	}
}

// RefundEscrow processes escrow refund
func (k Keeper) RefundEscrow(ctx sdk.Context, escrowID string, reason string) error {
	escrow, found := k.GetEscrow(ctx, escrowID)
	if !found {
		return types.ErrEscrowNotFound
	}
	
	if escrow.Status != types.EscrowStatus_ESCROW_STATUS_ACTIVE {
		return fmt.Errorf("escrow not in active status")
	}
	
	// Get escrow module account
	escrowAddr := k.accountKeeper.GetModuleAddress(types.EscrowModuleName)
	refundAddr, err := sdk.AccAddressFromBech32(escrow.RefundTo)
	if err != nil {
		return err
	}
	
	// Refund full amount including platform fees
	totalRefund := escrow.Amount.Add(escrow.PlatformFee)
	if err := k.bankKeeper.SendCoins(ctx, escrowAddr, refundAddr, sdk.NewCoins(totalRefund)); err != nil {
		return err
	}
	
	// Update escrow status
	escrow.Status = types.EscrowStatus_ESCROW_STATUS_REFUNDED
	k.SetEscrow(ctx, escrow)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeEscrowRefunded,
			sdk.NewAttribute("escrow_id", escrowID),
			sdk.NewAttribute("amount", totalRefund.String()),
			sdk.NewAttribute("reason", reason),
		),
	)
	
	return nil
}

// CreateDispute creates a new dispute for a trade
func (k Keeper) CreateDispute(ctx sdk.Context, tradeID string, raisedBy sdk.AccAddress, reason string) (*types.Dispute, error) {
	trade, found := k.GetP2PTrade(ctx, tradeID)
	if !found {
		return nil, types.ErrTradeNotFound
	}
	
	// Check if user is party to the trade
	if raisedBy.String() != trade.Buyer && raisedBy.String() != trade.Seller {
		return nil, types.ErrUnauthorized
	}
	
	// Determine who dispute is against
	raisedAgainst := trade.Seller
	if raisedBy.String() == trade.Seller {
		raisedAgainst = trade.Buyer
	}
	
	// Create dispute
	dispute := &types.Dispute{
		DisputeId:     k.generateDisputeID(ctx),
		EscrowId:      "", // Will be set if escrow exists
		TradeId:       tradeID,
		RaisedBy:      raisedBy.String(),
		RaisedAgainst: raisedAgainst,
		Reason:        reason,
		Status:        types.DisputeStatus_DISPUTE_STATUS_OPEN,
		CreatedAt:     ctx.BlockTime(),
	}
	
	// Find associated escrow
	if escrow, found := k.GetEscrowByTradeID(ctx, tradeID); found {
		dispute.EscrowId = escrow.EscrowId
		escrow.Status = types.EscrowStatus_ESCROW_STATUS_DISPUTED
		k.SetEscrow(ctx, escrow)
	}
	
	// Update trade status
	trade.Status = types.P2PTradeStatus_TRADE_STATUS_DISPUTED
	trade.Disputed = true
	trade.DisputeReason = reason
	k.SetP2PTrade(ctx, trade)
	
	// Store dispute
	k.SetDispute(ctx, dispute)
	
	// Update user stats
	k.RecordDispute(ctx, trade.Buyer, false) // Will be updated when resolved
	k.RecordDispute(ctx, trade.Seller, false)
	
	return dispute, nil
}

// SetDispute stores a dispute
func (k Keeper) SetDispute(ctx sdk.Context, dispute *types.Dispute) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetDisputeKey(dispute.DisputeId)
	value := k.cdc.MustMarshal(dispute)
	store.Set(key, value)
	
	// Set indexes
	if dispute.EscrowId != "" {
		indexKey := types.GetDisputeEscrowIndexKey(dispute.EscrowId)
		store.Set(indexKey, []byte(dispute.DisputeId))
	}
}

// generateDisputeID generates a unique dispute ID
func (k Keeper) generateDisputeID(ctx sdk.Context) string {
	return fmt.Sprintf("DISP-%d-%s", ctx.BlockHeight(), k.generateRandomString(8))
}

// GetEscrowByTradeID retrieves escrow by trade ID
func (k Keeper) GetEscrowByTradeID(ctx sdk.Context, tradeID string) (*types.Escrow, bool) {
	store := ctx.KVStore(k.storeKey)
	indexKey := types.GetEscrowTradeIndexKey(tradeID)
	escrowID := store.Get(indexKey)
	if escrowID == nil {
		return nil, false
	}
	
	return k.GetEscrow(ctx, string(escrowID))
}