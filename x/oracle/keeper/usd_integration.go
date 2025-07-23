package keeper

import (
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/deshchain/x/oracle/types"
)

// USDOracleManager manages USD-specific oracle operations for DUSD integration
type USDOracleManager struct {
	keeper *Keeper
}

// NewUSDOracleManager creates a new USD oracle manager
func NewUSDOracleManager(keeper *Keeper) *USDOracleManager {
	return &USDOracleManager{
		keeper: keeper,
	}
}

// USD Oracle Sources Configuration
var USDOracleSources = map[string]types.OracleSourceConfig{
	"chainlink": {
		Name:        "Chainlink",
		Priority:    1,
		Weight:      30,
		Endpoint:    "https://api.chain.link/v1/feeds/usd-inr",
		UpdateFreq:  60, // 1 minute
		Reliability: 99.9,
	},
	"federal_reserve": {
		Name:        "Federal Reserve",
		Priority:    1,
		Weight:      25,
		Endpoint:    "https://api.stlouisfed.org/fred/series/observations?series_id=DEXINUS",
		UpdateFreq:  300, // 5 minutes
		Reliability: 99.5,
	},
	"band_protocol": {
		Name:        "Band Protocol", 
		Priority:    2,
		Weight:      20,
		Endpoint:    "https://api.bandchain.org/v1/oracle/symbols/USD",
		UpdateFreq:  30, // 30 seconds
		Reliability: 99.8,
	},
	"pyth_network": {
		Name:        "Pyth Network",
		Priority:    2,
		Weight:      15,
		Endpoint:    "https://pyth.network/api/price_feeds/crypto.USD/INR",
		UpdateFreq:  15, // 15 seconds
		Reliability: 99.7,
	},
	"bloomberg": {
		Name:        "Bloomberg Terminal",
		Priority:    3,
		Weight:      10,
		Endpoint:    "https://api.bloomberg.com/v1/currencies/USDINR",
		UpdateFreq:  60, // 1 minute
		Reliability: 99.9,
	},
}

// GetUSDPrice returns current USD price with enhanced validation
func (uom *USDOracleManager) GetUSDPrice(ctx sdk.Context) (types.PriceData, error) {
	// Try to get existing USD price data
	priceData, found := uom.keeper.GetPriceData(ctx, "USD")
	if !found {
		return types.PriceData{}, types.ErrPriceDataNotFound
	}
	
	// Check price freshness (should be within 5 minutes for USD)
	params := uom.keeper.GetParams(ctx)
	freshnessDuration := time.Duration(params.StalenessThreshold) * time.Second
	if ctx.BlockTime().Sub(priceData.LastUpdated) > freshnessDuration {
		return types.PriceData{}, types.ErrStalePriceData
	}
	
	// Validate price reasonableness (USD/INR should be between 70-100 typically)
	if priceData.Price.LT(sdk.NewDec(50)) || priceData.Price.GT(sdk.NewDec(120)) {
		return types.PriceData{}, types.ErrUnreasonablePrice
	}
	
	return priceData, nil
}

// SubmitUSDPrice submits USD price from multiple oracle sources
func (uom *USDOracleManager) SubmitUSDPrice(ctx sdk.Context, validator string, price sdk.Dec, source string) error {
	// Validate USD-specific price bounds
	if err := uom.validateUSDPrice(price); err != nil {
		return err
	}
	
	// Check if source is supported
	sourceConfig, found := USDOracleSources[strings.ToLower(source)]
	if !found {
		return fmt.Errorf("unsupported USD oracle source: %s", source)
	}
	
	// Submit price using base oracle functionality with USD-specific validation
	return uom.keeper.SubmitPrice(ctx, validator, "USD", price, sourceConfig.Name, ctx.BlockTime())
}

// validateUSDPrice validates USD price submissions
func (uom *USDOracleManager) validateUSDPrice(price sdk.Dec) error {
	// USD/INR reasonable bounds
	minUSDPrice := sdk.NewDec(50)  // Minimum reasonable USD/INR rate
	maxUSDPrice := sdk.NewDec(120) // Maximum reasonable USD/INR rate
	
	if price.LT(minUSDPrice) {
		return fmt.Errorf("USD price too low: %s < %s", price.String(), minUSDPrice.String())
	}
	
	if price.GT(maxUSDPrice) {
		return fmt.Errorf("USD price too high: %s > %s", price.String(), maxUSDPrice.String())
	}
	
	return nil
}

// GetUSDExchangeRate returns USD to target currency exchange rate
func (uom *USDOracleManager) GetUSDExchangeRate(ctx sdk.Context, targetCurrency string) (sdk.Dec, error) {
	// Get USD price in INR
	usdPrice, err := uom.GetUSDPrice(ctx)
	if err != nil {
		return sdk.ZeroDec(), err
	}
	
	if targetCurrency == "INR" || targetCurrency == "DINR" {
		return usdPrice.Price, nil
	}
	
	// For other currencies, try to get direct exchange rate first
	exchangeRate, found := uom.keeper.GetExchangeRate(ctx, "USD", targetCurrency)
	if found {
		// Check exchange rate freshness
		params := uom.keeper.GetParams(ctx)
		freshnessDuration := time.Duration(params.StalenessThreshold) * time.Second
		if ctx.BlockTime().Sub(exchangeRate.LastUpdated) <= freshnessDuration {
			return exchangeRate.Rate, nil
		}
	}
	
	// Try cross-rate calculation: USD->INR->TargetCurrency
	targetPrice, found := uom.keeper.GetPriceData(ctx, targetCurrency)
	if !found {
		return sdk.ZeroDec(), fmt.Errorf("exchange rate not available for USD/%s", targetCurrency)
	}
	
	// Calculate cross rate: USD/INR / TargetCurrency/INR = USD/TargetCurrency
	if targetPrice.Price.IsZero() {
		return sdk.ZeroDec(), fmt.Errorf("invalid target currency price for %s", targetCurrency)
	}
	
	crossRate := usdPrice.Price.Quo(targetPrice.Price)
	return crossRate, nil
}

// AggregateUSDPrices aggregates USD prices from multiple sources with enhanced logic
func (uom *USDOracleManager) AggregateUSDPrices(ctx sdk.Context, submissions []types.ValidatorPriceSubmission) (sdk.Dec, error) {
	if len(submissions) == 0 {
		return sdk.ZeroDec(), fmt.Errorf("no USD price submissions found")
	}
	
	// Weight prices by oracle source reliability
	weightedSum := sdk.ZeroDec()
	totalWeight := sdk.ZeroDec()
	
	for _, submission := range submissions {
		// Get source configuration
		sourceKey := strings.ToLower(submission.Source)
		sourceConfig, found := USDOracleSources[sourceKey]
		if !found {
			// Default weight for unknown sources
			sourceConfig = types.OracleSourceConfig{Weight: 5}
		}
		
		// Validate price bounds
		if err := uom.validateUSDPrice(submission.Price); err != nil {
			uom.keeper.Logger(ctx).Warn("invalid USD price submission", 
				"validator", submission.Validator,
				"price", submission.Price.String(),
				"error", err,
			)
			continue
		}
		
		// Apply weight
		weight := sdk.NewDec(int64(sourceConfig.Weight))
		weightedSum = weightedSum.Add(submission.Price.Mul(weight))
		totalWeight = totalWeight.Add(weight)
	}
	
	if totalWeight.IsZero() {
		return sdk.ZeroDec(), fmt.Errorf("no valid USD price submissions after filtering")
	}
	
	aggregatedPrice := weightedSum.Quo(totalWeight)
	
	// Final validation of aggregated price
	if err := uom.validateUSDPrice(aggregatedPrice); err != nil {
		return sdk.ZeroDec(), fmt.Errorf("aggregated USD price failed validation: %w", err)
	}
	
	return aggregatedPrice, nil
}

// MonitorUSDPriceStability monitors USD price for unusual volatility
func (uom *USDOracleManager) MonitorUSDPriceStability(ctx sdk.Context) error {
	// Get current USD price
	currentPrice, err := uom.GetUSDPrice(ctx)
	if err != nil {
		return err
	}
	
	// Get historical prices for comparison
	historicalPrices := uom.keeper.GetPriceHistory(ctx, "USD", 10) // Last 10 prices
	if len(historicalPrices) < 2 {
		return nil // Not enough history
	}
	
	// Calculate price volatility
	var priceChanges []sdk.Dec
	for i := 1; i < len(historicalPrices); i++ {
		prevPrice := historicalPrices[i].Price
		currentPriceForCalc := historicalPrices[i-1].Price
		
		if !prevPrice.IsZero() {
			change := currentPriceForCalc.Sub(prevPrice).Quo(prevPrice).Abs()
			priceChanges = append(priceChanges, change)
		}
	}
	
	if len(priceChanges) == 0 {
		return nil
	}
	
	// Calculate average volatility
	totalChange := sdk.ZeroDec()
	for _, change := range priceChanges {
		totalChange = totalChange.Add(change)
	}
	avgVolatility := totalChange.Quo(sdk.NewDec(int64(len(priceChanges))))
	
	// Check if current price change exceeds normal volatility significantly
	if len(historicalPrices) > 0 {
		lastPrice := historicalPrices[0].Price
		currentChange := currentPrice.Price.Sub(lastPrice).Quo(lastPrice).Abs()
		
		// Alert if change is more than 3x average volatility and > 2%
		alertThreshold := avgVolatility.Mul(sdk.NewDec(3))
		if currentChange.GT(alertThreshold) && currentChange.GT(sdk.NewDecWithPrec(2, 2)) {
			// Emit volatility alert
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					"usd_price_volatility_alert",
					sdk.NewAttribute("current_price", currentPrice.Price.String()),
					sdk.NewAttribute("last_price", lastPrice.String()),
					sdk.NewAttribute("price_change", currentChange.String()),
					sdk.NewAttribute("avg_volatility", avgVolatility.String()),
					sdk.NewAttribute("alert_threshold", alertThreshold.String()),
				),
			)
			
			uom.keeper.Logger(ctx).Warn("USD price volatility alert",
				"current_price", currentPrice.Price.String(),
				"price_change", currentChange.String(),
				"avg_volatility", avgVolatility.String(),
			)
		}
	}
	
	return nil
}

// GetCrossRate calculates cross rates between currencies via USD
func (uom *USDOracleManager) GetCrossRate(ctx sdk.Context, baseCurrency, quoteCurrency string) (sdk.Dec, error) {
	// If one of the currencies is USD, use direct rates
	if baseCurrency == "USD" {
		return uom.GetUSDExchangeRate(ctx, quoteCurrency)
	}
	if quoteCurrency == "USD" {
		rate, err := uom.GetUSDExchangeRate(ctx, baseCurrency)
		if err != nil {
			return sdk.ZeroDec(), err
		}
		return sdk.OneDec().Quo(rate), nil // Inverse rate
	}
	
	// Calculate cross rate via USD: Base/USD * USD/Quote = Base/Quote
	baseToUSD, err := uom.GetUSDExchangeRate(ctx, baseCurrency)
	if err != nil {
		return sdk.ZeroDec(), err
	}
	
	usdToQuote, err := uom.GetUSDExchangeRate(ctx, quoteCurrency)
	if err != nil {
		return sdk.ZeroDec(), err
	}
	
	if usdToQuote.IsZero() {
		return sdk.ZeroDec(), fmt.Errorf("zero exchange rate for USD/%s", quoteCurrency)
	}
	
	// Cross rate = (Base/USD) / (Quote/USD) = Base/Quote
	crossRate := baseToUSD.Quo(usdToQuote)
	return crossRate, nil
}

// ValidateUSDOracleHealth checks the health of USD oracle network
func (uom *USDOracleManager) ValidateUSDOracleHealth(ctx sdk.Context) error {
	// Get all oracle validators
	validators := uom.keeper.GetAllOracleValidators(ctx)
	
	activeValidators := 0
	recentSubmissions := 0
	
	// Check validator activity in last 10 blocks
	recentThreshold := ctx.BlockTime().Add(-10 * time.Minute)
	
	for _, validator := range validators {
		if validator.Active {
			activeValidators++
			
			// Check if validator has submitted recently
			if validator.LastSubmission.After(recentThreshold) {
				recentSubmissions++
			}
		}
	}
	
	// Ensure minimum validator participation
	params := uom.keeper.GetParams(ctx)
	if uint64(activeValidators) < params.MinValidators {
		return fmt.Errorf("insufficient active USD oracle validators: %d < %d", 
			activeValidators, params.MinValidators)
	}
	
	// Ensure minimum recent activity (at least 50% of active validators)
	minRecentSubmissions := activeValidators / 2
	if recentSubmissions < minRecentSubmissions {
		uom.keeper.Logger(ctx).Warn("low USD oracle activity",
			"active_validators", activeValidators,
			"recent_submissions", recentSubmissions,
			"threshold", minRecentSubmissions,
		)
		
		// Emit health warning but don't fail
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				"usd_oracle_health_warning",
				sdk.NewAttribute("active_validators", sdk.NewInt(int64(activeValidators)).String()),
				sdk.NewAttribute("recent_submissions", sdk.NewInt(int64(recentSubmissions)).String()),
				sdk.NewAttribute("min_required", sdk.NewInt(int64(minRecentSubmissions)).String()),
			),
		)
	}
	
	return nil
}

// GetPrice is a convenience method that integrates with existing oracle keeper interface
func (k Keeper) GetPrice(ctx sdk.Context, base, quote string) (types.PriceData, error) {
	// Handle USD-specific logic
	if base == "USD" || quote == "USD" {
		usdManager := NewUSDOracleManager(&k)
		
		if base == "USD" && quote == "DINR" {
			// Direct USD price
			return usdManager.GetUSDPrice(ctx)
		}
		
		if base == "USD" {
			// USD to other currency
			rate, err := usdManager.GetUSDExchangeRate(ctx, quote)
			if err != nil {
				return types.PriceData{}, err
			}
			
			return types.PriceData{
				Symbol:      fmt.Sprintf("%s/%s", base, quote),
				Price:       rate,
				LastUpdated: ctx.BlockTime(),
				Source:      "usd_oracle_manager",
			}, nil
		}
		
		if quote == "USD" {
			// Other currency to USD
			rate, err := usdManager.GetUSDExchangeRate(ctx, base)
			if err != nil {
				return types.PriceData{}, err
			}
			
			// Return inverse rate
			if rate.IsZero() {
				return types.PriceData{}, fmt.Errorf("zero rate for %s/USD", base)
			}
			
			return types.PriceData{
				Symbol:      fmt.Sprintf("%s/%s", base, quote),
				Price:       sdk.OneDec().Quo(rate),
				LastUpdated: ctx.BlockTime(),
				Source:      "usd_oracle_manager",
			}, nil
		}
	}
	
	// For non-USD pairs, use existing logic
	return k.GetPriceData(ctx, fmt.Sprintf("%s/%s", base, quote))
}