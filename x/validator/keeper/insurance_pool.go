package keeper

import (
    "fmt"
    "time"
    
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/deshchain/namo/x/validator/types"
)

// InsurancePool manages the validator insurance pool
type InsurancePool struct {
    TotalFunds       sdk.Int   `json:"total_funds"`
    ReservedFunds    sdk.Int   `json:"reserved_funds"`
    AvailableFunds   sdk.Int   `json:"available_funds"`
    TotalClaims      uint64    `json:"total_claims"`
    ApprovedClaims   uint64    `json:"approved_claims"`
    TotalPaidOut     sdk.Int   `json:"total_paid_out"`
    LastUpdateTime   time.Time `json:"last_update_time"`
}

// InsuranceClaim represents a claim against the insurance pool
type InsuranceClaim struct {
    ClaimID          uint64    `json:"claim_id"`
    Claimant         string    `json:"claimant"`
    ValidatorAddr    string    `json:"validator_address"`
    ClaimAmount      sdk.Int   `json:"claim_amount"`
    DumpAmount       sdk.Int   `json:"dump_amount"`
    PriceImpact      sdk.Dec   `json:"price_impact"`
    Evidence         string    `json:"evidence"`
    SubmissionTime   time.Time `json:"submission_time"`
    Status           string    `json:"status"` // pending, approved, rejected
    VotingEndTime    time.Time `json:"voting_end_time"`
    ApprovalVotes    uint64    `json:"approval_votes"`
    RejectionVotes   uint64    `json:"rejection_votes"`
    PayoutAmount     sdk.Int   `json:"payout_amount"`
    PayoutTime       time.Time `json:"payout_time"`
}

// InitializeInsurancePool sets up the insurance pool
func (k Keeper) InitializeInsurancePool(ctx sdk.Context) error {
    pool := InsurancePool{
        TotalFunds:     sdk.ZeroInt(),
        ReservedFunds:  sdk.ZeroInt(),
        AvailableFunds: sdk.ZeroInt(),
        TotalClaims:    0,
        ApprovedClaims: 0,
        TotalPaidOut:   sdk.ZeroInt(),
        LastUpdateTime: ctx.BlockTime(),
    }
    
    k.SetInsurancePool(ctx, pool)
    return nil
}

// FundInsurancePool adds funds to the insurance pool
func (k Keeper) FundInsurancePool(ctx sdk.Context, amount sdk.Int) error {
    pool := k.GetInsurancePool(ctx)
    
    pool.TotalFunds = pool.TotalFunds.Add(amount)
    pool.AvailableFunds = pool.AvailableFunds.Add(amount)
    pool.LastUpdateTime = ctx.BlockTime()
    
    k.SetInsurancePool(ctx, pool)
    
    ctx.EventManager().EmitEvent(
        sdk.NewEvent(
            "insurance_pool_funded",
            sdk.NewAttribute("amount", amount.String()),
            sdk.NewAttribute("total_funds", pool.TotalFunds.String()),
        ),
    )
    
    return nil
}

// FileInsuranceClaim creates a new insurance claim
func (k Keeper) FileInsuranceClaim(
    ctx sdk.Context,
    claimant sdk.AccAddress,
    validatorAddr string,
    dumpAmount sdk.Int,
    priceImpact sdk.Dec,
    evidence string,
) (uint64, error) {
    // Validate claim conditions
    if dumpAmount.LT(sdk.NewInt(100000000000)) { // 100K NAMO minimum
        return 0, fmt.Errorf("dump amount below minimum threshold")
    }
    
    if priceImpact.LT(sdk.NewDecWithPrec(10, 2)) { // 10% minimum impact
        return 0, fmt.Errorf("price impact below minimum threshold")
    }
    
    // Calculate claim amount (up to 20% of dump value)
    claimAmount := dumpAmount.ToDec().Mul(sdk.NewDecWithPrec(20, 2)).TruncateInt()
    
    // Apply tier-based coverage limits
    stake, found := k.GetValidatorStake(ctx, validatorAddr)
    if !found {
        return 0, fmt.Errorf("validator stake not found")
    }
    
    var maxCoverage sdk.Int
    switch stake.Tier {
    case 1: // Tier 1: Up to $100K
        maxCoverage = sdk.NewInt(100000000000) // Assuming $0.10 per NAMO
    case 2: // Tier 2: Up to $250K
        maxCoverage = sdk.NewInt(250000000000)
    case 3: // Tier 3: Up to $500K
        maxCoverage = sdk.NewInt(500000000000)
    default:
        maxCoverage = sdk.NewInt(50000000000) // $50K default
    }
    
    if claimAmount.GT(maxCoverage) {
        claimAmount = maxCoverage
    }
    
    // Apply deductible (10%)
    deductible := claimAmount.ToDec().Mul(sdk.NewDecWithPrec(10, 2)).TruncateInt()
    claimAmount = claimAmount.Sub(deductible)
    
    // Create claim
    claimID := k.GetNextClaimID(ctx)
    claim := InsuranceClaim{
        ClaimID:        claimID,
        Claimant:       claimant.String(),
        ValidatorAddr:  validatorAddr,
        ClaimAmount:    claimAmount,
        DumpAmount:     dumpAmount,
        PriceImpact:    priceImpact,
        Evidence:       evidence,
        SubmissionTime: ctx.BlockTime(),
        Status:         "pending",
        VotingEndTime:  ctx.BlockTime().Add(7 * 24 * time.Hour), // 7 days voting
        ApprovalVotes:  0,
        RejectionVotes: 0,
        PayoutAmount:   sdk.ZeroInt(),
    }
    
    k.SetInsuranceClaim(ctx, claim)
    k.SetNextClaimID(ctx, claimID+1)
    
    // Reserve funds
    pool := k.GetInsurancePool(ctx)
    pool.ReservedFunds = pool.ReservedFunds.Add(claimAmount)
    pool.AvailableFunds = pool.AvailableFunds.Sub(claimAmount)
    pool.TotalClaims++
    k.SetInsurancePool(ctx, pool)
    
    ctx.EventManager().EmitEvent(
        sdk.NewEvent(
            types.EventTypeInsuranceClaimFiled,
            sdk.NewAttribute(types.AttributeKeyClaimID, fmt.Sprintf("%d", claimID)),
            sdk.NewAttribute("claimant", claimant.String()),
            sdk.NewAttribute("validator", validatorAddr),
            sdk.NewAttribute(types.AttributeKeyClaimAmount, claimAmount.String()),
            sdk.NewAttribute("dump_amount", dumpAmount.String()),
            sdk.NewAttribute("price_impact", priceImpact.String()),
        ),
    )
    
    return claimID, nil
}

// ProcessInsuranceClaim handles the approval/rejection of a claim
func (k Keeper) ProcessInsuranceClaim(
    ctx sdk.Context,
    claimID uint64,
    approved bool,
) error {
    claim, found := k.GetInsuranceClaim(ctx, claimID)
    if !found {
        return fmt.Errorf("claim %d not found", claimID)
    }
    
    if claim.Status != "pending" {
        return fmt.Errorf("claim already processed")
    }
    
    pool := k.GetInsurancePool(ctx)
    
    if approved {
        // Check if enough funds available
        if pool.AvailableFunds.LT(claim.ClaimAmount) {
            return fmt.Errorf("insufficient funds in insurance pool")
        }
        
        // Pay out claim over 6 months
        claim.Status = "approved"
        claim.PayoutAmount = claim.ClaimAmount
        claim.PayoutTime = ctx.BlockTime()
        
        // Update pool
        pool.ReservedFunds = pool.ReservedFunds.Sub(claim.ClaimAmount)
        pool.TotalPaidOut = pool.TotalPaidOut.Add(claim.ClaimAmount)
        pool.ApprovedClaims++
        
        // Schedule monthly payouts
        monthlyPayout := claim.ClaimAmount.QuoRaw(6)
        k.ScheduleMonthlyPayouts(ctx, claim.Claimant, monthlyPayout, 6)
        
        ctx.EventManager().EmitEvent(
            sdk.NewEvent(
                types.EventTypeInsuranceClaimApproved,
                sdk.NewAttribute(types.AttributeKeyClaimID, fmt.Sprintf("%d", claimID)),
                sdk.NewAttribute("payout_amount", claim.PayoutAmount.String()),
            ),
        )
    } else {
        // Reject claim and release reserved funds
        claim.Status = "rejected"
        pool.ReservedFunds = pool.ReservedFunds.Sub(claim.ClaimAmount)
        pool.AvailableFunds = pool.AvailableFunds.Add(claim.ClaimAmount)
    }
    
    k.SetInsuranceClaim(ctx, claim)
    k.SetInsurancePool(ctx, pool)
    
    return nil
}

// GetInsurancePoolMetrics returns current pool statistics
func (k Keeper) GetInsurancePoolMetrics(ctx sdk.Context) InsurancePoolMetrics {
    pool := k.GetInsurancePool(ctx)
    
    // Calculate coverage ratio
    totalStakes := sdk.ZeroInt()
    stakes := k.GetAllValidatorStakes(ctx)
    for _, stake := range stakes {
        totalStakes = totalStakes.Add(stake.NAMOTokensStaked)
    }
    
    coverageRatio := sdk.ZeroDec()
    if !totalStakes.IsZero() {
        coverageRatio = pool.AvailableFunds.ToDec().Quo(totalStakes.ToDec())
    }
    
    return InsurancePoolMetrics{
        TotalFunds:      pool.TotalFunds,
        AvailableFunds:  pool.AvailableFunds,
        ReservedFunds:   pool.ReservedFunds,
        TotalClaims:     pool.TotalClaims,
        ApprovedClaims:  pool.ApprovedClaims,
        RejectionRate:   sdk.NewDec(int64(pool.TotalClaims - pool.ApprovedClaims)).Quo(sdk.NewDec(int64(pool.TotalClaims))),
        AverageClaimSize: pool.TotalPaidOut.ToDec().Quo(sdk.NewDec(int64(pool.ApprovedClaims))),
        CoverageRatio:   coverageRatio,
    }
}

// Storage functions

func (k Keeper) GetInsurancePool(ctx sdk.Context) InsurancePool {
    store := ctx.KVStore(k.storeKey)
    bz := store.Get([]byte("insurance_pool"))
    
    var pool InsurancePool
    if bz != nil {
        k.cdc.MustUnmarshal(bz, &pool)
    }
    
    return pool
}

func (k Keeper) SetInsurancePool(ctx sdk.Context, pool InsurancePool) {
    store := ctx.KVStore(k.storeKey)
    bz := k.cdc.MustMarshal(&pool)
    store.Set([]byte("insurance_pool"), bz)
}

func (k Keeper) GetInsuranceClaim(ctx sdk.Context, claimID uint64) (InsuranceClaim, bool) {
    store := ctx.KVStore(k.storeKey)
    bz := store.Get(types.GetInsuranceClaimKey(claimID))
    
    if bz == nil {
        return InsuranceClaim{}, false
    }
    
    var claim InsuranceClaim
    k.cdc.MustUnmarshal(bz, &claim)
    return claim, true
}

func (k Keeper) SetInsuranceClaim(ctx sdk.Context, claim InsuranceClaim) {
    store := ctx.KVStore(k.storeKey)
    bz := k.cdc.MustMarshal(&claim)
    store.Set(types.GetInsuranceClaimKey(claim.ClaimID), bz)
}

func (k Keeper) GetNextClaimID(ctx sdk.Context) uint64 {
    store := ctx.KVStore(k.storeKey)
    bz := store.Get([]byte("next_claim_id"))
    if bz == nil {
        return 1
    }
    return sdk.BigEndianToUint64(bz)
}

func (k Keeper) SetNextClaimID(ctx sdk.Context, id uint64) {
    store := ctx.KVStore(k.storeKey)
    store.Set([]byte("next_claim_id"), sdk.Uint64ToBigEndian(id))
}

func (k Keeper) ScheduleMonthlyPayouts(ctx sdk.Context, recipient string, amount sdk.Int, months int) {
    // This would integrate with a scheduled task system
    // For now, we'll emit an event
    ctx.EventManager().EmitEvent(
        sdk.NewEvent(
            "insurance_payout_scheduled",
            sdk.NewAttribute("recipient", recipient),
            sdk.NewAttribute("monthly_amount", amount.String()),
            sdk.NewAttribute("months", fmt.Sprintf("%d", months)),
        ),
    )
}

// InsurancePoolMetrics for reporting
type InsurancePoolMetrics struct {
    TotalFunds       sdk.Int `json:"total_funds"`
    AvailableFunds   sdk.Int `json:"available_funds"`
    ReservedFunds    sdk.Int `json:"reserved_funds"`
    TotalClaims      uint64  `json:"total_claims"`
    ApprovedClaims   uint64  `json:"approved_claims"`
    RejectionRate    sdk.Dec `json:"rejection_rate"`
    AverageClaimSize sdk.Dec `json:"average_claim_size"`
    CoverageRatio    sdk.Dec `json:"coverage_ratio"`
}