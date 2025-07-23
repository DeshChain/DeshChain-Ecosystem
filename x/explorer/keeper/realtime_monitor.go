package keeper

import (
	"context"
	"fmt"
	"sync"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/namo/x/explorer/types"
)

// RealTimeMonitor handles real-time blockchain monitoring and notifications
type RealTimeMonitor struct {
	keeper        Keeper
	subscribers   map[string]*Subscriber
	watchlists    map[string]*Watchlist
	alerts        map[string]*Alert
	metrics       *RealtimeMetrics
	mu            sync.RWMutex
	ctx           context.Context
	cancel        context.CancelFunc
	eventChannels map[string]chan types.MonitorEvent
}

// NewRealTimeMonitor creates a new real-time monitor
func NewRealTimeMonitor(keeper Keeper) *RealTimeMonitor {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &RealTimeMonitor{
		keeper:        keeper,
		subscribers:   make(map[string]*Subscriber),
		watchlists:    make(map[string]*Watchlist),
		alerts:        make(map[string]*Alert),
		metrics:       NewRealtimeMetrics(),
		eventChannels: make(map[string]chan types.MonitorEvent),
		ctx:           ctx,
		cancel:        cancel,
	}
}

// Subscriber represents a real-time data subscriber
type Subscriber struct {
	ID              string                 `json:"id"`
	Address         string                 `json:"address"`
	SubscriptionType string                `json:"subscription_type"` // blocks, transactions, addresses, tokens
	Filters         types.MonitorFilters   `json:"filters"`
	WebhookURL      string                 `json:"webhook_url,omitempty"`
	EventChannel    chan types.MonitorEvent `json:"-"`
	LastSeen        time.Time              `json:"last_seen"`
	IsActive        bool                   `json:"is_active"`
	CreatedAt       time.Time              `json:"created_at"`
}

// Watchlist represents a set of addresses/transactions to monitor
type Watchlist struct {
	ID          string                  `json:"id"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Owner       string                  `json:"owner"`
	Addresses   []types.WatchedAddress  `json:"addresses"`
	Tokens      []types.WatchedToken    `json:"tokens"`
	Thresholds  types.AlertThresholds   `json:"thresholds"`
	Notifications types.NotificationSettings `json:"notifications"`
	IsPublic    bool                    `json:"is_public"`
	CreatedAt   time.Time               `json:"created_at"`
	UpdatedAt   time.Time               `json:"updated_at"`
}

// Alert represents a triggered alert
type Alert struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"` // threshold, pattern, anomaly
	Severity    string                 `json:"severity"` // low, medium, high, critical
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	TriggerData types.AlertTriggerData `json:"trigger_data"`
	WatchlistID string                 `json:"watchlist_id,omitempty"`
	Status      string                 `json:"status"` // active, acknowledged, resolved
	CreatedAt   time.Time              `json:"created_at"`
	ResolvedAt  *time.Time             `json:"resolved_at,omitempty"`
}

// RealtimeMetrics tracks real-time monitoring metrics
type RealtimeMetrics struct {
	TotalSubscribers    int64     `json:"total_subscribers"`
	ActiveSubscribers   int64     `json:"active_subscribers"`
	TotalWatchlists     int64     `json:"total_watchlists"`
	ActiveAlerts        int64     `json:"active_alerts"`
	EventsProcessed     int64     `json:"events_processed"`
	LastBlockProcessed  int64     `json:"last_block_processed"`
	ProcessingLatency   time.Duration `json:"processing_latency"`
	EventsPerSecond     float64   `json:"events_per_second"`
	LastUpdated         time.Time `json:"last_updated"`
	mu                  sync.RWMutex
}

// Start begins real-time monitoring
func (rtm *RealTimeMonitor) Start(ctx sdk.Context) error {
	rtm.keeper.Logger(ctx).Info("Starting real-time monitor")
	
	// Start block monitoring goroutine
	go rtm.monitorBlocks()
	
	// Start transaction monitoring goroutine
	go rtm.monitorTransactions()
	
	// Start alert processing goroutine
	go rtm.processAlerts()
	
	// Start metrics collection goroutine
	go rtm.collectMetrics()
	
	// Start cleanup goroutine
	go rtm.cleanup()
	
	return nil
}

// Stop halts real-time monitoring
func (rtm *RealTimeMonitor) Stop() {
	rtm.cancel()
	
	// Close all event channels
	rtm.mu.Lock()
	for _, channel := range rtm.eventChannels {
		close(channel)
	}
	rtm.mu.Unlock()
}

// SubscribeToBlocks subscribes to real-time block updates
func (rtm *RealTimeMonitor) SubscribeToBlocks(ctx sdk.Context, subscriberAddr string, filters types.MonitorFilters) (string, chan types.MonitorEvent, error) {
	subscriberID := rtm.generateSubscriberID(ctx, subscriberAddr, "blocks")
	
	eventChannel := make(chan types.MonitorEvent, 1000) // Buffered channel
	
	subscriber := &Subscriber{
		ID:               subscriberID,
		Address:          subscriberAddr,
		SubscriptionType: "blocks",
		Filters:          filters,
		EventChannel:     eventChannel,
		LastSeen:         ctx.BlockTime(),
		IsActive:         true,
		CreatedAt:        ctx.BlockTime(),
	}
	
	rtm.mu.Lock()
	rtm.subscribers[subscriberID] = subscriber
	rtm.eventChannels[subscriberID] = eventChannel
	rtm.mu.Unlock()
	
	rtm.metrics.mu.Lock()
	rtm.metrics.TotalSubscribers++
	rtm.metrics.ActiveSubscribers++
	rtm.metrics.mu.Unlock()
	
	// Emit subscription event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeBlockSubscription,
			sdk.NewAttribute(types.AttributeKeySubscriberID, subscriberID),
			sdk.NewAttribute(types.AttributeKeySubscriberAddress, subscriberAddr),
		),
	)
	
	return subscriberID, eventChannel, nil
}

// SubscribeToTransactions subscribes to real-time transaction updates
func (rtm *RealTimeMonitor) SubscribeToTransactions(ctx sdk.Context, subscriberAddr string, filters types.MonitorFilters) (string, chan types.MonitorEvent, error) {
	subscriberID := rtm.generateSubscriberID(ctx, subscriberAddr, "transactions")
	
	eventChannel := make(chan types.MonitorEvent, 1000)
	
	subscriber := &Subscriber{
		ID:               subscriberID,
		Address:          subscriberAddr,
		SubscriptionType: "transactions",
		Filters:          filters,
		EventChannel:     eventChannel,
		LastSeen:         ctx.BlockTime(),
		IsActive:         true,
		CreatedAt:        ctx.BlockTime(),
	}
	
	rtm.mu.Lock()
	rtm.subscribers[subscriberID] = subscriber
	rtm.eventChannels[subscriberID] = eventChannel
	rtm.mu.Unlock()
	
	rtm.metrics.mu.Lock()
	rtm.metrics.TotalSubscribers++
	rtm.metrics.ActiveSubscribers++
	rtm.metrics.mu.Unlock()
	
	return subscriberID, eventChannel, nil
}

// CreateWatchlist creates a new address/token watchlist
func (rtm *RealTimeMonitor) CreateWatchlist(ctx sdk.Context, request types.CreateWatchlistRequest) (*Watchlist, error) {
	watchlistID := rtm.generateWatchlistID(ctx, request.Owner)
	
	watchlist := &Watchlist{
		ID:          watchlistID,
		Name:        request.Name,
		Description: request.Description,
		Owner:       request.Owner,
		Addresses:   request.Addresses,
		Tokens:      request.Tokens,
		Thresholds:  request.Thresholds,
		Notifications: request.Notifications,
		IsPublic:    request.IsPublic,
		CreatedAt:   ctx.BlockTime(),
		UpdatedAt:   ctx.BlockTime(),
	}
	
	rtm.mu.Lock()
	rtm.watchlists[watchlistID] = watchlist
	rtm.mu.Unlock()
	
	rtm.metrics.mu.Lock()
	rtm.metrics.TotalWatchlists++
	rtm.metrics.mu.Unlock()
	
	// Store watchlist
	rtm.keeper.SetWatchlist(ctx, *watchlist)
	
	// Emit watchlist creation event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeWatchlistCreated,
			sdk.NewAttribute(types.AttributeKeyWatchlistID, watchlistID),
			sdk.NewAttribute(types.AttributeKeyWatchlistOwner, request.Owner),
			sdk.NewAttribute(types.AttributeKeyWatchlistName, request.Name),
		),
	)
	
	return watchlist, nil
}

// monitorBlocks monitors new blocks and notifies subscribers
func (rtm *RealTimeMonitor) monitorBlocks() {
	ticker := time.NewTicker(1 * time.Second) // Check for new blocks every second
	defer ticker.Stop()
	
	var lastProcessedHeight int64 = 0
	
	for {
		select {
		case <-rtm.ctx.Done():
			return
		case <-ticker.C:
			// Get latest block height
			latestHeight := rtm.keeper.GetLatestBlockHeight()
			
			// Process new blocks
			for height := lastProcessedHeight + 1; height <= latestHeight; height++ {
				if block := rtm.keeper.GetBlockByHeight(sdk.Context{}, height); block != nil {
					rtm.processNewBlock(*block)
					lastProcessedHeight = height
					
					rtm.metrics.mu.Lock()
					rtm.metrics.LastBlockProcessed = height
					rtm.metrics.EventsProcessed++
					rtm.metrics.mu.Unlock()
				}
			}
		}
	}
}

// monitorTransactions monitors new transactions and notifies subscribers
func (rtm *RealTimeMonitor) monitorTransactions() {
	ticker := time.NewTicker(500 * time.Millisecond) // Check for new transactions more frequently
	defer ticker.Stop()
	
	var lastProcessedTxIndex int64 = 0
	
	for {
		select {
		case <-rtm.ctx.Done():
			return
		case <-ticker.C:
			// Get latest transactions
			newTxs := rtm.keeper.GetTransactionsAfterIndex(lastProcessedTxIndex)
			
			for _, tx := range newTxs {
				rtm.processNewTransaction(tx)
				lastProcessedTxIndex = tx.Index
				
				rtm.metrics.mu.Lock()
				rtm.metrics.EventsProcessed++
				rtm.metrics.mu.Unlock()
			}
		}
	}
}

// processNewBlock processes a new block and notifies relevant subscribers
func (rtm *RealTimeMonitor) processNewBlock(block types.Block) {
	startTime := time.Now()
	
	event := types.MonitorEvent{
		Type:      "new_block",
		BlockData: &block,
		Timestamp: time.Now(),
	}
	
	rtm.mu.RLock()
	defer rtm.mu.RUnlock()
	
	// Notify block subscribers
	for _, subscriber := range rtm.subscribers {
		if subscriber.SubscriptionType == "blocks" && subscriber.IsActive {
			if rtm.blockMatchesFilters(block, subscriber.Filters) {
				select {
				case subscriber.EventChannel <- event:
					// Successfully sent
				default:
					// Channel full, subscriber may be slow
					rtm.keeper.Logger(sdk.Context{}).Warn("Subscriber channel full", "subscriber_id", subscriber.ID)
				}
			}
		}
	}
	
	// Check watchlist triggers
	rtm.checkWatchlistTriggers(block)
	
	// Update processing latency
	rtm.metrics.mu.Lock()
	rtm.metrics.ProcessingLatency = time.Since(startTime)
	rtm.metrics.mu.Unlock()
}

// processNewTransaction processes a new transaction and notifies relevant subscribers
func (rtm *RealTimeMonitor) processNewTransaction(tx types.Transaction) {
	event := types.MonitorEvent{
		Type:            "new_transaction",
		TransactionData: &tx,
		Timestamp:       time.Now(),
	}
	
	rtm.mu.RLock()
	defer rtm.mu.RUnlock()
	
	// Notify transaction subscribers
	for _, subscriber := range rtm.subscribers {
		if subscriber.SubscriptionType == "transactions" && subscriber.IsActive {
			if rtm.transactionMatchesFilters(tx, subscriber.Filters) {
				select {
				case subscriber.EventChannel <- event:
					// Successfully sent
				default:
					// Channel full
					rtm.keeper.Logger(sdk.Context{}).Warn("Subscriber channel full", "subscriber_id", subscriber.ID)
				}
			}
		}
	}
	
	// Check for address-specific notifications
	rtm.checkAddressNotifications(tx)
	
	// Check for large transaction alerts
	rtm.checkLargeTransactionAlerts(tx)
}

// processAlerts processes and manages alerts
func (rtm *RealTimeMonitor) processAlerts() {
	ticker := time.NewTicker(5 * time.Second) // Process alerts every 5 seconds
	defer ticker.Stop()
	
	for {
		select {
		case <-rtm.ctx.Done():
			return
		case <-ticker.C:
			rtm.processActiveAlerts()
		}
	}
}

// collectMetrics collects and updates real-time metrics
func (rtm *RealTimeMonitor) collectMetrics() {
	ticker := time.NewTicker(10 * time.Second) // Update metrics every 10 seconds
	defer ticker.Stop()
	
	for {
		select {
		case <-rtm.ctx.Done():
			return
		case <-ticker.C:
			rtm.updateMetrics()
		}
	}
}

// cleanup performs periodic cleanup of inactive subscribers and old data
func (rtm *RealTimeMonitor) cleanup() {
	ticker := time.NewTicker(1 * time.Minute) // Cleanup every minute
	defer ticker.Stop()
	
	for {
		select {
		case <-rtm.ctx.Done():
			return
		case <-ticker.C:
			rtm.cleanupInactiveSubscribers()
			rtm.cleanupOldAlerts()
		}
	}
}

// Helper functions

func (rtm *RealTimeMonitor) blockMatchesFilters(block types.Block, filters types.MonitorFilters) bool {
	// Check height range
	if filters.MinHeight > 0 && block.Height < filters.MinHeight {
		return false
	}
	if filters.MaxHeight > 0 && block.Height > filters.MaxHeight {
		return false
	}
	
	// Check transaction count
	if filters.MinTxCount > 0 && int64(len(block.Transactions)) < filters.MinTxCount {
		return false
	}
	if filters.MaxTxCount > 0 && int64(len(block.Transactions)) > filters.MaxTxCount {
		return false
	}
	
	// Check proposer address
	if len(filters.ProposerAddresses) > 0 {
		found := false
		for _, addr := range filters.ProposerAddresses {
			if block.ProposerAddress == addr {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	
	return true
}

func (rtm *RealTimeMonitor) transactionMatchesFilters(tx types.Transaction, filters types.MonitorFilters) bool {
	// Check transaction types
	if len(filters.TxTypes) > 0 {
		found := false
		for _, txType := range filters.TxTypes {
			if tx.Type == txType {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	
	// Check addresses
	if len(filters.Addresses) > 0 {
		found := false
		for _, addr := range filters.Addresses {
			if rtm.transactionInvolvesAddress(tx, addr) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	
	// Check amount range
	if filters.MinAmount != nil {
		totalAmount := rtm.calculateTransactionTotal(tx)
		if totalAmount.LT(*filters.MinAmount) {
			return false
		}
	}
	if filters.MaxAmount != nil {
		totalAmount := rtm.calculateTransactionTotal(tx)
		if totalAmount.GT(*filters.MaxAmount) {
			return false
		}
	}
	
	// Check gas range
	if filters.MinGas > 0 && tx.GasUsed < filters.MinGas {
		return false
	}
	if filters.MaxGas > 0 && tx.GasUsed > filters.MaxGas {
		return false
	}
	
	return true
}

func (rtm *RealTimeMonitor) checkWatchlistTriggers(block types.Block) {
	for _, watchlist := range rtm.watchlists {
		// Check if any watched addresses are involved in block transactions
		for _, tx := range block.Transactions {
			for _, watchedAddr := range watchlist.Addresses {
				if rtm.transactionInvolvesAddress(tx, watchedAddr.Address) {
					rtm.triggerWatchlistAlert(watchlist, "address_activity", tx)
				}
			}
		}
	}
}

func (rtm *RealTimeMonitor) checkAddressNotifications(tx types.Transaction) {
	for _, watchlist := range rtm.watchlists {
		for _, watchedAddr := range watchlist.Addresses {
			if rtm.transactionInvolvesAddress(tx, watchedAddr.Address) {
				// Check thresholds
				if rtm.checkAddressThresholds(tx, watchedAddr, watchlist.Thresholds) {
					rtm.triggerAlert("threshold_exceeded", "high", fmt.Sprintf("Address %s exceeded threshold", watchedAddr.Address), tx)
				}
			}
		}
	}
}

func (rtm *RealTimeMonitor) checkLargeTransactionAlerts(tx types.Transaction) {
	totalAmount := rtm.calculateTransactionTotal(tx)
	largeTransactionThreshold := sdk.NewInt(1000000) // 10 lakh threshold
	
	if totalAmount.GT(largeTransactionThreshold) {
		rtm.triggerAlert("large_transaction", "medium", fmt.Sprintf("Large transaction detected: %s", totalAmount.String()), tx)
	}
}

func (rtm *RealTimeMonitor) triggerWatchlistAlert(watchlist *Watchlist, alertType string, tx types.Transaction) {
	alert := &Alert{
		ID:          rtm.generateAlertID(),
		Type:        alertType,
		Severity:    "medium",
		Title:       fmt.Sprintf("Watchlist Alert: %s", watchlist.Name),
		Description: fmt.Sprintf("Activity detected for watchlist %s", watchlist.Name),
		TriggerData: types.AlertTriggerData{
			TransactionHash: tx.Hash,
			BlockHeight:     tx.Height,
			Amount:          rtm.calculateTransactionTotal(tx).String(),
		},
		WatchlistID: watchlist.ID,
		Status:      "active",
		CreatedAt:   time.Now(),
	}
	
	rtm.mu.Lock()
	rtm.alerts[alert.ID] = alert
	rtm.mu.Unlock()
	
	rtm.metrics.mu.Lock()
	rtm.metrics.ActiveAlerts++
	rtm.metrics.mu.Unlock()
}

func (rtm *RealTimeMonitor) triggerAlert(alertType, severity, description string, tx types.Transaction) {
	alert := &Alert{
		ID:          rtm.generateAlertID(),
		Type:        alertType,
		Severity:    severity,
		Title:       fmt.Sprintf("Transaction Alert: %s", alertType),
		Description: description,
		TriggerData: types.AlertTriggerData{
			TransactionHash: tx.Hash,
			BlockHeight:     tx.Height,
			Amount:          rtm.calculateTransactionTotal(tx).String(),
		},
		Status:    "active",
		CreatedAt: time.Now(),
	}
	
	rtm.mu.Lock()
	rtm.alerts[alert.ID] = alert
	rtm.mu.Unlock()
	
	rtm.metrics.mu.Lock()
	rtm.metrics.ActiveAlerts++
	rtm.metrics.mu.Unlock()
}

func (rtm *RealTimeMonitor) processActiveAlerts() {
	rtm.mu.RLock()
	defer rtm.mu.RUnlock()
	
	for _, alert := range rtm.alerts {
		if alert.Status == "active" {
			// Send notifications for active alerts
			rtm.sendAlertNotifications(alert)
		}
	}
}

func (rtm *RealTimeMonitor) sendAlertNotifications(alert *Alert) {
	// Implementation would send notifications via webhook, email, etc.
	// For now, just log the alert
	rtm.keeper.Logger(sdk.Context{}).Info("Alert triggered", 
		"alert_id", alert.ID, 
		"type", alert.Type, 
		"severity", alert.Severity,
		"description", alert.Description)
}

func (rtm *RealTimeMonitor) updateMetrics() {
	rtm.metrics.mu.Lock()
	defer rtm.metrics.mu.Unlock()
	
	// Count active subscribers
	activeCount := int64(0)
	rtm.mu.RLock()
	for _, subscriber := range rtm.subscribers {
		if subscriber.IsActive {
			activeCount++
		}
	}
	rtm.mu.RUnlock()
	
	rtm.metrics.ActiveSubscribers = activeCount
	rtm.metrics.LastUpdated = time.Now()
	
	// Calculate events per second
	if rtm.metrics.LastUpdated.Sub(rtm.metrics.LastUpdated) > 0 {
		duration := time.Since(rtm.metrics.LastUpdated).Seconds()
		rtm.metrics.EventsPerSecond = float64(rtm.metrics.EventsProcessed) / duration
	}
}

func (rtm *RealTimeMonitor) cleanupInactiveSubscribers() {
	rtm.mu.Lock()
	defer rtm.mu.Unlock()
	
	inactiveThreshold := time.Now().Add(-30 * time.Minute) // 30 minutes
	
	for id, subscriber := range rtm.subscribers {
		if subscriber.LastSeen.Before(inactiveThreshold) {
			subscriber.IsActive = false
			close(subscriber.EventChannel)
			delete(rtm.eventChannels, id)
			
			rtm.metrics.mu.Lock()
			rtm.metrics.ActiveSubscribers--
			rtm.metrics.mu.Unlock()
		}
	}
}

func (rtm *RealTimeMonitor) cleanupOldAlerts() {
	rtm.mu.Lock()
	defer rtm.mu.Unlock()
	
	oldThreshold := time.Now().Add(-24 * time.Hour) // 24 hours
	
	for id, alert := range rtm.alerts {
		if alert.CreatedAt.Before(oldThreshold) && alert.Status == "resolved" {
			delete(rtm.alerts, id)
		}
	}
}

// Utility functions

func (rtm *RealTimeMonitor) generateSubscriberID(ctx sdk.Context, addr, subType string) string {
	timestamp := ctx.BlockTime().Unix()
	return fmt.Sprintf("SUB-%s-%s-%d", addr[:8], subType, timestamp)
}

func (rtm *RealTimeMonitor) generateWatchlistID(ctx sdk.Context, owner string) string {
	timestamp := ctx.BlockTime().Unix()
	return fmt.Sprintf("WL-%s-%d", owner[:8], timestamp)
}

func (rtm *RealTimeMonitor) generateAlertID() string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("ALERT-%d", timestamp)
}

func (rtm *RealTimeMonitor) transactionInvolvesAddress(tx types.Transaction, address string) bool {
	if tx.Sender == address {
		return true
	}
	
	// Check in message recipients
	for _, msg := range tx.Messages {
		if msg.To == address || msg.From == address {
			return true
		}
	}
	
	return false
}

func (rtm *RealTimeMonitor) calculateTransactionTotal(tx types.Transaction) sdk.Int {
	total := sdk.ZeroInt()
	for _, msg := range tx.Messages {
		if msg.Amount.Amount.IsPositive() {
			total = total.Add(msg.Amount.Amount)
		}
	}
	return total
}

func (rtm *RealTimeMonitor) checkAddressThresholds(tx types.Transaction, watchedAddr types.WatchedAddress, thresholds types.AlertThresholds) bool {
	if thresholds.MaxTransactionAmount != nil {
		txAmount := rtm.calculateTransactionTotal(tx)
		if txAmount.GT(*thresholds.MaxTransactionAmount) {
			return true
		}
	}
	
	// Add more threshold checks as needed
	
	return false
}

func NewRealtimeMetrics() *RealtimeMetrics {
	return &RealtimeMetrics{
		LastUpdated: time.Now(),
	}
}

// GetMetrics returns current real-time metrics
func (rtm *RealTimeMonitor) GetMetrics() *RealtimeMetrics {
	rtm.metrics.mu.RLock()
	defer rtm.metrics.mu.RUnlock()
	
	// Return a copy to avoid race conditions
	return &RealtimeMetrics{
		TotalSubscribers:   rtm.metrics.TotalSubscribers,
		ActiveSubscribers:  rtm.metrics.ActiveSubscribers,
		TotalWatchlists:    rtm.metrics.TotalWatchlists,
		ActiveAlerts:       rtm.metrics.ActiveAlerts,
		EventsProcessed:    rtm.metrics.EventsProcessed,
		LastBlockProcessed: rtm.metrics.LastBlockProcessed,
		ProcessingLatency:  rtm.metrics.ProcessingLatency,
		EventsPerSecond:    rtm.metrics.EventsPerSecond,
		LastUpdated:        rtm.metrics.LastUpdated,
	}
}