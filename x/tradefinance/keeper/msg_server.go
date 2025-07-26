package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/tradefinance/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the trade finance MsgServer interface
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// RegisterParty implements types.MsgServer
func (k msgServer) RegisterParty(goCtx context.Context, msg *types.MsgRegisterParty) (*types.MsgRegisterPartyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Create trade party
	party := types.TradeParty{
		PartyType:   msg.PartyType,
		Name:        msg.Name,
		Address:     msg.Address,
		DeshAddress: msg.Creator,
		Country:     msg.Country,
		TaxId:       msg.TaxId,
		KycLevel:    msg.KycLevel,
	}

	// Register party
	partyID, err := k.Keeper.RegisterParty(ctx, party)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypePartyRegistered,
			sdk.NewAttribute(types.AttributeKeyPartyId, partyID),
			sdk.NewAttribute(types.AttributeKeyPartyType, msg.PartyType),
			sdk.NewAttribute(types.AttributeKeyName, msg.Name),
			sdk.NewAttribute(types.AttributeKeyDeshAddress, msg.Creator),
			sdk.NewAttribute(types.AttributeKeyCountry, msg.Country),
		),
	)

	return &types.MsgRegisterPartyResponse{
		PartyId: partyID,
	}, nil
}

// IssueLc implements types.MsgServer
func (k msgServer) IssueLc(goCtx context.Context, msg *types.MsgIssueLc) (*types.MsgIssueLcResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	lcID, lcNumber, err := k.Keeper.IssueLc(ctx, msg)
	if err != nil {
		return nil, err
	}

	return &types.MsgIssueLcResponse{
		LcId:     lcID,
		LcNumber: lcNumber,
	}, nil
}

// AcceptLc implements types.MsgServer
func (k msgServer) AcceptLc(goCtx context.Context, msg *types.MsgAcceptLc) (*types.MsgAcceptLcResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := k.Keeper.AcceptLc(ctx, msg.LcId, msg.Beneficiary)
	if err != nil {
		return nil, err
	}

	return &types.MsgAcceptLcResponse{}, nil
}

// SubmitDocuments implements types.MsgServer
func (k msgServer) SubmitDocuments(goCtx context.Context, msg *types.MsgSubmitDocuments) (*types.MsgSubmitDocumentsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	documentIDs, err := k.Keeper.SubmitDocuments(ctx, msg.LcId, msg.Submitter, msg.Documents)
	if err != nil {
		return nil, err
	}

	return &types.MsgSubmitDocumentsResponse{
		DocumentIds: documentIDs,
	}, nil
}

// VerifyDocument implements types.MsgServer
func (k msgServer) VerifyDocument(goCtx context.Context, msg *types.MsgVerifyDocument) (*types.MsgVerifyDocumentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := k.Keeper.VerifyDocument(ctx, msg.DocumentId, msg.Verifier, msg.Approved, msg.RejectionReason)
	if err != nil {
		return nil, err
	}

	return &types.MsgVerifyDocumentResponse{}, nil
}

// RequestPayment implements types.MsgServer
func (k msgServer) RequestPayment(goCtx context.Context, msg *types.MsgRequestPayment) (*types.MsgRequestPaymentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	instructionID, err := k.Keeper.RequestPayment(ctx, msg.LcId, msg.Beneficiary, msg.Amount)
	if err != nil {
		return nil, err
	}

	return &types.MsgRequestPaymentResponse{
		PaymentInstructionId: instructionID,
	}, nil
}

// MakePayment implements types.MsgServer
func (k msgServer) MakePayment(goCtx context.Context, msg *types.MsgMakePayment) (*types.MsgMakePaymentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := k.Keeper.MakePayment(ctx, msg.PaymentInstructionId, msg.Payer)
	if err != nil {
		return nil, err
	}

	// Get the updated payment instruction to return transaction hash
	instruction, found := k.Keeper.GetPaymentInstruction(ctx, msg.PaymentInstructionId)
	if !found {
		return nil, types.ErrPaymentInstructionNotFound
	}

	return &types.MsgMakePaymentResponse{
		TransactionHash: instruction.TransactionHash,
	}, nil
}

// AmendLc implements types.MsgServer
func (k msgServer) AmendLc(goCtx context.Context, msg *types.MsgAmendLc) (*types.MsgAmendLcResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get LC
	lc, found := k.Keeper.GetLetterOfCredit(ctx, msg.LcId)
	if !found {
		return nil, types.ErrLCNotFound
	}

	// Validate issuing bank
	issuingBankPartyID := k.Keeper.GetPartyIDByAddress(ctx, msg.IssuingBank)
	if issuingBankPartyID != lc.IssuingBankId {
		return nil, types.ErrUnauthorized
	}

	// Update LC with amendment (simplified)
	lc.UpdatedAt = ctx.BlockTime()
	k.Keeper.SetLetterOfCredit(ctx, lc)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeLcAmended,
			sdk.NewAttribute(types.AttributeKeyLcId, msg.LcId),
			sdk.NewAttribute(types.AttributeKeyAmendmentType, msg.AmendmentType),
			sdk.NewAttribute(types.AttributeKeyAmendedBy, issuingBankPartyID),
		),
	)

	return &types.MsgAmendLcResponse{}, nil
}

// CancelLc implements types.MsgServer
func (k msgServer) CancelLc(goCtx context.Context, msg *types.MsgCancelLc) (*types.MsgCancelLcResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get LC
	lc, found := k.Keeper.GetLetterOfCredit(ctx, msg.LcId)
	if !found {
		return nil, types.ErrLCNotFound
	}

	// Validate issuing bank
	issuingBankPartyID := k.Keeper.GetPartyIDByAddress(ctx, msg.IssuingBank)
	if issuingBankPartyID != lc.IssuingBankId {
		return nil, types.ErrUnauthorized
	}

	// Cancel LC
	lc.Status = "cancelled"
	lc.UpdatedAt = ctx.BlockTime()
	k.Keeper.SetLetterOfCredit(ctx, lc)

	// Release collateral
	issuingBank, found := k.Keeper.GetTradeParty(ctx, lc.IssuingBankId)
	if found {
		issuingBankAddr, err := sdk.AccAddressFromBech32(issuingBank.DeshAddress)
		if err == nil {
			k.Keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, issuingBankAddr, sdk.NewCoins(lc.Collateral))
		}
	}

	// Update stats
	stats := k.Keeper.GetTradeFinanceStats(ctx)
	stats.ActiveLcs--
	k.Keeper.SetTradeFinanceStats(ctx, stats)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeLcCancelled,
			sdk.NewAttribute(types.AttributeKeyLcId, msg.LcId),
			sdk.NewAttribute(types.AttributeKeyCancelledBy, issuingBankPartyID),
			sdk.NewAttribute(types.AttributeKeyReason, msg.Reason),
		),
	)

	return &types.MsgCancelLcResponse{}, nil
}

// CreateInsurancePolicy implements types.MsgServer
func (k msgServer) CreateInsurancePolicy(goCtx context.Context, msg *types.MsgCreateInsurancePolicy) (*types.MsgCreateInsurancePolicyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	policyID, err := k.Keeper.CreateInsurancePolicy(ctx, msg)
	if err != nil {
		return nil, err
	}

	return &types.MsgCreateInsurancePolicyResponse{
		PolicyId: policyID,
	}, nil
}

// UpdateShipment implements types.MsgServer
func (k msgServer) UpdateShipment(goCtx context.Context, msg *types.MsgUpdateShipment) (*types.MsgUpdateShipmentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := k.Keeper.UpdateShipment(ctx, msg)
	if err != nil {
		return nil, err
	}

	return &types.MsgUpdateShipmentResponse{}, nil
}

// UpdateParams implements types.MsgServer
func (k msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Verify authority
	if msg.Authority != k.authority {
		return nil, types.ErrUnauthorized
	}

	if err := k.SetParams(ctx, msg.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}