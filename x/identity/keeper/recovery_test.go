package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/namo/x/identity/keeper"
	"github.com/namo/x/identity/types"
	"github.com/namo/testutil"
)

type RecoveryTestSuite struct {
	suite.Suite

	ctx    sdk.Context
	keeper *keeper.Keeper
	addrs  []sdk.AccAddress
}

func (suite *RecoveryTestSuite) SetupTest() {
	suite.ctx, suite.keeper = testutil.IdentityKeeperTestSetup(suite.T())
	suite.addrs = testutil.CreateIncrementalAccounts(10)
}

func TestRecoveryTestSuite(t *testing.T) {
	suite.Run(t, new(RecoveryTestSuite))
}

// TestCreateIdentityBackup tests creating a backup
func (suite *RecoveryTestSuite) TestCreateIdentityBackup() {
	ctx := suite.ctx
	k := suite.keeper
	holderDID := "did:desh:backup123"
	
	// Create holder identity
	identity := types.Identity{
		DID:        holderDID,
		Controller: suite.addrs[0].String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)
	
	// Define recovery methods
	recoveryMethods := []types.RecoveryMethod{
		{
			MethodID:           "mnemonic_1",
			MethodType:         types.RecoveryMethodType_MNEMONIC_PHRASE,
			MethodName:         "12-word mnemonic",
			Configuration:      map[string]interface{}{"words": 12},
			TrustLevel:         "high",
			RequiredConfidence: 90,
			Enabled:            true,
			CreatedAt:          ctx.BlockTime(),
			UsageCount:         0,
			Metadata:           map[string]interface{}{},
		},
		{
			MethodID:           "social_1",
			MethodType:         types.RecoveryMethodType_SOCIAL_RECOVERY,
			MethodName:         "Social recovery",
			Configuration:      map[string]interface{}{"threshold": 3},
			TrustLevel:         "medium",
			RequiredConfidence: 75,
			Enabled:            true,
			CreatedAt:          ctx.BlockTime(),
			UsageCount:         0,
			Metadata:           map[string]interface{}{},
		},
	}
	
	encryptionKey := []byte("test-encryption-key-32-bytes!!")
	retentionPeriod := 365 * 24 * time.Hour
	
	// Create backup
	backup, err := k.CreateIdentityBackup(
		ctx,
		holderDID,
		recoveryMethods,
		encryptionKey,
		retentionPeriod,
	)
	
	suite.NoError(err)
	suite.NotNil(backup)
	suite.NotEmpty(backup.BackupID)
	suite.Equal(holderDID, backup.HolderDID)
	suite.Equal(int64(1), backup.BackupVersion)
	suite.Len(backup.RecoveryMethods, 2)
	suite.Equal(types.BackupStatus_ACTIVE, backup.Status)
	suite.Equal("AES-256-GCM", backup.EncryptionMethod)
	suite.NotEmpty(backup.IntegrityHash)
	
	// Verify backup is stored
	storedBackup, found := k.GetIdentityBackup(ctx, backup.BackupID)
	suite.True(found)
	suite.Equal(backup.BackupID, storedBackup.BackupID)
	
	// Verify backup is valid
	suite.True(backup.IsValid())
	suite.False(backup.IsExpired())
}

// TestInitiateRecovery tests initiating recovery
func (suite *RecoveryTestSuite) TestInitiateRecovery() {
	ctx := suite.ctx
	k := suite.keeper
	holderDID := "did:desh:recovery123"
	requester := suite.addrs[0]
	
	// Create holder identity and backup first
	identity := types.Identity{
		DID:        holderDID,
		Controller: requester.String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)
	
	recoveryMethods := []types.RecoveryMethod{
		{
			MethodID:           "mnemonic_1",
			MethodType:         types.RecoveryMethodType_MNEMONIC_PHRASE,
			MethodName:         "Mnemonic phrase",
			RequiredConfidence: 90,
			Enabled:            true,
			CreatedAt:          ctx.BlockTime(),
		},
	}
	
	backup, err := k.CreateIdentityBackup(
		ctx,
		holderDID,
		recoveryMethods,
		[]byte("test-key"),
		365*24*time.Hour,
	)
	suite.NoError(err)
	
	// Initiate recovery
	recoveryRequest, err := k.InitiateRecovery(
		ctx,
		requester,
		holderDID,
		backup.BackupID,
		"Device lost, need to recover identity",
	)
	
	suite.NoError(err)
	suite.NotNil(recoveryRequest)
	suite.NotEmpty(recoveryRequest.RequestID)
	suite.Equal(holderDID, recoveryRequest.HolderDID)
	suite.Equal(backup.BackupID, recoveryRequest.BackupID)
	suite.Equal(requester.String(), recoveryRequest.RequestedBy)
	suite.Equal(types.RecoveryRequestStatus_PENDING, recoveryRequest.Status)
	suite.True(recoveryRequest.RequiredScore > 0)
	suite.Equal(5, recoveryRequest.MaxAttempts)
	suite.Equal(0, recoveryRequest.AttemptCount)
	
	// Verify recovery request is stored
	storedRequest, found := k.GetRecoveryRequest(ctx, recoveryRequest.RequestID)
	suite.True(found)
	suite.Equal(recoveryRequest.RequestID, storedRequest.RequestID)
	
	// Verify recovery request can attempt recovery
	suite.True(recoveryRequest.CanAttemptRecovery())
	suite.False(recoveryRequest.HasSufficientConfidence())
}

// TestSubmitRecoveryProof tests submitting recovery proof
func (suite *RecoveryTestSuite) TestSubmitRecoveryProof() {
	ctx := suite.ctx
	k := suite.keeper
	holderDID := "did:desh:proof123"
	requester := suite.addrs[0]
	
	// Setup identity, backup, and recovery request
	identity := types.Identity{
		DID:        holderDID,
		Controller: requester.String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)
	
	recoveryMethods := []types.RecoveryMethod{
		{
			MethodID:           "mnemonic_1",
			MethodType:         types.RecoveryMethodType_MNEMONIC_PHRASE,
			MethodName:         "Mnemonic phrase",
			RequiredConfidence: 90,
			Enabled:            true,
			CreatedAt:          ctx.BlockTime(),
		},
	}
	
	backup, err := k.CreateIdentityBackup(
		ctx,
		holderDID,
		recoveryMethods,
		[]byte("test-key"),
		365*24*time.Hour,
	)
	suite.NoError(err)
	
	recoveryRequest, err := k.InitiateRecovery(
		ctx,
		requester,
		holderDID,
		backup.BackupID,
		"Testing proof submission",
	)
	suite.NoError(err)
	
	// Submit recovery proof
	proofData := []byte("valid-mnemonic-phrase-proof")
	verificationData := map[string]interface{}{
		"method": "mnemonic",
		"words":  12,
	}
	
	attempt, err := k.SubmitRecoveryProof(
		ctx,
		requester,
		recoveryRequest.RequestID,
		"mnemonic_1",
		proofData,
		verificationData,
	)
	
	suite.NoError(err)
	suite.NotNil(attempt)
	suite.NotEmpty(attempt.AttemptID)
	suite.Equal("mnemonic_1", attempt.MethodID)
	suite.Equal(types.RecoveryMethodType_MNEMONIC_PHRASE, attempt.MethodType)
	suite.Equal(types.AttemptStatus_VERIFIED, attempt.Status)
	suite.Equal(90, attempt.Confidence) // Mocked confidence from keeper
	
	// Verify recovery request is updated
	updatedRequest, found := k.GetRecoveryRequest(ctx, recoveryRequest.RequestID)
	suite.True(found)
	suite.Equal(1, updatedRequest.AttemptCount)
	suite.Len(updatedRequest.RecoveryMethods, 1)
	suite.Equal(90, updatedRequest.ConfidenceScore)
	
	// Should be approved now since confidence >= required
	suite.Equal(types.RecoveryRequestStatus_APPROVED, updatedRequest.Status)
	suite.True(updatedRequest.HasSufficientConfidence())
}

// TestExecuteRecovery tests executing recovery
func (suite *RecoveryTestSuite) TestExecuteRecovery() {
	ctx := suite.ctx
	k := suite.keeper
	holderDID := "did:desh:execute123"
	requester := suite.addrs[0]
	newController := suite.addrs[1]
	
	// Setup and approve recovery request
	identity := types.Identity{
		DID:        holderDID,
		Controller: requester.String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)
	
	recoveryMethods := []types.RecoveryMethod{
		{
			MethodID:           "mnemonic_1",
			MethodType:         types.RecoveryMethodType_MNEMONIC_PHRASE,
			RequiredConfidence: 90,
			Enabled:            true,
			CreatedAt:          ctx.BlockTime(),
		},
	}
	
	backup, _ := k.CreateIdentityBackup(ctx, holderDID, recoveryMethods, []byte("test-key"), 365*24*time.Hour)
	recoveryRequest, _ := k.InitiateRecovery(ctx, requester, holderDID, backup.BackupID, "Test execution")
	
	// Submit proof to approve request
	_, _ = k.SubmitRecoveryProof(ctx, requester, recoveryRequest.RequestID, "mnemonic_1", []byte("proof"), nil)
	
	// Execute recovery
	decryptionKey := []byte("decryption-key")
	err := k.ExecuteRecovery(
		ctx,
		requester,
		recoveryRequest.RequestID,
		newController,
		decryptionKey,
	)
	
	suite.NoError(err)
	
	// Verify recovery request is completed
	completedRequest, found := k.GetRecoveryRequest(ctx, recoveryRequest.RequestID)
	suite.True(found)
	suite.Equal(types.RecoveryRequestStatus_COMPLETED, completedRequest.Status)
}

// TestAddSocialRecoveryGuardian tests adding guardians
func (suite *RecoveryTestSuite) TestAddSocialRecoveryGuardian() {
	ctx := suite.ctx
	k := suite.keeper
	holderDID := "did:desh:guardian123"
	guardianDID := "did:desh:guardian456"
	
	// Create identities
	holderIdentity := types.Identity{
		DID:        holderDID,
		Controller: suite.addrs[0].String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, holderIdentity)
	
	guardianIdentity := types.Identity{
		DID:        guardianDID,
		Controller: suite.addrs[1].String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, guardianIdentity)
	
	// Add guardian
	guardian, err := k.AddSocialRecoveryGuardian(
		ctx,
		holderDID,
		guardianDID,
		suite.addrs[1].String(),
		"John Doe",
		10,
		"john@example.com",
		"guardian-public-key",
	)
	
	suite.NoError(err)
	suite.NotNil(guardian)
	suite.NotEmpty(guardian.GuardianID)
	suite.Equal(guardianDID, guardian.GuardianDID)
	suite.Equal(suite.addrs[1].String(), guardian.GuardianAddress)
	suite.Equal("John Doe", guardian.GuardianName)
	suite.Equal(10, guardian.Weight)
	suite.Equal("medium", guardian.TrustLevel)
	suite.Equal(types.GuardianStatus_ACTIVE, guardian.Status)
	suite.Contains(guardian.ContactInfo, "encrypted_") // Mock encryption
	
	// Verify guardian is stored
	storedGuardian, found := k.GetSocialRecoveryGuardian(ctx, guardian.GuardianID)
	suite.True(found)
	suite.Equal(guardian.GuardianID, storedGuardian.GuardianID)
	suite.True(guardian.IsActive())
}

// TestSubmitGuardianVote tests guardian voting
func (suite *RecoveryTestSuite) TestSubmitGuardianVote() {
	ctx := suite.ctx
	k := suite.keeper
	holderDID := "did:desh:vote123"
	guardianDID := "did:desh:guardianvote456"
	guardianAddr := suite.addrs[1]
	
	// Setup identities and guardian
	holderIdentity := types.Identity{
		DID:        holderDID,
		Controller: suite.addrs[0].String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, holderIdentity)
	
	guardianIdentity := types.Identity{
		DID:        guardianDID,
		Controller: guardianAddr.String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, guardianIdentity)
	
	guardian, _ := k.AddSocialRecoveryGuardian(
		ctx, holderDID, guardianDID, guardianAddr.String(),
		"Guardian", 10, "contact", "pubkey",
	)
	
	// Create recovery request
	recoveryMethods := []types.RecoveryMethod{{
		MethodID:           "social_1",
		MethodType:         types.RecoveryMethodType_SOCIAL_RECOVERY,
		RequiredConfidence: 75,
		Enabled:            true,
		CreatedAt:          ctx.BlockTime(),
	}}
	
	backup, _ := k.CreateIdentityBackup(ctx, holderDID, recoveryMethods, []byte("key"), 365*24*time.Hour)
	recoveryRequest, _ := k.InitiateRecovery(ctx, suite.addrs[0], holderDID, backup.BackupID, "Test voting")
	
	// Submit guardian vote
	vote, err := k.SubmitGuardianVote(
		ctx,
		guardianAddr,
		recoveryRequest.RequestID,
		types.VoteType_APPROVE,
		"Legitimate recovery request",
		"guardian-signature",
	)
	
	suite.NoError(err)
	suite.NotNil(vote)
	suite.NotEmpty(vote.VoteID)
	suite.Equal(recoveryRequest.RequestID, vote.RequestID)
	suite.Equal(guardian.GuardianID, vote.GuardianID)
	suite.Equal(types.VoteType_APPROVE, vote.Vote)
	suite.Equal("Legitimate recovery request", vote.Reason)
	suite.Equal(10, vote.Weight)
	
	// Verify vote is stored
	// Note: In production, votes would be indexed for easy retrieval
}

// TestVerifyBackupIntegrity tests backup verification
func (suite *RecoveryTestSuite) TestVerifyBackupIntegrity() {
	ctx := suite.ctx
	k := suite.keeper
	holderDID := "did:desh:verify123"
	
	// Create identity and backup
	identity := types.Identity{
		DID:        holderDID,
		Controller: suite.addrs[0].String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)
	
	recoveryMethods := []types.RecoveryMethod{{
		MethodID:           "method_1",
		MethodType:         types.RecoveryMethodType_MNEMONIC_PHRASE,
		RequiredConfidence: 90,
		Enabled:            true,
		CreatedAt:          ctx.BlockTime(),
	}}
	
	backup, err := k.CreateIdentityBackup(
		ctx,
		holderDID,
		recoveryMethods,
		[]byte("verification-key"),
		365*24*time.Hour,
	)
	suite.NoError(err)
	
	// Verify backup integrity
	verificationKey := []byte("verification-key")
	result, err := k.VerifyBackupIntegrity(
		ctx,
		backup.BackupID,
		verificationKey,
	)
	
	suite.NoError(err)
	suite.NotNil(result)
	suite.NotEmpty(result.VerificationID)
	suite.Equal(backup.BackupID, result.BackupID)
	suite.True(result.IntegrityValid)    // Mock implementation returns true
	suite.True(result.DecryptionValid)   // Mock implementation returns true
	suite.Equal(100.0, result.DataCompleteness) // Mock completeness
	suite.True(result.RecoverabilityScore > 0)
	
	// Should have no issues for a fresh backup
	suite.Len(result.IssuesFound, 0)
	suite.Len(result.Recommendations, 0)
}

// TestRecoveryRequestExpiry tests expired recovery requests
func (suite *RecoveryTestSuite) TestRecoveryRequestExpiry() {
	ctx := suite.ctx
	k := suite.keeper
	holderDID := "did:desh:expiry123"
	
	// Create identity and backup
	identity := types.Identity{
		DID:        holderDID,
		Controller: suite.addrs[0].String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)
	
	recoveryMethods := []types.RecoveryMethod{{
		MethodID:           "method_1",
		MethodType:         types.RecoveryMethodType_MNEMONIC_PHRASE,
		RequiredConfidence: 90,
		Enabled:            true,
		CreatedAt:          ctx.BlockTime(),
	}}
	
	backup, _ := k.CreateIdentityBackup(ctx, holderDID, recoveryMethods, []byte("key"), 365*24*time.Hour)
	
	// Create recovery request that expires quickly
	recoveryRequest, _ := k.InitiateRecovery(ctx, suite.addrs[0], holderDID, backup.BackupID, "Test expiry")
	
	// Manually set expiry to past (simulating time passage)
	recoveryRequest.ExpiresAt = ctx.BlockTime().Add(-time.Hour)
	k.SetRecoveryRequest(ctx, *recoveryRequest)
	
	// Verify request is expired
	suite.True(recoveryRequest.IsExpired())
	suite.False(recoveryRequest.CanAttemptRecovery())
	
	// Try to submit proof on expired request
	_, err := k.SubmitRecoveryProof(
		ctx,
		suite.addrs[0],
		recoveryRequest.RequestID,
		"method_1",
		[]byte("proof"),
		nil,
	)
	
	suite.Error(err)
	suite.Contains(err.Error(), "expired")
}

// TestBackupExpiry tests expired backups
func (suite *RecoveryTestSuite) TestBackupExpiry() {
	ctx := suite.ctx
	k := suite.keeper
	holderDID := "did:desh:backupexpiry123"
	
	// Create identity
	identity := types.Identity{
		DID:        holderDID,
		Controller: suite.addrs[0].String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)
	
	recoveryMethods := []types.RecoveryMethod{{
		MethodID:           "method_1",
		MethodType:         types.RecoveryMethodType_MNEMONIC_PHRASE,
		RequiredConfidence: 90,
		Enabled:            true,
		CreatedAt:          ctx.BlockTime(),
	}}
	
	// Create backup with short expiry
	backup, err := k.CreateIdentityBackup(
		ctx,
		holderDID,
		recoveryMethods,
		[]byte("key"),
		time.Nanosecond, // Very short expiry
	)
	suite.NoError(err)
	
	// Verify backup is expired
	suite.True(backup.IsExpired())
	suite.False(backup.IsValid())
	
	// Try to initiate recovery with expired backup
	_, err = k.InitiateRecovery(
		ctx,
		suite.addrs[0],
		holderDID,
		backup.BackupID,
		"Test with expired backup",
	)
	
	suite.Error(err)
	suite.Contains(err.Error(), "expired")
}

// TestInvalidRecoveryScenarios tests various invalid scenarios
func (suite *RecoveryTestSuite) TestInvalidRecoveryScenarios() {
	ctx := suite.ctx
	k := suite.keeper
	
	// Test with non-existent identity
	_, err := k.CreateIdentityBackup(
		ctx,
		"did:desh:nonexistent",
		[]types.RecoveryMethod{},
		[]byte("key"),
		time.Hour,
	)
	suite.Error(err)
	suite.Contains(err.Error(), "not found")
	
	// Test with non-existent backup
	_, err = k.InitiateRecovery(
		ctx,
		suite.addrs[0],
		"did:desh:holder",
		"non-existent-backup",
		"reason",
	)
	suite.Error(err)
	
	// Test with non-existent recovery request
	_, err = k.SubmitRecoveryProof(
		ctx,
		suite.addrs[0],
		"non-existent-request",
		"method",
		[]byte("proof"),
		nil,
	)
	suite.Error(err)
}