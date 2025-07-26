package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/tradefinance/types"
)

// SubmitDocuments submits trade documents for an LC
func (k Keeper) SubmitDocuments(ctx sdk.Context, lcID string, submitterAddr string, documents []types.DocumentSubmission) ([]string, error) {
	// Get LC
	lc, found := k.GetLetterOfCredit(ctx, lcID)
	if !found {
		return nil, types.ErrLCNotFound
	}

	// Validate LC status
	if lc.Status != "accepted" && lc.Status != "documents_presented" {
		return nil, types.ErrInvalidLCStatus
	}

	// Validate submitter is beneficiary
	submitterPartyID := k.GetPartyIDByAddress(ctx, submitterAddr)
	if submitterPartyID != lc.BeneficiaryId {
		return nil, types.ErrUnauthorized
	}

	// Check if past shipment deadline
	if ctx.BlockTime().After(lc.LatestShipmentDate) {
		return nil, types.ErrShipmentDeadlinePassed
	}

	params := k.GetParams(ctx)
	var documentIDs []string

	// Process each document
	for _, docSubmission := range documents {
		// Validate document type
		validType := false
		for _, supportedType := range params.SupportedDocumentTypes {
			if docSubmission.DocumentType == supportedType {
				validType = true
				break
			}
		}
		if !validType {
			return nil, types.ErrInvalidDocumentType
		}

		// Generate document ID
		docID := k.GetNextDocumentID(ctx)
		docIDStr := fmt.Sprintf("DOC%08d", docID)

		// Create document
		doc := types.TradeDocument{
			DocumentId:     docIDStr,
			LcId:           lcID,
			DocumentType:   docSubmission.DocumentType,
			DocumentHash:   docSubmission.DocumentHash, // IPFS hash
			Issuer:         docSubmission.Issuer,
			IssuedDate:     docSubmission.IssuedDate,
			IsVerified:     false,
			VerifiedBy:     "",
			VerifiedAt:     ctx.BlockTime(), // Will be updated when verified
			Status:         "submitted",
			RejectionReason: "",
		}

		// Save document
		k.SetTradeDocument(ctx, doc)
		k.AddDocumentToLcIndex(ctx, lcID, docIDStr)
		documentIDs = append(documentIDs, docIDStr)

		// Emit event
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeDocumentSubmitted,
				sdk.NewAttribute(types.AttributeKeyDocumentId, docIDStr),
				sdk.NewAttribute(types.AttributeKeyLcId, lcID),
				sdk.NewAttribute(types.AttributeKeyDocumentType, doc.DocumentType),
				sdk.NewAttribute(types.AttributeKeySubmitter, submitterAddr),
			),
		)

		k.SetNextDocumentID(ctx, docID+1)
	}

	// Update LC status
	if lc.Status == "accepted" {
		lc.Status = "documents_presented"
		lc.UpdatedAt = ctx.BlockTime()
		k.SetLetterOfCredit(ctx, lc)
	}

	// Update stats
	stats := k.GetTradeFinanceStats(ctx)
	stats.DocumentsVerified += uint64(len(documentIDs))
	k.SetTradeFinanceStats(ctx, stats)

	return documentIDs, nil
}

// VerifyDocument verifies or rejects a submitted document
func (k Keeper) VerifyDocument(ctx sdk.Context, documentID string, verifierAddr string, approved bool, rejectionReason string) error {
	// Get document
	doc, found := k.GetTradeDocument(ctx, documentID)
	if !found {
		return types.ErrDocumentNotFound
	}

	// Check if already verified
	if doc.Status == "verified" || doc.Status == "rejected" {
		return types.ErrDocumentAlreadyVerified
	}

	// Get LC
	lc, found := k.GetLetterOfCredit(ctx, doc.LcId)
	if !found {
		return types.ErrLCNotFound
	}

	// Validate verifier is issuing bank
	verifierPartyID := k.GetPartyIDByAddress(ctx, verifierAddr)
	if verifierPartyID != lc.IssuingBankId {
		return types.ErrUnauthorized
	}

	// Update document
	doc.IsVerified = approved
	doc.VerifiedBy = verifierPartyID
	doc.VerifiedAt = ctx.BlockTime()
	
	if approved {
		doc.Status = "verified"
		doc.RejectionReason = ""
	} else {
		doc.Status = "rejected"
		doc.RejectionReason = rejectionReason
	}

	k.SetTradeDocument(ctx, doc)

	// Check if all required documents are verified
	if approved {
		allVerified := k.CheckAllDocumentsVerified(ctx, lc.LcId, lc.RequiredDocuments)
		if allVerified {
			// Update LC status to allow payment
			lc.Status = "documents_verified"
			lc.UpdatedAt = ctx.BlockTime()
			k.SetLetterOfCredit(ctx, lc)
		}
	} else {
		// Mark LC as discrepant
		if lc.Status != "discrepant" {
			lc.Status = "discrepant"
			lc.UpdatedAt = ctx.BlockTime()
			k.SetLetterOfCredit(ctx, lc)
		}
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDocumentVerified,
			sdk.NewAttribute(types.AttributeKeyDocumentId, documentID),
			sdk.NewAttribute(types.AttributeKeyLcId, doc.LcId),
			sdk.NewAttribute(types.AttributeKeyVerifier, verifierPartyID),
			sdk.NewAttribute(types.AttributeKeyApproved, fmt.Sprintf("%t", approved)),
			sdk.NewAttribute(types.AttributeKeyRejectionReason, rejectionReason),
		),
	)

	return nil
}

// GetTradeDocument returns a document by ID
func (k Keeper) GetTradeDocument(ctx sdk.Context, documentID string) (types.TradeDocument, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.TradeDocumentPrefix)
	
	bz := store.Get([]byte(documentID))
	if bz == nil {
		return types.TradeDocument{}, false
	}

	var doc types.TradeDocument
	k.cdc.MustUnmarshal(bz, &doc)
	return doc, true
}

// SetTradeDocument saves a document
func (k Keeper) SetTradeDocument(ctx sdk.Context, doc types.TradeDocument) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.TradeDocumentPrefix)
	bz := k.cdc.MustMarshal(&doc)
	store.Set([]byte(doc.DocumentId), bz)
}

// GetDocumentsByLc returns all documents for an LC
func (k Keeper) GetDocumentsByLc(ctx sdk.Context, lcID string) []types.TradeDocument {
	docIDs := k.GetDocumentIDsByLc(ctx, lcID)
	
	var documents []types.TradeDocument
	for _, docID := range docIDs {
		doc, found := k.GetTradeDocument(ctx, docID)
		if found {
			documents = append(documents, doc)
		}
	}
	
	return documents
}

// AddDocumentToLcIndex adds a document to LC's index
func (k Keeper) AddDocumentToLcIndex(ctx sdk.Context, lcID, documentID string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DocumentByLcPrefix)
	key := append([]byte(lcID), []byte(documentID)...)
	store.Set(key, []byte{1})
}

// GetDocumentIDsByLc returns document IDs for an LC
func (k Keeper) GetDocumentIDsByLc(ctx sdk.Context, lcID string) []string {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DocumentByLcPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte(lcID))
	defer iterator.Close()
	
	var docIDs []string
	for ; iterator.Valid(); iterator.Next() {
		// Extract document ID from key
		key := iterator.Key()
		docID := string(key[len(lcID):])
		docIDs = append(docIDs, docID)
	}
	
	return docIDs
}

// CheckAllDocumentsVerified checks if all required documents are verified
func (k Keeper) CheckAllDocumentsVerified(ctx sdk.Context, lcID string, requiredDocTypes []string) bool {
	documents := k.GetDocumentsByLc(ctx, lcID)
	
	// Create a map of required document types
	requiredMap := make(map[string]bool)
	for _, docType := range requiredDocTypes {
		requiredMap[docType] = false
	}
	
	// Check which documents are verified
	for _, doc := range documents {
		if doc.Status == "verified" {
			if _, required := requiredMap[doc.DocumentType]; required {
				requiredMap[doc.DocumentType] = true
			}
		}
	}
	
	// Check if all required documents are verified
	for _, verified := range requiredMap {
		if !verified {
			return false
		}
	}
	
	return true
}

// GetNextDocumentID returns the next document ID
func (k Keeper) GetNextDocumentID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.NextDocumentIDKey)
	
	if bz == nil {
		return 1
	}
	
	return sdk.BigEndianToUint64(bz)
}

// SetNextDocumentID sets the next document ID
func (k Keeper) SetNextDocumentID(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.NextDocumentIDKey, sdk.Uint64ToBigEndian(id))
}