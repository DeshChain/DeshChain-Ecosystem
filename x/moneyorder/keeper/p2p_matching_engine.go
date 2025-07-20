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
	"math"
	"sort"
	"strings"
	"sync"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/deshchain/x/moneyorder/types"
)

// Enhanced P2P Matching Engine with real-time order book management

// OrderBook maintains active P2P orders for efficient matching
type OrderBook struct {
	sync.RWMutex
	buyOrders  map[string]*types.P2POrder  // orderID -> order
	sellOrders map[string]*types.P2POrder  // orderID -> order
	
	// Indexes for efficient queries
	postalIndex    map[string][]string      // postalCode -> []orderID
	districtIndex  map[string][]string      // district -> []orderID
	stateIndex     map[string][]string      // state -> []orderID
	currencyIndex  map[string][]string      // currency -> []orderID
	
	// Price indexes for market data
	buyPriceIndex  map[string]*PriceLevel   // currency -> price levels
	sellPriceIndex map[string]*PriceLevel   // currency -> price levels
}

// PriceLevel represents orders at a specific price point
type PriceLevel struct {
	sync.RWMutex
	levels map[string][]*types.P2POrder // price -> orders
}

// MatchingEngine handles P2P order matching with advanced algorithms
type MatchingEngine struct {
	keeper    *Keeper
	orderBook *OrderBook
	
	// Matching parameters
	maxDistanceKm      int32
	minTrustScore      int32
	maxPriceDeviation  sdk.Dec
	languageBonus      float64
	paymentMethodBonus float64
	
	// Performance metrics
	matchCount         int64
	avgMatchTime       time.Duration
	lastMatchTimestamp time.Time
}

// NewMatchingEngine creates a new P2P matching engine
func (k *Keeper) NewMatchingEngine() *MatchingEngine {
	return &MatchingEngine{
		keeper: k,
		orderBook: &OrderBook{
			buyOrders:      make(map[string]*types.P2POrder),
			sellOrders:     make(map[string]*types.P2POrder),
			postalIndex:    make(map[string][]string),
			districtIndex:  make(map[string][]string),
			stateIndex:     make(map[string][]string),
			currencyIndex:  make(map[string][]string),
			buyPriceIndex:  make(map[string]*PriceLevel),
			sellPriceIndex: make(map[string]*PriceLevel),
		},
		maxDistanceKm:     50,
		minTrustScore:     50,
		maxPriceDeviation: sdk.NewDecWithPrec(5, 2), // 5%
		languageBonus:     10.0,
		paymentMethodBonus: 15.0,
	}
}

// AddOrderToBook adds a new order to the order book with indexing
func (me *MatchingEngine) AddOrderToBook(order *types.P2POrder) {
	me.orderBook.Lock()
	defer me.orderBook.Unlock()
	
	// Add to main storage
	if order.OrderType == types.OrderType_BUY_NAMO {
		me.orderBook.buyOrders[order.OrderId] = order
	} else {
		me.orderBook.sellOrders[order.OrderId] = order
	}
	
	// Update indexes
	me.updateIndexes(order, true)
	
	// Update price levels
	me.updatePriceLevels(order, true)
}

// RemoveOrderFromBook removes an order from the order book
func (me *MatchingEngine) RemoveOrderFromBook(orderID string) {
	me.orderBook.Lock()
	defer me.orderBook.Unlock()
	
	var order *types.P2POrder
	var found bool
	
	// Find and remove from main storage
	if order, found = me.orderBook.buyOrders[orderID]; found {
		delete(me.orderBook.buyOrders, orderID)
	} else if order, found = me.orderBook.sellOrders[orderID]; found {
		delete(me.orderBook.sellOrders, orderID)
	}
	
	if found && order != nil {
		// Update indexes
		me.updateIndexes(order, false)
		
		// Update price levels
		me.updatePriceLevels(order, false)
	}
}

// updateIndexes updates all indexes for an order
func (me *MatchingEngine) updateIndexes(order *types.P2POrder, add bool) {
	indexes := []struct {
		index map[string][]string
		key   string
	}{
		{me.orderBook.postalIndex, order.PostalCode},
		{me.orderBook.districtIndex, order.District},
		{me.orderBook.stateIndex, order.State},
		{me.orderBook.currencyIndex, order.FiatCurrency},
	}
	
	for _, idx := range indexes {
		if add {
			idx.index[idx.key] = append(idx.index[idx.key], order.OrderId)
		} else {
			// Remove from index
			orderIDs := idx.index[idx.key]
			for i, id := range orderIDs {
				if id == order.OrderId {
					idx.index[idx.key] = append(orderIDs[:i], orderIDs[i+1:]...)
					break
				}
			}
		}
	}
}

// updatePriceLevels updates price level indexes
func (me *MatchingEngine) updatePriceLevels(order *types.P2POrder, add bool) {
	price := me.calculatePrice(order)
	priceStr := price.String()
	
	var priceIndex *PriceLevel
	if order.OrderType == types.OrderType_BUY_NAMO {
		if me.orderBook.buyPriceIndex[order.FiatCurrency] == nil {
			me.orderBook.buyPriceIndex[order.FiatCurrency] = &PriceLevel{
				levels: make(map[string][]*types.P2POrder),
			}
		}
		priceIndex = me.orderBook.buyPriceIndex[order.FiatCurrency]
	} else {
		if me.orderBook.sellPriceIndex[order.FiatCurrency] == nil {
			me.orderBook.sellPriceIndex[order.FiatCurrency] = &PriceLevel{
				levels: make(map[string][]*types.P2POrder),
			}
		}
		priceIndex = me.orderBook.sellPriceIndex[order.FiatCurrency]
	}
	
	priceIndex.Lock()
	defer priceIndex.Unlock()
	
	if add {
		priceIndex.levels[priceStr] = append(priceIndex.levels[priceStr], order)
	} else {
		// Remove from price level
		orders := priceIndex.levels[priceStr]
		for i, o := range orders {
			if o.OrderId == order.OrderId {
				priceIndex.levels[priceStr] = append(orders[:i], orders[i+1:]...)
				break
			}
		}
	}
}

// FindBestMatches finds the best matching orders using advanced algorithm
func (me *MatchingEngine) FindBestMatches(ctx sdk.Context, order *types.P2POrder, limit int) []*MatchResult {
	startTime := time.Now()
	defer func() {
		me.avgMatchTime = time.Since(startTime)
		me.lastMatchTimestamp = time.Now()
	}()
	
	// Get candidate orders based on location proximity
	candidates := me.getCandidateOrders(order)
	
	// Score and rank candidates
	var matches []*MatchResult
	for _, candidate := range candidates {
		if score := me.calculateAdvancedMatchScore(ctx, order, candidate); score > 0 {
			matches = append(matches, &MatchResult{
				Order:         candidate,
				Score:         score,
				Distance:      me.calculateDistance(order.PostalCode, candidate.PostalCode),
				PriceMatch:    me.calculatePriceMatch(order, candidate),
				TrustScore:    me.keeper.GetUserTrustScore(candidate.Creator),
			})
		}
	}
	
	// Sort by score
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Score > matches[j].Score
	})
	
	// Return top matches
	if len(matches) > limit {
		matches = matches[:limit]
	}
	
	return matches
}

// getCandidateOrders retrieves orders that could potentially match
func (me *MatchingEngine) getCandidateOrders(order *types.P2POrder) []*types.P2POrder {
	me.orderBook.RLock()
	defer me.orderBook.RUnlock()
	
	var candidates []*types.P2POrder
	candidateMap := make(map[string]bool)
	
	// Determine which order map to search
	var orderMap map[string]*types.P2POrder
	if order.OrderType == types.OrderType_BUY_NAMO {
		orderMap = me.orderBook.sellOrders
	} else {
		orderMap = me.orderBook.buyOrders
	}
	
	// Search strategies in order of preference
	// 1. Same postal code
	for _, orderID := range me.orderBook.postalIndex[order.PostalCode] {
		if candidate, found := orderMap[orderID]; found && !candidateMap[orderID] {
			candidates = append(candidates, candidate)
			candidateMap[orderID] = true
		}
	}
	
	// 2. Same district (if not enough matches)
	if len(candidates) < 10 {
		for _, orderID := range me.orderBook.districtIndex[order.District] {
			if candidate, found := orderMap[orderID]; found && !candidateMap[orderID] {
				candidates = append(candidates, candidate)
				candidateMap[orderID] = true
			}
		}
	}
	
	// 3. Same state (if still not enough matches)
	if len(candidates) < 20 {
		for _, orderID := range me.orderBook.stateIndex[order.State] {
			if candidate, found := orderMap[orderID]; found && !candidateMap[orderID] {
				// Check distance constraint
				if me.calculateDistance(order.PostalCode, candidate.PostalCode) <= float64(order.MaxDistanceKm) {
					candidates = append(candidates, candidate)
					candidateMap[orderID] = true
				}
			}
		}
	}
	
	return candidates
}

// calculateAdvancedMatchScore calculates a comprehensive match score
func (me *MatchingEngine) calculateAdvancedMatchScore(ctx sdk.Context, order1, order2 *types.P2POrder) float64 {
	// Base compatibility check
	if !me.areOrdersCompatible(order1, order2) {
		return 0
	}
	
	score := 100.0
	
	// 1. Distance Score (0-30 points)
	distance := me.calculateDistance(order1.PostalCode, order2.PostalCode)
	distanceScore := me.calculateDistanceScore(distance, float64(order1.MaxDistanceKm))
	score += distanceScore * 30
	
	// 2. Price Match Score (0-25 points)
	priceScore := me.calculatePriceMatch(order1, order2)
	score += priceScore * 25
	
	// 3. Trust Score (0-20 points)
	user1Trust := me.keeper.GetUserTrustScore(order1.Creator)
	user2Trust := me.keeper.GetUserTrustScore(order2.Creator)
	avgTrust := float64(user1Trust+user2Trust) / 2.0
	score += (avgTrust / 100.0) * 20
	
	// 4. Payment Method Score (0-15 points)
	paymentScore := me.calculatePaymentMethodScore(order1.PaymentMethods, order2.PaymentMethods)
	score += paymentScore * 15
	
	// 5. Language Score (0-10 points)
	if me.hasCommonLanguage(order1.PreferredLanguages, order2.PreferredLanguages) {
		score += me.languageBonus
	}
	
	// 6. Time Decay Bonus (older orders get priority)
	ageHours := time.Since(order2.CreatedAt).Hours()
	timeBonus := math.Min(ageHours*0.5, 10) // Max 10 points
	score += timeBonus
	
	// 7. Volume Bonus (larger orders get slight preference)
	volumeBonus := me.calculateVolumeBonus(order1.Amount, order2.Amount)
	score += volumeBonus
	
	// 8. User Stats Bonus
	userStats := me.keeper.GetUserStats(ctx, order2.Creator)
	if userStats != nil {
		// Completion rate bonus
		if userStats.TotalTrades > 0 {
			completionRate := float64(userStats.SuccessfulTrades) / float64(userStats.TotalTrades)
			score += completionRate * 5
		}
		
		// Low dispute bonus
		if userStats.DisputesLost == 0 && userStats.TotalTrades > 10 {
			score += 5
		}
	}
	
	return score
}

// calculateDistance calculates distance between two Indian postal codes
func (me *MatchingEngine) calculateDistance(pincode1, pincode2 string) float64 {
	// Indian postal codes structure:
	// First digit: region/zone
	// First 2 digits: sub-region  
	// First 3 digits: sorting district
	// Last 3 digits: specific post office
	
	if pincode1 == pincode2 {
		return 0
	}
	
	// Same sorting district (first 3 digits)
	if pincode1[:3] == pincode2[:3] {
		return 5.0 // ~5km average
	}
	
	// Same sub-region (first 2 digits)
	if pincode1[:2] == pincode2[:2] {
		return 25.0 // ~25km average
	}
	
	// Same region (first digit)
	if pincode1[0] == pincode2[0] {
		return 100.0 // ~100km average
	}
	
	// Different regions
	return 500.0 // ~500km average
}

// calculateDistanceScore converts distance to a normalized score
func (me *MatchingEngine) calculateDistanceScore(distance, maxDistance float64) float64 {
	if distance == 0 {
		return 1.0 // Perfect score for same location
	}
	
	if distance > maxDistance {
		return 0
	}
	
	// Linear decay
	return 1.0 - (distance / maxDistance)
}

// calculatePriceMatch calculates how well prices match
func (me *MatchingEngine) calculatePriceMatch(order1, order2 *types.P2POrder) float64 {
	price1 := me.calculatePrice(order1)
	price2 := me.calculatePrice(order2)
	
	if price1.IsZero() || price2.IsZero() {
		return 0
	}
	
	// Calculate deviation
	deviation := price1.Sub(price2).Abs().Quo(price1)
	
	if deviation.GT(me.maxPriceDeviation) {
		return 0
	}
	
	// Convert to score (closer prices = higher score)
	return (me.maxPriceDeviation.Sub(deviation)).Quo(me.maxPriceDeviation).MustFloat64()
}

// calculatePrice calculates price per NAMO
func (me *MatchingEngine) calculatePrice(order *types.P2POrder) sdk.Dec {
	if order.Amount.IsZero() {
		return sdk.ZeroDec()
	}
	
	fiatDec := sdk.NewDecFromInt(order.FiatAmount.Amount)
	namoDec := sdk.NewDecFromInt(order.Amount.Amount)
	
	return fiatDec.Quo(namoDec)
}

// calculatePaymentMethodScore scores payment method compatibility
func (me *MatchingEngine) calculatePaymentMethodScore(methods1, methods2 []types.PaymentMethod) float64 {
	if len(methods1) == 0 || len(methods2) == 0 {
		return 0
	}
	
	commonMethods := 0
	totalMethods := len(methods1) + len(methods2)
	
	for _, m1 := range methods1 {
		for _, m2 := range methods2 {
			if me.arePaymentMethodsCompatible(m1, m2) {
				commonMethods++
			}
		}
	}
	
	if commonMethods == 0 {
		return 0
	}
	
	// Score based on overlap percentage
	return float64(commonMethods*2) / float64(totalMethods)
}

// arePaymentMethodsCompatible checks if two payment methods are compatible
func (me *MatchingEngine) arePaymentMethodsCompatible(m1, m2 types.PaymentMethod) bool {
	// Same method type
	if m1.MethodType != m2.MethodType {
		return false
	}
	
	// For UPI, any provider works
	if m1.MethodType == "UPI" {
		return true
	}
	
	// For bank transfers, same or compatible banks
	if m1.MethodType == "IMPS" || m1.MethodType == "NEFT" {
		return true // All banks support these
	}
	
	// For specific providers
	return m1.Provider == m2.Provider
}

// hasCommonLanguage checks for language overlap
func (me *MatchingEngine) hasCommonLanguage(langs1, langs2 []string) bool {
	// Hindi is default fallback
	if len(langs1) == 0 || len(langs2) == 0 {
		return true
	}
	
	langMap := make(map[string]bool)
	for _, lang := range langs1 {
		langMap[strings.ToLower(lang)] = true
	}
	
	for _, lang := range langs2 {
		if langMap[strings.ToLower(lang)] {
			return true
		}
	}
	
	return false
}

// calculateVolumeBonus gives bonus for order size
func (me *MatchingEngine) calculateVolumeBonus(amount1, amount2 sdk.Coin) float64 {
	// Average volume
	avgVolume := amount1.Amount.Add(amount2.Amount).QuoRaw(2)
	
	// Logarithmic scale bonus (larger orders get bonus but with diminishing returns)
	volumeFloat := float64(avgVolume.Int64()) / 1000000 // Convert to millions
	if volumeFloat > 1 {
		return math.Log10(volumeFloat) * 2 // Max ~10 points for very large orders
	}
	
	return 0
}

// areOrdersCompatible performs basic compatibility checks
func (me *MatchingEngine) areOrdersCompatible(order1, order2 *types.P2POrder) bool {
	// Opposite order types
	if order1.OrderType == order2.OrderType {
		return false
	}
	
	// Both active
	if order1.Status != types.P2POrderStatus_P2P_STATUS_ACTIVE ||
		order2.Status != types.P2POrderStatus_P2P_STATUS_ACTIVE {
		return false
	}
	
	// Same fiat currency
	if order1.FiatCurrency != order2.FiatCurrency {
		return false
	}
	
	// Not expired
	now := time.Now()
	if now.After(order1.ExpiresAt) || now.After(order2.ExpiresAt) {
		return false
	}
	
	// Amount overlap
	if !me.hasAmountOverlap(order1, order2) {
		return false
	}
	
	// Trust score requirements
	user1Trust := me.keeper.GetUserTrustScore(order1.Creator)
	user2Trust := me.keeper.GetUserTrustScore(order2.Creator)
	
	if user1Trust < order2.MinTrustScore || user2Trust < order1.MinTrustScore {
		return false
	}
	
	// KYC requirements
	if order1.RequireKyc && !me.keeper.IsKYCVerified(order2.Creator) {
		return false
	}
	if order2.RequireKyc && !me.keeper.IsKYCVerified(order1.Creator) {
		return false
	}
	
	return true
}

// hasAmountOverlap checks if order amounts overlap
func (me *MatchingEngine) hasAmountOverlap(order1, order2 *types.P2POrder) bool {
	// For full amount orders
	if order1.MinAmount.IsZero() && order2.MinAmount.IsZero() {
		return order1.Amount.Equal(order2.Amount)
	}
	
	// Check if ranges overlap
	min1 := order1.MinAmount
	if min1.IsZero() {
		min1 = order1.Amount
	}
	max1 := order1.MaxAmount
	if max1.IsZero() {
		max1 = order1.Amount
	}
	
	min2 := order2.MinAmount
	if min2.IsZero() {
		min2 = order2.Amount
	}
	max2 := order2.MaxAmount
	if max2.IsZero() {
		max2 = order2.Amount
	}
	
	// Check overlap
	return min1.IsLTE(max2) && min2.IsLTE(max1)
}

// GetMarketDepth returns current market depth for a currency pair
func (me *MatchingEngine) GetMarketDepth(currency string, levels int) *MarketDepth {
	me.orderBook.RLock()
	defer me.orderBook.RUnlock()
	
	depth := &MarketDepth{
		Currency:  currency,
		Timestamp: time.Now(),
		BuyLevels: make([]PriceDepth, 0),
		SellLevels: make([]PriceDepth, 0),
	}
	
	// Get buy side depth
	if buyIndex := me.orderBook.buyPriceIndex[currency]; buyIndex != nil {
		depth.BuyLevels = me.getPriceLevels(buyIndex, levels, true)
	}
	
	// Get sell side depth  
	if sellIndex := me.orderBook.sellPriceIndex[currency]; sellIndex != nil {
		depth.SellLevels = me.getPriceLevels(sellIndex, levels, false)
	}
	
	return depth
}

// getPriceLevels extracts price levels from index
func (me *MatchingEngine) getPriceLevels(priceIndex *PriceLevel, maxLevels int, descending bool) []PriceDepth {
	priceIndex.RLock()
	defer priceIndex.RUnlock()
	
	// Get all price points
	var prices []string
	for price := range priceIndex.levels {
		prices = append(prices, price)
	}
	
	// Sort prices
	sort.Slice(prices, func(i, j int) bool {
		p1, _ := sdk.NewDecFromStr(prices[i])
		p2, _ := sdk.NewDecFromStr(prices[j])
		if descending {
			return p1.GT(p2)
		}
		return p1.LT(p2)
	})
	
	// Build depth levels
	var levels []PriceDepth
	for i, price := range prices {
		if i >= maxLevels {
			break
		}
		
		orders := priceIndex.levels[price]
		totalVolume := sdk.ZeroInt()
		for _, order := range orders {
			totalVolume = totalVolume.Add(order.Amount.Amount)
		}
		
		priceDec, _ := sdk.NewDecFromStr(price)
		levels = append(levels, PriceDepth{
			Price:      priceDec,
			Volume:     sdk.NewCoin(types.DefaultDenom, totalVolume),
			OrderCount: len(orders),
		})
	}
	
	return levels
}

// MatchResult represents a potential match with scoring details
type MatchResult struct {
	Order      *types.P2POrder
	Score      float64
	Distance   float64
	PriceMatch float64
	TrustScore int32
}

// MarketDepth represents order book depth
type MarketDepth struct {
	Currency   string
	Timestamp  time.Time
	BuyLevels  []PriceDepth
	SellLevels []PriceDepth
}

// PriceDepth represents orders at a price level
type PriceDepth struct {
	Price      sdk.Dec
	Volume     sdk.Coin
	OrderCount int
}

// GetOrderBookStats returns statistics about the order book
func (me *MatchingEngine) GetOrderBookStats() *OrderBookStats {
	me.orderBook.RLock()
	defer me.orderBook.RUnlock()
	
	stats := &OrderBookStats{
		TotalBuyOrders:  len(me.orderBook.buyOrders),
		TotalSellOrders: len(me.orderBook.sellOrders),
		Currencies:      make(map[string]int),
		Districts:       make(map[string]int),
		AvgMatchTime:    me.avgMatchTime,
		TotalMatches:    me.matchCount,
		LastMatchTime:   me.lastMatchTimestamp,
	}
	
	// Count by currency
	for currency, orderIDs := range me.orderBook.currencyIndex {
		stats.Currencies[currency] = len(orderIDs)
	}
	
	// Count by district
	for district, orderIDs := range me.orderBook.districtIndex {
		stats.Districts[district] = len(orderIDs)
	}
	
	return stats
}

// OrderBookStats represents order book statistics
type OrderBookStats struct {
	TotalBuyOrders  int
	TotalSellOrders int
	Currencies      map[string]int
	Districts       map[string]int
	AvgMatchTime    time.Duration
	TotalMatches    int64
	LastMatchTime   time.Time
}