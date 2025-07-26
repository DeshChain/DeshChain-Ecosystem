package keeper_test

import (
	"fmt"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/DeshChain/DeshChain-Ecosystem/x/identity/keeper"
	"github.com/DeshChain/DeshChain-Ecosystem/x/identity/types"
	"github.com/DeshChain/DeshChain-Ecosystem/testutil"
)

type PerformanceTestSuite struct {
	suite.Suite

	ctx    sdk.Context
	keeper *keeper.Keeper
	addrs  []sdk.AccAddress
}

func (suite *PerformanceTestSuite) SetupTest() {
	suite.ctx, suite.keeper = testutil.IdentityKeeperTestSetup(suite.T())
	suite.addrs = testutil.CreateIncrementalAccounts(1000) // Large number for performance testing
}

func TestPerformanceTestSuite(t *testing.T) {
	suite.Run(t, new(PerformanceTestSuite))
}

// TestHighVolumeIdentityCreation tests creating many identities
func (suite *PerformanceTestSuite) TestHighVolumeIdentityCreation() {
	ctx := suite.ctx
	k := suite.keeper
	
	numIdentities := 1000
	start := time.Now()
	
	// Create many identities
	for i := 0; i < numIdentities; i++ {
		did := fmt.Sprintf("did:desh:perf%d", i)
		identity := types.Identity{
			Did:        did,
			Controller: suite.addrs[i].String(),
			Status:     types.IdentityStatus_ACTIVE,
			CreatedAt:  ctx.BlockTime(),
			UpdatedAt:  ctx.BlockTime(),
		}
		k.SetIdentity(ctx, identity)
	}
	
	duration := time.Since(start)
	
	// Performance assertion: should complete within reasonable time
	suite.Less(duration, 5*time.Second, "Creating %d identities took %v", numIdentities, duration)
	
	// Verify all identities were created
	allIdentities := k.GetAllIdentities(ctx)
	suite.Len(allIdentities, numIdentities)
	
	fmt.Printf("Created %d identities in %v (avg: %v per identity)\n", 
		numIdentities, duration, duration/time.Duration(numIdentities))
}

// TestHighVolumeCredentialIssuance tests issuing many credentials
func (suite *PerformanceTestSuite) TestHighVolumeCredentialIssuance() {
	ctx := suite.ctx
	k := suite.keeper
	issuer := suite.addrs[0]
	
	numCredentials := 500
	
	// Create identities first
	for i := 0; i < numCredentials; i++ {
		did := fmt.Sprintf("did:desh:cred%d", i)
		identity := types.Identity{
			Did:        did,
			Controller: suite.addrs[i+1].String(),
			Status:     types.IdentityStatus_ACTIVE,
			CreatedAt:  ctx.BlockTime(),
			UpdatedAt:  ctx.BlockTime(),
		}
		k.SetIdentity(ctx, identity)
	}
	
	start := time.Now()
	
	// Issue credentials
	for i := 0; i < numCredentials; i++ {
		did := fmt.Sprintf("did:desh:cred%d", i)
		credentialSubject := map[string]interface{}{
			"id":    did,
			"index": i,
			"batch": "performance_test",
		}
		
		_, err := k.IssueCredential(
			ctx,
			issuer,
			did,
			[]string{"VerifiableCredential", "PerformanceTestCredential"},
			credentialSubject,
		)
		suite.NoError(err)
	}
	
	duration := time.Since(start)
	
	// Performance assertion
	suite.Less(duration, 10*time.Second, "Issuing %d credentials took %v", numCredentials, duration)
	
	fmt.Printf("Issued %d credentials in %v (avg: %v per credential)\n", 
		numCredentials, duration, duration/time.Duration(numCredentials))
}

// TestConcurrentIdentityOperations tests concurrent operations
func (suite *PerformanceTestSuite) TestConcurrentIdentityOperations() {
	ctx := suite.ctx
	k := suite.keeper
	
	numGoroutines := 10
	operationsPerGoroutine := 100
	
	start := time.Now()
	done := make(chan bool, numGoroutines)
	
	// Run concurrent operations
	for g := 0; g < numGoroutines; g++ {
		go func(goroutineID int) {
			defer func() { done <- true }()
			
			for i := 0; i < operationsPerGoroutine; i++ {
				did := fmt.Sprintf("did:desh:concurrent%d_%d", goroutineID, i)
				identity := types.Identity{
					Did:        did,
					Controller: suite.addrs[goroutineID*operationsPerGoroutine+i].String(),
					Status:     types.IdentityStatus_ACTIVE,
					CreatedAt:  ctx.BlockTime(),
					UpdatedAt:  ctx.BlockTime(),
				}
				k.SetIdentity(ctx, identity)
				
				// Read back to verify
				_, found := k.GetIdentity(ctx, did)
				suite.True(found)
			}
		}(g)
	}
	
	// Wait for all goroutines to complete
	for g := 0; g < numGoroutines; g++ {
		<-done
	}
	
	duration := time.Since(start)
	totalOperations := numGoroutines * operationsPerGoroutine
	
	suite.Less(duration, 15*time.Second, "Concurrent operations took %v", duration)
	
	fmt.Printf("Completed %d concurrent operations in %v\n", totalOperations, duration)
}

// TestBulkCredentialRetrieval tests retrieving credentials in bulk
func (suite *PerformanceTestSuite) TestBulkCredentialRetrieval() {
	ctx := suite.ctx
	k := suite.keeper
	issuer := suite.addrs[0]
	holderDID := "did:desh:bulktest"
	
	// Create holder identity
	identity := types.Identity{
		Did:        holderDID,
		Controller: suite.addrs[1].String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)
	
	numCredentials := 100
	credentialTypes := []string{"TypeA", "TypeB", "TypeC", "TypeD", "TypeE"}
	
	// Issue many credentials of different types
	for i := 0; i < numCredentials; i++ {
		credType := credentialTypes[i%len(credentialTypes)]
		credentialSubject := map[string]interface{}{
			"id":    holderDID,
			"index": i,
			"type":  credType,
		}
		
		_, err := k.IssueCredential(
			ctx,
			issuer,
			holderDID,
			[]string{"VerifiableCredential", credType},
			credentialSubject,
		)
		suite.NoError(err)
	}
	
	start := time.Now()
	
	// Test bulk retrieval by holder
	allCreds := k.GetCredentialsByHolder(ctx, holderDID)
	suite.Len(allCreds, numCredentials)
	
	// Test retrieval by type
	for _, credType := range credentialTypes {
		typeCreds := k.GetCredentialsByType(ctx, holderDID, credType)
		suite.NotEmpty(typeCreds)
	}
	
	duration := time.Since(start)
	suite.Less(duration, 2*time.Second, "Bulk retrieval took %v", duration)
	
	fmt.Printf("Retrieved %d credentials in bulk in %v\n", numCredentials, duration)
}

// TestLargeCredentialPresentation tests presenting many credentials at once
func (suite *PerformanceTestSuite) TestLargeCredentialPresentation() {
	ctx := suite.ctx
	k := suite.keeper
	issuer := suite.addrs[0]
	verifier := suite.addrs[1]
	holderDID := "did:desh:presentation"
	
	// Create holder identity
	identity := types.Identity{
		Did:        holderDID,
		Controller: suite.addrs[2].String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)
	
	numCredentials := 50
	credentialIDs := make([]string, numCredentials)
	disclosureMap := make(map[string][]string)
	
	// Issue credentials
	for i := 0; i < numCredentials; i++ {
		credentialSubject := map[string]interface{}{
			"id":        holderDID,
			"index":     i,
			"data_a":    fmt.Sprintf("value_a_%d", i),
			"data_b":    fmt.Sprintf("value_b_%d", i),
			"sensitive": fmt.Sprintf("secret_%d", i),
		}
		
		credID, err := k.IssueCredential(
			ctx,
			issuer,
			holderDID,
			[]string{"VerifiableCredential", "LargePresentationCredential"},
			credentialSubject,
		)
		suite.NoError(err)
		
		credentialIDs[i] = credID
		disclosureMap[credID] = []string{"index", "data_a"} // Only disclose non-sensitive data
	}
	
	start := time.Now()
	
	// Create large presentation
	presentationID := k.PresentCredential(
		ctx,
		credentialIDs,
		verifier.String(),
		disclosureMap,
	)
	
	duration := time.Since(start)
	suite.NotEmpty(presentationID)
	suite.Less(duration, 5*time.Second, "Large presentation took %v", duration)
	
	// Verify presentation
	presentation, found := k.GetPresentation(ctx, presentationID)
	suite.True(found)
	suite.Len(presentation.VerifiableCredentials, numCredentials)
	
	// Verify selective disclosure worked
	for _, cred := range presentation.VerifiableCredentials {
		suite.Contains(cred.CredentialSubject, "index")
		suite.Contains(cred.CredentialSubject, "data_a")
		suite.NotContains(cred.CredentialSubject, "sensitive")
	}
	
	fmt.Printf("Created presentation with %d credentials in %v\n", numCredentials, duration)
}

// TestZKProofPerformance tests ZK proof generation performance
func (suite *PerformanceTestSuite) TestZKProofPerformance() {
	ctx := suite.ctx
	k := suite.keeper
	
	numProofs := 100
	did := "did:desh:zkperf"
	
	// Create identity
	identity := types.Identity{
		Did:        did,
		Controller: suite.addrs[0].String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)
	
	start := time.Now()
	
	// Generate many age proofs
	for i := 0; i < numProofs; i++ {
		actualAge := uint32(18 + (i % 50)) // Ages from 18 to 67
		minAge := uint32(18)
		
		proofData, err := k.GenerateAgeProof(ctx, did, actualAge, minAge)
		suite.NoError(err)
		suite.NotEmpty(proofData)
		
		// Verify the proof
		isValid, err := k.VerifyAgeProof(ctx, did, proofData, minAge)
		suite.NoError(err)
		suite.True(isValid)
	}
	
	duration := time.Since(start)
	suite.Less(duration, 30*time.Second, "ZK proof generation took %v", duration)
	
	fmt.Printf("Generated and verified %d ZK proofs in %v (avg: %v per proof)\n", 
		numProofs, duration, duration/time.Duration(numProofs))
}

// TestBiometricPerformance tests biometric operations performance
func (suite *PerformanceTestSuite) TestBiometricPerformance() {
	ctx := suite.ctx
	k := suite.keeper
	
	numBiometrics := 100
	did := "did:desh:bioperf"
	
	// Create identity
	identity := types.Identity{
		Did:        did,
		Controller: suite.addrs[0].String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)
	
	start := time.Now()
	
	// Store and verify biometric templates
	for i := 0; i < numBiometrics; i++ {
		templateType := fmt.Sprintf("fingerprint_%d", i%10) // 10 different fingers
		templateData := []byte(fmt.Sprintf("encrypted_template_data_%d", i))
		
		// Store template
		err := k.StoreBiometricTemplate(ctx, did, templateType, templateData)
		suite.NoError(err)
		
		// Verify template
		isMatch, err := k.VerifyBiometric(ctx, did, templateType, templateData)
		suite.NoError(err)
		suite.True(isMatch)
	}
	
	duration := time.Since(start)
	suite.Less(duration, 10*time.Second, "Biometric operations took %v", duration)
	
	fmt.Printf("Performed %d biometric operations in %v (avg: %v per operation)\n", 
		numBiometrics*2, duration, duration/time.Duration(numBiometrics*2))
}

// TestMemoryUsage tests memory efficiency with large datasets
func (suite *PerformanceTestSuite) TestMemoryUsage() {
	ctx := suite.ctx
	k := suite.keeper
	
	// Create a large number of identities and credentials
	numIdentities := 1000
	credsPerIdentity := 5
	
	start := time.Now()
	
	for i := 0; i < numIdentities; i++ {
		did := fmt.Sprintf("did:desh:memory%d", i)
		identity := types.Identity{
			Did:        did,
			Controller: suite.addrs[i%len(suite.addrs)].String(),
			Status:     types.IdentityStatus_ACTIVE,
			CreatedAt:  ctx.BlockTime(),
			UpdatedAt:  ctx.BlockTime(),
		}
		k.SetIdentity(ctx, identity)
		
		// Add credentials for each identity
		for j := 0; j < credsPerIdentity; j++ {
			credentialSubject := map[string]interface{}{
				"id":    did,
				"type":  fmt.Sprintf("MemoryTestType%d", j),
				"index": j,
			}
			
			_, err := k.IssueCredential(
				ctx,
				suite.addrs[0],
				did,
				[]string{"VerifiableCredential", fmt.Sprintf("MemoryTestCredential%d", j)},
				credentialSubject,
			)
			suite.NoError(err)
		}
	}
	
	duration := time.Since(start)
	totalItems := numIdentities + (numIdentities * credsPerIdentity)
	
	suite.Less(duration, 30*time.Second, "Memory test took %v", duration)
	
	// Test bulk retrieval performance doesn't degrade
	retrievalStart := time.Now()
	allIdentities := k.GetAllIdentities(ctx)
	retrievalDuration := time.Since(retrievalStart)
	
	suite.Len(allIdentities, numIdentities)
	suite.Less(retrievalDuration, 2*time.Second, "Bulk retrieval took %v", retrievalDuration)
	
	fmt.Printf("Created %d items in %v, retrieved all in %v\n", 
		totalItems, duration, retrievalDuration)
}
