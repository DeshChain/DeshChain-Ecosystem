package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/deshchain/x/remittance/types"
)

// ========================= Settlement Management =========================

// SetSettlement stores a settlement record
func (k Keeper) SetSettlement(ctx context.Context, settlement types.Settlement) error {
	store := k.GetStore(ctx)
	key := types.SettlementKey(settlement.TransferId)
	bz := k.cdc.MustMarshal(&settlement)
	store.Set(key, bz)
	return nil
}

// GetSettlement retrieves a settlement record
func (k Keeper) GetSettlement(ctx context.Context, transferID string) (types.Settlement, error) {
	store := k.GetStore(ctx)
	key := types.SettlementKey(transferID)
	bz := store.Get(key)
	if bz == nil {
		return types.Settlement{}, types.ErrSettlementNotFound
	}

	var settlement types.Settlement
	k.cdc.MustUnmarshal(bz, &settlement)
	return settlement, nil
}

// ProcessSettlement processes settlement for a transfer
func (k Keeper) ProcessSettlement(
	ctx context.Context,
	authority string,
	transferID string,
	partnerID string,
	partnerReference string,
	settlementProof string,
	settledAmount sdk.Coin,
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

	// Validate transfer status
	if transfer.Status != types.TRANSFER_STATUS_PROCESSING {
		return types.ErrTransferNotProcessing
	}

	// Get partner
	partner, err := k.GetCorridorPartner(ctx, partnerID)
	if err != nil {
		return err
	}

	if !partner.IsActive {
		return types.ErrPartnerInactive
	}

	// Validate settlement method
	if !k.isSettlementMethodSupported(transfer.SettlementMethod, partner.SupportedMethods) {
		return types.ErrInvalidSettlementMethod
	}

	// Create settlement record
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	settlement := types.Settlement{
		Id:               k.GetNextSettlementID(ctx),
		TransferId:       transferID,
		PartnerId:        partnerID,
		PartnerReference: partnerReference,
		SettledAmount:    settledAmount,
		SettlementTime:   sdkCtx.BlockTime(),
		Status:           "processing",
		Method:           transfer.SettlementMethod,
		ProofHash:        settlementProof,
	}

	// Process based on settlement method
	switch transfer.SettlementMethod {
	case "bank_transfer":
		err = k.processBankTransfer(ctx, transfer, partner, &settlement)
	case "mobile_wallet":
		err = k.processMobileWallet(ctx, transfer, partner, &settlement)
	case "cash_pickup":
		err = k.processCashPickup(ctx, transfer, partner, &settlement)
	default:
		return types.ErrInvalidSettlementMethod
	}

	if err != nil {
		settlement.Status = "failed"
		settlement.ErrorMessage = err.Error()
		k.SetSettlement(ctx, settlement)
		return err
	}

	// Mark settlement as completed
	settlement.Status = "completed"

	// Update transfer status
	transfer.Status = types.TRANSFER_STATUS_PROCESSING // Will be completed when recipient confirms
	transfer.UpdatedAt = sdkCtx.BlockTime()
	transfer.PartnerId = partnerID

	// Add status update
	statusUpdate := types.StatusUpdate{
		Status:    types.TRANSFER_STATUS_PROCESSING,
		Timestamp: sdkCtx.BlockTime(),
		Message:   fmt.Sprintf("Settlement processed by partner %s", partnerID),
		UpdatedBy: authority,
	}
	transfer.StatusHistory = append(transfer.StatusHistory, statusUpdate)

	// Save records
	if err := k.SetSettlement(ctx, settlement); err != nil {
		return err
	}

	if err := k.SetRemittanceTransfer(ctx, transfer); err != nil {
		return err
	}

	// Emit event
	k.emitSettlementEvent(ctx, types.EventTypeProcessSettlement, settlement)

	return nil
}

// processBankTransfer handles bank transfer settlements
func (k Keeper) processBankTransfer(
	ctx context.Context,
	transfer types.RemittanceTransfer,
	partner types.CorridorPartner,
	settlement *types.Settlement,
) error {
	// Validate bank details
	if transfer.SettlementDetails == nil {
		return types.ErrInvalidSettlementDetails
	}

	details := transfer.SettlementDetails
	if details.BankName == "" || details.AccountNumber == "" {
		return types.ErrInvalidBankDetails
	}

	// Log settlement processing
	k.Logger(ctx).Info("Processing bank transfer",
		"transferID", transfer.Id,
		"partnerID", partner.Id,
		"bankName", details.BankName,
		"amount", settlement.SettledAmount,
	)

	// In a real implementation, this would:
	// 1. Validate bank details with partner systems
	// 2. Submit transfer request to banking APIs
	// 3. Monitor transfer status
	// 4. Handle success/failure responses
	
	// For now, we simulate successful processing
	settlement.Notes = fmt.Sprintf("Bank transfer to %s, Account: %s", 
		details.BankName, maskAccountNumber(details.AccountNumber))

	return nil
}

// processMobileWallet handles mobile wallet settlements
func (k Keeper) processMobileWallet(
	ctx context.Context,
	transfer types.RemittanceTransfer,
	partner types.CorridorPartner,
	settlement *types.Settlement,
) error {
	// Validate mobile wallet details
	if transfer.SettlementDetails == nil {
		return types.ErrInvalidSettlementDetails
	}

	details := transfer.SettlementDetails
	if details.MobileNumber == "" {
		return types.ErrInvalidPhoneNumber
	}

	// Log settlement processing
	k.Logger(ctx).Info("Processing mobile wallet transfer",
		"transferID", transfer.Id,
		"partnerID", partner.Id,
		"mobileProvider", details.MobileProvider,
		"amount", settlement.SettledAmount,
	)

	// In a real implementation, this would:
	// 1. Validate mobile number format
	// 2. Check wallet provider compatibility
	// 3. Submit transfer to mobile money APIs
	// 4. Handle wallet-specific requirements
	
	// For now, we simulate successful processing
	settlement.Notes = fmt.Sprintf("Mobile wallet transfer to %s via %s", 
		maskPhoneNumber(details.MobileNumber), details.MobileProvider)

	return nil
}

// processCashPickup handles cash pickup settlements
func (k Keeper) processCashPickup(
	ctx context.Context,
	transfer types.RemittanceTransfer,
	partner types.CorridorPartner,
	settlement *types.Settlement,
) error {
	// Validate pickup details
	if transfer.SettlementDetails == nil {
		return types.ErrInvalidSettlementDetails
	}

	details := transfer.SettlementDetails
	if details.PickupLocationId == "" {
		return types.ErrInvalidPickupLocation
	}

	// Generate pickup code
	pickupCode := k.generatePickupCode(ctx, transfer.Id)

	// Log settlement processing
	k.Logger(ctx).Info("Processing cash pickup",
		"transferID", transfer.Id,
		"partnerID", partner.Id,
		"locationID", details.PickupLocationId,
		"pickupCode", pickupCode,
		"amount", settlement.SettledAmount,
	)

	// In a real implementation, this would:
	// 1. Validate pickup location
	// 2. Reserve cash at pickup location
	// 3. Generate secure pickup codes
	// 4. Send notification to recipient
	
	// For now, we simulate successful processing
	settlement.Notes = fmt.Sprintf("Cash available for pickup at %s, Code: %s", 
		details.PickupLocationName, pickupCode)
	settlement.PickupCode = pickupCode

	return nil
}

// ========================= Partner Management =========================

// SetCorridorPartner stores a corridor partner
func (k Keeper) SetCorridorPartner(ctx context.Context, partner types.CorridorPartner) error {
	store := k.GetStore(ctx)
	key := types.PartnerKey(partner.Id)
	bz := k.cdc.MustMarshal(&partner)
	store.Set(key, bz)
	return nil
}

// GetCorridorPartner retrieves a corridor partner
func (k Keeper) GetCorridorPartner(ctx context.Context, partnerID string) (types.CorridorPartner, error) {
	store := k.GetStore(ctx)
	key := types.PartnerKey(partnerID)
	bz := store.Get(key)
	if bz == nil {
		return types.CorridorPartner{}, types.ErrPartnerNotFound
	}

	var partner types.CorridorPartner
	k.cdc.MustUnmarshal(bz, &partner)
	return partner, nil
}

// GetAllCorridorPartners returns all corridor partners
func (k Keeper) GetAllCorridorPartners(ctx context.Context) ([]types.CorridorPartner, error) {
	store := k.GetStore(ctx)
	iterator := store.Iterator(types.PartnerKeyPrefix, nil)
	defer iterator.Close()

	var partners []types.CorridorPartner
	for ; iterator.Valid(); iterator.Next() {
		var partner types.CorridorPartner
		k.cdc.MustUnmarshal(iterator.Value(), &partner)
		partners = append(partners, partner)
	}

	return partners, nil
}

// RegisterPartner registers a new settlement partner
func (k Keeper) RegisterPartner(
	ctx context.Context,
	authority string,
	partnerName string,
	partnerType string,
	country string,
	supportedMethods []string,
	settlementCurrency string,
	feeStructure sdk.Dec,
	contactInfo types.ContactInfo,
) (string, error) {
	// Validate authority
	if authority != k.GetAuthority() {
		return "", types.ErrUnauthorized
	}

	// Validate inputs
	if partnerName == "" {
		return "", fmt.Errorf("partner name cannot be empty")
	}

	if country == "" {
		return "", types.ErrInvalidCountryCode
	}

	if len(supportedMethods) == 0 {
		return "", fmt.Errorf("at least one settlement method must be supported")
	}

	// Generate partner ID
	partnerID := k.GetNextPartnerID(ctx)

	// Create partner record
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	partner := types.CorridorPartner{
		Id:                 partnerID,
		Name:               partnerName,
		Type:               partnerType,
		Country:            country,
		SupportedMethods:   supportedMethods,
		SettlementCurrency: settlementCurrency,
		FeeStructure:       feeStructure,
		IsActive:           true,
		ContactInfo:        &contactInfo,
		CreatedAt:          sdkCtx.BlockTime(),
		UpdatedAt:          sdkCtx.BlockTime(),
	}

	// Save partner
	if err := k.SetCorridorPartner(ctx, partner); err != nil {
		return "", err
	}

	// Emit event
	k.emitPartnerEvent(ctx, types.EventTypeRegisterPartner, partner)

	return partnerID, nil
}

// UpdatePartner updates partner information
func (k Keeper) UpdatePartner(
	ctx context.Context,
	authority string,
	partnerID string,
	isActive bool,
	supportedMethods []string,
	feeStructure sdk.Dec,
	contactInfo types.ContactInfo,
) error {
	// Validate authority
	if authority != k.GetAuthority() {
		return types.ErrUnauthorized
	}

	// Get existing partner
	partner, err := k.GetCorridorPartner(ctx, partnerID)
	if err != nil {
		return err
	}

	// Update fields
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	partner.IsActive = isActive
	partner.SupportedMethods = supportedMethods
	partner.FeeStructure = feeStructure
	partner.ContactInfo = &contactInfo
	partner.UpdatedAt = sdkCtx.BlockTime()

	// Save updated partner
	if err := k.SetCorridorPartner(ctx, partner); err != nil {
		return err
	}

	return nil
}

// GetPartnersByCountry returns partners for a specific country
func (k Keeper) GetPartnersByCountry(ctx context.Context, country string) ([]types.CorridorPartner, error) {
	partners, err := k.GetAllCorridorPartners(ctx)
	if err != nil {
		return nil, err
	}

	var countryPartners []types.CorridorPartner
	for _, partner := range partners {
		if partner.Country == country && partner.IsActive {
			countryPartners = append(countryPartners, partner)
		}
	}

	return countryPartners, nil
}

// Helper functions

// isSettlementMethodSupported checks if a settlement method is supported by a partner
func (k Keeper) isSettlementMethodSupported(method string, supportedMethods []string) bool {
	for _, supported := range supportedMethods {
		if supported == method {
			return true
		}
	}
	return false
}

// maskAccountNumber masks sensitive account number information
func maskAccountNumber(accountNumber string) string {
	if len(accountNumber) <= 4 {
		return "****"
	}
	return "****" + accountNumber[len(accountNumber)-4:]
}

// maskPhoneNumber masks sensitive phone number information
func maskPhoneNumber(phoneNumber string) string {
	if len(phoneNumber) <= 4 {
		return "****"
	}
	return "****" + phoneNumber[len(phoneNumber)-4:]
}

// generatePickupCode generates a secure pickup code for cash pickup
func (k Keeper) generatePickupCode(ctx context.Context, transferID string) string {
	// In a real implementation, this would generate a cryptographically secure code
	// For now, use a simple format based on transfer ID
	return fmt.Sprintf("PU-%s", transferID[len(transferID)-6:])
}

// GetSettlementsByStatus returns settlements filtered by status
func (k Keeper) GetSettlementsByStatus(ctx context.Context, status string) ([]types.Settlement, error) {
	settlements, err := k.GetAllSettlements(ctx)
	if err != nil {
		return nil, err
	}

	var filtered []types.Settlement
	for _, settlement := range settlements {
		if settlement.Status == status {
			filtered = append(filtered, settlement)
		}
	}

	return filtered, nil
}

// GetAllSettlements returns all settlement records
func (k Keeper) GetAllSettlements(ctx context.Context) ([]types.Settlement, error) {
	store := k.GetStore(ctx)
	iterator := store.Iterator(types.SettlementKeyPrefix, nil)
	defer iterator.Close()

	var settlements []types.Settlement
	for ; iterator.Valid(); iterator.Next() {
		var settlement types.Settlement
		k.cdc.MustUnmarshal(iterator.Value(), &settlement)
		settlements = append(settlements, settlement)
	}

	return settlements, nil
}

// GetPartnerStatistics returns statistics for a partner
func (k Keeper) GetPartnerStatistics(ctx context.Context, partnerID string) (types.PartnerStatistics, error) {
	// Get partner
	partner, err := k.GetCorridorPartner(ctx, partnerID)
	if err != nil {
		return types.PartnerStatistics{}, err
	}

	// Get all settlements for this partner
	settlements, err := k.GetAllSettlements(ctx)
	if err != nil {
		return types.PartnerStatistics{}, err
	}

	// Calculate statistics
	var totalTransfers uint64
	var totalVolume sdk.Int = sdk.ZeroInt()
	var successfulTransfers uint64
	var failedTransfers uint64

	for _, settlement := range settlements {
		if settlement.PartnerId == partnerID {
			totalTransfers++
			totalVolume = totalVolume.Add(settlement.SettledAmount.Amount)
			
			switch settlement.Status {
			case "completed":
				successfulTransfers++
			case "failed":
				failedTransfers++
			}
		}
	}

	var successRate sdk.Dec = sdk.ZeroDec()
	if totalTransfers > 0 {
		successRate = sdk.NewDec(int64(successfulTransfers)).Quo(sdk.NewDec(int64(totalTransfers)))
	}

	stats := types.PartnerStatistics{
		PartnerId:           partner.Id,
		PartnerName:         partner.Name,
		PartnerType:         partner.Type,
		Country:             partner.Country,
		TotalTransfers:      totalTransfers,
		SuccessfulTransfers: successfulTransfers,
		FailedTransfers:     failedTransfers,
		TotalVolume:         sdk.NewCoin("NAMO", totalVolume),
		SuccessRate:         successRate,
		IsActive:            partner.IsActive,
	}

	return stats, nil
}

// emitSettlementEvent emits blockchain events for settlements
func (k Keeper) emitSettlementEvent(ctx context.Context, eventType string, settlement types.Settlement) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			eventType,
			sdk.NewAttribute(types.AttributeKeyTransferID, settlement.TransferId),
			sdk.NewAttribute(types.AttributeKeyPartnerID, settlement.PartnerId),
			sdk.NewAttribute(types.AttributeKeyAmount, settlement.SettledAmount.String()),
			sdk.NewAttribute(types.AttributeKeySettlementMethod, settlement.Method),
			sdk.NewAttribute(types.AttributeKeySettlementTime, settlement.SettlementTime.String()),
		),
	)
}

// emitPartnerEvent emits blockchain events for partners
func (k Keeper) emitPartnerEvent(ctx context.Context, eventType string, partner types.CorridorPartner) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			eventType,
			sdk.NewAttribute(types.AttributeKeyPartnerID, partner.Id),
			sdk.NewAttribute("partner_name", partner.Name),
			sdk.NewAttribute("partner_type", partner.Type),
			sdk.NewAttribute("country", partner.Country),
		),
	)
}