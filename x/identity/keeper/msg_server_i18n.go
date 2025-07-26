package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/deshchain/deshchain/x/identity/types"
)

// SetUserLanguage sets the user's preferred language
func (k msgServer) SetUserLanguage(goCtx context.Context, msg *types.MsgSetUserLanguage) (*types.MsgSetUserLanguageResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate that the signer has permission to set language for this user
	// In a production system, you might want additional authorization checks
	if err := k.validateUserPermission(ctx, msg.Signer, msg.UserDid); err != nil {
		return nil, err
	}

	// Set the user's language preference
	langCode := types.LanguageCode(msg.LanguageCode)
	if err := k.i18nKeeper.SetUserLanguagePreference(ctx, msg.UserDid, langCode); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to set user language preference")
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUserLanguageSet,
			sdk.NewAttribute(types.AttributeKeyUserDID, msg.UserDid),
			sdk.NewAttribute(types.AttributeKeyLanguageCode, msg.LanguageCode),
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer),
		),
	)

	return &types.MsgSetUserLanguageResponse{
		Success: true,
		Message: k.i18nKeeper.GetLocalizedMessage(ctx, types.MsgLanguageChanged, langCode),
	}, nil
}

// AddCustomMessage adds a custom localized message
func (k msgServer) AddCustomMessage(goCtx context.Context, msg *types.MsgAddCustomMessage) (*types.MsgAddCustomMessageResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the signer has admin permissions for adding custom messages
	if err := k.validateAdminPermission(ctx, msg.Signer); err != nil {
		return nil, err
	}

	// Add the custom message
	if err := k.i18nKeeper.AddCustomMessage(ctx, msg.Key, msg.Category, msg.Description, msg.Text); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to add custom message")
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCustomMessageAdded,
			sdk.NewAttribute(types.AttributeKeyMessageKey, msg.Key),
			sdk.NewAttribute(types.AttributeKeyCategory, msg.Category),
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer),
		),
	)

	return &types.MsgAddCustomMessageResponse{
		Success: true,
		Message: "Custom message added successfully",
	}, nil
}

// UpdateLocalizationConfig updates the localization configuration
func (k msgServer) UpdateLocalizationConfig(goCtx context.Context, msg *types.MsgUpdateLocalizationConfig) (*types.MsgUpdateLocalizationConfigResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the signer has admin permissions for updating configuration
	if err := k.validateAdminPermission(ctx, msg.Signer); err != nil {
		return nil, err
	}

	// Update the localization configuration
	if err := k.i18nKeeper.SetLocalizationConfig(ctx, msg.Config); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to update localization configuration")
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeLocalizationConfigUpdated,
			sdk.NewAttribute(types.AttributeKeyDefaultLanguage, string(msg.Config.DefaultLanguage)),
			sdk.NewAttribute(types.AttributeKeyRegion, msg.Config.Region),
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer),
		),
	)

	return &types.MsgUpdateLocalizationConfigResponse{
		Success: true,
		Message: "Localization configuration updated successfully",
	}, nil
}

// ImportMessages imports messages from a catalog
func (k msgServer) ImportMessages(goCtx context.Context, msg *types.MsgImportMessages) (*types.MsgImportMessagesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the signer has admin permissions for importing messages
	if err := k.validateAdminPermission(ctx, msg.Signer); err != nil {
		return nil, err
	}

	// Import the messages
	if err := k.i18nKeeper.ImportMessages(ctx, msg.Catalog); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to import messages")
	}

	// Count imported messages
	importedCount := len(msg.Catalog.Messages)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMessagesImported,
			sdk.NewAttribute(types.AttributeKeyMessageCount, string(rune(importedCount))),
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer),
		),
	)

	return &types.MsgImportMessagesResponse{
		Success:        true,
		Message:        "Messages imported successfully",
		ImportedCount:  int32(importedCount),
	}, nil
}

// Helper functions

// validateUserPermission validates that the signer has permission to modify user settings
func (k msgServer) validateUserPermission(ctx sdk.Context, signer, userDID string) error {
	// Convert signer address
	signerAddr, err := sdk.AccAddressFromBech32(signer)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrInvalidAddress, "invalid signer address: %v", err)
	}

	// Check if the signer is the user themselves
	// In a more sophisticated system, you might have delegation or guardian permissions
	identity, err := k.GetIdentityByDID(ctx, userDID)
	if err != nil {
		// If identity doesn't exist yet, allow the signer to set up their own identity
		if types.ErrIdentityNotFound.Is(err) {
			return nil
		}
		return err
	}

	// Check if the signer owns this identity
	if identity.Controller != signer {
		return sdkerrors.Wrap(types.ErrUnauthorized, "signer is not authorized to modify this user's settings")
	}

	return nil
}

// validateAdminPermission validates that the signer has admin permissions
func (k msgServer) validateAdminPermission(ctx sdk.Context, signer string) error {
	// Convert signer address
	signerAddr, err := sdk.AccAddressFromBech32(signer)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrInvalidAddress, "invalid signer address: %v", err)
	}

	// Check if the signer has admin role
	// This could be implemented through various mechanisms:
	// 1. Governance-controlled admin list
	// 2. Multi-signature admin accounts
	// 3. Module parameters with admin addresses
	
	// For now, we'll use module parameters to define admin addresses
	params := k.GetParams(ctx)
	if params.AdminAddresses == nil {
		return sdkerrors.Wrap(types.ErrUnauthorized, "no admin addresses configured")
	}

	// Check if the signer is in the admin list
	for _, adminAddr := range params.AdminAddresses {
		if adminAddr == signer {
			return nil
		}
	}

	return sdkerrors.Wrap(types.ErrUnauthorized, "signer is not authorized for admin operations")
}

// Helper to get identity by DID (this would typically be in the main keeper)
func (k msgServer) GetIdentityByDID(ctx sdk.Context, did string) (*types.Identity, error) {
	store := ctx.KVStore(k.storeKey)
	
	// Get identity by DID (simplified implementation)
	indexKey := types.GetIdentityByDIDIndexKey(did)
	bz := store.Get(indexKey)
	if bz == nil {
		return nil, types.ErrIdentityNotFound
	}

	// The index stores the identity address, now get the full identity
	identityKey := types.GetIdentityKey(string(bz))
	identityBz := store.Get(identityKey)
	if identityBz == nil {
		return nil, types.ErrIdentityNotFound
	}

	var identity types.Identity
	if err := k.cdc.Unmarshal(identityBz, &identity); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to unmarshal identity")
	}

	return &identity, nil
}

// GetParams retrieves module parameters (this would typically be in the main keeper)
func (k msgServer) GetParams(ctx sdk.Context) types.Params {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ModuleParamsPrefix)
	if bz == nil {
		return types.DefaultParams()
	}

	var params types.Params
	if err := k.cdc.Unmarshal(bz, &params); err != nil {
		return types.DefaultParams()
	}

	return params
}