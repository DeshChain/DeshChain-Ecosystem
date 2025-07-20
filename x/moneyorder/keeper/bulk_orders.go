package keeper

import (
	"fmt"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/deshchain/x/moneyorder/types"
)

// BulkOrderManager handles bulk money order operations for businesses
type BulkOrderManager struct {
	keeper *Keeper
}

// NewBulkOrderManager creates a new bulk order manager
func NewBulkOrderManager(keeper *Keeper) *BulkOrderManager {
	return &BulkOrderManager{
		keeper: keeper,
	}
}

// CreateBulkOrder creates multiple money orders in a single transaction
func (bom *BulkOrderManager) CreateBulkOrder(
	ctx sdk.Context,
	businessAddress string,
	orders []types.BulkOrderItem,
	batchMetadata types.BulkOrderMetadata,
) (*types.BulkOrderResult, error) {
	// Validate business account
	business, found := bom.keeper.GetBusinessAccount(ctx, businessAddress)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrUserNotFound, "business account not found")
	}

	if !business.IsActive || !business.BulkOrdersEnabled {
		return nil, sdkerrors.Wrap(types.ErrUnauthorized, "business not authorized for bulk orders")
	}

	// Validate batch size limits
	maxBatchSize := bom.keeper.GetParams(ctx).MaxBulkOrderSize
	if len(orders) > int(maxBatchSize) {
		return nil, sdkerrors.Wrapf(types.ErrInvalidInput, "batch size %d exceeds maximum %d", len(orders), maxBatchSize)
	}

	if len(orders) == 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalidInput, "no orders provided")
	}

	// Calculate total amount and validate individual orders
	totalAmount := sdk.ZeroInt()
	var validatedOrders []types.ValidatedBulkOrderItem
	
	for i, order := range orders {
		// Validate individual order
		validatedOrder, err := bom.validateBulkOrderItem(ctx, businessAddress, order, i)
		if err != nil {
			return nil, sdkerrors.Wrapf(err, "validation failed for order %d", i)
		}
		
		validatedOrders = append(validatedOrders, validatedOrder)
		totalAmount = totalAmount.Add(validatedOrder.Amount)
	}

	// Check business balance and limits
	if err := bom.validateBusinessLimits(ctx, business, totalAmount, len(orders)); err != nil {
		return nil, err
	}

	// Generate bulk order ID
	bulkOrderID := bom.generateBulkOrderID(ctx, businessAddress)

	// Create bulk order record
	bulkOrder := types.BulkOrder{
		ID:              bulkOrderID,
		BusinessAddress: businessAddress,
		TotalOrders:     int64(len(orders)),
		TotalAmount:     totalAmount,
		Status:          types.BulkOrderStatus_PROCESSING,
		CreatedAt:       time.Now(),
		Metadata:        batchMetadata,
		ProcessingStats: types.BulkOrderProcessingStats{
			TotalOrders:    int64(len(orders)),
			ProcessedOrders: 0,
			SuccessfulOrders: 0,
			FailedOrders:    0,
			StartTime:      time.Now(),
		},
	}

	// Save bulk order
	bom.saveBulkOrder(ctx, bulkOrder)

	// Process orders in batches
	result := &types.BulkOrderResult{
		BulkOrderID:     bulkOrderID,
		TotalOrders:     int64(len(orders)),
		SuccessfulOrders: []types.OrderResult{},
		FailedOrders:    []types.OrderFailure{},
		ProcessingTime:  time.Duration(0),
	}

	startTime := time.Now()
	batchSize := bom.keeper.GetParams(ctx).BulkProcessingBatchSize

	for i := 0; i < len(validatedOrders); i += int(batchSize) {
		end := i + int(batchSize)
		if end > len(validatedOrders) {
			end = len(validatedOrders)
		}

		batchOrders := validatedOrders[i:end]
		batchResults := bom.processBatch(ctx, businessAddress, bulkOrderID, batchOrders, i)

		// Aggregate results
		result.SuccessfulOrders = append(result.SuccessfulOrders, batchResults.SuccessfulOrders...)
		result.FailedOrders = append(result.FailedOrders, batchResults.FailedOrders...)

		// Update bulk order progress
		bulkOrder.ProcessingStats.ProcessedOrders += int64(len(batchOrders))
		bulkOrder.ProcessingStats.SuccessfulOrders += int64(len(batchResults.SuccessfulOrders))
		bulkOrder.ProcessingStats.FailedOrders += int64(len(batchResults.FailedOrders))
		bom.saveBulkOrder(ctx, bulkOrder)

		// Emit progress event
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeBulkOrderProgress,
				sdk.NewAttribute(types.AttributeKeyBulkOrderID, bulkOrderID),
				sdk.NewAttribute(types.AttributeKeyProcessedOrders, fmt.Sprintf("%d", bulkOrder.ProcessingStats.ProcessedOrders)),
				sdk.NewAttribute(types.AttributeKeyTotalOrders, fmt.Sprintf("%d", bulkOrder.ProcessingStats.TotalOrders)),
			),
		)
	}

	// Finalize bulk order
	processingTime := time.Since(startTime)
	result.ProcessingTime = processingTime

	bulkOrder.Status = types.BulkOrderStatus_COMPLETED
	bulkOrder.CompletedAt = time.Now()
	bulkOrder.ProcessingStats.EndTime = time.Now()
	bulkOrder.ProcessingStats.ProcessingDuration = processingTime
	bom.saveBulkOrder(ctx, bulkOrder)

	// Update business statistics
	bom.updateBusinessStats(ctx, businessAddress, result)

	// Calculate and deduct fees
	totalFees := bom.calculateBulkOrderFees(ctx, business, result)
	if err := bom.deductBulkOrderFees(ctx, businessAddress, totalFees); err != nil {
		// Log error but don't fail the transaction
		ctx.Logger().Error("Failed to deduct bulk order fees", "error", err)
	}

	// Emit completion event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeBulkOrderCompleted,
			sdk.NewAttribute(types.AttributeKeyBulkOrderID, bulkOrderID),
			sdk.NewAttribute(types.AttributeKeyBusinessAddress, businessAddress),
			sdk.NewAttribute(types.AttributeKeyTotalOrders, fmt.Sprintf("%d", result.TotalOrders)),
			sdk.NewAttribute(types.AttributeKeySuccessfulOrders, fmt.Sprintf("%d", len(result.SuccessfulOrders))),
			sdk.NewAttribute(types.AttributeKeyFailedOrders, fmt.Sprintf("%d", len(result.FailedOrders))),
			sdk.NewAttribute(types.AttributeKeyProcessingTime, processingTime.String()),
		),
	)

	return result, nil
}

// GetBulkOrderStatus returns the status of a bulk order
func (bom *BulkOrderManager) GetBulkOrderStatus(
	ctx sdk.Context,
	bulkOrderID string,
) (*types.BulkOrderStatusResponse, error) {
	bulkOrder, found := bom.getBulkOrder(ctx, bulkOrderID)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrOrderNotFound, "bulk order not found")
	}

	// Get individual order details
	orders := bom.getBulkOrderItems(ctx, bulkOrderID)

	return &types.BulkOrderStatusResponse{
		BulkOrder:    bulkOrder,
		Orders:       orders,
		Summary:      bom.generateBulkOrderSummary(bulkOrder, orders),
	}, nil
}

// GetBusinessBulkOrders returns all bulk orders for a business
func (bom *BulkOrderManager) GetBusinessBulkOrders(
	ctx sdk.Context,
	businessAddress string,
	limit int,
	offset int,
) ([]types.BulkOrder, error) {
	var bulkOrders []types.BulkOrder
	
	store := prefix.NewStore(ctx.KVStore(bom.keeper.storeKey), types.BulkOrderPrefix)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	count := 0
	skipped := 0

	for ; iterator.Valid(); iterator.Next() {
		var bulkOrder types.BulkOrder
		bom.keeper.cdc.MustUnmarshal(iterator.Value(), &bulkOrder)
		
		if bulkOrder.BusinessAddress == businessAddress {
			if skipped < offset {
				skipped++
				continue
			}
			
			if count >= limit {
				break
			}
			
			bulkOrders = append(bulkOrders, bulkOrder)
			count++
		}
	}

	return bulkOrders, nil
}

// CancelBulkOrder cancels a bulk order if it's still processing
func (bom *BulkOrderManager) CancelBulkOrder(
	ctx sdk.Context,
	businessAddress string,
	bulkOrderID string,
	reason string,
) error {
	bulkOrder, found := bom.getBulkOrder(ctx, bulkOrderID)
	if !found {
		return sdkerrors.Wrap(types.ErrOrderNotFound, "bulk order not found")
	}

	if bulkOrder.BusinessAddress != businessAddress {
		return sdkerrors.Wrap(types.ErrUnauthorized, "not authorized to cancel this bulk order")
	}

	if bulkOrder.Status != types.BulkOrderStatus_PROCESSING {
		return sdkerrors.Wrap(types.ErrInvalidOrderStatus, "bulk order cannot be cancelled")
	}

	// Update bulk order status
	bulkOrder.Status = types.BulkOrderStatus_CANCELLED
	bulkOrder.CancelledAt = time.Now()
	bulkOrder.CancellationReason = reason
	bom.saveBulkOrder(ctx, bulkOrder)

	// Cancel pending individual orders
	orders := bom.getBulkOrderItems(ctx, bulkOrderID)
	for _, order := range orders {
		if order.Status == types.MoneyOrderStatus_PENDING || order.Status == types.MoneyOrderStatus_PROCESSING {
			bom.keeper.CancelMoneyOrder(ctx, order.OrderID, "Bulk order cancelled")
		}
	}

	// Emit cancellation event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeBulkOrderCancelled,
			sdk.NewAttribute(types.AttributeKeyBulkOrderID, bulkOrderID),
			sdk.NewAttribute(types.AttributeKeyBusinessAddress, businessAddress),
			sdk.NewAttribute(types.AttributeKeyReason, reason),
		),
	)

	return nil
}

// ValidateBulkOrderTemplate validates a bulk order template before processing
func (bom *BulkOrderManager) ValidateBulkOrderTemplate(
	ctx sdk.Context,
	businessAddress string,
	template types.BulkOrderTemplate,
) (*types.BulkOrderValidationResult, error) {
	business, found := bom.keeper.GetBusinessAccount(ctx, businessAddress)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrUserNotFound, "business account not found")
	}

	validationResult := &types.BulkOrderValidationResult{
		IsValid:      true,
		TotalOrders:  int64(len(template.Orders)),
		TotalAmount:  sdk.ZeroInt(),
		Warnings:     []string{},
		Errors:       []string{},
		ValidOrders:  0,
		InvalidOrders: 0,
	}

	// Validate each order in the template
	for i, order := range template.Orders {
		if err := bom.validateBulkOrderItemBasic(order); err != nil {
			validationResult.Errors = append(validationResult.Errors, 
				fmt.Sprintf("Order %d: %s", i+1, err.Error()))
			validationResult.InvalidOrders++
		} else {
			validationResult.ValidOrders++
			validationResult.TotalAmount = validationResult.TotalAmount.Add(order.Amount)
		}
	}

	// Check business limits
	if validationResult.TotalOrders > int64(bom.keeper.GetParams(ctx).MaxBulkOrderSize) {
		validationResult.Warnings = append(validationResult.Warnings, 
			"Batch size exceeds recommended maximum")
	}

	if validationResult.TotalAmount.GT(business.DailyLimit) {
		validationResult.Errors = append(validationResult.Errors, 
			"Total amount exceeds daily limit")
		validationResult.IsValid = false
	}

	if validationResult.InvalidOrders > 0 {
		validationResult.IsValid = false
	}

	return validationResult, nil
}

// Helper functions

func (bom *BulkOrderManager) validateBulkOrderItem(
	ctx sdk.Context,
	businessAddress string,
	order types.BulkOrderItem,
	index int,
) (types.ValidatedBulkOrderItem, error) {
	// Basic validation
	if err := bom.validateBulkOrderItemBasic(order); err != nil {
		return types.ValidatedBulkOrderItem{}, err
	}

	// Extended validation
	validated := types.ValidatedBulkOrderItem{
		Index:           index,
		OriginalOrder:   order,
		Amount:          order.Amount,
		RecipientAddress: order.RecipientAddress,
		SenderAddress:   businessAddress,
		Memo:            order.Memo,
		Priority:        order.Priority,
	}

	// Validate recipient address
	if err := bom.keeper.ValidateAddress(ctx, order.RecipientAddress); err != nil {
		return validated, sdkerrors.Wrap(types.ErrInvalidRecipient, "invalid recipient address")
	}

	// Check if recipient is not blacklisted
	if bom.keeper.IsAddressBlacklisted(ctx, order.RecipientAddress) {
		return validated, sdkerrors.Wrap(types.ErrUserBlacklisted, "recipient is blacklisted")
	}

	return validated, nil
}

func (bom *BulkOrderManager) validateBulkOrderItemBasic(order types.BulkOrderItem) error {
	if order.Amount.IsZero() || order.Amount.IsNegative() {
		return sdkerrors.Wrap(types.ErrInvalidAmount, "amount must be positive")
	}

	if order.RecipientAddress == "" {
		return sdkerrors.Wrap(types.ErrInvalidRecipient, "recipient address is required")
	}

	if len(order.Memo) > 256 {
		return sdkerrors.Wrap(types.ErrInvalidInput, "memo too long")
	}

	return nil
}

func (bom *BulkOrderManager) validateBusinessLimits(
	ctx sdk.Context,
	business types.BusinessAccount,
	totalAmount sdk.Int,
	orderCount int,
) error {
	// Check balance
	balance := bom.keeper.GetBalance(ctx, business.Address)
	if balance.LT(totalAmount) {
		return sdkerrors.Wrap(types.ErrInsufficientFunds, "insufficient balance for bulk order")
	}

	// Check daily limits
	todayUsage := bom.getDailyUsage(ctx, business.Address)
	if todayUsage.Add(totalAmount).GT(business.DailyLimit) {
		return sdkerrors.Wrap(types.ErrDailyLimitExceeded, "bulk order would exceed daily limit")
	}

	// Check hourly rate limits
	hourlyOrderCount := bom.getHourlyOrderCount(ctx, business.Address)
	maxHourlyOrders := bom.keeper.GetParams(ctx).MaxHourlyBulkOrders
	if hourlyOrderCount+int64(orderCount) > maxHourlyOrders {
		return sdkerrors.Wrap(types.ErrInvalidInput, "would exceed hourly order limit")
	}

	return nil
}

func (bom *BulkOrderManager) processBatch(
	ctx sdk.Context,
	businessAddress string,
	bulkOrderID string,
	orders []types.ValidatedBulkOrderItem,
	batchStartIndex int,
) types.BatchProcessingResult {
	result := types.BatchProcessingResult{
		SuccessfulOrders: []types.OrderResult{},
		FailedOrders:     []types.OrderFailure{},
	}

	for _, order := range orders {
		orderResult, err := bom.processIndividualOrder(ctx, businessAddress, bulkOrderID, order)
		if err != nil {
			result.FailedOrders = append(result.FailedOrders, types.OrderFailure{
				Index:   order.Index,
				OrderID: orderResult.OrderID,
				Error:   err.Error(),
				Amount:  order.Amount,
			})
		} else {
			result.SuccessfulOrders = append(result.SuccessfulOrders, orderResult)
		}
	}

	return result
}

func (bom *BulkOrderManager) processIndividualOrder(
	ctx sdk.Context,
	businessAddress string,
	bulkOrderID string,
	order types.ValidatedBulkOrderItem,
) (types.OrderResult, error) {
	// Create individual money order
	moneyOrder := types.MoneyOrder{
		SenderAddress:    businessAddress,
		RecipientAddress: order.RecipientAddress,
		Amount:           order.Amount,
		Memo:             order.Memo,
		CreatedAt:        time.Now(),
		ExpiresAt:        time.Now().Add(24 * time.Hour), // 24 hour expiry
		Status:           types.MoneyOrderStatus_PENDING,
		BulkOrderID:      bulkOrderID,
		Priority:         order.Priority,
	}

	// Generate order ID
	orderID := bom.keeper.GenerateMoneyOrderID(ctx)
	moneyOrder.OrderID = orderID

	// Save money order
	bom.keeper.SetMoneyOrder(ctx, moneyOrder)

	// Create bulk order item record
	bulkOrderItem := types.BulkOrderItem{
		BulkOrderID:      bulkOrderID,
		OrderID:          orderID,
		Index:            order.Index,
		RecipientAddress: order.RecipientAddress,
		Amount:           order.Amount,
		Status:           types.MoneyOrderStatus_PENDING,
		CreatedAt:        time.Now(),
		Memo:             order.Memo,
		Priority:         order.Priority,
	}

	bom.saveBulkOrderItem(ctx, bulkOrderItem)

	return types.OrderResult{
		Index:            order.Index,
		OrderID:          orderID,
		RecipientAddress: order.RecipientAddress,
		Amount:           order.Amount,
		Status:           types.MoneyOrderStatus_PENDING,
		CreatedAt:        time.Now(),
	}, nil
}

func (bom *BulkOrderManager) generateBulkOrderID(ctx sdk.Context, businessAddress string) string {
	timestamp := time.Now().Unix()
	sequence := bom.getNextBulkOrderSequence(ctx, businessAddress)
	return fmt.Sprintf("BULK_%s_%d_%d", 
		strings.ToUpper(businessAddress[len(businessAddress)-6:]), timestamp, sequence)
}

// Additional helper functions would be implemented here for:
// - saveBulkOrder
// - getBulkOrder
// - saveBulkOrderItem
// - getBulkOrderItems
// - updateBusinessStats
// - calculateBulkOrderFees
// - deductBulkOrderFees
// - getDailyUsage
// - getHourlyOrderCount
// - getNextBulkOrderSequence
// - generateBulkOrderSummary

// Placeholder implementations
func (bom *BulkOrderManager) saveBulkOrder(ctx sdk.Context, bulkOrder types.BulkOrder) {
	store := prefix.NewStore(ctx.KVStore(bom.keeper.storeKey), types.BulkOrderPrefix)
	bz := bom.keeper.cdc.MustMarshal(&bulkOrder)
	store.Set([]byte(bulkOrder.ID), bz)
}

func (bom *BulkOrderManager) getBulkOrder(ctx sdk.Context, bulkOrderID string) (types.BulkOrder, bool) {
	store := prefix.NewStore(ctx.KVStore(bom.keeper.storeKey), types.BulkOrderPrefix)
	bz := store.Get([]byte(bulkOrderID))
	if bz == nil {
		return types.BulkOrder{}, false
	}

	var bulkOrder types.BulkOrder
	bom.keeper.cdc.MustUnmarshal(bz, &bulkOrder)
	return bulkOrder, true
}

func (bom *BulkOrderManager) saveBulkOrderItem(ctx sdk.Context, item types.BulkOrderItem) {
	store := prefix.NewStore(ctx.KVStore(bom.keeper.storeKey), types.BulkOrderItemPrefix)
	key := fmt.Sprintf("%s_%s", item.BulkOrderID, item.OrderID)
	bz := bom.keeper.cdc.MustMarshal(&item)
	store.Set([]byte(key), bz)
}

func (bom *BulkOrderManager) getBulkOrderItems(ctx sdk.Context, bulkOrderID string) []types.BulkOrderItem {
	var items []types.BulkOrderItem
	store := prefix.NewStore(ctx.KVStore(bom.keeper.storeKey), types.BulkOrderItemPrefix)
	prefixKey := bulkOrderID + "_"
	iterator := sdk.KVStorePrefixIterator(store, []byte(prefixKey))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var item types.BulkOrderItem
		bom.keeper.cdc.MustUnmarshal(iterator.Value(), &item)
		items = append(items, item)
	}

	return items
}

func (bom *BulkOrderManager) updateBusinessStats(ctx sdk.Context, businessAddress string, result *types.BulkOrderResult) {
	// Implementation for updating business statistics
}

func (bom *BulkOrderManager) calculateBulkOrderFees(ctx sdk.Context, business types.BusinessAccount, result *types.BulkOrderResult) sdk.Int {
	// Implementation for calculating bulk order fees
	return sdk.ZeroInt()
}

func (bom *BulkOrderManager) deductBulkOrderFees(ctx sdk.Context, businessAddress string, fees sdk.Int) error {
	// Implementation for deducting fees
	return nil
}

func (bom *BulkOrderManager) getDailyUsage(ctx sdk.Context, businessAddress string) sdk.Int {
	// Implementation for getting daily usage
	return sdk.ZeroInt()
}

func (bom *BulkOrderManager) getHourlyOrderCount(ctx sdk.Context, businessAddress string) int64 {
	// Implementation for getting hourly order count
	return 0
}

func (bom *BulkOrderManager) getNextBulkOrderSequence(ctx sdk.Context, businessAddress string) int64 {
	// Implementation for getting next sequence number
	return 1
}

func (bom *BulkOrderManager) generateBulkOrderSummary(bulkOrder types.BulkOrder, orders []types.BulkOrderItem) types.BulkOrderSummary {
	// Implementation for generating summary
	return types.BulkOrderSummary{}
}