package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/tradefinance/types"
	oracletypes "github.com/DeshChain/DeshChain-Ecosystem/x/oracle/types"
)

// MultiCurrencyTradeManager handles multi-currency trade finance operations
type MultiCurrencyTradeManager struct {
	keeper *Keeper
}

// NewMultiCurrencyTradeManager creates a new multi-currency trade manager
func NewMultiCurrencyTradeManager(keeper *Keeper) *MultiCurrencyTradeManager {
	return &MultiCurrencyTradeManager{
		keeper: keeper,
	}
}

// Enhanced Letter of Credit with Multi-Currency Support
type EnhancedLetterOfCredit struct {
	// Existing LC fields
	types.LetterOfCredit
	
	// Multi-currency enhancement fields
	OriginalCurrency    string    `json:"original_currency"`    // "USD", "EUR", "SGD"
	OriginalAmount      sdk.Coin  `json:"original_amount"`      // Original trade amount
	SettlementCurrency  string    `json:"settlement_currency"`  // "DUSD", "DEUR", "DSGD"
	SettlementAmount    sdk.Coin  `json:"settlement_amount"`    // Stablecoin equivalent
	LocalCurrency       string    `json:"local_currency"`       // "DINR" for Indian recipients
	LocalAmount         sdk.Coin  `json:"local_amount"`         // Final local amount
	
	// Enhanced exchange rates and fees
	ExchangeRates       map[string]sdk.Dec `json:"exchange_rates"`
	ConversionFees      map[string]sdk.Dec `json:"conversion_fees"`
	TotalSavings        sdk.Coin           `json:"total_savings"` // vs traditional banking
	
	// Multi-currency metadata
	CurrencyRoute       []string  `json:"currency_route"`       // ["USD", "DUSD", "DINR"]
	ProcessingTime      time.Duration `json:"processing_time"`   // vs traditional methods
	TraditionalCost     sdk.Coin  `json:"traditional_cost"`     // Comparison cost
	DeshChainCost       sdk.Coin  `json:"deshchain_cost"`       // Actual cost
}

// ProcessMultiCurrencyLC processes an LC with multi-currency support
func (mctm *MultiCurrencyTradeManager) ProcessMultiCurrencyLC(
	ctx sdk.Context,
	lcID string,
	originalCurrency string,
	originalAmount sdk.Coin,
	settlementStablecoin string,
) (*EnhancedLetterOfCredit, error) {
	
	// 1. Validate currency support
	if !mctm.IsCurrencySupported(ctx, originalCurrency) {
		return nil, fmt.Errorf("unsupported currency: %s", originalCurrency)
	}
	
	if !mctm.IsStablecoinSupported(ctx, settlementStablecoin) {
		return nil, fmt.Errorf("unsupported settlement stablecoin: %s", settlementStablecoin)
	}
	
	// 2. Get real-time exchange rates from oracle
	exchangeRates := make(map[string]sdk.Dec)
	
	// Get original currency to stablecoin rate
	originalToStablecoinRate, err := mctm.getExchangeRate(ctx, originalCurrency, settlementStablecoin)
	if err != nil {
		return nil, fmt.Errorf("failed to get %s/%s exchange rate: %w", originalCurrency, settlementStablecoin, err)
	}
	exchangeRates[fmt.Sprintf("%s/%s", originalCurrency, settlementStablecoin)] = originalToStablecoinRate
	
	// Get stablecoin to DINR rate
	stablecoinToDINRRate, err := mctm.getExchangeRate(ctx, settlementStablecoin, "DINR")
	if err != nil {
		return nil, fmt.Errorf("failed to get %s/DINR exchange rate: %w", settlementStablecoin, err)
	}
	exchangeRates[fmt.Sprintf("%s/DINR", settlementStablecoin)] = stablecoinToDINRRate
	
	// 3. Calculate amounts and fees
	settlementAmount := sdk.NewCoin(settlementStablecoin, originalAmount.Amount.ToLegacyDec().Mul(originalToStablecoinRate).TruncateInt())
	localAmount := sdk.NewCoin("DINR", settlementAmount.Amount.ToLegacyDec().Mul(stablecoinToDINRRate).TruncateInt())
	
	// Calculate conversion fees (0.25% for multi-currency conversions)
	conversionFees := make(map[string]sdk.Dec)
	conversionFeeRate := sdk.NewDecWithPrec(25, 4) // 0.25%
	
	originalToStablecoinFee := originalAmount.Amount.ToLegacyDec().Mul(conversionFeeRate)
	conversionFees[fmt.Sprintf("%s/%s", originalCurrency, settlementStablecoin)] = originalToStablecoinFee
	
	stablecoinToDINRFee := settlementAmount.Amount.ToLegacyDec().Mul(conversionFeeRate)
	conversionFees[fmt.Sprintf("%s/DINR", settlementStablecoin)] = stablecoinToDINRFee
	
	// 4. Calculate traditional vs DeshChain costs
	traditionalCost := mctm.calculateTraditionalTradeCost(originalAmount)
	deshchainCost := sdk.NewCoin(originalAmount.Denom, originalToStablecoinFee.Add(stablecoinToDINRFee).TruncateInt())
	totalSavings := traditionalCost.Sub(deshchainCost)
	
	// 5. Lock stablecoin collateral
	if err := mctm.lockStablecoinCollateral(ctx, settlementStablecoin, settlementAmount); err != nil {
		return nil, fmt.Errorf("failed to lock stablecoin collateral: %w", err)
	}
	
	// 6. Get base LC (using existing LC processing)
	baseLC, found := mctm.keeper.GetLetterOfCredit(ctx, lcID)
	if !found {
		return nil, fmt.Errorf("letter of credit not found: %s", lcID)
	}
	
	// 7. Create enhanced LC
	enhancedLC := &EnhancedLetterOfCredit{
		LetterOfCredit:     baseLC,
		OriginalCurrency:   originalCurrency,
		OriginalAmount:     originalAmount,
		SettlementCurrency: settlementStablecoin,
		SettlementAmount:   settlementAmount,
		LocalCurrency:      "DINR",
		LocalAmount:        localAmount,
		ExchangeRates:      exchangeRates,
		ConversionFees:     conversionFees,
		TotalSavings:       totalSavings,
		CurrencyRoute:      []string{originalCurrency, settlementStablecoin, "DINR"},
		ProcessingTime:     5 * time.Minute, // vs 5-7 days traditional
		TraditionalCost:    traditionalCost,
		DeshChainCost:      deshchainCost,
	}
	
	// 8. Store enhanced LC data
	if err := mctm.storeEnhancedLC(ctx, enhancedLC); err != nil {
		return nil, fmt.Errorf("failed to store enhanced LC: %w", err)
	}
	
	// 9. Update LC status
	baseLC.Status = "multi_currency_processed"
	baseLC.UpdatedAt = ctx.BlockTime()
	mctm.keeper.SetLetterOfCredit(ctx, baseLC)
	
	// 10. Emit multi-currency LC event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"multi_currency_lc_processed",
			sdk.NewAttribute("lc_id", lcID),
			sdk.NewAttribute("original_currency", originalCurrency),
			sdk.NewAttribute("settlement_currency", settlementStablecoin),
			sdk.NewAttribute("original_amount", originalAmount.String()),
			sdk.NewAttribute("settlement_amount", settlementAmount.String()),
			sdk.NewAttribute("local_amount", localAmount.String()),
			sdk.NewAttribute("total_savings", totalSavings.String()),
			sdk.NewAttribute("processing_time", enhancedLC.ProcessingTime.String()),
		),
	)
	
	mctm.keeper.Logger(ctx).Info("processed multi-currency LC",
		"lc_id", lcID,
		"original_currency", originalCurrency,
		"settlement_currency", settlementStablecoin,
		"total_savings", totalSavings.String(),
	)
	
	return enhancedLC, nil
}

// getExchangeRate gets exchange rate between two currencies
func (mctm *MultiCurrencyTradeManager) getExchangeRate(ctx sdk.Context, fromCurrency, toCurrency string) (sdk.Dec, error) {
	// This would integrate with the oracle keeper to get real-time rates
	// For now, we'll use placeholder logic
	
	// Handle DUSD rates (using USD as proxy)
	if fromCurrency == "USD" && toCurrency == "DUSD" {
		return sdk.OneDec(), nil // 1:1 peg
	}
	if fromCurrency == "DUSD" && toCurrency == "USD" {
		return sdk.OneDec(), nil // 1:1 peg
	}
	
	// Handle USD to INR/DINR
	if fromCurrency == "USD" && (toCurrency == "INR" || toCurrency == "DINR") {
		// Use oracle to get USD/INR rate (example: 83.5)
		return sdk.NewDecWithPrec(835, 1), nil // 83.5 INR per USD
	}
	if fromCurrency == "DUSD" && toCurrency == "DINR" {
		// DUSD pegged to USD, so same rate as USD/INR
		return sdk.NewDecWithPrec(835, 1), nil // 83.5 DINR per DUSD
	}
	
	// Handle EUR rates
	if fromCurrency == "EUR" && toCurrency == "DEUR" {
		return sdk.OneDec(), nil // 1:1 peg
	}
	if fromCurrency == "EUR" && (toCurrency == "INR" || toCurrency == "DINR") {
		// EUR/INR rate (example: 91.2)
		return sdk.NewDecWithPrec(912, 1), nil // 91.2 INR per EUR
	}
	
	// Handle SGD rates
	if fromCurrency == "SGD" && toCurrency == "DSGD" {
		return sdk.OneDec(), nil // 1:1 peg
	}
	if fromCurrency == "SGD" && (toCurrency == "INR" || toCurrency == "DINR") {
		// SGD/INR rate (example: 62.1)
		return sdk.NewDecWithPrec(621, 1), nil // 62.1 INR per SGD
	}
	
	// Cross rates via USD
	return mctm.calculateCrossRate(ctx, fromCurrency, toCurrency)
}

// calculateCrossRate calculates cross rates via USD
func (mctm *MultiCurrencyTradeManager) calculateCrossRate(ctx sdk.Context, fromCurrency, toCurrency string) (sdk.Dec, error) {
	// Get from currency to USD rate
	fromToUSDRate, err := mctm.getDirectRate(ctx, fromCurrency, "USD")
	if err != nil {
		return sdk.ZeroDec(), err
	}
	
	// Get USD to target currency rate
	usdToTargetRate, err := mctm.getDirectRate(ctx, "USD", toCurrency)
	if err != nil {
		return sdk.ZeroDec(), err
	}
	
	// Calculate cross rate
	crossRate := fromToUSDRate.Mul(usdToTargetRate)
	return crossRate, nil
}

// getDirectRate gets direct exchange rate (placeholder for oracle integration)
func (mctm *MultiCurrencyTradeManager) getDirectRate(ctx sdk.Context, fromCurrency, toCurrency string) (sdk.Dec, error) {
	// This would integrate with the enhanced oracle system
	// For now, return sample rates
	
	rates := map[string]sdk.Dec{
		"EUR/USD": sdk.NewDecWithPrec(109, 2),  // 1.09
		"SGD/USD": sdk.NewDecWithPrec(74, 2),   // 0.74
		"GBP/USD": sdk.NewDecWithPrec(127, 2),  // 1.27
		"USD/INR": sdk.NewDecWithPrec(835, 1),  // 83.5
	}
	
	rateKey := fmt.Sprintf("%s/%s", fromCurrency, toCurrency)
	if rate, found := rates[rateKey]; found {
		return rate, nil
	}
	
	// Try inverse rate
	inverseKey := fmt.Sprintf("%s/%s", toCurrency, fromCurrency)
	if rate, found := rates[inverseKey]; found {
		return sdk.OneDec().Quo(rate), nil
	}
	
	return sdk.ZeroDec(), fmt.Errorf("exchange rate not available for %s/%s", fromCurrency, toCurrency)
}

// calculateTraditionalTradeCost calculates traditional banking cost for comparison
func (mctm *MultiCurrencyTradeManager) calculateTraditionalTradeCost(amount sdk.Coin) sdk.Coin {
	// Traditional LC costs: 2-4% of trade value
	traditionalFeeRate := sdk.NewDecWithPrec(3, 2) // 3%
	traditionalFee := amount.Amount.ToLegacyDec().Mul(traditionalFeeRate)
	return sdk.NewCoin(amount.Denom, traditionalFee.TruncateInt())
}

// lockStablecoinCollateral locks stablecoin collateral for LC
func (mctm *MultiCurrencyTradeManager) lockStablecoinCollateral(ctx sdk.Context, stablecoin string, amount sdk.Coin) error {
	// Lock collateral in escrow for LC guarantee
	// This would integrate with the stablecoin modules (DUSD, DEUR, etc.)
	
	mctm.keeper.Logger(ctx).Info("locking stablecoin collateral",
		"stablecoin", stablecoin,
		"amount", amount.String(),
	)
	
	// For now, just log the operation
	// In full implementation, this would:
	// 1. Transfer stablecoin to escrow account
	// 2. Lock the funds until LC completion
	// 3. Create collateral record for tracking
	
	return nil
}

// storeEnhancedLC stores enhanced LC data
func (mctm *MultiCurrencyTradeManager) storeEnhancedLC(ctx sdk.Context, enhancedLC *EnhancedLetterOfCredit) error {
	store := ctx.KVStore(mctm.keeper.storeKey)
	key := types.GetEnhancedLCKey(enhancedLC.LcId)
	
	bz, err := mctm.keeper.cdc.Marshal(enhancedLC)
	if err != nil {
		return err
	}
	
	store.Set(key, bz)
	return nil
}

// GetEnhancedLC retrieves enhanced LC data
func (mctm *MultiCurrencyTradeManager) GetEnhancedLC(ctx sdk.Context, lcID string) (*EnhancedLetterOfCredit, bool) {
	store := ctx.KVStore(mctm.keeper.storeKey)
	key := types.GetEnhancedLCKey(lcID)
	
	bz := store.Get(key)
	if bz == nil {
		return nil, false
	}
	
	var enhancedLC EnhancedLetterOfCredit
	if err := mctm.keeper.cdc.Unmarshal(bz, &enhancedLC); err != nil {
		return nil, false
	}
	
	return &enhancedLC, true
}

// IsCurrencySupported checks if a currency is supported for trade finance
func (mctm *MultiCurrencyTradeManager) IsCurrencySupported(ctx sdk.Context, currency string) bool {
	supportedCurrencies := []string{"USD", "EUR", "SGD", "GBP", "JPY", "INR", "DINR"}
	
	for _, supported := range supportedCurrencies {
		if currency == supported {
			return true
		}
	}
	
	return false
}

// IsStablecoinSupported checks if a stablecoin is supported for settlement
func (mctm *MultiCurrencyTradeManager) IsStablecoinSupported(ctx sdk.Context, stablecoin string) bool {
	supportedStablecoins := []string{"DUSD", "DEUR", "DSGD", "DGBP", "DJPY", "DINR"}
	
	for _, supported := range supportedStablecoins {
		if stablecoin == supported {
			return true
		}
	}
	
	return false
}

// GetMultiCurrencyStats returns statistics for multi-currency trade finance
func (mctm *MultiCurrencyTradeManager) GetMultiCurrencyStats(ctx sdk.Context) types.MultiCurrencyStats {
	// This would aggregate statistics from all enhanced LCs
	// For now, return placeholder data
	
	return types.MultiCurrencyStats{
		TotalLCs:           100,
		TotalVolume:        sdk.NewCoin("USD", sdk.NewInt(10000000)), // $10M
		TotalSavings:       sdk.NewCoin("USD", sdk.NewInt(300000)),   // $300K saved
		AverageProcessingTime: 5 * time.Minute,
		CurrencyBreakdown: map[string]int64{
			"USD": 60,
			"EUR": 25,
			"SGD": 10,
			"GBP": 5,
		},
		StablecoinBreakdown: map[string]int64{
			"DUSD": 60,
			"DEUR": 25,
			"DSGD": 10,
			"DGBP": 5,
		},
	}
}