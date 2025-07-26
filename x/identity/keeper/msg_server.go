package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/DeshChain/DeshChain-Ecosystem/x/identity/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the identity MsgServer interface
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// CreateIdentity creates a new identity
func (k msgServer) CreateIdentity(goCtx context.Context, msg *types.MsgCreateIdentity) (*types.MsgCreateIdentityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Validate identity creation
	if err := k.ValidateIdentityCreation(ctx, msg.Creator); err != nil {
		return nil, err
	}
	
	// Generate DID
	did := types.GenerateDID(msg.Creator)
	
	// Create DID document
	didDoc := types.DIDDocument{
		Context: []string{
			"https://www.w3.org/ns/did/v1",
			types.ContextDeshChain,
		},
		ID:         did,
		Controller: did,
		Created:    ctx.BlockTime(),
		Updated:    ctx.BlockTime(),
		Service:    msg.ServiceEndpoints,
	}
	
	// Add verification method from public key
	if msg.PublicKey != "" {
		vm := types.VerificationMethod{
			ID:                 fmt.Sprintf("%s#key-1", did),
			Type:               "Ed25519VerificationKey2020",
			Controller:         did,
			PublicKeyMultibase: msg.PublicKey,
			Created:            ctx.BlockTime(),
		}
		didDoc.VerificationMethod = []types.VerificationMethod{vm}
		didDoc.Authentication = []interface{}{vm.ID}
	}
	
	// Validate and store DID document
	if err := k.ValidateDIDRegistration(ctx, did, msg.Creator); err != nil {
		return nil, err
	}
	k.SetDIDDocument(ctx, didDoc)
	
	// Create identity
	identity := types.Identity{
		Address:         msg.Creator,
		DID:             did,
		DIDDocument:     &didDoc,
		Credentials:     []string{},
		Status:          types.IdentityStatus_ACTIVE,
		RecoveryMethods: msg.RecoveryMethods,
		Consents:        msg.InitialConsents,
		CreatedAt:       ctx.BlockTime(),
		UpdatedAt:       ctx.BlockTime(),
		LastActivityAt:  ctx.BlockTime(),
		Metadata:        msg.Metadata,
	}
	
	// Store identity
	k.SetIdentity(ctx, identity)
	k.SetIdentityDIDIndex(ctx, did, msg.Creator)
	
	// Log operation
	k.LogIdentityOperation(ctx, "create_identity", msg.Creator, map[string]string{
		"did": did,
	})
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeIdentityCreated,
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Creator),
			sdk.NewAttribute(types.AttributeKeyDID, did),
		),
	)
	
	return &types.MsgCreateIdentityResponse{
		Did: did,
	}, nil
}

// UpdateIdentity updates an existing identity
func (k msgServer) UpdateIdentity(goCtx context.Context, msg *types.MsgUpdateIdentity) (*types.MsgUpdateIdentityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Get identity
	identity, found := k.GetIdentity(ctx, msg.Creator)
	if !found {
		return nil, types.ErrIdentityNotFound
	}
	
	// Check if identity is active
	if !identity.IsActive() {
		return nil, types.ErrIdentityInactive
	}
	
	// Update fields
	if len(msg.ServiceEndpoints) > 0 {
		// Update DID document services
		didDoc, found := k.GetDIDDocument(ctx, identity.DID)
		if found {
			didDoc.Service = msg.ServiceEndpoints
			didDoc.Updated = ctx.BlockTime()
			k.SetDIDDocument(ctx, didDoc)
		}
	}
	
	if len(msg.RecoveryMethods) > 0 {
		identity.RecoveryMethods = msg.RecoveryMethods
	}
	
	if msg.Metadata != nil {
		identity.Metadata = msg.Metadata
	}
	
	identity.UpdatedAt = ctx.BlockTime()
	identity.LastActivityAt = ctx.BlockTime()
	
	// Store updated identity
	k.SetIdentity(ctx, identity)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeIdentityUpdated,
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Creator),
			sdk.NewAttribute(types.AttributeKeyDID, identity.DID),
		),
	)
	
	return &types.MsgUpdateIdentityResponse{}, nil
}

// RevokeIdentity revokes an identity
func (k msgServer) RevokeIdentity(goCtx context.Context, msg *types.MsgRevokeIdentity) (*types.MsgRevokeIdentityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Get identity
	identity, found := k.GetIdentity(ctx, msg.Creator)
	if !found {
		return nil, types.ErrIdentityNotFound
	}
	
	// Update status
	if err := k.UpdateIdentityStatus(ctx, msg.Creator, types.IdentityStatus_REVOKED); err != nil {
		return nil, err
	}
	
	// Deactivate DID
	if err := k.DeactivateDID(ctx, identity.DID); err != nil {
		return nil, err
	}
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeIdentityRevoked,
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Creator),
			sdk.NewAttribute(types.AttributeKeyDID, identity.DID),
			sdk.NewAttribute("reason", msg.Reason),
		),
	)
	
	return &types.MsgRevokeIdentityResponse{}, nil
}

// RegisterDID registers a new DID
func (k msgServer) RegisterDID(goCtx context.Context, msg *types.MsgRegisterDID) (*types.MsgRegisterDIDResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Validate DID document
	if err := types.ValidateDIDDocument(msg.DIDDocument); err != nil {
		return nil, err
	}
	
	// Validate DID registration
	if err := k.ValidateDIDRegistration(ctx, msg.DIDDocument.ID, msg.Creator); err != nil {
		return nil, err
	}
	
	// Validate document size
	if err := k.ValidateDIDDocumentSize(ctx, msg.DIDDocument); err != nil {
		return nil, err
	}
	
	// Store DID document
	msg.DIDDocument.Created = ctx.BlockTime()
	msg.DIDDocument.Updated = ctx.BlockTime()
	k.SetDIDDocument(ctx, *msg.DIDDocument)
	
	// Create or update identity if needed
	if !k.HasIdentity(ctx, msg.Creator) {
		identity := types.Identity{
			Address:        msg.Creator,
			DID:            msg.DIDDocument.ID,
			DIDDocument:    msg.DIDDocument,
			Status:         types.IdentityStatus_ACTIVE,
			CreatedAt:      ctx.BlockTime(),
			UpdatedAt:      ctx.BlockTime(),
			LastActivityAt: ctx.BlockTime(),
		}
		k.SetIdentity(ctx, identity)
		k.SetIdentityDIDIndex(ctx, msg.DIDDocument.ID, msg.Creator)
	}
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDIDRegistered,
			sdk.NewAttribute(types.AttributeKeyDID, msg.DIDDocument.ID),
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Creator),
		),
	)
	
	return &types.MsgRegisterDIDResponse{}, nil
}

// UpdateDID updates a DID document
func (k msgServer) UpdateDID(goCtx context.Context, msg *types.MsgUpdateDID) (*types.MsgUpdateDIDResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Get DID document
	didDoc, found := k.GetDIDDocument(ctx, msg.DID)
	if !found {
		return nil, types.ErrDIDNotFound
	}
	
	// Verify authorization
	controller, err := k.GetDIDController(ctx, msg.DID)
	if err != nil {
		return nil, err
	}
	
	// Check if creator controls the DID
	identity, found := k.GetIdentity(ctx, msg.Creator)
	if !found || identity.DID != controller {
		return nil, types.ErrUnauthorized
	}
	
	// Add verification methods
	for _, vm := range msg.VerificationMethods {
		vm.Created = ctx.BlockTime()
		if err := k.AddVerificationMethod(ctx, msg.DID, vm); err != nil {
			return nil, err
		}
	}
	
	// Remove verification methods
	for _, vmID := range msg.RemoveVerificationMethods {
		if err := k.RemoveVerificationMethod(ctx, msg.DID, vmID); err != nil {
			return nil, err
		}
	}
	
	// Add services
	for _, svc := range msg.Services {
		if err := k.AddService(ctx, msg.DID, svc); err != nil {
			return nil, err
		}
	}
	
	// Remove services
	for _, svcID := range msg.RemoveServices {
		if err := k.RemoveService(ctx, msg.DID, svcID); err != nil {
			return nil, err
		}
	}
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDIDUpdated,
			sdk.NewAttribute(types.AttributeKeyDID, msg.DID),
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Creator),
		),
	)
	
	return &types.MsgUpdateDIDResponse{}, nil
}

// DeactivateDID deactivates a DID
func (k msgServer) DeactivateDID(goCtx context.Context, msg *types.MsgDeactivateDID) (*types.MsgDeactivateDIDResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Verify authorization
	controller, err := k.GetDIDController(ctx, msg.DID)
	if err != nil {
		return nil, err
	}
	
	identity, found := k.GetIdentity(ctx, msg.Creator)
	if !found || identity.DID != controller {
		return nil, types.ErrUnauthorized
	}
	
	// Deactivate DID
	if err := k.Keeper.DeactivateDID(ctx, msg.DID); err != nil {
		return nil, err
	}
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"did_deactivated",
			sdk.NewAttribute(types.AttributeKeyDID, msg.DID),
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Creator),
			sdk.NewAttribute("reason", msg.Reason),
		),
	)
	
	return &types.MsgDeactivateDIDResponse{}, nil
}

// IssueCredential issues a new verifiable credential
func (k msgServer) IssueCredential(goCtx context.Context, msg *types.MsgIssueCredential) (*types.MsgIssueCredentialResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Check if consent is required
	if msg.RequireConsent {
		hasConsent, err := k.CheckConsent(ctx, msg.Holder, types.ConsentType_DATA_COLLECTION, msg.Issuer)
		if err != nil || !hasConsent {
			return nil, types.ErrConsentNotGiven
		}
	}
	
	// Issue credential
	credential, err := k.Keeper.IssueCredential(
		ctx,
		msg.Issuer,
		msg.Holder,
		msg.CredentialType,
		msg.Claims,
		msg.ExpirationDays,
	)
	if err != nil {
		return nil, err
	}
	
	// Add evidence if provided
	if len(msg.Evidence) > 0 {
		credential.Evidence = msg.Evidence
	}
	
	// Add metadata if provided
	if msg.Metadata != nil {
		credential.Metadata = make(map[string]interface{})
		for k, v := range msg.Metadata {
			credential.Metadata[k] = v
		}
	}
	
	// Update credential
	k.SetCredential(ctx, *credential)
	
	// Log operation
	k.LogIdentityOperation(ctx, "issue_credential", msg.Holder, map[string]string{
		"credential_id":   credential.ID,
		"credential_type": msg.CredentialType,
		"issuer":          msg.Issuer,
	})
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCredentialIssued,
			sdk.NewAttribute(types.AttributeKeyCredentialID, credential.ID),
			sdk.NewAttribute(types.AttributeKeyCredentialType, msg.CredentialType),
			sdk.NewAttribute(types.AttributeKeyIssuer, msg.Issuer),
			sdk.NewAttribute(types.AttributeKeyHolder, msg.Holder),
		),
	)
	
	return &types.MsgIssueCredentialResponse{
		CredentialId: credential.ID,
	}, nil
}

// RevokeCredential revokes a credential
func (k msgServer) RevokeCredential(goCtx context.Context, msg *types.MsgRevokeCredential) (*types.MsgRevokeCredentialResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Get credential
	credential, found := k.GetCredential(ctx, msg.CredentialID)
	if !found {
		return nil, types.ErrCredentialNotFound
	}
	
	// Verify issuer
	if credential.Issuer != msg.Issuer {
		return nil, types.ErrUnauthorized
	}
	
	// Revoke credential
	if err := k.Keeper.RevokeCredential(ctx, msg.CredentialID, msg.Reason); err != nil {
		return nil, err
	}
	
	// Log operation
	k.LogIdentityOperation(ctx, "revoke_credential", msg.Issuer, map[string]string{
		"credential_id": msg.CredentialID,
		"reason":        msg.Reason,
	})
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCredentialRevoked,
			sdk.NewAttribute(types.AttributeKeyCredentialID, msg.CredentialID),
			sdk.NewAttribute(types.AttributeKeyIssuer, msg.Issuer),
			sdk.NewAttribute("reason", msg.Reason),
		),
	)
	
	return &types.MsgRevokeCredentialResponse{}, nil
}

// PresentCredential creates a verifiable presentation
func (k msgServer) PresentCredential(goCtx context.Context, msg *types.MsgPresentCredential) (*types.MsgPresentCredentialResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Create presentation
	presentation, err := k.Keeper.PresentCredential(
		ctx,
		msg.Holder,
		msg.CredentialIDs,
		msg.Verifier,
		msg.Challenge,
		msg.Domain,
	)
	if err != nil {
		return nil, err
	}
	
	// Apply selective disclosure if requested
	if len(msg.RevealedClaims) > 0 {
		// Filter credentials to only include revealed claims
		// This is a simplified implementation
		for credID, revealedClaims := range msg.RevealedClaims {
			disclosure, err := k.CreateSelectiveDisclosureProof(ctx, credID, revealedClaims, msg.Holder)
			if err != nil {
				return nil, err
			}
			
			if presentation.Metadata == nil {
				presentation.Metadata = make(map[string]interface{})
			}
			presentation.Metadata[fmt.Sprintf("disclosure_%s", credID)] = disclosure
		}
	}
	
	// Log operation
	k.LogIdentityOperation(ctx, "present_credential", msg.Holder, map[string]string{
		"presentation_id": presentation.ID,
		"verifier":        msg.Verifier,
		"credential_count": fmt.Sprintf("%d", len(msg.CredentialIDs)),
	})
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCredentialPresented,
			sdk.NewAttribute("presentation_id", presentation.ID),
			sdk.NewAttribute(types.AttributeKeyHolder, msg.Holder),
			sdk.NewAttribute("verifier", msg.Verifier),
		),
	)
	
	return &types.MsgPresentCredentialResponse{
		PresentationId: presentation.ID,
	}, nil
}

// CreateZKProof creates a zero-knowledge proof
func (k msgServer) CreateZKProof(goCtx context.Context, msg *types.MsgCreateZKProof) (*types.MsgCreateZKProofResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Validate credentials exist
	for _, credID := range msg.CredentialIDs {
		if !k.HasCredential(ctx, credID) {
			return nil, types.ErrCredentialNotFound
		}
		
		// Verify creator owns the credential
		credential, _ := k.GetCredential(ctx, credID)
		subjectID, err := credential.GetSubjectID()
		if err != nil || subjectID != msg.Creator {
			return nil, types.ErrUnauthorized
		}
	}
	
	// Create ZK proof
	proof, err := k.Keeper.CreateZKProof(
		ctx,
		msg.Creator,
		msg.ProofType,
		msg.Statement,
		msg.ProofData,
		msg.ExpiryMinutes,
	)
	if err != nil {
		return nil, err
	}
	
	// Add public inputs
	if len(msg.PublicInputs) > 0 {
		proof.PublicInputs = msg.PublicInputs
	}
	
	// Update proof
	k.SetZKProof(ctx, *proof)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeZKProofCreated,
			sdk.NewAttribute(types.AttributeKeyProofType, msg.ProofType.String()),
			sdk.NewAttribute("proof_id", proof.ID),
			sdk.NewAttribute("creator", msg.Creator),
		),
	)
	
	return &types.MsgCreateZKProofResponse{
		ProofId: proof.ID,
	}, nil
}

// VerifyZKProof verifies a zero-knowledge proof
func (k msgServer) VerifyZKProof(goCtx context.Context, msg *types.MsgVerifyZKProof) (*types.MsgVerifyZKProofResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Verify proof
	if err := k.Keeper.VerifyZKProof(ctx, msg.ProofID, msg.Verifier); err != nil {
		return &types.MsgVerifyZKProofResponse{
			Valid: false,
			Error: err.Error(),
		}, nil
	}
	
	return &types.MsgVerifyZKProofResponse{
		Valid: true,
	}, nil
}

// LinkAadhaar links Aadhaar to identity
func (k msgServer) LinkAadhaar(goCtx context.Context, msg *types.MsgLinkAadhaar) (*types.MsgLinkAadhaarResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Link Aadhaar
	if err := k.Keeper.LinkAadhaar(
		ctx,
		msg.Creator,
		msg.AadhaarHash,
		msg.DemographicHash,
		msg.BiometricHash,
		msg.VerificationMethod,
		msg.ConsentArtefact,
	); err != nil {
		return nil, err
	}
	
	return &types.MsgLinkAadhaarResponse{}, nil
}

// ConnectDigiLocker connects DigiLocker account
func (k msgServer) ConnectDigiLocker(goCtx context.Context, msg *types.MsgConnectDigiLocker) (*types.MsgConnectDigiLockerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Connect DigiLocker
	if err := k.Keeper.ConnectDigiLocker(
		ctx,
		msg.Creator,
		msg.AuthToken,
		msg.ConsentID,
		msg.DocumentTypes,
	); err != nil {
		return nil, err
	}
	
	return &types.MsgConnectDigiLockerResponse{
		DocumentCount: int32(len(msg.DocumentTypes)),
	}, nil
}

// LinkUPI links UPI ID to identity
func (k msgServer) LinkUPI(goCtx context.Context, msg *types.MsgLinkUPI) (*types.MsgLinkUPIResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Link UPI
	if err := k.Keeper.LinkUPI(
		ctx,
		msg.Creator,
		msg.VPAURI,
		msg.PSPProvider,
		msg.AuthToken,
	); err != nil {
		return nil, err
	}
	
	return &types.MsgLinkUPIResponse{}, nil
}

// GiveConsent gives consent for data usage
func (k msgServer) GiveConsent(goCtx context.Context, msg *types.MsgGiveConsent) (*types.MsgGiveConsentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Give consent
	consent, err := k.Keeper.GiveConsent(
		ctx,
		msg.Creator,
		msg.ConsentType,
		msg.Purpose,
		msg.DataController,
		msg.DataCategories,
		msg.ProcessingTypes,
		msg.ExpirationDays,
	)
	if err != nil {
		return nil, err
	}
	
	return &types.MsgGiveConsentResponse{
		ConsentId: consent.ID,
	}, nil
}

// WithdrawConsent withdraws consent
func (k msgServer) WithdrawConsent(goCtx context.Context, msg *types.MsgWithdrawConsent) (*types.MsgWithdrawConsentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Withdraw consent
	if err := k.Keeper.WithdrawConsent(ctx, msg.Creator, msg.ConsentID, msg.Reason); err != nil {
		return nil, err
	}
	
	return &types.MsgWithdrawConsentResponse{}, nil
}

// AddRecoveryMethod adds a recovery method
func (k msgServer) AddRecoveryMethod(goCtx context.Context, msg *types.MsgAddRecoveryMethod) (*types.MsgAddRecoveryMethodResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Create recovery method
	method := types.RecoveryMethod{
		Type:     msg.Type,
		Value:    msg.Value, // Should be encrypted/hashed in production
		AddedAt:  ctx.BlockTime(),
		IsActive: true,
	}
	
	// Add recovery method
	if err := k.Keeper.AddRecoveryMethod(ctx, msg.Creator, method); err != nil {
		return nil, err
	}
	
	return &types.MsgAddRecoveryMethodResponse{}, nil
}

// InitiateRecovery initiates account recovery
func (k msgServer) InitiateRecovery(goCtx context.Context, msg *types.MsgInitiateRecovery) (*types.MsgInitiateRecoveryResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Initiate recovery
	request, err := k.Keeper.InitiateRecovery(
		ctx,
		msg.Creator,
		msg.TargetAddress,
		msg.RecoveryType,
		msg.RecoveryData,
	)
	if err != nil {
		return nil, err
	}
	
	return &types.MsgInitiateRecoveryResponse{
		RecoveryId: request.ID,
	}, nil
}

// CompleteRecovery completes account recovery
func (k msgServer) CompleteRecovery(goCtx context.Context, msg *types.MsgCompleteRecovery) (*types.MsgCompleteRecoveryResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Complete recovery
	if err := k.Keeper.CompleteRecovery(
		ctx,
		msg.RecoveryID,
		msg.NewPublicKey,
		msg.ProofData,
	); err != nil {
		return nil, err
	}
	
	return &types.MsgCompleteRecoveryResponse{}, nil
}

// UpdatePrivacySettings updates privacy settings
func (k msgServer) UpdatePrivacySettings(goCtx context.Context, msg *types.MsgUpdatePrivacySettings) (*types.MsgUpdatePrivacySettingsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Create privacy settings
	settings := types.PrivacySettings{
		UserAddress:              msg.Creator,
		DefaultDisclosureLevel:   msg.DefaultDisclosureLevel,
		AllowAnonymousUsage:      msg.AllowAnonymousUsage,
		RequireExplicitConsent:   msg.RequireExplicitConsent,
		DataMinimization:         msg.DataMinimization,
		AutoDeleteAfterDays:      msg.AutoDeleteAfterDays,
		AllowDerivedCredentials:  msg.AllowDerivedCredentials,
		PreferredProofSystems:    msg.PreferredProofSystems,
		BlacklistedVerifiers:     msg.BlacklistedVerifiers,
		UpdatedAt:                ctx.BlockTime(),
	}
	
	// Store settings
	k.SetPrivacySettings(ctx, settings)
	
	return &types.MsgUpdatePrivacySettingsResponse{}, nil
}