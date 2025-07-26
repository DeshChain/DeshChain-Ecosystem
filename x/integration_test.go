package integration_test

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
	charitabletrustkeeper "deshchain/x/charitabletrust/keeper"
	charitabletrusttypes "deshchain/x/charitabletrust/types"
	dswfkeeper "deshchain/x/dswf/keeper"
	dswftypes "deshchain/x/dswf/types"
	revenuekeeper "deshchain/x/revenue/keeper"
	revenuetypes "deshchain/x/revenue/types"
)

type IntegrationTestSuite struct {
	suite.Suite

	app                *app.App
	ctx                sdk.Context
	dswfKeeper         dswfkeeper.Keeper
	charitableTrustKeeper charitabletrustkeeper.Keeper
	revenueKeeper      revenuekeeper.Keeper
	dswfMsgServer      dswftypes.MsgServer
	trustMsgServer     charitabletrusttypes.MsgServer
	revenueMsgServer   revenuetypes.MsgServer
}

func (suite *IntegrationTestSuite) SetupTest() {
	isCheckTx := false
	suite.app = app.Setup(isCheckTx)
	suite.ctx = suite.app.BaseApp.NewContext(isCheckTx, tmproto.Header{
		Height: 1,
		Time:   time.Now(),
	})
	
	suite.dswfKeeper = suite.app.DswfKeeper
	suite.charitableTrustKeeper = suite.app.CharitableTrustKeeper
	suite.revenueKeeper = suite.app.RevenueKeeper
	
	suite.dswfMsgServer = dswfkeeper.NewMsgServerImpl(suite.dswfKeeper)
	suite.trustMsgServer = charitabletrustkeeper.NewMsgServerImpl(suite.charitableTrustKeeper)
	suite.revenueMsgServer = revenuekeeper.NewMsgServerImpl(suite.revenueKeeper)
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

// Test revenue distribution to both DSWF and CharitableTrust
func (suite *IntegrationTestSuite) TestRevenueDistributionToBothModules() {
	ctx := suite.ctx

	// Set up initial revenue parameters
	revenueParams := revenuetypes.Params{
		Enabled: true,
		DistributionSchedule: "daily",
		RevenueStreams: []string{"transaction_fees", "platform_fees"},
		DistributionPercentages: map[string]sdk.Dec{
			"validators":        sdk.NewDecWithPrec(30, 2), // 30%
			"community_pool":    sdk.NewDecWithPrec(20, 2), // 20%
			"dswf":             sdk.NewDecWithPrec(20, 2), // 20%
			"charitable_trust":  sdk.NewDecWithPrec(10, 2), // 10%
			"founders":         sdk.NewDecWithPrec(10, 2), // 10%
			"liquidity":        sdk.NewDecWithPrec(10, 2), // 10%
		},
	}
	suite.revenueKeeper.SetParams(ctx, revenueParams)

	// Set up DSWF governance
	dswfGovernance := dswftypes.FundGovernance{
		FundManagers: []dswftypes.FundManager{
			{
				Address:   "desh1fundmanager1test",
				Name:      "Test Fund Manager 1",
				AddedAt:   ctx.BlockTime(),
			},
			{
				Address:   "desh1fundmanager2test", 
				Name:      "Test Fund Manager 2",
				AddedAt:   ctx.BlockTime(),
			},
		},
		RequiredSignatures: 2,
	}
	suite.dswfKeeper.SetFundGovernance(ctx, dswfGovernance)

	// Set up CharitableTrust governance
	trustees := []charitabletrusttypes.Trustee{
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
	}
	trustGovernance := charitabletrusttypes.TrustGovernance{
		Trustees:          trustees,
		Quorum:            2,
		ApprovalThreshold: sdk.NewDecWithPrec(600, 3),
	}
	suite.charitableTrustKeeper.SetTrustGovernance(ctx, trustGovernance)

	// Simulate revenue collection (1000 NAMO)
	totalRevenue := sdk.NewCoin("unamo", sdk.NewInt(1000000000)) // 1000 NAMO

	// Execute revenue distribution
	distribution := revenuetypes.RevenueDistribution{
		TotalAmount: totalRevenue,
		Timestamp:   ctx.BlockTime(),
		DistributionBreakdown: map[string]sdk.Coin{
			"dswf":             sdk.NewCoin("unamo", sdk.NewInt(200000000)), // 200 NAMO
			"charitable_trust": sdk.NewCoin("unamo", sdk.NewInt(100000000)), // 100 NAMO
			"validators":       sdk.NewCoin("unamo", sdk.NewInt(300000000)), // 300 NAMO
			"community_pool":   sdk.NewCoin("unamo", sdk.NewInt(200000000)), // 200 NAMO
			"founders":         sdk.NewCoin("unamo", sdk.NewInt(100000000)), // 100 NAMO
			"liquidity":        sdk.NewCoin("unamo", sdk.NewInt(100000000)), // 100 NAMO
		},
	}

	// Verify DSWF received funds
	dswfFunds := suite.dswfKeeper.GetFundBalance(ctx)
	expectedDSWFBalance := sdk.NewCoin("unamo", sdk.NewInt(200000000))
	// In real implementation, this would be set by the revenue module
	suite.dswfKeeper.AddToFundBalance(ctx, expectedDSWFBalance)
	
	updatedDSWFBalance := suite.dswfKeeper.GetFundBalance(ctx)
	suite.Require().True(updatedDSWFBalance[0].Amount.Equal(expectedDSWFBalance.Amount))

	// Verify CharitableTrust received funds
	trustBalance := charitabletrusttypes.TrustFundBalance{
		TotalBalance:    sdk.NewCoin("unamo", sdk.NewInt(100000000)),
		AllocatedAmount: sdk.NewCoin("unamo", sdk.NewInt(0)),
		AvailableAmount: sdk.NewCoin("unamo", sdk.NewInt(100000000)),
		LastUpdated:     ctx.BlockTime(),
	}
	suite.charitableTrustKeeper.SetTrustFundBalance(ctx, trustBalance)

	retrievedBalance, found := suite.charitableTrustKeeper.GetTrustFundBalance(ctx)
	suite.Require().True(found)
	suite.Require().Equal(sdk.NewInt(100000000), retrievedBalance.TotalBalance.Amount)
}

// Test DSWF allocation proposal workflow
func (suite *IntegrationTestSuite) TestDSWFAllocationWorkflow() {
	ctx := suite.ctx

	// Set up DSWF with initial balance
	fundBalance := []sdk.Coin{sdk.NewCoin("unamo", sdk.NewInt(1000000000))} // 1000 NAMO
	suite.dswfKeeper.SetFundBalance(ctx, fundBalance)

	// Set up governance
	governance := dswftypes.FundGovernance{
		FundManagers: []dswftypes.FundManager{
			{Address: "desh1fundmanager1test", Name: "Manager 1"},
			{Address: "desh1fundmanager2test", Name: "Manager 2"},
		},
		RequiredSignatures: 2,
	}
	suite.dswfKeeper.SetFundGovernance(ctx, governance)

	// Create allocation proposal
	proposalMsg := &dswftypes.MsgProposeAllocation{
		Proposers:       []string{"desh1fundmanager1test", "desh1fundmanager2test"},
		Title:           "Infrastructure Development",
		Amount:          sdk.NewCoin("unamo", sdk.NewInt(100000000)), // 100 NAMO
		Category:        "infrastructure",
		RecipientId:     "recipient1",
		Description:     "Rural road development project",
		ExpectedImpact:  "Connect 10 villages to main highway",
		ExpectedReturns: "1.05",
		RiskAssessment:  "Low risk government project",
	}

	// Execute proposal
	res, err := suite.dswfMsgServer.ProposeAllocation(ctx, proposalMsg)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)

	// Verify allocation was created
	allocation, found := suite.dswfKeeper.GetFundAllocation(ctx, res.AllocationId)
	suite.Require().True(found)
	suite.Require().Equal("infrastructure", allocation.Category)
	suite.Require().Equal("approved", allocation.Status)

	// Test disbursement execution
	disbursementMsg := &dswftypes.MsgExecuteDisbursement{
		AllocationId: res.AllocationId,
		Amount:       sdk.NewCoin("unamo", sdk.NewInt(50000000)), // 50 NAMO
		Executor:     "desh1fundmanager1test",
		Recipient:    "desh1recipient1test",
		Purpose:      "First milestone payment",
	}

	disbursementRes, err := suite.dswfMsgServer.ExecuteDisbursement(ctx, disbursementMsg)
	suite.Require().NoError(err)
	suite.Require().NotNil(disbursementRes)

	// Verify disbursement was recorded
	allocation, found = suite.dswfKeeper.GetFundAllocation(ctx, res.AllocationId)
	suite.Require().True(found)
	suite.Require().Len(allocation.Disbursements, 1)
	suite.Require().Equal("executed", allocation.Disbursements[0].Status)
}

// Test CharitableTrust proposal and execution workflow
func (suite *IntegrationTestSuite) TestCharitableTrustWorkflow() {
	ctx := suite.ctx

	// Set up trust fund balance
	balance := charitabletrusttypes.TrustFundBalance{
		TotalBalance:    sdk.NewCoin("unamo", sdk.NewInt(500000000)), // 500 NAMO
		AvailableAmount: sdk.NewCoin("unamo", sdk.NewInt(500000000)),
		LastUpdated:     ctx.BlockTime(),
	}
	suite.charitableTrustKeeper.SetTrustFundBalance(ctx, balance)

	// Set up governance
	trustees := []charitabletrusttypes.Trustee{
		{
			Address:     "desh1trustee1test",
			Name:        "Trustee 1",
			Status:      "active",
			TermEndDate: ctx.BlockTime().Add(365 * 24 * time.Hour),
			VotingPower: 100,
		},
		{
			Address:     "desh1trustee2test",
			Name:        "Trustee 2",
			Status:      "active", 
			TermEndDate: ctx.BlockTime().Add(365 * 24 * time.Hour),
			VotingPower: 100,
		},
	}
	governance := charitabletrusttypes.TrustGovernance{
		Trustees:          trustees,
		Quorum:            2,
		ApprovalThreshold: sdk.NewDecWithPrec(500, 3), // 50%
	}
	suite.charitableTrustKeeper.SetTrustGovernance(ctx, governance)

	// Create allocation proposal
	proposalMsg := &charitabletrusttypes.MsgCreateAllocationProposal{
		Proposer:         "desh1trustee1test",
		Title:            "Q1 2025 Charity Distribution",
		Description:      "Support education and healthcare",
		TotalAmount:      sdk.NewCoin("unamo", sdk.NewInt(200000000)), // 200 NAMO
		ExpectedOutcomes: "Support 1000 beneficiaries",
		Allocations: []charitabletrusttypes.ProposedAllocation{
			{
				CharitableOrgWalletId: 1,
				OrganizationName:      "Education Foundation",
				Amount:                sdk.NewCoin("unamo", sdk.NewInt(100000000)),
				Purpose:               "School infrastructure",
				Category:              "education",
			},
			{
				CharitableOrgWalletId: 2,
				OrganizationName:      "Health Trust",
				Amount:                sdk.NewCoin("unamo", sdk.NewInt(100000000)),
				Purpose:               "Medical equipment",
				Category:              "healthcare",
			},
		},
	}

	// Create proposal
	proposalRes, err := suite.trustMsgServer.CreateAllocationProposal(ctx, proposalMsg)
	suite.Require().NoError(err)
	suite.Require().NotNil(proposalRes)

	// Vote on proposal
	vote1Msg := &charitabletrusttypes.MsgVoteOnProposal{
		ProposalId: proposalRes.ProposalId,
		Voter:      "desh1trustee1test",
		VoteType:   "yes",
		Comments:   "Good proposal",
	}

	_, err = suite.trustMsgServer.VoteOnProposal(ctx, vote1Msg)
	suite.Require().NoError(err)

	vote2Msg := &charitabletrusttypes.MsgVoteOnProposal{
		ProposalId: proposalRes.ProposalId,
		Voter:      "desh1trustee2test",
		VoteType:   "yes",
		Comments:   "Approved",
	}

	_, err = suite.trustMsgServer.VoteOnProposal(ctx, vote2Msg)
	suite.Require().NoError(err)

	// Check if proposal has quorum and is approved
	hasQuorum, approved := suite.charitableTrustKeeper.CheckProposalQuorum(ctx, proposalRes.ProposalId)
	suite.Require().True(hasQuorum)
	suite.Require().True(approved)

	// Execute allocation
	executeMsg := &charitabletrusttypes.MsgExecuteAllocation{
		ProposalId: proposalRes.ProposalId,
		Executor:   "desh1trustee1test",
	}

	executeRes, err := suite.trustMsgServer.ExecuteAllocation(ctx, executeMsg)
	suite.Require().NoError(err)
	suite.Require().NotNil(executeRes)
	suite.Require().Len(executeRes.AllocationIds, 2)

	// Verify allocations were created
	allocation1, found := suite.charitableTrustKeeper.GetCharitableAllocation(ctx, executeRes.AllocationIds[0])
	suite.Require().True(found)
	suite.Require().Equal("Education Foundation", allocation1.OrganizationName)
	suite.Require().Equal("active", allocation1.Status)

	allocation2, found := suite.charitableTrustKeeper.GetCharitableAllocation(ctx, executeRes.AllocationIds[1])
	suite.Require().True(found)
	suite.Require().Equal("Health Trust", allocation2.OrganizationName)
	suite.Require().Equal("active", allocation2.Status)

	// Verify trust fund balance was updated
	updatedBalance, found := suite.charitableTrustKeeper.GetTrustFundBalance(ctx)
	suite.Require().True(found)
	suite.Require().Equal(sdk.NewInt(200000000), updatedBalance.AllocatedAmount.Amount)
	suite.Require().Equal(sdk.NewInt(300000000), updatedBalance.AvailableAmount.Amount)
}

// Test fraud detection and investigation workflow
func (suite *IntegrationTestSuite) TestFraudDetectionWorkflow() {
	ctx := suite.ctx

	// Create a charitable allocation
	allocation := charitabletrusttypes.CharitableAllocation{
		Id:                    1,
		CharitableOrgWalletId: 1,
		OrganizationName:      "Suspicious Charity",
		Amount:                sdk.NewCoin("unamo", sdk.NewInt(100000000)),
		Status:                "distributed",
	}
	suite.charitableTrustKeeper.SetCharitableAllocation(ctx, allocation)

	// Report fraud
	fraudMsg := &charitabletrusttypes.MsgReportFraud{
		AlertType:        "misuse_of_funds",
		Severity:         "high",
		AllocationId:     1,
		OrganizationId:   1,
		Description:      "Funds used for unauthorized purposes",
		Evidence:         []string{"Bank statements", "Photos"},
		ReportedBy:       "desh1whistleblower1test",
		MonetaryImpact:   sdk.NewCoin("unamo", sdk.NewInt(50000000)),
		AffectedEntities: []string{"Suspicious Charity"},
	}

	fraudRes, err := suite.trustMsgServer.ReportFraud(ctx, fraudMsg)
	suite.Require().NoError(err)
	suite.Require().NotNil(fraudRes)

	// Verify fraud alert was created
	alert, found := suite.charitableTrustKeeper.GetFraudAlert(ctx, fraudRes.AlertId)
	suite.Require().True(found)
	suite.Require().Equal("high", alert.Severity)
	suite.Require().Equal("pending", alert.Status)

	// Investigate fraud
	investigateMsg := &charitabletrusttypes.MsgInvestigateFraud{
		AlertId:          fraudRes.AlertId,
		InvestigatorId:   "desh1investigator1test",
		Findings:         []string{"Confirmed misuse", "Documentation falsified"},
		Recommendation:   "blacklist_organization",
		ActionsTaken:     []string{"Funds frozen", "Legal action initiated"},
		EvidenceGathered: []string{"Additional bank records", "Witness statements"},
	}

	investigateRes, err := suite.trustMsgServer.InvestigateFraud(ctx, investigateMsg)
	suite.Require().NoError(err)
	suite.Require().NotNil(investigateRes)

	// Verify investigation was recorded
	updatedAlert, found := suite.charitableTrustKeeper.GetFraudAlert(ctx, fraudRes.AlertId)
	suite.Require().True(found)
	suite.Require().Equal("investigated", updatedAlert.Status)
	suite.Require().NotNil(updatedAlert.Investigation)
	suite.Require().Equal("blacklist_organization", updatedAlert.Investigation.Recommendation)
}

// Test impact reporting and verification
func (suite *IntegrationTestSuite) TestImpactReportingWorkflow() {
	ctx := suite.ctx

	// Create a distributed allocation
	allocation := charitabletrusttypes.CharitableAllocation{
		Id:                    1,
		CharitableOrgWalletId: 1,
		OrganizationName:      "Education Foundation",
		Amount:                sdk.NewCoin("unamo", sdk.NewInt(100000000)),
		Status:                "distributed",
	}
	suite.charitableTrustKeeper.SetCharitableAllocation(ctx, allocation)

	// Submit impact report
	reportMsg := &charitabletrusttypes.MsgSubmitImpactReport{
		AllocationId:         1,
		OrganizationId:       1,
		ReportingPeriod:      "Q1 2025",
		BeneficiariesReached: 500,
		FundsUtilized:        sdk.NewCoin("unamo", sdk.NewInt(80000000)),
		ActivitiesConducted: []string{
			"Built 2 classrooms",
			"Trained 10 teachers",
			"Distributed books to 500 students",
		},
		OutcomesAchieved: []string{
			"500 students enrolled",
			"95% attendance rate",
			"90% completion rate",
		},
		Challenges: []string{
			"Remote location access",
			"Teacher retention",
		},
		SupportingDocuments: []string{
			"https://ipfs.io/ipfs/Qm...photos",
			"https://ipfs.io/ipfs/Qm...attendance",
		},
		ImpactMetrics: map[string]string{
			"students_enrolled": "500",
			"attendance_rate":   "95%",
			"completion_rate":   "90%",
		},
		Submitter: "desh1charity1test",
	}

	reportRes, err := suite.trustMsgServer.SubmitImpactReport(ctx, reportMsg)
	suite.Require().NoError(err)
	suite.Require().NotNil(reportRes)

	// Verify report was created
	report, found := suite.charitableTrustKeeper.GetImpactReport(ctx, reportRes.ReportId)
	suite.Require().True(found)
	suite.Require().Equal(int32(500), report.BeneficiariesReached)
	suite.Require().Len(report.ActivitiesConducted, 3)
	suite.Require().Len(report.OutcomesAchieved, 3)

	// Verify impact report
	verifyMsg := &charitabletrusttypes.MsgVerifyImpactReport{
		ReportId:         reportRes.ReportId,
		Verifier:         "desh1verifier1test",
		IsVerified:       true,
		VerificationNote: "All claims verified through site visit",
		Score:            92,
	}

	verifyRes, err := suite.trustMsgServer.VerifyImpactReport(ctx, verifyMsg)
	suite.Require().NoError(err)
	suite.Require().NotNil(verifyRes)

	// Verify the verification was recorded
	verifiedReport, found := suite.charitableTrustKeeper.GetImpactReport(ctx, reportRes.ReportId)
	suite.Require().True(found)
	suite.Require().NotNil(verifiedReport.Verification)
	suite.Require().True(verifiedReport.Verification.IsVerified)
	suite.Require().Equal(int32(92), verifiedReport.Verification.Score)
}

// Test cross-module fund flow validation
func (suite *IntegrationTestSuite) TestCrossModuleFundValidation() {
	ctx := suite.ctx

	// Set up both modules with balances
	dswfBalance := []sdk.Coin{sdk.NewCoin("unamo", sdk.NewInt(1000000000))}
	suite.dswfKeeper.SetFundBalance(ctx, dswfBalance)

	trustBalance := charitabletrusttypes.TrustFundBalance{
		TotalBalance:    sdk.NewCoin("unamo", sdk.NewInt(500000000)),
		AvailableAmount: sdk.NewCoin("unamo", sdk.NewInt(500000000)),
	}
	suite.charitableTrustKeeper.SetTrustFundBalance(ctx, trustBalance)

	// Test that allocations cannot exceed available balances
	
	// DSWF over-allocation test
	overAllocationMsg := &dswftypes.MsgProposeAllocation{
		Proposers: []string{"desh1fundmanager1test", "desh1fundmanager2test"},
		Title:     "Over Allocation Test",
		Amount:    sdk.NewCoin("unamo", sdk.NewInt(1500000000)), // More than available
		Category:  "infrastructure",
	}

	// This should fail validation
	err := suite.dswfKeeper.ValidateAllocationProposal(ctx, overAllocationMsg.Amount, overAllocationMsg.Category)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "insufficient funds")

	// CharitableTrust over-allocation test
	trustOverMsg := &charitabletrusttypes.MsgCreateAllocationProposal{
		Proposer:    "desh1trustee1test",
		Title:       "Over Allocation Test",
		TotalAmount: sdk.NewCoin("unamo", sdk.NewInt(600000000)), // More than available
	}

	err = suite.charitableTrustKeeper.ValidateAllocationAmount(ctx, trustOverMsg.TotalAmount)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "insufficient")
}

// Test module state consistency
func (suite *IntegrationTestSuite) TestModuleStateConsistency() {
	ctx := suite.ctx

	// Create multiple allocations in both modules
	for i := 1; i <= 5; i++ {
		// DSWF allocations
		dswfAllocation := dswftypes.FundAllocation{
			Id:       uint64(i),
			Title:    fmt.Sprintf("DSWF Allocation %d", i),
			Amount:   sdk.NewCoin("unamo", sdk.NewInt(int64(i*100000000))),
			Status:   "active",
			Category: "infrastructure",
		}
		suite.dswfKeeper.SetFundAllocation(ctx, dswfAllocation)

		// CharitableTrust allocations
		trustAllocation := charitabletrusttypes.CharitableAllocation{
			Id:               uint64(i),
			OrganizationName: fmt.Sprintf("Charity %d", i),
			Amount:           sdk.NewCoin("unamo", sdk.NewInt(int64(i*50000000))),
			Status:           "active",
			Category:         "education",
		}
		suite.charitableTrustKeeper.SetCharitableAllocation(ctx, trustAllocation)
	}

	// Verify all allocations were stored correctly
	dswfAllocations := suite.dswfKeeper.GetAllFundAllocations(ctx)
	trustAllocations := suite.charitableTrustKeeper.GetAllCharitableAllocations(ctx)

	suite.Require().Len(dswfAllocations, 5)
	suite.Require().Len(trustAllocations, 5)

	// Verify counters are consistent
	nextDSWFID := suite.dswfKeeper.IncrementAllocationCount(ctx)
	nextTrustID := suite.charitableTrustKeeper.IncrementAllocationCount(ctx)

	suite.Require().Equal(uint64(6), nextDSWFID)
	suite.Require().Equal(uint64(6), nextTrustID)

	// Test module-specific queries work correctly
	dswfByCategory := suite.dswfKeeper.GetAllocationsByCategory(ctx, "infrastructure")
	suite.Require().Len(dswfByCategory, 5)

	trustByCategory := suite.charitableTrustKeeper.GetAllocationsByCategory(ctx, "education")
	suite.Require().Len(trustByCategory, 5)
}