package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/namo/x/shikshaamitra/types"
)

// InstitutionManager handles educational institution partnerships and verification
type InstitutionManager struct {
	keeper Keeper
}

// NewInstitutionManager creates a new institution manager
func NewInstitutionManager(keeper Keeper) *InstitutionManager {
	return &InstitutionManager{
		keeper: keeper,
	}
}

// PartnerInstitution represents a partnered educational institution
type PartnerInstitution struct {
	InstitutionID        string                      `json:"institution_id"`
	InstitutionInfo      types.InstitutionInfo       `json:"institution_info"`
	PartnershipLevel     string                      `json:"partnership_level"` // GOLD, SILVER, BRONZE, STANDARD
	AccreditationStatus  types.AccreditationStatus   `json:"accreditation_status"`
	PlacementRecord      types.PlacementRecord       `json:"placement_record"`
	FinancialHealth      types.InstitutionFinancials `json:"financial_health"`
	AcademicRanking      types.AcademicRanking       `json:"academic_ranking"`
	LoanPerformance      types.LoanPerformance       `json:"loan_performance"`
	SpecialPrograms      []types.SpecialProgram      `json:"special_programs"`
	DiscountStructure    types.DiscountStructure     `json:"discount_structure"`
	DisbursementTerms    types.DisbursementTerms     `json:"disbursement_terms"`
	PartnershipAgreement types.PartnershipAgreement  `json:"partnership_agreement"`
	Status               string                      `json:"status"`
	CreatedAt            time.Time                   `json:"created_at"`
	UpdatedAt            time.Time                   `json:"updated_at"`
}

// VerifyInstitution verifies institution eligibility and partnership status
func (im *InstitutionManager) VerifyInstitution(ctx sdk.Context, institutionID string) (types.InstitutionStatus, error) {
	// Get institution details
	institution, found := im.keeper.GetPartnerInstitution(ctx, institutionID)
	if !found {
		// Check if institution is in general database
		institutionInfo, found := im.keeper.GetInstitutionInfo(ctx, institutionID)
		if !found {
			return types.InstitutionStatus{}, fmt.Errorf("institution not found: %s", institutionID)
		}

		// Perform basic verification for non-partner institutions
		return im.performBasicInstitutionVerification(ctx, institutionInfo)
	}

	// Enhanced verification for partner institutions
	return im.performPartnerInstitutionVerification(ctx, institution)
}

// OnboardPartnerInstitution onboards a new educational institution as partner
func (im *InstitutionManager) OnboardPartnerInstitution(ctx sdk.Context, request types.PartnershipRequest) (*PartnerInstitution, error) {
	// Validate institution eligibility
	eligible, reason := im.checkInstitutionEligibility(ctx, request)
	if !eligible {
		return nil, fmt.Errorf("institution not eligible for partnership: %s", reason)
	}

	// Perform due diligence
	dueDiligence, err := im.performInstitutionDueDiligence(ctx, request.InstitutionID)
	if err != nil {
		return nil, fmt.Errorf("due diligence failed: %w", err)
	}

	// Determine partnership level
	partnershipLevel := im.determinePartnershipLevel(ctx, dueDiligence)

	// Create partnership agreement
	agreement := im.createPartnershipAgreement(ctx, request, partnershipLevel)

	// Generate institution ID if new
	institutionID := request.InstitutionID
	if institutionID == "" {
		institutionID = im.generateInstitutionID(ctx, request.InstitutionName)
	}

	// Create partner institution
	partner := &PartnerInstitution{
		InstitutionID:   institutionID,
		InstitutionInfo: request.InstitutionInfo,
		PartnershipLevel: partnershipLevel,
		AccreditationStatus: dueDiligence.AccreditationStatus,
		PlacementRecord:     dueDiligence.PlacementRecord,
		FinancialHealth:     dueDiligence.FinancialHealth,
		AcademicRanking:     dueDiligence.AcademicRanking,
		LoanPerformance: types.LoanPerformance{
			InstitutionID:        institutionID,
			TotalLoansProcessed:  0,
			DefaultRate:         sdk.ZeroDec(),
			AverageRepaymentTime: 0,
			PlacementSuccess:    sdk.ZeroDec(),
		},
		SpecialPrograms:   im.createSpecialPrograms(ctx, partnershipLevel, request.CourseCategories),
		DiscountStructure: im.createDiscountStructure(ctx, partnershipLevel),
		DisbursementTerms: im.createDisbursementTerms(ctx, partnershipLevel),
		PartnershipAgreement: agreement,
		Status:    "ACTIVE",
		CreatedAt: ctx.BlockTime(),
		UpdatedAt: ctx.BlockTime(),
	}

	// Store partner institution
	im.keeper.SetPartnerInstitution(ctx, *partner)

	// Emit partnership event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeInstitutionPartnered,
			sdk.NewAttribute(types.AttributeKeyInstitutionID, institutionID),
			sdk.NewAttribute(types.AttributeKeyInstitutionName, request.InstitutionName),
			sdk.NewAttribute(types.AttributeKeyPartnershipLevel, partnershipLevel),
		),
	)

	return partner, nil
}

// UpdatePartnershipLevel updates institution partnership level based on performance
func (im *InstitutionManager) UpdatePartnershipLevel(ctx sdk.Context, institutionID string) error {
	// Get current partner institution
	partner, found := im.keeper.GetPartnerInstitution(ctx, institutionID)
	if !found {
		return fmt.Errorf("partner institution not found: %s", institutionID)
	}

	// Evaluate current performance
	performance := im.evaluateInstitutionPerformance(ctx, institutionID)

	// Determine new partnership level
	newLevel := im.determinePartnershipLevelFromPerformance(performance)

	// Update if level changed
	if newLevel != partner.PartnershipLevel {
		oldLevel := partner.PartnershipLevel
		partner.PartnershipLevel = newLevel
		partner.DiscountStructure = im.createDiscountStructure(ctx, newLevel)
		partner.DisbursementTerms = im.createDisbursementTerms(ctx, newLevel)
		partner.UpdatedAt = ctx.BlockTime()

		// Store updated partner
		im.keeper.SetPartnerInstitution(ctx, partner)

		// Emit level change event
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypePartnershipLevelUpdated,
				sdk.NewAttribute(types.AttributeKeyInstitutionID, institutionID),
				sdk.NewAttribute(types.AttributeKeyOldLevel, oldLevel),
				sdk.NewAttribute(types.AttributeKeyNewLevel, newLevel),
			),
		)
	}

	return nil
}

// ProcessInstitutionPayment handles direct payments to educational institutions
func (im *InstitutionManager) ProcessInstitutionPayment(ctx sdk.Context, institutionID string, amount sdk.Coin, purpose string) error {
	// Get institution payment details
	partner, found := im.keeper.GetPartnerInstitution(ctx, institutionID)
	if !found {
		return fmt.Errorf("partner institution not found for payment: %s", institutionID)
	}

	// Get institution payment address
	institutionAddr, err := sdk.AccAddressFromBech32(partner.InstitutionInfo.PaymentAddress)
	if err != nil {
		return fmt.Errorf("invalid institution payment address: %s", partner.InstitutionInfo.PaymentAddress)
	}

	// Calculate payment processing fee for institution
	processingFee := im.calculateInstitutionProcessingFee(ctx, amount, partner.PartnershipLevel)
	netPayment := amount.Sub(processingFee)

	// Transfer from ShikshaMitra module to institution
	err = im.keeper.bankKeeper.SendCoinsFromModuleToAccount(
		ctx,
		types.ModuleName,
		institutionAddr,
		sdk.NewCoins(netPayment),
	)
	if err != nil {
		return fmt.Errorf("failed to transfer payment to institution: %w", err)
	}

	// Record payment
	payment := types.InstitutionPayment{
		PaymentID:       im.generatePaymentID(ctx, institutionID),
		InstitutionID:   institutionID,
		Amount:          amount,
		NetAmount:       netPayment,
		ProcessingFee:   processingFee,
		Purpose:         purpose,
		PaymentDate:     ctx.BlockTime(),
		Status:          "COMPLETED",
	}

	im.keeper.SetInstitutionPayment(ctx, payment)

	// Update institution payment statistics
	im.updateInstitutionPaymentStats(ctx, institutionID, payment)

	// Emit payment event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeInstitutionPaymentProcessed,
			sdk.NewAttribute(types.AttributeKeyInstitutionID, institutionID),
			sdk.NewAttribute(types.AttributeKeyPaymentAmount, netPayment.String()),
			sdk.NewAttribute(types.AttributeKeyPaymentPurpose, purpose),
		),
	)

	return nil
}

// Helper functions for institution management

func (im *InstitutionManager) performBasicInstitutionVerification(ctx sdk.Context, info types.InstitutionInfo) (types.InstitutionStatus, error) {
	status := types.InstitutionStatus{
		InstitutionID:   info.InstitutionID,
		InstitutionInfo: info,
		IsAccredited:    false,
		PartnershipLevel: "NONE",
		IsEligible:      false,
	}

	// Check basic accreditation
	accreditation := im.keeper.GetInstitutionAccreditation(ctx, info.InstitutionID)
	if accreditation != nil && accreditation.IsValid {
		status.IsAccredited = true
		status.AccreditationLevel = accreditation.Level
		status.AccreditingBody = accreditation.AccreditingBody
	}

	// Check government recognition
	if im.isGovernmentRecognized(ctx, info.InstitutionID) {
		status.GovernmentRecognized = true
	}

	// Basic eligibility check
	if status.IsAccredited && status.GovernmentRecognized {
		status.IsEligible = true
		status.EligibilityReason = "Accredited and government recognized"
	} else {
		status.EligibilityReason = "Not accredited or government recognized"
	}

	return status, nil
}

func (im *InstitutionManager) performPartnerInstitutionVerification(ctx sdk.Context, partner PartnerInstitution) (types.InstitutionStatus, error) {
	status := types.InstitutionStatus{
		InstitutionID:        partner.InstitutionID,
		InstitutionInfo:      partner.InstitutionInfo,
		IsAccredited:         partner.AccreditationStatus.IsAccredited,
		AccreditationLevel:   partner.AccreditationStatus.Level,
		AccreditingBody:      partner.AccreditationStatus.AccreditingBody,
		GovernmentRecognized: partner.AccreditationStatus.GovernmentRecognized,
		PartnershipLevel:     partner.PartnershipLevel,
		IsEligible:          true,
		EligibilityReason:   "Verified partner institution",
		PartnerBenefits:     im.getPartnerBenefits(partner.PartnershipLevel),
		PlacementRate:       partner.PlacementRecord.PlacementRate,
		AveragePackage:      partner.PlacementRecord.AveragePackage,
	}

	return status, nil
}

func (im *InstitutionManager) checkInstitutionEligibility(ctx sdk.Context, request types.PartnershipRequest) (bool, string) {
	// Check accreditation
	accreditation := im.keeper.GetInstitutionAccreditation(ctx, request.InstitutionID)
	if accreditation == nil || !accreditation.IsValid {
		return false, "Institution not properly accredited"
	}

	// Check government recognition
	if !im.isGovernmentRecognized(ctx, request.InstitutionID) {
		return false, "Institution not government recognized"
	}

	// Check minimum years of operation
	if request.YearsOfOperation < 5 {
		return false, "Institution must have at least 5 years of operation"
	}

	// Check minimum student strength
	if request.StudentStrength < 500 {
		return false, "Institution must have at least 500 students"
	}

	// Check placement rate
	placementData := im.keeper.GetInstitutionPlacementData(ctx, request.InstitutionID)
	if placementData != nil && placementData.PlacementRate.LT(sdk.NewDecWithPrec(6, 1)) { // < 60%
		return false, "Institution placement rate below minimum threshold"
	}

	return true, ""
}

func (im *InstitutionManager) performInstitutionDueDiligence(ctx sdk.Context, institutionID string) (types.InstitutionDueDiligence, error) {
	dueDiligence := types.InstitutionDueDiligence{
		InstitutionID: institutionID,
	}

	// Accreditation verification
	accreditation := im.keeper.GetInstitutionAccreditation(ctx, institutionID)
	if accreditation != nil {
		dueDiligence.AccreditationStatus = types.AccreditationStatus{
			IsAccredited:         accreditation.IsValid,
			Level:               accreditation.Level,
			AccreditingBody:     accreditation.AccreditingBody,
			ExpiryDate:          accreditation.ExpiryDate,
			GovernmentRecognized: im.isGovernmentRecognized(ctx, institutionID),
		}
	}

	// Placement record verification
	placementData := im.keeper.GetInstitutionPlacementData(ctx, institutionID)
	if placementData != nil {
		dueDiligence.PlacementRecord = types.PlacementRecord{
			InstitutionID:    institutionID,
			PlacementRate:    placementData.PlacementRate,
			AveragePackage:   placementData.AveragePackage,
			MedianPackage:    placementData.MedianPackage,
			TopRecruiters:    placementData.TopRecruiters,
			PlacementTrends:  placementData.YearlyTrends,
		}
	}

	// Financial health assessment
	financialData := im.keeper.GetInstitutionFinancialData(ctx, institutionID)
	if financialData != nil {
		dueDiligence.FinancialHealth = types.InstitutionFinancials{
			InstitutionID:    institutionID,
			AnnualRevenue:    financialData.Revenue,
			Profitability:    financialData.Profit.ToDec().Quo(financialData.Revenue.Amount.ToDec()),
			DebtLevels:       financialData.TotalDebt,
			LiquidityRatio:   financialData.CurrentAssets.Amount.ToDec().Quo(financialData.CurrentLiabilities.Amount.ToDec()),
			FinancialRating:  im.calculateFinancialRating(financialData),
		}
	}

	// Academic ranking
	ranking := im.keeper.GetInstitutionRanking(ctx, institutionID)
	if ranking != nil {
		dueDiligence.AcademicRanking = types.AcademicRanking{
			InstitutionID:     institutionID,
			NIRFRanking:      ranking.NIRFRank,
			CategoryRanking:   ranking.CategoryRank,
			GlobalRanking:     ranking.GlobalRank,
			ResearchRating:    ranking.ResearchScore,
			TeachingRating:    ranking.TeachingScore,
			IndustryRating:    ranking.IndustryScore,
		}
	}

	return dueDiligence, nil
}

func (im *InstitutionManager) determinePartnershipLevel(ctx sdk.Context, dueDiligence types.InstitutionDueDiligence) string {
	score := 0

	// Academic ranking score
	if dueDiligence.AcademicRanking.NIRFRanking > 0 && dueDiligence.AcademicRanking.NIRFRanking <= 50 {
		score += 30 // Top 50 NIRF
	} else if dueDiligence.AcademicRanking.NIRFRanking <= 100 {
		score += 20 // Top 100 NIRF
	} else if dueDiligence.AcademicRanking.NIRFRanking <= 200 {
		score += 10 // Top 200 NIRF
	}

	// Placement rate score
	if dueDiligence.PlacementRecord.PlacementRate.GTE(sdk.NewDecWithPrec(9, 1)) { // >= 90%
		score += 25
	} else if dueDiligence.PlacementRecord.PlacementRate.GTE(sdk.NewDecWithPrec(8, 1)) { // >= 80%
		score += 20
	} else if dueDiligence.PlacementRecord.PlacementRate.GTE(sdk.NewDecWithPrec(7, 1)) { // >= 70%
		score += 15
	}

	// Financial health score
	if dueDiligence.FinancialHealth.FinancialRating == "AAA" || dueDiligence.FinancialHealth.FinancialRating == "AA" {
		score += 20
	} else if dueDiligence.FinancialHealth.FinancialRating == "A" {
		score += 15
	} else if dueDiligence.FinancialHealth.FinancialRating == "BBB" {
		score += 10
	}

	// Accreditation score
	if dueDiligence.AccreditationStatus.Level == "A++" {
		score += 15
	} else if dueDiligence.AccreditationStatus.Level == "A+" {
		score += 12
	} else if dueDiligence.AccreditationStatus.Level == "A" {
		score += 10
	}

	// Research rating score
	if dueDiligence.AcademicRanking.ResearchRating.GTE(sdk.NewDec(9)) {
		score += 10
	} else if dueDiligence.AcademicRanking.ResearchRating.GTE(sdk.NewDec(8)) {
		score += 8
	}

	// Determine partnership level based on total score
	if score >= 80 {
		return "GOLD"
	} else if score >= 60 {
		return "SILVER"
	} else if score >= 40 {
		return "BRONZE"
	} else {
		return "STANDARD"
	}
}

func (im *InstitutionManager) createDiscountStructure(ctx sdk.Context, partnershipLevel string) types.DiscountStructure {
	structure := types.DiscountStructure{
		PartnershipLevel: partnershipLevel,
	}

	switch partnershipLevel {
	case "GOLD":
		structure.InterestDiscount = sdk.NewDecWithPrec(2, 2)     // 2% discount
		structure.ProcessingFeeDiscount = sdk.NewDecWithPrec(5, 1) // 50% discount
		structure.SpecialOffers = []string{"Zero processing fee for merit students", "Additional scholarship programs"}
	case "SILVER":
		structure.InterestDiscount = sdk.NewDecWithPrec(15, 3)    // 1.5% discount
		structure.ProcessingFeeDiscount = sdk.NewDecWithPrec(3, 1) // 30% discount
		structure.SpecialOffers = []string{"Reduced processing fee for merit students"}
	case "BRONZE":
		structure.InterestDiscount = sdk.NewDecWithPrec(1, 2)     // 1% discount
		structure.ProcessingFeeDiscount = sdk.NewDecWithPrec(2, 1) // 20% discount
		structure.SpecialOffers = []string{"Merit-based fee reductions"}
	default: // STANDARD
		structure.InterestDiscount = sdk.NewDecWithPrec(5, 3)     // 0.5% discount
		structure.ProcessingFeeDiscount = sdk.NewDecWithPrec(1, 1) // 10% discount
		structure.SpecialOffers = []string{"Standard partnership benefits"}
	}

	return structure
}

func (im *InstitutionManager) createDisbursementTerms(ctx sdk.Context, partnershipLevel string) types.DisbursementTerms {
	terms := types.DisbursementTerms{
		PartnershipLevel: partnershipLevel,
	}

	switch partnershipLevel {
	case "GOLD":
		terms.DirectInstitutionPayment = true
		terms.PaymentSchedule = "SEMESTER_WISE"
		terms.PaymentTiming = "ADVANCE"
		terms.DocumentationRequired = []string{"Enrollment confirmation", "Fee structure"}
	case "SILVER":
		terms.DirectInstitutionPayment = true
		terms.PaymentSchedule = "SEMESTER_WISE"
		terms.PaymentTiming = "ON_ENROLLMENT"
		terms.DocumentationRequired = []string{"Enrollment confirmation", "Fee structure", "Academic progress"}
	case "BRONZE":
		terms.DirectInstitutionPayment = true
		terms.PaymentSchedule = "QUARTERLY"
		terms.PaymentTiming = "ON_ENROLLMENT"
		terms.DocumentationRequired = []string{"Enrollment confirmation", "Fee structure", "Academic progress", "Attendance verification"}
	default: // STANDARD
		terms.DirectInstitutionPayment = false
		terms.PaymentSchedule = "MANUAL"
		terms.PaymentTiming = "ON_REQUEST"
		terms.DocumentationRequired = []string{"Enrollment confirmation", "Fee structure", "Academic progress", "Attendance verification", "Institution verification"}
	}

	return terms
}

// Additional utility functions
func (im *InstitutionManager) generateInstitutionID(ctx sdk.Context, institutionName string) string {
	timestamp := ctx.BlockTime().Unix()
	return fmt.Sprintf("INST-%s-%d", institutionName[:4], timestamp)
}

func (im *InstitutionManager) generatePaymentID(ctx sdk.Context, institutionID string) string {
	timestamp := ctx.BlockTime().Unix()
	return fmt.Sprintf("PAY-%s-%d", institutionID[:8], timestamp)
}

func (im *InstitutionManager) calculateInstitutionProcessingFee(ctx sdk.Context, amount sdk.Coin, partnershipLevel string) sdk.Coin {
	params := im.keeper.GetParams(ctx)
	baseFeeRate := params.InstitutionProcessingFeeRate

	// Apply partnership discount
	discountStructure := im.createDiscountStructure(ctx, partnershipLevel)
	discountedRate := baseFeeRate.Mul(sdk.OneDec().Sub(discountStructure.ProcessingFeeDiscount))

	feeAmount := amount.Amount.ToDec().Mul(discountedRate).TruncateInt()
	return sdk.NewCoin(amount.Denom, feeAmount)
}

func (im *InstitutionManager) isGovernmentRecognized(ctx sdk.Context, institutionID string) bool {
	// Check government recognition database
	recognition := im.keeper.GetGovernmentRecognition(ctx, institutionID)
	return recognition != nil && recognition.IsRecognized
}

func (im *InstitutionManager) calculateFinancialRating(data types.InstitutionFinancialData) string {
	// Simplified financial rating calculation
	score := 0

	// Profitability
	profitMargin := data.Profit.ToDec().Quo(data.Revenue.Amount.ToDec())
	if profitMargin.GTE(sdk.NewDecWithPrec(2, 1)) { // >= 20%
		score += 25
	} else if profitMargin.GTE(sdk.NewDecWithPrec(15, 2)) { // >= 15%
		score += 20
	} else if profitMargin.GTE(sdk.NewDecWithPrec(1, 1)) { // >= 10%
		score += 15
	}

	// Liquidity
	liquidityRatio := data.CurrentAssets.Amount.ToDec().Quo(data.CurrentLiabilities.Amount.ToDec())
	if liquidityRatio.GTE(sdk.NewDecWithPrec(15, 1)) { // >= 1.5
		score += 25
	} else if liquidityRatio.GTE(sdk.OneDec()) { // >= 1.0
		score += 20
	}

	// Debt levels
	debtRatio := data.TotalDebt.Amount.ToDec().Quo(data.Revenue.Amount.ToDec())
	if debtRatio.LTE(sdk.NewDecWithPrec(2, 1)) { // <= 20%
		score += 25
	} else if debtRatio.LTE(sdk.NewDecWithPrec(4, 1)) { // <= 40%
		score += 20
	}

	// Revenue growth
	if data.RevenueGrowth.GTE(sdk.NewDecWithPrec(1, 1)) { // >= 10%
		score += 25
	} else if data.RevenueGrowth.GTE(sdk.NewDecWithPrec(5, 2)) { // >= 5%
		score += 20
	}

	// Assign rating based on score
	if score >= 90 {
		return "AAA"
	} else if score >= 80 {
		return "AA"
	} else if score >= 70 {
		return "A"
	} else if score >= 60 {
		return "BBB"
	} else if score >= 50 {
		return "BB"
	} else {
		return "B"
	}
}

// Additional helper functions would include:
// - evaluateInstitutionPerformance
// - determinePartnershipLevelFromPerformance
// - createPartnershipAgreement
// - createSpecialPrograms
// - getPartnerBenefits
// - updateInstitutionPaymentStats