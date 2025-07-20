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
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/deshchain/x/moneyorder/types"
)

// P2P Order Storage and Retrieval Functions

// SetP2POrder stores a P2P order
func (k Keeper) SetP2POrder(ctx sdk.Context, order *types.P2POrder) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetP2POrderKey(order.OrderId)
	value := k.cdc.MustMarshal(order)
	store.Set(key, value)
	
	// Update indexes
	k.SetP2POrderTypeIndex(ctx, order)
	k.SetP2POrderPostalIndex(ctx, order)
	k.SetP2POrderUserIndex(ctx, order)
	
	// Add to matching engine if active
	if order.Status == types.P2POrderStatus_P2P_STATUS_ACTIVE {
		if k.matchingEngine != nil {
			k.matchingEngine.AddOrderToBook(order)
		}
	}
}

// GetP2POrder retrieves a P2P order by ID
func (k Keeper) GetP2POrder(ctx sdk.Context, orderID string) (*types.P2POrder, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetP2POrderKey(orderID)
	value := store.Get(key)
	if value == nil {
		return nil, false
	}
	
	var order types.P2POrder
	k.cdc.MustUnmarshal(value, &order)
	return &order, true
}

// DeleteP2POrder removes a P2P order
func (k Keeper) DeleteP2POrder(ctx sdk.Context, orderID string) {
	order, found := k.GetP2POrder(ctx, orderID)
	if !found {
		return
	}
	
	store := ctx.KVStore(k.storeKey)
	key := types.GetP2POrderKey(orderID)
	store.Delete(key)
	
	// Remove from indexes
	k.RemoveP2POrderIndexes(ctx, order)
	
	// Remove from matching engine
	if k.matchingEngine != nil {
		k.matchingEngine.RemoveOrderFromBook(orderID)
	}
}

// SetP2POrderTypeIndex indexes order by type
func (k Keeper) SetP2POrderTypeIndex(ctx sdk.Context, order *types.P2POrder) {
	store := ctx.KVStore(k.storeKey)
	key := k.getP2POrderTypeIndexKey(order.OrderType, order.OrderId)
	store.Set(key, []byte{1})
}

// SetP2POrderPostalIndex indexes order by postal code
func (k Keeper) SetP2POrderPostalIndex(ctx sdk.Context, order *types.P2POrder) {
	store := ctx.KVStore(k.storeKey)
	
	// Index by postal code
	postalKey := k.getP2POrderPostalIndexKey(order.PostalCode, order.OrderId)
	store.Set(postalKey, []byte{1})
	
	// Index by district
	districtKey := k.getP2POrderDistrictIndexKey(order.District, order.OrderId)
	store.Set(districtKey, []byte{1})
	
	// Index by state
	stateKey := k.getP2POrderStateIndexKey(order.State, order.OrderId)
	store.Set(stateKey, []byte{1})
}

// SetP2POrderUserIndex indexes order by user
func (k Keeper) SetP2POrderUserIndex(ctx sdk.Context, order *types.P2POrder) {
	store := ctx.KVStore(k.storeKey)
	key := k.getP2POrderUserIndexKey(order.Creator, order.OrderId)
	store.Set(key, []byte{1})
}

// RemoveP2POrderIndexes removes all indexes for an order
func (k Keeper) RemoveP2POrderIndexes(ctx sdk.Context, order *types.P2POrder) {
	store := ctx.KVStore(k.storeKey)
	
	// Remove type index
	typeKey := k.getP2POrderTypeIndexKey(order.OrderType, order.OrderId)
	store.Delete(typeKey)
	
	// Remove postal indexes
	postalKey := k.getP2POrderPostalIndexKey(order.PostalCode, order.OrderId)
	store.Delete(postalKey)
	
	districtKey := k.getP2POrderDistrictIndexKey(order.District, order.OrderId)
	store.Delete(districtKey)
	
	stateKey := k.getP2POrderStateIndexKey(order.State, order.OrderId)
	store.Delete(stateKey)
	
	// Remove user index
	userKey := k.getP2POrderUserIndexKey(order.Creator, order.OrderId)
	store.Delete(userKey)
}

// GetP2POrdersByType returns all orders of a specific type
func (k Keeper) GetP2POrdersByType(ctx sdk.Context, orderType types.OrderType) []*types.P2POrder {
	store := ctx.KVStore(k.storeKey)
	prefix := k.getP2POrderTypeIndexPrefix(orderType)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()
	
	var orders []*types.P2POrder
	for ; iterator.Valid(); iterator.Next() {
		// Extract order ID from key
		orderID := k.extractOrderIDFromIndexKey(iterator.Key(), prefix)
		if order, found := k.GetP2POrder(ctx, orderID); found {
			orders = append(orders, order)
		}
	}
	
	return orders
}

// GetP2POrdersByPostalCode returns orders in a postal code
func (k Keeper) GetP2POrdersByPostalCode(ctx sdk.Context, postalCode string) []*types.P2POrder {
	store := ctx.KVStore(k.storeKey)
	prefix := k.getP2POrderPostalIndexPrefix(postalCode)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()
	
	var orders []*types.P2POrder
	for ; iterator.Valid(); iterator.Next() {
		orderID := k.extractOrderIDFromIndexKey(iterator.Key(), prefix)
		if order, found := k.GetP2POrder(ctx, orderID); found {
			orders = append(orders, order)
		}
	}
	
	return orders
}

// GetP2POrdersByDistrict returns orders in a district
func (k Keeper) GetP2POrdersByDistrict(ctx sdk.Context, district string) []*types.P2POrder {
	store := ctx.KVStore(k.storeKey)
	prefix := k.getP2POrderDistrictIndexPrefix(district)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()
	
	var orders []*types.P2POrder
	for ; iterator.Valid(); iterator.Next() {
		orderID := k.extractOrderIDFromIndexKey(iterator.Key(), prefix)
		if order, found := k.GetP2POrder(ctx, orderID); found {
			orders = append(orders, order)
		}
	}
	
	return orders
}

// GetP2POrdersByState returns orders in a state
func (k Keeper) GetP2POrdersByState(ctx sdk.Context, state string) []*types.P2POrder {
	store := ctx.KVStore(k.storeKey)
	prefix := k.getP2POrderStateIndexPrefix(state)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()
	
	var orders []*types.P2POrder
	for ; iterator.Valid(); iterator.Next() {
		orderID := k.extractOrderIDFromIndexKey(iterator.Key(), prefix)
		if order, found := k.GetP2POrder(ctx, orderID); found {
			orders = append(orders, order)
		}
	}
	
	return orders
}

// GetP2POrdersByUser returns orders created by a user
func (k Keeper) GetP2POrdersByUser(ctx sdk.Context, userAddr string) []*types.P2POrder {
	store := ctx.KVStore(k.storeKey)
	prefix := k.getP2POrderUserIndexPrefix(userAddr)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()
	
	var orders []*types.P2POrder
	for ; iterator.Valid(); iterator.Next() {
		orderID := k.extractOrderIDFromIndexKey(iterator.Key(), prefix)
		if order, found := k.GetP2POrder(ctx, orderID); found {
			orders = append(orders, order)
		}
	}
	
	return orders
}

// GetP2POrdersByTypeAndArea returns orders of a type within an area
func (k Keeper) GetP2POrdersByTypeAndArea(ctx sdk.Context, orderType types.OrderType, postalCode string, maxDistanceKm int32) []*types.P2POrder {
	var orders []*types.P2POrder
	orderMap := make(map[string]bool)
	
	// Get orders by type
	typeOrders := k.GetP2POrdersByType(ctx, orderType)
	
	// Filter by area
	for _, order := range typeOrders {
		// Skip if not active
		if order.Status != types.P2POrderStatus_P2P_STATUS_ACTIVE {
			continue
		}
		
		// Check distance
		if k.isWithinDistance(postalCode, order.PostalCode, maxDistanceKm) {
			if !orderMap[order.OrderId] {
				orders = append(orders, order)
				orderMap[order.OrderId] = true
			}
		}
	}
	
	return orders
}

// GetActiveP2POrders returns all active P2P orders
func (k Keeper) GetActiveP2POrders(ctx sdk.Context) []*types.P2POrder {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.P2POrderPrefix)
	defer iterator.Close()
	
	var orders []*types.P2POrder
	for ; iterator.Valid(); iterator.Next() {
		var order types.P2POrder
		k.cdc.MustUnmarshal(iterator.Value(), &order)
		
		if order.Status == types.P2POrderStatus_P2P_STATUS_ACTIVE {
			orders = append(orders, &order)
		}
	}
	
	return orders
}

// GetExpiredP2POrders returns orders that have expired
func (k Keeper) GetExpiredP2POrders(ctx sdk.Context, currentTime time.Time) []*types.P2POrder {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.P2POrderPrefix)
	defer iterator.Close()
	
	var orders []*types.P2POrder
	for ; iterator.Valid(); iterator.Next() {
		var order types.P2POrder
		k.cdc.MustUnmarshal(iterator.Value(), &order)
		
		if order.Status == types.P2POrderStatus_P2P_STATUS_ACTIVE && currentTime.After(order.ExpiresAt) {
			orders = append(orders, &order)
		}
	}
	
	return orders
}

// Index key generation functions

func (k Keeper) getP2POrderTypeIndexKey(orderType types.OrderType, orderID string) []byte {
	return append(k.getP2POrderTypeIndexPrefix(orderType), []byte(orderID)...)
}

func (k Keeper) getP2POrderTypeIndexPrefix(orderType types.OrderType) []byte {
	return append(types.P2POrderTypeIndexPrefix, byte(orderType))
}

func (k Keeper) getP2POrderPostalIndexKey(postalCode, orderID string) []byte {
	return append(k.getP2POrderPostalIndexPrefix(postalCode), []byte(orderID)...)
}

func (k Keeper) getP2POrderPostalIndexPrefix(postalCode string) []byte {
	return append([]byte("p2p-postal-"), []byte(postalCode)...)
}

func (k Keeper) getP2POrderDistrictIndexKey(district, orderID string) []byte {
	return append(k.getP2POrderDistrictIndexPrefix(district), []byte(orderID)...)
}

func (k Keeper) getP2POrderDistrictIndexPrefix(district string) []byte {
	return append([]byte("p2p-district-"), []byte(district)...)
}

func (k Keeper) getP2POrderStateIndexKey(state, orderID string) []byte {
	return append(k.getP2POrderStateIndexPrefix(state), []byte(orderID)...)
}

func (k Keeper) getP2POrderStateIndexPrefix(state string) []byte {
	return append([]byte("p2p-state-"), []byte(state)...)
}

func (k Keeper) getP2POrderUserIndexKey(userAddr, orderID string) []byte {
	return append(k.getP2POrderUserIndexPrefix(userAddr), []byte(orderID)...)
}

func (k Keeper) getP2POrderUserIndexPrefix(userAddr string) []byte {
	return append([]byte("p2p-user-"), []byte(userAddr)...)
}

func (k Keeper) extractOrderIDFromIndexKey(key, prefix []byte) string {
	return string(key[len(prefix):])
}

// P2P Trade Storage

// SetP2PTrade stores a P2P trade
func (k Keeper) SetP2PTrade(ctx sdk.Context, trade *types.P2PTrade) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetP2PTradeKey(trade.TradeId)
	value := k.cdc.MustMarshal(trade)
	store.Set(key, value)
	
	// Update indexes
	k.SetP2PTradeUserIndex(ctx, trade)
}

// GetP2PTrade retrieves a P2P trade by ID
func (k Keeper) GetP2PTrade(ctx sdk.Context, tradeID string) (*types.P2PTrade, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetP2PTradeKey(tradeID)
	value := store.Get(key)
	if value == nil {
		return nil, false
	}
	
	var trade types.P2PTrade
	k.cdc.MustUnmarshal(value, &trade)
	return &trade, true
}

// SetP2PTradeUserIndex indexes trade by users
func (k Keeper) SetP2PTradeUserIndex(ctx sdk.Context, trade *types.P2PTrade) {
	store := ctx.KVStore(k.storeKey)
	
	// Index by buyer
	buyerKey := k.getP2PTradeUserIndexKey(trade.Buyer, trade.TradeId)
	store.Set(buyerKey, []byte{1})
	
	// Index by seller
	sellerKey := k.getP2PTradeUserIndexKey(trade.Seller, trade.TradeId)
	store.Set(sellerKey, []byte{1})
}

func (k Keeper) getP2PTradeUserIndexKey(userAddr, tradeID string) []byte {
	return append([]byte(fmt.Sprintf("p2p-trade-user-%s-", userAddr)), []byte(tradeID)...)
}

// GetP2PTradesByUser returns trades involving a user
func (k Keeper) GetP2PTradesByUser(ctx sdk.Context, userAddr string) []*types.P2PTrade {
	store := ctx.KVStore(k.storeKey)
	prefix := []byte(fmt.Sprintf("p2p-trade-user-%s-", userAddr))
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()
	
	var trades []*types.P2PTrade
	tradeMap := make(map[string]bool)
	
	for ; iterator.Valid(); iterator.Next() {
		tradeID := string(iterator.Key()[len(prefix):])
		if !tradeMap[tradeID] {
			if trade, found := k.GetP2PTrade(ctx, tradeID); found {
				trades = append(trades, trade)
				tradeMap[tradeID] = true
			}
		}
	}
	
	return trades
}

// Helper functions

// generateP2POrderID generates a unique P2P order ID
func (k Keeper) generateP2POrderID(ctx sdk.Context, creator string) string {
	return fmt.Sprintf("P2P-%d-%s-%s", ctx.BlockHeight(), creator[:8], k.generateRandomString(6))
}

// generateTradeID generates a unique trade ID
func (k Keeper) generateTradeID(ctx sdk.Context) string {
	return fmt.Sprintf("TRADE-%d-%s", ctx.BlockHeight(), k.generateRandomString(8))
}

// generateEscrowAddress generates an escrow address for an order
func (k Keeper) generateEscrowAddress(orderID string) sdk.AccAddress {
	// In production, would derive from module account
	// For now, use deterministic address based on order ID
	return sdk.AccAddress(sdk.Keccak256([]byte(fmt.Sprintf("escrow-%s", orderID)))[:20])
}

// AddToRefundQueue adds an order to the refund queue
func (k Keeper) AddToRefundQueue(ctx sdk.Context, orderID string, refundAt time.Time) {
	store := ctx.KVStore(k.storeKey)
	key := append(types.RefundQueuePrefix, sdk.FormatTimeBytes(refundAt)...)
	key = append(key, []byte(orderID)...)
	store.Set(key, []byte{1})
}

// GetRefundQueueItems returns orders due for refund
func (k Keeper) GetRefundQueueItems(ctx sdk.Context, currentTime time.Time) []string {
	store := ctx.KVStore(k.storeKey)
	endKey := append(types.RefundQueuePrefix, sdk.FormatTimeBytes(currentTime)...)
	iterator := store.Iterator(types.RefundQueuePrefix, endKey)
	defer iterator.Close()
	
	var orderIDs []string
	for ; iterator.Valid(); iterator.Next() {
		// Extract order ID from key
		key := iterator.Key()
		timeLen := len(types.RefundQueuePrefix) + 8 // 8 bytes for time
		if len(key) > timeLen {
			orderID := string(key[timeLen:])
			orderIDs = append(orderIDs, orderID)
		}
	}
	
	return orderIDs
}

// notifyP2PMatch sends notifications for a new match
func (k Keeper) notifyP2PMatch(ctx sdk.Context, trade *types.P2PTrade) {
	// In production, would send notifications via SMS/WhatsApp
	// For now, emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"p2p_match_notification",
			sdk.NewAttribute("trade_id", trade.TradeId),
			sdk.NewAttribute("buyer", trade.Buyer),
			sdk.NewAttribute("seller", trade.Seller),
			sdk.NewAttribute("notification_sent", "true"),
		),
	)
}

// updateUserStats updates user statistics after trade
func (k Keeper) updateUserStats(ctx sdk.Context, buyer, seller string, amount sdk.Coin) {
	// Update buyer stats
	buyerStats := k.GetUserStats(ctx, buyer)
	if buyerStats == nil {
		buyerStats = &types.UserP2PStats{
			Address: buyer,
		}
	}
	buyerStats.TotalTrades++
	buyerStats.SuccessfulTrades++
	buyerStats.TotalVolume = buyerStats.TotalVolume.Add(amount)
	buyerStats.LastTradeAt = ctx.BlockTime()
	k.SetUserStats(ctx, buyerStats)
	
	// Update seller stats
	sellerStats := k.GetUserStats(ctx, seller)
	if sellerStats == nil {
		sellerStats = &types.UserP2PStats{
			Address: seller,
		}
	}
	sellerStats.TotalTrades++
	sellerStats.SuccessfulTrades++
	sellerStats.TotalVolume = sellerStats.TotalVolume.Add(amount)
	sellerStats.LastTradeAt = ctx.BlockTime()
	k.SetUserStats(ctx, sellerStats)
}