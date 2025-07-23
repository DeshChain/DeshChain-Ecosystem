package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/deshchain/x/krishimitra/types"
)

// Keeper of the krishimitra store
type Keeper struct {
	cdc               codec.BinaryCodec
	storeKey          sdk.StoreKey
	memKey            sdk.StoreKey
	paramspace        types.ParamSubspace
	bankKeeper        types.BankKeeper
	accountKeeper     types.AccountKeeper
	dhanpataKeeper    types.DhanPataKeeper
	liquidityKeeper   types.LiquidityManagerKeeper // REVOLUTIONARY: Member-only lending verification
}

// NewKeeper creates a new krishimitra Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	ps types.ParamSubspace,
	bankKeeper types.BankKeeper,
	accountKeeper types.AccountKeeper,
	dhanpataKeeper types.DhanPataKeeper,
	liquidityKeeper types.LiquidityManagerKeeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:             cdc,
		storeKey:        storeKey,
		memKey:          memKey,
		paramspace:      ps,
		bankKeeper:      bankKeeper,
		accountKeeper:   accountKeeper,
		dhanpataKeeper:  dhanpataKeeper,
		liquidityKeeper: liquidityKeeper,
	}
}

// Logger returns a module-specific logger
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// SetLoan stores a loan in the store
func (k Keeper) SetLoan(ctx sdk.Context, loan types.Loan) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&loan)
	store.Set(types.GetLoanKey(loan.ID), bz)
}

// GetLoan retrieves a loan from the store
func (k Keeper) GetLoan(ctx sdk.Context, loanID string) (loan types.Loan, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetLoanKey(loanID))
	if bz == nil {
		return loan, false
	}
	k.cdc.MustUnmarshal(bz, &loan)
	return loan, true
}

// SetLoanApplication stores a loan application
func (k Keeper) SetLoanApplication(ctx sdk.Context, application types.LoanApplication) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&application)
	store.Set(types.GetApplicationKey(application.ID), bz)
}

// GetLoanApplication retrieves a loan application
func (k Keeper) GetLoanApplication(ctx sdk.Context, applicationID string) (application types.LoanApplication, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetApplicationKey(applicationID))
	if bz == nil {
		return application, false
	}
	k.cdc.MustUnmarshal(bz, &application)
	return application, true
}

// SetPINCodeEligibility stores PIN code eligibility criteria
func (k Keeper) SetPINCodeEligibility(ctx sdk.Context, eligibility types.PINCodeEligibility) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&eligibility)
	store.Set(types.GetPINCodeKey(eligibility.PINCode), bz)
}

// GetPINCodeEligibility retrieves PIN code eligibility
func (k Keeper) GetPINCodeEligibility(ctx sdk.Context, pinCode string) (eligibility types.PINCodeEligibility, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetPINCodeKey(pinCode))
	if bz == nil {
		return eligibility, false
	}
	k.cdc.MustUnmarshal(bz, &eligibility)
	return eligibility, true
}

// CheckEligibility verifies if an applicant is eligible for a loan
func (k Keeper) CheckEligibility(ctx sdk.Context, application types.LoanApplication) (bool, string) {
	// Check PIN code eligibility
	pinEligibility, found := k.GetPINCodeEligibility(ctx, application.PINCode)
	if !found || !pinEligibility.IsEligible {
		return false, "PIN code not eligible for agricultural loans"
	}

	// Check loan amount limits
	if application.RequestedAmount.Amount.GT(pinEligibility.MaxLoanAmount.Amount) {
		return false, fmt.Sprintf("Requested amount exceeds maximum limit for PIN code %s", application.PINCode)
	}

	// Check credit score
	if application.CreditScore < 500 {
		return false, "Credit score too low for loan approval"
	}

	// Check DhanPata address verification
	isDhanPataVerified := k.dhanpataKeeper.IsAddressVerified(ctx, application.DhanPataAddress)
	if !isDhanPataVerified {
		return false, "DhanPata address not verified"
	}

	// Check land area minimum (at least 0.5 acres)
	minLandArea := sdk.NewDecWithPrec(5, 1) // 0.5
	if application.LandArea.LT(minLandArea) {
		return false, "Land area too small for agricultural loan"
	}

	return true, ""
}

// CalculateInterestRate calculates the interest rate based on various factors
func (k Keeper) CalculateInterestRate(ctx sdk.Context, application types.LoanApplication) sdk.Dec {
	// Get base rate from PIN code
	pinEligibility, found := k.GetPINCodeEligibility(ctx, application.PINCode)
	if !found {
		// Default to max rate if PIN code not found
		maxRate, _ := sdk.NewDecFromStr(types.MaxInterestRate)
		return maxRate
	}

	baseRate := pinEligibility.BaseInterestRate

	// Adjust based on credit score
	creditScoreAdjustment := sdk.ZeroDec()
	if application.CreditScore >= 750 {
		creditScoreAdjustment = sdk.NewDecWithPrec(-5, 3) // -0.5%
	} else if application.CreditScore >= 650 {
		creditScoreAdjustment = sdk.NewDecWithPrec(-25, 4) // -0.25%
	} else if application.CreditScore < 550 {
		creditScoreAdjustment = sdk.NewDecWithPrec(5, 3) // +0.5%
	}

	// Priority district bonus
	priorityBonus := sdk.ZeroDec()
	if pinEligibility.PriorityDistrict {
		priorityBonus = sdk.NewDecWithPrec(-1, 2) // -1%
	}

	// Calculate final rate
	finalRate := baseRate.Add(creditScoreAdjustment).Add(priorityBonus)

	// Ensure within bounds
	minRate, _ := sdk.NewDecFromStr(types.MinInterestRate)
	maxRate, _ := sdk.NewDecFromStr(types.MaxInterestRate)

	if finalRate.LT(minRate) {
		return minRate
	}
	if finalRate.GT(maxRate) {
		return maxRate
	}

	return finalRate
}
// Additional helper methods for gRPC queries

// GetWeatherData retrieves weather data for a pincode
func (k Keeper) GetWeatherData(ctx sdk.Context, pincode string) (types.WeatherInfo, bool) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(append([]byte("weather:"), []byte(pincode)...))
	if b == nil {
		return types.WeatherInfo{}, false
	}

	var weather types.WeatherInfo
	k.cdc.MustUnmarshal(b, &weather)
	return weather, true
}

// GetCropRecommendations returns crop recommendations based on weather
func (k Keeper) GetCropRecommendations(ctx sdk.Context, pincode string, weather types.WeatherInfo) []string {
	// Simple recommendations based on weather
	recommendations := []string{}
	
	// Parse temperature
	if weather.Temperature \!= "" {
		// Extract numeric value (assuming format like "28°C")
		if len(weather.Temperature) > 2 {
			tempStr := weather.Temperature[:2]
			if temp, err := sdk.NewDecFromStr(tempStr); err == nil {
				if temp.GT(sdk.NewDec(25)) && temp.LT(sdk.NewDec(35)) {
					recommendations = append(recommendations, "Rice", "Cotton", "Sugarcane")
				} else if temp.LT(sdk.NewDec(25)) {
					recommendations = append(recommendations, "Wheat", "Mustard", "Potato")
				}
			}
		}
	}
	
	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Rice", "Wheat", "Pulses") // Default crops
	}
	
	return recommendations
}

// CalculateLoanStatistics calculates loan statistics for a given period
func (k Keeper) CalculateLoanStatistics(ctx sdk.Context, period string) types.LoanStatistics {
	store := ctx.KVStore(k.storeKey)
	loanStore := prefix.NewStore(store, types.LoanKeyPrefix)
	
	stats := types.LoanStatistics{
		LoansByCrop: make(map[string]int64),
		LoansByPincode: make(map[string]int64),
	}
	
	iterator := loanStore.Iterator(nil, nil)
	defer iterator.Close()
	
	totalAmount := sdk.ZeroInt()
	totalRepaid := sdk.ZeroInt()
	totalInterest := sdk.ZeroDec()
	defaultedCount := int64(0)
	successfulHarvests := int64(0)
	
	for ; iterator.Valid(); iterator.Next() {
		var loan types.AgriculturalLoan
		if err := k.cdc.Unmarshal(iterator.Value(), &loan); err \!= nil {
			continue
		}
		
		stats.TotalLoans++
		totalAmount = totalAmount.Add(loan.LoanAmount.Amount)
		totalRepaid = totalRepaid.Add(loan.RepaidAmount.Amount)
		
		rate, _ := sdk.NewDecFromStr(loan.InterestRate)
		totalInterest = totalInterest.Add(rate)
		
		if loan.Status == types.LoanStatus_LOAN_STATUS_ACTIVE {
			stats.ActiveLoans++
		} else if loan.Status == types.LoanStatus_LOAN_STATUS_DEFAULTED {
			defaultedCount++
		}
		
		// Count by crop type
		cropTypeStr := loan.CropType.String()
		stats.LoansByCrop[cropTypeStr]++
		
		// Count by pincode
		stats.LoansByPincode[loan.Pincode]++
		
		// Count successful harvests
		if loan.HarvestStatus == "SUCCESSFUL" {
			successfulHarvests++
		}
	}
	
	stats.TotalDisbursed = totalAmount.String()
	stats.TotalRepaid = totalRepaid.String()
	stats.DefaultedLoans = defaultedCount
	stats.SuccessfulHarvests = successfulHarvests
	
	if stats.TotalLoans > 0 {
		stats.AverageLoanAmount = totalAmount.Quo(sdk.NewInt(stats.TotalLoans)).String()
		stats.AverageInterestRate = totalInterest.Quo(sdk.NewDec(stats.TotalLoans)).String()
		stats.DefaultRate = sdk.NewDec(defaultedCount).Mul(sdk.NewDec(100)).Quo(sdk.NewDec(stats.TotalLoans)).String() + "%"
	} else {
		stats.AverageLoanAmount = "0"
		stats.AverageInterestRate = "0"
		stats.DefaultRate = "0%"
	}
	
	return stats
}

// GetActiveFestivalOffers returns currently active festival offers
func (k Keeper) GetActiveFestivalOffers(ctx sdk.Context) []types.FestivalOffer {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.FestivalOfferPrefix)
	defer iterator.Close()

	var offers []types.FestivalOffer
	currentTime := ctx.BlockTime()
	
	for ; iterator.Valid(); iterator.Next() {
		var offer types.FestivalOffer
		k.cdc.MustUnmarshal(iterator.Value(), &offer)

		if currentTime.After(offer.StartDate) && currentTime.Before(offer.EndDate) {
			offers = append(offers, offer)
		}
	}
	
	return offers
}

// GetNextFestivalOffer returns the next upcoming festival offer
func (k Keeper) GetNextFestivalOffer(ctx sdk.Context) *types.FestivalOffer {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.FestivalOfferPrefix)
	defer iterator.Close()

	var nextOffer *types.FestivalOffer
	currentTime := ctx.BlockTime()
	
	for ; iterator.Valid(); iterator.Next() {
		var offer types.FestivalOffer
		k.cdc.MustUnmarshal(iterator.Value(), &offer)

		if currentTime.Before(offer.StartDate) {
			if nextOffer == nil || offer.StartDate.Before(nextOffer.StartDate) {
				nextOffer = &offer
			}
		}
	}
	
	return nextOffer
}

// CheckFarmerEligibility checks if a farmer is eligible for a loan
func (k Keeper) CheckFarmerEligibility(ctx sdk.Context, farmer string, amount sdk.Coin, cropType string, landSize string) (bool, sdk.Coin, []string) {
	reasons := []string{}
	params := k.GetParams(ctx)
	
	// Check minimum and maximum loan amounts
	if amount.IsLT(params.MinLoanAmount) {
		reasons = append(reasons, fmt.Sprintf("Loan amount below minimum %s", params.MinLoanAmount))
		return false, params.MaxLoanAmount, reasons
	}
	
	if amount.IsGT(params.MaxLoanAmount) {
		reasons = append(reasons, fmt.Sprintf("Loan amount exceeds maximum %s", params.MaxLoanAmount))
		return false, params.MaxLoanAmount, reasons
	}
	
	// Check farmer profile
	profile, found := k.GetFarmerProfile(ctx, farmer)
	if \!found {
		reasons = append(reasons, "Farmer profile not found. Please complete KYC")
		return false, params.MaxLoanAmount, reasons
	}
	
	// Check active loans
	if profile.ActiveLoans >= 2 {
		reasons = append(reasons, "Maximum 2 active loans allowed")
		return false, sdk.ZeroCoin(amount.Denom), reasons
	}
	
	// Check default history
	if profile.DefaultedLoans > 0 {
		reasons = append(reasons, "Previous loan defaults found")
		return false, sdk.ZeroCoin(amount.Denom), reasons
	}
	
	return true, params.MaxLoanAmount, []string{"Eligible for agricultural loan"}
}

// EstimateInterestRate estimates the interest rate for a farmer
func (k Keeper) EstimateInterestRate(ctx sdk.Context, farmer string, cropType string, amount sdk.Coin) sdk.Dec {
	params := k.GetParams(ctx)
	baseRate := sdk.MustNewDecFromStr(params.MaxInterestRate)
	
	// Get farmer profile
	profile, found := k.GetFarmerProfile(ctx, farmer)
	if \!found {
		return baseRate
	}
	
	// Apply discounts based on profile
	if profile.TotalLoans > 0 && profile.DefaultedLoans == 0 {
		// Good repayment history - 1% discount
		baseRate = baseRate.Sub(sdk.MustNewDecFromStr("0.01"))
	}
	
	// Small farmer discount
	if profile.LandSize.LT(sdk.NewDec(5)) {
		smallFarmerDiscount := sdk.MustNewDecFromStr(params.SmallFarmerDiscount)
		baseRate = baseRate.Sub(smallFarmerDiscount)
	}
	
	// Ensure rate is within bounds
	minRate := sdk.MustNewDecFromStr(params.MinInterestRate)
	if baseRate.LT(minRate) {
		baseRate = minRate
	}
	
	return baseRate
}

// GetApplicableDiscounts returns applicable discounts for a farmer
func (k Keeper) GetApplicableDiscounts(ctx sdk.Context, farmer string, cropType string) []string {
	discounts := []string{}
	params := k.GetParams(ctx)
	
	profile, found := k.GetFarmerProfile(ctx, farmer)
	if \!found {
		return discounts
	}
	
	if profile.TotalLoans > 0 && profile.DefaultedLoans == 0 {
		discounts = append(discounts, "Good repayment history: 1% discount")
	}
	
	if profile.LandSize.LT(sdk.NewDec(5)) {
		discounts = append(discounts, fmt.Sprintf("Small farmer: %s discount", params.SmallFarmerDiscount))
	}
	
	if profile.IsWomenFarmer {
		discounts = append(discounts, fmt.Sprintf("Women farmer: %s discount", params.WomenFarmerDiscount))
	}
	
	// Check for festival offers
	festivalOffers := k.GetActiveFestivalOffers(ctx)
	for _, offer := range festivalOffers {
		discounts = append(discounts, fmt.Sprintf("%s festival: %s discount", offer.FestivalName, offer.InterestReduction))
	}
	
	return discounts
}

// GetRequiredDocuments returns required documents for a crop type
func (k Keeper) GetRequiredDocuments(cropType string) []string {
	docs := []string{
		"Aadhaar Card",
		"PAN Card",
		"Land ownership documents",
		"Bank statements (last 6 months)",
		"Kisan Credit Card (if available)",
	}
	
	// Add crop-specific documents
	switch cropType {
	case "RICE", "WHEAT":
		docs = append(docs, "Previous harvest records")
	case "COTTON", "SUGARCANE":
		docs = append(docs, "Water availability certificate")
	case "HORTICULTURE":
		docs = append(docs, "Soil testing report")
	}
	
	return docs
}

// GetCreditScoringEngine returns a new credit scoring engine
func (k Keeper) GetCreditScoringEngine() *CreditScoringEngine {
	return NewCreditScoringEngine(k)
}

// GetInsuranceEngine returns a new insurance engine
func (k Keeper) GetInsuranceEngine() *InsuranceEngine {
	return NewInsuranceEngine(k)
}

// GetLoanProcessor returns a new loan processor
func (k Keeper) GetLoanProcessor() *LoanProcessor {
	return NewLoanProcessor(k)
}

// ProcessAgriculturalLoan processes a complete agricultural loan from application to disbursement
func (k Keeper) ProcessAgriculturalLoan(ctx sdk.Context, applicationID string) error {
	loanProcessor := k.GetLoanProcessor()
	
	// Process loan application
	loan, err := loanProcessor.ProcessLoanApplication(ctx, applicationID)
	if err != nil {
		return fmt.Errorf("failed to process loan application: %w", err)
	}
	
	// Auto-disburse if all conditions are met
	if loan.Status == types.LoanStatusApproved && !loan.CollateralRequired && 
	   (!loan.InsuranceRequired || loan.InsurancePolicyID != "") {
		err = loanProcessor.DisburseLoan(ctx, loan.LoanID)
		if err != nil {
			k.Logger(ctx).Error("Failed to auto-disburse loan", "loan_id", loan.LoanID, "error", err)
		}
	}
	
	return nil
}

// CreateComprehensiveCropInsurance creates crop insurance with weather derivatives
func (k Keeper) CreateComprehensiveCropInsurance(ctx sdk.Context, farmerID, cropType string, loanID string) (string, error) {
	insuranceEngine := k.GetInsuranceEngine()
	
	// Get farmer profile for insurance calculation
	farmerProfile, found := k.GetFarmerProfile(ctx, farmerID)
	if !found {
		return "", fmt.Errorf("farmer profile not found: %s", farmerID)
	}
	
	// Create insurance policy request
	request := &types.InsurancePolicyRequest{
		FarmerID:     farmerID,
		LoanID:       loanID,
		CropType:     cropType,
		CropVariety:  "STANDARD", // Could be enhanced to specify variety
		SeasonType:   k.determineCurrentSeason(ctx),
		SowingArea:   farmerProfile.TotalLandArea,
		CoverageType: "COMPREHENSIVE",
		SowingDate:   ctx.BlockTime(),
		HarvestDate:  k.calculateHarvestDate(ctx, cropType),
	}
	
	// Create insurance policy
	policy, err := insuranceEngine.CreateCropInsurancePolicy(ctx, request)
	if err != nil {
		return "", fmt.Errorf("failed to create insurance policy: %w", err)
	}
	
	return policy.PolicyID, nil
}

// Helper functions for KrishiMitra operations

func (k Keeper) determineCurrentSeason(ctx sdk.Context) string {
	currentTime := ctx.BlockTime()
	month := currentTime.Month()
	
	if month >= 6 && month <= 11 { // June to November
		return "KHARIF"
	} else if month >= 11 || month <= 3 { // November to March
		return "RABI"
	} else { // March to June
		return "ZAID"
	}
}

func (k Keeper) calculateHarvestDate(ctx sdk.Context, cropType string) time.Time {
	currentTime := ctx.BlockTime()
	
	// Crop duration map (in months)
	cropDurations := map[string]int{
		"RICE":       4,
		"WHEAT":      5,
		"COTTON":     6,
		"SUGARCANE":  12,
		"VEGETABLES": 3,
		"PULSES":     4,
	}
	
	if duration, found := cropDurations[cropType]; found {
		return currentTime.AddDate(0, duration, 0)
	}
	
	return currentTime.AddDate(0, 5, 0) // Default 5 months
}
