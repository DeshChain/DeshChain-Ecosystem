package keeper

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/deshchain/x/identity/types"
)

// SetConsentRecord stores a consent record
func (k Keeper) SetConsentRecord(ctx sdk.Context, consent types.ConsentRecord) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ConsentRecordPrefix)
	b := k.cdc.MustMarshal(&consent)
	store.Set([]byte(consent.ID), b)
	
	// Create index for consent by type
	k.SetConsentTypeIndex(ctx, consent.DataController, consent.Type, consent.ID)
}

// GetConsentRecord retrieves a consent record
func (k Keeper) GetConsentRecord(ctx sdk.Context, consentID string) (types.ConsentRecord, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ConsentRecordPrefix)
	b := store.Get([]byte(consentID))
	if b == nil {
		return types.ConsentRecord{}, false
	}
	
	var consent types.ConsentRecord
	k.cdc.MustUnmarshal(b, &consent)
	return consent, true
}

// SetConsentTypeIndex creates an index for consents by type
func (k Keeper) SetConsentTypeIndex(ctx sdk.Context, address string, consentType types.ConsentType, consentID string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ConsentIndexPrefix)
	key := types.GetConsentByTypeIndexKey(address, consentType.String())
	store.Set(append(key, []byte(consentID)...), []byte{1})
}

// RemoveConsentTypeIndex removes a consent from type index
func (k Keeper) RemoveConsentTypeIndex(ctx sdk.Context, address string, consentType types.ConsentType, consentID string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ConsentIndexPrefix)
	key := types.GetConsentByTypeIndexKey(address, consentType.String())
	store.Delete(append(key, []byte(consentID)...))
}

// GetConsentsByType returns all consents of a specific type for an address
func (k Keeper) GetConsentsByType(ctx sdk.Context, address string, consentType types.ConsentType) []types.ConsentRecord {
	var consents []types.ConsentRecord
	
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ConsentIndexPrefix)
	iterator := sdk.KVStorePrefixIterator(store, types.GetConsentByTypeIndexKey(address, consentType.String()))
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		// Extract consent ID from the key
		key := string(iterator.Key())
		baseKey := string(types.GetConsentByTypeIndexKey(address, consentType.String()))
		consentID := key[len(baseKey):]
		
		if consent, found := k.GetConsentRecord(ctx, consentID); found {
			consents = append(consents, consent)
		}
	}
	
	return consents
}

// GiveConsent creates a new consent record
func (k Keeper) GiveConsent(
	ctx sdk.Context,
	userAddress string,
	consentType types.ConsentType,
	purpose string,
	dataController string,
	dataCategories []string,
	processingTypes []string,
	expirationDays int32,
) (*types.ConsentRecord, error) {
	// Validate user has identity
	if !k.HasIdentity(ctx, userAddress) {
		return nil, types.ErrIdentityNotFound
	}
	
	// Check if similar consent already exists
	existingConsents := k.GetConsentsByType(ctx, userAddress, consentType)
	for _, existing := range existingConsents {
		if existing.DataController == dataController && existing.Given && existing.WithdrawnAt == nil {
			return nil, types.ErrConsentAlreadyExists
		}
	}
	
	// Calculate expiration
	var expiresAt *time.Time
	if expirationDays > 0 {
		expiry := ctx.BlockTime().AddDate(0, 0, int(expirationDays))
		expiresAt = &expiry
	} else {
		// Use default max consent duration
		maxDays := k.MaxConsentDurationDays(ctx)
		expiry := ctx.BlockTime().AddDate(0, 0, int(maxDays))
		expiresAt = &expiry
	}
	
	// Create consent record
	consent := types.ConsentRecord{
		ID:              fmt.Sprintf("consent:%s:%s:%s", userAddress, consentType.String(), sdk.NewRand().Str(16)),
		Type:            consentType,
		Purpose:         purpose,
		DataController:  dataController,
		DataCategories:  dataCategories,
		ProcessingTypes: processingTypes,
		Given:           true,
		GivenAt:         ctx.BlockTime(),
		ExpiresAt:       expiresAt,
		Version:         "1.0",
	}
	
	// Store consent
	k.SetConsentRecord(ctx, consent)
	
	// Add to identity
	if err := k.AddIdentityConsent(ctx, userAddress, consent); err != nil {
		return nil, err
	}
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeConsentGiven,
			sdk.NewAttribute(types.AttributeKeyAddress, userAddress),
			sdk.NewAttribute(types.AttributeKeyConsentType, consentType.String()),
			sdk.NewAttribute("data_controller", dataController),
			sdk.NewAttribute("purpose", purpose),
		),
	)
	
	return &consent, nil
}

// WithdrawConsent withdraws a given consent
func (k Keeper) WithdrawConsent(ctx sdk.Context, userAddress string, consentID string, reason string) error {
	// Get consent record
	consent, found := k.GetConsentRecord(ctx, consentID)
	if !found {
		return types.ErrInvalidRequest
	}
	
	// Verify ownership (consent should be in user's identity)
	identity, found := k.GetIdentity(ctx, userAddress)
	if !found {
		return types.ErrIdentityNotFound
	}
	
	// Check if consent belongs to user
	consentFound := false
	for _, identityConsent := range identity.Consents {
		if identityConsent.ID == consentID {
			consentFound = true
			break
		}
	}
	
	if !consentFound {
		return types.ErrUnauthorized
	}
	
	// Check if already withdrawn
	if !consent.Given || consent.WithdrawnAt != nil {
		return types.ErrConsentWithdrawn
	}
	
	// Update consent record
	withdrawnAt := ctx.BlockTime()
	consent.Given = false
	consent.WithdrawnAt = &withdrawnAt
	k.SetConsentRecord(ctx, consent)
	
	// Update identity consent
	if err := k.WithdrawIdentityConsent(ctx, userAddress, consentID); err != nil {
		return err
	}
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeConsentWithdrawn,
			sdk.NewAttribute(types.AttributeKeyAddress, userAddress),
			sdk.NewAttribute(types.AttributeKeyConsentType, consent.Type.String()),
			sdk.NewAttribute("consent_id", consentID),
			sdk.NewAttribute("reason", reason),
		),
	)
	
	return nil
}

// CheckConsent verifies if a user has given consent for a specific purpose
func (k Keeper) CheckConsent(
	ctx sdk.Context,
	userAddress string,
	consentType types.ConsentType,
	dataController string,
) (bool, error) {
	identity, found := k.GetIdentity(ctx, userAddress)
	if !found {
		return false, types.ErrIdentityNotFound
	}
	
	// Get privacy settings
	privacySettings, _ := k.GetPrivacySettings(ctx, userAddress)
	
	// Check if explicit consent is required
	if privacySettings.RequireExplicitConsent {
		// Look for active consent
		for _, consent := range identity.GetActiveConsents() {
			if consent.Type == consentType && consent.DataController == dataController {
				// Check if not expired
				if consent.ExpiresAt == nil || consent.ExpiresAt.After(ctx.BlockTime()) {
					return true, nil
				}
			}
		}
		return false, types.ErrConsentNotGiven
	}
	
	// If explicit consent not required, check if type is allowed by default
	// This would depend on the consent type and privacy settings
	return true, nil
}

// GetActiveConsentsForController returns all active consents for a data controller
func (k Keeper) GetActiveConsentsForController(ctx sdk.Context, dataController string) []types.ConsentRecord {
	var activeConsents []types.ConsentRecord
	
	// Iterate through all identities
	k.IterateIdentities(ctx, func(identity types.Identity) bool {
		for _, consent := range identity.GetActiveConsents() {
			if consent.DataController == dataController {
				activeConsents = append(activeConsents, consent)
			}
		}
		return false
	})
	
	return activeConsents
}

// ValidateDataAccess validates if data access is allowed based on consent
func (k Keeper) ValidateDataAccess(
	ctx sdk.Context,
	dataSubject string,
	dataController string,
	dataCategory string,
	processingType string,
) error {
	identity, found := k.GetIdentity(ctx, dataSubject)
	if !found {
		return types.ErrIdentityNotFound
	}
	
	// Check for active consent
	consentFound := false
	for _, consent := range identity.GetActiveConsents() {
		if consent.DataController != dataController {
			continue
		}
		
		// Check data category
		categoryAllowed := false
		for _, category := range consent.DataCategories {
			if category == dataCategory || category == "*" {
				categoryAllowed = true
				break
			}
		}
		
		if !categoryAllowed {
			continue
		}
		
		// Check processing type
		processingAllowed := false
		for _, procType := range consent.ProcessingTypes {
			if procType == processingType || procType == "*" {
				processingAllowed = true
				break
			}
		}
		
		if processingAllowed {
			consentFound = true
			break
		}
	}
	
	if !consentFound {
		return types.ErrConsentNotGiven
	}
	
	// Log data access for audit
	k.logDataAccess(ctx, dataSubject, dataController, dataCategory, processingType)
	
	return nil
}

// GetConsentHistory returns the consent history for a user
func (k Keeper) GetConsentHistory(ctx sdk.Context, userAddress string) []types.ConsentRecord {
	identity, found := k.GetIdentity(ctx, userAddress)
	if !found {
		return []types.ConsentRecord{}
	}
	
	// Return all consents (including withdrawn ones)
	return identity.Consents
}

// RevokeAllConsentsForController revokes all consents for a specific data controller
func (k Keeper) RevokeAllConsentsForController(ctx sdk.Context, dataController string, reason string) error {
	revokedCount := 0
	
	// Iterate through all identities
	k.IterateIdentities(ctx, func(identity types.Identity) bool {
		for _, consent := range identity.GetActiveConsents() {
			if consent.DataController == dataController {
				// Withdraw consent
				if err := k.WithdrawConsent(ctx, identity.Address, consent.ID, reason); err == nil {
					revokedCount++
				}
			}
		}
		return false
	})
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"consents_bulk_revoked",
			sdk.NewAttribute("data_controller", dataController),
			sdk.NewAttribute("revoked_count", fmt.Sprintf("%d", revokedCount)),
			sdk.NewAttribute("reason", reason),
		),
	)
	
	return nil
}

// Helper function to log data access
func (k Keeper) logDataAccess(ctx sdk.Context, dataSubject, dataController, dataCategory, processingType string) {
	// In production, this would create an audit log entry
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"data_access_logged",
			sdk.NewAttribute("data_subject", dataSubject),
			sdk.NewAttribute("data_controller", dataController),
			sdk.NewAttribute("data_category", dataCategory),
			sdk.NewAttribute("processing_type", processingType),
			sdk.NewAttribute("timestamp", ctx.BlockTime().Format(time.RFC3339)),
		),
	)
}