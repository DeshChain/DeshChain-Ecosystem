package types

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
)

// ExplorerIndexer provides indexing functionality for the explorer
type ExplorerIndexer struct {
	blockIndex        map[int64]*BlockInfo
	transactionIndex  map[string]*TransactionInfo
	validatorIndex    map[string]*ValidatorInfo
	addressIndex      map[string]*AddressInfo
	heightIndex       map[int64][]string // height -> transaction hashes
	addressTxIndex    map[string][]string // address -> transaction hashes
	searchIndex       map[string][]SearchResult
	culturalQuoteIndex map[uint64]*CulturalQuoteDisplay
	patriotismIndex   map[string]*PatriotismLeaderboard
	holderIndex       map[string]*HolderRanking
	donationIndex     map[string]*DonationTracker
	burnIndex         map[uint64]*TokenBurnInfo
	activityIndex     []RecentActivity
	networkStats      *NetworkStats
	taxStats          *TaxStatistics
	lastUpdateTime    time.Time
}

// NewExplorerIndexer creates a new explorer indexer
func NewExplorerIndexer() *ExplorerIndexer {
	return &ExplorerIndexer{
		blockIndex:         make(map[int64]*BlockInfo),
		transactionIndex:   make(map[string]*TransactionInfo),
		validatorIndex:     make(map[string]*ValidatorInfo),
		addressIndex:       make(map[string]*AddressInfo),
		heightIndex:        make(map[int64][]string),
		addressTxIndex:     make(map[string][]string),
		searchIndex:        make(map[string][]SearchResult),
		culturalQuoteIndex: make(map[uint64]*CulturalQuoteDisplay),
		patriotismIndex:    make(map[string]*PatriotismLeaderboard),
		holderIndex:        make(map[string]*HolderRanking),
		donationIndex:      make(map[string]*DonationTracker),
		burnIndex:          make(map[uint64]*TokenBurnInfo),
		activityIndex:      make([]RecentActivity, 0),
		networkStats:       &NetworkStats{},
		taxStats:           &TaxStatistics{},
		lastUpdateTime:     time.Now(),
	}
}

// IndexBlock indexes a block and its transactions
func (idx *ExplorerIndexer) IndexBlock(ctx context.Context, block *BlockInfo) error {
	if block == nil {
		return ErrInvalidData
	}

	// Index the block
	idx.blockIndex[block.Height] = block

	// Record activity
	activity := RecentActivity{
		ActivityType:    ActivityTypeBlock,
		Description:     fmt.Sprintf("Block %d proposed by %s", block.Height, block.ProposerAddress),
		Timestamp:       block.Timestamp,
		TransactionHash: "",
		BlockHeight:     block.Height,
		Amount:          block.Reward.Amount,
		FromAddress:     "",
		ToAddress:       block.ProposerAddress,
		AdditionalInfo: map[string]string{
			"transaction_count": fmt.Sprintf("%d", block.TransactionCount),
			"gas_used":          fmt.Sprintf("%d", block.GasUsed),
			"gas_limit":         fmt.Sprintf("%d", block.GasLimit),
			"size":              fmt.Sprintf("%d", block.Size),
			"cultural_quote_id": fmt.Sprintf("%d", block.CulturalQuoteId),
			"tax_collected":     block.TaxCollected.String(),
			"donations_count":   fmt.Sprintf("%d", block.DonationsCount),
			"burned_amount":     block.BurnedAmount.String(),
		},
	}
	idx.AddActivity(activity)

	// Update network stats
	idx.UpdateNetworkStats(block)

	// Update search index
	idx.UpdateSearchIndex(block)

	return nil
}

// IndexTransaction indexes a transaction
func (idx *ExplorerIndexer) IndexTransaction(ctx context.Context, tx *TransactionInfo) error {
	if tx == nil {
		return ErrInvalidData
	}

	// Index the transaction
	idx.transactionIndex[tx.Hash] = tx

	// Update height index
	if idx.heightIndex[tx.Height] == nil {
		idx.heightIndex[tx.Height] = make([]string, 0)
	}
	idx.heightIndex[tx.Height] = append(idx.heightIndex[tx.Height], tx.Hash)

	// Update address indexes
	if tx.FromAddress != "" {
		if idx.addressTxIndex[tx.FromAddress] == nil {
			idx.addressTxIndex[tx.FromAddress] = make([]string, 0)
		}
		idx.addressTxIndex[tx.FromAddress] = append(idx.addressTxIndex[tx.FromAddress], tx.Hash)
		idx.UpdateAddressInfo(tx.FromAddress, tx)
	}

	if tx.ToAddress != "" {
		if idx.addressTxIndex[tx.ToAddress] == nil {
			idx.addressTxIndex[tx.ToAddress] = make([]string, 0)
		}
		idx.addressTxIndex[tx.ToAddress] = append(idx.addressTxIndex[tx.ToAddress], tx.Hash)
		idx.UpdateAddressInfo(tx.ToAddress, tx)
	}

	// Record activity
	activity := RecentActivity{
		ActivityType:    ActivityTypeTransaction,
		Description:     fmt.Sprintf("Transaction %s", tx.MessageType),
		Timestamp:       tx.Timestamp,
		TransactionHash: tx.Hash,
		BlockHeight:     tx.Height,
		Amount:          tx.Amount.Amount,
		FromAddress:     tx.FromAddress,
		ToAddress:       tx.ToAddress,
		AdditionalInfo: map[string]string{
			"message_type":    tx.MessageType,
			"status":          tx.Status,
			"gas_used":        fmt.Sprintf("%d", tx.GasUsed),
			"gas_limit":       fmt.Sprintf("%d", tx.GasLimit),
			"fee":             tx.Fee.String(),
			"tax_amount":      tx.TaxAmount.String(),
			"is_donation":     fmt.Sprintf("%t", tx.IsDonation),
			"patriotism_score": fmt.Sprintf("%d", tx.PatriotismScore),
			"burn_amount":     tx.BurnAmount.String(),
		},
	}
	idx.AddActivity(activity)

	// Update cultural quote index if present
	if tx.CulturalQuoteId > 0 {
		quote := &CulturalQuoteDisplay{
			QuoteId:        tx.CulturalQuoteId,
			QuoteText:      tx.CulturalQuoteText,
			QuoteAuthor:    tx.CulturalQuoteAuthor,
			TransactionHash: tx.Hash,
			BlockHeight:    tx.Height,
			Timestamp:      tx.Timestamp,
		}
		idx.culturalQuoteIndex[tx.CulturalQuoteId] = quote
	}

	// Update patriotism scores
	if tx.PatriotismScore > 0 {
		idx.UpdatePatriotismScore(tx.FromAddress, tx.PatriotismScore)
	}

	// Update donation tracking
	if tx.IsDonation {
		idx.UpdateDonationTracking(tx)
	}

	// Update tax statistics
	if !tx.TaxAmount.IsZero() {
		idx.UpdateTaxStatistics(tx)
	}

	// Update burn tracking
	if !tx.BurnAmount.IsZero() {
		burnInfo := &TokenBurnInfo{
			Id:              uint64(len(idx.burnIndex) + 1),
			TransactionHash: tx.Hash,
			BurnerAddress:   tx.FromAddress,
			Amount:          tx.BurnAmount,
			Timestamp:       tx.Timestamp,
			Height:          tx.Height,
			Reason:          "Transaction burn",
			BurnType:        BurnTypeAutomatic,
			CulturalQuoteId: tx.CulturalQuoteId,
		}
		idx.burnIndex[burnInfo.Id] = burnInfo
	}

	// Update search index
	idx.UpdateSearchIndexForTransaction(tx)

	return nil
}

// IndexValidator indexes validator information
func (idx *ExplorerIndexer) IndexValidator(ctx context.Context, validator *ValidatorInfo) error {
	if validator == nil {
		return ErrInvalidData
	}

	idx.validatorIndex[validator.OperatorAddress] = validator

	// Update search index
	idx.UpdateSearchIndexForValidator(validator)

	return nil
}

// GetBlock retrieves a block by height
func (idx *ExplorerIndexer) GetBlock(height int64) (*BlockInfo, error) {
	block, exists := idx.blockIndex[height]
	if !exists {
		return nil, ErrBlockNotFound
	}
	return block, nil
}

// GetTransaction retrieves a transaction by hash
func (idx *ExplorerIndexer) GetTransaction(hash string) (*TransactionInfo, error) {
	tx, exists := idx.transactionIndex[hash]
	if !exists {
		return nil, ErrTransactionNotFound
	}
	return tx, nil
}

// GetValidator retrieves a validator by operator address
func (idx *ExplorerIndexer) GetValidator(operatorAddress string) (*ValidatorInfo, error) {
	validator, exists := idx.validatorIndex[operatorAddress]
	if !exists {
		return nil, ErrValidatorNotFound
	}
	return validator, nil
}

// GetAddress retrieves address information
func (idx *ExplorerIndexer) GetAddress(address string) (*AddressInfo, error) {
	info, exists := idx.addressIndex[address]
	if !exists {
		return nil, ErrAddressNotFound
	}
	return info, nil
}

// GetBlocksPaginated retrieves blocks with pagination
func (idx *ExplorerIndexer) GetBlocksPaginated(req *query.PageRequest) ([]*BlockInfo, *query.PageResponse, error) {
	// Get all block heights and sort them
	heights := make([]int64, 0, len(idx.blockIndex))
	for height := range idx.blockIndex {
		heights = append(heights, height)
	}
	sort.Slice(heights, func(i, j int) bool {
		return heights[i] > heights[j] // Latest first
	})

	// Apply pagination
	limit := DefaultMaxSearchResults
	if req != nil && req.Limit > 0 {
		limit = int(req.Limit)
		if limit > MaxBlocksPerPage {
			limit = MaxBlocksPerPage
		}
	}

	offset := 0
	if req != nil && req.Offset > 0 {
		offset = int(req.Offset)
	}

	// Get blocks for the requested page
	blocks := make([]*BlockInfo, 0, limit)
	for i := offset; i < len(heights) && len(blocks) < limit; i++ {
		if block, exists := idx.blockIndex[heights[i]]; exists {
			blocks = append(blocks, block)
		}
	}

	// Create page response
	pageResponse := &query.PageResponse{
		Total: uint64(len(heights)),
	}
	if offset+limit < len(heights) {
		pageResponse.NextKey = []byte(fmt.Sprintf("%d", offset+limit))
	}

	return blocks, pageResponse, nil
}

// GetTransactionsPaginated retrieves transactions with pagination
func (idx *ExplorerIndexer) GetTransactionsPaginated(req *query.PageRequest) ([]*TransactionInfo, *query.PageResponse, error) {
	// Get all transaction hashes and sort by timestamp
	type txEntry struct {
		hash      string
		timestamp int64
	}
	
	entries := make([]txEntry, 0, len(idx.transactionIndex))
	for hash, tx := range idx.transactionIndex {
		entries = append(entries, txEntry{hash: hash, timestamp: tx.Timestamp})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].timestamp > entries[j].timestamp // Latest first
	})

	// Apply pagination
	limit := DefaultMaxSearchResults
	if req != nil && req.Limit > 0 {
		limit = int(req.Limit)
		if limit > MaxTransactionsPerPage {
			limit = MaxTransactionsPerPage
		}
	}

	offset := 0
	if req != nil && req.Offset > 0 {
		offset = int(req.Offset)
	}

	// Get transactions for the requested page
	transactions := make([]*TransactionInfo, 0, limit)
	for i := offset; i < len(entries) && len(transactions) < limit; i++ {
		if tx, exists := idx.transactionIndex[entries[i].hash]; exists {
			transactions = append(transactions, tx)
		}
	}

	// Create page response
	pageResponse := &query.PageResponse{
		Total: uint64(len(entries)),
	}
	if offset+limit < len(entries) {
		pageResponse.NextKey = []byte(fmt.Sprintf("%d", offset+limit))
	}

	return transactions, pageResponse, nil
}

// GetValidatorsPaginated retrieves validators with pagination
func (idx *ExplorerIndexer) GetValidatorsPaginated(req *query.PageRequest) ([]*ValidatorInfo, *query.PageResponse, error) {
	// Get all validators and sort by voting power
	validators := make([]*ValidatorInfo, 0, len(idx.validatorIndex))
	for _, validator := range idx.validatorIndex {
		validators = append(validators, validator)
	}
	sort.Slice(validators, func(i, j int) bool {
		return validators[i].Rank < validators[j].Rank // Best rank first
	})

	// Apply pagination
	limit := DefaultMaxSearchResults
	if req != nil && req.Limit > 0 {
		limit = int(req.Limit)
		if limit > MaxValidatorsPerPage {
			limit = MaxValidatorsPerPage
		}
	}

	offset := 0
	if req != nil && req.Offset > 0 {
		offset = int(req.Offset)
	}

	// Get validators for the requested page
	result := make([]*ValidatorInfo, 0, limit)
	for i := offset; i < len(validators) && len(result) < limit; i++ {
		result = append(result, validators[i])
	}

	// Create page response
	pageResponse := &query.PageResponse{
		Total: uint64(len(validators)),
	}
	if offset+limit < len(validators) {
		pageResponse.NextKey = []byte(fmt.Sprintf("%d", offset+limit))
	}

	return result, pageResponse, nil
}

// GetAddressTransactions retrieves transactions for a specific address
func (idx *ExplorerIndexer) GetAddressTransactions(address string, req *query.PageRequest) ([]*TransactionInfo, *query.PageResponse, error) {
	hashes, exists := idx.addressTxIndex[address]
	if !exists {
		return []*TransactionInfo{}, &query.PageResponse{Total: 0}, nil
	}

	// Sort transactions by timestamp (latest first)
	type txEntry struct {
		hash      string
		timestamp int64
	}
	
	entries := make([]txEntry, 0, len(hashes))
	for _, hash := range hashes {
		if tx, exists := idx.transactionIndex[hash]; exists {
			entries = append(entries, txEntry{hash: hash, timestamp: tx.Timestamp})
		}
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].timestamp > entries[j].timestamp
	})

	// Apply pagination
	limit := DefaultMaxSearchResults
	if req != nil && req.Limit > 0 {
		limit = int(req.Limit)
		if limit > MaxTransactionsPerPage {
			limit = MaxTransactionsPerPage
		}
	}

	offset := 0
	if req != nil && req.Offset > 0 {
		offset = int(req.Offset)
	}

	// Get transactions for the requested page
	transactions := make([]*TransactionInfo, 0, limit)
	for i := offset; i < len(entries) && len(transactions) < limit; i++ {
		if tx, exists := idx.transactionIndex[entries[i].hash]; exists {
			transactions = append(transactions, tx)
		}
	}

	// Create page response
	pageResponse := &query.PageResponse{
		Total: uint64(len(entries)),
	}
	if offset+limit < len(entries) {
		pageResponse.NextKey = []byte(fmt.Sprintf("%d", offset+limit))
	}

	return transactions, pageResponse, nil
}

// Search performs a search across indexed data
func (idx *ExplorerIndexer) Search(query string, searchType string, limit uint32) ([]SearchResult, error) {
	if query == "" {
		return []SearchResult{}, ErrInvalidSearchQuery
	}

	query = strings.ToLower(strings.TrimSpace(query))
	results := make([]SearchResult, 0)

	// Search based on type
	switch searchType {
	case SearchTypeAll:
		results = append(results, idx.searchBlocks(query)...)
		results = append(results, idx.searchTransactions(query)...)
		results = append(results, idx.searchValidators(query)...)
		results = append(results, idx.searchAddresses(query)...)
	case SearchTypeBlocks:
		results = idx.searchBlocks(query)
	case SearchTypeTransactions:
		results = idx.searchTransactions(query)
	case SearchTypeValidators:
		results = idx.searchValidators(query)
	case SearchTypeAddresses:
		results = idx.searchAddresses(query)
	default:
		return []SearchResult{}, ErrInvalidSearchType
	}

	// Sort by relevance score
	sort.Slice(results, func(i, j int) bool {
		return results[i].RelevanceScore > results[j].RelevanceScore
	})

	// Apply limit
	if limit > 0 && len(results) > int(limit) {
		results = results[:limit]
	}

	return results, nil
}

// Helper methods for search functionality

func (idx *ExplorerIndexer) searchBlocks(query string) []SearchResult {
	results := make([]SearchResult, 0)
	
	for height, block := range idx.blockIndex {
		score := 0
		
		// Exact height match
		if fmt.Sprintf("%d", height) == query {
			score = 100
		}
		
		// Hash match
		if strings.Contains(strings.ToLower(block.Hash), query) {
			score = 90
		}
		
		// Proposer address match
		if strings.Contains(strings.ToLower(block.ProposerAddress), query) {
			score = 80
		}
		
		if score > 0 {
			results = append(results, SearchResult{
				ResultType:     ResultTypeBlock,
				Title:          fmt.Sprintf("Block %d", height),
				Description:    fmt.Sprintf("Block %d proposed by %s", height, block.ProposerAddress),
				Url:            fmt.Sprintf("/block/%d", height),
				RelevanceScore: int32(score),
				AdditionalInfo: map[string]string{
					"height":       fmt.Sprintf("%d", height),
					"hash":         block.Hash,
					"proposer":     block.ProposerAddress,
					"tx_count":     fmt.Sprintf("%d", block.TransactionCount),
					"timestamp":    fmt.Sprintf("%d", block.Timestamp),
				},
			})
		}
	}
	
	return results
}

func (idx *ExplorerIndexer) searchTransactions(query string) []SearchResult {
	results := make([]SearchResult, 0)
	
	for hash, tx := range idx.transactionIndex {
		score := 0
		
		// Exact hash match
		if strings.Contains(strings.ToLower(hash), query) {
			score = 100
		}
		
		// Address matches
		if strings.Contains(strings.ToLower(tx.FromAddress), query) ||
		   strings.Contains(strings.ToLower(tx.ToAddress), query) {
			score = 90
		}
		
		// Message type match
		if strings.Contains(strings.ToLower(tx.MessageType), query) {
			score = 70
		}
		
		// Memo match
		if strings.Contains(strings.ToLower(tx.Memo), query) {
			score = 60
		}
		
		if score > 0 {
			results = append(results, SearchResult{
				ResultType:     ResultTypeTransaction,
				Title:          fmt.Sprintf("Transaction %s", hash[:8]),
				Description:    fmt.Sprintf("%s transaction from %s to %s", tx.MessageType, tx.FromAddress, tx.ToAddress),
				Url:            fmt.Sprintf("/transaction/%s", hash),
				RelevanceScore: int32(score),
				AdditionalInfo: map[string]string{
					"hash":         hash,
					"from":         tx.FromAddress,
					"to":           tx.ToAddress,
					"amount":       tx.Amount.String(),
					"message_type": tx.MessageType,
					"status":       tx.Status,
					"timestamp":    fmt.Sprintf("%d", tx.Timestamp),
				},
			})
		}
	}
	
	return results
}

func (idx *ExplorerIndexer) searchValidators(query string) []SearchResult {
	results := make([]SearchResult, 0)
	
	for address, validator := range idx.validatorIndex {
		score := 0
		
		// Exact address match
		if strings.Contains(strings.ToLower(address), query) {
			score = 100
		}
		
		// Moniker match
		if strings.Contains(strings.ToLower(validator.Moniker), query) {
			score = 90
		}
		
		// Identity match
		if strings.Contains(strings.ToLower(validator.Identity), query) {
			score = 80
		}
		
		// Website match
		if strings.Contains(strings.ToLower(validator.Website), query) {
			score = 70
		}
		
		if score > 0 {
			results = append(results, SearchResult{
				ResultType:     ResultTypeValidator,
				Title:          fmt.Sprintf("Validator %s", validator.Moniker),
				Description:    fmt.Sprintf("Validator %s with %s voting power", validator.Moniker, validator.VotingPower),
				Url:            fmt.Sprintf("/validator/%s", address),
				RelevanceScore: int32(score),
				AdditionalInfo: map[string]string{
					"address":      address,
					"moniker":      validator.Moniker,
					"voting_power": validator.VotingPower,
					"commission":   validator.CommissionRate,
					"status":       validator.Status,
					"jailed":       fmt.Sprintf("%t", validator.Jailed),
				},
			})
		}
	}
	
	return results
}

func (idx *ExplorerIndexer) searchAddresses(query string) []SearchResult {
	results := make([]SearchResult, 0)
	
	for address, info := range idx.addressIndex {
		score := 0
		
		// Exact address match
		if strings.Contains(strings.ToLower(address), query) {
			score = 100
		}
		
		// NGO name match
		if strings.Contains(strings.ToLower(info.NgoName), query) {
			score = 90
		}
		
		// Account type match
		if strings.Contains(strings.ToLower(info.AccountType), query) {
			score = 80
		}
		
		if score > 0 {
			results = append(results, SearchResult{
				ResultType:     ResultTypeAddress,
				Title:          fmt.Sprintf("Address %s", address[:8]),
				Description:    fmt.Sprintf("Address with balance %s", info.Balance.String()),
				Url:            fmt.Sprintf("/address/%s", address),
				RelevanceScore: int32(score),
				AdditionalInfo: map[string]string{
					"address":        address,
					"balance":        info.Balance.String(),
					"account_type":   info.AccountType,
					"tx_count":       fmt.Sprintf("%d", info.TransactionCount),
					"patriotism_score": fmt.Sprintf("%d", info.PatriotismScore),
					"holder_rank":    fmt.Sprintf("%d", info.HolderRank),
				},
			})
		}
	}
	
	return results
}

// Helper methods for updating indexes

func (idx *ExplorerIndexer) UpdateAddressInfo(address string, tx *TransactionInfo) {
	info, exists := idx.addressIndex[address]
	if !exists {
		info = &AddressInfo{
			Address:            address,
			Balance:            sdk.NewCoin("namo", sdk.ZeroInt()),
			TransactionCount:   0,
			FirstSeen:          tx.Timestamp,
			LastSeen:           tx.Timestamp,
			TotalSent:          sdk.NewCoin("namo", sdk.ZeroInt()),
			TotalReceived:      sdk.NewCoin("namo", sdk.ZeroInt()),
			TotalFeesPaid:      sdk.NewCoin("namo", sdk.ZeroInt()),
			TotalDonations:     sdk.NewCoin("namo", sdk.ZeroInt()),
			TotalTaxPaid:       sdk.NewCoin("namo", sdk.ZeroInt()),
			PatriotismScore:    DefaultPatriotismScore,
			PatriotismRank:     0,
			HolderRank:         0,
			IsValidator:        false,
			IsNgo:              false,
			NgoName:            "",
			AccountType:        AccountTypeRegular,
			CulturalQuotesReceived: 0,
			FavoriteQuoteCategories: make([]string, 0),
		}
		idx.addressIndex[address] = info
	}

	// Update transaction count
	info.TransactionCount++
	info.LastSeen = tx.Timestamp

	// Update amounts based on transaction direction
	if tx.FromAddress == address {
		info.TotalSent = info.TotalSent.Add(tx.Amount)
		info.TotalFeesPaid = info.TotalFeesPaid.Add(tx.Fee)
		info.TotalTaxPaid = info.TotalTaxPaid.Add(tx.TaxAmount)
	}
	
	if tx.ToAddress == address {
		info.TotalReceived = info.TotalReceived.Add(tx.Amount)
	}

	// Update donation amount
	if tx.IsDonation && tx.FromAddress == address {
		info.TotalDonations = info.TotalDonations.Add(tx.Amount)
	}

	// Update cultural quotes received
	if tx.CulturalQuoteId > 0 {
		info.CulturalQuotesReceived++
	}

	// Update patriotism score
	if tx.PatriotismScore > 0 && tx.FromAddress == address {
		info.PatriotismScore += tx.PatriotismScore
	}
}

func (idx *ExplorerIndexer) UpdateNetworkStats(block *BlockInfo) {
	idx.networkStats.CurrentHeight = block.Height
	idx.networkStats.TotalTransactions += block.TransactionCount
	idx.networkStats.TotalDonations = idx.networkStats.TotalDonations.Add(block.TotalDonations)
	idx.networkStats.TotalTaxCollected = idx.networkStats.TotalTaxCollected.Add(block.TaxCollected)
	idx.networkStats.BurnedSupply = idx.networkStats.BurnedSupply.Add(block.BurnedAmount)
	idx.networkStats.QuotesDisplayedToday++
}

func (idx *ExplorerIndexer) UpdatePatriotismScore(address string, score int32) {
	entry, exists := idx.patriotismIndex[address]
	if !exists {
		entry = &PatriotismLeaderboard{
			Address:          address,
			PatriotismScore:  score,
			Rank:             0,
			TotalDonations:   sdk.NewCoin("namo", sdk.ZeroInt()),
			DonationCount:    0,
			CulturalEngagement: 0,
			ConsistencyScore: 0,
			LastDonationTime: time.Now().Unix(),
			BadgeLevel:       BadgeLevelBronze,
			Achievements:     make([]string, 0),
		}
		idx.patriotismIndex[address] = entry
	}
	
	entry.PatriotismScore += score
	entry.CulturalEngagement++
	
	// Update badge level based on score
	if entry.PatriotismScore >= LegendThreshold {
		entry.BadgeLevel = BadgeLevelLegend
	} else if entry.PatriotismScore >= DiamondThreshold {
		entry.BadgeLevel = BadgeLevelDiamond
	} else if entry.PatriotismScore >= PlatinumThreshold {
		entry.BadgeLevel = BadgeLevelPlatinum
	} else if entry.PatriotismScore >= GoldThreshold {
		entry.BadgeLevel = BadgeLevelGold
	} else if entry.PatriotismScore >= SilverThreshold {
		entry.BadgeLevel = BadgeLevelSilver
	}
}

func (idx *ExplorerIndexer) UpdateDonationTracking(tx *TransactionInfo) {
	// Update donation statistics
	tracker, exists := idx.donationIndex["global"]
	if !exists {
		tracker = &DonationTracker{
			TotalDonations:      sdk.NewCoin("namo", sdk.ZeroInt()),
			DonationCount:       0,
			ActiveNgos:          0,
			TotalBeneficiaries:  0,
			FundsDistributed:    sdk.NewCoin("namo", sdk.ZeroInt()),
			AverageDonationSize: sdk.NewCoin("namo", sdk.ZeroInt()),
			TransparencyScore:   "10.0",
			TopDonationCategories: make([]CategoryStats, 0),
			RecentDonations:     make([]RecentDonation, 0),
			DonationGrowthRate:  "0.0",
		}
		idx.donationIndex["global"] = tracker
	}
	
	tracker.TotalDonations = tracker.TotalDonations.Add(tx.Amount)
	tracker.DonationCount++
	
	// Add to recent donations
	recentDonation := RecentDonation{
		TransactionHash: tx.Hash,
		Amount:          tx.Amount,
		NgoName:         tx.DonationNgoName,
		Purpose:         tx.DonationPurpose,
		Timestamp:       tx.Timestamp,
		IsAnonymous:     false, // Based on transaction visibility
		CulturalQuoteText: tx.CulturalQuoteText,
	}
	
	tracker.RecentDonations = append(tracker.RecentDonations, recentDonation)
	
	// Keep only recent donations (last 100)
	if len(tracker.RecentDonations) > 100 {
		tracker.RecentDonations = tracker.RecentDonations[len(tracker.RecentDonations)-100:]
	}
}

func (idx *ExplorerIndexer) UpdateTaxStatistics(tx *TransactionInfo) {
	if idx.taxStats.TaxTransactionsCount == 0 {
		idx.taxStats.TotalTaxCollected = "0"
		idx.taxStats.CurrentTaxRate = "2.5"
		idx.taxStats.TaxTransactionsCount = 0
		idx.taxStats.AverageTaxPerTransaction = "0"
		idx.taxStats.TaxDistribution = make(map[string]string)
		idx.taxStats.TaxBurnPercentage = "50.0"
		idx.taxStats.TaxGrowthRate = "0.0"
		idx.taxStats.VolumeBasedReductions = "0"
		idx.taxStats.UsersWithTaxCap = 0
	}
	
	idx.taxStats.TaxTransactionsCount++
	// Additional tax statistics updates would go here
}

func (idx *ExplorerIndexer) UpdateSearchIndex(block *BlockInfo) {
	// Update search index for the block
	results := make([]SearchResult, 0)
	
	// Add block to search index
	result := SearchResult{
		ResultType:     ResultTypeBlock,
		Title:          fmt.Sprintf("Block %d", block.Height),
		Description:    fmt.Sprintf("Block %d proposed by %s", block.Height, block.ProposerAddress),
		Url:            fmt.Sprintf("/block/%d", block.Height),
		RelevanceScore: 100,
		AdditionalInfo: map[string]string{
			"height":     fmt.Sprintf("%d", block.Height),
			"hash":       block.Hash,
			"proposer":   block.ProposerAddress,
			"tx_count":   fmt.Sprintf("%d", block.TransactionCount),
			"timestamp":  fmt.Sprintf("%d", block.Timestamp),
		},
	}
	results = append(results, result)
	
	// Index by height
	heightKey := fmt.Sprintf("%d", block.Height)
	idx.searchIndex[heightKey] = results
	
	// Index by hash
	hashKey := strings.ToLower(block.Hash)
	idx.searchIndex[hashKey] = results
}

func (idx *ExplorerIndexer) UpdateSearchIndexForTransaction(tx *TransactionInfo) {
	results := make([]SearchResult, 0)
	
	result := SearchResult{
		ResultType:     ResultTypeTransaction,
		Title:          fmt.Sprintf("Transaction %s", tx.Hash[:8]),
		Description:    fmt.Sprintf("%s transaction from %s to %s", tx.MessageType, tx.FromAddress, tx.ToAddress),
		Url:            fmt.Sprintf("/transaction/%s", tx.Hash),
		RelevanceScore: 100,
		AdditionalInfo: map[string]string{
			"hash":         tx.Hash,
			"from":         tx.FromAddress,
			"to":           tx.ToAddress,
			"amount":       tx.Amount.String(),
			"message_type": tx.MessageType,
			"status":       tx.Status,
			"timestamp":    fmt.Sprintf("%d", tx.Timestamp),
		},
	}
	results = append(results, result)
	
	// Index by hash
	hashKey := strings.ToLower(tx.Hash)
	idx.searchIndex[hashKey] = results
}

func (idx *ExplorerIndexer) UpdateSearchIndexForValidator(validator *ValidatorInfo) {
	results := make([]SearchResult, 0)
	
	result := SearchResult{
		ResultType:     ResultTypeValidator,
		Title:          fmt.Sprintf("Validator %s", validator.Moniker),
		Description:    fmt.Sprintf("Validator %s with %s voting power", validator.Moniker, validator.VotingPower),
		Url:            fmt.Sprintf("/validator/%s", validator.OperatorAddress),
		RelevanceScore: 100,
		AdditionalInfo: map[string]string{
			"address":      validator.OperatorAddress,
			"moniker":      validator.Moniker,
			"voting_power": validator.VotingPower,
			"commission":   validator.CommissionRate,
			"status":       validator.Status,
			"jailed":       fmt.Sprintf("%t", validator.Jailed),
		},
	}
	results = append(results, result)
	
	// Index by address
	addressKey := strings.ToLower(validator.OperatorAddress)
	idx.searchIndex[addressKey] = results
	
	// Index by moniker
	monikerKey := strings.ToLower(validator.Moniker)
	idx.searchIndex[monikerKey] = results
}

func (idx *ExplorerIndexer) AddActivity(activity RecentActivity) {
	idx.activityIndex = append(idx.activityIndex, activity)
	
	// Keep only recent activities (last 1000)
	if len(idx.activityIndex) > 1000 {
		idx.activityIndex = idx.activityIndex[len(idx.activityIndex)-1000:]
	}
}

func (idx *ExplorerIndexer) GetRecentActivities(limit uint32) []RecentActivity {
	if limit == 0 {
		limit = MaxRecentActivities
	}
	
	start := 0
	if len(idx.activityIndex) > int(limit) {
		start = len(idx.activityIndex) - int(limit)
	}
	
	activities := make([]RecentActivity, 0, limit)
	for i := start; i < len(idx.activityIndex); i++ {
		activities = append(activities, idx.activityIndex[i])
	}
	
	// Reverse to get latest first
	for i := 0; i < len(activities)/2; i++ {
		j := len(activities) - 1 - i
		activities[i], activities[j] = activities[j], activities[i]
	}
	
	return activities
}

func (idx *ExplorerIndexer) GetNetworkStats() *NetworkStats {
	return idx.networkStats
}

func (idx *ExplorerIndexer) GetTaxStatistics() *TaxStatistics {
	return idx.taxStats
}

func (idx *ExplorerIndexer) GetPatriotismLeaderboard(limit uint32) []PatriotismLeaderboard {
	if limit == 0 {
		limit = MaxLeaderboardEntries
	}
	
	// Get all entries and sort by score
	entries := make([]PatriotismLeaderboard, 0, len(idx.patriotismIndex))
	for _, entry := range idx.patriotismIndex {
		entries = append(entries, *entry)
	}
	
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].PatriotismScore > entries[j].PatriotismScore
	})
	
	// Update ranks
	for i := range entries {
		entries[i].Rank = uint64(i + 1)
	}
	
	// Apply limit
	if len(entries) > int(limit) {
		entries = entries[:limit]
	}
	
	return entries
}

func (idx *ExplorerIndexer) GetHolderRankings(limit uint32) []HolderRanking {
	if limit == 0 {
		limit = MaxHolderRankings
	}
	
	// Get all addresses and sort by balance
	rankings := make([]HolderRanking, 0, len(idx.addressIndex))
	for address, info := range idx.addressIndex {
		ranking := HolderRanking{
			Address:              address,
			Balance:              info.Balance,
			Rank:                 0,
			PercentageOfSupply:   "0.0",
			AccountType:          info.AccountType,
			FirstSeen:            info.FirstSeen,
			IsActive:             info.LastSeen > time.Now().Unix()-86400*30, // Active in last 30 days
			LastTransactionTime:  info.LastSeen,
		}
		rankings = append(rankings, ranking)
	}
	
	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].Balance.Amount.GT(rankings[j].Balance.Amount)
	})
	
	// Update ranks
	for i := range rankings {
		rankings[i].Rank = uint64(i + 1)
	}
	
	// Apply limit
	if len(rankings) > int(limit) {
		rankings = rankings[:limit]
	}
	
	return rankings
}

func (idx *ExplorerIndexer) GetDonationTracking() *DonationTracker {
	tracker, exists := idx.donationIndex["global"]
	if !exists {
		return &DonationTracker{
			TotalDonations:      sdk.NewCoin("namo", sdk.ZeroInt()),
			DonationCount:       0,
			ActiveNgos:          0,
			TotalBeneficiaries:  0,
			FundsDistributed:    sdk.NewCoin("namo", sdk.ZeroInt()),
			AverageDonationSize: sdk.NewCoin("namo", sdk.ZeroInt()),
			TransparencyScore:   "10.0",
			TopDonationCategories: make([]CategoryStats, 0),
			RecentDonations:     make([]RecentDonation, 0),
			DonationGrowthRate:  "0.0",
		}
	}
	return tracker
}

func (idx *ExplorerIndexer) GetTokenBurns(limit uint32) []TokenBurnInfo {
	if limit == 0 {
		limit = 100
	}
	
	// Get all burns and sort by timestamp
	burns := make([]TokenBurnInfo, 0, len(idx.burnIndex))
	for _, burn := range idx.burnIndex {
		burns = append(burns, *burn)
	}
	
	sort.Slice(burns, func(i, j int) bool {
		return burns[i].Timestamp > burns[j].Timestamp
	})
	
	// Apply limit
	if len(burns) > int(limit) {
		burns = burns[:limit]
	}
	
	return burns
}

func (idx *ExplorerIndexer) GetCulturalQuotes(limit uint32) []CulturalQuoteDisplay {
	if limit == 0 {
		limit = MaxCulturalQuotesPerPage
	}
	
	// Get all quotes and sort by timestamp
	quotes := make([]CulturalQuoteDisplay, 0, len(idx.culturalQuoteIndex))
	for _, quote := range idx.culturalQuoteIndex {
		quotes = append(quotes, *quote)
	}
	
	sort.Slice(quotes, func(i, j int) bool {
		return quotes[i].Timestamp > quotes[j].Timestamp
	})
	
	// Apply limit
	if len(quotes) > int(limit) {
		quotes = quotes[:limit]
	}
	
	return quotes
}

func (idx *ExplorerIndexer) GetLastUpdateTime() time.Time {
	return idx.lastUpdateTime
}

func (idx *ExplorerIndexer) SetLastUpdateTime(t time.Time) {
	idx.lastUpdateTime = t
}