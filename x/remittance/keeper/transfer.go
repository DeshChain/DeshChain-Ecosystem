package keeper

import (
	"context"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/remittance/types"
)

// ========================= Transfer Lifecycle Management =========================

// InitiateTransfer starts a new remittance transfer with optional Sewa Mitra routing
func (k Keeper) InitiateTransfer(
	ctx context.Context,
	sender string,
	recipientAddress string,
	senderCountry string,
	recipientCountry string,
	amount sdk.Coin,
	sourceCurrency string,
	destinationCurrency string,
	settlementMethod string,
	purposeCode string,
	memo string,
	expiresAt time.Time,
	recipientInfo types.RecipientInfo,
	settlementDetails types.SettlementDetails,
	preferSewaMitra bool,
	sewaMitraLocation string,
) (string, error) {
	// Validate basic inputs
	if _, err := sdk.AccAddressFromBech32(sender); err != nil {
		return "", types.ErrInvalidAddress
	}

	params, err := k.GetParams(ctx)
	if err != nil {
		return "", err
	}

	// Check if transfers are enabled
	if !params.TransferEnabled {
		return "", types.ErrServiceUnavailable
	}

	// Check emergency mode
	if params.EmergencyMode {
		return "", types.ErrMaintenanceMode
	}

	// Validate amount limits
	if amount.Amount.LT(params.MinTransferAmount) {
		return "", types.ErrAmountBelowMinimum
	}
	if amount.Amount.GT(params.MaxTransferAmount) {
		return "", types.ErrAmountAboveMaximum
	}

	// Check supported currencies
	if !k.isCurrencySupported(sourceCurrency, params.SupportedCurrencies) ||
		!k.isCurrencySupported(destinationCurrency, params.SupportedCurrencies) {
		return "", types.ErrCurrencyNotSupported
	}

	// Check supported countries
	if !k.isCountrySupported(senderCountry, params.SupportedCountries) ||
		!k.isCountrySupported(recipientCountry, params.SupportedCountries) {
		return "", types.ErrCountryNotSupported
	}

	// Check daily limits
	if err := k.checkDailyLimits(ctx, sender, amount, params); err != nil {
		return "", err
	}

	// Perform KYC checks
	if err := k.performKYCChecks(ctx, sender, amount, params); err != nil {
		return "", err
	}

	// Find or validate corridor
	corridorID, err := k.findCorridorForTransfer(ctx, senderCountry, recipientCountry, sourceCurrency, destinationCurrency)
	if err != nil {
		return "", err
	}

	// Check liquidity availability
	if !k.CheckLiquidityAvailability(ctx, amount.Amount, sourceCurrency, destinationCurrency) {
		return "", types.ErrInsufficientLiquidity
	}

	// Check for Sewa Mitra routing if requested
	var sewaMitraAgent types.SewaMitraAgent
	var sewaMitraCommission sdk.Coin
	usesSewaMitra := false
	
	params, _ := k.GetParams(ctx)
	if preferSewaMitra && params.EnableSewaMitra && settlementMethod == "cash_pickup" {
		// Try to find a suitable Sewa Mitra agent
		agent, err := k.FindNearestSewaMitraAgent(ctx, recipientCountry, sewaMitraLocation, destinationCurrency)
		if err == nil {
			sewaMitraAgent = agent
			usesSewaMitra = true
			
			// Calculate Sewa Mitra commission
			baseCommission, volumeBonus, err := k.CalculateSewaMitraCommission(ctx, agent.AgentId, amount)
			if err == nil {
				sewaMitraCommission = baseCommission.Add(volumeBonus)
			}
		}
	}

	// Calculate fees and exchange rate (including Sewa Mitra commission if applicable)
	exchangeRate, fees, recipientAmount, err := k.calculateTransferCosts(
		ctx, amount, sourceCurrency, destinationCurrency, settlementMethod, sewaMitraCommission,
	)
	if err != nil {
		return "", err
	}

	// Generate transfer ID
	transferID := k.GetNextTransferID(ctx)

	// Create transfer record
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	transfer := types.RemittanceTransfer{
		Id:                     transferID,
		SenderAddress:          sender,
		RecipientAddress:       recipientAddress,
		SenderCountry:          senderCountry,
		RecipientCountry:       recipientCountry,
		Amount:                 amount,
		SourceCurrency:         sourceCurrency,
		DestinationCurrency:    destinationCurrency,
		ExchangeRate:           exchangeRate,
		Fees:                   fees,
		RecipientAmount:        recipientAmount,
		Status:                 types.TRANSFER_STATUS_PENDING,
		CorridorId:             corridorID,
		SettlementMethod:       settlementMethod,
		PurposeCode:            purposeCode,
		Memo:                   memo,
		CreatedAt:              sdkCtx.BlockTime(),
		UpdatedAt:              sdkCtx.BlockTime(),
		ExpiresAt:              expiresAt,
		RecipientInfo:          &recipientInfo,
		SettlementDetails:      &settlementDetails,
		ComplianceInfo:         &types.ComplianceInfo{},
		StatusHistory:          []types.StatusUpdate{},
		// Sewa Mitra integration
		UsesSewaMitra:          usesSewaMitra,
		SewaMitraAgentId:       sewaMitraAgent.AgentId,
		SewaMitraCommission:    sewaMitraCommission,
		SewaMitraLocation:      sewaMitraLocation,
	}

	// Add initial status update
	statusUpdate := types.StatusUpdate{
		Status:    types.TRANSFER_STATUS_PENDING,
		Timestamp: sdkCtx.BlockTime(),
		Message:   "Transfer initiated",
		UpdatedBy: "system",
	}
	transfer.StatusHistory = append(transfer.StatusHistory, statusUpdate)

	// Lock sender funds in escrow
	if err := k.lockFundsInEscrow(ctx, sender, amount); err != nil {
		return "", err
	}

	// Save transfer
	if err := k.SetRemittanceTransfer(ctx, transfer); err != nil {
		// Unlock funds if save fails
		k.unlockFundsFromEscrow(ctx, sender, amount)
		return "", err
	}

	// Emit event
	k.emitTransferEvent(ctx, types.EventTypeInitiateTransfer, transfer)

	return transferID, nil
}

// ConfirmTransfer confirms receipt of a transfer by the recipient
func (k Keeper) ConfirmTransfer(
	ctx context.Context,
	recipient string,
	transferID string,
	confirmationCode string,
	settlementProof string,
) error {
	// Get transfer
	transfer, err := k.GetRemittanceTransfer(ctx, transferID)
	if err != nil {
		return err
	}

	// Validate recipient
	if transfer.RecipientAddress != recipient {
		return types.ErrNotRecipient
	}

	// Check current status
	if transfer.Status != types.TRANSFER_STATUS_PROCESSING {
		return types.ErrTransferNotProcessing
	}

	// Validate confirmation code
	if err := k.validateConfirmationCode(ctx, transfer, confirmationCode); err != nil {
		return err
	}

	// Update transfer status
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	transfer.Status = types.TRANSFER_STATUS_COMPLETED
	transfer.UpdatedAt = sdkCtx.BlockTime()
	transfer.CompletedAt = &sdkCtx.BlockTime()

	// Add status update
	statusUpdate := types.StatusUpdate{
		Status:    types.TRANSFER_STATUS_COMPLETED,
		Timestamp: sdkCtx.BlockTime(),
		Message:   "Transfer confirmed by recipient",
		UpdatedBy: recipient,
	}
	transfer.StatusHistory = append(transfer.StatusHistory, statusUpdate)

	// Create settlement record
	settlement := types.Settlement{
		Id:               k.GetNextSettlementID(ctx),
		TransferId:       transferID,
		PartnerId:        "", // Will be set by settlement processor
		PartnerReference: settlementProof,
		SettledAmount:    transfer.RecipientAmount,
		SettlementTime:   sdkCtx.BlockTime(),
		Status:           "completed",
		Method:           transfer.SettlementMethod,
	}

	// Save settlement
	if err := k.SetSettlement(ctx, settlement); err != nil {
		return err
	}

	// Save updated transfer
	if err := k.SetRemittanceTransfer(ctx, transfer); err != nil {
		return err
	}

	// Record Sewa Mitra commission if applicable
	if transfer.UsesSewaMitra && !transfer.SewaMitraCommission.IsZero() {
		baseCommission, volumeBonus, err := k.CalculateSewaMitraCommission(ctx, transfer.SewaMitraAgentId, transfer.Amount)
		if err == nil {
			if err := k.RecordSewaMitraCommission(ctx, transferID, transfer.SewaMitraAgentId, baseCommission, volumeBonus); err != nil {
				k.Logger(ctx).Error("Failed to record Sewa Mitra commission", "transferID", transferID, "error", err)
			}
			
			// Update agent statistics
			if err := k.UpdateSewaMitraAgentStats(ctx, transfer.SewaMitraAgentId, transfer.Amount, 0, true); err != nil {
				k.Logger(ctx).Error("Failed to update agent stats", "agentID", transfer.SewaMitraAgentId, "error", err)
			}
			
			// Update global Sewa Mitra counters
			counters := k.GetCounters(ctx)
			counters.TotalSevaMitraTransactions++
			k.SetCounters(ctx, counters)
		}
	}

	// Release escrowed funds
	if err := k.releaseFundsFromEscrow(ctx, transferID); err != nil {
		k.Logger(ctx).Error("Failed to release escrowed funds", "transferID", transferID, "error", err)
	}

	// Update counters
	k.incrementTransferCounters(ctx, transfer)

	// Emit event
	k.emitTransferEvent(ctx, types.EventTypeConfirmTransfer, transfer)

	return nil
}

// CancelTransfer cancels a pending transfer
func (k Keeper) CancelTransfer(
	ctx context.Context,
	sender string,
	transferID string,
	reason string,
) error {
	// Get transfer
	transfer, err := k.GetRemittanceTransfer(ctx, transferID)
	if err != nil {
		return err
	}

	// Validate sender
	if transfer.SenderAddress != sender {
		return types.ErrNotOwner
	}

	// Check if transfer can be cancelled
	if transfer.Status != types.TRANSFER_STATUS_PENDING {
		return types.ErrCannotCancelTransfer
	}

	// Update transfer status
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	transfer.Status = types.TRANSFER_STATUS_CANCELLED
	transfer.UpdatedAt = sdkCtx.BlockTime()
	transfer.CancelledAt = &sdkCtx.BlockTime()

	// Add status update
	statusUpdate := types.StatusUpdate{
		Status:    types.TRANSFER_STATUS_CANCELLED,
		Timestamp: sdkCtx.BlockTime(),
		Message:   fmt.Sprintf("Transfer cancelled: %s", reason),
		UpdatedBy: sender,
	}
	transfer.StatusHistory = append(transfer.StatusHistory, statusUpdate)

	// Save updated transfer
	if err := k.SetRemittanceTransfer(ctx, transfer); err != nil {
		return err
	}

	// Refund locked funds
	if err := k.refundLockedFunds(ctx, sender, transfer.Amount); err != nil {
		k.Logger(ctx).Error("Failed to refund locked funds", "transferID", transferID, "error", err)
	}

	// Emit event
	k.emitTransferEvent(ctx, types.EventTypeCancelTransfer, transfer)

	return nil
}

// ProcessTransfer processes a transfer (moves to processing status)
func (k Keeper) ProcessTransfer(ctx context.Context, transferID string, partnerID string) error {
	// Get transfer
	transfer, err := k.GetRemittanceTransfer(ctx, transferID)
	if err != nil {
		return err
	}

	// Check current status
	if transfer.Status != types.TRANSFER_STATUS_PENDING {
		return types.ErrTransferNotPending
	}

	// Validate partner
	partner, err := k.GetCorridorPartner(ctx, partnerID)
	if err != nil {
		return err
	}

	if !partner.IsActive {
		return types.ErrPartnerInactive
	}

	// Update transfer status
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	transfer.Status = types.TRANSFER_STATUS_PROCESSING
	transfer.UpdatedAt = sdkCtx.BlockTime()
	transfer.PartnerId = partnerID

	// Add status update
	statusUpdate := types.StatusUpdate{
		Status:    types.TRANSFER_STATUS_PROCESSING,
		Timestamp: sdkCtx.BlockTime(),
		Message:   fmt.Sprintf("Transfer processing started by partner %s", partnerID),
		UpdatedBy: "system",
	}
	transfer.StatusHistory = append(transfer.StatusHistory, statusUpdate)

	// Save updated transfer
	if err := k.SetRemittanceTransfer(ctx, transfer); err != nil {
		return err
	}

	// Emit event
	k.emitTransferEvent(ctx, types.EventTypeConfirmTransfer, transfer)

	return nil
}

// RefundTransfer processes a refund for a failed transfer
func (k Keeper) RefundTransfer(
	ctx context.Context,
	authority string,
	transferID string,
	reason string,
	refundAmount sdk.Coin,
) error {
	// Validate authority
	if authority != k.GetAuthority() {
		return types.ErrUnauthorized
	}

	// Get transfer
	transfer, err := k.GetRemittanceTransfer(ctx, transferID)
	if err != nil {
		return err
	}

	// Check if refund is valid
	if transfer.Status == types.TRANSFER_STATUS_COMPLETED {
		return fmt.Errorf("cannot refund completed transfer")
	}

	// Update transfer status
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	transfer.Status = types.TRANSFER_STATUS_REFUNDED
	transfer.UpdatedAt = sdkCtx.BlockTime()
	transfer.RefundedAt = &sdkCtx.BlockTime()

	// Add status update
	statusUpdate := types.StatusUpdate{
		Status:    types.TRANSFER_STATUS_REFUNDED,
		Timestamp: sdkCtx.BlockTime(),
		Message:   fmt.Sprintf("Transfer refunded: %s", reason),
		UpdatedBy: authority,
	}
	transfer.StatusHistory = append(transfer.StatusHistory, statusUpdate)

	// Save updated transfer
	if err := k.SetRemittanceTransfer(ctx, transfer); err != nil {
		return err
	}

	// Process refund
	if err := k.processRefund(ctx, transfer.SenderAddress, refundAmount); err != nil {
		return err
	}

	// Emit event
	k.emitTransferEvent(ctx, types.EventTypeRefundTransfer, transfer)

	return nil
}

// Helper functions

// isCurrencySupported checks if a currency is in the supported list
func (k Keeper) isCurrencySupported(currency string, supportedCurrencies []string) bool {
	for _, supported := range supportedCurrencies {
		if supported == currency {
			return true
		}
	}
	return false
}

// isCountrySupported checks if a country is in the supported list
func (k Keeper) isCountrySupported(country string, supportedCountries []string) bool {
	for _, supported := range supportedCountries {
		if supported == country {
			return true
		}
	}
	return false
}

// checkDailyLimits checks if the sender has exceeded daily limits
func (k Keeper) checkDailyLimits(ctx context.Context, sender string, amount sdk.Coin, params types.RemittanceParams) error {
	// Get today's transfers for this sender
	transfers, err := k.GetTransfersBySender(ctx, sender)
	if err != nil {
		return err
	}

	// Count today's transfers
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	today := sdkCtx.BlockTime().Truncate(24 * time.Hour)
	todayCount := 0
	todayVolume := sdk.ZeroInt()

	for _, transfer := range transfers {
		if transfer.CreatedAt.Truncate(24*time.Hour).Equal(today) && 
		   transfer.Status != types.TRANSFER_STATUS_CANCELLED &&
		   transfer.Status != types.TRANSFER_STATUS_REFUNDED {
			todayCount++
			todayVolume = todayVolume.Add(transfer.Amount.Amount)
		}
	}

	// Check limits
	if uint32(todayCount) >= params.MaxDailyTransfers {
		return types.ErrDailyLimitExceeded
	}

	// TODO: Add volume-based daily limits
	return nil
}

// performKYCChecks performs KYC validation
func (k Keeper) performKYCChecks(ctx context.Context, sender string, amount sdk.Coin, params types.RemittanceParams) error {
	if !params.RequireKyc {
		return nil
	}

	// TODO: Implement actual KYC checking logic
	// This would integrate with the compliance keeper
	
	// For now, just check if sender is valid address
	if _, err := sdk.AccAddressFromBech32(sender); err != nil {
		return types.ErrKYCRequired
	}

	return nil
}

// findCorridorForTransfer finds the appropriate corridor for a transfer
func (k Keeper) findCorridorForTransfer(ctx context.Context, sourceCountry, destCountry, sourceCurrency, destCurrency string) (string, error) {
	// TODO: Implement corridor lookup logic
	// For now, generate a default corridor ID
	corridorID := fmt.Sprintf("%s-%s-%s-%s", sourceCountry, destCountry, sourceCurrency, destCurrency)
	return corridorID, nil
}

// calculateTransferCosts calculates exchange rate, fees, and recipient amount
func (k Keeper) calculateTransferCosts(
	ctx context.Context,
	amount sdk.Coin,
	sourceCurrency, destCurrency string,
	settlementMethod string,
	sewaMitraCommission sdk.Coin,
) (sdk.Dec, []types.Fee, sdk.Coin, error) {
	// Get exchange rate
	exchangeRate, err := k.GetExchangeRate(ctx, sourceCurrency, destCurrency)
	if err != nil {
		// Use 1:1 rate if same currency
		if sourceCurrency == destCurrency {
			exchangeRate = sdk.OneDec()
		} else {
			return sdk.ZeroDec(), nil, sdk.Coin{}, err
		}
	}

	// Calculate base recipient amount
	recipientAmount := sdk.NewDecFromInt(amount.Amount).Mul(exchangeRate).TruncateInt()

	// Calculate fees
	params, err := k.GetParams(ctx)
	if err != nil {
		return sdk.ZeroDec(), nil, sdk.Coin{}, err
	}

	transferFee := sdk.NewDecFromInt(amount.Amount).Mul(params.TransferFeeRate).TruncateInt()
	
	fees := []types.Fee{
		{
			Type:        "transfer_fee",
			Amount:      sdk.NewCoin(amount.Denom, transferFee),
			Description: "Base transfer fee",
		},
	}

	// Add Sewa Mitra commission if applicable
	if !sewaMitraCommission.IsZero() {
		fees = append(fees, types.Fee{
			Type:        "seva_mitra_commission",
			Amount:      sewaMitraCommission,
			Description: "Sewa Mitra agent commission",
		})
	}

	// Subtract fees from recipient amount
	totalFees := transferFee
	if !sewaMitraCommission.IsZero() {
		totalFees = totalFees.Add(sewaMitraCommission.Amount)
	}
	finalRecipientAmount := recipientAmount.Sub(totalFees)

	if finalRecipientAmount.IsNegative() {
		return sdk.ZeroDec(), nil, sdk.Coin{}, types.ErrInvalidAmount
	}

	recipientCoin := sdk.NewCoin(destCurrency, finalRecipientAmount)
	return exchangeRate, fees, recipientCoin, nil
}

// lockFundsInEscrow locks sender funds
func (k Keeper) lockFundsInEscrow(ctx context.Context, sender string, amount sdk.Coin) error {
	senderAddr, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return err
	}

	return k.bankKeeper.SendCoinsFromAccountToModule(
		ctx, senderAddr, types.ModuleName, sdk.NewCoins(amount),
	)
}

// unlockFundsFromEscrow unlocks sender funds (for failed transactions)
func (k Keeper) unlockFundsFromEscrow(ctx context.Context, sender string, amount sdk.Coin) error {
	senderAddr, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return err
	}

	return k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx, types.ModuleName, senderAddr, sdk.NewCoins(amount),
	)
}

// releaseFundsFromEscrow releases funds after successful transfer
func (k Keeper) releaseFundsFromEscrow(ctx context.Context, transferID string) error {
	// For completed transfers, funds are considered delivered
	// In a real implementation, this might involve sending to settlement partners
	return nil
}

// refundLockedFunds refunds cancelled transfer funds
func (k Keeper) refundLockedFunds(ctx context.Context, sender string, amount sdk.Coin) error {
	return k.unlockFundsFromEscrow(ctx, sender, amount)
}

// processRefund processes refund payments
func (k Keeper) processRefund(ctx context.Context, recipient string, amount sdk.Coin) error {
	recipientAddr, err := sdk.AccAddressFromBech32(recipient)
	if err != nil {
		return err
	}

	return k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx, types.ModuleName, recipientAddr, sdk.NewCoins(amount),
	)
}

// validateConfirmationCode validates transfer confirmation codes
func (k Keeper) validateConfirmationCode(ctx context.Context, transfer types.RemittanceTransfer, code string) error {
	// TODO: Implement proper confirmation code validation
	// For now, accept any non-empty code
	if code == "" {
		return fmt.Errorf("confirmation code required")
	}
	return nil
}

// incrementTransferCounters updates global transfer statistics
func (k Keeper) incrementTransferCounters(ctx context.Context, transfer types.RemittanceTransfer) {
	counters := k.GetCounters(ctx)
	counters.TotalTransfers++
	
	// Convert to USD equivalent for global tracking
	// TODO: Implement proper USD conversion
	usdAmount := transfer.Amount.Amount.Uint64()
	counters.TotalVolumeUsd += usdAmount
	
	k.SetCounters(ctx, counters)
}

// emitTransferEvent emits blockchain events for transfers
func (k Keeper) emitTransferEvent(ctx context.Context, eventType string, transfer types.RemittanceTransfer) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			eventType,
			sdk.NewAttribute(types.AttributeKeyTransferID, transfer.Id),
			sdk.NewAttribute(types.AttributeKeySender, transfer.SenderAddress),
			sdk.NewAttribute(types.AttributeKeyRecipient, transfer.RecipientAddress),
			sdk.NewAttribute(types.AttributeKeyAmount, transfer.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyStatus, transfer.Status.String()),
		),
	)
}