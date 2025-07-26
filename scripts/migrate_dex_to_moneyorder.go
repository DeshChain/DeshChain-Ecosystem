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

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"

	"github.com/DeshChain/DeshChain-Ecosystem/app"
	dextypes "github.com/DeshChain/DeshChain-Ecosystem/x/dex/types"
	moneyordertypes "github.com/DeshChain/DeshChain-Ecosystem/x/moneyorder/types"
)

// MigrationConfig holds configuration for the migration
type MigrationConfig struct {
	InputGenesisFile  string `json:"input_genesis_file"`
	OutputGenesisFile string `json:"output_genesis_file"`
	ChainID           string `json:"chain_id"`
	OverwriteTime     bool   `json:"overwrite_time"`
}

// PoolMigrationMapping maps old DEX pools to new Money Order pools
type PoolMigrationMapping struct {
	DEXPoolID     uint64 `json:"dex_pool_id"`
	PoolType      string `json:"pool_type"`      // "amm", "fixed_rate", "village"
	TokenA        string `json:"token_a"`
	TokenB        string `json:"token_b"`
	ExchangeRate  string `json:"exchange_rate,omitempty"`
	VillageCode   string `json:"village_code,omitempty"`
	CulturalTheme string `json:"cultural_theme,omitempty"`
}

// migrateDEXToMoneyOrder performs the migration from x/dex to x/moneyorder
func migrateDEXToMoneyOrder(appState map[string]json.RawMessage, clientCtx client.Context, config MigrationConfig) map[string]json.RawMessage {
	cdc := clientCtx.Codec

	// Extract DEX state
	var dexGenState dextypes.GenesisState
	if appState[dextypes.ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[dextypes.ModuleName], &dexGenState)
	}

	// Create new Money Order genesis state
	moneyOrderGenState := moneyordertypes.DefaultGenesis()
	
	// Migrate pools
	poolMappings := getPoolMigrationMappings()
	migratedPools := migratePools(dexGenState.Pools, poolMappings)
	
	// Set migrated pools in Money Order state
	moneyOrderGenState.AmmPools = migratedPools.AMMPools
	moneyOrderGenState.FixedRatePools = migratedPools.FixedRatePools
	moneyOrderGenState.VillagePools = migratedPools.VillagePools

	// Migrate liquidity positions
	moneyOrderGenState.LiquidityPositions = migrateLiquidityPositions(dexGenState.Positions)

	// Migrate trading pairs
	moneyOrderGenState.TradingPairs = migrateTradingPairs(dexGenState.TradingPairs)

	// Create cultural integration data
	moneyOrderGenState.CulturalQuotes = createCulturalQuotes()
	moneyOrderGenState.FestivalPeriods = createFestivalPeriods()

	// Set enhanced parameters
	moneyOrderGenState.Params = createEnhancedParams(dexGenState.Params)

	// Marshal the new state
	appState[moneyordertypes.ModuleName] = cdc.MustMarshalJSON(moneyOrderGenState)

	// Remove old DEX state
	delete(appState, dextypes.ModuleName)

	// Update app version and other metadata
	updateAppMetadata(appState, cdc, config)

	return appState
}

// getPoolMigrationMappings returns predefined mappings for pool migration
func getPoolMigrationMappings() []PoolMigrationMapping {
	return []PoolMigrationMapping{
		{
			DEXPoolID:     1,
			PoolType:      "fixed_rate",
			TokenA:        "unamo",
			TokenB:        "inr",
			ExchangeRate:  "0.075", // 1 NAMO = 0.075 INR
			CulturalTheme: "independence",
		},
		{
			DEXPoolID:     2,
			PoolType:      "amm",
			TokenA:        "unamo",
			TokenB:        "usdt",
			CulturalTheme: "prosperity",
		},
		{
			DEXPoolID:     3,
			PoolType:      "village",
			TokenA:        "unamo",
			TokenB:        "inr",
			VillageCode:   "110001",
			CulturalTheme: "community",
		},
	}
}

// MigratedPools holds the result of pool migration
type MigratedPools struct {
	AMMPools       []moneyordertypes.AMMPool       `json:"amm_pools"`
	FixedRatePools []moneyordertypes.FixedRatePool `json:"fixed_rate_pools"`
	VillagePools   []moneyordertypes.VillagePool   `json:"village_pools"`
}

// migratePools converts DEX pools to Money Order pools based on mappings
func migratePools(dexPools []dextypes.Pool, mappings []PoolMigrationMapping) MigratedPools {
	result := MigratedPools{
		AMMPools:       []moneyordertypes.AMMPool{},
		FixedRatePools: []moneyordertypes.FixedRatePool{},
		VillagePools:   []moneyordertypes.VillagePool{},
	}

	for _, dexPool := range dexPools {
		mapping := findMappingForPool(dexPool.Id, mappings)
		if mapping == nil {
			log.Printf("No mapping found for DEX pool %d, skipping", dexPool.Id)
			continue
		}

		switch mapping.PoolType {
		case "amm":
			ammPool := migrateToAMMPool(dexPool, *mapping)
			result.AMMPools = append(result.AMMPools, ammPool)

		case "fixed_rate":
			fixedPool := migrateToFixedRatePool(dexPool, *mapping)
			result.FixedRatePools = append(result.FixedRatePools, fixedPool)

		case "village":
			villagePool := migrateToVillagePool(dexPool, *mapping)
			result.VillagePools = append(result.VillagePools, villagePool)
		}
	}

	return result
}

// findMappingForPool finds the migration mapping for a given pool ID
func findMappingForPool(poolID uint64, mappings []PoolMigrationMapping) *PoolMigrationMapping {
	for _, mapping := range mappings {
		if mapping.DEXPoolID == poolID {
			return &mapping
		}
	}
	return nil
}

// migrateToAMMPool converts a DEX pool to AMM pool
func migrateToAMMPool(dexPool dextypes.Pool, mapping PoolMigrationMapping) moneyordertypes.AMMPool {
	return moneyordertypes.AMMPool{
		PoolId:              dexPool.Id,
		TokenA:              mapping.TokenA,
		TokenB:              mapping.TokenB,
		ReserveA:            dexPool.ReserveA,
		ReserveB:            dexPool.ReserveB,
		TotalShares:         dexPool.TotalShares,
		SwapFee:             dexPool.SwapFee,
		ExitFee:             sdk.NewDecWithPrec(5, 3), // 0.5% exit fee
		CreatedAt:           dexPool.CreatedAt,
		LastTradeAt:         dexPool.LastTradeAt,
		Status:              dexPool.Status,
		CulturalTheme:       mapping.CulturalTheme,
		FestivalMultiplier:  sdk.NewDecWithPrec(105, 2), // 5% festival bonus
		PatriotismScore:     sdk.NewDecWithPrec(95, 2),  // 95% patriotism score
		TradingVolume24h:    sdk.ZeroInt(),
		APY:                 sdk.NewDecWithPrec(15, 2), // 15% APY
	}
}

// migrateToFixedRatePool converts a DEX pool to fixed rate pool
func migrateToFixedRatePool(dexPool dextypes.Pool, mapping PoolMigrationMapping) moneyordertypes.FixedRatePool {
	exchangeRate, _ := sdk.NewDecFromStr(mapping.ExchangeRate)
	
	return moneyordertypes.FixedRatePool{
		PoolId:              dexPool.Id,
		TokenA:              mapping.TokenA,
		TokenB:              mapping.TokenB,
		ExchangeRate:        exchangeRate,
		ReserveA:            dexPool.ReserveA,
		ReserveB:            dexPool.ReserveB,
		TotalVolume:         sdk.ZeroInt(),
		CreatedBy:           sdk.AccAddress{}, // Will be updated during migration
		CreatedAt:           dexPool.CreatedAt,
		LastUpdated:         dexPool.LastTradeAt,
		Status:              dexPool.Status,
		CulturalTheme:       mapping.CulturalTheme,
		PriceStability:      sdk.NewDecWithPrec(99, 2), // 99% stability
		TrustScore:          sdk.NewDecWithPrec(98, 2), // 98% trust score
		MonthlyVolume:       sdk.ZeroInt(),
		PatriotismQuote:     getCulturalQuoteForTheme(mapping.CulturalTheme),
	}
}

// migrateToVillagePool converts a DEX pool to village pool
func migrateToVillagePool(dexPool dextypes.Pool, mapping PoolMigrationMapping) moneyordertypes.VillagePool {
	return moneyordertypes.VillagePool{
		PoolId:              dexPool.Id,
		VillageName:         getVillageNameFromCode(mapping.VillageCode),
		PostalCode:          mapping.VillageCode,
		TokenA:              mapping.TokenA,
		TokenB:              mapping.TokenB,
		ReserveA:            dexPool.ReserveA,
		ReserveB:            dexPool.ReserveB,
		TotalShares:         dexPool.TotalShares,
		Members:             []string{}, // Will be populated during migration
		Coordinator:         sdk.AccAddress{}, // Will be assigned
		SwapFee:             sdk.NewDecWithPrec(15, 3), // 0.15% reduced fee
		CreatedAt:           dexPool.CreatedAt,
		LastActivity:        dexPool.LastTradeAt,
		Status:              dexPool.Status,
		Verified:            true,
		CulturalTheme:       mapping.CulturalTheme,
		TrustScore:          sdk.NewDecWithPrec(96, 2), // 96% trust score
		CommunityImpact:     sdk.NewDecWithPrec(88, 2), // 88% community impact
		LocalEconomyBoost:   sdk.NewDecWithPrec(92, 2), // 92% local economy boost
		MonthlyTransactions: 0,
		PatriotismQuote:     getCulturalQuoteForTheme("community"),
	}
}

// migrateLiquidityPositions converts DEX positions to Money Order positions
func migrateLiquidityPositions(dexPositions []dextypes.Position) []moneyordertypes.LiquidityPosition {
	positions := make([]moneyordertypes.LiquidityPosition, 0, len(dexPositions))

	for _, dexPos := range dexPositions {
		position := moneyordertypes.LiquidityPosition{
			PositionId:          dexPos.Id,
			PoolId:              dexPos.PoolId,
			Owner:               dexPos.Owner,
			Shares:              dexPos.Shares,
			TokenAAmount:        dexPos.TokenAAmount,
			TokenBAmount:        dexPos.TokenBAmount,
			RewardsAccumulated:  sdk.NewCoins(),
			CreatedAt:           dexPos.CreatedAt,
			LastClaimAt:         dexPos.LastClaimAt,
			Status:              dexPos.Status,
			CulturalBonus:       sdk.NewCoin("unamo", sdk.ZeroInt()),
			PatriotismReward:    sdk.NewCoin("unamo", sdk.ZeroInt()),
			CommunityContribution: sdk.NewDecWithPrec(75, 2), // 75% community contribution
		}
		positions = append(positions, position)
	}

	return positions
}

// migrateTradingPairs converts DEX trading pairs to Money Order pairs
func migrateTradingPairs(dexPairs []dextypes.TradingPair) []moneyordertypes.TradingPair {
	pairs := make([]moneyordertypes.TradingPair, 0, len(dexPairs))

	for _, dexPair := range dexPairs {
		pair := moneyordertypes.TradingPair{
			PairId:              dexPair.Id,
			TokenA:              dexPair.TokenA,
			TokenB:              dexPair.TokenB,
			Active:              dexPair.Active,
			MinTradeSize:        dexPair.MinTradeSize,
			MaxTradeSize:        dexPair.MaxTradeSize,
			CreatedAt:           dexPair.CreatedAt,
			LastTradeAt:         dexPair.LastTradeAt,
			TotalVolume:         dexPair.TotalVolume,
			CulturalSignificance: "moderate",
			RegionalPopularity:   sdk.NewDecWithPrec(65, 2), // 65% regional popularity
			EducationalValue:     sdk.NewDecWithPrec(70, 2), // 70% educational value
			CommunitySupport:     sdk.NewDecWithPrec(80, 2), // 80% community support
		}
		pairs = append(pairs, pair)
	}

	return pairs
}

// createCulturalQuotes generates cultural quotes for the Money Order system
func createCulturalQuotes() []moneyordertypes.CulturalQuote {
	return []moneyordertypes.CulturalQuote{
		{
			QuoteId:    1,
			Text:       "सत्यमेव जयते - Truth alone triumphs",
			Author:     "Mundaka Upanishad",
			Category:   "truth",
			Language:   "sanskrit",
			Occasion:   "general",
			Active:     true,
			Weight:     10,
		},
		{
			QuoteId:    2,
			Text:       "Be the change you wish to see in the world",
			Author:     "Mahatma Gandhi",
			Category:   "change",
			Language:   "english",
			Occasion:   "motivation",
			Active:     true,
			Weight:     10,
		},
		{
			QuoteId:    3,
			Text:       "यत्र नार्यस्तु पूज्यन्ते रमन्ते तत्र देवताः - Where women are honored, divinity blossoms there",
			Author:     "Manusmriti",
			Category:   "women_empowerment",
			Language:   "sanskrit",
			Occasion:   "women_day",
			Active:     true,
			Weight:     10,
		},
	}
}

// createFestivalPeriods generates festival periods for cultural integration
func createFestivalPeriods() []moneyordertypes.FestivalPeriod {
	return []moneyordertypes.FestivalPeriod{
		{
			FestivalId:    1,
			Name:          "Diwali",
			Description:   "Festival of Lights",
			StartDate:     "2024-11-01",
			EndDate:       "2024-11-05",
			BonusRate:     sdk.NewDecWithPrec(15, 2), // 15% bonus
			CulturalTheme: "prosperity",
			Region:        "pan_india",
			Active:        true,
		},
		{
			FestivalId:    2,
			Name:          "Holi",
			Description:   "Festival of Colors",
			StartDate:     "2024-03-13",
			EndDate:       "2024-03-14",
			BonusRate:     sdk.NewDecWithPrec(10, 2), // 10% bonus
			CulturalTheme: "unity",
			Region:        "north_india",
			Active:        true,
		},
		{
			FestivalId:    3,
			Name:          "Independence Day",
			Description:   "National Independence Day",
			StartDate:     "2024-08-15",
			EndDate:       "2024-08-15",
			BonusRate:     sdk.NewDecWithPrec(20, 2), // 20% patriotism bonus
			CulturalTheme: "patriotism",
			Region:        "pan_india",
			Active:        true,
		},
	}
}

// createEnhancedParams creates enhanced parameters for Money Order module
func createEnhancedParams(dexParams dextypes.Params) moneyordertypes.Params {
	return moneyordertypes.Params{
		TradingFee:              dexParams.TradingFee,
		ProtocolFee:             sdk.NewDecWithPrec(5, 3), // 0.5% protocol fee
		MaxSlippageTolerance:    sdk.NewDecWithPrec(5, 2), // 5% max slippage
		MinLiquidityThreshold:   sdk.NewInt(1000000),      // 1M minimum liquidity
		CulturalBonusRate:       sdk.NewDecWithPrec(5, 2), // 5% cultural bonus
		FestivalBonusMultiplier: sdk.NewDecWithPrec(15, 2), // 15% festival multiplier
		PatriotismThreshold:     sdk.NewDecWithPrec(70, 2), // 70% patriotism threshold
		VillagePoolDiscount:     sdk.NewDecWithPrec(25, 2), // 25% village discount
		MoneyOrderFee:           sdk.NewDecWithPrec(1, 3),  // 0.1% money order fee
		MaxDailyVolume:          sdk.NewInt(10000000000),   // 10B daily volume limit
		EnableCulturalFeatures:  true,
		EnableFestivalBonuses:   true,
		EnablePatriotismRewards: true,
	}
}

// Helper functions

func getCulturalQuoteForTheme(theme string) string {
	quotes := map[string]string{
		"independence": "स्वतंत्रता हमारा जन्मसिद्ध अधिकार है - Freedom is our birthright",
		"prosperity":   "सर्वे भवन्तु सुखिनः - May all beings be happy",
		"community":    "वसुधैव कुटुम्बकम् - The world is one family",
		"unity":        "एकता में शक्ति है - There is strength in unity",
	}
	
	if quote, exists := quotes[theme]; exists {
		return quote
	}
	return "सत्यमेव जयते - Truth alone triumphs"
}

func getVillageNameFromCode(postalCode string) string {
	villages := map[string]string{
		"110001": "Central Delhi Village",
		"400001": "Mumbai Fort Village",
		"560001": "Bangalore Central Village",
		"500001": "Hyderabad Central Village",
	}
	
	if name, exists := villages[postalCode]; exists {
		return name
	}
	return fmt.Sprintf("Village %s", postalCode)
}

func updateAppMetadata(appState map[string]json.RawMessage, cdc codec.Codec, config MigrationConfig) {
	// Update genesis metadata
	var genDoc genutiltypes.GenesisDoc
	if appState["genutil"] != nil {
		var genState genutil.GenesisState
		cdc.MustUnmarshalJSON(appState["genutil"], &genState)
		
		// Add migration metadata
		genState.GenTxs = append(genState.GenTxs, json.RawMessage(fmt.Sprintf(`{
			"migration": {
				"from_module": "dex",
				"to_module": "moneyorder",
				"migration_time": "%s",
				"version": "v2.0.0"
			}
		}`, time.Now().Format(time.RFC3339))))
		
		appState["genutil"] = cdc.MustMarshalJSON(&genState)
	}
}

// main function to run the migration
func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: migrate_dex_to_moneyorder <config_file>")
	}

	configFile := os.Args[1]
	configData, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	var config MigrationConfig
	if err := json.Unmarshal(configData, &config); err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	// Read genesis file
	genesisData, err := os.ReadFile(config.InputGenesisFile)
	if err != nil {
		log.Fatalf("Failed to read genesis file: %v", err)
	}

	var genesisDoc genutiltypes.GenesisDoc
	if err := json.Unmarshal(genesisData, &genesisDoc); err != nil {
		log.Fatalf("Failed to parse genesis: %v", err)
	}

	// Setup codec
	config := app.MakeEncodingConfig()
	cdc := config.Marshaler
	clientCtx := client.Context{}.WithCodec(cdc)

	// Parse app state
	var appState map[string]json.RawMessage
	if err := json.Unmarshal(genesisDoc.AppState, &appState); err != nil {
		log.Fatalf("Failed to parse app state: %v", err)
	}

	// Perform migration
	migratedAppState := migrateDEXToMoneyOrder(appState, clientCtx, config)

	// Update genesis doc
	newAppState, err := json.Marshal(migratedAppState)
	if err != nil {
		log.Fatalf("Failed to marshal app state: %v", err)
	}
	genesisDoc.AppState = newAppState

	if config.OverwriteTime {
		genesisDoc.GenesisTime = time.Now()
	}

	if config.ChainID != "" {
		genesisDoc.ChainID = config.ChainID
	}

	// Write migrated genesis
	migratedGenesis, err := json.MarshalIndent(genesisDoc, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal genesis: %v", err)
	}

	if err := os.WriteFile(config.OutputGenesisFile, migratedGenesis, 0644); err != nil {
		log.Fatalf("Failed to write genesis file: %v", err)
	}

	fmt.Printf("Successfully migrated genesis from DEX to Money Order\n")
	fmt.Printf("Input: %s\n", config.InputGenesisFile)
	fmt.Printf("Output: %s\n", config.OutputGenesisFile)
	fmt.Printf("Migration completed at: %s\n", time.Now().Format(time.RFC3339))
}