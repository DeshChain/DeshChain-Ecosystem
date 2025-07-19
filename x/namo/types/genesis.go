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

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:              DefaultParams(),
		TokenSupply:         DefaultTokenSupply(),
		VestingSchedules:    []VestingSchedule{},
		DistributionEvents:  []TokenDistributionEvent{},
	}
}

// DefaultTokenSupply returns the default token supply configuration
func DefaultTokenSupply() TokenSupply {
	return TokenSupply{
		TotalSupply:           sdk.NewInt(int64(TotalSupply * 1_000_000)), // 6 decimals
		PublicSaleAllocation:  sdk.NewInt(int64(PublicSaleAllocation * 1_000_000)),
		LiquidityAllocation:   sdk.NewInt(int64(LiquidityAllocation * 1_000_000)),
		TeamAllocation:        sdk.NewInt(int64(TeamAllocation * 1_000_000)),
		DevelopmentAllocation: sdk.NewInt(int64(DevelopmentAllocation * 1_000_000)),
		CommunityAllocation:   sdk.NewInt(int64(CommunityAllocation * 1_000_000)),
		DaoTreasuryAllocation: sdk.NewInt(int64(DAOTreasuryAllocation * 1_000_000)),
		InitialBurnAllocation: sdk.NewInt(int64(InitialBurnAllocation * 1_000_000)),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Validate params
	if err := gs.Params.Validate(); err != nil {
		return fmt.Errorf("invalid params: %w", err)
	}

	// Validate token supply
	if err := gs.TokenSupply.Validate(); err != nil {
		return fmt.Errorf("invalid token supply: %w", err)
	}

	// Validate vesting schedules
	for i, schedule := range gs.VestingSchedules {
		if err := schedule.Validate(); err != nil {
			return fmt.Errorf("invalid vesting schedule %d: %w", i, err)
		}
	}

	// Validate distribution events
	for i, event := range gs.DistributionEvents {
		if err := event.Validate(); err != nil {
			return fmt.Errorf("invalid distribution event %d: %w", i, err)
		}
	}

	return nil
}

// Validate validates the TokenSupply
func (ts TokenSupply) Validate() error {
	if ts.TotalSupply.IsNil() || ts.TotalSupply.LTE(sdk.ZeroInt()) {
		return fmt.Errorf("total supply must be positive")
	}

	// Validate all allocations are positive
	allocations := []sdk.Int{
		ts.PublicSaleAllocation,
		ts.LiquidityAllocation,
		ts.TeamAllocation,
		ts.DevelopmentAllocation,
		ts.CommunityAllocation,
		ts.DaoTreasuryAllocation,
		ts.InitialBurnAllocation,
	}

	for i, allocation := range allocations {
		if allocation.IsNil() || allocation.LT(sdk.ZeroInt()) {
			return fmt.Errorf("allocation %d must be non-negative", i)
		}
	}

	// Validate total allocations equal total supply
	totalAllocations := sdk.ZeroInt()
	for _, allocation := range allocations {
		totalAllocations = totalAllocations.Add(allocation)
	}

	if !totalAllocations.Equal(ts.TotalSupply) {
		return fmt.Errorf("total allocations (%s) must equal total supply (%s)", totalAllocations, ts.TotalSupply)
	}

	return nil
}

// Validate validates the VestingSchedule
func (vs VestingSchedule) Validate() error {
	if _, err := sdk.AccAddressFromBech32(vs.Recipient); err != nil {
		return fmt.Errorf("invalid recipient address: %w", err)
	}

	if vs.TotalAmount.IsNil() || vs.TotalAmount.LTE(sdk.ZeroInt()) {
		return fmt.Errorf("total amount must be positive")
	}

	if vs.VestingPeriodMonths <= 0 {
		return fmt.Errorf("vesting period must be positive")
	}

	if vs.CliffPeriodMonths < 0 {
		return fmt.Errorf("cliff period cannot be negative")
	}

	if vs.CliffPeriodMonths > vs.VestingPeriodMonths {
		return fmt.Errorf("cliff period cannot be longer than vesting period")
	}

	if vs.StartTime <= 0 {
		return fmt.Errorf("start time must be positive")
	}

	if vs.VestedAmount.IsNil() {
		vs.VestedAmount = sdk.ZeroInt()
	}

	if vs.VestedAmount.IsNegative() {
		return fmt.Errorf("vested amount cannot be negative")
	}

	if vs.VestedAmount.GT(vs.TotalAmount) {
		return fmt.Errorf("vested amount cannot exceed total amount")
	}

	return nil
}

// Validate validates the TokenDistributionEvent
func (tde TokenDistributionEvent) Validate() error {
	if len(tde.EventType) == 0 {
		return fmt.Errorf("event type cannot be empty")
	}

	if _, err := sdk.AccAddressFromBech32(tde.Recipient); err != nil {
		return fmt.Errorf("invalid recipient address: %w", err)
	}

	if err := tde.Amount.Validate(); err != nil {
		return fmt.Errorf("invalid amount: %w", err)
	}

	if tde.Amount.IsNil() || tde.Amount.Amount.LTE(sdk.ZeroInt()) {
		return fmt.Errorf("amount must be positive")
	}

	if tde.Timestamp <= 0 {
		return fmt.Errorf("timestamp must be positive")
	}

	return nil
}

// GetVestableAmount calculates the amount that can be vested at a given time
func (vs VestingSchedule) GetVestableAmount(currentTime time.Time) sdk.Int {
	if currentTime.Unix() < vs.StartTime {
		return sdk.ZeroInt()
	}

	// Check if we're still in the cliff period
	cliffEndTime := time.Unix(vs.StartTime, 0).AddDate(0, int(vs.CliffPeriodMonths), 0)
	if currentTime.Before(cliffEndTime) {
		return sdk.ZeroInt()
	}

	// Calculate vested amount based on elapsed time
	vestingEndTime := time.Unix(vs.StartTime, 0).AddDate(0, int(vs.VestingPeriodMonths), 0)
	if currentTime.After(vestingEndTime) {
		// Fully vested
		return vs.TotalAmount
	}

	// Calculate proportional vesting
	totalVestingDuration := vestingEndTime.Sub(time.Unix(vs.StartTime, 0))
	elapsedDuration := currentTime.Sub(time.Unix(vs.StartTime, 0))
	
	vestingRatio := sdk.NewDecFromBigInt(elapsedDuration.Nanoseconds()).
		Quo(sdk.NewDecFromBigInt(totalVestingDuration.Nanoseconds()))
	
	vestedAmount := vestingRatio.MulInt(vs.TotalAmount).TruncateInt()
	
	return vestedAmount
}

// GetClaimableAmount calculates the amount that can be claimed (vested - already claimed)
func (vs VestingSchedule) GetClaimableAmount(currentTime time.Time) sdk.Int {
	vestableAmount := vs.GetVestableAmount(currentTime)
	claimableAmount := vestableAmount.Sub(vs.VestedAmount)
	
	if claimableAmount.IsNegative() {
		return sdk.ZeroInt()
	}
	
	return claimableAmount
}