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
	"sort"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/deshchain/x/moneyorder/types"
)

// P2P Matching Engine for NAMO <-> Fiat trades

// CreateP2POrder creates a new P2P buy/sell order
func (k Keeper) CreateP2POrder(ctx sdk.Context, msg *types.MsgCreateP2POrder) (*types.P2POrder, error) {
	// Validate postal code
	if !k.ValidateIndianPincode(msg.PostalCode) {
		return nil, types.ErrInvalidPostalCode
	}

	// Generate order ID
	orderID := k.generateP2POrderID(ctx, msg.Creator)

	// Create escrow address
	escrowAddr := k.generateEscrowAddress(orderID)

	// Determine escrow amount based on order type
	var escrowAmount sdk.Coin
	if msg.OrderType == types.OrderType_SELL_NAMO {
		// Seller deposits NAMO
		escrowAmount = msg.Amount
	} else {
		// Buyer deposits platform fee as commitment
		escrowAmount = k.calculatePlatformFee(msg.Amount)
	}

	// Create order
	order := &types.P2POrder{
		OrderId:      orderID,
		Creator:      msg.Creator,
		OrderType:    msg.OrderType,
		Amount:       msg.Amount,
		FiatAmount:   msg.FiatAmount,
		FiatCurrency: msg.FiatCurrency,
		PostalCode:   msg.PostalCode,
		District:     k.getDistrictFromPincode(msg.PostalCode),
		State:        k.GetStateName(msg.PostalCode),
		PaymentMethods: msg.PaymentMethods,
		MinAmount:    msg.MinAmount,
		MaxAmount:    msg.MaxAmount,
		CreatedAt:    ctx.BlockTime(),
		ExpiresAt:    ctx.BlockTime().Add(24 * time.Hour), // 24 hour expiry
		Status:       types.P2POrderStatus_P2P_STATUS_PENDING,
		EscrowAddress: escrowAddr.String(),
		EscrowAmount: escrowAmount,
		MaxDistanceKm: msg.MaxDistanceKm,
		PreferredLanguages: msg.PreferredLanguages,
		MinTrustScore: msg.MinTrustScore,
		RequireKyc:   msg.RequireKyc,
	}

	// Store order
	k.SetP2POrder(ctx, order)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeP2POrderCreated,
			sdk.NewAttribute(types.AttributeKeyOrderID, orderID),
			sdk.NewAttribute("order_type", order.OrderType.String()),
			sdk.NewAttribute("postal_code", order.PostalCode),
			sdk.NewAttribute("amount", order.Amount.String()),
		),
	)

	return order, nil
}

// DepositToEscrow deposits funds to escrow and activates order
func (k Keeper) DepositToEscrow(ctx sdk.Context, orderID string, depositor sdk.AccAddress, amount sdk.Coin) error {
	order, found := k.GetP2POrder(ctx, orderID)
	if !found {
		return types.ErrOrderNotFound
	}

	if order.Status != types.P2POrderStatus_P2P_STATUS_PENDING {
		return types.ErrInvalidOrderStatus
	}

	// Verify amount matches escrow requirement
	if !amount.Equal(order.EscrowAmount) {
		return types.ErrInvalidAmount
	}

	// Create escrow record
	escrow := &types.Escrow{
		EscrowId:    k.generateEscrowID(ctx),
		OrderId:     orderID,
		Depositor:   depositor.String(),
		Amount:      amount,
		PlatformFee: k.calculatePlatformFee(amount),
		Status:      types.EscrowStatus_ESCROW_STATUS_ACTIVE,
		CreatedAt:   ctx.BlockTime(),
		ExpiresAt:   order.ExpiresAt,
		RefundTo:    depositor.String(),
	}

	// Transfer funds to escrow
	escrowAddr, _ := sdk.AccAddressFromBech32(order.EscrowAddress)
	err := k.bankKeeper.SendCoins(ctx, depositor, escrowAddr, sdk.NewCoins(amount))
	if err != nil {
		return err
	}

	// Update order status
	order.Status = types.P2POrderStatus_P2P_STATUS_ACTIVE
	k.SetP2POrder(ctx, order)
	k.SetEscrow(ctx, escrow)

	// Start matching process
	go k.findP2PMatches(ctx, order)

	return nil
}

// findP2PMatches finds matching orders in the same postal area
func (k Keeper) findP2PMatches(ctx sdk.Context, order *types.P2POrder) {
	// Get all active orders of opposite type
	oppositeType := types.OrderType_BUY_NAMO
	if order.OrderType == types.OrderType_BUY_NAMO {
		oppositeType = types.OrderType_SELL_NAMO
	}

	matches := k.GetP2POrdersByTypeAndArea(ctx, oppositeType, order.PostalCode, order.MaxDistanceKm)

	// Score and sort matches
	scoredMatches := k.scoreP2PMatches(order, matches)
	sort.Slice(scoredMatches, func(i, j int) bool {
		return scoredMatches[i].score > scoredMatches[j].score
	})

	// Try to match with best candidate
	for _, match := range scoredMatches {
		if k.tryCreateP2PTrade(ctx, order, match.order) {
			return
		}
	}

	// No match found - set up auto-refund after expiry
	k.scheduleAutoRefund(ctx, order)
}

// P2PMatch represents a potential match with score
type P2PMatch struct {
	order *types.P2POrder
	score float64
}

// scoreP2PMatches scores potential matches
func (k Keeper) scoreP2PMatches(order *types.P2POrder, candidates []*types.P2POrder) []P2PMatch {
	var matches []P2PMatch

	for _, candidate := range candidates {
		score := k.calculateMatchScore(order, candidate)
		if score > 0 {
			matches = append(matches, P2PMatch{
				order: candidate,
				score: score,
			})
		}
	}

	return matches
}

// calculateMatchScore calculates compatibility score between two orders
func (k Keeper) calculateMatchScore(order1, order2 *types.P2POrder) float64 {
	score := 100.0

	// Check basic compatibility
	if !k.areOrdersCompatible(order1, order2) {
		return 0
	}

	// Distance penalty (same pincode = 100, same district = 80, same state = 60)
	if order1.PostalCode == order2.PostalCode {
		score += 20
	} else if order1.District == order2.District {
		score -= 10
	} else if order1.State == order2.State {
		score -= 20
	} else {
		score -= 40
	}

	// Payment method compatibility
	commonMethods := k.getCommonPaymentMethods(order1.PaymentMethods, order2.PaymentMethods)
	if len(commonMethods) == 0 {
		return 0
	}
	score += float64(len(commonMethods)) * 5

	// Language compatibility
	if k.hasCommonLanguage(order1.PreferredLanguages, order2.PreferredLanguages) {
		score += 10
	}

	// Trust score check
	user1Trust := k.GetUserTrustScore(order1.Creator)
	user2Trust := k.GetUserTrustScore(order2.Creator)
	
	if user1Trust < order2.MinTrustScore || user2Trust < order1.MinTrustScore {
		return 0
	}

	// Amount overlap
	amountOverlap := k.calculateAmountOverlap(order1, order2)
	if amountOverlap == 0 {
		return 0
	}
	score *= amountOverlap

	// Time bonus for older orders
	ageBonus := time.Since(order2.CreatedAt).Hours()
	score += math.Min(ageBonus*2, 20) // Max 20 points for age

	return score
}

// areOrdersCompatible checks basic compatibility
func (k Keeper) areOrdersCompatible(order1, order2 *types.P2POrder) bool {
	// Opposite order types
	if order1.OrderType == order2.OrderType {
		return false
	}

	// Both active
	if order1.Status != types.P2POrderStatus_P2P_STATUS_ACTIVE || 
	   order2.Status != types.P2POrderStatus_P2P_STATUS_ACTIVE {
		return false
	}

	// Same fiat currency
	if order1.FiatCurrency != order2.FiatCurrency {
		return false
	}

	// KYC requirements
	if order1.RequireKyc && !k.IsKYCVerified(order2.Creator) {
		return false
	}
	if order2.RequireKyc && !k.IsKYCVerified(order1.Creator) {
		return false
	}

	return true
}

// tryCreateP2PTrade attempts to create a trade between two orders
func (k Keeper) tryCreateP2PTrade(ctx sdk.Context, buyOrder, sellOrder *types.P2POrder) bool {
	// Determine trade amount
	tradeAmount := k.determineTradeAmount(buyOrder, sellOrder)
	if tradeAmount.IsZero() {
		return false
	}

	// Calculate fiat amount based on rate
	fiatAmount := k.calculateFiatAmount(tradeAmount, buyOrder.FiatAmount, buyOrder.Amount)

	// Select payment method
	paymentMethod := k.selectPaymentMethod(buyOrder.PaymentMethods, sellOrder.PaymentMethods)
	if paymentMethod == nil {
		return false
	}

	// Create trade
	trade := &types.P2PTrade{
		TradeId:       k.generateTradeID(ctx),
		BuyerOrderId:  buyOrder.OrderId,
		SellerOrderId: sellOrder.OrderId,
		Buyer:         buyOrder.Creator,
		Seller:        sellOrder.Creator,
		NamoAmount:    tradeAmount,
		FiatAmount:    fiatAmount,
		FiatCurrency:  buyOrder.FiatCurrency,
		PaymentMethod: *paymentMethod,
		CreatedAt:     ctx.BlockTime(),
		ExpiresAt:     ctx.BlockTime().Add(2 * time.Hour), // 2 hour timeout
		Status:        types.P2PTradeStatus_TRADE_STATUS_MATCHED,
		EscrowAddress: sellOrder.EscrowAddress, // Use seller's escrow
	}

	// Update order statuses
	buyOrder.Status = types.P2POrderStatus_P2P_STATUS_MATCHED
	sellOrder.Status = types.P2POrderStatus_P2P_STATUS_MATCHED
	k.SetP2POrder(ctx, buyOrder)
	k.SetP2POrder(ctx, sellOrder)
	k.SetP2PTrade(ctx, trade)

	// Send notifications
	k.notifyP2PMatch(ctx, trade)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeP2PTradeMatched,
			sdk.NewAttribute("trade_id", trade.TradeId),
			sdk.NewAttribute("buyer", trade.Buyer),
			sdk.NewAttribute("seller", trade.Seller),
			sdk.NewAttribute("amount", trade.NamoAmount.String()),
		),
	)

	return true
}

// CompleteP2PTrade completes a P2P trade after fiat payment confirmation
func (k Keeper) CompleteP2PTrade(ctx sdk.Context, tradeID string, confirmedBy sdk.AccAddress) error {
	trade, found := k.GetP2PTrade(ctx, tradeID)
	if !found {
		return types.ErrTradeNotFound
	}

	// Only seller can confirm payment received
	if confirmedBy.String() != trade.Seller {
		return types.ErrUnauthorized
	}

	if trade.Status != types.P2PTradeStatus_TRADE_STATUS_PAYMENT_PENDING {
		return types.ErrInvalidTradeStatus
	}

	// Release escrow to buyer
	escrowAddr, _ := sdk.AccAddressFromBech32(trade.EscrowAddress)
	buyerAddr, _ := sdk.AccAddressFromBech32(trade.Buyer)
	
	err := k.bankKeeper.SendCoins(ctx, escrowAddr, buyerAddr, sdk.NewCoins(trade.NamoAmount))
	if err != nil {
		return err
	}

	// Update trade status
	trade.Status = types.P2PTradeStatus_TRADE_STATUS_COMPLETED
	k.SetP2PTrade(ctx, trade)

	// Update user stats
	k.updateUserStats(ctx, trade.Buyer, trade.Seller, trade.NamoAmount)

	// Update order statuses
	buyOrder, _ := k.GetP2POrder(ctx, trade.BuyerOrderId)
	sellOrder, _ := k.GetP2POrder(ctx, trade.SellerOrderId)
	buyOrder.Status = types.P2POrderStatus_P2P_STATUS_COMPLETED
	sellOrder.Status = types.P2POrderStatus_P2P_STATUS_COMPLETED
	k.SetP2POrder(ctx, buyOrder)
	k.SetP2POrder(ctx, sellOrder)

	return nil
}

// RefundP2POrder refunds an expired or unmatched order
func (k Keeper) RefundP2POrder(ctx sdk.Context, orderID string) error {
	order, found := k.GetP2POrder(ctx, orderID)
	if !found {
		return types.ErrOrderNotFound
	}

	// Check if order can be refunded
	if order.Status != types.P2POrderStatus_P2P_STATUS_ACTIVE {
		return types.ErrInvalidOrderStatus
	}

	// Check if expired
	if ctx.BlockTime().Before(order.ExpiresAt) {
		return types.ErrOrderNotExpired
	}

	// Get escrow
	escrow, found := k.GetEscrowByOrderID(ctx, orderID)
	if !found {
		return types.ErrEscrowNotFound
	}

	// Process refund including platform fees
	escrowAddr, _ := sdk.AccAddressFromBech32(order.EscrowAddress)
	refundAddr, _ := sdk.AccAddressFromBech32(escrow.RefundTo)
	
	// Refund full amount including fees
	totalRefund := escrow.Amount.Add(escrow.PlatformFee)
	err := k.bankKeeper.SendCoins(ctx, escrowAddr, refundAddr, sdk.NewCoins(totalRefund))
	if err != nil {
		return err
	}

	// Update statuses
	order.Status = types.P2POrderStatus_P2P_STATUS_REFUNDED
	escrow.Status = types.EscrowStatus_ESCROW_STATUS_REFUNDED
	k.SetP2POrder(ctx, order)
	k.SetEscrow(ctx, escrow)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeP2POrderRefunded,
			sdk.NewAttribute("order_id", orderID),
			sdk.NewAttribute("refund_amount", totalRefund.String()),
			sdk.NewAttribute("reason", "no_match_found"),
		),
	)

	return nil
}

// SevaMitra Management

// RegisterSevaMitra registers a new cash-in/out seva mitra
func (k Keeper) RegisterSevaMitra(ctx sdk.Context, msg *types.MsgRegisterSevaMitra) (*types.SevaMitra, error) {
	// Validate postal code
	if !k.ValidateIndianPincode(msg.PostalCode) {
		return nil, types.ErrInvalidPostalCode
	}

	// Check if already registered
	if k.IsSevaMitraRegistered(ctx, msg.Address) {
		return nil, types.ErrSevaMitraAlreadyRegistered
	}

	// Create seva mitra
	sevaMitra := &types.SevaMitra{
		MitraId:              k.generateSevaMitraID(ctx),
		Address:              msg.Address,
		BusinessName:         msg.BusinessName,
		RegistrationNumber:   msg.RegistrationNumber,
		PostalCode:          msg.PostalCode,
		FullAddress:         msg.FullAddress,
		District:            k.getDistrictFromPincode(msg.PostalCode),
		State:               k.GetStateName(msg.PostalCode),
		Latitude:            msg.Latitude,
		Longitude:           msg.Longitude,
		Phone:               msg.Phone,
		Email:               msg.Email,
		Languages:           msg.Languages,
		OperatingHours:      msg.OperatingHours,
		Services:            msg.Services,
		DailyLimit:          msg.DailyLimit,
		PerTransactionLimit: msg.PerTransactionLimit,
		Status:              types.SevaMitraStatus_MITRA_STATUS_PENDING,
		KycVerified:         false,
		Stats:               &types.SevaMitraStats{},
		CommissionRate:      "2.5", // Default 2.5%
		CreatedAt:           ctx.BlockTime(),
	}

	// Store seva mitra
	k.SetSevaMitra(ctx, sevaMitra)

	return sevaMitra, nil
}

// FindNearbySevaMitras finds seva mitras near a postal code
func (k Keeper) FindNearbySevaMitras(ctx sdk.Context, postalCode string, maxDistanceKm int32, services []types.SevaMitraService) []*types.SevaMitra {
	var nearbySevaMitras []*types.SevaMitra

	// Get all seva mitras in the same district first
	district := k.getDistrictFromPincode(postalCode)
	sevaMitras := k.GetSevaMitrasByDistrict(ctx, district)

	for _, sevaMitra := range sevaMitras {
		// Check if seva mitra is active
		if sevaMitra.Status != types.SevaMitraStatus_MITRA_STATUS_ACTIVE {
			continue
		}

		// Check if seva mitra provides required services
		if !k.sevaMitraProvidesServices(sevaMitra, services) {
			continue
		}

		// Simple distance check based on postal code similarity
		if k.isWithinDistance(postalCode, sevaMitra.PostalCode, maxDistanceKm) {
			nearbySevaMitras = append(nearbySevaMitras, sevaMitra)
		}
	}

	// Sort by trust score
	sort.Slice(nearbySevaMitras, func(i, j int) bool {
		return k.calculateSevaMitraTrustScore(nearbySevaMitras[i]) > k.calculateSevaMitraTrustScore(nearbySevaMitras[j])
	})

	return nearbySevaMitras
}

// Helper functions

func (k Keeper) calculatePlatformFee(amount sdk.Coin) sdk.Coin {
	// 0.5% platform fee
	fee := amount.Amount.MulRaw(5).QuoRaw(1000)
	return sdk.NewCoin(amount.Denom, fee)
}

func (k Keeper) getDistrictFromPincode(pincode string) string {
	// Simplified - in production would use actual mapping
	return fmt.Sprintf("District-%s", pincode[:3])
}

func (k Keeper) scheduleAutoRefund(ctx sdk.Context, order *types.P2POrder) {
	// In production, would use a proper scheduler
	// For now, mark for refund check
	k.AddToRefundQueue(ctx, order.OrderId, order.ExpiresAt)
}

func (k Keeper) GetUserTrustScore(address string) int32 {
	// Simplified trust score calculation
	// In production, would consider transaction history, disputes, etc.
	return 75 // Default trust score
}

func (k Keeper) IsKYCVerified(address string) bool {
	// Check KYC status
	// In production, would integrate with KYC system
	return true // Simplified for now
}

func (k Keeper) hasCommonLanguage(langs1, langs2 []string) bool {
	for _, l1 := range langs1 {
		for _, l2 := range langs2 {
			if l1 == l2 {
				return true
			}
		}
	}
	return false
}

func (k Keeper) calculateSevaMitraTrustScore(sevaMitra *types.SevaMitra) float64 {
	if sevaMitra.Stats == nil {
		return 0
	}
	
	successRate := float64(sevaMitra.Stats.SuccessfulTransactions) / float64(sevaMitra.Stats.TotalTransactions)
	return successRate * float64(sevaMitra.Stats.AverageRating) * 20
}