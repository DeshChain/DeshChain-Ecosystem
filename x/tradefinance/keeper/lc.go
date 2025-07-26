package keeper

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/tradefinance/types"
)

// IssueLc issues a new Letter of Credit
func (k Keeper) IssueLc(ctx sdk.Context, msg *types.MsgIssueLc) (string, string, error) {
	params := k.GetParams(ctx)
	
	// Check if module is enabled
	if !params.ModuleEnabled {
		return "", "", types.ErrModuleDisabled
	}

	// Validate parties
	applicant, found := k.GetTradeParty(ctx, msg.ApplicantId)
	if !found || applicant.PartyType != "importer" {
		return "", "", types.ErrPartyNotFound
	}

	beneficiary, found := k.GetTradeParty(ctx, msg.BeneficiaryId)
	if !found || beneficiary.PartyType != "exporter" {
		return "", "", types.ErrPartyNotFound
	}

	issuingBank, found := k.GetTradeParty(ctx, msg.IssuingBankId)
	if !found || issuingBank.PartyType != "bank" {
		return "", "", types.ErrPartyNotFound
	}

	// Validate amount
	if msg.Amount.Amount.LT(sdk.NewInt(int64(params.MinLcAmount))) {
		return "", "", types.ErrMinimumLCAmount
	}

	// Validate dates
	currentTime := ctx.BlockTime()
	if msg.ExpiryDate.Before(currentTime) {
		return "", "", types.ErrLCExpired
	}

	maxDuration := time.Duration(params.MaxLcDurationDays) * 24 * time.Hour
	if msg.ExpiryDate.Sub(currentTime) > maxDuration {
		return "", "", types.ErrMaximumLCDuration
	}

	// Check collateral
	issuingBankAddr, err := sdk.AccAddressFromBech32(msg.IssuingBank)
	if err != nil {
		return "", "", err
	}

	requiredCollateral := k.calculateRequiredCollateral(msg.Amount, params.CollateralRatio)
	if msg.Collateral.Amount.LT(requiredCollateral.Amount) {
		return "", "", types.ErrInsufficientCollateral
	}

	// Transfer collateral to module account
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, issuingBankAddr, types.ModuleName, sdk.NewCoins(msg.Collateral))
	if err != nil {
		return "", "", err
	}

	// Calculate and collect fees
	issuanceFee := k.calculateIssuanceFee(msg.Amount, params.Fees.LcIssuanceFee)
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, issuingBankAddr, types.ModuleName, sdk.NewCoins(issuanceFee))
	if err != nil {
		return "", "", err
	}

	// Generate LC ID and number
	lcID := k.GetNextLcID(ctx)
	lcIDStr := fmt.Sprintf("LC%08d", lcID)
	lcNumber := fmt.Sprintf("DESH/LC/%s/%d", currentTime.Format("2006"), lcID)

	// Create LC
	lc := types.LetterOfCredit{
		LcId:                lcIDStr,
		LcNumber:            lcNumber,
		Status:              "issued",
		ApplicantId:         msg.ApplicantId,
		BeneficiaryId:       msg.BeneficiaryId,
		IssuingBankId:       msg.IssuingBankId,
		AdvisingBankId:      msg.AdvisingBankId,
		ConfirmingBankId:    "", // Can be added later
		Amount:              msg.Amount,
		Currency:            msg.Amount.Denom,
		IssueDate:           currentTime,
		ExpiryDate:          msg.ExpiryDate,
		LatestShipmentDate:  msg.LatestShipmentDate,
		PaymentTerms:        msg.PaymentTerms,
		DeferredPaymentDays: msg.DeferredPaymentDays,
		Incoterms:           msg.Incoterms,
		PortOfLoading:       msg.PortOfLoading,
		PortOfDischarge:     msg.PortOfDischarge,
		PartialShipmentAllowed: msg.PartialShipmentAllowed,
		TransshipmentAllowed:   msg.TransshipmentAllowed,
		GoodsDescription:    msg.GoodsDescription,
		RequiredDocuments:   msg.RequiredDocuments,
		Collateral:          msg.Collateral,
		FeesPaid:            issuanceFee,
		CreatedAt:           currentTime,
		UpdatedAt:           currentTime,
	}

	// Save LC
	k.SetLetterOfCredit(ctx, lc)
	k.SetNextLcID(ctx, lcID+1)

	// Create indexes
	k.AddLcToPartyIndex(ctx, msg.ApplicantId, lcIDStr)
	k.AddLcToPartyIndex(ctx, msg.BeneficiaryId, lcIDStr)
	k.AddLcToPartyIndex(ctx, msg.IssuingBankId, lcIDStr)

	// Update stats
	k.UpdateLcStats(ctx, lc, true)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeLcIssued,
			sdk.NewAttribute(types.AttributeKeyLcId, lcIDStr),
			sdk.NewAttribute(types.AttributeKeyLcNumber, lcNumber),
			sdk.NewAttribute(types.AttributeKeyIssuingBank, msg.IssuingBankId),
			sdk.NewAttribute(types.AttributeKeyApplicant, msg.ApplicantId),
			sdk.NewAttribute(types.AttributeKeyBeneficiary, msg.BeneficiaryId),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.String()),
		),
	)

	return lcIDStr, lcNumber, nil
}

// AcceptLc accepts a Letter of Credit
func (k Keeper) AcceptLc(ctx sdk.Context, lcID string, beneficiaryAddr string) error {
	lc, found := k.GetLetterOfCredit(ctx, lcID)
	if !found {
		return types.ErrLCNotFound
	}

	// Validate status
	if lc.Status != "issued" {
		return types.ErrInvalidLCStatus
	}

	// Validate beneficiary
	beneficiaryPartyID := k.GetPartyIDByAddress(ctx, beneficiaryAddr)
	if beneficiaryPartyID != lc.BeneficiaryId {
		return types.ErrUnauthorized
	}

	// Update status
	lc.Status = "accepted"
	lc.UpdatedAt = ctx.BlockTime()
	k.SetLetterOfCredit(ctx, lc)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeLcAccepted,
			sdk.NewAttribute(types.AttributeKeyLcId, lcID),
			sdk.NewAttribute(types.AttributeKeyBeneficiary, lc.BeneficiaryId),
		),
	)

	return nil
}

// GetLetterOfCredit returns an LC by ID
func (k Keeper) GetLetterOfCredit(ctx sdk.Context, lcID string) (types.LetterOfCredit, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.LetterOfCreditPrefix)
	
	bz := store.Get([]byte(lcID))
	if bz == nil {
		return types.LetterOfCredit{}, false
	}

	var lc types.LetterOfCredit
	k.cdc.MustUnmarshal(bz, &lc)
	return lc, true
}

// SetLetterOfCredit saves an LC
func (k Keeper) SetLetterOfCredit(ctx sdk.Context, lc types.LetterOfCredit) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.LetterOfCreditPrefix)
	bz := k.cdc.MustMarshal(&lc)
	store.Set([]byte(lc.LcId), bz)
}

// GetAllLettersOfCredit returns all LCs
func (k Keeper) GetAllLettersOfCredit(ctx sdk.Context) []types.LetterOfCredit {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.LetterOfCreditPrefix)
	
	var lcs []types.LetterOfCredit
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var lc types.LetterOfCredit
		k.cdc.MustUnmarshal(iterator.Value(), &lc)
		lcs = append(lcs, lc)
	}
	
	return lcs
}

// GetLcsByParty returns LCs for a specific party
func (k Keeper) GetLcsByParty(ctx sdk.Context, partyID string) []types.LetterOfCredit {
	lcIDs := k.GetLcIDsByParty(ctx, partyID)
	
	var lcs []types.LetterOfCredit
	for _, lcID := range lcIDs {
		lc, found := k.GetLetterOfCredit(ctx, lcID)
		if found {
			lcs = append(lcs, lc)
		}
	}
	
	return lcs
}

// Helper functions

func (k Keeper) calculateRequiredCollateral(lcAmount sdk.Coin, collateralRatio uint64) sdk.Coin {
	// collateralRatio is in basis points (e.g., 11000 = 110%)
	requiredAmount := lcAmount.Amount.Mul(sdk.NewInt(int64(collateralRatio))).Quo(sdk.NewInt(10000))
	return sdk.NewCoin(lcAmount.Denom, requiredAmount)
}

func (k Keeper) calculateIssuanceFee(lcAmount sdk.Coin, feeRate uint64) sdk.Coin {
	// feeRate is in basis points
	feeAmount := lcAmount.Amount.Mul(sdk.NewInt(int64(feeRate))).Quo(sdk.NewInt(10000))
	return sdk.NewCoin("dinr", feeAmount) // Fees in DINR
}

// GetNextLcID returns the next LC ID
func (k Keeper) GetNextLcID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.NextLcIDKey)
	
	if bz == nil {
		return 1
	}
	
	return sdk.BigEndianToUint64(bz)
}

// SetNextLcID sets the next LC ID
func (k Keeper) SetNextLcID(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.NextLcIDKey, sdk.Uint64ToBigEndian(id))
}

// AddLcToPartyIndex adds an LC to a party's index
func (k Keeper) AddLcToPartyIndex(ctx sdk.Context, partyID, lcID string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.LcByPartyPrefix)
	key := append([]byte(partyID), []byte(lcID)...)
	store.Set(key, []byte{1})
}

// GetLcIDsByParty returns LC IDs for a party
func (k Keeper) GetLcIDsByParty(ctx sdk.Context, partyID string) []string {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.LcByPartyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte(partyID))
	defer iterator.Close()
	
	var lcIDs []string
	for ; iterator.Valid(); iterator.Next() {
		// Extract LC ID from key
		key := iterator.Key()
		lcID := string(key[len(partyID):])
		lcIDs = append(lcIDs, lcID)
	}
	
	return lcIDs
}

// UpdateLcStats updates statistics
func (k Keeper) UpdateLcStats(ctx sdk.Context, lc types.LetterOfCredit, isNew bool) {
	stats := k.GetTradeFinanceStats(ctx)
	
	if isNew {
		stats.TotalLcsIssued++
		stats.ActiveLcs++
		stats.TotalTradeValue = stats.TotalTradeValue.Add(lc.Amount)
		stats.TotalFeesCollected = stats.TotalFeesCollected.Add(lc.FeesPaid)
	}
	
	stats.LastUpdate = ctx.BlockTime()
	k.SetTradeFinanceStats(ctx, stats)
}