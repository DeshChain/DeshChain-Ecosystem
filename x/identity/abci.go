package identity

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	
	"github.com/deshchain/x/identity/keeper"
	"github.com/deshchain/x/identity/types"
)

// BeginBlocker handles block beginning logic for identity module
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	// Process expired credentials
	processExpiredCredentials(ctx, k)
	
	// Process expired consents
	processExpiredConsents(ctx, k)
	
	// Process expired proofs
	processExpiredProofs(ctx, k)
	
	// Update identity activity tracking
	updateIdentityActivity(ctx, k)
}

// EndBlocker handles block ending logic for identity module
func EndBlocker(ctx sdk.Context, req abci.RequestEndBlock, k keeper.Keeper) []abci.ValidatorUpdate {
	// Process pending recovery requests
	processPendingRecoveries(ctx, k)
	
	// Update analytics and metrics
	updateAnalytics(ctx, k)
	
	return []abci.ValidatorUpdate{}
}

// processExpiredCredentials revokes credentials that have expired
func processExpiredCredentials(ctx sdk.Context, k keeper.Keeper) {
	currentTime := ctx.BlockTime()
	
	k.IterateCredentials(ctx, func(credential types.VerifiableCredential) bool {
		if credential.IsExpired() {
			// Add to revocation list
			k.RevokeCredential(ctx, credential.ID, "Credential expired")
			
			// Emit event
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeCredentialRevoked,
					sdk.NewAttribute(types.AttributeKeyCredentialID, credential.ID),
					sdk.NewAttribute(types.AttributeKeyIssuer, credential.Issuer),
					sdk.NewAttribute("reason", "expired"),
				),
			)
		}
		return false
	})
}

// processExpiredConsents withdraws consents that have expired
func processExpiredConsents(ctx sdk.Context, k keeper.Keeper) {
	currentTime := ctx.BlockTime()
	
	k.IterateIdentities(ctx, func(identity types.Identity) bool {
		for i, consent := range identity.Consents {
			if consent.Given && consent.ExpiresAt != nil && consent.ExpiresAt.Before(currentTime) {
				// Mark consent as withdrawn
				identity.Consents[i].Given = false
				identity.Consents[i].WithdrawnAt = &currentTime
				
				// Update identity
				k.SetIdentity(ctx, identity)
				
				// Emit event
				ctx.EventManager().EmitEvent(
					sdk.NewEvent(
						types.EventTypeConsentWithdrawn,
						sdk.NewAttribute(types.AttributeKeyAddress, identity.Address),
						sdk.NewAttribute(types.AttributeKeyConsentType, consent.Type.String()),
						sdk.NewAttribute("reason", "expired"),
					),
				)
			}
		}
		return false
	})
}

// processExpiredProofs removes expired zero-knowledge proofs
func processExpiredProofs(ctx sdk.Context, k keeper.Keeper) {
	currentTime := ctx.BlockTime()
	
	k.IterateZKProofs(ctx, func(proof types.ZKProof) bool {
		if proof.ExpiresAt != nil && proof.ExpiresAt.Before(currentTime) {
			// Remove expired proof
			k.DeleteZKProof(ctx, proof.ID)
			
			// Emit event
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					"zk_proof_expired",
					sdk.NewAttribute(types.AttributeKeyProofType, proof.Type.String()),
					sdk.NewAttribute("proof_id", proof.ID),
				),
			)
		}
		return false
	})
}

// updateIdentityActivity updates last activity timestamps
func updateIdentityActivity(ctx sdk.Context, k keeper.Keeper) {
	// This could be optimized to only update identities that had activity
	// For now, we'll skip automatic updates to avoid unnecessary writes
}

// processPendingRecoveries processes account recovery requests
func processPendingRecoveries(ctx sdk.Context, k keeper.Keeper) {
	k.IterateRecoveryRequests(ctx, func(request types.RecoveryRequest) bool {
		// Check if recovery period has elapsed
		recoveryPeriod := 48 * time.Hour // 48 hours recovery period
		
		if ctx.BlockTime().Sub(request.InitiatedAt) > recoveryPeriod {
			if request.Status == "pending" {
				// Recovery period elapsed, allow completion
				request.Status = "ready"
				k.SetRecoveryRequest(ctx, request)
				
				// Emit event
				ctx.EventManager().EmitEvent(
					sdk.NewEvent(
						"recovery_ready",
						sdk.NewAttribute("recovery_id", request.ID),
						sdk.NewAttribute("target_address", request.TargetAddress),
					),
				)
			}
		}
		return false
	})
}

// updateAnalytics updates module analytics and usage metrics
func updateAnalytics(ctx sdk.Context, k keeper.Keeper) {
	// Update daily active identities
	k.UpdateDailyActiveIdentities(ctx)
	
	// Update credential issuance metrics
	k.UpdateCredentialMetrics(ctx)
	
	// Update verification metrics
	k.UpdateVerificationMetrics(ctx)
}