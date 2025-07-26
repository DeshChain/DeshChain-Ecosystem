package keeper

import (
	"fmt"

	"cosmossdk.io/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"cosmossdk.io/store/prefix"
	"github.com/deshchain/deshchain/x/shikshamitra/types"
)

// Keeper of the shikshamitra store
type Keeper struct {
	cdc             codec.BinaryCodec
	storeKey        sdk.StoreKey
	memKey          sdk.StoreKey
	paramspace      types.ParamSubspace
	bankKeeper      types.BankKeeper
	accountKeeper   types.AccountKeeper
	dhanpataKeeper  types.DhanPataKeeper
	liquidityKeeper types.LiquidityManagerKeeper // REVOLUTIONARY: Member-only lending verification
}

// NewKeeper creates a new shikshamitra Keeper instance
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
	return log.NewLogger(ctx.Logger()).With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// SetEducationLoan stores an education loan in the store
func (k Keeper) SetEducationLoan(ctx sdk.Context, loan types.EducationLoan) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&loan)
	store.Set(types.GetLoanKey(loan.ID), bz)
}

// REVOLUTIONARY STAGED LOAN FUNCTIONS

// CalculateTotalSemesters calculates total semesters based on course type
func (k Keeper) CalculateTotalSemesters(courseType string, duration int32) int32 {
	switch courseType {
	case "BACHELORS":
		return 6  // 3 years × 2 semesters
	case "ENGINEERING", "B_TECH":
		return 8  // 4 years × 2 semesters  
	case "MASTERS", "MBA", "M_TECH":
		return 4  // 2 years × 2 semesters
	case "BACHELORS_MASTERS_INTEGRATED":
		return 10 // 5 years × 2 semesters
	case "DOCTORATE", "PHD":
		return duration / 6 // Custom duration in months ÷ 6 months per semester
	default:
		return duration / 6 // Fallback: total months ÷ 6
	}
}

// CalculateSemesterFunding calculates platform and student portions
func (k Keeper) CalculateSemesterFunding(semesterFee sdk.Coin) (platformPortion, studentDeposit, processingFee sdk.Coin) {
	// 80% platform funding, 20% student deposit
	platformAmount := semesterFee.Amount.MulRaw(80).QuoRaw(100)
	studentAmount := semesterFee.Amount.MulRaw(20).QuoRaw(100)
	
	// 1% processing fee on platform portion, capped at ₹2500
	feeAmount := platformAmount.MulRaw(1).QuoRaw(100)
	maxFee := sdk.NewInt(250000000) // ₹2500 in smallest denomination
	if feeAmount.GT(maxFee) {
		feeAmount = maxFee
	}
	
	platformPortion = sdk.NewCoin(semesterFee.Denom, platformAmount)
	studentDeposit = sdk.NewCoin(semesterFee.Denom, studentAmount)
	processingFee = sdk.NewCoin(semesterFee.Denom, feeAmount)
	
	return platformPortion, studentDeposit, processingFee
}

// ValidateAcademicProgress validates if student can proceed to next semester
func (k Keeper) ValidateAcademicProgress(ctx sdk.Context, studentID string, currentSemester int32) (canProceed bool, reason string) {
	// Get student's academic records
	profile, found := k.GetStudentProfile(ctx, studentID)
	if !found {
		return false, "Student profile not found"
	}
	
	// For first semester, just check admission
	if currentSemester <= 1 {
		return true, "Initial admission verified"
	}
	
	// Check if previous semester was passed
	if len(profile.AcademicRecords) == 0 {
		return false, "No academic records found for previous semester"
	}
	
	lastRecord := profile.AcademicRecords[len(profile.AcademicRecords)-1]
	
	// Check passing criteria: >= 50% marks and >= 75% attendance
	if lastRecord.Percentage.LT(sdk.NewDec(50)) {
		return false, fmt.Sprintf("Failed previous semester with %.1f%% marks (50%% required)", lastRecord.Percentage.MustFloat64())
	}
	
	attendance, _ := sdk.NewDecFromStr(lastRecord.Attendance)
	if attendance.LT(sdk.NewDec(75)) {
		return false, fmt.Sprintf("Insufficient attendance: %.1f%% (75%% required)", attendance.MustFloat64())
	}
	
	return true, "Previous semester passed successfully"
}

// REVOLUTIONARY PERFORMANCE-BASED INTEREST INCENTIVE SYSTEM

// CalculatePerformanceBasedInterestReduction calculates interest reduction based on academic performance
func (k Keeper) CalculatePerformanceBasedInterestReduction(ctx sdk.Context, loanID string, semesterResults []types.SemesterResult) (semesterReductions []sdk.Dec, totalWaiver sdk.Dec, performanceReport string) {
	semesterReductions = make([]sdk.Dec, len(semesterResults))
	allSemestersAbove90 := true
	excellentSemesterCount := 0
	
	performanceDetails := []string{}
	
	// Check each semester for 90%+ performance
	for i, result := range semesterResults {
		if result.Percentage.GTE(sdk.NewDec(90)) && result.Attendance.GTE(sdk.NewDec(75)) {
			// 0.25% interest reduction for this semester
			semesterReductions[i] = sdk.NewDecWithPrec(25, 4) // 0.0025 = 0.25%
			excellentSemesterCount++
			performanceDetails = append(performanceDetails, 
				fmt.Sprintf("Semester %d: %.1f%% marks → 0.25%% interest reduction", 
					i+1, result.Percentage.MustFloat64()))
		} else {
			semesterReductions[i] = sdk.ZeroDec()
			allSemestersAbove90 = false
			if result.Percentage.LT(sdk.NewDec(90)) {
				performanceDetails = append(performanceDetails, 
					fmt.Sprintf("Semester %d: %.1f%% marks → no reduction (90%% required)", 
						i+1, result.Percentage.MustFloat64()))
			}
		}
	}
	
	// REVOLUTIONARY BONUS: 1% total interest waiver for consistent 90%+ performance
	if allSemestersAbove90 && len(semesterResults) > 0 {
		totalWaiver = sdk.NewDecWithPrec(1, 2) // 0.01 = 1%
		performanceReport = fmt.Sprintf(
			"🏆 ACADEMIC EXCELLENCE ACHIEVEMENT!\n"+
			"✅ All %d semesters with 90%+ marks\n"+
			"🎯 Per-semester reductions: %d × 0.25%% = %.2f%%\n"+
			"🚀 Excellence bonus: 1%% total interest waiver\n"+
			"💰 Total benefit: %.2f%% + 1%% waiver",
			len(semesterResults), excellentSemesterCount, 
			float64(excellentSemesterCount)*0.25,
			float64(excellentSemesterCount)*0.25)
	} else {
		totalWaiver = sdk.ZeroDec()
		performanceReport = fmt.Sprintf(
			"📊 PERFORMANCE SUMMARY:\n"+
			"✅ Excellent semesters (90%+): %d/%d\n"+
			"🎯 Interest reductions earned: %.2f%%\n"+
			"📈 Excellence bonus status: %s\n"+
			"💡 Tip: Maintain 90%+ in all semesters for 1%% total waiver!",
			excellentSemesterCount, len(semesterResults),
			float64(excellentSemesterCount)*0.25,
			func() string {
				if allSemestersAbove90 {
					return "QUALIFIED (1% waiver applied)"
				} else {
					return fmt.Sprintf("IN PROGRESS (%d/%d semesters completed)", excellentSemesterCount, len(semesterResults))
				}
			}())
	}
	
	return semesterReductions, totalWaiver, performanceReport
}

// CalculateEffectiveInterestRate calculates the final interest rate after performance incentives
func (k Keeper) CalculateEffectiveInterestRate(ctx sdk.Context, loanID string, baseInterestRate sdk.Dec, semesterResults []types.SemesterResult) (effectiveRate sdk.Dec, savings sdk.Dec, incentiveBreakdown string) {
	// Get performance-based reductions
	semesterReductions, totalWaiver, performanceReport := k.CalculatePerformanceBasedInterestReduction(ctx, loanID, semesterResults)
	
	// Calculate total semester-wise reduction
	totalSemesterReduction := sdk.ZeroDec()
	for _, reduction := range semesterReductions {
		totalSemesterReduction = totalSemesterReduction.Add(reduction)
	}
	
	// Apply semester-wise reductions to base rate
	effectiveRate = baseInterestRate.Sub(totalSemesterReduction)
	
	// Apply total waiver if qualified (further reduces the already reduced rate)
	if !totalWaiver.IsZero() {
		effectiveRate = effectiveRate.Sub(totalWaiver)
	}
	
	// Ensure rate doesn't go below minimum (2% floor for platform sustainability)
	minRate := sdk.NewDecWithPrec(2, 2) // 2%
	if effectiveRate.LT(minRate) {
		effectiveRate = minRate
	}
	
	// Calculate total savings
	savings = baseInterestRate.Sub(effectiveRate)
	
	// Create detailed breakdown
	incentiveBreakdown = fmt.Sprintf(
		"REVOLUTIONARY PERFORMANCE INCENTIVE BREAKDOWN:\n\n"+
		"📊 Base Interest Rate: %.2f%%\n"+
		"🎯 Semester-wise Reductions: -%.2f%%\n"+
		"🏆 Excellence Waiver: -%.2f%%\n"+
		"💰 Final Effective Rate: %.2f%%\n"+
		"🚀 Total Savings: %.2f%%\n\n"+
		"%s",
		baseInterestRate.Mul(sdk.NewDec(100)).MustFloat64(),
		totalSemesterReduction.Mul(sdk.NewDec(100)).MustFloat64(),
		totalWaiver.Mul(sdk.NewDec(100)).MustFloat64(),
		effectiveRate.Mul(sdk.NewDec(100)).MustFloat64(),
		savings.Mul(sdk.NewDec(100)).MustFloat64(),
		performanceReport)
	
	return effectiveRate, savings, incentiveBreakdown
}

// UpdateLoanWithPerformanceIncentives updates loan with performance-based interest adjustments
func (k Keeper) UpdateLoanWithPerformanceIncentives(ctx sdk.Context, loanID string) {
	loan, found := k.GetEducationLoan(ctx, loanID)
	if !found {
		return
	}
	
	// Get student's semester results
	profile, found := k.GetStudentProfile(ctx, loan.StudentID)
	if !found {
		return
	}
	
	// Convert academic records to semester results
	semesterResults := make([]types.SemesterResult, len(profile.AcademicRecords))
	for i, record := range profile.AcademicRecords {
		attendance, _ := sdk.NewDecFromStr(record.Attendance)
		semesterResults[i] = types.SemesterResult{
			Semester:   int32(i + 1),
			Percentage: record.Percentage,
			Attendance: attendance,
		}
	}
	
	// Calculate effective interest rate with performance incentives
	baseRate := sdk.MustNewDecFromStr(loan.InterestRate)
	effectiveRate, savings, breakdown := k.CalculateEffectiveInterestRate(ctx, loanID, baseRate, semesterResults)
	
	// Update loan with new effective rate
	loan.EffectiveInterestRate = effectiveRate.String()
	loan.PerformanceSavings = savings.String()
	loan.IncentiveBreakdown = breakdown
	
	k.SetEducationLoan(ctx, loan)
	
	// Emit performance incentive event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"performance_incentive_applied",
			sdk.NewAttribute("loan_id", loanID),
			sdk.NewAttribute("base_rate", baseRate.Mul(sdk.NewDec(100)).String()+"%"),
			sdk.NewAttribute("effective_rate", effectiveRate.Mul(sdk.NewDec(100)).String()+"%"),
			sdk.NewAttribute("total_savings", savings.Mul(sdk.NewDec(100)).String()+"%"),
			sdk.NewAttribute("excellent_semesters", fmt.Sprintf("%d", k.CountExcellentSemesters(semesterResults))),
			sdk.NewAttribute("excellence_bonus", func() string {
				if k.IsEligibleForExcellenceBonus(semesterResults) {
					return "QUALIFIED - 1% waiver applied"
				}
				return "NOT_YET - maintain 90%+ in all semesters"
			}()),
		),
	)
}

// CountExcellentSemesters counts semesters with 90%+ performance
func (k Keeper) CountExcellentSemesters(semesterResults []types.SemesterResult) int {
	count := 0
	for _, result := range semesterResults {
		if result.Percentage.GTE(sdk.NewDec(90)) && result.Attendance.GTE(sdk.NewDec(75)) {
			count++
		}
	}
	return count
}

// IsEligibleForExcellenceBonus checks if student qualifies for 1% total waiver
func (k Keeper) IsEligibleForExcellenceBonus(semesterResults []types.SemesterResult) bool {
	if len(semesterResults) == 0 {
		return false
	}
	
	for _, result := range semesterResults {
		if result.Percentage.LT(sdk.NewDec(90)) || result.Attendance.LT(sdk.NewDec(75)) {
			return false
		}
	}
	return true
}

// CreateSemesterDisbursement creates a disbursement record for a semester
func (k Keeper) CreateSemesterDisbursement(ctx sdk.Context, loanID string, semester int32, semesterFee sdk.Coin, studentDepositPaid bool) (*types.SemesterDisbursement, error) {
	// Validate academic progress
	loan, found := k.GetEducationLoan(ctx, loanID)
	if !found {
		return nil, fmt.Errorf("loan not found")
	}
	
	canProceed, reason := k.ValidateAcademicProgress(ctx, loan.StudentID, semester)
	if !canProceed {
		// REVOLUTIONARY: Terminate loan due to academic failure
		k.TerminateLoanForAcademicFailure(ctx, loanID, reason)
		return nil, fmt.Errorf("academic validation failed: %s", reason)
	}
	
	// Check if student has paid 20% deposit
	if !studentDepositPaid {
		return nil, fmt.Errorf("student must pay 20%% deposit (₹%s) before platform disbursement", 
			semesterFee.Amount.MulRaw(20).QuoRaw(100).String())
	}
	
	// Calculate funding breakdown
	platformPortion, studentDeposit, processingFee := k.CalculateSemesterFunding(semesterFee)
	
	// Create disbursement record
	disbursement := &types.SemesterDisbursement{
		LoanID:           loanID,
		Semester:         semester,
		SemesterFee:      semesterFee,
		PlatformPortion:  platformPortion,
		StudentDeposit:   studentDeposit,
		ProcessingFee:    processingFee,
		DisbursedAmount:  platformPortion.Sub(processingFee), // 80% - 1% fee
		Status:           "APPROVED",
		DisbursementDate: ctx.BlockTime(),
		ValidatedBy:      "academic_progress_gate",
		Remarks:          fmt.Sprintf("Semester %d: %s", semester, reason),
	}
	
	return disbursement, nil
}

// TerminateLoanForAcademicFailure terminates loan when student fails academically
func (k Keeper) TerminateLoanForAcademicFailure(ctx sdk.Context, loanID, reason string) {
	loan, found := k.GetEducationLoan(ctx, loanID)
	if !found {
		return
	}
	
	// REVOLUTIONARY COMMUNITY GOVERNANCE: Check if student deserves a second chance
	excellenceHistory := k.CalculateExcellenceHistory(ctx, loanID)
	if k.IsEligibleForCommunityVote(excellenceHistory) {
		// Instead of immediate termination, trigger community vote
		k.InitiateCommunityVoteForContinuation(ctx, loanID, reason, excellenceHistory)
		return
	}
	
	// Regular termination for students without excellent track record
	k.ProcessStandardTermination(ctx, loanID, reason)
}

// REVOLUTIONARY COMMUNITY VOTING SYSTEM

// CalculateExcellenceHistory calculates student's academic excellence track record
func (k Keeper) CalculateExcellenceHistory(ctx sdk.Context, loanID string) types.ExcellenceHistory {
	loan, found := k.GetEducationLoan(ctx, loanID)
	if !found {
		return types.ExcellenceHistory{}
	}
	
	profile, found := k.GetStudentProfile(ctx, loan.StudentID)
	if !found {
		return types.ExcellenceHistory{}
	}
	
	excellentSemesters := 0
	totalSemesters := len(profile.AcademicRecords)
	averageMarks := sdk.ZeroDec()
	
	// Analyze each semester's performance
	for _, record := range profile.AcademicRecords {
		if record.Percentage.GTE(sdk.NewDec(90)) {
			excellentSemesters++
		}
		averageMarks = averageMarks.Add(record.Percentage)
	}
	
	if totalSemesters > 0 {
		averageMarks = averageMarks.Quo(sdk.NewDec(int64(totalSemesters)))
	}
	
	excellencePercentage := sdk.ZeroDec()
	if totalSemesters > 0 {
		excellencePercentage = sdk.NewDec(int64(excellentSemesters)).Quo(sdk.NewDec(int64(totalSemesters))).Mul(sdk.NewDec(100))
	}
	
	return types.ExcellenceHistory{
		LoanID:              loanID,
		TotalSemesters:      int32(totalSemesters),
		ExcellentSemesters:  int32(excellentSemesters),
		AverageMarks:        averageMarks,
		ExcellencePercent:   excellencePercentage,
		HasExcellenceRecord: excellentSemesters >= 2 && excellencePercentage.GTE(sdk.NewDec(60)), // 60% of semesters excellent
	}
}

// IsEligibleForCommunityVote checks if student qualifies for community vote
func (k Keeper) IsEligibleForCommunityVote(history types.ExcellenceHistory) bool {
	// Criteria for community vote eligibility:
	// 1. At least 2 excellent semesters (90%+)
	// 2. At least 60% of semesters were excellent
	// 3. Average marks ≥ 80%
	
	return history.HasExcellenceRecord && 
		   history.ExcellentSemesters >= 2 && 
		   history.ExcellencePercent.GTE(sdk.NewDec(60)) && 
		   history.AverageMarks.GTE(sdk.NewDec(80))
}

// InitiateCommunityVoteForContinuation starts community voting process
func (k Keeper) InitiateCommunityVoteForContinuation(ctx sdk.Context, loanID, failureReason string, history types.ExcellenceHistory) {
	loan, found := k.GetEducationLoan(ctx, loanID)
	if !found {
		return
	}
	
	// Create community vote proposal
	voteID := k.GenerateVoteID(ctx)
	votingPeriod := ctx.BlockTime().AddDate(0, 0, 7) // 7 days voting period
	
	vote := types.CommunityVote{
		VoteID:          voteID,
		LoanID:          loanID,
		StudentID:       loan.StudentID,
		ProposalType:    "LOAN_CONTINUATION",
		FailureReason:   failureReason,
		ExcellenceHistory: &history,
		VotingStartDate: ctx.BlockTime(),
		VotingEndDate:   votingPeriod,
		Status:          "ACTIVE",
		YesVotes:        sdk.ZeroInt(),
		NoVotes:         sdk.ZeroInt(),
		TotalVoters:     0,
		RequiredQuorum:  sdk.NewDecWithPrec(30, 2), // 30% quorum required
		PassingThreshold: sdk.NewDecWithPrec(60, 2), // 60% yes votes to pass
	}
	
	// Update loan status to pending community vote
	loan.Status = "PENDING_COMMUNITY_VOTE"
	loan.CommunityVote = &vote
	loan.TerminationReason = failureReason
	k.SetEducationLoan(ctx, loan)
	
	// Store vote separately for easy querying
	k.SetCommunityVote(ctx, vote)
	
	// Emit community vote initiation event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"community_vote_initiated",
			sdk.NewAttribute("vote_id", voteID),
			sdk.NewAttribute("loan_id", loanID),
			sdk.NewAttribute("student_id", loan.StudentID),
			sdk.NewAttribute("failure_reason", failureReason),
			sdk.NewAttribute("excellent_semesters", fmt.Sprintf("%d/%d", history.ExcellentSemesters, history.TotalSemesters)),
			sdk.NewAttribute("average_marks", fmt.Sprintf("%.1f%%", history.AverageMarks.MustFloat64())),
			sdk.NewAttribute("excellence_percentage", fmt.Sprintf("%.1f%%", history.ExcellencePercent.MustFloat64())),
			sdk.NewAttribute("voting_end_date", votingPeriod.Format("2006-01-02")),
			sdk.NewAttribute("message", "🗳️ Community will decide: Does this excellent student deserve a second chance?"),
		),
	)
}

// ProcessStandardTermination handles regular termination for non-excellent students
func (k Keeper) ProcessStandardTermination(ctx sdk.Context, loanID, reason string) {
	loan, found := k.GetEducationLoan(ctx, loanID)
	if !found {
		return
	}
	
	// Update loan status
	loan.Status = "TERMINATED_ACADEMIC_FAILURE"
	loan.TerminationDate = &ctx.BlockTime()
	loan.TerminationReason = reason
	
	// Calculate total amount to be repaid (only disbursed semesters)
	totalDisbursed := sdk.ZeroInt()
	for _, disbursement := range loan.SemesterDisbursements {
		if disbursement.Status == "DISBURSED" {
			totalDisbursed = totalDisbursed.Add(disbursement.DisbursedAmount.Amount)
		}
	}
	
	// REVOLUTIONARY REPAYMENT LOGIC: Immediate repayment for failed students
	// No interest charged for social consideration, but immediate repayment required
	loan.TotalRepayment = sdk.NewCoin("NAMO", totalDisbursed)
	loan.RepaymentStartDate = &ctx.BlockTime() // IMMEDIATE REPAYMENT
	loan.RepaymentStatus = "ACTIVE"
	
	k.SetEducationLoan(ctx, loan)
	
	// Emit termination event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"education_loan_terminated_failure",
			sdk.NewAttribute("loan_id", loanID),
			sdk.NewAttribute("reason", reason),
			sdk.NewAttribute("total_disbursed", totalDisbursed.String()),
			sdk.NewAttribute("repayment_amount", loan.TotalRepayment.String()),
			sdk.NewAttribute("repayment_start", "IMMEDIATE - no grace period for academic failure"),
			sdk.NewAttribute("interest_charged", "NO - social consideration for failed students"),
		),
	)
}

// SubmitCommunityVote allows community members to vote on loan continuation
func (k Keeper) SubmitCommunityVote(ctx sdk.Context, voteID string, voterAddr sdk.AccAddress, voteChoice string, voterStake sdk.Coin) error {
	vote, found := k.GetCommunityVote(ctx, voteID)
	if !found {
		return fmt.Errorf("community vote not found")
	}
	
	// Check if voting period is active
	if ctx.BlockTime().After(vote.VotingEndDate) {
		return fmt.Errorf("voting period has ended")
	}
	
	if vote.Status != "ACTIVE" {
		return fmt.Errorf("vote is not active")
	}
	
	// Check if voter has already voted
	if k.HasVoterVoted(ctx, voteID, voterAddr.String()) {
		return fmt.Errorf("voter has already voted")
	}
	
	// Verify voter is a pool member with stake
	if !k.liquidityKeeper.IsPoolMember(ctx, voterAddr) {
		return fmt.Errorf("only pool members can vote")
	}
	
	// Calculate vote weight based on stake (1 NAMO = 1 vote, minimum 100 NAMO to vote)
	minStake := sdk.NewCoin("NAMO", sdk.NewInt(10000000000)) // 100 NAMO minimum
	if voterStake.IsLT(minStake) {
		return fmt.Errorf("minimum stake of 100 NAMO required to vote")
	}
	
	voteWeight := voterStake.Amount
	
	// Record vote
	voterRecord := types.VoterRecord{
		VoteID:      voteID,
		VoterAddr:   voterAddr.String(),
		VoteChoice:  voteChoice,
		VoteWeight:  voteWeight,
		VoteDate:    ctx.BlockTime(),
		VoterStake:  voterStake,
	}
	
	k.SetVoterRecord(ctx, voterRecord)
	
	// Update vote tallies
	if voteChoice == "YES" {
		vote.YesVotes = vote.YesVotes.Add(voteWeight)
	} else if voteChoice == "NO" {
		vote.NoVotes = vote.NoVotes.Add(voteWeight)
	} else {
		return fmt.Errorf("invalid vote choice: %s", voteChoice)
	}
	
	vote.TotalVoters++
	vote.TotalStakeVoted = vote.TotalStakeVoted.Add(voteWeight)
	
	k.SetCommunityVote(ctx, vote)
	
	// Emit vote submission event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"community_vote_submitted",
			sdk.NewAttribute("vote_id", voteID),
			sdk.NewAttribute("voter", voterAddr.String()),
			sdk.NewAttribute("choice", voteChoice),
			sdk.NewAttribute("weight", voteWeight.String()),
			sdk.NewAttribute("current_yes", vote.YesVotes.String()),
			sdk.NewAttribute("current_no", vote.NoVotes.String()),
		),
	)
	
	return nil
}

// FinalizeCommunityVote processes the final vote result
func (k Keeper) FinalizeCommunityVote(ctx sdk.Context, voteID string) {
	vote, found := k.GetCommunityVote(ctx, voteID)
	if !found {
		return
	}
	
	// Check if voting period has ended
	if ctx.BlockTime().Before(vote.VotingEndDate) {
		return // Voting still active
	}
	
	if vote.Status != "ACTIVE" {
		return // Already finalized
	}
	
	loan, found := k.GetEducationLoan(ctx, vote.LoanID)
	if !found {
		return
	}
	
	// Calculate vote results
	totalVotes := vote.YesVotes.Add(vote.NoVotes)
	totalPoolStake := k.liquidityKeeper.GetTotalPoolStake(ctx)
	quorumMet := vote.TotalStakeVoted.ToDec().Quo(totalPoolStake.ToDec()).GTE(vote.RequiredQuorum)
	
	var voteResult string
	var voteApproved bool
	
	if !quorumMet {
		voteResult = "FAILED_QUORUM"
		voteApproved = false
	} else if totalVotes.IsZero() {
		voteResult = "NO_VOTES"
		voteApproved = false
	} else {
		yesPercentage := vote.YesVotes.ToDec().Quo(totalVotes.ToDec())
		if yesPercentage.GTE(vote.PassingThreshold) {
			voteResult = "APPROVED"
			voteApproved = true
		} else {
			voteResult = "REJECTED"
			voteApproved = false
		}
	}
	
	// Update vote status
	vote.Status = voteResult
	vote.FinalizedDate = &ctx.BlockTime()
	k.SetCommunityVote(ctx, vote)
	
	// Process result
	if voteApproved {
		// Grant second chance with conditions
		k.GrantSecondChance(ctx, vote.LoanID, vote)
	} else {
		// Proceed with termination
		k.ProcessStandardTermination(ctx, vote.LoanID, vote.FailureReason)
	}
	
	// Emit finalization event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"community_vote_finalized",
			sdk.NewAttribute("vote_id", voteID),
			sdk.NewAttribute("loan_id", vote.LoanID),
			sdk.NewAttribute("result", voteResult),
			sdk.NewAttribute("yes_votes", vote.YesVotes.String()),
			sdk.NewAttribute("no_votes", vote.NoVotes.String()),
			sdk.NewAttribute("quorum_met", fmt.Sprintf("%t", quorumMet)),
			sdk.NewAttribute("approved", fmt.Sprintf("%t", voteApproved)),
		),
	)
}

// GrantSecondChance grants continuation with strict conditions
func (k Keeper) GrantSecondChance(ctx sdk.Context, loanID string, vote types.CommunityVote) {
	loan, found := k.GetEducationLoan(ctx, loanID)
	if !found {
		return
	}
	
	// Grant second chance with CONDITIONS
	loan.Status = "SECOND_CHANCE_GRANTED"
	loan.SecondChanceGranted = true
	loan.SecondChanceConditions = []string{
		"Must achieve ≥80% marks in next semester",
		"Must maintain ≥90% attendance",
		"Additional 1% interest penalty applied",
		"No further community votes if failed again",
		"Monthly progress reporting required",
	}
	loan.SecondChanceDeadline = ctx.BlockTime().AddDate(0, 6, 0) // 6 months to prove
	
	// Apply 1% interest penalty for the second chance
	currentRate := sdk.MustNewDecFromStr(loan.InterestRate)
	penaltyRate := currentRate.Add(sdk.NewDecWithPrec(1, 2)) // +1%
	loan.InterestRate = penaltyRate.String()
	
	k.SetEducationLoan(ctx, loan)
	
	// Emit second chance granted event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"second_chance_granted",
			sdk.NewAttribute("loan_id", loanID),
			sdk.NewAttribute("student_id", loan.StudentID),
			sdk.NewAttribute("deadline", loan.SecondChanceDeadline.Format("2006-01-02")),
			sdk.NewAttribute("penalty_rate", penaltyRate.Mul(sdk.NewDec(100)).String()+"%"),
			sdk.NewAttribute("message", "🙏 Community believes in your potential - don't let them down!"),
		),
	)
}

// REVOLUTIONARY LOAN WRITE-OFF SYSTEM FOR ECOSYSTEM SUSTAINABILITY

// ClassifyNonPerformingAsset classifies loans as NPAs based on overdue status
func (k Keeper) ClassifyNonPerformingAsset(ctx sdk.Context, loanID string) types.NPAClassification {
	loan, found := k.GetEducationLoan(ctx, loanID)
	if !found {
		return types.NPAClassification{Classification: "NOT_FOUND"}
	}
	
	// Only classify loans that are in repayment status
	if loan.RepaymentStatus != "ACTIVE" {
		return types.NPAClassification{Classification: "NOT_ELIGIBLE"}
	}
	
	// Calculate days overdue
	currentDate := ctx.BlockTime()
	var daysOverdue int64
	
	if loan.RepaymentStartDate != nil && currentDate.After(*loan.RepaymentStartDate) {
		// Simplified calculation - in production, use actual missed payment dates
		monthsOverdue := int64(currentDate.Sub(*loan.RepaymentStartDate).Hours() / (24 * 30))
		daysOverdue = monthsOverdue * 30
	}
	
	// Classify based on Basel III norms adapted for education loans
	classification := types.NPAClassification{
		LoanID:      loanID,
		DaysOverdue: daysOverdue,
		AssessmentDate: currentDate,
	}
	
	if daysOverdue >= 540 { // 18 months - Loss Assets
		classification.Classification = "LOSS_ASSET"
		classification.ProvisionRequired = sdk.NewDecWithPrec(100, 2) // 100%
		classification.WriteOffEligible = true
	} else if daysOverdue >= 360 { // 12 months - Doubtful Assets
		classification.Classification = "DOUBTFUL_ASSET"
		classification.ProvisionRequired = sdk.NewDecWithPrec(75, 2) // 75%
		classification.WriteOffEligible = true
	} else if daysOverdue >= 180 { // 6 months - Substandard Assets
		classification.Classification = "SUBSTANDARD"
		classification.ProvisionRequired = sdk.NewDecWithPrec(25, 2) // 25%
		classification.WriteOffEligible = false
	} else if daysOverdue >= 90 { // 3 months - Special Mention
		classification.Classification = "SPECIAL_MENTION"
		classification.ProvisionRequired = sdk.NewDecWithPrec(5, 2) // 5%
		classification.WriteOffEligible = false
	} else {
		classification.Classification = "STANDARD"
		classification.ProvisionRequired = sdk.NewDecWithPrec(1, 2) // 1%
		classification.WriteOffEligible = false
	}
	
	return classification
}

// InitiateWriteOffVote starts community voting for loan write-off
func (k Keeper) InitiateWriteOffVote(ctx sdk.Context, loanID string, writeOffReason string) error {
	loan, found := k.GetEducationLoan(ctx, loanID)
	if !found {
		return fmt.Errorf("loan not found")
	}
	
	// Check NPA classification
	npaClass := k.ClassifyNonPerformingAsset(ctx, loanID)
	if !npaClass.WriteOffEligible {
		return fmt.Errorf("loan not eligible for write-off. Classification: %s", npaClass.Classification)
	}
	
	// Calculate financial impact
	impact := k.CalculateWriteOffImpact(ctx, loanID)
	
	// Generate write-off vote ID
	writeOffVoteID := k.GenerateWriteOffVoteID(ctx)
	votingPeriod := ctx.BlockTime().AddDate(0, 0, 14) // 14 days for critical decisions
	
	// Create write-off vote
	writeOffVote := types.WriteOffVote{
		VoteID:           writeOffVoteID,
		LoanID:           loanID,
		StudentID:        loan.StudentID,
		ProposalType:     "LOAN_WRITEOFF",
		WriteOffReason:   writeOffReason,
		NPAClassification: &npaClass,
		FinancialImpact:  &impact,
		VotingStartDate:  ctx.BlockTime(),
		VotingEndDate:    votingPeriod,
		Status:           "ACTIVE",
		YesVotes:         sdk.ZeroInt(),
		NoVotes:          sdk.ZeroInt(),
		TotalVoters:      0,
		RequiredQuorum:   sdk.NewDecWithPrec(50, 2), // 50% quorum for critical decisions
		PassingThreshold: sdk.NewDecWithPrec(80, 2), // 80% YES votes required
	}
	
	// Update loan status
	loan.Status = "PENDING_WRITEOFF_VOTE"
	loan.WriteOffVote = &writeOffVote
	k.SetEducationLoan(ctx, loan)
	
	// Store vote separately
	k.SetWriteOffVote(ctx, writeOffVote)
	
	// Emit write-off vote initiation event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"writeoff_vote_initiated",
			sdk.NewAttribute("vote_id", writeOffVoteID),
			sdk.NewAttribute("loan_id", loanID),
			sdk.NewAttribute("npa_classification", npaClass.Classification),
			sdk.NewAttribute("days_overdue", fmt.Sprintf("%d", npaClass.DaysOverdue)),
			sdk.NewAttribute("outstanding_amount", impact.OutstandingAmount.String()),
			sdk.NewAttribute("ecosystem_impact", fmt.Sprintf("%.2f%%", impact.EcosystemImpactPercent.MustFloat64())),
			sdk.NewAttribute("voting_end_date", votingPeriod.Format("2006-01-02")),
			sdk.NewAttribute("required_threshold", "80% YES votes"),
			sdk.NewAttribute("message", "⚠️ CRITICAL DECISION: Community must decide on loan write-off"),
		),
	)
	
	return nil
}

// CalculateWriteOffImpact calculates the financial impact of writing off a loan
func (k Keeper) CalculateWriteOffImpact(ctx sdk.Context, loanID string) types.WriteOffImpact {
	loan, found := k.GetEducationLoan(ctx, loanID)
	if !found {
		return types.WriteOffImpact{}
	}
	
	// Calculate outstanding amount
	outstandingAmount := loan.TotalRepayment.Sub(loan.RepaidAmount)
	
	// Calculate total ecosystem pool value
	totalPoolValue := k.liquidityKeeper.GetTotalPoolValue(ctx)
	
	// Calculate impact percentage
	impactPercent := sdk.ZeroDec()
	if totalPoolValue.IsPositive() {
		impactPercent = outstandingAmount.Amount.ToDec().Quo(totalPoolValue.ToDec()).Mul(sdk.NewDec(100))
	}
	
	// Calculate provision coverage
	npaClass := k.ClassifyNonPerformingAsset(ctx, loanID)
	provisionAmount := outstandingAmount.Amount.ToDec().Mul(npaClass.ProvisionRequired).TruncateInt()
	
	// Get current ecosystem metrics
	totalNPAs := k.GetTotalNPACount(ctx)
	npaRatio := k.CalculateNPARatio(ctx)
	
	return types.WriteOffImpact{
		LoanID:                 loanID,
		OutstandingAmount:      outstandingAmount,
		ProvisionAmount:        sdk.NewCoin(outstandingAmount.Denom, provisionAmount),
		EcosystemImpactPercent: impactPercent,
		TotalPoolValue:         totalPoolValue,
		CurrentNPACount:        totalNPAs,
		CurrentNPARatio:        npaRatio,
		ProjectedNPARatio:      k.CalculateProjectedNPARatio(ctx, outstandingAmount),
		RecommendedAction:      k.GetWriteOffRecommendation(impactPercent, npaClass.Classification),
	}
}

// SubmitWriteOffVote allows community members to vote on loan write-offs
func (k Keeper) SubmitWriteOffVote(ctx sdk.Context, voteID string, voterAddr sdk.AccAddress, voteChoice string, voterStake sdk.Coin) error {
	vote, found := k.GetWriteOffVote(ctx, voteID)
	if !found {
		return fmt.Errorf("write-off vote not found")
	}
	
	// Check voting period
	if ctx.BlockTime().After(vote.VotingEndDate) {
		return fmt.Errorf("voting period has ended")
	}
	
	if vote.Status != "ACTIVE" {
		return fmt.Errorf("vote is not active")
	}
	
	// Check if voter already voted
	if k.HasWriteOffVoterVoted(ctx, voteID, voterAddr.String()) {
		return fmt.Errorf("voter has already voted")
	}
	
	// Verify voter is pool member with sufficient stake
	if !k.liquidityKeeper.IsPoolMember(ctx, voterAddr) {
		return fmt.Errorf("only pool members can vote on write-offs")
	}
	
	// Higher minimum stake for write-off votes (500 NAMO)
	minStake := sdk.NewCoin("NAMO", sdk.NewInt(50000000000)) // 500 NAMO minimum
	if voterStake.IsLT(minStake) {
		return fmt.Errorf("minimum stake of 500 NAMO required for write-off votes")
	}
	
	voteWeight := voterStake.Amount
	
	// Record vote
	voterRecord := types.WriteOffVoterRecord{
		VoteID:      voteID,
		VoterAddr:   voterAddr.String(),
		VoteChoice:  voteChoice,
		VoteWeight:  voteWeight,
		VoteDate:    ctx.BlockTime(),
		VoterStake:  voterStake,
	}
	
	k.SetWriteOffVoterRecord(ctx, voterRecord)
	
	// Update vote tallies
	if voteChoice == "YES" {
		vote.YesVotes = vote.YesVotes.Add(voteWeight)
	} else if voteChoice == "NO" {
		vote.NoVotes = vote.NoVotes.Add(voteWeight)
	} else {
		return fmt.Errorf("invalid vote choice: %s", voteChoice)
	}
	
	vote.TotalVoters++
	vote.TotalStakeVoted = vote.TotalStakeVoted.Add(voteWeight)
	
	k.SetWriteOffVote(ctx, vote)
	
	// Emit vote submission event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"writeoff_vote_submitted",
			sdk.NewAttribute("vote_id", voteID),
			sdk.NewAttribute("voter", voterAddr.String()),
			sdk.NewAttribute("choice", voteChoice),
			sdk.NewAttribute("weight", voteWeight.String()),
			sdk.NewAttribute("current_yes", vote.YesVotes.String()),
			sdk.NewAttribute("current_no", vote.NoVotes.String()),
			sdk.NewAttribute("threshold_progress", fmt.Sprintf("%.1f%% of 80%% required", 
				vote.YesVotes.ToDec().Quo(vote.YesVotes.Add(vote.NoVotes).ToDec()).Mul(sdk.NewDec(100)).MustFloat64())),
		),
	)
	
	return nil
}

// FinalizeWriteOffVote processes the final write-off vote result
func (k Keeper) FinalizeWriteOffVote(ctx sdk.Context, voteID string) {
	vote, found := k.GetWriteOffVote(ctx, voteID)
	if !found {
		return
	}
	
	// Check if voting period has ended
	if ctx.BlockTime().Before(vote.VotingEndDate) {
		return // Voting still active
	}
	
	if vote.Status != "ACTIVE" {
		return // Already finalized
	}
	
	loan, found := k.GetEducationLoan(ctx, vote.LoanID)
	if !found {
		return
	}
	
	// Calculate vote results
	totalVotes := vote.YesVotes.Add(vote.NoVotes)
	totalPoolStake := k.liquidityKeeper.GetTotalPoolStake(ctx)
	quorumMet := vote.TotalStakeVoted.ToDec().Quo(totalPoolStake.ToDec()).GTE(vote.RequiredQuorum)
	
	var voteResult string
	var writeOffApproved bool
	
	if !quorumMet {
		voteResult = "FAILED_QUORUM"
		writeOffApproved = false
	} else if totalVotes.IsZero() {
		voteResult = "NO_VOTES"
		writeOffApproved = false
	} else {
		yesPercentage := vote.YesVotes.ToDec().Quo(totalVotes.ToDec())
		if yesPercentage.GTE(vote.PassingThreshold) {
			voteResult = "APPROVED"
			writeOffApproved = true
		} else {
			voteResult = "REJECTED"
			writeOffApproved = false
		}
	}
	
	// Update vote status
	vote.Status = voteResult
	vote.FinalizedDate = &ctx.BlockTime()
	k.SetWriteOffVote(ctx, vote)
	
	// Process result
	if writeOffApproved {
		k.ProcessLoanWriteOff(ctx, vote.LoanID, vote)
	} else {
		// Keep loan as NPA, continue collection efforts
		loan.Status = "NPA_COLLECTION_ACTIVE"
		k.SetEducationLoan(ctx, loan)
	}
	
	// Emit finalization event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"writeoff_vote_finalized",
			sdk.NewAttribute("vote_id", voteID),
			sdk.NewAttribute("loan_id", vote.LoanID),
			sdk.NewAttribute("result", voteResult),
			sdk.NewAttribute("yes_votes", vote.YesVotes.String()),
			sdk.NewAttribute("no_votes", vote.NoVotes.String()),
			sdk.NewAttribute("yes_percentage", fmt.Sprintf("%.1f%%", 
				vote.YesVotes.ToDec().Quo(totalVotes.ToDec()).Mul(sdk.NewDec(100)).MustFloat64())),
			sdk.NewAttribute("threshold_required", "80.0%"),
			sdk.NewAttribute("quorum_met", fmt.Sprintf("%t", quorumMet)),
			sdk.NewAttribute("approved", fmt.Sprintf("%t", writeOffApproved)),
		),
	)
}

// ProcessLoanWriteOff executes the actual loan write-off
func (k Keeper) ProcessLoanWriteOff(ctx sdk.Context, loanID string, vote types.WriteOffVote) {
	loan, found := k.GetEducationLoan(ctx, loanID)
	if !found {
		return
	}
	
	// Calculate write-off amounts
	outstandingAmount := loan.TotalRepayment.Sub(loan.RepaidAmount)
	writeOffDate := ctx.BlockTime()
	
	// Update loan status
	loan.Status = "WRITTEN_OFF"
	loan.WriteOffDate = &writeOffDate
	loan.WriteOffAmount = outstandingAmount
	loan.WriteOffReason = vote.WriteOffReason
	loan.WriteOffVoteID = vote.VoteID
	
	// Transfer from provision pool (already reserved)
	provisionAmount := outstandingAmount.Amount.ToDec().Mul(vote.NPAClassification.ProvisionRequired).TruncateInt()
	actualWriteOff := sdk.NewCoin(outstandingAmount.Denom, provisionAmount)
	
	k.SetEducationLoan(ctx, loan)
	
	// Update ecosystem metrics
	k.UpdateEcosystemNPAMetrics(ctx, outstandingAmount)
	
	// Emit write-off execution event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"loan_written_off",
			sdk.NewAttribute("loan_id", loanID),
			sdk.NewAttribute("student_id", loan.StudentID),
			sdk.NewAttribute("outstanding_amount", outstandingAmount.String()),
			sdk.NewAttribute("provision_used", actualWriteOff.String()),
			sdk.NewAttribute("write_off_date", writeOffDate.Format("2006-01-02")),
			sdk.NewAttribute("vote_id", vote.VoteID),
			sdk.NewAttribute("ecosystem_impact", "NPA removed from active portfolio"),
			sdk.NewAttribute("message", "💔 Community decision: Loan written off for ecosystem health"),
		),
	)
}

// Helper functions for write-off system

func (k Keeper) GetTotalNPACount(ctx sdk.Context) int64 {
	// Implementation to count total NPAs
	return 0 // Placeholder
}

func (k Keeper) CalculateNPARatio(ctx sdk.Context) sdk.Dec {
	// Implementation to calculate NPA ratio
	return sdk.ZeroDec() // Placeholder
}

func (k Keeper) CalculateProjectedNPARatio(ctx sdk.Context, writeOffAmount sdk.Coin) sdk.Dec {
	// Implementation to calculate projected NPA ratio after write-off
	return sdk.ZeroDec() // Placeholder
}

func (k Keeper) GetWriteOffRecommendation(impactPercent sdk.Dec, classification string) string {
	if impactPercent.GT(sdk.NewDec(5)) {
		return "HIGH_IMPACT - Consider restructuring before write-off"
	} else if classification == "LOSS_ASSET" {
		return "RECOMMENDED - Loss asset eligible for write-off"
	} else {
		return "EVALUATE - Continue collection efforts"
	}
}

func (k Keeper) UpdateEcosystemNPAMetrics(ctx sdk.Context, writeOffAmount sdk.Coin) {
	// Implementation to update ecosystem-wide NPA metrics
}

// Storage functions for write-off votes
func (k Keeper) SetWriteOffVote(ctx sdk.Context, vote types.WriteOffVote) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&vote)
	store.Set(types.GetWriteOffVoteKey(vote.VoteID), bz)
}

func (k Keeper) GetWriteOffVote(ctx sdk.Context, voteID string) (types.WriteOffVote, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetWriteOffVoteKey(voteID))
	if bz == nil {
		return types.WriteOffVote{}, false
	}
	
	var vote types.WriteOffVote
	k.cdc.MustUnmarshal(bz, &vote)
	return vote, true
}

func (k Keeper) SetWriteOffVoterRecord(ctx sdk.Context, record types.WriteOffVoterRecord) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&record)
	store.Set(types.GetWriteOffVoterKey(record.VoteID, record.VoterAddr), bz)
}

func (k Keeper) HasWriteOffVoterVoted(ctx sdk.Context, voteID, voterAddr string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetWriteOffVoterKey(voteID, voterAddr))
}

func (k Keeper) GenerateWriteOffVoteID(ctx sdk.Context) string {
	return fmt.Sprintf("writeoff_%d_%d", ctx.BlockHeight(), ctx.BlockTime().Unix())
}

// CompleteLoanAfterGraduation completes loan when student graduates successfully
func (k Keeper) CompleteLoanAfterGraduation(ctx sdk.Context, loanID string) {
	loan, found := k.GetEducationLoan(ctx, loanID)
	if !found {
		return
	}
	
	// Update loan status
	loan.Status = "GRADUATED_SUCCESSFULLY"
	graduationDate := ctx.BlockTime()
	loan.GraduationDate = &graduationDate
	
	// Calculate total amount to be repaid (all disbursed semesters + PERFORMANCE-BASED interest)
	totalDisbursed := sdk.ZeroInt()
	for _, disbursement := range loan.SemesterDisbursements {
		if disbursement.Status == "DISBURSED" {
			totalDisbursed = totalDisbursed.Add(disbursement.DisbursedAmount.Amount)
		}
	}
	
	// REVOLUTIONARY: Apply performance-based interest incentives
	k.UpdateLoanWithPerformanceIncentives(ctx, loanID)
	
	// Reload loan to get updated effective interest rate
	loan, _ = k.GetEducationLoan(ctx, loanID)
	
	// Use effective interest rate (after performance incentives) instead of base rate
	var effectiveInterestRate sdk.Dec
	if loan.EffectiveInterestRate != "" {
		effectiveInterestRate = sdk.MustNewDecFromStr(loan.EffectiveInterestRate)
	} else {
		effectiveInterestRate = sdk.MustNewDecFromStr(loan.InterestRate)
	}
	
	// Calculate average loan duration (disbursement periods)
	avgDurationMonths := int64(len(loan.SemesterDisbursements)) * 6 / 2 // Average time money was borrowed
	
	// REVOLUTIONARY INTEREST CALCULATION: Using performance-adjusted rate
	interest := totalDisbursed.ToDec().Mul(effectiveInterestRate).Mul(sdk.NewDec(avgDurationMonths)).Quo(sdk.NewDec(12))
	totalRepayment := totalDisbursed.Add(interest.TruncateInt())
	
	// REVOLUTIONARY REPAYMENT LOGIC: 6-month grace period after graduation
	repaymentStartDate := graduationDate.AddDate(0, 6, 0) // 6 months after graduation
	
	loan.TotalRepayment = sdk.NewCoin("NAMO", totalRepayment)
	loan.RepaymentStartDate = &repaymentStartDate // 6 MONTHS GRACE PERIOD
	loan.RepaymentStatus = "GRACE_PERIOD"
	
	k.SetEducationLoan(ctx, loan)
	
	// Calculate savings vs base rate
	baseRate := sdk.MustNewDecFromStr(loan.InterestRate)
	baseInterest := totalDisbursed.ToDec().Mul(baseRate).Mul(sdk.NewDec(avgDurationMonths)).Quo(sdk.NewDec(12))
	interestSavings := baseInterest.Sub(interest)
	
	// Emit comprehensive graduation event with performance incentives
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"education_loan_graduated_with_incentives",
			sdk.NewAttribute("loan_id", loanID),
			sdk.NewAttribute("graduation_date", graduationDate.Format("2006-01-02")),
			sdk.NewAttribute("total_disbursed", totalDisbursed.String()),
			sdk.NewAttribute("base_interest_rate", baseRate.Mul(sdk.NewDec(100)).String()+"%"),
			sdk.NewAttribute("effective_interest_rate", effectiveInterestRate.Mul(sdk.NewDec(100)).String()+"%"),
			sdk.NewAttribute("interest_amount", interest.TruncateInt().String()),
			sdk.NewAttribute("interest_savings", interestSavings.TruncateInt().String()),
			sdk.NewAttribute("total_repayment", loan.TotalRepayment.String()),
			sdk.NewAttribute("grace_period", "6 months"),
			sdk.NewAttribute("repayment_start", repaymentStartDate.Format("2006-01-02")),
			sdk.NewAttribute("performance_bonus", loan.IncentiveBreakdown),
			sdk.NewAttribute("success_message", "🎓 Congratulations! Your academic excellence has earned you interest savings!"),
		),
	)
}

// ActivateRepaymentAfterGracePeriod activates repayment after grace period ends
func (k Keeper) ActivateRepaymentAfterGracePeriod(ctx sdk.Context, loanID string) {
	loan, found := k.GetEducationLoan(ctx, loanID)
	if !found {
		return
	}
	
	// Check if grace period has ended
	if loan.RepaymentStartDate != nil && ctx.BlockTime().After(*loan.RepaymentStartDate) {
		if loan.RepaymentStatus == "GRACE_PERIOD" {
			loan.RepaymentStatus = "ACTIVE"
			k.SetEducationLoan(ctx, loan)
			
			// Emit repayment activation event
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					"repayment_activated",
					sdk.NewAttribute("loan_id", loanID),
					sdk.NewAttribute("grace_period_ended", "6 months completed"),
					sdk.NewAttribute("repayment_amount", loan.TotalRepayment.String()),
					sdk.NewAttribute("message", "Employment-based EMI calculation available"),
				),
			)
		}
	}
}

// CheckCollegeFeeDirectPayment validates that disbursement goes directly to college
func (k Keeper) CheckCollegeFeeDirectPayment(ctx sdk.Context, collegeAddress string, feeStructure sdk.Coin) bool {
	// Verify college is registered and authorized to receive fee payments
	// This ensures funds are used only for education, not other expenses
	
	// Check if college address is verified
	college, found := k.GetVerifiedCollege(ctx, collegeAddress)
	if !found {
		return false
	}
	
	// Verify fee structure matches college's declared fee
	if !college.SemesterFee.IsEqual(feeStructure) {
		return false
	}
	
	return true
}

// GetVerifiedCollege retrieves a verified college record
func (k Keeper) GetVerifiedCollege(ctx sdk.Context, collegeAddress string) (types.VerifiedCollege, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetCollegeKey(collegeAddress))
	if bz == nil {
		return types.VerifiedCollege{}, false
	}
	
	var college types.VerifiedCollege
	k.cdc.MustUnmarshal(bz, &college)
	return college, true
}

// CalculateMaximumExposureReduction calculates risk reduction vs traditional loans
func (k Keeper) CalculateMaximumExposureReduction(courseType string, totalCourseFee sdk.Coin) (traditional, staged, riskReduction sdk.Dec) {
	// Traditional model: Full course fee at risk
	traditionalExposure := totalCourseFee.Amount.ToDec()
	
	// Staged model: Only 1 semester at risk (80% of semester fee)
	semesterFee := traditionalExposure.Quo(sdk.NewDec(int64(k.CalculateTotalSemesters(courseType, 0))))
	stagedExposure := semesterFee.Mul(sdk.NewDecWithPrec(80, 2)) // 80% platform portion
	
	// Risk reduction percentage
	riskReduction = traditionalExposure.Sub(stagedExposure).Quo(traditionalExposure).Mul(sdk.NewDec(100))
	
	return traditionalExposure, stagedExposure, riskReduction
}

// GetEducationLoan retrieves an education loan from the store
func (k Keeper) GetEducationLoan(ctx sdk.Context, loanID string) (loan types.EducationLoan, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetLoanKey(loanID))
	if bz == nil {
		return loan, false
	}
	k.cdc.MustUnmarshal(bz, &loan)
	return loan, true
}

// SetStudentProfile stores a student profile
func (k Keeper) SetStudentProfile(ctx sdk.Context, profile types.StudentProfile) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&profile)
	store.Set(types.GetStudentProfileKey(profile.ID), bz)
}

// GetStudentProfile retrieves a student profile
func (k Keeper) GetStudentProfile(ctx sdk.Context, studentID string) (profile types.StudentProfile, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetStudentProfileKey(studentID))
	if bz == nil {
		return profile, false
	}
	k.cdc.MustUnmarshal(bz, &profile)
	return profile, true
}

// GetAllStudentProfiles retrieves all student profiles
func (k Keeper) GetAllStudentProfiles(ctx sdk.Context) []types.StudentProfile {
	store := ctx.KVStore(k.storeKey)
	iterator := store.Iterator(types.StudentProfilePrefix, nil)
	defer iterator.Close()

	var profiles []types.StudentProfile
	for ; iterator.Valid(); iterator.Next() {
		var profile types.StudentProfile
		k.cdc.MustUnmarshal(iterator.Value(), &profile)
		profiles = append(profiles, profile)
	}
	return profiles
}

// GetAllEducationLoans retrieves all education loans
func (k Keeper) GetAllEducationLoans(ctx sdk.Context) []types.EducationLoan {
	store := ctx.KVStore(k.storeKey)
	iterator := store.Iterator(types.LoanKeyPrefix, nil)
	defer iterator.Close()

	var loans []types.EducationLoan
	for ; iterator.Valid(); iterator.Next() {
		var loan types.EducationLoan
		k.cdc.MustUnmarshal(iterator.Value(), &loan)
		loans = append(loans, loan)
	}
	return loans
}

// SetInstitution stores an institution
func (k Keeper) SetInstitution(ctx sdk.Context, institution types.Institution) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&institution)
	store.Set(types.GetInstitutionKey(institution.ID), bz)
}

// GetInstitution retrieves an institution
func (k Keeper) GetInstitution(ctx sdk.Context, institutionID string) (institution types.Institution, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetInstitutionKey(institutionID))
	if bz == nil {
		return institution, false
	}
	k.cdc.MustUnmarshal(bz, &institution)
	return institution, true
}

// SetCourse stores a course
func (k Keeper) SetCourse(ctx sdk.Context, course types.Course) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&course)
	store.Set(types.GetCourseKey(course.ID), bz)
}

// GetCourse retrieves a course
func (k Keeper) GetCourse(ctx sdk.Context, courseID string) (course types.Course, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetCourseKey(courseID))
	if bz == nil {
		return course, false
	}
	k.cdc.MustUnmarshal(bz, &course)
	return course, true
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

// CheckEligibility verifies if a student is eligible for an education loan
func (k Keeper) CheckEligibility(ctx sdk.Context, application types.LoanApplication) (bool, string) {
	// Get student profile
	profile, found := k.GetStudentProfile(ctx, application.StudentID)
	if !found {
		return false, "Student profile not found"
	}

	// Check verification status
	if profile.VerificationStatus != "verified" {
		return false, "Student profile not verified"
	}

	// Check DhanPata address verification
	isDhanPataVerified := k.dhanpataKeeper.IsAddressVerified(ctx, application.DhanPataAddress)
	if !isDhanPataVerified {
		return false, "DhanPata address not verified"
	}

	// Check age (must be under 35 for most courses)
	age := ctx.BlockTime().Year() - profile.DateOfBirth.Year()
	if age > 35 {
		course, _ := k.GetCourse(ctx, application.CourseID)
		if course.MaxAge > 0 && int32(age) > course.MaxAge {
			return false, "Age exceeds maximum limit for this course"
		}
	}

	// Check admission status
	if application.AdmissionStatus != "confirmed" {
		return false, "Admission not confirmed by institution"
	}

	// Check institution recognition
	institution, found := k.GetInstitution(ctx, application.InstitutionID)
	if !found || !institution.IsRecognized {
		return false, "Institution not recognized for education loans"
	}

	// Check co-applicant credit score
	if application.CoApplicantDetails.CreditScore < 650 {
		return false, "Co-applicant credit score too low"
	}

	// Check loan amount limits
	if application.RequestedAmount.Amount.GT(institution.MaxLoanAmount.Amount) {
		return false, "Requested amount exceeds institution's maximum loan limit"
	}

	// Check active loans
	if profile.ActiveLoans >= 2 {
		return false, "Maximum active education loans limit reached"
	}

	return true, ""
}

// CalculateInterestRate calculates interest rate based on various factors
func (k Keeper) CalculateInterestRate(ctx sdk.Context, application types.LoanApplication) sdk.Dec {
	// Base rate
	baseRate, _ := sdk.NewDecFromStr("0.055") // 5.5%

	// Get student profile
	profile, found := k.GetStudentProfile(ctx, application.StudentID)
	if !found {
		maxRate, _ := sdk.NewDecFromStr(types.MaxInterestRate)
		return maxRate
	}

	// Get institution
	institution, found := k.GetInstitution(ctx, application.InstitutionID)
	if !found {
		maxRate, _ := sdk.NewDecFromStr(types.MaxInterestRate)
		return maxRate
	}

	// Merit-based reduction
	meritReduction := sdk.ZeroDec()
	if len(profile.AcademicRecords) > 0 {
		latestRecord := profile.AcademicRecords[len(profile.AcademicRecords)-1]
		if latestRecord.Percentage.GTE(sdk.NewDec(90)) {
			meritReduction, _ = sdk.NewDecFromStr(types.MeritReduction90Plus)
		} else if latestRecord.Percentage.GTE(sdk.NewDec(80)) {
			meritReduction, _ = sdk.NewDecFromStr(types.MeritReduction80Plus)
		} else if latestRecord.Percentage.GTE(sdk.NewDec(70)) {
			meritReduction, _ = sdk.NewDecFromStr(types.MeritReduction70Plus)
		}
	}

	// Institution type reduction
	institutionReduction := sdk.ZeroDec()
	switch institution.Type {
	case types.InstitutionType_IIT:
		institutionReduction = sdk.NewDecWithPrec(-15, 3) // -1.5%
	case types.InstitutionType_IIM:
		institutionReduction = sdk.NewDecWithPrec(-15, 3) // -1.5%
	case types.InstitutionType_NIT:
		institutionReduction = sdk.NewDecWithPrec(-1, 2) // -1%
	case types.InstitutionType_CENTRAL_UNIVERSITY:
		institutionReduction = sdk.NewDecWithPrec(-5, 3) // -0.5%
	}

	// Course type adjustment
	courseAdjustment := sdk.ZeroDec()
	course, found := k.GetCourse(ctx, application.CourseID)
	if found {
		switch course.Type {
		case types.CourseType_PROFESSIONAL:
			courseAdjustment = sdk.NewDecWithPrec(-5, 3) // -0.5%
		case types.CourseType_VOCATIONAL:
			courseAdjustment = sdk.NewDecWithPrec(5, 3) // +0.5%
		case types.CourseType_DOCTORATE:
			courseAdjustment = sdk.NewDecWithPrec(-1, 2) // -1%
		}
	}

	// Co-applicant credit score adjustment
	coApplicantAdjustment := sdk.ZeroDec()
	if application.CoApplicantDetails.CreditScore >= 750 {
		coApplicantAdjustment = sdk.NewDecWithPrec(-5, 3) // -0.5%
	} else if application.CoApplicantDetails.CreditScore < 700 {
		coApplicantAdjustment = sdk.NewDecWithPrec(5, 3) // +0.5%
	}

	// Calculate final rate
	finalRate := baseRate.Sub(meritReduction).Add(institutionReduction).Add(courseAdjustment).Add(coApplicantAdjustment)

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

// CalculateMoratoriumPeriod calculates the grace period before repayment starts
func (k Keeper) CalculateMoratoriumPeriod(ctx sdk.Context, loan types.EducationLoan) int64 {
	// Course duration + grace period
	course, found := k.GetCourse(ctx, loan.CourseID)
	if !found {
		return types.GracePeriodMonths
	}

	// Course duration + 6 months grace period
	return int64(course.Duration) + types.GracePeriodMonths
}

// UpdateEmploymentRecord updates student's employment status for loan tracking
func (k Keeper) UpdateEmploymentRecord(ctx sdk.Context, record types.EmploymentRecord) {
	// Update student profile with employment status
	profile, found := k.GetStudentProfile(ctx, record.StudentID)
	if !found {
		return
	}

	// Store employment record
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&record)
	store.Set(append(types.EmploymentRecordPrefix, []byte(record.StudentID)...), bz)

	// If student has active loans, consider adjusting repayment terms
	// This could trigger automatic repayment plan adjustments based on salary
}

// CalculateLoanStatistics computes overall statistics for education loans
func (k Keeper) CalculateLoanStatistics(ctx sdk.Context) types.LoanStatistics {
	stats := types.LoanStatistics{
		TotalLoansDisbursed: sdk.ZeroInt(),
		ActiveLoans:         0,
		DefaultedLoans:      0,
		AverageInterestRate: sdk.ZeroDec(),
		TotalInterestEarned: sdk.ZeroInt(),
	}

	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.LoanKeyPrefix)
	defer iterator.Close()

	totalRate := sdk.ZeroDec()
	loanCount := 0

	for ; iterator.Valid(); iterator.Next() {
		var loan types.EducationLoan
		k.cdc.MustUnmarshal(iterator.Value(), &loan)

		stats.TotalLoansDisbursed = stats.TotalLoansDisbursed.Add(loan.Amount.Amount)
		
		switch loan.Status {
		case "active":
			stats.ActiveLoans++
		case "defaulted":
			stats.DefaultedLoans++
		}

		// Calculate average interest rate
		rate, _ := sdk.NewDecFromStr(loan.InterestRate)
		totalRate = totalRate.Add(rate)
		loanCount++

		// Calculate total interest earned
		if loan.TotalInterestPaid.IsPositive() {
			stats.TotalInterestEarned = stats.TotalInterestEarned.Add(loan.TotalInterestPaid.Amount)
		}
	}

	if loanCount > 0 {
		stats.AverageInterestRate = totalRate.Quo(sdk.NewDec(int64(loanCount)))
	}

	// Add institution-wise breakdown
	stats.InstitutionBreakdown = k.GetInstitutionBreakdown(ctx)
	stats.CourseTypeBreakdown = k.GetCourseTypeBreakdown(ctx)

	return stats
}

// GetActiveFestivalOffers returns currently active festival offers
func (k Keeper) GetActiveFestivalOffers(ctx sdk.Context) []types.FestivalOffer {
	offers := []types.FestivalOffer{}
	
	// Get current time
	currentTime := ctx.BlockTime()
	
	// Define festival offers with dates
	festivalOffers := []types.FestivalOffer{
		{
			Name:             "Diwali Special",
			FestivalID:       "diwali",
			InterestReduction: sdk.NewDecWithPrec(5, 3), // 0.5% reduction
			ProcessingFeeWaiver: sdk.NewDecWithPrec(50, 2), // 50% waiver
			ValidFrom:         currentTime.AddDate(0, -1, 0), // 1 month before
			ValidTo:           currentTime.AddDate(0, 1, 0),  // 1 month after
			Description:       "Special Diwali offer with reduced interest rates",
		},
		{
			Name:             "Independence Day Scholarship",
			FestivalID:       "independence_day",
			InterestReduction: sdk.NewDecWithPrec(75, 4), // 0.75% reduction
			ProcessingFeeWaiver: sdk.NewDecWithPrec(100, 2), // 100% waiver
			ValidFrom:         currentTime.AddDate(0, -1, 0),
			ValidTo:           currentTime.AddDate(0, 1, 0),
			Description:       "Freedom to learn - Special Independence Day offer",
		},
		{
			Name:             "Teacher's Day Special",
			FestivalID:       "teachers_day",
			InterestReduction: sdk.NewDecWithPrec(1, 2), // 1% reduction
			ProcessingFeeWaiver: sdk.NewDecWithPrec(75, 2), // 75% waiver
			ValidFrom:         currentTime.AddDate(0, -1, 0),
			ValidTo:           currentTime.AddDate(0, 1, 0),
			Description:       "Honoring educators with special loan offers",
		},
	}
	
	// Filter active offers
	for _, offer := range festivalOffers {
		if currentTime.After(offer.ValidFrom) && currentTime.Before(offer.ValidTo) {
			offers = append(offers, offer)
		}
	}
	
	return offers
}

// GetNextFestivalOffer returns the next upcoming festival offer
func (k Keeper) GetNextFestivalOffer(ctx sdk.Context) *types.FestivalOffer {
	currentTime := ctx.BlockTime()
	
	// Get all festival offers
	allOffers := []types.FestivalOffer{
		{
			Name:             "Holi Colorful Future",
			FestivalID:       "holi",
			InterestReduction: sdk.NewDecWithPrec(5, 3),
			ProcessingFeeWaiver: sdk.NewDecWithPrec(25, 2),
			ValidFrom:         currentTime.AddDate(0, 2, 0), // 2 months from now
			ValidTo:           currentTime.AddDate(0, 3, 0),
			Description:       "Add colors to your education journey",
		},
		{
			Name:             "Ganesh Chaturthi Wisdom",
			FestivalID:       "ganesh_chaturthi",
			InterestReduction: sdk.NewDecWithPrec(75, 4),
			ProcessingFeeWaiver: sdk.NewDecWithPrec(50, 2),
			ValidFrom:         currentTime.AddDate(0, 3, 0),
			ValidTo:           currentTime.AddDate(0, 4, 0),
			Description:       "Lord Ganesha's blessings for knowledge seekers",
		},
	}
	
	// Find the next upcoming offer
	var nextOffer *types.FestivalOffer
	minTime := currentTime.AddDate(100, 0, 0) // Far future
	
	for i := range allOffers {
		if allOffers[i].ValidFrom.After(currentTime) && allOffers[i].ValidFrom.Before(minTime) {
			nextOffer = &allOffers[i]
			minTime = allOffers[i].ValidFrom
		}
	}
	
	return nextOffer
}

// CheckStudentEligibility performs comprehensive eligibility check
func (k Keeper) CheckStudentEligibility(ctx sdk.Context, studentID string, institutionType, courseType string) (bool, string, error) {
	// Get student profile
	profile, found := k.GetStudentProfile(ctx, studentID)
	if !found {
		return false, "Student profile not found", nil
	}
	
	// Check basic eligibility
	if profile.VerificationStatus != "verified" {
		return false, "Student profile not verified", nil
	}
	
	// Check age limits based on course type
	age := ctx.BlockTime().Year() - profile.DateOfBirth.Year()
	maxAge := 35 // default
	
	switch courseType {
	case "doctorate":
		maxAge = 40
	case "professional":
		maxAge = 35
	case "vocational":
		maxAge = 45
	}
	
	if age > maxAge {
		return false, fmt.Sprintf("Age exceeds maximum limit of %d years for %s courses", maxAge, courseType), nil
	}
	
	// Check academic requirements
	if len(profile.AcademicRecords) == 0 {
		return false, "No academic records found", nil
	}
	
	latestRecord := profile.AcademicRecords[len(profile.AcademicRecords)-1]
	minPercentage := sdk.NewDec(50) // default minimum
	
	switch institutionType {
	case "iit", "iim":
		minPercentage = sdk.NewDec(75)
	case "nit", "central_university":
		minPercentage = sdk.NewDec(65)
	case "state_university":
		minPercentage = sdk.NewDec(55)
	}
	
	if latestRecord.Percentage.LT(minPercentage) {
		return false, fmt.Sprintf("Academic percentage below minimum requirement of %s%%", minPercentage.String()), nil
	}
	
	// Check for active loan limits
	if profile.ActiveLoans >= 2 {
		return false, "Maximum active education loans limit reached", nil
	}
	
	return true, "Eligible", nil
}

// EstimateInterestRate provides estimated interest rate based on parameters
func (k Keeper) EstimateInterestRate(ctx sdk.Context, studentID, institutionType, courseType string, academicScore, familyIncome sdk.Dec) sdk.Dec {
	// Base rate
	baseRate := sdk.NewDecWithPrec(55, 3) // 5.5%
	
	// Academic score adjustment
	academicAdjustment := sdk.ZeroDec()
	if academicScore.GTE(sdk.NewDec(90)) {
		academicAdjustment = sdk.NewDecWithPrec(-15, 3) // -1.5%
	} else if academicScore.GTE(sdk.NewDec(80)) {
		academicAdjustment = sdk.NewDecWithPrec(-10, 3) // -1.0%
	} else if academicScore.GTE(sdk.NewDec(70)) {
		academicAdjustment = sdk.NewDecWithPrec(-5, 3) // -0.5%
	}
	
	// Institution type adjustment
	institutionAdjustment := sdk.ZeroDec()
	switch institutionType {
	case "iit", "iim":
		institutionAdjustment = sdk.NewDecWithPrec(-15, 3) // -1.5%
	case "nit":
		institutionAdjustment = sdk.NewDecWithPrec(-10, 3) // -1.0%
	case "central_university":
		institutionAdjustment = sdk.NewDecWithPrec(-5, 3) // -0.5%
	}
	
	// Course type adjustment
	courseAdjustment := sdk.ZeroDec()
	switch courseType {
	case "professional":
		courseAdjustment = sdk.NewDecWithPrec(-5, 3) // -0.5%
	case "doctorate":
		courseAdjustment = sdk.NewDecWithPrec(-10, 3) // -1.0%
	case "vocational":
		courseAdjustment = sdk.NewDecWithPrec(5, 3) // +0.5%
	}
	
	// Family income adjustment (economically weaker sections get better rates)
	incomeAdjustment := sdk.ZeroDec()
	annualIncomeLakhs := familyIncome.Quo(sdk.NewDec(100000))
	if annualIncomeLakhs.LTE(sdk.NewDec(3)) { // Less than 3 lakhs
		incomeAdjustment = sdk.NewDecWithPrec(-10, 3) // -1.0%
	} else if annualIncomeLakhs.LTE(sdk.NewDec(6)) { // 3-6 lakhs
		incomeAdjustment = sdk.NewDecWithPrec(-5, 3) // -0.5%
	}
	
	// Calculate final rate
	finalRate := baseRate.Add(academicAdjustment).Add(institutionAdjustment).Add(courseAdjustment).Add(incomeAdjustment)
	
	// Apply bounds
	minRate := sdk.NewDecWithPrec(40, 3) // 4.0%
	maxRate := sdk.NewDecWithPrec(70, 3) // 7.0%
	
	if finalRate.LT(minRate) {
		return minRate
	}
	if finalRate.GT(maxRate) {
		return maxRate
	}
	
	return finalRate
}

// GetApplicableDiscounts returns all applicable discounts for a student
func (k Keeper) GetApplicableDiscounts(ctx sdk.Context, studentID, institutionType, courseType string, academicScore sdk.Dec) []types.DiscountInfo {
	discounts := []types.DiscountInfo{}
	
	// Academic excellence discount
	if academicScore.GTE(sdk.NewDec(95)) {
		discounts = append(discounts, types.DiscountInfo{
			Type:        "Academic Excellence",
			Description: "Top performer discount",
			Value:       sdk.NewDecWithPrec(20, 2), // 20%
		})
	} else if academicScore.GTE(sdk.NewDec(90)) {
		discounts = append(discounts, types.DiscountInfo{
			Type:        "Academic Merit",
			Description: "Merit-based discount",
			Value:       sdk.NewDecWithPrec(15, 2), // 15%
		})
	}
	
	// Premier institution discount
	if institutionType == "iit" || institutionType == "iim" {
		discounts = append(discounts, types.DiscountInfo{
			Type:        "Premier Institution",
			Description: "IIT/IIM student discount",
			Value:       sdk.NewDecWithPrec(10, 2), // 10%
		})
	}
	
	// Women empowerment discount
	profile, found := k.GetStudentProfile(ctx, studentID)
	if found && profile.Gender == "female" {
		discounts = append(discounts, types.DiscountInfo{
			Type:        "Women Empowerment",
			Description: "Special discount for women students",
			Value:       sdk.NewDecWithPrec(5, 2), // 5%
		})
	}
	
	// Rural area discount
	if found && k.IsRuralPincode(profile.Pincode) {
		discounts = append(discounts, types.DiscountInfo{
			Type:        "Rural Development",
			Description: "Support for rural students",
			Value:       sdk.NewDecWithPrec(10, 2), // 10%
		})
	}
	
	// Festival discount (if any active)
	activeFestivals := k.GetActiveFestivalOffers(ctx)
	for _, festival := range activeFestivals {
		discounts = append(discounts, types.DiscountInfo{
			Type:        "Festival Offer",
			Description: festival.Name,
			Value:       festival.ProcessingFeeWaiver,
		})
	}
	
	return discounts
}

// GetRequiredDocuments returns list of required documents for loan application
func (k Keeper) GetRequiredDocuments(ctx sdk.Context, studentID, institutionType, courseType string) []types.DocumentRequirement {
	docs := []types.DocumentRequirement{
		{
			Type:        "Identity Proof",
			Description: "Aadhaar Card / PAN Card / Passport",
			Mandatory:   true,
		},
		{
			Type:        "Address Proof",
			Description: "Aadhaar Card / Utility Bill / Rental Agreement",
			Mandatory:   true,
		},
		{
			Type:        "Academic Records",
			Description: "10th, 12th and graduation marksheets",
			Mandatory:   true,
		},
		{
			Type:        "Admission Letter",
			Description: "Official admission letter from institution",
			Mandatory:   true,
		},
		{
			Type:        "Fee Structure",
			Description: "Complete fee structure from institution",
			Mandatory:   true,
		},
		{
			Type:        "Co-applicant Income Proof",
			Description: "Salary slips / ITR / Bank statements",
			Mandatory:   true,
		},
		{
			Type:        "Co-applicant Identity",
			Description: "Co-applicant's Aadhaar/PAN",
			Mandatory:   true,
		},
	}
	
	// Additional documents for higher loan amounts
	if courseType == "professional" || institutionType == "foreign" {
		docs = append(docs, types.DocumentRequirement{
			Type:        "Entrance Exam Score",
			Description: "CAT/GMAT/GRE/JEE scorecard",
			Mandatory:   true,
		})
	}
	
	// For foreign education
	if institutionType == "foreign" {
		docs = append(docs, []types.DocumentRequirement{
			{
				Type:        "Visa Documents",
				Description: "Student visa approval",
				Mandatory:   true,
			},
			{
				Type:        "I-20 Form",
				Description: "For US universities",
				Mandatory:   false,
			},
			{
				Type:        "Language Test Score",
				Description: "TOEFL/IELTS scorecard",
				Mandatory:   true,
			},
		}...)
	}
	
	return docs
}

// GetCourseDuration returns the duration of a course in months
func (k Keeper) GetCourseDuration(ctx sdk.Context, courseID string) (int32, error) {
	course, found := k.GetCourse(ctx, courseID)
	if !found {
		// Return default durations based on course type
		// This would be better with actual course data
		return 24, fmt.Errorf("course not found, using default duration")
	}
	
	return course.Duration, nil
}

// Helper functions

// GetInstitutionBreakdown returns loan distribution by institution type
func (k Keeper) GetInstitutionBreakdown(ctx sdk.Context) map[string]sdk.Int {
	breakdown := make(map[string]sdk.Int)
	
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.LoanKeyPrefix)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var loan types.EducationLoan
		k.cdc.MustUnmarshal(iterator.Value(), &loan)
		
		institution, found := k.GetInstitution(ctx, loan.InstitutionID)
		if found {
			if _, exists := breakdown[string(institution.Type)]; !exists {
				breakdown[string(institution.Type)] = sdk.ZeroInt()
			}
			breakdown[string(institution.Type)] = breakdown[string(institution.Type)].Add(loan.Amount.Amount)
		}
	}
	
	return breakdown
}

// GetCourseTypeBreakdown returns loan distribution by course type
func (k Keeper) GetCourseTypeBreakdown(ctx sdk.Context) map[string]sdk.Int {
	breakdown := make(map[string]sdk.Int)
	
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.LoanKeyPrefix)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var loan types.EducationLoan
		k.cdc.MustUnmarshal(iterator.Value(), &loan)
		
		course, found := k.GetCourse(ctx, loan.CourseID)
		if found {
			if _, exists := breakdown[string(course.Type)]; !exists {
				breakdown[string(course.Type)] = sdk.ZeroInt()
			}
			breakdown[string(course.Type)] = breakdown[string(course.Type)].Add(loan.Amount.Amount)
		}
	}
	
	return breakdown
}

// IsRuralPincode checks if a pincode is in rural area
func (k Keeper) IsRuralPincode(pincode string) bool {
	// Simplified logic - in production this would check against actual rural pincode database
	// Pincodes starting with 1-3 are generally more urban, 4-9 more rural
	if len(pincode) > 0 {
		firstDigit := pincode[0]
		return firstDigit >= '4' && firstDigit <= '9'
	}
	return false
}

// GetAllLoansByStudent returns all loans for a specific student
func (k Keeper) GetAllLoansByStudent(ctx sdk.Context, studentID string) []types.EducationLoan {
	loans := []types.EducationLoan{}
	
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.LoanKeyPrefix)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var loan types.EducationLoan
		k.cdc.MustUnmarshal(iterator.Value(), &loan)
		
		if loan.StudentID == studentID {
			loans = append(loans, loan)
		}
	}
	
	return loans
}

// GetActiveLoansByStudent returns only active loans for a student
func (k Keeper) GetActiveLoansByStudent(ctx sdk.Context, studentID string) []types.EducationLoan {
	allLoans := k.GetAllLoansByStudent(ctx, studentID)
	activeLoans := []types.EducationLoan{}
	
	for _, loan := range allLoans {
		if loan.Status == "active" || loan.Status == "moratorium" {
			activeLoans = append(activeLoans, loan)
		}
	}
	
	return activeLoans
}

// GetAllScholarships returns all available scholarships
func (k Keeper) GetAllScholarships(ctx sdk.Context) []types.Scholarship {
	scholarships := []types.Scholarship{}
	
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ScholarshipKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var scholarship types.Scholarship
		k.cdc.MustUnmarshal(iterator.Value(), &scholarship)
		scholarships = append(scholarships, scholarship)
	}
	
	return scholarships
}

// GetEligibleScholarships returns scholarships a student is eligible for
func (k Keeper) GetEligibleScholarships(ctx sdk.Context, studentID string) []types.Scholarship {
	allScholarships := k.GetAllScholarships(ctx)
	eligibleScholarships := []types.Scholarship{}
	
	profile, found := k.GetStudentProfile(ctx, studentID)
	if !found {
		return eligibleScholarships
	}
	
	for _, scholarship := range allScholarships {
		if k.IsEligibleForScholarship(ctx, profile, scholarship) {
			eligibleScholarships = append(eligibleScholarships, scholarship)
		}
	}
	
	return eligibleScholarships
}

// IsEligibleForScholarship checks if a student meets scholarship criteria
func (k Keeper) IsEligibleForScholarship(ctx sdk.Context, profile types.StudentProfile, scholarship types.Scholarship) bool {
	// Check academic requirements
	if len(profile.AcademicRecords) > 0 {
		latestRecord := profile.AcademicRecords[len(profile.AcademicRecords)-1]
		if latestRecord.Percentage.LT(scholarship.MinPercentage) {
			return false
		}
	}
	
	// Check income criteria
	if scholarship.MaxFamilyIncome.IsPositive() && profile.FamilyIncome.GT(scholarship.MaxFamilyIncome) {
		return false
	}
	
	// Check gender criteria
	if scholarship.GenderCriteria != "" && scholarship.GenderCriteria != "all" && profile.Gender != scholarship.GenderCriteria {
		return false
	}
	
	// Check category criteria
	if scholarship.CategoryCriteria != "" && scholarship.CategoryCriteria != "all" && profile.Category != scholarship.CategoryCriteria {
		return false
	}
	
	// Check location criteria
	if scholarship.LocationCriteria != "" {
		if scholarship.LocationCriteria == "rural" && !k.IsRuralPincode(profile.Pincode) {
			return false
		}
		if scholarship.LocationCriteria == "urban" && k.IsRuralPincode(profile.Pincode) {
			return false
		}
	}
	
	return true
}