package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/deshchain/x/dswf/keeper"
	"github.com/deshchain/deshchain/x/dswf/types"
)

// Test ProposeAllocation
func (suite *KeeperTestSuite) TestMsgProposeAllocation() {
	// Setup governance
	governance := types.FundGovernance{
		FundManagers: []types.FundManager{
			{Address: "desh1manager1..."},
			{Address: "desh1manager2..."},
			{Address: "desh1manager3..."},
		},
		RequiredSignatures: 2,
	}
	suite.keeper.SetFundGovernance(suite.ctx, governance)

	// Setup fund balance
	fundBalance := sdk.NewCoins(sdk.NewCoin("unamo", sdk.NewInt(10000000000)))
	err := suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, fundBalance)
	suite.Require().NoError(err)

	msgServer := keeper.NewMsgServerImpl(suite.keeper)

	// Test valid proposal
	msg := &types.MsgProposeAllocation{
		Proposers:    []string{"desh1manager1...", "desh1manager2..."},
		Purpose:      "Education Infrastructure Development",
		Amount:       sdk.NewCoin("unamo", sdk.NewInt(500000000)), // 500 NAMO
		Category:     "infrastructure",
		Recipient:    "desh1education123...",
		Justification: "Build 10 new schools in rural areas",
		ExpectedImpact: "Provide education access to 5000 students",
		ExpectedReturns: "1.08", // 8% expected return
		RiskAssessment: "Low risk government-backed project",
		Disbursements: []types.DisbursementSchedule{
			{
				Amount:        sdk.NewCoin("unamo", sdk.NewInt(250000000)),
				ScheduledDate: suite.ctx.BlockTime().Add(7 * 24 * time.Hour),
				Milestone:     "Phase 1 - Land acquisition and foundation",
				Conditions:    []string{"Land titles verified", "Environmental clearance obtained"},
			},
			{
				Amount:        sdk.NewCoin("unamo", sdk.NewInt(250000000)),
				ScheduledDate: suite.ctx.BlockTime().Add(60 * 24 * time.Hour),
				Milestone:     "Phase 2 - Construction completion",
				Conditions:    []string{"Phase 1 completed", "Quality inspection passed"},
			},
		},
	}

	resp, err := msgServer.ProposeAllocation(sdk.WrapSDKContext(suite.ctx), msg)
	suite.Require().NoError(err)
	suite.Require().NotNil(resp)
	suite.Require().Greater(resp.AllocationId, uint64(0))

	// Verify allocation created
	allocation, found := suite.keeper.GetFundAllocation(suite.ctx, resp.AllocationId)
	suite.Require().True(found)
	suite.Require().Equal("proposed", allocation.Status)
	suite.Require().Equal(msg.Amount, allocation.Amount)

	// Test insufficient signers
	msgInvalid := &types.MsgProposeAllocation{
		Proposers:    []string{"desh1manager1..."}, // Only 1 signer
		Purpose:      "Test",
		Amount:       sdk.NewCoin("unamo", sdk.NewInt(100000000)),
		Category:     "infrastructure",
		Recipient:    "desh1test...",
		Justification: "Test",
	}

	_, err = msgServer.ProposeAllocation(sdk.WrapSDKContext(suite.ctx), msgInvalid)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "insufficient signatures")
}

// Test ApproveAllocation
func (suite *KeeperTestSuite) TestMsgApproveAllocation() {
	// Setup governance
	governance := types.FundGovernance{
		FundManagers: []types.FundManager{
			{Address: "desh1manager1..."},
			{Address: "desh1manager2..."},
			{Address: "desh1manager3..."},
		},
		RequiredSignatures: 2,
		ApprovalThreshold:  sdk.MustNewDecFromStr("0.67"),
	}
	suite.keeper.SetFundGovernance(suite.ctx, governance)

	// Create a proposed allocation
	allocation := types.FundAllocation{
		Id:       1,
		Purpose:  "Test Allocation",
		Amount:   sdk.NewCoin("unamo", sdk.NewInt(100000000)),
		Category: "infrastructure",
		Status:   "proposed",
	}
	suite.keeper.SetFundAllocation(suite.ctx, allocation)

	msgServer := keeper.NewMsgServerImpl(suite.keeper)

	// Test approval with sufficient signers
	msg := &types.MsgApproveAllocation{
		AllocationId: 1,
		Approvers:    []string{"desh1manager1...", "desh1manager2..."},
		Decision:     "approve",
		Comments:     "Project meets all criteria",
	}

	resp, err := msgServer.ApproveAllocation(sdk.WrapSDKContext(suite.ctx), msg)
	suite.Require().NoError(err)
	suite.Require().True(resp.Success)
	suite.Require().Equal("approved", resp.Status)

	// Verify allocation status updated
	updatedAllocation, found := suite.keeper.GetFundAllocation(suite.ctx, 1)
	suite.Require().True(found)
	suite.Require().Equal("approved", updatedAllocation.Status)

	// Test rejection
	allocation2 := types.FundAllocation{
		Id:       2,
		Purpose:  "Test Allocation 2",
		Amount:   sdk.NewCoin("unamo", sdk.NewInt(200000000)),
		Category: "education",
		Status:   "proposed",
	}
	suite.keeper.SetFundAllocation(suite.ctx, allocation2)

	msgReject := &types.MsgApproveAllocation{
		AllocationId: 2,
		Approvers:    []string{"desh1manager1...", "desh1manager2...", "desh1manager3..."},
		Decision:     "reject",
		Comments:     "Insufficient documentation",
	}

	resp, err = msgServer.ApproveAllocation(sdk.WrapSDKContext(suite.ctx), msgReject)
	suite.Require().NoError(err)
	suite.Require().True(resp.Success)
	suite.Require().Equal("rejected", resp.Status)
}

// Test ExecuteDisbursement
func (suite *KeeperTestSuite) TestMsgExecuteDisbursement() {
	// Setup fund balance
	fundBalance := sdk.NewCoins(sdk.NewCoin("unamo", sdk.NewInt(10000000000)))
	err := suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, fundBalance)
	suite.Require().NoError(err)

	// Create approved allocation with disbursements
	allocation := types.FundAllocation{
		Id:         1,
		Purpose:    "Test Project",
		Amount:     sdk.NewCoin("unamo", sdk.NewInt(1000000000)),
		Category:   "infrastructure",
		Recipient:  "desh1recipient...",
		Status:     "approved",
		Disbursements: []types.Disbursement{
			{
				Amount:        sdk.NewCoin("unamo", sdk.NewInt(500000000)),
				ScheduledDate: suite.ctx.BlockTime().Add(-1 * time.Hour), // Past date
				Status:        "pending",
				Milestone:     "Phase 1",
			},
			{
				Amount:        sdk.NewCoin("unamo", sdk.NewInt(500000000)),
				ScheduledDate: suite.ctx.BlockTime().Add(30 * 24 * time.Hour), // Future date
				Status:        "pending",
				Milestone:     "Phase 2",
			},
		},
	}
	suite.keeper.SetFundAllocation(suite.ctx, allocation)

	msgServer := keeper.NewMsgServerImpl(suite.keeper)

	// Test executing a disbursement
	msg := &types.MsgExecuteDisbursement{
		AllocationId:      1,
		DisbursementIndex: 0,
		Executor:          "desh1executor...",
		VerificationNotes: "Milestone 1 completed successfully",
	}

	resp, err := msgServer.ExecuteDisbursement(sdk.WrapSDKContext(suite.ctx), msg)
	suite.Require().NoError(err)
	suite.Require().True(resp.Success)
	suite.Require().NotEmpty(resp.TxHash)

	// Verify disbursement status
	updatedAllocation, found := suite.keeper.GetFundAllocation(suite.ctx, 1)
	suite.Require().True(found)
	suite.Require().Equal("completed", updatedAllocation.Disbursements[0].Status)
	suite.Require().NotNil(updatedAllocation.Disbursements[0].ExecutedAt)

	// Test executing future disbursement (should fail)
	msgFuture := &types.MsgExecuteDisbursement{
		AllocationId:      1,
		DisbursementIndex: 1,
		Executor:          "desh1executor...",
	}

	_, err = msgServer.ExecuteDisbursement(sdk.WrapSDKContext(suite.ctx), msgFuture)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "not yet due")
}

// Test UpdatePortfolio
func (suite *KeeperTestSuite) TestMsgUpdatePortfolio() {
	// Setup governance with investment committee
	governance := types.FundGovernance{
		InvestmentCommittee: []string{
			"desh1committee1...",
			"desh1committee2...",
			"desh1committee3...",
		},
		RiskOfficer: "desh1risk...",
	}
	suite.keeper.SetFundGovernance(suite.ctx, governance)

	msgServer := keeper.NewMsgServerImpl(suite.keeper)

	// Test portfolio update
	msg := &types.MsgUpdatePortfolio{
		Authority: "desh1committee1...",
		TotalValue: sdk.NewCoin("unamo", sdk.NewInt(10000000000)),
		Components: []types.PortfolioComponent{
			{
				AssetType:    "treasury_bonds",
				Amount:       sdk.NewCoin("unamo", sdk.NewInt(3000000000)),
				CurrentValue: sdk.NewCoin("unamo", sdk.NewInt(3240000000)),
				ReturnRate:   sdk.MustNewDecFromStr("0.08"),
				RiskRating:   "AAA",
				Maturity:     suite.ctx.BlockTime().Add(365 * 24 * time.Hour),
			},
			{
				AssetType:    "equity_index",
				Amount:       sdk.NewCoin("unamo", sdk.NewInt(4000000000)),
				CurrentValue: sdk.NewCoin("unamo", sdk.NewInt(4600000000)),
				ReturnRate:   sdk.MustNewDecFromStr("0.15"),
				RiskRating:   "A",
			},
			{
				AssetType:    "innovation_fund",
				Amount:       sdk.NewCoin("unamo", sdk.NewInt(2000000000)),
				CurrentValue: sdk.NewCoin("unamo", sdk.NewInt(2200000000)),
				ReturnRate:   sdk.MustNewDecFromStr("0.10"),
				RiskRating:   "B",
			},
			{
				AssetType:    "cash_reserve",
				Amount:       sdk.NewCoin("unamo", sdk.NewInt(1000000000)),
				CurrentValue: sdk.NewCoin("unamo", sdk.NewInt(1000000000)),
				ReturnRate:   sdk.MustNewDecFromStr("0.00"),
				RiskRating:   "AAA",
			},
		},
		TotalReturns:     sdk.NewCoin("unamo", sdk.NewInt(1040000000)),
		AnnualReturnRate: "0.104", // 10.4%
		RiskScore:        3,
		RebalanceReason:  "Quarterly rebalancing - shift to growth strategy",
	}

	resp, err := msgServer.UpdatePortfolio(sdk.WrapSDKContext(suite.ctx), msg)
	suite.Require().NoError(err)
	suite.Require().True(resp.Success)

	// Verify portfolio updated
	portfolio, found := suite.keeper.GetInvestmentPortfolio(suite.ctx)
	suite.Require().True(found)
	suite.Require().Equal(msg.TotalValue, portfolio.TotalValue)
	suite.Require().Len(portfolio.Components, 4)

	// Test unauthorized update
	msgUnauth := &types.MsgUpdatePortfolio{
		Authority:  "desh1unauthorized...",
		TotalValue: sdk.NewCoin("unamo", sdk.NewInt(10000000000)),
		Components: []types.PortfolioComponent{},
	}

	_, err = msgServer.UpdatePortfolio(sdk.WrapSDKContext(suite.ctx), msgUnauth)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "unauthorized")
}

// Test SubmitMonthlyReport
func (suite *KeeperTestSuite) TestMsgSubmitMonthlyReport() {
	// Setup governance
	governance := types.FundGovernance{
		ComplianceOfficer: "desh1compliance...",
	}
	suite.keeper.SetFundGovernance(suite.ctx, governance)

	msgServer := keeper.NewMsgServerImpl(suite.keeper)

	// Test report submission
	msg := &types.MsgSubmitMonthlyReport{
		Reporter:        "desh1compliance...",
		Period:          "2025-01",
		OpeningBalance:  sdk.NewCoin("unamo", sdk.NewInt(9000000000)),
		ClosingBalance:  sdk.NewCoin("unamo", sdk.NewInt(10000000000)),
		TotalInflows:    sdk.NewCoin("unamo", sdk.NewInt(1500000000)),
		TotalOutflows:   sdk.NewCoin("unamo", sdk.NewInt(500000000)),
		AllocationsApproved: 3,
		AllocationsCompleted: 2,
		AllocationsActive: 1,
		TotalAllocated: sdk.NewCoin("unamo", sdk.NewInt(1500000000)),
		TotalDisbursed: sdk.NewCoin("unamo", sdk.NewInt(800000000)),
		TotalReturns:   sdk.NewCoin("unamo", sdk.NewInt(120000000)),
		AverageReturnRate: "0.08",
		CategoryBreakdown: []types.CategoryAllocation{
			{
				Category: "infrastructure",
				Amount:   sdk.NewCoin("unamo", sdk.NewInt(800000000)),
				Count:    2,
			},
			{
				Category: "education",
				Amount:   sdk.NewCoin("unamo", sdk.NewInt(500000000)),
				Count:    1,
			},
			{
				Category: "healthcare",
				Amount:   sdk.NewCoin("unamo", sdk.NewInt(200000000)),
				Count:    1,
			},
		},
		RiskMetrics: types.RiskMetrics{
			OverallRiskScore:     3,
			LiquidityRatio:       sdk.MustNewDecFromStr("0.30"),
			ConcentrationRisk:    sdk.MustNewDecFromStr("0.15"),
			DefaultRate:          sdk.MustNewDecFromStr("0.00"),
			RecoveryRate:         sdk.MustNewDecFromStr("1.00"),
			ValueAtRisk:          sdk.NewCoin("unamo", sdk.NewInt(500000000)),
			StressTestResults:    "All scenarios passed",
		},
		Highlights: []string{
			"Successfully completed 2 major infrastructure projects",
			"Portfolio returns exceeded target by 2%",
			"Zero defaults on disbursements",
		},
		Challenges: []string{
			"Delays in one education project due to weather",
			"Need to increase allocation to healthcare sector",
		},
		NextSteps: []string{
			"Review and approve 5 pending proposals",
			"Rebalance portfolio to increase growth allocation",
			"Conduct site visits for active projects",
		},
	}

	resp, err := msgServer.SubmitMonthlyReport(sdk.WrapSDKContext(suite.ctx), msg)
	suite.Require().NoError(err)
	suite.Require().True(resp.Success)

	// Verify report saved
	report, found := suite.keeper.GetMonthlyReport(suite.ctx, "2025-01")
	suite.Require().True(found)
	suite.Require().Equal(msg.Period, report.Period)
	suite.Require().Equal(msg.TotalReturns, report.TotalReturns)

	// Test duplicate report
	_, err = msgServer.SubmitMonthlyReport(sdk.WrapSDKContext(suite.ctx), msg)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "already exists")
}

// Test UpdateGovernance
func (suite *KeeperTestSuite) TestMsgUpdateGovernance() {
	// Set initial governance
	suite.keeper.SetFundGovernance(suite.ctx, types.FundGovernance{
		FundManagers: []types.FundManager{
			{Address: "desh1manager1..."},
		},
		RequiredSignatures: 1,
	})

	msgServer := keeper.NewMsgServerImpl(suite.keeper)

	// Test governance update
	msg := &types.MsgUpdateGovernance{
		Authority: suite.keeper.GetAuthority(),
		NewGovernance: types.FundGovernance{
			FundManagers: []types.FundManager{
				{
					Address:   "desh1manager1...",
					Name:      "Fund Manager 1",
					Expertise: "Fixed Income",
				},
				{
					Address:   "desh1manager2...",
					Name:      "Fund Manager 2",
					Expertise: "Equities",
				},
				{
					Address:   "desh1manager3...",
					Name:      "Fund Manager 3",
					Expertise: "Alternative Investments",
				},
			},
			RequiredSignatures: 2,
			ApprovalThreshold:  sdk.MustNewDecFromStr("0.67"),
			InvestmentCommittee: []string{
				"desh1committee1...",
				"desh1committee2...",
			},
			RiskOfficer:       "desh1risk...",
			ComplianceOfficer: "desh1compliance...",
			AuditSchedule:     180,
		},
	}

	resp, err := msgServer.UpdateGovernance(sdk.WrapSDKContext(suite.ctx), msg)
	suite.Require().NoError(err)
	suite.Require().True(resp.Success)

	// Verify governance updated
	governance, found := suite.keeper.GetFundGovernance(suite.ctx)
	suite.Require().True(found)
	suite.Require().Len(governance.FundManagers, 3)
	suite.Require().Equal(int32(2), governance.RequiredSignatures)

	// Test unauthorized update
	msgUnauth := &types.MsgUpdateGovernance{
		Authority:     "desh1unauthorized...",
		NewGovernance: types.FundGovernance{},
	}

	_, err = msgServer.UpdateGovernance(sdk.WrapSDKContext(suite.ctx), msgUnauth)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "unauthorized")
}

// Test RecordReturns
func (suite *KeeperTestSuite) TestMsgRecordReturns() {
	// Create an active allocation
	allocation := types.FundAllocation{
		Id:       1,
		Purpose:  "Test Investment",
		Amount:   sdk.NewCoin("unamo", sdk.NewInt(1000000000)),
		Category: "infrastructure",
		Status:   "active",
		ExpectedReturns: sdk.MustNewDecFromStr("1.08"),
	}
	suite.keeper.SetFundAllocation(suite.ctx, allocation)

	// Setup fund balance for returns
	returnAmount := sdk.NewCoins(sdk.NewCoin("unamo", sdk.NewInt(80000000)))
	err := suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, returnAmount)
	suite.Require().NoError(err)

	msgServer := keeper.NewMsgServerImpl(suite.keeper)

	// Test recording returns
	msg := &types.MsgRecordReturns{
		AllocationId:   1,
		ActualReturns:  sdk.NewCoin("unamo", sdk.NewInt(80000000)), // 8% return
		Period:         30,
		Reporter:       "desh1reporter...",
		Documentation:  []string{"audit_report.pdf", "bank_statement.pdf"},
	}

	resp, err := msgServer.RecordReturns(sdk.WrapSDKContext(suite.ctx), msg)
	suite.Require().NoError(err)
	suite.Require().True(resp.Success)
	suite.Require().Equal("0.08", resp.ReturnRate)

	// Verify allocation updated
	updatedAllocation, found := suite.keeper.GetFundAllocation(suite.ctx, 1)
	suite.Require().True(found)
	suite.Require().NotNil(updatedAllocation.Returns)
	suite.Require().Equal(msg.ActualReturns, updatedAllocation.Returns.ActualAmount)

	// Test recording returns for non-existent allocation
	msgInvalid := &types.MsgRecordReturns{
		AllocationId:  999,
		ActualReturns: sdk.NewCoin("unamo", sdk.NewInt(100000000)),
		Period:        30,
		Reporter:      "desh1reporter...",
	}

	_, err = msgServer.RecordReturns(sdk.WrapSDKContext(suite.ctx), msgInvalid)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "not found")
}