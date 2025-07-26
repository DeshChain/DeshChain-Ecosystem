package keeper

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Metrics contains all prometheus metrics for the CharitableTrust module
type Metrics struct {
	// Trust fund metrics
	TrustFundBalance      prometheus.Gauge
	TrustAllocatedAmount  prometheus.Gauge
	TrustAvailableAmount  prometheus.Gauge
	TrustTotalDistributed prometheus.Gauge

	// Allocation metrics
	CharitableAllocationsTotal      prometheus.Counter
	CharitableAllocationsActive     prometheus.Gauge
	CharitableAllocationsDistributed prometheus.Counter
	CharitableAllocationAmounts     prometheus.Histogram
	AllocationsByCategory           prometheus.GaugeVec
	OrganizationCount               prometheus.Gauge

	// Proposal metrics
	ProposalsCreated      prometheus.Counter
	ProposalsApproved     prometheus.Counter
	ProposalsRejected     prometheus.Counter
	ProposalsExecuted     prometheus.Counter
	ProposalVotingTime    prometheus.Histogram
	ProposalApprovalRate  prometheus.Gauge
	VotesPerProposal      prometheus.Histogram

	// Impact metrics
	ImpactReportsSubmitted     prometheus.Counter
	ImpactReportsVerified      prometheus.Counter
	BeneficiariesReached       prometheus.Counter
	FundsUtilizationRate       prometheus.Gauge
	ImpactScoreDistribution    prometheus.Histogram
	VerificationCompletionTime prometheus.Histogram

	// Fraud metrics
	FraudAlertsCreated      prometheus.Counter
	FraudAlertsInvestigated prometheus.Counter
	FraudAlertsSeverity     prometheus.CounterVec
	FraudInvestigationTime  prometheus.Histogram
	FraudConfirmedCases     prometheus.Counter

	// Governance metrics
	TrusteesActive         prometheus.Gauge
	TrusteeVotingRate      prometheus.Gauge
	GovernanceQuorum       prometheus.Gauge
	AdvisoryCommitteeSize  prometheus.Gauge

	// Performance metrics
	TransactionDuration prometheus.HistogramVec
	QueryDuration       prometheus.HistogramVec
	ValidationErrors    prometheus.CounterVec

	// Transparency metrics
	TransparencyScore           prometheus.Gauge
	PublicReportsPublished      prometheus.Counter
	StakeholderEngagementScore  prometheus.Gauge
	ComplianceViolations        prometheus.Counter
}

// NewMetrics creates a new metrics instance
func NewMetrics(namespace string) *Metrics {
	return &Metrics{
		// Trust fund metrics
		TrustFundBalance: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "charitabletrust",
			Name:      "trust_fund_balance",
			Help:      "Total balance in the charitable trust fund in unamo",
		}),
		TrustAllocatedAmount: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "charitabletrust",
			Name:      "trust_allocated_amount",
			Help:      "Total allocated amount pending distribution in unamo",
		}),
		TrustAvailableAmount: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "charitabletrust",
			Name:      "trust_available_amount",
			Help:      "Available amount for new allocations in unamo",
		}),
		TrustTotalDistributed: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "charitabletrust",
			Name:      "trust_total_distributed",
			Help:      "Total amount distributed to charities in unamo",
		}),

		// Allocation metrics
		CharitableAllocationsTotal: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "charitabletrust",
			Name:      "allocations_total",
			Help:      "Total number of charitable allocations created",
		}),
		CharitableAllocationsActive: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "charitabletrust",
			Name:      "allocations_active",
			Help:      "Number of currently active allocations",
		}),
		CharitableAllocationsDistributed: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "charitabletrust",
			Name:      "allocations_distributed",
			Help:      "Total number of allocations distributed",
		}),
		CharitableAllocationAmounts: promauto.NewHistogram(prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: "charitabletrust",
			Name:      "allocation_amounts",
			Help:      "Distribution of allocation amounts in unamo",
			Buckets:   prometheus.ExponentialBuckets(100000, 10, 8), // 0.1 NAMO to 1M NAMO
		}),
		AllocationsByCategory: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "charitabletrust",
			Name:      "allocations_by_category",
			Help:      "Number of allocations by category",
		}, []string{"category"}),
		OrganizationCount: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "charitabletrust",
			Name:      "organizations_registered",
			Help:      "Total number of registered charitable organizations",
		}),

		// Proposal metrics
		ProposalsCreated: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "charitabletrust",
			Name:      "proposals_created",
			Help:      "Total number of allocation proposals created",
		}),
		ProposalsApproved: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "charitabletrust",
			Name:      "proposals_approved",
			Help:      "Total number of proposals approved",
		}),
		ProposalsRejected: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "charitabletrust",
			Name:      "proposals_rejected",
			Help:      "Total number of proposals rejected",
		}),
		ProposalsExecuted: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "charitabletrust",
			Name:      "proposals_executed",
			Help:      "Total number of proposals executed",
		}),
		ProposalVotingTime: promauto.NewHistogram(prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: "charitabletrust",
			Name:      "proposal_voting_time_hours",
			Help:      "Time taken for proposal voting in hours",
			Buckets:   prometheus.LinearBuckets(0, 24, 8), // 0 to 7 days
		}),
		ProposalApprovalRate: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "charitabletrust",
			Name:      "proposal_approval_rate",
			Help:      "Percentage of proposals approved",
		}),
		VotesPerProposal: promauto.NewHistogram(prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: "charitabletrust",
			Name:      "votes_per_proposal",
			Help:      "Distribution of votes per proposal",
			Buckets:   prometheus.LinearBuckets(0, 1, 8), // 0 to 7 votes
		}),

		// Impact metrics
		ImpactReportsSubmitted: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "charitabletrust",
			Name:      "impact_reports_submitted",
			Help:      "Total number of impact reports submitted",
		}),
		ImpactReportsVerified: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "charitabletrust",
			Name:      "impact_reports_verified",
			Help:      "Total number of impact reports verified",
		}),
		BeneficiariesReached: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "charitabletrust",
			Name:      "beneficiaries_reached_total",
			Help:      "Total number of beneficiaries reached",
		}),
		FundsUtilizationRate: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "charitabletrust",
			Name:      "funds_utilization_rate",
			Help:      "Average funds utilization rate across all allocations",
		}),
		ImpactScoreDistribution: promauto.NewHistogram(prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: "charitabletrust",
			Name:      "impact_score_distribution",
			Help:      "Distribution of impact scores",
			Buckets:   prometheus.LinearBuckets(0, 10, 11), // 0 to 100 in steps of 10
		}),
		VerificationCompletionTime: promauto.NewHistogram(prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: "charitabletrust",
			Name:      "verification_completion_time_days",
			Help:      "Time taken to verify impact reports in days",
			Buckets:   prometheus.LinearBuckets(0, 5, 12), // 0 to 60 days
		}),

		// Fraud metrics
		FraudAlertsCreated: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "charitabletrust",
			Name:      "fraud_alerts_created",
			Help:      "Total number of fraud alerts created",
		}),
		FraudAlertsInvestigated: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "charitabletrust",
			Name:      "fraud_alerts_investigated",
			Help:      "Total number of fraud alerts investigated",
		}),
		FraudAlertsSeverity: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "charitabletrust",
			Name:      "fraud_alerts_by_severity",
			Help:      "Fraud alerts by severity level",
		}, []string{"severity"}),
		FraudInvestigationTime: promauto.NewHistogram(prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: "charitabletrust",
			Name:      "fraud_investigation_time_days",
			Help:      "Time taken to investigate fraud alerts in days",
			Buckets:   prometheus.LinearBuckets(0, 2, 15), // 0 to 30 days
		}),
		FraudConfirmedCases: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "charitabletrust",
			Name:      "fraud_confirmed_cases",
			Help:      "Total number of confirmed fraud cases",
		}),

		// Governance metrics
		TrusteesActive: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "charitabletrust",
			Name:      "trustees_active",
			Help:      "Number of active trustees",
		}),
		TrusteeVotingRate: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "charitabletrust",
			Name:      "trustee_voting_rate",
			Help:      "Average trustee voting participation rate",
		}),
		GovernanceQuorum: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "charitabletrust",
			Name:      "governance_quorum",
			Help:      "Current governance quorum requirement",
		}),
		AdvisoryCommitteeSize: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "charitabletrust",
			Name:      "advisory_committee_size",
			Help:      "Number of advisory committee members",
		}),

		// Performance metrics
		TransactionDuration: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: "charitabletrust",
			Name:      "transaction_duration_ms",
			Help:      "Transaction processing duration in milliseconds",
			Buckets:   prometheus.ExponentialBuckets(1, 2, 10), // 1ms to 1s
		}, []string{"msg_type"}),
		QueryDuration: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: "charitabletrust",
			Name:      "query_duration_ms",
			Help:      "Query processing duration in milliseconds",
			Buckets:   prometheus.ExponentialBuckets(0.1, 2, 10), // 0.1ms to 100ms
		}, []string{"query_type"}),
		ValidationErrors: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "charitabletrust",
			Name:      "validation_errors_total",
			Help:      "Total number of validation errors by type",
		}, []string{"error_type"}),

		// Transparency metrics
		TransparencyScore: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "charitabletrust",
			Name:      "transparency_score",
			Help:      "Overall transparency score (0-100)",
		}),
		PublicReportsPublished: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "charitabletrust",
			Name:      "public_reports_published",
			Help:      "Total number of public reports published",
		}),
		StakeholderEngagementScore: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "charitabletrust",
			Name:      "stakeholder_engagement_score",
			Help:      "Stakeholder engagement score (0-100)",
		}),
		ComplianceViolations: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "charitabletrust",
			Name:      "compliance_violations_total",
			Help:      "Total number of compliance violations",
		}),
	}
}

// UpdateTrustFundMetrics updates all trust fund related metrics
func (k Keeper) UpdateTrustFundMetrics(ctx sdk.Context) {
	if k.metrics == nil {
		return
	}

	// Update trust fund balance
	balance, found := k.GetTrustFundBalance(ctx)
	if found {
		k.metrics.TrustFundBalance.Set(float64(balance.TotalBalance.Amount.Int64()))
		k.metrics.TrustAllocatedAmount.Set(float64(balance.AllocatedAmount.Amount.Int64()))
		k.metrics.TrustAvailableAmount.Set(float64(balance.AvailableAmount.Amount.Int64()))
		k.metrics.TrustTotalDistributed.Set(float64(balance.TotalDistributed.Amount.Int64()))
	}

	// Update governance metrics
	governance, found := k.GetTrustGovernance(ctx)
	if found {
		activeTrustees := 0
		for _, trustee := range governance.Trustees {
			if trustee.Status == "active" {
				activeTrustees++
			}
		}
		k.metrics.TrusteesActive.Set(float64(activeTrustees))
		k.metrics.GovernanceQuorum.Set(float64(governance.Quorum))
		k.metrics.AdvisoryCommitteeSize.Set(float64(len(governance.AdvisoryCommittee)))
	}

	// Count allocations by category
	categoryCount := make(map[string]int)
	activeCount := 0
	organizationMap := make(map[uint64]bool)

	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.CharitableAllocationKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var allocation types.CharitableAllocation
		k.cdc.MustUnmarshal(iterator.Value(), &allocation)
		
		categoryCount[allocation.Category]++
		organizationMap[allocation.CharitableOrgWalletId] = true
		
		if allocation.Status == "active" {
			activeCount++
		}
	}

	// Update category metrics
	for category, count := range categoryCount {
		k.metrics.AllocationsByCategory.WithLabelValues(category).Set(float64(count))
	}

	k.metrics.CharitableAllocationsActive.Set(float64(activeCount))
	k.metrics.OrganizationCount.Set(float64(len(organizationMap)))
}

// RecordAllocationMetrics records metrics for a new allocation
func (k Keeper) RecordAllocationMetrics(ctx sdk.Context, allocation types.CharitableAllocation) {
	if k.metrics == nil {
		return
	}

	k.metrics.CharitableAllocationsTotal.Inc()
	k.metrics.CharitableAllocationAmounts.Observe(float64(allocation.Amount.Amount.Int64()))

	if allocation.Status == "distributed" {
		k.metrics.CharitableAllocationsDistributed.Inc()
	}
}

// RecordProposalMetrics records proposal-related metrics
func (k Keeper) RecordProposalMetrics(ctx sdk.Context, proposal types.AllocationProposal, votingStarted, votingEnded time.Time) {
	if k.metrics == nil {
		return
	}

	k.metrics.ProposalsCreated.Inc()
	k.metrics.VotesPerProposal.Observe(float64(len(proposal.Votes)))

	if !votingEnded.IsZero() {
		votingTime := votingEnded.Sub(votingStarted)
		k.metrics.ProposalVotingTime.Observe(votingTime.Hours())
	}

	switch proposal.Status {
	case "approved":
		k.metrics.ProposalsApproved.Inc()
	case "rejected":
		k.metrics.ProposalsRejected.Inc()
	case "executed":
		k.metrics.ProposalsExecuted.Inc()
	}

	// Update approval rate
	k.updateApprovalRate(ctx)
}

// RecordImpactReportMetrics records impact report metrics
func (k Keeper) RecordImpactReportMetrics(ctx sdk.Context, report types.ImpactReport, impactScore int) {
	if k.metrics == nil {
		return
	}

	k.metrics.ImpactReportsSubmitted.Inc()
	k.metrics.BeneficiariesReached.Add(float64(report.BeneficiariesReached))
	k.metrics.ImpactScoreDistribution.Observe(float64(impactScore))

	if report.Verification != nil && report.Verification.IsVerified {
		k.metrics.ImpactReportsVerified.Inc()
		
		verificationTime := report.Verification.VerifiedAt.Sub(report.SubmittedAt)
		k.metrics.VerificationCompletionTime.Observe(verificationTime.Hours() / 24)
	}

	// Update funds utilization rate
	k.updateFundsUtilizationRate(ctx)
}

// RecordFraudAlertMetrics records fraud alert metrics
func (k Keeper) RecordFraudAlertMetrics(ctx sdk.Context, alert types.FraudAlert) {
	if k.metrics == nil {
		return
	}

	k.metrics.FraudAlertsCreated.Inc()
	k.metrics.FraudAlertsSeverity.WithLabelValues(alert.Severity).Inc()

	if alert.Investigation != nil && !alert.Investigation.CompletedAt.IsZero() {
		k.metrics.FraudAlertsInvestigated.Inc()
		
		investigationTime := alert.Investigation.CompletedAt.Sub(alert.Investigation.StartedAt)
		k.metrics.FraudInvestigationTime.Observe(investigationTime.Hours() / 24)

		if alert.Investigation.Recommendation == "confirmed_fraud" {
			k.metrics.FraudConfirmedCases.Inc()
		}
	}
}

// RecordTransactionMetrics records transaction processing metrics
func (k Keeper) RecordTransactionMetrics(msgType string, duration time.Duration, err error) {
	if k.metrics == nil {
		return
	}

	k.metrics.TransactionDuration.WithLabelValues(msgType).Observe(float64(duration.Milliseconds()))
	
	if err != nil {
		errorType := "unknown"
		switch {
		case types.ErrInsufficientFunds.Is(err):
			errorType = "insufficient_funds"
		case types.ErrNotTrustee.Is(err):
			errorType = "not_trustee"
		case types.ErrDuplicateVote.Is(err):
			errorType = "duplicate_vote"
		case types.ErrVotingPeriodExpired.Is(err):
			errorType = "voting_expired"
		}
		k.metrics.ValidationErrors.WithLabelValues(errorType).Inc()
	}
}

// RecordQueryMetrics records query processing metrics
func (k Keeper) RecordQueryMetrics(queryType string, duration time.Duration) {
	if k.metrics == nil {
		return
	}

	k.metrics.QueryDuration.WithLabelValues(queryType).Observe(float64(duration.Microseconds()) / 1000)
}

// updateApprovalRate calculates and updates the proposal approval rate
func (k Keeper) updateApprovalRate(ctx sdk.Context) {
	totalProposals := 0
	approvedProposals := 0

	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.AllocationProposalKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var proposal types.AllocationProposal
		k.cdc.MustUnmarshal(iterator.Value(), &proposal)
		
		if proposal.Status != "pending" {
			totalProposals++
			if proposal.Status == "approved" || proposal.Status == "executed" {
				approvedProposals++
			}
		}
	}

	if totalProposals > 0 {
		approvalRate := float64(approvedProposals) / float64(totalProposals) * 100
		k.metrics.ProposalApprovalRate.Set(approvalRate)
	}
}

// updateFundsUtilizationRate calculates average funds utilization across allocations
func (k Keeper) updateFundsUtilizationRate(ctx sdk.Context) {
	totalAllocated := sdk.ZeroInt()
	totalUtilized := sdk.ZeroInt()
	reportCount := 0

	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ImpactReportKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var report types.ImpactReport
		k.cdc.MustUnmarshal(iterator.Value(), &report)
		
		allocation, found := k.GetCharitableAllocation(ctx, report.AllocationId)
		if found {
			totalAllocated = totalAllocated.Add(allocation.Amount.Amount)
			totalUtilized = totalUtilized.Add(report.FundsUtilized.Amount)
			reportCount++
		}
	}

	if !totalAllocated.IsZero() {
		utilizationRate := sdk.NewDecFromInt(totalUtilized).Quo(sdk.NewDecFromInt(totalAllocated)).MulInt64(100)
		k.metrics.FundsUtilizationRate.Set(utilizationRate.MustFloat64())
	}
}

// UpdateTransparencyScore calculates and updates the transparency score
func (k Keeper) UpdateTransparencyScore(ctx sdk.Context) {
	if k.metrics == nil {
		return
	}

	score := float64(0)
	factors := 0

	// Factor 1: Impact report submission rate
	allocations := k.GetAllCharitableAllocations(ctx)
	if len(allocations) > 0 {
		reportsSubmitted := 0
		for _, allocation := range allocations {
			if allocation.Status == "distributed" {
				reports := k.GetImpactReportsByAllocation(ctx, allocation.Id)
				if len(reports) > 0 {
					reportsSubmitted++
				}
			}
		}
		reportRate := float64(reportsSubmitted) / float64(len(allocations)) * 100
		score += reportRate
		factors++
	}

	// Factor 2: Verification rate
	totalReports := k.GetTotalImpactReports(ctx)
	verifiedReports := k.GetVerifiedImpactReports(ctx)
	if totalReports > 0 {
		verificationRate := float64(verifiedReports) / float64(totalReports) * 100
		score += verificationRate
		factors++
	}

	// Factor 3: Trustee voting participation
	governance, found := k.GetTrustGovernance(ctx)
	if found {
		avgVotingRate := k.calculateAverageVotingRate(ctx, governance)
		score += avgVotingRate
		factors++
		k.metrics.TrusteeVotingRate.Set(avgVotingRate)
	}

	// Calculate final score
	if factors > 0 {
		finalScore := score / float64(factors)
		k.metrics.TransparencyScore.Set(finalScore)
	}
}

// calculateAverageVotingRate calculates the average voting participation rate
func (k Keeper) calculateAverageVotingRate(ctx sdk.Context, governance types.TrustGovernance) float64 {
	totalVotes := 0
	totalPossibleVotes := 0

	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.AllocationProposalKey)
	defer iterator.Close()

	proposalCount := 0
	for ; iterator.Valid(); iterator.Next() {
		var proposal types.AllocationProposal
		k.cdc.MustUnmarshal(iterator.Value(), &proposal)
		
		if proposal.Status != "pending" {
			proposalCount++
			totalVotes += len(proposal.Votes)
			totalPossibleVotes += len(governance.Trustees)
		}
	}

	if totalPossibleVotes > 0 {
		return float64(totalVotes) / float64(totalPossibleVotes) * 100
	}
	return 0
}