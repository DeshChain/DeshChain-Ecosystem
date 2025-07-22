package keeper

import (
	"context"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/deshchain/x/tradefinance/types"
)

// UCP600MessageServer implements UCP 600 specific message handling
type UCP600MessageServer struct {
	keeper *Keeper
	engine *UCP600ComplianceEngine
}

// NewUCP600MessageServer creates a new UCP 600 message server
func NewUCP600MessageServer(k *Keeper) *UCP600MessageServer {
	return &UCP600MessageServer{
		keeper: k,
		engine: NewUCP600ComplianceEngine(k),
	}
}

// SubmitDocumentPresentation handles document presentation submissions
func (s *UCP600MessageServer) SubmitDocumentPresentation(
	goCtx context.Context,
	req *types.MsgSubmitDocumentPresentation,
) (*types.MsgSubmitDocumentPresentationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate the request
	if err := req.ValidateBasic(); err != nil {
		return nil, err
	}

	// Check if LC exists and is in valid status
	lc, found := s.keeper.GetLetterOfCredit(ctx, req.LcId)
	if !found {
		return nil, types.ErrLCNotFound
	}

	// Validate LC status allows document presentation
	validStatuses := []string{"issued", "accepted", "confirmed"}
	if !s.isValidStatus(lc.Status, validStatuses) {
		return nil, fmt.Errorf("LC status %s does not allow document presentation", lc.Status)
	}

	// Check presentation deadline
	if s.isPresentationPastDeadline(ctx, lc) {
		return nil, types.ErrPresentationDeadlinePassed
	}

	// Validate presenter authority
	if err := s.validatePresenterAuthority(ctx, req.Presentor, lc); err != nil {
		return nil, err
	}

	// Generate presentation ID
	presentationID := s.generatePresentationID(ctx)

	// Create document presentation
	presentation := DocumentPresentation{
		ID:               presentationID,
		LcID:             req.LcId,
		PresentorID:      req.Presentor,
		PresentingBank:   req.PresentingBank,
		Documents:        req.Documents,
		PresentationDate: ctx.BlockTime(),
		Status:           "pending_examination",
		CreatedAt:        ctx.BlockTime(),
		UpdatedAt:        ctx.BlockTime(),
	}

	// Perform initial document structure validation
	validationErrors := s.validateDocumentStructure(presentation.Documents)
	if len(validationErrors) > 0 {
		presentation.Status = "rejected"
		// Store rejection reasons
	}

	// Store the presentation
	s.keeper.SetDocumentPresentation(ctx, presentation)

	// Calculate and collect presentation fees
	presentationFee := s.keeper.CalculateUCP600ExaminationFee(ctx, lc.Amount)
	presenterAddr, err := sdk.AccAddressFromBech32(req.Presentor)
	if err != nil {
		return nil, err
	}

	if err := s.keeper.bankKeeper.SendCoinsFromAccountToModule(
		ctx, presenterAddr, types.ModuleName, sdk.NewCoins(presentationFee),
	); err != nil {
		return nil, fmt.Errorf("failed to collect presentation fee: %w", err)
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"document_presentation_submitted",
			sdk.NewAttribute("lc_id", req.LcId),
			sdk.NewAttribute("presentation_id", presentationID),
			sdk.NewAttribute("presentor", req.Presentor),
			sdk.NewAttribute("document_count", fmt.Sprintf("%d", len(req.Documents))),
		),
	)

	return &types.MsgSubmitDocumentPresentationResponse{
		PresentationId:   presentationID,
		ExaminationFee:   presentationFee,
		ExaminationStart: ctx.BlockTime(),
		ExaminationEnd:   ctx.BlockTime().Add(5 * 24 * time.Hour), // 5 banking days
	}, nil
}

// ExamineDocuments handles document examination requests
func (s *UCP600MessageServer) ExamineDocuments(
	goCtx context.Context,
	req *types.MsgExamineDocuments,
) (*types.MsgExamineDocumentsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate examiner authority
	if err := s.validateExaminerAuthority(ctx, req.Examiner, req.BankId); err != nil {
		return nil, err
	}

	// Get presentation
	presentation, found := s.keeper.GetDocumentPresentation(ctx, req.PresentationId)
	if !found {
		return nil, types.ErrPresentationNotFound
	}

	// Validate presentation status
	if presentation.Status != "pending_examination" {
		return nil, fmt.Errorf("presentation status %s does not allow examination", presentation.Status)
	}

	// Check examination deadline (5 banking days from presentation)
	examinationDeadline := presentation.PresentationDate.Add(5 * 24 * time.Hour)
	if ctx.BlockTime().After(examinationDeadline) {
		return nil, fmt.Errorf("examination deadline has passed")
	}

	// Perform UCP 600 compliance check
	complianceResult, err := s.engine.PerformUCP600Compliance(ctx, presentation.LcID, req.PresentationId)
	if err != nil {
		return nil, fmt.Errorf("compliance check failed: %w", err)
	}

	// Update presentation with examination results
	presentation.Status = complianceResult.OverallStatus
	presentation.IsCompliant = complianceResult.IsCompliant
	presentation.ComplianceResult = complianceResult
	presentation.ExaminerID = req.Examiner
	examinationTime := ctx.BlockTime()
	presentation.ExaminationDate = &examinationTime
	presentation.UpdatedAt = ctx.BlockTime()

	s.keeper.SetDocumentPresentation(ctx, presentation)

	// Store compliance record
	complianceRecord := UCP600ComplianceRecord{
		ID:               s.generateComplianceRecordID(ctx),
		LcID:             presentation.LcID,
		PresentationID:   req.PresentationId,
		ExaminationDate:  ctx.BlockTime(),
		ExaminerID:       req.Examiner,
		IsCompliant:      complianceResult.IsCompliant,
		ComplianceScore:  complianceResult.ComplianceScore,
		DiscrepancyCount: len(complianceResult.Discrepancies),
		ProcessingTime:   complianceResult.ProcessingTime.Milliseconds(),
		Status:           complianceResult.OverallStatus,
		UCP600Version:    "UCP 600 2007",
		CreatedAt:        ctx.BlockTime(),
		UpdatedAt:        ctx.BlockTime(),
	}

	s.keeper.SetUCP600ComplianceRecord(ctx, complianceRecord)

	// Process the compliance result
	if err := s.keeper.ProcessUCP600ComplianceResult(
		ctx, presentation.LcID, req.PresentationId, complianceResult,
	); err != nil {
		return nil, err
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"documents_examined",
			sdk.NewAttribute("lc_id", presentation.LcID),
			sdk.NewAttribute("presentation_id", req.PresentationId),
			sdk.NewAttribute("examiner", req.Examiner),
			sdk.NewAttribute("compliance_status", complianceResult.OverallStatus),
			sdk.NewAttribute("compliance_score", fmt.Sprintf("%d", complianceResult.ComplianceScore)),
			sdk.NewAttribute("discrepancy_count", fmt.Sprintf("%d", len(complianceResult.Discrepancies))),
		),
	)

	return &types.MsgExamineDocumentsResponse{
		ComplianceResult: complianceResult,
		RecordId:        complianceRecord.ID,
	}, nil
}

// AcceptDiscrepancy handles discrepancy acceptance by applicant
func (s *UCP600MessageServer) AcceptDiscrepancy(
	goCtx context.Context,
	req *types.MsgAcceptDiscrepancy,
) (*types.MsgAcceptDiscrepancyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get LC
	lc, found := s.keeper.GetLetterOfCredit(ctx, req.LcId)
	if !found {
		return nil, types.ErrLCNotFound
	}

	// Validate applicant authority
	if req.Applicant != lc.ApplicantId {
		return nil, types.ErrUnauthorized
	}

	// Get presentation
	presentation, found := s.keeper.GetDocumentPresentation(ctx, req.PresentationId)
	if !found {
		return nil, types.ErrPresentationNotFound
	}

	// Validate presentation status
	if presentation.Status != "discrepant" {
		return nil, fmt.Errorf("only discrepant presentations can have discrepancies accepted")
	}

	// Update LC status to indicate discrepancy acceptance
	lc.Status = "discrepancy_accepted"
	lc.UpdatedAt = ctx.BlockTime()
	s.keeper.SetLetterOfCredit(ctx, lc)

	// Update presentation status
	presentation.Status = "accepted_with_discrepancies"
	presentation.UpdatedAt = ctx.BlockTime()
	s.keeper.SetDocumentPresentation(ctx, presentation)

	// Process payment if auto-payment is enabled
	if req.AuthorizePayment {
		if err := s.processPayment(ctx, lc, presentation); err != nil {
			return nil, fmt.Errorf("payment processing failed: %w", err)
		}
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"discrepancy_accepted",
			sdk.NewAttribute("lc_id", req.LcId),
			sdk.NewAttribute("presentation_id", req.PresentationId),
			sdk.NewAttribute("applicant", req.Applicant),
			sdk.NewAttribute("authorize_payment", fmt.Sprintf("%t", req.AuthorizePayment)),
		),
	)

	return &types.MsgAcceptDiscrepancyResponse{
		Success:       true,
		PaymentStatus: s.getPaymentStatus(ctx, lc.LcId),
	}, nil
}

// RejectDocuments handles document rejection
func (s *UCP600MessageServer) RejectDocuments(
	goCtx context.Context,
	req *types.MsgRejectDocuments,
) (*types.MsgRejectDocumentsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate rejecting bank authority
	if err := s.validateBankAuthority(ctx, req.RejectingBank, req.BankId); err != nil {
		return nil, err
	}

	// Get presentation
	presentation, found := s.keeper.GetDocumentPresentation(ctx, req.PresentationId)
	if !found {
		return nil, types.ErrPresentationNotFound
	}

	// Validate presentation status
	validStatuses := []string{"discrepant", "under_review"}
	if !s.isValidStatus(presentation.Status, validStatuses) {
		return nil, fmt.Errorf("presentation status %s does not allow rejection", presentation.Status)
	}

	// Update presentation status
	presentation.Status = "rejected"
	presentation.UpdatedAt = ctx.BlockTime()
	s.keeper.SetDocumentPresentation(ctx, presentation)

	// Update LC status
	lc, found := s.keeper.GetLetterOfCredit(ctx, presentation.LcID)
	if found {
		lc.Status = "documents_rejected"
		lc.UpdatedAt = ctx.BlockTime()
		s.keeper.SetLetterOfCredit(ctx, lc)
	}

	// Process refund if applicable
	if req.ProcessRefund {
		if err := s.processRefund(ctx, presentation); err != nil {
			return nil, fmt.Errorf("refund processing failed: %w", err)
		}
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"documents_rejected",
			sdk.NewAttribute("lc_id", presentation.LcID),
			sdk.NewAttribute("presentation_id", req.PresentationId),
			sdk.NewAttribute("rejecting_bank", req.RejectingBank),
			sdk.NewAttribute("rejection_reason", req.RejectionReason),
		),
	)

	return &types.MsgRejectDocumentsResponse{
		Success:      true,
		RefundStatus: "pending", // Would be determined by refund processing
	}, nil
}

// Helper methods

func (s *UCP600MessageServer) isValidStatus(status string, validStatuses []string) bool {
	for _, valid := range validStatuses {
		if status == valid {
			return true
		}
	}
	return false
}

func (s *UCP600MessageServer) isPresentationPastDeadline(ctx sdk.Context, lc types.LetterOfCredit) bool {
	presentationDeadline := lc.ExpiryDate
	if !lc.LatestShipmentDate.IsZero() {
		shipmentDeadline := lc.LatestShipmentDate.AddDate(0, 0, 21) // 21 days after shipment
		if shipmentDeadline.Before(presentationDeadline) {
			presentationDeadline = shipmentDeadline
		}
	}
	return ctx.BlockTime().After(presentationDeadline)
}

func (s *UCP600MessageServer) validatePresenterAuthority(
	ctx sdk.Context,
	presenter string,
	lc types.LetterOfCredit,
) error {
	// Presenter must be beneficiary or authorized party
	presenterPartyID := s.keeper.GetPartyIDByAddress(ctx, presenter)
	if presenterPartyID != lc.BeneficiaryId {
		// Check if presenter is authorized by beneficiary
		// For now, only allow direct beneficiary presentation
		return fmt.Errorf("only beneficiary can present documents")
	}
	return nil
}

func (s *UCP600MessageServer) validateExaminerAuthority(
	ctx sdk.Context,
	examiner string,
	bankID string,
) error {
	// Examiner must be authorized by the examining bank
	examinerPartyID := s.keeper.GetPartyIDByAddress(ctx, examiner)
	if examinerPartyID != bankID {
		return fmt.Errorf("examiner not authorized by bank")
	}

	// Bank must be a registered trade finance participant
	bank, found := s.keeper.GetTradeParty(ctx, bankID)
	if !found || bank.PartyType != "bank" {
		return fmt.Errorf("invalid examining bank")
	}

	return nil
}

func (s *UCP600MessageServer) validateBankAuthority(
	ctx sdk.Context,
	banker string,
	bankID string,
) error {
	// Similar to examiner validation
	return s.validateExaminerAuthority(ctx, banker, bankID)
}

func (s *UCP600MessageServer) validateDocumentStructure(docs []types.TradeDocument) []string {
	var errors []string
	
	for _, doc := range docs {
		docErrors := s.keeper.ValidateUCP600DocumentStructure(doc)
		errors = append(errors, docErrors...)
	}
	
	return errors
}

func (s *UCP600MessageServer) generatePresentationID(ctx sdk.Context) string {
	// Generate unique presentation ID
	blockHeight := ctx.BlockHeight()
	timestamp := ctx.BlockTime().Unix()
	return fmt.Sprintf("PRES-%d-%d", blockHeight, timestamp)
}

func (s *UCP600MessageServer) generateComplianceRecordID(ctx sdk.Context) string {
	// Generate unique compliance record ID
	blockHeight := ctx.BlockHeight()
	timestamp := ctx.BlockTime().Unix()
	return fmt.Sprintf("UCP600-%d-%d", blockHeight, timestamp)
}

func (s *UCP600MessageServer) processPayment(
	ctx sdk.Context,
	lc types.LetterOfCredit,
	presentation DocumentPresentation,
) error {
	// Process LC payment to beneficiary
	// This would involve complex payment processing logic
	// For now, implement basic escrow release

	beneficiaryAddr, err := sdk.AccAddressFromBech32(presentation.PresentorID)
	if err != nil {
		return err
	}

	// Release LC amount to beneficiary
	if err := s.keeper.bankKeeper.SendCoinsFromModuleToAccount(
		ctx, types.ModuleName, beneficiaryAddr, sdk.NewCoins(lc.Amount),
	); err != nil {
		return err
	}

	// Update LC status
	lc.Status = "payment_processed"
	lc.UpdatedAt = ctx.BlockTime()
	s.keeper.SetLetterOfCredit(ctx, lc)

	return nil
}

func (s *UCP600MessageServer) processRefund(
	ctx sdk.Context,
	presentation DocumentPresentation,
) error {
	// Process refund logic
	// This would typically involve refunding fees to the presenter
	presenterAddr, err := sdk.AccAddressFromBech32(presentation.PresentorID)
	if err != nil {
		return err
	}

	// Calculate refund amount (e.g., partial fee refund)
	lc, found := s.keeper.GetLetterOfCredit(ctx, presentation.LcID)
	if !found {
		return types.ErrLCNotFound
	}

	refundAmount := s.keeper.CalculateUCP600ExaminationFee(ctx, lc.Amount)
	refundAmount.Amount = refundAmount.Amount.Quo(sdk.NewInt(2)) // 50% refund

	return s.keeper.bankKeeper.SendCoinsFromModuleToAccount(
		ctx, types.ModuleName, presenterAddr, sdk.NewCoins(refundAmount),
	)
}

func (s *UCP600MessageServer) getPaymentStatus(ctx sdk.Context, lcID string) string {
	lc, found := s.keeper.GetLetterOfCredit(ctx, lcID)
	if !found {
		return "unknown"
	}

	switch lc.Status {
	case "payment_processed":
		return "completed"
	case "discrepancy_accepted":
		return "authorized"
	default:
		return "pending"
	}
}