package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"deshchain/app"
	"deshchain/x/charitabletrust/keeper"
	"deshchain/x/charitabletrust/types"
)

type KeeperTestSuite struct {
	suite.Suite

	app    *app.App
	ctx    sdk.Context
	keeper keeper.Keeper
	msgServer types.MsgServer
	queryClient types.QueryClient
}

func (suite *KeeperTestSuite) SetupTest() {
	isCheckTx := false
	suite.app = app.Setup(isCheckTx)
	suite.ctx = suite.app.BaseApp.NewContext(isCheckTx, tmproto.Header{
		Height: 1,
		Time:   time.Now(),
	})
	suite.keeper = suite.app.CharitableTrustKeeper
	suite.msgServer = keeper.NewMsgServerImpl(suite.keeper)

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.keeper)
	suite.queryClient = types.NewQueryClient(queryHelper)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

// Test Trust Fund Balance Operations
func (suite *KeeperTestSuite) TestTrustFundBalance() {
	ctx := suite.ctx
	k := suite.keeper

	// Test getting empty balance
	balance, found := k.GetTrustFundBalance(ctx)
	suite.Require().False(found)

	// Set balance
	testBalance := types.TrustFundBalance{
		TotalBalance:     sdk.NewCoin("unamo", sdk.NewInt(1000000000)),
		AllocatedAmount:  sdk.NewCoin("unamo", sdk.NewInt(200000000)),
		AvailableAmount:  sdk.NewCoin("unamo", sdk.NewInt(800000000)),
		TotalDistributed: sdk.NewCoin("unamo", sdk.NewInt(500000000)),
		LastUpdated:      ctx.BlockTime(),
	}
	k.SetTrustFundBalance(ctx, testBalance)

	// Get balance
	balance, found = k.GetTrustFundBalance(ctx)
	suite.Require().True(found)
	suite.Require().Equal(testBalance.TotalBalance, balance.TotalBalance)
	suite.Require().Equal(testBalance.AllocatedAmount, balance.AllocatedAmount)
	suite.Require().Equal(testBalance.AvailableAmount, balance.AvailableAmount)
	suite.Require().Equal(testBalance.TotalDistributed, balance.TotalDistributed)
}

// Test Governance Operations
func (suite *KeeperTestSuite) TestTrustGovernance() {
	ctx := suite.ctx
	k := suite.keeper

	// Create test trustees
	trustees := []types.Trustee{
		{
			Address:      "desh1trustee1test",
			Name:         "Test Trustee 1",
			Role:         "Chairman",
			Expertise:    "Social Impact",
			AppointedAt:  ctx.BlockTime(),
			TermEndDate:  ctx.BlockTime().Add(2 * 365 * 24 * time.Hour),
			Status:       "active",
			VotingPower:  100,
		},
		{
			Address:      "desh1trustee2test",
			Name:         "Test Trustee 2",
			Role:         "Secretary",
			Expertise:    "Finance",
			AppointedAt:  ctx.BlockTime(),
			TermEndDate:  ctx.BlockTime().Add(2 * 365 * 24 * time.Hour),
			Status:       "active",
			VotingPower:  100,
		},
	}

	// Set governance
	governance := types.TrustGovernance{
		Trustees:            trustees,
		Quorum:              2,
		ApprovalThreshold:   sdk.NewDecWithPrec(667, 3), // 66.7%
		AdvisoryCommittee:   []string{},
		TransparencyOfficer: "desh1transparency1test",
		NextElection:        ctx.BlockTime().Add(2 * 365 * 24 * time.Hour),
	}
	k.SetTrustGovernance(ctx, governance)

	// Get governance
	retrievedGov, found := k.GetTrustGovernance(ctx)
	suite.Require().True(found)
	suite.Require().Equal(len(governance.Trustees), len(retrievedGov.Trustees))
	suite.Require().Equal(governance.Quorum, retrievedGov.Quorum)

	// Test IsTrustee
	suite.Require().True(k.IsTrustee(ctx, "desh1trustee1test"))
	suite.Require().True(k.IsTrustee(ctx, "desh1trustee2test"))
	suite.Require().False(k.IsTrustee(ctx, "desh1nottrusteestest"))
}

// Test Charitable Allocation Operations
func (suite *KeeperTestSuite) TestCharitableAllocation() {
	ctx := suite.ctx
	k := suite.keeper

	// Create test allocation
	allocation := types.CharitableAllocation{
		Id:                     1,
		CharitableOrgWalletId:  1,
		OrganizationName:       "Test Charity",
		Amount:                 sdk.NewCoin("unamo", sdk.NewInt(100000000)),
		Purpose:                "Test donation",
		Category:               "education",
		ProposalId:             1,
		ApprovedBy:             []string{"desh1trustee1test", "desh1trustee2test"},
		AllocatedAt:            ctx.BlockTime(),
		ExpectedImpact:         "Help 100 students",
		Monitoring: &types.MonitoringRequirements{
			ReportingFrequency:     30,
			RequiredReports:        []string{"impact", "financial"},
			Kpis:                   []string{"students_helped", "funds_utilized"},
			MonitoringDuration:     180,
			SiteVisitsRequired:     false,
			FinancialAuditRequired: true,
		},
		Status: "active",
		Distribution: &types.Distribution{
			TxHash:        "",
			DistributedAt: time.Time{},
			DistributedBy: "",
		},
	}

	// Set allocation
	k.SetCharitableAllocation(ctx, allocation)

	// Get allocation
	retrieved, found := k.GetCharitableAllocation(ctx, 1)
	suite.Require().True(found)
	suite.Require().Equal(allocation.Id, retrieved.Id)
	suite.Require().Equal(allocation.OrganizationName, retrieved.OrganizationName)
	suite.Require().Equal(allocation.Amount, retrieved.Amount)
	suite.Require().Equal(allocation.Category, retrieved.Category)

	// Test GetAllCharitableAllocations
	allocations := k.GetAllCharitableAllocations(ctx)
	suite.Require().Len(allocations, 1)
	suite.Require().Equal(allocation.Id, allocations[0].Id)

	// Test increment allocation count
	nextID := k.IncrementAllocationCount(ctx)
	suite.Require().Equal(uint64(2), nextID)
}

// Test Allocation Proposal Operations
func (suite *KeeperTestSuite) TestAllocationProposal() {
	ctx := suite.ctx
	k := suite.keeper

	// Create test proposal
	proposal := types.AllocationProposal{
		Id:               1,
		Title:            "Q1 2025 Charity Distribution",
		Description:      "Quarterly distribution to charities",
		TotalAmount:      sdk.NewCoin("unamo", sdk.NewInt(1000000000)),
		Allocations:      []types.ProposedAllocation{},
		ProposedBy:       "desh1trustee1test",
		ProposedAt:       ctx.BlockTime(),
		VotingStartTime:  ctx.BlockTime(),
		VotingEndTime:    ctx.BlockTime().Add(7 * 24 * time.Hour),
		Status:           "pending",
		Votes:            []types.Vote{},
		Documents:        []string{},
		ExpectedOutcomes: "Impact 10,000 beneficiaries",
	}

	// Set proposal
	k.SetAllocationProposal(ctx, proposal)

	// Get proposal
	retrieved, found := k.GetAllocationProposal(ctx, 1)
	suite.Require().True(found)
	suite.Require().Equal(proposal.Id, retrieved.Id)
	suite.Require().Equal(proposal.Title, retrieved.Title)
	suite.Require().Equal(proposal.TotalAmount, retrieved.TotalAmount)

	// Test GetAllAllocationProposals
	proposals := k.GetAllAllocationProposals(ctx)
	suite.Require().Len(proposals, 1)
	suite.Require().Equal(proposal.Id, proposals[0].Id)

	// Test AddVoteToProposal
	vote := types.Vote{
		Voter:     "desh1trustee1test",
		VoteType:  "yes",
		VotedAt:   ctx.BlockTime(),
		Comments:  "Approved",
		VoteValue: 100,
	}
	err := k.AddVoteToProposal(ctx, 1, vote)
	suite.Require().NoError(err)

	// Verify vote was added
	retrieved, found = k.GetAllocationProposal(ctx, 1)
	suite.Require().True(found)
	suite.Require().Len(retrieved.Votes, 1)
	suite.Require().Equal(vote.Voter, retrieved.Votes[0].Voter)

	// Test duplicate vote
	err = k.AddVoteToProposal(ctx, 1, vote)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "already voted")

	// Test increment proposal count
	nextID := k.IncrementProposalCount(ctx)
	suite.Require().Equal(uint64(2), nextID)
}

// Test Impact Report Operations
func (suite *KeeperTestSuite) TestImpactReport() {
	ctx := suite.ctx
	k := suite.keeper

	// Create test impact report
	report := types.ImpactReport{
		Id:                   1,
		AllocationId:         1,
		OrganizationId:       1,
		ReportingPeriod:      "Q1 2025",
		SubmittedAt:          ctx.BlockTime(),
		BeneficiariesReached: 150,
		FundsUtilized:        sdk.NewCoin("unamo", sdk.NewInt(80000000)),
		ActivitiesConducted:  []string{"Distributed books", "Conducted classes"},
		OutcomesAchieved:     []string{"150 students enrolled", "95% attendance"},
		Challenges:           []string{"Remote area access"},
		SupportingDocuments:  []string{"https://ipfs.io/ipfs/Qm...report"},
		Verification:         nil,
		ImpactMetrics: map[string]string{
			"students_enrolled": "150",
			"attendance_rate":   "95%",
			"completion_rate":   "90%",
		},
	}

	// Set impact report
	k.SetImpactReport(ctx, report)

	// Get impact report
	retrieved, found := k.GetImpactReport(ctx, 1)
	suite.Require().True(found)
	suite.Require().Equal(report.Id, retrieved.Id)
	suite.Require().Equal(report.AllocationId, retrieved.AllocationId)
	suite.Require().Equal(report.BeneficiariesReached, retrieved.BeneficiariesReached)
	suite.Require().Equal(report.FundsUtilized, retrieved.FundsUtilized)

	// Test GetImpactReportsByAllocation
	reports := k.GetImpactReportsByAllocation(ctx, 1)
	suite.Require().Len(reports, 1)
	suite.Require().Equal(report.Id, reports[0].Id)

	// Test VerifyImpactReport
	verification := types.Verification{
		VerifiedBy:       "desh1verifier1test",
		VerifiedAt:       ctx.BlockTime(),
		IsVerified:       true,
		VerificationNote: "All claims verified",
		SiteVisitReport:  "",
		Score:            95,
	}
	err := k.VerifyImpactReport(ctx, 1, verification)
	suite.Require().NoError(err)

	// Verify the verification was added
	retrieved, found = k.GetImpactReport(ctx, 1)
	suite.Require().True(found)
	suite.Require().NotNil(retrieved.Verification)
	suite.Require().Equal(verification.IsVerified, retrieved.Verification.IsVerified)
	suite.Require().Equal(verification.Score, retrieved.Verification.Score)

	// Test increment report count
	nextID := k.IncrementReportCount(ctx)
	suite.Require().Equal(uint64(2), nextID)
}

// Test Fraud Alert Operations
func (suite *KeeperTestSuite) TestFraudAlert() {
	ctx := suite.ctx
	k := suite.keeper

	// Create test fraud alert
	alert := types.FraudAlert{
		Id:               1,
		AlertType:        "misuse_of_funds",
		Severity:         "high",
		AllocationId:     1,
		OrganizationId:   1,
		ReportedBy:       "desh1whistleblower1test",
		ReportedAt:       ctx.BlockTime(),
		Description:      "Funds used for personal expenses",
		Evidence:         []string{"Bank statements", "Receipts"},
		Status:           "pending",
		Investigation:    nil,
		Resolution:       nil,
		MonetaryImpact:   sdk.NewCoin("unamo", sdk.NewInt(10000000)),
		AffectedEntities: []string{"Test Charity"},
	}

	// Set fraud alert
	k.SetFraudAlert(ctx, alert)

	// Get fraud alert
	retrieved, found := k.GetFraudAlert(ctx, 1)
	suite.Require().True(found)
	suite.Require().Equal(alert.Id, retrieved.Id)
	suite.Require().Equal(alert.AlertType, retrieved.AlertType)
	suite.Require().Equal(alert.Severity, retrieved.Severity)
	suite.Require().Equal(alert.Description, retrieved.Description)

	// Test GetFraudAlertsByOrganization
	alerts := k.GetFraudAlertsByOrganization(ctx, 1)
	suite.Require().Len(alerts, 1)
	suite.Require().Equal(alert.Id, alerts[0].Id)

	// Test UpdateFraudAlertStatus
	err := k.UpdateFraudAlertStatus(ctx, 1, "investigating")
	suite.Require().NoError(err)

	retrieved, found = k.GetFraudAlert(ctx, 1)
	suite.Require().True(found)
	suite.Require().Equal("investigating", retrieved.Status)

	// Test AddInvestigationToAlert
	investigation := types.Investigation{
		InvestigatorId:   "desh1investigator1test",
		StartedAt:        ctx.BlockTime(),
		CompletedAt:      ctx.BlockTime().Add(7 * 24 * time.Hour),
		Findings:         []string{"Funds misused", "False receipts"},
		Recommendation:   "Blacklist organization",
		ActionsTaken:     []string{"Funds frozen", "Legal notice sent"},
		EvidenceGathered: []string{"Bank records", "Witness statements"},
	}
	err = k.AddInvestigationToAlert(ctx, 1, investigation)
	suite.Require().NoError(err)

	// Verify investigation was added
	retrieved, found = k.GetFraudAlert(ctx, 1)
	suite.Require().True(found)
	suite.Require().NotNil(retrieved.Investigation)
	suite.Require().Equal(investigation.InvestigatorId, retrieved.Investigation.InvestigatorId)
	suite.Require().Equal(investigation.Recommendation, retrieved.Investigation.Recommendation)

	// Test increment alert count
	nextID := k.IncrementAlertCount(ctx)
	suite.Require().Equal(uint64(2), nextID)
}

// Test parameter operations
func (suite *KeeperTestSuite) TestParams() {
	ctx := suite.ctx
	k := suite.keeper

	// Test default params
	params := k.GetParams(ctx)
	suite.Require().True(params.Enabled)
	suite.Require().Equal("unamo", params.MinAllocationAmount.Denom)

	// Update params
	newParams := types.Params{
		Enabled:                       false,
		MinAllocationAmount:           sdk.NewCoin("unamo", sdk.NewInt(200000000)),
		MaxMonthlyAllocationPerOrg:    sdk.NewCoin("unamo", sdk.NewInt(200000000000)),
		ProposalVotingPeriod:          14 * 24 * 60 * 60, // 14 days
		FraudInvestigationPeriod:      60,
		ImpactReportFrequency:         60,
		DistributionCategories:        []string{"education", "healthcare", "environment"},
	}
	k.SetParams(ctx, newParams)

	// Verify params updated
	updatedParams := k.GetParams(ctx)
	suite.Require().False(updatedParams.Enabled)
	suite.Require().Equal(newParams.MinAllocationAmount, updatedParams.MinAllocationAmount)
	suite.Require().Equal(newParams.ProposalVotingPeriod, updatedParams.ProposalVotingPeriod)
	suite.Require().Len(updatedParams.DistributionCategories, 3)
}

// Test validation functions
func (suite *KeeperTestSuite) TestValidations() {
	ctx := suite.ctx
	k := suite.keeper

	// Set up governance with trustees
	trustees := []types.Trustee{
		{
			Address:     "desh1trustee1test",
			Name:        "Test Trustee 1",
			Status:      "active",
			TermEndDate: ctx.BlockTime().Add(365 * 24 * time.Hour),
			VotingPower: 100,
		},
		{
			Address:     "desh1trustee2test",
			Name:        "Test Trustee 2",
			Status:      "active",
			TermEndDate: ctx.BlockTime().Add(365 * 24 * time.Hour),
			VotingPower: 100,
		},
		{
			Address:     "desh1trustee3test",
			Name:        "Test Trustee 3",
			Status:      "active",
			TermEndDate: ctx.BlockTime().Add(365 * 24 * time.Hour),
			VotingPower: 100,
		},
		{
			Address:     "desh1trustee4test",
			Name:        "Test Trustee 4",
			Status:      "active",
			TermEndDate: ctx.BlockTime().Add(365 * 24 * time.Hour),
			VotingPower: 100,
		},
	}
	governance := types.TrustGovernance{
		Trustees:          trustees,
		Quorum:            3,
		ApprovalThreshold: sdk.NewDecWithPrec(600, 3), // 60%
	}
	k.SetTrustGovernance(ctx, governance)

	// Test ValidateOrganization
	err := k.ValidateOrganization(ctx, 1)
	suite.Require().NoError(err) // Assuming organization validation passes for now

	// Test ValidateAllocationAmount
	minAmount := sdk.NewCoin("unamo", sdk.NewInt(100000000))
	validAmount := sdk.NewCoin("unamo", sdk.NewInt(200000000))
	invalidAmount := sdk.NewCoin("unamo", sdk.NewInt(50000000))

	k.SetParams(ctx, types.Params{
		Enabled:             true,
		MinAllocationAmount: minAmount,
	})

	err = k.ValidateAllocationAmount(ctx, validAmount)
	suite.Require().NoError(err)

	err = k.ValidateAllocationAmount(ctx, invalidAmount)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "below minimum")

	// Test CheckProposalQuorum
	proposal := types.AllocationProposal{
		Id:     1,
		Status: "pending",
		Votes: []types.Vote{
			{Voter: "desh1trustee1test", VoteType: "yes", VoteValue: 100},
			{Voter: "desh1trustee2test", VoteType: "yes", VoteValue: 100},
			{Voter: "desh1trustee3test", VoteType: "no", VoteValue: 100},
		},
	}
	k.SetAllocationProposal(ctx, proposal)

	hasQuorum, approved := k.CheckProposalQuorum(ctx, 1)
	suite.Require().True(hasQuorum)  // 3 votes >= quorum of 3
	suite.Require().True(approved)   // 200 yes vs 100 no = 66.7% > 60% threshold
}

// Test metric tracking
func (suite *KeeperTestSuite) TestMetrics() {
	ctx := suite.ctx
	k := suite.keeper

	// Test GetTotalImpactReports
	report1 := types.ImpactReport{Id: 1}
	report2 := types.ImpactReport{Id: 2}
	k.SetImpactReport(ctx, report1)
	k.SetImpactReport(ctx, report2)

	total := k.GetTotalImpactReports(ctx)
	suite.Require().Equal(2, total)

	// Test GetVerifiedImpactReports
	verifiedReport := types.ImpactReport{
		Id: 3,
		Verification: &types.Verification{
			IsVerified: true,
		},
	}
	k.SetImpactReport(ctx, verifiedReport)

	verifiedCount := k.GetVerifiedImpactReports(ctx)
	suite.Require().Equal(1, verifiedCount)
}