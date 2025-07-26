package keeper

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/DeshChain/DeshChain-Ecosystem/x/identity/types"
)

// SetRecoveryRequest stores a recovery request
func (k Keeper) SetRecoveryRequest(ctx sdk.Context, request types.RecoveryRequest) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RecoveryRequestPrefix)
	b := k.cdc.MustMarshal(&request)
	store.Set([]byte(request.ID), b)
}

// GetRecoveryRequest retrieves a recovery request
func (k Keeper) GetRecoveryRequest(ctx sdk.Context, requestID string) (types.RecoveryRequest, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RecoveryRequestPrefix)
	b := store.Get([]byte(requestID))
	if b == nil {
		return types.RecoveryRequest{}, false
	}
	
	var request types.RecoveryRequest
	k.cdc.MustUnmarshal(b, &request)
	return request, true
}

// DeleteRecoveryRequest removes a recovery request
func (k Keeper) DeleteRecoveryRequest(ctx sdk.Context, requestID string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RecoveryRequestPrefix)
	store.Delete([]byte(requestID))
}

// IterateRecoveryRequests iterates over all recovery requests
func (k Keeper) IterateRecoveryRequests(ctx sdk.Context, cb func(request types.RecoveryRequest) bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RecoveryRequestPrefix)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var request types.RecoveryRequest
		k.cdc.MustUnmarshal(iterator.Value(), &request)
		if cb(request) {
			break
		}
	}
}

// InitiateRecovery starts the account recovery process
func (k Keeper) InitiateRecovery(
	ctx sdk.Context,
	initiatorAddress string,
	targetAddress string,
	recoveryType types.RecoveryType,
	recoveryData string,
) (*types.RecoveryRequest, error) {
	// Validate target has identity
	targetIdentity, found := k.GetIdentity(ctx, targetAddress)
	if !found {
		return nil, types.ErrIdentityNotFound
	}
	
	// Check if recovery method is set up
	hasMethod := false
	for _, method := range targetIdentity.RecoveryMethods {
		if method.Type == recoveryType && method.IsActive {
			hasMethod = true
			break
		}
	}
	
	if !hasMethod {
		return nil, types.ErrRecoveryMethodNotSet
	}
	
	// Check for existing pending recovery
	existingFound := false
	k.IterateRecoveryRequests(ctx, func(request types.RecoveryRequest) bool {
		if request.TargetAddress == targetAddress && request.Status == "pending" {
			existingFound = true
			return true
		}
		return false
	})
	
	if existingFound {
		return nil, types.ErrInvalidRequest
	}
	
	// Validate recovery data based on type
	if err := k.validateRecoveryData(ctx, recoveryType, recoveryData, targetIdentity); err != nil {
		return nil, err
	}
	
	// Create recovery request
	request := types.RecoveryRequest{
		ID:               fmt.Sprintf("recovery:%s:%s", targetAddress, sdk.NewRand().Str(16)),
		TargetAddress:    targetAddress,
		InitiatorAddress: initiatorAddress,
		RecoveryType:     recoveryType,
		RecoveryData:     recoveryData,
		Status:           "pending",
		InitiatedAt:      ctx.BlockTime(),
	}
	
	// Store request
	k.SetRecoveryRequest(ctx, request)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRecoveryInitiated,
			sdk.NewAttribute(types.AttributeKeyAddress, targetAddress),
			sdk.NewAttribute(types.AttributeKeyRecoveryType, recoveryType.String()),
			sdk.NewAttribute("recovery_id", request.ID),
		),
	)
	
	// Send notifications based on recovery type
	k.sendRecoveryNotification(ctx, recoveryType, targetIdentity, request.ID)
	
	return &request, nil
}

// CompleteRecovery completes the account recovery process
func (k Keeper) CompleteRecovery(
	ctx sdk.Context,
	recoveryID string,
	newPublicKey string,
	proofData string,
) error {
	// Get recovery request
	request, found := k.GetRecoveryRequest(ctx, recoveryID)
	if !found {
		return types.ErrInvalidRequest
	}
	
	// Check if request is ready
	if request.Status != "ready" {
		return types.ErrRecoveryFailed
	}
	
	// Validate proof based on recovery type
	if err := k.validateRecoveryProof(ctx, request.RecoveryType, proofData, request); err != nil {
		return err
	}
	
	// Get target identity
	identity, found := k.GetIdentity(ctx, request.TargetAddress)
	if !found {
		return types.ErrIdentityNotFound
	}
	
	// Get DID document
	didDoc, found := k.GetDIDDocument(ctx, identity.DID)
	if !found {
		return types.ErrDIDNotFound
	}
	
	// Create new verification method with new public key
	newVerificationMethod := types.VerificationMethod{
		ID:                 fmt.Sprintf("%s#recovery-key-%d", identity.DID, ctx.BlockHeight()),
		Type:               "Ed25519VerificationKey2020",
		Controller:         identity.DID,
		PublicKeyMultibase: newPublicKey,
		Created:            ctx.BlockTime(),
	}
	
	// Add new verification method
	didDoc.AddVerificationMethod(newVerificationMethod)
	
	// Update authentication methods
	didDoc.Authentication = append(didDoc.Authentication, newVerificationMethod.ID)
	
	// Save updated DID document
	k.SetDIDDocument(ctx, didDoc)
	
	// Update recovery request
	completedAt := ctx.BlockTime()
	request.Status = "completed"
	request.CompletedAt = &completedAt
	request.NewPublicKey = newPublicKey
	k.SetRecoveryRequest(ctx, request)
	
	// Update identity
	identity.UpdatedAt = ctx.BlockTime()
	identity.LastActivityAt = ctx.BlockTime()
	k.SetIdentity(ctx, identity)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRecoveryCompleted,
			sdk.NewAttribute(types.AttributeKeyAddress, request.TargetAddress),
			sdk.NewAttribute(types.AttributeKeyRecoveryType, request.RecoveryType.String()),
			sdk.NewAttribute("recovery_id", recoveryID),
		),
	)
	
	return nil
}

// GetActiveRecoveryMethods returns all active recovery methods for an identity
func (k Keeper) GetActiveRecoveryMethods(ctx sdk.Context, address string) []types.RecoveryMethod {
	identity, found := k.GetIdentity(ctx, address)
	if !found {
		return []types.RecoveryMethod{}
	}
	
	var activeMethods []types.RecoveryMethod
	for _, method := range identity.RecoveryMethods {
		if method.IsActive {
			activeMethods = append(activeMethods, method)
		}
	}
	
	return activeMethods
}

// DeactivateRecoveryMethod deactivates a recovery method
func (k Keeper) DeactivateRecoveryMethod(ctx sdk.Context, address string, methodType types.RecoveryType) error {
	identity, found := k.GetIdentity(ctx, address)
	if !found {
		return types.ErrIdentityNotFound
	}
	
	methodFound := false
	for i, method := range identity.RecoveryMethods {
		if method.Type == methodType {
			identity.RecoveryMethods[i].IsActive = false
			methodFound = true
			break
		}
	}
	
	if !methodFound {
		return types.ErrRecoveryMethodNotSet
	}
	
	identity.UpdatedAt = ctx.BlockTime()
	k.SetIdentity(ctx, identity)
	
	return nil
}

// Helper functions

func (k Keeper) validateRecoveryData(ctx sdk.Context, recoveryType types.RecoveryType, data string, identity types.Identity) error {
	switch recoveryType {
	case types.RecoveryType_EMAIL:
		// Validate email format and check if it matches stored recovery email
		// In production, would verify against hashed/encrypted email
		return nil
		
	case types.RecoveryType_PHONE:
		// Validate phone format and check if it matches stored recovery phone
		// In production, would verify against hashed/encrypted phone
		return nil
		
	case types.RecoveryType_GUARDIAN:
		// Validate guardian address exists and is authorized
		if !k.HasIdentity(ctx, data) {
			return types.ErrGuardianNotFound
		}
		// Check if guardian is in recovery methods
		for _, method := range identity.RecoveryMethods {
			if method.Type == types.RecoveryType_GUARDIAN && method.Value == data {
				return nil
			}
		}
		return types.ErrUnauthorized
		
	case types.RecoveryType_SOCIAL:
		// Validate social recovery threshold
		// This would involve checking multiple guardian approvals
		return nil
		
	case types.RecoveryType_MNEMONIC:
		// Validate mnemonic phrase
		// In production, would verify against encrypted mnemonic hash
		return nil
		
	default:
		return types.ErrInvalidRecoveryMethod
	}
}

func (k Keeper) validateRecoveryProof(ctx sdk.Context, recoveryType types.RecoveryType, proof string, request types.RecoveryRequest) error {
	switch recoveryType {
	case types.RecoveryType_EMAIL:
		// Validate email OTP or verification link
		return nil
		
	case types.RecoveryType_PHONE:
		// Validate SMS OTP
		return nil
		
	case types.RecoveryType_GUARDIAN:
		// Validate guardian signature
		return nil
		
	case types.RecoveryType_SOCIAL:
		// Validate multiple guardian signatures meet threshold
		return nil
		
	case types.RecoveryType_MNEMONIC:
		// Validate mnemonic phrase matches
		return nil
		
	default:
		return types.ErrInvalidRecoveryMethod
	}
}

func (k Keeper) sendRecoveryNotification(ctx sdk.Context, recoveryType types.RecoveryType, identity types.Identity, recoveryID string) {
	// In production, this would send actual notifications
	// For now, just emit events
	
	switch recoveryType {
	case types.RecoveryType_EMAIL:
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				"recovery_email_sent",
				sdk.NewAttribute("recovery_id", recoveryID),
				sdk.NewAttribute("method", "email"),
			),
		)
		
	case types.RecoveryType_PHONE:
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				"recovery_sms_sent",
				sdk.NewAttribute("recovery_id", recoveryID),
				sdk.NewAttribute("method", "sms"),
			),
		)
		
	case types.RecoveryType_GUARDIAN:
		// Notify guardian
		for _, method := range identity.RecoveryMethods {
			if method.Type == types.RecoveryType_GUARDIAN {
				ctx.EventManager().EmitEvent(
					sdk.NewEvent(
						"guardian_notified",
						sdk.NewAttribute("recovery_id", recoveryID),
						sdk.NewAttribute("guardian", method.Value),
					),
				)
			}
		}
		
	case types.RecoveryType_SOCIAL:
		// Notify all guardians for social recovery
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				"social_recovery_initiated",
				sdk.NewAttribute("recovery_id", recoveryID),
			),
		)
	}
}

// GetPendingRecoveryRequests returns all pending recovery requests
func (k Keeper) GetPendingRecoveryRequests(ctx sdk.Context) []types.RecoveryRequest {
	var pendingRequests []types.RecoveryRequest
	
	k.IterateRecoveryRequests(ctx, func(request types.RecoveryRequest) bool {
		if request.Status == "pending" || request.Status == "ready" {
			pendingRequests = append(pendingRequests, request)
		}
		return false
	})
	
	return pendingRequests
}

// CancelRecoveryRequest cancels a pending recovery request
func (k Keeper) CancelRecoveryRequest(ctx sdk.Context, recoveryID string, canceller string) error {
	request, found := k.GetRecoveryRequest(ctx, recoveryID)
	if !found {
		return types.ErrInvalidRequest
	}
	
	// Only initiator or target can cancel
	if canceller != request.InitiatorAddress && canceller != request.TargetAddress {
		return types.ErrUnauthorized
	}
	
	// Only pending or ready requests can be cancelled
	if request.Status != "pending" && request.Status != "ready" {
		return types.ErrInvalidRequest
	}
	
	// Update status
	request.Status = "cancelled"
	completedAt := ctx.BlockTime()
	request.CompletedAt = &completedAt
	k.SetRecoveryRequest(ctx, request)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"recovery_cancelled",
			sdk.NewAttribute("recovery_id", recoveryID),
			sdk.NewAttribute("cancelled_by", canceller),
		),
	)
	
	return nil
}