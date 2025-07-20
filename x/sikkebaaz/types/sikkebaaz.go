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

// TokenLaunch represents a memecoin launch on Sikkebaaz platform
type TokenLaunch struct {
	LaunchID        string    `json:"launch_id" yaml:"launch_id"`
	Creator         string    `json:"creator" yaml:"creator"`
	CreatorPincode  string    `json:"creator_pincode" yaml:"creator_pincode"`
	
	// Token details
	TokenName       string    `json:"token_name" yaml:"token_name"`
	TokenSymbol     string    `json:"token_symbol" yaml:"token_symbol"`
	TokenDescription string   `json:"token_description" yaml:"token_description"`
	TotalSupply     sdk.Int   `json:"total_supply" yaml:"total_supply"`
	Decimals        uint32    `json:"decimals" yaml:"decimals"`
	
	// Launch configuration
	LaunchType      string    `json:"launch_type" yaml:"launch_type"`
	TargetAmount    sdk.Int   `json:"target_amount" yaml:"target_amount"`
	RaisedAmount    sdk.Int   `json:"raised_amount" yaml:"raised_amount"`
	MinContribution sdk.Int   `json:"min_contribution" yaml:"min_contribution"`
	MaxContribution sdk.Int   `json:"max_contribution" yaml:"max_contribution"`
	
	// Timing
	StartTime       time.Time `json:"start_time" yaml:"start_time"`
	EndTime         time.Time `json:"end_time" yaml:"end_time"`
	TradingDelay    int64     `json:"trading_delay" yaml:"trading_delay"` // Seconds after launch
	
	// Anti-pump protection
	AntiPumpConfig  AntiPumpConfig `json:"anti_pump_config" yaml:"anti_pump_config"`
	
	// Cultural features
	CulturalQuote   string    `json:"cultural_quote" yaml:"cultural_quote"`
	FestivalBonus   bool      `json:"festival_bonus" yaml:"festival_bonus"`
	PatriotismScore int64     `json:"patriotism_score" yaml:"patriotism_score"`
	
	// Financial
	LaunchFee       sdk.Int   `json:"launch_fee" yaml:"launch_fee"`
	CharityAllocation sdk.Int `json:"charity_allocation" yaml:"charity_allocation"`
	CreatorAllocation sdk.Int `json:"creator_allocation" yaml:"creator_allocation"`
	
	// Status and metadata
	Status          string              `json:"status" yaml:"status"`
	ParticipantCount uint64             `json:"participant_count" yaml:"participant_count"`
	Whitelist       []string            `json:"whitelist" yaml:"whitelist"`
	Metadata        map[string]string   `json:"metadata" yaml:"metadata"`
	
	// Timestamps
	CreatedAt       time.Time `json:"created_at" yaml:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" yaml:"updated_at"`
	CompletedAt     *time.Time `json:"completed_at,omitempty" yaml:"completed_at,omitempty"`
}

// AntiPumpConfig defines anti-pump and dump protection settings
type AntiPumpConfig struct {
	// Wallet limits
	MaxWalletPercent24h  uint32 `json:"max_wallet_percent_24h" yaml:"max_wallet_percent_24h"`   // 5%
	MaxWalletPercentAfter uint32 `json:"max_wallet_percent_after" yaml:"max_wallet_percent_after"` // 10%
	
	// Trading restrictions
	TradingDelay         int64  `json:"trading_delay" yaml:"trading_delay"`                     // Seconds
	MinBlocksBetweenTx   uint32 `json:"min_blocks_between_tx" yaml:"min_blocks_between_tx"`     // Anti-bot
	MaxTxPerBlock        uint32 `json:"max_tx_per_block" yaml:"max_tx_per_block"`               // Rate limiting
	
	// Liquidity protection
	LiquidityLockDays    uint32 `json:"liquidity_lock_days" yaml:"liquidity_lock_days"`         // 365 minimum
	MinLiquidityPercent  uint32 `json:"min_liquidity_percent" yaml:"min_liquidity_percent"`     // 80% minimum
	
	// Price protection
	MaxPriceImpact       sdk.Dec `json:"max_price_impact" yaml:"max_price_impact"`               // 10% max
	MaxSlippage          sdk.Dec `json:"max_slippage" yaml:"max_slippage"`                       // 5% max
	
	// Time-based restrictions
	CooldownPeriod       int64   `json:"cooldown_period" yaml:"cooldown_period"`                 // Between sells
	GradualReleaseEnabled bool   `json:"gradual_release_enabled" yaml:"gradual_release_enabled"` // Vesting
}

// LiquidityLock represents locked liquidity for a token
type LiquidityLock struct {
	TokenAddress    string    `json:"token_address" yaml:"token_address"`
	LockOwner       string    `json:"lock_owner" yaml:"lock_owner"`
	LPTokenAddress  string    `json:"lp_token_address" yaml:"lp_token_address"`
	LockedAmount    sdk.Int   `json:"locked_amount" yaml:"locked_amount"`
	LockDate        time.Time `json:"lock_date" yaml:"lock_date"`
	UnlockDate      time.Time `json:"unlock_date" yaml:"unlock_date"`
	IsWithdrawn     bool      `json:"is_withdrawn" yaml:"is_withdrawn"`
	WithdrawnAt     *time.Time `json:"withdrawn_at,omitempty" yaml:"withdrawn_at,omitempty"`
}

// CreatorReward tracks rewards for token creators
type CreatorReward struct {
	Creator         string    `json:"creator" yaml:"creator"`
	TokenAddress    string    `json:"token_address" yaml:"token_address"`
	RewardRate      sdk.Dec   `json:"reward_rate" yaml:"reward_rate"`         // 2% of trading volume
	AccumulatedReward sdk.Int `json:"accumulated_reward" yaml:"accumulated_reward"`
	LastClaimedAt   time.Time `json:"last_claimed_at" yaml:"last_claimed_at"`
	TotalClaimed    sdk.Int   `json:"total_claimed" yaml:"total_claimed"`
	IsActive        bool      `json:"is_active" yaml:"is_active"`
}

// CommunityVeto represents community voting against a launch
type CommunityVeto struct {
	LaunchID        string            `json:"launch_id" yaml:"launch_id"`
	InitiatedBy     string            `json:"initiated_by" yaml:"initiated_by"`
	VoteStartTime   time.Time         `json:"vote_start_time" yaml:"vote_start_time"`
	VoteEndTime     time.Time         `json:"vote_end_time" yaml:"vote_end_time"`
	Votes           map[string]bool   `json:"votes" yaml:"votes"`          // voter -> yes/no
	VotingPower     map[string]sdk.Int `json:"voting_power" yaml:"voting_power"` // voter -> power
	TotalVotingPower sdk.Int          `json:"total_voting_power" yaml:"total_voting_power"`
	VetoThreshold   sdk.Dec           `json:"veto_threshold" yaml:"veto_threshold"` // 70%
	Status          string            `json:"status" yaml:"status"`        // active, passed, failed
	Reason          string            `json:"reason" yaml:"reason"`
}

// WalletLimits tracks wallet limits for anti-pump protection
type WalletLimits struct {
	TokenAddress    string    `json:"token_address" yaml:"token_address"`
	WalletAddress   string    `json:"wallet_address" yaml:"wallet_address"`
	MaxAmount       sdk.Int   `json:"max_amount" yaml:"max_amount"`
	CurrentAmount   sdk.Int   `json:"current_amount" yaml:"current_amount"`
	LastTxTime      time.Time `json:"last_tx_time" yaml:"last_tx_time"`
	LastTxBlock     int64     `json:"last_tx_block" yaml:"last_tx_block"`
	ViolationCount  uint32    `json:"violation_count" yaml:"violation_count"`
	IsRestricted    bool      `json:"is_restricted" yaml:"is_restricted"`
}

// TradingMetrics tracks trading statistics for tokens
type TradingMetrics struct {
	TokenAddress      string    `json:"token_address" yaml:"token_address"`
	TotalVolume       sdk.Int   `json:"total_volume" yaml:"total_volume"`
	DailyVolume       sdk.Int   `json:"daily_volume" yaml:"daily_volume"`
	TotalTrades       uint64    `json:"total_trades" yaml:"total_trades"`
	DailyTrades       uint64    `json:"daily_trades" yaml:"daily_trades"`
	UniqueTraders     uint64    `json:"unique_traders" yaml:"unique_traders"`
	CurrentPrice      sdk.Dec   `json:"current_price" yaml:"current_price"`
	PriceChange24h    sdk.Dec   `json:"price_change_24h" yaml:"price_change_24h"`
	MarketCap         sdk.Int   `json:"market_cap" yaml:"market_cap"`
	Liquidity         sdk.Int   `json:"liquidity" yaml:"liquidity"`
	LastUpdated       time.Time `json:"last_updated" yaml:"last_updated"`
}

// LaunchParticipation represents user participation in a launch
type LaunchParticipation struct {
	LaunchID        string    `json:"launch_id" yaml:"launch_id"`
	Participant     string    `json:"participant" yaml:"participant"`
	ContributedAmount sdk.Int `json:"contributed_amount" yaml:"contributed_amount"`
	TokensAllocated sdk.Int   `json:"tokens_allocated" yaml:"tokens_allocated"`
	TokensClaimed   sdk.Int   `json:"tokens_claimed" yaml:"tokens_claimed"`
	ParticipatedAt  time.Time `json:"participated_at" yaml:"participated_at"`
	ClaimedAt       *time.Time `json:"claimed_at,omitempty" yaml:"claimed_at,omitempty"`
	IsRefunded      bool      `json:"is_refunded" yaml:"is_refunded"`
	RefundedAt      *time.Time `json:"refunded_at,omitempty" yaml:"refunded_at,omitempty"`
}

// SecurityAudit represents security audit results for a token
type SecurityAudit struct {
	TokenAddress    string              `json:"token_address" yaml:"token_address"`
	AuditorAddress  string              `json:"auditor_address" yaml:"auditor_address"`
	AuditDate       time.Time           `json:"audit_date" yaml:"audit_date"`
	SecurityScore   uint32              `json:"security_score" yaml:"security_score"` // 0-100
	RiskLevel       string              `json:"risk_level" yaml:"risk_level"`         // low, medium, high
	Findings        []string            `json:"findings" yaml:"findings"`
	Recommendations []string            `json:"recommendations" yaml:"recommendations"`
	CertificateHash string              `json:"certificate_hash" yaml:"certificate_hash"`
	IsApproved      bool                `json:"is_approved" yaml:"is_approved"`
	Metadata        map[string]string   `json:"metadata" yaml:"metadata"`
}

// FestivalBonus represents festival-based bonuses
type FestivalBonus struct {
	LaunchID        string    `json:"launch_id" yaml:"launch_id"`
	FestivalName    string    `json:"festival_name" yaml:"festival_name"`
	BonusRate       sdk.Dec   `json:"bonus_rate" yaml:"bonus_rate"`     // 10%
	BonusAmount     sdk.Int   `json:"bonus_amount" yaml:"bonus_amount"`
	AppliedAt       time.Time `json:"applied_at" yaml:"applied_at"`
	CulturalQuote   string    `json:"cultural_quote" yaml:"cultural_quote"`
}

// EmergencyControl represents emergency controls for problematic tokens
type EmergencyControl struct {
	TokenAddress    string              `json:"token_address" yaml:"token_address"`
	ControlType     string              `json:"control_type" yaml:"control_type"`     // pause, blacklist, etc.
	InitiatedBy     string              `json:"initiated_by" yaml:"initiated_by"`
	Reason          string              `json:"reason" yaml:"reason"`
	ActivatedAt     time.Time           `json:"activated_at" yaml:"activated_at"`
	ExpiresAt       *time.Time          `json:"expires_at,omitempty" yaml:"expires_at,omitempty"`
	IsActive        bool                `json:"is_active" yaml:"is_active"`
	Metadata        map[string]string   `json:"metadata" yaml:"metadata"`
}

// Validation functions

// ValidateTokenLaunch validates a token launch configuration
func ValidateTokenLaunch(launch TokenLaunch) error {
	if launch.TokenName == "" {
		return ErrInvalidTokenName
	}
	if launch.TokenSymbol == "" {
		return ErrInvalidTokenSymbol
	}
	if launch.TotalSupply.IsZero() || launch.TotalSupply.IsNegative() {
		return ErrInvalidTotalSupply
	}
	if launch.TargetAmount.IsZero() || launch.TargetAmount.IsNegative() {
		return ErrInvalidTargetAmount
	}
	return nil
}

// ValidateAntiPumpConfig validates anti-pump configuration
func ValidateAntiPumpConfig(config AntiPumpConfig) error {
	if config.MaxWalletPercent24h > 1000 { // 10% max
		return ErrInvalidWalletLimit
	}
	if config.LiquidityLockDays < 365 { // 1 year minimum
		return ErrInsufficientLiquidityLock
	}
	if config.MinLiquidityPercent < 8000 { // 80% minimum
		return ErrInsufficientLiquidity
	}
	return nil
}

// CalculateLaunchFee calculates the total launch fee
func CalculateLaunchFee(targetAmount sdk.Int) sdk.Int {
	baseFee, _ := sdk.NewIntFromString(BaseLaunchFee)
	variableFee := targetAmount.MulRaw(5).QuoRaw(100) // 5%
	return baseFee.Add(variableFee)
}

// CalculateCharityAllocation calculates charity allocation (1% of raised)
func CalculateCharityAllocation(raisedAmount sdk.Int) sdk.Int {
	return raisedAmount.MulRaw(1).QuoRaw(100) // 1%
}

// IsLaunchActive checks if a launch is currently active
func (l TokenLaunch) IsLaunchActive(currentTime time.Time) bool {
	return l.Status == LaunchStatusActive &&
		currentTime.After(l.StartTime) &&
		currentTime.Before(l.EndTime) &&
		l.RaisedAmount.LT(l.TargetAmount)
}

// CalculateTokenAllocation calculates token allocation for a contribution
func (l TokenLaunch) CalculateTokenAllocation(contribution sdk.Int) sdk.Int {
	if l.TargetAmount.IsZero() {
		return sdk.ZeroInt()
	}
	// Simple allocation: (contribution / target) * total_supply
	allocation := contribution.Mul(l.TotalSupply).Quo(l.TargetAmount)
	return allocation
}

// GetCurrentWalletLimit returns current wallet limit based on time since launch
func (l TokenLaunch) GetCurrentWalletLimit(currentTime time.Time) uint32 {
	if l.CompletedAt == nil {
		return l.AntiPumpConfig.MaxWalletPercent24h
	}
	
	hoursSinceLaunch := currentTime.Sub(*l.CompletedAt).Hours()
	if hoursSinceLaunch < 24 {
		return l.AntiPumpConfig.MaxWalletPercent24h
	}
	return l.AntiPumpConfig.MaxWalletPercentAfter
}