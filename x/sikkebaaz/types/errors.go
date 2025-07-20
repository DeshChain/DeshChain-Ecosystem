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
	"cosmossdk.io/errors"
)

// Error codes for Sikkebaaz module
var (
	// Token launch errors
	ErrInvalidTokenName        = errors.Register(ModuleName, 1, "invalid token name")
	ErrInvalidTokenSymbol      = errors.Register(ModuleName, 2, "invalid token symbol")
	ErrInvalidTotalSupply      = errors.Register(ModuleName, 3, "invalid total supply")
	ErrInvalidTargetAmount     = errors.Register(ModuleName, 4, "invalid target amount")
	ErrTokenAlreadyExists      = errors.Register(ModuleName, 5, "token already exists")
	ErrLaunchNotFound          = errors.Register(ModuleName, 6, "launch not found")
	ErrLaunchNotActive         = errors.Register(ModuleName, 7, "launch not active")
	ErrLaunchAlreadyCompleted  = errors.Register(ModuleName, 8, "launch already completed")
	ErrLaunchAlreadyCancelled  = errors.Register(ModuleName, 9, "launch already cancelled")
	
	// Anti-pump protection errors
	ErrInvalidWalletLimit      = errors.Register(ModuleName, 10, "invalid wallet limit percentage")
	ErrInsufficientLiquidityLock = errors.Register(ModuleName, 11, "insufficient liquidity lock period")
	ErrInsufficientLiquidity   = errors.Register(ModuleName, 12, "insufficient liquidity percentage")
	ErrWalletLimitExceeded     = errors.Register(ModuleName, 13, "wallet limit exceeded")
	ErrTradingNotStarted       = errors.Register(ModuleName, 14, "trading not started yet")
	ErrBotDetected             = errors.Register(ModuleName, 15, "bot activity detected")
	ErrMaxTransactionsExceeded = errors.Register(ModuleName, 16, "maximum transactions per block exceeded")
	ErrCooldownPeriodActive    = errors.Register(ModuleName, 17, "cooldown period still active")
	ErrPriceImpactTooHigh      = errors.Register(ModuleName, 18, "price impact too high")
	
	// Participation errors
	ErrInvalidContribution     = errors.Register(ModuleName, 20, "invalid contribution amount")
	ErrContributionTooLow      = errors.Register(ModuleName, 21, "contribution below minimum")
	ErrContributionTooHigh     = errors.Register(ModuleName, 22, "contribution above maximum")
	ErrAlreadyParticipated     = errors.Register(ModuleName, 23, "already participated in launch")
	ErrNotWhitelisted          = errors.Register(ModuleName, 24, "not whitelisted for launch")
	ErrTargetReached           = errors.Register(ModuleName, 25, "target amount already reached")
	ErrLaunchExpired           = errors.Register(ModuleName, 26, "launch period expired")
	ErrInsufficientBalance     = errors.Register(ModuleName, 27, "insufficient balance")
	
	// Cultural and regional errors
	ErrInvalidPincode          = errors.Register(ModuleName, 30, "invalid PIN code format")
	ErrInvalidCulturalQuote    = errors.Register(ModuleName, 31, "invalid cultural quote")
	ErrInsufficientPatriotismScore = errors.Register(ModuleName, 32, "insufficient patriotism score")
	ErrFestivalNotActive       = errors.Register(ModuleName, 33, "festival not currently active")
	ErrRegionNotSupported      = errors.Register(ModuleName, 34, "region not supported")
	
	// Community governance errors
	ErrCommunityVetoActive     = errors.Register(ModuleName, 40, "community veto process active")
	ErrCommunityVetoPassed     = errors.Register(ModuleName, 41, "community veto passed")
	ErrInvalidVoteWeight       = errors.Register(ModuleName, 42, "invalid voting weight")
	ErrVotingPeriodExpired     = errors.Register(ModuleName, 43, "voting period expired")
	ErrAlreadyVoted            = errors.Register(ModuleName, 44, "already voted on this proposal")
	ErrInsufficientVotingPower = errors.Register(ModuleName, 45, "insufficient voting power")
	
	// Liquidity and trading errors
	ErrLiquidityLocked         = errors.Register(ModuleName, 50, "liquidity is locked")
	ErrLiquidityNotLocked      = errors.Register(ModuleName, 51, "liquidity not locked")
	ErrInvalidLockPeriod       = errors.Register(ModuleName, 52, "invalid lock period")
	ErrLockNotExpired          = errors.Register(ModuleName, 53, "lock period not expired")
	ErrAlreadyWithdrawn        = errors.Register(ModuleName, 54, "liquidity already withdrawn")
	ErrInvalidTradingPair      = errors.Register(ModuleName, 55, "invalid trading pair")
	ErrNoLiquidity             = errors.Register(ModuleName, 56, "no liquidity available")
	
	// Creator and reward errors
	ErrNotCreator              = errors.Register(ModuleName, 60, "not the token creator")
	ErrCreatorRewardNotFound   = errors.Register(ModuleName, 61, "creator reward not found")
	ErrNoRewardsToClaim        = errors.Register(ModuleName, 62, "no rewards to claim")
	ErrRewardAlreadyClaimed    = errors.Register(ModuleName, 63, "reward already claimed")
	ErrCreatorRewardInactive   = errors.Register(ModuleName, 64, "creator reward inactive")
	
	// Security and audit errors
	ErrSecurityAuditRequired   = errors.Register(ModuleName, 70, "security audit required")
	ErrSecurityAuditFailed     = errors.Register(ModuleName, 71, "security audit failed")
	ErrHighRiskToken           = errors.Register(ModuleName, 72, "high risk token")
	ErrEmergencyStop           = errors.Register(ModuleName, 73, "emergency stop activated")
	ErrTokenBlacklisted        = errors.Register(ModuleName, 74, "token blacklisted")
	ErrUnauthorizedAccess      = errors.Register(ModuleName, 75, "unauthorized access")
	
	// Fee and payment errors
	ErrInsufficientFees        = errors.Register(ModuleName, 80, "insufficient fees")
	ErrInvalidFeeAmount        = errors.Register(ModuleName, 81, "invalid fee amount")
	ErrFeePaymentFailed        = errors.Register(ModuleName, 82, "fee payment failed")
	ErrRefundFailed            = errors.Register(ModuleName, 83, "refund failed")
	
	// Configuration errors
	ErrInvalidConfiguration    = errors.Register(ModuleName, 90, "invalid configuration")
	ErrParameterOutOfRange     = errors.Register(ModuleName, 91, "parameter out of valid range")
	ErrIncompatibleSettings    = errors.Register(ModuleName, 92, "incompatible settings")
	ErrFeatureNotEnabled       = errors.Register(ModuleName, 93, "feature not enabled")
	ErrModuleDisabled          = errors.Register(ModuleName, 94, "module disabled")
	
	// General errors
	ErrInternalError           = errors.Register(ModuleName, 99, "internal error")
)