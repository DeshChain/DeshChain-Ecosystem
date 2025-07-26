package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/tradefinance/types"
)

// ProcessLcExpiries checks and updates expired LCs
func (k Keeper) ProcessLcExpiries(ctx sdk.Context) {
	lcs := k.GetAllLettersOfCredit(ctx)
	currentTime := ctx.BlockTime()

	for _, lc := range lcs {
		// Check if LC is active and has expired
		if (lc.Status == "issued" || lc.Status == "accepted" || lc.Status == "documents_presented") && 
			currentTime.After(lc.ExpiryDate) {
			
			// Mark LC as expired
			lc.Status = "expired"
			lc.UpdatedAt = currentTime
			k.SetLetterOfCredit(ctx, lc)

			// Release collateral back to issuing bank
			issuingBank, found := k.GetTradeParty(ctx, lc.IssuingBankId)
			if found {
				issuingBankAddr, err := sdk.AccAddressFromBech32(issuingBank.DeshAddress)
				if err == nil {
					k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, issuingBankAddr, sdk.NewCoins(lc.Collateral))
				}
			}

			// Update stats
			stats := k.GetTradeFinanceStats(ctx)
			stats.ActiveLcs--
			k.SetTradeFinanceStats(ctx, stats)

			// Emit event
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					"lc_expired",
					sdk.NewAttribute(types.AttributeKeyLcId, lc.LcId),
					sdk.NewAttribute(types.AttributeKeyStatus, "expired"),
				),
			)
		}
	}
}

// ProcessPaymentDueDates checks for overdue payments
func (k Keeper) ProcessPaymentDueDates(ctx sdk.Context) {
	instructions := k.GetAllPaymentInstructions(ctx)
	currentTime := ctx.BlockTime()
	params := k.GetParams(ctx)

	for _, instruction := range instructions {
		// Check if payment is pending and overdue
		if instruction.Status == "pending" && currentTime.After(instruction.DueDate) {
			// Apply late payment penalty if configured
			if params.Fees.LatePaymentPenalty > 0 {
				// Calculate penalty
				penaltyAmount := instruction.Amount.Amount.Mul(sdk.NewInt(int64(params.Fees.LatePaymentPenalty))).Quo(sdk.NewInt(10000))
				
				// Update instruction with penalty (in real implementation, would create a separate penalty instruction)
				instruction.Status = "overdue"
				k.SetPaymentInstruction(ctx, instruction)

				// Emit event
				ctx.EventManager().EmitEvent(
					sdk.NewEvent(
						"payment_overdue",
						sdk.NewAttribute(types.AttributeKeyPaymentInstructionId, instruction.InstructionId),
						sdk.NewAttribute(types.AttributeKeyLcId, instruction.LcId),
						sdk.NewAttribute("penalty_amount", penaltyAmount.String()),
					),
				)
			}
		}
	}
}

// UpdateStatistics updates module-wide statistics
func (k Keeper) UpdateStatistics(ctx sdk.Context) {
	stats := k.GetTradeFinanceStats(ctx)
	
	// Count active LCs
	activeLcs := uint64(0)
	completedTrades := uint64(0)
	lcs := k.GetAllLettersOfCredit(ctx)
	
	for _, lc := range lcs {
		switch lc.Status {
		case "issued", "accepted", "documents_presented", "documents_verified":
			activeLcs++
		case "paid":
			completedTrades++
		}
	}
	
	stats.ActiveLcs = activeLcs
	stats.CompletedTrades = completedTrades
	stats.LastUpdate = ctx.BlockTime()
	
	// Calculate average processing time (simplified)
	if completedTrades > 0 {
		// In real implementation, would track actual processing times
		stats.AverageProcessingHours = 4 // Default 4 hours
	}
	
	k.SetTradeFinanceStats(ctx, stats)
}

// GetAllTradeDocuments returns all trade documents (for genesis export)
func (k Keeper) GetAllTradeDocuments(ctx sdk.Context) []types.TradeDocument {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.TradeDocumentPrefix)
	
	var documents []types.TradeDocument
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var doc types.TradeDocument
		k.cdc.MustUnmarshal(iterator.Value(), &doc)
		documents = append(documents, doc)
	}
	
	return documents
}

// GetAllPaymentInstructions returns all payment instructions (for genesis export)
func (k Keeper) GetAllPaymentInstructions(ctx sdk.Context) []types.PaymentInstruction {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PaymentInstructionPrefix)
	
	var instructions []types.PaymentInstruction
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var instruction types.PaymentInstruction
		k.cdc.MustUnmarshal(iterator.Value(), &instruction)
		instructions = append(instructions, instruction)
	}
	
	return instructions
}