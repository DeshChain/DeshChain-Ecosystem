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

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/moneyorder/types"
)

// SetEscrow stores an escrow
func (k Keeper) SetEscrow(ctx sdk.Context, escrow *types.Escrow) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetEscrowKey(escrow.EscrowId)
	value := k.cdc.MustMarshal(escrow)
	store.Set(key, value)
	
	// Index by order ID if present
	if escrow.OrderId != "" {
		k.SetEscrowOrderIndex(ctx, escrow.OrderId, escrow.EscrowId)
	}
	
	// Index by trade ID if present
	if escrow.TradeId != "" {
		k.SetEscrowTradeIndex(ctx, escrow.TradeId, escrow.EscrowId)
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

// GetEscrowByOrderID retrieves escrow by order ID
func (k Keeper) GetEscrowByOrderID(ctx sdk.Context, orderID string) (*types.Escrow, bool) {
	escrowID := k.GetEscrowOrderIndex(ctx, orderID)
	if escrowID == "" {
		return nil, false
	}
	return k.GetEscrow(ctx, escrowID)
}

// GetEscrowByTradeID retrieves escrow by trade ID
func (k Keeper) GetEscrowByTradeID(ctx sdk.Context, tradeID string) (*types.Escrow, bool) {
	escrowID := k.GetEscrowTradeIndex(ctx, tradeID)
	if escrowID == "" {
		return nil, false
	}
	return k.GetEscrow(ctx, escrowID)
}

// SetEscrowOrderIndex stores escrow ID by order ID
func (k Keeper) SetEscrowOrderIndex(ctx sdk.Context, orderID, escrowID string) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetEscrowOrderIndexKey(orderID)
	store.Set(key, []byte(escrowID))
}

// GetEscrowOrderIndex retrieves escrow ID by order ID
func (k Keeper) GetEscrowOrderIndex(ctx sdk.Context, orderID string) string {
	store := ctx.KVStore(k.storeKey)
	key := types.GetEscrowOrderIndexKey(orderID)
	value := store.Get(key)
	if value == nil {
		return ""
	}
	return string(value)
}

// SetEscrowTradeIndex stores escrow ID by trade ID
func (k Keeper) SetEscrowTradeIndex(ctx sdk.Context, tradeID, escrowID string) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetEscrowTradeIndexKey(tradeID)
	store.Set(key, []byte(escrowID))
}

// GetEscrowTradeIndex retrieves escrow ID by trade ID
func (k Keeper) GetEscrowTradeIndex(ctx sdk.Context, tradeID string) string {
	store := ctx.KVStore(k.storeKey)
	key := types.GetEscrowTradeIndexKey(tradeID)
	value := store.Get(key)
	if value == nil {
		return ""
	}
	return string(value)
}

// SetDispute stores a dispute
func (k Keeper) SetDispute(ctx sdk.Context, dispute *types.Dispute) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetDisputeKey(dispute.DisputeId)
	value := k.cdc.MustMarshal(dispute)
	store.Set(key, value)
	
	// Index by escrow ID
	k.SetDisputeEscrowIndex(ctx, dispute.EscrowId, dispute.DisputeId)
}

// GetDispute retrieves a dispute by ID
func (k Keeper) GetDispute(ctx sdk.Context, disputeID string) (*types.Dispute, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetDisputeKey(disputeID)
	value := store.Get(key)
	if value == nil {
		return nil, false
	}
	
	var dispute types.Dispute
	k.cdc.MustUnmarshal(value, &dispute)
	return &dispute, true
}

// SetDisputeEscrowIndex stores dispute ID by escrow ID
func (k Keeper) SetDisputeEscrowIndex(ctx sdk.Context, escrowID, disputeID string) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetDisputeEscrowIndexKey(escrowID)
	store.Set(key, []byte(disputeID))
}

// GetDisputeByEscrowID retrieves dispute by escrow ID
func (k Keeper) GetDisputeByEscrowID(ctx sdk.Context, escrowID string) (*types.Dispute, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetDisputeEscrowIndexKey(escrowID)
	value := store.Get(key)
	if value == nil {
		return nil, false
	}
	
	disputeID := string(value)
	return k.GetDispute(ctx, disputeID)
}

// AddToExpiryQueue adds escrow to expiry check queue
func (k Keeper) AddToExpiryQueue(ctx sdk.Context, escrowID string, expiresAt time.Time) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetEscrowExpiryQueueKey(expiresAt, escrowID)
	store.Set(key, []byte{1})
}

// GetExpiredEscrows returns escrows that have expired
func (k Keeper) GetExpiredEscrows(ctx sdk.Context, currentTime time.Time) []string {
	store := ctx.KVStore(k.storeKey)
	prefix := types.EscrowExpiryQueuePrefix
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()
	
	var expiredEscrows []string
	
	for ; iterator.Valid(); iterator.Next() {
		// Parse time from key
		key := iterator.Key()
		timeBytes := key[len(prefix):len(prefix)+8]
		escrowIDBytes := key[len(prefix)+8:]
		
		// Check if expired
		expiryTime := types.ParseTimeFromBytes(timeBytes)
		if currentTime.After(expiryTime) {
			expiredEscrows = append(expiredEscrows, string(escrowIDBytes))
			// Remove from queue
			store.Delete(key)
		}
	}
	
	return expiredEscrows
}

// GetEscrowStats returns escrow statistics
func (k Keeper) GetEscrowStats(ctx sdk.Context) *types.EscrowStats {
	stats := &types.EscrowStats{
		TotalVolume:       sdk.NewCoin(types.DefaultDenom, sdk.ZeroInt()),
		TotalFeesCollected: sdk.NewCoin(types.DefaultDenom, sdk.ZeroInt()),
	}
	
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.EscrowPrefix)
	defer iterator.Close()
	
	var totalDuration time.Duration
	var completedCount int64
	
	for ; iterator.Valid(); iterator.Next() {
		var escrow types.Escrow
		k.cdc.MustUnmarshal(iterator.Value(), &escrow)
		
		stats.TotalEscrows++
		
		switch escrow.Status {
		case types.EscrowStatus_ESCROW_STATUS_ACTIVE:
			stats.ActiveEscrows++
		case types.EscrowStatus_ESCROW_STATUS_RELEASED:
			stats.CompletedEscrows++
			completedCount++
			duration := escrow.ExpiresAt.Sub(escrow.CreatedAt)
			totalDuration += duration
			stats.TotalVolume = stats.TotalVolume.Add(escrow.Amount)
			stats.TotalFeesCollected = stats.TotalFeesCollected.Add(escrow.PlatformFee)
		case types.EscrowStatus_ESCROW_STATUS_REFUNDED:
			stats.RefundedEscrows++
		case types.EscrowStatus_ESCROW_STATUS_DISPUTED:
			stats.DisputedEscrows++
		}
	}
	
	// Calculate averages
	if completedCount > 0 {
		stats.AvgEscrowDurationHours = totalDuration.Hours() / float64(completedCount)
	}
	
	if stats.TotalEscrows > 0 {
		stats.DisputeRate = float64(stats.DisputedEscrows) / float64(stats.TotalEscrows) * 100
		stats.SuccessfulCompletionRate = float64(stats.CompletedEscrows) / float64(stats.TotalEscrows) * 100
	}
	
	return stats
}