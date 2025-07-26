package keeper

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/deshchain/x/identity/types"
)

// SetDIDDocument stores a DID document
func (k Keeper) SetDIDDocument(ctx sdk.Context, didDoc types.DIDDocument) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DIDDocumentPrefix)
	b := k.cdc.MustMarshal(&didDoc)
	store.Set([]byte(didDoc.ID), b)
	
	// Store version history
	k.SetDIDVersion(ctx, didDoc.ID, didDoc.Updated.Format(time.RFC3339), didDoc)
}

// GetDIDDocument retrieves a DID document
func (k Keeper) GetDIDDocument(ctx sdk.Context, did string) (types.DIDDocument, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DIDDocumentPrefix)
	b := store.Get([]byte(did))
	if b == nil {
		return types.DIDDocument{}, false
	}
	
	var didDoc types.DIDDocument
	k.cdc.MustUnmarshal(b, &didDoc)
	return didDoc, true
}

// HasDIDDocument checks if a DID document exists
func (k Keeper) HasDIDDocument(ctx sdk.Context, did string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DIDDocumentPrefix)
	return store.Has([]byte(did))
}

// DeleteDIDDocument removes a DID document
func (k Keeper) DeleteDIDDocument(ctx sdk.Context, did string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DIDDocumentPrefix)
	store.Delete([]byte(did))
}

// SetDIDVersion stores a specific version of a DID document
func (k Keeper) SetDIDVersion(ctx sdk.Context, did string, version string, didDoc types.DIDDocument) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DIDVersionPrefix)
	b := k.cdc.MustMarshal(&didDoc)
	store.Set(types.GetDIDVersionKey(did, version), b)
}

// GetDIDVersion retrieves a specific version of a DID document
func (k Keeper) GetDIDVersion(ctx sdk.Context, did string, version string) (types.DIDDocument, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DIDVersionPrefix)
	b := store.Get(types.GetDIDVersionKey(did, version))
	if b == nil {
		return types.DIDDocument{}, false
	}
	
	var didDoc types.DIDDocument
	k.cdc.MustUnmarshal(b, &didDoc)
	return didDoc, true
}

// IterateDIDDocuments iterates over all DID documents
func (k Keeper) IterateDIDDocuments(ctx sdk.Context, cb func(didDoc types.DIDDocument) bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DIDDocumentPrefix)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var didDoc types.DIDDocument
		k.cdc.MustUnmarshal(iterator.Value(), &didDoc)
		if cb(didDoc) {
			break
		}
	}
}

// ResolveDID resolves a DID to its document with metadata
func (k Keeper) ResolveDID(ctx sdk.Context, did string) (*types.DIDResolutionResult, error) {
	// Parse DID
	method, identifier, err := types.ParseDID(did)
	if err != nil {
		return nil, types.ErrInvalidDID
	}
	
	// Check if method is supported
	params := k.GetParams(ctx)
	if !params.IsDIDMethodSupported(method) {
		return nil, types.ErrDIDMethodNotSupported
	}
	
	// Get DID document
	didDoc, found := k.GetDIDDocument(ctx, did)
	if !found {
		return &types.DIDResolutionResult{
			DidDocument: nil,
			DidResolutionMetadata: &types.DIDResolutionMetadata{
				ResolutionTime: ctx.BlockTime(),
				Error:          "notFound",
				ErrorMessage:   "DID not found",
			},
			DidDocumentMetadata: nil,
		}, nil
	}
	
	// Check if DID is deactivated
	deactivated := false
	if didDoc.Metadata != nil {
		if deact, ok := didDoc.Metadata["deactivated"].(bool); ok {
			deactivated = deact
		}
	}
	
	// Build resolution result
	result := &types.DIDResolutionResult{
		DidDocument: &didDoc,
		DidResolutionMetadata: &types.DIDResolutionMetadata{
			ContentType:    "application/did+json",
			ResolutionTime: ctx.BlockTime(),
		},
		DidDocumentMetadata: &types.DIDDocumentMetadata{
			Created:     didDoc.Created,
			Updated:     didDoc.Updated,
			Deactivated: deactivated,
			VersionID:   didDoc.Updated.Format(time.RFC3339),
		},
	}
	
	return result, nil
}

// ValidateDIDRegistration validates if a DID can be registered
func (k Keeper) ValidateDIDRegistration(ctx sdk.Context, did string, creator string) error {
	// Check if DID already exists
	if k.HasDIDDocument(ctx, did) {
		return types.ErrDIDAlreadyExists
	}
	
	// Parse DID
	method, identifier, err := types.ParseDID(did)
	if err != nil {
		return types.ErrInvalidDID
	}
	
	// Check if method is supported
	params := k.GetParams(ctx)
	if !params.IsDIDMethodSupported(method) {
		return types.ErrDIDMethodNotSupported
	}
	
	// For did:desh method, verify the identifier matches creator
	if method == types.DIDMethod {
		// The identifier should be derived from the creator address
		expectedIdentifier := sdk.MustAccAddressFromBech32(creator).String()
		if identifier != expectedIdentifier {
			return types.ErrUnauthorized
		}
	}
	
	return nil
}

// AddVerificationMethod adds a verification method to a DID document
func (k Keeper) AddVerificationMethod(ctx sdk.Context, did string, vm types.VerificationMethod) error {
	didDoc, found := k.GetDIDDocument(ctx, did)
	if !found {
		return types.ErrDIDNotFound
	}
	
	// Check if verification method ID already exists
	for _, existingVM := range didDoc.VerificationMethod {
		if existingVM.ID == vm.ID {
			return types.ErrInvalidRequest
		}
	}
	
	didDoc.AddVerificationMethod(vm)
	k.SetDIDDocument(ctx, didDoc)
	
	return nil
}

// RemoveVerificationMethod removes a verification method from a DID document
func (k Keeper) RemoveVerificationMethod(ctx sdk.Context, did string, vmID string) error {
	didDoc, found := k.GetDIDDocument(ctx, did)
	if !found {
		return types.ErrDIDNotFound
	}
	
	// Find and remove verification method
	newVMs := []types.VerificationMethod{}
	found = false
	for _, vm := range didDoc.VerificationMethod {
		if vm.ID != vmID {
			newVMs = append(newVMs, vm)
		} else {
			found = true
		}
	}
	
	if !found {
		return types.ErrInvalidRequest
	}
	
	// Ensure at least one verification method remains
	if len(newVMs) == 0 {
		return types.ErrInvalidRequest
	}
	
	didDoc.VerificationMethod = newVMs
	didDoc.Updated = ctx.BlockTime()
	k.SetDIDDocument(ctx, didDoc)
	
	return nil
}

// AddService adds a service endpoint to a DID document
func (k Keeper) AddService(ctx sdk.Context, did string, service types.Service) error {
	didDoc, found := k.GetDIDDocument(ctx, did)
	if !found {
		return types.ErrDIDNotFound
	}
	
	// Check if service ID already exists
	for _, existingSvc := range didDoc.Service {
		if existingSvc.ID == service.ID {
			return types.ErrInvalidRequest
		}
	}
	
	didDoc.AddService(service)
	k.SetDIDDocument(ctx, didDoc)
	
	return nil
}

// RemoveService removes a service endpoint from a DID document
func (k Keeper) RemoveService(ctx sdk.Context, did string, serviceID string) error {
	didDoc, found := k.GetDIDDocument(ctx, did)
	if !found {
		return types.ErrDIDNotFound
	}
	
	// Find and remove service
	newServices := []types.Service{}
	found = false
	for _, svc := range didDoc.Service {
		if svc.ID != serviceID {
			newServices = append(newServices, svc)
		} else {
			found = true
		}
	}
	
	if !found {
		return types.ErrInvalidRequest
	}
	
	didDoc.Service = newServices
	didDoc.Updated = ctx.BlockTime()
	k.SetDIDDocument(ctx, didDoc)
	
	return nil
}

// DeactivateDID deactivates a DID document
func (k Keeper) DeactivateDID(ctx sdk.Context, did string) error {
	didDoc, found := k.GetDIDDocument(ctx, did)
	if !found {
		return types.ErrDIDNotFound
	}
	
	// Mark as deactivated in metadata
	if didDoc.Metadata == nil {
		didDoc.Metadata = make(map[string]interface{})
	}
	didDoc.Metadata["deactivated"] = true
	didDoc.Metadata["deactivatedAt"] = ctx.BlockTime().Format(time.RFC3339)
	didDoc.Updated = ctx.BlockTime()
	
	k.SetDIDDocument(ctx, didDoc)
	
	return nil
}

// GetDIDController returns the controller of a DID
func (k Keeper) GetDIDController(ctx sdk.Context, did string) (string, error) {
	didDoc, found := k.GetDIDDocument(ctx, did)
	if !found {
		return "", types.ErrDIDNotFound
	}
	
	if didDoc.Controller != "" {
		return didDoc.Controller, nil
	}
	
	// If no explicit controller, the DID controls itself
	return did, nil
}

// ValidateDIDDocumentSize validates the size of a DID document
func (k Keeper) ValidateDIDDocumentSize(ctx sdk.Context, didDoc *types.DIDDocument) error {
	bz, err := didDoc.ToJSON()
	if err != nil {
		return err
	}
	
	maxSize := k.MaxDIDDocumentSize(ctx)
	if uint64(len(bz)) > maxSize {
		return types.ErrDataTooLarge
	}
	
	return nil
}