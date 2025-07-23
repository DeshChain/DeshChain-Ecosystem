package types

import (
	"fmt"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "dusd"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_dusd"

	// DUSD denomination
	DUSDDenom = "dusd"
)

// DUSD Fee Constants (USD equivalent)
var (
	// Base fee: $0.10 USD
	BaseFeeUSD = sdk.NewDecWithPrec(10, 2) // 0.10
	
	// Max fee: $1.00 USD  
	MaxFeeUSD = sdk.NewDecWithPrec(100, 2) // 1.00
	
	// Default fee percentage: 0.25%
	DefaultFeePercentage = sdk.NewDecWithPrec(25, 4) // 0.0025
	
	// Target price: $1.00 USD
	TargetPriceUSD = sdk.OneDec()
	
	// Price tolerance: 1%
	DefaultPriceTolerance = sdk.NewDecWithPrec(1, 2) // 0.01
	
	// Min collateral ratio: 150%
	DefaultMinCollateralRatio = sdk.NewDecWithPrec(150, 2) // 1.50
	
	// Liquidation ratio: 120%
	DefaultLiquidationRatio = sdk.NewDecWithPrec(120, 2) // 1.20
)

// Key prefixes for store
var (
	PositionKey        = []byte{0x01}
	ParamsKey          = []byte{0x02}
	PriceDataKey       = []byte{0x03}
	StabilityActionKey = []byte{0x04}
	ReserveStatsKey    = []byte{0x05}
	SupplyStatsKey     = []byte{0x06}
)

// GetPositionStoreKey returns the store key for a position
func GetPositionStoreKey(positionID string) []byte {
	return append(PositionKey, []byte(positionID)...)
}

// GetPriceDataStoreKey returns the store key for price data
func GetPriceDataStoreKey(source string, timestamp time.Time) []byte {
	key := append(PriceDataKey, []byte(source)...)
	return append(key, sdk.Uint64ToBigEndian(uint64(timestamp.Unix()))...)
}

// GetStabilityActionStoreKey returns the store key for stability actions
func GetStabilityActionStoreKey(actionID string) []byte {
	return append(StabilityActionKey, []byte(actionID)...)
}

// CalculateFee calculates DUSD transaction fee using same logic as DINR
func CalculateFee(amount sdk.Coin) sdk.Coin {
	// Convert amount to USD equivalent for fee calculation
	amountDec := amount.Amount.ToLegacyDec()
	
	// Calculate percentage fee
	percentageFee := amountDec.Mul(DefaultFeePercentage)
	
	// Compare with base fee (minimum)
	if percentageFee.LT(BaseFeeUSD) {
		percentageFee = BaseFeeUSD
	}
	
	// Cap at max fee
	if percentageFee.GT(MaxFeeUSD) {
		percentageFee = MaxFeeUSD
	}
	
	// Convert back to integer amount
	feeAmount := percentageFee.TruncateInt()
	
	return sdk.NewCoin(DUSDDenom, feeAmount)
}

// CalculateHealthFactor calculates position health factor (same logic as DINR)
func CalculateHealthFactor(collateralValue sdk.Dec, debtValue sdk.Dec, liquidationRatio sdk.Dec) sdk.Dec {
	if debtValue.IsZero() {
		return sdk.NewDec(1000) // Max health factor for positions with no debt
	}
	
	// Health Factor = (Collateral Value * Liquidation Ratio) / Debt Value
	return collateralValue.Mul(liquidationRatio).Quo(debtValue)
}

// IsPositionLiquidatable checks if a position can be liquidated
func IsPositionLiquidatable(healthFactor sdk.Dec) bool {
	return healthFactor.LT(sdk.OneDec())
}

// CalculateCollateralRequired calculates required collateral for DUSD amount
func CalculateCollateralRequired(dusdAmount sdk.Dec, collateralPrice sdk.Dec, minCollateralRatio sdk.Dec) sdk.Dec {
	// Required Collateral = (DUSD Amount / Collateral Price) * Min Collateral Ratio
	return dusdAmount.Quo(collateralPrice).Mul(minCollateralRatio)
}

// ValidateAddress validates a cosmos address
func ValidateAddress(address string) error {
	_, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return fmt.Errorf("invalid address %s: %w", address, err)
	}
	return nil
}

// ValidateDenom validates a coin denomination
func ValidateDenom(denom string) error {
	if denom == "" {
		return fmt.Errorf("denomination cannot be empty")
	}
	if len(denom) > 128 {
		return fmt.Errorf("denomination too long: %d > 128", len(denom))
	}
	return nil
}

// ValidatePositiveAmount validates that an amount is positive
func ValidatePositiveAmount(amount sdk.Coin) error {
	if amount.Amount.IsNil() || amount.Amount.LTE(sdk.ZeroInt()) {
		return fmt.Errorf("amount must be positive, got %s", amount)
	}
	return nil
}

// GeneratePositionID generates a unique position ID
func GeneratePositionID(owner string, timestamp time.Time) string {
	return fmt.Sprintf("%s-%d", owner, timestamp.Unix())
}

// GenerateActionID generates a unique action ID
func GenerateActionID(actionType string, timestamp time.Time) string {
	return fmt.Sprintf("%s-%d", actionType, timestamp.Unix())
}

// DefaultParams returns default module parameters
func DefaultParams() DUSDParams {
	return DUSDParams{
		TargetPrice:              TargetPriceUSD.String(),
		PriceTolerance:           DefaultPriceTolerance.String(),
		MinCollateralRatio:       DefaultMinCollateralRatio.String(),
		LiquidationRatio:         DefaultLiquidationRatio.String(),
		BaseFeeUsd:              BaseFeeUSD.String(),
		MaxFeeUsd:               MaxFeeUSD.String(),
		FeePercentage:           DefaultFeePercentage.String(),
		OracleSources:           []string{"chainlink", "federal_reserve", "band", "pyth"},
		PriceDeviationThreshold: sdk.NewDecWithPrec(5, 3).String(), // 0.5%
		OracleTimeoutSeconds:    300,                               // 5 minutes
		RebalanceThreshold:      sdk.NewDecWithPrec(5, 3).String(), // 0.5%
		EmergencyThreshold:      sdk.NewDecWithPrec(2, 2).String(), // 2%
		CircuitBreakerEnabled:   true,
		ReserveRatio:            sdk.NewDecWithPrec(20, 2).String(), // 20%
		AcceptedCollateral:      []string{"NAMO", "USDC", "USDT"},
	}
}