package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/deshchain/x/tradefinance/types"
)

// Additional types for UCP 600 compliance that extend the existing trade finance types

// DocumentPresentation represents a collection of documents presented for LC payment
type DocumentPresentation struct {
	ID               string                  `json:"id"`
	LcID             string                  `json:"lc_id"`
	PresentorID      string                  `json:"presentor_id"`
	PresentingBank   string                  `json:"presenting_bank"`
	Documents        []types.TradeDocument   `json:"documents"`
	PresentationDate time.Time               `json:"presentation_date"`
	Status           string                  `json:"status"` // pending, examined, accepted, rejected
	ExaminerID       string                  `json:"examiner_id"`
	ExaminationDate  *time.Time              `json:"examination_date,omitempty"`
	IsCompliant      bool                    `json:"is_compliant"`
	ComplianceResult *UCP600ComplianceResult `json:"compliance_result,omitempty"`
	CreatedAt        time.Time               `json:"created_at"`
	UpdatedAt        time.Time               `json:"updated_at"`
}

// UCP600ComplianceRecord stores compliance check results in the blockchain
type UCP600ComplianceRecord struct {
	ID               string        `json:"id"`
	LcID             string        `json:"lc_id"`
	PresentationID   string        `json:"presentation_id"`
	ExaminationDate  time.Time     `json:"examination_date"`
	ExaminerID       string        `json:"examiner_id"`
	IsCompliant      bool          `json:"is_compliant"`
	ComplianceScore  int           `json:"compliance_score"`
	DiscrepancyCount int           `json:"discrepancy_count"`
	ProcessingTime   int64         `json:"processing_time"` // milliseconds
	Status           string        `json:"status"`          // compliant, discrepant, rejected
	UCP600Version    string        `json:"ucp600_version"`  // e.g., "UCP 600 2007"
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`
}

// Extended TradeDocument fields for UCP 600 compliance
// These would be added to the existing TradeDocument type or used as extensions

// DocumentMetadata contains UCP 600 specific document metadata
type DocumentMetadata struct {
	DocumentNumber       string    `json:"document_number"`
	IssueDate           time.Time `json:"issue_date"`
	IssuingAuthority    string    `json:"issuing_authority"`
	IsOriginal          bool      `json:"is_original"`
	NumberOfOriginals   int       `json:"number_of_originals"`
	NumberOfCopies      int       `json:"number_of_copies"`
	IsClean             bool      `json:"is_clean"`
	HasEndorsement      bool      `json:"has_endorsement"`
	EndorsementDetails  string    `json:"endorsement_details"`
	Signatures          []string  `json:"signatures"`
	NotarizedBy         string    `json:"notarized_by,omitempty"`
	AuthenticationProof string    `json:"authentication_proof,omitempty"`
}

// TransportDocumentDetails contains transport-specific information
type TransportDocumentDetails struct {
	VesselName          string    `json:"vessel_name,omitempty"`
	VoyageNumber        string    `json:"voyage_number,omitempty"`
	FlightNumber        string    `json:"flight_number,omitempty"`
	TruckNumber         string    `json:"truck_number,omitempty"`
	TrainNumber         string    `json:"train_number,omitempty"`
	PortOfLoading       string    `json:"port_of_loading"`
	PortOfDischarge     string    `json:"port_of_discharge"`
	DepartureAirport    string    `json:"departure_airport,omitempty"`
	DestinationAirport  string    `json:"destination_airport,omitempty"`
	OnBoardDate         time.Time `json:"on_board_date"`
	ShipmentDate        time.Time `json:"shipment_date"`
	FreightStatus       string    `json:"freight_status"` // paid, collect, prepaid
	IsClean             bool      `json:"is_clean"`
	TransshipmentAllowed bool     `json:"transshipment_allowed"`
	PartialShipmentAllowed bool   `json:"partial_shipment_allowed"`
}

// InsuranceDocumentDetails contains insurance-specific information
type InsuranceDocumentDetails struct {
	PolicyNumber      string    `json:"policy_number"`
	CertificateNumber string    `json:"certificate_number"`
	InsuredAmount     sdk.Coin  `json:"insured_amount"`
	Currency          string    `json:"currency"`
	CoverageType      string    `json:"coverage_type"` // all risks, named perils, etc.
	Deductible        sdk.Coin  `json:"deductible"`
	EffectiveDate     time.Time `json:"effective_date"`
	ExpiryDate        time.Time `json:"expiry_date"`
	InsuredParty      string    `json:"insured_party"`
	Beneficiary       string    `json:"beneficiary"`
	RisksSpecified    []string  `json:"risks_specified"`
	Exclusions        []string  `json:"exclusions"`
}

// Commercial Invoice specific fields
type CommercialInvoiceDetails struct {
	InvoiceNumber     string    `json:"invoice_number"`
	InvoiceDate       time.Time `json:"invoice_date"`
	BuyerInfo         string    `json:"buyer_info"`
	SellerInfo        string    `json:"seller_info"`
	PaymentTerms      string    `json:"payment_terms"`
	Currency          string    `json:"currency"`
	TotalAmount       sdk.Coin  `json:"total_amount"`
	NetAmount         sdk.Coin  `json:"net_amount"`
	TaxAmount         sdk.Coin  `json:"tax_amount"`
	DiscountAmount    sdk.Coin  `json:"discount_amount"`
	FreightCharges    sdk.Coin  `json:"freight_charges"`
	InsuranceCharges  sdk.Coin  `json:"insurance_charges"`
	OtherCharges      sdk.Coin  `json:"other_charges"`
	GoodsDescription  string    `json:"goods_description"`
	CountryOfOrigin   string    `json:"country_of_origin"`
	HSCodes          []string   `json:"hs_codes"`
	PackingDetails    string    `json:"packing_details"`
}

// Helper methods for the keeper

// GetDocumentPresentation retrieves a document presentation
func (k Keeper) GetDocumentPresentation(ctx sdk.Context, presentationID string) (DocumentPresentation, bool) {
	store := ctx.KVStore(k.storeKey)
	key := []byte("presentation:" + presentationID)
	bz := store.Get(key)
	if bz == nil {
		return DocumentPresentation{}, false
	}

	var presentation DocumentPresentation
	k.cdc.MustUnmarshal(bz, &presentation)
	return presentation, true
}

// SetDocumentPresentation stores a document presentation
func (k Keeper) SetDocumentPresentation(ctx sdk.Context, presentation DocumentPresentation) {
	store := ctx.KVStore(k.storeKey)
	key := []byte("presentation:" + presentation.ID)
	bz := k.cdc.MustMarshal(&presentation)
	store.Set(key, bz)
}

// GetUCP600ComplianceRecord retrieves a compliance record
func (k Keeper) GetUCP600ComplianceRecord(ctx sdk.Context, recordID string) (UCP600ComplianceRecord, bool) {
	store := ctx.KVStore(k.storeKey)
	key := []byte("ucp600:" + recordID)
	bz := store.Get(key)
	if bz == nil {
		return UCP600ComplianceRecord{}, false
	}

	var record UCP600ComplianceRecord
	k.cdc.MustUnmarshal(bz, &record)
	return record, true
}

// SetUCP600ComplianceRecord stores a compliance record
func (k Keeper) SetUCP600ComplianceRecord(ctx sdk.Context, record UCP600ComplianceRecord) {
	store := ctx.KVStore(k.storeKey)
	key := []byte("ucp600:" + record.ID)
	bz := k.cdc.MustMarshal(&record)
	store.Set(key, bz)
}

// GetPartyIDByAddress returns party ID for a given address
func (k Keeper) GetPartyIDByAddress(ctx sdk.Context, address string) string {
	// This would typically involve looking up party information by address
	// For now, return the address itself as party ID
	return address
}

// GetAllDocumentPresentations returns all document presentations for an LC
func (k Keeper) GetAllDocumentPresentations(ctx sdk.Context, lcID string) []DocumentPresentation {
	store := ctx.KVStore(k.storeKey)
	prefix := []byte("presentation:")
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	var presentations []DocumentPresentation
	for ; iterator.Valid(); iterator.Next() {
		var presentation DocumentPresentation
		k.cdc.MustUnmarshal(iterator.Value(), &presentation)
		if presentation.LcID == lcID {
			presentations = append(presentations, presentation)
		}
	}

	return presentations
}

// GetComplianceRecordsByLC returns all compliance records for an LC
func (k Keeper) GetComplianceRecordsByLC(ctx sdk.Context, lcID string) []UCP600ComplianceRecord {
	store := ctx.KVStore(k.storeKey)
	prefix := []byte("ucp600:")
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	var records []UCP600ComplianceRecord
	for ; iterator.Valid(); iterator.Next() {
		var record UCP600ComplianceRecord
		k.cdc.MustUnmarshal(iterator.Value(), &record)
		if record.LcID == lcID {
			records = append(records, record)
		}
	}

	return records
}

// IsUCP600CompliantBank checks if a bank follows UCP 600 standards
func (k Keeper) IsUCP600CompliantBank(ctx sdk.Context, bankID string) bool {
	// This would check if the bank is registered as UCP 600 compliant
	// For now, assume all registered banks are compliant
	_, found := k.GetTradeParty(ctx, bankID)
	return found
}

// ValidateUCP600DocumentStructure validates document structure according to UCP 600
func (k Keeper) ValidateUCP600DocumentStructure(doc types.TradeDocument) []string {
	var issues []string

	// Basic document validation
	if doc.DocumentNumber == "" {
		issues = append(issues, "Document number is required")
	}

	if doc.IssueDate.IsZero() {
		issues = append(issues, "Issue date is required")
	}

	// Document type specific validation
	switch doc.DocumentType {
	case "commercial_invoice":
		if doc.Amount.IsZero() {
			issues = append(issues, "Commercial invoice must have an amount")
		}
		if doc.BuyerInfo == "" {
			issues = append(issues, "Commercial invoice must specify buyer")
		}
		if doc.SellerInfo == "" {
			issues = append(issues, "Commercial invoice must specify seller")
		}

	case "bill_of_lading":
		if doc.PortOfLoading == "" {
			issues = append(issues, "Bill of lading must specify port of loading")
		}
		if doc.PortOfDischarge == "" {
			issues = append(issues, "Bill of lading must specify port of discharge")
		}
		if doc.OnBoardDate.IsZero() {
			issues = append(issues, "Bill of lading must have on board date")
		}

	case "insurance_certificate", "insurance_policy":
		if doc.InsuranceAmount.IsZero() {
			issues = append(issues, "Insurance document must specify coverage amount")
		}
		if doc.InsuredParty == "" {
			issues = append(issues, "Insurance document must specify insured party")
		}
	}

	return issues
}

// CalculateUCP600ExaminationFee calculates examination fees per UCP 600 standards
func (k Keeper) CalculateUCP600ExaminationFee(ctx sdk.Context, lcAmount sdk.Coin) sdk.Coin {
	params := k.GetParams(ctx)
	
	// Typical examination fee is 0.125% (12.5 basis points) of LC amount
	// with minimum and maximum limits
	feeRate := sdk.NewDecWithPrec(125, 6) // 0.000125 = 0.0125%
	feeAmount := sdk.NewDecFromInt(lcAmount.Amount).Mul(feeRate).TruncateInt()
	
	// Apply minimum fee
	minFee := sdk.NewInt(int64(params.Fees.MinExaminationFee))
	if feeAmount.LT(minFee) {
		feeAmount = minFee
	}
	
	// Apply maximum fee
	maxFee := sdk.NewInt(int64(params.Fees.MaxExaminationFee))
	if feeAmount.GT(maxFee) {
		feeAmount = maxFee
	}
	
	return sdk.NewCoin("dinr", feeAmount)
}

// GenerateUCP600Report generates a comprehensive UCP 600 compliance report
func (k Keeper) GenerateUCP600Report(ctx sdk.Context, lcID string) (UCP600ComplianceReport, error) {
	lc, found := k.GetLetterOfCredit(ctx, lcID)
	if !found {
		return UCP600ComplianceReport{}, types.ErrLCNotFound
	}

	presentations := k.GetAllDocumentPresentations(ctx, lcID)
	records := k.GetComplianceRecordsByLC(ctx, lcID)

	report := UCP600ComplianceReport{
		LcID:                lcID,
		LcNumber:           lc.LcNumber,
		GeneratedAt:        ctx.BlockTime(),
		TotalPresentations: len(presentations),
		ComplianceRecords:  records,
		OverallCompliance:  k.calculateOverallCompliance(records),
	}

	return report, nil
}

// UCP600ComplianceReport represents a comprehensive compliance report
type UCP600ComplianceReport struct {
	LcID                string                    `json:"lc_id"`
	LcNumber            string                    `json:"lc_number"`
	GeneratedAt         time.Time                 `json:"generated_at"`
	TotalPresentations  int                       `json:"total_presentations"`
	CompliantPresentations int                    `json:"compliant_presentations"`
	DiscrepantPresentations int                   `json:"discrepant_presentations"`
	ComplianceRecords   []UCP600ComplianceRecord  `json:"compliance_records"`
	OverallCompliance   ComplianceMetrics         `json:"overall_compliance"`
}

// ComplianceMetrics provides statistical analysis of compliance
type ComplianceMetrics struct {
	ComplianceRate      float64 `json:"compliance_rate"`      // percentage
	AverageScore        float64 `json:"average_score"`        // 0-100
	CommonDiscrepancies []string `json:"common_discrepancies"`
	ProcessingTimeAvg   int64   `json:"processing_time_avg"`  // milliseconds
	FirstTimeCompliance float64 `json:"first_time_compliance"` // percentage
}

// Helper function to calculate overall compliance metrics
func (k Keeper) calculateOverallCompliance(records []UCP600ComplianceRecord) ComplianceMetrics {
	if len(records) == 0 {
		return ComplianceMetrics{}
	}

	compliantCount := 0
	totalScore := 0
	totalProcessingTime := int64(0)

	for _, record := range records {
		if record.IsCompliant {
			compliantCount++
		}
		totalScore += record.ComplianceScore
		totalProcessingTime += record.ProcessingTime
	}

	return ComplianceMetrics{
		ComplianceRate:      float64(compliantCount) / float64(len(records)) * 100,
		AverageScore:        float64(totalScore) / float64(len(records)),
		ProcessingTimeAvg:   totalProcessingTime / int64(len(records)),
		FirstTimeCompliance: float64(compliantCount) / float64(len(records)) * 100, // Simplified
		CommonDiscrepancies: []string{}, // Would require more complex analysis
	}
}