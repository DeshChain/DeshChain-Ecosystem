package keeper_test

import (
	"strings"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/DeshChain/DeshChain-Ecosystem/x/identity/keeper"
	"github.com/DeshChain/DeshChain-Ecosystem/x/identity/types"
	"github.com/DeshChain/DeshChain-Ecosystem/testutil"
)

type EdgeCasesTestSuite struct {
	suite.Suite

	ctx    sdk.Context
	keeper *keeper.Keeper
	addrs  []sdk.AccAddress
}

func (suite *EdgeCasesTestSuite) SetupTest() {
	suite.ctx, suite.keeper = testutil.IdentityKeeperTestSetup(suite.T())
	suite.addrs = testutil.CreateIncrementalAccounts(10)
}

func TestEdgeCasesTestSuite(t *testing.T) {
	suite.Run(t, new(EdgeCasesTestSuite))
}

// TestDuplicateIdentityCreation tests handling of duplicate identity creation
func (suite *EdgeCasesTestSuite) TestDuplicateIdentityCreation() {
	ctx := suite.ctx
	k := suite.keeper
	addr := suite.addrs[0]
	
	did := "did:desh:duplicate123"
	identity := types.Identity{
		Did:        did,
		Controller: addr.String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	
	// Create first identity
	k.SetIdentity(ctx, identity)
	
	// Verify it exists
	storedIdentity, found := k.GetIdentity(ctx, did)
	suite.True(found)
	suite.Equal(identity.Did, storedIdentity.Did)
	
	// Try to create duplicate with different controller
	duplicateIdentity := identity
	duplicateIdentity.Controller = suite.addrs[1].String()
	duplicateIdentity.UpdatedAt = ctx.BlockTime().Add(time.Hour)
	
	// SetIdentity should overwrite (this is expected behavior for updates)
	k.SetIdentity(ctx, duplicateIdentity)
	
	// Verify the identity was updated
	updatedIdentity, found := k.GetIdentity(ctx, did)
	suite.True(found)
	suite.Equal(suite.addrs[1].String(), updatedIdentity.Controller)
	suite.Equal(duplicateIdentity.UpdatedAt, updatedIdentity.UpdatedAt)
}

// TestInvalidDIDFormats tests handling of invalid DID formats
func (suite *EdgeCasesTestSuite) TestInvalidDIDFormats() {
	ctx := suite.ctx
	k := suite.keeper
	addr := suite.addrs[0]
	
	invalidDIDs := []string{
		"",                           // Empty DID
		"not-a-did",                // Not a DID format
		"did:",                      // Incomplete DID
		"did:invalid",              // Missing identifier
		"did:desh:",                // Empty identifier
		"did::missing-method",      // Missing method
		"invalid:desh:test",        // Wrong scheme
		strings.Repeat("x", 1000),  // Extremely long string
		"did:desh:" + strings.Repeat("x", 500), // Long but valid structure
	}
	
	for _, invalidDID := range invalidDIDs {
		identity := types.Identity{
			Did:        invalidDID,
			Controller: addr.String(),
			Status:     types.IdentityStatus_ACTIVE,
			CreatedAt:  ctx.BlockTime(),
			UpdatedAt:  ctx.BlockTime(),
		}
		
		// Store the identity (keeper doesn't validate DID format)
		k.SetIdentity(ctx, identity)
		
		// Verify it can be retrieved
		storedIdentity, found := k.GetIdentity(ctx, invalidDID)
		if invalidDID != "" { // Empty string might not be stored
			suite.True(found, "Should be able to store DID: %s", invalidDID)
			suite.Equal(invalidDID, storedIdentity.Did)
		}
	}
}

// TestCredentialWithEmptySubject tests credentials with empty or nil subjects
func (suite *EdgeCasesTestSuite) TestCredentialWithEmptySubject() {
	ctx := suite.ctx
	k := suite.keeper
	issuer := suite.addrs[0]
	holderDID := "did:desh:emptysubject"
	
	// Create holder identity
	identity := types.Identity{
		Did:        holderDID,
		Controller: suite.addrs[1].String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)
	
	// Test with nil subject
	credentialID1, err := k.IssueCredential(
		ctx,
		issuer,
		holderDID,
		[]string{"VerifiableCredential", "EmptySubjectCredential"},
		nil,
	)
	suite.NoError(err)
	suite.NotEmpty(credentialID1)
	
	// Test with empty subject
	emptySubject := make(map[string]interface{})
	credentialID2, err := k.IssueCredential(
		ctx,
		issuer,
		holderDID,
		[]string{"VerifiableCredential", "EmptySubjectCredential"},
		emptySubject,
	)
	suite.NoError(err)
	suite.NotEmpty(credentialID2)
	
	// Verify credentials exist and can be retrieved
	cred1, found1 := k.GetCredential(ctx, credentialID1)
	suite.True(found1)
	suite.NotNil(cred1.CredentialSubject) // Should be initialized
	
	cred2, found2 := k.GetCredential(ctx, credentialID2)
	suite.True(found2)
	suite.NotNil(cred2.CredentialSubject)
	suite.Empty(cred2.CredentialSubject)
}

// TestCredentialWithLargeSubject tests credentials with very large subjects
func (suite *EdgeCasesTestSuite) TestCredentialWithLargeSubject() {
	ctx := suite.ctx
	k := suite.keeper
	issuer := suite.addrs[0]
	holderDID := "did:desh:largesubject"
	
	// Create holder identity
	identity := types.Identity{
		Did:        holderDID,
		Controller: suite.addrs[1].String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)
	
	// Create large credential subject
	largeSubject := map[string]interface{}{
		"id": holderDID,
	}
	
	// Add many fields
	for i := 0; i < 1000; i++ {
		largeSubject["field_"+fmt.Sprintf("%d", i)] = strings.Repeat("data", 100)
	}
	
	// Add deeply nested object
	nestedData := make(map[string]interface{})
	for depth := 0; depth < 10; depth++ {
		level := make(map[string]interface{})
		for j := 0; j < 50; j++ {
			level["item_"+fmt.Sprintf("%d", j)] = "value_" + fmt.Sprintf("%d", j)
		}
		nestedData["level_"+fmt.Sprintf("%d", depth)] = level
	}
	largeSubject["nested_data"] = nestedData
	
	// Issue credential with large subject
	credentialID, err := k.IssueCredential(
		ctx,
		issuer,
		holderDID,
		[]string{"VerifiableCredential", "LargeSubjectCredential"},
		largeSubject,
	)
	suite.NoError(err)
	suite.NotEmpty(credentialID)
	
	// Verify credential can be retrieved
	credential, found := k.GetCredential(ctx, credentialID)
	suite.True(found)
	suite.Equal(holderDID, credential.CredentialSubject["id"])
	suite.Contains(credential.CredentialSubject, "nested_data")
}

// TestRevokeNonExistentCredential tests revoking a credential that doesn't exist
func (suite *EdgeCasesTestSuite) TestRevokeNonExistentCredential() {
	ctx := suite.ctx
	k := suite.keeper
	issuer := suite.addrs[0]
	
	nonExistentCredentialID := "non-existent-credential-id"
	
	// Try to revoke non-existent credential
	err := k.RevokeCredential(ctx, issuer, nonExistentCredentialID, "Testing non-existent")
	suite.Error(err)
	suite.Contains(err.Error(), "not found")
}

// TestVerifyRevokedCredential tests verifying an already revoked credential
func (suite *EdgeCasesTestSuite) TestVerifyRevokedCredential() {
	ctx := suite.ctx
	k := suite.keeper
	issuer := suite.addrs[0]
	verifier := suite.addrs[1]
	holderDID := "did:desh:revoked123"
	
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
		"data": "test data",
	}
	
	credentialID, err := k.IssueCredential(
		ctx,
		issuer,
		holderDID,
		[]string{"VerifiableCredential", "TestCredential"},
		credentialSubject,
	)
	suite.NoError(err)
	
	// Revoke credential
	err = k.RevokeCredential(ctx, issuer, credentialID, "Test revocation")
	suite.NoError(err)
	
	// Try to revoke again
	err = k.RevokeCredential(ctx, issuer, credentialID, "Second revocation")
	suite.Error(err)
	suite.Contains(err.Error(), "already revoked")
	
	// Verify credential is invalid
	isValid, err := k.VerifyCredential(ctx, credentialID, verifier)
	suite.NoError(err)
	suite.False(isValid)
}

// TestPresentRevokedCredential tests presenting a revoked credential
func (suite *EdgeCasesTestSuite) TestPresentRevokedCredential() {
	ctx := suite.ctx
	k := suite.keeper
	issuer := suite.addrs[0]
	verifier := suite.addrs[1]
	holderDID := "did:desh:presentrevoked"
	
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
		"name": "John Doe",
		"age":  30,
	}
	
	credentialID, err := k.IssueCredential(
		ctx,
		issuer,
		holderDID,
		[]string{"VerifiableCredential", "TestCredential"},
		credentialSubject,
	)
	suite.NoError(err)
	
	// Revoke credential
	err = k.RevokeCredential(ctx, issuer, credentialID, "Test revocation")
	suite.NoError(err)
	
	// Try to present revoked credential
	disclosureMap := map[string][]string{
		credentialID: {"name", "age"},
	}
	
	presentationID := k.PresentCredential(
		ctx,
		[]string{credentialID},
		verifier.String(),
		disclosureMap,
	)
	
	// Presentation should be created but contain revoked credential
	suite.NotEmpty(presentationID)
	
	presentation, found := k.GetPresentation(ctx, presentationID)
	suite.True(found)
	suite.Len(presentation.VerifiableCredentials, 1)
	
	// Verify the credential in presentation is marked as revoked
	presentedCred := presentation.VerifiableCredentials[0]
	suite.Equal(types.CredentialStatus_REVOKED, presentedCred.Status)
}

// TestZKProofWithInvalidAge tests ZK proof generation with edge case ages
func (suite *EdgeCasesTestSuite) TestZKProofWithInvalidAge() {
	ctx := suite.ctx
	k := suite.keeper
	did := "did:desh:invalidage"
	
	// Create identity
	identity := types.Identity{
		Did:        did,
		Controller: suite.addrs[0].String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)
	
	// Test with zero age
	_, err := k.GenerateAgeProof(ctx, did, 0, 18)
	suite.Error(err)
	
	// Test with negative age (will be treated as very large uint32)
	_, err = k.GenerateAgeProof(ctx, did, uint32(-1), 18)
	suite.NoError(err) // Very large number should work
	
	// Test with min age greater than actual age (should fail verification)
	proofData, err := k.GenerateAgeProof(ctx, did, 16, 18)
	suite.NoError(err)
	
	isValid, err := k.VerifyAgeProof(ctx, did, proofData, 18)
	suite.NoError(err)
	suite.False(isValid)
	
	// Test with same age as minimum
	proofData, err = k.GenerateAgeProof(ctx, did, 18, 18)
	suite.NoError(err)
	
	isValid, err = k.VerifyAgeProof(ctx, did, proofData, 18)
	suite.NoError(err)
	suite.True(isValid)
}

// TestBiometricWithInvalidData tests biometric operations with invalid data
func (suite *EdgeCasesTestSuite) TestBiometricWithInvalidData() {
	ctx := suite.ctx
	k := suite.keeper
	did := "did:desh:invalidbio"
	
	// Create identity
	identity := types.Identity{
		Did:        did,
		Controller: suite.addrs[0].String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)
	
	// Test with empty template data
	err := k.StoreBiometricTemplate(ctx, did, "fingerprint", []byte{})
	suite.NoError(err) // Empty data should be allowed
	
	// Test with nil template data
	err = k.StoreBiometricTemplate(ctx, did, "face", nil)
	suite.NoError(err) // Nil data should be allowed
	
	// Test with empty template type
	err = k.StoreBiometricTemplate(ctx, did, "", []byte("data"))
	suite.NoError(err) // Empty type should be allowed
	
	// Test verification with mismatched data
	err = k.StoreBiometricTemplate(ctx, did, "iris", []byte("original_data"))
	suite.NoError(err)
	
	isMatch, err := k.VerifyBiometric(ctx, did, "iris", []byte("different_data"))
	suite.NoError(err)
	suite.False(isMatch)
	
	// Test verification with non-existent template type
	isMatch, err = k.VerifyBiometric(ctx, did, "non_existent_type", []byte("data"))
	suite.Error(err)
	suite.False(isMatch)
}

// TestConsentWithInvalidData tests consent management with edge cases
func (suite *EdgeCasesTestSuite) TestConsentWithInvalidData() {
	ctx := suite.ctx
	k := suite.keeper
	did := "did:desh:invalidconsent"
	
	// Create identity
	identity := types.Identity{
		Did:        did,
		Controller: suite.addrs[0].String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)
	
	// Test with empty purposes
	consentID, err := k.GrantConsent(
		ctx,
		did,
		suite.addrs[1].String(),
		[]string{}, // Empty purposes
		[]string{"name", "age"},
		ctx.BlockTime().Add(24*time.Hour),
	)
	suite.NoError(err) // Should allow empty purposes
	suite.NotEmpty(consentID)
	
	// Test with empty data types
	consentID2, err := k.GrantConsent(
		ctx,
		did,
		suite.addrs[1].String(),
		[]string{"KYC verification"},
		[]string{}, // Empty data types
		ctx.BlockTime().Add(24*time.Hour),
	)
	suite.NoError(err) // Should allow empty data types
	suite.NotEmpty(consentID2)
	
	// Test with past expiry time
	pastTime := ctx.BlockTime().Add(-24 * time.Hour)
	consentID3, err := k.GrantConsent(
		ctx,
		did,
		suite.addrs[1].String(),
		[]string{"KYC verification"},
		[]string{"name"},
		pastTime,
	)
	suite.NoError(err) // Should allow past time (immediate expiry)
	suite.NotEmpty(consentID3)
	
	// Verify expired consent is not valid
	hasConsent := k.HasValidConsent(ctx, did, suite.addrs[1].String(), "KYC verification", "name")
	suite.False(hasConsent)
	
	// Test revoking non-existent consent
	err = k.RevokeConsent(ctx, did, "non-existent-consent-id")
	suite.Error(err)
	suite.Contains(err.Error(), "not found")
}

// TestMaximumCredentialTypes tests system limits
func (suite *EdgeCasesTestSuite) TestMaximumCredentialTypes() {
	ctx := suite.ctx
	k := suite.keeper
	issuer := suite.addrs[0]
	holderDID := "did:desh:maxtypes"
	
	// Create holder identity
	identity := types.Identity{
		Did:        holderDID,
		Controller: suite.addrs[1].String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)
	
	// Issue credential with many types
	manyTypes := make([]string, 100)
	manyTypes[0] = "VerifiableCredential" // First type must be VerifiableCredential
	for i := 1; i < 100; i++ {
		manyTypes[i] = "Type" + fmt.Sprintf("%d", i)
	}
	
	credentialSubject := map[string]interface{}{
		"id":   holderDID,
		"data": "test",
	}
	
	credentialID, err := k.IssueCredential(
		ctx,
		issuer,
		holderDID,
		manyTypes,
		credentialSubject,
	)
	suite.NoError(err)
	suite.NotEmpty(credentialID)
	
	// Verify credential was stored with all types
	credential, found := k.GetCredential(ctx, credentialID)
	suite.True(found)
	suite.Len(credential.Type, 100)
	suite.Equal("VerifiableCredential", credential.Type[0])
}

// TestIdentityStatusTransitions tests various identity status transitions
func (suite *EdgeCasesTestSuite) TestIdentityStatusTransitions() {
	ctx := suite.ctx
	k := suite.keeper
	did := "did:desh:statustransition"
	
	// Test all possible status transitions
	statuses := []types.IdentityStatus{
		types.IdentityStatus_ACTIVE,
		types.IdentityStatus_SUSPENDED,
		types.IdentityStatus_REVOKED,
		types.IdentityStatus_PENDING,
	}
	
	for i, status := range statuses {
		identity := types.Identity{
			Did:        did,
			Controller: suite.addrs[0].String(),
			Status:     status,
			CreatedAt:  ctx.BlockTime(),
			UpdatedAt:  ctx.BlockTime().Add(time.Duration(i) * time.Hour),
		}
		
		k.SetIdentity(ctx, identity)
		
		// Verify status was set
		storedIdentity, found := k.GetIdentity(ctx, did)
		suite.True(found)
		suite.Equal(status, storedIdentity.Status)
	}
}
