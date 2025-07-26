package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/deshchain/x/identity/types"
)

// Service Registry Functions

// RegisterService registers a new service
func (k Keeper) RegisterService(ctx sdk.Context, service types.ServiceRegistryEntry) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ServiceRegistryPrefix)
	b := k.cdc.MustMarshal(&service)
	store.Set([]byte(service.ID), b)
	
	// Create type index
	k.SetServiceTypeIndex(ctx, service.Type, service.ID)
}

// GetService retrieves a service by ID
func (k Keeper) GetService(ctx sdk.Context, serviceID string) (types.ServiceRegistryEntry, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ServiceRegistryPrefix)
	b := store.Get([]byte(serviceID))
	if b == nil {
		return types.ServiceRegistryEntry{}, false
	}
	
	var service types.ServiceRegistryEntry
	k.cdc.MustUnmarshal(b, &service)
	return service, true
}

// UpdateService updates a service entry
func (k Keeper) UpdateService(ctx sdk.Context, service types.ServiceRegistryEntry) error {
	if !k.HasService(ctx, service.ID) {
		return types.ErrServiceNotFound
	}
	
	k.RegisterService(ctx, service)
	return nil
}

// HasService checks if a service exists
func (k Keeper) HasService(ctx sdk.Context, serviceID string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ServiceRegistryPrefix)
	return store.Has([]byte(serviceID))
}

// DeactivateService deactivates a service
func (k Keeper) DeactivateService(ctx sdk.Context, serviceID string) error {
	service, found := k.GetService(ctx, serviceID)
	if !found {
		return types.ErrServiceNotFound
	}
	
	service.IsActive = false
	k.RegisterService(ctx, service)
	
	return nil
}

// SetServiceTypeIndex creates an index for services by type
func (k Keeper) SetServiceTypeIndex(ctx sdk.Context, serviceType string, serviceID string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ServiceIndexPrefix)
	store.Set(types.GetServiceByTypeIndexKey(serviceType, serviceID), []byte{1})
}

// RemoveServiceTypeIndex removes a service from type index
func (k Keeper) RemoveServiceTypeIndex(ctx sdk.Context, serviceType string, serviceID string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ServiceIndexPrefix)
	store.Delete(types.GetServiceByTypeIndexKey(serviceType, serviceID))
}

// GetServicesByType returns all services of a specific type
func (k Keeper) GetServicesByType(ctx sdk.Context, serviceType string) []types.ServiceRegistryEntry {
	var services []types.ServiceRegistryEntry
	
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ServiceIndexPrefix)
	iterator := sdk.KVStorePrefixIterator(store, append([]byte("type:"), []byte(serviceType)...))
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		// Extract service ID from the key
		key := string(iterator.Key())
		serviceID := key[len("type:"+serviceType):]
		
		if service, found := k.GetService(ctx, serviceID); found {
			services = append(services, service)
		}
	}
	
	return services
}

// GetActiveServices returns all active services
func (k Keeper) GetActiveServices(ctx sdk.Context) []types.ServiceRegistryEntry {
	var activeServices []types.ServiceRegistryEntry
	
	k.IterateServiceRegistry(ctx, func(service types.ServiceRegistryEntry) bool {
		if service.IsActive {
			activeServices = append(activeServices, service)
		}
		return false
	})
	
	return activeServices
}

// IterateServiceRegistry iterates over all services
func (k Keeper) IterateServiceRegistry(ctx sdk.Context, cb func(service types.ServiceRegistryEntry) bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ServiceRegistryPrefix)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var service types.ServiceRegistryEntry
		k.cdc.MustUnmarshal(iterator.Value(), &service)
		if cb(service) {
			break
		}
	}
}

// IsServiceAuthorized checks if a service is authorized for an operation
func (k Keeper) IsServiceAuthorized(ctx sdk.Context, serviceID string, operation string) bool {
	service, found := k.GetService(ctx, serviceID)
	if !found || !service.IsActive {
		return false
	}
	
	for _, allowedOp := range service.AllowedOperations {
		if allowedOp == operation || allowedOp == "*" {
			return true
		}
	}
	
	return false
}

// Issuer Registry Functions

// RegisterIssuer registers a new credential issuer
func (k Keeper) RegisterIssuer(ctx sdk.Context, issuer types.IssuerRegistryEntry) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.IssuerRegistryPrefix)
	b := k.cdc.MustMarshal(&issuer)
	store.Set([]byte(issuer.DID), b)
	
	// Create type index
	k.SetIssuerTypeIndex(ctx, issuer.Type, issuer.DID)
}

// GetIssuer retrieves an issuer by DID
func (k Keeper) GetIssuer(ctx sdk.Context, issuerDID string) (types.IssuerRegistryEntry, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.IssuerRegistryPrefix)
	b := store.Get([]byte(issuerDID))
	if b == nil {
		return types.IssuerRegistryEntry{}, false
	}
	
	var issuer types.IssuerRegistryEntry
	k.cdc.MustUnmarshal(b, &issuer)
	return issuer, true
}

// UpdateIssuer updates an issuer entry
func (k Keeper) UpdateIssuer(ctx sdk.Context, issuer types.IssuerRegistryEntry) error {
	if !k.HasIssuer(ctx, issuer.DID) {
		return types.ErrInvalidIssuer
	}
	
	k.RegisterIssuer(ctx, issuer)
	return nil
}

// HasIssuer checks if an issuer exists
func (k Keeper) HasIssuer(ctx sdk.Context, issuerDID string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.IssuerRegistryPrefix)
	return store.Has([]byte(issuerDID))
}

// DeactivateIssuer deactivates an issuer
func (k Keeper) DeactivateIssuer(ctx sdk.Context, issuerDID string) error {
	issuer, found := k.GetIssuer(ctx, issuerDID)
	if !found {
		return types.ErrInvalidIssuer
	}
	
	issuer.IsActive = false
	k.RegisterIssuer(ctx, issuer)
	
	return nil
}

// SetIssuerTypeIndex creates an index for issuers by type
func (k Keeper) SetIssuerTypeIndex(ctx sdk.Context, issuerType string, issuerDID string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.IssuerIndexPrefix)
	key := append(append([]byte("type:"), []byte(issuerType)...), []byte(issuerDID)...)
	store.Set(key, []byte{1})
}

// RemoveIssuerTypeIndex removes an issuer from type index
func (k Keeper) RemoveIssuerTypeIndex(ctx sdk.Context, issuerType string, issuerDID string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.IssuerIndexPrefix)
	key := append(append([]byte("type:"), []byte(issuerType)...), []byte(issuerDID)...)
	store.Delete(key)
}

// GetIssuersByType returns all issuers of a specific type
func (k Keeper) GetIssuersByType(ctx sdk.Context, issuerType string) []types.IssuerRegistryEntry {
	var issuers []types.IssuerRegistryEntry
	
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.IssuerIndexPrefix)
	iterator := sdk.KVStorePrefixIterator(store, append([]byte("type:"), []byte(issuerType)...))
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		// Extract issuer DID from the key
		key := string(iterator.Key())
		issuerDID := key[len("type:"+issuerType):]
		
		if issuer, found := k.GetIssuer(ctx, issuerDID); found {
			issuers = append(issuers, issuer)
		}
	}
	
	return issuers
}

// GetActiveIssuers returns all active issuers
func (k Keeper) GetActiveIssuers(ctx sdk.Context) []types.IssuerRegistryEntry {
	var activeIssuers []types.IssuerRegistryEntry
	
	k.IterateIssuerRegistry(ctx, func(issuer types.IssuerRegistryEntry) bool {
		if issuer.IsActive {
			activeIssuers = append(activeIssuers, issuer)
		}
		return false
	})
	
	return activeIssuers
}

// IterateIssuerRegistry iterates over all issuers
func (k Keeper) IterateIssuerRegistry(ctx sdk.Context, cb func(issuer types.IssuerRegistryEntry) bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.IssuerRegistryPrefix)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var issuer types.IssuerRegistryEntry
		k.cdc.MustUnmarshal(iterator.Value(), &issuer)
		if cb(issuer) {
			break
		}
	}
}

// IsRegisteredIssuer checks if an issuer is registered and active
func (k Keeper) IsRegisteredIssuer(ctx sdk.Context, issuerDID string) bool {
	issuer, found := k.GetIssuer(ctx, issuerDID)
	return found && issuer.IsActive
}

// IsIssuerAuthorizedForType checks if an issuer can issue a specific credential type
func (k Keeper) IsIssuerAuthorizedForType(ctx sdk.Context, issuerDID string, credentialType string) bool {
	issuer, found := k.GetIssuer(ctx, issuerDID)
	if !found || !issuer.IsActive {
		return false
	}
	
	for _, allowedType := range issuer.AllowedCredentialTypes {
		if allowedType == credentialType || allowedType == "*" {
			return true
		}
	}
	
	return false
}

// UpdateIssuerTrustLevel updates the trust level of an issuer
func (k Keeper) UpdateIssuerTrustLevel(ctx sdk.Context, issuerDID string, newTrustLevel int32) error {
	issuer, found := k.GetIssuer(ctx, issuerDID)
	if !found {
		return types.ErrInvalidIssuer
	}
	
	if newTrustLevel < 0 || newTrustLevel > 100 {
		return types.ErrInvalidRequest
	}
	
	issuer.TrustLevel = newTrustLevel
	k.RegisterIssuer(ctx, issuer)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"issuer_trust_updated",
			sdk.NewAttribute("issuer_did", issuerDID),
			sdk.NewAttribute("new_trust_level", fmt.Sprintf("%d", newTrustLevel)),
		),
	)
	
	return nil
}

// GetHighTrustIssuers returns issuers with trust level above threshold
func (k Keeper) GetHighTrustIssuers(ctx sdk.Context, minTrustLevel int32) []types.IssuerRegistryEntry {
	var highTrustIssuers []types.IssuerRegistryEntry
	
	k.IterateIssuerRegistry(ctx, func(issuer types.IssuerRegistryEntry) bool {
		if issuer.IsActive && issuer.TrustLevel >= minTrustLevel {
			highTrustIssuers = append(highTrustIssuers, issuer)
		}
		return false
	})
	
	return highTrustIssuers
}