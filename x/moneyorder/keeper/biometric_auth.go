package keeper

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/DeshChain/DeshChain-Ecosystem/x/moneyorder/types"
)

// BiometricAuthManager handles biometric authentication for Money Orders
type BiometricAuthManager struct {
	keeper *Keeper
}

// NewBiometricAuthManager creates a new biometric authentication manager
func NewBiometricAuthManager(keeper *Keeper) *BiometricAuthManager {
	return &BiometricAuthManager{
		keeper: keeper,
	}
}

// RegisterBiometric registers a new biometric template for a user
func (bam *BiometricAuthManager) RegisterBiometric(
	ctx sdk.Context,
	userAddress string,
	biometricType types.BiometricType,
	templateHash string,
	deviceID string,
) error {
	// Validate input parameters
	if userAddress == "" {
		return sdkerrors.Wrap(types.ErrInvalidAddress, "user address cannot be empty")
	}
	
	if templateHash == "" {
		return sdkerrors.Wrap(types.ErrInvalidBiometric, "template hash cannot be empty")
	}
	
	if deviceID == "" {
		return sdkerrors.Wrap(types.ErrInvalidDevice, "device ID cannot be empty")
	}

	// Check if user exists and is active
	_, found := bam.keeper.GetUser(ctx, userAddress)
	if !found {
		return sdkerrors.Wrap(types.ErrUserNotFound, "user not found for biometric registration")
	}

	// Create biometric registration
	biometric := types.BiometricRegistration{
		UserAddress:   userAddress,
		BiometricType: biometricType,
		TemplateHash:  templateHash,
		DeviceId:      deviceID,
		RegisteredAt:  time.Now(),
		IsActive:      true,
		FailCount:     0,
		LastUsed:      time.Now(),
	}

	// Generate unique biometric ID
	biometricID := bam.generateBiometricID(userAddress, deviceID, biometricType)
	biometric.BiometricId = biometricID

	// Store biometric registration
	store := prefix.NewStore(ctx.KVStore(bam.keeper.storeKey), types.BiometricRegistrationPrefix)
	bz := bam.keeper.cdc.MustMarshal(&biometric)
	store.Set([]byte(biometricID), bz)

	// Update user's biometric status
	bam.updateUserBiometricStatus(ctx, userAddress, true)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeBiometricRegistered,
			sdk.NewAttribute(types.AttributeKeyUserAddress, userAddress),
			sdk.NewAttribute(types.AttributeKeyBiometricType, biometricType.String()),
			sdk.NewAttribute(types.AttributeKeyDeviceID, deviceID),
			sdk.NewAttribute(types.AttributeKeyBiometricID, biometricID),
		),
	)

	return nil
}

// AuthenticateBiometric verifies a biometric authentication attempt
func (bam *BiometricAuthManager) AuthenticateBiometric(
	ctx sdk.Context,
	userAddress string,
	biometricType types.BiometricType,
	templateHash string,
	deviceID string,
) (*types.BiometricAuthResult, error) {
	// Find biometric registration
	biometricID := bam.generateBiometricID(userAddress, deviceID, biometricType)
	biometric, found := bam.getBiometricRegistration(ctx, biometricID)
	if !found {
		return &types.BiometricAuthResult{
			Success:      false,
			ErrorMessage: "Biometric not registered for this device",
			AuthScore:    0,
		}, nil
	}

	// Check if biometric is active
	if !biometric.IsActive {
		return &types.BiometricAuthResult{
			Success:      false,
			ErrorMessage: "Biometric authentication is disabled",
			AuthScore:    0,
		}, nil
	}

	// Check fail count limits
	if biometric.FailCount >= 5 {
		// Temporarily disable biometric after 5 failed attempts
		bam.disableBiometric(ctx, biometricID, "Too many failed attempts")
		return &types.BiometricAuthResult{
			Success:      false,
			ErrorMessage: "Biometric authentication locked due to too many failed attempts",
			AuthScore:    0,
		}, nil
	}

	// Verify template hash
	authScore := bam.calculateBiometricScore(biometric.TemplateHash, templateHash)
	success := authScore >= 0.85 // 85% match threshold

	if success {
		// Reset fail count and update last used
		biometric.FailCount = 0
		biometric.LastUsed = time.Now()
		bam.updateBiometricRegistration(ctx, biometricID, biometric)

		// Record successful authentication
		bam.recordAuthenticationAttempt(ctx, userAddress, biometricID, true, authScore)

		// Emit success event
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeBiometricAuthSuccess,
				sdk.NewAttribute(types.AttributeKeyUserAddress, userAddress),
				sdk.NewAttribute(types.AttributeKeyBiometricType, biometricType.String()),
				sdk.NewAttribute(types.AttributeKeyAuthScore, fmt.Sprintf("%.2f", authScore)),
			),
		)

		return &types.BiometricAuthResult{
			Success:   true,
			AuthScore: authScore,
			BiometricId: biometricID,
		}, nil
	} else {
		// Increment fail count
		biometric.FailCount++
		bam.updateBiometricRegistration(ctx, biometricID, biometric)

		// Record failed authentication
		bam.recordAuthenticationAttempt(ctx, userAddress, biometricID, false, authScore)

		// Emit failure event
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeBiometricAuthFailed,
				sdk.NewAttribute(types.AttributeKeyUserAddress, userAddress),
				sdk.NewAttribute(types.AttributeKeyBiometricType, biometricType.String()),
				sdk.NewAttribute(types.AttributeKeyAuthScore, fmt.Sprintf("%.2f", authScore)),
				sdk.NewAttribute(types.AttributeKeyFailCount, fmt.Sprintf("%d", biometric.FailCount)),
			),
		)

		return &types.BiometricAuthResult{
			Success:      false,
			ErrorMessage: "Biometric authentication failed",
			AuthScore:    authScore,
		}, nil
	}
}

// AuthorizeMoneyOrderWithBiometric authorizes a money order transaction using biometric authentication
func (bam *BiometricAuthManager) AuthorizeMoneyOrderWithBiometric(
	ctx sdk.Context,
	orderID string,
	userAddress string,
	biometricType types.BiometricType,
	templateHash string,
	deviceID string,
) error {
	// Get money order
	order, found := bam.keeper.GetMoneyOrder(ctx, orderID)
	if !found {
		return sdkerrors.Wrap(types.ErrOrderNotFound, "money order not found")
	}

	// Verify user is authorized to perform this action
	if order.SenderAddress != userAddress {
		return sdkerrors.Wrap(types.ErrUnauthorized, "user not authorized for this money order")
	}

	// Check if order is in pending state
	if order.Status != types.MoneyOrderStatus_PENDING {
		return sdkerrors.Wrap(types.ErrInvalidOrderStatus, "order must be pending for biometric authorization")
	}

	// Perform biometric authentication
	authResult, err := bam.AuthenticateBiometric(ctx, userAddress, biometricType, templateHash, deviceID)
	if err != nil {
		return sdkerrors.Wrap(err, "biometric authentication error")
	}

	if !authResult.Success {
		return sdkerrors.Wrap(types.ErrBiometricAuthFailed, authResult.ErrorMessage)
	}

	// Update order with biometric authorization
	order.BiometricAuthRequired = true
	order.BiometricAuthCompleted = true
	order.BiometricAuthScore = authResult.AuthScore
	order.BiometricAuthTimestamp = time.Now()
	order.AuthorizedDeviceId = deviceID

	// Update order status to confirmed if all requirements met
	if bam.isOrderFullyAuthorized(order) {
		order.Status = types.MoneyOrderStatus_CONFIRMED
		order.ConfirmedAt = time.Now()
	}

	// Save updated order
	bam.keeper.SetMoneyOrder(ctx, order)

	// Emit authorization event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMoneyOrderBiometricAuth,
			sdk.NewAttribute(types.AttributeKeyOrderID, orderID),
			sdk.NewAttribute(types.AttributeKeyUserAddress, userAddress),
			sdk.NewAttribute(types.AttributeKeyBiometricType, biometricType.String()),
			sdk.NewAttribute(types.AttributeKeyAuthScore, fmt.Sprintf("%.2f", authResult.AuthScore)),
			sdk.NewAttribute(types.AttributeKeyOrderStatus, order.Status.String()),
		),
	)

	return nil
}

// DisableBiometric disables a biometric registration
func (bam *BiometricAuthManager) DisableBiometric(
	ctx sdk.Context,
	userAddress string,
	biometricID string,
	reason string,
) error {
	// Verify user ownership
	biometric, found := bam.getBiometricRegistration(ctx, biometricID)
	if !found {
		return sdkerrors.Wrap(types.ErrBiometricNotFound, "biometric registration not found")
	}

	if biometric.UserAddress != userAddress {
		return sdkerrors.Wrap(types.ErrUnauthorized, "user not authorized to disable this biometric")
	}

	// Disable biometric
	bam.disableBiometric(ctx, biometricID, reason)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeBiometricDisabled,
			sdk.NewAttribute(types.AttributeKeyUserAddress, userAddress),
			sdk.NewAttribute(types.AttributeKeyBiometricID, biometricID),
			sdk.NewAttribute(types.AttributeKeyReason, reason),
		),
	)

	return nil
}

// GetUserBiometrics returns all biometric registrations for a user
func (bam *BiometricAuthManager) GetUserBiometrics(
	ctx sdk.Context,
	userAddress string,
) []types.BiometricRegistration {
	var biometrics []types.BiometricRegistration
	
	store := prefix.NewStore(ctx.KVStore(bam.keeper.storeKey), types.BiometricRegistrationPrefix)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var biometric types.BiometricRegistration
		bam.keeper.cdc.MustUnmarshal(iterator.Value(), &biometric)
		
		if biometric.UserAddress == userAddress {
			biometrics = append(biometrics, biometric)
		}
	}

	return biometrics
}

// Internal helper functions

func (bam *BiometricAuthManager) generateBiometricID(userAddress, deviceID string, biometricType types.BiometricType) string {
	data := fmt.Sprintf("%s:%s:%s", userAddress, deviceID, biometricType.String())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])[:16] // Use first 16 characters
}

func (bam *BiometricAuthManager) getBiometricRegistration(ctx sdk.Context, biometricID string) (types.BiometricRegistration, bool) {
	store := prefix.NewStore(ctx.KVStore(bam.keeper.storeKey), types.BiometricRegistrationPrefix)
	bz := store.Get([]byte(biometricID))
	if bz == nil {
		return types.BiometricRegistration{}, false
	}

	var biometric types.BiometricRegistration
	bam.keeper.cdc.MustUnmarshal(bz, &biometric)
	return biometric, true
}

func (bam *BiometricAuthManager) updateBiometricRegistration(ctx sdk.Context, biometricID string, biometric types.BiometricRegistration) {
	store := prefix.NewStore(ctx.KVStore(bam.keeper.storeKey), types.BiometricRegistrationPrefix)
	bz := bam.keeper.cdc.MustMarshal(&biometric)
	store.Set([]byte(biometricID), bz)
}

func (bam *BiometricAuthManager) disableBiometric(ctx sdk.Context, biometricID string, reason string) {
	biometric, found := bam.getBiometricRegistration(ctx, biometricID)
	if found {
		biometric.IsActive = false
		biometric.DisabledAt = time.Now()
		biometric.DisabledReason = reason
		bam.updateBiometricRegistration(ctx, biometricID, biometric)
	}
}

func (bam *BiometricAuthManager) updateUserBiometricStatus(ctx sdk.Context, userAddress string, hasBiometric bool) {
	user, found := bam.keeper.GetUser(ctx, userAddress)
	if found {
		user.BiometricEnabled = hasBiometric
		bam.keeper.SetUser(ctx, user)
	}
}

func (bam *BiometricAuthManager) calculateBiometricScore(storedHash, providedHash string) float64 {
	// Simple similarity calculation - in production, use proper biometric matching algorithms
	if storedHash == providedHash {
		return 1.0
	}
	
	// Calculate Hamming distance for demonstration
	if len(storedHash) != len(providedHash) {
		return 0.0
	}
	
	matches := 0
	for i := 0; i < len(storedHash); i++ {
		if storedHash[i] == providedHash[i] {
			matches++
		}
	}
	
	return float64(matches) / float64(len(storedHash))
}

func (bam *BiometricAuthManager) recordAuthenticationAttempt(
	ctx sdk.Context,
	userAddress string,
	biometricID string,
	success bool,
	score float64,
) {
	attempt := types.BiometricAuthAttempt{
		UserAddress:   userAddress,
		BiometricId:   biometricID,
		Success:       success,
		AuthScore:     score,
		AttemptTime:   time.Now(),
		ClientIP:      ctx.Context().Value("client_ip").(string), // If available
	}

	// Store authentication attempt for audit trail
	store := prefix.NewStore(ctx.KVStore(bam.keeper.storeKey), types.BiometricAuthAttemptPrefix)
	attemptID := fmt.Sprintf("%s:%d", biometricID, time.Now().Unix())
	bz := bam.keeper.cdc.MustMarshal(&attempt)
	store.Set([]byte(attemptID), bz)
}

func (bam *BiometricAuthManager) isOrderFullyAuthorized(order types.MoneyOrder) bool {
	// Check if all required authorizations are complete
	if order.BiometricAuthRequired && !order.BiometricAuthCompleted {
		return false
	}
	
	// Add other authorization checks as needed (PIN, 2FA, etc.)
	
	return true
}

// RequireBiometricForOrder marks an order as requiring biometric authentication
func (bam *BiometricAuthManager) RequireBiometricForOrder(
	ctx sdk.Context,
	orderID string,
	required bool,
) error {
	order, found := bam.keeper.GetMoneyOrder(ctx, orderID)
	if !found {
		return sdkerrors.Wrap(types.ErrOrderNotFound, "money order not found")
	}

	order.BiometricAuthRequired = required
	if !required {
		order.BiometricAuthCompleted = true // If not required, mark as completed
	}

	bam.keeper.SetMoneyOrder(ctx, order)
	return nil
}

// GetBiometricSecurityLevel returns the security level based on biometric configuration
func (bam *BiometricAuthManager) GetBiometricSecurityLevel(ctx sdk.Context, userAddress string) types.SecurityLevel {
	biometrics := bam.GetUserBiometrics(ctx, userAddress)
	
	if len(biometrics) == 0 {
		return types.SecurityLevel_BASIC
	}
	
	activeCount := 0
	hasMultipleTypes := false
	typesMap := make(map[types.BiometricType]bool)
	
	for _, bio := range biometrics {
		if bio.IsActive {
			activeCount++
			typesMap[bio.BiometricType] = true
		}
	}
	
	if len(typesMap) > 1 {
		hasMultipleTypes = true
	}
	
	if activeCount >= 2 && hasMultipleTypes {
		return types.SecurityLevel_PREMIUM
	} else if activeCount >= 1 {
		return types.SecurityLevel_ENHANCED
	}
	
	return types.SecurityLevel_BASIC
}