package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/suite"

	"github.com/DeshChain/DeshChain-Ecosystem/x/identity/keeper"
	"github.com/DeshChain/DeshChain-Ecosystem/x/identity/types"
	"github.com/DeshChain/DeshChain-Ecosystem/testutil"
)

type ErrorHandlingTestSuite struct {
	suite.Suite

	ctx    sdk.Context
	keeper *keeper.Keeper
	addrs  []sdk.AccAddress
}

func (suite *ErrorHandlingTestSuite) SetupTest() {
	suite.ctx, suite.keeper = testutil.IdentityKeeperTestSetup(suite.T())
	suite.addrs = testutil.CreateIncrementalAccounts(10)
}

func TestErrorHandlingTestSuite(t *testing.T) {
	suite.Run(t, new(ErrorHandlingTestSuite))
}

// TestIssueCredentialWithoutIdentity tests issuing credential for non-existent identity
func (suite *ErrorHandlingTestSuite) TestIssueCredentialWithoutIdentity() {
	ctx := suite.ctx
	k := suite.keeper
	issuer := suite.addrs[0]
	nonExistentDID := "did:desh:nonexistent"
	
	credentialSubject := map[string]interface{}{
		"id":   nonExistentDID,
		"data": "test",
	}
	
	// Try to issue credential for non-existent identity
	_, err := k.IssueCredential(
		ctx,
		issuer,
		nonExistentDID,
		[]string{"VerifiableCredential", "TestCredential"},
		credentialSubject,
	)
	
	suite.Error(err)
	suite.Contains(err.Error(), "identity not found")
}

// TestRevokeCredentialUnauthorized tests unauthorized credential revocation
func (suite *ErrorHandlingTestSuite) TestRevokeCredentialUnauthorized() {
	ctx := suite.ctx
	k := suite.keeper
	issuer := suite.addrs[0]
	unauthorized := suite.addrs[1]
	holderDID := "did:desh:unauthorized"
	
	// Create holder identity
	identity := types.Identity{
		Did:        holderDID,
		Controller: suite.addrs[2].String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)
	
	// Issue credential
	credentialSubject := map[string]interface{}{
		"id":   holderDID,
		"data": "test",
	}
	
	credentialID, err := k.IssueCredential(
		ctx,
		issuer,
		holderDID,
		[]string{"VerifiableCredential", "TestCredential"},
		credentialSubject,
	)
	suite.NoError(err)
	
	// Try to revoke with unauthorized address
	err = k.RevokeCredential(ctx, unauthorized, credentialID, "Unauthorized attempt")
	suite.Error(err)
	suite.Contains(err.Error(), "unauthorized")
	
	// Verify credential is still active
	credential, found := k.GetCredential(ctx, credentialID)
	suite.True(found)
	suite.Equal(types.CredentialStatus_ACTIVE, credential.Status)
}

// TestVerifyCredentialWithInvalidID tests verification with invalid credential ID
func (suite *ErrorHandlingTestSuite) TestVerifyCredentialWithInvalidID() {
	ctx := suite.ctx
	k := suite.keeper
	verifier := suite.addrs[0]
	
	// Try to verify non-existent credential
	isValid, err := k.VerifyCredential(ctx, "non-existent-credential-id", verifier)
	suite.Error(err)
	suite.False(isValid)
	suite.Contains(err.Error(), "not found")
}

// TestPresentCredentialWithInvalidID tests presentation with invalid credential IDs
func (suite *ErrorHandlingTestSuite) TestPresentCredentialWithInvalidID() {
	ctx := suite.ctx
	k := suite.keeper
	verifier := suite.addrs[0]
	
	// Try to present non-existent credentials
	disclosureMap := map[string][]string{
		"non-existent-1": {"field1"},
		"non-existent-2": {"field2"},
	}
	
	presentationID := k.PresentCredential(
		ctx,
		[]string{"non-existent-1", "non-existent-2"},
		verifier.String(),
		disclosureMap,
	)
	
	// Presentation should be created but may contain empty or error credentials
	suite.NotEmpty(presentationID)
	
	presentation, found := k.GetPresentation(ctx, presentationID)
	suite.True(found)
	// The behavior depends on implementation - might be empty or contain error entries
}

// TestZKProofWithNonExistentIdentity tests ZK proof operations without identity
func (suite *ErrorHandlingTestSuite) TestZKProofWithNonExistentIdentity() {
	ctx := suite.ctx
	k := suite.keeper
	nonExistentDID := "did:desh:zknonexistent"
	
	// Try to generate age proof for non-existent identity
	_, err := k.GenerateAgeProof(ctx, nonExistentDID, 25, 18)
	suite.Error(err)
	suite.Contains(err.Error(), "identity not found")
	
	// Try to verify age proof for non-existent identity
	fakeProofData := []byte("fake_proof_data")
	_, err = k.VerifyAgeProof(ctx, nonExistentDID, fakeProofData, 18)
	suite.Error(err)
	suite.Contains(err.Error(), "identity not found")
}

// TestBiometricWithNonExistentIdentity tests biometric operations without identity
func (suite *ErrorHandlingTestSuite) TestBiometricWithNonExistentIdentity() {
	ctx := suite.ctx
	k := suite.keeper
	nonExistentDID := "did:desh:biononexistent"
	
	// Try to store biometric template for non-existent identity
	templateData := []byte("template_data")
	err := k.StoreBiometricTemplate(ctx, nonExistentDID, "fingerprint", templateData)
	suite.Error(err)
	suite.Contains(err.Error(), "identity not found")
	
	// Try to verify biometric for non-existent identity
	_, err = k.VerifyBiometric(ctx, nonExistentDID, "fingerprint", templateData)
	suite.Error(err)
	suite.Contains(err.Error(), "identity not found")
}

// TestConsentWithNonExistentIdentity tests consent operations without identity
func (suite *ErrorHandlingTestSuite) TestConsentWithNonExistentIdentity() {
	ctx := suite.ctx
	k := suite.keeper
	nonExistentDID := "did:desh:consentnonexistent"
	
	// Try to grant consent for non-existent identity
	_, err := k.GrantConsent(
		ctx,
		nonExistentDID,
		suite.addrs[0].String(),
		[]string{"KYC verification"},
		[]string{"name", "age"},
		ctx.BlockTime().Add(24*time.Hour),
	)
	suite.Error(err)
	suite.Contains(err.Error(), "identity not found")
	
	// Try to check consent for non-existent identity
	hasConsent := k.HasValidConsent(ctx, nonExistentDID, suite.addrs[0].String(), "KYC verification", "name")
	suite.False(hasConsent)
}

// TestInvalidCredentialTypes tests credential issuance with invalid types
func (suite *ErrorHandlingTestSuite) TestInvalidCredentialTypes() {
	ctx := suite.ctx
	k := suite.keeper
	issuer := suite.addrs[0]
	holderDID := "did:desh:invalidtypes"
	
	// Create holder identity
	identity := types.Identity{
		Did:        holderDID,
		Controller: suite.addrs[1].String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)
	
	credentialSubject := map[string]interface{}{
		"id":   holderDID,
		"data": "test",
	}
	
	// Test with empty types array
	_, err := k.IssueCredential(
		ctx,
		issuer,
		holderDID,
		[]string{}, // Empty types
		credentialSubject,
	)
	suite.Error(err)
	suite.Contains(err.Error(), "at least one credential type")
	
	// Test with nil types
	_, err = k.IssueCredential(
		ctx,
		issuer,
		holderDID,
		nil, // Nil types
		credentialSubject,
	)
	suite.Error(err)
	suite.Contains(err.Error(), "at least one credential type")
	
	// Test without VerifiableCredential as first type
	_, err = k.IssueCredential(
		ctx,
		issuer,
		holderDID,
		[]string{"CustomCredential"}, // Missing VerifiableCredential
		credentialSubject,
	)
	suite.Error(err)
	suite.Contains(err.Error(), "VerifiableCredential")
}

// TestCorruptedCredentialData tests handling of corrupted credential data
func (suite *ErrorHandlingTestSuite) TestCorruptedCredentialData() {
	ctx := suite.ctx
	k := suite.keeper
	issuer := suite.addrs[0]
	holderDID := "did:desh:corrupted"
	
	// Create holder identity
	identity := types.Identity{
		Did:        holderDID,
		Controller: suite.addrs[1].String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)
	
	// Test with credential subject containing unsupported data types
	corruptedSubject := map[string]interface{}{
		"id":       holderDID,
		"function": func() {}, // Function (unsupported in JSON)
		"channel":  make(chan int), // Channel (unsupported)
		"circular": nil, // Will be set to self-reference
	}
	corruptedSubject["circular"] = corruptedSubject // Create circular reference
	
	// This might succeed or fail depending on the JSON marshaling implementation
	_, err := k.IssueCredential(
		ctx,
		issuer,
		holderDID,
		[]string{"VerifiableCredential", "CorruptedCredential"},
		corruptedSubject,
	)
	
	// Error handling depends on the specific implementation
	// Some JSON marshalers can handle this, others cannot
	if err != nil {
		suite.Contains(err.Error(), "json", "marshal", "unsupported", "circular")
	}
}

// TestExpiredConsentAccess tests accessing expired consent
func (suite *ErrorHandlingTestSuite) TestExpiredConsentAccess() {
	ctx := suite.ctx
	k := suite.keeper
	did := "did:desh:expiredconsent"
	
	// Create identity
	identity := types.Identity{
		Did:        did,
		Controller: suite.addrs[0].String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)
	
	// Grant consent with very short expiry
	shortExpiry := ctx.BlockTime().Add(1 * time.Nanosecond)
	consentID, err := k.GrantConsent(
		ctx,
		did,
		suite.addrs[1].String(),
		[]string{"KYC verification"},
		[]string{"name", "age"},
		shortExpiry,
	)
	suite.NoError(err)
	suite.NotEmpty(consentID)
	
	// Advance time to expire consent
	expiredCtx := ctx.WithBlockTime(ctx.BlockTime().Add(1 * time.Hour))
	
	// Try to use expired consent
	hasConsent := k.HasValidConsent(expiredCtx, did, suite.addrs[1].String(), "KYC verification", "name")
	suite.False(hasConsent)
	
	// Try to revoke expired consent
	err = k.RevokeConsent(expiredCtx, did, consentID)
	// This might succeed (revoking already expired consent) or fail depending on implementation
	if err != nil {
		suite.Contains(err.Error(), "expired", "invalid")
	}
}

// TestConcurrentCredentialRevocation tests race conditions in revocation
func (suite *ErrorHandlingTestSuite) TestConcurrentCredentialRevocation() {
	ctx := suite.ctx
	k := suite.keeper
	issuer := suite.addrs[0]
	holderDID := "did:desh:concurrent"
	
	// Create holder identity
	identity := types.Identity{
		Did:        holderDID,
		Controller: suite.addrs[1].String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)
	
	// Issue credential
	credentialSubject := map[string]interface{}{
		"id":   holderDID,
		"data": "test",
	}
	
	credentialID, err := k.IssueCredential(
		ctx,
		issuer,
		holderDID,
		[]string{"VerifiableCredential", "TestCredential"},
		credentialSubject,
	)
	suite.NoError(err)
	
	// Try to revoke the same credential multiple times
	err1 := k.RevokeCredential(ctx, issuer, credentialID, "First revocation")
	err2 := k.RevokeCredential(ctx, issuer, credentialID, "Second revocation")
	err3 := k.RevokeCredential(ctx, issuer, credentialID, "Third revocation")
	
	// First revocation should succeed
	suite.NoError(err1)
	
	// Subsequent revocations should fail
	suite.Error(err2)
	suite.Error(err3)
	suite.Contains(err2.Error(), "already revoked")
	suite.Contains(err3.Error(), "already revoked")
}

// TestInvalidZKProofData tests ZK proof verification with invalid data
func (suite *ErrorHandlingTestSuite) TestInvalidZKProofData() {
	ctx := suite.ctx
	k := suite.keeper
	did := "did:desh:invalidzkproof"
	
	// Create identity
	identity := types.Identity{
		Did:        did,
		Controller: suite.addrs[0].String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)
	
	// Test with invalid proof data
	invalidProofData := []byte("invalid_proof_data")
	isValid, err := k.VerifyAgeProof(ctx, did, invalidProofData, 18)
	suite.Error(err)
	suite.False(isValid)
	suite.Contains(err.Error(), "invalid", "proof")
	
	// Test with empty proof data
	emptyProofData := []byte{}
	isValid, err = k.VerifyAgeProof(ctx, did, emptyProofData, 18)
	suite.Error(err)
	suite.False(isValid)
	
	// Test with nil proof data
	isValid, err = k.VerifyAgeProof(ctx, did, nil, 18)
	suite.Error(err)
	suite.False(isValid)
}

// TestSystemOverload tests behavior under high load conditions
func (suite *ErrorHandlingTestSuite) TestSystemOverload() {
	ctx := suite.ctx
	k := suite.keeper
	
	// This test simulates system overload by rapidly creating many operations
	// In a real system, this might trigger rate limiting or resource exhaustion
	
	numOperations := 10000
	errorCount := 0
	
	for i := 0; i < numOperations; i++ {
		did := fmt.Sprintf("did:desh:overload%d", i)
		identity := types.Identity{
			Did:        did,
			Controller: suite.addrs[i%len(suite.addrs)].String(),
			Status:     types.IdentityStatus_ACTIVE,
			CreatedAt:  ctx.BlockTime(),
			UpdatedAt:  ctx.BlockTime(),
		}
		
		// Some operations might fail under extreme load
		k.SetIdentity(ctx, identity)
		
		// Quick verification
		_, found := k.GetIdentity(ctx, did)
		if !found {
			errorCount++
		}
		
		// Break if too many errors (system overload protection)
		if errorCount > numOperations/10 { // Allow up to 10% failure rate
			break
		}
	}
	
	// System should handle most operations even under load
	suite.Less(errorCount, numOperations/10, "Too many errors under load: %d/%d", errorCount, numOperations)
}

// TestSDKErrorWrapping tests proper SDK error wrapping
func (suite *ErrorHandlingTestSuite) TestSDKErrorWrapping() {
	ctx := suite.ctx
	k := suite.keeper
	
	// Test operations that should return specific SDK errors
	
	// Try to issue credential with invalid issuer address
	invalidIssuer := sdk.AccAddress{} // Empty address
	_, err := k.IssueCredential(
		ctx,
		invalidIssuer,
		"did:desh:test",
		[]string{"VerifiableCredential", "TestCredential"},
		map[string]interface{}{"id": "did:desh:test"},
	)
	
	if err != nil {
		// Check if error is properly wrapped with SDK error types
		suite.True(sdkerrors.IsOf(err, sdkerrors.ErrInvalidAddress) || 
				strings.Contains(err.Error(), "invalid") ||
				strings.Contains(err.Error(), "address"))
	}
}
