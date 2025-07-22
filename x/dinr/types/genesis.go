package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultGenesisTime returns a default genesis time
func DefaultGenesisTime() time.Time {
	return time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
}

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
		UserPositions: []UserPosition{},
		StabilityData: StabilityData{
			CurrentPrice:          "1.00",
			TargetPrice:           "1.00",
			PriceDeviation:        0,
			TotalSupply:           sdk.NewCoin(DINRDenom, sdk.ZeroInt()),
			TotalCollateralValue:  sdk.NewCoin("inr", sdk.ZeroInt()),
			GlobalCollateralRatio: 0,
			LastUpdate:            DefaultGenesisTime(),
		},
		InsuranceFund: InsuranceFund{
			Balance:       sdk.NewCoin(DINRDenom, sdk.ZeroInt()),
			TargetRatio:   DefaultInsuranceFundRatio,
			Assets:        []sdk.Coin{},
			LastRebalance: DefaultGenesisTime(),
		},
		YieldStrategies: []YieldStrategy{
			{
				Id:             "gramsuraksha",
				Name:           "Gram Suraksha Pool Integration",
				Protocol:       "gramsuraksha",
				Allocation:     3000, // 30%
				ExpectedApy:    "12.0",
				IsActive:       true,
				DeployedAmount: sdk.NewCoin(DINRDenom, sdk.ZeroInt()),
			},
			{
				Id:             "urbansuraksha",
				Name:           "Urban Suraksha Pool Integration",
				Protocol:       "urbansuraksha",
				Allocation:     2000, // 20%
				ExpectedApy:    "10.0",
				IsActive:       true,
				DeployedAmount: sdk.NewCoin(DINRDenom, sdk.ZeroInt()),
			},
			{
				Id:             "internal_lending",
				Name:           "Internal Lending Protocols",
				Protocol:       "internal_lending",
				Allocation:     3000, // 30%
				ExpectedApy:    "8.0",
				IsActive:       true,
				DeployedAmount: sdk.NewCoin(DINRDenom, sdk.ZeroInt()),
			},
		},
		TotalDinrMinted:   "0",
		TotalFeesCollected: "0",
	}
}

// ValidateGenesis performs basic genesis state validation
func ValidateGenesis(data GenesisState) error {
	// Validate params
	if err := data.Params.Validate(); err != nil {
		return err
	}

	// Validate user positions
	userMap := make(map[string]bool)
	for _, position := range data.UserPositions {
		if userMap[position.Address] {
			return fmt.Errorf("duplicate user position for address %s", position.Address)
		}
		userMap[position.Address] = true

		// Validate position
		if err := validateUserPosition(position); err != nil {
			return err
		}
	}

	// Validate stability data
	if err := validateStabilityData(data.StabilityData); err != nil {
		return err
	}

	// Validate insurance fund
	if err := validateInsuranceFund(data.InsuranceFund); err != nil {
		return err
	}

	// Validate yield strategies
	strategyMap := make(map[string]bool)
	totalAllocation := uint64(0)
	for _, strategy := range data.YieldStrategies {
		if strategyMap[strategy.Id] {
			return fmt.Errorf("duplicate yield strategy ID %s", strategy.Id)
		}
		strategyMap[strategy.Id] = true

		if err := validateYieldStrategy(strategy); err != nil {
			return err
		}

		totalAllocation += strategy.Allocation
	}

	// Total allocation can be less than 100% but not more
	if totalAllocation > 10000 {
		return fmt.Errorf("total yield strategy allocation exceeds 100%%")
	}

	// Validate total minted and fees
	totalMinted, ok := sdk.NewIntFromString(data.TotalDinrMinted)
	if !ok || totalMinted.IsNegative() {
		return fmt.Errorf("invalid total DINR minted amount")
	}

	totalFees, ok := sdk.NewIntFromString(data.TotalFeesCollected)
	if !ok || totalFees.IsNegative() {
		return fmt.Errorf("invalid total fees collected amount")
	}

	return nil
}

func validateUserPosition(position UserPosition) error {
	_, err := sdk.AccAddressFromBech32(position.Address)
	if err != nil {
		return fmt.Errorf("invalid user address: %w", err)
	}

	// Validate collateral
	if err := position.Collateral.Validate(); err != nil {
		return fmt.Errorf("invalid collateral: %w", err)
	}

	// Validate DINR minted
	if !position.DinrMinted.IsValid() || position.DinrMinted.IsNegative() {
		return fmt.Errorf("invalid DINR minted amount")
	}

	if position.DinrMinted.Denom != DINRDenom {
		return fmt.Errorf("DINR minted denom must be %s", DINRDenom)
	}

	// Validate health factor
	healthFactor, err := sdk.NewDecFromStr(position.HealthFactor)
	if err != nil || healthFactor.IsNegative() {
		return fmt.Errorf("invalid health factor")
	}

	return nil
}

func validateStabilityData(data StabilityData) error {
	currentPrice, err := sdk.NewDecFromStr(data.CurrentPrice)
	if err != nil || currentPrice.IsNegative() || currentPrice.IsZero() {
		return fmt.Errorf("invalid current price")
	}

	targetPrice, err := sdk.NewDecFromStr(data.TargetPrice)
	if err != nil || targetPrice.IsNegative() || targetPrice.IsZero() {
		return fmt.Errorf("invalid target price")
	}

	if !data.TotalSupply.IsValid() || data.TotalSupply.IsNegative() {
		return fmt.Errorf("invalid total supply")
	}

	if data.TotalSupply.Denom != DINRDenom {
		return fmt.Errorf("total supply denom must be %s", DINRDenom)
	}

	if !data.TotalCollateralValue.IsValid() || data.TotalCollateralValue.IsNegative() {
		return fmt.Errorf("invalid total collateral value")
	}

	return nil
}

func validateInsuranceFund(fund InsuranceFund) error {
	if !fund.Balance.IsValid() || fund.Balance.IsNegative() {
		return fmt.Errorf("invalid insurance fund balance")
	}

	if fund.Balance.Denom != DINRDenom {
		return fmt.Errorf("insurance fund balance denom must be %s", DINRDenom)
	}

	targetRatio, err := sdk.NewDecFromStr(fund.TargetRatio)
	if err != nil || targetRatio.IsNegative() || targetRatio.GT(sdk.OneDec()) {
		return fmt.Errorf("invalid target ratio")
	}

	if err := fund.Assets.Validate(); err != nil {
		return fmt.Errorf("invalid insurance fund assets: %w", err)
	}

	return nil
}

func validateYieldStrategy(strategy YieldStrategy) error {
	if strategy.Id == "" {
		return fmt.Errorf("yield strategy ID cannot be empty")
	}

	if strategy.Name == "" {
		return fmt.Errorf("yield strategy name cannot be empty")
	}

	if strategy.Protocol == "" {
		return fmt.Errorf("yield strategy protocol cannot be empty")
	}

	if strategy.Allocation > 10000 {
		return fmt.Errorf("yield strategy allocation cannot exceed 100%%")
	}

	apy, err := sdk.NewDecFromStr(strategy.ExpectedApy)
	if err != nil || apy.IsNegative() {
		return fmt.Errorf("invalid expected APY")
	}

	if !strategy.DeployedAmount.IsValid() || strategy.DeployedAmount.IsNegative() {
		return fmt.Errorf("invalid deployed amount")
	}

	return nil
}