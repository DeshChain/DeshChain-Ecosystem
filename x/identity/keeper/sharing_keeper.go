package keeper

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/DeshChain/DeshChain-Ecosystem/x/identity/types"
)

// Cross-Module Identity Sharing Protocol Implementation

// CreateShareRequest creates a new identity sharing request
func (k Keeper) CreateShareRequest(
	ctx sdk.Context,
	requesterModule string,
	providerModule string,
	holderDID string,
	requestedData []types.DataRequest,
	purpose string,
	justification string,
	ttl time.Duration,
) (*types.IdentityShareRequest, error) {
	// Validate identity exists
	_, found := k.GetIdentity(ctx, holderDID)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrIdentityNotFound, "holder DID: %s", holderDID)
	}

	// Validate module capabilities
	if err := k.validateModuleCapabilities(ctx, requesterModule, providerModule, requestedData); err != nil {
		return nil, err
	}

	// Generate unique request ID
	requestID, err := k.generateRequestID()
	if err != nil {
		return nil, err
	}

	// Create share request
	request := &types.IdentityShareRequest{
		RequestID:       requestID,
		RequesterModule: requesterModule,
		ProviderModule:  providerModule,
		HolderDID:       holderDID,
		RequestedData:   requestedData,
		Purpose:         purpose,
		Justification:   justification,
		TTL:             ttl,
		RequestedAt:     ctx.BlockTime(),
		Status:          types.ShareRequestStatus_PENDING,
		Metadata:        make(map[string]interface{}),
	}

	// Validate request
	if err := request.ValidateBasic(); err != nil {
		return nil, err
	}

	// Check if request can be auto-approved
	if agreement, found := k.getActiveAgreement(ctx, requesterModule, providerModule); found {
		if agreement.CanAutoApprove(request) {
			request.Status = types.ShareRequestStatus_APPROVED
		}
	}

	// Check access policies
	if err := k.checkAccessPolicies(ctx, request); err != nil {
		return nil, err
	}

	// Store the request
	k.SetShareRequest(ctx, *request)

	// Create audit log
	k.logShareAudit(ctx, request.RequestID, holderDID, requesterModule, providerModule,
		types.ShareAuditAction_REQUEST_CREATED, []string{}, purpose)

	return request, nil
}

// ApproveShareRequest approves a pending share request
func (k Keeper) ApproveShareRequest(
	ctx sdk.Context,
	authority sdk.AccAddress,
	requestID string,
	accessToken string,
) error {
	request, found := k.GetShareRequest(ctx, requestID)
	if !found {
		return types.ErrShareRequestNotFound
	}

	if request.IsExpired() {
		return types.ErrShareRequestExpired
	}

	if request.Status != types.ShareRequestStatus_PENDING {
		return sdkerrors.Wrapf(types.ErrInvalidShareRequest, "request status: %v", request.Status)
	}

	// Validate authority (could be holder or authorized party)
	if err := k.validateShareAuthority(ctx, authority, request.HolderDID); err != nil {
		return err
	}

	// Update request status
	request.Status = types.ShareRequestStatus_APPROVED
	k.SetShareRequest(ctx, request)

	// Create response with shared data
	response := &types.IdentityShareResponse{
		RequestID:   requestID,
		Status:      types.ShareRequestStatus_APPROVED,
		ExpiresAt:   ctx.BlockTime().Add(request.TTL),
		ResponseAt:  ctx.BlockTime(),
		AccessToken: accessToken,
		Metadata:    make(map[string]interface{}),
	}

	// Collect and prepare shared data
	sharedData, err := k.collectSharedData(ctx, &request)
	if err != nil {
		return err
	}
	response.SharedData = sharedData

	// Store the response
	k.SetShareResponse(ctx, *response)

	// Create audit log
	sharedDataTypes := make([]string, len(sharedData))
	for i, data := range sharedData {
		sharedDataTypes[i] = data.CredentialType
	}
	k.logShareAudit(ctx, requestID, request.HolderDID, request.RequesterModule, request.ProviderModule,
		types.ShareAuditAction_REQUEST_APPROVED, sharedDataTypes, request.Purpose)

	return nil
}

// DenyShareRequest denies a pending share request
func (k Keeper) DenyShareRequest(
	ctx sdk.Context,
	authority sdk.AccAddress,
	requestID string,
	denialReason string,
) error {
	request, found := k.GetShareRequest(ctx, requestID)
	if !found {
		return types.ErrShareRequestNotFound
	}

	if request.Status != types.ShareRequestStatus_PENDING {
		return sdkerrors.Wrapf(types.ErrInvalidShareRequest, "request status: %v", request.Status)
	}

	// Validate authority
	if err := k.validateShareAuthority(ctx, authority, request.HolderDID); err != nil {
		return err
	}

	// Update request status
	request.Status = types.ShareRequestStatus_DENIED
	k.SetShareRequest(ctx, request)

	// Create response
	response := &types.IdentityShareResponse{
		RequestID:    requestID,
		Status:       types.ShareRequestStatus_DENIED,
		DenialReason: denialReason,
		ResponseAt:   ctx.BlockTime(),
		Metadata:     make(map[string]interface{}),
	}

	k.SetShareResponse(ctx, *response)

	// Create audit log
	k.logShareAudit(ctx, requestID, request.HolderDID, request.RequesterModule, request.ProviderModule,
		types.ShareAuditAction_REQUEST_DENIED, []string{}, request.Purpose)

	return nil
}

// GetSharedData retrieves shared data for an approved request
func (k Keeper) GetSharedData(
	ctx sdk.Context,
	requesterModule string,
	requestID string,
	accessToken string,
) (*types.IdentityShareResponse, error) {
	request, found := k.GetShareRequest(ctx, requestID)
	if !found {
		return nil, types.ErrShareRequestNotFound
	}

	response, found := k.GetShareResponse(ctx, requestID)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrShareRequestNotFound, "response not found")
	}

	// Validate requester
	if request.RequesterModule != requesterModule {
		return nil, types.ErrUnauthorizedModule
	}

	// Validate access token
	if response.AccessToken != accessToken {
		return nil, sdkerrors.Wrap(types.ErrUnauthorizedModule, "invalid access token")
	}

	// Check if response is expired
	if ctx.BlockTime().After(response.ExpiresAt) {
		return nil, types.ErrShareRequestExpired
	}

	// Create audit log for data access
	sharedDataTypes := make([]string, len(response.SharedData))
	for i, data := range response.SharedData {
		sharedDataTypes[i] = data.CredentialType
	}
	k.logShareAudit(ctx, requestID, request.HolderDID, request.RequesterModule, request.ProviderModule,
		types.ShareAuditAction_DATA_ACCESSED, sharedDataTypes, request.Purpose)

	return &response, nil
}

// CreateSharingAgreement creates a standing agreement between modules
func (k Keeper) CreateSharingAgreement(
	ctx sdk.Context,
	authority sdk.AccAddress,
	requesterModule string,
	providerModule string,
	allowedDataTypes []string,
	purposes []string,
	trustLevel string,
	autoApprove bool,
	maxTTL time.Duration,
	validityPeriod time.Duration,
) (*types.IdentityShareAgreement, error) {
	// Generate agreement ID
	agreementID, err := k.generateAgreementID()
	if err != nil {
		return nil, err
	}

	// Create agreement
	agreement := &types.IdentityShareAgreement{
		AgreementID:      agreementID,
		RequesterModule:  requesterModule,
		ProviderModule:   providerModule,
		AllowedDataTypes: allowedDataTypes,
		Purposes:         purposes,
		TrustLevel:       trustLevel,
		AutoApprove:      autoApprove,
		MaxTTL:           maxTTL,
		CreatedAt:        ctx.BlockTime(),
		ExpiresAt:        ctx.BlockTime().Add(validityPeriod),
		Status:           types.AgreementStatus_ACTIVE,
		Metadata:         make(map[string]interface{}),
	}

	// Validate agreement
	if err := agreement.ValidateBasic(); err != nil {
		return nil, err
	}

	// Store agreement
	k.SetSharingAgreement(ctx, *agreement)

	return agreement, nil
}

// CreateAccessPolicy creates an access policy for a holder
func (k Keeper) CreateAccessPolicy(
	ctx sdk.Context,
	holderDID string,
	allowedModules []string,
	deniedModules []string,
	dataRestrictions map[string][]string,
	purposeRestrictions []string,
	timeRestrictions types.TimeRestriction,
	maxSharesPerDay int,
	requireExplicitConsent bool,
) (*types.IdentityAccessPolicy, error) {
	// Validate identity exists
	_, found := k.GetIdentity(ctx, holderDID)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrIdentityNotFound, "holder DID: %s", holderDID)
	}

	// Generate policy ID
	policyID, err := k.generatePolicyID()
	if err != nil {
		return nil, err
	}

	// Create policy
	policy := &types.IdentityAccessPolicy{
		PolicyID:               policyID,
		HolderDID:              holderDID,
		AllowedModules:         allowedModules,
		DeniedModules:          deniedModules,
		DataRestrictions:       dataRestrictions,
		PurposeRestrictions:    purposeRestrictions,
		TimeRestrictions:       timeRestrictions,
		MaxSharesPerDay:        maxSharesPerDay,
		RequireExplicitConsent: requireExplicitConsent,
		CreatedAt:              ctx.BlockTime(),
		UpdatedAt:              ctx.BlockTime(),
		Metadata:               make(map[string]interface{}),
	}

	// Validate policy
	if err := policy.ValidateBasic(); err != nil {
		return nil, err
	}

	// Store policy
	k.SetAccessPolicy(ctx, *policy)

	return policy, nil
}

// Storage functions

// SetShareRequest stores a share request
func (k Keeper) SetShareRequest(ctx sdk.Context, request types.IdentityShareRequest) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&request)
	store.Set(types.ShareRequestKey(request.RequestID), bz)
}

// GetShareRequest retrieves a share request
func (k Keeper) GetShareRequest(ctx sdk.Context, requestID string) (types.IdentityShareRequest, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ShareRequestKey(requestID))
	if bz == nil {
		return types.IdentityShareRequest{}, false
	}

	var request types.IdentityShareRequest
	k.cdc.MustUnmarshal(bz, &request)
	return request, true
}

// SetShareResponse stores a share response
func (k Keeper) SetShareResponse(ctx sdk.Context, response types.IdentityShareResponse) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&response)
	store.Set(types.ShareResponseKey(response.RequestID), bz)
}

// GetShareResponse retrieves a share response
func (k Keeper) GetShareResponse(ctx sdk.Context, requestID string) (types.IdentityShareResponse, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ShareResponseKey(requestID))
	if bz == nil {
		return types.IdentityShareResponse{}, false
	}

	var response types.IdentityShareResponse
	k.cdc.MustUnmarshal(bz, &response)
	return response, true
}

// SetSharingAgreement stores a sharing agreement
func (k Keeper) SetSharingAgreement(ctx sdk.Context, agreement types.IdentityShareAgreement) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&agreement)
	store.Set(types.SharingAgreementKey(agreement.AgreementID), bz)
}

// GetSharingAgreement retrieves a sharing agreement
func (k Keeper) GetSharingAgreement(ctx sdk.Context, agreementID string) (types.IdentityShareAgreement, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.SharingAgreementKey(agreementID))
	if bz == nil {
		return types.IdentityShareAgreement{}, false
	}

	var agreement types.IdentityShareAgreement
	k.cdc.MustUnmarshal(bz, &agreement)
	return agreement, true
}

// SetAccessPolicy stores an access policy
func (k Keeper) SetAccessPolicy(ctx sdk.Context, policy types.IdentityAccessPolicy) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&policy)
	store.Set(types.AccessPolicyKey(policy.PolicyID), bz)
	
	// Also index by holder DID for efficient lookups
	store.Set(types.AccessPolicyByHolderKey(policy.HolderDID, policy.PolicyID), []byte(policy.PolicyID))
}

// GetAccessPolicy retrieves an access policy
func (k Keeper) GetAccessPolicy(ctx sdk.Context, policyID string) (types.IdentityAccessPolicy, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.AccessPolicyKey(policyID))
	if bz == nil {
		return types.IdentityAccessPolicy{}, false
	}

	var policy types.IdentityAccessPolicy
	k.cdc.MustUnmarshal(bz, &policy)
	return policy, true
}

// GetAccessPoliciesByHolder retrieves all access policies for a holder
func (k Keeper) GetAccessPoliciesByHolder(ctx sdk.Context, holderDID string) []types.IdentityAccessPolicy {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.AccessPolicyByHolderPrefix(holderDID))
	defer iterator.Close()

	var policies []types.IdentityAccessPolicy
	for ; iterator.Valid(); iterator.Next() {
		policyID := string(iterator.Value())
		policy, found := k.GetAccessPolicy(ctx, policyID)
		if found {
			policies = append(policies, policy)
		}
	}

	return policies
}

// Helper functions

// generateRequestID generates a unique request ID
func (k Keeper) generateRequestID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "share_req_" + hex.EncodeToString(bytes), nil
}

// generateAgreementID generates a unique agreement ID
func (k Keeper) generateAgreementID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "share_agr_" + hex.EncodeToString(bytes), nil
}

// generatePolicyID generates a unique policy ID
func (k Keeper) generatePolicyID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "access_pol_" + hex.EncodeToString(bytes), nil
}

// validateModuleCapabilities validates that modules can request/provide the specified data
func (k Keeper) validateModuleCapabilities(ctx sdk.Context, requesterModule, providerModule string, requestedData []types.DataRequest) error {
	// Get module capabilities
	requesterCaps, found := k.GetModuleCapabilities(ctx, requesterModule)
	if !found {
		return sdkerrors.Wrapf(types.ErrUnauthorizedModule, "requester module not found: %s", requesterModule)
	}

	providerCaps, found := k.GetModuleCapabilities(ctx, providerModule)
	if !found {
		return sdkerrors.Wrapf(types.ErrUnauthorizedModule, "provider module not found: %s", providerModule)
	}

	// Validate each data request
	for _, dataReq := range requestedData {
		// Check if requester can request this data type
		canRequest := false
		for _, canReqType := range requesterCaps.CanRequest {
			if canReqType == dataReq.CredentialType {
				canRequest = true
				break
			}
		}
		if !canRequest {
			return sdkerrors.Wrapf(types.ErrUnauthorizedModule, 
				"requester %s cannot request %s", requesterModule, dataReq.CredentialType)
		}

		// Check if provider can provide this data type
		canProvide := false
		for _, canProvType := range providerCaps.CanProvide {
			if canProvType == dataReq.CredentialType {
				canProvide = true
				break
			}
		}
		if !canProvide {
			return sdkerrors.Wrapf(types.ErrUnauthorizedModule, 
				"provider %s cannot provide %s", providerModule, dataReq.CredentialType)
		}
	}

	return nil
}

// getActiveAgreement retrieves an active agreement between modules
func (k Keeper) getActiveAgreement(ctx sdk.Context, requesterModule, providerModule string) (types.IdentityShareAgreement, bool) {
	// This would iterate through agreements and find active ones
	// For now, return empty (implementation depends on indexing strategy)
	return types.IdentityShareAgreement{}, false
}

// checkAccessPolicies checks if the request complies with access policies
func (k Keeper) checkAccessPolicies(ctx sdk.Context, request *types.IdentityShareRequest) error {
	policies := k.GetAccessPoliciesByHolder(ctx, request.HolderDID)
	
	for _, policy := range policies {
		// Check denied modules
		for _, deniedModule := range policy.DeniedModules {
			if deniedModule == request.RequesterModule {
				return types.ErrAccessPolicyViolation
			}
		}
		
		// Check allowed modules (if specified)
		if len(policy.AllowedModules) > 0 {
			allowed := false
			for _, allowedModule := range policy.AllowedModules {
				if allowedModule == request.RequesterModule {
					allowed = true
					break
				}
			}
			if !allowed {
				return types.ErrAccessPolicyViolation
			}
		}
		
		// Check purpose restrictions
		if len(policy.PurposeRestrictions) > 0 {
			purposeAllowed := false
			for _, allowedPurpose := range policy.PurposeRestrictions {
				if allowedPurpose == request.Purpose {
					purposeAllowed = true
					break
				}
			}
			if !purposeAllowed {
				return types.ErrAccessPolicyViolation
			}
		}
		
		// Check daily limits
		if policy.MaxSharesPerDay > 0 {
			todayCount := k.getShareCountToday(ctx, request.HolderDID)
			if todayCount >= policy.MaxSharesPerDay {
				return types.ErrDailyLimitExceeded
			}
		}
	}
	
	return nil
}

// validateShareAuthority validates that the authority can approve/deny the request
func (k Keeper) validateShareAuthority(ctx sdk.Context, authority sdk.AccAddress, holderDID string) error {
	// Get identity
	identity, found := k.GetIdentity(ctx, holderDID)
	if !found {
		return types.ErrIdentityNotFound
	}
	
	// Check if authority is the controller
	if identity.Controller == authority.String() {
		return nil
	}
	
	// Check for delegation (this would need additional implementation)
	// For now, only controller can approve
	return sdkerrors.Wrap(types.ErrUnauthorizedModule, "only identity controller can approve")
}

// collectSharedData collects the requested data from credentials
func (k Keeper) collectSharedData(ctx sdk.Context, request *types.IdentityShareRequest) ([]types.SharedCredential, error) {
	var sharedData []types.SharedCredential
	
	for _, dataReq := range request.RequestedData {
		// Get credentials of the requested type
		credentials := k.GetCredentialsByType(ctx, request.HolderDID, dataReq.CredentialType)
		
		if len(credentials) == 0 && dataReq.Required {
			return nil, sdkerrors.Wrapf(types.ErrInvalidShareRequest, 
				"required credential type %s not found", dataReq.CredentialType)
		}
		
		// Use the most recent valid credential
		for _, cred := range credentials {
			if cred.Status != CredentialStatus_ACTIVE {
				continue
			}
			
			// Extract only requested attributes
			sharedSubject := make(map[string]interface{})
			for _, attr := range dataReq.Attributes {
				if value, exists := cred.CredentialSubject[attr]; exists {
					sharedSubject[attr] = value
				}
			}
			
			sharedCred := types.SharedCredential{
				CredentialID:   cred.Id,
				CredentialType: dataReq.CredentialType,
				Issuer:         cred.Issuer,
				IssuedAt:       cred.IssuanceDate,
				SharedData:     sharedSubject,
				TrustLevel:     "high", // This could be calculated based on issuer
			}
			
			sharedData = append(sharedData, sharedCred)
			break // Use only the first valid credential
		}
	}
	
	return sharedData, nil
}

// logShareAudit creates an audit log entry
func (k Keeper) logShareAudit(ctx sdk.Context, requestID, holderDID, requesterModule, providerModule string, action types.ShareAuditAction, sharedData []string, purpose string) {
	// Generate log ID
	logID := fmt.Sprintf("audit_%d_%s", ctx.BlockHeight(), requestID)
	
	auditLog := types.IdentityShareAuditLog{
		LogID:           logID,
		RequestID:       requestID,
		HolderDID:       holderDID,
		RequesterModule: requesterModule,
		ProviderModule:  providerModule,
		Action:          action,
		SharedData:      sharedData,
		Purpose:         purpose,
		Timestamp:       ctx.BlockTime(),
		Metadata:        make(map[string]interface{}),
	}
	
	// Store audit log
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&auditLog)
	store.Set(types.ShareAuditLogKey(logID), bz)
	
	// Index by holder DID for efficient queries
	store.Set(types.ShareAuditLogByHolderKey(holderDID, logID), []byte(logID))
}

// getShareCountToday gets the number of shares for a holder today
func (k Keeper) getShareCountToday(ctx sdk.Context, holderDID string) int {
	// This would count audit logs for today
	// For now, return 0 (implementation depends on indexing strategy)
	return 0
}

// GetModuleCapabilities retrieves module capabilities (placeholder)
func (k Keeper) GetModuleCapabilities(ctx sdk.Context, moduleName string) (types.ModuleCapabilities, bool) {
	// This would be implemented based on module registration
	// For now, return default capabilities
	defaultCaps := types.ModuleCapabilities{
		ModuleName:     moduleName,
		CanRequest:     []string{"KYCCredential", "BiometricCredential"},
		CanProvide:     []string{"KYCCredential", "BiometricCredential"},
		TrustLevel:     "medium",
		Certifications: []string{"basic"},
		Version:        "1.0.0",
		LastUpdated:    ctx.BlockTime(),
	}
	return defaultCaps, true
}
