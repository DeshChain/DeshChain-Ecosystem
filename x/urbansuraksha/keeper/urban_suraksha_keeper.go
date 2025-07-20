/*
Copyright 2024 DeshChain Foundation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/deshchain/deshchain/x/urbansuraksha/types"
)

// Keeper handles Urban Pension Scheme operations
type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   storetypes.StoreKey
	paramstore paramtypes.Subspace
	bankKeeper types.BankKeeper
	
	// Integration with other modules
	moneyOrderKeeper types.MoneyOrderKeeper
	revenueKeeper    types.RevenueKeeper
}

// NewKeeper creates a new Urban Pension keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	bankKeeper types.BankKeeper,
	moneyOrderKeeper types.MoneyOrderKeeper,
	revenueKeeper types.RevenueKeeper,
) Keeper {
	return Keeper{
		cdc:              cdc,
		storeKey:         storeKey,
		paramstore:       ps,
		bankKeeper:       bankKeeper,
		moneyOrderKeeper: moneyOrderKeeper,
		revenueKeeper:    revenueKeeper,
	}
}

// CreateUrbanSurakshaScheme creates a new urban pension scheme
func (k Keeper) CreateUrbanSurakshaScheme(
	ctx sdk.Context,
	contributorAddr sdk.AccAddress,
	cityCode string,
	referrerAddr sdk.AccAddress,
) (*types.UrbanSurakshaScheme, error) {
	// Generate scheme ID
	schemeID := fmt.Sprintf("UPS-%s-%d", cityCode, ctx.BlockTime().Unix())
	accountID := fmt.Sprintf("UPA-%s-%d", contributorAddr.String()[:8], ctx.BlockTime().Unix())
	
	// Calculate contribution amount (₹2500 worth of NAMO)
	namoPrice := k.getNAMOPrice(ctx) // Get current NAMO price
	contributionAmount := sdk.NewDec(2500).Quo(namoPrice).TruncateInt()
	monthlyContribution := sdk.NewCoin("unamo", contributionAmount)
	
	// Create scheme
	scheme := types.UrbanSurakshaScheme{
		SchemeID:            schemeID,
		AccountID:           accountID,
		ContributorAddress:  contributorAddr,
		MonthlyContribution: monthlyContribution,
		ContributionPeriod:  types.UrbanSurakshaContributionPeriod,
		TotalContributions:  sdk.NewCoin("unamo", sdk.ZeroInt()),
		StartDate:           ctx.BlockTime(),
		MaturityDate:        ctx.BlockTime().AddDate(0, types.UrbanSurakshaMaturityMonth, 0),
		Status:              types.StatusActive,
		PaidContributions:   0,
		MissedContributions: 0,
		NextContribution:    ctx.BlockTime().AddDate(0, 1, 0),
		ReturnPercentage:    sdk.NewDecWithPrec(35, 2), // 35% base return
		MaturityPaid:        false,
		ReferrerAddress:     referrerAddr,
		PerformanceScore:    sdk.OneDec(),
		CommunityRating:     sdk.OneDec(),
	}
	
	// Set default insurance coverage
	scheme.LifeInsuranceCover = sdk.NewCoin("inr", sdk.NewInt(types.DefaultLifeCover))
	scheme.HealthInsuranceCover = sdk.NewCoin("inr", sdk.NewInt(types.DefaultHealthCover))
	
	// Get or create urban pool
	poolID, err := k.getOrCreateUrbanPool(ctx, cityCode)
	if err != nil {
		return nil, err
	}
	scheme.UrbanPoolID = poolID
	
	// Calculate pool contributions (70% to pool, 30% to reserve)
	poolContribution := monthlyContribution.Amount.ToDec().Mul(sdk.NewDecWithPrec(70, 2)).TruncateInt()
	reserveContribution := monthlyContribution.Amount.Sub(poolContribution)
	
	scheme.PoolContribution = sdk.NewCoin("unamo", poolContribution)
	scheme.ReserveContribution = sdk.NewCoin("unamo", reserveContribution)
	
	// Store scheme
	k.SetUrbanSurakshaScheme(ctx, scheme)
	
	// Process referral reward if applicable
	if !referrerAddr.Empty() {
		k.processReferralReward(ctx, referrerAddr, contributorAddr, "pension")
	}
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"urban_pension_created",
			sdk.NewAttribute("scheme_id", schemeID),
			sdk.NewAttribute("contributor", contributorAddr.String()),
			sdk.NewAttribute("monthly_contribution", monthlyContribution.String()),
			sdk.NewAttribute("city_code", cityCode),
		),
	)
	
	return &scheme, nil
}

// ProcessMonthlyContribution processes monthly pension contribution
func (k Keeper) ProcessMonthlyContribution(
	ctx sdk.Context,
	schemeID string,
	contributorAddr sdk.AccAddress,
	amount sdk.Coin,
) error {
	// Get scheme
	scheme, found := k.GetUrbanSurakshaScheme(ctx, schemeID)
	if !found {
		return types.ErrSchemeNotFound
	}
	
	// Validate contributor
	if !scheme.ContributorAddress.Equals(contributorAddr) {
		return types.ErrUnauthorizedContributor
	}
	
	// Validate amount
	if !amount.IsEqual(scheme.MonthlyContribution) {
		return types.ErrInvalidContributionAmount
	}
	
	// Check if contribution is due
	if ctx.BlockTime().Before(scheme.NextContribution.AddDate(0, 0, -7)) { // 7 days early allowed
		return types.ErrContributionTooEarly
	}
	
	// Transfer contribution from contributor
	if err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx, contributorAddr, types.ModuleName, sdk.NewCoins(amount),
	); err != nil {
		return err
	}
	
	// Update scheme
	scheme.PaidContributions++
	scheme.TotalContributions = scheme.TotalContributions.Add(amount)
	scheme.LastContribution = ctx.BlockTime()
	scheme.NextContribution = ctx.BlockTime().AddDate(0, 1, 0)
	
	// Update performance score
	k.updatePerformanceScore(ctx, &scheme)
	
	// Add to urban pool
	err := k.addToUrbanPool(ctx, scheme.UrbanPoolID, scheme.PoolContribution)
	if err != nil {
		return err
	}
	
	// Auto-deduct insurance premium if applicable
	k.processInsurancePremium(ctx, &scheme)
	
	// Store updated scheme
	k.SetUrbanSurakshaScheme(ctx, scheme)
	
	// Process monthly hook for unified liquidity pool integration
	k.moneyOrderKeeper.Hooks().AfterSurakshaContribution(
		ctx, scheme.AccountID, contributorAddr, amount, scheme.UrbanPoolID,
	)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"urban_suraksha_contribution",
			sdk.NewAttribute("scheme_id", schemeID),
			sdk.NewAttribute("contributor", contributorAddr.String()),
			sdk.NewAttribute("amount", amount.String()),
			sdk.NewAttribute("paid_contributions", fmt.Sprintf("%d", scheme.PaidContributions)),
		),
	)
	
	return nil
}

// ApplyForEducationLoan applies for education loan
func (k Keeper) ApplyForEducationLoan(
	ctx sdk.Context,
	schemeID string,
	applicantAddr sdk.AccAddress,
	loanAmount sdk.Coin,
	institutionName string,
	courseType string,
	courseDuration uint32,
) (*types.UrbanEducationLoan, error) {
	// Get pension scheme
	scheme, found := k.GetUrbanSurakshaScheme(ctx, schemeID)
	if !found {
		return nil, types.ErrSchemeNotFound
	}
	
	// Check eligibility
	if !scheme.IsEligibleForEducationLoan() {
		return nil, types.ErrNotEligibleForLoan
	}
	
	// Validate loan amount (max 10x monthly contribution)
	maxLoan := scheme.MonthlyContribution.Amount.MulRaw(10)
	if loanAmount.Amount.GT(maxLoan) {
		return nil, types.ErrLoanAmountTooHigh
	}
	
	// Calculate interest rate based on pension performance
	interestRate := scheme.CalculateEducationLoanRate()
	
	// Generate loan ID
	loanID := fmt.Sprintf("UEL-%s-%d", schemeID[:8], ctx.BlockTime().Unix())
	
	// Create education loan
	loan := types.UrbanEducationLoan{
		LoanID:              loanID,
		PensionAccountID:    scheme.AccountID,
		BorrowerAddress:     applicantAddr,
		LoanAmount:          loanAmount,
		InterestRate:        interestRate,
		LoanTerm:            courseDuration + 12, // Course + 1 year grace period
		InstitutionName:     institutionName,
		CourseType:          courseType,
		CourseDuration:      courseDuration,
		ExpectedCompletion:  ctx.BlockTime().AddDate(0, int(courseDuration), 0),
		RepaymentStartDate:  ctx.BlockTime().AddDate(0, int(courseDuration+6), 0), // 6 months after completion
		OutstandingAmount:   loanAmount,
		Status:              types.LoanStatusPending,
		CreatedAt:           ctx.BlockTime(),
		LastUpdated:         ctx.BlockTime(),
	}
	
	// Calculate monthly EMI
	loan.MonthlyEMI = k.calculateEMI(loanAmount, interestRate, loan.LoanTerm)
	
	// Store loan application
	k.SetUrbanEducationLoan(ctx, loan)
	
	// Update pension scheme
	scheme.EducationLoanTaken = true
	scheme.EducationLoanAmount = loanAmount
	k.SetUrbanSurakshaScheme(ctx, scheme)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"education_loan_applied",
			sdk.NewAttribute("loan_id", loanID),
			sdk.NewAttribute("scheme_id", schemeID),
			sdk.NewAttribute("amount", loanAmount.String()),
			sdk.NewAttribute("interest_rate", interestRate.String()),
			sdk.NewAttribute("institution", institutionName),
		),
	)
	
	return &loan, nil
}

// ProcessSurakshaMaturity processes pension maturity and payout
func (k Keeper) ProcessSurakshaMaturity(
	ctx sdk.Context,
	schemeID string,
) error {
	// Get scheme
	scheme, found := k.GetUrbanSurakshaScheme(ctx, schemeID)
	if !found {
		return types.ErrSchemeNotFound
	}
	
	// Check if maturity date reached
	if ctx.BlockTime().Before(scheme.MaturityDate) {
		return types.ErrMaturityNotReached
	}
	
	// Check if already paid
	if scheme.MaturityPaid {
		return types.ErrMaturityAlreadyPaid
	}
	
	// Get pool performance for return calculation
	pool, found := k.GetUrbanUnifiedPool(ctx, scheme.UrbanPoolID)
	if !found {
		return types.ErrPoolNotFound
	}
	
	// Calculate return amount
	maturityAmount := scheme.CalculateExpectedReturn(pool.PerformanceScore)
	
	// Check if sufficient funds available
	if err := k.validateMaturityFunds(ctx, scheme.UrbanPoolID, maturityAmount); err != nil {
		return err
	}
	
	// Transfer maturity amount to contributor
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx, types.ModuleName, scheme.ContributorAddress, sdk.NewCoins(maturityAmount),
	); err != nil {
		return err
	}
	
	// Update scheme
	scheme.MaturityPaid = true
	scheme.MaturityAmount = maturityAmount
	scheme.Status = types.StatusMatured
	k.SetUrbanSurakshaScheme(ctx, scheme)
	
	// Update pool
	k.updatePoolAfterMaturity(ctx, scheme.UrbanPoolID, maturityAmount)
	
	// Process referral rewards for referrer
	if !scheme.ReferrerAddress.Empty() {
		k.processMaturityReferralBonus(ctx, scheme.ReferrerAddress, maturityAmount)
	}
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"urban_pension_matured",
			sdk.NewAttribute("scheme_id", schemeID),
			sdk.NewAttribute("contributor", scheme.ContributorAddress.String()),
			sdk.NewAttribute("maturity_amount", maturityAmount.String()),
			sdk.NewAttribute("return_percentage", scheme.ReturnPercentage.String()),
		),
	)
	
	return nil
}

// Helper functions

func (k Keeper) getNAMOPrice(ctx sdk.Context) sdk.Dec {
	// Get current NAMO to INR exchange rate
	// This would integrate with price oracle in production
	return sdk.NewDecWithPrec(75, 3) // 0.075 INR per NAMO (placeholder)
}

func (k Keeper) getOrCreateUrbanPool(ctx sdk.Context, cityCode string) (uint64, error) {
	// Check if pool exists for city
	pools := k.GetAllUrbanUnifiedPools(ctx)
	for _, pool := range pools {
		if pool.CityCode == cityCode {
			return pool.PoolID, nil
		}
	}
	
	// Create new pool
	poolID := uint64(len(pools) + 1)
	pool := types.UrbanUnifiedPool{
		PoolID:      poolID,
		PoolName:    fmt.Sprintf("Urban Pool - %s", cityCode),
		CityCode:    cityCode,
		CityName:    k.getCityName(cityCode),
		Status:      types.StatusActive,
		CreatedAt:   ctx.BlockTime(),
		LastUpdated: ctx.BlockTime(),
	}
	
	k.SetUrbanUnifiedPool(ctx, pool)
	return poolID, nil
}

func (k Keeper) addToUrbanPool(ctx sdk.Context, poolID uint64, amount sdk.Coin) error {
	pool, found := k.GetUrbanUnifiedPool(ctx, poolID)
	if !found {
		return types.ErrPoolNotFound
	}
	
	// Add to total liquidity
	pool.TotalLiquidity = pool.TotalLiquidity.Add(amount)
	pool.MonthlyInflow = pool.MonthlyInflow.Add(amount)
	
	// Allocate according to percentages
	totalAmount := amount.Amount.ToDec()
	
	// Pension Reserve (25%)
	pensionAlloc := totalAmount.Mul(sdk.NewDecWithPrec(25, 2)).TruncateInt()
	pool.SurakshaReserve = pool.SurakshaReserve.Add(sdk.NewCoin(amount.Denom, pensionAlloc))
	
	// Education Loan Pool (35%)
	eduAlloc := totalAmount.Mul(sdk.NewDecWithPrec(35, 2)).TruncateInt()
	pool.EducationLoanPool = pool.EducationLoanPool.Add(sdk.NewCoin(amount.Denom, eduAlloc))
	
	// Insurance Reserve (15%)
	insAlloc := totalAmount.Mul(sdk.NewDecWithPrec(15, 2)).TruncateInt()
	pool.InsuranceReserve = pool.InsuranceReserve.Add(sdk.NewCoin(amount.Denom, insAlloc))
	
	// Investment Pool (20%)
	invAlloc := totalAmount.Mul(sdk.NewDecWithPrec(20, 2)).TruncateInt()
	pool.InvestmentPool = pool.InvestmentPool.Add(sdk.NewCoin(amount.Denom, invAlloc))
	
	// Emergency Reserve (5%)
	emergAlloc := totalAmount.Sub(pensionAlloc.ToDec()).Sub(eduAlloc.ToDec()).Sub(insAlloc.ToDec()).Sub(invAlloc.ToDec()).TruncateInt()
	pool.EmergencyReserve = pool.EmergencyReserve.Add(sdk.NewCoin(amount.Denom, emergAlloc))
	
	pool.LastUpdated = ctx.BlockTime()
	k.SetUrbanUnifiedPool(ctx, pool)
	
	return nil
}

func (k Keeper) processReferralReward(ctx sdk.Context, referrerAddr, refereeAddr sdk.AccAddress, rewardType string) {
	// Get referrer's referral count
	referralCount := k.GetReferralCount(ctx, referrerAddr)
	
	// Calculate reward amount
	reward := types.CalculateReferralReward(rewardType, referralCount)
	
	// Transfer reward
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx, types.ModuleName, referrerAddr, sdk.NewCoins(reward),
	); err == nil {
		// Update referral count
		k.SetReferralCount(ctx, referrerAddr, referralCount+1)
		
		// Create referral record
		referralReward := types.ReferralReward{
			ReferrerAddress:  referrerAddr,
			RefereeAddress:   refereeAddr,
			RewardType:       rewardType,
			RewardAmount:     reward,
			RewardDate:       ctx.BlockTime(),
			RewardStatus:     "paid",
			MilestoneLevel:   referralCount + 1,
		}
		k.SetReferralReward(ctx, referralReward)
		
		// Emit event
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				"referral_reward_paid",
				sdk.NewAttribute("referrer", referrerAddr.String()),
				sdk.NewAttribute("referee", refereeAddr.String()),
				sdk.NewAttribute("reward_type", rewardType),
				sdk.NewAttribute("amount", reward.String()),
			),
		)
	}
}

func (k Keeper) calculateEMI(principal sdk.Coin, rate sdk.Dec, termMonths uint32) sdk.Coin {
	// EMI = P * r * (1+r)^n / ((1+r)^n - 1)
	p := principal.Amount.ToDec()
	r := rate.Quo(sdk.NewDec(12)) // Monthly rate
	n := sdk.NewDec(int64(termMonths))
	
	// For simplicity, using simple EMI calculation
	// In production, use compound interest formula
	monthlyPrincipal := p.Quo(n)
	monthlyInterest := p.Mul(r)
	emi := monthlyPrincipal.Add(monthlyInterest)
	
	return sdk.NewCoin(principal.Denom, emi.TruncateInt())
}

func (k Keeper) updatePerformanceScore(ctx sdk.Context, scheme *types.UrbanSurakshaScheme) {
	// Calculate performance based on payment consistency
	totalExpected := scheme.PaidContributions + scheme.MissedContributions
	if totalExpected == 0 {
		return
	}
	
	paymentRatio := sdk.NewDec(int64(scheme.PaidContributions)).Quo(sdk.NewDec(int64(totalExpected)))
	
	// Bonus for early payments, penalty for late
	timeBonusPoints := sdk.ZeroDec()
	if ctx.BlockTime().Before(scheme.NextContribution) {
		timeBonusPoints = sdk.NewDecWithPrec(5, 2) // 5% bonus for early payment
	}
	
	scheme.PerformanceScore = paymentRatio.Add(timeBonusPoints)
	if scheme.PerformanceScore.GT(sdk.OneDec()) {
		scheme.PerformanceScore = sdk.OneDec()
	}
}

func (k Keeper) processInsurancePremium(ctx sdk.Context, scheme *types.UrbanSurakshaScheme) {
	premium := scheme.CalculateInsurancePremium(scheme.LifeInsuranceCover, scheme.HealthInsuranceCover)
	
	// Auto-deduct premium from contribution
	if scheme.MonthlyContribution.Amount.GTE(premium.Amount) {
		scheme.InsurancePremiumPaid = scheme.InsurancePremiumPaid.Add(premium)
		
		// Create/update insurance policy
		k.updateInsurancePolicy(ctx, scheme)
	}
}

func (k Keeper) updateInsurancePolicy(ctx sdk.Context, scheme *types.UrbanSurakshaScheme) {
	// Implementation for insurance policy management
	// This would create or update the insurance policy based on the scheme
}

func (k Keeper) getCityName(cityCode string) string {
	// Map city codes to names
	cityMap := map[string]string{
		"DEL": "Delhi",
		"MUM": "Mumbai", 
		"BLR": "Bangalore",
		"HYD": "Hyderabad",
		"CHE": "Chennai",
		"KOL": "Kolkata",
		"PUN": "Pune",
		"AHM": "Ahmedabad",
	}
	
	if name, exists := cityMap[cityCode]; exists {
		return name
	}
	return cityCode
}

// Storage functions will be implemented here (SetUrbanSurakshaScheme, GetUrbanSurakshaScheme, etc.)