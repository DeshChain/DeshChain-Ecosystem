package keeper

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/DeshChain/DeshChain-Ecosystem/x/identity/types"
)

// Identity Backup and Recovery Message Server Implementation

// CreateIdentityBackup handles creating a complete backup of an identity
func (k msgServer) CreateIdentityBackup(goCtx context.Context, msg *types.MsgCreateIdentityBackup) (*types.MsgCreateIdentityBackupResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Decode encryption key
	encryptionKey, err := base64.StdEncoding.DecodeString(msg.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("invalid encryption key: %w", err)
	}

	// Convert retention days to duration
	retentionPeriod := time.Duration(msg.RetentionDays) * 24 * time.Hour

	// Create the backup
	backup, err := k.CreateIdentityBackup(
		ctx,
		msg.HolderDID,
		msg.RecoveryMethods,
		encryptionKey,
		retentionPeriod,
	)
	if err != nil {
		return nil, err
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"identity_backup_created",
			sdk.NewAttribute("backup_id", backup.BackupID),
			sdk.NewAttribute("holder_did", backup.HolderDID),
			sdk.NewAttribute("backup_version", fmt.Sprintf("%d", backup.BackupVersion)),
			sdk.NewAttribute("recovery_methods_count", fmt.Sprintf("%d", len(backup.RecoveryMethods))),
			sdk.NewAttribute("expires_at", backup.ExpiresAt.Format(time.RFC3339)),
		),
	)

	return &types.MsgCreateIdentityBackupResponse{
		BackupID:      backup.BackupID,
		BackupVersion: backup.BackupVersion,
		CreatedAt:     backup.CreatedAt,
		ExpiresAt:     backup.ExpiresAt,
	}, nil
}

// InitiateRecovery handles initiating an identity recovery process
func (k msgServer) InitiateRecovery(goCtx context.Context, msg *types.MsgInitiateRecovery) (*types.MsgInitiateRecoveryResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the requester address
	requesterAddress, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return nil, err
	}

	// Initiate the recovery
	recoveryRequest, err := k.InitiateRecovery(
		ctx,
		requesterAddress,
		msg.HolderDID,
		msg.BackupID,
		msg.Reason,
	)
	if err != nil {
		return nil, err
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"identity_recovery_initiated",
			sdk.NewAttribute("request_id", recoveryRequest.RequestID),
			sdk.NewAttribute("holder_did", recoveryRequest.HolderDID),
			sdk.NewAttribute("backup_id", recoveryRequest.BackupID),
			sdk.NewAttribute("requested_by", recoveryRequest.RequestedBy),
			sdk.NewAttribute("reason", recoveryRequest.Reason),
			sdk.NewAttribute("expires_at", recoveryRequest.ExpiresAt.Format(time.RFC3339)),
		),
	)

	return &types.MsgInitiateRecoveryResponse{
		RequestID:     recoveryRequest.RequestID,
		RequiredScore: recoveryRequest.RequiredScore,
		ExpiresAt:     recoveryRequest.ExpiresAt,
		MaxAttempts:   recoveryRequest.MaxAttempts,
	}, nil
}

// SubmitRecoveryProof handles submitting proof for a recovery method
func (k msgServer) SubmitRecoveryProof(goCtx context.Context, msg *types.MsgSubmitRecoveryProof) (*types.MsgSubmitRecoveryProofResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the requester address
	requesterAddress, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return nil, err
	}

	// Decode proof data
	proofData, err := base64.StdEncoding.DecodeString(msg.ProofData)
	if err != nil {
		return nil, fmt.Errorf("invalid proof data: %w", err)
	}

	// Submit the recovery proof
	attempt, err := k.SubmitRecoveryProof(
		ctx,
		requesterAddress,
		msg.RequestID,
		msg.MethodID,
		proofData,
		msg.VerificationData,
	)
	if err != nil {
		return nil, err
	}

	// Get the updated recovery request
	request, found := k.GetRecoveryRequest(ctx, msg.RequestID)
	if !found {
		return nil, types.ErrRecoveryRequestNotFound
	}

	// Emit event
	eventType := "identity_recovery_proof_submitted"
	if attempt.Status == types.AttemptStatus_VERIFIED {
		eventType = "identity_recovery_proof_verified"
	} else if attempt.Status == types.AttemptStatus_FAILED {
		eventType = "identity_recovery_proof_failed"
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			eventType,
			sdk.NewAttribute("request_id", msg.RequestID),
			sdk.NewAttribute("attempt_id", attempt.AttemptID),
			sdk.NewAttribute("method_id", msg.MethodID),
			sdk.NewAttribute("method_type", attempt.MethodType.String()),
			sdk.NewAttribute("status", attempt.Status.String()),
			sdk.NewAttribute("confidence", fmt.Sprintf("%d", attempt.Confidence)),
			sdk.NewAttribute("current_score", fmt.Sprintf("%d", request.ConfidenceScore)),
			sdk.NewAttribute("required_score", fmt.Sprintf("%d", request.RequiredScore)),
		),
	)

	return &types.MsgSubmitRecoveryProofResponse{
		AttemptID:      attempt.AttemptID,
		Status:         attempt.Status,
		Confidence:     attempt.Confidence,
		CurrentScore:   request.ConfidenceScore,
		RequiredScore:  request.RequiredScore,
		CanAttemptMore: request.CanAttemptRecovery(),
	}, nil
}

// ExecuteRecovery handles executing the recovery process
func (k msgServer) ExecuteRecovery(goCtx context.Context, msg *types.MsgExecuteRecovery) (*types.MsgExecuteRecoveryResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the requester address
	requesterAddress, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return nil, err
	}

	// Get the new controller address
	newControllerAddress, err := sdk.AccAddressFromBech32(msg.NewControllerAddress)
	if err != nil {
		return nil, err
	}

	// Decode decryption key
	decryptionKey, err := base64.StdEncoding.DecodeString(msg.DecryptionKey)
	if err != nil {
		return nil, fmt.Errorf("invalid decryption key: %w", err)
	}

	// Execute the recovery
	err = k.ExecuteRecovery(
		ctx,
		requesterAddress,
		msg.RequestID,
		newControllerAddress,
		decryptionKey,
	)
	if err != nil {
		return nil, err
	}

	recoveredAt := ctx.BlockTime()

	// Get the recovery request for event details
	request, found := k.GetRecoveryRequest(ctx, msg.RequestID)
	if !found {
		return nil, types.ErrRecoveryRequestNotFound
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"identity_recovery_executed",
			sdk.NewAttribute("request_id", msg.RequestID),
			sdk.NewAttribute("holder_did", request.HolderDID),
			sdk.NewAttribute("old_controller", requesterAddress.String()),
			sdk.NewAttribute("new_controller", msg.NewControllerAddress),
			sdk.NewAttribute("recovered_at", recoveredAt.Format(time.RFC3339)),
		),
	)

	return &types.MsgExecuteRecoveryResponse{
		RecoveredAt:          recoveredAt,
		NewControllerAddress: msg.NewControllerAddress,
		DataRestored:         true, // Placeholder - should reflect actual restoration status
	}, nil
}

// AddSocialRecoveryGuardian handles adding a guardian for social recovery
func (k msgServer) AddSocialRecoveryGuardian(goCtx context.Context, msg *types.MsgAddSocialRecoveryGuardian) (*types.MsgAddSocialRecoveryGuardianResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Add the guardian
	guardian, err := k.AddSocialRecoveryGuardian(
		ctx,
		msg.HolderDID,
		msg.GuardianDID,
		msg.GuardianAddress,
		msg.GuardianName,
		msg.Weight,
		msg.ContactInfo,
		msg.PublicKey,
	)
	if err != nil {
		return nil, err
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"social_recovery_guardian_added",
			sdk.NewAttribute("guardian_id", guardian.GuardianID),
			sdk.NewAttribute("holder_did", msg.HolderDID),
			sdk.NewAttribute("guardian_did", msg.GuardianDID),
			sdk.NewAttribute("guardian_address", msg.GuardianAddress),
			sdk.NewAttribute("guardian_name", msg.GuardianName),
			sdk.NewAttribute("weight", fmt.Sprintf("%d", msg.Weight)),
			sdk.NewAttribute("added_at", guardian.AddedAt.Format(time.RFC3339)),
		),
	)

	return &types.MsgAddSocialRecoveryGuardianResponse{
		GuardianID: guardian.GuardianID,
		AddedAt:    guardian.AddedAt,
	}, nil
}

// SubmitGuardianVote handles submitting a guardian vote for recovery
func (k msgServer) SubmitGuardianVote(goCtx context.Context, msg *types.MsgSubmitGuardianVote) (*types.MsgSubmitGuardianVoteResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the guardian address
	guardianAddress, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return nil, err
	}

	// Submit the vote
	vote, err := k.SubmitGuardianVote(
		ctx,
		guardianAddress,
		msg.RequestID,
		msg.Vote,
		msg.Reason,
		msg.Signature,
	)
	if err != nil {
		return nil, err
	}

	// Get the recovery request to check voting status
	request, found := k.GetRecoveryRequest(ctx, msg.RequestID)
	if !found {
		return nil, types.ErrRecoveryRequestNotFound
	}

	// Calculate current and required votes (placeholder logic)
	currentVotes := len(request.RecoveryMethods) // Simplified
	requiredVotes := 3                          // Placeholder

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"guardian_vote_submitted",
			sdk.NewAttribute("vote_id", vote.VoteID),
			sdk.NewAttribute("request_id", msg.RequestID),
			sdk.NewAttribute("guardian_id", vote.GuardianID),
			sdk.NewAttribute("vote", vote.Vote.String()),
			sdk.NewAttribute("weight", fmt.Sprintf("%d", vote.Weight)),
			sdk.NewAttribute("current_votes", fmt.Sprintf("%d", currentVotes)),
			sdk.NewAttribute("required_votes", fmt.Sprintf("%d", requiredVotes)),
			sdk.NewAttribute("voted_at", vote.VotedAt.Format(time.RFC3339)),
		),
	)

	return &types.MsgSubmitGuardianVoteResponse{
		VoteID:        vote.VoteID,
		VotedAt:       vote.VotedAt,
		Weight:        vote.Weight,
		CurrentVotes:  currentVotes,
		RequiredVotes: requiredVotes,
	}, nil
}

// VerifyBackupIntegrity handles verifying the integrity of a backup
func (k msgServer) VerifyBackupIntegrity(goCtx context.Context, msg *types.MsgVerifyBackupIntegrity) (*types.MsgVerifyBackupIntegrityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Decode verification key if provided
	var verificationKey []byte
	var err error
	if msg.VerificationKey != "" {
		verificationKey, err = base64.StdEncoding.DecodeString(msg.VerificationKey)
		if err != nil {
			return nil, fmt.Errorf("invalid verification key: %w", err)
		}
	}

	// Verify the backup integrity
	result, err := k.VerifyBackupIntegrity(
		ctx,
		msg.BackupID,
		verificationKey,
	)
	if err != nil {
		return nil, err
	}

	// Get the backup for event details
	backup, found := k.GetIdentityBackup(ctx, msg.BackupID)
	if !found {
		return nil, types.ErrBackupNotFound
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"backup_integrity_verified",
			sdk.NewAttribute("backup_id", msg.BackupID),
			sdk.NewAttribute("verification_id", result.VerificationID),
			sdk.NewAttribute("holder_did", backup.HolderDID),
			sdk.NewAttribute("integrity_valid", fmt.Sprintf("%t", result.IntegrityValid)),
			sdk.NewAttribute("decryption_valid", fmt.Sprintf("%t", result.DecryptionValid)),
			sdk.NewAttribute("data_completeness", fmt.Sprintf("%.2f", result.DataCompleteness)),
			sdk.NewAttribute("recoverability_score", fmt.Sprintf("%d", result.RecoverabilityScore)),
			sdk.NewAttribute("issues_count", fmt.Sprintf("%d", len(result.IssuesFound))),
			sdk.NewAttribute("verified_at", result.VerifiedAt.Format(time.RFC3339)),
		),
	)

	return &types.MsgVerifyBackupIntegrityResponse{
		VerificationID:      result.VerificationID,
		IntegrityValid:      result.IntegrityValid,
		DecryptionValid:     result.DecryptionValid,
		DataCompleteness:    result.DataCompleteness,
		RecoverabilityScore: result.RecoverabilityScore,
		IssuesFound:         result.IssuesFound,
		Recommendations:     result.Recommendations,
	}, nil
}