package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/deshchain/deshchain/x/identity/keeper"
	"github.com/deshchain/deshchain/x/identity/types"
	"github.com/deshchain/deshchain/testutil"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx    sdk.Context
	keeper *keeper.Keeper
	addrs  []sdk.AccAddress
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.ctx, suite.keeper = testutil.IdentityKeeperTestSetup(suite.T())
	suite.addrs = testutil.CreateIncrementalAccounts(5)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

// TestCreateIdentity tests basic identity creation
func (suite *KeeperTestSuite) TestCreateIdentity() {
	ctx := suite.ctx
	k := suite.keeper
	addr := suite.addrs[0]

	// Test creating a new identity
	did := "did:desh:test123"
	identity := types.Identity{
		Did:        did,
		Controller: addr.String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}

	// Store the identity
	k.SetIdentity(ctx, identity)

	// Retrieve and verify
	storedIdentity, found := k.GetIdentity(ctx, did)
	suite.True(found)
	suite.Equal(identity.Did, storedIdentity.Did)
	suite.Equal(identity.Controller, storedIdentity.Controller)
	suite.Equal(identity.Status, storedIdentity.Status)

	// Test retrieval by address
	identityByAddr, found := k.GetIdentityByAddress(ctx, addr)
	suite.True(found)
	suite.Equal(identity.Did, identityByAddr.Did)
}

// TestIdentityNotFound tests handling of non-existent identities
func (suite *KeeperTestSuite) TestIdentityNotFound() {
	ctx := suite.ctx
	k := suite.keeper

	// Test non-existent DID
	_, found := k.GetIdentity(ctx, "did:desh:nonexistent")
	suite.False(found)

	// Test non-existent address
	_, found = k.GetIdentityByAddress(ctx, suite.addrs[0])
	suite.False(found)
}

// TestUpdateIdentity tests identity updates
func (suite *KeeperTestSuite) TestUpdateIdentity() {
	ctx := suite.ctx
	k := suite.keeper
	addr := suite.addrs[0]

	// Create initial identity
	did := "did:desh:update123"
	identity := types.Identity{
		Did:        did,
		Controller: addr.String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)

	// Update the identity
	updatedTime := ctx.BlockTime().Add(time.Hour)
	ctx = ctx.WithBlockTime(updatedTime)
	identity.Status = types.IdentityStatus_SUSPENDED
	identity.UpdatedAt = updatedTime
	k.SetIdentity(ctx, identity)

	// Verify update
	storedIdentity, found := k.GetIdentity(ctx, did)
	suite.True(found)
	suite.Equal(types.IdentityStatus_SUSPENDED, storedIdentity.Status)
	suite.Equal(updatedTime, storedIdentity.UpdatedAt)
}

// TestIssueCredential tests credential issuance
func (suite *KeeperTestSuite) TestIssueCredential() {
	ctx := suite.ctx
	k := suite.keeper
	issuer := suite.addrs[0]
	holder := suite.addrs[1]

	// Create holder identity first
	holderDID := "did:desh:holder123"
	identity := types.Identity{
		Did:        holderDID,
		Controller: holder.String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)

	// Issue a credential
	credentialSubject := map[string]interface{}{
		"id":       holderDID,
		"name":     "John Doe",
		"age":      30,
		"verified": true,
	}

	credentialID, err := k.IssueCredential(
		ctx,
		issuer,
		holderDID,
		[]string{"VerifiableCredential", "KYCCredential"},
		credentialSubject,
	)

	suite.NoError(err)
	suite.NotEmpty(credentialID)

	// Retrieve and verify credential
	credential, found := k.GetCredential(ctx, credentialID)
	suite.True(found)
	suite.Equal([]string{"VerifiableCredential", "KYCCredential"}, credential.Type)
	suite.Equal(types.CredentialStatus_ACTIVE, credential.Status)
	suite.Equal(issuer.String(), credential.Issuer)
	suite.Equal(holderDID, credential.CredentialSubject["id"])
}

// TestCredentialsByType tests retrieving credentials by type
func (suite *KeeperTestSuite) TestCredentialsByType() {
	ctx := suite.ctx
	k := suite.keeper
	issuer := suite.addrs[0]
	holder := suite.addrs[1]

	// Create holder identity
	holderDID := "did:desh:holder456"
	identity := types.Identity{
		Did:        holderDID,
		Controller: holder.String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)

	// Issue multiple credentials of different types
	kycSubject := map[string]interface{}{
		"id":         holderDID,
		"kyc_level":  "enhanced",
		"verified":   true,
	}
	
	biometricSubject := map[string]interface{}{
		"id":              holderDID,
		"biometric_hash":  "hash123",
		"template_type":   "fingerprint",
	}

	// Issue KYC credential
	kycCredID, err := k.IssueCredential(
		ctx,
		issuer,
		holderDID,
		[]string{"VerifiableCredential", "KYCCredential"},
		kycSubject,
	)
	suite.NoError(err)

	// Issue Biometric credential
	bioCredID, err := k.IssueCredential(
		ctx,
		issuer,
		holderDID,
		[]string{"VerifiableCredential", "BiometricCredential"},
		biometricSubject,
	)
	suite.NoError(err)

	// Test retrieval by type
	kycCreds := k.GetCredentialsByType(ctx, holderDID, "KYCCredential")
	suite.Len(kycCreds, 1)
	suite.Equal(kycCredID, kycCreds[0].Id)

	bioCreds := k.GetCredentialsByType(ctx, holderDID, "BiometricCredential")
	suite.Len(bioCreds, 1)
	suite.Equal(bioCredID, bioCreds[0].Id)

	// Test non-existent type
	nonExistentCreds := k.GetCredentialsByType(ctx, holderDID, "NonExistentCredential")
	suite.Len(nonExistentCreds, 0)
}

// TestRevokeCredential tests credential revocation
func (suite *KeeperTestSuite) TestRevokeCredential() {
	ctx := suite.ctx
	k := suite.keeper
	issuer := suite.addrs[0]
	holder := suite.addrs[1]

	// Create holder identity and issue credential
	holderDID := "did:desh:revoke123"
	identity := types.Identity{
		Did:        holderDID,
		Controller: holder.String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)

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

	// Verify credential is active
	credential, found := k.GetCredential(ctx, credentialID)
	suite.True(found)
	suite.Equal(types.CredentialStatus_ACTIVE, credential.Status)

	// Revoke the credential
	err = k.RevokeCredential(ctx, issuer, credentialID, "Testing revocation")
	suite.NoError(err)

	// Verify credential is revoked
	revokedCredential, found := k.GetCredential(ctx, credentialID)
	suite.True(found)
	suite.Equal(types.CredentialStatus_REVOKED, revokedCredential.Status)
	suite.Equal("Testing revocation", revokedCredential.RevocationReason)
}

// TestUnauthorizedRevocation tests that only issuers can revoke credentials
func (suite *KeeperTestSuite) TestUnauthorizedRevocation() {
	ctx := suite.ctx
	k := suite.keeper
	issuer := suite.addrs[0]
	holder := suite.addrs[1]
	unauthorizedAddr := suite.addrs[2]

	// Create holder identity and issue credential
	holderDID := "did:desh:unauthorized123"
	identity := types.Identity{
		Did:        holderDID,
		Controller: holder.String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)

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

	// Try to revoke with unauthorized address
	err = k.RevokeCredential(ctx, unauthorizedAddr, credentialID, "Unauthorized revocation")
	suite.Error(err)
	suite.Contains(err.Error(), "unauthorized")

	// Verify credential is still active
	credential, found := k.GetCredential(ctx, credentialID)
	suite.True(found)
	suite.Equal(types.CredentialStatus_ACTIVE, credential.Status)
}

// TestVerifyCredential tests credential verification
func (suite *KeeperTestSuite) TestVerifyCredential() {
	ctx := suite.ctx
	k := suite.keeper
	issuer := suite.addrs[0]
	holder := suite.addrs[1]
	verifier := suite.addrs[2]

	// Create holder identity and issue credential
	holderDID := "did:desh:verify123"
	identity := types.Identity{
		Did:        holderDID,
		Controller: holder.String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)

	credentialSubject := map[string]interface{}{
		"id":       holderDID,
		"verified": true,
		"level":    "high",
	}

	credentialID, err := k.IssueCredential(
		ctx,
		issuer,
		holderDID,
		[]string{"VerifiableCredential", "TestCredential"},
		credentialSubject,
	)
	suite.NoError(err)

	// Verify the credential
	isValid, err := k.VerifyCredential(ctx, credentialID, verifier)
	suite.NoError(err)
	suite.True(isValid)

	// Test verification of revoked credential
	err = k.RevokeCredential(ctx, issuer, credentialID, "Testing")
	suite.NoError(err)

	isValid, err = k.VerifyCredential(ctx, credentialID, verifier)
	suite.NoError(err)
	suite.False(isValid)
}

// TestPresentCredential tests selective disclosure
func (suite *KeeperTestSuite) TestPresentCredential() {
	ctx := suite.ctx
	k := suite.keeper
	issuer := suite.addrs[0]
	holder := suite.addrs[1]
	verifier := suite.addrs[2]

	// Create holder identity and issue credential
	holderDID := "did:desh:present123"
	identity := types.Identity{
		Did:        holderDID,
		Controller: holder.String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)

	credentialSubject := map[string]interface{}{
		"id":       holderDID,
		"name":     "John Doe",
		"age":      30,
		"email":    "john@example.com",
		"verified": true,
	}

	credentialID, err := k.IssueCredential(
		ctx,
		issuer,
		holderDID,
		[]string{"VerifiableCredential", "PersonalCredential"},
		credentialSubject,
	)
	suite.NoError(err)

	// Present credential with selective disclosure
	disclosureMap := map[string][]string{
		credentialID: {"name", "verified"}, // Only disclose name and verified status
	}

	presentationID := k.PresentCredential(
		ctx,
		[]string{credentialID},
		verifier.String(),
		disclosureMap,
	)

	suite.NotEmpty(presentationID)

	// Retrieve and verify presentation
	presentation, found := k.GetPresentation(ctx, presentationID)
	suite.True(found)
	suite.Equal(verifier.String(), presentation.VerifierDid)
	suite.Len(presentation.VerifiableCredentials, 1)

	// Verify only disclosed attributes are present
	presentedCred := presentation.VerifiableCredentials[0]
	suite.Contains(presentedCred.CredentialSubject, "name")
	suite.Contains(presentedCred.CredentialSubject, "verified")
	suite.NotContains(presentedCred.CredentialSubject, "age")
	suite.NotContains(presentedCred.CredentialSubject, "email")
}

// TestGetAllIdentities tests bulk identity retrieval
func (suite *KeeperTestSuite) TestGetAllIdentities() {
	ctx := suite.ctx
	k := suite.keeper

	// Create multiple identities
	identities := []types.Identity{
		{
			Did:        "did:desh:bulk1",
			Controller: suite.addrs[0].String(),
			Status:     types.IdentityStatus_ACTIVE,
			CreatedAt:  ctx.BlockTime(),
			UpdatedAt:  ctx.BlockTime(),
		},
		{
			Did:        "did:desh:bulk2",
			Controller: suite.addrs[1].String(),
			Status:     types.IdentityStatus_ACTIVE,
			CreatedAt:  ctx.BlockTime(),
			UpdatedAt:  ctx.BlockTime(),
		},
		{
			Did:        "did:desh:bulk3",
			Controller: suite.addrs[2].String(),
			Status:     types.IdentityStatus_SUSPENDED,
			CreatedAt:  ctx.BlockTime(),
			UpdatedAt:  ctx.BlockTime(),
		},
	}

	// Store all identities
	for _, identity := range identities {
		k.SetIdentity(ctx, identity)
	}

	// Retrieve all identities
	allIdentities := k.GetAllIdentities(ctx)
	suite.Len(allIdentities, 3)

	// Verify all identities are present
	dids := make(map[string]bool)
	for _, identity := range allIdentities {
		dids[identity.Did] = true
	}

	for _, originalIdentity := range identities {
		suite.True(dids[originalIdentity.Did], "Identity %s not found", originalIdentity.Did)
	}
}

// TestZeroKnowledgeProof tests ZK proof generation and verification
func (suite *KeeperTestSuite) TestZeroKnowledgeProof() {
	ctx := suite.ctx
	k := suite.keeper
	prover := suite.addrs[0]

	// Create identity
	did := "did:desh:zkproof123"
	identity := types.Identity{
		Did:        did,
		Controller: prover.String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)

	// Test age proof (proving age >= 18 without revealing actual age)
	actualAge := uint32(25)
	minAge := uint32(18)

	proofData, err := k.GenerateAgeProof(ctx, did, actualAge, minAge)
	suite.NoError(err)
	suite.NotEmpty(proofData)

	// Verify the proof
	isValid, err := k.VerifyAgeProof(ctx, did, proofData, minAge)
	suite.NoError(err)
	suite.True(isValid)

	// Test with age below threshold
	youngAge := uint32(16)
	proofDataYoung, err := k.GenerateAgeProof(ctx, did, youngAge, minAge)
	suite.NoError(err)

	isValidYoung, err := k.VerifyAgeProof(ctx, did, proofDataYoung, minAge)
	suite.NoError(err)
	suite.False(isValidYoung)
}

// TestBiometricTemplate tests biometric template storage and verification
func (suite *KeeperTestSuite) TestBiometricTemplate() {
	ctx := suite.ctx
	k := suite.keeper
	user := suite.addrs[0]

	// Create identity
	did := "did:desh:biometric123"
	identity := types.Identity{
		Did:        did,
		Controller: user.String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)

	// Store biometric template
	templateData := []byte("encrypted_biometric_template_data")
	templateType := "fingerprint"

	err := k.StoreBiometricTemplate(ctx, did, templateType, templateData)
	suite.NoError(err)

	// Verify template storage
	storedTemplate, found := k.GetBiometricTemplate(ctx, did, templateType)
	suite.True(found)
	suite.Equal(templateData, storedTemplate.TemplateData)
	suite.Equal(templateType, storedTemplate.TemplateType)

	// Test biometric verification
	candidateData := templateData // In real scenario, this would be processed biometric data
	isMatch, err := k.VerifyBiometric(ctx, did, templateType, candidateData)
	suite.NoError(err)
	suite.True(isMatch)

	// Test with non-matching data
	nonMatchingData := []byte("different_biometric_data")
	isMatch, err = k.VerifyBiometric(ctx, did, templateType, nonMatchingData)
	suite.NoError(err)
	suite.False(isMatch)
}

// TestConsentManagement tests consent mechanisms
func (suite *KeeperTestSuite) TestConsentManagement() {
	ctx := suite.ctx
	k := suite.keeper
	dataSubject := suite.addrs[0]
	dataProcessor := suite.addrs[1]

	// Create identity
	did := "did:desh:consent123"
	identity := types.Identity{
		Did:        did,
		Controller: dataSubject.String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)

	// Grant consent
	purposes := []string{"KYC verification", "Risk assessment"}
	dataTypes := []string{"name", "age", "address"}
	
	consentID, err := k.GrantConsent(
		ctx,
		did,
		dataProcessor.String(),
		purposes,
		dataTypes,
		ctx.BlockTime().Add(24*time.Hour), // Valid for 24 hours
	)
	suite.NoError(err)
	suite.NotEmpty(consentID)

	// Verify consent
	hasConsent := k.HasValidConsent(ctx, did, dataProcessor.String(), "KYC verification", "name")
	suite.True(hasConsent)

	// Test consent for non-granted purpose
	hasConsentOther := k.HasValidConsent(ctx, did, dataProcessor.String(), "Marketing", "name")
	suite.False(hasConsentOther)

	// Revoke consent
	err = k.RevokeConsent(ctx, did, consentID)
	suite.NoError(err)

	// Verify consent is revoked
	hasConsentAfterRevoke := k.HasValidConsent(ctx, did, dataProcessor.String(), "KYC verification", "name")
	suite.False(hasConsentAfterRevoke)
}

// TestCredentialExpiry tests credential expiration handling
func (suite *KeeperTestSuite) TestCredentialExpiry() {
	ctx := suite.ctx
	k := suite.keeper
	issuer := suite.addrs[0]
	holder := suite.addrs[1]

	// Create holder identity
	holderDID := "did:desh:expiry123"
	identity := types.Identity{
		Did:        holderDID,
		Controller: holder.String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)

	// Issue credential with expiry
	expiryTime := ctx.BlockTime().Add(time.Hour)
	credentialSubject := map[string]interface{}{
		"id":     holderDID,
		"data":   "test data",
		"expiry": expiryTime.Format(time.RFC3339),
	}

	credentialID, err := k.IssueCredential(
		ctx,
		issuer,
		holderDID,
		[]string{"VerifiableCredential", "ExpiringCredential"},
		credentialSubject,
	)
	suite.NoError(err)

	// Verify credential is valid before expiry
	isValid, err := k.VerifyCredential(ctx, credentialID, issuer)
	suite.NoError(err)
	suite.True(isValid)

	// Advance time past expiry
	expiredCtx := ctx.WithBlockTime(expiryTime.Add(time.Minute))

	// Verify credential is invalid after expiry
	isValidAfterExpiry, err := k.VerifyCredential(expiredCtx, credentialID, issuer)
	suite.NoError(err)
	suite.False(isValidAfterExpiry)
}