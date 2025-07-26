package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/DeshChain/DeshChain-Ecosystem/x/vyavasayamitra/types"
)

// Keeper of the vyavasayamitra store
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

// NewKeeper creates a new vyavasayamitra Keeper instance
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

// SetBusinessLoan stores a business loan in the store
func (k Keeper) SetBusinessLoan(ctx sdk.Context, loan types.BusinessLoan) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&loan)
	store.Set(types.GetLoanKey(loan.ID), bz)
}

// GetBusinessLoan retrieves a business loan from the store
func (k Keeper) GetBusinessLoan(ctx sdk.Context, loanID string) (loan types.BusinessLoan, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetLoanKey(loanID))
	if bz == nil {
		return loan, false
	}
	k.cdc.MustUnmarshal(bz, &loan)
	return loan, true
}

// SetBusinessProfile stores a business profile
func (k Keeper) SetBusinessProfile(ctx sdk.Context, profile types.BusinessProfile) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&profile)
	store.Set(types.GetBusinessProfileKey(profile.ID), bz)
}

// GetBusinessProfile retrieves a business profile
func (k Keeper) GetBusinessProfile(ctx sdk.Context, businessID string) (profile types.BusinessProfile, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetBusinessProfileKey(businessID))
	if bz == nil {
		return profile, false
	}
	k.cdc.MustUnmarshal(bz, &profile)
	return profile, true
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

// SetCreditLine stores a credit line
func (k Keeper) SetCreditLine(ctx sdk.Context, creditLine types.CreditLine) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&creditLine)
	store.Set(types.GetCreditLineKey(creditLine.ID), bz)
}

// GetCreditLine retrieves a credit line
func (k Keeper) GetCreditLine(ctx sdk.Context, creditLineID string) (creditLine types.CreditLine, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetCreditLineKey(creditLineID))
	if bz == nil {
		return creditLine, false
	}
	k.cdc.MustUnmarshal(bz, &creditLine)
	return creditLine, true
}

// SetMerchantRating stores merchant rating
func (k Keeper) SetMerchantRating(ctx sdk.Context, rating types.MerchantRating) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&rating)
	store.Set(types.GetMerchantRatingKey(rating.BusinessID), bz)
}

// GetMerchantRating retrieves merchant rating
func (k Keeper) GetMerchantRating(ctx sdk.Context, businessID string) (rating types.MerchantRating, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetMerchantRatingKey(businessID))
	if bz == nil {
		return rating, false
	}
	k.cdc.MustUnmarshal(bz, &rating)
	return rating, true
}

// CheckEligibility verifies if a business is eligible for a loan
func (k Keeper) CheckEligibility(ctx sdk.Context, application types.LoanApplication) (bool, string) {
	// Get business profile
	profile, found := k.GetBusinessProfile(ctx, application.BusinessID)
	if !found {
		return false, "Business profile not found"
	}

	// Check verification status
	if profile.VerificationStatus != "verified" {
		return false, "Business not verified"
	}

	// Check credit score
	if profile.CreditScore < 600 {
		return false, "Credit score too low for business loan"
	}

	// Check DhanPata address verification
	isDhanPataVerified := k.dhanpataKeeper.IsAddressVerified(ctx, application.DhanPataAddress)
	if !isDhanPataVerified {
		return false, "DhanPata address not verified"
	}

	// Check merchant rating
	rating, found := k.GetMerchantRating(ctx, application.BusinessID)
	if found && rating.OverallRating.LT(sdk.NewDecWithPrec(3, 0)) {
		return false, "Merchant rating too low"
	}

	// Check business age (minimum 1 year)
	oneYearAgo := ctx.BlockTime().AddDate(-1, 0, 0)
	if profile.EstablishedDate.After(oneYearAgo) {
		return false, "Business must be at least 1 year old"
	}

	// Check active loans limit
	if profile.ActiveLoans >= 3 {
		return false, "Maximum active loans limit reached"
	}

	return true, ""
}

// CalculateInterestRate calculates interest rate based on various factors
func (k Keeper) CalculateInterestRate(ctx sdk.Context, application types.LoanApplication) sdk.Dec {
	// Base rate
	baseRate, _ := sdk.NewDecFromStr("0.10") // 10%

	// Get business profile
	profile, found := k.GetBusinessProfile(ctx, application.BusinessID)
	if !found {
		maxRate, _ := sdk.NewDecFromStr(types.MaxInterestRate)
		return maxRate
	}

	// Credit score adjustment
	creditScoreAdjustment := sdk.ZeroDec()
	if profile.CreditScore >= 800 {
		creditScoreAdjustment = sdk.NewDecWithPrec(-2, 2) // -2%
	} else if profile.CreditScore >= 700 {
		creditScoreAdjustment = sdk.NewDecWithPrec(-1, 2) // -1%
	} else if profile.CreditScore < 650 {
		creditScoreAdjustment = sdk.NewDecWithPrec(1, 2) // +1%
	}

	// Business type adjustment
	businessTypeAdjustment := sdk.ZeroDec()
	switch profile.Category {
	case types.CategoryManufacturing:
		businessTypeAdjustment = sdk.NewDecWithPrec(-5, 3) // -0.5%
	case types.CategoryExportImport:
		businessTypeAdjustment = sdk.NewDecWithPrec(-1, 2) // -1%
	case types.CategoryTechnology:
		if profile.BusinessType == types.BusinessType_STARTUP {
			businessTypeAdjustment = sdk.NewDecWithPrec(1, 2) // +1%
		}
	}

	// Merchant rating adjustment
	rating, found := k.GetMerchantRating(ctx, profile.ID)
	if found {
		if rating.OverallRating.GTE(sdk.NewDecWithPrec(45, 1)) { // >= 4.5
			creditScoreAdjustment = creditScoreAdjustment.Add(sdk.NewDecWithPrec(-5, 3)) // -0.5%
		}
	}

	// Calculate final rate
	finalRate := baseRate.Add(creditScoreAdjustment).Add(businessTypeAdjustment)

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

// UpdateMerchantRating updates merchant rating based on repayment behavior
func (k Keeper) UpdateMerchantRating(ctx sdk.Context, businessID string, repayment types.Repayment) {
	rating, found := k.GetMerchantRating(ctx, businessID)
	if !found {
		// Initialize new rating
		rating = types.MerchantRating{
			BusinessID:        businessID,
			OverallRating:     sdk.NewDecWithPrec(35, 1), // 3.5
			PaymentScore:      sdk.NewDecWithPrec(35, 1), // 3.5
			BusinessScore:     sdk.NewDecWithPrec(35, 1), // 3.5
			FinancialScore:    sdk.NewDecWithPrec(35, 1), // 3.5
			TotalTransactions: 0,
			OnTimePayments:    0,
			DelayedPayments:   0,
			DefaultedPayments: 0,
		}
	}

	rating.TotalTransactions++
	
	// Update payment behavior
	loan, _ := k.GetBusinessLoan(ctx, repayment.LoanID)
	if loan.MaturityDate != nil && repayment.PaidAt.After(*loan.MaturityDate) {
		rating.DelayedPayments++
		// Decrease payment score
		rating.PaymentScore = rating.PaymentScore.Sub(sdk.NewDecWithPrec(1, 1))
	} else {
		rating.OnTimePayments++
		// Increase payment score
		rating.PaymentScore = rating.PaymentScore.Add(sdk.NewDecWithPrec(1, 1))
	}

	// Ensure scores are within 1-5 range
	if rating.PaymentScore.GT(sdk.NewDec(5)) {
		rating.PaymentScore = sdk.NewDec(5)
	}
	if rating.PaymentScore.LT(sdk.NewDec(1)) {
		rating.PaymentScore = sdk.NewDec(1)
	}

	// Update overall rating (average of all scores)
	rating.OverallRating = rating.PaymentScore.Add(rating.BusinessScore).Add(rating.FinancialScore).Quo(sdk.NewDec(3))
	rating.LastUpdated = ctx.BlockTime()

	k.SetMerchantRating(ctx, rating)
}
// Additional helper methods for gRPC queries

// CalculateLoanStatistics calculates loan statistics for a given period
func (k Keeper) CalculateLoanStatistics(ctx sdk.Context, period string) types.LoanStatistics {
	store := ctx.KVStore(k.storeKey)
	loanStore := prefix.NewStore(store, types.LoanKeyPrefix)
	
	stats := types.LoanStatistics{
		LoansByBusinessType: make(map[string]int64),
		LoansByPurpose: make(map[string]int64),
	}
	
	iterator := loanStore.Iterator(nil, nil)
	defer iterator.Close()
	
	totalAmount := sdk.ZeroInt()
	totalRepaid := sdk.ZeroInt()
	totalInterest := sdk.ZeroDec()
	totalCreditLines := sdk.ZeroInt()
	totalInvoiceFinancing := sdk.ZeroInt()
	totalBusinessAge := int64(0)
	businessCount := int64(0)
	
	for ; iterator.Valid(); iterator.Next() {
		var loan types.BusinessLoan
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
		}
		
		// Count by business type
		businessTypeStr := loan.BusinessType.String()
		stats.LoansByBusinessType[businessTypeStr]++
		
		// Count by purpose
		purposeStr := loan.LoanPurpose.String()
		stats.LoansByPurpose[purposeStr]++
		
		// Sum credit lines
		if loan.CreditLine \!= nil && loan.CreditLine.IsActive {
			totalCreditLines = totalCreditLines.Add(loan.CreditLine.CreditLimit.Amount)
		}
		
		// Sum invoice financing
		for _, inv := range loan.InvoiceFinancings {
			totalInvoiceFinancing = totalInvoiceFinancing.Add(inv.FinancedAmount.Amount)
		}
		
		// Sum business age
		totalBusinessAge += int64(loan.BusinessAge)
		businessCount++
	}
	
	stats.TotalDisbursed = totalAmount.String()
	stats.TotalRepaid = totalRepaid.String()
	stats.TotalCreditLines = totalCreditLines.String()
	stats.InvoiceFinancingVolume = totalInvoiceFinancing.String()
	
	if stats.TotalLoans > 0 {
		stats.AverageLoanAmount = totalAmount.Quo(sdk.NewInt(stats.TotalLoans)).String()
		stats.AverageInterestRate = totalInterest.Quo(sdk.NewDec(stats.TotalLoans)).String()
		defaultRate := k.CalculateDefaultRate(ctx)
		stats.DefaultRate = defaultRate
	} else {
		stats.AverageLoanAmount = "0"
		stats.AverageInterestRate = "0"
		stats.DefaultRate = "0%"
	}
	
	if businessCount > 0 {
		stats.AverageBusinessAge = fmt.Sprintf("%d months", totalBusinessAge/businessCount)
	} else {
		stats.AverageBusinessAge = "0 months"
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

// CheckBusinessEligibility checks if a business is eligible for a loan
func (k Keeper) CheckBusinessEligibility(ctx sdk.Context, business string, amount sdk.Coin, businessType string, annualRevenue sdk.Coin) (bool, sdk.Coin, []string) {
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
	
	// Check business profile
	profile, found := k.GetBusinessProfile(ctx, business)
	if \!found {
		reasons = append(reasons, "Business profile not found. Please complete KYC")
		return false, params.MaxLoanAmount, reasons
	}
	
	// Check active loans
	if profile.ActiveLoans >= 3 {
		reasons = append(reasons, "Maximum 3 active loans allowed")
		return false, sdk.ZeroCoin(amount.Denom), reasons
	}
	
	// Check credit score
	if profile.CreditScore < 650 {
		reasons = append(reasons, "Credit score too low (minimum 650 required)")
		return false, sdk.ZeroCoin(amount.Denom), reasons
	}
	
	// Check loan to revenue ratio
	maxLoanAmount := annualRevenue.Amount.Mul(sdk.NewInt(3)).Quo(sdk.NewInt(10)) // 30% of annual revenue
	if amount.Amount.GT(maxLoanAmount) {
		reasons = append(reasons, "Loan amount exceeds 30% of annual revenue")
		return false, sdk.NewCoin(amount.Denom, maxLoanAmount), reasons
	}
	
	return true, sdk.NewCoin(amount.Denom, maxLoanAmount), []string{"Eligible for business loan"}
}

// EstimateInterestRate estimates the interest rate for a business
func (k Keeper) EstimateInterestRate(ctx sdk.Context, business string, businessType string, amount sdk.Coin, creditScore string) sdk.Dec {
	params := k.GetParams(ctx)
	baseRate := sdk.MustNewDecFromStr(params.MaxInterestRate)
	
	// Get business profile
	profile, found := k.GetBusinessProfile(ctx, business)
	if \!found {
		return baseRate
	}
	
	// Apply credit score discount
	score, _ := sdk.NewDecFromStr(creditScore)
	if score.GTE(sdk.NewDec(750)) {
		baseRate = baseRate.Sub(sdk.MustNewDecFromStr("0.015")) // 1.5% discount
	} else if score.GTE(sdk.NewDec(700)) {
		baseRate = baseRate.Sub(sdk.MustNewDecFromStr("0.01")) // 1% discount
	}
	
	// Apply business type discount
	if businessType == "STARTUP" {
		startupDiscount := sdk.MustNewDecFromStr(params.StartupDiscount)
		baseRate = baseRate.Sub(startupDiscount)
	}
	
	// Apply women entrepreneur discount
	if profile.IsWomenOwned {
		womenDiscount := sdk.MustNewDecFromStr(params.WomenEntrepreneurDiscount)
		baseRate = baseRate.Sub(womenDiscount)
	}
	
	// Ensure rate is within bounds
	minRate := sdk.MustNewDecFromStr(params.MinInterestRate)
	if baseRate.LT(minRate) {
		baseRate = minRate
	}
	
	return baseRate
}

// GetApplicableDiscounts returns applicable discounts for a business
func (k Keeper) GetApplicableDiscounts(ctx sdk.Context, business string, businessType string) []string {
	discounts := []string{}
	params := k.GetParams(ctx)
	
	profile, found := k.GetBusinessProfile(ctx, business)
	if \!found {
		return discounts
	}
	
	if profile.CreditScore >= 750 {
		discounts = append(discounts, "Excellent credit score: 1.5% discount")
	} else if profile.CreditScore >= 700 {
		discounts = append(discounts, "Good credit score: 1% discount")
	}
	
	if businessType == "STARTUP" {
		discounts = append(discounts, fmt.Sprintf("Startup business: %s discount", params.StartupDiscount))
	}
	
	if profile.IsWomenOwned {
		discounts = append(discounts, fmt.Sprintf("Women entrepreneur: %s discount", params.WomenEntrepreneurDiscount))
	}
	
	// Check for festival offers
	festivalOffers := k.GetActiveFestivalOffers(ctx)
	for _, offer := range festivalOffers {
		discounts = append(discounts, fmt.Sprintf("%s festival: %s discount", offer.FestivalName, offer.InterestReduction))
	}
	
	return discounts
}

// IsEligibleForCreditLine checks if a business is eligible for a credit line
func (k Keeper) IsEligibleForCreditLine(ctx sdk.Context, business string, annualRevenue sdk.Coin) bool {
	profile, found := k.GetBusinessProfile(ctx, business)
	if \!found {
		return false
	}
	
	// Requirements for credit line
	if profile.CreditScore < 700 {
		return false
	}
	
	if profile.YearsInBusiness < 2 {
		return false
	}
	
	// Minimum revenue requirement
	minRevenue := sdk.NewCoin(annualRevenue.Denom, sdk.NewInt(1000000)) // 10 lakh
	if annualRevenue.IsLT(minRevenue) {
		return false
	}
	
	return true
}

// CalculateMaxCreditLine calculates maximum credit line based on revenue
func (k Keeper) CalculateMaxCreditLine(annualRevenue sdk.Coin) sdk.Coin {
	// 20% of annual revenue as credit line
	maxCredit := annualRevenue.Amount.Mul(sdk.NewInt(20)).Quo(sdk.NewInt(100))
	return sdk.NewCoin(annualRevenue.Denom, maxCredit)
}

// GetRequiredDocuments returns required documents for a business type
func (k Keeper) GetRequiredDocuments(businessType string) []string {
	docs := []string{
		"GST Registration Certificate",
		"PAN Card",
		"Aadhaar Card",
		"Bank statements (last 12 months)",
		"ITR for last 2 years",
		"Business registration documents",
		"Financial statements (P&L, Balance Sheet)",
	}
	
	// Add business type specific documents
	switch businessType {
	case "MANUFACTURING":
		docs = append(docs, "Factory license", "Pollution control certificate")
	case "EXPORT":
		docs = append(docs, "Import Export Code", "Export orders")
	case "STARTUP":
		docs = append(docs, "Business plan", "Pitch deck")
	}
	
	return docs
}

// CalculateDefaultRate calculates the default rate for business loans
func (k Keeper) CalculateDefaultRate(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	loanStore := prefix.NewStore(store, types.LoanKeyPrefix)
	
	totalLoans := int64(0)
	defaultedLoans := int64(0)
	
	iterator := loanStore.Iterator(nil, nil)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var loan types.BusinessLoan
		if err := k.cdc.Unmarshal(iterator.Value(), &loan); err \!= nil {
			continue
		}
		
		totalLoans++
		if loan.Status == types.LoanStatus_LOAN_STATUS_DEFAULTED {
			defaultedLoans++
		}
	}
	
	if totalLoans == 0 {
		return "0%"
	}
	
	rate := sdk.NewDec(defaultedLoans).Mul(sdk.NewDec(100)).Quo(sdk.NewDec(totalLoans))
	return rate.String() + "%"
}

// GetAllBusinessProfiles retrieves all business profiles
func (k Keeper) GetAllBusinessProfiles(ctx sdk.Context) []types.BusinessProfile {
	store := ctx.KVStore(k.storeKey)
	iterator := store.Iterator(types.BusinessProfilePrefix, nil)
	defer iterator.Close()

	var profiles []types.BusinessProfile
	for ; iterator.Valid(); iterator.Next() {
		var profile types.BusinessProfile
		k.cdc.MustUnmarshal(iterator.Value(), &profile)
		profiles = append(profiles, profile)
	}
	return profiles
}

// GetAllBusinessLoans retrieves all business loans
func (k Keeper) GetAllBusinessLoans(ctx sdk.Context) []types.BusinessLoan {
	store := ctx.KVStore(k.storeKey)
	iterator := store.Iterator(types.LoanKeyPrefix, nil)
	defer iterator.Close()

	var loans []types.BusinessLoan
	for ; iterator.Valid(); iterator.Next() {
		var loan types.BusinessLoan
		k.cdc.MustUnmarshal(iterator.Value(), &loan)
		loans = append(loans, loan)
	}
	return loans
}
