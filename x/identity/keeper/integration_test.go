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

type IntegrationTestSuite struct {
	suite.Suite

	ctx            sdk.Context
	identityKeeper *keeper.Keeper
	addrs          []sdk.AccAddress
}

func (suite *IntegrationTestSuite) SetupTest() {
	suite.ctx, suite.identityKeeper = testutil.IdentityKeeperTestSetup(suite.T())
	suite.addrs = testutil.CreateIncrementalAccounts(10)
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

// TestKYCWorkflow tests the complete KYC workflow
func (suite *IntegrationTestSuite) TestKYCWorkflow() {
	ctx := suite.ctx
	k := suite.identityKeeper
	user := suite.addrs[0]
	kycProvider := suite.addrs[1]

	// Step 1: Create user identity
	userDID := "did:desh:user12345"
	identity := types.Identity{
		Did:        userDID,
		Controller: user.String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)

	// Step 2: Basic KYC credential issuance
	basicKYCSubject := map[string]interface{}{
		"id":         userDID,
		"name":       "John Doe",
		"country":    "IN",
		"kyc_level":  "basic",
		"verified":   true,
		"issued_at":  ctx.BlockTime().Format(time.RFC3339),
	}

	basicKYCID, err := k.IssueCredential(
		ctx,
		kycProvider,
		userDID,
		[]string{"VerifiableCredential", "KYCCredential"},
		basicKYCSubject,
	)
	suite.NoError(err)

	// Step 3: Enhanced KYC upgrade
	enhancedKYCSubject := map[string]interface{}{
		"id":                userDID,
		"kyc_level":         "enhanced",
		"income_verified":   true,
		"address_verified":  true,
		"documents":         []string{"passport", "utility_bill"},
		"risk_score":        "low",
		"upgraded_at":       ctx.BlockTime().Add(time.Hour).Format(time.RFC3339),
	}

	enhancedKYCID, err := k.IssueCredential(
		ctx,
		kycProvider,
		userDID,
		[]string{"VerifiableCredential", "EnhancedKYCCredential"},
		enhancedKYCSubject,
	)
	suite.NoError(err)

	// Step 4: Verify KYC progression
	kycCreds := k.GetCredentialsByType(ctx, userDID, "KYCCredential")
	suite.Len(kycCreds, 1)
	suite.Equal(basicKYCID, kycCreds[0].Id)

	enhancedCreds := k.GetCredentialsByType(ctx, userDID, "EnhancedKYCCredential")
	suite.Len(enhancedCreds, 1)
	suite.Equal(enhancedKYCID, enhancedCreds[0].Id)

	// Step 5: Verify user can present different KYC levels
	basicVerifier := suite.addrs[2]
	enhancedVerifier := suite.addrs[3]

	// Present basic KYC
	basicDisclosure := map[string][]string{
		basicKYCID: {"name", "country", "verified"},
	}
	basicPresentationID := k.PresentCredential(ctx, []string{basicKYCID}, basicVerifier.String(), basicDisclosure)
	suite.NotEmpty(basicPresentationID)

	// Present enhanced KYC
	enhancedDisclosure := map[string][]string{
		enhancedKYCID: {"kyc_level", "risk_score", "income_verified"},
	}
	enhancedPresentationID := k.PresentCredential(ctx, []string{enhancedKYCID}, enhancedVerifier.String(), enhancedDisclosure)
	suite.NotEmpty(enhancedPresentationID)

	// Verify presentations exist
	basicPresentation, found := k.GetPresentation(ctx, basicPresentationID)
	suite.True(found)
	suite.Equal(basicVerifier.String(), basicPresentation.VerifierDid)

	enhancedPresentation, found := k.GetPresentation(ctx, enhancedPresentationID)
	suite.True(found)
	suite.Equal(enhancedVerifier.String(), enhancedPresentation.VerifierDid)
}

// TestMultiModuleCredentialFlow tests credentials across different modules
func (suite *IntegrationTestSuite) TestMultiModuleCredentialFlow() {
	ctx := suite.ctx
	k := suite.identityKeeper
	user := suite.addrs[0]
	
	// Different module authorities
	tradeFinanceAuth := suite.addrs[1]
	gramSurakshaAuth := suite.addrs[2]
	validatorAuth := suite.addrs[3]
	remittanceAuth := suite.addrs[4]

	// Create user identity
	userDID := "did:desh:multimodule123"
	identity := types.Identity{
		Did:        userDID,
		Controller: user.String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)

	// Issue credentials from different modules
	
	// 1. TradeFinance KYC Credential
	tradeKYCSubject := map[string]interface{}{
		"id":              userDID,
		"trader_id":       "TDR-001",
		"business_type":   "importer",
		"kyc_level":       "enhanced",
		"trade_license":   "TL123456",
		"risk_category":   "low",
	}
	
	tradeKYCID, err := k.IssueCredential(
		ctx,
		tradeFinanceAuth,
		userDID,
		[]string{"VerifiableCredential", "TradeFinanceKYCCredential"},
		tradeKYCSubject,
	)
	suite.NoError(err)

	// 2. GramSuraksha Participant Credential
	gramParticipantSubject := map[string]interface{}{
		"id":                userDID,
		"participant_id":    "GP-001",
		"village_code":      "560001",
		"age_verified":      true,
		"income_category":   "rural",
		"contribution_tier": "gold",
	}
	
	gramParticipantID, err := k.IssueCredential(
		ctx,
		gramSurakshaAuth,
		userDID,
		[]string{"VerifiableCredential", "GramSurakshaParticipantCredential"},
		gramParticipantSubject,
	)
	suite.NoError(err)

	// 3. Validator Credential
	validatorSubject := map[string]interface{}{
		"id":             userDID,
		"validator_rank": 5,
		"stake_amount":   "300000000000",
		"nft_bound":      true,
		"tier":           1,
	}
	
	validatorCredID, err := k.IssueCredential(
		ctx,
		validatorAuth,
		userDID,
		[]string{"VerifiableCredential", "ValidatorCredential"},
		validatorSubject,
	)
	suite.NoError(err)

	// 4. Remittance Sender Credential
	remittanceSenderSubject := map[string]interface{}{
		"id":                 userDID,
		"kyc_level":          "enhanced",
		"aml_verified":       true,
		"sanctions_checked":  true,
		"max_transfer_limit": "50000usd",
		"source_of_funds":    "business",
		"risk_level":         "low",
	}
	
	remittanceSenderID, err := k.IssueCredential(
		ctx,
		remittanceAuth,
		userDID,
		[]string{"VerifiableCredential", "RemittanceSenderCredential"},
		remittanceSenderSubject,
	)
	suite.NoError(err)

	// Verify all credentials are issued
	allCreds := k.GetCredentialsByHolder(ctx, userDID)
	suite.Len(allCreds, 4)

	// Test cross-module verification scenarios
	
	// Scenario 1: Validator applying for trade finance
	// Should be able to present validator credential + KYC for enhanced trust
	tradeVerifier := suite.addrs[5]
	crossModuleDisclosure := map[string][]string{
		tradeKYCID:      {"business_type", "kyc_level", "risk_category"},
		validatorCredID: {"validator_rank", "stake_amount"},
	}
	
	crossModulePresentationID := k.PresentCredential(
		ctx,
		[]string{tradeKYCID, validatorCredID},
		tradeVerifier.String(),
		crossModuleDisclosure,
	)
	suite.NotEmpty(crossModulePresentationID)

	// Scenario 2: Remittance sender with rural background verification
	remittanceVerifier := suite.addrs[6]
	ruralRemittanceDisclosure := map[string][]string{
		remittanceSenderID: {"kyc_level", "aml_verified", "max_transfer_limit"},
		gramParticipantID:  {"village_code", "income_category"},
	}
	
	ruralRemittancePresentationID := k.PresentCredential(
		ctx,
		[]string{remittanceSenderID, gramParticipantID},
		remittanceVerifier.String(),
		ruralRemittanceDisclosure,
	)
	suite.NotEmpty(ruralRemittancePresentationID)

	// Verify presentations
	crossModulePresentation, found := k.GetPresentation(ctx, crossModulePresentationID)
	suite.True(found)
	suite.Len(crossModulePresentation.VerifiableCredentials, 2)

	ruralRemittancePresentation, found := k.GetPresentation(ctx, ruralRemittancePresentationID)
	suite.True(found)
	suite.Len(ruralRemittancePresentation.VerifiableCredentials, 2)
}

// TestCredentialRevocationFlow tests credential revocation across modules
func (suite *IntegrationTestSuite) TestCredentialRevocationFlow() {
	ctx := suite.ctx
	k := suite.identityKeeper
	user := suite.addrs[0]
	issuer := suite.addrs[1]
	verifier := suite.addrs[2]

	// Create user identity
	userDID := "did:desh:revocation123"
	identity := types.Identity{
		Did:        userDID,
		Controller: user.String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)

	// Issue credential
	credentialSubject := map[string]interface{}{
		"id":       userDID,
		"status":   "active",
		"valid_until": ctx.BlockTime().Add(24*time.Hour).Format(time.RFC3339),
	}
	
	credentialID, err := k.IssueCredential(
		ctx,
		issuer,
		userDID,
		[]string{"VerifiableCredential", "TestCredential"},
		credentialSubject,
	)
	suite.NoError(err)

	// Create presentation before revocation
	disclosure := map[string][]string{
		credentialID: {"status"},
	}
	
	presentationID := k.PresentCredential(ctx, []string{credentialID}, verifier.String(), disclosure)
	suite.NotEmpty(presentationID)

	// Verify credential is valid
	isValid, err := k.VerifyCredential(ctx, credentialID, verifier)
	suite.NoError(err)
	suite.True(isValid)

	// Revoke credential
	err = k.RevokeCredential(ctx, issuer, credentialID, "Compliance violation detected")
	suite.NoError(err)

	// Verify credential is now invalid
	isValidAfterRevocation, err := k.VerifyCredential(ctx, credentialID, verifier)
	suite.NoError(err)
	suite.False(isValidAfterRevocation)

	// Verify presentation is also invalidated
	presentation, found := k.GetPresentation(ctx, presentationID)
	suite.True(found)
	
	// Check that the presentation contains revoked credential
	presentedCred := presentation.VerifiableCredentials[0]
	suite.Equal(types.CredentialStatus_REVOKED, presentedCred.Status)
}

// TestBiometricIntegrationFlow tests biometric integration across modules
func (suite *IntegrationTestSuite) TestBiometricIntegrationFlow() {
	ctx := suite.ctx
	k := suite.identityKeeper
	user := suite.addrs[0]
	biometricProvider := suite.addrs[1]
	moneyOrderVerifier := suite.addrs[2]

	// Create user identity
	userDID := "did:desh:biometric123"
	identity := types.Identity{
		Did:        userDID,
		Controller: user.String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)

	// Store biometric templates
	fingerprintTemplate := []byte("encrypted_fingerprint_template")
	faceTemplate := []byte("encrypted_face_template")

	err := k.StoreBiometricTemplate(ctx, userDID, "fingerprint", fingerprintTemplate)
	suite.NoError(err)

	err = k.StoreBiometricTemplate(ctx, userDID, "face", faceTemplate)
	suite.NoError(err)

	// Issue biometric credential
	biometricSubject := map[string]interface{}{
		"id":                userDID,
		"fingerprint_hash":  "fp_hash_123",
		"face_hash":         "face_hash_456",
		"templates_stored":  2,
		"verification_level": "high",
	}
	
	biometricCredID, err := k.IssueCredential(
		ctx,
		biometricProvider,
		userDID,
		[]string{"VerifiableCredential", "BiometricCredential"},
		biometricSubject,
	)
	suite.NoError(err)

	// Test high-value transaction scenario (MoneyOrder)
	// Create presentation for money order verification
	disclosure := map[string][]string{
		biometricCredID: {"verification_level", "templates_stored"},
	}
	
	moneyOrderPresentationID := k.PresentCredential(
		ctx,
		[]string{biometricCredID},
		moneyOrderVerifier.String(),
		disclosure,
	)
	suite.NotEmpty(moneyOrderPresentationID)

	// Verify biometric authentication for transaction
	candidateFingerprint := fingerprintTemplate // Simulating matching biometric
	isMatch, err := k.VerifyBiometric(ctx, userDID, "fingerprint", candidateFingerprint)
	suite.NoError(err)
	suite.True(isMatch)

	// Test non-matching biometric
	wrongFingerprint := []byte("wrong_fingerprint_data")
	isMatch, err = k.VerifyBiometric(ctx, userDID, "fingerprint", wrongFingerprint)
	suite.NoError(err)
	suite.False(isMatch)

	// Verify presentation exists and contains correct data
	presentation, found := k.GetPresentation(ctx, moneyOrderPresentationID)
	suite.True(found)
	suite.Equal(moneyOrderVerifier.String(), presentation.VerifierDid)
	
	presentedCred := presentation.VerifiableCredentials[0]
	suite.Equal("high", presentedCred.CredentialSubject["verification_level"])
	suite.Equal(float64(2), presentedCred.CredentialSubject["templates_stored"])
}

// TestConsentAcrossModules tests consent management across different modules
func (suite *IntegrationTestSuite) TestConsentAcrossModules() {
	ctx := suite.ctx
	k := suite.identityKeeper
	user := suite.addrs[0]
	
	// Different modules requesting data access
	tradeFinanceModule := suite.addrs[1]
	remittanceModule := suite.addrs[2]
	validatorModule := suite.addrs[3]

	// Create user identity
	userDID := "did:desh:consent123"
	identity := types.Identity{
		Did:        userDID,
		Controller: user.String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)

	// Grant consent for trade finance (KYC verification)
	tradeConsentID, err := k.GrantConsent(
		ctx,
		userDID,
		tradeFinanceModule.String(),
		[]string{"KYC verification", "Risk assessment"},
		[]string{"name", "business_type", "trade_license"},
		ctx.BlockTime().Add(30*24*time.Hour), // 30 days
	)
	suite.NoError(err)

	// Grant consent for remittance (AML compliance)
	remittanceConsentID, err := k.GrantConsent(
		ctx,
		userDID,
		remittanceModule.String(),
		[]string{"AML compliance", "Sanctions screening"},
		[]string{"name", "country", "source_of_funds"},
		ctx.BlockTime().Add(7*24*time.Hour), // 7 days
	)
	suite.NoError(err)

	// Grant consent for validator operations (stake verification)
	validatorConsentID, err := k.GrantConsent(
		ctx,
		userDID,
		validatorModule.String(),
		[]string{"Stake verification", "Performance monitoring"},
		[]string{"validator_rank", "stake_amount", "performance_metrics"},
		ctx.BlockTime().Add(365*24*time.Hour), // 1 year
	)
	suite.NoError(err)

	// Test consent verification for each module
	
	// Trade Finance access
	hasTradeConsent := k.HasValidConsent(ctx, userDID, tradeFinanceModule.String(), "KYC verification", "name")
	suite.True(hasTradeConsent)
	
	hasTradeConsentBusiness := k.HasValidConsent(ctx, userDID, tradeFinanceModule.String(), "Risk assessment", "business_type")
	suite.True(hasTradeConsentBusiness)

	// Remittance access
	hasRemittanceConsent := k.HasValidConsent(ctx, userDID, remittanceModule.String(), "AML compliance", "country")
	suite.True(hasRemittanceConsent)

	// Validator access
	hasValidatorConsent := k.HasValidConsent(ctx, userDID, validatorModule.String(), "Stake verification", "stake_amount")
	suite.True(hasValidatorConsent)

	// Test unauthorized access attempts
	hasUnauthorizedConsent := k.HasValidConsent(ctx, userDID, tradeFinanceModule.String(), "AML compliance", "country")
	suite.False(hasUnauthorizedConsent) // Trade finance shouldn't have AML consent

	hasUnauthorizedData := k.HasValidConsent(ctx, userDID, remittanceModule.String(), "AML compliance", "stake_amount")
	suite.False(hasUnauthorizedData) // Remittance shouldn't access validator data

	// Test consent revocation
	err = k.RevokeConsent(ctx, userDID, remittanceConsentID)
	suite.NoError(err)

	// Verify revoked consent is no longer valid
	hasRevokedConsent := k.HasValidConsent(ctx, userDID, remittanceModule.String(), "AML compliance", "country")
	suite.False(hasRevokedConsent)

	// Verify other consents are still valid
	hasTradeConsentAfterRevoke := k.HasValidConsent(ctx, userDID, tradeFinanceModule.String(), "KYC verification", "name")
	suite.True(hasTradeConsentAfterRevoke)

	hasValidatorConsentAfterRevoke := k.HasValidConsent(ctx, userDID, validatorModule.String(), "Stake verification", "stake_amount")
	suite.True(hasValidatorConsentAfterRevoke)

	// Verify we have the expected number of active consents
	allConsents := k.GetConsentsByHolder(ctx, userDID)
	activeConsents := 0
	for _, consent := range allConsents {
		if consent.Status == types.ConsentStatus_GRANTED {
			activeConsents++
		}
	}
	suite.Equal(2, activeConsents) // Trade and Validator consents should be active
}

// TestZKProofIntegration tests zero-knowledge proofs in practical scenarios
func (suite *IntegrationTestSuite) TestZKProofIntegration() {
	ctx := suite.ctx
	k := suite.identityKeeper
	user := suite.addrs[0]
	ageVerifier := suite.addrs[1]
	incomeVerifier := suite.addrs[2]

	// Create user identity
	userDID := "did:desh:zkproof123"
	identity := types.Identity{
		Did:        userDID,
		Controller: user.String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)

	// Test age verification for GramSuraksha (18+ requirement)
	actualAge := uint32(25)
	minAgeGramSuraksha := uint32(18)

	ageProofGram, err := k.GenerateAgeProof(ctx, userDID, actualAge, minAgeGramSuraksha)
	suite.NoError(err)

	isValidAgeGram, err := k.VerifyAgeProof(ctx, userDID, ageProofGram, minAgeGramSuraksha)
	suite.NoError(err)
	suite.True(isValidAgeGram)

	// Test age verification for senior benefits (60+ requirement)
	minAgeSenior := uint32(60)

	ageProofSenior, err := k.GenerateAgeProof(ctx, userDID, actualAge, minAgeSenior)
	suite.NoError(err)

	isValidAgeSenior, err := k.VerifyAgeProof(ctx, userDID, ageProofSenior, minAgeSenior)
	suite.NoError(err)
	suite.False(isValidAgeSenior) // 25 is not >= 60

	// Test income range proof for loan eligibility
	actualIncome := uint64(500000) // ₹5 lakhs
	minIncome := uint64(300000)    // ₹3 lakhs minimum
	maxIncome := uint64(1000000)   // ₹10 lakhs maximum

	incomeProof, err := k.GenerateIncomeRangeProof(ctx, userDID, actualIncome, minIncome, maxIncome)
	suite.NoError(err)

	isValidIncome, err := k.VerifyIncomeRangeProof(ctx, userDID, incomeProof, minIncome, maxIncome)
	suite.NoError(err)
	suite.True(isValidIncome)

	// Test income proof with out-of-range value
	lowIncome := uint64(200000) // ₹2 lakhs (below minimum)
	lowIncomeProof, err := k.GenerateIncomeRangeProof(ctx, userDID, lowIncome, minIncome, maxIncome)
	suite.NoError(err)

	isValidLowIncome, err := k.VerifyIncomeRangeProof(ctx, userDID, lowIncomeProof, minIncome, maxIncome)
	suite.NoError(err)
	suite.False(isValidLowIncome)

	// Create credentials with ZK proof data
	zkCredentialSubject := map[string]interface{}{
		"id":                 userDID,
		"age_proof_18plus":   string(ageProofGram),
		"income_proof_range": string(incomeProof),
		"verification_date":  ctx.BlockTime().Format(time.RFC3339),
	}

	zkCredID, err := k.IssueCredential(
		ctx,
		ageVerifier,
		userDID,
		[]string{"VerifiableCredential", "ZKProofCredential"},
		zkCredentialSubject,
	)
	suite.NoError(err)

	// Present ZK proofs to different verifiers
	ageDisclosure := map[string][]string{
		zkCredID: {"age_proof_18plus"},
	}

	agePresentationID := k.PresentCredential(ctx, []string{zkCredID}, ageVerifier.String(), ageDisclosure)
	suite.NotEmpty(agePresentationID)

	incomeDisclosure := map[string][]string{
		zkCredID: {"income_proof_range"},
	}

	incomePresentationID := k.PresentCredential(ctx, []string{zkCredID}, incomeVerifier.String(), incomeDisclosure)
	suite.NotEmpty(incomePresentationID)

	// Verify presentations exist and contain only disclosed attributes
	agePresentation, found := k.GetPresentation(ctx, agePresentationID)
	suite.True(found)
	
	ageCredential := agePresentation.VerifiableCredentials[0]
	suite.Contains(ageCredential.CredentialSubject, "age_proof_18plus")
	suite.NotContains(ageCredential.CredentialSubject, "income_proof_range")

	incomePresentation, found := k.GetPresentation(ctx, incomePresentationID)
	suite.True(found)
	
	incomeCredential := incomePresentation.VerifiableCredentials[0]
	suite.Contains(incomeCredential.CredentialSubject, "income_proof_range")
	suite.NotContains(incomeCredential.CredentialSubject, "age_proof_18plus")
}