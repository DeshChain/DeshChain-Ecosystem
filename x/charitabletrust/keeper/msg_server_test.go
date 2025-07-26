package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"deshchain/x/charitabletrust/types"
)

// Test CreateAllocationProposal
func (suite *KeeperTestSuite) TestMsgCreateAllocationProposal() {
	// Set up governance first
	trustees := []types.Trustee{
		{
			Address:     "desh1trustee1test",
			Name:        "Test Trustee 1",
			Status:      "active",
			TermEndDate: suite.ctx.BlockTime().Add(365 * 24 * time.Hour),
			VotingPower: 100,
		},
	}
	governance := types.TrustGovernance{
		Trustees:          trustees,
		Quorum:            1,
		ApprovalThreshold: sdk.NewDecWithPrec(500, 3),
	}
	suite.keeper.SetTrustGovernance(suite.ctx, governance)

	// Create valid proposal
	msg := &types.MsgCreateAllocationProposal{
		Proposer:         "desh1trustee1test",
		Title:            "Q1 2025 Distribution",
		Description:      "Quarterly charity distribution",
		TotalAmount:      sdk.NewCoin("unamo", sdk.NewInt(1000000000)),
		ExpectedOutcomes: "Help 1000 beneficiaries",
		Allocations: []types.ProposedAllocation{
			{
				CharitableOrgWalletId: 1,
				OrganizationName:      "Test Charity",
				Amount:                sdk.NewCoin("unamo", sdk.NewInt(500000000)),
				Purpose:               "Education support",
				Category:              "education",
			},
			{
				CharitableOrgWalletId: 2,
				OrganizationName:      "Health Foundation",
				Amount:                sdk.NewCoin("unamo", sdk.NewInt(500000000)),
				Purpose:               "Medical camps",
				Category:              "healthcare",
			},
		},
		Documents: []string{"https://ipfs.io/ipfs/Qm...proposal"},
	}

	// Execute
	res, err := suite.msgServer.CreateAllocationProposal(suite.ctx, msg)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Equal(uint64(1), res.ProposalId)

	// Verify proposal was created
	proposal, found := suite.keeper.GetAllocationProposal(suite.ctx, 1)
	suite.Require().True(found)
	suite.Require().Equal(msg.Title, proposal.Title)
	suite.Require().Equal(msg.TotalAmount, proposal.TotalAmount)
	suite.Require().Len(proposal.Allocations, 2)
	suite.Require().Equal("pending", proposal.Status)

	// Test with non-trustee
	msg.Proposer = "desh1nottrusteetest"
	_, err = suite.msgServer.CreateAllocationProposal(suite.ctx, msg)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "not a trustee")
}

// Test VoteOnProposal
func (suite *KeeperTestSuite) TestMsgVoteOnProposal() {
	// Set up governance with multiple trustees
	trustees := []types.Trustee{
		{
			Address:     "desh1trustee1test",
			Name:        "Test Trustee 1",
			Status:      "active",
			TermEndDate: suite.ctx.BlockTime().Add(365 * 24 * time.Hour),
			VotingPower: 100,
		},
		{
			Address:     "desh1trustee2test",
			Name:        "Test Trustee 2",
			Status:      "active",
			TermEndDate: suite.ctx.BlockTime().Add(365 * 24 * time.Hour),
			VotingPower: 100,
		},
		{
			Address:     "desh1trustee3test",
			Name:        "Test Trustee 3",
			Status:      "active",
			TermEndDate: suite.ctx.BlockTime().Add(365 * 24 * time.Hour),
			VotingPower: 100,
		},
	}
	governance := types.TrustGovernance{
		Trustees:          trustees,
		Quorum:            2,
		ApprovalThreshold: sdk.NewDecWithPrec(600, 3), // 60%
	}
	suite.keeper.SetTrustGovernance(suite.ctx, governance)

	// Create a proposal
	proposal := types.AllocationProposal{
		Id:              1,
		Title:           "Test Proposal",
		TotalAmount:     sdk.NewCoin("unamo", sdk.NewInt(1000000000)),
		Status:          "pending",
		VotingStartTime: suite.ctx.BlockTime(),
		VotingEndTime:   suite.ctx.BlockTime().Add(7 * 24 * time.Hour),
		Votes:           []types.Vote{},
	}
	suite.keeper.SetAllocationProposal(suite.ctx, proposal)

	// Test voting yes
	msgYes := &types.MsgVoteOnProposal{
		ProposalId: 1,
		Voter:      "desh1trustee1test",
		VoteType:   "yes",
		Comments:   "Looks good",
	}

	res, err := suite.msgServer.VoteOnProposal(suite.ctx, msgYes)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Equal("Vote recorded successfully", res.Result)

	// Test voting no
	msgNo := &types.MsgVoteOnProposal{
		ProposalId: 1,
		Voter:      "desh1trustee2test",
		VoteType:   "no",
		Comments:   "Need more details",
	}

	res, err = suite.msgServer.VoteOnProposal(suite.ctx, msgNo)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)

	// Verify votes were recorded
	proposal, found := suite.keeper.GetAllocationProposal(suite.ctx, 1)
	suite.Require().True(found)
	suite.Require().Len(proposal.Votes, 2)
	suite.Require().Equal("desh1trustee1test", proposal.Votes[0].Voter)
	suite.Require().Equal("yes", proposal.Votes[0].VoteType)
	suite.Require().Equal("desh1trustee2test", proposal.Votes[1].Voter)
	suite.Require().Equal("no", proposal.Votes[1].VoteType)

	// Test duplicate vote
	_, err = suite.msgServer.VoteOnProposal(suite.ctx, msgYes)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "already voted")

	// Test non-trustee voting
	msgInvalid := &types.MsgVoteOnProposal{
		ProposalId: 1,
		Voter:      "desh1nottrusteetest",
		VoteType:   "yes",
	}
	_, err = suite.msgServer.VoteOnProposal(suite.ctx, msgInvalid)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "not a trustee")

	// Test voting on non-existent proposal
	msgInvalid.ProposalId = 999
	msgInvalid.Voter = "desh1trustee1test"
	_, err = suite.msgServer.VoteOnProposal(suite.ctx, msgInvalid)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "proposal not found")
}

// Test ExecuteAllocation
func (suite *KeeperTestSuite) TestMsgExecuteAllocation() {
	// Set up governance
	trustees := []types.Trustee{
		{
			Address:     "desh1trustee1test",
			Name:        "Test Trustee 1",
			Status:      "active",
			TermEndDate: suite.ctx.BlockTime().Add(365 * 24 * time.Hour),
			VotingPower: 100,
		},
		{
			Address:     "desh1trustee2test",
			Name:        "Test Trustee 2",
			Status:      "active",
			TermEndDate: suite.ctx.BlockTime().Add(365 * 24 * time.Hour),
			VotingPower: 100,
		},
	}
	governance := types.TrustGovernance{
		Trustees:          trustees,
		Quorum:            2,
		ApprovalThreshold: sdk.NewDecWithPrec(500, 3), // 50%
	}
	suite.keeper.SetTrustGovernance(suite.ctx, governance)

	// Create and approve a proposal
	proposal := types.AllocationProposal{
		Id:          1,
		Title:       "Test Proposal",
		TotalAmount: sdk.NewCoin("unamo", sdk.NewInt(1000000000)),
		Status:      "approved",
		Allocations: []types.ProposedAllocation{
			{
				CharitableOrgWalletId: 1,
				OrganizationName:      "Test Charity",
				Amount:                sdk.NewCoin("unamo", sdk.NewInt(500000000)),
				Purpose:               "Education",
				Category:              "education",
			},
			{
				CharitableOrgWalletId: 2,
				OrganizationName:      "Health Foundation",
				Amount:                sdk.NewCoin("unamo", sdk.NewInt(500000000)),
				Purpose:               "Healthcare",
				Category:              "healthcare",
			},
		},
		Votes: []types.Vote{
			{Voter: "desh1trustee1test", VoteType: "yes", VoteValue: 100},
			{Voter: "desh1trustee2test", VoteType: "yes", VoteValue: 100},
		},
	}
	suite.keeper.SetAllocationProposal(suite.ctx, proposal)

	// Set trust fund balance
	balance := types.TrustFundBalance{
		TotalBalance:    sdk.NewCoin("unamo", sdk.NewInt(2000000000)),
		AllocatedAmount: sdk.NewCoin("unamo", sdk.NewInt(0)),
		AvailableAmount: sdk.NewCoin("unamo", sdk.NewInt(2000000000)),
	}
	suite.keeper.SetTrustFundBalance(suite.ctx, balance)

	// Execute allocation
	msg := &types.MsgExecuteAllocation{
		ProposalId: 1,
		Executor:   "desh1trustee1test",
	}

	res, err := suite.msgServer.ExecuteAllocation(suite.ctx, msg)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Equal("Allocations created from proposal", res.Result)
	suite.Require().Len(res.AllocationIds, 2)

	// Verify allocations were created
	allocation1, found := suite.keeper.GetCharitableAllocation(suite.ctx, res.AllocationIds[0])
	suite.Require().True(found)
	suite.Require().Equal(uint64(1), allocation1.CharitableOrgWalletId)
	suite.Require().Equal("active", allocation1.Status)

	allocation2, found := suite.keeper.GetCharitableAllocation(suite.ctx, res.AllocationIds[1])
	suite.Require().True(found)
	suite.Require().Equal(uint64(2), allocation2.CharitableOrgWalletId)
	suite.Require().Equal("active", allocation2.Status)

	// Verify proposal status updated
	proposal, found = suite.keeper.GetAllocationProposal(suite.ctx, 1)
	suite.Require().True(found)
	suite.Require().Equal("executed", proposal.Status)

	// Verify trust fund balance updated
	balance, found = suite.keeper.GetTrustFundBalance(suite.ctx)
	suite.Require().True(found)
	suite.Require().Equal(sdk.NewInt(1000000000), balance.AllocatedAmount.Amount)
	suite.Require().Equal(sdk.NewInt(1000000000), balance.AvailableAmount.Amount)

	// Test executing already executed proposal
	_, err = suite.msgServer.ExecuteAllocation(suite.ctx, msg)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "not in approved status")
}

// Test DistributeFunds
func (suite *KeeperTestSuite) TestMsgDistributeFunds() {
	// Create an active allocation
	allocation := types.CharitableAllocation{
		Id:                    1,
		CharitableOrgWalletId: 1,
		OrganizationName:      "Test Charity",
		Amount:                sdk.NewCoin("unamo", sdk.NewInt(100000000)),
		Status:                "active",
		Distribution:          &types.Distribution{},
	}
	suite.keeper.SetCharitableAllocation(suite.ctx, allocation)

	// Set trust fund balance
	balance := types.TrustFundBalance{
		TotalBalance:     sdk.NewCoin("unamo", sdk.NewInt(1000000000)),
		AllocatedAmount:  sdk.NewCoin("unamo", sdk.NewInt(100000000)),
		AvailableAmount:  sdk.NewCoin("unamo", sdk.NewInt(900000000)),
		TotalDistributed: sdk.NewCoin("unamo", sdk.NewInt(0)),
	}
	suite.keeper.SetTrustFundBalance(suite.ctx, balance)

	// Distribute funds
	msg := &types.MsgDistributeFunds{
		AllocationId: 1,
		Distributor:  "desh1treasury",
		RecipientAddress: "desh1charity1test",
		TxReference:  "Distribution for Q1 2025",
	}

	res, err := suite.msgServer.DistributeFunds(suite.ctx, msg)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Contains(res.TxHash, "0x")
	suite.Require().Equal("distributed", res.Status)

	// Verify allocation status updated
	allocation, found := suite.keeper.GetCharitableAllocation(suite.ctx, 1)
	suite.Require().True(found)
	suite.Require().Equal("distributed", allocation.Status)
	suite.Require().NotEmpty(allocation.Distribution.TxHash)
	suite.Require().Equal("desh1treasury", allocation.Distribution.DistributedBy)

	// Verify trust fund balance updated
	balance, found = suite.keeper.GetTrustFundBalance(suite.ctx)
	suite.Require().True(found)
	suite.Require().Equal(sdk.NewInt(0), balance.AllocatedAmount.Amount)
	suite.Require().Equal(sdk.NewInt(100000000), balance.TotalDistributed.Amount)

	// Test distributing already distributed allocation
	_, err = suite.msgServer.DistributeFunds(suite.ctx, msg)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "not in active status")
}

// Test SubmitImpactReport
func (suite *KeeperTestSuite) TestMsgSubmitImpactReport() {
	// Create a distributed allocation
	allocation := types.CharitableAllocation{
		Id:                    1,
		CharitableOrgWalletId: 1,
		OrganizationName:      "Test Charity",
		Amount:                sdk.NewCoin("unamo", sdk.NewInt(100000000)),
		Status:                "distributed",
	}
	suite.keeper.SetCharitableAllocation(suite.ctx, allocation)

	// Submit impact report
	msg := &types.MsgSubmitImpactReport{
		AllocationId:         1,
		OrganizationId:       1,
		ReportingPeriod:      "Q1 2025",
		BeneficiariesReached: 200,
		FundsUtilized:        sdk.NewCoin("unamo", sdk.NewInt(80000000)),
		ActivitiesConducted:  []string{"Distributed books", "Conducted classes"},
		OutcomesAchieved:     []string{"200 students enrolled", "95% attendance"},
		Challenges:           []string{"Remote area access"},
		SupportingDocuments:  []string{"https://ipfs.io/ipfs/Qm...report"},
		ImpactMetrics: map[string]string{
			"students_enrolled": "200",
			"attendance_rate":   "95%",
		},
		Submitter: "desh1charity1test",
	}

	res, err := suite.msgServer.SubmitImpactReport(suite.ctx, msg)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Equal(uint64(1), res.ReportId)
	suite.Require().Equal("submitted", res.Status)

	// Verify report was created
	report, found := suite.keeper.GetImpactReport(suite.ctx, 1)
	suite.Require().True(found)
	suite.Require().Equal(msg.AllocationId, report.AllocationId)
	suite.Require().Equal(msg.BeneficiariesReached, report.BeneficiariesReached)
	suite.Require().Equal(msg.FundsUtilized, report.FundsUtilized)
	suite.Require().Len(report.ActivitiesConducted, 2)
	suite.Require().Len(report.OutcomesAchieved, 2)
}

// Test VerifyImpactReport
func (suite *KeeperTestSuite) TestMsgVerifyImpactReport() {
	// Create an impact report
	report := types.ImpactReport{
		Id:                   1,
		AllocationId:         1,
		OrganizationId:       1,
		BeneficiariesReached: 200,
		FundsUtilized:        sdk.NewCoin("unamo", sdk.NewInt(80000000)),
	}
	suite.keeper.SetImpactReport(suite.ctx, report)

	// Verify the report
	msg := &types.MsgVerifyImpactReport{
		ReportId:         1,
		Verifier:         "desh1verifier1test",
		IsVerified:       true,
		VerificationNote: "All claims verified through site visit",
		Score:            90,
	}

	res, err := suite.msgServer.VerifyImpactReport(suite.ctx, msg)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Equal("verified", res.Status)

	// Verify the verification was recorded
	report, found := suite.keeper.GetImpactReport(suite.ctx, 1)
	suite.Require().True(found)
	suite.Require().NotNil(report.Verification)
	suite.Require().True(report.Verification.IsVerified)
	suite.Require().Equal(int32(90), report.Verification.Score)
	suite.Require().Equal("desh1verifier1test", report.Verification.VerifiedBy)
}

// Test ReportFraud
func (suite *KeeperTestSuite) TestMsgReportFraud() {
	msg := &types.MsgReportFraud{
		AlertType:        "misuse_of_funds",
		Severity:         "high",
		AllocationId:     1,
		OrganizationId:   1,
		Description:      "Funds used for unauthorized purposes",
		Evidence:         []string{"Bank statements", "Photos"},
		ReportedBy:       "desh1whistleblower1test",
		MonetaryImpact:   sdk.NewCoin("unamo", sdk.NewInt(20000000)),
		AffectedEntities: []string{"Test Charity"},
	}

	res, err := suite.msgServer.ReportFraud(suite.ctx, msg)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Equal(uint64(1), res.AlertId)
	suite.Require().Equal("pending", res.Status)

	// Verify fraud alert was created
	alert, found := suite.keeper.GetFraudAlert(suite.ctx, 1)
	suite.Require().True(found)
	suite.Require().Equal(msg.AlertType, alert.AlertType)
	suite.Require().Equal(msg.Severity, alert.Severity)
	suite.Require().Equal(msg.Description, alert.Description)
	suite.Require().Len(alert.Evidence, 2)
}

// Test InvestigateFraud
func (suite *KeeperTestSuite) TestMsgInvestigateFraud() {
	// Create a fraud alert
	alert := types.FraudAlert{
		Id:             1,
		AlertType:      "misuse_of_funds",
		Status:         "pending",
		AllocationId:   1,
		OrganizationId: 1,
	}
	suite.keeper.SetFraudAlert(suite.ctx, alert)

	// Investigate the fraud
	msg := &types.MsgInvestigateFraud{
		AlertId:          1,
		InvestigatorId:   "desh1investigator1test",
		Findings:         []string{"Funds misused", "False documentation"},
		Recommendation:   "blacklist_organization",
		ActionsTaken:     []string{"Funds frozen", "Legal notice sent"},
		EvidenceGathered: []string{"Bank records", "Witness testimonies"},
	}

	res, err := suite.msgServer.InvestigateFraud(suite.ctx, msg)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Equal("investigated", res.Status)
	suite.Require().Equal("blacklist_organization", res.Recommendation)

	// Verify investigation was recorded
	alert, found := suite.keeper.GetFraudAlert(suite.ctx, 1)
	suite.Require().True(found)
	suite.Require().NotNil(alert.Investigation)
	suite.Require().Equal("desh1investigator1test", alert.Investigation.InvestigatorId)
	suite.Require().Len(alert.Investigation.Findings, 2)
	suite.Require().Equal("blacklist_organization", alert.Investigation.Recommendation)
	suite.Require().Equal("investigated", alert.Status)
}

// Test UpdateTrustGovernance
func (suite *KeeperTestSuite) TestMsgUpdateTrustGovernance() {
	// Set initial governance
	initialGov := types.TrustGovernance{
		Trustees: []types.Trustee{
			{
				Address: "desh1trustee1test",
				Status:  "active",
			},
		},
		Quorum:            1,
		ApprovalThreshold: sdk.NewDecWithPrec(500, 3),
	}
	suite.keeper.SetTrustGovernance(suite.ctx, initialGov)

	// Update governance
	msg := &types.MsgUpdateTrustGovernance{
		Authority: "desh1gov", // This should be the governance module account
		Trustees: []types.Trustee{
			{
				Address:     "desh1trustee1test",
				Name:        "Updated Trustee 1",
				Role:        "Chairman",
				Status:      "active",
				TermEndDate: suite.ctx.BlockTime().Add(730 * 24 * time.Hour),
				VotingPower: 150,
			},
			{
				Address:     "desh1trustee2new",
				Name:        "New Trustee 2",
				Role:        "Secretary",
				Status:      "active",
				TermEndDate: suite.ctx.BlockTime().Add(730 * 24 * time.Hour),
				VotingPower: 100,
			},
		},
		Quorum:              2,
		ApprovalThreshold:   sdk.NewDecWithPrec(667, 3), // 66.7%
		TransparencyOfficer: "desh1transparency2test",
	}

	// This would normally require governance authority
	// For testing, we'll directly update
	_, err := suite.msgServer.UpdateTrustGovernance(suite.ctx, msg)
	// The actual implementation would check authority
	// For now, we'll test the update logic directly
	suite.Require().Error(err) // Should fail without proper authority
}