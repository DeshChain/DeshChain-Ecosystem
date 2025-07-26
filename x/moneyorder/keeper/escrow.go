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
	"github.com/DeshChain/DeshChain-Ecosystem/x/moneyorder/types"
)

// CreateEscrow creates a new escrow for secure fund holding
func (k Keeper) CreateEscrow(ctx sdk.Context, depositor sdk.AccAddress, amount sdk.Coin, orderID string, tradeID string, expiresAt time.Time) (*types.Escrow, error) {
	// Generate escrow ID
	escrowID := k.generateEscrowID(ctx)
	
	// Create escrow address
	escrowAddr := k.generateEscrowAddress(escrowID)
	
	// Calculate platform fee
	platformFee := k.calculatePlatformFee(amount)
	
	// Create escrow record
	escrow := &types.Escrow{
		EscrowId:    escrowID,
		OrderId:     orderID,
		TradeId:     tradeID,
		Depositor:   depositor.String(),
		Amount:      amount,
		PlatformFee: platformFee,
		Status:      types.EscrowStatus_ESCROW_STATUS_ACTIVE,
		CreatedAt:   ctx.BlockTime(),
		ExpiresAt:   expiresAt,
		RefundTo:    depositor.String(),
	}
	
	// Transfer funds to escrow address
	err := k.bankKeeper.SendCoins(ctx, depositor, escrowAddr, sdk.NewCoins(amount.Add(platformFee)))
	if err != nil {
		return nil, err
	}
	
	// Store escrow
	k.SetEscrow(ctx, escrow)
	
	// Schedule auto-refund check
	k.ScheduleEscrowCheck(ctx, escrowID, expiresAt)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeEscrowCreated,
			sdk.NewAttribute("escrow_id", escrowID),
			sdk.NewAttribute("depositor", depositor.String()),
			sdk.NewAttribute("amount", amount.String()),
			sdk.NewAttribute("expires_at", expiresAt.Format(time.RFC3339)),
		),
	)
	
	return escrow, nil
}

// ReleaseEscrow releases escrow funds to the recipient
func (k Keeper) ReleaseEscrow(ctx sdk.Context, escrowID string, releaseTo sdk.AccAddress, authorizedBy sdk.AccAddress) error {
	escrow, found := k.GetEscrow(ctx, escrowID)
	if !found {
		return types.ErrEscrowNotFound
	}
	
	// Check status
	if escrow.Status != types.EscrowStatus_ESCROW_STATUS_ACTIVE {
		return types.ErrInvalidEscrowStatus
	}
	
	// Check authorization
	if !k.isAuthorizedToRelease(ctx, escrow, authorizedBy) {
		return types.ErrUnauthorized
	}
	
	// Get escrow address
	escrowAddr := k.generateEscrowAddress(escrowID)
	
	// Release funds (minus platform fee)
	err := k.bankKeeper.SendCoins(ctx, escrowAddr, releaseTo, sdk.NewCoins(escrow.Amount))
	if err != nil {
		return err
	}
	
	// Transfer platform fee to fee collector
	if escrow.PlatformFee.IsPositive() {
		feeCollector := k.authKeeper.GetModuleAddress(types.FeeCollectorName)
		err = k.bankKeeper.SendCoins(ctx, escrowAddr, feeCollector, sdk.NewCoins(escrow.PlatformFee))
		if err != nil {
			return err
		}
	}
	
	// Update status
	escrow.Status = types.EscrowStatus_ESCROW_STATUS_RELEASED
	escrow.ReleaseTo = releaseTo.String()
	k.SetEscrow(ctx, escrow)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeEscrowReleased,
			sdk.NewAttribute("escrow_id", escrowID),
			sdk.NewAttribute("released_to", releaseTo.String()),
			sdk.NewAttribute("amount", escrow.Amount.String()),
		),
	)
	
	return nil
}

// RefundEscrow refunds escrow including platform fees
func (k Keeper) RefundEscrow(ctx sdk.Context, escrowID string, reason string) error {
	escrow, found := k.GetEscrow(ctx, escrowID)
	if !found {
		return types.ErrEscrowNotFound
	}
	
	// Check status
	if escrow.Status != types.EscrowStatus_ESCROW_STATUS_ACTIVE {
		return types.ErrInvalidEscrowStatus
	}
	
	// Get addresses
	escrowAddr := k.generateEscrowAddress(escrowID)
	refundAddr, err := sdk.AccAddressFromBech32(escrow.RefundTo)
	if err != nil {
		return err
	}
	
	// Refund full amount including platform fee
	totalRefund := escrow.Amount.Add(escrow.PlatformFee)
	err = k.bankKeeper.SendCoins(ctx, escrowAddr, refundAddr, sdk.NewCoins(totalRefund))
	if err != nil {
		return err
	}
	
	// Update status
	escrow.Status = types.EscrowStatus_ESCROW_STATUS_REFUNDED
	k.SetEscrow(ctx, escrow)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeEscrowRefunded,
			sdk.NewAttribute("escrow_id", escrowID),
			sdk.NewAttribute("refund_to", escrow.RefundTo),
			sdk.NewAttribute("amount", totalRefund.String()),
			sdk.NewAttribute("reason", reason),
		),
	)
	
	return nil
}

// DisputeEscrow marks escrow as disputed
func (k Keeper) DisputeEscrow(ctx sdk.Context, escrowID string, disputedBy sdk.AccAddress, reason string) error {
	escrow, found := k.GetEscrow(ctx, escrowID)
	if !found {
		return types.ErrEscrowNotFound
	}
	
	// Check status
	if escrow.Status != types.EscrowStatus_ESCROW_STATUS_ACTIVE {
		return types.ErrInvalidEscrowStatus
	}
	
	// Check if disputer is involved in the trade
	if !k.isPartyToEscrow(ctx, escrow, disputedBy) {
		return types.ErrUnauthorized
	}
	
	// Update status
	escrow.Status = types.EscrowStatus_ESCROW_STATUS_DISPUTED
	k.SetEscrow(ctx, escrow)
	
	// Create dispute record
	dispute := &types.Dispute{
		DisputeId:   k.generateDisputeID(ctx),
		EscrowId:    escrowID,
		TradeId:     escrow.TradeId,
		DisputedBy:  disputedBy.String(),
		Reason:      reason,
		Status:      types.DisputeStatus_DISPUTE_STATUS_OPEN,
		CreatedAt:   ctx.BlockTime(),
	}
	k.SetDispute(ctx, dispute)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeEscrowDisputed,
			sdk.NewAttribute("escrow_id", escrowID),
			sdk.NewAttribute("disputed_by", disputedBy.String()),
			sdk.NewAttribute("reason", reason),
		),
	)
	
	return nil
}

// ResolveDispute resolves a disputed escrow
func (k Keeper) ResolveDispute(ctx sdk.Context, disputeID string, resolution types.DisputeResolution, resolvedBy sdk.AccAddress) error {
	dispute, found := k.GetDispute(ctx, disputeID)
	if !found {
		return types.ErrDisputeNotFound
	}
	
	if dispute.Status != types.DisputeStatus_DISPUTE_STATUS_OPEN {
		return types.ErrInvalidDisputeStatus
	}
	
	// Check authorization (must be validator or governance)
	if !k.isAuthorizedResolver(ctx, resolvedBy) {
		return types.ErrUnauthorized
	}
	
	escrow, found := k.GetEscrow(ctx, dispute.EscrowId)
	if !found {
		return types.ErrEscrowNotFound
	}
	
	// Apply resolution
	switch resolution.Decision {
	case types.DisputeDecision_RELEASE_TO_BUYER:
		trade, _ := k.GetP2PTrade(ctx, escrow.TradeId)
		buyerAddr, _ := sdk.AccAddressFromBech32(trade.Buyer)
		err := k.ReleaseEscrow(ctx, dispute.EscrowId, buyerAddr, resolvedBy)
		if err != nil {
			return err
		}
		
	case types.DisputeDecision_RELEASE_TO_SELLER:
		trade, _ := k.GetP2PTrade(ctx, escrow.TradeId)
		sellerAddr, _ := sdk.AccAddressFromBech32(trade.Seller)
		err := k.ReleaseEscrow(ctx, dispute.EscrowId, sellerAddr, resolvedBy)
		if err != nil {
			return err
		}
		
	case types.DisputeDecision_SPLIT:
		// Split funds between parties
		trade, _ := k.GetP2PTrade(ctx, escrow.TradeId)
		buyerAddr, _ := sdk.AccAddressFromBech32(trade.Buyer)
		sellerAddr, _ := sdk.AccAddressFromBech32(trade.Seller)
		
		// Calculate split amounts
		buyerAmount := escrow.Amount.Amount.MulRaw(resolution.BuyerPercentage).QuoRaw(100)
		sellerAmount := escrow.Amount.Amount.Sub(buyerAmount)
		
		// Release split amounts
		escrowAddr := k.generateEscrowAddress(dispute.EscrowId)
		k.bankKeeper.SendCoins(ctx, escrowAddr, buyerAddr, sdk.NewCoins(sdk.NewCoin(escrow.Amount.Denom, buyerAmount)))
		k.bankKeeper.SendCoins(ctx, escrowAddr, sellerAddr, sdk.NewCoins(sdk.NewCoin(escrow.Amount.Denom, sellerAmount)))
		
	case types.DisputeDecision_REFUND:
		err := k.RefundEscrow(ctx, dispute.EscrowId, "Dispute resolved with refund")
		if err != nil {
			return err
		}
	}
	
	// Update dispute status
	dispute.Status = types.DisputeStatus_DISPUTE_STATUS_RESOLVED
	dispute.Resolution = &resolution
	dispute.ResolvedAt = ctx.BlockTime()
	dispute.ResolvedBy = resolvedBy.String()
	k.SetDispute(ctx, dispute)
	
	// Update user trust scores based on resolution
	k.updateDisputeScores(ctx, dispute, resolution)
	
	return nil
}

// CheckExpiredEscrows checks and refunds expired escrows
func (k Keeper) CheckExpiredEscrows(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.EscrowPrefix)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var escrow types.Escrow
		k.cdc.MustUnmarshal(iterator.Value(), &escrow)
		
		// Skip non-active escrows
		if escrow.Status != types.EscrowStatus_ESCROW_STATUS_ACTIVE {
			continue
		}
		
		// Check if expired
		if ctx.BlockTime().After(escrow.ExpiresAt) {
			// Auto-refund expired escrow
			err := k.RefundEscrow(ctx, escrow.EscrowId, "Escrow expired without match")
			if err != nil {
				ctx.Logger().Error("Failed to refund expired escrow", "escrow_id", escrow.EscrowId, "error", err)
			}
		}
	}
}

// Helper functions

func (k Keeper) generateEscrowID(ctx sdk.Context) string {
	return fmt.Sprintf("ESC-%d-%s", ctx.BlockHeight(), k.generateRandomString(8))
}

func (k Keeper) generateEscrowAddress(escrowID string) sdk.AccAddress {
	// Deterministic escrow address generation
	return sdk.AccAddress([]byte("escrow-" + escrowID))
}

func (k Keeper) generateDisputeID(ctx sdk.Context) string {
	return fmt.Sprintf("DISP-%d-%s", ctx.BlockHeight(), k.generateRandomString(6))
}

func (k Keeper) isAuthorizedToRelease(ctx sdk.Context, escrow *types.Escrow, authorizer sdk.AccAddress) bool {
	// For P2P trades, seller confirms payment received
	if escrow.TradeId != "" {
		trade, found := k.GetP2PTrade(ctx, escrow.TradeId)
		if found && trade.Seller == authorizer.String() {
			return true
		}
	}
	
	// System admins can release
	return k.isSystemAdmin(ctx, authorizer)
}

func (k Keeper) isPartyToEscrow(ctx sdk.Context, escrow *types.Escrow, party sdk.AccAddress) bool {
	// Check if party is depositor
	if escrow.Depositor == party.String() {
		return true
	}
	
	// Check if party is in the trade
	if escrow.TradeId != "" {
		trade, found := k.GetP2PTrade(ctx, escrow.TradeId)
		if found && (trade.Buyer == party.String() || trade.Seller == party.String()) {
			return true
		}
	}
	
	return false
}

func (k Keeper) isAuthorizedResolver(ctx sdk.Context, resolver sdk.AccAddress) bool {
	// Check if resolver is a validator
	// In production, would check validator set or governance module
	return k.isSystemAdmin(ctx, resolver)
}

func (k Keeper) isSystemAdmin(ctx sdk.Context, addr sdk.AccAddress) bool {
	// Simplified admin check
	// In production, would check against governance or admin module
	return true
}

func (k Keeper) updateDisputeScores(ctx sdk.Context, dispute *types.Dispute, resolution types.DisputeResolution) {
	// Update trust scores based on dispute outcome
	// This affects future P2P matching
	
	trade, found := k.GetP2PTrade(ctx, dispute.TradeId)
	if !found {
		return
	}
	
	// Apply score adjustments based on resolution
	switch resolution.Decision {
	case types.DisputeDecision_RELEASE_TO_BUYER:
		// Seller at fault
		k.decreaseTrustScore(ctx, trade.Seller, 10)
		k.increaseTrustScore(ctx, trade.Buyer, 5)
		
	case types.DisputeDecision_RELEASE_TO_SELLER:
		// Buyer at fault
		k.decreaseTrustScore(ctx, trade.Buyer, 10)
		k.increaseTrustScore(ctx, trade.Seller, 5)
		
	case types.DisputeDecision_SPLIT:
		// Both parties partially at fault
		k.decreaseTrustScore(ctx, trade.Buyer, 5)
		k.decreaseTrustScore(ctx, trade.Seller, 5)
	}
}

// ScheduleEscrowCheck schedules an escrow for expiry check
func (k Keeper) ScheduleEscrowCheck(ctx sdk.Context, escrowID string, expiresAt time.Time) {
	// In production, would use a proper scheduler or cron job
	// For now, we'll check during BeginBlock
	k.AddToExpiryQueue(ctx, escrowID, expiresAt)
}