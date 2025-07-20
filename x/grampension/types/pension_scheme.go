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
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// PensionScheme represents a pension scheme configuration
type PensionScheme struct {
	SchemeID                string    `json:"scheme_id"`
	SchemeName              string    `json:"scheme_name"`
	Description             string    `json:"description"`
	MonthlyContribution     sdk.Coin  `json:"monthly_contribution"`
	ContributionPeriod      uint32    `json:"contribution_period"`      // in months
	MaturityBonus           sdk.Dec   `json:"maturity_bonus"`           // percentage (e.g., 0.50 for 50%)
	MinAge                  uint32    `json:"min_age"`
	MaxAge                  uint32    `json:"max_age"`
	MaxParticipants         uint64    `json:"max_participants"`
	CurrentParticipants     uint64    `json:"current_participants"`
	Status                  string    `json:"status"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
	GracePeriodDays         uint32    `json:"grace_period_days"`
	EarlyWithdrawalPenalty  sdk.Dec   `json:"early_withdrawal_penalty"`
	LatePaymentPenalty      sdk.Dec   `json:"late_payment_penalty"`
	ReferralRewardPercent   sdk.Dec   `json:"referral_reward_percent"`
	OnTimeBonusPercent      sdk.Dec   `json:"on_time_bonus_percent"`
	KYCRequired             bool      `json:"kyc_required"`
	AutoRenewalEnabled      bool      `json:"auto_renewal_enabled"`
	CulturalIntegration     bool      `json:"cultural_integration"`
	LiquidityProvision      bool      `json:"liquidity_provision"`      // Enable unified pool integration
	LiquidityUtilization    sdk.Dec   `json:"liquidity_utilization"`    // Percentage used for liquidity (e.g., 0.80 for 80%)
}

// PensionParticipant represents a participant in a pension scheme
type PensionParticipant struct {
	ParticipantID       string         `json:"participant_id"`
	SchemeID            string         `json:"scheme_id"`
	Address             sdk.AccAddress `json:"address"`
	Name                string         `json:"name"`
	Age                 uint32         `json:"age"`
	EnrollmentDate      time.Time      `json:"enrollment_date"`
	MaturityDate        time.Time      `json:"maturity_date"`
	Status              string         `json:"status"`
	TotalContributed    sdk.Coin       `json:"total_contributed"`
	MissedPayments      uint32         `json:"missed_payments"`
	OnTimePayments      uint32         `json:"on_time_payments"`
	BonusEarned         sdk.Coin       `json:"bonus_earned"`
	PenaltyIncurred     sdk.Coin       `json:"penalty_incurred"`
	ReferrerAddress     sdk.AccAddress `json:"referrer_address,omitempty"`
	ReferralRewards     sdk.Coin       `json:"referral_rewards"`
	KYCStatus           string         `json:"kyc_status"`
	PatriotismScore     uint64         `json:"patriotism_score"`
	LastContribution    time.Time      `json:"last_contribution"`
	VillagePostalCode   string         `json:"village_postal_code"`
	LiquidityContributed sdk.Coin      `json:"liquidity_contributed"`
}

// PensionContribution represents a single contribution
type PensionContribution struct {
	ContributionID   string         `json:"contribution_id"`
	ParticipantID    string         `json:"participant_id"`
	SchemeID         string         `json:"scheme_id"`
	Address          sdk.AccAddress `json:"address"`
	Amount           sdk.Coin       `json:"amount"`
	Month            uint32         `json:"month"` // which month of the scheme (1-12)
	ContributionDate time.Time      `json:"contribution_date"`
	TransactionHash  string         `json:"transaction_hash"`
	Status           string         `json:"status"`
	BonusApplied     sdk.Coin       `json:"bonus_applied,omitempty"`
	PenaltyApplied   sdk.Coin       `json:"penalty_applied,omitempty"`
	CulturalQuote    string         `json:"cultural_quote,omitempty"`
	LiquidityProvided bool          `json:"liquidity_provided"`
}

// PensionMaturity represents a matured pension account
type PensionMaturity struct {
	MaturityID         string         `json:"maturity_id"`
	ParticipantID      string         `json:"participant_id"`
	SchemeID           string         `json:"scheme_id"`
	Address            sdk.AccAddress `json:"address"`
	MaturityDate       time.Time      `json:"maturity_date"`
	TotalContributed   sdk.Coin       `json:"total_contributed"`
	MaturityBonus      sdk.Coin       `json:"maturity_bonus"`
	TotalPayout        sdk.Coin       `json:"total_payout"`
	Status             string         `json:"status"`
	ProcessedDate      time.Time      `json:"processed_date,omitempty"`
	TransactionHash    string         `json:"transaction_hash,omitempty"`
	LiquidityReturned  sdk.Coin       `json:"liquidity_returned"`
	LiquidityEarnings  sdk.Coin       `json:"liquidity_earnings"`
}

// PensionWithdrawal represents an early withdrawal request
type PensionWithdrawal struct {
	WithdrawalID    string         `json:"withdrawal_id"`
	ParticipantID   string         `json:"participant_id"`
	SchemeID        string         `json:"scheme_id"`
	Address         sdk.AccAddress `json:"address"`
	RequestedAmount sdk.Coin       `json:"requested_amount"`
	PenaltyAmount   sdk.Coin       `json:"penalty_amount"`
	NetAmount       sdk.Coin       `json:"net_amount"`
	Reason          string         `json:"reason"`
	RequestDate     time.Time      `json:"request_date"`
	ProcessedDate   time.Time      `json:"processed_date,omitempty"`
	Status          string         `json:"status"`
	TransactionHash string         `json:"transaction_hash,omitempty"`
}

// PensionStatistics represents scheme-wide statistics
type PensionStatistics struct {
	SchemeID              string    `json:"scheme_id"`
	TotalParticipants     uint64    `json:"total_participants"`
	ActiveParticipants    uint64    `json:"active_participants"`
	MaturedParticipants   uint64    `json:"matured_participants"`
	WithdrawnParticipants uint64    `json:"withdrawn_participants"`
	DefaultedParticipants uint64    `json:"defaulted_participants"`
	TotalContributed      sdk.Coin  `json:"total_contributed"`
	TotalMaturityPaid     sdk.Coin  `json:"total_maturity_paid"`
	TotalBonusPaid        sdk.Coin  `json:"total_bonus_paid"`
	TotalPenaltyCollected sdk.Coin  `json:"total_penalty_collected"`
	TotalReferralRewards  sdk.Coin  `json:"total_referral_rewards"`
	AveragePatriotismScore uint64   `json:"average_patriotism_score"`
	CompletionRate        sdk.Dec   `json:"completion_rate"`
	DefaultRate           sdk.Dec   `json:"default_rate"`
	SustainabilityScore   sdk.Dec   `json:"sustainability_score"`
	LastUpdated           time.Time `json:"last_updated"`
	TotalLiquidityProvided sdk.Coin `json:"total_liquidity_provided"`
	LiquidityEarnings     sdk.Coin  `json:"liquidity_earnings"`
}

// Validate performs basic validation on PensionScheme
func (ps PensionScheme) Validate() error {
	if ps.SchemeID == "" {
		return ErrInvalidSchemeID
	}
	if ps.SchemeName == "" {
		return ErrInvalidSchemeName
	}
	if !ps.MonthlyContribution.IsPositive() {
		return ErrInvalidContribution
	}
	if ps.ContributionPeriod == 0 || ps.ContributionPeriod > 120 {
		return ErrInvalidContributionPeriod
	}
	if ps.MaturityBonus.IsNegative() || ps.MaturityBonus.GT(sdk.NewDec(1)) {
		return ErrInvalidMaturityBonus
	}
	if ps.MinAge == 0 || ps.MaxAge == 0 || ps.MinAge > ps.MaxAge {
		return ErrInvalidAgeRange
	}
	if ps.LiquidityUtilization.IsNegative() || ps.LiquidityUtilization.GT(sdk.NewDec(1)) {
		return ErrInvalidLiquidityUtilization
	}
	return nil
}

// Validate performs basic validation on PensionParticipant
func (pp PensionParticipant) Validate() error {
	if pp.ParticipantID == "" {
		return ErrInvalidParticipantID
	}
	if pp.SchemeID == "" {
		return ErrInvalidSchemeID
	}
	if pp.Address.Empty() {
		return ErrInvalidAddress
	}
	if pp.Age == 0 {
		return ErrInvalidAge
	}
	return nil
}

// CalculateMaturityAmount calculates the total maturity amount including bonus
func (ps PensionScheme) CalculateMaturityAmount(totalContributed sdk.Coin) sdk.Coin {
	bonusAmount := totalContributed.Amount.ToDec().Mul(ps.MaturityBonus).TruncateInt()
	totalAmount := totalContributed.Amount.Add(bonusAmount)
	return sdk.NewCoin(totalContributed.Denom, totalAmount)
}

// CalculateEarlyWithdrawalAmount calculates amount after penalty
func (ps PensionScheme) CalculateEarlyWithdrawalAmount(totalContributed sdk.Coin) sdk.Coin {
	penaltyAmount := totalContributed.Amount.ToDec().Mul(ps.EarlyWithdrawalPenalty).TruncateInt()
	netAmount := totalContributed.Amount.Sub(penaltyAmount)
	if netAmount.IsNegative() {
		netAmount = sdk.ZeroInt()
	}
	return sdk.NewCoin(totalContributed.Denom, netAmount)
}

// IsEligible checks if a person is eligible for the scheme
func (ps PensionScheme) IsEligible(age uint32) bool {
	return age >= ps.MinAge && age <= ps.MaxAge && ps.CurrentParticipants < ps.MaxParticipants && ps.Status == StatusActive
}

// GetLiquidityAmount calculates the amount available for liquidity provision
func (ps PensionScheme) GetLiquidityAmount(contribution sdk.Coin) sdk.Coin {
	if !ps.LiquidityProvision {
		return sdk.NewCoin(contribution.Denom, sdk.ZeroInt())
	}
	liquidityAmount := contribution.Amount.ToDec().Mul(ps.LiquidityUtilization).TruncateInt()
	return sdk.NewCoin(contribution.Denom, liquidityAmount)
}