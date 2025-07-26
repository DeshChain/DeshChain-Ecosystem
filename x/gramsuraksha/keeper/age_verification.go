package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/gramsuraksha/types"
)

// VerifyParticipantAgeWithZKP verifies participant age using zero-knowledge proof
func (k Keeper) VerifyParticipantAgeWithZKP(
	ctx sdk.Context,
	participantAddress sdk.AccAddress,
	minAge, maxAge int32,
) (bool, error) {
	if k.identityAdapter == nil {
		// Fall back to traditional age verification
		return k.verifyAgeTraditional(ctx, participantAddress, minAge, maxAge)
	}

	// Use identity adapter for ZK proof verification
	return k.identityAdapter.VerifyAgeWithZKProof(ctx, participantAddress, minAge, maxAge)
}

// verifyAgeTraditional performs traditional age verification
func (k Keeper) verifyAgeTraditional(
	ctx sdk.Context,
	participantAddress sdk.AccAddress,
	minAge, maxAge int32,
) (bool, error) {
	// Check if we have participant data
	participants := k.GetParticipantsByAddress(ctx, participantAddress)
	if len(participants) == 0 {
		return false, fmt.Errorf("no participant data found")
	}

	// Check age from existing participant data
	for _, participant := range participants {
		if participant.Age >= minAge && participant.Age <= maxAge {
			return true, nil
		}
	}

	return false, nil
}

// GetParticipantIdentityStatus returns the identity status of a participant
func (k Keeper) GetParticipantIdentityStatus(
	ctx sdk.Context,
	participantAddress sdk.AccAddress,
) (*ParticipantIdentityInfo, error) {
	info := &ParticipantIdentityInfo{
		Address:           participantAddress.String(),
		HasTraditionalKYC: false,
		HasIdentity:       false,
	}

	// Check traditional KYC
	if k.kycKeeper != nil && k.kycKeeper.IsKYCVerified(ctx, participantAddress) {
		info.HasTraditionalKYC = true
	}

	// Check identity-based verification
	if k.identityAdapter != nil {
		identityStatus, err := k.identityAdapter.VerifyParticipantIdentity(ctx, participantAddress, "")
		if err == nil {
			info.HasIdentity = identityStatus.HasIdentity
			info.DID = identityStatus.DID
			info.IsKYCVerified = identityStatus.IsKYCVerified
			info.KYCLevel = identityStatus.KYCLevel
			
			// Get all GramSuraksha credentials
			creds, _ := k.identityAdapter.GetParticipantCredentials(ctx, participantAddress)
			info.GramSurakshaCredentials = make([]GramSurakshaCredentialInfo, len(creds))
			for i, cred := range creds {
				info.GramSurakshaCredentials[i] = GramSurakshaCredentialInfo{
					CredentialID: cred.CredentialID,
					SchemeID:     cred.SchemeID,
					Status:       cred.Status,
					IsValid:      cred.IsValid,
				}
			}
		}
	}

	// Get enrolled schemes
	participants := k.GetParticipantsByAddress(ctx, participantAddress)
	info.EnrolledSchemes = make([]string, len(participants))
	for i, p := range participants {
		info.EnrolledSchemes[i] = p.SchemeID
	}

	return info, nil
}

// UpdateParticipantIdentityCredential updates the identity credential when participant status changes
func (k Keeper) UpdateParticipantIdentityCredential(
	ctx sdk.Context,
	participantID string,
	newStatus string,
) error {
	// Get participant
	participant, found := k.GetParticipant(ctx, participantID)
	if !found {
		return types.ErrParticipantNotFound
	}

	// Update traditional status
	participant.Status = newStatus
	k.SetParticipant(ctx, participant)

	// Update identity credential if available
	if k.identityAdapter != nil {
		reason := fmt.Sprintf("Status changed to %s", newStatus)
		if err := k.identityAdapter.UpdateParticipantCredentialStatus(
			ctx, participantID, newStatus, reason,
		); err != nil {
			k.Logger(ctx).Error("Failed to update credential status", "error", err)
			// Don't fail the operation if credential update fails
		}
	}

	return nil
}

// MigrateToIdentitySystem migrates all existing participants to the identity system
func (k Keeper) MigrateToIdentitySystem(ctx sdk.Context) error {
	if k.identityAdapter == nil {
		return fmt.Errorf("identity adapter not initialized")
	}

	migrated, err := k.identityAdapter.MigrateExistingParticipants(ctx)
	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	k.Logger(ctx).Info("GramSuraksha identity migration completed", 
		"migrated_participants", migrated)

	// Emit migration event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"gramsuraksha_identity_migration",
			sdk.NewAttribute("migrated_count", fmt.Sprintf("%d", migrated)),
			sdk.NewAttribute("timestamp", ctx.BlockTime().String()),
		),
	)

	return nil
}

// Types for identity information

type ParticipantIdentityInfo struct {
	Address                 string
	HasTraditionalKYC       bool
	HasIdentity             bool
	DID                     string
	IsKYCVerified           bool
	KYCLevel                string
	EnrolledSchemes         []string
	GramSurakshaCredentials []GramSurakshaCredentialInfo
}

type GramSurakshaCredentialInfo struct {
	CredentialID string
	SchemeID     string
	Status       string
	IsValid      bool
}