package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/deshchain/deshchain/x/dswf/types"
)

// Test FundStatus Query
func (suite *KeeperTestSuite) TestQueryFundStatus() {
	// Setup fund with various allocations
	fundBalance := sdk.NewCoins(sdk.NewCoin("unamo", sdk.NewInt(10000000000))) // 10000 NAMO
	err := suite.app.BankKeeper.MintCoins(suite.ctx, types.ModuleName, fundBalance)
	suite.Require().NoError(err)

	// Create portfolio
	portfolio := types.InvestmentPortfolio{
		TotalValue:       fundBalance[0],
		AllocatedAmount:  sdk.NewCoin("unamo", sdk.NewInt(1500000000)), // 1500 NAMO
		AvailableAmount:  sdk.NewCoin("unamo", sdk.NewInt(8500000000)), // 8500 NAMO
		InvestedAmount:   sdk.NewCoin("unamo", sdk.NewInt(6000000000)), // 6000 NAMO
		TotalReturns:     sdk.NewCoin("unamo", sdk.NewInt(500000000)),  // 500 NAMO
		AnnualReturnRate: sdk.MustNewDecFromStr("0.083"),               // 8.3%
		LastRebalanced:   suite.ctx.BlockTime(),
	}
	suite.keeper.SetInvestmentPortfolio(suite.ctx, portfolio)

	// Create some allocations
	allocations := []types.FundAllocation{
		{Id: 1, Status: "active", Amount: sdk.NewCoin("unamo", sdk.NewInt(500000000))},
		{Id: 2, Status: "active", Amount: sdk.NewCoin("unamo", sdk.NewInt(300000000))},
		{Id: 3, Status: "completed", Amount: sdk.NewCoin("unamo", sdk.NewInt(700000000))},
		{Id: 4, Status: "completed", Amount: sdk.NewCoin("unamo", sdk.NewInt(200000000))},
	}

	for _, alloc := range allocations {
		suite.keeper.SetFundAllocation(suite.ctx, alloc)
	}

	// Query fund status
	req := &types.QueryFundStatusRequest{}
	resp, err := suite.keeper.FundStatus(sdk.WrapSDKContext(suite.ctx), req)
	suite.Require().NoError(err)
	suite.Require().NotNil(resp)

	suite.Require().Equal(fundBalance[0], resp.TotalBalance)
	suite.Require().Equal(portfolio.AllocatedAmount, resp.AllocatedAmount)
	suite.Require().Equal(sdk.NewCoin("unamo", sdk.NewInt(8500000000)), resp.AvailableAmount)
	suite.Require().Equal(portfolio.InvestedAmount, resp.InvestedAmount)
	suite.Require().Equal(portfolio.TotalReturns, resp.TotalReturns)
	suite.Require().Equal("0.083000000000000000", resp.AnnualReturnRate)
	suite.Require().Equal(int32(2), resp.ActiveAllocations)
	suite.Require().Equal(int32(2), resp.CompletedAllocations)
}

// Test Allocation Query
func (suite *KeeperTestSuite) TestQueryAllocation() {
	// Create test allocation
	allocation := types.FundAllocation{
		Id:              1,
		Purpose:         "Test Project",
		Amount:          sdk.NewCoin("unamo", sdk.NewInt(1000000000)),
		Category:        "infrastructure",
		Recipient:       "desh1recipient...",
		Status:          "active",
		ExpectedReturns: sdk.MustNewDecFromStr("1.08"),
		RiskCategory:    "medium",
	}
	suite.keeper.SetFundAllocation(suite.ctx, allocation)

	// Query existing allocation
	req := &types.QueryAllocationRequest{AllocationId: 1}
	resp, err := suite.keeper.Allocation(sdk.WrapSDKContext(suite.ctx), req)
	suite.Require().NoError(err)
	suite.Require().NotNil(resp)
	suite.Require().Equal(allocation, *resp.Allocation)

	// Query non-existent allocation
	reqNotFound := &types.QueryAllocationRequest{AllocationId: 999}
	_, err = suite.keeper.Allocation(sdk.WrapSDKContext(suite.ctx), reqNotFound)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "not found")
}

// Test Allocations Query with Pagination
func (suite *KeeperTestSuite) TestQueryAllocations() {
	// Create multiple allocations
	for i := 1; i <= 10; i++ {
		status := "active"
		if i%3 == 0 {
			status = "completed"
		}
		allocation := types.FundAllocation{
			Id:       uint64(i),
			Purpose:  fmt.Sprintf("Project %d", i),
			Amount:   sdk.NewCoin("unamo", sdk.NewInt(int64(i*100000000))),
			Category: "infrastructure",
			Status:   status,
		}
		suite.keeper.SetFundAllocation(suite.ctx, allocation)
	}

	// Query all allocations with pagination
	req := &types.QueryAllocationsRequest{
		Pagination: &query.PageRequest{
			Limit: 5,
		},
	}
	resp, err := suite.keeper.Allocations(sdk.WrapSDKContext(suite.ctx), req)
	suite.Require().NoError(err)
	suite.Require().NotNil(resp)
	suite.Require().Len(resp.Allocations, 5)
	suite.Require().NotNil(resp.Pagination.NextKey)

	// Query by status
	reqByStatus := &types.QueryAllocationsRequest{
		Status: "active",
		Pagination: &query.PageRequest{
			Limit: 10,
		},
	}
	respByStatus, err := suite.keeper.Allocations(sdk.WrapSDKContext(suite.ctx), reqByStatus)
	suite.Require().NoError(err)
	suite.Require().NotNil(respByStatus)
	// Should have 7 active allocations (1,2,4,5,7,8,10)
	suite.Require().Len(respByStatus.Allocations, 7)
}

// Test Portfolio Query
func (suite *KeeperTestSuite) TestQueryPortfolio() {
	// Create portfolio
	portfolio := types.InvestmentPortfolio{
		TotalValue:      sdk.NewCoin("unamo", sdk.NewInt(10000000000)),
		LiquidAssets:    sdk.NewCoin("unamo", sdk.NewInt(3000000000)),
		InvestedAssets:  sdk.NewCoin("unamo", sdk.NewInt(6000000000)),
		ReservedAssets:  sdk.NewCoin("unamo", sdk.NewInt(1000000000)),
		Components: []types.PortfolioComponent{
			{
				AssetType:    "treasury_bonds",
				Amount:       sdk.NewCoin("unamo", sdk.NewInt(3000000000)),
				CurrentValue: sdk.NewCoin("unamo", sdk.NewInt(3240000000)),
				ReturnRate:   sdk.MustNewDecFromStr("0.08"),
				RiskRating:   "AAA",
			},
			{
				AssetType:    "equity_index",
				Amount:       sdk.NewCoin("unamo", sdk.NewInt(3000000000)),
				CurrentValue: sdk.NewCoin("unamo", sdk.NewInt(3450000000)),
				ReturnRate:   sdk.MustNewDecFromStr("0.15"),
				RiskRating:   "A",
			},
		},
		TotalReturns:     sdk.NewCoin("unamo", sdk.NewInt(690000000)),
		AnnualReturnRate: sdk.MustNewDecFromStr("0.115"),
		RiskScore:        3,
		LastRebalanced:   suite.ctx.BlockTime(),
	}
	suite.keeper.SetInvestmentPortfolio(suite.ctx, portfolio)

	// Query portfolio
	req := &types.QueryPortfolioRequest{}
	resp, err := suite.keeper.Portfolio(sdk.WrapSDKContext(suite.ctx), req)
	suite.Require().NoError(err)
	suite.Require().NotNil(resp)
	suite.Require().Equal(portfolio, *resp.Portfolio)
}

// Test MonthlyReports Query
func (suite *KeeperTestSuite) TestQueryMonthlyReports() {
	// Create monthly reports
	reports := []types.MonthlyReport{
		{
			Period:               "2024-10",
			OpeningBalance:       sdk.NewCoin("unamo", sdk.NewInt(8000000000)),
			ClosingBalance:       sdk.NewCoin("unamo", sdk.NewInt(8500000000)),
			TotalAllocated:       sdk.NewCoin("unamo", sdk.NewInt(1000000000)),
			TotalReturns:         sdk.NewCoin("unamo", sdk.NewInt(100000000)),
			AllocationsApproved:  5,
			AllocationsCompleted: 3,
		},
		{
			Period:               "2024-11",
			OpeningBalance:       sdk.NewCoin("unamo", sdk.NewInt(8500000000)),
			ClosingBalance:       sdk.NewCoin("unamo", sdk.NewInt(9000000000)),
			TotalAllocated:       sdk.NewCoin("unamo", sdk.NewInt(1200000000)),
			TotalReturns:         sdk.NewCoin("unamo", sdk.NewInt(120000000)),
			AllocationsApproved:  4,
			AllocationsCompleted: 2,
		},
		{
			Period:               "2024-12",
			OpeningBalance:       sdk.NewCoin("unamo", sdk.NewInt(9000000000)),
			ClosingBalance:       sdk.NewCoin("unamo", sdk.NewInt(10000000000)),
			TotalAllocated:       sdk.NewCoin("unamo", sdk.NewInt(800000000)),
			TotalReturns:         sdk.NewCoin("unamo", sdk.NewInt(150000000)),
			AllocationsApproved:  3,
			AllocationsCompleted: 4,
		},
	}

	for _, report := range reports {
		suite.keeper.SetMonthlyReport(suite.ctx, report)
	}

	// Query all reports
	req := &types.QueryMonthlyReportsRequest{}
	resp, err := suite.keeper.MonthlyReports(sdk.WrapSDKContext(suite.ctx), req)
	suite.Require().NoError(err)
	suite.Require().NotNil(resp)
	suite.Require().Len(resp.Reports, 3)

	// Query with period filter
	reqFiltered := &types.QueryMonthlyReportsRequest{
		FromPeriod: "2024-11",
		ToPeriod:   "2024-12",
	}
	respFiltered, err := suite.keeper.MonthlyReports(sdk.WrapSDKContext(suite.ctx), reqFiltered)
	suite.Require().NoError(err)
	suite.Require().NotNil(respFiltered)
	suite.Require().Len(respFiltered.Reports, 2)
}

// Test Governance Query
func (suite *KeeperTestSuite) TestQueryGovernance() {
	// Create governance
	governance := types.FundGovernance{
		FundManagers: []types.FundManager{
			{
				Address:     "desh1manager1...",
				Name:        "Fund Manager 1",
				Expertise:   "Fixed Income",
				Performance: sdk.MustNewDecFromStr("1.08"),
			},
			{
				Address:     "desh1manager2...",
				Name:        "Fund Manager 2",
				Expertise:   "Equities",
				Performance: sdk.MustNewDecFromStr("1.12"),
			},
		},
		RequiredSignatures:  2,
		ApprovalThreshold:   sdk.MustNewDecFromStr("0.67"),
		InvestmentCommittee: []string{"desh1committee1...", "desh1committee2..."},
		RiskOfficer:        "desh1risk...",
		ComplianceOfficer:  "desh1compliance...",
		AuditSchedule:      180,
	}
	suite.keeper.SetFundGovernance(suite.ctx, governance)

	// Query governance
	req := &types.QueryGovernanceRequest{}
	resp, err := suite.keeper.Governance(sdk.WrapSDKContext(suite.ctx), req)
	suite.Require().NoError(err)
	suite.Require().NotNil(resp)
	suite.Require().Equal(governance, *resp.Governance)
}

// Test AllocationsByCategory Query
func (suite *KeeperTestSuite) TestQueryAllocationsByCategory() {
	// Create allocations in different categories
	allocations := []types.FundAllocation{
		{
			Id:       1,
			Category: "infrastructure",
			Amount:   sdk.NewCoin("unamo", sdk.NewInt(500000000)),
			Status:   "active",
		},
		{
			Id:       2,
			Category: "infrastructure",
			Amount:   sdk.NewCoin("unamo", sdk.NewInt(300000000)),
			Status:   "completed",
		},
		{
			Id:       3,
			Category: "education",
			Amount:   sdk.NewCoin("unamo", sdk.NewInt(400000000)),
			Status:   "active",
		},
		{
			Id:       4,
			Category: "infrastructure",
			Amount:   sdk.NewCoin("unamo", sdk.NewInt(200000000)),
			Status:   "active",
		},
	}

	for _, alloc := range allocations {
		suite.keeper.SetFundAllocation(suite.ctx, alloc)
	}

	// Query infrastructure allocations
	req := &types.QueryAllocationsByCategoryRequest{
		Category: "infrastructure",
	}
	resp, err := suite.keeper.AllocationsByCategory(sdk.WrapSDKContext(suite.ctx), req)
	suite.Require().NoError(err)
	suite.Require().NotNil(resp)
	suite.Require().Len(resp.Allocations, 3)
	suite.Require().Equal(sdk.NewCoin("unamo", sdk.NewInt(1000000000)), resp.TotalAllocated)

	// Query with status filter
	reqWithStatus := &types.QueryAllocationsByCategoryRequest{
		Category: "infrastructure",
		Status:   "active",
	}
	respWithStatus, err := suite.keeper.AllocationsByCategory(sdk.WrapSDKContext(suite.ctx), reqWithStatus)
	suite.Require().NoError(err)
	suite.Require().NotNil(respWithStatus)
	suite.Require().Len(respWithStatus.Allocations, 2)
	suite.Require().Equal(sdk.NewCoin("unamo", sdk.NewInt(700000000)), respWithStatus.TotalAllocated)
}

// Test PendingDisbursements Query
func (suite *KeeperTestSuite) TestQueryPendingDisbursements() {
	// Create allocations with disbursements
	currentTime := suite.ctx.BlockTime()
	
	allocations := []types.FundAllocation{
		{
			Id:       1,
			Purpose:  "Project A",
			Amount:   sdk.NewCoin("unamo", sdk.NewInt(1000000000)),
			Status:   "active",
			Recipient: "desh1projecta...",
			Disbursements: []types.Disbursement{
				{
					Amount:        sdk.NewCoin("unamo", sdk.NewInt(500000000)),
					ScheduledDate: currentTime.Add(5 * 24 * time.Hour), // 5 days future
					Status:        "pending",
					Milestone:     "Phase 1",
				},
				{
					Amount:        sdk.NewCoin("unamo", sdk.NewInt(500000000)),
					ScheduledDate: currentTime.Add(60 * 24 * time.Hour), // 60 days future
					Status:        "pending",
					Milestone:     "Phase 2",
				},
			},
		},
		{
			Id:       2,
			Purpose:  "Project B",
			Amount:   sdk.NewCoin("unamo", sdk.NewInt(600000000)),
			Status:   "active",
			Recipient: "desh1projectb...",
			Disbursements: []types.Disbursement{
				{
					Amount:        sdk.NewCoin("unamo", sdk.NewInt(300000000)),
					ScheduledDate: currentTime.Add(10 * 24 * time.Hour), // 10 days future
					Status:        "pending",
					Milestone:     "Initial Payment",
				},
				{
					Amount:        sdk.NewCoin("unamo", sdk.NewInt(300000000)),
					ScheduledDate: currentTime.Add(-5 * 24 * time.Hour), // Past date
					Status:        "completed",
					Milestone:     "Completed Payment",
				},
			},
		},
	}

	for _, alloc := range allocations {
		suite.keeper.SetFundAllocation(suite.ctx, alloc)
	}

	// Query pending disbursements (next 30 days)
	req := &types.QueryPendingDisbursementsRequest{
		Pagination: &query.PageRequest{
			Limit: 10,
		},
	}
	resp, err := suite.keeper.PendingDisbursements(sdk.WrapSDKContext(suite.ctx), req)
	suite.Require().NoError(err)
	suite.Require().NotNil(resp)
	suite.Require().Len(resp.Disbursements, 2) // Only 2 disbursements in next 30 days
	suite.Require().Equal(sdk.NewCoin("unamo", sdk.NewInt(800000000)), resp.TotalPending)

	// Verify disbursement details
	for _, disbursement := range resp.Disbursements {
		suite.Require().NotEmpty(disbursement.AllocationPurpose)
		suite.Require().NotEmpty(disbursement.Milestone)
		suite.Require().True(disbursement.Amount.IsPositive())
	}
}

// Test Params Query
func (suite *KeeperTestSuite) TestQueryParams() {
	// Set custom params
	params := types.Params{
		MinFundBalance:        sdk.NewCoin("unamo", sdk.NewInt(1000000000)),
		MaxAllocationPercentage: sdk.MustNewDecFromStr("0.10"),
		MinLiquidityRatio:     sdk.MustNewDecFromStr("0.20"),
		RebalancingFrequency:  90,
		AllocationCategories: []string{
			"infrastructure",
			"education",
			"healthcare",
			"technology",
			"agriculture",
		},
		InvestmentHorizon:     3650, // 10 years
		TargetReturnRate:      sdk.MustNewDecFromStr("0.08"),
		MaxRiskScore:          5,
		DisbursementBatchSize: 10,
		ReportingFrequency:    30,
		AuditRequirement:      true,
	}
	suite.keeper.SetParams(suite.ctx, params)

	// Query params
	req := &types.QueryParamsRequest{}
	resp, err := suite.keeper.Params(sdk.WrapSDKContext(suite.ctx), req)
	suite.Require().NoError(err)
	suite.Require().NotNil(resp)
	suite.Require().Equal(params, resp.Params)
}