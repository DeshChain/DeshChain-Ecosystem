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

package types

import (
	"fmt"
	"time"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// FixedRatePool represents a pool with fixed exchange rates
// Similar to traditional money order counters with fixed rates
type FixedRatePool struct {
	// Pool identification
	PoolId          uint64         `json:"pool_id" yaml:"pool_id"`
	PoolName        string         `json:"pool_name" yaml:"pool_name"`
	Description     string         `json:"description" yaml:"description"`
	
	// Token pair configuration
	Token0Denom     string         `json:"token0_denom" yaml:"token0_denom"` // e.g., "unamo"
	Token1Denom     string         `json:"token1_denom" yaml:"token1_denom"` // e.g., "uinr"
	ExchangeRate    sdk.Dec        `json:"exchange_rate" yaml:"exchange_rate"` // Fixed rate
	ReverseRate     sdk.Dec        `json:"reverse_rate" yaml:"reverse_rate"` // 1/ExchangeRate
	
	// Pool limits (like traditional money order limits)
	MinOrderAmount  sdk.Int        `json:"min_order_amount" yaml:"min_order_amount"`
	MaxOrderAmount  sdk.Int        `json:"max_order_amount" yaml:"max_order_amount"`
	DailyLimit      sdk.Int        `json:"daily_limit" yaml:"daily_limit"`
	MonthlyLimit    sdk.Int        `json:"monthly_limit" yaml:"monthly_limit"`
	
	// Liquidity management
	Token0Balance   sdk.Int        `json:"token0_balance" yaml:"token0_balance"`
	Token1Balance   sdk.Int        `json:"token1_balance" yaml:"token1_balance"`
	ReservedBalance sdk.Int        `json:"reserved_balance" yaml:"reserved_balance"`
	
	// Geographic and KYC requirements
	SupportedRegions []string      `json:"supported_regions" yaml:"supported_regions"` // Postal codes
	RequiresKYC     bool           `json:"requires_kyc" yaml:"requires_kyc"`
	KYCThreshold    sdk.Int        `json:"kyc_threshold" yaml:"kyc_threshold"`
	
	// Fee structure
	BaseFee         sdk.Dec        `json:"base_fee" yaml:"base_fee"` // Base fee percentage
	ExpressFee      sdk.Dec        `json:"express_fee" yaml:"express_fee"` // Additional for instant
	BulkDiscount    sdk.Dec        `json:"bulk_discount" yaml:"bulk_discount"` // Discount for bulk orders
	
	// Cultural features
	CulturalPair    bool           `json:"cultural_pair" yaml:"cultural_pair"`
	FestivalBonus   bool           `json:"festival_bonus" yaml:"festival_bonus"`
	VillagePriority bool           `json:"village_priority" yaml:"village_priority"`
	
	// Pool status
	Active          bool           `json:"active" yaml:"active"`
	MaintenanceMode bool           `json:"maintenance_mode" yaml:"maintenance_mode"`
	CreatedBy       sdk.AccAddress `json:"created_by" yaml:"created_by"`
	CreatedAt       time.Time      `json:"created_at" yaml:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at" yaml:"updated_at"`
	
	// Statistics
	TotalOrders     uint64         `json:"total_orders" yaml:"total_orders"`
	TotalVolume     sdk.Int        `json:"total_volume" yaml:"total_volume"`
	DailyVolume     sdk.Int        `json:"daily_volume" yaml:"daily_volume"`
	MonthlyVolume   sdk.Int        `json:"monthly_volume" yaml:"monthly_volume"`
}

// FixedRateOrder represents an order in fixed rate pool
type FixedRateOrder struct {
	OrderId         string         `json:"order_id" yaml:"order_id"`
	PoolId          uint64         `json:"pool_id" yaml:"pool_id"`
	Sender          sdk.AccAddress `json:"sender" yaml:"sender"`
	Receiver        sdk.AccAddress `json:"receiver" yaml:"receiver"`
	
	// Order details
	InputAmount     sdk.Coin       `json:"input_amount" yaml:"input_amount"`
	OutputAmount    sdk.Coin       `json:"output_amount" yaml:"output_amount"`
	ExchangeRate    sdk.Dec        `json:"exchange_rate" yaml:"exchange_rate"`
	Fees            sdk.Coin       `json:"fees" yaml:"fees"`
	
	// Order type
	OrderType       string         `json:"order_type" yaml:"order_type"` // "normal", "express", "bulk"
	Priority        uint32         `json:"priority" yaml:"priority"`
	
	// Status tracking
	Status          string         `json:"status" yaml:"status"`
	CreatedAt       time.Time      `json:"created_at" yaml:"created_at"`
	ExpiresAt       time.Time      `json:"expires_at" yaml:"expires_at"`
	CompletedAt     time.Time      `json:"completed_at" yaml:"completed_at"`
}

// NewFixedRatePool creates a new fixed rate pool
func NewFixedRatePool(
	poolId uint64,
	token0Denom string,
	token1Denom string,
	exchangeRate sdk.Dec,
	creator sdk.AccAddress,
) *FixedRatePool {
	now := time.Now()
	
	return &FixedRatePool{
		PoolId:          poolId,
		PoolName:        fmt.Sprintf("%s-%s Fixed Rate", token0Denom, token1Denom),
		Token0Denom:     token0Denom,
		Token1Denom:     token1Denom,
		ExchangeRate:    exchangeRate,
		ReverseRate:     sdk.OneDec().Quo(exchangeRate),
		
		// Default limits (like traditional money orders)
		MinOrderAmount:  MinMoneyOrderAmount,
		MaxOrderAmount:  MaxMoneyOrderAmount,
		DailyLimit:      MaxDailyUserLimit,
		MonthlyLimit:    MaxDailyUserLimit.Mul(sdk.NewInt(30)),
		
		// Initial balances
		Token0Balance:   sdk.ZeroInt(),
		Token1Balance:   sdk.ZeroInt(),
		ReservedBalance: sdk.ZeroInt(),
		
		// Default fees
		BaseFee:         sdk.MustNewDecFromStr(MoneyOrderBaseFee),
		ExpressFee:      sdk.MustNewDecFromStr(ExpressDeliveryFee),
		BulkDiscount:    sdk.MustNewDecFromStr(BulkOrderDiscount),
		
		// KYC settings
		RequiresKYC:     true,
		KYCThreshold:    KYCRequiredThreshold,
		
		// Status
		Active:          true,
		CreatedBy:       creator,
		CreatedAt:       now,
		UpdatedAt:       now,
		
		// Initialize statistics
		TotalOrders:     0,
		TotalVolume:     sdk.ZeroInt(),
		DailyVolume:     sdk.ZeroInt(),
		MonthlyVolume:   sdk.ZeroInt(),
	}
}

// CalculateOutput calculates output amount for given input
func (p *FixedRatePool) CalculateOutput(inputAmount sdk.Int, isForward bool) (sdk.Int, sdk.Int, error) {
	if !p.Active {
		return sdk.ZeroInt(), sdk.ZeroInt(), fmt.Errorf("pool is not active")
	}
	
	// Check limits
	if inputAmount.LT(p.MinOrderAmount) {
		return sdk.ZeroInt(), sdk.ZeroInt(), fmt.Errorf("amount below minimum: %s", p.MinOrderAmount)
	}
	
	if inputAmount.GT(p.MaxOrderAmount) {
		return sdk.ZeroInt(), sdk.ZeroInt(), fmt.Errorf("amount exceeds maximum: %s", p.MaxOrderAmount)
	}
	
	// Calculate base output
	var outputAmount sdk.Dec
	if isForward {
		outputAmount = p.ExchangeRate.MulInt(inputAmount)
	} else {
		outputAmount = p.ReverseRate.MulInt(inputAmount)
	}
	
	// Calculate fees
	feeAmount := p.BaseFee.MulInt(inputAmount)
	
	// Apply fee to output
	netOutput := outputAmount.Sub(feeAmount)
	
	if netOutput.IsNegative() {
		return sdk.ZeroInt(), sdk.ZeroInt(), fmt.Errorf("fees exceed output amount")
	}
	
	return netOutput.TruncateInt(), feeAmount.TruncateInt(), nil
}

// CanFulfillOrder checks if pool has sufficient liquidity
func (p *FixedRatePool) CanFulfillOrder(outputAmount sdk.Int, outputDenom string) bool {
	if outputDenom == p.Token0Denom {
		return p.Token0Balance.Sub(p.ReservedBalance).GTE(outputAmount)
	} else if outputDenom == p.Token1Denom {
		return p.Token1Balance.Sub(p.ReservedBalance).GTE(outputAmount)
	}
	return false
}

// UpdateVolume updates pool volume statistics
func (p *FixedRatePool) UpdateVolume(amount sdk.Int) {
	p.TotalVolume = p.TotalVolume.Add(amount)
	p.DailyVolume = p.DailyVolume.Add(amount)
	p.MonthlyVolume = p.MonthlyVolume.Add(amount)
	p.TotalOrders++
	p.UpdatedAt = time.Now()
}

// GetFeeForOrderType returns fee based on order type
func (p *FixedRatePool) GetFeeForOrderType(orderType string, amount sdk.Int) sdk.Dec {
	baseFee := p.BaseFee
	
	switch orderType {
	case "express":
		return baseFee.Add(p.ExpressFee)
	case "bulk":
		discount := baseFee.Mul(p.BulkDiscount)
		return baseFee.Sub(discount)
	default:
		return baseFee
	}
}

// ValidatePool ensures pool parameters are valid
func (p *FixedRatePool) ValidatePool() error {
	if p.Token0Denom == "" || p.Token1Denom == "" {
		return fmt.Errorf("token denoms cannot be empty")
	}
	
	if p.Token0Denom == p.Token1Denom {
		return fmt.Errorf("token denoms must be different")
	}
	
	if p.ExchangeRate.IsZero() || p.ExchangeRate.IsNegative() {
		return fmt.Errorf("invalid exchange rate")
	}
	
	if p.MinOrderAmount.GT(p.MaxOrderAmount) {
		return fmt.Errorf("min order amount cannot exceed max order amount")
	}
	
	if p.BaseFee.IsNegative() || p.BaseFee.GTE(sdk.OneDec()) {
		return fmt.Errorf("invalid base fee")
	}
	
	return nil
}

// IsEligibleForCulturalBonus checks if pool qualifies for cultural bonuses
func (p *FixedRatePool) IsEligibleForCulturalBonus() bool {
	return p.CulturalPair && p.FestivalBonus && p.Active
}

// GetEffectiveFee returns the effective fee after all discounts
func (p *FixedRatePool) GetEffectiveFee(
	orderType string,
	amount sdk.Int,
	isFestival bool,
	isVillageUser bool,
) sdk.Dec {
	// Start with base fee for order type
	fee := p.GetFeeForOrderType(orderType, amount)
	
	// Apply cultural discounts
	if p.IsEligibleForCulturalBonus() {
		if isFestival {
			festivalDiscount := sdk.MustNewDecFromStr(FestivalDiscount)
			fee = fee.Mul(sdk.OneDec().Sub(festivalDiscount))
		}
		
		if isVillageUser && p.VillagePriority {
			villageDiscount := sdk.MustNewDecFromStr(VillagePoolDiscount)
			fee = fee.Mul(sdk.OneDec().Sub(villageDiscount))
		}
	}
	
	// Ensure fee doesn't go below minimum
	minFee := sdk.MustNewDecFromStr("0.0001") // 0.01% minimum
	if fee.LT(minFee) {
		fee = minFee
	}
	
	return fee
}