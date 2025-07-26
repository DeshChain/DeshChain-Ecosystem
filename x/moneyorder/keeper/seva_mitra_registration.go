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

// IMPORTANT: Seva Mitras are independent community members providing voluntary
// peer-to-peer services. They are NOT agents, employees, or representatives of
// DeshChain. This system facilitates voluntary community cooperation without
// creating any agency relationship.

package keeper

import (
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/moneyorder/types"
)

// RegisterSevaMitra registers a new cash-in/out sevaMitra
func (k Keeper) RegisterSevaMitra(ctx sdk.Context, msg *types.MsgRegisterSevaMitra) (*types.SevaMitra, error) {
	// Validate postal code
	if !k.ValidateIndianPincode(msg.PostalCode) {
		return nil, types.ErrInvalidPostalCode
	}
	
	// Check if already registered
	if k.IsSevaMitraRegistered(ctx, msg.Address) {
		return nil, types.ErrSevaMitraAlreadyRegistered
	}
	
	// Validate business registration
	if !k.validateBusinessRegistration(msg.BusinessName, msg.RegistrationNumber) {
		return nil, types.ErrInvalidBusinessRegistration
	}
	
	// Validate services
	if len(msg.Services) == 0 {
		return nil, types.ErrNoServicesProvided
	}
	
	// Create sevaMitra
	sevaMitra := &types.SevaMitra{
		MitraId:              k.generateSevaMitraID(ctx),
		Address:              msg.Address,
		BusinessName:         msg.BusinessName,
		RegistrationNumber:   msg.RegistrationNumber,
		PostalCode:          msg.PostalCode,
		FullAddress:         msg.FullAddress,
		District:            k.getDistrictFromPincode(msg.PostalCode),
		State:               k.GetStateName(msg.PostalCode),
		Latitude:            msg.Latitude,
		Longitude:           msg.Longitude,
		Phone:               msg.Phone,
		Email:               msg.Email,
		Languages:           msg.Languages,
		OperatingHours:      msg.OperatingHours,
		Services:            msg.Services,
		DailyLimit:          msg.DailyLimit,
		PerTransactionLimit: msg.PerTransactionLimit,
		Status:              types.SevaMitraStatus_MITRA_STATUS_PENDING,
		KycVerified:         false,
		Stats:               k.initializeSevaMitraStats(),
		CommissionRate:      k.getDefaultCommissionRate(msg.Services),
		CreatedAt:           ctx.BlockTime(),
		SecurityDeposit:     msg.SecurityDeposit,
		BankDetails:         msg.BankDetails,
	}
	
	// Require security deposit
	depositor, _ := sdk.AccAddressFromBech32(msg.Address)
	err := k.collectSecurityDeposit(ctx, depositor, msg.SecurityDeposit)
	if err != nil {
		return nil, err
	}
	
	// Store sevaMitra
	k.SetSevaMitra(ctx, sevaMitra)
	
	// Add to district index
	k.AddSevaMitraToDistrictIndex(ctx, sevaMitra.District, sevaMitra.MitraId)
	
	// Schedule KYC verification
	k.ScheduleKYCVerification(ctx, sevaMitra.MitraId)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSevaMitraRegistered,
			sdk.NewAttribute("sevaMitra_id", sevaMitra.MitraId),
			sdk.NewAttribute("business_name", sevaMitra.BusinessName),
			sdk.NewAttribute("postal_code", sevaMitra.PostalCode),
			sdk.NewAttribute("services", fmt.Sprintf("%v", sevaMitra.Services)),
		),
	)
	
	return sevaMitra, nil
}

// CompleteSevaMitraKYC completes KYC verification for an sevaMitra
func (k Keeper) CompleteSevaMitraKYC(ctx sdk.Context, sevaMitraID string, verifier sdk.AccAddress, kycData *types.KYCData) error {
	sevaMitra, found := k.GetSevaMitra(ctx, sevaMitraID)
	if !found {
		return types.ErrSevaMitraNotFound
	}
	
	// Check if verifier is authorized
	if !k.isAuthorizedKYCVerifier(ctx, verifier) {
		return types.ErrUnauthorizedKYCVerifier
	}
	
	// Validate KYC data
	if !k.validateKYCData(kycData) {
		return types.ErrInvalidKYCData
	}
	
	// Update sevaMitra status
	sevaMitra.KycVerified = true
	sevaMitra.KycData = kycData
	sevaMitra.KycVerifiedAt = ctx.BlockTime()
	sevaMitra.KycVerifiedBy = verifier.String()
	sevaMitra.Status = types.SevaMitraStatus_MITRA_STATUS_ACTIVE
	
	// Store updated sevaMitra
	k.SetSevaMitra(ctx, sevaMitra)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSevaMitraKYCCompleted,
			sdk.NewAttribute("sevaMitra_id", sevaMitraID),
			sdk.NewAttribute("verified_by", verifier.String()),
		),
	)
	
	return nil
}

// UpdateSevaMitraStatus updates sevaMitra's operational status
func (k Keeper) UpdateSevaMitraStatus(ctx sdk.Context, sevaMitraID string, newStatus types.SevaMitraStatus, reason string) error {
	sevaMitra, found := k.GetSevaMitra(ctx, sevaMitraID)
	if !found {
		return types.ErrSevaMitraNotFound
	}
	
	// Validate status transition
	if !k.isValidStatusTransition(sevaMitra.Status, newStatus) {
		return types.ErrInvalidStatusTransition
	}
	
	// Update status
	sevaMitra.Status = newStatus
	sevaMitra.StatusUpdatedAt = ctx.BlockTime()
	sevaMitra.StatusReason = reason
	
	// Handle suspension
	if newStatus == types.SevaMitraStatus_MITRA_STATUS_SUSPENDED {
		sevaMitra.SuspendedUntil = ctx.BlockTime().Add(30 * 24 * time.Hour) // 30 days
	}
	
	// Store updated sevaMitra
	k.SetSevaMitra(ctx, sevaMitra)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSevaMitraStatusUpdated,
			sdk.NewAttribute("sevaMitra_id", sevaMitraID),
			sdk.NewAttribute("new_status", newStatus.String()),
			sdk.NewAttribute("reason", reason),
		),
	)
	
	return nil
}

// UpdateSevaMitraLimits updates transaction limits for an sevaMitra
func (k Keeper) UpdateSevaMitraLimits(ctx sdk.Context, sevaMitraID string, dailyLimit, perTransactionLimit sdk.Coin) error {
	sevaMitra, found := k.GetSevaMitra(ctx, sevaMitraID)
	if !found {
		return types.ErrSevaMitraNotFound
	}
	
	// Validate limits
	if !k.validateSevaMitraLimits(dailyLimit, perTransactionLimit, sevaMitra.Stats) {
		return types.ErrInvalidLimits
	}
	
	// Update limits
	sevaMitra.DailyLimit = dailyLimit
	sevaMitra.PerTransactionLimit = perTransactionLimit
	
	// Store updated sevaMitra
	k.SetSevaMitra(ctx, sevaMitra)
	
	return nil
}

// RecordSevaMitraTransaction records a transaction for sevaMitra stats
func (k Keeper) RecordSevaMitraTransaction(ctx sdk.Context, sevaMitraID string, txType types.SevaMitraService, amount sdk.Coin, success bool) error {
	sevaMitra, found := k.GetSevaMitra(ctx, sevaMitraID)
	if !found {
		return types.ErrSevaMitraNotFound
	}
	
	// Update stats
	sevaMitra.Stats.TotalTransactions++
	if success {
		sevaMitra.Stats.SuccessfulTransactions++
	}
	sevaMitra.Stats.TotalVolume = sevaMitra.Stats.TotalVolume.Add(amount)
	sevaMitra.Stats.LastActive = ctx.BlockTime()
	
	// Update daily stats
	k.updateSevaMitraDailyStats(ctx, sevaMitra, amount)
	
	// Calculate new commission rate based on performance
	sevaMitra.CommissionRate = k.calculateSevaMitraCommissionRate(sevaMitra.Stats)
	
	// Store updated sevaMitra
	k.SetSevaMitra(ctx, sevaMitra)
	
	return nil
}

// RateSevaMitra allows users to rate an sevaMitra
func (k Keeper) RateSevaMitra(ctx sdk.Context, sevaMitraID string, rater sdk.AccAddress, rating int32, comment string) error {
	sevaMitra, found := k.GetSevaMitra(ctx, sevaMitraID)
	if !found {
		return types.ErrSevaMitraNotFound
	}
	
	// Validate rating
	if rating < 1 || rating > 5 {
		return types.ErrInvalidRating
	}
	
	// Check if user has transacted with sevaMitra
	if !k.hasTransactedWithSevaMitra(ctx, rater, sevaMitraID) {
		return types.ErrNoTransactionHistory
	}
	
	// Update rating
	totalRating := sevaMitra.Stats.AverageRating * float64(sevaMitra.Stats.TotalRatings)
	totalRating += float64(rating)
	sevaMitra.Stats.TotalRatings++
	sevaMitra.Stats.AverageRating = totalRating / float64(sevaMitra.Stats.TotalRatings)
	
	// Store rating details
	k.StoreSevaMitraRating(ctx, sevaMitraID, rater, rating, comment)
	
	// Update trust score
	k.updateSevaMitraTrustScore(ctx, sevaMitra)
	
	// Store updated sevaMitra
	k.SetSevaMitra(ctx, sevaMitra)
	
	return nil
}

// Helper functions

func (k Keeper) generateSevaMitraID(ctx sdk.Context) string {
	return fmt.Sprintf("SM-%d-%s", ctx.BlockHeight(), k.generateRandomString(6))
}

func (k Keeper) validateBusinessRegistration(businessName, registrationNumber string) bool {
	// Basic validation
	if len(businessName) < 3 || len(registrationNumber) < 5 {
		return false
	}
	
	// Check for blacklisted names
	blacklist := []string{"test", "demo", "sample"}
	lowerName := strings.ToLower(businessName)
	for _, blocked := range blacklist {
		if strings.Contains(lowerName, blocked) {
			return false
		}
	}
	
	return true
}

func (k Keeper) initializeSevaMitraStats() *types.SevaMitraStats {
	return &types.SevaMitraStats{
		TotalTransactions:      0,
		SuccessfulTransactions: 0,
		TotalVolume:           sdk.NewCoin(types.DefaultDenom, sdk.ZeroInt()),
		AverageRating:         0,
		TotalRatings:          0,
		DisputesResolved:      0,
		DisputesLost:          0,
		LastActive:            time.Time{},
	}
}

func (k Keeper) getDefaultCommissionRate(services []types.SevaMitraService) string {
	// Base rate depends on services offered
	hasRemittance := false
	hasCashOut := false
	
	for _, service := range services {
		if service == types.SevaMitraService_SERVICE_REMITTANCE {
			hasRemittance = true
		}
		if service == types.SevaMitraService_SERVICE_CASH_OUT {
			hasCashOut = true
		}
	}
	
	if hasRemittance && hasCashOut {
		return "2.0" // Lower rate for full service
	} else if hasRemittance || hasCashOut {
		return "2.5" // Medium rate
	}
	
	return "3.0" // Higher rate for basic services
}

func (k Keeper) collectSecurityDeposit(ctx sdk.Context, depositor sdk.AccAddress, amount sdk.Coin) error {
	// Minimum deposit requirement
	minDeposit := sdk.NewCoin(types.DefaultDenom, sdk.NewInt(100000000000)) // 100,000 NAMO
	if amount.IsLT(minDeposit) {
		return types.ErrInsufficientDeposit
	}
	
	// Transfer to sevaMitra security module account
	securityAddr := k.authKeeper.GetModuleAddress(types.SevaMitraSecurityPool)
	return k.bankKeeper.SendCoins(ctx, depositor, securityAddr, sdk.NewCoins(amount))
}

func (k Keeper) isAuthorizedKYCVerifier(ctx sdk.Context, verifier sdk.AccAddress) bool {
	// Check if verifier is in authorized list
	// In production, would check against KYC provider registry
	return k.HasRole(ctx, verifier, "kyc_verifier")
}

func (k Keeper) validateKYCData(kycData *types.KYCData) bool {
	// Validate required fields
	if kycData.AadhaarHash == "" || kycData.PanHash == "" {
		return false
	}
	
	// Validate document hashes are proper format
	if len(kycData.AadhaarHash) != 64 || len(kycData.PanHash) != 64 {
		return false
	}
	
	// Check verification score
	if kycData.VerificationScore < 80 {
		return false
	}
	
	return true
}

func (k Keeper) isValidStatusTransition(current, new types.SevaMitraStatus) bool {
	// Define valid transitions
	validTransitions := map[types.SevaMitraStatus][]types.SevaMitraStatus{
		types.SevaMitraStatus_MITRA_STATUS_PENDING: {
			types.SevaMitraStatus_MITRA_STATUS_ACTIVE,
			types.SevaMitraStatus_MITRA_STATUS_INACTIVE,
		},
		types.SevaMitraStatus_MITRA_STATUS_ACTIVE: {
			types.SevaMitraStatus_MITRA_STATUS_SUSPENDED,
			types.SevaMitraStatus_MITRA_STATUS_INACTIVE,
		},
		types.SevaMitraStatus_MITRA_STATUS_SUSPENDED: {
			types.SevaMitraStatus_MITRA_STATUS_ACTIVE,
			types.SevaMitraStatus_MITRA_STATUS_INACTIVE,
		},
		types.SevaMitraStatus_MITRA_STATUS_INACTIVE: {
			types.SevaMitraStatus_MITRA_STATUS_ACTIVE,
		},
	}
	
	allowed, exists := validTransitions[current]
	if !exists {
		return false
	}
	
	for _, status := range allowed {
		if status == new {
			return true
		}
	}
	
	return false
}

func (k Keeper) validateSevaMitraLimits(dailyLimit, perTransactionLimit sdk.Coin, stats *types.SevaMitraStats) bool {
	// Per transaction cannot exceed daily limit
	if perTransactionLimit.IsGTE(dailyLimit) {
		return false
	}
	
	// Limits should be reasonable based on sevaMitra's history
	if stats.TotalTransactions > 100 {
		// Experienced sevaMitras can have higher limits
		maxDaily := sdk.NewCoin(types.DefaultDenom, sdk.NewInt(10000000000000)) // 10M NAMO
		if dailyLimit.IsGTE(maxDaily) {
			return false
		}
	} else {
		// New sevaMitras have lower limits
		maxDaily := sdk.NewCoin(types.DefaultDenom, sdk.NewInt(1000000000000)) // 1M NAMO
		if dailyLimit.IsGTE(maxDaily) {
			return false
		}
	}
	
	return true
}

func (k Keeper) updateSevaMitraDailyStats(ctx sdk.Context, sevaMitra *types.SevaMitra, amount sdk.Coin) {
	// Track daily volume for limit enforcement
	dateKey := ctx.BlockTime().Format("2006-01-02")
	dailyVolume := k.GetSevaMitraDailyVolume(ctx, sevaMitra.MitraId, dateKey)
	
	newVolume := dailyVolume.Add(amount)
	if newVolume.IsGTE(sevaMitra.DailyLimit) {
		// SevaMitra has reached daily limit
		sevaMitra.DailyLimitReached = true
		sevaMitra.DailyLimitResetAt = ctx.BlockTime().Add(24 * time.Hour)
	}
	
	k.SetSevaMitraDailyVolume(ctx, sevaMitra.MitraId, dateKey, newVolume)
}

func (k Keeper) calculateSevaMitraCommissionRate(stats *types.SevaMitraStats) string {
	// Performance-based commission
	baseRate := 2.5
	
	// Volume bonus (up to -0.5%)
	if stats.TotalVolume.Amount.GT(sdk.NewInt(100000000000000)) { // 100M NAMO
		baseRate -= 0.5
	} else if stats.TotalVolume.Amount.GT(sdk.NewInt(10000000000000)) { // 10M NAMO
		baseRate -= 0.3
	}
	
	// Rating bonus (up to -0.3%)
	if stats.AverageRating >= 4.5 {
		baseRate -= 0.3
	} else if stats.AverageRating >= 4.0 {
		baseRate -= 0.2
	}
	
	// Success rate bonus (up to -0.2%)
	if stats.TotalTransactions > 0 {
		successRate := float64(stats.SuccessfulTransactions) / float64(stats.TotalTransactions)
		if successRate >= 0.98 {
			baseRate -= 0.2
		} else if successRate >= 0.95 {
			baseRate -= 0.1
		}
	}
	
	// Minimum rate is 1.5%
	if baseRate < 1.5 {
		baseRate = 1.5
	}
	
	return fmt.Sprintf("%.1f", baseRate)
}

func (k Keeper) hasTransactedWithSevaMitra(ctx sdk.Context, user sdk.AccAddress, sevaMitraID string) bool {
	// Check transaction history
	// In production, would query transaction index
	return true // Simplified for now
}

func (k Keeper) updateSevaMitraTrustScore(ctx sdk.Context, sevaMitra *types.SevaMitra) {
	// Calculate trust score based on multiple factors
	score := float64(50) // Base score
	
	// Rating component (0-30 points)
	score += sevaMitra.Stats.AverageRating * 6
	
	// Volume component (0-20 points)
	volumeScore := float64(sevaMitra.Stats.TotalVolume.Amount.Int64()) / 1000000000000
	if volumeScore > 20 {
		volumeScore = 20
	}
	score += volumeScore
	
	// Success rate component (0-20 points)
	if sevaMitra.Stats.TotalTransactions > 0 {
		successRate := float64(sevaMitra.Stats.SuccessfulTransactions) / float64(sevaMitra.Stats.TotalTransactions)
		score += successRate * 20
	}
	
	// Dispute component (-20 to +10 points)
	if sevaMitra.Stats.DisputesResolved > 0 || sevaMitra.Stats.DisputesLost > 0 {
		disputeRate := float64(sevaMitra.Stats.DisputesLost) / float64(sevaMitra.Stats.DisputesResolved + sevaMitra.Stats.DisputesLost)
		score -= disputeRate * 20
		score += float64(sevaMitra.Stats.DisputesResolved-sevaMitra.Stats.DisputesLost) * 2
	}
	
	// Cap between 0-100
	if score < 0 {
		score = 0
	} else if score > 100 {
		score = 100
	}
	
	sevaMitra.TrustScore = int32(score)
}

// ScheduleKYCVerification schedules KYC verification for an sevaMitra
func (k Keeper) ScheduleKYCVerification(ctx sdk.Context, sevaMitraID string) {
	// In production, would integrate with KYC provider
	// For now, add to verification queue
	k.AddToKYCQueue(ctx, sevaMitraID, ctx.BlockTime().Add(24*time.Hour))
}