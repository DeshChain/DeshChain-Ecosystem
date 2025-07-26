package keeper

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Metrics contains all prometheus metrics for the DSWF module
type Metrics struct {
	// Fund metrics
	FundTotalBalance     prometheus.Gauge
	FundAllocatedAmount  prometheus.Gauge
	FundAvailableAmount  prometheus.Gauge
	FundInvestedAmount   prometheus.Gauge
	FundTotalReturns     prometheus.Gauge
	FundReturnRate       prometheus.Gauge

	// Allocation metrics
	AllocationsTotal     prometheus.Counter
	AllocationsActive    prometheus.Gauge
	AllocationsCompleted prometheus.Counter
	AllocationsRejected  prometheus.Counter
	AllocationAmounts    prometheus.Histogram

	// Disbursement metrics
	DisbursementsTotal     prometheus.Counter
	DisbursementAmounts    prometheus.Histogram
	DisbursementDelays     prometheus.Histogram
	DisbursementsScheduled prometheus.Gauge

	// Portfolio metrics
	PortfolioRiskScore      prometheus.Gauge
	PortfolioComponents     prometheus.GaugeVec
	PortfolioRebalances     prometheus.Counter
	PortfolioReturnsByAsset prometheus.GaugeVec

	// Governance metrics
	ProposalsTotal       prometheus.Counter
	ProposalVotingTime   prometheus.Histogram
	GovernanceQuorum     prometheus.Gauge
	FundManagersActive   prometheus.Gauge
	MultiSigValidations  prometheus.Counter

	// Performance metrics
	TransactionDuration prometheus.HistogramVec
	QueryDuration       prometheus.HistogramVec
	ValidationErrors    prometheus.CounterVec

	// Business metrics
	MonthlyReportsSubmitted prometheus.Counter
	AuditScheduleAdherence  prometheus.Gauge
	InvestmentHorizon       prometheus.Gauge
}

// NewMetrics creates a new metrics instance
func NewMetrics(namespace string) *Metrics {
	return &Metrics{
		// Fund metrics
		FundTotalBalance: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "dswf",
			Name:      "fund_total_balance",
			Help:      "Total balance in the DSWF in unamo",
		}),
		FundAllocatedAmount: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "dswf",
			Name:      "fund_allocated_amount",
			Help:      "Total allocated amount in unamo",
		}),
		FundAvailableAmount: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "dswf",
			Name:      "fund_available_amount",
			Help:      "Available amount for new allocations in unamo",
		}),
		FundInvestedAmount: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "dswf",
			Name:      "fund_invested_amount",
			Help:      "Total amount currently invested in unamo",
		}),
		FundTotalReturns: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "dswf",
			Name:      "fund_total_returns",
			Help:      "Total returns generated in unamo",
		}),
		FundReturnRate: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "dswf",
			Name:      "fund_return_rate",
			Help:      "Annual return rate as a percentage",
		}),

		// Allocation metrics
		AllocationsTotal: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "dswf",
			Name:      "allocations_total",
			Help:      "Total number of allocations created",
		}),
		AllocationsActive: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "dswf",
			Name:      "allocations_active",
			Help:      "Number of currently active allocations",
		}),
		AllocationsCompleted: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "dswf",
			Name:      "allocations_completed",
			Help:      "Total number of completed allocations",
		}),
		AllocationsRejected: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "dswf",
			Name:      "allocations_rejected",
			Help:      "Total number of rejected allocations",
		}),
		AllocationAmounts: promauto.NewHistogram(prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: "dswf",
			Name:      "allocation_amounts",
			Help:      "Distribution of allocation amounts in unamo",
			Buckets:   prometheus.ExponentialBuckets(1000000, 10, 8), // 1 NAMO to 10M NAMO
		}),

		// Disbursement metrics
		DisbursementsTotal: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "dswf",
			Name:      "disbursements_total",
			Help:      "Total number of disbursements executed",
		}),
		DisbursementAmounts: promauto.NewHistogram(prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: "dswf",
			Name:      "disbursement_amounts",
			Help:      "Distribution of disbursement amounts in unamo",
			Buckets:   prometheus.ExponentialBuckets(100000, 10, 8), // 0.1 NAMO to 1M NAMO
		}),
		DisbursementDelays: promauto.NewHistogram(prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: "dswf",
			Name:      "disbursement_delays_days",
			Help:      "Delays in disbursement execution in days",
			Buckets:   prometheus.LinearBuckets(0, 1, 30), // 0 to 30 days
		}),
		DisbursementsScheduled: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "dswf",
			Name:      "disbursements_scheduled",
			Help:      "Number of pending scheduled disbursements",
		}),

		// Portfolio metrics
		PortfolioRiskScore: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "dswf",
			Name:      "portfolio_risk_score",
			Help:      "Current portfolio risk score (1-10)",
		}),
		PortfolioComponents: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "dswf",
			Name:      "portfolio_component_value",
			Help:      "Value of each portfolio component in unamo",
		}, []string{"asset_type", "risk_rating"}),
		PortfolioRebalances: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "dswf",
			Name:      "portfolio_rebalances_total",
			Help:      "Total number of portfolio rebalances",
		}),
		PortfolioReturnsByAsset: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "dswf",
			Name:      "portfolio_returns_by_asset",
			Help:      "Returns by asset type as percentage",
		}, []string{"asset_type"}),

		// Governance metrics
		ProposalsTotal: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "dswf",
			Name:      "proposals_total",
			Help:      "Total number of allocation proposals",
		}),
		ProposalVotingTime: promauto.NewHistogram(prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: "dswf",
			Name:      "proposal_voting_time_hours",
			Help:      "Time taken for proposal approval in hours",
			Buckets:   prometheus.LinearBuckets(0, 6, 20), // 0 to 120 hours
		}),
		GovernanceQuorum: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "dswf",
			Name:      "governance_quorum",
			Help:      "Current governance quorum requirement",
		}),
		FundManagersActive: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "dswf",
			Name:      "fund_managers_active",
			Help:      "Number of active fund managers",
		}),
		MultiSigValidations: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "dswf",
			Name:      "multisig_validations_total",
			Help:      "Total number of multi-signature validations",
		}),

		// Performance metrics
		TransactionDuration: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: "dswf",
			Name:      "transaction_duration_ms",
			Help:      "Transaction processing duration in milliseconds",
			Buckets:   prometheus.ExponentialBuckets(1, 2, 10), // 1ms to 1s
		}, []string{"msg_type"}),
		QueryDuration: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: "dswf",
			Name:      "query_duration_ms",
			Help:      "Query processing duration in milliseconds",
			Buckets:   prometheus.ExponentialBuckets(0.1, 2, 10), // 0.1ms to 100ms
		}, []string{"query_type"}),
		ValidationErrors: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "dswf",
			Name:      "validation_errors_total",
			Help:      "Total number of validation errors by type",
		}, []string{"error_type"}),

		// Business metrics
		MonthlyReportsSubmitted: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "dswf",
			Name:      "monthly_reports_submitted",
			Help:      "Total number of monthly reports submitted",
		}),
		AuditScheduleAdherence: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "dswf",
			Name:      "audit_schedule_adherence",
			Help:      "Percentage of audits completed on schedule",
		}),
		InvestmentHorizon: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "dswf",
			Name:      "investment_horizon_days",
			Help:      "Current investment horizon in days",
		}),
	}
}

// UpdateFundMetrics updates all fund-related metrics
func (k Keeper) UpdateFundMetrics(ctx sdk.Context) {
	if k.metrics == nil {
		return
	}

	// Get fund balance
	balance := k.GetFundBalance(ctx)
	if len(balance) > 0 {
		k.metrics.FundTotalBalance.Set(float64(balance[0].Amount.Int64()))
	}

	// Get portfolio
	portfolio, found := k.GetInvestmentPortfolio(ctx)
	if found {
		k.metrics.FundAllocatedAmount.Set(float64(portfolio.AllocatedAmount.Amount.Int64()))
		k.metrics.FundAvailableAmount.Set(float64(portfolio.AvailableAmount.Amount.Int64()))
		k.metrics.FundInvestedAmount.Set(float64(portfolio.InvestedAmount.Amount.Int64()))
		k.metrics.FundTotalReturns.Set(float64(portfolio.TotalReturns.Amount.Int64()))
		k.metrics.FundReturnRate.Set(portfolio.AnnualReturnRate.MustFloat64() * 100)
		k.metrics.PortfolioRiskScore.Set(float64(portfolio.RiskScore))

		// Update portfolio components
		for _, component := range portfolio.Components {
			k.metrics.PortfolioComponents.WithLabelValues(
				component.AssetType,
				component.RiskRating,
			).Set(float64(component.CurrentValue.Amount.Int64()))

			k.metrics.PortfolioReturnsByAsset.WithLabelValues(
				component.AssetType,
			).Set(component.ReturnRate.MustFloat64() * 100)
		}
	}

	// Update governance metrics
	governance, found := k.GetFundGovernance(ctx)
	if found {
		k.metrics.FundManagersActive.Set(float64(len(governance.FundManagers)))
		k.metrics.GovernanceQuorum.Set(float64(governance.RequiredSignatures))
	}

	// Count active allocations
	activeCount := 0
	scheduledDisbursements := 0
	
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.FundAllocationKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var allocation types.FundAllocation
		k.cdc.MustUnmarshal(iterator.Value(), &allocation)
		
		if allocation.Status == "active" {
			activeCount++
			
			// Count scheduled disbursements
			for _, disbursement := range allocation.Disbursements {
				if disbursement.Status == "pending" {
					scheduledDisbursements++
				}
			}
		}
	}

	k.metrics.AllocationsActive.Set(float64(activeCount))
	k.metrics.DisbursementsScheduled.Set(float64(scheduledDisbursements))
}

// RecordAllocationMetrics records metrics for a new allocation
func (k Keeper) RecordAllocationMetrics(ctx sdk.Context, allocation types.FundAllocation) {
	if k.metrics == nil {
		return
	}

	k.metrics.AllocationsTotal.Inc()
	k.metrics.AllocationAmounts.Observe(float64(allocation.Amount.Amount.Int64()))

	// Record by status
	switch allocation.Status {
	case "completed":
		k.metrics.AllocationsCompleted.Inc()
	case "rejected":
		k.metrics.AllocationsRejected.Inc()
	}
}

// RecordDisbursementMetrics records metrics for a disbursement
func (k Keeper) RecordDisbursementMetrics(ctx sdk.Context, disbursement types.Disbursement, scheduledDate time.Time) {
	if k.metrics == nil {
		return
	}

	k.metrics.DisbursementsTotal.Inc()
	k.metrics.DisbursementAmounts.Observe(float64(disbursement.Amount.Amount.Int64()))

	// Calculate delay
	delay := ctx.BlockTime().Sub(scheduledDate)
	delayDays := delay.Hours() / 24
	k.metrics.DisbursementDelays.Observe(delayDays)
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
		case sdk.ErrInsufficientFunds.Is(err):
			errorType = "insufficient_funds"
		case types.ErrInvalidAmount.Is(err):
			errorType = "invalid_amount"
		case types.ErrUnauthorized.Is(err):
			errorType = "unauthorized"
		case types.ErrInvalidCategory.Is(err):
			errorType = "invalid_category"
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

// RecordProposalMetrics records proposal-related metrics
func (k Keeper) RecordProposalMetrics(ctx sdk.Context, proposalID uint64, created, approved time.Time) {
	if k.metrics == nil {
		return
	}

	k.metrics.ProposalsTotal.Inc()
	
	if !approved.IsZero() {
		votingTime := approved.Sub(created)
		k.metrics.ProposalVotingTime.Observe(votingTime.Hours())
	}
}

// RecordPortfolioRebalance records portfolio rebalancing metrics
func (k Keeper) RecordPortfolioRebalance(ctx sdk.Context) {
	if k.metrics == nil {
		return
	}

	k.metrics.PortfolioRebalances.Inc()
}

// RecordMonthlyReport records monthly report submission
func (k Keeper) RecordMonthlyReport(ctx sdk.Context) {
	if k.metrics == nil {
		return
	}

	k.metrics.MonthlyReportsSubmitted.Inc()
}

// RecordMultiSigValidation records multi-signature validation attempts
func (k Keeper) RecordMultiSigValidation(ctx sdk.Context, success bool) {
	if k.metrics == nil {
		return
	}

	k.metrics.MultiSigValidations.Inc()
}