package keeper

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/namo/x/explorer/types"
)

// TransactionSearchEngine handles comprehensive transaction search and filtering
type TransactionSearchEngine struct {
	keeper Keeper
}

// NewTransactionSearchEngine creates a new transaction search engine
func NewTransactionSearchEngine(keeper Keeper) *TransactionSearchEngine {
	return &TransactionSearchEngine{
		keeper: keeper,
	}
}

// SearchResult represents a comprehensive search result
type SearchResult struct {
	SearchID        string                    `json:"search_id"`
	SearchQuery     types.SearchQuery         `json:"search_query"`
	TotalResults    int64                     `json:"total_results"`
	FilteredResults int64                     `json:"filtered_results"`
	Transactions    []types.EnhancedTransaction `json:"transactions"`
	Blocks          []types.BlockSummary      `json:"blocks"`
	Addresses       []types.AddressInfo       `json:"addresses"`
	Tokens          []types.TokenInfo         `json:"tokens"`
	Modules         []types.ModuleActivity    `json:"modules"`
	Statistics      types.SearchStatistics   `json:"statistics"`
	SearchTime      time.Duration             `json:"search_time"`
	CreatedAt       time.Time                 `json:"created_at"`
}

// AdvancedSearchFilters represents comprehensive search filters
type AdvancedSearchFilters struct {
	// Time filters
	StartTime    *time.Time `json:"start_time,omitempty"`
	EndTime      *time.Time `json:"end_time,omitempty"`
	BlockRange   *types.BlockRange `json:"block_range,omitempty"`
	
	// Transaction filters
	TxTypes      []string `json:"tx_types,omitempty"`
	MessageTypes []string `json:"message_types,omitempty"`
	TxStatus     []string `json:"tx_status,omitempty"` // success, failed, pending
	
	// Address filters
	Addresses    []string `json:"addresses,omitempty"`
	AddressTypes []string `json:"address_types,omitempty"` // user, contract, module
	
	// Amount filters
	AmountRange  *types.AmountRange `json:"amount_range,omitempty"`
	Tokens       []string `json:"tokens,omitempty"`
	
	// Module filters
	Modules      []string `json:"modules,omitempty"`
	ModuleActions []string `json:"module_actions,omitempty"`
	
	// Advanced filters
	GasRange     *types.GasRange `json:"gas_range,omitempty"`
	FeeRange     *types.FeeRange `json:"fee_range,omitempty"`
	Memo         string `json:"memo,omitempty"`
	
	// Cultural filters (unique to DeshChain)
	CulturalEvents []string `json:"cultural_events,omitempty"`
	Festivals      []string `json:"festivals,omitempty"`
	Languages      []string `json:"languages,omitempty"`
	
	// Pagination
	Page         int64 `json:"page"`
	PageSize     int64 `json:"page_size"`
	SortBy       string `json:"sort_by"`
	SortOrder    string `json:"sort_order"` // asc, desc
}

// SearchTransactions performs comprehensive transaction search
func (tse *TransactionSearchEngine) SearchTransactions(ctx sdk.Context, query types.SearchQuery, filters AdvancedSearchFilters) (*SearchResult, error) {
	startTime := time.Now()
	
	// Generate search ID
	searchID := tse.generateSearchID(ctx, query.Query)
	
	result := &SearchResult{
		SearchID:    searchID,
		SearchQuery: query,
		CreatedAt:   ctx.BlockTime(),
	}
	
	// Parse and validate search query
	parsedQuery, err := tse.parseSearchQuery(query.Query)
	if err != nil {
		return nil, fmt.Errorf("failed to parse search query: %w", err)
	}
	
	// Determine search strategy based on query type
	searchStrategy := tse.determineSearchStrategy(parsedQuery, filters)
	
	// Execute search based on strategy
	switch searchStrategy {
	case "hash_search":
		err = tse.executeHashSearch(ctx, result, parsedQuery, filters)
	case "address_search":
		err = tse.executeAddressSearch(ctx, result, parsedQuery, filters)
	case "block_search":
		err = tse.executeBlockSearch(ctx, result, parsedQuery, filters)
	case "advanced_search":
		err = tse.executeAdvancedSearch(ctx, result, parsedQuery, filters)
	case "module_search":
		err = tse.executeModuleSearch(ctx, result, parsedQuery, filters)
	case "cultural_search":
		err = tse.executeCulturalSearch(ctx, result, parsedQuery, filters)
	default:
		err = tse.executeFullTextSearch(ctx, result, parsedQuery, filters)
	}
	
	if err != nil {
		return nil, fmt.Errorf("search execution failed: %w", err)
	}
	
	// Apply additional filters
	tse.applyAdvancedFilters(result, filters)
	
	// Sort results
	tse.sortResults(result, filters.SortBy, filters.SortOrder)
	
	// Apply pagination
	tse.applyPagination(result, filters.Page, filters.PageSize)
	
	// Calculate statistics
	result.Statistics = tse.calculateSearchStatistics(result)
	
	// Record search time
	result.SearchTime = time.Since(startTime)
	
	// Store search result for caching
	tse.keeper.SetSearchResult(ctx, *result)
	
	// Emit search event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTransactionSearch,
			sdk.NewAttribute(types.AttributeKeySearchID, searchID),
			sdk.NewAttribute(types.AttributeKeyQuery, query.Query),
			sdk.NewAttribute(types.AttributeKeyResultCount, fmt.Sprintf("%d", result.FilteredResults)),
			sdk.NewAttribute(types.AttributeKeySearchTime, result.SearchTime.String()),
		),
	)
	
	return result, nil
}

// executeHashSearch searches by transaction hash or block hash
func (tse *TransactionSearchEngine) executeHashSearch(ctx sdk.Context, result *SearchResult, query types.ParsedQuery, filters AdvancedSearchFilters) error {
	hash := query.Value
	
	// Try transaction hash first
	tx, found := tse.keeper.GetTransactionByHash(ctx, hash)
	if found {
		enhancedTx := tse.enhanceTransaction(ctx, tx)
		result.Transactions = []types.EnhancedTransaction{enhancedTx}
		result.TotalResults = 1
		result.FilteredResults = 1
		return nil
	}
	
	// Try block hash
	block, found := tse.keeper.GetBlockByHash(ctx, hash)
	if found {
		blockSummary := tse.createBlockSummary(ctx, block)
		result.Blocks = []types.BlockSummary{blockSummary}
		result.TotalResults = 1
		result.FilteredResults = 1
		
		// Also get transactions in this block
		txs := tse.keeper.GetTransactionsInBlock(ctx, block.Height)
		for _, tx := range txs {
			enhanced := tse.enhanceTransaction(ctx, tx)
			result.Transactions = append(result.Transactions, enhanced)
		}
		return nil
	}
	
	return fmt.Errorf("hash not found: %s", hash)
}

// executeAddressSearch searches by address
func (tse *TransactionSearchEngine) executeAddressSearch(ctx sdk.Context, result *SearchResult, query types.ParsedQuery, filters AdvancedSearchFilters) error {
	address := query.Value
	
	// Get address information
	addressInfo := tse.keeper.GetAddressInfo(ctx, address)
	if addressInfo != nil {
		result.Addresses = []types.AddressInfo{*addressInfo}
	}
	
	// Get transactions involving this address
	txs := tse.keeper.GetTransactionsByAddress(ctx, address, filters.StartTime, filters.EndTime)
	
	// Enhance transactions
	for _, tx := range txs {
		enhanced := tse.enhanceTransaction(ctx, tx)
		result.Transactions = append(result.Transactions, enhanced)
	}
	
	result.TotalResults = int64(len(result.Transactions))
	result.FilteredResults = result.TotalResults
	
	return nil
}

// executeBlockSearch searches by block height or range
func (tse *TransactionSearchEngine) executeBlockSearch(ctx sdk.Context, result *SearchResult, query types.ParsedQuery, filters AdvancedSearchFilters) error {
	// Parse block height or range
	if query.Type == "block_height" {
		height := query.IntValue
		block, found := tse.keeper.GetBlockByHeight(ctx, height)
		if !found {
			return fmt.Errorf("block not found at height: %d", height)
		}
		
		blockSummary := tse.createBlockSummary(ctx, block)
		result.Blocks = []types.BlockSummary{blockSummary}
		result.TotalResults = 1
		result.FilteredResults = 1
		
		// Get transactions in this block
		txs := tse.keeper.GetTransactionsInBlock(ctx, height)
		for _, tx := range txs {
			enhanced := tse.enhanceTransaction(ctx, tx)
			result.Transactions = append(result.Transactions, enhanced)
		}
	} else if filters.BlockRange != nil {
		// Search block range
		blocks := tse.keeper.GetBlocksInRange(ctx, filters.BlockRange.Start, filters.BlockRange.End)
		for _, block := range blocks {
			blockSummary := tse.createBlockSummary(ctx, block)
			result.Blocks = append(result.Blocks, blockSummary)
			
			// Get transactions in each block
			txs := tse.keeper.GetTransactionsInBlock(ctx, block.Height)
			for _, tx := range txs {
				enhanced := tse.enhanceTransaction(ctx, tx)
				result.Transactions = append(result.Transactions, enhanced)
			}
		}
		
		result.TotalResults = int64(len(result.Blocks))
		result.FilteredResults = result.TotalResults
	}
	
	return nil
}

// executeAdvancedSearch performs advanced multi-criteria search
func (tse *TransactionSearchEngine) executeAdvancedSearch(ctx sdk.Context, result *SearchResult, query types.ParsedQuery, filters AdvancedSearchFilters) error {
	// Build search criteria
	criteria := types.SearchCriteria{
		StartTime:     filters.StartTime,
		EndTime:       filters.EndTime,
		TxTypes:       filters.TxTypes,
		MessageTypes:  filters.MessageTypes,
		Addresses:     filters.Addresses,
		Modules:       filters.Modules,
		AmountRange:   filters.AmountRange,
		GasRange:      filters.GasRange,
		FeeRange:      filters.FeeRange,
	}
	
	// Execute multi-criteria search
	txs := tse.keeper.SearchTransactionsByCriteria(ctx, criteria)
	
	// Enhance and filter results
	for _, tx := range txs {
		enhanced := tse.enhanceTransaction(ctx, tx)
		if tse.matchesAdvancedFilters(enhanced, filters) {
			result.Transactions = append(result.Transactions, enhanced)
		}
	}
	
	result.TotalResults = int64(len(txs))
	result.FilteredResults = int64(len(result.Transactions))
	
	return nil
}

// executeModuleSearch searches module-specific transactions
func (tse *TransactionSearchEngine) executeModuleSearch(ctx sdk.Context, result *SearchResult, query types.ParsedQuery, filters AdvancedSearchFilters) error {
	moduleName := query.Value
	
	// Get module activity
	moduleActivity := tse.keeper.GetModuleActivity(ctx, moduleName, filters.StartTime, filters.EndTime)
	if moduleActivity != nil {
		result.Modules = []types.ModuleActivity{*moduleActivity}
	}
	
	// Get module-specific transactions
	txs := tse.keeper.GetTransactionsByModule(ctx, moduleName, filters.StartTime, filters.EndTime)
	
	// Filter by module actions if specified
	for _, tx := range txs {
		enhanced := tse.enhanceTransaction(ctx, tx)
		if len(filters.ModuleActions) == 0 || tse.containsModuleAction(enhanced, filters.ModuleActions) {
			result.Transactions = append(result.Transactions, enhanced)
		}
	}
	
	result.TotalResults = int64(len(result.Transactions))
	result.FilteredResults = result.TotalResults
	
	return nil
}

// executeCulturalSearch searches cultural and festival-related transactions
func (tse *TransactionSearchEngine) executeCulturalSearch(ctx sdk.Context, result *SearchResult, query types.ParsedQuery, filters AdvancedSearchFilters) error {
	// Search cultural module transactions
	culturalTxs := tse.keeper.GetCulturalTransactions(ctx, filters.CulturalEvents, filters.Festivals, filters.Languages)
	
	// Enhanced cultural search
	for _, tx := range culturalTxs {
		enhanced := tse.enhanceTransaction(ctx, tx)
		if tse.matchesCulturalFilters(enhanced, filters) {
			result.Transactions = append(result.Transactions, enhanced)
		}
	}
	
	result.TotalResults = int64(len(result.Transactions))
	result.FilteredResults = result.TotalResults
	
	return nil
}

// executeFullTextSearch performs full-text search across all indexed content
func (tse *TransactionSearchEngine) executeFullTextSearch(ctx sdk.Context, result *SearchResult, query types.ParsedQuery, filters AdvancedSearchFilters) error {
	searchTerms := strings.Fields(query.Value)
	
	// Search across different content types
	txs := tse.keeper.FullTextSearchTransactions(ctx, searchTerms)
	addresses := tse.keeper.FullTextSearchAddresses(ctx, searchTerms)
	tokens := tse.keeper.FullTextSearchTokens(ctx, searchTerms)
	
	// Enhance transaction results
	for _, tx := range txs {
		enhanced := tse.enhanceTransaction(ctx, tx)
		result.Transactions = append(result.Transactions, enhanced)
	}
	
	result.Addresses = addresses
	result.Tokens = tokens
	result.TotalResults = int64(len(txs) + len(addresses) + len(tokens))
	result.FilteredResults = result.TotalResults
	
	return nil
}

// Helper functions

func (tse *TransactionSearchEngine) parseSearchQuery(query string) (types.ParsedQuery, error) {
	parsed := types.ParsedQuery{
		Original: query,
		Value:    strings.TrimSpace(query),
	}
	
	// Detect query type
	if tse.isTransactionHash(query) {
		parsed.Type = "tx_hash"
	} else if tse.isBlockHash(query) {
		parsed.Type = "block_hash"
	} else if tse.isAddress(query) {
		parsed.Type = "address"
	} else if tse.isBlockHeight(query) {
		parsed.Type = "block_height"
		// Parse block height
		if height, err := sdk.ParseUint(query); err == nil {
			parsed.IntValue = int64(height)
		}
	} else if tse.isModuleName(query) {
		parsed.Type = "module"
	} else {
		parsed.Type = "full_text"
	}
	
	return parsed, nil
}

func (tse *TransactionSearchEngine) determineSearchStrategy(query types.ParsedQuery, filters AdvancedSearchFilters) string {
	// Determine search strategy based on query type and filters
	if query.Type == "tx_hash" || query.Type == "block_hash" {
		return "hash_search"
	}
	if query.Type == "address" {
		return "address_search"
	}
	if query.Type == "block_height" || filters.BlockRange != nil {
		return "block_search"
	}
	if query.Type == "module" || len(filters.Modules) > 0 {
		return "module_search"
	}
	if len(filters.CulturalEvents) > 0 || len(filters.Festivals) > 0 {
		return "cultural_search"
	}
	if tse.hasAdvancedFilters(filters) {
		return "advanced_search"
	}
	return "full_text_search"
}

func (tse *TransactionSearchEngine) enhanceTransaction(ctx sdk.Context, tx types.Transaction) types.EnhancedTransaction {
	enhanced := types.EnhancedTransaction{
		Transaction: tx,
		Timestamp:   ctx.BlockTime(),
	}
	
	// Add block information
	if block, found := tse.keeper.GetBlockByHeight(ctx, tx.Height); found {
		enhanced.BlockHash = block.Hash
		enhanced.BlockTime = block.Time
	}
	
	// Add gas information
	enhanced.GasUsed = tx.GasUsed
	enhanced.GasWanted = tx.GasWanted
	enhanced.GasPrice = tse.calculateGasPrice(tx)
	
	// Add fee information
	enhanced.Fees = tx.Fees
	enhanced.FeePayer = tx.FeePayer
	
	// Add message details
	enhanced.MessageTypes = tse.extractMessageTypes(tx)
	enhanced.MessageCount = int64(len(tx.Messages))
	
	// Add module interactions
	enhanced.ModulesInvolved = tse.extractModulesInvolved(tx)
	
	// Add token transfers
	enhanced.TokenTransfers = tse.extractTokenTransfers(ctx, tx)
	
	// Add cultural context if applicable
	enhanced.CulturalContext = tse.extractCulturalContext(ctx, tx)
	
	// Add success/failure details
	enhanced.Success = tx.Code == 0
	if !enhanced.Success {
		enhanced.ErrorMessage = tx.RawLog
	}
	
	return enhanced
}

func (tse *TransactionSearchEngine) createBlockSummary(ctx sdk.Context, block types.Block) types.BlockSummary {
	summary := types.BlockSummary{
		Height:          block.Height,
		Hash:            block.Hash,
		Time:            block.Time,
		ProposerAddress: block.ProposerAddress,
		TransactionCount: int64(len(block.Transactions)),
	}
	
	// Calculate block statistics
	var totalGasUsed, totalFees uint64
	for _, tx := range block.Transactions {
		totalGasUsed += tx.GasUsed
		for _, fee := range tx.Fees {
			if fee.Denom == "namo" {
				totalFees += fee.Amount.Uint64()
			}
		}
	}
	
	summary.TotalGasUsed = totalGasUsed
	summary.TotalFees = sdk.NewCoin("namo", sdk.NewIntFromUint64(totalFees))
	
	// Add validator information
	if validator, found := tse.keeper.GetValidatorByAddress(ctx, block.ProposerAddress); found {
		summary.ValidatorName = validator.Description.Moniker
		summary.ValidatorOperator = validator.OperatorAddress
	}
	
	return summary
}

func (tse *TransactionSearchEngine) applyAdvancedFilters(result *SearchResult, filters AdvancedSearchFilters) {
	// Filter transactions
	filteredTxs := []types.EnhancedTransaction{}
	for _, tx := range result.Transactions {
		if tse.matchesAdvancedFilters(tx, filters) {
			filteredTxs = append(filteredTxs, tx)
		}
	}
	result.Transactions = filteredTxs
	result.FilteredResults = int64(len(filteredTxs))
}

func (tse *TransactionSearchEngine) matchesAdvancedFilters(tx types.EnhancedTransaction, filters AdvancedSearchFilters) bool {
	// Time filters
	if filters.StartTime != nil && tx.Timestamp.Before(*filters.StartTime) {
		return false
	}
	if filters.EndTime != nil && tx.Timestamp.After(*filters.EndTime) {
		return false
	}
	
	// Transaction type filters
	if len(filters.TxTypes) > 0 && !tse.containsString(filters.TxTypes, tx.Transaction.Type) {
		return false
	}
	
	// Message type filters
	if len(filters.MessageTypes) > 0 {
		hasMatchingType := false
		for _, msgType := range tx.MessageTypes {
			if tse.containsString(filters.MessageTypes, msgType) {
				hasMatchingType = true
				break
			}
		}
		if !hasMatchingType {
			return false
		}
	}
	
	// Status filters
	if len(filters.TxStatus) > 0 {
		status := "success"
		if !tx.Success {
			status = "failed"
		}
		if !tse.containsString(filters.TxStatus, status) {
			return false
		}
	}
	
	// Address filters
	if len(filters.Addresses) > 0 {
		hasMatchingAddress := false
		for _, addr := range filters.Addresses {
			if tse.transactionInvolvesAddress(tx, addr) {
				hasMatchingAddress = true
				break
			}
		}
		if !hasMatchingAddress {
			return false
		}
	}
	
	// Amount filters
	if filters.AmountRange != nil {
		totalAmount := tse.calculateTotalTransactionAmount(tx)
		if totalAmount.LT(filters.AmountRange.Min) || totalAmount.GT(filters.AmountRange.Max) {
			return false
		}
	}
	
	// Gas filters
	if filters.GasRange != nil {
		if tx.GasUsed < filters.GasRange.Min || tx.GasUsed > filters.GasRange.Max {
			return false
		}
	}
	
	// Fee filters
	if filters.FeeRange != nil {
		totalFee := tse.calculateTotalFee(tx.Fees)
		if totalFee.LT(filters.FeeRange.Min) || totalFee.GT(filters.FeeRange.Max) {
			return false
		}
	}
	
	// Memo filter
	if filters.Memo != "" {
		if !strings.Contains(strings.ToLower(tx.Transaction.Memo), strings.ToLower(filters.Memo)) {
			return false
		}
	}
	
	return true
}

func (tse *TransactionSearchEngine) sortResults(result *SearchResult, sortBy, sortOrder string) {
	if sortBy == "" {
		sortBy = "timestamp"
	}
	if sortOrder == "" {
		sortOrder = "desc"
	}
	
	// Sort transactions
	sort.Slice(result.Transactions, func(i, j int) bool {
		var less bool
		switch sortBy {
		case "timestamp":
			less = result.Transactions[i].Timestamp.Before(result.Transactions[j].Timestamp)
		case "height":
			less = result.Transactions[i].Height < result.Transactions[j].Height
		case "gas_used":
			less = result.Transactions[i].GasUsed < result.Transactions[j].GasUsed
		case "fee":
			feeI := tse.calculateTotalFee(result.Transactions[i].Fees)
			feeJ := tse.calculateTotalFee(result.Transactions[j].Fees)
			less = feeI.LT(feeJ)
		default:
			less = result.Transactions[i].Timestamp.Before(result.Transactions[j].Timestamp)
		}
		
		if sortOrder == "desc" {
			return !less
		}
		return less
	})
	
	// Sort blocks
	sort.Slice(result.Blocks, func(i, j int) bool {
		less := result.Blocks[i].Height < result.Blocks[j].Height
		if sortOrder == "desc" {
			return !less
		}
		return less
	})
}

func (tse *TransactionSearchEngine) applyPagination(result *SearchResult, page, pageSize int64) {
	if pageSize <= 0 {
		pageSize = 50 // Default page size
	}
	if page <= 0 {
		page = 1
	}
	
	// Paginate transactions
	start := (page - 1) * pageSize
	end := start + pageSize
	
	if start < int64(len(result.Transactions)) {
		if end > int64(len(result.Transactions)) {
			end = int64(len(result.Transactions))
		}
		result.Transactions = result.Transactions[start:end]
	} else {
		result.Transactions = []types.EnhancedTransaction{}
	}
	
	// Paginate blocks
	if start < int64(len(result.Blocks)) {
		if end > int64(len(result.Blocks)) {
			end = int64(len(result.Blocks))
		}
		result.Blocks = result.Blocks[start:end]
	} else {
		result.Blocks = []types.BlockSummary{}
	}
}

func (tse *TransactionSearchEngine) calculateSearchStatistics(result *SearchResult) types.SearchStatistics {
	stats := types.SearchStatistics{
		TotalTransactions: int64(len(result.Transactions)),
		TotalBlocks:       int64(len(result.Blocks)),
		TotalAddresses:    int64(len(result.Addresses)),
		TotalTokens:       int64(len(result.Tokens)),
		TotalModules:      int64(len(result.Modules)),
	}
	
	// Calculate transaction statistics
	var totalGasUsed uint64
	var totalFees sdk.Int = sdk.ZeroInt()
	successCount := int64(0)
	
	for _, tx := range result.Transactions {
		totalGasUsed += tx.GasUsed
		totalFees = totalFees.Add(tse.calculateTotalFee(tx.Fees))
		if tx.Success {
			successCount++
		}
	}
	
	stats.TotalGasUsed = totalGasUsed
	stats.TotalFees = sdk.NewCoin("namo", totalFees)
	stats.SuccessRate = float64(successCount) / float64(len(result.Transactions))
	
	if len(result.Transactions) > 0 {
		stats.AverageGasUsed = totalGasUsed / uint64(len(result.Transactions))
		stats.AverageFee = sdk.NewCoin("namo", totalFees.QuoRaw(int64(len(result.Transactions))))
	}
	
	return stats
}

// Utility functions
func (tse *TransactionSearchEngine) generateSearchID(ctx sdk.Context, query string) string {
	timestamp := ctx.BlockTime().Unix()
	return fmt.Sprintf("SEARCH-%d-%s", timestamp, tse.hashString(query)[:8])
}

func (tse *TransactionSearchEngine) hashString(s string) string {
	// Simple hash function for demo
	return fmt.Sprintf("%x", sdk.AccAddress(s).String())[:16]
}

func (tse *TransactionSearchEngine) isTransactionHash(s string) bool {
	// Check if string matches transaction hash pattern (64 hex characters)
	matched, _ := regexp.MatchString("^[0-9a-fA-F]{64}$", s)
	return matched
}

func (tse *TransactionSearchEngine) isBlockHash(s string) bool {
	// Check if string matches block hash pattern (64 hex characters)
	matched, _ := regexp.MatchString("^[0-9a-fA-F]{64}$", s)
	return matched
}

func (tse *TransactionSearchEngine) isAddress(s string) bool {
	// Check if string matches address pattern
	_, err := sdk.AccAddressFromBech32(s)
	return err == nil
}

func (tse *TransactionSearchEngine) isBlockHeight(s string) bool {
	// Check if string is a valid block height (positive integer)
	matched, _ := regexp.MatchString("^[1-9][0-9]*$", s)
	return matched
}

func (tse *TransactionSearchEngine) isModuleName(s string) bool {
	// Check if string matches known module names
	knownModules := []string{"bank", "staking", "gov", "cultural", "namo", "krishimitra", "vyavasayamitra", "shikshaamitra"}
	return tse.containsString(knownModules, s)
}

func (tse *TransactionSearchEngine) hasAdvancedFilters(filters AdvancedSearchFilters) bool {
	return len(filters.TxTypes) > 0 || len(filters.MessageTypes) > 0 || 
		   len(filters.Addresses) > 0 || filters.AmountRange != nil ||
		   filters.GasRange != nil || filters.FeeRange != nil ||
		   filters.StartTime != nil || filters.EndTime != nil
}

func (tse *TransactionSearchEngine) containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func (tse *TransactionSearchEngine) containsModuleAction(tx types.EnhancedTransaction, actions []string) bool {
	for _, action := range actions {
		for _, msgType := range tx.MessageTypes {
			if strings.Contains(msgType, action) {
				return true
			}
		}
	}
	return false
}

func (tse *TransactionSearchEngine) matchesCulturalFilters(tx types.EnhancedTransaction, filters AdvancedSearchFilters) bool {
	if len(filters.CulturalEvents) > 0 {
		if tx.CulturalContext == nil || !tse.containsString(filters.CulturalEvents, tx.CulturalContext.Event) {
			return false
		}
	}
	if len(filters.Festivals) > 0 {
		if tx.CulturalContext == nil || !tse.containsString(filters.Festivals, tx.CulturalContext.Festival) {
			return false
		}
	}
	return true
}

func (tse *TransactionSearchEngine) transactionInvolvesAddress(tx types.EnhancedTransaction, address string) bool {
	// Check if transaction involves the specified address
	if tx.Transaction.Sender == address {
		return true
	}
	
	// Check in token transfers
	for _, transfer := range tx.TokenTransfers {
		if transfer.From == address || transfer.To == address {
			return true
		}
	}
	
	return false
}

func (tse *TransactionSearchEngine) calculateTotalTransactionAmount(tx types.EnhancedTransaction) sdk.Int {
	total := sdk.ZeroInt()
	for _, transfer := range tx.TokenTransfers {
		total = total.Add(transfer.Amount.Amount)
	}
	return total
}

func (tse *TransactionSearchEngine) calculateTotalFee(fees []sdk.Coin) sdk.Int {
	total := sdk.ZeroInt()
	for _, fee := range fees {
		if fee.Denom == "namo" {
			total = total.Add(fee.Amount)
		}
	}
	return total
}

func (tse *TransactionSearchEngine) calculateGasPrice(tx types.Transaction) sdk.Dec {
	if tx.GasWanted == 0 {
		return sdk.ZeroDec()
	}
	totalFee := tse.calculateTotalFee(tx.Fees)
	return totalFee.ToDec().QuoInt64(int64(tx.GasWanted))
}

func (tse *TransactionSearchEngine) extractMessageTypes(tx types.Transaction) []string {
	types := []string{}
	for _, msg := range tx.Messages {
		types = append(types, msg.Type)
	}
	return types
}

func (tse *TransactionSearchEngine) extractModulesInvolved(tx types.Transaction) []string {
	modules := []string{}
	moduleSet := make(map[string]bool)
	
	for _, msg := range tx.Messages {
		// Extract module name from message type
		parts := strings.Split(msg.Type, ".")
		if len(parts) > 0 {
			module := parts[0]
			if !moduleSet[module] {
				modules = append(modules, module)
				moduleSet[module] = true
			}
		}
	}
	
	return modules
}

func (tse *TransactionSearchEngine) extractTokenTransfers(ctx sdk.Context, tx types.Transaction) []types.TokenTransfer {
	transfers := []types.TokenTransfer{}
	
	// Extract transfers from bank messages
	for _, msg := range tx.Messages {
		if msg.Type == "/cosmos.bank.v1beta1.MsgSend" {
			// Parse bank send message
			transfer := types.TokenTransfer{
				From:   msg.From,
				To:     msg.To,
				Amount: msg.Amount,
				Denom:  msg.Amount.Denom,
			}
			transfers = append(transfers, transfer)
		}
	}
	
	return transfers
}

func (tse *TransactionSearchEngine) extractCulturalContext(ctx sdk.Context, tx types.Transaction) *types.CulturalContext {
	// Check if transaction has cultural context
	for _, msg := range tx.Messages {
		if strings.Contains(msg.Type, "cultural") {
			return &types.CulturalContext{
				Event:    "cultural_transaction",
				Festival: "unknown",
				Language: "hindi",
			}
		}
	}
	return nil
}