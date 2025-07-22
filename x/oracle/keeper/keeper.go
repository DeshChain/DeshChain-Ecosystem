package keeper

import (
	"fmt"
	"sort"
	"time"

	"github.com/cometbft/cometbft/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/deshchain/deshchain/x/oracle/types"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		memKey     storetypes.StoreKey
		paramstore paramtypes.Subspace

		accountKeeper types.AccountKeeper
		bankKeeper    types.BankKeeper
		stakingKeeper types.StakingKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		memKey:        memKey,
		paramstore:    ps,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		stakingKeeper: stakingKeeper,
	}
}

// StoreKey returns the store key
func (k Keeper) StoreKey() storetypes.StoreKey {
	return k.storeKey
}

// Codec returns the codec
func (k Keeper) Codec() codec.BinaryCodec {
	return k.cdc
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetParams get all parameters as types.OracleParams
func (k Keeper) GetParams(ctx sdk.Context) types.OracleParams {
	var p types.OracleParams
	k.paramstore.GetParamSet(ctx, &p)
	return p
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.OracleParams) {
	k.paramstore.SetParamSet(ctx, &params)
}

// GetAuthority returns the module's authority
func (k Keeper) GetAuthority() string {
	// In Cosmos SDK, this is typically the governance module account
	return k.accountKeeper.GetModuleAddress("gov").String()
}

// SubmitExchangeRate handles exchange rate submission from oracle validators
func (k Keeper) SubmitExchangeRate(ctx sdk.Context, validator, base, target string, rate sdk.Dec, source string, timestamp time.Time) error {
	// Check if validator is authorized
	oracleValidator, found := k.GetOracleValidator(ctx, validator)
	if !found {
		return types.ErrUnauthorizedValidator
	}

	if !oracleValidator.Active {
		return types.ErrValidatorNotActive
	}

	// Validate exchange rate submission
	if err := k.validateExchangeRateSubmission(ctx, base, target, rate, timestamp); err != nil {
		return err
	}

	// Store the exchange rate (simplified - could be enhanced with aggregation)
	exchangeRate := types.ExchangeRate{
		Base:        base,
		Target:      target,
		Rate:        rate,
		LastUpdated: timestamp,
		Source:      source,
		Validator:   validator,
	}

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&exchangeRate)
	store.Set(types.ExchangeRateKey(base, target), bz)

	// Update validator stats
	oracleValidator.SuccessfulSubmissions++
	oracleValidator.LastSubmission = timestamp
	k.SetOracleValidator(ctx, oracleValidator)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeExchangeRateUpdate,
			sdk.NewAttribute(types.AttributeKeyValidator, validator),
			sdk.NewAttribute(types.AttributeKeyBase, base),
			sdk.NewAttribute(types.AttributeKeyTarget, target),
			sdk.NewAttribute(types.AttributeKeyRate, rate.String()),
			sdk.NewAttribute(types.AttributeKeySource, source),
			sdk.NewAttribute(types.AttributeKeyTimestamp, timestamp.Format(time.RFC3339)),
		),
	)

	return nil
}

// validateExchangeRateSubmission validates an exchange rate submission
func (k Keeper) validateExchangeRateSubmission(ctx sdk.Context, base, target string, rate sdk.Dec, timestamp time.Time) error {
	// Check if currencies are valid (basic validation)
	if len(base) < 2 || len(base) > 10 {
		return types.ErrInvalidCurrency
	}
	
	if len(target) < 2 || len(target) > 10 {
		return types.ErrInvalidCurrency
	}

	if base == target {
		return types.ErrInvalidCurrency
	}

	// Check rate bounds
	if rate.GT(sdk.NewDec(1000000000)) { // 1 billion max
		return types.ErrInvalidExchangeRate
	}

	// Check timestamp freshness
	blockTime := ctx.BlockTime()
	if timestamp.After(blockTime.Add(5 * time.Minute)) {
		return types.ErrInvalidTimestamp
	}

	if timestamp.Before(blockTime.Add(-1 * time.Hour)) {
		return types.ErrInvalidTimestamp
	}

	return nil
}

// SubmitPrice handles price submission from oracle validators
func (k Keeper) SubmitPrice(ctx sdk.Context, validator, symbol string, price sdk.Dec, source string, timestamp time.Time) error {
	// Check if validator is authorized
	oracleValidator, found := k.GetOracleValidator(ctx, validator)
	if !found {
		return types.ErrUnauthorizedValidator
	}

	if !oracleValidator.Active {
		return types.ErrValidatorNotActive
	}

	// Validate price submission
	if err := k.validatePriceSubmission(ctx, symbol, price, timestamp); err != nil {
		return err
	}

	// Check for duplicate submission in current window
	currentHeight := uint64(ctx.BlockHeight())
	params := k.GetParams(ctx)
	windowStart := currentHeight - (currentHeight % params.AggregationWindow)

	submissionKey := types.ValidatorSubmissionKey(validator, symbol, currentHeight)
	store := ctx.KVStore(k.storeKey)

	if store.Has(submissionKey) {
		return types.ErrDuplicateSubmission
	}

	// Store the submission
	submission := types.ValidatorPriceSubmission{
		Validator:   validator,
		Symbol:      symbol,
		Price:       price,
		BlockHeight: currentHeight,
		Timestamp:   timestamp,
		Source:      source,
	}

	bz := k.cdc.MustMarshal(&submission)
	store.Set(submissionKey, bz)

	// Update validator stats
	oracleValidator.SuccessfulSubmissions++
	oracleValidator.LastSubmission = timestamp
	k.SetOracleValidator(ctx, oracleValidator)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypePriceSubmission,
			sdk.NewAttribute(types.AttributeKeyValidator, validator),
			sdk.NewAttribute(types.AttributeKeySymbol, symbol),
			sdk.NewAttribute(types.AttributeKeyPrice, price.String()),
			sdk.NewAttribute(types.AttributeKeySource, source),
			sdk.NewAttribute(types.AttributeKeyTimestamp, timestamp.Format(time.RFC3339)),
		),
	)

	return nil
}

// GetPriceData retrieves price data for a symbol
func (k Keeper) GetPriceData(ctx sdk.Context, symbol string) (types.PriceData, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.PriceDataKey(symbol))
	if bz == nil {
		return types.PriceData{}, false
	}

	var priceData types.PriceData
	k.cdc.MustUnmarshal(bz, &priceData)
	return priceData, true
}

// SetPriceData stores price data for a symbol
func (k Keeper) SetPriceData(ctx sdk.Context, priceData types.PriceData) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&priceData)
	store.Set(types.PriceDataKey(priceData.Symbol), bz)
}

// GetAllPriceData retrieves all price data
func (k Keeper) GetAllPriceData(ctx sdk.Context) []types.PriceData {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.PriceDataKeyPrefix)
	defer iterator.Close()

	var priceDataList []types.PriceData
	for ; iterator.Valid(); iterator.Next() {
		var priceData types.PriceData
		k.cdc.MustUnmarshal(iterator.Value(), &priceData)
		priceDataList = append(priceDataList, priceData)
	}

	return priceDataList
}

// GetOracleValidator retrieves an oracle validator
func (k Keeper) GetOracleValidator(ctx sdk.Context, validator string) (types.OracleValidator, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.OracleValidatorKey(validator))
	if bz == nil {
		return types.OracleValidator{}, false
	}

	var oracleValidator types.OracleValidator
	k.cdc.MustUnmarshal(bz, &oracleValidator)
	return oracleValidator, true
}

// SetOracleValidator stores an oracle validator
func (k Keeper) SetOracleValidator(ctx sdk.Context, oracleValidator types.OracleValidator) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&oracleValidator)
	store.Set(types.OracleValidatorKey(oracleValidator.Validator), bz)
}

// GetAllOracleValidators retrieves all oracle validators
func (k Keeper) GetAllOracleValidators(ctx sdk.Context) []types.OracleValidator {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.OracleValidatorKeyPrefix)
	defer iterator.Close()

	var validators []types.OracleValidator
	for ; iterator.Valid(); iterator.Next() {
		var validator types.OracleValidator
		k.cdc.MustUnmarshal(iterator.Value(), &validator)
		validators = append(validators, validator)
	}

	return validators
}

// RegisterOracleValidator registers a new oracle validator
func (k Keeper) RegisterOracleValidator(ctx sdk.Context, validatorAddr string, power uint64, description string) error {
	// Check if validator already exists
	if _, found := k.GetOracleValidator(ctx, validatorAddr); found {
		return types.ErrValidatorAlreadyExists
	}

	// Verify the validator exists in staking module
	valAddr, err := sdk.ValAddressFromBech32(validatorAddr)
	if err != nil {
		return types.ErrInvalidValidator
	}

	validator, found := k.stakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		return types.ErrValidatorNotFound
	}

	// Create oracle validator
	oracleValidator := types.OracleValidator{
		Validator:             validatorAddr,
		Power:                 power,
		Active:                true,
		SuccessfulSubmissions: 0,
		FailedSubmissions:     0,
		LastSubmission:        time.Time{},
		SlashCount:            0,
	}

	k.SetOracleValidator(ctx, oracleValidator)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeValidatorRegistered,
			sdk.NewAttribute(types.AttributeKeyValidator, validatorAddr),
			sdk.NewAttribute(types.AttributeKeyPower, fmt.Sprintf("%d", power)),
			sdk.NewAttribute(types.AttributeKeyActive, "true"),
		),
	)

	return nil
}

// ProcessAggregationWindow processes price submissions for the current aggregation window
func (k Keeper) ProcessAggregationWindow(ctx sdk.Context) error {
	params := k.GetParams(ctx)
	currentHeight := uint64(ctx.BlockHeight())
	
	// Only process at the end of aggregation windows
	if currentHeight%params.AggregationWindow != 0 {
		return nil
	}

	windowStart := currentHeight - params.AggregationWindow
	windowEnd := currentHeight

	// Get all tracked symbols
	allPriceData := k.GetAllPriceData(ctx)
	symbols := make(map[string]bool)
	for _, priceData := range allPriceData {
		symbols[priceData.Symbol] = true
	}

	// Get symbols from recent submissions if no existing price data
	if len(symbols) == 0 {
		symbols = k.getSymbolsFromRecentSubmissions(ctx, windowStart, windowEnd)
	}

	// Process each symbol
	for symbol := range symbols {
		if err := k.aggregatePricesForSymbol(ctx, symbol, windowStart, windowEnd); err != nil {
			k.Logger(ctx).Error("failed to aggregate prices for symbol", "symbol", symbol, "error", err)
			continue
		}
	}

	return nil
}

// validatePriceSubmission validates a price submission
func (k Keeper) validatePriceSubmission(ctx sdk.Context, symbol string, price sdk.Dec, timestamp time.Time) error {
	// Check if symbol is supported (basic validation)
	if len(symbol) < 2 || len(symbol) > 10 {
		return types.ErrInvalidSymbol
	}

	// Check price bounds (reasonable price range)
	if price.GT(sdk.NewDec(1000000000)) { // 1 billion INR max
		return types.ErrInvalidPrice
	}

	// Check timestamp freshness
	blockTime := ctx.BlockTime()
	if timestamp.After(blockTime.Add(5 * time.Minute)) {
		return types.ErrInvalidTimestamp
	}

	if timestamp.Before(blockTime.Add(-1 * time.Hour)) {
		return types.ErrInvalidTimestamp
	}

	return nil
}

// aggregatePricesForSymbol aggregates prices for a specific symbol in the given window
func (k Keeper) aggregatePricesForSymbol(ctx sdk.Context, symbol string, windowStart, windowEnd uint64) error {
	submissions := k.getSubmissionsForWindow(ctx, symbol, windowStart, windowEnd)
	params := k.GetParams(ctx)

	// Check if we have enough validators
	if uint64(len(submissions)) < params.MinValidators {
		return types.ErrInsufficientValidators
	}

	// Extract prices and calculate statistics
	prices := make([]sdk.Dec, len(submissions))
	validators := make([]string, len(submissions))
	totalPower := uint64(0)

	for i, submission := range submissions {
		prices[i] = submission.Price
		validators[i] = submission.Validator

		// Get validator power for weighted calculations
		if oracleVal, found := k.GetOracleValidator(ctx, submission.Validator); found {
			totalPower += oracleVal.Power
		}
	}

	// Sort prices for median calculation
	sort.Slice(prices, func(i, j int) bool {
		return prices[i].LT(prices[j])
	})

	// Calculate median
	var median sdk.Dec
	if len(prices)%2 == 0 {
		median = prices[len(prices)/2-1].Add(prices[len(prices)/2]).Quo(sdk.NewDec(2))
	} else {
		median = prices[len(prices)/2]
	}

	// Calculate mean
	sum := sdk.ZeroDec()
	for _, price := range prices {
		sum = sum.Add(price)
	}
	mean := sum.Quo(sdk.NewDec(int64(len(prices))))

	// Calculate standard deviation
	variance := sdk.ZeroDec()
	for _, price := range prices {
		diff := price.Sub(mean)
		variance = variance.Add(diff.Mul(diff))
	}
	variance = variance.Quo(sdk.NewDec(int64(len(prices))))
	stdDev := sdk.NewDecFromInt(variance.TruncateInt()).ApproxRoot(2) // Approximate square root

	// Check deviation against parameters
	maxDeviation := mean.Mul(params.MaxPriceDeviation)
	if stdDev.GT(maxDeviation) {
		// Emit deviation event but don't fail
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypePriceDeviation,
				sdk.NewAttribute(types.AttributeKeySymbol, symbol),
				sdk.NewAttribute(types.AttributeKeyMeanPrice, mean.String()),
				sdk.NewAttribute(types.AttributeKeyDeviation, stdDev.String()),
				sdk.NewAttribute(types.AttributeKeyMaxDeviation, maxDeviation.String()),
				sdk.NewAttribute(types.AttributeKeyDeviationExceeded, "true"),
			),
		)
	}

	// Update price data
	currentPrice, found := k.GetPriceData(ctx, symbol)
	oldPrice := sdk.ZeroDec()
	if found {
		oldPrice = currentPrice.Price
	}

	newPriceData := types.PriceData{
		Symbol:         symbol,
		Price:          median, // Use median as the final price
		LastUpdated:    ctx.BlockTime(),
		Source:         "aggregated",
		ValidatorCount: uint64(len(submissions)),
		Deviation:      stdDev,
	}

	k.SetPriceData(ctx, newPriceData)

	// Store historical price
	k.storeHistoricalPrice(ctx, symbol, median, uint64(ctx.BlockHeight()))

	// Emit price update event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypePriceUpdate,
			sdk.NewAttribute(types.AttributeKeySymbol, symbol),
			sdk.NewAttribute(types.AttributeKeyOldPrice, oldPrice.String()),
			sdk.NewAttribute(types.AttributeKeyNewPrice, median.String()),
			sdk.NewAttribute(types.AttributeKeyValidatorCount, fmt.Sprintf("%d", len(submissions))),
			sdk.NewAttribute(types.AttributeKeyMedianPrice, median.String()),
			sdk.NewAttribute(types.AttributeKeyMeanPrice, mean.String()),
			sdk.NewAttribute(types.AttributeKeyDeviation, stdDev.String()),
		),
	)

	return nil
}

// getSubmissionsForWindow retrieves all submissions for a symbol in the given window
func (k Keeper) getSubmissionsForWindow(ctx sdk.Context, symbol string, windowStart, windowEnd uint64) []types.ValidatorPriceSubmission {
	store := ctx.KVStore(k.storeKey)
	var submissions []types.ValidatorPriceSubmission

	// Iterate through the height range
	for height := windowStart; height < windowEnd; height++ {
		// Get all validators' submissions for this height and symbol
		validators := k.GetAllOracleValidators(ctx)
		for _, validator := range validators {
			key := types.ValidatorSubmissionKey(validator.Validator, symbol, height)
			bz := store.Get(key)
			if bz != nil {
				var submission types.ValidatorPriceSubmission
				k.cdc.MustUnmarshal(bz, &submission)
				submissions = append(submissions, submission)
			}
		}
	}

	return submissions
}

// getSymbolsFromRecentSubmissions gets symbols from recent submissions
func (k Keeper) getSymbolsFromRecentSubmissions(ctx sdk.Context, windowStart, windowEnd uint64) map[string]bool {
	symbols := make(map[string]bool)
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ValidatorSubmissionKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var submission types.ValidatorPriceSubmission
		k.cdc.MustUnmarshal(iterator.Value(), &submission)

		if submission.BlockHeight >= windowStart && submission.BlockHeight < windowEnd {
			symbols[submission.Symbol] = true
		}
	}

	return symbols
}

// storeHistoricalPrice stores a price in the historical data
func (k Keeper) storeHistoricalPrice(ctx sdk.Context, symbol string, price sdk.Dec, blockHeight uint64) {
	historicalPrice := types.HistoricalPrice{
		Symbol:      symbol,
		Price:       price,
		BlockHeight: blockHeight,
		Timestamp:   ctx.BlockTime(),
		ValidatorCount: 0, // This could be enhanced to track validator count
	}

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&historicalPrice)
	store.Set(types.PriceHistoryKey(symbol, blockHeight), bz)
}

// GetExchangeRate retrieves an exchange rate for a currency pair
func (k Keeper) GetExchangeRate(ctx sdk.Context, base, target string) (types.ExchangeRate, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ExchangeRateKey(base, target))
	if bz == nil {
		return types.ExchangeRate{}, false
	}

	var exchangeRate types.ExchangeRate
	k.cdc.MustUnmarshal(bz, &exchangeRate)
	return exchangeRate, true
}

// GetAllExchangeRates retrieves all exchange rates
func (k Keeper) GetAllExchangeRates(ctx sdk.Context) []types.ExchangeRate {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ExchangeRateKeyPrefix)
	defer iterator.Close()

	var exchangeRates []types.ExchangeRate
	for ; iterator.Valid(); iterator.Next() {
		var exchangeRate types.ExchangeRate
		k.cdc.MustUnmarshal(iterator.Value(), &exchangeRate)
		exchangeRates = append(exchangeRates, exchangeRate)
	}

	return exchangeRates
}

// GetPriceHistory retrieves historical prices for a symbol
func (k Keeper) GetPriceHistory(ctx sdk.Context, symbol string, limit uint32) []types.HistoricalPrice {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, append(types.PriceHistoryKeyPrefix, []byte(symbol)...))
	defer iterator.Close()

	var historicalPrices []types.HistoricalPrice
	count := uint32(0)
	
	// Iterate in reverse to get latest prices first
	for ; iterator.Valid() && (limit == 0 || count < limit); iterator.Next() {
		var historicalPrice types.HistoricalPrice
		k.cdc.MustUnmarshal(iterator.Value(), &historicalPrice)
		historicalPrices = append(historicalPrices, historicalPrice)
		count++
	}

	return historicalPrices
}

// cleanupOldData removes old submissions and historical data
func (k Keeper) cleanupOldData(ctx sdk.Context) {
	params := k.GetParams(ctx)
	currentHeight := uint64(ctx.BlockHeight())
	
	// Only cleanup every 1000 blocks to avoid performance issues
	if currentHeight%1000 != 0 {
		return
	}

	// Remove submissions older than retention period (default 1 week in blocks)
	retentionBlocks := uint64(7 * 24 * 60 * 10) // ~1 week assuming 6s blocks
	if currentHeight > retentionBlocks {
		cutoffHeight := currentHeight - retentionBlocks
		k.removeOldSubmissions(ctx, cutoffHeight)
	}
}

// removeOldSubmissions removes validator submissions older than cutoff height
func (k Keeper) removeOldSubmissions(ctx sdk.Context, cutoffHeight uint64) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ValidatorSubmissionKeyPrefix)
	defer iterator.Close()

	var keysToDelete [][]byte
	for ; iterator.Valid(); iterator.Next() {
		var submission types.ValidatorPriceSubmission
		k.cdc.MustUnmarshal(iterator.Value(), &submission)
		
		if submission.BlockHeight < cutoffHeight {
			keysToDelete = append(keysToDelete, iterator.Key())
		}
	}

	// Delete old submissions
	for _, key := range keysToDelete {
		store.Delete(key)
	}
}

// checkStalePrices checks for stale price data and emits events
func (k Keeper) checkStalePrices(ctx sdk.Context) {
	params := k.GetParams(ctx)
	currentTime := ctx.BlockTime()
	
	allPrices := k.GetAllPriceData(ctx)
	for _, priceData := range allPrices {
		// Check if price is stale (older than staleness threshold)
		stalenessThreshold := time.Duration(params.StalenessThreshold) * time.Second
		if currentTime.Sub(priceData.LastUpdated) > stalenessThreshold {
			// Emit stale price event
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					"stale_price_detected",
					sdk.NewAttribute(types.AttributeKeySymbol, priceData.Symbol),
					sdk.NewAttribute(types.AttributeKeyPrice, priceData.Price.String()),
					sdk.NewAttribute("last_updated", priceData.LastUpdated.Format(time.RFC3339)),
					sdk.NewAttribute("staleness_duration", currentTime.Sub(priceData.LastUpdated).String()),
				),
			)
		}
	}
}

// updateValidatorStats updates validator performance statistics
func (k Keeper) updateValidatorStats(ctx sdk.Context) {
	// This could be enhanced to track validator performance over time
	// For now, we just log active validator count
	validators := k.GetAllOracleValidators(ctx)
	activeCount := 0
	
	for _, validator := range validators {
		if validator.Active {
			activeCount++
		}
	}

	// Emit validator stats event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"oracle_validator_stats",
			sdk.NewAttribute("total_validators", sdk.NewInt(int64(len(validators))).String()),
			sdk.NewAttribute("active_validators", sdk.NewInt(int64(activeCount)).String()),
			sdk.NewAttribute("block_height", sdk.NewInt(ctx.BlockHeight()).String()),
		),
	)
}