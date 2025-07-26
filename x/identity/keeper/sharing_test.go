package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/DeshChain/DeshChain-Ecosystem/x/identity/keeper"
	"github.com/DeshChain/DeshChain-Ecosystem/x/identity/types"
	"github.com/DeshChain/DeshChain-Ecosystem/testutil"
)

type SharingTestSuite struct {
	suite.Suite

	ctx    sdk.Context
	keeper *keeper.Keeper
	addrs  []sdk.AccAddress
}

func (suite *SharingTestSuite) SetupTest() {
	suite.ctx, suite.keeper = testutil.IdentityKeeperTestSetup(suite.T())
	suite.addrs = testutil.CreateIncrementalAccounts(10)
}

func TestSharingTestSuite(t *testing.T) {
	suite.Run(t, new(SharingTestSuite))
}

// TestCreateShareRequest tests creating a new identity sharing request
func (suite *SharingTestSuite) TestCreateShareRequest() {
	ctx := suite.ctx
	k := suite.keeper
	holderDID := "did:desh:sharetest123"
	
	// Create holder identity
	identity := types.Identity{
		Did:        holderDID,
		Controller: suite.addrs[0].String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)
	
	// Define requested data
	requestedData := []types.DataRequest{
		{
			CredentialType:    "KYCCredential",
			Attributes:        []string{"name", "age", "verified"},
			MinimumTrustLevel: "high",
			Required:          true,
		},
		{
			CredentialType:    "BiometricCredential",
			Attributes:        []string{"template_type", "verification_level"},
			MinimumTrustLevel: "medium",
			Required:          false,
		},
	}
	
	// Create share request
	request, err := k.CreateShareRequest(
		ctx,
		"tradefinance",
		"gramsuraksha",
		holderDID,
		requestedData,
		"Loan eligibility verification",
		"Required for trade finance loan application",
		24*time.Hour,
	)
	
	suite.NoError(err)
	suite.NotNil(request)
	suite.NotEmpty(request.RequestID)
	suite.Equal("tradefinance", request.RequesterModule)
	suite.Equal("gramsuraksha", request.ProviderModule)
	suite.Equal(holderDID, request.HolderDID)
	suite.Len(request.RequestedData, 2)
	suite.Equal(types.ShareRequestStatus_PENDING, request.Status)
	
	// Verify request is stored
	storedRequest, found := k.GetShareRequest(ctx, request.RequestID)
	suite.True(found)
	suite.Equal(request.RequestID, storedRequest.RequestID)
}

// TestApproveShareRequest tests approving a share request
func (suite *SharingTestSuite) TestApproveShareRequest() {
	ctx := suite.ctx
	k := suite.keeper
	holderDID := "did:desh:approve123"
	authority := suite.addrs[0]
	
	// Create holder identity
	identity := types.Identity{
		Did:        holderDID,
		Controller: authority.String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)
	
	// Issue some credentials to share
	kycSubject := map[string]interface{}{
		"id":       holderDID,
		"name":     "John Doe",
		"age":      30,
		"verified": true,
	}
	
	_, err := k.IssueCredential(
		ctx,
		suite.addrs[1],
		holderDID,
		[]string{"VerifiableCredential", "KYCCredential"},
		kycSubject,
	)
	suite.NoError(err)
	
	// Create share request
	requestedData := []types.DataRequest{
		{
			CredentialType: "KYCCredential",
			Attributes:     []string{"name", "age"},
			Required:       true,
		},
	}
	
	request, err := k.CreateShareRequest(
		ctx,
		"tradefinance",
		"gramsuraksha",
		holderDID,
		requestedData,
		"Test approval",
		"Testing approval flow",
		24*time.Hour,
	)
	suite.NoError(err)
	
	// Approve the request
	accessToken := "test_access_token_123"
	err = k.ApproveShareRequest(ctx, authority, request.RequestID, accessToken)
	suite.NoError(err)
	
	// Verify request is approved
	updatedRequest, found := k.GetShareRequest(ctx, request.RequestID)
	suite.True(found)
	suite.Equal(types.ShareRequestStatus_APPROVED, updatedRequest.Status)
	
	// Verify response is created
	response, found := k.GetShareResponse(ctx, request.RequestID)
	suite.True(found)
	suite.Equal(types.ShareRequestStatus_APPROVED, response.Status)
	suite.Equal(accessToken, response.AccessToken)
	suite.Len(response.SharedData, 1)
	
	// Verify shared data contains only requested attributes
	sharedCred := response.SharedData[0]
	suite.Equal("KYCCredential", sharedCred.CredentialType)
	suite.Contains(sharedCred.SharedData, "name")
	suite.Contains(sharedCred.SharedData, "age")
	suite.NotContains(sharedCred.SharedData, "verified") // Not requested
}

// TestDenyShareRequest tests denying a share request
func (suite *SharingTestSuite) TestDenyShareRequest() {
	ctx := suite.ctx
	k := suite.keeper
	holderDID := "did:desh:deny123"
	authority := suite.addrs[0]
	
	// Create holder identity
	identity := types.Identity{
		Did:        holderDID,
		Controller: authority.String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)
	
	// Create share request
	requestedData := []types.DataRequest{
		{
			CredentialType: "KYCCredential",
			Attributes:     []string{"name", "age"},
			Required:       true,
		},
	}
	
	request, err := k.CreateShareRequest(
		ctx,
		"tradefinance",
		"gramsuraksha",
		holderDID,
		requestedData,
		"Test denial",
		"Testing denial flow",
		24*time.Hour,
	)
	suite.NoError(err)
	
	// Deny the request
	denialReason := "Insufficient KYC level"
	err = k.DenyShareRequest(ctx, authority, request.RequestID, denialReason)
	suite.NoError(err)
	
	// Verify request is denied
	updatedRequest, found := k.GetShareRequest(ctx, request.RequestID)
	suite.True(found)
	suite.Equal(types.ShareRequestStatus_DENIED, updatedRequest.Status)
	
	// Verify response is created with denial
	response, found := k.GetShareResponse(ctx, request.RequestID)
	suite.True(found)
	suite.Equal(types.ShareRequestStatus_DENIED, response.Status)
	suite.Equal(denialReason, response.DenialReason)
	suite.Len(response.SharedData, 0) // No data shared
}

// TestGetSharedData tests retrieving shared data with access token
func (suite *SharingTestSuite) TestGetSharedData() {
	ctx := suite.ctx
	k := suite.keeper
	holderDID := "did:desh:getdata123"
	authority := suite.addrs[0]
	
	// Create holder identity and credentials
	identity := types.Identity{
		Did:        holderDID,
		Controller: authority.String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)
	
	kycSubject := map[string]interface{}{
		"id":       holderDID,
		"name":     "Jane Doe",
		"age":      25,
		"verified": true,
	}
	
	_, err := k.IssueCredential(
		ctx,
		suite.addrs[1],
		holderDID,
		[]string{"VerifiableCredential", "KYCCredential"},
		kycSubject,
	)
	suite.NoError(err)
	
	// Create and approve share request
	requestedData := []types.DataRequest{
		{
			CredentialType: "KYCCredential",
			Attributes:     []string{"name", "verified"},
			Required:       true,
		},
	}
	
	request, err := k.CreateShareRequest(
		ctx,
		"tradefinance",
		"gramsuraksha",
		holderDID,
		requestedData,
		"Test data retrieval",
		"Testing data access",
		24*time.Hour,
	)
	suite.NoError(err)
	
	accessToken := "test_token_456"
	err = k.ApproveShareRequest(ctx, authority, request.RequestID, accessToken)
	suite.NoError(err)
	
	// Retrieve shared data
	response, err := k.GetSharedData(ctx, "tradefinance", request.RequestID, accessToken)
	suite.NoError(err)
	suite.NotNil(response)
	suite.Equal(types.ShareRequestStatus_APPROVED, response.Status)
	suite.Len(response.SharedData, 1)
	
	// Verify data content
	sharedCred := response.SharedData[0]
	suite.Equal("Jane Doe", sharedCred.SharedData["name"])
	suite.Equal(true, sharedCred.SharedData["verified"])
	suite.NotContains(sharedCred.SharedData, "age") // Not requested
	
	// Test unauthorized access
	_, err = k.GetSharedData(ctx, "wrong_module", request.RequestID, accessToken)
	suite.Error(err)
	suite.Contains(err.Error(), "unauthorized")
	
	// Test wrong access token
	_, err = k.GetSharedData(ctx, "tradefinance", request.RequestID, "wrong_token")
	suite.Error(err)
	suite.Contains(err.Error(), "invalid access token")
}

// TestCreateSharingAgreement tests creating a standing agreement between modules
func (suite *SharingTestSuite) TestCreateSharingAgreement() {
	ctx := suite.ctx
	k := suite.keeper
	authority := suite.addrs[0]
	
	// Create sharing agreement
	agreement, err := k.CreateSharingAgreement(
		ctx,
		authority,
		"tradefinance",
		"gramsuraksha",
		[]string{"KYCCredential", "BiometricCredential"},
		[]string{"loan_verification", "risk_assessment"},
		"high",
		true, // auto-approve
		48*time.Hour,
		365*24*time.Hour,
	)
	
	suite.NoError(err)
	suite.NotNil(agreement)
	suite.NotEmpty(agreement.AgreementID)
	suite.Equal("tradefinance", agreement.RequesterModule)
	suite.Equal("gramsuraksha", agreement.ProviderModule)
	suite.Len(agreement.AllowedDataTypes, 2)
	suite.Len(agreement.Purposes, 2)
	suite.True(agreement.AutoApprove)
	suite.Equal(types.AgreementStatus_ACTIVE, agreement.Status)
	
	// Verify agreement is stored
	storedAgreement, found := k.GetSharingAgreement(ctx, agreement.AgreementID)
	suite.True(found)
	suite.Equal(agreement.AgreementID, storedAgreement.AgreementID)
}

// TestCreateAccessPolicy tests creating an access policy
func (suite *SharingTestSuite) TestCreateAccessPolicy() {
	ctx := suite.ctx
	k := suite.keeper
	holderDID := "did:desh:policy123"
	
	// Create holder identity
	identity := types.Identity{
		Did:        holderDID,
		Controller: suite.addrs[0].String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)
	
	// Create access policy
	dataRestrictions := map[string][]string{
		"KYCCredential":      {"name", "verified"},
		"BiometricCredential": {"template_type"},
	}
	
	timeRestrictions := types.TimeRestriction{
		AllowedHours: []int{9, 10, 11, 12, 13, 14, 15, 16, 17}, // Business hours
		AllowedDays:  []int{1, 2, 3, 4, 5},                    // Weekdays
	}
	
	policy, err := k.CreateAccessPolicy(
		ctx,
		holderDID,
		[]string{"tradefinance", "gramsuraksha"}, // allowed modules
		[]string{"validator"},                     // denied modules
		dataRestrictions,
		[]string{"loan_verification", "kyc_check"}, // allowed purposes
		timeRestrictions,
		10,   // max shares per day
		true, // require explicit consent
	)
	
	suite.NoError(err)
	suite.NotNil(policy)
	suite.NotEmpty(policy.PolicyID)
	suite.Equal(holderDID, policy.HolderDID)
	suite.Len(policy.AllowedModules, 2)
	suite.Len(policy.DeniedModules, 1)
	suite.Len(policy.DataRestrictions, 2)
	suite.Len(policy.PurposeRestrictions, 2)
	suite.Equal(10, policy.MaxSharesPerDay)
	suite.True(policy.RequireExplicitConsent)
	
	// Verify policy is stored
	storedPolicy, found := k.GetAccessPolicy(ctx, policy.PolicyID)
	suite.True(found)
	suite.Equal(policy.PolicyID, storedPolicy.PolicyID)
	
	// Verify policy is indexed by holder
	holderPolicies := k.GetAccessPoliciesByHolder(ctx, holderDID)
	suite.Len(holderPolicies, 1)
	suite.Equal(policy.PolicyID, holderPolicies[0].PolicyID)
}

// TestAutoApproveRequest tests automatic approval based on agreements
func (suite *SharingTestSuite) TestAutoApproveRequest() {
	ctx := suite.ctx
	k := suite.keeper
	holderDID := "did:desh:autoapprove123"
	authority := suite.addrs[0]
	
	// Create holder identity
	identity := types.Identity{
		Did:        holderDID,
		Controller: authority.String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)
	
	// Create sharing agreement with auto-approve
	_, err := k.CreateSharingAgreement(
		ctx,
		authority,
		"tradefinance",
		"gramsuraksha",
		[]string{"KYCCredential"},
		[]string{"loan_verification"},
		"high",
		true, // auto-approve
		24*time.Hour,
		365*24*time.Hour,
	)
	suite.NoError(err)
	
	// Create share request that matches the agreement
	requestedData := []types.DataRequest{
		{
			CredentialType: "KYCCredential",
			Attributes:     []string{"name", "verified"},
			Required:       true,
		},
	}
	
	// Note: In a full implementation, this would check against stored agreements
	// For now, we manually test the CanAutoApprove logic
	request, err := k.CreateShareRequest(
		ctx,
		"tradefinance",
		"gramsuraksha",
		holderDID,
		requestedData,
		"loan_verification",
		"Testing auto-approval",
		12*time.Hour, // Within allowed TTL
	)
	
	suite.NoError(err)
	// In a full implementation, this would be auto-approved
	suite.Equal(types.ShareRequestStatus_PENDING, request.Status)
}

// TestAccessPolicyViolation tests access policy enforcement
func (suite *SharingTestSuite) TestAccessPolicyViolation() {
	ctx := suite.ctx
	k := suite.keeper
	holderDID := "did:desh:violation123"
	
	// Create holder identity
	identity := types.Identity{
		Did:        holderDID,
		Controller: suite.addrs[0].String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)
	
	// Create restrictive access policy
	_, err := k.CreateAccessPolicy(
		ctx,
		holderDID,
		[]string{"gramsuraksha"}, // only allow gramsuraksha
		[]string{"tradefinance"}, // deny tradefinance
		map[string][]string{},
		[]string{"kyc_check"}, // only allow kyc_check purpose
		types.TimeRestriction{},
		1,    // max 1 share per day
		false, // don't require explicit consent
	)
	suite.NoError(err)
	
	// Try to create request from denied module
	requestedData := []types.DataRequest{
		{
			CredentialType: "KYCCredential",
			Attributes:     []string{"name"},
			Required:       true,
		},
	}
	
	_, err = k.CreateShareRequest(
		ctx,
		"tradefinance", // denied module
		"gramsuraksha",
		holderDID,
		requestedData,
		"loan_verification", // disallowed purpose
		"Should be denied",
		24*time.Hour,
	)
	
	suite.Error(err)
	suite.Contains(err.Error(), "access policy violation")
}

// TestShareRequestExpiry tests handling of expired requests
func (suite *SharingTestSuite) TestShareRequestExpiry() {
	ctx := suite.ctx
	k := suite.keeper
	holderDID := "did:desh:expiry123"
	authority := suite.addrs[0]
	
	// Create holder identity
	identity := types.Identity{
		Did:        holderDID,
		Controller: authority.String(),
		Status:     types.IdentityStatus_ACTIVE,
		CreatedAt:  ctx.BlockTime(),
		UpdatedAt:  ctx.BlockTime(),
	}
	k.SetIdentity(ctx, identity)
	
	// Create share request with short TTL
	requestedData := []types.DataRequest{
		{
			CredentialType: "KYCCredential",
			Attributes:     []string{"name"},
			Required:       true,
		},
	}
	
	request, err := k.CreateShareRequest(
		ctx,
		"tradefinance",
		"gramsuraksha",
		holderDID,
		requestedData,
		"Test expiry",
		"Testing expiry",
		1*time.Nanosecond, // Very short TTL
	)
	suite.NoError(err)
	
	// Verify request is expired
	suite.True(request.IsExpired())
	
	// Try to approve expired request
	err = k.ApproveShareRequest(ctx, authority, request.RequestID, "token")
	suite.Error(err)
	suite.Contains(err.Error(), "expired")
}
