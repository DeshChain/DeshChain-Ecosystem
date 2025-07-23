package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/deshchain/x/remittance/types"
)

// USDCorridorManager manages USD-specific remittance corridors for DUSD integration
type USDCorridorManager struct {
	keeper *Keeper
}

// NewUSDCorridorManager creates a new USD corridor manager
func NewUSDCorridorManager(keeper *Keeper) *USDCorridorManager {
	return &USDCorridorManager{
		keeper: keeper,
	}
}

// Enhanced Remittance Transfer with Multi-Currency Support
type EnhancedRemittanceTransfer struct {
	// Existing transfer fields
	types.RemittanceTransfer
	
	// Multi-currency enhancement fields
	SourceCurrency       string            `json:"source_currency"`       // "USD", "EUR", "SGD"
	SourceAmount         sdk.Coin          `json:"source_amount"`         // Original amount
	RoutingCurrency      string            `json:"routing_currency"`      // "DUSD", "DEUR", "DSGD"
	RoutingAmount        sdk.Coin          `json:"routing_amount"`        // Stablecoin routing amount
	DestinationCurrency  string            `json:"destination_currency"`  // "DINR"
	DestinationAmount    sdk.Coin          `json:"destination_amount"`    // Final amount
	
	// Enhanced cost analysis
	TraditionalCost      sdk.Coin          `json:"traditional_cost"`      // Cost via traditional banking
	DeshChainCost        sdk.Coin          `json:"deshchain_cost"`        // Cost via DeshChain
	TotalSavings         sdk.Coin          `json:"total_savings"`         // Customer savings
	ProcessingTime       time.Duration     `json:"processing_time"`       // vs 1-3 days traditional
	
	// Routing optimization
	OptimalRoute         []string          `json:"optimal_route"`         // Currency conversion path
	ExchangeRates        map[string]sdk.Dec `json:"exchange_rates"`       // All rates used
	CorridorEfficiency   sdk.Dec           `json:"corridor_efficiency"`   // Efficiency score 0-1
}

// OptimizeRemittanceCorridor finds the optimal stablecoin routing currency
func (ucm *USDCorridorManager) OptimizeRemittanceCorridor(
	ctx sdk.Context,
	sourceCurrency string,
	destinationCurrency string,
	amount sdk.Coin,
) (string, *types.CorridorOptimization, error) {
	
	// Available stablecoin corridors
	corridors := []string{"DUSD", "DEUR", "DSGD", "DGBP"}
	
	bestCorridor := ""
	lowestCost := sdk.NewCoin("namo", sdk.NewInt(1000000000)) // Initialize with high value
	var bestOptimization *types.CorridorOptimization
	
	for _, corridor := range corridors {
		optimization, err := ucm.AnalyzeCorridorEfficiency(ctx, sourceCurrency, corridor, destinationCurrency, amount)
		if err != nil {
			ucm.keeper.logger.Error("failed to analyze corridor",
				"source", sourceCurrency,
				"corridor", corridor,
				"destination", destinationCurrency,
				"error", err,
			)
			continue
		}
		
		if optimization.TotalCost.Amount.LT(lowestCost.Amount) {
			lowestCost = optimization.TotalCost
			bestCorridor = corridor
			bestOptimization = optimization
		}
	}
	
	if bestCorridor == "" {
		return "", nil, fmt.Errorf("no suitable corridor found for %s to %s", sourceCurrency, destinationCurrency)
	}
	
	return bestCorridor, bestOptimization, nil
}

// AnalyzeCorridorEfficiency analyzes the efficiency of a specific corridor
func (ucm *USDCorridorManager) AnalyzeCorridorEfficiency(
	ctx sdk.Context,
	sourceCurrency, corridorCurrency, destinationCurrency string,
	amount sdk.Coin,
) (*types.CorridorOptimization, error) {
	
	// Get exchange rates
	sourceToCorridorRate, err := ucm.getExchangeRate(ctx, sourceCurrency, corridorCurrency)
	if err != nil {
		return nil, fmt.Errorf("failed to get %s/%s rate: %w", sourceCurrency, corridorCurrency, err)
	}
	
	corridorToDestinationRate, err := ucm.getExchangeRate(ctx, corridorCurrency, destinationCurrency)
	if err != nil {
		return nil, fmt.Errorf("failed to get %s/%s rate: %w", corridorCurrency, destinationCurrency, err)
	}
	
	// Calculate amounts at each step
	corridorAmount := sdk.NewCoin(corridorCurrency, amount.Amount.ToLegacyDec().Mul(sourceToCorridorRate).TruncateInt())
	destinationAmount := sdk.NewCoin(destinationCurrency, corridorAmount.Amount.ToLegacyDec().Mul(corridorToDestinationRate).TruncateInt())
	
	// Calculate fees (0.25% for each conversion + 0.05% corridor fee)
	conversionFeeRate := sdk.NewDecWithPrec(25, 4)  // 0.25%
	corridorFeeRate := sdk.NewDecWithPrec(5, 4)     // 0.05%
	
	sourceFee := amount.Amount.ToLegacyDec().Mul(conversionFeeRate)
	corridorFee := corridorAmount.Amount.ToLegacyDec().Mul(corridorFeeRate)
	destinationFee := corridorAmount.Amount.ToLegacyDec().Mul(conversionFeeRate)
	
	totalFee := sourceFee.Add(corridorFee).Add(destinationFee)
	
	// Calculate traditional cost for comparison (6-8% typical)
	traditionalFeeRate := sdk.NewDecWithPrec(7, 2) // 7%
	traditionalCost := amount.Amount.ToLegacyDec().Mul(traditionalFeeRate)
	
	// Calculate processing time (instant vs 1-3 days)
	processingTime := 30 * time.Second // Near-instant
	
	// Calculate efficiency score
	costEfficiency := sdk.OneDec().Sub(totalFee.Quo(traditionalCost))
	timeEfficiency := sdk.NewDecWithPrec(99, 2) // 99% time savings
	
	// Overall efficiency (weighted average)
	efficiency := costEfficiency.Mul(sdk.NewDecWithPrec(7, 1)).Add(timeEfficiency.Mul(sdk.NewDecWithPrec(3, 1)))
	
	optimization := &types.CorridorOptimization{
		SourceCurrency:         sourceCurrency,
		CorridorCurrency:      corridorCurrency,
		DestinationCurrency:   destinationCurrency,
		SourceAmount:          amount,
		CorridorAmount:        corridorAmount,
		DestinationAmount:     destinationAmount,
		SourceToCorridorRate:  sourceToCorridorRate,
		CorridorToDestRate:    corridorToDestinationRate,
		TotalCost:            sdk.NewCoin(amount.Denom, totalFee.TruncateInt()),
		TraditionalCost:      sdk.NewCoin(amount.Denom, traditionalCost.TruncateInt()),
		Savings:              sdk.NewCoin(amount.Denom, traditionalCost.Sub(totalFee).TruncateInt()),
		ProcessingTime:       processingTime,
		EfficiencyScore:      efficiency,
		Route:                []string{sourceCurrency, corridorCurrency, destinationCurrency},
	}
	
	return optimization, nil
}

// ProcessEnhancedRemittance processes a remittance with multi-currency optimization
func (ucm *USDCorridorManager) ProcessEnhancedRemittance(
	ctx sdk.Context,
	transfer types.RemittanceTransfer,
	sourceCurrency string,
) (*EnhancedRemittanceTransfer, error) {
	
	// 1. Validate basic transfer
	if err := ucm.validateBasicTransfer(ctx, &transfer); err != nil {
		return nil, fmt.Errorf("transfer validation failed: %w", err)
	}
	
	// 2. Optimize corridor
	optimalCorridor, optimization, err := ucm.OptimizeRemittanceCorridor(
		ctx, 
		sourceCurrency, 
		"DINR", 
		transfer.Amount,
	)
	if err != nil {
		return nil, fmt.Errorf("corridor optimization failed: %w", err)
	}
	
	// 3. Create enhanced transfer
	enhancedTransfer := &EnhancedRemittanceTransfer{
		RemittanceTransfer:   transfer,
		SourceCurrency:       sourceCurrency,
		SourceAmount:         transfer.Amount,
		RoutingCurrency:      optimalCorridor,
		RoutingAmount:        optimization.CorridorAmount,
		DestinationCurrency:  "DINR",
		DestinationAmount:    optimization.DestinationAmount,
		TraditionalCost:      optimization.TraditionalCost,
		DeshChainCost:        optimization.TotalCost,
		TotalSavings:         optimization.Savings,
		ProcessingTime:       optimization.ProcessingTime,
		OptimalRoute:         optimization.Route,
		ExchangeRates:        make(map[string]sdk.Dec),
		CorridorEfficiency:   optimization.EfficiencyScore,
	}
	
	// 4. Store exchange rates used
	enhancedTransfer.ExchangeRates[fmt.Sprintf("%s/%s", sourceCurrency, optimalCorridor)] = optimization.SourceToCorridorRate
	enhancedTransfer.ExchangeRates[fmt.Sprintf("%s/DINR", optimalCorridor)] = optimization.CorridorToDestRate
	
	// 5. Execute the transfer
	if err := ucm.executeMultiCurrencyTransfer(ctx, enhancedTransfer); err != nil {
		return nil, fmt.Errorf("transfer execution failed: %w", err)
	}
	
	// 6. Update transfer status
	enhancedTransfer.Status = "completed"
	enhancedTransfer.UpdatedAt = ctx.BlockTime()
	
	// 7. Store enhanced transfer data
	if err := ucm.storeEnhancedTransfer(ctx, enhancedTransfer); err != nil {
		return nil, fmt.Errorf("failed to store enhanced transfer: %w", err)
	}
	
	// 8. Emit enhanced remittance event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"enhanced_remittance_processed",
			sdk.NewAttribute("transfer_id", transfer.Id),
			sdk.NewAttribute("source_currency", sourceCurrency),
			sdk.NewAttribute("routing_currency", optimalCorridor),
			sdk.NewAttribute("source_amount", transfer.Amount.String()),
			sdk.NewAttribute("destination_amount", optimization.DestinationAmount.String()),
			sdk.NewAttribute("total_savings", optimization.Savings.String()),
			sdk.NewAttribute("processing_time", optimization.ProcessingTime.String()),
			sdk.NewAttribute("efficiency_score", optimization.EfficiencyScore.String()),
		),
	)
	
	ucm.keeper.logger.Info("processed enhanced remittance",
		"transfer_id", transfer.Id,
		"source_currency", sourceCurrency,
		"routing_currency", optimalCorridor,
		"total_savings", optimization.Savings.String(),
		"efficiency_score", optimization.EfficiencyScore.String(),
	)
	
	return enhancedTransfer, nil
}

// executeMultiCurrencyTransfer executes the actual currency conversions
func (ucm *USDCorridorManager) executeMultiCurrencyTransfer(ctx sdk.Context, transfer *EnhancedRemittanceTransfer) error {
	// Step 1: Convert source currency to routing currency
	if err := ucm.convertCurrency(ctx, transfer.SourceCurrency, transfer.RoutingCurrency, transfer.SourceAmount); err != nil {
		return fmt.Errorf("source to routing conversion failed: %w", err)
	}
	
	// Step 2: Route through stablecoin network
	if err := ucm.routeThroughStablecoin(ctx, transfer.RoutingCurrency, transfer.RoutingAmount); err != nil {
		return fmt.Errorf("stablecoin routing failed: %w", err)
	}
	
	// Step 3: Convert routing currency to destination currency
	if err := ucm.convertCurrency(ctx, transfer.RoutingCurrency, transfer.DestinationCurrency, transfer.RoutingAmount); err != nil {
		return fmt.Errorf("routing to destination conversion failed: %w", err)
	}
	
	// Step 4: Handle Sewa Mitra delivery if applicable
	if transfer.SewaMitraId != "" {
		if err := ucm.processSewaMitraDelivery(ctx, transfer); err != nil {
			return fmt.Errorf("Sewa Mitra delivery failed: %w", err)
		}
	}
	
	return nil
}

// convertCurrency handles currency conversion
func (ucm *USDCorridorManager) convertCurrency(ctx sdk.Context, fromCurrency, toCurrency string, amount sdk.Coin) error {
	ucm.keeper.logger.Info("converting currency",
		"from", fromCurrency,
		"to", toCurrency,
		"amount", amount.String(),
	)
	
	// In full implementation, this would:
	// 1. Lock source currency
	// 2. Get real-time exchange rate
	// 3. Execute conversion through DEX or oracle
	// 4. Mint/burn stablecoins as needed
	// 5. Transfer converted amount
	
	return nil
}

// routeThroughStablecoin routes payment through stablecoin network
func (ucm *USDCorridorManager) routeThroughStablecoin(ctx sdk.Context, stablecoin string, amount sdk.Coin) error {
	ucm.keeper.logger.Info("routing through stablecoin",
		"stablecoin", stablecoin,
		"amount", amount.String(),
	)
	
	// In full implementation, this would:
	// 1. Lock stablecoin collateral
	// 2. Route through fastest path (direct or via liquidity pools)
	// 3. Handle any necessary rebalancing
	// 4. Ensure atomic execution
	
	return nil
}

// processSewaMitraDelivery handles Sewa Mitra delivery for USD corridors
func (ucm *USDCorridorManager) processSewaMitraDelivery(ctx sdk.Context, transfer *EnhancedRemittanceTransfer) error {
	ucm.keeper.logger.Info("processing Sewa Mitra delivery",
		"transfer_id", transfer.Id,
		"sewa_mitra_id", transfer.SewaMitraId,
		"destination_amount", transfer.DestinationAmount.String(),
	)
	
	// In full implementation, this would:
	// 1. Notify Sewa Mitra agent
	// 2. Lock funds in escrow
	// 3. Generate pickup code
	// 4. Handle completion confirmation
	// 5. Release funds to agent
	
	return nil
}

// getExchangeRate gets exchange rate between currencies (placeholder for oracle integration)
func (ucm *USDCorridorManager) getExchangeRate(ctx sdk.Context, fromCurrency, toCurrency string) (sdk.Dec, error) {
	// Sample exchange rates for USD corridors
	rates := map[string]sdk.Dec{
		"USD/DUSD": sdk.OneDec(),                    // 1:1 peg
		"DUSD/DINR": sdk.NewDecWithPrec(835, 1),    // 83.5 DINR per DUSD
		"EUR/DEUR": sdk.OneDec(),                    // 1:1 peg
		"DEUR/DINR": sdk.NewDecWithPrec(912, 1),    // 91.2 DINR per DEUR
		"SGD/DSGD": sdk.OneDec(),                    // 1:1 peg
		"DSGD/DINR": sdk.NewDecWithPrec(621, 1),    // 62.1 DINR per DSGD
		"GBP/DGBP": sdk.OneDec(),                    // 1:1 peg
		"DGBP/DINR": sdk.NewDecWithPrec(1055, 1),   // 105.5 DINR per DGBP
	}
	
	rateKey := fmt.Sprintf("%s/%s", fromCurrency, toCurrency)
	if rate, found := rates[rateKey]; found {
		return rate, nil
	}
	
	return sdk.ZeroDec(), fmt.Errorf("exchange rate not available for %s/%s", fromCurrency, toCurrency)
}

// validateBasicTransfer validates basic transfer requirements
func (ucm *USDCorridorManager) validateBasicTransfer(ctx sdk.Context, transfer *types.RemittanceTransfer) error {
	// Validate sender
	if transfer.SenderId == "" {
		return fmt.Errorf("sender ID cannot be empty")
	}
	
	// Validate recipient
	if transfer.RecipientId == "" {
		return fmt.Errorf("recipient ID cannot be empty")
	}
	
	// Validate amount
	if transfer.Amount.Amount.IsZero() || transfer.Amount.Amount.IsNegative() {
		return fmt.Errorf("invalid transfer amount: %s", transfer.Amount.String())
	}
	
	// Validate currency support
	supportedCurrencies := []string{"USD", "EUR", "SGD", "GBP", "DINR"}
	currencySupported := false
	for _, supported := range supportedCurrencies {
		if transfer.Amount.Denom == supported {
			currencySupported = true
			break
		}
	}
	if !currencySupported {
		return fmt.Errorf("unsupported currency: %s", transfer.Amount.Denom)
	}
	
	return nil
}

// storeEnhancedTransfer stores enhanced transfer data
func (ucm *USDCorridorManager) storeEnhancedTransfer(ctx sdk.Context, transfer *EnhancedRemittanceTransfer) error {
	store := ucm.keeper.storeService.OpenKVStore(ctx)
	key := types.GetEnhancedTransferKey(transfer.Id)
	
	bz, err := ucm.keeper.cdc.Marshal(transfer)
	if err != nil {
		return err
	}
	
	return store.Set(key, bz)
}

// GetEnhancedTransfer retrieves enhanced transfer data
func (ucm *USDCorridorManager) GetEnhancedTransfer(ctx sdk.Context, transferID string) (*EnhancedRemittanceTransfer, bool) {
	store := ucm.keeper.storeService.OpenKVStore(ctx)
	key := types.GetEnhancedTransferKey(transferID)
	
	bz, err := store.Get(key)
	if err != nil || bz == nil {
		return nil, false
	}
	
	var enhancedTransfer EnhancedRemittanceTransfer
	if err := ucm.keeper.cdc.Unmarshal(bz, &enhancedTransfer); err != nil {
		return nil, false
	}
	
	return &enhancedTransfer, true
}

// GetCorridorStats returns statistics for specific USD corridors
func (ucm *USDCorridorManager) GetCorridorStats(ctx sdk.Context) map[string]types.CorridorStats {
	// This would aggregate real statistics from enhanced transfers
	// For now, return sample data
	
	return map[string]types.CorridorStats{
		"USD-INR": {
			Corridor:              "USD-INR",
			TotalVolume:           sdk.NewCoin("USD", sdk.NewInt(100000000)), // $100M
			TotalTransactions:     10000,
			AverageAmount:         sdk.NewCoin("USD", sdk.NewInt(10000)),     // $10K
			TotalSavings:          sdk.NewCoin("USD", sdk.NewInt(6000000)),   // $6M saved
			AverageProcessingTime: 30 * time.Second,
			TraditionalCost:       sdk.NewCoin("USD", sdk.NewInt(7000000)),   // $7M traditional
			DeshChainCost:         sdk.NewCoin("USD", sdk.NewInt(1000000)),   // $1M DeshChain
			CostSavingPercent:     sdk.NewDecWithPrec(857, 3),                // 85.7%
		},
		"EUR-INR": {
			Corridor:              "EUR-INR",
			TotalVolume:           sdk.NewCoin("EUR", sdk.NewInt(50000000)),  // €50M
			TotalTransactions:     5000,
			AverageAmount:         sdk.NewCoin("EUR", sdk.NewInt(10000)),     // €10K
			TotalSavings:          sdk.NewCoin("EUR", sdk.NewInt(3000000)),   // €3M saved
			AverageProcessingTime: 30 * time.Second,
			TraditionalCost:       sdk.NewCoin("EUR", sdk.NewInt(3500000)),   // €3.5M traditional
			DeshChainCost:         sdk.NewCoin("EUR", sdk.NewInt(500000)),    // €0.5M DeshChain
			CostSavingPercent:     sdk.NewDecWithPrec(857, 3),                // 85.7%
		},
		"SGD-INR": {
			Corridor:              "SGD-INR",
			TotalVolume:           sdk.NewCoin("SGD", sdk.NewInt(25000000)),  // S$25M
			TotalTransactions:     2500,
			AverageAmount:         sdk.NewCoin("SGD", sdk.NewInt(10000)),     // S$10K
			TotalSavings:          sdk.NewCoin("SGD", sdk.NewInt(1500000)),   // S$1.5M saved
			AverageProcessingTime: 30 * time.Second,
			TraditionalCost:       sdk.NewCoin("SGD", sdk.NewInt(1750000)),   // S$1.75M traditional
			DeshChainCost:         sdk.NewCoin("SGD", sdk.NewInt(250000)),    // S$0.25M DeshChain
			CostSavingPercent:     sdk.NewDecWithPrec(857, 3),                // 85.7%
		},
	}
}