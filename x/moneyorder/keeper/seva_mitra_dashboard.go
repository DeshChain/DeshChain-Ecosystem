package keeper

import (
	"fmt"
	"sort"
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/deshchain/x/moneyorder/types"
)

// SevaMitraDashboard manages dashboard functionality for Seva Mitras
type SevaMitraDashboard struct {
	keeper *Keeper
}

// NewSevaMitraDashboard creates a new Seva Mitra dashboard manager
func NewSevaMitraDashboard(keeper *Keeper) *SevaMitraDashboard {
	return &SevaMitraDashboard{
		keeper: keeper,
	}
}

// GetDashboardData returns comprehensive dashboard data for a Seva Mitra
func (smd *SevaMitraDashboard) GetDashboardData(
	ctx sdk.Context,
	mitraAddress string,
) (*types.SevaMitraDashboardData, error) {
	// Verify Seva Mitra exists and is active
	mitra, found := smd.keeper.GetSevaMitra(ctx, mitraAddress)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrNotFound, "seva mitra not found")
	}

	if !mitra.IsActive {
		return nil, sdkerrors.Wrap(types.ErrUserNotActive, "seva mitra not active")
	}

	// Gather dashboard data
	dashboardData := &types.SevaMitraDashboardData{
		MitraInfo:        mitra,
		Summary:          smd.getDashboardSummary(ctx, mitraAddress),
		EarningsData:     smd.getEarningsData(ctx, mitraAddress),
		ServiceRequests:  smd.getRecentServiceRequests(ctx, mitraAddress, 10),
		PerformanceStats: smd.getPerformanceStats(ctx, mitraAddress),
		Analytics:        smd.getAnalytics(ctx, mitraAddress),
		Notifications:    smd.getNotifications(ctx, mitraAddress),
		Rankings:         smd.getRankings(ctx, mitraAddress),
	}

	return dashboardData, nil
}

// getDashboardSummary returns key summary metrics
func (smd *SevaMitraDashboard) getDashboardSummary(
	ctx sdk.Context,
	mitraAddress string,
) types.DashboardSummary {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	thisMonthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	return types.DashboardSummary{
		TotalServices:     smd.getTotalServicesCount(ctx, mitraAddress),
		TodayServices:     smd.getServicesCount(ctx, mitraAddress, todayStart, now),
		MonthlyServices:   smd.getServicesCount(ctx, mitraAddress, thisMonthStart, now),
		TotalEarnings:     smd.getTotalEarnings(ctx, mitraAddress),
		TodayEarnings:     smd.getEarnings(ctx, mitraAddress, todayStart, now),
		MonthlyEarnings:   smd.getEarnings(ctx, mitraAddress, thisMonthStart, now),
		PendingRequests:   smd.getPendingRequestsCount(ctx, mitraAddress),
		TrustScore:        smd.getCurrentTrustScore(ctx, mitraAddress),
		ResponseTime:      smd.getAverageResponseTime(ctx, mitraAddress),
		CustomerRating:    smd.getAverageCustomerRating(ctx, mitraAddress),
		OnlineStatus:      smd.isOnline(ctx, mitraAddress),
		LastActiveTime:    smd.getLastActiveTime(ctx, mitraAddress),
	}
}

// getEarningsData returns detailed earnings breakdown
func (smd *SevaMitraDashboard) getEarningsData(
	ctx sdk.Context,
	mitraAddress string,
) types.EarningsData {
	now := time.Now()
	
	return types.EarningsData{
		TotalEarnings:     smd.getTotalEarnings(ctx, mitraAddress),
		WeeklyEarnings:    smd.getWeeklyEarnings(ctx, mitraAddress),
		MonthlyEarnings:   smd.getMonthlyEarnings(ctx, mitraAddress),
		ServiceTypeBreakdown: smd.getEarningsByServiceType(ctx, mitraAddress),
		RecentTransactions:   smd.getRecentEarningsTransactions(ctx, mitraAddress, 20),
		PendingPayments:      smd.getPendingPayments(ctx, mitraAddress),
		PaymentHistory:       smd.getPaymentHistory(ctx, mitraAddress, 30), // Last 30 days
	}
}

// getRecentServiceRequests returns recent service requests
func (smd *SevaMitraDashboard) getRecentServiceRequests(
	ctx sdk.Context,
	mitraAddress string,
	limit int,
) []types.ServiceRequest {
	var requests []types.ServiceRequest
	
	store := prefix.NewStore(ctx.KVStore(smd.keeper.storeKey), types.ServiceRequestPrefix)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var request types.ServiceRequest
		smd.keeper.cdc.MustUnmarshal(iterator.Value(), &request)
		
		if request.SevaMitraAddress == mitraAddress {
			requests = append(requests, request)
		}
	}

	// Sort by created time (most recent first)
	sort.Slice(requests, func(i, j int) bool {
		return requests[i].CreatedAt.After(requests[j].CreatedAt)
	})

	// Return limited results
	if len(requests) > limit {
		requests = requests[:limit]
	}

	return requests
}

// getPerformanceStats returns performance statistics
func (smd *SevaMitraDashboard) getPerformanceStats(
	ctx sdk.Context,
	mitraAddress string,
) types.PerformanceStats {
	now := time.Now()
	thisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	lastMonth := thisMonth.AddDate(0, -1, 0)
	
	return types.PerformanceStats{
		CompletionRate:        smd.getCompletionRate(ctx, mitraAddress),
		AverageResponseTime:   smd.getAverageResponseTime(ctx, mitraAddress),
		CustomerSatisfaction:  smd.getAverageCustomerRating(ctx, mitraAddress),
		RepeatedCustomers:     smd.getRepeatedCustomersCount(ctx, mitraAddress),
		DisputeRate:          smd.getDisputeRate(ctx, mitraAddress),
		OnTimeDelivery:       smd.getOnTimeDeliveryRate(ctx, mitraAddress),
		MonthlyGrowth:        smd.calculateGrowthRate(ctx, mitraAddress, lastMonth, thisMonth),
		ServiceReliability:   smd.getServiceReliabilityScore(ctx, mitraAddress),
		TrustScoreHistory:    smd.getTrustScoreHistory(ctx, mitraAddress, 30), // Last 30 days
	}
}

// getAnalytics returns analytical data
func (smd *SevaMitraDashboard) getAnalytics(
	ctx sdk.Context,
	mitraAddress string,
) types.AnalyticsData {
	return types.AnalyticsData{
		ServiceTypesChart:     smd.getServiceTypesDistribution(ctx, mitraAddress),
		HourlyActivity:        smd.getHourlyActivityPattern(ctx, mitraAddress),
		GeographicDistribution: smd.getGeographicDistribution(ctx, mitraAddress),
		CustomerDemographics:  smd.getCustomerDemographics(ctx, mitraAddress),
		SeasonalTrends:        smd.getSeasonalTrends(ctx, mitraAddress),
		CompetitorAnalysis:    smd.getCompetitorAnalysis(ctx, mitraAddress),
		MarketOpportunities:   smd.getMarketOpportunities(ctx, mitraAddress),
	}
}

// getNotifications returns important notifications for the Seva Mitra
func (smd *SevaMitraDashboard) getNotifications(
	ctx sdk.Context,
	mitraAddress string,
) []types.Notification {
	var notifications []types.Notification
	
	// Check for urgent service requests
	urgentRequests := smd.getUrgentServiceRequests(ctx, mitraAddress)
	for _, request := range urgentRequests {
		notifications = append(notifications, types.Notification{
			ID:       fmt.Sprintf("urgent_%s", request.ID),
			Type:     types.NotificationType_URGENT_REQUEST,
			Title:    "Urgent Service Request",
			Message:  fmt.Sprintf("New urgent request for %s service", request.ServiceType),
			Priority: types.NotificationPriority_HIGH,
			CreatedAt: request.CreatedAt,
			Data:     map[string]interface{}{"request_id": request.ID},
		})
	}
	
	// Check for low ratings alerts
	if smd.hasRecentLowRatings(ctx, mitraAddress) {
		notifications = append(notifications, types.Notification{
			ID:       "low_rating_alert",
			Type:     types.NotificationType_RATING_ALERT,
			Title:    "Rating Alert",
			Message:  "Your recent ratings are below average. Consider improving service quality.",
			Priority: types.NotificationPriority_MEDIUM,
			CreatedAt: time.Now(),
		})
	}
	
	// Check for payment notifications
	pendingPayments := smd.getPendingPayments(ctx, mitraAddress)
	if len(pendingPayments) > 0 {
		totalPending := sdk.ZeroInt()
		for _, payment := range pendingPayments {
			totalPending = totalPending.Add(payment.Amount)
		}
		
		notifications = append(notifications, types.Notification{
			ID:       "pending_payments",
			Type:     types.NotificationType_PAYMENT,
			Title:    "Pending Payments",
			Message:  fmt.Sprintf("You have ₹%s in pending payments", totalPending.String()),
			Priority: types.NotificationPriority_MEDIUM,
			CreatedAt: time.Now(),
			Data:     map[string]interface{}{"amount": totalPending.String()},
		})
	}
	
	// Check for performance milestones
	if smd.hasReachedMilestone(ctx, mitraAddress) {
		notifications = append(notifications, types.Notification{
			ID:       "milestone_achieved",
			Type:     types.NotificationType_ACHIEVEMENT,
			Title:    "Milestone Achieved!",
			Message:  "Congratulations! You've reached a new service milestone.",
			Priority: types.NotificationPriority_LOW,
			CreatedAt: time.Now(),
		})
	}
	
	return notifications
}

// getRankings returns ranking information
func (smd *SevaMitraDashboard) getRankings(
	ctx sdk.Context,
	mitraAddress string,
) types.RankingData {
	return types.RankingData{
		LocalRank:     smd.getLocalRanking(ctx, mitraAddress),
		RegionalRank:  smd.getRegionalRanking(ctx, mitraAddress),
		NationalRank:  smd.getNationalRanking(ctx, mitraAddress),
		CategoryRank:  smd.getCategoryRanking(ctx, mitraAddress),
		TrustScoreRank: smd.getTrustScoreRanking(ctx, mitraAddress),
		EarningsRank:  smd.getEarningsRanking(ctx, mitraAddress),
		TotalMitras:   smd.getTotalMitrasCount(ctx),
	}
}

// UpdateOnlineStatus updates the online status of a Seva Mitra
func (smd *SevaMitraDashboard) UpdateOnlineStatus(
	ctx sdk.Context,
	mitraAddress string,
	isOnline bool,
) error {
	mitra, found := smd.keeper.GetSevaMitra(ctx, mitraAddress)
	if !found {
		return sdkerrors.Wrap(types.ErrNotFound, "seva mitra not found")
	}

	mitra.IsOnline = isOnline
	mitra.LastActiveTime = time.Now()
	
	if isOnline {
		mitra.OnlineSince = time.Now()
	}

	smd.keeper.SetSevaMitra(ctx, mitra)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSevaMitraStatusUpdate,
			sdk.NewAttribute(types.AttributeKeySevaMitraAddress, mitraAddress),
			sdk.NewAttribute(types.AttributeKeyOnlineStatus, fmt.Sprintf("%t", isOnline)),
			sdk.NewAttribute(types.AttributeKeyTimestamp, time.Now().Format(time.RFC3339)),
		),
	)

	return nil
}

// UpdateServiceAvailability updates service availability
func (smd *SevaMitraDashboard) UpdateServiceAvailability(
	ctx sdk.Context,
	mitraAddress string,
	serviceType string,
	available bool,
) error {
	mitra, found := smd.keeper.GetSevaMitra(ctx, mitraAddress)
	if !found {
		return sdkerrors.Wrap(types.ErrNotFound, "seva mitra not found")
	}

	// Update service availability
	if mitra.ServiceAvailability == nil {
		mitra.ServiceAvailability = make(map[string]bool)
	}
	mitra.ServiceAvailability[serviceType] = available

	smd.keeper.SetSevaMitra(ctx, mitra)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSevaMitraServiceUpdate,
			sdk.NewAttribute(types.AttributeKeySevaMitraAddress, mitraAddress),
			sdk.NewAttribute(types.AttributeKeyServiceType, serviceType),
			sdk.NewAttribute(types.AttributeKeyAvailable, fmt.Sprintf("%t", available)),
		),
	)

	return nil
}

// AcceptServiceRequest allows Seva Mitra to accept a service request
func (smd *SevaMitraDashboard) AcceptServiceRequest(
	ctx sdk.Context,
	mitraAddress string,
	requestID string,
) error {
	// Get service request
	request, found := smd.getServiceRequest(ctx, requestID)
	if !found {
		return sdkerrors.Wrap(types.ErrNotFound, "service request not found")
	}

	// Verify Seva Mitra is assigned to this request
	if request.SevaMitraAddress != mitraAddress {
		return sdkerrors.Wrap(types.ErrUnauthorized, "seva mitra not assigned to this request")
	}

	// Check if request is in pending state
	if request.Status != types.ServiceRequestStatus_PENDING {
		return sdkerrors.Wrap(types.ErrInvalidOrderStatus, "request is not in pending state")
	}

	// Update request status
	request.Status = types.ServiceRequestStatus_ACCEPTED
	request.AcceptedAt = time.Now()
	request.EstimatedCompletionTime = smd.calculateEstimatedCompletion(request.ServiceType)

	// Save updated request
	smd.saveServiceRequest(ctx, request)

	// Update Seva Mitra metrics
	smd.updateMitraMetrics(ctx, mitraAddress, "request_accepted")

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeServiceRequestAccepted,
			sdk.NewAttribute(types.AttributeKeyRequestID, requestID),
			sdk.NewAttribute(types.AttributeKeySevaMitraAddress, mitraAddress),
			sdk.NewAttribute(types.AttributeKeyTimestamp, time.Now().Format(time.RFC3339)),
		),
	)

	return nil
}

// CompleteServiceRequest marks a service as completed
func (smd *SevaMitraDashboard) CompleteServiceRequest(
	ctx sdk.Context,
	mitraAddress string,
	requestID string,
	completionNote string,
) error {
	// Get service request
	request, found := smd.getServiceRequest(ctx, requestID)
	if !found {
		return sdkerrors.Wrap(types.ErrNotFound, "service request not found")
	}

	// Verify Seva Mitra is assigned to this request
	if request.SevaMitraAddress != mitraAddress {
		return sdkerrors.Wrap(types.ErrUnauthorized, "seva mitra not assigned to this request")
	}

	// Check if request is in progress
	if request.Status != types.ServiceRequestStatus_IN_PROGRESS && request.Status != types.ServiceRequestStatus_ACCEPTED {
		return sdkerrors.Wrap(types.ErrInvalidOrderStatus, "request is not in progress")
	}

	// Update request status
	request.Status = types.ServiceRequestStatus_COMPLETED
	request.CompletedAt = time.Now()
	request.CompletionNote = completionNote

	// Calculate service duration
	if !request.AcceptedAt.IsZero() {
		request.ServiceDuration = time.Since(request.AcceptedAt)
	}

	// Calculate earnings
	earnings := smd.calculateServiceEarnings(request)
	request.EarningsAmount = earnings

	// Save updated request
	smd.saveServiceRequest(ctx, request)

	// Update Seva Mitra earnings and metrics
	smd.updateMitraEarnings(ctx, mitraAddress, earnings)
	smd.updateMitraMetrics(ctx, mitraAddress, "service_completed")

	// Create earnings transaction
	smd.createEarningsTransaction(ctx, mitraAddress, requestID, earnings)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeServiceRequestCompleted,
			sdk.NewAttribute(types.AttributeKeyRequestID, requestID),
			sdk.NewAttribute(types.AttributeKeySevaMitraAddress, mitraAddress),
			sdk.NewAttribute(types.AttributeKeyEarnings, earnings.String()),
			sdk.NewAttribute(types.AttributeKeyTimestamp, time.Now().Format(time.RFC3339)),
		),
	)

	return nil
}

// Helper functions (implementations would be provided based on specific requirements)

func (smd *SevaMitraDashboard) getTotalServicesCount(ctx sdk.Context, mitraAddress string) int64 {
	// Implementation to count total services by this Seva Mitra
	return 0 // Placeholder
}

func (smd *SevaMitraDashboard) getServicesCount(ctx sdk.Context, mitraAddress string, start, end time.Time) int64 {
	// Implementation to count services in time range
	return 0 // Placeholder
}

func (smd *SevaMitraDashboard) getTotalEarnings(ctx sdk.Context, mitraAddress string) sdk.Int {
	// Implementation to get total earnings
	return sdk.ZeroInt() // Placeholder
}

func (smd *SevaMitraDashboard) getEarnings(ctx sdk.Context, mitraAddress string, start, end time.Time) sdk.Int {
	// Implementation to get earnings in time range
	return sdk.ZeroInt() // Placeholder
}

func (smd *SevaMitraDashboard) getPendingRequestsCount(ctx sdk.Context, mitraAddress string) int64 {
	// Implementation to count pending requests
	return 0 // Placeholder
}

func (smd *SevaMitraDashboard) getCurrentTrustScore(ctx sdk.Context, mitraAddress string) float64 {
	mitra, found := smd.keeper.GetSevaMitra(ctx, mitraAddress)
	if !found {
		return 0.0
	}
	return mitra.TrustScore
}

func (smd *SevaMitraDashboard) getAverageResponseTime(ctx sdk.Context, mitraAddress string) time.Duration {
	// Implementation to calculate average response time
	return time.Minute * 5 // Placeholder
}

func (smd *SevaMitraDashboard) getAverageCustomerRating(ctx sdk.Context, mitraAddress string) float64 {
	// Implementation to calculate average customer rating
	return 4.5 // Placeholder
}

func (smd *SevaMitraDashboard) isOnline(ctx sdk.Context, mitraAddress string) bool {
	mitra, found := smd.keeper.GetSevaMitra(ctx, mitraAddress)
	if !found {
		return false
	}
	return mitra.IsOnline
}

func (smd *SevaMitraDashboard) getLastActiveTime(ctx sdk.Context, mitraAddress string) time.Time {
	mitra, found := smd.keeper.GetSevaMitra(ctx, mitraAddress)
	if !found {
		return time.Time{}
	}
	return mitra.LastActiveTime
}

// Additional helper function implementations would continue here...
// This is a comprehensive framework for the Seva Mitra dashboard