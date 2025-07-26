package keeper

import (
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/tradefinance/types"
)

// UCP600ComplianceEngine implements Uniform Customs and Practice for Documentary Credits (UCP 600)
type UCP600ComplianceEngine struct {
	keeper *Keeper
}

// NewUCP600ComplianceEngine creates a new UCP 600 compliance engine
func NewUCP600ComplianceEngine(k *Keeper) *UCP600ComplianceEngine {
	return &UCP600ComplianceEngine{
		keeper: k,
	}
}

// UCP600ComplianceResult represents the result of UCP 600 compliance check
type UCP600ComplianceResult struct {
	IsCompliant      bool                    `json:"is_compliant"`
	OverallStatus    string                  `json:"overall_status"` // compliant, discrepant, rejected
	Discrepancies    []UCP600Discrepancy     `json:"discrepancies"`
	DocumentReviews  []DocumentReview        `json:"document_reviews"`
	ComplianceScore  int                     `json:"compliance_score"` // 0-100
	ProcessingTime   time.Duration           `json:"processing_time"`
	ExaminationDate  time.Time               `json:"examination_date"`
	DeadlineStatus   DeadlineComplianceCheck `json:"deadline_status"`
}

// UCP600Discrepancy represents a discrepancy found during examination
type UCP600Discrepancy struct {
	ID          string                 `json:"id"`
	Type        UCP600DiscrepancyType  `json:"type"`
	Category    string                 `json:"category"` // major, minor, technical
	Severity    UCP600DiscrepancySeverity `json:"severity"`
	DocumentID  string                 `json:"document_id"`
	Article     string                 `json:"article"`     // UCP 600 article reference
	Description string                 `json:"description"`
	Suggestion  string                 `json:"suggestion"`
	IsWaivable  bool                   `json:"is_waivable"`
	Impact      string                 `json:"impact"`
	DetectedAt  time.Time              `json:"detected_at"`
}

// DocumentReview represents the review of a single document
type DocumentReview struct {
	DocumentID     string                    `json:"document_id"`
	DocumentType   string                    `json:"document_type"`
	Status         DocumentComplianceStatus  `json:"status"`
	IsCompliant    bool                      `json:"is_compliant"`
	Findings       []ComplianceFinding       `json:"findings"`
	ReviewTime     time.Duration             `json:"review_time"`
	ReviewerNotes  string                    `json:"reviewer_notes"`
	UCP600Articles []string                  `json:"ucp600_articles"`
}

// ComplianceFinding represents a specific finding during document review
type ComplianceFinding struct {
	Type        string    `json:"type"`        // requirement, discrepancy, warning
	Description string    `json:"description"`
	Field       string    `json:"field"`
	Expected    string    `json:"expected"`
	Actual      string    `json:"actual"`
	Severity    string    `json:"severity"`
	Article     string    `json:"article"`
	FoundAt     time.Time `json:"found_at"`
}

// DeadlineComplianceCheck represents compliance with presentation deadlines
type DeadlineComplianceCheck struct {
	PresentationDeadline      time.Time `json:"presentation_deadline"`
	ActualPresentationTime    time.Time `json:"actual_presentation_time"`
	IsWithinDeadline         bool      `json:"is_within_deadline"`
	TimeRemaining            int64     `json:"time_remaining"`      // seconds
	ExpiryDate               time.Time `json:"expiry_date"`
	ShipmentDeadline         time.Time `json:"shipment_deadline"`
	IsShipmentDeadlineValid  bool      `json:"is_shipment_deadline_valid"`
}

// Enums for UCP 600 compliance
type UCP600DiscrepancyType int

const (
	DiscrepancyDocumentMissing UCP600DiscrepancyType = iota
	DiscrepancyDocumentIncomplete
	DiscrepancyDocumentInconsistent
	DiscrepancyAmountExcessive
	DiscrepancyDatesInconsistent
	DiscrepancyDescriptionMismatch
	DiscrepancyTransportInconsistent
	DiscrepancyInsuranceInsufficient
	DiscrepancySignatureMissing
	DiscrepancyEndorsementMissing
	DiscrepancyPresentationLate
	DiscrepancyExpiryPassed
)

type UCP600DiscrepancySeverity int

const (
	SeverityMinor UCP600DiscrepancySeverity = iota
	SeverityMajor
	SeverityCritical
)

type DocumentComplianceStatus int

const (
	StatusPending DocumentComplianceStatus = iota
	StatusCompliant
	StatusDiscrepant
	StatusRejected
	StatusUnderReview
)

// PerformUCP600Compliance performs comprehensive UCP 600 compliance check
func (uce *UCP600ComplianceEngine) PerformUCP600Compliance(
	ctx sdk.Context,
	lcID string,
	presentationID string,
) (*UCP600ComplianceResult, error) {
	startTime := time.Now()
	
	// Get LC and presentation
	lc, found := uce.keeper.GetLetterOfCredit(ctx, lcID)
	if !found {
		return nil, types.ErrLCNotFound
	}

	presentation, found := uce.keeper.GetDocumentPresentation(ctx, presentationID)
	if !found {
		return nil, types.ErrPresentationNotFound
	}

	result := &UCP600ComplianceResult{
		IsCompliant:     true,
		OverallStatus:   "under_review",
		Discrepancies:   []UCP600Discrepancy{},
		DocumentReviews: []DocumentReview{},
		ExaminationDate: ctx.BlockTime(),
	}

	// Article 14: Standard for examination of documents
	// Banks must examine documents within a maximum of 5 banking days
	result.DeadlineStatus = uce.checkPresentationDeadlines(ctx, lc, presentation)

	// Article 15: Complying presentation
	if err := uce.checkComplyingPresentation(ctx, lc, presentation, result); err != nil {
		return nil, err
	}

	// Article 16: Discrepant documents
	uce.examineDocuments(ctx, lc, presentation, result)

	// Article 8: Interpretation of data
	uce.checkDataInterpretation(ctx, lc, presentation, result)

	// Article 20: Commercial invoice requirements
	uce.validateCommercialInvoice(ctx, lc, presentation, result)

	// Article 19-27: Transport documents
	uce.validateTransportDocuments(ctx, lc, presentation, result)

	// Article 28: Insurance documents
	uce.validateInsuranceDocuments(ctx, lc, presentation, result)

	// Calculate overall compliance
	result.ComplianceScore = uce.calculateComplianceScore(result)
	result.ProcessingTime = time.Since(startTime)

	// Determine final status
	if len(result.Discrepancies) == 0 {
		result.OverallStatus = "compliant"
		result.IsCompliant = true
	} else {
		criticalDiscrepancies := uce.countCriticalDiscrepancies(result.Discrepancies)
		if criticalDiscrepancies > 0 {
			result.OverallStatus = "rejected"
			result.IsCompliant = false
		} else {
			result.OverallStatus = "discrepant"
			result.IsCompliant = false
		}
	}

	return result, nil
}

// checkPresentationDeadlines validates Article 6 & 14 compliance (presentation deadlines)
func (uce *UCP600ComplianceEngine) checkPresentationDeadlines(
	ctx sdk.Context,
	lc types.LetterOfCredit,
	presentation types.DocumentPresentation,
) DeadlineComplianceCheck {
	currentTime := ctx.BlockTime()
	
	// UCP 600 Article 6(d): 21 days from shipment date or expiry date, whichever is earlier
	presentationDeadline := lc.ExpiryDate
	if !lc.LatestShipmentDate.IsZero() {
		shipmentDeadline := lc.LatestShipmentDate.AddDate(0, 0, 21)
		if shipmentDeadline.Before(presentationDeadline) {
			presentationDeadline = shipmentDeadline
		}
	}

	return DeadlineComplianceCheck{
		PresentationDeadline:     presentationDeadline,
		ActualPresentationTime:   presentation.PresentationDate,
		IsWithinDeadline:        presentation.PresentationDate.Before(presentationDeadline),
		TimeRemaining:           presentationDeadline.Sub(currentTime).Nanoseconds() / 1e9,
		ExpiryDate:              lc.ExpiryDate,
		ShipmentDeadline:        lc.LatestShipmentDate,
		IsShipmentDeadlineValid: !lc.LatestShipmentDate.IsZero() && lc.LatestShipmentDate.After(currentTime),
	}
}

// checkComplyingPresentation validates Article 15 requirements
func (uce *UCP600ComplianceEngine) checkComplyingPresentation(
	ctx sdk.Context,
	lc types.LetterOfCredit,
	presentation types.DocumentPresentation,
	result *UCP600ComplianceResult,
) error {
	// Check if presentation is made by or on behalf of the beneficiary
	if presentation.PresentorId != lc.BeneficiaryId {
		result.Discrepancies = append(result.Discrepancies, UCP600Discrepancy{
			ID:          fmt.Sprintf("DISC-%d", len(result.Discrepancies)+1),
			Type:        DiscrepancyDocumentInconsistent,
			Category:    "major",
			Severity:    SeverityMajor,
			Article:     "Article 15",
			Description: "Presentation not made by or on behalf of beneficiary",
			Suggestion:  "Ensure presentation is made by the beneficiary or authorized party",
			IsWaivable:  false,
			Impact:      "Presentation may be rejected",
			DetectedAt:  ctx.BlockTime(),
		})
	}

	// Check required documents are present
	requiredDocs := lc.RequiredDocuments
	presentedDocs := uce.extractPresentedDocumentTypes(presentation.Documents)

	for _, required := range requiredDocs {
		if !uce.containsDocumentType(presentedDocs, required) {
			result.Discrepancies = append(result.Discrepancies, UCP600Discrepancy{
				ID:          fmt.Sprintf("DISC-%d", len(result.Discrepancies)+1),
				Type:        DiscrepancyDocumentMissing,
				Category:    "critical",
				Severity:    SeverityCritical,
				Article:     "Article 15",
				Description: fmt.Sprintf("Required document missing: %s", required),
				Suggestion:  fmt.Sprintf("Include the required %s document", required),
				IsWaivable:  false,
				Impact:      "Presentation will be rejected",
				DetectedAt:  ctx.BlockTime(),
			})
		}
	}

	return nil
}

// examineDocuments performs detailed document examination per Article 14
func (uce *UCP600ComplianceEngine) examineDocuments(
	ctx sdk.Context,
	lc types.LetterOfCredit,
	presentation types.DocumentPresentation,
	result *UCP600ComplianceResult,
) {
	for _, doc := range presentation.Documents {
		review := DocumentReview{
			DocumentID:    doc.Id,
			DocumentType:  doc.DocumentType,
			Status:        StatusUnderReview,
			IsCompliant:   true,
			Findings:      []ComplianceFinding{},
			ReviewerNotes: "",
			UCP600Articles: []string{},
		}

		startReview := time.Now()

		// Examine document based on type
		switch strings.ToLower(doc.DocumentType) {
		case "commercial_invoice":
			uce.examineCommercialInvoice(ctx, lc, doc, &review)
		case "bill_of_lading", "sea_waybill":
			uce.examineMarineTransportDocument(ctx, lc, doc, &review)
		case "air_waybill":
			uce.examineAirTransportDocument(ctx, lc, doc, &review)
		case "insurance_certificate", "insurance_policy":
			uce.examineInsuranceDocument(ctx, lc, doc, &review)
		case "certificate_of_origin":
			uce.examineCertificateOfOrigin(ctx, lc, doc, &review)
		case "packing_list":
			uce.examinePackingList(ctx, lc, doc, &review)
		default:
			uce.examineGenericDocument(ctx, lc, doc, &review)
		}

		review.ReviewTime = time.Since(startReview)

		// Determine final status for this document
		if len(review.Findings) == 0 {
			review.Status = StatusCompliant
		} else {
			hasRejectableFindings := false
			for _, finding := range review.Findings {
				if finding.Severity == "critical" {
					hasRejectableFindings = true
					break
				}
			}
			if hasRejectableFindings {
				review.Status = StatusRejected
				review.IsCompliant = false
			} else {
				review.Status = StatusDiscrepant
				review.IsCompliant = false
			}
		}

		result.DocumentReviews = append(result.DocumentReviews, review)
	}
}

// examineCommercialInvoice validates commercial invoice per Article 20
func (uce *UCP600ComplianceEngine) examineCommercialInvoice(
	ctx sdk.Context,
	lc types.LetterOfCredit,
	doc types.TradeDocument,
	review *DocumentReview,
) {
	review.UCP600Articles = append(review.UCP600Articles, "Article 20")

	// Check invoice amount does not exceed LC amount
	if doc.Amount.IsGT(lc.Amount) {
		review.Findings = append(review.Findings, ComplianceFinding{
			Type:        "discrepancy",
			Description: "Invoice amount exceeds LC amount",
			Field:       "amount",
			Expected:    lc.Amount.String(),
			Actual:      doc.Amount.String(),
			Severity:    "major",
			Article:     "Article 20",
			FoundAt:     ctx.BlockTime(),
		})
	}

	// Check goods description matches
	if !uce.isGoodsDescriptionMatching(lc.GoodsDescription, doc.GoodsDescription) {
		review.Findings = append(review.Findings, ComplianceFinding{
			Type:        "discrepancy",
			Description: "Goods description does not match LC requirements",
			Field:       "goods_description",
			Expected:    lc.GoodsDescription,
			Actual:      doc.GoodsDescription,
			Severity:    "major",
			Article:     "Article 20",
			FoundAt:     ctx.BlockTime(),
		})
	}

	// Check parties (applicant and beneficiary)
	if doc.BuyerInfo != lc.ApplicantId {
		review.Findings = append(review.Findings, ComplianceFinding{
			Type:        "discrepancy",
			Description: "Buyer information does not match applicant",
			Field:       "buyer_info",
			Expected:    lc.ApplicantId,
			Actual:      doc.BuyerInfo,
			Severity:    "major",
			Article:     "Article 20",
			FoundAt:     ctx.BlockTime(),
		})
	}

	if doc.SellerInfo != lc.BeneficiaryId {
		review.Findings = append(review.Findings, ComplianceFinding{
			Type:        "discrepancy",
			Description: "Seller information does not match beneficiary",
			Field:       "seller_info",
			Expected:    lc.BeneficiaryId,
			Actual:      doc.SellerInfo,
			Severity:    "major",
			Article:     "Article 20",
			FoundAt:     ctx.BlockTime(),
		})
	}
}

// examineMarineTransportDocument validates marine transport documents per Articles 20-23
func (uce *UCP600ComplianceEngine) examineMarineTransportDocument(
	ctx sdk.Context,
	lc types.LetterOfCredit,
	doc types.TradeDocument,
	review *DocumentReview,
) {
	review.UCP600Articles = append(review.UCP600Articles, "Article 20", "Article 21", "Article 22", "Article 23")

	// Check port of loading
	if lc.PortOfLoading != "" && !strings.Contains(doc.PortOfLoading, lc.PortOfLoading) {
		review.Findings = append(review.Findings, ComplianceFinding{
			Type:        "discrepancy",
			Description: "Port of loading does not match LC requirements",
			Field:       "port_of_loading",
			Expected:    lc.PortOfLoading,
			Actual:      doc.PortOfLoading,
			Severity:    "major",
			Article:     "Article 20",
			FoundAt:     ctx.BlockTime(),
		})
	}

	// Check port of discharge
	if lc.PortOfDischarge != "" && !strings.Contains(doc.PortOfDischarge, lc.PortOfDischarge) {
		review.Findings = append(review.Findings, ComplianceFinding{
			Type:        "discrepancy",
			Description: "Port of discharge does not match LC requirements",
			Field:       "port_of_discharge",
			Expected:    lc.PortOfDischarge,
			Actual:      doc.PortOfDischarge,
			Severity:    "major",
			Article:     "Article 20",
			FoundAt:     ctx.BlockTime(),
		})
	}

	// Check clean on board notation (for ocean bills of lading)
	if doc.DocumentType == "bill_of_lading" && !doc.IsClean {
		review.Findings = append(review.Findings, ComplianceFinding{
			Type:        "discrepancy",
			Description: "Bill of lading must be clean (no adverse remarks)",
			Field:       "clean_status",
			Expected:    "clean",
			Actual:      "claused",
			Severity:    "major",
			Article:     "Article 21",
			FoundAt:     ctx.BlockTime(),
		})
	}

	// Check on board date against shipment deadline
	if !lc.LatestShipmentDate.IsZero() && doc.OnBoardDate.After(lc.LatestShipmentDate) {
		review.Findings = append(review.Findings, ComplianceFinding{
			Type:        "discrepancy",
			Description: "On board date exceeds latest shipment date",
			Field:       "on_board_date",
			Expected:    lc.LatestShipmentDate.Format("2006-01-02"),
			Actual:      doc.OnBoardDate.Format("2006-01-02"),
			Severity:    "critical",
			Article:     "Article 21",
			FoundAt:     ctx.BlockTime(),
		})
	}
}

// examineAirTransportDocument validates air transport documents per Article 24
func (uce *UCP600ComplianceEngine) examineAirTransportDocument(
	ctx sdk.Context,
	lc types.LetterOfCredit,
	doc types.TradeDocument,
	review *DocumentReview,
) {
	review.UCP600Articles = append(review.UCP600Articles, "Article 24")

	// Check departure airport
	if lc.PortOfLoading != "" && !strings.Contains(doc.DepartureAirport, lc.PortOfLoading) {
		review.Findings = append(review.Findings, ComplianceFinding{
			Type:        "discrepancy",
			Description: "Departure airport does not match LC requirements",
			Field:       "departure_airport",
			Expected:    lc.PortOfLoading,
			Actual:      doc.DepartureAirport,
			Severity:    "major",
			Article:     "Article 24",
			FoundAt:     ctx.BlockTime(),
		})
	}

	// Check destination airport
	if lc.PortOfDischarge != "" && !strings.Contains(doc.DestinationAirport, lc.PortOfDischarge) {
		review.Findings = append(review.Findings, ComplianceFinding{
			Type:        "discrepancy",
			Description: "Destination airport does not match LC requirements",
			Field:       "destination_airport",
			Expected:    lc.PortOfDischarge,
			Actual:      doc.DestinationAirport,
			Severity:    "major",
			Article:     "Article 24",
			FoundAt:     ctx.BlockTime(),
		})
	}

	// Air waybills must show that freight has been paid or is payable at destination
	if doc.FreightStatus == "" {
		review.Findings = append(review.Findings, ComplianceFinding{
			Type:        "requirement",
			Description: "Air waybill must indicate freight payment status",
			Field:       "freight_status",
			Expected:    "freight paid or payable at destination",
			Actual:      "not specified",
			Severity:    "minor",
			Article:     "Article 24",
			FoundAt:     ctx.BlockTime(),
		})
	}
}

// examineInsuranceDocument validates insurance documents per Article 28
func (uce *UCP600ComplianceEngine) examineInsuranceDocument(
	ctx sdk.Context,
	lc types.LetterOfCredit,
	doc types.TradeDocument,
	review *DocumentReview,
) {
	review.UCP600Articles = append(review.UCP600Articles, "Article 28")

	// Check insurance amount (typically 110% of invoice value)
	minInsuranceAmount := lc.Amount.Amount.Mul(sdk.NewInt(110)).Quo(sdk.NewInt(100))
	if doc.InsuranceAmount.LT(sdk.NewCoin(lc.Currency, minInsuranceAmount)) {
		review.Findings = append(review.Findings, ComplianceFinding{
			Type:        "discrepancy",
			Description: "Insurance amount insufficient (should be minimum 110% of LC value)",
			Field:       "insurance_amount",
			Expected:    sdk.NewCoin(lc.Currency, minInsuranceAmount).String(),
			Actual:      doc.InsuranceAmount.String(),
			Severity:    "major",
			Article:     "Article 28",
			FoundAt:     ctx.BlockTime(),
		})
	}

	// Check coverage includes the goods as described in the LC
	if !uce.isGoodsDescriptionMatching(lc.GoodsDescription, doc.GoodsDescription) {
		review.Findings = append(review.Findings, ComplianceFinding{
			Type:        "discrepancy",
			Description: "Insurance coverage does not match goods in LC",
			Field:       "goods_coverage",
			Expected:    lc.GoodsDescription,
			Actual:      doc.GoodsDescription,
			Severity:    "major",
			Article:     "Article 28",
			FoundAt:     ctx.BlockTime(),
		})
	}
}

// Additional examination methods for other document types
func (uce *UCP600ComplianceEngine) examineCertificateOfOrigin(
	ctx sdk.Context,
	lc types.LetterOfCredit,
	doc types.TradeDocument,
	review *DocumentReview,
) {
	// Check country of origin if specified in LC
	if lc.CountryOfOrigin != "" && doc.CountryOfOrigin != lc.CountryOfOrigin {
		review.Findings = append(review.Findings, ComplianceFinding{
			Type:        "discrepancy",
			Description: "Country of origin does not match LC requirements",
			Field:       "country_of_origin",
			Expected:    lc.CountryOfOrigin,
			Actual:      doc.CountryOfOrigin,
			Severity:    "major",
			Article:     "General",
			FoundAt:     ctx.BlockTime(),
		})
	}
}

func (uce *UCP600ComplianceEngine) examinePackingList(
	ctx sdk.Context,
	lc types.LetterOfCredit,
	doc types.TradeDocument,
	review *DocumentReview,
) {
	// Basic validation for packing lists
	if doc.PackingDetails == "" {
		review.Findings = append(review.Findings, ComplianceFinding{
			Type:        "requirement",
			Description: "Packing list should contain packing details",
			Field:       "packing_details",
			Expected:    "detailed packing information",
			Actual:      "missing",
			Severity:    "minor",
			Article:     "General",
			FoundAt:     ctx.BlockTime(),
		})
	}
}

func (uce *UCP600ComplianceEngine) examineGenericDocument(
	ctx sdk.Context,
	lc types.LetterOfCredit,
	doc types.TradeDocument,
	review *DocumentReview,
) {
	// Basic validation for any document
	if doc.DocumentNumber == "" {
		review.Findings = append(review.Findings, ComplianceFinding{
			Type:        "requirement",
			Description: "Document should have a reference number",
			Field:       "document_number",
			Expected:    "document reference number",
			Actual:      "missing",
			Severity:    "minor",
			Article:     "General",
			FoundAt:     ctx.BlockTime(),
		})
	}
}

// Helper methods

func (uce *UCP600ComplianceEngine) checkDataInterpretation(
	ctx sdk.Context,
	lc types.LetterOfCredit,
	presentation types.DocumentPresentation,
	result *UCP600ComplianceResult,
) {
	// Article 8: Data interpretation and spelling variations
	// This would include checks for acceptable variations in:
	// - Company names and addresses
	// - Goods descriptions
	// - Amounts (in words vs figures)
	// - Dates
	
	// Implementation would include fuzzy matching algorithms
	// For now, implementing basic checks
}

func (uce *UCP600ComplianceEngine) validateCommercialInvoice(
	ctx sdk.Context,
	lc types.LetterOfCredit,
	presentation types.DocumentPresentation,
	result *UCP600ComplianceResult,
) {
	// Article 20: Commercial invoice requirements
	hasCommercialInvoice := false
	
	for _, doc := range presentation.Documents {
		if strings.ToLower(doc.DocumentType) == "commercial_invoice" {
			hasCommercialInvoice = true
			break
		}
	}
	
	if !hasCommercialInvoice {
		result.Discrepancies = append(result.Discrepancies, UCP600Discrepancy{
			ID:          fmt.Sprintf("DISC-%d", len(result.Discrepancies)+1),
			Type:        DiscrepancyDocumentMissing,
			Category:    "critical",
			Severity:    SeverityCritical,
			Article:     "Article 20",
			Description: "Commercial invoice is mandatory and missing",
			Suggestion:  "Include a properly prepared commercial invoice",
			IsWaivable:  false,
			Impact:      "Presentation will be rejected",
			DetectedAt:  ctx.BlockTime(),
		})
	}
}

func (uce *UCP600ComplianceEngine) validateTransportDocuments(
	ctx sdk.Context,
	lc types.LetterOfCredit,
	presentation types.DocumentPresentation,
	result *UCP600ComplianceResult,
) {
	// Articles 19-27: Transport documents validation
	hasTransportDoc := false
	acceptableTransportDocs := []string{"bill_of_lading", "sea_waybill", "air_waybill", "road_waybill", "rail_waybill"}
	
	for _, doc := range presentation.Documents {
		for _, acceptable := range acceptableTransportDocs {
			if strings.ToLower(doc.DocumentType) == acceptable {
				hasTransportDoc = true
				break
			}
		}
		if hasTransportDoc {
			break
		}
	}
	
	if !hasTransportDoc {
		result.Discrepancies = append(result.Discrepancies, UCP600Discrepancy{
			ID:          fmt.Sprintf("DISC-%d", len(result.Discrepancies)+1),
			Type:        DiscrepancyDocumentMissing,
			Category:    "critical",
			Severity:    SeverityCritical,
			Article:     "Articles 19-27",
			Description: "Transport document is required but missing",
			Suggestion:  "Include an appropriate transport document",
			IsWaivable:  false,
			Impact:      "Presentation will be rejected",
			DetectedAt:  ctx.BlockTime(),
		})
	}
}

func (uce *UCP600ComplianceEngine) validateInsuranceDocuments(
	ctx sdk.Context,
	lc types.LetterOfCredit,
	presentation types.DocumentPresentation,
	result *UCP600ComplianceResult,
) {
	// Article 28: Insurance documents
	// Check if insurance is required by checking LC terms
	if strings.Contains(strings.ToLower(lc.Incoterms), "cif") || 
	   strings.Contains(strings.ToLower(lc.Incoterms), "cip") {
		
		hasInsuranceDoc := false
		for _, doc := range presentation.Documents {
			if strings.Contains(strings.ToLower(doc.DocumentType), "insurance") {
				hasInsuranceDoc = true
				break
			}
		}
		
		if !hasInsuranceDoc {
			result.Discrepancies = append(result.Discrepancies, UCP600Discrepancy{
				ID:          fmt.Sprintf("DISC-%d", len(result.Discrepancies)+1),
				Type:        DiscrepancyDocumentMissing,
				Category:    "major",
				Severity:    SeverityMajor,
				Article:     "Article 28",
				Description: fmt.Sprintf("Insurance document required for %s terms", lc.Incoterms),
				Suggestion:  "Include insurance certificate or policy",
				IsWaivable:  true,
				Impact:      "May cause presentation to be discrepant",
				DetectedAt:  ctx.BlockTime(),
			})
		}
	}
}

// Utility functions
func (uce *UCP600ComplianceEngine) calculateComplianceScore(result *UCP600ComplianceResult) int {
	if len(result.Discrepancies) == 0 {
		return 100
	}
	
	totalPenalty := 0
	for _, discrepancy := range result.Discrepancies {
		switch discrepancy.Severity {
		case SeverityCritical:
			totalPenalty += 30
		case SeverityMajor:
			totalPenalty += 15
		case SeverityMinor:
			totalPenalty += 5
		}
	}
	
	score := 100 - totalPenalty
	if score < 0 {
		score = 0
	}
	
	return score
}

func (uce *UCP600ComplianceEngine) countCriticalDiscrepancies(discrepancies []UCP600Discrepancy) int {
	count := 0
	for _, d := range discrepancies {
		if d.Severity == SeverityCritical {
			count++
		}
	}
	return count
}

func (uce *UCP600ComplianceEngine) extractPresentedDocumentTypes(docs []types.TradeDocument) []string {
	var types []string
	for _, doc := range docs {
		types = append(types, doc.DocumentType)
	}
	return types
}

func (uce *UCP600ComplianceEngine) containsDocumentType(presented []string, required string) bool {
	for _, p := range presented {
		if strings.EqualFold(p, required) {
			return true
		}
	}
	return false
}

func (uce *UCP600ComplianceEngine) isGoodsDescriptionMatching(lcDescription, docDescription string) bool {
	// Implement fuzzy matching logic for goods description
	// For now, simple case-insensitive contains check
	return strings.Contains(strings.ToLower(docDescription), strings.ToLower(lcDescription))
}

// ProcessUCP600ComplianceResult processes the compliance result and updates LC status
func (k Keeper) ProcessUCP600ComplianceResult(
	ctx sdk.Context,
	lcID string,
	presentationID string,
	result *UCP600ComplianceResult,
) error {
	// Update LC status based on compliance result
	lc, found := k.GetLetterOfCredit(ctx, lcID)
	if !found {
		return types.ErrLCNotFound
	}
	
	if result.IsCompliant {
		lc.Status = "documents_compliant"
	} else {
		lc.Status = "documents_discrepant"
	}
	
	lc.UpdatedAt = ctx.BlockTime()
	k.SetLetterOfCredit(ctx, lc)
	
	// Store compliance result
	complianceRecord := types.UCP600ComplianceRecord{
		LcId:            lcID,
		PresentationId:  presentationID,
		ExaminationDate: result.ExaminationDate,
		IsCompliant:     result.IsCompliant,
		ComplianceScore: result.ComplianceScore,
		DiscrepancyCount: len(result.Discrepancies),
		ProcessingTime:  result.ProcessingTime.Milliseconds(),
		Status:          result.OverallStatus,
	}
	
	k.SetUCP600ComplianceRecord(ctx, complianceRecord)
	
	return nil
}