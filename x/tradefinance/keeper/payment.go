package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/tradefinance/types"
)

// RequestPayment creates a payment request under an LC
func (k Keeper) RequestPayment(ctx sdk.Context, lcID string, beneficiaryAddr string, amount sdk.Coin) (string, error) {
	// Get LC
	lc, found := k.GetLetterOfCredit(ctx, lcID)
	if !found {
		return "", types.ErrLCNotFound
	}

	// Validate LC status
	if lc.Status != "documents_verified" && lc.Status != "paid" {
		return "", types.ErrInvalidLCStatus
	}

	// Check if LC has expired
	if ctx.BlockTime().After(lc.ExpiryDate) {
		return "", types.ErrLCExpired
	}

	// Validate beneficiary
	beneficiaryPartyID := k.GetPartyIDByAddress(ctx, beneficiaryAddr)
	if beneficiaryPartyID != lc.BeneficiaryId {
		return "", types.ErrUnauthorized
	}

	// Validate amount doesn't exceed LC balance
	totalPaid := k.GetTotalPaidAmount(ctx, lcID)
	remainingAmount := lc.Amount.Sub(totalPaid)
	if amount.Amount.GT(remainingAmount.Amount) {
		return "", types.ErrAmountExceedsLC
	}

	// Generate payment instruction ID
	instructionID := k.GetNextInstructionID(ctx)
	instructionIDStr := fmt.Sprintf("PAY%08d", instructionID)

	// Determine due date based on payment terms
	dueDate := ctx.BlockTime()
	if lc.PaymentTerms == "deferred" {
		dueDate = dueDate.AddDate(0, 0, int(lc.DeferredPaymentDays))
	}

	// Create payment instruction
	instruction := types.PaymentInstruction{
		InstructionId:   instructionIDStr,
		LcId:            lcID,
		Payer:           lc.ApplicantId,
		Payee:           lc.BeneficiaryId,
		Amount:          amount,
		PaymentType:     lc.PaymentTerms,
		DueDate:         dueDate,
		Status:          "pending",
		TransactionHash: "",
		PaidAt:          ctx.BlockTime(), // Will be updated when actually paid
	}

	// Save payment instruction
	k.SetPaymentInstruction(ctx, instruction)
	k.AddPaymentToLcIndex(ctx, lcID, instructionIDStr)
	k.SetNextInstructionID(ctx, instructionID+1)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypePaymentRequested,
			sdk.NewAttribute(types.AttributeKeyPaymentInstructionId, instructionIDStr),
			sdk.NewAttribute(types.AttributeKeyLcId, lcID),
			sdk.NewAttribute(types.AttributeKeyBeneficiary, beneficiaryPartyID),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	)

	return instructionIDStr, nil
}

// MakePayment processes a payment instruction
func (k Keeper) MakePayment(ctx sdk.Context, instructionID string, payerAddr string) error {
	// Get payment instruction
	instruction, found := k.GetPaymentInstruction(ctx, instructionID)
	if !found {
		return types.ErrPaymentInstructionNotFound
	}

	// Check if already paid
	if instruction.Status == "completed" {
		return types.ErrPaymentAlreadyCompleted
	}

	// Get LC
	lc, found := k.GetLetterOfCredit(ctx, instruction.LcId)
	if !found {
		return types.ErrLCNotFound
	}

	// Validate payer
	payerPartyID := k.GetPartyIDByAddress(ctx, payerAddr)
	
	// Check if payer is the applicant or issuing bank
	if payerPartyID != lc.ApplicantId && payerPartyID != lc.IssuingBankId {
		return types.ErrUnauthorized
	}

	// Get addresses
	payerAccAddr, err := sdk.AccAddressFromBech32(payerAddr)
	if err != nil {
		return err
	}

	beneficiaryParty, found := k.GetTradeParty(ctx, lc.BeneficiaryId)
	if !found {
		return types.ErrPartyNotFound
	}

	beneficiaryAddr, err := sdk.AccAddressFromBech32(beneficiaryParty.DeshAddress)
	if err != nil {
		return err
	}

	// Transfer payment from payer to beneficiary
	err = k.bankKeeper.SendCoins(ctx, payerAccAddr, beneficiaryAddr, sdk.NewCoins(instruction.Amount))
	if err != nil {
		return err
	}

	// Update payment instruction
	instruction.Status = "completed"
	instruction.TransactionHash = fmt.Sprintf("TX_%s_%d", ctx.TxBytes(), ctx.BlockHeight())
	instruction.PaidAt = ctx.BlockTime()
	k.SetPaymentInstruction(ctx, instruction)

	// Check if LC is fully paid
	totalPaid := k.GetTotalPaidAmount(ctx, lc.LcId)
	if totalPaid.Equal(lc.Amount) {
		// Release collateral back to issuing bank
		issuingBankParty, found := k.GetTradeParty(ctx, lc.IssuingBankId)
		if found {
			issuingBankAddr, err := sdk.AccAddressFromBech32(issuingBankParty.DeshAddress)
			if err == nil {
				k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, issuingBankAddr, sdk.NewCoins(lc.Collateral))
			}
		}

		// Update LC status
		lc.Status = "paid"
		lc.UpdatedAt = ctx.BlockTime()
		k.SetLetterOfCredit(ctx, lc)

		// Update stats
		stats := k.GetTradeFinanceStats(ctx)
		stats.CompletedTrades++
		stats.ActiveLcs--
		k.SetTradeFinanceStats(ctx, stats)
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypePaymentCompleted,
			sdk.NewAttribute(types.AttributeKeyPaymentInstructionId, instructionID),
			sdk.NewAttribute(types.AttributeKeyLcId, instruction.LcId),
			sdk.NewAttribute(types.AttributeKeyPayer, payerPartyID),
			sdk.NewAttribute(types.AttributeKeyPayee, instruction.Payee),
			sdk.NewAttribute(types.AttributeKeyAmount, instruction.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyTransactionHash, instruction.TransactionHash),
		),
	)

	return nil
}

// GetPaymentInstruction returns a payment instruction by ID
func (k Keeper) GetPaymentInstruction(ctx sdk.Context, instructionID string) (types.PaymentInstruction, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PaymentInstructionPrefix)
	
	bz := store.Get([]byte(instructionID))
	if bz == nil {
		return types.PaymentInstruction{}, false
	}

	var instruction types.PaymentInstruction
	k.cdc.MustUnmarshal(bz, &instruction)
	return instruction, true
}

// SetPaymentInstruction saves a payment instruction
func (k Keeper) SetPaymentInstruction(ctx sdk.Context, instruction types.PaymentInstruction) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PaymentInstructionPrefix)
	bz := k.cdc.MustMarshal(&instruction)
	store.Set([]byte(instruction.InstructionId), bz)
}

// GetPaymentsByLc returns all payment instructions for an LC
func (k Keeper) GetPaymentsByLc(ctx sdk.Context, lcID string) []types.PaymentInstruction {
	paymentIDs := k.GetPaymentIDsByLc(ctx, lcID)
	
	var payments []types.PaymentInstruction
	for _, paymentID := range paymentIDs {
		payment, found := k.GetPaymentInstruction(ctx, paymentID)
		if found {
			payments = append(payments, payment)
		}
	}
	
	return payments
}

// GetTotalPaidAmount calculates the total amount paid for an LC
func (k Keeper) GetTotalPaidAmount(ctx sdk.Context, lcID string) sdk.Coin {
	payments := k.GetPaymentsByLc(ctx, lcID)
	
	total := sdk.NewCoin("dinr", sdk.ZeroInt())
	for _, payment := range payments {
		if payment.Status == "completed" {
			total = total.Add(payment.Amount)
		}
	}
	
	return total
}

// AddPaymentToLcIndex adds a payment to LC's index
func (k Keeper) AddPaymentToLcIndex(ctx sdk.Context, lcID, paymentID string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PaymentByLcPrefix)
	key := append([]byte(lcID), []byte(paymentID)...)
	store.Set(key, []byte{1})
}

// GetPaymentIDsByLc returns payment IDs for an LC
func (k Keeper) GetPaymentIDsByLc(ctx sdk.Context, lcID string) []string {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PaymentByLcPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte(lcID))
	defer iterator.Close()
	
	var paymentIDs []string
	for ; iterator.Valid(); iterator.Next() {
		// Extract payment ID from key
		key := iterator.Key()
		paymentID := string(key[len(lcID):])
		paymentIDs = append(paymentIDs, paymentID)
	}
	
	return paymentIDs
}

// GetNextInstructionID returns the next payment instruction ID
func (k Keeper) GetNextInstructionID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.NextInstructionIDKey)
	
	if bz == nil {
		return 1
	}
	
	return sdk.BigEndianToUint64(bz)
}

// SetNextInstructionID sets the next payment instruction ID
func (k Keeper) SetNextInstructionID(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.NextInstructionIDKey, sdk.Uint64ToBigEndian(id))
}

// EstimateFees estimates the fees for LC operations
func (k Keeper) EstimateFees(ctx sdk.Context, lcAmount sdk.Coin, paymentTerms string, insuranceRequired bool) (issuanceFee, documentFees, insuranceFee, totalFees sdk.Coin, processingTimeHours uint64) {
	params := k.GetParams(ctx)
	
	// Calculate issuance fee
	issuanceFee = k.calculateIssuanceFee(lcAmount, params.Fees.LcIssuanceFee)
	
	// Calculate document fees (fixed per document type)
	numDocs := uint64(len(params.SupportedDocumentTypes))
	documentFees = sdk.NewCoin("dinr", sdk.NewInt(int64(params.Fees.DocumentVerificationFee * numDocs)))
	
	// Calculate insurance fee if required
	if insuranceRequired {
		insuranceAmount := lcAmount.Amount.Mul(sdk.NewInt(int64(params.Fees.InsuranceProcessingFee))).Quo(sdk.NewInt(10000))
		insuranceFee = sdk.NewCoin("dinr", insuranceAmount)
	} else {
		insuranceFee = sdk.NewCoin("dinr", sdk.ZeroInt())
	}
	
	// Total fees
	totalFees = issuanceFee.Add(documentFees).Add(insuranceFee)
	
	// Estimate processing time based on payment terms
	switch paymentTerms {
	case "sight":
		processingTimeHours = 2 // 2 hours for sight LC
	case "deferred":
		processingTimeHours = 4 // 4 hours for deferred payment
	default:
		processingTimeHours = 3 // 3 hours average
	}
	
	return
}