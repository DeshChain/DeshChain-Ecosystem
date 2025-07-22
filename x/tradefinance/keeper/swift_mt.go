package keeper

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/deshchain/x/tradefinance/types"
)

// SWIFTMTProcessor handles SWIFT MT message processing for trade finance
type SWIFTMTProcessor struct {
	keeper *Keeper
}

// NewSWIFTMTProcessor creates a new SWIFT MT message processor
func NewSWIFTMTProcessor(k *Keeper) *SWIFTMTProcessor {
	return &SWIFTMTProcessor{
		keeper: k,
	}
}

// MTMessage represents a generic SWIFT MT message structure
type MTMessage struct {
	MessageType        string            `json:"message_type"`        // MT700, MT710, etc.
	MessageReference   string            `json:"message_reference"`   // Transaction reference
	RelatedReference   string            `json:"related_reference"`   // Related message reference
	DateTimeStamp      time.Time         `json:"datetime_stamp"`      // Message timestamp
	SenderBIC          string            `json:"sender_bic"`          // Sender's BIC code
	ReceiverBIC        string            `json:"receiver_bic"`        // Receiver's BIC code
	Priority           string            `json:"priority"`            // Normal, Urgent
	Fields             map[string]string `json:"fields"`              // MT field values
	BlockchainTxHash   string            `json:"blockchain_tx_hash"`  // Associated blockchain transaction
	ProcessingStatus   string            `json:"processing_status"`   // pending, processed, failed
	ValidationErrors   []string          `json:"validation_errors"`
	CreatedAt          time.Time         `json:"created_at"`
	ProcessedAt        *time.Time        `json:"processed_at,omitempty"`
}

// MT700Message represents an Issue of Documentary Credit message
type MT700Message struct {
	MTMessage
	// MT700 specific fields
	FormOfDocumentaryCredit string    `json:"form_of_documentary_credit"` // Field 40A
	ApplicableRules         string    `json:"applicable_rules"`           // Field 40E (UCP 600)
	DateOfIssue            time.Time `json:"date_of_issue"`              // Field 31C
	DateOfExpiry           time.Time `json:"date_of_expiry"`             // Field 31D
	PlaceOfExpiry          string    `json:"place_of_expiry"`            // Field 31D
	ApplicantBank          string    `json:"applicant_bank"`             // Field 50
	Applicant              string    `json:"applicant"`                   // Field 50
	Beneficiary            string    `json:"beneficiary"`                 // Field 59
	CurrencyCode           string    `json:"currency_code"`               // Field 32B
	Amount                 string    `json:"amount"`                      // Field 32B
	AvailableWithBy        string    `json:"available_with_by"`           // Field 41a
	DraftsAt               string    `json:"drafts_at"`                   // Field 42C
	DraweeBank             string    `json:"drawee_bank"`                 // Field 42a
	PartialShipments       string    `json:"partial_shipments"`          // Field 43P
	Transhipment           string    `json:"transhipment"`                // Field 43T
	ShipmentPeriod         string    `json:"shipment_period"`             // Field 44A/B/C
	GoodsServicesDescription string  `json:"goods_services_description"`  // Field 45A
	DocumentsRequired      string    `json:"documents_required"`          // Field 46A
	AdditionalConditions   string    `json:"additional_conditions"`       // Field 47A
	PresentationPeriod     string    `json:"presentation_period"`         // Field 48
	ConfirmationInstructions string  `json:"confirmation_instructions"`   // Field 49
	AdvisingBank           string    `json:"advising_bank"`               // Field 57a
	SenderToReceiverInfo   string    `json:"sender_to_receiver_info"`     // Field 72Z
}

// MT710Message represents Advice of Third Bank's Documentary Credit
type MT710Message struct {
	MTMessage
	// MT710 specific fields
	AdvisingBankReference string    `json:"advising_bank_reference"`    // Field 21
	IssuingBankReference  string    `json:"issuing_bank_reference"`     // Field 20
	DateOfAdvice          time.Time `json:"date_of_advice"`             // Field 31C
	IssuingBank           string    `json:"issuing_bank"`               // Field 51A
	OriginalMT700         MT700Message `json:"original_mt700"`          // Embedded MT700 data
}

// MT720Message represents Transfer of Documentary Credit
type MT720Message struct {
	MTMessage
	// MT720 specific fields
	TransferringBankRef string    `json:"transferring_bank_ref"`      // Field 20
	OriginalCreditRef   string    `json:"original_credit_ref"`        // Field 21
	TransferDate        time.Time `json:"transfer_date"`              // Field 31C
	TransferringBank    string    `json:"transferring_bank"`          // Field 51A
	SecondBeneficiary   string    `json:"second_beneficiary"`         // Field 59
	NewApplicant        string    `json:"new_applicant"`              // Field 50 (if changed)
	TransferConditions  string    `json:"transfer_conditions"`        // Field 47A
	AmountTransferred   string    `json:"amount_transferred"`         // Field 32B
	PercentageTransferred float64 `json:"percentage_transferred"`     // Calculated field
}

// MT750Message represents Advice of Discrepancy
type MT750Message struct {
	MTMessage
	// MT750 specific fields
	DocumentaryCreditsRef string      `json:"documentary_credits_ref"`    // Field 20
	DateOfDiscrepancy    time.Time   `json:"date_of_discrepancy"`        // Field 31C
	DiscrepancyDetails   string      `json:"discrepancy_details"`        // Field 77A
	ActionTaken          string      `json:"action_taken"`               // Field 78
	DocumentsDisposition string      `json:"documents_disposition"`      // Field 79
	BankToBank           string      `json:"bank_to_bank"`               // Field 72Z
	DiscrepancyType      []string    `json:"discrepancy_type"`           // Parsed from details
	Severity             string      `json:"severity"`                   // major, minor, critical
}

// ProcessMT700 processes an MT700 (Issue of Documentary Credit) message
func (smt *SWIFTMTProcessor) ProcessMT700(
	ctx sdk.Context,
	messageData string,
	senderBIC string,
	receiverBIC string,
) (*MT700Message, error) {
	// Parse the MT700 message
	mt700, err := smt.parseMT700(messageData, senderBIC, receiverBIC)
	if err != nil {
		return nil, fmt.Errorf("failed to parse MT700: %w", err)
	}

	// Validate the message
	if err := smt.validateMT700(ctx, mt700); err != nil {
		mt700.ValidationErrors = append(mt700.ValidationErrors, err.Error())
		mt700.ProcessingStatus = "validation_failed"
		return mt700, err
	}

	// Create corresponding LC in the blockchain
	lcID, err := smt.createLCFromMT700(ctx, mt700)
	if err != nil {
		mt700.ProcessingStatus = "failed"
		return mt700, fmt.Errorf("failed to create LC: %w", err)
	}

	// Store MT700 message
	if err := smt.storeMTMessage(ctx, mt700.MTMessage); err != nil {
		smt.keeper.Logger(ctx).Error("Failed to store MT700 message", "error", err)
	}

	// Mark as processed
	mt700.ProcessingStatus = "processed"
	processedTime := ctx.BlockTime()
	mt700.ProcessedAt = &processedTime

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"swift_mt700_processed",
			sdk.NewAttribute("message_ref", mt700.MessageReference),
			sdk.NewAttribute("lc_id", lcID),
			sdk.NewAttribute("sender_bic", mt700.SenderBIC),
			sdk.NewAttribute("receiver_bic", mt700.ReceiverBIC),
			sdk.NewAttribute("amount", fmt.Sprintf("%s %s", mt700.CurrencyCode, mt700.Amount)),
		),
	)

	return mt700, nil
}

// ProcessMT710 processes an MT710 (Advice of Third Bank's Documentary Credit) message
func (smt *SWIFTMTProcessor) ProcessMT710(
	ctx sdk.Context,
	messageData string,
	senderBIC string,
	receiverBIC string,
) (*MT710Message, error) {
	// Parse the MT710 message
	mt710, err := smt.parseMT710(messageData, senderBIC, receiverBIC)
	if err != nil {
		return nil, fmt.Errorf("failed to parse MT710: %w", err)
	}

	// Validate the message
	if err := smt.validateMT710(ctx, mt710); err != nil {
		mt710.ValidationErrors = append(mt710.ValidationErrors, err.Error())
		mt710.ProcessingStatus = "validation_failed"
		return mt710, err
	}

	// Find related LC
	lc, err := smt.findLCByReference(ctx, mt710.IssuingBankReference)
	if err != nil {
		return mt710, fmt.Errorf("failed to find related LC: %w", err)
	}

	// Update LC with advising bank information
	if err := smt.updateLCWithAdvisingBank(ctx, lc.LcId, senderBIC, mt710.DateOfAdvice); err != nil {
		return mt710, fmt.Errorf("failed to update LC: %w", err)
	}

	// Store MT710 message
	if err := smt.storeMTMessage(ctx, mt710.MTMessage); err != nil {
		smt.keeper.Logger(ctx).Error("Failed to store MT710 message", "error", err)
	}

	mt710.ProcessingStatus = "processed"
	processedTime := ctx.BlockTime()
	mt710.ProcessedAt = &processedTime

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"swift_mt710_processed",
			sdk.NewAttribute("message_ref", mt710.MessageReference),
			sdk.NewAttribute("related_lc_ref", mt710.IssuingBankReference),
			sdk.NewAttribute("advising_bank", senderBIC),
		),
	)

	return mt710, nil
}

// ProcessMT720 processes an MT720 (Transfer of Documentary Credit) message
func (smt *SWIFTMTProcessor) ProcessMT720(
	ctx sdk.Context,
	messageData string,
	senderBIC string,
	receiverBIC string,
) (*MT720Message, error) {
	// Parse the MT720 message
	mt720, err := smt.parseMT720(messageData, senderBIC, receiverBIC)
	if err != nil {
		return nil, fmt.Errorf("failed to parse MT720: %w", err)
	}

	// Validate the message
	if err := smt.validateMT720(ctx, mt720); err != nil {
		mt720.ValidationErrors = append(mt720.ValidationErrors, err.Error())
		mt720.ProcessingStatus = "validation_failed"
		return mt720, err
	}

	// Find original LC
	originalLC, err := smt.findLCByReference(ctx, mt720.OriginalCreditRef)
	if err != nil {
		return mt720, fmt.Errorf("failed to find original LC: %w", err)
	}

	// Create transferred LC
	transferredLCID, err := smt.createTransferredLC(ctx, originalLC, mt720)
	if err != nil {
		return mt720, fmt.Errorf("failed to create transferred LC: %w", err)
	}

	// Store MT720 message
	if err := smt.storeMTMessage(ctx, mt720.MTMessage); err != nil {
		smt.keeper.Logger(ctx).Error("Failed to store MT720 message", "error", err)
	}

	mt720.ProcessingStatus = "processed"
	processedTime := ctx.BlockTime()
	mt720.ProcessedAt = &processedTime

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"swift_mt720_processed",
			sdk.NewAttribute("message_ref", mt720.MessageReference),
			sdk.NewAttribute("original_lc_ref", mt720.OriginalCreditRef),
			sdk.NewAttribute("transferred_lc_id", transferredLCID),
			sdk.NewAttribute("transferring_bank", senderBIC),
			sdk.NewAttribute("second_beneficiary", mt720.SecondBeneficiary),
		),
	)

	return mt720, nil
}

// ProcessMT750 processes an MT750 (Advice of Discrepancy) message
func (smt *SWIFTMTProcessor) ProcessMT750(
	ctx sdk.Context,
	messageData string,
	senderBIC string,
	receiverBIC string,
) (*MT750Message, error) {
	// Parse the MT750 message
	mt750, err := smt.parseMT750(messageData, senderBIC, receiverBIC)
	if err != nil {
		return nil, fmt.Errorf("failed to parse MT750: %w", err)
	}

	// Validate the message
	if err := smt.validateMT750(ctx, mt750); err != nil {
		mt750.ValidationErrors = append(mt750.ValidationErrors, err.Error())
		mt750.ProcessingStatus = "validation_failed"
		return mt750, err
	}

	// Find related LC
	lc, err := smt.findLCByReference(ctx, mt750.DocumentaryCreditsRef)
	if err != nil {
		return mt750, fmt.Errorf("failed to find related LC: %w", err)
	}

	// Update LC status with discrepancy information
	if err := smt.updateLCWithDiscrepancy(ctx, lc.LcId, mt750); err != nil {
		return mt750, fmt.Errorf("failed to update LC with discrepancy: %w", err)
	}

	// Store MT750 message
	if err := smt.storeMTMessage(ctx, mt750.MTMessage); err != nil {
		smt.keeper.Logger(ctx).Error("Failed to store MT750 message", "error", err)
	}

	mt750.ProcessingStatus = "processed"
	processedTime := ctx.BlockTime()
	mt750.ProcessedAt = &processedTime

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"swift_mt750_processed",
			sdk.NewAttribute("message_ref", mt750.MessageReference),
			sdk.NewAttribute("lc_ref", mt750.DocumentaryCreditsRef),
			sdk.NewAttribute("discrepancy_date", mt750.DateOfDiscrepancy.Format("2006-01-02")),
			sdk.NewAttribute("severity", mt750.Severity),
		),
	)

	return mt750, nil
}

// Message parsing functions

func (smt *SWIFTMTProcessor) parseMT700(messageData, senderBIC, receiverBIC string) (*MT700Message, error) {
	fields := smt.parseMessageFields(messageData)
	
	mt700 := &MT700Message{
		MTMessage: MTMessage{
			MessageType:      "MT700",
			MessageReference: fields["20"],
			RelatedReference: fields["21"],
			DateTimeStamp:    time.Now(),
			SenderBIC:        senderBIC,
			ReceiverBIC:      receiverBIC,
			Fields:           fields,
			ProcessingStatus: "pending",
			CreatedAt:        time.Now(),
		},
	}

	// Parse specific MT700 fields
	mt700.FormOfDocumentaryCredit = fields["40A"]
	mt700.ApplicableRules = fields["40E"]
	
	// Parse dates
	if dateStr := fields["31C"]; dateStr != "" {
		if date, err := smt.parseDate(dateStr); err == nil {
			mt700.DateOfIssue = date
		}
	}
	if dateStr := fields["31D"]; dateStr != "" {
		if date, err := smt.parseDate(dateStr); err == nil {
			mt700.DateOfExpiry = date
		}
		// Extract place from 31D field (format: YYMMDDCCCC where CCCC is place code)
		if len(dateStr) > 6 {
			mt700.PlaceOfExpiry = dateStr[6:]
		}
	}

	// Parse currency and amount from field 32B
	if amountStr := fields["32B"]; amountStr != "" {
		if len(amountStr) >= 3 {
			mt700.CurrencyCode = amountStr[:3]
			mt700.Amount = amountStr[3:]
		}
	}

	mt700.ApplicantBank = fields["50"]
	mt700.Applicant = fields["50"]
	mt700.Beneficiary = fields["59"]
	mt700.AvailableWithBy = fields["41A"]
	mt700.DraftsAt = fields["42C"]
	mt700.DraweeBank = fields["42A"]
	mt700.PartialShipments = fields["43P"]
	mt700.Transhipment = fields["43T"]
	mt700.ShipmentPeriod = fields["44A"]
	if mt700.ShipmentPeriod == "" {
		mt700.ShipmentPeriod = fields["44B"]
	}
	if mt700.ShipmentPeriod == "" {
		mt700.ShipmentPeriod = fields["44C"]
	}
	mt700.GoodsServicesDescription = fields["45A"]
	mt700.DocumentsRequired = fields["46A"]
	mt700.AdditionalConditions = fields["47A"]
	mt700.PresentationPeriod = fields["48"]
	mt700.ConfirmationInstructions = fields["49"]
	mt700.AdvisingBank = fields["57A"]
	mt700.SenderToReceiverInfo = fields["72Z"]

	return mt700, nil
}

func (smt *SWIFTMTProcessor) parseMT710(messageData, senderBIC, receiverBIC string) (*MT710Message, error) {
	fields := smt.parseMessageFields(messageData)
	
	mt710 := &MT710Message{
		MTMessage: MTMessage{
			MessageType:      "MT710",
			MessageReference: fields["21"], // Advising bank's reference
			RelatedReference: fields["20"], // Issuing bank's reference
			DateTimeStamp:    time.Now(),
			SenderBIC:        senderBIC,
			ReceiverBIC:      receiverBIC,
			Fields:           fields,
			ProcessingStatus: "pending",
			CreatedAt:        time.Now(),
		},
		AdvisingBankReference: fields["21"],
		IssuingBankReference:  fields["20"],
		IssuingBank:          fields["51A"],
	}

	// Parse advice date
	if dateStr := fields["31C"]; dateStr != "" {
		if date, err := smt.parseDate(dateStr); err == nil {
			mt710.DateOfAdvice = date
		}
	}

	return mt710, nil
}

func (smt *SWIFTMTProcessor) parseMT720(messageData, senderBIC, receiverBIC string) (*MT720Message, error) {
	fields := smt.parseMessageFields(messageData)
	
	mt720 := &MT720Message{
		MTMessage: MTMessage{
			MessageType:      "MT720",
			MessageReference: fields["20"],
			RelatedReference: fields["21"],
			DateTimeStamp:    time.Now(),
			SenderBIC:        senderBIC,
			ReceiverBIC:      receiverBIC,
			Fields:           fields,
			ProcessingStatus: "pending",
			CreatedAt:        time.Now(),
		},
		TransferringBankRef: fields["20"],
		OriginalCreditRef:   fields["21"],
		TransferringBank:    senderBIC,
		SecondBeneficiary:   fields["59"],
		NewApplicant:        fields["50"],
		TransferConditions:  fields["47A"],
		AmountTransferred:   fields["32B"],
	}

	// Parse transfer date
	if dateStr := fields["31C"]; dateStr != "" {
		if date, err := smt.parseDate(dateStr); err == nil {
			mt720.TransferDate = date
		}
	}

	// Calculate percentage transferred (if needed)
	if fields["39A"] != "" {
		if percentage, err := strconv.ParseFloat(fields["39A"], 64); err == nil {
			mt720.PercentageTransferred = percentage
		}
	}

	return mt720, nil
}

func (smt *SWIFTMTProcessor) parseMT750(messageData, senderBIC, receiverBIC string) (*MT750Message, error) {
	fields := smt.parseMessageFields(messageData)
	
	mt750 := &MT750Message{
		MTMessage: MTMessage{
			MessageType:      "MT750",
			MessageReference: fields["20"],
			DateTimeStamp:    time.Now(),
			SenderBIC:        senderBIC,
			ReceiverBIC:      receiverBIC,
			Fields:           fields,
			ProcessingStatus: "pending",
			CreatedAt:        time.Now(),
		},
		DocumentaryCreditsRef: fields["20"],
		DiscrepancyDetails:    fields["77A"],
		ActionTaken:          fields["78"],
		DocumentsDisposition: fields["79"],
		BankToBank:           fields["72Z"],
	}

	// Parse discrepancy date
	if dateStr := fields["31C"]; dateStr != "" {
		if date, err := smt.parseDate(dateStr); err == nil {
			mt750.DateOfDiscrepancy = date
		}
	}

	// Analyze discrepancy severity from details
	mt750.Severity = smt.analyzeDiscrepancySeverity(mt750.DiscrepancyDetails)
	mt750.DiscrepancyType = smt.extractDiscrepancyTypes(mt750.DiscrepancyDetails)

	return mt750, nil
}

// Helper functions

func (smt *SWIFTMTProcessor) parseMessageFields(messageData string) map[string]string {
	fields := make(map[string]string)
	
	// Simple field parser - in reality, this would be more sophisticated
	// SWIFT messages use field tags like :20:, :21:, etc.
	re := regexp.MustCompile(`:(\d+\w*):([^:]+)`)
	matches := re.FindAllStringSubmatch(messageData, -1)
	
	for _, match := range matches {
		if len(match) >= 3 {
			fieldTag := match[1]
			fieldValue := strings.TrimSpace(match[2])
			fields[fieldTag] = fieldValue
		}
	}
	
	return fields
}

func (smt *SWIFTMTProcessor) parseDate(dateStr string) (time.Time, error) {
	// SWIFT date format is usually YYMMDD
	if len(dateStr) >= 6 {
		return time.Parse("060102", dateStr[:6])
	}
	return time.Time{}, fmt.Errorf("invalid date format: %s", dateStr)
}

func (smt *SWIFTMTProcessor) analyzeDiscrepancySeverity(details string) string {
	details = strings.ToLower(details)
	
	criticalKeywords := []string{"reject", "refus", "not accept", "critical", "major discrepancy"}
	majorKeywords := []string{"discrepancy", "inconsisten", "missing", "incorrect"}
	
	for _, keyword := range criticalKeywords {
		if strings.Contains(details, keyword) {
			return "critical"
		}
	}
	
	for _, keyword := range majorKeywords {
		if strings.Contains(details, keyword) {
			return "major"
		}
	}
	
	return "minor"
}

func (smt *SWIFTMTProcessor) extractDiscrepancyTypes(details string) []string {
	var types []string
	details = strings.ToLower(details)
	
	typeMap := map[string]string{
		"amount":         "amount_discrepancy",
		"date":           "date_discrepancy", 
		"document":       "document_discrepancy",
		"signature":      "signature_missing",
		"endorsement":    "endorsement_missing",
		"description":    "description_mismatch",
		"transport":      "transport_inconsistent",
		"insurance":      "insurance_insufficient",
		"late present":   "presentation_late",
		"expired":        "expiry_passed",
	}
	
	for keyword, discType := range typeMap {
		if strings.Contains(details, keyword) {
			types = append(types, discType)
		}
	}
	
	if len(types) == 0 {
		types = append(types, "general_discrepancy")
	}
	
	return types
}

// Blockchain integration functions

func (smt *SWIFTMTProcessor) createLCFromMT700(ctx sdk.Context, mt700 *MT700Message) (string, error) {
	// Convert SWIFT BIC to party IDs (this would require BIC-to-party mapping)
	issuingBankID := smt.bicToPartyID(ctx, mt700.SenderBIC)
	advisingBankID := smt.bicToPartyID(ctx, mt700.ReceiverBIC)
	
	// Parse amount
	amountInt, err := smt.parseAmountToSDK(mt700.Amount, mt700.CurrencyCode)
	if err != nil {
		return "", err
	}

	// Create LC using existing keeper method
	// This would need to be adapted to use the existing IssueLc method
	params := smt.keeper.GetParams(ctx)
	
	// Create a minimal collateral (this would be handled differently in production)
	collateral := sdk.NewCoin("dinr", amountInt.Amount.Quo(sdk.NewInt(10))) // 10% collateral
	
	// Create LC message structure
	msg := &types.MsgIssueLc{
		Creator:              issuingBankID,
		IssuingBank:          issuingBankID,
		ApplicantId:          smt.parsePartyFromField(mt700.Applicant),
		BeneficiaryId:        smt.parsePartyFromField(mt700.Beneficiary),
		AdvisingBankId:       advisingBankID,
		Amount:               amountInt,
		ExpiryDate:           mt700.DateOfExpiry,
		LatestShipmentDate:   smt.parseShipmentDate(mt700.ShipmentPeriod),
		PaymentTerms:         mt700.DraftsAt,
		DeferredPaymentDays:  0, // Parse from payment terms
		Incoterms:           smt.extractIncoterms(mt700.ShipmentPeriod),
		PortOfLoading:       smt.extractPortOfLoading(mt700.ShipmentPeriod),
		PortOfDischarge:     smt.extractPortOfDischarge(mt700.ShipmentPeriod),
		PartialShipmentAllowed: smt.parseBooleanField(mt700.PartialShipments),
		TransshipmentAllowed:   smt.parseBooleanField(mt700.Transhipment),
		GoodsDescription:    mt700.GoodsServicesDescription,
		RequiredDocuments:   strings.Split(mt700.DocumentsRequired, "\n"),
		Collateral:          collateral,
	}

	lcID, _, err := smt.keeper.IssueLc(ctx, msg)
	if err != nil {
		return "", err
	}

	// Store SWIFT reference mapping
	smt.storeSwiftReference(ctx, mt700.MessageReference, lcID)
	
	return lcID, nil
}

func (smt *SWIFTMTProcessor) findLCByReference(ctx sdk.Context, swiftRef string) (types.LetterOfCredit, error) {
	// Look up LC by SWIFT reference
	lcID := smt.getSwiftReference(ctx, swiftRef)
	if lcID == "" {
		return types.LetterOfCredit{}, fmt.Errorf("LC not found for SWIFT reference: %s", swiftRef)
	}
	
	lc, found := smt.keeper.GetLetterOfCredit(ctx, lcID)
	if !found {
		return types.LetterOfCredit{}, types.ErrLCNotFound
	}
	
	return lc, nil
}

// Storage and validation methods

func (smt *SWIFTMTProcessor) validateMT700(ctx sdk.Context, mt700 *MT700Message) error {
	var errors []string
	
	// Basic field validation
	if mt700.MessageReference == "" {
		errors = append(errors, "Message reference (field 20) is required")
	}
	if mt700.CurrencyCode == "" || len(mt700.CurrencyCode) != 3 {
		errors = append(errors, "Valid currency code (field 32B) is required")
	}
	if mt700.Amount == "" {
		errors = append(errors, "Amount (field 32B) is required")
	}
	if mt700.DateOfExpiry.IsZero() {
		errors = append(errors, "Expiry date (field 31D) is required")
	}
	if mt700.Applicant == "" {
		errors = append(errors, "Applicant (field 50) is required")
	}
	if mt700.Beneficiary == "" {
		errors = append(errors, "Beneficiary (field 59) is required")
	}
	
	// Business logic validation
	if mt700.DateOfExpiry.Before(time.Now()) {
		errors = append(errors, "Expiry date cannot be in the past")
	}
	
	if len(errors) > 0 {
		return fmt.Errorf("MT700 validation failed: %s", strings.Join(errors, "; "))
	}
	
	return nil
}

func (smt *SWIFTMTProcessor) validateMT710(ctx sdk.Context, mt710 *MT710Message) error {
	if mt710.IssuingBankReference == "" {
		return fmt.Errorf("issuing bank reference is required")
	}
	if mt710.AdvisingBankReference == "" {
		return fmt.Errorf("advising bank reference is required")
	}
	return nil
}

func (smt *SWIFTMTProcessor) validateMT720(ctx sdk.Context, mt720 *MT720Message) error {
	if mt720.OriginalCreditRef == "" {
		return fmt.Errorf("original credit reference is required")
	}
	if mt720.SecondBeneficiary == "" {
		return fmt.Errorf("second beneficiary is required")
	}
	return nil
}

func (smt *SWIFTMTProcessor) validateMT750(ctx sdk.Context, mt750 *MT750Message) error {
	if mt750.DocumentaryCreditsRef == "" {
		return fmt.Errorf("documentary credit reference is required")
	}
	if mt750.DiscrepancyDetails == "" {
		return fmt.Errorf("discrepancy details are required")
	}
	return nil
}

func (smt *SWIFTMTProcessor) storeMTMessage(ctx sdk.Context, msg MTMessage) error {
	store := ctx.KVStore(smt.keeper.storeKey)
	key := []byte(fmt.Sprintf("swift_mt:%s:%s", msg.MessageType, msg.MessageReference))
	bz := smt.keeper.cdc.MustMarshal(&msg)
	store.Set(key, bz)
	return nil
}

// Utility methods for conversion and mapping

func (smt *SWIFTMTProcessor) bicToPartyID(ctx sdk.Context, bic string) string {
	// In a real implementation, this would maintain a BIC-to-party mapping
	// For now, return BIC as party ID
	return bic
}

func (smt *SWIFTMTProcessor) parseAmountToSDK(amountStr, currencyCode string) (sdk.Coin, error) {
	// Parse SWIFT amount format to SDK coin
	// Remove comma separators and parse decimal
	cleanAmount := strings.ReplaceAll(amountStr, ",", "")
	if cleanAmount == "" {
		return sdk.Coin{}, fmt.Errorf("empty amount")
	}
	
	// Convert to integer (assuming 2 decimal places for most currencies)
	amountFloat, err := strconv.ParseFloat(cleanAmount, 64)
	if err != nil {
		return sdk.Coin{}, err
	}
	
	// Convert to smallest unit (multiply by 100 for 2 decimal places)
	amountInt := int64(amountFloat * 100)
	
	return sdk.NewInt64Coin(strings.ToLower(currencyCode), amountInt), nil
}

func (smt *SWIFTMTProcessor) parsePartyFromField(fieldValue string) string {
	// Parse party information from SWIFT field format
	// This would extract party ID from formatted field
	return fieldValue // Simplified
}

func (smt *SWIFTMTProcessor) parseShipmentDate(shipmentPeriod string) time.Time {
	// Parse shipment date from period string
	// This would parse various SWIFT date formats
	return time.Now().AddDate(0, 1, 0) // Default to 1 month from now
}

func (smt *SWIFTMTProcessor) extractIncoterms(fieldValue string) string {
	// Extract Incoterms from field value
	incoterms := []string{"FOB", "CIF", "CFR", "EXW", "FCA", "CPT", "CIP", "DAT", "DAP", "DDP"}
	fieldUpper := strings.ToUpper(fieldValue)
	
	for _, term := range incoterms {
		if strings.Contains(fieldUpper, term) {
			return term
		}
	}
	return "FOB" // Default
}

func (smt *SWIFTMTProcessor) extractPortOfLoading(fieldValue string) string {
	// Extract port of loading from field value
	// This would use port code databases
	return "DEFAULT_PORT"
}

func (smt *SWIFTMTProcessor) extractPortOfDischarge(fieldValue string) string {
	// Extract port of discharge from field value
	return "DEFAULT_PORT"
}

func (smt *SWIFTMTProcessor) parseBooleanField(fieldValue string) bool {
	return !strings.Contains(strings.ToUpper(fieldValue), "NOT ALLOWED")
}

func (smt *SWIFTMTProcessor) storeSwiftReference(ctx sdk.Context, swiftRef, lcID string) {
	store := ctx.KVStore(smt.keeper.storeKey)
	key := []byte("swift_ref:" + swiftRef)
	store.Set(key, []byte(lcID))
}

func (smt *SWIFTMTProcessor) getSwiftReference(ctx sdk.Context, swiftRef string) string {
	store := ctx.KVStore(smt.keeper.storeKey)
	key := []byte("swift_ref:" + swiftRef)
	bz := store.Get(key)
	if bz == nil {
		return ""
	}
	return string(bz)
}

// Additional helper methods for MT720, MT750 processing would go here...

func (smt *SWIFTMTProcessor) updateLCWithAdvisingBank(ctx sdk.Context, lcID, advisingBankBIC string, adviceDate time.Time) error {
	lc, found := smt.keeper.GetLetterOfCredit(ctx, lcID)
	if !found {
		return types.ErrLCNotFound
	}
	
	lc.AdvisingBankId = smt.bicToPartyID(ctx, advisingBankBIC)
	lc.Status = "advised"
	lc.UpdatedAt = ctx.BlockTime()
	
	smt.keeper.SetLetterOfCredit(ctx, lc)
	return nil
}

func (smt *SWIFTMTProcessor) createTransferredLC(ctx sdk.Context, originalLC types.LetterOfCredit, mt720 *MT720Message) (string, error) {
	// Create a new LC based on the transfer
	// This would involve complex business logic for transfers
	return "TRANSFERRED_LC_ID", nil
}

func (smt *SWIFTMTProcessor) updateLCWithDiscrepancy(ctx sdk.Context, lcID string, mt750 *MT750Message) error {
	lc, found := smt.keeper.GetLetterOfCredit(ctx, lcID)
	if !found {
		return types.ErrLCNotFound
	}
	
	lc.Status = "discrepancy_advised"
	lc.UpdatedAt = ctx.BlockTime()
	
	smt.keeper.SetLetterOfCredit(ctx, lc)
	return nil
}